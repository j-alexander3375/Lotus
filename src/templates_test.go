package main

import (
	"testing"
)

func TestTemplateDeclarationStructure(t *testing.T) {
	// Test that template declaration structure is correct
	params := []TemplateParameter{
		{Name: "T", IsType: true},
	}

	fnDef := &FunctionDefinition{
		Name:       "testFunc",
		ReturnType: TokenIdentifier,
		Parameters: []FunctionParam{
			{Name: "a", Type: TokenIdentifier},
		},
	}

	template := &TemplateDeclaration{
		Parameters:  params,
		Declaration: fnDef,
	}

	if len(template.Parameters) != 1 {
		t.Errorf("Expected 1 template parameter, got %d", len(template.Parameters))
	}

	if template.Parameters[0].Name != "T" {
		t.Errorf("Expected parameter name 'T', got '%s'", template.Parameters[0].Name)
	}

	if !template.Parameters[0].IsType {
		t.Error("Expected IsType to be true")
	}
}

func TestTemplateParameterTypes(t *testing.T) {
	typeParam := TemplateParameter{
		Name:   "T",
		IsType: true,
	}

	valueParam := TemplateParameter{
		Name:   "N",
		IsType: false,
		Type:   TokenTypeInt32,
	}

	if !typeParam.IsType {
		t.Error("Type parameter should have IsType true")
	}

	if valueParam.IsType {
		t.Error("Value parameter should have IsType false")
	}

	if valueParam.Type != TokenTypeInt32 {
		t.Errorf("Expected value type Int32, got %v", valueParam.Type)
	}
}

func TestNamespaceDeclarationStructure(t *testing.T) {
	ns := &NamespaceDeclaration{
		Name: "test",
		Body: []ASTNode{
			&FunctionDefinition{Name: "func1"},
			&FunctionDefinition{Name: "func2"},
		},
	}

	if ns.Name != "test" {
		t.Errorf("Expected namespace name 'test', got '%s'", ns.Name)
	}

	if len(ns.Body) != 2 {
		t.Errorf("Expected 2 body items, got %d", len(ns.Body))
	}
}

func TestNamespaceStructure(t *testing.T) {
	ns := &Namespace{
		Name:      "math",
		Functions: make(map[string][]*FunctionDefinition),
		Variables: make(map[string]*VariableDeclaration),
		Constants: make(map[string]*ConstantDeclaration),
		Children:  make(map[string]*Namespace),
	}

	// Add a function
	fn := &FunctionDefinition{Name: "add"}
	ns.Functions["add"] = append(ns.Functions["add"], fn)

	if len(ns.Functions["add"]) != 1 {
		t.Errorf("Expected 1 function, got %d", len(ns.Functions["add"]))
	}

	// Test child namespace
	child := &Namespace{Name: "utils"}
	ns.Children["utils"] = child

	if ns.Children["utils"].Name != "utils" {
		t.Error("Child namespace not added correctly")
	}
}

func TestUsingDeclaration(t *testing.T) {
	using := &UsingDeclaration{
		IsNamespace: false,
		Namespaces:  []string{"math"},
		Name:        "sqrt",
	}

	if using.Namespaces[0] != "math" {
		t.Errorf("Expected namespace 'math', got '%s'", using.Namespaces[0])
	}

	if using.Name != "sqrt" {
		t.Errorf("Expected name 'sqrt', got '%s'", using.Name)
	}

	// Test namespace using
	nsUsing := &UsingDeclaration{
		IsNamespace: true,
		Namespaces:  []string{"std", "collections"},
	}

	if !nsUsing.IsNamespace {
		t.Error("Expected IsNamespace to be true")
	}

	if len(nsUsing.Namespaces) != 2 {
		t.Errorf("Expected 2 namespaces, got %d", len(nsUsing.Namespaces))
	}
}
