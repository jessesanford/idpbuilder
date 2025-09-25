# INTEGRATION_BRANCH_SETUP State Rules

## State Purpose
Prepare the integration branch environment, including deletion of stale branches from previous attempts.

## Entry Conditions
- Requirements loaded from INTEGRATION_LOAD_REQUIREMENTS
- Integration type and target branch name determined
- Attempt number known (1 for first, >1 for re-attempts)

## Required Actions

### 1. Check for Existing Integration Branch
```bash
# Check local branches
git branch -a | grep "${TARGET_BRANCH}"

# Check remote branches
git ls-remote --heads origin "${TARGET_BRANCH}"
```

### 2. Determine Stale Branch Handling
```bash
if [[ ${ATTEMPT_NUMBER} -gt 1 ]]; then
    echo "[CRITICAL] This is attempt ${ATTEMPT_NUMBER} - stale branch MUST be deleted"
    MUST_DELETE_STALE=true
elif git show-ref --verify "refs/heads/${TARGET_BRANCH}" 2>/dev/null; then
    echo "[WARNING] Integration branch exists from unknown previous attempt"
    MUST_DELETE_STALE=true
else
    echo "[INFO] First attempt, no stale branches"
    MUST_DELETE_STALE=false
fi
```

### 3. Transition Decision
- If `MUST_DELETE_STALE=true` → INTEGRATION_DELETE_STALE
- If `MUST_DELETE_STALE=false` → INTEGRATION_CREATE_BRANCH

## Critical Rules

### R327 - Cascade Recreation of Stale Integrations (BLOCKING)
**MANDATORY**: After fixes are applied, integration branches become stale and MUST be recreated:
```bash
# This is MANDATORY - no exceptions
delete_stale_integration() {
    local BRANCH="$1"

    # Delete local branch
    git branch -D "${BRANCH}" 2>/dev/null || true

    # Delete remote branch
    git push origin --delete "${BRANCH}" 2>/dev/null || true

    echo "✅ Deleted stale integration branch: ${BRANCH}"
}
```

### R321 - Integration Branches are Read-Only
- Integration branches NEVER receive direct fixes
- All fixes go to source branches
- After fixes, integration must be recreated

### R352 - Overlapping Cascade Support
- Multiple fix cascades may be in progress
- Track which fixes have been applied
- Each re-attempt includes all accumulated fixes

## State Tracking Updates
```json
{
  "branch_setup": {
    "target_branch": "wave1-integration",
    "base_branch": "main",
    "stale_detected": true,
    "deletion_required": true,
    "attempt_number": 2
  }
}
```

## Validation Checks

### Pre-Deletion Validation
```bash
# Ensure we're not on the branch we're about to delete
CURRENT_BRANCH=$(git branch --show-current)
if [[ "${CURRENT_BRANCH}" == "${TARGET_BRANCH}" ]]; then
    git checkout main
fi

# Ensure no uncommitted changes
if ! git diff --quiet; then
    echo "ERROR: Uncommitted changes detected"
    exit 1
fi
```

### Post-Setup Validation
- Target branch must not exist (if deletion was required)
- Workspace must be clean
- Base branch must be checked out

## Error Handling
- If deletion fails → Log warning but continue
- If workspace dirty → INTEGRATION_ERROR
- If base branch missing → INTEGRATION_ABORT

## Logging Requirements
```bash
echo "[INTEGRATION_BRANCH_SETUP] Checking for stale branches"
echo "[INTEGRATION_BRANCH_SETUP] Attempt: ${ATTEMPT_NUMBER}"
echo "[INTEGRATION_BRANCH_SETUP] Target: ${TARGET_BRANCH}"
echo "[INTEGRATION_BRANCH_SETUP] Stale detected: ${MUST_DELETE_STALE}"
```

## Metrics to Track
- Stale branch detection rate
- Deletion success rate
- Average attempts before success
- Time since last successful integration

## Common Patterns

### Pattern 1: First Attempt
```
No stale branch → Direct to CREATE_BRANCH
```

### Pattern 2: Re-attempt After Fix
```
Stale branch exists → DELETE_STALE → CREATE_BRANCH
```

### Pattern 3: Abandoned Integration Recovery
```
Old branch found → DELETE_STALE → CREATE_BRANCH
```

## Success Criteria
✅ Stale branches identified correctly
✅ Deletion path determined
✅ Workspace prepared for integration
✅ Ready to create fresh branch