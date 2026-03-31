package main

// ClassDefinition represents a class type definition
type ClassDefinition struct {
	BaseNode
	Name    string
	Fields  []ClassField
	Methods []ClassMethod
}

func (c *ClassDefinition) astNode() {}

// ClassField represents a field in a class
type ClassField struct {
	Name   string
	Type   TokenType
	Offset int // byte offset from class instance start
}

// ClassMethod represents a method in a class
type ClassMethod struct {
	Name       string
	Params     []FunctionParam
	ReturnType TokenType
	Body       []ASTNode
	Label      string // label for method
}

// ClassLiteral represents class instantiation (new ClassName)
type ClassLiteral struct {
	BaseNode
	ClassName string
	Fields    map[string]ASTNode // field name -> initial value
}

func (c *ClassLiteral) astNode() {}

// MethodCall represents a method call: obj.method() or obj->method()
type MethodCall struct {
	BaseNode
	Object     ASTNode
	MethodName string
	Args       []ASTNode
	IsPointer  bool // true for ->, false for .
}

func (m *MethodCall) astNode() {}

// ClassRegistry stores defined class types
var ClassRegistry = make(map[string]*ClassDefinition)

// getClassFieldOffset looks up field offset in class
func getClassFieldOffset(className, fieldName string) (int, bool) {
	classDef, exists := ClassRegistry[className]
	if !exists {
		return 0, false
	}

	for _, field := range classDef.Fields {
		if field.Name == fieldName {
			return field.Offset, true
		}
	}

	return 0, false
}
