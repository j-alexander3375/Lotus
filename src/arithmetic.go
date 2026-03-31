package main

// BinaryOp represents a binary operation expression
type BinaryOp struct {
	BaseNode
	Left     ASTNode
	Operator TokenType // TokenPlus, TokenMinus, TokenStar, TokenSlash, TokenPercent
	Right    ASTNode
}

func (b *BinaryOp) astNode() {}

// UnaryOp represents a unary operation expression
type UnaryOp struct {
	BaseNode
	Operator TokenType // TokenMinus, TokenExclaim, TokenAmpersand, TokenStar, TokenTilde
	Operand  ASTNode
}

func (u *UnaryOp) astNode() {}

// BitwiseOp represents bitwise operations (&, |, ^, <<, >>)
type BitwiseOp struct {
	BaseNode
	Left     ASTNode
	Operator TokenType // TokenAmpersand, TokenPipe, TokenCaret, TokenLShift, TokenRShift
	Right    ASTNode
}

func (b *BitwiseOp) astNode() {}

// IncrementOp represents increment/decrement (++, --)
type IncrementOp struct {
	BaseNode
	Operand  ASTNode
	IsPrefix bool      // true for ++x, false for x++
	Operator TokenType // TokenPlusPlus, TokenMinusMinus
}

func (i *IncrementOp) astNode() {}
