package main

import "fmt"

// Reference represents a pointer/reference to a value
type Reference struct {
	Target ASTNode
}

func (r *Reference) astNode() {}

// Dereference represents dereferencing a pointer
type Dereference struct {
	Pointer ASTNode
}

func (d *Dereference) astNode() {}

// Assignment represents variable assignment (reassignment)
type Assignment struct {
	Target ASTNode
	Value  ASTNode
}

func (a *Assignment) astNode() {}

// CompoundAssignment represents compound assignment (+=, -=, *=, /=, %=)
type CompoundAssignment struct {
	Target   ASTNode
	Operator TokenType // TokenPlusEq, TokenMinusEq, TokenStarEq, TokenSlashEq, TokenPercentEq
	Value    ASTNode
}

func (c *CompoundAssignment) astNode() {}

// generateAssignment generates assembly for variable assignment
func (cg *CodeGenerator) generateAssignment(assign *Assignment) {
	if id, ok := assign.Target.(*Identifier); ok {
		if v, exists := cg.variables[id.Name]; exists {
			// Evaluate right side into rax
			cg.generateExpressionToReg(assign.Value, "rax")
			// Store into variable location
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, -%d(%%rbp)\n", v.Offset))

			// Update string length if assigning string
			if str, ok := assign.Value.(*StringLiteral); ok {
				cg.stringLengths[id.Name] = len(str.Value)
			}
		}
	}
}

// generateCompoundAssignment generates assembly for compound assignment operators
func (cg *CodeGenerator) generateCompoundAssignment(compAssign *CompoundAssignment) {
	if id, ok := compAssign.Target.(*Identifier); ok {
		if v, exists := cg.variables[id.Name]; exists {
			// Load current value into rax
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rax\n", v.Offset))

			// Save current value
			cg.textSection.WriteString("    pushq %rax\n")

			// Evaluate right side into rax
			cg.generateExpressionToReg(compAssign.Value, "rax")

			// Pop left value to rcx
			cg.textSection.WriteString("    movq %rax, %rcx\n")
			cg.textSection.WriteString("    popq %rax\n")

			// Perform operation
			switch compAssign.Operator {
			case TokenPlusEq:
				cg.textSection.WriteString("    addq %rcx, %rax\n")
			case TokenMinusEq:
				cg.textSection.WriteString("    subq %rcx, %rax\n")
			case TokenStarEq:
				cg.textSection.WriteString("    imulq %rcx, %rax\n")
			case TokenSlashEq:
				cg.textSection.WriteString("    cqo\n")
				cg.textSection.WriteString("    idivq %rcx\n")
			case TokenPercentEq:
				cg.textSection.WriteString("    cqo\n")
				cg.textSection.WriteString("    idivq %rcx\n")
				cg.textSection.WriteString("    movq %rdx, %rax\n")
			}

			// Store result back
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, -%d(%%rbp)\n", v.Offset))
		}
	}
}

// generateReference generates assembly for taking a reference (&x)
func (cg *CodeGenerator) generateReference(ref *Reference) {
	if id, ok := ref.Target.(*Identifier); ok {
		if v, exists := cg.variables[id.Name]; exists {
			// lea (load effective address) gets the address
			cg.textSection.WriteString(fmt.Sprintf("    leaq -%d(%%rbp), %%rax\n", v.Offset))
		}
	}
}

// generateDereference generates assembly for dereferencing a pointer (*ptr)
func (cg *CodeGenerator) generateDereference(deref *Dereference) {
	// Evaluate pointer into rax
	cg.generateExpressionToReg(deref.Pointer, "rax")
	// Dereference: load value from address in rax
	cg.textSection.WriteString("    movq (%rax), %rax\n")
}
