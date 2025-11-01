# ORCHESTRATOR STATE: SPAWN_ARCHITECT_MASTER_PLANNING

## 🔴🔴🔴 CRITICAL: STATE FILES PROPOSE, STATE MANAGER DECIDES 🔴🔴🔴

**R288 SUPREME LAW - YOU NEVER CALL update_state DIRECTLY!**

This state file uses the **BOOKEND PATTERN** (R600):
1. **START**: Set `PROPOSED_NEXT_STATE` and `TRANSITION_REASON` variables
2. **WORK**: Execute state-specific logic
3. **END**: Follow 10-step completion checklist to exit properly

**NEVER CALL update_state() DIRECTLY - IT IS PROHIBITED!**
- ❌ `update_state "NEXT_STATE" "reason"` = SYSTEM VIOLATION
- ✅ `PROPOSED_NEXT_STATE="WAITING_FOR_MASTER_ARCHITECTURE"` = CORRECT
- ✅ `TRANSITION_REASON="reason"` = CORRECT

The State Manager (`run-software-factory.sh`) handles ALL state transitions.

**See:**
- `$CLAUDE_PROJECT_DIR/rule-library/R288-state-transition-authority.md` (SUPREME LAW)
- `$CLAUDE_PROJECT_DIR/rule-library/R600-orchestrator-bookend-pattern.md`

---

# PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R308** - Architect Integration Role
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R308-incremental-branching-strategy.md`

5. **R287** - TODO Persistence Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`

6. **R288** - State File Update Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`

7. **R304** - Mandatory Line Counter Usage (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`

8. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`

9. **R324** - State Transition Validation (SUPREME LAW)
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
This state spawns the Architect to create the master architecture that defines the entire project structure.

## Entry Criteria
- No existing PROJECT-ARCHITECTURE-PLAN.md at standard R550 location
- Starting a brand new project
- Transitioned from INIT state

## State Responsibilities

### 1. Verify No Existing Master Plan (R550 Compliant)
```bash
# R550: Use standardized planning_files location (NOT legacy planning_artifacts)
ARCH_PLAN_FILE=$(jq -r '.planning_files.project.architecture_plan // "planning/project/PROJECT-ARCHITECTURE-PLAN.md"' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")
ARCH_PLAN_PATH="$CLAUDE_PROJECT_DIR/$ARCH_PLAN_FILE"

# Check for existing architecture at R550 standardized location
if [ -f "$ARCH_PLAN_PATH" ]; then
    echo "⚠️ Project architecture plan already exists at: $ARCH_PLAN_PATH"
    PROPOSED_NEXT_STATE="WAITING_FOR_MASTER_ARCHITECTURE"
    TRANSITION_REASON="Project architecture plan already exists, skipping to test planning"
else
    echo "✅ No project architecture plan found at: $ARCH_PLAN_PATH"
    echo "✅ Proceeding with architect spawn per R550"
fi
```

### 2. Initialize Planning Files Metadata (R550 Pattern)
```bash
# R550: Initialize planning_files.project section (NEW standard)
# NOTE: planning_artifacts is deprecated but maintained for backward compatibility
jq '.planning_files = (.planning_files // {}) |
    .planning_files.project = (.planning_files.project // {}) |
    .planning_files.project.architecture_plan = "planning/project/PROJECT-ARCHITECTURE-PLAN.md"' \
   "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" > "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.tmp" && \
   mv "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.tmp" "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"

echo "✅ R550: Planning files metadata initialized for project architecture"
```

### 3. Spawn Architect Agent with R550 Standardized Output Location
```bash
# R550: Spawn architect with EXPLICIT instructions about R550-compliant output location
/spawn-agent architect \
    --state MASTER_PLANNING \
    --task "Create project architecture plan for entire project" \
    --output-location "$ARCH_PLAN_FILE" \
    --instructions "Create PROJECT-ARCHITECTURE-PLAN.md at planning/project/PROJECT-ARCHITECTURE-PLAN.md (R550 standardized location per orchestrator-state-v3.json planning_files.project.architecture_plan). This is the PROJECT-LEVEL architecture plan."
```

### 4. Record Spawn in State File
- Update orchestrator-state-v3.json with spawned agent
- Track timestamp per R151
- Record spawn ID

### 5. Transition to Waiting State
```bash
PROPOSED_NEXT_STATE="WAITING_FOR_MASTER_ARCHITECTURE"
TRANSITION_REASON="Spawned architect for master planning"
save_todos "Spawned architect for master planning"
git add orchestrator-state-v3.json todos/
git commit -m "state: spawned architect for master planning"
git push
```

## Exit Criteria
- Architect spawned successfully
- State file updated to WAITING_FOR_MASTER_ARCHITECTURE
- R313: MUST stop after spawning

## Success Metrics
- ✅ Architect spawned within 5s (R151)
- ✅ State file updated before stop (R324)
- ✅ TODOs saved (R287)



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete SPAWN_ARCHITECT_MASTER_PLANNING:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

**CRITICAL**: During state work, set these variables (DO NOT call update_state):
```bash
PROPOSED_NEXT_STATE="[state determined by logic]"
TRANSITION_REASON="[reason for transition]"
```

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Verify Proposed Transition Set
```bash
# Verify PROPOSED_NEXT_STATE was set during state work
if [ -z "$PROPOSED_NEXT_STATE" ]; then
    echo "❌ CRITICAL: PROPOSED_NEXT_STATE not set!"
    echo "State work must set PROPOSED_NEXT_STATE and TRANSITION_REASON"
    exit 288
fi

echo "✅ Proposed transition: SPAWN_ARCHITECT_MASTER_PLANNING → $PROPOSED_NEXT_STATE"
echo "   Reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Save Proposed Transition to File (R288)
```bash
# Write proposed transition to file for State Manager
cat > "$CLAUDE_PROJECT_DIR/.proposed-transition" <<EOF
PROPOSED_NEXT_STATE="$PROPOSED_NEXT_STATE"
TRANSITION_REASON="$TRANSITION_REASON"
CURRENT_STATE="SPAWN_ARCHITECT_MASTER_PLANNING"
TIMESTAMP="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
EOF

echo "✅ Proposed transition saved to .proposed-transition"
```

---

### ✅ Step 4: Validate Current State File (R324)
```bash
# Validate state file integrity BEFORE proposing transition
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
echo "✅ Current state file validated"
```

---

### ✅ Step 5: Commit Work Products (if any)
```bash
# Commit any work products created during this state
# (State file itself is committed by State Manager)
if [ -n "$(git status --porcelain | grep -v orchestrator-state-v3.json)" ]; then
    git add .
    git commit -m "work: SPAWN_ARCHITECT_MASTER_PLANNING work products [R288]"
    git push
    echo "✅ Work products committed and pushed"
else
    echo "✅ No work products to commit"
fi
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "SPAWN_ARCHITECT_MASTER_PLANNING_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - SPAWN_ARCHITECT_MASTER_PLANNING complete [R287]"; then
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

### ✅ Step 7: Output Transition Proposal (R288)
```bash
# Output the proposed transition for State Manager
echo "════════════════════════════════════════════════════════"
echo "PROPOSED STATE TRANSITION (R288):"
echo "  FROM: SPAWN_ARCHITECT_MASTER_PLANNING"
echo "  TO:   $PROPOSED_NEXT_STATE"
echo "  REASON: $TRANSITION_REASON"
echo "════════════════════════════════════════════════════════"
```

---

### ✅ Step 8: Output Continuation Flag (R405 - SUPREME LAW)
```bash
# Output continuation flag (R405)
# TRUE = state work complete, transition proposed
# FALSE = error occurred, manual intervention needed

if [ "$PROPOSED_NEXT_STATE" = "ERROR_RECOVERY" ]; then
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
else
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
fi
```

**⚠️ THIS MUST BE THE LAST OUTPUT BEFORE STOP! ⚠️**

---

### ✅ Step 9: State Manager Transition Checkpoint
```
🔄 STATE MANAGER TAKES CONTROL HERE (R288)

The State Manager (run-software-factory.sh) will:
1. Read .proposed-transition file
2. Validate proposed transition against state machine
3. Update orchestrator-state-v3.json with new state
4. Commit and push state file
5. Re-invoke orchestrator in new state (if CONTINUE-SOFTWARE-FACTORY=TRUE)

DO NOT PROCEED PAST THIS POINT - STATE MANAGER HANDLES TRANSITION
```

---

### ✅ Step 10: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for State Manager handoff (R322)
echo "🛑 State work complete - State Manager will handle transition"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 1: State work incomplete
- Missing Step 2: No proposed transition = stuck forever (R288 violation, -100%)
- Missing Step 3: State Manager can't read proposal = broken automation (R288 violation, -100%)
- Missing Step 4: Invalid state = corruption (R324 violation)
- Missing Step 5: Work products lost (data loss)
- Missing Step 6: TODOs lost on compaction (R287 violation, -20% to -100%)
- Missing Step 7: No transition visibility (R288 violation)
- Missing Step 8: Automation doesn't know if it can continue (R405 violation, -100%)
- Missing Step 9: State Manager handoff failed (R288 violation, -100%)
- Missing Step 10: Context corruption (R322 violation, -100%)

**ALL 10 STEPS ARE MANDATORY - NO EXCEPTIONS**

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
PROPOSED_NEXT_STATE="WAITING_FOR_MASTER_ARCHITECTURE"
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
- R151: Parallel agent timing
- R287: TODO persistence
- R313: Mandatory stop after spawn
- R324: State update before stop
- R341: TDD requirement

