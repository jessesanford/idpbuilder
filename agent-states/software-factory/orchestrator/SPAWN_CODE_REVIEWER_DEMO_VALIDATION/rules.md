# Orchestrator State: SPAWN_CODE_REVIEWER_DEMO_VALIDATION

## Purpose
Spawn Code Reviewer to validate integration demos by RUNNING them and verifying
they execute successfully. This is the R291 GATE 4 enforcement mechanism.

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


## Entry Conditions
- Integration code review passed (APPROVED)
- Integration infrastructure exists
- Demo files should be present (created by integration agent)

## State Responsibilities

### 1. Verify Demo Infrastructure
Check that demos directory exists and demo scripts are present:

```bash
# Determine integration level and demo path
INTEGRATE_WAVE_EFFORTS_TYPE=$(jq -r '.current_integration_type' orchestrator-state-v3.json)
CURRENT_PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
CURRENT_WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

# Set demo paths based on level
if [ "$INTEGRATE_WAVE_EFFORTS_TYPE" = "wave" ]; then
    DEMO_DIR="demos/phase${CURRENT_PHASE}/wave${CURRENT_WAVE}/integration"
elif [ "$INTEGRATE_WAVE_EFFORTS_TYPE" = "phase" ]; then
    DEMO_DIR="demos/phase${CURRENT_PHASE}/integration"
elif [ "$INTEGRATE_WAVE_EFFORTS_TYPE" = "project" ]; then
    DEMO_DIR="demos/project"
else
    echo "🔴 Unknown integration type: $INTEGRATE_WAVE_EFFORTS_TYPE"
    exit 291
fi

# Verify demo infrastructure exists
if [ ! -d "$DEMO_DIR" ]; then
    echo "⚠️ WARNING: Demo directory missing: $DEMO_DIR"
    echo "This may indicate integration agent did not create demos per R291"
    echo "Code Reviewer will attempt to find and execute any demos present"
fi
```

### 2. Prepare Code Reviewer Spawn Instructions

Create clear instructions for the Code Reviewer to run demos:

```bash
# Create demo validation task file
TASK_FILE=".software-factory/demo-validation-task.json"

cat > "$TASK_FILE" << EOF
{
  "task_name": "Validate ${INTEGRATE_WAVE_EFFORTS_TYPE} integration demos",
  "task_description": "Run and verify integration demos execute successfully per R291 Gate 4",
  "integration_type": "${INTEGRATE_WAVE_EFFORTS_TYPE}",
  "demo_directory": "${DEMO_DIR}",
  "requirements": {
    "execute_all_demos": true,
    "capture_outputs": true,
    "verify_success": true,
    "create_report": true,
    "save_logs": true
  },
  "report_location": ".software-factory/phase${CURRENT_PHASE}/${INTEGRATE_WAVE_EFFORTS_TYPE}/demo-evaluation-report.md",
  "log_location": ".software-factory/phase${CURRENT_PHASE}/${INTEGRATE_WAVE_EFFORTS_TYPE}/demo-evaluation.log",
  "state_to_enter": "DEMO_VALIDATION",
  "spawned_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
}
EOF

echo "✅ Created demo validation task file: $TASK_FILE"
```

### 3. Spawn Code Reviewer for Demo Validation

```bash
# Spawn Code Reviewer with demo validation instructions
echo "🚀 Spawning Code Reviewer for demo validation..."

# Code Reviewer will:
# 1. Navigate to demos directory
# 2. Execute all demo scripts
# 3. Capture outputs
# 4. Verify successful execution
# 5. Create demo-evaluation-report.md
# 6. Save logs to .software-factory/.../demo-evaluation.log

# SPAWN COMMAND (actual implementation depends on agent spawning mechanism)
spawn_code_reviewer \
    --state DEMO_VALIDATION \
    --task-file "$TASK_FILE" \
    --integration-type "$INTEGRATE_WAVE_EFFORTS_TYPE" \
    --demo-directory "$DEMO_DIR"
```

### 4. Update Orchestrator State

Record the code reviewer spawn and transition to waiting state:

```bash
# Update state to reflect demo validation in progress
jq ".waiting_for = \"demo_validation\"" -i orchestrator-state-v3.json
jq ".demo_validation = {
  \"integration_type\": \"$INTEGRATE_WAVE_EFFORTS_TYPE\",
  \"demo_directory\": \"$DEMO_DIR\",
  \"validation_status\": \"in_progress\",
  \"code_reviewer_spawned_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
  \"evaluation_report\": \".software-factory/phase${CURRENT_PHASE}/${INTEGRATE_WAVE_EFFORTS_TYPE}/demo-evaluation-report.md\",
  \"evaluation_log\": \".software-factory/phase${CURRENT_PHASE}/${INTEGRATE_WAVE_EFFORTS_TYPE}/demo-evaluation.log\"
}" -i orchestrator-state-v3.json

# Commit state update
git add orchestrator-state-v3.json "$TASK_FILE"
git commit -m "state: spawn code reviewer for demo validation (R291 Gate 4)"
git push
```

## Exit Conditions

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE
Conditions for successful completion:
- Code Reviewer spawned successfully
- Demo infrastructure verified (or noted as missing)
- State updated to WAITING_FOR_DEMO_VALIDATION
- Task file created with clear instructions

**Next State**: WAITING_FOR_DEMO_VALIDATION

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE
Conditions requiring manual intervention:
- Demo infrastructure completely missing (no recovery possible)
- Cannot spawn Code Reviewer (infrastructure broken)
- State file corruption detected
- Integration type unknown/invalid

**Next State**: ERROR_RECOVERY (manual intervention required)

## Next States
- **WAITING_FOR_DEMO_VALIDATION** (normal path)
- **ERROR_RECOVERY** (if infrastructure fundamentally broken)

## R291 Compliance

This state enforces R291 GATE 4 (Demo Verification). It CANNOT be skipped.

**Per R291 line 44:**
> "Marking integration complete without passing build/test/demo = IMMEDIATE DISQUALIFICATION"

**Penalty for skipping this state**: -100% (AUTOMATIC FAILURE)

**State machine enforcement**:
- Direct transitions from code review to completion are PROHIBITED
- Demo validation is a MANDATORY step
- Cannot bypass via state manipulation

## State Machine Reference

This state is part of the R291 enforcement chain:

```
WAITING_FOR_REVIEW_WAVE_INTEGRATION (approved)
  ↓ (REQUIRED)
SPAWN_CODE_REVIEWER_DEMO_VALIDATION  ← YOU ARE HERE
  ↓ (REQUIRED)
WAITING_FOR_DEMO_VALIDATION
  ↓ (based on results)
  ├─ Demos passed → WAVE/PHASE/PROJECT_COMPLETE
  └─ Demos failed → ERROR_RECOVERY (MANDATORY)
```

## Critical Rules Referenced
- **R291**: Integration Demo Requirement (BLOCKING)
- **R322**: Mandatory Orchestrator Checkpoints (SUPREME LAW)
- **R263**: Integration Documentation Requirements
- **R265**: Integration Testing Requirements

## Remember

**Demos are not optional.** They are PROOF that the integration actually works.

Per R291:
- Build + Tests + Demo MUST all pass
- No exceptions
- No shortcuts
- No "we'll do it later"

**This state ensures that promise is kept.**

CONTINUE-SOFTWARE-FACTORY=TRUE


## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete SPAWN_CODE_REVIEWER_DEMO_VALIDATION:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="SPAWN_CODE_REVIEWER_DEMO_VALIDATION complete - [describe what was accomplished]"
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
  "state_completed": "SPAWN_CODE_REVIEWER_DEMO_VALIDATION",
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
  --current-state "SPAWN_CODE_REVIEWER_DEMO_VALIDATION" \
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
save_todos "SPAWN_CODE_REVIEWER_DEMO_VALIDATION_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - SPAWN_CODE_REVIEWER_DEMO_VALIDATION complete [R287]"; then
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

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
