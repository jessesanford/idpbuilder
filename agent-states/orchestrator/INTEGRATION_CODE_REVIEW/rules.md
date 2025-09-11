# INTEGRATION_CODE_REVIEW State Rules

## State Purpose
Spawn Code Reviewer to perform quality review of integrated code after wave integration is complete. This state ensures the merged code maintains quality standards and catches integration-introduced issues before architect review.

## Entry Criteria
- **From**: MONITORING_INTEGRATION
- **Condition**: Integration Agent has completed wave integration successfully
- **Required**: 
  - Integration workspace exists with merged code
  - All effort branches merged into wave integration branch
  - Basic merge conflicts resolved by Integration Agent

## State Actions

### 1. IMMEDIATE: Spawn Code Reviewer for Integration Review
```bash
# Spawn Code Reviewer to review integrated code
/spawn agent-code-reviewer INTEGRATION_CODE_REVIEW \
  --wave "W${current_wave}" \
  --branch "${wave_integration_branch}" \
  --focus "integration-quality"
```

### 2. Code Reviewer Responsibilities
The spawned Code Reviewer will:
- Review the quality of integrated code
- Check for integration-introduced bugs
- Validate merge conflict resolutions
- Ensure code consistency across merged efforts
- Verify tests still pass after integration
- Check for duplicate code or functionality conflicts
- **Verify demo scripts are present and executable (R330/R291)**
- **Check demo coverage of integrated features**
- **Validate demo documentation completeness**
- **Review demo execution results from integration testing**
- Create INTEGRATION_CODE_REVIEW_REPORT.md

### 3. Update State File
```json
{
  "current_state": "INTEGRATION_CODE_REVIEW",
  "phase": "integration",
  "review_status": "code_review_in_progress",
  "integration_review": {
    "reviewer": "agent-code-reviewer",
    "branch": "${wave_integration_branch}",
    "focus": "integration_quality",
    "started_at": "timestamp"
  }
}
```

## Exit Criteria

### Success Path → WAITING_FOR_INTEGRATION_CODE_REVIEW
- Code Reviewer spawned successfully
- State file updated with review details
- Transition to waiting state for results

### Failure Scenarios
- **Spawn Failure** → ERROR_RECOVERY
- **Invalid Integration State** → ERROR_RECOVERY

## Key Differences from Other Reviews
- **Focus**: Integration quality, not individual effort quality
- **Scope**: Merged code interactions and conflicts
- **Timing**: After integration, before architect review
- **Purpose**: Catch integration-specific issues

## Rules Enforced
- R233: Immediate action upon state entry
- R313: Stop after spawning agent
- R238: Monitor for review reports
- R321: Any fixes require immediate backport

## Report Expected
The Code Reviewer will create:
- `INTEGRATION_CODE_REVIEW_REPORT.md` with:
  - Integration quality assessment
  - Merge conflict resolution review
  - Cross-effort consistency check
  - Test status after integration
  - **Demo verification results (R330/R291)**
    - Demo scripts present: YES/NO
    - Demo scripts executable: YES/NO
    - Demo coverage adequate: YES/NO
    - Demo documentation complete: YES/NO
    - Demo execution results: PASSED/FAILED
  - Issues found (if any)
  - Recommendation: PASS/FAIL

### Demo Verification Checklist (R330/R291)
```markdown
## Demo Verification (R330/R291 Compliance)
- [ ] All effort demo scripts present per R330 plans
- [ ] Wave integration demo script exists
- [ ] Demo documentation (DEMO.md files) present
- [ ] Demo scripts are executable (chmod +x)
- [ ] Demo test data files available
- [ ] Demo execution logs captured in demo-results/
- [ ] Demo failures documented with root cause
- [ ] Demo coverage matches integrated features
```

## Transition Rules
- **ALWAYS** → WAITING_FOR_INTEGRATION_CODE_REVIEW (after spawn)
- **NEVER** skip directly to WAVE_REVIEW
- **NEVER** proceed without Code Reviewer spawn