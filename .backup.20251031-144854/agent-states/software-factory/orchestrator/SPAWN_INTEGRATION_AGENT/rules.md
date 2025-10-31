# Orchestrator - SPAWN_INTEGRATION_AGENT State Rules

# PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R361** - Integration Conflict Resolution Only (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R361-integration-conflict-resolution-only.md`
   - Integration Agent can ONLY resolve conflicts, NO new code/packages

3. **R362** - No Architectural Rewrites Without Approval (SUPREME LAW #7)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R362-no-architectural-rewrites.md`
   - Integration Agent MUST NOT change architecture, remove libraries, or deviate from plan

4. **R269** - WAVE Integration Merge Plan Protocol (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R269-code-reviewer-merge-plan-no-execution.md`

5. **R270** - PHASE Integration Merge Plan Protocol (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R270-no-integration-branches-as-sources.md`

6. **R287** - TODO Persistence Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`

7. **R288** - State File Update Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`

8. **R304** - Mandatory Line Counter Usage (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`

9. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`

10. **R324** - State Transition Validation (SUPREME LAW)
    - File: `$CLAUDE_PROJECT_DIR/rule-library/R324-mandatory-line-counter-auto-detection.md`


# 🔴🔴🔴 MANDATORY: R322 STOP + R405 CONTINUATION FLAG 🔴🔴🔴

**CRITICAL FOR SPAWN STATES - READ THIS FIRST OR FAIL TEST 2!**

## 🚨 THE PATTERN THAT FAILED TEST 2 🚨

**WHAT HAPPENED IN TEST 2:**
- Orchestrator spawned Code Reviewers ✅ (correct)
- Orchestrator stopped per R322 ✅ (correct)
- Orchestrator **DID NOT emit `CONTINUE-SOFTWARE-FACTORY=TRUE`** ❌ (WRONG!)
- Test framework saw no continuation flag → stopped automation
- Test 2 FAILED at iteration 8

**ROOT CAUSE:** Confusion between R322 "stop" and R405 continuation flag

## 🔴 CRITICAL DISTINCTION: TWO INDEPENDENT DECISIONS 🔴

### Decision 1: Should Agent Stop? (R322 - Context Preservation)
**YES - ALWAYS stop after spawning for context preservation**

- **Purpose**: Prevent context overflow between states
- **Action**: `exit 0` to end conversation
- **User Experience**: User sees "/continue-orchestrating" as next step
- **This is NORMAL!** Not an error!

### Decision 2: Should Factory Continue? (R405 - Automation Control)
**YES - ALWAYS emit TRUE for normal spawning operations**

- **Purpose**: Tell automation whether it CAN restart
- **Action**: `echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"` (LAST output before exit)
- **Automation**: Framework will auto-restart orchestrator
- **This is NORMAL!** Designed behavior!

## ✅ REQUIRED PATTERN FOR ALL SPAWN STATES

```bash
# 1. Complete spawning work
echo "✅ Spawned [agent type] for [purpose]"

# 2. Update state file per R324/R288
update_state "[NEXT_STATE]"
commit_state_files_per_r288()

# 3. Save TODOs per R287
save_todos "SPAWNED_[AGENT_TYPE]"

# 4. R322: Stop conversation (context preservation)
echo "🛑 R322: Stopping after spawn for context preservation"

# 5. R405: CONTINUATION FLAG - MUST BE TRUE FOR SPAWNING!
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# 6. Exit to end conversation
exit 0
```

## ❌ WRONG PATTERN (CAUSES TEST FAILURES)

```bash
# ❌ THIS KILLS AUTOMATION - DO NOT DO THIS!
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
exit 0

# Result: Test framework stops, Test 2 fails at iteration 8
```

## 🎯 WHY TRUE IS CORRECT FOR SPAWNING

**Spawning is NORMAL operation:**
- ✅ System knows next state (from state machine)
- ✅ Automation can continue (designed workflow)
- ✅ No manual intervention needed
- ✅ Context preservation ≠ error condition

**The orchestrator stopping (`exit 0`) is for:**
- Preserving context between conversation turns
- Allowing state file commits
- Creating clean state boundaries

**The TRUE flag indicates:**
- Automation CAN restart the conversation
- System knows what to do next (check state file)
- Normal operation is proceeding

## 🔴 WHEN TO USE FALSE (NOT FOR SPAWNING!)

**FALSE should ONLY be used for catastrophic failures:**
- ❌ State file corrupted beyond parsing
- ❌ Critical infrastructure destroyed
- ❌ Unrecoverable system errors
- ❌ **NEVER for normal spawning operations!**

## 📋 SPAWN STATE CHECKLIST

**Before exiting this spawn state, verify:**
1. [ ] All agents spawned successfully
2. [ ] State file updated to next state per R324
3. [ ] State files committed per R288
4. [ ] TODOs saved per R287
5. [ ] R322 stop message displayed
6. [ ] **CONTINUE-SOFTWARE-FACTORY=TRUE emitted** ← Critical!
7. [ ] Exited with `exit 0`

**Missing step 6 = Test 2 failure = -100% grade**

---


## 🔴🔴🔴 R322 MANDATORY: STOP BEFORE STATE TRANSITIONS 🔴🔴🔴

**CRITICAL REQUIREMENT PER R322:**
After spawning ANY agents in this state, you MUST:
1. Record what was spawned in state file
2. Save TODOs per R287
3. Commit and push state changes
4. Display stop message with continuation instructions
5. EXIT immediately with code 0

**VIOLATION = AUTOMATIC -100% FAILURE**

See: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`

---


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_INTEGRATION_AGENT STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_INTEGRATION_AGENT
echo "$(date +%s) - Rules read and acknowledged for SPAWN_INTEGRATION_AGENT" > .state_rules_read_orchestrator_SPAWN_INTEGRATION_AGENT
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWN_INTEGRATION_AGENT WORK UNTIL RULES ARE READ:
- ❌ Start spawn integration specialist
- ❌ Start coordinate merging
- ❌ Start manage integration
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R151, R208, R053..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all SPAWN_INTEGRATION_AGENT rules"
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

### ✅ CORRECT PATTERN FOR SPAWN_INTEGRATION_AGENT:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SPAWN_INTEGRATION_AGENT work until:**
### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SPAWN_INTEGRATION_AGENT work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING_SWE_PROGRESS YOUR READ TOOL CALLS!**

## 🔴🔴🔴 R208 SUPREME LAW: CD BEFORE SPAWN 🔴🔴🔴

**YOU MUST CD TO INTEGRATE_WAVE_EFFORTS DIRECTORY BEFORE SPAWNING THE INTEGRATE_WAVE_EFFORTS AGENT!**
- Violation = -100% GRADE = AUTOMATIC FAILURE
- NO EXCEPTIONS, NO SHORTCUTS, NO WORKAROUNDS

## State Definition
The orchestrator spawns an Integration Agent to execute the merge plan created by Code Reviewer. The merge plan must already exist.

## Required Actions

### 1. Verify Merge Plan Exists
```bash
# Check that Code Reviewer completed the merge plan
PHASE=$(jq '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
WAVE=$(jq '.project_progression.current_wave.wave_number' orchestrator-state-v3.json)
INTEGRATE_WAVE_EFFORTS_DIR="/efforts/phase${PHASE}/wave${WAVE}/integration-workspace"

cd "$INTEGRATE_WAVE_EFFORTS_DIR"

if [ ! -f "WAVE-MERGE-PLAN.md" ]; then
    echo "❌ Cannot spawn Integration Agent - no merge plan!"
    echo "🔍 Code Reviewer must complete WAVE-MERGE-PLAN.md first"
    exit 1
fi

echo "✅ Found WAVE-MERGE-PLAN.md"

# 🔴🔴🔴 R312 EXCEPTION: UNLOCK CONFIG FOR INTEGRATE_WAVE_EFFORTS 🔴🔴🔴
echo ""
echo "🔓 R312 EXCEPTION: Unlocking git config for INTEGRATE_WAVE_EFFORTS agent"
echo "Integration agents NEED to pull from multiple branches"

# Check if config is locked
if [ ! -w .git/config ]; then
    echo "📋 Config is currently locked (as expected for efforts)"
    
    # Store current permissions and ownership for audit
    BEFORE_PERMS=$(stat -c %a .git/config 2>/dev/null || stat -f %A .git/config)
    BEFORE_OWNER=$(stat -c %U:%G .git/config 2>/dev/null || stat -f %Su:%Sg .git/config)
    
    # Check current ownership
    CURRENT_OWNER=$(stat -c %U:%G .git/config 2>/dev/null || stat -f %Su:%Sg .git/config)
    
    # Unlock for integration work - handle root ownership
    if [ "$CURRENT_OWNER" = "root:root" ]; then
        # Need sudo to change from root ownership
        if command -v sudo >/dev/null 2>&1; then
            echo "🔓 Restoring user ownership from root..."
            sudo chown $(id -u):$(id -g) .git/config
            sudo chmod 644 .git/config
        else
            echo "❌ ERROR: Config is root-owned but sudo not available!"
            echo "Cannot unlock config for integration"
            exit 312
        fi
    else
        # Simple permission change
        chmod 644 .git/config
    fi
    
    # Verify unlock succeeded
    if [ ! -w .git/config ]; then
        echo "❌ ERROR: Failed to unlock config for integration!"
        echo "Integration agent needs writable config to merge branches"
        exit 312
    fi
    
    # Create exception marker
    cat > .git/R312_INTEGRATE_WAVE_EFFORTS_EXCEPTION << EOF
Timestamp: $(date '+%Y-%m-%d %H:%M:%S')
Unlocked by: orchestrator
State: SPAWN_INTEGRATION_AGENT
Phase: ${PHASE}
Wave: ${WAVE}
Previous ownership: $BEFORE_OWNER
Current ownership: $(stat -c %U:%G .git/config 2>/dev/null || echo 'unknown')
Previous permissions: $BEFORE_PERMS
Current permissions: 644 (writable)
Purpose: Integration requires ability to merge from multiple branches
EOF
    
    echo "✅ R312 Exception Applied: Config unlocked for integration"
    echo "   Owner: $BEFORE_OWNER → $(stat -c %U:%G .git/config 2>/dev/null || echo 'unknown')"
    echo "   Permissions: $BEFORE_PERMS → 644"
    echo "📝 Integration agent can now:"
    echo "   ✅ Pull from multiple effort branches"
    echo "   ✅ Create integration branches"
    echo "   ✅ Merge efforts together"
else
    echo "✅ Config already writable (integration workspace)"
fi

# Quick validation of merge plan
MERGE_PLAN=$(ls -t .software-factory/phase${PHASE}/wave${WAVE}/integration/WAVE-MERGE-PLAN--*.md 2>/dev/null | head -1)
MERGE_COUNT=$([ -n "$MERGE_PLAN" ] && grep -c "git merge origin/" "$MERGE_PLAN" 2>/dev/null || echo "0")
echo "📊 Merge plan contains $MERGE_COUNT merge operations"

if [[ $MERGE_COUNT -eq 0 ]]; then
    echo "⚠️ Warning: No merge commands found in plan!"
fi
```

### 2. Spawn Integration Agent
```bash
# Prepare spawn command for Integration Agent
CURRENT_BRANCH=$(git branch --show-current)

cat > /tmp/integration-agent-task.md << EOF
Execute integration merges for Phase ${PHASE} Wave ${WAVE}.

🔴🔴🔴 R361 SUPREME LAW: CONFLICT RESOLUTION ONLY 🔴🔴🔴
- NO new packages or directories
- NO adapter or wrapper code
- NO "glue code" or compatibility layers
- Maximum 50 lines of changes total (excluding merges)
- Integration = conflict resolution ONLY

CRITICAL REQUIREMENTS (R260):
1. You are in INTEGRATE_WAVE_EFFORTS_DIR: ${INTEGRATE_WAVE_EFFORTS_DIR}
2. IMMEDIATELY acknowledge and set INTEGRATE_WAVE_EFFORTS_DIR variable
3. Verify you're in the correct directory
4. Read and follow WAVE-MERGE-PLAN.md EXACTLY
5. Execute merges in specified order
6. Handle conflicts as directed in plan (R361: choose versions, don't create new code)
7. Run tests after each merge
8. Document everything in work-log.md

🎬 DEMO REQUIREMENTS (R291/R330):
9. Execute effort demos after each merge (see merge plan)
10. Run integrated wave demo after all merges
11. Capture all demo outputs in demo-results/
12. Document demo status in INTEGRATE_WAVE_EFFORTS_REPORT.md
13. If ANY demo fails, mark Demo Status: FAILED (triggers ERROR_RECOVERY)

R291 GATES YOU MUST ENFORCE:
- BUILD GATE: Code must compile
- TEST GATE: All tests must pass
- DEMO GATE: Demo scripts must execute
- ARTIFACT GATE: Build outputs must exist

GRADING CRITERIA (R267):
- 50% Completeness of Integration (including demos)
- 50% Meticulous Tracking and Documentation

Your working directory: ${INTEGRATE_WAVE_EFFORTS_DIR}
Current branch: ${CURRENT_BRANCH}
Merge plan to follow: WAVE-MERGE-PLAN.md

You are spawned into state: INIT
EOF

# 🔴🔴🔴 R208 SUPREME LAW: CD BEFORE SPAWN 🔴🔴🔴
echo "🔴 R208 SUPREME LAW: CD'ing to integration directory"
cd "$INTEGRATE_WAVE_EFFORTS_DIR" || {
    echo "❌ R208 VIOLATION: Failed to CD to $INTEGRATE_WAVE_EFFORTS_DIR"
    echo "❌ GRADE: -100% (AUTOMATIC FAILURE)"
    exit 1
}

echo "📍 R208 PWD VERIFICATION: $(pwd)"  # MUST show integration directory
echo "✅ R208: Confirmed in correct directory for Integration Agent spawn"

echo "🚀 Spawning Integration Agent for merge execution..."

/spawn integration-agent INIT "$(cat /tmp/integration-agent-task.md)"

# R208: Return to orchestrator directory after spawn
cd /workspaces/project
echo "📍 R208: Returned to orchestrator directory"
```

### 3. Update State Tracking
```yaml
# Update orchestrator-state-v3.json
integration_status:
  phase: ${PHASE}
  wave: ${WAVE}
  merge_plan_ready: true
  integration_agent_spawned: true
  integration_started_at: "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  waiting_for: "integration-agent-completion"
```

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## Critical Requirements

### R260 Compliance - INTEGRATE_WAVE_EFFORTS_DIR Acknowledgment
The Integration Agent MUST:
1. Acknowledge the INTEGRATE_WAVE_EFFORTS_DIR immediately upon startup
2. Set INTEGRATE_WAVE_EFFORTS_DIR environment variable
3. Verify current directory matches INTEGRATE_WAVE_EFFORTS_DIR
4. Exit with error if in wrong directory

### Working Directory Setup
**CRITICAL**: The orchestrator MUST cd into INTEGRATE_WAVE_EFFORTS_DIR before spawning!
```bash
cd "$INTEGRATE_WAVE_EFFORTS_DIR"  # MANDATORY before spawn
```

## Transition Rules

### 🔴🔴🔴 CRITICAL: Update State BEFORE Stopping! 🔴🔴🔴
Per R322, you MUST update `current_state` to the next state BEFORE stopping:

```bash
# After spawning integration agent successfully:
echo "📝 Updating state file for transition..."
jq '.transition_time = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
git add orchestrator-state-v3.json
git commit -m "state: transition from SPAWN_INTEGRATION_AGENT to MONITORING_INTEGRATE_WAVE_EFFORTS"
git push

# THEN stop per R322
echo "🛑 Stopping before MONITORING_INTEGRATE_WAVE_EFFORTS state (per R322)"
```

- Next state: MONITORING_INTEGRATE_WAVE_EFFORTS (UPDATE STATE FIRST!)
- Cannot transition if: No merge plan exists
- Must be in integration directory when spawning

## Success Criteria
- Merge plan verified to exist
- Integration Agent spawned with INTEGRATE_WAVE_EFFORTS_DIR
- Working directory set correctly
- Clear grading criteria communicated
- State tracking updated



## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**


### 🔴🔴🔴 MANDATORY VALIDATION REQUIREMENT 🔴🔴🔴

**Per R288 and R324**: ALL state file updates MUST be validated before commit:

```bash
# After ANY update to orchestrator-state-v3.json:
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
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



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete SPAWN_INTEGRATION_AGENT:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="SPAWN_INTEGRATION_AGENT complete - [describe what was accomplished]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for SHUTDOWN_CONSULTATION
```bash
# State Manager validates transition and updates state files (SF 3.0 Pattern)
echo "🔄 Spawning State Manager for SHUTDOWN_CONSULTATION..."

# Prepare work results summary
WORK_RESULTS=$(cat <<EOF
{
  "state_completed": "SPAWN_INTEGRATION_AGENT",
  "work_accomplished": [
    "[List accomplishments from state work]"
  ],
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "SPAWN_INTEGRATION_AGENT" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON" \
  --work-results "$WORK_RESULTS"

# State Manager will:
# 1. Validate PROPOSED_NEXT_STATE exists and transition is valid
# 2. Update all 4 state files atomically (R288)
# 3. Commit and push state files
# 4. Return REQUIRED_NEXT_STATE (usually same as proposed unless invalid)

echo "✅ State Manager consultation complete"
echo "✅ State files updated by State Manager"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "SPAWN_INTEGRATION_AGENT_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - SPAWN_INTEGRATION_AGENT complete [R287]"; then
    echo "❌ ERROR: Failed to commit TODO files"
    echo "This is non-fatal but TODOs may be lost in compaction"
    echo "Proceeding with state execution..."
    # Don't exit - TODO commit failure is not fatal
fi

git push || echo "⚠️ WARNING: TODO push failed - committed locally"
echo "✅ TODOs saved and committed"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 5: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 6: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No state update = state machine broken (R288 violation, -100%)
- Missing Step 4: Invalid state = corruption (R324 violation)
- Missing Step 5: No commit = state lost on compaction (R288 violation, -100%)
- Missing Step 6: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 7: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 8: No exit = R322 violation (-100%)

- Missing Step 2: No proposed next state = State Manager can't proceed
- Missing Step 3: No State Manager consultation = bypassing bookend pattern (-100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern) - NO EXCEPTIONS**

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS (SF 3.0)

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="NEXT_STATE"
TRANSITION_REASON="State work complete"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"
# State Manager updates all 4 state files atomically

# 4. Save TODOs
save_todos "R322_CHECKPOINT"

# 5. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# 6. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
### 🚨 SPAWN STATE PATTERN - R322 + R405 USAGE 🚨

**Spawning operations require R322 stop for context preservation:**
```bash
# After spawning agent(s)
echo "✅ Spawned agents for work"

# R322 checkpoint (context preservation)
echo "🛑 R322: Stopping after spawn for context preservation"

# Flag? → MUST BE TRUE (normal operation!)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# Stop inference
exit 0
```

**Why TRUE is correct:**
- Spawning is NORMAL operation
- System knows next state
- Automation can continue
- **Context preservation ≠ manual intervention needed!**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

