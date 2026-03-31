package main

// Reference represents a pointer/reference to a value
type Reference struct {
	BaseNode
	Target ASTNode
}

func (r *Reference) astNode() {}

// Dereference represents dereferencing a pointer
type Dereference struct {
	BaseNode
	Pointer ASTNode
}

func (d *Dereference) astNode() {}

// Assignment represents variable assignment (reassignment)
type Assignment struct {
	BaseNode
	Target ASTNode
	Value  ASTNode
}

func (a *Assignment) astNode() {}

// CompoundAssignment represents compound assignment (+=, -=, *=, /=, %=)
type CompoundAssignment struct {
	BaseNode
	Target   ASTNode
	Operator TokenType // TokenPlusEq, TokenMinusEq, TokenStarEq, TokenSlashEq, TokenPercentEq
	Value    ASTNode
}

func (c *CompoundAssignment) astNode() {}
