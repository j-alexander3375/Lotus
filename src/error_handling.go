package main

import (
	"fmt"
)

// error_handling.go - Exception handling and null safety
// This file implements try/catch/finally blocks and null literal support.

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

// generateTryStatement generates code for try-catch-finally blocks.
// Uses label-based control flow for exception handling.
func (cg *CodeGenerator) generateTryStatement(stmt *TryStatement) {
	tryLabel := fmt.Sprintf(".try_%d", cg.labelCount)
	catchLabel := fmt.Sprintf(".catch_%d", cg.labelCount)
	finallyLabel := fmt.Sprintf(".finally_%d", cg.labelCount)
	endLabel := fmt.Sprintf(".try_end_%d", cg.labelCount)
	cg.labelCount++

	cg.textSection.WriteString("    # Try-catch-finally block\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", tryLabel))

	// Generate try block
	for _, node := range stmt.TryBlock {
		cg.generateStatement(node)
	}

	// If no exception, jump to finally (or end if no finally)
	if len(stmt.FinallyBlock) > 0 {
		cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", finallyLabel))
	} else {
		cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", endLabel))
	}

	// Generate catch blocks
	if len(stmt.CatchClauses) > 0 {
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", catchLabel))
		for i, catchClause := range stmt.CatchClauses {
			if catchClause.ExceptionVar != "" {
				// Store exception in variable
				cg.stackOffset += 8
				cg.variables[catchClause.ExceptionVar] = Variable{
					Name:   catchClause.ExceptionVar,
					Type:   TokenTypeInt, // Exception code as int
					Offset: cg.stackOffset,
				}
				cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, -%d(%%rbp)  # Store exception: %s\n",
					cg.stackOffset, catchClause.ExceptionVar))
			}

			// Generate catch body
			for _, node := range catchClause.Body {
				cg.generateStatement(node)
			}

			// Jump to finally or end
			if i == len(stmt.CatchClauses)-1 {
				if len(stmt.FinallyBlock) > 0 {
					cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", finallyLabel))
				} else {
					cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", endLabel))
				}
			}
		}
	}

	// Generate finally block
	if len(stmt.FinallyBlock) > 0 {
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", finallyLabel))
		for _, node := range stmt.FinallyBlock {
			cg.generateStatement(node)
		}
	}

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLabel))
}

// generateThrowStatement generates code for throwing an exception
func (cg *CodeGenerator) generateThrowStatement(stmt *ThrowStatement) {
	cg.textSection.WriteString("    # Throw exception\n")

	// Evaluate exception expression to rax
	cg.generateExpressionToReg(stmt.Exception, "rax")

	// For now, just exit with the exception code
	// In a full implementation, this would unwind the stack to the nearest catch
	cg.textSection.WriteString("    # Exception thrown - exiting with error code\n")
	cg.textSection.WriteString("    movq %rax, %rdi\n")
	cg.textSection.WriteString("    movq $60, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// generateNullLiteral generates code for null value
func (cg *CodeGenerator) generateNullLiteral(lit *NullLiteral) {
	cg.textSection.WriteString("    # null literal\n")
	cg.textSection.WriteString("    movq $0, %rax\n")
}

// isNullCheck generates a null check for an expression
func (cg *CodeGenerator) generateNullCheck(expr ASTNode) {
	checkLabel := fmt.Sprintf(".null_check_%d", cg.labelCount)
	okLabel := fmt.Sprintf(".null_ok_%d", cg.labelCount)
	cg.labelCount++

	cg.textSection.WriteString("    # Null check\n")
	cg.generateExpressionToReg(expr, "rax")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s\n", okLabel))

	// Handle null - throw null pointer exception
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", checkLabel))
	cg.textSection.WriteString("    # Null pointer exception\n")
	cg.textSection.WriteString("    movq $1, %rdi  # Exit code 1 for null pointer\n")
	cg.textSection.WriteString("    movq $60, %rax\n")
	cg.textSection.WriteString("    syscall\n")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", okLabel))
}
