# 🚨🚨🚨 BLOCKING RULE R238: Integration Report Evaluation Protocol

## Criticality: BLOCKING
**Failure to check integration reports = -100% GRADE**

## Description
The orchestrator MUST check for and evaluate integration reports from the Integration Agent to determine if fixes are needed.

## Requirements

### 1. Report Detection (MANDATORY)
In MONITORING_INTEGRATION state, the orchestrator MUST:
```bash
# Check for integration report
REPORT_FILE="efforts/phase${PHASE}/wave${WAVE}/integration-workspace/INTEGRATION_REPORT.md"
if [ ! -f "$REPORT_FILE" ]; then
    # ERROR: No report found
    transition_to ERROR_RECOVERY
fi
```

### 2. Status Evaluation (CRITICAL)
Parse the integration report for:
- Integration Status (SUCCESS/FAILED/BLOCKED)
- Build Status (PASSING/FAILED/BLOCKED_BY_DEPENDENCIES)
- Test Status (PASSING/FAILED/NOT_RUN)
- Failed Branches (list of branches that failed)
- Missing Dependencies (if any)

### 3. Decision Logic (MANDATORY)
```
IF status == SUCCESS AND build == PASSING AND tests == PASSING:
    → WAVE_REVIEW (proceed normally)
ELIF status == FAILED OR build == BLOCKED_BY_DEPENDENCIES:
    → INTEGRATION_FEEDBACK_REVIEW (initiate fix cycle)
ELSE:
    → ERROR_RECOVERY (unexpected state)
```

### 4. Phase Integration Evaluation
Same protocol applies to MONITORING_PHASE_INTEGRATION with PHASE_INTEGRATION_REPORT.md

## Violations

### AUTOMATIC FAILURE (-100%)
- Ignoring integration failures and marking complete
- Not checking for integration report
- Transitioning to WAVE_REVIEW when integration failed

### MAJOR VIOLATIONS (-50%)
- Not parsing all status fields
- Missing dependency detection
- Incorrect state transition

## Implementation Example

```bash
# In MONITORING_INTEGRATION state
check_integration_report() {
    local report="$1"
    
    # Extract statuses
    INT_STATUS=$(grep "Integration Status:" "$report" | cut -d: -f2)
    BUILD_STATUS=$(grep "Build Status:" "$report" | cut -d: -f2)
    
    # Evaluate and transition
    if [[ "$INT_STATUS" == "SUCCESS" ]]; then
        next_state="WAVE_REVIEW"
    else
        next_state="INTEGRATION_FEEDBACK_REVIEW"
    fi
    
    update_state "$next_state"
}
```

## Related Rules
- R260: Integration Agent Core Requirements
- R263: Integration Documentation Requirements
- R239: Fix Plan Distribution Protocol
- R300: Comprehensive Fix Management Protocol

## Grading Impact
- **Correct evaluation**: +20% compliance bonus
- **Ignored failures**: -100% automatic failure
- **Partial compliance**: -50% major violation