package main

import "fmt"

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

// generateArrayLiteral generates assembly for array literals
func (cg *CodeGenerator) generateArrayLiteral(arr *ArrayLiteral) {
	elemSize := getTypeSizeFromToken(arr.ElemType)
	totalSize := len(arr.Elements) * elemSize

	// Allocate memory for array on heap
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdi\n", totalSize))
	cg.textSection.WriteString("    call malloc@PLT\n")
	cg.textSection.WriteString("    pushq %rax\n") // Save array pointer

	// Fill array with elements
	for i, elem := range arr.Elements {
		// Evaluate element
		cg.generateExpressionToReg(elem, "rcx")

		// Get array pointer
		cg.textSection.WriteString("    movq (%rsp), %rax\n")

		// Store element at offset
		offset := i * elemSize
		switch elemSize {
		case 1:
			cg.textSection.WriteString(fmt.Sprintf("    movb %%cl, %d(%%rax)\n", offset))
		case 2:
			cg.textSection.WriteString(fmt.Sprintf("    movw %%cx, %d(%%rax)\n", offset))
		case 4:
			cg.textSection.WriteString(fmt.Sprintf("    movl %%ecx, %d(%%rax)\n", offset))
		case 8:
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rcx, %d(%%rax)\n", offset))
		}
	}

	// Pop array pointer back to rax
	cg.textSection.WriteString("    popq %rax\n")
}

// generateArrayAccess generates assembly for array indexing
func (cg *CodeGenerator) generateArrayAccess(access *ArrayAccess) {
	// Get array base address into rax
	cg.generateExpressionToReg(access.Array, "rax")

	// Save array pointer
	cg.textSection.WriteString("    pushq %rax\n")

	// Evaluate index into rcx
	cg.generateExpressionToReg(access.Index, "rcx")

	// Pop array pointer
	cg.textSection.WriteString("    popq %rax\n")

	// Calculate element size (assume 8 bytes for now, should be determined by array type)
	elemSize := 8

	// Calculate offset: index * elemSize
	if elemSize == 8 {
		cg.textSection.WriteString("    shlq $3, %rcx\n") // multiply by 8 (2^3)
	} else {
		cg.textSection.WriteString(fmt.Sprintf("    imulq $%d, %%rcx\n", elemSize))
	}

	// Add offset to base address
	cg.textSection.WriteString("    addq %rcx, %rax\n")

	// Load value from memory address
	cg.textSection.WriteString("    movq (%rax), %rax\n")
}

// generateArrayDeclaration generates assembly for dynamic array declaration
func (cg *CodeGenerator) generateArrayDeclaration(decl *ArrayDeclaration) {
	elemSize := getTypeSizeFromToken(decl.ElemType)

	var arraySize int
	if decl.Size != nil {
		// Static size specified
		if intLit, ok := decl.Size.(*IntLiteral); ok {
			arraySize = intLit.Value
		} else {
			// Dynamic size - evaluate at runtime
			cg.generateExpressionToReg(decl.Size, "rcx")
			cg.textSection.WriteString(fmt.Sprintf("    imulq $%d, %%rcx\n", elemSize))
			cg.textSection.WriteString("    movq %rcx, %rdi\n")
			cg.textSection.WriteString("    call malloc@PLT\n")

			// Store array pointer in variable
			cg.stackOffset += 8
			cg.variables[decl.Name] = Variable{
				Name:   decl.Name,
				Type:   TokenTypeUint64, // pointer type
				Offset: cg.stackOffset,
			}
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, -%d(%%rbp)\n", cg.stackOffset))
			return
		}
	} else {
		// Use initial values to determine size
		arraySize = len(decl.Initial)
	}

	totalSize := arraySize * elemSize

	// Allocate array on heap
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdi\n", totalSize))
	cg.textSection.WriteString("    call malloc@PLT\n")

	// Store array pointer
	cg.stackOffset += 8
	cg.variables[decl.Name] = Variable{
		Name:   decl.Name,
		Type:   TokenTypeUint64, // pointer type
		Offset: cg.stackOffset,
	}
	cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, -%d(%%rbp)\n", cg.stackOffset))

	// Initialize with values if provided
	if len(decl.Initial) > 0 {
		for i, elem := range decl.Initial {
			cg.generateExpressionToReg(elem, "rcx")
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rax\n", cg.stackOffset))
			offset := i * elemSize

			switch elemSize {
			case 1:
				cg.textSection.WriteString(fmt.Sprintf("    movb %%cl, %d(%%rax)\n", offset))
			case 2:
				cg.textSection.WriteString(fmt.Sprintf("    movw %%cx, %d(%%rax)\n", offset))
			case 4:
				cg.textSection.WriteString(fmt.Sprintf("    movl %%ecx, %d(%%rax)\n", offset))
			case 8:
				cg.textSection.WriteString(fmt.Sprintf("    movq %%rcx, %d(%%rax)\n", offset))
			}
		}
	}
}

// Helper function to resize dynamic arrays (would be called from runtime)
func generateArrayResize(cg *CodeGenerator) {
	// This would be a runtime function to resize arrays
	// Implementation would involve:
	// 1. Allocate new larger array
	// 2. Copy existing elements
	// 3. Free old array
	// 4. Update pointer, capacity
}
