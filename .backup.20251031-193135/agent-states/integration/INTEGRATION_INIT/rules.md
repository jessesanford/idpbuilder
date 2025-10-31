# INTEGRATE_WAVE_EFFORTS_INIT State Rules

## State Purpose
Initialize the integration sub-state machine and prepare for integration operations.

## Entry Conditions
- Entered from main state machine (WAVE_COMPLETE, COMPLETE_PHASE, or PROJECT_INTEGRATE_WAVE_EFFORTS)
- Parent state has provided integration requirements
- Integration state file needs to be created

## Required Actions

### 1. Create Integration State File
```bash
# Create state file with unique identifier
INTEGRATE_WAVE_EFFORTS_TYPE="${1:-WAVE}"  # WAVE, PHASE, or PROJECT
IDENTIFIER="${2:-wave1}"
STATE_FILE="integration-${INTEGRATE_WAVE_EFFORTS_TYPE}-${IDENTIFIER}-state.json"
```

### 2. Initialize State Structure
```json
{
  "sub_state_type": "INTEGRATE_WAVE_EFFORTS",
  "current_state": "INTEGRATE_WAVE_EFFORTS_INIT",
  "integration_config": {
    "type": "WAVE",
    "identifier": "wave1",
    "started_at": "2025-01-21T10:00:00Z",
    "attempt": 1,
    "max_attempts": 10
  },
  "parent_state_machine": {
    "state_file": "orchestrator-state-v3.json",
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
- Transition to INTEGRATE_WAVE_EFFORTS_LOAD_REQUIREMENTS

## Validation Rules

### R375 - Dual State File Pattern
- Main state file remains `orchestrator-state-v3.json`
- Integration state in `integration-[type]-[id]-state.json`
- No mixing of state data between files

### R327 - Stale Integration Prevention
- Check for stale integration branches
- Mark for deletion if found
- Prevent reuse of old integrations

## Error Handling
- If state file creation fails → INTEGRATE_WAVE_EFFORTS_ERROR
- If requirements missing → INTEGRATE_WAVE_EFFORTS_ABORT
- If max attempts exceeded → INTEGRATE_WAVE_EFFORTS_ABORT

## Logging Requirements
```bash
echo "[INTEGRATE_WAVE_EFFORTS_INIT] Starting integration sub-state machine"
echo "[INTEGRATE_WAVE_EFFORTS_INIT] Type: ${INTEGRATE_WAVE_EFFORTS_TYPE}, Identifier: ${IDENTIFIER}"
echo "[INTEGRATE_WAVE_EFFORTS_INIT] Attempt: ${ATTEMPT_NUMBER}/${MAX_ATTEMPTS}"
echo "[INTEGRATE_WAVE_EFFORTS_INIT] State file: ${STATE_FILE}"
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

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

