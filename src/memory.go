package main

import "fmt"

// MallocCall represents a memory allocation call
type MallocCall struct {
	Size ASTNode // size in bytes to allocate
}

func (m *MallocCall) astNode() {}

// FreeCall represents a memory deallocation call
type FreeCall struct {
	Pointer ASTNode // pointer to free
}

func (f *FreeCall) astNode() {}

// SizeofExpr represents a sizeof expression
type SizeofExpr struct {
	TypeOrExpr ASTNode // type or expression to get size of
}

func (s *SizeofExpr) astNode() {}

// generateMallocCall generates assembly for malloc call
func (cg *CodeGenerator) generateMallocCall(malloc *MallocCall) {
	// Evaluate size into rdi (first argument for malloc)
	cg.generateExpressionToReg(malloc.Size, "rdi")

	// Align stack to 16 bytes before call
	cg.textSection.WriteString("    pushq %rbp\n")
	cg.textSection.WriteString("    movq %rsp, %rbp\n")
	cg.textSection.WriteString("    andq $-16, %rsp\n")

	// Call malloc (system library function)
	cg.textSection.WriteString("    call malloc@PLT\n")

	// Restore stack
	cg.textSection.WriteString("    movq %rbp, %rsp\n")
	cg.textSection.WriteString("    popq %rbp\n")

	// Result is in rax (pointer to allocated memory)
}

// generateFreeCall generates assembly for free call
func (cg *CodeGenerator) generateFreeCall(free *FreeCall) {
	// Evaluate pointer into rdi (first argument for free)
	cg.generateExpressionToReg(free.Pointer, "rdi")

	// Align stack to 16 bytes before call
	cg.textSection.WriteString("    pushq %rbp\n")
	cg.textSection.WriteString("    movq %rsp, %rbp\n")
	cg.textSection.WriteString("    andq $-16, %rsp\n")

	// Call free (system library function)
	cg.textSection.WriteString("    call free@PLT\n")

	// Restore stack
	cg.textSection.WriteString("    movq %rbp, %rsp\n")
	cg.textSection.WriteString("    popq %rbp\n")
}

// generateSizeofExpr generates assembly for sizeof expression
func (cg *CodeGenerator) generateSizeofExpr(sizeof *SizeofExpr) {
	// Determine size based on type
	size := cg.getSizeOfType(sizeof.TypeOrExpr)
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rax\n", size))
}

// getSizeOfType returns the size in bytes of a type
func (cg *CodeGenerator) getSizeOfType(node ASTNode) int {
	switch t := node.(type) {
	case *Identifier:
		// Check if it's a type name
		typeName := t.Name
		return cg.getTypeSize(typeName)
	default:
		// Default to pointer size (8 bytes on x86-64)
		return 8
	}
}

// getTypeSize returns the size of a type by name
func (cg *CodeGenerator) getTypeSize(typeName string) int {
	switch typeName {
	case "int8", "uint8", "bool":
		return 1
	case "int16", "uint16":
		return 2
	case "int32", "uint32", "float":
		return 4
	case "int64", "uint64", "int", "uint":
		return 8
	case "string": // pointer to string data
		return 8
	default:
		// For user-defined types, would need type registry
		return 8
	}
}

// getTypeSizeFromToken returns the size of a type from TokenType
func getTypeSizeFromToken(tokenType TokenType) int {
	switch tokenType {
	case TokenTypeInt8, TokenTypeUint8, TokenTypeBool:
		return 1
	case TokenTypeInt16, TokenTypeUint16:
		return 2
	case TokenTypeInt32, TokenTypeUint32, TokenTypeFloat:
		return 4
	case TokenTypeInt64, TokenTypeUint64, TokenTypeInt, TokenTypeUint:
		return 8
	case TokenTypeString:
		return 8 // pointer
	default:
		return 8
	}
}
