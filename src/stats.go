package main

import (
	"fmt"
	"time"
)

// stats.go - Compilation statistics and metrics tracking
// This file provides functionality to track and report compilation metrics
// including timing, token counts, AST node counts, and output sizes.

// CompilationStats tracks various metrics during compilation
type CompilationStats struct {
	// Timing information
	StartTime    time.Time
	TokenizeTime time.Duration
	ParseTime    time.Duration
	CodegenTime  time.Duration
	AssembleTime time.Duration
	LinkTime     time.Duration
	TotalTime    time.Duration

	// Source code metrics
	SourceFile  string
	SourceLines int
	SourceBytes int

	// Lexical analysis metrics
	TokenCount   int
	CommentCount int

	// Syntax analysis metrics
	ASTNodeCount  int
	FunctionCount int
	VariableCount int
	ConstantCount int

	// Code generation metrics
	AssemblyLines   int
	AssemblyBytes   int
	DataSectionSize int
	TextSectionSize int

	// Import metrics
	ImportCount int
	StdlibCalls int

	// Output metrics
	OutputFile  string
	OutputBytes int
}

// NewCompilationStats creates a new statistics tracker
func NewCompilationStats(sourceFile string) *CompilationStats {
	return &CompilationStats{
		StartTime:  time.Now(),
		SourceFile: sourceFile,
	}
}

// RecordTokenization records lexical analysis metrics
func (cs *CompilationStats) RecordTokenization(duration time.Duration, tokenCount, commentCount int) {
	cs.TokenizeTime = duration
	cs.TokenCount = tokenCount
	cs.CommentCount = commentCount
}

// RecordParsing records syntax analysis metrics
func (cs *CompilationStats) RecordParsing(duration time.Duration, astNodeCount, funcCount, varCount, constCount int) {
	cs.ParseTime = duration
	cs.ASTNodeCount = astNodeCount
	cs.FunctionCount = funcCount
	cs.VariableCount = varCount
	cs.ConstantCount = constCount
}

// RecordCodegen records code generation metrics
func (cs *CompilationStats) RecordCodegen(duration time.Duration, asmLines, asmBytes, dataSize, textSize int) {
	cs.CodegenTime = duration
	cs.AssemblyLines = asmLines
	cs.AssemblyBytes = asmBytes
	cs.DataSectionSize = dataSize
	cs.TextSectionSize = textSize
}

// RecordAssemble records assembly phase metrics
func (cs *CompilationStats) RecordAssemble(duration time.Duration) {
	cs.AssembleTime = duration
}

// RecordLink records linking phase metrics
func (cs *CompilationStats) RecordLink(duration time.Duration, outputFile string, outputBytes int) {
	cs.LinkTime = duration
	cs.OutputFile = outputFile
	cs.OutputBytes = outputBytes
}

// RecordImports records import system metrics
func (cs *CompilationStats) RecordImports(importCount, stdlibCalls int) {
	cs.ImportCount = importCount
	cs.StdlibCalls = stdlibCalls
}

// Finalize calculates total compilation time
func (cs *CompilationStats) Finalize() {
	cs.TotalTime = time.Since(cs.StartTime)
}

// Print outputs a formatted statistics report
func (cs *CompilationStats) Print() {
	fmt.Println("\n=== Compilation Statistics ===")
	fmt.Printf("Source: %s\n", cs.SourceFile)

	// Source metrics
	if cs.SourceLines > 0 {
		fmt.Printf("  Lines: %d\n", cs.SourceLines)
	}
	if cs.SourceBytes > 0 {
		fmt.Printf("  Size: %s\n", formatBytes(cs.SourceBytes))
	}

	// Compilation phases
	fmt.Println("\nPhases:")
	if cs.TokenizeTime > 0 {
		fmt.Printf("  Tokenize: %s (%d tokens, %d comments)\n",
			cs.TokenizeTime, cs.TokenCount, cs.CommentCount)
	}
	if cs.ParseTime > 0 {
		fmt.Printf("  Parse:    %s (%d AST nodes, %d functions, %d vars, %d consts)\n",
			cs.ParseTime, cs.ASTNodeCount, cs.FunctionCount, cs.VariableCount, cs.ConstantCount)
	}
	if cs.CodegenTime > 0 {
		fmt.Printf("  Codegen:  %s (%d lines, %s)\n",
			cs.CodegenTime, cs.AssemblyLines, formatBytes(cs.AssemblyBytes))
	}
	if cs.AssembleTime > 0 {
		fmt.Printf("  Assemble: %s\n", cs.AssembleTime)
	}
	if cs.LinkTime > 0 {
		fmt.Printf("  Link:     %s\n", cs.LinkTime)
	}

	// Import info
	if cs.ImportCount > 0 || cs.StdlibCalls > 0 {
		fmt.Printf("\nImports: %d modules, %d stdlib calls\n", cs.ImportCount, cs.StdlibCalls)
	}

	// Output
	if cs.OutputFile != "" {
		fmt.Printf("\nOutput: %s", cs.OutputFile)
		if cs.OutputBytes > 0 {
			fmt.Printf(" (%s)", formatBytes(cs.OutputBytes))
		}
		fmt.Println()
	}

	// Total
	fmt.Printf("\nTotal Time: %s\n", cs.TotalTime)
	fmt.Println("==============================")
}

// PrintCompact outputs a single-line summary
func (cs *CompilationStats) PrintCompact() {
	fmt.Printf("Compiled %s in %s (%d tokens → %d AST nodes → %d asm lines)\n",
		cs.SourceFile, cs.TotalTime, cs.TokenCount, cs.ASTNodeCount, cs.AssemblyLines)
}

// formatBytes converts bytes to human-readable format
func formatBytes(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
