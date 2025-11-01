# SPAWN_ARCHITECT_PHASE_PLANNING State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES.
**I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS**
*YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE!!!

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE
THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

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


## 📋 PRIMARY DIRECTIVES FOR SPAWN_ARCHITECT_PHASE_PLANNING STATE

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

4. **🔴🔴🔴 R510** - STATE EXECUTION CHECKLIST COMPLIANCE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R510-state-execution-checklist-compliance.md`
   - Criticality: SUPREME LAW
   - Summary: MUST complete and acknowledge every checklist item

5. **🔴🔴🔴 R405** - AUTOMATION CONTINUATION FLAG (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md`
   - Criticality: SUPREME LAW
   - Summary: MUST output CONTINUE-SOFTWARE-FACTORY flag before exit

6. **🔴🔴🔴 R322** - MANDATORY STOP BEFORE STATE TRANSITIONS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

7. **🚨🚨🚨 R340** - ARCHITECTURE PLAN QUALITY GATES (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R340-architecture-plan-quality-gates.md`
   - Criticality: BLOCKING - Defines pseudocode fidelity requirements
   - Summary: Phase architecture must be PSEUDOCODE level (high-level patterns, NO real code)

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. Spawn Architect agent for phase architecture planning
  - Agent: architect
  - Purpose: Create phase architecture plan with PSEUDOCODE fidelity
  - Output File: `planning/PHASE-N-ARCHITECTURE.md`
  - Template: `templates/PHASE-ARCHITECTURE-TEMPLATE.md`
  - Fidelity Level: PSEUDOCODE (high-level patterns, library choices, NO real code)
  - Parameters: --phase [phase-number] --template phase-architecture
  - Validation: Agent ID returned with timestamp
  - **BLOCKING**: Cannot proceed without phase architecture plan

### STANDARD EXECUTION TASKS (Required)

- [ ] 2. Update state file with architect spawn metadata
  - Fields: architect_id, architect_phase, architect_purpose, started_at
  - Source: Architect spawn response
  - Purpose: Track architect activity for monitoring

- [ ] 3. Display guidance for Architect agent
  - Show: Fidelity level (PSEUDOCODE)
  - Show: Expected outputs (patterns, library choices, adaptation notes)
  - Show: Template location
  - Expected: Clear guidance message

### EXIT REQUIREMENTS (Must complete before transition)

**NOTE**: These are STANDARD across ALL states - copy exactly

- [ ] 4. Set proposed next state and transition reason
  - Proposed Next State: `WAITING_FOR_PHASE_ARCHITECTURE`
  - Transition Reason: "Architect spawned for phase architecture planning"
  - Validation: Variables set correctly

- [ ] 5. Spawn State Manager for SHUTDOWN_CONSULTATION (SF 3.0)
  - Purpose: Validate transition and update all 4 state files atomically
  - Parameters: --current-state SPAWN_ARCHITECT_PHASE_PLANNING --proposed-next-state WAITING_FOR_PHASE_ARCHITECTURE
  - Work Results: Architect spawn metadata
  - Validation: State Manager returns REQUIRED_NEXT_STATE
  - **NOTE**: State Manager owns state file updates (R288)

- [ ] 6. Save TODOs per R287 (within 30s of last TodoWrite)
  - Trigger: "R510_CHECKLIST_COMPLETE"
  - Format: `todos/orchestrator-SPAWN_ARCHITECT_PHASE_PLANNING-[YYYYMMDD-HHMMSS].todo`
  - Validation: TODO file exists and contains current state
  - Commit: Within 60 seconds with message "todo: orchestrator - SPAWN_ARCHITECT_PHASE_PLANNING complete [R287]"

- [ ] 7. Set CONTINUE-SOFTWARE-FACTORY flag per R405
  - Value: `TRUE` (successful spawn, factory can continue)
  - Context: "Architect spawned for phase planning - factory continues"
  - **NOTE**: Agent stops (R322) but factory continues (R405)
  - **POSITION**: Must be LAST output line before exit

- [ ] 8. Stop execution (exit 0)
  - Command: `exit 0`
  - Timing: After ALL above items complete
  - Message: "🛑 Stopping per R322 - use /continue-orchestrating to resume"
  - Per: R322 (spawn state), R313

## State Purpose

This state spawns the Architect agent to create a phase-level architecture plan. The phase architecture provides PSEUDOCODE-level guidance (high-level patterns, library choices, design approaches) without concrete implementation details. This is the first step of the SF 3.0 progressive planning process.

**Primary Goal:** Spawn Architect to create phase architecture document
**Key Actions:** Spawn architect agent with phase-specific parameters
**Success Outcome:** Architect spawned successfully, ready to create phase architecture plan
**Fidelity Level:** PSEUDOCODE (NO real code, just patterns and approaches)

## Entry Criteria

- **From**: START_PHASE_ITERATION (when starting a new phase)
- **Condition**: No phase architecture plan exists for current phase
- **Required**:
  - orchestrator-state-v3.json exists with current phase number
  - PROJECT-IMPLEMENTATION-PLAN.md exists
  - templates/PHASE-ARCHITECTURE-TEMPLATE.md exists
  - Current phase has no existing planning/PHASE-N-ARCHITECTURE.md

## State Actions

### 1. Spawn Architect Agent for Phase Architecture Planning

**Implementation:**
- Read current phase number from orchestrator-state-v3.json
- Verify phase architecture file does not exist: `planning/PHASE-{N}-ARCHITECTURE.md`
- Spawn Architect agent with parameters:
  - `--phase [phase-number]`
  - `--template phase-architecture`
  - `--fidelity pseudocode`
- Capture agent ID with timestamp

**Fidelity Guidance for Architect:**
- **DO**: High-level design patterns (e.g., "Use Factory pattern for client creation")
- **DO**: Library/framework choices (e.g., "Use boto3 for AWS interactions")
- **DO**: Architecture approaches (e.g., "Implement async event handling")
- **DO**: Pseudocode examples (e.g., "when X happens, do Y")
- **DON'T**: Real function signatures (e.g., "def create_client(config: dict) -> Client:")
- **DON'T**: Actual code implementation
- **DON'T**: Detailed interface definitions

**Validation:**
- Agent ID returned in format: `agent-architect-[phase-id]-[YYYYMMDD-HHMMSS]`
- Spawn command exits with code 0
- Agent metadata recorded in state file

### 2. Update State Metadata

**Implementation:**
- Update orchestrator-state-v3.json (via State Manager in exit protocol)
- Record: architect_id, phase_number, spawn_timestamp, fidelity_level
- This metadata enables monitoring in WAITING_FOR_PHASE_ARCHITECTURE state

**Validation:**
- State file contains architect spawn metadata
- All required fields present

## Exit Criteria

### Success Path → WAITING_FOR_PHASE_ARCHITECTURE

- Architect agent spawned successfully
- Agent ID recorded in state file (via State Manager)
- TODOs saved per R287
- State transition performed by State Manager (SF 3.0)
- CONTINUE-SOFTWARE-FACTORY=TRUE flag set
- Execution stopped with exit 0

### Failure Scenarios

- **Architect Spawn Fails** → ERROR_RECOVERY
  - Condition: Spawn command returns non-zero exit code
  - Action: Log error details
  - Recovery: Review spawn parameters, check agent availability

- **Template Missing** → ERROR_RECOVERY
  - Condition: templates/PHASE-ARCHITECTURE-TEMPLATE.md does not exist
  - Action: Log missing template path
  - Recovery: Create template or restore from repository

- **Phase Number Invalid** → ERROR_RECOVERY
  - Condition: Current phase number not found in state file
  - Action: Log state file corruption
  - Recovery: Review orchestrator-state-v3.json integrity

## Rules Enforced

- R510: State Execution Checklist Compliance (this checklist)
- R511: Checklist Creation Protocol (checklist design)
- R288: State File Updates (via State Manager in SF 3.0)
- R287: TODO Persistence (save before transition)
- R322: Mandatory Stop After Spawn (must exit after spawning)
- R405: Automation Continuation Flag (set TRUE for normal completion)
- R340: Architecture Plan Quality Gates (pseudocode fidelity)
- R006: Orchestrator Never Writes Code (only spawns agents)

## Transition Rules

- **ALWAYS** → WAITING_FOR_PHASE_ARCHITECTURE (after successful spawn)
- **NEVER** → SPAWN_CODE_REVIEWER_PHASE_IMPL (must wait for architect first)
- **NEVER** → WAVE_START (cannot skip phase planning)
- **ERROR** → ERROR_RECOVERY (if spawn fails)

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at spawn states**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 spawn states = TRUE (99.9% of cases)**

### THE PATTERN AT THIS SPAWN STATE (SF 3.0)

```bash
# 1. Complete state work (spawn architect)
echo "✅ Architect spawned successfully"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="WAITING_FOR_PHASE_ARCHITECTURE"
TRANSITION_REASON="Architect spawned for phase architecture planning"

# 3. Spawn State Manager for state transition (SF 3.0)
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "SPAWN_ARCHITECT_PHASE_PLANNING" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"
# State Manager updates all 4 state files atomically

# 4. Save TODOs
save_todos "R510_CHECKLIST_COMPLETE"

# 5. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# 6. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ Architect spawned successfully
- ✅ Ready for WAITING_FOR_PHASE_ARCHITECTURE
- ✅ State work completed successfully
- ✅ R322 checkpoint reached (NORMAL)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for normal spawn failures (use ERROR_RECOVERY instead)

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

## Additional Context

### SF 3.0 Progressive Planning Process

Phase architecture is the FIRST step of SF 3.0's progressive refinement:

1. **Phase Architecture** (THIS STATE): PSEUDOCODE level
   - High-level patterns and approaches
   - Library/framework choices
   - NO real code

2. **Phase Implementation Plan** (next states): HIGH-LEVEL wave list
   - Wave names and descriptions
   - NO detailed effort definitions yet

3. **Wave Architecture** (later): CONCRETE designs
   - Real code examples from previous waves
   - Actual interface definitions

4. **Wave Implementation Plan** (later): EXACT specifications
   - Detailed effort definitions with R213 metadata
   - File lists and dependencies

### Adaptation Protocol

If this is NOT the first phase:
- Architect should review previous phase architecture
- Note what patterns worked well
- Document what changed or evolved
- Include adaptation notes section in output

### Common Pitfalls

- **DON'T** spawn with concrete fidelity - phase architecture is PSEUDOCODE only
- **DON'T** continue to next state - R322 requires stop after spawn
- **DON'T** set CONTINUE-SOFTWARE-FACTORY=FALSE - this is normal spawn (use TRUE)
- **DON'T** skip State Manager consultation - SF 3.0 requires bookend pattern
