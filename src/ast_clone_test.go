package main

import (
	"testing"
)

func TestCloneIntLiteral(t *testing.T) {
	original := &IntLiteral{Value: 42}
	cloned := CloneASTNode(original)

	if cloned == nil {
		t.Fatal("Cloned node is nil")
	}

	clonedInt, ok := cloned.(*IntLiteral)
	if !ok {
		t.Fatalf("Expected *IntLiteral, got %T", cloned)
	}

	if clonedInt.Value != original.Value {
		t.Errorf("Expected value %d, got %d", original.Value, clonedInt.Value)
	}

	// Verify it's a different instance
	if clonedInt == original {
		t.Error("Cloned node is same instance as original")
	}
}

func TestCloneFloatLiteral(t *testing.T) {
	original := &FloatLiteral{Value: 3140} // Stored as int64 * 1000
	cloned := CloneASTNode(original)

	clonedFloat, ok := cloned.(*FloatLiteral)
	if !ok {
		t.Fatalf("Expected *FloatLiteral, got %T", cloned)
	}

	if clonedFloat.Value != original.Value {
		t.Errorf("Expected value %d, got %d", original.Value, clonedFloat.Value)
	}
}

func TestCloneBinaryOp(t *testing.T) {
	original := &BinaryOp{
		Left:     &IntLiteral{Value: 5},
		Operator: TokenPlus,
		Right:    &IntLiteral{Value: 10},
	}

	cloned := CloneASTNode(original)
	clonedOp, ok := cloned.(*BinaryOp)
	if !ok {
		t.Fatalf("Expected *BinaryOp, got %T", cloned)
	}

	if clonedOp.Operator != original.Operator {
		t.Errorf("Expected operator %v, got %v", original.Operator, clonedOp.Operator)
	}

	// Check that children are also cloned
	if clonedOp.Left == original.Left {
		t.Error("Left child was not cloned")
	}

	if clonedOp.Right == original.Right {
		t.Error("Right child was not cloned")
	}
}

func TestCloneFunctionDefinition(t *testing.T) {
	original := &FunctionDefinition{
		Name:       "testFunc",
		ReturnType: TokenTypeInt32,
		Parameters: []FunctionParam{
			{Name: "a", Type: TokenTypeInt32},
			{Name: "b", Type: TokenTypeInt32},
		},
		Body: []ASTNode{
			&ReturnStatement{Value: &IntLiteral{Value: 42}},
		},
	}

	cloned := CloneFunctionDefinition(original)

	if cloned.Name != original.Name {
		t.Errorf("Expected name %s, got %s", original.Name, cloned.Name)
	}

	if cloned.ReturnType != original.ReturnType {
		t.Errorf("Expected return type %v, got %v", original.ReturnType, cloned.ReturnType)
	}

	if len(cloned.Parameters) != len(original.Parameters) {
		t.Errorf("Expected %d parameters, got %d", len(original.Parameters), len(cloned.Parameters))
	}

	if len(cloned.Body) != len(original.Body) {
		t.Errorf("Expected %d body statements, got %d", len(original.Body), len(cloned.Body))
	}

	// Verify it's a different instance
	if cloned == original {
		t.Error("Cloned function is same instance as original")
	}
}

func TestCloneIfStatement(t *testing.T) {
	original := &IfStatement{
		Condition: &BoolLiteral{Value: true},
		ThenBody: []ASTNode{
			&ReturnStatement{Value: &IntLiteral{Value: 1}},
		},
		ElseBody: []ASTNode{
			&ReturnStatement{Value: &IntLiteral{Value: 0}},
		},
	}

	cloned := CloneASTNode(original)
	clonedIf, ok := cloned.(*IfStatement)
	if !ok {
		t.Fatalf("Expected *IfStatement, got %T", cloned)
	}

	if clonedIf.Condition == original.Condition {
		t.Error("Condition was not cloned")
	}

	if len(clonedIf.ThenBody) != len(original.ThenBody) {
		t.Errorf("Expected %d then statements, got %d", len(original.ThenBody), len(clonedIf.ThenBody))
	}

	if len(clonedIf.ElseBody) != len(original.ElseBody) {
		t.Errorf("Expected %d else statements, got %d", len(original.ElseBody), len(clonedIf.ElseBody))
	}
}

func TestCloneAssignment(t *testing.T) {
	original := &Assignment{
		Target: &Identifier{Name: "x"},
		Value:  &IntLiteral{Value: 42},
	}

	cloned := CloneASTNode(original)
	clonedAssign, ok := cloned.(*Assignment)
	if !ok {
		t.Fatalf("Expected *Assignment, got %T", cloned)
	}

	if clonedAssign.Target == original.Target {
		t.Error("Target was not cloned")
	}

	if clonedAssign.Value == original.Value {
		t.Error("Value was not cloned")
	}
}

func TestCloneFunctionCall(t *testing.T) {
	original := &FunctionCall{
		Name: "foo",
		Args: []ASTNode{
			&IntLiteral{Value: 1},
			&IntLiteral{Value: 2},
		},
	}

	cloned := CloneASTNode(original)
	clonedCall, ok := cloned.(*FunctionCall)
	if !ok {
		t.Fatalf("Expected *FunctionCall, got %T", cloned)
	}

	if clonedCall.Name != original.Name {
		t.Errorf("Expected name %s, got %s", original.Name, clonedCall.Name)
	}

	if len(clonedCall.Args) != len(original.Args) {
		t.Errorf("Expected %d args, got %d", len(original.Args), len(clonedCall.Args))
	}

	// Verify args are cloned
	for i := range original.Args {
		if clonedCall.Args[i] == original.Args[i] {
			t.Errorf("Arg %d was not cloned", i)
		}
	}
}
