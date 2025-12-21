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

// Comparison represents a comparison expression
type Comparison struct {
	Left     ASTNode
	Operator TokenType // TokenEqual, TokenNotEqual, TokenLess, TokenLessEq, TokenGreater, TokenGreaterEq
	Right    ASTNode
}

func (c *Comparison) astNode() {}

// LogicalOp represents logical operations (&&, ||)
type LogicalOp struct {
	Left     ASTNode
	Operator TokenType // TokenAnd, TokenOr
	Right    ASTNode
}

func (l *LogicalOp) astNode() {}

// TernaryOp represents ternary conditional operator (condition ? trueExpr : falseExpr)
type TernaryOp struct {
	Condition ASTNode
	TrueExpr  ASTNode
	FalseExpr ASTNode
}

func (t *TernaryOp) astNode() {}

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
	case TokenEqual:
		cg.textSection.WriteString("    sete %al\n")
	case TokenNotEqual:
		cg.textSection.WriteString("    setne %al\n")
	case TokenLess:
		cg.textSection.WriteString("    setl %al\n")
	case TokenLessEq:
		cg.textSection.WriteString("    setle %al\n")
	case TokenGreater:
		cg.textSection.WriteString("    setg %al\n")
	case TokenGreaterEq:
		cg.textSection.WriteString("    setge %al\n")
	}

	// Zero-extend result to 64 bits
	cg.textSection.WriteString("    movzbl %al, %eax\n")
}

// generateLogicalOp generates assembly for logical operations (&&, ||)
func (cg *CodeGenerator) generateLogicalOp(logOp *LogicalOp) {
	if logOp.Operator == TokenAnd {
		// Short-circuit AND: if left is false, result is false
		endLabel := cg.getLabel("and_end")
		
		// Evaluate left
		cg.generateExpressionToReg(logOp.Left, "rax")
		cg.textSection.WriteString("    testq %rax, %rax\n")
		cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", endLabel))
		
		// Left is true, evaluate right
		cg.generateExpressionToReg(logOp.Right, "rax")
		
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLabel))
	} else if logOp.Operator == TokenOr {
		// Short-circuit OR: if left is true, result is true
		trueLabel := cg.getLabel("or_true")
		endLabel := cg.getLabel("or_end")
		
		// Evaluate left
		cg.generateExpressionToReg(logOp.Left, "rax")
		cg.textSection.WriteString("    testq %rax, %rax\n")
		cg.textSection.WriteString(fmt.Sprintf("    jnz %s\n", trueLabel))
		
		// Left is false, evaluate right
		cg.generateExpressionToReg(logOp.Right, "rax")
		cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", endLabel))
		
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", trueLabel))
		cg.textSection.WriteString("    movq $1, %rax\n")
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLabel))
	}
}

// generateTernaryOp generates assembly for ternary conditional operator
func (cg *CodeGenerator) generateTernaryOp(ternary *TernaryOp) {
	falseLabel := cg.getLabel("ternary_false")
	endLabel := cg.getLabel("ternary_end")
	
	// Evaluate condition
	cg.generateExpressionToReg(ternary.Condition, "rax")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", falseLabel))
	
	// True branch
	cg.generateExpressionToReg(ternary.TrueExpr, "rax")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", endLabel))
	
	// False branch
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", falseLabel))
	cg.generateExpressionToReg(ternary.FalseExpr, "rax")
	
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLabel))
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
