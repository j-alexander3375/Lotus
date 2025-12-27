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

	// Track line and column for error messages
	line := 1
	col := 1
	startLine := 1
	startCol := 1

	// Helper to create a token with position
	makeToken := func(t TokenType, v string) Token {
		return Token{Type: t, Value: v, Line: startLine, Column: startCol}
	}

	for i := 0; i < len(runes); i++ {
		c := runes[i]
		startLine = line
		startCol = col

		if c == '\n' {
			tokens = append(tokens, makeToken(TokenNewline, ""))
			line++
			col = 1
			continue
		}

		if c == '"' {
			// String literal (UTF-8 aware)
			col++
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
					case 'u':
						// Unicode escape: \uXXXX (4 hex digits)
						if i+4 < len(runes) {
							hexStr := string(runes[i+1 : i+5])
							var codepoint rune
							fmt.Sscanf(hexStr, "%x", &codepoint)
							buf.WriteRune(codepoint)
							i += 4
						}
					case 'U':
						// Unicode escape: \UXXXXXXXX (8 hex digits)
						if i+8 < len(runes) {
							hexStr := string(runes[i+1 : i+9])
							var codepoint rune
							fmt.Sscanf(hexStr, "%x", &codepoint)
							buf.WriteRune(codepoint)
							i += 8
						}
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
			tokens = append(tokens, makeToken(TokenString, buf.String()))
			buf.Reset()
		} else if c == '\'' {
			// Character literal (single Unicode character, UTF-8 aware)
			i++
			buf.Reset()
			if i < len(runes) {
				if runes[i] == '\\' && i+1 < len(runes) {
					// Handle escape sequences in character literal
					i++
					switch runes[i] {
					case 'n':
						buf.WriteRune('\n')
					case 't':
						buf.WriteRune('\t')
					case '\\':
						buf.WriteRune('\\')
					case '\'':
						buf.WriteRune('\'')
					case 'u':
						// Unicode escape: \uXXXX (4 hex digits)
						if i+4 < len(runes) {
							hexStr := string(runes[i+1 : i+5])
							var codepoint rune
							fmt.Sscanf(hexStr, "%x", &codepoint)
							buf.WriteRune(codepoint)
							i += 4
						}
					case 'U':
						// Unicode escape: \UXXXXXXXX (8 hex digits)
						if i+8 < len(runes) {
							hexStr := string(runes[i+1 : i+9])
							var codepoint rune
							fmt.Sscanf(hexStr, "%x", &codepoint)
							buf.WriteRune(codepoint)
							i += 8
						}
					default:
						buf.WriteRune(runes[i])
					}
				} else {
					buf.WriteRune(runes[i])
				}
				i++
			}
			if i >= len(runes) || runes[i] != '\'' {
				fmt.Fprintf(os.Stderr, "Unterminated character literal\n")
				return []Token{}
			}
			charStr := buf.String()
			if len([]rune(charStr)) != 1 {
				fmt.Fprintf(os.Stderr, "Character literal must contain exactly one character\n")
				return []Token{}
			}
			tokens = append(tokens, makeToken(TokenChar, charStr))
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
				tokens = append(tokens, makeToken(TokenRet, ""))
			case "return":
				tokens = append(tokens, makeToken(TokenReturn, ""))
			case "const":
				tokens = append(tokens, makeToken(TokenConst, ""))
			case "true", "false":
				tokens = append(tokens, makeToken(TokenBool, word))
			case "int":
				tokens = append(tokens, makeToken(TokenTypeInt, ""))
			case "int8":
				tokens = append(tokens, makeToken(TokenTypeInt8, ""))
			case "int16":
				tokens = append(tokens, makeToken(TokenTypeInt16, ""))
			case "int32":
				tokens = append(tokens, makeToken(TokenTypeInt32, ""))
			case "int64":
				tokens = append(tokens, makeToken(TokenTypeInt64, ""))
			case "uint":
				tokens = append(tokens, makeToken(TokenTypeUint, ""))
			case "uint8":
				tokens = append(tokens, makeToken(TokenTypeUint8, ""))
			case "uint16":
				tokens = append(tokens, makeToken(TokenTypeUint16, ""))
			case "uint32":
				tokens = append(tokens, makeToken(TokenTypeUint32, ""))
			case "uint64":
				tokens = append(tokens, makeToken(TokenTypeUint64, ""))
			case "string":
				tokens = append(tokens, makeToken(TokenTypeString, ""))
			case "char":
				tokens = append(tokens, makeToken(TokenTypeChar, ""))
			case "bool":
				tokens = append(tokens, makeToken(TokenTypeBool, ""))
			case "float":
				tokens = append(tokens, makeToken(TokenTypeFloat, ""))
			case "struct":
				tokens = append(tokens, makeToken(TokenStruct, ""))
			case "enum":
				tokens = append(tokens, makeToken(TokenEnum, ""))
			case "class":
				tokens = append(tokens, makeToken(TokenClass, ""))
			case "new":
				tokens = append(tokens, makeToken(TokenNew, ""))
			case "if":
				tokens = append(tokens, makeToken(TokenIf, ""))
			case "else":
				tokens = append(tokens, makeToken(TokenElse, ""))
			case "while":
				tokens = append(tokens, makeToken(TokenWhile, ""))
			case "for":
				tokens = append(tokens, makeToken(TokenFor, ""))
			case "fn":
				tokens = append(tokens, makeToken(TokenFn, ""))
			case "use":
				tokens = append(tokens, makeToken(TokenUse, ""))
			case "as":
				tokens = append(tokens, makeToken(TokenAs, ""))
			case "try":
				tokens = append(tokens, makeToken(TokenTry, ""))
			case "catch":
				tokens = append(tokens, makeToken(TokenCatch, ""))
			case "finally":
				tokens = append(tokens, makeToken(TokenFinally, ""))
			case "throw":
				tokens = append(tokens, makeToken(TokenThrow, ""))
			case "null":
				tokens = append(tokens, makeToken(TokenNull, ""))
			case "Printf":
				tokens = append(tokens, makeToken(TokenPrintf, ""))
			case "FPrintf":
				tokens = append(tokens, makeToken(TokenFPrintf, ""))
			case "Println":
				tokens = append(tokens, makeToken(TokenPrintln, ""))
			case "SPrint":
				tokens = append(tokens, makeToken(TokenSPrint, ""))
			case "SPrintf":
				tokens = append(tokens, makeToken(TokenSPrintf, ""))
			case "Sprintln":
				tokens = append(tokens, makeToken(TokenSPrintln, ""))
			case "Fatalf":
				tokens = append(tokens, makeToken(TokenFatalf, ""))
			case "Fatalln":
				tokens = append(tokens, makeToken(TokenFatalln, ""))
			case "Logf":
				tokens = append(tokens, makeToken(TokenLogf, ""))
			case "Logln":
				tokens = append(tokens, makeToken(TokenLogln, ""))
			default:
				tokens = append(tokens, makeToken(TokenIdentifier, word))
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
				tokens = append(tokens, makeToken(TokenFloat, numStr))
			} else {
				tokens = append(tokens, makeToken(TokenInt, numStr))
			}
			buf.Reset()
		} else if unicode.IsSpace(c) {
			// Skip other whitespace (newlines handled above)
			col++
		} else if c == ';' {
			tokens = append(tokens, makeToken(TokenSemi, ""))
		} else if c == '=' {
			// Check for ==
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, makeToken(TokenEqual, ""))
				i++
			} else {
				tokens = append(tokens, makeToken(TokenAssign, ""))
			}
		} else if c == '!' {
			// Check for !=
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, makeToken(TokenNotEqual, ""))
				i++
			} else {
				tokens = append(tokens, makeToken(TokenExclaim, ""))
			}
		} else if c == '<' {
			// Check for <= or <<
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, makeToken(TokenLessEq, ""))
				i++
			} else if i+1 < len(runes) && runes[i+1] == '<' {
				tokens = append(tokens, makeToken(TokenLShift, ""))
				i++
			} else {
				tokens = append(tokens, makeToken(TokenLess, ""))
			}
		} else if c == '>' {
			// Check for >= or >>
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, makeToken(TokenGreaterEq, ""))
				i++
			} else if i+1 < len(runes) && runes[i+1] == '>' {
				tokens = append(tokens, makeToken(TokenRShift, ""))
				i++
			} else {
				tokens = append(tokens, makeToken(TokenGreater, ""))
			}
		} else if c == '+' {
			// Check for ++ or +=
			if i+1 < len(runes) && runes[i+1] == '+' {
				tokens = append(tokens, makeToken(TokenPlusPlus, ""))
				i++
			} else if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, makeToken(TokenPlusEq, ""))
				i++
			} else {
				tokens = append(tokens, makeToken(TokenPlus, ""))
			}
		} else if c == '-' {
			// Check for --, ->, or -=
			if i+1 < len(runes) && runes[i+1] == '-' {
				tokens = append(tokens, makeToken(TokenMinusMinus, ""))
				i++
			} else if i+1 < len(runes) && runes[i+1] == '>' {
				tokens = append(tokens, makeToken(TokenArrow, ""))
				i++
			} else if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, makeToken(TokenMinusEq, ""))
				i++
			} else {
				tokens = append(tokens, makeToken(TokenMinus, ""))
			}
		} else if c == '*' {
			// Check for *=
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, makeToken(TokenStarEq, ""))
				i++
			} else {
				tokens = append(tokens, makeToken(TokenStar, ""))
			}
		} else if c == '/' {
			// Check for // comment
			if i+1 < len(runes) && runes[i+1] == '/' {
				// Skip until end of line
				i += 2
				for i < len(runes) && runes[i] != '\n' {
					i++
				}
				// Don't skip the newline itself, let it be processed
				if i < len(runes) {
					i--
				}
			} else if i+1 < len(runes) && runes[i+1] == '=' {
				// Check for /=
				tokens = append(tokens, makeToken(TokenSlashEq, ""))
				i++
			} else {
				tokens = append(tokens, makeToken(TokenSlash, ""))
			}
		} else if c == '%' {
			// Check for %=
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, makeToken(TokenPercentEq, ""))
				i++
			} else {
				tokens = append(tokens, makeToken(TokenPercent, ""))
			}
		} else if c == '&' {
			// Check for &&
			if i+1 < len(runes) && runes[i+1] == '&' {
				tokens = append(tokens, makeToken(TokenAnd, ""))
				i++
			} else {
				tokens = append(tokens, makeToken(TokenAmpersand, ""))
			}
		} else if c == '[' {
			tokens = append(tokens, makeToken(TokenLBracket, ""))
		} else if c == ']' {
			tokens = append(tokens, makeToken(TokenRBracket, ""))
		} else if c == '.' {
			tokens = append(tokens, makeToken(TokenDot, ""))
		} else if c == ':' {
			tokens = append(tokens, makeToken(TokenColon, ""))
		} else if c == '(' {
			tokens = append(tokens, makeToken(TokenLParen, ""))
		} else if c == ')' {
			tokens = append(tokens, makeToken(TokenRParen, ""))
		} else if c == '{' {
			tokens = append(tokens, makeToken(TokenLBrace, ""))
		} else if c == '}' {
			tokens = append(tokens, makeToken(TokenRBrace, ""))
		} else if c == ',' {
			tokens = append(tokens, makeToken(TokenComma, ""))
		} else if c == '|' {
			// Check for ||
			if i+1 < len(runes) && runes[i+1] == '|' {
				tokens = append(tokens, makeToken(TokenOr, ""))
				i++
			} else {
				tokens = append(tokens, makeToken(TokenPipe, ""))
			}
		} else if c == '^' {
			tokens = append(tokens, makeToken(TokenCaret, ""))
		} else if c == '~' {
			tokens = append(tokens, makeToken(TokenTilde, ""))
		} else if c == '?' {
			tokens = append(tokens, makeToken(TokenQuestion, ""))
		} else {
			fmt.Fprintf(os.Stderr, "line %d, col %d: unable to parse '%c'\n", line, col, c)
			return []Token{}
		}

		// Update column for next iteration
		col++
	}

	tokens = append(tokens, makeToken(TokenEOF, ""))
	return tokens
}

// TokenValue returns a string representation of a token's value
func TokenValue(t Token) string {
	switch t.Type {
	case TokenRet:
		return "ret"
	case TokenReturn:
		return "return"
	case TokenInt:
		return t.Value
	case TokenSemi:
		return ";"
	case TokenEOF:
		return "EOF"
	case TokenString:
		return fmt.Sprintf("\"%s\"", t.Value)
	case TokenChar:
		return fmt.Sprintf("'%s'", t.Value)
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
	case TokenTypeInt8:
		return "int8"
	case TokenTypeInt16:
		return "int16"
	case TokenTypeInt32:
		return "int32"
	case TokenTypeInt64:
		return "int64"
	case TokenTypeUint:
		return "uint"
	case TokenTypeUint8:
		return "uint8"
	case TokenTypeUint16:
		return "uint16"
	case TokenTypeUint32:
		return "uint32"
	case TokenTypeUint64:
		return "uint64"
	case TokenTypeString:
		return "string"
	case TokenTypeBool:
		return "bool"
	case TokenTypeFloat:
		return "float"
	case TokenStruct:
		return "struct"
	case TokenEnum:
		return "enum"
	case TokenClass:
		return "class"
	case TokenNew:
		return "new"
	case TokenPrintf:
		return "Printf"
	case TokenFPrintf:
		return "FPrintf"
	case TokenPrintln:
		return "Println"
	case TokenSPrint:
		return "SPrint"
	case TokenSPrintln:
		return "SPrintln"
	case TokenSPrintf:
		return "SPrintf"
	case TokenFatalf:
		return "Fatalf"
	case TokenFatalln:
		return "Fatalln"
	case TokenLogf:
		return "Logf"
	case TokenLogln:
		return "Logln"
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenComma:
		return ","
	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenStar:
		return "*"
	case TokenSlash:
		return "/"
	case TokenPercent:
		return "%"
	case TokenEqual:
		return "=="
	case TokenNotEqual:
		return "!="
	case TokenLess:
		return "<"
	case TokenLessEq:
		return "<="
	case TokenGreater:
		return ">"
	case TokenGreaterEq:
		return ">="
	case TokenAmpersand:
		return "&"
	case TokenExclaim:
		return "!"
	case TokenLBrace:
		return "{"
	case TokenRBrace:
		return "}"
	case TokenLBracket:
		return "["
	case TokenRBracket:
		return "]"
	case TokenDot:
		return "."
	case TokenColon:
		return ":"
	case TokenArrow:
		return "->"
	case TokenPipe:
		return "|"
	case TokenCaret:
		return "^"
	case TokenTilde:
		return "~"
	case TokenAnd:
		return "&&"
	case TokenOr:
		return "||"
	case TokenLShift:
		return "<<"
	case TokenRShift:
		return ">>"
	case TokenPlusPlus:
		return "++"
	case TokenMinusMinus:
		return "--"
	case TokenPlusEq:
		return "+="
	case TokenMinusEq:
		return "-="
	case TokenStarEq:
		return "*="
	case TokenSlashEq:
		return "/="
	case TokenPercentEq:
		return "%="
	case TokenQuestion:
		return "?"
	case TokenIf:
		return "if"
	case TokenElse:
		return "else"
	case TokenWhile:
		return "while"
	case TokenFor:
		return "for"
	case TokenFn:
		return "fn"
	default:
		return "unknown"
	}
}
