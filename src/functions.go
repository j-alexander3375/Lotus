package main

import "fmt"

// FunctionDefinition represents a user-defined function
type FunctionDefinition struct {
	BaseNode
	Name       string
	Parameters []FunctionParam
	ReturnType TokenType
	Body       []ASTNode
}

func (f *FunctionDefinition) astNode() {}

// FunctionParam represents a function parameter
type FunctionParam struct {
	Name string
	Type TokenType
}

// FunctionContext holds information about a function during code generation
type FunctionContext struct {
	Name        string
	Parameters  map[string]Variable
	LocalVars   map[string]Variable
	StackSize   int
	ReturnLabel string
}

// Global registry of function definitions
var UserDefinedFunctions = make(map[string]*FunctionDefinition)

// calculateStackSize calculates the required stack space for a function
func (cg *CodeGenerator) calculateStackSize(funcDef *FunctionDefinition) int {
	stackNeeded := 0

	// Each parameter takes 8 bytes (we pass them via registers but store on stack)
	stackNeeded += len(funcDef.Parameters) * 8

	// Count local variables in the function body
	stackNeeded += cg.countLocalVariablesInBody(funcDef.Body) * 8

	// Ensure 16-byte alignment for x86-64 ABI (add padding if needed)
	if stackNeeded%16 != 0 {
		stackNeeded += 16 - (stackNeeded % 16)
	}

	// Minimum 8192 bytes for recursive functions, maximum 1MB to prevent stack overflow
	if stackNeeded < 8192 {
		stackNeeded = 8192
	}
	if stackNeeded > 1048576 {
		stackNeeded = 1048576
	}

	return stackNeeded
}

// countLocalVariablesInBody counts unique variables declared in function body
func (cg *CodeGenerator) countLocalVariablesInBody(body []ASTNode) int {
	varNames := make(map[string]bool)
	cg.walkBodyForVars(body, varNames)
	return len(varNames)
}

// walkBodyForVars recursively finds all variable declarations
func (cg *CodeGenerator) walkBodyForVars(body []ASTNode, vars map[string]bool) {
	for _, node := range body {
		switch n := node.(type) {
		case *VariableDeclaration:
			vars[n.Name] = true
		case *IfStatement:
			cg.walkBodyForVars(n.ThenBody, vars)
			cg.walkBodyForVars(n.ElseBody, vars)
		case *WhileLoop:
			cg.walkBodyForVars(n.Body, vars)
		case *ForLoop:
			if n.Init != nil {
				if decl, ok := n.Init.(*VariableDeclaration); ok {
					vars[decl.Name] = true
				}
			}
			cg.walkBodyForVars(n.Body, vars)
		}
	}
}

// generateFunctionDefinition processes a function definition and stores it
func (cg *CodeGenerator) generateFunctionDefinition(funcDef *FunctionDefinition) {
	// Register the function for later use
	UserDefinedFunctions[funcDef.Name] = funcDef

	// Generate the function assembly
	funcLabel := cg.getFunctionLabel(funcDef.Name)
	returnLabel := cg.getLabel("return")

	// Calculate required stack size
	stackSize := cg.calculateStackSize(funcDef)

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", funcLabel))
	cg.textSection.WriteString("    # Function prologue\n")
	cg.textSection.WriteString("    pushq %rbp\n")
	cg.textSection.WriteString("    movq %rsp, %rbp\n")
	cg.textSection.WriteString(fmt.Sprintf("    subq $%d, %%rsp\n", stackSize)) // Dynamic stack allocation

	// Save current state and create new scope
	savedVars := cg.variables
	cg.variables = make(map[string]Variable)
	savedStackOffset := cg.stackOffset
	cg.stackOffset = 0
	savedInFunction := cg.inFunction
	savedReturnLbl := cg.currentFunctionReturnLbl
	cg.inFunction = true
	cg.currentFunctionReturnLbl = returnLabel

	// Set up parameters (System V AMD64 ABI: rdi, rsi, rdx, rcx, r8, r9)
	paramRegs := []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}
	for i, param := range funcDef.Parameters {
		cg.stackOffset += 8
		if i < len(paramRegs) {
			cg.textSection.WriteString(fmt.Sprintf("    movq %%%s, -%d(%%rbp)\n", paramRegs[i], cg.stackOffset))
		}
		cg.variables[param.Name] = Variable{
			Name:   param.Name,
			Type:   param.Type,
			Offset: cg.stackOffset,
		}
	}

	// Generate function body
	for _, stmt := range funcDef.Body {
		cg.generateStatement(stmt)
	}

	// Function epilogue
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", returnLabel))
	cg.textSection.WriteString("    # Function epilogue\n")
	if funcDef.Name == "main" {
		// Exit directly from main using return value in rax
		cg.textSection.WriteString("    # Exit from main\n")
		cg.textSection.WriteString("    movq %rax, %rdi\n")
		cg.textSection.WriteString("    movq $60, %rax\n")
		cg.textSection.WriteString("    syscall\n")
	} else {
		cg.textSection.WriteString("    movq %rbp, %rsp\n") // Restore stack pointer
		cg.textSection.WriteString("    popq %rbp\n")
		cg.textSection.WriteString("    ret\n")
	}

	// Restore variable scope
	cg.variables = savedVars
	cg.stackOffset = savedStackOffset
	cg.inFunction = savedInFunction
	cg.currentFunctionReturnLbl = savedReturnLbl
}

// generateUserFunctionCall generates assembly for calling a user-defined function
func (cg *CodeGenerator) generateUserFunctionCall(funcCall *FunctionCall) bool {
	_, exists := UserDefinedFunctions[funcCall.Name]
	if !exists {
		return false
	}

	// System V AMD64 ABI: rdi, rsi, rdx, rcx, r8, r9
	paramRegs := []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

	// Evaluate arguments and place in registers
	for i, arg := range funcCall.Args {
		if i >= len(paramRegs) {
			// Additional args go on stack
			cg.generateExpressionToReg(arg, "rax")
			cg.textSection.WriteString("    pushq %rax\n")
		} else {
			cg.generateExpressionToReg(arg, paramRegs[i])
		}
	}

	// Call function
	funcLabel := cg.getFunctionLabel(funcCall.Name)
	cg.textSection.WriteString(fmt.Sprintf("    call %s\n", funcLabel))

	// Clean up stack if there were extra arguments
	extraArgs := len(funcCall.Args) - len(paramRegs)
	if extraArgs > 0 {
		cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsp\n", extraArgs*8))
	}

	// Return value is in rax
	return true
}

// getFunctionLabel gets the assembly label for a function
func (cg *CodeGenerator) getFunctionLabel(funcName string) string {
	return "." + funcName
}
