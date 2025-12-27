package main

import (
	"fmt"
	"io"
	"os"
)

// main.go - Entry point for the Lotus compiler
// This file handles command-line parsing, version display, and compilation orchestration.

func main() {
	os.Exit(run())
}

// run orchestrates CLI parsing and compilation, returning a process exit code.
func run() int {
	// Phase 1: Parse command-line flags
	opts, args, err := ParseFlags()
	if err != nil {
		// Flag parsing errors already printed to stderr
		return 2
	}

	// Phase 2: Handle special flags
	if opts.ShowVersion {
		fmt.Printf("Lotus compiler version %s\n", Version)
		return 0
	}

	if opts.ShowDocs || opts.DocsSection != "" {
		PrintDocs(opts.DocsSection)
		return 0
	}

	// Phase 3: Validate input
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Error: no input file specified")
		printUsage(os.Stderr)
		return 1
	}

	// Phase 4: Compile the source file
	compiler := NewCompiler(opts)
	if err := compiler.CompileFile(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Compilation failed: %v\n", err)
		return 1
	}

	return 0
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: lotus [flags] <file>")
	fmt.Fprintln(w, "Run 'lotus -h' for help")
}
