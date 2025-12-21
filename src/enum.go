package main

import "fmt"

// EnumDefinition represents an enum type definition
type EnumDefinition struct {
	Name   string
	Values []EnumValue
}

func (e *EnumDefinition) astNode() {}

// EnumValue represents a single enum constant
type EnumValue struct {
	Name  string
	Value int // explicit value or auto-assigned
}

// EnumLiteral represents an enum value reference
type EnumLiteral struct {
	EnumName  string
	ValueName string
}

func (e *EnumLiteral) astNode() {}

// EnumRegistry stores defined enum types
var EnumRegistry = make(map[string]*EnumDefinition)

// generateEnumDefinition registers an enum type
func (cg *CodeGenerator) generateEnumDefinition(def *EnumDefinition) {
	// Auto-assign values if not explicitly set
	nextValue := 0
	for i := range def.Values {
		if def.Values[i].Value == 0 && i > 0 {
			def.Values[i].Value = nextValue
		}
		nextValue = def.Values[i].Value + 1
	}
	
	// Register enum in global registry
	EnumRegistry[def.Name] = def
	
	// Enums compile to constants, no runtime code needed
}

// generateEnumLiteral generates assembly for enum value reference
func (cg *CodeGenerator) generateEnumLiteral(lit *EnumLiteral) {
	enumDef, exists := EnumRegistry[lit.EnumName]
	if !exists {
		fmt.Printf("Error: enum type '%s' not defined\n", lit.EnumName)
		return
	}
	
	// Find the enum value
	for _, val := range enumDef.Values {
		if val.Name == lit.ValueName {
			// Load the constant value into rax
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rax\n", val.Value))
			return
		}
	}
	
	fmt.Printf("Error: enum value '%s' not found in enum '%s'\n", lit.ValueName, lit.EnumName)
}

// getEnumValue looks up an enum value by name
func getEnumValue(enumName, valueName string) (int, bool) {
	enumDef, exists := EnumRegistry[enumName]
	if !exists {
		return 0, false
	}
	
	for _, val := range enumDef.Values {
		if val.Name == valueName {
			return val.Value, true
		}
	}
	
	return 0, false
}

// Example enum syntax:
// enum Color {
//     Red = 0,
//     Green = 1,
//     Blue = 2
// }
//
// enum Status {
//     Ok,        // 0
//     Error,     // 1
//     Pending    // 2
// }
//
// Usage:
// int color = Color.Red;
// int status = Status.Ok;
