package main

import "fmt"

// ArithmeticOperator represents binary arithmetic operations
type ArithmeticOperator int

const (
	OpAdd ArithmeticOperator = iota
	OpSubtract
	OpMultiply
	OpDivide
	OpModulo
)

// BinaryOp represents a binary operation expression
type BinaryOp struct {
	Left     ASTNode
	Operator ArithmeticOperator
	Right    ASTNode
}

func (b *BinaryOp) astNode() {}

// UnaryOp represents a unary operation expression
type UnaryOp struct {
	Operator string // "+", "-", "!"
	Operand  ASTNode
}

func (u *UnaryOp) astNode() {}

// generateBinaryOp generates assembly for binary operations
func (cg *CodeGenerator) generateBinaryOp(binop *BinaryOp) {
	// Evaluate left operand into rax
	cg.generateExpressionToReg(binop.Left, "rax")

	// Save left result on stack temporarily
	cg.textSection.WriteString("    pushq %rax\n")

	// Evaluate right operand into rax
	cg.generateExpressionToReg(binop.Right, "rax")

	// Pop left operand from stack to rcx
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString("    popq %rax\n")

	// Perform operation
	switch binop.Operator {
	case OpAdd:
		cg.textSection.WriteString("    addq %rcx, %rax\n")
	case OpSubtract:
		cg.textSection.WriteString("    subq %rcx, %rax\n")
	case OpMultiply:
		cg.textSection.WriteString("    imulq %rcx, %rax\n")
	case OpDivide:
		cg.textSection.WriteString("    cqo\n")
		cg.textSection.WriteString("    idivq %rcx\n")
	case OpModulo:
		cg.textSection.WriteString("    cqo\n")
		cg.textSection.WriteString("    idivq %rcx\n")
		cg.textSection.WriteString("    movq %rdx, %rax\n") // remainder in rdx
	}
}

// generateUnaryOp generates assembly for unary operations
func (cg *CodeGenerator) generateUnaryOp(unop *UnaryOp) {
	cg.generateExpressionToReg(unop.Operand, "rax")

	switch unop.Operator {
	case "-":
		cg.textSection.WriteString("    negq %rax\n")
	case "+":
		// No-op for unary plus
	case "!":
		cg.textSection.WriteString("    testq %rax, %rax\n")
		cg.textSection.WriteString("    setzb %al\n")
		cg.textSection.WriteString("    movzbl %al, %eax\n")
	}
}

// generateExpressionToReg generates code to evaluate an expression and store result in the specified register
func (cg *CodeGenerator) generateExpressionToReg(expr ASTNode, reg string) {
	switch e := expr.(type) {
	case *IntLiteral:
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%%s\n", e.Value, reg))
	case *Identifier:
		if v, exists := cg.variables[e.Name]; exists {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%%s\n", v.Offset, reg))
		}
	case *BinaryOp:
		cg.generateBinaryOp(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	case *UnaryOp:
		cg.generateUnaryOp(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	}
}
