package main

// llvm_format.go - Enhanced printf formatting for LLVM backend
//
// This file implements advanced printf formatting capabilities including:
// - Width specifiers: %5d, %10s
// - Padding flags: %05d (zero-pad), %-10s (left-align)
// - Precision: %.3f, %.5s
// - Sign flags: %+d (always show sign), % d (space for positive)
// - Alternative forms: %#x (0x prefix), %#o (0 prefix)
//
// Format syntax: %[flags][width][.precision][length]specifier

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

// FormatFlags represents printf format flags
type FormatFlags struct {
	LeftAlign bool // '-' flag
	ShowSign  bool // '+' flag
	SpaceSign bool // ' ' flag
	AltForm   bool // '#' flag
	ZeroPad   bool // '0' flag
	Width     int  // minimum field width
	Precision int  // precision for floats/strings
	HasWidth  bool
	HasPrec   bool
}

// ============================================================================
// ENHANCED PRINTF GENERATION
// ============================================================================

// generateFormattedPrintf generates printf with enhanced format specifiers
func (cg *LLVMCodeGenerator) generateFormattedPrintf(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) == 0 {
		return llvm.Value{}, nil
	}

	// Get the format string
	formatArg := call.Args[0]
	formatStr, isString := formatArg.(*StringLiteral)
	if !isString {
		// Fall back to regular printf if format isn't a literal
		return cg.generatePrintf(call)
	}

	// Parse the format string and generate appropriate calls
	return cg.generateFormattedOutput(formatStr.Value, call.Args[1:])
}

// generateFormattedOutput processes a format string with arguments
func (cg *LLVMCodeGenerator) generateFormattedOutput(format string, args []ASTNode) (llvm.Value, error) {
	// For now, use C's printf which handles all standard format specifiers
	// The format string is already in the correct format

	allArgs := make([]llvm.Value, 0, len(args)+1)

	// Add format string
	fmtStr := cg.createGlobalString(format)
	allArgs = append(allArgs, fmtStr)

	// Add remaining arguments
	for _, arg := range args {
		val, err := cg.generateExpression(arg)
		if err != nil {
			return llvm.Value{}, err
		}
		allArgs = append(allArgs, val)
	}

	printf := cg.functions["printf"]
	return cg.builder.CreateCall(printf.GlobalValueType(), printf, allArgs, "printftmp"), nil
}

// ============================================================================
// STRING FORMATTING HELPERS
// ============================================================================

// generateSprintf generates formatted string output
// sprintf(buffer, format, args...) -> int (bytes written)
func (cg *LLVMCodeGenerator) generateSprintf(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) < 2 {
		return llvm.Value{}, nil
	}

	// Declare sprintf if not already
	if _, ok := cg.functions["sprintf"]; !ok {
		sprintfType := llvm.FunctionType(
			cg.context.Int32Type(),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0), // buffer
				llvm.PointerType(cg.context.Int8Type(), 0), // format
			},
			true, // variadic
		)
		sprintfFn := llvm.AddFunction(cg.module, "sprintf", sprintfType)
		sprintfFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["sprintf"] = sprintfFn
	}

	args := make([]llvm.Value, len(call.Args))
	for i, arg := range call.Args {
		val, err := cg.generateExpression(arg)
		if err != nil {
			return llvm.Value{}, err
		}
		args[i] = val
	}

	sprintf := cg.functions["sprintf"]
	result := cg.builder.CreateCall(sprintf.GlobalValueType(), sprintf, args, "sprintftmp")
	return cg.builder.CreateZExt(result, cg.context.Int64Type(), "sprintfext"), nil
}

// generateSnprintf generates formatted string output with size limit
// snprintf(buffer, size, format, args...) -> int (bytes written)
func (cg *LLVMCodeGenerator) generateSnprintf(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) < 3 {
		return llvm.Value{}, nil
	}

	// Declare snprintf if not already
	if _, ok := cg.functions["snprintf"]; !ok {
		snprintfType := llvm.FunctionType(
			cg.context.Int32Type(),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0), // buffer
				cg.context.Int64Type(),                     // size
				llvm.PointerType(cg.context.Int8Type(), 0), // format
			},
			true, // variadic
		)
		snprintfFn := llvm.AddFunction(cg.module, "snprintf", snprintfType)
		snprintfFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["snprintf"] = snprintfFn
	}

	args := make([]llvm.Value, len(call.Args))
	for i, arg := range call.Args {
		val, err := cg.generateExpression(arg)
		if err != nil {
			return llvm.Value{}, err
		}
		// Convert size argument (index 1) to i64 if needed
		if i == 1 && args[0].Type().IntTypeWidth() != 64 {
			args[i] = cg.builder.CreateZExt(val, cg.context.Int64Type(), "sizeext")
		} else {
			args[i] = val
		}
	}

	snprintf := cg.functions["snprintf"]
	result := cg.builder.CreateCall(snprintf.GlobalValueType(), snprintf, args, "snprintftmp")
	return cg.builder.CreateZExt(result, cg.context.Int64Type(), "snprintfext"), nil
}

// ============================================================================
// NUMBER FORMATTING
// ============================================================================

// generateFormatInt formats an integer with specified width and flags
// format_int(value, width, flags) -> string
func (cg *LLVMCodeGenerator) generateFormatInt(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) < 1 {
		return llvm.Value{}, nil
	}

	val, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Allocate buffer (32 bytes for int representation)
	bufSize := llvm.ConstInt(cg.context.Int64Type(), 32, false)
	malloc := cg.functions["malloc"]
	buf := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{bufSize}, "intbuf")

	// Get width (default 0)
	width := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	if len(call.Args) >= 2 {
		width, err = cg.generateExpression(call.Args[1])
		if err != nil {
			return llvm.Value{}, err
		}
	}

	// Format using sprintf with appropriate format string
	// Build format string based on width
	fmtStr := cg.createGlobalString("%ld")
	_ = width // Would use this for dynamic format string building

	// Declare sprintf if needed
	if _, ok := cg.functions["sprintf"]; !ok {
		sprintfType := llvm.FunctionType(
			cg.context.Int32Type(),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0),
				llvm.PointerType(cg.context.Int8Type(), 0),
			},
			true,
		)
		sprintfFn := llvm.AddFunction(cg.module, "sprintf", sprintfType)
		sprintfFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["sprintf"] = sprintfFn
	}

	sprintf := cg.functions["sprintf"]
	cg.builder.CreateCall(sprintf.GlobalValueType(), sprintf, []llvm.Value{buf, fmtStr, val}, "")

	return buf, nil
}

// generateFormatHex formats an integer as hexadecimal
// format_hex(value, uppercase, prefix) -> string
func (cg *LLVMCodeGenerator) generateFormatHex(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) < 1 {
		return llvm.Value{}, nil
	}

	val, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Check for uppercase flag
	uppercase := false
	if len(call.Args) >= 2 {
		if boolLit, ok := call.Args[1].(*BoolLiteral); ok {
			uppercase = boolLit.Value
		}
	}

	// Check for prefix flag
	prefix := false
	if len(call.Args) >= 3 {
		if boolLit, ok := call.Args[2].(*BoolLiteral); ok {
			prefix = boolLit.Value
		}
	}

	// Allocate buffer
	bufSize := llvm.ConstInt(cg.context.Int64Type(), 32, false)
	malloc := cg.functions["malloc"]
	buf := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{bufSize}, "hexbuf")

	// Build format string
	var fmtStr llvm.Value
	if prefix && uppercase {
		fmtStr = cg.createGlobalString("0x%lX")
	} else if prefix {
		fmtStr = cg.createGlobalString("0x%lx")
	} else if uppercase {
		fmtStr = cg.createGlobalString("%lX")
	} else {
		fmtStr = cg.createGlobalString("%lx")
	}

	// Declare sprintf if needed
	cg.declareSprintf()

	sprintf := cg.functions["sprintf"]
	cg.builder.CreateCall(sprintf.GlobalValueType(), sprintf, []llvm.Value{buf, fmtStr, val}, "")

	return buf, nil
}

// generateFormatBinary formats an integer as binary.
// format_bin(value) -> string (e.g. format_bin(5) -> "101"; format_bin(0) -> "0")
// Extracts bits from the top down with a real loop (no leading zeros, and no
// sprintf("%ld") shortcut - see P1-2 in FIXER_HANDOFF.md).
func (cg *LLVMCodeGenerator) generateFormatBinary(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) < 1 {
		return llvm.Value{}, fmt.Errorf("format_bin requires 1 argument")
	}

	val, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	val = cg.coerceToType(val, cg.context.Int64Type())

	i64 := cg.context.Int64Type()
	i8 := cg.context.Int8Type()
	i1 := cg.context.Int1Type()

	// Buffer: 64 bits + null terminator.
	bufSize := llvm.ConstInt(i64, 65, false)
	malloc := cg.functions["malloc"]
	buf := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{bufSize}, "binbuf")

	bitAlloca := cg.builder.CreateAlloca(i64, "bin_bit")
	idxAlloca := cg.builder.CreateAlloca(i64, "bin_idx")
	startedAlloca := cg.builder.CreateAlloca(i1, "bin_started")
	cg.builder.CreateStore(llvm.ConstInt(i64, 63, false), bitAlloca)
	cg.builder.CreateStore(llvm.ConstInt(i64, 0, false), idxAlloca)
	cg.builder.CreateStore(llvm.ConstInt(i1, 0, false), startedAlloca)

	condBlock := llvm.AddBasicBlock(cg.currentFn, "bin_cond")
	bodyBlock := llvm.AddBasicBlock(cg.currentFn, "bin_body")
	writeBlock := llvm.AddBasicBlock(cg.currentFn, "bin_write")
	nextBlock := llvm.AddBasicBlock(cg.currentFn, "bin_next")
	doneBlock := llvm.AddBasicBlock(cg.currentFn, "bin_done")

	cg.builder.CreateBr(condBlock)

	// while (bit >= 0)
	cg.builder.SetInsertPointAtEnd(condBlock)
	bit := cg.builder.CreateLoad(i64, bitAlloca, "bit")
	cond := cg.builder.CreateICmp(llvm.IntSGE, bit, llvm.ConstInt(i64, 0, false), "bincond")
	cg.builder.CreateCondBr(cond, bodyBlock, doneBlock)

	// bitVal = (val >> bit) & 1 ; write if it's a 1, we've already started, or
	// this is the last (lowest) bit and nothing was written yet (value == 0).
	cg.builder.SetInsertPointAtEnd(bodyBlock)
	bit = cg.builder.CreateLoad(i64, bitAlloca, "bit")
	shifted := cg.builder.CreateLShr(val, bit, "binshift")
	bitVal := cg.builder.CreateAnd(shifted, llvm.ConstInt(i64, 1, false), "binmask")
	isOne := cg.builder.CreateICmp(llvm.IntNE, bitVal, llvm.ConstInt(i64, 0, false), "binisone")
	started := cg.builder.CreateLoad(i1, startedAlloca, "started")
	isLastBit := cg.builder.CreateICmp(llvm.IntEQ, bit, llvm.ConstInt(i64, 0, false), "binislast")
	shouldWrite := cg.builder.CreateOr(cg.builder.CreateOr(isOne, started, "binor1"), isLastBit, "binor2")
	cg.builder.CreateCondBr(shouldWrite, writeBlock, nextBlock)

	cg.builder.SetInsertPointAtEnd(writeBlock)
	cg.builder.CreateStore(llvm.ConstInt(i1, 1, false), startedAlloca)
	idx := cg.builder.CreateLoad(i64, idxAlloca, "idx")
	digit := cg.builder.CreateSelect(isOne, llvm.ConstInt(i8, '1', false), llvm.ConstInt(i8, '0', false), "digit")
	digitPtr := cg.builder.CreateGEP(i8, buf, []llvm.Value{idx}, "digitptr")
	cg.builder.CreateStore(digit, digitPtr)
	cg.builder.CreateStore(cg.builder.CreateAdd(idx, llvm.ConstInt(i64, 1, false), "idxnext"), idxAlloca)
	cg.builder.CreateBr(nextBlock)

	cg.builder.SetInsertPointAtEnd(nextBlock)
	bit = cg.builder.CreateLoad(i64, bitAlloca, "bit")
	cg.builder.CreateStore(cg.builder.CreateSub(bit, llvm.ConstInt(i64, 1, false), "bitnext"), bitAlloca)
	cg.builder.CreateBr(condBlock)

	cg.builder.SetInsertPointAtEnd(doneBlock)
	finalIdx := cg.builder.CreateLoad(i64, idxAlloca, "finalidx")
	endPtr := cg.builder.CreateGEP(i8, buf, []llvm.Value{finalIdx}, "binend")
	cg.builder.CreateStore(llvm.ConstInt(i8, 0, false), endPtr)

	return buf, nil
}

// declareSprintf declares sprintf if not already declared
func (cg *LLVMCodeGenerator) declareSprintf() {
	if _, ok := cg.functions["sprintf"]; !ok {
		sprintfType := llvm.FunctionType(
			cg.context.Int32Type(),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0),
				llvm.PointerType(cg.context.Int8Type(), 0),
			},
			true,
		)
		sprintfFn := llvm.AddFunction(cg.module, "sprintf", sprintfType)
		sprintfFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["sprintf"] = sprintfFn
	}
}

// ============================================================================
// PAD/ALIGN HELPERS
// ============================================================================

// declareMemsetMemcpy declares the libc memset/memcpy externs used by the
// pad helpers below, if not already declared.
func (cg *LLVMCodeGenerator) declareMemsetMemcpy() {
	i8ptr := llvm.PointerType(cg.context.Int8Type(), 0)
	if _, ok := cg.functions["memset"]; !ok {
		memsetType := llvm.FunctionType(i8ptr, []llvm.Type{i8ptr, cg.context.Int32Type(), cg.context.Int64Type()}, false)
		fn := llvm.AddFunction(cg.module, "memset", memsetType)
		fn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["memset"] = fn
	}
	if _, ok := cg.functions["memcpy"]; !ok {
		memcpyType := llvm.FunctionType(i8ptr, []llvm.Type{i8ptr, i8ptr, cg.context.Int64Type()}, false)
		fn := llvm.AddFunction(cg.module, "memcpy", memcpyType)
		fn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["memcpy"] = fn
	}
}

// padArgs evaluates the (str, width, char) arguments shared by pad_left and
// pad_right. width and char accept any runtime expression, not just literals
// - the value is already an llvm.Value, there is no reason to require a
// literal here (see P1-2 in FIXER_HANDOFF.md).
func (cg *LLVMCodeGenerator) padArgs(call *FunctionCall) (str, width, padChar llvm.Value, err error) {
	str, err = cg.generateExpression(call.Args[0])
	if err != nil {
		return
	}
	widthVal, err2 := cg.generateExpression(call.Args[1])
	if err2 != nil {
		err = err2
		return
	}
	width = cg.coerceToType(widthVal, cg.context.Int64Type())

	padChar = llvm.ConstInt(cg.context.Int32Type(), ' ', false)
	if len(call.Args) >= 3 {
		c, err3 := cg.generateExpression(call.Args[2])
		if err3 != nil {
			err = err3
			return
		}
		padChar = cg.coerceToType(c, cg.context.Int32Type())
	}
	return
}

// generatePadLeft pads a string on the left to reach target width.
// pad_left(str, width, char) -> string
// Buffer is sized max(width, strlen)+1 so a string longer than width cannot
// overflow it; the string is never truncated.
func (cg *LLVMCodeGenerator) generatePadLeft(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) < 2 {
		return llvm.Value{}, fmt.Errorf("pad_left requires at least 2 arguments (str, width)")
	}
	str, width, padChar, err := cg.padArgs(call)
	if err != nil {
		return llvm.Value{}, err
	}

	cg.declareStringHelpers()
	cg.declareMemsetMemcpy()
	strlenFn := cg.functions["strlen"]
	strLen := cg.builder.CreateCall(strlenFn.GlobalValueType(), strlenFn, []llvm.Value{str}, "strlen")

	i64 := cg.context.Int64Type()
	zero := llvm.ConstInt(i64, 0, false)
	one := llvm.ConstInt(i64, 1, false)

	// padNeeded = max(width - strlen, 0)
	rawPad := cg.builder.CreateSub(width, strLen, "rawpad")
	padNeg := cg.builder.CreateICmp(llvm.IntSLT, rawPad, zero, "padneg")
	padNeeded := cg.builder.CreateSelect(padNeg, zero, rawPad, "padneeded")

	// bufLen = max(width, strlen)
	widthLess := cg.builder.CreateICmp(llvm.IntSLT, width, strLen, "widthless")
	bufLen := cg.builder.CreateSelect(widthLess, strLen, width, "buflen")
	bufSize := cg.builder.CreateAdd(bufLen, one, "bufsize")

	malloc := cg.functions["malloc"]
	buf := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{bufSize}, "padbuf")

	memset := cg.functions["memset"]
	cg.builder.CreateCall(memset.GlobalValueType(), memset, []llvm.Value{buf, padChar, padNeeded}, "")

	// Copy the string right after the padding: buf + padNeeded
	dst := cg.builder.CreateGEP(cg.context.Int8Type(), buf, []llvm.Value{padNeeded}, "paddst")
	strLenPlus1 := cg.builder.CreateAdd(strLen, one, "strlenp1")
	memcpy := cg.functions["memcpy"]
	cg.builder.CreateCall(memcpy.GlobalValueType(), memcpy, []llvm.Value{dst, str, strLenPlus1}, "")

	return buf, nil
}

// generatePadRight pads a string on the right to reach target width.
// pad_right(str, width, char) -> string
// Mirror of pad_left: copy the string first, then fill the remaining width
// with the pad character (overwriting the copied null terminator, then
// re-terminating at the end).
func (cg *LLVMCodeGenerator) generatePadRight(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) < 2 {
		return llvm.Value{}, fmt.Errorf("pad_right requires at least 2 arguments (str, width)")
	}
	str, width, padChar, err := cg.padArgs(call)
	if err != nil {
		return llvm.Value{}, err
	}

	cg.declareStringHelpers()
	cg.declareMemsetMemcpy()
	strlenFn := cg.functions["strlen"]
	strLen := cg.builder.CreateCall(strlenFn.GlobalValueType(), strlenFn, []llvm.Value{str}, "strlen")

	i64 := cg.context.Int64Type()
	zero := llvm.ConstInt(i64, 0, false)
	one := llvm.ConstInt(i64, 1, false)

	rawPad := cg.builder.CreateSub(width, strLen, "rawpad")
	padNeg := cg.builder.CreateICmp(llvm.IntSLT, rawPad, zero, "padneg")
	padNeeded := cg.builder.CreateSelect(padNeg, zero, rawPad, "padneeded")

	widthLess := cg.builder.CreateICmp(llvm.IntSLT, width, strLen, "widthless")
	bufLen := cg.builder.CreateSelect(widthLess, strLen, width, "buflen")
	bufSize := cg.builder.CreateAdd(bufLen, one, "bufsize")

	malloc := cg.functions["malloc"]
	buf := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{bufSize}, "padbuf")

	memcpy := cg.functions["memcpy"]
	cg.builder.CreateCall(memcpy.GlobalValueType(), memcpy, []llvm.Value{buf, str, strLen}, "")

	// Fill [strLen, strLen+padNeeded) with padChar, then null-terminate at bufLen
	fillDst := cg.builder.CreateGEP(cg.context.Int8Type(), buf, []llvm.Value{strLen}, "filldst")
	memset := cg.functions["memset"]
	cg.builder.CreateCall(memset.GlobalValueType(), memset, []llvm.Value{fillDst, padChar, padNeeded}, "")

	endPtr := cg.builder.CreateGEP(cg.context.Int8Type(), buf, []llvm.Value{bufLen}, "endptr")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int8Type(), 0, false), endPtr)

	return buf, nil
}
