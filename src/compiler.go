package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// compiler.go - High-level compilation pipeline orchestration
// This file manages the end-to-end compilation process from source to binary.

// Compiler encapsulates the complete compilation pipeline
type Compiler struct {
	Options *CompilerOptions // Configuration and command-line options
}

// NewCompiler creates a new compiler instance with the given options
func NewCompiler(opts *CompilerOptions) *Compiler {
	return &Compiler{Options: opts}
}

// CompileFile compiles a single Lotus source file through the full pipeline:
// Source → Tokens → AST → Assembly → Binary
func (c *Compiler) CompileFile(inputPath string) error {
	if c.Options.Verbose {
		log.Printf("Compiling: input=%s output=%s includes=%v trimpath=%q",
			inputPath, c.Options.OutPath, c.Options.IncludeDirs, c.Options.Trimpath)
	}

	// Phase 1: Read source file
	contents, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Phase 2: Lexical analysis (tokenization)
	tokens := Tokenize(string(contents))
	if len(tokens) == 0 {
		return fmt.Errorf("tokenization produced no tokens")
	}

	// Handle token dump mode (debugging)
	if c.Options.TokenDump {
		fmt.Println("=== Token Stream ===")
		for i, token := range tokens {
			fmt.Printf("[%d] Type: %v, Value: %s\n", i, token.Type, TokenValue(token))
		}
		return nil
	}

	// Phase 3: Syntax analysis and code generation
	asm, err := GenerateAssembly(tokens)
	if err != nil {
		return err
	}

	// Handle assembly output mode (-S flag)
	if c.Options.PrintAsm {
		return c.writeAssembly(asm)
	}

	// Phase 4: Assemble and link to binary
	if err := c.buildBinary(asm); err != nil {
		return err
	}

	// Phase 5: Optionally run the compiled binary (--run flag)
	if c.Options.RunAfterBuild {
		return c.runBinary()
	}

	return nil
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
	cmd := exec.Command("gcc", "-nostartfiles", "-no-pie", "-o", c.Options.OutPath, tmpAsm)

	if c.Options.Verbose {
		log.Printf("Assembling: %s", strings.Join(cmd.Args, " "))
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return fmt.Errorf("assembly failed:\n%s", string(out))
		}
		return fmt.Errorf("assembly failed: %w", err)
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
