# SPAWN_CODE_REVIEWER_INTEGRATION_FIX_PLAN State Rules

## State Purpose
Spawn Code Reviewer to create fix plans for integration issues identified during code review. This state is reached when integration code review fails.

## Entry Criteria
- **From**: WAITING_FOR_INTEGRATION_CODE_REVIEW
- **Condition**: Integration code review failed with issues identified
- **Required**: 
  - INTEGRATION_CODE_REVIEW_REPORT.md exists with FAIL status
  - Integration issues documented
  - Fix plans needed

## State Actions

### 1. IMMEDIATE: Spawn Code Reviewer for Fix Planning
```bash
# Spawn Code Reviewer to create integration fix plans
/spawn agent-code-reviewer CREATE_INTEGRATION_FIX_PLAN \
  --wave "W${current_wave}" \
  --branch "${wave_integration_branch}" \
  --issues-file "INTEGRATION_CODE_REVIEW_REPORT.md" \
  --focus "integration-fixes"
```

### 2. Code Reviewer Responsibilities
The spawned Code Reviewer will:
- Analyze integration code review failures
- Create detailed fix plans for each issue
- Determine which source branches need fixes
- Identify backport requirements per R321
- Create fix distribution strategy
- Generate INTEGRATION_FIX_PLAN.md

### 3. Update State File
```json
{
  "current_state": "SPAWN_CODE_REVIEWER_INTEGRATION_FIX_PLAN",
  "phase": "integration_fixes",
  "fix_planning": {
    "reviewer": "agent-code-reviewer",
    "issues_from": "INTEGRATION_CODE_REVIEW_REPORT.md",
    "target_branch": "${wave_integration_branch}",
    "started_at": "timestamp"
  }
}
```

## Exit Criteria

### Success Path → WAITING_FOR_INTEGRATION_FIX_PLANS
- Code Reviewer spawned successfully
- State file updated
- Transition to waiting for fix plans

### Failure Scenarios
- **Spawn Failure** → ERROR_RECOVERY
- **Missing Review Report** → ERROR_RECOVERY

## Rules Enforced
- R233: Immediate action upon state entry
- R313: Stop after spawning agent
- R321: Fixes must go to source branches
- R300: Comprehensive fix management

## Expected Output
The Code Reviewer will create:
- `INTEGRATION_FIX_PLAN.md` with:
  - Issue categorization
  - Fix assignments by effort
  - Backport requirements
  - Fix sequence and dependencies
  - Validation criteria

## Transition Rules
- **ALWAYS** → WAITING_FOR_INTEGRATION_FIX_PLANS (after spawn)
- **NEVER** continue without spawning reviewer