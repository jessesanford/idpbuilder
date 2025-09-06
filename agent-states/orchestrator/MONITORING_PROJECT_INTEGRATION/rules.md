# Orchestrator - MONITORING_PROJECT_INTEGRATION State Rules

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