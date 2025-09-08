# Orchestrator - MONITOR_FIXES State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.json with new state
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

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED MONITOR_FIXES STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_MONITOR_FIXES
echo "$(date +%s) - Rules read and acknowledged for MONITOR_FIXES" > .state_rules_read_orchestrator_MONITOR_FIXES
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## 📋 PRIMARY DIRECTIVES FOR MONITOR_FIXES STATE

### Core Mandatory Rules

### 🚨🚨🚨 R006 - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality**: BLOCKING - Automatic termination, 0% grade
**Summary**: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents

### 🚨🚨🚨 R319 - ORCHESTRATOR NEVER MEASURES CODE (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-code.md`
**Criticality**: BLOCKING - Orchestrator MUST NOT use line-counter.sh
**Summary**: Code Reviewers measure code size, NOT orchestrators

### ⚠️⚠️⚠️ R317 - Working Directory Restrictions (WARNING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R317-working-directory-restrictions.md`
**Criticality**: WARNING - -25% for violations
**Summary**: MUST NOT enter agent working directories - operate from root only

### 🚨🚨🚨 R318 - Agent Failure Escalation Protocol (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R318-agent-failure-escalation-protocol.md`
**Criticality**: BLOCKING - -40% for attempting forbidden fixes
**Summary**: NEVER fix agent failures directly - respawn with better instructions

### 🚨🚨🚨 R287 - TODO Save Frequency Requirements (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-save-frequency.md`
**Criticality**: BLOCKING - Save every 15 minutes/10 messages
**Summary**: Mandatory TODO saves during monitoring

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state.json on all state changes

### State-Specific Rules

### 🚨🚨🚨 R008 - Monitoring Frequency (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R008-monitoring-frequency.md`
**Criticality**: BLOCKING - Monitor every 5 messages/10 minutes
**Summary**: Continuous monitoring of all active agents

### 🚨🚨🚨 R254 - Agent Error Reporting (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R254-AGENT-ERROR-REPORTING.md`
**Criticality**: BLOCKING - Report and handle agent errors
**Summary**: Detect and report agent failures immediately

### 🚨🚨🚨 R255 - Post-Agent Work Verification (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R255-POST-AGENT-WORK-VERIFICATION.md`
**Criticality**: BLOCKING - Verify all work locations
**Summary**: Confirm agents worked in correct directories and branches

### 🔴🔴🔴 R321 - Immediate Backport During Integration (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-during-integration.md`
**Criticality**: SUPREME LAW - Integration fixes must be backported immediately
**Summary**: No deferred backporting - fix source branches immediately

## 🔴🔴🔴 CRITICAL: MONITOR_FIXES IS A VERB - START MONITORING FIXES NOW! 🔴🔴🔴

**MONITOR_FIXES MEANS ACTIVELY MONITORING FIX PROGRESS RIGHT NOW!**
- ❌ NOT "I'm in monitor fixes state"  
- ❌ NOT "Ready to monitor fixes"
- ✅ YES "I'm checking SW Engineer fixing E3.1.2 NOW"
- ✅ YES "I'm verifying fix progress in E3.1.3 NOW"
- ✅ YES "I'm detecting fix completion NOW"

## State Context
MONITOR_FIXES = You ARE ACTIVELY monitoring spawned SW Engineers fixing review issues or build failures THIS INSTANT. Focus specifically on fix progress and verification.

## 🚨🚨🚨 IMMEDIATE ACTIONS UPON ENTERING MONITOR_FIXES STATE 🚨🚨🚨

**THE INSTANT YOU ENTER MONITOR_FIXES STATE, DO THIS:**

```bash
# ✅ CORRECT - IMMEDIATE ACTION
echo "🔍 MONITORING FIXES: Checking all active SW Engineers fixing issues NOW..."

# Step 1: List all SW Engineers fixing issues (DO NOW!)
echo "📊 Active SW Engineers performing fixes:"
for effort in $(jq '.efforts_in_progress[].name' orchestrator-state.json); do
    FIX_STATUS=$(jq '.efforts_in_progress[] | select(.name == \"$effort\") | .fix_status' orchestrator-state.json)
    
    if [ "$FIX_STATUS" = "IN_PROGRESS" ]; then
        echo "  - $effort: Checking fix progress..."
        FIX_DIR="/efforts/phase${PHASE}/wave${WAVE}/${effort}"
        
        # Check for fix progress
        if [ -f "$FIX_DIR/FIX-PROGRESS.md" ]; then
            echo "    📄 Fix progress found - reading status..."
            cat "$FIX_DIR/FIX-PROGRESS.md"
        fi
    fi
done

# Step 2: Check for completed fixes (DO NOW!)
echo "🔍 Checking for completed fixes..."
for effort in $(jq '.efforts_in_progress[].name' orchestrator-state.json); do
    FIX_STATUS=$(jq '.efforts_in_progress[] | select(.name == \"$effort\") | .fix_status' orchestrator-state.json)
    FIX_DIR="/efforts/phase${PHASE}/wave${WAVE}/${effort}"
    
    if [ "$FIX_STATUS" = "IN_PROGRESS" ] && [ -f "$FIX_DIR/FIXES-COMPLETE.marker" ]; then
        echo "✅ Fixes complete for $effort!"
        jq '.efforts_in_progress[] |= select(.name == \"$effort\") |= .fix_status = \"COMPLETE\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
        
        # Determine if re-review needed
        echo "📋 Fixes complete - need re-review"
        echo "➡️ Must spawn Code Reviewer to verify fixes"
    fi
done

# Step 3: Check for backport requirements (R321) (DO NOW!)
echo "🔄 Checking for backport requirements (R321)..."
INTEGRATION_CONTEXT=$(jq '.integration_context // "none"' orchestrator-state.json)

if [ "$INTEGRATION_CONTEXT" = "active" ]; then
    echo "🔴🔴🔴 R321: INTEGRATION CONTEXT ACTIVE - IMMEDIATE BACKPORT REQUIRED!"
    
    for effort in $(jq '.efforts_in_progress[].name' orchestrator-state.json); do
        FIX_STATUS=$(jq '.efforts_in_progress[] | select(.name == \"$effort\") | .fix_status' orchestrator-state.json)
        
        if [ "$FIX_STATUS" = "COMPLETE" ]; then
            BACKPORTED=$(jq '.efforts_in_progress[] | select(.name == \"$effort\") | .backported // false' orchestrator-state.json)
            
            if [ "$BACKPORTED" != "true" ]; then
                echo "⚠️ $effort: Fixes complete but NOT BACKPORTED!"
                echo "🚨 R321 VIOLATION: Must backport immediately!"
                echo "➡️ Transitioning to IMMEDIATE_BACKPORT_REQUIRED"
                break
            fi
        fi
    done
fi

# Step 4: Check for blocked fixes (DO NOW!)
echo "🚧 Checking for blocked fix attempts..."
for effort in $(jq '.efforts_in_progress[].name' orchestrator-state.json); do
    FIX_DIR="/efforts/phase${PHASE}/wave${WAVE}/${effort}"
    
    if [ -f "$FIX_DIR/FIX-BLOCKED.marker" ]; then
        echo "⚠️ Fix blocked for $effort!"
        # Read blocker details
        if [ -f "$FIX_DIR/FIX-BLOCKER-DETAILS.md" ]; then
            cat "$FIX_DIR/FIX-BLOCKER-DETAILS.md"
        fi
        echo "➡️ May need architectural intervention"
    fi
done

# Step 5: Track fix types (DO NOW!)
echo "📊 Categorizing active fixes..."
REVIEW_FIXES=0
BUILD_FIXES=0
INTEGRATION_FIXES=0

for effort in $(jq '.efforts_in_progress[].name' orchestrator-state.json); do
    FIX_TYPE=$(jq '.efforts_in_progress[] | select(.name == \"$effort\") | .fix_type // \"unknown\"' orchestrator-state.json)
    
    case "$FIX_TYPE" in
        "review_issues") ((REVIEW_FIXES++)) ;;
        "build_failures") ((BUILD_FIXES++)) ;;
        "integration_conflicts") ((INTEGRATION_FIXES++)) ;;
    esac
done

echo "📈 Fix Statistics:"
echo "  - Review Issue Fixes: $REVIEW_FIXES"
echo "  - Build Failure Fixes: $BUILD_FIXES"
echo "  - Integration Conflict Fixes: $INTEGRATION_FIXES"

# Step 6: Determine next action (DO NOW!)
echo "🎯 Determining immediate next action based on fix monitoring..."
```

## Monitoring Fix Focus Areas

### 1. Fix Progress Tracking
Monitor SW Engineers actively fixing issues:
- Track commits addressing specific issues
- Verify fixes match review requirements
- Monitor test results after fixes
- Check build status improvements

### 2. Fix Type Management
Different fix types require different handling:

#### Review Issue Fixes:
- Addressing Code Reviewer findings
- Implementing requested changes
- Adding missing tests
- Fixing pattern violations
- **Next**: Re-review by Code Reviewer

#### Build Failure Fixes:
- Resolving compilation errors
- Fixing link issues
- Adding missing dependencies
- Correcting configurations
- **Next**: Re-run build validation

#### Integration Conflict Fixes:
- Resolving merge conflicts
- Fixing API incompatibilities
- Addressing version mismatches
- **CRITICAL (R321)**: Must backport immediately!

### 3. Backport Enforcement (R321)
**SUPREME LAW**: During integration, ALL fixes must be backported immediately:

```bash
# R321 Enforcement
if [ "$INTEGRATION_CONTEXT" = "active" ]; then
    for fix in completed_fixes; do
        if ! is_backported($fix); then
            echo "🔴 R321 VIOLATION DETECTED!"
            echo "Fix $fix not backported to source branch!"
            transition_to_state "IMMEDIATE_BACKPORT_REQUIRED"
            exit 1
        fi
    done
fi
```

### 4. Fix Completion Verification
When fixes complete:
1. **Verify all requested changes made**
2. **Tests pass with fixes**
3. **Build succeeds if applicable**
4. **Prepare for re-review**
5. **Track backport requirements**

### 5. Blocked Fix Handling
When SW Engineer reports blocked:
- Identify root cause of blockage
- Determine if architectural change needed
- Consider spawning Architect for guidance
- Document blocker for escalation

## State Transitions

From MONITOR_FIXES state:
- **FIXES_COMPLETE** → SPAWN_CODE_REVIEWERS_FOR_REVIEW (Re-review needed)
- **BACKPORT_REQUIRED** → IMMEDIATE_BACKPORT_REQUIRED (R321 enforcement)
- **FIX_BLOCKED** → ERROR_RECOVERY (Cannot proceed with fix)
- **BUILD_FIXES_COMPLETE** → BUILD_VALIDATION (Re-run build)
- **FIXES_ACTIVE** → MONITOR_FIXES (Continue monitoring)

## 🚨🚨🚨 R321 CRITICAL ENFORCEMENT 🚨🚨🚨

**During Integration Context:**
- **EVERY fix must be backported immediately**
- **NO deferred backporting allowed**
- **Integration branches are READ-ONLY for code**
- **Source branches must be fixed first**

**Violation Detection:**
```bash
# Check for R321 violation
if [[ "$INTEGRATION_ACTIVE" == "true" ]] && [[ "$FIX_NOT_BACKPORTED" == "true" ]]; then
    echo "🔴🔴🔴 R321 VIOLATION: Fix completed but not backported!"
    echo "MANDATORY: Transition to IMMEDIATE_BACKPORT_REQUIRED"
    exit 1
fi
```

## Fix Verification Requirements

Before marking fixes complete:
1. **Original Issues Resolved**
   - All review findings addressed
   - All build errors fixed
   - All conflicts resolved

2. **Quality Maintained**
   - Tests still pass
   - No new issues introduced
   - Code quality improved

3. **Documentation Updated**
   - Fix notes documented
   - Backport requirements noted
   - Review ready for re-check

## Success Criteria
- ✅ All SW Engineers fixing issues tracked
- ✅ Fix progress monitored continuously
- ✅ Completed fixes verified thoroughly
- ✅ Backport requirements enforced (R321)
- ✅ Re-reviews triggered when needed
- ✅ Blocked fixes escalated appropriately

## Failure Triggers
- ❌ Auto-transition without stop = R322 VIOLATION
- ❌ Miss fix completion = -50% penalty
- ❌ Skip backport during integration = R321 VIOLATION (-100%)
- ❌ Accept incomplete fixes = -40% penalty
- ❌ Fail to trigger re-review = -30% penalty

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**