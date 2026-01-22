# Release Notes for Lotus v1.5.2

## Bug Fixes

### Float Type Constants Support
- **Fixed**: Added missing support for `float32` and `float64` constant declarations in the code generator
- Constants using `float32` and `float64` types can now be properly declared and used
- Resolves compilation errors when using floating-point type names in constant declarations

### Type System Improvements
- **Fixed**: Added `TokenTypeChar` to the `isTypeToken()` function for proper character type recognition
- Character types are now correctly validated in type-checking contexts

## Testing

### Comprehensive Test Suite Added
- **New**: Added extensive unit tests for tokenizer (`tokenizer_test.go`)
  - Tests for all keywords, type keywords, operators, literals, and edge cases
  - 100+ test cases covering the entire tokenization pipeline

- **New**: Added unit tests for type system (`types_test.go`)
  - Tests for type size calculations, type checking functions
  - Validation of integer, float, numeric, and primitive type classifications

- **New**: Added parser unit tests (`parser_test.go`)
  - Tests for variable declarations, constants, functions, control flow
  - Binary expression parsing, type name parsing validation

- **New**: Added integration tests (`integration_test.go`)
  - End-to-end compilation tests for various language features
  - Validation of example files and test suite compilation

All tests pass successfully, ensuring code quality and preventing regressions.

## Technical Details

### Modified Files
- `src/codegen.go`: Added `TokenTypeFloat32` and `TokenTypeFloat64` cases to `generateConstantDeclaration()`
- `src/parser.go`: Added `TokenTypeChar` to `isTypeToken()` function
- `src/constants.go`: Updated `CompilerVersion` from "1.5.1" to "1.5.2"

### Test Coverage
- Tokenizer: Keywords, operators, literals, comments, edge cases
- Parser: Declarations, expressions, statements, type system
- Integration: Full compilation pipeline validation
- Types: Size calculations, type classification functions

## Installation

Download the release tarball and build:
```bash
tar xzf lotus-lang-1.5.2.tar.gz
cd Lotus-1.5.2
cd src && go build -o ../lotus .
sudo install -Dm755 lotus /usr/bin/lotus
```

Or use the AUR package (Arch Linux):
```bash
yay -S lotus-lang
```

## Contributors

- Joshua Alexander (@j-alexander3375)
- Claude Sonnet 4.5 (Testing framework and bug fixes)

---

**Full Changelog**: https://github.com/j-alexander3375/Lotus/compare/1.5.1...1.5.2
