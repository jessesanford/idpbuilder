#!/bin/bash
# Comprehensive SF 3.0 Transition Validator
# Validates specific transitions and scans for invalid state references
#
# USAGE:
#   bash tools/validate-sf-transitions.sh CASCADE_REINTEGRATION INTEGRATE_WAVE_EFFORTS
#   bash tools/validate-sf-transitions.sh CASCADE_REINTEGRATION --all
#   bash tools/validate-sf-transitions.sh --scan-rules agent-states/**/*.md
#   bash tools/validate-sf-transitions.sh --comprehensive
#
# EXIT CODES:
#   0 - Validation passed
#   1 - Validation failed
#   2 - Usage error
#
# RULES ENFORCED:
#   - R206: State transition validation
#   - R516: State naming conventions

set -euo pipefail

# =============================================================================
# CONFIGURATION
# =============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(cd "$SCRIPT_DIR/.." && pwd)}"
VALIDATION_LIB="${PROJECT_DIR}/lib/state-validation-lib.sh"

# =============================================================================
# HELPER FUNCTIONS
# =============================================================================

usage() {
    cat <<EOF
Usage: $0 <from_state> <to_state>       # Validate specific transition
       $0 <from_state> --all             # Validate all transitions from state
       $0 --scan-rules <pattern>         # Scan files for invalid state refs
       $0 --comprehensive                # Run all validations

Comprehensive SF 3.0 Transition Validator

MODES:
    Single Transition:
        Validate a specific state transition
        Example: $0 CASCADE_REINTEGRATION INTEGRATE_WAVE_EFFORTS

    All Transitions:
        Validate all allowed transitions from a state
        Example: $0 CASCADE_REINTEGRATION --all

    Scan Rules:
        Scan rule files for invalid state references
        Example: $0 --scan-rules 'agent-states/**/*.md'

    Comprehensive:
        Run all validations (transition checks + rule scans)
        Example: $0 --comprehensive

EXIT CODES:
    0  Validation passed
    1  Validation failed
    2  Usage error or missing dependencies

EOF
}

log_info() {
    echo "  $*"
}

log_error() {
    echo "  ❌ ERROR: $*" >&2
}

log_success() {
    echo "  ✅ $*"
}

log_warning() {
    echo "  ⚠️  WARNING: $*" >&2
}

# =============================================================================
# VALIDATION FUNCTIONS
# =============================================================================

validate_single_transition() {
    local from_state="$1"
    local to_state="$2"

    echo "=================================="
    echo "Transition Validation"
    echo "=================================="
    echo "From: $from_state"
    echo "To:   $to_state"
    echo ""

    # Source library
    source "$VALIDATION_LIB"

    # Validate from_state exists
    log_info "Checking source state exists..."
    if ! validate_state_exists_in_machine "$from_state" 2>&1; then
        log_error "Source state '$from_state' does not exist in state machine"
        return 1
    fi
    log_success "Source state exists"

    # Validate to_state exists
    log_info "Checking target state exists..."
    if ! validate_state_exists_in_machine "$to_state" 2>&1; then
        log_error "Target state '$to_state' does not exist in state machine"
        return 1
    fi
    log_success "Target state exists"

    # Validate transition is allowed
    log_info "Checking transition is allowed..."
    if ! validate_transition_exists "$from_state" "$to_state" 2>&1; then
        log_error "Transition $from_state → $to_state is not allowed"

        # Show allowed transitions
        local allowed
        allowed=$(get_allowed_transitions "$from_state")
        if [ -n "$allowed" ]; then
            echo ""
            echo "  Allowed transitions from '$from_state':"
            for target in $allowed; do
                echo "    - $target"
            done
        fi

        return 1
    fi
    log_success "Transition is allowed"

    echo ""
    echo "✅✅✅ TRANSITION VALID ✅✅✅"
    return 0
}

validate_all_transitions_from_state() {
    local from_state="$1"

    echo "=================================="
    echo "All Transitions from State"
    echo "=================================="
    echo "State: $from_state"
    echo ""

    # Source library
    source "$VALIDATION_LIB"

    # Validate from_state exists
    if ! validate_state_exists_in_machine "$from_state" >/dev/null 2>&1; then
        log_error "State '$from_state' does not exist in state machine"
        return 1
    fi

    # Get all allowed transitions
    local allowed
    allowed=$(get_allowed_transitions "$from_state")

    if [ -z "$allowed" ]; then
        log_warning "State '$from_state' has no allowed transitions (terminal state?)"
        return 0
    fi

    echo "  Found $(echo "$allowed" | wc -w) allowed transition(s)"
    echo ""

    local errors=0
    local checked=0

    for target in $allowed; do
        ((checked++))

        # Skip ANY_STATE
        if [ "$target" = "ANY_STATE" ]; then
            log_info "$checked. $from_state → $target (special case, always valid)"
            continue
        fi

        # Validate target exists
        if validate_state_exists_in_machine "$target" >/dev/null 2>&1; then
            log_success "$checked. $from_state → $target (valid)"
        else
            log_error "$checked. $from_state → $target (target does not exist)"
            ((errors++))
        fi
    done

    echo ""
    echo "=================================="
    echo "Summary"
    echo "=================================="
    echo "Checked: $checked transitions"
    echo "Valid:   $((checked - errors))"
    echo "Invalid: $errors"
    echo ""

    if [ $errors -eq 0 ]; then
        echo "✅✅✅ ALL TRANSITIONS VALID ✅✅✅"
        return 0
    else
        echo "❌❌❌ FOUND $errors INVALID TRANSITION(S) ❌❌❌"
        return 1
    fi
}

scan_rules_for_invalid_references() {
    local pattern="$1"

    echo "=================================="
    echo "Scan Rules for Invalid State References"
    echo "=================================="
    echo "Pattern: $pattern"
    echo ""

    # Source library
    source "$VALIDATION_LIB"

    # Find matching files
    local files
    files=$(find "$PROJECT_DIR" -path "$PROJECT_DIR/$pattern" 2>/dev/null || true)

    if [ -z "$files" ]; then
        log_warning "No files found matching pattern: $pattern"
        return 1
    fi

    local total_files=0
    local files_with_issues=0
    local total_issues=0

    while IFS= read -r file; do
        [ ! -f "$file" ] && continue
        ((total_files++))

        log_info "Scanning: ${file#$PROJECT_DIR/}"

        # Capture validation output
        local output
        output=$(validate_state_references_in_file "$file" 2>&1)
        local result=$?

        if [ $result -ne 0 ]; then
            ((files_with_issues++))

            # Count warnings in output
            local file_issues
            file_issues=$(echo "$output" | grep -c "^WARNING:" || true)
            total_issues=$((total_issues + file_issues))

            # Show warnings
            echo "$output" | grep "^WARNING:" | sed 's/^/    /'
        fi
    done <<< "$files"

    echo ""
    echo "=================================="
    echo "Scan Summary"
    echo "=================================="
    echo "Files scanned:     $total_files"
    echo "Files with issues: $files_with_issues"
    echo "Total issues:      $total_issues"
    echo ""

    if [ $total_issues -eq 0 ]; then
        echo "✅✅✅ NO INVALID REFERENCES FOUND ✅✅✅"
        return 0
    else
        echo "⚠️⚠️⚠️ FOUND $total_issues POTENTIAL ISSUE(S) ⚠️⚠️⚠️"
        echo "(These may be false positives - manual review recommended)"
        return 1
    fi
}

run_comprehensive_validation() {
    echo "=================================="
    echo "Comprehensive SF 3.0 Validation"
    echo "=================================="
    echo ""

    local errors=0

    # 1. Validate state machine integrity
    echo "1. Validating state machine integrity..."
    if bash "${SCRIPT_DIR}/validate-state-machine.sh" --fast >/dev/null 2>&1; then
        log_success "State machine integrity check passed"
    else
        log_error "State machine integrity check failed"
        echo "   Run: bash tools/validate-state-machine.sh for details"
        ((errors++))
    fi
    echo ""

    # 2. Scan orchestrator rules
    echo "2. Scanning orchestrator state rules..."
    if scan_rules_for_invalid_references "agent-states/software-factory/orchestrator/*/rules.md" >/dev/null 2>&1; then
        log_success "Orchestrator rules scan passed"
    else
        log_warning "Orchestrator rules scan found potential issues (may be false positives)"
    fi
    echo ""

    # 3. Scan state manager rules
    echo "3. Scanning state manager rules..."
    if scan_rules_for_invalid_references "agent-states/state-manager/*/rules.md" >/dev/null 2>&1; then
        log_success "State manager rules scan passed"
    else
        log_warning "State manager rules scan found potential issues (may be false positives)"
    fi
    echo ""

    # Summary
    echo "=================================="
    echo "Comprehensive Validation Summary"
    echo "=================================="
    echo ""

    if [ $errors -eq 0 ]; then
        echo "✅✅✅ COMPREHENSIVE VALIDATION PASSED ✅✅✅"
        echo ""
        echo "All critical validations passed."
        echo "Review any warnings above manually."
        return 0
    else
        echo "❌❌❌ COMPREHENSIVE VALIDATION FAILED ❌❌❌"
        echo ""
        echo "Found $errors critical error(s)."
        echo "Run specific validation tools for details."
        return 1
    fi
}

# =============================================================================
# MAIN
# =============================================================================

main() {
    # Check prerequisites
    if [ ! -f "$VALIDATION_LIB" ]; then
        echo "ERROR: Validation library not found: $VALIDATION_LIB" >&2
        exit 2
    fi

    # Parse arguments
    if [ $# -eq 0 ]; then
        usage
        exit 2
    fi

    case "$1" in
        --help|-h)
            usage
            exit 0
            ;;
        --comprehensive)
            run_comprehensive_validation
            exit $?
            ;;
        --scan-rules)
            if [ $# -lt 2 ]; then
                echo "ERROR: --scan-rules requires a pattern argument" >&2
                usage
                exit 2
            fi
            scan_rules_for_invalid_references "$2"
            exit $?
            ;;
        *)
            # Single or all transitions mode
            if [ $# -lt 2 ]; then
                echo "ERROR: Missing arguments" >&2
                usage
                exit 2
            fi

            local from_state="$1"
            local to_state="$2"

            if [ "$to_state" = "--all" ]; then
                validate_all_transitions_from_state "$from_state"
            else
                validate_single_transition "$from_state" "$to_state"
            fi
            exit $?
            ;;
    esac
}

main "$@"
