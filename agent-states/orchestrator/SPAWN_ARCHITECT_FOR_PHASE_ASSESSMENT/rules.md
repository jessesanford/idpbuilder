# Orchestrator - SPAWN_ARCHITECT_FOR_PHASE_ASSESSMENT State Rules

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
Spawn Architect agent to assess completed phase and determine readiness for next phase.

## Primary Rules

### 🚨🚨🚨 R208 - CD Before Spawn (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R208-cd-before-spawn.md`
**Criticality**: SUPREME LAW - Must CD to effort directory before spawn
**Summary**: Always CD to the correct working directory before spawning agents

### 🚨🚨🚨 R257 - Mandatory Phase Assessment Report (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R257-mandatory-phase-assessment-report.md`
**Criticality**: BLOCKING - Architect must produce phase assessment
**Summary**: Phase assessment must generate PHASE-ASSESSMENT-REPORT.md

### 🚨🚨🚨 R285 - Mandatory Phase Integration Before Assessment (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R285-phase-integration-before-assessment.md`
**Criticality**: BLOCKING - Phase must be integrated before assessment
**Summary**: Cannot assess phase until all waves are integrated

### 🚨🚨🚨 R216 - Bash Execution Syntax Protocol (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R216-bash-execution-syntax-protocol.md`
**Criticality**: BLOCKING - Incorrect syntax causes failures
**Summary**: Use parentheses for subshells, proper variable syntax

### 🚨🚨🚨 R235 - Pre-flight Verification Checklist (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R235-pre-flight-verification-checklist.md`
**Criticality**: BLOCKING - Must verify environment before spawning
**Summary**: Check directories, permissions, branches before agent spawn

## Primary Actions

1. **Prepare Phase Assessment Context**:
   - Verify phase integration complete
   - Gather all wave integration results
   - Prepare assessment criteria
   
2. **Spawn Architect**:
   - CD to phase directory first (R208)
   - Spawn Architect with PHASE_ASSESSMENT state
   - Provide phase context and criteria
   
3. **Update State File**:
   - Record architect spawn
   - Track phase assessment status
   - Commit and push state changes

## State Transitions

- **SUCCESS** → MONITOR_ARCHITECT_REVIEW (architect spawned)
- **FAILURE** → ERROR_RECOVERY (spawn failed)

## Success Criteria

Before transitioning:
1. ✅ Architect spawned with correct state
2. ✅ Phase context provided to architect
3. ✅ State file updated
4. ✅ Clear continuation instructions provided