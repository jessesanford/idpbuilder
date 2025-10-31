# ⚠️ DEPRECATED - SUPERSEDED BY R300 ⚠️
**This rule has been consolidated into R300: Comprehensive Fix Management Protocol**
**Please use R300 for all fix management requirements**

# 🔴🔴🔴 SUPREME RULE R299: Fix Application to Effort Branches Protocol

## Criticality: SUPREME LAW
**Applying fixes to integration branches instead of effort branches = -100% AUTOMATIC FAILURE**

## Description
ALL fixes discovered during integration, review, or error recovery MUST be applied to the ORIGINAL EFFORT BRANCHES, not to integration branches. Integration branches are temporary and get recreated - fixes in them will be lost.

## Critical Context
This rule exists because:
1. Integration branches are created fresh from main each time (per R271)
2. Fixes applied only to integration branches are LOST when new integration branches are created
3. The same issues will reappear in subsequent integration attempts
4. This creates an infinite loop of fixing the same issues

## Requirements

### 1. FIX APPLICATION LOCATION (ABSOLUTE)
Software Engineers in FIX_ISSUES state MUST:
```bash
# MANDATORY: Work in the effort branch, NOT integration
cd /efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}

# Verify you're on the effort branch
CURRENT_BRANCH=$(git branch --show-current)
if [[ ! "$CURRENT_BRANCH" =~ ^effort- ]]; then
    echo "🔴 CRITICAL ERROR: Not on effort branch!"
    echo "Current branch: $CURRENT_BRANCH"
    echo "Expected pattern: effort-*"
    exit 1
fi

# Apply fixes to the effort branch
# Make changes, test, commit
git add -A
git commit -m "fix: [description of fix applied]"
git push origin "$CURRENT_BRANCH"
```

### 2. NEVER APPLY FIXES TO INTEGRATE_WAVE_EFFORTS BRANCHES
```bash
# FORBIDDEN - AUTOMATIC FAILURE
cd /efforts/phase${PHASE}/wave${WAVE}/integration  # ❌ WRONG!
# Making fixes here will be LOST

# FORBIDDEN - AUTOMATIC FAILURE  
git checkout integration-phase-1-wave-1  # ❌ WRONG!
# Fixes here disappear on next integration attempt
```

### 3. VERIFICATION AFTER FIX APPLICATION
After applying fixes, SW Engineer MUST verify:
```bash
# Confirm fixes are in effort branch
git log --oneline -5  # Should show fix commits
git branch --show-current  # Must be effort-* branch
git push origin HEAD  # Push to effort branch remote

# Create completion marker IN EFFORT DIRECTORY
echo "Fixes applied to effort branch $(git branch --show-current)" > FIX_COMPLETE.flag
```

### 4. ORCHESTRATOR VERIFICATION
Orchestrator MUST verify fixes are in effort branches:
```bash
for effort in "${EFFORTS[@]}"; do
    cd "/efforts/phase${PHASE}/wave${WAVE}/${effort}"
    
    # Check for fix commits in effort branch
    if git log --oneline --since="1 hour ago" | grep -q "fix:"; then
        echo "✅ Fixes found in effort branch for $effort"
    else
        echo "❌ No fixes found in effort branch for $effort"
        echo "❌ Fixes may have been applied to wrong branch!"
    fi
done
```

### 5. INTEGRATE_WAVE_EFFORTS RETRY PROTOCOL
When retrying integration after fixes:
1. Create NEW integration branch from main (R271)
2. Merge UPDATED effort branches (with fixes)
3. Fixes now included via effort branches
4. No manual fix application in integration

## Common Mistakes (AUTOMATIC FAILURES)

### ❌ Applying fixes in integration branch
```bash
# WRONG - Loses fixes on next attempt
cd integration-phase-1-wave-1
vim broken_file.go  # Fix here
git commit -m "fix: resolved issue"
# This fix is LOST when integration branch recreated!
```

### ❌ Cherry-picking fixes from integration
```bash
# WRONG - Still working in wrong branch
git checkout integration-phase-1-wave-1
git cherry-pick abc123  # Getting fix from somewhere
# Still in integration branch - will be lost!
```

### ✅ CORRECT: Apply to effort, then re-integrate
```bash
# RIGHT - Fixes persist
cd /efforts/phase1/wave1/effort-api
git checkout effort-api-phase-1-wave-1
vim api.go  # Fix here
git commit -m "fix: resolved integration issue"
git push

# Later, orchestrator creates new integration
cd /efforts/phase1/wave1/integration
git checkout -b integration-phase-1-wave-1-retry2 main
git merge effort-api-phase-1-wave-1  # Includes fix!
```

## Enforcement Checkpoints

### 1. Before Fix Application
- Verify working directory is effort directory
- Verify current branch is effort-* branch
- Verify remote tracking is set correctly

### 2. During Fix Application  
- All edits in effort directory only
- All commits to effort branch only
- Push to effort branch remote only

### 3. After Fix Application
- Verify fixes exist in effort branch history
- Verify effort branch pushed to remote
- Verify NO changes in integration branch

## State-Specific Requirements

### ERROR_RECOVERY State
When fixes are identified during ERROR_RECOVERY:
1. Parse fix requirements from reports
2. Spawn SW Engineers with EXPLICIT instructions to work in effort branches
3. Monitor that fixes appear in effort branches
4. NEVER attempt fixes in integration branches

### FIX_ISSUES State (SW Engineer)
When entering FIX_ISSUES:
1. IMMEDIATELY verify you're in effort directory
2. IMMEDIATELY verify you're on effort branch
3. Apply ALL fixes to effort branch
4. Push effort branch to remote
5. Mark completion in effort directory

## Violations

### AUTOMATIC FAILURE (-100%)
- Any fix applied to integration branch
- Any commit made to integration branch for fixes
- Working in integration directory for fixes
- Not verifying effort branch before fixes

### MAJOR VIOLATIONS (-50%)
- Fixes not pushed to effort branch remote
- Missing verification of fix location
- Unclear instructions to engineers about fix location

## Related Rules
- R240: Integration Fix Execution Protocol
- R239: Fix Plan Distribution Protocol
- R271: Integration Testing Branch Creation
- R209: Effort Directory Isolation Protocol
- R197: One Agent Per Effort

## Grading Impact
- **Correct fix application to efforts**: +25% compliance bonus
- **Any fix in integration branch**: -100% AUTOMATIC FAILURE
- **Missing verification**: -30% violation
- **Fixes lost due to wrong branch**: -100% CRITICAL FAILURE

## Implementation Verification
```bash
verify_fix_application_location() {
    local phase=$1
    local wave=$2
    
    echo "🔍 Verifying fix application locations..."
    
    # Check integration branch for unwanted fixes
    cd "/efforts/phase${phase}/wave${wave}/integration"
    if git log --oneline --grep="^fix:" --since="1 hour ago" | grep -q "fix:"; then
        echo "🔴 CRITICAL: Fixes found in integration branch!"
        echo "🔴 This violates R299 - fixes must go to effort branches!"
        return 1
    fi
    
    # Check effort branches for required fixes
    for effort_dir in /efforts/phase${phase}/wave${wave}/effort-*/; do
        effort=$(basename "$effort_dir")
        cd "$effort_dir"
        
        if [ -f "FIX_COMPLETE.flag" ]; then
            if ! git log --oneline --grep="^fix:" --since="1 hour ago" | grep -q "fix:"; then
                echo "❌ $effort marked complete but no fixes in effort branch"
                return 1
            fi
            echo "✅ $effort has fixes in effort branch"
        fi
    done
    
    echo "✅ All fixes properly applied to effort branches"
    return 0
}
```

## Critical Reminder
**FIXES IN INTEGRATE_WAVE_EFFORTS BRANCHES ARE LOST!**
- Integration branches are temporary
- They are recreated from main each time
- Only fixes in effort branches persist
- This is why R299 is a SUPREME LAW

Without this rule, the system enters an infinite loop of:
1. Find issue during integration
2. Fix in integration branch
3. Create new integration branch (loses fix)
4. Find same issue again
5. Repeat forever

**ALWAYS FIX IN EFFORT BRANCHES!**