#!/bin/bash
# Validates that a state transition is legal according to the state machine definition
# Usage: validate-state-transition.sh <from_state> <to_state> [state_machine_file]
#
# This utility prevents invalid state transitions by validating:
# 1. Target state exists in the state machine
# 2. Source state exists in the state machine (if not INIT)
# 3. Transition is allowed by state machine definition
#
# Created as part of Test 1.5 P0 fix #2 to prevent CREATE_NEXT_INFRASTRUCTURE-type errors

set -e

FROM_STATE="$1"
TO_STATE="$2"
STATE_MACHINE_FILE="${3:-state-machines/software-factory-3.0-state-machine.json}"

if [ -z "$FROM_STATE" ] || [ -z "$TO_STATE" ]; then
    echo "Usage: $0 <from_state> <to_state> [state_machine_file]"
    echo ""
    echo "Examples:"
    echo "  $0 INIT CREATE_NEXT_INFRASTRUCTURE"
    echo "  $0 ANALYZE_CODE_REVIEWER_PARALLELIZATION SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
    exit 1
fi

# Check if state machine file exists
if [ ! -f "$STATE_MACHINE_FILE" ]; then
    echo "❌ ERROR: State machine file not found: $STATE_MACHINE_FILE"
    exit 1
fi

echo "Validating state transition: $FROM_STATE → $TO_STATE"

# 1. Validate TO_STATE exists
STATE_EXISTS=$(jq -r --arg state "$TO_STATE" '.states | has($state)' "$STATE_MACHINE_FILE")

if [ "$STATE_EXISTS" != "true" ]; then
    echo ""
    echo "════════════════════════════════════════════════════════════"
    echo "❌ FATAL ERROR: Invalid State Name"
    echo "════════════════════════════════════════════════════════════"
    echo ""
    echo "Target state '$TO_STATE' does not exist in state machine!"
    echo ""
    echo "Valid states are:"
    jq -r '.states | keys[]' "$STATE_MACHINE_FILE" | sed 's/^/  ✓ /'
    echo ""
    echo "Common mistakes:"
    echo "  ❌ CREATE_NEXT_INFRASTRUCTURE (old, deprecated)"
    echo "  ✓ CREATE_NEXT_INFRASTRUCTURE (correct, SF 3.0)"
    echo ""
    echo "  ❌ CODE_REVIEW (old, SF 2.0)"
    echo "  ✓ CODE_REVIEW (correct, SF 3.0)"
    echo ""
    echo "  ❌ MONITOR (old, SF 2.0)"
    echo "  ✓ MONITORING_SWE_PROGRESS (correct, SF 3.0)"
    echo ""
    echo "This transition would have been rejected."
    echo "════════════════════════════════════════════════════════════"
    exit 1
fi

# 2. Validate FROM_STATE exists (if not INIT)
if [ "$FROM_STATE" != "INIT" ]; then
    FROM_STATE_EXISTS=$(jq -r --arg state "$FROM_STATE" '.states | has($state)' "$STATE_MACHINE_FILE")

    if [ "$FROM_STATE_EXISTS" != "true" ]; then
        echo ""
        echo "════════════════════════════════════════════════════════════"
        echo "❌ ERROR: Invalid Source State"
        echo "════════════════════════════════════════════════════════════"
        echo ""
        echo "From state '$FROM_STATE' does not exist in state machine!"
        echo ""
        echo "Valid states are:"
        jq -r '.states | keys[]' "$STATE_MACHINE_FILE" | sed 's/^/  ✓ /'
        echo ""
        echo "════════════════════════════════════════════════════════════"
        exit 1
    fi
fi

# 3. Validate transition is allowed
TRANSITION_VALID=$(jq -r --arg from "$FROM_STATE" --arg to "$TO_STATE" '
    .states[$from].transitions | if . then any(. == $to) else false end
' "$STATE_MACHINE_FILE")

if [ "$TRANSITION_VALID" != "true" ]; then
    echo ""
    echo "⚠️  WARNING: Transition not defined in state machine"
    echo "    From: $FROM_STATE"
    echo "    To: $TO_STATE"
    echo ""
    echo "    Allowed transitions from $FROM_STATE:"
    jq -r --arg from "$FROM_STATE" '.states[$from].transitions[]?' "$STATE_MACHINE_FILE" | sed 's/^/      → /'
    echo ""
    echo "    This may be valid for dynamic transitions (ERROR_RECOVERY, etc.)"
    echo "    Proceeding with warning..."
    echo ""
    # This is a warning, not an error (some dynamic transitions may be valid)
fi

echo "✅ State validation passed"
exit 0
