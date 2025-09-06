# Orchestrator - WAITING_FOR_PROJECT_VALIDATION State Rules

## State Context

**Purpose:**
Wait for Code Reviewer to complete comprehensive validation of the integrated project.

## Primary Actions

1. **Check for Validation Report**:
   - Look for PROJECT-VALIDATION-REPORT.md
   - Review validation results
2. **Evaluate Results**:
   - PASS: All phases work together correctly
   - FAIL: Inter-phase issues found
   - BLOCKED: Cannot validate due to errors
3. **Document Issues** if validation fails

## Valid State Transitions

- **PASS** → CREATE_INTEGRATION_TESTING (project validated)
- **FAIL** → ERROR_RECOVERY (validation failed)
- **BLOCKED** → ERROR_RECOVERY (cannot validate)