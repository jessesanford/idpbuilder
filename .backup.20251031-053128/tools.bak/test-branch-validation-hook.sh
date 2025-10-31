#!/bin/bash

# Test script for branch validation hook
# Tests various valid and invalid branch name patterns

set -euo pipefail

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m'

# Functions
print_test() {
    echo -e "${BLUE}Testing:${NC} $1"
}

print_pass() {
    echo -e "  ${GREEN}✅ PASS:${NC} $1"
}

print_fail() {
    echo -e "  ${RED}❌ FAIL:${NC} $1"
}

print_section() {
    echo ""
    echo -e "${BOLD}═══════════════════════════════════════════════════════${NC}"
    echo -e "${BOLD}$1${NC}"
    echo -e "${BOLD}═══════════════════════════════════════════════════════${NC}"
}

# Test the validation hook directly
test_branch_name() {
    local branch="$1"
    local expected="$2"  # "valid" or "invalid"
    local description="${3:-}"

    print_test "$branch"

    # Create a temporary test directory
    TEST_DIR=$(mktemp -d)
    cd "$TEST_DIR"

    # Initialize git repo
    git init --quiet
    git config user.email "test@test.com"
    git config user.name "Test User"

    # Create initial commit on main
    echo "test" > test.txt
    git add .
    git commit -m "initial" --quiet

    # Install the hook
    cp "$HOOK_SCRIPT" .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit

    # Try to create and commit on the test branch
    git checkout -b "$branch" --quiet 2>/dev/null || true
    echo "test2" > test2.txt
    git add .

    # Attempt commit and capture result
    if git commit -m "test commit" --quiet 2>/dev/null; then
        RESULT="valid"
    else
        RESULT="invalid"
    fi

    # Clean up
    cd - > /dev/null
    rm -rf "$TEST_DIR"

    # Check result
    if [ "$RESULT" = "$expected" ]; then
        print_pass "$description"
        return 0
    else
        print_fail "Expected $expected but got $RESULT - $description"
        return 1
    fi
}

# Main test execution
main() {
    # Check if hook script exists
    # Try new location first, fall back to old if needed
    if [ -f "$CLAUDE_PROJECT_DIR/tools/git-commit-hooks/effort-hooks/branch-name-validation.hook" ]; then
        HOOK_SCRIPT="$CLAUDE_PROJECT_DIR/tools/git-commit-hooks/effort-hooks/branch-name-validation.hook"
    elif [ -f "$CLAUDE_PROJECT_DIR/utilities/branch-name-validation-hook.sh" ]; then
        HOOK_SCRIPT="$CLAUDE_PROJECT_DIR/utilities/branch-name-validation-hook.sh"
    else
        HOOK_SCRIPT="/home/vscode/software-factory-template/tools/git-commit-hooks/effort-hooks/branch-name-validation.hook"
    fi
    if [ ! -f "$HOOK_SCRIPT" ]; then
        echo "❌ Hook script not found: $HOOK_SCRIPT"
        exit 1
    fi

    TOTAL_TESTS=0
    PASSED_TESTS=0
    FAILED_TESTS=0

    print_section "BRANCH NAME VALIDATION HOOK TEST SUITE"

    # Test valid patterns WITH project prefix
    print_section "Valid Patterns WITH Project Prefix"

    # Create test config with project prefix
    TEST_CONFIG=$(mktemp)
    cat > "$TEST_CONFIG" << EOF
branch_naming:
  project_prefix: "testproject"
EOF
    export CLAUDE_PROJECT_DIR=$(dirname "$TEST_CONFIG")
    cp "$TEST_CONFIG" "$CLAUDE_PROJECT_DIR/target-repo-config.yaml"

    # Valid effort branches with prefix
    test_branch_name "testproject/phase1/wave1/feature-a" "valid" "Effort with project prefix" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "testproject/phase2/wave3/complex-feature-name" "valid" "Complex effort name with prefix" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    # Valid split branches with prefix
    test_branch_name "testproject/phase1/wave1/feature-a--split-001" "valid" "Split 001 with prefix" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "testproject/phase2/wave1/api--split-999" "valid" "Split 999 with prefix" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    # Valid fix branches with prefix
    test_branch_name "testproject/phase1/wave1/feature-a-fix" "valid" "Fix branch with prefix" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    # Valid integration branches with prefix
    test_branch_name "testproject/phase1/wave1/integration" "valid" "Wave integration with prefix" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "testproject/phase1/integration" "valid" "Phase integration with prefix" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "testproject/integration" "valid" "Project integration with prefix" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    # Test valid patterns WITHOUT project prefix
    print_section "Valid Patterns WITHOUT Project Prefix"

    # Remove project prefix
    rm -f "$CLAUDE_PROJECT_DIR/target-repo-config.yaml"

    # Valid effort branches without prefix
    test_branch_name "phase1/wave1/feature-a" "valid" "Effort without prefix" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "phase10/wave99/long-feature-name" "valid" "Large phase/wave numbers" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    # Valid split branches without prefix
    test_branch_name "phase1/wave1/feature--split-001" "valid" "Split without prefix" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    # Valid integration branches without prefix
    test_branch_name "phase1/wave1/integration" "valid" "Wave integration without prefix" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "phase2/integration" "valid" "Phase integration without prefix" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    # Test special allowed branches
    print_section "Special Allowed Branches"

    test_branch_name "main" "valid" "Main branch" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "master" "valid" "Master branch" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "develop" "valid" "Develop branch" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "feature/new-feature" "valid" "Traditional feature branch" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "hotfix/urgent-fix" "valid" "Hotfix branch" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "release/v1.0.0" "valid" "Release branch" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    # Test INVALID patterns
    print_section "Invalid Branch Patterns (Should Be Blocked)"

    test_branch_name "phaseX/waveY/effort" "invalid" "Letter instead of number in phase" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "phase1/wavey/effort" "invalid" "Typo: wavey instead of wave" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "phase1/wave1/integration/repo" "invalid" "Extra path after integration" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "phase1/wave1/effort/effort" "invalid" "Duplicate effort path" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "random-branch-name" "invalid" "Random branch name" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "phase1-wave1-effort" "invalid" "Dashes instead of slashes" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "phase1/effort" "invalid" "Missing wave" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "wave1/effort" "invalid" "Missing phase" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "testproject/phase1/wave1" "invalid" "Missing effort name" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    # Test edge cases
    print_section "Edge Cases"

    test_branch_name "phase1/wave1/feature--split-000" "invalid" "Split with 000 (should start at 001)" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "phase1/wave1/feature--split-1" "invalid" "Split without 3 digits" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "phase1/wave1/UPPERCASE" "invalid" "Uppercase effort name" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    test_branch_name "phase1/wave1/feature_with_underscore" "invalid" "Underscore in effort name" && ((PASSED_TESTS++)) || ((FAILED_TESTS++))
    ((TOTAL_TESTS++))

    # Clean up
    rm -rf "$CLAUDE_PROJECT_DIR"

    # Print summary
    print_section "TEST SUMMARY"
    echo -e "${BOLD}Total Tests:${NC} $TOTAL_TESTS"
    echo -e "${GREEN}Passed:${NC} $PASSED_TESTS"
    echo -e "${RED}Failed:${NC} $FAILED_TESTS"

    if [ $FAILED_TESTS -eq 0 ]; then
        echo ""
        echo -e "${GREEN}${BOLD}✅ ALL TESTS PASSED!${NC}"
        exit 0
    else
        echo ""
        echo -e "${RED}${BOLD}❌ SOME TESTS FAILED!${NC}"
        exit 1
    fi
}

# Show usage if --help
if [ "${1:-}" = "--help" ] || [ "${1:-}" = "-h" ]; then
    echo "Branch Validation Hook Test Suite"
    echo ""
    echo "Usage: $0"
    echo ""
    echo "Tests the branch name validation hook against various"
    echo "valid and invalid branch name patterns to ensure it"
    echo "correctly enforces Software Factory naming conventions."
    echo ""
    exit 0
fi

# Run tests
main "$@"