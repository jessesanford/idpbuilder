# Orchestrator - WAITING_FOR_PROJECT_MERGE_PLAN State Rules

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

## State Context

**Purpose:**
Monitor Code Reviewer creating the project-level merge plan that will integrate all phase branches.

## Primary Actions

1. **Check for Merge Plan**:
   - Look for PROJECT-MERGE-PLAN.md in project integration workspace
   - Verify plan includes all phases in correct order (R270)
2. **Validate Plan Completeness**:
   - All phase branches listed
   - Merge order specified
   - Conflict resolution strategy documented
3. **Update State** when plan is ready

## Valid State Transitions

- **SUCCESS** → SPAWN_INTEGRATION_AGENT_PROJECT (plan ready)
- **TIMEOUT** → ERROR_RECOVERY (plan not created)
- **FAILURE** → ERROR_RECOVERY (invalid plan)