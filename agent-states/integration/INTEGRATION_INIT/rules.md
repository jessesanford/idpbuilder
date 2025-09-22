# INTEGRATION_INIT State Rules

## State Purpose
Initialize the integration sub-state machine and prepare for integration operations.

## Entry Conditions
- Entered from main state machine (WAVE_COMPLETE, PHASE_COMPLETE, or PROJECT_INTEGRATION)
- Parent state has provided integration requirements
- Integration state file needs to be created

## Required Actions

### 1. Create Integration State File
```bash
# Create state file with unique identifier
INTEGRATION_TYPE="${1:-WAVE}"  # WAVE, PHASE, or PROJECT
IDENTIFIER="${2:-wave1}"
STATE_FILE="integration-${INTEGRATION_TYPE}-${IDENTIFIER}-state.json"
```

### 2. Initialize State Structure
```json
{
  "sub_state_type": "INTEGRATION",
  "current_state": "INTEGRATION_INIT",
  "integration_config": {
    "type": "WAVE",
    "identifier": "wave1",
    "started_at": "2025-01-21T10:00:00Z",
    "attempt": 1,
    "max_attempts": 10
  },
  "parent_state_machine": {
    "state_file": "orchestrator-state.json",
    "return_state": "WAVE_COMPLETE",
    "nested_level": 1
  },
  "cycle_tracking": {
    "current_attempt": 1,
    "cycle_history": []
  }
}
```

### 3. Load Integration Requirements
- Read requirements from parent state
- Identify branches to integrate
- Set validation requirements
- Determine base branch

### 4. Check for Previous Attempts
- Look for existing integration state files
- Load cycle history if re-entering
- Increment attempt counter if applicable

## Exit Conditions
- State file created successfully
- Requirements loaded
- Transition to INTEGRATION_LOAD_REQUIREMENTS

## Validation Rules

### R375 - Dual State File Pattern
- Main state file remains `orchestrator-state.json`
- Integration state in `integration-[type]-[id]-state.json`
- No mixing of state data between files

### R327 - Stale Integration Prevention
- Check for stale integration branches
- Mark for deletion if found
- Prevent reuse of old integrations

## Error Handling
- If state file creation fails → INTEGRATION_ERROR
- If requirements missing → INTEGRATION_ABORT
- If max attempts exceeded → INTEGRATION_ABORT

## Logging Requirements
```bash
echo "[INTEGRATION_INIT] Starting integration sub-state machine"
echo "[INTEGRATION_INIT] Type: ${INTEGRATION_TYPE}, Identifier: ${IDENTIFIER}"
echo "[INTEGRATION_INIT] Attempt: ${ATTEMPT_NUMBER}/${MAX_ATTEMPTS}"
echo "[INTEGRATION_INIT] State file: ${STATE_FILE}"
```

## Metrics to Track
- Entry timestamp
- Integration type distribution
- Re-entry frequency
- Previous attempt count

## Common Issues
1. **Missing Requirements**: Parent didn't provide branches list
2. **State File Conflicts**: Previous integration not cleaned up
3. **Max Attempts**: Already exceeded retry limit

## Success Criteria
✅ State file created and initialized
✅ Requirements loaded from parent
✅ Attempt tracking initialized
✅ Ready to proceed with integration