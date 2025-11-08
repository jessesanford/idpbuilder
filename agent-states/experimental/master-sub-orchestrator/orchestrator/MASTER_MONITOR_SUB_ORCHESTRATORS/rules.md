# MASTER_MONITOR_SUB_ORCHESTRATORS State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
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

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Update state file
jq '.state_machine.current_state = "NEXT_STATE"' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# 3. Save TODOs
save_todos "R322_CHECKPOINT"

# 4. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 5. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

