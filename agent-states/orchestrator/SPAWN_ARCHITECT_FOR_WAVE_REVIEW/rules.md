# Orchestrator - SPAWN_ARCHITECT_FOR_WAVE_REVIEW State Rules

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

## 🛑🛑🛑 R313 MANDATORY STOP AFTER SPAWNING AGENTS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

You MUST STOP IMMEDIATELY after spawning Architect to preserve context.

## State Context

**Purpose:**
Spawn Architect agent to review completed wave and assess readiness for integration.

## Primary Rules

### 🚨🚨🚨 R208 - CD Before Spawn (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R208-cd-before-spawn.md`
**Criticality**: SUPREME LAW - Must CD to effort directory before spawn
**Summary**: Always CD to the correct working directory before spawning agents

### 🚨🚨🚨 R258 - Mandatory Wave Review Report (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R258-mandatory-wave-review-report.md`
**Criticality**: BLOCKING - Architect must produce wave review report
**Summary**: Wave review must generate WAVE-REVIEW-REPORT.md

### 🚨🚨🚨 R256 - Wave Review Protocol (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R256-wave-review-protocol.md`
**Criticality**: BLOCKING - Must follow wave review protocol
**Summary**: Architect reviews all efforts in wave for completeness

### 🚨🚨🚨 R216 - Bash Execution Syntax Protocol (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R216-bash-execution-syntax-protocol.md`
**Criticality**: BLOCKING - Incorrect syntax causes failures
**Summary**: Use parentheses for subshells, proper variable syntax

### 🚨🚨🚨 R235 - Pre-flight Verification Checklist (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R235-pre-flight-verification-checklist.md`
**Criticality**: BLOCKING - Must verify environment before spawning
**Summary**: Check directories, permissions, branches before agent spawn

## Primary Actions

1. **Prepare Wave Review Context**:
   - Gather all wave efforts completed
   - Identify integration branch
   - Prepare review instructions
   
2. **Spawn Architect**:
   - CD to wave directory first (R208)
   - Spawn Architect with WAVE_REVIEW state
   - Provide wave context and review criteria
   
3. **Update State File**:
   - Record architect spawn
   - Track wave review status
   - Commit and push state changes

## State Transitions

- **SUCCESS** → MONITOR_ARCHITECT_REVIEW (architect spawned)
- **FAILURE** → ERROR_RECOVERY (spawn failed)

## Success Criteria

Before transitioning:
1. ✅ Architect spawned with correct state
2. ✅ Wave context provided to architect
3. ✅ State file updated
4. ✅ Clear continuation instructions provided