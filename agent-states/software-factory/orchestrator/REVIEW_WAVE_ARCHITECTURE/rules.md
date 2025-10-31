# REVIEW_WAVE_ARCHITECTURE State Rules


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

---

## 📋 PRIMARY DIRECTIVES FOR REVIEW_WAVE_ARCHITECTURE STATE

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
   - Summary: MUST update all 4 state files atomically before EVERY state transition

4. **🔴🔴🔴 R510** - STATE EXECUTION CHECKLIST COMPLIANCE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R510-state-execution-checklist-compliance.md`
   - Criticality: SUPREME LAW
   - Summary: MUST complete and acknowledge every checklist item

5. **🔴🔴🔴 R405** - AUTOMATION CONTINUATION FLAG (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md`
   - Criticality: SUPREME - Required for all states
   - Summary: MUST set CONTINUE-SOFTWARE-FACTORY flag as last line of output

### State-Specific Rules:

6. **🚨🚨🚨 R232** - Spawn Agent Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R232-spawn-agent-protocol.md`
   - Criticality: BLOCKING
   - Summary: Spawn architect agent with proper context and workspace

7. **🔴🔴🔴 R233** - Architect Review Requirements
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-architect-review-requirements.md`
   - Criticality: SUPREME LAW
   - Summary: Architect validates pattern compliance and system coherence

---

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### PRE-EXECUTION: Check for Post-CASCADE Scenario (BUG #4 FIX)

**Before executing checklist, determine if this is a post-CASCADE scenario:**

```bash
# Check for post-CASCADE scenario
PREVIOUS_STATE=$(jq -r '.state_machine.state_history[-2].to_state // "UNKNOWN"' orchestrator-state-v3.json)

if [ "$PREVIOUS_STATE" = "CASCADE_REINTEGRATION" ]; then
    echo "✅ POST-CASCADE SCENARIO DETECTED"
    echo "   Previous state: CASCADE_REINTEGRATION"
    echo "   Wave architecture comes from cascaded integration - using alternate validation"

    # Set flag for alternate validation path
    POST_CASCADE_MODE=true

    # In post-CASCADE mode:
    # - Skip bug-free verification (CASCADE already validated)
    # - Skip build/test verification (CASCADE already validated)
    # - Skip architect spawn (architecture was reviewed during CASCADE)
    # - Transition directly to completion

    echo "📋 Post-CASCADE checklist mode enabled"
    echo "   Skipping redundant validation checks"
    echo "   Transitioning to: BUILD_VALIDATION (CASCADE validation sufficient)"

    # Update state file and transition
    PROPOSED_NEXT_STATE="BUILD_VALIDATION"
    TRANSITION_REASON="Post-CASCADE architecture review - proceeding to build validation"

    # Set continuation flag and exit
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
    exit 0
else
    echo "📋 Standard wave architecture review (not post-CASCADE)"
    POST_CASCADE_MODE=false
fi
```

**If POST_CASCADE_MODE=true**, skip to exit requirements. Otherwise, continue with standard checklist below.

---

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. Verify integration is bug-free
  - Check: bugs_found == 0 from REVIEW_WAVE_INTEGRATION
  - Source: integration-containers.json iteration_history[N].bugs_found
  - Validation: Most recent iteration has bugs_found = 0
  - **BLOCKING**: Cannot proceed to architecture review with bugs
  - **POST-CASCADE**: ⚠️ SKIP if POST_CASCADE_MODE=true

- [ ] 2. Verify build and tests pass on integration branch
  - Branch: wave-${phase_id}-${wave_id}-integration
  - Check: Build succeeds
  - Check: All tests pass (unit, integration, smoke)
  - Validation: Build exit code = 0, test suite green
  - **BLOCKING**: Cannot review broken code

- [ ] 3. Spawn Architect agent for wave integration assessment
  - Agent: architect
  - Context: Wave integration branch, pattern compliance review
  - Task: WAVE_INTEGRATE_WAVE_EFFORTS_REVIEW
  - Validation: Architect agent spawned, returns agent ID
  - **BLOCKING**: Architecture review requires architect

- [ ] 4. Wait for Architect to complete assessment
  - Monitor: Architect progress
  - Timeout: Reasonable time limit
  - Output: Architecture assessment report
  - Validation: Assessment report exists with decision
  - **BLOCKING**: Cannot proceed without architect decision

- [ ] 5. Record architect decision in orchestrator-state-v3.json
  - Decision values: PROCEED or CHANGES_REQUIRED
  - Record: Architect agent ID, assessment timestamp
  - Validation: Decision recorded
  - **BLOCKING**: Decision determines next state

### STANDARD EXECUTION TASKS (Required)

- [ ] 6. Update integration-containers.json with architecture review outcome
  - Field: iteration_history[N].architecture_review
  - Include: Architect decision, issues identified (if any)
  - Purpose: Audit trail for wave completion

- [ ] 7. If CHANGES_REQUIRED, record architectural issues
  - Source: Architect assessment report
  - Format: Similar to bugs but category = ARCHITECTURE
  - Location: May use bug-tracking.json or separate tracking
  - Purpose: Issues must be addressed before wave complete

### EXIT REQUIREMENTS (Must complete before transition)

- [ ] 8. Determine next state based on architect decision
  - If decision == PROCEED: Propose COMPLETE_WAVE
  - If decision == CHANGES_REQUIRED: Propose CREATE_WAVE_FIX_PLAN
  - Validation: Next state matches guard condition

- [ ] 9. Update state file to [next_state] per R288
  - Spawn: State Manager agent for SHUTDOWN_CONSULTATION
  - Provide: Work report with architect decision
  - Proposed next state: `COMPLETE_WAVE` or `CREATE_WAVE_FIX_PLAN`
  - State Manager validates transition and updates all 4 files atomically
  - Validation: State Manager returns validation_result with CONTINUE flag

- [ ] 10. Save TODOs per R287 (within 30s of last TodoWrite)
  - Trigger: "REVIEW_WAVE_ARCHITECTURE_COMPLETE"
  - Format: `todos/orchestrator-REVIEW_WAVE_ARCHITECTURE-[YYYYMMDD-HHMMSS].todo`
  - Validation: TODO file exists and contains current state

- [ ] 11. Commit all changes with descriptive message
  - Include: Architect decision
  - Include: Review outcome
  - Include: Rule compliance references (R288, R287, R510, R233)
  - Format: Multi-line commit message with context

- [ ] 12. Push changes to remote
  - Remote: `origin`
  - Branch: Current branch
  - Validation: `git status` shows "up to date with origin"

- [ ] 13. Set CONTINUE-SOFTWARE-FACTORY flag per R405
  - Value: `TRUE` (architect review complete, proceed based on decision)
  - Context: Architecture reviewed, transitioning to complete or fix plan
  - **NOTE**: This is NOT an R322 checkpoint - factory continues automatically

- [ ] 14. Stop execution (exit 0)
  - Command: `exit 0`
  - Timing: After ALL above items complete
  - Next: /continue-software-factory will proceed to next state based on guard

---

## State Purpose

REVIEW_WAVE_ARCHITECTURE performs architect review of the clean wave integration to validate pattern compliance and system coherence. This state spawns an Architect agent, waits for assessment, and determines whether the wave can be completed or if architectural changes are required.

**Primary Goal:** Validate wave integration architecture and design patterns
**Key Actions:** Spawn architect, assess integration, determine completion readiness
**Success Outcome:** Architect decision obtained, path determined (complete or changes needed)

---

## Entry Criteria

- **From**: REVIEW_WAVE_INTEGRATION (when bugs_found == 0)
- **Condition**: Integration is bug-free from code review
- **Required**:
  - bugs_found == 0 in most recent iteration
  - Build and tests passing on integration branch
  - Code review clean
  - Ready for architecture assessment

---

## State Actions

### 1. Verify Integration is Bug-Free

Confirm that the integration has zero bugs before architecture review.

**Implementation:**
```bash
CONTAINER_ID="wave-${PHASE_ID}-${WAVE_ID}"
WAVE_BRANCH="wave-${PHASE_ID}-${WAVE_ID}-integration"

echo "🔍 Verifying integration is bug-free..."

BUGS_FOUND=$(echo "✅ State file updated to: $NEXT_STATE"
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

if ! git commit -m "state: REVIEW_WAVE_ARCHITECTURE → $NEXT_STATE - REVIEW_WAVE_ARCHITECTURE complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: REVIEW_WAVE_ARCHITECTURE"
    echo "Attempted transition from: REVIEW_WAVE_ARCHITECTURE"
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
save_todos "REVIEW_WAVE_ARCHITECTURE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - REVIEW_WAVE_ARCHITECTURE complete [R287]"; then
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

**SUPREME LAW - MUST BE LAST LINE OF OUTPUT**

After completing this state's checklist, you MUST output EXACTLY ONE of these lines as the ABSOLUTE LAST LINE:

```
CONTINUE-SOFTWARE-FACTORY=TRUE
```

OR

```
CONTINUE-SOFTWARE-FACTORY=FALSE
```

**When to use TRUE:**
- ✅ Architect assessment completed
- ✅ Decision recorded (PROCEED or CHANGES_REQUIRED)
- ✅ State transition validated by State Manager
- ✅ Ready to proceed (complete wave or fix architectural issues)

**When to use FALSE:**
- ❌ Architect review failed
- ❌ Critical errors during assessment
- ❌ Cannot determine next path
- ❌ Requires human intervention

**DEFAULT for this state: TRUE** (assessment typically completes, failures go to ERROR_RECOVERY)

**IMPORTANT:** This is NOT an R322 checkpoint state. Factory continues automatically after architecture review.

---

## Additional Context

### Guard Conditions

This state has TWO possible exit paths based on architect decision:

```
if decision == PROCEED:
    next_state = COMPLETE_WAVE
else if decision == CHANGES_REQUIRED:
    next_state = CREATE_WAVE_FIX_PLAN
```

State Manager enforces these guards during SHUTDOWN_CONSULTATION.

### Architecture Review vs Code Review

**Code Review (REVIEW_WAVE_INTEGRATION):**
- Identifies bugs and defects
- Functional correctness
- Test coverage
- Code quality

**Architecture Review (REVIEW_WAVE_ARCHITECTURE):**
- Pattern compliance
- System coherence
- Design consistency
- Architectural soundness

Both must pass for wave to be complete.

### Architect Decisions

**PROCEED:**
- Architecture patterns followed correctly
- System coherence maintained
- No architectural concerns
- Wave can be completed

**CHANGES_REQUIRED:**
- Pattern violations detected
- System coherence issues
- Design improvements needed
- Must address before wave complete

### Common Pitfalls

1. **Skipping architecture review**: Code clean ≠ architecture sound
2. **Confusing with code review**: Different focus, different agent
3. **Not handling CHANGES_REQUIRED**: Must create fix plan like bugs
4. **Wrong guard logic**: decision value determines path

---

**State created per R516 State Creation Protocol**
**Template version: 1.0**
**Last updated: 2025-10-08**
