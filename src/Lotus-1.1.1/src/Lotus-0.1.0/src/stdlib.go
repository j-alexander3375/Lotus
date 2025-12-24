package main

import (
	"fmt"
)

// stdlib.go - Standard Library module system for Lotus
// Defines available stdlib modules and their exported functions/types.

// StdlibModule represents a module in the standard library
type StdlibModule struct {
	Name      string                     // Module name (e.g., "io", "math")
	Functions map[string]*StdlibFunction // Available functions in this module
	Types     map[string]TokenType       // Available types (future)
}

// StdlibFunction represents a function available in the stdlib
type StdlibFunction struct {
	Name     string
	Module   string // Module it belongs to
	NumArgs  int    // -1 for variadic
	ArgTypes []TokenType
	RetType  TokenType
	CodeGen  func(*CodeGenerator, []ASTNode) // Code generation function
}

// StandardLibrary holds all available stdlib modules
var StandardLibrary = map[string]*StdlibModule{
	"io":   createIOModule(),
	"mem":  createMemoryModule(),
	"math": createMathModule(),
	"str":  createStringModule(),
}

// createIOModule creates the I/O standard library module
func createIOModule() *StdlibModule {
	return &StdlibModule{
		Name: "io",
		Functions: map[string]*StdlibFunction{
			"print": {
				Name:    "print",
				Module:  "io",
				NumArgs: -1,
				CodeGen: generatePrintCode,
			},
			"println": {
				Name:    "println",
				Module:  "io",
				NumArgs: -1,
				CodeGen: generatePrintlnCode,
			},
			"printf": {
				Name:    "printf",
				Module:  "io",
				NumArgs: -1,
				CodeGen: generateIOPrintf,
			},
			"fprintf": {
				Name:    "fprintf",
				Module:  "io",
				NumArgs: -1,
				CodeGen: generateIOFprintf,
			},
			"sprint": {
				Name:    "sprint",
				Module:  "io",
				NumArgs: -1,
				CodeGen: generateIOSprint,
			},
			"sprintf": {
				Name:    "sprintf",
				Module:  "io",
				NumArgs: -1,
				CodeGen: generateIOSprintf,
			},
			"sprintln": {
				Name:    "sprintln",
				Module:  "io",
				NumArgs: -1,
				CodeGen: generateIOSprintln,
			},
		},
		Types: map[string]TokenType{},
	}
}

// createMemoryModule creates the memory management stdlib module
func createMemoryModule() *StdlibModule {
	return &StdlibModule{
		Name: "mem",
		Functions: map[string]*StdlibFunction{
			"malloc": {
				Name:    "malloc",
				Module:  "mem",
				NumArgs: 1,
				CodeGen: func(cg *CodeGenerator, args []ASTNode) {
					// Generate malloc call - will be handled in codegen
				},
			},
			"free": {
				Name:    "free",
				Module:  "mem",
				NumArgs: 1,
				CodeGen: func(cg *CodeGenerator, args []ASTNode) {
					// Generate free call
				},
			},
			"sizeof": {
				Name:    "sizeof",
				Module:  "mem",
				NumArgs: 1,
				CodeGen: func(cg *CodeGenerator, args []ASTNode) {
					// Generate sizeof calculation
				},
			},
		},
		Types: map[string]TokenType{},
	}
}

// createMathModule creates the math stdlib module
func createMathModule() *StdlibModule {
	return &StdlibModule{
		Name: "math",
		Functions: map[string]*StdlibFunction{
			"abs": {
				Name:    "abs",
				Module:  "math",
				NumArgs: 1,
				CodeGen: generateMathAbs,
			},
			"min": {
				Name:    "min",
				Module:  "math",
				NumArgs: 2,
				CodeGen: generateMathMin,
			},
			"max": {
				Name:    "max",
				Module:  "math",
				NumArgs: 2,
				CodeGen: generateMathMax,
			},
			"sqrt": {
				Name:    "sqrt",
				Module:  "math",
				NumArgs: 1,
				CodeGen: generateMathSqrt,
			},
			"pow": {
				Name:    "pow",
				Module:  "math",
				NumArgs: 2,
				CodeGen: generateMathPow,
			},
		},
		Types: map[string]TokenType{},
	}
}

// createStringModule creates the string manipulation stdlib module
func createStringModule() *StdlibModule {
	return &StdlibModule{
		Name: "str",
		Functions: map[string]*StdlibFunction{
			"len": {
				Name:    "len",
				Module:  "str",
				NumArgs: 1,
				CodeGen: generateStringLen,
			},
			"concat": {
				Name:    "concat",
				Module:  "str",
				NumArgs: -1,
				CodeGen: generateStringConcat,
			},
			"compare": {
				Name:    "compare",
				Module:  "str",
				NumArgs: 2,
				CodeGen: generateStringCompare,
			},
			"copy": {
				Name:    "copy",
				Module:  "str",
				NumArgs: 2,
				CodeGen: generateStringCopy,
			},
		},
		Types: map[string]TokenType{},
	}
}

// GetModuleFunction retrieves a function from a module
func GetModuleFunction(moduleName, funcName string) *StdlibFunction {
	if module, ok := StandardLibrary[moduleName]; ok {
		if fn, ok := module.Functions[funcName]; ok {
			return fn
		}
	}
	return nil
}

// ImportContext tracks what has been imported in the current compilation
type ImportContext struct {
	ImportedModules   map[string]string          // Maps alias to module name
	ImportedFunctions map[string]*StdlibFunction // Maps function name to function
	UseWildcard       bool                       // true if using wildcard import
}

// NewImportContext creates a new import tracking context
func NewImportContext() *ImportContext {
	return &ImportContext{
		ImportedModules:   make(map[string]string),
		ImportedFunctions: make(map[string]*StdlibFunction),
	}
}

// ProcessImport processes an import statement and adds exported items to context
func (ic *ImportContext) ProcessImport(stmt *ImportStatement) error {
	module, exists := StandardLibrary[stmt.Module]
	if !exists {
		return fmt.Errorf("module '%s' not found in standard library", stmt.Module)
	}

	alias := stmt.Alias
	if alias == "" {
		alias = stmt.Module
	}

	ic.ImportedModules[alias] = stmt.Module

	// Process specific imports
	if stmt.IsWildcard {
		// Import all functions from module
		for name, fn := range module.Functions {
			ic.ImportedFunctions[name] = fn
		}
	} else if len(stmt.Items) > 0 {
		// Import specific items
		for _, item := range stmt.Items {
			if fn, ok := module.Functions[item]; ok {
				ic.ImportedFunctions[item] = fn
			}
		}
	} else {
		// No specific items means import all
		for name, fn := range module.Functions {
			ic.ImportedFunctions[name] = fn
		}
	}

	return nil
}

// ============================================================================
// Placeholder code generation functions for stdlib functions
// ============================================================================

func generateMathAbs(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	// abs(x): if x < 0 then -x else x
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    movq %rax, %rbx\n")
	cg.textSection.WriteString("    sarq $63, %rbx\n")
	cg.textSection.WriteString("    xorq %rbx, %rax\n")
	cg.textSection.WriteString("    subq %rbx, %rax\n")
}

func generateMathMin(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	labelGT := cg.getLabel("min_gt")
	labelEnd := cg.getLabel("min_end")
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    pushq %rax\n")
	cg.generateExpressionToReg(args[1], "rax")
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString("    popq %rax\n")
	cg.textSection.WriteString("    cmpq %rcx, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jle %s\n", labelGT))
	// rcx smaller
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", labelEnd))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelGT))
	// rax already min
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelEnd))
}

func generateMathMax(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	labelLT := cg.getLabel("max_lt")
	labelEnd := cg.getLabel("max_end")
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    pushq %rax\n")
	cg.generateExpressionToReg(args[1], "rax")
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString("    popq %rax\n")
	cg.textSection.WriteString("    cmpq %rcx, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", labelLT))
	// rcx larger
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", labelEnd))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelLT))
	// rax already max
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelEnd))
}

func generateMathSqrt(cg *CodeGenerator, args []ASTNode) {
	// TODO: Implement sqrt(x)
}

func generateMathPow(cg *CodeGenerator, args []ASTNode) {
	// TODO: Implement pow(base, exp)
}

func generateStringLen(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	switch v := args[0].(type) {
	case *StringLiteral:
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rax\n", len(v.Value)))
	case *Identifier:
		if l, ok := cg.stringLengths[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rax\n", l))
			return
		}
		cg.textSection.WriteString("    movq $0, %rax\n")
	default:
		cg.textSection.WriteString("    movq $0, %rax\n")
	}
}

func generateStringConcat(cg *CodeGenerator, args []ASTNode) {
	// TODO: Implement concat(str1, str2, ...)
}

func generateStringCompare(cg *CodeGenerator, args []ASTNode) {
	// TODO: Implement compare(str1, str2)
}

func generateStringCopy(cg *CodeGenerator, args []ASTNode) {
	// TODO: Implement copy(src, dst)
}

// printfuncs.go wrappers (will use existing implementations from printfuncs.go)
func generatePrintCode(cg *CodeGenerator, args []ASTNode) {
	// References existing print logic from printfuncs.go
}

// Note: Printf and Printf in module system delegate to printfuncs implementations
// The io module functions reference the existing printf/fprintf/sprintf implementations

func generateStringMemCopy(cg *CodeGenerator, args []ASTNode) {
	// TODO: Implement memcpy(dst, src, size)
}

// IO module wrapper functions - delegate to printfuncs.go implementations
func generateIOPrintf(cg *CodeGenerator, args []ASTNode) {
	generatePrintfCode(cg, args)
}

func generateIOFprintf(cg *CodeGenerator, args []ASTNode) {
	generateFprintfCode(cg, args)
}

func generateIOSprint(cg *CodeGenerator, args []ASTNode) {
	generateSprintCode(cg, args)
}

func generateIOSprintf(cg *CodeGenerator, args []ASTNode) {
	generateSprintfCode(cg, args)
}

func generateIOSprintln(cg *CodeGenerator, args []ASTNode) {
	generateSprintlnCode(cg, args)
}
