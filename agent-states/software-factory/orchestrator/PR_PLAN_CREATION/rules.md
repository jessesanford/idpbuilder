# PR_PLAN_CREATION State Rules

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

## State Overview
**State**: PR_PLAN_CREATION
**Type**: action
**Description**: Generating MASTER-PR-PLAN.md for human PRs (R279)

## Entry Conditions
- All validation complete
- Production-ready validation passed

## Exit Conditions
- PR plan created
- MASTER-PR-PLAN.md generated

## Required Actions
1. Generate comprehensive PR plan
2. Document all branches ready for PR
3. Create step-by-step instructions for human
4. Include validation results
5. Save MASTER-PR-PLAN.md

## Valid Transitions
- PROJECT_DONE: When PR plan is successfully created

## Critical Rules
- **R279**: PR plan generation protocol
- **R206**: State machine validation
- **R287**: TODO persistence requirements

## SF 3.0 State Management

This state reads from and updates orchestrator-state-v3.json:
- Reads final integration status from `state_machine.current_state` to confirm readiness for PR creation
- Accesses `integration-containers.json` to identify all successfully converged integration containers
- References `metadata_locations` in orchestrator-state-v3.json to locate validation reports and test results
- Updates `state_machine.current_state` to reflect PR plan creation completion per R288

The MASTER-PR-PLAN.md is created with complete information for human operators to execute pull requests to the main branch, as SF 3.0 never merges to main automatically (R280).

## Implementation Notes
This state is responsible for creating the final PR plan that a human will use to create actual pull requests. The plan should include:
- List of all branches ready for PR
- Order of PR submission
- PR titles and descriptions
- Any special instructions or dependencies
- Validation results summary

## Monitoring Requirements
- Verify MASTER-PR-PLAN.md is created
- Ensure all branches are documented
- Confirm validation results are included

## Error Handling
If PR plan generation fails:
1. Document the error
2. Transition to ERROR_RECOVERY
3. Preserve any partial work

## TODO Persistence
Must save TODOs before transitioning per R287.
## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete PR_PLAN_CREATION:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, set variables for State Manager
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="[REASON_FOR_TRANSITION]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for State Transition (R288 - SUPREME LAW)
```bash
# State Manager handles ALL state file updates atomically
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE_NAME" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"

# State Manager will:
# 1. Validate the transition against state machine
# 2. Update all 4 state tracking locations atomically:
#    - orchestrator-state-v3.json
#    - orchestrator-state-demo.json
#    - .cascade-state-backup.json
#    - .orchestrator-state-v3.json
# 3. Commit and push all changes
# 4. Return control

echo "✅ State Manager completed transition"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before exit (R287 trigger)
save_todos "STATE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - state complete [R287]"; then
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

### ✅ Step 5: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# State Manager determines the correct flag value based on next state
# PR_PLAN_CREATION → PROJECT_DONE (terminal), so State Manager sets FALSE

# Use the flag returned by State Manager's validation result
# For terminal state transitions, this will be FALSE
echo "CONTINUE-SOFTWARE-FACTORY=$CONTINUE_FLAG"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**
**NOTE**: State Manager sets CONTINUE_FLAG based on transition target (PROJECT_DONE is terminal → FALSE)

---

### ✅ Step 6: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No State Manager spawn = state machine broken (R288 violation, -100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)


