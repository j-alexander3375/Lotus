package main

import (
	"fmt"
	"strings"
)

// ast_utils.go - AST utility functions for debugging and analysis
// This file provides functions for inspecting, dumping, and analyzing AST structures.

// DumpAST prints a human-readable representation of the AST
func DumpAST(program []ASTNode) {
	fmt.Println("=== Abstract Syntax Tree ===")
	for i, node := range program {
		fmt.Printf("\nNode %d:\n", i)
		dumpNode(node, 1)
	}
	fmt.Println("\n============================")
}

// dumpNode recursively prints an AST node with indentation
func dumpNode(node ASTNode, indent int) {
	prefix := strings.Repeat("  ", indent)

	switch n := node.(type) {
	case *VariableDeclaration:
		fmt.Printf("%sVariableDeclaration\n", prefix)
		fmt.Printf("%s  Type: %d\n", prefix, n.Type)
		fmt.Printf("%s  Name: %s\n", prefix, n.Name)
		if n.Value != nil {
			fmt.Printf("%s  Value:\n", prefix)
			dumpNode(n.Value, indent+2)
		}

	case *ConstantDeclaration:
		fmt.Printf("%sConstantDeclaration\n", prefix)
		fmt.Printf("%s  Type: %d\n", prefix, n.Type)
		fmt.Printf("%s  Name: %s\n", prefix, n.Name)
		if n.Value != nil {
			fmt.Printf("%s  Value:\n", prefix)
			dumpNode(n.Value, indent+2)
		}

	case *FunctionDefinition:
		fmt.Printf("%sFunctionDefinition\n", prefix)
		fmt.Printf("%s  Name: %s\n", prefix, n.Name)
		fmt.Printf("%s  ReturnType: %d\n", prefix, n.ReturnType)
		fmt.Printf("%s  Parameters: %d\n", prefix, len(n.Parameters))
		for i, param := range n.Parameters {
			fmt.Printf("%s    Param %d: %d %s\n", prefix, i, param.Type, param.Name)
		}
		fmt.Printf("%s  Body: %d statements\n", prefix, len(n.Body))
		for i, stmt := range n.Body {
			fmt.Printf("%s    Statement %d:\n", prefix, i)
			dumpNode(stmt, indent+3)
		}

	case *FunctionCall:
		fmt.Printf("%sFunctionCall\n", prefix)
		fmt.Printf("%s  Name: %s\n", prefix, n.Name)
		fmt.Printf("%s  Arguments: %d\n", prefix, len(n.Args))
		for i, arg := range n.Args {
			fmt.Printf("%s    Arg %d:\n", prefix, i)
			dumpNode(arg, indent+2)
		}

	case *ReturnStatement:
		fmt.Printf("%sReturnStatement\n", prefix)
		if n.Value != nil {
			fmt.Printf("%s  Value:\n", prefix)
			dumpNode(n.Value, indent+2)
		}

	case *IfStatement:
		fmt.Printf("%sIfStatement\n", prefix)
		fmt.Printf("%s  Condition:\n", prefix)
		dumpNode(n.Condition, indent+2)
		fmt.Printf("%s  Then: %d statements\n", prefix, len(n.ThenBody))
		if len(n.ElseBody) > 0 {
			fmt.Printf("%s  Else: %d statements\n", prefix, len(n.ElseBody))
		}

	case *WhileLoop:
		fmt.Printf("%sWhileLoop\n", prefix)
		fmt.Printf("%s  Condition:\n", prefix)
		dumpNode(n.Condition, indent+2)
		fmt.Printf("%s  Body: %d statements\n", prefix, len(n.Body))

	case *ForLoop:
		fmt.Printf("%sForLoop\n", prefix)
		if n.Init != nil {
			fmt.Printf("%s  Init:\n", prefix)
			dumpNode(n.Init, indent+2)
		}
		if n.Condition != nil {
			fmt.Printf("%s  Condition:\n", prefix)
			dumpNode(n.Condition, indent+2)
		}
		if n.Update != nil {
			fmt.Printf("%s  Update:\n", prefix)
			dumpNode(n.Update, indent+2)
		}
		fmt.Printf("%s  Body: %d statements\n", prefix, len(n.Body))

	case *Comparison:
		fmt.Printf("%sComparison\n", prefix)
		fmt.Printf("%s  Operator: %v\n", prefix, n.Operator)
		fmt.Printf("%s  Left:\n", prefix)
		dumpNode(n.Left, indent+2)
		fmt.Printf("%s  Right:\n", prefix)
		dumpNode(n.Right, indent+2)

	case *Assignment:
		fmt.Printf("%sAssignment\n", prefix)
		fmt.Printf("%s  Target: %s\n", prefix, n.Target)
		fmt.Printf("%s  Value:\n", prefix)
		dumpNode(n.Value, indent+2)

	case *IntLiteral:
		fmt.Printf("%sIntLiteral: %d\n", prefix, n.Value)

	case *FloatLiteral:
		fmt.Printf("%sFloatLiteral: %d\n", prefix, n.Value)

	case *StringLiteral:
		fmt.Printf("%sStringLiteral: %q\n", prefix, n.Value)

	case *BoolLiteral:
		fmt.Printf("%sBoolLiteral: %t\n", prefix, n.Value)

	case *Identifier:
		fmt.Printf("%sIdentifier: %s\n", prefix, n.Name)

	case *ImportStatement:
		fmt.Printf("%sImportStatement\n", prefix)
		fmt.Printf("%s  Module: %s\n", prefix, n.Module)
		if n.Alias != "" {
			fmt.Printf("%s  Alias: %s\n", prefix, n.Alias)
		}
		if len(n.Items) > 0 {
			fmt.Printf("%s  Items: %v\n", prefix, n.Items)
		}

	default:
		fmt.Printf("%sUnknown node type: %T\n", prefix, n)
	}
}

// CountASTNodes counts the total number of nodes in an AST
func CountASTNodes(program []ASTNode) int {
	count := 0
	for _, node := range program {
		count += countNode(node)
	}
	return count
}

// countNode recursively counts nodes in a single AST node
func countNode(node ASTNode) int {
	count := 1 // Count the node itself

	switch n := node.(type) {
	case *VariableDeclaration:
		if n.Value != nil {
			count += countNode(n.Value)
		}
	case *ConstantDeclaration:
		if n.Value != nil {
			count += countNode(n.Value)
		}
	case *FunctionDefinition:
		for _, stmt := range n.Body {
			count += countNode(stmt)
		}
	case *FunctionCall:
		for _, arg := range n.Args {
			count += countNode(arg)
		}
	case *ReturnStatement:
		if n.Value != nil {
			count += countNode(n.Value)
		}
	case *IfStatement:
		count += countNode(n.Condition)
		for _, stmt := range n.ThenBody {
			count += countNode(stmt)
		}
		for _, stmt := range n.ElseBody {
			count += countNode(stmt)
		}
	case *WhileLoop:
		count += countNode(n.Condition)
		for _, stmt := range n.Body {
			count += countNode(stmt)
		}
	case *ForLoop:
		if n.Init != nil {
			count += countNode(n.Init)
		}
		if n.Condition != nil {
			count += countNode(n.Condition)
		}
		if n.Update != nil {
			count += countNode(n.Update)
		}
		for _, stmt := range n.Body {
			count += countNode(stmt)
		}
	case *Comparison:
		count += countNode(n.Left)
		count += countNode(n.Right)
	case *BinaryOp:
		count += countNode(n.Left)
		count += countNode(n.Right)
	case *BitwiseOp:
		count += countNode(n.Left)
		count += countNode(n.Right)
	case *UnaryOp:
		count += countNode(n.Operand)
	case *Assignment:
		count += countNode(n.Target)
		count += countNode(n.Value)
	}

	return count
}

// AnalyzeAST provides statistics about the AST structure
func AnalyzeAST(program []ASTNode) (funcCount, varCount, constCount int) {
	for _, node := range program {
		switch node.(type) {
		case *FunctionDefinition:
			funcCount++
		case *VariableDeclaration:
			varCount++
		case *ConstantDeclaration:
			constCount++
		}
	}
	return
}
