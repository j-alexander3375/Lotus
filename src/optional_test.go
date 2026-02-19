package main

import (
	"testing"
)

func TestOptionExpressionStructure(t *testing.T) {
	// Test Some(value)
	someExpr := &OptionExpression{
		IsSome: true,
		Value:  &IntLiteral{Value: 42},
	}

	if !someExpr.IsSome {
		t.Error("Expected IsSome to be true")
	}

	if someExpr.Value == nil {
		t.Error("Expected Value to be non-nil for Some")
	}

	// Test None
	noneExpr := &OptionExpression{
		IsSome: false,
		Value:  nil,
	}

	if noneExpr.IsSome {
		t.Error("Expected IsSome to be false for None")
	}

	if noneExpr.Value != nil {
		t.Error("Expected Value to be nil for None")
	}
}

func TestParseOptionExpressionSome(t *testing.T) {
	input := "Some(42)"
	tokens := Tokenize(input)
	parser := NewParser(tokens)

	expr, err := parser.parseExpression()
	if err != nil {
		t.Fatalf("Failed to parse Some expression: %v", err)
	}

	optExpr, ok := expr.(*OptionExpression)
	if !ok {
		t.Fatalf("Expected OptionExpression, got %T", expr)
	}

	if !optExpr.IsSome {
		t.Error("Expected IsSome to be true")
	}

	if optExpr.Value == nil {
		t.Error("Expected Value to be non-nil")
	}
}

func TestParseOptionExpressionNone(t *testing.T) {
	input := "None"
	tokens := Tokenize(input)
	parser := NewParser(tokens)

	expr, err := parser.parseExpression()
	if err != nil {
		t.Fatalf("Failed to parse None expression: %v", err)
	}

	optExpr, ok := expr.(*OptionExpression)
	if !ok {
		t.Fatalf("Expected OptionExpression, got %T", expr)
	}

	if optExpr.IsSome {
		t.Error("Expected IsSome to be false for None")
	}

	if optExpr.Value != nil {
		t.Error("Expected Value to be nil for None")
	}
}

func TestParseOptionWithDifferentTypes(t *testing.T) {
	tests := []struct {
		input    string
		isSome   bool
		hasValue bool
	}{
		{"Some(100)", true, true},
		{"Some(true)", true, true},
		{"Some(3.14)", true, true},
		{`Some("hello")`, true, true},
		{"None", false, false},
	}

	for _, tt := range tests {
		tokens := Tokenize(tt.input)
		parser := NewParser(tokens)

		expr, err := parser.parseExpression()
		if err != nil {
			t.Errorf("Failed to parse %s: %v", tt.input, err)
			continue
		}

		optExpr, ok := expr.(*OptionExpression)
		if !ok {
			t.Errorf("Expected OptionExpression for %s, got %T", tt.input, expr)
			continue
		}

		if optExpr.IsSome != tt.isSome {
			t.Errorf("Input %s: expected IsSome=%v, got %v", tt.input, tt.isSome, optExpr.IsSome)
		}

		if tt.hasValue && optExpr.Value == nil {
			t.Errorf("Input %s: expected non-nil value", tt.input)
		}

		if !tt.hasValue && optExpr.Value != nil {
			t.Errorf("Input %s: expected nil value", tt.input)
		}
	}
}

func TestOptionExpressionInVariableDeclaration(t *testing.T) {
	input := "int x = Some(42);"
	tokens := Tokenize(input)
	parser := NewParser(tokens)

	stmt, err := parser.parseStatement()
	if err != nil {
		t.Fatalf("Failed to parse variable declaration with Some: %v", err)
	}

	varDecl, ok := stmt.(*VariableDeclaration)
	if !ok {
		t.Fatalf("Expected VariableDeclaration, got %T", stmt)
	}

	optExpr, ok := varDecl.Value.(*OptionExpression)
	if !ok {
		t.Fatalf("Expected OptionExpression in initializer, got %T", varDecl.Value)
	}

	if !optExpr.IsSome {
		t.Error("Expected IsSome to be true")
	}
}
