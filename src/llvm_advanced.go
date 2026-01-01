package main

// llvm_advanced.go - LLVM codegen for advanced Lotus features
// Includes: structs, enums, classes, arrays, pointers, ternary, compound assignment

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

// ============================================================================
// TERNARY OPERATOR
// ============================================================================

func (cg *LLVMCodeGenerator) generateTernary(t *TernaryOp) (llvm.Value, error) {
	cond, err := cg.generateExpression(t.Condition)
	if err != nil {
		return llvm.Value{}, err
	}

	// Ensure condition is i1
	if cond.Type().IntTypeWidth() != 1 {
		zero := llvm.ConstInt(cond.Type(), 0, false)
		cond = cg.builder.CreateICmp(llvm.IntNE, cond, zero, "ternarycond")
	}

	currentFn := cg.builder.GetInsertBlock().Parent()
	thenBB := cg.context.AddBasicBlock(currentFn, "ternary_then")
	elseBB := cg.context.AddBasicBlock(currentFn, "ternary_else")
	mergeBB := cg.context.AddBasicBlock(currentFn, "ternary_merge")

	cg.builder.CreateCondBr(cond, thenBB, elseBB)

	// Then value
	cg.builder.SetInsertPointAtEnd(thenBB)
	thenVal, err := cg.generateExpression(t.TrueExpr)
	if err != nil {
		return llvm.Value{}, err
	}
	thenBB = cg.builder.GetInsertBlock() // Update in case of nested expressions
	cg.builder.CreateBr(mergeBB)

	// Else value
	cg.builder.SetInsertPointAtEnd(elseBB)
	elseVal, err := cg.generateExpression(t.FalseExpr)
	if err != nil {
		return llvm.Value{}, err
	}
	elseBB = cg.builder.GetInsertBlock()
	cg.builder.CreateBr(mergeBB)

	// Merge with PHI
	cg.builder.SetInsertPointAtEnd(mergeBB)
	phi := cg.builder.CreatePHI(thenVal.Type(), "ternary_result")
	phi.AddIncoming([]llvm.Value{thenVal, elseVal}, []llvm.BasicBlock{thenBB, elseBB})

	return phi, nil
}

// ============================================================================
// COMPOUND ASSIGNMENT (+=, -=, *=, /=, etc.)
// ============================================================================

func (cg *LLVMCodeGenerator) generateCompoundAssignment(ca *CompoundAssignment) error {
	// Get the variable
	varName := ""
	if id, ok := ca.Target.(*Identifier); ok {
		varName = id.Name
	} else {
		return fmt.Errorf("compound assignment target must be an identifier")
	}

	variable, ok := cg.namedValues[varName]
	if !ok {
		return fmt.Errorf("undefined variable: %s", varName)
	}

	// Load current value
	currentVal := cg.builder.CreateLoad(variable.ElementType, variable.Alloca, varName)

	// Generate the right-hand side
	rhsVal, err := cg.generateExpression(ca.Value)
	if err != nil {
		return err
	}

	// Perform the operation
	var newVal llvm.Value
	switch ca.Operator {
	case TokenPlusEq:
		newVal = cg.builder.CreateAdd(currentVal, rhsVal, "addassign")
	case TokenMinusEq:
		newVal = cg.builder.CreateSub(currentVal, rhsVal, "subassign")
	case TokenStarEq:
		newVal = cg.builder.CreateMul(currentVal, rhsVal, "mulassign")
	case TokenSlashEq:
		newVal = cg.builder.CreateSDiv(currentVal, rhsVal, "divassign")
	case TokenPercentEq:
		newVal = cg.builder.CreateSRem(currentVal, rhsVal, "modassign")
	case TokenAmpersand:
		newVal = cg.builder.CreateAnd(currentVal, rhsVal, "andassign")
	case TokenPipe:
		newVal = cg.builder.CreateOr(currentVal, rhsVal, "orassign")
	case TokenCaret:
		newVal = cg.builder.CreateXor(currentVal, rhsVal, "xorassign")
	case TokenLShift:
		newVal = cg.builder.CreateShl(currentVal, rhsVal, "shlassign")
	case TokenRShift:
		newVal = cg.builder.CreateAShr(currentVal, rhsVal, "shrassign")
	default:
		return fmt.Errorf("unknown compound assignment operator: %v", ca.Operator)
	}

	// Store the new value
	cg.builder.CreateStore(newVal, variable.Alloca)
	return nil
}

// ============================================================================
// INCREMENT/DECREMENT (++, --)
// ============================================================================

func (cg *LLVMCodeGenerator) generateIncrementOp(inc *IncrementOp) error {
	varName := ""
	if id, ok := inc.Operand.(*Identifier); ok {
		varName = id.Name
	} else {
		return fmt.Errorf("increment/decrement operand must be an identifier")
	}

	variable, ok := cg.namedValues[varName]
	if !ok {
		return fmt.Errorf("undefined variable: %s", varName)
	}

	currentVal := cg.builder.CreateLoad(variable.ElementType, variable.Alloca, varName)
	one := llvm.ConstInt(variable.ElementType, 1, false)

	var newVal llvm.Value
	if inc.Operator == TokenPlusPlus {
		newVal = cg.builder.CreateAdd(currentVal, one, "incr")
	} else {
		newVal = cg.builder.CreateSub(currentVal, one, "decr")
	}

	cg.builder.CreateStore(newVal, variable.Alloca)
	return nil
}

// ============================================================================
// ARRAY OPERATIONS
// ============================================================================

func (cg *LLVMCodeGenerator) generateArrayDecl(a *ArrayDeclaration) error {
	elemType := cg.tokenTypeToLLVM(a.ElemType)
	
	// Get size - if Size is an IntLiteral, use its value, otherwise default to 0
	arraySize := 0
	if a.Size != nil {
		if intLit, ok := a.Size.(*IntLiteral); ok {
			arraySize = intLit.Value
		}
	}
	if arraySize == 0 && len(a.Initial) > 0 {
		arraySize = len(a.Initial)
	}
	
	arrayType := llvm.ArrayType(elemType, arraySize)

	alloca := cg.builder.CreateAlloca(arrayType, a.Name)
	cg.namedValues[a.Name] = LLVMVariable{
		Alloca:      alloca,
		ElementType: arrayType,
	}

	// Initialize elements if provided
	for i, elem := range a.Initial {
		val, err := cg.generateExpression(elem)
		if err != nil {
			return err
		}
		indices := []llvm.Value{
			llvm.ConstInt(cg.context.Int32Type(), 0, false),
			llvm.ConstInt(cg.context.Int32Type(), uint64(i), false),
		}
		elemPtr := cg.builder.CreateGEP(arrayType, alloca, indices, "elemptr")
		cg.builder.CreateStore(val, elemPtr)
	}

	return nil
}

func (cg *LLVMCodeGenerator) generateArrayAccess(a *ArrayAccess) (llvm.Value, error) {
	// Get the array variable
	varName := ""
	if id, ok := a.Array.(*Identifier); ok {
		varName = id.Name
	} else {
		return llvm.Value{}, fmt.Errorf("array access target must be an identifier")
	}

	variable, ok := cg.namedValues[varName]
	if !ok {
		return llvm.Value{}, fmt.Errorf("undefined array: %s", varName)
	}

	// Generate index
	index, err := cg.generateExpression(a.Index)
	if err != nil {
		return llvm.Value{}, err
	}

	// Handle both array types and pointer types
	var elemPtr llvm.Value
	if variable.ElementType.TypeKind() == llvm.ArrayTypeKind {
		indices := []llvm.Value{
			llvm.ConstInt(cg.context.Int32Type(), 0, false),
			index,
		}
		elemPtr = cg.builder.CreateGEP(variable.ElementType, variable.Alloca, indices, "elemptr")
	} else {
		// Pointer type - load pointer then GEP
		ptr := cg.builder.CreateLoad(variable.ElementType, variable.Alloca, "arrptr")
		elemType := cg.context.Int64Type() // Default element type
		elemPtr = cg.builder.CreateGEP(elemType, ptr, []llvm.Value{index}, "elemptr")
	}

	// Get element type for load
	elemType := variable.ElementType
	if elemType.TypeKind() == llvm.ArrayTypeKind {
		elemType = elemType.ElementType()
	}

	return cg.builder.CreateLoad(elemType, elemPtr, "elem"), nil
}

func (cg *LLVMCodeGenerator) generateArrayLiteral(a *ArrayLiteral) (llvm.Value, error) {
	if len(a.Elements) == 0 {
		return llvm.Value{}, fmt.Errorf("empty array literal")
	}

	// Determine element type from first element
	firstElem, err := cg.generateExpression(a.Elements[0])
	if err != nil {
		return llvm.Value{}, err
	}
	elemType := firstElem.Type()
	arrayType := llvm.ArrayType(elemType, len(a.Elements))

	// Allocate array
	alloca := cg.builder.CreateAlloca(arrayType, "arrayliteral")

	// Store first element
	indices := []llvm.Value{
		llvm.ConstInt(cg.context.Int32Type(), 0, false),
		llvm.ConstInt(cg.context.Int32Type(), 0, false),
	}
	elemPtr := cg.builder.CreateGEP(arrayType, alloca, indices, "elem0ptr")
	cg.builder.CreateStore(firstElem, elemPtr)

	// Store remaining elements
	for i := 1; i < len(a.Elements); i++ {
		val, err := cg.generateExpression(a.Elements[i])
		if err != nil {
			return llvm.Value{}, err
		}
		indices := []llvm.Value{
			llvm.ConstInt(cg.context.Int32Type(), 0, false),
			llvm.ConstInt(cg.context.Int32Type(), uint64(i), false),
		}
		elemPtr := cg.builder.CreateGEP(arrayType, alloca, indices, "elemptr")
		cg.builder.CreateStore(val, elemPtr)
	}

	// Return pointer to first element
	return cg.builder.CreateGEP(arrayType, alloca, []llvm.Value{
		llvm.ConstInt(cg.context.Int32Type(), 0, false),
		llvm.ConstInt(cg.context.Int32Type(), 0, false),
	}, "arrayptr"), nil
}

// ============================================================================
// STRUCT OPERATIONS
// ============================================================================

func (cg *LLVMCodeGenerator) generateStructDef(s *StructDefinition) error {
	// Create LLVM struct type
	fieldTypes := make([]llvm.Type, len(s.Fields))
	for i, field := range s.Fields {
		fieldTypes[i] = cg.tokenTypeToLLVM(field.Type)
	}

	structType := cg.context.StructType(fieldTypes, false)
	cg.structTypes[s.Name] = structType
	cg.structDefs[s.Name] = s

	return nil
}

func (cg *LLVMCodeGenerator) generateStructLiteral(s *StructLiteral) (llvm.Value, error) {
	structType, ok := cg.structTypes[s.StructName]
	if !ok {
		return llvm.Value{}, fmt.Errorf("undefined struct type: %s", s.StructName)
	}

	structDef, ok := cg.structDefs[s.StructName]
	if !ok {
		return llvm.Value{}, fmt.Errorf("struct definition not found: %s", s.StructName)
	}

	// Allocate struct
	alloca := cg.builder.CreateAlloca(structType, s.StructName+"_inst")

	// Initialize fields
	for fieldName, valueExpr := range s.Fields {
		// Find field index
		fieldIdx := -1
		for i, field := range structDef.Fields {
			if field.Name == fieldName {
				fieldIdx = i
				break
			}
		}
		if fieldIdx == -1 {
			return llvm.Value{}, fmt.Errorf("unknown field: %s in struct %s", fieldName, s.StructName)
		}

		val, err := cg.generateExpression(valueExpr)
		if err != nil {
			return llvm.Value{}, err
		}

		fieldPtr := cg.builder.CreateStructGEP(structType, alloca, fieldIdx, fieldName+"_ptr")
		cg.builder.CreateStore(val, fieldPtr)
	}

	return alloca, nil
}

func (cg *LLVMCodeGenerator) generateFieldAccess(f *FieldAccess) (llvm.Value, error) {
	// Generate the object expression
	obj, err := cg.generateExpression(f.Object)
	if err != nil {
		return llvm.Value{}, err
	}

	// Get the struct type name from the object
	// For now, we need to look up based on the identifier
	var structName string
	if id, ok := f.Object.(*Identifier); ok {
		// Look up the variable's type
		if variable, ok := cg.namedValues[id.Name]; ok {
			// Try to find the struct type
			for name, sType := range cg.structTypes {
				if sType == variable.ElementType {
					structName = name
					break
				}
			}
		}
	}

	if structName == "" {
		// Try class types too
		for name := range cg.classDefs {
			structName = name
			break
		}
	}

	structDef, ok := cg.structDefs[structName]
	if !ok {
		// Try class
		classDef, ok := cg.classDefs[structName]
		if !ok {
			return llvm.Value{}, fmt.Errorf("cannot determine struct type for field access")
		}
		// Handle class field access
		return cg.generateClassFieldAccess(obj, classDef, f.FieldName, f.IsPointer)
	}

	// Find field index
	fieldIdx := -1
	var fieldType llvm.Type
	for i, field := range structDef.Fields {
		if field.Name == f.FieldName {
			fieldIdx = i
			fieldType = cg.tokenTypeToLLVM(field.Type)
			break
		}
	}
	if fieldIdx == -1 {
		return llvm.Value{}, fmt.Errorf("unknown field: %s", f.FieldName)
	}

	structType := cg.structTypes[structName]

	// If pointer access (->), dereference first
	if f.IsPointer {
		obj = cg.builder.CreateLoad(llvm.PointerType(structType, 0), obj, "deref")
	}

	fieldPtr := cg.builder.CreateStructGEP(structType, obj, fieldIdx, f.FieldName+"_ptr")
	return cg.builder.CreateLoad(fieldType, fieldPtr, f.FieldName), nil
}

func (cg *LLVMCodeGenerator) generateClassFieldAccess(obj llvm.Value, classDef *ClassDefinition, fieldName string, isPointer bool) (llvm.Value, error) {
	classType := cg.classTypes[classDef.Name]

	// Find field index
	fieldIdx := -1
	var fieldType llvm.Type
	for i, field := range classDef.Fields {
		if field.Name == fieldName {
			fieldIdx = i
			fieldType = cg.tokenTypeToLLVM(field.Type)
			break
		}
	}
	if fieldIdx == -1 {
		return llvm.Value{}, fmt.Errorf("unknown field: %s in class %s", fieldName, classDef.Name)
	}

	if isPointer {
		obj = cg.builder.CreateLoad(llvm.PointerType(classType, 0), obj, "deref")
	}

	fieldPtr := cg.builder.CreateStructGEP(classType, obj, fieldIdx, fieldName+"_ptr")
	return cg.builder.CreateLoad(fieldType, fieldPtr, fieldName), nil
}

// ============================================================================
// ENUM OPERATIONS
// ============================================================================

func (cg *LLVMCodeGenerator) generateEnumDef(e *EnumDefinition) error {
	// Auto-assign values if not explicitly set
	nextValue := 0
	for i := range e.Values {
		if e.Values[i].Value == 0 && i > 0 {
			e.Values[i].Value = nextValue
		}
		nextValue = e.Values[i].Value + 1
	}

	cg.enumDefs[e.Name] = e
	return nil
}

func (cg *LLVMCodeGenerator) generateEnumLiteral(e *EnumLiteral) (llvm.Value, error) {
	enumDef, ok := cg.enumDefs[e.EnumName]
	if !ok {
		return llvm.Value{}, fmt.Errorf("undefined enum: %s", e.EnumName)
	}

	for _, val := range enumDef.Values {
		if val.Name == e.ValueName {
			return llvm.ConstInt(cg.context.Int64Type(), uint64(val.Value), false), nil
		}
	}

	return llvm.Value{}, fmt.Errorf("unknown enum value: %s.%s", e.EnumName, e.ValueName)
}

// ============================================================================
// CLASS OPERATIONS
// ============================================================================

func (cg *LLVMCodeGenerator) generateClassDef(c *ClassDefinition) error {
	// Create LLVM struct type for class (includes vtable pointer if has virtual methods)
	fieldTypes := make([]llvm.Type, len(c.Fields))
	for i, field := range c.Fields {
		fieldTypes[i] = cg.tokenTypeToLLVM(field.Type)
	}

	classType := cg.context.StructType(fieldTypes, false)
	cg.classTypes[c.Name] = classType
	cg.classDefs[c.Name] = c

	// Generate methods as functions with implicit 'this' parameter
	for _, method := range c.Methods {
		methodName := fmt.Sprintf("_class_%s_%s", c.Name, method.Name)

		// Build parameter types: 'this' pointer + explicit params
		paramTypes := make([]llvm.Type, len(method.Params)+1)
		paramTypes[0] = llvm.PointerType(classType, 0) // 'this' pointer
		for i, param := range method.Params {
			paramTypes[i+1] = cg.tokenTypeToLLVM(param.Type)
		}

		retType := cg.tokenTypeToLLVM(method.ReturnType)
		fnType := llvm.FunctionType(retType, paramTypes, false)
		fn := llvm.AddFunction(cg.module, methodName, fnType)
		cg.functions[methodName] = fn

		// Generate method body
		entry := llvm.AddBasicBlock(fn, "entry")
		cg.builder.SetInsertPointAtEnd(entry)

		// Store 'this' pointer
		oldNamedValues := cg.namedValues
		cg.namedValues = make(map[string]LLVMVariable)

		thisAlloca := cg.builder.CreateAlloca(llvm.PointerType(classType, 0), "this")
		cg.builder.CreateStore(fn.Param(0), thisAlloca)
		cg.namedValues["this"] = LLVMVariable{Alloca: thisAlloca, ElementType: llvm.PointerType(classType, 0)}

		// Store other parameters
		for i, param := range method.Params {
			paramType := cg.tokenTypeToLLVM(param.Type)
			alloca := cg.builder.CreateAlloca(paramType, param.Name)
			cg.builder.CreateStore(fn.Param(i+1), alloca)
			cg.namedValues[param.Name] = LLVMVariable{Alloca: alloca, ElementType: paramType}
		}

		// Generate method body
		for _, stmt := range method.Body {
			if err := cg.generateStatement(stmt); err != nil {
				return err
			}
		}

		// Add implicit return if needed
		lastBlock := cg.builder.GetInsertBlock()
		if lastBlock.LastInstruction().IsNil() || lastBlock.LastInstruction().InstructionOpcode() != llvm.Ret {
			if method.ReturnType == TokenTypeVoid {
				cg.builder.CreateRetVoid()
			} else {
				cg.builder.CreateRet(llvm.ConstNull(retType))
			}
		}

		cg.namedValues = oldNamedValues
	}

	return nil
}

func (cg *LLVMCodeGenerator) generateClassLiteral(c *ClassLiteral) (llvm.Value, error) {
	classDef, ok := cg.classDefs[c.ClassName]
	if !ok {
		return llvm.Value{}, fmt.Errorf("undefined class: %s", c.ClassName)
	}

	classType := cg.classTypes[c.ClassName]

	// Calculate size based on fields
	var typeSize uint64 = 0
	for _, field := range classDef.Fields {
		typeSize += uint64(GetTypeSize(field.Type))
	}
	if typeSize == 0 {
		typeSize = 8 // Minimum size
	}
	sizeVal := llvm.ConstInt(cg.context.Int64Type(), typeSize, false)

	malloc := cg.functions["malloc"]
	if malloc.IsNil() {
		mallocType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{cg.context.Int64Type()},
			false,
		)
		malloc = llvm.AddFunction(cg.module, "malloc", mallocType)
		malloc.SetLinkage(llvm.ExternalLinkage)
		cg.functions["malloc"] = malloc
	}

	rawPtr := cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{sizeVal}, "classalloc")
	instancePtr := cg.builder.CreateBitCast(rawPtr, llvm.PointerType(classType, 0), "instance")

	// Initialize fields
	for fieldName, valueExpr := range c.Fields {
		fieldIdx := -1
		for i, field := range classDef.Fields {
			if field.Name == fieldName {
				fieldIdx = i
				break
			}
		}
		if fieldIdx == -1 {
			return llvm.Value{}, fmt.Errorf("unknown field: %s in class %s", fieldName, c.ClassName)
		}

		val, err := cg.generateExpression(valueExpr)
		if err != nil {
			return llvm.Value{}, err
		}

		fieldPtr := cg.builder.CreateStructGEP(classType, instancePtr, fieldIdx, fieldName+"_ptr")
		cg.builder.CreateStore(val, fieldPtr)
	}

	return instancePtr, nil
}

func (cg *LLVMCodeGenerator) generateMethodCall(m *MethodCall) (llvm.Value, error) {
	// Generate object expression
	obj, err := cg.generateExpression(m.Object)
	if err != nil {
		return llvm.Value{}, err
	}

	// Determine class name
	var className string
	if _, ok := m.Object.(*Identifier); ok {
		// Look up variable type
		for name := range cg.classDefs {
			className = name // For now, use first class found
			break
		}
	}

	if className == "" {
		return llvm.Value{}, fmt.Errorf("cannot determine class for method call")
	}

	methodName := fmt.Sprintf("_class_%s_%s", className, m.MethodName)
	method, ok := cg.functions[methodName]
	if !ok {
		return llvm.Value{}, fmt.Errorf("undefined method: %s.%s", className, m.MethodName)
	}

	// If pointer access (->), dereference
	classType := cg.classTypes[className]
	if m.IsPointer {
		obj = cg.builder.CreateLoad(llvm.PointerType(classType, 0), obj, "deref")
	}

	// Build arguments: 'this' + explicit args
	args := make([]llvm.Value, len(m.Args)+1)
	args[0] = obj
	for i, arg := range m.Args {
		val, err := cg.generateExpression(arg)
		if err != nil {
			return llvm.Value{}, err
		}
		args[i+1] = val
	}

	// Call method
	retType := method.GlobalValueType().ReturnType()
	callName := ""
	if retType.TypeKind() != llvm.VoidTypeKind {
		callName = m.MethodName + "_ret"
	}

	return cg.builder.CreateCall(method.GlobalValueType(), method, args, callName), nil
}

// ============================================================================
// POINTER OPERATIONS
// ============================================================================

func (cg *LLVMCodeGenerator) generateReference(r *Reference) (llvm.Value, error) {
	// Get address of variable
	if id, ok := r.Target.(*Identifier); ok {
		variable, ok := cg.namedValues[id.Name]
		if !ok {
			return llvm.Value{}, fmt.Errorf("undefined variable: %s", id.Name)
		}
		return variable.Alloca, nil
	}
	return llvm.Value{}, fmt.Errorf("can only take address of identifiers")
}

func (cg *LLVMCodeGenerator) generateDereference(d *Dereference) (llvm.Value, error) {
	ptr, err := cg.generateExpression(d.Pointer)
	if err != nil {
		return llvm.Value{}, err
	}

	// Determine pointed-to type
	elemType := cg.context.Int64Type() // Default
	if ptr.Type().TypeKind() == llvm.PointerTypeKind {
		// For opaque pointers, we need to infer the type
		elemType = cg.context.Int64Type()
	}

	return cg.builder.CreateLoad(elemType, ptr, "deref"), nil
}

func (cg *LLVMCodeGenerator) generateSizeof(s *SizeofExpr) (llvm.Value, error) {
	// SizeofExpr has TypeOrExpr which can be an identifier (type name) or expression
	var size uint64 = 8 // Default to pointer size
	
	if s.TypeOrExpr != nil {
		if id, ok := s.TypeOrExpr.(*Identifier); ok {
			// sizeof(type)
			switch id.Name {
			case "int", "int64", "uint64":
				size = 8
			case "int32", "uint32", "float":
				size = 4
			case "int16", "uint16":
				size = 2
			case "int8", "uint8", "bool":
				size = 1
			case "string":
				size = 8 // pointer size
			default:
				// Check structs and classes
				if structDef, ok := cg.structDefs[id.Name]; ok {
					size = 0
					for _, field := range structDef.Fields {
						size += uint64(GetTypeSize(field.Type))
					}
				} else if classDef, ok := cg.classDefs[id.Name]; ok {
					size = 0
					for _, field := range classDef.Fields {
						size += uint64(GetTypeSize(field.Type))
					}
				} else {
					return llvm.Value{}, fmt.Errorf("unknown type for sizeof: %s", id.Name)
				}
			}
		} else {
			// sizeof(expr) - evaluate and get type size
			val, err := cg.generateExpression(s.TypeOrExpr)
			if err != nil {
				return llvm.Value{}, err
			}
			// Get size based on LLVM type
			switch val.Type().TypeKind() {
			case llvm.IntegerTypeKind:
				size = uint64(val.Type().IntTypeWidth() / 8)
			case llvm.DoubleTypeKind:
				size = 8
			case llvm.FloatTypeKind:
				size = 4
			case llvm.PointerTypeKind:
				size = 8
			default:
				size = 8
			}
		}
	} else {
		return llvm.Value{}, fmt.Errorf("sizeof requires type or expression")
	}

	return llvm.ConstInt(cg.context.Int64Type(), size, false), nil
}

func (cg *LLVMCodeGenerator) generateMallocExpr(m *MallocCall) (llvm.Value, error) {
	size, err := cg.generateExpression(m.Size)
	if err != nil {
		return llvm.Value{}, err
	}

	malloc := cg.functions["malloc"]
	if malloc.IsNil() {
		mallocType := llvm.FunctionType(
			llvm.PointerType(cg.context.Int8Type(), 0),
			[]llvm.Type{cg.context.Int64Type()},
			false,
		)
		malloc = llvm.AddFunction(cg.module, "malloc", mallocType)
		malloc.SetLinkage(llvm.ExternalLinkage)
		cg.functions["malloc"] = malloc
	}

	return cg.builder.CreateCall(malloc.GlobalValueType(), malloc, []llvm.Value{size}, "mallocptr"), nil
}

func (cg *LLVMCodeGenerator) generateFreeStmt(f *FreeCall) error {
	ptr, err := cg.generateExpression(f.Pointer)
	if err != nil {
		return err
	}

	free := cg.functions["free"]
	if free.IsNil() {
		freeType := llvm.FunctionType(
			cg.context.VoidType(),
			[]llvm.Type{llvm.PointerType(cg.context.Int8Type(), 0)},
			false,
		)
		free = llvm.AddFunction(cg.module, "free", freeType)
		free.SetLinkage(llvm.ExternalLinkage)
		cg.functions["free"] = free
	}

	// Cast to void* if needed
	if ptr.Type() != llvm.PointerType(cg.context.Int8Type(), 0) {
		ptr = cg.builder.CreateBitCast(ptr, llvm.PointerType(cg.context.Int8Type(), 0), "freecast")
	}

	cg.builder.CreateCall(free.GlobalValueType(), free, []llvm.Value{ptr}, "")
	return nil
}
