package main

import (
	"fmt"
	"strconv"
	"strings"
)

// GenerateAssembly generates x86-64 GNU assembly from tokens
func GenerateAssembly(tokens []Token) string {
	var b strings.Builder
	var dataSection strings.Builder
	var textSection strings.Builder

	variables := make(map[string]Variable)
	stackOffset := 0
	stringCount := 0
	exitCode := 0

	// Data section for strings
	dataSection.WriteString("    .section .data\n")

	// Parse tokens and generate assembly
	i := 0
	for i < len(tokens) {
		token := tokens[i]

		switch token.Type {
		case TokenTypeInt, TokenTypeString, TokenTypeBool, TokenTypeFloat:
			// Variable declaration: type name = value;
			if i+4 < len(tokens) && tokens[i+1].Type == TokenIdentifier && tokens[i+2].Type == TokenAssign {
				varName := tokens[i+1].Value
				varType := token.Type
				valueToken := tokens[i+3]

				stackOffset += 8
				variables[varName] = Variable{
					Name:   varName,
					Type:   varType,
					Offset: stackOffset,
				}

				switch varType {
				case TokenTypeInt:
					if valueToken.Type == TokenInt {
						if val, err := strconv.Atoi(valueToken.Value); err == nil {
							textSection.WriteString(fmt.Sprintf("    # int %s = %d\n", varName, val))
							textSection.WriteString(fmt.Sprintf("    movq $%d, -%d(%%rbp)\n", val, stackOffset))
						}
					}
				case TokenTypeFloat:
					if valueToken.Type == TokenFloat {
						// Store float as integer representation (simplified)
						if val, err := strconv.ParseFloat(valueToken.Value, 64); err == nil {
							intVal := int64(val * 1000) // store as fixed point
							textSection.WriteString(fmt.Sprintf("    # float %s = %s\n", varName, valueToken.Value))
							textSection.WriteString(fmt.Sprintf("    movq $%d, -%d(%%rbp)\n", intVal, stackOffset))
						}
					}
				case TokenTypeBool:
					if valueToken.Type == TokenBool {
						boolVal := 0
						if valueToken.Value == "true" {
							boolVal = 1
						}
						textSection.WriteString(fmt.Sprintf("    # bool %s = %s\n", varName, valueToken.Value))
						textSection.WriteString(fmt.Sprintf("    movq $%d, -%d(%%rbp)\n", boolVal, stackOffset))
					}
				case TokenTypeString:
					if valueToken.Type == TokenString {
						label := fmt.Sprintf(".str%d", stringCount)
						stringCount++
						dataSection.WriteString(fmt.Sprintf("%s:\n    .asciz \"%s\"\n", label, valueToken.Value))
						textSection.WriteString(fmt.Sprintf("    # string %s = \"%s\"\n", varName, valueToken.Value))
						textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rax\n", label))
						textSection.WriteString(fmt.Sprintf("    movq %%rax, -%d(%%rbp)\n", stackOffset))
					}
				}

				i += 5 // skip type, name, =, value, ;
				continue
			}

		case TokenRet:
			// ret value;
			if i+1 < len(tokens) && tokens[i+1].Type == TokenInt {
				if val, err := strconv.Atoi(tokens[i+1].Value); err == nil {
					exitCode = val
				}
			}
			i++

		default:
			i++
		}
	}

	// Build final assembly
	b.WriteString(dataSection.String())
	b.WriteString("\n")
	b.WriteString("    .global _start\n")
	b.WriteString("    .text\n")
	b.WriteString("_start:\n")
	b.WriteString("    # Set up stack frame\n")
	b.WriteString("    pushq %rbp\n")
	b.WriteString("    movq %rsp, %rbp\n")
	if stackOffset > 0 {
		b.WriteString(fmt.Sprintf("    subq $%d, %%rsp  # Allocate stack space for variables\n", (stackOffset+15)&^15))
	}
	b.WriteString("\n")
	b.WriteString(textSection.String())
	b.WriteString("\n")
	b.WriteString("    # Exit\n")
	b.WriteString("    movq $60, %rax\n")
	b.WriteString(fmt.Sprintf("    movq $%d, %%rdi\n", exitCode))
	b.WriteString("    syscall\n")

	return b.String()
}
