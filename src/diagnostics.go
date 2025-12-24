package main

import (
	"fmt"
	"os"
	"strings"
)

type DiagnosticLevel int

const (
	DiagnosticError DiagnosticLevel = iota
	DiagnosticWarning
	DiagnosticInfo
	DiagnosticHint
)

type Diagnostic struct {
	Level    DiagnosticLevel
	Message  string
	FilePath string
	Line     int
	Column   int
	Context  string
}

type DiagnosticManager struct {
	Diagnostics []Diagnostic
	ErrorCount  int
	WarnCount   int
}

func NewDiagnosticManager() *DiagnosticManager {
	return &DiagnosticManager{
		Diagnostics: make([]Diagnostic, 0),
	}
}

func (dm *DiagnosticManager) AddError(message, filePath string, line, column int, context string) {
	dm.Diagnostics = append(dm.Diagnostics, Diagnostic{
		Level:    DiagnosticError,
		Message:  message,
		FilePath: filePath,
		Line:     line,
		Column:   column,
		Context:  context,
	})
	dm.ErrorCount++
}

func (dm *DiagnosticManager) AddWarning(message, filePath string, line, column int, context string) {
	dm.Diagnostics = append(dm.Diagnostics, Diagnostic{
		Level:    DiagnosticWarning,
		Message:  message,
		FilePath: filePath,
		Line:     line,
		Column:   column,
		Context:  context,
	})
	dm.WarnCount++
}

func (dm *DiagnosticManager) HasErrors() bool {
	return dm.ErrorCount > 0
}

func (dm *DiagnosticManager) Print() {
	for _, diag := range dm.Diagnostics {
		dm.printDiagnostic(diag)
	}
	if dm.ErrorCount > 0 || dm.WarnCount > 0 {
		fmt.Fprintf(os.Stderr, "\n")
		if dm.ErrorCount > 0 {
			fmt.Fprintf(os.Stderr, "%d error(s) generated.\n", dm.ErrorCount)
		}
		if dm.WarnCount > 0 {
			fmt.Fprintf(os.Stderr, "%d warning(s) generated.\n", dm.WarnCount)
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

func (dm *DiagnosticManager) printDiagnostic(diag Diagnostic) {
	levelStr := ""
	colorCode := ""
	resetColor := "\033[0m"

	switch diag.Level {
	case DiagnosticError:
		levelStr = "error"
		colorCode = "\033[1;31m"
	case DiagnosticWarning:
		levelStr = "warning"
		colorCode = "\033[1;33m"
	case DiagnosticInfo:
		levelStr = "info"
		colorCode = "\033[1;36m"
	case DiagnosticHint:
		levelStr = "hint"
		colorCode = "\033[1;32m"
	}

	if diag.FilePath != "" {
		fmt.Fprintf(os.Stderr, "%s%s:%d:%d: %s: %s%s\n",
			colorCode, diag.FilePath, diag.Line, diag.Column, levelStr, diag.Message, resetColor)
	} else {
		fmt.Fprintf(os.Stderr, "%s%s: %s%s\n",
			colorCode, levelStr, diag.Message, resetColor)
	}

	if diag.Context != "" {
		lines := strings.Split(diag.Context, "\n")
		for i, line := range lines {
			if line != "" {
				lineNum := diag.Line + i
				fmt.Fprintf(os.Stderr, "%5d | %s\n", lineNum, line)
			}
		}
		if diag.Column > 0 {
			spaces := strings.Repeat(" ", 8+diag.Column-1)
			fmt.Fprintf(os.Stderr, "%s%s^%s\n", spaces, colorCode, resetColor)
		}
	}
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
