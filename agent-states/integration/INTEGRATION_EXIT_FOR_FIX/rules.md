# INTEGRATION_EXIT_FOR_FIX State Rules

## State Purpose
Prepare to exit the integration sub-state machine when fixes are required, preserving state for re-entry after fixes are applied.

## Entry Conditions
- Integration discovered issues requiring fixes
- Issues could be: build failures, test failures, performance issues
- Cannot proceed without fixing source branches
- Must exit to trigger FIX_CASCADE in main state machine

## Required Actions

### 1. Document Integration Issues
```bash
# Create detailed issue report
cat > "INTEGRATION_ISSUES_${ATTEMPT_NUMBER}.md" << EOF
# Integration Issues - Attempt ${ATTEMPT_NUMBER}

## Summary
- Integration Type: ${INTEGRATION_TYPE}
- Target Branch: ${TARGET_BRANCH}
- Failed At: ${FAILURE_POINT}
- Timestamp: $(date -Iseconds)

## Issues Found
${ISSUES_DESCRIPTION}

## Branches Requiring Fixes
${AFFECTED_BRANCHES}

## Recommended Fix Order
${FIX_PRIORITY_ORDER}
EOF
```

### 2. Save Integration State
```json
{
  "exit_context": {
    "reason": "FIX_REQUIRED",
    "attempt_number": 2,
    "failure_point": "BUILD_VALIDATION",
    "issues_found": [
      {
        "type": "BUILD_FAILURE",
        "severity": "BLOCKING",
        "description": "Undefined symbol: authenticate",
        "affected_files": ["src/auth/auth.c"],
        "affected_branches": ["effort-E1.2"]
      }
    ],
    "integration_branch_state": "wave1-integration-attempt-2",
    "fixes_required": true,
    "can_retry": true
  }
}
```

### 3. Prepare Fix Requirements for Parent
```json
{
  "fix_cascade_trigger": {
    "source": "INTEGRATION",
    "integration_type": "WAVE",
    "branches_to_fix": ["effort-E1.2", "effort-E1.3"],
    "fix_specifications": [
      {
        "branch": "effort-E1.2",
        "issue": "Missing auth implementation",
        "files": ["src/auth/auth.c"],
        "fix_type": "IMPLEMENTATION"
      }
    ],
    "validation_after_fix": {
      "must_rebuild": true,
      "must_retest": true,
      "specific_tests": ["auth_test", "integration_test"]
    }
  }
}
```

### 4. Increment Attempt Counter
```bash
# Update attempt tracking
NEXT_ATTEMPT=$((CURRENT_ATTEMPT + 1))

# Check if we're approaching limit
if [[ ${NEXT_ATTEMPT} -gt ${MAX_ATTEMPTS} ]]; then
    echo "[WARNING] Next attempt will exceed maximum (${MAX_ATTEMPTS})"
    echo "[WARNING] Consider manual intervention"
fi

# Update state file
jq --arg attempt "${NEXT_ATTEMPT}" \
   '.cycle_tracking.next_attempt = ($attempt | tonumber)' \
   "${STATE_FILE}" > tmp.json && mv tmp.json "${STATE_FILE}"
```

### 5. Create Re-Entry Checkpoint
```json
{
  "re_entry_checkpoint": {
    "branches_already_merged": ["effort-E1.1"],
    "branches_pending_merge": ["effort-E1.2", "effort-E1.3"],
    "validation_completed": {
      "merge": "PARTIAL",
      "build": "FAILED",
      "tests": "NOT_RUN"
    },
    "resume_from": "INTEGRATION_BRANCH_SETUP",
    "delete_stale_required": true
  }
}
```

## Critical Rules

### R327 - Stale Integration Marking
```bash
# Mark current integration branch as stale
echo "${TARGET_BRANCH}" >> .stale_integration_branches

# Document why it's stale
echo "Attempt ${ATTEMPT_NUMBER}: ${FAILURE_REASON}" >> .stale_reasons
```

### R321 - No Direct Fixes to Integration
- NEVER attempt to fix issues in integration branch
- ALL fixes must go to source branches
- Integration branch will be deleted on re-entry

### R352 - Support Overlapping Cascades
- Multiple fix cascades may be triggered
- Track dependencies between fixes
- Ensure all fixes are included on retry

## Exit Protocol

### 1. Update Cycle History
```json
{
  "cycle_history": [
    {
      "attempt": 2,
      "started": "2025-01-21T10:00:00Z",
      "ended": "2025-01-21T10:30:00Z",
      "duration": "30m",
      "result": "BUILD_FAILED",
      "issues_count": 3,
      "exit_state": "INTEGRATION_EXIT_FOR_FIX"
    }
  ]
}
```

### 2. Prepare Parent State Update
```bash
# Update parent state machine
jq --arg reason "FIX_REQUIRED" \
   --arg attempts "${CURRENT_ATTEMPT}" \
   '.sub_state_machine.exit_reason = $reason |
    .sub_state_machine.attempts_made = ($attempts | tonumber) |
    .sub_state_machine.ready_for_fix_cascade = true' \
   orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
```

### 3. Generate Exit Summary
```markdown
## Integration Exit Summary

**Reason**: FIX_REQUIRED
**Attempt**: 2/10
**Duration**: 30 minutes
**Next Action**: Trigger FIX_CASCADE for identified issues

### Issues to Fix:
1. Build failure in auth module
2. Undefined symbols in API
3. Test compilation errors

### Branches Affected:
- effort-E1.2 (3 issues)
- effort-E1.3 (1 issue)

**Re-entry**: After fixes are applied and reviewed
```

## Validation Before Exit

### Ensure State Preserved
```bash
# Verify state file is complete
required_fields=("exit_context" "re_entry_checkpoint" "cycle_history")
for field in "${required_fields[@]}"; do
    if ! jq -e ".${field}" "${STATE_FILE}" > /dev/null; then
        echo "ERROR: Missing required field: ${field}"
        exit 1
    fi
done
```

### Ensure Parent Notified
```bash
# Parent must know we're exiting for fixes
if [[ $(jq -r '.sub_state_machine.ready_for_fix_cascade' orchestrator-state.json) != "true" ]]; then
    echo "ERROR: Parent not properly notified of fix requirement"
    exit 1
fi
```

## Error Handling
- If state save fails → Retry with backup
- If parent update fails → Log error, continue
- If max attempts exceeded → Exit with ABORT instead

## Logging Requirements
```bash
echo "[INTEGRATION_EXIT_FOR_FIX] Exiting for fixes"
echo "[INTEGRATION_EXIT_FOR_FIX] Attempt: ${ATTEMPT_NUMBER}/${MAX_ATTEMPTS}"
echo "[INTEGRATION_EXIT_FOR_FIX] Issues found: ${ISSUE_COUNT}"
echo "[INTEGRATION_EXIT_FOR_FIX] Branches to fix: ${BRANCHES_TO_FIX}"
echo "[INTEGRATION_EXIT_FOR_FIX] Will retry after fixes"
```

## Metrics to Track
- Exit frequency per attempt number
- Most common exit reasons
- Time to fix and return
- Success rate after fix cycles

## Common Patterns

### Pattern 1: First Attempt Build Failure
```
Attempt 1 → Build fails immediately → Exit for fixes
Fix time: ~30 minutes
Re-entry success rate: 85%
```

### Pattern 2: Test Failure After Build Success
```
Attempt 1 → Build OK → Tests fail → Exit for fixes
Fix time: ~45 minutes
Re-entry success rate: 70%
```

### Pattern 3: Multiple Cascading Issues
```
Attempt 1 → Build fails → Fix
Attempt 2 → Different build issue → Fix
Attempt 3 → Tests fail → Fix
Attempt 4 → Success
Total time: 2-3 hours
```

## Success Criteria
✅ Issues fully documented
✅ State preserved for re-entry
✅ Parent notified of fix requirement
✅ Ready to trigger FIX_CASCADE

## Next State
- Transitions to INTEGRATION_PREPARE_EXIT
- Then exits to parent state machine
- Parent triggers FIX_CASCADE
- Will re-enter at INTEGRATION_BRANCH_SETUP after fixes