package main

import "fmt"

// FunctionDefinition represents a user-defined function
type FunctionDefinition struct {
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

// generateFunctionDefinition processes a function definition and stores it
func (cg *CodeGenerator) generateFunctionDefinition(funcDef *FunctionDefinition) {
	// Register the function for later use
	UserDefinedFunctions[funcDef.Name] = funcDef

	// Generate the function assembly
	funcLabel := cg.getFunctionLabel(funcDef.Name)
	returnLabel := cg.getLabel("return")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", funcLabel))
	cg.textSection.WriteString("    # Function prologue\n")
	cg.textSection.WriteString("    pushq %rbp\n")
	cg.textSection.WriteString("    movq %rsp, %rbp\n")

	// Save current variable state and create new scope
	savedVars := cg.variables
	cg.variables = make(map[string]Variable)
	savedStackOffset := cg.stackOffset
	cg.stackOffset = 0

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
	cg.textSection.WriteString("    popq %rbp\n")
	cg.textSection.WriteString("    ret\n")

	// Restore variable scope
	cg.variables = savedVars
	cg.stackOffset = savedStackOffset
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
