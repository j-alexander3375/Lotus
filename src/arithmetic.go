package main

import "fmt"

// BinaryOp represents a binary operation expression
type BinaryOp struct {
	Left     ASTNode
	Operator TokenType // TokenPlus, TokenMinus, TokenStar, TokenSlash, TokenPercent
	Right    ASTNode
}

func (b *BinaryOp) astNode() {}

// UnaryOp represents a unary operation expression
type UnaryOp struct {
	Operator TokenType // TokenMinus, TokenExclaim, TokenAmpersand, TokenStar
	Operand  ASTNode
}

func (u *UnaryOp) astNode() {}

// BitwiseOp represents bitwise operations (&, |, ^, <<, >>)
type BitwiseOp struct {
	Left     ASTNode
	Operator TokenType // TokenAmpersand, TokenPipe, TokenCaret, TokenLShift, TokenRShift
	Right    ASTNode
}

func (b *BitwiseOp) astNode() {}

// IncrementOp represents increment/decrement (++, --)
type IncrementOp struct {
	Operand  ASTNode
	IsPrefix bool      // true for ++x, false for x++
	Operator TokenType // TokenPlusPlus, TokenMinusMinus
}

func (i *IncrementOp) astNode() {}

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
	case TokenPlus:
		cg.textSection.WriteString("    addq %rcx, %rax\n")
	case TokenMinus:
		cg.textSection.WriteString("    subq %rcx, %rax\n")
	case TokenStar:
		cg.textSection.WriteString("    imulq %rcx, %rax\n")
	case TokenSlash:
		cg.textSection.WriteString("    cqo\n")        // sign extend rax to rdx:rax
		cg.textSection.WriteString("    idivq %rcx\n") // divide rax by rcx
	case TokenPercent:
		cg.textSection.WriteString("    cqo\n")             // sign extend rax to rdx:rax
		cg.textSection.WriteString("    idivq %rcx\n")      // divide rax by rcx
		cg.textSection.WriteString("    movq %rdx, %rax\n") // move remainder to rax
	default:
		fmt.Printf("Unknown binary operator: %v\n", binop.Operator)
	}
}

// generateUnaryOp generates assembly for unary operations
func (cg *CodeGenerator) generateUnaryOp(unop *UnaryOp) {
	cg.generateExpressionToReg(unop.Operand, "rax")

	switch unop.Operator {
	case TokenMinus:
		cg.textSection.WriteString("    negq %rax\n")
	case TokenExclaim:
		cg.textSection.WriteString("    testq %rax, %rax\n")
		cg.textSection.WriteString("    setzb %al\n")
		cg.textSection.WriteString("    movzbl %al, %eax\n")
	case TokenAmpersand:
		// Unary & (address-of) - handled in references.go
	case TokenStar:
		// Unary * (dereference) - handled in references.go
	case TokenTilde:
		// Bitwise NOT
		cg.textSection.WriteString("    notq %rax\n")
	}
}

// generateBitwiseOp generates assembly for bitwise operations
func (cg *CodeGenerator) generateBitwiseOp(bitOp *BitwiseOp) {
	// Evaluate left operand into rax
	cg.generateExpressionToReg(bitOp.Left, "rax")

	// Save left result on stack temporarily
	cg.textSection.WriteString("    pushq %rax\n")

	// Evaluate right operand into rax
	cg.generateExpressionToReg(bitOp.Right, "rax")

	// Pop left operand from stack to rcx
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString("    popq %rax\n")

	// Perform operation
	switch bitOp.Operator {
	case TokenAmpersand:
		cg.textSection.WriteString("    andq %rcx, %rax\n")
	case TokenPipe:
		cg.textSection.WriteString("    orq %rcx, %rax\n")
	case TokenCaret:
		cg.textSection.WriteString("    xorq %rcx, %rax\n")
	case TokenLShift:
		// Shift left: rax << rcx
		cg.textSection.WriteString("    shlq %cl, %rax\n")
	case TokenRShift:
		// Shift right: rax >> rcx
		cg.textSection.WriteString("    shrq %cl, %rax\n")
	default:
		fmt.Printf("Unknown bitwise operator: %v\n", bitOp.Operator)
	}
}

// generateIncrementOp generates assembly for increment/decrement operations
func (cg *CodeGenerator) generateIncrementOp(incOp *IncrementOp) {
	// Get variable location
	if ident, ok := incOp.Operand.(*Identifier); ok {
		if v, exists := cg.variables[ident.Name]; exists {
			if incOp.IsPrefix {
				// Prefix: increment then use value
				if incOp.Operator == TokenPlusPlus {
					cg.textSection.WriteString(fmt.Sprintf("    incq -%d(%%rbp)\n", v.Offset))
				} else {
					cg.textSection.WriteString(fmt.Sprintf("    decq -%d(%%rbp)\n", v.Offset))
				}
				cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rax\n", v.Offset))
			} else {
				// Postfix: use value then increment
				cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rax\n", v.Offset))
				if incOp.Operator == TokenPlusPlus {
					cg.textSection.WriteString(fmt.Sprintf("    incq -%d(%%rbp)\n", v.Offset))
				} else {
					cg.textSection.WriteString(fmt.Sprintf("    decq -%d(%%rbp)\n", v.Offset))
				}
			}
		}
	}
}

// generateExpressionToReg generates code to evaluate an expression and store result in the specified register
func (cg *CodeGenerator) generateExpressionToReg(expr ASTNode, reg string) {
	switch e := expr.(type) {
	case *IntLiteral:
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%%s\n", e.Value, reg))
	case *Identifier:
		// Check if it's a variable first
		if v, exists := cg.variables[e.Name]; exists {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%%s\n", v.Offset, reg))
		} else if c, exists := cg.constants[e.Name]; exists {
			// It's a constant - load from data section
			label := fmt.Sprintf(".const_%s", c.Name)
			if c.Type == TokenTypeString {
				// String constants: load address
				cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%%s\n", label, reg))
			} else {
				// Numeric/bool constants: load value
				cg.textSection.WriteString(fmt.Sprintf("    movq %s(%%rip), %%%s\n", label, reg))
			}
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
	case *ArrayAccess:
		cg.generateArrayAccess(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	case *ArrayLiteral:
		cg.generateArrayLiteral(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	case *FieldAccess:
		cg.generateFieldAccess(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	case *StructLiteral:
		cg.generateStructLiteral(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	case *ClassLiteral:
		cg.generateClassLiteral(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	case *EnumLiteral:
		cg.generateEnumLiteral(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	case *MallocCall:
		cg.generateMallocCall(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	case *SizeofExpr:
		cg.generateSizeofExpr(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	case *MethodCall:
		cg.generateMethodCall(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	case *BitwiseOp:
		cg.generateBitwiseOp(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	case *LogicalOp:
		cg.generateLogicalOp(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	case *TernaryOp:
		cg.generateTernaryOp(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	case *IncrementOp:
		cg.generateIncrementOp(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	case *NullLiteral:
		cg.generateNullLiteral(e)
		if reg != "rax" {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, %%%s\n", reg))
		}
	}
}
