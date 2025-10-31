# SETUP_PHASE_INFRASTRUCTURE State Rules


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

---

## 📋 PRIMARY DIRECTIVES FOR SETUP_PHASE_INFRASTRUCTURE STATE

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

6. **🚨🚨🚨 R514** - Infrastructure Creation Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R514-infrastructure-creation-protocol.md`
   - Criticality: BLOCKING
   - Summary: Create phase integration branches following cascade pattern

7. **🔴🔴🔴 R336** - Mandatory Wave Integration Before Next Wave
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R336-mandatory-wave-integration-before-next-wave.md`
   - Criticality: SUPREME LAW
   - Summary: Wave integration creates iteration containers that must converge before next wave

---

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. Verify all waves in phase completed
  - Check: All wave integration branches merged and passing tests
  - Check: All wave architecture reviews passed
  - Validation: Query orchestrator-state-v3.json for phase completion status
  - **BLOCKING**: Cannot setup phase integration without completed waves

- [ ] 2. Create phase integration branch from main
  - Branch name: `phase-${phase_id}-integration`
  - Base: `origin/main` (latest)
  - Validation: `git branch --show-current` returns phase integration branch name
  - **BLOCKING**: Phase integration requires dedicated branch

- [ ] 3. Validate wave directory structure (R507 - Directory Validation)
  - Check: All wave integration directories exist and are accessible
  - Check: Each wave integration branch exists and is valid
  - Validation: Verify wave integration branches in Git and infrastructure state
  - **BLOCKING**: Cannot proceed with invalid or missing wave infrastructure

- [ ] 4. Initialize iteration counter in integration-containers.json
  - Container ID: `phase-${phase_id}`
  - iteration: 0 (will be incremented in START_PHASE_ITERATION)
  - max_iterations: 10 (default)
  - Validation: `echo "✅ State file updated to: $NEXT_STATE"
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

if ! git commit -m "state: SETUP_PHASE_INFRASTRUCTURE → $NEXT_STATE - SETUP_PHASE_INFRASTRUCTURE complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: SETUP_PHASE_INFRASTRUCTURE"
    echo "Attempted transition from: SETUP_PHASE_INFRASTRUCTURE"
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
save_todos "SETUP_PHASE_INFRASTRUCTURE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - SETUP_PHASE_INFRASTRUCTURE complete [R287]"; then
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
- ✅ All BLOCKING checklist items completed successfully
- ✅ Phase infrastructure setup complete
- ✅ Iteration container initialized
- ✅ State transition validated by State Manager
- ✅ Ready to proceed to START_PHASE_ITERATION

**When to use FALSE:**
- ❌ Infrastructure setup failed
- ❌ Critical error encountered
- ❌ Cannot proceed with phase integration
- ❌ Requires human intervention

**DEFAULT for this state: TRUE** (infrastructure setup is automated)

**IMPORTANT:** This is NOT an R322 checkpoint state. Factory continues automatically after infrastructure setup.

---

## Additional Context

### Phase Integration Container Pattern

Phase-level integration is the SECOND level of integration iteration containers in SF 3.0. The pattern is:
1. SETUP_PHASE_INFRASTRUCTURE (this state) - Create infrastructure
2. START_PHASE_ITERATION - Increment iteration counter
3. INTEGRATE_PHASE_WAVES - Perform integration
4. REVIEW_PHASE_INTEGRATION - Code review
5. CREATE_PHASE_FIX_PLAN / REVIEW_PHASE_ARCHITECTURE - Handle results
6. FIX_PHASE_UPSTREAM_BUGS - Fix and re-integrate (back to step 2)
7. COMPLETE_PHASE - Mark phase complete

### Iteration Container Philosophy

Per SF 3.0 Architecture Part 4:
- Integration is EXPECTED to require multiple iterations
- Fix-reintegrate cycles are NORMAL, not exceptional
- Convergence tracking ensures forward progress
- Maximum iterations prevent infinite loops

### Common Pitfalls

1. **Creating branch from wrong base**: Always use `origin/main`, never wave branches
2. **Forgetting to push branch**: Remote tracking required for agent coordination
3. **Initializing iteration > 0**: First iteration hasn't started yet, must be 0
4. **Missing convergence structures**: Required for tracking progress

---

**State created per R516 State Creation Protocol**
**Template version: 1.0**
**Last updated: 2025-10-08**
