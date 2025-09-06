# Orchestrator - MONITORING_PHASE_INTEGRATION State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.yaml with new state
3. ✅ Committing and pushing the state file  
4. ✅ Providing work summary

### YOU MUST NOT:
- ❌ Continue to the next state automatically
- ❌ Start work for the new state
- ❌ Spawn agents for the new state
- ❌ Assume permission to continue

### STOP PROTOCOL:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: CURRENT_STATE → NEXT_STATE

### ✅ Current State Work Completed:
- [List completed work]

### 📊 Current Status:
- Current State: CURRENT_STATE
- Next State: NEXT_STATE
- TODOs Completed: X/Y
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to NEXT_STATE. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---

# ORCHESTRATOR STATE: MONITORING_PHASE_INTEGRATION

## 🔴🔴🔴 SUPREME DIRECTIVE: PHASE INTEGRATION FEEDBACK ENFORCEMENT 🔴🔴🔴

**YOU MUST CHECK FOR PHASE INTEGRATION REPORTS AND ACT ON FAILURES!**

## State Overview

In MONITORING_PHASE_INTEGRATION, you are monitoring the Integration Agent's phase-level integration progress and MUST check for phase integration reports.

## Required Actions

### 1. Monitor Phase Integration Agent Progress
```bash
# Check if Integration Agent is still running
INTEGRATION_PID=$(pgrep -f "integration-agent.*phase" || echo "")
if [ -n "$INTEGRATION_PID" ]; then
    echo "Phase Integration Agent still running (PID: $INTEGRATION_PID)"
    # Stay in MONITORING_PHASE_INTEGRATION
    sleep 5
    continue
fi
```

### 2. 🔴🔴🔴 CHECK PHASE INTEGRATION REPORT AND ENFORCE R291 GATES 🔴🔴🔴
```bash
# MANDATORY: Check for phase integration report AND enforce build/test gates
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
REPORT_FILE="efforts/phase${PHASE}/phase-integration/PHASE_INTEGRATION_REPORT.md"

if [ ! -f "$REPORT_FILE" ]; then
    echo "🔴 CRITICAL: No phase integration report found at $REPORT_FILE"
    # NO REPORT = IMMEDIATE ERROR_RECOVERY
    UPDATE_STATE="ERROR_RECOVERY"
    UPDATE_REASON="No phase integration report - R291 gates cannot be verified"
else
    echo "✅ Found phase integration report, enforcing R291 gates..."
    
    # Extract status from report
    INTEGRATION_STATUS=$(grep "^Phase Integration Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    BUILD_STATUS=$(grep "^Phase Build Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    TEST_STATUS=$(grep "^Phase Test Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    DEMO_STATUS=$(grep "^Phase Demo Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ' || echo "NOT_RUN")
    CONFLICTS=$(grep "^Unresolved Conflicts:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    
    echo "🔍 Phase Gate Status Check (R291 Enforcement):"
    echo "  Phase Build: $BUILD_STATUS"
    echo "  Phase Tests: $TEST_STATUS"
    echo "  Phase Demo: $DEMO_STATUS"
    echo "  Phase Integration: $INTEGRATION_STATUS"
    echo "  Conflicts: $CONFLICTS"
    
    # 🔴🔴🔴 R291 SUPREME GATE ENFORCEMENT FOR PHASE 🔴🔴🔴
    
    # BUILD GATE CHECK - PHASE LEVEL
    if [[ "$BUILD_STATUS" != "PASSING" ]] && [[ "$BUILD_STATUS" != "SUCCESS" ]]; then
        echo "🔴🔴🔴 PHASE BUILD GATE FAILED - MANDATORY ERROR_RECOVERY 🔴🔴🔴"
        echo "R291 VIOLATION: Phase build did not pass ($BUILD_STATUS)"
        echo "Cannot proceed to phase assessment without successful build!"
        UPDATE_STATE="ERROR_RECOVERY"
        UPDATE_REASON="R291 PHASE BUILD GATE FAILURE: $BUILD_STATUS"
        
    # TEST GATE CHECK - PHASE LEVEL
    elif [[ "$TEST_STATUS" != "PASSING" ]] && [[ "$TEST_STATUS" != "SUCCESS" ]]; then
        echo "🔴🔴🔴 PHASE TEST GATE FAILED - MANDATORY ERROR_RECOVERY 🔴🔴🔴"
        echo "R291 VIOLATION: Phase tests did not pass ($TEST_STATUS)"
        echo "Cannot proceed to phase assessment without passing tests!"
        UPDATE_STATE="ERROR_RECOVERY"
        UPDATE_REASON="R291 PHASE TEST GATE FAILURE: $TEST_STATUS"
        
    # DEMO GATE CHECK - PHASE LEVEL (if present)
    elif [[ "$DEMO_STATUS" != "NOT_RUN" ]] && \
         [[ "$DEMO_STATUS" != "PASSING" ]] && \
         [[ "$DEMO_STATUS" != "SUCCESS" ]]; then
        echo "🔴🔴🔴 PHASE DEMO GATE FAILED - MANDATORY ERROR_RECOVERY 🔴🔴🔴"
        echo "R291 VIOLATION: Phase demo did not pass ($DEMO_STATUS)"
        UPDATE_STATE="ERROR_RECOVERY"
        UPDATE_REASON="R291 PHASE DEMO GATE FAILURE: $DEMO_STATUS"
        
    # CONFLICTS CHECK
    elif [[ "$CONFLICTS" != "0" ]] && [[ -n "$CONFLICTS" ]]; then
        echo "🔴 Unresolved conflicts detected - cannot proceed"
        UPDATE_STATE="PHASE_INTEGRATION_FEEDBACK_REVIEW"
        UPDATE_REASON="Phase has $CONFLICTS unresolved conflicts"
        
    # INTEGRATION STATUS CHECK
    elif [[ "$INTEGRATION_STATUS" != "SUCCESS" ]]; then
        echo "🔴 Phase integration failed - review needed"
        UPDATE_STATE="PHASE_INTEGRATION_FEEDBACK_REVIEW"
        UPDATE_REASON="Phase integration status: $INTEGRATION_STATUS"
        
    # ALL GATES PASSED - PHASE CAN BE ASSESSED
    else
        echo "✅✅✅ ALL PHASE R291 GATES PASSED ✅✅✅"
        echo "  ✅ Phase Build: $BUILD_STATUS"
        echo "  ✅ Phase Tests: $TEST_STATUS"
        echo "  ✅ Phase Demo: $DEMO_STATUS (if run)"
        echo "  ✅ Phase Integration: $INTEGRATION_STATUS"
        echo "  ✅ Conflicts: $CONFLICTS"
        echo "Phase is ready for architect assessment!"
        UPDATE_STATE="SPAWN_ARCHITECT_PHASE_ASSESSMENT"
        UPDATE_REASON="All phase gates passed - ready for assessment"
    fi
fi

echo ""
echo "🎯 DECISION: Transitioning to $UPDATE_STATE"
echo "📝 REASON: $UPDATE_REASON"
```

### 3. Update State File with Phase Integration Results
```bash
# Update orchestrator state
yq eval ".current_state = \"$UPDATE_STATE\"" -i orchestrator-state.yaml
yq eval ".phase_integration_status.phase${PHASE} = \"$INTEGRATION_STATUS\"" -i orchestrator-state.yaml
yq eval ".state_transition_history += [{\"from\": \"MONITORING_PHASE_INTEGRATION\", \"to\": \"$UPDATE_STATE\", \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\", \"reason\": \"$UPDATE_REASON\"}]" -i orchestrator-state.yaml

# Commit state change
git add orchestrator-state.yaml
git commit -m "state: MONITORING_PHASE_INTEGRATION → $UPDATE_STATE - $UPDATE_REASON"
git push
```

## Valid Transitions

Based on phase integration report analysis:

1. **SUCCESS Path**: `MONITORING_PHASE_INTEGRATION` → `SPAWN_ARCHITECT_PHASE_ASSESSMENT`
   - When: Phase integration, build, tests all pass, no conflicts
   
2. **FAILURE Path**: `MONITORING_PHASE_INTEGRATION` → `PHASE_INTEGRATION_FEEDBACK_REVIEW`
   - When: Phase integration failed, conflicts exist, or build/tests fail
   
3. **ERROR Path**: `MONITORING_PHASE_INTEGRATION` → `ERROR_RECOVERY`
   - When: No report found or unexpected status

## Grading Criteria

- ✅ **+20%**: Properly check for phase integration report
- ✅ **+20%**: Correctly parse phase-level status fields
- ✅ **+20%**: Check for unresolved conflicts
- ✅ **+20%**: Transition to PHASE_INTEGRATION_FEEDBACK_REVIEW on failures
- ✅ **+20%**: Update phase integration status in state file

## Common Violations

- ❌ **-100%**: Ignoring phase integration failures
- ❌ **-50%**: Not checking for phase integration report
- ❌ **-50%**: Proceeding to phase assessment with conflicts
- ❌ **-30%**: Not updating phase integration status

## Related Rules

- R282: Phase Integration Protocol
- R259: Mandatory Phase Integration After Fixes
- R260: Integration Agent Core Requirements
- R263: Integration Documentation Requirements
- R206: State Machine Transition Validation

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

### R304: Mandatory Line Counter Tool Enforcement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
**Criticality**: BLOCKING - Manual counting = AUTOMATIC -100% FAILURE

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
