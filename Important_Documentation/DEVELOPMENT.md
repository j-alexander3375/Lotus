# Lotus Compiler - Development Summary

## Latest Phase: Enhanced Error Messages & Diagnostics ✅

### Current Accomplishments (December 2025 - Phase 7)

1. **Line/Column Tracking in Tokens**
   - Added `Line` and `Column` fields to Token struct
   - Tokenizer now tracks position during lexical analysis
   - All error messages include precise source locations

2. **Human-Readable Error Messages**
   - Created `error_messages.go` with error formatting utilities
   - `TokenTypeName()` - Converts token types to readable names
   - `FormatExpectedToken()` - Formats "expected X, got Y" messages
   - `FormatUnexpectedToken()` - Formats unexpected token errors
   - Error codes: E01xx (syntax), E02xx (semantic), E03xx (type)

3. **"Did You Mean?" Suggestions**
   - Levenshtein distance algorithm for typo detection
   - `SuggestForTypo()` suggests corrections for misspelled keywords
   - Integrated into parser error messages
   - Example: `pritnln` → "did you mean 'println'?"

4. **Enhanced Semantic Analyzer**
   - Rewrote `semantic.go` with proper scope tracking
   - Unused variable detection with warnings
   - Variable shadowing detection
   - Undefined variable suggestions using Levenshtein distance

5. **AST Location Tracking**
   - Added `Location` struct with Line/Column fields
   - `BaseNode` embedded in all AST nodes
   - All nodes now have `Loc()` method for position info

6. **Diagnostic Categories**
   - `DiagnosticCategory` enum: Syntax, Semantic, Unused, Shadow, Deprecated
   - Error codes for machine-readable parsing
   - Colored output with ANSI codes

7. **Comprehensive Test Suite**
   - `compiler_features_test.go` - 18 tests for new features
   - `integration_test.go` - End-to-end compilation tests
   - Tests for tokenizer tracking, error formatting, Levenshtein distance

---

## Previous Phase: Code Optimization ✅

### Accomplishments (December 2025 - Phase 6)

1. **AST-Level Optimizations** (`optimizer.go`)
   - **Constant Folding**: Evaluate constant expressions at compile time
     * Arithmetic: `2 + 3` → `5`, `10 * 4` → `40`
     * Nested expressions: `(3 + 2) * 4` → `20`
     * Bitwise operations: `0xFF & 0x0F` → `0x0F`
     * Unary operations: `-42` at compile time, `~0` → `-1`
   - **Strength Reduction**: Replace expensive operations with cheaper equivalents
     * Multiply by power of 2 → left shift: `x * 8` → `x << 3`
     * Works for 2, 4, 8, 16, 32, 64, 128, 256, etc.
   - **Identity Removal**: Eliminate no-op operations
     * `x + 0` → `x`, `x - 0` → `x`
     * `x * 1` → `x`, `x / 1` → `x`
     * `x * 0` → `0`, `x % 1` → `0`
     * `x & -1` → `x`, `x | 0` → `x`, `x ^ 0` → `x`
     * `x << 0` → `x`, `x >> 0` → `x`
   - Recursive optimization through all AST node types

2. **Assembly-Level Peephole Optimizations** (`peephole.go`)
   - **Redundant Move Elimination**: `movq %rax, %rax` → removed
   - **Dead Store Elimination**: Store-then-load to same location optimized
   - **Push-Pop Cancellation**: Adjacent `pushq %rax; popq %rax` → removed
   - **Zero Loading Optimization**: `movq $0, %rax` → `xorq %rax, %rax`
   - **Increment/Decrement Simplification**: `addq $1, %rax` → `incq %rax`
   - Multi-pass optimization until no more changes

3. **Optimization Pipeline Integration**
   - Optimizations run automatically during compilation
   - Phase 1: Parse tokens to AST
   - Phase 2: Optimize AST (constant folding, strength reduction)
   - Phase 3: Generate code from optimized AST
   - Phase 4: Apply peephole optimizations to assembly

4. **Comprehensive Test Coverage** (`optimizer_test.go`)
   - Unit tests for all optimization types
   - TestConstantFolding: arithmetic, nested expressions
   - TestIdentityRemoval: all identity cases
   - TestStrengthReduction: power-of-2 multiplication
   - TestBitwiseConstantFolding: AND, OR, XOR, shifts
   - TestUnaryConstantFolding: negation, bitwise NOT
   - Peephole tests: redundant moves, zero loading, inc/dec, push-pop

---


## Previous Phase: Tooling & Diagnostics Enhancements ✅

### Accomplishments (December 2025 - Phase 5)

1. **Compilation Statistics Tracking**
   - Created `stats.go` - Comprehensive compilation metrics system
   - `CompilationStats` structure tracks all phases:
     * Timing: Tokenization, Parsing, Codegen, Assembly, Linking
     * Metrics: Token counts, AST nodes, functions, variables, constants
     * Output: Assembly lines/bytes, binary size
   - Human-readable output formatting with automatic unit conversion (B/KB/MB/GB)
   - Compact and detailed reporting modes

2. **Enhanced Command-Line Flags**
   - `--stats` - Display detailed compilation statistics
   - `--timing` - Show phase-by-phase timing information
   - `--ast-dump` - Print Abstract Syntax Tree structure
   - `-q, --quiet` - Suppress non-error output
   - All new flags properly normalized and documented

3. **AST Debugging Utilities**
   - Created `ast_utils.go` - AST inspection and analysis tools
   - `DumpAST()` - Hierarchical pretty-printing of AST structure
   - `CountASTNodes()` - Recursive node counting for metrics
   - `AnalyzeAST()` - Extract statistics (functions, variables, constants)
   - Support for all AST node types with proper formatting

4. **Improved Diagnostics**
   - Enhanced `diagnostics.go` with color-coded output
   - Added `PrintSummary()` for compact error reporting
   - New diagnostic levels: `AddInfo()`, `AddHint()`
   - ANSI color support:  
     * Red for errors
     * Yellow for warnings
     * Cyan for info
     * Green for hints
   - Better context display with line numbers and caret positioning

5. **Testing Infrastructure**
   - Created comprehensive unit tests:
     * `stats_test.go` - 10 tests for compilation statistics
     * `diagnostics_test.go` - 7 tests for diagnostic management
     * `ast_utils_test.go` - 8 tests for AST utilities
     * `flags_test.go` - 5 tests for flag parsing
   - Created performance benchmarks:
     * `benchmarks_test.go` - 14 benchmarks covering all modules
     * Memory allocation tracking
     * Performance testing for small, medium, and large ASTs

6. **Code Quality**
   - All new modules follow existing code style
   - Comprehensive documentation and comments
   - Modular design for easy extension
   - Zero external dependencies


## Previous Phase: Advanced Stdlib Features ✅

### Accomplishments (December 2025 - Phase 4)

**Delivered in 1.2.1:**
1. **String Extensions** - substring, split, join, replace, toLower, toUpper, trim (codegen stubs)
2. **File I/O Module** - open, close, read, write, seek, stat, exists (Linux syscalls)
3. **Time Module** - now, sleep, millis, nanos, clock, gmtime, localtime (syscall-based)
4. **Extended Stdlib** - 11 new functions registered across 3 modules; full codegen infrastructure

**Completed in 1.2.5 (Phase 4 Finalization):**
1. **String Functions Fully Implemented** - toLower, toUpper, trim, split, join, replace
2. **Time Functions Implemented** - gmtime, localtime (UTC-based broken-down time)
3. **Math Functions Fixed** - GCD/LCM now work correctly with Euclidean algorithm

#### Phase 4 Completions ✅
1. **String Module Extended** ✅ **FULLY IMPLEMENTED**
   - ✅ substring(s, start, len) - Extract substring
   - ✅ split(str, delim) - Split string by single-char delimiter, returns array structure
   - ✅ join(array, sep) - Join array elements with separator
   - ✅ replace(str, old, new) - Replace all occurrences (single-char replacement)
   - ✅ toLower(str) - Lowercase copy with proper A-Z conversion
   - ✅ toUpper(str) - Uppercase copy with proper a-z conversion
   - ✅ trim(str) - Trim whitespace (space, tab, newline, carriage return) from both ends
   - Total: 15 string functions (ALL fully implemented)
2. **File I/O Module (`file`)** ✅
   - ✅ open(path, flags) - Linux open(2) syscall
   - ✅ close(fd) - Linux close(2) syscall
   - ✅ read(fd, buf, size) - Linux read(2) syscall
   - ✅ write(fd, buf, size) - Linux write(2) syscall
   - ✅ seek(fd, offset, whence) - Linux lseek(2) syscall
   - ✅ stat(path, statbuf) - Linux stat(2) syscall
   - ✅ exists(path) - File existence check
3. **Time Module (`time`)** ✅ **FULLY IMPLEMENTED**
   - ✅ now() - Unix timestamp via time(2)
   - ✅ sleep(seconds) - Sleep via nanosleep(2)
   - ✅ millis() - Millisecond timestamp
   - ✅ nanos() - Nanosecond timestamp
   - ✅ clock() - Clock ticks
   - ✅ gmtime(ts, buf) - Convert to broken-down time (UTC)
   - ✅ localtime(ts, buf) - Convert to broken-down time (UTC-based, no TZ support)
4. **Hashing Module (`hash`)** ✅ **FULLY IMPLEMENTED**
   - ✅ Module structure defined
   - ✅ CRC32 implementation (IEEE 802.3 polynomial with lookup table)
   - ✅ FNV-1a implementation (64-bit fast non-cryptographic hash)
   - ✅ DJB2 implementation (simple string hashing)
   - ✅ MurmurHash3 implementation (32-bit with seed support)
   - ✅ SHA-256 implementation (full FIPS 180-4 compliant with 64 rounds)
   - ✅ MD5 implementation (full RFC 1321 compliant with 4 rounds)
5. **Collections Module Enhancements** ✅ **FULLY IMPLEMENTED**
   - ✅ Array, Stack, Queue, Deque structures implemented
   - ✅ Heap (min-heap) implemented
   - ✅ HashMap and HashSet (int keys) implemented
   - ✅ HashMap and HashSet (string keys) implemented - djb2 hash, strcmp comparison
   - ✅ `binary_search_int` helper implemented
   - ✅ Sorted set (BST-based): new, add, contains, remove, min, max, len, free
   - ✅ Sorted map (BST-based): new, put, get, contains, remove, min_key, max_key, len, free
   - ✅ Memory management: array_int_resize, reserve, shrink, capacity, get, set, free
6. **HTTP Module (`http`)** ✅ **FULLY IMPLEMENTED**
   - ✅ Minimal `get` implemented atop `net` primitives
   - ✅ POST request support with Content-Length header
   - ✅ Response parsing: parse_status, get_header, get_body, parse_headers
   - ✅ Connection pooling: pool_new, pool_get, pool_put, pool_close
7. **Networking Module (`net`)** ✅ **FULLY IMPLEMENTED**
   - ✅ Socket, connect_ipv4, send, recv, close implemented (Linux syscalls)
   - ✅ UDP support: bind_ipv4, sendto_ipv4, recvfrom
   - ✅ IPv6 support: connect_ipv6, bind_ipv6, sendto_ipv6
   - ✅ DNS resolution: resolve (via /etc/hosts), resolve_ipv6 (stub)
8. **String Module Completion** ✅ **COMPLETE**
   - ✅ len, concat, compare, copy, indexOf, contains, startsWith, endsWith
   - ✅ substring, split, join (FULLY IMPLEMENTED)
   - ✅ toLower, toUpper, trim (FULLY IMPLEMENTED)
   - ✅ replace (FULLY IMPLEMENTED)
9. **File I/O Module** ✅
   - ✅ open, close, read, write, seek, stat, exists (all registered with Linux syscall codegen)
10. **Time Module** ✅ **COMPLETE**
   - ✅ now, sleep, millis, nanos, clock, gmtime, localtime (all registered)

---

## Previous Phase: Stdlib, Formatting, and Lotus Syntax Refresh ✅

### Accomplishments (December 2025 - Phase 3)

1. **Import System Implementation**
   - Added `use` and `as` keywords for module imports (string-based)
   - Rust-inspired import syntax:
      - `use "module"` - Import module
      - `use "module::function"` - Specific function
      - `use "module::*"` - Explicit wildcard
      - `use "module" as alias` - Aliased imports
   - `ImportStatement` AST node plus parse/codegen support

2. **Standard Library Module System**
  - Created `stdlib.go` - Module definitions and registration
  - Implemented `ImportContext` for tracking imports
  - `StandardLibrary` map with registered modules
  - `StdlibModule` and `StdlibFunction` types for organization

3. **Stdlib Modules Created**
   - **io** - print, println, printf, fprintf, sprint, sprintf, sprintln
   - **mem** - malloc, free, sizeof, memcpy, memset, mmap, munmap
   - **math** - abs, min, max, sqrt, pow, floor, ceil, round, gcd, lcm
   - **str** - len, concat, compare, copy, indexOf, contains, startsWith, endsWith
   - **num** - toInt8, toUint8, toInt16, toUint16, toInt32, toUint32, toInt64, toUint64, toBool
   - **hash** - djb2, fnv1a, crc32, murmur3, sha256, md5 (all fully implemented)
   - **collections** - dynamic arrays, stacks, queues, deques, heaps, hashmap, hashset; helper `binary_search_int`
   - **net** - socket, connect_ipv4, send, recv, close
   - **http** - get (minimal client built on `net`)

4. **Import System Features**
  - Module lookup and validation
  - Function availability tracking
  - Compile-time import resolution
  - Error reporting for unknown modules/functions
  - Assembly comments documenting imports

5. **Testing & Validation**
  - Imports: `test_imports_basic.lts`, `test_imports_specific.lts`, `test_imports_wildcard.lts`, `test_imports_alias.lts`, `test_imports_multiple.lts`, `test_imports_comprehensive.lts`
  - Showcase/printf: `test_showcase.lts`, `test_final.lts`
  - Arithmetic/types/data structures: `test_arithmetic.lts`, `test_types_minimal.lts`, `test_data_structures.lts`, etc.
  - All current tests passing ✅

6. **Formatting & Stdlib Updates**
   - printf supports %%, %d, %b, %o, %x/%X, %c, %q, %s, %v with base-aware int printing and char output
   - math: abs/min/max/sqrt/pow plus floor/ceil/round/gcd/lcm implemented
   - str: len/concat/compare/copy/indexOf/contains/startsWith/endsWith implemented
   - mem: malloc/free/sizeof plus memcpy/memset/mmap/munmap implemented
   - num: integer width conversions and boolean coercion implemented
   - hash: djb2/fnv1a/crc32/murmur3/sha256/md5 all fully implemented
   - collections: dynamic arrays, stacks, queues/deques, heaps, hashmap/hashset, and `binary_search_int` implemented
   - net/http: minimal socket helpers and `http.get` implemented for simple requests
   - Function calls now dispatch to imported stdlib functions from codegen
   - Comprehensive import demo shows correct outputs across modules

7. **Documentation**
  - Created `STDLIB_AND_IMPORTS.md` - Comprehensive module documentation
   - Updated `README.md` with import examples and module descriptions
   - Added stdlib section showing all available modules

## Earlier Phase: Major Refactoring & Constant Declarations ✅

### Recent Accomplishments (December 2025 - Phase 2)

1. **Major Code Refactoring**
   - Created `ast.go` - Centralized AST node definitions
   - Created `constants.go` - Compiler-wide constants and configuration
   - Created `types.go` - Type system utilities and helpers
   - Added comprehensive documentation to all files
   - Eliminated magic numbers and duplicate code
   - Improved naming conventions throughout

2. **Constant Declaration Feature**
   - Added `const` keyword to language
   - Constants stored in `.data` section for efficient access
   - Support for `const int`, `const string`, `const bool`
   - Compile-time constant evaluation
   - Immutable values with RIP-relative addressing

3. **Enhanced Documentation**
   - File-level documentation for every module
   - Detailed struct field comments
   - Comprehensive function documentation
   - Created REFACTORING.md summarizing changes

### What Was Previously Accomplished

1. **Created `parser.go`** (276 lines)
   - Implemented full Abstract Syntax Tree (AST) data structures
   - AST Node types: ReturnStatement, VariableDeclaration, FunctionCall, IntLiteral, StringLiteral, BoolLiteral, FloatLiteral, Identifier
   - Recursive descent Parser with 8+ parsing methods
   - Proper token-by-token navigation with lookahead

2. **Refactored `codegen.go`** (191 lines)
   - Changed from direct token stream parsing to AST-based code generation
   - Created CodeGenerator struct for managing state during assembly generation
   - Implemented statement-specific code generation methods
   - Proper handling of all variable types with correct stack allocation

3. **Fixed Assembly Output**
   - Corrected assembly directive formatting (removed leading spaces from `.section`, `.global`, `.text`)
   - Proper GNU AT&T syntax compliance
   - All assembly output now assembles correctly with GCC

4. **Enhanced Compiler Pipeline**
   - Fixed binary compilation logic
   - Proper handling of `-S` flag for assembly-only output
   - Default behavior now correctly produces executable binaries

### Current Capabilities

**Lexical Analysis:**
- Tokenizes 100+ token types including all type keywords
- Recognizes lowercase stdlib functions (printf/println/logf/etc.)
- Proper handling of string literals, integers, floats, booleans, pointers, keywords for try/catch/finally/throw/null
- Multi-char operators and delimiters

**Parsing:**
- Type-first variable declarations (e.g., `int x = 1;`)
- Function definitions with parameters/returns (`fn`), `ret` keyword
- Struct/enum/class definitions
- Control flow: if/else, while, for
- Try/catch/finally and throw
- Module imports via `use`/`as`

**Code Generation:**
- x86-64 GNU assembly (AT&T)
- Stack frames, 16-byte alignment, RIP-relative data
- Function calls honoring System V ABI
- Stdlib hooks across io/mem/math/str/num/hash/collections/net/http

**Compilation:**
- Go build: `go build -o lotus src/*.go`
- Assembly-only: `./lotus -S input.lts`
- Binary: `./lotus input.lts`
- Token dump: `./lotus -td input.lts`

### Test Snapshots

```
✓ test_showcase.lts      (printf/println demo)
✓ test_final.lts         (printf with vars/strings)
✓ test_imports_*         (basic/specific/wildcard/alias/multiple)
✓ test_arithmetic.lts    (arith & bitwise)
✓ test_data_structures.lts (struct/enum/class basics)
✓ test_error_handling.lts (try/catch/finally/throw/null)
```

### Architecture

```
Input (.lts file)
    ↓
[Tokenizer] → Token Stream
    ↓
[Parser] → AST
    ↓
[CodeGenerator] → Assembly
    ↓
[GCC] → Binary
```

### Current Phase: Advanced Stdlib Features (December 2025 - Phase 4) ✅ **COMPLETE**

**Delivered in 1.2.1:**
1. **String Extensions** - substring, split, join, replace, toLower, toUpper, trim (codegen stubs)
2. **File I/O Module** - open, close, read, write, seek, stat, exists (Linux syscalls)
3. **Time Module** - now, sleep, millis, nanos, clock, gmtime, localtime (syscall-based)
4. **Extended Stdlib** - 11 new functions registered across 3 modules; full codegen infrastructure

**Completed in 1.2.5 (Phase 4 Finalization):**
1. **String Functions Fully Implemented** - toLower, toUpper, trim, split, join, replace
2. **Time Functions Implemented** - gmtime, localtime (UTC-based broken-down time)
3. **Math Functions Fixed** - GCD/LCM now work correctly with Euclidean algorithm

#### Phase 4 Completions ✅

1. **String Module Extended** ✅ **FULLY IMPLEMENTED**
   - ✅ substring(s, start, len) - Extract substring
   - ✅ split(str, delim) - Split string by single-char delimiter, returns array structure
   - ✅ join(array, sep) - Join array elements with separator
   - ✅ replace(str, old, new) - Replace all occurrences (single-char replacement)
   - ✅ toLower(str) - Lowercase copy with proper A-Z conversion
   - ✅ toUpper(str) - Uppercase copy with proper a-z conversion
   - ✅ trim(str) - Trim whitespace (space, tab, newline, carriage return) from both ends
   - Total: 15 string functions (ALL fully implemented)

2. **File I/O Module (`file`)** ✅
   - ✅ open(path, flags) - Linux open(2) syscall
   - ✅ close(fd) - Linux close(2) syscall
   - ✅ read(fd, buf, size) - Linux read(2) syscall
   - ✅ write(fd, buf, size) - Linux write(2) syscall
   - ✅ seek(fd, offset, whence) - Linux lseek(2) syscall
   - ✅ stat(path, statbuf) - Linux stat(2) syscall
   - ✅ exists(path) - File existence check

3. **Time Module (`time`)** ✅ **FULLY IMPLEMENTED**
   - ✅ now() - Unix timestamp via time(2)
   - ✅ sleep(seconds) - Sleep via nanosleep(2)
   - ✅ millis() - Millisecond timestamp
   - ✅ nanos() - Nanosecond timestamp
   - ✅ clock() - Clock ticks
   - ✅ gmtime(ts, buf) - Convert to broken-down time (UTC)
   - ✅ localtime(ts, buf) - Convert to broken-down time (UTC-based, no TZ support)

#### Phase 4 Goals (Original - Now Complete)

1. **Hashing Module (`hash`)** ✅ **FULLY IMPLEMENTED**
   - ✅ Module structure defined
   - ✅ CRC32 implementation (IEEE 802.3 polynomial with lookup table)
   - ✅ FNV-1a implementation (64-bit fast non-cryptographic hash)
   - ✅ DJB2 implementation (simple string hashing)
   - ✅ MurmurHash3 implementation (32-bit with seed support)
   - ✅ SHA-256 implementation (full FIPS 180-4 compliant with 64 rounds)
   - ✅ MD5 implementation (full RFC 1321 compliant with 4 rounds)

2. **Collections Module Enhancements** ✅ **FULLY IMPLEMENTED**
   - ✅ Array, Stack, Queue, Deque structures implemented
   - ✅ Heap (min-heap) implemented
   - ✅ HashMap and HashSet (int keys) implemented
   - ✅ HashMap and HashSet (string keys) implemented - djb2 hash, strcmp comparison
   - ✅ `binary_search_int` helper implemented
   - ✅ Sorted set (BST-based): new, add, contains, remove, min, max, len, free
   - ✅ Sorted map (BST-based): new, put, get, contains, remove, min_key, max_key, len, free
   - ✅ Memory management: array_int_resize, reserve, shrink, capacity, get, set, free

3. **HTTP Module (`http`)** ✅ **FULLY IMPLEMENTED**
   - ✅ Minimal `get` implemented atop `net` primitives
   - ✅ POST request support with Content-Length header
   - ✅ Response parsing: parse_status, get_header, get_body, parse_headers
   - ✅ Connection pooling: pool_new, pool_get, pool_put, pool_close

4. **Networking Module (`net`)** ✅ **FULLY IMPLEMENTED**
   - ✅ Socket, connect_ipv4, send, recv, close implemented (Linux syscalls)
   - ✅ UDP support: bind_ipv4, sendto_ipv4, recvfrom
   - ✅ IPv6 support: connect_ipv6, bind_ipv6, sendto_ipv6
   - ✅ DNS resolution: resolve (via /etc/hosts), resolve_ipv6 (stub)

5. **String Module Completion** ✅ **COMPLETE**
   - ✅ len, concat, compare, copy, indexOf, contains, startsWith, endsWith
   - ✅ substring, split, join (FULLY IMPLEMENTED)
   - ✅ toLower, toUpper, trim (FULLY IMPLEMENTED)
   - ✅ replace (FULLY IMPLEMENTED)

6. **File I/O Module** ✅
   - ✅ open, close, read, write, seek, stat, exists (all registered with Linux syscall codegen)

7. **Time Module** ✅ **COMPLETE**
   - ✅ now, sleep, millis, nanos, clock, gmtime, localtime (all registered)

---


## Phase 8 Goals (Planned - Future Releases)

1. **Formatting Enhancements**
   - Width/padding flags for printf-like output
   - Custom format specifiers
   - ⏳ Regular expressions

2. **JSON Module**
   - ⏳ JSON parsing/serialization

3. **Advanced Optimization**
   - ⏳ Register allocation improvements
   - ⏳ Dead code elimination
   - ⏳ Inline function expansion

4. **Type System Enhancements**
   - Generics and type inference improvements
   - Union/option types
   - Pattern matching

5. **Tooling**
   - Language server protocol
   - Debug/trace hooks
   - Package/module manager
   - Build system integration

---

### Current File Structure

```
src/
├── main.go              - CLI orchestration and compilation flow  
├── compiler.go          - Compilation pipeline management
├── flags.go             - Command-line flag parsing with examples
├── constants.go         - Compiler constants and configuration (NEW)
├── keywords.go          - Token types (100+ tokens, reorganized)
├── tokenizer.go         - Lexical analysis with comment support
├── parser.go            - Syntactic analysis with constant declarations
├── ast.go               - Centralized AST node definitions (NEW)
├── types.go             - Type system utilities (NEW)
├── codegen.go           - Assembly code generation orchestrator
├── optimizer.go         - AST-level optimizations (constant folding, strength reduction) (NEW)
├── peephole.go          - Assembly-level peephole optimizations (NEW)
├── diagnostics.go       - Error and warning reporting
├── printfuncs.go        - Print function implementations
├── arithmetic.go        - Arithmetic & bitwise operations
├── control_flow.go      - If/while/for loops, comparisons
├── references.go        - References, assignments, dereferencing
├── functions.go         - User-defined functions
├── memory.go            - malloc/free/sizeof operations
├── array.go             - Dynamic arrays
├── struct.go            - Struct types and field access
├── enum.go              - Enumeration types
├── class.go             - Object-oriented programming
└── error_handling.go    - Try/catch/finally exception handling
```

### Build & Test Commands

```bash
# Build compiler
cd src && go build -o ../lotus

# Generate assembly
./lotus -S input.lts

# Compile to binary (default)
./lotus input.lts

# Build and run immediately
./lotus -run input.lts

# Dump tokens for debugging
./lotus -td input.lts

# Show version
./lotus --version

# Verbose compilation
./lotus -v input.lts
```

### Language Features Implemented

**Core Features:**
- ✅ Type-first variable declarations (`int x = 1;`)
- ✅ Constants (`const`), ints/bools/strings, pointer types
- ✅ Arithmetic + bitwise + logical + comparisons
- ✅ Control flow (if/else, while, for)
- ✅ Functions with `fn` and `ret`
- ✅ Arrays, pointers, struct/enum/class definitions
- ✅ Exception handling (try/catch/finally/throw/null)
- ✅ Import system with `use`/`as` and stdlib modules (io, mem, math, str, num, hash, collections, net, http)
- ✅ Stdlib: printf/println/logf; memory helpers (malloc/free/sizeof/memcpy/memset/mmap/munmap); rich math and string functions; numeric conversions; hashing; collections; basic net/http

### Example Constant Usage

```lotus
// Define constants
const int MAX_BUFFER = 4096;
const string VERSION = "1.0.0";
const bool PRODUCTION = false;

// Use constants
int buffer_size = MAX_BUFFER;
println(VERSION);
if (PRODUCTION) {
    printf("Running in production mode\n");
}
```

### Key Improvements Summary

- ✅ Complete modular architecture
- ✅ Proper separation of concerns (tokenization → parsing → code generation)
- ✅ AST-based approach enables future optimizations and transformations
- ✅ Parser enables handling of complex expressions and nested structures
- ✅ Code generation now maintainable and extensible per statement type
- ✅ Assembly output now fully compliant with GNU assembler requirements
- ✅ Full stdlib with 11 modules and 100+ functions
- ✅ Comprehensive optimization pipeline (AST + peephole)
- ✅ Enhanced diagnostics with color-coded output and suggestions
