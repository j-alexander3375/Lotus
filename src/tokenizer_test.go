package main

import (
	"testing"
)

// TestTokenizeKeywords tests tokenization of language keywords
func TestTokenizeKeywords(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"ret", TokenRet},
		{"return", TokenReturn},
		{"const", TokenConst},
		{"if", TokenIf},
		{"else", TokenElse},
		{"while", TokenWhile},
		{"for", TokenFor},
		{"break", TokenBreak},
		{"continue", TokenContinue},
		{"fn", TokenFn},
		{"inline", TokenInline},
		{"struct", TokenStruct},
		{"enum", TokenEnum},
		{"class", TokenClass},
		{"new", TokenNew},
		{"vrt", TokenVirtual},
		{"override", TokenOverride},
		{"static", TokenStatic},
		{"lcl", TokenLocal},
		{"gbl", TokenGlobal},
		{"partial", TokenPartial},
		{"wrap", TokenWrap},
		{"use", TokenUse},
		{"as", TokenAs},
		{"try", TokenTry},
		{"catch", TokenCatch},
		{"finally", TokenFinally},
		{"throw", TokenThrow},
		{"null", TokenNull},
		{"bitcast", TokenBitcast},
		{"transmute", TokenTransmute},
		{"match", TokenMatch},
		{"case", TokenCase},
		{"default", TokenDefault},
		{"when", TokenWhen},
		{"to", TokenTo},
		{"union", TokenUnion},
		{"option", TokenOption},
		{"Some", TokenSome},
		{"None", TokenNone},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := Tokenize(tt.input)
			if len(tokens) < 1 {
				t.Fatalf("Expected at least 1 token, got %d", len(tokens))
			}
			if tokens[0].Type != tt.expected {
				t.Errorf("Expected token type %v, got %v", tt.expected, tokens[0].Type)
			}
		})
	}
}

// TestTokenizeTypeKeywords tests tokenization of type keywords
func TestTokenizeTypeKeywords(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"int", TokenTypeInt},
		{"int8", TokenTypeInt8},
		{"int16", TokenTypeInt16},
		{"int32", TokenTypeInt32},
		{"int64", TokenTypeInt64},
		{"uint", TokenTypeUint},
		{"uint8", TokenTypeUint8},
		{"uint16", TokenTypeUint16},
		{"uint32", TokenTypeUint32},
		{"uint64", TokenTypeUint64},
		{"float32", TokenTypeFloat32},
		{"float64", TokenTypeFloat64},
		{"string", TokenTypeString},
		{"char", TokenTypeChar},
		{"bool", TokenTypeBool},
		{"void", TokenTypeVoid},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := Tokenize(tt.input)
			if len(tokens) < 1 {
				t.Fatalf("Expected at least 1 token, got %d", len(tokens))
			}
			if tokens[0].Type != tt.expected {
				t.Errorf("Expected token type %v for '%s', got %v", tt.expected, tt.input, tokens[0].Type)
			}
		})
	}
}

// TestTokenizeOperators tests tokenization of operators
func TestTokenizeOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected []TokenType
	}{
		{"+", []TokenType{TokenPlus}},
		{"-", []TokenType{TokenMinus}},
		{"*", []TokenType{TokenStar}},
		{"/", []TokenType{TokenSlash}},
		{"%", []TokenType{TokenPercent}},
		{"==", []TokenType{TokenEqual}},
		{"!=", []TokenType{TokenNotEqual}},
		{"<", []TokenType{TokenLess}},
		{"<=", []TokenType{TokenLessEq}},
		{">", []TokenType{TokenGreater}},
		{">=", []TokenType{TokenGreaterEq}},
		{"&&", []TokenType{TokenAnd}},
		{"||", []TokenType{TokenOr}},
		{"&", []TokenType{TokenAmpersand}},
		{"|", []TokenType{TokenPipe}},
		{"^", []TokenType{TokenCaret}},
		{"~", []TokenType{TokenTilde}},
		{"!", []TokenType{TokenExclaim}},
		{"<<", []TokenType{TokenLShift}},
		{">>", []TokenType{TokenRShift}},
		{"++", []TokenType{TokenPlusPlus}},
		{"--", []TokenType{TokenMinusMinus}},
		{"+=", []TokenType{TokenPlusEq}},
		{"-=", []TokenType{TokenMinusEq}},
		{"*=", []TokenType{TokenStarEq}},
		{"/=", []TokenType{TokenSlashEq}},
		{"%=", []TokenType{TokenPercentEq}},
		{"->", []TokenType{TokenArrow}},
		{"=>", []TokenType{TokenLambda}},
		{"|>", []TokenType{TokenPipe2}},
		{"..", []TokenType{TokenRange}},
		{"@", []TokenType{TokenAt}},
		{"?", []TokenType{TokenQuestion}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := Tokenize(tt.input)
			if len(tokens) < len(tt.expected) {
				t.Fatalf("Expected at least %d tokens, got %d", len(tt.expected), len(tokens))
			}
			for i, expectedType := range tt.expected {
				if tokens[i].Type != expectedType {
					t.Errorf("Token %d: expected %v, got %v", i, expectedType, tokens[i].Type)
				}
			}
		})
	}
}

// TestTokenizeIntegerLiterals tests tokenization of integer literals
func TestTokenizeIntegerLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0", "0"},
		{"123", "123"},
		{"42", "42"},
		{"0x1A", "0x1A"},
		{"0xFF", "0xFF"},
		{"0o777", "0o777"},
		{"0b1010", "0b1010"},
		{"0b11111111", "0b11111111"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := Tokenize(tt.input)
			if len(tokens) < 1 {
				t.Fatalf("Expected at least 1 token, got %d", len(tokens))
			}
			if tokens[0].Type != TokenInt {
				t.Errorf("Expected TokenInt, got %v", tokens[0].Type)
			}
			if tokens[0].Value != tt.expected {
				t.Errorf("Expected value '%s', got '%s'", tt.expected, tokens[0].Value)
			}
		})
	}
}

// TestTokenizeFloatLiterals tests tokenization of float literals
func TestTokenizeFloatLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0.0", "0.0"},
		{"1.5", "1.5"},
		{"3.14159", "3.14159"},
		{"0.5", "0.5"},
		{"123.456", "123.456"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := Tokenize(tt.input)
			if len(tokens) < 1 {
				t.Fatalf("Expected at least 1 token, got %d", len(tokens))
			}
			if tokens[0].Type != TokenFloat {
				t.Errorf("Expected TokenFloat, got %v", tokens[0].Type)
			}
			if tokens[0].Value != tt.expected {
				t.Errorf("Expected value '%s', got '%s'", tt.expected, tokens[0].Value)
			}
		})
	}
}

// TestTokenizeStringLiterals tests tokenization of string literals
func TestTokenizeStringLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`"world"`, "world"},
		{`"hello\nworld"`, "hello\nworld"},
		{`"tab\there"`, "tab\there"},
		{`"quote\"test"`, "quote\"test"},
		{`"backslash\\"`, "backslash\\"},
		{`""`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := Tokenize(tt.input)
			if len(tokens) < 1 {
				t.Fatalf("Expected at least 1 token, got %d", len(tokens))
			}
			if tokens[0].Type != TokenString {
				t.Errorf("Expected TokenString, got %v", tokens[0].Type)
			}
			if tokens[0].Value != tt.expected {
				t.Errorf("Expected value '%s', got '%s'", tt.expected, tokens[0].Value)
			}
		})
	}
}

// TestTokenizeCharLiterals tests tokenization of character literals
func TestTokenizeCharLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"'a'", "a"},
		{"'Z'", "Z"},
		{"'\\n'", "\n"},
		{"'\\t'", "\t"},
		{"'\\\\'", "\\"},
		{"'\\''", "'"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := Tokenize(tt.input)
			if len(tokens) < 1 {
				t.Fatalf("Expected at least 1 token, got %d", len(tokens))
			}
			if tokens[0].Type != TokenChar {
				t.Errorf("Expected TokenChar, got %v", tokens[0].Type)
			}
			if tokens[0].Value != tt.expected {
				t.Errorf("Expected value '%s', got '%s'", tt.expected, tokens[0].Value)
			}
		})
	}
}

// TestTokenizeComments tests that comments are properly skipped
func TestTokenizeComments(t *testing.T) {
	input := `int x = 5; // this is a comment
int y = 10;`
	tokens := Tokenize(input)

	// Should have: int, x, =, 5, ;, newline, int, y, =, 10, ;, EOF
	expectedTypes := []TokenType{
		TokenTypeInt, TokenIdentifier, TokenAssign, TokenInt, TokenSemi,
		TokenNewline,
		TokenTypeInt, TokenIdentifier, TokenAssign, TokenInt, TokenSemi,
		TokenEOF,
	}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, expectedType := range expectedTypes {
		if tokens[i].Type != expectedType {
			t.Errorf("Token %d: expected %v, got %v", i, expectedType, tokens[i].Type)
		}
	}
}

// TestTokenizeInterpolatedStrings tests interpolated string literals
func TestTokenizeInterpolatedStrings(t *testing.T) {
	input := `$"Hello {name}"`
	tokens := Tokenize(input)

	if len(tokens) < 1 {
		t.Fatalf("Expected at least 1 token, got %d", len(tokens))
	}
	if tokens[0].Type != TokenInterpolatedString {
		t.Errorf("Expected TokenInterpolatedString, got %v", tokens[0].Type)
	}
	if tokens[0].Value != "Hello {name}" {
		t.Errorf("Expected 'Hello {name}', got '%s'", tokens[0].Value)
	}
}

// TestTokenizeComplexExpression tests a complex expression
func TestTokenizeComplexExpression(t *testing.T) {
	input := `fn int add(int a, int b) { ret a + b; }`
	tokens := Tokenize(input)

	expectedTypes := []TokenType{
		TokenFn, TokenTypeInt, TokenIdentifier, TokenLParen,
		TokenTypeInt, TokenIdentifier, TokenComma, TokenTypeInt, TokenIdentifier,
		TokenRParen, TokenLBrace, TokenRet, TokenIdentifier, TokenPlus,
		TokenIdentifier, TokenSemi, TokenRBrace, TokenEOF,
	}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, expectedType := range expectedTypes {
		if tokens[i].Type != expectedType {
			t.Errorf("Token %d: expected %v, got %v", i, expectedType, tokens[i].Type)
		}
	}
}

// TestTokenizeBoolLiterals tests tokenization of boolean literals
func TestTokenizeBoolLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"true", "true"},
		{"false", "false"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := Tokenize(tt.input)
			if len(tokens) < 1 {
				t.Fatalf("Expected at least 1 token, got %d", len(tokens))
			}
			if tokens[0].Type != TokenBool {
				t.Errorf("Expected TokenBool, got %v", tokens[0].Type)
			}
			if tokens[0].Value != tt.expected {
				t.Errorf("Expected value '%s', got '%s'", tt.expected, tokens[0].Value)
			}
		})
	}
}

// TestTokenizeDelimiters tests tokenization of delimiters
func TestTokenizeDelimiters(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"(", TokenLParen},
		{")", TokenRParen},
		{"{", TokenLBrace},
		{"}", TokenRBrace},
		{"[", TokenLBracket},
		{"]", TokenRBracket},
		{",", TokenComma},
		{";", TokenSemi},
		{".", TokenDot},
		{":", TokenColon},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := Tokenize(tt.input)
			if len(tokens) < 1 {
				t.Fatalf("Expected at least 1 token, got %d", len(tokens))
			}
			if tokens[0].Type != tt.expected {
				t.Errorf("Expected token type %v, got %v", tt.expected, tokens[0].Type)
			}
		})
	}
}

// TestTokenizeLineAndColumn tests that line and column tracking works
func TestTokenizeLineAndColumn(t *testing.T) {
	input := `int x
int y`
	tokens := Tokenize(input)

	// First 'int' should be at line 1, col 1
	if tokens[0].Line != 1 || tokens[0].Column != 1 {
		t.Errorf("First token: expected line 1 col 1, got line %d col %d", tokens[0].Line, tokens[0].Column)
	}

	// Second 'int' should be at line 2, col 1
	// tokens[0] = int, tokens[1] = x, tokens[2] = newline, tokens[3] = int
	if len(tokens) > 3 {
		if tokens[3].Line != 2 || tokens[3].Column != 1 {
			t.Errorf("Second int token: expected line 2 col 1, got line %d col %d", tokens[3].Line, tokens[3].Column)
		}
	}
}

// TestTokenizeRangeOperator tests that range operator is not confused with float
func TestTokenizeRangeOperator(t *testing.T) {
	input := `1..10`
	tokens := Tokenize(input)

	expectedTypes := []TokenType{TokenInt, TokenRange, TokenInt, TokenEOF}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, expectedType := range expectedTypes {
		if tokens[i].Type != expectedType {
			t.Errorf("Token %d: expected %v, got %v", i, expectedType, tokens[i].Type)
		}
	}
}

// TestTokenizeIdentifiers tests tokenization of identifiers
func TestTokenizeIdentifiers(t *testing.T) {
	tests := []string{
		"x",
		"variable",
		"my_var",
		"var123",
		"_private",
		"camelCase",
		"snake_case",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			tokens := Tokenize(tt)
			if len(tokens) < 1 {
				t.Fatalf("Expected at least 1 token, got %d", len(tokens))
			}
			if tokens[0].Type != TokenIdentifier {
				t.Errorf("Expected TokenIdentifier, got %v", tokens[0].Type)
			}
			if tokens[0].Value != tt {
				t.Errorf("Expected value '%s', got '%s'", tt, tokens[0].Value)
			}
		})
	}
}
