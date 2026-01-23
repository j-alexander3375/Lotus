package main

// llvm_stdlib.go - LLVM IR generation for Lotus standard library functions
//
// This file implements stdlib functions for the LLVM backend using inline
// LLVM IR generation. These are compiled directly into the generated code
// rather than relying on external library calls.

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

// ============================================================================
// MATH FUNCTIONS
// ============================================================================

// generateAbs generates inline code for abs(x)
func (cg *LLVMCodeGenerator) generateAbs(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	val, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// abs(x) = x < 0 ? -x : x
	zero := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	isNeg := cg.builder.CreateICmp(llvm.IntSLT, val, zero, "isneg")
	neg := cg.builder.CreateNeg(val, "negval")
	return cg.builder.CreateSelect(isNeg, neg, val, "abstmp"), nil
}

// generateMin generates inline code for min(a, b)
func (cg *LLVMCodeGenerator) generateMin(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, nil
	}

	a, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	b, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// min(a, b) = a < b ? a : b
	cmp := cg.builder.CreateICmp(llvm.IntSLT, a, b, "mincmp")
	return cg.builder.CreateSelect(cmp, a, b, "mintmp"), nil
}

// generateMax generates inline code for max(a, b)
func (cg *LLVMCodeGenerator) generateMax(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, nil
	}

	a, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	b, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// max(a, b) = a > b ? a : b
	cmp := cg.builder.CreateICmp(llvm.IntSGT, a, b, "maxcmp")
	return cg.builder.CreateSelect(cmp, a, b, "maxtmp"), nil
}

// generateSqrt generates inline code for integer sqrt(x) using Newton's method
func (cg *LLVMCodeGenerator) generateSqrt(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	x, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// For integer sqrt, we'll call the C library sqrtf and convert
	// First, declare sqrtf if not already declared
	if _, ok := cg.functions["llvm.sqrt.f64"]; !ok {
		sqrtType := llvm.FunctionType(
			cg.context.DoubleType(),
			[]llvm.Type{cg.context.DoubleType()},
			false,
		)
		sqrtFn := llvm.AddFunction(cg.module, "llvm.sqrt.f64", sqrtType)
		cg.functions["llvm.sqrt.f64"] = sqrtFn
	}

	// Convert int to double
	xf := cg.builder.CreateSIToFP(x, cg.context.DoubleType(), "sqrtconv")

	// Call sqrt intrinsic
	sqrtFn := cg.functions["llvm.sqrt.f64"]
	result := cg.builder.CreateCall(sqrtFn.GlobalValueType(), sqrtFn, []llvm.Value{xf}, "sqrtcall")

	// Convert back to int
	return cg.builder.CreateFPToSI(result, cg.context.Int64Type(), "sqrtint"), nil
}

// generatePow generates inline code for integer pow(base, exp) using a loop
func (cg *LLVMCodeGenerator) generatePow(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("pow requires 2 arguments")
	}

	base, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	exp, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Pure integer pow using a loop (binary exponentiation)
	currentFn := cg.builder.GetInsertBlock().Parent()
	int64Type := cg.context.Int64Type()

	// Create basic blocks
	loopBB := cg.context.AddBasicBlock(currentFn, "pow_loop")
	bodyBB := cg.context.AddBasicBlock(currentFn, "pow_body")
	doneBB := cg.context.AddBasicBlock(currentFn, "pow_done")

	// Allocate result and counters
	resultPtr := cg.builder.CreateAlloca(int64Type, "pow_result")
	expPtr := cg.builder.CreateAlloca(int64Type, "pow_exp")
	basePtr := cg.builder.CreateAlloca(int64Type, "pow_base")

	// Initialize: result=1, exp=exp, base=base
	cg.builder.CreateStore(llvm.ConstInt(int64Type, 1, false), resultPtr)
	cg.builder.CreateStore(exp, expPtr)
	cg.builder.CreateStore(base, basePtr)

	cg.builder.CreateBr(loopBB)

	// Loop: while exp > 0
	cg.builder.SetInsertPointAtEnd(loopBB)
	expVal := cg.builder.CreateLoad(int64Type, expPtr, "exp_val")
	zero := llvm.ConstInt(int64Type, 0, false)
	cond := cg.builder.CreateICmp(llvm.IntSGT, expVal, zero, "pow_cond")
	cg.builder.CreateCondBr(cond, bodyBB, doneBB)

	// Body: if exp & 1, result *= base; base *= base; exp >>= 1
	cg.builder.SetInsertPointAtEnd(bodyBB)
	expVal2 := cg.builder.CreateLoad(int64Type, expPtr, "exp_val2")
	baseVal := cg.builder.CreateLoad(int64Type, basePtr, "base_val")
	resultVal := cg.builder.CreateLoad(int64Type, resultPtr, "result_val")

	// Check if exp is odd (exp & 1)
	one := llvm.ConstInt(int64Type, 1, false)
	oddBit := cg.builder.CreateAnd(expVal2, one, "odd_bit")
	isOdd := cg.builder.CreateICmp(llvm.IntNE, oddBit, zero, "is_odd")

	// result = isOdd ? result * base : result
	newResult := cg.builder.CreateMul(resultVal, baseVal, "new_result")
	resultSelect := cg.builder.CreateSelect(isOdd, newResult, resultVal, "result_select")
	cg.builder.CreateStore(resultSelect, resultPtr)

	// base = base * base
	newBase := cg.builder.CreateMul(baseVal, baseVal, "new_base")
	cg.builder.CreateStore(newBase, basePtr)

	// exp = exp >> 1
	newExp := cg.builder.CreateAShr(expVal2, one, "new_exp")
	cg.builder.CreateStore(newExp, expPtr)

	cg.builder.CreateBr(loopBB)

	// Done: return result
	cg.builder.SetInsertPointAtEnd(doneBB)
	return cg.builder.CreateLoad(int64Type, resultPtr, "pow_final"), nil
}

// generateGcd generates inline code for gcd(a, b) using Euclidean algorithm
func (cg *LLVMCodeGenerator) generateGcd(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("gcd requires 2 arguments")
	}

	a, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	b, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Create basic blocks for loop
	currentFn := cg.builder.GetInsertBlock().Parent()
	loopBB := llvm.AddBasicBlock(currentFn, "gcdloop")
	doneBB := llvm.AddBasicBlock(currentFn, "gcddone")
	_ = cg.builder.GetInsertBlock() // entry block (used to set up control flow)

	// Allocate variables for a and b
	aPtr := cg.builder.CreateAlloca(cg.context.Int64Type(), "gcd_a")
	bPtr := cg.builder.CreateAlloca(cg.context.Int64Type(), "gcd_b")
	cg.builder.CreateStore(a, aPtr)
	cg.builder.CreateStore(b, bPtr)

	cg.builder.CreateBr(loopBB)

	// Loop block: while b != 0
	cg.builder.SetInsertPointAtEnd(loopBB)
	bVal := cg.builder.CreateLoad(cg.context.Int64Type(), bPtr, "bval")
	zero := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	cond := cg.builder.CreateICmp(llvm.IntNE, bVal, zero, "gcdcond")

	bodyBB := llvm.AddBasicBlock(currentFn, "gcdbody")
	cg.builder.CreateCondBr(cond, bodyBB, doneBB)

	// Body: t = b; b = a % b; a = t
	cg.builder.SetInsertPointAtEnd(bodyBB)
	aVal := cg.builder.CreateLoad(cg.context.Int64Type(), aPtr, "aval")
	bVal2 := cg.builder.CreateLoad(cg.context.Int64Type(), bPtr, "bval2")
	rem := cg.builder.CreateSRem(aVal, bVal2, "gcdrem")
	cg.builder.CreateStore(bVal2, aPtr)
	cg.builder.CreateStore(rem, bPtr)
	cg.builder.CreateBr(loopBB)

	// Done: return a
	cg.builder.SetInsertPointAtEnd(doneBB)
	result := cg.builder.CreateLoad(cg.context.Int64Type(), aPtr, "gcdresult")
	return result, nil
}

// generateLcm generates inline code for lcm(a, b) = |a * b| / gcd(a, b)
func (cg *LLVMCodeGenerator) generateLcm(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("lcm requires 2 arguments")
	}

	// Get gcd first
	gcdVal, err := cg.generateGcd(call)
	if err != nil {
		return llvm.Value{}, err
	}

	// Get a and b again
	a, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	b, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// lcm = (a / gcd) * b to avoid overflow
	aDiv := cg.builder.CreateSDiv(a, gcdVal, "lcmdiv")
	result := cg.builder.CreateMul(aDiv, b, "lcmresult")

	// Return absolute value
	zero := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	isNeg := cg.builder.CreateICmp(llvm.IntSLT, result, zero, "lcmisneg")
	neg := cg.builder.CreateNeg(result, "lcmneg")
	return cg.builder.CreateSelect(isNeg, neg, result, "lcmabs"), nil
}

// generateFloor generates code for floor(x) - already an int, so just return it
func (cg *LLVMCodeGenerator) generateFloor(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}
	// For integer input, floor is identity
	return cg.generateExpression(call.Args[0])
}

// generateCeil generates code for ceil(x) - for int, just return it
func (cg *LLVMCodeGenerator) generateCeil(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}
	// For integer input, ceil is identity
	return cg.generateExpression(call.Args[0])
}

// generateRound generates code for round(x) - for int, just return it
func (cg *LLVMCodeGenerator) generateRound(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}
	// For integer input, round is identity
	return cg.generateExpression(call.Args[0])
}

// ============================================================================
// STRING FUNCTIONS
// ============================================================================

// generateStrlen generates code for len(str)
func (cg *LLVMCodeGenerator) generateStrlen(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	str, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Declare strlen if not already
	if _, ok := cg.functions["strlen"]; !ok {
		strlenType := llvm.FunctionType(
			cg.context.Int64Type(),
			[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0)},
			false,
		)
		strlenFn := llvm.AddFunction(cg.module, "strlen", strlenType)
		strlenFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strlen"] = strlenFn
	}

	strlen := cg.functions["strlen"]
	return cg.builder.CreateCall(strlen.GlobalValueType(), strlen, []llvm.Value{str}, "lentmp"), nil
}

// generateConcat generates code for concat(a, b)
func (cg *LLVMCodeGenerator) generateConcat(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, nil
	}

	a, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	b, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Declare helper functions
	cg.declareStringHelpers()

	// Get lengths
	strlen := cg.functions["strlen"]
	lenA := cg.builder.CreateCall(strlen.GlobalValueType(), strlen, []llvm.Value{a}, "lena")
	lenB := cg.builder.CreateCall(strlen.GlobalValueType(), strlen, []llvm.Value{b}, "lenb")

	// Allocate buffer for result (lenA + lenB + 1)
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	totalLen := cg.builder.CreateAdd(lenA, lenB, "totallen")
	bufLen := cg.builder.CreateAdd(totalLen, one, "buflen")

	malloc := cg.functions["malloc"]
	buf := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{bufLen}, "concatbuf")

	// Copy first string
	strcpy := cg.functions["strcpy"]
	cg.builder.CreateCall(strcpy.GlobalValueType(), strcpy, []llvm.Value{buf, a}, "")

	// Append second string
	strcat := cg.functions["strcat"]
	cg.builder.CreateCall(strcat.GlobalValueType(), strcat, []llvm.Value{buf, b}, "")

	return buf, nil
}

// generateContains generates code for contains(haystack, needle)
func (cg *LLVMCodeGenerator) generateContains(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, nil
	}

	haystack, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	needle, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Declare strstr if not already
	if _, ok := cg.functions["strstr"]; !ok {
		strstrType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0),
				llvm.PointerType(cg.context.Int8Type(), 0),
			},
			false,
		)
		strstrFn := llvm.AddFunction(cg.module, "strstr", strstrType)
		strstrFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strstr"] = strstrFn
	}

	strstr := cg.functions["strstr"]
	result := cg.builder.CreateCall(strstr.GlobalValueType(), strstr, []llvm.Value{haystack, needle}, "strstrres")

	// Convert pointer to bool (1 if not null, 0 if null)
	null := llvm.ConstNull(llvm.PointerType(cg.context.Int8Type(), 0))
	cmp := cg.builder.CreateICmp(llvm.IntNE, result, null, "containscmp")
	return cg.builder.CreateZExt(cmp, cg.context.Int64Type(), "containsext"), nil
}

// declareStringHelpers declares C string helper functions
func (cg *LLVMCodeGenerator) declareStringHelpers() {
	i8ptr := llvm.PointerType(cg.context.Int8Type(), 0)

	if _, ok := cg.functions["strlen"]; !ok {
		strlenType := llvm.FunctionType(cg.context.Int64Type(), []llvm.Type{i8ptr}, false)
		strlenFn := llvm.AddFunction(cg.module, "strlen", strlenType)
		strlenFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strlen"] = strlenFn
	}

	if _, ok := cg.functions["strcpy"]; !ok {
		strcpyType := llvm.FunctionType(i8ptr, []llvm.Type{i8ptr, i8ptr}, false)
		strcpyFn := llvm.AddFunction(cg.module, "strcpy", strcpyType)
		strcpyFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strcpy"] = strcpyFn
	}

	if _, ok := cg.functions["strcat"]; !ok {
		strcatType := llvm.FunctionType(i8ptr, []llvm.Type{i8ptr, i8ptr}, false)
		strcatFn := llvm.AddFunction(cg.module, "strcat", strcatType)
		strcatFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strcat"] = strcatFn
	}
}

// ============================================================================
// MEMORY FUNCTIONS
// ============================================================================

// generateMalloc generates code for malloc(size)
func (cg *LLVMCodeGenerator) generateMalloc(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	size, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	malloc := cg.functions["malloc"]
	ptr := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{size}, "malloctmp")
	// Cast to i64 for Lotus int type
	return cg.builder.CreatePtrToInt(ptr, cg.context.Int64Type(), "mallocint"), nil
}

// generateFree generates code for free(ptr)
func (cg *LLVMCodeGenerator) generateFree(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	ptr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Convert int back to pointer
	ptrVal := cg.builder.CreateIntToPtr(ptr, llvm.PointerType(cg.context.Int8Type(), 0), "freeptr")

	free := cg.functions["free"]
	cg.builder.CreateCall(free.GlobalValueType(), free, []llvm.Value{ptrVal}, "")
	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

// generateMemset generates code for memset(ptr, value, size)
func (cg *LLVMCodeGenerator) generateMemset(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, fmt.Errorf("memset requires 3 arguments")
	}

	ptr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	val, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	size, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}

	// Declare memset if not already
	if _, ok := cg.functions["memset"]; !ok {
		memsetType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0),
				cg.context.Int32Type(),
				cg.context.Int64Type(),
			},
			false,
		)
		memsetFn := llvm.AddFunction(cg.module, "memset", memsetType)
		memsetFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["memset"] = memsetFn
	}

	// Convert ptr (int) to pointer
	ptrVal := cg.builder.CreateIntToPtr(ptr, llvm.PointerType(cg.context.Int8Type(), 0), "memsetptr")
	// Truncate val to i32 (byte value)
	valI32 := cg.builder.CreateTrunc(val, cg.context.Int32Type(), "memsetval")

	memset := cg.functions["memset"]
	result := cg.builder.CreateCall(memset.GlobalValueType(), memset, []llvm.Value{ptrVal, valI32, size}, "memsetcall")
	return cg.builder.CreatePtrToInt(result, cg.context.Int64Type(), "memsetresult"), nil
}

// generateMemcpy generates code for memcpy(dest, src, size)
func (cg *LLVMCodeGenerator) generateMemcpy(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, fmt.Errorf("memcpy requires 3 arguments")
	}

	dest, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	src, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	size, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}

	// Declare memcpy if not already
	if _, ok := cg.functions["memcpy"]; !ok {
		memcpyType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0),
				llvm.PointerType(cg.context.Int8Type(), 0),
				cg.context.Int64Type(),
			},
			false,
		)
		memcpyFn := llvm.AddFunction(cg.module, "memcpy", memcpyType)
		memcpyFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["memcpy"] = memcpyFn
	}

	// Convert ptrs (ints) to pointers
	destPtr := cg.builder.CreateIntToPtr(dest, llvm.PointerType(cg.context.Int8Type(), 0), "memcpydest")
	srcPtr := cg.builder.CreateIntToPtr(src, llvm.PointerType(cg.context.Int8Type(), 0), "memcpysrc")

	memcpy := cg.functions["memcpy"]
	result := cg.builder.CreateCall(memcpy.GlobalValueType(), memcpy, []llvm.Value{destPtr, srcPtr, size}, "memcpycall")
	return cg.builder.CreatePtrToInt(result, cg.context.Int64Type(), "memcpyresult"), nil
}

// ============================================================================
// COLLECTIONS - DYNAMIC ARRAY (using malloc/free for simplicity)
// Array structure in memory: [capacity:8][length:8][data...]
// ============================================================================

// generateArrayIntNew creates a new dynamic array
func (cg *LLVMCodeGenerator) generateArrayIntNew(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	capacity, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Allocate: 16 bytes for header + capacity * 8 bytes for data
	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	sixteen := llvm.ConstInt(cg.context.Int64Type(), 16, false)
	dataSize := cg.builder.CreateMul(capacity, eight, "datasize")
	totalSize := cg.builder.CreateAdd(sixteen, dataSize, "totalsize")

	malloc := cg.functions["malloc"]
	ptr := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{totalSize}, "arraymalloc")
	ptrInt := cg.builder.CreatePtrToInt(ptr, cg.context.Int64Type(), "arrayptr")

	// Store capacity at offset 0
	capPtr := cg.builder.CreateIntToPtr(ptrInt, llvm.PointerType(cg.context.Int64Type(), 0), "capptr")
	cg.builder.CreateStore(capacity, capPtr)

	// Store length (0) at offset 8
	lenOffset := cg.builder.CreateAdd(ptrInt, eight, "lenoffset")
	lenPtr := cg.builder.CreateIntToPtr(lenOffset, llvm.PointerType(cg.context.Int64Type(), 0), "lenptr")
	zero := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	cg.builder.CreateStore(zero, lenPtr)

	return ptrInt, nil
}

// generateArrayIntLen returns the length of the array
func (cg *LLVMCodeGenerator) generateArrayIntLen(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	arrPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Length is at offset 8
	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	lenOffset := cg.builder.CreateAdd(arrPtr, eight, "lenoffset")
	lenPtr := cg.builder.CreateIntToPtr(lenOffset, llvm.PointerType(cg.context.Int64Type(), 0), "lenptr")
	return cg.builder.CreateLoad(cg.context.Int64Type(), lenPtr, "arrlen"), nil
}

// generateArrayIntPush appends an element to the array
func (cg *LLVMCodeGenerator) generateArrayIntPush(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, nil
	}

	arrPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	value, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	sixteen := llvm.ConstInt(cg.context.Int64Type(), 16, false)

	// Get current length
	lenOffset := cg.builder.CreateAdd(arrPtr, eight, "lenoffset")
	lenPtr := cg.builder.CreateIntToPtr(lenOffset, llvm.PointerType(cg.context.Int64Type(), 0), "lenptr")
	length := cg.builder.CreateLoad(cg.context.Int64Type(), lenPtr, "arrlen")

	// Calculate data position: arrPtr + 16 + length * 8
	dataOffset := cg.builder.CreateMul(length, eight, "dataoffset")
	dataBase := cg.builder.CreateAdd(arrPtr, sixteen, "database")
	elemOffset := cg.builder.CreateAdd(dataBase, dataOffset, "elemoffset")
	elemPtr := cg.builder.CreateIntToPtr(elemOffset, llvm.PointerType(cg.context.Int64Type(), 0), "elemptr")

	// Store value
	cg.builder.CreateStore(value, elemPtr)

	// Increment length
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	newLen := cg.builder.CreateAdd(length, one, "newlen")
	cg.builder.CreateStore(newLen, lenPtr)

	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

// generateArrayIntPop removes and returns the last element
func (cg *LLVMCodeGenerator) generateArrayIntPop(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	arrPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	sixteen := llvm.ConstInt(cg.context.Int64Type(), 16, false)
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)

	// Get current length
	lenOffset := cg.builder.CreateAdd(arrPtr, eight, "lenoffset")
	lenPtr := cg.builder.CreateIntToPtr(lenOffset, llvm.PointerType(cg.context.Int64Type(), 0), "lenptr")
	length := cg.builder.CreateLoad(cg.context.Int64Type(), lenPtr, "arrlen")

	// Decrement length first
	newLen := cg.builder.CreateSub(length, one, "newlen")
	cg.builder.CreateStore(newLen, lenPtr)

	// Get element at new length position
	dataOffset := cg.builder.CreateMul(newLen, eight, "dataoffset")
	dataBase := cg.builder.CreateAdd(arrPtr, sixteen, "database")
	elemOffset := cg.builder.CreateAdd(dataBase, dataOffset, "elemoffset")
	elemPtr := cg.builder.CreateIntToPtr(elemOffset, llvm.PointerType(cg.context.Int64Type(), 0), "elemptr")

	return cg.builder.CreateLoad(cg.context.Int64Type(), elemPtr, "popval"), nil
}

// generateArrayIntGet gets an element at index
func (cg *LLVMCodeGenerator) generateArrayIntGet(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, nil
	}

	arrPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	index, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	sixteen := llvm.ConstInt(cg.context.Int64Type(), 16, false)

	// Calculate element position: arrPtr + 16 + index * 8
	dataOffset := cg.builder.CreateMul(index, eight, "dataoffset")
	dataBase := cg.builder.CreateAdd(arrPtr, sixteen, "database")
	elemOffset := cg.builder.CreateAdd(dataBase, dataOffset, "elemoffset")
	elemPtr := cg.builder.CreateIntToPtr(elemOffset, llvm.PointerType(cg.context.Int64Type(), 0), "elemptr")

	return cg.builder.CreateLoad(cg.context.Int64Type(), elemPtr, "getval"), nil
}

// generateArrayIntSet sets an element at index
func (cg *LLVMCodeGenerator) generateArrayIntSet(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, nil
	}

	arrPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	index, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	value, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}

	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	sixteen := llvm.ConstInt(cg.context.Int64Type(), 16, false)

	// Calculate element position
	dataOffset := cg.builder.CreateMul(index, eight, "dataoffset")
	dataBase := cg.builder.CreateAdd(arrPtr, sixteen, "database")
	elemOffset := cg.builder.CreateAdd(dataBase, dataOffset, "elemoffset")
	elemPtr := cg.builder.CreateIntToPtr(elemOffset, llvm.PointerType(cg.context.Int64Type(), 0), "elemptr")

	cg.builder.CreateStore(value, elemPtr)
	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

// generateArrayIntCapacity returns the capacity of the array
func (cg *LLVMCodeGenerator) generateArrayIntCapacity(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	arrPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Capacity is at offset 0
	capPtr := cg.builder.CreateIntToPtr(arrPtr, llvm.PointerType(cg.context.Int64Type(), 0), "capptr")
	return cg.builder.CreateLoad(cg.context.Int64Type(), capPtr, "arrcap"), nil
}

// generateArrayIntResize resizes the array (no-op for now, capacity is fixed)
func (cg *LLVMCodeGenerator) generateArrayIntResize(call *FunctionCall) (llvm.Value, error) {
	// Simplified: just return 0 (no actual resize in this implementation)
	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

// generateArrayIntReserve reserves capacity (no-op for now)
func (cg *LLVMCodeGenerator) generateArrayIntReserve(call *FunctionCall) (llvm.Value, error) {
	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

// generateArrayIntShrink shrinks to fit (no-op for now)
func (cg *LLVMCodeGenerator) generateArrayIntShrink(call *FunctionCall) (llvm.Value, error) {
	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

// generateArrayIntFree frees the array
func (cg *LLVMCodeGenerator) generateArrayIntFree(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	arrPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	ptr := cg.builder.CreateIntToPtr(arrPtr, llvm.PointerType(cg.context.Int8Type(), 0), "freeptr")
	free := cg.functions["free"]
	cg.builder.CreateCall(free.GlobalValueType(), free, []llvm.Value{ptr}, "")
	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

// ============================================================================
// STACK (same structure as array - LIFO access)
// ============================================================================

func (cg *LLVMCodeGenerator) generateStackIntNew(call *FunctionCall) (llvm.Value, error) {
	return cg.generateArrayIntNew(call) // Same implementation
}

func (cg *LLVMCodeGenerator) generateStackIntPush(call *FunctionCall) (llvm.Value, error) {
	return cg.generateArrayIntPush(call) // Same implementation
}

func (cg *LLVMCodeGenerator) generateStackIntPop(call *FunctionCall) (llvm.Value, error) {
	return cg.generateArrayIntPop(call) // Same implementation
}

func (cg *LLVMCodeGenerator) generateStackIntLen(call *FunctionCall) (llvm.Value, error) {
	return cg.generateArrayIntLen(call) // Same implementation
}

// ============================================================================
// QUEUE (ring buffer with head/tail pointers)
// Structure: [capacity:8][length:8][head:8][tail:8][data...]
// Simplified: just use array semantics for now
// ============================================================================

func (cg *LLVMCodeGenerator) generateQueueIntNew(call *FunctionCall) (llvm.Value, error) {
	return cg.generateArrayIntNew(call)
}

func (cg *LLVMCodeGenerator) generateQueueIntEnqueue(call *FunctionCall) (llvm.Value, error) {
	return cg.generateArrayIntPush(call)
}

// generateQueueIntDequeue removes from front and shifts all elements left
func (cg *LLVMCodeGenerator) generateQueueIntDequeue(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	arrPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	sixteen := llvm.ConstInt(cg.context.Int64Type(), 16, false)
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)

	// Get data pointer
	dataBase := cg.builder.CreateAdd(arrPtr, sixteen, "database")
	dataPtr := cg.builder.CreateIntToPtr(dataBase, llvm.PointerType(cg.context.Int64Type(), 0), "dataptr")

	// Save first element to return
	firstVal := cg.builder.CreateLoad(cg.context.Int64Type(), dataPtr, "firstval")

	// Get and decrement length
	lenOffset := cg.builder.CreateAdd(arrPtr, eight, "lenoffset")
	lenPtr := cg.builder.CreateIntToPtr(lenOffset, llvm.PointerType(cg.context.Int64Type(), 0), "lenptr")
	length := cg.builder.CreateLoad(cg.context.Int64Type(), lenPtr, "arrlen")
	newLen := cg.builder.CreateSub(length, one, "newlen")
	cg.builder.CreateStore(newLen, lenPtr)

	// Shift all elements left by 1
	currentFn := cg.builder.GetInsertBlock().Parent()
	shiftLoopBB := cg.context.AddBasicBlock(currentFn, "deq_shift_loop")
	shiftBodyBB := cg.context.AddBasicBlock(currentFn, "deq_shift_body")
	shiftDoneBB := cg.context.AddBasicBlock(currentFn, "deq_shift_done")

	idxPtr := cg.builder.CreateAlloca(cg.context.Int64Type(), "shiftidx")
	zero := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	cg.builder.CreateStore(zero, idxPtr)
	cg.builder.CreateBr(shiftLoopBB)

	cg.builder.SetInsertPointAtEnd(shiftLoopBB)
	idx := cg.builder.CreateLoad(cg.context.Int64Type(), idxPtr, "idx")
	ltNewLen := cg.builder.CreateICmp(llvm.IntULT, idx, newLen, "ltnewlen")
	cg.builder.CreateCondBr(ltNewLen, shiftBodyBB, shiftDoneBB)

	cg.builder.SetInsertPointAtEnd(shiftBodyBB)
	// data[idx] = data[idx + 1]
	nextIdx := cg.builder.CreateAdd(idx, one, "nextidx")
	srcPtr := cg.builder.CreateGEP(cg.context.Int64Type(), dataPtr, []llvm.Value{nextIdx}, "srcptr")
	dstPtr := cg.builder.CreateGEP(cg.context.Int64Type(), dataPtr, []llvm.Value{idx}, "dstptr")
	elem := cg.builder.CreateLoad(cg.context.Int64Type(), srcPtr, "elem")
	cg.builder.CreateStore(elem, dstPtr)
	cg.builder.CreateStore(nextIdx, idxPtr)
	cg.builder.CreateBr(shiftLoopBB)

	cg.builder.SetInsertPointAtEnd(shiftDoneBB)
	return firstVal, nil
}

func (cg *LLVMCodeGenerator) generateQueueIntLen(call *FunctionCall) (llvm.Value, error) {
	return cg.generateArrayIntLen(call)
}

// ============================================================================
// DEQUE (double-ended queue)
// ============================================================================

func (cg *LLVMCodeGenerator) generateDequeIntNew(call *FunctionCall) (llvm.Value, error) {
	return cg.generateArrayIntNew(call)
}

func (cg *LLVMCodeGenerator) generateDequeIntPushBack(call *FunctionCall) (llvm.Value, error) {
	return cg.generateArrayIntPush(call)
}

func (cg *LLVMCodeGenerator) generateDequeIntPushFront(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, nil
	}

	arrPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	val, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Load current length from offset 8
	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	lenAddr := cg.builder.CreateAdd(arrPtr, eight, "lenaddr")
	lenPtr := cg.builder.CreateIntToPtr(lenAddr, llvm.PointerType(cg.context.Int64Type(), 0), "lenptr")
	length := cg.builder.CreateLoad(cg.context.Int64Type(), lenPtr, "len")

	// Get data base at offset 16
	sixteen := llvm.ConstInt(cg.context.Int64Type(), 16, false)
	dataBase := cg.builder.CreateAdd(arrPtr, sixteen, "database")
	dataPtr := cg.builder.CreateIntToPtr(dataBase, llvm.PointerType(cg.context.Int64Type(), 0), "dataptr")

	// Shift all elements right by 1 (from end to start to avoid overwriting)
	currentFn := cg.builder.GetInsertBlock().Parent()
	shiftLoopBB := cg.context.AddBasicBlock(currentFn, "shift_loop")
	shiftBodyBB := cg.context.AddBasicBlock(currentFn, "shift_body")
	shiftDoneBB := cg.context.AddBasicBlock(currentFn, "shift_done")

	idxPtr := cg.builder.CreateAlloca(cg.context.Int64Type(), "shiftidx")
	cg.builder.CreateStore(length, idxPtr)
	cg.builder.CreateBr(shiftLoopBB)

	cg.builder.SetInsertPointAtEnd(shiftLoopBB)
	idx := cg.builder.CreateLoad(cg.context.Int64Type(), idxPtr, "idx")
	zero := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	gtZero := cg.builder.CreateICmp(llvm.IntSGT, idx, zero, "gtzero")
	cg.builder.CreateCondBr(gtZero, shiftBodyBB, shiftDoneBB)

	cg.builder.SetInsertPointAtEnd(shiftBodyBB)
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	prevIdx := cg.builder.CreateSub(idx, one, "previdx")
	srcPtr := cg.builder.CreateGEP(cg.context.Int64Type(), dataPtr, []llvm.Value{prevIdx}, "srcptr")
	dstPtr := cg.builder.CreateGEP(cg.context.Int64Type(), dataPtr, []llvm.Value{idx}, "dstptr")
	elem := cg.builder.CreateLoad(cg.context.Int64Type(), srcPtr, "elem")
	cg.builder.CreateStore(elem, dstPtr)
	cg.builder.CreateStore(prevIdx, idxPtr)
	cg.builder.CreateBr(shiftLoopBB)

	cg.builder.SetInsertPointAtEnd(shiftDoneBB)
	// Store new value at index 0
	cg.builder.CreateStore(val, dataPtr)

	// Increment length
	newLen := cg.builder.CreateAdd(length, one, "newlen")
	cg.builder.CreateStore(newLen, lenPtr)

	return llvm.Value{}, nil
}

func (cg *LLVMCodeGenerator) generateDequeIntPopBack(call *FunctionCall) (llvm.Value, error) {
	return cg.generateArrayIntPop(call)
}

func (cg *LLVMCodeGenerator) generateDequeIntPopFront(call *FunctionCall) (llvm.Value, error) {
	return cg.generateQueueIntDequeue(call)
}

func (cg *LLVMCodeGenerator) generateDequeIntLen(call *FunctionCall) (llvm.Value, error) {
	return cg.generateArrayIntLen(call)
}

// ============================================================================
// HEAP (min-heap using array)
// ============================================================================

func (cg *LLVMCodeGenerator) generateHeapIntNew(call *FunctionCall) (llvm.Value, error) {
	return cg.generateArrayIntNew(call)
}

func (cg *LLVMCodeGenerator) generateHeapIntPush(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, nil
	}

	arrPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	val, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Load current length
	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	lenAddr := cg.builder.CreateAdd(arrPtr, eight, "lenaddr")
	lenPtr := cg.builder.CreateIntToPtr(lenAddr, llvm.PointerType(cg.context.Int64Type(), 0), "lenptr")
	length := cg.builder.CreateLoad(cg.context.Int64Type(), lenPtr, "len")

	// Get data pointer
	sixteen := llvm.ConstInt(cg.context.Int64Type(), 16, false)
	dataBase := cg.builder.CreateAdd(arrPtr, sixteen, "database")
	dataPtr := cg.builder.CreateIntToPtr(dataBase, llvm.PointerType(cg.context.Int64Type(), 0), "dataptr")

	// Add element at the end
	elemPtr := cg.builder.CreateGEP(cg.context.Int64Type(), dataPtr, []llvm.Value{length}, "elemptr")
	cg.builder.CreateStore(val, elemPtr)

	// Increment length
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	newLen := cg.builder.CreateAdd(length, one, "newlen")
	cg.builder.CreateStore(newLen, lenPtr)

	// Bubble up: while idx > 0 && heap[idx] < heap[parent], swap
	currentFn := cg.builder.GetInsertBlock().Parent()
	bubbleLoopBB := cg.context.AddBasicBlock(currentFn, "bubble_loop")
	bubbleCheckBB := cg.context.AddBasicBlock(currentFn, "bubble_check")
	bubbleSwapBB := cg.context.AddBasicBlock(currentFn, "bubble_swap")
	bubbleDoneBB := cg.context.AddBasicBlock(currentFn, "bubble_done")

	idxPtr := cg.builder.CreateAlloca(cg.context.Int64Type(), "bubbleidx")
	cg.builder.CreateStore(length, idxPtr)
	cg.builder.CreateBr(bubbleLoopBB)

	cg.builder.SetInsertPointAtEnd(bubbleLoopBB)
	idx := cg.builder.CreateLoad(cg.context.Int64Type(), idxPtr, "idx")
	zero := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	gtZero := cg.builder.CreateICmp(llvm.IntSGT, idx, zero, "gtzero")
	cg.builder.CreateCondBr(gtZero, bubbleCheckBB, bubbleDoneBB)

	cg.builder.SetInsertPointAtEnd(bubbleCheckBB)
	// Parent = (idx - 1) / 2
	idxMinus1 := cg.builder.CreateSub(idx, one, "idxm1")
	two := llvm.ConstInt(cg.context.Int64Type(), 2, false)
	parentIdx := cg.builder.CreateUDiv(idxMinus1, two, "parentidx")

	// Load current and parent values
	currentPtr := cg.builder.CreateGEP(cg.context.Int64Type(), dataPtr, []llvm.Value{idx}, "currptr")
	parentPtr := cg.builder.CreateGEP(cg.context.Int64Type(), dataPtr, []llvm.Value{parentIdx}, "parptr")
	currentVal := cg.builder.CreateLoad(cg.context.Int64Type(), currentPtr, "currval")
	parentVal := cg.builder.CreateLoad(cg.context.Int64Type(), parentPtr, "parval")

	// If current < parent, swap (min-heap)
	needsSwap := cg.builder.CreateICmp(llvm.IntSLT, currentVal, parentVal, "needsswap")
	cg.builder.CreateCondBr(needsSwap, bubbleSwapBB, bubbleDoneBB)

	cg.builder.SetInsertPointAtEnd(bubbleSwapBB)
	cg.builder.CreateStore(currentVal, parentPtr)
	cg.builder.CreateStore(parentVal, currentPtr)
	cg.builder.CreateStore(parentIdx, idxPtr)
	cg.builder.CreateBr(bubbleLoopBB)

	cg.builder.SetInsertPointAtEnd(bubbleDoneBB)
	return llvm.Value{}, nil
}

func (cg *LLVMCodeGenerator) generateHeapIntPop(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	arrPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Load current length
	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	lenAddr := cg.builder.CreateAdd(arrPtr, eight, "lenaddr")
	lenPtr := cg.builder.CreateIntToPtr(lenAddr, llvm.PointerType(cg.context.Int64Type(), 0), "lenptr")
	length := cg.builder.CreateLoad(cg.context.Int64Type(), lenPtr, "len")

	// Get data pointer
	sixteen := llvm.ConstInt(cg.context.Int64Type(), 16, false)
	dataBase := cg.builder.CreateAdd(arrPtr, sixteen, "database")
	dataPtr := cg.builder.CreateIntToPtr(dataBase, llvm.PointerType(cg.context.Int64Type(), 0), "dataptr")

	// Save root value (the minimum)
	rootVal := cg.builder.CreateLoad(cg.context.Int64Type(), dataPtr, "rootval")

	// Decrement length
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	newLen := cg.builder.CreateSub(length, one, "newlen")
	cg.builder.CreateStore(newLen, lenPtr)

	// Move last element to root
	lastPtr := cg.builder.CreateGEP(cg.context.Int64Type(), dataPtr, []llvm.Value{newLen}, "lastptr")
	lastVal := cg.builder.CreateLoad(cg.context.Int64Type(), lastPtr, "lastval")
	cg.builder.CreateStore(lastVal, dataPtr)

	// Sift down: while has children and parent > smallest child, swap
	currentFn := cg.builder.GetInsertBlock().Parent()
	siftLoopBB := cg.context.AddBasicBlock(currentFn, "sift_loop")
	siftCheckBB := cg.context.AddBasicBlock(currentFn, "sift_check")
	siftSwapBB := cg.context.AddBasicBlock(currentFn, "sift_swap")
	siftDoneBB := cg.context.AddBasicBlock(currentFn, "sift_done")

	idxPtr := cg.builder.CreateAlloca(cg.context.Int64Type(), "siftidx")
	zero := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	cg.builder.CreateStore(zero, idxPtr)
	cg.builder.CreateBr(siftLoopBB)

	cg.builder.SetInsertPointAtEnd(siftLoopBB)
	idx := cg.builder.CreateLoad(cg.context.Int64Type(), idxPtr, "idx")
	// Left child = 2*idx + 1
	two := llvm.ConstInt(cg.context.Int64Type(), 2, false)
	leftIdx := cg.builder.CreateMul(idx, two, "left1")
	leftIdx = cg.builder.CreateAdd(leftIdx, one, "leftidx")
	hasLeft := cg.builder.CreateICmp(llvm.IntULT, leftIdx, newLen, "hasleft")
	cg.builder.CreateCondBr(hasLeft, siftCheckBB, siftDoneBB)

	cg.builder.SetInsertPointAtEnd(siftCheckBB)
	// Right child = 2*idx + 2
	rightIdx := cg.builder.CreateAdd(leftIdx, one, "rightidx")

	// Find smallest child
	smallestPtr := cg.builder.CreateAlloca(cg.context.Int64Type(), "smallest")
	cg.builder.CreateStore(leftIdx, smallestPtr)

	// Check if right exists and is smaller
	hasRight := cg.builder.CreateICmp(llvm.IntULT, rightIdx, newLen, "hasright")

	checkRightBB := cg.context.AddBasicBlock(currentFn, "check_right")
	afterRightBB := cg.context.AddBasicBlock(currentFn, "after_right")
	cg.builder.CreateCondBr(hasRight, checkRightBB, afterRightBB)

	cg.builder.SetInsertPointAtEnd(checkRightBB)
	leftPtr := cg.builder.CreateGEP(cg.context.Int64Type(), dataPtr, []llvm.Value{leftIdx}, "leftptr")
	rightPtr := cg.builder.CreateGEP(cg.context.Int64Type(), dataPtr, []llvm.Value{rightIdx}, "rightptr")
	leftVal := cg.builder.CreateLoad(cg.context.Int64Type(), leftPtr, "leftval")
	rightVal := cg.builder.CreateLoad(cg.context.Int64Type(), rightPtr, "rightval")
	rightSmaller := cg.builder.CreateICmp(llvm.IntSLT, rightVal, leftVal, "rightsmaller")
	selectedIdx := cg.builder.CreateSelect(rightSmaller, rightIdx, leftIdx, "selectedidx")
	cg.builder.CreateStore(selectedIdx, smallestPtr)
	cg.builder.CreateBr(afterRightBB)

	cg.builder.SetInsertPointAtEnd(afterRightBB)
	smallest := cg.builder.CreateLoad(cg.context.Int64Type(), smallestPtr, "smallest")
	smallestChildPtr := cg.builder.CreateGEP(cg.context.Int64Type(), dataPtr, []llvm.Value{smallest}, "smallestptr")
	smallestVal := cg.builder.CreateLoad(cg.context.Int64Type(), smallestChildPtr, "smallestval")

	currentIdxPtr := cg.builder.CreateGEP(cg.context.Int64Type(), dataPtr, []llvm.Value{idx}, "curridxptr")
	currentVal := cg.builder.CreateLoad(cg.context.Int64Type(), currentIdxPtr, "currval")

	needsSwap := cg.builder.CreateICmp(llvm.IntSGT, currentVal, smallestVal, "needsswap")
	cg.builder.CreateCondBr(needsSwap, siftSwapBB, siftDoneBB)

	cg.builder.SetInsertPointAtEnd(siftSwapBB)
	cg.builder.CreateStore(currentVal, smallestChildPtr)
	cg.builder.CreateStore(smallestVal, currentIdxPtr)
	cg.builder.CreateStore(smallest, idxPtr)
	cg.builder.CreateBr(siftLoopBB)

	cg.builder.SetInsertPointAtEnd(siftDoneBB)
	return rootVal, nil
}

func (cg *LLVMCodeGenerator) generateHeapIntPeek(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	arrPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	sixteen := llvm.ConstInt(cg.context.Int64Type(), 16, false)
	dataBase := cg.builder.CreateAdd(arrPtr, sixteen, "database")
	firstPtr := cg.builder.CreateIntToPtr(dataBase, llvm.PointerType(cg.context.Int64Type(), 0), "firstptr")
	return cg.builder.CreateLoad(cg.context.Int64Type(), firstPtr, "peekval"), nil
}

func (cg *LLVMCodeGenerator) generateHeapIntLen(call *FunctionCall) (llvm.Value, error) {
	return cg.generateArrayIntLen(call)
}

// ============================================================================
// NET/SOCKET FUNCTIONS (using Linux syscalls)
// ============================================================================

// generateSocket creates a socket using syscall
func (cg *LLVMCodeGenerator) generateSocket(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, nil
	}

	domain, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	sockType, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	protocol, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}

	// Declare syscall if not already
	cg.declareSyscall()

	// socket() is syscall 41 on x86_64 Linux
	syscallNum := llvm.ConstInt(cg.context.Int64Type(), 41, false)
	syscall := cg.functions["syscall"]

	return cg.builder.CreateCall(syscall.GlobalValueType(), syscall,
		[]llvm.Value{syscallNum, domain, sockType, protocol}, "sockettmp"), nil
}

// generateClose closes a file descriptor using syscall
func (cg *LLVMCodeGenerator) generateClose(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	fd, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	cg.declareSyscall()

	// close() is syscall 3 on x86_64 Linux
	syscallNum := llvm.ConstInt(cg.context.Int64Type(), 3, false)
	zero := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	syscall := cg.functions["syscall"]

	return cg.builder.CreateCall(syscall.GlobalValueType(), syscall,
		[]llvm.Value{syscallNum, fd, zero, zero}, "closetmp"), nil
}

// ============================================================================
// NETWORKING FUNCTIONS (TCP/IP using Linux syscalls)
// ============================================================================

// generateConnectIPv4 connects to an IPv4 address
// connect_ipv4(fd, ip_addr, port) where ip_addr is packed as uint32 (e.g., 127.0.0.1 = 0x7f000001)
func (cg *LLVMCodeGenerator) generateConnectIPv4(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, fmt.Errorf("connect_ipv4 requires 3 arguments (fd, ip, port)")
	}

	fd, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	ipAddr, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	port, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}

	cg.declareSyscall()

	// Build sockaddr_in structure on stack (16 bytes)
	// struct sockaddr_in { sa_family_t sin_family; in_port_t sin_port; struct in_addr sin_addr; char sin_zero[8]; }
	sockaddrType := llvm.ArrayType(cg.context.Int8Type(), 16)
	sockaddr := cg.builder.CreateAlloca(sockaddrType, "sockaddr")

	// Zero out the structure
	memset := cg.functions["memset"]
	if memset.IsNil() {
		memsetType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0), cg.context.Int32Type(), cg.context.Int64Type()},
			false,
		)
		memset = llvm.AddFunction(cg.module, "memset", memsetType)
		cg.functions["memset"] = memset
	}
	sockaddrPtr := cg.builder.CreateBitCast(sockaddr, llvm.PointerType(cg.context.Int8Type(), 0), "sockaddrptr")
	cg.builder.CreateCall(memset.GlobalValueType(), memset,
		[]llvm.Value{sockaddrPtr, llvm.ConstInt(cg.context.Int32Type(), 0, false), llvm.ConstInt(cg.context.Int64Type(), 16, false)}, "")

	// Set sin_family = AF_INET (2) at offset 0 (2 bytes)
	familyPtr := cg.builder.CreateGEP(cg.context.Int8Type(), sockaddrPtr, []llvm.Value{llvm.ConstInt(cg.context.Int32Type(), 0, false)}, "familyptr")
	familyPtr16 := cg.builder.CreateBitCast(familyPtr, llvm.PointerType(cg.context.Int16Type(), 0), "familyptr16")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int16Type(), 2, false), familyPtr16)

	// Set sin_port at offset 2 (2 bytes, network byte order = big endian)
	portPtr := cg.builder.CreateGEP(cg.context.Int8Type(), sockaddrPtr, []llvm.Value{llvm.ConstInt(cg.context.Int32Type(), 2, false)}, "portptr")
	portPtr16 := cg.builder.CreateBitCast(portPtr, llvm.PointerType(cg.context.Int16Type(), 0), "portptr16")
	// Convert port to network byte order (big endian) using bswap
	port16 := cg.builder.CreateTrunc(port, cg.context.Int16Type(), "port16")
	// Manual byte swap: (port << 8) | (port >> 8)
	port16Ext := cg.builder.CreateZExt(port16, cg.context.Int32Type(), "port16ext")
	shifted := cg.builder.CreateShl(port16Ext, llvm.ConstInt(cg.context.Int32Type(), 8, false), "shifted")
	unshifted := cg.builder.CreateLShr(port16Ext, llvm.ConstInt(cg.context.Int32Type(), 8, false), "unshifted")
	portBE := cg.builder.CreateOr(shifted, unshifted, "portbe")
	portBE16 := cg.builder.CreateTrunc(portBE, cg.context.Int16Type(), "portbe16")
	cg.builder.CreateStore(portBE16, portPtr16)

	// Set sin_addr at offset 4 (convert to network byte order - big endian)
	addrPtr := cg.builder.CreateGEP(cg.context.Int8Type(), sockaddrPtr, []llvm.Value{llvm.ConstInt(cg.context.Int32Type(), 4, false)}, "addrptr")
	addrPtr32 := cg.builder.CreateBitCast(addrPtr, llvm.PointerType(cg.context.Int32Type(), 0), "addrptr32")
	ipAddr32 := cg.builder.CreateTrunc(ipAddr, cg.context.Int32Type(), "ipaddr32")
	// Byte swap: ABCD -> DCBA
	byte0 := cg.builder.CreateAnd(ipAddr32, llvm.ConstInt(cg.context.Int32Type(), 0xFF, false), "byte0")
	byte0Shifted := cg.builder.CreateShl(byte0, llvm.ConstInt(cg.context.Int32Type(), 24, false), "byte0shifted")
	byte1 := cg.builder.CreateAnd(cg.builder.CreateLShr(ipAddr32, llvm.ConstInt(cg.context.Int32Type(), 8, false), "tmp1"), llvm.ConstInt(cg.context.Int32Type(), 0xFF, false), "byte1")
	byte1Shifted := cg.builder.CreateShl(byte1, llvm.ConstInt(cg.context.Int32Type(), 16, false), "byte1shifted")
	byte2 := cg.builder.CreateAnd(cg.builder.CreateLShr(ipAddr32, llvm.ConstInt(cg.context.Int32Type(), 16, false), "tmp2"), llvm.ConstInt(cg.context.Int32Type(), 0xFF, false), "byte2")
	byte2Shifted := cg.builder.CreateShl(byte2, llvm.ConstInt(cg.context.Int32Type(), 8, false), "byte2shifted")
	byte3 := cg.builder.CreateLShr(ipAddr32, llvm.ConstInt(cg.context.Int32Type(), 24, false), "byte3")
	ipAddrBE := cg.builder.CreateOr(cg.builder.CreateOr(byte0Shifted, byte1Shifted, "tmp3"), cg.builder.CreateOr(byte2Shifted, byte3, "tmp4"), "ipaddrbe")
	cg.builder.CreateStore(ipAddrBE, addrPtr32)

	// Call connect syscall (42 on x86_64 Linux)
	syscallNum := llvm.ConstInt(cg.context.Int64Type(), 42, false)
	sockaddrInt := cg.builder.CreatePtrToInt(sockaddrPtr, cg.context.Int64Type(), "sockaddrint")
	addrLen := llvm.ConstInt(cg.context.Int64Type(), 16, false)
	syscall := cg.functions["syscall"]

	return cg.builder.CreateCall(syscall.GlobalValueType(), syscall,
		[]llvm.Value{syscallNum, fd, sockaddrInt, addrLen}, "connecttmp"), nil
}

// generateSend sends data on a connected socket
// send(fd, buf, len) -> bytes_sent
func (cg *LLVMCodeGenerator) generateSend(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, fmt.Errorf("send requires 3 arguments (fd, buf, len)")
	}

	fd, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	buf, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	length, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}

	cg.declareSyscall()

	// sendto syscall (44 on x86_64) with NULL dest_addr and 0 addrlen works as send()
	// Or we can use write() syscall (1) which works on sockets too
	syscallNum := llvm.ConstInt(cg.context.Int64Type(), 1, false) // write syscall

	var bufInt llvm.Value
	if buf.Type().TypeKind() == llvm.PointerTypeKind {
		bufInt = cg.builder.CreatePtrToInt(buf, cg.context.Int64Type(), "bufint")
	} else {
		bufInt = buf
	}

	syscall := cg.functions["syscall"]
	return cg.builder.CreateCall(syscall.GlobalValueType(), syscall,
		[]llvm.Value{syscallNum, fd, bufInt, length}, "sendtmp"), nil
}

// generateRecv receives data from a connected socket
// recv(fd, buf, maxlen) -> bytes_received
func (cg *LLVMCodeGenerator) generateRecv(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, fmt.Errorf("recv requires 3 arguments (fd, buf, maxlen)")
	}

	fd, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	buf, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	maxLen, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}

	cg.declareSyscall()

	// Use read syscall (0) which works on sockets
	syscallNum := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	syscall := cg.functions["syscall"]

	return cg.builder.CreateCall(syscall.GlobalValueType(), syscall,
		[]llvm.Value{syscallNum, fd, buf, maxLen}, "recvtmp"), nil
}

// generateBindIPv4 binds a socket to an IPv4 address and port
// bind_ipv4(fd, ip_addr, port) -> status
func (cg *LLVMCodeGenerator) generateBindIPv4(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, fmt.Errorf("bind_ipv4 requires 3 arguments (fd, ip, port)")
	}

	fd, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	ipAddr, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	port, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}

	cg.declareSyscall()

	// Build sockaddr_in structure (same as connect)
	sockaddrType := llvm.ArrayType(cg.context.Int8Type(), 16)
	sockaddr := cg.builder.CreateAlloca(sockaddrType, "sockaddr")

	memset := cg.functions["memset"]
	if memset.IsNil() {
		memsetType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0), cg.context.Int32Type(), cg.context.Int64Type()},
			false,
		)
		memset = llvm.AddFunction(cg.module, "memset", memsetType)
		cg.functions["memset"] = memset
	}
	sockaddrPtr := cg.builder.CreateBitCast(sockaddr, llvm.PointerType(cg.context.Int8Type(), 0), "sockaddrptr")
	cg.builder.CreateCall(memset.GlobalValueType(), memset,
		[]llvm.Value{sockaddrPtr, llvm.ConstInt(cg.context.Int32Type(), 0, false), llvm.ConstInt(cg.context.Int64Type(), 16, false)}, "")

	// Set sin_family = AF_INET (2)
	familyPtr := cg.builder.CreateGEP(cg.context.Int8Type(), sockaddrPtr, []llvm.Value{llvm.ConstInt(cg.context.Int32Type(), 0, false)}, "familyptr")
	familyPtr16 := cg.builder.CreateBitCast(familyPtr, llvm.PointerType(cg.context.Int16Type(), 0), "familyptr16")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int16Type(), 2, false), familyPtr16)

	// Set sin_port (network byte order)
	portPtr := cg.builder.CreateGEP(cg.context.Int8Type(), sockaddrPtr, []llvm.Value{llvm.ConstInt(cg.context.Int32Type(), 2, false)}, "portptr")
	portPtr16 := cg.builder.CreateBitCast(portPtr, llvm.PointerType(cg.context.Int16Type(), 0), "portptr16")
	port16 := cg.builder.CreateTrunc(port, cg.context.Int16Type(), "port16")
	port16Ext := cg.builder.CreateZExt(port16, cg.context.Int32Type(), "port16ext")
	shifted := cg.builder.CreateShl(port16Ext, llvm.ConstInt(cg.context.Int32Type(), 8, false), "shifted")
	unshifted := cg.builder.CreateLShr(port16Ext, llvm.ConstInt(cg.context.Int32Type(), 8, false), "unshifted")
	portBE := cg.builder.CreateOr(shifted, unshifted, "portbe")
	portBE16 := cg.builder.CreateTrunc(portBE, cg.context.Int16Type(), "portbe16")
	cg.builder.CreateStore(portBE16, portPtr16)

	// Set sin_addr (convert to network byte order - big endian)
	addrPtr := cg.builder.CreateGEP(cg.context.Int8Type(), sockaddrPtr, []llvm.Value{llvm.ConstInt(cg.context.Int32Type(), 4, false)}, "addrptr")
	addrPtr32 := cg.builder.CreateBitCast(addrPtr, llvm.PointerType(cg.context.Int32Type(), 0), "addrptr32")
	ipAddr32 := cg.builder.CreateTrunc(ipAddr, cg.context.Int32Type(), "ipaddr32")
	// Byte swap: ABCD -> DCBA
	byte0 := cg.builder.CreateAnd(ipAddr32, llvm.ConstInt(cg.context.Int32Type(), 0xFF, false), "byte0")
	byte0Shifted := cg.builder.CreateShl(byte0, llvm.ConstInt(cg.context.Int32Type(), 24, false), "byte0shifted")
	byte1 := cg.builder.CreateAnd(cg.builder.CreateLShr(ipAddr32, llvm.ConstInt(cg.context.Int32Type(), 8, false), "tmp1"), llvm.ConstInt(cg.context.Int32Type(), 0xFF, false), "byte1")
	byte1Shifted := cg.builder.CreateShl(byte1, llvm.ConstInt(cg.context.Int32Type(), 16, false), "byte1shifted")
	byte2 := cg.builder.CreateAnd(cg.builder.CreateLShr(ipAddr32, llvm.ConstInt(cg.context.Int32Type(), 16, false), "tmp2"), llvm.ConstInt(cg.context.Int32Type(), 0xFF, false), "byte2")
	byte2Shifted := cg.builder.CreateShl(byte2, llvm.ConstInt(cg.context.Int32Type(), 8, false), "byte2shifted")
	byte3 := cg.builder.CreateLShr(ipAddr32, llvm.ConstInt(cg.context.Int32Type(), 24, false), "byte3")
	ipAddrBE := cg.builder.CreateOr(cg.builder.CreateOr(byte0Shifted, byte1Shifted, "tmp3"), cg.builder.CreateOr(byte2Shifted, byte3, "tmp4"), "ipaddrbe")
	cg.builder.CreateStore(ipAddrBE, addrPtr32)

	// Call bind syscall (49 on x86_64 Linux)
	syscallNum := llvm.ConstInt(cg.context.Int64Type(), 49, false)
	sockaddrInt := cg.builder.CreatePtrToInt(sockaddrPtr, cg.context.Int64Type(), "sockaddrint")
	addrLen := llvm.ConstInt(cg.context.Int64Type(), 16, false)
	syscall := cg.functions["syscall"]

	return cg.builder.CreateCall(syscall.GlobalValueType(), syscall,
		[]llvm.Value{syscallNum, fd, sockaddrInt, addrLen}, "bindtmp"), nil
}

// generateListen marks a socket as listening for connections
// listen(fd, backlog) -> status
func (cg *LLVMCodeGenerator) generateListen(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("listen requires 2 arguments (fd, backlog)")
	}

	fd, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	backlog, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	cg.declareSyscall()

	// listen syscall (50 on x86_64 Linux)
	syscallNum := llvm.ConstInt(cg.context.Int64Type(), 50, false)
	zero := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	syscall := cg.functions["syscall"]

	return cg.builder.CreateCall(syscall.GlobalValueType(), syscall,
		[]llvm.Value{syscallNum, fd, backlog, zero}, "listentmp"), nil
}

// generateAccept accepts a connection on a listening socket
// accept(fd) -> new_fd (client socket)
func (cg *LLVMCodeGenerator) generateAccept(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("accept requires 1 argument (fd)")
	}

	fd, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	cg.declareSyscall()

	// accept syscall (43 on x86_64 Linux) with NULL addr and addrlen
	syscallNum := llvm.ConstInt(cg.context.Int64Type(), 43, false)
	zero := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	syscall := cg.functions["syscall"]

	return cg.builder.CreateCall(syscall.GlobalValueType(), syscall,
		[]llvm.Value{syscallNum, fd, zero, zero}, "accepttmp"), nil
}

// generateSetSockOpt sets socket options (useful for SO_REUSEADDR)
// setsockopt(fd, level, optname, optval) -> status
func (cg *LLVMCodeGenerator) generateSetSockOpt(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 4 {
		return llvm.Value{}, fmt.Errorf("setsockopt requires 4 arguments (fd, level, optname, optval)")
	}

	fd, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	level, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	optname, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}
	optval, err := cg.generateExpression(call.Args[3])
	if err != nil {
		return llvm.Value{}, err
	}

	// Use libc setsockopt instead of syscall (simpler)
	setsockoptFn := cg.functions["setsockopt"]
	if setsockoptFn.IsNil() {
		// int setsockopt(int sockfd, int level, int optname, const void *optval, socklen_t optlen)
		setsockoptType := llvm.FunctionType(
			cg.context.Int32Type(),
			[]llvm.Type{
				cg.context.Int32Type(),
				cg.context.Int32Type(),
				cg.context.Int32Type(),
				llvm.PointerType(cg.context.Int8Type(), 0),
				cg.context.Int32Type(),
			},
			false,
		)
		setsockoptFn = llvm.AddFunction(cg.module, "setsockopt", setsockoptType)
		setsockoptFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["setsockopt"] = setsockoptFn
	}

	// Create int on stack for optval
	optvalPtr := cg.builder.CreateAlloca(cg.context.Int32Type(), "optvalptr")
	optval32 := cg.builder.CreateTrunc(optval, cg.context.Int32Type(), "optval32")
	cg.builder.CreateStore(optval32, optvalPtr)

	fd32 := cg.builder.CreateTrunc(fd, cg.context.Int32Type(), "fd32")
	level32 := cg.builder.CreateTrunc(level, cg.context.Int32Type(), "level32")
	optname32 := cg.builder.CreateTrunc(optname, cg.context.Int32Type(), "optname32")
	optvalVoid := cg.builder.CreateBitCast(optvalPtr, llvm.PointerType(cg.context.Int8Type(), 0), "optvalvoid")
	optlen := llvm.ConstInt(cg.context.Int32Type(), 4, false) // sizeof(int)

	result := cg.builder.CreateCall(setsockoptFn.GlobalValueType(), setsockoptFn,
		[]llvm.Value{fd32, level32, optname32, optvalVoid, optlen}, "setsockopttmp")
	return cg.builder.CreateSExt(result, cg.context.Int64Type(), "setsockoptres"), nil
}

// ============================================================================
// FILE I/O FUNCTIONS
// ============================================================================

func (cg *LLVMCodeGenerator) generateOpen(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) < 2 {
		return llvm.Value{}, nil
	}

	path, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	flags, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	cg.declareSyscall()

	// Convert path pointer for syscall
	pathInt := cg.builder.CreatePtrToInt(path, cg.context.Int64Type(), "pathint")

	// open() is syscall 2 on x86_64 Linux, mode 0644 = 420
	syscallNum := llvm.ConstInt(cg.context.Int64Type(), 2, false)
	mode := llvm.ConstInt(cg.context.Int64Type(), 420, false) // 0644
	syscall := cg.functions["syscall"]

	return cg.builder.CreateCall(syscall.GlobalValueType(), syscall,
		[]llvm.Value{syscallNum, pathInt, flags, mode}, "opentmp"), nil
}

func (cg *LLVMCodeGenerator) generateRead(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, nil
	}

	fd, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	buf, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	count, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}

	cg.declareSyscall()

	// read() is syscall 0 on x86_64 Linux
	syscallNum := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	syscall := cg.functions["syscall"]

	return cg.builder.CreateCall(syscall.GlobalValueType(), syscall,
		[]llvm.Value{syscallNum, fd, buf, count}, "readtmp"), nil
}

func (cg *LLVMCodeGenerator) generateWrite(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, nil
	}

	fd, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	buf, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	count, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}

	cg.declareSyscall()

	// write() is syscall 1 on x86_64 Linux
	syscallNum := llvm.ConstInt(cg.context.Int64Type(), 1, false)

	// Convert buf (string pointer) to int for syscall
	var bufInt llvm.Value
	if buf.Type().TypeKind() == llvm.PointerTypeKind {
		bufInt = cg.builder.CreatePtrToInt(buf, cg.context.Int64Type(), "bufint")
	} else {
		bufInt = buf
	}

	syscall := cg.functions["syscall"]
	return cg.builder.CreateCall(syscall.GlobalValueType(), syscall,
		[]llvm.Value{syscallNum, fd, bufInt, count}, "writetmp"), nil
}

// ============================================================================
// TIME FUNCTIONS
// ============================================================================

func (cg *LLVMCodeGenerator) generateNow(call *FunctionCall) (llvm.Value, error) {
	cg.declareSyscall()

	// time(NULL) is syscall 201 on x86_64 Linux
	syscallNum := llvm.ConstInt(cg.context.Int64Type(), 201, false)
	zero := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	syscall := cg.functions["syscall"]

	return cg.builder.CreateCall(syscall.GlobalValueType(), syscall,
		[]llvm.Value{syscallNum, zero, zero, zero}, "nowtmp"), nil
}

func (cg *LLVMCodeGenerator) generateMillis(call *FunctionCall) (llvm.Value, error) {
	// Simplified: just return now() * 1000 (not accurate, but works for testing)
	now, err := cg.generateNow(call)
	if err != nil {
		return llvm.Value{}, err
	}
	thousand := llvm.ConstInt(cg.context.Int64Type(), 1000, false)
	return cg.builder.CreateMul(now, thousand, "millistmp"), nil
}

func (cg *LLVMCodeGenerator) generateNanos(call *FunctionCall) (llvm.Value, error) {
	// Simplified: just return now() * 1000000000
	now, err := cg.generateNow(call)
	if err != nil {
		return llvm.Value{}, err
	}
	billion := llvm.ConstInt(cg.context.Int64Type(), 1000000000, false)
	return cg.builder.CreateMul(now, billion, "nanostmp"), nil
}

func (cg *LLVMCodeGenerator) generateSleep(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	seconds, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Declare C sleep function
	if _, ok := cg.functions["sleep"]; !ok {
		sleepType := llvm.FunctionType(
			cg.context.Int32Type(),
			[]llvm.Type{cg.context.Int32Type()},
			false,
		)
		sleepFn := llvm.AddFunction(cg.module, "sleep", sleepType)
		sleepFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["sleep"] = sleepFn
	}

	sleep := cg.functions["sleep"]
	// Truncate i64 to i32
	sec32 := cg.builder.CreateTrunc(seconds, cg.context.Int32Type(), "sec32")
	result := cg.builder.CreateCall(sleep.GlobalValueType(), sleep, []llvm.Value{sec32}, "sleeptmp")
	return cg.builder.CreateZExt(result, cg.context.Int64Type(), "sleepext"), nil
}

// ============================================================================
// HTTP POOL FUNCTIONS (stub implementations)
// ============================================================================

func (cg *LLVMCodeGenerator) generatePoolNew(call *FunctionCall) (llvm.Value, error) {
	// Just allocate a simple counter for pool management
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	size, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Allocate size * 8 bytes for connection pool
	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	allocSize := cg.builder.CreateMul(size, eight, "poolsize")

	malloc := cg.functions["malloc"]
	ptr := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{allocSize}, "poolalloc")
	return cg.builder.CreatePtrToInt(ptr, cg.context.Int64Type(), "poolptr"), nil
}

func (cg *LLVMCodeGenerator) generatePoolGet(call *FunctionCall) (llvm.Value, error) {
	// Stub: return 0 as a fake connection
	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

func (cg *LLVMCodeGenerator) generatePoolPut(call *FunctionCall) (llvm.Value, error) {
	// Stub: no-op
	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

func (cg *LLVMCodeGenerator) generatePoolClose(call *FunctionCall) (llvm.Value, error) {
	// Free the pool
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	poolPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	ptr := cg.builder.CreateIntToPtr(poolPtr, llvm.PointerType(cg.context.Int8Type(), 0), "freeptr")
	free := cg.functions["free"]
	cg.builder.CreateCall(free.GlobalValueType(), free, []llvm.Value{ptr}, "")
	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

// ============================================================================
// RANDOM FUNCTIONS
// Uses xorshift64 PRNG for fast, good-quality random numbers
// ============================================================================

// declareRandomState ensures the random state global variable exists
func (cg *LLVMCodeGenerator) declareRandomState() llvm.Value {
	stateName := "__lotus_random_state"
	if state, ok := cg.globalVars[stateName]; ok {
		return state
	}

	// Create global variable for PRNG state, initialized with a default seed
	stateType := cg.context.Int64Type()
	state := llvm.AddGlobal(cg.module, stateType, stateName)
	state.SetLinkage(llvm.InternalLinkage)
	// Default seed based on a prime number
	state.SetInitializer(llvm.ConstInt(stateType, 88172645463325252, false))
	cg.globalVars[stateName] = state
	return state
}

// xorshift64 implements the xorshift64 PRNG algorithm
// state ^= state << 13
// state ^= state >> 7
// state ^= state << 17
func (cg *LLVMCodeGenerator) xorshift64(state llvm.Value) llvm.Value {
	stateVal := cg.builder.CreateLoad(cg.context.Int64Type(), state, "randstate")

	// state ^= state << 13
	shift1 := cg.builder.CreateShl(stateVal, llvm.ConstInt(cg.context.Int64Type(), 13, false), "shl13")
	xor1 := cg.builder.CreateXor(stateVal, shift1, "xor1")

	// state ^= state >> 7
	shift2 := cg.builder.CreateLShr(xor1, llvm.ConstInt(cg.context.Int64Type(), 7, false), "lshr7")
	xor2 := cg.builder.CreateXor(xor1, shift2, "xor2")

	// state ^= state << 17
	shift3 := cg.builder.CreateShl(xor2, llvm.ConstInt(cg.context.Int64Type(), 17, false), "shl17")
	xor3 := cg.builder.CreateXor(xor2, shift3, "xor3")

	// Store new state
	cg.builder.CreateStore(xor3, state)

	return xor3
}

// generateRand generates code for rand() -> random int64
func (cg *LLVMCodeGenerator) generateRand(call *FunctionCall) (llvm.Value, error) {
	state := cg.declareRandomState()
	return cg.xorshift64(state), nil
}

// generateRandRange generates code for rand_range(min, max) -> random int in [min, max]
func (cg *LLVMCodeGenerator) generateRandRange(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, nil
	}

	minVal, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	maxVal, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	state := cg.declareRandomState()
	randVal := cg.xorshift64(state)

	// range = max - min + 1
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	rangeSize := cg.builder.CreateSub(maxVal, minVal, "rangesub")
	rangeSize = cg.builder.CreateAdd(rangeSize, one, "rangeadd")

	// result = min + (rand % range)
	// Make rand positive first
	absRand := cg.builder.CreateAnd(randVal, llvm.ConstInt(cg.context.Int64Type(), 0x7FFFFFFFFFFFFFFF, false), "absrand")
	modVal := cg.builder.CreateURem(absRand, rangeSize, "randmod")
	result := cg.builder.CreateAdd(minVal, modVal, "randrange")

	return result, nil
}

// generateSeed generates code for seed(n) -> void
func (cg *LLVMCodeGenerator) generateSeed(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	seedVal, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	state := cg.declareRandomState()

	// Ensure seed is not zero (xorshift doesn't work with zero)
	zero := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	defaultSeed := llvm.ConstInt(cg.context.Int64Type(), 88172645463325252, false)
	isZero := cg.builder.CreateICmp(llvm.IntEQ, seedVal, zero, "iszero")
	finalSeed := cg.builder.CreateSelect(isZero, defaultSeed, seedVal, "finalseed")

	cg.builder.CreateStore(finalSeed, state)
	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

// generateRandFloat generates code for rand_float() -> random float [0.0, 1.0)
func (cg *LLVMCodeGenerator) generateRandFloat(call *FunctionCall) (llvm.Value, error) {
	state := cg.declareRandomState()
	randVal := cg.xorshift64(state)

	// Make positive
	absRand := cg.builder.CreateAnd(randVal, llvm.ConstInt(cg.context.Int64Type(), 0x7FFFFFFFFFFFFFFF, false), "absrand")

	// Convert to double and divide by max int64
	floatVal := cg.builder.CreateUIToFP(absRand, cg.context.DoubleType(), "randfloat")
	maxVal := llvm.ConstFloat(cg.context.DoubleType(), float64(0x7FFFFFFFFFFFFFFF))
	result := cg.builder.CreateFDiv(floatVal, maxVal, "randfloatdiv")

	// Convert back to i64 representation (as Lotus uses i64 for everything)
	return cg.builder.CreateBitCast(result, cg.context.Int64Type(), "floatbits"), nil
}

// generateRandBool generates code for rand_bool() -> random boolean (0 or 1)
func (cg *LLVMCodeGenerator) generateRandBool(call *FunctionCall) (llvm.Value, error) {
	state := cg.declareRandomState()
	randVal := cg.xorshift64(state)

	// Get least significant bit
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	return cg.builder.CreateAnd(randVal, one, "randbool"), nil
}

// generateRandBytes generates code for rand_bytes(buf, len) -> fills buffer with random bytes
func (cg *LLVMCodeGenerator) generateRandBytes(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, nil
	}

	bufPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	length, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	state := cg.declareRandomState()

	// Create loop to fill buffer
	currentFn := cg.builder.GetInsertBlock().Parent()
	loopBB := cg.context.AddBasicBlock(currentFn, "randbytes_loop")
	bodyBB := cg.context.AddBasicBlock(currentFn, "randbytes_body")
	exitBB := cg.context.AddBasicBlock(currentFn, "randbytes_exit")

	// Allocate loop counter
	counter := cg.builder.CreateAlloca(cg.context.Int64Type(), "randbytes_i")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int64Type(), 0, false), counter)

	cg.builder.CreateBr(loopBB)

	// Loop condition
	cg.builder.SetInsertPointAtEnd(loopBB)
	i := cg.builder.CreateLoad(cg.context.Int64Type(), counter, "i")
	cond := cg.builder.CreateICmp(llvm.IntULT, i, length, "loopcond")
	cg.builder.CreateCondBr(cond, bodyBB, exitBB)

	// Loop body
	cg.builder.SetInsertPointAtEnd(bodyBB)
	randVal := cg.xorshift64(state)
	randByte := cg.builder.CreateTrunc(randVal, cg.context.Int8Type(), "randbyte")

	// Store byte at buf[i]
	buf := cg.builder.CreateIntToPtr(bufPtr, llvm.PointerType(cg.context.Int8Type(), 0), "bufptr")
	elemPtr := cg.builder.CreateGEP(cg.context.Int8Type(), buf, []llvm.Value{i}, "elemptr")
	cg.builder.CreateStore(randByte, elemPtr)

	// Increment counter
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	nextI := cg.builder.CreateAdd(i, one, "nexti")
	cg.builder.CreateStore(nextI, counter)
	cg.builder.CreateBr(loopBB)

	// Exit
	cg.builder.SetInsertPointAtEnd(exitBB)
	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

// generateShuffle generates code for shuffle(arr, len) -> shuffles int array in place (Fisher-Yates)
func (cg *LLVMCodeGenerator) generateShuffle(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, nil
	}

	arrPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	length, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	state := cg.declareRandomState()

	// Fisher-Yates shuffle
	currentFn := cg.builder.GetInsertBlock().Parent()
	loopBB := cg.context.AddBasicBlock(currentFn, "shuffle_loop")
	bodyBB := cg.context.AddBasicBlock(currentFn, "shuffle_body")
	exitBB := cg.context.AddBasicBlock(currentFn, "shuffle_exit")

	// i = length - 1
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	counter := cg.builder.CreateAlloca(cg.context.Int64Type(), "shuffle_i")
	startI := cg.builder.CreateSub(length, one, "starti")
	cg.builder.CreateStore(startI, counter)

	cg.builder.CreateBr(loopBB)

	// Loop condition: i > 0
	cg.builder.SetInsertPointAtEnd(loopBB)
	i := cg.builder.CreateLoad(cg.context.Int64Type(), counter, "i")
	zero := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	cond := cg.builder.CreateICmp(llvm.IntSGT, i, zero, "loopcond")
	cg.builder.CreateCondBr(cond, bodyBB, exitBB)

	// Loop body
	cg.builder.SetInsertPointAtEnd(bodyBB)

	// j = rand() % (i + 1)
	randVal := cg.xorshift64(state)
	absRand := cg.builder.CreateAnd(randVal, llvm.ConstInt(cg.context.Int64Type(), 0x7FFFFFFFFFFFFFFF, false), "absrand")
	iPlusOne := cg.builder.CreateAdd(i, one, "iplus1")
	j := cg.builder.CreateURem(absRand, iPlusOne, "j")

	// Swap arr[i] and arr[j]
	arr := cg.builder.CreateIntToPtr(arrPtr, llvm.PointerType(cg.context.Int64Type(), 0), "arrptr")

	// Get arr[i]
	iOffset := cg.builder.CreateMul(i, eight, "ioffset")
	iPtrInt := cg.builder.CreateAdd(arrPtr, iOffset, "iptrint")
	iPtr := cg.builder.CreateIntToPtr(iPtrInt, llvm.PointerType(cg.context.Int64Type(), 0), "iptr")
	valI := cg.builder.CreateLoad(cg.context.Int64Type(), iPtr, "vali")

	// Get arr[j]
	jOffset := cg.builder.CreateMul(j, eight, "joffset")
	jPtrInt := cg.builder.CreateAdd(arrPtr, jOffset, "jptrint")
	jPtr := cg.builder.CreateIntToPtr(jPtrInt, llvm.PointerType(cg.context.Int64Type(), 0), "jptr")
	valJ := cg.builder.CreateLoad(cg.context.Int64Type(), jPtr, "valj")

	// Swap
	cg.builder.CreateStore(valJ, iPtr)
	cg.builder.CreateStore(valI, jPtr)

	// Decrement i
	nextI := cg.builder.CreateSub(i, one, "nexti")
	cg.builder.CreateStore(nextI, counter)
	cg.builder.CreateBr(loopBB)

	// Exit
	cg.builder.SetInsertPointAtEnd(exitBB)

	// Suppress unused variable warning
	_ = arr

	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

// generateChoice generates code for choice(arr, len) -> random element from int array
func (cg *LLVMCodeGenerator) generateChoice(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, nil
	}

	arrPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	length, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	state := cg.declareRandomState()
	randVal := cg.xorshift64(state)

	// index = abs(rand) % len
	absRand := cg.builder.CreateAnd(randVal, llvm.ConstInt(cg.context.Int64Type(), 0x7FFFFFFFFFFFFFFF, false), "absrand")
	index := cg.builder.CreateURem(absRand, length, "choiceindex")

	// Get arr[index]
	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	offset := cg.builder.CreateMul(index, eight, "offset")
	elemPtrInt := cg.builder.CreateAdd(arrPtr, offset, "elemptrint")
	elemPtr := cg.builder.CreateIntToPtr(elemPtrInt, llvm.PointerType(cg.context.Int64Type(), 0), "elemptr")

	return cg.builder.CreateLoad(cg.context.Int64Type(), elemPtr, "choiceval"), nil
}

// generateRandN generates code for rand_n(n) -> random int in [0, n)
func (cg *LLVMCodeGenerator) generateRandN(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	n, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	state := cg.declareRandomState()
	randVal := cg.xorshift64(state)

	// result = abs(rand) % n
	absRand := cg.builder.CreateAnd(randVal, llvm.ConstInt(cg.context.Int64Type(), 0x7FFFFFFFFFFFFFFF, false), "absrand")
	return cg.builder.CreateURem(absRand, n, "randn"), nil
}

// generateRandString generates code for rand_string(buf, len) -> random alphanumeric string
func (cg *LLVMCodeGenerator) generateRandString(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, nil
	}

	bufPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	length, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	state := cg.declareRandomState()

	// Alphanumeric characters: 0-9, A-Z, a-z (62 chars)
	// We'll use: index % 62 to pick a character

	currentFn := cg.builder.GetInsertBlock().Parent()
	loopBB := cg.context.AddBasicBlock(currentFn, "randstr_loop")
	bodyBB := cg.context.AddBasicBlock(currentFn, "randstr_body")
	exitBB := cg.context.AddBasicBlock(currentFn, "randstr_exit")

	counter := cg.builder.CreateAlloca(cg.context.Int64Type(), "randstr_i")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int64Type(), 0, false), counter)

	cg.builder.CreateBr(loopBB)

	// Loop condition
	cg.builder.SetInsertPointAtEnd(loopBB)
	i := cg.builder.CreateLoad(cg.context.Int64Type(), counter, "i")
	cond := cg.builder.CreateICmp(llvm.IntULT, i, length, "loopcond")
	cg.builder.CreateCondBr(cond, bodyBB, exitBB)

	// Loop body
	cg.builder.SetInsertPointAtEnd(bodyBB)
	randVal := cg.xorshift64(state)
	absRand := cg.builder.CreateAnd(randVal, llvm.ConstInt(cg.context.Int64Type(), 0x7FFFFFFFFFFFFFFF, false), "absrand")

	// charIndex = rand % 62
	sixtytwo := llvm.ConstInt(cg.context.Int64Type(), 62, false)
	charIndex := cg.builder.CreateURem(absRand, sixtytwo, "charindex")

	// Convert index to character:
	// 0-9 -> '0' (48) + index
	// 10-35 -> 'A' (65) + (index - 10)
	// 36-61 -> 'a' (97) + (index - 36)

	ten := llvm.ConstInt(cg.context.Int64Type(), 10, false)
	thirtysix := llvm.ConstInt(cg.context.Int64Type(), 36, false)
	fortyeight := llvm.ConstInt(cg.context.Int64Type(), 48, false) // '0'
	fiftyfive := llvm.ConstInt(cg.context.Int64Type(), 55, false)  // 'A' - 10
	sixtyone := llvm.ConstInt(cg.context.Int64Type(), 61, false)   // 'a' - 36

	// Check which range
	isDigit := cg.builder.CreateICmp(llvm.IntULT, charIndex, ten, "isdigit")
	isUpper := cg.builder.CreateICmp(llvm.IntULT, charIndex, thirtysix, "isupper")

	// Compute character for each case
	digitChar := cg.builder.CreateAdd(charIndex, fortyeight, "digitchar")
	upperChar := cg.builder.CreateAdd(charIndex, fiftyfive, "upperchar")
	lowerChar := cg.builder.CreateAdd(charIndex, sixtyone, "lowerchar")

	// Select based on range
	upperOrLower := cg.builder.CreateSelect(isUpper, upperChar, lowerChar, "upperorlower")
	finalChar := cg.builder.CreateSelect(isDigit, digitChar, upperOrLower, "finalchar")
	charByte := cg.builder.CreateTrunc(finalChar, cg.context.Int8Type(), "charbyte")

	// Store char at buf[i]
	buf := cg.builder.CreateIntToPtr(bufPtr, llvm.PointerType(cg.context.Int8Type(), 0), "bufptr")
	elemPtr := cg.builder.CreateGEP(cg.context.Int8Type(), buf, []llvm.Value{i}, "elemptr")
	cg.builder.CreateStore(charByte, elemPtr)

	// Increment counter
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	nextI := cg.builder.CreateAdd(i, one, "nexti")
	cg.builder.CreateStore(nextI, counter)
	cg.builder.CreateBr(loopBB)

	// Exit - add null terminator
	cg.builder.SetInsertPointAtEnd(exitBB)
	buf2 := cg.builder.CreateIntToPtr(bufPtr, llvm.PointerType(cg.context.Int8Type(), 0), "bufptr2")
	nullPtr := cg.builder.CreateGEP(cg.context.Int8Type(), buf2, []llvm.Value{length}, "nullptr")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int8Type(), 0, false), nullPtr)

	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

// ============================================================================
// TYPE CONVERSION FUNCTIONS
// ============================================================================

// generateToUint32 converts an int64 to uint32 (truncation)
func (cg *LLVMCodeGenerator) generateToUint32(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("toUint32 requires 1 argument")
	}

	val, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Truncate to 32-bit, then zero-extend back to 64-bit for consistency
	truncated := cg.builder.CreateTrunc(val, cg.context.Int32Type(), "touint32_trunc")
	extended := cg.builder.CreateZExt(truncated, cg.context.Int64Type(), "touint32_ext")
	return extended, nil
}

// generateToBool converts a value to boolean (0 = false, non-zero = true)
func (cg *LLVMCodeGenerator) generateToBool(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("toBool requires 1 argument")
	}

	val, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Compare with 0, return 1 if non-zero, 0 if zero
	zero := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	cmp := cg.builder.CreateICmp(llvm.IntNE, val, zero, "tobool_cmp")
	return cg.builder.CreateZExt(cmp, cg.context.Int64Type(), "tobool_result"), nil
}

// generateToInt converts a float to int64
func (cg *LLVMCodeGenerator) generateToInt(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("toInt requires 1 argument")
	}

	val, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// If already int, return as-is. If float, convert.
	if val.Type().TypeKind() == llvm.DoubleTypeKind {
		return cg.builder.CreateFPToSI(val, cg.context.Int64Type(), "toint_conv"), nil
	}
	// Already an integer
	return val, nil
}

// generateToFloat converts an int to float64
func (cg *LLVMCodeGenerator) generateToFloat(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("toFloat requires 1 argument")
	}

	val, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// If already float, return as-is. If int, convert.
	if val.Type().TypeKind() == llvm.IntegerTypeKind {
		return cg.builder.CreateSIToFP(val, cg.context.DoubleType(), "tofloat_conv"), nil
	}
	// Already a float
	return val, nil
}

// ============================================================================
// STRING MANIPULATION FUNCTIONS
// ============================================================================

// generateStrCopy creates a copy of a string using malloc and strcpy
func (cg *LLVMCodeGenerator) generateStrCopy(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("copy requires 1 argument for strings")
	}

	str, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Get string length
	strlen := cg.functions["strlen"]
	if strlen.IsNil() {
		strlenType := llvm.FunctionType(
			cg.context.Int64Type(),
			[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0)},
			false,
		)
		strlen = llvm.AddFunction(cg.module, "strlen", strlenType)
		strlen.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strlen"] = strlen
	}

	length := cg.builder.CreateCall(strlen.GlobalValueType(), strlen, []llvm.Value{str}, "copylen")

	// Allocate len + 1 bytes for null terminator
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	allocSize := cg.builder.CreateAdd(length, one, "allocsize")

	// Declare malloc if needed
	malloc := cg.functions["malloc"]
	if malloc.IsNil() {
		mallocType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{cg.context.Int64Type()},
			false,
		)
		malloc = llvm.AddFunction(cg.module, "malloc", mallocType)
		malloc.SetLinkage(llvm.ExternalLinkage)
		cg.functions["malloc"] = malloc
	}

	newStr := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{allocSize}, "copybuf")

	// Declare strcpy if needed
	strcpy := cg.functions["strcpy"]
	if strcpy.IsNil() {
		strcpyType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0),
				llvm.PointerType(cg.context.Int8Type(), 0),
			},
			false,
		)
		strcpy = llvm.AddFunction(cg.module, "strcpy", strcpyType)
		strcpy.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strcpy"] = strcpy
	}

	cg.builder.CreateCall(strcpy.GlobalValueType(), strcpy, []llvm.Value{newStr, str}, "strcpytmp")

	return newStr, nil
}

// generateStrCompare compares two strings using strcmp
func (cg *LLVMCodeGenerator) generateStrCompare(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("compare requires 2 arguments")
	}

	str1, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	str2, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Declare strcmp if needed
	strcmp := cg.functions["strcmp"]
	if strcmp.IsNil() {
		strcmpType := llvm.FunctionType(
			cg.context.Int32Type(),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0),
				llvm.PointerType(cg.context.Int8Type(), 0),
			},
			false,
		)
		strcmp = llvm.AddFunction(cg.module, "strcmp", strcmpType)
		strcmp.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strcmp"] = strcmp
	}

	result := cg.builder.CreateCall(strcmp.GlobalValueType(), strcmp, []llvm.Value{str1, str2}, "cmptmp")
	// Sign-extend to i64 for consistency
	return cg.builder.CreateSExt(result, cg.context.Int64Type(), "cmpext"), nil
}

// generateIndexOf finds the first occurrence of needle in haystack using strstr
func (cg *LLVMCodeGenerator) generateIndexOf(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("indexOf requires 2 arguments")
	}

	haystack, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	needle, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Declare strstr if needed
	strstr := cg.functions["strstr"]
	if strstr.IsNil() {
		strstrType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0),
				llvm.PointerType(cg.context.Int8Type(), 0),
			},
			false,
		)
		strstr = llvm.AddFunction(cg.module, "strstr", strstrType)
		strstr.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strstr"] = strstr
	}

	found := cg.builder.CreateCall(strstr.GlobalValueType(), strstr, []llvm.Value{haystack, needle}, "strstr_res")

	// If found is null, return -1, else return (found - haystack)
	currentFn := cg.builder.GetInsertBlock().Parent()
	foundBB := cg.context.AddBasicBlock(currentFn, "indexof_found")
	notFoundBB := cg.context.AddBasicBlock(currentFn, "indexof_notfound")
	mergeBB := cg.context.AddBasicBlock(currentFn, "indexof_merge")

	// Check if found is null
	nullPtr := llvm.ConstNull(llvm.PointerType(cg.context.Int8Type(), 0))
	isNull := cg.builder.CreateICmp(llvm.IntEQ, found, nullPtr, "isnull")
	cg.builder.CreateCondBr(isNull, notFoundBB, foundBB)

	// Not found: return -1
	cg.builder.SetInsertPointAtEnd(notFoundBB)
	negOne := llvm.ConstInt(cg.context.Int64Type(), 0xFFFFFFFFFFFFFFFF, false) // -1
	cg.builder.CreateBr(mergeBB)

	// Found: return (found - haystack)
	cg.builder.SetInsertPointAtEnd(foundBB)
	haystackInt := cg.builder.CreatePtrToInt(haystack, cg.context.Int64Type(), "haystackint")
	foundInt := cg.builder.CreatePtrToInt(found, cg.context.Int64Type(), "foundint")
	offset := cg.builder.CreateSub(foundInt, haystackInt, "offset")
	cg.builder.CreateBr(mergeBB)

	// Merge
	cg.builder.SetInsertPointAtEnd(mergeBB)
	phi := cg.builder.CreatePHI(cg.context.Int64Type(), "indexof_result")
	phi.AddIncoming([]llvm.Value{negOne, offset}, []llvm.BasicBlock{notFoundBB, foundBB})

	return phi, nil
}

// generateSubstring extracts a substring from a string
func (cg *LLVMCodeGenerator) generateSubstring(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, fmt.Errorf("substring requires 3 arguments (str, start, length)")
	}

	str, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	start, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	length, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}

	// Allocate length + 1 bytes
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	allocSize := cg.builder.CreateAdd(length, one, "substralloc")

	malloc := cg.functions["malloc"]
	if malloc.IsNil() {
		mallocType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{cg.context.Int64Type()},
			false,
		)
		malloc = llvm.AddFunction(cg.module, "malloc", mallocType)
		malloc.SetLinkage(llvm.ExternalLinkage)
		cg.functions["malloc"] = malloc
	}

	dest := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{allocSize}, "substrbuf")

	// Get pointer to str + start
	src := cg.builder.CreateGEP(cg.context.Int8Type(), str, []llvm.Value{start}, "substrsrc")

	// Declare memcpy
	memcpy := cg.functions["memcpy"]
	if memcpy.IsNil() {
		memcpyType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0),
				llvm.PointerType(cg.context.Int8Type(), 0),
				cg.context.Int64Type(),
			},
			false,
		)
		memcpy = llvm.AddFunction(cg.module, "memcpy", memcpyType)
		memcpy.SetLinkage(llvm.ExternalLinkage)
		cg.functions["memcpy"] = memcpy
	}

	cg.builder.CreateCall(memcpy.GlobalValueType(), memcpy, []llvm.Value{dest, src, length}, "memcpytmp")

	// Add null terminator
	nullPos := cg.builder.CreateGEP(cg.context.Int8Type(), dest, []llvm.Value{length}, "nullpos")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int8Type(), 0, false), nullPos)

	return dest, nil
}

// generateToUpper converts a string to uppercase
func (cg *LLVMCodeGenerator) generateToUpper(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("toUpper requires 1 argument")
	}

	str, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Get length and allocate new buffer
	strlen := cg.functions["strlen"]
	if strlen.IsNil() {
		strlenType := llvm.FunctionType(
			cg.context.Int64Type(),
			[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0)},
			false,
		)
		strlen = llvm.AddFunction(cg.module, "strlen", strlenType)
		strlen.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strlen"] = strlen
	}
	length := cg.builder.CreateCall(strlen.GlobalValueType(), strlen, []llvm.Value{str}, "upperlen")

	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	allocSize := cg.builder.CreateAdd(length, one, "upperalloc")

	malloc := cg.functions["malloc"]
	if malloc.IsNil() {
		mallocType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{cg.context.Int64Type()},
			false,
		)
		malloc = llvm.AddFunction(cg.module, "malloc", mallocType)
		malloc.SetLinkage(llvm.ExternalLinkage)
		cg.functions["malloc"] = malloc
	}

	dest := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{allocSize}, "upperbuf")

	// Copy and convert to uppercase in a loop
	currentFn := cg.builder.GetInsertBlock().Parent()
	loopBB := cg.context.AddBasicBlock(currentFn, "upper_loop")
	bodyBB := cg.context.AddBasicBlock(currentFn, "upper_body")
	exitBB := cg.context.AddBasicBlock(currentFn, "upper_exit")

	counter := cg.builder.CreateAlloca(cg.context.Int64Type(), "upper_i")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int64Type(), 0, false), counter)
	cg.builder.CreateBr(loopBB)

	cg.builder.SetInsertPointAtEnd(loopBB)
	i := cg.builder.CreateLoad(cg.context.Int64Type(), counter, "i")
	cond := cg.builder.CreateICmp(llvm.IntULT, i, length, "uppercond")
	cg.builder.CreateCondBr(cond, bodyBB, exitBB)

	cg.builder.SetInsertPointAtEnd(bodyBB)
	srcPtr := cg.builder.CreateGEP(cg.context.Int8Type(), str, []llvm.Value{i}, "srcptr")
	ch := cg.builder.CreateLoad(cg.context.Int8Type(), srcPtr, "ch")

	// If 'a' <= ch <= 'z', subtract 32
	chInt := cg.builder.CreateZExt(ch, cg.context.Int64Type(), "chint")
	isLowerA := cg.builder.CreateICmp(llvm.IntUGE, chInt, llvm.ConstInt(cg.context.Int64Type(), 97, false), "islowera")
	isLowerZ := cg.builder.CreateICmp(llvm.IntULE, chInt, llvm.ConstInt(cg.context.Int64Type(), 122, false), "islowerz")
	isLower := cg.builder.CreateAnd(isLowerA, isLowerZ, "islower")

	upperInt := cg.builder.CreateSub(chInt, llvm.ConstInt(cg.context.Int64Type(), 32, false), "upperint")
	resultInt := cg.builder.CreateSelect(isLower, upperInt, chInt, "resultint")
	resultCh := cg.builder.CreateTrunc(resultInt, cg.context.Int8Type(), "resultch")

	destPtr := cg.builder.CreateGEP(cg.context.Int8Type(), dest, []llvm.Value{i}, "destptr")
	cg.builder.CreateStore(resultCh, destPtr)

	nextI := cg.builder.CreateAdd(i, one, "nexti")
	cg.builder.CreateStore(nextI, counter)
	cg.builder.CreateBr(loopBB)

	cg.builder.SetInsertPointAtEnd(exitBB)
	nullPos := cg.builder.CreateGEP(cg.context.Int8Type(), dest, []llvm.Value{length}, "nullpos")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int8Type(), 0, false), nullPos)

	return dest, nil
}

// generateToLower converts a string to lowercase
func (cg *LLVMCodeGenerator) generateToLower(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("toLower requires 1 argument")
	}

	str, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Get length and allocate new buffer
	strlen := cg.functions["strlen"]
	if strlen.IsNil() {
		strlenType := llvm.FunctionType(
			cg.context.Int64Type(),
			[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0)},
			false,
		)
		strlen = llvm.AddFunction(cg.module, "strlen", strlenType)
		strlen.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strlen"] = strlen
	}
	length := cg.builder.CreateCall(strlen.GlobalValueType(), strlen, []llvm.Value{str}, "lowerlen")

	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	allocSize := cg.builder.CreateAdd(length, one, "loweralloc")

	malloc := cg.functions["malloc"]
	if malloc.IsNil() {
		mallocType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{cg.context.Int64Type()},
			false,
		)
		malloc = llvm.AddFunction(cg.module, "malloc", mallocType)
		malloc.SetLinkage(llvm.ExternalLinkage)
		cg.functions["malloc"] = malloc
	}

	dest := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{allocSize}, "lowerbuf")

	// Copy and convert to lowercase in a loop
	currentFn := cg.builder.GetInsertBlock().Parent()
	loopBB := cg.context.AddBasicBlock(currentFn, "lower_loop")
	bodyBB := cg.context.AddBasicBlock(currentFn, "lower_body")
	exitBB := cg.context.AddBasicBlock(currentFn, "lower_exit")

	counter := cg.builder.CreateAlloca(cg.context.Int64Type(), "lower_i")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int64Type(), 0, false), counter)
	cg.builder.CreateBr(loopBB)

	cg.builder.SetInsertPointAtEnd(loopBB)
	i := cg.builder.CreateLoad(cg.context.Int64Type(), counter, "i")
	cond := cg.builder.CreateICmp(llvm.IntULT, i, length, "lowercond")
	cg.builder.CreateCondBr(cond, bodyBB, exitBB)

	cg.builder.SetInsertPointAtEnd(bodyBB)
	srcPtr := cg.builder.CreateGEP(cg.context.Int8Type(), str, []llvm.Value{i}, "srcptr")
	ch := cg.builder.CreateLoad(cg.context.Int8Type(), srcPtr, "ch")

	// If 'A' <= ch <= 'Z', add 32
	chInt := cg.builder.CreateZExt(ch, cg.context.Int64Type(), "chint")
	isUpperA := cg.builder.CreateICmp(llvm.IntUGE, chInt, llvm.ConstInt(cg.context.Int64Type(), 65, false), "isuppera")
	isUpperZ := cg.builder.CreateICmp(llvm.IntULE, chInt, llvm.ConstInt(cg.context.Int64Type(), 90, false), "isupperz")
	isUpper := cg.builder.CreateAnd(isUpperA, isUpperZ, "isupper")

	lowerInt := cg.builder.CreateAdd(chInt, llvm.ConstInt(cg.context.Int64Type(), 32, false), "lowerint")
	resultInt := cg.builder.CreateSelect(isUpper, lowerInt, chInt, "resultint")
	resultCh := cg.builder.CreateTrunc(resultInt, cg.context.Int8Type(), "resultch")

	destPtr := cg.builder.CreateGEP(cg.context.Int8Type(), dest, []llvm.Value{i}, "destptr")
	cg.builder.CreateStore(resultCh, destPtr)

	nextI := cg.builder.CreateAdd(i, one, "nexti")
	cg.builder.CreateStore(nextI, counter)
	cg.builder.CreateBr(loopBB)

	cg.builder.SetInsertPointAtEnd(exitBB)
	nullPos := cg.builder.CreateGEP(cg.context.Int8Type(), dest, []llvm.Value{length}, "nullpos")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int8Type(), 0, false), nullPos)

	return dest, nil
}

// generateTrim removes leading and trailing whitespace from a string
func (cg *LLVMCodeGenerator) generateTrim(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("trim requires 1 argument")
	}

	str, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Get length
	strlen := cg.functions["strlen"]
	if strlen.IsNil() {
		strlenType := llvm.FunctionType(
			cg.context.Int64Type(),
			[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0)},
			false,
		)
		strlen = llvm.AddFunction(cg.module, "strlen", strlenType)
		strlen.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strlen"] = strlen
	}
	length := cg.builder.CreateCall(strlen.GlobalValueType(), strlen, []llvm.Value{str}, "trimlen")

	// We'll find start and end indices
	currentFn := cg.builder.GetInsertBlock().Parent()

	// Find start (skip leading whitespace)
	startLoopBB := cg.context.AddBasicBlock(currentFn, "trim_start_loop")
	startBodyBB := cg.context.AddBasicBlock(currentFn, "trim_start_body")
	startExitBB := cg.context.AddBasicBlock(currentFn, "trim_start_exit")

	startIdx := cg.builder.CreateAlloca(cg.context.Int64Type(), "start_idx")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int64Type(), 0, false), startIdx)
	cg.builder.CreateBr(startLoopBB)

	cg.builder.SetInsertPointAtEnd(startLoopBB)
	startI := cg.builder.CreateLoad(cg.context.Int64Type(), startIdx, "starti")
	inBounds := cg.builder.CreateICmp(llvm.IntULT, startI, length, "startinbounds")
	cg.builder.CreateCondBr(inBounds, startBodyBB, startExitBB)

	cg.builder.SetInsertPointAtEnd(startBodyBB)
	charPtr := cg.builder.CreateGEP(cg.context.Int8Type(), str, []llvm.Value{startI}, "startcharptr")
	ch := cg.builder.CreateLoad(cg.context.Int8Type(), charPtr, "startch")
	chInt := cg.builder.CreateZExt(ch, cg.context.Int64Type(), "startchint")

	// Whitespace: space(32), tab(9), newline(10), carriage return(13)
	isSpace := cg.builder.CreateICmp(llvm.IntEQ, chInt, llvm.ConstInt(cg.context.Int64Type(), 32, false), "isspace")
	isTab := cg.builder.CreateICmp(llvm.IntEQ, chInt, llvm.ConstInt(cg.context.Int64Type(), 9, false), "istab")
	isNewline := cg.builder.CreateICmp(llvm.IntEQ, chInt, llvm.ConstInt(cg.context.Int64Type(), 10, false), "isnl")
	isCR := cg.builder.CreateICmp(llvm.IntEQ, chInt, llvm.ConstInt(cg.context.Int64Type(), 13, false), "iscr")

	ws1 := cg.builder.CreateOr(isSpace, isTab, "ws1")
	ws2 := cg.builder.CreateOr(isNewline, isCR, "ws2")
	isWS := cg.builder.CreateOr(ws1, ws2, "isws")

	// Only increment and continue if it's whitespace
	startIncrBB := cg.context.AddBasicBlock(currentFn, "trim_start_incr")
	cg.builder.CreateCondBr(isWS, startIncrBB, startExitBB)

	cg.builder.SetInsertPointAtEnd(startIncrBB)
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	nextStart := cg.builder.CreateAdd(startI, one, "nextstart")
	cg.builder.CreateStore(nextStart, startIdx)
	cg.builder.CreateBr(startLoopBB)

	cg.builder.SetInsertPointAtEnd(startExitBB)
	finalStart := cg.builder.CreateLoad(cg.context.Int64Type(), startIdx, "finalstart")

	// Find end (skip trailing whitespace)
	endLoopBB := cg.context.AddBasicBlock(currentFn, "trim_end_loop")
	endBodyBB := cg.context.AddBasicBlock(currentFn, "trim_end_body")
	endExitBB := cg.context.AddBasicBlock(currentFn, "trim_end_exit")

	endIdx := cg.builder.CreateAlloca(cg.context.Int64Type(), "end_idx")
	cg.builder.CreateStore(length, endIdx)
	cg.builder.CreateBr(endLoopBB)

	cg.builder.SetInsertPointAtEnd(endLoopBB)
	endI := cg.builder.CreateLoad(cg.context.Int64Type(), endIdx, "endi")
	gtStart := cg.builder.CreateICmp(llvm.IntUGT, endI, finalStart, "gtstart")
	cg.builder.CreateCondBr(gtStart, endBodyBB, endExitBB)

	cg.builder.SetInsertPointAtEnd(endBodyBB)
	prevEnd := cg.builder.CreateSub(endI, one, "prevend")
	endCharPtr := cg.builder.CreateGEP(cg.context.Int8Type(), str, []llvm.Value{prevEnd}, "endcharptr")
	endCh := cg.builder.CreateLoad(cg.context.Int8Type(), endCharPtr, "endch")
	endChInt := cg.builder.CreateZExt(endCh, cg.context.Int64Type(), "endchint")

	isSpace2 := cg.builder.CreateICmp(llvm.IntEQ, endChInt, llvm.ConstInt(cg.context.Int64Type(), 32, false), "isspace2")
	isTab2 := cg.builder.CreateICmp(llvm.IntEQ, endChInt, llvm.ConstInt(cg.context.Int64Type(), 9, false), "istab2")
	isNewline2 := cg.builder.CreateICmp(llvm.IntEQ, endChInt, llvm.ConstInt(cg.context.Int64Type(), 10, false), "isnl2")
	isCR2 := cg.builder.CreateICmp(llvm.IntEQ, endChInt, llvm.ConstInt(cg.context.Int64Type(), 13, false), "iscr2")

	ws3 := cg.builder.CreateOr(isSpace2, isTab2, "ws3")
	ws4 := cg.builder.CreateOr(isNewline2, isCR2, "ws4")
	isWS2 := cg.builder.CreateOr(ws3, ws4, "isws2")

	// Only decrement if it's whitespace, otherwise exit with current endI
	endDecrBB := cg.context.AddBasicBlock(currentFn, "trim_end_decr")
	cg.builder.CreateCondBr(isWS2, endDecrBB, endExitBB)

	cg.builder.SetInsertPointAtEnd(endDecrBB)
	cg.builder.CreateStore(prevEnd, endIdx)
	cg.builder.CreateBr(endLoopBB)

	cg.builder.SetInsertPointAtEnd(endExitBB)
	finalEnd := cg.builder.CreateLoad(cg.context.Int64Type(), endIdx, "finalend")

	// Calculate new length and copy substring
	newLen := cg.builder.CreateSub(finalEnd, finalStart, "newlen")
	allocSize := cg.builder.CreateAdd(newLen, one, "trimalloc")

	malloc := cg.functions["malloc"]
	if malloc.IsNil() {
		mallocType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{cg.context.Int64Type()},
			false,
		)
		malloc = llvm.AddFunction(cg.module, "malloc", mallocType)
		malloc.SetLinkage(llvm.ExternalLinkage)
		cg.functions["malloc"] = malloc
	}

	dest := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{allocSize}, "trimbuf")
	src := cg.builder.CreateGEP(cg.context.Int8Type(), str, []llvm.Value{finalStart}, "trimsrc")

	memcpy := cg.functions["memcpy"]
	if memcpy.IsNil() {
		memcpyType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0),
				llvm.PointerType(cg.context.Int8Type(), 0),
				cg.context.Int64Type(),
			},
			false,
		)
		memcpy = llvm.AddFunction(cg.module, "memcpy", memcpyType)
		memcpy.SetLinkage(llvm.ExternalLinkage)
		cg.functions["memcpy"] = memcpy
	}

	cg.builder.CreateCall(memcpy.GlobalValueType(), memcpy, []llvm.Value{dest, src, newLen}, "trimmemcpy")

	nullPos := cg.builder.CreateGEP(cg.context.Int8Type(), dest, []llvm.Value{newLen}, "trimnullpos")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int8Type(), 0, false), nullPos)

	return dest, nil
}

// generateStartsWith checks if a string starts with a prefix
func (cg *LLVMCodeGenerator) generateStartsWith(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("startsWith requires 2 arguments")
	}

	str, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	prefix, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Get prefix length
	strlen := cg.functions["strlen"]
	if strlen.IsNil() {
		strlenType := llvm.FunctionType(
			cg.context.Int64Type(),
			[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0)},
			false,
		)
		strlen = llvm.AddFunction(cg.module, "strlen", strlenType)
		strlen.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strlen"] = strlen
	}
	prefixLen := cg.builder.CreateCall(strlen.GlobalValueType(), strlen, []llvm.Value{prefix}, "prefixlen")

	// Use strncmp
	strncmp := cg.functions["strncmp"]
	if strncmp.IsNil() {
		strncmpType := llvm.FunctionType(
			cg.context.Int32Type(),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0),
				llvm.PointerType(cg.context.Int8Type(), 0),
				cg.context.Int64Type(),
			},
			false,
		)
		strncmp = llvm.AddFunction(cg.module, "strncmp", strncmpType)
		strncmp.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strncmp"] = strncmp
	}

	result := cg.builder.CreateCall(strncmp.GlobalValueType(), strncmp, []llvm.Value{str, prefix, prefixLen}, "strncmptmp")
	zero := llvm.ConstInt(cg.context.Int32Type(), 0, false)
	isEqual := cg.builder.CreateICmp(llvm.IntEQ, result, zero, "startswith_cmp")
	return cg.builder.CreateZExt(isEqual, cg.context.Int64Type(), "startswith_result"), nil
}

// generateEndsWith checks if a string ends with a suffix
func (cg *LLVMCodeGenerator) generateEndsWith(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("endsWith requires 2 arguments")
	}

	str, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	suffix, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	strlen := cg.functions["strlen"]
	if strlen.IsNil() {
		strlenType := llvm.FunctionType(
			cg.context.Int64Type(),
			[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0)},
			false,
		)
		strlen = llvm.AddFunction(cg.module, "strlen", strlenType)
		strlen.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strlen"] = strlen
	}

	strLen := cg.builder.CreateCall(strlen.GlobalValueType(), strlen, []llvm.Value{str}, "strlen")
	suffixLen := cg.builder.CreateCall(strlen.GlobalValueType(), strlen, []llvm.Value{suffix}, "suffixlen")

	currentFn := cg.builder.GetInsertBlock().Parent()
	checkBB := cg.context.AddBasicBlock(currentFn, "endswith_check")
	falseBB := cg.context.AddBasicBlock(currentFn, "endswith_false")
	mergeBB := cg.context.AddBasicBlock(currentFn, "endswith_merge")

	// If suffixLen > strLen, return false
	canCheck := cg.builder.CreateICmp(llvm.IntULE, suffixLen, strLen, "cancheck")
	cg.builder.CreateCondBr(canCheck, checkBB, falseBB)

	cg.builder.SetInsertPointAtEnd(falseBB)
	cg.builder.CreateBr(mergeBB)

	cg.builder.SetInsertPointAtEnd(checkBB)
	// Get pointer to str + (strLen - suffixLen)
	offset := cg.builder.CreateSub(strLen, suffixLen, "endswithoffset")
	strEnd := cg.builder.CreateGEP(cg.context.Int8Type(), str, []llvm.Value{offset}, "strend")

	// Compare using strcmp
	strcmp := cg.functions["strcmp"]
	if strcmp.IsNil() {
		strcmpType := llvm.FunctionType(
			cg.context.Int32Type(),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0),
				llvm.PointerType(cg.context.Int8Type(), 0),
			},
			false,
		)
		strcmp = llvm.AddFunction(cg.module, "strcmp", strcmpType)
		strcmp.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strcmp"] = strcmp
	}

	result := cg.builder.CreateCall(strcmp.GlobalValueType(), strcmp, []llvm.Value{strEnd, suffix}, "strcmptmp")
	zero := llvm.ConstInt(cg.context.Int32Type(), 0, false)
	isEqual := cg.builder.CreateICmp(llvm.IntEQ, result, zero, "endswith_cmp")
	trueVal := cg.builder.CreateZExt(isEqual, cg.context.Int64Type(), "endswith_true")
	cg.builder.CreateBr(mergeBB)

	cg.builder.SetInsertPointAtEnd(mergeBB)
	phi := cg.builder.CreatePHI(cg.context.Int64Type(), "endswith_result")
	phi.AddIncoming(
		[]llvm.Value{llvm.ConstInt(cg.context.Int64Type(), 0, false), trueVal},
		[]llvm.BasicBlock{falseBB, checkBB},
	)

	return phi, nil
}

// generateSplit splits a string by a delimiter (returns first token for simplicity)
func (cg *LLVMCodeGenerator) generateSplit(call *FunctionCall) (llvm.Value, error) {
	// For a full split implementation we'd need to return an array
	// For now, return the first token using strtok
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("split requires 2 arguments")
	}

	delim, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Copy the string first (strtok modifies in place)
	copyCall := &FunctionCall{Name: "copy", Args: []ASTNode{call.Args[0]}}
	strCopy, err := cg.generateStrCopy(copyCall)
	if err != nil {
		return llvm.Value{}, err
	}

	// Declare strtok
	strtok := cg.functions["strtok"]
	if strtok.IsNil() {
		strtokType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0),
				llvm.PointerType(cg.context.Int8Type(), 0),
			},
			false,
		)
		strtok = llvm.AddFunction(cg.module, "strtok", strtokType)
		strtok.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strtok"] = strtok
	}

	// Get first token
	token := cg.builder.CreateCall(strtok.GlobalValueType(), strtok, []llvm.Value{strCopy, delim}, "token")

	return token, nil
}

// generateReplace replaces occurrences of old with new in a string (first occurrence only)
func (cg *LLVMCodeGenerator) generateReplace(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, fmt.Errorf("replace requires 3 arguments (str, old, new)")
	}

	str, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	oldStr, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	newStr, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}

	// Get lengths
	strlen := cg.functions["strlen"]
	if strlen.IsNil() {
		strlenType := llvm.FunctionType(
			cg.context.Int64Type(),
			[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0)},
			false,
		)
		strlen = llvm.AddFunction(cg.module, "strlen", strlenType)
		strlen.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strlen"] = strlen
	}

	strLen := cg.builder.CreateCall(strlen.GlobalValueType(), strlen, []llvm.Value{str}, "strlen")
	oldLen := cg.builder.CreateCall(strlen.GlobalValueType(), strlen, []llvm.Value{oldStr}, "oldlen")
	newLen := cg.builder.CreateCall(strlen.GlobalValueType(), strlen, []llvm.Value{newStr}, "newlen")

	// Find old in str using strstr
	strstr := cg.functions["strstr"]
	if strstr.IsNil() {
		strstrType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0),
				llvm.PointerType(cg.context.Int8Type(), 0),
			},
			false,
		)
		strstr = llvm.AddFunction(cg.module, "strstr", strstrType)
		strstr.SetLinkage(llvm.ExternalLinkage)
		cg.functions["strstr"] = strstr
	}

	found := cg.builder.CreateCall(strstr.GlobalValueType(), strstr, []llvm.Value{str, oldStr}, "found")

	currentFn := cg.builder.GetInsertBlock().Parent()
	replaceBB := cg.context.AddBasicBlock(currentFn, "replace_do")
	copyBB := cg.context.AddBasicBlock(currentFn, "replace_copy")
	mergeBB := cg.context.AddBasicBlock(currentFn, "replace_merge")

	nullPtr := llvm.ConstNull(llvm.PointerType(cg.context.Int8Type(), 0))
	isNull := cg.builder.CreateICmp(llvm.IntEQ, found, nullPtr, "isfoundnull")
	cg.builder.CreateCondBr(isNull, copyBB, replaceBB)

	// Not found: just copy the string
	cg.builder.SetInsertPointAtEnd(copyBB)
	copiedStr, err := cg.generateStrCopy(&FunctionCall{Name: "copy", Args: []ASTNode{call.Args[0]}})
	if err != nil {
		return llvm.Value{}, err
	}
	cg.builder.CreateBr(mergeBB)

	// Found: do the replacement
	cg.builder.SetInsertPointAtEnd(replaceBB)
	// Calculate new size: strLen - oldLen + newLen + 1
	diff := cg.builder.CreateSub(newLen, oldLen, "diff")
	resultLen := cg.builder.CreateAdd(strLen, diff, "resultlen")
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	allocSize := cg.builder.CreateAdd(resultLen, one, "replacealloc")

	malloc := cg.functions["malloc"]
	if malloc.IsNil() {
		mallocType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{cg.context.Int64Type()},
			false,
		)
		malloc = llvm.AddFunction(cg.module, "malloc", mallocType)
		malloc.SetLinkage(llvm.ExternalLinkage)
		cg.functions["malloc"] = malloc
	}

	result := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{allocSize}, "replacebuf")

	memcpy := cg.functions["memcpy"]
	if memcpy.IsNil() {
		memcpyType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0),
				llvm.PointerType(cg.context.Int8Type(), 0),
				cg.context.Int64Type(),
			},
			false,
		)
		memcpy = llvm.AddFunction(cg.module, "memcpy", memcpyType)
		memcpy.SetLinkage(llvm.ExternalLinkage)
		cg.functions["memcpy"] = memcpy
	}

	// Copy prefix (str to found)
	strInt := cg.builder.CreatePtrToInt(str, cg.context.Int64Type(), "strint")
	foundInt := cg.builder.CreatePtrToInt(found, cg.context.Int64Type(), "foundint")
	prefixLen := cg.builder.CreateSub(foundInt, strInt, "prefixlen")
	cg.builder.CreateCall(memcpy.GlobalValueType(), memcpy, []llvm.Value{result, str, prefixLen}, "cpyprefix")

	// Copy new string
	dest1 := cg.builder.CreateGEP(cg.context.Int8Type(), result, []llvm.Value{prefixLen}, "dest1")
	cg.builder.CreateCall(memcpy.GlobalValueType(), memcpy, []llvm.Value{dest1, newStr, newLen}, "cpynew")

	// Copy suffix
	suffixStart := cg.builder.CreateGEP(cg.context.Int8Type(), found, []llvm.Value{oldLen}, "suffixstart")
	suffixLen := cg.builder.CreateCall(strlen.GlobalValueType(), strlen, []llvm.Value{suffixStart}, "suffixlen")
	dest2Off := cg.builder.CreateAdd(prefixLen, newLen, "dest2off")
	dest2 := cg.builder.CreateGEP(cg.context.Int8Type(), result, []llvm.Value{dest2Off}, "dest2")
	cg.builder.CreateCall(memcpy.GlobalValueType(), memcpy, []llvm.Value{dest2, suffixStart, suffixLen}, "cpysuffix")

	// Add null terminator
	nullPos := cg.builder.CreateGEP(cg.context.Int8Type(), result, []llvm.Value{resultLen}, "replacenullpos")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int8Type(), 0, false), nullPos)

	cg.builder.CreateBr(mergeBB)

	cg.builder.SetInsertPointAtEnd(mergeBB)
	phi := cg.builder.CreatePHI(llvm.PointerType(cg.context.Int8Type(), 0), "replace_result")
	phi.AddIncoming([]llvm.Value{copiedStr, result}, []llvm.BasicBlock{copyBB, replaceBB})

	return phi, nil
}

// ============================================================================
// REGEX MODULE - LLVM Implementation
// Uses POSIX regex functions: regcomp, regexec, regfree
// ============================================================================

// declareRegexFunctions declares the POSIX regex functions from libc
func (cg *LLVMCodeGenerator) declareRegexFunctions() {
	// regex_t is a structure, we use an opaque pointer type for simplicity
	// In practice, regex_t on Linux is ~64 bytes, we allocate 128 to be safe

	if _, ok := cg.functions["regcomp"]; !ok {
		// int regcomp(regex_t *preg, const char *pattern, int cflags)
		regcompType := llvm.FunctionType(
			cg.context.Int32Type(),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0), // regex_t*
				llvm.PointerType(cg.context.Int8Type(), 0), // pattern
				cg.context.Int32Type(),                     // cflags
			},
			false,
		)
		regcompFn := llvm.AddFunction(cg.module, "regcomp", regcompType)
		regcompFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["regcomp"] = regcompFn
	}

	if _, ok := cg.functions["regexec"]; !ok {
		// int regexec(const regex_t *preg, const char *string, size_t nmatch, regmatch_t *pmatch, int eflags)
		regexecType := llvm.FunctionType(
			cg.context.Int32Type(),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0), // regex_t*
				llvm.PointerType(cg.context.Int8Type(), 0), // string
				cg.context.Int64Type(),                     // nmatch
				llvm.PointerType(cg.context.Int8Type(), 0), // regmatch_t*
				cg.context.Int32Type(),                     // eflags
			},
			false,
		)
		regexecFn := llvm.AddFunction(cg.module, "regexec", regexecType)
		regexecFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["regexec"] = regexecFn
	}

	if _, ok := cg.functions["regfree"]; !ok {
		// void regfree(regex_t *preg)
		regfreeType := llvm.FunctionType(
			cg.context.VoidType(),
			[]llvm.Type{
				llvm.PointerType(cg.context.Int8Type(), 0), // regex_t*
			},
			false,
		)
		regfreeFn := llvm.AddFunction(cg.module, "regfree", regfreeType)
		regfreeFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["regfree"] = regfreeFn
	}
}

// generateRegexMatch generates code for regex::match(pattern, string) -> bool
// Returns 1 if pattern matches string, 0 otherwise
func (cg *LLVMCodeGenerator) generateRegexMatch(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) < 2 {
		return llvm.Value{}, fmt.Errorf("regex::match requires 2 arguments (pattern, string)")
	}

	cg.declareRegexFunctions()

	pattern, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	str, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Allocate space for regex_t (128 bytes to be safe)
	size := llvm.ConstInt(cg.context.Int64Type(), 128, false)
	malloc := cg.functions["malloc"]
	regex := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{size}, "regex")

	// Compile the regex with REG_EXTENDED (1)
	regcomp := cg.functions["regcomp"]
	cflags := llvm.ConstInt(cg.context.Int32Type(), 1, false) // REG_EXTENDED
	compResult := cg.builder.CreateCall(regcomp.GlobalValueType(), regcomp, []llvm.Value{regex, pattern, cflags}, "compresult")

	// Check if compilation succeeded
	currentFn := cg.builder.GetInsertBlock().Parent()
	okBB := cg.context.AddBasicBlock(currentFn, "regex_comp_ok")
	failBB := cg.context.AddBasicBlock(currentFn, "regex_comp_fail")
	mergeBB := cg.context.AddBasicBlock(currentFn, "regex_match_merge")

	zero := llvm.ConstInt(cg.context.Int32Type(), 0, false)
	isOk := cg.builder.CreateICmp(llvm.IntEQ, compResult, zero, "compok")
	cg.builder.CreateCondBr(isOk, okBB, failBB)

	// Compilation failed - return false
	cg.builder.SetInsertPointAtEnd(failBB)
	free := cg.functions["free"]
	cg.builder.CreateCall(free.GlobalValueType(), free, []llvm.Value{regex}, "")
	falseval := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	cg.builder.CreateBr(mergeBB)

	// Compilation succeeded - execute regex
	cg.builder.SetInsertPointAtEnd(okBB)
	regexec := cg.functions["regexec"]
	nmatch := llvm.ConstInt(cg.context.Int64Type(), 0, false)
	nullPtr := llvm.ConstNull(llvm.PointerType(cg.context.Int8Type(), 0))
	eflags := llvm.ConstInt(cg.context.Int32Type(), 0, false)
	execResult := cg.builder.CreateCall(regexec.GlobalValueType(), regexec,
		[]llvm.Value{regex, str, nmatch, nullPtr, eflags}, "execresult")

	// Free the regex
	regfree := cg.functions["regfree"]
	cg.builder.CreateCall(regfree.GlobalValueType(), regfree, []llvm.Value{regex}, "")
	cg.builder.CreateCall(free.GlobalValueType(), free, []llvm.Value{regex}, "")

	// execResult == 0 means match found
	matched := cg.builder.CreateICmp(llvm.IntEQ, execResult, zero, "matched")
	trueval := cg.builder.CreateZExt(matched, cg.context.Int64Type(), "matchresult")
	cg.builder.CreateBr(mergeBB)

	// Merge block
	cg.builder.SetInsertPointAtEnd(mergeBB)
	phi := cg.builder.CreatePHI(cg.context.Int64Type(), "regex_match_result")
	phi.AddIncoming([]llvm.Value{falseval, trueval}, []llvm.BasicBlock{failBB, okBB})

	return phi, nil
}

// generateRegexFind generates code for regex::find(pattern, string) -> int
// Returns the position of the first match, or -1 if no match
func (cg *LLVMCodeGenerator) generateRegexFind(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) < 2 {
		return llvm.Value{}, fmt.Errorf("regex::find requires 2 arguments (pattern, string)")
	}

	cg.declareRegexFunctions()

	pattern, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	str, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Allocate space for regex_t (128 bytes)
	malloc := cg.functions["malloc"]
	regexSize := llvm.ConstInt(cg.context.Int64Type(), 128, false)
	regex := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{regexSize}, "regex")

	// Allocate space for regmatch_t[1] (each regmatch_t is 16 bytes: rm_so and rm_eo as int64)
	matchSize := llvm.ConstInt(cg.context.Int64Type(), 16, false)
	pmatch := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{matchSize}, "pmatch")

	// Compile the regex
	regcomp := cg.functions["regcomp"]
	cflags := llvm.ConstInt(cg.context.Int32Type(), 1, false) // REG_EXTENDED
	compResult := cg.builder.CreateCall(regcomp.GlobalValueType(), regcomp, []llvm.Value{regex, pattern, cflags}, "compresult")

	currentFn := cg.builder.GetInsertBlock().Parent()
	okBB := cg.context.AddBasicBlock(currentFn, "regex_find_ok")
	failBB := cg.context.AddBasicBlock(currentFn, "regex_find_fail")
	foundBB := cg.context.AddBasicBlock(currentFn, "regex_found")
	notFoundBB := cg.context.AddBasicBlock(currentFn, "regex_not_found")
	mergeBB := cg.context.AddBasicBlock(currentFn, "regex_find_merge")

	zero := llvm.ConstInt(cg.context.Int32Type(), 0, false)
	isOk := cg.builder.CreateICmp(llvm.IntEQ, compResult, zero, "compok")
	cg.builder.CreateCondBr(isOk, okBB, failBB)

	// Compilation failed
	cg.builder.SetInsertPointAtEnd(failBB)
	free := cg.functions["free"]
	cg.builder.CreateCall(free.GlobalValueType(), free, []llvm.Value{regex}, "")
	cg.builder.CreateCall(free.GlobalValueType(), free, []llvm.Value{pmatch}, "")
	notFound := llvm.ConstInt(cg.context.Int64Type(), 0xFFFFFFFFFFFFFFFF, true) // -1
	cg.builder.CreateBr(mergeBB)

	// Execute regex
	cg.builder.SetInsertPointAtEnd(okBB)
	regexec := cg.functions["regexec"]
	nmatch := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	eflags := llvm.ConstInt(cg.context.Int32Type(), 0, false)
	execResult := cg.builder.CreateCall(regexec.GlobalValueType(), regexec,
		[]llvm.Value{regex, str, nmatch, pmatch, eflags}, "execresult")

	regfree := cg.functions["regfree"]
	cg.builder.CreateCall(regfree.GlobalValueType(), regfree, []llvm.Value{regex}, "")
	cg.builder.CreateCall(free.GlobalValueType(), free, []llvm.Value{regex}, "")

	matched := cg.builder.CreateICmp(llvm.IntEQ, execResult, zero, "matched")
	cg.builder.CreateCondBr(matched, foundBB, notFoundBB)

	// Match found - get rm_so (start offset)
	cg.builder.SetInsertPointAtEnd(foundBB)
	// Cast pmatch to int64* to read rm_so
	pmatchI64 := cg.builder.CreateBitCast(pmatch, llvm.PointerType(cg.context.Int64Type(), 0), "pmatchi64")
	rmSo := cg.builder.CreateLoad(cg.context.Int64Type(), pmatchI64, "rm_so")
	cg.builder.CreateCall(free.GlobalValueType(), free, []llvm.Value{pmatch}, "")
	cg.builder.CreateBr(mergeBB)

	// Not found
	cg.builder.SetInsertPointAtEnd(notFoundBB)
	cg.builder.CreateCall(free.GlobalValueType(), free, []llvm.Value{pmatch}, "")
	notFound2 := llvm.ConstInt(cg.context.Int64Type(), 0xFFFFFFFFFFFFFFFF, true) // -1
	cg.builder.CreateBr(mergeBB)

	// Merge
	cg.builder.SetInsertPointAtEnd(mergeBB)
	phi := cg.builder.CreatePHI(cg.context.Int64Type(), "regex_find_result")
	phi.AddIncoming([]llvm.Value{notFound, rmSo, notFound2}, []llvm.BasicBlock{failBB, foundBB, notFoundBB})

	return phi, nil
}

// generateRegexReplace generates code for regex::replace(pattern, replacement, string) -> string
// Replaces the first match with the replacement string
func (cg *LLVMCodeGenerator) generateRegexReplace(call *FunctionCall) (llvm.Value, error) {
	// For now, return the original string (stub)
	// Full implementation would be complex due to memory management
	if len(call.Args) < 3 {
		return llvm.Value{}, fmt.Errorf("regex::replace requires 3 arguments (pattern, replacement, string)")
	}

	// Return the original string for now
	return cg.generateExpression(call.Args[2])
}

// generateRegexReplaceAll generates code for regex::replace_all(pattern, replacement, string) -> string
func (cg *LLVMCodeGenerator) generateRegexReplaceAll(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) < 3 {
		return llvm.Value{}, fmt.Errorf("regex::replace_all requires 3 arguments (pattern, replacement, string)")
	}

	// Return the original string for now (stub)
	return cg.generateExpression(call.Args[2])
}

// generateRegexSplit generates code for regex::split(pattern, string) -> returns first token
func (cg *LLVMCodeGenerator) generateRegexSplit(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) < 2 {
		return llvm.Value{}, fmt.Errorf("regex::split requires 2 arguments (pattern, string)")
	}

	// Return the original string for now (stub)
	return cg.generateExpression(call.Args[1])
}

// generateRegexFindAll generates code for regex::find_all(pattern, string) -> returns first match position
func (cg *LLVMCodeGenerator) generateRegexFindAll(call *FunctionCall) (llvm.Value, error) {
	// For now, just call find to get the first match
	return cg.generateRegexFind(call)
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// declareSyscall declares the syscall function for Linux system calls
func (cg *LLVMCodeGenerator) declareSyscall() {
	if _, ok := cg.functions["syscall"]; !ok {
		// syscall(number, arg1, arg2, arg3, ...) -> long
		syscallType := llvm.FunctionType(
			cg.context.Int64Type(),
			[]llvm.Type{
				cg.context.Int64Type(), // syscall number
				cg.context.Int64Type(), // arg1
				cg.context.Int64Type(), // arg2
				cg.context.Int64Type(), // arg3
			},
			true, // variadic for more args
		)
		syscallFn := llvm.AddFunction(cg.module, "syscall", syscallType)
		syscallFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["syscall"] = syscallFn
	}
}

// ============================================================================
// JSON MODULE
// ============================================================================
//
// This module implements JSON parsing and serialization for the LLVM backend.
// The JSON module provides functions to:
// - Parse JSON strings into structured data
// - Serialize data structures to JSON strings
// - Access JSON values by path
// - Create and manipulate JSON objects and arrays
//
// JSON Types:
// - JSON_NULL   = 0
// - JSON_BOOL   = 1
// - JSON_NUMBER = 2
// - JSON_STRING = 3
// - JSON_ARRAY  = 4
// - JSON_OBJECT = 5

const (
	JSONNull   = 0
	JSONBool   = 1
	JSONNumber = 2
	JSONString = 3
	JSONArray  = 4
	JSONObject = 5
)

// JSONValue structure in memory:
// Offset 0:  type (i64)      - JSON type enum
// Offset 8:  value (i64)     - value or pointer to data
// Offset 16: size (i64)      - size for arrays/objects
// Offset 24: capacity (i64)  - allocated capacity
// Total: 32 bytes

const JSONValueSize = 32

// JSONValue is the Go representation for testing
type JSONValue struct {
	Type      int
	IntValue  int64
	StrValue  string
	BoolValue bool
	Elements  []*JSONValue
	Fields    map[string]*JSONValue
}

// ============================================================================
// JSON FUNCTION GENERATORS
// ============================================================================

// generateJSONParse parses a JSON string and returns a JSONValue pointer
// json_parse(str) -> json_value_ptr
func (cg *LLVMCodeGenerator) generateJSONParse(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	str, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Allocate a JSONValue structure
	valueSize := llvm.ConstInt(cg.context.Int64Type(), JSONValueSize, false)
	malloc := cg.functions["malloc"]
	jsonPtr := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{valueSize}, "jsonalloc")
	jsonPtrInt := cg.builder.CreatePtrToInt(jsonPtr, cg.context.Int64Type(), "jsonptr")

	// Call our JSON parser helper (simplified inline parsing)
	cg.generateJSONParseInline(str, jsonPtrInt)

	return jsonPtrInt, nil
}

// generateJSONParseInline generates inline code to parse JSON
func (cg *LLVMCodeGenerator) generateJSONParseInline(str llvm.Value, jsonPtr llvm.Value) {
	// Get first character
	firstChar := cg.builder.CreateLoad(cg.context.Int8Type(), str, "firstchar")

	// Create blocks for type detection
	fn := cg.currentFn
	nullBlock := llvm.AddBasicBlock(fn, "json_null")
	boolBlock := llvm.AddBasicBlock(fn, "json_bool")
	numBlock := llvm.AddBasicBlock(fn, "json_num")
	strBlock := llvm.AddBasicBlock(fn, "json_str")
	arrBlock := llvm.AddBasicBlock(fn, "json_arr")
	objBlock := llvm.AddBasicBlock(fn, "json_obj")
	mergeBlock := llvm.AddBasicBlock(fn, "json_merge")

	// Check for null
	isN := cg.builder.CreateICmp(llvm.IntEQ, firstChar, llvm.ConstInt(cg.context.Int8Type(), 'n', false), "isn")
	cg.builder.CreateCondBr(isN, nullBlock, boolBlock)

	// Null block
	cg.builder.SetInsertPointAtEnd(nullBlock)
	typePtr := cg.builder.CreateIntToPtr(jsonPtr, llvm.PointerType(cg.context.Int64Type(), 0), "typeptr")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int64Type(), JSONNull, false), typePtr)
	cg.builder.CreateBr(mergeBlock)

	// Bool detection (t for true, f for false)
	cg.builder.SetInsertPointAtEnd(boolBlock)
	isT := cg.builder.CreateICmp(llvm.IntEQ, firstChar, llvm.ConstInt(cg.context.Int8Type(), 't', false), "ist")
	isF := cg.builder.CreateICmp(llvm.IntEQ, firstChar, llvm.ConstInt(cg.context.Int8Type(), 'f', false), "isf")
	isBool := cg.builder.CreateOr(isT, isF, "isbool")
	cg.builder.CreateCondBr(isBool, objBlock, numBlock)

	// Number detection (digit or minus)
	cg.builder.SetInsertPointAtEnd(numBlock)
	isMinus := cg.builder.CreateICmp(llvm.IntEQ, firstChar, llvm.ConstInt(cg.context.Int8Type(), '-', false), "isminus")
	isDigit := cg.builder.CreateAnd(
		cg.builder.CreateICmp(llvm.IntSGE, firstChar, llvm.ConstInt(cg.context.Int8Type(), '0', false), "ge0"),
		cg.builder.CreateICmp(llvm.IntSLE, firstChar, llvm.ConstInt(cg.context.Int8Type(), '9', false), "le9"),
		"isdigit")
	isNum := cg.builder.CreateOr(isMinus, isDigit, "isnum")
	cg.builder.CreateCondBr(isNum, strBlock, strBlock)

	// String detection (quote)
	cg.builder.SetInsertPointAtEnd(strBlock)
	isQuote := cg.builder.CreateICmp(llvm.IntEQ, firstChar, llvm.ConstInt(cg.context.Int8Type(), '"', false), "isquote")
	cg.builder.CreateCondBr(isQuote, arrBlock, arrBlock)

	// Array detection (bracket)
	cg.builder.SetInsertPointAtEnd(arrBlock)
	isBracket := cg.builder.CreateICmp(llvm.IntEQ, firstChar, llvm.ConstInt(cg.context.Int8Type(), '[', false), "isbracket")
	cg.builder.CreateCondBr(isBracket, objBlock, objBlock)

	// Object detection (brace) - default case
	cg.builder.SetInsertPointAtEnd(objBlock)
	objTypePtr := cg.builder.CreateIntToPtr(jsonPtr, llvm.PointerType(cg.context.Int64Type(), 0), "objtypeptr")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int64Type(), JSONObject, false), objTypePtr)
	cg.builder.CreateBr(mergeBlock)

	// Merge block
	cg.builder.SetInsertPointAtEnd(mergeBlock)
}

// generateJSONStringify converts a JSONValue to a string
// json_stringify(json_value_ptr) -> string
func (cg *LLVMCodeGenerator) generateJSONStringify(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	jsonPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Allocate output buffer (1024 bytes)
	bufSize := llvm.ConstInt(cg.context.Int64Type(), 1024, false)
	malloc := cg.functions["malloc"]
	buf := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{bufSize}, "strbuf")

	// Get the JSON type
	typePtr := cg.builder.CreateIntToPtr(jsonPtr, llvm.PointerType(cg.context.Int64Type(), 0), "typeptr")
	jsonType := cg.builder.CreateLoad(cg.context.Int64Type(), typePtr, "jsontype")

	// Generate string based on type
	cg.generateJSONStringifyByType(buf, jsonPtr, jsonType)

	return buf, nil
}

// generateJSONStringifyByType generates the string representation
func (cg *LLVMCodeGenerator) generateJSONStringifyByType(buf, jsonPtr, jsonType llvm.Value) {
	fn := cg.currentFn
	nullBlock := llvm.AddBasicBlock(fn, "stringify_null")
	otherBlock := llvm.AddBasicBlock(fn, "stringify_other")
	mergeBlock := llvm.AddBasicBlock(fn, "stringify_merge")

	isNull := cg.builder.CreateICmp(llvm.IntEQ, jsonType, llvm.ConstInt(cg.context.Int64Type(), JSONNull, false), "isnull")
	cg.builder.CreateCondBr(isNull, nullBlock, otherBlock)

	// Null block - write "null"
	cg.builder.SetInsertPointAtEnd(nullBlock)
	nullStr := cg.createGlobalString("null")
	cg.declareStringHelpers()
	strcpy := cg.functions["strcpy"]
	cg.builder.CreateCall(strcpy.GlobalValueType(), strcpy, []llvm.Value{buf, nullStr}, "")
	cg.builder.CreateBr(mergeBlock)

	// Other types - write "{}"
	cg.builder.SetInsertPointAtEnd(otherBlock)
	objStr := cg.createGlobalString("{}")
	cg.builder.CreateCall(strcpy.GlobalValueType(), strcpy, []llvm.Value{buf, objStr}, "")
	cg.builder.CreateBr(mergeBlock)

	cg.builder.SetInsertPointAtEnd(mergeBlock)
}

// generateJSONGet gets a value from a JSON object by key
// json_get(json_obj, key) -> json_value_ptr
func (cg *LLVMCodeGenerator) generateJSONGet(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, nil
	}

	jsonPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	_, err = cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Simplified: return the same pointer (would need proper lookup)
	return jsonPtr, nil
}

// generateJSONGetType returns the type of a JSON value
// json_type(json_value_ptr) -> int
func (cg *LLVMCodeGenerator) generateJSONGetType(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	jsonPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Type is at offset 0
	typePtr := cg.builder.CreateIntToPtr(jsonPtr, llvm.PointerType(cg.context.Int64Type(), 0), "typeptr")
	return cg.builder.CreateLoad(cg.context.Int64Type(), typePtr, "jsontype"), nil
}

// generateJSONGetInt gets an integer value from a JSON number
// json_int(json_value_ptr) -> int
func (cg *LLVMCodeGenerator) generateJSONGetInt(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	jsonPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Value is at offset 8
	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	valueOffset := cg.builder.CreateAdd(jsonPtr, eight, "valoffset")
	valuePtr := cg.builder.CreateIntToPtr(valueOffset, llvm.PointerType(cg.context.Int64Type(), 0), "valptr")
	return cg.builder.CreateLoad(cg.context.Int64Type(), valuePtr, "jsonint"), nil
}

// generateJSONGetString gets a string value from a JSON string
// json_str(json_value_ptr) -> string
func (cg *LLVMCodeGenerator) generateJSONGetString(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	jsonPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// String pointer is at offset 8
	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	valueOffset := cg.builder.CreateAdd(jsonPtr, eight, "valoffset")
	valuePtr := cg.builder.CreateIntToPtr(valueOffset, llvm.PointerType(cg.context.Int64Type(), 0), "valptr")
	strPtr := cg.builder.CreateLoad(cg.context.Int64Type(), valuePtr, "strptr")
	return cg.builder.CreateIntToPtr(strPtr, llvm.PointerType(cg.context.Int8Type(), 0), "jsonstr"), nil
}

// generateJSONGetBool gets a boolean value from a JSON boolean
// json_bool(json_value_ptr) -> bool
func (cg *LLVMCodeGenerator) generateJSONGetBool(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	jsonPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Boolean value is at offset 8
	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	valueOffset := cg.builder.CreateAdd(jsonPtr, eight, "valoffset")
	valuePtr := cg.builder.CreateIntToPtr(valueOffset, llvm.PointerType(cg.context.Int64Type(), 0), "valptr")
	return cg.builder.CreateLoad(cg.context.Int64Type(), valuePtr, "jsonbool"), nil
}

// generateJSONArrayLen gets the length of a JSON array
// json_array_len(json_value_ptr) -> int
func (cg *LLVMCodeGenerator) generateJSONArrayLen(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	jsonPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Size is at offset 16
	sixteen := llvm.ConstInt(cg.context.Int64Type(), 16, false)
	sizeOffset := cg.builder.CreateAdd(jsonPtr, sixteen, "sizeoffset")
	sizePtr := cg.builder.CreateIntToPtr(sizeOffset, llvm.PointerType(cg.context.Int64Type(), 0), "sizeptr")
	return cg.builder.CreateLoad(cg.context.Int64Type(), sizePtr, "arrlen"), nil
}

// generateJSONArrayGet gets an element from a JSON array by index
// json_array_get(json_value_ptr, index) -> json_value_ptr
func (cg *LLVMCodeGenerator) generateJSONArrayGet(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, nil
	}

	jsonPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	index, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Array data starts at offset 8 (pointer to array of JSONValue)
	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	arrOffset := cg.builder.CreateAdd(jsonPtr, eight, "arroffset")
	arrPtr := cg.builder.CreateIntToPtr(arrOffset, llvm.PointerType(cg.context.Int64Type(), 0), "arrptr")
	arrBase := cg.builder.CreateLoad(cg.context.Int64Type(), arrPtr, "arrbase")

	// Calculate element offset: base + index * JSONValueSize
	elemSize := llvm.ConstInt(cg.context.Int64Type(), JSONValueSize, false)
	elemOffset := cg.builder.CreateMul(index, elemSize, "elemoffset")
	elemAddr := cg.builder.CreateAdd(arrBase, elemOffset, "elemaddr")

	return elemAddr, nil
}

// generateJSONFree frees a JSON value and its contents
// json_free(json_value_ptr)
func (cg *LLVMCodeGenerator) generateJSONFree(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	jsonPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	ptr := cg.builder.CreateIntToPtr(jsonPtr, llvm.PointerType(cg.context.Int8Type(), 0), "freeptr")
	free := cg.functions["free"]
	cg.builder.CreateCall(free.GlobalValueType(), free, []llvm.Value{ptr}, "")

	return llvm.ConstInt(cg.context.Int64Type(), 0, false), nil
}

// generateJSONNew creates a new JSON value of the given type
// json_new(type) -> json_value_ptr
func (cg *LLVMCodeGenerator) generateJSONNew(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}

	jsonType, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Allocate JSONValue structure
	valueSize := llvm.ConstInt(cg.context.Int64Type(), JSONValueSize, false)
	malloc := cg.functions["malloc"]
	ptr := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{valueSize}, "jsonalloc")
	ptrInt := cg.builder.CreatePtrToInt(ptr, cg.context.Int64Type(), "jsonptr")

	// Store type
	typePtr := cg.builder.CreateIntToPtr(ptrInt, llvm.PointerType(cg.context.Int64Type(), 0), "typeptr")
	cg.builder.CreateStore(jsonType, typePtr)

	// Initialize value to 0
	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	valueOffset := cg.builder.CreateAdd(ptrInt, eight, "valoffset")
	valuePtr := cg.builder.CreateIntToPtr(valueOffset, llvm.PointerType(cg.context.Int64Type(), 0), "valptr")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int64Type(), 0, false), valuePtr)

	return ptrInt, nil
}

// ============================================================================
// HASH FUNCTIONS
// ============================================================================

// generateDjb2 generates DJB2 hash for a string
func (cg *LLVMCodeGenerator) generateDjb2(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("djb2 requires 1 argument")
	}

	str, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// DJB2: hash = 5381; for each c: hash = hash * 33 + c
	currentFn := cg.builder.GetInsertBlock().Parent()
	loopBB := llvm.AddBasicBlock(currentFn, "djb2loop")
	bodyBB := llvm.AddBasicBlock(currentFn, "djb2body")
	doneBB := llvm.AddBasicBlock(currentFn, "djb2done")

	// Initialize hash = 5381
	hashPtr := cg.builder.CreateAlloca(cg.context.Int64Type(), "djb2hash")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int64Type(), 5381, false), hashPtr)

	// Initialize pointer
	ptrPtr := cg.builder.CreateAlloca(llvm.PointerType(cg.context.Int8Type(), 0), "djb2ptr")
	strPtr := cg.builder.CreateIntToPtr(str, llvm.PointerType(cg.context.Int8Type(), 0), "strptr")
	cg.builder.CreateStore(strPtr, ptrPtr)

	cg.builder.CreateBr(loopBB)

	// Loop: check if *ptr != 0
	cg.builder.SetInsertPointAtEnd(loopBB)
	ptr := cg.builder.CreateLoad(llvm.PointerType(cg.context.Int8Type(), 0), ptrPtr, "ptr")
	c := cg.builder.CreateLoad(cg.context.Int8Type(), ptr, "char")
	zero := llvm.ConstInt(cg.context.Int8Type(), 0, false)
	cond := cg.builder.CreateICmp(llvm.IntNE, c, zero, "djb2cond")
	cg.builder.CreateCondBr(cond, bodyBB, doneBB)

	// Body: hash = hash * 33 + c
	cg.builder.SetInsertPointAtEnd(bodyBB)
	hash := cg.builder.CreateLoad(cg.context.Int64Type(), hashPtr, "hash")
	thirtythree := llvm.ConstInt(cg.context.Int64Type(), 33, false)
	hash33 := cg.builder.CreateMul(hash, thirtythree, "hash33")
	cExt := cg.builder.CreateZExt(c, cg.context.Int64Type(), "cext")
	newHash := cg.builder.CreateAdd(hash33, cExt, "newhash")
	cg.builder.CreateStore(newHash, hashPtr)

	// Increment pointer
	one := llvm.ConstInt(cg.context.Int32Type(), 1, false)
	nextPtr := cg.builder.CreateGEP(cg.context.Int8Type(), ptr, []llvm.Value{one}, "nextptr")
	cg.builder.CreateStore(nextPtr, ptrPtr)
	cg.builder.CreateBr(loopBB)

	// Done
	cg.builder.SetInsertPointAtEnd(doneBB)
	result := cg.builder.CreateLoad(cg.context.Int64Type(), hashPtr, "djb2result")
	return result, nil
}

// generateFnv1a generates FNV-1a hash for data
func (cg *LLVMCodeGenerator) generateFnv1a(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("fnv1a requires 2 arguments (data, len)")
	}

	data, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	length, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// FNV-1a 64-bit: hash = 14695981039346656037; for each byte: hash ^= byte; hash *= 1099511628211
	currentFn := cg.builder.GetInsertBlock().Parent()
	loopBB := llvm.AddBasicBlock(currentFn, "fnvloop")
	bodyBB := llvm.AddBasicBlock(currentFn, "fnvbody")
	doneBB := llvm.AddBasicBlock(currentFn, "fnvdone")

	// FNV offset basis
	fnvOffset := llvm.ConstInt(cg.context.Int64Type(), 14695981039346656037, false)
	fnvPrime := llvm.ConstInt(cg.context.Int64Type(), 1099511628211, false)

	hashPtr := cg.builder.CreateAlloca(cg.context.Int64Type(), "fnvhash")
	cg.builder.CreateStore(fnvOffset, hashPtr)

	idxPtr := cg.builder.CreateAlloca(cg.context.Int64Type(), "fnvidx")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int64Type(), 0, false), idxPtr)

	dataPtr := cg.builder.CreateIntToPtr(data, llvm.PointerType(cg.context.Int8Type(), 0), "dataptr")

	cg.builder.CreateBr(loopBB)

	// Loop: while idx < length
	cg.builder.SetInsertPointAtEnd(loopBB)
	idx := cg.builder.CreateLoad(cg.context.Int64Type(), idxPtr, "idx")
	cond := cg.builder.CreateICmp(llvm.IntSLT, idx, length, "fnvcond")
	cg.builder.CreateCondBr(cond, bodyBB, doneBB)

	// Body
	cg.builder.SetInsertPointAtEnd(bodyBB)
	bytePtr := cg.builder.CreateGEP(cg.context.Int8Type(), dataPtr, []llvm.Value{idx}, "byteptr")
	byteVal := cg.builder.CreateLoad(cg.context.Int8Type(), bytePtr, "byteval")
	byteExt := cg.builder.CreateZExt(byteVal, cg.context.Int64Type(), "byteext")

	hash := cg.builder.CreateLoad(cg.context.Int64Type(), hashPtr, "hash")
	hashXor := cg.builder.CreateXor(hash, byteExt, "hashxor")
	hashMul := cg.builder.CreateMul(hashXor, fnvPrime, "hashmul")
	cg.builder.CreateStore(hashMul, hashPtr)

	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	nextIdx := cg.builder.CreateAdd(idx, one, "nextidx")
	cg.builder.CreateStore(nextIdx, idxPtr)
	cg.builder.CreateBr(loopBB)

	// Done
	cg.builder.SetInsertPointAtEnd(doneBB)
	result := cg.builder.CreateLoad(cg.context.Int64Type(), hashPtr, "fnvresult")
	return result, nil
}

// generateCrc32 generates CRC32 hash (simple implementation)
func (cg *LLVMCodeGenerator) generateCrc32(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("crc32 requires 2 arguments (data, len)")
	}

	data, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	length, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Simple CRC32 using polynomial 0xEDB88320
	currentFn := cg.builder.GetInsertBlock().Parent()
	outerLoopBB := llvm.AddBasicBlock(currentFn, "crcouterloop")
	outerBodyBB := llvm.AddBasicBlock(currentFn, "crcouterbody")
	innerLoopBB := llvm.AddBasicBlock(currentFn, "crcinnerloop")
	innerBodyBB := llvm.AddBasicBlock(currentFn, "crcinnerbody")
	innerDoneBB := llvm.AddBasicBlock(currentFn, "crcinnerdone")
	doneBB := llvm.AddBasicBlock(currentFn, "crcdone")

	// Initialize CRC = 0xFFFFFFFF
	crcPtr := cg.builder.CreateAlloca(cg.context.Int64Type(), "crc")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int64Type(), 0xFFFFFFFF, false), crcPtr)

	idxPtr := cg.builder.CreateAlloca(cg.context.Int64Type(), "crcidx")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int64Type(), 0, false), idxPtr)

	bitPtr := cg.builder.CreateAlloca(cg.context.Int64Type(), "crcbit")

	dataPtr := cg.builder.CreateIntToPtr(data, llvm.PointerType(cg.context.Int8Type(), 0), "crcdataptr")

	cg.builder.CreateBr(outerLoopBB)

	// Outer loop: for each byte
	cg.builder.SetInsertPointAtEnd(outerLoopBB)
	idx := cg.builder.CreateLoad(cg.context.Int64Type(), idxPtr, "idx")
	outerCond := cg.builder.CreateICmp(llvm.IntSLT, idx, length, "outercond")
	cg.builder.CreateCondBr(outerCond, outerBodyBB, doneBB)

	cg.builder.SetInsertPointAtEnd(outerBodyBB)
	bytePtr := cg.builder.CreateGEP(cg.context.Int8Type(), dataPtr, []llvm.Value{idx}, "byteptr")
	byteVal := cg.builder.CreateLoad(cg.context.Int8Type(), bytePtr, "byteval")
	byteExt := cg.builder.CreateZExt(byteVal, cg.context.Int64Type(), "byteext")

	crc := cg.builder.CreateLoad(cg.context.Int64Type(), crcPtr, "crc")
	crcXor := cg.builder.CreateXor(crc, byteExt, "crcxor")
	cg.builder.CreateStore(crcXor, crcPtr)

	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int64Type(), 0, false), bitPtr)
	cg.builder.CreateBr(innerLoopBB)

	// Inner loop: for 8 bits
	cg.builder.SetInsertPointAtEnd(innerLoopBB)
	bit := cg.builder.CreateLoad(cg.context.Int64Type(), bitPtr, "bit")
	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	innerCond := cg.builder.CreateICmp(llvm.IntSLT, bit, eight, "innercond")
	cg.builder.CreateCondBr(innerCond, innerBodyBB, innerDoneBB)

	cg.builder.SetInsertPointAtEnd(innerBodyBB)
	crcVal := cg.builder.CreateLoad(cg.context.Int64Type(), crcPtr, "crcval")
	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	lsb := cg.builder.CreateAnd(crcVal, one, "lsb")
	isOne := cg.builder.CreateICmp(llvm.IntEQ, lsb, one, "isone")

	crcShift := cg.builder.CreateLShr(crcVal, one, "crcshift")
	poly := llvm.ConstInt(cg.context.Int64Type(), 0xEDB88320, false)
	crcXorPoly := cg.builder.CreateXor(crcShift, poly, "crcxorpoly")
	newCrc := cg.builder.CreateSelect(isOne, crcXorPoly, crcShift, "newcrc")
	cg.builder.CreateStore(newCrc, crcPtr)

	nextBit := cg.builder.CreateAdd(bit, one, "nextbit")
	cg.builder.CreateStore(nextBit, bitPtr)
	cg.builder.CreateBr(innerLoopBB)

	cg.builder.SetInsertPointAtEnd(innerDoneBB)
	nextIdx := cg.builder.CreateAdd(idx, one, "nextidx")
	cg.builder.CreateStore(nextIdx, idxPtr)
	cg.builder.CreateBr(outerLoopBB)

	// Done: return ~crc
	cg.builder.SetInsertPointAtEnd(doneBB)
	finalCrc := cg.builder.CreateLoad(cg.context.Int64Type(), crcPtr, "finalcrc")
	result := cg.builder.CreateNot(finalCrc, "crcresult")
	// Mask to 32 bits
	mask := llvm.ConstInt(cg.context.Int64Type(), 0xFFFFFFFF, false)
	return cg.builder.CreateAnd(result, mask, "crc32final"), nil
}

// generateMurmur generates MurmurHash3 (simplified 32-bit version)
func (cg *LLVMCodeGenerator) generateMurmur(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 3 {
		return llvm.Value{}, fmt.Errorf("murmur requires 3 arguments (data, len, seed)")
	}

	data, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	length, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	seed, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}

	// Simplified MurmurHash3 - process as bytes with mixing
	currentFn := cg.builder.GetInsertBlock().Parent()
	loopBB := llvm.AddBasicBlock(currentFn, "murmurloop")
	bodyBB := llvm.AddBasicBlock(currentFn, "murmurbody")
	doneBB := llvm.AddBasicBlock(currentFn, "murmurdone")

	hashPtr := cg.builder.CreateAlloca(cg.context.Int64Type(), "murmurhash")
	cg.builder.CreateStore(seed, hashPtr)

	idxPtr := cg.builder.CreateAlloca(cg.context.Int64Type(), "murmuridx")
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int64Type(), 0, false), idxPtr)

	dataPtr := cg.builder.CreateIntToPtr(data, llvm.PointerType(cg.context.Int8Type(), 0), "murmurdataptr")

	c1 := llvm.ConstInt(cg.context.Int64Type(), 0xcc9e2d51, false)
	c2 := llvm.ConstInt(cg.context.Int64Type(), 0x1b873593, false)

	cg.builder.CreateBr(loopBB)

	// Loop
	cg.builder.SetInsertPointAtEnd(loopBB)
	idx := cg.builder.CreateLoad(cg.context.Int64Type(), idxPtr, "idx")
	cond := cg.builder.CreateICmp(llvm.IntSLT, idx, length, "murmurcond")
	cg.builder.CreateCondBr(cond, bodyBB, doneBB)

	// Body: mix each byte
	cg.builder.SetInsertPointAtEnd(bodyBB)
	bytePtr := cg.builder.CreateGEP(cg.context.Int8Type(), dataPtr, []llvm.Value{idx}, "byteptr")
	byteVal := cg.builder.CreateLoad(cg.context.Int8Type(), bytePtr, "byteval")
	k := cg.builder.CreateZExt(byteVal, cg.context.Int64Type(), "k")

	// k *= c1; k = rotl(k, 15); k *= c2
	k1 := cg.builder.CreateMul(k, c1, "k1")
	fifteen := llvm.ConstInt(cg.context.Int64Type(), 15, false)
	fortyNine := llvm.ConstInt(cg.context.Int64Type(), 49, false)
	k2l := cg.builder.CreateShl(k1, fifteen, "k2l")
	k2r := cg.builder.CreateLShr(k1, fortyNine, "k2r")
	k2 := cg.builder.CreateOr(k2l, k2r, "k2")
	k3 := cg.builder.CreateMul(k2, c2, "k3")

	// hash ^= k; hash = rotl(hash, 13); hash = hash * 5 + 0xe6546b64
	hash := cg.builder.CreateLoad(cg.context.Int64Type(), hashPtr, "hash")
	hash1 := cg.builder.CreateXor(hash, k3, "hash1")
	thirteen := llvm.ConstInt(cg.context.Int64Type(), 13, false)
	fiftyOne := llvm.ConstInt(cg.context.Int64Type(), 51, false)
	hash2l := cg.builder.CreateShl(hash1, thirteen, "hash2l")
	hash2r := cg.builder.CreateLShr(hash1, fiftyOne, "hash2r")
	hash2 := cg.builder.CreateOr(hash2l, hash2r, "hash2")
	five := llvm.ConstInt(cg.context.Int64Type(), 5, false)
	hash3 := cg.builder.CreateMul(hash2, five, "hash3")
	magic := llvm.ConstInt(cg.context.Int64Type(), 0xe6546b64, false)
	hash4 := cg.builder.CreateAdd(hash3, magic, "hash4")
	cg.builder.CreateStore(hash4, hashPtr)

	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	nextIdx := cg.builder.CreateAdd(idx, one, "nextidx")
	cg.builder.CreateStore(nextIdx, idxPtr)
	cg.builder.CreateBr(loopBB)

	// Done: finalization
	cg.builder.SetInsertPointAtEnd(doneBB)
	finalHash := cg.builder.CreateLoad(cg.context.Int64Type(), hashPtr, "finalhash")
	h1 := cg.builder.CreateXor(finalHash, length, "h1")
	// fmix: h ^= h >> 16; h *= 0x85ebca6b; h ^= h >> 13; h *= 0xc2b2ae35; h ^= h >> 16
	sixteen := llvm.ConstInt(cg.context.Int64Type(), 16, false)
	h2 := cg.builder.CreateLShr(h1, sixteen, "h2")
	h3 := cg.builder.CreateXor(h1, h2, "h3")
	m1 := llvm.ConstInt(cg.context.Int64Type(), 0x85ebca6b, false)
	h4 := cg.builder.CreateMul(h3, m1, "h4")
	h5 := cg.builder.CreateLShr(h4, thirteen, "h5")
	h6 := cg.builder.CreateXor(h4, h5, "h6")
	m2 := llvm.ConstInt(cg.context.Int64Type(), 0xc2b2ae35, false)
	h7 := cg.builder.CreateMul(h6, m2, "h7")
	h8 := cg.builder.CreateLShr(h7, sixteen, "h8")
	result := cg.builder.CreateXor(h7, h8, "murmurresult")
	// Mask to 32 bits
	mask := llvm.ConstInt(cg.context.Int64Type(), 0xFFFFFFFF, false)
	return cg.builder.CreateAnd(result, mask, "murmur32"), nil
}

// ============================================================================
// TIME FUNCTIONS (additional)
// ============================================================================

// generateClock generates clock ticks
func (cg *LLVMCodeGenerator) generateClock(call *FunctionCall) (llvm.Value, error) {
	// Declare clock if not already
	if _, ok := cg.functions["clock"]; !ok {
		clockType := llvm.FunctionType(cg.context.Int64Type(), []llvm.Type{}, false)
		clockFn := llvm.AddFunction(cg.module, "clock", clockType)
		clockFn.SetLinkage(llvm.ExternalLinkage)
		cg.functions["clock"] = clockFn
	}

	clockFn := cg.functions["clock"]
	return cg.builder.CreateCall(clockFn.GlobalValueType(), clockFn, []llvm.Value{}, "clocktmp"), nil
}

// ============================================================================
// NUMBER CONVERSION FUNCTIONS (additional)
// ============================================================================

// generateToInt8 truncates to int8
func (cg *LLVMCodeGenerator) generateToInt8(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}
	val, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	trunc := cg.builder.CreateTrunc(val, cg.context.Int8Type(), "toint8")
	return cg.builder.CreateSExt(trunc, cg.context.Int64Type(), "int8ext"), nil
}

// generateToUint8 truncates to uint8
func (cg *LLVMCodeGenerator) generateToUint8(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}
	val, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	mask := llvm.ConstInt(cg.context.Int64Type(), 0xFF, false)
	return cg.builder.CreateAnd(val, mask, "touint8"), nil
}

// generateToInt16 truncates to int16
func (cg *LLVMCodeGenerator) generateToInt16(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}
	val, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	trunc := cg.builder.CreateTrunc(val, cg.context.Int16Type(), "toint16")
	return cg.builder.CreateSExt(trunc, cg.context.Int64Type(), "int16ext"), nil
}

// generateToUint16 truncates to uint16
func (cg *LLVMCodeGenerator) generateToUint16(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}
	val, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	mask := llvm.ConstInt(cg.context.Int64Type(), 0xFFFF, false)
	return cg.builder.CreateAnd(val, mask, "touint16"), nil
}

// generateToInt32 truncates to int32
func (cg *LLVMCodeGenerator) generateToInt32(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}
	val, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	trunc := cg.builder.CreateTrunc(val, cg.context.Int32Type(), "toint32")
	return cg.builder.CreateSExt(trunc, cg.context.Int64Type(), "int32ext"), nil
}

// generateToInt64 is identity for int64
func (cg *LLVMCodeGenerator) generateToInt64(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}
	return cg.generateExpression(call.Args[0])
}

// generateToUint64 is identity (just reinterpret)
func (cg *LLVMCodeGenerator) generateToUint64(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, nil
	}
	return cg.generateExpression(call.Args[0])
}

// ============================================================================
// SDL3 FUNCTIONS
// ============================================================================

// generateSDL3Init calls SDL_Init(flags) -> bool (SDL3: returns bool)
func (cg *LLVMCodeGenerator) generateSDL3Init(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("SDL3::init requires 1 argument")
	}

	flags, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Convert to int32 for SDL_Init
	flags32 := cg.builder.CreateTrunc(flags, cg.context.Int32Type(), "flags32")

	sdlInit := cg.functions["SDL_Init"]
	if sdlInit.IsNil() {
		return llvm.Value{}, fmt.Errorf("SDL_Init not declared")
	}

	// SDL3 returns bool (Int1Type)
	result := cg.builder.CreateCall(sdlInit.GlobalValueType(), sdlInit, []llvm.Value{flags32}, "sdlinit")
	// Convert bool to int64 (0 or 1)
	return cg.builder.CreateZExt(result, cg.context.Int64Type(), "sdlinitext"), nil
}

// generateSDL3Quit calls SDL_Quit()
func (cg *LLVMCodeGenerator) generateSDL3Quit(call *FunctionCall) (llvm.Value, error) {
	sdlQuit := cg.functions["SDL_Quit"]
	if sdlQuit.IsNil() {
		return llvm.Value{}, fmt.Errorf("SDL_Quit not declared")
	}

	cg.builder.CreateCall(sdlQuit.GlobalValueType(), sdlQuit, []llvm.Value{}, "")
	return llvm.Value{}, nil
}

// generateSDL3CreateWindow calls SDL_CreateWindow(title, w, h, flags) -> window* (SDL3: no x, y)
func (cg *LLVMCodeGenerator) generateSDL3CreateWindow(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 4 {
		return llvm.Value{}, fmt.Errorf("SDL3::create_window requires 4 arguments (title, w, h, flags)")
	}

	title, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	w, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	h, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}
	flags, err := cg.generateExpression(call.Args[3])
	if err != nil {
		return llvm.Value{}, err
	}

	// Convert to int32
	w32 := cg.builder.CreateTrunc(w, cg.context.Int32Type(), "w32")
	h32 := cg.builder.CreateTrunc(h, cg.context.Int32Type(), "h32")
	flags32 := cg.builder.CreateTrunc(flags, cg.context.Int32Type(), "flags32")

	sdlCreateWindow := cg.functions["SDL_CreateWindow"]
	if sdlCreateWindow.IsNil() {
		return llvm.Value{}, fmt.Errorf("SDL_CreateWindow not declared")
	}

	window := cg.builder.CreateCall(sdlCreateWindow.GlobalValueType(), sdlCreateWindow,
		[]llvm.Value{title, w32, h32, flags32}, "sdlwindow")
	return cg.builder.CreatePtrToInt(window, cg.context.Int64Type(), "windowptr"), nil
}

// generateSDL3DestroyWindow calls SDL_DestroyWindow(window*)
func (cg *LLVMCodeGenerator) generateSDL3DestroyWindow(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("SDL3::destroy_window requires 1 argument")
	}

	windowPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Convert int64 to pointer
	window := cg.builder.CreateIntToPtr(windowPtr, llvm.PointerType(cg.context.Int8Type(), 0), "window")

	sdlDestroyWindow := cg.functions["SDL_DestroyWindow"]
	if sdlDestroyWindow.IsNil() {
		return llvm.Value{}, fmt.Errorf("SDL_DestroyWindow not declared")
	}

	cg.builder.CreateCall(sdlDestroyWindow.GlobalValueType(), sdlDestroyWindow, []llvm.Value{window}, "")
	return llvm.Value{}, nil
}

// generateSDL2CreateRenderer calls SDL_CreateRenderer(window*, index, flags) -> renderer*
func (cg *LLVMCodeGenerator) generateSDL2CreateRenderer(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 2 {
		return llvm.Value{}, fmt.Errorf("SDL3::create_renderer requires 2 arguments")
	}

	windowPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	flags, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}

	// Convert int64 to pointer
	window := cg.builder.CreateIntToPtr(windowPtr, llvm.PointerType(cg.context.Int8Type(), 0), "window")
	flags32 := cg.builder.CreateTrunc(flags, cg.context.Int32Type(), "flags32")
	// -1 for default driver (SDL_RENDERER_DRIVER_INDEX)
	index32 := llvm.ConstInt(cg.context.Int32Type(), 4294967295, false) // 0xFFFFFFFF as uint64

	sdlCreateRenderer := cg.functions["SDL_CreateRenderer"]
	if sdlCreateRenderer.IsNil() {
		return llvm.Value{}, fmt.Errorf("SDL_CreateRenderer not declared")
	}

	renderer := cg.builder.CreateCall(sdlCreateRenderer.GlobalValueType(), sdlCreateRenderer,
		[]llvm.Value{window, index32, flags32}, "sdlrenderer")
	return cg.builder.CreatePtrToInt(renderer, cg.context.Int64Type(), "rendererptr"), nil
}

// generateSDL3DestroyRenderer calls SDL_DestroyRenderer(renderer*)
func (cg *LLVMCodeGenerator) generateSDL3DestroyRenderer(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("SDL3::destroy_renderer requires 1 argument")
	}

	rendererPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Convert int64 to pointer
	renderer := cg.builder.CreateIntToPtr(rendererPtr, llvm.PointerType(cg.context.Int8Type(), 0), "renderer")

	sdlDestroyRenderer := cg.functions["SDL_DestroyRenderer"]
	if sdlDestroyRenderer.IsNil() {
		return llvm.Value{}, fmt.Errorf("SDL_DestroyRenderer not declared")
	}

	cg.builder.CreateCall(sdlDestroyRenderer.GlobalValueType(), sdlDestroyRenderer, []llvm.Value{renderer}, "")
	return llvm.Value{}, nil
}

// generateSDL3RenderClear calls SDL_RenderClear(renderer*) -> bool (SDL3: returns bool)
func (cg *LLVMCodeGenerator) generateSDL3RenderClear(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("SDL3::render_clear requires 1 argument")
	}

	rendererPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Convert int64 to pointer
	renderer := cg.builder.CreateIntToPtr(rendererPtr, llvm.PointerType(cg.context.Int8Type(), 0), "renderer")

	sdlRenderClear := cg.functions["SDL_RenderClear"]
	if sdlRenderClear.IsNil() {
		return llvm.Value{}, fmt.Errorf("SDL_RenderClear not declared")
	}

	// SDL3 returns bool (Int1Type)
	result := cg.builder.CreateCall(sdlRenderClear.GlobalValueType(), sdlRenderClear, []llvm.Value{renderer}, "sdlclear")
	// Convert bool to int64 (0 or 1)
	return cg.builder.CreateZExt(result, cg.context.Int64Type(), "sdlclearext"), nil
}

// generateSDL3RenderPresent calls SDL_RenderPresent(renderer*)
func (cg *LLVMCodeGenerator) generateSDL3RenderPresent(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("SDL3::render_present requires 1 argument")
	}

	rendererPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Convert int64 to pointer
	renderer := cg.builder.CreateIntToPtr(rendererPtr, llvm.PointerType(cg.context.Int8Type(), 0), "renderer")

	sdlRenderPresent := cg.functions["SDL_RenderPresent"]
	if sdlRenderPresent.IsNil() {
		return llvm.Value{}, fmt.Errorf("SDL_RenderPresent not declared")
	}

	cg.builder.CreateCall(sdlRenderPresent.GlobalValueType(), sdlRenderPresent, []llvm.Value{renderer}, "")
	return llvm.Value{}, nil
}

// generateSDL3SetRenderDrawColor calls SDL_SetRenderDrawColor(renderer*, r, g, b, a) -> bool (SDL3: returns bool)
func (cg *LLVMCodeGenerator) generateSDL3SetRenderDrawColor(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 5 {
		return llvm.Value{}, fmt.Errorf("SDL3::set_render_draw_color requires 5 arguments")
	}

	rendererPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	r, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	g, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}
	b, err := cg.generateExpression(call.Args[3])
	if err != nil {
		return llvm.Value{}, err
	}
	a, err := cg.generateExpression(call.Args[4])
	if err != nil {
		return llvm.Value{}, err
	}

	// Convert int64 to pointer and truncate colors to uint8
	renderer := cg.builder.CreateIntToPtr(rendererPtr, llvm.PointerType(cg.context.Int8Type(), 0), "renderer")
	r8 := cg.builder.CreateTrunc(r, cg.context.Int8Type(), "r8")
	g8 := cg.builder.CreateTrunc(g, cg.context.Int8Type(), "g8")
	b8 := cg.builder.CreateTrunc(b, cg.context.Int8Type(), "b8")
	a8 := cg.builder.CreateTrunc(a, cg.context.Int8Type(), "a8")

	sdlSetRenderDrawColor := cg.functions["SDL_SetRenderDrawColor"]
	if sdlSetRenderDrawColor.IsNil() {
		return llvm.Value{}, fmt.Errorf("SDL_SetRenderDrawColor not declared")
	}

	// SDL3 returns bool (Int1Type)
	result := cg.builder.CreateCall(sdlSetRenderDrawColor.GlobalValueType(), sdlSetRenderDrawColor,
		[]llvm.Value{renderer, r8, g8, b8, a8}, "sdlcolor")
	// Convert bool to int64 (0 or 1)
	return cg.builder.CreateZExt(result, cg.context.Int64Type(), "sdlcolorext"), nil
}

// generateSDL3RenderDrawLine calls SDL_RenderDrawLine(renderer*, x1, y1, x2, y2) -> bool (SDL3: returns bool)
func (cg *LLVMCodeGenerator) generateSDL3RenderDrawLine(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 5 {
		return llvm.Value{}, fmt.Errorf("SDL3::render_draw_line requires 5 arguments")
	}

	rendererPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	x1, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	y1, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}
	x2, err := cg.generateExpression(call.Args[3])
	if err != nil {
		return llvm.Value{}, err
	}
	y2, err := cg.generateExpression(call.Args[4])
	if err != nil {
		return llvm.Value{}, err
	}

	// Convert int64 to pointer and truncate coordinates to int32
	renderer := cg.builder.CreateIntToPtr(rendererPtr, llvm.PointerType(cg.context.Int8Type(), 0), "renderer")
	x1_32 := cg.builder.CreateTrunc(x1, cg.context.Int32Type(), "x1_32")
	y1_32 := cg.builder.CreateTrunc(y1, cg.context.Int32Type(), "y1_32")
	x2_32 := cg.builder.CreateTrunc(x2, cg.context.Int32Type(), "x2_32")
	y2_32 := cg.builder.CreateTrunc(y2, cg.context.Int32Type(), "y2_32")

	sdlRenderDrawLine := cg.functions["SDL_RenderDrawLine"]
	if sdlRenderDrawLine.IsNil() {
		return llvm.Value{}, fmt.Errorf("SDL_RenderDrawLine not declared")
	}

	// SDL3 returns bool (Int1Type)
	result := cg.builder.CreateCall(sdlRenderDrawLine.GlobalValueType(), sdlRenderDrawLine,
		[]llvm.Value{renderer, x1_32, y1_32, x2_32, y2_32}, "sdlline")
	// Convert bool to int64 (0 or 1)
	return cg.builder.CreateZExt(result, cg.context.Int64Type(), "sdllineext"), nil
}

// generateSDL2RenderDrawRect calls SDL_RenderDrawRect(renderer*, x, y, w, h) -> int
// Creates SDL_Rect structure on stack and passes pointer to SDL2
func (cg *LLVMCodeGenerator) generateSDL2RenderDrawRect(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 5 {
		return llvm.Value{}, fmt.Errorf("SDL3::render_draw_rect requires 5 arguments")
	}

	rendererPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	x, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	y, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}
	w, err := cg.generateExpression(call.Args[3])
	if err != nil {
		return llvm.Value{}, err
	}
	h, err := cg.generateExpression(call.Args[4])
	if err != nil {
		return llvm.Value{}, err
	}

	// Create SDL_Rect structure: {int x, int y, int w, int h}
	sdlRectType := llvm.StructType([]llvm.Type{
		cg.context.Int32Type(), // x
		cg.context.Int32Type(), // y
		cg.context.Int32Type(), // w
		cg.context.Int32Type(), // h
	}, false)

	// Allocate SDL_Rect on stack
	rectPtr := cg.builder.CreateAlloca(sdlRectType, "rect")

	// Convert int64 to int32 and store in rect structure
	x32 := cg.builder.CreateTrunc(x, cg.context.Int32Type(), "x32")
	y32 := cg.builder.CreateTrunc(y, cg.context.Int32Type(), "y32")
	w32 := cg.builder.CreateTrunc(w, cg.context.Int32Type(), "w32")
	h32 := cg.builder.CreateTrunc(h, cg.context.Int32Type(), "h32")

	// Store values into rect structure
	xGEP := cg.builder.CreateStructGEP(sdlRectType, rectPtr, 0, "xgep")
	yGEP := cg.builder.CreateStructGEP(sdlRectType, rectPtr, 1, "ygep")
	wGEP := cg.builder.CreateStructGEP(sdlRectType, rectPtr, 2, "wgep")
	hGEP := cg.builder.CreateStructGEP(sdlRectType, rectPtr, 3, "hgep")

	cg.builder.CreateStore(x32, xGEP)
	cg.builder.CreateStore(y32, yGEP)
	cg.builder.CreateStore(w32, wGEP)
	cg.builder.CreateStore(h32, hGEP)

	// Convert int64 to pointer
	renderer := cg.builder.CreateIntToPtr(rendererPtr, llvm.PointerType(cg.context.Int8Type(), 0), "renderer")

	sdlRenderDrawRect := cg.functions["SDL_RenderDrawRect"]
	if sdlRenderDrawRect.IsNil() {
		return llvm.Value{}, fmt.Errorf("SDL_RenderDrawRect not declared")
	}

	// SDL3 returns bool (Int1Type)
	result := cg.builder.CreateCall(sdlRenderDrawRect.GlobalValueType(), sdlRenderDrawRect,
		[]llvm.Value{renderer, rectPtr}, "sdlrect")
	// Convert bool to int64 (0 or 1)
	return cg.builder.CreateZExt(result, cg.context.Int64Type(), "sdlrectext"), nil
}

// generateSDL2RenderFillRect calls SDL_RenderFillRect(renderer*, x, y, w, h) -> int
// Creates SDL_Rect structure on stack and passes pointer to SDL2
func (cg *LLVMCodeGenerator) generateSDL2RenderFillRect(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 5 {
		return llvm.Value{}, fmt.Errorf("SDL3::render_fill_rect requires 5 arguments")
	}

	rendererPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}
	x, err := cg.generateExpression(call.Args[1])
	if err != nil {
		return llvm.Value{}, err
	}
	y, err := cg.generateExpression(call.Args[2])
	if err != nil {
		return llvm.Value{}, err
	}
	w, err := cg.generateExpression(call.Args[3])
	if err != nil {
		return llvm.Value{}, err
	}
	h, err := cg.generateExpression(call.Args[4])
	if err != nil {
		return llvm.Value{}, err
	}

	// Create SDL_Rect structure: {int x, int y, int w, int h}
	sdlRectType := llvm.StructType([]llvm.Type{
		cg.context.Int32Type(), // x
		cg.context.Int32Type(), // y
		cg.context.Int32Type(), // w
		cg.context.Int32Type(), // h
	}, false)

	// Allocate SDL_Rect on stack
	rectPtr := cg.builder.CreateAlloca(sdlRectType, "rect")

	// Convert int64 to int32 and store in rect structure
	x32 := cg.builder.CreateTrunc(x, cg.context.Int32Type(), "x32")
	y32 := cg.builder.CreateTrunc(y, cg.context.Int32Type(), "y32")
	w32 := cg.builder.CreateTrunc(w, cg.context.Int32Type(), "w32")
	h32 := cg.builder.CreateTrunc(h, cg.context.Int32Type(), "h32")

	// Store values into rect structure
	xGEP := cg.builder.CreateStructGEP(sdlRectType, rectPtr, 0, "xgep")
	yGEP := cg.builder.CreateStructGEP(sdlRectType, rectPtr, 1, "ygep")
	wGEP := cg.builder.CreateStructGEP(sdlRectType, rectPtr, 2, "wgep")
	hGEP := cg.builder.CreateStructGEP(sdlRectType, rectPtr, 3, "hgep")

	cg.builder.CreateStore(x32, xGEP)
	cg.builder.CreateStore(y32, yGEP)
	cg.builder.CreateStore(w32, wGEP)
	cg.builder.CreateStore(h32, hGEP)

	// Convert int64 to pointer
	renderer := cg.builder.CreateIntToPtr(rendererPtr, llvm.PointerType(cg.context.Int8Type(), 0), "renderer")

	sdlRenderFillRect := cg.functions["SDL_RenderFillRect"]
	if sdlRenderFillRect.IsNil() {
		return llvm.Value{}, fmt.Errorf("SDL_RenderFillRect not declared")
	}

	// SDL3 returns bool (Int1Type)
	result := cg.builder.CreateCall(sdlRenderFillRect.GlobalValueType(), sdlRenderFillRect,
		[]llvm.Value{renderer, rectPtr}, "sdlfillrect")
	// Convert bool to int64 (0 or 1)
	return cg.builder.CreateZExt(result, cg.context.Int64Type(), "sdlfillrectext"), nil
}

// generateSDL2PollEvent calls SDL_PollEvent(event*) -> int (returns 1 if event, 0 if none)
func (cg *LLVMCodeGenerator) generateSDL2PollEvent(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("SDL3::poll_event requires 1 argument")
	}

	eventPtr, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Convert int64 to pointer
	event := cg.builder.CreateIntToPtr(eventPtr, llvm.PointerType(cg.context.Int8Type(), 0), "event")

	sdlPollEvent := cg.functions["SDL_PollEvent"]
	if sdlPollEvent.IsNil() {
		return llvm.Value{}, fmt.Errorf("SDL_PollEvent not declared")
	}

	result := cg.builder.CreateCall(sdlPollEvent.GlobalValueType(), sdlPollEvent, []llvm.Value{event}, "sdlevent")
	return cg.builder.CreateSExt(result, cg.context.Int64Type(), "sdleventext"), nil
}

// generateSDL2Delay calls SDL_Delay(ms)
func (cg *LLVMCodeGenerator) generateSDL2Delay(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("SDL3::delay requires 1 argument")
	}

	ms, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	// Convert to int32
	ms32 := cg.builder.CreateTrunc(ms, cg.context.Int32Type(), "ms32")

	sdlDelay := cg.functions["SDL_Delay"]
	if sdlDelay.IsNil() {
		return llvm.Value{}, fmt.Errorf("SDL_Delay not declared")
	}

	cg.builder.CreateCall(sdlDelay.GlobalValueType(), sdlDelay, []llvm.Value{ms32}, "")
	return llvm.Value{}, nil
}

// generateSDL3GetTicks calls SDL_GetTicks() -> int
func (cg *LLVMCodeGenerator) generateSDL3GetTicks(call *FunctionCall) (llvm.Value, error) {
	sdlGetTicks := cg.functions["SDL_GetTicks"]
	if sdlGetTicks.IsNil() {
		return llvm.Value{}, fmt.Errorf("SDL_GetTicks not declared")
	}

	result := cg.builder.CreateCall(sdlGetTicks.GlobalValueType(), sdlGetTicks, []llvm.Value{}, "sdlticks")
	return cg.builder.CreateSExt(result, cg.context.Int64Type(), "sdlsticksext"), nil
}
