package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Tokenize converts source code into a slice of tokens
func Tokenize(input string) []Token {
	var tokens []Token
	var buf strings.Builder
	runes := []rune(input)

	for i := 0; i < len(runes); i++ {
		c := runes[i]

		if c == '"' {
			// String literal
			i++
			buf.Reset()
			for i < len(runes) && runes[i] != '"' {
				if runes[i] == '\\' && i+1 < len(runes) {
					// Handle escape sequences
					i++
					switch runes[i] {
					case 'n':
						buf.WriteRune('\n')
					case 't':
						buf.WriteRune('\t')
					case '\\':
						buf.WriteRune('\\')
					case '"':
						buf.WriteRune('"')
					default:
						buf.WriteRune(runes[i])
					}
				} else {
					buf.WriteRune(runes[i])
				}
				i++
			}
			if i >= len(runes) {
				fmt.Fprintf(os.Stderr, "Unterminated string literal\n")
				return []Token{}
			}
			tokens = append(tokens, Token{Type: TokenString, Value: buf.String()})
			buf.Reset()
		} else if unicode.IsLetter(c) || c == '_' {
			// Identifier or keyword
			buf.WriteRune(c)
			i++
			for i < len(runes) && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i]) || runes[i] == '_') {
				buf.WriteRune(runes[i])
				i++
			}
			i-- // Compensate for the loop increment

			word := buf.String()
			switch word {
			case "ret":
				tokens = append(tokens, Token{Type: TokenRet, Value: ""})
			case "true", "false":
				tokens = append(tokens, Token{Type: TokenBool, Value: word})
			case "int":
				tokens = append(tokens, Token{Type: TokenTypeInt, Value: ""})
			case "string":
				tokens = append(tokens, Token{Type: TokenTypeString, Value: ""})
			case "bool":
				tokens = append(tokens, Token{Type: TokenTypeBool, Value: ""})
			case "float":
				tokens = append(tokens, Token{Type: TokenTypeFloat, Value: ""})
			default:
				tokens = append(tokens, Token{Type: TokenIdentifier, Value: word})
			}
			buf.Reset()
		} else if unicode.IsDigit(c) {
			// Number literal (int or float)
			buf.WriteRune(c)
			i++
			isFloat := false
			for i < len(runes) && (unicode.IsDigit(runes[i]) || runes[i] == '.') {
				if runes[i] == '.' {
					isFloat = true
				}
				buf.WriteRune(runes[i])
				i++
			}
			i-- // Compensate for the loop increment

			numStr := buf.String()
			if isFloat {
				tokens = append(tokens, Token{Type: TokenFloat, Value: numStr})
			} else {
				tokens = append(tokens, Token{Type: TokenInt, Value: numStr})
			}
			buf.Reset()
		} else if c == '\n' {
			tokens = append(tokens, Token{Type: TokenNewline, Value: ""})
		} else if unicode.IsSpace(c) {
			// Skip other whitespace
		} else if c == ';' {
			tokens = append(tokens, Token{Type: TokenSemi, Value: ""})
		} else if c == '=' {
			tokens = append(tokens, Token{Type: TokenAssign, Value: ""})
		} else {
			fmt.Fprintf(os.Stderr, "Unable to parse %c\n", c)
			return []Token{}
		}
	}

	tokens = append(tokens, Token{Type: TokenEOF, Value: ""})
	return tokens
}

// TokenValue returns a string representation of a token's value
func TokenValue(t Token) string {
	switch t.Type {
	case TokenRet:
		return "ret"
	case TokenInt:
		return t.Value
	case TokenSemi:
		return ";"
	case TokenEOF:
		return "EOF"
	case TokenString:
		return fmt.Sprintf("\"%s\"", t.Value)
	case TokenBool:
		return t.Value
	case TokenNewline:
		return "\\n"
	case TokenAssign:
		return "="
	case TokenFloat:
		return t.Value
	case TokenIdentifier:
		return t.Value
	case TokenTypeInt:
		return "int"
	case TokenTypeString:
		return "string"
	case TokenTypeBool:
		return "bool"
	case TokenTypeFloat:
		return "float"
	default:
		return "unknown"
	}
}
