# Lotus Compiler - Fresh Systems Language

*** CLAUDE IS USED SOLELY FOR RELEASE WORKFLOWS AND TEST OPTIMIZATION ***

**Lotus** is a systems programming language with deliberate design choices: module imports are string-based (`use "io";`), returns use `ret`, and declarations default to type-first bindings (`int n = 42;`). The compiler uses **LLVM** as its default backend for cross-platform compilation and advanced optimizations.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://go.dev)
[![LLVM](https://img.shields.io/badge/LLVM-Backend-orange.svg)](https://llvm.org)

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Language Overview](#language-overview)
- [Standard Library](#standard-library)
- [Compiler Options](#compiler-options)
- [Documentation](#documentation)
- [Examples](#examples)
- [Contributing](#contributing)

## Features

### Language Features
- **Type-first bindings** - `int count = 42;` with explicit `ret` for returns
- **String-based imports** - `use "module";` with Rust-like aliasing (`as`)
- **Virtual functions** - `vrt fn` for virtual methods, `override fn` for overrides
- **Scope modifiers** - `static`, `lcl` (local), `gbl` (global) for explicit storage control
- **Structs, enums, classes** - snake_case identifiers with full OOP support
- **Error handling** - `try`/`catch`/`finally` and `throw` for exceptions
- **Functional programming** - Pipe operator (`|>`), wrappers/decorators (`@`), currying (`partial`)
- **Void return type** - `fn void foo()` for functions with no return value
- **Type reinterpretation** - `bitcast<Type>(expr)` for low-level bit manipulation

### Compiler Features
- **LLVM backend** (default) - Cross-platform: x86, ARM, RISC-V, WebAssembly
- **Advanced optimizations** - Dead code elimination, constant folding, inlining
- **Enhanced diagnostics** - Line/column tracking, "did you mean?" suggestions
- **Semantic analysis** - Unused variable detection, shadowing warnings
- **Language Server Protocol** - Real-time IDE support with `lotus --lsp`
- **Legacy GCC backend** - Available with `--gcc` flag

### Standard Library Modules
- **io** - printf, println, file operations
- **mem** - malloc, free, sizeof, memcpy, memset
- **math** - abs, min, max, sqrt, pow, gcd, lcm, floor, ceil, round
- **str** - len, concat, compare, indexOf, contains, toUpper, toLower, trim, split, replace
- **hash** - djb2, fnv1a, crc32, murmur
- **net** - socket, connect_ipv4, bind_ipv4, listen, accept, send, recv, setsockopt
- **collections** - dynamic arrays, stacks, queues, deques, heaps
- **json** - parse, stringify, get, array operations
- **format** - sprintf, snprintf, pad_left, pad_right
- **random** - rand, rand_range, seed, rand_float, rand_bool, shuffle, choice

## Quick Start

```lotus
use "io";
use "math";

fn int main() {
    int max_val = max(10, 20);
    printf("Maximum: %d\n", max_val);
    ret 0;
}
```

```bash
# Compile and run
./lotus --run program.lts

# Just compile
./lotus program.lts -o myprogram
./myprogram
```

## Installation

### From Source (Recommended)

```bash
# Clone repository
git clone https://github.com/yourusername/lotus.git
cd lotus

# Build compiler
go build -o lotus ./src

# Verify installation
./lotus --version
```

### Arch Linux (AUR)

```bash
# Using yay
yay -S lotus-lang

# Or with paru
paru -S lotus-lang
```

### Dependencies

- **Go 1.21+** - Compiler is written in Go
- **LLVM 15+** - Backend (automatically linked via go-llvm)
- **GCC/Clang** - For linking final binaries

## Language Overview

### Variables and Constants

```lotus
// Type-first declarations
int count = 42;
string name = "Lotus";
bool active = true;
float pi = 3.14159;

// Constants
const int MAX_SIZE = 1024;

// Scope modifiers
static int persistent = 0;    // Persists across function calls
gbl int shared = 100;         // Explicit global
lcl int temp = 5;             // Explicit local (stack)
```

### Functions

```lotus
// Regular function
fn int add(int a, int b) {
    ret a + b;
}

// Virtual function (for polymorphism)
vrt fn int compute() {
    ret 42;
}

// Override function
override fn int compute() {
    ret 100;
}

// Static function (file-local)
static fn int helper() {
    ret 0;
}
```

### Control Flow

```lotus
// If/else
if x > 0 {
    printf("positive\n");
} else {
    printf("non-positive\n");
}

// While loop
while n > 0 {
    n = n - 1;
}

// For loop
for int i = 0; i < 10; i++ {
    printf("%d\n", i);
}
```

### Structs and Enums

```lotus
struct point {
    int x;
    int y;
}

enum status {
    ok = 0,
    error = -1
}

fn int main() {
    point p;
    p.x = 10;
    p.y = 20;
    ret status::ok;
}
```

### Error Handling

```lotus
fn int risky_operation() {
    try {
        // Code that might fail
        int result = dangerous_call();
        ret result;
    } catch {
        printf("Error occurred\n");
        ret -1;
    } finally {
        cleanup();
    }
}
```

### Functional Programming

```lotus
// Pipe operator - chain function calls
int result = 5 |> double |> addOne;  // Becomes addOne(double(5))

// Wrappers (decorators) - wrap functions with common behavior
wrap fn logging(fn wrapped) {
    printf("[LOG] Entering function\n");
    wrapped();
    printf("[LOG] Exiting function\n");
}

// Apply wrapper to a function
@logging
fn void greet() {
    printf("Hello, World!\n");
}

// Currying with partial application
fn int add(int a, int b) {
    ret a + b;
}

fn int main() {
    // Call partially applied function
    int result = partial(add, 5)(3);  // 8
    ret 0;
}
```

## Standard Library

### io Module
```lotus
use "io";

printf("Hello, %s\n", "World");
println("Auto newline");
fprintf(stderr, "Error: %s\n", msg);
```

**Printf verbs:** `%%`, `%d`, `%b`, `%o`, `%x`/`%X`, `%c`, `%q`, `%s`, `%v`

### mem Module
```lotus
use "mem";

int* buffer = malloc(sizeof(int) * 100);
memset(buffer, 0, 100);
free(buffer);
```

### math Module
```lotus
use "math";

int m = max(10, 20);      // 20
int r = sqrt(16);         // 4
int p = pow(2, 8);        // 256
int g = gcd(48, 18);      // 6
```

### str Module
```lotus
use "str";

int length = len("Hello");           // 5
bool has = contains("Hello", "ell"); // true
int pos = indexOf("Hello", "l");     // 2
string upper = toUpper("hello");     // "HELLO"
string lower = toLower("WORLD");     // "world"
string trimmed = trim("  hi  ");     // "hi"
```

### random Module
```lotus
use "random";

seed(12345);                         // Set random seed
int r = rand();                      // Random int
int n = rand_range(1, 100);          // Random int in [1, 100]
int m = rand_n(10);                  // Random int in [0, 10)
bool b = rand_bool();                // Random boolean
shuffle(arr, 10);                    // Shuffle array in place
int pick = choice(arr, 10);          // Random element from array
```

## Compiler Options

```
Usage: lotus [options] <source.lts>

Output Options:
  -o <file>        Output file name (default: a.out)
  -S               Emit assembly instead of binary
  --emit-llvm      Emit LLVM IR
  --run            Compile and run immediately

Backend Options:
  --gcc            Use legacy GCC/assembly backend
  --target <triple> Cross-compile for target (e.g., aarch64-linux-gnu)

Optimization:
  -O0              No optimization
  -O1              Basic optimization
  -O2              Standard optimization (default)
  -O3              Aggressive optimization

Debugging:
  --ast-dump       Print AST structure
  --stats          Show compilation statistics
  --timing         Show phase timing
  -v, --verbose    Verbose output

Other:
  --version        Show version
  --help           Show help
```

## Documentation

| Document | Description |
|----------|-------------|
| [STYLE_GUIDE.md](Important_Documentation/STYLE_GUIDE.md) | Naming conventions, formatting, idioms |
| [STDLIB_AND_IMPORTS.md](Important_Documentation/STDLIB_AND_IMPORTS.md) | Import patterns and module usage |
| [STDLIB_FINAL_SUMMARY.md](Important_Documentation/STDLIB_FINAL_SUMMARY.md) | Complete stdlib reference |
| [DEVELOPMENT.md](Important_Documentation/DEVELOPMENT.md) | Contributor guide and architecture |

## Examples

See the [examples/](examples/) directory for complete programs:

- `control_flow_if.lts` - Conditionals
- `control_flow_for.lts` - Loops
- `control_flow_while.lts` - While loops

See [tests/](tests/) for more comprehensive examples.

## Project Structure

```
lotus/
├── src/                    # Compiler source code
│   ├── main.go            # Entry point
│   ├── compiler.go        # Compilation pipeline
│   ├── tokenizer.go       # Lexical analysis
│   ├── parser.go          # Parsing
│   ├── llvm_codegen.go    # LLVM code generation
│   ├── llvm_optimizer.go  # LLVM optimizations
│   ├── llvm_stdlib.go     # LLVM stdlib implementations
│   └── ...
├── ext/                    # VS Code extension
├── tests/                  # Test files
├── examples/              # Example programs
└── Important_Documentation/
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open a Pull Request

See [DEVELOPMENT.md](Important_Documentation/DEVELOPMENT.md) for architecture details.

## License

MIT License - see [LICENSE](LICENSE) for details.
