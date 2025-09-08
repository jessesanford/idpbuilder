# Orchestrator - MONITOR State Rules

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

### 🚨🚨🚨 R006 - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality**: BLOCKING - Automatic termination, 0% grade
**Summary**: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents
**CRITICAL**: Even when checking agent work - NO code operations allowed!

### ⚠️⚠️⚠️ R317 - Working Directory Restrictions (WARNING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R317-working-directory-restrictions.md`
**Criticality**: WARNING - -25% for violations
**Summary**: MUST NOT enter agent working directories - operate from root only
**CRITICAL**: Check agent work using absolute paths, never CD into their directories!

### 🚨🚨🚨 R318 - Agent Failure Escalation Protocol (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R318-agent-failure-escalation-protocol.md`
**Criticality**: BLOCKING - -40% for attempting forbidden fixes
**Summary**: NEVER fix agent failures directly - respawn with better instructions or escalate
**CRITICAL**: If agent fails, respawn or escalate - NEVER attempt DIY fixes!

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

### 🔴🔴🔴 R322 - Mandatory Stop Before State Transitions (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
**Criticality**: SUPREME LAW - Violation = -100% failure
**Summary**: MUST STOP and summarize before transitioning to next state

### 🚨🚨🚨 R287 - TODO Save Frequency Requirements
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-save-frequency.md`
**Criticality**: BLOCKING - Save every 15 minutes/10 messages
**Summary**: Mandatory TODO saves during monitoring

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state.json on all state changes

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
# 🔴🔴🔴 R319 CRITICAL: ORCHESTRATOR MUST NEVER MEASURE CODE! 🔴🔴🔴
# ORCHESTRATOR CANNOT USE line-counter.sh - That's Code Reviewer work!
# Instead, check if Code Reviewers have been spawned for completed efforts
for effort in $(yq '.efforts_in_progress[].name' orchestrator-state.json); do
    IMPL_STATUS=$(yq ".efforts_in_progress[] | select(.name == \"$effort\") | .implementation_status" orchestrator-state.json)
    REVIEW_STATUS=$(yq ".efforts_in_progress[] | select(.name == \"$effort\") | .review_status" orchestrator-state.json)
    
    if [ "$IMPL_STATUS" = "COMPLETE" ] && [ "$REVIEW_STATUS" != "IN_PROGRESS" ]; then
        echo "🚨 Effort $effort needs Code Reviewer for size assessment!"
        echo "🚀 Spawning Code Reviewer to measure and review..."
        # Spawn Code Reviewer who will use line-counter.sh
        # The Code Reviewer will determine if size violations exist
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
for effort in $(yq '.efforts_in_progress[]' orchestrator-state.json); do
    IMPL_STATUS=$(yq ".efforts_in_progress[] | select(.name == \"$effort\") | .implementation_status" orchestrator-state.json)
    REVIEW_STATUS=$(yq ".efforts_in_progress[] | select(.name == \"$effort\") | .review_status" orchestrator-state.json)
    
    if [ "$IMPL_STATUS" = "COMPLETE" ] && [ "$REVIEW_STATUS" != "IN_PROGRESS" ] && [ "$REVIEW_STATUS" != "PASSED" ]; then
        echo "🚨 CRITICAL: Effort $effort implementation COMPLETE but review NOT STARTED!"
        echo "📝 SPAWNING CODE REVIEWER IMMEDIATELY for $effort"
        
        # Spawn the Code Reviewer NOW
        cd /efforts/phase${PHASE}/wave${WAVE}/${effort}
        Task: subagent_type="code-reviewer" \
              prompt="Review implementation in ${effort}. 
              CRITICAL: CD into effort directory first - it's a separate git repo!
              Get branch name with 'git branch --show-current' (NOT directory name!).
              Check: Size compliance (<800 lines using line-counter.sh - auto-detects base), Code quality, Tests pass.
              Create review report in: .software-factory/phase${PHASE}/wave${WAVE}/${effort}/reports/CODE-REVIEW-REPORT-$(date +%Y%m%d-%H%M%S).md with status: PASSED/FAILED/NEEDS_SPLIT.
              If NEEDS_SPLIT, create split plan in: .software-factory/phase${PHASE}/wave${WAVE}/${effort}/plans/SPLIT-PLAN-$(date +%Y%m%d-%H%M%S).md
              Also create backward-compatible symlinks: ln -s to CODE-REVIEW-REPORT.md and SPLIT-PLAN.md in root." \
              description="Code Review ${effort}"
        
        # Update state file
        yq -i ".efforts_in_progress[] |= select(.name == \"$effort\") |= .review_status = \"IN_PROGRESS\"" orchestrator-state.json
        
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
5. **DELEGATE size assessment to Code Reviewers** (R006: Orchestrator NEVER measures!)
6. Identify blocked/stalled agents
7. Track completion status

**🔴 R006 CRITICAL: The orchestrator NEVER runs line-counter.sh! 🔴**
- Code Reviewers measure sizes
- Code Reviewers determine violations  
- Orchestrator only reads their reports

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
1. **Check if split infrastructure needed** (transition to CREATE_NEXT_SPLIT_INFRASTRUCTURE)
2. **Update split_tracking section** in orchestrator-state.json
3. **Track current split** status (ACTIVE, COMPLETED, REVIEWED)
4. **Detect when current split complete** and next split needed
5. **Record line counts** for each completed split
6. **Mark original as SPLIT_DEPRECATED** when all splits done
7. **Update integration planning** to use split branches

## 🔴🔴🔴 ORCHESTRATOR CREATES SPLIT INFRASTRUCTURE JUST-IN-TIME (R204) 🔴🔴🔴

**CRITICAL: When Code Reviewer creates split plans, ORCHESTRATOR MUST:**

1. **WAIT for Code Reviewer to complete split planning**
   - SPLIT-INVENTORY.md created in too-large branch
   - All SPLIT-PLAN-XXX.md files created
   - Plans committed and pushed to remote

2. **ORCHESTRATOR TRANSITIONS TO CREATE_NEXT_SPLIT_INFRASTRUCTURE:**
```bash
# When split plans detected, orchestrator must transition:
detect_split_plan_needs_infrastructure() {
    local EFFORT_NAME="$1"
    local TOO_LARGE_DIR="/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
    
    echo "🔴🔴🔴 R204: ORCHESTRATOR CREATING SPLIT INFRASTRUCTURE"
    
    # Load branch naming helpers for consistent naming
    source "$CLAUDE_PROJECT_DIR/utilities/branch-naming-helpers.sh"
    
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
    
    # 🔴 R308 CRITICAL: Determine base branch for splits!
    # Splits use SAME incremental base as original effort
    echo "🔴 R308: Determining incremental base for splits..."
    PHASE=$(yq '.current_phase' orchestrator-state.json)
    WAVE=$(yq '.current_wave' orchestrator-state.json)
    
    # Use the R308 function to get correct base
    if [[ $PHASE -eq 1 && $WAVE -eq 1 ]]; then
        SPLIT_BASE="main"
    elif [[ $WAVE -eq 1 ]]; then
        PREV_PHASE=$((PHASE - 1))
        SPLIT_BASE="phase${PREV_PHASE}-integration"
        echo "🔴 R308: Phase $PHASE splits from $SPLIT_BASE (NOT main!)"
    else
        PREV_WAVE=$((WAVE - 1))
        SPLIT_BASE="phase${PHASE}-wave${PREV_WAVE}-integration"
    fi
    
    echo "✅ R308: Splits will use incremental base: $SPLIT_BASE"
    
    # Get original branch name with prefix using helper function
    ORIGINAL_BRANCH=$(get_effort_branch_name "$PHASE" "$WAVE" "$EFFORT_NAME")
    echo "📌 Original branch (with prefix): $ORIGINAL_BRANCH"
    
    # Check if first split infrastructure exists
    SPLIT_001_DIR="/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}-SPLIT-001"
    
    if [ ! -d "$SPLIT_001_DIR" ]; then
        echo "🔴 First split infrastructure not created yet"
        echo "➡️ Transitioning to CREATE_NEXT_SPLIT_INFRASTRUCTURE"
        transition_to_state "CREATE_NEXT_SPLIT_INFRASTRUCTURE"
        return
    fi
    
    # Check if current split is complete and next needs creation
    CURRENT_SPLIT=$(yq ".split_tracking.\"$EFFORT_NAME\".current_split // 0" orchestrator-state.json)
    SPLIT_STATUS=$(yq ".split_tracking.\"$EFFORT_NAME\".splits[$CURRENT_SPLIT].status" orchestrator-state.json)
    
    if [ "$SPLIT_STATUS" = "COMPLETED" ]; then
        NEXT_SPLIT=$((CURRENT_SPLIT + 1))
        if [ $NEXT_SPLIT -le $TOTAL_SPLITS ]; then
            echo "🔴 Split $CURRENT_SPLIT complete, need infrastructure for split $NEXT_SPLIT"
            echo "➡️ Transitioning to CREATE_NEXT_SPLIT_INFRASTRUCTURE"
            transition_to_state "CREATE_NEXT_SPLIT_INFRASTRUCTURE"
            return
        fi
    fi
        SPLIT_NAME=$(printf "%03d" $split_num)
        # Directory uses -SPLIT- (uppercase, single hyphen) per R204
        SPLIT_DIR="/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}-SPLIT-${SPLIT_NAME}"
        # Branch uses --split- (lowercase, double hyphen) per helper function
        SPLIT_BRANCH=$(get_split_branch_name "$ORIGINAL_BRANCH" "$SPLIT_NAME")
        
        echo "Creating split $SPLIT_NAME infrastructure..."
        echo "  Directory: $SPLIT_DIR"
        echo "  Branch: $SPLIT_BRANCH"
        
        # 1. Create split directory
        mkdir -p "$SPLIT_DIR"
        
        # 2. Clone target repo (R308 incremental base!)
        if [ $split_num -eq 1 ]; then
            # First split from same base as original
            git clone --branch "$BASE_BRANCH" "$TARGET_REPO_URL" "$SPLIT_DIR"
        else
            # Subsequent splits from previous split
            PREV_SPLIT=$(printf "%03d" $((split_num - 1)))
            PREV_BRANCH=$(get_split_branch_name "$ORIGINAL_BRANCH" "$PREV_SPLIT")
            echo "  Basing on previous: $PREV_BRANCH"
            git clone --branch "$PREV_BRANCH" "$TARGET_REPO_URL" "$SPLIT_DIR"
        fi
        
        # 3. Create split branch
        cd "$SPLIT_DIR"
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

1. **Monitor until state transition ready** - Then follow R322 to stop
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
- ❌ Auto-transition without stop = R322 VIOLATION
- ❌ Forget TODO saves = R287 VIOLATION
- ❌ Accept wrong location work = R255 VIOLATION
- ❌ Miss agent failures = R254 VIOLATION

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
