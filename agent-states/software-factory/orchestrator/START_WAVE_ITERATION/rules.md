# START_WAVE_ITERATION State Rules

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

7. **🔴🔴🔴 R531** - Integration Iteration Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R531-integration-iteration-protocol.md`
   - Criticality: SUPREME LAW
   - Summary: Iteration counter management and re-integration procedures

8. **🔴🔴🔴 R532** - Backport Attempt Limits
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R532-backport-attempt-limits.md`
   - Criticality: SUPREME LAW
   - Summary: Backport attempt counter limits to prevent infinite loops within same iteration

9. **🔴🔴🔴 R615** - Progress-Based Iteration Limits
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R615-progress-based-iteration-limits.md`
   - Criticality: SUPREME LAW
   - Summary: Two-tiered iteration limits (5 no-progress, 10 some-progress) based on actual bug fixes

10. **🔴🔴🔴 R616** - Bug Lifecycle Tracking Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R616-bug-lifecycle-tracking.md`
   - Criticality: SUPREME LAW
   - Summary: Bug identification, lifecycle states (OPEN/CLOSED/REOPENED), and tracking requirements

---

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. Conditionally increment wave iteration counter based on previous state
  - Tool: `tools/iteration-manager.sh increment_iteration WAVE` (only if from SETUP or FIX states)
  - Tool: `tools/iteration-manager.sh get_iteration_count WAVE` (if retrying same attempt)
  - Validation: Increment ONLY when from SETUP_WAVE_INFRASTRUCTURE or FIX_WAVE_UPSTREAM_BUGS
  - Validation: Do NOT increment when from IMMEDIATE_BACKPORT_REQUIRED, MONITORING_BACKPORT_PROGRESS, or ERROR_RECOVERY
  - Proof: `echo "✅ CHECKLIST[1]: Wave iteration [incremented to N | remains at N]"`

- [ ] 1b. Reset or check backport attempt counter (R532)
  - Tool: `tools/backport-attempt-manager.sh reset_backport_attempts WAVE` (if new iteration)
  - Tool: `tools/backport-attempt-manager.sh check_max_backport_attempts WAVE` (if retry)
  - Validation: Reset when R531 counter incremented (new iteration)
  - Validation: Check limit when R531 counter stayed same (retry)
  - Escalate: To ERROR_RECOVERY if backport attempts exceeded (exit 532)
  - Proof: `echo "✅ CHECKLIST[1b]: Backport attempts [reset to 0 | check passed N/M]"`

- [ ] 1c. Analyze bug progress and apply two-tiered iteration limits (R615, R616)
  - Tool: `bash tools/bug-progress-analyzer.sh should_continue_or_escalate WAVE $CURRENT_ITERATION`
  - Purpose: Determine if making ACTUAL progress (bugs fixed) vs just churning
  - Decision Logic:
    - PURE_PROGRESS/DISCOVERY_PROGRESS: Reset stall counter to 0, continue
    - SLOW_PROGRESS: Keep stall counter, continue
    - STALL/REGRESSION: Increment stall counter
    - FLAPPING (reopened bugs): Increment stall counter
  - Two-Tiered Limits:
    - Tier 1 (No-Progress): stall_counter >= 5 → ERROR_RECOVERY
    - Tier 2 (Some-Progress): iteration >= 10 → ERROR_RECOVERY
  - Validation: Tool returns decision: "CONTINUE" or "ERROR_RECOVERY"
  - On ERROR_RECOVERY: Exit with reason (no-progress or some-progress limit)
  - **BLOCKING**: Must analyze progress before checking blind iteration limit
  - Proof: `echo "✅ CHECKLIST[1c]: Progress analysis: [CATEGORY] (stalls: N/5, iterations: M/10, decision: [CONTINUE|ERROR_RECOVERY])"`

- [ ] 2. Check max iterations not exceeded (legacy R531 check - superseded by R615 progress-aware limits)
  - Tool: `tools/iteration-manager.sh check_max_iterations WAVE`
  - Validation: Must return "WITHIN_LIMIT"
  - On EXCEEDED: Escalate to ERROR_RECOVERY state
  - Note: This is now a backup check; R615 progress analysis (item 1c) is the primary enforcement
  - Proof: `echo "✅ CHECKLIST[2]: Max iterations check passed (N/10)"`

- [ ] 3. Determine next state
  - Success path: INTEGRATE_WAVE_EFFORTS
  - Failure path: ERROR_RECOVERY (if max iterations exceeded)
  - Proof: `echo "✅ CHECKLIST[3]: Next state determined: $NEXT_STATE"`

---

### ✅ Step 1: Conditionally Increment Wave Iteration Counter (R531/R336)
```bash
# Determine if iteration counter should increment based on previous state
# INCREMENT when:
#   - Coming from SETUP_WAVE_INFRASTRUCTURE (first iteration: 0->1)
#   - Coming from FIX_WAVE_UPSTREAM_BUGS (new attempt after full cycle: N->N+1)
# DO NOT INCREMENT when:
#   - Coming from IMMEDIATE_BACKPORT_REQUIRED (retrying same attempt)
#   - Coming from MONITORING_BACKPORT_PROGRESS (retrying same attempt)
#   - Coming from ERROR_RECOVERY (resuming same attempt)

PREVIOUS_STATE=$(jq -r '.state_machine.previous_state // "SETUP_WAVE_INFRASTRUCTURE"' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")

case "$PREVIOUS_STATE" in
    SETUP_WAVE_INFRASTRUCTURE|FIX_WAVE_UPSTREAM_BUGS)
        # These states indicate a NEW integration attempt - increment counter
        echo "📊 Previous state: $PREVIOUS_STATE - incrementing iteration counter (new attempt)"
        NEW_ITERATION=$(bash "$CLAUDE_PROJECT_DIR/tools/iteration-manager.sh" increment_iteration WAVE)

        if [ $? -ne 0 ]; then
            echo "❌ Failed to increment wave iteration counter"
            PROPOSED_NEXT_STATE="ERROR_RECOVERY"
            echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
            exit 1
        fi

        echo "✅ CHECKLIST[1]: Wave iteration incremented to ${NEW_ITERATION}"
        ;;

    IMMEDIATE_BACKPORT_REQUIRED|MONITORING_BACKPORT_PROGRESS|ERROR_RECOVERY)
        # These states indicate RETRYING same attempt - do NOT increment counter
        echo "📊 Previous state: $PREVIOUS_STATE - retrying same iteration (no increment)"
        CURRENT_ITERATION=$(bash "$CLAUDE_PROJECT_DIR/tools/iteration-manager.sh" get_iteration_count WAVE)
        echo "✅ CHECKLIST[1]: Wave iteration remains at ${CURRENT_ITERATION} (retry of same attempt)"
        ;;

    *)
        # Unknown previous state - log warning and increment to be safe
        echo "⚠️ WARNING: Unknown previous state '$PREVIOUS_STATE' - defaulting to increment"
        NEW_ITERATION=$(bash "$CLAUDE_PROJECT_DIR/tools/iteration-manager.sh" increment_iteration WAVE)
        echo "✅ CHECKLIST[1]: Wave iteration incremented to ${NEW_ITERATION} (default behavior)"
        ;;
esac
```

---

### ✅ Step 1b: Reset or Check Backport Attempt Counter (R532)
```bash
# R532: Manage backport_attempts_this_iteration to prevent infinite loops
# This counter tracks retries WITHIN the same iteration

case "$PREVIOUS_STATE" in
    SETUP_WAVE_INFRASTRUCTURE|FIX_WAVE_UPSTREAM_BUGS)
        # New iteration starting (R531 counter incremented above)
        # R532: Reset backport counter to 0 for fresh iteration
        echo "🔄 R532: Resetting backport attempts counter (new iteration)"
        bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" reset_backport_attempts WAVE

        if [ $? -ne 0 ]; then
            echo "❌ Failed to reset backport attempts counter"
            PROPOSED_NEXT_STATE="ERROR_RECOVERY"
            echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
            exit 1
        fi

        echo "✅ CHECKLIST[1b]: Backport attempts counter reset to 0 (new iteration per R532)"
        ;;

    IMMEDIATE_BACKPORT_REQUIRED|MONITORING_BACKPORT_PROGRESS|ERROR_RECOVERY)
        # Same iteration retry (R531 counter did NOT increment)
        # R532: Check if backport attempts have exceeded limit
        echo "🔍 R532: Checking backport attempts limit (retrying same iteration)"

        BACKPORT_STATUS=$(bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" \
            check_max_backport_attempts WAVE)

        if [ "$BACKPORT_STATUS" = "EXCEEDED" ]; then
            BACKPORT_COUNT=$(bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" \
                get_backport_attempt_count WAVE)
            MAX_BACKPORT=$(jq -r '.project_progression.current_wave.max_backport_attempts_per_iteration // 3' \
                "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")
            CURRENT_ITERATION=$(bash "$CLAUDE_PROJECT_DIR/tools/iteration-manager.sh" \
                get_iteration_count WAVE)

            echo "❌ R532 VIOLATION: Max backport attempts exceeded"
            echo "   Current iteration: ${CURRENT_ITERATION}"
            echo "   Backport attempts this iteration: ${BACKPORT_COUNT}/${MAX_BACKPORT}"
            echo "   Previous state: ${PREVIOUS_STATE}"
            echo ""
            echo "This indicates:"
            echo "- Backport fixes are not effective"
            echo "- Same issues reoccurring despite fixes"
            echo "- Upstream branches have systemic problems"
            echo "- Possible architecture issues"
            echo ""
            echo "REQUIRED ACTION: Escalate to ERROR_RECOVERY"

            PROPOSED_NEXT_STATE="ERROR_RECOVERY"
            echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MAX_BACKPORT_ATTEMPTS_EXCEEDED"
            exit 532
        fi

        BACKPORT_COUNT=$(bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" \
            get_backport_attempt_count WAVE)
        MAX_BACKPORT=$(jq -r '.project_progression.current_wave.max_backport_attempts_per_iteration // 3' \
            "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")
        echo "✅ CHECKLIST[1b]: Backport attempts check passed (${BACKPORT_COUNT}/${MAX_BACKPORT} per R532)"
        ;;

    *)
        # Unknown previous state - reset to be safe
        echo "⚠️ WARNING: Unknown previous state '$PREVIOUS_STATE' - resetting backport counter to be safe"
        bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" reset_backport_attempts WAVE
        echo "✅ CHECKLIST[1b]: Backport attempts counter reset (unknown previous state)"
        ;;
esac
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
