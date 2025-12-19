package main

import "fmt"

// ASTNode represents a node in the abstract syntax tree
type ASTNode interface {
	astNode()
}

// Statement nodes
type ReturnStatement struct {
	Value ASTNode
}

func (r *ReturnStatement) astNode() {}

type VariableDeclaration struct {
	Name  string
	Type  TokenType
	Value ASTNode
}

func (v *VariableDeclaration) astNode() {}

type FunctionCall struct {
	Name string
	Args []ASTNode
}

func (f *FunctionCall) astNode() {}

// Expression nodes
type IntLiteral struct {
	Value int
}

func (i *IntLiteral) astNode() {}

type StringLiteral struct {
	Value string
}

func (s *StringLiteral) astNode() {}

type BoolLiteral struct {
	Value bool
}

func (b *BoolLiteral) astNode() {}

type FloatLiteral struct {
	Value int64 // stored as int * 1000
}

func (f *FloatLiteral) astNode() {}

type Identifier struct {
	Name string
}

func (id *Identifier) astNode() {}

// Parser holds parsing state
type Parser struct {
	tokens []Token
	pos    int
}

// NewParser creates a new parser for the given tokens
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

// current returns the current token
func (p *Parser) current() Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return Token{Type: TokenEOF}
}

// peek returns the next token without consuming it
func (p *Parser) peek() Token {
	if p.pos+1 < len(p.tokens) {
		return p.tokens[p.pos+1]
	}
	return Token{Type: TokenEOF}
}

// advance moves to the next token
func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

// expect consumes a token of the expected type or returns an error
func (p *Parser) expect(expected TokenType) error {
	if p.current().Type != expected {
		return fmt.Errorf("expected token type %d, got %d", expected, p.current().Type)
	}
	p.advance()
	return nil
}

// Parse parses the token stream and returns an AST
func (p *Parser) Parse() ([]ASTNode, error) {
	var statements []ASTNode

	for p.current().Type != TokenEOF {
		// Skip newlines
		if p.current().Type == TokenNewline {
			p.advance()
			continue
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}

		// Expect semicolon after statement
		if p.current().Type == TokenSemi {
			p.advance()
		}
	}

	return statements, nil
}

// parseStatement parses a single statement
func (p *Parser) parseStatement() (ASTNode, error) {
	switch p.current().Type {
	case TokenRet:
		return p.parseReturnStatement()
	case TokenTypeInt, TokenTypeString, TokenTypeBool, TokenTypeFloat:
		return p.parseVariableDeclaration()
	case TokenPrintString, TokenPrintf, TokenFPrintf, TokenPrintln, TokenSPrint, TokenSPrintf, TokenSPrintln, TokenFatalf, TokenFatalln, TokenLogf, TokenLogln:
		return p.parseFunctionCall()
	case TokenIdentifier:
		// Could be a function call or variable reference
		name := p.current().Value
		p.advance()
		if p.current().Type == TokenLParen {
			p.pos-- // Back up to re-parse as function call
			return p.parseFunctionCall()
		}
		// Otherwise it's just an identifier (not currently used)
		return &Identifier{Name: name}, nil
	default:
		return nil, fmt.Errorf("unexpected token type %d at position %d", p.current().Type, p.pos)
	}
}

// parseReturnStatement parses a return statement
func (p *Parser) parseReturnStatement() (*ReturnStatement, error) {
	if err := p.expect(TokenRet); err != nil {
		return nil, err
	}

	var value ASTNode
	if p.current().Type == TokenInt {
		val, _ := parseIntToken(p.current().Value)
		value = &IntLiteral{Value: val}
		p.advance()
	}

	return &ReturnStatement{Value: value}, nil
}

// parseVariableDeclaration parses a variable declaration
func (p *Parser) parseVariableDeclaration() (*VariableDeclaration, error) {
	varType := p.current().Type
	p.advance()

	if p.current().Type != TokenIdentifier {
		return nil, fmt.Errorf("expected identifier, got token type %d", p.current().Type)
	}
	varName := p.current().Value
	p.advance()

	if err := p.expect(TokenAssign); err != nil {
		return nil, err
	}

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &VariableDeclaration{
		Name:  varName,
		Type:  varType,
		Value: value,
	}, nil
}

// parseFunctionCall parses a function call
func (p *Parser) parseFunctionCall() (*FunctionCall, error) {
	name := p.current().Value
	p.advance()

	if err := p.expect(TokenLParen); err != nil {
		return nil, err
	}

	var args []ASTNode
	for p.current().Type != TokenRParen {
		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		if p.current().Type == TokenComma {
			p.advance()
		}
	}

	if err := p.expect(TokenRParen); err != nil {
		return nil, err
	}

	return &FunctionCall{
		Name: name,
		Args: args,
	}, nil
}

// parseExpression parses an expression with operator precedence
func (p *Parser) parseExpression() (ASTNode, error) {
	return p.parseAddSubtract()
}

// parseAddSubtract handles + and - operators
func (p *Parser) parseAddSubtract() (ASTNode, error) {
	left, err := p.parseMultiplyDivide()
	if err != nil {
		return nil, err
	}

	for p.current().Type == TokenPlus || p.current().Type == TokenMinus {
		op := p.current().Type
		p.advance()
		right, err := p.parseMultiplyDivide()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

// parseMultiplyDivide handles * / % operators
func (p *Parser) parseMultiplyDivide() (ASTNode, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.current().Type == TokenStar || p.current().Type == TokenSlash || p.current().Type == TokenPercent {
		op := p.current().Type
		p.advance()
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

// parseComparison handles comparison operators
func (p *Parser) parseComparison() (ASTNode, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.current().Type == TokenEqual || p.current().Type == TokenNotEqual ||
		p.current().Type == TokenLess || p.current().Type == TokenLessEq ||
		p.current().Type == TokenGreater || p.current().Type == TokenGreaterEq {
		op := p.current().Type
		p.advance()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &Comparison{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

// parseUnary handles unary operators
func (p *Parser) parseUnary() (ASTNode, error) {
	if p.current().Type == TokenMinus || p.current().Type == TokenExclaim || p.current().Type == TokenAmpersand || p.current().Type == TokenStar {
		op := p.current().Type
		p.advance()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &UnaryOp{
			Operator: op,
			Operand:  operand,
		}, nil
	}

	return p.parsePrimary()
}

// parsePrimary handles primary expressions (literals, identifiers, parentheses)
func (p *Parser) parsePrimary() (ASTNode, error) {
	switch p.current().Type {
	case TokenInt:
		val, _ := parseIntToken(p.current().Value)
		p.advance()
		return &IntLiteral{Value: val}, nil
	case TokenString:
		val := p.current().Value
		p.advance()
		return &StringLiteral{Value: val}, nil
	case TokenBool:
		val := p.current().Value == "true"
		p.advance()
		return &BoolLiteral{Value: val}, nil
	case TokenFloat:
		val, _ := parseFloatToken(p.current().Value)
		p.advance()
		return &FloatLiteral{Value: val}, nil
	case TokenIdentifier:
		name := p.current().Value
		p.advance()
		return &Identifier{Name: name}, nil
	case TokenLParen:
		p.advance()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(TokenRParen); err != nil {
			return nil, err
		}
		return expr, nil
	default:
		return nil, fmt.Errorf("unexpected token in expression: type %d", p.current().Type)
	}
}

// Helper functions
func parseIntToken(s string) (int, error) {
	var val int
	_, err := fmt.Sscanf(s, "%d", &val)
	return val, err
}

func parseFloatToken(s string) (int64, error) {
	var val float64
	_, err := fmt.Sscanf(s, "%f", &val)
	return int64(val * 1000), err
}
