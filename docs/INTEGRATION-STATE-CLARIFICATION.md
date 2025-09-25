# Integration State Clarification

## Purpose
This document clarifies the correct responsibilities of WAVE_COMPLETE vs INTEGRATION states to prevent confusion about where integration branch creation should occur.

## CORRECT State Responsibilities

### WAVE_COMPLETE State
**Purpose**: Validate wave completion and prepare for integration

**MUST DO**:
- Verify all efforts in the wave are marked as completed
- Confirm all code reviews have passed
- Check no pending size violations exist
- Update state file to mark wave as complete
- Set pending_integration_wave metadata

**MUST NOT DO**:
- Create integration branches
- Set up integration infrastructure
- Start any integration work
- Modify git branches

### INTEGRATION State
**Purpose**: Set up integration infrastructure and prepare for merging

**MUST DO**:
- Create the wave-${WAVE_NUM}-integration branch
- Check out from main and pull latest
- Update state file with integration_in_progress = true
- Set integration_branch name in state
- Prepare for spawning merge plan creation

**MUST NOT DO**:
- Re-validate wave completion (already done in WAVE_COMPLETE)
- Attempt to create branch if it already exists

## State Transition Flow

```
MONITOR (all reviews pass)
    ↓
WAVE_COMPLETE (validate completion)
    ↓
INTEGRATION (create integration branch)
    ↓
SPAWN_CODE_REVIEWER_MERGE_PLAN
```

## Common Mistakes to Avoid

1. **DON'T** put integration branch creation in WAVE_COMPLETE
2. **DON'T** duplicate branch creation logic in multiple states
3. **DON'T** skip WAVE_COMPLETE validation before integration
4. **DON'T** try to create the same branch twice

## Validation Code

### In WAVE_COMPLETE:
```bash
# Validate but don't create branches
WAVE_NUM=$(jq '.current_wave' orchestrator-state.json)
INCOMPLETE_COUNT=$(jq '.efforts_in_progress | length' orchestrator-state.json)

if [ "$INCOMPLETE_COUNT" -gt 0 ]; then
    echo "ERROR: Still have ${INCOMPLETE_COUNT} efforts in progress"
    exit 1
fi

# Just prepare metadata
jq ".pending_integration_wave = ${WAVE_NUM}" orchestrator-state.json
```

### In INTEGRATION:
```bash
# Create the integration branch HERE
WAVE_NUM=$(jq '.pending_integration_wave' orchestrator-state.json)

git checkout main
git pull origin main
git checkout -b "wave-${WAVE_NUM}-integration"

jq '.integration_in_progress = true' orchestrator-state.json
jq ".integration_branch = \"wave-${WAVE_NUM}-integration\"" orchestrator-state.json
```

## Reference
- State Machine: SOFTWARE-FACTORY-STATE-MACHINE.md (lines 327-328)
- Rule R104: Integration branch creation (if exists)
- Rule R234: Mandatory state traversal

## Last Updated
2025-08-30

## Verified By
Software Factory Manager Agent