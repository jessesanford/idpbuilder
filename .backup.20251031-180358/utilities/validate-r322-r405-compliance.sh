#!/bin/bash
# R322/R405 Compliance Validation Script
#
# PURPOSE:
#   Validate that state rules and agent outputs comply with R322/R405 interaction
#   Checks for correct continuation flag usage at R322 checkpoints
#
# USAGE:
#   bash utilities/validate-r322-r405-compliance.sh [state-rules-file]
#   bash utilities/validate-r322-r405-compliance.sh [agent-output-file]
#
# EXIT CODES:
#   0 - Fully compliant
#   1 - Violations found
#   2 - Invalid arguments

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0

log_check() {
    ((TOTAL_CHECKS++))
}

log_pass() {
    ((PASSED_CHECKS++))
    echo -e "${GREEN}✅ PASS${NC}: $1"
}

log_fail() {
    ((FAILED_CHECKS++))
    echo -e "${RED}❌ FAIL${NC}: $1"
}

log_warn() {
    echo -e "${YELLOW}⚠️  WARN${NC}: $1"
}

log_info() {
    echo -e "ℹ️  INFO: $1"
}

# Validate state rules file for R322/R405 compliance
validate_state_rules() {
    local rules_file="$1"

    if [ ! -f "$rules_file" ]; then
        log_fail "State rules file not found: $rules_file"
        return 1
    fi

    log_info "Validating state rules: $rules_file"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    # Check 1: R405 reference exists
    log_check
    if grep -q "R405" "$rules_file"; then
        log_pass "R405 reference found in state rules"
    else
        log_fail "R405 reference MISSING from state rules"
    fi

    # Check 2: R322 reference exists (if this is a checkpoint state)
    local state_name=$(basename "$(dirname "$rules_file")")
    if [[ "$state_name" =~ SPAWN_|WAITING_FOR_|MONITORING_ ]]; then
        log_check
        if grep -q "R322" "$rules_file"; then
            log_pass "R322 reference found (checkpoint state: $state_name)"
        else
            log_warn "R322 reference missing (checkpoint state: $state_name)"
        fi
    fi

    # Check 3: Continuation flag format specified
    log_check
    if grep -q "CONTINUE-SOFTWARE-FACTORY" "$rules_file"; then
        log_pass "Continuation flag format specified in rules"
    else
        log_fail "Continuation flag format NOT specified in rules"
    fi

    # Check 4: Enhanced format with checkpoint context mentioned (if R322 state)
    if [[ "$state_name" =~ SPAWN_ ]]; then
        log_check
        if grep -q "CHECKPOINT=R322" "$rules_file"; then
            log_pass "Enhanced format with CHECKPOINT=R322 specified"
        else
            log_warn "Enhanced format with CHECKPOINT=R322 NOT specified (recommended)"
        fi
    fi

    # Check 5: EXIT REQUIREMENTS section includes R405
    log_check
    if grep -A 20 "EXIT REQUIREMENTS" "$rules_file" | grep -q "R405"; then
        log_pass "EXIT REQUIREMENTS includes R405 continuation flag"
    else
        log_fail "EXIT REQUIREMENTS missing R405 continuation flag"
    fi

    # Check 6: No FALSE flag recommended for R322 checkpoints
    if [[ "$state_name" =~ SPAWN_ ]]; then
        log_check
        if grep -q "CONTINUE-SOFTWARE-FACTORY=FALSE" "$rules_file" | grep -v "ERROR\|EXCEPTIONAL"; then
            log_warn "State may be recommending FALSE flag at R322 checkpoint"
        else
            log_pass "No incorrect FALSE flag usage at R322 checkpoint"
        fi
    fi

    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

# Validate agent output for R322/R405 compliance
validate_agent_output() {
    local output_file="$1"

    if [ ! -f "$output_file" ]; then
        log_fail "Agent output file not found: $output_file"
        return 1
    fi

    log_info "Validating agent output: $output_file"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    # Check 1: Continuation flag exists
    log_check
    if grep -q "CONTINUE-SOFTWARE-FACTORY=" "$output_file"; then
        log_pass "Continuation flag found in output"
    else
        log_fail "Continuation flag MISSING from output (R405 VIOLATION)"
        return 1
    fi

    # Check 2: Flag is TRUE or FALSE (not other values)
    log_check
    local flag_value=$(grep "CONTINUE-SOFTWARE-FACTORY=" "$output_file" | tail -1 | grep -o "CONTINUE-SOFTWARE-FACTORY=[A-Z]*" | cut -d= -f2)
    if [[ "$flag_value" == "TRUE" ]] || [[ "$flag_value" == "FALSE" ]]; then
        log_pass "Flag value is valid: $flag_value"
    else
        log_fail "Flag value is INVALID: $flag_value (must be TRUE or FALSE)"
    fi

    # Check 3: Flag appears at end of output (R405 requirement)
    log_check
    local flag_line=$(grep -n "CONTINUE-SOFTWARE-FACTORY=" "$output_file" | tail -1 | cut -d: -f1)
    local total_lines=$(wc -l < "$output_file")
    local lines_after_flag=$((total_lines - flag_line))

    if [ $lines_after_flag -le 2 ]; then
        log_pass "Flag is at end of output (within last 2 lines)"
    else
        log_warn "Flag is NOT at end of output ($lines_after_flag lines after flag)"
    fi

    # Check 4: Checkpoint context (if enhanced format)
    log_check
    if grep "CONTINUE-SOFTWARE-FACTORY=TRUE" "$output_file" | grep -q "CHECKPOINT="; then
        local checkpoint_type=$(grep "CONTINUE-SOFTWARE-FACTORY=TRUE" "$output_file" | grep -o "CHECKPOINT=[A-Z0-9_]*" | cut -d= -f2)
        log_pass "Enhanced format with checkpoint context: $checkpoint_type"

        # Check if checkpoint type is valid
        if [[ "$checkpoint_type" =~ ^(R322|NONE)$ ]]; then
            log_pass "Checkpoint type is valid: $checkpoint_type"
        else
            log_warn "Checkpoint type may be non-standard: $checkpoint_type"
        fi
    else
        log_info "Basic format (no checkpoint context) - backward compatible"
    fi

    # Check 5: Error reason (if FALSE)
    log_check
    if grep "CONTINUE-SOFTWARE-FACTORY=FALSE" "$output_file" | grep -q "REASON="; then
        local error_reason=$(grep "CONTINUE-SOFTWARE-FACTORY=FALSE" "$output_file" | grep -o "REASON=[A-Z_]*" | cut -d= -f2)
        log_pass "Enhanced format with error reason: $error_reason"

        # Check if reason type is valid
        if [[ "$error_reason" =~ ^(ERROR|STATE_CORRUPTION|MISSING_FILES|RECURSIVE_SPLIT|ITERATION_OVERFLOW|WRONG_LOCATION)$ ]]; then
            log_pass "Error reason is valid: $error_reason"
        else
            log_warn "Error reason may be non-standard: $error_reason"
        fi
    fi

    # Check 6: R322 checkpoint with TRUE flag (not FALSE)
    log_check
    if grep -q "R322.*Stopping\|R322.*checkpoint" "$output_file"; then
        if grep "CONTINUE-SOFTWARE-FACTORY=TRUE" "$output_file" | grep -q "CHECKPOINT=R322"; then
            log_pass "R322 checkpoint uses TRUE CHECKPOINT=R322 (CORRECT!)"
        elif grep -q "CONTINUE-SOFTWARE-FACTORY=TRUE" "$output_file"; then
            log_warn "R322 checkpoint uses TRUE but missing CHECKPOINT=R322 context"
        elif grep -q "CONTINUE-SOFTWARE-FACTORY=FALSE" "$output_file"; then
            log_fail "R322 checkpoint uses FALSE (WRONG! Should use TRUE CHECKPOINT=R322)"
        fi
    fi

    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

# Main validation logic
main() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "R322/R405 COMPLIANCE VALIDATION"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    if [ $# -eq 0 ]; then
        echo "USAGE: $0 [file-to-validate]"
        echo ""
        echo "EXAMPLES:"
        echo "  $0 agent-states/software-factory/orchestrator/SPAWN_SW_ENGINEERS/rules.md"
        echo "  $0 test-workspace/orchestrator-output.json"
        echo ""
        exit 2
    fi

    local file="$1"

    # Determine file type and validate accordingly
    if [[ "$file" =~ rules\.md$ ]]; then
        validate_state_rules "$file"
    elif [[ "$file" =~ \.json$|\.log$|\.txt$ ]]; then
        validate_agent_output "$file"
    else
        log_fail "Unknown file type: $file"
        echo "Expected: rules.md (state rules) or .json/.log/.txt (agent output)"
        exit 2
    fi

    # Summary
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "VALIDATION SUMMARY"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Total Checks:  $TOTAL_CHECKS"
    echo -e "${GREEN}Passed:        $PASSED_CHECKS${NC}"
    echo -e "${RED}Failed:        $FAILED_CHECKS${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    if [ $FAILED_CHECKS -eq 0 ]; then
        echo -e "${GREEN}✅ ALL CHECKS PASSED${NC}"
        exit 0
    else
        echo -e "${RED}❌ VIOLATIONS DETECTED${NC}"
        exit 1
    fi
}

main "$@"
