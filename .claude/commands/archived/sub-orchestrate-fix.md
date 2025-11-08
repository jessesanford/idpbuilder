# /sub-orchestrate-fix Command

## Purpose
Entry point for FIX_CASCADE sub-orchestrator. Handles fix propagation across multiple branches independently from main orchestrator using output signals for communication.

## Usage
```bash
claude -p "$CLAUDE_PROJECT_DIR" \
  --command "/sub-orchestrate-fix" \
  --params "file=/tmp/params-uuid.json" \
  --state "sub-state-fix-uuid.json"
```

## Parameter File Format
```json
{
  "sub_orchestrator_type": "FIX_CASCADE",
  "unique_id": "uuid-v4",
  "master_state_file": "orchestrator-state-v3.json",
  "input_parameters": {
    "fix_id": "fix-bug-123",
    "branches_to_fix": ["branch1", "branch2", "branch3"],
    "issue_description": "Critical bug in authentication module",
    "validation_requirements": {
      "tests_must_pass": true,
      "linting_required": true,
      "review_required": false
    },
    "source_fix": {
      "branch": "fix/original-branch",
      "commit": "abc123def456"
    }
  },
  "output_location": "/tmp/sub-orch-uuid/output.json",
  "state_file": "/tmp/sub-orch-uuid/state.json",
  "max_duration_seconds": 3600
}
```

## Execution Flow

### 1. INITIALIZATION WITH SIGNALS
```bash
#!/bin/bash

# Output startup signals
echo "SUB_ORCHESTRATOR_PID: $$"
echo "SUB_ORCHESTRATOR_STARTED: $(date -Iseconds)"
echo "SUB_ORCHESTRATOR_TYPE: FIX_CASCADE"

# Parse parameters
PARAM_FILE="${1:-/tmp/params.json}"
PARAMS=$(cat "$PARAM_FILE")
SUB_ID=$(echo "$PARAMS" | jq -r '.unique_id')

# Set up environment
export SUB_ORCHESTRATOR_ID="$SUB_ID"
export SUB_TYPE="FIX_CASCADE"
export OUTPUT_FILE=$(echo "$PARAMS" | jq -r '.output_location')
export STATE_FILE=$(echo "$PARAMS" | jq -r '.state_file')

# Initialize progress
echo "SUB_ORCHESTRATOR_PROGRESS: 0% - Initializing FIX_CASCADE"

# Set up completion handler
trap 'handle_completion' EXIT

handle_completion() {
  if [ $? -eq 0 ]; then
    echo "SUB_ORCHESTRATOR_COMPLETE: Success"
  else
    echo "SUB_ORCHESTRATOR_FAILED: Exit code $?"
  fi
}
```

### 2. MAIN EXECUTION WITH PROGRESS SIGNALS
```bash
# Load fix cascade state machine
echo "SUB_ORCHESTRATOR_PROGRESS: 5% - Loading state machine"
load_state_machine "FIX_CASCADE"

# Execute states with progress reporting
execute_sub_orchestration() {
  case "$CURRENT_STATE" in
    "FIX_CASCADE_INIT")
      echo "SUB_ORCHESTRATOR_PROGRESS: 10% - Validating fix requirements"
      validate_fix_requirements
      ;;
    "FIX_CASCADE_SPAWN_SW_ENGINEERS")
      echo "SUB_ORCHESTRATOR_PROGRESS: 25% - Spawning fix agents for branches"
      spawn_fix_agents_for_branches
      ;;
    "FIX_CASCADE_MONITORING_EFFORT_FIXES")
      echo "SUB_ORCHESTRATOR_PROGRESS: 50% - Monitoring fix progress"
      monitor_fix_progress
      # Update progress based on branches done
      PERCENT=$((50 + (BRANCHES_DONE * 30 / BRANCHES_TOTAL)))
      echo "SUB_ORCHESTRATOR_PROGRESS: ${PERCENT}% - Fixed ${BRANCHES_DONE}/${BRANCHES_TOTAL} branches"
      ;;
    "FIX_CASCADE_VALIDATE")
      echo "SUB_ORCHESTRATOR_PROGRESS: 80% - Validating all fixes"
      validate_all_fixes
      ;;
    "FIX_CASCADE_COMPLETE")
      echo "SUB_ORCHESTRATOR_PROGRESS: 95% - Writing completion output"
      write_completion_output
      echo "SUB_ORCHESTRATOR_PROGRESS: 100% - FIX_CASCADE complete"
      ;;
  esac
}
```

### 3. STATE PERSISTENCE (Replaces Checkpointing)
```bash
# Save state for recovery
save_state() {
  local STATE="$1"
  local DATA="$2"

  cat > "$STATE_FILE" <<EOF
{
  "timestamp": "$(date -Iseconds)",
  "current_state": "$STATE",
  "data": $DATA,
  "progress": $PROGRESS,
  "can_resume": true,
  "branches_processed": $BRANCHES_DONE,
  "branches_total": $BRANCHES_TOTAL
}
EOF

  echo "State saved: $STATE (${PROGRESS}%)"
}

# Main loop with state saves
while [[ "$CURRENT_STATE" != "COMPLETE" ]]; do
  # Execute current state
  execute_state "$CURRENT_STATE"

  # Save state after each transition
  save_state "$CURRENT_STATE" "$STATE_DATA"

  # Report any errors via signals
  if [[ -n "$ERROR" ]]; then
    echo "SUB_ORCHESTRATOR_ERROR: $ERROR"
  fi

  # Transition to next state
  CURRENT_STATE=$(get_next_state "$CURRENT_STATE")
done
```

### 4. OUTPUT GENERATION
```bash
# Write final output
write_output() {
  cat > "$OUTPUT_FILE" <<EOF
{
  "status": "$FINAL_STATUS",
  "completion_timestamp": "$(date -Iseconds)",
  "duration_seconds": $DURATION,
  "results": {
    "branches_fixed": $BRANCHES_FIXED,
    "test_results": $TEST_RESULTS,
    "commits_created": $COMMITS,
    "validation_passed": $VALIDATION
  },
  "errors": $ERRORS,
  "next_action": "$NEXT_ACTION",
  "state_transitions_executed": $TRANSITIONS,
  "artifacts_created": $ARTIFACTS
}
EOF

  echo "Output written to $OUTPUT_FILE"
}
```

## State Machine States
- `FIX_CASCADE_INIT` - Initialize and validate
- `FIX_CASCADE_SPAWN_SW_ENGINEERS` - Spawn SW engineers for fixes
- `FIX_CASCADE_MONITORING_EFFORT_FIXES` - Monitor fix progress
- `FIX_CASCADE_VALIDATE` - Validate all fixes
- `FIX_CASCADE_BACKPORT` - Backport fixes if needed
- `FIX_CASCADE_TEST` - Run comprehensive tests
- `FIX_CASCADE_COMPLETE` - Finalize and report

## Recovery Support
```bash
# Resume from state file if exists
if [[ -f "$STATE_FILE" ]]; then
  echo "SUB_ORCHESTRATOR_PROGRESS: Resuming from saved state"
  RESUME_STATE=$(jq -r '.current_state' "$STATE_FILE")
  RESUME_DATA=$(jq '.data' "$STATE_FILE")
  BRANCHES_DONE=$(jq -r '.branches_processed' "$STATE_FILE")

  echo "Resuming from state: $RESUME_STATE (${BRANCHES_DONE} branches done)"
  resume_from_state "$RESUME_STATE" "$RESUME_DATA"
fi
```

## Error Handling with Signals
```bash
# Trap errors and signal via output
trap 'handle_error $?' ERR

handle_error() {
  local EXIT_CODE="$1"

  # Output error signal
  echo "SUB_ORCHESTRATOR_ERROR: Process failed with exit code $EXIT_CODE"
  echo "SUB_ORCHESTRATOR_ERROR: $ERROR_MSG"

  # Save error state
  cat > "$STATE_FILE" <<EOF
{
  "status": "FAILED",
  "error_code": $EXIT_CODE,
  "error_message": "$ERROR_MSG",
  "last_state": "$CURRENT_STATE",
  "timestamp": "$(date -Iseconds)"
}
EOF

  # Write failure output
  write_failure_output

  # Signal failure
  echo "SUB_ORCHESTRATOR_FAILED: $ERROR_MSG"

  exit $EXIT_CODE
}
```

## Monitoring Integration
- Progress signals output continuously
- State saved at each transition
- Error signals on failures
- Completion signals on success
- No file polling needed - master monitors via BashOutput

## Standard Output Signals
```bash
# Standard signals the master orchestrator expects:
echo "SUB_ORCHESTRATOR_PID: $$"                    # Process ID
echo "SUB_ORCHESTRATOR_STARTED: <timestamp>"       # Start time
echo "SUB_ORCHESTRATOR_PROGRESS: <N>% - <desc>"    # Progress updates
echo "SUB_ORCHESTRATOR_ERROR: <message>"           # Error conditions
echo "SUB_ORCHESTRATOR_STUCK: <reason>"            # Stuck detection
echo "SUB_ORCHESTRATOR_COMPLETE: Success"          # Successful completion
echo "SUB_ORCHESTRATOR_FAILED: <reason>"           # Failed completion
```

## Rules Applied
- R377: Communication Protocol (Output Signals)
- R378: Lifecycle Management (Background Processes)
- R379: Process-Based Monitoring
- R375: Fix State File Management
- R376: Fix Cascade Quality Gates

## Success Criteria
- All branches successfully fixed
- All tests passing
- Validation requirements met
- Output file written
- COMPLETE signal sent
- Clean exit with status 0

## Common Issues
1. Branch conflicts during backport
2. Test failures on specific branches
3. Resource exhaustion with many branches
4. Timeout on large codebases

## Example Launch (Background)
```python
# Master orchestrator spawns this in background
import json

# Prepare parameters
params = {
    "sub_orchestrator_type": "FIX_CASCADE",
    "unique_id": "fix-123-uuid",
    "input_parameters": {
        "fix_id": "bug-123",
        "branches_to_fix": ["release/1.0", "release/1.1"],
        "issue_description": "Memory leak in cache handler"
    },
    "output_location": "/tmp/sub-orch/fix-123/output.json",
    "state_file": "/tmp/sub-orch/fix-123/state.json"
}

# Write params file
param_file = "/tmp/fix-params.json"
with open(param_file, 'w') as f:
    json.dump(params, f)

# Launch in background
result = Bash(
    command=f"claude -p . --command '/sub-orchestrate-fix' --params 'file={param_file}'",
    run_in_background=True
)

# Monitor via shell_id
shell_id = result['shell_id']
```