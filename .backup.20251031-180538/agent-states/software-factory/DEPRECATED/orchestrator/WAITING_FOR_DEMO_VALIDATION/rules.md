# Orchestrator State: WAITING_FOR_DEMO_VALIDATION

## Purpose
Wait for Code Reviewer to complete demo validation, then enforce R291 Gate 4
based on results. This state is the FINAL ENFORCEMENT POINT for demo requirements.

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

echo "📊 Demo Validation Results:"
echo "   Status: $DEMO_STATUS"
echo "   Passed: $DEMOS_PASSED"
echo "   Failed: $DEMOS_FAILED"

# Read demo execution log for detailed analysis
if [ -f "$DEMO_LOG" ]; then
    echo "📋 Demo execution log available at: $DEMO_LOG"
fi
```

### 3. Enforce R291 Gate 4

This is the CRITICAL ENFORCEMENT POINT. Demo failures MUST trigger ERROR_RECOVERY:

```bash
# R291 GATE 4 ENFORCEMENT - THIS IS ABSOLUTE
echo "🔴🔴🔴 R291 GATE 4 ENFORCEMENT 🔴🔴🔴"

if [ "$DEMO_STATUS" = "PASSED" ] && [ "$DEMOS_FAILED" = "0" ]; then
    echo "✅ R291 GATE 4: PASSED"
    echo "   All demos executed successfully"
    echo "   Integration may proceed to completion"

    # Update state with success
    jq ".demo_validation.validation_status = \"passed\"" -i orchestrator-state-v3.json
    jq ".demo_validation.demos_passed = $DEMOS_PASSED" -i orchestrator-state-v3.json
    jq ".demo_validation.demos_failed = 0" -i orchestrator-state-v3.json
    jq ".demo_validation.validated_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state-v3.json

    # Determine next state based on integration type
    case "$INTEGRATE_WAVE_EFFORTS_TYPE" in
        wave)
            NEXT_STATE="WAVE_COMPLETE"
            ;;
        phase)
            NEXT_STATE="COMPLETE_PHASE"
            ;;
        project)
            NEXT_STATE="PROJECT_INTEGRATE_WAVE_EFFORTS_FINALIZATION"
            ;;
        *)
            echo "🔴 Unknown integration type: $INTEGRATE_WAVE_EFFORTS_TYPE"
            exit 291
            ;;
    esac

    echo "✅ Transitioning to: $NEXT_STATE"
    jq ".state_machine.current_state = \"$NEXT_STATE\"" -i orchestrator-state-v3.json

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

    # Update state with failure details
    jq ".demo_validation.validation_status = \"failed\"" -i orchestrator-state-v3.json
    jq ".demo_validation.demos_passed = $DEMOS_PASSED" -i orchestrator-state-v3.json
    jq ".demo_validation.demos_failed = $DEMOS_FAILED" -i orchestrator-state-v3.json
    jq ".demo_validation.validated_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state-v3.json

    # MANDATORY transition to ERROR_RECOVERY
    jq ".state_machine.current_state = \"ERROR_RECOVERY\"" -i orchestrator-state-v3.json
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

    echo "🔴 State updated to ERROR_RECOVERY"
    echo "   Review $DEMO_REPORT for failure details"
    echo "   Review $DEMO_LOG for execution logs"

    NEXT_STATE="ERROR_RECOVERY"

    # R291 violation detected - demos MUST work
    exit 291

else
    echo "⚠️ Demo validation status unclear or missing: '$DEMO_STATUS'"
    echo "   Expected: 'PASSED' or 'FAILED'"
    echo "   Treating as failure (safe default per R291)"

    # Unclear status = treat as failure for safety
    jq ".demo_validation.validation_status = \"unclear\"" -i orchestrator-state-v3.json
    jq ".demo_validation.validated_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state-v3.json

    jq ".state_machine.current_state = \"ERROR_RECOVERY\"" -i orchestrator-state-v3.json
    jq ".error_recovery = {
      \"trigger\": \"R291_DEMO_STATUS_UNCLEAR\",
      \"reason\": \"Demo validation status unclear: '$DEMO_STATUS'\",
      \"integration_type\": \"$INTEGRATE_WAVE_EFFORTS_TYPE\",
      \"demo_report\": \"$DEMO_REPORT\",
      \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
      \"required_action\": \"Investigate demo validation report format\"
    }" -i orchestrator-state-v3.json

    NEXT_STATE="ERROR_RECOVERY"
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

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
