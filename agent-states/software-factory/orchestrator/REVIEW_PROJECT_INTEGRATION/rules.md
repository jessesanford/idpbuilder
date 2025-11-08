# REVIEW_PROJECT_INTEGRATION State Rules

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

## 📋 PRIMARY DIRECTIVES FOR REVIEW_PROJECT_INTEGRATION STATE

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
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R232-enforcement-examples.md`
   - Criticality: BLOCKING
   - Summary: Spawn code reviewer agent with proper context and workspace

7. **🔴🔴🔴 R313** - Bug Tracking Requirements
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R313-mandatory-stop-after-spawn.md`
   - Criticality: SUPREME LAW
   - Summary: Record all bugs in bug-tracking.json with proper categorization

---

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. Verify integration build passes
  - Check: Project integration branch builds successfully
  - Check: Basic smoke tests pass
  - Validation: Build exit code = 0
  - **BLOCKING**: Cannot review broken code

- [ ] 2. Spawn Code Reviewer agent for integration review
  - Agent: code-reviewer
  - Context: Project integration branch (project-${phase_id}-${project_id}-integration)
  - Task: Identify bugs and integration issues
  - Validation: Code reviewer agent spawned, returns agent ID
  - **BLOCKING**: Review required to identify bugs

- [ ] 3. Wait for Code Reviewer to complete review
  - Monitor: Code reviewer progress
  - Timeout: Reasonable time limit (project-specific)
  - Output: Review report with bug findings
  - Validation: Review report exists with bugs_found count
  - **BLOCKING**: Cannot proceed without review results

- [ ] 4. Record bugs in bug-tracking.json
  - Source: Code reviewer report
  - Fields: bug_id, severity, category, affected_branch, description
  - Location: bug-tracking.json bugs array
  - Validation: All bugs from review recorded
  - **BLOCKING**: Bug tracking required for fix planning

- [ ] 5. Update iteration history with bugs_found count
  - Field: integration-containers.json iteration_history[N].bugs_found
  - Value: Count from bug-tracking.json
  - Validation: bugs_found = number of bugs recorded
  - **BLOCKING**: Iteration tracking requires bug count

### STANDARD EXECUTION TASKS (Required)

- [ ] 6. Update orchestrator-state-v3.json with review outcome
  - Record: Review completion timestamp
  - Record: bugs_found count
  - Record: Reviewer agent ID

- [ ] 7. Categorize bugs by type and severity
  - Integration-specific: Bugs from effort interaction
  - Upstream: Bugs in individual effort code
  - Severity: Critical, High, Medium, Low
  - Purpose: Informs fix planning strategy

### EXIT REQUIREMENTS (Must complete before transition)

- [ ] 8. Determine next state based on bugs_found
  - If bugs_found > 0: Propose CREATE_PROJECT_FIX_PLAN
  - If bugs_found == 0: Propose REVIEW_PROJECT_ARCHITECTURE
  - Validation: Next state matches guard conditions

- [ ] 9. Update state file to [next_state] per R288
  - Spawn: State Manager agent for SHUTDOWN_CONSULTATION
  - Provide: Work report with review results and bugs_found
  - Proposed next state: `CREATE_PROJECT_FIX_PLAN` or `REVIEW_PROJECT_ARCHITECTURE`
  - State Manager validates transition and updates all 4 files atomically
  - Validation: State Manager returns validation_result with CONTINUE flag

- [ ] 10. Save TODOs per R287 (within 30s of last TodoWrite)
  - Trigger: "REVIEW_PROJECT_INTEGRATION_COMPLETE"
  - Format: `todos/orchestrator-REVIEW_PROJECT_INTEGRATION-[YYYYMMDD-HHMMSS].todo`
  - Validation: TODO file exists and contains current state

- [ ] 11. Commit all changes with descriptive message
  - Include: Review outcome summary
  - Include: Bugs found count
  - Include: Rule compliance references (R288, R287, R510, R313)
  - Format: Multi-line commit message with context

- [ ] 12. Push changes to remote
  - Remote: `origin`
  - Branch: Current branch
  - Validation: `git status` shows "up to date with origin"

- [ ] 13. Set CONTINUE-SOFTWARE-FACTORY flag per R405
  - Value: `TRUE` (review complete, proceed based on bugs_found)
  - Context: Review done, transitioning to fix plan or architecture review
  - **NOTE**: This is NOT an R322 checkpoint - factory continues automatically

- [ ] 14. Stop execution (exit 0)
  - Command: `exit 0`
  - Timing: After ALL above items complete
  - Next: /continue-software-factory will proceed to next state based on guard

---

## State Purpose

REVIEW_PROJECT_INTEGRATION performs code review of the project integration branch to identify bugs and integration issues. This state spawns a Code Reviewer agent, waits for review completion, records all bugs in bug-tracking.json, and determines whether fix planning or architecture review is needed next.

**Primary Goal:** Identify bugs in project integration through code review
**Key Actions:** Spawn code reviewer, record bugs, determine next path
**Success Outcome:** Review complete, bugs recorded, path determined (fix or architect review)

---

## Entry Criteria

- **From**: INTEGRATE_PROJECT_EFFORTS
- **Condition**: Integration complete, build passes
- **Required**:
  - Project integration branch exists with all phases merged
  - Build successful on integration branch
  - Basic tests passing
  - Integration branch pushed to remote

---

## State Actions

### 1. Verify Integration Build Status

Confirm that the integrated code builds successfully before review.

**Implementation:**
```bash
PROJECT_BRANCH="project-${PHASE_ID}-${PROJECT_ID}-integration"
git checkout "$PROJECT_BRANCH"

echo "🏗️ Verifying build status before review..."
if ! make build && make test-basic; then
    echo "❌ Build broken, cannot review"
    # Transition to ERROR_RECOVERY
    exit 1
fi

echo "✅ Build passes, ready for code review"
```

**Validation:**
- Build succeeds
- Tests pass
- Ready for review

### 2. Spawn Code Reviewer Agent

Spawn code reviewer agent following R232 protocol.

**Implementation:**
```bash
echo "🔍 Spawning Code Reviewer for project integration review..."

# Prepare reviewer context
REVIEWER_CONTEXT="Project ${PHASE_ID}.${PROJECT_ID} integration review"
REVIEW_BRANCH="$PROJECT_BRANCH"
CONTAINER_ID="project-${PHASE_ID}-${PROJECT_ID}"

# Spawn code reviewer agent (orchestrator spawns, doesn't implement)
# Actual spawn mechanism project-specific
REVIEWER_AGENT_ID=$(spawn_agent "code-reviewer" \
    --context "$REVIEWER_CONTEXT" \
    --branch "$REVIEW_BRANCH" \
    --task "INTEGRATE_WAVE_EFFORTS_REVIEW")

echo "✅ Code Reviewer spawned: $REVIEWER_AGENT_ID"
```

**Validation:**
- Agent spawned successfully
- Agent ID returned
- Agent has access to integration branch

### 3. Wait for Review Completion

Monitor code reviewer progress and wait for review report.

**Implementation:**
```bash
echo "⏳ Waiting for Code Reviewer to complete review..."

# Poll for review completion (project-specific mechanism)
while ! review_complete "$REVIEWER_AGENT_ID"; do
    sleep 30
    echo "Still reviewing..."
done

# Get review report
REVIEW_REPORT=$(get_review_report "$REVIEWER_AGENT_ID")
BUGS_FOUND=$(echo "$REVIEW_REPORT" | jq -r '.bugs_found')

echo "✅ Review complete: $BUGS_FOUND bugs found"
```

**Validation:**
- Review report retrieved
- Bug count extracted
- Ready to record bugs

### 4. Record Bugs in bug-tracking.json

Record all bugs found by code reviewer following R313 requirements.

**Implementation:**
```bash
echo "📝 Recording $BUGS_FOUND bugs in bug-tracking.json..."

# Extract bugs from review report and add to bug-tracking.json
# (Implementation depends on bug tracking schema)
for bug in $(echo "$REVIEW_REPORT" | jq -c '.bugs[]'); do
    BUG_ID=$(echo "$bug" | jq -r '.bug_id')
    SEVERITY=$(echo "$bug" | jq -r '.severity')
    CATEGORY=$(echo "$bug" | jq -r '.category')
    DESCRIPTION=$(echo "$bug" | jq -r '.description')

    # Add bug to bug-tracking.json (R313)
    jq ".bugs += [{
        bug_id: \"$BUG_ID\",
        severity: \"$SEVERITY\",
        category: \"$CATEGORY\",
        affected_branch: \"$PROJECT_BRANCH\",
        description: \"$DESCRIPTION\",
        found_in_iteration: $CURRENT_ITERATION,
        status: \"OPEN\"
    }]" bug-tracking.json > bug-tracking.json.tmp
    mv bug-tracking.json.tmp bug-tracking.json
done

echo "✅ All bugs recorded in bug-tracking.json"
```

**Validation:**
- All bugs from review added
- Bug tracking file updated
- Ready for fix planning

### 5. Update Iteration History

Update iteration_history with bugs_found count for convergence tracking.

**Implementation:**
```bash
echo "📊 Updating iteration history with bugs_found count..."

# Update integration-containers.json with bugs_found
CONTAINER_ID="project-${PROJECT_ID}"
jq "(.containers[] | select(.container_id == \"$CONTAINER_ID\") | .iteration_history[-1].bugs_found) = $BUGS_FOUND" \
    integration-containers.json > integration-containers.json.tmp
mv integration-containers.json.tmp integration-containers.json

echo "✅ Iteration history updated: bugs_found = $BUGS_FOUND"
```

**Validation:**
- Iteration history updated
- Convergence trackable
- Ready for next state determination

---

## Exit Transition

### 1. Determine Next State

Based on bugs_found, determine the next state per guard conditions.

**Implementation:**
```bash
if [ "$BUGS_FOUND" -gt 0 ]; then
    PROPOSED_NEXT_STATE="SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING"
    echo "🔧 Bugs found, proceeding to fix planning"
else
    PROPOSED_NEXT_STATE="REVIEW_PROJECT_ARCHITECTURE"
    echo "✅ No bugs found, project complete!"
fi
```

---

### 2. Spawn State Manager for State Transition

Per R288 SF 3.0, spawn State Manager for SHUTDOWN_CONSULTATION.

**Implementation:**
```bash
echo "🔄 Spawning State Manager for state transition validation..."

# Prepare work report for State Manager
WORK_REPORT="REVIEW_PROJECT_INTEGRATION complete
- Review completed successfully
- Bugs found: $BUGS_FOUND
- All bugs recorded in bug-tracking.json
- Iteration history updated
- Proposed next state: $PROPOSED_NEXT_STATE"

# Spawn State Manager for SHUTDOWN_CONSULTATION (R288)
/spawn agent-state-manager SHUTDOWN_CONSULTATION \
    --current-state "REVIEW_PROJECT_INTEGRATION" \
    --proposed-next-state "$PROPOSED_NEXT_STATE" \
    --work-report "$WORK_REPORT"

echo "✅ State Manager spawned for transition validation"
```

**State Manager will:**
- Validate transition guard conditions
- Update orchestrator-state-v3.json
- Update integration-containers.json
- Update bug-tracking.json
- Update fix-cascade-state.json
- Return validation_result with CONTINUE flag

---

### ✅ Step 3: Wait for State Manager Validation

```bash
# Wait for State Manager to complete validation and return result
# State Manager will update all 4 state files atomically
echo "⏳ Waiting for State Manager validation..."

# Poll for validation result
while [ ! -f ".state-manager-result.json" ]; do
    sleep 2
done

VALIDATION_RESULT=$(cat .state-manager-result.json)
NEXT_STATE=$(echo "$VALIDATION_RESULT" | jq -r '.validated_next_state')
CONTINUE_FLAG=$(echo "$VALIDATION_RESULT" | jq -r '.continue_flag')

rm .state-manager-result.json

echo "✅ State file updated to: $NEXT_STATE"
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

if ! git commit -m "state: REVIEW_PROJECT_INTEGRATION → $NEXT_STATE - REVIEW_PROJECT_INTEGRATION complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: REVIEW_PROJECT_INTEGRATION"
    echo "Attempted transition from: REVIEW_PROJECT_INTEGRATION"
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
save_todos "REVIEW_PROJECT_INTEGRATION_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - REVIEW_PROJECT_INTEGRATION complete [R287]"; then
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
- ✅ Code review completed successfully
- ✅ All bugs recorded in bug-tracking.json
- ✅ Iteration history updated
- ✅ State transition validated by State Manager
- ✅ Ready to proceed (fix plan or architect review)

**When to use FALSE:**
- ❌ Code review failed
- ❌ Critical errors during review
- ❌ Cannot determine next path
- ❌ Requires human intervention

**DEFAULT for this state: TRUE** (review typically completes, failures go to ERROR_RECOVERY)

**IMPORTANT:** This is NOT an R322 checkpoint state. Factory continues automatically after review.

---

## Additional Context

### Guard Conditions

This state has TWO possible exit paths based on bugs_found:

```
if bugs_found > 0:
    next_state = CREATE_PROJECT_FIX_PLAN
else:
    next_state = REVIEW_PROJECT_ARCHITECTURE
```

State Manager enforces these guards during SHUTDOWN_CONSULTATION.

### Bug Categories

Bugs are typically categorized as:
- **Integration bugs**: Issues from effort interaction
- **Upstream bugs**: Issues in individual effort code
- **Build bugs**: Issues with build/compilation
- **Test bugs**: Test failures

### Convergence Indicator

bugs_found in each iteration should decrease over iterations:
- Iteration 1: N bugs
- Iteration 2: < N bugs (after fixes)
- Iteration 3: << N bugs (convergence)

### Common Pitfalls

1. **Not waiting for review completion**: Must poll until reviewer finishes
2. **Missing bugs in recording**: All bugs must be tracked
3. **Wrong guard logic**: bugs_found == 0 vs > 0 determines path
4. **Not updating iteration history**: Required for convergence tracking

---

**State created per R516 State Creation Protocol**
**Template version: 1.0**
**Last updated: 2025-10-08**
