# Release Notes for Lotus v1.9.0

## Highlights

This minor release expands Lotus functional programming support and stabilizes core functional workflows end-to-end.
The `func` module now includes applicative-style helpers and monoid helpers, while Option/Result functor and monad-style combinators are validated through updated examples and tests.

## New Features

### Language Features
- Improved functional programming workflows across Option/Result usage with validated examples.
- Lambda expressions are demonstrated in first-class function flows via `call1(...)` usage.

### Standard Library
- Added `func::apply(f, x)` for applicative-style function application.
- Added `func::ap(f, x)` as an alias of `apply`.
- Added `func::mempty()` (additive integer identity).
- Added `func::mappend(a, b)` (additive integer combine).

### Compiler Features
- LLVM dispatch and codegen now handle new `func` module primitives: `apply`, `ap`, `mempty`, and `mappend`.

## Improvements

### Performance
- No major performance-focused changes in this release.

### Code Generation
- Added dedicated LLVM generators for new functional helpers in `func` module.
- Improved functional coverage in real-world examples used for regression checking.

### Error Messages
- No major diagnostics-focused changes in this release.

### Documentation
- Updated functional examples to include Applicative and Monoid usage.
- Refreshed optional advanced example text to align with implemented capabilities.

## Bug Fixes

- Fixed release validation gaps by adding runtime checks for functional combinations used in examples.
- Resolved stale-binary validation confusion by standardizing release verification on a fresh `go build` output.

## Breaking Changes

None.

## API Changes

- Added new `func` module APIs:
	- `func::apply(f, x)`
	- `func::ap(f, x)`
	- `func::mempty()`
	- `func::mappend(a, b)`

## Deprecations

None.

## Technical Details

- Added and wired LLVM functional generation paths for new stdlib operations.
- Expanded `func` module unit tests to verify registration and arity for all functional helpers.
- Updated functional example programs to exercise:
	- Functor style: `.map(...)`
	- Monad style: `.and_then(...)`
	- Applicative style: `apply(...)`, `ap(...)`
	- Monoid style: `mempty()`, `mappend(...)`

## Dependencies

No dependency changes.

## Example Usage

```lotus
use "io";
use "func";

fn int double(int x) { ret x * 2; }

fn int main() {
		int a = Some(21).map(fn double).unwrap();
		int b = Some(21).and_then(fn double).unwrap_or(0);
		int c = func::apply(fn double, 11);
		int z = func::mempty();
		int s = func::mappend(7, 8);
		printf("%d %d %d %d %d\n", a, b, c, z, s);
		ret 0;
}
```

## Contributors

- Core Lotus maintainers and contributors across parser, LLVM codegen, and stdlib updates.

## Installation

### From Source
```bash
git clone https://github.com/j-alexander3375/Lotus
cd Lotus
git checkout v1.9.0
cd src
go build -o ../lotus .
```

### Arch Linux (AUR)
```bash
yay -S lotus-lang
# or
paru -S lotus-lang
```

### Binary Downloads
Download pre-built binaries from the [releases page](https://github.com/j-alexander3375/Lotus/releases/tag/v1.9.0).

## Checksums

SHA256 checksums for release artifacts are available in the release assets.
