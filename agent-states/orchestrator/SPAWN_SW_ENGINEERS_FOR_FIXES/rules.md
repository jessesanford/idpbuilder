# Orchestrator - SPAWN_SW_ENGINEERS_FOR_FIXES State Rules

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

## 🛑🛑🛑 R313 MANDATORY STOP AFTER SPAWNING AGENTS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

You MUST STOP IMMEDIATELY after spawning SW Engineers to preserve context.

## State Context

**Purpose:**
Spawn Software Engineer agents to fix issues identified by Code Reviewers.

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

### 🚨🚨🚨 R295 - SW Engineer Spawn Clarity Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R295-sw-engineer-spawn-clarity-protocol.md`
**Criticality**: SUPREME LAW - Must provide clear instructions
**Summary**: Every spawn MUST include state name, exact plan file, and clear instructions

### 🚨🚨🚨 R197 - One Agent Per Effort (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R197-one-agent-per-effort.md`
**Criticality**: BLOCKING - Never spawn multiple agents for same effort
**Summary**: Only one SW Engineer per effort at a time

### 🚨🚨🚨 R151 - Parallel Spawning Timestamp Requirement (CRITICAL)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-spawning-timing.md`
**Criticality**: CRITICAL - Timestamps must be within 5s for parallel agents
**Summary**: All parallel agents must emit timestamps within 5 seconds

## Primary Actions

1. **Identify Fixes Needed**:
   - Parse review reports for required fixes
   - Prioritize critical issues first
   
2. **Spawn SW Engineers**:
   - CD to effort directory first (R208)
   - Spawn SW Engineer with FIX_ISSUES state
   - Provide specific fix instructions
   
3. **Update State File**:
   - Record spawned agents
   - Track fix implementation status
   - Commit and push state changes

## State Transitions

- **SUCCESS** → MONITOR_FIX_IMPLEMENTATION (agents spawned)
- **FAILURE** → ERROR_RECOVERY (spawn failed)

## Success Criteria

Before transitioning:
1. ✅ All fix engineers spawned
2. ✅ State file updated with agent info
3. ✅ TODOs saved and committed
4. ✅ Clear continuation instructions provided