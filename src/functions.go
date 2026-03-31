package main

// FunctionDefinition represents a user-defined function
type FunctionDefinition struct {
	BaseNode
	Name       string
	Parameters []FunctionParam
	ReturnType TokenType
	Body       []ASTNode
	IsVirtual  bool         // true if marked with 'vrt' keyword
	IsOverride bool         // true if marked with 'override' keyword
	IsStatic   bool         // true if marked with 'static' keyword
	IsInline   bool         // true if marked with 'inline' keyword
	Storage    StorageClass // Storage class (static makes function file-local)
}

func (f *FunctionDefinition) astNode() {}

// LambdaExpression represents an anonymous function / lambda
// Syntax: |params| => expr  or  |params| => { body }
type LambdaExpression struct {
	BaseNode
	Parameters []FunctionParam
	ReturnType TokenType
	Body       []ASTNode // For block body { ... }
	Expression ASTNode   // For single expression body (either Body or Expression is set)
}

func (l *LambdaExpression) astNode() {}

// FunctionReference represents a reference to a function as a value
// Syntax: fn funcName or &funcName
type FunctionReference struct {
	BaseNode
	FunctionName string
}

func (f *FunctionReference) astNode() {}

// FunctionParam represents a function parameter
type FunctionParam struct {
	Name string
	Type TokenType
}

// VTableEntry represents an entry in a virtual method table
type VTableEntry struct {
	MethodName   string    // Method name
	ReturnType   TokenType // Return type
	Parameters   []FunctionParam
	ImplLabel    string // Label for implementation
	VTableOffset int    // Offset in vtable
}

// VTable represents a virtual method table for a class
type VTable struct {
	ClassName string
	Entries   []VTableEntry
	Label     string // Label for vtable
}

// VTableRegistry maps class names to their vtables
var VTableRegistry = make(map[string]*VTable)

// FunctionContext holds information about a function during code generation
type FunctionContext struct {
	Name        string
	Parameters  map[string]Variable
	LocalVars   map[string]Variable
	StackSize   int
	ReturnLabel string
}

// Global registry of function definitions
var UserDefinedFunctions = make(map[string]*FunctionDefinition)
