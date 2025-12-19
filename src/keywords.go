package main

type TokenType int

const (
	TokenRet TokenType = iota
	TokenInt           // integer literal
	TokenSemi
	TokenEOF
	TokenString // string literal
	TokenBool   // bool literal
	TokenNewline
	TokenAssign
	TokenFloat      // float literal
	TokenIdentifier // variable name
	TokenTypeInt    // type keyword: int
	TokenTypeString // type keyword: string
	TokenTypeBool   // type keyword: bool
	TokenTypeFloat  // type keyword: float
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
	TokenLParen    // (
	TokenRParen    // )
	TokenComma     // ,
	TokenPlus      // +
	TokenMinus     // -
	TokenStar      // *
	TokenSlash     // /
	TokenPercent   // %
	TokenEqual     // ==
	TokenNotEqual  // !=
	TokenLess      // <
	TokenLessEq    // <=
	TokenGreater   // >
	TokenGreaterEq // >=
	TokenAmpersand // &
	TokenExclaim   // !
	TokenLBrace    // {
	TokenRBrace    // }
	TokenIf        // if
	TokenElse      // else
	TokenWhile     // while
	TokenFor       // for
	TokenFn        // fn
	TokenReturn    // return
	TokenUnknown
)

type Token struct {
	Type  TokenType
	Value string
}

// Variable represents a variable in the symbol table
type Variable struct {
	Name   string
	Type   TokenType
	Offset int // stack offset from rbp
}
