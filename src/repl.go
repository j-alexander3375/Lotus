package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
)

// repl.go - GHCi-style interactive REPL for the Lotus interpreter.
// Prompt: "lts ~ "   Continuation: "      | "

const replBanner = `Lotus %s REPL  (interpreter)
Type :help for commands, :quit to exit.
`

const replHelp = `Commands:
  :q,  :quit          Exit the REPL
  :h,  :help          Show this message
  :reset              Reset interpreter state (clears all bindings)
  :load <file>        Load and evaluate a .lts source file
  :type <expr>        Show the type of an expression
  :! <cmd>            Run a shell command and print its output

Multi-line: open a '{' to start a block; close it with '}' to evaluate.
`

// RunREPL starts the interactive REPL. It never returns normally — it exits via
// :quit, Ctrl-D, or Ctrl-C.
func RunREPL() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:            "lts ~ ",
		HistoryFile:       os.ExpandEnv("$HOME/.lotus_history"),
		HistorySearchFold: true,
		InterruptPrompt:   "^C",
		EOFPrompt:         ":quit",
	})
	if err != nil {
		// Non-TTY or readline unavailable — fall back to plain stdin.
		runReplPlain()
		return
	}
	defer rl.Close()

	fmt.Printf(replBanner, Version)

	interp := NewInterpreter()
	runLoop(interp, rl)
}

// runLoop is the core read-eval-print loop. Separated so runReplPlain can reuse
// the evaluation logic with a trivial readline stub.
func runLoop(interp *Interpreter, rl *readline.Instance) {
	var buf strings.Builder
	depth := 0

	for {
		prompt := "lts ~ "
		if buf.Len() > 0 {
			prompt = "      | "
		}
		rl.SetPrompt(prompt)

		line, err := rl.Readline()
		if err != nil {
			// Ctrl-D or closed stdin
			fmt.Println("\nGoodbye!")
			return
		}

		trimmed := strings.TrimRight(line, " \t")

		// Colon commands are only accepted at the start of a new input (not mid-block).
		if buf.Len() == 0 && strings.HasPrefix(trimmed, ":") {
			if dispatchCommand(trimmed, interp, rl) {
				return // :quit
			}
			continue
		}

		buf.WriteString(line)
		buf.WriteByte('\n')
		depth += netBraceDepth(line)

		// Keep collecting if we're still inside an open block.
		if depth > 0 {
			continue
		}
		if depth < 0 {
			depth = 0
		}

		input := buf.String()
		buf.Reset()

		if strings.TrimSpace(input) == "" {
			continue
		}

		evalAndPrint(input, interp)
	}
}

// evalAndPrint tokenizes, parses, and evaluates input, printing the result.
func evalAndPrint(input string, interp *Interpreter) {
	tokens := Tokenize(input)
	parser := NewParser(tokens)
	stmts, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  parse error: %v\n", err)
		return
	}
	val, err := interp.EvalStatements(stmts, interp.globals)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  error: %v\n", err)
		return
	}
	if val != nil && val.Kind != KindVoid {
		fmt.Println(val.Display())
	}
}

// dispatchCommand handles a leading-colon command. Returns true when the REPL
// should exit.
func dispatchCommand(line string, interp *Interpreter, rl *readline.Instance) bool {
	// :! passes everything after the prefix directly to the shell.
	if strings.HasPrefix(line, ":!") {
		shellCmd := strings.TrimSpace(line[2:])
		if shellCmd == "" {
			fmt.Fprintln(os.Stderr, "  usage: :! <shell command>")
			return false
		}
		replRunShell(shellCmd)
		return false
	}

	parts := strings.Fields(line)
	if len(parts) == 0 {
		return false
	}
	switch parts[0] {
	case ":q", ":quit":
		fmt.Println("Goodbye!")
		return true

	case ":h", ":help":
		fmt.Print(replHelp)

	case ":reset":
		*interp = *NewInterpreter()
		fmt.Println("State reset.")

	case ":load":
		if len(parts) < 2 {
			fmt.Fprintln(os.Stderr, "  usage: :load <file>")
			break
		}
		replLoadFile(parts[1], interp)

	case ":type":
		if len(parts) < 2 {
			fmt.Fprintln(os.Stderr, "  usage: :type <expr>")
			break
		}
		expr := strings.Join(parts[1:], " ")
		tokens := Tokenize(expr)
		parser := NewParser(tokens)
		stmts, err := parser.Parse()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  parse error: %v\n", err)
			break
		}
		if len(stmts) == 0 {
			break
		}
		val, err := interp.evalNode(stmts[0], interp.globals)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  error: %v\n", err)
			break
		}
		fmt.Println(val.TypeName())

	default:
		fmt.Fprintf(os.Stderr, "  unknown command %q — try :help\n", parts[0])
	}
	return false
}

// replLoadFile reads a .lts file and evaluates it in the current interpreter.
func replLoadFile(path string, interp *Interpreter) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  error: %v\n", err)
		return
	}
	tokens := Tokenize(string(data))
	parser := NewParser(tokens)
	stmts, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  parse error in %s: %v\n", path, err)
		return
	}
	_, err = interp.EvalStatements(stmts, interp.globals)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  error in %s: %v\n", path, err)
		return
	}
	fmt.Printf("Loaded %s\n", path)
}

// replRunShell runs a shell command, wiring its stdout/stderr directly to the
// terminal so output (including colour, progress bars, etc.) streams live.
func replRunShell(command string) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// Only print the error when it's not just a non-zero exit code
		// (the command itself will have printed its own error message).
		if _, isExit := err.(*exec.ExitError); !isExit {
			fmt.Fprintf(os.Stderr, "  shell error: %v\n", err)
		}
	}
}

// netBraceDepth returns the net { minus } count for a line, ignoring strings.
func netBraceDepth(line string) int {
	depth := 0
	inStr := false
	escape := false
	for _, c := range line {
		if escape {
			escape = false
			continue
		}
		if c == '\\' && inStr {
			escape = true
			continue
		}
		if c == '"' {
			inStr = !inStr
			continue
		}
		if inStr {
			continue
		}
		if c == '{' {
			depth++
		} else if c == '}' {
			depth--
		}
	}
	return depth
}

// runReplPlain is a minimal fallback for non-TTY environments (no readline).
func runReplPlain() {
	// Build a readline.Instance backed by stdin with no history or editing.
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "lts ~ ",
		Stdin:       os.Stdin,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		FuncIsTerminal: func() bool { return false },
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "REPL unavailable:", err)
		return
	}
	defer rl.Close()
	fmt.Printf(replBanner, Version)
	runLoop(NewInterpreter(), rl)
}
