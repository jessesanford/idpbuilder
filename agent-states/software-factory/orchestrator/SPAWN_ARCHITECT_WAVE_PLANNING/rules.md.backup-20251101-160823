# SPAWN_ARCHITECT_WAVE_PLANNING State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES.

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


## 📋 PRIMARY DIRECTIVES

### Core Mandatory Rules:

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE
2. **🔴🔴🔴 R287** - TODO PERSISTENCE COMPREHENSIVE  
3. **🔴🔴🔴 R288** - STATE FILE UPDATE REQUIREMENTS
4. **🔴🔴🔴 R510** - STATE EXECUTION CHECKLIST COMPLIANCE
5. **🔴🔴🔴 R405** - AUTOMATION CONTINUATION FLAG
6. **🔴🔴🔴 R322** - MANDATORY STOP BEFORE STATE TRANSITIONS

### State-Specific Rules:

7. **🚨🚨🚨 R340** - ARCHITECTURE PLAN QUALITY GATES
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R340-architecture-plan-quality-gates.md`
   - Summary: Wave architecture must be CONCRETE level (REAL CODE examples from previous waves)

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

### BLOCKING REQUIREMENTS

- [ ] 1. Spawn Architect agent for wave architecture planning
  - Agent: architect
  - Purpose: Create wave architecture with CONCRETE fidelity (REAL CODE examples)
  - Output File: `planning/phase{N}/wave{M}/WAVE-{N}.{M}-ARCHITECTURE.md` (per R550)
  - Template: `templates/WAVE-ARCHITECTURE-TEMPLATE.md`
  - Fidelity Level: CONCRETE (real code from previous waves, actual interfaces)
  - Parameters: --wave [wave-number] --template wave-architecture
  - State Tracking: Update `orchestrator-state-v3.json` field `planning_files.phases.phase{N}.waves.wave{M}.architecture_plan`
  - **BLOCKING**: Cannot proceed without wave architecture with real code

### STANDARD EXECUTION TASKS

- [ ] 2. Update state file with architect spawn metadata
- [ ] 3. Display guidance for Architect agent
  - Fidelity: CONCRETE (NOT pseudocode, use real code from previous waves)
  - Show examples of what "real code" means

### EXIT REQUIREMENTS

- [ ] 4. Set proposed next state: `WAITING_FOR_WAVE_ARCHITECTURE`
- [ ] 5. Spawn State Manager for SHUTDOWN_CONSULTATION
- [ ] 6. Save TODOs per R287
- [ ] 7. Set CONTINUE-SOFTWARE-FACTORY=TRUE
- [ ] 8. Stop execution (exit 0)

## State Purpose

Spawn Architect to create wave-level architecture with CONCRETE fidelity. Unlike phase architecture (pseudocode), wave architecture uses REAL CODE EXAMPLES from previously completed waves to show actual patterns and interfaces.

**Fidelity Level:** CONCRETE (real code examples, actual interfaces, working usage patterns)

## Fidelity Guidance

**Wave Architecture Must Include:**
- Real function signatures from previous waves: `def create_user(username: str, email: str) -> User:`
- Actual interface definitions with types
- Working code examples showing usage
- Real patterns observed in previous waves
- Concrete designs (not abstract patterns)

**Wave Architecture Must NOT:**
- Use pseudocode like "when X happens, do Y"
- Use abstract patterns without concrete code
- Skip showing actual interfaces
- Be high-level like phase architecture

**Example Progression:**
- Phase Architecture: "Use Factory pattern for user creation"
- **Wave Architecture**: `class UserFactory:\n    def create_user(self, username: str) -> User:\n        return User(username=username)`

## Entry Criteria

- **From**: WAVE_START (beginning new wave)
- **Condition**: Wave has no architecture plan yet
- **Required**: Phase architecture exists, previous wave(s) may provide code examples

## Exit Criteria

**Success** → WAITING_FOR_WAVE_ARCHITECTURE
**Failure** → ERROR_RECOVERY

## Rules Enforced

- R510, R511, R288, R287, R322, R405, R340, R006

## Additional Context

This is step 3 of SF 3.0 progressive planning:
1. Phase Architecture (pseudocode)
2. Phase Implementation (wave list)
3. **Wave Architecture** ← THIS STATE (concrete/real code)
4. Wave Implementation (exact specifications)

Architect should:
- Review previous wave implementations for real code
- Copy actual function signatures
- Show working usage examples
- Document real patterns observed
