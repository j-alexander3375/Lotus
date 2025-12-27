package main

import (
	"fmt"
	"os"
	"strings"
)

// DiagnosticLevel represents the severity of a diagnostic message
type DiagnosticLevel int

const (
	DiagnosticError DiagnosticLevel = iota
	DiagnosticWarning
	DiagnosticInfo
	DiagnosticHint
)

// DiagnosticCategory represents the category of warning/error for filtering
type DiagnosticCategory string

const (
	CategorySyntax     DiagnosticCategory = "syntax"
	CategorySemantic   DiagnosticCategory = "semantic"
	CategoryType       DiagnosticCategory = "type"
	CategoryUnused     DiagnosticCategory = "unused"
	CategoryShadow     DiagnosticCategory = "shadow"
	CategoryImplicit   DiagnosticCategory = "implicit"
	CategoryDeprecated DiagnosticCategory = "deprecated"
	CategoryStyle      DiagnosticCategory = "style"
	CategoryMemory     DiagnosticCategory = "memory"
	CategoryGeneral    DiagnosticCategory = "general"
)

// Diagnostic represents a single compiler diagnostic message
type Diagnostic struct {
	Level      DiagnosticLevel
	Category   DiagnosticCategory
	Code       string // e.g., "E0001", "W0042"
	Message    string
	FilePath   string
	Line       int
	Column     int
	EndLine    int // For multi-line diagnostics
	EndColumn  int
	Context    string   // Source line(s) for context
	Suggestion string   // Suggested fix
	Notes      []string // Additional notes/hints
}

// DiagnosticManager collects and reports compiler diagnostics
type DiagnosticManager struct {
	Diagnostics   []Diagnostic
	ErrorCount    int
	WarnCount     int
	MaxErrors     int
	TreatWarnErr  bool                // -Werror: treat warnings as errors
	SuppressWarns bool                // -w: suppress warnings
	UseColor      bool                // Enable colored output
	SourceLines   map[string][]string // Cache of source file lines
}

// NewDiagnosticManager creates a new diagnostic manager with defaults
func NewDiagnosticManager() *DiagnosticManager {
	return &DiagnosticManager{
		Diagnostics:   make([]Diagnostic, 0),
		MaxErrors:     20,
		TreatWarnErr:  false,
		SuppressWarns: false,
		UseColor:      true,
		SourceLines:   make(map[string][]string),
	}
}

// SetSourceLines caches source lines for a file for better error context
func (dm *DiagnosticManager) SetSourceLines(filePath string, source string) {
	dm.SourceLines[filePath] = strings.Split(source, "\n")
}

// getSourceLine returns a specific line from the cached source
func (dm *DiagnosticManager) getSourceLine(filePath string, line int) string {
	if lines, ok := dm.SourceLines[filePath]; ok {
		if line > 0 && line <= len(lines) {
			return lines[line-1]
		}
	}
	return ""
}

// AddError adds an error diagnostic
func (dm *DiagnosticManager) AddError(message, filePath string, line, column int, context string) {
	dm.AddErrorWithCode("", CategoryGeneral, message, filePath, line, column, context)
}

// AddErrorWithCode adds an error with a specific error code and category
func (dm *DiagnosticManager) AddErrorWithCode(code string, category DiagnosticCategory, message, filePath string, line, column int, context string) {
	if dm.ErrorCount >= dm.MaxErrors {
		return // Stop collecting after max errors
	}

	dm.Diagnostics = append(dm.Diagnostics, Diagnostic{
		Level:    DiagnosticError,
		Category: category,
		Code:     code,
		Message:  message,
		FilePath: filePath,
		Line:     line,
		Column:   column,
		Context:  context,
	})
	dm.ErrorCount++
}

// AddErrorWithSuggestion adds an error with a suggested fix
func (dm *DiagnosticManager) AddErrorWithSuggestion(message, suggestion, filePath string, line, column int) {
	if dm.ErrorCount >= dm.MaxErrors {
		return
	}

	context := dm.getSourceLine(filePath, line)
	dm.Diagnostics = append(dm.Diagnostics, Diagnostic{
		Level:      DiagnosticError,
		Category:   CategoryGeneral,
		Message:    message,
		FilePath:   filePath,
		Line:       line,
		Column:     column,
		Context:    context,
		Suggestion: suggestion,
	})
	dm.ErrorCount++
}

// AddWarning adds a warning diagnostic
func (dm *DiagnosticManager) AddWarning(message, filePath string, line, column int, context string) {
	dm.AddWarningWithCategory(CategoryGeneral, message, filePath, line, column, context)
}

// AddWarningWithCategory adds a warning with a specific category
func (dm *DiagnosticManager) AddWarningWithCategory(category DiagnosticCategory, message, filePath string, line, column int, context string) {
	if dm.SuppressWarns {
		return
	}

	level := DiagnosticWarning
	if dm.TreatWarnErr {
		level = DiagnosticError
		dm.ErrorCount++
	} else {
		dm.WarnCount++
	}

	dm.Diagnostics = append(dm.Diagnostics, Diagnostic{
		Level:    level,
		Category: category,
		Message:  message,
		FilePath: filePath,
		Line:     line,
		Column:   column,
		Context:  context,
	})
}

func (dm *DiagnosticManager) HasErrors() bool {
	return dm.ErrorCount > 0
}

// ReachedMaxErrors returns true if max error limit has been reached
func (dm *DiagnosticManager) ReachedMaxErrors() bool {
	return dm.ErrorCount >= dm.MaxErrors
}

// Print prints all diagnostics with rich formatting
func (dm *DiagnosticManager) Print() {
	for _, diag := range dm.Diagnostics {
		dm.printDiagnostic(diag)
	}

	// Print summary
	if dm.ErrorCount > 0 || dm.WarnCount > 0 {
		fmt.Fprintf(os.Stderr, "\n")

		summaryColor := ""
		resetColor := ""
		if dm.UseColor {
			if dm.ErrorCount > 0 {
				summaryColor = "\033[1;31m" // Red for errors
			} else {
				summaryColor = "\033[1;33m" // Yellow for warnings only
			}
			resetColor = "\033[0m"
		}

		fmt.Fprintf(os.Stderr, "%s", summaryColor)
		if dm.ErrorCount > 0 {
			fmt.Fprintf(os.Stderr, "%d error(s)", dm.ErrorCount)
			if dm.WarnCount > 0 {
				fmt.Fprintf(os.Stderr, " and ")
			}
		}
		if dm.WarnCount > 0 {
			fmt.Fprintf(os.Stderr, "%d warning(s)", dm.WarnCount)
		}
		fmt.Fprintf(os.Stderr, " generated.%s\n", resetColor)

		if dm.ReachedMaxErrors() {
			fmt.Fprintf(os.Stderr, "note: compilation stopped after %d errors\n", dm.MaxErrors)
		}
	}
}

// PrintSummary prints a compact summary without full diagnostics
func (dm *DiagnosticManager) PrintSummary() {
	if dm.ErrorCount > 0 || dm.WarnCount > 0 {
		if dm.ErrorCount > 0 {
			fmt.Fprintf(os.Stderr, "%d error(s)", dm.ErrorCount)
		}
		if dm.WarnCount > 0 {
			if dm.ErrorCount > 0 {
				fmt.Fprintf(os.Stderr, ", ")
			}
			fmt.Fprintf(os.Stderr, "%d warning(s)", dm.WarnCount)
		}
		fmt.Fprintf(os.Stderr, " found\n")
	}
}

// printDiagnostic prints a single diagnostic with rich context
func (dm *DiagnosticManager) printDiagnostic(diag Diagnostic) {
	// Color codes
	var levelStr, colorCode, boldCode, cyanCode, resetColor string
	if dm.UseColor {
		resetColor = "\033[0m"
		boldCode = "\033[1m"
		cyanCode = "\033[36m"
	}

	switch diag.Level {
	case DiagnosticError:
		levelStr = "error"
		if dm.UseColor {
			colorCode = "\033[1;31m" // Bold red
		}
	case DiagnosticWarning:
		levelStr = "warning"
		if dm.UseColor {
			colorCode = "\033[1;33m" // Bold yellow
		}
	case DiagnosticInfo:
		levelStr = "info"
		if dm.UseColor {
			colorCode = "\033[1;36m" // Bold cyan
		}
	case DiagnosticHint:
		levelStr = "hint"
		if dm.UseColor {
			colorCode = "\033[1;32m" // Bold green
		}
	}

	// Add error code if present
	codeStr := ""
	if diag.Code != "" {
		codeStr = fmt.Sprintf("[%s] ", diag.Code)
	}

	// Print location and message
	if diag.FilePath != "" {
		fmt.Fprintf(os.Stderr, "%s%s:%d:%d:%s %s%s%s:%s %s\n",
			boldCode, diag.FilePath, diag.Line, diag.Column, resetColor,
			colorCode, levelStr, resetColor,
			codeStr, diag.Message)
	} else {
		fmt.Fprintf(os.Stderr, "%s%s%s: %s%s\n",
			colorCode, levelStr, resetColor, codeStr, diag.Message)
	}

	// Print source context with line numbers
	if diag.Context != "" {
		lines := strings.Split(diag.Context, "\n")
		lineNumWidth := len(fmt.Sprintf("%d", diag.Line+len(lines)))

		for i, line := range lines {
			if line != "" {
				lineNum := diag.Line + i
				fmt.Fprintf(os.Stderr, " %s%*d |%s %s\n",
					cyanCode, lineNumWidth, lineNum, resetColor, line)
			}
		}

		// Print caret pointing to error column
		if diag.Column > 0 && len(lines) > 0 {
			padding := strings.Repeat(" ", lineNumWidth+3+diag.Column-1)

			// Calculate underline length
			underlineLen := 1
			if diag.EndColumn > diag.Column {
				underlineLen = diag.EndColumn - diag.Column
			}
			underline := strings.Repeat("^", underlineLen)

			fmt.Fprintf(os.Stderr, " %s%s%s%s%s\n", padding, colorCode, underline, resetColor, "")
		}
	}

	// Print suggestion if present
	if diag.Suggestion != "" {
		suggColor := ""
		if dm.UseColor {
			suggColor = "\033[1;32m" // Green for suggestions
		}
		fmt.Fprintf(os.Stderr, "   %ssuggestion:%s %s\n", suggColor, resetColor, diag.Suggestion)
	}

	// Print additional notes
	for _, note := range diag.Notes {
		noteColor := ""
		if dm.UseColor {
			noteColor = "\033[36m" // Cyan for notes
		}
		fmt.Fprintf(os.Stderr, "   %snote:%s %s\n", noteColor, resetColor, note)
	}

	fmt.Fprintf(os.Stderr, "\n")
}

// AddInfo adds an informational diagnostic
func (dm *DiagnosticManager) AddInfo(message, filePath string, line, column int, context string) {
	dm.Diagnostics = append(dm.Diagnostics, Diagnostic{
		Level:    DiagnosticInfo,
		Message:  message,
		FilePath: filePath,
		Line:     line,
		Column:   column,
		Context:  context,
	})
}

// AddHint adds a hint diagnostic
func (dm *DiagnosticManager) AddHint(message, filePath string, line, column int, context string) {
	dm.Diagnostics = append(dm.Diagnostics, Diagnostic{
		Level:    DiagnosticHint,
		Message:  message,
		FilePath: filePath,
		Line:     line,
		Column:   column,
		Context:  context,
	})
}
