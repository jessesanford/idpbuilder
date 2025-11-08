#!/usr/bin/env bash
#
# validate-test-state-names.sh
# Purpose: Validate that all state names referenced in runtime tests exist in SF 3.0 state machine
# Usage: bash utilities/validate-test-state-names.sh
#
# Exit 0: All state names are valid
# Exit 1: Invalid state names found

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
STATE_MACHINE_FILE="$PROJECT_ROOT/state-machines/software-factory-3.0-state-machine.json"
TEST_DIR="$PROJECT_ROOT/tests"

echo "=== Runtime Test State Name Validation ==="
echo "State Machine: $STATE_MACHINE_FILE"
echo "Test Directory: $TEST_DIR"
echo ""

# Check if state machine exists
if [ ! -f "$STATE_MACHINE_FILE" ]; then
    echo "ERROR: State machine file not found: $STATE_MACHINE_FILE"
    exit 1
fi

# Extract all state names from state machine (keys under "states" object)
echo "Extracting valid state names from state machine..."
VALID_STATES=$(jq -r '.states | keys[]' "$STATE_MACHINE_FILE" | sort)
VALID_STATE_COUNT=$(echo "$VALID_STATES" | wc -l)
echo "Found $VALID_STATE_COUNT valid states in state machine"
echo ""

# Create temporary file for valid states
VALID_STATES_FILE=$(mktemp)
echo "$VALID_STATES" > "$VALID_STATES_FILE"

# Find all runtime test scripts
TEST_SCRIPTS=$(find "$TEST_DIR" -name "runtime-test-*.sh" -type f | grep -v "framework" | grep -v "monitor" | sort)
TEST_COUNT=$(echo "$TEST_SCRIPTS" | wc -l)

if [ -z "$TEST_SCRIPTS" ]; then
    echo "ERROR: No runtime test scripts found in $TEST_DIR"
    rm -f "$VALID_STATES_FILE"
    exit 1
fi

echo "Found $TEST_COUNT runtime test scripts to validate"
echo ""

# Track validation results
INVALID_FOUND=0
TOTAL_STATES_CHECKED=0

# Process each test script
for test_script in $TEST_SCRIPTS; do
    test_name=$(basename "$test_script")
    echo "Validating: $test_name"

    # Extract state names from nominal_path arrays in test scripts
    # Look for lines like: "STATE_NAME"
    # Inside nominal_path=( ... ) declarations

    # Use awk to extract states from nominal_path arrays
    SCRIPT_STATES=$(awk '
        /nominal_path=\(/ { in_array=1; next }
        in_array && /^\s*\)/ { in_array=0; next }
        in_array && /^\s*"[A-Z_]+"/ {
            gsub(/^[[:space:]]*"/, "")
            gsub(/"[[:space:]]*$/, "")
            gsub(/".*/, "")
            print
        }
    ' "$test_script" | sort -u)

    if [ -z "$SCRIPT_STATES" ]; then
        echo "  WARNING: No states found in nominal_path array"
        continue
    fi

    STATE_COUNT=$(echo "$SCRIPT_STATES" | wc -l)
    echo "  Found $STATE_COUNT unique states in nominal_path"
    TOTAL_STATES_CHECKED=$((TOTAL_STATES_CHECKED + STATE_COUNT))

    # Check each state against valid states
    SCRIPT_INVALID=0
    while IFS= read -r state; do
        if ! grep -Fxq "$state" "$VALID_STATES_FILE"; then
            echo "  ERROR: Invalid state name: $state"
            SCRIPT_INVALID=$((SCRIPT_INVALID + 1))
            INVALID_FOUND=$((INVALID_FOUND + 1))
        fi
    done <<< "$SCRIPT_STATES"

    if [ $SCRIPT_INVALID -eq 0 ]; then
        echo "  OK: All $STATE_COUNT states are valid"
    else
        echo "  FAIL: $SCRIPT_INVALID invalid states found"
    fi
    echo ""
done

# Cleanup
rm -f "$VALID_STATES_FILE"

# Summary
echo "=== Validation Summary ==="
echo "Total test scripts checked: $TEST_COUNT"
echo "Total state references checked: $TOTAL_STATES_CHECKED"
echo "Invalid state names found: $INVALID_FOUND"
echo ""

if [ $INVALID_FOUND -eq 0 ]; then
    echo "PROJECT_DONE: All runtime test state names are valid SF 3.0 states"
    exit 0
else
    echo "FAILURE: $INVALID_FOUND invalid state names must be fixed"
    exit 1
fi
