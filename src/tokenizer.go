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
			case "return":
				tokens = append(tokens, Token{Type: TokenReturn, Value: ""})
			case "const":
				tokens = append(tokens, Token{Type: TokenConst, Value: ""})
			case "true", "false":
				tokens = append(tokens, Token{Type: TokenBool, Value: word})
			case "int":
				tokens = append(tokens, Token{Type: TokenTypeInt, Value: ""})
			case "int8":
				tokens = append(tokens, Token{Type: TokenTypeInt8, Value: ""})
			case "int16":
				tokens = append(tokens, Token{Type: TokenTypeInt16, Value: ""})
			case "int32":
				tokens = append(tokens, Token{Type: TokenTypeInt32, Value: ""})
			case "int64":
				tokens = append(tokens, Token{Type: TokenTypeInt64, Value: ""})
			case "uint":
				tokens = append(tokens, Token{Type: TokenTypeUint, Value: ""})
			case "uint8":
				tokens = append(tokens, Token{Type: TokenTypeUint8, Value: ""})
			case "uint16":
				tokens = append(tokens, Token{Type: TokenTypeUint16, Value: ""})
			case "uint32":
				tokens = append(tokens, Token{Type: TokenTypeUint32, Value: ""})
			case "uint64":
				tokens = append(tokens, Token{Type: TokenTypeUint64, Value: ""})
			case "string":
				tokens = append(tokens, Token{Type: TokenTypeString, Value: ""})
			case "bool":
				tokens = append(tokens, Token{Type: TokenTypeBool, Value: ""})
			case "float":
				tokens = append(tokens, Token{Type: TokenTypeFloat, Value: ""})
			case "struct":
				tokens = append(tokens, Token{Type: TokenStruct, Value: ""})
			case "enum":
				tokens = append(tokens, Token{Type: TokenEnum, Value: ""})
			case "class":
				tokens = append(tokens, Token{Type: TokenClass, Value: ""})
			case "new":
				tokens = append(tokens, Token{Type: TokenNew, Value: ""})
			case "malloc":
				tokens = append(tokens, Token{Type: TokenMalloc, Value: ""})
			case "free":
				tokens = append(tokens, Token{Type: TokenFree, Value: ""})
			case "sizeof":
				tokens = append(tokens, Token{Type: TokenSizeof, Value: ""})
			case "if":
				tokens = append(tokens, Token{Type: TokenIf, Value: ""})
			case "else":
				tokens = append(tokens, Token{Type: TokenElse, Value: ""})
			case "while":
				tokens = append(tokens, Token{Type: TokenWhile, Value: ""})
			case "for":
				tokens = append(tokens, Token{Type: TokenFor, Value: ""})
			case "fn":
				tokens = append(tokens, Token{Type: TokenFn, Value: ""})
			case "use":
				tokens = append(tokens, Token{Type: TokenUse, Value: ""})
			case "as":
				tokens = append(tokens, Token{Type: TokenAs, Value: ""})
			case "try":
				tokens = append(tokens, Token{Type: TokenTry, Value: ""})
			case "catch":
				tokens = append(tokens, Token{Type: TokenCatch, Value: ""})
			case "finally":
				tokens = append(tokens, Token{Type: TokenFinally, Value: ""})
			case "throw":
				tokens = append(tokens, Token{Type: TokenThrow, Value: ""})
			case "null":
				tokens = append(tokens, Token{Type: TokenNull, Value: ""})
			case "Printf":
				tokens = append(tokens, Token{Type: TokenPrintf, Value: ""})
			case "FPrintf":
				tokens = append(tokens, Token{Type: TokenFPrintf, Value: ""})
			case "Println":
				tokens = append(tokens, Token{Type: TokenPrintln, Value: ""})
			case "SPrint":
				tokens = append(tokens, Token{Type: TokenSPrint, Value: ""})
			case "SPrintf":
				tokens = append(tokens, Token{Type: TokenSPrintf, Value: ""})
			case "Sprintln":
				tokens = append(tokens, Token{Type: TokenSPrintln, Value: ""})
			case "Fatalf":
				tokens = append(tokens, Token{Type: TokenFatalf, Value: ""})
			case "Fatalln":
				tokens = append(tokens, Token{Type: TokenFatalln, Value: ""})
			case "Logf":
				tokens = append(tokens, Token{Type: TokenLogf, Value: ""})
			case "Logln":
				tokens = append(tokens, Token{Type: TokenLogln, Value: ""})
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
			// Check for ==
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: TokenEqual, Value: ""})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenAssign, Value: ""})
			}
		} else if c == '!' {
			// Check for !=
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: TokenNotEqual, Value: ""})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenExclaim, Value: ""})
			}
		} else if c == '<' {
			// Check for <= or <<
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: TokenLessEq, Value: ""})
				i++
			} else if i+1 < len(runes) && runes[i+1] == '<' {
				tokens = append(tokens, Token{Type: TokenLShift, Value: ""})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenLess, Value: ""})
			}
		} else if c == '>' {
			// Check for >= or >>
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: TokenGreaterEq, Value: ""})
				i++
			} else if i+1 < len(runes) && runes[i+1] == '>' {
				tokens = append(tokens, Token{Type: TokenRShift, Value: ""})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenGreater, Value: ""})
			}
		} else if c == '+' {
			// Check for ++ or +=
			if i+1 < len(runes) && runes[i+1] == '+' {
				tokens = append(tokens, Token{Type: TokenPlusPlus, Value: ""})
				i++
			} else if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: TokenPlusEq, Value: ""})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenPlus, Value: ""})
			}
		} else if c == '-' {
			// Check for --, ->, or -=
			if i+1 < len(runes) && runes[i+1] == '-' {
				tokens = append(tokens, Token{Type: TokenMinusMinus, Value: ""})
				i++
			} else if i+1 < len(runes) && runes[i+1] == '>' {
				tokens = append(tokens, Token{Type: TokenArrow, Value: ""})
				i++
			} else if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: TokenMinusEq, Value: ""})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenMinus, Value: ""})
			}
		} else if c == '*' {
			// Check for *=
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: TokenStarEq, Value: ""})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenStar, Value: ""})
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
				tokens = append(tokens, Token{Type: TokenSlashEq, Value: ""})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenSlash, Value: ""})
			}
		} else if c == '%' {
			// Check for %=
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: TokenPercentEq, Value: ""})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenPercent, Value: ""})
			}
		} else if c == '&' {
			// Check for &&
			if i+1 < len(runes) && runes[i+1] == '&' {
				tokens = append(tokens, Token{Type: TokenAnd, Value: ""})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenAmpersand, Value: ""})
			}
		} else if c == '[' {
			tokens = append(tokens, Token{Type: TokenLBracket, Value: ""})
		} else if c == ']' {
			tokens = append(tokens, Token{Type: TokenRBracket, Value: ""})
		} else if c == '.' {
			tokens = append(tokens, Token{Type: TokenDot, Value: ""})
		} else if c == ':' {
			tokens = append(tokens, Token{Type: TokenColon, Value: ""})
		} else if c == '(' {
			tokens = append(tokens, Token{Type: TokenLParen, Value: ""})
		} else if c == ')' {
			tokens = append(tokens, Token{Type: TokenRParen, Value: ""})
		} else if c == '{' {
			tokens = append(tokens, Token{Type: TokenLBrace, Value: ""})
		} else if c == '}' {
			tokens = append(tokens, Token{Type: TokenRBrace, Value: ""})
		} else if c == ',' {
			tokens = append(tokens, Token{Type: TokenComma, Value: ""})
		} else if c == '|' {
			// Check for ||
			if i+1 < len(runes) && runes[i+1] == '|' {
				tokens = append(tokens, Token{Type: TokenOr, Value: ""})
				i++
			} else {
				tokens = append(tokens, Token{Type: TokenPipe, Value: ""})
			}
		} else if c == '^' {
			tokens = append(tokens, Token{Type: TokenCaret, Value: ""})
		} else if c == '~' {
			tokens = append(tokens, Token{Type: TokenTilde, Value: ""})
		} else if c == '?' {
			tokens = append(tokens, Token{Type: TokenQuestion, Value: ""})
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
	case TokenMalloc:
		return "malloc"
	case TokenFree:
		return "free"
	case TokenSizeof:
		return "sizeof"
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
