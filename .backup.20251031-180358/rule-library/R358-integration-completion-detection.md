# 🔴🔴🔴 RULE R358: INTEGRATE_WAVE_EFFORTS COMPLETION DETECTION AND AUTOMATIC TRANSITION (SUPREME LAW)

## Summary
MONITORING_INTEGRATE_WAVE_EFFORTS state MUST continuously check integration status and transition IMMEDIATELY upon completion. The orchestrator cannot remain in MONITORING_INTEGRATE_WAVE_EFFORTS after integration completes.

## Criticality Level
**SUPREME LAW** - Failure to transition on completion = -50% grade penalty

## Description

The MONITORING_INTEGRATE_WAVE_EFFORTS state exists solely to monitor an active integration process. Once integration completes (successfully or with failure), the orchestrator MUST:

1. **Detect completion within 30 seconds**
2. **Determine the appropriate next state based on outcome**
3. **Update state file immediately**
4. **Transition to the new state**

### Detection Mechanism

The orchestrator MUST check for integration completion using ALL of these methods:

```bash
# Method 1: Check integration_status field
INTEGRATE_WAVE_EFFORTS_STATUS=$(jq -r '.integration_status.status // "unknown"' orchestrator-state-v3.json)
COMPLETED_AT=$(jq -r '.integration_status.completed_at // ""' orchestrator-state-v3.json)

# Method 2: Check for integration report file
REPORT_FILE=$(jq -r '.metadata_locations.integration_reports."phase${PHASE}_wave${WAVE}".file_path // ""' orchestrator-state-v3.json)

# Method 3: Check integration agent process
INTEGRATE_WAVE_EFFORTS_PID=$(pgrep -f "integration-agent" || echo "")

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
| PROJECT_DONE | PASSING/PROJECT_DONE | PASSING/PROJECT_DONE | REVIEW_WAVE_INTEGRATION | All gates passed - need review |
| PROJECT_DONE (cascade mode) | PASSING/PROJECT_DONE | PASSING/PROJECT_DONE | CASCADE_REINTEGRATION | Continue cascade operations |
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
    INTEGRATE_WAVE_EFFORTS_COMPLETE=false
    INTEGRATE_WAVE_EFFORTS_OUTCOME=""

    # Method 1: Check state file
    STATUS=$(jq -r '.integration_status.status // "unknown"' orchestrator-state-v3.json)
    COMPLETED_AT=$(jq -r '.integration_status.completed_at // ""' orchestrator-state-v3.json)
    BUILD_STATUS=$(jq -r '.integration_status.build_status // "unknown"' orchestrator-state-v3.json)

    if [[ "$STATUS" == "completed" ]] || [[ -n "$COMPLETED_AT" ]]; then
        INTEGRATE_WAVE_EFFORTS_COMPLETE=true
        INTEGRATE_WAVE_EFFORTS_OUTCOME="$BUILD_STATUS"
    fi

    # Method 2: Check for report
    if [[ -f "$REPORT_FILE" ]]; then
        INTEGRATE_WAVE_EFFORTS_COMPLETE=true
        # Parse outcome from report
        INTEGRATE_WAVE_EFFORTS_OUTCOME=$(grep "^Integration Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    fi

    # Method 3: Check process
    if [[ -z "$INTEGRATE_WAVE_EFFORTS_PID" ]]; then
        # Integration agent not running - must have completed
        INTEGRATE_WAVE_EFFORTS_COMPLETE=true
        if [[ -z "$INTEGRATE_WAVE_EFFORTS_OUTCOME" ]]; then
            INTEGRATE_WAVE_EFFORTS_OUTCOME="PROCESS_ENDED"
        fi
    fi

    # TRANSITION IMMEDIATELY IF COMPLETE
    if [[ "$INTEGRATE_WAVE_EFFORTS_COMPLETE" == "true" ]]; then
        echo "🔴🔴🔴 R358: INTEGRATE_WAVE_EFFORTS COMPLETED - MUST TRANSITION 🔴🔴🔴"
        determine_next_state "$INTEGRATE_WAVE_EFFORTS_OUTCOME"
        update_state_and_stop
    fi
done
```

## Violations

### BLOCKING Violations (-50% to -100%):
- ❌ Remaining in MONITORING_INTEGRATE_WAVE_EFFORTS after integration completes
- ❌ Not checking integration status regularly
- ❌ Ignoring integration completion signals
- ❌ Manual checking instead of continuous monitoring
- ❌ Transitioning to wrong state based on outcome

### WARNING Violations (-20%):
- ⚠️ Checking interval > 30 seconds
- ⚠️ Not checking all three detection methods
- ⚠️ Missing transition reason in state file

## Implementation Requirements

### 1. Entry to MONITORING_INTEGRATE_WAVE_EFFORTS
```bash
# On entering MONITORING_INTEGRATE_WAVE_EFFORTS, start monitoring loop
echo "📊 Starting R358 integration monitoring loop..."
monitor_integration_completion &
MONITOR_PID=$!
```

### 2. Determining Next State
```bash
determine_next_state() {
    local outcome="$1"

    # Check cascade mode first (R351)
    CASCADE_MODE=$(jq -r '.cascade_coordination.cascade_mode // false' orchestrator-state-v3.json)

    # Check for stale integrations (R327)
    STALE_INTEGRATE_WAVE_EFFORTSS=$(jq -r '.stale_integration_tracking.stale_integrations[]? | select(.recreation_required == true)' orchestrator-state-v3.json)

    if [[ -n "$STALE_INTEGRATE_WAVE_EFFORTSS" ]]; then
        NEXT_STATE="CASCADE_REINTEGRATION"
        REASON="R327 - Stale integrations detected"
    elif [[ "$outcome" == "PROJECT_DONE" ]] || [[ "$outcome" == "PASSING" ]]; then
        if [[ "$CASCADE_MODE" == "true" ]]; then
            NEXT_STATE="CASCADE_REINTEGRATION"
            REASON="R351 - Continue cascade operations"
        else
            NEXT_STATE="REVIEW_WAVE_INTEGRATION"
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
        .state_machine.current_state = $state |
        .state_transition_history += [{
            "from": "MONITORING_INTEGRATE_WAVE_EFFORTS",
            "to": $state,
            "timestamp": now | todate,
            "reason": $reason,
            "r358_enforced": true
        }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

    # Commit change
    git add orchestrator-state-v3.json
    git commit -m "state: MONITORING_INTEGRATE_WAVE_EFFORTS → $NEXT_STATE - $REASON (R358 enforced)"
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
1. Update `integration_status` in orchestrator-state-v3.json when starting
2. Update `integration_status` with completion status and timestamp
3. Create integration report at tracked location (R344)
4. Update state file before exiting

### Orchestrator Responsibilities:
1. Start monitoring loop on entering MONITORING_INTEGRATE_WAVE_EFFORTS
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
    "build_status": "PROJECT_DONE",
    "test_status": "PASSING"
  }
}
```
**Action**: Transition to REVIEW_WAVE_INTEGRATION

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

**Remember**: MONITORING_SWE_PROGRESS states exist to MONITOR, not to wait passively. Detection and transition must be AUTOMATIC and IMMEDIATE!