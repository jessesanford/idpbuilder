# Orchestrator State: WAITING_FOR_DEMO_VALIDATION

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
## Purpose
Wait for Code Reviewer to complete demo validation, then enforce R291 Gate 4 + R331
based on results. This state is the FINAL ENFORCEMENT POINT for demo requirements.

**🚨 CRITICAL: Enforces BOTH R291 (demos exist and pass) AND R331 (demos are real, not simulated) 🚨**

**R322 Human Checkpoint**: Demos require manual human validation to verify they demonstrate real working functionality. Code Reviewer performs automated checks, but ultimate sign-off requires human approval that demos are authentic and demonstrate value.

## Entry Conditions
- Code Reviewer spawned for demo validation
- Demos should exist in demos/ directory
- Waiting for demo evaluation report

## State Responsibilities

### 1. Monitor for Demo Validation Completion

Check for the demo evaluation report that Code Reviewer creates:

```bash
# Extract demo validation metadata from state
INTEGRATE_WAVE_EFFORTS_TYPE=$(jq -r '.demo_validation.integration_type' orchestrator-state-v3.json)
DEMO_REPORT=$(jq -r '.demo_validation.evaluation_report' orchestrator-state-v3.json)
DEMO_LOG=$(jq -r '.demo_validation.evaluation_log' orchestrator-state-v3.json)

echo "📊 Monitoring for demo validation completion..."
echo "   Integration Type: $INTEGRATE_WAVE_EFFORTS_TYPE"
echo "   Expected Report: $DEMO_REPORT"
echo "   Expected Log: $DEMO_LOG"

# Check if report exists
if [ ! -f "$DEMO_REPORT" ]; then
    echo "⏳ Demo validation report not yet available"
    echo "   Code Reviewer still working..."
    echo "   Will check again when /continue-orchestrating is invoked"
    exit 0  # Not an error, just not ready yet
fi

echo "✅ Demo validation report found!"
```

### 2. Read Demo Validation Results

Extract the validation status and demo results:

```bash
# Read demo evaluation report
echo "📖 Reading demo validation results from: $DEMO_REPORT"

# Extract key fields from report
DEMO_STATUS=$(grep "^Demo Validation Status:" "$DEMO_REPORT" | cut -d: -f2 | tr -d ' ')
DEMOS_PASSED=$(grep "^- Demos Passed:" "$DEMO_REPORT" | cut -d: -f2 | tr -d ' ')
DEMOS_FAILED=$(grep "^- Demos Failed:" "$DEMO_REPORT" | cut -d: -f2 | tr -d ' ')

# R331: Check for R331 validation status
R331_STATUS=$(grep "^## R331 Validation Status:" "$DEMO_REPORT" | cut -d: -f2 | tr -d ' ')
if [ -n "$R331_STATUS" ]; then
    echo "🔍 R331 Validation Status: $R331_STATUS"
fi

echo "📊 Demo Validation Results:"
echo "   Status: $DEMO_STATUS"
echo "   Passed: $DEMOS_PASSED"
echo "   Failed: $DEMOS_FAILED"
echo "   R331 Compliance: ${R331_STATUS:-UNKNOWN}"

# Read demo execution log for detailed analysis
if [ -f "$DEMO_LOG" ]; then
    echo "📋 Demo execution log available at: $DEMO_LOG"
fi
```

### 3. Enforce R291 Gate 4 + R331 Compliance

This is the CRITICAL ENFORCEMENT POINT. Demo failures OR R331 violations MUST trigger ERROR_RECOVERY:

```bash
# R291 + R331 GATE 4 ENFORCEMENT - THIS IS ABSOLUTE
echo "🔴🔴🔴 R291 + R331 GATE 4 ENFORCEMENT 🔴🔴🔴"

# R331: Check for validation failures FIRST (higher priority)
if [ "$R331_STATUS" = "FAILED" ]; then
    echo "🔴 R331 VIOLATION: Demos are SIMULATED or INCOMPLETE!"
    echo "   PENALTY: -100% (simulation) or -50% to -75% (other violations)"
    echo "   Demos do not meet R331 requirements"
    echo "   This is a BLOCKING failure per R331"
    echo ""
    echo "🔴 MANDATORY: Transitioning to ERROR_RECOVERY"

    # Update tracking fields (ALLOWED - orchestrator maintains this data)
    jq ".error_recovery.trigger = \"R331_DEMO_VALIDATION_FAILURE\"" -i orchestrator-state-v3.json
    jq ".error_recovery.reason = \"Demos violate R331 - simulated or incomplete\"" -i orchestrator-state-v3.json
    jq ".error_recovery.timestamp = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state-v3.json
    jq ".error_recovery.demo_report = \"$DEMO_REPORT\"" -i orchestrator-state-v3.json

    # Set proposed next state (State Manager will update state_machine fields)
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="R331 demo validation failure - demos simulated/incomplete"
    # State Manager consultation happens in Step 3 of completion checklist

    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
    exit 331
fi

# R291: Check for demo execution failures
if [ "$DEMO_STATUS" = "PASSED" ] && [ "$DEMOS_FAILED" = "0" ] && [ "$R331_STATUS" != "FAILED" ]; then
    echo "✅ R291 GATE 4: PASSED"
    echo "✅ R331 COMPLIANCE: PASSED"
    echo "   All demos executed successfully"
    echo "   Demos are real working implementations (not simulated)"
    echo "   Integration may proceed to completion"

    # Update tracking fields (ALLOWED - orchestrator maintains this data)
    jq ".demo_validation.validation_status = \"passed\"" -i orchestrator-state-v3.json
    jq ".demo_validation.demos_passed = $DEMOS_PASSED" -i orchestrator-state-v3.json
    jq ".demo_validation.demos_failed = 0" -i orchestrator-state-v3.json
    jq ".demo_validation.validated_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state-v3.json

    # Determine next state based on integration type
    case "$INTEGRATE_WAVE_EFFORTS_TYPE" in
        wave)
            PROPOSED_NEXT_STATE="COMPLETE_WAVE"
            ;;
        phase)
            PROPOSED_NEXT_STATE="COMPLETE_WAVE"
            ;;
        project)
            PROPOSED_NEXT_STATE="COMPLETE_WAVE"
            ;;
        *)
            echo "🔴 Unknown integration type: $INTEGRATE_WAVE_EFFORTS_TYPE"
            exit 291
            ;;
    esac

    echo "✅ Transitioning to: $PROPOSED_NEXT_STATE"
    TRANSITION_REASON="Demo validation passed (R291 Gate 4 + R331)"
    # State Manager consultation happens in Step 3 of completion checklist

elif [ "$DEMO_STATUS" = "FAILED" ]; then
    echo "🔴 R291 GATE 4: FAILED"
    echo "   Demos did not execute successfully"
    echo "   MANDATORY: Must enter ERROR_RECOVERY per R291 lines 29-46"
    echo ""
    echo "⚠️⚠️⚠️ CRITICAL R291 VIOLATION ⚠️⚠️⚠️"
    echo "Per R291 line 44: Marking integration complete without"
    echo "passing build/test/demo = IMMEDIATE DISQUALIFICATION"
    echo ""
    echo "Demos MUST work before integration can proceed!"

    # Update tracking fields (ALLOWED - orchestrator maintains this data)
    jq ".demo_validation.validation_status = \"failed\"" -i orchestrator-state-v3.json
    jq ".demo_validation.demos_passed = $DEMOS_PASSED" -i orchestrator-state-v3.json
    jq ".demo_validation.demos_failed = $DEMOS_FAILED" -i orchestrator-state-v3.json
    jq ".demo_validation.validated_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state-v3.json

    jq ".error_recovery = {
      \"trigger\": \"R291_DEMO_GATE_FAILURE\",
      \"reason\": \"Integration demos failed - ${DEMOS_FAILED} demo(s) did not execute successfully\",
      \"integration_type\": \"$INTEGRATE_WAVE_EFFORTS_TYPE\",
      \"failed_demos\": $DEMOS_FAILED,
      \"demo_report\": \"$DEMO_REPORT\",
      \"demo_log\": \"$DEMO_LOG\",
      \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
      \"required_action\": \"Review demo failures, fix issues in effort branches, re-integrate per R291\"
    }" -i orchestrator-state-v3.json

    echo "🔴 Transitioning to ERROR_RECOVERY"
    echo "   Review $DEMO_REPORT for failure details"
    echo "   Review $DEMO_LOG for execution logs"

    # Set proposed next state (State Manager will update state_machine fields)
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="R291 demo gate failure - ${DEMOS_FAILED} demo(s) failed"
    # State Manager consultation happens in Step 3 of completion checklist

    # R291 violation detected - demos MUST work
    exit 291

else
    echo "⚠️ Demo validation status unclear or missing: '$DEMO_STATUS'"
    echo "   Expected: 'PASSED' or 'FAILED'"
    echo "   Treating as failure (safe default per R291)"

    # Update tracking fields (ALLOWED - orchestrator maintains this data)
    jq ".demo_validation.validation_status = \"unclear\"" -i orchestrator-state-v3.json
    jq ".demo_validation.validated_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state-v3.json

    jq ".error_recovery = {
      \"trigger\": \"R291_DEMO_STATUS_UNCLEAR\",
      \"reason\": \"Demo validation status unclear: '$DEMO_STATUS'\",
      \"integration_type\": \"$INTEGRATE_WAVE_EFFORTS_TYPE\",
      \"demo_report\": \"$DEMO_REPORT\",
      \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
      \"required_action\": \"Investigate demo validation report format\"
    }" -i orchestrator-state-v3.json

    # Set proposed next state (State Manager will update state_machine fields)
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Demo validation status unclear"
    # State Manager consultation happens in Step 3 of completion checklist

    exit 291
fi
```

### 4. Update and Commit State

Commit the state changes reflecting demo validation results:

```bash
# Commit state update
git add orchestrator-state-v3.json
if [ "$DEMO_STATUS" = "PASSED" ]; then
    git commit -m "state: demo validation passed - proceeding to $NEXT_STATE (R291 Gate 4)"
else
    git commit -m "error: demo validation failed - entering ERROR_RECOVERY (R291 Gate 4)"
fi
git push

echo "✅ State file updated and committed"
```

## Exit Conditions

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE
Conditions for continuation (whether passed or failed):

**Demos PASSED:**
- All demos executed successfully
- DEMOS_FAILED = 0
- State updated to appropriate completion state
- Continue to: WAVE_COMPLETE / COMPLETE_PHASE / PROJECT_INTEGRATE_WAVE_EFFORTS_FINALIZATION

**Demos FAILED (still continues, but to ERROR_RECOVERY):**
- Demo execution failed
- DEMOS_FAILED > 0
- State updated to ERROR_RECOVERY
- Error recovery metadata populated
- Continue to: ERROR_RECOVERY

**Both cases allow continuation** because:
- Passed demos → normal completion flow
- Failed demos → recoverable via ERROR_RECOVERY (fix and retry)

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE
Conditions requiring manual intervention:
- Demo report missing/corrupt (cannot determine status)
- State file corruption detected
- Integration type unknown/invalid
- State machine corruption (no valid next state)

**Manual intervention needed** because these indicate systemic issues,
not just demo failures.

## Next States

Based on demo validation results:

**If demos PASSED:**
- **WAVE_COMPLETE** (if wave integration)
- **COMPLETE_PHASE** (if phase integration)
- **PROJECT_INTEGRATE_WAVE_EFFORTS_FINALIZATION** (if project integration)

**If demos FAILED:**
- **ERROR_RECOVERY** (MANDATORY per R291 lines 29-46)

**Illegal Next States:**
- ANY completion state if demos failed
- Skipping ERROR_RECOVERY when demos failed = -100% penalty

## R291 Enforcement

This state IS the R291 Gate 4 enforcement mechanism.

**BLOCKING STATE** - Cannot proceed without demo validation

**Per R291 line 44:**
> "Marking integration complete without passing build/test/demo = IMMEDIATE DISQUALIFICATION"

**Enforcement mechanism:**
1. Read demo evaluation results
2. If PASSED and DEMOS_FAILED=0 → Allow completion
3. If FAILED or DEMOS_FAILED>0 → MUST go to ERROR_RECOVERY
4. If unclear → Treat as failed (safe default)

**Attempting to bypass this check = -100% penalty**

**Per R291 lines 29-46:** Demo failures MUST trigger ERROR_RECOVERY.
This is not optional. This is not negotiable. This is MANDATORY.

## State Machine Reference

```
SPAWN_CODE_REVIEWER_DEMO_VALIDATION
  ↓ (Code Reviewer executes demos)
WAITING_FOR_DEMO_VALIDATION  ← YOU ARE HERE
  ↓ (enforce R291 Gate 4)
  ├─ Demos passed (0 failures) → WAVE/PHASE/PROJECT_COMPLETE
  └─ Demos failed (>0 failures) → ERROR_RECOVERY (MANDATORY!)
```

## Critical Rules Referenced
- **R233**: Active Monitoring Pattern (BLOCKING - NOT passive waiting)
- **R291**: Integration Demo Requirement (BLOCKING - this state enforces it)
- **R322**: Mandatory Orchestrator Checkpoints (SUPREME LAW)
- **R234**: Mandatory State Traversal (cannot skip this state)
- **R300**: Comprehensive Fix Management (for handling failures)

## Remember

**This state is the last line of defense against incomplete integrations.**

Demos are not suggestions. They are REQUIREMENTS.

If demos fail:
- ✅ DO go to ERROR_RECOVERY
- ✅ DO fix issues in effort branches
- ✅ DO re-integrate
- ✅ DO re-validate demos

If demos fail:
- ❌ DO NOT proceed to completion
- ❌ DO NOT skip ERROR_RECOVERY
- ❌ DO NOT mark integration as complete
- ❌ DO NOT assume "it's probably fine"

**Build + Tests + Demo = ALL THREE MUST PASS**

No exceptions.

CONTINUE-SOFTWARE-FACTORY=TRUE

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="WAITING_FOR_DEMO_VALIDATION complete - [accomplishment description]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for SHUTDOWN_CONSULTATION
```bash
# State Manager validates transition and updates state files (SF 3.0 Pattern)
echo "🔄 Spawning State Manager for SHUTDOWN_CONSULTATION..."

# Prepare work results summary
WORK_RESULTS=$(cat <<EOF
{
  "state_completed": "WAITING_FOR_DEMO_VALIDATION",
  "work_accomplished": [
    "Read demo validation results",
    "Enforced R291 Gate 4 + R331 compliance",
    "Validated all demos passed or transitioned to ERROR_RECOVERY"
  ],
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "WAITING_FOR_DEMO_VALIDATION" \
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

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "WAITING_FOR_DEMO_VALIDATION_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - WAITING_FOR_DEMO_VALIDATION complete [R287]"; then
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

### ✅ Step 5: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 6: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No proposed next state = State Manager can't proceed
- Missing Step 3: No State Manager consultation = bypassing bookend pattern (-100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
