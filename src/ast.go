package main

// ast.go - Central AST node definitions for the Lotus compiler
// This file consolidates all Abstract Syntax Tree node definitions in one place
// for better maintainability and organization.

// ASTNode is the interface that all AST nodes must implement
type ASTNode interface {
	astNode()
}

// ============================================================================
// Statement Nodes
// ============================================================================

// ReturnStatement represents a return statement with an optional value
type ReturnStatement struct {
	Value ASTNode
}

func (r *ReturnStatement) astNode() {}

// VariableDeclaration represents a variable declaration with type and initial value
type VariableDeclaration struct {
	Name  string
	Type  TokenType
	Value ASTNode
}

func (v *VariableDeclaration) astNode() {}

// ConstantDeclaration represents a constant declaration with type and value
// Constants are immutable and their values must be compile-time evaluable
type ConstantDeclaration struct {
	Name  string
	Type  TokenType
	Value ASTNode
}

func (c *ConstantDeclaration) astNode() {}

// ============================================================================
// Literal Expression Nodes
// ============================================================================

// IntLiteral represents an integer constant
type IntLiteral struct {
	Value int
}

func (i *IntLiteral) astNode() {}

// StringLiteral represents a string constant
type StringLiteral struct {
	Value string
}

func (s *StringLiteral) astNode() {}

// BoolLiteral represents a boolean constant (true/false)
type BoolLiteral struct {
	Value bool
}

func (b *BoolLiteral) astNode() {}

// FloatLiteral represents a floating-point constant
// Value is stored as int * 1000 for precision
type FloatLiteral struct {
	Value int64
}

func (f *FloatLiteral) astNode() {}

// NullLiteral represents a null value
type NullLiteral struct{}

func (n *NullLiteral) astNode() {}

// ============================================================================
// Identifier and Function Call Nodes
// ============================================================================

// Identifier represents a variable or symbol name reference
type Identifier struct {
	Name string
}

func (id *Identifier) astNode() {}

// FunctionCall represents a function invocation with arguments
type FunctionCall struct {
	Name string
	Args []ASTNode
}

func (f *FunctionCall) astNode() {}
