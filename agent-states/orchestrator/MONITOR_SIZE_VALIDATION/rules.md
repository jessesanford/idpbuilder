# Orchestrator - MONITOR_SIZE_VALIDATION State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.json with new state
3. ✅ Committing and pushing the state file  
4. ✅ Providing work summary

### YOU MUST NOT:
- ❌ Continue to the next state automatically
- ❌ Start work for the new state
- ❌ Spawn agents for the new state
- ❌ Assume permission to continue

### STOP PROTOCOL:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: CURRENT_STATE → NEXT_STATE

### ✅ Current State Work Completed:
- [List completed work]

### 📊 Current Status:
- Current State: CURRENT_STATE
- Next State: NEXT_STATE
- TODOs Completed: X/Y
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to NEXT_STATE. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---

## Primary Rules

### 🚨🚨🚨 R319 - Orchestrator NEVER Measures or Assesses Code (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-or-assesses.md`
**Criticality**: BLOCKING - Orchestrator coordinates, never analyzes
**Summary**: Orchestrator must delegate all code analysis to appropriate agents

## State Context

**Purpose:**
Monitor size validation checks for efforts.

## Primary Actions

1. Check line count results
2. Determine if splits needed
3. Update tracking

## State Transitions

- **SUCCESS** → SPAWN_AGENTS or CREATE_SPLIT_PLAN
- **FAILURE** → ERROR_RECOVERY

## Success Criteria

Before transitioning:
1. ✅ State work completed
2. ✅ State file updated
3. ✅ TODOs saved and committed
4. ✅ Clear continuation instructions provided
