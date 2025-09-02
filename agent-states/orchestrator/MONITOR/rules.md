# Orchestrator - MONITOR State Rules

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED MONITOR STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_MONITOR
echo "$(date +%s) - Rules read and acknowledged for MONITOR" > .state_rules_read_orchestrator_MONITOR
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY MONITOR WORK UNTIL RULES ARE READ:
- ❌ Start check agent progress
- ❌ Start monitor size limits
- ❌ Start track implementation status
- ❌ Start collect metrics
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R151, R208, R053..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all MONITOR rules"
   (YOU Must READ AND ACKNOWLEDGE EACH rule individually)
   ```

3. **Silent Reading**:
   ```
   ❌ WRONG: [Reads rules but doesn't acknowledge]
   "Now I've read the rules, let me start work..."
   (MUST explicitly acknowledge EACH rule)
   ```

4. **Reading From Memory**:
   ```
   ❌ WRONG: "I know R208 requires CD before spawn..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR MONITOR:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute MONITOR work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY MONITOR work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute MONITOR work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with MONITOR work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY MONITOR work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 📋 PRIMARY DIRECTIVES FOR MONITOR STATE

### 🚨🚨🚨 R008 - Monitoring Frequency
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R008-monitoring-frequency.md`
**Criticality**: BLOCKING - Monitor every 5 messages/10 minutes
**Summary**: Continuous monitoring of all active agents

### 🚨🚨🚨 R254 - Agent Error Reporting
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R254-AGENT-ERROR-REPORTING.md`
**Criticality**: BLOCKING - Report and handle agent errors
**Summary**: Detect and report agent failures immediately

### 🚨🚨🚨 R255 - Post-Agent Work Verification
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R255-POST-AGENT-WORK-VERIFICATION.md`
**Criticality**: BLOCKING - Verify all work locations
**Summary**: Confirm agents worked in correct directories and branches

### 🔴🔴🔴 R021 - Orchestrator Never Stops (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R021-orchestrator-never-stops.md`
**Criticality**: SUPREME LAW - Violation = -100% failure
**Summary**: Continue monitoring until all agents complete

### 🚨🚨🚨 R287 - TODO Save Frequency Requirements
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-save-frequency.md`
**Criticality**: BLOCKING - Save every 15 minutes/10 messages
**Summary**: Mandatory TODO saves during monitoring

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state.yaml on all state changes

## 🔴🔴🔴 CRITICAL: MONITOR IS A VERB - START MONITORING IMMEDIATELY! 🔴🔴🔴

**MONITOR MEANS ACTIVELY MONITORING RIGHT NOW!**
- ❌ NOT "I'm in monitor state"  
- ❌ NOT "Ready to monitor"
- ❌ NOT "Monitoring mode activated"
- ✅ YES "I'm checking agent E3.1.2 status NOW"
- ✅ YES "I'm verifying E3.1.3 line count NOW"
- ✅ YES "I'm detecting blocking conditions NOW"

## State Context
MONITOR = You ARE ACTIVELY monitoring spawned agents THIS INSTANT. Not preparing to monitor, not ready to monitor, but MONITORING NOW!

## 🚨🚨🚨 IMMEDIATE ACTIONS UPON ENTERING MONITOR STATE 🚨🚨🚨

**THE INSTANT YOU ENTER MONITOR STATE, DO THIS:**

```bash
# ✅ CORRECT - IMMEDIATE ACTION
echo "🔍 MONITORING: Checking all spawned agents NOW..."

# Step 1: List all agents being monitored (DO NOW!)
echo "📊 Active agents under monitoring:"
for agent in "${SPAWNED_AGENTS[@]}"; do
    echo "  - $agent: checking status..."
    check_agent_status "$agent"
done

# Step 2: Check for completed agents (DO NOW!)
echo "🔍 Checking for completed efforts..."
check_completed_efforts

# Step 3: Check for violations (DO NOW!)
echo "⚠️ Checking for size violations..."
# CRITICAL: Each effort is a SEPARATE git repository!
# Must CD into effort dir and use BRANCH NAMES not directory names!
for effort in $(yq '.efforts_in_progress[].name' orchestrator-state.yaml); do
    EFFORT_DIR="/efforts/phase${PHASE}/wave${WAVE}/${effort}"
    if [ -d "$EFFORT_DIR" ]; then
        cd "$EFFORT_DIR"
        # Get ACTUAL branch name (not directory name!)
        CURRENT_BRANCH=$(git branch --show-current)
        BASE_BRANCH="phase${PHASE}/integration"  # Or check git branch -a
        # Run line-counter with BRANCH names
        ../../tools/line-counter.sh -b "$BASE_BRANCH" -c "$CURRENT_BRANCH"
    fi
done

# Step 4: Check for blocked agents (DO NOW!)
echo "🚧 Checking for blocked agents..."
detect_blocking_conditions

# Step 5: Determine next action (DO NOW!)
echo "🎯 Determining immediate next action based on monitoring..."
determine_next_action_from_monitoring
```

**❌❌❌ VIOLATIONS THAT CAUSE AUTOMATIC FAILURE:**

```bash
# ❌ CATASTROPHIC - Stopping after transition
transition_to_state "MONITOR"
echo "STATE TRANSITION COMPLETE: Now in MONITOR State"
# [stops here] - AUTOMATIC FAILURE!

# ❌ CATASTROPHIC - Announcing state without action
echo "Successfully entered MONITOR state"
echo "Ready to begin monitoring when needed..."
# NO! Start monitoring NOW!

# ❌ CATASTROPHIC - Summarizing instead of monitoring
echo "📊 Summary: We're now monitoring agents"
echo "Previous actions completed successfully"
# STOP TALKING, START MONITORING!
```

## Monitoring Implementation

### 🚨🚨🚨 MANDATORY: Spawn Code Reviewers When Implementation Complete 🚨🚨🚨

**CRITICAL R222 ENFORCEMENT - YOU WILL BE GRADED ON THIS!**

When MONITOR detects ANY effort has:
- implementation_status: COMPLETE
- review_status: NOT_STARTED or null or undefined

**YOU MUST IMMEDIATELY:**
1. **STOP** monitoring that effort
2. **SPAWN** Code Reviewer for that specific effort 
3. **UPDATE** review_status to IN_PROGRESS in state file
4. **TRACK** the spawned reviewer in state file
5. **CONTINUE** monitoring other efforts

**VIOLATION = -100% GRADE**: Allowing WAVE_COMPLETE without spawning reviewers

#### Implementation Complete Detection Logic:
```bash
# CRITICAL: Check EVERY effort for completion without review
for effort in $(yq '.efforts_in_progress[]' orchestrator-state.yaml); do
    IMPL_STATUS=$(yq ".efforts_in_progress[] | select(.name == \"$effort\") | .implementation_status" orchestrator-state.yaml)
    REVIEW_STATUS=$(yq ".efforts_in_progress[] | select(.name == \"$effort\") | .review_status" orchestrator-state.yaml)
    
    if [ "$IMPL_STATUS" = "COMPLETE" ] && [ "$REVIEW_STATUS" != "IN_PROGRESS" ] && [ "$REVIEW_STATUS" != "PASSED" ]; then
        echo "🚨 CRITICAL: Effort $effort implementation COMPLETE but review NOT STARTED!"
        echo "📝 SPAWNING CODE REVIEWER IMMEDIATELY for $effort"
        
        # Spawn the Code Reviewer NOW
        cd /efforts/phase${PHASE}/wave${WAVE}/${effort}
        Task: subagent_type="code-reviewer" \
              prompt="Review implementation in ${effort}. 
              CRITICAL: CD into effort directory first - it's a separate git repo!
              Get branch name with 'git branch --show-current' (NOT directory name!).
              Check: Size compliance (<800 lines using line-counter.sh -b phase${PHASE}/integration -c ACTUAL_BRANCH_NAME), Code quality, Tests pass.
              Create CODE-REVIEW-REPORT.md with status: PASSED/FAILED/NEEDS_SPLIT.
              If NEEDS_SPLIT, create SPLIT-PLAN.md." \
              description="Code Review ${effort}"
        
        # Update state file
        yq -i ".efforts_in_progress[] |= select(.name == \"$effort\") |= .review_status = \"IN_PROGRESS\"" orchestrator-state.yaml
        
        echo "✅ Code Reviewer spawned for $effort"
    fi
done
```

### Agent Progress Tracking
Monitor all agents continuously:
1. Check agent status every 10 minutes
2. **DETECT IMPLEMENTATION COMPLETION IMMEDIATELY**
3. **SPAWN CODE REVIEWERS FOR COMPLETED IMPLEMENTATIONS**
4. Validate progress against expectations
5. Detect size violations immediately (MUST use line-counter.sh with BRANCH names in effort repos!)
6. Identify blocked/stalled agents
7. Track completion status

### Dependency Coordination
When monitoring dependent efforts:
1. Track prerequisite completion status
2. Notify dependent agents when prerequisites ready
3. Prevent premature starts
4. Optimize start times for maximum parallelization
5. Handle dependency failures gracefully

### Intervention Triggers

**IMMEDIATE (Stop monitoring, take action):**
- Agent unresponsive >15 minutes
- Size limit exceeded
- Critical test failures
- Build system failure

**WARNING (Alert, continue monitoring):**
- Progress <70% of expected
- Timeline utilization >80%
- Agent reporting difficulties

**OPTIMIZATION (Suggest improvements):**
- Progress significantly ahead
- Resource underutilization
- Potential for increased parallelization

## State Transitions

From MONITOR state:
- **IMPLEMENTATION_COMPLETE + NO_REVIEW** → SPAWN_CODE_REVIEWER (Must spawn reviewer immediately!)
- **ALL_COMPLETE + ALL_REVIEWS_PASS** → WAVE_COMPLETE (All agents finished successfully)
- **REVIEW_FAILED** → FIX_ISSUES (SW Engineer must fix)
- **CRITICAL_ERROR** → ERROR_RECOVERY (Major failure requiring intervention)
- **SIZE_VIOLATIONS** → CREATE_SPLIT_PLAN (Code Reviewer creates split plan)
- **AGENTS_ACTIVE** → MONITOR (Continue monitoring)

**CRITICAL FLOW FOR REVIEWS:**
```
MONITOR detects implementation_status: COMPLETE
    ↓ (MUST spawn Code Reviewer)
SPAWN_CODE_REVIEWER for that effort
    ↓ (Review in progress)
MONITOR continues, tracking review_status
    ↓ (Review completes)
If PASSED → Check if all efforts complete
If FAILED → Spawn SW Engineer to FIX_ISSUES
If NEEDS_SPLIT → Spawn Code Reviewer for SPLIT_PLAN
    ↓ (After Code Reviewer creates split plans)
ORCHESTRATOR CREATES SPLIT INFRASTRUCTURE (R204)
    ↓ (Infrastructure ready)
Then spawn SW Engineer for sequential splits
```

## 🚨🚨🚨 SPLIT TRACKING REQUIREMENTS (R302) 🚨🚨🚨

**MANDATORY: Track all split operations meticulously:**

When monitoring splits:
1. **Update split_tracking section** in orchestrator-state.yaml
2. **Track each split branch** status (ACTIVE, COMPLETED, REVIEWED)
3. **Record line counts** for each completed split
4. **Mark original as SPLIT_DEPRECATED** when all splits done
5. **Update integration planning** to use split branches

## 🔴🔴🔴 ORCHESTRATOR CREATES SPLIT INFRASTRUCTURE (R204) 🔴🔴🔴

**CRITICAL: When Code Reviewer creates split plans, ORCHESTRATOR MUST:**

1. **WAIT for Code Reviewer to complete split planning**
   - SPLIT-INVENTORY.md created in too-large branch
   - All SPLIT-PLAN-XXX.md files created
   - Plans committed and pushed to remote

2. **ORCHESTRATOR CREATES ALL SPLIT INFRASTRUCTURE:**
```bash
# When split plans detected, orchestrator must:
create_split_infrastructure_after_plan() {
    local EFFORT_NAME="$1"
    local TOO_LARGE_DIR="/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
    
    echo "🔴🔴🔴 R204: ORCHESTRATOR CREATING SPLIT INFRASTRUCTURE"
    
    # Pull latest from too-large branch to get split plans
    cd "$TOO_LARGE_DIR"
    git pull
    
    # Check split plans exist
    if [ ! -f "SPLIT-INVENTORY.md" ]; then
        echo "❌ Waiting for Code Reviewer to create split plans..."
        return 1
    fi
    
    # Count splits needed
    TOTAL_SPLITS=$(grep -c "^| [0-9]" SPLIT-INVENTORY.md)
    echo "📋 Found $TOTAL_SPLITS splits to create"
    
    # Create infrastructure for EACH split
    for split_num in $(seq 1 $TOTAL_SPLITS); do
        SPLIT_NAME=$(printf "%03d" $split_num)
        SPLIT_DIR="/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}-SPLIT-${SPLIT_NAME}"
        
        echo "Creating split $SPLIT_NAME infrastructure..."
        
        # 1. Create split directory
        mkdir -p "$SPLIT_DIR"
        
        # 2. Clone target repo (sequential branching!)
        if [ $split_num -eq 1 ]; then
            # First split from same base as original
            git clone --branch "$BASE_BRANCH" "$TARGET_REPO_URL" "$SPLIT_DIR"
        else
            # Subsequent splits from previous split
            PREV_SPLIT=$(printf "%03d" $((split_num - 1)))
            PREV_BRANCH="${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}-SPLIT-${PREV_SPLIT}"
            git clone --branch "$PREV_BRANCH" "$TARGET_REPO_URL" "$SPLIT_DIR"
        fi
        
        # 3. Create split branch
        cd "$SPLIT_DIR"
        SPLIT_BRANCH="${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}-SPLIT-${SPLIT_NAME}"
        git checkout -b "$SPLIT_BRANCH"
        git push -u origin "$SPLIT_BRANCH"
        
        # 4. Copy split plan from too-large branch
        cp "$TOO_LARGE_DIR/SPLIT-PLAN-${SPLIT_NAME}.md" .
        
        # 5. Commit initial setup
        git add -A
        git commit -m "chore: initialize split-${SPLIT_NAME} infrastructure"
        git push
        
        echo "✅ Split $SPLIT_NAME infrastructure ready"
    done
    
    echo "✅ ALL SPLIT INFRASTRUCTURE CREATED - Ready for SW Engineer"
}
```

3. **THEN AND ONLY THEN spawn SW Engineer for splits**

```bash
# Check for split completion markers
if [ -f "/tmp/splits-complete-${EFFORT_NAME}.marker" ]; then
    echo "✅ All splits complete for $EFFORT_NAME"
    # Update split_tracking in state file
    # Mark original branch as SPLIT_DEPRECATED
    # List all replacement splits
fi
```

## Critical Requirements Summary

1. **Never stop monitoring** - R021 violation = -100%
2. **Save TODOs every 15 minutes** - R287 violation = -15% per occurrence
3. **Verify work locations** - R255 violation = -100%
4. **Report agent failures** - R254 violation = -50%
5. **Update state continuously** - R288 violation = -50%
6. **Track split operations** - R302 violation = -30%

## Success Criteria
- ✅ All agents tracked continuously
- ✅ Progress validated every 10 minutes
- ✅ TODOs saved every 15 minutes
- ✅ All work in correct locations
- ✅ All reviews passed before completion

## Failure Triggers
- ❌ Stop monitoring = R021 VIOLATION
- ❌ Forget TODO saves = R287 VIOLATION
- ❌ Accept wrong location work = R255 VIOLATION
- ❌ Miss agent failures = R254 VIOLATION
