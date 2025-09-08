# Orchestrator - SPAWN_CODE_REVIEWERS_FOR_SPLITS State Rules

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

You MUST STOP IMMEDIATELY after spawning Code Reviewers to preserve context:
- Record what was spawned in state file
- Save TODOs and commit state changes
- EXIT with clear continuation instructions

## State Context

**Purpose:**
Spawn Code Reviewer agents to review split implementations when effort splits are required.

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

### 🚨🚨🚨 R269 - Code Reviewer Merge Plan No Execution (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R269-code-reviewer-merge-plan-no-execution.md`
**Criticality**: BLOCKING - Code Reviewer only plans, never executes
**Summary**: Code Reviewer creates plan, Integration Agent executes

### 🚨🚨🚨 R151 - Parallel Spawning Timestamp Requirement (CRITICAL)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-spawning-timing.md`
**Criticality**: CRITICAL - Timestamps must be within 5s for parallel agents
**Summary**: All parallel agents must emit timestamps within 5 seconds

## Primary Actions

1. **Identify Split Reviews Needed**:
   - Check which splits need review
   - Prepare review instructions for each split
   
2. **Spawn Code Reviewers**:
   - CD to effort directory first (R208)
   - Spawn Code Reviewer for each split
   - Provide split-specific instructions
   
3. **Update State File**:
   - Record spawned agents
   - Update split_tracking with review status
   - Commit and push state changes

## State Transitions

- **SUCCESS** → MONITOR_CODE_REVIEW (agents spawned)
- **FAILURE** → ERROR_RECOVERY (spawn failed)

## Success Criteria

Before transitioning:
1. ✅ All split reviewers spawned
2. ✅ State file updated with agent info
3. ✅ TODOs saved and committed
4. ✅ Clear continuation instructions provided