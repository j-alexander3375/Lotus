package main

import "fmt"

// ClassDefinition represents a class type definition
type ClassDefinition struct {
	BaseNode
	Name    string
	Fields  []ClassField
	Methods []ClassMethod
}

func (c *ClassDefinition) astNode() {}

// ClassField represents a field in a class
type ClassField struct {
	Name   string
	Type   TokenType
	Offset int // byte offset from class instance start
}

// ClassMethod represents a method in a class
type ClassMethod struct {
	Name       string
	Params     []FunctionParam
	ReturnType TokenType
	Body       []ASTNode
	Label      string // assembly label for method
}

// ClassLiteral represents class instantiation (new ClassName)
type ClassLiteral struct {
	BaseNode
	ClassName string
	Fields    map[string]ASTNode // field name -> initial value
}

func (c *ClassLiteral) astNode() {}

// MethodCall represents a method call: obj.method() or obj->method()
type MethodCall struct {
	BaseNode
	Object     ASTNode
	MethodName string
	Args       []ASTNode
	IsPointer  bool // true for ->, false for .
}

func (m *MethodCall) astNode() {}

// ClassRegistry stores defined class types
var ClassRegistry = make(map[string]*ClassDefinition)

// generateClassDefinition registers a class type and generates method code
func (cg *CodeGenerator) generateClassDefinition(def *ClassDefinition) {
	// Calculate field offsets
	offset := 0
	for i := range def.Fields {
		def.Fields[i].Offset = offset
		fieldSize := getTypeSizeFromToken(def.Fields[i].Type)
		offset += fieldSize
	}

	// Register class in global registry
	ClassRegistry[def.Name] = def

	// Generate assembly for each method
	for i := range def.Methods {
		method := &def.Methods[i]
		method.Label = fmt.Sprintf("_class_%s_%s", def.Name, method.Name)

		// Generate method code (similar to function generation)
		cg.textSection.WriteString(fmt.Sprintf("\n.globl %s\n", method.Label))
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", method.Label))

		// Method prologue
		cg.textSection.WriteString("    pushq %rbp\n")
		cg.textSection.WriteString("    movq %rsp, %rbp\n")

		// First parameter is always 'this' pointer (in rdi)
		// Store 'this' pointer on stack
		cg.stackOffset += 8
		cg.variables["this"] = Variable{
			Name:   "this",
			Type:   TokenTypeUint64,
			Offset: cg.stackOffset,
		}
		cg.textSection.WriteString(fmt.Sprintf("    movq %%rdi, -%d(%%rbp)\n", cg.stackOffset))

		// Handle other parameters (rsi, rdx, rcx, r8, r9)
		paramRegs := []string{"rsi", "rdx", "rcx", "r8", "r9"}
		for j, param := range method.Params {
			if j < len(paramRegs) {
				cg.stackOffset += 8
				cg.variables[param.Name] = Variable{
					Name:   param.Name,
					Type:   param.Type,
					Offset: cg.stackOffset,
				}
				cg.textSection.WriteString(fmt.Sprintf("    movq %%%s, -%d(%%rbp)\n", paramRegs[j], cg.stackOffset))
			}
		}

		// Generate method body
		for _, stmt := range method.Body {
			cg.generateStatement(stmt)
		}

		// Method epilogue
		cg.textSection.WriteString("    movq %rbp, %rsp\n")
		cg.textSection.WriteString("    popq %rbp\n")
		cg.textSection.WriteString("    ret\n")

		// Clear variables for next method
		cg.variables = make(map[string]Variable)
		cg.stackOffset = 0
	}
}

// generateClassLiteral generates assembly for class instantiation
func (cg *CodeGenerator) generateClassLiteral(lit *ClassLiteral) {
	classDef, exists := ClassRegistry[lit.ClassName]
	if !exists {
		fmt.Printf("Error: class type '%s' not defined\n", lit.ClassName)
		return
	}

	// Calculate total class size
	totalSize := 0
	for _, field := range classDef.Fields {
		fieldSize := getTypeSizeFromToken(field.Type)
		totalSize += fieldSize
	}

	// Allocate memory for class instance on heap
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdi\n", totalSize))
	cg.textSection.WriteString("    call malloc@PLT\n")
	cg.textSection.WriteString("    pushq %rax\n") // Save class instance pointer

	// Initialize each field
	for fieldName, valueExpr := range lit.Fields {
		// Find field definition
		var fieldDef *ClassField
		for i := range classDef.Fields {
			if classDef.Fields[i].Name == fieldName {
				fieldDef = &classDef.Fields[i]
				break
			}
		}

		if fieldDef == nil {
			fmt.Printf("Error: field '%s' not found in class '%s'\n", fieldName, lit.ClassName)
			continue
		}

		// Evaluate field value
		cg.generateExpressionToReg(valueExpr, "rcx")

		// Get class instance pointer
		cg.textSection.WriteString("    movq (%rsp), %rax\n")

		// Store value at field offset
		fieldSize := getTypeSizeFromToken(fieldDef.Type)
		switch fieldSize {
		case 1:
			cg.textSection.WriteString(fmt.Sprintf("    movb %%cl, %d(%%rax)\n", fieldDef.Offset))
		case 2:
			cg.textSection.WriteString(fmt.Sprintf("    movw %%cx, %d(%%rax)\n", fieldDef.Offset))
		case 4:
			cg.textSection.WriteString(fmt.Sprintf("    movl %%ecx, %d(%%rax)\n", fieldDef.Offset))
		case 8:
			cg.textSection.WriteString(fmt.Sprintf("    movq %%rcx, %d(%%rax)\n", fieldDef.Offset))
		}
	}

	// Pop class instance pointer back to rax
	cg.textSection.WriteString("    popq %rax\n")
}

// generateMethodCall generates assembly for method calls
func (cg *CodeGenerator) generateMethodCall(call *MethodCall) {
	// Evaluate object to get class instance pointer (this pointer)
	cg.generateExpressionToReg(call.Object, "rdi")

	// Get class type from object (requires type tracking)
	// For now, we'll assume we know the class type

	// Find method in class definition
	// Simplified - would need proper type resolution

	// Push arguments in reverse order (after this pointer)
	if len(call.Args) > 0 {
		// Save this pointer
		cg.textSection.WriteString("    pushq %rdi\n")

		// System V AMD64 ABI: rdi (this), rsi, rdx, rcx, r8, r9, then stack
		argRegs := []string{"rsi", "rdx", "rcx", "r8", "r9"}

		for i, arg := range call.Args {
			if i < len(argRegs) {
				cg.generateExpressionToReg(arg, argRegs[i])
				cg.textSection.WriteString(fmt.Sprintf("    pushq %%%s\n", argRegs[i]))
			}
		}

		// Restore arguments in reverse
		for i := len(call.Args) - 1; i >= 0; i-- {
			if i < len(argRegs) {
				cg.textSection.WriteString(fmt.Sprintf("    popq %%%s\n", argRegs[i]))
			}
		}

		// Restore this pointer
		cg.textSection.WriteString("    popq %rdi\n")
	}

	// Call method (need to generate label from class and method name)
	// This would be looked up from the class definition
	methodLabel := fmt.Sprintf("_class_%s_%s", "UnknownClass", call.MethodName)
	cg.textSection.WriteString(fmt.Sprintf("    call %s\n", methodLabel))

	// Result is in rax
}

// getClassFieldOffset looks up field offset in class
func getClassFieldOffset(className, fieldName string) (int, bool) {
	classDef, exists := ClassRegistry[className]
	if !exists {
		return 0, false
	}

	for _, field := range classDef.Fields {
		if field.Name == fieldName {
			return field.Offset, true
		}
	}

	return 0, false
}

// Example class syntax:
// class Person {
//     int age;
//     string name;
//
//     fn init(int a, string n) {
//         this.age = a;
//         this.name = n;
//     }
//
//     fn greet() {
//         printf("Hello, I'm %s\n", this.name);
//     }
// }
//
// Usage:
// Person* p = new Person;
// p->init(25, "Alice");
// p->greet();
