# Lotus Compiler - Fresh Systems Language

Lotus is a systems programming language with a few deliberate twists: module imports are string-based (use "io";), returns use ret, and declarations default to type-first bindings (int n = 42;). The compiler emits x86-64 assembly with a modular standard library.

## Documentation

| Document | Purpose |
|----------|---------|
| [STYLE_GUIDE.md](Important_Documentation/STYLE_GUIDE.md) | Naming, formatting, and idioms tuned for Lotus (snake_case for functions/structs/enums; ret preferred for returns; type-first bindings). |
| [STDLIB_AND_IMPORTS.md](Important_Documentation/STDLIB_AND_IMPORTS.md) | How to pull in stdlib modules via use, with import patterns and examples. |
| [STDLIB_IMPLEMENTATION.md](Important_Documentation/STDLIB_IMPLEMENTATION.md) | How stdlib modules are registered and wired into the compiler. |
| [STDLIB_FINAL_SUMMARY.md](Important_Documentation/STDLIB_FINAL_SUMMARY.md) | One-page reference for all stdlib functions across io, mem, math, and str. |
| [DEVELOPMENT.md](Important_Documentation/DEVELOPMENT.md) | Contributor guide and architecture overview. |

## Quick Example (Lotus-flavored)

```lotus
use "io";
use "math";

fn int main() {
  int max_val = max(10, 20);
  printf("Maximum: %d\n", max_val);
  ret 0;
}
```

## Key Features

- Type-first bindings and explicit ret for a distinct Lotus feel
- String-based imports with use "module"; and Rust-like aliasing
- Standard library modules for I/O, memory, math, and strings
- Printf-style formatting verbs: %%, %d, %b, %o, %x/%X, %c, %q, %s, %v
- Structs, enums, and classes with snake_case identifiers
- Error handling via try/catch/finally and throw
- Direct x86-64 GNU assembly output (System V AMD64 ABI)

## Project Structure
```
src/
├── main.go              # Entry point and CLI orchestration
├── compiler.go          # Compilation pipeline management
├── flags.go             # Command-line argument parsing
├── constants.go         # Compiler constants and configuration
│
├── keywords.go          # Token type definitions (100+ tokens)
├── tokenizer.go         # Lexical analysis with multi-char operators
├── parser.go            # Recursive descent parser with precedence
├── ast.go               # Central AST node definitions
├── types.go             # Type system utilities and helpers
├── codegen.go           # Code generation orchestrator
├── diagnostics.go       # Error and warning reporting
│
├── printfuncs.go        # Print function implementations (printf/println/log)
├── memory.go            # Memory management helpers (malloc/free/sizeof)
├── arithmetic.go        # Arithmetic & bitwise operations
├── control_flow.go      # Branching and loops
├── references.go        # References, deref, and assignments
├── functions.go         # User-defined functions and ABI compliance
├── array.go             # Arrays and indexing
├── struct.go            # Struct definitions and field access
├── enum.go              # Enum definitions and resolution
├── class.go             # Classes with fields/methods and new
├── error_handling.go    # try/catch/finally/throw/null
└── stdlib.go            # Stdlib module registration and imports
```

## Standard Library Modules

**io**
```lotus
use "io";
printf("Hello, %s\\n", "Lotus");
println("Auto newline included");
```

Supported printf verbs: %%, %d, %b, %o, %x/%X, %c, %q, %s, %v (ints, chars, strings, quoted strings)

**mem**
```lotus
use "mem";
int* ptr = malloc(sizeof(int) * 4);
free(ptr);
```

**math**
```lotus
use "math";
int top = max(10, 20);
int root = sqrt(16);
```
Implemented: abs, min, max, sqrt, pow, floor, ceil, round, gcd, lcm.

**str**
```lotus
use "str";
int length = len("Hello");
```
Implemented: len, concat, compare, copy, indexOf, contains, startsWith, endsWith.

## Language Notes

- Naming: snake_case for functions/structs/enums; constants stay UPPER_SNAKE_CASE; variables in snake_case.
- Returns: prefer `ret expr;` and `fn int main() { ... ret 0; }` for entrypoints.
- Declarations: type-first by default (`int count = 42;`). Pointers use postfix `*` (`int* buffer`).
- Functions: C-style signatures with `fn <return_type> name(<type> <name>, ...)`.
- Control flow: `if/else`, `while`, and `for (init; cond; update)` blocks require braces.
- Modules: import via `use "module";` and alias with `as` (`use "io::printf" as io_print;`).
- Ownership: allocate with `malloc`, release with `free`; document lifetimes when transferring ownership.

## Sample Patterns

**Function with explicit ret**
```lotus
fn int add(int a, int b) {
  ret a + b;
}
```

**Looping and conditionals**
```lotus
fn int countdown(int start) {
  int n = start;
  while n > 0 {
    printf("%d...\n", n);
    n = n - 1;
  }
  println("liftoff!");
  ret 0;
}
```

**Structs and enums**
```lotus
struct point { int x, y }

enum status { ok = 0, error = -1 }

fn move_point(*point p, int dx, dy) {
    p->x = p->x + dx;
    p->y = p->y + dy;
}
```

## Print Functions (printfuncs.go)
All print functions ultimately call Linux write():
```lotus
printf("Hello, %s\\n", name);
println("Auto newline");
logf("Debug: %d\\n", value);
fatalf("Error: %s\\n", msg);
```

## Arithmetic & Bitwise
```lotus
int sum = 10 + 5;
int diff = x - y;
int masked = flags & 0xFF;
int shifted = value << 2;
```

## Compilation Pipeline
```
Source (.lts)
  ↓
Tokenizer
  ↓
Parser → AST
  ↓
CodeGen → x86-64 assembly
  ↓
Assembler/Linker → Binary (ELF)
```

## Build & Run

```bash
# Build compiler
go build -o lotus ./src

# Compile Lotus source to a.out (default)
./lotus program.lts

# Emit assembly instead of a binary
./lotus -S -o program.s program.lts
```
