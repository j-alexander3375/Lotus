package main

import (
	"fmt"
)

// parser_templates.go - Parser functions for template and namespace syntax

// ============================================================================
// Template Parsing
// ============================================================================

// parseTemplateDeclaration parses a template<...> declaration
// Syntax: template<typename T, int N> fn T max(T a, T b) { ... }
func (p *Parser) parseTemplateDeclaration() (*TemplateDeclaration, error) {
	startLoc := p.current().Line

	// Expect 'template'
	if err := p.expect(TokenTemplate); err != nil {
		return nil, err
	}

	// Expect '<'
	if err := p.expect(TokenLess); err != nil {
		return nil, err
	}

	// Parse template parameters
	var params []TemplateParameter
	for p.current().Type != TokenGreater {
		param, err := p.parseTemplateParameter()
		if err != nil {
			return nil, err
		}
		params = append(params, param)

		// Register type parameters in parser context
		if param.IsType {
			p.templateParams[param.Name] = true
		}

		// Check for comma
		if p.current().Type == TokenComma {
			p.advance()
		} else if p.current().Type != TokenGreater {
			return nil, p.formatError("expected ',' or '>' in template parameter list")
		}
	}

	// Expect '>'
	if err := p.expect(TokenGreater); err != nil {
		return nil, err
	}

	// Skip newlines after template parameters
	for p.current().Type == TokenNewline {
		p.advance()
	}

	// Parse the templated declaration (function, struct, or class)
	var decl ASTNode
	var err error

	switch p.current().Type {
	case TokenFn:
		decl, err = p.parseFunctionDefinition()
	case TokenStruct:
		decl, err = p.parseStructDefinition()
	case TokenClass:
		return nil, p.formatError("template classes not yet supported")
	default:
		return nil, p.formatError(fmt.Sprintf("expected fn or struct after template declaration, got token type %d", p.current().Type))
	}

	// Clear template parameters from context after parsing
	for _, param := range params {
		if param.IsType {
			delete(p.templateParams, param.Name)
		}
	}

	if err != nil {
		return nil, err
	}

	return &TemplateDeclaration{
		BaseNode:    BaseNode{Location: Location{Line: startLoc}},
		Parameters:  params,
		Declaration: decl,
	}, nil
}

// parseTemplateParameter parses a single template parameter
// Syntax: typename T, int N, typename T = int
func (p *Parser) parseTemplateParameter() (TemplateParameter, error) {
	param := TemplateParameter{}

	// Check if it's a type parameter (typename) or value parameter (int, etc.)
	if p.current().Type == TokenTypename {
		param.IsType = true
		p.advance()

		// Get the parameter name
		if p.current().Type != TokenIdentifier {
			return param, p.formatError("expected identifier after 'typename'")
		}
		param.Name = p.current().Value
		p.advance()

		// Check for default type
		if p.current().Type == TokenAssign {
			p.advance()
			// Parse default type (currently just store as string)
			if !isTypeToken(p.current().Type) {
				return param, p.formatError("expected type after '=' in template parameter")
			}
			param.Default = TokenTypeName(p.current().Type)
			p.advance()
		}
	} else if isTypeToken(p.current().Type) {
		// Value parameter: int N, float Scale, etc.
		param.IsType = false
		param.Type = p.current().Type
		p.advance()

		// Get the parameter name
		if p.current().Type != TokenIdentifier {
			return param, p.formatError(fmt.Sprintf("expected identifier after type in template parameter"))
		}
		param.Name = p.current().Value
		p.advance()

		// Check for default value
		if p.current().Type == TokenAssign {
			p.advance()
			// Parse default value
			defaultValue, err := p.parseExpression()
			if err != nil {
				return param, err
			}
			param.Default = defaultValue
		}
	} else {
		return param, p.formatError(fmt.Sprintf("expected 'typename' or type in template parameter, got %s", TokenTypeName(p.current().Type)))
	}

	return param, nil
}

// parseTemplateSpecialization parses a template specialization
// Syntax: template<> fn void foo<int>() { ... }
func (p *Parser) parseTemplateSpecialization() (*TemplateSpecialization, error) {
	startLoc := p.current().Line

	// Expect 'template'
	if err := p.expect(TokenTemplate); err != nil {
		return nil, err
	}

	// Expect '<>'
	if err := p.expect(TokenLess); err != nil {
		return nil, err
	}
	if err := p.expect(TokenGreater); err != nil {
		return nil, err
	}

	// Parse the specialized declaration
	var decl ASTNode
	var templateName string
	var typeArgs []string
	var err error

	switch p.current().Type {
	case TokenFn:
		// fn RetType name<Args>(...) { ... }
		p.advance()

		// Return type
		if !isTypeToken(p.current().Type) {
			return nil, p.formatError("expected return type in template specialization")
		}
		p.advance()

		// Function name
		if p.current().Type != TokenIdentifier {
			return nil, p.formatError("expected function name in template specialization")
		}
		templateName = p.current().Value
		p.advance()

		// Parse template arguments
		typeArgs, err = p.parseTemplateArguments()
		if err != nil {
			return nil, err
		}

		// Parse function (back up to re-parse properly)
		p.pos -= 2 // Go back before function name
		decl, err = p.parseFunctionDefinition()
		if err != nil {
			return nil, err
		}

	default:
		return nil, p.formatError("expected fn in template specialization")
	}

	return &TemplateSpecialization{
		BaseNode:     BaseNode{Location: Location{Line: startLoc}},
		TemplateName: templateName,
		TypeArgs:     typeArgs,
		Declaration:  decl,
		IsPartial:    false,
	}, nil
}

// parseTemplateArguments parses template arguments in angle brackets
// Syntax: <int>, <T, int>, <float, 10>
func (p *Parser) parseTemplateArguments() ([]string, error) {
	if p.current().Type != TokenLess {
		return nil, nil // No template arguments
	}
	p.advance()

	var args []string
	for p.current().Type != TokenGreater {
		// Try to parse as type first
		if isTypeToken(p.current().Type) {
			args = append(args, TokenTypeName(p.current().Type))
			p.advance()
		} else if p.current().Type == TokenIdentifier {
			// Could be a type name or value expression
			args = append(args, p.current().Value)
			p.advance()
		} else {
			return nil, p.formatError(fmt.Sprintf("expected type or identifier in template arguments, got %s", TokenTypeName(p.current().Type)))
		}

		if p.current().Type == TokenComma {
			p.advance()
		} else if p.current().Type != TokenGreater {
			return nil, p.formatError("expected ',' or '>' in template arguments")
		}
	}

	if err := p.expect(TokenGreater); err != nil {
		return nil, err
	}

	return args, nil
}

// parseTemplateInstantiation parses a template being instantiated
// This is called when we see identifier<...> in an expression or type context
func (p *Parser) parseTemplateInstantiation(name string) (*TemplateInstantiation, error) {
	startLoc := p.current().Line

	// Parse template arguments
	typeArgs, err := p.parseTemplateArguments()
	if err != nil {
		return nil, err
	}

	return &TemplateInstantiation{
		BaseNode:     BaseNode{Location: Location{Line: startLoc}},
		TemplateName: name,
		TypeArgs:     typeArgs,
	}, nil
}

// ============================================================================
// Namespace Parsing
// ============================================================================

// parseNamespaceDeclaration parses a namespace { ... } block
// Syntax: namespace name { ... }
//
//	inline namespace name { ... }
func (p *Parser) parseNamespaceDeclaration() (*NamespaceDeclaration, error) {
	startLoc := p.current().Line
	isInline := false

	// Check for inline namespace
	if p.current().Type == TokenInline {
		isInline = true
		p.advance()
	}

	// Expect 'namespace'
	if err := p.expect(TokenNamespace); err != nil {
		return nil, err
	}

	// Namespace name (optional for anonymous namespace)
	name := ""
	if p.current().Type == TokenIdentifier {
		name = p.current().Value
		p.advance()
	}

	// Parse nested namespace: namespace foo::bar::baz { ... }
	for p.current().Type == TokenColon && p.peek().Type == TokenColon {
		p.advance() // skip first :
		p.advance() // skip second :
		if p.current().Type != TokenIdentifier {
			return nil, p.formatError("expected identifier after '::' in namespace declaration")
		}
		name += "::" + p.current().Value
		p.advance()
	}

	// Parse body
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

	return &NamespaceDeclaration{
		BaseNode: BaseNode{Location: Location{Line: startLoc}},
		Name:     name,
		Body:     body,
		IsInline: isInline,
	}, nil
}

// parseUsingDeclaration parses a using statement
// Syntax: using std::vector;
//
//	using namespace std;
//	using std::vector as vec;
func (p *Parser) parseUsingDeclaration() (*UsingDeclaration, error) {
	startLoc := p.current().Line

	// Expect 'using'
	if err := p.expect(TokenUsing); err != nil {
		return nil, err
	}

	// Check for 'using namespace'
	isNamespace := false
	if p.current().Type == TokenNamespace {
		isNamespace = true
		p.advance()
	}

	// Parse qualified name
	var namespaces []string
	var name string

	if p.current().Type != TokenIdentifier {
		return nil, p.formatError("expected identifier after 'using'")
	}

	// Could be namespace::name or just name
	firstPart := p.current().Value
	p.advance()

	// Check for ::
	for p.current().Type == TokenColon && p.peek().Type == TokenColon {
		namespaces = append(namespaces, firstPart)
		p.advance() // skip first :
		p.advance() // skip second :

		if p.current().Type != TokenIdentifier {
			return nil, p.formatError("expected identifier after '::'")
		}
		firstPart = p.current().Value
		p.advance()
	}

	if isNamespace {
		// using namespace std; - firstPart is the final namespace
		namespaces = append(namespaces, firstPart)
		name = ""
	} else {
		// using std::vector; - firstPart is the name
		name = firstPart
	}

	// Check for alias: using std::vector as vec;
	alias := ""
	if p.current().Type == TokenAs {
		p.advance()
		if p.current().Type != TokenIdentifier {
			return nil, p.formatError("expected identifier after 'as'")
		}
		alias = p.current().Value
		p.advance()
	}

	// Expect semicolon
	if err := p.expect(TokenSemi); err != nil {
		return nil, err
	}

	return &UsingDeclaration{
		BaseNode:    BaseNode{Location: Location{Line: startLoc}},
		IsNamespace: isNamespace,
		Namespaces:  namespaces,
		Name:        name,
		Alias:       alias,
	}, nil
}

// parseQualifiedName parses a namespace-qualified identifier
// Syntax: std::vector, math::max, foo::bar::baz
func (p *Parser) parseQualifiedName() (*QualifiedName, error) {
	startLoc := p.current().Line

	var namespaces []string

	// First identifier
	if p.current().Type != TokenIdentifier {
		return nil, p.formatError("expected identifier in qualified name")
	}
	firstPart := p.current().Value
	p.advance()

	// Check for ::
	if p.current().Type != TokenColon || p.peek().Type != TokenColon {
		// Not a qualified name, just a simple identifier
		return nil, p.formatError("expected '::' in qualified name")
	}

	// Parse namespace::name
	for p.current().Type == TokenColon && p.peek().Type == TokenColon {
		namespaces = append(namespaces, firstPart)
		p.advance() // skip first :
		p.advance() // skip second :

		if p.current().Type != TokenIdentifier {
			return nil, p.formatError("expected identifier after '::'")
		}
		firstPart = p.current().Value
		p.advance()

		// Check if there's another ::
		if p.current().Type != TokenColon || p.peek().Type != TokenColon {
			break
		}
	}

	// firstPart is now the final name
	name := firstPart

	return &QualifiedName{
		BaseNode:   BaseNode{Location: Location{Line: startLoc}},
		Namespaces: namespaces,
		Name:       name,
	}, nil
}

// Note: TokenTypeName and isTypeToken are already defined in error_messages.go and parser.go
