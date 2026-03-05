#!/bin/bash

# fix_test_files.sh - Fix known issues in test .lts files

set -e

REPO_DIR="/mnt/c/Users/joshu/develLotus"

log_info() {
    echo "[FIX] $1"
}

cd "$REPO_DIR"

log_info "Fixing const_namespace_test.lts - add namespace qualifiers"
if [ -f "examples/const_namespace_test.lts" ]; then
    sed -i 's/float area = circle_area/float area = math::circle_area/' examples/const_namespace_test.lts
    sed -i 's/int max = get_max/int max = math::get_max/' examples/const_namespace_test.lts
    log_info "✓ Fixed const_namespace_test.lts"
fi

log_info "Fixing templates_simple_test.lts - add namespace qualifiers"
if [ -f "examples/templates_simple_test.lts" ]; then
    # Check if multiply needs namespace qualifier
    if grep -q "fn.*multiply" examples/templates_simple_test.lts; then
        if ! grep -q "::" examples/templates_simple_test.lts; then
            # Add namespace qualifier if calling multiply without it
            sed -i 's/int result = multiply/int result = math::multiply/g' examples/templates_simple_test.lts
            log_info "✓ Fixed templates_simple_test.lts"
        fi
    fi
fi

log_info "Checking for other namespace-related issues..."
for file in examples/*.lts tests/*.lts; do
    if [ -f "$file" ]; then
        # Check for function calls that might need namespace qualifiers
        if grep -q "Compilation failed.*undefined function" "$REPO_DIR/test_results.log" 2>/dev/null; then
            func_name=$(grep "undefined function" "$REPO_DIR/test_results.log" | grep "$file" | head -1 | sed 's/.*undefined function: \([a-zA-Z_][a-zA-Z0-9_]*\).*/\1/')
            if [ -n "$func_name" ]; then
                log_info "File $file may need namespace qualifier for: $func_name"
            fi
        fi
    fi
done

log_info "Test file fixes complete"
