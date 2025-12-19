package main

import (
	"fmt"
	"os"
)

func main() {
	// Parse command line flags
	opts, args, err := ParseFlags()
	if err != nil {
		os.Exit(2)
	}

	// Handle version flag
	if opts.ShowVersion {
		fmt.Printf("lotus compiler %s\n", Version)
		os.Exit(0)
	}

	// Validate input file argument
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: lotus [flags] <file>")
		fmt.Fprintln(os.Stderr, "Run 'lotus -h' for help")
		os.Exit(1)
	}

	// Create compiler and compile file
	compiler := NewCompiler(opts)
	if err := compiler.CompileFile(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
