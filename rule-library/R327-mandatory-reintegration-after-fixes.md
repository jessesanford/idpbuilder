# 🔴🔴🔴 SUPREME RULE R327: Mandatory Re-Integration After Fixes

## Criticality: SUPREME LAW
**Violation = -100% AUTOMATIC FAILURE**

## Description
This rule mandates CASCADING RE-INTEGRATION: When ANY fixes are applied to source branches, ALL integration branches containing those sources MUST be deleted and recreated IN CASCADE ORDER. This is DETERMINISTIC and NON-NEGOTIABLE:
- Fix in effort branch → DELETE & RECREATE wave integration
- Recreated wave → DELETE & RECREATE phase integration  
- Recreated phase → DELETE & RECREATE project integration

This prevents the catastrophic problem of integration branches containing broken code while fixes exist only in upstream branches, leading to unbuildable binaries and false "completion" states.

## 🔴🔴🔴 THE ABSOLUTE CASCADE LAW 🔴🔴🔴

**AFTER ANY FIX AT ANY LEVEL, YOU MUST CASCADE ALL DEPENDENCIES AND RECREATE ALL CONTAINING INTEGRATIONS!**

### 🚨🚨🚨 COMPLETE CASCADE IS MANDATORY - NO EXCEPTIONS 🚨🚨🚨

**THE COMPLETE CASCADE RULE (AUTOMATIC FAILURE IF VIOLATED):**
1. **Effort fix** → All dependent efforts MUST be rebased (R350/R351)
2. **Dependent effort rebases** → Wave integration MUST be deleted & recreated
3. **Wave recreation** → Next wave's efforts MUST be rebased on new integration
4. **Next wave rebases** → Phase integration MUST be deleted & recreated  
5. **Phase recreation** → Next phase's efforts MUST be rebased
6. **All cascades complete** → Project integration MUST be deleted & recreated
7. **NO PARTIAL CASCADES** → Every dependency must be satisfied

**THIS IS DETERMINISTIC - ENFORCED BY R350 DEPENDENCY GRAPH AND R351 EXECUTION PROTOCOL!**

### 🔴 EFFORT-TO-EFFORT CASCADE REQUIREMENT 🔴

When fixes are applied to an effort branch, the cascade includes:
- **Within-wave dependencies**: effort2 depends on effort1, effort3 on effort2, etc.
- **Cross-wave dependencies**: Next wave efforts depend on current wave integration
- **Complete rebase chain**: Every dependent effort must be rebased in order

See R350 for dependency tracking and R351 for execution protocol.

### The Problem This Solves
```
❌ BROKEN FLOW (WHAT ALMOST HAPPENED TO YOU):
1. Wave integrations created at 03:24 and 17:53
2. Fixes applied to effort branches at 21:00
3. Orchestrator tries to use STALE wave integrations
4. Phase integration would have BROKEN CODE
5. Project would be UNBUILDABLE

✅ CORRECT CASCADE FLOW (MANDATORY):
1. Effort branch gets fix at 21:00
2. Wave integration from 03:24 is STALE → DELETE & RECREATE
3. Phase integration (if exists) is STALE → DELETE & RECREATE
4. Project integration (if exists) is STALE → DELETE & RECREATE
5. ALL integrations now contain fixes
6. Binary builds successfully
```

### 🔴 EXPLICIT CASCADE EXAMPLES 🔴

**Example 1: Complete Effort-to-Effort Cascade**
```
phase1/wave1/effort1 gets fix at 14:30
  ↓ CASCADE MANDATORY (R350/R351)
phase1/wave1/effort2 (depends on effort1) → REBASE on fixed effort1
  ↓ CASCADE MANDATORY
phase1/wave1/effort3 (depends on effort2) → REBASE on rebased effort2
  ↓ CASCADE MANDATORY
phase1/wave1/integration → DELETE & RECREATE with all rebased efforts
  ↓ CASCADE MANDATORY
phase1/wave2/effort1 (based on wave1/integration) → REBASE on new integration
  ↓ CASCADE MANDATORY
phase1/wave2/effort2 (depends on wave2/effort1) → REBASE on rebased effort1
  ↓ CASCADE MANDATORY
phase1/wave2/integration → DELETE & RECREATE
  ↓ CASCADE MANDATORY
phase1/integration → DELETE & RECREATE
  ↓ CASCADE MANDATORY
project/integration → DELETE & RECREATE
```

**Example 2: Multiple Efforts in Same Wave**
```
effort-2 and effort-3 get fixes
  ↓ CASCADE MANDATORY
effort-4 (depends on effort-3) → REBASE
effort-5 (depends on effort-4) → REBASE
  ↓ CASCADE MANDATORY
wave2-integration → DELETE & RECREATE (with ALL rebased efforts)
  ↓ CASCADE MANDATORY
wave3 efforts → REBASE on new wave2-integration
  ↓ CASCADE MANDATORY
phase-integration → DELETE & RECREATE
  ↓ CASCADE MANDATORY  
project-integration → DELETE & RECREATE
```

**Example 3: Cross-Wave Cascade Chain**
```
phase1/wave1/effort1 gets fix
  ↓ DEPENDENCY GRAPH (R350) CALCULATES:
  - wave1/effort2 depends on effort1
  - wave1/effort3 depends on effort2
  - wave2/effort1 depends on wave1/integration
  - wave2/effort2 depends on wave2/effort1
  - phase2/wave1/effort1 depends on phase1/integration
  
  ↓ EXECUTION PROTOCOL (R351) EXECUTES IN ORDER:
  1. REBASE wave1/effort2 on fixed effort1
  2. REBASE wave1/effort3 on rebased effort2
  3. RECREATE wave1/integration
  4. REBASE wave2/effort1 on new wave1/integration
  5. REBASE wave2/effort2 on rebased wave2/effort1
  6. RECREATE wave2/integration
  7. RECREATE phase1/integration
  8. REBASE phase2/wave1/effort1 on new phase1/integration
  9. Continue cascade...
```

## Stale Integration Tracking Mechanism

### Comprehensive Tracking Structure
The orchestrator-state.json MUST maintain detailed tracking of stale integrations:

```json
"stale_integration_tracking": {
  "stale_integrations": [...],     // List of stale integrations
  "staleness_cascade": [...],      // Cascade requirements
  "fix_tracking": {                // Comprehensive fix tracking
    "fixes_applied": [...],        // All fixes with integration status
    "fixes_pending_integration": [...]  // Fixes not yet integrated
  },
  "validation_history": [...]      // History of freshness checks
}
```

**See**: `/docs/STALE-INTEGRATION-TRACKING-MECHANISM.md` for full structure
**Example**: `/docs/STALE-TRACKING-EXAMPLE.md` for practical usage

## Core Requirements

### 1. MANDATORY CASCADE RE-INTEGRATION AT ALL LEVELS

#### 🔴🔴🔴 TIMESTAMP VALIDATION REQUIRED 🔴🔴🔴
```bash
# BEFORE ANY INTEGRATION, VALIDATE TIMESTAMPS:
validate_integration_timestamps() {
    local INTEGRATION_TYPE=$1  # wave, phase, project
    
    echo "🔍 R327 CASCADE VALIDATION: Checking timestamps"
    
    # Get integration creation time
    INTEGRATION_TIME=$(git log -1 --format=%ct "${INTEGRATION_TYPE}-integration")
    
    # Check ALL source branches for newer fixes
    SOURCES=$(get_source_branches_for "$INTEGRATION_TYPE")
    
    for source in $SOURCES; do
        LAST_FIX=$(git log -1 --grep="fix:" --format=%ct "$source" 2>/dev/null || echo 0)
        
        if [ "$LAST_FIX" -gt "$INTEGRATION_TIME" ]; then
            echo "❌ R327 CASCADE VIOLATION DETECTED!"
            echo "   Integration created: $(date -d @$INTEGRATION_TIME)"
            echo "   Fix applied to $source: $(date -d @$LAST_FIX)"
            echo "   🔴 MUST DELETE AND RECREATE INTEGRATION!"
            return 1
        fi
    done
    
    echo "✅ Integration is current (no fixes after creation)"
    return 0
}

# MANDATORY CHECK BEFORE USING ANY INTEGRATION
if ! validate_integration_timestamps "wave"; then
    echo "🔴🔴🔴 R327 CASCADE ENFORCEMENT ACTIVATED 🔴🔴🔴"
    delete_and_recreate_integration "wave"
    # CASCADE UP!
    delete_and_recreate_integration "phase"
    delete_and_recreate_integration "project"
fi
```

#### Wave Level Re-Integration (TRIGGERS CASCADE)
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
    
    echo "🔴 R327 CASCADE: Wave recreated, MUST recreate phase!"
    echo "🔴 THIS IS NOT OPTIONAL - CASCADE IS MANDATORY!"
    handle_phase_reintegration $PHASE  # CASCADE UP!
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
    
    echo "🔴 R327 CASCADE: Phase recreated, MUST recreate project!"
    echo "🔴 THIS IS NOT OPTIONAL - CASCADE IS MANDATORY!"
    handle_project_reintegration  # CASCADE UP!
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

#### 🔴🔴🔴 CASCADE_REINTEGRATION STATE - MANDATORY ENFORCEMENT 🔴🔴🔴

**When stale integrations are detected, the orchestrator MUST transition to CASCADE_REINTEGRATION state!**

This state:
- **BLOCKS** all other work until cascades complete
- **ENFORCES** proper recreation order (wave → phase → project)
- **TRACKS** cascade progress in state file
- **CANNOT BE SKIPPED** - it's a trap state until cascades complete

#### Required State Transitions After Fixes
```yaml
# CASCADE DETECTION - From ANY monitoring state:
MONITORING_INTEGRATION:
  stale_integrations_detected: true
  next_state: CASCADE_REINTEGRATION  # MANDATORY

WAVE_COMPLETE:
  fixes_applied_during_wave: true
  next_state: CASCADE_REINTEGRATION  # MANDATORY

PHASE_COMPLETE:
  pending_cascades_exist: true
  next_state: CASCADE_REINTEGRATION  # BLOCKS phase completion

# CASCADE EXECUTION:
CASCADE_REINTEGRATION:
  description: "Enforces cascade recreation of all stale integrations"
  transitions:
    - to: INTEGRATION  # For wave recreation
    - to: PHASE_INTEGRATION  # For phase recreation
    - to: PROJECT_INTEGRATION  # For project recreation
    - to: CASCADE_REINTEGRATION  # Loop until all cascades done
    - to: INTEGRATION_CODE_REVIEW  # When ALL cascades complete

# WAVE LEVEL - After fixes reviewed:
MONITOR_REVIEWS:
  all_fixes_reviewed: true
  next_state: WAVE_COMPLETE

WAVE_COMPLETE:
  description: "Marks wave ready for re-integration"
  checks_for_cascades: true  # R327 enforcement
  next_state: CASCADE_REINTEGRATION or INTEGRATION

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

### 3. MANDATORY VALIDATION GATES

#### 🔴🔴🔴 ORCHESTRATOR MUST-CHECK PROTOCOL 🔴🔴🔴
```bash
# ORCHESTRATOR MUST RUN THIS BEFORE ANY INTEGRATION USE:
orchestrator_cascade_check() {
    echo "🔍 R327 CASCADE CHECK: Validating all integrations"
    
    # Check 1: Are integrations newer than their sources?
    for integration in wave phase project; do
        if [ -d "${integration}-integration" ]; then
            INTEGRATION_TIME=$(git log -1 --format=%ct "${integration}-integration")
            
            # Find newest fix in any source
            NEWEST_FIX=0
            for source in $(get_sources_for "$integration"); do
                FIX_TIME=$(git log -1 --grep="fix:" --format=%ct "$source" 2>/dev/null || echo 0)
                [ "$FIX_TIME" -gt "$NEWEST_FIX" ] && NEWEST_FIX=$FIX_TIME
            done
            
            if [ "$NEWEST_FIX" -gt "$INTEGRATION_TIME" ]; then
                echo "❌ R327 CASCADE VIOLATION!"
                echo "   $integration-integration is STALE"
                echo "   Created: $(date -d @$INTEGRATION_TIME)"
                echo "   Newest fix: $(date -d @$NEWEST_FIX)"
                echo "   🔴 MUST CASCADE DELETE AND RECREATE!"
                return 1
            fi
        fi
    done
    
    # Check 2: Do integrations contain fix commits?
    echo "🔍 Checking for fix commit presence..."
    FIX_COMMITS=$(git log --all --grep="fix:" --since="6 hours ago" --format=%H)
    
    for commit in $FIX_COMMITS; do
        for integration in wave phase project; do
            if [ -d "${integration}-integration" ]; then
                if ! git log "${integration}-integration" --format=%H | grep -q "$commit"; then
                    echo "❌ R327 CASCADE VIOLATION!"
                    echo "   Fix commit $commit missing from ${integration}-integration"
                    echo "   🔴 MUST CASCADE DELETE AND RECREATE!"
                    return 1
                fi
            fi
        done
    done
    
    echo "✅ All integrations are current and contain fixes"
    return 0
}

# THIS MUST BE RUN BEFORE EVERY INTEGRATION OPERATION
if ! orchestrator_cascade_check; then
    echo "🔴🔴🔴 AUTOMATIC FAILURE TRIGGERED 🔴🔴🔴"
    echo "Attempting to use stale integrations = -100% GRADE"
    exit 1
fi
```

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

## Common Violations (ALL RESULT IN AUTOMATIC FAILURE)

### ❌ VIOLATION 1: Skipping Cascade Re-Integration
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

### ❌ VIOLATION 2: Incomplete Cascade
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

### ❌ VIOLATION 3: Using Stale Integration
```bash
# WRONG (WHAT ALMOST HAPPENED):
# Wave integration created at 03:24
# Fixes applied at 21:00
# Trying to use 03:24 integration for phase merge
cd phase-integration
git merge wave1-integration  # STALE! Missing fixes!
```

### ✅ CORRECTION 3: Check Timestamps First
```bash
# RIGHT:
# Check if wave integration is newer than fixes
WAVE_TIME=$(git log -1 --format=%ct wave1-integration)
FIX_TIME=$(git log -1 --grep="fix:" --format=%ct effort-branches)

if [ "$FIX_TIME" -gt "$WAVE_TIME" ]; then
    echo "🔴 Wave integration is STALE!"
    delete_and_recreate_wave_integration
    delete_and_recreate_phase_integration  # CASCADE!
    delete_and_recreate_project_integration  # CASCADE!
fi
```

### ❌ VIOLATION 4: Not Cascading Up
```bash
# WRONG:
cd project-integration
git pull  # Trying to "update" broken integration
```

### ✅ CORRECTION 4: Always Cascade Up
```bash
# RIGHT:
rm -rf project-integration
git push origin --delete project-integration
# Create fresh from main
```

## Grading Impact

### AUTOMATIC FAILURE (-100%)
- Using stale wave integration after effort fixes
- Using stale phase integration after wave recreation
- Using stale project integration after phase recreation
- Not cascading recreations up the hierarchy
- Skipping timestamp validation before integration use
- Proceeding with integration older than its source fixes
- Claiming completion with stale integrations

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

### Depends on
- R321: Immediate backports to source branches
- R350: Complete Cascade Dependency Graph (tracks all dependencies)
- R351: Cascade Execution Protocol (executes cascades in order)

### Enables
- R271: Fresh branches from main
- R348: Cascade State Transitions (uses CASCADE_REINTEGRATION)

### Critical for
- R266: Project bug detection (ensures bugs fixed in final integration)
- R328: Integration Freshness Validation (validates freshness)

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

**"CASCADE IS LAW - NO EXCEPTIONS"**
**"Fix in effort → Recreate wave → Recreate phase → Recreate project"**
**"Timestamp validation is MANDATORY before ANY integration use"**
**"Stale integration = AUTOMATIC FAILURE"**
**"When in doubt, CASCADE DELETE AND RECREATE"**

### 🔴🔴🔴 THE CASCADE MANTRA 🔴🔴🔴
```
A fix at any level,
Cascades to all above.
No integration left behind,
No stale branches to shove.

Check the timestamps first,
Validate commits are there.
If anything is stale,
CASCADE WITHOUT A CARE!
```

The goal: DETERMINISTIC CASCADE - Every fix triggers automatic recreation of ALL containing integrations. This is NOT optional, NOT subject to interpretation, and NOT negotiable.