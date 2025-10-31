# ORCHESTRATOR STATE: WAITING_FOR_PROJECT_TEST_PLAN


## 🚨 State Manager Bookend Pattern (MANDATORY)

**BEFORE this state**:
- State Manager validated transition via STARTUP_CONSULTATION
- You are here because State Manager directed you here
- orchestrator-state-v3.json shows validated_by: "state-manager"

**DURING this state**:
- Perform state-specific work
- NEVER call update_state directly
- Prepare results for State Manager
- Propose next state (don't decide!)

**AFTER this state**:
- Spawn State Manager SHUTDOWN_CONSULTATION
- Provide results and proposed next state
- State Manager validates and decides actual next state
- Transition to State Manager's required_next_state

**CRITICAL**: The orchestrator PROPOSES, the State Manager DECIDES!

---

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_PROJECT_TEST_PLAN STATE

### Core Mandatory Rules (ALL orchestrator states must have these):

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Criticality: BLOCKING - Automatic termination, 0% grade
   - Summary: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents

2. **🔴🔴🔴 R287** - TODO PERSISTENCE COMPREHENSIVE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Criticality: SUPREME - -20% to -100% penalty for violations
   - Summary: MUST save TODOs within 30s after write, every 10 messages, before transitions

3. **🔴🔴🔴 R288** - STATE FILE UPDATE REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state-v3.json before EVERY state transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

5. **🔴🔴🔴 R233** - Immediate Action On State Entry
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-all-states-immediate-action.md`
   - Criticality: SUPREME LAW - Must act immediately on entering state
   - Summary: WAITING states require active monitoring, not passive waiting

## Overview
This state waits for the Code Reviewer to complete project-level test planning, then enforces R342 by transitioning to early branch creation.

## Entry Criteria
- Transitioned from SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING
- Code Reviewer spawned for project test planning

## State Responsibilities

### 1. Check for Test Plan Completion
```bash
# Look for completed project test plan
if [ -f "PROJECT-TEST-PLAN.md" ]; then
    echo "✅ Project test plan completed"
    TEST_READY=true
else
    echo "⏳ Waiting for project test plan..."
    TEST_READY=false
fi
```

### 2. Verify Test Artifacts
When tests are ready, verify:
```bash
# Check for required test artifacts
[ -f "PROJECT-TEST-PLAN.md" ] || { echo "❌ Missing test plan"; exit 1; }
[ -f "PROJECT-TEST-HARNESS.sh" ] || { echo "⚠️ Missing test harness"; }
[ -d "project-tests/" ] || { echo "⚠️ Missing test directory"; }
```

### 3. Update Planning File Tracking (R340)
```bash
# Parse test plan metadata from code reviewer output
TEST_PLAN_PATH=$(grep "📋 Test Plan:" code_reviewer_output.txt | awk '{print $NF}')
TEST_HARNESS_PATH=$(grep "📋 Test Harness:" code_reviewer_output.txt | awk '{print $NF}')
DEMO_SCENARIOS_PATH=$(grep "📋 Demo Scenarios:" code_reviewer_output.txt | awk '{print $NF}')

# Record test plan location and metadata in R340 format
jq --arg test_plan "$TEST_PLAN_PATH" \
   --arg test_harness "$TEST_HARNESS_PATH" \
   --arg demo_scenarios "$DEMO_SCENARIOS_PATH" \
   --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
   --arg creator "code-reviewer" \
   '
   .planning_repo_files.test_plans.project = {
     "test_plan_path": $test_plan,
     "test_harness_path": $test_harness,
     "demo_scenarios_path": $demo_scenarios,
     "created_at": $timestamp,
     "created_by": $creator,
     "status": "active"
   }
   ' orchestrator-state-v3.json > orchestrator-state-v3.json.tmp
mv orchestrator-state-v3.json.tmp orchestrator-state-v3.json

# Commit R340 update
git add orchestrator-state-v3.json
git commit -m "state: record project test plan location per R340"
git push
```

### 4. Prepare for R342 Early Branch Creation
```bash
if [ "$TEST_READY" = true ]; then
    echo "🔴🔴🔴 R342 ENFORCEMENT: Test plan complete, preparing for branch creation"
    PROPOSED_NEXT_STATE="CREATE_PROJECT_INTEGRATION_BRANCH_EARLY"
    TRANSITION_REASON="R342 requires early integration branch creation for test storage"
    echo "📍 Will propose: $PROPOSED_NEXT_STATE to State Manager"
    # NOTE: Actual transition via State Manager in completion checklist
fi
```

## Exit Criteria
- Project test plan exists
- Test artifacts verified
- R340 metadata tracking complete
- State transitions to CREATE_PROJECT_INTEGRATION_BRANCH_EARLY (R342)

## Success Metrics
- ✅ PROJECT-TEST-PLAN.md created
- ✅ Test harness present
- ✅ R340 planning file tracking complete
- ✅ R342 early branch creation enforced



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete WAITING_FOR_PROJECT_TEST_PLAN:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, set variables for State Manager
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="[REASON_FOR_TRANSITION]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for State Transition (R288 - SUPREME LAW)
```bash
# State Manager handles ALL state file updates atomically
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE_NAME" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"

# State Manager will:
# 1. Validate the transition against state machine
# 2. Update all 4 state tracking locations atomically:
#    - orchestrator-state-v3.json
#    - orchestrator-state-demo.json
#    - .cascade-state-backup.json
#    - .orchestrator-state-v3.json
# 3. Commit and push all changes
# 4. Return control

echo "✅ State Manager completed transition"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before exit (R287 trigger)
save_todos "STATE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - state complete [R287]"; then
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
- Missing Step 3: No State Manager spawn = state machine broken (R288 violation, -100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)

**CRITICAL**: Steps 2-5 enforce State Manager bookend pattern - orchestrator PROPOSES, State Manager DECIDES!

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
PROPOSED_NEXT_STATE="CREATE_PROJECT_INTEGRATION_BRANCH_EARLY"
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
### 🚨 WAITING STATE PATTERN - CRITICAL UNDERSTANDING 🚨

**This is a WAITING state. Common source of incorrect FALSE usage!**

**WRONG interpretation:**
> "R322 mandates stop before transition"
> "State work is complete (validation done)"
> "User needs to invoke /continue-orchestrating"
> "Therefore I must set CONTINUE-SOFTWARE-FACTORY=FALSE"

**CORRECT interpretation:**
> "R322 checkpoint is NORMAL procedure for context preservation"
> "State work completed successfully = NORMAL outcome"
> "Waiting for /continue is DESIGNED user experience"
> "System KNOWS next step from state file"
> "NO manual intervention required, just normal continuation"
> "Therefore set CONTINUE-SOFTWARE-FACTORY=TRUE"

**The key distinction:**
- **Stopping inference** (`exit 0`) = Context management (ALWAYS at R322 points)
- **Continuation flag** = Can automation proceed? (TRUE unless catastrophic failure)

**ONLY use FALSE if:**
- ❌ The thing we're waiting for completely disappeared (agents crashed with no recovery)
- ❌ Results arrived but are completely corrupted/unreadable
- ❌ State file corruption prevents determining what to wait for
- ❌ System deadlock with no automated resolution
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**


## Related Rules
- R340: Planning file metadata tracking
- R341: TDD - tests created before implementation
- R342: Early integration branch creation (SUPREME LAW)
- R287: TODO persistence
