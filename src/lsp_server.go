package main

// lsp_server.go - Language Server Protocol implementation for Lotus
// Implements LSP 3.17 specification for IDE integration

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

// ============================================================================
// LSP PROTOCOL TYPES (prefixed to avoid conflicts)
// ============================================================================

// LSPMessage represents an LSP message with JSON-RPC 2.0 structure
type LSPMessage struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Method  string           `json:"method,omitempty"`
	Params  json.RawMessage  `json:"params,omitempty"`
	Result  interface{}      `json:"result,omitempty"`
	Error   *LSPError        `json:"error,omitempty"`
}

// LSPError represents an LSP error response
type LSPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LSPPosition represents a position in a text document
type LSPPosition struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// LSPRange represents a range in a text document
type LSPRange struct {
	Start LSPPosition `json:"start"`
	End   LSPPosition `json:"end"`
}

// LSPLocation represents a location inside a resource
type LSPLocation struct {
	URI   string   `json:"uri"`
	Range LSPRange `json:"range"`
}

// LSPTextDocumentIdentifier identifies a text document
type LSPTextDocumentIdentifier struct {
	URI string `json:"uri"`
}

// LSPTextDocumentPositionParams contains a text document and position
type LSPTextDocumentPositionParams struct {
	TextDocument LSPTextDocumentIdentifier `json:"textDocument"`
	Position     LSPPosition               `json:"position"`
}

// LSPDiagnostic represents a diagnostic (error, warning, hint)
type LSPDiagnostic struct {
	Range    LSPRange `json:"range"`
	Severity int      `json:"severity,omitempty"` // 1=Error, 2=Warning, 3=Info, 4=Hint
	Code     string   `json:"code,omitempty"`
	Source   string   `json:"source,omitempty"`
	Message  string   `json:"message"`
}

// LSPPublishDiagnosticsParams for textDocument/publishDiagnostics
type LSPPublishDiagnosticsParams struct {
	URI         string          `json:"uri"`
	Diagnostics []LSPDiagnostic `json:"diagnostics"`
}

// LSPCompletionItem represents a completion item
type LSPCompletionItem struct {
	Label         string `json:"label"`
	Kind          int    `json:"kind,omitempty"` // CompletionItemKind
	Detail        string `json:"detail,omitempty"`
	Documentation string `json:"documentation,omitempty"`
	InsertText    string `json:"insertText,omitempty"`
}

// CompletionItemKind enum
const (
	LSPCompletionKindText          = 1
	LSPCompletionKindMethod        = 2
	LSPCompletionKindFunction      = 3
	LSPCompletionKindConstructor   = 4
	LSPCompletionKindField         = 5
	LSPCompletionKindVariable      = 6
	LSPCompletionKindClass         = 7
	LSPCompletionKindInterface     = 8
	LSPCompletionKindModule        = 9
	LSPCompletionKindProperty      = 10
	LSPCompletionKindKeyword       = 14
	LSPCompletionKindSnippet       = 15
	LSPCompletionKindTypeParameter = 25
)

// LSPHover represents hover information
type LSPHover struct {
	Contents LSPMarkupContent `json:"contents"`
	Range    *LSPRange        `json:"range,omitempty"`
}

// LSPMarkupContent for hover/completion documentation
type LSPMarkupContent struct {
	Kind  string `json:"kind"` // "plaintext" or "markdown"
	Value string `json:"value"`
}

// ============================================================================
// INITIALIZATION TYPES
// ============================================================================

// LSPInitializeParams for initialize request
type LSPInitializeParams struct {
	ProcessID    int                   `json:"processId"`
	RootURI      string                `json:"rootUri"`
	Capabilities LSPClientCapabilities `json:"capabilities"`
}

// LSPClientCapabilities describes client capabilities
type LSPClientCapabilities struct {
	TextDocument LSPTextDocumentClientCapabilities `json:"textDocument"`
}

// LSPTextDocumentClientCapabilities for text document capabilities
type LSPTextDocumentClientCapabilities struct {
	Completion LSPCompletionClientCapabilities `json:"completion"`
	Hover      LSPHoverClientCapabilities      `json:"hover"`
}

// LSPCompletionClientCapabilities for completion support
type LSPCompletionClientCapabilities struct {
	CompletionItem LSPCompletionItemClientCapabilities `json:"completionItem"`
}

// LSPCompletionItemClientCapabilities for completion item support
type LSPCompletionItemClientCapabilities struct {
	SnippetSupport bool `json:"snippetSupport"`
}

// LSPHoverClientCapabilities for hover support
type LSPHoverClientCapabilities struct {
	ContentFormat []string `json:"contentFormat"`
}

// LSPServerCapabilities describes server capabilities
type LSPServerCapabilities struct {
	TextDocumentSync   int                   `json:"textDocumentSync"`
	CompletionProvider *LSPCompletionOptions `json:"completionProvider,omitempty"`
	HoverProvider      bool                  `json:"hoverProvider"`
	DefinitionProvider bool                  `json:"definitionProvider"`
}

// LSPCompletionOptions for completion settings
type LSPCompletionOptions struct {
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}

// LSPInitializeResult for initialize response
type LSPInitializeResult struct {
	Capabilities LSPServerCapabilities `json:"capabilities"`
	ServerInfo   LSPServerInfo         `json:"serverInfo,omitempty"`
}

// LSPServerInfo about the server
type LSPServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// ============================================================================
// TEXT DOCUMENT SYNC TYPES
// ============================================================================

// LSPDidOpenTextDocumentParams for textDocument/didOpen
type LSPDidOpenTextDocumentParams struct {
	TextDocument LSPTextDocumentItem `json:"textDocument"`
}

// LSPTextDocumentItem describes a text document
type LSPTextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

// LSPDidChangeTextDocumentParams for textDocument/didChange
type LSPDidChangeTextDocumentParams struct {
	TextDocument   LSPVersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []LSPTextDocumentContentChangeEvent `json:"contentChanges"`
}

// LSPVersionedTextDocumentIdentifier includes document version
type LSPVersionedTextDocumentIdentifier struct {
	URI     string `json:"uri"`
	Version int    `json:"version"`
}

// LSPTextDocumentContentChangeEvent describes a content change
type LSPTextDocumentContentChangeEvent struct {
	Text string `json:"text"` // Full content for full sync
}

// LSPDidCloseTextDocumentParams for textDocument/didClose
type LSPDidCloseTextDocumentParams struct {
	TextDocument LSPTextDocumentIdentifier `json:"textDocument"`
}

// ============================================================================
// LSP SERVER
// ============================================================================

// LotusLSPServer is the main LSP server for Lotus
type LotusLSPServer struct {
	documents   map[string]string          // URI -> content
	diagnostics map[string][]LSPDiagnostic // URI -> diagnostics
	symbols     map[string][]LSPSymbol     // URI -> symbols
	mutex       sync.RWMutex
	initialized bool
	shutdown    bool
	logFile     *os.File
}

// LSPSymbol represents a symbol in the document
type LSPSymbol struct {
	Name     string
	Kind     int // SymbolKind
	Location LSPLocation
	Type     string
}

// NewLotusLSPServer creates a new LSP server
func NewLotusLSPServer() *LotusLSPServer {
	logFile, err := os.OpenFile("/tmp/lotus-lsp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logFile = nil
	}

	return &LotusLSPServer{
		documents:   make(map[string]string),
		diagnostics: make(map[string][]LSPDiagnostic),
		symbols:     make(map[string][]LSPSymbol),
		logFile:     logFile,
	}
}

// log writes to the log file
func (s *LotusLSPServer) log(format string, args ...interface{}) {
	if s.logFile != nil {
		fmt.Fprintf(s.logFile, format+"\n", args...)
	}
}

// Run starts the LSP server on stdin/stdout
func (s *LotusLSPServer) Run() error {
	reader := bufio.NewReader(os.Stdin)
	writer := os.Stdout

	s.log("Lotus LSP Server started")

	for {
		msg, err := s.readMessage(reader)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			s.log("Error reading message: %v", err)
			continue
		}

		s.log("Received: %s", msg.Method)

		response := s.handleMessage(msg)
		if response != nil {
			if err := s.writeMessage(writer, response); err != nil {
				s.log("Error writing response: %v", err)
			}
		}

		if s.shutdown {
			return nil
		}
	}
}

// readMessage reads an LSP message from the reader
func (s *LotusLSPServer) readMessage(reader *bufio.Reader) (*LSPMessage, error) {
	var contentLength int
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "Content-Length:") {
			lengthStr := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
			contentLength, _ = strconv.Atoi(lengthStr)
		}
	}

	content := make([]byte, contentLength)
	_, err := io.ReadFull(reader, content)
	if err != nil {
		return nil, err
	}

	var msg LSPMessage
	if err := json.Unmarshal(content, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// writeMessage writes an LSP message to the writer
func (s *LotusLSPServer) writeMessage(writer io.Writer, msg *LSPMessage) error {
	content, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(content))
	if _, err := writer.Write([]byte(header)); err != nil {
		return err
	}
	if _, err := writer.Write(content); err != nil {
		return err
	}

	return nil
}

// handleMessage handles an incoming LSP message
func (s *LotusLSPServer) handleMessage(msg *LSPMessage) *LSPMessage {
	switch msg.Method {
	case "initialize":
		return s.handleInitialize(msg)
	case "initialized":
		s.initialized = true
		return nil
	case "shutdown":
		s.shutdown = true
		return s.createResponse(msg.ID, nil)
	case "exit":
		os.Exit(0)
		return nil
	case "textDocument/didOpen":
		s.handleDidOpen(msg)
		return nil
	case "textDocument/didChange":
		s.handleDidChange(msg)
		return nil
	case "textDocument/didClose":
		s.handleDidClose(msg)
		return nil
	case "textDocument/completion":
		return s.handleCompletion(msg)
	case "textDocument/hover":
		return s.handleHover(msg)
	case "textDocument/definition":
		return s.handleDefinition(msg)
	default:
		s.log("Unhandled method: %s", msg.Method)
		return nil
	}
}

// handleInitialize handles the initialize request
func (s *LotusLSPServer) handleInitialize(msg *LSPMessage) *LSPMessage {
	result := LSPInitializeResult{
		Capabilities: LSPServerCapabilities{
			TextDocumentSync: 1, // Full sync
			CompletionProvider: &LSPCompletionOptions{
				TriggerCharacters: []string{".", ":"},
			},
			HoverProvider:      true,
			DefinitionProvider: true,
		},
		ServerInfo: LSPServerInfo{
			Name:    "lotus-lsp",
			Version: "1.0.0",
		},
	}
	return s.createResponse(msg.ID, result)
}

// handleDidOpen handles textDocument/didOpen
func (s *LotusLSPServer) handleDidOpen(msg *LSPMessage) {
	var params LSPDidOpenTextDocumentParams
	json.Unmarshal(msg.Params, &params)

	s.mutex.Lock()
	s.documents[params.TextDocument.URI] = params.TextDocument.Text
	s.mutex.Unlock()

	s.analyzeDocument(params.TextDocument.URI)
}

// handleDidChange handles textDocument/didChange
func (s *LotusLSPServer) handleDidChange(msg *LSPMessage) {
	var params LSPDidChangeTextDocumentParams
	json.Unmarshal(msg.Params, &params)

	s.mutex.Lock()
	if len(params.ContentChanges) > 0 {
		s.documents[params.TextDocument.URI] = params.ContentChanges[0].Text
	}
	s.mutex.Unlock()

	s.analyzeDocument(params.TextDocument.URI)
}

// handleDidClose handles textDocument/didClose
func (s *LotusLSPServer) handleDidClose(msg *LSPMessage) {
	var params LSPDidCloseTextDocumentParams
	json.Unmarshal(msg.Params, &params)

	s.mutex.Lock()
	delete(s.documents, params.TextDocument.URI)
	delete(s.diagnostics, params.TextDocument.URI)
	delete(s.symbols, params.TextDocument.URI)
	s.mutex.Unlock()
}

// handleCompletion handles textDocument/completion
func (s *LotusLSPServer) handleCompletion(msg *LSPMessage) *LSPMessage {
	var params LSPTextDocumentPositionParams
	json.Unmarshal(msg.Params, &params)

	items := s.getCompletions(params.TextDocument.URI, params.Position)
	return s.createResponse(msg.ID, items)
}

// handleHover handles textDocument/hover
func (s *LotusLSPServer) handleHover(msg *LSPMessage) *LSPMessage {
	var params LSPTextDocumentPositionParams
	json.Unmarshal(msg.Params, &params)

	hover := s.getHover(params.TextDocument.URI, params.Position)
	if hover == nil {
		return s.createResponse(msg.ID, nil)
	}
	return s.createResponse(msg.ID, hover)
}

// handleDefinition handles textDocument/definition
func (s *LotusLSPServer) handleDefinition(msg *LSPMessage) *LSPMessage {
	var params LSPTextDocumentPositionParams
	json.Unmarshal(msg.Params, &params)

	location := s.getDefinition(params.TextDocument.URI, params.Position)
	return s.createResponse(msg.ID, location)
}

// createResponse creates an LSP response message
func (s *LotusLSPServer) createResponse(id *json.RawMessage, result interface{}) *LSPMessage {
	return &LSPMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// ============================================================================
// DOCUMENT ANALYSIS
// ============================================================================

// analyzeDocument analyzes a document and publishes diagnostics
func (s *LotusLSPServer) analyzeDocument(uri string) {
	s.mutex.RLock()
	content, ok := s.documents[uri]
	s.mutex.RUnlock()

	if !ok {
		return
	}

	diagnostics := s.getDiagnostics(content)
	symbols := s.extractSymbols(content, uri)

	s.mutex.Lock()
	s.diagnostics[uri] = diagnostics
	s.symbols[uri] = symbols
	s.mutex.Unlock()

	s.publishDiagnostics(uri, diagnostics)
}

// getDiagnostics tokenizes and parses the content to find errors
func (s *LotusLSPServer) getDiagnostics(content string) []LSPDiagnostic {
	var diagnostics []LSPDiagnostic

	// Tokenize
	tokens := Tokenize(content)

	// Parse
	parser := NewParser(tokens)
	_, err := parser.Parse()
	if err != nil {
		line, col := extractLineCol(err.Error())
		diagnostics = append(diagnostics, LSPDiagnostic{
			Range: LSPRange{
				Start: LSPPosition{Line: line, Character: col},
				End:   LSPPosition{Line: line, Character: col + 10},
			},
			Severity: 1,
			Source:   "lotus",
			Message:  err.Error(),
		})
	}

	return diagnostics
}

// extractLineCol attempts to extract line/col from an error message
func extractLineCol(errMsg string) (int, int) {
	line := 0
	col := 0

	if idx := strings.Index(errMsg, "line "); idx != -1 {
		rest := errMsg[idx+5:]
		if commaIdx := strings.Index(rest, ","); commaIdx != -1 {
			lineStr := rest[:commaIdx]
			if l, err := strconv.Atoi(lineStr); err == nil {
				line = l - 1 // LSP is 0-indexed
			}
			if colIdx := strings.Index(rest, "col "); colIdx != -1 {
				colRest := rest[colIdx+4:]
				if colonIdx := strings.Index(colRest, ":"); colonIdx != -1 {
					colStr := colRest[:colonIdx]
					if c, err := strconv.Atoi(colStr); err == nil {
						col = c - 1
					}
				}
			}
		}
	}

	return line, col
}

// extractSymbols extracts symbols (functions, variables) from content
func (s *LotusLSPServer) extractSymbols(content string, uri string) []LSPSymbol {
	var symbols []LSPSymbol

	tokens := Tokenize(content)
	parser := NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		return symbols
	}

	for _, stmt := range statements {
		switch n := stmt.(type) {
		case *FunctionDefinition:
			symbols = append(symbols, LSPSymbol{
				Name: n.Name,
				Kind: 12, // Function
				Type: TokenTypeName(n.ReturnType),
				Location: LSPLocation{
					URI: uri,
					Range: LSPRange{
						Start: LSPPosition{Line: n.Line - 1, Character: 0},
						End:   LSPPosition{Line: n.Line - 1, Character: len(n.Name)},
					},
				},
			})
		case *VariableDeclaration:
			symbols = append(symbols, LSPSymbol{
				Name: n.Name,
				Kind: 13, // Variable
				Type: TokenTypeName(n.Type),
				Location: LSPLocation{
					URI: uri,
					Range: LSPRange{
						Start: LSPPosition{Line: n.Line - 1, Character: 0},
						End:   LSPPosition{Line: n.Line - 1, Character: len(n.Name)},
					},
				},
			})
		case *ConstantDeclaration:
			symbols = append(symbols, LSPSymbol{
				Name: n.Name,
				Kind: 14, // Constant
				Type: TokenTypeName(n.Type),
				Location: LSPLocation{
					URI: uri,
					Range: LSPRange{
						Start: LSPPosition{Line: n.Line - 1, Character: 0},
						End:   LSPPosition{Line: n.Line - 1, Character: len(n.Name)},
					},
				},
			})
		case *StructDefinition:
			symbols = append(symbols, LSPSymbol{
				Name: n.Name,
				Kind: 23, // Struct
				Location: LSPLocation{
					URI: uri,
					Range: LSPRange{
						Start: LSPPosition{Line: n.Line - 1, Character: 0},
						End:   LSPPosition{Line: n.Line - 1, Character: len(n.Name)},
					},
				},
			})
		case *EnumDefinition:
			symbols = append(symbols, LSPSymbol{
				Name: n.Name,
				Kind: 10, // Enum
				Location: LSPLocation{
					URI: uri,
					Range: LSPRange{
						Start: LSPPosition{Line: n.Line - 1, Character: 0},
						End:   LSPPosition{Line: n.Line - 1, Character: len(n.Name)},
					},
				},
			})
		}
	}

	return symbols
}

// publishDiagnostics sends diagnostics to the client
func (s *LotusLSPServer) publishDiagnostics(uri string, diagnostics []LSPDiagnostic) {
	params := LSPPublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	}

	msg := &LSPMessage{
		JSONRPC: "2.0",
		Method:  "textDocument/publishDiagnostics",
	}
	msg.Params, _ = json.Marshal(params)

	s.writeMessage(os.Stdout, msg)
}

// ============================================================================
// COMPLETIONS
// ============================================================================

// getCompletions returns completion items for the given position
func (s *LotusLSPServer) getCompletions(uri string, pos LSPPosition) []LSPCompletionItem {
	var items []LSPCompletionItem

	// Add Lotus keywords
	keywords := []string{
		"fn", "ret", "return", "if", "else", "while", "for", "break", "continue",
		"const", "struct", "enum", "class", "match", "case", "default",
		"use", "as", "static", "lcl", "gbl", "virtual", "override",
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float", "bool", "string", "char", "void",
		"true", "false", "null",
		"bitcast", "transmute", // Type reinterpretation
		"try", "catch", "finally", "throw", // Error handling
		"partial", "wrap", // Functional programming
	}

	for _, kw := range keywords {
		items = append(items, LSPCompletionItem{
			Label:  kw,
			Kind:   LSPCompletionKindKeyword,
			Detail: "keyword",
		})
	}

	// Add built-in functions
	builtins := map[string]string{
		"printf":  "void printf(string format, ...)",
		"println": "void println(string s)",
		"sprintf": "string sprintf(string format, ...)",
		"malloc":  "ptr malloc(int size)",
		"free":    "void free(ptr p)",
		"sizeof":  "int sizeof(type)",
		"len":     "int len(array a)",
		"append":  "array append(array a, any value)",
		"rand":    "int rand()",
		"seed":    "void seed(int n)",
		"exit":    "void exit(int code)",
	}

	for name, sig := range builtins {
		items = append(items, LSPCompletionItem{
			Label:  name,
			Kind:   LSPCompletionKindFunction,
			Detail: sig,
		})
	}

	// Add symbols from document
	s.mutex.RLock()
	if symbols, ok := s.symbols[uri]; ok {
		for _, sym := range symbols {
			kind := LSPCompletionKindVariable
			if sym.Kind == 12 {
				kind = LSPCompletionKindFunction
			} else if sym.Kind == 23 {
				kind = LSPCompletionKindClass
			}
			items = append(items, LSPCompletionItem{
				Label:  sym.Name,
				Kind:   kind,
				Detail: sym.Type,
			})
		}
	}
	s.mutex.RUnlock()

	return items
}

// ============================================================================
// HOVER
// ============================================================================

// getHover returns hover information for the given position
func (s *LotusLSPServer) getHover(uri string, pos LSPPosition) *LSPHover {
	s.mutex.RLock()
	content, ok := s.documents[uri]
	symbols := s.symbols[uri]
	s.mutex.RUnlock()

	if !ok {
		return nil
	}

	word := getWordAtPosition(content, pos)
	if word == "" {
		return nil
	}

	// Check if it's a symbol
	for _, sym := range symbols {
		if sym.Name == word {
			return &LSPHover{
				Contents: LSPMarkupContent{
					Kind:  "markdown",
					Value: fmt.Sprintf("**%s**\n\n```lotus\n%s %s\n```", sym.Name, sym.Type, sym.Name),
				},
			}
		}
	}

	// Check if it's a keyword
	keywordDocs := map[string]string{
		"fn":     "Define a function\n\n```lotus\nfn returnType name(params) { body }\n```",
		"ret":    "Return from function\n\n```lotus\nret value;\n```",
		"if":     "Conditional statement\n\n```lotus\nif condition { body } else { body }\n```",
		"while":  "While loop\n\n```lotus\nwhile condition { body }\n```",
		"for":    "For loop\n\n```lotus\nfor init; condition; update { body }\n```",
		"match":  "Pattern matching\n\n```lotus\nmatch value { case pattern => body }\n```",
		"struct": "Define a struct type\n\n```lotus\nstruct Name { field1 type1, field2 type2 }\n```",
		"enum":   "Define an enum type\n\n```lotus\nenum Name { Value1, Value2 }\n```",
		"const":  "Declare a constant\n\n```lotus\nconst type name = value;\n```",
	}

	if doc, ok := keywordDocs[word]; ok {
		return &LSPHover{
			Contents: LSPMarkupContent{
				Kind:  "markdown",
				Value: doc,
			},
		}
	}

	return nil
}

// getWordAtPosition gets the word at the given position
func getWordAtPosition(content string, pos LSPPosition) string {
	lines := strings.Split(content, "\n")
	if pos.Line >= len(lines) {
		return ""
	}

	line := lines[pos.Line]
	if pos.Character >= len(line) {
		return ""
	}

	start := pos.Character
	end := pos.Character

	for start > 0 && isWordChar(line[start-1]) {
		start--
	}
	for end < len(line) && isWordChar(line[end]) {
		end++
	}

	if start >= end {
		return ""
	}

	return line[start:end]
}

// isWordChar returns true if c is a word character
func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

// ============================================================================
// GO TO DEFINITION
// ============================================================================

// getDefinition returns the definition location for the given position
func (s *LotusLSPServer) getDefinition(uri string, pos LSPPosition) *LSPLocation {
	s.mutex.RLock()
	content, ok := s.documents[uri]
	symbols := s.symbols[uri]
	s.mutex.RUnlock()

	if !ok {
		return nil
	}

	word := getWordAtPosition(content, pos)
	if word == "" {
		return nil
	}

	for _, sym := range symbols {
		if sym.Name == word {
			return &sym.Location
		}
	}

	return nil
}

// ============================================================================
// MAIN ENTRY POINT
// ============================================================================

// RunLSPServer starts the LSP server
func RunLSPServer() {
	server := NewLotusLSPServer()
	if err := server.Run(); err != nil {
		log.Fatalf("LSP server error: %v", err)
	}
}
