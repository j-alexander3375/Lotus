package main

// types.go - Type registry and type system utilities

// Variable represents a variable in the symbol table with its metadata
type Variable struct {
	Name   string    // Variable identifier
	Type   TokenType // Data type
	Offset int       // Stack offset from rbp for local variables
}

// TypeRegistry manages all type definitions in the compiler
type TypeRegistry struct {
	Structs map[string]*StructDefinition
	Enums   map[string]*EnumDefinition
	Classes map[string]*ClassDefinition
}

// NewTypeRegistry creates a new type registry instance
func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{
		Structs: make(map[string]*StructDefinition),
		Enums:   make(map[string]*EnumDefinition),
		Classes: make(map[string]*ClassDefinition),
	}
}

// GetTypeSize returns the size in bytes of a given type
func GetTypeSize(tokenType TokenType) int {
	switch tokenType {
	case TokenTypeInt8, TokenTypeUint8, TokenTypeBool:
		return Int8Size
	case TokenTypeInt16, TokenTypeUint16:
		return Int16Size
	case TokenTypeInt, TokenTypeInt32, TokenTypeUint, TokenTypeUint32, TokenTypeFloat:
		return Int32Size
	case TokenTypeInt64, TokenTypeUint64, TokenTypeString:
		return Int64Size
	default:
		return PointerSize // Default to pointer size for unknown types
	}
}

// IsIntegerType checks if a token type represents an integer
func IsIntegerType(tokenType TokenType) bool {
	switch tokenType {
	case TokenTypeInt, TokenTypeInt8, TokenTypeInt16, TokenTypeInt32, TokenTypeInt64,
		TokenTypeUint, TokenTypeUint8, TokenTypeUint16, TokenTypeUint32, TokenTypeUint64:
		return true
	default:
		return false
	}
}

// IsNumericType checks if a token type represents a numeric type (int or float)
func IsNumericType(tokenType TokenType) bool {
	return IsIntegerType(tokenType) || tokenType == TokenTypeFloat
}

// IsPrimitiveType checks if a token type is a primitive type
func IsPrimitiveType(tokenType TokenType) bool {
	return IsNumericType(tokenType) || tokenType == TokenTypeBool || tokenType == TokenTypeString
}
