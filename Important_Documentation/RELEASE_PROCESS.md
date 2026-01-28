# Lotus Release Process

This document describes the release process for Lotus, including version management, automated builds, and distribution.

## Table of Contents

1. [Overview](#overview)
2. [Version Numbering](#version-numbering)
3. [Release Workflow](#release-workflow)
4. [Automated Version Bumping](#automated-version-bumping)
5. [Creating a Release](#creating-a-release)
6. [Post-Release Tasks](#post-release-tasks)
7. [Troubleshooting](#troubleshooting)

## Overview

Lotus uses GitHub Actions for automated releases, including:
- Multi-platform binary builds (Linux, macOS, Windows)
- Source tarballs for package managers
- Automated checksums
- Release notes integration

## Version Numbering

Lotus follows [Semantic Versioning](https://semver.org/) (SemVer):

```
MAJOR.MINOR.PATCH
```

- **MAJOR**: Incompatible API changes
- **MINOR**: New features, backwards-compatible
- **PATCH**: Bug fixes, backwards-compatible

### Examples
- `1.5.3` → `1.5.4`: Patch release (bug fixes)
- `1.5.4` → `1.6.0`: Minor release (new features)
- `1.6.0` → `2.0.0`: Major release (breaking changes)

## Release Workflow

### Architecture

```
Developer              GitHub Actions                 Distributions
    |                        |                              |
    | git tag v1.5.5         |                              |
    |----------------------->|                              |
    |                        | Build binaries               |
    |                        | (Linux, macOS, Windows)      |
    |                        |                              |
    |                        | Create GitHub Release        |
    |                        |                              |
    |                        | Upload artifacts             |
    |                        |------------------------------>
    |                        |                              | AUR update
    |                        |                              | brew update
    |                        |                              | etc.
```

### Triggered By

The release workflow is triggered when a tag matching these patterns is pushed:
- `v*.*.*` (e.g., `v1.5.5`)
- `*.*.*` (e.g., `1.5.5`)

## Automated Version Bumping

Use the `scripts/bump_version.sh` script to automate version updates.

### Usage

```bash
# Dry run (preview changes)
./scripts/bump_version.sh patch --dry-run

# Bump patch version (1.5.4 → 1.5.5)
./scripts/bump_version.sh patch

# Bump minor version (1.5.5 → 1.6.0)
./scripts/bump_version.sh minor

# Bump major version (1.6.0 → 2.0.0)
./scripts/bump_version.sh major
```

### What It Does

The script automatically:
1. Updates `src/constants.go` (CompilerVersion)
2. Updates `PKGBUILD` (pkgver)
3. Updates `.SRCINFO` (pkgver)
4. Creates `RELEASE_NOTES_X.Y.Z.md` from template
5. Shows next steps for completing the release

## Creating a Release

### Step-by-Step Process

#### 1. Prepare for Release

```bash
# Ensure you're on master and up to date
git checkout master
git pull origin master

# Run tests
cd src
go test -timeout 0 ./...
cd ..

# Build and test examples
./lotus --run examples/control_flow_if.lts
./lotus --run examples/control_flow_for.lts
```

#### 2. Bump Version

```bash
# Preview changes
./scripts/bump_version.sh patch --dry-run

# Execute version bump
./scripts/bump_version.sh patch
```

#### 3. Edit Release Notes

Open the generated `RELEASE_NOTES_X.Y.Z.md` and fill in:
- New features
- Improvements
- Bug fixes
- Breaking changes
- Examples

```bash
# Example for 1.5.5
vi RELEASE_NOTES_1.5.5.md
```

#### 4. Commit and Tag

```bash
# Review changes
git diff

# Commit version bump
git add -A
git commit -m "Bump version to X.Y.Z"

# Create and push tag
git tag vX.Y.Z
git push origin master --tags
```

**Important**: Pushing the tag triggers the GitHub Actions release workflow.

#### 5. Monitor Release Build

1. Go to GitHub Actions tab: https://github.com/j-alexander3375/Lotus/actions
2. Watch the "Release" workflow
3. Verify all builds complete successfully
4. Check that release is created with all artifacts

#### 6. Update PKGBUILD Checksums

After the GitHub release is created:

```bash
# Update checksums (requires arch-linux tools)
updpkgsums

# Regenerate .SRCINFO
makepkg --printsrcinfo > .SRCINFO

# Commit updates
git add PKGBUILD .SRCINFO
git commit -m "Update PKGBUILD checksums for vX.Y.Z"
git push origin master
```

#### 7. Update AUR (if applicable)

```bash
cd ~/aur-lotus-lang  # Your AUR repository
cp /path/to/lotus/PKGBUILD .
cp /path/to/lotus/.SRCINFO .

# Test build
makepkg -si

# Commit and push to AUR
git add PKGBUILD .SRCINFO
git commit -m "Update to vX.Y.Z"
git push
```

## Post-Release Tasks

### Verify Release

1. **Check GitHub Release**: Visit https://github.com/j-alexander3375/Lotus/releases
   - Verify release notes are displayed
   - Check all binary artifacts are uploaded
   - Verify checksums file is present

2. **Test Binary Downloads**:
   ```bash
   # Download and test a binary
   wget https://github.com/j-alexander3375/Lotus/releases/download/vX.Y.Z/lotus-X.Y.Z-linux-amd64.tar.gz
   tar xzf lotus-X.Y.Z-linux-amd64.tar.gz
   cd lotus-X.Y.Z-linux-amd64
   ./lotus --version
   ```

3. **Test Source Tarball**:
   ```bash
   wget https://github.com/j-alexander3375/Lotus/releases/download/vX.Y.Z/lotus-lang-X.Y.Z.tar.gz
   tar xzf lotus-lang-X.Y.Z.tar.gz
   cd Lotus-X.Y.Z/src
   go build -o ../lotus .
   ../lotus --version
   ```

### Announce Release

1. Update project README if needed
2. Post announcement on relevant forums/communities
3. Update documentation website (if applicable)
4. Tweet/social media announcement

### Update Development Branch

```bash
# Create next development version marker
git checkout -b develop
# Continue development...
```

## Troubleshooting

### Build Failures

**Problem**: GitHub Actions build fails

**Solutions**:
1. Check the Actions log for specific errors
2. Common issues:
   - LLVM installation failure: Check LLVM version compatibility
   - Go build failure: Verify Go version in workflow
   - Test failure: Run tests locally first

### Wrong Version in Binary

**Problem**: Compiled binary shows wrong version

**Solution**:
```bash
# Verify constants.go was updated
grep CompilerVersion src/constants.go

# Rebuild
cd src
go build -o ../lotus .
../lotus --version
```

### Missing Release Assets

**Problem**: Some binaries didn't upload

**Solution**:
1. Check if specific platform build failed
2. Re-run failed jobs in GitHub Actions
3. Manually build and upload if necessary:
   ```bash
   cd src
   GOOS=linux GOARCH=amd64 go build -o ../lotus-linux-amd64 .
   # Upload to release manually via GitHub UI
   ```

### PKGBUILD Checksum Mismatch

**Problem**: AUR users report checksum mismatch

**Solution**:
```bash
# Re-download tarball and regenerate
rm lotus-lang-*.tar.gz
wget https://github.com/j-alexander3375/Lotus/archive/refs/tags/vX.Y.Z.tar.gz -O lotus-lang-X.Y.Z.tar.gz
updpkgsums
makepkg --printsrcinfo > .SRCINFO
git commit -am "Fix checksums for vX.Y.Z"
```

## Workflow Files

### `.github/workflows/release.yml`
- Triggered by version tags
- Builds multi-platform binaries
- Creates GitHub release
- Uploads artifacts

### `.github/workflows/ci.yml`
- Runs on every push/PR
- Tests on multiple platforms
- Linting and formatting checks
- Build verification

## Version Files

- `src/constants.go`: Compiler version constant
- `PKGBUILD`: Arch Linux package version
- `.SRCINFO`: AUR metadata
- `RELEASE_NOTES_X.Y.Z.md`: Release-specific notes

## Best Practices

1. **Always test before releasing**:
   - Run full test suite
   - Build and test examples
   - Test on multiple platforms if possible

2. **Write clear release notes**:
   - List all significant changes
   - Include examples for new features
   - Document breaking changes prominently
   - Credit contributors

3. **Use semantic versioning**:
   - Patch: Bug fixes only
   - Minor: New features, no breaking changes
   - Major: Breaking changes

4. **Verify builds**:
   - Check all platforms build successfully
   - Download and test at least one binary
   - Verify version string is correct

5. **Keep documentation updated**:
   - Update README if needed
   - Update DEVELOPMENT.md for new features
   - Keep examples current

## Quick Reference

```bash
# Complete release checklist
□ Run tests: cd src && go test -timeout 0 ./...
□ Bump version: ./scripts/bump_version.sh [patch|minor|major]
□ Edit release notes: vi RELEASE_NOTES_X.Y.Z.md
□ Review changes: git diff
□ Commit: git add -A && git commit -m "Bump version to X.Y.Z"
□ Tag: git tag vX.Y.Z
□ Push: git push origin master --tags
□ Monitor GitHub Actions
□ Verify release artifacts
□ Update PKGBUILD: updpkgsums && makepkg --printsrcinfo > .SRCINFO
□ Update AUR (if applicable)
□ Announce release
```

## Support

For issues with the release process:
- GitHub Issues: https://github.com/j-alexander3375/Lotus/issues
- Check GitHub Actions logs for build failures
- Review this document for common solutions
