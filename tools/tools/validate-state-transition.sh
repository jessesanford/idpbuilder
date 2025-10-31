#!/bin/bash
# tools/validate-state-transition.sh
# Validates state transitions against state machine allowed_transitions
# Used by State Manager and pre-commit hooks

set -euo pipefail

CURRENT_STATE="${1:-}"
PROPOSED_NEXT_STATE="${2:-}"
STATE_MACHINE="${3:-state-machines/software-factory-3.0-state-machine.json}"

if [ -z "$CURRENT_STATE" ] || [ -z "$PROPOSED_NEXT_STATE" ]; then
    echo "Usage: $0 <current_state> <proposed_next_state> [state_machine_file]"
    exit 1
fi

# Validate state machine file exists
if [ ! -f "$STATE_MACHINE" ]; then
    echo "❌ CRITICAL: State machine file not found: $STATE_MACHINE"
    exit 1
fi

# Log validation attempt
LOG_FILE="state-transition-validation.log"
echo "$(date -Iseconds): Validating $CURRENT_STATE → $PROPOSED_NEXT_STATE" >> "$LOG_FILE"

# Check current state exists in state machine
if ! jq -e ".states.\"$CURRENT_STATE\"" "$STATE_MACHINE" > /dev/null 2>&1; then
    echo "❌ CRITICAL: Current state '$CURRENT_STATE' not found in state machine!"
    echo "$(date -Iseconds): ERROR - Current state not found: $CURRENT_STATE" >> "$LOG_FILE"
    exit 1
fi

# Check proposed state exists in state machine
if ! jq -e ".states.\"$PROPOSED_NEXT_STATE\"" "$STATE_MACHINE" > /dev/null 2>&1; then
    echo "❌ CRITICAL: Proposed state '$PROPOSED_NEXT_STATE' not found in state machine!"
    echo "$(date -Iseconds): ERROR - Proposed state not found: $PROPOSED_NEXT_STATE" >> "$LOG_FILE"
    exit 1
fi

# Get allowed transitions
ALLOWED_TRANSITIONS=$(jq -r ".states.\"$CURRENT_STATE\".allowed_transitions[]" "$STATE_MACHINE" 2>/dev/null || echo "")

if [ -z "$ALLOWED_TRANSITIONS" ]; then
    echo "❌ CRITICAL: No allowed_transitions defined for state '$CURRENT_STATE'!"
    echo "$(date -Iseconds): ERROR - No allowed_transitions: $CURRENT_STATE" >> "$LOG_FILE"
    exit 1
fi

# Check if proposed transition is allowed
if ! jq -e ".states.\"$CURRENT_STATE\".allowed_transitions | contains([\"$PROPOSED_NEXT_STATE\"])" \
        "$STATE_MACHINE" > /dev/null 2>&1; then

    echo "❌❌❌ CRITICAL: ILLEGAL STATE TRANSITION! ❌❌❌"
    echo ""
    echo "Current state: $CURRENT_STATE"
    echo "Proposed next state: $PROPOSED_NEXT_STATE"
    echo ""
    echo "Allowed transitions from state machine:"
    echo "$ALLOWED_TRANSITIONS" | while read -r allowed; do
        echo "  ✅ $allowed"
    done
    echo ""
    echo "The transition $CURRENT_STATE → $PROPOSED_NEXT_STATE is NOT ALLOWED!"
    echo ""
    echo "FIX: Choose one of the allowed transitions above"
    echo ""

    # Log error
    echo "$(date -Iseconds): BLOCKED ILLEGAL TRANSITION: $CURRENT_STATE → $PROPOSED_NEXT_STATE" >> "$LOG_FILE"
    echo "  Allowed: $ALLOWED_TRANSITIONS" >> "$LOG_FILE"

    exit 1
fi

# Validation passed
echo "✅ STATE TRANSITION VALIDATED"
echo "  From: $CURRENT_STATE"
echo "  To: $PROPOSED_NEXT_STATE"
echo "  State machine: $STATE_MACHINE"

echo "$(date -Iseconds): VALIDATED: $CURRENT_STATE → $PROPOSED_NEXT_STATE" >> "$LOG_FILE"

exit 0
