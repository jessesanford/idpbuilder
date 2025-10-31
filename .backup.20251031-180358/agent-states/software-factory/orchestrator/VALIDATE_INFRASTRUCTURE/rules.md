# ORCHESTRATOR STATE: VALIDATE_INFRASTRUCTURE


## 🚨🚨🚨 STATE PURPOSE [BLOCKING]
Validate ALL infrastructure configuration against authoritative sources BEFORE any implementation work begins. This is a MANDATORY gate that prevents catastrophic failures.

## SF 3.0 Validation Context

This validation state ensures infrastructure integrity in SF 3.0:
- Reads pre-planned infrastructure from orchestrator-state-v3.json `pre_planned_infrastructure` section
- Validates all effort directories, branches, and repository URLs match target configuration
- Updates `state_machine.current_state` with validation results before any agent spawning
- Records validation status in orchestrator-state-v3.json per R288 atomic update protocol
- Prevents catastrophic failures from infrastructure pointing to wrong repositories (R508 enforcement)

## 🔴🔴🔴 SUPREME LAW ENFORCEMENT
This state enforces:
- **R507**: Mandatory Infrastructure Validation [BLOCKING]
- **R508**: Target Repository Enforcement [SUPREME LAW]

ANY violation requires immediate transition to ERROR_RECOVERY.

## 🚨🚨🚨 ENTRY CONDITION VALIDATION 🚨🚨🚨

**CRITICAL**: Before performing validation, verify entry conditions are met.

### MANDATORY ENTRY CHECKS

Upon entering VALIDATE_INFRASTRUCTURE, orchestrator MUST verify:

```bash
echo "🔍 VALIDATING ENTRY CONDITIONS FOR VALIDATE_INFRASTRUCTURE..."

# Entry Check 1: All efforts claim creation complete
echo "📋 Checking all efforts are marked as created..."
CURRENT_PHASE=$(jq -r '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
CURRENT_WAVE=$(jq -r '.project_progression.current_wave.wave_number' orchestrator-state-v3.json)

UNCREATED_EFFORTS=$(jq -r '[.pre_planned_infrastructure.efforts | to_entries[] | select(.value.phase == "phase'$CURRENT_PHASE'" and .value.wave == "wave'$CURRENT_WAVE'" and .value.created == false)] | length' orchestrator-state-v3.json)

if [ "$UNCREATED_EFFORTS" -gt 0 ]; then
    # Check if this is expected during multi-effort creation
    # Allow controlled looping with a maximum iteration count (BUG #4 FIX)
    VALIDATE_INFRA_LOOPS=$(jq -r '.state_machine.loop_detection.validate_infrastructure_loops // 0' orchestrator-state-v3.json)

    if [ "$VALIDATE_INFRA_LOOPS" -lt 5 ]; then
        echo "⚠️ Some efforts not yet created ($UNCREATED_EFFORTS remaining)"
        echo "   Returning to CREATE_NEXT_INFRASTRUCTURE (iteration $((VALIDATE_INFRA_LOOPS + 1))/5)"

        # Increment loop counter
        jq ".state_machine.loop_detection.validate_infrastructure_loops = ($VALIDATE_INFRA_LOOPS + 1)" \
            orchestrator-state-v3.json > tmp-state.json && mv tmp-state.json orchestrator-state-v3.json

        PROPOSED_NEXT_STATE="CREATE_NEXT_INFRASTRUCTURE"
        TRANSITION_REASON="Some efforts not created yet - returning to create remaining (loop $((VALIDATE_INFRA_LOOPS + 1)))"
        exit 1
    else
        echo "❌ ENTRY VIOLATION: Still have $UNCREATED_EFFORTS uncreated efforts after 5 loops"
        echo "   This indicates a stuck state - transitioning to ERROR_RECOVERY"
        PROPOSED_NEXT_STATE="ERROR_RECOVERY"
        TRANSITION_REASON="Entry condition failed - efforts not created after maximum loop iterations"
        exit 1
    fi
fi

# If all efforts created, reset loop counter
jq 'del(.state_machine.loop_detection.validate_infrastructure_loops)' \
    orchestrator-state-v3.json > tmp-state.json && mv tmp-state.json orchestrator-state-v3.json
echo "✅ All efforts marked as created"

# Entry Check 2: No efforts should have validation_failure_reason set from creation
echo "📋 Checking no efforts have creation failures..."
FAILED_EFFORTS=$(jq -r '[.pre_planned_infrastructure.efforts | to_entries[] | select(.value.phase == "phase'$CURRENT_PHASE'" and .value.wave == "wave'$CURRENT_WAVE'" and .value.validation_failure_reason != null)] | length' orchestrator-state-v3.json)

if [ "$FAILED_EFFORTS" -gt 0 ]; then
    echo "❌ ENTRY VIOLATION: $FAILED_EFFORTS efforts have validation failures from creation"
    echo "   These should have been caught in CREATE_NEXT_INFRASTRUCTURE"
    echo "   Transitioning to ERROR_RECOVERY to investigate"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Entry condition failed - efforts have pre-existing validation failures"
    exit 1
fi

echo "✅ No pre-existing validation failures"

# Entry Check 3: Verify we're in correct mandatory sequence
echo "📋 Checking mandatory sequence position..."
# Should be in wave_execution sequence, after CREATE_NEXT_INFRASTRUCTURE loop
CURRENT_STATE=$(jq -r '.state_machine.current_state' orchestrator-state-v3.json)

if [ "$CURRENT_STATE" != "VALIDATE_INFRASTRUCTURE" ]; then
    echo "⚠️  WARNING: State file shows different current_state: $CURRENT_STATE"
fi

echo "✅ All entry conditions validated"
echo ""
```

**Purpose**:
- Catches issues BEFORE validation runs
- Provides clearer error messages
- Prevents ERROR_RECOVERY during validation for issues that should have been caught earlier
- Ensures VALIDATE_INFRASTRUCTURE only runs when infrastructure is actually ready

**Penalty for skipping entry checks**: -30% (causes confusing validation failures)

## 🚨🚨🚨 ENTRY CONDITIONS [BLOCKING]
MUST have:
1. ✅ Just completed CREATE_NEXT_INFRASTRUCTURE
2. ✅ Infrastructure directories created
3. ✅ Need to verify configuration before spawning agents

## 🚨🚨🚨 REQUIRED ACTIONS [MANDATORY]

### 1. RUN VALIDATION SCRIPT
```bash
echo "🔍 Starting infrastructure validation (R507/R508)..."
bash $CLAUDE_PROJECT_DIR/utilities/validate-infrastructure.sh

VALIDATION_RESULT=$?

if [ $VALIDATION_RESULT -eq 0 ]; then
    echo "✅ Infrastructure validation PASSED"
    # Can proceed to next state
elif [ $VALIDATION_RESULT -eq 911 ]; then
    echo "🔴🔴🔴 CATASTROPHIC: Wrong repository configured (R508 violation)"
    # MUST transition to ERROR_RECOVERY
else
    echo "❌ Infrastructure validation FAILED"
    # MUST transition to ERROR_RECOVERY
fi
```

### 2. CHECK EACH INFRASTRUCTURE COMPONENT
For EVERY effort/split/integration directory:

#### A. VALIDATE EFFORT INFRASTRUCTURE
```bash
# For each effort in pre_planned_infrastructure
for effort_key in $(echo "✅ State file updated to: $NEXT_STATE"
```

---

### ✅ Step 4: Validate State File (R324)
```bash
# Validate state file before committing
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
echo "✅ State file validated"
```

---

### ✅ Step 5: Commit State File (R288)
```bash
# Commit and push state file immediately
git add orchestrator-state-v3.json

if ! git commit -m "state: VALIDATE_INFRASTRUCTURE → $NEXT_STATE - VALIDATE_INFRASTRUCTURE complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: VALIDATE_INFRASTRUCTURE"
    echo "Attempted transition from: VALIDATE_INFRASTRUCTURE"
    echo ""
    echo "Common causes:"
    echo "  - Schema validation failure (check pre-commit hook output above)"
    echo "  - Missing required fields in JSON files"
    echo "  - Invalid JSON syntax"
    echo ""
    echo "🛑 Cannot proceed - manual intervention required"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=SCHEMA_VALIDATION"
    exit 1
fi

git push || echo "⚠️ WARNING: Push failed - committed locally"
echo "✅ State file committed and pushed"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "VALIDATE_INFRASTRUCTURE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - VALIDATE_INFRASTRUCTURE complete [R287]"; then
    echo "❌ ERROR: Failed to commit TODO files"
    echo "This is non-fatal but TODOs may be lost in compaction"
    echo "Proceeding with state execution..."
    # Don't exit - TODO commit failure is not fatal
fi

git push || echo "⚠️ WARNING: TODO push failed - committed locally"
echo "✅ TODOs saved and committed"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 7: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 8: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No state update = state machine broken (R288 violation, -100%)
- Missing Step 4: Invalid state = corruption (R324 violation)
- Missing Step 5: No commit = state lost on compaction (R288 violation, -100%)
- Missing Step 6: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 7: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 8: No exit = R322 violation (-100%)

**ALL 8 STEPS ARE MANDATORY - NO EXCEPTIONS**

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
