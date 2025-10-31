# Rule R206: State Machine Transition Validation Protocol

## Rule Statement
Agents MUST NEVER attempt to transition to states that are not defined in the software-factory-3.0-state-machine.json. Before ANY state transition, agents MUST read the state machine definition, verify the target state exists for their agent type, and ONLY update current_state if the transition is valid.

## Criticality Level
```bash
validate_state_transition() {
    local AGENT_TYPE="$1"
    local CURRENT_STATE="$2"
    local TARGET_STATE="$3"

    # STEP 1: Read state machine definition
    STATE_MACHINE_FILE="$CLAUDE_PROJECT_DIR/state-machines/software-factory-3.0-state-machine.json"
    if [ ! -f "$STATE_MACHINE_FILE" ]; then
        echo "❌ FATAL: State machine missing!"
        return 1
    fi

    # STEP 2: Verify target state exists for this agent type
    if ! jq -e ".agents.\"${AGENT_TYPE}\".states.\"${TARGET_STATE}\"" "$STATE_MACHINE_FILE" >/dev/null 2>&1; then
        echo "❌ Invalid state '${TARGET_STATE}' for agent '${AGENT_TYPE}'"
        return 1
    fi

    # STEP 3: Check if transition is valid
    VALID_TRANSITIONS=$(jq -r ".transition_matrix.\"${AGENT_TYPE}\".\"${CURRENT_STATE}\"[]?" "$STATE_MACHINE_FILE" 2>/dev/null)
    if [ -z "$VALID_TRANSITIONS" ]; then
        echo "❌ No valid transitions from '${CURRENT_STATE}'"
        return 1
    fi

    if ! echo "$VALID_TRANSITIONS" | grep -q "^${TARGET_STATE}$"; then
        echo "❌ Invalid transition: ${CURRENT_STATE} -> ${TARGET_STATE}"
        echo "Valid transitions: $(echo $VALID_TRANSITIONS | tr '\n' ', ')"
        return 1
    fi

    echo "✅ Valid transition: ${CURRENT_STATE} -> ${TARGET_STATE}"
    return 0
}
```

### List Valid States for Agent
```bash
list_valid_states() {
    local AGENT_TYPE="$1"

    if [ ! -f "$CLAUDE_PROJECT_DIR/state-machines/software-factory-3.0-state-machine.json" ]; then
        echo "❌ FATAL: State machine definition missing!";
        return 1
    fi

    echo "Valid states for ${AGENT_TYPE}:"
    jq -r ".agents.\"${AGENT_TYPE}\".states | keys[]" "$CLAUDE_PROJECT_DIR/state-machines/software-factory-3.0-state-machine.json" 2>/dev/null | sed 's/^/  /'
}
```

## State Manager Coordination (SF 3.0)

State Manager validates all transitions during shutdown consultation:
- **Checks allowed_transitions** from state-machine JSON before committing new state
- **Verifies guard conditions** (requires.conditions) are met
- **Ensures atomic 4-file updates** via `tools/atomic-state-update.sh`
- **Rolls back on validation failure** to prevent inconsistent state

The State Manager's SHUTDOWN_CONSULTATION state implements this rule automatically.

See: R288 (atomic updates), `agent-states/state-manager/SHUTDOWN_CONSULTATION/rules.md`, `state-machines/software-factory-3.0-state-machine.json`
