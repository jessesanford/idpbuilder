# ⚠️ DEPRECATED - SUPERSEDED BY R300 ⚠️
**This rule has been consolidated into R300: Comprehensive Fix Management Protocol**
**Please use R300 for all fix management requirements**

# 🔴🔴🔴 RULE R298: Fix Backporting Verification Protocol (SUPREME LAW)

## Classification
- **Category**: Integration Management
- **Criticality Level**: 🔴🔴🔴 SUPREME LAW
- **Enforcement**: AUTOMATIC VERIFICATION REQUIRED
- **Penalty**: -100% for proceeding without verification

## The Rule

**EVERY fix made during ERROR_RECOVERY or integration failure MUST be verified to exist in the original effort branch before ANY re-integration attempt.**

## 🔴🔴🔴 SUPREME REQUIREMENT - NO EXCEPTIONS 🔴🔴🔴

### MANDATORY VERIFICATION PROTOCOL:
```bash
# BEFORE ANY RE-INTEGRATE_WAVE_EFFORTS, YOU MUST:
verify_fixes_in_effort_branches() {
    local PHASE=$1
    local WAVE=$2
    local VERIFICATION_FAILED=false
    
    echo "🔍 VERIFYING ALL FIXES ARE IN EFFORT BRANCHES (R298)"
    
    # For each effort that had fixes
    for effort in $(jq '.efforts_in_progress | keys | .[]' orchestrator-state-v3.json); do
        EFFORT_BRANCH="phase${PHASE}-wave${WAVE}-${effort}"
        
        # Check if effort had fixes applied
        FIX_COMMIT=$(git log origin/${EFFORT_BRANCH} --oneline --grep="fix:" --since="1 hour ago" | head -1)
        
        if [ -z "$FIX_COMMIT" ]; then
            echo "❌ CRITICAL: No fix commits found in ${EFFORT_BRANCH}!"
            echo "   Effort ${effort} was supposed to have fixes but branch has no fix commits"
            VERIFICATION_FAILED=true
        else
            echo "✅ Found fix in ${EFFORT_BRANCH}: ${FIX_COMMIT}"
            
            # Verify fix is pushed to remote
            git fetch origin ${EFFORT_BRANCH}
            LOCAL_SHA=$(git rev-parse ${EFFORT_BRANCH})
            REMOTE_SHA=$(git rev-parse origin/${EFFORT_BRANCH})
            
            if [ "$LOCAL_SHA" != "$REMOTE_SHA" ]; then
                echo "❌ CRITICAL: Fix not pushed to remote ${EFFORT_BRANCH}!"
                echo "   Local: ${LOCAL_SHA}"
                echo "   Remote: ${REMOTE_SHA}"
                VERIFICATION_FAILED=true
            fi
        fi
    done
    
    if [ "$VERIFICATION_FAILED" = true ]; then
        echo "🔴🔴🔴 VERIFICATION FAILED - CANNOT PROCEED WITH INTEGRATE_WAVE_EFFORTS 🔴🔴🔴"
        echo "Fixes are missing from effort branches!"
        echo "This violates R298 - Fix Backporting Verification Protocol"
        exit 1
    fi
    
    echo "✅ All fixes verified in effort branches - safe to proceed"
}
```

## Why This Is A Supreme Law

### THE PROBLEM THIS SOLVES:
1. **Fixes applied to integration branches get lost** when creating new integration branches
2. **Same issues recur** in subsequent integration attempts
3. **Effort branches remain broken** causing perpetual integration failures
4. **Infinite fix loops** waste time and resources

### THE SOLUTION:
1. **Verify fixes exist in effort branches** before re-integration
2. **Block integration** if fixes are missing
3. **Force proper fix propagation** to source branches
4. **Ensure fixes persist** across integration attempts

## Enforcement Points

### 1. SW Engineer - After Applying Fixes
```bash
# SW Engineer MUST verify fix is in effort branch
verify_my_fixes() {
    echo "🔍 Verifying my fixes are in effort branch..."
    
    # Check current branch
    CURRENT_BRANCH=$(git branch --show-current)
    if [[ "$CURRENT_BRANCH" == *"integration"* ]]; then
        echo "❌ CRITICAL: You're in integration branch! Fixes must be in effort branch!"
        exit 1
    fi
    
    # Verify fix commit exists
    FIX_COMMIT=$(git log --oneline --grep="fix:" -1)
    if [ -z "$FIX_COMMIT" ]; then
        echo "❌ No fix commit found!"
        exit 1
    fi
    
    echo "✅ Fix commit: $FIX_COMMIT"
    
    # Push to remote
    git push origin $CURRENT_BRANCH
    echo "✅ Fix pushed to remote effort branch"
}
```

### 2. Orchestrator - Before Re-Integration
```bash
# Orchestrator MUST verify all fixes before spawning integration
spawn_integration_after_fixes() {
    # R298 MANDATORY CHECK
    verify_fixes_in_effort_branches $PHASE $WAVE || {
        echo "❌ R298 VIOLATION: Cannot spawn integration without verified fixes"
        transition_to_state "ERROR_RECOVERY"
        exit 1
    }
    
    # Only proceed if verification passes
    spawn_integration_agent
}
```

### 3. Integration Agent - Before Creating Integration Branch
```bash
# Integration Agent MUST verify fixes before integration
create_integration_branch() {
    # R298 CHECK: Verify all effort branches have fixes
    for effort_branch in $EFFORT_BRANCHES; do
        echo "Checking $effort_branch for recent fixes..."
        
        git fetch origin $effort_branch
        RECENT_FIXES=$(git log origin/$effort_branch --oneline --grep="fix:" --since="2 hours ago")
        
        if [ -n "$RECENT_FIXES" ]; then
            echo "✅ Found fixes in $effort_branch:"
            echo "$RECENT_FIXES"
        fi
    done
    
    # Create fresh integration branch from main (R271)
    git checkout main
    git pull origin main
    git checkout -b integration-phase${PHASE}-wave${WAVE}-$(date +%s)
    
    # Merge effort branches WITH their fixes
    for effort_branch in $EFFORT_BRANCHES; do
        git merge origin/$effort_branch --no-ff -m "integrate: $effort_branch with fixes"
    done
}
```

## Verification Checklist

Before ANY re-integration attempt, verify:
- ✅ Each effort branch has fix commits
- ✅ Fix commits are pushed to remote
- ✅ Remote branches are up-to-date
- ✅ No fixes exist only in integration branches
- ✅ Integration branch will be created fresh from main
- ✅ All effort branches will be merged with their fixes

## Grading Impact

### AUTOMATIC FAILURES (-100%):
- Creating integration branch without verifying fixes
- Fixes exist only in integration branch
- Claiming fixes complete but effort branches unchanged
- Re-integration attempts without fix verification

### MAJOR VIOLATIONS (-50%):
- Incomplete fix verification
- Missing push to remote
- Partial fix application

## Common Violations

### ❌ VIOLATION 1: Fixing in Integration Branch
```bash
# WRONG - Fix applied to integration branch
git checkout integration-wave-1
vim src/broken.ts
git commit -m "fix: resolve integration issue"
# Effort branch still broken!
```

### ❌ VIOLATION 2: Not Pushing to Remote
```bash
# WRONG - Fix only in local effort branch
git checkout feature/effort-1
vim src/broken.ts
git commit -m "fix: resolve issue"
# Forgot to push! Remote doesn't have fix
```

### ❌ VIOLATION 3: Proceeding Without Verification
```bash
# WRONG - No verification before re-integration
# Just assuming fixes are there
spawn_integration_agent
# Fixes might be missing!
```

## Correct Implementation

### ✅ CORRECT: Full Fix Propagation
```bash
# 1. Fix in effort branch
git checkout feature/effort-1
vim src/broken.ts
git add -A
git commit -m "fix: resolve integration failure"
git push origin feature/effort-1

# 2. Verify fix is in remote
git fetch origin
git log origin/feature/effort-1 --oneline -1

# 3. Only then proceed with integration
verify_fixes_in_effort_branches
spawn_integration_agent
```

## Related Rules
- R292: Integration Fixes in Effort Branches (BLOCKING)
- R271: Single Branch Full Checkout (BLOCKING)
- R240: Integration Fix Execution Protocol (BLOCKING)
- R291: Integration Demo Requirement (BLOCKING)

## Remember

**"No fix left behind - Every fix must reach its home branch"**
**"Verify before you integrate - Assume nothing"**
**"The effort branch is the source of truth"**

Integration without verification is GUARANTEED to fail again!