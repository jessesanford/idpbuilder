# PROJECT_INTEGRATION_CODE_REVIEW State Rules

## State Purpose
Spawn Code Reviewer to perform comprehensive quality review of the fully integrated project after all phases are merged. This is the final code quality gate before project completion.

## Entry Criteria
- **From**: MONITORING_PROJECT_INTEGRATION
- **Condition**: Integration Agent has completed project integration successfully
- **Required**: 
  - Project integration workspace exists with all phases merged
  - All phase branches merged into project branch
  - Project-level merge conflicts resolved
  - No build failures from integration

## State Actions

### 1. IMMEDIATE: Spawn Code Reviewer for Project Integration Review
```bash
# Spawn Code Reviewer for comprehensive project review
/spawn agent-code-reviewer PROJECT_INTEGRATION_CODE_REVIEW \
  --project "${project_name}" \
  --branch "${project_integration_branch}" \
  --focus "project-integration-quality" \
  --comprehensive true
```

### 2. Code Reviewer Responsibilities
The spawned Code Reviewer will:
- Perform comprehensive project-wide code review
- Validate all phases integrate correctly
- Check for cross-phase conflicts or duplications
- Ensure architectural consistency across entire project
- Verify all project requirements are met
- Validate test suite completeness and passing
- Check for performance regressions
- Review security implications
- Assess technical debt introduced
- **Verify end-to-end demo scenarios (R330/R291)**
- **Validate production readiness demonstrations**
- **Check comprehensive feature coverage in demos**
- **Review full project demo execution results**
- **Ensure demo documentation is production-ready**
- Create PROJECT_INTEGRATION_CODE_REVIEW_REPORT.md

### 3. Update State File
```json
{
  "current_state": "PROJECT_INTEGRATION_CODE_REVIEW",
  "phase": "project_integration",
  "review_status": "project_code_review_in_progress",
  "project_integration_review": {
    "reviewer": "agent-code-reviewer",
    "branch": "${project_integration_branch}",
    "focus": "project_integration_quality",
    "phases_integrated": ["P1", "P2", "P3"],
    "comprehensive_review": true,
    "started_at": "timestamp"
  }
}
```

## Exit Criteria

### Success Path → WAITING_FOR_PROJECT_INTEGRATION_CODE_REVIEW
- Code Reviewer spawned successfully
- State file updated with project review details
- Transition to waiting state for comprehensive results

### Failure Scenarios
- **Spawn Failure** → ERROR_RECOVERY
- **Invalid Project State** → ERROR_RECOVERY

## Key Differences from Phase Integration Review
- **Scope**: Entire project (all phases)
- **Depth**: Most comprehensive review
- **Focus**: Project-wide consistency and completeness
- **Validation**: All requirements met
- **Final Gate**: Last quality check before completion

## Rules Enforced
- R233: Immediate action upon state entry
- R313: Stop after spawning agent
- R238: Monitor for review reports
- R283: Project must include all phases
- R266: Comprehensive validation required
- R321: Any fixes require immediate backport

## Report Expected
The Code Reviewer will create:
- `PROJECT_INTEGRATION_CODE_REVIEW_REPORT.md` with:
  - Project-wide quality assessment
  - Cross-phase integration verification
  - Requirements completion checklist
  - Architectural consistency review
  - Performance baseline comparison
  - Security audit results
  - Technical debt assessment
  - Test coverage report
  - Critical issues (if any)
  - Risk assessment
  - Recommendation: PASS/FAIL/CONDITIONAL_PASS

## Special Considerations
- This is the FINAL code review gate
- More stringent than phase reviews
- May require multiple iterations
- Focus on production readiness

## Transition Rules
- **ALWAYS** → WAITING_FOR_PROJECT_INTEGRATION_CODE_REVIEW (after spawn)
- **NEVER** skip directly to SPAWN_CODE_REVIEWER_PROJECT_VALIDATION
- **NEVER** proceed without comprehensive review