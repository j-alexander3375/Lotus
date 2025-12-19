package main

import "fmt"

// IfStatement represents an if/else control flow structure
type IfStatement struct {
	Condition ASTNode
	ThenBody  []ASTNode
	ElseBody  []ASTNode
}

func (i *IfStatement) astNode() {}

// WhileLoop represents a while loop structure
type WhileLoop struct {
	Condition ASTNode
	Body      []ASTNode
}

func (w *WhileLoop) astNode() {}

// ForLoop represents a for loop structure
type ForLoop struct {
	Init      ASTNode
	Condition ASTNode
	Update    ASTNode
	Body      []ASTNode
}

func (f *ForLoop) astNode() {}

// ComparisonOp represents comparison operators
type ComparisonOp string

const (
	CmpEqual       ComparisonOp = "=="
	CmpNotEqual    ComparisonOp = "!="
	CmpLess        ComparisonOp = "<"
	CmpLessEq      ComparisonOp = "<="
	CmpGreater     ComparisonOp = ">"
	CmpGreaterEq   ComparisonOp = ">="
)

// Comparison represents a comparison expression
type Comparison struct {
	Left     ASTNode
	Operator ComparisonOp
	Right    ASTNode
}

func (c *Comparison) astNode() {}

// generateIfStatement generates assembly for if/else statements
func (cg *CodeGenerator) generateIfStatement(ifStmt *IfStatement) {
	ifLabel := cg.getLabel("if_then")
	endLabel := cg.getLabel("if_end")
	elseLabel := ""
	if len(ifStmt.ElseBody) > 0 {
		elseLabel = cg.getLabel("if_else")
	}

	// Evaluate condition into rax
	cg.generateConditionToReg(ifStmt.Condition, "rax")

	// Jump to else or end if condition is false
	if len(ifStmt.ElseBody) > 0 {
		cg.textSection.WriteString(fmt.Sprintf("    testq %%rax, %%rax\n    jz %s\n", elseLabel))
	} else {
		cg.textSection.WriteString(fmt.Sprintf("    testq %%rax, %%rax\n    jz %s\n", endLabel))
	}

	// Then body
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", ifLabel))
	for _, stmt := range ifStmt.ThenBody {
		cg.generateStatement(stmt)
	}
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", endLabel))

	// Else body
	if len(ifStmt.ElseBody) > 0 {
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", elseLabel))
		for _, stmt := range ifStmt.ElseBody {
			cg.generateStatement(stmt)
		}
	}

	// End label
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLabel))
}

// generateWhileLoop generates assembly for while loops
func (cg *CodeGenerator) generateWhileLoop(whileLoop *WhileLoop) {
	loopLabel := cg.getLabel("while_loop")
	endLabel := cg.getLabel("while_end")

	// Loop start
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", loopLabel))

	// Evaluate condition
	cg.generateConditionToReg(whileLoop.Condition, "rax")

	// Exit loop if condition is false
	cg.textSection.WriteString(fmt.Sprintf("    testq %%rax, %%rax\n    jz %s\n", endLabel))

	// Loop body
	for _, stmt := range whileLoop.Body {
		cg.generateStatement(stmt)
	}

	// Jump back to condition
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", loopLabel))

	// End label
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLabel))
}

// generateForLoop generates assembly for for loops
func (cg *CodeGenerator) generateForLoop(forLoop *ForLoop) {
	loopLabel := cg.getLabel("for_loop")
	endLabel := cg.getLabel("for_end")

	// Initialize
	if forLoop.Init != nil {
		cg.generateStatement(forLoop.Init)
	}

	// Loop start
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", loopLabel))

	// Evaluate condition
	if forLoop.Condition != nil {
		cg.generateConditionToReg(forLoop.Condition, "rax")
		cg.textSection.WriteString(fmt.Sprintf("    testq %%rax, %%rax\n    jz %s\n", endLabel))
	}

	// Loop body
	for _, stmt := range forLoop.Body {
		cg.generateStatement(stmt)
	}

	// Update
	if forLoop.Update != nil {
		cg.generateStatement(forLoop.Update)
	}

	// Jump back to condition
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", loopLabel))

	// End label
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLabel))
}

// generateComparison generates assembly for comparison operations
func (cg *CodeGenerator) generateComparison(cmp *Comparison) {
	// Evaluate left into rax
	cg.generateExpressionToReg(cmp.Left, "rax")

	// Save to stack
	cg.textSection.WriteString("    pushq %rax\n")

	// Evaluate right into rax
	cg.generateExpressionToReg(cmp.Right, "rax")

	// Move left back to rcx
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString("    popq %rax\n")

	// Compare
	cg.textSection.WriteString("    cmpq %rcx, %rax\n")

	// Set result based on condition
	switch cmp.Operator {
	case CmpEqual:
		cg.textSection.WriteString("    sete %al\n")
	case CmpNotEqual:
		cg.textSection.WriteString("    setne %al\n")
	case CmpLess:
		cg.textSection.WriteString("    setl %al\n")
	case CmpLessEq:
		cg.textSection.WriteString("    setle %al\n")
	case CmpGreater:
		cg.textSection.WriteString("    setg %al\n")
	case CmpGreaterEq:
		cg.textSection.WriteString("    setge %al\n")
	}

	// Zero-extend result to 64 bits
	cg.textSection.WriteString("    movzbl %al, %eax\n")
}

// generateConditionToReg evaluates a condition and stores result in register
func (cg *CodeGenerator) generateConditionToReg(cond ASTNode, reg string) {
	switch c := cond.(type) {
	case *Comparison:
		cg.generateComparison(c)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	case *BinaryOp:
		cg.generateBinaryOp(c)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	default:
		cg.generateExpressionToReg(cond, reg)
	}
}

// Label management
var labelCounter int = 0

func (cg *CodeGenerator) getLabel(prefix string) string {
	label := fmt.Sprintf(".%s_%d", prefix, labelCounter)
	labelCounter++
	return label
}
