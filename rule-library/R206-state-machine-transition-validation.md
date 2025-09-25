# Rule R206: State Machine Transition Validation Protocol

## Rule Statement
Agents MUST NEVER attempt to transition to states that are not defined in the SOFTWARE-FACTORY-STATE-MACHINE.md. Before ANY state transition, agents MUST read the state machine definition, verify the target state exists for their agent type, and ONLY update current_state if the transition is valid.

## Criticality Level
**BLOCKING** - Invalid state transitions cause system failures and lost work

## Enforcement Mechanism
- **Technical**: Validate target state exists before transition
- **Behavioral**: Exit with error on invalid state attempt
- **Grading**: -50% for invalid state transition (System integrity failure)

## Core Principle

```
State Transition = Read Definition → Validate Target → Verify Path → THEN Transition → IMMEDIATELY CONTINUE!
NEVER transition to a state that doesn't exist in the state machine
NEVER guess or create new states
NEVER STOP after a transition - ALWAYS CONTINUE TO NEXT ACTION!

🚨🚨🚨 CRITICAL: STATE TRANSITIONS ARE WAYPOINTS, NOT DESTINATIONS! 🚨🚨🚨
```

## Detailed Requirements

### MANDATORY State Validation Before Transition

```bash
# ✅✅✅ CORRECT - Validate state before transition
validate_state_transition() {
    local CURRENT_STATE="$1"
    local TARGET_STATE="$2"
    local AGENT_TYPE="$3"  # orchestrator|sw-engineer|code-reviewer|architect
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "🔍 VALIDATING STATE TRANSITION"
    echo "Agent: $AGENT_TYPE"
    echo "Current: $CURRENT_STATE → Target: $TARGET_STATE"
    echo "═══════════════════════════════════════════════════════════════"
    
    # STEP 1: Read state machine definition
    STATE_MACHINE_FILE="SOFTWARE-FACTORY-STATE-MACHINE.md"
    if [ ! -f "$STATE_MACHINE_FILE" ]; then 
        echo "❌ FATAL: State machine definition not found!"; 
        echo "Cannot validate transitions without definition"; 
        exit 1; 
    fi
    
    # STEP 2: Extract valid states for this agent type
    echo "Reading valid states for $AGENT_TYPE..."
    
    # Get the section for this agent's states
    case "$AGENT_TYPE" in
        orchestrator)
            VALID_STATES=$(grep -A 50 "## Orchestrator States" "$STATE_MACHINE_FILE" | 
                          grep "^- " | grep -oE "[A-Z_]+" | sort -u)
            ;;
        sw-engineer)
            VALID_STATES=$(grep -A 50 "## SW Engineer States" "$STATE_MACHINE_FILE" | 
                          grep "^- " | grep -oE "[A-Z_]+" | sort -u)
            ;;
        code-reviewer)
            VALID_STATES=$(grep -A 50 "## Code Reviewer States" "$STATE_MACHINE_FILE" | 
                          grep "^- " | grep -oE "[A-Z_]+" | sort -u)
            ;;
        architect)
            VALID_STATES=$(grep -A 50 "## Architect States" "$STATE_MACHINE_FILE" | 
                          grep "^- " | grep -oE "[A-Z_]+" | sort -u)
            ;;
        *)
            echo "❌ FATAL: Unknown agent type: $AGENT_TYPE"; 
            exit 1
            ;;
    esac
    
    # STEP 3: Verify target state exists
    if echo "$VALID_STATES" | grep -q "^${TARGET_STATE}$"; then 
        echo "✅ Target state '$TARGET_STATE' is valid for $AGENT_TYPE"; 
    else 
        echo "❌ FATAL: State '$TARGET_STATE' does NOT exist for $AGENT_TYPE!"; 
        echo "Valid states are:"; 
        echo "$VALID_STATES" | sed 's/^/  - /'; 
        exit 1; 
    fi
    
    # STEP 4: Verify transition path is allowed
    # Check if there's a defined transition from current to target
    TRANSITION_PATTERN="${CURRENT_STATE}.*→.*${TARGET_STATE}"
    if grep -q "$TRANSITION_PATTERN" "$STATE_MACHINE_FILE"; then 
        echo "✅ Transition path verified: $CURRENT_STATE → $TARGET_STATE"; 
    else 
        echo "⚠️ WARNING: No explicit transition path found"; 
        echo "Verify this transition is allowed in state machine"; 
    fi
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "✅ STATE TRANSITION VALIDATED - Safe to proceed"
    echo "🚀 CONTINUING IMMEDIATELY TO $TARGET_STATE WORK (R021 - NEVER STOP!)"
    echo "⚠️ State transitions are CONTINUATION points, NOT stopping points!"
    echo "═══════════════════════════════════════════════════════════════"
    
    return 0
}

# USAGE: Agents validate their OWN state transitions
# R215: Only orchestrator updates orchestrator-state.json

# For ORCHESTRATOR:
validate_state_transition "WAVE_COMPLETE" "INTEGRATION" "orchestrator"
if [ $? -eq 0 ]; then 
    # Orchestrator updates its own state
    echo "current_state: INTEGRATION" >> orchestrator-state.json; 
fi

# For OTHER AGENTS (SW Eng, Code Reviewer, Architect):
validate_state_transition "INIT" "IMPLEMENTATION" "sw-engineer"
if [ $? -eq 0 ]; then 
    # Agent updates its INTERNAL state (NOT orchestrator-state.json!)
    echo "agent_state: IMPLEMENTATION" > agent-internal-state.yaml; 
    # Report to orchestrator via status file
    echo "status: IMPLEMENTATION" > agent-status.yaml; 
fi
```

### Valid State Lists by Agent Type

```yaml
# ORCHESTRATOR Valid States (Updated with R210, R211)
orchestrator_states:
  - INIT
  - WAVE_START
  - SPAWN_AGENTS
  - SPAWN_ARCHITECT_PHASE_PLANNING  # NEW (R210)
  - SPAWN_ARCHITECT_WAVE_PLANNING   # NEW (R210)
  - SPAWN_CODE_REVIEWER_PHASE_IMPL  # NEW (R211)
  - SPAWN_CODE_REVIEWER_WAVE_IMPL   # NEW (R211)
  - WAITING_FOR_ARCHITECTURE_PLAN   # NEW
  - WAITING_FOR_IMPLEMENTATION_PLAN # NEW
  - INJECT_WAVE_METADATA            # NEW (R213)
  - MONITOR
  - WAVE_COMPLETE
  - INTEGRATION
  - WAVE_REVIEW
  - ERROR_RECOVERY
  - PLANNING
  - SUCCESS
  - HARD_STOP

# SW ENGINEER Valid States
sw_engineer_states:
  - INIT
  - IMPLEMENTATION
  - MEASURE_SIZE
  - FIX_ISSUES
  - SPLIT_IMPLEMENTATION
  - TEST_WRITING
  - COMPLETED
  - BLOCKED

# CODE REVIEWER Valid States (Updated with R211, R214)
code_reviewer_states:
  - INIT
  - PHASE_IMPLEMENTATION_PLANNING   # NEW (R211)
  - WAVE_IMPLEMENTATION_PLANNING    # NEW (R211)
  - WAVE_DIRECTORY_ACKNOWLEDGMENT   # NEW (R214)
  - EFFORT_PLAN_CREATION
  - CODE_REVIEW
  - CREATE_SPLIT_PLAN
  - SPLIT_REVIEW
  - VALIDATION
  - COMPLETED

# ARCHITECT Valid States (Updated with R210, R212)
architect_states:
  - INIT
  - PHASE_ARCHITECTURE_PLANNING     # NEW (R210)
  - WAVE_ARCHITECTURE_PLANNING      # NEW (R210)
  - PHASE_DIRECTORY_ACKNOWLEDGMENT  # NEW (R212)
  - WAVE_REVIEW
  - PHASE_ASSESSMENT
  - INTEGRATION_REVIEW
  - ARCHITECTURE_AUDIT
  - ARCHITECTURE_VALIDATION         # NEW
  - DECISION
```

### State Update Function with Validation

```bash
# Safe state update function
update_state_safely() {
    local FIELD="$1"
    local NEW_VALUE="$2"
    
    # Special handling for current_state field
    if [ "$FIELD" == "current_state" ]; then 
        # Get current state
        CURRENT=$(grep "current_state:" orchestrator-state.json | awk '{print $2}'); 
        
        # Determine agent type from context
        AGENT_TYPE="orchestrator"  # Or detect from prompt/context
        
        # Validate transition
        validate_state_transition "$CURRENT" "$NEW_VALUE" "$AGENT_TYPE"; 
        if [ $? -ne 0 ]; then 
            echo "❌ State transition blocked by R206!"; 
            exit 1; 
        fi; 
    fi
    
    # Proceed with update
    sed -i "s/^${FIELD}:.*/${FIELD}: ${NEW_VALUE}/" orchestrator-state.json
    echo "✅ Updated $FIELD to $NEW_VALUE"
}
```

## Common Violations to Avoid

### ❌ Transitioning to Non-Existent State
```bash
# WRONG - Making up a state
update_state "current_state" "WAITING_FOR_COFFEE"  # This state doesn't exist!

# WRONG - Typo in state name
update_state "current_state" "WAVE_COMPELTE"  # Should be WAVE_COMPLETE
```

### ❌ Skipping Validation
```bash
# WRONG - Direct update without validation
echo "current_state: SOME_STATE" >> orchestrator-state.json
```

### ❌ Wrong Agent State
```bash
# WRONG - SW Engineer trying to use Orchestrator state
# SW Engineer attempting:
update_state "current_state" "WAVE_COMPLETE"  # Only orchestrator has this state!
```

### ✅ Correct State Transition
```bash
# RIGHT - Validate first, then transition
validate_state_transition "$CURRENT_STATE" "INTEGRATION" "orchestrator"
if [ $? -eq 0 ]; then 
    update_state "current_state" "INTEGRATION"; 
else 
    echo "Cannot transition to INTEGRATION - invalid state"; 
    exit 1; 
fi
```

## State Machine Reading Protocol

```bash
# Every agent MUST read the state machine on startup
read_state_machine() {
    local AGENT_TYPE="$1"
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "📖 READING STATE MACHINE DEFINITION (R206)"
    echo "═══════════════════════════════════════════════════════════════"
    
    if [ ! -f "SOFTWARE-FACTORY-STATE-MACHINE.md" ]; then 
        echo "❌ FATAL: State machine definition missing!"; 
        exit 1; 
    fi
    
    # Extract and display valid states for this agent
    echo "Valid states for $AGENT_TYPE:"
    
    case "$AGENT_TYPE" in
        orchestrator)
            grep -A 20 "## Orchestrator States" SOFTWARE-FACTORY-STATE-MACHINE.md | 
                grep "^- " | sed 's/^/  /'
            ;;
        sw-engineer)
            grep -A 20 "## SW Engineer States" SOFTWARE-FACTORY-STATE-MACHINE.md | 
                grep "^- " | sed 's/^/  /'
            ;;
        code-reviewer)
            grep -A 20 "## Code Reviewer States" SOFTWARE-FACTORY-STATE-MACHINE.md | 
                grep "^- " | sed 's/^/  /'
            ;;
        architect)
            grep -A 20 "## Architect States" SOFTWARE-FACTORY-STATE-MACHINE.md | 
                grep "^- " | sed 's/^/  /'
            ;;
    esac
    
    echo "═══════════════════════════════════════════════════════════════"
}

# Call on agent startup
read_state_machine "orchestrator"
```

## Integration with Other Rules

- **R203**: State-aware startup includes validation
- **R186**: Compaction recovery validates resumed state
- **State Machine Definition**: Single source of truth for valid states
- **R001**: Preflight checks include state validation

## Grading Impact

- **Invalid state transition**: -50% (System integrity failure)
- **No validation before transition**: -30% (Protocol violation)
- **Using wrong agent's states**: -40% (Role violation)
- **Creating new states**: -45% (Architecture violation)
- **Typos in state names**: -25% (Quality failure)

## Python Validation Example

```python
def validate_state_transition(current_state, target_state, agent_type):
    """Validate state transition against state machine definition"""
    
    # Define valid states per agent type
    VALID_STATES = {
        'orchestrator': [
            'INIT', 'WAVE_START', 'SPAWN_AGENTS', 'MONITOR',
            'WAVE_COMPLETE', 'INTEGRATION', 'WAVE_REVIEW',
            'ERROR_RECOVERY', 'PLANNING', 'SUCCESS', 'HARD_STOP'
        ],
        'sw-engineer': [
            'INIT', 'IMPLEMENTATION', 'MEASURE_SIZE', 'FIX_ISSUES',
            'SPLIT_IMPLEMENTATION', 'TEST_WRITING', 'COMPLETED', 'BLOCKED'
        ],
        'code-reviewer': [
            'INIT', 'EFFORT_PLAN_CREATION', 'CODE_REVIEW',
            'CREATE_SPLIT_PLAN', 'SPLIT_REVIEW', 'VALIDATION', 'COMPLETED'
        ],
        'architect': [
            'INIT', 'WAVE_REVIEW', 'PHASE_ASSESSMENT',
            'INTEGRATION_REVIEW', 'ARCHITECTURE_AUDIT', 'DECISION'
        ]
    }
    
    # Check agent type is valid
    if agent_type not in VALID_STATES:
        raise ValueError(f"Unknown agent type: {agent_type}")
    
    # Check target state exists for this agent
    valid_for_agent = VALID_STATES[agent_type]
    if target_state not in valid_for_agent:
        raise ValueError(
            f"State '{target_state}' is not valid for {agent_type}. "
            f"Valid states are: {', '.join(valid_for_agent)}"
        )
    
    print(f"✅ Transition validated: {current_state} → {target_state}")
    return True

# Usage
try:
    validate_state_transition("WAVE_COMPLETE", "INTEGRATION", "orchestrator")
    # Safe to proceed with transition
except ValueError as e:
    print(f"❌ Invalid transition: {e}")
    sys.exit(1)
```

## Summary

**Remember**:
- NEVER transition to a state not in the state machine
- ALWAYS validate target state exists for your agent type
- READ the state machine definition before transitions
- VERIFY the transition path is allowed
- EXIT with error if validation fails
- **NEVER STOP after transition - IMMEDIATELY CONTINUE** (R021)
- **State transitions are WAYPOINTS, not DESTINATIONS**
- **Stopping after transition = AUTOMATIC FAILURE**
- This prevents system corruption and lost work