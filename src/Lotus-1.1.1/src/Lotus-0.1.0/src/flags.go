package main

import (
	"flag"
	"fmt"
	"os"
)

// flags.go - Command-line flag parsing and compiler options
// This file manages all compiler configuration and command-line arguments.

// CompilerOptions holds all compiler configuration settings
type CompilerOptions struct {
	OutPath       string   // Output file path (-o)
	Verbose       bool     // Enable verbose logging (-v)
	TokenDump     bool     // Print tokens and exit (-td, --token-dump)
	PrintAsm      bool     // Emit assembly instead of binary (-S)
	RunAfterBuild bool     // Build and run the binary (-run)
	Trimpath      string   // Remove prefix from recorded file paths (--trimpath)
	ShowVersion   bool     // Print version and exit (--version)
	IncludeDirs   []string // Include directories for imports (-I)
}

// Version is the current compiler version
const Version = CompilerVersion

// ParseFlags parses command-line arguments and returns compiler options.
// Returns (options, positional args, error)
func ParseFlags() (*CompilerOptions, []string, error) {
	opts := &CompilerOptions{}

	fs := flag.NewFlagSet("lotus", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	// Output options
	fs.StringVar(&opts.OutPath, "o", "a.out", "write output to `file`")
	fs.BoolVar(&opts.PrintAsm, "S", false, "emit assembly to -o path (or a.s)")

	// Debug options
	fs.BoolVar(&opts.Verbose, "v", false, "enable verbose logging")
	fs.BoolVar(&opts.TokenDump, "td", false, "print tokens and exit")
	fs.BoolVar(&opts.TokenDump, "token-dump", false, "print tokens and exit")

	// Execution options
	fs.BoolVar(&opts.RunAfterBuild, "run", false, "build and run the compiled binary")

	// Path options
	fs.StringVar(&opts.Trimpath, "trimpath", "", "remove `prefix` from recorded file paths")
	fs.Func("I", "add include `dir` to search path (repeatable)", func(val string) error {
		if val != "" {
			opts.IncludeDirs = append(opts.IncludeDirs, val)
		}
		return nil
	})

	// Version
	fs.BoolVar(&opts.ShowVersion, "version", false, "print compiler version and exit")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: lotus [flags] <file>")
		fmt.Fprintln(os.Stderr, "\nFlags:")
		fs.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  lotus program.lts              # Compile to a.out")
		fmt.Fprintln(os.Stderr, "  lotus -o myapp program.lts     # Compile to myapp")
		fmt.Fprintln(os.Stderr, "  lotus -S program.lts           # Generate assembly")
		fmt.Fprintln(os.Stderr, "  lotus -run program.lts         # Compile and run")
		fmt.Fprintln(os.Stderr, "  lotus -td program.lts          # Dump tokens")
	}

	// Normalize args to accept --token-dump
	raw := os.Args[1:]
	norm := make([]string, 0, len(raw))
	for _, a := range raw {
		if a == "--token-dump" {
			norm = append(norm, "-token-dump")
		} else {
			norm = append(norm, a)
		}
	}

	if err := fs.Parse(norm); err != nil {
		return nil, nil, err
	}

	return opts, fs.Args(), nil
}
