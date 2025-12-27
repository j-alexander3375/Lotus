package main

import (
	"fmt"
	"strings"
)

// semantic.go - Semantic analysis for the Lotus compiler
// Performs basic scope tracking, usage analysis, and generates warnings.

// SemanticAnalyzer performs semantic analysis on the AST
type SemanticAnalyzer struct {
	diagnostics *DiagnosticManager
	options     *CompilerOptions
	scopes      []map[string]*SymbolInfo // Stack of scopes
	filePath    string
	sourceLines []string
	currentLine int // Approximate line tracking
}

// SymbolInfo holds information about a declared symbol
type SymbolInfo struct {
	Name         string
	Kind         SymbolKind
	TypeName     string
	DeclLine     int
	IsUsed       bool
	IsMutable    bool
	ShadowsOuter bool
}

// SymbolKind represents the kind of symbol
type SymbolKind int

const (
	SymbolVariable SymbolKind = iota
	SymbolConstant
	SymbolFunction
	SymbolParameter
	SymbolType
)

// NewSemanticAnalyzer creates a new semantic analyzer
func NewSemanticAnalyzer(diagnostics *DiagnosticManager, opts *CompilerOptions, filePath string, source string) *SemanticAnalyzer {
	sa := &SemanticAnalyzer{
		diagnostics: diagnostics,
		options:     opts,
		scopes:      make([]map[string]*SymbolInfo, 0),
		filePath:    filePath,
		sourceLines: strings.Split(source, "\n"),
		currentLine: 1,
	}
	// Push global scope
	sa.pushScope()
	return sa
}

// pushScope creates a new nested scope
func (sa *SemanticAnalyzer) pushScope() {
	sa.scopes = append(sa.scopes, make(map[string]*SymbolInfo))
}

// popScope removes the current scope and checks for unused variables
func (sa *SemanticAnalyzer) popScope() {
	if len(sa.scopes) == 0 {
		return
	}

	currentScope := sa.scopes[len(sa.scopes)-1]

	// Check for unused variables in this scope
	if sa.shouldWarn(CategoryUnused) {
		for name, info := range currentScope {
			if !info.IsUsed && (info.Kind == SymbolVariable || info.Kind == SymbolParameter) {
				// Don't warn about parameters starting with _ (intentionally unused)
				if strings.HasPrefix(name, "_") {
					continue
				}
				context := sa.getSourceLine(info.DeclLine)
				sa.diagnostics.AddWarningWithCategory(
					CategoryUnused,
					fmt.Sprintf("unused variable '%s'", name),
					sa.filePath,
					info.DeclLine,
					1,
					context,
				)
			}
		}
	}

	sa.scopes = sa.scopes[:len(sa.scopes)-1]
}

// declareSymbol adds a symbol to the current scope
func (sa *SemanticAnalyzer) declareSymbol(name string, kind SymbolKind, typeName string, line int) {
	if len(sa.scopes) == 0 {
		return
	}

	currentScope := sa.scopes[len(sa.scopes)-1]

	// Check for redeclaration in current scope
	if existing, ok := currentScope[name]; ok {
		sa.diagnostics.AddErrorWithCode(
			string(ErrRedefinition),
			CategorySemantic,
			fmt.Sprintf("redeclaration of '%s' (previously declared at line %d)", name, existing.DeclLine),
			sa.filePath,
			line,
			1,
			sa.getSourceLine(line),
		)
		return
	}

	// Check for shadowing
	shadowsOuter := false
	if sa.shouldWarn(CategoryShadow) {
		for i := len(sa.scopes) - 2; i >= 0; i-- {
			if outer, ok := sa.scopes[i][name]; ok {
				shadowsOuter = true
				context := sa.getSourceLine(line)
				sa.diagnostics.AddWarningWithCategory(
					CategoryShadow,
					fmt.Sprintf("variable '%s' shadows outer variable declared at line %d", name, outer.DeclLine),
					sa.filePath,
					line,
					1,
					context,
				)
				break
			}
		}
	}

	currentScope[name] = &SymbolInfo{
		Name:         name,
		Kind:         kind,
		TypeName:     typeName,
		DeclLine:     line,
		IsUsed:       false,
		IsMutable:    kind == SymbolVariable,
		ShadowsOuter: shadowsOuter,
	}
}

// useSymbol marks a symbol as used
func (sa *SemanticAnalyzer) useSymbol(name string) {
	// Search from innermost to outermost scope
	for i := len(sa.scopes) - 1; i >= 0; i-- {
		if info, ok := sa.scopes[i][name]; ok {
			info.IsUsed = true
			return
		}
	}
	// Symbol not found - could be undefined or a built-in
}

// lookupSymbol finds a symbol in any visible scope
func (sa *SemanticAnalyzer) lookupSymbol(name string) *SymbolInfo {
	for i := len(sa.scopes) - 1; i >= 0; i-- {
		if info, ok := sa.scopes[i][name]; ok {
			return info
		}
	}
	return nil
}

// shouldWarn checks if a warning category is enabled
func (sa *SemanticAnalyzer) shouldWarn(category DiagnosticCategory) bool {
	if sa.options == nil {
		return true
	}
	if sa.options.NoWarn {
		return false
	}

	switch category {
	case CategoryUnused:
		return sa.options.Wall || sa.options.WarnUnused
	case CategoryShadow:
		return sa.options.Wall || sa.options.WarnShadow
	case CategoryDeprecated:
		return sa.options.Wall || sa.options.WarnDeprecated
	default:
		return sa.options.Wall
	}
}

// getSourceLine returns the source line at the given line number
func (sa *SemanticAnalyzer) getSourceLine(line int) string {
	if line < 1 || line > len(sa.sourceLines) {
		return ""
	}
	return sa.sourceLines[line-1]
}

// Analyze performs semantic analysis on the AST
func (sa *SemanticAnalyzer) Analyze(statements []ASTNode) {
	for _, stmt := range statements {
		sa.analyzeNode(stmt)
	}

	// Pop global scope and check for unused at top level
	sa.popScope()
}

// analyzeNode analyzes a single AST node
func (sa *SemanticAnalyzer) analyzeNode(node ASTNode) {
	if node == nil {
		return
	}

	// Get line from location if available
	loc := node.Loc()
	if loc.Line > 0 {
		sa.currentLine = loc.Line
	}

	switch n := node.(type) {
	case *FunctionDefinition:
		sa.analyzeFunctionDefinition(n)
	case *VariableDeclaration:
		sa.analyzeVariableDeclaration(n)
	case *ConstantDeclaration:
		sa.analyzeConstantDeclaration(n)
	case *Assignment:
		sa.analyzeAssignment(n)
	case *Identifier:
		sa.useSymbol(n.Name)
	case *BinaryOp:
		sa.analyzeNode(n.Left)
		sa.analyzeNode(n.Right)
	case *UnaryOp:
		sa.analyzeNode(n.Operand)
	case *FunctionCall:
		sa.analyzeFunctionCall(n)
	case *IfStatement:
		sa.analyzeIfStatement(n)
	case *WhileLoop:
		sa.analyzeWhileLoop(n)
	case *ForLoop:
		sa.analyzeForLoop(n)
	case *ReturnStatement:
		if n.Value != nil {
			sa.analyzeNode(n.Value)
		}
	case *ArrayLiteral:
		for _, elem := range n.Elements {
			sa.analyzeNode(elem)
		}
	case *ArrayAccess:
		sa.analyzeNode(n.Array)
		sa.analyzeNode(n.Index)
	case *Comparison:
		sa.analyzeNode(n.Left)
		sa.analyzeNode(n.Right)
	case *LogicalOp:
		sa.analyzeNode(n.Left)
		sa.analyzeNode(n.Right)
	}
}

func (sa *SemanticAnalyzer) analyzeFunctionDefinition(fn *FunctionDefinition) {
	// Declare function in current scope
	line := fn.Loc().Line
	if line == 0 {
		line = sa.currentLine
	}
	sa.declareSymbol(fn.Name, SymbolFunction, TokenTypeName(fn.ReturnType), line)

	// Create new scope for function body
	sa.pushScope()

	// Declare parameters
	for _, param := range fn.Parameters {
		sa.declareSymbol(param.Name, SymbolParameter, TokenTypeName(param.Type), line)
		// Parameters are always "used" (passed by caller)
		if len(sa.scopes) > 0 {
			if info, ok := sa.scopes[len(sa.scopes)-1][param.Name]; ok {
				info.IsUsed = true
			}
		}
	}

	// Analyze function body
	for _, stmt := range fn.Body {
		sa.analyzeNode(stmt)
	}

	sa.popScope()
}

func (sa *SemanticAnalyzer) analyzeVariableDeclaration(decl *VariableDeclaration) {
	// Analyze initializer first (before declaring the variable)
	if decl.Value != nil {
		sa.analyzeNode(decl.Value)
	}

	// Declare the variable
	line := decl.Loc().Line
	if line == 0 {
		line = sa.currentLine
	}
	sa.declareSymbol(decl.Name, SymbolVariable, TokenTypeName(decl.Type), line)
}

func (sa *SemanticAnalyzer) analyzeConstantDeclaration(decl *ConstantDeclaration) {
	// Analyze initializer first
	if decl.Value != nil {
		sa.analyzeNode(decl.Value)
	}

	// Declare the constant
	line := decl.Loc().Line
	if line == 0 {
		line = sa.currentLine
	}
	sa.declareSymbol(decl.Name, SymbolConstant, TokenTypeName(decl.Type), line)
}

func (sa *SemanticAnalyzer) analyzeAssignment(assign *Assignment) {
	// Check that target exists
	if ident, ok := assign.Target.(*Identifier); ok {
		sa.useSymbol(ident.Name)
	}

	// Analyze the value
	if assign.Value != nil {
		sa.analyzeNode(assign.Value)
	}
}

func (sa *SemanticAnalyzer) analyzeFunctionCall(call *FunctionCall) {
	// Analyze arguments
	for _, arg := range call.Args {
		sa.analyzeNode(arg)
	}

	// Check for deprecated functions
	if sa.shouldWarn(CategoryDeprecated) {
		deprecatedFuncs := map[string]string{
			"gets": "use 'fgets' or 'readline' instead",
		}

		if replacement, ok := deprecatedFuncs[call.Name]; ok {
			sa.diagnostics.AddWarningWithCategory(
				CategoryDeprecated,
				fmt.Sprintf("'%s' is deprecated: %s", call.Name, replacement),
				sa.filePath,
				sa.currentLine,
				1,
				"",
			)
		}
	}
}

func (sa *SemanticAnalyzer) analyzeIfStatement(ifStmt *IfStatement) {
	sa.analyzeNode(ifStmt.Condition)

	sa.pushScope()
	for _, stmt := range ifStmt.ThenBody {
		sa.analyzeNode(stmt)
	}
	sa.popScope()

	if len(ifStmt.ElseBody) > 0 {
		sa.pushScope()
		for _, stmt := range ifStmt.ElseBody {
			sa.analyzeNode(stmt)
		}
		sa.popScope()
	}
}

func (sa *SemanticAnalyzer) analyzeWhileLoop(loop *WhileLoop) {
	sa.analyzeNode(loop.Condition)

	sa.pushScope()
	for _, stmt := range loop.Body {
		sa.analyzeNode(stmt)
	}
	sa.popScope()
}

func (sa *SemanticAnalyzer) analyzeForLoop(loop *ForLoop) {
	sa.pushScope()

	if loop.Init != nil {
		sa.analyzeNode(loop.Init)
	}
	if loop.Condition != nil {
		sa.analyzeNode(loop.Condition)
	}
	if loop.Update != nil {
		sa.analyzeNode(loop.Update)
	}

	for _, stmt := range loop.Body {
		sa.analyzeNode(stmt)
	}

	sa.popScope()
}

// levenshteinDistance computes the edit distance between two strings
// Used for "did you mean?" suggestions
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create distance matrix
	d := make([][]int, len(s1)+1)
	for i := range d {
		d[i] = make([]int, len(s2)+1)
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}

	// Fill in the matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			deletion := d[i-1][j] + 1
			insertion := d[i][j-1] + 1
			substitution := d[i-1][j-1] + cost

			// Use min of three
			minVal := deletion
			if insertion < minVal {
				minVal = insertion
			}
			if substitution < minVal {
				minVal = substitution
			}
			d[i][j] = minVal
		}
	}

	return d[len(s1)][len(s2)]
}

// FindSimilarSymbol finds a similar symbol name for "did you mean" suggestions
func (sa *SemanticAnalyzer) FindSimilarSymbol(name string) string {
	bestMatch := ""
	bestDist := 3 // Max distance for suggestions

	for i := len(sa.scopes) - 1; i >= 0; i-- {
		for symName := range sa.scopes[i] {
			dist := levenshteinDistance(strings.ToLower(name), strings.ToLower(symName))
			if dist < bestDist {
				bestDist = dist
				bestMatch = symName
			}
		}
	}

	return bestMatch
}
