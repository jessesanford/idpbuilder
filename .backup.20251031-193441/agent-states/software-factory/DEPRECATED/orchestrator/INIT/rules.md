# INIT State Rules

## State Purpose
Initial orchestrator startup - establishes system baseline, validates environment, and determines entry point into state machine.

## Entry Conditions
- Fresh system start OR
- Recovery from ERROR_RECOVERY with clean slate OR
- Explicit reset to INIT state

## Mandatory Validations

### 1. Environment Validation
```bash
# Validate critical directories exist
echo "Validating Software Factory environment..."

# Check for state machine definition
if [[ ! -f "$CLAUDE_PROJECT_DIR/state-machines/software-factory-3.0-state-machine.json" ]]; then
    echo "ERROR: State machine definition not found!"
    exit 1
fi

# Check for orchestrator state file
if [[ ! -f "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" ]]; then
    echo "WARNING: orchestrator-state-v3.json not found - will create from template"
    if [[ -f "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json.example" ]]; then
        cp "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json.example" "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
        echo "Created orchestrator-state-v3.json from template"
    else
        echo "ERROR: No orchestrator-state-v3.json or template found!"
        exit 1
    fi
fi

# Validate agent configs exist
for agent in orchestrator sw-engineer code-reviewer architect; do
    if [[ ! -f "$CLAUDE_PROJECT_DIR/.claude/agents/${agent}.md" ]]; then
        echo "ERROR: Agent config missing for ${agent}"
        exit 1
    fi
done
```

### 2. Load System Metadata
```bash
# Load project configuration
if [[ -f "$CLAUDE_PROJECT_DIR/target-repo-config.yaml" ]]; then
    echo "Loading target repository configuration..."
    TARGET_REPO=$(yq '.repository' "$CLAUDE_PROJECT_DIR/target-repo-config.yaml")
    echo "Target repository: ${TARGET_REPO}"
fi

# Get current state from orchestrator-state-v3.json
CURRENT_STATE=$(jq -r '.state_machine.current_state' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")
echo "Current state in file: ${CURRENT_STATE}"

# If state is already non-INIT, this is a recovery scenario
if [[ "$CURRENT_STATE" != "INIT" ]] && [[ "$CURRENT_STATE" != "null" ]]; then
    echo "WARNING: State file shows ${CURRENT_STATE}, but we're in INIT"
    echo "This indicates a restart or recovery scenario"
fi
```

### 3. Determine System Status
```bash
# Check for existing work
CURRENT_PHASE=$(jq -r '.current_phase' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")
CURRENT_WAVE=$(jq -r '.current_wave' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")

if [[ "$CURRENT_PHASE" == "null" ]] || [[ "$CURRENT_PHASE" == "0" ]]; then
    echo "System is uninitialized - need to start Phase 1"
    NEXT_STATE="START_PHASE_ITERATION"
    UPDATE_FIELDS='.current_phase = 1 | .current_wave = 0'
else
    echo "System has existing progress:"
    echo "  Current Phase: ${CURRENT_PHASE}"
    echo "  Current Wave: ${CURRENT_WAVE}"

    # Check if we're resuming from a specific state
    if [[ "$CURRENT_STATE" != "INIT" ]] && [[ "$CURRENT_STATE" != "null" ]]; then
        echo "Resuming from state: ${CURRENT_STATE}"
        NEXT_STATE="$CURRENT_STATE"
        UPDATE_FIELDS='.'  # No field updates needed
    else
        # Determine next state based on phase/wave
        if [[ "$CURRENT_WAVE" == "0" ]] || [[ "$CURRENT_WAVE" == "null" ]]; then
            NEXT_STATE="START_PHASE_ITERATION"
        else
            NEXT_STATE="WAVE_START"
        fi
        UPDATE_FIELDS='.'  # No field updates needed
    fi
fi
```

## State Actions

### 1. Initialize State File
```bash
# Ensure state file has proper structure
echo "Initializing orchestrator state file..."

# Update state to INIT (temporary)
jq '.state_machine.current_state = "INIT"' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" > /tmp/state.json
mv /tmp/state.json "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"

# Apply any necessary field updates
if [[ "$UPDATE_FIELDS" != "." ]]; then
    jq "$UPDATE_FIELDS" "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" > /tmp/state.json
    mv /tmp/state.json "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
fi

# Ensure required directories exist
mkdir -p "$CLAUDE_PROJECT_DIR/planning"
mkdir -p "$CLAUDE_PROJECT_DIR/protocols"
mkdir -p "$CLAUDE_PROJECT_DIR/efforts"
mkdir -p "$CLAUDE_PROJECT_DIR/todos"

echo "State file initialized"
```

### 2. Create Startup Report
```bash
cat <<EOF
================================================================
ORCHESTRATOR INITIALIZATION COMPLETE
================================================================
Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')
State Machine: Loaded
Configuration: Valid
Current Phase: ${CURRENT_PHASE:-1}
Current Wave: ${CURRENT_WAVE:-0}
Next State: ${NEXT_STATE}
================================================================
EOF
```

## Exit Conditions

### Success Criteria
- Environment validated
- State file initialized/validated
- Next state determined
- All required directories created

### State Transitions
- **START_PHASE_ITERATION**: When starting new phase (Phase 1 or resuming)
- **WAVE_START**: When resuming mid-phase with existing wave
- **ERROR_RECOVERY**: If critical validation fails
- **[CURRENT_STATE]**: If resuming from specific state

### State Update Requirements
```bash
# Update state file with next state
update_state() {
    local next_state="$1"
    local notes="${2:-Transitioning from INIT}"

    jq --arg state "$next_state" \
       --arg notes "$notes" \
       --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       '.state_machine.current_state = $state |
        .last_transition = {
            from: "INIT",
            to: $state,
            timestamp: $timestamp,
            notes: $notes
        }' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" > /tmp/state.json

    mv /tmp/state.json "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
    echo "State updated to: ${next_state}"
}

# Transition to next state
update_state "${NEXT_STATE}" "Initialization complete"
```

## Associated Rules
- **R290**: State rule reading verification (SUPREME LAW)
- **R203**: State-aware agent startup
- **R206**: State machine validation
- **R288**: State file update requirements
- **R233**: Single operation per state (SUPREME LAW)

## Prohibitions
- ❌ Skip environment validation
- ❌ Proceed without valid state file
- ❌ Transition without determining proper next state

## Automation Flag

```bash
# After successful initialization and validation:
echo "✅ Initialization complete, proceeding to next state"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue to appropriate next state
```
- ❌ Ignore existing work in progress
- ❌ Create duplicate work

## Notes
- INIT is the universal entry point for fresh starts
- Handles both new projects and recovery scenarios
- Must gracefully handle partial state files
- Always validates before proceeding
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
