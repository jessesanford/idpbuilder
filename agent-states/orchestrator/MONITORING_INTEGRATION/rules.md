# ORCHESTRATOR STATE: MONITORING_INTEGRATION

## 🔴🔴🔴 SUPREME DIRECTIVE: INTEGRATION FEEDBACK ENFORCEMENT 🔴🔴🔴

**YOU MUST CHECK FOR INTEGRATION REPORTS AND ACT ON FAILURES!**

## State Overview

In MONITORING_INTEGRATION, you are monitoring the Integration Agent's progress and MUST check for integration reports to determine next state.

## Required Actions

### 1. Monitor Integration Agent Progress
```bash
# Check if Integration Agent is still running
INTEGRATION_PID=$(pgrep -f "integration-agent" || echo "")
if [ -n "$INTEGRATION_PID" ]; then
    echo "Integration Agent still running (PID: $INTEGRATION_PID)"
    # Stay in MONITORING_INTEGRATION
    sleep 5
    continue
fi
```

### 2. 🚨🚨🚨 CHECK FOR INTEGRATION REPORT AND DEMO (CRITICAL) 🚨🚨🚨

**Per R291 and R292: Integration is NOT complete until demo passes!**

```bash
# MANDATORY: Check for integration report
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
WAVE=$(yq '.current_wave' orchestrator-state.yaml)
REPORT_FILE="efforts/phase${PHASE}/wave${WAVE}/integration-workspace/INTEGRATION_REPORT.md"
DEMO_STATUS_FILE="efforts/phase${PHASE}/wave${WAVE}/integration-workspace/DEMO_STATUS.md"

if [ ! -f "$REPORT_FILE" ]; then
    echo "❌ CRITICAL: No integration report found at $REPORT_FILE"
    # Transition to ERROR_RECOVERY
    UPDATE_STATE="ERROR_RECOVERY"
    UPDATE_REASON="No integration report generated"
else
    echo "✅ Found integration report, analyzing..."
    
    # Extract status from report
    INTEGRATION_STATUS=$(grep "^Integration Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    BUILD_STATUS=$(grep "^Build Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    TEST_STATUS=$(grep "^Test Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    DEMO_STATUS=$(grep "^Demo Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ' || echo "NOT_RUN")
    
    echo "Integration Status: $INTEGRATION_STATUS"
    echo "Build Status: $BUILD_STATUS"
    echo "Test Status: $TEST_STATUS"
    echo "Demo Status: $DEMO_STATUS"
    
    # R291 ENFORCEMENT: Demo MUST pass
    if [[ "$DEMO_STATUS" != "PASSING" ]] && [[ "$DEMO_STATUS" != "SUCCESS" ]]; then
        echo "🚨 DEMO NOT PASSING - Integration BLOCKED per R291!"
        echo "Demo must build, run, and pass before integration is complete"
        UPDATE_STATE="INTEGRATION_FEEDBACK_REVIEW"
        UPDATE_REASON="Demo not passing - fixes required in effort branches (R292)"
    # All must pass for success
    elif [[ "$INTEGRATION_STATUS" == "SUCCESS" ]] && \
         [[ "$BUILD_STATUS" == "PASSING" ]] && \
         [[ "$TEST_STATUS" == "PASSING" ]] && \
         [[ "$DEMO_STATUS" == "PASSING" || "$DEMO_STATUS" == "SUCCESS" ]]; then
        echo "✅ Integration successful with passing demo - proceeding to WAVE_REVIEW"
        UPDATE_STATE="WAVE_REVIEW"
        UPDATE_REASON="Integration complete with successful demo"
    # Any failure triggers feedback review
    elif [[ "$BUILD_STATUS" == "BLOCKED_BY_DEPENDENCIES" ]] || \
         [[ "$INTEGRATION_STATUS" == "FAILED" ]] || \
         [[ "$BUILD_STATUS" == "FAILED" ]] || \
         [[ "$TEST_STATUS" == "FAILED" ]]; then
        echo "🚨 Integration has issues - transitioning to INTEGRATION_FEEDBACK_REVIEW"
        UPDATE_STATE="INTEGRATION_FEEDBACK_REVIEW"
        UPDATE_REASON="Integration failed - fixes required in effort branches (R292)"
    else
        echo "⚠️ Unexpected status combination - going to ERROR_RECOVERY"
        UPDATE_STATE="ERROR_RECOVERY"
        UPDATE_REASON="Unexpected integration status"
    fi
fi
```

### 3. Update State File
```bash
# Update orchestrator state
yq eval ".current_state = \"$UPDATE_STATE\"" -i orchestrator-state.yaml
yq eval ".state_transition_history += [{\"from\": \"MONITORING_INTEGRATION\", \"to\": \"$UPDATE_STATE\", \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\", \"reason\": \"$UPDATE_REASON\"}]" -i orchestrator-state.yaml

# Commit state change
git add orchestrator-state.yaml
git commit -m "state: MONITORING_INTEGRATION → $UPDATE_STATE - $UPDATE_REASON"
git push
```

## Valid Transitions

Based on integration report analysis:

1. **SUCCESS Path**: `MONITORING_INTEGRATION` → `WAVE_REVIEW`
   - When: Integration, build, and tests all pass
   
2. **FAILURE Path**: `MONITORING_INTEGRATION` → `INTEGRATION_FEEDBACK_REVIEW`
   - When: Integration failed, build blocked, or tests fail
   
3. **ERROR Path**: `MONITORING_INTEGRATION` → `ERROR_RECOVERY`
   - When: No report found or unexpected status

## Grading Criteria

- ✅ **+20%**: Properly check for integration report
- ✅ **+20%**: Correctly parse report status fields
- ✅ **+20%**: Transition to INTEGRATION_FEEDBACK_REVIEW on failures
- ✅ **+20%**: Never ignore integration failures
- ✅ **+20%**: Update state file with proper reason

## Common Violations

- ❌ **-100%**: Ignoring integration failures and marking COMPLETE
- ❌ **-50%**: Not checking for integration report
- ❌ **-50%**: Transitioning to WAVE_REVIEW when integration failed
- ❌ **-30%**: Not parsing report status fields

## Related Rules

- R291: Integration Demo Requirement (CRITICAL - demo must pass)
- R292: Integration Fixes MUST Be In Effort Branches
- R238: Integration Report Evaluation Protocol (to be created)
- R260: Integration Agent Core Requirements
- R263: Integration Documentation Requirements
- R206: State Machine Transition Validation