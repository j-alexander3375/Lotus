package main

// llvm_types.go - LLVM type system utilities for Lotus
//
// This file provides type mapping and conversion utilities between
// Lotus types and LLVM IR types.

import (
	"tinygo.org/x/go-llvm"
)

// LotusTypeInfo represents metadata about a Lotus type
type LotusTypeInfo struct {
	Name      string    // Lotus type name
	LLVMType  llvm.Type // Corresponding LLVM type
	Size      int       // Size in bytes
	Alignment int       // Alignment in bytes
	IsSigned  bool      // Whether it's a signed integer type
	IsPointer bool      // Whether it's a pointer type
	IsFloat   bool      // Whether it's a floating point type
}

// TypeMapper handles type conversions between Lotus and LLVM
type TypeMapper struct {
	context   llvm.Context
	typeCache map[string]LotusTypeInfo
}

// NewTypeMapper creates a new type mapper
func NewTypeMapper(ctx llvm.Context) *TypeMapper {
	tm := &TypeMapper{
		context:   ctx,
		typeCache: make(map[string]LotusTypeInfo),
	}
	tm.initBuiltinTypes()
	return tm
}

// initBuiltinTypes initializes all builtin Lotus types
func (tm *TypeMapper) initBuiltinTypes() {
	// Void type
	tm.typeCache["void"] = LotusTypeInfo{
		Name:      "void",
		LLVMType:  tm.context.VoidType(),
		Size:      0,
		Alignment: 0,
		IsSigned:  false,
		IsPointer: false,
		IsFloat:   false,
	}

	// Boolean type
	tm.typeCache["bool"] = LotusTypeInfo{
		Name:      "bool",
		LLVMType:  tm.context.Int1Type(),
		Size:      1,
		Alignment: 1,
		IsSigned:  false,
		IsPointer: false,
		IsFloat:   false,
	}

	// Integer types
	tm.typeCache["int8"] = LotusTypeInfo{
		Name:      "int8",
		LLVMType:  tm.context.Int8Type(),
		Size:      1,
		Alignment: 1,
		IsSigned:  true,
		IsPointer: false,
		IsFloat:   false,
	}
	tm.typeCache["uint8"] = LotusTypeInfo{
		Name:      "uint8",
		LLVMType:  tm.context.Int8Type(),
		Size:      1,
		Alignment: 1,
		IsSigned:  false,
		IsPointer: false,
		IsFloat:   false,
	}
	tm.typeCache["char"] = tm.typeCache["int8"] // Alias

	tm.typeCache["int16"] = LotusTypeInfo{
		Name:      "int16",
		LLVMType:  tm.context.Int16Type(),
		Size:      2,
		Alignment: 2,
		IsSigned:  true,
		IsPointer: false,
		IsFloat:   false,
	}
	tm.typeCache["uint16"] = LotusTypeInfo{
		Name:      "uint16",
		LLVMType:  tm.context.Int16Type(),
		Size:      2,
		Alignment: 2,
		IsSigned:  false,
		IsPointer: false,
		IsFloat:   false,
	}

	tm.typeCache["int32"] = LotusTypeInfo{
		Name:      "int32",
		LLVMType:  tm.context.Int32Type(),
		Size:      4,
		Alignment: 4,
		IsSigned:  true,
		IsPointer: false,
		IsFloat:   false,
	}
	tm.typeCache["uint32"] = LotusTypeInfo{
		Name:      "uint32",
		LLVMType:  tm.context.Int32Type(),
		Size:      4,
		Alignment: 4,
		IsSigned:  false,
		IsPointer: false,
		IsFloat:   false,
	}

	tm.typeCache["int64"] = LotusTypeInfo{
		Name:      "int64",
		LLVMType:  tm.context.Int64Type(),
		Size:      8,
		Alignment: 8,
		IsSigned:  true,
		IsPointer: false,
		IsFloat:   false,
	}
	tm.typeCache["uint64"] = LotusTypeInfo{
		Name:      "uint64",
		LLVMType:  tm.context.Int64Type(),
		Size:      8,
		Alignment: 8,
		IsSigned:  false,
		IsPointer: false,
		IsFloat:   false,
	}

	// Default 'int' is int64
	tm.typeCache["int"] = tm.typeCache["int64"]
	tm.typeCache["uint"] = tm.typeCache["uint64"]

	// Floating point types
	tm.typeCache["float32"] = LotusTypeInfo{
		Name:      "float32",
		LLVMType:  tm.context.FloatType(),
		Size:      4,
		Alignment: 4,
		IsSigned:  true,
		IsPointer: false,
		IsFloat:   true,
	}
	tm.typeCache["float64"] = LotusTypeInfo{
		Name:      "float64",
		LLVMType:  tm.context.DoubleType(),
		Size:      8,
		Alignment: 8,
		IsSigned:  true,
		IsPointer: false,
		IsFloat:   true,
	}
	tm.typeCache["float"] = tm.typeCache["float64"] // Alias

	// String type (pointer to i8)
	tm.typeCache["string"] = LotusTypeInfo{
		Name:      "string",
		LLVMType:  llvm.PointerType(tm.context.Int8Type(), 0),
		Size:      8,
		Alignment: 8,
		IsSigned:  false,
		IsPointer: true,
		IsFloat:   false,
	}
}

// GetType returns the LLVM type for a Lotus type name
func (tm *TypeMapper) GetType(name string) (llvm.Type, bool) {
	info, ok := tm.typeCache[name]
	if ok {
		return info.LLVMType, true
	}
	return llvm.Type{}, false
}

// GetTypeInfo returns full type information for a Lotus type
func (tm *TypeMapper) GetTypeInfo(name string) (LotusTypeInfo, bool) {
	info, ok := tm.typeCache[name]
	return info, ok
}

// GetPointerType returns a pointer type to the given base type
func (tm *TypeMapper) GetPointerType(baseType string) llvm.Type {
	if info, ok := tm.typeCache[baseType]; ok {
		return llvm.PointerType(info.LLVMType, 0)
	}
	// Default to i8*
	return llvm.PointerType(tm.context.Int8Type(), 0)
}

// GetArrayType returns an array type
func (tm *TypeMapper) GetArrayType(elementType string, size int) llvm.Type {
	if info, ok := tm.typeCache[elementType]; ok {
		return llvm.ArrayType(info.LLVMType, size)
	}
	// Default to i64 array
	return llvm.ArrayType(tm.context.Int64Type(), size)
}

// CreateStructType creates a new struct type
func (tm *TypeMapper) CreateStructType(name string, fieldTypes []llvm.Type) llvm.Type {
	structType := tm.context.StructCreateNamed(name)
	structType.StructSetBody(fieldTypes, false)
	return structType
}

// IsIntegerType checks if a type is an integer type
func (tm *TypeMapper) IsIntegerType(name string) bool {
	info, ok := tm.typeCache[name]
	if !ok {
		return false
	}
	return !info.IsFloat && !info.IsPointer && name != "void" && name != "bool"
}

// IsSignedType checks if a type is a signed integer type
func (tm *TypeMapper) IsSignedType(name string) bool {
	info, ok := tm.typeCache[name]
	if !ok {
		return false
	}
	return info.IsSigned
}

// IsFloatType checks if a type is a floating point type
func (tm *TypeMapper) IsFloatType(name string) bool {
	info, ok := tm.typeCache[name]
	if !ok {
		return false
	}
	return info.IsFloat
}

// IsPointerType checks if a type is a pointer type
func (tm *TypeMapper) IsPointerType(name string) bool {
	info, ok := tm.typeCache[name]
	if !ok {
		return false
	}
	return info.IsPointer
}

// GetTypeSize returns the size of a type in bytes
func (tm *TypeMapper) GetTypeSize(name string) int {
	info, ok := tm.typeCache[name]
	if !ok {
		return 8 // Default to 8 bytes
	}
	return info.Size
}

// GetTypeAlignment returns the alignment of a type in bytes
func (tm *TypeMapper) GetTypeAlignment(name string) int {
	info, ok := tm.typeCache[name]
	if !ok {
		return 8 // Default to 8 bytes
	}
	return info.Alignment
}

// TypeConversion represents a type conversion operation
type TypeConversion struct {
	FromType  string
	ToType    string
	Converter func(llvm.Builder, llvm.Value) llvm.Value
}

// GetConversion returns a conversion function between two types
func (tm *TypeMapper) GetConversion(from, to string) (func(llvm.Builder, llvm.Value) llvm.Value, bool) {
	fromInfo, fromOk := tm.typeCache[from]
	toInfo, toOk := tm.typeCache[to]

	if !fromOk || !toOk {
		return nil, false
	}

	// Same type - no conversion needed
	if from == to {
		return func(b llvm.Builder, v llvm.Value) llvm.Value { return v }, true
	}

	// Integer to integer
	if tm.IsIntegerType(from) && tm.IsIntegerType(to) {
		if fromInfo.Size < toInfo.Size {
			// Extend
			if fromInfo.IsSigned {
				return func(b llvm.Builder, v llvm.Value) llvm.Value {
					return b.CreateSExt(v, toInfo.LLVMType, "sext")
				}, true
			}
			return func(b llvm.Builder, v llvm.Value) llvm.Value {
				return b.CreateZExt(v, toInfo.LLVMType, "zext")
			}, true
		}
		// Truncate
		return func(b llvm.Builder, v llvm.Value) llvm.Value {
			return b.CreateTrunc(v, toInfo.LLVMType, "trunc")
		}, true
	}

	// Integer to float
	if tm.IsIntegerType(from) && tm.IsFloatType(to) {
		if fromInfo.IsSigned {
			return func(b llvm.Builder, v llvm.Value) llvm.Value {
				return b.CreateSIToFP(v, toInfo.LLVMType, "sitofp")
			}, true
		}
		return func(b llvm.Builder, v llvm.Value) llvm.Value {
			return b.CreateUIToFP(v, toInfo.LLVMType, "uitofp")
		}, true
	}

	// Float to integer
	if tm.IsFloatType(from) && tm.IsIntegerType(to) {
		if toInfo.IsSigned {
			return func(b llvm.Builder, v llvm.Value) llvm.Value {
				return b.CreateFPToSI(v, toInfo.LLVMType, "fptosi")
			}, true
		}
		return func(b llvm.Builder, v llvm.Value) llvm.Value {
			return b.CreateFPToUI(v, toInfo.LLVMType, "fptoui")
		}, true
	}

	// Float to float
	if tm.IsFloatType(from) && tm.IsFloatType(to) {
		if fromInfo.Size < toInfo.Size {
			return func(b llvm.Builder, v llvm.Value) llvm.Value {
				return b.CreateFPExt(v, toInfo.LLVMType, "fpext")
			}, true
		}
		return func(b llvm.Builder, v llvm.Value) llvm.Value {
			return b.CreateFPTrunc(v, toInfo.LLVMType, "fptrunc")
		}, true
	}

	return nil, false
}
