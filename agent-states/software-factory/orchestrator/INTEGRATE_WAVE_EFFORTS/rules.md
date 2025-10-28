# INTEGRATE_WAVE_EFFORTS State Rules


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

---

## 📋 PRIMARY DIRECTIVES FOR INTEGRATE_WAVE_EFFORTS STATE

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

6. **🔴🔴🔴 R307** - Integration Iteration Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R307-integration-iteration-protocol.md`
   - Criticality: SUPREME LAW
   - Summary: Sequential merge patterns and conflict resolution procedures

7. **🔴🔴🔴 R308** - Incremental Branching Strategy
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R308-incremental-branching-strategy.md`
   - Criticality: SUPREME LAW
   - Summary: Merge effort branches sequentially following progressive integration

---

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. Verify all effort branches exist and are ready
  - Check: All effort branches in wave completed
  - Source: integration-containers.json effort_branches array
  - Validation: Each branch exists in git and has passing tests
  - **BLOCKING**: Cannot integrate missing or broken efforts

- [ ] 2. Ensure wave integration branch is clean
  - Branch: wave-${phase_id}-${wave_id}-integration
  - State: Clean working directory, no uncommitted changes
  - Base: Reset to origin/main (if re-iteration)
  - Validation: `git status --porcelain` returns empty
  - **BLOCKING**: Dirty state would corrupt integration

- [ ] 3. Merge effort branches sequentially into integration branch
  - Order: Sequential (effort-1, effort-2, effort-3, ...)
  - Method: `git merge --no-ff origin/[effort-branch]`
  - Conflict handling: Auto-resolve if possible, record conflicts
  - Validation: Each merge commits successfully
  - **BLOCKING**: Integration requires all efforts merged

- [ ] 4. Resolve merge conflicts if any
  - Detection: Check for conflict markers
  - Action: Apply conflict resolution strategy
  - Validation: All conflicts resolved, no markers remain
  - **BLOCKING**: Unresolved conflicts break build

- [ ] 5. Execute comprehensive integration testing checkpoint (EXPLICIT - R265)
  - **Build Verification**:
    * Command: Project-specific build command (e.g., `make build`, `go build ./...`, `npm run build`)
    * Capture: Build output and timing
    * Document: Build success/failure with error details
  - **Test Suite Execution**:
    * Command: Run ALL available tests (not just smoke tests)
    * Examples: `make test`, `go test ./...`, `npm test`, `pytest`
    * Capture: Test output, pass/fail counts, coverage %
    * Document: Failed tests with file:line locations
  - **Static Analysis** (if available):
    * Linting: `golangci-lint`, `eslint`, `pylint`, etc.
    * Type checking: `mypy`, `tsc --noEmit`, etc.
    * Document: Critical issues found
  - **Test Coverage Analysis** (if available):
    * Generate coverage report: `go test -cover`, `npm run coverage`, `pytest --cov`
    * Document: Coverage percentage and uncovered critical paths
  - **Integration Testing Report**:
    * Create: `INTEGRATE_WAVE_EFFORTS-TEST-RESULTS-wave-${phase_id}-${wave_id}.md`
    * Include: Build status, test counts, coverage %, failures list
    * Format: Per R265 integration testing documentation requirements
  - **Validation**: Build succeeds AND all tests pass (or failures documented)
  - **BLOCKING**: Cannot proceed to review with broken build or undocumented test failures
  - **Note**: This is an EXPLICIT testing checkpoint, not implicit - preserves SF 2.0 INTEGRATE_WAVE_EFFORTS_TESTING state functionality

### STANDARD EXECUTION TASKS (Required)

- [ ] 6. Record integration outcome in integration-containers.json
  - Field: iteration_history[N].integration_status
  - Values: "PROJECT_DONE" or "FAILED"
  - Include: Merge details, conflicts resolved, build status

- [ ] 7. Update orchestrator-state-v3.json with integration progress
  - Record: Integration complete timestamp
  - Record: Branches merged count
  - Record: Ready for code review

- [ ] 8. Push integrated wave branch to remote
  - Remote: origin
  - Branch: wave-${phase_id}-${wave_id}-integration
  - Validation: Remote branch updated with all merges

### EXIT REQUIREMENTS (Must complete before transition)

- [ ] 9. Update state file to REVIEW_WAVE_INTEGRATION per R288
  - Spawn: State Manager agent for SHUTDOWN_CONSULTATION
  - Provide: Work report with integration outcome
  - Proposed next state: `REVIEW_WAVE_INTEGRATION`
  - State Manager validates transition and updates all 4 files atomically
  - Validation: State Manager returns validation_result with CONTINUE flag

- [ ] 10. Save TODOs per R287 (within 30s of last TodoWrite)
  - Trigger: "INTEGRATE_WAVE_EFFORTS_COMPLETE"
  - Format: `todos/orchestrator-INTEGRATE_WAVE_EFFORTS-[YYYYMMDD-HHMMSS].todo`
  - Validation: TODO file exists and contains current state

- [ ] 11. Commit all changes with descriptive message
  - Include: Integration outcome summary
  - Include: Efforts merged list
  - Include: Rule compliance references (R288, R287, R510, R307, R308)
  - Format: Multi-line commit message with context

- [ ] 12. Push changes to remote
  - Remote: `origin`
  - Branch: Wave integration branch
  - Validation: `git status` shows "up to date with origin"

- [ ] 13. Set CONTINUE-SOFTWARE-FACTORY flag per R405
  - Value: `TRUE` (integration complete, proceed to review)
  - Context: Wave efforts integrated, ready for code review
  - **NOTE**: This is NOT an R322 checkpoint - factory continues automatically

- [ ] 14. Stop execution (exit 0)
  - Command: `exit 0`
  - Timing: After ALL above items complete
  - Next: /continue-software-factory will proceed to REVIEW_WAVE_INTEGRATION

---

## State Purpose

INTEGRATE_WAVE_EFFORTS merges all effort branches in the current wave into the wave integration branch. This state performs sequential merges following R308 incremental branching strategy, resolves conflicts, and validates that the integrated codebase builds successfully.

**Primary Goal:** Merge all wave efforts into single integration branch
**Key Actions:** Sequential merge, conflict resolution, build validation
**Success Outcome:** Wave integration branch contains all efforts, builds successfully

---

## Entry Criteria

- **From**: START_WAVE_ITERATION
- **Condition**: Iteration started, integration branch ready
- **Required**:
  - All effort branches completed and passing tests
  - Wave integration branch clean (reset if re-iteration)
  - Iteration counter incremented
  - No blocking upstream issues

---

## State Actions

### 1. Verify Effort Branch Readiness

Check that all effort branches are complete and ready for integration.

**Implementation:**
```bash
CONTAINER_ID="wave-${PHASE_ID}-${WAVE_ID}"
EFFORT_BRANCHES=$(echo "✅ State file updated to: $NEXT_STATE"
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

if ! git commit -m "state: INTEGRATE_WAVE_EFFORTS → $NEXT_STATE - INTEGRATE_WAVE_EFFORTS complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: INTEGRATE_WAVE_EFFORTS"
    echo "Attempted transition from: INTEGRATE_WAVE_EFFORTS"
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
save_todos "INTEGRATE_WAVE_EFFORTS_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - INTEGRATE_WAVE_EFFORTS complete [R287]"; then
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
- ✅ All effort branches merged successfully
- ✅ Conflicts resolved (if any)
- ✅ Build passes
- ✅ State transition validated by State Manager
- ✅ Ready to proceed to REVIEW_WAVE_INTEGRATION

**When to use FALSE:**
- ❌ Integration failed
- ❌ Build broken after merge
- ❌ Critical conflicts unresolvable
- ❌ Requires human intervention

**DEFAULT for this state: TRUE** (integration typically succeeds, failures go to ERROR_RECOVERY)

**IMPORTANT:** This is NOT an R322 checkpoint state. Factory continues automatically after integration.

---

## Additional Context

### Sequential Merge Pattern

Per R308, merges happen sequentially to:
- Detect conflicts incrementally
- Identify which effort causes issues
- Make debugging easier
- Preserve clean history

### Integration vs Re-Integration

**First Integration:**
- Fresh merges from effort branches
- Likely to succeed if efforts tested individually

**Re-Integration:**
- After upstream bugs fixed
- Integration branch was reset to clean state
- Should have fewer/no bugs than previous iteration

### Common Pitfalls

1. **Parallel merges**: Don't merge all at once, sequential is required
2. **Skipping build validation**: Always verify integrated code builds
3. **Not recording conflicts**: Document what was resolved for review
4. **Forgetting to push**: Remote must be updated for code reviewer access

---

**State created per R516 State Creation Protocol**
**Template version: 1.0**
**Last updated: 2025-10-08**
