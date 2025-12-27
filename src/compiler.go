package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// compiler.go - High-level compilation pipeline orchestration
// This file manages the end-to-end compilation process from source to binary.

// Compiler encapsulates the complete compilation pipeline
type Compiler struct {
	Options *CompilerOptions  // Configuration and command-line options
	Stats   *CompilationStats // Compilation statistics
}

// NewCompiler creates a new compiler instance with the given options
func NewCompiler(opts *CompilerOptions) *Compiler {
	return &Compiler{Options: opts}
}

// CompileFile compiles a single Lotus source file through the full pipeline:
// Source → Tokens → AST → Assembly → Binary
func (c *Compiler) CompileFile(inputPath string) error {
	// Initialize stats tracking
	c.Stats = NewCompilationStats(inputPath)

	if c.Options.Verbose {
		log.Printf("Compiling: input=%s output=%s includes=%v trimpath=%q",
			inputPath, c.Options.OutPath, c.Options.IncludeDirs, c.Options.Trimpath)
	}

	// Phase 1: Read source file
	contents, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}
	c.Stats.SourceBytes = len(contents)
	c.Stats.SourceLines = strings.Count(string(contents), "\n") + 1

	// Phase 2: Lexical analysis (tokenization)
	tokenStart := time.Now()
	tokens := Tokenize(string(contents))
	tokenDuration := time.Since(tokenStart)

	if len(tokens) == 0 {
		return fmt.Errorf("tokenization produced no tokens")
	}
	c.Stats.RecordTokenization(tokenDuration, len(tokens), 0)

	// Handle token dump mode (debugging)
	if c.Options.TokenDump {
		fmt.Println("=== Token Stream ===")
		for i, token := range tokens {
			fmt.Printf("[%d] Type: %v, Value: %s\n", i, token.Type, TokenValue(token))
		}
		return nil
	}

	// Phase 3: Syntax analysis and code generation
	codegenStart := time.Now()
	asm, err := GenerateAssembly(tokens)
	codegenDuration := time.Since(codegenStart)
	if err != nil {
		return err
	}
	asmLines := strings.Count(asm, "\n")
	c.Stats.RecordCodegen(codegenDuration, asmLines, len(asm), 0, 0)

	// Handle assembly output mode (-S flag)
	if c.Options.PrintAsm {
		err := c.writeAssembly(asm)
		c.printStats()
		return err
	}

	// Phase 4: Assemble and link to binary
	if err := c.buildBinary(asm); err != nil {
		return err
	}

	// Phase 5: Optionally run the compiled binary (--run flag)
	if c.Options.RunAfterBuild {
		c.printStats()
		return c.runBinary()
	}

	c.printStats()
	return nil
}

// printStats outputs timing and statistics if enabled
func (c *Compiler) printStats() {
	c.Stats.Finalize()

	// Show timing info
	if c.Options.TimingInfo {
		fmt.Fprintf(os.Stderr, "\n=== Timing ===\n")
		fmt.Fprintf(os.Stderr, "  Tokenize: %v\n", c.Stats.TokenizeTime)
		fmt.Fprintf(os.Stderr, "  Codegen:  %v\n", c.Stats.CodegenTime)
		fmt.Fprintf(os.Stderr, "  Assemble: %v\n", c.Stats.AssembleTime)
		fmt.Fprintf(os.Stderr, "  Link:     %v\n", c.Stats.LinkTime)
		fmt.Fprintf(os.Stderr, "  Total:    %v\n", c.Stats.TotalTime)
	}

	// Show detailed stats
	if c.Options.ShowStats {
		c.Stats.Print()
	}

	// Show file and memory stats (-stat flag)
	if c.Options.ShowFileStat {
		c.printFileStat()
	}
}

// printFileStat outputs file size and memory usage
func (c *Compiler) printFileStat() {
	fmt.Fprintf(os.Stderr, "\n=== Statistics ===\n")

	// Source file info
	fmt.Fprintf(os.Stderr, "Source: %s (%s, %d lines)\n",
		c.Stats.SourceFile, formatBytes(c.Stats.SourceBytes), c.Stats.SourceLines)

	// Token and assembly info
	fmt.Fprintf(os.Stderr, "Tokens: %d\n", c.Stats.TokenCount)
	fmt.Fprintf(os.Stderr, "Assembly: %d lines (%s)\n", c.Stats.AssemblyLines, formatBytes(c.Stats.AssemblyBytes))

	// Output file info
	if c.Stats.OutputFile != "" {
		if info, err := os.Stat(c.Stats.OutputFile); err == nil {
			fmt.Fprintf(os.Stderr, "Output: %s (%s)\n", c.Stats.OutputFile, formatBytes(int(info.Size())))
		}
	}

	// Memory usage
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Fprintf(os.Stderr, "Memory: %s allocated, %s from system\n",
		formatBytes(int(m.Alloc)), formatBytes(int(m.Sys)))
}

// writeAssembly writes assembly output to a file with appropriate extension
func (c *Compiler) writeAssembly(asm string) error {
	asmOut := c.Options.OutPath

	// Auto-generate .s extension if needed
	if asmOut == "a.out" {
		asmOut = "a.s"
	} else if filepath.Ext(asmOut) == "" {
		asmOut = asmOut + ".s"
	}

	if err := os.WriteFile(asmOut, []byte(asm), 0644); err != nil {
		return fmt.Errorf("failed to write assembly file: %w", err)
	}

	if c.Options.Verbose {
		log.Printf("Assembly written to: %s", asmOut)
	}

	return nil
}

// buildBinary assembles and links the assembly to produce an executable binary
func (c *Compiler) buildBinary(asm string) error {
	// Write assembly to temporary file
	tmpAsm := filepath.Join(os.TempDir(), "lotus_tmp.s")
	if err := os.WriteFile(tmpAsm, []byte(asm), 0644); err != nil {
		return fmt.Errorf("failed to write temporary assembly: %w", err)
	}
	defer os.Remove(tmpAsm) // Clean up temp file

	// Invoke GCC to assemble and link
	assembleStart := time.Now()
	cmd := exec.Command("gcc", "-nostartfiles", "-no-pie", "-o", c.Options.OutPath, tmpAsm)

	if c.Options.Verbose {
		log.Printf("Assembling: %s", strings.Join(cmd.Args, " "))
	}

	out, err := cmd.CombinedOutput()
	assembleDuration := time.Since(assembleStart)
	c.Stats.RecordAssemble(assembleDuration)

	if err != nil {
		if len(out) > 0 {
			return fmt.Errorf("assembly failed:\n%s", string(out))
		}
		return fmt.Errorf("assembly failed: %w", err)
	}

	// Record output file info
	if info, statErr := os.Stat(c.Options.OutPath); statErr == nil {
		c.Stats.RecordLink(0, c.Options.OutPath, int(info.Size()))
	} else {
		c.Stats.RecordLink(0, c.Options.OutPath, 0)
	}

	if c.Options.Verbose {
		if len(out) > 0 {
			log.Printf("Assembler output:\n%s", string(out))
		}
		log.Printf("Binary written to: %s", c.Options.OutPath)
	}

	return nil
}

// runBinary executes the compiled binary and streams its output
func (c *Compiler) runBinary() error {
	if c.Options.Verbose {
		log.Printf("Executing: %s", c.Options.OutPath)
	}

	// Execute the binary with inherited stdio
	cmd := exec.Command("./" + c.Options.OutPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Program exited with non-zero status (expected behavior)
			if c.Options.Verbose {
				log.Printf("Program exited with code: %d", exitErr.ExitCode())
			}
			return nil // Don't treat non-zero exit as compiler error
		}
		return fmt.Errorf("failed to execute binary: %w", err)
	}

	return nil
}
