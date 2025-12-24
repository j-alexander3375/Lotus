package main

// keywords.go - Token type definitions and lexical token structure
// This file defines all token types used by the lexer and parser.

// TokenType represents the type of a lexical token
type TokenType int

// Token type constants - organized by category
const (
	// Keywords and control flow
	TokenRet    TokenType = iota
	TokenReturn           // return
	TokenConst            // const
	TokenIf               // if
	TokenElse             // else
	TokenWhile            // while
	TokenFor              // for
	TokenFn               // fn
	TokenUse              // use (imports)
	TokenAs               // as (aliasing)

	// Literals
	TokenInt        // integer literal
	TokenString     // string literal
	TokenBool       // bool literal
	TokenFloat      // float literal
	TokenIdentifier // variable name

	// Type keywords - integers
	TokenTypeInt    // type keyword: int
	TokenTypeInt8   // type keyword: int8
	TokenTypeInt16  // type keyword: int16
	TokenTypeInt32  // type keyword: int32
	TokenTypeInt64  // type keyword: int64
	TokenTypeUint   // type keyword: uint
	TokenTypeUint8  // type keyword: uint8
	TokenTypeUint16 // type keyword: uint16
	TokenTypeUint32 // type keyword: uint32
	TokenTypeUint64 // type keyword: uint64

	// Type keywords - other
	TokenTypeString // type keyword: string
	TokenTypeBool   // type keyword: bool
	TokenTypeFloat  // type keyword: float

	// Built-in print functions
	TokenPrintString
	TokenPrintf
	TokenFPrintf
	TokenPrintln
	TokenSPrint
	TokenSPrintf
	TokenSPrintln
	TokenFatalf
	TokenFatalln
	TokenLogf
	TokenLogln

	// Operators - arithmetic
	TokenPlus    // +
	TokenMinus   // -
	TokenStar    // *
	TokenSlash   // /
	TokenPercent // %

	// Operators - comparison
	TokenEqual     // ==
	TokenNotEqual  // !=
	TokenLess      // <
	TokenLessEq    // <=
	TokenGreater   // >
	TokenGreaterEq // >=

	// Operators - logical
	TokenAmpersand // &
	TokenExclaim   // !
	TokenPipe      // |
	TokenCaret     // ^
	TokenTilde     // ~
	TokenAnd       // &&
	TokenOr        // ||

	// Operators - bitwise
	TokenLShift // <<
	TokenRShift // >>

	// Operators - increment/decrement
	TokenPlusPlus   // ++
	TokenMinusMinus // --

	// Operators - compound assignment
	TokenPlusEq    // +=
	TokenMinusEq   // -=
	TokenStarEq    // *=
	TokenSlashEq   // /=
	TokenPercentEq // %=

	// Delimiters
	TokenLParen   // (
	TokenRParen   // )
	TokenLBrace   // {
	TokenRBrace   // }
	TokenLBracket // [
	TokenRBracket // ]

	// Punctuation
	TokenComma    // ,
	TokenSemi     // ;
	TokenDot      // .
	TokenColon    // :
	TokenArrow    // ->
	TokenQuestion // ?
	TokenNewline  // newline
	TokenAssign   // =

	// Data structures
	TokenStruct // struct
	TokenEnum   // enum
	TokenClass  // class

	// Memory management
	TokenNew // new

	// Error handling
	TokenTry     // try
	TokenCatch   // catch
	TokenFinally // finally
	TokenThrow   // throw
	TokenNull    // null

	// Special tokens
	TokenEOF
	TokenUnknown
)

// Token represents a lexical token with its type and value
type Token struct {
	Type  TokenType // The token type
	Value string    // The raw string value from source code
}
