package main

import (
	"fmt"
	"os"
)

// PrintFunction represents a print function that can be called in Lotus code
type PrintFunction struct {
	Name     string
	NumArgs  int // -1 for variadic
	ArgTypes []TokenType
	CodeGen  func(*CodeGenerator, []ASTNode) // Code generation function
}

// RegisteredPrintFunctions maps function names to their implementations
var RegisteredPrintFunctions = map[string]*PrintFunction{
	"printf":   {Name: "printf", NumArgs: -1, CodeGen: generatePrintfCode},
	"println":  {Name: "println", NumArgs: -1, CodeGen: generatePrintlnCode},
	"fprintf":  {Name: "fprintf", NumArgs: -1, CodeGen: generateFprintfCode},
	"sprint":   {Name: "sprint", NumArgs: -1, CodeGen: generateSprintCode},
	"sprintf":  {Name: "sprintf", NumArgs: -1, CodeGen: generateSprintfCode},
	"sprintln": {Name: "sprintln", NumArgs: -1, CodeGen: generateSprintlnCode},
	"fatalf":   {Name: "fatalf", NumArgs: -1, CodeGen: generateFatalfCode},
	"fatalln":  {Name: "fatalln", NumArgs: -1, CodeGen: generateFatallnCode},
	"logf":     {Name: "logf", NumArgs: -1, CodeGen: generateLogfCode},
	"logln":    {Name: "logln", NumArgs: -1, CodeGen: generateLoglnCode},
}

// generatePrintfCode generates assembly for printf(str) - outputs to stdout via write() syscall
func generatePrintfCode(cg *CodeGenerator, args []ASTNode) {
	if len(args) < 1 {
		return
	}

	cg.textSection.WriteString("    # Printf\n")

	// Handle string literal
	if lit, ok := args[0].(*StringLiteral); ok {
		label := fmt.Sprintf(".str%d", cg.stringCount)
		length := len(lit.Value)
		cg.stringCount++
		escapedStr := escapeAssemblyString(lit.Value)
		cg.dataSection.WriteString(fmt.Sprintf("%s:\n    .asciz \"%s\"\n", label, escapedStr))
		// write() syscall: rax=1, rdi=fd(1), rsi=buf, rdx=count
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rsi\n", label))
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdx\n", length))
		cg.textSection.WriteString("    movq $1, %rdi\n") // stdout fd
		cg.textSection.WriteString("    movq $1, %rax\n") // write syscall
		cg.textSection.WriteString("    syscall\n")
		return
	}

	// Handle identifier (variable reference)
	if id, ok := args[0].(*Identifier); ok {
		if v, exists := cg.variables[id.Name]; exists {
			cg.textSection.WriteString("    # Printf with variable\n")
			// Load string pointer from stack
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rsi\n", v.Offset))
			// Use tracked string length if available, otherwise use max
			length := 100
			if len, ok := cg.stringLengths[id.Name]; ok {
				length = len
			}
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdx\n", length))
			cg.textSection.WriteString("    movq $1, %rdi\n") // stdout fd
			cg.textSection.WriteString("    movq $1, %rax\n") // write syscall
			cg.textSection.WriteString("    syscall\n")
		}
	}
}

// generatePrintlnCode generates assembly for println(str) - outputs to stdout with newline
func generatePrintlnCode(cg *CodeGenerator, args []ASTNode) {
	// Print the content
	generatePrintfCode(cg, args)

	// Add newline
	cg.textSection.WriteString("    # Println - add newline\n")
	label := fmt.Sprintf(".newline%d", cg.stringCount)
	cg.stringCount++
	cg.dataSection.WriteString(fmt.Sprintf("%s:\n    .asciz \"\\n\"\n", label))
	cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rsi\n", label))
	cg.textSection.WriteString("    movq $1, %rdx\n")
	cg.textSection.WriteString("    movq $1, %rdi\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// generateFprintfCode generates assembly for fprintf(fd, str)
func generateFprintfCode(cg *CodeGenerator, args []ASTNode) {
	// fprintf requires file descriptor as first argument
	// For now, just call printf (future: handle fd parameter)
	generatePrintfCode(cg, args[1:])
}

// generateSprintCode generates code for sprint(str) - returns string
// Note: This is complex for pure assembly, would need buffer management
func generateSprintCode(cg *CodeGenerator, args []ASTNode) {
	cg.textSection.WriteString("    # Sprint - returning string buffer (stub)\n")
	cg.textSection.WriteString("    # This requires dynamic memory allocation\n")
}

// generateSprintfCode generates code for sprintf(format, args...)
// Note: This is complex for pure assembly, would need formatting
func generateSprintfCode(cg *CodeGenerator, args []ASTNode) {
	cg.textSection.WriteString("    # Sprintf - returning formatted string (stub)\n")
	cg.textSection.WriteString("    # This requires printf-like formatting logic\n")
}

// generateSprintlnCode generates code for sprintln(args...)
func generateSprintlnCode(cg *CodeGenerator, args []ASTNode) {
	cg.textSection.WriteString("    # Sprintln - returning string with newline (stub)\n")
}

// generateFatalfCode generates code for fatalf(format, args...) - exits with error
func generateFatalfCode(cg *CodeGenerator, args []ASTNode) {
	// Print error message
	generatePrintfCode(cg, args)
	// Exit with code 1
	cg.textSection.WriteString("    # Fatalf - exit with code 1\n")
	cg.textSection.WriteString("    movq $60, %rax\n")
	cg.textSection.WriteString("    movq $1, %rdi\n")
	cg.textSection.WriteString("    syscall\n")
}

// generateFatallnCode generates code for fatalln(args...) - exits with error and newline
func generateFatallnCode(cg *CodeGenerator, args []ASTNode) {
	// Print error message with newline
	generatePrintlnCode(cg, args)
	// Exit with code 1
	cg.textSection.WriteString("    # Fatalln - exit with code 1\n")
	cg.textSection.WriteString("    movq $60, %rax\n")
	cg.textSection.WriteString("    movq $1, %rdi\n")
	cg.textSection.WriteString("    syscall\n")
}

// generateLogfCode generates code for logf(format, args...) - logging with newline
func generateLogfCode(cg *CodeGenerator, args []ASTNode) {
	generatePrintfCode(cg, args)
}

// generateLoglnCode generates code for logln(args...) - logging with newline
func generateLoglnCode(cg *CodeGenerator, args []ASTNode) {
	generatePrintlnCode(cg, args)
}

// Legacy runtime implementations (used when not compiling to assembly)
func printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func fprintf(w *os.File, format string, a ...interface{}) {
	fmt.Fprintf(w, format, a...)
}

func println(a ...interface{}) {
	fmt.Println(a...)
}

func sprint(a ...interface{}) string {
	return fmt.Sprint(a...)
}

func sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func sprintln(a ...interface{}) string {
	return fmt.Sprintln(a...)
}

func fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func fatalln(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func logf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func logln(a ...interface{}) {
	fmt.Println(a...)
}
