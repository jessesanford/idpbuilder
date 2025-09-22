# MASTER_RECOVER_SUB_FAILURE State Rules

## State Purpose
Handle failed sub-orchestrators through recovery strategies including restart, checkpoint recovery, or escalation to manual intervention.

## Entry Criteria
- Sub-orchestrator failure detected
- Failure type categorized
- Recovery strategy determined

## Actions Required

### 1. Analyze Failure Type
```bash
analyze_failure() {
  local SUB_ID="$1"
  local FAILURE_DATA=$(get_failure_data "$SUB_ID")

  # Categorize failure
  local FAILURE_TYPE=$(categorize_failure "$FAILURE_DATA")

  case "$FAILURE_TYPE" in
    "TRANSIENT")
      # Network, temporary resource issues
      strategy="RETRY_IMMEDIATE"
      ;;
    "RECOVERABLE")
      # Has checkpoint, can resume
      strategy="CHECKPOINT_RECOVERY"
      ;;
    "SYSTEMATIC")
      # Repeated same error
      strategy="ESCALATE"
      ;;
    "RESOURCE")
      # Out of memory, disk, etc.
      strategy="RETRY_WITH_LIMITS"
      ;;
    *)
      strategy="ESCALATE"
      ;;
  esac
}
```

### 2. Execute Recovery Strategy

#### RETRY_IMMEDIATE Strategy
```bash
retry_immediate() {
  local SUB_ID="$1"
  local RETRY_COUNT=$(get_retry_count "$SUB_ID")

  if [[ $RETRY_COUNT -lt 3 ]]; then
    # Clean up failed attempt
    cleanup_failed_sub "$SUB_ID"

    # Respawn with same parameters
    respawn_sub_orchestrator "$SUB_ID"

    # Update retry tracking
    increment_retry_count "$SUB_ID"
  else
    # Max retries exceeded
    escalate_failure "$SUB_ID" "MAX_RETRIES_EXCEEDED"
  fi
}
```

#### CHECKPOINT_RECOVERY Strategy
```bash
recover_from_checkpoint() {
  local SUB_ID="$1"
  local CHECKPOINT_FILE="/tmp/sub-orch-${SUB_ID}/checkpoint.json"

  if [[ -f "$CHECKPOINT_FILE" ]]; then
    local CHECKPOINT=$(cat "$CHECKPOINT_FILE")

    # Prepare recovery parameters
    local RECOVERY_PARAMS=$(build_recovery_params "$CHECKPOINT")

    # Spawn recovery sub-orchestrator
    spawn_recovery_sub "$SUB_ID" "$RECOVERY_PARAMS"
  else
    # No checkpoint available, full restart
    retry_immediate "$SUB_ID"
  fi
}
```

#### RETRY_WITH_LIMITS Strategy
```bash
retry_with_resource_limits() {
  local SUB_ID="$1"

  # Apply resource constraints
  local LIMITED_PARAMS=$(cat <<EOF
{
  "resource_limits": {
    "max_memory_mb": 1024,
    "max_cpu_percent": 50,
    "timeout_seconds": 1800
  },
  "original_params": $(get_original_params "$SUB_ID")
}
EOF
)

  # Respawn with limits
  spawn_limited_sub "$SUB_ID" "$LIMITED_PARAMS"
}
```

#### ESCALATE Strategy
```bash
escalate_failure() {
  local SUB_ID="$1"
  local REASON="$2"

  # Preserve failure context
  preserve_failure_artifacts "$SUB_ID"

  # Update state for manual intervention
  update_state_for_escalation "$SUB_ID" "$REASON"

  # Notify and stop
  create_escalation_report "$SUB_ID"
}
```

### 3. Recovery Monitoring
```bash
monitor_recovery_attempt() {
  local SUB_ID="$1"
  local RECOVERY_TYPE="$2"

  # Track recovery metrics
  update_recovery_metrics "$SUB_ID" "$RECOVERY_TYPE"

  # Set enhanced monitoring
  enable_detailed_monitoring "$SUB_ID"

  # Set recovery timeout
  set_recovery_timeout "$SUB_ID" 3600  # 1 hour max
}
```

### 4. Update Master State
```json
{
  "recovery_attempts": {
    "sub_id": "uuid-123",
    "failure_type": "TRANSIENT",
    "recovery_strategy": "RETRY_IMMEDIATE",
    "attempt_number": 2,
    "max_attempts": 3,
    "recovery_started": "timestamp"
  }
}
```

### 5. Failure Pattern Detection
```bash
detect_failure_patterns() {
  local SUB_ID="$1"
  local HISTORY=$(get_failure_history "$SUB_ID")

  # Check for repeated failures
  if has_repeated_failures "$HISTORY"; then
    return 0  # Systematic issue detected
  fi

  # Check for cascade failures
  if has_cascade_pattern "$HISTORY"; then
    return 0  # Underlying issue causing cascades
  fi

  return 1  # Isolated failure
}
```

## Exit Criteria
- Recovery strategy executed
- Sub-orchestrator restarted OR
- Failure escalated for manual intervention
- State updated with recovery status

## Success Transitions
- Recovery successful → `MASTER_MONITOR_SUB_ORCHESTRATORS`
- Escalated → `ERROR_RECOVERY`
- Abandoned → Continue without sub result

## Failure Transitions
- Recovery failed → `ERROR_RECOVERY`
- System error → `HARD_STOP`

## Rules Applied
- R379: Monitoring and Recovery
- R378: Lifecycle Management
- R377: Communication Protocol
- R206: State Machine Validation

## Recovery Decision Matrix

| Failure Type | Retry Count | Has Checkpoint | Action |
|-------------|-------------|----------------|--------|
| TRANSIENT | < 3 | - | Retry immediate |
| TRANSIENT | >= 3 | - | Escalate |
| RECOVERABLE | < 2 | Yes | Checkpoint recovery |
| RECOVERABLE | < 2 | No | Full restart |
| RESOURCE | < 2 | - | Retry with limits |
| SYSTEMATIC | Any | - | Escalate |
| UNKNOWN | < 1 | - | Retry once |
| UNKNOWN | >= 1 | - | Escalate |

## Escalation Report Template
```markdown
# SUB-ORCHESTRATOR FAILURE ESCALATION

## Failure Summary
- Sub ID: uuid-123
- Type: FIX_CASCADE
- Started: 2025-01-21T10:00:00Z
- Failed: 2025-01-21T10:45:00Z
- Attempts: 3

## Failure Details
- Type: SYSTEMATIC
- Error: State corruption in FIX_CASCADE_BACKPORT
- Pattern: Same error on all retry attempts

## Preserved Artifacts
- Logs: /archives/failed/uuid-123/
- State: /archives/failed/uuid-123/state.json
- Checkpoint: /archives/failed/uuid-123/checkpoint.json

## Recommended Actions
1. Manual investigation required
2. Check for environmental issues
3. Review state machine logic
4. Consider alternative approach

## Impact
- Blocked workflows: [list]
- Dependent tasks: [list]
```

## State Updates Required
```json
{
  "sub_orchestrator_failures": [
    {
      "id": "uuid-123",
      "failure_count": 3,
      "last_failure": "timestamp",
      "recovery_attempted": true,
      "escalated": false
    }
  ],
  "recovery_in_progress": {
    "sub_id": "uuid-123",
    "strategy": "CHECKPOINT_RECOVERY",
    "started": "timestamp"
  }
}
```

## Recovery Best Practices
1. **Quick retry** for transient failures
2. **Checkpoint recovery** when available
3. **Resource limits** for exhaustion issues
4. **Escalate early** for systematic problems
5. **Preserve context** for debugging

## Common Issues
1. **Recovery loops** - Same failure recurring
2. **Checkpoint corruption** - Invalid checkpoint data
3. **Resource starvation** - System under pressure
4. **Timeout cascades** - Recovery taking too long

## Notes
- Always preserve failure artifacts
- Consider impact on dependent work
- Clean up failed attempts properly
- Document patterns for future prevention