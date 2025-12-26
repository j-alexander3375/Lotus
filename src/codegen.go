package main

import (
	"fmt"
	"strings"
)

// codegen.go - Assembly code generation from AST
// This file contains the main CodeGenerator struct and orchestrates
// the conversion from Abstract Syntax Tree to x86-64 GNU assembly.

// escapeAssemblyString escapes special characters in a string for assembly output
// It handles common escape sequences like newlines, tabs, quotes, and backslashes.
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

// CodeGenerator orchestrates the conversion of AST nodes to x86-64 assembly.
// It maintains state including variable symbol tables, string constants,
// and manages label generation for control flow structures.
type CodeGenerator struct {
	// Assembly section builders
	dataSection strings.Builder // Accumulates .data section (constants, strings)
	textSection strings.Builder // Accumulates .text section (executable code)

	// Symbol table and state
	variables     map[string]Variable // Maps variable names to their metadata
	constants     map[string]Variable // Maps constant names to their metadata (immutable)
	stringLengths map[string]int      // Maps variable names to string lengths
	stackOffset   int                 // Current stack offset from rbp

	// Module and import tracking
	imports *ImportContext // Tracks imported modules and functions

	// Counters for unique label generation
	stringCount int // Counter for .str labels
	labelCount  int // Counter for control flow labels

	// Program exit state
	exitCode int // Exit code from return statement

	// Diagnostic and error reporting
	diagnostics *DiagnosticManager

	// Function generation context
	inFunction               bool   // true when generating inside a function body
	currentFunctionReturnLbl string // label to jump to for function returns
}

// NewCodeGenerator creates and initializes a new code generator instance.
// All state is initialized to default values with empty maps and zero counters.
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		variables:                make(map[string]Variable),
		constants:                make(map[string]Variable),
		stringLengths:            make(map[string]int),
		stackOffset:              0,
		imports:                  NewImportContext(),
		stringCount:              0,
		labelCount:               0,
		exitCode:                 0,
		diagnostics:              NewDiagnosticManager(),
		inFunction:               false,
		currentFunctionReturnLbl: "",
	}
}

// GenerateAssembly is the main entry point for assembly generation.
// It takes a token stream, parses it into an AST, and generates x86-64 assembly.
// Returns the complete assembly program as a string.
func GenerateAssembly(tokens []Token) string {
	// Phase 1: Parse tokens to AST
	parser := NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		return fmt.Sprintf("# Parse error: %v\n", err)
	}

	// Phase 2: Optimize AST (constant folding, strength reduction, etc.)
	statements = OptimizeAST(statements)

	// Phase 3: Generate code from optimized AST
	gen := NewCodeGenerator()
	gen.dataSection.WriteString(DataSectionDirective + "\n")

	for _, stmt := range statements {
		gen.generateStatement(stmt)
	}

	// Phase 4: Apply peephole optimizations to generated assembly
	assembly := gen.buildFinalAssembly()
	return ApplyPeepholeOptimizations(assembly)
}

// generateStatement dispatches AST nodes to their appropriate code generation methods.
// This is the main router for statement-level code generation.
func (cg *CodeGenerator) generateStatement(stmt ASTNode) {
	switch s := stmt.(type) {
	case *ImportStatement:
		cg.generateImportStatement(s)
	case *VariableDeclaration:
		cg.generateVariableDeclaration(s)
	case *ConstantDeclaration:
		cg.generateConstantDeclaration(s)
	case *ReturnStatement:
		cg.generateReturnStatement(s)
	case *FunctionCall:
		cg.generateFunctionCall(s)
	case *Assignment:
		cg.generateAssignment(s)
	case *CompoundAssignment:
		cg.generateCompoundAssignment(s)
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
	case *TryStatement:
		cg.generateTryStatement(s)
	case *ThrowStatement:
		cg.generateThrowStatement(s)
	}
}

// generateVariableDeclaration allocates stack space and generates code for variable initialization.
// Handles all primitive types including integers, floats, booleans, and strings.
func (cg *CodeGenerator) generateVariableDeclaration(decl *VariableDeclaration) {
	// Allocate stack space (always 8 bytes for alignment)
	cg.stackOffset += PointerSize
	cg.variables[decl.Name] = Variable{
		Name:   decl.Name,
		Type:   decl.Type,
		Offset: cg.stackOffset,
	}

	switch decl.Type {
	case TokenTypeInt, TokenTypeInt8, TokenTypeInt16, TokenTypeInt32, TokenTypeInt64,
		TokenTypeUint, TokenTypeUint8, TokenTypeUint16, TokenTypeUint32, TokenTypeUint64:
		if lit, ok := decl.Value.(*IntLiteral); ok {
			cg.textSection.WriteString(fmt.Sprintf("    # int-type %s = %d\n", decl.Name, lit.Value))
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, -%d(%%rbp)\n", lit.Value, cg.stackOffset))
			return
		}
	case TokenTypeFloat:
		if lit, ok := decl.Value.(*FloatLiteral); ok {
			cg.textSection.WriteString(fmt.Sprintf("    # float %s\n", decl.Name))
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, -%d(%%rbp)\n", lit.Value, cg.stackOffset))
			return
		}
	case TokenTypeBool:
		if lit, ok := decl.Value.(*BoolLiteral); ok {
			boolVal := 0
			if lit.Value {
				boolVal = 1
			}
			cg.textSection.WriteString(fmt.Sprintf("    # bool %s = %v\n", decl.Name, lit.Value))
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, -%d(%%rbp)\n", boolVal, cg.stackOffset))
			return
		}
	case TokenTypeString:
		if lit, ok := decl.Value.(*StringLiteral); ok {
			label := fmt.Sprintf("%s%d", StringLabelPrefix, cg.stringCount)
			cg.stringCount++
			escapedStr := escapeAssemblyString(lit.Value)
			cg.dataSection.WriteString(fmt.Sprintf("%s:\n    .asciz \"%s\"\n", label, escapedStr))
			cg.textSection.WriteString(fmt.Sprintf("    # string %s = \"%s\"\n", decl.Name, escapedStr))
			cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rax\n", label))
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, -%d(%%rbp)\n", cg.stackOffset))
			// Track the string length for this variable
			cg.stringLengths[decl.Name] = len(lit.Value)
			return
		}
	}

	// Fallback: evaluate expression into rax and store
	cg.textSection.WriteString(fmt.Sprintf("    # %s = <expr>\n", decl.Name))
	cg.generateExpressionToReg(decl.Value, "rax")
	cg.textSection.WriteString(fmt.Sprintf("    movq %%rax, -%d(%%rbp)\n", cg.stackOffset))
}

// generateImportStatement processes an import/use statement
// This function loads and registers imported modules and functions
func (cg *CodeGenerator) generateImportStatement(stmt *ImportStatement) {
	if err := cg.imports.ProcessImport(stmt); err != nil {
		cg.diagnostics.AddError(fmt.Sprintf("Import error: %v", err), "", 0, 0, "")
		return
	}

	// Record import in assembly comments for debugging
	cg.textSection.WriteString(fmt.Sprintf("    # use \"%s\"\n", stmt.Module))
}

// generateConstantDeclaration generates code for a constant declaration.
// Constants are stored in the data section for numeric/bool types, or as labels for strings.
// They can be referenced like variables but cannot be modified.
func (cg *CodeGenerator) generateConstantDeclaration(decl *ConstantDeclaration) {
	switch decl.Type {
	case TokenTypeInt, TokenTypeInt8, TokenTypeInt16, TokenTypeInt32, TokenTypeInt64,
		TokenTypeUint, TokenTypeUint8, TokenTypeUint16, TokenTypeUint32, TokenTypeUint64:
		if lit, ok := decl.Value.(*IntLiteral); ok {
			// Generate a label for the constant in data section
			label := fmt.Sprintf(".const_%s", decl.Name)
			cg.dataSection.WriteString(fmt.Sprintf("%s:\n    .quad %d\n", label, lit.Value))

			// Store constant metadata for later reference
			cg.constants[decl.Name] = Variable{
				Name:   decl.Name,
				Type:   decl.Type,
				Offset: -1, // Constants don't use stack offsets
			}

			cg.textSection.WriteString(fmt.Sprintf("    # const int-type %s = %d\n", decl.Name, lit.Value))
		}
	case TokenTypeBool:
		if lit, ok := decl.Value.(*BoolLiteral); ok {
			boolVal := 0
			if lit.Value {
				boolVal = 1
			}
			// Generate a label for the constant in data section
			label := fmt.Sprintf(".const_%s", decl.Name)
			cg.dataSection.WriteString(fmt.Sprintf("%s:\n    .quad %d\n", label, boolVal))

			cg.constants[decl.Name] = Variable{
				Name:   decl.Name,
				Type:   decl.Type,
				Offset: -1,
			}

			cg.textSection.WriteString(fmt.Sprintf("    # const bool %s = %v\n", decl.Name, lit.Value))
		}
	case TokenTypeString:
		if lit, ok := decl.Value.(*StringLiteral); ok {
			// String constants are already stored as labels
			label := fmt.Sprintf(".const_%s", decl.Name)
			escapedStr := escapeAssemblyString(lit.Value)
			cg.dataSection.WriteString(fmt.Sprintf("%s:\n    .asciz \"%s\"\n", label, escapedStr))

			cg.constants[decl.Name] = Variable{
				Name:   decl.Name,
				Type:   decl.Type,
				Offset: -1,
			}
			cg.stringLengths[decl.Name] = len(lit.Value)

			cg.textSection.WriteString(fmt.Sprintf("    # const string %s = \"%s\"\n", decl.Name, escapedStr))
		}
	}
}

// generateReturnStatement processes a return statement and sets the exit code.
// Currently only handles integer literals; more complex expressions will be supported later.
func (cg *CodeGenerator) generateReturnStatement(ret *ReturnStatement) {
	if cg.inFunction {
		// Evaluate return expression into RAX if provided
		if ret.Value != nil {
			cg.generateExpressionToReg(ret.Value, "rax")
		} else {
			// Default return value 0
			cg.textSection.WriteString("    xorq %rax, %rax\n")
		}
		// Jump to function epilogue
		if cg.currentFunctionReturnLbl != "" {
			cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", cg.currentFunctionReturnLbl))
		} else {
			cg.textSection.WriteString("    ret\n")
		}
		return
	}

	// Top-level program return sets exit code (supports int literal only)
	if lit, ok := ret.Value.(*IntLiteral); ok {
		cg.exitCode = lit.Value
	}
}

// generateFunctionCall dispatches function calls to user-defined or built-in functions.
// First checks for user-defined functions, then for registered print functions.
func (cg *CodeGenerator) generateFunctionCall(call *FunctionCall) {
	// Check if it's a user-defined function first
	if cg.generateUserFunctionCall(call) {
		return
	}

	// Check imported stdlib functions
	if cg.imports != nil {
		if fn, ok := cg.imports.ImportedFunctions[call.Name]; ok && fn != nil {
			fn.CodeGen(cg, call.Args)
			return
		}
	}

	// Check for module-qualified function calls (module::function)
	if strings.Contains(call.Name, "::") {
		parts := strings.SplitN(call.Name, "::", 2)
		if len(parts) == 2 {
			moduleName := parts[0]
			funcName := parts[1]
			// Look up the module and function using GetModuleFunction
			if fn := GetModuleFunction(moduleName, funcName); fn != nil {
				fn.CodeGen(cg, call.Args)
				return
			}
		}
	}

	// Check if it's a registered print function
	if printFunc, ok := RegisteredPrintFunctions[call.Name]; ok {
		printFunc.CodeGen(cg, call.Args)
		return
	}

	// If not a recognized print function, emit a comment (future: user-defined functions)
	cg.textSection.WriteString(fmt.Sprintf("    # Unknown function call: %s\n", call.Name))
}

// buildFinalAssembly constructs the complete assembly program with proper sections and entry point.
// It combines the data section, text section, and generates the program prologue and epilogue.
func (cg *CodeGenerator) buildFinalAssembly() string {
	var b strings.Builder

	// Data section with constants and strings
	b.WriteString(cg.dataSection.String())
	b.WriteString("\n")

	// Text section with code
	b.WriteString(fmt.Sprintf("%s %s\n", GlobalDirective, EntryPointLabel))
	b.WriteString(TextSectionDirective + "\n")
	b.WriteString(EntryPointLabel + ":\n")

	// Program prologue
	b.WriteString("    # Program start\n")
	b.WriteString("    movq %rsp, %rbp\n") // Set up base pointer
	b.WriteString("    subq $256, %rsp\n") // Allocate stack space (256 bytes for locals)
	b.WriteString("\n")

	// Call user-defined main if present
	if _, exists := UserDefinedFunctions["main"]; exists {
		b.WriteString("    # Call user-defined main\n")
		b.WriteString("    call .main\n")
		b.WriteString("\n")
	}

	// Program code (function bodies and statements)
	b.WriteString(cg.textSection.String())
	b.WriteString("\n")

	// Program epilogue - exit syscall (only when no user-defined main)
	if _, exists := UserDefinedFunctions["main"]; !exists {
		b.WriteString("    # Exit program\n")
		b.WriteString(fmt.Sprintf("    movq $%d, %%rax  # syscall: exit\n", SyscallExit))
		b.WriteString(fmt.Sprintf("    movq $%d, %%rdi  # exit code\n", cg.exitCode))
		b.WriteString("    syscall\n")
	}

	return b.String()
}
