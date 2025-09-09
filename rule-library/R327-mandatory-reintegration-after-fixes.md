# 🔴🔴🔴 SUPREME RULE R327: Mandatory Re-Integration After Fixes

## Criticality: SUPREME LAW
**Violation = -100% AUTOMATIC FAILURE**

## Description
This rule mandates that after ANY fixes are applied to source branches during integration failures, the ENTIRE integration MUST be deleted and re-run from scratch. This prevents the catastrophic problem of integration branches containing broken code while fixes exist only in upstream branches, leading to unbuildable binaries and false "completion" states.

## 🔴🔴🔴 THE ABSOLUTE REQUIREMENT 🔴🔴🔴

**AFTER FIXES, YOU MUST DELETE AND RE-CREATE THE ENTIRE INTEGRATION!**

### The Problem This Solves
```
❌ BROKEN FLOW (leads to unbuildable binaries):
1. Project integration merges all phases
2. Bugs found (e.g., wrong import paths)
3. Fixes applied to source branches (phase/wave/effort)
4. Integration branch STILL HAS BROKEN CODE
5. Binary cannot be built
6. Project appears "complete" but is broken

✅ CORRECT FLOW (ensures working binaries):
1. Project integration merges all phases
2. Bugs found
3. Fixes applied to source branches
4. DELETE old integration branch (has broken code)
5. CREATE fresh integration from main
6. RE-MERGE all branches (now with fixes)
7. Binary builds successfully
```

## Core Requirements

### 1. MANDATORY RE-INTEGRATION AT ALL LEVELS

#### Wave Level Re-Integration
```bash
# After wave fixes complete in effort branches:
handle_wave_reintegration() {
    local PHASE=$1
    local WAVE=$2
    
    echo "🔴 R327: Deleting broken wave integration"
    git push origin --delete "phase${PHASE}-wave${WAVE}-integration"
    rm -rf "/efforts/phase${PHASE}/wave${WAVE}/wave-integration"
    
    echo "✅ R327: Creating fresh wave integration"
    mkdir -p "/efforts/phase${PHASE}/wave${WAVE}/wave-integration"
    cd "/efforts/phase${PHASE}/wave${WAVE}/wave-integration"
    git clone "$REPO_URL" .
    git checkout -b "phase${PHASE}-wave${WAVE}-integration"
    
    echo "✅ R327: Re-merging all efforts with fixes"
    # Re-execute ENTIRE merge plan with fixed sources
}
```

#### Phase Level Re-Integration
```bash
# After phase fixes complete in wave branches:
handle_phase_reintegration() {
    local PHASE=$1
    
    echo "🔴 R327: Deleting broken phase integration"
    git push origin --delete "phase${PHASE}-integration"
    rm -rf "/efforts/phase${PHASE}/phase-integration"
    
    echo "✅ R327: Creating fresh phase integration"
    mkdir -p "/efforts/phase${PHASE}/phase-integration"
    cd "/efforts/phase${PHASE}/phase-integration"
    git clone "$REPO_URL" .
    git checkout -b "phase${PHASE}-integration"
    
    echo "✅ R327: Re-merging all waves with fixes"
    # Re-execute ENTIRE phase merge plan
}
```

#### Project Level Re-Integration
```bash
# After project fixes complete in phase branches:
handle_project_reintegration() {
    echo "🔴 R327: Deleting broken project integration"
    git push origin --delete "project-integration"
    rm -rf "/efforts/project-integration"
    
    echo "✅ R327: Creating fresh project integration"
    mkdir -p "/efforts/project-integration"
    cd "/efforts/project-integration"
    git clone "$REPO_URL" .
    git checkout -b "project-integration"
    
    echo "✅ R327: Re-merging all phases with fixes"
    # Re-execute ENTIRE project merge plan
}
```

### 2. STATE MACHINE ENFORCEMENT

#### Required State Transitions After Fixes
```yaml
# WAVE LEVEL - After fixes reviewed:
MONITOR_REVIEWS:
  all_fixes_reviewed: true
  next_state: WAVE_COMPLETE

WAVE_COMPLETE:
  description: "Marks wave ready for re-integration"
  next_state: INTEGRATION  # DELETE & RE-CREATE

# PHASE LEVEL - After fixes complete:
MONITORING_FIX_PROGRESS:
  phase_fixes_complete: true
  next_state: PHASE_INTEGRATION  # DELETE & RE-CREATE

# PROJECT LEVEL - After fixes complete:
MONITORING_PROJECT_FIXES:
  all_fixes_complete: true
  all_fixes_reviewed: true
  next_state: PROJECT_INTEGRATION  # DELETE & RE-CREATE
```

### 3. VALIDATION GATES

#### Pre-Re-Integration Validation
```bash
# BEFORE re-integration, verify ALL fixes are in source branches:
validate_fixes_in_sources() {
    local LEVEL=$1  # wave, phase, project
    
    echo "🔍 R327 Validation: Checking fixes in source branches"
    
    # List all fix commits
    FIX_COMMITS=$(git log --all --grep="fix:" --since="4 hours ago" --oneline)
    
    # Verify each fix is in a source branch, NOT integration
    while IFS= read -r commit; do
        BRANCHES=$(git branch -r --contains "${commit%% *}")
        
        if echo "$BRANCHES" | grep -q "integration"; then
            if ! echo "$BRANCHES" | grep -qv "integration"; then
                echo "❌ R327 VIOLATION: Fix only in integration branch!"
                echo "   Commit: $commit"
                echo "   This fix MUST be in source branch first!"
                exit 1
            fi
        fi
    done <<< "$FIX_COMMITS"
    
    echo "✅ All fixes verified in source branches"
}
```

#### Post-Re-Integration Validation
```bash
# AFTER re-integration, verify it builds:
validate_reintegration_success() {
    local INTEGRATION_DIR=$1
    
    cd "$INTEGRATION_DIR"
    
    echo "🔍 R327 Validation: Testing re-integrated code"
    
    # Must build successfully
    if ! make build; then
        echo "❌ R327 FAILURE: Re-integration still doesn't build!"
        echo "More fixes needed in source branches"
        return 1
    fi
    
    # Must pass tests
    if ! make test; then
        echo "❌ R327 FAILURE: Re-integration tests fail!"
        echo "More fixes needed in source branches"
        return 1
    fi
    
    echo "✅ R327 SUCCESS: Re-integration builds and tests pass!"
}
```

## Common Violations

### ❌ VIOLATION 1: Skipping Re-Integration
```bash
# WRONG:
MONITORING_PROJECT_FIXES → SPAWN_CODE_REVIEWER_PROJECT_VALIDATION
# Integration branch still has broken code!
```

### ✅ CORRECTION 1: Mandatory Re-Integration
```bash
# RIGHT:
MONITORING_PROJECT_FIXES → PROJECT_INTEGRATION → [full re-merge cycle]
# Fresh integration with all fixes
```

### ❌ VIOLATION 2: Partial Re-Integration
```bash
# WRONG:
# Only re-merging the "fixed" branches
git merge phase1-wave2-effort3  # Just the fixed one
```

### ✅ CORRECTION 2: Full Re-Integration
```bash
# RIGHT:
# Re-merge ALL branches in order (they may depend on each other)
for effort in effort1 effort2 effort3; do
    git merge "phase1-wave2-${effort}"
done
```

### ❌ VIOLATION 3: Keeping Old Integration
```bash
# WRONG:
cd project-integration
git pull  # Trying to "update" broken integration
```

### ✅ CORRECTION 3: Delete and Recreate
```bash
# RIGHT:
rm -rf project-integration
git push origin --delete project-integration
# Create fresh from main
```

## Grading Impact

### AUTOMATIC FAILURE (-100%)
- Skipping re-integration after fixes
- Proceeding to validation with broken integration branch
- Claiming project complete without re-integration
- Binary cannot be built due to skipped re-integration

### MAJOR VIOLATIONS (-50%)
- Partial re-integration (only some branches)
- Keeping old integration branch
- Not validating build after re-integration

### COMPLIANCE BONUS (+30%)
- Full re-integration at all levels
- Clean deletion and recreation
- Build validation after each re-integration
- Clear documentation of re-integration cycles

## Relationship to Other Rules

### Depends on R321
- R321: Immediate backports to source branches
- R327: Re-integration after those backports

### Enables R271
- R271: Fresh branches from main
- R327: Ensures fresh re-integration

### Critical for R266
- R266: Project bug detection
- R327: Ensures bugs actually get fixed in final integration

## Quick Reference

### Check if Re-Integration Needed
```bash
# Are there fixes not in integration?
INTEGRATION_COMMIT=$(git log integration -1 --format=%H)
FIXES_AFTER=$(git log --all --grep="fix:" --after="$INTEGRATION_COMMIT" --oneline)
if [ -n "$FIXES_AFTER" ]; then
    echo "🔴 R327: Re-integration required!"
fi
```

### Force Re-Integration
```bash
# Nuclear option - force complete re-integration
force_reintegration() {
    echo "🔴 R327 ENFORCEMENT: Forcing complete re-integration"
    
    # Delete ALL integration branches
    for branch in $(git branch -r | grep "integration"); do
        git push origin --delete "${branch#origin/}"
    done
    
    # Restart integration from scratch
    transition_to_state "PROJECT_INTEGRATION"
}
```

## Remember

**"Fixed sources, broken integration = Failed project"**
**"Delete the broken, create the fresh"**
**"Every fix requires full re-integration"**
**"If it doesn't build, it's not done"**

The goal: Every integration at every level contains the latest fixes from all source branches. Never proceed with broken integration branches.