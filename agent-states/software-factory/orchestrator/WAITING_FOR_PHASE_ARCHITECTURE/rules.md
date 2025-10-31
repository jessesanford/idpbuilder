# WAITING_FOR_PHASE_ARCHITECTURE State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES.
**I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS**
*YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE!!!

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE
THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_PHASE_ARCHITECTURE STATE

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
   - Criticality: SUPREME LAW - Must stop at checkpoints
   - Summary: Waiting states are checkpoints requiring stop

### State-Specific Rules:

7. **🚨🚨🚨 R233** - ACTIVE MONITORING PROTOCOL (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-active-monitoring-protocol.md`
   - Criticality: BLOCKING - Defines monitoring requirements for waiting states
   - Summary: MUST actively monitor Architect progress every 5 messages

8. **🚨🚨🚨 R340** - ARCHITECTURE PLAN QUALITY GATES (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R340-architecture-plan-quality-gates.md`
   - Criticality: BLOCKING - Defines quality requirements for phase architecture
   - Summary: Phase architecture must be PSEUDOCODE level with specific quality gates

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. Verify phase architecture plan file exists
  - File Path: `planning/PHASE-N-ARCHITECTURE.md`
  - Check: File exists and is not empty
  - Validation: `test -f planning/PHASE-N-ARCHITECTURE.md && [ -s planning/PHASE-N-ARCHITECTURE.md ]`
  - **BLOCKING**: Cannot proceed without completed phase architecture

- [ ] 2. Validate phase architecture quality per R340
  - Check: File contains pseudocode patterns (NOT real code)
  - Check: File contains library/framework choices
  - Check: File contains adaptation notes (if not first phase)
  - Check: File does NOT contain real function signatures or implementations
  - Validation: `grep -q "pseudocode\|pattern\|library\|framework" planning/PHASE-N-ARCHITECTURE.md`
  - **BLOCKING**: Quality gates must pass before using plan

### STANDARD EXECUTION TASKS (Required)

- [ ] 3. Display phase architecture summary
  - Show: Key patterns and approaches documented
  - Show: Library choices justified
  - Show: Adaptation notes present (if applicable)
  - Expected: Clear summary message for user review

- [ ] 4. Display next steps guidance
  - Show: Next state will be SPAWN_CODE_REVIEWER_PHASE_IMPL
  - Show: Code Reviewer will create wave list (high-level only)
  - Show: Use `/continue-orchestrating` to proceed
  - Expected: Clear guidance message

### EXIT REQUIREMENTS (Must complete before transition)

**NOTE**: These are STANDARD across ALL states - copy exactly

- [ ] 5. Set proposed next state and transition reason
  - Proposed Next State: `SPAWN_CODE_REVIEWER_PHASE_IMPL`
  - Transition Reason: "Phase architecture complete and validated"
  - Validation: Variables set correctly

- [ ] 6. Spawn State Manager for SHUTDOWN_CONSULTATION (SF 3.0)
  - Purpose: Validate transition and update all 4 state files atomically
  - Parameters: --current-state WAITING_FOR_PHASE_ARCHITECTURE --proposed-next-state SPAWN_CODE_REVIEWER_PHASE_IMPL
  - Work Results: Architecture validation results
  - Validation: State Manager returns REQUIRED_NEXT_STATE
  - **NOTE**: State Manager owns state file updates (R288)

- [ ] 7. Save TODOs per R287 (within 30s of last TodoWrite)
  - Trigger: "R510_CHECKLIST_COMPLETE"
  - Format: `todos/orchestrator-WAITING_FOR_PHASE_ARCHITECTURE-[YYYYMMDD-HHMMSS].todo`
  - Validation: TODO file exists and contains current state
  - Commit: Within 60 seconds with message "todo: orchestrator - WAITING_FOR_PHASE_ARCHITECTURE complete [R287]"

- [ ] 8. Set CONTINUE-SOFTWARE-FACTORY flag per R405
  - Value: `TRUE` (architecture validated, factory can continue)
  - Context: "Phase architecture validated - proceeding to implementation planning"
  - **NOTE**: R322 checkpoint - agent stops but factory continues
  - **POSITION**: Must be LAST output line before exit

- [ ] 9. Stop execution (exit 0)
  - Command: `exit 0`
  - Timing: After ALL above items complete
  - Message: "🛑 Stopping per R322 - use /continue-orchestrating to resume"
  - Per: R322 (checkpoint state)

## State Purpose

This state monitors the Architect agent during phase architecture planning and validates the resulting plan meets R340 quality gates (PSEUDOCODE fidelity). It is a WAITING checkpoint where the orchestrator verifies the architecture before proceeding to implementation planning.

**Primary Goal:** Validate phase architecture plan meets quality requirements
**Key Actions:** Monitor Architect progress, validate plan quality per R340
**Success Outcome:** Phase architecture validated, ready for phase implementation planning
**Quality Gates:** Pseudocode patterns present, NO real code, library choices justified

## Entry Criteria

- **From**: SPAWN_ARCHITECT_PHASE_PLANNING (after Architect spawned)
- **Condition**: Architect is working on phase architecture plan
- **Required**:
  - Architect agent spawned successfully in previous state
  - orchestrator-state-v3.json contains architect_id metadata
  - planning/ directory exists

## State Actions

### 1. Active Monitoring (R233)

**Implementation:**
- Check Architect progress every 5 messages (R233)
- Monitor for planning/PHASE-N-ARCHITECTURE.md file creation
- Check if Architect is stuck or needs guidance
- Provide status updates to user

**Validation:**
- Architect has not stalled
- Work is progressing normally
- File will be created soon or already exists

### 2. Phase Architecture Quality Validation (R340)

**Implementation:**
When planning/PHASE-N-ARCHITECTURE.md exists:
- **Check 1**: File contains pseudocode patterns (NOT real code)
  - Look for: High-level design descriptions
  - Look for: Pattern names (Factory, Observer, etc.)
  - Reject: Real function definitions like `def foo(x: int) -> str:`
  
- **Check 2**: File contains library/framework choices
  - Look for: Specific libraries mentioned (boto3, requests, FastAPI, etc.)
  - Look for: Justification for choices
  
- **Check 3**: File contains adaptation notes (if not first phase)
  - Look for: "Adaptation Notes" section
  - Look for: References to previous phase architecture
  - Look for: What changed or evolved
  
- **Check 4**: File does NOT contain implementation details
  - Reject: Real code snippets
  - Reject: Detailed function signatures
  - Reject: Exact interface definitions (those come in wave architecture)

**Validation:**
- All 4 quality checks pass
- R340 requirements met
- Ready for phase implementation planning

### 3. Display Status and Next Steps

**Implementation:**
- Show user what was validated
- Explain next state (SPAWN_CODE_REVIEWER_PHASE_IMPL)
- Clarify fidelity will remain HIGH-LEVEL (wave list only, no detailed plans yet)

**Validation:**
- User sees clear status message
- User understands next steps

## Exit Criteria

### Success Path → SPAWN_CODE_REVIEWER_PHASE_IMPL

- Phase architecture file exists: `planning/PHASE-N-ARCHITECTURE.md`
- R340 quality gates passed (pseudocode fidelity confirmed)
- TODOs saved per R287
- State transition performed by State Manager (SF 3.0)
- CONTINUE-SOFTWARE-FACTORY=TRUE flag set
- Execution stopped with exit 0

### Failure Scenarios

- **Architecture File Missing** → ERROR_RECOVERY
  - Condition: Architect completed but no file created
  - Action: Log issue, check Architect output
  - Recovery: Re-spawn Architect or investigate failure

- **Quality Gates Failed** → ERROR_RECOVERY
  - Condition: File exists but contains real code (not pseudocode)
  - Action: Log validation failures
  - Recovery: Request Architect to revise plan to pseudocode level

- **Architect Stalled** → ERROR_RECOVERY
  - Condition: No progress after extended time (R233 monitoring detects)
  - Action: Check Architect status, logs
  - Recovery: Terminate and re-spawn Architect

## Rules Enforced

- R510: State Execution Checklist Compliance (this checklist)
- R511: Checklist Creation Protocol (checklist design)
- R288: State File Updates (via State Manager in SF 3.0)
- R287: TODO Persistence (save before transition)
- R322: Mandatory Stop at Checkpoints (must exit with checkpoint message)
- R405: Automation Continuation Flag (set TRUE for normal completion)
- R233: Active Monitoring Protocol (monitor Architect every 5 messages)
- R340: Architecture Plan Quality Gates (validate pseudocode fidelity)
- R006: Orchestrator Never Writes Code (only validates plans)

## Transition Rules

- **ALWAYS** → SPAWN_CODE_REVIEWER_PHASE_IMPL (after architecture validated)
- **NEVER** → WAVE_START (must create phase implementation plan first)
- **NEVER** → SPAWN_ARCHITECT_WAVE_PLANNING (phase planning must complete first)
- **ERROR** → ERROR_RECOVERY (if validation fails or Architect stalls)

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoint/waiting states**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoint states = TRUE (99.9% of cases)**

### THE PATTERN AT THIS CHECKPOINT STATE (SF 3.0)

```bash
# 1. Complete state work (validate architecture)
echo "✅ Phase architecture validated per R340"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="SPAWN_CODE_REVIEWER_PHASE_IMPL"
TRANSITION_REASON="Phase architecture complete and validated"

# 3. Spawn State Manager for state transition (SF 3.0)
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "WAITING_FOR_PHASE_ARCHITECTURE" \
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
- ✅ Phase architecture validated successfully
- ✅ Ready for SPAWN_CODE_REVIEWER_PHASE_IMPL
- ✅ State work completed successfully
- ✅ R322 checkpoint reached (NORMAL)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for validation failures (use ERROR_RECOVERY instead)

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

## Additional Context

### SF 3.0 Progressive Planning - Phase Level

This state validates the FIRST artifact in SF 3.0's progressive planning:

**Phase Architecture (THIS STATE validates)**: PSEUDOCODE level
- High-level patterns and approaches
- Library/framework choices  
- NO real code, NO concrete designs

**Next Steps**:
1. SPAWN_CODE_REVIEWER_PHASE_IMPL: Create wave list (names + descriptions only)
2. WAITING_FOR_PHASE_IMPLEMENTATION_PLAN: Validate wave list
3. WAVE_START: Begin first wave (which will have its own planning)

### R340 Quality Gate Examples

**✅ PASSES R340 (Pseudocode Level)**:
```
Use Factory pattern to create cloud provider clients
Implement async event handling for notifications
Use boto3 library for AWS operations
Pattern: Observer for state change notifications
```

**❌ FAILS R340 (Too Concrete - Real Code)**:
```python
def create_client(provider: str, config: dict) -> CloudClient:
    if provider == "aws":
        return boto3.client('s3', **config)
    elif provider == "gcp":
        return storage.Client(**config)
```

**❌ FAILS R340 (Too Concrete - Exact Interfaces)**:
```
class CloudClient:
    def __init__(self, config: dict) -> None
    def upload(self, key: str, data: bytes) -> UploadResult
    def download(self, key: str) -> bytes
```

### Common Pitfalls

- **DON'T** accept architecture with real code - must be pseudocode only
- **DON'T** skip R340 validation - quality gates are mandatory
- **DON'T** set CONTINUE-SOFTWARE-FACTORY=FALSE for validation failures - use ERROR_RECOVERY
- **DON'T** skip State Manager consultation - SF 3.0 requires bookend pattern
- **DON'T** continue to next state without checkpoint stop - R322 violation
