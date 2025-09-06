# Orchestrator - WAITING_FOR_PROJECT_VALIDATION State Rules

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
Wait for Code Reviewer to complete comprehensive validation of the integrated project.

## Primary Actions

1. **Check for Validation Report**:
   - Look for PROJECT-VALIDATION-REPORT.md
   - Review validation results
2. **Evaluate Results**:
   - PASS: All phases work together correctly
   - FAIL: Inter-phase issues found
   - BLOCKED: Cannot validate due to errors
3. **Document Issues** if validation fails

## Valid State Transitions

- **PASS** → CREATE_INTEGRATION_TESTING (project validated)
- **FAIL** → ERROR_RECOVERY (validation failed)
- **BLOCKED** → ERROR_RECOVERY (cannot validate)