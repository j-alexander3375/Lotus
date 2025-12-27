package main

// ast.go - Central AST node definitions for the Lotus compiler
// This file consolidates all Abstract Syntax Tree node definitions in one place
// for better maintainability and organization.

// Location stores source position information for an AST node
type Location struct {
	Line   int // Line number (1-based)
	Column int // Column number (1-based)
}

// ASTNode is the interface that all AST nodes must implement
type ASTNode interface {
	astNode()
	Loc() Location // Returns the source location of the node
}

// BaseNode provides common functionality for all AST nodes
type BaseNode struct {
	Location
}

// Loc returns the source location
func (b BaseNode) Loc() Location {
	return b.Location
}

// ============================================================================
// Statement Nodes
// ============================================================================

// ReturnStatement represents a return statement with an optional value
type ReturnStatement struct {
	BaseNode
	Value ASTNode
}

func (r *ReturnStatement) astNode() {}

// VariableDeclaration represents a variable declaration with type and initial value
type VariableDeclaration struct {
	BaseNode
	Name  string
	Type  TokenType
	Value ASTNode
}

func (v *VariableDeclaration) astNode() {}

// ConstantDeclaration represents a constant declaration with type and value
// Constants are immutable and their values must be compile-time evaluable
type ConstantDeclaration struct {
	BaseNode
	Name  string
	Type  TokenType
	Value ASTNode
}

func (c *ConstantDeclaration) astNode() {}

// ============================================================================
// Import/Module Nodes
// ============================================================================

// ImportStatement represents a use/import statement
// Examples:
//
//	use "io" - imports all io functions
//	use "math::sqrt" - imports specific function
//	use "math::*" - wildcard import
//	use "io" as io_module - aliased import
type ImportStatement struct {
	BaseNode
	Module     string   // e.g., "io", "math", "std::collections"
	Items      []string // Specific items to import, nil for all
	Alias      string   // Optional alias name
	IsWildcard bool     // true if use "module::*"
}

func (i *ImportStatement) astNode() {}

// ============================================================================
// Literal Expression Nodes
// ============================================================================

// IntLiteral represents an integer constant
type IntLiteral struct {
	BaseNode
	Value int
}

func (i *IntLiteral) astNode() {}

// StringLiteral represents a string constant
type StringLiteral struct {
	BaseNode
	Value string
}

func (s *StringLiteral) astNode() {}

// CharLiteral represents a single Unicode character (32-bit code point)
type CharLiteral struct {
	BaseNode
	Value string // Single Unicode character as string
}

func (c *CharLiteral) astNode() {}

// BoolLiteral represents a boolean constant (true/false)
type BoolLiteral struct {
	BaseNode
	Value bool
}

func (b *BoolLiteral) astNode() {}

// FloatLiteral represents a floating-point constant
// Value is stored as int * 1000 for precision
type FloatLiteral struct {
	BaseNode
	Value int64
}

func (f *FloatLiteral) astNode() {}

// NullLiteral represents a null value
type NullLiteral struct {
	BaseNode
}

func (n *NullLiteral) astNode() {}

// ============================================================================
// Identifier and Function Call Nodes
// ============================================================================

// Identifier represents a variable or symbol name reference
type Identifier struct {
	BaseNode
	Name string
}

func (id *Identifier) astNode() {}

// FunctionCall represents a function invocation with arguments
type FunctionCall struct {
	BaseNode
	Name string
	Args []ASTNode
}

func (f *FunctionCall) astNode() {}
