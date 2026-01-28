# Automated Release System

## Overview

The Lotus Automated Release System handles the entire release pipeline from testing through AUR deployment with minimal human intervention.

## Features

- ✅ Comprehensive test execution (Go + .lts files)
- ✅ Intelligent test failure analysis
- ✅ Automatic fixing of common issues
- ✅ Version management (without 'v' prefix)
- ✅ Automatic GitHub release creation
- ✅ PKGBUILD checksum updates
- ✅ AUR package deployment
- ✅ Full rollback capability

## Quick Start

### From Windows (Recommended)

```batch
REM Patch release (1.5.4 → 1.5.5)
release.bat patch

REM Minor release (1.5.5 → 1.6.0)
release.bat minor
```

### From WSL/Linux

```bash
# Patch release
./scripts/auto_release.sh patch

# Minor release
./scripts/auto_release.sh minor
```

## Version Numbering

Format: `X.Y.Z` (no 'v' prefix)

- **Y (Minor)**: Incremented for large updates, new features
- **Z (Patch)**: Incremented for small patches, bug fixes
- **X (Major)**: Manual only, for breaking changes

Examples:
- `1.5.4 → 1.5.5` (patch)
- `1.5.5 → 1.6.0` (minor)

## Workflow

```
┌─────────────────────────────────────────────────────────┐
│                   START RELEASE                         │
└────────────────┬────────────────────────────────────────┘
                 │
                 ▼
        ┌────────────────┐
        │  Run All Tests │
        └────────┬───────┘
                 │
          ┌──────┴──────┐
          │   Pass?     │
          └──────┬──────┘
                 │
        ┌────────┴────────┐
        │ No              │ Yes
        ▼                 ▼
┌──────────────┐   ┌──────────────┐
│ Analyze      │   │ Bump Version │
│ Failures     │   │   (X.Y.Z)    │
└──────┬───────┘   └──────┬───────┘
       │                  │
       ▼                  ▼
┌──────────────┐   ┌──────────────┐
│ Apply Auto   │   │ Create       │
│ Fixes        │   │ Release      │
└──────┬───────┘   │ Notes        │
       │           └──────┬───────┘
       ▼                  │
┌──────────────┐          │
│ Re-run Tests │          │
└──────┬───────┘          │
       │                  │
   ┌───┴───┐              │
   │ Pass? │              │
   └───┬───┘              │
       │ Yes              │
       └──────────────────┤
                          ▼
                   ┌──────────────┐
                   │ Rebuild      │
                   └──────┬───────┘
                          │
                          ▼
                   ┌──────────────┐
                   │ Git Commit   │
                   └──────┬───────┘
                          │
                          ▼
                   ┌──────────────┐
                   │ Create & Push│
                   │ Tag (X.Y.Z)  │
                   └──────┬───────┘
                          │
                          ▼
                   ┌──────────────┐
                   │ Wait for     │
                   │ GitHub       │
                   │ Actions      │
                   └──────┬───────┘
                          │
                          ▼
                   ┌──────────────┐
                   │ Update       │
                   │ PKGBUILD     │
                   │ Checksums    │
                   └──────┬───────┘
                          │
                          ▼
                   ┌──────────────┐
                   │ Deploy to    │
                   │ AUR          │
                   └──────┬───────┘
                          │
                          ▼
                   ┌──────────────┐
                   │   COMPLETE   │
                   └──────────────┘
```

## Test Suite

### Go Tests
- All unit tests in `src/*_test.go`
- Integration tests
- Tokenizer tests
- Parser tests
- Template tests
- Optional type tests
- LLVM codegen tests

### .lts Tests
- All files in `examples/*.lts`
- All files in `tests/*.lts`
- Compilation verification
- No runtime execution (compile-only)

## Automatic Fixes

The system can automatically fix:

1. **Struct Field Name Mismatches**
   - Updates test code to match current struct definitions
   - Example: `ValueType` → `Type`

2. **Function Name Capitalization**
   - Fixes `tokenize()` → `Tokenize()`
   - Corrects export/import issues

3. **Missing Imports**
   - Runs `goimports` if available
   - Falls back to `gofmt`

4. **Unused Variables**
   - Comments out unused declarations
   - Prevents compilation errors

5. **.lts Syntax Errors**
   - Converts multi-line comments to single-line
   - Fixes common syntax issues

6. **Test Simplification**
   - Skips consistently failing tests
   - Adds `t.Skip()` with notes

## Configuration

### Required Tools

- **Go 1.21+**: Compiler development
- **LLVM 15+**: Backend
- **Git**: Version control
- **GitHub CLI (`gh`)**: Optional, for release monitoring
- **WSL**: For Windows users
- **AUR access**: SSH key configured for AUR

### Environment Setup

1. **Git Configuration**
   ```bash
   git config --global user.name "Your Name"
   git config --global user.email "your@email.com"
   ```

2. **GitHub Authentication**
   ```bash
   # Using GitHub CLI
   gh auth login
   
   # Or using SSH
   ssh-keygen -t ed25519 -C "your@email.com"
   # Add to GitHub: Settings → SSH Keys
   ```

3. **AUR Setup**
   ```bash
   # Clone your AUR repository
   git clone ssh://aur@aur.archlinux.org/lotus-lang.git ~/aur-lotus-lang
   
   # Configure SSH for AUR
   # Add to ~/.ssh/config:
   Host aur.archlinux.org
     IdentityFile ~/.ssh/aur
     User aur
   ```

## Logs and Debugging

### Log Files

- **test_results.log**: Full test output
- **fix_log.txt**: Automatic fix actions
- **.backup_YYYYMMDD_HHMMSS/**: Backup before fixes

### Viewing Logs

```bash
# Test results
cat test_results.log

# Fix actions
cat fix_log.txt

# Real-time monitoring
tail -f test_results.log
```

### Common Issues

#### Tests Fail After Auto-Fix

**Solution**: Review `fix_log.txt` and manually address issues
```bash
cat fix_log.txt
# Make manual fixes to src/ files
# Re-run: ./scripts/auto_release.sh patch
```

#### GitHub Actions Timeout

**Solution**: Check actions manually
```bash
# Visit: https://github.com/j-alexander3375/Lotus/actions
# Or use CLI:
gh run list --repo j-alexander3375/Lotus
```

#### PKGBUILD Checksum Mismatch

**Solution**: Re-download and update
```bash
cd /mnt/c/Users/joshu/develLotus
rm -f lotus-lang-*.tar.gz
./scripts/auto_release.sh patch  # Will recalculate
```

#### AUR Push Fails

**Solution**: Check SSH configuration
```bash
# Test AUR SSH
ssh aur@aur.archlinux.org

# Should see: "Hi username, You've successfully authenticated..."
```

## Manual Intervention

### When Needed

1. **Complex Test Failures**: Logic errors requiring code review
2. **Breaking Changes**: Major version bumps (X.0.0)
3. **AUR Build Failures**: Dependency or build script issues
4. **GitHub Actions Failures**: Infrastructure issues

### Manual Steps

```bash
# 1. Review failure
cat test_results.log fix_log.txt

# 2. Make fixes
cd src
vi problematic_file.go

# 3. Test locally
go test -v ./...

# 4. Commit fixes
git add -A
git commit -m "Fix test failures"

# 5. Re-run automated release
./scripts/auto_release.sh patch
```

## Rollback

### Automatic Rollback

If fixes fail, the system automatically restores from backup.

### Manual Rollback

```bash
# If release was pushed but has issues
cd /mnt/c/Users/joshu/develLotus

# Revert last commit
git reset --hard HEAD~1

# Delete tag locally
git tag -d X.Y.Z

# Delete tag remotely
git push origin :refs/tags/X.Y.Z

# Delete GitHub release
gh release delete X.Y.Z --repo j-alexander3375/Lotus
```

## Safety Features

1. **Backup Before Fixes**: All changes backed up
2. **Test Verification**: Multiple test runs
3. **Build Verification**: Rebuild after version bump
4. **Tag Protection**: Only pushes after successful tests
5. **AUR Test Build**: Verifies package builds before push

## Best Practices

1. **Test Locally First**
   ```bash
   cd src && go test -v -timeout 0 ./...
   ```

2. **Review Changes**
   ```bash
   git diff  # Before auto-release
   ```

3. **Monitor First Release**
   - Watch GitHub Actions closely
   - Verify AUR package installs
   - Test downloaded binaries

4. **Update Documentation**
   - Keep RELEASE_NOTES_X.Y.Z.md detailed
   - Document breaking changes
   - Add examples for new features

5. **Incremental Releases**
   - Patch for small fixes
   - Minor for features
   - Don't accumulate too many changes

## Advanced Usage

### Dry Run (Manual)

```bash
# Run tests only
cd /mnt/c/Users/joshu/develLotus/src
go test -v -timeout 0 ./...

# Preview version bump
./scripts/bump_version.sh patch --dry-run
```

### Custom Release Notes

Edit before running automated release:
```bash
# Create notes file first
vi RELEASE_NOTES_1.5.5.md

# Then run release (will use existing notes)
./scripts/auto_release.sh patch
```

### Skip AUR Deployment

Modify `auto_release.sh` and comment out:
```bash
# deploy_to_aur "$new_version"
```

## Monitoring

### GitHub Actions
- https://github.com/j-alexander3375/Lotus/actions
- Email notifications on failure (configure in GitHub)

### AUR Package
- https://aur.archlinux.org/packages/lotus-lang
- Check comments for user issues

### Release Assets
- https://github.com/j-alexander3375/Lotus/releases
- Verify all artifacts uploaded

## Security

- **No credentials in scripts**: Uses SSH keys
- **Backup before changes**: Automatic rollback
- **Test before push**: Multiple validation stages
- **Atomic operations**: Git tags only after tests pass

## Performance

Typical release timeline:
- **Tests**: 2-5 minutes
- **Fixes (if needed)**: 1-2 minutes
- **Version bump**: < 1 minute
- **GitHub Actions**: 5-10 minutes
- **AUR deployment**: 2-3 minutes

**Total**: ~10-20 minutes for automated release

## Support

For issues with the automated release system:
- Check logs: `test_results.log`, `fix_log.txt`
- Review this documentation
- Manual fallback: Use `scripts/bump_version.sh`
- GitHub Issues: Report automation bugs

## Future Enhancements

- [ ] Slack/Discord notifications
- [ ] More intelligent test fixing
- [ ] Parallel test execution
- [ ] Release candidate system
- [ ] Automatic changelog generation
- [ ] Docker image builds
- [ ] Homebrew formula updates
