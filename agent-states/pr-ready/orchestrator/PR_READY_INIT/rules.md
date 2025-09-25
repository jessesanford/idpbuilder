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