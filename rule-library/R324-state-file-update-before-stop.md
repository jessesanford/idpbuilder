# 🔴🔴🔴 RULE R324: STATE FILE UPDATE BEFORE STOP 🔴🔴🔴

## Rule Identifier
**Rule ID:** R324  
**Category:** State Machine Control  
**Criticality:** 🔴🔴🔴 SUPREME LAW  
**Introduced:** Version 2.0.324  
**Related:** R322, R206, R288

## Rule Statement

**THE ORCHESTRATOR MUST UPDATE current_state TO THE NEW STATE BEFORE STOPPING!**

Without this, the orchestrator will be stuck in infinite loops, repeatedly executing the same state work forever.

## The Critical Problem This Solves

### ❌ THE LOOP BUG (WHAT HAPPENS WITHOUT THIS RULE):
1. Orchestrator completes work in STATE_A
2. Orchestrator says "Next state: STATE_B" 
3. Orchestrator stops per R322 (but current_state still says STATE_A!)
4. User runs /continue-orchestrating
5. Orchestrator reads current_state: STATE_A
6. Orchestrator repeats STATE_A work again
7. INFINITE LOOP - Never progresses!

### ✅ THE CORRECT PATTERN:
1. Orchestrator completes work in STATE_A
2. Orchestrator updates current_state to STATE_B in file
3. Orchestrator commits and pushes the change
4. Orchestrator stops per R322
5. User runs /continue-orchestrating
6. Orchestrator reads current_state: STATE_B
7. Orchestrator continues from STATE_B correctly

## Mandatory Implementation Pattern

### 🔴🔴🔴 THIS EXACT SEQUENCE IS REQUIRED 🔴🔴🔴

```bash
# Step 1: Complete all work for current state
echo "✅ Completed all work for CURRENT_STATE"

# Step 2: UPDATE STATE FILE FIRST (BEFORE STOPPING!)
echo "📝 Updating state file for transition..."
yq -i '.current_state = "NEXT_STATE"' orchestrator-state.yaml
yq -i '.previous_state = "CURRENT_STATE"' orchestrator-state.yaml
yq -i ".transition_time = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.yaml

# Step 3: Commit and push the state change
git add orchestrator-state.yaml
git commit -m "state: transition from CURRENT_STATE to NEXT_STATE"
git push

# Step 4: NOW stop per R322
echo "🛑 STATE TRANSITION CHECKPOINT: CURRENT_STATE → NEXT_STATE"
echo "📊 State file updated to: NEXT_STATE"
echo "⏸️ STOPPED - Ready to continue in NEXT_STATE"
echo "When restarted, will continue from NEXT_STATE"
# EXIT HERE
```

## Common State Transitions That MUST Follow This Pattern

All of these MUST update current_state before stopping:

- `WAITING_FOR_MERGE_PLAN` → `SPAWN_INTEGRATION_AGENT`
- `SPAWN_INTEGRATION_AGENT` → `MONITORING_INTEGRATION`
- `MONITORING_INTEGRATION` → `INTEGRATION_FEEDBACK_REVIEW` or `WAVE_REVIEW`
- `WAVE_COMPLETE` → `INTEGRATION`
- `INTEGRATION` → `SPAWN_CODE_REVIEWER_MERGE_PLAN`
- `PHASE_COMPLETE` → `INIT` (next phase)
- `SPAWN_AGENTS` → `MONITOR`
- ALL OTHER STATE TRANSITIONS!

## Enforcement

### Success Criteria
- ✅ current_state updated BEFORE stop
- ✅ State file committed and pushed
- ✅ Clear message about new state
- ✅ No loops when restarted

### Failure Conditions
- ❌ Stop without updating current_state = INFINITE LOOP
- ❌ Update state after stopping = TOO LATE
- ❌ Forget to commit = State lost
- ❌ Continue after updating = R322 violation

## Examples

### ✅ CORRECT: Update state, then stop
```bash
# Complete SPAWN_AGENTS work
echo "All agents spawned successfully"

# UPDATE STATE FIRST!
yq -i '.current_state = "MONITOR"' orchestrator-state.yaml
yq -i '.previous_state = "SPAWN_AGENTS"' orchestrator-state.yaml
git add orchestrator-state.yaml
git commit -m "state: transition to MONITOR"
git push

# THEN STOP
echo "🛑 Stopping before MONITOR state"
echo "State updated to: MONITOR"
```

### ❌ WRONG: Stop without updating (CAUSES LOOP!)
```bash
# Complete SPAWN_AGENTS work
echo "All agents spawned successfully"

# WRONG - Stopping without updating state!
echo "🛑 Stopping before MONITOR state"
echo "Next state: MONITOR"  # Just saying it doesn't update the file!

# PROBLEM: current_state still says SPAWN_AGENTS
# When restarted, will repeat SPAWN_AGENTS forever!
```

### ❌ WRONG: Update after stopping (TOO LATE!)
```bash
# Complete SPAWN_AGENTS work
echo "All agents spawned successfully"

# WRONG ORDER!
echo "🛑 Stopping now"
exit 0

# This never executes!
yq -i '.current_state = "MONITOR"' orchestrator-state.yaml
```

## Debugging Loops

If the orchestrator is stuck in a loop:

1. **Check the state file**:
   ```bash
   grep current_state orchestrator-state.yaml
   ```

2. **Look for the pattern**:
   - Does it keep saying the same state?
   - Does the work keep repeating?

3. **Fix it**:
   ```bash
   # Manually update to the correct next state
   yq -i '.current_state = "CORRECT_NEXT_STATE"' orchestrator-state.yaml
   git add orchestrator-state.yaml
   git commit -m "fix: manual state correction to break loop"
   git push
   ```

## Integration with R322

This rule ENHANCES R322, not replaces it:

- **R322**: Says STOP before state transitions
- **R324**: Says UPDATE STATE FILE before that stop
- **Together**: Update state → Stop → Wait for user → Continue in new state

## Critical Reminders

1. **The state file is the ONLY memory between runs**
2. **current_state determines where to continue**
3. **Without updating it, you're stuck forever**
4. **This is not optional - it's MANDATORY**

## Penalty for Violations

- **Infinite loop detected**: -100% IMMEDIATE FAILURE
- **State not updated before stop**: -50% penalty
- **State update not committed**: -30% penalty
- **Wrong state in file**: -40% penalty

## Notes

- This rule prevents the most common orchestrator bug
- Every state transition MUST follow this pattern
- No exceptions, no shortcuts
- The orchestrator has no other way to remember its state