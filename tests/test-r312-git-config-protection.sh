#!/bin/bash

# Test script for R312 Git Config Immutability Protocol

set -euo pipefail

echo "═══════════════════════════════════════════════════════════════════════"
echo "🧪 R312 Git Config Immutability Protocol - Test Suite"
echo "═══════════════════════════════════════════════════════════════════════"
echo ""

# Store test results
TESTS_PASSED=0
TESTS_FAILED=0
FAILURES=""

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test functions
run_test() {
    local TEST_NAME="$1"
    local TEST_COMMAND="$2"
    local EXPECTED_RESULT="$3"  # "pass" or "fail"
    
    echo -n "Testing: $TEST_NAME ... "
    
    if [ "$EXPECTED_RESULT" = "pass" ]; then
        if eval "$TEST_COMMAND" 2>/dev/null; then
            echo -e "${GREEN}✅ PASSED${NC}"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            echo -e "${RED}❌ FAILED${NC}"
            TESTS_FAILED=$((TESTS_FAILED + 1))
            FAILURES="$FAILURES\n  - $TEST_NAME: Expected to pass but failed"
        fi
    else
        if eval "$TEST_COMMAND" 2>/dev/null; then
            echo -e "${RED}❌ FAILED${NC}"
            TESTS_FAILED=$((TESTS_FAILED + 1))
            FAILURES="$FAILURES\n  - $TEST_NAME: Expected to fail but passed"
        else
            echo -e "${GREEN}✅ PASSED${NC} (correctly failed)"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        fi
    fi
}

# Create test repository
TEST_DIR="/tmp/r312-test-$$"
echo "📂 Creating test repository in: $TEST_DIR"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"
git init --quiet
git config user.email "test@example.com"
git config user.name "Test User"
echo "test" > test.txt
git add test.txt
git commit -m "Initial commit" --quiet

echo ""
echo "🔧 Setting up test environment..."
echo "──────────────────────────────────"

# Test 1: Config should be writable initially
echo ""
echo "📋 Test Group 1: Initial State"
echo "──────────────────────────────────"
run_test "Config is initially writable" "[ -w .git/config ]" "pass"
run_test "Can modify config when writable" "git config --local test.value 'test' && git config --unset test.value" "pass"

# Test 2: Lock the config per R312
echo ""
echo "📋 Test Group 2: Locking Config (R312)"
echo "──────────────────────────────────"
echo "🔒 Applying R312 lock..."
chmod 444 .git/config

# Create lock marker
cat > .git/R312_CONFIG_LOCKED << EOF
Timestamp: $(date '+%Y-%m-%d %H:%M:%S')
Locked by: test-script
Purpose: Testing R312 protection
EOF

run_test "Config is readonly after lock" "[ ! -w .git/config ]" "pass"
run_test "Lock marker exists" "[ -f .git/R312_CONFIG_LOCKED ]" "pass"

# Test 3: Protected operations should fail
echo ""
echo "📋 Test Group 3: Protected Operations (Should Fail)"
echo "──────────────────────────────────"
run_test "Cannot modify config values" "git config --local test.value 'test'" "fail"
run_test "Cannot change remote URL" "git remote set-url origin https://new.url" "fail"
run_test "Cannot add new remote" "git remote add upstream https://upstream.url" "fail"
run_test "Cannot change branch upstream" "git branch --set-upstream-to=origin/main" "fail"
run_test "Cannot checkout new branch" "git checkout -b test-branch" "fail"

# Test 4: Allowed operations should still work
echo ""
echo "📋 Test Group 4: Allowed Operations (Should Pass)"
echo "──────────────────────────────────"
echo "test2" > test2.txt
run_test "Can stage files" "git add test2.txt" "pass"
run_test "Can commit changes" "git commit -m 'Test commit' --quiet" "pass"
run_test "Can view status" "git status --short" "pass"
run_test "Can view log" "git log --oneline -1" "pass"
run_test "Can view diff" "git diff HEAD~1" "pass"

# Test 5: Integration exception
echo ""
echo "📋 Test Group 5: Integration Exception (R312)"
echo "──────────────────────────────────"
echo "🔓 Applying integration exception..."
chmod 644 .git/config

# Create exception marker
cat > .git/R312_INTEGRATION_EXCEPTION << EOF
Timestamp: $(date '+%Y-%m-%d %H:%M:%S')
Unlocked by: test-script
Purpose: Testing integration exception
EOF

run_test "Config is writable with exception" "[ -w .git/config ]" "pass"
run_test "Exception marker exists" "[ -f .git/R312_INTEGRATION_EXCEPTION ]" "pass"
run_test "Can modify config with exception" "git config --local test.integration 'allowed' && git config --unset test.integration" "pass"

# Test 6: Verify protection detection
echo ""
echo "📋 Test Group 6: SW Engineer Validation"
echo "──────────────────────────────────"

# Relock for SW engineer test
chmod 444 .git/config

# Simulate SW engineer validation
validate_r312() {
    if [ ! -f .git/config ]; then
        return 1  # No config found
    fi
    
    if [ -w .git/config ]; then
        return 1  # Config is writable (violation)
    fi
    
    # Test that config modifications fail
    if git config --local test.value "test" 2>/dev/null; then
        return 1  # Config modification succeeded (shouldn't happen)
    fi
    
    return 0  # All validations passed
}

run_test "SW Engineer R312 validation passes" "validate_r312" "pass"

# Test 7: Verify lock survives git operations
echo ""
echo "📋 Test Group 7: Lock Persistence"
echo "──────────────────────────────────"
echo "test3" > test3.txt
git add test3.txt
git commit -m "Test lock persistence" --quiet
run_test "Lock persists after commit" "[ ! -w .git/config ]" "pass"

# Cleanup
echo ""
echo "🧹 Cleaning up test repository..."
cd /
rm -rf "$TEST_DIR"

# Results
echo ""
echo "═══════════════════════════════════════════════════════════════════════"
echo "📊 TEST RESULTS"
echo "──────────────────────────────────"
echo -e "  Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "  Failed: ${RED}$TESTS_FAILED${NC}"

if [ $TESTS_FAILED -gt 0 ]; then
    echo ""
    echo -e "${RED}Failed Tests:${NC}"
    echo -e "$FAILURES"
    echo ""
    echo -e "${RED}❌ TEST SUITE FAILED${NC}"
    exit 1
else
    echo ""
    echo -e "${GREEN}✅ ALL TESTS PASSED!${NC}"
    echo ""
    echo "R312 Git Config Immutability Protocol is working correctly!"
fi

echo "═══════════════════════════════════════════════════════════════════════"