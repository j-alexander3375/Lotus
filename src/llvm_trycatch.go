package main

// llvm_trycatch.go - try/catch/finally/throw codegen.
//
// Implemented with setjmp/longjmp (see runtimec/lotus_runtime.c) rather than
// real LLVM exception handling (landingpad/invoke/a personality routine):
// Lotus has no stable "exception object" type to hand a personality
// function, and setjmp/longjmp gets the practically important behavior - a
// throw anywhere on the call stack, including deep inside unrelated
// function calls, transfers control straight back to the nearest enclosing
// try - without needing every intervening function to cooperate.
//
// Two independent mechanisms are involved, and it's worth being explicit
// about the boundary between them:
//
//  1. `throw` propagation across function calls is handled ENTIRELY at
//     runtime via the jmp_buf stack in lotus_runtime.c. A function that
//     throws needs no cooperation from its callers, and a function that
//     calls something that might throw needs no cooperation either - the
//     longjmp physically restores the enclosing try's stack pointer,
//     skipping over any number of intervening call frames for free. No
//     compile-time bookkeeping is needed for this case.
//  2. `ret`/`break`/`continue` written DIRECTLY inside a try/catch/finally
//     body in the SAME function do NOT go through longjmp at all - they're
//     plain LLVM branches/returns. For these, `finally` blocks are tracked
//     with the compile-time cg.finallyStack and inlined at each early-exit
//     point before the actual branch/return (see generateReturn/
//     generateBreak/generateContinue in llvm_codegen.go).
//
// Exception values: `throw expr` evaluates expr and marshals it into a
// single i64 payload slot (pointers via ptrtoint, floats via bitcast, ints
// via the usual width coercion). A catch clause's declared type, if any
// (`catch (string e) { ... }`), is a static reinterpretation of that
// payload - like the language's existing `bitcast<Type>(...)` - not a
// runtime type check: Lotus has no runtime type tags on values, so it
// cannot actually verify the thrown value matches. This is also why only
// one catch clause per try is accepted (see parseTryStatement) - a second
// clause could never be genuinely selected over the first.

import (
	"tinygo.org/x/go-llvm"
)

// declareTryRuntime declares the externs used by try/catch/finally/throw,
// if not already declared.
func (cg *LLVMCodeGenerator) declareTryRuntime() {
	i64 := cg.context.Int64Type()
	i32 := cg.context.Int32Type()
	voidTy := cg.context.VoidType()
	i8ptr := llvm.PointerType(cg.context.Int8Type(), 0)

	if _, ok := cg.functions["setjmp"]; !ok {
		fnType := llvm.FunctionType(i32, []llvm.Type{i8ptr}, false)
		fn := llvm.AddFunction(cg.module, "setjmp", fnType)
		fn.SetLinkage(llvm.ExternalLinkage)
		// setjmp returns more than once (once normally, once per longjmp
		// that targets it) - without this attribute, an optimizing build
		// (anything above -O0) is entitled to assume it returns once and
		// cache/reorder values computed before the call incorrectly across
		// it, corrupting locals read after control resumes via longjmp.
		kindID := llvm.AttributeKindID("returns_twice")
		fn.AddFunctionAttr(cg.context.CreateEnumAttribute(kindID, 0))
		cg.functions["setjmp"] = fn
	}
	if _, ok := cg.functions["lotus_try_push"]; !ok {
		fn := llvm.AddFunction(cg.module, "lotus_try_push", llvm.FunctionType(i8ptr, nil, false))
		fn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["lotus_try_push"] = fn
	}
	if _, ok := cg.functions["lotus_try_pop"]; !ok {
		fn := llvm.AddFunction(cg.module, "lotus_try_pop", llvm.FunctionType(voidTy, nil, false))
		fn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["lotus_try_pop"] = fn
	}
	if _, ok := cg.functions["lotus_get_exception"]; !ok {
		fn := llvm.AddFunction(cg.module, "lotus_get_exception", llvm.FunctionType(i64, nil, false))
		fn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["lotus_get_exception"] = fn
	}
	if _, ok := cg.functions["lotus_throw"]; !ok {
		fn := llvm.AddFunction(cg.module, "lotus_throw", llvm.FunctionType(voidTy, []llvm.Type{i64}, false))
		fn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["lotus_throw"] = fn
	}
}

// exceptionValueToI64 marshals a thrown value into the single i64 payload
// slot lotus_throw carries. Pointers (strings, heap handles) round-trip via
// ptrtoint/inttoptr elsewhere in this codegen already; floats need a
// bit-preserving reinterpretation (not a numeric coerceToType conversion,
// which would truncate/round) since the payload slot is an integer.
func (cg *LLVMCodeGenerator) exceptionValueToI64(val llvm.Value) llvm.Value {
	switch val.Type().TypeKind() {
	case llvm.PointerTypeKind:
		return cg.builder.CreatePtrToInt(val, cg.context.Int64Type(), "excptr")
	case llvm.DoubleTypeKind:
		return cg.builder.CreateBitCast(val, cg.context.Int64Type(), "excbits")
	case llvm.FloatTypeKind:
		d := cg.builder.CreateFPExt(val, cg.context.DoubleType(), "excfpext")
		return cg.builder.CreateBitCast(d, cg.context.Int64Type(), "excbits")
	default:
		return cg.coerceToType(val, cg.context.Int64Type())
	}
}

// exceptionI64FromType reinterprets a raw i64 exception payload as the
// catch clause's declared type (see the file doc comment on why this is a
// static cast, not a runtime-checked type match). typeName is one of the
// canonical names catchExceptionTypeName produces ("string", "bool",
// "float", "char"), or "" / "int" for the default int64 interpretation.
func (cg *LLVMCodeGenerator) exceptionI64FromType(raw llvm.Value, typeName string) llvm.Value {
	switch typeName {
	case "string":
		return cg.builder.CreateIntToPtr(raw, llvm.PointerType(cg.context.Int8Type(), 0), "excstr")
	case "bool":
		return cg.builder.CreateTrunc(raw, cg.context.Int1Type(), "excbool")
	case "float":
		return cg.builder.CreateBitCast(raw, cg.context.DoubleType(), "excfloat")
	case "char":
		return cg.builder.CreateTrunc(raw, cg.context.Int32Type(), "excchar")
	default:
		return raw
	}
}

// generateThrowStatement implements `throw expr;`.
func (cg *LLVMCodeGenerator) generateThrowStatement(t *ThrowStatement) error {
	cg.declareTryRuntime()
	val, err := cg.generateExpression(t.Exception)
	if err != nil {
		return err
	}
	payload := cg.exceptionValueToI64(val)
	throwFn := cg.functions["lotus_throw"]
	cg.builder.CreateCall(throwFn.GlobalValueType(), throwFn, []llvm.Value{payload}, "")
	// lotus_throw never returns (it either longjmps to an enclosing try or
	// calls exit(1)) - mark the block accordingly rather than leaving it
	// without a terminator.
	cg.builder.CreateUnreachable()
	return nil
}

// rethrowCurrentException re-raises whatever exception is currently being
// handled and terminates the current block. Used when a finally-only try
// (no catch clause) finishes running its finally block after an exception,
// and must propagate it to the next enclosing try.
func (cg *LLVMCodeGenerator) rethrowCurrentException() {
	excFn := cg.functions["lotus_get_exception"]
	exc := cg.builder.CreateCall(excFn.GlobalValueType(), excFn, nil, "exc")
	throwFn := cg.functions["lotus_throw"]
	cg.builder.CreateCall(throwFn.GlobalValueType(), throwFn, []llvm.Value{exc}, "")
	cg.builder.CreateUnreachable()
}

// generateTryStatement implements try/catch/finally. See the file doc
// comment for the overall design.
func (cg *LLVMCodeGenerator) generateTryStatement(t *TryStatement) error {
	cg.declareTryRuntime()
	fn := cg.currentFn

	tryBlock := llvm.AddBasicBlock(fn, "try.body")
	dispatchBlock := llvm.AddBasicBlock(fn, "try.dispatch")
	afterBlock := llvm.AddBasicBlock(fn, "try.after")

	hasCatch := len(t.CatchClauses) == 1
	hasFinally := len(t.FinallyBlock) > 0

	var catchBlock, finallyBlock llvm.BasicBlock
	if hasCatch {
		catchBlock = llvm.AddBasicBlock(fn, "try.catch")
	}
	if hasFinally {
		finallyBlock = llvm.AddBasicBlock(fn, "try.finally")
	}

	// Tracks which blocks branch into finallyBlock, and whether that path
	// must re-throw once the finally body completes (only the "exception
	// occurred and no catch clause claimed it" path does).
	type finallyEdge struct {
		block   llvm.BasicBlock
		rethrow llvm.Value // i1 constant
	}
	var finallyEdges []finallyEdge
	falseV := llvm.ConstInt(cg.context.Int1Type(), 0, false)
	trueV := llvm.ConstInt(cg.context.Int1Type(), 1, false)

	// Push a runtime frame and call setjmp directly here - the C standard
	// requires setjmp's result to be inspected in the same function
	// invocation (and while it's still on the stack) that called it, so
	// this can't be hidden behind a Go-runtime wrapper that returns before
	// inspecting the result.
	pushFn := cg.functions["lotus_try_push"]
	buf := cg.builder.CreateCall(pushFn.GlobalValueType(), pushFn, nil, "trybuf")
	setjmpFn := cg.functions["setjmp"]
	r := cg.builder.CreateCall(setjmpFn.GlobalValueType(), setjmpFn, []llvm.Value{buf}, "setjmpres")
	wasThrown := cg.builder.CreateICmp(llvm.IntNE, r, llvm.ConstInt(cg.context.Int32Type(), 0, false), "wasthrown")
	cg.builder.CreateCondBr(wasThrown, dispatchBlock, tryBlock)

	// Register this try's finally frame for the DURATION of the try+catch
	// body generation, so a ret/break/continue written directly inside
	// either one unwinds through it (see generateReturn/generateBreak/
	// generateContinue). Popped again below, right before generating the
	// finally block "for real" - but runPendingFinally (triggered by such
	// an early exit) may ALREADY have popped it as part of its own unwind,
	// so the identity of the pushed frame is tracked here and the pop
	// below only fires if it's actually still on top (a plain unconditional
	// pop would double-pop and panic on an unrelated frame/empty stack).
	var myFrame *finallyFrame
	if hasFinally {
		myFrame = &finallyFrame{
			body:           t.FinallyBlock,
			loopExitAtPush: cg.loopExitBlock,
		}
		cg.finallyStack = append(cg.finallyStack, myFrame)
	}

	// --- try body (normal, non-thrown path) ---
	cg.builder.SetInsertPointAtEnd(tryBlock)
	for _, stmt := range t.TryBlock {
		if err := cg.generateStatement(stmt); err != nil {
			return err
		}
		if cg.blockTerminated() {
			break
		}
	}
	if !cg.blockTerminated() {
		popFn := cg.functions["lotus_try_pop"]
		cg.builder.CreateCall(popFn.GlobalValueType(), popFn, nil, "")
		exitBlock := cg.builder.GetInsertBlock()
		if hasFinally {
			cg.builder.CreateBr(finallyBlock)
			finallyEdges = append(finallyEdges, finallyEdge{exitBlock, falseV})
		} else {
			cg.builder.CreateBr(afterBlock)
		}
	}

	// --- dispatch: control resumes here via longjmp ---
	cg.builder.SetInsertPointAtEnd(dispatchBlock)
	switch {
	case hasCatch:
		cg.builder.CreateBr(catchBlock)
	case hasFinally:
		exitBlock := cg.builder.GetInsertBlock()
		cg.builder.CreateBr(finallyBlock)
		finallyEdges = append(finallyEdges, finallyEdge{exitBlock, trueV})
	default:
		// parseTryStatement requires at least one of catch/finally, so this
		// is unreachable in practice; handle defensively.
		cg.rethrowCurrentException()
	}

	// --- catch body ---
	if hasCatch {
		cg.builder.SetInsertPointAtEnd(catchBlock)
		clause := t.CatchClauses[0]
		excFn := cg.functions["lotus_get_exception"]
		rawExc := cg.builder.CreateCall(excFn.GlobalValueType(), excFn, nil, "exc")

		var shadowed LLVMVariable
		var hadShadowed bool
		if clause.ExceptionVar != "" {
			typedExc := cg.exceptionI64FromType(rawExc, clause.ExceptionType)
			shadowed, hadShadowed = cg.namedValues[clause.ExceptionVar]
			alloca := cg.builder.CreateAlloca(typedExc.Type(), clause.ExceptionVar)
			cg.builder.CreateStore(typedExc, alloca)
			cg.namedValues[clause.ExceptionVar] = LLVMVariable{Alloca: alloca, ElementType: typedExc.Type()}
		}

		for _, stmt := range clause.Body {
			if err := cg.generateStatement(stmt); err != nil {
				return err
			}
			if cg.blockTerminated() {
				break
			}
		}

		if clause.ExceptionVar != "" {
			if hadShadowed {
				cg.namedValues[clause.ExceptionVar] = shadowed
			} else {
				delete(cg.namedValues, clause.ExceptionVar)
			}
		}

		if !cg.blockTerminated() {
			exitBlock := cg.builder.GetInsertBlock()
			if hasFinally {
				cg.builder.CreateBr(finallyBlock)
				finallyEdges = append(finallyEdges, finallyEdge{exitBlock, falseV})
			} else {
				cg.builder.CreateBr(afterBlock)
			}
		}
	}

	// --- finally body ---
	if hasFinally {
		// No longer "pending" - we're generating it for real now, so a
		// ret/break/continue inside the finally body itself must NOT try to
		// re-run it (see generateReturn/generateBreak/generateContinue).
		// Only pop if our frame is still on top: an early ret/break/continue
		// inside the try/catch body above may have already consumed it via
		// runPendingFinally's own unwind.
		if n := len(cg.finallyStack); n > 0 && cg.finallyStack[n-1] == myFrame {
			cg.finallyStack = cg.finallyStack[:n-1]
		}

		cg.builder.SetInsertPointAtEnd(finallyBlock)

		if len(finallyEdges) == 0 {
			// Both try and catch always exited early (e.g. every path
			// ends in `ret`) - this block is unreachable via normal flow.
			cg.builder.CreateUnreachable()
		} else {
			var rethrowFlag llvm.Value
			if len(finallyEdges) == 1 {
				rethrowFlag = finallyEdges[0].rethrow
			} else {
				phi := cg.builder.CreatePHI(cg.context.Int1Type(), "try.rethrow")
				blocks := make([]llvm.BasicBlock, len(finallyEdges))
				values := make([]llvm.Value, len(finallyEdges))
				for i, e := range finallyEdges {
					blocks[i] = e.block
					values[i] = e.rethrow
				}
				phi.AddIncoming(values, blocks)
				rethrowFlag = phi
			}

			for _, stmt := range t.FinallyBlock {
				if err := cg.generateStatement(stmt); err != nil {
					return err
				}
				if cg.blockTerminated() {
					break
				}
			}

			if !cg.blockTerminated() {
				rethrowBlock := llvm.AddBasicBlock(fn, "try.rethrow")
				cg.builder.CreateCondBr(rethrowFlag, rethrowBlock, afterBlock)
				cg.builder.SetInsertPointAtEnd(rethrowBlock)
				cg.rethrowCurrentException()
			}
		}
	}

	cg.builder.SetInsertPointAtEnd(afterBlock)
	return nil
}

// runPendingFinally replays the bodies of finally frames pending in
// finallyStack, from innermost to outermost, at the current insert point -
// used by generateReturn/generateBreak/generateContinue for an early exit
// written directly inside a try/catch/finally. Each matching frame is
// popped from the stack before its body is generated, so if that body
// itself contains a ret/break/continue, the resulting recursive call only
// sees the remaining outer frames (no frame runs twice), and so a `ret`
// inside a `finally` correctly does not try to re-run its own enclosing
// finally.
//
// match reports whether a given frame should run for this particular exit:
// generateReturn passes a match that accepts every frame (return exits the
// whole function, past every enclosing finally); generateBreak/
// generateContinue pass one that only accepts frames opened during the
// loop actually being exited.
func (cg *LLVMCodeGenerator) runPendingFinally(match func(*finallyFrame) bool) error {
	for len(cg.finallyStack) > 0 {
		top := cg.finallyStack[len(cg.finallyStack)-1]
		if !match(top) {
			return nil
		}
		cg.finallyStack = cg.finallyStack[:len(cg.finallyStack)-1]
		for _, stmt := range top.body {
			if err := cg.generateStatement(stmt); err != nil {
				return err
			}
			if cg.blockTerminated() {
				// The finally body itself exited early (its own ret/break/
				// continue/throw) - that supersedes the original exit this
				// unwind was for; don't run anything further.
				return nil
			}
		}
	}
	return nil
}
