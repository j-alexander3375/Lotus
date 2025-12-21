package main

import (
	"fmt"
	"strings"
)

// escapeAssemblyString escapes a string for use in assembly code
func escapeAssemblyString(s string) string {
	result := strings.Builder{}
	for _, ch := range s {
		switch ch {
		case '\n':
			result.WriteString("\\n")
		case '\t':
			result.WriteString("\\t")
		case '\r':
			result.WriteString("\\r")
		case '"':
			result.WriteString("\\\"")
		case '\\':
			result.WriteString("\\\\")
		default:
			result.WriteRune(ch)
		}
	}
	return result.String()
}

// CodeGenerator handles AST to assembly conversion
type CodeGenerator struct {
	dataSection   strings.Builder
	textSection   strings.Builder
	variables     map[string]Variable
	stringLengths map[string]int // maps variable name to string length
	stackOffset   int
	stringCount   int
	exitCode      int
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		variables:     make(map[string]Variable),
		stringLengths: make(map[string]int),
		stackOffset:   0,
		stringCount:   0,
		exitCode:      0,
	}
}

// GenerateAssembly generates x86-64 GNU assembly from tokens
func GenerateAssembly(tokens []Token) string {
	// Parse tokens to AST
	parser := NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		return fmt.Sprintf("# Parse error: %v\n", err)
	}

	// Generate code from AST
	gen := NewCodeGenerator()
	gen.dataSection.WriteString(".section .data\n")

	for _, stmt := range statements {
		gen.generateStatement(stmt)
	}

	return gen.buildFinalAssembly()
}

// generateStatement generates code for a statement
func (cg *CodeGenerator) generateStatement(stmt ASTNode) {
	switch s := stmt.(type) {
	case *VariableDeclaration:
		cg.generateVariableDeclaration(s)
	case *ReturnStatement:
		cg.generateReturnStatement(s)
	case *FunctionCall:
		cg.generateFunctionCall(s)
	case *Assignment:
		cg.generateAssignment(s)
	case *IfStatement:
		cg.generateIfStatement(s)
	case *WhileLoop:
		cg.generateWhileLoop(s)
	case *ForLoop:
		cg.generateForLoop(s)
	case *FunctionDefinition:
		cg.generateFunctionDefinition(s)
	case *StructDefinition:
		cg.generateStructDefinition(s)
	case *EnumDefinition:
		cg.generateEnumDefinition(s)
	case *ClassDefinition:
		cg.generateClassDefinition(s)
	case *ArrayDeclaration:
		cg.generateArrayDeclaration(s)
	case *MallocCall:
		cg.generateMallocCall(s)
	case *FreeCall:
		cg.generateFreeCall(s)
	}
}

// generateVariableDeclaration generates code for a variable declaration
func (cg *CodeGenerator) generateVariableDeclaration(decl *VariableDeclaration) {
	cg.stackOffset += 8
	cg.variables[decl.Name] = Variable{
		Name:   decl.Name,
		Type:   decl.Type,
		Offset: cg.stackOffset,
	}

	switch decl.Type {
	case TokenTypeInt:
		if lit, ok := decl.Value.(*IntLiteral); ok {
			cg.textSection.WriteString(fmt.Sprintf("    # int %s = %d\n", decl.Name, lit.Value))
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, -%d(%%rbp)\n", lit.Value, cg.stackOffset))
		}
	case TokenTypeFloat:
		if lit, ok := decl.Value.(*FloatLiteral); ok {
			cg.textSection.WriteString(fmt.Sprintf("    # float %s\n", decl.Name))
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, -%d(%%rbp)\n", lit.Value, cg.stackOffset))
		}
	case TokenTypeBool:
		if lit, ok := decl.Value.(*BoolLiteral); ok {
			boolVal := 0
			if lit.Value {
				boolVal = 1
			}
			cg.textSection.WriteString(fmt.Sprintf("    # bool %s = %v\n", decl.Name, lit.Value))
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, -%d(%%rbp)\n", boolVal, cg.stackOffset))
		}
	case TokenTypeString:
		if lit, ok := decl.Value.(*StringLiteral); ok {
			label := fmt.Sprintf(".str%d", cg.stringCount)
			cg.stringCount++
			escapedStr := escapeAssemblyString(lit.Value)
			cg.dataSection.WriteString(fmt.Sprintf("%s:\n    .asciz \"%s\"\n", label, escapedStr))
			cg.textSection.WriteString(fmt.Sprintf("    # string %s = \"%s\"\n", decl.Name, escapedStr))
			cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rax\n", label))
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, -%d(%%rbp)\n", cg.stackOffset))
			// Track the string length for this variable
			cg.stringLengths[decl.Name] = len(lit.Value)
		}
	}
}

// generateReturnStatement generates code for a return statement
func (cg *CodeGenerator) generateReturnStatement(ret *ReturnStatement) {
	if lit, ok := ret.Value.(*IntLiteral); ok {
		cg.exitCode = lit.Value
	}
}

// generateFunctionCall generates code for a function call
func (cg *CodeGenerator) generateFunctionCall(call *FunctionCall) {
	// Check if it's a user-defined function first
	if cg.generateUserFunctionCall(call) {
		return
	}

	// Check if it's a registered print function
	if printFunc, ok := RegisteredPrintFunctions[call.Name]; ok {
		printFunc.CodeGen(cg, call.Args)
		return
	}

	// If not a recognized print function, emit a comment (future: user-defined functions)
	cg.textSection.WriteString(fmt.Sprintf("    # Unknown function call: %s\n", call.Name))
}

// buildFinalAssembly builds the final assembly program
func (cg *CodeGenerator) buildFinalAssembly() string {
	var b strings.Builder

	b.WriteString(cg.dataSection.String())
	b.WriteString("\n")
	b.WriteString(".global _start\n")
	b.WriteString(".text\n")
	b.WriteString("_start:\n")
	b.WriteString("    # Set up stack frame\n")
	b.WriteString("    pushq %rbp\n")
	b.WriteString("    movq %rsp, %rbp\n")
	if cg.stackOffset > 0 {
		b.WriteString(fmt.Sprintf("    subq $%d, %%rsp  # Allocate stack space for variables\n", (cg.stackOffset+15)&^15))
	}
	b.WriteString("\n")
	b.WriteString(cg.textSection.String())
	b.WriteString("\n")
	b.WriteString("    # Exit\n")
	b.WriteString("    movq $60, %rax\n")
	b.WriteString(fmt.Sprintf("    movq $%d, %%rdi\n", cg.exitCode))
	b.WriteString("    syscall\n")

	return b.String()
}
