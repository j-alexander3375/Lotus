# Lotus v1.1.0 - Hash Module Implementation

## 🎉 What's New

### Hash Module (`hash`)
A complete hashing module with both cryptographic and non-cryptographic hash functions:

- **CRC32**: IEEE 802.3 polynomial-based checksum with optimized 256-entry lookup table
- **FNV-1a**: 64-bit Fowler-Noll-Vo hash algorithm for fast, non-cryptographic hashing
- **DJB2**: Simple and efficient string hashing algorithm
- **MurmurHash3**: 32-bit hash with seed support for consistent hashing

### Usage Example
```lotus
use "hash"

fn int main() {
    const string data = "Hello, World!";
    
    int crc = crc32(data, 13);      // Checksum validation
    int fnv = fnv1a(data, 13);      // Fast non-crypto hash
    int djb = djb2(data);           // Simple string hash
    int murmur = murmur(data, 13, 42);  // Seeded hash
    
    ret 0;
}
```

## 🔧 Infrastructure Improvements

- Added label tracking system (`hasLabel`, `markLabel`) for efficient code generation
- Enhanced code generator with reusable label management
- Improved stdlib module system

## 📚 Documentation

- Updated DEVELOPMENT.md with Phase 4 goals and progress tracking
- Comprehensive hash module documentation
- Test coverage for all hash functions

## 🧪 Testing

- `test_hash.lts` - Basic hash function testing
- `test_hash_full.lts` - Comprehensive hash module tests
- All unit tests passing ✅

## 📦 Installation

### Arch Linux (AUR)
```bash
yay -S lotus-lang
```

### From Source
```bash
git clone https://github.com/j-alexander3375/Lotus
cd Lotus/src
go build -o ../lotus .
```

## 🔗 Links

- [Documentation](https://github.com/j-alexander3375/Lotus/tree/master/Important_Documentation)
- [Examples](https://github.com/j-alexander3375/Lotus/tree/master/examples)
- [Test Suite](https://github.com/j-alexander3375/Lotus/tree/master/tests)

## 📊 Stats

- 4 new hash functions implemented
- 200+ lines of optimized x86-64 assembly generation
- Full test coverage

---

**Full Changelog**: https://github.com/j-alexander3375/Lotus/compare/0.1.3...1.1.0
