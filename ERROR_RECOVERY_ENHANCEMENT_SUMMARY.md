# ERROR_RECOVERY State Enhancement Summary

## Date: 2025-08-27
## Manager: software-factory-manager

## Problem Identified

The ERROR_RECOVERY state rules did NOT explicitly handle architect review feedback scenarios:
1. WAVE_REVIEW with CHANGES_REQUIRED decision (R258 compliance)
2. PHASE_ASSESSMENT with NEEDS_WORK decision (R257 compliance)

The rules only mentioned generic "architecture violations" without specific instructions for reading and processing the mandatory architect reports.

## Solution Implemented

### 1. Enhanced ERROR_RECOVERY Rules (`agent-states/orchestrator/ERROR_RECOVERY/rules.md`)

#### Added Specific Error Types to Recovery Decision Matrix:

**WAVE_REVIEW_CHANGES_REQUIRED**:
- Severity: HIGH
- Strategy: ARCHITECT_DIRECTED_FIXES
- Actions:
  - Read wave review report per R258
  - Extract Issues Identified section
  - Parse Required Fixes from report
  - Map each fix to appropriate agent type
  - Spawn agents for each required fix
  - Track completion against report requirements
  - Return to INTEGRATION for re-review
- Target Time: 120 minutes
- Report Location: `wave-reviews/phase{N}/wave{W}/PHASE-{N}-WAVE-{W}-REVIEW-REPORT.md`

**PHASE_ASSESSMENT_NEEDS_WORK**:
- Severity: CRITICAL
- Strategy: PHASE_LEVEL_REMEDIATION
- Actions:
  - Read phase assessment report per R257
  - Extract all Issues Identified
  - Parse Required Fixes section
  - Create comprehensive fix plan
  - Spawn multiple agents if needed
  - Coordinate cross-wave fixes
  - Track against assessment criteria
  - Return to PHASE_COMPLETE for reassessment
- Target Time: 240 minutes
- Report Location: `phase-assessments/phase{N}/PHASE-{N}-ASSESSMENT-REPORT.md`

#### Added Report Reading Functions:

```bash
read_wave_review_issues()      # Extracts issues from wave review reports
read_phase_assessment_issues()  # Extracts issues from phase assessment reports
parse_and_assign_fixes()        # Maps issues to appropriate agent types
```

#### Added Recovery Validation Functions:

```python
validate_wave_review_fixes()      # Validates all wave review issues addressed
validate_phase_assessment_fixes()  # Validates all phase assessment issues addressed
```

#### Updated State Transition Clarifications:

For WAVE_REVIEW errors:
- `ERROR_RECOVERY` (wave fixes complete) → `INTEGRATION` → `WAVE_REVIEW` (re-review)

For PHASE_ASSESSMENT errors:
- `ERROR_RECOVERY` (phase fixes complete) → `PHASE_COMPLETE` → `WAITING_FOR_PHASE_ASSESSMENT` (reassessment)

### 2. Enhanced WAVE_REVIEW State (`agent-states/orchestrator/WAVE_REVIEW/rules.md`)

Updated CHANGES_REQUIRED and WAVE_FAILED handlers to:
- Setup ERROR_RECOVERY context with report location
- Extract issues from report sections
- Store issues in state file for ERROR_RECOVERY processing
- Clearly indicate ERROR_RECOVERY must read the report

### 3. Enhanced WAITING_FOR_PHASE_ASSESSMENT State

Updated NEEDS_WORK handler to:
- Extract Priority 1 and Priority 2 fixes from report
- Store issues in state file
- Setup ERROR_RECOVERY context with report location
- Indicate ERROR_RECOVERY must process report-documented issues

## Key Benefits

1. **Explicit Report Compliance**: ERROR_RECOVERY now explicitly reads R257/R258 mandatory reports
2. **No Verbal Feedback**: All architect feedback must come from documented reports
3. **Structured Fix Assignment**: Issues are parsed and assigned to appropriate agents
4. **Clear State Transitions**: Recovery knows exactly which state to return to after fixes
5. **Audit Trail**: All decisions and fixes are tracked against report requirements

## Rules Cross-Referenced

- **R257**: Mandatory Phase Assessment Report
- **R258**: Mandatory Wave Review Report
- **R021**: Orchestrator Never Stops (continues through recovery)
- **R156**: Error Recovery Time targets
- **R233**: All States Require Immediate Action

## Verification Steps

To verify these changes are working correctly:

1. When architect returns CHANGES_REQUIRED from wave review:
   - Check that ERROR_RECOVERY reads the wave review report
   - Verify issues are extracted from "Issues Identified" section
   - Confirm appropriate agents are spawned for fixes
   - Validate return to INTEGRATION for re-review

2. When architect returns NEEDS_WORK from phase assessment:
   - Check that ERROR_RECOVERY reads the phase assessment report
   - Verify Priority 1 fixes are extracted
   - Confirm comprehensive fix plan created
   - Validate return to PHASE_COMPLETE for reassessment

## Files Modified

1. `/workspaces/software-factory-2.0-template/agent-states/orchestrator/ERROR_RECOVERY/rules.md`
   - Added explicit handlers for architect feedback
   - Added report reading functions
   - Enhanced recovery validation

2. `/workspaces/software-factory-2.0-template/agent-states/orchestrator/WAVE_REVIEW/rules.md`
   - Enhanced CHANGES_REQUIRED handling
   - Enhanced WAVE_FAILED handling
   - Added ERROR_RECOVERY context setup

3. `/workspaces/software-factory-2.0-template/agent-states/orchestrator/WAITING_FOR_PHASE_ASSESSMENT/rules.md`
   - Enhanced NEEDS_WORK handling
   - Added issue extraction from report
   - Added ERROR_RECOVERY context setup

## Commit Information

- Commit Hash: df336f3
- Branch: rule-moves
- Timestamp: 2025-08-27 19:56:00 UTC

## Compliance Status

✅ ERROR_RECOVERY now fully compliant with R257 and R258
✅ All architect feedback flows through documented reports
✅ Complete audit trail maintained
✅ Clear recovery paths established

## Recommendation

All orchestrator implementations should be tested to ensure:
1. They transition to ERROR_RECOVERY when receiving CHANGES_REQUIRED or NEEDS_WORK
2. ERROR_RECOVERY reads the actual report files (not verbal feedback)
3. Fixes are tracked against report requirements
4. Proper state transitions occur after recovery

---

**Factory Manager Sign-off**: ERROR_RECOVERY state now explicitly handles architect review feedback per R257 and R258 requirements. No architect decision can be processed without reading the mandatory report files.