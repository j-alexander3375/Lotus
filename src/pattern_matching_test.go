package main

import (
	"testing"
)

func TestMatchExpressionStructure(t *testing.T) {
	// Test that match expression structure is correct
	value := &IntLiteral{Value: 42}

	cases := []MatchCase{
		{
			Pattern: &LiteralPattern{Value: &IntLiteral{Value: 42}},
			Body:    []ASTNode{&IntLiteral{Value: 1}},
		},
		{
			Pattern:   &WildcardPattern{},
			Body:      []ASTNode{&IntLiteral{Value: 0}},
			IsDefault: false,
		},
	}

	match := &MatchExpression{
		Value: value,
		Cases: cases,
	}

	if len(match.Cases) != 2 {
		t.Errorf("Expected 2 cases, got %d", len(match.Cases))
	}

	if match.Value.(*IntLiteral).Value != 42 {
		t.Errorf("Expected match value 42, got %d", match.Value.(*IntLiteral).Value)
	}
}

func TestWildcardPattern(t *testing.T) {
	pattern := &WildcardPattern{}

	// Wildcard patterns should always match
	if pattern == nil {
		t.Error("Wildcard pattern should not be nil")
	}
}

func TestLiteralPattern(t *testing.T) {
	intPattern := &LiteralPattern{
		Value: &IntLiteral{Value: 42},
	}

	if intPattern.Value.(*IntLiteral).Value != 42 {
		t.Errorf("Expected pattern value 42, got %d", intPattern.Value.(*IntLiteral).Value)
	}

	strPattern := &LiteralPattern{
		Value: &StringLiteral{Value: "hello"},
	}

	if strPattern.Value.(*StringLiteral).Value != "hello" {
		t.Errorf("Expected pattern value 'hello', got '%s'", strPattern.Value.(*StringLiteral).Value)
	}
}

func TestBindingPattern(t *testing.T) {
	pattern := &BindingPattern{
		Name: "x",
	}

	if pattern.Name != "x" {
		t.Errorf("Expected binding name 'x', got '%s'", pattern.Name)
	}
}

func TestRangePattern(t *testing.T) {
	pattern := &RangePattern{
		Start: &IntLiteral{Value: 1},
		End:   &IntLiteral{Value: 10},
	}

	if pattern.Start.(*IntLiteral).Value != 1 {
		t.Errorf("Expected range start 1, got %d", pattern.Start.(*IntLiteral).Value)
	}

	if pattern.End.(*IntLiteral).Value != 10 {
		t.Errorf("Expected range end 10, got %d", pattern.End.(*IntLiteral).Value)
	}
}

func TestMatchCaseWithGuard(t *testing.T) {
	matchCase := MatchCase{
		Pattern: &BindingPattern{Name: "x"},
		Guard: &BinaryOp{
			Left:     &Identifier{Name: "x"},
			Operator: TokenGreater,
			Right:    &IntLiteral{Value: 10},
		},
		Body: []ASTNode{
			&IntLiteral{Value: 1},
		},
	}

	if matchCase.Pattern.(*BindingPattern).Name != "x" {
		t.Errorf("Expected pattern binding 'x', got '%s'", matchCase.Pattern.(*BindingPattern).Name)
	}

	if matchCase.Guard == nil {
		t.Error("Expected guard to be present")
	}

	if len(matchCase.Body) != 1 {
		t.Errorf("Expected 1 body statement, got %d", len(matchCase.Body))
	}
}

func TestMatchCaseDefault(t *testing.T) {
	matchCase := MatchCase{
		IsDefault: true,
		Body:      []ASTNode{&IntLiteral{Value: 0}},
	}

	if !matchCase.IsDefault {
		t.Error("Expected IsDefault to be true")
	}
}

func TestMultipleMatchCases(t *testing.T) {
	cases := []MatchCase{
		{
			Pattern: &LiteralPattern{Value: &IntLiteral{Value: 1}},
			Body:    []ASTNode{&StringLiteral{Value: "one"}},
		},
		{
			Pattern: &LiteralPattern{Value: &IntLiteral{Value: 2}},
			Body:    []ASTNode{&StringLiteral{Value: "two"}},
		},
		{
			Pattern: &RangePattern{
				Start: &IntLiteral{Value: 3},
				End:   &IntLiteral{Value: 10},
			},
			Body: []ASTNode{&StringLiteral{Value: "range"}},
		},
		{
			Pattern: &WildcardPattern{},
			Body:    []ASTNode{&StringLiteral{Value: "default"}},
		},
	}

	if len(cases) != 4 {
		t.Errorf("Expected 4 cases, got %d", len(cases))
	}

	// Check first case is literal
	if _, ok := cases[0].Pattern.(*LiteralPattern); !ok {
		t.Error("First case should be LiteralPattern")
	}

	// Check third case is range
	if _, ok := cases[2].Pattern.(*RangePattern); !ok {
		t.Error("Third case should be RangePattern")
	}

	// Check last case is wildcard
	if _, ok := cases[3].Pattern.(*WildcardPattern); !ok {
		t.Error("Last case should be WildcardPattern")
	}
}

func TestOptionExpression(t *testing.T) {
	// Test Some variant
	someExpr := &OptionExpression{
		IsSome: true,
		Value:  &IntLiteral{Value: 42},
	}

	if !someExpr.IsSome {
		t.Error("Expected IsSome to be true")
	}

	if someExpr.Value.(*IntLiteral).Value != 42 {
		t.Errorf("Expected Some value 42, got %d", someExpr.Value.(*IntLiteral).Value)
	}

	// Test None variant
	noneExpr := &OptionExpression{
		IsSome: false,
		Value:  nil,
	}

	if noneExpr.IsSome {
		t.Error("Expected IsSome to be false for None")
	}

	if noneExpr.Value != nil {
		t.Error("Expected None value to be nil")
	}
}
