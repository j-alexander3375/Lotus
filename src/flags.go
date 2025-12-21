package main

import (
	"flag"
	"fmt"
	"os"
)

// CompilerOptions holds all compiler configuration
type CompilerOptions struct {
	OutPath     string
	Verbose     bool
	TokenDump   bool
	PrintAsm    bool
	Run         bool
	Trimpath    string
	ShowVersion bool
	IncludeDirs []string
}

// ParseFlags parses command line arguments and returns compiler options
func ParseFlags() (*CompilerOptions, []string, error) {
	opts := &CompilerOptions{}

	fs := flag.NewFlagSet("lotus", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	fs.StringVar(&opts.OutPath, "o", "a.out", "write output to `file`")
	fs.BoolVar(&opts.Verbose, "v", false, "enable verbose logging")
	fs.BoolVar(&opts.TokenDump, "td", false, "print tokens and exit")
	fs.BoolVar(&opts.TokenDump, "token-dump", false, "print tokens and exit")
	fs.BoolVar(&opts.PrintAsm, "S", false, "emit assembly to -o path (or a.s)")
	fs.BoolVar(&opts.Run, "run", false, "build and run the compiled binary")
	fs.StringVar(&opts.Trimpath, "trimpath", "", "remove `prefix` from recorded file paths")
	fs.BoolVar(&opts.ShowVersion, "version", false, "print compiler version and exit")

	fs.Func("I", "add include `dir` to search path (repeatable)", func(val string) error {
		if val != "" {
			opts.IncludeDirs = append(opts.IncludeDirs, val)
		}
		return nil
	})

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: lotus [flags] <file>")
		fmt.Fprintln(os.Stderr, "Flags:")
		fs.PrintDefaults()
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
