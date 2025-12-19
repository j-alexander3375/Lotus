package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const Version = "b0.1.0 pre-release"

// Compiler encapsulates the compilation process
type Compiler struct {
	Options *CompilerOptions
}

// NewCompiler creates a new compiler instance
func NewCompiler(opts *CompilerOptions) *Compiler {
	return &Compiler{Options: opts}
}

// CompileFile compiles a single source file
func (c *Compiler) CompileFile(inputPath string) error {
	if c.Options.Verbose {
		log.Printf("input=%s out=%s includes=%v trimpath=%q",
			inputPath, c.Options.OutPath, c.Options.IncludeDirs, c.Options.Trimpath)
	}

	// Read source file
	contents, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	// Tokenize
	tokens := Tokenize(string(contents))
	if len(tokens) == 0 {
		return fmt.Errorf("tokenization failed")
	}

	// Handle token dump mode
	if c.Options.TokenDump {
		fmt.Println(tokens)
		for _, token := range tokens {
			fmt.Printf("Token Type: %v, Value: %s\n", token.Type, TokenValue(token))
		}
		return nil
	}

	// Generate assembly
	asm := GenerateAssembly(tokens)

	// Handle assembly output mode
	if c.Options.PrintAsm {
		return c.writeAssembly(asm)
	}

	// Handle binary output mode
	if c.Options.OutPath != "a.out" {
		return c.buildBinary(asm)
	}

	// Default: just report tokenization status
	if tokens[len(tokens)-1].Type == TokenEOF {
		log.Println("Tokenization completed successfully.")
		log.Print("End of File reached. Closing File Buffer.")
	} else {
		return fmt.Errorf("tokenization did not complete successfully")
	}

	return nil
}

// writeAssembly writes assembly output to file
func (c *Compiler) writeAssembly(asm string) error {
	asmOut := c.Options.OutPath
	if asmOut == "a.out" {
		asmOut = "a.s"
	} else if filepath.Ext(asmOut) == "" {
		asmOut = asmOut + ".s"
	}

	if err := os.WriteFile(asmOut, []byte(asm), 0644); err != nil {
		return fmt.Errorf("error writing assembly: %w", err)
	}

	if c.Options.Verbose {
		log.Printf("Wrote assembly to %s", asmOut)
	}

	return nil
}

// buildBinary assembles and links to produce an executable
func (c *Compiler) buildBinary(asm string) error {
	// Write temp assembly
	tmpAsm := filepath.Join(os.TempDir(), "lotus_tmp.s")
	if err := os.WriteFile(tmpAsm, []byte(asm), 0644); err != nil {
		return fmt.Errorf("error writing temp assembly: %w", err)
	}

	// Assemble and link
	cmd := exec.Command("gcc", "-nostartfiles", "-no-pie", "-o", c.Options.OutPath, tmpAsm)
	if c.Options.Verbose {
		log.Printf("Running: %s", strings.Join(cmd.Args, " "))
	}

	out, err := cmd.CombinedOutput()
	if c.Options.Verbose && len(out) > 0 {
		log.Printf("gcc output:\n%s", string(out))
	}

	if err != nil {
		return fmt.Errorf("failed to build binary: %w", err)
	}

	if c.Options.Verbose {
		log.Printf("Wrote binary to %s", c.Options.OutPath)
	}

	return nil
}
