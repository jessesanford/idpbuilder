# Orchestrator - MONITORING_PROJECT_INTEGRATION State Rules

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
Monitor the Integration Agent as it merges all phase branches into the project integration branch per R283.

## Primary Actions

1. **Monitor Integration Progress**:
   - Check for PROJECT-INTEGRATION-REPORT.md
   - Track merge completion for each phase
   - Monitor for conflicts or failures
2. **Validate Integration Success**:
   - All phases merged successfully
   - No unresolved conflicts
   - Build/tests pass in integrated state
3. **Handle Failures per R321**:
   - If integration fails, must fix in source branches
   - Trigger IMMEDIATE_BACKPORT_REQUIRED if needed

## Valid State Transitions

- **SUCCESS** → SPAWN_CODE_REVIEWER_PROJECT_VALIDATION (all phases merged)
- **FAILURE** → ERROR_RECOVERY (integration failed)
- **CONFLICTS** → IMMEDIATE_BACKPORT_REQUIRED (R321 enforcement)