package main

import "fmt"

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

// generateStructDefinition registers a struct type
func (cg *CodeGenerator) generateStructDefinition(def *StructDefinition) {
	// Calculate field offsets
	offset := 0
	for i := range def.Fields {
		def.Fields[i].Offset = offset
		fieldSize := getTypeSizeFromToken(def.Fields[i].Type)
		offset += fieldSize
	}

	// Register struct in global registry
	StructRegistry[def.Name] = def

	// No assembly generation needed for type definition
}

// generateStructLiteral generates assembly for struct literal initialization
func (cg *CodeGenerator) generateStructLiteral(lit *StructLiteral) {
	structDef, exists := StructRegistry[lit.StructName]
	if !exists {
		fmt.Printf("Error: struct type '%s' not defined\n", lit.StructName)
		return
	}

	// Calculate total struct size
	totalSize := 0
	for _, field := range structDef.Fields {
		fieldSize := getTypeSizeFromToken(field.Type)
		totalSize += fieldSize
	}

	// Allocate memory for struct on heap
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdi\n", totalSize))
	cg.textSection.WriteString("    call malloc@PLT\n")
	cg.textSection.WriteString("    pushq %rax\n") // Save struct pointer

	// Initialize each field
	for fieldName, valueExpr := range lit.Fields {
		// Find field definition
		var fieldDef *StructField
		for i := range structDef.Fields {
			if structDef.Fields[i].Name == fieldName {
				fieldDef = &structDef.Fields[i]
				break
			}
		}

		if fieldDef == nil {
			fmt.Printf("Error: field '%s' not found in struct '%s'\n", fieldName, lit.StructName)
			continue
		}

		// Evaluate field value
		cg.generateExpressionToReg(valueExpr, "rcx")

		// Get struct pointer
		cg.textSection.WriteString("    movq (%rsp), %rax\n")

		// Store value at field offset
		fieldSize := getTypeSizeFromToken(fieldDef.Type)
		switch fieldSize {
		case 1:
			cg.textSection.WriteString(fmt.Sprintf("    movb %%cl, %d(%%rax)\n", fieldDef.Offset))
		case 2:
			cg.textSection.WriteString(fmt.Sprintf("    movw %%cx, %d(%%rax)\n", fieldDef.Offset))
		case 4:
			cg.textSection.WriteString(fmt.Sprintf("    movl %%ecx, %d(%%rax)\n", fieldDef.Offset))
		case 8:
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rcx, %d(%%rax)\n", fieldDef.Offset))
		}
	}

	// Pop struct pointer back to rax
	cg.textSection.WriteString("    popq %rax\n")
}

// generateFieldAccess generates assembly for struct field access
func (cg *CodeGenerator) generateFieldAccess(access *FieldAccess) {
	// Evaluate object to get struct pointer
	cg.generateExpressionToReg(access.Object, "rax")

	// If using -> operator, rax already contains pointer
	// If using . operator, need to get address
	if !access.IsPointer {
		// For . operator, object should be on stack
		// This is a simplification - full implementation would need type information
		if ident, ok := access.Object.(*Identifier); ok {
			if v, exists := cg.variables[ident.Name]; exists {
				cg.textSection.WriteString(fmt.Sprintf("    leaq -%d(%%rbp), %%rax\n", v.Offset))
			}
		}
	}

	// Get struct type from object
	// For now, we'll need to track this in the type system
	// Simplified: assume we know the struct type

	// Find field offset (this requires type information)
	// For demonstration, we'll use a placeholder
	fieldOffset := 0 // Would need to look this up from struct definition

	// Load field value
	cg.textSection.WriteString(fmt.Sprintf("    movq %d(%%rax), %%rax\n", fieldOffset))
}

// Helper to get field offset in struct
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
