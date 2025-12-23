# Lotus Language - Bitwise and Math Functions Support

## Summary of Changes

This update adds comprehensive bitwise operator support and additional math functions to the Lotus programming language compiler.

## New Features

### 1. Bitwise Operators

All standard bitwise operators have been added with proper operator precedence:

#### Binary Bitwise Operators
- **`&` (Bitwise AND)**: Performs bitwise AND operation
  - Example: `int result = a & b;`
- **`|` (Bitwise OR)**: Performs bitwise OR operation
  - Example: `int result = a | b;`
- **`^` (Bitwise XOR)**: Performs bitwise exclusive OR
  - Example: `int result = a ^ b;`
- **`<<` (Left Shift)**: Shifts bits left by specified amount
  - Example: `int result = a << 2;`
- **`>>` (Right Shift)**: Shifts bits right by specified amount
  - Example: `int result = a >> 2;`

#### Unary Bitwise Operators
- **`~` (Bitwise NOT)**: Performs bitwise NOT (one's complement)
  - Example: `int result = ~a;`

### 2. Logical Operators

Logical operators are now properly integrated with operator precedence:

- **`&&` (Logical AND)**: Short-circuit AND operation
- **`||` (Logical OR)**: Short-circuit OR operation

### 3. Enhanced Math Module

The math module now includes the following functions:

#### Existing Functions
- `abs(x)` - Absolute value
- `min(a, b)` - Minimum of two values
- `max(a, b)` - Maximum of two values
- `sqrt(x)` - Square root (TODO: full implementation)
- `pow(base, exp)` - Power function (TODO: full implementation)

#### New Functions
- `floor(x)` - Floor function (rounds down)
- `ceil(x)` - Ceiling function (rounds up)
- `round(x)` - Rounding function
- `gcd(a, b)` - Greatest Common Divisor (Euclidean algorithm)
- `lcm(a, b)` - Least Common Multiple

### 4. Operator Precedence

The complete operator precedence hierarchy (from lowest to highest):

1. Logical OR (`||`)
2. Logical AND (`&&`)
3. Bitwise OR (`|`)
4. Bitwise XOR (`^`)
5. Bitwise AND (`&`)
6. Comparison (`==`, `!=`, `<`, `>`, `<=`, `>=`)
7. Bit Shifts (`<<`, `>>`)
8. Addition/Subtraction (`+`, `-`)
9. Multiplication/Division/Modulo (`*`, `/`, `%`)
10. Unary operators (`-`, `!`, `~`, `&` address-of, `*` dereference)

## Implementation Details

### Files Modified

1. **parser.go**
   - Added bitwise operator parsing with proper precedence
   - Added logical operator parsing with short-circuit evaluation support
   - Functions added:
     - `parseLogicalOr()` - Handles `||` operator
     - `parseLogicalAnd()` - Handles `&&` operator
     - `parseBitwiseOr()` - Handles `|` operator
     - `parseBitwiseXor()` - Handles `^` operator
     - `parseBitwiseAnd()` - Handles binary `&` operator
     - `parseShift()` - Handles `<<` and `>>` operators
     - `isUnaryContext()` - Disambiguates `&` as unary (address-of) vs binary

2. **arithmetic.go**
   - Updated UnaryOp comment to include TokenTilde
   - `generateBitwiseOp()` function already existed and handles all binary bitwise operations
   - `generateUnaryOp()` function already supported TokenTilde for bitwise NOT

3. **stdlib.go**
   - Expanded math module in `createMathModule()`:
     - Added: `floor`, `ceil`, `round`, `gcd`, `lcm`
   - Added code generation functions:
     - `generateMathFloor()` - Integer floor function
     - `generateMathCeil()` - Integer ceiling function
     - `generateMathRound()` - Integer rounding function
     - `generateMathGcd()` - GCD using Euclidean algorithm
     - `generateMathLcm()` - LCM calculation

### Control Flow

The `control_flow.go` file already contained:
- `LogicalOp` struct definition
- `generateLogicalOp()` implementation with short-circuit evaluation

These were integrated with the new parser precedence.

## Usage Examples

### Bitwise Operations

```lotus
// Bitwise AND
int a = 12;  // 1100
int b = 10;  // 1010
int and_result = a & b;  // 8 (1000)
println(and_result);

// Bitwise OR
int or_result = a | b;  // 14 (1110)
println(or_result);

// Bitwise XOR
int xor_result = a ^ b;  // 6 (0110)
println(xor_result);

// Bitwise NOT
int not_result = ~a;  // Flips all bits
println(not_result);

// Shift operations
int left_shift = a << 1;  // 24
int right_shift = a >> 1;  // 6
```

### Logical Operations

```lotus
int x = 5;
int y = 0;

// Short-circuit AND
if (x && y) {
    println("Both truthy");
}

// Short-circuit OR
if (x || y) {
    println("At least one truthy");
}
```

### Math Functions

```lotus
use "math";

int abs_val = math::abs(-42);           // 42
int minimum = math::min(10, 5);         // 5
int maximum = math::max(10, 5);         // 10
int gcd_result = math::gcd(48, 18);     // 6
int lcm_result = math::lcm(12, 8);      // 24
```

## Testing

A comprehensive test file has been created at `tests/test_bitwise_math.lts` that covers:
- All bitwise operations
- Logical operations
- Math module functions
- Both direct function calls and module-qualified calls

## Assembly Code Generation

The compiler generates proper x86-64 assembly for all operations:
- Bitwise operations use: `andq`, `orq`, `xorq`, `notq`, `shlq`, `shrq`
- Logical operations use conditional jumps with short-circuit evaluation
- Math functions are generated based on their implementation

## Future Enhancements

The following placeholders are marked for future implementation:
- `sqrt()` - Square root calculation (TODO)
- `pow()` - Power function (TODO)
- More trigonometric functions
- Floating-point math functions

## Compatibility

All changes maintain backward compatibility with existing Lotus code. The new operators and functions are additions that don't modify existing functionality.

