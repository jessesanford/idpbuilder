#!/usr/bin/env bash
#
# validate-state-schema.sh
#
# Validates orchestrator state files against JSON schema
# Part of SF 3.0 state integrity system
#
# Usage:
#   bash tools/validate-state-schema.sh [state-file.json]
#   bash tools/validate-state-schema.sh (validates all state files)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

SCHEMA_FILE="$PROJECT_ROOT/schemas/orchestrator-state-v3.schema.json"
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Validation results
VALIDATION_PASSED=0
VALIDATION_FAILED=0
VALIDATION_ERRORS=()

echo "=================================================="
echo "  SOFTWARE FACTORY 3.0 - STATE SCHEMA VALIDATOR"
echo "=================================================="
echo ""
echo "Schema: $SCHEMA_FILE"
echo ""

# Check if ajv-cli is available
if ! command -v ajv &> /dev/null; then
    echo -e "${YELLOW}WARNING: ajv-cli not found. Installing...${NC}"
    npm install -g ajv-cli ajv-formats 2>&1 | grep -v "npm WARN" || true

    if ! command -v ajv &> /dev/null; then
        echo -e "${RED}ERROR: Could not install ajv-cli${NC}"
        echo "Please install manually: npm install -g ajv-cli ajv-formats"
        exit 1
    fi
fi

# Function to validate a single state file
validate_state_file() {
    local state_file="$1"
    local file_name=$(basename "$state_file")

    echo -n "Validating: $file_name ... "

    # Check if file exists
    if [[ ! -f "$state_file" ]]; then
        echo -e "${RED}SKIP (file not found)${NC}"
        return
    fi

    # Check if file is valid JSON
    if ! jq empty "$state_file" 2>/dev/null; then
        echo -e "${RED}FAIL (invalid JSON)${NC}"
        VALIDATION_FAILED=$((VALIDATION_FAILED + 1))
        VALIDATION_ERRORS+=("$file_name: Invalid JSON syntax")
        return
    fi

    # Validate against schema using ajv
    if ajv validate -s "$SCHEMA_FILE" -d "$state_file" --strict=false 2>&1 | grep -q "valid"; then
        echo -e "${GREEN}PASS${NC}"
        VALIDATION_PASSED=$((VALIDATION_PASSED + 1))
    else
        echo -e "${RED}FAIL${NC}"
        VALIDATION_FAILED=$((VALIDATION_FAILED + 1))

        # Capture detailed error
        local error_msg=$(ajv validate -s "$SCHEMA_FILE" -d "$state_file" --strict=false 2>&1 || true)
        VALIDATION_ERRORS+=("$file_name: $error_msg")
    fi
}

# Function to check specific schema violations we know about
check_known_violations() {
    local state_file="$1"
    local file_name=$(basename "$state_file")
    local violations=()

    # Check if waves_completed is an array (should be integer)
    local waves_completed_type=$(jq -r '.project_progression.current_phase.waves_completed | type' "$state_file" 2>/dev/null || echo "null")
    if [[ "$waves_completed_type" == "array" ]]; then
        violations+=("waves_completed is an array (should be integer count)")
    fi

    # Check if total_waves_in_phase is missing
    if ! jq -e '.project_progression.current_phase.total_waves_in_phase' "$state_file" >/dev/null 2>&1; then
        violations+=("total_waves_in_phase field is missing (REQUIRED)")
    fi

    # Check if waves_completed is negative or not a number
    if [[ "$waves_completed_type" == "number" ]]; then
        local waves_completed=$(jq -r '.project_progression.current_phase.waves_completed' "$state_file" 2>/dev/null)
        if [[ "$waves_completed" -lt 0 ]]; then
            violations+=("waves_completed is negative (must be >= 0)")
        fi
    fi

    # Check if total_waves_in_phase exists and is valid
    if jq -e '.project_progression.current_phase.total_waves_in_phase' "$state_file" >/dev/null 2>&1; then
        local total_waves=$(jq -r '.project_progression.current_phase.total_waves_in_phase' "$state_file" 2>/dev/null)
        if [[ "$total_waves" -lt 1 ]]; then
            violations+=("total_waves_in_phase is < 1 (must be >= 1)")
        fi
    fi

    # Print violations if any found
    if [[ ${#violations[@]} -gt 0 ]]; then
        echo ""
        echo -e "${YELLOW}Known violations in $file_name:${NC}"
        for violation in "${violations[@]}"; do
            echo -e "  ${RED}✗${NC} $violation"
        done
        return 1
    fi

    return 0
}

# Main validation logic
main() {
    local target_file="${1:-}"

    if [[ -n "$target_file" ]]; then
        # Validate single file
        echo "Validating single file: $target_file"
        echo ""
        validate_state_file "$target_file"
        check_known_violations "$target_file"
    else
        # Validate all known state files
        echo "Validating all state files in project..."
        echo ""

        local state_files=(
            "$PROJECT_ROOT/orchestrator-state-v3.json"
            "$PROJECT_ROOT/bug-tracking.json"
            "$PROJECT_ROOT/integration-containers.json"
            "$PROJECT_ROOT/fix-cascade-state.json"
        )

        for state_file in "${state_files[@]}"; do
            if [[ -f "$state_file" ]]; then
                validate_state_file "$state_file"
                check_known_violations "$state_file" || true
                echo ""
            fi
        done

        # Check backups too
        echo "Checking recent backups..."
        if [[ -d "$PROJECT_ROOT/.state-backup" ]]; then
            local backup_count=0
            while IFS= read -r backup_file; do
                if [[ $backup_count -lt 3 ]]; then  # Only check 3 most recent
                    validate_state_file "$backup_file"
                    check_known_violations "$backup_file" || true
                    echo ""
                    backup_count=$((backup_count + 1))
                fi
            done < <(find "$PROJECT_ROOT/.state-backup" -name "orchestrator-state-v3.json" -type f | sort -r)
        fi
    fi

    # Print summary
    echo "=================================================="
    echo "  VALIDATION SUMMARY"
    echo "=================================================="
    echo -e "Passed: ${GREEN}$VALIDATION_PASSED${NC}"
    echo -e "Failed: ${RED}$VALIDATION_FAILED${NC}"
    echo ""

    if [[ $VALIDATION_FAILED -gt 0 ]]; then
        echo -e "${RED}VALIDATION ERRORS:${NC}"
        for error in "${VALIDATION_ERRORS[@]}"; do
            echo -e "  ${RED}✗${NC} $error"
        done
        echo ""
        echo -e "${RED}RESULT: VALIDATION FAILED${NC}"
        exit 1
    else
        echo -e "${GREEN}RESULT: ALL VALIDATIONS PASSED${NC}"
        exit 0
    fi
}

# Run main function
main "$@"
