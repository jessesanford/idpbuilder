# START_PHASE_ITERATION State Rules

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
Initialize new phase, validate phase metadata, and determine if phase-level planning is needed per R210.

## Entry Conditions
- Previous phase complete (COMPLETE_PHASE) OR starting Phase 1 (INIT)
- orchestrator-state-v3.json current_phase field set
- Phase metadata exists in orchestrator-state-v3.json

## Mandatory Validations

### 1. Validate Phase Metadata
```bash
# Get current phase number
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
echo "Initializing Phase ${PHASE}..."

# Validate phase metadata exists
PHASE_STATUS=$(jq -r ".phases.phase_${PHASE}.status" orchestrator-state-v3.json)

if [ "$PHASE_STATUS" = "null" ]; then
    echo "ERROR: Phase ${PHASE} metadata missing in orchestrator-state-v3.json"
    update_state "ERROR_RECOVERY" "Missing phase metadata for Phase ${PHASE}"
    exit 1
fi

# Validate phase status is "planned" or "in_progress"
if [ "$PHASE_STATUS" != "planned" ] && [ "$PHASE_STATUS" != "in_progress" ]; then
    echo "ERROR: Phase ${PHASE} has invalid status: ${PHASE_STATUS}"
    update_state "ERROR_RECOVERY" "Invalid phase status: ${PHASE_STATUS}"
    exit 1
fi
```

### 2. Check for Existing Phase Plans (R502)
```bash
# Check if phase plans already exist
echo "Checking for existing phase plans..."

# Updated path structure per standard: planning/phaseX/
PLANNING_DIR="$CLAUDE_PROJECT_DIR/planning/phase${PHASE}"
mkdir -p "$PLANNING_DIR"

# Check for phase architecture plan
ARCH_PLAN=$(ls -t "${PLANNING_DIR}/PHASE-${PHASE}-ARCHITECTURE-PLAN--"*.md 2>/dev/null | head -1)

# Check for phase implementation plan
IMPL_PLAN=$(ls -t "${PLANNING_DIR}/PHASE-${PHASE}-PLAN--"*.md 2>/dev/null | head -1)

if [ -n "$ARCH_PLAN" ] && [ -n "$IMPL_PLAN" ]; then
    echo "Phase plans already exist:"
    echo "  Architecture: $(basename "$ARCH_PLAN")"
    echo "  Implementation: $(basename "$IMPL_PLAN")"

    # Update state file with plan locations
    jq --arg arch "$ARCH_PLAN" --arg impl "$IMPL_PLAN" \
       ".phases.phase_${PHASE}.architecture_plan = \$arch | .phases.phase_${PHASE}.implementation_plan = \$impl" \
       orchestrator-state-v3.json > orchestrator-state.tmp && mv orchestrator-state.tmp orchestrator-state-v3.json

    echo "Proceeding directly to WAVE_START"
    update_state "WAVE_START" "Phase ${PHASE} plans already exist"
    exit 0
fi

# If either plan is missing, both must be created
if [ -n "$ARCH_PLAN" ] || [ -n "$IMPL_PLAN" ]; then
    echo "WARNING: Partial phase plans detected. Both must exist or neither."
    echo "  Architecture exists: $([ -n "$ARCH_PLAN" ] && echo "yes" || echo "no")"
    echo "  Implementation exists: $([ -n "$IMPL_PLAN" ] && echo "yes" || echo "no")"
    echo "Spawning architect to create complete phase plans..."
fi
```

### 3. Spawn Architect for Phase Planning (R210)
```bash
echo "═══════════════════════════════════════════════════════"
echo "R210: Spawning architect for Phase ${PHASE} planning..."
echo "═══════════════════════════════════════════════════════"

# Prepare architect instructions
cat > /tmp/architect-phase-planning-instructions.md << 'EOF'
# Phase Planning Instructions for Architect

You are being spawned to create phase-level architecture and implementation plans per R210.

## Your Task:
Create comprehensive phase planning documents for Phase ${PHASE}.

## Required Outputs:
1. **Phase Architecture Plan**: `$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/PHASE-${PHASE}-ARCHITECTURE-PLAN--$(date +%Y%m%d-%H%M%S).md`
   - Define phase-wide architectural decisions
   - Specify APIs, contracts, and interfaces
   - Document shared abstractions and patterns
   - List reusable libraries from previous phases

2. **Phase Implementation Plan**: `$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/PHASE-${PHASE}-PLAN--$(date +%Y%m%d-%H%M%S).md`
   - Break down phase into waves
   - Define wave objectives and dependencies
   - Specify effort sequencing within waves
   - Document parallelization opportunities

## Required Reading:
1. Master PROJECT-IMPLEMENTATION-PLAN.md
2. orchestrator-state-v3.json for phase metadata
3. Previous phase assessments (if any)
4. Existing code and architecture from completed phases

## Mission Critical per R210:
This phase planning MUST be complete before any wave planning can begin.
Wave planning depends on the phase architecture established here.

## State Transitions:
After creating both plans, the orchestrator will transition to WAVE_START.
EOF

# Spawn the architect
/spawn architect PHASE_ARCHITECTURE_PLANNING /tmp/architect-phase-planning-instructions.md

# Update phase status
jq ".phases.phase_${PHASE}.status = \"planning\"" orchestrator-state-v3.json > orchestrator-state.tmp && \
    mv orchestrator-state.tmp orchestrator-state-v3.json

# Transition to waiting state
update_state "WAITING_FOR_PHASE_PLANS" "Architect spawned for Phase ${PHASE} planning"
```

## Exit Conditions
- Phase metadata validated
- Decision made: spawn architect OR proceed to WAVE_START
- orchestrator-state-v3.json updated with decision
- State transitioned appropriately

## State Transitions
- **SPAWN_ARCHITECT_PHASE_PLANNING**: When phase plans don't exist
- **WAVE_START**: When valid phase plans already exist
- **ERROR_RECOVERY**: When validation fails or critical error occurs

## Error Handling
```bash
# Common error conditions
handle_phase_start_error() {
    local error_type="$1"
    local error_details="$2"

    case "$error_type" in
        "missing_metadata")
            echo "ERROR: Phase metadata not found in orchestrator-state-v3.json"
            echo "Details: $error_details"
            update_state "ERROR_RECOVERY" "Missing phase metadata: $error_details"
            ;;
        "invalid_phase")
            echo "ERROR: Invalid phase number or phase out of sequence"
            echo "Details: $error_details"
            update_state "ERROR_RECOVERY" "Invalid phase: $error_details"
            ;;
        "corrupted_state")
            echo "ERROR: orchestrator-state-v3.json is corrupted or invalid"
            echo "Details: $error_details"
            update_state "ERROR_RECOVERY" "Corrupted state: $error_details"
            ;;
        *)
            echo "ERROR: Unknown error in START_PHASE_ITERATION"
            echo "Type: $error_type"
            echo "Details: $error_details"
            update_state "ERROR_RECOVERY" "Unknown error: $error_type - $error_details"
            ;;
    esac
}
```

## Grading Criteria Compliance
- R210: Phase architecture planning happens BEFORE wave planning
- R502: Phase plan validation gates are enforced
- R206: State machine transitions are validated
- R313: Agent stops after spawning architect

## Automation Flag
```bash
# After all operations complete successfully:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
