# COMPLETE_PROJECT State Rules


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

---

## 📋 PRIMARY DIRECTIVES FOR COMPLETE_PROJECT STATE

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

6. **🔴🔴🔴 R336** - Mandatory Wave Integration Before Next Wave
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R336-mandatory-wave-integration-before-next-wave.md`
   - Criticality: SUPREME LAW
   - Summary: Wave integration creates iteration containers that must converge before next wave

7. **🔴🔴🔴 R234** - Mandatory State Traversal
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal.md`
   - Criticality: SUPREME LAW
   - Summary: Determine correct next state (next project or phase infrastructure)

---

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. Verify project integration fully approved
  - Check: bugs_found == 0 from code review
  - Check: Architect decision == PROCEED
  - Check: All tests passing
  - Validation: Project meets all completion criteria
  - **BLOCKING**: Cannot complete project without full approval

- [ ] 2. Mark project as completed in orchestrator-state-v3.json
  - Field: projects[project-${phase_id}-${project_id}].status = "COMPLETED"
  - Record: Completion timestamp
  - Record: Final iteration count
  - Validation: Project status updated
  - **BLOCKING**: State tracking required

- [ ] 3. Update integration-containers.json with final status
  - Field: Container status = "COMPLETED"
  - Field: completed_at = current timestamp
  - Record: Final convergence metrics
  - Validation: Container marked complete
  - **BLOCKING**: Iteration container must be closed

- [ ] 4. Calculate and record convergence metrics
  - Metrics: Total iterations, bugs per iteration, convergence rate
  - Purpose: Measure integration efficiency
  - Storage: integration-containers.json convergence_metrics
  - Validation: Metrics calculated and recorded
  - **BLOCKING**: Required for project analytics

- [ ] 5. Close iteration container
  - Action: Move from active_integrations to completed_integrations
  - Record: Container history preserved
  - Validation: Container no longer in active list
  - **BLOCKING**: Container cleanup required

### STANDARD EXECUTION TASKS (Required)

- [ ] 6. Determine if more projects in current phase
  - Check: orchestrator-state-v3.json phase definition
  - Check: Current project number vs total projects in phase
  - Result: more_projects_in_phase = true/false
  - Purpose: Determines next state

- [ ] 7. Generate project completion summary
  - Include: Iteration count, bugs found/fixed, effort branches
  - Include: Convergence metrics
  - Purpose: Project documentation and analytics

### EXIT REQUIREMENTS (Must complete before transition)

- [ ] 8. Determine next state based on phase progress
  - If more_projects_in_phase: Propose SETUP_PROJECT_INFRASTRUCTURE (next project)
  - If all_projects_complete_in_phase: Propose SETUP_PHASE_INFRASTRUCTURE
  - Validation: Next state matches guard condition

- [ ] 9. Update state file to [next_state] per R288
  - Spawn: State Manager agent for SHUTDOWN_CONSULTATION
  - Provide: Work report with project completion details
  - Proposed next state: `SETUP_PROJECT_INFRASTRUCTURE` or `SETUP_PHASE_INFRASTRUCTURE`
  - State Manager validates transition and updates all 4 files atomically
  - Validation: State Manager returns validation_result with CONTINUE flag

- [ ] 10. Save TODOs per R287 (within 30s of last TodoWrite)
  - Trigger: "COMPLETE_PROJECT_COMPLETE"
  - Format: `todos/orchestrator-COMPLETE_PROJECT-[YYYYMMDD-HHMMSS].todo`
  - Validation: TODO file exists and contains current state

- [ ] 11. Commit all changes with descriptive message
  - Include: Project completion summary
  - Include: Convergence metrics
  - Include: Rule compliance references (R288, R287, R510, R336, R234)
  - Format: Multi-line commit message with context

- [ ] 12. Push changes to remote
  - Remote: `origin`
  - Branch: Current branch
  - Validation: `git status` shows "up to date with origin"

- [ ] 13. Set CONTINUE-SOFTWARE-FACTORY flag per R405
  - Value: `TRUE` (project complete, proceed to next project or phase integration)
  - Context: Project successfully integrated, moving forward
  - **NOTE**: This is NOT an R322 checkpoint - factory continues automatically

- [ ] 14. Stop execution (exit 0)
  - Command: `exit 0`
  - Timing: After ALL above items complete
  - Next: /continue-software-factory will proceed to next state based on guard

---

## State Purpose

COMPLETE_PROJECT marks the current project as successfully integrated, closes the iteration container, records convergence metrics, and determines the next step in the project workflow (either starting the next project or proceeding to phase-level integration).

**Primary Goal:** Mark project as complete and determine next workflow step
**Key Actions:** Update project status, close container, record metrics, determine next state
**Success Outcome:** Project completed, container closed, ready for next project or phase integration

---

## Entry Criteria

- **From**: REVIEW_PROJECT_ARCHITECTURE (when decision == PROCEED)
- **Condition**: Project integration approved by both code review and architect
- **Required**:
  - bugs_found == 0 from code review
  - Architect decision = PROCEED
  - All tests passing on integration branch
  - Integration clean and ready

---

## State Actions

### 1. Verify Project Completion Criteria

Confirm all requirements met for project completion.

**Implementation:**
```bash
CONTAINER_ID="project-${PHASE_ID}-${PROJECT_ID}"
PROJECT_BRANCH="project-${PHASE_ID}-${PROJECT_ID}-integration"

echo "🔍 Verifying project completion criteria..."

# Check bugs_found == 0
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

if ! git commit -m "state: COMPLETE_PROJECT → $NEXT_STATE - COMPLETE_PROJECT complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: COMPLETE_PROJECT"
    echo "Attempted transition from: COMPLETE_PROJECT"
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
save_todos "COMPLETE_PROJECT_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - COMPLETE_PROJECT complete [R287]"; then
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
- ✅ Project marked as completed
- ✅ Container closed and archived
- ✅ Convergence metrics recorded
- ✅ Next state determined
- ✅ State transition validated by State Manager
- ✅ Ready to proceed (next project or phase integration)

**When to use FALSE:**
- ❌ Project completion failed
- ❌ Container closure errors
- ❌ Critical state corruption
- ❌ Requires human intervention

**DEFAULT for this state: TRUE** (completion typically succeeds)

**IMPORTANT:** This is NOT an R322 checkpoint state. Factory continues automatically after project completion.

---

## Additional Context

### Guard Conditions

This state has TWO possible exit paths based on phase progress:

```
if projects_completed < total_projects:
    next_state = SETUP_PROJECT_INFRASTRUCTURE  # Next project
else:
    next_state = SETUP_PHASE_INFRASTRUCTURE  # Phase integration
```

State Manager enforces these guards during SHUTDOWN_CONSULTATION.

### Project Completion Significance

Completing a project means:
- All phases in project successfully integrated
- Code review passed (bugs_found = 0)
- Architecture review passed (PROCEED)
- Tests passing
- Project integration branch ready

This project's integration branch will be input to phase-level integration.

### Convergence Metrics Purpose

Metrics track integration iteration efficiency:
- **Total iterations**: How many attempts to converge
- **Total bugs**: How many bugs found across all iterations
- **Convergence rate**: How quickly bugs decreased

Good metrics:
- Low iteration count (1-3 iterations typical)
- High convergence rate (bugs decrease each iteration)
- Zero bugs in final iteration

### Common Pitfalls

1. **Wrong next state**: Must check projects_completed vs total_projects
2. **Not closing container**: Container must move to completed
3. **Missing metrics**: Convergence data required for analytics
4. **Forgetting to increment projects_completed**: Breaks phase tracking

---

**State created per R516 State Creation Protocol**
**Template version: 1.0**
**Last updated: 2025-10-08**
