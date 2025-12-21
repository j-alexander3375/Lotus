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
	TokenTypeInt8   // type keyword: int8
	TokenTypeInt16  // type keyword: int16
	TokenTypeInt32  // type keyword: int32
	TokenTypeInt64  // type keyword: int64
	TokenTypeUint   // type keyword: uint
	TokenTypeUint8  // type keyword: uint8
	TokenTypeUint16 // type keyword: uint16
	TokenTypeUint32 // type keyword: uint32
	TokenTypeUint64 // type keyword: uint64
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
	TokenPipe      // |
	TokenCaret     // ^
	TokenTilde     // ~
	TokenAnd       // &&
	TokenOr        // ||
	TokenLShift    // <<
	TokenRShift    // >>
	TokenPlusPlus  // ++
	TokenMinusMinus // --
	TokenPlusEq    // +=
	TokenMinusEq   // -=
	TokenStarEq    // *=
	TokenSlashEq   // /=
	TokenPercentEq // %=
	TokenQuestion  // ?
	TokenIf        // if
	TokenElse      // else
	TokenWhile     // while
	TokenFor       // for
	TokenFn        // fn
	TokenReturn    // return
	TokenStruct    // struct
	TokenEnum      // enum
	TokenClass     // class
	TokenNew       // new
	TokenMalloc    // malloc
	TokenFree      // free
	TokenSizeof    // sizeof
	TokenLBracket  // [
	TokenRBracket  // ]
	TokenDot       // .
	TokenColon     // :
	TokenArrow     // ->
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
