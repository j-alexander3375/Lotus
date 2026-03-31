package main

// memory.go - Memory management and sizeof operations
// This file defines AST nodes for dynamic memory allocation and sizeof expressions.

// MallocCall represents a memory allocation call: malloc(size)
type MallocCall struct {
	BaseNode
	Size ASTNode // Size in bytes to allocate
}

func (m *MallocCall) astNode() {}

// FreeCall represents a memory deallocation call: free(ptr)
type FreeCall struct {
	BaseNode
	Pointer ASTNode // Pointer to memory to free
}

func (f *FreeCall) astNode() {}

// SizeofExpr represents a sizeof expression: sizeof(type)
type SizeofExpr struct {
	BaseNode
	TypeOrExpr ASTNode // Type or expression to get size of
}

func (s *SizeofExpr) astNode() {}
