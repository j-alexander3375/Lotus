package main

import (
	"fmt"
	"os"
)

// main.go - Entry point for the Lotus compiler
// This file handles command-line parsing, version display, and compilation orchestration.

func main() {
	// Phase 1: Parse command-line flags
	opts, args, err := ParseFlags()
	if err != nil {
		// Flag parsing errors already printed to stderr
		os.Exit(2)
	}

	// Phase 2: Handle special flags
	if opts.ShowVersion {
		fmt.Printf("Lotus compiler version %s\n", Version)
		os.Exit(0)
	}

	// Phase 3: Validate input
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Error: no input file specified")
		fmt.Fprintln(os.Stderr, "Usage: lotus [flags] <file>")
		fmt.Fprintln(os.Stderr, "Run 'lotus -h' for help")
		os.Exit(1)
	}

	// Phase 4: Compile the source file
	compiler := NewCompiler(opts)
	if err := compiler.CompileFile(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Compilation failed: %v\n", err)
		os.Exit(1)
	}

	// Success
	os.Exit(0)
}
