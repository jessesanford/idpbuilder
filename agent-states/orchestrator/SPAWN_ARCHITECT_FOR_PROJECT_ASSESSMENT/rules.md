# Orchestrator - SPAWN_ARCHITECT_FOR_PROJECT_ASSESSMENT State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.yaml with new state
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

## State Context

**Purpose:**
Spawn Architect agent to assess entire project completion.

## Primary Rules

### 🚨🚨🚨 R208 - CD Before Spawn (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R208-cd-before-spawn.md`
**Criticality**: SUPREME LAW - Must CD to effort directory before spawn
**Summary**: Always CD to the correct working directory before spawning agents

### 🚨🚨🚨 R216 - Bash Execution Syntax Protocol (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R216-bash-execution-syntax-protocol.md`
**Criticality**: BLOCKING - Incorrect syntax causes failures
**Summary**: Use parentheses for subshells, proper variable syntax

### 🚨🚨🚨 R235 - Pre-flight Verification Checklist (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R235-pre-flight-verification-checklist.md`
**Criticality**: BLOCKING - Must verify environment before spawning
**Summary**: Check directories, permissions, branches before agent spawn

## Primary Actions

1. Prepare project assessment context
2. Spawn Architect with PROJECT_ASSESSMENT state
3. Update state file

## State Transitions

- **SUCCESS** → MONITOR_ARCHITECT_REVIEW
- **FAILURE** → ERROR_RECOVERY

## Success Criteria

Before transitioning:
1. ✅ State work completed
2. ✅ State file updated
3. ✅ TODOs saved and committed
4. ✅ Clear continuation instructions provided
