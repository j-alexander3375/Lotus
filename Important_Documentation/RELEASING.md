# Quick Release Guide

**For detailed instructions, see [RELEASE_PROCESS.md](Important_Documentation/RELEASE_PROCESS.md)**

## Quick Steps

### 1. Prepare
```bash
git checkout master && git pull
cd src && go test -timeout 0 ./... && cd ..
```

### 2. Bump Version
```bash
# Preview
./scripts/bump_version.sh patch --dry-run

# Execute (patch/minor/major)
./scripts/bump_version.sh patch
```

### 3. Edit Release Notes
```bash
vi RELEASE_NOTES_X.Y.Z.md  # Fill in details
```

### 4. Commit and Tag
```bash
git add -A
git commit -m "Bump version to X.Y.Z"
git tag vX.Y.Z
git push origin master --tags
```

### 5. Monitor Build
- Go to: https://github.com/j-alexander3375/Lotus/actions
- Wait for "Release" workflow to complete
- Verify all artifacts uploaded

### 6. Update PKGBUILD
```bash
updpkgsums
makepkg --printsrcinfo > .SRCINFO
git add PKGBUILD .SRCINFO
git commit -m "Update PKGBUILD checksums for vX.Y.Z"
git push origin master
```

### 7. Update AUR (if applicable)
```bash
cd ~/aur-lotus-lang
cp /path/to/lotus/{PKGBUILD,.SRCINFO} .
makepkg -si  # Test build
git add PKGBUILD .SRCINFO
git commit -m "Update to vX.Y.Z"
git push
```

## Version Types

- **patch**: Bug fixes (1.5.4 → 1.5.5)
- **minor**: New features (1.5.5 → 1.6.0)
- **major**: Breaking changes (1.6.0 → 2.0.0)

## Files Updated by bump_version.sh

- `src/constants.go` - CompilerVersion
- `PKGBUILD` - pkgver
- `.SRCINFO` - pkgver
- `RELEASE_NOTES_X.Y.Z.md` - Created from template

## Checklist

```
□ Tests pass
□ Version bumped
□ Release notes filled
□ Changes committed
□ Tag created and pushed
□ GitHub Actions succeeded
□ PKGBUILD updated
□ AUR updated (if applicable)
□ Release verified
□ Release announced
```

## Troubleshooting

- **Build fails**: Check GitHub Actions logs
- **Wrong version**: Verify `src/constants.go`
- **Checksum mismatch**: Run `updpkgsums` again after release is live

## Resources

- [Full Documentation](Important_Documentation/RELEASE_PROCESS.md)
- [GitHub Actions](https://github.com/j-alexander3375/Lotus/actions)
- [Releases](https://github.com/j-alexander3375/Lotus/releases)
