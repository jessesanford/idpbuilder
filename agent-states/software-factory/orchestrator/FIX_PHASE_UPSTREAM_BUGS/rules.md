# FIX_PHASE_UPSTREAM_BUGS State Rules


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

---

## 📋 PRIMARY DIRECTIVES FOR FIX_PHASE_UPSTREAM_BUGS STATE

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

6. **🔴🔴🔴 R321** - Immediate Backport Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-protocol.md`
   - Criticality: SUPREME LAW
   - Summary: Fix bugs in upstream wave branches, never in integration branch

7. **🚨🚨🚨 R232** - Spawn Agent Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R232-spawn-agent-protocol.md`
   - Criticality: BLOCKING
   - Summary: Spawn SW Engineer agents to fix bugs in upstream branches

8. **🔴🔴🔴 R313** - Bug Tracking Requirements
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R313-bug-tracking-requirements.md`
   - Criticality: SUPREME LAW
   - Summary: Update bug status as fixes are applied

---

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. Load approved fix plan from bug-tracking.json
  - Source: bug-tracking.json bugs with fix_plan defined
  - Filter: Status = IDENTIFIED, fix_plan.strategy = R321_BACKPORT
  - Validation: Fix plan exists for all upstream bugs
  - **BLOCKING**: Cannot fix without approved plan

- [ ] 2. Group bugs by affected upstream branch
  - Purpose: Spawn one SW Engineer per affected branch
  - Grouping: bugs → wave branch mapping
  - Validation: Each upstream branch has bug assignment list
  - **BLOCKING**: Required for efficient fix execution

- [ ] 3. Spawn SW Engineer for each affected upstream branch
  - Agent: sw-engineer
  - Task: Fix assigned bugs in wave branch
  - Context: Bug list, wave branch, fix instructions
  - Validation: SW Engineer agent spawned, returns agent ID
  - **BLOCKING**: Bugs cannot be fixed without agents

- [ ] 4. Monitor SW Engineers until all fixes complete
  - Track: Each agent's progress
  - Timeout: Reasonable time limit per bug
  - Completion: All agents report fixes committed and pushed
  - Validation: All upstream branches have fixes applied
  - **BLOCKING**: Must wait for all fixes before re-integration

- [ ] 5. Update bug status to FIXED in bug-tracking.json
  - Status: IDENTIFIED → FIXED
  - Record: Fixed timestamp, fixing agent ID, commit hash
  - Validation: All bugs marked FIXED
  - **BLOCKING**: Bug tracking required for convergence metrics

### STANDARD EXECUTION TASKS (Required)

- [ ] 6. Update orchestrator-state-v3.json with fix progress
  - Record: Bugs fixed count
  - Record: Upstream branches updated
  - Record: Ready for re-integration

- [ ] 7. Verify fixes committed and pushed to upstream branches
  - Check: Each affected branch has new commits
  - Check: Commits reference bug IDs
  - Validation: Remote branches updated

- [ ] 8. Record re-integration trigger in integration-containers.json
  - Note: Fixes complete, re-integration required
  - Purpose: Audit trail for iteration cycle

### EXIT REQUIREMENTS (Must complete before transition)

- [ ] 9. Update state file to START_PHASE_ITERATION per R288
  - Spawn: State Manager agent for SHUTDOWN_CONSULTATION
  - Provide: Work report with fixes complete, re-integration needed
  - Proposed next state: `START_PHASE_ITERATION`
  - State Manager validates transition and updates all 4 files atomically
  - Validation: State Manager returns validation_result with CONTINUE flag

- [ ] 10. Save TODOs per R287 (within 30s of last TodoWrite)
  - Trigger: "FIX_PHASE_UPSTREAM_BUGS_COMPLETE"
  - Format: `todos/orchestrator-FIX_PHASE_UPSTREAM_BUGS-[YYYYMMDD-HHMMSS].todo`
  - Validation: TODO file exists and contains current state

- [ ] 11. Commit all changes with descriptive message
  - Include: Bugs fixed summary
  - Include: Upstream branches updated
  - Include: Rule compliance references (R288, R287, R510, R321, R313)
  - Format: Multi-line commit message with context

- [ ] 12. Push changes to remote
  - Remote: `origin`
  - Branch: Current branch
  - Validation: `git status` shows "up to date with origin"

- [ ] 13. Set CONTINUE-SOFTWARE-FACTORY flag per R405
  - Value: `TRUE` (fixes complete, proceed to re-integration)
  - Context: Upstream bugs fixed, ready for new integration iteration
  - **NOTE**: This is NOT an R322 checkpoint - factory continues automatically

- [ ] 14. Stop execution (exit 0)
  - Command: `exit 0`
  - Timing: After ALL above items complete
  - Next: /continue-software-factory will proceed to START_PHASE_ITERATION (re-iteration)

---

## State Purpose

FIX_PHASE_UPSTREAM_BUGS executes the approved fix plan by spawning SW Engineer agents to fix bugs in upstream wave branches following R321 immediate backport protocol. This state monitors fix progress, updates bug tracking, and triggers re-integration once all fixes are applied.

**Primary Goal:** Fix bugs in upstream wave branches per R321
**Key Actions:** Spawn SW Engineers, monitor fixes, update bug status, trigger re-integration
**Success Outcome:** All bugs fixed in upstream branches, ready for re-integration

---

## Entry Criteria

- **From**: CREATE_PHASE_FIX_PLAN
- **Condition**: Fix plan approved by user
- **Required**:
  - Fix plan exists in bug-tracking.json
  - Affected upstream branches identified
  - User approved fix plan (R322 checkpoint passed)
  - Bugs categorized as UPSTREAM with R321_BACKPORT strategy

---

## State Actions

### 1. Load and Group Fix Assignments

Load approved fix plan and group bugs by upstream branch.

**Implementation:**
```bash
CONTAINER_ID="phase-${PHASE_ID}-${PHASE_ID}"
PHASE_BRANCH="phase-${PHASE_ID}-${PHASE_ID}-integration"

echo "📋 Loading fix plan..."
BUGS_TO_FIX=$(echo "✅ State file updated to: $NEXT_STATE"
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

if ! git commit -m "state: FIX_PHASE_UPSTREAM_BUGS → $NEXT_STATE - FIX_PHASE_UPSTREAM_BUGS complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: FIX_PHASE_UPSTREAM_BUGS"
    echo "Attempted transition from: FIX_PHASE_UPSTREAM_BUGS"
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
save_todos "FIX_PHASE_UPSTREAM_BUGS_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - FIX_PHASE_UPSTREAM_BUGS complete [R287]"; then
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
- ✅ All upstream bugs fixed successfully
- ✅ Bug status updated to FIXED
- ✅ Upstream branches committed and pushed
- ✅ State transition validated by State Manager
- ✅ Ready to proceed to START_PHASE_ITERATION (re-integration)

**When to use FALSE:**
- ❌ Some fixes failed
- ❌ SW Engineers encountered critical errors
- ❌ Cannot proceed with re-integration
- ❌ Requires human intervention

**DEFAULT for this state: TRUE** (fixes typically succeed, failures go to ERROR_RECOVERY)

**IMPORTANT:** This is NOT an R322 checkpoint state. Factory continues automatically after fixes.

---

## Additional Context

### The Fix-Reintegrate Loop

This state completes the iteration cycle:

```
START_PHASE_ITERATION → INTEGRATE_PHASE_EFFORTS → REVIEW_PHASE_INTEGRATION
                                                          ↓
                                                    bugs_found > 0
                                                          ↓
                                                 CREATE_PHASE_FIX_PLAN
                                                          ↓
                                                  user approves
                                                          ↓
                                              FIX_PHASE_UPSTREAM_BUGS ←┐
                                                          ↓           │
                                                  fixes complete      │
                                                          ↓           │
                                              START_PHASE_ITERATION ──┘
                                                 (re-integration)
```

### R321 Immediate Backport

Critical: Bugs are ALWAYS fixed in upstream wave branches:
- Maintains wave branches as source of truth
- Prevents divergence with integration branch
- Ensures re-integration is clean merge
- Makes debugging easier (bug fixed at source)

### Convergence Pattern

Each iteration should have fewer bugs:
- Iteration 1: N bugs → Fix N bugs upstream
- Iteration 2: < N bugs (hopefully fewer!) → Fix remaining
- Iteration 3: 0 bugs (convergence!) → Architecture review

### Common Pitfalls

1. **Fixing in integration branch**: Violates R321, creates divergence
2. **Not waiting for all fixes**: Partial fixes = broken re-integration
3. **Missing bug status update**: Breaks convergence tracking
4. **Forgetting to push**: Upstream fixes must be on remote

---

**State created per R516 State Creation Protocol**
**Template version: 1.0**
**Last updated: 2025-10-08**
