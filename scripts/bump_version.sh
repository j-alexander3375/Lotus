#!/bin/bash

# bump_version.sh - Automated version bumping script for Lotus
# Usage: ./scripts/bump_version.sh <major|minor|patch> [--dry-run]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Files to update
CONSTANTS_FILE="src/constants.go"
PKGBUILD_FILE="PKGBUILD"
SRCINFO_FILE=".SRCINFO"

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to extract current version
get_current_version() {
    if [ ! -f "$CONSTANTS_FILE" ]; then
        print_error "Constants file not found: $CONSTANTS_FILE"
        exit 1
    fi

    grep 'CompilerVersion = ' "$CONSTANTS_FILE" | sed 's/.*"\(.*\)".*/\1/'
}

# Function to parse version into components
parse_version() {
    local version=$1
    IFS='.' read -r -a parts <<< "$version"
    MAJOR="${parts[0]}"
    MINOR="${parts[1]}"
    PATCH="${parts[2]}"
}

# Function to bump version
bump_version() {
    local bump_type=$1
    local current_version=$2

    parse_version "$current_version"

    case "$bump_type" in
        major)
            MAJOR=$((MAJOR + 1))
            MINOR=0
            PATCH=0
            ;;
        minor)
            MINOR=$((MINOR + 1))
            PATCH=0
            ;;
        patch)
            PATCH=$((PATCH + 1))
            ;;
        *)
            print_error "Invalid bump type: $bump_type"
            echo "Usage: $0 <major|minor|patch> [--dry-run]"
            exit 1
            ;;
    esac

    echo "${MAJOR}.${MINOR}.${PATCH}"
}

# Function to update file
update_file() {
    local file=$1
    local old_version=$2
    local new_version=$3
    local dry_run=$4

    if [ ! -f "$file" ]; then
        print_warn "File not found: $file (skipping)"
        return
    fi

    if [ "$dry_run" = true ]; then
        print_info "Would update $file: $old_version -> $new_version"
        return
    fi

    # Create backup
    cp "$file" "${file}.bak"

    # Update the file
    sed -i "s/$old_version/$new_version/g" "$file"

    print_info "Updated $file: $old_version -> $new_version"
}

# Function to create release notes file
create_release_notes() {
    local version=$1
    local dry_run=$2

    local notes_file="RELEASE_NOTES_${version}.md"

    if [ -f "$notes_file" ]; then
        print_warn "Release notes already exist: $notes_file"
        return
    fi

    if [ "$dry_run" = true ]; then
        print_info "Would create release notes file: $notes_file"
        return
    fi

    # Copy template and replace version
    if [ -f ".github/RELEASE_TEMPLATE.md" ]; then
        sed "s/{VERSION}/$version/g" .github/RELEASE_TEMPLATE.md > "$notes_file"
        print_info "Created release notes file: $notes_file"
        print_info "Please edit $notes_file with release details"
    else
        print_warn "Release template not found, creating basic release notes"
        cat > "$notes_file" <<EOF
# Release Notes for Lotus v${version}

## New Features

<!-- Add new features here -->

## Improvements

<!-- Add improvements here -->

## Bug Fixes

<!-- Add bug fixes here -->

## Breaking Changes

<!-- Add breaking changes here -->
EOF
        print_info "Created release notes file: $notes_file"
    fi
}

# Function to update PKGBUILD sha256sum
update_pkgbuild_sha() {
    local version=$1
    local dry_run=$2

    if [ ! -f "$PKGBUILD_FILE" ]; then
        print_warn "PKGBUILD not found (skipping)"
        return
    fi

    if [ "$dry_run" = true ]; then
        print_info "Would update PKGBUILD sha256sum for version $version"
        print_warn "Remember to run: updpkgsums after creating the GitHub release"
        return
    fi

    print_warn "PKGBUILD updated with new version"
    print_warn "After creating the GitHub release, run: updpkgsums"
    print_warn "Then update .SRCINFO with: makepkg --printsrcinfo > .SRCINFO"
}

# Main script
main() {
    local bump_type=$1
    local dry_run=false

    # Check for dry-run flag
    if [ "$2" = "--dry-run" ]; then
        dry_run=true
        print_info "Running in DRY-RUN mode (no files will be modified)"
    fi

    # Validate arguments
    if [ -z "$bump_type" ]; then
        print_error "Missing version bump type"
        echo "Usage: $0 <major|minor|patch> [--dry-run]"
        exit 1
    fi

    # Get current version
    current_version=$(get_current_version)
    print_info "Current version: $current_version"

    # Calculate new version
    new_version=$(bump_version "$bump_type" "$current_version")
    print_info "New version: $new_version"

    # Update files
    update_file "$CONSTANTS_FILE" "$current_version" "$new_version" $dry_run
    update_file "$PKGBUILD_FILE" "$current_version" "$new_version" $dry_run
    update_file "$SRCINFO_FILE" "$current_version" "$new_version" $dry_run

    # Create release notes
    create_release_notes "$new_version" $dry_run

    # Update PKGBUILD checksums
    update_pkgbuild_sha "$new_version" $dry_run

    if [ "$dry_run" = false ]; then
        print_info ""
        print_info "Version bump complete! Next steps:"
        print_info "1. Edit RELEASE_NOTES_${new_version}.md with release details"
        print_info "2. Review changes: git diff"
        print_info "3. Commit changes: git add -A && git commit -m 'Bump version to $new_version'"
        print_info "4. Create and push tag: git tag v${new_version} && git push origin master --tags"
        print_info "5. After GitHub release is created, update PKGBUILD:"
        print_info "   - Run: updpkgsums"
        print_info "   - Run: makepkg --printsrcinfo > .SRCINFO"
        print_info "   - Commit: git add PKGBUILD .SRCINFO && git commit -m 'Update PKGBUILD checksums for v${new_version}'"
    fi
}

main "$@"
