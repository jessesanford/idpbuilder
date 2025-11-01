# WAVE_START State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## State Purpose
Initialize a new wave within the current phase, validate wave metadata, and prepare for effort planning and execution.

## Entry Conditions
- From START_PHASE_ITERATION after phase planning complete OR
- From WAVE_COMPLETE after previous wave integrated OR
- From WAITING_FOR_WAVE_PLANS after plans created

## Mandatory Validations

### 1. Validate Wave Context
```bash
# Get current phase and wave
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

echo "Initializing Phase ${PHASE}, Wave ${WAVE}..."

# Validate wave metadata exists
WAVE_STATUS=$(jq -r ".phases.phase_${PHASE}.waves.wave_${WAVE}.status" orchestrator-state-v3.json)

if [ "$WAVE_STATUS" = "null" ]; then
    echo "ERROR: Wave ${WAVE} metadata missing for Phase ${PHASE}"
    update_state "ERROR_RECOVERY" "Missing wave metadata"
    exit 1
fi

# Validate wave status
if [ "$WAVE_STATUS" != "planned" ] && [ "$WAVE_STATUS" != "in_progress" ]; then
    echo "ERROR: Wave ${WAVE} has invalid status: ${WAVE_STATUS}"
    update_state "ERROR_RECOVERY" "Invalid wave status"
    exit 1
fi
```

### 2. Check for Wave Plans (R507)
```bash
# Check if wave plans exist
PLANNING_DIR="$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/wave${WAVE}"
mkdir -p "$PLANNING_DIR"

# Look for wave plan
WAVE_PLAN=$(ls -t "${PLANNING_DIR}/WAVE-${WAVE}-PLAN--"*.md 2>/dev/null | head -1)

if [[ -z "$WAVE_PLAN" ]]; then
    echo "No wave plan found - need architect planning"
    NEEDS_PLANNING=true
else
    echo "Found wave plan: $(basename "$WAVE_PLAN")"
    NEEDS_PLANNING=false
fi
```

### 3. Check for Effort Plans
```bash
# Count planned efforts for this wave
EFFORT_COUNT=$(jq -r ".phases.phase_${PHASE}.waves.wave_${WAVE}.efforts | length" orchestrator-state-v3.json)

if [[ "$EFFORT_COUNT" -eq 0 ]] || [[ "$EFFORT_COUNT" == "null" ]]; then
    echo "No efforts defined for Wave ${WAVE}"
    NEEDS_EFFORT_PLANNING=true
else
    echo "Found ${EFFORT_COUNT} efforts defined"
    NEEDS_EFFORT_PLANNING=false
fi
```

## State Actions

### 1. Determine Next Action
```bash
if [[ "$NEEDS_PLANNING" == "true" ]]; then
    echo "Wave needs architect planning"
    NEXT_STATE="SPAWN_ARCHITECT_WAVE_PLANNING"
elif [[ "$NEEDS_EFFORT_PLANNING" == "true" ]]; then
    echo "Wave needs effort planning"
    NEXT_STATE="SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
else
    echo "Wave is ready for implementation"
    NEXT_STATE="ANALYZE_IMPLEMENTATION_PARALLELIZATION"
fi
```

### 2. Update Wave Status
```bash
# Mark wave as in_progress
jq --arg phase "$PHASE" \
   --arg wave "$WAVE" \
   '.phases["phase_" + $phase].waves["wave_" + $wave].status = "in_progress" |
    .phases["phase_" + $phase].waves["wave_" + $wave].started_at = now | todate' \
   orchestrator-state-v3.json > /tmp/state.json

mv /tmp/state.json orchestrator-state-v3.json
```

### 3. Create Wave Report
```bash
cat <<EOF
================================================================
WAVE ${WAVE} INITIALIZATION (Phase ${PHASE})
================================================================
Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')
Wave Status: ${WAVE_STATUS}
Wave Planning: $(if [[ "$NEEDS_PLANNING" == "true" ]]; then echo "REQUIRED"; else echo "COMPLETE"; fi)
Effort Planning: $(if [[ "$NEEDS_EFFORT_PLANNING" == "true" ]]; then echo "REQUIRED"; else echo "COMPLETE"; fi)
Effort Count: ${EFFORT_COUNT:-0}
Next State: ${NEXT_STATE}
================================================================
EOF
```

## Exit Conditions

### Success Criteria
- Wave metadata validated
- Planning status determined
- Next state identified
- Wave marked as in_progress

### State Transitions
- **SPAWN_ARCHITECT_WAVE_PLANNING**: If wave needs architect planning
- **SPAWN_CODE_REVIEWERS_EFFORT_PLANNING**: If wave needs effort planning
- **ANALYZE_IMPLEMENTATION_PARALLELIZATION**: If ready for implementation
- **ERROR_RECOVERY**: If validation fails

### State Update Requirements
```bash
# Update to next state
update_state() {
    local next_state="$1"
    local notes="${2:-Wave initialization complete}"

    jq --arg state "$next_state" \
       --arg notes "$notes" \
       --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       '.state_machine.current_state = $state |
        .last_transition = {
            from: "WAVE_START",
            to: $state,
            timestamp: $timestamp,
            notes: $notes
        }' orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json
    echo "State updated to: ${next_state}"
}

# Transition to next state
update_state "${NEXT_STATE}" "Wave ${WAVE} initialization complete"
```

## Associated Rules
- **R290**: State rule reading verification (SUPREME LAW)
- **R507**: Wave planning requirements
- **R233**: Single operation per state (SUPREME LAW)
- **R313**: Wave completion requirements
- **R322**: Effort planning protocol

## Prohibitions
- ❌ Start implementation without plans
- ❌ Skip wave validation
- ❌ Mix waves (work on multiple waves simultaneously)
- ❌ Proceed with invalid metadata
- ❌ Create duplicate efforts

## Automation Flag

```bash
# After successful wave initialization:
echo "✅ Wave ${WAVE} initialized, proceeding to next action"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue to wave planning or infrastructure
```

## Notes
- Each wave must be fully planned before implementation
- Efforts within a wave can be parallelized per plan
- Wave completion requires all efforts complete and integrated
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
