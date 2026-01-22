package main

import (
	"testing"
)

// TestParseSimpleVariableDeclaration tests parsing of variable declarations
func TestParseSimpleVariableDeclaration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		varType  TokenType
		varName  string
		hasValue bool
	}{
		{
			name:     "int variable",
			input:    "int x = 5;",
			varType:  TokenTypeInt,
			varName:  "x",
			hasValue: true,
		},
		{
			name:     "float32 variable",
			input:    "float32 y = 3.14;",
			varType:  TokenTypeFloat32,
			varName:  "y",
			hasValue: true,
		},
		{
			name:     "float64 variable",
			input:    "float64 z = 2.71828;",
			varType:  TokenTypeFloat64,
			varName:  "z",
			hasValue: true,
		},
		{
			name:     "bool variable",
			input:    "bool flag = true;",
			varType:  TokenTypeBool,
			varName:  "flag",
			hasValue: true,
		},
		{
			name:     "string variable",
			input:    `string msg = "hello";`,
			varType:  TokenTypeString,
			varName:  "msg",
			hasValue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := Tokenize(tt.input)
			parser := NewParser(tokens)
			stmt, err := parser.parseStatement()

			if err != nil {
				t.Fatalf("parseStatement() error = %v", err)
			}

			varDecl, ok := stmt.(*VariableDeclaration)
			if !ok {
				t.Fatalf("Expected *VariableDeclaration, got %T", stmt)
			}

			if varDecl.Type != tt.varType {
				t.Errorf("Expected type %v, got %v", tt.varType, varDecl.Type)
			}

			if varDecl.Name != tt.varName {
				t.Errorf("Expected name '%s', got '%s'", tt.varName, varDecl.Name)
			}

			if tt.hasValue && varDecl.Value == nil {
				t.Error("Expected value, got nil")
			}
		})
	}
}

// TestParseConstantDeclaration tests parsing of constant declarations
func TestParseConstantDeclaration(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		constType TokenType
		constName string
	}{
		{
			name:      "int constant",
			input:     "const int MAX = 100;",
			constType: TokenTypeInt,
			constName: "MAX",
		},
		{
			name:      "float32 constant",
			input:     "const float32 PI = 3.14;",
			constType: TokenTypeFloat32,
			constName: "PI",
		},
		{
			name:      "float64 constant",
			input:     "const float64 E = 2.71828;",
			constType: TokenTypeFloat64,
			constName: "E",
		},
		{
			name:      "string constant",
			input:     `const string GREETING = "Hello";`,
			constType: TokenTypeString,
			constName: "GREETING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := Tokenize(tt.input)
			parser := NewParser(tokens)
			stmt, err := parser.parseStatement()

			if err != nil {
				t.Fatalf("parseStatement() error = %v", err)
			}

			constDecl, ok := stmt.(*ConstantDeclaration)
			if !ok {
				t.Fatalf("Expected *ConstantDeclaration, got %T", stmt)
			}

			if constDecl.Type != tt.constType {
				t.Errorf("Expected type %v, got %v", tt.constType, constDecl.Type)
			}

			if constDecl.Name != tt.constName {
				t.Errorf("Expected name '%s', got '%s'", tt.constName, constDecl.Name)
			}

			if constDecl.Value == nil {
				t.Error("Expected value, got nil")
			}
		})
	}
}

// TestParseFunctionDeclaration tests parsing of function declarations
func TestParseFunctionDeclaration(t *testing.T) {
	input := `fn int add(int a, int b) {
		ret a + b;
	}`

	tokens := Tokenize(input)
	parser := NewParser(tokens)
	stmt, err := parser.parseStatement()

	if err != nil {
		t.Fatalf("parseStatement() error = %v", err)
	}

	fnDecl, ok := stmt.(*FunctionDefinition)
	if !ok {
		t.Fatalf("Expected *FunctionDefinition, got %T", stmt)
	}

	if fnDecl.Name != "add" {
		t.Errorf("Expected function name 'add', got '%s'", fnDecl.Name)
	}

	if fnDecl.ReturnType != TokenTypeInt {
		t.Errorf("Expected return type TokenTypeInt, got %v", fnDecl.ReturnType)
	}

	if len(fnDecl.Parameters) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(fnDecl.Parameters))
	}

	if len(fnDecl.Parameters) >= 2 {
		if fnDecl.Parameters[0].Name != "a" {
			t.Errorf("Expected first param name 'a', got '%s'", fnDecl.Parameters[0].Name)
		}
		if fnDecl.Parameters[1].Name != "b" {
			t.Errorf("Expected second param name 'b', got '%s'", fnDecl.Parameters[1].Name)
		}
	}
}

// TestParseIfStatement tests parsing of if statements
func TestParseIfStatement(t *testing.T) {
	input := `if (x > 0) {
		ret 1;
	}`

	tokens := Tokenize(input)
	parser := NewParser(tokens)
	stmt, err := parser.parseStatement()

	if err != nil {
		t.Fatalf("parseStatement() error = %v", err)
	}

	ifStmt, ok := stmt.(*IfStatement)
	if !ok {
		t.Fatalf("Expected *IfStatement, got %T", stmt)
	}

	if ifStmt.Condition == nil {
		t.Error("Expected condition, got nil")
	}

	if ifStmt.ThenBody == nil {
		t.Error("Expected then body, got nil")
	}
}

// TestParseWhileStatement tests parsing of while loops
func TestParseWhileStatement(t *testing.T) {
	input := `while (x < 10) {
		x = x + 1;
	}`

	tokens := Tokenize(input)
	parser := NewParser(tokens)
	stmt, err := parser.parseStatement()

	if err != nil {
		t.Fatalf("parseStatement() error = %v", err)
	}

	whileStmt, ok := stmt.(*WhileLoop)
	if !ok {
		t.Fatalf("Expected *WhileLoop, got %T", stmt)
	}

	if whileStmt.Condition == nil {
		t.Error("Expected condition, got nil")
	}

	if whileStmt.Body == nil {
		t.Error("Expected body, got nil")
	}
}

// TestParseForStatement tests parsing of for loops
func TestParseForStatement(t *testing.T) {
	input := `for (int i = 0; i < 10; i = i + 1) {
		println("test");
	}`

	tokens := Tokenize(input)
	parser := NewParser(tokens)
	stmt, err := parser.parseStatement()

	if err != nil {
		t.Fatalf("parseStatement() error = %v", err)
	}

	forStmt, ok := stmt.(*ForLoop)
	if !ok {
		t.Fatalf("Expected *ForLoop, got %T", stmt)
	}

	if forStmt.Init == nil {
		t.Error("Expected init, got nil")
	}

	if forStmt.Condition == nil {
		t.Error("Expected condition, got nil")
	}

	if forStmt.Update == nil {
		t.Error("Expected update, got nil")
	}

	if forStmt.Body == nil {
		t.Error("Expected body, got nil")
	}
}

// TestParseReturnStatement tests parsing of return statements
func TestParseReturnStatement(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasValue bool
	}{
		{
			name:     "return with value",
			input:    "ret 42;",
			hasValue: true,
		},
		{
			name:     "return keyword",
			input:    "return 100;",
			hasValue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := Tokenize(tt.input)
			parser := NewParser(tokens)
			stmt, err := parser.parseStatement()

			if err != nil {
				t.Fatalf("parseStatement() error = %v", err)
			}

			retStmt, ok := stmt.(*ReturnStatement)
			if !ok {
				t.Fatalf("Expected *ReturnStatement, got %T", stmt)
			}

			if tt.hasValue && retStmt.Value == nil {
				t.Error("Expected return value, got nil")
			}
		})
	}
}

// TestParseBinaryExpression tests parsing of binary expressions
func TestParseBinaryExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		operator TokenType
	}{
		{"addition", "a + b", TokenPlus},
		{"subtraction", "a - b", TokenMinus},
		{"multiplication", "a * b", TokenStar},
		{"division", "a / b", TokenSlash},
		{"modulo", "a % b", TokenPercent},
		{"equality", "a == b", TokenEqual},
		{"inequality", "a != b", TokenNotEqual},
		{"less than", "a < b", TokenLess},
		{"less or equal", "a <= b", TokenLessEq},
		{"greater than", "a > b", TokenGreater},
		{"greater or equal", "a >= b", TokenGreaterEq},
		{"logical and", "a && b", TokenAnd},
		{"logical or", "a || b", TokenOr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Wrap in a statement context
			input := "int x = " + tt.input + ";"
			tokens := Tokenize(input)
			parser := NewParser(tokens)
			stmt, err := parser.parseStatement()

			if err != nil {
				t.Fatalf("parseStatement() error = %v", err)
			}

			varDecl, ok := stmt.(*VariableDeclaration)
			if !ok {
				t.Fatalf("Expected *VariableDeclaration, got %T", stmt)
			}

			// Check if the value is a binary operation or comparison
			switch expr := varDecl.Value.(type) {
			case *BinaryOp:
				if expr.Operator != tt.operator {
					t.Errorf("Expected operator %v, got %v", tt.operator, expr.Operator)
				}
			case *Comparison:
				if expr.Operator != tt.operator {
					t.Errorf("Expected operator %v, got %v", tt.operator, expr.Operator)
				}
			case *LogicalOp:
				if expr.Operator != tt.operator {
					t.Errorf("Expected operator %v, got %v", tt.operator, expr.Operator)
				}
			default:
				t.Errorf("Expected binary/comparison/logical operation, got %T", expr)
			}
		})
	}
}

// TestParseStructDefinition tests parsing of struct definitions
func TestParseStructDefinition(t *testing.T) {
	t.Skip("Struct parsing not yet implemented in parseStatement - requires adding TokenStruct case")
}

// TestParseTypeName tests the parseTypeName function
func TestParseTypeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"int", "int"},
		{"int8", "int8"},
		{"int16", "int16"},
		{"int32", "int32"},
		{"int64", "int64"},
		{"uint", "uint"},
		{"uint8", "uint8"},
		{"uint16", "uint16"},
		{"uint32", "uint32"},
		{"uint64", "uint64"},
		{"float32", "float32"},
		{"float64", "float64"},
		{"string", "string"},
		{"char", "char"},
		{"bool", "bool"},
		{"void", "void"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := Tokenize(tt.input)
			parser := NewParser(tokens)
			typeName, err := parser.parseTypeName()

			if err != nil {
				t.Fatalf("parseTypeName() error = %v", err)
			}

			if typeName != tt.expected {
				t.Errorf("Expected type name '%s', got '%s'", tt.expected, typeName)
			}
		})
	}
}

// TestIsTypeToken tests the isTypeToken helper function
func TestIsTypeToken(t *testing.T) {
	typeTokens := []TokenType{
		TokenTypeInt, TokenTypeInt8, TokenTypeInt16, TokenTypeInt32, TokenTypeInt64,
		TokenTypeUint, TokenTypeUint8, TokenTypeUint16, TokenTypeUint32, TokenTypeUint64,
		TokenTypeFloat32, TokenTypeFloat64,
		TokenTypeString, TokenTypeBool, TokenTypeChar, TokenTypeVoid,
	}

	nonTypeTokens := []TokenType{
		TokenPlus, TokenMinus, TokenStar, TokenSlash,
		TokenIdentifier, TokenInt, TokenFloat,
		TokenLParen, TokenRParen, TokenLBrace, TokenRBrace,
	}

	for _, tt := range typeTokens {
		t.Run("TypeToken_"+TokenTypeName(tt), func(t *testing.T) {
			if !isTypeToken(tt) {
				t.Errorf("isTypeToken(%v) = false, expected true", tt)
			}
		})
	}

	for _, tt := range nonTypeTokens {
		t.Run("NonTypeToken_"+TokenTypeName(tt), func(t *testing.T) {
			if isTypeToken(tt) {
				t.Errorf("isTypeToken(%v) = true, expected false", tt)
			}
		})
	}
}

// TestParseArrayLiteral tests parsing of array literals
func TestParseArrayLiteral(t *testing.T) {
	input := "int arr = [1, 2, 3, 4, 5];"

	tokens := Tokenize(input)
	parser := NewParser(tokens)
	stmt, err := parser.parseStatement()

	if err != nil {
		t.Fatalf("parseStatement() error = %v", err)
	}

	varDecl, ok := stmt.(*VariableDeclaration)
	if !ok {
		t.Fatalf("Expected *VariableDeclaration, got %T", stmt)
	}

	arrLit, ok := varDecl.Value.(*ArrayLiteral)
	if !ok {
		t.Fatalf("Expected *ArrayLiteral, got %T", varDecl.Value)
	}

	if len(arrLit.Elements) != 5 {
		t.Errorf("Expected 5 elements, got %d", len(arrLit.Elements))
	}
}

// TestParseEnumDefinition tests parsing of enum definitions
func TestParseEnumDefinition(t *testing.T) {
	t.Skip("Enum parsing not yet implemented in parseStatement - requires adding TokenEnum case")
}
