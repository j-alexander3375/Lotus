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

// ============================================================================
// DEAD CODE ELIMINATION
// ============================================================================
//
// Dead code elimination removes code that can never be executed or has no effect:
// 1. Code after unconditional return/break/continue statements
// 2. Empty if/else branches
// 3. Unreachable code in if(false) branches
// 4. Unused variable assignments (when variable is never read)

// EliminateDeadCode removes unreachable code from a list of statements.
// It processes each function and removes dead code within function bodies.
func EliminateDeadCode(statements []ASTNode) []ASTNode {
	result := make([]ASTNode, 0, len(statements))
	for _, stmt := range statements {
		cleaned := eliminateDeadCodeInNode(stmt)
		if cleaned != nil {
			result = append(result, cleaned)
		}
	}
	return result
}

// eliminateDeadCodeInNode removes dead code from a single AST node.
func eliminateDeadCodeInNode(node ASTNode) ASTNode {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *FunctionDefinition:
		n.Body = eliminateDeadCodeInBlock(n.Body)
		return n

	case *IfStatement:
		// Check if condition is constant
		if cond := evaluateConstantBool(n.Condition); cond != nil {
			if *cond {
				// if(true) - just use the then body as statements
				// For now, keep the structure but eliminate else
				n.ElseBody = nil
			} else {
				// if(false) - eliminate then body, use else body
				n.ThenBody = n.ElseBody
				n.ElseBody = nil
				n.Condition = &BoolLiteral{Value: true}
			}
		}
		n.ThenBody = eliminateDeadCodeInBlock(n.ThenBody)
		n.ElseBody = eliminateDeadCodeInBlock(n.ElseBody)

		// Remove empty if statements
		if len(n.ThenBody) == 0 && len(n.ElseBody) == 0 {
			return nil
		}
		return n

	case *WhileLoop:
		// Check for while(false) - dead loop
		if cond := evaluateConstantBool(n.Condition); cond != nil && !*cond {
			return nil // while(false) never executes
		}
		n.Body = eliminateDeadCodeInBlock(n.Body)
		return n

	case *ForLoop:
		// Check for for(;false;) - dead loop
		if cond := evaluateConstantBool(n.Condition); cond != nil && !*cond {
			// Only the init statement might have side effects
			if n.Init != nil {
				return n.Init
			}
			return nil
		}
		n.Body = eliminateDeadCodeInBlock(n.Body)
		return n

	case *TryStatement:
		n.TryBlock = eliminateDeadCodeInBlock(n.TryBlock)
		for _, clause := range n.CatchClauses {
			clause.Body = eliminateDeadCodeInBlock(clause.Body)
		}
		n.FinallyBlock = eliminateDeadCodeInBlock(n.FinallyBlock)
		return n

	default:
		return node
	}
}

// eliminateDeadCodeInBlock removes dead code from a block of statements.
// Specifically, it removes any code that appears after a return, break, or continue.
func eliminateDeadCodeInBlock(block []ASTNode) []ASTNode {
	if block == nil {
		return nil
	}

	result := make([]ASTNode, 0, len(block))
	for _, stmt := range block {
		// First, recursively process the statement
		cleaned := eliminateDeadCodeInNode(stmt)
		if cleaned == nil {
			continue
		}
		result = append(result, cleaned)

		// Check if this is a terminating statement
		if isTerminatingStatement(cleaned) {
			// Everything after this is dead code
			break
		}
	}
	return result
}

// isTerminatingStatement returns true if the statement unconditionally
// terminates the current block (return, break, continue, throw).
func isTerminatingStatement(stmt ASTNode) bool {
	switch s := stmt.(type) {
	case *ReturnStatement:
		return true
	case *BreakStatement:
		return true
	case *ContinueStatement:
		return true
	case *ThrowStatement:
		return true
	case *IfStatement:
		// if both branches terminate, the whole if terminates
		if len(s.ThenBody) > 0 && len(s.ElseBody) > 0 {
			thenTerminates := len(s.ThenBody) > 0 && isTerminatingStatement(s.ThenBody[len(s.ThenBody)-1])
			elseTerminates := len(s.ElseBody) > 0 && isTerminatingStatement(s.ElseBody[len(s.ElseBody)-1])
			return thenTerminates && elseTerminates
		}
		return false
	}
	return false
}

// ============================================================================
// UNUSED VARIABLE DETECTION
// ============================================================================

// VariableUsage tracks how a variable is used in the code.
type VariableUsage struct {
	Name       string
	DeclNode   ASTNode // The declaration node
	WriteCount int     // Number of times written to
	ReadCount  int     // Number of times read from
}

// ============================================================================
// UNUSED FUNCTION ELIMINATION
// ============================================================================

// EliminateUnusedFunctions removes functions that are never called.
// It preserves "main" and any functions called directly or indirectly from main.
func EliminateUnusedFunctions(statements []ASTNode) []ASTNode {
	// Collect all function definitions
	functions := make(map[string]*FunctionDefinition)
	var nonFunctions []ASTNode

	for _, stmt := range statements {
		if fn, ok := stmt.(*FunctionDefinition); ok {
			functions[fn.Name] = fn
		} else {
			nonFunctions = append(nonFunctions, stmt)
		}
	}

	// Find all called functions starting from "main"
	calledFunctions := make(map[string]bool)
	if mainFn, ok := functions["main"]; ok {
		markCalledFunctions(mainFn, functions, calledFunctions)
	}

	// Always keep main
	calledFunctions["main"] = true

	// Build result with only used functions
	result := make([]ASTNode, 0, len(statements))
	result = append(result, nonFunctions...)
	for name, fn := range functions {
		if calledFunctions[name] {
			result = append(result, fn)
		}
	}

	return result
}

// markCalledFunctions recursively marks all functions called from the given function.
func markCalledFunctions(fn *FunctionDefinition, allFunctions map[string]*FunctionDefinition, called map[string]bool) {
	if fn == nil || called[fn.Name] {
		return
	}
	called[fn.Name] = true

	// Find all function calls in this function's body
	callsInBody := findFunctionCalls(fn.Body)
	for _, callName := range callsInBody {
		if calledFn, ok := allFunctions[callName]; ok {
			markCalledFunctions(calledFn, allFunctions, called)
		}
	}
}

// findFunctionCalls finds all function names called in a list of statements.
func findFunctionCalls(statements []ASTNode) []string {
	var calls []string
	for _, stmt := range statements {
		calls = append(calls, findCallsInNode(stmt)...)
	}
	return calls
}

// findCallsInNode finds all function calls in a single AST node.
func findCallsInNode(node ASTNode) []string {
	if node == nil {
		return nil
	}

	var calls []string
	switch n := node.(type) {
	case *FunctionCall:
		calls = append(calls, n.Name)
		for _, arg := range n.Args {
			calls = append(calls, findCallsInNode(arg)...)
		}
	case *VariableDeclaration:
		calls = append(calls, findCallsInNode(n.Value)...)
	case *Assignment:
		calls = append(calls, findCallsInNode(n.Value)...)
	case *ReturnStatement:
		calls = append(calls, findCallsInNode(n.Value)...)
	case *IfStatement:
		calls = append(calls, findCallsInNode(n.Condition)...)
		calls = append(calls, findFunctionCalls(n.ThenBody)...)
		calls = append(calls, findFunctionCalls(n.ElseBody)...)
	case *WhileLoop:
		calls = append(calls, findCallsInNode(n.Condition)...)
		calls = append(calls, findFunctionCalls(n.Body)...)
	case *ForLoop:
		if n.Init != nil {
			calls = append(calls, findCallsInNode(n.Init)...)
		}
		calls = append(calls, findCallsInNode(n.Condition)...)
		if n.Update != nil {
			calls = append(calls, findCallsInNode(n.Update)...)
		}
		calls = append(calls, findFunctionCalls(n.Body)...)
	case *BinaryOp:
		calls = append(calls, findCallsInNode(n.Left)...)
		calls = append(calls, findCallsInNode(n.Right)...)
	case *UnaryOp:
		calls = append(calls, findCallsInNode(n.Operand)...)
	case *Comparison:
		calls = append(calls, findCallsInNode(n.Left)...)
		calls = append(calls, findCallsInNode(n.Right)...)
	case *LogicalOp:
		calls = append(calls, findCallsInNode(n.Left)...)
		calls = append(calls, findCallsInNode(n.Right)...)
	case *TernaryOp:
		calls = append(calls, findCallsInNode(n.Condition)...)
		calls = append(calls, findCallsInNode(n.TrueExpr)...)
		calls = append(calls, findCallsInNode(n.FalseExpr)...)
	case *PipeExpression:
		calls = append(calls, findCallsInNode(n.Left)...)
		calls = append(calls, n.Function) // The function being piped into
	case *MatchExpression:
		calls = append(calls, findCallsInNode(n.Value)...)
		for _, c := range n.Cases {
			calls = append(calls, findFunctionCalls(c.Body)...)
		}
	}
	return calls
}

// FindUnusedVariables analyzes statements and returns a list of unused variable names.
func FindUnusedVariables(statements []ASTNode) []string {
	usage := make(map[string]*VariableUsage)

	// First pass: collect all variable declarations
	collectVariableDeclarations(statements, usage)

	// Second pass: count reads and writes
	countVariableUsage(statements, usage)

	// Find unused variables (declared but never read)
	var unused []string
	for name, u := range usage {
		if u.ReadCount == 0 {
			unused = append(unused, name)
		}
	}
	return unused
}

// collectVariableDeclarations finds all variable declarations in the statements.
func collectVariableDeclarations(statements []ASTNode, usage map[string]*VariableUsage) {
	for _, stmt := range statements {
		collectDeclsInNode(stmt, usage)
	}
}

func collectDeclsInNode(node ASTNode, usage map[string]*VariableUsage) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *VariableDeclaration:
		usage[n.Name] = &VariableUsage{Name: n.Name, DeclNode: n, WriteCount: 1}
	case *FunctionDefinition:
		// Create a new scope for function parameters
		for _, param := range n.Parameters {
			usage[param.Name] = &VariableUsage{Name: param.Name, DeclNode: nil, WriteCount: 1}
		}
		collectVariableDeclarations(n.Body, usage)
	case *IfStatement:
		collectVariableDeclarations(n.ThenBody, usage)
		collectVariableDeclarations(n.ElseBody, usage)
	case *WhileLoop:
		collectVariableDeclarations(n.Body, usage)
	case *ForLoop:
		if n.Init != nil {
			collectDeclsInNode(n.Init, usage)
		}
		collectVariableDeclarations(n.Body, usage)
	}
}

// countVariableUsage counts reads and writes to variables.
func countVariableUsage(statements []ASTNode, usage map[string]*VariableUsage) {
	for _, stmt := range statements {
		countUsageInNode(stmt, usage)
	}
}

func countUsageInNode(node ASTNode, usage map[string]*VariableUsage) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *Identifier:
		if u, ok := usage[n.Name]; ok {
			u.ReadCount++
		}
	case *Assignment:
		if id, ok := n.Target.(*Identifier); ok {
			if u, exists := usage[id.Name]; exists {
				u.WriteCount++
			}
		}
		countUsageInExpr(n.Value, usage)
	case *CompoundAssignment:
		if id, ok := n.Target.(*Identifier); ok {
			if u, exists := usage[id.Name]; exists {
				u.WriteCount++
				u.ReadCount++ // Also reads for compound assignment
			}
		}
		countUsageInExpr(n.Value, usage)
	case *VariableDeclaration:
		countUsageInExpr(n.Value, usage)
	case *ReturnStatement:
		countUsageInExpr(n.Value, usage)
	case *FunctionCall:
		for _, arg := range n.Args {
			countUsageInExpr(arg, usage)
		}
	case *FunctionDefinition:
		countVariableUsage(n.Body, usage)
	case *IfStatement:
		countUsageInExpr(n.Condition, usage)
		countVariableUsage(n.ThenBody, usage)
		countVariableUsage(n.ElseBody, usage)
	case *WhileLoop:
		countUsageInExpr(n.Condition, usage)
		countVariableUsage(n.Body, usage)
	case *ForLoop:
		if n.Init != nil {
			countUsageInNode(n.Init, usage)
		}
		countUsageInExpr(n.Condition, usage)
		if n.Update != nil {
			countUsageInNode(n.Update, usage)
		}
		countVariableUsage(n.Body, usage)
	case *BinaryOp:
		countUsageInExpr(n.Left, usage)
		countUsageInExpr(n.Right, usage)
	case *Comparison:
		countUsageInExpr(n.Left, usage)
		countUsageInExpr(n.Right, usage)
	case *LogicalOp:
		countUsageInExpr(n.Left, usage)
		countUsageInExpr(n.Right, usage)
	}
}

func countUsageInExpr(expr ASTNode, usage map[string]*VariableUsage) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *Identifier:
		if u, ok := usage[e.Name]; ok {
			u.ReadCount++
		}
	case *BinaryOp:
		countUsageInExpr(e.Left, usage)
		countUsageInExpr(e.Right, usage)
	case *UnaryOp:
		countUsageInExpr(e.Operand, usage)
	case *BitwiseOp:
		countUsageInExpr(e.Left, usage)
		countUsageInExpr(e.Right, usage)
	case *Comparison:
		countUsageInExpr(e.Left, usage)
		countUsageInExpr(e.Right, usage)
	case *LogicalOp:
		countUsageInExpr(e.Left, usage)
		countUsageInExpr(e.Right, usage)
	case *FunctionCall:
		for _, arg := range e.Args {
			countUsageInExpr(arg, usage)
		}
	case *ArrayAccess:
		countUsageInExpr(e.Index, usage)
	case *TernaryOp:
		countUsageInExpr(e.Condition, usage)
		countUsageInExpr(e.TrueExpr, usage)
		countUsageInExpr(e.FalseExpr, usage)
	}
}
