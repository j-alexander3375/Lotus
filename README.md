# Lotus Compiler - Modular Architecture

## Overview
The Lotus compiler has been restructured into a highly modular architecture with separate modules for different language features. This design enables extensibility and maintainability while keeping code concerns separated.

## Project Structure

### Core Modules
```
src/
├── main.go              # Entry point and CLI
├── compiler.go          # Compilation pipeline
├── flags.go             # Command-line argument parsing
├── keywords.go          # Token type definitions
├── tokenizer.go         # Lexical analysis
├── parser.go            # Syntactic analysis (AST building)
├── codegen.go           # Code generation orchestrator
│
├── printfuncs.go        # Print function implementations
│                         # - printf, println, fprintf
│                         # - fatalf, fatalln
│                         # - logf, logln
│                         # Uses write() syscall (fd=1 for stdout)
│
├── arithmetic.go        # Arithmetic operations
│                         # - BinaryOp: +, -, *, /, %
│                         # - UnaryOp: +, -, !
│
├── control_flow.go      # Control flow structures
│                         # - IfStatement
│                         # - WhileLoop
│                         # - ForLoop
│                         # - Comparison operators
│
├── references.go        # Memory references
│                         # - Reference (&)
│                         # - Dereference (*)
│                         # - Assignment
│                         # - ArrayAccess
│
└── functions.go         # User-defined functions
                          # - FunctionDefinition
                          # - Function calls
                          # - System V ABI compliance
```

## Language Features

### Supported Types
- `int` - 64-bit integers
- `string` - UTF-8 strings
- `bool` - Boolean values (true/false)
- `float` - Floating-point numbers (fixed-point: *1000)

### Print Functions (printfuncs.go)
All print functions use Linux `write()` syscall:
```lotus
printf("Hello, %s\n", name);      // Formatted output to stdout
println("Line with newline");      // Output with automatic newline
fprintf(fd, "Format", args...);    // Write to file descriptor
fatalf("Error: %s\n", msg);        // Print error and exit(1)
logf("Log: %s\n", data);           // Logging output
```

### Arithmetic Operations (arithmetic.go)
```lotus
int x = 10 + 5;      // Addition
int y = x - 3;       // Subtraction
int z = x * 2;       // Multiplication
int a = x / 2;       // Division
int b = x % 3;       // Modulo
int c = -x;          // Negation
```

### Control Flow (control_flow.go)
```lotus
if (x > 10) {
  printf("x is greater than 10\n");
} else {
  printf("x is 10 or less\n");
}

while (x > 0) {
  x = x - 1;
}

for (int i = 0; i < 10; i = i + 1) {
  printf("i = %d\n", i);
}
```

### Comparisons
- `==` - Equal
- `!=` - Not equal
- `<` - Less than
- `<=` - Less than or equal
- `>` - Greater than
- `>=` - Greater than or equal

### References (references.go)
```lotus
int x = 42;
int* ptr = &x;        // Take address of x
int y = *ptr;         // Dereference pointer
x = 100;              // Assignment/reassignment
```

### User-Defined Functions (functions.go)
```lotus
fn add(int a, int b) -> int {
  ret a + b;
}

fn greet(string name) {
  printf("Hello, ");
  printf(name);
  printf("\n");
}

int result = add(5, 3);
greet("World");
```

Features:
- Multiple parameters (System V AMD64 ABI)
- Return types
- Local variable scoping
- Recursive calls supported

## Code Generation

### Assembly Target
- x86-64 GNU AT&T syntax
- Direct syscall usage (no libc dependency for core functions)
- Proper stack frame management
- 16-byte stack alignment for syscalls

### Example: Arithmetic Expression
```lotus
int sum = 10 + 5;
```
Generates:
```asm
    movq $10, %rax
    pushq %rax
    movq $5, %rax
    movq %rax, %rcx
    popq %rax
    addq %rcx, %rax
    movq %rax, -8(%rbp)
```

### Example: Control Flow
```lotus
if (x > 0) {
  printf("positive");
}
```
Generates:
```asm
    cmpq $0, %rax
    setg %al
    movzbl %al, %eax
    testq %rax, %rax
    jz .if_end_0
.if_then_0:
    # printf code
    jmp .if_end_0
.if_end_0:
```

## Compilation Pipeline

```
Source Code (.lts)
       ↓
[Tokenizer]  → Token Stream (lexical analysis)
       ↓
[Parser]     → Abstract Syntax Tree (syntactic analysis)
       ↓
[CodeGen]    → Assembly instructions
  ├─ arithmetic.go   (expression evaluation)
  ├─ control_flow.go (jumps and labels)
  ├─ references.go   (memory operations)
  ├─ functions.go    (function calls)
  └─ printfuncs.go   (syscall writes)
       ↓
[GCC]        → Binary executable
       ↓
Output (./a.out)
```

## Build & Test

```bash
# Build compiler
cd src
go build -o ../lotus

# Compile Lotus source
../lotus program.lts

# Run compiled binary
./a.out

# View assembly
../lotus -S program.lts

# Debug with token dump
../lotus -td program.lts
```

## Future Extensions

Easy to add:
1. **New AST nodes** - Create node types in appropriate module, add to parser, add codegen
2. **New operators** - Add to arithmetic/control_flow, update tokenizer/keywords
3. **New print functions** - Add generator to printfuncs.go RegisteredPrintFunctions
4. **Data structures** - Create arrays.go, structs.go modules
5. **Memory** - Create memory.go for heap allocation

## Design Principles

1. **Separation of Concerns** - Each module handles one language feature
2. **Extensibility** - Easy to add new features without modifying existing code
3. **Modularity** - Independent compilation units that can be tested separately
4. **Assembly-First** - Direct generation of efficient x86-64 code
5. **System ABI Compliance** - Follows System V AMD64 calling conventions

## Performance

- Direct syscalls (no libc overhead for core functions)
- x86-64 optimized code generation
- Single-pass compilation
- No intermediate representations beyond AST
