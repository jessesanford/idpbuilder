# SPAWN_CODE_REVIEWER_WAVE_IMPL State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

## 📋 PRIMARY DIRECTIVES

### Core Mandatory Rules:
1. R006, R287, R288, R510, R405, R322

### State-Specific Rules:
7. **🚨🚨🚨 R502** - Implementation Plan Quality Gates
8. **🚨🚨🚨 R213** - Effort Definition Metadata Requirements
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R213-effort-metadata-schema.md`
   - Summary: Wave implementation plans must include R213 metadata for each effort

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


## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

### BLOCKING REQUIREMENTS

- [ ] 1. Spawn Code Reviewer agent for wave implementation planning
  - Agent: code-reviewer
  - Purpose: Create wave implementation plan with EXACT fidelity + R213 metadata
  - Output: `planning/phase{N}/wave{M}/WAVE-IMPLEMENTATION-PLAN.md` (per R550)
  - Template: `templates/WAVE-IMPLEMENTATION-TEMPLATE.md`
  - Fidelity: EXACT (detailed effort definitions, file lists, R213 metadata)
  - State Tracking: Update `orchestrator-state-v3.json` field `planning_files.phases.phase{N}.waves.wave{M}.implementation_plan`
  - **BLOCKING**: Must have detailed efforts with R213 metadata

### STANDARD EXECUTION TASKS

- [ ] 2. Update state file with reviewer spawn metadata
- [ ] 3. Display guidance for Code Reviewer
  - Fidelity: EXACT (NOT wave list, include ALL details)
  - Show R213 metadata requirements
  - Show effort definition format

### EXIT REQUIREMENTS

- [ ] 4. Set proposed next state: `WAITING_FOR_IMPLEMENTATION_PLAN`
- [ ] 5. Spawn State Manager for SHUTDOWN_CONSULTATION
- [ ] 6. Save TODOs per R287
- [ ] 7. Set CONTINUE-SOFTWARE-FACTORY=TRUE
- [ ] 8. Stop execution (exit 0)

## State Purpose

Spawn Code Reviewer to create wave implementation plan with EXACT fidelity. This includes detailed effort definitions, file lists, dependencies, and R213 metadata for each effort.

**Fidelity Level:** EXACT (complete specifications, R213 metadata)

## Fidelity Guidance

**Wave Implementation Plan Must Include:**
- Detailed effort definitions (e.g., "Effort 1.1: Create User Model")
- File paths for each effort (e.g., "src/models/user.py")
- R213 metadata block for EACH effort:
  ```yaml
  effort_id: "1.1"
  estimated_lines: 150
  dependencies: ["effort:1.0"]
  files_touched: ["src/models/user.py", "tests/test_user.py"]
  branch_name: "effort/1.1-user-model"
  ```
- Complete task lists
- Exact specifications (not general guidance)

**Wave Implementation Must NOT:**
- Be wave-level only (that was phase implementation)
- Skip R213 metadata
- Use general descriptions
- Miss file paths or dependencies

## Entry Criteria

- **From**: WAITING_FOR_WAVE_ARCHITECTURE
- **Required**: Wave architecture validated, ready for detailed planning

## Exit Criteria

**Success** → WAITING_FOR_IMPLEMENTATION_PLAN
**Failure** → ERROR_RECOVERY

## Rules Enforced

- R510, R502, R213, R288, R287, R322, R405, R006

## Additional Context

SF 3.0 Progressive Planning Final Step:
1. Phase Architecture (pseudocode) ✅
2. Phase Implementation (wave list) ✅
3. Wave Architecture (real code) ✅
4. **Wave Implementation** ← THIS STATE (exact specs + R213)

This is the MOST DETAILED planning level:
- Exact effort definitions
- Complete file lists
- R213 metadata for tracking
- Dependency graphs
- Branch naming schemes

Common Mistakes:
- Missing R213 metadata → BLOCKING
- Wave-level descriptions → Wrong fidelity
- No file paths → Incomplete
