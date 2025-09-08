# 🔴🔴🔴 RULE R322: MANDATORY STOP BEFORE STATE TRANSITIONS 🔴🔴🔴

## Rule Identifier
**Rule ID:** R322  
**Category:** State Machine Control  
**Criticality:** 🔴🔴🔴 SUPREME LAW  
**Introduced:** Version 2.0.322  

## Rule Statement

**THE ORCHESTRATOR MUST STOP BEFORE EVERY STATE TRANSITION!**

After completing ALL work for the current state:
1. **STOP** - Do not automatically transition
2. **SUMMARIZE** - Provide clear summary of completed work
3. **SAVE** - Persist TODOs and commit state files (per R287/R288)
4. **WAIT** - For user to explicitly continue

## Detailed Requirements

### 1. Completion of Current State
Before stopping, the orchestrator MUST:
- ✅ Complete ALL TODOs for the current state
- ✅ Ensure all spawned agents have completed their tasks
- ✅ Verify all state requirements are met
- ✅ Update state files with completed work

### 2. State Update BEFORE Stop Protocol
When ready to transition, the orchestrator MUST:

#### 🔴🔴🔴 CRITICAL: UPDATE current_state FIRST OR GET STUCK IN INFINITE LOOP! 🔴🔴🔴

**⚠️⚠️⚠️ WARNING: FAILURE TO UPDATE current_state = INFINITE LOOP BUG! ⚠️⚠️⚠️**

The orchestrator WILL be stuck repeating the same state work forever if you don't update current_state!

```bash
# 🚨🚨🚨 MANDATORY SEQUENCE - DO NOT SKIP ANY STEP! 🚨🚨🚨
# STEP 1: Determine your next state
NEXT_STATE="[YOUR_NEXT_STATE_HERE]"  # e.g., "MONITOR", "SPAWN_AGENTS", etc.
CURRENT_STATE="[YOUR_CURRENT_STATE]"  # What state you're leaving

# STEP 2: UPDATE THE STATE FILE (THIS IS THE CRITICAL PART!)
echo "📝 CRITICAL: Updating current_state to prevent infinite loop..."
yq -i ".current_state = \"$NEXT_STATE\"" orchestrator-state.json
yq -i ".previous_state = \"$CURRENT_STATE\"" orchestrator-state.json
yq -i ".transition_time = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.json

# STEP 3: Verify the update worked
echo "✅ State file updated:"
grep "current_state:" orchestrator-state.json

# STEP 4: Commit and push IMMEDIATELY
git add orchestrator-state.json
git commit -m "state: transition from $CURRENT_STATE to $NEXT_STATE (R324 compliance)"
git push

# STEP 5: NOW you can stop (state is safely persisted)
echo "✅ State transition persisted - safe to stop"
```

**SEE ALSO: R324 - State File Update Before Stop (PREVENTS INFINITE LOOPS)**

#### THEN Stop and Summarize:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: [CURRENT_STATE] → [NEXT_STATE]

### ✅ Current State Work Completed:
- [Summary of work done in current state]
- [Key accomplishments]
- [Any issues or blockers encountered]

### 📊 Current Status:
- Current State: [NEXT_STATE] ← UPDATED IN FILE!
- Previous State: [CURRENT_STATE]
- TODOs Completed: [X/Y]
- State Files: Updated and committed ✅

### 📝 State Persistence:
- TODOs saved to: todos/orchestrator-[STATE]-[timestamp].todo
- State file updated with NEW state: [NEXT_STATE]
- State file committed: [commit hash]

### ⏸️ STOPPED - Ready to Continue in [NEXT_STATE]
When restarted, will continue from [NEXT_STATE]. Please use the appropriate continuation command.
```

### 3. Save and Commit Requirements
Before stopping, ALWAYS:
- Save TODOs using R287 protocol
- Commit state files using R288 protocol
- Push all changes to remote
- Include transition readiness in commit message

### 4. User Continuation Required
The orchestrator MUST NOT:
- ❌ Automatically transition to the next state
- ❌ Continue without explicit user instruction
- ❌ Skip the summary and stop protocol
- ❌ Assume permission to continue

The orchestrator MUST:
- ✅ Wait for user to explicitly continue
- ✅ Provide clear guidance on next steps
- ✅ Maintain state readiness for continuation

## Replaces Previous Rules

This rule SUPERSEDES and REPLACES:
- **R021** - Orchestrator Never Stops (DEPRECATED by R322)
- **R231** - Continuous Operation Through Transitions (DEPRECATED by R322)
- **R313** - Mandatory Stop After Spawn (CONSOLIDATED into R322)

## Enforcement

### Success Criteria
- ✅ Orchestrator stops at EVERY state transition
- ✅ Clear summary provided for each transition
- ✅ State files and TODOs properly persisted
- ✅ User explicitly continues each transition

### Failure Conditions
- ❌ Automatic transition without stopping = -50% penalty
- ❌ Missing summary at transition = -20% penalty
- ❌ Not saving state/TODOs = -30% penalty
- ❌ Continuing without user permission = -100% FAILURE

## Examples

### ✅ CORRECT: Update state THEN stop
```bash
# First, update the state file
echo "Updating state file for transition..."
yq -i '.current_state = "MONITOR"' orchestrator-state.json
yq -i '.previous_state = "SPAWN_AGENTS"' orchestrator-state.json
yq -i ".transition_time = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.json
git add orchestrator-state.json
git commit -m "state: transition from SPAWN_AGENTS to MONITOR"
git push
```

```
## 🛑 STATE TRANSITION CHECKPOINT: SPAWN_AGENTS → MONITOR

### ✅ Current State Work Completed:
- Spawned 3 SW-Engineer agents for parallel implementation
- All agents confirmed startup in correct directories
- Timestamps verified within 5s window

### 📊 Current Status:
- Current State: MONITOR (UPDATED IN FILE!)
- Previous State: SPAWN_AGENTS
- TODOs Completed: 5/5
- State Files: Updated with new state and committed ✅

### ⏸️ STOPPED - Ready to Continue in MONITOR
When restarted, will continue from MONITOR state.
```

### ❌ WRONG: Stop without updating state (CAUSES LOOPS!)
```
Completing SPAWN_AGENTS state work...
All agents spawned successfully.

## 🛑 STATE TRANSITION CHECKPOINT: SPAWN_AGENTS → MONITOR
- Current State: SPAWN_AGENTS  # WRONG - File still says SPAWN_AGENTS!
- Next State: MONITOR
STOPPED - Awaiting continuation

# PROBLEM: When restarted, orchestrator reads current_state: SPAWN_AGENTS
# and repeats SPAWN_AGENTS work forever!
```

### ❌ WRONG: Automatic transition without stop
```
Completing SPAWN_AGENTS state work...
All agents spawned successfully.
Transitioning to MONITOR state...  # WRONG - No stop!
```

## Related Rules
- **R324** - State File Update Before Stop (CRITICAL - prevents loops!)
- **R287** - TODO Persistence Requirements (still required)
- **R288** - State File Update Requirements (still required)
- **R234** - Mandatory State Traversal (still required)
- **R206** - State Machine Validation (still required)

## Notes
- This rule fundamentally changes orchestrator behavior from continuous to checkpoint-based operation
- Each state transition becomes a natural breakpoint for user review and control
- Improves debuggability and user oversight of the orchestration process
- Preserves context better by creating explicit checkpoints