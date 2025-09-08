# ⚠️ DEPRECATED - SUPERSEDED BY R300 ⚠️
**This rule has been consolidated into R300: Comprehensive Fix Management Protocol**
**Please use R300 for all fix management requirements**

# 🚨🚨🚨 BLOCKING RULE R240: Integration Fix Execution Protocol

## Criticality: BLOCKING
**Orchestrator executing fixes = -100% AUTOMATIC FAILURE**

## Description
Integration fixes MUST be executed by Software Engineers, NEVER by the orchestrator. The orchestrator only coordinates.

## Requirements

### 1. Orchestrator Responsibilities (COORDINATION ONLY)
The orchestrator MUST:
- Spawn SW Engineers for fixes
- Monitor fix progress
- Check for completion flags
- Spawn Code Reviewers for review

The orchestrator MUST NOT:
- Execute any fixes directly
- Write any code
- Modify effort files
- Run build commands

### 2. SW Engineer Fix Execution
Engineers in FIX_INTEGRATION_ISSUES state MUST:
```bash
# 🔴🔴🔴 CRITICAL: Apply fixes to EFFORT BRANCH per R299 🔴🔴🔴
# NEVER apply fixes to integration branches - they will be LOST!

# 1. Verify correct location (MANDATORY per R299)
cd /efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}
CURRENT_BRANCH=$(git branch --show-current)
if [[ ! "$CURRENT_BRANCH" =~ ^effort- ]]; then
    echo "🔴 VIOLATION OF R299: Not on effort branch!"
    exit 1
fi

# 2. Read fix plan
FIX_PLAN=$(cat FIX_PLAN_LOCATION.txt)
cat "$FIX_PLAN"

# 3. Implement fixes IN EFFORT BRANCH
- Apply all code changes to effort branch
- Install missing dependencies in effort
- Resolve conflicts in effort code
- Fix build issues in effort branch

# 4. Verify fixes
- Run verification commands in effort directory
- Ensure build passes for effort
- Run tests in effort context

# 5. Mark completion and push to effort branch
rm FIX_REQUIRED.flag
echo "Fixes applied to effort branch: $CURRENT_BRANCH" > FIX_COMPLETE.flag
git add -A
git commit -m "fix: Integration issues resolved in effort branch"
git push origin "$CURRENT_BRANCH"  # Push to EFFORT branch, not integration!
```

### 3. Completion Detection
Orchestrator in MONITORING_FIX_PROGRESS MUST detect:
- Presence of `FIX_COMPLETE.flag`
- Absence of `FIX_REQUIRED.flag`
- Commit history showing fixes

### 4. Fix Review Requirement
After fixes complete:
1. Orchestrator spawns Code Reviewers
2. Reviews verify fixes are correct
3. Only then proceed to retry integration

## Violations

### AUTOMATIC FAILURE (-100%)
- Orchestrator executing fixes directly
- Orchestrator writing code
- Orchestrator modifying effort files

### MAJOR VIOLATIONS (-50%)
- Not spawning engineers for fixes
- Missing completion detection
- Skipping fix review

## Retry Limits

Maximum retry attempts to prevent infinite loops:
- Wave integration: 3 attempts
- Phase integration: 2 attempts
- After max retries: transition to ERROR_RECOVERY

## Implementation Example

```bash
# CORRECT: Orchestrator spawns engineer
spawn_fix_engineer() {
    local effort="$1"
    echo "@agent-software-engineer Fix integration issues in $effort"
    # Create command file
    # Spawn engineer
    # Monitor progress
}

# WRONG: Orchestrator fixes directly
fix_integration_directly() {  # NEVER DO THIS!
    cd "$effort_dir"
    vim broken_file.go  # VIOLATION!
    make build          # VIOLATION!
}
```

## Related Rules
- R299: Fix Application to Effort Branches Protocol (SUPREME LAW)
- R006: Orchestrator Never Writes Code
- R197: One Agent Per Effort
- R238: Integration Report Evaluation
- R239: Fix Plan Distribution Protocol

## Grading Impact
- **Proper delegation**: +20% compliance bonus
- **Orchestrator writes code**: -100% AUTOMATIC FAILURE
- **Missing engineers**: -50% major violation