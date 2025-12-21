# Constant Declaration Feature - Implementation Summary

## Overview
Added support for constant declarations in Lotus language, allowing developers to define immutable compile-time values that are stored efficiently in the data section.

## Syntax
```lotus
const <type> <name> = <value>;
```

### Examples
```lotus
const int MAX_SIZE = 1024;
const string APP_NAME = "Lotus";
const bool DEBUG_MODE = true;
```

## Implementation Details

### 1. Token Addition
**File**: `src/keywords.go`
- Added `TokenConst` to token types enumeration
- Positioned after `TokenReturn` in control flow keywords section

### 2. Tokenizer Support
**File**: `src/tokenizer.go`
- Added `"const"` keyword recognition in tokenizer switch statement
- Maps `const` keyword to `TokenConst` token type

### 3. AST Node Definition
**File**: `src/ast.go`
- Created `ConstantDeclaration` struct with:
  - `Name` (string) - constant identifier
  - `Type` (TokenType) - data type
  - `Value` (ASTNode) - constant value expression
- Implements `astNode()` interface method

### 4. Parser Implementation
**File**: `src/parser.go`

Added `parseConstantDeclaration()` method:
- Expects `const` keyword
- Requires type specification (int, string, bool, etc.)
- Requires identifier name
- Expects `=` assignment operator
- Parses constant value expression
- Returns `*ConstantDeclaration` AST node

Updated `parseStatement()`:
- Added `case TokenConst:` to route to constant parser
- Returns parsed constant declaration

Added helper:
- `isTypeToken()` function to validate type tokens

### 5. Code Generation
**File**: `src/codegen.go`

Updated `CodeGenerator` struct:
- Added `constants map[string]Variable` field to track declared constants
- Initialized in `NewCodeGenerator()`

Added `generateConstantDeclaration()` method:
- **Integer constants**: Generate `.quad` directive in data section with label `.const_<name>`
- **Boolean constants**: Convert bool to int (0/1), store as `.quad`
- **String constants**: Generate `.asciz` directive with escaped string
- Store metadata in `constants` map for later reference
- Comments generated in text section for debugging

Updated `generateStatement()`:
- Added `case *ConstantDeclaration:` to dispatch to constant generator

### 6. Constant Reference Support
**File**: `src/arithmetic.go`

Updated `generateExpressionToReg()`:
- Enhanced `case *Identifier:` to check both variables and constants
- For constants:
  - String constants: Use `leaq` to load address
  - Numeric/bool constants: Use `movq` to load value
- Uses RIP-relative addressing for position-independent code

## Assembly Output

### Integer Constant
```lotus
const int MAX = 100;
```
Generates:
```asm
.section .data
.const_MAX:
    .quad 100
```

### String Constant
```lotus
const string NAME = "Lotus";
```
Generates:
```asm
.section .data
.const_NAME:
    .asciz "Lotus"
```

### Boolean Constant
```lotus
const bool ENABLED = true;
```
Generates:
```asm
.section .data
.const_ENABLED:
    .quad 1
```

## Usage in Code

Constants can be referenced like variables:
```lotus
const int SIZE = 1024;
int buffer = SIZE;  // Loads value from .const_SIZE label
```

## Benefits

1. **Performance**: Constants stored in data section, no stack allocation needed
2. **Type Safety**: Type checking at compile time
3. **Immutability**: Cannot be modified after declaration
4. **Memory Efficient**: Single data section entry, no runtime overhead
5. **Position Independent**: Uses RIP-relative addressing

## Limitations (Current)

1. **No Constant Expressions**: Value must be a literal (no arithmetic in declaration)
2. **No Type Inference**: Type must be explicitly specified
3. **No Constant Folding**: Complex expressions not evaluated at compile time

## Future Enhancements

1. **Constant Expression Evaluation**
   ```lotus
   const int SIZE = 1024 * 1024;  // Evaluate at compile time
   ```

2. **Derived Constants**
   ```lotus
   const int HALF_SIZE = SIZE / 2;  // Reference other constants
   ```

3. **Enum-style Constants**
   ```lotus
   const int RED = 0;
   const int GREEN = 1;
   const int BLUE = 2;
   ```

4. **Constant Arrays**
   ```lotus
   const int[] PRIMES = [2, 3, 5, 7, 11];
   ```

## Testing

Test files created:
- `tests/test_constants.lts` - Basic constant declarations
- `tests/test_constants_comprehensive.lts` - All constant types

All tests pass successfully:
```bash
$ ./lotus -run tests/test_constants.lts
Testing constants in Lotus

$ ./lotus -run tests/test_constants_comprehensive.lts
=== Lotus Constant Declaration Test ===
...
All constant declarations successful!
```

## Documentation Updates

Updated files:
- `README.md` - Added constants section with syntax and examples
- `DEVELOPMENT.md` - Listed constant declarations as completed feature
- `src/ast.go` - Documented ConstantDeclaration node
- `src/codegen.go` - Documented constant generation logic
- `src/parser.go` - Documented parseConstantDeclaration method

## Compiler Version
This feature is included in Lotus compiler version **0.3.0**.
