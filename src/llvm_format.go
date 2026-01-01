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

// generateFormatBinary formats an integer as binary
// format_bin(value, prefix) -> string
func (cg *LLVMCodeGenerator) generateFormatBinary(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) < 1 {
		return llvm.Value{}, nil
	}

	val, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Allocate buffer (72 bytes for 64-bit binary + "0b" prefix + null)
	bufSize := llvm.ConstInt(cg.context.Int64Type(), 72, false)
	malloc := cg.functions["malloc"]
	buf := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{bufSize}, "binbuf")
	bufInt := cg.builder.CreatePtrToInt(buf, cg.context.Int64Type(), "bufint")

	// Generate binary string manually
	// This is a loop that extracts bits and writes '0' or '1'
	cg.generateBinaryStringLoop(val, bufInt)

	return buf, nil
}

// generateBinaryStringLoop generates a loop to convert int to binary string
func (cg *LLVMCodeGenerator) generateBinaryStringLoop(val, bufPtr llvm.Value) {
	// Simplified: just call sprintf with %ld for now
	// A real implementation would extract bits in a loop

	fmtStr := cg.createGlobalString("%ld")
	buf := cg.builder.CreateIntToPtr(bufPtr, llvm.PointerType(cg.context.Int8Type(), 0), "buf")

	cg.declareSprintf()
	sprintf := cg.functions["sprintf"]
	cg.builder.CreateCall(sprintf.GlobalValueType(), sprintf, []llvm.Value{buf, fmtStr, val}, "")
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

// generatePadLeft pads a string on the left to reach target width
// pad_left(str, width, char) -> string
func (cg *LLVMCodeGenerator) generatePadLeft(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) < 2 {
		return llvm.Value{}, nil
	}

	str, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	width, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Get padding character (default space)
	padChar := llvm.ConstInt(cg.context.Int8Type(), ' ', false)
	if len(call.Args) >= 3 {
		if charLit, ok := call.Args[2].(*IntLiteral); ok {
			padChar = llvm.ConstInt(cg.context.Int8Type(), uint64(charLit.Value), false)
		}
	}

	// Get string length
	cg.declareStringHelpers()
	strlen := cg.functions["strlen"]
	strLen := cg.builder.CreateCall(strlen.GlobalValueType(), strlen, []llvm.Value{str}, "strlen")

	// Allocate buffer (width + 1 for null)
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	bufSize := cg.builder.CreateAdd(width, one, "bufsize")
	malloc := cg.functions["malloc"]
	buf := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{bufSize}, "padbuf")

	// Calculate padding needed
	padNeeded := cg.builder.CreateSub(width, strLen, "padneeded")
	_ = padChar
	_ = padNeeded

	// For now, just copy the string (simplified)
	strcpy := cg.functions["strcpy"]
	cg.builder.CreateCall(strcpy.GlobalValueType(), strcpy, []llvm.Value{buf, str}, "")

	return buf, nil
}

// generatePadRight pads a string on the right to reach target width
// pad_right(str, width, char) -> string
func (cg *LLVMCodeGenerator) generatePadRight(call *FunctionCall) (llvm.Value, error) {
	// Similar to padLeft but pads on the right
	return cg.generatePadLeft(call) // Simplified for now
}
