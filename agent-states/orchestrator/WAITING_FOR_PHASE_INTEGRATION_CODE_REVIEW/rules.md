# WAITING_FOR_PHASE_INTEGRATION_CODE_REVIEW State Rules

## State Purpose
Monitor Code Reviewer progress on phase integration code review and process results when complete.

## Entry Criteria
- **From**: PHASE_INTEGRATION_CODE_REVIEW
- **Condition**: Code Reviewer spawned for phase integration review
- **Required**: State file shows phase review in progress

## State Actions

### 1. IMMEDIATE: Check for Phase Review Completion
```bash
# Check for phase integration code review report
if [ -f "PHASE_INTEGRATION_CODE_REVIEW_REPORT.md" ]; then
    echo "Phase integration code review complete"
    process_phase_review_results
else
    echo "Waiting for phase integration code review to complete"
    exit 0  # Will be re-invoked later
fi
```

### 2. Process Phase Review Results
When report exists:
- Parse PHASE_INTEGRATION_CODE_REVIEW_REPORT.md
- Extract PASS/FAIL status
- Identify phase-level integration issues
- Check feature completeness
- Determine next state based on results

### 3. Update State File
```json
{
  "current_state": "WAITING_FOR_PHASE_INTEGRATION_CODE_REVIEW",
  "phase_integration_review": {
    "status": "complete",
    "result": "PASS|FAIL",
    "phase_issues_found": [],
    "feature_complete": true,
    "cross_wave_conflicts": [],
    "report": "PHASE_INTEGRATION_CODE_REVIEW_REPORT.md"
  }
}
```

## Exit Criteria

### Success Path (PASS) → SPAWN_ARCHITECT_PHASE_ASSESSMENT
- Review passed with no critical issues
- Phase integration quality verified
- Ready for architect assessment

### Failure Path (FAIL) → SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN
- Critical phase integration issues found
- Cross-wave conflicts detected
- Need comprehensive fix plan

### Waiting Path → WAITING_FOR_PHASE_INTEGRATION_CODE_REVIEW
- Review still in progress
- Exit and wait for next check

## Report Processing
Parse PHASE_INTEGRATION_CODE_REVIEW_REPORT.md for:
- Overall phase status (PASS/FAIL)
- Critical cross-wave issues
- Feature completeness assessment
- Architectural violations
- Performance regressions
- Test coverage gaps

## Rules Enforced
- R233: Immediate check upon entry
- R238: Monitor for completion
- R285: Phase completeness validation
- R321: Fixes require backport planning

## Transition Rules
- **If PASS** → SPAWN_ARCHITECT_PHASE_ASSESSMENT
- **If FAIL** → SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN
- **If pending** → WAITING_FOR_PHASE_INTEGRATION_CODE_REVIEW (self)