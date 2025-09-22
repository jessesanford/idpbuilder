# SPAWN_CODE_REVIEWER_INTEGRATION_FIX_PLAN State Rules

# PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R256** - Fix Planning Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R256-fix-planning-protocol.md`

4. **R287** - TODO Persistence Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`

5. **R288** - State File Update Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-requirements.md`

6. **R304** - Mandatory Line Counter Usage (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-usage.md`

7. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-after-spawn.md`

8. **R324** - State Transition Validation (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R324-state-transition-validation.md`


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

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
