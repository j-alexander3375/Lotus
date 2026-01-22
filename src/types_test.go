package main

import (
	"testing"
)

// TestGetTypeSize tests that type sizes are correctly returned
func TestGetTypeSize(t *testing.T) {
	tests := []struct {
		tokenType TokenType
		expected  int
	}{
		{TokenTypeInt8, Int8Size},
		{TokenTypeUint8, Int8Size},
		{TokenTypeBool, Int8Size},
		{TokenTypeInt16, Int16Size},
		{TokenTypeUint16, Int16Size},
		{TokenTypeInt, Int32Size},
		{TokenTypeInt32, Int32Size},
		{TokenTypeUint, Int32Size},
		{TokenTypeUint32, Int32Size},
		{TokenTypeFloat32, Int32Size},
		{TokenTypeChar, Int32Size},
		{TokenTypeInt64, Int64Size},
		{TokenTypeUint64, Int64Size},
		{TokenTypeFloat64, Int64Size},
		{TokenTypeString, Int64Size},
	}

	for _, tt := range tests {
		t.Run(TokenTypeName(tt.tokenType), func(t *testing.T) {
			size := GetTypeSize(tt.tokenType)
			if size != tt.expected {
				t.Errorf("GetTypeSize(%v) = %d, expected %d", tt.tokenType, size, tt.expected)
			}
		})
	}
}

// TestIsIntegerType tests integer type checking
func TestIsIntegerType(t *testing.T) {
	integerTypes := []TokenType{
		TokenTypeInt, TokenTypeInt8, TokenTypeInt16, TokenTypeInt32, TokenTypeInt64,
		TokenTypeUint, TokenTypeUint8, TokenTypeUint16, TokenTypeUint32, TokenTypeUint64,
	}

	nonIntegerTypes := []TokenType{
		TokenTypeFloat32, TokenTypeFloat64, TokenTypeBool, TokenTypeString, TokenTypeChar, TokenTypeVoid,
	}

	for _, tt := range integerTypes {
		t.Run("Integer_"+TokenTypeName(tt), func(t *testing.T) {
			if !IsIntegerType(tt) {
				t.Errorf("IsIntegerType(%v) = false, expected true", tt)
			}
		})
	}

	for _, tt := range nonIntegerTypes {
		t.Run("NonInteger_"+TokenTypeName(tt), func(t *testing.T) {
			if IsIntegerType(tt) {
				t.Errorf("IsIntegerType(%v) = true, expected false", tt)
			}
		})
	}
}

// TestIsFloatType tests float type checking
func TestIsFloatType(t *testing.T) {
	floatTypes := []TokenType{
		TokenTypeFloat32, TokenTypeFloat64,
	}

	nonFloatTypes := []TokenType{
		TokenTypeInt, TokenTypeInt8, TokenTypeInt16, TokenTypeInt32, TokenTypeInt64,
		TokenTypeUint, TokenTypeUint8, TokenTypeUint16, TokenTypeUint32, TokenTypeUint64,
		TokenTypeBool, TokenTypeString, TokenTypeChar, TokenTypeVoid,
	}

	for _, tt := range floatTypes {
		t.Run("Float_"+TokenTypeName(tt), func(t *testing.T) {
			if !IsFloatType(tt) {
				t.Errorf("IsFloatType(%v) = false, expected true", tt)
			}
		})
	}

	for _, tt := range nonFloatTypes {
		t.Run("NonFloat_"+TokenTypeName(tt), func(t *testing.T) {
			if IsFloatType(tt) {
				t.Errorf("IsFloatType(%v) = true, expected false", tt)
			}
		})
	}
}

// TestIsNumericType tests numeric type checking (int or float)
func TestIsNumericType(t *testing.T) {
	numericTypes := []TokenType{
		TokenTypeInt, TokenTypeInt8, TokenTypeInt16, TokenTypeInt32, TokenTypeInt64,
		TokenTypeUint, TokenTypeUint8, TokenTypeUint16, TokenTypeUint32, TokenTypeUint64,
		TokenTypeFloat32, TokenTypeFloat64,
	}

	nonNumericTypes := []TokenType{
		TokenTypeBool, TokenTypeString, TokenTypeChar, TokenTypeVoid,
	}

	for _, tt := range numericTypes {
		t.Run("Numeric_"+TokenTypeName(tt), func(t *testing.T) {
			if !IsNumericType(tt) {
				t.Errorf("IsNumericType(%v) = false, expected true", tt)
			}
		})
	}

	for _, tt := range nonNumericTypes {
		t.Run("NonNumeric_"+TokenTypeName(tt), func(t *testing.T) {
			if IsNumericType(tt) {
				t.Errorf("IsNumericType(%v) = true, expected false", tt)
			}
		})
	}
}

// TestIsPrimitiveType tests primitive type checking
func TestIsPrimitiveType(t *testing.T) {
	primitiveTypes := []TokenType{
		TokenTypeInt, TokenTypeInt8, TokenTypeInt16, TokenTypeInt32, TokenTypeInt64,
		TokenTypeUint, TokenTypeUint8, TokenTypeUint16, TokenTypeUint32, TokenTypeUint64,
		TokenTypeFloat32, TokenTypeFloat64,
		TokenTypeBool, TokenTypeString, TokenTypeChar,
	}

	nonPrimitiveTypes := []TokenType{
		TokenTypeVoid, TokenTypeFn,
	}

	for _, tt := range primitiveTypes {
		t.Run("Primitive_"+TokenTypeName(tt), func(t *testing.T) {
			if !IsPrimitiveType(tt) {
				t.Errorf("IsPrimitiveType(%v) = false, expected true", tt)
			}
		})
	}

	for _, tt := range nonPrimitiveTypes {
		t.Run("NonPrimitive_"+TokenTypeName(tt), func(t *testing.T) {
			if IsPrimitiveType(tt) {
				t.Errorf("IsPrimitiveType(%v) = true, expected false", tt)
			}
		})
	}
}

// TestTypeRegistry tests the type registry creation and operations
func TestTypeRegistry(t *testing.T) {
	registry := NewTypeRegistry()

	if registry == nil {
		t.Fatal("NewTypeRegistry() returned nil")
	}

	if registry.Structs == nil {
		t.Error("TypeRegistry.Structs is nil")
	}

	if registry.Enums == nil {
		t.Error("TypeRegistry.Enums is nil")
	}

	if registry.Classes == nil {
		t.Error("TypeRegistry.Classes is nil")
	}

	// Test that maps are initialized and empty
	if len(registry.Structs) != 0 {
		t.Errorf("Expected empty Structs map, got %d entries", len(registry.Structs))
	}

	if len(registry.Enums) != 0 {
		t.Errorf("Expected empty Enums map, got %d entries", len(registry.Enums))
	}

	if len(registry.Classes) != 0 {
		t.Errorf("Expected empty Classes map, got %d entries", len(registry.Classes))
	}
}

// TestVariable tests the Variable struct
func TestVariable(t *testing.T) {
	v := Variable{
		Name:   "test_var",
		Type:   TokenTypeInt32,
		Offset: 8,
	}

	if v.Name != "test_var" {
		t.Errorf("Expected Name 'test_var', got '%s'", v.Name)
	}

	if v.Type != TokenTypeInt32 {
		t.Errorf("Expected Type TokenTypeInt32, got %v", v.Type)
	}

	if v.Offset != 8 {
		t.Errorf("Expected Offset 8, got %d", v.Offset)
	}
}
