# State Manager Agent - SHUTDOWN_CONSULTATION State Rules

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
**State**: SHUTDOWN_CONSULTATION
**Agent**: state-manager
**Purpose**: Validate all state files and provide validation report at end of work session

---

## State Overview

This state is activated when:
- Orchestrator completes work session and prepares to exit
- State Manager invoked via bookend pattern (shutdown phase)
- All state files need final validation before session end

**Entry Conditions:**
- Orchestrator has completed work (implementation/integration/fixes)
- State files may have been updated during session
- Ready to validate consistency and commit state

**Exit Conditions:**
- All 4 state files validated successfully
- validation_result report generated
- Orchestrator receives approval to finalize session

---

## Core Responsibilities

### 1. Validate All State Files

**MANDATORY validation sequence:**

```bash
# Validate all 4 state files using schema validator
for file in orchestrator-state-v3.json bug-tracking.json integration-containers.json fix-cascade-state.json; do
    if [ -f "$file" ]; then
        bash tools/validate-state-file.sh "$file" || VALIDATION_FAILED=true
    fi
done
```

### 2. Check State Machine Consistency

**Verify state transitions are valid:**

```bash
# Extract current_state from orchestrator-state-v3.json
CURRENT_STATE=$(jq -r '.state_machine.current_state' orchestrator-state-v3.json)

# Verify state exists in state machine
STATE_EXISTS=$(jq --arg state "$CURRENT_STATE" '.states | has($state)' state-machines/software-factory-3.0-state-machine.json)

if [ "$STATE_EXISTS" != "true" ]; then
    VALIDATION_FAILED=true
fi
```

### 2.5. Validate Transition is Allowed (CRITICAL FIX)

**🚨🚨🚨 MANDATORY: Verify transition exists in allowed_transitions**

This validation was MISSING and caused Test 1 to allow invalid transitions!

```bash
validate_transition_allowed() {
    local from_state="$1"
    local to_state="$2"
    local state_machine="state-machines/software-factory-3.0-state-machine.json"

    echo "🔍 Validating transition: $from_state → $to_state"

    # Check if transition is in allowed_transitions list
    ALLOWED=$(jq --arg from "$from_state" --arg to "$to_state" \
      '.states[$from].allowed_transitions | contains([$to])' \
      "$state_machine" 2>/dev/null)

    if [ "$ALLOWED" != "true" ]; then
        echo "❌ INVALID TRANSITION: $from_state → $to_state"
        echo "   This transition is NOT in allowed_transitions list"

        # Show what IS allowed
        ALLOWED_LIST=$(jq -r --arg from "$from_state" \
          '.states[$from].allowed_transitions | join(", ")' \
          "$state_machine" 2>/dev/null || echo "UNKNOWN")
        echo "   Allowed from $from_state: $ALLOWED_LIST"

        # Log critical error
        echo "$(date -Iseconds): CRITICAL - Invalid transition: $from_state → $to_state (allowed: $ALLOWED_LIST)" >> state-manager-errors.log

        return 1
    fi

    echo "✅ Transition is allowed by state machine"
    return 0
}

# REQUIRED: Validate transition before accepting proposal
PREVIOUS_STATE=$(jq -r '.state_machine.previous_state // "INIT"' orchestrator-state-v3.json)
PROPOSED_NEXT_STATE=$(jq -r '.state_machine.current_state' orchestrator-state-v3.json)

if ! validate_transition_allowed "$PREVIOUS_STATE" "$PROPOSED_NEXT_STATE"; then
    echo "⚠️  WARNING: State machine violation detected"
    echo "   Proposed transition: $PREVIOUS_STATE → $PROPOSED_NEXT_STATE"
    echo "   This transition will be recorded but flagged as invalid"
    TRANSITION_INVALID=true
else
    TRANSITION_INVALID=false
fi
```

**Why this is critical:**
- Test 1 allowed SPAWN_ARCHITECT_MASTER_PLANNING → SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING
- This transition is NOT in allowed_transitions (should visit WAITING_FOR_MASTER_ARCHITECTURE first)
- Without this validation, state machine rules can be violated silently

### 2.75. Record State Transition in state_history (CRITICAL FIX)

**🚨🚨🚨 MANDATORY: Append transition to state_history array**

This recording was COMPLETELY MISSING and caused Test 1 state history gap!

```bash
record_state_transition() {
    local orchestrator_state_file="orchestrator-state-v3.json"
    local previous_state="$1"
    local current_state="$2"
    local orchestrator_proposal="${3:-$current_state}"
    local proposal_accepted="${4:-true}"
    local consultation_id="${5:-shutdown-$(date +%s)}"
    local transition_invalid="${6:-false}"

    echo "📝 Recording state transition: $previous_state → $current_state"

    # Create new state_history entry with all metadata
    NEW_HISTORY_ENTRY=$(jq -n \
      --arg from "$previous_state" \
      --arg to "$current_state" \
      --arg ts "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
      --arg validated "state-manager" \
      --arg cid "$consultation_id" \
      --arg proposal "$orchestrator_proposal" \
      --argjson accepted "$proposal_accepted" \
      --argjson invalid "$transition_invalid" \
      '{
        "from_state": $from,
        "to_state": $to,
        "timestamp": $ts,
        "validated_by": $validated,
        "consultation_id": $cid,
        "orchestrator_proposal": $proposal,
        "proposal_accepted": $accepted,
        "transition_invalid": $invalid
      }')

    # Append to state_history array
    TMP_FILE=$(mktemp)
    if jq ".state_machine.state_history += [$NEW_HISTORY_ENTRY]" \
      "$orchestrator_state_file" > "$TMP_FILE" 2>/dev/null; then
        mv "$TMP_FILE" "$orchestrator_state_file"
        echo "✅ State transition recorded in state_history"
        echo "   Entry: $previous_state → $current_state @ $(date -u +%Y-%m-%dT%H:%M:%SZ)"
        return 0
    else
        echo "❌ ERROR: Failed to record state transition"
        rm -f "$TMP_FILE"
        return 1
    fi
}

# REQUIRED: Call this function during shutdown consultation
PREVIOUS_STATE=$(jq -r '.state_machine.previous_state // "INIT"' orchestrator-state-v3.json)
CURRENT_STATE=$(jq -r '.state_machine.current_state' orchestrator-state-v3.json)
ORCHESTRATOR_PROPOSAL="${ORCHESTRATOR_PROPOSAL:-$CURRENT_STATE}"
PROPOSAL_ACCEPTED="${PROPOSAL_ACCEPTED:-true}"
CONSULTATION_ID="${CONSULTATION_ID:-shutdown-$(date +%s)}"

if ! record_state_transition "$PREVIOUS_STATE" "$CURRENT_STATE" \
  "$ORCHESTRATOR_PROPOSAL" "$PROPOSAL_ACCEPTED" "$CONSULTATION_ID" "$TRANSITION_INVALID"; then
    echo "⚠️  WARNING: State transition recording failed"
    echo "   This will cause state_history gaps (Test 1 failure mode)"
    VALIDATION_FAILED=true
fi
```

**Why this is critical:**
- Test 1 had **6 missing state transitions** because recording stopped
- **1 hour 57 minute gap** in state_history (23:26:46Z to 01:24:19Z)
- Auto-stop mechanism relies on state_history to detect target state
- Without recording, tests cannot validate state coverage

**What this adds:**
- Appends new entry to `.state_machine.state_history[]` array
- Includes timestamp, validation metadata, proposal tracking
- Records whether transition was valid per state machine
- Enables Test 1 auto-stop to detect WAVE_START

### 3. Generate validation_result Report

**REQUIRED format (structured markdown):**

```markdown
# State Manager Shutdown Consultation Report

**Date**: [TIMESTAMP_UTC]
**Session**: [SESSION_ID or SESSION_COUNT]
**Agent**: state-manager

---

## Validation Results

### orchestrator-state-v3.json
- **Schema Validation**: [PASS/FAIL]
- **State Machine Consistency**: [PASS/FAIL]
- **Current State**: [state_name]
- **Errors**: [List or "None"]

### bug-tracking.json
- **Schema Validation**: [PASS/FAIL]
- **Open Bugs**: [count]
- **Errors**: [List or "None"]

### integration-containers.json
- **Schema Validation**: [PASS/FAIL]
- **Active Containers**: [count]
- **Errors**: [List or "None"]

### fix-cascade-state.json
- **Schema Validation**: [PASS/FAIL]
- **Active Cascades**: [count or "N/A if file doesn't exist"]
- **Errors**: [List or "None"]

---

## Consistency Checks

### Cross-File References
- **Bug IDs in orchestrator-state**: [CONSISTENT/INCONSISTENT]
- **Container IDs in orchestrator-state**: [CONSISTENT/INCONSISTENT]
- **Cascade IDs in orchestrator-state**: [CONSISTENT/INCONSISTENT]

### State Integrity
- **Orphaned references**: [count or "None"]
- **Duplicate IDs**: [count or "None"]
- **Missing required fields**: [count or "None"]

---

## Session Summary

### Work Completed This Session
- **States Transitioned**: [list state transitions]
- **Efforts Modified**: [count or "None"]
- **Bugs Created**: [count or "None"]
- **Containers Updated**: [count or "None"]

### State File Changes
- **Files Modified**: [list]
- **Commits Made**: [count]
- **Last Commit**: [commit hash + message]

---

## Validation Directive

### Status: [APPROVED / NEEDS_FIXES]

**APPROVED** - All validations passed, safe to finalize session:
- ✅ All state files schema-valid
- ✅ State machine consistent
- ✅ No cross-file reference errors
- ✅ No orphaned data

**NEEDS_FIXES** - Validation failures detected:
- ❌ [List each validation failure]
- 🔧 [Recommended fix for each]
- ⏸️  Do NOT finalize session until fixed

### Required Actions
[If APPROVED:]
1. Commit final state updates (if any pending)
2. Push all commits to remote
3. Set CONTINUE-SOFTWARE-FACTORY flag appropriately
4. Exit cleanly

[If NEEDS_FIXES:]
1. STOP session finalization
2. Fix validation errors listed above
3. Re-run shutdown consultation
4. Only proceed when APPROVED

### Next State Recommendation
- **If APPROVED**: [Next state from state machine]
- **If NEEDS_FIXES**: ERROR_RECOVERY or stay in current state

---

## Consultation Complete

**Report Generated**: [TIMESTAMP_UTC]
**Validation Status**: [APPROVED/NEEDS_FIXES]
**Safe to Finalize**: [YES/NO]
```

---

## Validation Logic Details

### Schema Validation

**Use existing validator:**

```bash
# Validate each file against its schema
bash tools/validate-state-file.sh orchestrator-state-v3.json
bash tools/validate-state-file.sh bug-tracking.json
bash tools/validate-state-file.sh integration-containers.json

# fix-cascade-state.json is optional (only exists if cascade active)
if [ -f fix-cascade-state.json ]; then
    bash tools/validate-state-file.sh fix-cascade-state.json
fi
```

**Exit code 0 = PASS, non-zero = FAIL**

### State Machine Consistency

**Check current state validity:**

```bash
# 1. Extract current_state
CURRENT_STATE=$(jq -r '.state_machine.current_state' orchestrator-state-v3.json)

# 2. Verify state exists in state machine
jq --arg state "$CURRENT_STATE" \
   'if .states | has($state) then "CONSISTENT" else "INCONSISTENT - Unknown state" end' \
   state-machines/software-factory-3.0-state-machine.json
```

### Cross-File Reference Validation

**Verify IDs match across files:**

```bash
# Example: Check bug IDs
BUGS_IN_TRACKING=$(jq -r '.bugs[].bug_id' bug-tracking.json)
BUGS_IN_ORCHESTRATOR=$(jq -r '.project_progression.bugs[]?' orchestrator-state-v3.json)

# Compare lists for orphans
# (Simplified - actual implementation would be more thorough)
```

---

## Example Scenarios

### Scenario 1: Clean Shutdown (All Valid)

**Input:**
- All 4 files exist and schema-valid
- Current state is valid in state machine
- No cross-file reference errors

**Output (validation_result):**

```markdown
## Validation Results

### orchestrator-state-v3.json
- **Schema Validation**: PASS
- **State Machine Consistency**: PASS
- **Current State**: INTEGRATE_WAVE_EFFORTS
- **Errors**: None

### bug-tracking.json
- **Schema Validation**: PASS
- **Open Bugs**: 3
- **Errors**: None

### integration-containers.json
- **Schema Validation**: PASS
- **Active Containers**: 1
- **Errors**: None

### fix-cascade-state.json
- **Schema Validation**: N/A (file does not exist - no active cascade)
- **Active Cascades**: 0
- **Errors**: None

## Validation Directive

### Status: APPROVED

**APPROVED** - All validations passed, safe to finalize session:
- ✅ All state files schema-valid
- ✅ State machine consistent
- ✅ No cross-file reference errors
- ✅ No orphaned data

### Required Actions
1. Commit final state updates (if any pending)
2. Push all commits to remote
3. Set CONTINUE-SOFTWARE-FACTORY=TRUE
4. Exit cleanly

### Next State Recommendation
- **If APPROVED**: REVIEW_WAVE_INTEGRATION
```

**Orchestrator Action:** Proceed with session finalization

---

### Scenario 2: Validation Failure (Schema Error)

**Input:**
- orchestrator-state-v3.json has invalid JSON
- Other files valid

**Output (validation_result):**

```markdown
## Validation Results

### orchestrator-state-v3.json
- **Schema Validation**: FAIL
- **State Machine Consistency**: UNABLE_TO_CHECK (schema invalid)
- **Current State**: UNKNOWN
- **Errors**:
  - Line 42: Missing closing brace
  - Invalid JSON syntax

### Validation Directive

### Status: NEEDS_FIXES

**NEEDS_FIXES** - Validation failures detected:
- ❌ orchestrator-state-v3.json schema validation failed
- 🔧 Fix JSON syntax error at line 42
- ⏸️  Do NOT finalize session until fixed

### Required Actions
1. STOP session finalization
2. Fix JSON syntax in orchestrator-state-v3.json
3. Re-run shutdown consultation
4. Only proceed when APPROVED

### Safe to Finalize: NO
```

**Orchestrator Action:** Fix errors, re-consult State Manager

---

### Scenario 3: State Machine Inconsistency

**Input:**
- All files schema-valid
- current_state = "INVALID_STATE_XYZ" (doesn't exist in state machine)

**Output (validation_result):**

```markdown
## Validation Results

### orchestrator-state-v3.json
- **Schema Validation**: PASS
- **State Machine Consistency**: FAIL
- **Current State**: INVALID_STATE_XYZ
- **Errors**:
  - State "INVALID_STATE_XYZ" not found in state-machines/software-factory-3.0-state-machine.json
  - Possible states: [list valid states]

### Validation Directive

### Status: NEEDS_FIXES

**NEEDS_FIXES** - Validation failures detected:
- ❌ Invalid current_state in orchestrator-state-v3.json
- 🔧 Update current_state to a valid state from state machine
- ⏸️  Do NOT finalize session until fixed

### Required Actions
1. STOP session finalization
2. Correct state_machine.state_machine.current_state in orchestrator-state-v3.json
3. Re-run shutdown consultation
4. Only proceed when APPROVED

### Safe to Finalize: NO
```

**Orchestrator Action:** Fix state reference, re-consult

---

## Rules Compliance

**This state enforces:**
- **R506**: Pre-commit validation (state files validated before commit)
- **R288**: Atomic state updates (all 4 files consistent)
- **R516**: State naming convention (state exists in state machine)
- **R203**: State-aware operation (consultation follows bookend pattern)

**Validation ensures:**
- No invalid state file commits
- State machine integrity maintained
- Cross-file references consistent
- Safe session finalization

---

## Orchestrator Integration

**Bookend pattern (shutdown phase):**

```markdown
1. Orchestrator completes work session
2. Orchestrator spawns State Manager in SHUTDOWN_CONSULTATION state
3. State Manager validates all state files
4. State Manager generates validation_result report
5. State Manager returns report to Orchestrator
6. If APPROVED: Orchestrator finalizes session
7. If NEEDS_FIXES: Orchestrator fixes errors and re-consults
```

**State Manager provides GATE FUNCTION** - prevents bad state commits.

---

## Success Criteria

Shutdown consultation succeeds when:

✅ All 4 state files validated (schema + consistency)
✅ State machine current_state is valid
✅ Cross-file references consistent
✅ validation_result report generated
✅ Status = APPROVED (or NEEDS_FIXES with clear fix guidance)
✅ Orchestrator receives actionable directive

---

## Related Files

- **Startup Consultation**: `agent-states/state-manager/STARTUP_CONSULTATION/rules.md`
- **Validator Tool**: `tools/validate-state-file.sh`
- **State Machine**: `state-machines/software-factory-3.0-state-machine.json`
- **Schemas**: `schemas/*-schema.json`

---

**REMEMBER**: Shutdown consultation is the last line of defense before committing state. If validation fails, STOP and fix - never bypass.
