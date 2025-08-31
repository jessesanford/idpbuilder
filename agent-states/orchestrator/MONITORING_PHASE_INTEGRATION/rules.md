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

### 2. 🚨🚨🚨 CHECK FOR PHASE INTEGRATION REPORT (CRITICAL) 🚨🚨🚨
```bash
# MANDATORY: Check for phase integration report
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
REPORT_FILE="efforts/phase${PHASE}/phase-integration/PHASE_INTEGRATION_REPORT.md"

if [ ! -f "$REPORT_FILE" ]; then
    echo "❌ CRITICAL: No phase integration report found at $REPORT_FILE"
    # Transition to ERROR_RECOVERY
    UPDATE_STATE="ERROR_RECOVERY"
    UPDATE_REASON="No phase integration report generated"
else
    echo "✅ Found phase integration report, analyzing..."
    
    # Extract status from report
    INTEGRATION_STATUS=$(grep "^Phase Integration Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    BUILD_STATUS=$(grep "^Phase Build Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    TEST_STATUS=$(grep "^Phase Test Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    CONFLICTS=$(grep "^Unresolved Conflicts:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    
    echo "Phase Integration Status: $INTEGRATION_STATUS"
    echo "Phase Build Status: $BUILD_STATUS"
    echo "Phase Test Status: $TEST_STATUS"
    echo "Unresolved Conflicts: $CONFLICTS"
    
    # Determine next state based on report
    if [[ "$INTEGRATION_STATUS" == "SUCCESS" ]] && \
       [[ "$BUILD_STATUS" == "PASSING" ]] && \
       [[ "$TEST_STATUS" == "PASSING" ]] && \
       [[ "$CONFLICTS" == "0" ]]; then
        echo "✅ Phase integration successful - proceeding to SPAWN_ARCHITECT_PHASE_ASSESSMENT"
        UPDATE_STATE="SPAWN_ARCHITECT_PHASE_ASSESSMENT"
        UPDATE_REASON="Phase integration complete and successful"
    elif [[ "$CONFLICTS" != "0" ]] || \
         [[ "$INTEGRATION_STATUS" == "FAILED" ]] || \
         [[ "$BUILD_STATUS" == "FAILED" ]]; then
        echo "🚨 Phase integration has issues - transitioning to PHASE_INTEGRATION_FEEDBACK_REVIEW"
        UPDATE_STATE="PHASE_INTEGRATION_FEEDBACK_REVIEW"
        UPDATE_REASON="Phase integration failed with conflicts or build issues"
    else
        echo "⚠️ Unexpected status combination - going to ERROR_RECOVERY"
        UPDATE_STATE="ERROR_RECOVERY"
        UPDATE_REASON="Unexpected phase integration status"
    fi
fi
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