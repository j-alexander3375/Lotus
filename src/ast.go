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

// StorageClass represents the storage class of a variable or function
type StorageClass int

const (
	StorageAuto   StorageClass = iota // Default: local for function vars, global for top-level
	StorageStatic                     // static: persistent storage, file-local scope
	StorageLocal                      // lcl: explicitly local scope (stack allocated)
	StorageGlobal                     // gbl: explicitly global scope (heap/data section)
)

// ReturnStatement represents a return statement with an optional value
type ReturnStatement struct {
	BaseNode
	Value ASTNode
}

func (r *ReturnStatement) astNode() {}

// VariableDeclaration represents a variable declaration with type and initial value
type VariableDeclaration struct {
	BaseNode
	Name    string
	Type    TokenType
	Value   ASTNode
	Storage StorageClass // Storage class modifier (static, lcl, gbl)
	IsArray bool         // Whether this is an array type declaration (e.g., int[])
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

// InterpolatedString represents a string with embedded expressions
// Example: $"Hello {name}, you are {age} years old"
type InterpolatedString struct {
	BaseNode
	Parts []InterpolatedPart // Alternating text and expression parts
}

func (i *InterpolatedString) astNode() {}

// InterpolatedPart is either a text segment or an expression
type InterpolatedPart struct {
	IsExpr bool    // true if this is an expression, false if text
	Text   string  // The text content (if IsExpr is false)
	Expr   ASTNode // The expression (if IsExpr is true)
}

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

// PartialApplication represents partial function application (currying)
// partial(fn_name, arg1, arg2, ...) returns a closure capturing the provided args
type PartialApplication struct {
	BaseNode
	FunctionName string    // The function to partially apply
	BoundArgs    []ASTNode // Arguments already bound
}

func (p *PartialApplication) astNode() {}

// WrapperDefinition defines a wrapper/decorator function
// wrap fn timing(fn wrapped) { ... before ...; wrapped(); ... after ... }
type WrapperDefinition struct {
	BaseNode
	Name       string    // Wrapper name (e.g., "timing")
	WrappedArg string    // Name of the wrapped function parameter
	Body       []ASTNode // Wrapper body
}

func (w *WrapperDefinition) astNode() {}

// DecoratedFunction represents a function with one or more decorators
// @timing @logging fn foo() { ... }
type DecoratedFunction struct {
	BaseNode
	Decorators []string            // List of decorator names to apply
	Function   *FunctionDefinition // The wrapped function
}

func (d *DecoratedFunction) astNode() {}

// PipeExpression represents the pipe operator for function chaining
// value |> fn1 |> fn2 becomes fn2(fn1(value))
type PipeExpression struct {
	BaseNode
	Left     ASTNode // The value being piped
	Function string  // The function to pipe into
}

func (p *PipeExpression) astNode() {}

// BitCastExpression represents a bitwise type reinterpretation
// bitcast<TargetType>(value) - reinterprets the bits of value as TargetType
// Example: bitcast<int64>(floatVal) to access float bits as integer
type BitCastExpression struct {
	BaseNode
	TargetType string  // The target type to cast to (e.g., "int64", "float")
	Value      ASTNode // The expression to bitcast
}

func (b *BitCastExpression) astNode() {}

// ============================================================================
// Generics Nodes
// ============================================================================

// TypeParameter represents a generic type parameter (e.g., T in Box<T>)
type TypeParameter struct {
	Name string // Type parameter name (e.g., "T", "K", "V")
}

// GenericFunctionDefinition represents a generic function with type parameters
// fn T max<T>(T a, T b) { ret a > b ? a : b; }
type GenericFunctionDefinition struct {
	BaseNode
	Name           string          // Function name
	TypeParams     []TypeParameter // Type parameters (e.g., [T, K, V])
	ReturnType     TokenType       // Return type (can use type parameter)
	ReturnTypeName string          // Name if return type is a type parameter
	Parameters     []FunctionParam // Function parameters
	Body           []ASTNode       // Function body
}

func (g *GenericFunctionDefinition) astNode() {}

// GenericStructDefinition represents a generic struct with type parameters
// struct Box<T> { T value; }
type GenericStructDefinition struct {
	BaseNode
	Name       string          // Struct name
	TypeParams []TypeParameter // Type parameters
	Fields     []StructField   // Struct fields (can use type parameters)
}

func (g *GenericStructDefinition) astNode() {}

// GenericInstantiation represents the use of a generic type with concrete types
// Box<int>, Array<string>, HashMap<string, int>
type GenericInstantiation struct {
	BaseNode
	Name       string   // Generic name (e.g., "Box", "Array")
	TypeArgs   []string // Type arguments (e.g., ["int"], ["string", "int"])
	IsStruct   bool     // true if this is a struct instantiation
	IsFunction bool     // true if this is a function call
}

func (g *GenericInstantiation) astNode() {}
