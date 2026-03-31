package main

// StructDefinition represents a struct type definition
type StructDefinition struct {
	BaseNode
	Name   string
	Fields []StructField
}

func (s *StructDefinition) astNode() {}

// StructField represents a field in a struct
type StructField struct {
	Name   string
	Type   TokenType
	Offset int // byte offset from struct start
}

// StructLiteral represents struct initialization
type StructLiteral struct {
	BaseNode
	StructName string
	Fields     map[string]ASTNode // field name -> value
}

func (s *StructLiteral) astNode() {}

// FieldAccess represents struct field access: obj.field or obj->field
type FieldAccess struct {
	BaseNode
	Object    ASTNode
	FieldName string
	IsPointer bool // true for ->, false for .
}

func (f *FieldAccess) astNode() {}

// StructRegistry stores defined struct types
var StructRegistry = make(map[string]*StructDefinition)

// getFieldOffset looks up field offset in struct
func getFieldOffset(structName, fieldName string) (int, bool) {
	structDef, exists := StructRegistry[structName]
	if !exists {
		return 0, false
	}

	for _, field := range structDef.Fields {
		if field.Name == fieldName {
			return field.Offset, true
		}
	}

	return 0, false
}
