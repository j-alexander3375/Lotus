# Lotus Standard Library and Import System

## Overview

Lotus now features a modular standard library with Rust-inspired import syntax. This allows organized, namespace-aware access to standard library functions across organized modules.

## Import Syntax

The `use` keyword imports modules and their functions. Syntax similar to Rust provides flexibility in what gets imported.

### Basic Import - Import All Functions

```lotus
use "io";

fn main() {
    println("Hello, World!");
}
```

Imports all functions from the `io` module into the current namespace.

### Specific Function Import

```lotus
use "io::printf";

fn main() {
    printf("Value: %d\n", 42);
}
```

Imports only the `printf` function from `io` module.

### Wildcard Import

```lotus
use "io::*";

fn main() {
    println("Using wildcard");
    printf("Also works: %d\n", 100);
}
```

Explicitly imports all functions (equivalent to basic import).

### Aliased Imports

```lotus
use "io" as output;

fn main() {
    output::println("Hello");
}
```

Imports module under an alias name for disambiguation.

### Multiple Imports

```lotus
use "io";
use "mem";
use "math";

fn main() {
    x: int = malloc(sizeof(int));
    println("Allocated memory");
    free(x);
}
```

Multiple import statements can load different modules.

## Standard Library Modules

### `io` - Input/Output Functions

Functions for formatted and unformatted output:

- **`print(args...)`** - Output to stdout (variadic)
- **`println(args...)`** - Output with newline (variadic)
- **`printf(format, args...)`** - Formatted output (variadic)
- **`fprintf(fd, format, args...)`** - Formatted output to file descriptor
- **`sprint(args...)`** - Format to string (variadic)
- **`sprintf(format, args...)`** - Format to string with format string
- **`sprintln(args...)`** - Format to string with newline (variadic)

Supported printf verbs: %%, %d, %b, %o, %x/%X, %c (byte), %q (quoted strings), %s, %v.

Example:
```lotus
use "io";

fn main() {
    printf("Number: %d\n", 42);
    println("Simple output");
}
```

### `mem` - Memory Management

Functions for dynamic memory allocation:

- **`malloc(size)`** - Allocate memory, returns pointer
- **`free(ptr)`** - Deallocate memory
- **`sizeof(type)`** - Get size of type in bytes

Example:
```lotus
use "mem";

fn main() {
    size: int = sizeof(int);
    ptr: int = malloc(size);
    free(ptr);
}
```

### `math` - Mathematical Functions

Numerical operations and calculations:

- **`abs(n)`** - Absolute value
- **`min(a, b)`** - Minimum of two values
- **`max(a, b)`** - Maximum of two values
- **`sqrt(x)`** - Square root
- **`pow(base, exp)`** - Power (base^exp)

Status: abs/min/max implemented; sqrt/pow pending.

Example:
```lotus
use "math";

fn main() {
    a: int = 10;
    b: int = 5;
    larger: int = max(a, b);
    printf("Max: %d\n", larger);
}
```

### `str` - String Manipulation

String operations and analysis:

- **`len(str)`** - Get string length
- **`concat(s1, s2, ...)`** - Concatenate strings (variadic)
- **`compare(s1, s2)`** - Compare strings (returns 0 if equal)
- **`copy(src, dst)`** - Copy string from source to destination

Status: len implemented; concat/compare/copy pending.

Example:
```lotus
use "str";

fn main() {
    s: string = "Hello";
    length: int = len(s);
    printf("Length: %d\n", length);
}
```

## Module Architecture

### System Structure

```
src/
├── stdlib.go          # Module definitions and registration
├── keywords.go        # use, as keywords
├── tokenizer.go       # Tokenization of use statements
├── parser.go          # Parsing of import statements
├── codegen.go         # Code generation with import tracking
│
└── printfuncs.go      # io module implementations
└── memory.go          # mem module implementations
```

### How Modules Work

1. **Definition**: Modules are defined in `stdlib.go` as `StdlibModule` structs
2. **Registration**: All modules register in the `StandardLibrary` map
3. **Parsing**: The `use` keyword is tokenized and parsed into `ImportStatement` AST nodes
4. **Loading**: Import statements populate the `ImportContext` with available functions
5. **Code Generation**: Function calls dispatch to registered code generators

### ImportContext

The `ImportContext` tracks what has been imported:

```go
type ImportContext struct {
    ImportedModules map[string]string          // alias -> module
    ImportedFunctions map[string]*StdlibFunction // function name -> implementation
}
```

When you import a function, it gets added to `ImportedFunctions` map for availability during codegen.

## Adding New Modules

To add a new stdlib module:

### 1. Create Module Functions

Implement code generation functions in a new file (e.g., `time.go`):

```go
func generateTimeNow(cg *CodeGenerator, args []ASTNode) {
    // Generate assembly for time::now()
}
```

### 2. Register in stdlib.go

Add module creation function:

```go
func createTimeModule() *StdlibModule {
    return &StdlibModule{
        Name: "time",
        Functions: map[string]*StdlibFunction{
            "now": {
                Name:    "now",
                Module:  "time",
                NumArgs: 0,
                CodeGen: generateTimeNow,
            },
        },
    }
}
```

### 3. Add to StandardLibrary Map

```go
var StandardLibrary = map[string]*StdlibModule{
    "io":   createIOModule(),
    "mem":  createMemoryModule(),
    "math": createMathModule(),
    "str":  createStringModule(),
    "time": createTimeModule(),  // New module
}
```

## Import Resolution

When parsing `use "module::function"`:

1. **Module lookup** - Find module in `StandardLibrary`
2. **Function lookup** - Find function in module's `Functions` map
3. **Registration** - Add to `ImportContext.ImportedFunctions`
4. **Code generation** - Call registered `CodeGen` function

Errors are reported if:
- Module doesn't exist
- Function doesn't exist in module
- Multiple imports create conflicts (future enhancement)

## Planned Enhancements

### Namespace Management
```lotus
use "io" as io_module;
use "math" as math_module;

fn main() {
    io_module::println("Output");
    math_module::sqrt(16);
}
```

### Re-exports
```lotus
use "io::{println, printf}";
```

### Conditional Imports (Future)
```lotus
#[cfg(feature = "debug")]
use "debug";
```

### Package System
Support for user-defined packages and modules:
```lotus
use "my_lib::utils::math";
```

## Compilation and Usage

### Building with Modules

```bash
cd src/
go build -o ../lotus
```

### Running Programs with Imports

```bash
./lotus -run program.lts
./lotus -run program.lts -S -o program.s  # Generate assembly
./lotus -run program.lts -o program       # Compile to binary
```

### Examples

```lotus
// example1.lts - Multiple module usage
use "io";
use "mem";
use "math";

fn main() {
    max_val: int = max(10, 20);
    printf("Maximum: %d\n", max_val);
    
    ptr: int = malloc(sizeof(int));
    println("Memory allocated");
    free(ptr);
}
```

```lotus
// example2.lts - String operations
use "str";
use "io";

fn main() {
    s: string = "Lotus";
    length: int = len(s);
    printf("String length: %d\n", length);
}
```

## Implementation Details

### TokenTypes Added
- `TokenUse` - use keyword
- `TokenAs` - as keyword

### AST Nodes Added
```go
type ImportStatement struct {
    Module     string    // Module name
    Items      []string  // Specific items, nil for all
    Alias      string    // Optional alias
    IsWildcard bool      // true for module::*
}
```

### CodeGenerator Changes
- Added `imports *ImportContext` field
- Added `generateImportStatement()` method
- Imports tracked during parsing, available for code generation

## Type Safety

Modules provide:
- **Compile-time checking** - Function availability verified at compile time
- **Namespace isolation** - No global function pollution
- **Clear dependencies** - Import statements show what's used
- **Error reporting** - Unknown module/function reported with context

## Performance Impact

- **Zero runtime overhead** - Imports resolved at compile time
- **No dynamic loading** - All module code is compiled in
- **Efficient dispatch** - Direct function pointers in codegen
- **Minimal binary size** - Only imported functions included in final binary

## Backward Compatibility

Existing Lotus code without import statements continues to work. The old function names (printf, println, malloc, etc.) remain available when no imports are specified, for compatibility.

Future versions may deprecate implicit global function availability in favor of explicit imports.
