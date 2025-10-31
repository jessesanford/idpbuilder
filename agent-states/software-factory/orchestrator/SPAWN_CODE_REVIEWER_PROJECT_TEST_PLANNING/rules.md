# ORCHESTRATOR STATE: SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING

## 🔴🔴🔴 STATE MANAGER BOOKEND REMINDER 🔴🔴🔴

**CRITICAL:** This state's completion checklist shows jq/yq commands to update state files.
**THESE ARE PROHIBITED IN SF 3.0!**

**MANDATORY PATTERN:**
```bash
# WRONG (SF 2.0 - DO NOT USE):

# CORRECT (SF 3.0 - ALWAYS USE):
# Let State Manager handle it via variables:
NEXT_STATE="WAITING_FOR_PROJECT_TEST_PLAN"
TRANSITION_REASON="Spawned code reviewer for project tests"
# Then use /continue-software-factory which calls State Manager
```

**YOU MUST:**
1. Set NEXT_STATE and TRANSITION_REASON variables
2. Complete state work
3. Set CONTINUE-SOFTWARE-FACTORY=TRUE
4. Exit and let State Manager update state file

**See:** $CLAUDE_PROJECT_DIR/rule-library/R600-state-manager-bookend-protocol.md

---

# PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R355** - Code Reviewer Test Planning
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R355-production-ready-code-enforcement-supreme-law.md`

4. **R287** - TODO Persistence Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`

5. **R288** - State File Update Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`

6. **R304** - Mandatory Line Counter Usage (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`

7. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`

8. **R324** - State Transition Validation (SUPREME LAW)
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


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

## Overview
This state spawns a Code Reviewer to create project-level tests from the master architecture, implementing TDD at the project level per R341.

## Entry Criteria
- Master architecture exists (PROJECT-ARCHITECTURE.md or MASTER-ARCHITECTURE.md)
- Transitioned from WAITING_FOR_MASTER_ARCHITECTURE
- No existing PROJECT-TEST-PLAN.md

## State Responsibilities

### 1. Verify Architecture Exists at R550 Standardized Location
```bash
# R550: Get standardized location from planning_files (NOT legacy planning_artifacts)
ARCH_PLAN_FILE=$(jq -r '.planning_files.project.architecture_plan // "planning/project/PROJECT-ARCHITECTURE-PLAN.md"' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")
ARCH_PLAN_PATH="$CLAUDE_PROJECT_DIR/$ARCH_PLAN_FILE"

# Ensure we have architecture to work from at R550 standardized location
if [ ! -f "$ARCH_PLAN_PATH" ]; then
    echo "❌ Cannot create tests without project architecture plan!"
    echo "   Expected R550 location: $ARCH_PLAN_PATH"
    echo "   (from orchestrator-state-v3.json planning_files.project.architecture_plan)"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Project architecture plan not found at R550 standardized location"
    exit 341
fi

echo "✅ Project architecture plan found at R550 location: $ARCH_PLAN_PATH"
```

### 2. Spawn Code Reviewer for Test Planning
```bash
/spawn-agent code-reviewer \
    --state PROJECT_TEST_PLANNING \
    --task "Create project-level tests from master architecture" \
    --tdd "Tests must be created BEFORE Phase 1 implementation per R341"
```

### 3. Record Spawn with R151 Timing
```bash
# Track spawn time for R151 compliance
SPAWN_TIME=$(date +%s)
yq -i '.spawned_agents.project_test_reviewer.id = "'$AGENT_ID'"' orchestrator-state-v3.json
yq -i '.spawned_agents.project_test_reviewer.timestamp = "'$SPAWN_TIME'"' orchestrator-state-v3.json
yq -i '.spawned_agents.project_test_reviewer.state = "PROJECT_TEST_PLANNING"' orchestrator-state-v3.json
```

### 4. Update State and Stop (R313)
```bash
NEXT_STATE="WAITING_FOR_PROJECT_TEST_PLAN"
TRANSITION_REASON="Spawned Code Reviewer for project test planning"
save_todos "Spawned Code Reviewer for project test planning"
git add todos/
git commit -m "todo: spawned code reviewer for project tests (R341 TDD)"
git push
echo "🛑 Stopping per R313 after spawn"
```

## Exit Criteria
- Code Reviewer spawned for PROJECT_TEST_PLANNING
- State file updated to WAITING_FOR_PROJECT_TEST_PLAN
- R313: MUST stop after spawning

## Success Metrics
- ✅ Code Reviewer spawned within 5s (R151)
- ✅ State updated before stop (R324)
- ✅ R341 TDD compliance tracked
- ✅ TODOs saved (R287)



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING complete - [describe what was accomplished]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Commit Work Products (if any)
```bash
# Commit any deliverables BEFORE state transition
# (Plans, reports, configurations, etc.)
git add [work-products]
git commit -m "feat/doc: [description of work products]"
git push
echo "✅ Work products committed"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING complete [R287]"; then
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

### ✅ Step 5: Verify State Transition Variables Set
```bash
# Verify required variables are set for State Manager
if [ -z "$NEXT_STATE" ] || [ -z "$TRANSITION_REASON" ]; then
    echo "❌ FATAL: State transition variables not set!"
    echo "NEXT_STATE='$NEXT_STATE'"
    echo "TRANSITION_REASON='$TRANSITION_REASON'"
    exit 1
fi
echo "✅ State transition variables verified"
echo "   NEXT_STATE=$NEXT_STATE"
echo "   TRANSITION_REASON=$TRANSITION_REASON"
```

---

### ✅ Step 6: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 7: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-software-factory to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**SF 3.0 uses 7-step exit (not 8-step) - State Manager handles state file updates!**

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever (variables not set)
- Missing Step 3: Work products not committed = lost deliverables
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: Variables not verified = silent failures
- Missing Step 6: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 7: No exit = R322 violation (-100%)

- Missing Step 2: No proposed next state = State Manager can't proceed
- Missing Step 3: No State Manager consultation = bypassing bookend pattern (-100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern) - NO EXCEPTIONS**

**ELIMINATED STEPS FROM SF 2.0:**
- ❌ Step 3 (old): Manual jq/yq state updates → State Manager handles this
- ❌ Step 4 (old): Manual state validation → State Manager validates
- ❌ Step 5 (old): Manual state commit → State Manager commits

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


## Related Rules
- R341: TDD - tests before implementation
- R342: Early branch creation for test storage
- R151: Parallel agent timing
- R313: Mandatory stop after spawn

