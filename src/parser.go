package main

import "fmt"

// parser.go - Recursive descent parser for Lotus language
// This file implements syntactic analysis, converting a token stream into an AST.

// Parser holds the state for parsing a token stream
// It maintains the current position and provides methods for lookahead and navigation.
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
	case TokenUse:
		return p.parseImportStatement()
	case TokenRet:
		return p.parseReturnStatement()
	case TokenConst:
		return p.parseConstantDeclaration()
	case TokenTypeInt, TokenTypeInt8, TokenTypeInt16, TokenTypeInt32, TokenTypeInt64,
		TokenTypeUint, TokenTypeUint8, TokenTypeUint16, TokenTypeUint32, TokenTypeUint64,
		TokenTypeString, TokenTypeBool, TokenTypeFloat:
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

// parseConstantDeclaration parses a constant declaration: const int MAX = 100;
func (p *Parser) parseConstantDeclaration() (*ConstantDeclaration, error) {
	if err := p.expect(TokenConst); err != nil {
		return nil, err
	}

	// Type is required for constants
	if !isTypeToken(p.current().Type) {
		return nil, fmt.Errorf("expected type after 'const', got token type %d", p.current().Type)
	}
	constType := p.current().Type
	p.advance()

	// Constant name
	if p.current().Type != TokenIdentifier {
		return nil, fmt.Errorf("expected identifier for constant name, got token type %d", p.current().Type)
	}
	constName := p.current().Value
	p.advance()

	// Assignment operator
	if err := p.expect(TokenAssign); err != nil {
		return nil, err
	}

	// Constant value (must be a literal or compile-time constant expression)
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ConstantDeclaration{
		Name:  constName,
		Type:  constType,
		Value: value,
	}, nil
}

// parseImportStatement parses a use/import statement
// Syntax: use "module"
//
//	use "module::function"
//	use "module::*"
//	use "module" as alias
func (p *Parser) parseImportStatement() (*ImportStatement, error) {
	if err := p.expect(TokenUse); err != nil {
		return nil, err
	}

	// Module name must be a string
	if p.current().Type != TokenString {
		return nil, fmt.Errorf("expected module name (string) after 'use', got token type %d", p.current().Type)
	}
	moduleName := p.current().Value
	p.advance()

	stmt := &ImportStatement{
		Module: moduleName,
		Items:  []string{},
	}

	// Check for specific imports (::function or ::*)
	if p.current().Type == TokenColon && p.peek().Type == TokenColon {
		p.advance() // consume first :
		p.advance() // consume second :

		if p.current().Type == TokenStar {
			// Wildcard import: use "module::*"
			stmt.IsWildcard = true
			p.advance()
		} else if p.current().Type == TokenIdentifier {
			// Specific import: use "module::function"
			stmt.Items = append(stmt.Items, p.current().Value)
			p.advance()
		}
	}

	// Check for alias: as identifier
	if p.current().Type == TokenAs {
		p.advance()
		if p.current().Type != TokenIdentifier {
			return nil, fmt.Errorf("expected identifier after 'as', got token type %d", p.current().Type)
		}
		stmt.Alias = p.current().Value
		p.advance()
	}

	return stmt, nil
}

// isTypeToken checks if a token type represents a data type
func isTypeToken(t TokenType) bool {
	switch t {
	case TokenTypeInt, TokenTypeInt8, TokenTypeInt16, TokenTypeInt32, TokenTypeInt64,
		TokenTypeUint, TokenTypeUint8, TokenTypeUint16, TokenTypeUint32, TokenTypeUint64,
		TokenTypeString, TokenTypeBool, TokenTypeFloat:
		return true
	default:
		return false
	}
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
