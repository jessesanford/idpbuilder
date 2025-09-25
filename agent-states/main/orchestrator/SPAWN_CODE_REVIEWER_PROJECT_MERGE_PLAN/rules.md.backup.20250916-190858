# Orchestrator - SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN State Rules

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

## 🛑🛑🛑 R322 Part A MANDATORY STOP AFTER SPAWNING AGENTS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

You MUST STOP IMMEDIATELY after spawning the Code Reviewer to preserve context:
- Record what was spawned in state file
- Save TODOs and commit state changes
- EXIT with clear continuation instructions
- This prevents agent responses from overflowing context

## State Context

**Purpose:**
Spawn Code Reviewer to create a comprehensive merge plan for integrating all phase branches into the project integration branch.

## Primary Actions

1. **Spawn Code Reviewer** with PROJECT_MERGE_PLANNING directive
2. **Provide Context**:
   - List of all phase integration branches
   - Target project integration branch
   - Dependency order requirements (R270)
3. **Update State File** with spawned agent details
4. **STOP per R322 Part A** - Exit immediately after spawning

## Valid State Transitions

- **SUCCESS** → WAITING_FOR_PROJECT_MERGE_PLAN (agent spawned)
- **FAILURE** → ERROR_RECOVERY (spawn failed)