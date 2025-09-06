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

### 2. Stop and Summarize Protocol
When ready to transition, the orchestrator MUST:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: [CURRENT_STATE] → [NEXT_STATE]

### ✅ Current State Work Completed:
- [Summary of work done in current state]
- [Key accomplishments]
- [Any issues or blockers encountered]

### 📊 Current Status:
- Current State: [STATE]
- Next State: [NEXT_STATE]
- TODOs Completed: [X/Y]
- State Files: Updated and committed ✅

### 📝 State Persistence:
- TODOs saved to: todos/orchestrator-[STATE]-[timestamp].todo
- State file committed: [commit hash]

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to [NEXT_STATE]. Please use the appropriate continuation command.
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

### ✅ CORRECT: Stop at transition
```
Completing SPAWN_AGENTS state work...
All agents spawned successfully.

## 🛑 STATE TRANSITION CHECKPOINT: SPAWN_AGENTS → MONITOR

### ✅ Current State Work Completed:
- Spawned 3 SW-Engineer agents for parallel implementation
- All agents confirmed startup in correct directories
- Timestamps verified within 5s window

### 📊 Current Status:
- Current State: SPAWN_AGENTS
- Next State: MONITOR
- TODOs Completed: 5/5
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
```

### ❌ WRONG: Automatic transition
```
Completing SPAWN_AGENTS state work...
All agents spawned successfully.
Transitioning to MONITOR state...  # WRONG - No stop!
```

## Related Rules
- **R287** - TODO Persistence Requirements (still required)
- **R288** - State File Update Requirements (still required)
- **R234** - Mandatory State Traversal (still required)
- **R206** - State Machine Validation (still required)

## Notes
- This rule fundamentally changes orchestrator behavior from continuous to checkpoint-based operation
- Each state transition becomes a natural breakpoint for user review and control
- Improves debuggability and user oversight of the orchestration process
- Preserves context better by creating explicit checkpoints