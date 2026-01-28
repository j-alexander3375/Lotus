# Release Notes for Lotus 1.6.0

## Highlights

Automated release 1.6.0 with comprehensive testing and fixes.

## Changes in This Release

### Improvements
- Automated release pipeline enhancements
- Comprehensive test suite validation
- Build system improvements

### Bug Fixes
- Various fixes identified through automated testing

## Testing
- All Go tests passed
- All .lts example files compile successfully
- LLVM backend verified
- Cross-platform compatibility confirmed

## Installation

### Binary Downloads
Download pre-built binaries from the [releases page](https://github.com/j-alexander3375/Lotus/releases/tag/1.6.0).

### From Source
```bash
git clone https://github.com/j-alexander3375/Lotus
cd Lotus
git checkout 1.6.0
cd src
go build -o ../lotus .
```

### Arch Linux (AUR)
```bash
yay -S lotus-lang
```

## Checksums
SHA256 checksums for release artifacts are available in the release assets.
