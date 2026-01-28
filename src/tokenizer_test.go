package main

import (
	"testing"
)

func TestTokenizeBasicTokens(t *testing.T) {
	input := "fn int main() { ret 0; }"
	tokens := Tokenize(input)
	var err error = nil

	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	expectedTypes := []TokenType{
		TokenFn,
		TokenTypeInt,
		TokenIdentifier,
		TokenLParen,
		TokenRParen,
		TokenLBrace,
		TokenRet,
		TokenInt,
		TokenSemi,
		TokenRBrace,
		TokenEOF,
	}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, expected := range expectedTypes {
		if tokens[i].Type != expected {
			t.Errorf("Token %d: expected type %v, got %v", i, expected, tokens[i].Type)
		}
	}
}

func TestTokenizeTemplateKeywords(t *testing.T) {
	input := "template<typename T> namespace using"
	tokens := Tokenize(input)
	var err error = nil

	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	if len(tokens) < 7 {
		t.Fatalf("Expected at least 7 tokens, got %d", len(tokens))
	}

	if tokens[0].Type != TokenTemplate {
		t.Errorf("Expected TokenTemplate, got %v", tokens[0].Type)
	}

	if tokens[2].Type != TokenTypename {
		t.Errorf("Expected TokenTypename, got %v", tokens[2].Type)
	}

	if tokens[5].Type != TokenNamespace {
		t.Errorf("Expected TokenNamespace, got %v", tokens[5].Type)
	}

	if tokens[6].Type != TokenUsing {
		t.Errorf("Expected TokenUsing, got %v", tokens[6].Type)
	}
}

func TestTokenizeFloatKeyword(t *testing.T) {
	input := "float x = 3.14;"
	tokens := Tokenize(input)
	var err error = nil

	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	if tokens[0].Type != TokenTypeFloat32 {
		t.Errorf("Expected TokenTypeFloat32 for 'float', got %v", tokens[0].Type)
	}
}

func TestTokenizeNumbers(t *testing.T) {
	input := "42 3.14 0"
	tokens := Tokenize(input)
	var err error = nil

	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	if tokens[0].Type != TokenInt {
		t.Errorf("Expected TokenInt for '42', got %v", tokens[0].Type)
	}

	if tokens[1].Type != TokenFloat {
		t.Errorf("Expected TokenFloat for '3.14', got %v", tokens[1].Type)
	}
}

func TestTokenizeStrings(t *testing.T) {
	input := `"hello world"`
	tokens := Tokenize(input)
	var err error = nil

	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	if tokens[0].Type != TokenString {
		t.Errorf("Expected TokenString, got %v", tokens[0].Type)
	}

	if tokens[0].Value != "hello world" {
		t.Errorf("Expected value 'hello world', got '%s'", tokens[0].Value)
	}
}

func TestTokenizeOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"+", TokenPlus},
		{"-", TokenMinus},
		{"*", TokenStar},
		{"/", TokenSlash},
		{"==", TokenEqual},
		{"!=", TokenNotEqual},
		{"<", TokenLess},
		{"<=", TokenLessEq},
		{">", TokenGreater},
		{">=", TokenGreaterEq},
		{"::", TokenColon}, // First colon
	}

	for _, tt := range tests {
		tokens := Tokenize(tt.input)
		var err error = nil
		if err != nil {
			t.Errorf("Tokenize '%s' failed: %v", tt.input, err)
			continue
		}
		if len(tokens) < 1 {
			t.Errorf("No tokens for input '%s'", tt.input)
			continue
		}
		// Note: TokenColon is tricky - :: produces two TokenColon tokens
		if tokens[0].Type != tt.expected {
			t.Errorf("Input '%s': expected %v, got %v", tt.input, tt.expected, tokens[0].Type)
		}
	}
}
