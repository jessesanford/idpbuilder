# WAITING_FOR_PHASE_IMPLEMENTATION_PLAN State Rules

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

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_PHASE_IMPLEMENTATION_PLAN STATE

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

### State-Specific Rules:

6. **🚨🚨🚨 R502** - IMPLEMENTATION PLAN QUALITY GATES (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R502-plan-quality-validation.md`
   - Criticality: BLOCKING - Validates phase implementation plan quality
   - Summary: Phase implementation plan must contain wave list ONLY (NO detailed effort definitions)

7. **🚨🚨🚨 R233** - ACTIVE MONITORING PROTOCOL (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-active-monitoring-protocol.md`
   - Criticality: BLOCKING - Requires active agent monitoring
   - Summary: MUST actively monitor spawned agents, not passive waiting

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. Wait for Code Reviewer to complete phase implementation plan
  - Agent: code-reviewer (spawned in SPAWN_CODE_REVIEWER_PHASE_IMPL)
  - Expected Output: `planning/PHASE-N-IMPLEMENTATION.md`
  - Polling Interval: Every 30 seconds per R233
  - Timeout: 30 minutes (escalate to ERROR_RECOVERY)
  - **BLOCKING**: Cannot proceed without completed wave list

- [ ] 2. Validate phase implementation plan exists
  - File: `planning/PHASE-{N}-IMPLEMENTATION.md`
  - Check: File exists and is not empty
  - Minimum Size: >500 bytes (basic validation)
  - **BLOCKING**: File must exist before quality validation

- [ ] 3. Validate phase implementation plan quality (R502)
  - Check: Contains 3-8 wave definitions
  - Check: Each wave has name and description
  - Check: NO effort definitions present (e.g., "Effort 1.1")
  - Check: NO R213 metadata present
  - Check: NO file lists or detailed tasks
  - **BLOCKING**: Must meet fidelity requirements (HIGH-LEVEL only)

- [ ] 4. Extract wave list for state file
  - Parse: Wave names from implementation plan
  - Format: Array of {wave_number, wave_name, wave_description}
  - Store: In orchestrator-state-v3.json for WAVE_START usage
  - **BLOCKING**: Wave list needed for wave iteration

### STANDARD EXECUTION TASKS (Required)

- [ ] 5. Record validation results in state file
  - Fields: plan_validated, validation_timestamp, wave_count
  - Source: Validation checks above
  - Purpose: Audit trail for planning quality

- [ ] 6. Display validation summary
  - Show: Wave count, fidelity level confirmed
  - Show: Validation timestamp
  - Show: Next state (WAVE_START typically)
  - Expected: Clear summary message

### EXIT REQUIREMENTS (Must complete before transition)

- [ ] 7. Set proposed next state and transition reason
  - Proposed Next State: `WAVE_START` (typically, or phase-specific logic)
  - Transition Reason: "Phase implementation plan validated - wave list extracted"
  - Validation: Variables set correctly

- [ ] 8. Spawn State Manager for SHUTDOWN_CONSULTATION (SF 3.0)
  - Purpose: Validate transition and update all 4 state files atomically
  - Parameters: --current-state WAITING_FOR_PHASE_IMPLEMENTATION_PLAN --proposed-next-state [NEXT_STATE]
  - Work Results: Validation results, wave list
  - Validation: State Manager returns REQUIRED_NEXT_STATE

- [ ] 9. Save TODOs per R287 (within 30s of last TodoWrite)
  - Trigger: "R510_CHECKLIST_COMPLETE"
  - Format: `todos/orchestrator-WAITING_FOR_PHASE_IMPLEMENTATION_PLAN-[YYYYMMDD-HHMMSS].todo`
  - Commit: Within 60 seconds

- [ ] 10. Set CONTINUE-SOFTWARE-FACTORY flag per R405
  - Value: `TRUE` (normal completion)
  - Context: "Phase implementation plan validated - proceeding to waves"
  - **POSITION**: Must be LAST output line before exit

- [ ] 11. Stop execution (exit 0)
  - Command: `exit 0`
  - Message: "✅ Phase implementation plan validated - ready for wave execution"
  - Per: R313

## State Purpose

This state actively monitors the Code Reviewer agent creating the phase implementation plan, validates the plan meets R502 quality gates (wave list ONLY, NO detailed efforts), and extracts the wave list for subsequent wave iteration.

**Primary Goal:** Validate phase implementation plan quality and extract wave list
**Key Actions:** Monitor agent, validate fidelity, extract waves
**Success Outcome:** Valid wave list ready for WAVE_START
**Fidelity Level Enforced:** HIGH-LEVEL (wave list only)

## Entry Criteria

- **From**: SPAWN_CODE_REVIEWER_PHASE_IMPL (after spawning code reviewer)
- **Condition**: Code Reviewer agent working on phase implementation plan
- **Required**:
  - orchestrator-state-v3.json contains reviewer_id
  - planning/PHASE-N-ARCHITECTURE.md exists (validated previously)
  - Code Reviewer agent spawned successfully

## State Actions

### 1. Active Monitoring (R233)

**Implementation:**
- Poll for planning/PHASE-{N}-IMPLEMENTATION.md existence every 30 seconds
- Check agent status if available
- Timeout after 30 minutes → ERROR_RECOVERY
- Log monitoring activity

**Validation:**
- File appears within timeout period
- File size >500 bytes minimum

### 2. Quality Validation (R502)

**Implementation:**
- Read planning/PHASE-{N}-IMPLEMENTATION.md
- Count wave definitions (expect 3-8 typically)
- Verify each wave has name + description
- Scan for forbidden content:
  - "Effort" or "effort" (indicates wrong fidelity)
  - File paths (e.g., src/models/user.py)
  - R213 metadata blocks
  - Detailed task lists

**Validation Criteria:**
- ✅ Contains 3-8 waves
- ✅ Each wave has name and description (2+ sentences)
- ✅ NO effort definitions
- ✅ NO R213 metadata
- ✅ HIGH-LEVEL descriptions only

**Failure Handling:**
- If fidelity wrong → Log issues, transition to ERROR_RECOVERY
- If too few/many waves → Flag for review
- If missing descriptions → Request revision

### 3. Wave List Extraction

**Implementation:**
- Parse wave names from implementation plan
- Extract descriptions
- Detect dependencies (e.g., "Depends on: Wave 1")
- Create structured wave list:

```json
{
  "waves": [
    {
      "wave_number": 1,
      "wave_name": "Core API Foundation",
      "wave_description": "Create the basic API structure...",
      "depends_on": []
    },
    {
      "wave_number": 2,
      "wave_name": "Business Logic Layer",
      "wave_description": "Implement core business rules...",
      "depends_on": [1]
    }
  ]
}
```

- Store in orchestrator-state-v3.json (via State Manager)

**Validation:**
- All waves extracted
- Dependencies parsed correctly
- Data structure valid

## Exit Criteria

### Success Path → WAVE_START (typically)

- Phase implementation plan exists and validated
- Wave list extracted successfully
- State file updated with wave list (via State Manager)
- TODOs saved per R287
- CONTINUE-SOFTWARE-FACTORY=TRUE flag set
- Execution stopped with exit 0

### Failure Scenarios

- **Plan Never Appears** → ERROR_RECOVERY
  - Condition: 30 minute timeout exceeded
  - Action: Check Code Reviewer agent status
  - Recovery: May need to respawn Code Reviewer

- **Fidelity Violation** → ERROR_RECOVERY
  - Condition: Plan contains effort definitions or R213 metadata
  - Action: Log specific violations
  - Recovery: Request Code Reviewer revision with corrected guidance

- **No Waves Found** → ERROR_RECOVERY
  - Condition: Plan exists but no waves detected
  - Action: Log parsing error
  - Recovery: Check plan format, may need manual review

## Rules Enforced

- R510: State Execution Checklist Compliance (this checklist)
- R502: Implementation Plan Quality Gates (HIGH-LEVEL fidelity)
- R233: Active Monitoring Protocol (30s polling)
- R288: State File Updates (via State Manager)
- R287: TODO Persistence
- R405: Automation Continuation Flag
- R006: Orchestrator Never Writes Code

## Transition Rules

- **TYPICALLY** → WAVE_START (begin first wave)
- **ALTERNATIVE** → ERROR_RECOVERY (if validation fails)
- **NEVER** → SPAWN_CODE_REVIEWERS_EFFORT_PLANNING (skip wave architecture/implementation)

## Additional Context

### Wave List Format Example

```markdown
# Phase 1 Implementation Plan

## Wave 1: Core API Foundation
Create the basic FastAPI application structure with authentication endpoints and user data models. This establishes the foundation for all subsequent development.

## Wave 2: Business Logic Layer  
Implement core business rules, validation logic, and service layer on top of the API foundation created in Wave 1. Depends on: Wave 1

## Wave 3: Integration & Testing
Connect external systems, implement comprehensive integration tests, and validate end-to-end workflows. Requires Waves 1-2 complete.
```

### Progressive Planning Context

This is step 2 validation of SF 3.0 progressive planning:
1. Phase Architecture (validated previously)
2. **Phase Implementation Plan** ← THIS STATE validates
3. Wave Architecture (comes next for each wave)
4. Wave Implementation Plan (comes after wave architecture)

### Common Validation Failures

1. **Too Much Detail**: Plan includes effort definitions
   - Fix: Code Reviewer must remove detailed efforts, keep wave-level only

2. **Wrong Fidelity**: Plan includes R213 metadata or file paths
   - Fix: Code Reviewer must remove metadata, use HIGH-LEVEL descriptions

3. **Missing Dependencies**: Waves don't specify what they depend on
   - Fix: Not blocking, but helpful to document

### Timeout Handling

If Code Reviewer takes >30 minutes:
1. Check if agent is still running
2. Check if there are error messages
3. Consider respawning with clearer guidance
4. Escalate to human if repeated failures
