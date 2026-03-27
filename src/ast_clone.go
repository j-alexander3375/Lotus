package main

// ast_clone.go - Deep cloning functions for AST nodes
// Used primarily for template instantiation where we need independent copies of AST subtrees

// CloneFunctionDefinition creates a deep copy of a function definition
func CloneFunctionDefinition(fn *FunctionDefinition) *FunctionDefinition {
	cloned := &FunctionDefinition{
		Name:       fn.Name,
		ReturnType: fn.ReturnType,
		IsVirtual:  fn.IsVirtual,
		IsOverride: fn.IsOverride,
		IsStatic:   fn.IsStatic,
		IsInline:   fn.IsInline,
		Storage:    fn.Storage,
		Parameters: make([]FunctionParam, len(fn.Parameters)),
		Body:       make([]ASTNode, len(fn.Body)),
	}

	// Clone parameters
	for i, param := range fn.Parameters {
		cloned.Parameters[i] = FunctionParam{
			Name: param.Name,
			Type: param.Type,
		}
	}

	// Clone body statements
	for i, stmt := range fn.Body {
		cloned.Body[i] = CloneASTNode(stmt)
	}

	return cloned
}

// CloneASTNode creates a deep copy of any AST node
func CloneASTNode(node ASTNode) ASTNode {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *IntLiteral:
		return &IntLiteral{Value: n.Value}
	case *FloatLiteral:
		return &FloatLiteral{Value: n.Value}
	case *StringLiteral:
		return &StringLiteral{Value: n.Value}
	case *BoolLiteral:
		return &BoolLiteral{Value: n.Value}
	case *CharLiteral:
		return &CharLiteral{Value: n.Value}
	case *NullLiteral:
		return &NullLiteral{}
	case *Identifier:
		return &Identifier{Name: n.Name}

	case *BinaryOp:
		return &BinaryOp{
			Left:     CloneASTNode(n.Left),
			Operator: n.Operator,
			Right:    CloneASTNode(n.Right),
		}

	case *UnaryOp:
		return &UnaryOp{
			Operator: n.Operator,
			Operand:  CloneASTNode(n.Operand),
		}

	case *FunctionCall:
		args := make([]ASTNode, len(n.Args))
		for i, arg := range n.Args {
			args[i] = CloneASTNode(arg)
		}
		return &FunctionCall{
			Name: n.Name,
			Args: args,
		}

	case *VariableDeclaration:
		return &VariableDeclaration{
			Name:    n.Name,
			Type:    n.Type,
			Value:   CloneASTNode(n.Value),
			Storage: n.Storage,
		}

	case *Assignment:
		return &Assignment{
			Target: CloneASTNode(n.Target),
			Value:  CloneASTNode(n.Value),
		}

	case *ReturnStatement:
		return &ReturnStatement{
			Value: CloneASTNode(n.Value),
		}

	case *IfStatement:
		thenBody := make([]ASTNode, len(n.ThenBody))
		for i, stmt := range n.ThenBody {
			thenBody[i] = CloneASTNode(stmt)
		}
		elseBody := make([]ASTNode, len(n.ElseBody))
		for i, stmt := range n.ElseBody {
			elseBody[i] = CloneASTNode(stmt)
		}
		return &IfStatement{
			Condition: CloneASTNode(n.Condition),
			ThenBody:  thenBody,
			ElseBody:  elseBody,
		}

	case *WhileLoop:
		body := make([]ASTNode, len(n.Body))
		for i, stmt := range n.Body {
			body[i] = CloneASTNode(stmt)
		}
		return &WhileLoop{
			Condition: CloneASTNode(n.Condition),
			Body:      body,
		}

	case *ForLoop:
		body := make([]ASTNode, len(n.Body))
		for i, stmt := range n.Body {
			body[i] = CloneASTNode(stmt)
		}
		return &ForLoop{
			Init:      CloneASTNode(n.Init),
			Condition: CloneASTNode(n.Condition),
			Update:    CloneASTNode(n.Update),
			Body:      body,
		}

	case *BreakStatement:
		return &BreakStatement{}

	case *ContinueStatement:
		return &ContinueStatement{}

	case *CompoundAssignment:
		return &CompoundAssignment{
			Target:   CloneASTNode(n.Target),
			Operator: n.Operator,
			Value:    CloneASTNode(n.Value),
		}

	case *IncrementOp:
		return &IncrementOp{
			Operand:  CloneASTNode(n.Operand),
			IsPrefix: n.IsPrefix,
			Operator: n.Operator,
		}

	case *OptionExpression:
		return &OptionExpression{
			IsSome: n.IsSome,
			Value:  CloneASTNode(n.Value),
		}

	case *ResultExpression:
		return &ResultExpression{
			IsOk:  n.IsOk,
			Value: CloneASTNode(n.Value),
		}

	case *MethodCall:
		args := make([]ASTNode, len(n.Args))
		for i, arg := range n.Args {
			args[i] = CloneASTNode(arg)
		}
		return &MethodCall{
			Object:     CloneASTNode(n.Object),
			MethodName: n.MethodName,
			Args:       args,
			IsPointer:  n.IsPointer,
		}

	case *FieldAccess:
		return &FieldAccess{
			Object:    CloneASTNode(n.Object),
			FieldName: n.FieldName,
			IsPointer: n.IsPointer,
		}

	// For unsupported node types, return the original
	// This is safe for immutable nodes but may cause issues with mutable ones
	default:
		return node
	}
}
