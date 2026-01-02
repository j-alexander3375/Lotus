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
// Or with LLVM: Source → Tokens → AST → LLVM IR → Binary
func (c *Compiler) CompileFile(inputPath string) error {
	// Initialize stats tracking
	c.Stats = NewCompilationStats(inputPath)

	if c.Options.Verbose {
		log.Printf("Compiling: input=%s output=%s includes=%v trimpath=%q gcc=%v",
			inputPath, c.Options.OutPath, c.Options.IncludeDirs, c.Options.Trimpath, c.Options.UseGCC)
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

	// Check if GCC backend is requested (LLVM is now the default)
	if c.Options.UseGCC {
		// Legacy GCC/assembly backend
		return c.compileWithGCC(inputPath, tokens)
	}

	// Default: Use LLVM backend
	return c.compileWithLLVM(inputPath, tokens)
}

// compileWithGCC compiles using the legacy GCC/assembly backend
func (c *Compiler) compileWithGCC(inputPath string, tokens []Token) error {
	// Phase 3: Syntax analysis and code generation (direct assembly)
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

// compileWithLLVM compiles using the LLVM backend
func (c *Compiler) compileWithLLVM(inputPath string, tokens []Token) error {
	if c.Options.Verbose {
		log.Printf("Using LLVM backend with optimization level %d", c.Options.OptLevel)
	}

	// Parse tokens to AST
	codegenStart := time.Now()
	parser := NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	// Apply AST-level optimizations
	statements = OptimizeAST(statements)
	statements = EliminateDeadCode(statements)
	statements = EliminateUnusedFunctions(statements)

	// Create LLVM code generator
	moduleName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	llvmGen := NewLLVMCodeGenerator(moduleName)
	defer llvmGen.Dispose()

	// Generate LLVM IR
	if err := llvmGen.Generate(statements); err != nil {
		return fmt.Errorf("LLVM codegen error: %w", err)
	}

	// Apply optimizations
	if c.Options.OptLevel > 0 {
		llvmGen.Optimize(c.Options.OptLevel)
	}

	codegenDuration := time.Since(codegenStart)
	ir := llvmGen.GetIR()
	c.Stats.RecordCodegen(codegenDuration, strings.Count(ir, "\n"), len(ir), 0, 0)

	// Handle LLVM IR output mode (--emit-llvm flag)
	if c.Options.EmitLLVMIR {
		return c.writeLLVMIR(ir)
	}

	// Compile to object file and link
	if err := c.buildBinaryWithLLVM(llvmGen); err != nil {
		return err
	}

	// Optionally run the compiled binary
	if c.Options.RunAfterBuild {
		c.printStats()
		return c.runBinary()
	}

	c.printStats()
	return nil
}

// writeLLVMIR writes LLVM IR to a file
func (c *Compiler) writeLLVMIR(ir string) error {
	irOut := c.Options.OutPath

	// Auto-generate .ll extension if needed
	if irOut == "a.out" {
		irOut = "a.ll"
	} else if filepath.Ext(irOut) == "" {
		irOut = irOut + ".ll"
	}

	if err := os.WriteFile(irOut, []byte(ir), 0644); err != nil {
		return fmt.Errorf("failed to write LLVM IR file: %w", err)
	}

	if c.Options.Verbose {
		log.Printf("LLVM IR written to: %s", irOut)
	}

	return nil
}

// buildBinaryWithLLVM compiles LLVM IR to a binary
func (c *Compiler) buildBinaryWithLLVM(llvmGen *LLVMCodeGenerator) error {
	// Determine target triple
	targetTriple := c.Options.TargetTriple
	if targetTriple == "" {
		// Use default host triple
		targetTriple = getDefaultTargetTriple()
	}

	// Write LLVM IR to temp file
	tmpIR := filepath.Join(os.TempDir(), "lotus_tmp.ll")
	if err := os.WriteFile(tmpIR, []byte(llvmGen.GetIR()), 0644); err != nil {
		return fmt.Errorf("failed to write temporary LLVM IR: %w", err)
	}
	defer os.Remove(tmpIR)

	// Use clang to compile LLVM IR to binary
	assembleStart := time.Now()
	optFlag := fmt.Sprintf("-O%d", c.Options.OptLevel)
	cmd := exec.Command("clang", optFlag, "-o", c.Options.OutPath, tmpIR)

	if c.Options.TargetTriple != "" {
		cmd.Args = append(cmd.Args, "-target", c.Options.TargetTriple)
	}

	if c.Options.Verbose {
		log.Printf("Compiling with clang: %s", strings.Join(cmd.Args, " "))
	}

	out, err := cmd.CombinedOutput()
	assembleDuration := time.Since(assembleStart)
	c.Stats.RecordAssemble(assembleDuration)

	if err != nil {
		if len(out) > 0 {
			return fmt.Errorf("LLVM compilation failed:\n%s", string(out))
		}
		return fmt.Errorf("LLVM compilation failed: %w", err)
	}

	// Record output file info
	if info, statErr := os.Stat(c.Options.OutPath); statErr == nil {
		c.Stats.RecordLink(0, c.Options.OutPath, int(info.Size()))
	}

	if c.Options.Verbose {
		log.Printf("Binary written to: %s", c.Options.OutPath)
	}

	return nil
}

// getDefaultTargetTriple returns the default target triple for the host system
func getDefaultTargetTriple() string {
	// Map Go's GOOS/GOARCH to LLVM target triples
	arch := runtime.GOARCH
	osName := runtime.GOOS

	var llvmArch string
	switch arch {
	case "amd64":
		llvmArch = "x86_64"
	case "arm64":
		llvmArch = "aarch64"
	case "arm":
		llvmArch = "arm"
	case "386":
		llvmArch = "i386"
	default:
		llvmArch = arch
	}

	var llvmOS string
	switch osName {
	case "linux":
		llvmOS = "linux-gnu"
	case "darwin":
		llvmOS = "apple-darwin"
	case "windows":
		llvmOS = "windows-msvc"
	default:
		llvmOS = osName
	}

	return fmt.Sprintf("%s-unknown-%s", llvmArch, llvmOS)
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
