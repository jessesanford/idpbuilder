# WAITING_FOR_INTEGRATION_CODE_REVIEW State Rules

## State Purpose
Monitor Code Reviewer progress on integration code review and process results when complete.

## Entry Criteria
- **From**: INTEGRATION_CODE_REVIEW
- **Condition**: Code Reviewer spawned for integration review
- **Required**: State file shows review in progress

## State Actions

### 1. IMMEDIATE: Check for Review Completion
```bash
# Check for integration code review report
if [ -f "INTEGRATION_CODE_REVIEW_REPORT.md" ]; then
    echo "Integration code review complete"
    process_review_results
else
    echo "Waiting for integration code review to complete"
    exit 0  # Will be re-invoked later
fi
```

### 2. Process Review Results
When report exists:
- Parse INTEGRATION_CODE_REVIEW_REPORT.md
- Extract PASS/FAIL status
- Identify any integration issues
- Determine next state based on results

### 3. Update State File
```json
{
  "current_state": "WAITING_FOR_INTEGRATION_CODE_REVIEW",
  "integration_review": {
    "status": "complete",
    "result": "PASS|FAIL",
    "issues_found": [],
    "report": "INTEGRATION_CODE_REVIEW_REPORT.md"
  }
}
```

## Exit Criteria

### Success Path (PASS) → WAVE_REVIEW
- Review passed with no critical issues
- Integration quality verified
- Ready for architect review

### Failure Path (FAIL) → SPAWN_CODE_REVIEWER_INTEGRATION_FIX_PLAN
- Critical integration issues found
- Fixes required before proceeding
- Need plan to address issues

### Waiting Path → WAITING_FOR_INTEGRATION_CODE_REVIEW
- Review still in progress
- Exit and wait for next check

## Report Processing
Parse INTEGRATION_CODE_REVIEW_REPORT.md for:
- Overall status (PASS/FAIL)
- Critical issues list
- Merge conflict concerns
- Test failures after integration
- Cross-effort conflicts

## Rules Enforced
- R233: Immediate check upon entry
- R238: Monitor for completion
- R321: Fixes require backport planning

## Transition Rules
- **If PASS** → WAVE_REVIEW
- **If FAIL** → SPAWN_CODE_REVIEWER_INTEGRATION_FIX_PLAN
- **If pending** → WAITING_FOR_INTEGRATION_CODE_REVIEW (self)