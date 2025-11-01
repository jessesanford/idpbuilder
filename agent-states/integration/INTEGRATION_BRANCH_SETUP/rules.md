# INTEGRATE_WAVE_EFFORTS_BRANCH_SETUP State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## State Purpose
Prepare the integration branch environment, including deletion of stale branches from previous attempts.

## Entry Conditions
- Requirements loaded from INTEGRATE_WAVE_EFFORTS_LOAD_REQUIREMENTS
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
- If `MUST_DELETE_STALE=true` → INTEGRATE_WAVE_EFFORTS_DELETE_STALE
- If `MUST_DELETE_STALE=false` → INTEGRATE_WAVE_EFFORTS_CREATE_BRANCH

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
- If workspace dirty → INTEGRATE_WAVE_EFFORTS_ERROR
- If base branch missing → INTEGRATE_WAVE_EFFORTS_ABORT

## Logging Requirements
```bash
echo "[INTEGRATE_WAVE_EFFORTS_BRANCH_SETUP] Checking for stale branches"
echo "[INTEGRATE_WAVE_EFFORTS_BRANCH_SETUP] Attempt: ${ATTEMPT_NUMBER}"
echo "[INTEGRATE_WAVE_EFFORTS_BRANCH_SETUP] Target: ${TARGET_BRANCH}"
echo "[INTEGRATE_WAVE_EFFORTS_BRANCH_SETUP] Stale detected: ${MUST_DELETE_STALE}"
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

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

