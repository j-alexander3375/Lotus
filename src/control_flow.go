package main

// IfStatement represents an if/else control flow structure
type IfStatement struct {
	BaseNode
	Condition ASTNode
	ThenBody  []ASTNode
	ElseBody  []ASTNode
}

func (i *IfStatement) astNode() {}

// WhileLoop represents a while loop structure
type WhileLoop struct {
	BaseNode
	Condition ASTNode
	Body      []ASTNode
}

func (w *WhileLoop) astNode() {}

// ForLoop represents a for loop structure
type ForLoop struct {
	BaseNode
	Init      ASTNode
	Condition ASTNode
	Update    ASTNode
	Body      []ASTNode
}

func (f *ForLoop) astNode() {}

// BreakStatement represents a break statement to exit a loop
type BreakStatement struct {
	BaseNode
}

func (b *BreakStatement) astNode() {}

// ContinueStatement represents a continue statement to skip to next iteration
type ContinueStatement struct {
	BaseNode
}

func (c *ContinueStatement) astNode() {}

// Comparison represents a comparison expression
type Comparison struct {
	BaseNode
	Left     ASTNode
	Operator TokenType // TokenEqual, TokenNotEqual, TokenLess, TokenLessEq, TokenGreater, TokenGreaterEq
	Right    ASTNode
}

func (c *Comparison) astNode() {}

// LogicalOp represents logical operations (&&, ||)
type LogicalOp struct {
	BaseNode
	Left     ASTNode
	Operator TokenType // TokenAnd, TokenOr
	Right    ASTNode
}

func (l *LogicalOp) astNode() {}

// TernaryOp represents ternary conditional operator (condition ? trueExpr : falseExpr)
type TernaryOp struct {
	BaseNode
	Condition ASTNode
	TrueExpr  ASTNode
	FalseExpr ASTNode
}

func (t *TernaryOp) astNode() {}

// ============================================================================
// PATTERN MATCHING
// ============================================================================

// MatchExpression represents a pattern matching expression
// Syntax: match expr { case pattern => body, ... }
type MatchExpression struct {
	BaseNode
	Value ASTNode     // The value to match against
	Cases []MatchCase // List of cases
}

func (m *MatchExpression) astNode() {}

// MatchYield represents the implicit "produce this value as the match
// expression's result" sugar for a single-expression arm (`case X => expr`).
// It is distinct from ReturnStatement, which the parser previously reused
// for this - conflating the sugar with a genuine explicit `ret` written
// inside a block-bodied arm (`case X => { ret 5; }`), so codegen could not
// tell them apart and the explicit `ret` never actually returned from the
// function. See SP-B-5 in FIXER_HANDOFF.md.
type MatchYield struct {
	BaseNode
	Value ASTNode
}

func (m *MatchYield) astNode() {}

// MatchCase represents a single case in a match expression
type MatchCase struct {
	Pattern   ASTNode   // The pattern to match (literal, identifier, wildcard)
	Guard     ASTNode   // Optional guard condition (when clause)
	Body      []ASTNode // Body to execute if matched
	IsDefault bool      // True if this is the default case
}

// Pattern types for matching
type WildcardPattern struct {
	BaseNode
}

func (w *WildcardPattern) astNode() {}

type LiteralPattern struct {
	BaseNode
	Value ASTNode // IntLiteral, StringLiteral, etc.
}

func (l *LiteralPattern) astNode() {}

type BindingPattern struct {
	BaseNode
	Name string // Variable name to bind the matched value
}

func (b *BindingPattern) astNode() {}

type RangePattern struct {
	BaseNode
	Start ASTNode // Start of range (inclusive)
	End   ASTNode // End of range (inclusive)
}

func (r *RangePattern) astNode() {}

// ============================================================================
// UNION AND OPTION TYPES
// ============================================================================

// UnionDefinition represents a union type definition
// Syntax: union Name { Variant1(Type), Variant2(Type1, Type2), Variant3 }
type UnionDefinition struct {
	BaseNode
	Name     string
	Variants []UnionVariant
}

func (u *UnionDefinition) astNode() {}

// UnionVariant represents a variant of a union type
type UnionVariant struct {
	Name   string      // Variant name
	Fields []TokenType // Field types (can be empty for unit variants)
}

// OptionExpression represents Some(value) or None
type OptionExpression struct {
	BaseNode
	IsSome bool    // true for Some, false for None
	Value  ASTNode // The value (only for Some)
}

func (o *OptionExpression) astNode() {}

// UnionLiteral represents instantiation of a union variant
// Syntax: VariantName(args...) or VariantName
type UnionLiteral struct {
	BaseNode
	UnionName   string    // The union type name
	VariantName string    // The variant being constructed
	Args        []ASTNode // Constructor arguments
}

func (u *UnionLiteral) astNode() {}
