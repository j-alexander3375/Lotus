package main

// llvm_json.go - JSON module implementation for LLVM backend
//
// This file implements JSON parsing and serialization for the LLVM backend.
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

import (
	"tinygo.org/x/go-llvm"
)

// ============================================================================
// JSON MODULE TYPES
// ============================================================================

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
	// For now, just detect the type and store it
	cg.generateJSONParseInline(str, jsonPtrInt)

	return jsonPtrInt, nil
}

// generateJSONParseInline generates inline code to parse JSON
func (cg *LLVMCodeGenerator) generateJSONParseInline(str llvm.Value, jsonPtr llvm.Value) {
	// This is a simplified parser that detects the JSON type
	// Full implementation would recursively parse nested structures

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
	cg.builder.CreateCondBr(isBool, objBlock, numBlock) // Simplified: skip to next check

	// Number detection (digit or minus)
	cg.builder.SetInsertPointAtEnd(numBlock)
	isMinus := cg.builder.CreateICmp(llvm.IntEQ, firstChar, llvm.ConstInt(cg.context.Int8Type(), '-', false), "isminus")
	isDigit := cg.builder.CreateAnd(
		cg.builder.CreateICmp(llvm.IntSGE, firstChar, llvm.ConstInt(cg.context.Int8Type(), '0', false), "ge0"),
		cg.builder.CreateICmp(llvm.IntSLE, firstChar, llvm.ConstInt(cg.context.Int8Type(), '9', false), "le9"),
		"isdigit")
	isNum := cg.builder.CreateOr(isMinus, isDigit, "isnum")
	cg.builder.CreateCondBr(isNum, strBlock, strBlock) // Continue to string check

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

	// Allocate output buffer (1024 bytes for now)
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
	// Simplified: just write "null" for now
	// Full implementation would handle all types

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

	// Other types
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
