package main

// llvm_codegen.go - LLVM IR code generation for Lotus
//
// This file provides an alternative backend that generates LLVM IR instead of
// direct x86-64 assembly. Benefits include:
// - Cross-platform compilation (x86, ARM, RISC-V, WebAssembly)
// - Advanced optimizations via LLVM's optimization passes
// - Better register allocation
// - Debug information support (DWARF)
// - Smaller, more efficient binaries

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

// LLVMVariable holds an alloca and its element type for correct loading
type LLVMVariable struct {
	Alloca      llvm.Value
	ElementType llvm.Type
}

// LLVMCodeGenerator generates LLVM IR from Lotus AST
type LLVMCodeGenerator struct {
	context llvm.Context
	module  llvm.Module
	builder llvm.Builder

	// Current function being generated
	currentFn llvm.Value

	// Symbol tables
	namedValues   map[string]LLVMVariable // Local variables with type info
	globalStrings map[string]llvm.Value   // String constants
	functions     map[string]llvm.Value   // Declared functions
	globalVars    map[string]llvm.Value   // Global and static variables
	staticVars    map[string]llvm.Value   // Static variables (function-local persistent)

	// Virtual function tables
	vtables        map[string]llvm.Value // Class name -> vtable global
	virtualMethods map[string]bool       // Set of virtual method names

	// Type registries for LLVM
	structTypes map[string]llvm.Type           // Struct name -> LLVM struct type
	structDefs  map[string]*StructDefinition   // Struct name -> definition
	enumDefs    map[string]*EnumDefinition     // Enum name -> definition
	classDefs   map[string]*ClassDefinition    // Class name -> definition
	classTypes  map[string]llvm.Type           // Class name -> LLVM struct type

	// Partial applications (currying)
	partialApplications map[string]*PartialApplication // Function name -> partial info

	// Import context for stdlib
	imports *ImportContext

	// String counter for unique names
	stringCounter int
}

// NewLLVMCodeGenerator creates a new LLVM code generator
func NewLLVMCodeGenerator(moduleName string) *LLVMCodeGenerator {
	ctx := llvm.NewContext()
	mod := ctx.NewModule(moduleName)
	builder := ctx.NewBuilder()

	cg := &LLVMCodeGenerator{
		context:             ctx,
		module:              mod,
		builder:             builder,
		namedValues:         make(map[string]LLVMVariable),
		globalStrings:       make(map[string]llvm.Value),
		functions:           make(map[string]llvm.Value),
		globalVars:          make(map[string]llvm.Value),
		staticVars:          make(map[string]llvm.Value),
		vtables:             make(map[string]llvm.Value),
		virtualMethods:      make(map[string]bool),
		structTypes:         make(map[string]llvm.Type),
		structDefs:          make(map[string]*StructDefinition),
		enumDefs:            make(map[string]*EnumDefinition),
		classDefs:           make(map[string]*ClassDefinition),
		classTypes:          make(map[string]llvm.Type),
		imports:             NewImportContext(),
		partialApplications: make(map[string]*PartialApplication),
		stringCounter:       0,
	}

	// Declare external C library functions
	cg.declareExternalFunctions()

	return cg
}

// declareExternalFunctions declares C library functions we'll use
func (cg *LLVMCodeGenerator) declareExternalFunctions() {
	// printf(const char* format, ...) -> int
	printfType := llvm.FunctionType(
		cg.context.Int32Type(),
		[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0)},
		true, // variadic
	)
	printf := llvm.AddFunction(cg.module, "printf", printfType)
	printf.SetLinkage(llvm.ExternalLinkage)
	cg.functions["printf"] = printf

	// puts(const char* s) -> int
	putsType := llvm.FunctionType(
		cg.context.Int32Type(),
		[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0)},
		false,
	)
	puts := llvm.AddFunction(cg.module, "puts", putsType)
	puts.SetLinkage(llvm.ExternalLinkage)
	cg.functions["puts"] = puts

	// malloc(size_t size) -> void*
	mallocType := llvm.FunctionType(
		llvm.PointerType(cg.context.Int8Type(), 0),
		[]llvm.Type{cg.context.Int64Type()},
		false,
	)
	malloc := llvm.AddFunction(cg.module, "malloc", mallocType)
	malloc.SetLinkage(llvm.ExternalLinkage)
	cg.functions["malloc"] = malloc

	// free(void* ptr)
	freeType := llvm.FunctionType(
		cg.context.VoidType(),
		[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0)},
		false,
	)
	free := llvm.AddFunction(cg.module, "free", freeType)
	free.SetLinkage(llvm.ExternalLinkage)
	cg.functions["free"] = free

	// exit(int status)
	exitType := llvm.FunctionType(
		cg.context.VoidType(),
		[]llvm.Type{cg.context.Int32Type()},
		false,
	)
	exit := llvm.AddFunction(cg.module, "exit", exitType)
	exit.SetLinkage(llvm.ExternalLinkage)
	cg.functions["exit"] = exit
}

// Generate generates LLVM IR from the AST
func (cg *LLVMCodeGenerator) Generate(statements []ASTNode) error {
	// First pass: declare all functions and top-level global variables
	for _, stmt := range statements {
		switch s := stmt.(type) {
		case *FunctionDefinition:
			if err := cg.declareFunction(s); err != nil {
				return err
			}
		case *VariableDeclaration:
			// Top-level variables are implicitly global
			if err := cg.declareTopLevelVar(s); err != nil {
				return err
			}
		}
	}

	// Second pass: generate function bodies (skip top-level vars, already handled)
	for _, stmt := range statements {
		switch stmt.(type) {
		case *VariableDeclaration:
			// Already handled in first pass
			continue
		default:
			if err := cg.generateStatement(stmt); err != nil {
				return err
			}
		}
	}

	// Verify the module
	if err := llvm.VerifyModule(cg.module, llvm.ReturnStatusAction); err != nil {
		return fmt.Errorf("LLVM module verification failed: %v", err)
	}

	return nil
}

// declareTopLevelVar declares a top-level (global) variable
func (cg *LLVMCodeGenerator) declareTopLevelVar(v *VariableDeclaration) error {
	varType := cg.tokenTypeToLLVM(v.Type)

	// Check if already created
	if _, ok := cg.globalVars[v.Name]; ok {
		return nil // Already declared
	}

	// Create global variable
	global := llvm.AddGlobal(cg.module, varType, v.Name)

	// Set linkage based on storage class
	if v.Storage == StorageStatic {
		global.SetLinkage(llvm.InternalLinkage)
	} else {
		global.SetLinkage(llvm.ExternalLinkage)
	}

	// Initialize with constant value if provided
	if v.Value != nil {
		if constVal := cg.tryGetConstantValue(v.Value, varType); !constVal.IsNil() {
			global.SetInitializer(constVal)
		} else {
			// Non-constant initializer - use zero init and set in main
			global.SetInitializer(llvm.ConstNull(varType))
		}
	} else {
		global.SetInitializer(llvm.ConstNull(varType))
	}

	cg.globalVars[v.Name] = global
	cg.namedValues[v.Name] = LLVMVariable{Alloca: global, ElementType: varType}
	return nil
}

// tryGetConstantValue attempts to evaluate an expression as a constant
func (cg *LLVMCodeGenerator) tryGetConstantValue(expr ASTNode, targetType llvm.Type) llvm.Value {
	switch e := expr.(type) {
	case *IntLiteral:
		return llvm.ConstInt(targetType, uint64(e.Value), true)
	case *BoolLiteral:
		if e.Value {
			return llvm.ConstInt(targetType, 1, false)
		}
		return llvm.ConstInt(targetType, 0, false)
	case *FloatLiteral:
		return llvm.ConstFloat(targetType, float64(e.Value)/1000.0)
	case *StringLiteral:
		// Create a global string constant using the builder-free method
		return cg.createGlobalStringConstant(e.Value)
	default:
		return llvm.Value{} // Not a constant
	}
}

// declareFunction creates a function declaration
func (cg *LLVMCodeGenerator) declareFunction(fn *FunctionDefinition) error {
	// Build parameter types
	paramTypes := make([]llvm.Type, len(fn.Parameters))
	for i, param := range fn.Parameters {
		paramTypes[i] = cg.tokenTypeToLLVM(param.Type)
	}

	// Build function type
	retType := cg.tokenTypeToLLVM(fn.ReturnType)
	fnType := llvm.FunctionType(retType, paramTypes, false)

	// Create function
	llvmFn := llvm.AddFunction(cg.module, fn.Name, fnType)

	// Set linkage based on static modifier
	if fn.IsStatic {
		llvmFn.SetLinkage(llvm.InternalLinkage) // Static functions are file-local
	} else {
		llvmFn.SetLinkage(llvm.ExternalLinkage)
	}

	// Name parameters
	for i, param := range fn.Parameters {
		llvmFn.Param(i).SetName(param.Name)
	}

	cg.functions[fn.Name] = llvmFn

	// Track virtual methods for vtable generation
	if fn.IsVirtual {
		cg.virtualMethods[fn.Name] = true
	}

	return nil
}

// generateStatement generates LLVM IR for a statement
func (cg *LLVMCodeGenerator) generateStatement(stmt ASTNode) error {
	switch s := stmt.(type) {
	case *FunctionDefinition:
		return cg.generateFunction(s)
	case *VariableDeclaration:
		return cg.generateVarDecl(s)
	case *ConstantDeclaration:
		return cg.generateConstDecl(s)
	case *ReturnStatement:
		return cg.generateReturn(s)
	case *FunctionCall:
		_, err := cg.generateCall(s)
		return err
	case *IfStatement:
		return cg.generateIf(s)
	case *WhileLoop:
		return cg.generateWhile(s)
	case *ForLoop:
		return cg.generateFor(s)
	case *Assignment:
		return cg.generateAssignment(s)
	case *ImportStatement:
		return cg.handleImport(s)
	case *WrapperDefinition:
		return cg.generateWrapper(s)
	case *DecoratedFunction:
		return cg.generateDecoratedFunction(s)
	case *PartialApplication:
		_, err := cg.generatePartialApplication(s)
		return err
	case *CompoundAssignment:
		return cg.generateCompoundAssignment(s)
	case *IncrementOp:
		return cg.generateIncrementOp(s)
	case *StructDefinition:
		return cg.generateStructDef(s)
	case *EnumDefinition:
		return cg.generateEnumDef(s)
	case *ClassDefinition:
		return cg.generateClassDef(s)
	case *ArrayDeclaration:
		return cg.generateArrayDecl(s)
	case *FreeCall:
		return cg.generateFreeStmt(s)
	default:
		return fmt.Errorf("unsupported statement type: %T", stmt)
	}
}

// generateFunction generates a complete function
func (cg *LLVMCodeGenerator) generateFunction(fn *FunctionDefinition) error {
	llvmFn, ok := cg.functions[fn.Name]
	if !ok {
		return fmt.Errorf("function %s not declared", fn.Name)
	}

	cg.currentFn = llvmFn
	cg.namedValues = make(map[string]LLVMVariable) // Fresh scope

	// Create entry block
	entry := llvm.AddBasicBlock(llvmFn, "entry")
	cg.builder.SetInsertPointAtEnd(entry)

	// Allocate space for parameters and store them
	for i, param := range fn.Parameters {
		paramType := cg.tokenTypeToLLVM(param.Type)
		alloca := cg.builder.CreateAlloca(paramType, param.Name)
		cg.builder.CreateStore(llvmFn.Param(i), alloca)
		cg.namedValues[param.Name] = LLVMVariable{Alloca: alloca, ElementType: paramType}
	}

	// Generate function body
	for _, stmt := range fn.Body {
		if err := cg.generateStatement(stmt); err != nil {
			return err
		}
	}

	// Add implicit return if needed
	lastBlock := cg.builder.GetInsertBlock()
	if lastBlock.LastInstruction().IsNil() || lastBlock.LastInstruction().InstructionOpcode() != llvm.Ret {
		// Use the function definition's return type instead of querying LLVM
		if fn.ReturnType == TokenTypeVoid || fn.ReturnType == TokenRet || fn.ReturnType == TokenReturn {
			cg.builder.CreateRetVoid()
		} else {
			retType := cg.tokenTypeToLLVM(fn.ReturnType)
			cg.builder.CreateRet(llvm.ConstNull(retType))
		}
	}

	return nil
}

// generateVarDecl generates a variable declaration
func (cg *LLVMCodeGenerator) generateVarDecl(v *VariableDeclaration) error {
	varType := cg.tokenTypeToLLVM(v.Type)

	switch v.Storage {
	case StorageStatic:
		// Static variables are allocated in the data section (global with internal linkage)
		return cg.generateStaticVar(v, varType)
	case StorageGlobal:
		// Global variables are allocated in the data section with external linkage
		return cg.generateGlobalVar(v, varType)
	case StorageLocal, StorageAuto:
		// Local/auto variables are stack allocated (default behavior)
		return cg.generateLocalVar(v, varType)
	default:
		return cg.generateLocalVar(v, varType)
	}
}

// generateConstDecl generates a constant declaration
// Constants are implemented as global variables with constant initializers
func (cg *LLVMCodeGenerator) generateConstDecl(c *ConstantDeclaration) error {
	constType := cg.tokenTypeToLLVM(c.Type)

	// Create a global constant
	global := llvm.AddGlobal(cg.module, constType, c.Name)
	global.SetLinkage(llvm.InternalLinkage)
	global.SetGlobalConstant(true)

	// Get the constant value
	if c.Value != nil {
		constVal := cg.tryGetConstantValue(c.Value, constType)
		if constVal.IsNil() {
			return fmt.Errorf("constant %s must have a compile-time constant value", c.Name)
		}
		global.SetInitializer(constVal)
	} else {
		return fmt.Errorf("constant %s must have a value", c.Name)
	}

	// Store in global vars map so it can be accessed
	cg.globalVars[c.Name] = global
	return nil
}

// generateLocalVar generates a stack-allocated local variable
func (cg *LLVMCodeGenerator) generateLocalVar(v *VariableDeclaration, varType llvm.Type) error {
	alloca := cg.builder.CreateAlloca(varType, v.Name)

	if v.Value != nil {
		val, err := cg.generateExpression(v.Value)
		if err != nil {
			return err
		}
		cg.builder.CreateStore(val, alloca)
	}

	cg.namedValues[v.Name] = LLVMVariable{Alloca: alloca, ElementType: varType}
	return nil
}

// generateStaticVar generates a static variable (persistent storage, internal linkage)
func (cg *LLVMCodeGenerator) generateStaticVar(v *VariableDeclaration, varType llvm.Type) error {
	// Create a unique name for the static variable (function-scoped)
	staticName := fmt.Sprintf("_static_%s_%s", cg.currentFn.Name(), v.Name)

	// Check if already created (for re-entry into same function)
	if existing, ok := cg.staticVars[staticName]; ok {
		cg.namedValues[v.Name] = LLVMVariable{Alloca: existing, ElementType: varType}
		return nil
	}

	// Create global variable with internal linkage
	global := llvm.AddGlobal(cg.module, varType, staticName)
	global.SetLinkage(llvm.InternalLinkage)

	// For static variables, use constant initializer if available
	if v.Value != nil {
		if constVal := cg.tryGetConstantValue(v.Value, varType); !constVal.IsNil() {
			global.SetInitializer(constVal)
		} else {
			// Non-constant initializer - use guard variable pattern
			global.SetInitializer(llvm.ConstNull(varType))
			cg.generateStaticInitGuard(staticName, global, v.Value)
		}
	} else {
		global.SetInitializer(llvm.ConstNull(varType))
	}

	// Store in static vars map
	cg.staticVars[staticName] = global
	cg.namedValues[v.Name] = LLVMVariable{Alloca: global, ElementType: varType}
	return nil
}

// generateStaticInitGuard generates a guard for one-time static initialization
func (cg *LLVMCodeGenerator) generateStaticInitGuard(staticName string, global llvm.Value, initValue ASTNode) {
	guardName := staticName + "_guard"

	// Check if guard exists
	if _, ok := cg.staticVars[guardName]; ok {
		return // Already initialized
	}

	// Create guard variable
	guardVar := llvm.AddGlobal(cg.module, cg.context.Int1Type(), guardName)
	guardVar.SetLinkage(llvm.InternalLinkage)
	guardVar.SetInitializer(llvm.ConstInt(cg.context.Int1Type(), 0, false))
	cg.staticVars[guardName] = guardVar

	// Generate: if (!guard) { static_var = init_val; guard = true; }
	initBlock := llvm.AddBasicBlock(cg.currentFn, "static_init")
	contBlock := llvm.AddBasicBlock(cg.currentFn, "static_cont")

	// Load guard and branch
	guardVal := cg.builder.CreateLoad(cg.context.Int1Type(), guardVar, "guard")
	cg.builder.CreateCondBr(guardVal, contBlock, initBlock)

	// Init block
	cg.builder.SetInsertPointAtEnd(initBlock)
	val, _ := cg.generateExpression(initValue)
	cg.builder.CreateStore(val, global)
	cg.builder.CreateStore(llvm.ConstInt(cg.context.Int1Type(), 1, false), guardVar)
	cg.builder.CreateBr(contBlock)

	// Continue block
	cg.builder.SetInsertPointAtEnd(contBlock)
}

// generateGlobalVar generates a global variable (external linkage)
func (cg *LLVMCodeGenerator) generateGlobalVar(v *VariableDeclaration, varType llvm.Type) error {
	// Check if already created
	if existing, ok := cg.globalVars[v.Name]; ok {
		cg.namedValues[v.Name] = LLVMVariable{Alloca: existing, ElementType: varType}
		return nil
	}

	// Create global variable with external linkage
	global := llvm.AddGlobal(cg.module, varType, v.Name)
	global.SetLinkage(llvm.ExternalLinkage)

	// Initialize with constant value if available, otherwise null
	if v.Value != nil {
		// Try to get a constant initializer
		initVal := cg.getConstantValue(v.Value, varType)
		global.SetInitializer(initVal)
	} else {
		global.SetInitializer(llvm.ConstNull(varType))
	}

	// Store in global vars map
	cg.globalVars[v.Name] = global
	cg.namedValues[v.Name] = LLVMVariable{Alloca: global, ElementType: varType}
	return nil
}

// getConstantValue tries to get a constant LLVM value from an AST node
func (cg *LLVMCodeGenerator) getConstantValue(node ASTNode, varType llvm.Type) llvm.Value {
	switch n := node.(type) {
	case *IntLiteral:
		return llvm.ConstInt(varType, uint64(n.Value), true)
	case *BoolLiteral:
		if n.Value {
			return llvm.ConstInt(varType, 1, false)
		}
		return llvm.ConstInt(varType, 0, false)
	default:
		// Default to null for non-constant expressions
		return llvm.ConstNull(varType)
	}
}

// generateReturn generates a return statement
func (cg *LLVMCodeGenerator) generateReturn(r *ReturnStatement) error {
	if r.Value == nil {
		cg.builder.CreateRetVoid()
		return nil
	}

	val, err := cg.generateExpression(r.Value)
	if err != nil {
		return err
	}
	cg.builder.CreateRet(val)
	return nil
}

// generateExpression generates LLVM IR for an expression
func (cg *LLVMCodeGenerator) generateExpression(expr ASTNode) (llvm.Value, error) {
	switch e := expr.(type) {
	case *IntLiteral:
		return llvm.ConstInt(cg.context.Int64Type(), uint64(e.Value), true), nil

	case *StringLiteral:
		return cg.createGlobalString(e.Value), nil

	case *BoolLiteral:
		val := uint64(0)
		if e.Value {
			val = 1
		}
		return llvm.ConstInt(cg.context.Int1Type(), val, false), nil

	case *FloatLiteral:
		// FloatLiteral stores value as int64 * 1000 for precision
		return llvm.ConstFloat(cg.context.DoubleType(), float64(e.Value)/1000.0), nil

	case *CharLiteral:
		// CharLiteral stores a single Unicode character as string
		charVal := uint64(0)
		if len(e.Value) > 0 {
			runes := []rune(e.Value)
			charVal = uint64(runes[0])
		}
		return llvm.ConstInt(cg.context.Int32Type(), charVal, false), nil

	case *NullLiteral:
		return llvm.ConstNull(llvm.PointerType(cg.context.Int8Type(), 0)), nil

	case *TernaryOp:
		return cg.generateTernary(e)

	case *Identifier:
		// First check local variables
		variable, ok := cg.namedValues[e.Name]
		if !ok {
			// Then check global variables
			if globalVar, gok := cg.globalVars[e.Name]; gok {
				varType := globalVar.GlobalValueType()
				return cg.builder.CreateLoad(varType, globalVar, e.Name), nil
			}
			return llvm.Value{}, fmt.Errorf("undefined variable: %s", e.Name)
		}
		return cg.builder.CreateLoad(variable.ElementType, variable.Alloca, e.Name), nil

	case *BinaryOp:
		return cg.generateBinaryExpr(e)

	case *BitwiseOp:
		return cg.generateBitwiseOp(e)

	case *UnaryOp:
		return cg.generateUnaryExpr(e)

	case *Comparison:
		return cg.generateComparison(e)

	case *LogicalOp:
		return cg.generateLogicalOp(e)

	case *FunctionCall:
		return cg.generateCall(e)

	case *PartialApplication:
		return cg.generatePartialApplication(e)

	case *PipeExpression:
		return cg.generatePipeExpression(e)

	case *pipeValueWrapper:
		return e.value, nil

	case *ArrayAccess:
		return cg.generateArrayAccess(e)

	case *ArrayLiteral:
		return cg.generateArrayLiteral(e)

	case *FieldAccess:
		return cg.generateFieldAccess(e)

	case *StructLiteral:
		return cg.generateStructLiteral(e)

	case *EnumLiteral:
		return cg.generateEnumLiteral(e)

	case *ClassLiteral:
		return cg.generateClassLiteral(e)

	case *MethodCall:
		return cg.generateMethodCall(e)

	case *Reference:
		return cg.generateReference(e)

	case *Dereference:
		return cg.generateDereference(e)

	case *SizeofExpr:
		return cg.generateSizeof(e)

	case *MallocCall:
		return cg.generateMallocExpr(e)

	default:
		return llvm.Value{}, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

// generateBinaryExpr generates a binary expression
func (cg *LLVMCodeGenerator) generateBinaryExpr(b *BinaryOp) (llvm.Value, error) {
	left, err := cg.generateExpression(b.Left)
	if err != nil {
		return llvm.Value{}, err
	}
	right, err := cg.generateExpression(b.Right)
	if err != nil {
		return llvm.Value{}, err
	}

	switch b.Operator {
	// Arithmetic
	case TokenPlus:
		return cg.builder.CreateAdd(left, right, "addtmp"), nil
	case TokenMinus:
		return cg.builder.CreateSub(left, right, "subtmp"), nil
	case TokenStar:
		return cg.builder.CreateMul(left, right, "multmp"), nil
	case TokenSlash:
		return cg.builder.CreateSDiv(left, right, "divtmp"), nil
	case TokenPercent:
		return cg.builder.CreateSRem(left, right, "modtmp"), nil

	// Bitwise
	case TokenAmpersand:
		return cg.builder.CreateAnd(left, right, "andtmp"), nil
	case TokenPipe:
		return cg.builder.CreateOr(left, right, "ortmp"), nil
	case TokenCaret:
		return cg.builder.CreateXor(left, right, "xortmp"), nil
	case TokenLShift:
		return cg.builder.CreateShl(left, right, "shltmp"), nil
	case TokenRShift:
		return cg.builder.CreateAShr(left, right, "shrtmp"), nil

	default:
		return llvm.Value{}, fmt.Errorf("unsupported binary operator: %v", b.Operator)
	}
}

// generateBitwiseOp generates a bitwise operation
func (cg *LLVMCodeGenerator) generateBitwiseOp(b *BitwiseOp) (llvm.Value, error) {
	left, err := cg.generateExpression(b.Left)
	if err != nil {
		return llvm.Value{}, err
	}
	right, err := cg.generateExpression(b.Right)
	if err != nil {
		return llvm.Value{}, err
	}

	switch b.Operator {
	case TokenAmpersand:
		return cg.builder.CreateAnd(left, right, "andtmp"), nil
	case TokenPipe:
		return cg.builder.CreateOr(left, right, "ortmp"), nil
	case TokenCaret:
		return cg.builder.CreateXor(left, right, "xortmp"), nil
	case TokenLShift:
		return cg.builder.CreateShl(left, right, "shltmp"), nil
	case TokenRShift:
		return cg.builder.CreateAShr(left, right, "shrtmp"), nil
	default:
		return llvm.Value{}, fmt.Errorf("unsupported bitwise operator: %v", b.Operator)
	}
}

// generateComparison generates a comparison expression
func (cg *LLVMCodeGenerator) generateComparison(c *Comparison) (llvm.Value, error) {
	left, err := cg.generateExpression(c.Left)
	if err != nil {
		return llvm.Value{}, err
	}
	right, err := cg.generateExpression(c.Right)
	if err != nil {
		return llvm.Value{}, err
	}

	var cmp llvm.Value
	switch c.Operator {
	case TokenEqual:
		cmp = cg.builder.CreateICmp(llvm.IntEQ, left, right, "eqtmp")
	case TokenNotEqual:
		cmp = cg.builder.CreateICmp(llvm.IntNE, left, right, "netmp")
	case TokenLess:
		cmp = cg.builder.CreateICmp(llvm.IntSLT, left, right, "lttmp")
	case TokenLessEq:
		cmp = cg.builder.CreateICmp(llvm.IntSLE, left, right, "letmp")
	case TokenGreater:
		cmp = cg.builder.CreateICmp(llvm.IntSGT, left, right, "gttmp")
	case TokenGreaterEq:
		cmp = cg.builder.CreateICmp(llvm.IntSGE, left, right, "getmp")
	default:
		return llvm.Value{}, fmt.Errorf("unsupported comparison operator: %v", c.Operator)
	}
	return cg.builder.CreateZExt(cmp, cg.context.Int64Type(), "cmpext"), nil
}

// generateLogicalOp generates a logical operation
func (cg *LLVMCodeGenerator) generateLogicalOp(l *LogicalOp) (llvm.Value, error) {
	left, err := cg.generateExpression(l.Left)
	if err != nil {
		return llvm.Value{}, err
	}
	right, err := cg.generateExpression(l.Right)
	if err != nil {
		return llvm.Value{}, err
	}

	leftBool := cg.builder.CreateICmp(llvm.IntNE, left, llvm.ConstInt(cg.context.Int64Type(), 0, false), "leftbool")
	rightBool := cg.builder.CreateICmp(llvm.IntNE, right, llvm.ConstInt(cg.context.Int64Type(), 0, false), "rightbool")

	var result llvm.Value
	switch l.Operator {
	case TokenAnd:
		result = cg.builder.CreateAnd(leftBool, rightBool, "andtmp")
	case TokenOr:
		result = cg.builder.CreateOr(leftBool, rightBool, "ortmp")
	default:
		return llvm.Value{}, fmt.Errorf("unsupported logical operator: %v", l.Operator)
	}
	return cg.builder.CreateZExt(result, cg.context.Int64Type(), "logext"), nil
}

// generateUnaryExpr generates a unary expression
func (cg *LLVMCodeGenerator) generateUnaryExpr(u *UnaryOp) (llvm.Value, error) {
	operand, err := cg.generateExpression(u.Operand)
	if err != nil {
		return llvm.Value{}, err
	}

	switch u.Operator {
	case TokenMinus:
		return cg.builder.CreateNeg(operand, "negtmp"), nil
	case TokenTilde:
		return cg.builder.CreateNot(operand, "nottmp"), nil
	case TokenExclaim:
		cmp := cg.builder.CreateICmp(llvm.IntEQ, operand, llvm.ConstInt(cg.context.Int64Type(), 0, false), "lnottmp")
		return cg.builder.CreateZExt(cmp, cg.context.Int64Type(), "lnotext"), nil
	default:
		return llvm.Value{}, fmt.Errorf("unsupported unary operator: %v", u.Operator)
	}
}

// generateCall generates a function call
func (cg *LLVMCodeGenerator) generateCall(call *FunctionCall) (llvm.Value, error) {
	// Handle builtin functions
	switch call.Name {
	case "println", "print":
		return cg.generatePrint(call)
	case "printf":
		return cg.generatePrintf(call)

	// Math stdlib functions
	case "abs":
		return cg.generateAbs(call)
	case "min":
		return cg.generateMin(call)
	case "max":
		return cg.generateMax(call)
	case "sqrt":
		return cg.generateSqrt(call)
	case "pow":
		return cg.generatePow(call)

	// Number conversion functions
	case "toUint32":
		return cg.generateToUint32(call)
	case "toBool":
		return cg.generateToBool(call)
	case "toInt":
		return cg.generateToInt(call)
	case "toFloat":
		return cg.generateToFloat(call)

	// String stdlib functions
	case "len":
		return cg.generateStrlen(call)
	case "concat":
		return cg.generateConcat(call)
	case "contains":
		return cg.generateContains(call)
	case "copy":
		return cg.generateStrCopy(call)
	case "compare":
		return cg.generateStrCompare(call)
	case "indexOf":
		return cg.generateIndexOf(call)
	case "substring":
		return cg.generateSubstring(call)
	case "toUpper":
		return cg.generateToUpper(call)
	case "toLower":
		return cg.generateToLower(call)
	case "trim":
		return cg.generateTrim(call)
	case "split":
		return cg.generateSplit(call)
	case "replace":
		return cg.generateReplace(call)
	case "startsWith":
		return cg.generateStartsWith(call)
	case "endsWith":
		return cg.generateEndsWith(call)

	// Memory stdlib functions
	case "malloc":
		return cg.generateMalloc(call)
	case "free":
		return cg.generateFree(call)

	// Collections stdlib functions
	case "array_int_new":
		return cg.generateArrayIntNew(call)
	case "array_int_push":
		return cg.generateArrayIntPush(call)
	case "array_int_pop":
		return cg.generateArrayIntPop(call)
	case "array_int_len":
		return cg.generateArrayIntLen(call)
	case "array_int_resize":
		return cg.generateArrayIntResize(call)
	case "array_int_reserve":
		return cg.generateArrayIntReserve(call)
	case "array_int_shrink":
		return cg.generateArrayIntShrink(call)
	case "array_int_capacity":
		return cg.generateArrayIntCapacity(call)
	case "array_int_get":
		return cg.generateArrayIntGet(call)
	case "array_int_set":
		return cg.generateArrayIntSet(call)
	case "array_int_free":
		return cg.generateArrayIntFree(call)

	// Stack functions
	case "stack_int_new":
		return cg.generateStackIntNew(call)
	case "stack_int_push":
		return cg.generateStackIntPush(call)
	case "stack_int_pop":
		return cg.generateStackIntPop(call)
	case "stack_int_len":
		return cg.generateStackIntLen(call)

	// Queue functions
	case "queue_int_new":
		return cg.generateQueueIntNew(call)
	case "queue_int_enqueue":
		return cg.generateQueueIntEnqueue(call)
	case "queue_int_dequeue":
		return cg.generateQueueIntDequeue(call)
	case "queue_int_len":
		return cg.generateQueueIntLen(call)

	// Deque functions
	case "deque_int_new":
		return cg.generateDequeIntNew(call)
	case "deque_int_push_front":
		return cg.generateDequeIntPushFront(call)
	case "deque_int_push_back":
		return cg.generateDequeIntPushBack(call)
	case "deque_int_pop_front":
		return cg.generateDequeIntPopFront(call)
	case "deque_int_pop_back":
		return cg.generateDequeIntPopBack(call)
	case "deque_int_len":
		return cg.generateDequeIntLen(call)

	// Heap functions
	case "heap_int_new":
		return cg.generateHeapIntNew(call)
	case "heap_int_push":
		return cg.generateHeapIntPush(call)
	case "heap_int_pop":
		return cg.generateHeapIntPop(call)
	case "heap_int_peek":
		return cg.generateHeapIntPeek(call)
	case "heap_int_len":
		return cg.generateHeapIntLen(call)

	// Net/socket functions
	case "socket":
		return cg.generateSocket(call)
	case "close":
		return cg.generateClose(call)

	// File I/O functions
	case "open":
		return cg.generateOpen(call)
	case "read":
		return cg.generateRead(call)
	case "write":
		return cg.generateWrite(call)

	// Time functions
	case "now":
		return cg.generateNow(call)
	case "millis":
		return cg.generateMillis(call)
	case "nanos":
		return cg.generateNanos(call)
	case "sleep":
		return cg.generateSleep(call)

	// HTTP pool functions
	case "pool_new":
		return cg.generatePoolNew(call)
	case "pool_get":
		return cg.generatePoolGet(call)
	case "pool_put":
		return cg.generatePoolPut(call)
	case "pool_close":
		return cg.generatePoolClose(call)

	// JSON module functions
	case "json_parse":
		return cg.generateJSONParse(call)
	case "json_stringify":
		return cg.generateJSONStringify(call)
	case "json_get":
		return cg.generateJSONGet(call)
	case "json_type":
		return cg.generateJSONGetType(call)
	case "json_int":
		return cg.generateJSONGetInt(call)
	case "json_str":
		return cg.generateJSONGetString(call)
	case "json_bool":
		return cg.generateJSONGetBool(call)
	case "json_array_len":
		return cg.generateJSONArrayLen(call)
	case "json_array_get":
		return cg.generateJSONArrayGet(call)
	case "json_new":
		return cg.generateJSONNew(call)
	case "json_free":
		return cg.generateJSONFree(call)

	// Formatting functions
	case "sprintf":
		return cg.generateSprintf(call)
	case "snprintf":
		return cg.generateSnprintf(call)
	case "format_int":
		return cg.generateFormatInt(call)
	case "format_hex":
		return cg.generateFormatHex(call)
	case "format_bin":
		return cg.generateFormatBinary(call)
	case "pad_left":
		return cg.generatePadLeft(call)
	case "pad_right":
		return cg.generatePadRight(call)

	// Random functions
	case "rand":
		return cg.generateRand(call)
	case "rand_range":
		return cg.generateRandRange(call)
	case "seed":
		return cg.generateSeed(call)
	case "rand_float":
		return cg.generateRandFloat(call)
	case "rand_bool":
		return cg.generateRandBool(call)
	case "rand_bytes":
		return cg.generateRandBytes(call)
	case "shuffle":
		return cg.generateShuffle(call)
	case "choice":
		return cg.generateChoice(call)
	case "rand_n":
		return cg.generateRandN(call)
	case "rand_string":
		return cg.generateRandString(call)
	}

	// Look up the function
	fn, ok := cg.functions[call.Name]
	if !ok {
		return llvm.Value{}, fmt.Errorf("undefined function: %s", call.Name)
	}

	// Generate arguments
	args := make([]llvm.Value, len(call.Args))
	for i, arg := range call.Args {
		val, err := cg.generateExpression(arg)
		if err != nil {
			return llvm.Value{}, err
		}
		args[i] = val
	}

	// Check if function returns void - void calls cannot have names
	fnType := fn.GlobalValueType()
	retType := fnType.ReturnType()
	if retType.TypeKind() == llvm.VoidTypeKind {
		return cg.builder.CreateCall(fnType, fn, args, ""), nil
	}
	return cg.builder.CreateCall(fnType, fn, args, "calltmp"), nil
}

// generatePrint generates a print/println call
func (cg *LLVMCodeGenerator) generatePrint(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) == 0 {
		return llvm.Value{}, nil
	}

	arg, err := cg.generateExpression(call.Args[0])
	if err != nil {
		return llvm.Value{}, err
	}

	puts := cg.functions["puts"]
	return cg.builder.CreateCall(puts.GlobalValueType(), puts, []llvm.Value{arg}, "putstmp"), nil
}

// generatePrintf generates a printf call
func (cg *LLVMCodeGenerator) generatePrintf(call *FunctionCall) (llvm.Value, error) {
	if len(call.Args) == 0 {
		return llvm.Value{}, fmt.Errorf("printf requires at least one argument")
	}

	args := make([]llvm.Value, len(call.Args))
	for i, arg := range call.Args {
		val, err := cg.generateExpression(arg)
		if err != nil {
			return llvm.Value{}, err
		}
		args[i] = val
	}

	printf := cg.functions["printf"]
	return cg.builder.CreateCall(printf.GlobalValueType(), printf, args, "printftmp"), nil
}

// generateIf generates an if statement
func (cg *LLVMCodeGenerator) generateIf(i *IfStatement) error {
	cond, err := cg.generateExpression(i.Condition)
	if err != nil {
		return err
	}

	// Convert to boolean
	condBool := cg.builder.CreateICmp(llvm.IntNE, cond, llvm.ConstInt(cg.context.Int64Type(), 0, false), "ifcond")

	// Create blocks
	thenBlock := llvm.AddBasicBlock(cg.currentFn, "then")
	elseBlock := llvm.AddBasicBlock(cg.currentFn, "else")
	mergeBlock := llvm.AddBasicBlock(cg.currentFn, "ifcont")

	cg.builder.CreateCondBr(condBool, thenBlock, elseBlock)

	// Generate then block
	cg.builder.SetInsertPointAtEnd(thenBlock)
	for _, stmt := range i.ThenBody {
		if err := cg.generateStatement(stmt); err != nil {
			return err
		}
	}
	if cg.builder.GetInsertBlock().LastInstruction().IsNil() ||
		cg.builder.GetInsertBlock().LastInstruction().InstructionOpcode() != llvm.Ret {
		cg.builder.CreateBr(mergeBlock)
	}

	// Generate else block
	cg.builder.SetInsertPointAtEnd(elseBlock)
	if len(i.ElseBody) > 0 {
		for _, stmt := range i.ElseBody {
			if err := cg.generateStatement(stmt); err != nil {
				return err
			}
		}
	}
	if cg.builder.GetInsertBlock().LastInstruction().IsNil() ||
		cg.builder.GetInsertBlock().LastInstruction().InstructionOpcode() != llvm.Ret {
		cg.builder.CreateBr(mergeBlock)
	}

	// Continue at merge block
	cg.builder.SetInsertPointAtEnd(mergeBlock)
	return nil
}

// generateWhile generates a while loop
func (cg *LLVMCodeGenerator) generateWhile(w *WhileLoop) error {
	condBlock := llvm.AddBasicBlock(cg.currentFn, "whilecond")
	bodyBlock := llvm.AddBasicBlock(cg.currentFn, "whilebody")
	exitBlock := llvm.AddBasicBlock(cg.currentFn, "whileexit")

	// Jump to condition
	cg.builder.CreateBr(condBlock)

	// Condition block
	cg.builder.SetInsertPointAtEnd(condBlock)
	cond, err := cg.generateExpression(w.Condition)
	if err != nil {
		return err
	}
	condBool := cg.builder.CreateICmp(llvm.IntNE, cond, llvm.ConstInt(cg.context.Int64Type(), 0, false), "whilecond")
	cg.builder.CreateCondBr(condBool, bodyBlock, exitBlock)

	// Body block
	cg.builder.SetInsertPointAtEnd(bodyBlock)
	for _, stmt := range w.Body {
		if err := cg.generateStatement(stmt); err != nil {
			return err
		}
	}
	cg.builder.CreateBr(condBlock)

	// Exit block
	cg.builder.SetInsertPointAtEnd(exitBlock)
	return nil
}

// generateFor generates a for loop
func (cg *LLVMCodeGenerator) generateFor(f *ForLoop) error {
	// Generate init
	if f.Init != nil {
		if err := cg.generateStatement(f.Init); err != nil {
			return err
		}
	}

	condBlock := llvm.AddBasicBlock(cg.currentFn, "forcond")
	bodyBlock := llvm.AddBasicBlock(cg.currentFn, "forbody")
	postBlock := llvm.AddBasicBlock(cg.currentFn, "forpost")
	exitBlock := llvm.AddBasicBlock(cg.currentFn, "forexit")

	cg.builder.CreateBr(condBlock)

	// Condition
	cg.builder.SetInsertPointAtEnd(condBlock)
	if f.Condition != nil {
		cond, err := cg.generateExpression(f.Condition)
		if err != nil {
			return err
		}
		condBool := cg.builder.CreateICmp(llvm.IntNE, cond, llvm.ConstInt(cg.context.Int64Type(), 0, false), "forcond")
		cg.builder.CreateCondBr(condBool, bodyBlock, exitBlock)
	} else {
		cg.builder.CreateBr(bodyBlock)
	}

	// Body
	cg.builder.SetInsertPointAtEnd(bodyBlock)
	for _, stmt := range f.Body {
		if err := cg.generateStatement(stmt); err != nil {
			return err
		}
	}
	cg.builder.CreateBr(postBlock)

	// Post (Update)
	cg.builder.SetInsertPointAtEnd(postBlock)
	if f.Update != nil {
		if err := cg.generateStatement(f.Update); err != nil {
			return err
		}
	}
	cg.builder.CreateBr(condBlock)

	// Exit
	cg.builder.SetInsertPointAtEnd(exitBlock)
	return nil
}

// generateAssignment generates an assignment
func (cg *LLVMCodeGenerator) generateAssignment(a *Assignment) error {
	val, err := cg.generateExpression(a.Value)
	if err != nil {
		return err
	}

	// Extract variable name from target (should be an Identifier)
	id, ok := a.Target.(*Identifier)
	if !ok {
		return fmt.Errorf("assignment target must be an identifier, got %T", a.Target)
	}

	// First check local variables
	variable, ok := cg.namedValues[id.Name]
	if ok {
		cg.builder.CreateStore(val, variable.Alloca)
		return nil
	}

	// Then check global variables
	if globalVar, gok := cg.globalVars[id.Name]; gok {
		cg.builder.CreateStore(val, globalVar)
		return nil
	}

	return fmt.Errorf("undefined variable: %s", id.Name)
}

// handleImport handles import statements
func (cg *LLVMCodeGenerator) handleImport(i *ImportStatement) error {
	return cg.imports.ProcessImport(i)
}

// createGlobalString creates a global string constant
func (cg *LLVMCodeGenerator) createGlobalString(s string) llvm.Value {
	if val, ok := cg.globalStrings[s]; ok {
		return val
	}

	name := fmt.Sprintf(".str.%d", cg.stringCounter)
	cg.stringCounter++

	globalStr := cg.builder.CreateGlobalStringPtr(s, name)
	cg.globalStrings[s] = globalStr
	return globalStr
}

// createGlobalStringConstant creates a global string constant without using the builder
// This is safe to call at module level (outside any function)
func (cg *LLVMCodeGenerator) createGlobalStringConstant(s string) llvm.Value {
	if val, ok := cg.globalStrings[s]; ok {
		return val
	}

	name := fmt.Sprintf(".str.%d", cg.stringCounter)
	cg.stringCounter++

	// Create the string data as a constant
	strData := llvm.ConstString(s, true) // true = null terminate

	// Create global for the string data
	strGlobal := llvm.AddGlobal(cg.module, strData.Type(), name)
	strGlobal.SetInitializer(strData)
	strGlobal.SetLinkage(llvm.PrivateLinkage)
	strGlobal.SetGlobalConstant(true)

	// Get pointer to the string (GEP to first element)
	zero := llvm.ConstInt(cg.context.Int32Type(), 0, false)
	indices := []llvm.Value{zero, zero}
	ptr := llvm.ConstGEP(strData.Type(), strGlobal, indices)

	cg.globalStrings[s] = ptr
	return ptr
}

// toLLVMType converts a Lotus type string to LLVM type
func (cg *LLVMCodeGenerator) toLLVMType(typeStr string) llvm.Type {
	switch typeStr {
	case "int", "int64":
		return cg.context.Int64Type()
	case "int32":
		return cg.context.Int32Type()
	case "int16":
		return cg.context.Int16Type()
	case "int8", "char":
		return cg.context.Int8Type()
	case "bool":
		return cg.context.Int1Type()
	case "float", "float64":
		return cg.context.DoubleType()
	case "float32":
		return cg.context.FloatType()
	case "string":
		return llvm.PointerType(cg.context.Int8Type(), 0)
	case "void":
		return cg.context.VoidType()
	default:
		// Default to i64 for unknown types
		return cg.context.Int64Type()
	}
}

// tokenTypeToLLVM converts a Lotus TokenType to LLVM type
func (cg *LLVMCodeGenerator) tokenTypeToLLVM(tokenType TokenType) llvm.Type {
	switch tokenType {
	case TokenTypeInt, TokenTypeInt64:
		return cg.context.Int64Type()
	case TokenTypeInt32:
		return cg.context.Int32Type()
	case TokenTypeInt16:
		return cg.context.Int16Type()
	case TokenTypeInt8:
		return cg.context.Int8Type()
	case TokenTypeUint, TokenTypeUint64:
		return cg.context.Int64Type()
	case TokenTypeUint32:
		return cg.context.Int32Type()
	case TokenTypeUint16:
		return cg.context.Int16Type()
	case TokenTypeUint8:
		return cg.context.Int8Type()
	case TokenTypeBool:
		return cg.context.Int1Type()
	case TokenTypeFloat:
		return cg.context.DoubleType()
	case TokenTypeString:
		return llvm.PointerType(cg.context.Int8Type(), 0)
	case TokenTypeChar:
		return cg.context.Int32Type() // Unicode code point (32-bit)
	case TokenTypeVoid, TokenRet, TokenReturn:
		return cg.context.VoidType() // void return type
	default:
		// Default to i64 for unknown types (matches int behavior)
		return cg.context.Int64Type()
	}
}

// GetIR returns the generated LLVM IR as a string
func (cg *LLVMCodeGenerator) GetIR() string {
	return cg.module.String()
}

// Optimize runs LLVM optimization passes on the module
func (cg *LLVMCodeGenerator) Optimize(level int) {
	// The go-llvm bindings don't expose individual pass methods directly.
	// For now, we just verify the module and leave optimization to clang/llc.
	// In the future, we can use the new pass manager via C API.
	if level > 0 {
		// Module verification is the main thing we can do here
		_ = llvm.VerifyModule(cg.module, llvm.ReturnStatusAction)
	}
}

// CompileToObject compiles the module to an object file
func (cg *LLVMCodeGenerator) CompileToObject(targetTriple string, outputPath string) error {
	// Initialize target
	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()

	// Get target
	target, err := llvm.GetTargetFromTriple(targetTriple)
	if err != nil {
		return fmt.Errorf("failed to get target: %v", err)
	}

	// Create target machine
	tm := target.CreateTargetMachine(
		targetTriple,
		"generic",
		"",
		llvm.CodeGenLevelDefault,
		llvm.RelocPIC,
		llvm.CodeModelDefault,
	)
	defer tm.Dispose()

	// Set target triple on module
	cg.module.SetTarget(targetTriple)
	cg.module.SetDataLayout(tm.CreateTargetData().String())

	// Emit object file
	mb, err := tm.EmitToMemoryBuffer(cg.module, llvm.ObjectFile)
	if err != nil {
		return fmt.Errorf("failed to emit object file: %v", err)
	}
	defer mb.Dispose()

	// Write to file
	return writeFile(outputPath, mb.Bytes())
}

// ============================================================================
// CURRYING AND WRAPPERS
// ============================================================================

// wrapperDefs stores wrapper definitions for later application
var wrapperDefs = make(map[string]*WrapperDefinition)

// generateWrapper generates code for a wrapper definition
// wrap fn timing(fn wrapped) { ... }
func (cg *LLVMCodeGenerator) generateWrapper(w *WrapperDefinition) error {
	// Store the wrapper definition for later use when decorators are applied
	wrapperDefs[w.Name] = w
	return nil
}

// generateDecoratedFunction generates a function with decorators applied
// @timing @logging fn foo() { ... }
// This creates: foo() -> timing's body -> logging's body -> __wrapped_foo()
func (cg *LLVMCodeGenerator) generateDecoratedFunction(d *DecoratedFunction) error {
	// First, generate the underlying function with a modified name
	originalName := d.Function.Name
	wrappedName := "__wrapped_" + originalName

	// Create a copy of the function with the wrapped name
	wrappedFn := &FunctionDefinition{
		BaseNode:   d.Function.BaseNode,
		Name:       wrappedName,
		ReturnType: d.Function.ReturnType,
		Parameters: d.Function.Parameters,
		Body:       d.Function.Body,
		IsVirtual:  d.Function.IsVirtual,
		IsOverride: d.Function.IsOverride,
		IsStatic:   d.Function.IsStatic,
		Storage:    d.Function.Storage,
	}

	// Declare and generate the wrapped function
	cg.declareFunction(wrappedFn)
	if err := cg.generateFunction(wrappedFn); err != nil {
		return err
	}

	// Build the decorator chain from inside out
	// @timing @logging fn foo() means: foo calls timing's body, which calls logging's body, which calls __wrapped_foo
	// We create intermediate functions for each decorator level
	
	currentCallTarget := wrappedName
	
	// Process decorators from innermost to outermost
	// For @timing @logging fn foo():
	// - First create __logging_foo which applies logging and calls __wrapped_foo
	// - Then create foo which applies timing and calls __logging_foo
	for i := len(d.Decorators) - 1; i >= 0; i-- {
		decoratorName := d.Decorators[i]
		wrapperDef, ok := wrapperDefs[decoratorName]
		if !ok {
			return fmt.Errorf("undefined wrapper: %s", decoratorName)
		}

		// Determine the name for this level of wrapping
		var thisLevelName string
		if i == 0 {
			// Outermost decorator - use the original function name
			thisLevelName = originalName
		} else {
			// Intermediate level - create a helper function
			thisLevelName = fmt.Sprintf("__%s_%s", decoratorName, originalName)
		}

		// Create a function definition for this level
		levelFn := &FunctionDefinition{
			BaseNode:   d.Function.BaseNode,
			Name:       thisLevelName,
			ReturnType: d.Function.ReturnType,
			Parameters: d.Function.Parameters,
		}

		// Declare the function
		cg.declareFunction(levelFn)

		llvmFn, ok := cg.functions[thisLevelName]
		if !ok {
			return fmt.Errorf("function %s not declared", thisLevelName)
		}

		cg.currentFn = llvmFn
		cg.namedValues = make(map[string]LLVMVariable)

		// Create entry block
		entry := llvm.AddBasicBlock(llvmFn, "entry")
		cg.builder.SetInsertPointAtEnd(entry)

		// Allocate space for parameters
		for j, param := range d.Function.Parameters {
			paramType := cg.tokenTypeToLLVM(param.Type)
			alloca := cg.builder.CreateAlloca(paramType, param.Name)
			cg.builder.CreateStore(llvmFn.Param(j), alloca)
			cg.namedValues[param.Name] = LLVMVariable{
				Alloca:      alloca,
				ElementType: paramType,
			}
		}

		// Generate the wrapper body, replacing wrapped() calls with currentCallTarget
		for _, stmt := range wrapperDef.Body {
			if err := cg.generateDecoratorStatement(stmt, currentCallTarget, wrapperDef.WrappedArg); err != nil {
				return err
			}
		}

		// Add return if needed
		if d.Function.ReturnType == TokenTypeVoid {
			cg.builder.CreateRetVoid()
		} else {
			retType := cg.tokenTypeToLLVM(d.Function.ReturnType)
			cg.builder.CreateRet(llvm.ConstNull(retType))
		}

		// Update current call target for next decorator level
		currentCallTarget = thisLevelName
	}

	return nil
}

// generateDecoratorStatement generates a statement from a wrapper body,
// replacing calls to the wrapped function parameter with the actual wrapped function
func (cg *LLVMCodeGenerator) generateDecoratorStatement(stmt ASTNode, wrappedFnName string, wrappedArgName string) error {
	switch s := stmt.(type) {
	case *FunctionCall:
		// Check if this is calling the wrapped function
		if s.Name == wrappedArgName {
			// Call the actual wrapped function instead
			fn, ok := cg.functions[wrappedFnName]
			if !ok {
				return fmt.Errorf("undefined function: %s", wrappedFnName)
			}
			fnType := fn.GlobalValueType()
			retType := fnType.ReturnType()
			if retType.TypeKind() == llvm.VoidTypeKind {
				cg.builder.CreateCall(fnType, fn, []llvm.Value{}, "")
			} else {
				cg.builder.CreateCall(fnType, fn, []llvm.Value{}, "wrappedcall")
			}
			return nil
		}
		// Regular function call
		_, err := cg.generateCall(s)
		return err
	default:
		// For other statements, use normal generation
		return cg.generateStatement(stmt)
	}
}

// generatePartialApplication generates code for partial function application
// partial(fn_name, arg1, arg2, ...) creates a closure
func (cg *LLVMCodeGenerator) generatePartialApplication(p *PartialApplication) (llvm.Value, error) {
	// For a simple implementation, we store bound args in a heap-allocated struct
	// and return a pointer to it. When called, it retrieves the args.

	// This is a simplified implementation that works for common cases:
	// We create a small struct containing:
	// - Function pointer (as i64)
	// - Number of bound args
	// - Bound arg values

	// Allocate closure struct: [fn_ptr:8][num_bound:8][arg0:8][arg1:8]...
	numBound := len(p.BoundArgs)
	structSize := llvm.ConstInt(cg.context.Int64Type(), uint64(16+numBound*8), false)

	malloc := cg.functions["malloc"]
	closurePtr := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{structSize}, "closurealloc")
	closureInt := cg.builder.CreatePtrToInt(closurePtr, cg.context.Int64Type(), "closureint")

	// Store function pointer (as identifier for now)
	// We'll use a hash of the function name as the identifier
	fnHash := uint64(0)
	for _, c := range p.FunctionName {
		fnHash = fnHash*31 + uint64(c)
	}
	fnPtrConst := llvm.ConstInt(cg.context.Int64Type(), fnHash, false)
	fnPtrAddr := cg.builder.CreateIntToPtr(closureInt, llvm.PointerType(cg.context.Int64Type(), 0), "fnptraddr")
	cg.builder.CreateStore(fnPtrConst, fnPtrAddr)

	// Store number of bound args
	eight := llvm.ConstInt(cg.context.Int64Type(), 8, false)
	numBoundOffset := cg.builder.CreateAdd(closureInt, eight, "numboundoffset")
	numBoundPtr := cg.builder.CreateIntToPtr(numBoundOffset, llvm.PointerType(cg.context.Int64Type(), 0), "numboundptr")
	numBoundVal := llvm.ConstInt(cg.context.Int64Type(), uint64(numBound), false)
	cg.builder.CreateStore(numBoundVal, numBoundPtr)

	// Store bound arguments
	for i, arg := range p.BoundArgs {
		argVal, err := cg.generateExpression(arg)
		if err != nil {
			return llvm.Value{}, err
		}

		offset := llvm.ConstInt(cg.context.Int64Type(), uint64(16+i*8), false)
		argOffset := cg.builder.CreateAdd(closureInt, offset, fmt.Sprintf("arg%doffset", i))
		argPtr := cg.builder.CreateIntToPtr(argOffset, llvm.PointerType(cg.context.Int64Type(), 0), fmt.Sprintf("arg%dptr", i))
		cg.builder.CreateStore(argVal, argPtr)
	}

	// Register this partial for later lookup when it's called
	// Store the mapping: closureInt -> (functionName, boundArgs)
	cg.partialApplications[p.FunctionName] = p

	return closureInt, nil
}

// generatePipeExpression generates code for the pipe operator
// value |> fn becomes fn(value)
func (cg *LLVMCodeGenerator) generatePipeExpression(p *PipeExpression) (llvm.Value, error) {
	// Evaluate the left side
	leftVal, err := cg.generateExpression(p.Left)
	if err != nil {
		return llvm.Value{}, err
	}

	// Create a function call with the left value as the first argument
	call := &FunctionCall{
		Name: p.Function,
		Args: []ASTNode{&pipeValueWrapper{value: leftVal}},
	}

	return cg.generateCall(call)
}

// pipeValueWrapper wraps an already-computed LLVM value to be used as a function argument
type pipeValueWrapper struct {
	BaseNode
	value llvm.Value
}

func (p *pipeValueWrapper) astNode() {}

// Dispose cleans up LLVM resources
func (cg *LLVMCodeGenerator) Dispose() {
	// Only dispose the builder - the module is owned by the context
	// and disposing the context will clean up everything
	cg.builder.Dispose()
	// Note: Don't dispose module separately - it's owned by context
	// cg.module.Dispose() // This can cause double-free
	// Note: Don't dispose context in defer - causes issues with go-llvm
	// The context will be cleaned up when the process exits
}

// Helper to write bytes to file
func writeFile(path string, data []byte) error {
	// Use os.WriteFile in production
	return nil // Placeholder - implement with os package
}
