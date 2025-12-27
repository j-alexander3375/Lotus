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

	// Tooling enhancements
	ShowStats  bool // Display compilation statistics (--stats)
	Quiet      bool // Suppress all non-error output (-q, --quiet)
	TimingInfo bool // Show detailed phase timing (--timing)
	ASTDump    bool // Print AST and exit (--ast-dump)

	// Documentation
	ShowDocs    bool   // Show offline documentation (-docs, --docs)
	DocsSection string // Specific documentation section to show

	// Warning and error control
	Wall           bool // Enable all warnings (-Wall)
	Werror         bool // Treat warnings as errors (-Werror)
	WarnUnused     bool // Warn about unused variables (-Wunused)
	WarnShadow     bool // Warn about variable shadowing (-Wshadow)
	WarnImplicit   bool // Warn about implicit conversions (-Wimplicit)
	WarnDeprecated bool // Warn about deprecated features (-Wdeprecated)
	NoWarn         bool // Suppress all warnings (-w)
	MaxErrors      int  // Maximum errors before stopping (--max-errors)
	ColorOutput    bool // Enable colored output (--color)
	NoColor        bool // Disable colored output (--no-color)
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
	fs.BoolVar(&opts.ASTDump, "ast-dump", false, "print AST and exit")

	// Tooling options
	fs.BoolVar(&opts.ShowStats, "stats", false, "display compilation statistics")
	fs.BoolVar(&opts.Quiet, "q", false, "suppress non-error output")
	fs.BoolVar(&opts.Quiet, "quiet", false, "suppress non-error output")
	fs.BoolVar(&opts.TimingInfo, "timing", false, "show detailed phase timing")

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

	// Documentation
	fs.BoolVar(&opts.ShowDocs, "docs", false, "show offline documentation")
	fs.StringVar(&opts.DocsSection, "docs-section", "", "show specific docs section (syntax, stdlib, types, examples)")

	// Warning flags
	fs.BoolVar(&opts.Wall, "Wall", false, "enable all warnings")
	fs.BoolVar(&opts.Werror, "Werror", false, "treat warnings as errors")
	fs.BoolVar(&opts.WarnUnused, "Wunused", false, "warn about unused variables")
	fs.BoolVar(&opts.WarnShadow, "Wshadow", false, "warn about variable shadowing")
	fs.BoolVar(&opts.WarnImplicit, "Wimplicit", false, "warn about implicit conversions")
	fs.BoolVar(&opts.WarnDeprecated, "Wdeprecated", false, "warn about deprecated features")
	fs.BoolVar(&opts.NoWarn, "w", false, "suppress all warnings")
	fs.IntVar(&opts.MaxErrors, "max-errors", 20, "maximum number of errors before stopping")
	fs.BoolVar(&opts.ColorOutput, "color", true, "enable colored output")
	fs.BoolVar(&opts.NoColor, "no-color", false, "disable colored output")

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
		fmt.Fprintln(os.Stderr, "  lotus --stats program.lts      # Show compilation stats")
		fmt.Fprintln(os.Stderr, "  lotus --timing program.lts     # Show phase timing")
		fmt.Fprintln(os.Stderr, "  lotus --ast-dump program.lts   # Dump AST structure")
		fmt.Fprintln(os.Stderr, "  lotus -docs                    # Show documentation")
		fmt.Fprintln(os.Stderr, "  lotus -docs-section stdlib     # Show stdlib docs")
		fmt.Fprintln(os.Stderr, "  lotus -Wall program.lts        # Enable all warnings")
		fmt.Fprintln(os.Stderr, "  lotus -Werror program.lts      # Warnings as errors")
	}

	// Normalize args to accept various flag formats
	raw := os.Args[1:]
	norm := make([]string, 0, len(raw))
	for _, a := range raw {
		switch a {
		case "--token-dump":
			norm = append(norm, "-token-dump")
		case "--ast-dump":
			norm = append(norm, "-ast-dump")
		case "--stats":
			norm = append(norm, "-stats")
		case "--quiet":
			norm = append(norm, "-quiet")
		case "--timing":
			norm = append(norm, "-timing")
		case "--docs":
			norm = append(norm, "-docs")
		case "--docs-section":
			norm = append(norm, "-docs-section")
		default:
			norm = append(norm, a)
		}
	}

	if err := fs.Parse(norm); err != nil {
		return nil, nil, err
	}

	return opts, fs.Args(), nil
}
