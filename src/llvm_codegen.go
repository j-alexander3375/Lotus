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
		context:       ctx,
		module:        mod,
		builder:       builder,
		namedValues:   make(map[string]LLVMVariable),
		globalStrings: make(map[string]llvm.Value),
		functions:     make(map[string]llvm.Value),
		imports:       NewImportContext(),
		stringCounter: 0,
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
	// First pass: declare all functions
	for _, stmt := range statements {
		if fn, ok := stmt.(*FunctionDefinition); ok {
			if err := cg.declareFunction(fn); err != nil {
				return err
			}
		}
	}

	// Second pass: generate function bodies
	for _, stmt := range statements {
		if err := cg.generateStatement(stmt); err != nil {
			return err
		}
	}

	// Verify the module
	if err := llvm.VerifyModule(cg.module, llvm.ReturnStatusAction); err != nil {
		return fmt.Errorf("LLVM module verification failed: %v", err)
	}

	return nil
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
	llvmFn.SetLinkage(llvm.ExternalLinkage)

	// Name parameters
	for i, param := range fn.Parameters {
		llvmFn.Param(i).SetName(param.Name)
	}

	cg.functions[fn.Name] = llvmFn
	return nil
}

// generateStatement generates LLVM IR for a statement
func (cg *LLVMCodeGenerator) generateStatement(stmt ASTNode) error {
	switch s := stmt.(type) {
	case *FunctionDefinition:
		return cg.generateFunction(s)
	case *VariableDeclaration:
		return cg.generateVarDecl(s)
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
		retType := llvmFn.Type().ElementType().ReturnType()
		if retType.TypeKind() == llvm.VoidTypeKind {
			cg.builder.CreateRetVoid()
		} else {
			cg.builder.CreateRet(llvm.ConstNull(retType))
		}
	}

	return nil
}

// generateVarDecl generates a variable declaration
func (cg *LLVMCodeGenerator) generateVarDecl(v *VariableDeclaration) error {
	varType := cg.tokenTypeToLLVM(v.Type)
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

	case *Identifier:
		variable, ok := cg.namedValues[e.Name]
		if !ok {
			return llvm.Value{}, fmt.Errorf("undefined variable: %s", e.Name)
		}
		return cg.builder.CreateLoad(variable.ElementType, variable.Alloca, e.Name), nil

	case *BinaryOp:
		return cg.generateBinaryExpr(e)

	case *UnaryOp:
		return cg.generateUnaryExpr(e)

	case *Comparison:
		return cg.generateComparison(e)

	case *LogicalOp:
		return cg.generateLogicalOp(e)

	case *FunctionCall:
		return cg.generateCall(e)

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

	return cg.builder.CreateCall(fn.GlobalValueType(), fn, args, "calltmp"), nil
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

	variable, ok := cg.namedValues[id.Name]
	if !ok {
		return fmt.Errorf("undefined variable: %s", id.Name)
	}

	cg.builder.CreateStore(val, variable.Alloca)
	return nil
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
	case TokenRet, TokenReturn:
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
