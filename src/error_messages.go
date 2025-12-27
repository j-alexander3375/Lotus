package main

import (
	"fmt"
	"strings"
)

// error_messages.go - Human-friendly error messages for the Lotus compiler
// Provides detailed error messages with context and suggestions.

// ErrorCode represents a unique error identifier
type ErrorCode string

// Error codes for different error categories
const (
	// Syntax errors (E01xx)
	ErrUnexpectedToken     ErrorCode = "E0101"
	ErrExpectedToken       ErrorCode = "E0102"
	ErrMissingBrace        ErrorCode = "E0103"
	ErrMissingSemicolon    ErrorCode = "E0104"
	ErrInvalidExpression   ErrorCode = "E0105"
	ErrMissingFunctionBody ErrorCode = "E0106"
	ErrInvalidDeclaration  ErrorCode = "E0107"

	// Semantic errors (E02xx)
	ErrUndefinedVariable ErrorCode = "E0201"
	ErrUndefinedFunction ErrorCode = "E0202"
	ErrRedefinition      ErrorCode = "E0203"
	ErrTypeMismatch      ErrorCode = "E0204"
	ErrInvalidOperation  ErrorCode = "E0205"

	// Type errors (E03xx)
	ErrIncompatibleTypes ErrorCode = "E0301"
	ErrInvalidCast       ErrorCode = "E0302"
	ErrArrayIndexType    ErrorCode = "E0303"

	// Import errors (E04xx)
	ErrModuleNotFound      ErrorCode = "E0401"
	ErrFunctionNotExported ErrorCode = "E0402"
	ErrCircularImport      ErrorCode = "E0403"
)

// TokenTypeName returns a human-readable name for a token type
func TokenTypeName(t TokenType) string {
	names := map[TokenType]string{
		TokenEOF:        "end of file",
		TokenNewline:    "newline",
		TokenIdentifier: "identifier",
		TokenInt:        "number",
		TokenString:     "string",
		TokenChar:       "character",
		TokenBool:       "boolean",
		TokenFloat:      "float",
		TokenSemi:       "';'",
		TokenComma:      "','",
		TokenColon:      "':'",
		TokenLParen:     "'('",
		TokenRParen:     "')'",
		TokenLBrace:     "'{'",
		TokenRBrace:     "'}'",
		TokenLBracket:   "'['",
		TokenRBracket:   "']'",
		TokenAssign:     "'='",
		TokenPlus:       "'+'",
		TokenMinus:      "'-'",
		TokenStar:       "'*'",
		TokenSlash:      "'/'",
		TokenPercent:    "'%'",
		TokenEqual:      "'=='",
		TokenNotEqual:   "'!='",
		TokenLess:       "'<'",
		TokenGreater:    "'>'",
		TokenLessEq:     "'<='",
		TokenGreaterEq:  "'>='",
		TokenAnd:        "'&&'",
		TokenOr:         "'||'",
		TokenExclaim:    "'!'",
		TokenAmpersand:  "'&'",
		TokenPipe:       "'|'",
		TokenCaret:      "'^'",
		TokenTilde:      "'~'",
		TokenLShift:     "'<<'",
		TokenRShift:     "'>>'",
		TokenFn:         "'fn'",
		TokenRet:        "'ret'",
		TokenReturn:     "'return'",
		TokenIf:         "'if'",
		TokenElse:       "'else'",
		TokenWhile:      "'while'",
		TokenFor:        "'for'",
		TokenConst:      "'const'",
		TokenUse:        "'use'",
		TokenAs:         "'as'",
		TokenStruct:     "'struct'",
		TokenEnum:       "'enum'",
		TokenClass:      "'class'",
		TokenNew:        "'new'",
		TokenTry:        "'try'",
		TokenCatch:      "'catch'",
		TokenNull:       "'null'",
		TokenTypeInt:    "'int'",
		TokenTypeInt8:   "'int8'",
		TokenTypeInt16:  "'int16'",
		TokenTypeInt32:  "'int32'",
		TokenTypeInt64:  "'int64'",
		TokenTypeUint:   "'uint'",
		TokenTypeUint8:  "'uint8'",
		TokenTypeUint16: "'uint16'",
		TokenTypeUint32: "'uint32'",
		TokenTypeUint64: "'uint64'",
		TokenTypeFloat:  "'float'",
		TokenTypeString: "'string'",
		TokenTypeBool:   "'bool'",
		TokenTypeChar:   "'char'",
	}

	if name, ok := names[t]; ok {
		return name
	}
	return fmt.Sprintf("token(%d)", t)
}

// ParseError represents a detailed parse error with context
type ParseError struct {
	Code       ErrorCode
	Message    string
	FilePath   string
	Line       int
	Column     int
	Context    string // The source line
	Suggestion string
	Notes      []string
}

func (e *ParseError) Error() string {
	return e.Message
}

// NewParseError creates a detailed parse error
func NewParseError(code ErrorCode, msg string, line, col int) *ParseError {
	return &ParseError{
		Code:    code,
		Message: msg,
		Line:    line,
		Column:  col,
	}
}

// WithContext adds source context to the error
func (e *ParseError) WithContext(context string) *ParseError {
	e.Context = context
	return e
}

// WithSuggestion adds a suggestion to the error
func (e *ParseError) WithSuggestion(suggestion string) *ParseError {
	e.Suggestion = suggestion
	return e
}

// WithNote adds a note to the error
func (e *ParseError) WithNote(note string) *ParseError {
	e.Notes = append(e.Notes, note)
	return e
}

// FormatExpectedToken creates a helpful "expected X, got Y" message
func FormatExpectedToken(expected, got TokenType, gotValue string) string {
	msg := fmt.Sprintf("expected %s, got %s", TokenTypeName(expected), TokenTypeName(got))
	if gotValue != "" && got == TokenIdentifier {
		msg += fmt.Sprintf(" '%s'", gotValue)
	}
	return msg
}

// FormatUnexpectedToken creates a message for unexpected tokens
func FormatUnexpectedToken(got TokenType, gotValue string, context string) string {
	msg := fmt.Sprintf("unexpected %s", TokenTypeName(got))
	if gotValue != "" && got == TokenIdentifier {
		msg += fmt.Sprintf(" '%s'", gotValue)
	}
	if context != "" {
		msg += fmt.Sprintf(" %s", context)
	}
	return msg
}

// SuggestForMissingBrace suggests fixes for brace-related errors
func SuggestForMissingBrace(openOrClose string) string {
	if openOrClose == "open" {
		return "add '{' to open the block"
	}
	return "add '}' to close the block"
}

// SuggestForTypo attempts to find a similar keyword for typos
func SuggestForTypo(typo string) string {
	keywords := []string{
		"fn", "ret", "return", "if", "elif", "else", "while", "for",
		"break", "continue", "use", "const", "true", "false", "nil",
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float", "float32", "float64", "string", "bool", "void",
		"struct", "enum", "class", "impl", "pub", "priv",
		"println", "printf", "print", "malloc", "free", "sizeof",
	}

	typoLower := strings.ToLower(typo)
	bestMatch := ""
	bestDist := 3 // Max distance for suggestions

	for _, kw := range keywords {
		dist := levenshteinDistance(typoLower, kw)
		if dist < bestDist {
			bestDist = dist
			bestMatch = kw
		}
	}

	if bestMatch != "" {
		return fmt.Sprintf("did you mean '%s'?", bestMatch)
	}
	return ""
}

// SuggestForMissingSemicolon provides context-aware semicolon suggestions
func SuggestForMissingSemicolon(context string) string {
	return "add ';' at the end of the statement"
}

// ErrorHelpText provides extended help for error codes
func ErrorHelpText(code ErrorCode) string {
	help := map[ErrorCode]string{
		ErrUnexpectedToken: `An unexpected token was encountered during parsing.
This usually means there's a syntax error in your code.
Check for:
  - Missing semicolons at the end of statements
  - Mismatched braces or parentheses
  - Typos in keywords`,

		ErrExpectedToken: `The parser expected a specific token but found something else.
Common causes:
  - Missing opening or closing brace
  - Missing parentheses in function calls
  - Forgetting the return type in function declarations`,

		ErrMissingBrace: `A brace '{' or '}' is missing.
Every opening brace must have a matching closing brace.
Check your function bodies, if statements, and loops.`,

		ErrUndefinedVariable: `A variable was used before it was declared.
In Lotus, variables must be declared before use:
  int x = 42;  // Declaration
  x = 10;      // Usage`,

		ErrUndefinedFunction: `A function was called that doesn't exist.
Check that:
  - The function is defined before it's called
  - The module containing the function is imported
  - The function name is spelled correctly`,

		ErrTypeMismatch: `The types in an expression don't match.
Lotus is statically typed - ensure both sides of an
operation have compatible types.`,

		ErrModuleNotFound: `The imported module could not be found.
Check that:
  - The module name is spelled correctly
  - The module is part of the standard library
  - Custom modules are in the include path`,
	}

	if text, ok := help[code]; ok {
		return text
	}
	return ""
}

// Common error message templates
var (
	MsgMissingReturnType    = "expected return type after 'fn'"
	MsgMissingFunctionName  = "expected function name after return type"
	MsgMissingParameterList = "expected '(' after function name"
	MsgMissingFunctionBody  = "expected '{' to start function body"
	MsgMissingCondition     = "expected condition expression"
	MsgMissingBlockOpen     = "expected '{' to open block"
	MsgMissingBlockClose    = "expected '}' to close block"
	MsgMissingExpression    = "expected expression"
	MsgMissingIdentifier    = "expected identifier"
	MsgMissingType          = "expected type"
	MsgMissingSemicolon     = "expected ';' after statement"
	MsgMissingAssignment    = "expected '=' in declaration"
	MsgInvalidEscapeSeq     = "invalid escape sequence"
	MsgUnterminatedString   = "unterminated string literal"
	MsgInvalidNumber        = "invalid number literal"
)
