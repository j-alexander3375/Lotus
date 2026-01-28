package main

// keywords.go - Token type definitions and lexical token structure
// This file defines all token types used by the lexer and parser.

// TokenType represents the type of a lexical token
type TokenType int

// Token type constants - organized by category
const (
	// Keywords and control flow
	TokenRet      TokenType = iota
	TokenReturn             // return
	TokenConst              // const
	TokenIf                 // if
	TokenElse               // else
	TokenWhile              // while
	TokenFor                // for
	TokenBreak              // break
	TokenContinue           // continue
	TokenFn                 // fn
	TokenInline             // inline (inline function hint)
	TokenLambda             // => (lambda arrow)
	TokenUse                // use (imports)
	TokenAs                 // as (aliasing)

	// Literals
	TokenInt                // integer literal
	TokenString             // string literal
	TokenInterpolatedString // interpolated string: $"Hello {name}"
	TokenChar               // character literal (single Unicode character)
	TokenBool               // bool literal
	TokenFloat              // float literal
	TokenIdentifier         // variable name

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

	// Type keywords - floats
	TokenTypeFloat32 // type keyword: float32 (32-bit single precision)
	TokenTypeFloat64 // type keyword: float64 (64-bit double precision)

	// Type keywords - other
	TokenTypeString // type keyword: string
	TokenTypeChar   // type keyword: char (Unicode character, 32-bit)
	TokenTypeBool   // type keyword: bool
	TokenTypeVoid   // type keyword: void (for functions with no return value)
	TokenTypeFn     // type keyword: fn (for function types / first-class functions)

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

	// Virtual functions and inheritance
	TokenVirtual  // vrt (virtual function marker)
	TokenOverride // override (method override marker)

	// Scope and storage modifiers
	TokenStatic // static (persistent storage, file-local scope)
	TokenLocal  // lcl (explicitly local scope)
	TokenGlobal // gbl (explicitly global scope)

	// Functional programming
	TokenPartial // partial (partial function application / currying)
	TokenWrap    // wrap (wrapper/decorator definition)
	TokenAt      // @ (decorator application)
	TokenPipe2   // |> (pipe operator for function chaining)

	// Pattern matching
	TokenMatch   // match (pattern matching expression)
	TokenCase    // case (pattern case)
	TokenDefault // default (default case)
	TokenWhen    // when (guard clause)
	TokenTo      // to (range pattern: 1 to 10)
	TokenRange   // .. (range operator: 1..10)

	// Union/Option types
	TokenUnion  // union (union type definition)
	TokenOption // option (optional type)
	TokenSome   // Some (option with value)
	TokenNone   // None (empty option)

	// Memory management
	TokenNew // new

	// Type casting and reinterpretation
	TokenBitcast   // bitcast (reinterpret bits as different type)
	TokenTransmute // transmute (alias for bitcast)

	// Error handling
	TokenTry     // try
	TokenCatch   // catch
	TokenFinally // finally
	TokenThrow   // throw
	TokenNull    // null

	// Templates (C++ style generics)
	TokenTemplate // template (template declaration)
	TokenTypename // typename (type parameter in template)

	// Namespaces
	TokenNamespace // namespace (namespace declaration)
	TokenUsing     // using (using declaration/directive)
	TokenInline2   // inline (inline namespace)

	// Special tokens
	TokenEOF
	TokenUnknown
)

// Token represents a lexical token with its type and value
type Token struct {
	Type   TokenType // The token type
	Value  string    // The raw string value from source code
	Line   int       // Line number in source (1-based)
	Column int       // Column number in source (1-based)
}
