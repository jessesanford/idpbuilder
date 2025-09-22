# 🚨🚨🚨 BLOCKING RULE R379: Process-Based Monitoring and Recovery 🚨🚨🚨

## Rule Definition
**ID**: R379
**Category**: ORCHESTRATION
**Criticality**: BLOCKING
**Created**: 2025-01-21
**Updated**: 2025-01-21 (Process-Based Architecture)
**Dependencies**: R377 (Communication), R378 (Lifecycle), R206 (State Validation)

## Description
Mandates continuous monitoring of sub-orchestrator background processes using BashOutput tool for real-time output monitoring and KillShell for process management. Ensures system resilience through output-based health checks and intelligent recovery strategies without file polling.

## 🔴 MONITORING REQUIREMENTS

### 1. PROCESS OUTPUT MONITORING PROTOCOL

```python
# Master monitoring loop using BashOutput
def monitor_processes():
    """Monitor all sub-orchestrators via their shell output"""
    while True:
        for sub in list_active_subs():
            check_sub_output(sub['shell_id'], sub['id'])
        time.sleep(15)  # Check every 15 seconds

def check_sub_output(shell_id, sub_id):
    """Check sub-orchestrator health via output"""
    try:
        # Get new output since last check
        output = BashOutput(bash_id=shell_id)

        if output:
            # Process has produced output
            update_last_output_time(sub_id)

            # Parse for health signals
            if "SUB_ORCHESTRATOR_ERROR:" in output:
                handle_error_signal(sub_id, output)
            elif "SUB_ORCHESTRATOR_STUCK:" in output:
                handle_stuck_signal(sub_id, output)
            elif "SUB_ORCHESTRATOR_COMPLETE" in output:
                handle_completion(sub_id, output)
            elif "SUB_ORCHESTRATOR_FAILED" in output:
                handle_failure(sub_id, output)

            # Update progress
            if "PROGRESS:" in output:
                update_progress(sub_id, extract_progress(output))

        else:
            # No output - check if stuck
            if time_since_last_output(sub_id) > 300:  # 5 minutes
                handle_no_output(sub_id, shell_id)

    except Exception as e:
        # Shell no longer exists
        handle_terminated_process(sub_id, str(e))
```

### 2. PROCESS HEALTH CHECKS

```python
# Multi-level health validation using process monitoring
def validate_sub_health(sub_id, shell_id):
    """Validate sub-orchestrator health via process checks"""

    # Level 1: Check if shell is alive
    if not is_shell_alive(shell_id):
        handle_dead_process(sub_id)
        return False

    # Level 2: Check output responsiveness
    test_output = BashOutput(bash_id=shell_id)
    if test_output is None:
        handle_terminated_shell(sub_id)
        return False

    # Level 3: Progress advancing
    if not check_progress_advancing(sub_id):
        # Send progress request signal
        request_progress_update(sub_id)
        return False

    # Level 4: Check for error signals
    if has_error_signals(test_output):
        handle_error_condition(sub_id, test_output)
        return False

    return True  # Healthy

def is_shell_alive(shell_id):
    """Check if background shell is still running"""
    try:
        # Attempt to get output - will fail if shell dead
        BashOutput(bash_id=shell_id)
        return True
    except:
        return False
```

### 3. PROGRESS TRACKING

```json
{
  "progress_tracking": {
    "sub_id": "uuid-123",
    "shell_id": "bash_abc123",
    "progress_history": [
      {"timestamp": "10:00:00", "percentage": "0%"},
      {"timestamp": "10:05:00", "percentage": "25%"},
      {"timestamp": "10:10:00", "percentage": "50%"},
      {"timestamp": "10:15:00", "percentage": "50%"}  // Stuck!
    ],
    "last_output": "10:15:30",
    "no_output_threshold_seconds": 300,
    "is_stuck": true
  }
}
```

## 🛑 FAILURE DETECTION

### 1. FAILURE TYPES (Process-Based)

```python
def categorize_failure(sub_id, symptoms):
    """Categorize failures based on process behavior"""

    if "shell_terminated" in symptoms:
        return "PROCESS_TERMINATED"
    elif "no_output" in symptoms:
        return "PROCESS_HUNG"
    elif "progress_stuck" in symptoms:
        return "STUCK_STATE"
    elif "error_signal" in symptoms:
        return "ERROR_REPORTED"
    elif "resource_exceeded" in symptoms:
        return "RESOURCE_EXHAUSTION"
    elif "FAILED" in symptoms:
        return "EXPLICIT_FAILURE"
    else:
        return "UNKNOWN_FAILURE"
```

### 2. DETECTION THRESHOLDS

| Condition | Threshold | Action |
|-----------|-----------|--------|
| No output | 5 minutes | Check if alive |
| Shell terminated | Immediate | Recovery |
| Progress stuck | 10 minutes | Send UNSTICK signal |
| ERROR signal | Immediate | Analyze error |
| FAILED signal | Immediate | Read failure state |
| Memory high | In output | Resource warning |

## 🔄 RECOVERY MECHANISMS

### 1. RECOVERY DECISION TREE

```python
def determine_recovery_strategy(failure_type, retry_count, has_state_file):
    """Determine recovery strategy based on failure type"""

    if retry_count >= 3:
        return "ESCALATE"  # Too many failures

    if failure_type == "PROCESS_TERMINATED":
        if has_state_file:
            return "RESUME_FROM_STATE"
        else:
            return "RESTART_FRESH"

    elif failure_type == "PROCESS_HUNG":
        return "KILL_AND_RESTART"

    elif failure_type == "STUCK_STATE":
        return "SEND_UNSTICK_SIGNAL"

    elif failure_type == "RESOURCE_EXHAUSTION":
        return "RESTART_WITH_LIMITS"

    else:
        return "ESCALATE"  # Unknown failure
```

### 2. STATE-BASED RECOVERY

```python
# Sub-orchestrator writes state to file
def write_state_file(state, data):
    """Write recovery state to file"""
    state_file = f"/tmp/sub-orchestrators/{SUB_ID}/state.json"
    with open(state_file, 'w') as f:
        json.dump({
            "timestamp": datetime.now().isoformat(),
            "state": state,
            "data": data,
            "can_resume": True
        }, f)

# Master performs recovery from state
def recover_from_state(sub_id):
    """Recover sub-orchestrator from saved state"""
    state_file = f"/tmp/sub-orchestrators/{sub_id}/state.json"

    if os.path.exists(state_file):
        with open(state_file, 'r') as f:
            state = json.load(f)

        if state.get('can_resume'):
            # Spawn recovery process using appropriate continue command
            sub_type = state.get('sub_type', 'orchestrating')
            if sub_type == 'FIX_CASCADE':
                command = '/fix-cascade'
            elif sub_type == 'SPLITTING':
                command = '/splitting'
            else:
                command = '/continue-orchestrating'

            result = Bash(
                command=f"claude -p . --command '{command}' --params 'resume=true state_file={state_file}'",
                run_in_background=True,
                description=f"Resume {sub_type} from saved state"
            )
            track_recovery(sub_id, result['shell_id'])
        else:
            restart_sub_orchestrator(sub_id)
```

### 3. GRACEFUL TERMINATION

```python
def attempt_graceful_recovery(sub_id, shell_id):
    """Attempt graceful recovery of stuck process"""

    # Step 1: Request checkpoint via output
    send_checkpoint_request(sub_id)
    time.sleep(5)

    # Step 2: Check for checkpoint confirmation
    output = BashOutput(bash_id=shell_id)
    if "CHECKPOINT_SAVED" in output:
        # Step 3: Terminate shell gracefully
        KillShell(shell_id=shell_id)

        # Step 4: Recover from checkpoint
        recover_from_state(sub_id)
        return True

    # Step 5: Force kill if no response
    KillShell(shell_id=shell_id)
    return False

def send_checkpoint_request(sub_id):
    """Send checkpoint request to sub-orchestrator"""
    # Could write to a control file that sub monitors
    control_file = f"/tmp/sub-orchestrators/{sub_id}/control"
    with open(control_file, 'w') as f:
        f.write("CHECKPOINT_REQUEST")
```

## ⚠️ ESCALATION PROCEDURES

### 1. ESCALATION TRIGGERS

```python
def should_escalate(sub_id):
    """Determine if escalation is needed"""

    retry_count = get_retry_count(sub_id)
    if retry_count >= 3:
        return "MAX_RETRIES_EXCEEDED"

    failure_pattern = get_failure_pattern(sub_id)
    if failure_pattern == "REPEATED_SAME_ERROR":
        return "SYSTEMATIC_FAILURE"

    runtime = get_runtime(sub_id)
    if runtime > 7200:  # 2 hours
        return "TIMEOUT_EXCEEDED"

    return None  # No escalation needed
```

### 2. ESCALATION ACTIONS

```json
{
  "escalation": {
    "level": "ERROR_RECOVERY",
    "sub_orchestrator_failed": {
      "id": "uuid-123",
      "shell_id": "bash_abc123",
      "type": "FIX_CASCADE",
      "failure_reason": "MAX_RETRIES_EXCEEDED",
      "attempts": 3,
      "last_output": "ERROR: State corruption detected"
    },
    "recommended_action": "MANUAL_INTERVENTION",
    "preserved_artifacts": [
      "/tmp/sub-orchestrators/uuid-123/",
      "archives/output-uuid-123.log",
      "states/uuid-123-last.json"
    ]
  }
}
```

## 🔴 MONITORING DASHBOARD

### Required Metrics:
```json
{
  "monitoring_metrics": {
    "active_shells": 3,
    "healthy_processes": 2,
    "warning_processes": 1,
    "failed_processes": 0,
    "recovery_in_progress": 1,
    "total_recoveries": 5,
    "success_rate": 0.95,
    "average_runtime": 1800,
    "shell_statistics": {
      "total_output_lines": 4500,
      "errors_detected": 2,
      "progress_updates": 45
    }
  }
}
```

## 🚨 IMPLEMENTATION REQUIREMENTS

### Monitoring Infrastructure:
```python
def setup_monitoring():
    """Initialize process monitoring"""
    # No daemon needed - use main orchestrator loop

    # Track active shells in state
    state = {
        "active_sub_orchestrators": [],
        "monitoring_config": {
            "check_interval": 15,
            "no_output_timeout": 300,
            "max_runtime": 7200
        }
    }

    # Start monitoring in orchestrator's main loop
    return state

def monitor_loop():
    """Main monitoring loop in orchestrator"""
    while has_active_subs():
        for sub in get_active_subs():
            # Check output
            output = BashOutput(bash_id=sub['shell_id'])
            process_output(sub['id'], output)

            # Check health
            if not validate_sub_health(sub['id'], sub['shell_id']):
                trigger_recovery(sub['id'])

        # Update metrics
        update_monitoring_dashboard()

        # Clean terminated shells
        clean_terminated_shells()

        time.sleep(15)
```

## 📋 Recovery Testing

### Test Scenarios:
1. **Process termination** - Kill shell unexpectedly
2. **No output** - Sub-orchestrator goes silent
3. **Error signals** - Inject ERROR output
4. **State recovery** - Resume from state file
5. **Resource limits** - Memory/CPU stress
6. **Graceful shutdown** - Request checkpoint
7. **Multiple failures** - Cascading issues

### Test Implementation:
```python
def test_recovery_scenarios():
    """Test various recovery scenarios"""

    # Test 1: Kill shell
    shell_id = spawn_test_sub()
    KillShell(shell_id=shell_id)
    assert recovery_triggered(shell_id)

    # Test 2: Simulate hang
    shell_id = spawn_silent_sub()
    time.sleep(360)  # Wait for timeout
    assert marked_as_hung(shell_id)

    # Test 3: Error signal
    shell_id = spawn_error_sub()
    assert error_handled(shell_id)
```

## 🚨 VIOLATIONS

### BLOCKING (-100% penalty):
1. No monitoring of active shells
2. Failed to detect terminated process
3. No recovery attempt after failure detection
4. Lost sub-orchestrator output/results
5. Shell leaks (not cleaned up)

### WARNING (-25% penalty):
1. Monitoring interval > 30 seconds
2. Missing state file capabilities
3. No output parsing for signals
4. Poor failure categorization
5. Incomplete escalation data

## 📋 Implementation Checklist

- [ ] BashOutput monitoring implemented
- [ ] Signal parsing for all standard signals
- [ ] KillShell cleanup on completion
- [ ] State file recovery mechanism
- [ ] No-output timeout detection
- [ ] Graceful termination attempt
- [ ] Escalation thresholds defined
- [ ] Monitoring metrics tracked

## 🔗 Related Rules
- R377: Communication Protocol (Output Signals)
- R378: Lifecycle Management (Background Processes)
- R313: Stop After Spawn
- R287: TODO Persistence

## 🔴 Key Changes in Process-Based Monitoring

1. **No Heartbeat Files** - Monitor via BashOutput
2. **Shell-Based Health** - Check shell alive, not PID
3. **Output Signals** - Parse standardized signals
4. **Clean Termination** - KillShell for cleanup
5. **Real-time Monitoring** - Direct output stream

---

**Remember**: Monitor background processes via BashOutput, not file polling. Every shell must be tracked and cleaned up properly!