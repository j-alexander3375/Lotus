package main

// llvm_runtime.go - support runtime bindings.
//
// A number of documented stdlib functions (P1-1: hashmaps, sorted
// collections, IPv6/DNS, file stat/seek, time, hashing, JSON, str::join) have
// no in-tree LLVM codegen. The old x86-asm backend implemented some of these
// with hand-written assembly; that approach doesn't fit the LLVM IR path.
// Since compiled Lotus binaries already link against libc via clang
// (compiler.go), the pragmatic route recommended by the audit is to link a
// small C support library and declare its functions as externs, which is
// what this file wires up.
//
// The C source lives in lotus_runtime.c, embedded into the compiler binary
// so it needs no separate install step; buildBinaryWithLLVM writes it to a
// temp file and passes it to clang alongside the generated IR.

import (
	_ "embed"
	"fmt"
	"strings"

	"tinygo.org/x/go-llvm"
)

//go:embed runtimec/lotus_runtime.c
var lotusRuntimeC string

// runtimeArgKind describes how a Lotus-side argument/return value should be
// marshalled across the C ABI boundary, where every runtime function takes
// and returns either int64_t or char*.
type runtimeArgKind int

const (
	rtInt runtimeArgKind = iota // passed/returned as i64 (also used for heap "pointers as i64")
	rtStr                       // passed/returned as i8* (C string)
)

// declareRuntimeFn declares an extern C function from lotus_runtime.c with
// the given argument kinds and return kind, if not already declared.
func (cg *LLVMCodeGenerator) declareRuntimeFn(cName string, argKinds []runtimeArgKind, retKind runtimeArgKind) llvm.Value {
	if fn, ok := cg.functions[cName]; ok {
		return fn
	}
	paramTypes := make([]llvm.Type, len(argKinds))
	for i, k := range argKinds {
		paramTypes[i] = cg.runtimeLLVMType(k)
	}
	fnType := llvm.FunctionType(cg.runtimeLLVMType(retKind), paramTypes, false)
	fn := llvm.AddFunction(cg.module, cName, fnType)
	fn.SetLinkage(llvm.ExternalLinkage)
	cg.functions[cName] = fn
	return fn
}

func (cg *LLVMCodeGenerator) runtimeLLVMType(k runtimeArgKind) llvm.Type {
	if k == rtStr {
		return llvm.PointerType(cg.context.Int8Type(), 0)
	}
	return cg.context.Int64Type()
}

// generateRuntimeCall generates a call into the embedded C runtime, coercing
// each Lotus argument to the declared kind (ptrtoint for a string argument
// expected as i64, ptr passthrough for a string argument expected as i8*,
// etc.) and coercing the result back to i64 if the C function returns char*.
func (cg *LLVMCodeGenerator) generateRuntimeCall(call *FunctionCall, cName string, argKinds []runtimeArgKind, retKind runtimeArgKind) (llvm.Value, error) {
	if len(call.Args) != len(argKinds) {
		return llvm.Value{}, fmt.Errorf("%s expects %d argument(s), got %d", call.Name, len(argKinds), len(call.Args))
	}

	fn := cg.declareRuntimeFn(cName, argKinds, retKind)

	args := make([]llvm.Value, len(call.Args))
	for i, argNode := range call.Args {
		val, err := cg.generateExpression(argNode)
		if err != nil {
			return llvm.Value{}, err
		}
		args[i] = cg.coerceRuntimeArg(val, argKinds[i])
	}

	// No return-value conversion needed: a C char* and a Lotus string are
	// both i8* in LLVM IR (see generateConcat et al.), and rtInt results are
	// already i64 on both sides (opaque heap handles, matching the
	// convention array/hashmap/json codegen uses elsewhere in this file).
	return cg.builder.CreateCall(fn.GlobalValueType(), fn, args, cName+"tmp"), nil
}

// generateFprintf implements io::fprintf(fp, fmt, args...): writes formatted
// output to the FILE* fp (represented as i64, matching the os::popen
// convention elsewhere in this file) using C's variadic fprintf. As with
// printf, the format argument should be a literal - see generatePrintf.
func (cg *LLVMCodeGenerator) generateFprintf(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) < 2 {
		return llvm.Value{}, fmt.Errorf("fprintf requires at least 2 arguments (fp, format)")
	}
	if _, ok := cg.functions["fprintf"]; !ok {
		fnType := llvm.FunctionType(
			cg.context.Int32Type(),
			[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0), llvm.PointerType(cg.context.Int8Type(), 0)},
			true,
		)
		f := llvm.AddFunction(cg.module, "fprintf", fnType)
		f.SetLinkage(llvm.ExternalLinkage)
		cg.functions["fprintf"] = f
	}

	fpVal, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	fp := cg.builder.CreateIntToPtr(fpVal, llvm.PointerType(cg.context.Int8Type(), 0), "fprintf_fp")

	args := []llvm.Value{fp}
	for _, argNode := range call.Args[1:] {
		val, err := cg.generateExpression(argNode)
		if err != nil {
			return llvm.Value{}, err
		}
		if val.Type().TypeKind() == llvm.FloatTypeKind {
			val = cg.builder.CreateFPExt(val, cg.context.DoubleType(), "fpext")
		}
		args = append(args, val)
	}

	fn := cg.functions["fprintf"]
	return cg.builder.CreateCall(fn.GlobalValueType(), fn, args, "fprintftmp"), nil
}

// generateSprint implements io::sprint/io::sprintln: format every argument
// like print/println (space-separated, per-type default format) but return
// the result as a freshly allocated string instead of writing it to stdout.
func (cg *LLVMCodeGenerator) generateSprint(call *FunctionCall, newline bool) (llvm.Value, error) {
	formatParts := make([]string, 0, len(call.Args))
	callArgs := make([]llvm.Value, 0, len(call.Args)+1)
	for _, argNode := range call.Args {
		arg, err := cg.generateExpression(argNode)
		if err != nil {
			return llvm.Value{}, err
		}
		switch arg.Type().TypeKind() {
		case llvm.PointerTypeKind:
			formatParts = append(formatParts, "%s")
		case llvm.DoubleTypeKind, llvm.FloatTypeKind:
			if arg.Type().TypeKind() == llvm.FloatTypeKind {
				arg = cg.builder.CreateFPExt(arg, cg.context.DoubleType(), "fpext")
			}
			formatParts = append(formatParts, "%f")
		default:
			if arg.Type().TypeKind() == llvm.IntegerTypeKind {
				arg = cg.coerceToType(arg, cg.context.Int64Type())
			}
			formatParts = append(formatParts, "%ld")
		}
		callArgs = append(callArgs, arg)
	}
	format := strings.Join(formatParts, " ")
	if newline {
		format += "\n"
	}

	// lotus_asprintf is a true C varargs function; declare it directly with
	// variadic=true rather than through declareRuntimeFn's fixed-arity path.
	asprintfName := "lotus_asprintf"
	if _, ok := cg.functions[asprintfName+"_variadic"]; !ok {
		fnType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0)},
			true,
		)
		f := llvm.AddFunction(cg.module, asprintfName, fnType)
		f.SetLinkage(llvm.ExternalLinkage)
		cg.functions[asprintfName+"_variadic"] = f
	}
	fn := cg.functions[asprintfName+"_variadic"]

	args := append([]llvm.Value{cg.createGlobalString(format)}, callArgs...)
	return cg.builder.CreateCall(fn.GlobalValueType(), fn, args, "sprinttmp"), nil
}

// generateMemSizeof implements mem::sizeof(x) -> byte size of x's LLVM type.
// Unlike the other stdlib functions this needs compile-time type
// information, not a runtime value, so it is resolved from the AST rather
// than dispatched into the C runtime.
func (cg *LLVMCodeGenerator) generateMemSizeof(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("sizeof requires exactly 1 argument")
	}
	val, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	bits := cg.getTypeSizeInBits(val.Type())
	if bits == 0 {
		bits = 64 // unknown aggregate/pointer-like type: assume a pointer width
	}
	return llvm.ConstInt(cg.context.Int64Type(), uint64(bits/8), false), nil
}

// coerceRuntimeArg adapts a generated Lotus value to the C parameter kind.
func (cg *LLVMCodeGenerator) coerceRuntimeArg(val llvm.Value, kind runtimeArgKind) llvm.Value {
	wantPtr := kind == rtStr
	isPtr := val.Type().TypeKind() == llvm.PointerTypeKind
	switch {
	case wantPtr && isPtr:
		return val
	case wantPtr && !isPtr:
		return cg.builder.CreateIntToPtr(val, llvm.PointerType(cg.context.Int8Type(), 0), "argptr")
	case !wantPtr && isPtr:
		return cg.builder.CreatePtrToInt(val, cg.context.Int64Type(), "argint")
	default: // !wantPtr && !isPtr
		return cg.coerceToType(val, cg.context.Int64Type())
	}
}
