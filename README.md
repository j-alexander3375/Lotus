# Lotus Compiler - Modern Modular Architecture

## Overview
Lotus is a powerful systems programming language compiler that generates optimized x86-64 assembly. Built with a modular architecture, Lotus supports advanced features including data structures, bitwise operations, logical operators, and comprehensive type systems.

## Project Structure

### Core Modules
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
├── printfuncs.go        # Print function implementations
│                         # - printf, println, fprintf, sprintf
│                         # - fatalf, fatalln, logf, logln
│                         # Uses write() syscall (fd=1 for stdout)
│
├── arithmetic.go        # Arithmetic & bitwise operations
│                         # - BinaryOp: +, -, *, /, %
│                         # - BitwiseOp: &, |, ^, <<, >>
│                         # - UnaryOp: +, -, !, ~
│                         # - IncrementOp: ++, -- (prefix/postfix)
│
├── control_flow.go      # Control flow & logical operations
│                         # - IfStatement, WhileLoop, ForLoop
│                         # - Comparison: ==, !=, <, <=, >, >=
│                         # - LogicalOp: &&, || (short-circuit)
│                         # - TernaryOp: ? :
│
├── references.go        # Memory references & assignments
│                         # - Reference (&), Dereference (*)
│                         # - Assignment (=)
│                         # - CompoundAssignment: +=, -=, *=, /=, %=
│
├── functions.go         # User-defined functions
│                         # - Function definitions with parameters
│                         # - System V AMD64 ABI compliance
│                         # - Recursive calls supported
│
├── memory.go            # Memory allocation
│                         # - malloc(), free(), sizeof()
│                         # - Type-aware size calculation
│                         # - Heap management via libc
│
├── array.go             # Dynamic arrays
│                         # - Array literals: [1, 2, 3]
│                         # - Array indexing: arr[index]
│                         # - Heap-allocated arrays
│
├── struct.go            # Struct types
│                         # - Field definitions with offsets
│                         # - Struct literals & initialization
│                         # - Field access: . and ->
│
├── enum.go              # Enumeration types
│                         # - Enum definitions with auto-values
│                         # - Compile-time constant resolution
│
└── class.go             # Object-oriented programming
                          # - Class definitions with fields/methods
                          # - Method calls with 'this' pointer
                          # - new operator for instantiation
```

## Language Features

### Type System
**Primitive Types:**
- `int`, `int8`, `int16`, `int32`, `int64` - Signed integers (1-8 bytes)
- `uint`, `uint8`, `uint16`, `uint32`, `uint64` - Unsigned integers
- `string` - UTF-8 strings (pointer type)
- `bool` - Boolean values (true/false)
- `float` - Floating-point numbers

**Constants:**
```lotus
const int MAX_VALUE = 1000;
const string APP_NAME = "Lotus";
const bool DEBUG = true;
```
Constants are immutable, stored in the data section, and accessed via labels.

**Complex Types:**
- Arrays: `int[]`, dynamic sizing
- Structs: User-defined composite types
- Enums: Named integer constants
- Classes: OOP with methods and inheritance
- Pointers: `int*`, `string*`, etc.

### Print Functions (printfuncs.go)
All print functions use Linux `write()` syscall:
```lotus
printf("Hello, %s\n", name);      // Formatted output
println("Auto newline");           // Output with newline
sprintf(buffer, "Format", args);   // String formatting
fatalf("Error: %s\n", msg);        // Print error and exit(1)
logf("Debug: %d\n", value);        // Logging output
```

### Arithmetic & Bitwise Operations (arithmetic.go)
```lotus
// Arithmetic
int sum = 10 + 5;
int diff = x - y;
int prod = a * b;
int quot = x / 2;
int mod = x % 3;

// Bitwise
int masked = flags & 0xFF;
int combined = a | b;
int toggled = x ^ mask;
int shifted = value << 2;
int unshifted = value >> 1;
int inverted = ~bits;

// Increment/Decrement
int x = 5;
x++;              // Postfix increment
++x;              // Prefix increment
y = x++;          // Use then increment
z = ++x;          // Increment then use
```

### Logical & Comparison Operators (control_flow.go)
```lotus
// Comparisons
if (x == 10) { }
if (x != 0) { }
if (x < 100) { }
if (x >= 50) { }

// Logical operations (short-circuit)
if (x > 0 && y < 10) { }
if (flag1 || flag2) { }

// Ternary operator
int max = (a > b) ? a : b;
string status = (active) ? "on" : "off";
```

### Constants
```lotus
// Declare constants (immutable values)
const int MAX_BUFFER = 4096;
const string APP_VERSION = "1.0.0";
const bool DEBUG = true;

// Use constants in expressions
int buffer_size = MAX_BUFFER;
if (DEBUG) {
  println("Debug mode enabled");
}

// Constants are stored in data section for efficiency
// They cannot be modified after declaration
```

### Control Flow (control_flow.go)
```lotus
// If-else
if (x > 10) {
  printf("Greater\n");
} else {
  printf("Lesser\n");
}

// While loops
while (x > 0) {
  x = x - 1;
}

// For loops
for (int i = 0; i < 10; i++) {
  printf("%d\n", i);
}
```

### Compound Assignment (references.go)
```lotus
x += 10;          // x = x + 10
y -= 5;           // y = y - 5
z *= 2;           // z = z * 2
w /= 3;           // w = w / 3
r %= 7;           // r = r % 7
```

### Memory Management (memory.go)
```lotus
// Allocate memory
int* ptr = malloc(sizeof(int) * 10);

// Use memory
ptr[0] = 42;
int value = ptr[5];

// Free memory
free(ptr);

// Get type sizes
int size = sizeof(int64);  // Returns 8
```

### Arrays (array.go)
```lotus
// Array literals
int[] nums = [1, 2, 3, 4, 5];

// Array indexing
int first = nums[0];
nums[2] = 100;

// Dynamic arrays
int[] dynamic = malloc(sizeof(int) * count);
```

### Structs (struct.go)
```lotus
struct Point {
  int x;
  int y;
}

// Struct initialization
Point p = {10, 20};

// Field access
int x_coord = p.x;
p.y = 30;

// Pointer access
Point* ptr = &p;
int val = ptr->x;
```

### Enums (enum.go)
```lotus
enum Color {
  Red,      // 0
  Green,    // 1
  Blue      // 2
}

enum Status {
  Ok = 0,
  Error = -1,
  Pending = 1
}

int color = Color.Red;
```

### Classes (class.go)
```lotus
class Person {
  int age;
  string name;
  
  fn init(int a, string n) {
    this.age = a;
    this.name = n;
  }
  
  fn greet() {
    printf("Hello, I'm %s\n", this.name);
  }
}

// Instantiation
Person* p = new Person;
p->init(25, "Alice");
p->greet();
```

### User-Defined Functions (functions.go)
```lotus
fn add(int a, int b) -> int {
  return a + b;
}

fn factorial(int n) -> int {
  if (n <= 1) {
    return 1;
  }
  return n * factorial(n - 1);
}

int result = add(5, 3);
int fact = factorial(5);
```

Features:
- Multiple parameters (System V AMD64 ABI)
- Return types with type checking
- Local variable scoping
- Recursive calls fully supported
- Up to 6 register-passed arguments (rdi, rsi, rdx, rcx, r8, r9)

### Assembly Target
- **Architecture**: x86-64 (AMD64)
- **Syntax**: GNU AT&T assembly
- **ABI**: System V AMD64 calling convention
- **Stack**: 16-byte aligned for syscalls
- **Optimization**: Direct instruction generation, no IR

### Example: Bitwise Operations
```lotus
int flags = 0x0F;
int masked = flags & 0xFF;
int shifted = flags << 4;
```
Generates:
```asm
    movq $15, %rax
    movq %rax, -8(%rbp)
    movq -8(%rbp), %rax
    pushq %rax
    movq $255, %rax
    movq %rax, %rcx
    popq %rax
    andq %rcx, %rax
    movq %rax, -16(%rbp)
```

### Example: Logical Short-Circuit
```lotus
if (x > 0 && y < 10) {
  printf("valid\n");
}
```
Generates:
```asm
    movq -8(%rbp), %rax
    testq %rax, %rax
    jz .and_end_0
    movq -16(%rbp), %rax
    # ... comparison code
.and_end_0:
```

### Example: Ternary Operator
```lotus
int max = (a > b) ? a : b;
```
Generates:
```asm
    movq -8(%rbp), %rax
    cmpq -16(%rbp), %rax
    setg %al
    testq %rax, %rax
    jz .ternary_false_0
    movq -8(%rbp), %rax
    jmp .ternary_end_0
.ternary_false_0:
    movq -16(%rbp), %rax
.ternary_end_0:
    movq %rax, -24(%rbp)
```

## Compilation Pipeline

```
Source Code (.lts)
       ↓
[Tokenizer]  → Token Stream (70+ token types)
       ↓        - Multi-character operators
       ↓        - Keyword recognition
       ↓
[Parser]     → Abstract Syntax Tree
       ↓        - Recursive descent
       ↓        - Operator precedence
       ↓
[CodeGen]    → x86-64 Assembly
  ├─ arithmetic.go   (expressions, bitwise, increment)
  ├─ control_flow.go (branches, loops, logical ops)
  ├─ references.go   (memory, compound assignment)
  ├─ functions.go    (calls, definitions, ABI)
  ├─ printfuncs.go   (syscall writes)
  ├─ memory.go       (malloc/free)
  ├─ array.go        (indexing, literals)
  ├─ struct.go       (field access)
  ├─ enum.go         (constant resolution)
  └─ class.go        (methods, this pointer)
       ↓
[GCC/AS]     → Binary Executable (ELF)
       ↓
Output (./a.out)
```

## Build & Run

```bash
# Build compiler
go build -o lotus src/*.go

# Compile Lotus source
./lotus program.lts

# Run compiled binary
./a.out

# View assembly only
./lotus -S program.lts

# Save assembly to file
./lotus -S -o output.s program.lts

# Verbose output
./lotus -v program.lts

# Show version
./lotus -v

# Help
./lotus -h
```

## Advanced Features

### Data Structures
```lotus
// Arrays
int[] numbers = [1, 2, 3, 4, 5];
numbers[2] = 100;

// Structs
struct Vec2 { int x; int y; }
Vec2 pos = {10, 20};

// Enums
enum State { Idle, Running, Stopped }
State current = State.Running;

// Classes
class Entity {
  int health;
  fn takeDamage(int dmg) {
    this.health -= dmg;
  }
}
```

### Memory Management
```lotus
// Allocate
int* buffer = malloc(1024);

// Use
for (int i = 0; i < 256; i++) {
  buffer[i] = i * 2;
}

// Free
free(buffer);
```

### Complex Expressions
```lotus
int result = (x > 0) ? (y << 2) | flags : -1;
value += (count++ * 3) & mask;
bool valid = (status == Ok) && (data != 0) || fallback;
```

## Future Goals

### Short-term (Next Iteration)
- [ ] Generic types support
- [ ] String manipulation library
- [ ] Standard library expansion
- [ ] Module system / imports
- [ ] Pattern matching

### Mid-term
- [ ] Traits/interfaces system
- [ ] Closures and lambdas
- [ ] Compile-time evaluation
- [ ] Inline assembly blocks
- [ ] LLVM backend option

### Long-term
- [ ] Self-hosting compiler
- [ ] Package manager integration
- [ ] IDE language server
- [ ] Cross-compilation support
- [ ] Optimization passes

## Design Principles

1. **Modularity** - Each feature in its own file with clear boundaries
2. **Extensibility** - Easy to add new operators, types, and constructs
3. **Performance** - Direct assembly generation, minimal overhead
4. **Safety** - Type checking, bounds checking (future)
5. **Simplicity** - Clean syntax, predictable behavior
6. **System Programming** - Direct syscall access, memory control

## Performance Characteristics

- **Compilation**: Single-pass, sub-second for small programs
- **Runtime**: Direct syscalls, no GC overhead
- **Binary Size**: Minimal (no runtime, static linking)
- **Optimization**: Currently unoptimized (future: peephole, register allocation)

## Extension Guide

### Adding New AST Nodes
1. Define node struct in appropriate module
2. Implement `astNode()` method
3. Add generation function `generateXXX()`
4. Update `generateExpressionToReg()` or `generateStatement()`
5. Add parser support

### Adding New Operators
1. Add `TokenXXX` constant to `keywords.go`
2. Update tokenizer recognition in `tokenizer.go`
3. Add `TokenValue` case for pretty printing
4. Create or update AST node type
5. Implement code generation
6. Update parser precedence if needed

### Adding New Types
1. Add `TokenTypeXXX` to `keywords.go`
2. Update tokenizer keyword recognition
3. Add to parser type checking
4. Update `getTypeSizeFromToken()` in `memory.go`
5. Add type-specific code generation

## Community & Contributing

Lotus is designed for extensibility. Contributions welcome in:
- New language features
- Optimization passes
- Standard library functions
- Documentation and examples
- Bug fixes and testing

## License

MIT License - See LICENSE file for details

