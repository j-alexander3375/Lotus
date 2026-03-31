package main

// error_handling.go - Exception handling AST nodes
// This file defines AST nodes for try/catch/finally blocks and throw statements.

// TryStatement represents a try-catch-finally block
type TryStatement struct {
	BaseNode
	TryBlock     []ASTNode      // Statements in the try block
	CatchClauses []*CatchClause // One or more catch blocks
	FinallyBlock []ASTNode      // Optional finally block
}

func (t *TryStatement) astNode() {}

// CatchClause represents a catch block with optional exception type
type CatchClause struct {
	BaseNode
	ExceptionType string    // Type of exception to catch (empty for catch-all)
	ExceptionVar  string    // Variable name to store the exception
	Body          []ASTNode // Statements in the catch block
}

func (c *CatchClause) astNode() {}

// ThrowStatement represents throwing an exception
type ThrowStatement struct {
	BaseNode
	Exception ASTNode // Expression to throw as exception
}

func (t *ThrowStatement) astNode() {}

// Note: NullLiteral is defined in ast.go
