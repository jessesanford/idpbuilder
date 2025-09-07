# Orchestrator - MONITOR_IMPLEMENTATION State Rules

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

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED MONITOR_IMPLEMENTATION STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_MONITOR_IMPLEMENTATION
echo "$(date +%s) - Rules read and acknowledged for MONITOR_IMPLEMENTATION" > .state_rules_read_orchestrator_MONITOR_IMPLEMENTATION
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## 📋 PRIMARY DIRECTIVES FOR MONITOR_IMPLEMENTATION STATE

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
**Summary**: Update orchestrator-state.yaml on all state changes

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

## 🔴🔴🔴 CRITICAL: MONITOR_IMPLEMENTATION IS A VERB - START MONITORING IMPLEMENTATIONS NOW! 🔴🔴🔴

**MONITOR_IMPLEMENTATION MEANS ACTIVELY MONITORING IMPLEMENTATIONS RIGHT NOW!**
- ❌ NOT "I'm in monitor implementation state"  
- ❌ NOT "Ready to monitor implementations"
- ✅ YES "I'm checking SW Engineer E3.1.2 progress NOW"
- ✅ YES "I'm verifying implementation in E3.1.3 NOW"
- ✅ YES "I'm detecting implementation completion NOW"

## State Context
MONITOR_IMPLEMENTATION = You ARE ACTIVELY monitoring spawned SW Engineers implementing features THIS INSTANT. Focus specifically on implementation progress, not reviews or fixes.

## 🚨🚨🚨 IMMEDIATE ACTIONS UPON ENTERING MONITOR_IMPLEMENTATION STATE 🚨🚨🚨

**THE INSTANT YOU ENTER MONITOR_IMPLEMENTATION STATE, DO THIS:**

```bash
# ✅ CORRECT - IMMEDIATE ACTION
echo "🔍 MONITORING IMPLEMENTATIONS: Checking all active SW Engineers NOW..."

# Step 1: List all SW Engineers being monitored (DO NOW!)
echo "📊 Active SW Engineers under monitoring:"
for effort in $(yq '.efforts_in_progress[].name' orchestrator-state.yaml); do
    IMPL_STATUS=$(yq ".efforts_in_progress[] | select(.name == \"$effort\") | .implementation_status" orchestrator-state.yaml)
    
    if [ "$IMPL_STATUS" = "IN_PROGRESS" ]; then
        echo "  - $effort: Checking implementation progress..."
        # Check agent working directory
        # Verify commits being made
        # Check for completion markers
    fi
done

# Step 2: Check for completed implementations (DO NOW!)
echo "🔍 Checking for completed implementations..."
for effort in $(yq '.efforts_in_progress[].name' orchestrator-state.yaml); do
    IMPL_STATUS=$(yq ".efforts_in_progress[] | select(.name == \"$effort\") | .implementation_status" orchestrator-state.yaml)
    
    if [ "$IMPL_STATUS" != "COMPLETE" ]; then
        # Check for completion markers
        if [ -f "/efforts/phase${PHASE}/wave${WAVE}/${effort}/IMPLEMENTATION-COMPLETE.marker" ]; then
            echo "✅ Implementation complete for $effort!"
            # Update state file
            yq -i ".efforts_in_progress[] |= select(.name == \"$effort\") |= .implementation_status = \"COMPLETE\"" orchestrator-state.yaml
            
            echo "🚨 CRITICAL: Implementation complete - must spawn Code Reviewer!"
            echo "➡️ Transitioning to SPAWN_CODE_REVIEWERS_FOR_REVIEW"
        fi
    fi
done

# Step 3: Check for blocked implementations (DO NOW!)
echo "🚧 Checking for blocked SW Engineers..."
for effort in $(yq '.efforts_in_progress[].name' orchestrator-state.yaml); do
    # Check for BLOCKED markers or stalled progress
    if [ -f "/efforts/phase${PHASE}/wave${WAVE}/${effort}/BLOCKED.marker" ]; then
        echo "⚠️ SW Engineer blocked on $effort!"
        # Determine intervention needed
    fi
done

# Step 4: Verify work locations (DO NOW!)
echo "📍 Verifying all work in correct locations (R255)..."
for effort in $(yq '.efforts_in_progress[].name' orchestrator-state.yaml); do
    EXPECTED_DIR="/efforts/phase${PHASE}/wave${WAVE}/${effort}"
    if [ -d "$EXPECTED_DIR" ]; then
        cd "$EXPECTED_DIR"
        CURRENT_BRANCH=$(git branch --show-current)
        echo "  - $effort: Branch $CURRENT_BRANCH in $EXPECTED_DIR ✅"
    else
        echo "  ❌ ERROR: Expected directory missing for $effort!"
    fi
done

# Step 5: Determine next action (DO NOW!)
echo "🎯 Determining immediate next action based on implementation monitoring..."
```

## Monitoring Implementation Focus Areas

### 1. Progress Tracking
Monitor SW Engineers actively implementing features:
- Check commit frequency (should see commits every 15-30 minutes)
- Verify work happening in correct directories
- Monitor for completion markers
- Track implementation velocity

### 2. Completion Detection
**CRITICAL**: When ANY implementation completes:
1. **IMMEDIATELY** update implementation_status to COMPLETE
2. **STOP** monitoring that implementation
3. **TRANSITION** to SPAWN_CODE_REVIEWERS_FOR_REVIEW
4. **DO NOT** wait for all implementations to complete

**🔴🔴🔴 SPLIT IMPLEMENTATION SPECIAL CASE 🔴🔴🔴**
When monitoring SPLIT implementations:
- Each split completion MUST trigger immediate review
- NEVER go directly to CREATE_NEXT_SPLIT_INFRASTRUCTURE
- ALWAYS: Split complete → Spawn reviewer → Wait for review → THEN decide next action
- Only create next split infrastructure if review passes AND more splits needed

### 3. Blocked Agent Detection
Identify and escalate blocked implementations:
- Agent hasn't committed in >30 minutes
- BLOCKED.marker file exists
- Error messages in agent output
- Dependency issues reported

### 4. Work Location Verification (R255)
Continuously verify:
- All work in correct effort directories
- Correct branches being used
- No work in wrong locations
- Remote pushes happening regularly

## State Transitions

From MONITOR_IMPLEMENTATION state:
- **IMPLEMENTATION_COMPLETE** → SPAWN_CODE_REVIEWERS_FOR_REVIEW (Any implementation done)
- **AGENT_BLOCKED** → ERROR_RECOVERY (SW Engineer stuck)
- **AGENT_FAILED** → SPAWN_AGENTS (Respawn with better instructions)
- **ALL_IMPLEMENTATIONS_COMPLETE** → MONITOR_REVIEWS (All done, reviews ongoing)
- **IMPLEMENTATIONS_ACTIVE** → MONITOR_IMPLEMENTATION (Continue monitoring)

## 🚨🚨🚨 CRITICAL ENFORCEMENT POINTS 🚨🚨🚨

### Implementation Complete = Immediate Review Spawn
**VIOLATION = -100% GRADE**: Allowing implementation to sit without review

When detecting implementation_status: COMPLETE:
1. **STOP** monitoring that effort
2. **SPAWN** Code Reviewer IMMEDIATELY
3. **UPDATE** review_status to IN_PROGRESS
4. **TRACK** spawned reviewer in state file

### Never Measure Code Yourself (R319)
**VIOLATION = -100% GRADE**: Using line-counter.sh as orchestrator

- ❌ NEVER run line-counter.sh yourself
- ❌ NEVER count lines manually
- ❌ NEVER check if splits are within size limits
- ✅ ALWAYS spawn Code Reviewer for ALL measurements
- ✅ Code Reviewer will determine size compliance
- ✅ Code Reviewer will verify split sizes

## Success Criteria
- ✅ All SW Engineers tracked continuously
- ✅ Implementation progress validated every 10 minutes
- ✅ Completion detected within 5 minutes of occurrence
- ✅ Code Reviewers spawned immediately on completion
- ✅ All work in correct locations
- ✅ TODOs saved every 15 minutes

## Failure Triggers
- ❌ Auto-transition without stop = R322 VIOLATION
- ❌ Measure code yourself = R319 VIOLATION
- ❌ Miss implementation completion = -50% penalty
- ❌ Delay review spawn = -100% penalty
- ❌ Accept wrong location work = R255 VIOLATION

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**