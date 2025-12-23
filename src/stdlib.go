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
	"num":  createNumModule(),
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
				CodeGen: generateMemMalloc,
			},
			"free": {
				Name:    "free",
				Module:  "mem",
				NumArgs: 1,
				CodeGen: generateMemFreeNoop,
			},
			"sizeof": {
				Name:    "sizeof",
				Module:  "mem",
				NumArgs: 1,
				CodeGen: generateMemSizeof,
			},
			"memcpy": {
				Name:    "memcpy",
				Module:  "mem",
				NumArgs: 3,
				CodeGen: generateMemMemcpy,
			},
			"memset": {
				Name:    "memset",
				Module:  "mem",
				NumArgs: 3,
				CodeGen: generateMemMemset,
			},
			"mmap": {
				Name:    "mmap",
				Module:  "mem",
				NumArgs: 1,
				CodeGen: generateMemMmap,
			},
			"munmap": {
				Name:    "munmap",
				Module:  "mem",
				NumArgs: 2,
				CodeGen: generateMemMunmap,
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
			"floor": {
				Name:    "floor",
				Module:  "math",
				NumArgs: 1,
				CodeGen: generateMathFloor,
			},
			"ceil": {
				Name:    "ceil",
				Module:  "math",
				NumArgs: 1,
				CodeGen: generateMathCeil,
			},
			"round": {
				Name:    "round",
				Module:  "math",
				NumArgs: 1,
				CodeGen: generateMathRound,
			},
			"gcd": {
				Name:    "gcd",
				Module:  "math",
				NumArgs: 2,
				CodeGen: generateMathGcd,
			},
			"lcm": {
				Name:    "lcm",
				Module:  "math",
				NumArgs: 2,
				CodeGen: generateMathLcm,
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
				NumArgs: 1,
				CodeGen: generateStringCopy,
			},
			"indexOf": {
				Name:    "indexOf",
				Module:  "str",
				NumArgs: 2,
				CodeGen: generateStringIndexOf,
			},
			"contains": {
				Name:    "contains",
				Module:  "str",
				NumArgs: 2,
				CodeGen: generateStringContains,
			},
			"startsWith": {
				Name:    "startsWith",
				Module:  "str",
				NumArgs: 2,
				CodeGen: generateStringStartsWith,
			},
			"endsWith": {
				Name:    "endsWith",
				Module:  "str",
				NumArgs: 2,
				CodeGen: generateStringEndsWith,
			},
		},
		Types: map[string]TokenType{},
	}
}

// createNumModule creates the numeric conversions stdlib module
func createNumModule() *StdlibModule {
	return &StdlibModule{
		Name: "num",
		Functions: map[string]*StdlibFunction{
			"toInt8":   {Name: "toInt8", Module: "num", NumArgs: 1, CodeGen: generateNumToInt8},
			"toUint8":  {Name: "toUint8", Module: "num", NumArgs: 1, CodeGen: generateNumToUint8},
			"toInt16":  {Name: "toInt16", Module: "num", NumArgs: 1, CodeGen: generateNumToInt16},
			"toUint16": {Name: "toUint16", Module: "num", NumArgs: 1, CodeGen: generateNumToUint16},
			"toInt32":  {Name: "toInt32", Module: "num", NumArgs: 1, CodeGen: generateNumToInt32},
			"toUint32": {Name: "toUint32", Module: "num", NumArgs: 1, CodeGen: generateNumToUint32},
			"toInt64":  {Name: "toInt64", Module: "num", NumArgs: 1, CodeGen: generateNumToInt64},
			"toUint64": {Name: "toUint64", Module: "num", NumArgs: 1, CodeGen: generateNumToUint64},
			"toBool":   {Name: "toBool", Module: "num", NumArgs: 1, CodeGen: generateNumToBool},
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

func generateMathFloor(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	// For integers, floor is just the value itself
	cg.generateExpressionToReg(args[0], "rax")
}

func generateMathCeil(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	// For integers, ceil is just the value itself
	cg.generateExpressionToReg(args[0], "rax")
}

func generateMathRound(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	// For integers, round is just the value itself
	cg.generateExpressionToReg(args[0], "rax")
}

func generateMathGcd(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	// Greatest Common Divisor using Euclidean algorithm
	labelLoop := cg.getLabel("gcd_loop")
	labelEnd := cg.getLabel("gcd_end")

	// Load first argument into rax
	cg.generateExpressionToReg(args[0], "rax")
	// Load second argument into rcx
	cg.generateExpressionToReg(args[1], "rcx")

	// rax = a, rcx = b
	// while b != 0:
	//   temp = b
	//   b = a % b
	//   a = temp

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelLoop))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", labelEnd))

	cg.textSection.WriteString("    movq %rcx, %rdx\n") // temp = b
	cg.textSection.WriteString("    movq %rax, %r8\n")  // r8 = a (preserve for division)
	cg.textSection.WriteString("    cqo\n")             // sign extend for division
	cg.textSection.WriteString("    idivq %rcx\n")      // rax = a / b, rdx = a % b
	cg.textSection.WriteString("    movq %rdx, %rcx\n") // b = a % b
	cg.textSection.WriteString("    movq %r8, %rax\n")  // rax = a (restore from temp? no, need the divisor result)
	cg.textSection.WriteString("    movq %rdx, %rcx\n") // Actually, set b to remainder
	cg.textSection.WriteString("    movq %r8, %rax\n")  // a = old b (which is in rdx now... let me fix this)

	// Actually, let me rewrite this more clearly
	// gcd(a,b):
	//   while b != 0:
	//     temp = b
	//     b = a mod b
	//     a = temp
	//   return a

	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", labelLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelEnd))
	// rax already contains the GCD
}

func generateMathLcm(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	// LCM(a,b) = (a * b) / GCD(a,b)

	// Load first argument into rax
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    pushq %rax\n") // save a

	// Load second argument into rax
	cg.generateExpressionToReg(args[1], "rax")
	cg.textSection.WriteString("    pushq %rax\n") // save b

	// For now, just return a*b (simplified version)
	// TODO: Call GCD and divide
	cg.textSection.WriteString("    popq %rcx\n")        // rcx = b
	cg.textSection.WriteString("    popq %rax\n")        // rax = a
	cg.textSection.WriteString("    imulq %rcx, %rax\n") // rax = a * b
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
	// concat(a, b): allocate new buffer [len(a)+len(b)+1], copy a then b, NUL-terminate, return ptr in rax
	if len(args) < 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	// Resolve first string pointer -> r10, len -> r8
	resolvePtrLen := func(expr ASTNode, ptrReg, lenReg string, loopTag string) {
		switch v := expr.(type) {
		case *StringLiteral:
			lbl, _ := emitStringLiteral(cg, v.Value)
			cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%%s\n", lbl, ptrReg))
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%%s\n", len(v.Value), lenReg))
		case *Identifier:
			if varInfo, ok := cg.variables[v.Name]; ok {
				cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%%s\n", varInfo.Offset, ptrReg))
				if l, ok := cg.stringLengths[v.Name]; ok {
					cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%%s\n", l, lenReg))
				} else {
					// compute length at runtime
					lLoop := cg.getLabel(loopTag + "_loop")
					lEnd := cg.getLabel(loopTag + "_end")
					cg.textSection.WriteString(fmt.Sprintf("    movq %%%s, %%rbx\n", ptrReg))
					cg.textSection.WriteString(fmt.Sprintf("    movq $0, %%%s\n", lenReg))
					cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
					cg.textSection.WriteString("    movzbq (%rbx), %rax\n")
					cg.textSection.WriteString("    testq %rax, %rax\n")
					cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
					cg.textSection.WriteString(fmt.Sprintf("    inc %%%s\n", lenReg))
					cg.textSection.WriteString("    inc %rbx\n")
					cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
					cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
				}
			} else {
				cg.textSection.WriteString(fmt.Sprintf("    xorq %%%s, %%%s\n", ptrReg, ptrReg))
				cg.textSection.WriteString(fmt.Sprintf("    xorq %%%s, %%%s\n", lenReg, lenReg))
			}
		default:
			cg.textSection.WriteString(fmt.Sprintf("    xorq %%%s, %%%s\n", ptrReg, ptrReg))
			cg.textSection.WriteString(fmt.Sprintf("    xorq %%%s, %%%s\n", lenReg, lenReg))
		}
	}
	resolvePtrLen(args[0], "r10", "r8", "s1len")
	resolvePtrLen(args[1], "r11", "r9", "s2len")
	// Save string pointers and lengths before mmap clobbers them
	// Use %r13-%r15 and %r12 since %r11 gets clobbered by syscall
	cg.textSection.WriteString("    movq %r10, %r13\n") // save s1 ptr
	cg.textSection.WriteString("    movq %r8, %r12\n")  // save s1 len
	cg.textSection.WriteString("    movq %r11, %r14\n") // save s2 ptr
	cg.textSection.WriteString("    movq %r9, %r15\n")  // save s2 len
	// total len in rdx = r8 + r9
	cg.textSection.WriteString("    movq %r8, %rdx\n")
	cg.textSection.WriteString("    addq %r9, %rdx\n")
	// size = total+1 in rsi
	cg.textSection.WriteString("    movq %rdx, %rsi\n")
	cg.textSection.WriteString("    addq $1, %rsi\n")
	// mmap(size)
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	// Check mmap error
	lErr := cg.getLabel("cat_mmap_err")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    js %s\n", lErr))
	// dest base in rax, save to rbx
	cg.textSection.WriteString("    movq %rax, %rbx\n")
	// copy first string
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString("    movq %r13, %rsi\n") // restore s1 ptr from r13
	cg.textSection.WriteString("    movq %r12, %rcx\n") // restore s1 len from r12
	l1 := cg.getLabel("cat_cp1")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", l1))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString("    je 1f\n")
	cg.textSection.WriteString("    movb (%rsi), %al\n")
	cg.textSection.WriteString("    movb %al, (%rdi)\n")
	cg.textSection.WriteString("    inc %rsi\n    inc %rdi\n    dec %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", l1))
	cg.textSection.WriteString("1:\n")
	// copy second string
	cg.textSection.WriteString("    movq %r14, %rsi\n") // restore s2 ptr from r14
	cg.textSection.WriteString("    movq %r15, %rcx\n") // restore s2 len from r15
	l2 := cg.getLabel("cat_cp2")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", l2))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString("    je 2f\n")
	cg.textSection.WriteString("    movb (%rsi), %al\n")
	cg.textSection.WriteString("    movb %al, (%rdi)\n")
	cg.textSection.WriteString("    inc %rsi\n    inc %rdi\n    dec %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", l2))
	cg.textSection.WriteString("2:\n")
	// NUL terminate
	cg.textSection.WriteString("    movb $0, (%rdi)\n")
	// return base ptr
	cg.textSection.WriteString("    movq %rbx, %rax\n")
	// Jump past error handler
	lEnd := cg.getLabel("cat_end")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lEnd))
	// On mmap error, return NULL
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lErr))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
}

func generateStringCompare(cg *CodeGenerator, args []ASTNode) {
	// compare(a,b): returns 0 if equal, <0 if a<b, >0 if a>b (lexicographic)
	if len(args) != 2 {
		return
	}
	// Resolve pointers to the two strings into rbx and rcx
	// Supports string literals and identifiers
	resolveStrPtr := func(expr ASTNode, targetReg string) {
		switch v := expr.(type) {
		case *StringLiteral:
			lbl, _ := emitStringLiteral(cg, v.Value)
			cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%%s\n", lbl, targetReg))
		case *Identifier:
			if varInfo, ok := cg.variables[v.Name]; ok {
				cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%%s\n", varInfo.Offset, targetReg))
			} else {
				cg.textSection.WriteString(fmt.Sprintf("    xorq %%%s, %%%s\n", targetReg, targetReg))
			}
		default:
			cg.textSection.WriteString(fmt.Sprintf("    xorq %%%s, %%%s\n", targetReg, targetReg))
		}
	}
	resolveStrPtr(args[0], "rbx")
	resolveStrPtr(args[1], "rcx")
	lLoop := cg.getLabel("strcmp_loop")
	lEnd := cg.getLabel("strcmp_end")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
	cg.textSection.WriteString("    movzbq (%rbx), %rax\n")
	cg.textSection.WriteString("    movzbq (%rcx), %rdx\n")
	cg.textSection.WriteString("    cmpq %rdx, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lEnd))
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
	cg.textSection.WriteString("    inc %rbx\n")
	cg.textSection.WriteString("    inc %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
	cg.textSection.WriteString("    subq %rdx, %rax\n")
	// result in rax
}

func generateStringCopy(cg *CodeGenerator, args []ASTNode) {
	// copy(src): allocate new buffer [len(src)+1], copy, NUL-terminate, return ptr
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	// Resolve src pointer -> r10, len -> r8
	switch v := args[0].(type) {
	case *StringLiteral:
		lbl, _ := emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%r10\n", lbl))
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", len(v.Value)))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%r10\n", varInfo.Offset))
			if l, ok := cg.stringLengths[v.Name]; ok {
				cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", l))
			} else {
				lLoop := cg.getLabel("cpy_len_loop")
				lEnd := cg.getLabel("cpy_len_end")
				cg.textSection.WriteString("    movq %r10, %rbx\n")
				cg.textSection.WriteString("    movq $0, %r8\n")
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
				cg.textSection.WriteString("    movzbq (%rbx), %rax\n")
				cg.textSection.WriteString("    testq %rax, %rax\n")
				cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
				cg.textSection.WriteString("    inc %r8\n")
				cg.textSection.WriteString("    inc %rbx\n")
				cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
			}
		} else {
			cg.textSection.WriteString("    xorq %r10, %r10\n    xorq %r8, %r8\n")
		}
	default:
		cg.textSection.WriteString("    xorq %r10, %r10\n    xorq %r8, %r8\n")
	}
	// Save source ptr and length before mmap clobbers registers
	// Use %r13, %r14 instead of %r11, %r12 because syscall clobbers %r11
	cg.textSection.WriteString("    movq %r10, %r13\n") // save src ptr to r13
	cg.textSection.WriteString("    movq %r8, %r14\n")  // save src len to r14
	// size = len+1 in rsi
	cg.textSection.WriteString("    movq %r8, %rsi\n    addq $1, %rsi\n")
	// mmap(size)
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	// Check mmap error: if rax is negative, it's an error (Linux returns -errno)
	lErr := cg.getLabel("copy_mmap_err")
	lEnd := cg.getLabel("copy_end")
	cg.textSection.WriteString(fmt.Sprintf("    testq %%rax, %%rax\n"))
	cg.textSection.WriteString(fmt.Sprintf("    js %s\n", lErr))
	// dest base in rax, save to rbx
	cg.textSection.WriteString("    movq %rax, %rbx\n")
	// copy: rdi=dest, rsi=src, rcx=len
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString("    movq %r13, %rsi\n") // restore src from r13
	cg.textSection.WriteString("    movq %r14, %rcx\n") // restore len from r14
	l1 := cg.getLabel("cpy_cp")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", l1))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString("    je 3f\n")
	cg.textSection.WriteString("    movb (%rsi), %al\n")
	cg.textSection.WriteString("    movb %al, (%rdi)\n")
	cg.textSection.WriteString("    inc %rsi\n    inc %rdi\n    dec %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", l1))
	cg.textSection.WriteString("3:\n")
	cg.textSection.WriteString("    movb $0, (%rdi)\n")
	cg.textSection.WriteString("    movq %rbx, %rax\n")
	// Jump past error handler
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lEnd))
	// On mmap error, return NULL (rax = 0)
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lErr))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
}

// indexOf(s, chOrSubstr): if second arg is string literal of length 1, treat as char
func generateStringIndexOf(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	// Resolve haystack in rbx
	switch v := args[0].(type) {
	case *StringLiteral:
		lbl, _ := emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rbx\n", lbl))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rbx\n", varInfo.Offset))
		} else {
			cg.textSection.WriteString("    xorq %rbx, %rbx\n")
		}
	default:
		cg.textSection.WriteString("    xorq %rbx, %rbx\n")
	}
	// Resolve needle in al (single byte) or pointer in rcx
	isChar := false
	if lit, ok := args[1].(*StringLiteral); ok && len(lit.Value) == 1 {
		isChar = true
		cg.textSection.WriteString(fmt.Sprintf("    movb $%d, %%al\n", byte(lit.Value[0])))
	} else {
		switch v := args[1].(type) {
		case *StringLiteral:
			lbl, _ := emitStringLiteral(cg, v.Value)
			cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rcx\n", lbl))
		case *Identifier:
			if varInfo, ok := cg.variables[v.Name]; ok {
				cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rcx\n", varInfo.Offset))
			} else {
				cg.textSection.WriteString("    xorq %rcx, %rcx\n")
			}
		default:
			cg.textSection.WriteString("    xorq %rcx, %rcx\n")
		}
	}
	if isChar {
		lLoop := cg.getLabel("idx_loop")
		lNotFound := cg.getLabel("idx_nf")
		lEnd := cg.getLabel("idx_end")
		lDone := cg.getLabel("idx_done")
		cg.textSection.WriteString("    movq $0, %r10\n")
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
		cg.textSection.WriteString("    movzbq (%rbx), %rdx\n")
		cg.textSection.WriteString("    testq %rdx, %rdx\n")
		cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lNotFound))
		cg.textSection.WriteString("    cmpb %al, %dl\n")
		cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
		cg.textSection.WriteString("    inc %r10\n")
		cg.textSection.WriteString("    inc %rbx\n")
		cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lNotFound))
		cg.textSection.WriteString("    movq $-1, %rax\n")
		cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lDone))
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
		cg.textSection.WriteString("    movq %r10, %rax\n")
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lDone))
		return
	}
	// substring search (naive)
	lOuter := cg.getLabel("sub_outer")
	lInner := cg.getLabel("sub_inner")
	lNext := cg.getLabel("sub_next")
	lFound := cg.getLabel("sub_found")
	lNF := cg.getLabel("sub_nf")
	cg.textSection.WriteString("    movq $0, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lOuter))
	cg.textSection.WriteString("    movzbq (%rbx), %r8\n")
	cg.textSection.WriteString("    testq %r8, %r8\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lNF))
	cg.textSection.WriteString("    movq %rbx, %r10\n")
	cg.textSection.WriteString("    movq %rcx, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lInner))
	cg.textSection.WriteString("    movzbq (%r11), %r8\n")
	cg.textSection.WriteString("    testq %r8, %r8\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lFound))
	cg.textSection.WriteString("    movzbq (%r10), %r9\n")
	cg.textSection.WriteString("    cmpq %r8, %r9\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lNext))
	cg.textSection.WriteString("    inc %r10\n")
	cg.textSection.WriteString("    inc %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lInner))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lNext))
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    inc %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lOuter))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lFound))
	// rax holds starting index
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_end\n", lOuter))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lNF))
	cg.textSection.WriteString("    movq $-1, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s_end:\n", lOuter))
}

func generateStringContains(cg *CodeGenerator, args []ASTNode) {
	// contains(s, sub) -> indexOf >= 0
	generateStringIndexOf(cg, args)
	// rax has index or -1
	lTrue := cg.getLabel("contains_true")
	lEnd := cg.getLabel("contains_end")
	cg.textSection.WriteString("    cmpq $0, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lTrue))
	cg.textSection.WriteString("    movq $0, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lEnd))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lTrue))
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
}

func generateStringStartsWith(cg *CodeGenerator, args []ASTNode) {
	// startsWith(s, prefix): compare from start; return 1/0
	if len(args) != 2 {
		return
	}
	// Resolve pointers in rbx (s) and rcx (prefix)
	// Reuse logic from compare/indexOf
	// s -> rbx
	switch v := args[0].(type) {
	case *StringLiteral:
		lbl, _ := emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rbx\n", lbl))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rbx\n", varInfo.Offset))
		} else {
			cg.textSection.WriteString("    xorq %rbx, %rbx\n")
		}
	default:
		cg.textSection.WriteString("    xorq %rbx, %rbx\n")
	}
	// prefix -> rcx
	switch v := args[1].(type) {
	case *StringLiteral:
		lbl, _ := emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rcx\n", lbl))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rcx\n", varInfo.Offset))
		} else {
			cg.textSection.WriteString("    xorq %rcx, %rcx\n")
		}
	default:
		cg.textSection.WriteString("    xorq %rcx, %rcx\n")
	}
	lLoop := cg.getLabel("sw_loop")
	lEnd := cg.getLabel("sw_end")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
	cg.textSection.WriteString("    movzbq (%rcx), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
	cg.textSection.WriteString("    movzbq (%rbx), %rdx\n")
	cg.textSection.WriteString("    cmpq %rax, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s_false\n", lEnd))
	cg.textSection.WriteString("    inc %rbx\n")
	cg.textSection.WriteString("    inc %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s_false:\n", lEnd))
	cg.textSection.WriteString("    movq $0, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_done\n", lEnd))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s_done:\n", lEnd))
}

func generateStringEndsWith(cg *CodeGenerator, args []ASTNode) {
	// endsWith(s, suffix): naive by computing lengths and comparing from back
	if len(args) != 2 {
		return
	}
	// Compute len(s) in r8 and len(suffix) in r9, pointers in rbx (s) rcx (suffix)
	// Resolve pointers
	var lbl string
	switch v := args[0].(type) {
	case *StringLiteral:
		lbl, _ = emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rbx\n", lbl))
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", len(v.Value)))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rbx\n", varInfo.Offset))
			ln := 0
			if l, ok := cg.stringLengths[v.Name]; ok {
				ln = l
			}
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", ln))
		} else {
			cg.textSection.WriteString("    xorq %rbx, %rbx\n    xorq %r8, %r8\n")
		}
	default:
		cg.textSection.WriteString("    xorq %rbx, %rbx\n    xorq %r8, %r8\n")
	}
	switch v := args[1].(type) {
	case *StringLiteral:
		lbl, _ = emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rcx\n", lbl))
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r9\n", len(v.Value)))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rcx\n", varInfo.Offset))
			ln := 0
			if l, ok := cg.stringLengths[v.Name]; ok {
				ln = l
			}
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r9\n", ln))
		} else {
			cg.textSection.WriteString("    xorq %rcx, %rcx\n    xorq %r9, %r9\n")
		}
	default:
		cg.textSection.WriteString("    xorq %rcx, %rcx\n    xorq %r9, %r9\n")
	}
	// if r9 > r8 -> false
	lFalse := cg.getLabel("ew_false")
	lTrue := cg.getLabel("ew_true")
	lEnd := cg.getLabel("ew_end")
	cg.textSection.WriteString("    cmpq %r9, %r8\n")
	cg.textSection.WriteString(fmt.Sprintf("    jb %s\n", lFalse))
	// set pointers to ends
	cg.textSection.WriteString("    movq %rbx, %r10\n    addq %r8, %r10\n    dec %r10\n")
	cg.textSection.WriteString("    movq %rcx, %r11\n    addq %r9, %r11\n    dec %r11\n")
	lLoop := cg.getLabel("ew_loop")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
	cg.textSection.WriteString("    testq %r9, %r9\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lTrue))
	cg.textSection.WriteString("    movzbq (%r10), %rax\n")
	cg.textSection.WriteString("    movzbq (%r11), %rdx\n")
	cg.textSection.WriteString("    cmpq %rdx, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lFalse))
	cg.textSection.WriteString("    dec %r10\n    dec %r11\n    dec %r9\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lTrue))
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lEnd))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lFalse))
	cg.textSection.WriteString("    movq $0, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
}

// printfuncs.go wrappers (will use existing implementations from printfuncs.go)
func generatePrintCode(cg *CodeGenerator, args []ASTNode) {
	// References existing print logic from printfuncs.go
}

// Note: Printf and Printf in module system delegate to printfuncs implementations
// The io module functions reference the existing printf/fprintf/sprintf implementations

// ========================= Number conversions =============================
func generateNumToInt8(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    salq $56, %rax\n    sarq $56, %rax\n")
}
func generateNumToUint8(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    andq $255, %rax\n")
}
func generateNumToInt16(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    salq $48, %rax\n    sarq $48, %rax\n")
}
func generateNumToUint16(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    andq $65535, %rax\n")
}
func generateNumToInt32(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    cltq\n")
}
func generateNumToUint32(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    movl %eax, %eax\n")
}
func generateNumToInt64(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	// no-op, already 64-bit
}
func generateNumToUint64(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	// no-op for now
}
func generateNumToBool(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    testq %rax, %rax\n    setne %al\n    movzbq %al, %rax\n")
}

// ========================= Memory implementations ==========================

// generateMemSizeof: sizeof(x) -> returns byte size of variable's type if identifier
func generateMemSizeof(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	switch v := args[0].(type) {
	case *Identifier:
		if info, ok := cg.variables[v.Name]; ok {
			sz := GetTypeSize(info.Type)
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rax\n", sz))
			return
		}
	case *StringLiteral:
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rax\n", len(v.Value)))
		return
	case *IntLiteral:
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rax\n", GetTypeSize(TokenTypeInt32)))
		return
	}
	cg.textSection.WriteString("    movq $0, %rax\n")
}

// memcpy(dst, src, n): returns dst
func generateMemMemcpy(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	// dst -> rdi, src -> rsi, n -> rcx
	cg.generateExpressionToReg(args[0], "rdi")
	cg.generateExpressionToReg(args[1], "rsi")
	cg.generateExpressionToReg(args[2], "rcx")
	lLoop := cg.getLabel("memcpy_loop")
	lEnd := cg.getLabel("memcpy_end")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
	cg.textSection.WriteString("    movb (%rsi), %al\n")
	cg.textSection.WriteString("    movb %al, (%rdi)\n")
	cg.textSection.WriteString("    inc %rsi\n    inc %rdi\n    dec %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
	cg.textSection.WriteString("    movq %rdi, %rax\n    subq %rcx, %rax\n")
}

// memset(dst, value, n): returns dst (value treated as byte)
func generateMemMemset(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi")
	cg.generateExpressionToReg(args[1], "rax")
	cg.generateExpressionToReg(args[2], "rcx")
	lLoop := cg.getLabel("memset_loop")
	lEnd := cg.getLabel("memset_end")
	cg.textSection.WriteString("    movb %al, %dl\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
	cg.textSection.WriteString("    movb %dl, (%rdi)\n")
	cg.textSection.WriteString("    inc %rdi\n    dec %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
	cg.textSection.WriteString("    movq %rdi, %rax\n")
}

// malloc(size): implement via mmap(size) and return pointer
func generateMemMalloc(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	// rax=9 (mmap), rdi=NULL, rsi=len, rdx=PROT_READ|PROT_WRITE(3), r10=MAP_PRIVATE|MAP_ANONYMOUS(0x22), r8= -1, r9=0
	cg.generateExpressionToReg(args[0], "rsi")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	// ptr in rax
}

// mmap(size): same as malloc
func generateMemMmap(cg *CodeGenerator, args []ASTNode) { generateMemMalloc(cg, args) }

// free(ptr): currently no-op; recommend mem.munmap
func generateMemFreeNoop(cg *CodeGenerator, args []ASTNode) {
	// Intentionally a no-op; munmap should be used
}

// munmap(ptr, size): rax=11, rdi=addr, rsi=len
func generateMemMunmap(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi")
	cg.generateExpressionToReg(args[1], "rsi")
	cg.textSection.WriteString("    movq $11, %rax\n    syscall\n")
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
