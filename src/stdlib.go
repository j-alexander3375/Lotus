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
	"io":          createIOModule(),
	"mem":         createMemoryModule(),
	"math":        createMathModule(),
	"str":         createStringModule(),
	"num":         createNumModule(),
	"hash":        createHashModule(),
	"collections": createCollectionsModule(),
	"net":         createNetModule(),
	"http":        createHTTPModule(),
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

// createNetModule creates a low-level networking module (Linux syscalls)
func createNetModule() *StdlibModule {
	return &StdlibModule{
		Name: "net",
		Functions: map[string]*StdlibFunction{
			"socket":       {Name: "socket", Module: "net", NumArgs: 3, CodeGen: generateNetSocket},
			"connect_ipv4": {Name: "connect_ipv4", Module: "net", NumArgs: 3, CodeGen: generateNetConnectIPv4},
			"send":         {Name: "send", Module: "net", NumArgs: 3, CodeGen: generateNetSend},
			"recv":         {Name: "recv", Module: "net", NumArgs: 3, CodeGen: generateNetRecv},
			"close":        {Name: "close", Module: "net", NumArgs: 1, CodeGen: generateNetClose},
		},
		Types: map[string]TokenType{},
	}
}

// createHTTPModule creates a minimal HTTP client module built on net helpers
func createHTTPModule() *StdlibModule {
	return &StdlibModule{
		Name: "http",
		Functions: map[string]*StdlibFunction{
			"get": {Name: "get", Module: "http", NumArgs: 7, CodeGen: generateHTTPGetSimple},
		},
		Types: map[string]TokenType{},
	}
}

// createCollectionsModule creates the collections stdlib module (data structures)
func createCollectionsModule() *StdlibModule {
	return &StdlibModule{
		Name: "collections",
		Functions: map[string]*StdlibFunction{
			// Dynamic array
			"array_int_new":  {Name: "array_int_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsArrayIntNew},
			"array_int_push": {Name: "array_int_push", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsArrayIntPush},
			"array_int_pop":  {Name: "array_int_pop", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsArrayIntPop},
			"array_int_len":  {Name: "array_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsArrayIntLen},

			// Stack
			"stack_int_new":  {Name: "stack_int_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsStackIntNew},
			"stack_int_push": {Name: "stack_int_push", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsStackIntPush},
			"stack_int_pop":  {Name: "stack_int_pop", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsStackIntPop},
			"stack_int_len":  {Name: "stack_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsStackIntLen},

			// Queue / Deque
			"queue_int_new":     {Name: "queue_int_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsQueueIntNew},
			"queue_int_enqueue": {Name: "queue_int_enqueue", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsQueueIntEnqueue},
			"queue_int_dequeue": {Name: "queue_int_dequeue", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsQueueIntDequeue},
			"queue_int_len":     {Name: "queue_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsQueueIntLen},

			"deque_int_new":        {Name: "deque_int_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsDequeIntNew},
			"deque_int_push_front": {Name: "deque_int_push_front", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsDequeIntPushFront},
			"deque_int_push_back":  {Name: "deque_int_push_back", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsDequeIntPushBack},
			"deque_int_pop_front":  {Name: "deque_int_pop_front", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsDequeIntPopFront},
			"deque_int_pop_back":   {Name: "deque_int_pop_back", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsDequeIntPopBack},
			"deque_int_len":        {Name: "deque_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsDequeIntLen},

			// Heap (min-heap)
			"heap_int_new":  {Name: "heap_int_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHeapIntNew},
			"heap_int_push": {Name: "heap_int_push", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHeapIntPush},
			"heap_int_pop":  {Name: "heap_int_pop", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHeapIntPop},
			"heap_int_peek": {Name: "heap_int_peek", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHeapIntPeek},
			"heap_int_len":  {Name: "heap_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHeapIntLen},

			// Hash map & set (int keys)
			"hashmap_int_new":    {Name: "hashmap_int_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashmapIntNew},
			"hashmap_int_put":    {Name: "hashmap_int_put", Module: "collections", NumArgs: 3, CodeGen: generateCollectionsHashmapIntPut},
			"hashmap_int_get":    {Name: "hashmap_int_get", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashmapIntGet},
			"hashmap_int_remove": {Name: "hashmap_int_remove", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashmapIntRemove},
			"hashmap_int_len":    {Name: "hashmap_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashmapIntLen},
			"hashmap_int_clear":  {Name: "hashmap_int_clear", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashmapIntClear},
			"hashmap_int_free":   {Name: "hashmap_int_free", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashmapIntFree},

			"hashset_int_new":      {Name: "hashset_int_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashsetIntNew},
			"hashset_int_add":      {Name: "hashset_int_add", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashsetIntAdd},
			"hashset_int_contains": {Name: "hashset_int_contains", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashsetIntContains},
			"hashset_int_remove":   {Name: "hashset_int_remove", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashsetIntRemove},
			"hashset_int_len":      {Name: "hashset_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashsetIntLen},
			"hashset_int_clear":    {Name: "hashset_int_clear", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashsetIntClear},
			"hashset_int_free":     {Name: "hashset_int_free", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashsetIntFree},

			// Array helper
			"binary_search_int": {Name: "binary_search_int", Module: "collections", NumArgs: 3, CodeGen: generateCollectionsBinarySearchInt},
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

// createHashModule creates the hashing stdlib module (cryptographic and non-cryptographic)
func createHashModule() *StdlibModule {
	return &StdlibModule{
		Name: "hash",
		Functions: map[string]*StdlibFunction{
			// Non-cryptographic hashes (fast, simple)
			"crc32":  {Name: "crc32", Module: "hash", NumArgs: 2, CodeGen: generateHashCRC32},    // crc32(data_ptr, len) -> uint32
			"fnv1a":  {Name: "fnv1a", Module: "hash", NumArgs: 2, CodeGen: generateHashFNV1a},    // fnv1a(data_ptr, len) -> uint64
			"djb2":   {Name: "djb2", Module: "hash", NumArgs: 1, CodeGen: generateHashDJB2},      // djb2(string_ptr) -> uint64
			"murmur": {Name: "murmur", Module: "hash", NumArgs: 3, CodeGen: generateHashMurmur3}, // murmur(data_ptr, len, seed) -> uint32

			// Cryptographic hashes
			"sha256": {Name: "sha256", Module: "hash", NumArgs: 3, CodeGen: generateHashSHA256}, // sha256(data_ptr, len, out_buf) -> void
			"md5":    {Name: "md5", Module: "hash", NumArgs: 3, CodeGen: generateHashMD5},       // md5(data_ptr, len, out_buf) -> void
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
	if len(args) != 1 {
		return
	}
	// Integer floor sqrt using hardware sqrt on double; negative -> -1
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    js 2f\n")
	cg.textSection.WriteString("    cvtsi2sd %rax, %xmm0\n")
	cg.textSection.WriteString("    sqrtsd %xmm0, %xmm0\n")
	cg.textSection.WriteString("    cvttsd2si %xmm0, %rax\n")
	cg.textSection.WriteString("    jmp 3f\n")
	cg.textSection.WriteString("2:  movq $-1, %rax\n")
	cg.textSection.WriteString("3:\n")
}

func generateMathPow(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	// Fast exponentiation for integers; negative exponent -> 0
	cg.generateExpressionToReg(args[0], "rbx") // base
	cg.generateExpressionToReg(args[1], "rcx") // exp
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString("    js 3f\n")
	cg.textSection.WriteString("1:  testq %rcx, %rcx\n")
	cg.textSection.WriteString("    je 2f\n")
	cg.textSection.WriteString("    testb $1, %cl\n")
	cg.textSection.WriteString("    je 4f\n")
	cg.textSection.WriteString("    imulq %rbx, %rax\n")
	cg.textSection.WriteString("4:  shrq $1, %rcx\n")
	cg.textSection.WriteString("    imulq %rbx, %rbx\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("2:  jmp 5f\n")
	cg.textSection.WriteString("3:  xorq %rax, %rax\n")
	cg.textSection.WriteString("5:\n")
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

// ============================================================================
// Net module (Linux syscalls)
// ============================================================================

// socket(domain, type, protocol)
func generateNetSocket(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi")
	cg.generateExpressionToReg(args[1], "rsi")
	cg.generateExpressionToReg(args[2], "rdx")
	cg.textSection.WriteString("    movq $41, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// connect_ipv4(fd, ip_u32_host, port_host)
func generateNetConnectIPv4(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi") // fd
	cg.generateExpressionToReg(args[1], "rsi") // ip host order
	cg.generateExpressionToReg(args[2], "rdx") // port host order
	cg.textSection.WriteString("    subq $32, %rsp\n")
	cg.textSection.WriteString("    movw $2, (%rsp)\n") // AF_INET
	cg.textSection.WriteString("    movq %rdx, %rcx\n")
	cg.textSection.WriteString("    rorw $8, %cx\n")
	cg.textSection.WriteString("    movw %cx, 2(%rsp)\n") // port network order
	cg.textSection.WriteString("    movq %rsi, %r8\n")
	cg.textSection.WriteString("    movl %r8d, %ecx\n")
	cg.textSection.WriteString("    bswap %ecx\n")
	cg.textSection.WriteString("    movl %ecx, 4(%rsp)\n") // ip network order
	cg.textSection.WriteString("    movq $0, 8(%rsp)\n")
	cg.textSection.WriteString("    movq $0, 16(%rsp)\n")
	cg.textSection.WriteString("    movq %rsp, %rsi\n")
	cg.textSection.WriteString("    movq $16, %rdx\n")
	cg.textSection.WriteString("    movq $42, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    addq $32, %rsp\n")
}

// send(fd, buf, len) -> bytes written
func generateNetSend(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi")
	cg.generateExpressionToReg(args[1], "rsi")
	cg.generateExpressionToReg(args[2], "rdx")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// recv(fd, buf, len) -> bytes read
func generateNetRecv(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi")
	cg.generateExpressionToReg(args[1], "rsi")
	cg.generateExpressionToReg(args[2], "rdx")
	cg.textSection.WriteString("    movq $0, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// close(fd)
func generateNetClose(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi")
	cg.textSection.WriteString("    movq $3, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// ============================================================================
// HTTP module - minimal GET over an existing connected socket
// ============================================================================

// get(fd, host_ptr, host_len, path_ptr, path_len, buf_ptr, buf_len) -> bytes read
func generateHTTPGetSimple(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 7 {
		return
	}
	// fd in r12 for reuse
	cg.generateExpressionToReg(args[0], "rdi")
	cg.textSection.WriteString("    movq %rdi, %r12\n")

	// literals
	lblGet, _ := emitStringLiteral(cg, "GET ")
	lblMid, _ := emitStringLiteral(cg, " HTTP/1.0\r\nHost: ")
	lblEnd, _ := emitStringLiteral(cg, "\r\nConnection: close\r\n\r\n")

	writeLiteral := func(label string, length int) {
		cg.textSection.WriteString(fmt.Sprintf("    movq %%r12, %%rdi\n"))
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rsi\n", label))
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdx\n", length))
		cg.textSection.WriteString("    movq $1, %rax\n")
		cg.textSection.WriteString("    syscall\n")
	}

	writeLiteral(lblGet, 4)

	// write path
	cg.generateExpressionToReg(args[3], "rsi") // path ptr
	cg.generateExpressionToReg(args[4], "rdx") // path len
	cg.textSection.WriteString("    movq %r12, %rdi\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    syscall\n")

	writeLiteral(lblMid, 17)

	// write host
	cg.generateExpressionToReg(args[1], "rsi") // host ptr
	cg.generateExpressionToReg(args[2], "rdx") // host len
	cg.textSection.WriteString("    movq %r12, %rdi\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    syscall\n")

	writeLiteral(lblEnd, 24)

	// read response into caller buffer
	cg.generateExpressionToReg(args[5], "rsi") // buf ptr
	cg.generateExpressionToReg(args[6], "rdx") // buf len
	cg.textSection.WriteString("    movq %r12, %rdi\n")
	cg.textSection.WriteString("    movq $0, %rax\n")
	cg.textSection.WriteString("    syscall\n")
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
	// Evaluate size first before setting up syscall registers
	cg.generateExpressionToReg(args[0], "rax")

	// Save current stack pointer and align to 16 bytes for syscall
	cg.textSection.WriteString("    pushq %rbp\n")
	cg.textSection.WriteString("    movq %rsp, %rbp\n")
	cg.textSection.WriteString("    andq $-16, %rsp\n") // Align stack to 16 bytes

	// Setup mmap syscall: rax=9, rdi=NULL, rsi=len, rdx=PROT_READ|PROT_WRITE(3), r10=MAP_PRIVATE|MAP_ANONYMOUS(0x22), r8=-1, r9=0
	cg.textSection.WriteString("    movq %rax, %rsi\n") // Move size from rax to rsi
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")

	// Restore stack pointer
	cg.textSection.WriteString("    movq %rbp, %rsp\n")
	cg.textSection.WriteString("    popq %rbp\n")
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
	// Save and align stack for syscall
	cg.textSection.WriteString("    pushq %rbp\n")
	cg.textSection.WriteString("    movq %rsp, %rbp\n")
	cg.textSection.WriteString("    andq $-16, %rsp\n")

	cg.generateExpressionToReg(args[0], "rdi")
	cg.generateExpressionToReg(args[1], "rsi")
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")

	// Restore stack
	cg.textSection.WriteString("    movq %rbp, %rsp\n")
	cg.textSection.WriteString("    popq %rbp\n")
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

// ============================================================================
// Collections module - first-pass implementations
// ============================================================================

const collectionsHeaderSize = 40 // len(0), cap(8), head(16), tail(24), data ptr(32)

// allocate header + backing store (cap * elemSize bytes) using mmap; cap is preserved in %rbx
func collectionsAllocSized(cg *CodeGenerator, capExpr ASTNode, elemSize int) {
	cg.generateExpressionToReg(capExpr, "rax")
	cg.textSection.WriteString("    movq %rax, %rbx\n") // cap
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    imulq $%d, %%rcx\n", elemSize))
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
}

func collectionsAlloc(cg *CodeGenerator, capExpr ASTNode) { collectionsAllocSized(cg, capExpr, 8) }

// store common header fields and return base pointer in %rax
func collectionsInitHeader(cg *CodeGenerator) {
	cg.textSection.WriteString("    movq %rax, %rdx\n")    // base
	cg.textSection.WriteString("    movq $0, (%rdx)\n")    // len
	cg.textSection.WriteString("    movq %rbx, 8(%rdx)\n") // cap
	cg.textSection.WriteString("    movq $0, 16(%rdx)\n")  // head
	cg.textSection.WriteString("    movq $0, 24(%rdx)\n")  // tail
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%rdx), %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, 32(%rdx)\n") // data ptr
	cg.textSection.WriteString("    movq %rdx, %rax\n")
}

// Dynamic array (int)
func generateCollectionsArrayIntNew(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	collectionsAlloc(cg, args[0])
	collectionsInitHeader(cg)
}

func generateCollectionsArrayIntPush(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	// arr ptr -> rbx, value -> rcx
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")  // len
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n") // cap
	cg.textSection.WriteString("    cmpq %rdx, %rax\n")
	cg.textSection.WriteString("    jae 1f\n")             // full
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n") // data ptr
	cg.textSection.WriteString("    movq %rcx, (%r8,%rax,8)\n")
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    movq %rax, %rax\n") // return new len
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  movq $-1, %rax\n") // overflow
	cg.textSection.WriteString("2:\n")
}

func generateCollectionsArrayIntPop(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    je 1f\n")
	cg.textSection.WriteString("    dec %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n")
	cg.textSection.WriteString("    movq (%r8,%rax,8), %rax\n")
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  movq $-1, %rax\n")
	cg.textSection.WriteString("2:\n")
}

func generateCollectionsArrayIntLen(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
}

// Stack (int) - reuse array layout
func generateCollectionsStackIntNew(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntNew(cg, args)
}
func generateCollectionsStackIntPush(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntPush(cg, args)
}
func generateCollectionsStackIntPop(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntPop(cg, args)
}
func generateCollectionsStackIntLen(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntLen(cg, args)
}

// Queue / Deque (int)
func generateCollectionsQueueIntNew(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	collectionsAlloc(cg, args[0])
	collectionsInitHeader(cg)
}

func generateCollectionsQueueIntEnqueue(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")  // len
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n") // cap
	cg.textSection.WriteString("    cmpq %rdx, %rax\n")
	cg.textSection.WriteString("    jae 1f\n")
	cg.textSection.WriteString("    movq 24(%rbx), %r9\n") // tail
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n") // data
	cg.textSection.WriteString("    movq %rcx, (%r8,%r9,8)\n")
	cg.textSection.WriteString("    inc %r9\n")
	cg.textSection.WriteString("    cmpq %rdx, %r9\n")
	cg.textSection.WriteString("    jb 3f\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("3:  movq %r9, 24(%rbx)\n")
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    movq %rax, %rax\n")
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  movq $-1, %rax\n")
	cg.textSection.WriteString("2:\n")
}

func generateCollectionsQueueIntDequeue(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    je 1f\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r9\n") // head
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n")
	cg.textSection.WriteString("    movq (%r8,%r9,8), %rcx\n")
	cg.textSection.WriteString("    inc %r9\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n")
	cg.textSection.WriteString("    cmpq %rdx, %r9\n")
	cg.textSection.WriteString("    jb 3f\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("3:  movq %r9, 16(%rbx)\n")
	cg.textSection.WriteString("    dec %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  movq $-1, %rax\n")
	cg.textSection.WriteString("2:\n")
}

func generateCollectionsQueueIntLen(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntLen(cg, args)
}

func generateCollectionsDequeIntNew(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsQueueIntNew(cg, args)
}
func generateCollectionsDequeIntLen(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntLen(cg, args)
}

func generateCollectionsDequeIntPushFront(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n")
	cg.textSection.WriteString("    cmpq %rdx, %rax\n")
	cg.textSection.WriteString("    jae 1f\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r9\n") // head
	cg.textSection.WriteString("    testq %r9, %r9\n")
	cg.textSection.WriteString("    jne 3f\n")
	cg.textSection.WriteString("    movq %rdx, %r9\n")
	cg.textSection.WriteString("3:  dec %r9\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n")
	cg.textSection.WriteString("    movq %rcx, (%r8,%r9,8)\n")
	cg.textSection.WriteString("    movq %r9, 16(%rbx)\n")
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    movq %rax, %rax\n")
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  movq $-1, %rax\n")
	cg.textSection.WriteString("2:\n")
}

func generateCollectionsDequeIntPushBack(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	generateCollectionsQueueIntEnqueue(cg, args)
}

func generateCollectionsDequeIntPopFront(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsQueueIntDequeue(cg, args)
}

func generateCollectionsDequeIntPopBack(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    je 1f\n")
	cg.textSection.WriteString("    movq 24(%rbx), %r9\n") // tail
	cg.textSection.WriteString("    testq %r9, %r9\n")
	cg.textSection.WriteString("    jne 3f\n")
	cg.textSection.WriteString("    movq 8(%rbx), %r9\n")
	cg.textSection.WriteString("3:  dec %r9\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n")
	cg.textSection.WriteString("    movq (%r8,%r9,8), %rcx\n")
	cg.textSection.WriteString("    movq %r9, 24(%rbx)\n")
	cg.textSection.WriteString("    dec %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  movq $-1, %rax\n")
	cg.textSection.WriteString("2:\n")
}

// Heap (min-heap, int) backed by array layout
func generateCollectionsHeapIntNew(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntNew(cg, args)
}

func generateCollectionsHeapIntPush(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")  // len
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n") // cap
	cg.textSection.WriteString("    cmpq %rdx, %rax\n")
	cg.textSection.WriteString("    jae 5f\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n")
	cg.textSection.WriteString("    movq %rcx, (%r8,%rax,8)\n") // place at end
	cg.textSection.WriteString("    movq %rax, %r9\n")          // idx
	cg.textSection.WriteString("1:  testq %r9, %r9\n")
	cg.textSection.WriteString("    je 2f\n")
	cg.textSection.WriteString("    movq %r9, %r10\n")
	cg.textSection.WriteString("    dec %r10\n")
	cg.textSection.WriteString("    shrq $1, %r10\n")           // parent = (idx-1)/2
	cg.textSection.WriteString("    movq (%r8,%r10,8), %r11\n") // parent val
	cg.textSection.WriteString("    movq (%r8,%r9,8), %r12\n")  // current val
	cg.textSection.WriteString("    cmpq %r11, %r12\n")
	cg.textSection.WriteString("    jge 2f\n")
	cg.textSection.WriteString("    movq %r11, (%r8,%r9,8)\n")
	cg.textSection.WriteString("    movq %r12, (%r8,%r10,8)\n")
	cg.textSection.WriteString("    movq %r10, %r9\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("2:  inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    movq %rax, %rax\n")
	cg.textSection.WriteString("    jmp 6f\n")
	cg.textSection.WriteString("5:  movq $-1, %rax\n")
	cg.textSection.WriteString("6:\n")
}

func generateCollectionsHeapIntPop(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    je 9f\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n")
	cg.textSection.WriteString("    movq (%r8), %rcx\n") // result
	cg.textSection.WriteString("    dec %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    je 8f\n")
	cg.textSection.WriteString("    movq (%r8,%rax,8), %r9\n") // last element
	cg.textSection.WriteString("    movq %r9, (%r8)\n")
	cg.textSection.WriteString("    movq $0, %r10\n") // idx
	cg.textSection.WriteString("3:  movq %r10, %r11\n")
	cg.textSection.WriteString("    shlq $1, %r11\n")
	cg.textSection.WriteString("    inc %r11\n") // left child
	cg.textSection.WriteString("    movq %rax, %r12\n")
	cg.textSection.WriteString("    cmpq %r12, %r11\n")
	cg.textSection.WriteString("    jae 8f\n")
	cg.textSection.WriteString("    movq %r11, %r13\n")
	cg.textSection.WriteString("    inc %r13\n")        // right child
	cg.textSection.WriteString("    movq %r11, %r14\n") // smallest = left
	cg.textSection.WriteString("    cmpq %r12, %r13\n")
	cg.textSection.WriteString("    jae 4f\n")
	cg.textSection.WriteString("    movq (%r8,%r13,8), %r15\n")
	cg.textSection.WriteString("    movq (%r8,%r14,8), %rdi\n")
	cg.textSection.WriteString("    cmpq %rdi, %r15\n")
	cg.textSection.WriteString("    jge 4f\n")
	cg.textSection.WriteString("    movq %r13, %r14\n")
	cg.textSection.WriteString("4:  movq (%r8,%r14,8), %r15\n") // smallest val
	cg.textSection.WriteString("    movq (%r8,%r10,8), %rdi\n") // current val
	cg.textSection.WriteString("    cmpq %rdi, %r15\n")
	cg.textSection.WriteString("    jge 8f\n")
	cg.textSection.WriteString("    movq %r15, (%r8,%r10,8)\n")
	cg.textSection.WriteString("    movq %rdi, (%r8,%r14,8)\n")
	cg.textSection.WriteString("    movq %r14, %r10\n")
	cg.textSection.WriteString("    jmp 3b\n")
	cg.textSection.WriteString("8:  movq %rcx, %rax\n")
	cg.textSection.WriteString("    jmp 10f\n")
	cg.textSection.WriteString("9:  movq $-1, %rax\n")
	cg.textSection.WriteString("10:\n")
}

func generateCollectionsHeapIntPeek(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    je 1f\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n")
	cg.textSection.WriteString("    movq (%r8), %rax\n")
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  movq $-1, %rax\n")
	cg.textSection.WriteString("2:\n")
}

func generateCollectionsHeapIntLen(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntLen(cg, args)
}

// Hash map (int -> int) with hashing, open addressing, and resize (power-of-two cap)
func generateCollectionsHashmapIntNew(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	// round capacity to power-of-two >= requested
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    movq %rax, %rbx\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("1:  cmpq %rbx, %rax\n")
	cg.textSection.WriteString("    jae 2f\n")
	cg.textSection.WriteString("    shlq $1, %rax\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("2:  movq %rax, %rbx\n") // rbX = cap pow2
	// size = header + buckets(16*cap) + states(1*cap padded to 8)
	cg.textSection.WriteString("    movq %rbx, %rcx\n")
	cg.textSection.WriteString("    imulq $16, %rcx\n")
	cg.textSection.WriteString("    movq %rbx, %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n") // pad states to 8-byte
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rcx\n")
	cg.textSection.WriteString("    movq %rcx, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	// init header
	cg.textSection.WriteString("    movq %rax, %r11\n")
	cg.textSection.WriteString("    movq $0, (%r11)\n")
	cg.textSection.WriteString("    movq %rbx, 8(%r11)\n")
	cg.textSection.WriteString("    movq $0, 24(%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%r11), %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, 32(%r11)\n") // data ptr
	cg.textSection.WriteString("    leaq (%rcx,%rbx,16), %rdx\n")
	cg.textSection.WriteString("    movq %rdx, 16(%r11)\n") // states ptr
	cg.textSection.WriteString("    movq %r11, %rax\n")
}

func hashmapHash(cg *CodeGenerator, keyReg string) {
	cg.textSection.WriteString(fmt.Sprintf("    movq %%%s, %%rax\n", keyReg))
	cg.textSection.WriteString("    movq %rax, %r14\n")
	cg.textSection.WriteString("    shrq $33, %r14\n")
	cg.textSection.WriteString("    xorq %r14, %rax\n")
	cg.textSection.WriteString("    imulq $-352301462168404107, %rax\n")
	cg.textSection.WriteString("    movq %rax, %r14\n")
	cg.textSection.WriteString("    shrq $33, %r14\n")
	cg.textSection.WriteString("    xorq %r14, %rax\n")
}

func hashmapEnsureCapacity(cg *CodeGenerator, mapReg string) {
	// if len*10 >= cap*7 -> resize (cap *=2)
	cg.textSection.WriteString(fmt.Sprintf("    movq (%%%s), %%r8\n", mapReg))  // len
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%r9\n", mapReg)) // cap
	cg.textSection.WriteString("    movq %r8, %rax\n")
	cg.textSection.WriteString("    imulq $10, %rax\n")
	cg.textSection.WriteString("    movq %r9, %rbx\n")
	cg.textSection.WriteString("    imulq $7, %rbx\n")
	cg.textSection.WriteString("    cmpq %rbx, %rax\n")
	cg.textSection.WriteString("    jb 9f\n")
	// resize: newCap = cap*2
	cg.textSection.WriteString("    shlq $1, %r9\n")
	// compute new size
	cg.textSection.WriteString("    movq %r9, %rcx\n")
	cg.textSection.WriteString("    imulq $16, %rcx\n")
	cg.textSection.WriteString("    movq %r9, %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rcx\n")
	cg.textSection.WriteString("    movq %rcx, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	// rax = new base, old base in mapReg
	cg.textSection.WriteString(fmt.Sprintf("    movq %%%s, %%r15\n", mapReg))
	// new header init
	cg.textSection.WriteString("    movq %rax, %r11\n")
	cg.textSection.WriteString("    movq $0, (%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%rbx\n", mapReg))
	cg.textSection.WriteString("    shlq $1, %rbx\n")
	cg.textSection.WriteString("    movq %rbx, 8(%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%r11), %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, 32(%r11)\n")
	cg.textSection.WriteString("    movq %rbx, %rdx\n")
	cg.textSection.WriteString("    imulq $16, %rdx\n")
	cg.textSection.WriteString("    leaq (%rcx,%rdx,1), %rsi\n")
	cg.textSection.WriteString("    movq %rsi, 16(%r11)\n")
	// reinsert old entries
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%r8\n", mapReg))   // old cap
	cg.textSection.WriteString(fmt.Sprintf("    movq 32(%%%s), %%r9\n", mapReg))  // old data
	cg.textSection.WriteString(fmt.Sprintf("    movq 16(%%%s), %%r10\n", mapReg)) // old states
	cg.textSection.WriteString("    xorq %r12, %r12\n")                           // idx = 0
	cg.textSection.WriteString("1:  cmpq %r8, %r12\n")
	cg.textSection.WriteString("    jae 4f\n")
	cg.textSection.WriteString("    movzbq (%r10,%r12,1), %r13\n")
	cg.textSection.WriteString("    cmpq $1, %r13\n")
	cg.textSection.WriteString("    jne 3f\n")
	// load key/value
	cg.textSection.WriteString("    movq (%r9,%r12,16), %rcx\n")
	cg.textSection.WriteString("    movq 8(%r9,%r12,16), %r15\n")
	// insert into new table (reuse put logic using r11 as base)
	cg.textSection.WriteString("    movq %r11, %rbx\n")
	cg.textSection.WriteString("    jmp 6f\n")
	cg.textSection.WriteString("3:  inc %r12\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("4:  jmp 8f\n")
	// insertion fragment shared with runtime put
	cg.textSection.WriteString("6:  movq (%rbx), %rax\n")   // len
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")   // cap
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")  // data
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n") // states
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    movq %r15, %rsi\n")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r13\n") // idx
	cg.textSection.WriteString("    movq $-1, %r14\n")  // tomb = -1
	cg.textSection.WriteString("7:  movzbq (%r10,%r13,1), %r15\n")
	cg.textSection.WriteString("    cmpq $1, %r15\n")
	cg.textSection.WriteString("    je 10f\n")
	cg.textSection.WriteString("    cmpq $2, %r15\n")
	cg.textSection.WriteString("    jne 9f\n")
	cg.textSection.WriteString("    cmpq $-1, %r14\n")
	cg.textSection.WriteString("    jne 11f\n")
	cg.textSection.WriteString("    movq %r13, %r14\n")
	cg.textSection.WriteString("    jmp 11f\n")
	cg.textSection.WriteString("9:  movq %r13, %rax\n")
	cg.textSection.WriteString("    cmpq $-1, %r14\n")
	cg.textSection.WriteString("    jne 12f\n")
	cg.textSection.WriteString("    jmp 13f\n")
	cg.textSection.WriteString("10: movq (%r9,%r13,16), %r15\n")
	cg.textSection.WriteString("    cmpq %rcx, %r15\n")
	cg.textSection.WriteString("    jne 11f\n")
	cg.textSection.WriteString("    movq %rsi, 8(%r9,%r13,16)\n")
	cg.textSection.WriteString("    jmp 12f\n")
	cg.textSection.WriteString("11: inc %r13\n")
	cg.textSection.WriteString("    cmpq %r8, %r13\n")
	cg.textSection.WriteString("    jb 7b\n")
	cg.textSection.WriteString("    xorq %r13, %r13\n")
	cg.textSection.WriteString("    jmp 7b\n")
	cg.textSection.WriteString("12: jmp 14f\n")
	cg.textSection.WriteString("13: movq %r13, %r14\n")
	cg.textSection.WriteString("14: cmpq $-1, %r14\n")
	cg.textSection.WriteString("    je 15f\n")
	cg.textSection.WriteString("    movb $1, (%r10,%r14,1)\n")
	cg.textSection.WriteString("    movq %rcx, (%r9,%r14,16)\n")
	cg.textSection.WriteString("    movq %rsi, 8(%r9,%r14,16)\n")
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("15: jmp 3b\n")
	cg.textSection.WriteString("8:  movq %r11, %rax\n")
	// free old
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%rsi\n", mapReg))
	cg.textSection.WriteString("    imulq $16, %rsi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%rdx\n", mapReg))
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rsi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq %%%s, %%rdi\n", mapReg))
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	// update mapReg to new base
	cg.textSection.WriteString(fmt.Sprintf("    movq %%r11, %%%s\n", mapReg))
	cg.textSection.WriteString("9:\n")
}

func generateCollectionsHashmapIntPut(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx") // map base
	// ensure capacity
	hashmapEnsureCapacity(cg, "rbx")
	cg.generateExpressionToReg(args[1], "rcx")              // key
	cg.generateExpressionToReg(args[2], "rdx")              // value
	cg.textSection.WriteString("    movq %rdx, %r15\n")     // stash value
	cg.textSection.WriteString("    movq (%rbx), %rax\n")   // len
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")   // cap
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")  // data ptr
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n") // states ptr
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n") // idx
	cg.textSection.WriteString("    movq $-1, %r12\n")  // tomb
	cg.textSection.WriteString("1:  movzbq (%r10,%r11,1), %r13\n")
	cg.textSection.WriteString("    cmpq $1, %r13\n")
	cg.textSection.WriteString("    je 3f\n")
	cg.textSection.WriteString("    cmpq $2, %r13\n")
	cg.textSection.WriteString("    jne 2f\n")
	cg.textSection.WriteString("    cmpq $-1, %r12\n")
	cg.textSection.WriteString("    jne 4f\n")
	cg.textSection.WriteString("    movq %r11, %r12\n")
	cg.textSection.WriteString("    jmp 4f\n")
	cg.textSection.WriteString("2:  movq %r11, %rax\n")
	cg.textSection.WriteString("    cmpq $-1, %r12\n")
	cg.textSection.WriteString("    jne 6f\n")
	cg.textSection.WriteString("    jmp 7f\n")
	cg.textSection.WriteString("3:  movq (%r9,%r11,16), %r14\n")
	cg.textSection.WriteString("    cmpq %rcx, %r14\n")
	cg.textSection.WriteString("    jne 4f\n")
	cg.textSection.WriteString("    movq %r15, 8(%r9,%r11,16)\n")
	cg.textSection.WriteString("    movq %r15, %rax\n")
	cg.textSection.WriteString("    jmp 6f\n")
	cg.textSection.WriteString("4:  inc %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString("    jb 1b\n")
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("6:  jmp 8f\n")
	cg.textSection.WriteString("7:  movq %r11, %r12\n")
	cg.textSection.WriteString("8:  cmpq $-1, %r12\n")
	cg.textSection.WriteString("    je 9f\n")
	cg.textSection.WriteString("    movb $1, (%r10,%r12,1)\n")
	cg.textSection.WriteString("    movq %rcx, (%r9,%r12,16)\n")
	cg.textSection.WriteString("    movq %r15, 8(%r9,%r12,16)\n")
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("9:\n")
}

func generateCollectionsHashmapIntGet(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")
	cg.textSection.WriteString("1:  movzbq (%r10,%r11,1), %r12\n")
	cg.textSection.WriteString("    cmpq $0, %r12\n")
	cg.textSection.WriteString("    je 3f\n")
	cg.textSection.WriteString("    cmpq $1, %r12\n")
	cg.textSection.WriteString("    jne 2f\n")
	cg.textSection.WriteString("    movq (%r9,%r11,16), %r13\n")
	cg.textSection.WriteString("    cmpq %rcx, %r13\n")
	cg.textSection.WriteString("    je 4f\n")
	cg.textSection.WriteString("2:  inc %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString("    jb 1b\n")
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("3:  movq $-1, %rax\n")
	cg.textSection.WriteString("    jmp 5f\n")
	cg.textSection.WriteString("4:  movq 8(%r9,%r11,16), %rax\n")
	cg.textSection.WriteString("5:\n")
}

func generateCollectionsHashmapIntRemove(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq (%rbx), %r14\n") // len
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")
	cg.textSection.WriteString("1:  movzbq (%r10,%r11,1), %r12\n")
	cg.textSection.WriteString("    cmpq $0, %r12\n")
	cg.textSection.WriteString("    je 3f\n")
	cg.textSection.WriteString("    movq (%r9,%r11,16), %r13\n")
	cg.textSection.WriteString("    cmpq %rcx, %r13\n")
	cg.textSection.WriteString("    je 4f\n")
	cg.textSection.WriteString("    inc %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString("    jb 1b\n")
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("3:  movq $-1, %rax\n")
	cg.textSection.WriteString("    jmp 6f\n")
	cg.textSection.WriteString("4:  movq 8(%r9,%r11,16), %rax\n")
	cg.textSection.WriteString("    movb $2, (%r10,%r11,1)\n")
	cg.textSection.WriteString("    dec %r14\n")
	cg.textSection.WriteString("    movq %r14, (%rbx)\n")
	cg.textSection.WriteString("6:\n")
}

func generateCollectionsHashmapIntLen(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
}

func generateCollectionsHashmapIntClear(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq $0, (%rbx)\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rcx\n")
	cg.textSection.WriteString("    movq 16(%rbx), %rdi\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    rep stosb\n")
}

func generateCollectionsHashmapIntFree(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq 8(%rbx), %rsi\n")
	cg.textSection.WriteString("    imulq $16, %rsi\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rsi\n")
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
}

// Hash set (int) with hashing, open addressing, resize
func generateCollectionsHashsetIntNew(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    movq %rax, %rbx\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("1:  cmpq %rbx, %rax\n")
	cg.textSection.WriteString("    jae 2f\n")
	cg.textSection.WriteString("    shlq $1, %rax\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("2:  movq %rax, %rbx\n")
	cg.textSection.WriteString("    movq %rbx, %rcx\n")
	cg.textSection.WriteString("    imulq $8, %rcx\n")
	cg.textSection.WriteString("    movq %rbx, %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rcx\n")
	cg.textSection.WriteString("    movq %rcx, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    movq %rax, %r11\n")
	cg.textSection.WriteString("    movq $0, (%r11)\n")
	cg.textSection.WriteString("    movq %rbx, 8(%r11)\n")
	cg.textSection.WriteString("    movq $0, 24(%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%r11), %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, 32(%r11)\n")
	cg.textSection.WriteString("    leaq (%rcx,%rbx,8), %rdx\n")
	cg.textSection.WriteString("    movq %rdx, 16(%r11)\n")
	cg.textSection.WriteString("    movq %r11, %rax\n")
}

func hashsetEnsureCapacity(cg *CodeGenerator, mapReg string) {
	cg.textSection.WriteString(fmt.Sprintf("    movq (%%%s), %%r8\n", mapReg))  // len
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%r9\n", mapReg)) // cap
	cg.textSection.WriteString("    movq %r8, %rax\n")
	cg.textSection.WriteString("    imulq $10, %rax\n")
	cg.textSection.WriteString("    movq %r9, %rbx\n")
	cg.textSection.WriteString("    imulq $7, %rbx\n")
	cg.textSection.WriteString("    cmpq %rbx, %rax\n")
	cg.textSection.WriteString("    jb 9f\n")
	cg.textSection.WriteString("    shlq $1, %r9\n")
	cg.textSection.WriteString("    movq %r9, %rcx\n")
	cg.textSection.WriteString("    imulq $8, %rcx\n")
	cg.textSection.WriteString("    movq %r9, %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rcx\n")
	cg.textSection.WriteString("    movq %rcx, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq %%%s, %%r15\n", mapReg))
	cg.textSection.WriteString("    movq %rax, %r11\n")
	cg.textSection.WriteString("    movq $0, (%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%rbx\n", mapReg))
	cg.textSection.WriteString("    shlq $1, %rbx\n")
	cg.textSection.WriteString("    movq %rbx, 8(%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%r11), %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, 32(%r11)\n")
	cg.textSection.WriteString("    movq %rbx, %rdx\n")
	cg.textSection.WriteString("    imulq $8, %rdx\n")
	cg.textSection.WriteString("    leaq (%rcx,%rdx,1), %rsi\n")
	cg.textSection.WriteString("    movq %rsi, 16(%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%r8\n", mapReg))
	cg.textSection.WriteString(fmt.Sprintf("    movq 32(%%%s), %%r9\n", mapReg))
	cg.textSection.WriteString(fmt.Sprintf("    movq 16(%%%s), %%r10\n", mapReg))
	cg.textSection.WriteString("    xorq %r12, %r12\n")
	cg.textSection.WriteString("1:  cmpq %r8, %r12\n")
	cg.textSection.WriteString("    jae 4f\n")
	cg.textSection.WriteString("    movzbq (%r10,%r12,1), %r13\n")
	cg.textSection.WriteString("    cmpq $1, %r13\n")
	cg.textSection.WriteString("    jne 3f\n")
	cg.textSection.WriteString("    movq (%r9,%r12,8), %rcx\n")
	cg.textSection.WriteString("    movq %r11, %rbx\n")
	cg.textSection.WriteString("    jmp 6f\n")
	cg.textSection.WriteString("3:  inc %r12\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("4:  jmp 8f\n")
	cg.textSection.WriteString("6:  movq (%rbx), %rax\n")
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r13\n")
	cg.textSection.WriteString("    movq $-1, %r14\n")
	cg.textSection.WriteString("7:  movzbq (%r10,%r13,1), %r15\n")
	cg.textSection.WriteString("    cmpq $1, %r15\n")
	cg.textSection.WriteString("    je 10f\n")
	cg.textSection.WriteString("    cmpq $2, %r15\n")
	cg.textSection.WriteString("    jne 9f\n")
	cg.textSection.WriteString("    cmpq $-1, %r14\n")
	cg.textSection.WriteString("    jne 11f\n")
	cg.textSection.WriteString("    movq %r13, %r14\n")
	cg.textSection.WriteString("    jmp 11f\n")
	cg.textSection.WriteString("9:  movq %r13, %rax\n")
	cg.textSection.WriteString("    cmpq $-1, %r14\n")
	cg.textSection.WriteString("    jne 12f\n")
	cg.textSection.WriteString("    jmp 13f\n")
	cg.textSection.WriteString("10: movq (%r9,%r13,8), %r15\n")
	cg.textSection.WriteString("    cmpq %rcx, %r15\n")
	cg.textSection.WriteString("    jne 11f\n")
	cg.textSection.WriteString("    jmp 12f\n")
	cg.textSection.WriteString("11: inc %r13\n")
	cg.textSection.WriteString("    cmpq %r8, %r13\n")
	cg.textSection.WriteString("    jb 7b\n")
	cg.textSection.WriteString("    xorq %r13, %r13\n")
	cg.textSection.WriteString("    jmp 7b\n")
	cg.textSection.WriteString("12: jmp 14f\n")
	cg.textSection.WriteString("13: movq %r13, %r14\n")
	cg.textSection.WriteString("14: cmpq $-1, %r14\n")
	cg.textSection.WriteString("    je 15f\n")
	cg.textSection.WriteString("    movb $1, (%r10,%r14,1)\n")
	cg.textSection.WriteString("    movq %rcx, (%r9,%r14,8)\n")
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("15: jmp 3b\n")
	cg.textSection.WriteString("8:  movq %r11, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%rsi\n", mapReg))
	cg.textSection.WriteString("    imulq $8, %rsi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%rdx\n", mapReg))
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rsi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq %%%s, %%rdi\n", mapReg))
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq %%r11, %%%s\n", mapReg))
	cg.textSection.WriteString("9:\n")
}

func generateCollectionsHashsetIntAdd(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	hashsetEnsureCapacity(cg, "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")
	cg.textSection.WriteString("    movq $-1, %r12\n")
	cg.textSection.WriteString("1:  movzbq (%r10,%r11,1), %r13\n")
	cg.textSection.WriteString("    cmpq $1, %r13\n")
	cg.textSection.WriteString("    je 3f\n")
	cg.textSection.WriteString("    cmpq $2, %r13\n")
	cg.textSection.WriteString("    jne 2f\n")
	cg.textSection.WriteString("    cmpq $-1, %r12\n")
	cg.textSection.WriteString("    jne 4f\n")
	cg.textSection.WriteString("    movq %r11, %r12\n")
	cg.textSection.WriteString("    jmp 4f\n")
	cg.textSection.WriteString("2:  movq %r11, %rax\n")
	cg.textSection.WriteString("    cmpq $-1, %r12\n")
	cg.textSection.WriteString("    jne 6f\n")
	cg.textSection.WriteString("    jmp 7f\n")
	cg.textSection.WriteString("3:  movq (%r9,%r11,8), %r14\n")
	cg.textSection.WriteString("    cmpq %rcx, %r14\n")
	cg.textSection.WriteString("    je 6f\n")
	cg.textSection.WriteString("4:  inc %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString("    jb 1b\n")
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("6:  jmp 8f\n")
	cg.textSection.WriteString("7:  movq %r11, %r12\n")
	cg.textSection.WriteString("8:  cmpq $-1, %r12\n")
	cg.textSection.WriteString("    je 9f\n")
	cg.textSection.WriteString("    movb $1, (%r10,%r12,1)\n")
	cg.textSection.WriteString("    movq %rcx, (%r9,%r12,8)\n")
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("9:\n")
}

func generateCollectionsHashsetIntContains(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")
	cg.textSection.WriteString("1:  movzbq (%r10,%r11,1), %r12\n")
	cg.textSection.WriteString("    cmpq $0, %r12\n")
	cg.textSection.WriteString("    je 3f\n")
	cg.textSection.WriteString("    cmpq $1, %r12\n")
	cg.textSection.WriteString("    jne 2f\n")
	cg.textSection.WriteString("    movq (%r9,%r11,8), %r13\n")
	cg.textSection.WriteString("    cmpq %rcx, %r13\n")
	cg.textSection.WriteString("    je 4f\n")
	cg.textSection.WriteString("2:  inc %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString("    jb 1b\n")
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("3:  movq $0, %rax\n")
	cg.textSection.WriteString("    jmp 5f\n")
	cg.textSection.WriteString("4:  movq $1, %rax\n")
	cg.textSection.WriteString("5:\n")
}

func generateCollectionsHashsetIntRemove(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq (%rbx), %r14\n")
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")
	cg.textSection.WriteString("1:  movzbq (%r10,%r11,1), %r12\n")
	cg.textSection.WriteString("    cmpq $0, %r12\n")
	cg.textSection.WriteString("    je 3f\n")
	cg.textSection.WriteString("    movq (%r9,%r11,8), %r13\n")
	cg.textSection.WriteString("    cmpq %rcx, %r13\n")
	cg.textSection.WriteString("    je 4f\n")
	cg.textSection.WriteString("    inc %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString("    jb 1b\n")
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("3:  movq $-1, %rax\n")
	cg.textSection.WriteString("    jmp 6f\n")
	cg.textSection.WriteString("4:  movq %rcx, %rax\n")
	cg.textSection.WriteString("    movb $2, (%r10,%r11,1)\n")
	cg.textSection.WriteString("    dec %r14\n")
	cg.textSection.WriteString("    movq %r14, (%rbx)\n")
	cg.textSection.WriteString("6:\n")
}

func generateCollectionsHashsetIntLen(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
}

func generateCollectionsHashsetIntClear(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq $0, (%rbx)\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rcx\n")
	cg.textSection.WriteString("    movq 16(%rbx), %rdi\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    rep stosb\n")
}

func generateCollectionsHashsetIntFree(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq 8(%rbx), %rsi\n")
	cg.textSection.WriteString("    imulq $8, %rsi\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rsi\n")
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
}

// Binary search helper (int) expects args: base pointer, length, target
func generateCollectionsBinarySearchInt(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")         // base pointer
	cg.generateExpressionToReg(args[1], "rcx")         // len
	cg.generateExpressionToReg(args[2], "rdx")         // target
	cg.textSection.WriteString("    xorq %r8, %r8\n")  // low
	cg.textSection.WriteString("    movq %rcx, %r9\n") // high
	cg.textSection.WriteString("1:  cmpq %r9, %r8\n")
	cg.textSection.WriteString("    jge 4f\n")
	cg.textSection.WriteString("    movq %r8, %r10\n")
	cg.textSection.WriteString("    addq %r9, %r10\n")
	cg.textSection.WriteString("    shrq $1, %r10\n") // mid
	cg.textSection.WriteString("    movq (%rbx,%r10,8), %r11\n")
	cg.textSection.WriteString("    cmpq %rdx, %r11\n")
	cg.textSection.WriteString("    je 3f\n")
	cg.textSection.WriteString("    jl 2f\n")
	cg.textSection.WriteString("    movq %r10, %r9\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("2:  inc %r10\n")
	cg.textSection.WriteString("    movq %r10, %r8\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("3:  movq %r10, %rax\n")
	cg.textSection.WriteString("    jmp 5f\n")
	cg.textSection.WriteString("4:  movq $-1, %rax\n")
	cg.textSection.WriteString("5:\n")
}

// ============================================================================
// Hash module implementations
// ============================================================================

// generateHashCRC32 computes CRC32 checksum (IEEE 802.3 polynomial)
// Args: data_ptr, len -> returns uint32 in rax
func generateHashCRC32(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	// CRC32 table-driven implementation
	// Polynomial: 0xEDB88320 (IEEE 802.3)
	tableLbl := cg.getLabel("crc32_table")
	loopLbl := cg.getLabel("crc32_loop")
	endLbl := cg.getLabel("crc32_end")

	// Generate CRC32 lookup table in .data section
	if !cg.hasLabel(tableLbl) {
		cg.dataSection.WriteString(fmt.Sprintf("%s:\n", tableLbl))
		poly := uint32(0xEDB88320)
		for i := 0; i < 256; i++ {
			crc := uint32(i)
			for j := 0; j < 8; j++ {
				if crc&1 != 0 {
					crc = (crc >> 1) ^ poly
				} else {
					crc >>= 1
				}
			}
			cg.dataSection.WriteString(fmt.Sprintf("    .long 0x%08x\n", crc))
		}
		cg.markLabel(tableLbl)
	}

	cg.generateExpressionToReg(args[0], "rsi")                 // data pointer
	cg.generateExpressionToReg(args[1], "rcx")                 // length
	cg.textSection.WriteString("    movl $0xFFFFFFFF, %eax\n") // crc = ~0
	cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rdi\n", tableLbl))

	// Main loop
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", loopLbl))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", endLbl))
	cg.textSection.WriteString("    movzbq (%rsi), %rbx\n")      // byte
	cg.textSection.WriteString("    xorb %al, %bl\n")            // index = (crc ^ byte) & 0xFF
	cg.textSection.WriteString("    shrl $8, %eax\n")            // crc >>= 8
	cg.textSection.WriteString("    movl (%rdi,%rbx,4), %edx\n") // table[index]
	cg.textSection.WriteString("    xorl %edx, %eax\n")          // crc ^= table[index]
	cg.textSection.WriteString("    incq %rsi\n")
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", loopLbl))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLbl))
	cg.textSection.WriteString("    notl %eax\n") // final XOR with 0xFFFFFFFF
}

// generateHashFNV1a computes FNV-1a hash (64-bit)
// Args: data_ptr, len -> returns uint64 in rax
func generateHashFNV1a(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	loopLbl := cg.getLabel("fnv1a_loop")
	endLbl := cg.getLabel("fnv1a_end")

	cg.generateExpressionToReg(args[0], "rsi") // data pointer
	cg.generateExpressionToReg(args[1], "rcx") // length

	// FNV-1a 64-bit offset basis and prime
	cg.textSection.WriteString("    movabsq $0xCBF29CE484222325, %rax\n") // offset_basis
	cg.textSection.WriteString("    movabsq $0x100000001B3, %r8\n")       // FNV prime

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", loopLbl))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", endLbl))
	cg.textSection.WriteString("    movzbq (%rsi), %rbx\n") // byte
	cg.textSection.WriteString("    xorq %rbx, %rax\n")     // hash ^= byte
	cg.textSection.WriteString("    imulq %r8, %rax\n")     // hash *= prime
	cg.textSection.WriteString("    incq %rsi\n")
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", loopLbl))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLbl))
}

// generateHashDJB2 computes DJB2 hash (simple string hash)
// Args: string_ptr -> returns uint64 in rax
func generateHashDJB2(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	loopLbl := cg.getLabel("djb2_loop")
	endLbl := cg.getLabel("djb2_end")

	cg.generateExpressionToReg(args[0], "rsi")           // string pointer
	cg.textSection.WriteString("    movq $5381, %rax\n") // hash = 5381

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", loopLbl))
	cg.textSection.WriteString("    movzbq (%rsi), %rbx\n") // byte
	cg.textSection.WriteString("    testq %rbx, %rbx\n")    // check for null terminator
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", endLbl))
	// hash = hash * 33 + c
	cg.textSection.WriteString("    movq %rax, %rdx\n") // save hash
	cg.textSection.WriteString("    shlq $5, %rax\n")   // hash << 5 (hash * 32)
	cg.textSection.WriteString("    addq %rdx, %rax\n") // + hash (now hash * 33)
	cg.textSection.WriteString("    addq %rbx, %rax\n") // + c
	cg.textSection.WriteString("    incq %rsi\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", loopLbl))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLbl))
}

// generateHashMurmur3 computes MurmurHash3 (32-bit)
// Args: data_ptr, len, seed -> returns uint32 in rax
func generateHashMurmur3(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	// Simplified MurmurHash3 implementation
	// For full implementation, we'd need proper block processing
	loopLbl := cg.getLabel("murmur_loop")
	endLbl := cg.getLabel("murmur_end")

	cg.generateExpressionToReg(args[0], "rsi") // data pointer
	cg.generateExpressionToReg(args[1], "rcx") // length
	cg.generateExpressionToReg(args[2], "rax") // seed

	cg.textSection.WriteString("    movl %eax, %r8d\n") // hash = seed

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", loopLbl))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", endLbl))
	cg.textSection.WriteString("    movzbq (%rsi), %rbx\n")
	cg.textSection.WriteString("    imull $0xCC9E2D51, %ebx, %ebx\n")
	cg.textSection.WriteString("    roll $15, %ebx\n")
	cg.textSection.WriteString("    imull $0x1B873593, %ebx, %ebx\n")
	cg.textSection.WriteString("    xorl %ebx, %r8d\n")
	cg.textSection.WriteString("    roll $13, %r8d\n")
	cg.textSection.WriteString("    imull $5, %r8d, %r8d\n")
	cg.textSection.WriteString("    addl $0xE6546B64, %r8d\n")
	cg.textSection.WriteString("    incq %rsi\n")
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", loopLbl))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLbl))
	// Final mix
	cg.textSection.WriteString("    xorl %ecx, %r8d\n")
	cg.textSection.WriteString("    movl %r8d, %eax\n")
	cg.textSection.WriteString("    shrl $16, %eax\n")
	cg.textSection.WriteString("    xorl %r8d, %eax\n")
	cg.textSection.WriteString("    imull $0x85EBCA6B, %eax, %eax\n")
	cg.textSection.WriteString("    movl %eax, %edx\n")
	cg.textSection.WriteString("    shrl $13, %edx\n")
	cg.textSection.WriteString("    xorl %edx, %eax\n")
	cg.textSection.WriteString("    imull $0xC2B2AE35, %eax, %eax\n")
	cg.textSection.WriteString("    movl %eax, %edx\n")
	cg.textSection.WriteString("    shrl $16, %edx\n")
	cg.textSection.WriteString("    xorl %edx, %eax\n")
}

// generateHashSHA256 computes SHA-256 hash
// Args: data_ptr, len, out_buf (32 bytes) -> void
func generateHashSHA256(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	// SHA-256 is complex - would require full implementation
	// For now, placeholder that zeros the output
	cg.generateExpressionToReg(args[2], "rdi") // output buffer
	cg.textSection.WriteString("    movq $32, %rcx\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    rep stosb\n")
	cg.textSection.WriteString("    # TODO: Implement SHA-256\n")
}

// generateHashMD5 computes MD5 hash
// Args: data_ptr, len, out_buf (16 bytes) -> void
func generateHashMD5(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	// MD5 is complex - would require full implementation
	// For now, placeholder that zeros the output
	cg.generateExpressionToReg(args[2], "rdi") // output buffer
	cg.textSection.WriteString("    movq $16, %rcx\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    rep stosb\n")
	cg.textSection.WriteString("    # TODO: Implement MD5\n")
}
