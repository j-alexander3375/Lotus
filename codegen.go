package main

import (
	"fmt"
	"strings"
)

// CodeGenerator handles AST to assembly conversion
type CodeGenerator struct {
	dataSection strings.Builder
	textSection strings.Builder
	variables   map[string]Variable
	stackOffset int
	stringCount int
	exitCode    int
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		variables:   make(map[string]Variable),
		stackOffset: 0,
		stringCount: 0,
		exitCode:    0,
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
			cg.dataSection.WriteString(fmt.Sprintf("%s:\n    .asciz \"%s\"\n", label, lit.Value))
			cg.textSection.WriteString(fmt.Sprintf("    # string %s = \"%s\"\n", decl.Name, lit.Value))
			cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rax\n", label))
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, -%d(%%rbp)\n", cg.stackOffset))
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
	switch call.Name {
	case "Printf":
		cg.generatePrintf(call.Args)
	case "PrintString":
		cg.generatePrintString(call.Args)
	case "PrintInt":
		cg.generatePrintInt(call.Args)
	case "Println":
		cg.generatePrintln(call.Args)
	// Add more print functions as needed
	}
}

// generatePrintString generates code for PrintString(str)
func (cg *CodeGenerator) generatePrintString(args []ASTNode) {
	if len(args) < 1 {
		return
	}

	if lit, ok := args[0].(*StringLiteral); ok {
		label := fmt.Sprintf(".str%d", cg.stringCount)
		cg.stringCount++
		cg.dataSection.WriteString(fmt.Sprintf("%s:\n    .asciz \"%s\"\n", label, lit.Value))

		cg.textSection.WriteString("    # PrintString\n")
		cg.textSection.WriteString("    leaq " + label + "(%rip), %rdi\n")
		cg.textSection.WriteString("    movq $1, %rsi\n") // length = 1 byte for now (simplified)
		// We could call libc write() here, but for now just comment
		cg.textSection.WriteString("    # call to write() would go here\n")
	}
}

// generatePrintf generates code for Printf(format, args...)
func (cg *CodeGenerator) generatePrintf(args []ASTNode) {
	if len(args) < 1 {
		return
	}

	cg.textSection.WriteString("    # Printf\n")
	// Printf implementation would link to libc printf
	// For now, just emit a comment
	cg.textSection.WriteString("    # call to printf() would go here\n")
}

// generatePrintInt generates code for PrintInt(val)
func (cg *CodeGenerator) generatePrintInt(args []ASTNode) {
	cg.textSection.WriteString("    # PrintInt\n")
	cg.textSection.WriteString("    # call to print int would go here\n")
}

// generatePrintln generates code for Println(args...)
func (cg *CodeGenerator) generatePrintln(args []ASTNode) {
	cg.textSection.WriteString("    # Println\n")
	cg.textSection.WriteString("    # call to println would go here\n")
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
