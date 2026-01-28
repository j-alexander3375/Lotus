# Quick Release Guide - Automated System

## One Command Release

### Windows

```batch
REM Patch release (bug fixes, small updates)
release.bat patch

REM Minor release (new features, larger updates)
release.bat minor
```

### Linux/WSL

```bash
# Patch release
./scripts/auto_release.sh patch

# Minor release
./scripts/auto_release.sh minor
```

## What Happens Automatically

1. ✅ **Runs all tests** (Go + .lts files)
2. ✅ **Fixes common issues** automatically
3. ✅ **Bumps version** (X.Y.Z format, no 'v')
4. ✅ **Creates release notes**
5. ✅ **Rebuilds compiler**
6. ✅ **Commits changes**
7. ✅ **Creates and pushes tag**
8. ✅ **Waits for GitHub Actions**
9. ✅ **Updates PKGBUILD checksums**
10. ✅ **Deploys to AUR**

## Version Format

- **1.5.4 → 1.5.5**: Patch (bug fixes)
- **1.5.5 → 1.6.0**: Minor (new features)

No 'v' prefix!

## Requirements

- WSL (for Windows users)
- Go 1.21+
- LLVM 15+
- Git configured
- GitHub authentication
- AUR SSH access

## Logs

- **test_results.log**: Test output
- **fix_log.txt**: Automatic fixes applied

## If Something Fails

The system will:
1. Show error in console
2. Keep backup of code
3. Log issues to fix_log.txt
4. Wait for manual fix

Then you can:
```bash
# View logs
cat test_results.log
cat fix_log.txt

# Make fixes
cd src
vi problematic_file.go

# Re-run
./scripts/auto_release.sh patch
```

## Monitor Progress

- **Console**: Live output
- **GitHub Actions**: https://github.com/j-alexander3375/Lotus/actions
- **Releases**: https://github.com/j-alexander3375/Lotus/releases
- **AUR**: https://aur.archlinux.org/packages/lotus-lang

## Safety

- Automatic backup before fixes
- Tests must pass before release
- Rollback on failure
- No credentials in scripts

## Full Documentation

See [AUTOMATED_RELEASE.md](Important_Documentation/AUTOMATED_RELEASE.md) for complete details.

## First Time Setup

```bash
# 1. Install GitHub CLI (optional but recommended)
# Ubuntu/Debian
sudo apt install gh
# or download from https://cli.github.com/

# 2. Authenticate
gh auth login

# 3. Configure AUR SSH
ssh-keygen -t ed25519 -C "your@email.com"
# Add public key to AUR account

# 4. Test
cd /mnt/c/Users/joshu/develLotus
./scripts/auto_release.sh patch
```

## That's It!

Just run `release.bat patch` and everything happens automatically! 🚀
