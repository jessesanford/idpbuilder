# MASTER_MONITOR_SUB_ORCHESTRATORS State Rules

## State Purpose
Continuously monitor active sub-orchestrators for health, progress, and completion using BashOutput tool to check process output in real-time.

## Entry Criteria
- At least one active sub-orchestrator running
- Shell IDs tracked in master state
- Background processes confirmed running

## Actions Required

### 1. Monitor All Active Sub-Orchestrators
```python
# Monitor all active sub-orchestrators using BashOutput
for sub in active_sub_orchestrators:
    output = BashOutput(bash_id=sub['shell_id'])
    process_sub_output(sub['id'], output)
    detect_completion(sub['id'], output)
```

### 2. Process Output Check Protocol
```python
def check_sub_output(sub_id, shell_id):
    # Get new output since last check
    output = BashOutput(bash_id=shell_id)

    # Parse output for status signals
    if "SUB_ORCHESTRATOR_STARTED" in output:
        update_sub_status(sub_id, "RUNNING")

    if "SUB_ORCHESTRATOR_PROGRESS:" in output:
        # Extract and track progress
        progress = parse_progress(output)
        update_sub_progress(sub_id, progress)

    if "SUB_ORCHESTRATOR_ERROR:" in output:
        # Handle errors
        error = parse_error(output)
        handle_sub_error(sub_id, error)

    if "SUB_ORCHESTRATOR_COMPLETE" in output:
        # Mark as complete
        mark_complete(sub_id)

    if "SUB_ORCHESTRATOR_FAILED" in output:
        # Mark as failed
        mark_failed(sub_id)

    return output
```

### 3. Progress Monitoring
```python
def monitor_sub_progress(sub_id, shell_id):
    output = BashOutput(bash_id=shell_id)

    # Look for progress markers in output
    progress_lines = [l for l in output.split('\n') if 'PROGRESS:' in l]

    if progress_lines:
        latest_progress = progress_lines[-1]
        update_progress(sub_id, latest_progress)

    # Check if process is still producing output
    if not output and time_since_last_output(sub_id) > 300:  # 5 minutes
        mark_as_stuck(sub_id)
        attempt_recovery(sub_id)
```

### 4. Completion Detection
```python
def detect_completion(sub_id, output):
    # Check for completion signals in output
    if "SUB_ORCHESTRATOR_COMPLETE" in output:
        # Read final state from state file
        state_file = get_state_file(sub_id)
        final_state = read_json(state_file)

        # Read output file
        output_file = get_output_file(sub_id)
        results = read_json(output_file)

        # Process completion
        handle_successful_completion(sub_id, results)

    elif "SUB_ORCHESTRATOR_FAILED" in output:
        # Extract failure reason
        failure_reason = parse_failure(output)
        handle_failed_completion(sub_id, failure_reason)

    elif "SUB_ORCHESTRATOR_PARTIAL" in output:
        # Handle partial completion
        partial_results = parse_partial(output)
        handle_partial_completion(sub_id, partial_results)
```

### 5. Update Monitoring Metrics
```json
{
  "monitoring_snapshot": {
    "timestamp": "2025-01-21T10:30:00Z",
    "active_subs": 3,
    "running": 2,
    "completed": 1,
    "failed": 0,
    "details": [
      {
        "id": "uuid-1",
        "shell_id": "bash_abc123",
        "status": "RUNNING",
        "progress": "75%",
        "last_output": "timestamp",
        "output_lines": 450
      }
    ]
  }
}
```

## Exit Criteria
- All sub-orchestrators completed OR
- Critical failure requiring intervention OR
- Timeout reached (configurable)

## Success Transitions
- All complete → `MASTER_HANDLE_SUB_COMPLETION`
- Partial complete → `MASTER_HANDLE_SUB_COMPLETION`
- Need more subs → `MASTER_SPAWN_SUB_ORCHESTRATOR`

## Failure Transitions
- Unrecoverable failure → `MASTER_RECOVER_SUB_FAILURE`
- System error → `ERROR_RECOVERY`
- Resource exhaustion → `RESOURCE_WAIT`

## Rules Applied
- R377: Communication Protocol (via output)
- R378: Sub-Orchestrator Lifecycle (Background Processes)
- R379: Process-Based Monitoring
- R206: State Machine Validation
- R287: TODO Persistence (monitor checkpoints)

## Monitoring Loop Implementation
```python
def monitor_loop():
    iterations = 0
    max_iterations = 240  # 1 hour max

    while has_active_subs() and iterations < max_iterations:
        # Check all subs
        for sub in list_active_subs():
            # Get output from background process
            output = BashOutput(bash_id=sub['shell_id'])

            # Process output
            if output:
                process_sub_output(sub['id'], output)
                detect_completion(sub['id'], output)
                update_last_output_time(sub['id'])

            # Check if process died unexpectedly
            if process_terminated(sub['shell_id']):
                handle_unexpected_termination(sub['id'])

        # Update metrics
        update_monitoring_dashboard()

        # Save checkpoint every minute
        if iterations % 4 == 0:
            save_monitoring_checkpoint()

        time.sleep(15)
        iterations += 1
```

## Process Management
```python
def process_terminated(shell_id):
    """Check if background process has terminated"""
    try:
        # BashOutput will indicate if process ended
        output = BashOutput(bash_id=shell_id)
        # Look for shell termination indicators
        return "Process completed" in output or "Process terminated" in output
    except Exception as e:
        # Shell no longer exists
        return True

def kill_stuck_process(shell_id):
    """Kill a stuck background process"""
    KillShell(shell_id=shell_id)
```

## Recovery Triggers
| Condition | Action | Priority |
|-----------|--------|----------|
| No output 5+ min | Check process alive | HIGH |
| Process terminated | Read final state/restart | CRITICAL |
| ERROR signal in output | Analyze and recover | HIGH |
| STUCK signal in output | Attempt unstick | MEDIUM |
| Resource exceeded | Kill and restart | MEDIUM |

## State Updates Required
```json
{
  "monitoring_state": {
    "last_check": "timestamp",
    "iterations": 10,
    "process_summary": {
      "running": 2,
      "completed": 1,
      "failed": 0
    },
    "sub_details": {
      "uuid-1": {
        "shell_id": "bash_abc123",
        "last_output": "timestamp",
        "status": "RUNNING",
        "progress": "75%"
      }
    },
    "pending_actions": [
      {
        "sub_id": "uuid-3",
        "action": "RESTART",
        "reason": "no_output_timeout"
      }
    ]
  }
}
```

## Performance Considerations
- Monitor check interval: 15 seconds
- No-output timeout: 5 minutes
- Progress stuck timeout: 10 minutes
- Maximum monitoring duration: 2 hours
- Checkpoint frequency: Every minute

## Common Issues
1. **Buffer overflow** - Very verbose sub-orchestrators may fill output buffer
2. **Process zombies** - Terminated processes not cleaned up properly
3. **Output parsing** - Ensure consistent signal format
4. **Race conditions** - Process may complete between checks

## Example Output Parsing
```python
def parse_sub_output(output):
    """Parse standardized sub-orchestrator output"""
    signals = {
        'started': False,
        'progress': None,
        'complete': False,
        'failed': False,
        'error': None
    }

    for line in output.split('\n'):
        if 'SUB_ORCHESTRATOR_STARTED' in line:
            signals['started'] = True
        elif 'PROGRESS:' in line:
            # Extract progress percentage
            signals['progress'] = extract_percentage(line)
        elif 'SUB_ORCHESTRATOR_COMPLETE' in line:
            signals['complete'] = True
        elif 'SUB_ORCHESTRATOR_FAILED' in line:
            signals['failed'] = True
        elif 'ERROR:' in line:
            signals['error'] = line.split('ERROR:')[1].strip()

    return signals
```

## Notes
- Monitor multiple subs concurrently via shell_id
- No file polling needed - direct process output
- Clean process termination with KillShell if needed
- Preserve all output for debugging