package main

// optimizer.go - AST-level optimization passes
// This file contains optimization passes that run on the AST before code generation.
// Optimizations include:
// - Constant folding: Evaluate constant expressions at compile time
// - Strength reduction: Replace expensive operations with cheaper equivalents
// - Identity removal: Remove no-op operations like x + 0, x * 1
// - Peephole optimizations: Local optimizations on small patterns

import (
	"math/bits"
)

// OptimizeAST applies optimization passes to an AST node tree.
// Returns an optimized version of the AST.
func OptimizeAST(statements []ASTNode) []ASTNode {
	result := make([]ASTNode, 0, len(statements))
	for _, stmt := range statements {
		optimized := optimizeNode(stmt)
		if optimized != nil {
			result = append(result, optimized)
		}
	}
	return result
}

// optimizeNode applies optimizations to a single AST node recursively.
func optimizeNode(node ASTNode) ASTNode {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *VariableDeclaration:
		n.Value = optimizeExpression(n.Value)
		return n
	case *ConstantDeclaration:
		n.Value = optimizeExpression(n.Value)
		return n
	case *Assignment:
		n.Value = optimizeExpression(n.Value)
		return n
	case *CompoundAssignment:
		n.Value = optimizeExpression(n.Value)
		return n
	case *ReturnStatement:
		n.Value = optimizeExpression(n.Value)
		return n
	case *FunctionCall:
		for i, arg := range n.Args {
			n.Args[i] = optimizeExpression(arg)
		}
		return n
	case *IfStatement:
		n.Condition = optimizeExpression(n.Condition)
		// Optimize then/else bodies
		for i, stmt := range n.ThenBody {
			n.ThenBody[i] = optimizeNode(stmt)
		}
		for i, stmt := range n.ElseBody {
			n.ElseBody[i] = optimizeNode(stmt)
		}
		// Constant condition elimination
		if cond := evaluateConstantBool(n.Condition); cond != nil {
			if *cond {
				// Condition is always true, inline then block
				// For now, just keep the if statement as-is for simplicity
			}
		}
		return n
	case *WhileLoop:
		n.Condition = optimizeExpression(n.Condition)
		for i, stmt := range n.Body {
			n.Body[i] = optimizeNode(stmt)
		}
		return n
	case *ForLoop:
		if n.Init != nil {
			n.Init = optimizeNode(n.Init)
		}
		n.Condition = optimizeExpression(n.Condition)
		if n.Update != nil {
			n.Update = optimizeNode(n.Update)
		}
		for i, stmt := range n.Body {
			n.Body[i] = optimizeNode(stmt)
		}
		return n
	case *FunctionDefinition:
		for i, stmt := range n.Body {
			n.Body[i] = optimizeNode(stmt)
		}
		return n
	case *TryStatement:
		for i, stmt := range n.TryBlock {
			n.TryBlock[i] = optimizeNode(stmt)
		}
		for _, clause := range n.CatchClauses {
			for i, stmt := range clause.Body {
				clause.Body[i] = optimizeNode(stmt)
			}
		}
		for i, stmt := range n.FinallyBlock {
			n.FinallyBlock[i] = optimizeNode(stmt)
		}
		return n
	default:
		return node
	}
}

// optimizeExpression applies expression-level optimizations.
// This includes constant folding, strength reduction, and identity removal.
func optimizeExpression(expr ASTNode) ASTNode {
	if expr == nil {
		return nil
	}

	switch e := expr.(type) {
	case *BinaryOp:
		// Recursively optimize operands first
		e.Left = optimizeExpression(e.Left)
		e.Right = optimizeExpression(e.Right)

		// Try constant folding
		if folded := foldBinaryOp(e); folded != nil {
			return folded
		}

		// Try strength reduction and identity removal
		if reduced := reduceBinaryOp(e); reduced != nil {
			return reduced
		}

		return e

	case *UnaryOp:
		e.Operand = optimizeExpression(e.Operand)

		// Try constant folding for unary ops
		if folded := foldUnaryOp(e); folded != nil {
			return folded
		}

		return e

	case *BitwiseOp:
		e.Left = optimizeExpression(e.Left)
		e.Right = optimizeExpression(e.Right)

		// Try constant folding
		if folded := foldBitwiseOp(e); folded != nil {
			return folded
		}

		// Try strength reduction for shifts
		if reduced := reduceBitwiseOp(e); reduced != nil {
			return reduced
		}

		return e

	case *FunctionCall:
		for i, arg := range e.Args {
			e.Args[i] = optimizeExpression(arg)
		}
		return e

	case *ArrayAccess:
		e.Index = optimizeExpression(e.Index)
		return e

	case *Comparison:
		e.Left = optimizeExpression(e.Left)
		e.Right = optimizeExpression(e.Right)

		// Try constant folding for comparisons
		if folded := foldComparison(e); folded != nil {
			return folded
		}
		return e

	case *LogicalOp:
		e.Left = optimizeExpression(e.Left)
		e.Right = optimizeExpression(e.Right)

		// Try constant folding for logical operations
		if folded := foldLogicalOp(e); folded != nil {
			return folded
		}
		return e

	default:
		return expr
	}
}

// foldBinaryOp attempts to evaluate a binary operation at compile time.
// Returns nil if folding is not possible (non-constant operands).
func foldBinaryOp(op *BinaryOp) ASTNode {
	left, leftOk := op.Left.(*IntLiteral)
	right, rightOk := op.Right.(*IntLiteral)

	if !leftOk || !rightOk {
		return nil
	}

	var result int
	switch op.Operator {
	case TokenPlus:
		result = left.Value + right.Value
	case TokenMinus:
		result = left.Value - right.Value
	case TokenStar:
		result = left.Value * right.Value
	case TokenSlash:
		if right.Value == 0 {
			return nil // Avoid division by zero
		}
		result = left.Value / right.Value
	case TokenPercent:
		if right.Value == 0 {
			return nil // Avoid division by zero
		}
		result = left.Value % right.Value
	default:
		return nil
	}

	return &IntLiteral{Value: result}
}

// foldUnaryOp attempts to evaluate a unary operation at compile time.
func foldUnaryOp(op *UnaryOp) ASTNode {
	switch operand := op.Operand.(type) {
	case *IntLiteral:
		switch op.Operator {
		case TokenMinus:
			return &IntLiteral{Value: -operand.Value}
		case TokenTilde:
			return &IntLiteral{Value: ^operand.Value}
		}
	case *BoolLiteral:
		if op.Operator == TokenExclaim {
			return &BoolLiteral{Value: !operand.Value}
		}
	}
	return nil
}

// foldBitwiseOp attempts to evaluate a bitwise operation at compile time.
func foldBitwiseOp(op *BitwiseOp) ASTNode {
	left, leftOk := op.Left.(*IntLiteral)
	right, rightOk := op.Right.(*IntLiteral)

	if !leftOk || !rightOk {
		return nil
	}

	var result int
	switch op.Operator {
	case TokenAmpersand:
		result = left.Value & right.Value
	case TokenPipe:
		result = left.Value | right.Value
	case TokenCaret:
		result = left.Value ^ right.Value
	case TokenLShift:
		if right.Value < 0 || right.Value >= 64 {
			return nil // Invalid shift amount
		}
		result = left.Value << uint(right.Value)
	case TokenRShift:
		if right.Value < 0 || right.Value >= 64 {
			return nil // Invalid shift amount
		}
		result = left.Value >> uint(right.Value)
	default:
		return nil
	}

	return &IntLiteral{Value: result}
}

// reduceBinaryOp applies strength reduction and identity removal.
// Examples:
// - x + 0 → x
// - x * 1 → x
// - x * 0 → 0
// - x * 2 → x << 1 (power of 2 multiplication)
// - x / 1 → x
func reduceBinaryOp(op *BinaryOp) ASTNode {
	// Check for constant on the right
	if right, ok := op.Right.(*IntLiteral); ok {
		switch op.Operator {
		case TokenPlus:
			if right.Value == 0 {
				return op.Left // x + 0 = x
			}
		case TokenMinus:
			if right.Value == 0 {
				return op.Left // x - 0 = x
			}
		case TokenStar:
			if right.Value == 0 {
				return &IntLiteral{Value: 0} // x * 0 = 0
			}
			if right.Value == 1 {
				return op.Left // x * 1 = x
			}
			// Strength reduction: multiply by power of 2 → shift
			if right.Value > 0 && isPowerOfTwo(right.Value) {
				shift := bits.TrailingZeros64(uint64(right.Value))
				return &BitwiseOp{
					Left:     op.Left,
					Operator: TokenLShift,
					Right:    &IntLiteral{Value: shift},
				}
			}
		case TokenSlash:
			if right.Value == 1 {
				return op.Left // x / 1 = x
			}
			// Strength reduction: divide by power of 2 → shift (for positive values)
			// Note: This is only valid for unsigned division; for signed we need extra care
			// For now, we'll skip this optimization for safety
		case TokenPercent:
			if right.Value == 1 {
				return &IntLiteral{Value: 0} // x % 1 = 0
			}
		}
	}

	// Check for constant on the left
	if left, ok := op.Left.(*IntLiteral); ok {
		switch op.Operator {
		case TokenPlus:
			if left.Value == 0 {
				return op.Right // 0 + x = x
			}
		case TokenStar:
			if left.Value == 0 {
				return &IntLiteral{Value: 0} // 0 * x = 0
			}
			if left.Value == 1 {
				return op.Right // 1 * x = x
			}
			// Strength reduction for power of 2
			if left.Value > 0 && isPowerOfTwo(left.Value) {
				shift := bits.TrailingZeros64(uint64(left.Value))
				return &BitwiseOp{
					Left:     op.Right,
					Operator: TokenLShift,
					Right:    &IntLiteral{Value: shift},
				}
			}
		}
	}

	return nil
}

// reduceBitwiseOp applies identity removal for bitwise operations.
// Examples:
// - x & 0 → 0
// - x | 0 → x
// - x ^ 0 → x
// - x << 0 → x
// - x >> 0 → x
func reduceBitwiseOp(op *BitwiseOp) ASTNode {
	if right, ok := op.Right.(*IntLiteral); ok {
		switch op.Operator {
		case TokenAmpersand:
			if right.Value == 0 {
				return &IntLiteral{Value: 0} // x & 0 = 0
			}
			if right.Value == -1 { // All bits set
				return op.Left // x & -1 = x
			}
		case TokenPipe:
			if right.Value == 0 {
				return op.Left // x | 0 = x
			}
			if right.Value == -1 {
				return &IntLiteral{Value: -1} // x | -1 = -1
			}
		case TokenCaret:
			if right.Value == 0 {
				return op.Left // x ^ 0 = x
			}
		case TokenLShift, TokenRShift:
			if right.Value == 0 {
				return op.Left // x << 0 = x, x >> 0 = x
			}
		}
	}

	return nil
}

// isPowerOfTwo checks if n is a power of two.
func isPowerOfTwo(n int) bool {
	if n <= 0 {
		return false
	}
	return n&(n-1) == 0
}

// foldComparison attempts to evaluate a comparison at compile time.
// Returns nil if folding is not possible.
func foldComparison(cmp *Comparison) ASTNode {
	left, leftOk := cmp.Left.(*IntLiteral)
	right, rightOk := cmp.Right.(*IntLiteral)

	if !leftOk || !rightOk {
		return nil
	}

	var result bool
	switch cmp.Operator {
	case TokenEqual:
		result = left.Value == right.Value
	case TokenNotEqual:
		result = left.Value != right.Value
	case TokenLess:
		result = left.Value < right.Value
	case TokenLessEq:
		result = left.Value <= right.Value
	case TokenGreater:
		result = left.Value > right.Value
	case TokenGreaterEq:
		result = left.Value >= right.Value
	default:
		return nil
	}

	return &BoolLiteral{Value: result}
}

// foldLogicalOp attempts to evaluate a logical operation at compile time.
// Also applies short-circuit optimizations.
func foldLogicalOp(op *LogicalOp) ASTNode {
	// Check for constant left operand (short-circuit optimization)
	if leftBool := evaluateConstantBool(op.Left); leftBool != nil {
		switch op.Operator {
		case TokenAnd:
			if !*leftBool {
				return &BoolLiteral{Value: false} // false && x = false
			}
			// true && x = x
			return op.Right
		case TokenOr:
			if *leftBool {
				return &BoolLiteral{Value: true} // true || x = true
			}
			// false || x = x
			return op.Right
		}
	}

	// Check for constant right operand
	if rightBool := evaluateConstantBool(op.Right); rightBool != nil {
		switch op.Operator {
		case TokenAnd:
			if !*rightBool {
				return &BoolLiteral{Value: false} // x && false = false
			}
			// x && true = x
			return op.Left
		case TokenOr:
			if *rightBool {
				return &BoolLiteral{Value: true} // x || true = true
			}
			// x || false = x
			return op.Left
		}
	}

	return nil
}

// evaluateConstantBool attempts to evaluate an expression to a constant boolean.
// Returns nil if the expression is not a constant.
func evaluateConstantBool(expr ASTNode) *bool {
	switch e := expr.(type) {
	case *BoolLiteral:
		return &e.Value
	case *IntLiteral:
		result := e.Value != 0
		return &result
	}
	return nil
}

// OptimizationStats tracks statistics about applied optimizations.
type OptimizationStats struct {
	ConstantsFolded   int
	StrengthReduced   int
	IdentitiesRemoved int
}

// String returns a human-readable summary of optimization statistics.
func (s *OptimizationStats) String() string {
	if s.ConstantsFolded == 0 && s.StrengthReduced == 0 && s.IdentitiesRemoved == 0 {
		return "No optimizations applied"
	}
	return ""
}
