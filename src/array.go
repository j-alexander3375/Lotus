package main

// ArrayLiteral represents an array literal [1, 2, 3]
type ArrayLiteral struct {
	BaseNode
	Elements []ASTNode
	ElemType TokenType
}

func (a *ArrayLiteral) astNode() {}

// ArrayAccess represents array indexing: arr[index]
type ArrayAccess struct {
	BaseNode
	Name  string // For backwards compatibility with semantic analyzer
	Array ASTNode
	Index ASTNode
}

func (a *ArrayAccess) astNode() {}

// ArrayDeclaration represents dynamic array declaration
type ArrayDeclaration struct {
	BaseNode
	Name     string
	ElemType TokenType
	Size     ASTNode   // size expression or nil for dynamic
	Initial  []ASTNode // initial values
}

func (a *ArrayDeclaration) astNode() {}

// DynamicArray represents a dynamic array structure
type DynamicArray struct {
	Data     int64 // pointer to data
	Length   int64 // current length
	Capacity int64 // allocated capacity
	ElemSize int   // size of each element
}
