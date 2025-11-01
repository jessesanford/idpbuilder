# MASTER_SPAWN_SUB_ORCHESTRATOR State Rules

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
Spawn a sub-orchestrator to handle a specific sub-state machine (FIX_CASCADE, INTEGRATE_WAVE_EFFORTS, SPLIT_COORDINATION) using background bash execution for non-blocking operation.

## Entry Criteria
- Decision made to delegate work to sub-orchestrator
- Parameters prepared for sub-orchestration
- No conflicting sub-orchestrators running for same scope

## Actions Required

### 1. Validate Spawn Conditions
```bash
# Check resource availability
# Verify no conflicts with existing sub-orchestrators
# Confirm state machine allows this spawn
```

### 2. Prepare Execution Environment
```bash
# Create working directory for sub-orchestrator
mkdir -p "/tmp/sub-orchestrators/$SUB_ID"

# Initialize tracking files (output and state only - no heartbeat needed)
touch "$OUTPUT_FILE"
touch "$STATE_FILE"
```

### 3. Generate Parameter File
```json
{
  "sub_orchestrator_type": "FIX_CASCADE|INTEGRATE_WAVE_EFFORTS|SPLIT_COORDINATION",
  "unique_id": "generated-uuid",
  "master_state_file": "orchestrator-state-v3.json",
  "input_parameters": {
    // Type-specific parameters
  },
  "output_location": "path/to/output.json",
  "state_file": "path/to/sub-state.json",
  "max_duration_seconds": 3600
}
```

### 4. Spawn Sub-Orchestrator Using Background Bash
```python
# Use Bash tool with run_in_background=true for non-blocking execution
result = Bash(
  command=f"claude -p '{CLAUDE_PROJECT_DIR}' --command '/sub-orchestrate-{type}' --params 'file=/tmp/params-{SUB_ID}.json'",
  run_in_background=True,
  description="Spawn sub-orchestrator in background"
)

# Extract shell_id from result for monitoring
SHELL_ID = result['shell_id']

# Extract PID from initial output if available
# Sub-orchestrator should echo its PID as first output line
```

### 5. Update Master State with Shell Tracking
```json
{
  "active_sub_orchestrators": [{
    "id": "uuid",
    "type": "FIX_CASCADE",
    "shell_id": "bash_abc123",
    "command": "claude -p . --command '/sub-orchestrate-fix'",
    "spawned_at": "timestamp",
    "status": "RUNNING",
    "parameter_file": "/tmp/params-uuid.json",
    "state_file": "/tmp/sub-orch/sub-state.json",
    "output_file": "/tmp/sub-orch/output.json",
    "last_output_check": "timestamp"
  }]
}
```

## Exit Criteria
- Sub-orchestrator successfully spawned in background
- Shell ID captured for monitoring
- Master state updated with tracking info
- Background process confirmed running

## Success Transition
→ `MASTER_MONITOR_SUB_ORCHESTRATORS`

## Failure Transitions
- Spawn failure → `ERROR_RECOVERY`
- Resource exhausted → `RESOURCE_WAIT`
- Invalid parameters → `PLANNING`

## Rules Applied
- R377: Master-Sub Communication Protocol
- R378: Sub-Orchestrator Lifecycle Management (Background Processes)
- R313: Mandatory Stop After Spawn
- R206: State Machine Validation

## Required State Updates
```json
{
  "spawn_result": {
    "success": true,
    "sub_id": "uuid",
    "shell_id": "bash_abc123",
    "type": "FIX_CASCADE",
    "spawned_at": "timestamp"
  }
}
```

## Monitoring Requirements
- Use BashOutput tool to check initial output
- Verify process started successfully (check for error output)
- Monitor for "SUB_ORCHESTRATOR_STARTED" signal in output
- Set up monitoring schedule using shell_id

## Common Issues
1. **Command not found** - Ensure claude command is in PATH
2. **Permission denied** - Check file/directory permissions
3. **Resource limits** - Verify system has capacity
4. **Parameter validation** - Ensure all required params present

## Example Implementation
```python
# Full spawn sequence using background bash
def spawn_fix_cascade_sub(fix_id, branches):
    import uuid
    import json

    # Generate unique ID
    sub_id = str(uuid.uuid4())

    # Prepare parameters
    params = {
        "sub_orchestrator_type": "FIX_CASCADE",
        "unique_id": sub_id,
        "input_parameters": {
            "fix_id": fix_id,
            "branches_to_fix": branches
        },
        "output_location": f"/tmp/sub-orchestrators/{sub_id}/output.json",
        "state_file": f"/tmp/sub-orchestrators/{sub_id}/state.json"
    }

    param_file = f"/tmp/params-{sub_id}.json"
    with open(param_file, 'w') as f:
        json.dump(params, f)

    # Spawn sub-orchestrator in background
    result = Bash(
        command=f"claude -p '{CLAUDE_PROJECT_DIR}' --command '/sub-orchestrate-fix' --params 'file={param_file}'",
        run_in_background=True,
        description="Spawn fix cascade sub-orchestrator"
    )

    # Update master state with shell_id
    update_master_state({
        "sub_id": sub_id,
        "shell_id": result['shell_id'],
        "type": "FIX_CASCADE"
    })

    return result['shell_id']
```

## Process Management
```python
# Check if sub-orchestrator is still running
def check_sub_process(shell_id):
    output = BashOutput(bash_id=shell_id)

    # Check for completion signal
    if "SUB_ORCHESTRATOR_COMPLETE" in output:
        return "COMPLETED"
    elif "SUB_ORCHESTRATOR_FAILED" in output:
        return "FAILED"
    else:
        return "RUNNING"
```

## Notes
- MUST stop after spawning (R313)
- Sub-orchestrator runs independently in background
- Monitor via BashOutput tool, not file polling
- Clean process lifecycle with shell_id tracking

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

