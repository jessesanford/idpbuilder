# Orchestrator - MONITOR_REVIEWS State Rules

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

**YOU HAVE ENTERED MONITOR_REVIEWS STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_MONITOR_REVIEWS
echo "$(date +%s) - Rules read and acknowledged for MONITOR_REVIEWS" > .state_rules_read_orchestrator_MONITOR_REVIEWS
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## 📋 PRIMARY DIRECTIVES FOR MONITOR_REVIEWS STATE

### Core Mandatory Rules

### 🚨🚨🚨 R006 - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality**: BLOCKING - Automatic termination, 0% grade
**Summary**: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents

### 🚨🚨🚨 R319 - ORCHESTRATOR NEVER MEASURES CODE (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-code.md`
**Criticality**: BLOCKING - Orchestrator MUST NOT use line-counter.sh
**Summary**: Code Reviewers measure code size, NOT orchestrators

### 🚨🚨🚨 R338 - MANDATORY LINE COUNT STATE TRACKING (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R338-mandatory-line-count-state-tracking.md`
**Criticality**: BLOCKING - -50% per missing tracking, -100% if none
**Summary**: MUST capture line counts from Code Reviewer reports and update orchestrator-state.json
**Action Required**: When Code Reviewer completes, extract "Implementation Lines:" and update line_count_tracking

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

### 🚨🚨🚨 R222 - Mandatory Code Review Spawn (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R222-mandatory-code-review-spawn.md`
**Criticality**: BLOCKING - Must spawn reviewers for all complete implementations
**Summary**: WAVE_COMPLETE forbidden without all reviews

### 🚨🚨🚨 R302 - Split Tracking Requirements (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R302-split-tracking-requirements.md`
**Criticality**: BLOCKING - Track all split operations
**Summary**: Meticulously track split creation and progress

### 🔴🔴🔴 R204 - Orchestrator Creates Split Infrastructure (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R204-orchestrator-creates-split-infrastructure.md`
**Criticality**: SUPREME LAW - Only orchestrator creates split infrastructure
**Summary**: Create split infrastructure just-in-time when needed

## 🔴🔴🔴 CRITICAL: MONITOR_REVIEWS IS A VERB - START MONITORING REVIEWS NOW! 🔴🔴🔴

**MONITOR_REVIEWS MEANS ACTIVELY MONITORING CODE REVIEWS RIGHT NOW!**
- ❌ NOT "I'm in monitor reviews state"  
- ❌ NOT "Ready to monitor reviews"
- ✅ YES "I'm checking Code Reviewer CR-E3.1.2 progress NOW"
- ✅ YES "I'm reading review report for E3.1.3 NOW"
- ✅ YES "I'm detecting review failures NOW"

## State Context
MONITOR_REVIEWS = You ARE ACTIVELY monitoring spawned Code Reviewers performing reviews THIS INSTANT. Focus specifically on review progress, results, and split requirements.

## 🚨🚨🚨 IMMEDIATE ACTIONS UPON ENTERING MONITOR_REVIEWS STATE 🚨🚨🚨

**THE INSTANT YOU ENTER MONITOR_REVIEWS STATE, DO THIS:**

```bash
# ✅ CORRECT - IMMEDIATE ACTION
echo "🔍 MONITORING REVIEWS: Checking all active Code Reviewers NOW..."

# Step 1: List all Code Reviewers being monitored (DO NOW!)
echo "📊 Active Code Reviewers under monitoring:"
for effort in $(jq '.efforts_in_progress[].name' orchestrator-state.json); do
    REVIEW_STATUS=$(jq '.efforts_in_progress[] | select(.name == \"$effort\") | .review_status' orchestrator-state.json)
    
    if [ "$REVIEW_STATUS" = "IN_PROGRESS" ]; then
        echo "  - $effort: Checking review progress..."
        REVIEW_DIR="/efforts/phase${PHASE}/wave${WAVE}/${effort}"
        
        # Check for review report
        if [ -f "$REVIEW_DIR/CODE-REVIEW-REPORT.md" ]; then
            echo "    📄 Review report found - reading results..."
            # Parse review results
        fi
    fi
done

# Step 2: Check for completed reviews (DO NOW!)
echo "🔍 Checking for completed reviews..."
for effort in $(jq '.efforts_in_progress[].name' orchestrator-state.json); do
    REVIEW_STATUS=$(jq '.efforts_in_progress[] | select(.name == \"$effort\") | .review_status' orchestrator-state.json)
    REVIEW_DIR="/efforts/phase${PHASE}/wave${WAVE}/${effort}"
    
    if [ "$REVIEW_STATUS" = "IN_PROGRESS" ] && [ -f "$REVIEW_DIR/CODE-REVIEW-REPORT.md" ]; then
        # Parse review result
        REVIEW_RESULT=$(grep "^REVIEW_STATUS:" "$REVIEW_DIR/CODE-REVIEW-REPORT.md" | cut -d: -f2 | xargs)
        
        if [ "$REVIEW_RESULT" = "PASSED" ]; then
            echo "✅ Review PASSED for $effort!"
            jq '.efforts_in_progress[] |= select(.name == \"$effort\") |= .review_status = \"PASSED\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
            
        elif [ "$REVIEW_RESULT" = "FAILED" ]; then
            echo "❌ Review FAILED for $effort - fixes needed!"
            jq '.efforts_in_progress[] |= select(.name == \"$effort\") |= .review_status = \"FAILED\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
            echo "➡️ Must spawn SW Engineer for fixes"
            
        elif [ "$REVIEW_RESULT" = "NEEDS_SPLIT" ]; then
            echo "🔀 Review requires SPLIT for $effort!"
            jq '.efforts_in_progress[] |= select(.name == \"$effort\") |= .review_status = \"NEEDS_SPLIT\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
            
            # Check if split plan exists
            if [ -f "$REVIEW_DIR/SPLIT-INVENTORY.md" ]; then
                echo "📋 Split inventory found - initializing split tracking"
                
                # 🔴🔴🔴 CRITICAL: Initialize split_tracking with CORRECT count! 🔴🔴🔴
                # Count actual splits from SPLIT-INVENTORY.md
                ACTUAL_SPLITS=$(grep -c "^| [0-9]" "$REVIEW_DIR/SPLIT-INVENTORY.md" || echo 0)
                echo "📊 Detected $ACTUAL_SPLITS splits needed from SPLIT-INVENTORY.md"
                
                # Initialize split tracking for this effort
                jq ".split_tracking.\"$effort\" = {
                    \"total_splits\": $ACTUAL_SPLITS,
                    \"current_split\": 0,
                    \"splits\": [],
                    \"original_branch\": \"$(jq -r ".efforts_in_progress[] | select(.name == \"$effort\") | .branch" orchestrator-state.json)\",
                    \"status\": \"SPLIT_PLANNED\",
                    \"split_date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
                }" orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
                
                echo "✅ Split tracking initialized for $effort with $ACTUAL_SPLITS splits"
                echo "➡️ Transitioning to CREATE_NEXT_SPLIT_INFRASTRUCTURE"
            else
                echo "❌ ERROR: NEEDS_SPLIT but no SPLIT-INVENTORY.md found!"
                echo "   Code Reviewer must create split inventory first"
            fi
        fi
    fi
done

# Step 3: Check for size violations (DO NOW!)
echo "📏 Checking for size violations in reviews..."
for effort in $(jq '.efforts_in_progress[].name' orchestrator-state.json); do
    REVIEW_DIR="/efforts/phase${PHASE}/wave${WAVE}/${effort}"
    
    if [ -f "$REVIEW_DIR/CODE-REVIEW-REPORT.md" ]; then
        # Check for size violations (Code Reviewer measured, not us!)
        if grep -q "SIZE_VIOLATION" "$REVIEW_DIR/CODE-REVIEW-REPORT.md"; then
            echo "⚠️ Size violation detected in $effort!"
            echo "📊 Code Reviewer's measurement shows >800 lines"
            
            # Check for split inventory
            if [ -f "$REVIEW_DIR/SPLIT-INVENTORY.md" ]; then
                echo "✅ Split inventory created by Code Reviewer"
                
                # Initialize split_tracking if not already done
                if ! jq -e ".split_tracking.\"$effort\"" orchestrator-state.json > /dev/null 2>&1; then
                    echo "📊 Initializing split tracking from size violation detection"
                    
                    # Count actual splits from SPLIT-INVENTORY.md
                    ACTUAL_SPLITS=$(grep -c "^| [0-9]" "$REVIEW_DIR/SPLIT-INVENTORY.md" || echo 0)
                    echo "📊 Detected $ACTUAL_SPLITS splits needed from SPLIT-INVENTORY.md"
                    
                    # Initialize split tracking
                    jq ".split_tracking.\"$effort\" = {
                        \"total_splits\": $ACTUAL_SPLITS,
                        \"current_split\": 0,
                        \"splits\": [],
                        \"original_branch\": \"$(jq -r ".efforts_in_progress[] | select(.name == \"$effort\") | .branch" orchestrator-state.json)\",
                        \"status\": \"SPLIT_PLANNED\",
                        \"split_reason\": \"SIZE_VIOLATION detected by Code Reviewer\",
                        \"split_date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
                    }" orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
                    
                    echo "✅ Split tracking initialized for $effort with $ACTUAL_SPLITS splits"
                fi
                
                echo "➡️ Need to create split infrastructure"
            fi
        fi
    fi
done

# Step 4: Check all reviews complete (DO NOW!)
echo "🎯 Checking if all reviews complete..."
ALL_REVIEWS_DONE=true
for effort in $(jq '.efforts_in_progress[].name' orchestrator-state.json); do
    REVIEW_STATUS=$(jq '.efforts_in_progress[] | select(.name == \"$effort\") | .review_status' orchestrator-state.json)
    
    if [ "$REVIEW_STATUS" != "PASSED" ] && [ "$REVIEW_STATUS" != "FAILED" ] && [ "$REVIEW_STATUS" != "NEEDS_SPLIT" ]; then
        ALL_REVIEWS_DONE=false
        echo "  - $effort: Review still in progress"
    fi
done

if [ "$ALL_REVIEWS_DONE" = true ]; then
    echo "✅ All reviews complete!"
    # Determine next state based on results
fi

# Step 5: Determine next action (DO NOW!)
echo "🎯 Determining immediate next action based on review monitoring..."
```

## Monitoring Review Focus Areas

### 1. Review Progress Tracking
Monitor Code Reviewers actively reviewing:
- Check for CODE-REVIEW-REPORT.md creation
- Monitor review completion status
- Track time spent on reviews
- Verify thorough analysis being done

### 2. Review Result Processing
**CRITICAL**: Process review results immediately:

#### PASSED Reviews:
- Update review_status to PASSED
- Check if all reviews complete
- Prepare for WAVE_COMPLETE if all passed

#### FAILED Reviews:
- Update review_status to FAILED
- Prepare to spawn SW Engineer for fixes
- Document required fixes from report
- Transition to SPAWN_ENGINEERS_FOR_FIXES

#### NEEDS_SPLIT Reviews:
- Update review_status to NEEDS_SPLIT
- Verify SPLIT-INVENTORY.md exists
- Verify SPLIT-PLAN-*.md files exist
- Transition to CREATE_NEXT_SPLIT_INFRASTRUCTURE

### 3. Split Detection and Handling (R302)
When Code Reviewer reports NEEDS_SPLIT:
1. **Verify split plans created** by Code Reviewer
2. **Track in split_tracking** section of state file
3. **Transition to CREATE_NEXT_SPLIT_INFRASTRUCTURE** (R204)
4. **Orchestrator creates infrastructure** (not Code Reviewer!)
5. **Then spawn SW Engineer** for sequential splits

### 4. All Reviews Complete Detection
Continuously check if all reviews done:
```bash
# Check if ready for WAVE_COMPLETE
ALL_PASSED=true
ANY_FAILED=false
ANY_NEEDS_SPLIT=false

for effort in $(jq '.efforts_in_progress[].name' orchestrator-state.json); do
    STATUS=$(jq '.efforts_in_progress[] | select(.name == \"$effort\") | .review_status' orchestrator-state.json)
    
    case "$STATUS" in
        "PASSED") ;;
        "FAILED") ANY_FAILED=true; ALL_PASSED=false ;;
        "NEEDS_SPLIT") ANY_NEEDS_SPLIT=true; ALL_PASSED=false ;;
        *) ALL_PASSED=false ;;
    esac
done

if [ "$ALL_PASSED" = true ]; then
    echo "🎉 All reviews PASSED - ready for WAVE_COMPLETE!"
elif [ "$ANY_FAILED" = true ]; then
    echo "🔧 Some reviews failed - need fixes"
elif [ "$ANY_NEEDS_SPLIT" = true ]; then
    echo "🔀 Some reviews need splits - create infrastructure"
fi
```

## State Transitions

From MONITOR_REVIEWS state:
- **REVIEW_FAILED** → SPAWN_ENGINEERS_FOR_FIXES (SW Engineer must fix issues)
- **NEEDS_SPLIT** → CREATE_NEXT_SPLIT_INFRASTRUCTURE (Create split infrastructure)
- **ALL_REVIEWS_PASSED** → WAVE_COMPLETE (All reviews successful)
- **REVIEWS_ACTIVE** → MONITOR_REVIEWS (Continue monitoring)
- **REVIEW_COMPLETE_MORE_PENDING** → MONITOR_REVIEWS (Some done, others ongoing)

## 🚨🚨🚨 CRITICAL SPLIT HANDLING (R204 + R302) 🚨🚨🚨

**When Code Reviewer creates split plans:**

1. **Code Reviewer Creates** (in too-large branch):
   - SPLIT-INVENTORY.md
   - SPLIT-PLAN-001.md, SPLIT-PLAN-002.md, etc.

2. **Orchestrator Detects** (MONITOR_REVIEWS):
   - Reads NEEDS_SPLIT status
   - Verifies split files exist
   - Updates split_tracking in state

3. **Orchestrator Creates Infrastructure** (CREATE_NEXT_SPLIT_INFRASTRUCTURE):
   - Creates split directories
   - Creates split branches
   - Copies split plans
   - **NEVER delegates this to Code Reviewer!**

4. **Orchestrator Spawns SW Engineer** (SPAWN_AGENTS):
   - Sequential implementation of splits
   - One split at a time
   - Each split gets own review

## Success Criteria
- ✅ All Code Reviewers tracked continuously
- ✅ Review results processed immediately
- ✅ Failed reviews trigger fix spawns
- ✅ Split needs trigger infrastructure creation
- ✅ All reviews complete before WAVE_COMPLETE
- ✅ Split tracking meticulously maintained

## Failure Triggers
- ❌ Auto-transition without stop = R322 VIOLATION
- ❌ Miss review completion = -50% penalty
- ❌ Allow WAVE_COMPLETE without all reviews = R222 VIOLATION
- ❌ Create split infrastructure yourself = R204 VIOLATION
- ❌ Fail to track splits = R302 VIOLATION

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**

### 🔴🔴🔴 MANDATORY VALIDATION REQUIREMENT 🔴🔴🔴

**Per R288 and R324**: ALL state file updates MUST be validated before commit:

```bash
# After ANY update to orchestrator-state.json:
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state.json || {
    echo "❌ State file validation failed!"
    exit 288
}
```

**Use helper functions for automatic validation:**
```bash
# Source the helper functions
source "$CLAUDE_PROJECT_DIR/utilities/state-file-update-functions.sh"

# Use safe functions that include validation:
safe_state_transition "NEW_STATE" "reason"
safe_update_field "field_name" "value"
```
