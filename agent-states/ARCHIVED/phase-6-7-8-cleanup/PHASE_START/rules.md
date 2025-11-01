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
## 🔴🔴🔴 CRITICAL: STATE FILES PROPOSE, STATE MANAGER DECIDES 🔴🔴🔴

**R288 SUPREME LAW - YOU NEVER CALL update_state DIRECTLY!**

This state file uses the **BOOKEND PATTERN** (R600):
1. **START**: Set `PROPOSED_NEXT_STATE` and `TRANSITION_REASON` variables
2. **WORK**: Execute state-specific logic
3. **END**: Follow 11-step completion checklist to exit properly

**NEVER CALL update_state() DIRECTLY - IT IS PROHIBITED!**
- ❌ `update_state "NEXT_STATE" "reason"` = SYSTEM VIOLATION
- ✅ `PROPOSED_NEXT_STATE="NEXT_STATE"` = CORRECT
- ✅ `TRANSITION_REASON="reason"` = CORRECT

The State Manager (`run-software-factory.sh`) handles ALL state transitions.

**See:**
- `$CLAUDE_PROJECT_DIR/rule-library/R288-state-transition-authority.md` (SUPREME LAW)
- `$CLAUDE_PROJECT_DIR/rule-library/R600-orchestrator-bookend-pattern.md`

---

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
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Missing phase metadata for Phase ${PHASE}"
    exit 1
fi

# Validate phase status is "planned" or "in_progress"
if [ "$PHASE_STATUS" != "planned" ] && [ "$PHASE_STATUS" != "in_progress" ]; then
    echo "ERROR: Phase ${PHASE} has invalid status: ${PHASE_STATUS}"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Invalid phase status: ${PHASE_STATUS}"
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
    PROPOSED_NEXT_STATE="WAVE_START"
    TRANSITION_REASON="Phase ${PHASE} plans already exist"
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
PROPOSED_NEXT_STATE="WAITING_FOR_PHASE_PLANS"
TRANSITION_REASON="Architect spawned for Phase ${PHASE} planning"
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
            PROPOSED_NEXT_STATE="ERROR_RECOVERY"
            TRANSITION_REASON="Missing phase metadata: $error_details"
            ;;
        "invalid_phase")
            echo "ERROR: Invalid phase number or phase out of sequence"
            echo "Details: $error_details"
            PROPOSED_NEXT_STATE="ERROR_RECOVERY"
            TRANSITION_REASON="Invalid phase: $error_details"
            ;;
        "corrupted_state")
            echo "ERROR: orchestrator-state-v3.json is corrupted or invalid"
            echo "Details: $error_details"
            PROPOSED_NEXT_STATE="ERROR_RECOVERY"
            TRANSITION_REASON="Corrupted state: $error_details"
            ;;
        *)
            echo "ERROR: Unknown error in START_PHASE_ITERATION"
            echo "Type: $error_type"
            echo "Details: $error_details"
            PROPOSED_NEXT_STATE="ERROR_RECOVERY"
            TRANSITION_REASON="Unknown error: $error_type - $error_details"
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

## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete START_PHASE_ITERATION:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

**CRITICAL**: During state work, set these variables (DO NOT call update_state):
```bash
PROPOSED_NEXT_STATE="[state determined by logic]"
TRANSITION_REASON="[reason for transition]"
```

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Verify Proposed Transition Set
```bash
# Verify PROPOSED_NEXT_STATE was set during state work
if [ -z "$PROPOSED_NEXT_STATE" ]; then
    echo "❌ CRITICAL: PROPOSED_NEXT_STATE not set!"
    echo "State work must set PROPOSED_NEXT_STATE and TRANSITION_REASON"
    exit 288
fi

echo "✅ Proposed transition: START_PHASE_ITERATION → $PROPOSED_NEXT_STATE"
echo "   Reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Update State File (R288 - SUPREME LAW)
```bash
# State Manager validates transition and updates state files (SF 3.0 Pattern)
echo "🔄 Spawning State Manager for SHUTDOWN_CONSULTATION..."

# Prepare work results summary
WORK_RESULTS=$(cat <<EOF
{
  "state_completed": "START_PHASE_ITERATION",
  "work_accomplished": [
    "[List accomplishments from state work]"
  ],
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "START_PHASE_ITERATION" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON" \
  --work-results "$WORK_RESULTS"

# State Manager will:
# 1. Validate PROPOSED_NEXT_STATE exists and transition is valid
# 2. Update all 4 state files atomically (R288)
# 3. Commit and push state files
# 4. Return REQUIRED_NEXT_STATE (usually same as proposed unless invalid)

echo "✅ State Manager consultation complete"
echo "✅ State files updated by State Manager"
```

---

### ✅ Step 4: Save Proposed Transition to File (R288)
```bash
# Write proposed transition to file for State Manager
cat > "$CLAUDE_PROJECT_DIR/.proposed-transition" <<EOF
PROPOSED_NEXT_STATE="$PROPOSED_NEXT_STATE"
TRANSITION_REASON="$TRANSITION_REASON"
CURRENT_STATE="START_PHASE_ITERATION"
TIMESTAMP="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
EOF

echo "✅ Proposed transition saved to .proposed-transition"
```

---

### ✅ Step 5: Validate Current State File (R324)
```bash
# Validate state file integrity BEFORE proposing transition
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
echo "✅ Current state file validated"
```

---

### ✅ Step 6: Commit Work Products (if any)
```bash
# Commit any work products created during this state
# (State file itself is committed by State Manager)
if [ -n "$(git status --porcelain | grep -v orchestrator-state-v3.json)" ]; then
    git add .
    git commit -m "work: START_PHASE_ITERATION work products [R288]"
    git push
    echo "✅ Work products committed and pushed"
else
    echo "✅ No work products to commit"
fi
```

---

### ✅ Step 7: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "START_PHASE_ITERATION_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo
git commit -m "todo: orchestrator - START_PHASE_ITERATION complete [R287]"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 8: Output Transition Proposal (R288)
```bash
# Output the proposed transition for State Manager
echo "════════════════════════════════════════════════════════"
echo "PROPOSED STATE TRANSITION (R288):"
echo "  FROM: START_PHASE_ITERATION"
echo "  TO:   $PROPOSED_NEXT_STATE"
echo "  REASON: $TRANSITION_REASON"
echo "════════════════════════════════════════════════════════"
```

---

### ✅ Step 9: Output Continuation Flag (R405 - SUPREME LAW)
```bash
# Output continuation flag (R405)
# TRUE = state work complete, transition proposed
# FALSE = error occurred, manual intervention needed

if [ "$PROPOSED_NEXT_STATE" = "ERROR_RECOVERY" ]; then
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
else
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
fi
```

**⚠️ THIS MUST BE THE LAST OUTPUT BEFORE STOP! ⚠️**

---

### ✅ Step 10: State Manager Transition Checkpoint
```
🔄 STATE MANAGER TAKES CONTROL HERE (R288)

The State Manager (run-software-factory.sh) will:
1. Read .proposed-transition file
2. Validate proposed transition against state machine
3. Update orchestrator-state-v3.json with new state
4. Commit and push state file
5. Re-invoke orchestrator in new state (if CONTINUE-SOFTWARE-FACTORY=TRUE)

DO NOT PROCEED PAST THIS POINT - STATE MANAGER HANDLES TRANSITION
```

---

### ✅ Step 11: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for State Manager handoff (R322)
echo "🛑 State work complete - State Manager will handle transition"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 1: State work incomplete
- Missing Step 2: No proposed transition = stuck forever (R288 violation, -100%)
- Missing Step 3: No State Manager consultation = broken SF 3.0 pattern (R288 violation, -100%)
- Missing Step 4: State Manager can't read proposal = broken automation (R288 violation, -100%)
- Missing Step 5: Invalid state = corruption (R324 violation)
- Missing Step 6: Work products lost (data loss)
- Missing Step 7: TODOs lost on compaction (R287 violation, -20% to -100%)
- Missing Step 8: No transition visibility (R288 violation)
- Missing Step 9: Automation doesn't know if it can continue (R405 violation, -100%)
- Missing Step 10: State Manager handoff failed (R288 violation, -100%)
- Missing Step 11: Context corruption (R322 violation, -100%)

**ALL 11 STEPS ARE MANDATORY - NO EXCEPTIONS**

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
