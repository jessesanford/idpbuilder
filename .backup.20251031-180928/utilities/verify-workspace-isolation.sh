#!/bin/bash

# 🚨 WORKSPACE ISOLATION VERIFICATION UTILITY
# Part of Software Factory 2.0 - RULE R176, R177, R178
# 
# Purpose: Verify that agents are working in isolated effort directories
# Usage: ./verify-workspace-isolation.sh [effort-name]
#
# Returns:
#   0 - Workspace properly isolated
#   1 - Workspace violation detected

set -e

echo "═══════════════════════════════════════════════════════════════"
echo "🔍 WORKSPACE ISOLATION VERIFICATION"
echo "═══════════════════════════════════════════════════════════════"
echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo "═══════════════════════════════════════════════════════════════"

# Function to check if in effort directory
check_effort_directory() {
    local current_dir=$(pwd)
    
    if [[ "$current_dir" != *"/efforts/phase"*"/wave"*"/"* ]]; then
        echo "❌ FAIL: Not in effort directory"
        echo "   Current: $current_dir"
        echo "   Expected: */efforts/phase*/wave*/[effort-name]"
        return 1
    fi
    
    echo "✅ PASS: In effort directory"
    echo "   Location: $current_dir"
    return 0
}

# Function to check if pkg directory exists
check_pkg_directory() {
    if [ ! -d "./pkg" ]; then
        echo "❌ FAIL: No isolated pkg directory"
        echo "   All code must be in ./pkg/"
        return 1
    fi
    
    echo "✅ PASS: Isolated pkg directory exists"
    echo "   Location: $(pwd)/pkg"
    return 0
}

# Function to check main pkg is empty
check_main_pkg_empty() {
    local project_root=$(git rev-parse --show-toplevel 2>/dev/null || echo "/workspaces/$(basename $(pwd))")
    local main_pkg="$project_root/pkg"
    
    if [ -d "$main_pkg" ] && [ -n "$(ls -A $main_pkg 2>/dev/null)" ]; then
        echo "⚠️ WARNING: Main /pkg contains files during implementation"
        echo "   This may indicate workspace violation"
        echo "   Main pkg should be empty until integration phase"
        
        # List what's in main pkg
        echo "   Files found in main pkg:"
        find "$main_pkg" -type f -name "*.go" 2>/dev/null | head -5
    else
        echo "✅ PASS: Main /pkg is properly empty"
    fi
}

# Function to verify effort structure
verify_effort_structure() {
    local required_files=(
        "IMPLEMENTATION-PLAN.md"
        "work-log.md"
    )
    
    echo ""
    echo "Checking effort structure..."
    
    for file in "${required_files[@]}"; do
        if [ -f "./$file" ]; then
            echo "✅ Found: $file"
        else
            echo "⚠️ Missing: $file"
        fi
    done
}

# Function to count code in effort
count_effort_code() {
    if [ -d "./pkg" ]; then
        local go_files=$(find ./pkg -name "*.go" -not -name "*_test.go" 2>/dev/null | wc -l)
        local test_files=$(find ./pkg -name "*_test.go" 2>/dev/null | wc -l)
        
        echo ""
        echo "Code in effort:"
        echo "  Production files: $go_files"
        echo "  Test files: $test_files"
        
        if [ $go_files -eq 0 ]; then
            echo "⚠️ No production code found yet"
        fi
    fi
}

# Main verification flow
main() {
    local effort_name="${1:-}"
    local exit_code=0
    
    if [ -n "$effort_name" ]; then
        echo "Verifying effort: $effort_name"
        echo ""
    fi
    
    # Run all checks
    if ! check_effort_directory; then
        exit_code=1
    fi
    
    if ! check_pkg_directory; then
        exit_code=1
    fi
    
    check_main_pkg_empty
    verify_effort_structure
    count_effort_code
    
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    
    if [ $exit_code -eq 0 ]; then
        echo "✅ WORKSPACE ISOLATION: VERIFIED"
        echo "Agent is working in proper isolated directory"
    else
        echo "❌ WORKSPACE ISOLATION: VIOLATED"
        echo "CRITICAL: This is a GRADING FAILURE (20% lost)"
        echo "Agent must work in isolated effort directory!"
    fi
    
    echo "═══════════════════════════════════════════════════════════════"
    
    return $exit_code
}

# Run main function
main "$@"