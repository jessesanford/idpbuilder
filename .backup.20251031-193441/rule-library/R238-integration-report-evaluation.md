# 🚨🚨🚨 BLOCKING RULE R238: Integration Report Evaluation Protocol

## Criticality: BLOCKING
**Failure to check integration reports = -100% GRADE**

## Description
The orchestrator MUST check for and evaluate integration reports from the Integration Agent to determine if fixes are needed.

## Requirements

### 1. Report Detection (MANDATORY)
In MONITORING_INTEGRATE_WAVE_EFFORTS state, the orchestrator MUST:
```bash
# Check for integration report
REPORT_FILE="efforts/phase${PHASE}/wave${WAVE}/integration-workspace/INTEGRATE_WAVE_EFFORTS_REPORT.md"
if [ ! -f "$REPORT_FILE" ]; then
    # ERROR: No report found
    transition_to ERROR_RECOVERY
fi
```

### 2. Status Evaluation (CRITICAL)
Parse the integration report for:
- Integration Status (PROJECT_DONE/FAILED/BLOCKED)
- Build Status (PASSING/FAILED/BLOCKED_BY_DEPENDENCIES)
- Test Status (PASSING/FAILED/NOT_RUN)
- Failed Branches (list of branches that failed)
- Missing Dependencies (if any)

### 3. Decision Logic (MANDATORY)
```
IF status == PROJECT_DONE AND build == PASSING AND tests == PASSING:
    → REVIEW_WAVE_ARCHITECTURE (proceed normally)
ELIF status == FAILED OR build == BLOCKED_BY_DEPENDENCIES:
    → INTEGRATE_WAVE_EFFORTS_FEEDBACK_REVIEW (initiate fix cycle)
ELSE:
    → ERROR_RECOVERY (unexpected state)
```

### 4. Phase Integration Evaluation
Same protocol applies to MONITORING_INTEGRATE_PHASE_WAVES with INTEGRATE_PHASE_WAVES_REPORT.md

## Violations

### AUTOMATIC FAILURE (-100%)
- Ignoring integration failures and marking complete
- Not checking for integration report
- Transitioning to REVIEW_WAVE_ARCHITECTURE when integration failed

### MAJOR VIOLATIONS (-50%)
- Not parsing all status fields
- Missing dependency detection
- Incorrect state transition

## Implementation Example

```bash
# In MONITORING_INTEGRATE_WAVE_EFFORTS state
check_integration_report() {
    local report="$1"
    
    # Extract statuses
    INT_STATUS=$(grep "Integration Status:" "$report" | cut -d: -f2)
    BUILD_STATUS=$(grep "Build Status:" "$report" | cut -d: -f2)
    
    # Evaluate and transition
    if [[ "$INT_STATUS" == "PROJECT_DONE" ]]; then
        next_state="REVIEW_WAVE_ARCHITECTURE"
    else
        next_state="INTEGRATE_WAVE_EFFORTS_FEEDBACK_REVIEW"
    fi
    
    update_state "$next_state"
}
```

## Iteration Container Integration (SF 3.0)

This rule operates within **iteration containers** that expect multiple integration-fix-reintegrate cycles:

### Convergence Expectation
- First integration: Often produces bugs/conflicts → triggers fix cycle
- Second integration: Fewer issues → may require another iteration
- Nth integration: Clean convergence → proceed to architecture review

### Container State Tracking
```json
{
  "container_id": "wave-1-integration",
  "current_iteration": 3,
  "convergence_status": "CONVERGING",
  "integration_attempts": [
    {"iteration": 1, "status": "FAILED", "bugs_found": 12},
    {"iteration": 2, "status": "FAILED", "bugs_found": 3},
    {"iteration": 3, "status": "PROJECT_DONE", "bugs_found": 0}
  ]
}
```

This report evaluation determines:
1. If iteration container continues (bugs found → fix cycle)
2. If convergence achieved (no bugs → architecture review)
3. If escalation needed (no progress → ERROR_RECOVERY)

## Related Rules
- R260: Integration Agent Core Requirements
- R263: Integration Documentation Requirements
- R239: Fix Plan Distribution Protocol
- R300: Comprehensive Fix Management Protocol
- R520: Integration Attempt Tracking (iteration container management)

## Grading Impact
- **Correct evaluation**: +20% compliance bonus
- **Ignored failures**: -100% automatic failure
- **Partial compliance**: -50% major violation