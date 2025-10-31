#!/bin/bash
# State Machine Integrity Validator
# Validates software-factory-3.0-state-machine.json for consistency and correctness
#
# USAGE:
#   bash tools/validate-state-machine.sh                # Full validation
#   bash tools/validate-state-machine.sh --fast         # Quick validation (pre-commit)
#   bash tools/validate-state-machine.sh --verbose      # Detailed output
#
# EXIT CODES:
#   0 - All validations passed
#   1 - Validation failures found
#   2 - Usage error or missing dependencies
#
# RULES ENFORCED:
#   - R516: State naming conventions
#   - R206: Transition validation
#   - All allowed_transitions point to existing states

set -euo pipefail

# =============================================================================
# CONFIGURATION
# =============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(cd "$SCRIPT_DIR/.." && pwd)}"
STATE_MACHINE="${PROJECT_DIR}/state-machines/software-factory-3.0-state-machine.json"
VALIDATION_LIB="${PROJECT_DIR}/lib/state-validation-lib.sh"

# Mode flags
FAST_MODE=false
VERBOSE=false

# =============================================================================
# HELPER FUNCTIONS
# =============================================================================

usage() {
    cat <<EOF
Usage: $0 [OPTIONS]

Validate Software Factory 3.0 state machine integrity.

OPTIONS:
    --fast      Fast validation mode (pre-commit hook)
    --verbose   Verbose output with detailed information
    --help      Show this help message

VALIDATIONS:
    1. JSON syntax validation
    2. Required fields presence
    3. State naming conventions (R516)
    4. Transition target existence
    5. Orphaned states detection
    6. Circular transition detection (full mode only)

EXIT CODES:
    0  All validations passed
    1  Validation failures found
    2  Usage error or missing dependencies

EXAMPLES:
    # Full validation
    $0

    # Quick validation for pre-commit hook
    $0 --fast

    # Detailed diagnostics
    $0 --verbose

EOF
}

log_info() {
    echo "  $*"
}

log_verbose() {
    if [ "$VERBOSE" = true ]; then
        echo "    [VERBOSE] $*"
    fi
}

log_error() {
    echo "  ❌ ERROR: $*" >&2
}

log_warning() {
    echo "  ⚠️  WARNING: $*" >&2
}

log_success() {
    echo "  ✅ $*"
}

# =============================================================================
# VALIDATION FUNCTIONS
# =============================================================================

validate_prerequisites() {
    local errors=0

    # Check jq available
    if ! command -v jq >/dev/null 2>&1; then
        log_error "jq is required but not installed"
        ((errors++))
    fi

    # Check state machine exists
    if [ ! -f "$STATE_MACHINE" ]; then
        log_error "State machine not found: $STATE_MACHINE"
        ((errors++))
    fi

    # Check validation library exists
    if [ ! -f "$VALIDATION_LIB" ]; then
        log_error "Validation library not found: $VALIDATION_LIB"
        ((errors++))
    fi

    return $errors
}

validate_json_syntax() {
    log_info "Validating JSON syntax..."

    if ! jq . "$STATE_MACHINE" >/dev/null 2>&1; then
        log_error "Invalid JSON syntax in state machine"
        jq . "$STATE_MACHINE" 2>&1 | head -10 | sed 's/^/    /'
        return 1
    fi

    log_success "JSON syntax valid"
    return 0
}

validate_required_fields() {
    log_info "Validating required fields..."
    local errors=0

    # Check metadata section
    if ! jq -e '.metadata' "$STATE_MACHINE" >/dev/null 2>&1; then
        log_error "Missing 'metadata' section"
        ((errors++))
    fi

    # Check states section
    if ! jq -e '.states' "$STATE_MACHINE" >/dev/null 2>&1; then
        log_error "Missing 'states' section"
        ((errors++))
    fi

    # Check each state has required fields
    local states
    states=$(jq -r '.states | keys[]' "$STATE_MACHINE" 2>/dev/null || true)

    while IFS= read -r state; do
        [ -z "$state" ] && continue

        log_verbose "Checking required fields for state: $state"

        # Check description
        if ! jq -e ".states.\"$state\".description" "$STATE_MACHINE" >/dev/null 2>&1; then
            log_warning "State '$state' missing description"
        fi

        # Check agent
        if ! jq -e ".states.\"$state\".agent" "$STATE_MACHINE" >/dev/null 2>&1; then
            log_error "State '$state' missing agent field"
            ((errors++))
        fi

        # Check allowed_transitions
        if ! jq -e ".states.\"$state\".allowed_transitions" "$STATE_MACHINE" >/dev/null 2>&1; then
            log_warning "State '$state' missing allowed_transitions (may be terminal)"
        fi
    done <<< "$states"

    if [ $errors -eq 0 ]; then
        log_success "All required fields present"
    fi

    return $errors
}

validate_state_naming_conventions() {
    log_info "Validating state naming conventions (R516)..."
    local errors=0

    # Source validation library
    source "$VALIDATION_LIB"

    local states
    states=$(jq -r '.states | keys[]' "$STATE_MACHINE" 2>/dev/null || true)

    while IFS= read -r state; do
        [ -z "$state" ] && continue

        log_verbose "Checking naming convention for: $state"

        if ! is_valid_state_name "$state"; then
            log_error "State '$state' violates R516 naming conventions"
            log_error "  Must be ALL_CAPS with underscores, no leading/trailing underscores"
            ((errors++))
        fi
    done <<< "$states"

    if [ $errors -eq 0 ]; then
        log_success "All state names follow R516 conventions"
    fi

    return $errors
}

validate_transition_targets() {
    log_info "Validating transition targets exist..."

    # Source validation library
    source "$VALIDATION_LIB"

    # Use library function (captures errors internally)
    local output
    output=$(validate_transitions_point_to_existing_states "$STATE_MACHINE" 2>&1)
    local result=$?

    if [ $result -ne 0 ]; then
        # Parse output and display errors
        echo "$output" | grep "^ERROR:" | while IFS= read -r line; do
            log_error "$(echo "$line" | sed 's/^ERROR: //')"
        done
        return 1
    fi

    log_success "All transition targets exist"
    return 0
}

validate_orphaned_states() {
    log_info "Checking for orphaned states..."
    local warnings=0

    local states
    states=$(jq -r '.states | keys[]' "$STATE_MACHINE" 2>/dev/null || true)

    # Build set of all states that are transition targets
    local all_targets=()
    while IFS= read -r state; do
        [ -z "$state" ] && continue

        local transitions
        transitions=$(jq -r ".states.\"$state\".allowed_transitions[]?" "$STATE_MACHINE" 2>/dev/null || true)

        while IFS= read -r target; do
            [ -z "$target" ] && continue
            [ "$target" = "ANY_STATE" ] && continue
            all_targets+=("$target")
        done <<< "$transitions"
    done <<< "$states"

    # Check each state (except INIT) is reachable
    while IFS= read -r state; do
        [ -z "$state" ] && continue
        [ "$state" = "INIT" ] && continue  # INIT is entry point
        [ "$state" = "STARTUP_CONSULTATION" ] && continue  # Special consultation state
        [ "$state" = "SHUTDOWN_CONSULTATION" ] && continue  # Special consultation state

        # Check if state is in targets array
        local is_target=false
        for target in "${all_targets[@]}"; do
            if [ "$target" = "$state" ]; then
                is_target=true
                break
            fi
        done

        if [ "$is_target" = false ]; then
            log_warning "State '$state' may be orphaned (no incoming transitions)"
            ((warnings++))
        fi
    done <<< "$states"

    if [ $warnings -eq 0 ]; then
        log_success "No orphaned states found"
    else
        log_info "Found $warnings potentially orphaned state(s) (may be intentional)"
    fi

    return 0  # Warnings don't fail validation
}

detect_circular_transitions() {
    if [ "$FAST_MODE" = true ]; then
        log_verbose "Skipping circular transition detection (fast mode)"
        return 0
    fi

    log_info "Detecting potential circular transitions..."
    local warnings=0

    # Simple check: states that can transition to themselves
    local states
    states=$(jq -r '.states | keys[]' "$STATE_MACHINE" 2>/dev/null || true)

    while IFS= read -r state; do
        [ -z "$state" ] && continue

        local transitions
        transitions=$(jq -r ".states.\"$state\".allowed_transitions[]?" "$STATE_MACHINE" 2>/dev/null || true)

        while IFS= read -r target; do
            [ -z "$target" ] && continue

            if [ "$target" = "$state" ]; then
                log_warning "State '$state' can transition to itself (potential infinite loop)"
                ((warnings++))
            fi
        done <<< "$transitions"
    done <<< "$states"

    if [ $warnings -eq 0 ]; then
        log_success "No obvious circular transitions detected"
    fi

    return 0  # Warnings don't fail validation
}

# =============================================================================
# MAIN VALIDATION FLOW
# =============================================================================

main() {
    local errors=0

    echo "=================================="
    echo "State Machine Integrity Validation"
    echo "=================================="
    echo "File: $STATE_MACHINE"
    echo "Mode: $([ "$FAST_MODE" = true ] && echo "FAST" || echo "FULL")"
    echo ""

    # Parse arguments
    while [ $# -gt 0 ]; do
        case "$1" in
            --fast)
                FAST_MODE=true
                shift
                ;;
            --verbose)
                VERBOSE=true
                shift
                ;;
            --help)
                usage
                exit 0
                ;;
            *)
                echo "ERROR: Unknown option: $1" >&2
                usage
                exit 2
                ;;
        esac
    done

    # Run validations
    validate_prerequisites || ((errors+=$?))
    [ $errors -gt 0 ] && exit 2  # Stop if prerequisites fail

    validate_json_syntax || ((errors+=$?))
    validate_required_fields || ((errors+=$?))
    validate_state_naming_conventions || ((errors+=$?))
    validate_transition_targets || ((errors+=$?))
    validate_orphaned_states || ((errors+=$?))

    if [ "$FAST_MODE" = false ]; then
        detect_circular_transitions || ((errors+=$?))
    fi

    # Summary
    echo ""
    echo "=================================="
    echo "Validation Summary"
    echo "=================================="

    if [ $errors -eq 0 ]; then
        echo "✅✅✅ ALL VALIDATIONS PASSED ✅✅✅"
        echo ""
        echo "State machine is consistent and ready to use."
        return 0
    else
        echo "❌❌❌ VALIDATION FAILED ❌❌❌"
        echo ""
        echo "Found $errors error(s) that must be fixed."
        echo "See errors above for details."
        return 1
    fi
}

# Run main function
main "$@"
exit $?
