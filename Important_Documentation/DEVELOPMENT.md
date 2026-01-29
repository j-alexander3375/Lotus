# Lotus Compiler - Development Summary

> Last Updated: January 2026 | Current Phase: 10 (Developer Tools & Optimization)

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Language Features](#language-features)
4. [Standard Library](#standard-library)
5. [Development History](#development-history)
6. [Project Structure](#project-structure)
7. [Build & Usage](#build--usage)
8. [Future Roadmap](#future-roadmap)

---

## Overview

Lotus is a compiled systems programming language with a modern syntax inspired by Rust, Go, and C. It features an LLVM-based backend (default), comprehensive standard library, and full IDE support via Language Server Protocol.

**Key Highlights:**
- LLVM backend with cross-compilation support
- 11 stdlib modules with 100+ functions
- LSP integration for VS Code
- AST-level and peephole optimizations
- Comprehensive error diagnostics with suggestions

---

## Architecture

```
Input (.lts file)
       ↓
  [Tokenizer]  →  Token Stream (100+ token types)
       ↓
    [Parser]   →  AST (with location tracking)
       ↓
  [Optimizer]  →  Optimized AST (constant folding, dead code elimination)
       ↓
[Code Generator] → LLVM IR (default) or x86-64 Assembly (--gcc flag)
       ↓
    [LLVM/GCC] →  Binary Executable
```

**Backend Options:**
- **LLVM (default)**: Full optimization pipeline, cross-compilation, `-O0` to `-O3`
- **GCC (legacy)**: x86-64 GNU assembly with System V ABI, `--gcc` flag

---

## Language Features

### Core Syntax

- Type-first variable declarations: `int x = 1;`
- Constants: `const int MAX = 100;`
- Functions: `fn int add(int a, int b) { ret a + b; }`
- Void functions: `fn void greet() { println("Hello"); }`

### Data Types

| Category | Types |
|----------|-------|
| Integers | `int8`, `int16`, `int32`, `int64`, `uint8`, `uint16`, `uint32`, `uint64` |
| Floats | `float` |
| Other | `bool`, `char`, `string`, `void` |
| Compound | `struct`, `enum`, `class`, `arrays`, `pointers` |

### Operators

| Category | Operators |
|----------|-----------|
| Arithmetic | `+`, `-`, `*`, `/`, `%` |
| Bitwise | `&`, `\|`, `^`, `~`, `<<`, `>>` |
| Comparison | `==`, `!=`, `<`, `<=`, `>`, `>=` |
| Logical | `&&`, `\|\|`, `!` |
| Assignment | `=`, `+=`, `-=`, `*=`, `/=`, `%=` |
| Other | `++`, `--`, `? :` (ternary), `\|>` (pipe) |

### Control Flow

- Conditionals: `if`, `else`
- Loops: `while`, `for`
- Exception handling: `try`, `catch`, `finally`, `throw`, `null`

### Advanced Features

| Feature | Syntax | Description |
|---------|--------|-------------|
| Pipe operator | `value \|> fn1 \|> fn2` | Chain function calls |
| Decorators | `@timing fn void foo() {}` | Apply function wrappers |
| Wrappers | `wrap fn timing(fn wrapped) {}` | Define decorators |
| Currying | `partial(add, 5)(3)` | Partial function application |
| Bitcast | `bitcast<int64>(floatVal)` | Bit reinterpretation |
| Virtual functions | `vrt fn foo()` / `override fn foo()` | Polymorphism |
| Storage classes | `static`, `lcl`, `gbl` | Variable scope control |
| Optional types | `Some(value)`, `None` | Type-safe null handling |
| Pattern matching | `match expr { case ... }` | Haskell-style pattern matching |
| Generics/Templates | `template<typename T>` | C++-style generic programming |

### Optional Types

Lotus provides built-in optional types for representing values that may or may not exist:

```lotus
// Creating optionals
Some(42)        // Optional with value
None            // Empty optional

// Optionals are represented internally as structs: { i1 has_value, T value }
// This allows type-safe handling of potentially missing values
```

### Pattern Matching

Haskell-style pattern matching with `match` expressions:

```lotus
// Literal matching
match value {
    case 42 => println("Found 42!"),
    case 100 => println("Found 100!"),
    default => println("Other")
}

// Range patterns
match score {
    case 90..100 => println("Grade: A"),
    case 80..89 => println("Grade: B"),
    default => println("Grade: F")
}

// Binding patterns with guards
match num {
    case x when x > 10 => println("Greater than 10"),
    case x when x < 10 => println("Less than 10"),
    default => println("Equal to 10")
}

// Wildcard pattern
match value {
    case _ => println("Matches anything")
}
```

**Pattern Types:**
- Literal patterns: `case 42 =>`, `case "hello" =>`
- Range patterns: `case 1..10 =>`
- Binding patterns: `case x =>` (binds value to variable)
- Wildcard patterns: `case _ =>` (matches anything)
- Guard clauses: `case x when condition =>`

### Generics/Templates

C++-style template system for generic programming:

```lotus
// Generic function
template<typename T>
fn T maximum(T a, T b) {
    if a > b {
        ret a;
    }
    ret b;
}

// Usage with type inference
int maxInt = maximum(10, 20);      // Infers T = int
float maxFloat = maximum(3.14, 2.71);  // Infers T = float

// Generic struct (planned)
template<typename T>
struct Box {
    T value;
}
```

**Features:**
- Type parameters: `template<typename T>`
- Type inference from arguments
- Multiple type parameters: `template<typename K, typename V>`
- Works with functions (structs partially implemented)

### Import System

```lotus
use "io"                    // Import module
use "math::sqrt"            // Specific function
use "str::*"                // Wildcard import
use "collections" as col    // Aliased import
```

---

## Standard Library

| Module | Description | Key Functions |
|--------|-------------|---------------|
| **io** | Input/output | `print`, `println`, `printf`, `fprintf`, `sprintf` |
| **mem** | Memory management | `malloc`, `free`, `sizeof`, `memcpy`, `memset`, `mmap`, `munmap` |
| **math** | Mathematics | `abs`, `min`, `max`, `sqrt`, `pow`, `floor`, `ceil`, `round`, `gcd`, `lcm` |
| **str** | String manipulation | `len`, `concat`, `compare`, `copy`, `indexOf`, `contains`, `substring`, `split`, `join`, `replace`, `toLower`, `toUpper`, `trim`, `startsWith`, `endsWith` |
| **num** | Type conversion | `toInt8`–`toUint64`, `toBool`, `toInt`, `toFloat` |
| **file** | File I/O | `open`, `close`, `read`, `write`, `seek`, `stat`, `exists` |
| **time** | Time operations | `now`, `sleep`, `millis`, `nanos`, `clock`, `gmtime`, `localtime` |
| **hash** | Hashing algorithms | `djb2`, `fnv1a`, `crc32`, `murmur3`, `sha256`, `md5` |
| **collections** | Data structures | Arrays, Stack, Queue, Deque, Heap, HashMap, HashSet, SortedSet, SortedMap |
| **net** | Networking | `socket`, `connect_ipv4`, `connect_ipv6`, `bind_ipv4`, `send`, `recv`, `close`, `resolve` |
| **http** | HTTP client | `get`, `post`, `parse_status`, `get_header`, `get_body`, connection pooling |
| **rand** | Random numbers | `rand`, `rand_range`, `rand_n`, `seed`, `rand_float`, `rand_bool`, `shuffle`, `choice` |

---

## Development History

### Phase 1: Foundation

*Initial compiler implementation*

- **Tokenizer**: 100+ token types, string/int/float/bool literals, multi-char operators
- **Parser**: Recursive descent, AST generation with all node types
- **Code Generator**: x86-64 GNU assembly (AT&T syntax), stack frames, System V ABI
- **Basic constructs**: Variables, functions (`fn`/`ret`), control flow, struct/enum/class

### Phase 2: Refactoring & Constants

*December 2025*

- Created modular file structure: `ast.go`, `constants.go`, `types.go`
- Added `const` keyword for compile-time constants
- Constants stored in `.data` section with RIP-relative addressing
- Eliminated magic numbers and improved naming conventions

### Phase 3: Import System & Stdlib

*December 2025*

- Implemented `use`/`as` import syntax (Rust-inspired)
- Created `stdlib.go` with module registration system
- Initial stdlib modules: io, mem, math, str, num
- `printf` format specifiers: `%%`, `%d`, `%b`, `%o`, `%x`, `%X`, `%c`, `%q`, `%s`, `%v`
- Import validation and compile-time resolution

### Phase 4: Extended Stdlib

*December 2025 (v1.2.1–v1.2.5)*

- **String module**: Full implementation of 15 functions
- **File I/O**: Linux syscalls (open, close, read, write, seek, stat, exists)
- **Time module**: Timestamps, sleep, gmtime/localtime
- **Hash module**: CRC32, FNV-1a, DJB2, MurmurHash3, SHA-256, MD5
- **Collections**: Dynamic arrays, stacks, queues, heaps, hashmaps (int/string keys), BST-based sorted sets/maps
- **Net/HTTP**: Socket primitives, IPv4/IPv6, UDP, DNS resolution, HTTP GET/POST

### Phase 5: Tooling & Diagnostics

*December 2025*

- **Compilation statistics**: `--stats` flag, timing per phase, AST metrics
- **AST utilities**: `DumpAST()`, `CountASTNodes()`, `AnalyzeAST()`
- **Enhanced diagnostics**: Color-coded output (ANSI), caret positioning
- **New flags**: `--stats`, `--timing`, `--ast-dump`, `-q`/`--quiet`
- Comprehensive test suite and benchmarks

### Phase 6: Code Optimization

*December 2025*

- **AST-level optimizations** (`optimizer.go`):
  - Constant folding: `2 + 3` → `5`
  - Strength reduction: `x * 8` → `x << 3`
  - Identity removal: `x + 0` → `x`, `x * 1` → `x`
- **Peephole optimizations** (`peephole.go`):
  - Redundant move elimination
  - Dead store elimination
  - Push-pop cancellation
  - Zero loading: `movq $0, %rax` → `xorq %rax, %rax`
  - Increment simplification: `addq $1, %rax` → `incq %rax`

### Phase 7: Error Messages & Diagnostics

*December 2025*

- Line/column tracking in all tokens
- Human-readable error messages with error codes (E01xx–E03xx)
- **"Did you mean?" suggestions**: Levenshtein distance for typos
- Enhanced semantic analyzer: Unused variable detection, shadowing warnings
- AST location tracking via `BaseNode` with `Loc()` method
- Diagnostic categories: Syntax, Semantic, Unused, Shadow, Deprecated

### Phase 8: LLVM Backend

*December 2025*

- **LLVM infrastructure**: `llvm_codegen.go`, `llvm_types.go`
- LLVM as **default** backend (tinygo.org/x/go-llvm bindings)
- Type mapping between Lotus and LLVM types
- **Backend flags**: `--gcc`, `--emit-llvm`, `-O0`–`-O3`, `--target`
- Cross-compilation: x86, ARM, RISC-V, WebAssembly targets
- **OOP features**: `vrt`/`override` for virtual functions
- **Storage classes**: `static`, `lcl`, `gbl`

### Phase 9: Functional Programming

*January 2026*

- **Pipe operator** (`|>`): Chain function calls
- **Decorators/Wrappers**: `@decorator` syntax, `wrap fn` definitions
- **Currying**: `partial(fn, arg)` for partial application
- **Void return type**: `fn void foo()`
- **Random module**: Full PRNG suite with xorshift64
- **Extended string functions**: copy, compare, indexOf, substring, etc.
- **LLVM advanced features**: Ternary, compound assignment, ++/--, sizeof

### Phase 10: Developer Tools ✅ (Current)

*January 2026*

- **Language Server Protocol**: Full LSP 3.17 in `lsp_server.go`
  - Real-time diagnostics
  - Auto-completion (keywords, builtins, symbols)
  - Hover documentation
  - Go-to-definition
  - Symbol outline
  - VS Code extension integration
- **Dead code elimination**: `EliminateDeadCode()`, `EliminateUnusedFunctions()`
- **Bitcast/Transmute**: Type reinterpretation for low-level manipulation
- **Float arithmetic fix**: Proper LLVM float instructions

### Release 1.5.1: Stdlib Bug Fixes & Networking

*January 2026*

- **Bug fixes**:
  - `trim()`: Fixed leading whitespace removal
  - `heap_int_push/pop`: Proper min-heap bubble-up/sift-down
  - `deque_int_push_front`: Proper element shifting
  - `queue_int_dequeue`: Proper FIFO with element shifting
  - IPv4 address byte ordering in `connect_ipv4` and `bind_ipv4`
- **New LLVM stdlib implementations**:
  - Math: `gcd`, `lcm`, `floor`, `ceil`, `round`
  - Memory: `memset`, `memcpy`
  - Hash: `djb2`, `fnv1a`, `crc32`, `murmur`
  - Time: `clock`
  - Number conversions: `toInt8`, `toInt16`, `toInt32`, `toInt64`, `toUint32`, `toBool`
- **Networking enhancements**:
  - Server functions: `bind_ipv4`, `listen`, `accept`, `setsockopt`
  - Full TCP client/server support on localhost
  - Proper network byte order handling for IP addresses

---

## Project Structure

```
src/
├── main.go              - CLI entry point
├── compiler.go          - Compilation pipeline
├── flags.go             - Command-line flag parsing
├── constants.go         - Compiler constants
├── keywords.go          - Token definitions (100+ types)
├── tokenizer.go         - Lexical analysis
├── parser.go            - Syntactic analysis
├── ast.go               - AST node definitions
├── types.go             - Type system utilities
├── semantic.go          - Semantic analysis
├── optimizer.go         - AST optimizations
├── peephole.go          - Assembly peephole optimizations
├── codegen.go           - Code generation orchestrator
├── llvm_codegen.go      - LLVM IR generation
├── llvm_types.go        - LLVM type mapping
├── llvm_advanced.go     - Advanced LLVM features
├── lsp_server.go        - Language Server Protocol
├── diagnostics.go       - Error reporting
├── error_messages.go    - Error formatting
├── stdlib.go            - Standard library registration
├── stats.go             - Compilation statistics
├── ast_utils.go         - AST utilities
├── printfuncs.go        - Print function codegen
├── arithmetic.go        - Arithmetic operations
├── control_flow.go      - Control flow codegen
├── references.go        - Pointer operations
├── functions.go         - Function codegen
├── memory.go            - Memory operations
├── array.go             - Array operations
├── struct.go            - Struct handling
├── enum.go              - Enum handling
├── class.go             - Class/OOP support
└── error_handling.go    - Exception handling
```

---

## Build & Usage

### Building the Compiler

```bash
cd src && go build -o ../lotus
```

### Command-Line Options

| Command | Description |
|---------|-------------|
| `./lotus input.lts` | Compile to binary (LLVM backend) |
| `./lotus -run input.lts` | Compile and run immediately |
| `./lotus -S input.lts` | Generate assembly only |
| `./lotus --emit-llvm input.lts` | Emit LLVM IR |
| `./lotus --gcc input.lts` | Use legacy GCC backend |
| `./lotus -O2 input.lts` | Optimization level (0–3) |
| `./lotus --target <triple> input.lts` | Cross-compile |
| `./lotus --lsp` | Start Language Server |
| `./lotus -td input.lts` | Dump tokens (debug) |
| `./lotus --ast-dump input.lts` | Dump AST (debug) |
| `./lotus --stats input.lts` | Show compilation statistics |
| `./lotus -v input.lts` | Verbose output |
| `./lotus --version` | Show version |

### Example Code

```lotus
use "io"
use "math"

const int MAX_SIZE = 100;

fn int factorial(int n) {
    if (n <= 1) {
        ret 1;
    }
    ret n * factorial(n - 1);
}

fn void main() {
    int result = factorial(5);
    printf("5! = %d\n", result);
    
    // Pipe operator example
    float val = 10 |> factorial |> sqrt;
}
```

---

## Future Roadmap

### Planned Features

| Feature | Status | Notes |
|---------|--------|-------|
| Regular expressions | ⏳ | Pattern matching in strings |
| JSON module | ⏳ | Parse/serialize JSON |
| Width/padding for printf | ⏳ | Format specifier enhancements |
| Register allocation improvements | ⏳ | Better code generation |
| Inline function expansion | ⏳ | Performance optimization |
| Generics/Templates | ✅ | C++-style templates with type inference |
| Union/option types | ✅ | `Some(value)`, `None` |
| Pattern matching | ✅ | Haskell-style `match` with literals, ranges, bindings, guards |
| Constructor patterns | ⏳ | Pattern matching on `Some(x)`, enum variants |
| Debug/trace hooks | ⏳ | Runtime debugging |
| Package manager | ⏳ | Dependency management |
| Build system integration | ⏳ | Make/CMake/Meson support |

---

*For detailed stdlib documentation, see [STDLIB_AND_IMPORTS.md](STDLIB_AND_IMPORTS.md)*

