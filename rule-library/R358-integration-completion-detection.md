# 🔴🔴🔴 RULE R358: INTEGRATION COMPLETION DETECTION AND AUTOMATIC TRANSITION (SUPREME LAW)

## Summary
MONITORING_INTEGRATION state MUST continuously check integration status and transition IMMEDIATELY upon completion. The orchestrator cannot remain in MONITORING_INTEGRATION after integration completes.

## Criticality Level
**SUPREME LAW** - Failure to transition on completion = -50% grade penalty

## Description

The MONITORING_INTEGRATION state exists solely to monitor an active integration process. Once integration completes (successfully or with failure), the orchestrator MUST:

1. **Detect completion within 30 seconds**
2. **Determine the appropriate next state based on outcome**
3. **Update state file immediately**
4. **Transition to the new state**

### Detection Mechanism

The orchestrator MUST check for integration completion using ALL of these methods:

```bash
# Method 1: Check integration_status field
INTEGRATION_STATUS=$(jq -r '.integration_status.status // "unknown"' orchestrator-state.json)
COMPLETED_AT=$(jq -r '.integration_status.completed_at // ""' orchestrator-state.json)

# Method 2: Check for integration report file
REPORT_FILE=$(jq -r '.metadata_locations.integration_reports."phase${PHASE}_wave${WAVE}".file_path // ""' orchestrator-state.json)

# Method 3: Check integration agent process
INTEGRATION_PID=$(pgrep -f "integration-agent" || echo "")

# Integration is complete if:
# - integration_status.status == "completed" OR
# - integration_status.completed_at is not empty OR
# - Integration report file exists OR
# - Integration agent process is not running
```

### Transition Mapping

Based on integration outcome, the orchestrator MUST transition as follows:

| Integration Outcome | Build Status | Test Status | Next State | Reason |
|-------------------|--------------|-------------|------------|---------|
| SUCCESS | PASSING/SUCCESS | PASSING/SUCCESS | INTEGRATION_CODE_REVIEW | All gates passed - need review |
| SUCCESS (cascade mode) | PASSING/SUCCESS | PASSING/SUCCESS | CASCADE_REINTEGRATION | Continue cascade operations |
| FAILED | Any failure | Any | IMMEDIATE_BACKPORT_REQUIRED | R321 - Fix in source branches |
| BLOCKED | blocked/failed_upstream | Any | IMMEDIATE_BACKPORT_REQUIRED | R321 - Fix upstream issue |
| ERROR | missing/error | Any | ERROR_RECOVERY | Integration process error |
| STALE_DETECTED | Any | Any | CASCADE_REINTEGRATION | R327 - Handle stale integrations |

### Continuous Monitoring Loop

```bash
while true; do
    # Check every 30 seconds
    sleep 30

    # Check for completion
    INTEGRATION_COMPLETE=false
    INTEGRATION_OUTCOME=""

    # Method 1: Check state file
    STATUS=$(jq -r '.integration_status.status // "unknown"' orchestrator-state.json)
    COMPLETED_AT=$(jq -r '.integration_status.completed_at // ""' orchestrator-state.json)
    BUILD_STATUS=$(jq -r '.integration_status.build_status // "unknown"' orchestrator-state.json)

    if [[ "$STATUS" == "completed" ]] || [[ -n "$COMPLETED_AT" ]]; then
        INTEGRATION_COMPLETE=true
        INTEGRATION_OUTCOME="$BUILD_STATUS"
    fi

    # Method 2: Check for report
    if [[ -f "$REPORT_FILE" ]]; then
        INTEGRATION_COMPLETE=true
        # Parse outcome from report
        INTEGRATION_OUTCOME=$(grep "^Integration Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    fi

    # Method 3: Check process
    if [[ -z "$INTEGRATION_PID" ]]; then
        # Integration agent not running - must have completed
        INTEGRATION_COMPLETE=true
        if [[ -z "$INTEGRATION_OUTCOME" ]]; then
            INTEGRATION_OUTCOME="PROCESS_ENDED"
        fi
    fi

    # TRANSITION IMMEDIATELY IF COMPLETE
    if [[ "$INTEGRATION_COMPLETE" == "true" ]]; then
        echo "🔴🔴🔴 R358: INTEGRATION COMPLETED - MUST TRANSITION 🔴🔴🔴"
        determine_next_state "$INTEGRATION_OUTCOME"
        update_state_and_stop
    fi
done
```

## Violations

### BLOCKING Violations (-50% to -100%):
- ❌ Remaining in MONITORING_INTEGRATION after integration completes
- ❌ Not checking integration status regularly
- ❌ Ignoring integration completion signals
- ❌ Manual checking instead of continuous monitoring
- ❌ Transitioning to wrong state based on outcome

### WARNING Violations (-20%):
- ⚠️ Checking interval > 30 seconds
- ⚠️ Not checking all three detection methods
- ⚠️ Missing transition reason in state file

## Implementation Requirements

### 1. Entry to MONITORING_INTEGRATION
```bash
# On entering MONITORING_INTEGRATION, start monitoring loop
echo "📊 Starting R358 integration monitoring loop..."
monitor_integration_completion &
MONITOR_PID=$!
```

### 2. Determining Next State
```bash
determine_next_state() {
    local outcome="$1"

    # Check cascade mode first (R351)
    CASCADE_MODE=$(jq -r '.cascade_coordination.cascade_mode // false' orchestrator-state.json)

    # Check for stale integrations (R327)
    STALE_INTEGRATIONS=$(jq -r '.stale_integration_tracking.stale_integrations[]? | select(.recreation_required == true)' orchestrator-state.json)

    if [[ -n "$STALE_INTEGRATIONS" ]]; then
        NEXT_STATE="CASCADE_REINTEGRATION"
        REASON="R327 - Stale integrations detected"
    elif [[ "$outcome" == "SUCCESS" ]] || [[ "$outcome" == "PASSING" ]]; then
        if [[ "$CASCADE_MODE" == "true" ]]; then
            NEXT_STATE="CASCADE_REINTEGRATION"
            REASON="R351 - Continue cascade operations"
        else
            NEXT_STATE="INTEGRATION_CODE_REVIEW"
            REASON="Integration successful - needs review"
        fi
    elif [[ "$outcome" == "FAILED" ]] || [[ "$outcome" == "BLOCKED" ]]; then
        NEXT_STATE="IMMEDIATE_BACKPORT_REQUIRED"
        REASON="R321 - Integration failed, fix in source branches"
    else
        NEXT_STATE="ERROR_RECOVERY"
        REASON="Unexpected integration outcome: $outcome"
    fi
}
```

### 3. Update and Transition
```bash
update_state_and_stop() {
    # Update state file
    jq --arg state "$NEXT_STATE" --arg reason "$REASON" '
        .current_state = $state |
        .state_transition_history += [{
            "from": "MONITORING_INTEGRATION",
            "to": $state,
            "timestamp": now | todate,
            "reason": $reason,
            "r358_enforced": true
        }]' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

    # Commit change
    git add orchestrator-state.json
    git commit -m "state: MONITORING_INTEGRATION → $NEXT_STATE - $REASON (R358 enforced)"
    git push

    # Stop per R322
    echo "🛑 R358: Stopping at state transition checkpoint"
    echo "Integration completed. Next state: $NEXT_STATE"
    echo "Use /continue-orchestrating to proceed"
    exit 0
}
```

## Related Rules

- **R321**: Immediate backport during integration (for failures)
- **R327**: Staleness cascade requirements
- **R351**: CASCADE mode skip patterns
- **R322**: Mandatory stop after state transitions
- **R288**: State file update requirements
- **R346**: Integration state tracking validation

## Responsibility Assignment

### Integration Agent Responsibilities:
1. Update `integration_status` in orchestrator-state.json when starting
2. Update `integration_status` with completion status and timestamp
3. Create integration report at tracked location (R344)
4. Update state file before exiting

### Orchestrator Responsibilities:
1. Start monitoring loop on entering MONITORING_INTEGRATION
2. Check all three detection methods continuously
3. Transition immediately upon detecting completion
4. Never assume integration is still running without verification
5. Handle all possible outcomes correctly

## Example Scenarios

### Scenario 1: Successful Integration
```json
{
  "integration_status": {
    "status": "completed",
    "completed_at": "2025-01-20T19:31:10Z",
    "build_status": "SUCCESS",
    "test_status": "PASSING"
  }
}
```
**Action**: Transition to INTEGRATION_CODE_REVIEW

### Scenario 2: Failed Build
```json
{
  "integration_status": {
    "status": "completed",
    "completed_at": "2025-01-20T19:31:10Z",
    "build_status": "failed_upstream_bug",
    "test_status": "not_run"
  }
}
```
**Action**: Transition to IMMEDIATE_BACKPORT_REQUIRED (R321)

### Scenario 3: Integration Agent Crashed
- Process not running
- No completion status in state file
- No report file
**Action**: Transition to ERROR_RECOVERY

## Enforcement

This rule is enforced by:
1. State machine validation (R206)
2. Orchestrator grading criteria
3. Automated state consistency checks
4. Integration test suites

**Penalty for violation**: -50% for stuck states, -100% for wrong transitions

---

**Remember**: MONITORING states exist to MONITOR, not to wait passively. Detection and transition must be AUTOMATIC and IMMEDIATE!