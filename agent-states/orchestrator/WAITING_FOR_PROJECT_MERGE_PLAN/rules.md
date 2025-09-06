# Orchestrator - WAITING_FOR_PROJECT_MERGE_PLAN State Rules

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