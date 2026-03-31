package main

// llvm_functional.go - LLVM IR generation for functional programming primitives
//
// Implements:
//   Option<T>  methods: .is_some(), .is_none(), .unwrap(), .unwrap_or(default), .map(f), .and_then(f)
//   Result<T>  methods: .is_ok(),   .is_err(),  .unwrap(), .unwrap_or(default), .map(f), .map_err(f), .and_then(f)
//
// Both types share the same struct layout: { i1 flag, T value }
//   Option: flag = has_value  (1 = Some, 0 = None)
//   Result: flag = is_ok      (1 = Ok,   0 = Err)

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

// ============================================================================
// Result<T> constructor codegen
// ============================================================================

// generateResultExpression generates code for Ok(value) or Err(value).
// Struct layout: { i1 is_ok, T value }
func (cg *LLVMCodeGenerator) generateResultExpression(r *ResultExpression) (llvm.Value, error) {
	value, err := cg.generateExpression(r.Value)
	if err != nil {
		return llvm.Value{}, err
	}

	valueType := value.Type()
	resultType := cg.context.StructType([]llvm.Type{
		cg.context.Int1Type(),
		valueType,
	}, false)

	flag := llvm.ConstInt(cg.context.Int1Type(), 0, false)
	if r.IsOk {
		flag = llvm.ConstInt(cg.context.Int1Type(), 1, false)
	}

	s := llvm.ConstNull(resultType)
	s = cg.builder.CreateInsertValue(s, flag, 0, "set_flag")
	s = cg.builder.CreateInsertValue(s, value, 1, "set_value")
	return s, nil
}

// ============================================================================
// Shared Option/Result type detection
// ============================================================================

// isOptionalStructType returns true if t is a struct whose first field is i1.
// This matches the { i1 flag, T value } layout shared by Option<T> and Result<T>.
func isOptionalStructType(t llvm.Type) bool {
	if t.TypeKind() != llvm.StructTypeKind {
		return false
	}
	fields := t.StructElementTypes()
	return len(fields) == 2 &&
		fields[0].TypeKind() == llvm.IntegerTypeKind &&
		fields[0].IntTypeWidth() == 1
}

// isOptionalMethodName returns true for method names that are Option/Result combinators.
func isOptionalMethodName(name string) bool {
	switch name {
	case "is_some", "is_none", "is_ok", "is_err",
		"unwrap", "unwrap_or",
		"map", "map_err", "and_then":
		return true
	}
	return false
}

// ============================================================================
// Option/Result method dispatch
// ============================================================================

// generateOptionalMethodCall handles method calls on Option and Result values.
func (cg *LLVMCodeGenerator) generateOptionalMethodCall(m *MethodCall, obj llvm.Value) (llvm.Value, error) {
	switch m.MethodName {
	case "is_some", "is_ok":
		return cg.builder.CreateExtractValue(obj, 0, "flag"), nil

	case "is_none", "is_err":
		flag := cg.builder.CreateExtractValue(obj, 0, "flag")
		return cg.builder.CreateNot(flag, "notflag"), nil

	case "unwrap":
		return cg.builder.CreateExtractValue(obj, 1, "unwrapped"), nil

	case "unwrap_or":
		if len(m.Args) != 1 {
			return llvm.Value{}, fmt.Errorf("unwrap_or requires exactly 1 argument")
		}
		defaultVal, err := cg.generateExpression(m.Args[0])
		if err != nil {
			return llvm.Value{}, err
		}
		flag := cg.builder.CreateExtractValue(obj, 0, "flag")
		inner := cg.builder.CreateExtractValue(obj, 1, "inner")
		if inner.Type() != defaultVal.Type() {
			defaultVal = cg.coerceToType(defaultVal, inner.Type())
		}
		return cg.builder.CreateSelect(flag, inner, defaultVal, "unwrap_or"), nil

	case "map":
		if len(m.Args) != 1 {
			return llvm.Value{}, fmt.Errorf("map requires exactly 1 argument (a function)")
		}
		fn, fnType, err := cg.resolveFunctionArg(m.Args[0])
		if err != nil {
			return llvm.Value{}, err
		}
		return cg.generateOptionalMap(obj, fn, fnType, true)

	case "map_err":
		if len(m.Args) != 1 {
			return llvm.Value{}, fmt.Errorf("map_err requires exactly 1 argument (a function)")
		}
		fn, fnType, err := cg.resolveFunctionArg(m.Args[0])
		if err != nil {
			return llvm.Value{}, err
		}
		return cg.generateOptionalMap(obj, fn, fnType, false)

	case "and_then":
		if len(m.Args) != 1 {
			return llvm.Value{}, fmt.Errorf("and_then requires exactly 1 argument (a function)")
		}
		fn, fnType, err := cg.resolveFunctionArg(m.Args[0])
		if err != nil {
			return llvm.Value{}, err
		}
		return cg.generateOptionalAndThen(obj, fn, fnType)

	default:
		return llvm.Value{}, fmt.Errorf("unknown optional method: %s", m.MethodName)
	}
}

// generateOptionalMap generates LLVM IR for .map(f) and .map_err(f).
//
//	mapTrueBranch=true  → transform when flag==1 (Some / Ok)
//	mapTrueBranch=false → transform when flag==0 (None / Err) — used by map_err
func (cg *LLVMCodeGenerator) generateOptionalMap(opt llvm.Value, fn llvm.Value, fnType llvm.Type, mapTrueBranch bool) (llvm.Value, error) {
	currentFn := cg.currentFn
	mapBB := llvm.AddBasicBlock(currentFn, "opt.map.do")
	skipBB := llvm.AddBasicBlock(currentFn, "opt.map.skip")
	mergeBB := llvm.AddBasicBlock(currentFn, "opt.map.merge")

	flag := cg.builder.CreateExtractValue(opt, 0, "flag")
	inner := cg.builder.CreateExtractValue(opt, 1, "inner")

	var condVal llvm.Value
	if mapTrueBranch {
		condVal = flag
	} else {
		condVal = cg.builder.CreateNot(flag, "notflag")
	}
	cg.builder.CreateCondBr(condVal, mapBB, skipBB)

	// map branch: call f(inner), build new { i1, mappedType } struct
	cg.builder.SetInsertPointAtEnd(mapBB)
	mapped := cg.builder.CreateCall(fnType, fn, []llvm.Value{inner}, "mapped")
	mappedType := fnType.ReturnType()
	newStructType := cg.context.StructType([]llvm.Type{cg.context.Int1Type(), mappedType}, false)
	mapFlag := llvm.ConstInt(cg.context.Int1Type(), 1, false)
	if !mapTrueBranch {
		mapFlag = llvm.ConstInt(cg.context.Int1Type(), 0, false)
	}
	newStructMap := llvm.ConstNull(newStructType)
	newStructMap = cg.builder.CreateInsertValue(newStructMap, mapFlag, 0, "map_flag")
	newStructMap = cg.builder.CreateInsertValue(newStructMap, mapped, 1, "map_val")
	cg.builder.CreateBr(mergeBB)
	mapEndBB := cg.builder.GetInsertBlock()

	// skip branch: propagate None/Err with zeroed value field
	cg.builder.SetInsertPointAtEnd(skipBB)
	keepStructType := cg.context.StructType([]llvm.Type{cg.context.Int1Type(), mappedType}, false)
	keepStruct := llvm.ConstNull(keepStructType)
	var skipFlag llvm.Value
	if mapTrueBranch {
		skipFlag = llvm.ConstInt(cg.context.Int1Type(), 0, false) // propagate None/Err
	} else {
		skipFlag = llvm.ConstInt(cg.context.Int1Type(), 1, false) // propagate Some/Ok
	}
	keepStruct = cg.builder.CreateInsertValue(keepStruct, skipFlag, 0, "skip_flag")
	cg.builder.CreateBr(mergeBB)
	skipEndBB := cg.builder.GetInsertBlock()

	// merge
	cg.builder.SetInsertPointAtEnd(mergeBB)
	phi := cg.builder.CreatePHI(newStructType, "opt.map.result")
	phi.AddIncoming([]llvm.Value{newStructMap, keepStruct}, []llvm.BasicBlock{mapEndBB, skipEndBB})
	return phi, nil
}

// generateOptionalAndThen generates LLVM IR for .and_then(f) (monadic bind).
//
// Lotus functions cannot declare an Option/Result return type, so f always
// returns a plain value.  We wrap it the same way .map() does so that the
// result remains an { i1 flag, T value } struct that the combinator chain
// can continue to operate on.
//
//	if flag==1: result = { 1, f(inner) }   (Some/Ok with transformed value)
//	else:       result = { 0, zero }        (None/Err propagated)
func (cg *LLVMCodeGenerator) generateOptionalAndThen(opt llvm.Value, fn llvm.Value, fnType llvm.Type) (llvm.Value, error) {
	return cg.generateOptionalMap(opt, fn, fnType, true)
}

// ============================================================================
// Helper: resolve a function argument to (llvm.Value fn, llvm.Type fnType)
// Handles: Identifier, FunctionReference, LambdaExpression
// ============================================================================

func (cg *LLVMCodeGenerator) resolveFunctionArg(arg ASTNode) (llvm.Value, llvm.Type, error) {
	switch a := arg.(type) {
	case *Identifier:
		fn, ok := cg.functions[a.Name]
		if !ok {
			return llvm.Value{}, llvm.Type{}, fmt.Errorf("undefined function '%s' passed as combinator argument", a.Name)
		}
		return fn, fn.GlobalValueType(), nil

	case *FunctionReference:
		fn, ok := cg.functions[a.FunctionName]
		if !ok {
			return llvm.Value{}, llvm.Type{}, fmt.Errorf("undefined function '%s' passed as combinator argument", a.FunctionName)
		}
		return fn, fn.GlobalValueType(), nil

	case *LambdaExpression:
		// Compile the lambda; then look it up by its auto-generated name.
		if _, err := cg.generateLambdaExpression(a); err != nil {
			return llvm.Value{}, llvm.Type{}, err
		}
		lambdaName := fmt.Sprintf("__lambda_%d", cg.stringCounter-1)
		fn, ok := cg.functions[lambdaName]
		if !ok {
			return llvm.Value{}, llvm.Type{}, fmt.Errorf("lambda combinator argument requires a named function; use 'fn name' syntax")
		}
		return fn, fn.GlobalValueType(), nil

	default:
		return llvm.Value{}, llvm.Type{}, fmt.Errorf("combinator argument must be a function name or reference, got %T", arg)
	}
}

// ============================================================================
// func module: compose, identity, flip, const_fn, apply/ap, mempty, mappend
// ============================================================================

// generateFuncIdentity generates LLVM IR for func::identity(x) → x.
func (cg *LLVMCodeGenerator) generateFuncIdentity(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("identity requires exactly 1 argument")
	}
	return cg.generateExpression(call.Args[0])
}

// generateFuncCompose generates LLVM IR for func::compose(f, g, x) → g(f(x)).
func (cg *LLVMCodeGenerator) generateFuncCompose(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, fmt.Errorf("compose requires exactly 3 arguments (f, g, x)")
	}
	f, fType, err := cg.resolveFunctionArg(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	g, gType, err := cg.resolveFunctionArg(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	x, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}
	intermediate := cg.builder.CreateCall(fType, f, []llvm.Value{x}, "compose_f")
	return cg.builder.CreateCall(gType, g, []llvm.Value{intermediate}, "compose_g"), nil
}

// generateFuncFlip generates LLVM IR for func::flip(f, a, b) → f(b, a).
func (cg *LLVMCodeGenerator) generateFuncFlip(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, fmt.Errorf("flip requires exactly 3 arguments (f, a, b)")
	}
	f, fType, err := cg.resolveFunctionArg(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	a, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	b, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}
	return cg.builder.CreateCall(fType, f, []llvm.Value{b, a}, "flip_result"), nil
}

// generateFuncConstFn generates LLVM IR for func::const_fn(x, _) → x.
// Both arguments are evaluated (for side effects); only x is returned.
func (cg *LLVMCodeGenerator) generateFuncConstFn(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("const_fn requires exactly 2 arguments")
	}
	x, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	if _, err = cg.generateExpression(call.Args[1]); err != nil {
		return llvm.Value{}, err
	}
	return x, nil
}

// generateFuncApply generates LLVM IR for func::apply(f, x) -> f(x).
func (cg *LLVMCodeGenerator) generateFuncApply(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("apply requires exactly 2 arguments (f, x)")
	}
	f, fType, err := cg.resolveFunctionArg(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	x, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	return cg.builder.CreateCall(fType, f, []llvm.Value{x}, "apply_result"), nil
}

// generateFuncMempty generates LLVM IR for func::mempty() -> 0.
func (cg *LLVMCodeGenerator) generateFuncMempty(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 0 {
		return llvm.Value{}, fmt.Errorf("mempty requires exactly 0 arguments")
	}
	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

// generateFuncMappend generates LLVM IR for func::mappend(a, b) -> a + b.
func (cg *LLVMCodeGenerator) generateFuncMappend(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("mappend requires exactly 2 arguments")
	}
	a, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	b, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	if a.Type().TypeKind() != llvm.IntegerTypeKind || b.Type().TypeKind() != llvm.IntegerTypeKind {
		return llvm.Value{}, fmt.Errorf("mappend currently supports integer operands only")
	}
	// Normalize to i64 so both int32/int64 inputs work predictably.
	if a.Type() != cg.context.Int64Type() {
		a = cg.builder.CreateSExt(a, cg.context.Int64Type(), "mappend_a_i64")
	}
	if b.Type() != cg.context.Int64Type() {
		b = cg.builder.CreateSExt(b, cg.context.Int64Type(), "mappend_b_i64")
	}
	return cg.builder.CreateAdd(a, b, "mappend_result"), nil
}

// generateFuncNegate generates LLVM IR for func::negate(f, x) → !f(x).
// Returns 1 (i64) if f(x)==0, 0 otherwise.
func (cg *LLVMCodeGenerator) generateFuncNegate(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("negate requires exactly 2 arguments (f, x)")
	}
	f, fType, err := cg.resolveFunctionArg(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	x, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	result := cg.builder.CreateCall(fType, f, []llvm.Value{x}, "negate_call")
	zero := llvm.ConstInt(result.Type(), 0, false)
	isZero := cg.builder.CreateICmp(llvm.IntEQ, result, zero, "negate_eq0")
	return cg.builder.CreateZExt(isZero, cg.context.Int64Type(), "negate_result"), nil
}

// generateFuncBoth generates LLVM IR for func::both(f, g, x) → f(x) != 0 && g(x) != 0.
// Returns 1 (i64) if both f(x) and g(x) are non-zero, 0 otherwise.
func (cg *LLVMCodeGenerator) generateFuncBoth(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, fmt.Errorf("both requires exactly 3 arguments (f, g, x)")
	}
	f, fType, err := cg.resolveFunctionArg(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	g, gType, err := cg.resolveFunctionArg(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	x, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}
	fResult := cg.builder.CreateCall(fType, f, []llvm.Value{x}, "both_f")
	gResult := cg.builder.CreateCall(gType, g, []llvm.Value{x}, "both_g")
	fZero := llvm.ConstInt(fResult.Type(), 0, false)
	gZero := llvm.ConstInt(gResult.Type(), 0, false)
	fBool := cg.builder.CreateICmp(llvm.IntNE, fResult, fZero, "both_f_bool")
	gBool := cg.builder.CreateICmp(llvm.IntNE, gResult, gZero, "both_g_bool")
	andResult := cg.builder.CreateAnd(fBool, gBool, "both_and")
	return cg.builder.CreateZExt(andResult, cg.context.Int64Type(), "both_result"), nil
}

// generateFuncEitherFn generates LLVM IR for func::either_fn(f, g, x) → f(x) != 0 || g(x) != 0.
// Returns 1 (i64) if either f(x) or g(x) is non-zero, 0 otherwise.
func (cg *LLVMCodeGenerator) generateFuncEitherFn(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, fmt.Errorf("either_fn requires exactly 3 arguments (f, g, x)")
	}
	f, fType, err := cg.resolveFunctionArg(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	g, gType, err := cg.resolveFunctionArg(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	x, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}
	fResult := cg.builder.CreateCall(fType, f, []llvm.Value{x}, "either_f")
	gResult := cg.builder.CreateCall(gType, g, []llvm.Value{x}, "either_g")
	fZero := llvm.ConstInt(fResult.Type(), 0, false)
	gZero := llvm.ConstInt(gResult.Type(), 0, false)
	fBool := cg.builder.CreateICmp(llvm.IntNE, fResult, fZero, "either_f_bool")
	gBool := cg.builder.CreateICmp(llvm.IntNE, gResult, gZero, "either_g_bool")
	orResult := cg.builder.CreateOr(fBool, gBool, "either_or")
	return cg.builder.CreateZExt(orResult, cg.context.Int64Type(), "either_result"), nil
}

// generateFuncOn generates LLVM IR for func::on(f, g, a, b) → f(g(a), g(b)).
// Applies g to both inputs, then combines with binary function f.
func (cg *LLVMCodeGenerator) generateFuncOn(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 4 {
		return llvm.Value{}, fmt.Errorf("on requires exactly 4 arguments (f, g, a, b)")
	}
	f, fType, err := cg.resolveFunctionArg(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	g, gType, err := cg.resolveFunctionArg(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	a, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}
	b, err := cg.generateExpression(call.Args[3])
	if err != nil {
		return llvm.Value{}, err
	}
	ga := cg.builder.CreateCall(gType, g, []llvm.Value{a}, "on_ga")
	gb := cg.builder.CreateCall(gType, g, []llvm.Value{b}, "on_gb")
	return cg.builder.CreateCall(fType, f, []llvm.Value{ga, gb}, "on_result"), nil
}

// generateFuncConverge generates LLVM IR for func::converge(f, g, h, x) → f(g(x), h(x)).
// Applies g and h to the same input, then merges results with binary function f.
func (cg *LLVMCodeGenerator) generateFuncConverge(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 4 {
		return llvm.Value{}, fmt.Errorf("converge requires exactly 4 arguments (f, g, h, x)")
	}
	f, fType, err := cg.resolveFunctionArg(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	g, gType, err := cg.resolveFunctionArg(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	h, hType, err := cg.resolveFunctionArg(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}
	x, err := cg.generateExpression(call.Args[3])
	if err != nil {
		return llvm.Value{}, err
	}
	gx := cg.builder.CreateCall(gType, g, []llvm.Value{x}, "conv_gx")
	hx := cg.builder.CreateCall(hType, h, []llvm.Value{x}, "conv_hx")
	return cg.builder.CreateCall(fType, f, []llvm.Value{gx, hx}, "conv_result"), nil
}

// generateFuncClamp generates LLVM IR for func::clamp(lo, hi, x) → max(lo, min(hi, x)).
// Clamps integer x to the inclusive range [lo, hi]. All values normalized to i64.
func (cg *LLVMCodeGenerator) generateFuncClamp(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, fmt.Errorf("clamp requires exactly 3 arguments (lo, hi, x)")
	}
	lo, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	hi, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	x, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}
	i64 := cg.context.Int64Type()
	if lo.Type() != i64 {
		lo = cg.builder.CreateSExt(lo, i64, "lo_i64")
	}
	if hi.Type() != i64 {
		hi = cg.builder.CreateSExt(hi, i64, "hi_i64")
	}
	if x.Type() != i64 {
		x = cg.builder.CreateSExt(x, i64, "x_i64")
	}
	// max(lo, x): if lo > x then clampedLo = lo, else clampedLo = x
	loGtX := cg.builder.CreateICmp(llvm.IntSGT, lo, x, "lo_gt_x")
	clampedLo := cg.builder.CreateSelect(loGtX, lo, x, "max_lo_x")
	// min(hi, clampedLo): if hi < clampedLo then result = hi, else result = clampedLo
	hiLtClamped := cg.builder.CreateICmp(llvm.IntSLT, hi, clampedLo, "hi_lt_clamped")
	return cg.builder.CreateSelect(hiLtClamped, hi, clampedLo, "clamp_result"), nil
}
