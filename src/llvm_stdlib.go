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

// generatePow generates inline code for integer pow(base, exp) using repeated multiplication
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

	// Use LLVM pow intrinsic with double, then convert back to int
	// First, declare llvm.pow.f64 if not already declared
	if _, ok := cg.functions["llvm.pow.f64"]; !ok {
		powType := llvm.FunctionType(
			cg.context.DoubleType(),
			[]llvm.Type{cg.context.DoubleType(), cg.context.DoubleType()},
			false,
		)
		powFn := llvm.AddFunction(cg.module, "llvm.pow.f64", powType)
		cg.functions["llvm.pow.f64"] = powFn
	}

	// Convert ints to doubles
	baseF := cg.builder.CreateSIToFP(base, cg.context.DoubleType(), "powbaseconv")
	expF := cg.builder.CreateSIToFP(exp, cg.context.DoubleType(), "powexpconv")

	// Call pow intrinsic
	powFn := cg.functions["llvm.pow.f64"]
	result := cg.builder.CreateCall(powFn.GlobalValueType(), powFn, []llvm.Value{baseF, expF}, "powcall")

	// Convert back to int
	return cg.builder.CreateFPToSI(result, cg.context.Int64Type(), "powint"), nil
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

// generateQueueIntDequeue removes from front (simplified - shifts all elements)
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

	// Get first element (at data offset 0)
	dataBase := cg.builder.CreateAdd(arrPtr, sixteen, "database")
	firstPtr := cg.builder.CreateIntToPtr(dataBase, llvm.PointerType(cg.context.Int64Type(), 0), "firstptr")
	firstVal := cg.builder.CreateLoad(cg.context.Int64Type(), firstPtr, "firstval")

	// Decrement length (simplified - doesn't actually shift elements)
	lenOffset := cg.builder.CreateAdd(arrPtr, eight, "lenoffset")
	lenPtr := cg.builder.CreateIntToPtr(lenOffset, llvm.PointerType(cg.context.Int64Type(), 0), "lenptr")
	length := cg.builder.CreateLoad(cg.context.Int64Type(), lenPtr, "arrlen")
	newLen := cg.builder.CreateSub(length, one, "newlen")
	cg.builder.CreateStore(newLen, lenPtr)

	// Note: A real implementation would shift all elements or use a ring buffer
	// This is a simplified version that works for basic tests
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
	// Simplified: just push to back for now
	return cg.generateArrayIntPush(call)
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
	// Simplified: just push without heapify for now
	return cg.generateArrayIntPush(call)
}

func (cg *LLVMCodeGenerator) generateHeapIntPop(call *FunctionCall) (llvm.Value, error) {
	// Simplified: pop from front (should be min in a proper min-heap)
	return cg.generateQueueIntDequeue(call)
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

	one := llvm.ConstInt(cg.context.Int64Type(), 1, false)
	nextStart := cg.builder.CreateAdd(startI, one, "nextstart")
	cg.builder.CreateStore(nextStart, startIdx)

	cg.builder.CreateCondBr(isWS, startLoopBB, startExitBB)

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

	cg.builder.CreateStore(prevEnd, endIdx)
	cg.builder.CreateCondBr(isWS2, endLoopBB, endExitBB)

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
