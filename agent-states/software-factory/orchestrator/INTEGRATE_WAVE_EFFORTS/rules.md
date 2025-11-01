# INTEGRATE_WAVE_EFFORTS State Rules

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

6. **🚨🚨🚨 R329** - Orchestrator NEVER Performs Git Merges (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R329-orchestrator-never-performs-merges.md`
   - Criticality: BLOCKING - Immediate termination, 0% grade
   - Summary: ALL merge operations MUST be delegated to Integration Agent (extends R006)

7. **🔴🔴🔴 R531** - Integration Iteration Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R531-integration-iteration-protocol.md`
   - Criticality: SUPREME LAW
   - Summary: Iteration counter management and convergence tracking

8. **🔴🔴🔴 R308** - Incremental Branching Strategy
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

- [ ] 3. Verify integration infrastructure is ready
  - Check: Wave integration branch exists and is clean
  - Check: All effort branches are pushed to remote and accessible
  - Check: Integration workspace directory exists
  - Validation: Infrastructure ready for integration agent
  - **BLOCKING**: Cannot spawn agent without proper infrastructure

- [ ] 4. Prepare integration instructions for agent
  - Document: List of effort branches to merge (in sequential order per R308)
  - Document: Merge strategy (--no-ff for non-linear history)
  - Document: Conflict resolution guidelines
  - Document: Testing requirements per R265
  - Output: Integration instructions file for agent
  - **BLOCKING**: Agent needs clear instructions

- [ ] 5. Spawn integration agent to perform integration work (R329 - ORCHESTRATOR NEVER MERGES)
  - **CRITICAL**: Per R006 and R329, orchestrator MUST NEVER perform git merges
  - Agent: integration-agent
  - State: EXECUTE_WAVE_INTEGRATION
  - Workspace: Wave integration workspace
  - Instructions: Pass integration instructions file
  - Task: Sequential merge of all effort branches, conflict resolution, build validation, comprehensive testing
  - **BLOCKING**: Integration requires agent execution (orchestrator cannot do merges)

- [ ] 6. Monitor integration agent completion
  - Check: Integration agent reports completion
  - Review: Integration report with merge details
  - Review: Build validation results
  - Review: Test execution results (per R265 requirements)
  - Review: Conflict resolution summary (if any)
  - Validation: All efforts merged successfully, builds pass, tests pass
  - **BLOCKING**: Cannot proceed without successful integration

### STANDARD EXECUTION TASKS (Required)

- [ ] 7. Record integration outcome in integration-containers.json
  - Field: iteration_history[N].integration_status
  - Values: "PROJECT_DONE" or "FAILED"
  - Include: Merge details, conflicts resolved, build status

- [ ] 8. Update orchestrator-state-v3.json with integration progress
  - Record: Integration complete timestamp
  - Record: Branches merged count
  - Record: Ready for code review

- [ ] 9. Push integrated wave branch to remote
  - Remote: origin
  - Branch: wave-${phase_id}-${wave_id}-integration
  - Validation: Remote branch updated with all merges

### EXIT REQUIREMENTS (Must complete before transition)

- [ ] 10. Update state file to REVIEW_WAVE_INTEGRATION per R288
  - Spawn: State Manager agent for SHUTDOWN_CONSULTATION
  - Provide: Work report with integration outcome
  - Proposed next state: `REVIEW_WAVE_INTEGRATION`
  - State Manager validates transition and updates all 4 files atomically
  - Validation: State Manager returns validation_result with CONTINUE flag

- [ ] 11. Save TODOs per R287 (within 30s of last TodoWrite)
  - Trigger: "INTEGRATE_WAVE_EFFORTS_COMPLETE"
  - Format: `todos/orchestrator-INTEGRATE_WAVE_EFFORTS-[YYYYMMDD-HHMMSS].todo`
  - Validation: TODO file exists and contains current state

- [ ] 12. Commit all changes with descriptive message
  - Include: Integration outcome summary
  - Include: Efforts merged list
  - Include: Rule compliance references (R288, R287, R510, R531, R308, R329)
  - Format: Multi-line commit message with context

- [ ] 13. Push changes to remote
  - Remote: `origin`
  - Branch: Wave integration branch
  - Validation: `git status` shows "up to date with origin"

- [ ] 14. Set CONTINUE-SOFTWARE-FACTORY flag per R405
  - Value: `TRUE` (integration complete, proceed to review)
  - Context: Wave efforts integrated, ready for code review
  - **NOTE**: This is NOT an R322 checkpoint - factory continues automatically

- [ ] 15. Stop execution (exit 0)
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
- Missing Steps 3-6: Orchestrator performs merges directly = R006/R329 violation (IMMEDIATE FAILURE, 0% grade)
- Missing Step 10: No state update = state machine broken (R288 violation, -100%)
- Missing Step 11: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 12-13: No commit/push = state lost on compaction (R288 violation, -100%)
- Missing Step 14: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 15: No exit = R322 violation (-100%)

**ALL STEPS ARE MANDATORY - NO EXCEPTIONS**
**R329 CRITICAL: Orchestrator must NEVER perform git merges - always spawn integration agent**

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
