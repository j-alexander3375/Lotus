# Release Notes for Lotus v1.11.0

## Highlights

This minor release completes the exception-handling system with real stack unwinding across
function calls, and moves the toolchain to modern LLVM (21/22) so the compiler builds cleanly
against current system toolchains.

## New Features

### Language Features
- `throw` now unwinds across function calls: an exception thrown inside a callee propagates
  to the nearest enclosing `try` anywhere up the call stack.
- `catch` binding forms: `catch { }` (no binding), `catch (e) { }` (binds as `int`), and
  `catch (TYPE e) { }` which reinterprets the thrown value as `int`, `string`, `bool`,
  `float`, or `char`. One `catch` clause per `try`.
- `finally` is guaranteed to run on normal completion, on caught and uncaught exceptions,
  and on early exits from the `try`/`catch` body via `ret`/`break`/`continue`.
- An uncaught `throw` prints an error and exits the program with status 1.
- New example: `examples/error_handling.lts` demonstrating unwinding across a function call.

### Compiler & Runtime
- New exception runtime (`src/llvm_runtime.go`, `src/llvm_trycatch.go`, `src/runtimec/`)
  backing the try/catch implementation.
- `println`/`print` now print all of their arguments (space-separated) instead of silently
  dropping extras.
- Object-file output fixed: `-o` reliably writes the compiled binary.
- Substantial fixes and improvements across LLVM codegen, the optimizer, parser, tokenizer,
  and interpreter.

## Build & Packaging
- Vendored `tinygo.org/x/go-llvm` updated with support for LLVM 19 through 22; the compiler
  now builds against system LLVM 22.
- PKGBUILD: `depends` pinned to `llvm>=22`/`llvm<23`; `check()` fixed to pass `-o` before the
  input file (flags after the filename are ignored).

## Upgrade Notes
- Compiler flags must precede the input file: `lotus -o out program.lts`.
- No source-level breaking changes are expected from 1.10.0.
