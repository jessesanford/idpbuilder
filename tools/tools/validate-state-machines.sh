#!/bin/bash

# Software Factory 2.0 - State Machine Validation Script
# Validates all state machine JSON files for correctness and consistency

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
STATE_MACHINES_DIR="$PROJECT_ROOT/state-machines"

# Counters
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0
WARNINGS=0

# Track validation results
declare -A VALIDATION_RESULTS

# Function to print section headers
print_section() {
    echo -e "\n${BLUE}${BOLD}$1${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

# Function to print subsection
print_subsection() {
    echo -e "\n${CYAN}$1${NC}"
    echo "────────────────────────────────────────────────────────────"
}

# Function to print success
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
    ((PASSED_CHECKS++))
    ((TOTAL_CHECKS++))
}

# Function to print error
print_error() {
    echo -e "${RED}❌ $1${NC}"
    ((FAILED_CHECKS++))
    ((TOTAL_CHECKS++))
}

# Function to print warning
print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
    ((WARNINGS++))
}

# Function to print info
print_info() {
    echo -e "${MAGENTA}ℹ️  $1${NC}"
}

# Function to validate JSON syntax
validate_json_syntax() {
    local file="$1"
    local name="$(basename "$file")"

    if jq -e . "$file" > /dev/null 2>&1; then
        print_success "$name: Valid JSON syntax"
        return 0
    else
        print_error "$name: Invalid JSON syntax"
        jq -e . "$file" 2>&1 | sed 's/^/    /'
        return 1
    fi
}

# Function to validate required metadata fields
validate_metadata() {
    local file="$1"
    local name="$(basename "$file")"

    local required_fields=("version" "name" "description" "created_at")
    local all_valid=true

    for field in "${required_fields[@]}"; do
        if jq -e ".metadata.$field" "$file" > /dev/null 2>&1; then
            local value=$(jq -r ".metadata.$field" "$file")
            if [[ -n "$value" && "$value" != "null" ]]; then
                print_success "$name: Has metadata.$field"
            else
                print_error "$name: metadata.$field is empty or null"
                all_valid=false
            fi
        else
            print_error "$name: Missing metadata.$field"
            all_valid=false
        fi
    done

    $all_valid && return 0 || return 1
}

# Function to validate state definitions
validate_states() {
    local file="$1"
    local name="$(basename "$file")"

    # Check if states array or object exists
    if jq -e '.states' "$file" > /dev/null 2>&1; then
        local state_count=$(jq '[.states | if type == "array" then .[] else keys[] end] | length' "$file")
        if [[ $state_count -gt 0 ]]; then
            print_success "$name: Has $state_count states defined"
            return 0
        else
            print_error "$name: No states defined"
            return 1
        fi
    else
        print_error "$name: Missing 'states' field"
        return 1
    fi
}

# Function to validate transition matrix
validate_transitions() {
    local file="$1"
    local name="$(basename "$file")"

    # Check for transition_matrix or transitions
    if jq -e '.transition_matrix' "$file" > /dev/null 2>&1; then
        local transition_count=$(jq '[.transition_matrix | to_entries | .[] | .value | if type == "array" then .[] else (to_entries | .[].value | if type == "array" then .[] else . end) end] | length' "$file" 2>/dev/null || echo "0")
        if [[ $transition_count -gt 0 ]]; then
            print_success "$name: Has transition matrix with transitions"
            return 0
        else
            print_warning "$name: Transition matrix exists but appears empty"
            return 0
        fi
    elif jq -e '.transitions' "$file" > /dev/null 2>&1; then
        print_success "$name: Has transitions field"
        return 0
    else
        print_warning "$name: No transition_matrix or transitions field"
        return 0
    fi
}

# Function to validate state consistency
validate_state_consistency() {
    local file="$1"
    local name="$(basename "$file")"

    # Get all states mentioned in states field
    local defined_states=$(jq -r '[
        .states |
        if type == "array" then
            .[]
        elif type == "object" then
            keys[]
        else
            empty
        end
    ] | sort | unique | .[]' "$file" 2>/dev/null || echo "")

    # Get all states mentioned in transitions
    # Handle both flat transition matrix and agent-organized matrix
    local transition_states=$(jq -r '[
        .transition_matrix |
        if . then
            to_entries |
            map(
                # Include the key itself if it looks like a state (uppercase with underscores)
                (if (.key | test("^[A-Z_]+$")) then .key else empty end),
                # Extract states from the value
                (.value |
                    if type == "object" then
                        # Could be state->states or agent->states
                        to_entries | map(
                            # Include nested keys if they look like states
                            (if (.key | test("^[A-Z_]+$")) then .key else empty end),
                            # Extract from nested values
                            (.value |
                                if type == "array" then
                                    .[]
                                elif type == "object" then
                                    to_entries | map(.value | if type == "array" then .[] else . end)
                                else
                                    .
                                end
                            )
                        )
                    elif type == "array" then
                        .[]
                    else
                        .
                    end
                )
            ) | flatten | .[]
        else
            empty
        end
    ] | select(. != null and . != "" and (. | type) == "string") | select(test("^[A-Z_]+$")) | sort | unique | .[]' "$file" 2>/dev/null || echo "")

    # Check for orphaned states in transitions
    local orphaned_found=false
    for state in $transition_states; do
        if [[ -n "$defined_states" ]] && ! echo "$defined_states" | grep -q "^${state}$"; then
            if ! $orphaned_found; then
                print_warning "$name: States in transitions not defined in states list:"
                orphaned_found=true
            fi
            echo "    - $state"
        fi
    done

    if ! $orphaned_found; then
        print_success "$name: All transition states are defined"
    fi

    return 0
}

# Function to validate specific state machine types
validate_state_machine_type() {
    local file="$1"
    local name="$(basename "$file")"

    case "$name" in
        "software-factory-3.0-state-machine.json")
            # Main state machine specific checks
            if jq -e '.agents.orchestrator' "$file" > /dev/null 2>&1; then
                print_success "$name: Has orchestrator agent definition"
            else
                print_error "$name: Missing orchestrator agent definition"
            fi
            ;;

        "initialization-state-machine.json")
            # Check for entry command
            if jq -e '.metadata.entry_command' "$file" > /dev/null 2>&1; then
                print_success "$name: Has entry_command defined"
            else
                print_warning "$name: No entry_command defined"
            fi
            ;;

        "pr-ready-state-machine.json")
            # Check for quality gates
            if jq -e '.quality_gates' "$file" > /dev/null 2>&1; then
                print_success "$name: Has quality_gates defined"
            else
                print_warning "$name: No quality_gates defined"
            fi
            ;;

        "fix-cascade-state-machine.json")
            # Check for registry requirements
            if jq -e '.registry_requirements' "$file" > /dev/null 2>&1; then
                print_success "$name: Has registry_requirements defined"
            else
                print_warning "$name: No registry_requirements defined"
            fi
            ;;

        "splitting-state-machine.json")
            # Check for split rules
            if jq -e '.split_rules' "$file" > /dev/null 2>&1; then
                print_success "$name: Has split_rules defined"
            else
                print_warning "$name: No split_rules defined"
            fi
            ;;

        "integration-state-machine.json")
            # Check for integration types
            if jq -e '.integration_types' "$file" > /dev/null 2>&1; then
                print_success "$name: Has integration_types defined"
            else
                print_warning "$name: No integration_types defined"
            fi
            ;;
    esac
}

# Function to check for sub-state machine markers
validate_sub_state_machine() {
    local file="$1"
    local name="$(basename "$file")"

    # Skip main state machine
    if [[ "$name" == "software-factory-3.0-state-machine.json" ]]; then
        return 0
    fi

    # Check if it's marked as a sub-state machine
    local type=$(jq -r '.metadata.type // ""' "$file" 2>/dev/null)
    if [[ "$type" == "SUB_STATE_MACHINE" ]]; then
        print_success "$name: Correctly marked as SUB_STATE_MACHINE"

        # Check for entry and exit points
        if jq -e '.entry_points' "$file" > /dev/null 2>&1; then
            print_success "$name: Has entry_points defined"
        else
            print_warning "$name: No entry_points defined for sub-state machine"
        fi

        if jq -e '.exit_points' "$file" > /dev/null 2>&1; then
            print_success "$name: Has exit_points defined"
        else
            print_warning "$name: No exit_points defined for sub-state machine"
        fi
    else
        print_info "$name: Not marked as SUB_STATE_MACHINE (may be intentional)"
    fi
}

# Main validation function
validate_state_machine() {
    local file="$1"
    local name="$(basename "$file")"

    print_subsection "Validating: $name"

    local validation_passed=true

    # Run all validations
    validate_json_syntax "$file" || validation_passed=false
    validate_metadata "$file" || validation_passed=false
    validate_states "$file" || validation_passed=false
    validate_transitions "$file"
    validate_state_consistency "$file"
    validate_state_machine_type "$file"
    validate_sub_state_machine "$file"

    # Store result
    if $validation_passed; then
        VALIDATION_RESULTS["$name"]="PASSED"
    else
        VALIDATION_RESULTS["$name"]="FAILED"
    fi
}

# Main execution
main() {
    print_section "SOFTWARE FACTORY 2.0 - STATE MACHINE VALIDATION"
    echo "Validating all state machines in: $STATE_MACHINES_DIR"

    # Check if state-machines directory exists
    if [[ ! -d "$STATE_MACHINES_DIR" ]]; then
        print_error "State machines directory not found: $STATE_MACHINES_DIR"
        exit 1
    fi

    # Find all JSON files
    local json_files=$(find "$STATE_MACHINES_DIR" -name "*.json" -type f | sort)
    if [[ -z "$json_files" ]]; then
        print_error "No JSON files found in $STATE_MACHINES_DIR"
        exit 1
    fi

    # Count files
    local file_count=$(echo "$json_files" | wc -l)
    print_info "Found $file_count state machine files to validate"

    # Validate each file
    for file in $json_files; do
        validate_state_machine "$file"
    done

    # Print summary
    print_section "VALIDATION SUMMARY"

    echo -e "\n${BOLD}Results by File:${NC}"
    for name in "${!VALIDATION_RESULTS[@]}"; do
        result="${VALIDATION_RESULTS[$name]}"
        if [[ "$result" == "PASSED" ]]; then
            echo -e "  ${GREEN}✅ $name${NC}"
        else
            echo -e "  ${RED}❌ $name${NC}"
        fi
    done | sort

    echo -e "\n${BOLD}Statistics:${NC}"
    echo "  Total Checks:   $TOTAL_CHECKS"
    echo -e "  ${GREEN}Passed:${NC}        $PASSED_CHECKS"
    echo -e "  ${RED}Failed:${NC}        $FAILED_CHECKS"
    echo -e "  ${YELLOW}Warnings:${NC}      $WARNINGS"

    # Calculate pass rate
    if [[ $TOTAL_CHECKS -gt 0 ]]; then
        local pass_rate=$((PASSED_CHECKS * 100 / TOTAL_CHECKS))
        echo -e "  ${BOLD}Pass Rate:${NC}     ${pass_rate}%"
    fi

    # Final result
    echo
    if [[ $FAILED_CHECKS -eq 0 ]]; then
        echo -e "${GREEN}${BOLD}✅ ALL STATE MACHINES VALID${NC}"
        exit 0
    else
        echo -e "${RED}${BOLD}❌ VALIDATION FAILED${NC}"
        echo "Please fix the errors above and run validation again."
        exit 1
    fi
}

# Run main function
main "$@"