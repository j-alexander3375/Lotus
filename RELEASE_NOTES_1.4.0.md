# Lotus 1.4.0 Release Notes

**Release Date:** January 2026

## Overview

Lotus 1.4.0 introduces powerful functional programming features including the pipe operator, wrappers/decorators, and partial function application. This release also adds a comprehensive random number generation module and significantly expands the string manipulation capabilities.

## New Features

### Functional Programming

#### Pipe Operator (`|>`)
Chain function calls in a readable left-to-right style:
```lotus
int result = 5 |> double |> addOne;  // Equivalent to addOne(double(5))
```

#### Wrappers/Decorators
Define reusable function wrappers and apply them with `@` syntax:
```lotus
wrap fn logging(fn wrapped) {
    printf("[LOG] Entering function\n");
    wrapped();
    printf("[LOG] Exiting function\n");
}

@logging
fn void greet() {
    printf("Hello, World!\n");
}
```

Multiple decorators can be stacked:
```lotus
@timing
@logging
fn void doWork() {
    // Both timing and logging wrappers applied
}
```

#### Currying (Partial Application)
Create partially applied functions:
```lotus
fn int add(int a, int b) {
    ret a + b;
}

int result = partial(add, 5)(3);  // Returns 8
```

### New Modules

#### Random Module (`use "random"`)
Full-featured random number generation using xorshift64 PRNG:
- `rand()` - Random int64
- `rand_range(min, max)` - Random int in range [min, max]
- `rand_n(n)` - Random int in range [0, n)
- `seed(n)` - Set PRNG seed
- `rand_float()` - Random float in [0.0, 1.0)
- `rand_bool()` - Random boolean
- `rand_bytes(buf, len)` - Fill buffer with random bytes
- `shuffle(arr, len)` - Fisher-Yates shuffle
- `choice(arr, len)` - Random element from array
- `rand_string(buf, len)` - Random alphanumeric string

### Extended String Functions

New string manipulation functions in the `str` module:
- `copy(str)` - Create string copy
- `compare(s1, s2)` - String comparison
- `indexOf(haystack, needle)` - Find substring position
- `substring(str, start, len)` - Extract substring
- `toUpper(str)` - Convert to uppercase
- `toLower(str)` - Convert to lowercase
- `trim(str)` - Remove leading/trailing whitespace
- `startsWith(str, prefix)` - Check string prefix
- `endsWith(str, suffix)` - Check string suffix
- `split(str, delim)` - Split string by delimiter
- `replace(str, old, new)` - Replace first occurrence

### Language Improvements

#### Void Return Type
Explicit void return type for functions:
```lotus
fn void cleanup() {
    // No return value
}
```

#### Type Conversion Functions
- `toUint32(val)` - Convert to 32-bit unsigned
- `toBool(val)` - Convert to boolean
- `toInt(val)` - Convert float to int
- `toFloat(val)` - Convert int to float

### LLVM Backend Enhancements

New code generation support in `llvm_advanced.go`:
- Ternary operator: `cond ? then : else`
- Compound assignment operators: `+=`, `-=`, `*=`, `/=`, `%=`
- Increment/decrement operators: `++`, `--`
- Full struct support with field access
- Enum definitions and literals
- Class definitions with methods
- Pointer operations: `&` (reference), `*` (dereference)
- `sizeof` expression
- `pow(base, exp)` using LLVM intrinsic

## Breaking Changes

None. This release is fully backward compatible with 1.3.x code.

## Bug Fixes

- Fixed token type mismatches in LLVM code generator
- Fixed FloatLiteral and CharLiteral handling in LLVM backend
- Improved error messages for functional programming syntax errors

## Installation

### Arch Linux (AUR)
```bash
yay -S lotus-lang
```

### From Source
```bash
git clone https://github.com/j-alexander3375/Lotus.git
cd Lotus/src
go build -o ../lotus .
```

## What's Next

- Pattern matching syntax
- Enhanced error handling with typed exceptions
- Standard library expansion (networking, file system)
- WebAssembly target improvements
