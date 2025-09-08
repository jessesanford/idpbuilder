# Orchestrator - SPAWN_CODE_REVIEWER_PROJECT_VALIDATION State Rules

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

You MUST STOP IMMEDIATELY after spawning the Code Reviewer.

## State Context

**Purpose:**
Spawn Code Reviewer to perform comprehensive validation of the integrated project, ensuring all phases work together correctly.

## Primary Actions

1. **Spawn Code Reviewer** with PROJECT_VALIDATION directive
2. **Provide Context**:
   - Project integration branch location
   - List of all integrated phases
   - Validation requirements per R271-R275
3. **Request Validation**:
   - Functional correctness across phases
   - No inter-phase conflicts
   - Production readiness
4. **Update State File** and **STOP per R313**

## Valid State Transitions

- **SUCCESS** → WAITING_FOR_PROJECT_VALIDATION (reviewer spawned)
- **FAILURE** → ERROR_RECOVERY (spawn failed)