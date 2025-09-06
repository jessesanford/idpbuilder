# Orchestrator - IMMEDIATE_BACKPORT_REQUIRED State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.yaml with new state
3. ✅ Committing and pushing the state file  
4. ✅ Providing work summary

### YOU MUST NOT:
- ❌ Continue to the next state automatically
- ❌ Start work for the new state
- ❌ Spawn agents for the new state
- ❌ Assume permission to continue

### STOP PROTOCOL:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: CURRENT_STATE → NEXT_STATE

### ✅ Current State Work Completed:
- [List completed work]

### 📊 Current Status:
- Current State: CURRENT_STATE
- Next State: NEXT_STATE
- TODOs Completed: X/Y
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to NEXT_STATE. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED IMMEDIATE_BACKPORT_REQUIRED STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R321 ENFORCEMENT STATE 🔴🔴🔴

**This state exists to enforce R321: Immediate Backport During Integration Protocol**
**You are here because an integration issue was detected that requires fixing in source branches BEFORE continuing.**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_IMMEDIATE_BACKPORT_REQUIRED
echo "$(date +%s) - Rules read and acknowledged for IMMEDIATE_BACKPORT_REQUIRED" > .state_rules_read_orchestrator_IMMEDIATE_BACKPORT_REQUIRED
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## 📋 PRIMARY DIRECTIVES FOR IMMEDIATE_BACKPORT_REQUIRED STATE

### 🔴🔴🔴 R321 - Immediate Backport During Integration Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-during-integration.md`
**Criticality**: SUPREME LAW - Immediate backporting required
**Summary**: ALL fixes must go to source branches IMMEDIATELY when found

### 🔴🔴🔴 R300 - Comprehensive Fix Management Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R300-comprehensive-fix-management-protocol.md`
**Criticality**: SUPREME LAW - Fixes go to effort branches
**Summary**: Effort branches are the source of truth for all fixes

### 🔴🔴🔴 R006 - Orchestrator Never Writes Code
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality**: SUPREME LAW - Orchestrator is manager only
**Summary**: ALL code changes must be done by SW Engineers

## 🚨 IMMEDIATE_BACKPORT_REQUIRED - SPAWN CODE REVIEWER FIRST! 🚨

### IMMEDIATE ACTIONS UPON ENTERING THIS STATE

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. **IDENTIFY** what fixes were made during integration
2. **DOCUMENT** all fixes and affected branches
3. **SPAWN CODE REVIEWER** to create backport plan
4. **TRANSITION** to WAITING_FOR_BACKPORT_PLAN
5. **WAIT** for Code Reviewer to complete plan
6. **THEN** spawn SW Engineers via proper flow

**NEW FLOW - CLEAR SEPARATION OF RESPONSIBILITIES:**
1. IMMEDIATE_BACKPORT_REQUIRED (document fixes)
2. → SPAWN_CODE_REVIEWER_BACKPORT_PLAN (spawn reviewer)
3. → WAITING_FOR_BACKPORT_PLAN (wait for plan)
4. → SPAWN_SW_ENGINEER_BACKPORT_FIXES (spawn engineers)
5. → MONITORING_BACKPORT_PROGRESS (monitor fixes)
6. → Back to integration retry with fixed sources

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in IMMEDIATE_BACKPORT_REQUIRED" [stops]
- ❌ "I'll handle the backporting myself" [orchestrator writes code]
- ❌ "Let me quickly fix this" [violates R006]
- ❌ "We can defer this until later" [violates R321]
- ❌ "Spawning SW Engineers directly" [skip Code Reviewer plan]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering IMMEDIATE_BACKPORT_REQUIRED, preparing to spawn Code Reviewer for backport planning..."
- ✅ "Creating fix instructions for affected efforts..."
- ✅ "Monitoring fix progress in source branches..."

## 📝 REQUIRED ACTIONS - NEW FLOW

### Step 1: Document Integration Fixes
```bash
# Gather all fix information
cd /efforts/integration-testing

# Create fix documentation for Code Reviewer
cat > FIX-MANIFEST-FOR-BACKPORT.md << 'EOF'
# Integration Fixes Requiring Immediate Backport

## Reason for IMMEDIATE_BACKPORT_REQUIRED
[Document why we entered this state - what failed/was fixed]

## Fixes Applied During Integration
[List all fixes that were made]

## Affected Effort Branches
[List which efforts need these fixes]

## Next Steps
Code Reviewer will analyze these fixes and create BACKPORT-PLAN.md
EOF
```

### Step 2: Transition to SPAWN_CODE_REVIEWER_BACKPORT_PLAN
```bash
# Update state to spawn Code Reviewer
cd $CLAUDE_PROJECT_DIR

cat > orchestrator-state.yaml << 'EOF'
current_state: SPAWN_CODE_REVIEWER_BACKPORT_PLAN
previous_state: IMMEDIATE_BACKPORT_REQUIRED
backport_trigger: "R321 - Immediate backport required"
fix_documentation: /efforts/integration-testing/FIX-MANIFEST-FOR-BACKPORT.md
next_action: "Spawn Code Reviewer to create backport plan"
EOF

git add orchestrator-state.yaml
git commit -m "state: transition to SPAWN_CODE_REVIEWER_BACKPORT_PLAN per R321"
git push
```

### CRITICAL: Follow the New Flow
The proper sequence is now:
1. **THIS STATE**: Document fixes
2. **SPAWN_CODE_REVIEWER_BACKPORT_PLAN**: Spawn reviewer
3. **WAITING_FOR_BACKPORT_PLAN**: Wait for plan
4. **SPAWN_SW_ENGINEER_BACKPORT_FIXES**: Spawn engineers
5. **MONITORING_BACKPORT_PROGRESS**: Monitor implementation

DO NOT skip directly to spawning SW Engineers!

## State Context

You are here because:
1. An integration attempt (wave or phase) encountered issues
2. These issues require fixes in source effort branches
3. R321 mandates these fixes happen IMMEDIATELY, not later
4. The integration cannot proceed until sources are fixed

## Entry Information Required

When entering this state, you should have:
- **Issue Report**: What failed during integration
- **Affected Efforts**: Which effort branches need fixes
- **Fix Requirements**: What needs to be fixed in each effort
- **Integration Context**: Wave or phase level integration

## Process Flow

### 1. Parse Integration Issues
```bash
# Identify what needs fixing from integration logs
parse_integration_issues() {
    echo "🔍 Analyzing integration failure..."
    
    # Check for build failures
    if [ -f "build-failure.log" ]; then
        echo "Build failures detected in:"
        grep "ERROR" build-failure.log | grep -o "effort-[^/]*"
    fi
    
    # Check for test failures  
    if [ -f "test-failure.log" ]; then
        echo "Test failures detected in:"
        grep "FAIL" test-failure.log | grep -o "effort-[^/]*"
    fi
    
    # Check for merge conflicts
    if git status | grep -q "Unmerged paths"; then
        echo "Merge conflicts in:"
        git status --porcelain | grep "^UU" | awk '{print $2}'
    fi
}
```

### 2. Create Fix Instructions
```bash
# Create detailed fix instructions for each effort
create_fix_instructions() {
    local EFFORT=$1
    local ISSUE_TYPE=$2
    local DETAILS=$3
    
    cat > "/tmp/fix-${EFFORT}.md" << EOF
# IMMEDIATE FIX REQUIRED - R321 ENFORCEMENT

## Effort: ${EFFORT}
## Issue Type: ${ISSUE_TYPE}
## Priority: BLOCKING INTEGRATION

### Problem Description
${DETAILS}

### Required Actions
1. CD to effort directory: /efforts/phase\${PHASE}/wave\${WAVE}/${EFFORT}
2. Checkout effort branch: phase\${PHASE}-wave\${WAVE}-${EFFORT}
3. Apply the following fixes:
   - [Specific fix instructions based on issue type]
4. Test locally to verify fix
5. Commit with message: "fix(${ISSUE_TYPE}): resolve integration blocking issue"
6. Push to remote
7. Create BACKPORT_COMPLETE.flag in effort directory

### Validation Required
- Branch must build successfully
- All tests must pass
- No linting errors
- Integration-specific requirements met

### R321 Compliance
This fix MUST be in the effort branch before integration can proceed.
The integration branch will be recreated fresh after this fix.
EOF
}
```

### 3. Spawn SW Engineers
```bash
# Spawn SW Engineer for each effort needing fixes
spawn_engineers_for_backports() {
    local EFFORTS_TO_FIX=("$@")
    
    for effort in "${EFFORTS_TO_FIX[@]}"; do
        echo "🚀 Spawning SW Engineer for ${effort}"
        
        TASK_FILE="/tmp/fix-${effort}.md"
        
        spawn_sw_engineer FIX_ISSUES "$(cat $TASK_FILE)" \
            --effort-name="${effort}" \
            --state="FIX_INTEGRATION_ISSUES" \
            --priority="BLOCKING"
            
        echo "✅ SW Engineer spawned for ${effort}"
    done
    
    echo "📊 Total engineers spawned: ${#EFFORTS_TO_FIX[@]}"
}
```

### 4. Monitor Fix Progress
```bash
# Monitor all SW Engineers fixing issues
monitor_backport_progress() {
    local ALL_COMPLETE=false
    
    while [ "$ALL_COMPLETE" = false ]; do
        echo "🔍 Checking fix progress..."
        ALL_COMPLETE=true
        
        for effort_dir in /efforts/phase${PHASE}/wave${WAVE}/*/; do
            EFFORT=$(basename "$effort_dir")
            
            if [ -f "${effort_dir}/BACKPORT_COMPLETE.flag" ]; then
                echo "✅ ${EFFORT}: Fix complete"
            else
                echo "⏳ ${EFFORT}: Still working..."
                ALL_COMPLETE=false
            fi
        done
        
        if [ "$ALL_COMPLETE" = false ]; then
            echo "Waiting 30 seconds before next check..."
            sleep 30
        fi
    done
    
    echo "✅ All fixes completed!"
}
```

### 5. Verify Fixed Sources
```bash
# R321 MANDATORY: Verify all source branches work independently
verify_all_sources_fixed() {
    local PHASE=$1
    local WAVE=$2
    local ALL_WORKING=true
    
    echo "🔍 R321 Final Verification: Testing all source branches"
    
    for effort_dir in /efforts/phase${PHASE}/wave${WAVE}/*/; do
        EFFORT=$(basename "$effort_dir")
        cd "$effort_dir"
        
        echo "Testing ${EFFORT}..."
        
        # Fetch latest
        git fetch origin
        git pull origin $(git branch --show-current)
        
        # Build test
        if ! npm run build 2>&1 | tee "/tmp/${EFFORT}-build-verify.log"; then
            echo "❌ ${EFFORT} still doesn't build!"
            ALL_WORKING=false
        fi
        
        # Test suite
        if ! npm test 2>&1 | tee "/tmp/${EFFORT}-test-verify.log"; then
            echo "❌ ${EFFORT} tests still fail!"
            ALL_WORKING=false
        fi
        
        # Check commit exists
        if ! git log --oneline -1 | grep -q "fix"; then
            echo "❌ ${EFFORT} missing fix commit!"
            ALL_WORKING=false
        fi
    done
    
    if [ "$ALL_WORKING" = false ]; then
        echo "🔴 Some sources still broken - cannot proceed!"
        exit 1
    fi
    
    echo "✅ All source branches verified working!"
    return 0
}
```

## State Transitions

From IMMEDIATE_BACKPORT_REQUIRED:
- **SUCCESS** → RETRY_INTEGRATION (retry with fixed sources)
- **PARTIAL_SUCCESS** → MONITORING_FIX_PROGRESS (some still working)
- **FAILURE** → ERROR_RECOVERY (fixes failed)

To IMMEDIATE_BACKPORT_REQUIRED:
- **INTEGRATION** → When integration issues detected
- **PHASE_INTEGRATION** → When phase integration issues detected
- **INTEGRATION_TESTING** → When integration tests fail

## Success Criteria

Before transitioning to RETRY_INTEGRATION:
1. ✅ All affected effort branches have fix commits
2. ✅ All effort branches build independently
3. ✅ All effort branches pass tests independently
4. ✅ All fixes pushed to remote
5. ✅ No fixes exist only in integration branches
6. ✅ Integration workspace ready for fresh retry

## Common Mistakes to Avoid

1. **Trying to fix as orchestrator**
   - ❌ WRONG: Orchestrator edits code
   - ✅ RIGHT: Spawn SW Engineers

2. **Fixing in integration branch**
   - ❌ WRONG: Apply fix to integration
   - ✅ RIGHT: Fix in source effort branch

3. **Deferring fixes**
   - ❌ WRONG: "We'll backport later"
   - ✅ RIGHT: Fix immediately before continuing

4. **Incomplete verification**
   - ❌ WRONG: Assume fixes work
   - ✅ RIGHT: Test each source independently

## R321 Compliance Checklist

- [ ] Identified all efforts needing fixes
- [ ] Created detailed fix instructions
- [ ] Spawned SW Engineers for each effort
- [ ] Monitored fix progress
- [ ] Verified each source builds
- [ ] Verified each source tests pass
- [ ] Confirmed fixes pushed to remote
- [ ] Ready to retry integration

## Remember

**"Fix at the source, immediately, always"**
**"Integration branches are temporary, effort branches are permanent"**
**"No progression with broken sources"**
**"The orchestrator manages, engineers implement"**

This state enforces the critical requirement that all fixes must be applied to source branches immediately when found, preventing the cascade of broken branches through integration stages.
## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
