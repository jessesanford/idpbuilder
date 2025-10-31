# CREATE_PROJECT_FIX_PLAN State Rules


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

---

## 📋 PRIMARY DIRECTIVES FOR CREATE_PROJECT_FIX_PLAN STATE

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

6. **🔴🔴🔴 R322** - Mandatory Checkpoints
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-checkpoints.md`
   - Criticality: SUPREME LAW
   - Summary: Present fix plan to user for approval before proceeding

7. **🔴🔴🔴 R313** - Bug Tracking Requirements
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R313-bug-tracking-requirements.md`
   - Criticality: SUPREME LAW
   - Summary: Analyze bugs and create fix assignments in bug-tracking.json

8. **🔴🔴🔴 R321** - Immediate Backport Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-protocol.md`
   - Criticality: SUPREME LAW
   - Summary: Fix bugs in upstream effort branches (immediate backport)

---

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. Load bugs from bug-tracking.json for current iteration
  - Source: bug-tracking.json bugs array
  - Filter: discovered_in = project integration branch, iteration = current
  - Validation: bugs_found count matches REVIEW_PROJECT_INTEGRATION count
  - **BLOCKING**: Cannot plan fixes without bug list

- [ ] 2. Analyze each bug for upstream vs integration-specific
  - Categorization: UPSTREAM (in effort code) vs INTEGRATE_WAVE_EFFORTS (interaction issue)
  - Upstream: Bug exists in individual effort branch
  - Integration: Bug only appears when phases combined
  - Validation: Each bug has category assigned
  - **BLOCKING**: Category determines fix strategy

- [ ] 3. Identify affected upstream branches for each bug
  - Analysis: Trace bug to specific effort branch(es)
  - Record: bug-tracking.json affected_branches array
  - Validation: Each upstream bug maps to effort branch(es)
  - **BLOCKING**: Required to know where to apply fixes

- [ ] 4. Create fix plan for each bug
  - Plan elements: bug_id, fix_strategy, affected_branches, priority
  - Strategy: UPSTREAM_FIX (R321 backport) or INTEGRATE_WAVE_EFFORTS_FIX
  - Priority: Based on severity (Critical > High > Medium > Low)
  - Validation: Each bug has fix_plan entry in bug-tracking.json
  - **BLOCKING**: Fix plan required for execution

- [ ] 5. Present fix plan to user for approval (R322 checkpoint)
  - Format: Clear summary of bugs and fix strategies
  - Include: Bug count, affected branches, estimated re-integration time
  - Await: User approval to proceed
  - Validation: User approves plan
  - **BLOCKING**: R322 requires user approval at checkpoints

### STANDARD EXECUTION TASKS (Required)

- [ ] 6. Update orchestrator-state-v3.json with fix plan details
  - Record: Fix plan created timestamp
  - Record: Bugs to fix count
  - Record: Affected branches list

- [ ] 7. Calculate re-integration estimate
  - Consider: Number of bugs, complexity, affected branches
  - Output: Estimated time to fix and re-integrate
  - Purpose: User visibility into iteration cycle

### EXIT REQUIREMENTS (Must complete before transition)

- [ ] 8. Update state file to FIX_PROJECT_UPSTREAM_BUGS per R288
  - Spawn: State Manager agent for SHUTDOWN_CONSULTATION
  - Provide: Work report with fix plan and user approval
  - Proposed next state: `FIX_PROJECT_UPSTREAM_BUGS`
  - State Manager validates transition and updates all 4 files atomically
  - Validation: State Manager returns validation_result with CONTINUE flag

- [ ] 9. Save TODOs per R287 (within 30s of last TodoWrite)
  - Trigger: "CREATE_PROJECT_FIX_PLAN_COMPLETE"
  - Format: `todos/orchestrator-CREATE_PROJECT_FIX_PLAN-[YYYYMMDD-HHMMSS].todo`
  - Validation: TODO file exists and contains current state

- [ ] 10. Commit all changes with descriptive message
  - Include: Fix plan summary
  - Include: Bugs to fix list
  - Include: Rule compliance references (R288, R287, R510, R322, R313, R321)
  - Format: Multi-line commit message with context

- [ ] 11. Push changes to remote
  - Remote: `origin`
  - Branch: Current branch
  - Validation: `git status` shows "up to date with origin"

- [ ] 12. Set CONTINUE-SOFTWARE-FACTORY flag per R405
  - Value: `TRUE` if user approved, `FALSE` if user rejected/deferred
  - Context: Fix plan created and approved, ready to fix bugs
  - **NOTE**: This IS an R322 checkpoint - flag depends on user approval

- [ ] 13. Display R322 checkpoint message
  - Message: "🚧 R322 CHECKPOINT: Project fix plan approval required"
  - Content: Fix plan details, bugs summary, next steps
  - Action: Wait for user input

- [ ] 14. Stop execution (exit 0)
  - Command: `exit 0`
  - Timing: After ALL above items complete and user approves
  - Next: /continue-software-factory will proceed to FIX_PROJECT_UPSTREAM_BUGS if approved

---

## State Purpose

CREATE_PROJECT_FIX_PLAN analyzes bugs found during project integration review, categorizes them as upstream or integration-specific, identifies affected effort branches, and creates a fix plan. This state presents the plan to the user for approval per R322 checkpoint requirements before proceeding to fix execution.

**Primary Goal:** Create and gain approval for project-level bug fix plan
**Key Actions:** Analyze bugs, identify upstream branches, create fix strategies, get user approval
**Success Outcome:** Fix plan approved, ready to fix bugs in upstream branches

---

## Entry Criteria

- **From**: REVIEW_PROJECT_INTEGRATION (when bugs_found > 0)
- **Condition**: Bugs identified and recorded in bug-tracking.json
- **Required**:
  - bugs_found > 0 from code review
  - All bugs recorded in bug-tracking.json
  - Bug categories and severities assigned
  - Iteration history updated with bugs_found

---

## State Actions

### 1. Load and Analyze Bugs

Load bugs from bug-tracking.json and analyze their nature.

**Implementation:**
```bash
CONTAINER_ID="project-${PHASE_ID}-${PROJECT_ID}"
PROJECT_BRANCH="project-${PHASE_ID}-${PROJECT_ID}-integration"
CURRENT_ITER=$(echo "✅ State file updated to: $NEXT_STATE"
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

if ! git commit -m "state: CREATE_PROJECT_FIX_PLAN → $NEXT_STATE - CREATE_PROJECT_FIX_PLAN complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: CREATE_PROJECT_FIX_PLAN"
    echo "Attempted transition from: CREATE_PROJECT_FIX_PLAN"
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
save_todos "CREATE_PROJECT_FIX_PLAN_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - CREATE_PROJECT_FIX_PLAN complete [R287]"; then
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
- ✅ Fix plan created successfully
- ✅ All bugs analyzed and categorized
- ✅ Affected branches identified
- ✅ **USER APPROVED FIX PLAN** (R322 requirement)
- ✅ Ready to proceed to FIX_PROJECT_UPSTREAM_BUGS

**When to use FALSE:**
- ❌ User rejected fix plan
- ❌ User deferred decision
- ❌ Critical errors in analysis
- ❌ Requires human intervention before proceeding

**DEFAULT for this state: Depends on user approval** (R322 checkpoint)

**IMPORTANT:** This IS an R322 checkpoint state. CONTINUE flag depends on user approval.

---

## Additional Context

### R321 Immediate Backport Philosophy

Per R321, bugs are ALWAYS fixed in upstream effort branches, never in the integration branch:
- Ensures effort branches remain authoritative
- Prevents divergence between effort and integration
- Makes re-integration clean

### Fix Planning Strategy

**Upstream bugs:**
1. Identify effort branch(es) containing bug
2. Spawn SW Engineer to fix in effort branch
3. Re-integrate project (new iteration)

**Integration bugs:**
1. Rare - usually indicates missing effort
2. May require new effort or refactoring
3. Requires different strategy

### Convergence Tracking

Each iteration should reduce bugs:
- Iteration 1: N bugs → Create fix plan
- Fix bugs in upstream
- Iteration 2: < N bugs → Convergence!

### Common Pitfalls

1. **Trying to fix in integration branch**: Violates R321, creates divergence
2. **Missing affected branches**: Cannot fix without knowing where
3. **Skipping user approval**: R322 violation, checkpoint required
4. **Wrong categorization**: Determines entire fix strategy

---

**State created per R516 State Creation Protocol**
**Template version: 1.0**
**Last updated: 2025-10-08**
