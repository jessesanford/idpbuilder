# 🚨🚨🚨 BLOCKING RULE R378: Sub-Orchestrator Lifecycle Management (Background Processes) 🚨🚨🚨

## Rule Definition
**ID**: R378
**Category**: ORCHESTRATION
**Criticality**: BLOCKING
**Created**: 2025-01-21
**Updated**: 2025-01-21 (Background Process Architecture)
**Dependencies**: R377 (Communication Protocol), R287 (TODO Persistence), R313 (Stop After Spawn)

## Description
Governs the complete lifecycle of sub-orchestrators from spawn to termination using background bash processes. Utilizes the Bash tool's `run_in_background` capability for non-blocking execution and the BashOutput tool for real-time monitoring. Ensures proper resource management and state consistency without file-based heartbeats.

## 🔴 LIFECYCLE PHASES

### 1. PRE-SPAWN PHASE
Master orchestrator MUST complete before spawning:

```python
# 1. Validate pre-conditions
def validate_can_spawn_sub_orchestrator():
    # Check no conflicting sub-orchestrators running
    # Verify resources available
    # Confirm state machine allows spawn
    pass

# 2. Prepare execution environment
def prepare_sub_environment(sub_id):
    # Create working directory
    os.makedirs(f"/tmp/sub-orchestrators/{sub_id}", exist_ok=True)

    # Initialize output and state files only (no heartbeat)
    Path(f"/tmp/sub-orchestrators/{sub_id}/output.json").touch()
    Path(f"/tmp/sub-orchestrators/{sub_id}/state.json").touch()

# 3. Generate unique identifiers
sub_id = str(uuid.uuid4())
process_tag = f"sub-orch-{type}-{sub_id}"
```

### 2. SPAWN PHASE (Using Background Bash)

```python
# Master state before spawn
{
  "pending_sub_orchestration": {
    "type": "FIX_CASCADE",
    "id": "uuid-123",
    "reason": "Fix critical bug in 3 branches",
    "prepared_at": "timestamp"
  }
}

# Spawn with background execution
def spawn_sub_orchestrator(type, params):
    sub_id = str(uuid.uuid4())
    param_file = f"/tmp/params-{sub_id}.json"

    # Write parameters
    with open(param_file, 'w') as f:
        json.dump({
            "sub_orchestrator_type": type,
            "unique_id": sub_id,
            "input_parameters": params,
            "output_location": f"/tmp/sub-orchestrators/{sub_id}/output.json",
            "state_file": f"/tmp/sub-orchestrators/{sub_id}/state.json"
        }, f)

    # Launch sub-orchestrator in background
    result = Bash(
        command=f"claude -p '{CLAUDE_PROJECT_DIR}' --command '/sub-orchestrate-{type}' --params 'file={param_file}'",
        run_in_background=True,
        description=f"Spawn {type} sub-orchestrator"
    )

    # Record in master state with shell_id
    record_sub_spawn(sub_id, result['shell_id'], type)

    return result['shell_id']

# Master state after spawn
{
  "active_sub_orchestrators": [{
    "id": "uuid-123",
    "type": "FIX_CASCADE",
    "shell_id": "bash_abc123",
    "spawned_at": "timestamp",
    "status": "RUNNING"
  }]
}
```

### 3. INITIALIZATION PHASE (Sub-Orchestrator)

Sub-orchestrator MUST output these signals on startup:

```bash
#!/bin/bash
# SUB-ORCHESTRATOR INITIALIZATION PROTOCOL

# 1. Echo PID as first output
echo "SUB_ORCHESTRATOR_PID: $$"
echo "SUB_ORCHESTRATOR_STARTED: $(date -Iseconds)"

# 2. Parse parameters
PARAM_FILE="${1:-/tmp/params.json}"
PARAMS=$(cat "$PARAM_FILE")
SUB_ID=$(echo "$PARAMS" | jq -r '.unique_id')
SUB_TYPE=$(echo "$PARAMS" | jq -r '.sub_orchestrator_type')

# 3. Set up completion handler
trap 'handle_completion' EXIT

handle_completion() {
  if [ $? -eq 0 ]; then
    echo "SUB_PROCESS_COMPLETE: Success"
  else
    echo "SUB_PROCESS_FAILED: Exit code $?"
  fi
}

# 4. Validate environment
echo "SUB_PROCESS_PROGRESS: 0% - Validating environment"
validate_sub_environment

# 5. Load state machine
echo "SUB_PROCESS_PROGRESS: 10% - Loading state machine"
load_sub_state_machine
```

### 4. EXECUTION PHASE (With Progress Signals)

```bash
# Sub-orchestrator main loop with output signals
execute_sub_orchestration() {
  local CURRENT_STATE="START"
  local ITERATION=0
  local PROGRESS=0

  while [[ "$CURRENT_STATE" != "COMPLETE" ]]; do
    # Output progress signal
    echo "SUB_ORCHESTRATOR_PROGRESS: ${PROGRESS}% - State: $CURRENT_STATE"

    # Execute state action
    case "$CURRENT_STATE" in
      "START")
        echo "Processing initialization tasks..."
        # Perform start actions
        NEXT_STATE="PROCESS"
        PROGRESS=25
        ;;
      "PROCESS")
        echo "Main processing phase..."
        # Main processing
        NEXT_STATE="VALIDATE"
        PROGRESS=75
        ;;
      "VALIDATE")
        echo "Validating results..."
        # Validation
        NEXT_STATE="COMPLETE"
        PROGRESS=90
        ;;
    esac

    # Write state file for recovery
    write_state_file "$CURRENT_STATE" "$NEXT_STATE" "$PROGRESS"

    # Transition
    CURRENT_STATE="$NEXT_STATE"
    ITERATION=$((ITERATION + 1))

    # Safety check
    if [[ $ITERATION -gt 100 ]]; then
      echo "SUB_ORCHESTRATOR_ERROR: Infinite loop detected"
      exit 1
    fi
  done

  echo "SUB_ORCHESTRATOR_PROGRESS: 100% - Completed"
}
```

### 5. MONITORING PHASE (Master Using BashOutput)

```python
def monitor_sub_orchestrators():
    """Monitor all active sub-orchestrators via BashOutput"""
    iterations = 0
    max_iterations = 240  # 1 hour

    while has_active_subs() and iterations < max_iterations:
        for sub in get_active_subs():
            # Get output from background process
            output = BashOutput(bash_id=sub['shell_id'])

            if output:
                # Parse output for signals
                if "SUB_ORCHESTRATOR_PROGRESS:" in output:
                    progress = extract_progress(output)
                    update_sub_progress(sub['id'], progress)

                if "SUB_ORCHESTRATOR_ERROR:" in output:
                    error = extract_error(output)
                    handle_sub_error(sub['id'], error)

                if "SUB_ORCHESTRATOR_COMPLETE" in output:
                    handle_completed_sub(sub['id'])

                if "SUB_ORCHESTRATOR_FAILED" in output:
                    handle_failed_sub(sub['id'])

                # Update last output time
                update_last_output(sub['id'])

            # Check for no output timeout
            elif time_since_last_output(sub['id']) > 300:  # 5 minutes
                handle_stuck_sub(sub['id'])

        time.sleep(15)  # Check every 15 seconds
        iterations += 1
```

### 6. COMPLETION PHASE

```python
# Sub-orchestrator completion (bash)
complete_sub_orchestration() {
  # 1. Finalize results
  cat > "$OUTPUT_FILE" <<EOF
{
  "status": "SUCCESS",
  "completed_at": "$(date -Iseconds)",
  "results": $RESULTS,
  "next_action": "CONTINUE"
}
EOF

  # 2. Output completion signal
  echo "SUB_ORCHESTRATOR_COMPLETE: Success"

  # 3. Clean up temp files
  cleanup_temp_files

  # 4. Exit successfully
  exit 0
}

# Master handling completion (python)
def handle_completed_sub(sub_id):
    # 1. Read results from output file
    output_file = f"/tmp/sub-orchestrators/{sub_id}/output.json"
    with open(output_file, 'r') as f:
        results = json.load(f)

    # 2. Update master state
    mark_sub_complete(sub_id, results)

    # 3. Determine next action
    next_action = results.get('next_action', 'CONTINUE')
    if next_action == 'CONTINUE':
        proceed_to_next_state()
    elif next_action == 'RETRY':
        schedule_retry(sub_id)
    elif next_action == 'ESCALATE':
        escalate_to_error_recovery()

    # 4. Clean up shell if still exists
    try:
        KillShell(shell_id=get_shell_id(sub_id))
    except:
        pass  # Shell already terminated
```

### 7. CLEANUP PHASE

```python
def cleanup_sub_orchestrator(sub_id):
    # 1. Archive outputs
    archive_dir = "archives/sub-orchestrators"
    os.makedirs(archive_dir, exist_ok=True)

    # Move output files
    shutil.move(
        f"/tmp/sub-orchestrators/{sub_id}/output.json",
        f"{archive_dir}/{sub_id}-output.json"
    )

    # 2. Remove temp files
    shutil.rmtree(f"/tmp/sub-orchestrators/{sub_id}", ignore_errors=True)
    os.remove(f"/tmp/params-{sub_id}.json")

    # 3. Update history
    add_to_sub_history(sub_id)

    # 4. Kill shell if still running
    shell_id = get_shell_id(sub_id)
    if shell_id:
        try:
            KillShell(shell_id=shell_id)
        except:
            pass
```

## 🛑 STATE TRACKING REQUIREMENTS

### Master State Tracking
```json
{
  "active_sub_orchestrators": [
    {
      "id": "uuid-123",
      "type": "FIX_CASCADE",
      "shell_id": "bash_abc123",
      "spawned_at": "2025-01-21T10:00:00Z",
      "last_output": "2025-01-21T10:15:30Z",
      "status": "RUNNING",
      "progress": "60%",
      "current_state": "BACKPORTING"
    }
  ],
  "sub_orchestrator_history": [
    {
      "id": "uuid-122",
      "type": "INTEGRATION",
      "started": "2025-01-21T09:00:00Z",
      "completed": "2025-01-21T09:45:00Z",
      "result": "SUCCESS",
      "duration_seconds": 2700
    }
  ],
  "pending_sub_orchestrations": [
    {
      "type": "SPLIT_COORDINATION",
      "reason": "Effort E1.1 exceeds size limit",
      "scheduled_for": "AFTER_CURRENT_COMPLETE"
    }
  ]
}
```

## ⚠️ RESOURCE LIMITS

### Per Sub-Orchestrator:
- Maximum runtime: 2 hours
- Maximum retries: 3
- No-output timeout: 5 minutes
- State file size: 10MB
- State save frequency: 30 seconds

### System-wide:
- Maximum concurrent sub-orchestrators: 5
- Maximum queued sub-orchestrations: 10
- Total disk usage: 1GB
- Archive retention: 7 days

## 🔴 ERROR CONDITIONS

### MUST Handle:
1. **Sub-orchestrator crash** - Detect via BashOutput error signals
2. **Timeout** - Kill shell after max runtime
3. **No output** - Handle stuck processes after 5 minutes
4. **Resource exhaustion** - Prevent spawn if limits exceeded
5. **Orphaned shells** - Clean up with KillShell

### Recovery Strategies:
```python
# State-based recovery
def recover_from_state(sub_id):
    state_file = f"/tmp/sub-orchestrators/{sub_id}/state.json"
    if os.path.exists(state_file):
        with open(state_file, 'r') as f:
            last_state = json.load(f)

        # Restart from last known state
        spawn_recovery_sub(sub_id, last_state)

# Full restart
def restart_sub_orchestrator(sub_id, original_params):
    # Kill existing shell if any
    shell_id = get_shell_id(sub_id)
    if shell_id:
        KillShell(shell_id=shell_id)

    # Clean up failed attempt
    cleanup_failed_sub(sub_id)

    # Respawn with same parameters
    spawn_sub_orchestrator(original_params['type'], original_params)
```

## 🚨 VIOLATIONS

### BLOCKING (-100% penalty):
1. Spawning without capturing shell_id
2. No BashOutput monitoring
3. Missing completion signals
4. Shell leaks (not killed on completion)
5. Lost sub-orchestrator results

### WARNING (-25% penalty):
1. Delayed progress updates (>45s)
2. Missing state file saves
3. Incomplete cleanup
4. No retry on transient failure
5. Poor signal formatting

## 📋 Implementation Requirements

- [ ] Pre-spawn validation complete
- [ ] Background bash spawn with shell_id tracking
- [ ] BashOutput monitoring active
- [ ] Signal-based completion detection
- [ ] KillShell cleanup implemented
- [ ] Resource limits enforced
- [ ] Recovery mechanisms tested
- [ ] Archives properly maintained

## 🔗 Related Rules
- R377: Communication Protocol (via output signals)
- R379: Process-Based Monitoring
- R287: TODO Persistence
- R313: Stop After Spawn

## 🔴 Key Changes in Background Process Architecture

1. **No Heartbeat Files** - Use BashOutput for monitoring
2. **Shell ID Tracking** - Track shell_id instead of PID
3. **Signal-Based Communication** - Standard output signals
4. **Clean Process Management** - KillShell for termination
5. **Non-Blocking Execution** - Master never waits for sub

---

**Critical**: Every sub-orchestrator spawn MUST use `run_in_background=true` and track the shell_id. Monitor via BashOutput, not file polling!