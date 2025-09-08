# 🚨🚨🚨 BLOCKING RULE R019: Error Recovery Protocol

## Rule Statement
All agents MUST follow a systematic error recovery protocol when encountering failures, preserving state and context while attempting automated recovery before escalation.

## Requirements

### 1. Error Classification
Agents MUST classify errors by severity:
- **CRITICAL**: System corruption, data loss risk, security breach, BUILD/TEST GATE FAILURES (R291)
- **HIGH**: Workflow blocked, multi-agent failure, integration broken, demo failures
- **MEDIUM**: Single agent failure, recoverable state issue, partial test failures
- **LOW**: Transient errors, network hiccups, retry-able operations

**🔴 SPECIAL CLASSIFICATION: R291 BUILD/TEST GATE FAILURES = ALWAYS CRITICAL 🔴**
- Build failure → CRITICAL → ERROR_RECOVERY
- Test failure → CRITICAL → ERROR_RECOVERY  
- Demo failure → CRITICAL → ERROR_RECOVERY
- No artifacts → CRITICAL → ERROR_RECOVERY

### 2. State Preservation
Before any recovery attempt:
```bash
# Save current state
save_error_state() {
    echo "ERROR_STATE: $(date '+%Y-%m-%d %H:%M:%S')" > error-state.yaml
    echo "agent: $AGENT_NAME" >> error-state.yaml
    echo "state: $CURRENT_STATE" >> error-state.yaml
    echo "error: $ERROR_MESSAGE" >> error-state.yaml
    echo "context: $CONTEXT" >> error-state.yaml
    
    # Save TODO state
    save_todos "ERROR_RECOVERY_CHECKPOINT"
    
    # Git commit error state
    git add -A
    git commit -m "error: checkpoint before recovery attempt"
    git push
}
```

### 3. Recovery Sequence
Follow this exact sequence:

#### Step 1: Immediate Assessment
```bash
# Assess error severity and impact
assess_error() {
    local error_type="$1"
    case "$error_type" in
        "git_corruption"|"data_loss"|"security_breach")
            echo "CRITICAL"
            ;;
        "workflow_blocked"|"integration_failed"|"multi_agent_failure")
            echo "HIGH"
            ;;
        "single_agent_failure"|"state_issue")
            echo "MEDIUM"
            ;;
        "network"|"timeout"|"transient")
            echo "LOW"
            ;;
    esac
}
```

#### Step 2: Automated Recovery Attempts
- **LOW**: Retry up to 3 times with exponential backoff
- **MEDIUM**: Reset agent state, reload configuration, retry once
- **HIGH**: Coordinate with orchestrator, attempt rollback
- **CRITICAL**: Preserve everything, alert immediately, no automated attempts

#### Step 3: Escalation Path
If automated recovery fails:
1. Save detailed error report
2. Update orchestrator-state.json to ERROR_RECOVERY
3. Alert user with recovery options
4. Wait for manual intervention or orchestrator directive

### 4. Recovery Validation
After successful recovery:
```bash
validate_recovery() {
    # Verify state consistency
    check_state_consistency
    
    # Verify git status
    git status --porcelain
    
    # Verify TODO state recovered
    verify_todo_recovery
    
    # Run health checks
    run_agent_health_checks
}
```

## Violations
- Attempting recovery without state preservation: -50% penalty
- Not classifying error severity: -20% penalty
- Skipping recovery validation: -30% penalty
- Continuing after CRITICAL error without explicit clearance: -100% FAIL

## Enforcement
- All agents MUST implement error handlers
- Recovery attempts MUST be logged
- State MUST be preserved before any recovery
- Orchestrator monitors all ERROR_RECOVERY states

## Integration Rules
- Works with R156 (Error Recovery Time Targets)
- Supports R021 (Orchestrator Never Stops)
- Enables R187-R190 (TODO persistence during errors)

## Examples

### Example 1: Git Corruption (CRITICAL)
```bash
# Detection
if git status 2>&1 | grep -q "fatal:"; then
    save_error_state
    echo "🚨 CRITICAL: Git corruption detected"
    # DO NOT attempt automated recovery
    # Preserve everything and wait
fi
```

### Example 2: Agent Crash (MEDIUM)
```bash
# Detection
if agent_health_check_failed; then
    save_error_state
    echo "⚠️ MEDIUM: Agent crash detected"
    # Attempt recovery
    reload_agent_config
    restart_agent
    validate_recovery
fi
```

### Example 3: Network Timeout (LOW)
```bash
# Detection
if curl_failed_with_timeout; then
    echo "📝 LOW: Network timeout"
    # Simple retry with backoff
    for i in 1 2 4; do
        sleep $i
        retry_operation && break
    done
fi
```

## Metrics
- Mean Time To Recovery (MTTR) by severity
- Recovery success rate by error type
- State preservation compliance rate
- Escalation frequency

## Update History
- 2025-01-29: Created rule based on system requirements
- Defines systematic error recovery protocol
- Ensures state preservation and proper escalation