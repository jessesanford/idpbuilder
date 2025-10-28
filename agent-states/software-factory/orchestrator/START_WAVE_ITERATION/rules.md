# START_WAVE_ITERATION State Rules


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

---

## 📋 PRIMARY DIRECTIVES FOR START_WAVE_ITERATION STATE

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

7. **🔴🔴🔴 R307** - Integration Iteration Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R307-integration-iteration-protocol.md`
   - Criticality: SUPREME LAW
   - Summary: Iteration counter management and re-integration procedures

---

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. Increment wave iteration counter
  - Tool: `tools/iteration-manager.sh increment_iteration WAVE`
  - Validation: Check return value is new iteration number
  - Proof: `echo "✅ CHECKLIST[1]: Wave iteration incremented to N"`

- [ ] 2. Check max iterations not exceeded
  - Tool: `tools/iteration-manager.sh check_max_iterations WAVE`
  - Validation: Must return "WITHIN_LIMIT"
  - On EXCEEDED: Escalate to ERROR_RECOVERY state
  - Proof: `echo "✅ CHECKLIST[2]: Max iterations check passed (N/10)"`

- [ ] 3. Determine next state
  - Success path: INTEGRATE_WAVE_EFFORTS
  - Failure path: ERROR_RECOVERY (if max iterations exceeded)
  - Proof: `echo "✅ CHECKLIST[3]: Next state determined: $NEXT_STATE"`

---

### ✅ Step 1: Increment Wave Iteration Counter (R307/R336)
```bash
# Increment wave iteration counter using iteration-manager tool
NEW_ITERATION=$(bash "$CLAUDE_PROJECT_DIR/tools/iteration-manager.sh" increment_iteration WAVE)

if [ $? -ne 0 ]; then
    echo "❌ Failed to increment wave iteration counter"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
    exit 1
fi

echo "✅ CHECKLIST[1]: Wave iteration incremented to ${NEW_ITERATION}"
```

---

### ✅ Step 2: Check Max Iterations (R336)
```bash
# Check if max iterations exceeded - if so, escalate to ERROR_RECOVERY
ITERATION_STATUS=$(bash "$CLAUDE_PROJECT_DIR/tools/iteration-manager.sh" check_max_iterations WAVE)

if [ "$ITERATION_STATUS" = "EXCEEDED" ]; then
    echo "❌ Max iterations exceeded for wave integration"
    echo "   Escalating to ERROR_RECOVERY per R336"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
    exit 1
fi

# Get current iteration for reporting
CURRENT_ITER=$(bash "$CLAUDE_PROJECT_DIR/tools/iteration-manager.sh" get_iteration_count WAVE)
MAX_ITER=$(jq -r '.project_progression.current_wave.max_iterations // 10' orchestrator-state-v3.json)

echo "✅ CHECKLIST[2]: Max iterations check passed (${CURRENT_ITER}/${MAX_ITER})"
```

---

### ✅ Step 3: Determine Next State
```bash
# Normal path: proceed to wave integration
PROPOSED_NEXT_STATE="INTEGRATE_WAVE_EFFORTS"

echo "✅ CHECKLIST[3]: Next state determined: ${PROPOSED_NEXT_STATE}"
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

if ! git commit -m "state: START_WAVE_ITERATION → $NEXT_STATE - START_WAVE_ITERATION complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: START_WAVE_ITERATION"
    echo "Attempted transition from: START_WAVE_ITERATION"
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
save_todos "START_WAVE_ITERATION_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - START_WAVE_ITERATION complete [R287]"; then
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
- ✅ Iteration counter incremented successfully
- ✅ Iteration limit not exceeded
- ✅ Integration branch ready
- ✅ State transition validated by State Manager
- ✅ Ready to proceed to INTEGRATE_WAVE_EFFORTS

**When to use FALSE:**
- ❌ Max iterations exceeded
- ❌ Critical error encountered
- ❌ Cannot proceed with integration
- ❌ Requires human intervention

**DEFAULT for this state: TRUE** (iteration start is automated unless max exceeded)

**IMPORTANT:** This is NOT an R322 checkpoint state. Factory continues automatically after iteration start.

---

## Additional Context

### First Iteration vs Re-Iteration

**First Iteration (iteration = 1):**
- Integration branch already clean from SETUP_WAVE_INFRASTRUCTURE
- No reset needed
- First attempt at integrating all efforts

**Re-Iterations (iteration > 1):**
- Previous integration had bugs
- Bugs fixed in upstream effort branches
- Integration branch MUST be reset to clean state
- Fresh integration attempt with fixed efforts

### Convergence Philosophy

Per SF 3.0 Architecture Part 4:
- Multiple iterations are EXPECTED and NORMAL
- Each iteration should find fewer bugs (convergence)
- Max iterations prevent infinite loops
- Iteration history tracks convergence progress

### Common Pitfalls

1. **Forgetting to reset branch on re-iteration**: Leads to conflicts and confusion
2. **Not checking max iterations**: Could create infinite loops
3. **Missing iteration history entry**: Breaks convergence tracking
4. **Incrementing wrong container**: Verify container_id carefully

---

**State created per R516 State Creation Protocol**
**Template version: 1.0**
**Last updated: 2025-10-08**
