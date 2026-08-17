# Lotus Compiler - Fresh Systems Language

**CLAUDE IS USED SOLELY FOR RELEASE WORKFLOWS AND TEST OPTIMIZATION**

**Lotus** is a systems programming language with deliberate design choices: module imports are string-based (`use "io";`), returns use `ret`, and declarations default to type-first bindings (`int n = 42;`). The compiler uses **LLVM** as its default backend for cross-platform compilation and advanced optimizations.
**Current Status:** Lotus is in active development with core features operational and advanced features being refined. See the [Language Idioms](#language-idioms-and-best-practices) section for production-ready patterns and the [Examples](#examples) section for current capabilities.
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://go.dev)
[![LLVM](https://img.shields.io/badge/LLVM-Backend-orange.svg)](https://llvm.org)

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Interactive REPL](#interactive-repl)
- [Installation](#installation)
- [Language Overview](#language-overview)
- [Standard Library](#standard-library)
- [Language Idioms and Best Practices](#language-idioms-and-best-practices)
- [Compiler Options](#compiler-options)
- [Documentation](#documentation)
- [Examples](#examples)
- [Contributing](#contributing)

## Features

### Language Features
- **Type-first bindings** - `int count = 42;` with explicit `ret` for returns
- **String-based imports** - `use "module";` with clear module boundaries
- **Pattern matching** - Comprehensive `match` with literals, ranges, bindings, guards, and wildcards
- **Generics/Templates** - `template<typename T>` with automatic type inference
- **Virtual functions** - `vrt fn` for virtual methods, `override fn` for overrides
- **Scope modifiers** - `static`, `lcl` (local), `gbl` (global) for explicit storage control
- **Basic structs and enums** - snake_case identifiers with fundamental OOP support
- **Error handling** - `try`/`catch`/`finally` and `throw` for exceptions, unwinding across function calls
- **Optional types** - `Some(value)` and `None` syntax (parser support, advanced operations in development)
- **Functional programming** - Pipe operator (`|>`), wrappers/decorators (`@`), currying (`partial`)
- **Void return type** - `fn void foo()` for functions with no return value
- **Type reinterpretation** - `bitcast<Type>(expr)` for low-level bit manipulation

**Features in Active Development:**
- Advanced struct operations (parameters, return types, complex nesting)
- Enhanced namespace and module system  
- Sophisticated template specialization and constraints
- Complete optional type integration with pattern matching

### Compiler Features
- **LLVM backend** - Cross-platform: x86, ARM, RISC-V, WebAssembly
- **Interactive REPL** - Run `lotus` with no arguments for a GHCi-style session (`lts ~` prompt)
- **Tree-walking interpreter** - Execute `.lts` files without compilation via `--interpret`
- **Advanced optimizations** - Dead code elimination, constant folding, inlining
- **Enhanced diagnostics** - Line/column tracking, "did you mean?" suggestions
- **Semantic analysis** - Unused variable detection, shadowing warnings
- **Language Server Protocol** - Real-time IDE support with `lotus --lsp`

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
    int first_val = 10;
    int second_val = 20;
    int max_val = max(first_val, second_val);
    
    printf("Maximum of %d and %d: %d\n", first_val, second_val, max_val);
    ret 0;
}
```

```bash
# Interactive REPL (no arguments)
./lotus

# Compile and run immediately
./lotus --run program.lts

# Run via interpreter (no compilation)
./lotus --interpret program.lts

# Compile to binary
./lotus program.lts -o myprogram
./myprogram

# Show compilation statistics
./lotus --stats program.lts
```

## Interactive REPL

Run `lotus` with no arguments to enter the interactive session. State (variables, functions) persists across inputs. Multi-line blocks are accumulated until the opening `{` is closed.

```text
$ lotus
Lotus 1.10.0 REPL  (interpreter)
Type :help for commands, :quit to exit.

lts ~ int x = 10;
lts ~ fn int double(int n) {
      |     return n * 2;
      | }
lts ~ println(double(x));
20
lts ~ :quit
Goodbye!
```

### REPL Commands

| Command | Description |
| --- | --- |
| `:q`, `:quit` | Exit the REPL |
| `:h`, `:help` | Show command reference |
| `:reset` | Clear all bindings and restart the interpreter |
| `:load <file>` | Load and evaluate a `.lts` source file |
| `:type <expr>` | Print the type of an expression |
| `:! <cmd>` | Run a shell command directly from the REPL |

```text
lts ~ :load examples/control_flow_if.lts
Loaded examples/control_flow_if.lts
lts ~ :! ls examples/*.lts | wc -l
24
lts ~ :type x
int
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
- **Clang** - For linking final binaries

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

### Optional Types

```lotus
// Optional types provide type-safe null handling
// Syntax: Some(value) for values that exist, None for absence

// Creating optionals
Some(42)        // Optional with integer value
Some("hello")   // Optional with string value
None            // Empty optional

// Optionals are represented as structs: { i1 has_value, T value }
```

### Pattern Matching

```lotus
use "io"

fn void gradeStudent(int score) {
    match score {
        case 90..100 => println("Grade: A"),
        case 80..89 => println("Grade: B"),
        case 70..79 => println("Grade: C"),
        case 60..69 => println("Grade: D"),
        default => println("Grade: F")
    }
}

fn void checkNumber(int num) {
    match num {
        case x when x > 0 => println("Positive"),
        case x when x < 0 => println("Negative"),
        case _ => println("Zero")
    }
}
```

**Pattern Types:**
- Literal: `case 42 =>`, `case "hello" =>`
- Range: `case 1..10 =>`
- Binding: `case x =>` (binds to variable)
- Wildcard: `case _ =>` (matches all)
- Guards: `case x when x > 10 =>`

### Generics/Templates

```lotus
use "io"

// Generic function with type inference
template<typename T>
fn T maximum(T a, T b) {
    if a > b {
        ret a;
    }
    ret b;
}

fn void main() {
    // Type is inferred from arguments
    int maxInt = maximum(10, 20);          // T = int
    float maxFloat = maximum(3.14, 2.71);  // T = float
    
    printf("Max int: %d\n", maxInt);
    printf("Max float: %f\n", maxFloat);
}
```

### Error Handling

```lotus
use "io";

fn int risky(int x) {
    if x < 0 {
        throw "negative input";
    }
    ret x * 2;
}

fn int main() {
    try {
        int r = risky(-1);
        println(r);
    } catch (string e) {
        println("caught:");
        println(e);
    } finally {
        println("cleanup runs either way");
    }
    ret 0;
}
```

- `try { ... } catch (...) { ... } finally { ... }` - at least one of `catch`/`finally` is required.
- `throw expr;` can appear anywhere, including inside a function called from the `try` block - it unwinds straight to the nearest enclosing `try`, wherever that is on the call stack.
- `catch { ... }` (no parentheses) catches without binding a variable; `catch (e) { ... }` binds the exception as `int` by default; `catch (TYPE e) { ... }` reinterprets the thrown value as `TYPE` (one of `int`, `string`, `bool`, `float`, `char`) - this is a static reinterpretation like `bitcast<Type>(...)`, not a runtime type check, since Lotus values carry no runtime type tag. For that reason **only one `catch` clause per `try` is allowed** - a second clause could never actually be selected over the first.
- `finally` always runs - on normal completion, when the exception is caught, when it isn't (after which it re-propagates to the next enclosing `try`), and even if the `try`/`catch` body exits early via `ret`/`break`/`continue`.
- An uncaught `throw` prints an error and exits the program with status 1.

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

## Language Idioms and Best Practices

This section covers idiomatic Lotus patterns and best practices. For detailed coding style and conventions, see the [STYLE_GUIDE.md](Important_Documentation/STYLE_GUIDE.md).

### 1. Function Return Patterns

**Always use explicit `ret` statements:**
```lotus
// ✅ Good: Explicit return
fn int calculate_sum(int a, int b) {
    int result = a + b;
    ret result;
}

// ✅ Good: Early return for error conditions
fn int safe_divide(int a, int b) {
    if b == 0 {
        ret -1;  // Error sentinel
    }
    ret a / b;
}
```

**Prefer void functions when no value is returned:**
```lotus
// ✅ Good: Clear that nothing is returned
fn void print_banner() {
    printf("=== Lotus Application ===\n");
}

// ✅ Good: Explicit void return type
fn void main() {
    print_banner();
}
```

### 2. Type-First Variable Declarations

**Leverage Lotus's type-first syntax:**
```lotus
// ✅ Good: Clear type declarations
fn void process_data() {
    int count = 0;
    string message = "Processing";
    bool is_complete = false;
    float progress = 0.0;
    
    // Logic here...
}
```

**Use meaningful variable names with types:**
```lotus
// ✅ Good: Descriptive names
fn void file_operations() {
    int file_size = 1024;
    string file_path = "/tmp/data.txt";
    bool is_readable = true;
}
```

### 3. Pattern Matching Idioms

**Use ranges for numerical classifications:**
```lotus
fn void classify_score(int score) {
    match score {
        case 90..100 => println("Excellent"),
        case 80..89 => println("Good"),
        case 70..79 => println("Average"),
        case 60..69 => println("Below Average"),
        default => println("Needs Improvement")
    }
}
```

**Use guards for complex conditions:**
```lotus
fn void analyze_number(int num) {
    match num {
        case x when x > 0 && x % 2 == 0 => println("Positive Even"),
        case x when x > 0 && x % 2 == 1 => println("Positive Odd"),
        case x when x < 0 => println("Negative"),
        default => println("Zero")
    }
}
```

**Prefer binding when you need the matched value:**
```lotus
fn void process_value(int input) {
    match input {
        case val when val > 100 => {
            printf("Large value: %d\n", val);
            // Process large value
        },
        case val => {
            printf("Normal value: %d\n", val);
            // Process normal value
        }
    }
}
```

### 4. Template/Generic Function Patterns

**Use descriptive type parameter names:**
```lotus
// ✅ Good: Clear what T represents
template<typename T>
fn T find_maximum(T first, T second) {
    if first > second {
        ret first;
    }
    ret second;
}
```

**Keep generic functions simple and focused:**
```lotus
template<typename T>
fn void swap(T* a, T* b) {
    T temp = *a;
    *a = *b;
    *b = temp;
}

// Usage with type inference
int x = 10, y = 20;
swap(&x, &y);  // T inferred as int
```

### 5. Error Handling Patterns

**Use sentinel values consistently:**
```lotus
// Pattern: Return -1 for errors in integer functions
fn int parse_number(string str) {
    // Parsing logic...
    if parsing_failed {
        ret -1;  // Consistent error sentinel
    }
    ret result;
}

// Pattern: Return null for pointer functions
fn int* allocate_array(int size) {
    if size <= 0 {
        ret null;  // Invalid input
    }
    
    int* array = malloc(sizeof(int) * size);
    if array == null {
        ret null;  // Allocation failed
    }
    
    ret array;
}
```

**Always check return values:**
```lotus
fn void safe_operations() {
    int* buffer = allocate_array(100);
    if buffer == null {
        printf("Error: Failed to allocate memory\n");
        ret;
    }
    
    // Use buffer...
    
    free(buffer);  // Always clean up
}
```

### 6. Module Import Best Practices

**Use specific imports to show dependencies:**
```lotus
use "io";     // For printf, println
use "mem";    // For malloc, free
use "math";   // For mathematical operations

fn int main() {
    printf("Starting application...\n");
    ret 0;
}
```

**Group related functionality:**
```lotus
// For string processing applications
use "io";
use "str";
use "mem";

fn void process_text(string input) {
    int length = len(input);
    string upper = toUpper(input);
    printf("Original: %s (length: %d)\n", input, length);
    printf("Uppercase: %s\n", upper);
}
```

### 7. Control Flow Idioms

**Prefer early returns for validation:**
```lotus
fn int validate_and_process(int* data, int size) {
    if data == null {
        ret -1;  // Invalid data
    }
    
    if size <= 0 {
        ret -2;  // Invalid size
    }
    
    // Main processing logic here
    ret process_data(data, size);
}
```

**Use for loops for known iterations:**
```lotus
fn void initialize_array(int* arr, int size) {
    for int i = 0; i < size; i++ {
        arr[i] = 0;
    }
}
```

**Use while loops for condition-based iteration:**
```lotus
fn int find_first_zero(int* arr, int size) {
    int index = 0;
    while index < size && arr[index] != 0 {
        index++;
    }
    ret index < size ? index : -1;
}
```

### 8. Memory Management Patterns

**Follow the allocate-use-free pattern:**
```lotus
fn void process_large_data() {
    int* buffer = malloc(sizeof(int) * 1000);
    if buffer == null {
        printf("Memory allocation failed\n");
        ret;
    }
    
    // Initialize buffer
    for int i = 0; i < 1000; i++ {
        buffer[i] = i * 2;
    }
    
    // Process data...
    
    free(buffer);  // Always free allocated memory
}
```

**Use RAII-style patterns where possible:**
```lotus
fn int process_with_cleanup(int size) {
    int* temp = malloc(sizeof(int) * size);
    if temp == null {
        ret -1;  // Early return, nothing to clean up
    }
    
    int result = do_processing(temp);
    free(temp);  // Always reached if malloc succeeded
    ret result;
}
```

### 9. Printf Format String Best Practices

**Use appropriate format specifiers:**
```lotus
fn void display_info(int count, string name, float ratio) {
    printf("Item: %s\n", name);           // %s for strings
    printf("Count: %d\n", count);         // %d for integers
    printf("Ratio: %.2f\n", ratio);       // %f for floats (with precision)
}
```

**Keep format strings literal when possible:**
```lotus
// ✅ Good: Literal format strings are easier to validate
printf("Processing %d items...\n", count);
printf("Error code: %d\n", error);

// ⚠️ Avoid: Dynamic format strings when not necessary
// printf(format_string, value);  // Harder to validate
```

### 10. Code Organization Idioms

**Structure functions logically:**
```lotus
// 1. Helper/utility functions first
fn bool is_valid_size(int size) {
    ret size > 0 && size < 10000;
}

// 2. Core logic functions
fn int* create_buffer(int size) {
    if !is_valid_size(size) {
        ret null;
    }
    ret malloc(sizeof(int) * size);
}

// 3. Main/entry point last
fn int main() {
    int* buffer = create_buffer(100);
    if buffer != null {
        free(buffer);
    }
    ret 0;
}
```

**Use meaningful function names:**
```lotus
// ✅ Good: Function names describe purpose
fn bool validate_user_input(string input);
fn int calculate_fibonacci_number(int n);
fn void cleanup_temporary_files();

// ❌ Avoid: Generic or unclear names
fn bool check(string s);
fn int calc(int n);
fn void clean();
```

## Compiler Options

```
Usage: lotus [options] <source.lts>

Output Options:
  -o <file>        Output file name (default: a.out)
  --emit-llvm      Emit LLVM IR instead of binary
  --run            Compile and run immediately
  --interpret      Run source file via tree-walking interpreter (no compilation)

Backend Options:
  --target <triple> Cross-compile for target (e.g., aarch64-linux-gnu)

Optimization:
  -O0              No optimization
  -O1              Basic optimization  
  -O2              Standard optimization (default)
  -O3              Aggressive optimization

Analysis & Debugging:
  --ast-dump       Print AST structure and exit
  --stats          Show compilation statistics
  --timing         Show detailed phase timing
  --token-dump     Print tokens and exit
  -v, --verbose    Verbose compilation output
  -q, --quiet      Suppress non-error output

Warning Control:
  -Wall            Enable all warnings
  -Werror          Treat warnings as errors
  -Wunused         Warn about unused variables
  -Wshadow         Warn about variable shadowing
  -w               Suppress all warnings

Development:
  --lsp            Run as Language Server Protocol server
  --docs [section] Show offline documentation
  -I <dir>         Add include directory for imports

Other:
  --version        Show version information
  --help           Show this help message
```

## Documentation

| Document | Description |
|----------|-------------|
| [STYLE_GUIDE.md](Important_Documentation/STYLE_GUIDE.md) | Naming conventions, formatting, code organization |
| [Language Idioms (this document)](#language-idioms-and-best-practices) | Idiomatic patterns and best practices |
| [STDLIB_AND_IMPORTS.md](Important_Documentation/STDLIB_AND_IMPORTS.md) | Import patterns and module usage |
| [STDLIB_FINAL_SUMMARY.md](Important_Documentation/STDLIB_FINAL_SUMMARY.md) | Complete stdlib reference |
| [DEVELOPMENT.md](Important_Documentation/DEVELOPMENT.md) | Contributor guide and architecture |
| [RELEASE_PROCESS.md](Important_Documentation/RELEASE_PROCESS.md) | Release workflow and version management |
| [AUTOMATED_RELEASE.md](Important_Documentation/AUTOMATED_RELEASE.md) | Fully automated release system |

## Examples

The language showcases several powerful features through practical examples:

**Basic Programs:**
- [control_flow_if.lts](examples/control_flow_if.lts) - Conditional logic and function calls
- [control_flow_for.lts](examples/control_flow_for.lts) - Loop iteration patterns  
- [control_flow_while.lts](examples/control_flow_while.lts) - Condition-based loops

**Advanced Features:**
- [pattern_matching.lts](examples/pattern_matching.lts) - Complete pattern matching with literals, ranges, guards, and bindings
- [template_comprehensive_test.lts](examples/template_comprehensive_test.lts) - Generic functions with type inference
- [optional_test.lts](examples/optional_test.lts) - Optional type syntax (`Some`/`None`)
- [error_handling.lts](examples/error_handling.lts) - `try`/`catch`/`finally`/`throw`, including unwinding across a function call

**Language Features in Development:**
- Complex struct operations with parameters and return types
- Advanced namespace and module systems  
- Sophisticated template specialization
- Enhanced functional programming constructs

For complete working examples demonstrating idioms and best practices, see the [Language Idioms and Best Practices](#language-idioms-and-best-practices) section above.

Additional test files and examples can be found in [tests/](tests/) directory.

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
