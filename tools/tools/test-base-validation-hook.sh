#!/bin/bash
# Test script for the base branch validation pre-commit hook
# Tests various scenarios to ensure the hook correctly validates branch bases

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Temporary test directory
TEST_DIR="/tmp/sf2-hook-test-$$"
HOOK_PATH="$(pwd)/tools/hooks/pre-commit-base-validation"

print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_test() {
    echo -e "\n${YELLOW}TEST $((++TESTS_RUN)):${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓ PASSED:${NC} $1"
    ((TESTS_PASSED++))
}

print_failure() {
    echo -e "${RED}✗ FAILED:${NC} $1"
    ((TESTS_FAILED++))
}

# Setup test environment
setup_test_env() {
    print_header "Setting up test environment"

    # Create test directory
    rm -rf "$TEST_DIR"
    mkdir -p "$TEST_DIR"
    cd "$TEST_DIR"

    # Initialize git repo
    git init
    git config user.name "Test User"
    git config user.email "test@example.com"

    # Create initial commit on main
    echo "initial" > README.md
    git add README.md
    git commit -m "Initial commit"

    # Create basic orchestrator-state-v3.json
    cat > orchestrator-state-v3.json <<'EOF'
{
  "current_phase": 1,
  "current_wave": 1,
  "current_state": "PLANNING",
  "previous_state": "INIT",
  "transition_time": "2025-09-25T16:00:00Z",
  "branch_head_tracking": {
    "last_sync": "2025-09-25T16:00:00Z",
    "sync_count": 1,
    "branches": {},
    "validation_rules": {
      "effort_base_pattern": "For first effort in wave: base on main. For subsequent: base on previous effort",
      "split_base_pattern": "First split: base on oversized effort. Subsequent: base on previous split",
      "integration_base_pattern": "Wave integration: base on last effort"
    }
  },
  "phases_planned": 2,
  "waves_per_phase": [2, 2],
  "efforts_completed": [],
  "efforts_in_progress": [],
  "efforts_pending": [],
  "project_info": {
    "name": "test-project",
    "repository": "test-repo"
  }
}
EOF

    git add orchestrator-state-v3.json
    git commit -m "Add orchestrator-state-v3.json"

    # Install the hook
    mkdir -p .git/hooks
    cp "$HOOK_PATH" .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit

    echo "Test environment created in: $TEST_DIR"
}

# Test 1: First effort based on main (should pass)
test_first_effort_on_main() {
    print_test "First effort correctly based on main"

    # Create effort directory structure
    mkdir -p efforts/phase1/wave1/auth-service
    cd efforts/phase1/wave1/auth-service

    # Create a new branch from main
    git checkout -b phase1/wave1/auth-service

    # Make a change and try to commit
    echo "auth code" > auth.go
    git add auth.go

    if git commit -m "Add auth service" 2>/dev/null; then
        print_success "First effort on main accepted"
    else
        print_failure "First effort on main rejected (should have passed)"
    fi

    cd "$TEST_DIR"
}

# Test 2: Second effort based on wrong branch (should fail)
test_second_effort_wrong_base() {
    print_test "Second effort incorrectly based on main (should fail)"

    # Update orchestrator-state to show first effort completed
    jq '.efforts_completed += [{
        "name": "auth-service",
        "phase": 1,
        "wave": 1,
        "current_branch": "phase1/wave1/auth-service",
        "branch": "phase1/wave1/auth-service"
    }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

    git add orchestrator-state-v3.json
    git commit -m "Update state with completed effort"

    # Create second effort directory
    mkdir -p efforts/phase1/wave1/user-service
    cd efforts/phase1/wave1/user-service

    # Incorrectly create branch from main instead of previous effort
    git checkout main
    git checkout -b phase1/wave1/user-service

    # Try to commit
    echo "user code" > user.go
    git add user.go

    if git commit -m "Add user service" 2>/dev/null; then
        print_failure "Second effort on wrong base was accepted (should have failed)"
    else
        print_success "Second effort on wrong base was correctly rejected"
    fi

    cd "$TEST_DIR"
}

# Test 3: Split branch based on correct effort (should pass)
test_split_correct_base() {
    print_test "Split branch correctly based on oversized effort"

    # Switch to the auth-service branch
    git checkout phase1/wave1/auth-service

    # Create split directory
    mkdir -p efforts/phase1/wave1/auth-service/split-001
    cd efforts/phase1/wave1/auth-service/split-001

    # Create split branch from effort branch
    git checkout -b phase1/wave1/auth-service--split-001

    # Try to commit
    echo "split 1 code" > split1.go
    git add split1.go

    if git commit -m "First split of auth service" 2>/dev/null; then
        print_success "Split branch on correct base accepted"
    else
        print_failure "Split branch on correct base rejected (should have passed)"
    fi

    cd "$TEST_DIR"
}

# Test 4: Second split based on first split (should pass)
test_second_split_correct() {
    print_test "Second split correctly based on first split"

    # Switch to first split
    git checkout phase1/wave1/auth-service--split-001

    # Create second split directory
    mkdir -p efforts/phase1/wave1/auth-service/split-002
    cd efforts/phase1/wave1/auth-service/split-002

    # Create second split from first split
    git checkout -b phase1/wave1/auth-service--split-002

    # Try to commit
    echo "split 2 code" > split2.go
    git add split2.go

    if git commit -m "Second split of auth service" 2>/dev/null; then
        print_success "Second split on first split accepted"
    else
        print_failure "Second split on first split rejected (should have passed)"
    fi

    cd "$TEST_DIR"
}

# Test 5: Integration branch validation
test_integration_branch() {
    print_test "Wave integration branch validation"

    # Update state to show auth-service as last completed effort
    jq '.efforts_completed[-1].current_branch = "phase1/wave1/auth-service--split-002"' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
    git add orchestrator-state-v3.json
    git commit -m "Update with split completion"

    # Create integration directory
    mkdir -p efforts/phase1/wave1-integration
    cd efforts/phase1/wave1-integration

    # Create integration branch from last split
    git checkout phase1/wave1/auth-service--split-002
    git checkout -b phase1/wave1-integration

    # Try to commit
    echo "integration" > integration.txt
    git add integration.txt

    if git commit -m "Wave 1 integration" 2>/dev/null; then
        print_success "Integration branch on correct base accepted"
    else
        print_failure "Integration branch rejected (should have passed)"
    fi

    cd "$TEST_DIR"
}

# Test 6: Branch with missing files from base (should fail)
test_missing_critical_files() {
    print_test "Branch missing critical files from base"

    # Create a branch that's missing files
    git checkout main
    git checkout -b incomplete-branch

    # Create effort directory but without proper base
    mkdir -p efforts/phase1/wave2/payment-service
    cd efforts/phase1/wave2/payment-service

    # Add a file but miss files from expected base
    echo "payment" > payment.go
    git add payment.go

    # This should fail because it's not properly based on wave1-integration
    if git commit -m "Payment service" 2>/dev/null; then
        print_failure "Branch with missing files was accepted (should have failed)"
    else
        print_success "Branch with missing files was correctly rejected"
    fi

    cd "$TEST_DIR"
}

# Test 7: Bypass with --no-verify
test_bypass_with_no_verify() {
    print_test "Bypass validation with --no-verify flag"

    # Stay in the incomplete branch scenario
    cd efforts/phase1/wave2/payment-service 2>/dev/null || {
        mkdir -p efforts/phase1/wave2/payment-service
        cd efforts/phase1/wave2/payment-service
        echo "payment" > payment.go
        git add payment.go
    }

    # Try with --no-verify (should pass)
    if git commit --no-verify -m "Payment service (bypassed)" 2>/dev/null; then
        print_success "Bypass with --no-verify works as expected"
    else
        print_failure "Could not bypass with --no-verify"
    fi

    cd "$TEST_DIR"
}

# Test 8: HEAD tracking update
test_head_tracking_update() {
    print_test "HEAD tracking updates in orchestrator-state-v3.json"

    # Check if the successful commits updated the tracking
    local tracked_branches=$(jq -r '.branch_head_tracking.branches | keys[]' orchestrator-state-v3.json 2>/dev/null | wc -l)

    if [ "$tracked_branches" -gt 0 ]; then
        print_success "HEAD tracking was updated ($tracked_branches branches tracked)"

        # Show tracked branches
        echo "  Tracked branches:"
        jq -r '.branch_head_tracking.branches | to_entries[] | "    - \(.key): \(.value.head_commit)"' orchestrator-state-v3.json
    else
        print_failure "HEAD tracking was not updated"
    fi
}

# Cleanup
cleanup() {
    print_header "Cleaning up test environment"
    cd /
    rm -rf "$TEST_DIR"
    echo "Test directory removed"
}

# Main execution
main() {
    print_header "Base Branch Validation Hook Test Suite"
    echo "Testing hook at: $HOOK_PATH"

    # Check if hook exists
    if [ ! -f "$HOOK_PATH" ]; then
        echo -e "${RED}ERROR: Hook not found at $HOOK_PATH${NC}"
        exit 1
    fi

    # Setup
    setup_test_env

    # Run tests
    test_first_effort_on_main
    test_second_effort_wrong_base
    test_split_correct_base
    test_second_split_correct
    test_integration_branch
    test_missing_critical_files
    test_bypass_with_no_verify
    test_head_tracking_update

    # Results
    print_header "Test Results"
    echo -e "Tests Run: ${TESTS_RUN}"
    echo -e "Tests Passed: ${GREEN}${TESTS_PASSED}${NC}"
    echo -e "Tests Failed: ${RED}${TESTS_FAILED}${NC}"

    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "\n${GREEN}ALL TESTS PASSED! ✓${NC}"
        cleanup
        exit 0
    else
        echo -e "\n${RED}SOME TESTS FAILED ✗${NC}"
        echo -e "Test directory preserved at: $TEST_DIR"
        exit 1
    fi
}

# Handle interrupts
trap cleanup INT TERM

# Run if executed directly
if [ "${BASH_SOURCE[0]}" == "${0}" ]; then
    main "$@"
fi