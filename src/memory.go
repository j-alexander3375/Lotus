package main

import "fmt"

// memory.go - Memory management and sizeof operations
// This file handles dynamic memory allocation (malloc/free) and sizeof expressions.

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

// generateMallocCall generates assembly for dynamic memory allocation.
// Calls the system malloc() function via PLT (Procedure Linkage Table).
func (cg *CodeGenerator) generateMallocCall(malloc *MallocCall) {
	// Evaluate size expression and place result in rdi (first argument)
	cg.generateExpressionToReg(malloc.Size, "rdi")

	// Align stack to 16 bytes before call (System V ABI requirement)
	cg.textSection.WriteString("    pushq %rbp\n")
	cg.textSection.WriteString("    movq %rsp, %rbp\n")
	cg.textSection.WriteString("    andq $-16, %rsp\n")

	// Call malloc from C standard library
	cg.textSection.WriteString("    call malloc@PLT\n")

	// Restore stack pointer
	cg.textSection.WriteString("    movq %rbp, %rsp\n")
	cg.textSection.WriteString("    popq %rbp\n")

	// Result pointer is in rax
}

// generateFreeCall generates assembly for memory deallocation.
// Calls the system free() function via PLT.
func (cg *CodeGenerator) generateFreeCall(free *FreeCall) {
	// Evaluate pointer expression and place result in rdi (first argument)
	cg.generateExpressionToReg(free.Pointer, "rdi")

	// Align stack to 16 bytes before call (System V ABI requirement)
	cg.textSection.WriteString("    pushq %rbp\n")
	cg.textSection.WriteString("    movq %rsp, %rbp\n")
	cg.textSection.WriteString("    andq $-16, %rsp\n")

	// Call free from C standard library
	cg.textSection.WriteString("    call free@PLT\n")

	// Restore stack pointer
	cg.textSection.WriteString("    movq %rbp, %rsp\n")
	cg.textSection.WriteString("    popq %rbp\n")
}

// generateSizeofExpr generates assembly for sizeof expression.
// Evaluates to a compile-time constant representing the size in bytes.
func (cg *CodeGenerator) generateSizeofExpr(sizeof *SizeofExpr) {
	// Determine size based on type
	size := cg.getSizeOfType(sizeof.TypeOrExpr)
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rax  # sizeof result\n", size))
}

// getSizeOfType returns the size in bytes of a type from an AST node
func (cg *CodeGenerator) getSizeOfType(node ASTNode) int {
	switch t := node.(type) {
	case *Identifier:
		// Identifier might be a type name
		return cg.getTypeSize(t.Name)
	default:
		// Default to pointer size for unknown types
		return PointerSize
	}
}

// getTypeSize returns the size of a type by its string name
func (cg *CodeGenerator) getTypeSize(typeName string) int {
	switch typeName {
	case "int8", "uint8", "bool":
		return Int8Size
	case "int16", "uint16":
		return Int16Size
	case "int32", "uint32", "float":
		return Int32Size
	case "int64", "uint64", "int", "uint", "string":
		return Int64Size
	default:
		// For user-defined types (structs, classes), would query type registry
		return PointerSize
	}
}

// getTypeSizeFromToken returns the size in bytes of a type from its TokenType.
// This is used during code generation for type-specific operations.
func getTypeSizeFromToken(tokenType TokenType) int {
	// Use the centralized type size utility
	return GetTypeSize(tokenType)
}
