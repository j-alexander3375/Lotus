package main

import (
	"fmt"
	"os"
	"strings"
)

// PrintFunction represents a print function that can be called in Lotus code
type PrintFunction struct {
	Name     string
	NumArgs  int // -1 for variadic
	ArgTypes []TokenType
	CodeGen  func(*CodeGenerator, []ASTNode) // Code generation function
}

// RegisteredPrintFunctions maps function names to their implementations
var RegisteredPrintFunctions = make(map[string]*PrintFunction)

func init() {
	RegisteredPrintFunctions["printf"] = &PrintFunction{Name: "printf", NumArgs: -1, CodeGen: generatePrintfCode}
	RegisteredPrintFunctions["println"] = &PrintFunction{Name: "println", NumArgs: -1, CodeGen: generatePrintlnCode}
	RegisteredPrintFunctions["fprintf"] = &PrintFunction{Name: "fprintf", NumArgs: -1, CodeGen: generateFprintfCode}
	RegisteredPrintFunctions["sprint"] = &PrintFunction{Name: "sprint", NumArgs: -1, CodeGen: generateSprintCode}
	RegisteredPrintFunctions["sprintf"] = &PrintFunction{Name: "sprintf", NumArgs: -1, CodeGen: generateSprintfCode}
	RegisteredPrintFunctions["sprintln"] = &PrintFunction{Name: "sprintln", NumArgs: -1, CodeGen: generateSprintlnCode}
	RegisteredPrintFunctions["fatalf"] = &PrintFunction{Name: "fatalf", NumArgs: -1, CodeGen: generateFatalfCode}
	RegisteredPrintFunctions["fatalln"] = &PrintFunction{Name: "fatalln", NumArgs: -1, CodeGen: generateFatallnCode}
	RegisteredPrintFunctions["logf"] = &PrintFunction{Name: "logf", NumArgs: -1, CodeGen: generateLogfCode}
	RegisteredPrintFunctions["logln"] = &PrintFunction{Name: "logln", NumArgs: -1, CodeGen: generateLoglnCode}
}

// generatePrintfCode generates assembly for printf(str) - outputs to stdout via write() syscall
func generatePrintfCode(cg *CodeGenerator, args []ASTNode) {
	if len(args) < 1 {
		return
	}

	cg.textSection.WriteString("    # Printf\n")

	// Format-aware path for literal strings with %%, %d, %s, %b, %o, %x, %X, %c, %q, %v
	if lit, ok := args[0].(*StringLiteral); ok {
		fmtStr := lit.Value
		textParts, placeholders := parsePlaceholders(fmtStr)
		if len(placeholders) == len(args)-1 {
			argIdx := 1
			for i, text := range textParts {
				if text != "" {
					lbl, ln := emitStringLiteral(cg, text)
					emitWriteLiteral(cg, lbl, ln)
				}
				if i < len(placeholders) {
					switch placeholders[i] {
					case 'd':
						emitPrintIntBase(cg, args[argIdx], 10, false, false)
					case 'b':
						emitPrintIntBase(cg, args[argIdx], 2, false, false)
					case 'o':
						emitPrintIntBase(cg, args[argIdx], 8, false, false)
					case 'x':
						emitPrintIntBase(cg, args[argIdx], 16, false, false)
					case 'X':
						emitPrintIntBase(cg, args[argIdx], 16, true, false)
					case 'c':
						emitPrintIntBase(cg, args[argIdx], 10, false, true)
					case 's':
						emitPrintString(cg, args[argIdx])
					case 'q':
						emitPrintStringQuoted(cg, args[argIdx])
					case 'v':
						emitPrintValue(cg, args[argIdx])
					}
					argIdx++
				}
			}
			return
		}
		// Fallback: just print the literal with no substitutions
		lbl, ln := emitStringLiteral(cg, fmtStr)
		emitWriteLiteral(cg, lbl, ln)
		return
	}

	// Handle identifier (variable reference)
	if id, ok := args[0].(*Identifier); ok {
		if v, exists := cg.variables[id.Name]; exists {
			if v.Type == TokenTypeString {
				// Print string variable
				cg.textSection.WriteString("    # Printf with string variable\n")
				cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rsi\n", v.Offset))
				length := 100
				if l, ok := cg.stringLengths[id.Name]; ok {
					length = l
				}
				cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdx\n", length))
				cg.textSection.WriteString("    movq $1, %rdi\n")
				cg.textSection.WriteString("    movq $1, %rax\n")
				cg.textSection.WriteString("    syscall\n")
			} else if isIntegerType(v.Type) {
				// Print integer variable - convert to string first
				emitPrintInteger(cg, v)
			} else {
				// For other types, generate a comment
				cg.textSection.WriteString(fmt.Sprintf("    # Printf with %v variable (not yet supported)\n", v.Type))
			}
		}
	}
}

// isIntegerType checks if a token type is an integer type
func isIntegerType(t TokenType) bool {
	switch t {
	case TokenTypeInt, TokenTypeInt8, TokenTypeInt16, TokenTypeInt32, TokenTypeInt64,
		TokenTypeUint, TokenTypeUint8, TokenTypeUint16, TokenTypeUint32, TokenTypeUint64:
		return true
	default:
		return false
	}
}

// emitPrintInteger generates code to print an integer variable
// This converts the integer to a decimal string and prints it
func emitPrintInteger(cg *CodeGenerator, v Variable) {
	cg.textSection.WriteString("    # Printf with integer variable - convert to string\n")

	// We'll use a buffer on the stack to store the converted number
	// Use offset 256 - 32 = 224 to avoid colliding with variables
	bufferOffset := 224 // Use a fixed safe location in the allocated stack

	cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rax\n", v.Offset))    // Load integer into rax
	cg.textSection.WriteString(fmt.Sprintf("    leaq -%d(%%rbp), %%r9\n", bufferOffset)) // r9 = buffer start

	// Handle negative numbers
	cg.textSection.WriteString("    movq $0, %rcx\n") // total length
	cg.textSection.WriteString("    cmpq $0, %rax\n")
	labelPositive := cg.getLabel("int_positive")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", labelPositive))

	// Negative number - write minus sign at buffer start
	cg.textSection.WriteString("    movb $45, (%r9)\n") // '-' character (ASCII 45)
	cg.textSection.WriteString("    movq $1, %rcx\n")   // Mark that we have a minus sign
	cg.textSection.WriteString("    negq %rax\n")       // Make positive

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelPositive))

	// Now convert positive number to string
	// Buffer position for digits: after minus sign (if any)
	cg.textSection.WriteString("    movq %r9, %rsi\n")  // rsi = buffer start
	cg.textSection.WriteString("    addq %rcx, %rsi\n") // rsi = buffer start + minus sign length
	cg.textSection.WriteString("    movq %rsi, %r10\n") // r10 = digit start
	cg.textSection.WriteString("    movq $0, %r8\n")    // r8 = digit count

	// If number is 0, handle specially
	labelCheckZero := cg.getLabel("int_check_zero")
	labelNotZero := cg.getLabel("int_not_zero")
	cg.textSection.WriteString("    cmpq $0, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", labelNotZero))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelCheckZero))
	cg.textSection.WriteString("    movb $48, (%rsi)\n") // '0' character
	cg.textSection.WriteString("    movq $1, %r8\n")     // digit count = 1
	labelSkipConvert := cg.getLabel("int_skip_convert")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", labelSkipConvert))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelNotZero))

	// Convert to string (digits stored in reverse)
	labelDigitLoop := cg.getLabel("int_digit_loop")
	labelDigitLoopEnd := cg.getLabel("int_digit_loop_end")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelDigitLoop))
	cg.textSection.WriteString("    cmpq $0, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", labelDigitLoopEnd))

	cg.textSection.WriteString("    movq $10, %rbx\n")
	cg.textSection.WriteString("    cqo\n")
	cg.textSection.WriteString("    idivq %rbx\n")       // rax = rax / 10, rdx = rax % 10
	cg.textSection.WriteString("    addb $48, %dl\n")    // convert to ASCII
	cg.textSection.WriteString("    movb %dl, (%rsi)\n") // store digit
	cg.textSection.WriteString("    addq $1, %rsi\n")    // move buffer pointer
	cg.textSection.WriteString("    addq $1, %r8\n")     // increment digit count
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", labelDigitLoop))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelDigitLoopEnd))

	// Now reverse the digits in place
	// r10 = start of digits, rsi = end of digits (past last digit)
	// We need to reverse from r10 to rsi-1
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelSkipConvert))

	labelReverseStart := cg.getLabel("int_reverse_start")
	labelReverseEnd := cg.getLabel("int_reverse_end")
	labelReverseLoop := cg.getLabel("int_reverse_loop")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelReverseStart))

	// Only reverse if we have more than 1 digit
	cg.textSection.WriteString("    cmpq $1, %r8\n")
	cg.textSection.WriteString(fmt.Sprintf("    jle %s\n", labelReverseEnd))

	// rdi = start of digits (in r10)
	// rsi = end of digits (in rsi - 1, but rsi was modified in the loop)
	// Recalculate: end = start + digit_count - 1
	cg.textSection.WriteString("    movq %r10, %rdi\n") // rdi = digit start
	cg.textSection.WriteString("    movq %r10, %rsi\n") // rsi = digit start
	cg.textSection.WriteString("    addq %r8, %rsi\n")  // rsi = digit start + digit count
	cg.textSection.WriteString("    subq $1, %rsi\n")   // rsi = last digit position

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelReverseLoop))
	cg.textSection.WriteString("    cmpq %rsi, %rdi\n") // if rdi >= rsi, we're done
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", labelReverseEnd))

	// Swap bytes at rdi and rsi
	cg.textSection.WriteString("    movb (%rdi), %al\n")
	cg.textSection.WriteString("    movb (%rsi), %bl\n")
	cg.textSection.WriteString("    movb %bl, (%rdi)\n")
	cg.textSection.WriteString("    movb %al, (%rsi)\n")
	cg.textSection.WriteString("    addq $1, %rdi\n")
	cg.textSection.WriteString("    subq $1, %rsi\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", labelReverseLoop))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelReverseEnd))

	// Print the converted string
	// Total length = minus sign length (in rcx) + digit count (in r8)
	cg.textSection.WriteString("    movq %r9, %rsi\n")  // buffer start
	cg.textSection.WriteString("    addq %r8, %rcx\n")  // rcx = total length (minus + digits)
	cg.textSection.WriteString("    movq %rcx, %rdx\n") // length
	cg.textSection.WriteString("    movq $1, %rdi\n")   // stdout fd
	cg.textSection.WriteString("    movq $1, %rax\n")   // write syscall
	cg.textSection.WriteString("    syscall\n")
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

// emitStringLiteral stores a string in the data section and returns its label and length
func emitStringLiteral(cg *CodeGenerator, s string) (string, int) {
	label := fmt.Sprintf(".str%d", cg.stringCount)
	cg.stringCount++
	escaped := escapeAssemblyString(s)
	cg.dataSection.WriteString(fmt.Sprintf("%s:\n    .asciz \"%s\"\n", label, escaped))
	return label, len(s)
}

// emitWriteLiteral emits a write syscall for a labeled string
func emitWriteLiteral(cg *CodeGenerator, label string, length int) {
	cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rsi\n", label))
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdx\n", length))
	cg.textSection.WriteString("    movq $1, %rdi\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// emitPrintIntBase prints an integer with the given base; if asChar is true, emit a single byte
func emitPrintIntBase(cg *CodeGenerator, expr ASTNode, base int, uppercase bool, asChar bool) {
	cg.generateExpressionToReg(expr, "rax")

	if asChar {
		bufLabel := fmt.Sprintf(".charbuf%d", cg.stringCount)
		cg.stringCount++
		cg.dataSection.WriteString(fmt.Sprintf("%s:\n    .byte 0\n", bufLabel))
		cg.textSection.WriteString("    movb %al, %dl\n")
		cg.textSection.WriteString(fmt.Sprintf("    movb %%dl, %s(%%rip)\n", bufLabel))
		emitWriteLiteral(cg, bufLabel, 1)
		return
	}

	bufLabel := fmt.Sprintf(".intbuf%d", cg.stringCount)
	cg.stringCount++
	cg.dataSection.WriteString(fmt.Sprintf("%s:\n    .space 32\n", bufLabel))

	loopLabel := cg.getLabel("itoa_loop")
	digitLabel := cg.getLabel("itoa_digit")
	storeLabel := cg.getLabel("itoa_store")
	endLabel := cg.getLabel("itoa_end")

	cg.textSection.WriteString("    # itoa into buffer\n")
	cg.textSection.WriteString("    movq %rax, %rbx\n")
	cg.textSection.WriteString("    movq $0, %r9\n")
	cg.textSection.WriteString("    cmpq $0, %rbx\n")
	cg.textSection.WriteString("    jge 1f\n")
	cg.textSection.WriteString("    negq %rbx\n")
	cg.textSection.WriteString("    movq $1, %r9\n")
	cg.textSection.WriteString("1:\n")
	cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rsi\n", bufLabel))
	cg.textSection.WriteString("    addq $31, %rsi\n")
	cg.textSection.WriteString("    movb $0, (%rsi)\n")
	cg.textSection.WriteString("    movq $0, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", base))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", loopLabel))
	cg.textSection.WriteString("    movq %rbx, %rax\n")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rax, %rbx\n")
	cg.textSection.WriteString("    movq %rdx, %rax\n")
	cg.textSection.WriteString("    cmpb $9, %al\n")
	upperChar := 'a'
	if uppercase {
		upperChar = 'A'
	}
	cg.textSection.WriteString(fmt.Sprintf("    jbe %s\n", digitLabel))
	cg.textSection.WriteString(fmt.Sprintf("    addb $%d, %%al\n", upperChar-10))
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", storeLabel))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", digitLabel))
	cg.textSection.WriteString("    addb $'0', %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", storeLabel))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", storeLabel))
	cg.textSection.WriteString("    dec %rsi\n")
	cg.textSection.WriteString("    movb %al, (%rsi)\n")
	cg.textSection.WriteString("    inc %rcx\n")
	cg.textSection.WriteString("    cmpq $0, %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", loopLabel))
	cg.textSection.WriteString("    cmpq $0, %r9\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", endLabel))
	cg.textSection.WriteString("    dec %rsi\n")
	cg.textSection.WriteString("    movb $'-', (%rsi)\n")
	cg.textSection.WriteString("    inc %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLabel))
	cg.textSection.WriteString("    movq %rcx, %rdx\n")
	cg.textSection.WriteString("    movq $1, %rdi\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// emitPrintString prints a string expression
func emitPrintString(cg *CodeGenerator, expr ASTNode) {
	switch v := expr.(type) {
	case *StringLiteral:
		lbl, ln := emitStringLiteral(cg, v.Value)
		emitWriteLiteral(cg, lbl, ln)
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString("    # print string variable\n")
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rsi\n", varInfo.Offset))
			if l, ok := cg.stringLengths[v.Name]; ok {
				cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdx\n", l))
			} else {
				// Compute length at runtime by scanning to NUL
				lLoop := cg.getLabel("slen_loop")
				lEnd := cg.getLabel("slen_end")
				cg.textSection.WriteString("    movq %rsi, %rbx\n")
				cg.textSection.WriteString("    movq $0, %rdx\n")
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
				cg.textSection.WriteString("    movzbq (%rbx), %rax\n")
				cg.textSection.WriteString("    testq %rax, %rax\n")
				cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
				cg.textSection.WriteString("    inc %rdx\n")
				cg.textSection.WriteString("    inc %rbx\n")
				cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
			}
			cg.textSection.WriteString("    movq $1, %rdi\n")
			cg.textSection.WriteString("    movq $1, %rax\n")
			cg.textSection.WriteString("    syscall\n")
		}
	case *FunctionCall:
		// Evaluate the function call to get result in rax
		cg.generateFunctionCall(v)
		// Result is in rax; treat it as a string pointer and compute length
		cg.textSection.WriteString("    # print function result string\n")
		cg.textSection.WriteString("    movq %rax, %rsi\n")
		lLoop := cg.getLabel("slen_loopc")
		lEnd := cg.getLabel("slen_endc")
		cg.textSection.WriteString("    movq %rsi, %rbx\n")
		cg.textSection.WriteString("    movq $0, %rdx\n")
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
		cg.textSection.WriteString("    movzbq (%rbx), %rax\n")
		cg.textSection.WriteString("    testq %rax, %rax\n")
		cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
		cg.textSection.WriteString("    inc %rdx\n")
		cg.textSection.WriteString("    inc %rbx\n")
		cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
		cg.textSection.WriteString("    movq $1, %rdi\n")
		cg.textSection.WriteString("    movq $1, %rax\n")
		cg.textSection.WriteString("    syscall\n")
	default:
		// Unsupported expression kind; no-op
	}
}

// emitPrintStringQuoted prints a string with surrounding quotes
func emitPrintStringQuoted(cg *CodeGenerator, expr ASTNode) {
	switch v := expr.(type) {
	case *StringLiteral:
		lbl, ln := emitStringLiteral(cg, fmt.Sprintf("\"%s\"", v.Value))
		emitWriteLiteral(cg, lbl, ln)
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString("    # print quoted string variable\n")
			// Leading quote
			lq, lqlen := emitStringLiteral(cg, "\"")
			emitWriteLiteral(cg, lq, lqlen)
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rsi\n", varInfo.Offset))
			if l, ok := cg.stringLengths[v.Name]; ok {
				cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdx\n", l))
			} else {
				lLoop := cg.getLabel("slen_loopq")
				lEnd := cg.getLabel("slen_endq")
				cg.textSection.WriteString("    movq %rsi, %rbx\n")
				cg.textSection.WriteString("    movq $0, %rdx\n")
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
				cg.textSection.WriteString("    movzbq (%rbx), %rax\n")
				cg.textSection.WriteString("    testq %rax, %rax\n")
				cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
				cg.textSection.WriteString("    inc %rdx\n")
				cg.textSection.WriteString("    inc %rbx\n")
				cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
			}
			cg.textSection.WriteString("    movq $1, %rdi\n")
			cg.textSection.WriteString("    movq $1, %rax\n")
			cg.textSection.WriteString("    syscall\n")
			// Trailing quote
			rq, rqlen := emitStringLiteral(cg, "\"")
			emitWriteLiteral(cg, rq, rqlen)
		}
	default:
		// Unsupported expression kind; no-op
	}
}

// emitPrintValue chooses int or string rendering for %v
func emitPrintValue(cg *CodeGenerator, expr ASTNode) {
	switch expr.(type) {
	case *StringLiteral, *Identifier:
		emitPrintString(cg, expr)
	default:
		emitPrintIntBase(cg, expr, 10, false, false)
	}
}

// parsePlaceholders splits a format string into text parts and placeholder runes
// Supports %%, %d, %s, %b, %o, %x, %X, %c, %q, %v
func parsePlaceholders(s string) ([]string, []rune) {
	var texts []string
	var placeholders []rune
	var sb strings.Builder
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if runes[i] == '%' {
			if i+1 < len(runes) {
				next := runes[i+1]
				if next == '%' {
					sb.WriteRune('%')
					i++
					continue
				}
				if strings.ContainsRune("dsboxXcvq", next) {
					texts = append(texts, sb.String())
					sb.Reset()
					placeholders = append(placeholders, next)
					i++
					continue
				}
			}
		}
		sb.WriteRune(runes[i])
	}
	texts = append(texts, sb.String())
	return texts, placeholders
}
