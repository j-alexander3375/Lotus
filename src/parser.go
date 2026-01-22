package main

import (
	"fmt"
	"strconv"
	"strings"
)

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
		return p.formatError(FormatExpectedToken(expected, p.current().Type, p.current().Value))
	}
	p.advance()
	return nil
}

// formatError creates a detailed error with context
func (p *Parser) formatError(msg string) error {
	tok := p.current()
	return fmt.Errorf("[%s] line %d, col %d: %s", ErrUnexpectedToken, tok.Line, tok.Column, msg)
}

// formatErrorWithCode creates an error with a specific error code
func (p *Parser) formatErrorWithCode(code ErrorCode, msg string) error {
	tok := p.current()
	return fmt.Errorf("[%s] line %d, col %d: %s", code, tok.Line, tok.Column, msg)
}

// formatErrorWithSuggestion creates an error with a suggestion
func (p *Parser) formatErrorWithSuggestion(msg, suggestion string) error {
	tok := p.current()
	if suggestion != "" {
		return fmt.Errorf("[%s] line %d, col %d: %s\n  help: %s", ErrUnexpectedToken, tok.Line, tok.Column, msg, suggestion)
	}
	return fmt.Errorf("[%s] line %d, col %d: %s", ErrUnexpectedToken, tok.Line, tok.Column, msg)
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
	case TokenFn:
		return p.parseFunctionDefinition()
	case TokenInline:
		return p.parseInlineFunctionDefinition()
	case TokenVirtual:
		return p.parseVirtualFunctionDefinition()
	case TokenOverride:
		return p.parseOverrideFunctionDefinition()
	case TokenStatic:
		return p.parseStaticDeclaration()
	case TokenLocal:
		return p.parseLocalDeclaration()
	case TokenGlobal:
		return p.parseGlobalDeclaration()
	case TokenPartial:
		return p.parsePartialApplication()
	case TokenWrap:
		return p.parseWrapperDefinition()
	case TokenAt:
		return p.parseDecoratedFunction()
	case TokenUse:
		return p.parseImportStatement()
	case TokenRet:
		return p.parseReturnStatement()
	case TokenReturn:
		return p.parseReturnStatement()
	case TokenIf:
		return p.parseIfStatement()
	case TokenWhile:
		return p.parseWhileLoop()
	case TokenFor:
		return p.parseForLoop()
	case TokenBreak:
		return p.parseBreakStatement()
	case TokenContinue:
		return p.parseContinueStatement()
	case TokenMatch:
		return p.parseMatchExpression()
	case TokenUnion:
		return p.parseUnionDefinition()
	case TokenConst:
		return p.parseConstantDeclaration()
	case TokenTypeInt, TokenTypeInt8, TokenTypeInt16, TokenTypeInt32, TokenTypeInt64,
		TokenTypeUint, TokenTypeUint8, TokenTypeUint16, TokenTypeUint32, TokenTypeUint64,
		TokenTypeString, TokenTypeBool, TokenTypeFloat32, TokenTypeFloat64:
		return p.parseVariableDeclaration()
	case TokenPrintString, TokenPrintf, TokenFPrintf, TokenPrintln, TokenSPrint, TokenSPrintf, TokenSPrintln, TokenFatalf, TokenFatalln, TokenLogf, TokenLogln:
		return p.parseFunctionCall()
	case TokenIdentifier:
		// Could be a function call or assignment or variable reference
		name := p.current().Value
		p.advance()
		switch p.current().Type {
		case TokenLParen:
			// Back up to re-parse as function call
			p.pos--
			return p.parseFunctionCall()
		case TokenColon:
			// Check for module-qualified function call: module::function()
			if p.peek().Type == TokenColon {
				p.advance() // skip first ':'
				p.advance() // skip second ':'
				if p.current().Type != TokenIdentifier {
					return nil, fmt.Errorf("expected function name after '::', got token type %d", p.current().Type)
				}
				funcName := p.current().Value
				p.advance()
				if p.current().Type != TokenLParen {
					return nil, fmt.Errorf("expected '(' after function name, got token type %d", p.current().Type)
				}
				p.advance() // skip '('

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
				p.advance() // skip ')'

				// Return a FunctionCall with module-qualified name
				return &FunctionCall{
					Name: name + "::" + funcName,
					Args: args,
				}, nil
			}
			return nil, fmt.Errorf("unexpected single ':' after identifier %s", name)
		case TokenAssign:
			// Simple assignment: identifier = expression
			p.advance()
			value, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			return &Assignment{Target: &Identifier{Name: name}, Value: value}, nil
		case TokenPlusEq, TokenMinusEq, TokenStarEq, TokenSlashEq, TokenPercentEq:
			// Compound assignment: identifier op= expression
			op := p.current().Type
			p.advance()
			value, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			return &CompoundAssignment{Target: &Identifier{Name: name}, Operator: op, Value: value}, nil
		case TokenPlusPlus:
			// Increment: identifier++
			p.advance()
			return &IncrementOp{Operand: &Identifier{Name: name}, IsPrefix: false, Operator: TokenPlusPlus}, nil
		case TokenMinusMinus:
			// Decrement: identifier--
			p.advance()
			return &IncrementOp{Operand: &Identifier{Name: name}, IsPrefix: false, Operator: TokenMinusMinus}, nil
		default:
			// Bare identifier expression
			return &Identifier{Name: name}, nil
		}
	default:
		suggestion := SuggestForTypo(p.current().Value)
		return nil, p.formatErrorWithSuggestion(FormatUnexpectedToken(p.current().Type, p.current().Value, "in statement"), suggestion)
	}
}

// parseBlock parses a sequence of statements enclosed in braces { ... }
func (p *Parser) parseBlock() ([]ASTNode, error) {
	if err := p.expect(TokenLBrace); err != nil {
		return nil, err
	}

	var body []ASTNode
	for p.current().Type != TokenRBrace && p.current().Type != TokenEOF {
		if p.current().Type == TokenNewline {
			p.advance()
			continue
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			body = append(body, stmt)
		}
		if p.current().Type == TokenSemi {
			p.advance()
		}
	}

	if err := p.expect(TokenRBrace); err != nil {
		return nil, err
	}
	return body, nil
}

// parseIfStatement parses an if/else statement: if <expr> { ... } [else { ... } | else if ...]
func (p *Parser) parseIfStatement() (*IfStatement, error) {
	if err := p.expect(TokenIf); err != nil {
		return nil, err
	}

	// Condition (no parentheses required)
	cond, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// Then block
	thenBody, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	var elseBody []ASTNode
	// Optional else
	if p.current().Type == TokenElse {
		p.advance()
		// Support else-if chaining
		if p.current().Type == TokenIf {
			nested, err := p.parseIfStatement()
			if err != nil {
				return nil, err
			}
			elseBody = []ASTNode{nested}
		} else {
			body, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			elseBody = body
		}
	}

	return &IfStatement{Condition: cond, ThenBody: thenBody, ElseBody: elseBody}, nil
}

// parseWhileLoop parses a while loop: while <expr> { ... }
func (p *Parser) parseWhileLoop() (*WhileLoop, error) {
	if err := p.expect(TokenWhile); err != nil {
		return nil, err
	}

	cond, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &WhileLoop{Condition: cond, Body: body}, nil
}

// parseForLoop parses a C-style for loop:
// for (init; cond; post) { ... }
// Parentheses are optional: for init; cond; post { ... }
func (p *Parser) parseForLoop() (*ForLoop, error) {
	if err := p.expect(TokenFor); err != nil {
		return nil, err
	}

	hasParen := false
	if p.current().Type == TokenLParen {
		hasParen = true
		p.advance()
	}

	// Init (optional)
	var init ASTNode
	if p.current().Type != TokenSemi {
		var err error
		init, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
	}
	if err := p.expect(TokenSemi); err != nil {
		return nil, err
	}

	// Condition (optional)
	var cond ASTNode
	if p.current().Type != TokenSemi {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		cond = expr
	}
	if err := p.expect(TokenSemi); err != nil {
		return nil, err
	}

	// Update (optional)
	var update ASTNode
	if hasParen {
		if p.current().Type != TokenRParen {
			var err error
			update, err = p.parseStatement()
			if err != nil {
				return nil, err
			}
		}
		if err := p.expect(TokenRParen); err != nil {
			return nil, err
		}
	} else {
		if p.current().Type != TokenLBrace {
			var err error
			update, err = p.parseStatement()
			if err != nil {
				return nil, err
			}
		}
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &ForLoop{Init: init, Condition: cond, Update: update, Body: body}, nil
}

// parseBreakStatement parses a break statement
func (p *Parser) parseBreakStatement() (*BreakStatement, error) {
	if err := p.expect(TokenBreak); err != nil {
		return nil, err
	}
	// Optional semicolon
	if p.current().Type == TokenSemi {
		p.advance()
	}
	return &BreakStatement{}, nil
}

// parseContinueStatement parses a continue statement
func (p *Parser) parseContinueStatement() (*ContinueStatement, error) {
	if err := p.expect(TokenContinue); err != nil {
		return nil, err
	}
	// Optional semicolon
	if p.current().Type == TokenSemi {
		p.advance()
	}
	return &ContinueStatement{}, nil
}

// parseInterpolatedString parses an interpolated string: $"Hello {name}"
// The tokenizer gives us the raw content, we parse out the {expr} parts
func (p *Parser) parseInterpolatedString() (*InterpolatedString, error) {
	content := p.current().Value
	p.advance()

	var parts []InterpolatedPart
	var textBuf strings.Builder
	runes := []rune(content)

	for i := 0; i < len(runes); i++ {
		if runes[i] == '{' {
			// Save accumulated text
			if textBuf.Len() > 0 {
				parts = append(parts, InterpolatedPart{IsExpr: false, Text: textBuf.String()})
				textBuf.Reset()
			}
			// Find matching closing brace
			i++
			braceCount := 1
			exprStart := i
			for i < len(runes) && braceCount > 0 {
				if runes[i] == '{' {
					braceCount++
				} else if runes[i] == '}' {
					braceCount--
				}
				if braceCount > 0 {
					i++
				}
			}
			if braceCount != 0 {
				return nil, p.formatErrorWithCode(ErrExpectedToken, "unmatched { in interpolated string")
			}
			exprStr := string(runes[exprStart:i])

			// Parse the expression string
			exprTokens := Tokenize(exprStr)
			// Filter out newlines and EOF for expression parsing
			var filteredTokens []Token
			for _, tok := range exprTokens {
				if tok.Type != TokenNewline && tok.Type != TokenEOF {
					filteredTokens = append(filteredTokens, tok)
				}
			}
			filteredTokens = append(filteredTokens, Token{Type: TokenEOF})

			exprParser := &Parser{tokens: filteredTokens, pos: 0}
			expr, err := exprParser.parseExpression()
			if err != nil {
				return nil, err
			}
			parts = append(parts, InterpolatedPart{IsExpr: true, Expr: expr})
		} else {
			textBuf.WriteRune(runes[i])
		}
	}

	// Save any remaining text
	if textBuf.Len() > 0 {
		parts = append(parts, InterpolatedPart{IsExpr: false, Text: textBuf.String()})
	}

	return &InterpolatedString{Parts: parts}, nil
}

// parseMatchExpression parses a pattern matching expression
// Syntax: match expr { case pattern => body, case pattern when guard => body, default => body }
func (p *Parser) parseMatchExpression() (*MatchExpression, error) {
	if err := p.expect(TokenMatch); err != nil {
		return nil, err
	}

	// Parse the value to match
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// Expect opening brace
	if err := p.expect(TokenLBrace); err != nil {
		return nil, err
	}

	var cases []MatchCase

	for p.current().Type != TokenRBrace && p.current().Type != TokenEOF {
		// Skip newlines
		if p.current().Type == TokenNewline {
			p.advance()
			continue
		}

		var matchCase MatchCase

		if p.current().Type == TokenDefault {
			// Default case
			p.advance()
			matchCase.IsDefault = true
		} else if p.current().Type == TokenCase {
			// Regular case
			p.advance()

			// Parse pattern
			pattern, err := p.parsePattern()
			if err != nil {
				return nil, err
			}
			matchCase.Pattern = pattern

			// Check for guard (when clause)
			if p.current().Type == TokenWhen {
				p.advance()
				guard, err := p.parseExpression()
				if err != nil {
					return nil, err
				}
				matchCase.Guard = guard
			}
		} else {
			return nil, p.formatErrorWithCode(ErrExpectedToken, "expected 'case' or 'default' in match expression")
		}

		// Expect =>
		if err := p.expect(TokenLambda); err != nil {
			return nil, err
		}

		// Parse body (either single expression or block)
		if p.current().Type == TokenLBrace {
			// Block body
			body, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			matchCase.Body = body
		} else {
			// Single expression
			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			matchCase.Body = []ASTNode{&ReturnStatement{Value: expr}}
		}

		cases = append(cases, matchCase)

		// Optional comma between cases
		if p.current().Type == TokenComma {
			p.advance()
		}
	}

	if err := p.expect(TokenRBrace); err != nil {
		return nil, err
	}

	return &MatchExpression{
		Value: value,
		Cases: cases,
	}, nil
}

// parsePattern parses a pattern for pattern matching
func (p *Parser) parsePattern() (ASTNode, error) {
	switch p.current().Type {
	case TokenInt:
		val, _ := parseIntToken(p.current().Value)
		p.advance()
		// Check for range pattern: start..end
		if p.current().Type == TokenRange {
			p.advance() // skip ..
			if p.current().Type != TokenInt {
				return nil, p.formatErrorWithCode(ErrExpectedToken, "expected integer in range pattern")
			}
			endVal, _ := parseIntToken(p.current().Value)
			p.advance()
			return &RangePattern{
				Start: &IntLiteral{Value: val},
				End:   &IntLiteral{Value: endVal},
			}, nil
		}
		return &LiteralPattern{Value: &IntLiteral{Value: val}}, nil
	case TokenString:
		val := p.current().Value
		p.advance()
		return &LiteralPattern{Value: &StringLiteral{Value: val}}, nil
	case TokenBool:
		val := p.current().Value == "true"
		p.advance()
		return &LiteralPattern{Value: &BoolLiteral{Value: val}}, nil
	case TokenIdentifier:
		name := p.current().Value
		p.advance()
		// Check if it's a wildcard pattern (_)
		if name == "_" {
			return &WildcardPattern{}, nil
		}
		// Otherwise it's a binding pattern
		return &BindingPattern{Name: name}, nil
	default:
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected pattern, got "+TokenTypeName(p.current().Type))
	}
}

// parseUnionDefinition parses a union type definition
// Syntax: union Name { Variant1(Type), Variant2(Type1, Type2), Variant3 }
func (p *Parser) parseUnionDefinition() (*UnionDefinition, error) {
	if err := p.expect(TokenUnion); err != nil {
		return nil, err
	}

	// Union name
	if p.current().Type != TokenIdentifier {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected union name")
	}
	name := p.current().Value
	p.advance()

	// Opening brace
	if err := p.expect(TokenLBrace); err != nil {
		return nil, err
	}

	var variants []UnionVariant

	for p.current().Type != TokenRBrace && p.current().Type != TokenEOF {
		// Skip newlines
		if p.current().Type == TokenNewline {
			p.advance()
			continue
		}

		// Variant name
		if p.current().Type != TokenIdentifier {
			return nil, p.formatErrorWithCode(ErrExpectedToken, "expected variant name")
		}
		variantName := p.current().Value
		p.advance()

		var fields []TokenType

		// Check for fields
		if p.current().Type == TokenLParen {
			p.advance()
			for p.current().Type != TokenRParen {
				if isTypeToken(p.current().Type) {
					fields = append(fields, p.current().Type)
					p.advance()
				}
				if p.current().Type == TokenComma {
					p.advance()
				}
			}
			if err := p.expect(TokenRParen); err != nil {
				return nil, err
			}
		}

		variants = append(variants, UnionVariant{
			Name:   variantName,
			Fields: fields,
		})

		// Optional comma
		if p.current().Type == TokenComma {
			p.advance()
		}
	}

	if err := p.expect(TokenRBrace); err != nil {
		return nil, err
	}

	return &UnionDefinition{
		Name:     name,
		Variants: variants,
	}, nil
}

// parseReturnStatement parses a return statement
func (p *Parser) parseReturnStatement() (*ReturnStatement, error) {
	// Support both 'ret' and 'return'
	if p.current().Type == TokenRet || p.current().Type == TokenReturn {
		p.advance()
	} else {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected 'ret' or 'return', got "+TokenTypeName(p.current().Type))
	}

	// Return value can be any expression (optional)
	var value ASTNode
	// If next token looks like start of an expression, parse it
	switch p.current().Type {
	case TokenInt, TokenString, TokenBool, TokenFloat, TokenIdentifier, TokenLParen,
		TokenMinus, TokenExclaim, TokenAmpersand, TokenStar, TokenMatch:
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		value = expr
	}

	return &ReturnStatement{Value: value}, nil
}

// parseVariableDeclaration parses a variable declaration
func (p *Parser) parseVariableDeclaration() (*VariableDeclaration, error) {
	varType := p.current().Type
	p.advance()

	// Check for array type: int[], string[], etc.
	isArray := false
	if p.current().Type == TokenLBracket {
		p.advance() // skip [
		if err := p.expect(TokenRBracket); err != nil {
			return nil, p.formatErrorWithCode(ErrExpectedToken, "expected ']' in array type declaration")
		}
		isArray = true
	}

	if p.current().Type != TokenIdentifier {
		return nil, p.formatErrorWithCode(ErrExpectedToken, MsgMissingIdentifier+", got "+TokenTypeName(p.current().Type))
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
		Name:    varName,
		Type:    varType,
		Value:   value,
		IsArray: isArray,
	}, nil
}

// parseConstantDeclaration parses a constant declaration: const int MAX = 100;
func (p *Parser) parseConstantDeclaration() (*ConstantDeclaration, error) {
	if err := p.expect(TokenConst); err != nil {
		return nil, err
	}

	// Type is required for constants
	if !isTypeToken(p.current().Type) {
		return nil, p.formatErrorWithCode(ErrExpectedToken, MsgMissingType+" after 'const', got "+TokenTypeName(p.current().Type))
	}
	constType := p.current().Type
	p.advance()

	// Constant name
	if p.current().Type != TokenIdentifier {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected constant name, got "+TokenTypeName(p.current().Type))
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
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected module name (string) after 'use', got "+TokenTypeName(p.current().Type))
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
			return nil, p.formatErrorWithCode(ErrExpectedToken, "expected alias name after 'as', got "+TokenTypeName(p.current().Type))
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
		TokenTypeString, TokenTypeBool, TokenTypeChar, TokenTypeFloat32, TokenTypeFloat64, TokenTypeVoid:
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
	return p.parsePipeExpression()
}

// parsePipeExpression handles |> pipe operator (lowest precedence)
// value |> fn1 |> fn2 becomes fn2(fn1(value))
func (p *Parser) parsePipeExpression() (ASTNode, error) {
	left, err := p.parseLogicalOr()
	if err != nil {
		return nil, err
	}

	for p.current().Type == TokenPipe2 {
		p.advance() // skip |>

		// Next must be a function name
		if p.current().Type != TokenIdentifier {
			return nil, p.formatErrorWithCode(ErrExpectedToken, "expected function name after '|>'")
		}
		funcName := p.current().Value
		p.advance()

		left = &PipeExpression{
			Left:     left,
			Function: funcName,
		}
	}

	return left, nil
}

// parseLogicalOr handles || operator
func (p *Parser) parseLogicalOr() (ASTNode, error) {
	left, err := p.parseLogicalAnd()
	if err != nil {
		return nil, err
	}

	for p.current().Type == TokenOr {
		op := p.current().Type
		p.advance()
		right, err := p.parseLogicalAnd()
		if err != nil {
			return nil, err
		}
		left = &LogicalOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

// parseLogicalAnd handles && operator
func (p *Parser) parseLogicalAnd() (ASTNode, error) {
	left, err := p.parseBitwiseOr()
	if err != nil {
		return nil, err
	}

	for p.current().Type == TokenAnd {
		op := p.current().Type
		p.advance()
		right, err := p.parseBitwiseOr()
		if err != nil {
			return nil, err
		}
		left = &LogicalOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

// parseBitwiseOr handles | operator
func (p *Parser) parseBitwiseOr() (ASTNode, error) {
	left, err := p.parseBitwiseXor()
	if err != nil {
		return nil, err
	}

	for p.current().Type == TokenPipe {
		op := p.current().Type
		p.advance()
		right, err := p.parseBitwiseXor()
		if err != nil {
			return nil, err
		}
		left = &BitwiseOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

// parseBitwiseXor handles ^ operator
func (p *Parser) parseBitwiseXor() (ASTNode, error) {
	left, err := p.parseBitwiseAnd()
	if err != nil {
		return nil, err
	}

	for p.current().Type == TokenCaret {
		op := p.current().Type
		p.advance()
		right, err := p.parseBitwiseAnd()
		if err != nil {
			return nil, err
		}
		left = &BitwiseOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

// parseBitwiseAnd handles & operator (when used as binary, not address-of)
func (p *Parser) parseBitwiseAnd() (ASTNode, error) {
	left, err := p.parseEquality()
	if err != nil {
		return nil, err
	}

	for p.current().Type == TokenAmpersand && !isUnaryContext(p, left) {
		op := p.current().Type
		p.advance()
		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}
		left = &BitwiseOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

// parseEquality handles == and != operators
func (p *Parser) parseEquality() (ASTNode, error) {
	left, err := p.parseRelational()
	if err != nil {
		return nil, err
	}

	for p.current().Type == TokenEqual || p.current().Type == TokenNotEqual {
		op := p.current().Type
		p.advance()
		right, err := p.parseRelational()
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

// parseRelational handles <, >, <=, >= operators
func (p *Parser) parseRelational() (ASTNode, error) {
	left, err := p.parseShift()
	if err != nil {
		return nil, err
	}

	for p.current().Type == TokenLess || p.current().Type == TokenLessEq ||
		p.current().Type == TokenGreater || p.current().Type == TokenGreaterEq {
		op := p.current().Type
		p.advance()
		right, err := p.parseShift()
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

// parseComparison handles comparison operators (now delegates to equality/relational)
func (p *Parser) parseComparison() (ASTNode, error) {
	return p.parseEquality()
}

// parseShift handles << and >> operators
func (p *Parser) parseShift() (ASTNode, error) {
	left, err := p.parseAddSubtract()
	if err != nil {
		return nil, err
	}

	for p.current().Type == TokenLShift || p.current().Type == TokenRShift {
		op := p.current().Type
		p.advance()
		right, err := p.parseAddSubtract()
		if err != nil {
			return nil, err
		}
		left = &BitwiseOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
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
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.current().Type == TokenStar || p.current().Type == TokenSlash || p.current().Type == TokenPercent {
		op := p.current().Type
		p.advance()
		right, err := p.parseUnary()
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

// parseUnary handles unary operators including bitwise NOT (~)
func (p *Parser) parseUnary() (ASTNode, error) {
	if p.current().Type == TokenMinus || p.current().Type == TokenExclaim || p.current().Type == TokenAmpersand || p.current().Type == TokenStar || p.current().Type == TokenTilde {
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

// isUnaryContext checks if & should be treated as unary (address-of) vs binary (bitwise AND)
// This is a simple heuristic: & is unary if it appears after operators or at statement start
func isUnaryContext(p *Parser, left ASTNode) bool {
	// If left is nil or we're at the start, it's unary
	if left == nil {
		return true
	}
	// For now, assume & after a primary expression is binary
	return false
}

// parsePrimary handles primary expressions (literals, identifiers, parentheses, lambdas)
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
	case TokenInterpolatedString:
		return p.parseInterpolatedString()
	case TokenChar:
		val := p.current().Value
		p.advance()
		return &CharLiteral{Value: val}, nil
	case TokenBool:
		val := p.current().Value == "true"
		p.advance()
		return &BoolLiteral{Value: val}, nil
	case TokenFloat:
		val, _ := parseFloatToken(p.current().Value)
		p.advance()
		return &FloatLiteral{Value: val}, nil
	case TokenPipe:
		// Lambda expression: |params| => expr
		return p.parseLambdaExpression()
	case TokenFn:
		// Function reference: fn functionName
		return p.parseFunctionReference()
	case TokenMatch:
		// Match expression
		return p.parseMatchExpression()
	case TokenIdentifier:
		name := p.current().Value
		p.advance()
		// If followed by '(', treat as a function call in expression context
		if p.current().Type == TokenLParen {
			// Step back to re-parse as a function call
			p.pos--
			return p.parseFunctionCall()
		}
		// If followed by '[', treat as array index access
		if p.current().Type == TokenLBracket {
			p.advance() // skip '['
			index, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			if err := p.expect(TokenRBracket); err != nil {
				return nil, p.formatErrorWithCode(ErrExpectedToken, "expected ']' after array index")
			}
			return &ArrayAccess{Array: &Identifier{Name: name}, Index: index}, nil
		}
		// Check for module-qualified function call: module::function()
		if p.current().Type == TokenColon && p.peek().Type == TokenColon {
			p.advance() // skip first ':'
			p.advance() // skip second ':'
			if p.current().Type != TokenIdentifier {
				return nil, p.formatErrorWithCode(ErrExpectedToken, "expected function name after '::', got "+TokenTypeName(p.current().Type))
			}
			funcName := p.current().Value
			p.advance()
			if p.current().Type != TokenLParen {
				return nil, p.formatErrorWithCode(ErrExpectedToken, MsgMissingParameterList+", got "+TokenTypeName(p.current().Type))
			}
			p.advance() // skip '('

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
			p.advance() // skip ')'

			// Return a FunctionCall with module-qualified name
			return &FunctionCall{
				Name: name + "::" + funcName,
				Args: args,
			}, nil
		}
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
	case TokenLBracket:
		// Array literal: [1, 2, 3]
		return p.parseArrayLiteral()
	case TokenPartial:
		return p.parsePartialApplication()
	case TokenBitcast, TokenTransmute:
		return p.parseBitCastExpression()
	default:
		return nil, p.formatErrorWithSuggestion(FormatUnexpectedToken(p.current().Type, p.current().Value, "in expression"), SuggestForTypo(p.current().Value))
	}
}

// Helper functions
func parseIntToken(s string) (int, error) {
	var val int64
	var err error

	// Handle hex, octal, binary, and decimal literals
	if len(s) > 2 && s[0] == '0' {
		switch s[1] {
		case 'x', 'X':
			// Hex literal
			_, err = fmt.Sscanf(s, "%v", &val)
			if err != nil {
				// Try parsing manually
				val, err = strconv.ParseInt(s[2:], 16, 64)
			}
		case 'o', 'O':
			// Octal literal
			val, err = strconv.ParseInt(s[2:], 8, 64)
		case 'b', 'B':
			// Binary literal
			val, err = strconv.ParseInt(s[2:], 2, 64)
		default:
			// Regular decimal
			_, err = fmt.Sscanf(s, "%d", &val)
		}
	} else {
		_, err = fmt.Sscanf(s, "%d", &val)
	}
	return int(val), err
}

func parseFloatToken(s string) (int64, error) {
	var val float64
	_, err := fmt.Sscanf(s, "%f", &val)
	return int64(val * 1000), err
}

// parseFunctionDefinition parses a C-style function declaration prefixed with 'fn'
// Syntax: fn <return_type> <name>(<param_type> <param_name>, ...) { <body> }
func (p *Parser) parseFunctionDefinition() (*FunctionDefinition, error) {
	// 'fn'
	if err := p.expect(TokenFn); err != nil {
		return nil, err
	}

	// Return type
	if !isTypeToken(p.current().Type) {
		return nil, p.formatErrorWithCode(ErrExpectedToken, MsgMissingReturnType+", got "+TokenTypeName(p.current().Type))
	}
	retType := p.current().Type
	p.advance()

	// Function name
	if p.current().Type != TokenIdentifier {
		return nil, p.formatErrorWithCode(ErrExpectedToken, MsgMissingFunctionName+", got "+TokenTypeName(p.current().Type))
	}
	name := p.current().Value
	p.advance()

	// Parameters
	if err := p.expect(TokenLParen); err != nil {
		return nil, err
	}

	var params []FunctionParam
	for p.current().Type != TokenRParen {
		// Parameter type
		if !isTypeToken(p.current().Type) {
			return nil, p.formatErrorWithCode(ErrExpectedToken, "expected parameter type, got "+TokenTypeName(p.current().Type))
		}
		pType := p.current().Type
		p.advance()

		// Parameter name
		if p.current().Type != TokenIdentifier {
			return nil, p.formatErrorWithCode(ErrExpectedToken, "expected parameter name, got "+TokenTypeName(p.current().Type))
		}
		pName := p.current().Value
		p.advance()

		params = append(params, FunctionParam{Name: pName, Type: pType})

		if p.current().Type == TokenComma {
			p.advance()
		}
	}

	if err := p.expect(TokenRParen); err != nil {
		return nil, err
	}

	// Function body
	if err := p.expect(TokenLBrace); err != nil {
		return nil, err
	}

	var body []ASTNode
	for p.current().Type != TokenRBrace && p.current().Type != TokenEOF {
		// Skip newlines inside body
		if p.current().Type == TokenNewline {
			p.advance()
			continue
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			body = append(body, stmt)
		}
		if p.current().Type == TokenSemi {
			p.advance()
		}
	}

	if err := p.expect(TokenRBrace); err != nil {
		return nil, err
	}

	return &FunctionDefinition{
		Name:       name,
		Parameters: params,
		ReturnType: retType,
		Body:       body,
	}, nil
}

// parseInlineFunctionDefinition parses an inline function definition
// Syntax: inline fn <return_type> <name>(<params>) { <body> }
func (p *Parser) parseInlineFunctionDefinition() (*FunctionDefinition, error) {
	// 'inline' keyword
	if err := p.expect(TokenInline); err != nil {
		return nil, err
	}

	// Parse the function definition part
	funcDef, err := p.parseFunctionDefinition()
	if err != nil {
		return nil, err
	}

	// Mark as inline
	funcDef.IsInline = true
	return funcDef, nil
}

// parseLambdaExpression parses a lambda expression
// Syntax: |params| => expr  or  |params| => { body }
// Also supports: |int x, int y| => x + y
func (p *Parser) parseLambdaExpression() (*LambdaExpression, error) {
	// Expect '|'
	if err := p.expect(TokenPipe); err != nil {
		return nil, err
	}

	// Parse parameters
	var params []FunctionParam
	for p.current().Type != TokenPipe {
		// Parameter type
		if !isTypeToken(p.current().Type) {
			return nil, p.formatErrorWithCode(ErrExpectedToken, "expected parameter type in lambda, got "+TokenTypeName(p.current().Type))
		}
		pType := p.current().Type
		p.advance()

		// Parameter name
		if p.current().Type != TokenIdentifier {
			return nil, p.formatErrorWithCode(ErrExpectedToken, "expected parameter name in lambda, got "+TokenTypeName(p.current().Type))
		}
		pName := p.current().Value
		p.advance()

		params = append(params, FunctionParam{Name: pName, Type: pType})

		if p.current().Type == TokenComma {
			p.advance()
		}
	}

	// Expect closing '|'
	if err := p.expect(TokenPipe); err != nil {
		return nil, err
	}

	// Expect '=>'
	if err := p.expect(TokenLambda); err != nil {
		return nil, err
	}

	lambda := &LambdaExpression{
		Parameters: params,
	}

	// Check for block body { ... } or single expression
	if p.current().Type == TokenLBrace {
		// Block body
		p.advance() // skip '{'
		var body []ASTNode
		for p.current().Type != TokenRBrace && p.current().Type != TokenEOF {
			if p.current().Type == TokenNewline {
				p.advance()
				continue
			}
			stmt, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			if stmt != nil {
				body = append(body, stmt)
			}
			if p.current().Type == TokenSemi {
				p.advance()
			}
		}
		if err := p.expect(TokenRBrace); err != nil {
			return nil, err
		}
		lambda.Body = body
	} else {
		// Single expression body
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		lambda.Expression = expr
	}

	return lambda, nil
}

// parseFunctionReference parses a function reference expression
// Syntax: fn functionName (returns the function as a value)
func (p *Parser) parseFunctionReference() (*FunctionReference, error) {
	// 'fn' keyword
	if err := p.expect(TokenFn); err != nil {
		return nil, err
	}

	// Function name
	if p.current().Type != TokenIdentifier {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected function name after 'fn', got "+TokenTypeName(p.current().Type))
	}
	name := p.current().Value
	p.advance()

	return &FunctionReference{FunctionName: name}, nil
}

// parseVirtualFunctionDefinition parses a virtual function definition
// Syntax: vrt fn <return_type> <name>(<params>) { <body> }
func (p *Parser) parseVirtualFunctionDefinition() (*FunctionDefinition, error) {
	// 'vrt' keyword
	if err := p.expect(TokenVirtual); err != nil {
		return nil, err
	}

	// Parse the function definition part
	funcDef, err := p.parseFunctionDefinition()
	if err != nil {
		return nil, err
	}

	// Mark as virtual
	funcDef.IsVirtual = true
	return funcDef, nil
}

// parseOverrideFunctionDefinition parses an override function definition
// Syntax: override fn <return_type> <name>(<params>) { <body> }
func (p *Parser) parseOverrideFunctionDefinition() (*FunctionDefinition, error) {
	// 'override' keyword
	if err := p.expect(TokenOverride); err != nil {
		return nil, err
	}

	// Parse the function definition part
	funcDef, err := p.parseFunctionDefinition()
	if err != nil {
		return nil, err
	}

	// Mark as override
	funcDef.IsOverride = true
	return funcDef, nil
}

// parseStaticDeclaration parses a static variable or function declaration
// Syntax: static <type> <name> = <value>; or static fn <return_type> <name>(<params>) { <body> }
func (p *Parser) parseStaticDeclaration() (ASTNode, error) {
	// 'static' keyword
	if err := p.expect(TokenStatic); err != nil {
		return nil, err
	}

	// Check if this is a static function or static variable
	if p.current().Type == TokenFn {
		// Static function
		funcDef, err := p.parseFunctionDefinition()
		if err != nil {
			return nil, err
		}
		funcDef.IsStatic = true
		funcDef.Storage = StorageStatic
		return funcDef, nil
	}

	// Static variable - must be a type token
	if !isTypeToken(p.current().Type) {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected type or 'fn' after 'static', got "+TokenTypeName(p.current().Type))
	}

	varDecl, err := p.parseVariableDeclaration()
	if err != nil {
		return nil, err
	}
	varDecl.Storage = StorageStatic
	return varDecl, nil
}

// parseLocalDeclaration parses an explicitly local variable declaration
// Syntax: lcl <type> <name> = <value>;
func (p *Parser) parseLocalDeclaration() (ASTNode, error) {
	// 'lcl' keyword
	if err := p.expect(TokenLocal); err != nil {
		return nil, err
	}

	// Must be followed by a type token
	if !isTypeToken(p.current().Type) {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected type after 'lcl', got "+TokenTypeName(p.current().Type))
	}

	varDecl, err := p.parseVariableDeclaration()
	if err != nil {
		return nil, err
	}
	varDecl.Storage = StorageLocal
	return varDecl, nil
}

// parseGlobalDeclaration parses an explicitly global variable declaration
// Syntax: gbl <type> <name> = <value>;
func (p *Parser) parseGlobalDeclaration() (ASTNode, error) {
	// 'gbl' keyword
	if err := p.expect(TokenGlobal); err != nil {
		return nil, err
	}

	// Must be followed by a type token
	if !isTypeToken(p.current().Type) {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected type after 'gbl', got "+TokenTypeName(p.current().Type))
	}

	varDecl, err := p.parseVariableDeclaration()
	if err != nil {
		return nil, err
	}
	varDecl.Storage = StorageGlobal
	return varDecl, nil
}

// parseArrayLiteral parses an array literal: [1, 2, 3, 4]
func (p *Parser) parseArrayLiteral() (*ArrayLiteral, error) {
	if err := p.expect(TokenLBracket); err != nil {
		return nil, err
	}

	var elements []ASTNode
	elemType := TokenTypeInt // Default, will be inferred from first element

	for p.current().Type != TokenRBracket && p.current().Type != TokenEOF {
		// Skip newlines
		if p.current().Type == TokenNewline {
			p.advance()
			continue
		}

		elem, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, elem)

		// Infer type from first element
		if len(elements) == 1 {
			switch elem.(type) {
			case *IntLiteral:
				elemType = TokenTypeInt64
			case *FloatLiteral:
				elemType = TokenTypeFloat64
			case *StringLiteral, *InterpolatedString:
				elemType = TokenTypeString
			case *BoolLiteral:
				elemType = TokenTypeBool
			}
		}

		// Expect comma or closing bracket
		if p.current().Type == TokenComma {
			p.advance()
		} else if p.current().Type != TokenRBracket && p.current().Type != TokenNewline {
			break
		}
	}

	if err := p.expect(TokenRBracket); err != nil {
		return nil, err
	}

	return &ArrayLiteral{
		Elements: elements,
		ElemType: elemType,
	}, nil
}

// parsePartialApplication parses a partial function application
// Syntax: partial(fn_name, arg1, arg2, ...)(remaining_args)
// Returns the result of calling the function with all arguments
func (p *Parser) parsePartialApplication() (ASTNode, error) {
	// 'partial' keyword
	if err := p.expect(TokenPartial); err != nil {
		return nil, err
	}

	// Expect '('
	if err := p.expect(TokenLParen); err != nil {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected '(' after 'partial'")
	}

	// First argument must be a function name (identifier)
	if p.current().Type != TokenIdentifier {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected function name as first argument to 'partial'")
	}
	funcName := p.current().Value
	p.advance()

	// Parse bound arguments
	var boundArgs []ASTNode
	for p.current().Type == TokenComma {
		p.advance() // skip comma
		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		boundArgs = append(boundArgs, arg)
	}

	// Expect ')'
	if err := p.expect(TokenRParen); err != nil {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected ')' after partial arguments")
	}

	// Check if this is immediately called: partial(fn, args)(more_args)
	if p.current().Type == TokenLParen {
		p.advance() // skip '('

		// Parse remaining arguments
		var remainingArgs []ASTNode
		for p.current().Type != TokenRParen {
			arg, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			remainingArgs = append(remainingArgs, arg)
			if p.current().Type == TokenComma {
				p.advance()
			}
		}

		if err := p.expect(TokenRParen); err != nil {
			return nil, p.formatErrorWithCode(ErrExpectedToken, "expected ')' after call arguments")
		}

		// Combine bound args with remaining args and create a FunctionCall
		allArgs := append(boundArgs, remainingArgs...)
		return &FunctionCall{
			Name: funcName,
			Args: allArgs,
		}, nil
	}

	// Return just the partial application (for storing in a variable later)
	return &PartialApplication{
		FunctionName: funcName,
		BoundArgs:    boundArgs,
	}, nil
}

// parseBitCastExpression parses a bitcast expression
// Syntax: bitcast<TargetType>(expression)
// Example: bitcast<int64>(floatValue) - reinterprets float bits as int64
func (p *Parser) parseBitCastExpression() (ASTNode, error) {
	// 'bitcast' or 'transmute' keyword
	if p.current().Type != TokenBitcast && p.current().Type != TokenTransmute {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected 'bitcast' or 'transmute'")
	}
	p.advance()

	// Expect '<'
	if p.current().Type != TokenLess {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected '<' after 'bitcast'")
	}
	p.advance()

	// Parse target type
	targetType, err := p.parseTypeName()
	if err != nil {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected type name in bitcast<type>")
	}

	// Expect '>'
	if p.current().Type != TokenGreater {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected '>' after type in bitcast<type>")
	}
	p.advance()

	// Expect '('
	if err := p.expect(TokenLParen); err != nil {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected '(' after bitcast<type>")
	}

	// Parse expression to cast
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// Expect ')'
	if err := p.expect(TokenRParen); err != nil {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected ')' after bitcast expression")
	}

	return &BitCastExpression{
		TargetType: targetType,
		Value:      value,
	}, nil
}

// parseTypeName parses a type name (int, int64, float, string, etc.)
func (p *Parser) parseTypeName() (string, error) {
	switch p.current().Type {
	case TokenTypeInt:
		p.advance()
		return "int", nil
	case TokenTypeInt8:
		p.advance()
		return "int8", nil
	case TokenTypeInt16:
		p.advance()
		return "int16", nil
	case TokenTypeInt32:
		p.advance()
		return "int32", nil
	case TokenTypeInt64:
		p.advance()
		return "int64", nil
	case TokenTypeUint:
		p.advance()
		return "uint", nil
	case TokenTypeUint8:
		p.advance()
		return "uint8", nil
	case TokenTypeUint16:
		p.advance()
		return "uint16", nil
	case TokenTypeUint32:
		p.advance()
		return "uint32", nil
	case TokenTypeUint64:
		p.advance()
		return "uint64", nil
	case TokenTypeFloat32:
		p.advance()
		return "float32", nil
	case TokenTypeFloat64:
		p.advance()
		return "float64", nil
	case TokenTypeString:
		p.advance()
		return "string", nil
	case TokenTypeChar:
		p.advance()
		return "char", nil
	case TokenTypeBool:
		p.advance()
		return "bool", nil
	case TokenTypeVoid:
		p.advance()
		return "void", nil
	case TokenIdentifier:
		// Custom type (struct, class, etc.)
		name := p.current().Value
		p.advance()
		return name, nil
	default:
		return "", fmt.Errorf("expected type name, got %s", TokenTypeName(p.current().Type))
	}
}

// parseWrapperDefinition parses a wrapper/decorator definition
// Syntax: wrap fn wrapper_name(fn wrapped) { ... }
func (p *Parser) parseWrapperDefinition() (ASTNode, error) {
	// 'wrap' keyword
	if err := p.expect(TokenWrap); err != nil {
		return nil, err
	}

	// 'fn' keyword
	if err := p.expect(TokenFn); err != nil {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected 'fn' after 'wrap'")
	}

	// Wrapper name
	if p.current().Type != TokenIdentifier {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected wrapper name")
	}
	wrapperName := p.current().Value
	p.advance()

	// Expect '('
	if err := p.expect(TokenLParen); err != nil {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected '(' after wrapper name")
	}

	// Expect 'fn' keyword for wrapped function parameter
	if err := p.expect(TokenFn); err != nil {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected 'fn' for wrapped function parameter")
	}

	// Wrapped function parameter name
	if p.current().Type != TokenIdentifier {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected wrapped function parameter name")
	}
	wrappedArg := p.current().Value
	p.advance()

	// Expect ')'
	if err := p.expect(TokenRParen); err != nil {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected ')' after wrapper parameters")
	}

	// Parse body
	if err := p.expect(TokenLBrace); err != nil {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected '{' for wrapper body")
	}

	var body []ASTNode
	for p.current().Type != TokenRBrace && p.current().Type != TokenEOF {
		// Skip newlines inside wrapper body
		if p.current().Type == TokenNewline {
			p.advance()
			continue
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		body = append(body, stmt)
		// Skip optional semicolons
		for p.current().Type == TokenSemi {
			p.advance()
		}
	}

	if err := p.expect(TokenRBrace); err != nil {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected '}' after wrapper body")
	}

	return &WrapperDefinition{
		Name:       wrapperName,
		WrappedArg: wrappedArg,
		Body:       body,
	}, nil
}

// parseDecoratedFunction parses a decorated function
// Syntax: @decorator1 @decorator2 fn name() { ... }
func (p *Parser) parseDecoratedFunction() (ASTNode, error) {
	var decorators []string

	// Collect all decorators
	for p.current().Type == TokenAt {
		p.advance() // skip '@'

		if p.current().Type != TokenIdentifier {
			return nil, p.formatErrorWithCode(ErrExpectedToken, "expected decorator name after '@'")
		}
		decorators = append(decorators, p.current().Value)
		p.advance()

		// Skip newlines between decorators
		for p.current().Type == TokenNewline {
			p.advance()
		}
	}

	// Now parse the function
	if p.current().Type != TokenFn {
		return nil, p.formatErrorWithCode(ErrExpectedToken, "expected 'fn' after decorators")
	}

	funcDef, err := p.parseFunctionDefinition()
	if err != nil {
		return nil, err
	}

	return &DecoratedFunction{
		Decorators: decorators,
		Function:   funcDef,
	}, nil
}
