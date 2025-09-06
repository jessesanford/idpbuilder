# 🔴🔴🔴 SUPREME RULE R321: Immediate Backport During Integration Protocol

## Criticality: SUPREME LAW
**Violation = -50% to -100% AUTOMATIC FAILURE**

## Description
This rule mandates that ANY fix discovered during wave or phase integration MUST be immediately backported to the source branches (effort branches) BEFORE continuing with integration. Integration branches are READ-ONLY for existing code - they can only receive merges, never direct edits. This prevents the catastrophic anti-pattern of broken source branches propagating through multiple integration stages.

## 🔴🔴🔴 THE ABSOLUTE REQUIREMENT 🔴🔴🔴

**INTEGRATION BRANCHES ARE READ-ONLY. ALL FIXES GO TO SOURCE BRANCHES FIRST.**

### The Problem This Solves
```
❌ CURRENT BROKEN FLOW (FORBIDDEN):
Effort1 (broken) ──┐
                   ├──> Wave Integration (broken) ──> Fix here only ──> Phase Integration (broken) ──> Fix here only
Effort2 (broken) ──┘                                        ↓                                              ↓
                                                    Effort branches                                 Effort branches
                                                    REMAIN BROKEN                                   STILL BROKEN

✅ REQUIRED CORRECT FLOW:
Effort1 (working) ──┐
                    ├──> Wave Integration (validation only) ──> Phase Integration (validation only)
Effort2 (working) ──┘
    ↑
    Fix applied HERE FIRST (immediately when found)
```

## Core Requirements

### 1. IMMEDIATE BACKPORT PROTOCOL (NO DEFERRALS)

#### When Integration Fails (Build/Test/Conflicts)
```bash
# THE MOMENT a fix is needed during integration:
handle_integration_issue() {
    local ISSUE_TYPE=$1  # conflict, build_fail, test_fail
    local AFFECTED_EFFORT=$2
    
    echo "🔴 STOP: Integration issue detected"
    echo "🔴 R321 MANDATES: Fix in source branch FIRST"
    
    # 1. IMMEDIATELY stop integration work
    git stash  # Save any pending integration work
    
    # 2. Identify source effort branch
    SOURCE_BRANCH="phase${PHASE}-wave${WAVE}-${AFFECTED_EFFORT}"
    
    # 3. Spawn SW Engineer to fix source branch
    cat > /tmp/immediate-backport-task.md << EOF
CRITICAL R321 IMMEDIATE BACKPORT REQUIRED

Issue Type: ${ISSUE_TYPE}
Source Branch: ${SOURCE_BRANCH}
Location: /efforts/phase${PHASE}/wave${WAVE}/${AFFECTED_EFFORT}

YOU MUST:
1. CD to effort directory
2. Checkout effort branch
3. Apply the fix
4. Test the fix locally
5. Commit and push
6. Mark completion with BACKPORT_COMPLETE.flag

This is BLOCKING integration - fix must be in source branch before integration can continue.
EOF
    
    spawn_sw_engineer FIX_ISSUES "$(cat /tmp/immediate-backport-task.md)"
    
    # 4. WAIT for backport completion
    wait_for_backport_completion "${AFFECTED_EFFORT}"
    
    # 5. ONLY THEN resume integration with fixed source
    echo "✅ Source branch fixed, resuming integration with working code"
}
```

### 2. VALIDATION GATES AT EVERY INTEGRATION LEVEL

#### Wave Integration Gate
```bash
# BEFORE marking wave integration complete:
validate_wave_integration_sources() {
    local PHASE=$1
    local WAVE=$2
    local ALL_SOURCES_WORKING=true
    
    echo "🔍 R321 VALIDATION: Verifying all source branches work independently"
    
    # Test EACH effort branch in isolation
    for effort_dir in /efforts/phase${PHASE}/wave${WAVE}/*/; do
        EFFORT_NAME=$(basename "$effort_dir")
        echo "Testing ${EFFORT_NAME} in isolation..."
        
        cd "$effort_dir"
        git checkout "phase${PHASE}-wave${WAVE}-${EFFORT_NAME}"
        
        # Branch must build independently
        if ! npm run build 2>&1 | tee "/tmp/${EFFORT_NAME}-build.log"; then
            echo "❌ ${EFFORT_NAME} doesn't build independently!"
            ALL_SOURCES_WORKING=false
        fi
        
        # Branch must pass its own tests
        if ! npm test 2>&1 | tee "/tmp/${EFFORT_NAME}-test.log"; then
            echo "❌ ${EFFORT_NAME} tests fail independently!"
            ALL_SOURCES_WORKING=false
        fi
    done
    
    if [ "$ALL_SOURCES_WORKING" = false ]; then
        echo "🔴🔴🔴 R321 VIOLATION: Source branches not independently working!"
        echo "Integration cannot proceed with broken source branches"
        exit 1
    fi
    
    echo "✅ All source branches validated as independently working"
}
```

#### Phase Integration Gate
```bash
# BEFORE marking phase integration complete:
validate_phase_integration_sources() {
    local PHASE=$1
    local ALL_WAVES_WORKING=true
    
    echo "🔍 R321 VALIDATION: Verifying all wave branches work independently"
    
    # Test EACH wave integration branch
    for wave_num in $(seq 1 10); do
        WAVE_DIR="/efforts/phase${PHASE}/wave${wave_num}"
        if [ ! -d "$WAVE_DIR" ]; then
            continue
        fi
        
        WAVE_BRANCH="phase${PHASE}-wave${wave_num}-integration"
        
        # Clone and test wave branch independently
        TEMP_DIR="/tmp/wave${wave_num}-validation"
        git clone --single-branch --branch "$WAVE_BRANCH" "$TARGET_REPO_URL" "$TEMP_DIR"
        
        cd "$TEMP_DIR"
        if ! npm run build; then
            echo "❌ Wave ${wave_num} doesn't build independently!"
            ALL_WAVES_WORKING=false
        fi
        
        if ! npm test; then
            echo "❌ Wave ${wave_num} tests fail!"
            ALL_WAVES_WORKING=false
        fi
        
        rm -rf "$TEMP_DIR"
    done
    
    if [ "$ALL_WAVES_WORKING" = false ]; then
        echo "🔴🔴🔴 R321 VIOLATION: Wave branches not working!"
        exit 1
    fi
}
```

### 3. INTEGRATION BRANCH MODIFICATION DETECTION

#### Automatic Detection and Rejection
```bash
# Monitor integration branches for direct modifications
detect_integration_modifications() {
    local INTEGRATION_BRANCH=$1
    
    # Check for any non-merge commits
    NON_MERGE_COMMITS=$(git log --oneline --no-merges origin/main..$INTEGRATION_BRANCH)
    
    if [ -n "$NON_MERGE_COMMITS" ]; then
        echo "🔴🔴🔴 R321 VIOLATION DETECTED!"
        echo "Direct modifications found in integration branch:"
        echo "$NON_MERGE_COMMITS"
        echo ""
        echo "Integration branches are READ-ONLY!"
        echo "These changes MUST be backported to source branches immediately!"
        
        # Reject the integration
        exit 1
    fi
}
```

### 4. STATE MACHINE ENFORCEMENT

#### New Required States After Integration Issues
```yaml
# When integration has issues, MUST transition through:
WAVE_INTEGRATION_ISSUE_DETECTED:
  next_states:
    - IMMEDIATE_BACKPORT_REQUIRED  # R321 Mandatory
    
IMMEDIATE_BACKPORT_REQUIRED:
  description: "R321 enforcement - fixing source branches"
  required_actions:
    - Identify affected effort branches
    - Spawn SW Engineers for each
    - Wait for all fixes in source
    - Verify source branches work
  next_states:
    - RETRY_WAVE_INTEGRATION  # Only after sources fixed
    
RETRY_WAVE_INTEGRATION:
  description: "Recreate integration with fixed sources"
  required_actions:
    - Delete old integration workspace
    - Create fresh integration from main
    - Merge all fixed source branches
    - Validate again
```

## Why This Is Critical

### The Cascade of Broken Branches (Current Problem)
```
1. Effort branches created (may have issues)
2. Wave integration created → Issues found → Fixed in integration only
3. Effort branches STILL BROKEN
4. Phase integration created → Pulls broken efforts → More issues
5. Fixed in phase integration only
6. Effort branches STILL BROKEN
7. Final integration → Everything breaks again
8. Massive rework because sources were never fixed
```

### The Clean Pipeline (With R321)
```
1. Effort branches created
2. Wave integration attempted → Issue found
3. STOP → Fix effort branch immediately
4. Retry integration with fixed effort → Success
5. Phase integration uses working wave branches → Success
6. No cascading failures, no deferred technical debt
```

## Enforcement Mechanisms

### 1. Integration Agent Restrictions
```bash
# Integration Agent init.sh MUST include:
echo "🔴 R321 ENFORCEMENT: This agent is READ-ONLY for code"
echo "I can only:"
echo "  ✅ Merge branches"
echo "  ✅ Resolve conflicts by choosing versions"
echo "  ✅ Run builds and tests"
echo "  ✅ Report issues"
echo "I CANNOT:"
echo "  ❌ Edit source code"
echo "  ❌ Fix bugs"
echo "  ❌ Modify logic"
echo "  ❌ Apply patches"
```

### 2. Orchestrator Monitoring
```bash
# Orchestrator MUST check after every integration:
monitor_integration_compliance() {
    # Check for unauthorized changes
    if git diff --name-only origin/main HEAD | grep -v "^MERGE"; then
        echo "🔴 R321 VIOLATION: Non-merge changes detected!"
        transition_to_state "IMMEDIATE_BACKPORT_REQUIRED"
    fi
}
```

### 3. Code Reviewer Validation
```bash
# Code Reviewer MUST verify during review:
review_integration_for_r321() {
    echo "🔍 R321 Compliance Check:"
    echo "- [ ] No direct code edits in integration branch"
    echo "- [ ] All source branches build independently"
    echo "- [ ] All source branches pass tests independently"
    echo "- [ ] Integration only contains merge commits"
}
```

## Common Violations and Corrections

### ❌ VIOLATION 1: Fixing Conflicts in Integration
```bash
# WRONG:
cd integration-workspace
git merge effort-1
# Conflict in api.go
vim api.go  # Manually "fixing" the code
git add api.go
git commit -m "Resolved conflicts"  # VIOLATION!
```

### ✅ CORRECTION 1: Fix in Source First
```bash
# RIGHT:
cd integration-workspace
git merge effort-1
# Conflict detected
echo "Conflict requires source fix"

# Go fix effort-1 branch FIRST
cd /efforts/phase1/wave1/effort-1
git checkout phase1-wave1-effort-1
vim api.go  # Fix the source
git commit -m "fix: resolve integration conflict"
git push

# THEN retry integration with fixed source
cd integration-workspace
git reset --hard HEAD~1  # Undo failed merge
git merge effort-1  # Now merges cleanly
```

### ❌ VIOLATION 2: Patching Build Failures
```bash
# WRONG:
cd integration-workspace
npm run build  # Fails
vim package.json  # "Quick fix" in integration
npm run build  # Works now
git commit -m "Fixed build"  # VIOLATION!
```

### ✅ CORRECTION 2: Fix Build in Source
```bash
# RIGHT:
cd integration-workspace
npm run build  # Fails
echo "Build failure from effort-2"

# Fix in source
spawn_sw_engineer FIX_ISSUES "Fix build in effort-2"
# Wait for completion
verify_effort_branch_builds effort-2

# Recreate integration with working source
rm -rf integration-workspace
create_fresh_integration
```

## Grading Impact

### AUTOMATIC FAILURE (-100%)
- Any code changes directly in integration branches
- Proceeding with broken source branches
- Deferring fixes to "later" (BACKPORT_FIXES state)
- Integration branches with non-merge commits

### MAJOR VIOLATIONS (-50%)
- Not validating source branches independently
- Incomplete backporting before retry
- Missing validation gates
- Allowing cascading broken branches

### COMPLIANCE BONUS (+20%)
- All sources validated before integration
- Immediate backporting on issue detection
- Clean integration branches (merges only)
- No deferred technical debt

## Relationship to Other Rules

### Strengthens R300
- R300: Fixes must go to effort branches (eventually)
- R321: Fixes must go to effort branches IMMEDIATELY (no deferral)

### Enables R271
- R271: Fresh integration branches from main
- R321: Ensures sources are ready for fresh integration

### Supports R006
- R006: Orchestrator never writes code
- R321: Enforces delegation for ALL fixes, even during integration

## Quick Reference Commands

### Check Integration Branch Compliance
```bash
# Is integration branch clean (merges only)?
git log --oneline --no-merges origin/main..HEAD
# Should return EMPTY
```

### Validate Source Branch Independence
```bash
# Can effort branch build alone?
cd /efforts/phase${P}/wave${W}/effort-${E}
git checkout phase${P}-wave${W}-effort-${E}
npm run build && npm test
# Must succeed
```

### Detect Deferred Fixes
```bash
# Are there fixes only in integration?
diff <(git log integration --oneline --grep="fix:") \
     <(git log effort-1 effort-2 --oneline --grep="fix:")
# Should show NO unique fixes in integration
```

## Remember

**"Fix at the source, not in the merge"**
**"Integration validates, never modifies"**
**"Broken sources create broken products"**
**"Today's integration hack is tomorrow's production bug"**

The goal: Every branch at every level should be independently functional, buildable, and testable. Integration branches merely prove that functional branches can work together.