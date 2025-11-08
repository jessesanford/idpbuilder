#!/usr/bin/env bash
# test-agent-cleanup.sh
#
# Test suite for R610/R611/R612 agent cleanup functionality
# Tests cleanup-completed-agents.sh utility and validates proper operation
#
# Usage: bash tools/test-agent-cleanup.sh

set -euo pipefail

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Test counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Logging
log_test() {
    echo -e "${YELLOW}TEST[$TESTS_RUN]:${NC} $*"
}

log_pass() {
    ((TESTS_PASSED++))
    echo -e "${GREEN}✅ PASS:${NC} $*"
}

log_fail() {
    ((TESTS_FAILED++))
    echo -e "${RED}❌ FAIL:${NC} $*"
}

# Setup test environment
setup_test_env() {
    TEST_DIR=$(mktemp -d)
    TEST_STATE_FILE="$TEST_DIR/orchestrator-state-v3.json"

    # Create minimal test state file
    cat > "$TEST_STATE_FILE" <<'EOF'
{
  "state_machine": {
    "current_state": "MONITORING_SWE_PROGRESS"
  },
  "active_agents": [],
  "agents_history": []
}
EOF

    echo "$TEST_DIR"
}

# Cleanup test environment
cleanup_test_env() {
    local test_dir="$1"
    rm -rf "$test_dir"
}

# Test 1: Cleanup utility exists and is executable
test_cleanup_utility_exists() {
    ((TESTS_RUN++))
    log_test "Cleanup utility exists and is executable"

    if [ -f "tools/cleanup-completed-agents.sh" ] && [ -x "tools/cleanup-completed-agents.sh" ]; then
        log_pass "Cleanup utility found and executable"
        return 0
    else
        log_fail "Cleanup utility not found or not executable"
        return 1
    fi
}

# Test 2: Validate mode works with clean state
test_validate_clean_state() {
    ((TESTS_RUN++))
    log_test "Validate mode with clean state (no completed agents)"

    TEST_DIR=$(setup_test_env)
    export CLAUDE_PROJECT_DIR="$TEST_DIR"

    if bash tools/cleanup-completed-agents.sh --validate > /dev/null 2>&1; then
        log_pass "Validation passed for clean state"
        cleanup_test_env "$TEST_DIR"
        return 0
    else
        log_fail "Validation failed for clean state"
        cleanup_test_env "$TEST_DIR"
        return 1
    fi
}

# Test 3: Detect completed agents
test_detect_completed_agents() {
    ((TESTS_RUN++))
    log_test "Detect completed agents in active_agents"

    TEST_DIR=$(setup_test_env)
    TEST_STATE_FILE="$TEST_DIR/orchestrator-state-v3.json"

    # Add completed agent
    cat > "$TEST_STATE_FILE" <<'EOF'
{
  "state_machine": {"current_state": "MONITORING_SWE_PROGRESS"},
  "active_agents": [
    {
      "agent_id": "swe-test-1",
      "agent_type": "sw-engineer",
      "state": "COMPLETE",
      "effort_id": "1.1.1"
    }
  ],
  "agents_history": []
}
EOF

    export CLAUDE_PROJECT_DIR="$TEST_DIR"

    if ! bash tools/cleanup-completed-agents.sh --validate > /dev/null 2>&1; then
        log_pass "Correctly detected completed agent"
        cleanup_test_env "$TEST_DIR"
        return 0
    else
        log_fail "Failed to detect completed agent"
        cleanup_test_env "$TEST_DIR"
        return 1
    fi
}

# Test 4: Schema files exist
test_schema_files_exist() {
    ((TESTS_RUN++))
    log_test "R612 schema file exists"

    if [ -f "schemas/agents-history-schema.json" ]; then
        log_pass "agents-history-schema.json exists"
        return 0
    else
        log_fail "agents-history-schema.json not found"
        return 1
    fi
}

# Test 5: Rule files exist
test_rule_files_exist() {
    ((TESTS_RUN++))
    log_test "R610/R611/R612/R613 rule files exist"

    local all_exist=true

    for rule in R610 R611 R612 R613; do
        if [ ! -f "rule-library/${rule}"*.md ]; then
            log_fail "Rule file ${rule}*.md not found"
            all_exist=false
        fi
    done

    if $all_exist; then
        log_pass "All rule files exist"
        return 0
    else
        return 1
    fi
}

# Test 6: RULE-REGISTRY updated
test_rule_registry_updated() {
    ((TESTS_RUN++))
    log_test "RULE-REGISTRY.md contains R610-R613 entries"

    local all_present=true

    for rule in R610 R611 R612 R613; do
        if ! grep -q "## $rule " rule-library/RULE-REGISTRY.md; then
            log_fail "Rule $rule not in RULE-REGISTRY.md"
            all_present=false
        fi
    done

    if $all_present; then
        log_pass "All rules in RULE-REGISTRY.md"
        return 0
    else
        return 1
    fi
}

# Test 7: Monitoring states updated
test_monitoring_states_updated() {
    ((TESTS_RUN++))
    log_test "Monitoring states have R610 cleanup sections"

    local all_updated=true

    if ! grep -q "R610" agent-states/software-factory/orchestrator/MONITORING_SWE_PROGRESS/rules.md; then
        log_fail "MONITORING_SWE_PROGRESS missing R610"
        all_updated=false
    fi

    if ! grep -q "R610" agent-states/software-factory/orchestrator/MONITORING_EFFORT_REVIEWS/rules.md; then
        log_fail "MONITORING_EFFORT_REVIEWS missing R610"
        all_updated=false
    fi

    if $all_updated; then
        log_pass "Monitoring states updated with R610"
        return 0
    else
        return 1
    fi
}

# Test 8: Boundary states updated
test_boundary_states_updated() {
    ((TESTS_RUN++))
    log_test "Boundary states have R610 validation sections"

    local all_updated=true

    if ! grep -q "R610" agent-states/software-factory/orchestrator/COMPLETE_WAVE/rules.md; then
        log_fail "COMPLETE_WAVE missing R610"
        all_updated=false
    fi

    if ! grep -q "R610" agent-states/software-factory/orchestrator/COMPLETE_PHASE/rules.md; then
        log_fail "COMPLETE_PHASE missing R610"
        all_updated=false
    fi

    if $all_updated; then
        log_pass "Boundary states updated with R610"
        return 0
    else
        return 1
    fi
}

# Main test execution
main() {
    echo "================================================"
    echo "R610/R611/R612/R613 Agent Cleanup Test Suite"
    echo "================================================"
    echo ""

    # Run all tests
    test_cleanup_utility_exists || true
    test_validate_clean_state || true
    test_detect_completed_agents || true
    test_schema_files_exist || true
    test_rule_files_exist || true
    test_rule_registry_updated || true
    test_monitoring_states_updated || true
    test_boundary_states_updated || true

    # Report results
    echo ""
    echo "================================================"
    echo "Test Results"
    echo "================================================"
    echo "Tests Run: $TESTS_RUN"
    echo -e "${GREEN}Tests Passed: $TESTS_PASSED${NC}"
    echo -e "${RED}Tests Failed: $TESTS_FAILED${NC}"
    echo ""

    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "${GREEN}✅ ALL TESTS PASSED${NC}"
        return 0
    else
        echo -e "${RED}❌ SOME TESTS FAILED${NC}"
        return 1
    fi
}

# Run tests
main "$@"
