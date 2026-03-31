package main

// EnumDefinition represents an enum type definition
type EnumDefinition struct {
	BaseNode
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
	BaseNode
	EnumName  string
	ValueName string
}

func (e *EnumLiteral) astNode() {}

// EnumRegistry stores defined enum types
var EnumRegistry = make(map[string]*EnumDefinition)

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
