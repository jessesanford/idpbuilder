# SPAWN_INTEGRATION_AGENT_PHASE State Rules

# PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R269** - WAVE Integration Merge Plan Protocol (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R269-code-reviewer-merge-plan-no-execution.md`

3. **R270** - PHASE Integration Merge Plan Protocol (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R270-no-integration-branches-as-sources.md`

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

## State Purpose
Spawn Integration Agent (SW Engineer in integration mode) to execute the phase merge plan, integrating all wave branches into the phase integration branch.

## Critical Rules

### 🔴🔴🔴 RULE R322: MANDATORY STOP BEFORE STATE TRANSITION (SUPREME LAW)
- **STOP** and save state before ANY transition
- **READ** orchestrator-state-v3.json to verify current state
- **VALIDATE** next state exists in software-factory-3.0-state-machine.json
- **VIOLATION = IMMEDIATE FAILURE**

### 🔴🔴🔴 RULE R322 Part A: MANDATORY STOP AFTER SPAWNING AGENTS (SUPREME LAW)
- **MUST STOP IMMEDIATELY** after spawning Integration Agent
- **RECORD** what was spawned in state file
- **SAVE** TODOs and commit state changes
- **EXIT** with clear continuation instructions
- **VIOLATION = CONTEXT LOSS AND RULE FORGETTING**

### 🚨🚨🚨 RULE R290: STATE RULE VERIFICATION (BLOCKING)
- **MUST** verify this rules file exists and is loaded
- **MUST** acknowledge all rules before proceeding
- **MUST** validate state transitions against state machine

### 🚨🚨🚨 RULE R208: SPAWN DIRECTORY PROTOCOL (BLOCKING)
- **MUST** spawn Integration Agent in phase integration directory
- **MUST** provide PHASE-MERGE-PLAN.md location
- **MUST** ensure agent has phase integration branch checked out
- **MUST** verify all wave branches are accessible

### 🚨🚨🚨 RULE R285: PHASE INTEGRATE_WAVE_EFFORTS REQUIREMENTS (BLOCKING)
- Integration MUST follow PHASE-MERGE-PLAN.md exactly
- Waves MUST be merged in sequential order
- Each merge MUST be tested before next merge
- Failed merges trigger IMMEDIATE_BACKPORT_REQUIRED per R321

### ⚠️⚠️⚠️ RULE R321: IMMEDIATE BACKPORT PROTOCOL (WARNING)
- ANY merge conflict or failure requires immediate fix in source branch
- Integration branches are READ-ONLY for code changes
- Fixes MUST go to effort branches first, then re-merge

### ⚠️⚠️⚠️ RULE R287: TODO PERSISTENCE (WARNING)
- **MUST** save TODOs before spawning
- **MUST** save TODOs after spawn completes
- **MUST** commit and push TODO state

## Required Actions

1. **Verify Phase Merge Plan Exists**
   ```bash
   # Check for merge plan
   ls -la phase-*-integration/PHASE-MERGE-PLAN.md
   
   # Verify plan contents
   cat phase-*-integration/PHASE-MERGE-PLAN.md
   ```

2. **Verify Phase Integration Infrastructure**
   ```bash
   # Check phase integration directory
   ls -la phase-*-integration/
   
   # Verify phase integration branch
   cd phase-*-integration/
   git branch --show-current  # Should be phase-X-integration
   git pull origin phase-X-integration
   ```

3. **Verify Wave Branches Ready**
   ```bash
   # List all wave branches for phase
   git ls-remote origin | grep "wave-${PHASE_NUM}-"
   
   # Verify all waves marked complete
   grep "waves_completed" orchestrator-state-v3.json
   ```

4. **Spawn Integration Agent**
   ```bash
   # Spawn SW Engineer as Integration Agent
   cd phase-X-integration/
   
   /spawn @agent-software-engineer INTEGRATE_PHASE_WAVES_EXECUTION \
     --merge-plan "PHASE-MERGE-PLAN.md" \
     --target-branch "phase-X-integration" \
     --wave-branches "wave-X-1,wave-X-2,..." \
     --output "PHASE-INTEGRATE_WAVE_EFFORTS-REPORT.md"
   ```

5. **Update State File**
   ```yaml
   current_state: SPAWN_INTEGRATION_AGENT_PHASE
   spawned_agents:
     - agent: sw-engineer
       role: integration-agent
       directory: phase-X-integration
       task: phase_integration_execution
       merge_plan: PHASE-MERGE-PLAN.md
       timestamp: YYYY-MM-DD HH:MM:SS
   phase_integration:
     status: IN_PROGRESS
     agent_spawned: true
   ```

6. **Save and Exit (R322 Part A MANDATORY)**
   ```bash
   # Save TODOs
   save_todos "SPAWN_INTEGRATION_AGENT_PHASE"
   
   # Commit state
   git add orchestrator-state-v3.json todos/*.todo
   git commit -m "state: spawned Integration Agent for phase integration"
   git push
   
   # EXIT IMMEDIATELY
   echo "Integration Agent spawned for phase merges. Use /continue orchestrator to resume."
   exit 0
   ```

## Transition Rules

### Valid Next States
- **MONITORING_INTEGRATE_PHASE_WAVES** - After spawning (MANDATORY per R322 Part A)

### Invalid Transitions
- ❌ Any state other than MONITORING_INTEGRATE_PHASE_WAVES
- ❌ Continuing work after spawn (violates R322 Part A)
- ❌ Attempting integration yourself (orchestrator doesn't merge)



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete SPAWN_INTEGRATION_AGENT_PHASE:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="SPAWN_INTEGRATION_AGENT_PHASE complete - [describe what was accomplished]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Update State File (R288 - SUPREME LAW)
```bash
# State Manager validates transition and updates state files (SF 3.0 Pattern)
echo "🔄 Spawning State Manager for SHUTDOWN_CONSULTATION..."

# Prepare work results summary
WORK_RESULTS=$(cat <<EOF
{
  "state_completed": "SPAWN_INTEGRATION_AGENT_PHASE",
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
  --current-state "SPAWN_INTEGRATION_AGENT_PHASE" \
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

### ✅ Step 4: Validate State File (R324)
```bash
# Validate state file before committing
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
echo "✅ State file validated"
```

---

### ✅ Step 5: Commit State File (R288)
```bash
# Commit and push state file immediately
git add orchestrator-state-v3.json

if ! git commit -m "state: SPAWN_INTEGRATION_AGENT_PHASE → $NEXT_STATE - SPAWN_INTEGRATION_AGENT_PHASE complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: SPAWN_INTEGRATION_AGENT_PHASE"
    echo "Attempted transition from: SPAWN_INTEGRATION_AGENT_PHASE"
    echo ""
    echo "Common causes:"
    echo "  - Schema validation failure (check pre-commit hook output above)"
    echo "  - Missing required fields in JSON files"
    echo "  - Invalid JSON syntax"
    echo ""
    echo "🛑 Cannot proceed - manual intervention required"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=SCHEMA_VALIDATION"
    exit 1
fi

git push || echo "⚠️ WARNING: Push failed - committed locally"
echo "✅ State file committed and pushed"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "SPAWN_INTEGRATION_AGENT_COMPLETE_PHASE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - SPAWN_INTEGRATION_AGENT_PHASE complete [R287]"; then
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

### ✅ Step 7: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 8: Stop Processing (R322 - SUPREME LAW)
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

**ALL 8 STEPS ARE MANDATORY - NO EXCEPTIONS**

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


## Common Violations to Avoid

1. **Not stopping after spawn** - Violates R322 Part A, causes context loss
2. **Missing PHASE-MERGE-PLAN.md** - Integration has no guidance
3. **Wrong directory for spawn** - Integration happens in wrong place
4. **Not verifying wave completeness** - Merging incomplete work
5. **Forgetting R321 immediate backport** - Creating fix commits in integration branch

## Verification Commands

```bash
# Verify state entry
echo "Entered SPAWN_INTEGRATION_AGENT_PHASE at $(date)"

# Verify merge plan exists
test -f phase-*/PHASE-MERGE-PLAN.md && echo "✓ Merge plan found" || echo "✗ Missing merge plan"

# Verify phase infrastructure
ls -la phase-*-integration/
git branch -r | grep "phase-.*-integration"

# Verify all waves complete
grep -c "status: COMPLETE" orchestrator-state-v3.json

# After spawn, verify stop
echo "STOPPING per R322 Part A - Integration Agent spawned for phase merges"
```

## Integration Failure Handling

If integration fails (detected in MONITORING_INTEGRATE_PHASE_WAVES):
1. Transition to **IMMEDIATE_BACKPORT_REQUIRED** (R321)
2. DO NOT attempt fixes in integration branch
3. Spawn engineers to fix source branches
4. Re-run entire phase integration after fixes

## References
- R322 Part A: rule-library/R322-mandatory-stop-before-state-transitions.md
- R208: rule-library/R208-orchestrator-spawn-cd-protocol.md
- R285: rule-library/R285-mandatory-phase-integration-before-assessment.md
- R321: rule-library/R321-immediate-backport-during-integration.md
- R287: rule-library/R287-todo-persistence-comprehensive.md
- R290: rule-library/R290-state-rule-reading-verification-supreme-law.md
- R322: rule-library/R322-mandatory-stop-before-state-transitions.md

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

