# ⚠️⚠️⚠️ WARNING RULE R156: Error Recovery Time Targets

## Rule Statement
All error recovery operations MUST complete within defined time targets based on error severity, with automatic escalation if targets are exceeded.

## Requirements

### 1. Recovery Time Targets by Severity

| Severity | Maximum Recovery Time | Escalation After | Penalty for Violation |
|----------|----------------------|------------------|----------------------|
| CRITICAL | 30 minutes | 15 minutes | -50% grade |
| HIGH | 60 minutes | 30 minutes | -30% grade |
| MEDIUM | 2 hours | 1 hour | -20% grade |
| LOW | 4 hours | 2 hours | -10% grade |

### 2. Time Tracking Requirements

Every error recovery MUST track:
```bash
# Start recovery timer
start_recovery() {
    local severity="$1"
    local start_time=$(date +%s)
    
    echo "recovery_start: $start_time" >> recovery-metrics.yaml
    echo "severity: $severity" >> recovery-metrics.yaml
    echo "target_minutes: $(get_target_minutes $severity)" >> recovery-metrics.yaml
    
    # Set escalation timer
    schedule_escalation $severity $start_time
}

# Check recovery time
check_recovery_time() {
    local start_time="$1"
    local current_time=$(date +%s)
    local elapsed=$((current_time - start_time))
    local elapsed_minutes=$((elapsed / 60))
    
    echo "Elapsed: ${elapsed_minutes} minutes"
    return $elapsed_minutes
}
```

### 3. Escalation Triggers

Automatic escalation MUST occur at:
- **CRITICAL**: 15 minutes (50% of target)
- **HIGH**: 30 minutes (50% of target)
- **MEDIUM**: 60 minutes (50% of target)
- **LOW**: 120 minutes (50% of target)

```bash
schedule_escalation() {
    local severity="$1"
    local start_time="$2"
    
    case "$severity" in
        "CRITICAL")
            set_timer 15 "escalate_critical_error"
            ;;
        "HIGH")
            set_timer 30 "escalate_high_error"
            ;;
        "MEDIUM")
            set_timer 60 "escalate_medium_error"
            ;;
        "LOW")
            set_timer 120 "escalate_low_error"
            ;;
    esac
}
```

### 4. Recovery Optimization Requirements

For faster recovery:
- **Parallel Recovery**: Attempt multiple recovery strategies simultaneously for CRITICAL/HIGH
- **Cached Solutions**: Maintain database of successful recovery patterns
- **Preemptive Checks**: Run health checks every 5 minutes to catch issues early
- **Fast Rollback**: Keep last 3 known-good states for instant rollback

### 5. Monitoring and Alerts

```bash
monitor_recovery_progress() {
    local severity="$1"
    local start_time="$2"
    local target_minutes=$(get_target_minutes $severity)
    
    while in_recovery; do
        local elapsed=$(check_recovery_time $start_time)
        local remaining=$((target_minutes - elapsed))
        
        # Alert at 25%, 50%, 75% of target
        case $((elapsed * 100 / target_minutes)) in
            25) echo "⚠️ 25% of recovery time used" ;;
            50) echo "⚠️⚠️ 50% of recovery time used - escalating" ;;
            75) echo "🚨 75% of recovery time used - critical" ;;
            100) echo "🚨🚨🚨 Recovery time exceeded - VIOLATION" ;;
        esac
        
        sleep 60  # Check every minute
    done
}
```

## Violations

### Major Violations (Immediate Penalty)
- Exceeding recovery time target: Penalty per severity table
- Not tracking recovery time: -20% penalty
- Ignoring escalation triggers: -30% penalty
- Not attempting parallel recovery for CRITICAL: -40% penalty

### Minor Violations (Warning)
- Not logging recovery metrics: Warning
- Missing progress updates: Warning
- Delayed escalation (within 10% buffer): Warning

## Enforcement

### Automatic Enforcement
```bash
enforce_recovery_time() {
    local severity="$1"
    local start_time="$2"
    local max_minutes=$(get_target_minutes $severity)
    
    # Hard stop at target time
    timeout ${max_minutes}m recovery_attempt || {
        echo "🚨🚨🚨 RECOVERY TIME EXCEEDED"
        force_escalation $severity
        record_violation "R156" $severity $max_minutes
    }
}
```

### Metrics Collection
- Track MTTR (Mean Time To Recovery) per severity
- Track escalation frequency
- Track recovery success rate within target
- Generate weekly recovery performance reports

## Integration with Other Rules

- **R019**: Provides time boundaries for error recovery protocol
- **R021**: Ensures orchestrator continues despite time pressures
- **R187-R190**: TODO state must be saved before timeout
- **R206**: State transitions to ERROR_RECOVERY respect time limits

## Examples

### Example 1: CRITICAL Error (30-minute target)
```bash
# Git corruption detected
start_recovery "CRITICAL"
# At 15 minutes: Automatic escalation
# At 30 minutes: Violation recorded, manual intervention required
```

### Example 2: HIGH Error (60-minute target)
```bash
# Integration test failure blocking deployment
start_recovery "HIGH"
# Parallel attempts: rollback, fix, bypass
# At 30 minutes: Escalate to orchestrator
# At 60 minutes: Violation if not resolved
```

### Example 3: MEDIUM Error (2-hour target)
```bash
# Single agent configuration issue
start_recovery "MEDIUM"
# Sequential attempts with state preservation
# At 1 hour: Escalation warning
# At 2 hours: Violation if unresolved
```

## Recovery Time Optimization Strategies

1. **Pre-computed Recovery Plans**: Maintain library of solutions
2. **Parallel Execution**: Run multiple fixes simultaneously
3. **Fast Failure Detection**: Fail fast to try next strategy
4. **State Caching**: Keep recent good states for quick rollback
5. **Automated Escalation**: Don't wait for human intervention

## Reporting Format
```yaml
recovery_report:
  timestamp: "2025-01-29T10:30:00Z"
  severity: "HIGH"
  error_type: "integration_failure"
  start_time: "2025-01-29T10:00:00Z"
  resolution_time: "2025-01-29T10:45:00Z"
  elapsed_minutes: 45
  target_minutes: 60
  within_target: true
  escalated: true
  escalation_time: "2025-01-29T10:30:00Z"
  recovery_method: "rollback_and_retry"
  violations: []
```

## Update History
- 2025-01-29: Created rule with specific time targets
- Defines escalation triggers at 50% of target time
- Integrates with error recovery protocol R019