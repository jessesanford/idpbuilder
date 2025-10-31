# PR_READY_INIT State Rules

## 🔴🔴🔴 STATE PURPOSE: Initialize PR-Ready Transformation 🔴🔴🔴

### MANDATORY ACTIONS (R233 COMPLIANT)
Upon entering this state, IMMEDIATELY:

1. **Load PR transformation configuration**
   ```bash
   # Read transformation requirements
   cat SOFTWARE-FACTORY-PR-READY-STATE-MACHINE.md
   ```

2. **Initialize PR-ready state file**
   ```json
   {
     "transformation_id": "pr-ready-<timestamp>",
     "current_state": "PR_READY_INIT",
     "phase": "DISCOVERY",
     "branches_to_transform": [],
     "artifacts_discovered": {},
     "cleanup_status": {},
     "validation_results": {},
     "errors": []
   }
   ```

3. **Verify workspace is clean**
   - Ensure on main branch
   - No uncommitted changes
   - Origin is accessible

4. **Document transformation scope**
   - List all effort branches
   - Identify target upstream
   - Set success criteria

### EXIT CRITERIA
✅ Configuration loaded
✅ State file initialized
✅ Workspace verified
✅ Scope documented

### TRANSITIONS
- Success → PR_DISCOVERY_ASSESSMENT
- Error → PR_ERROR_DETECTED

### PROHIBITED ACTIONS
❌ Do NOT write any code
❌ Do NOT modify branches
❌ Do NOT spawn agents yet
❌ Do NOT skip initialization

### ERROR HANDLING
- Missing state machine file → ABORT
- Dirty workspace → ABORT with instructions
- No effort branches found → ABORT

### TIMING REQUIREMENTS
- Complete within 60 seconds
- Save state immediately after init
- Commit state file before transition

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Update state file
jq '.state_machine.current_state = "NEXT_STATE"' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# 3. Save TODOs
save_todos "R322_CHECKPOINT"

# 4. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 5. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

