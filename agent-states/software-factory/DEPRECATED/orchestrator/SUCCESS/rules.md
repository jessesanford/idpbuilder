# PROJECT_DONE State Rules

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
## 🔴🔴🔴 CRITICAL: PROJECT_DONE MEANS TOTAL PROJECT COMPLETION! 🔴🔴🔴

PROJECT_DONE is the TERMINAL state indicating the ENTIRE Software Factory project is complete!

## 🚨🚨🚨 MANDATORY ENTRY CONDITIONS - ALL MUST BE TRUE 🚨🚨🚨

The PROJECT_DONE state can ONLY be entered when **ALL** of the following are true:

### 1. Phase Completion Requirements
- ✅ **ALL 5 PHASES COMPLETE** (or total_phases if different)
  - Phase 1: Complete and integrated
  - Phase 2: Complete and integrated
  - Phase 3: Complete and integrated
  - Phase 4: Complete and integrated
  - Phase 5: Complete and integrated
- ✅ current_phase == project_info.total_phases
- ✅ No phases skipped or incomplete

### 2. Wave Completion Requirements
- ✅ All waves in ALL phases complete
- ✅ All wave integration branches created and merged
- ✅ No pending waves in any phase

### 3. Effort Completion Requirements
- ✅ All efforts across ALL phases implemented
- ✅ All effort branches merged to wave branches
- ✅ All split efforts completed and integrated
- ✅ Zero efforts in "in_progress" state

### 4. Integration Requirements
- ✅ All wave integrations complete
- ✅ All phase integrations complete
- ✅ Project integration branch created
- ✅ Final project integration to main ready

### 5. Quality Requirements
- ✅ All code reviews passed
- ✅ All tests passing
- ✅ All build validations successful
- ✅ No outstanding issues

### 6. Documentation Requirements
- ✅ MASTER-PR-PLAN.md contains ALL project PRs
- ✅ All phase documentation complete
- ✅ Implementation report generated

## ❌ FORBIDDEN ENTRY CONDITIONS

**NEVER** enter PROJECT_DONE state if:
- ❌ Any phase remains (e.g., only Phase 1 of 5 complete)
- ❌ Any wave in any phase is incomplete
- ❌ Any effort is still in progress
- ❌ Integration steps are pending
- ❌ Tests are failing
- ❌ Code reviews are outstanding

## State Validation Logic

```python
def can_enter_success_state(state_file):
    """
    Validates if PROJECT_DONE state can be entered.
    Returns (bool, reason_string)
    """

    # Check 1: All phases complete
    current_phase = state_file.get('current_phase', 0)
    total_phases = state_file.get('project_info', {}).get('total_phases', 5)

    if current_phase != total_phases:
        return False, f"Only completed {current_phase}/{total_phases} phases"

    # Check 2: All waves complete
    for phase_num in range(1, total_phases + 1):
        phase_key = f"phase_{phase_num}"
        phase_data = state_file.get(phase_key, {})

        if not phase_data.get('completed', False):
            return False, f"Phase {phase_num} not marked complete"

        # Check all waves in phase
        waves = phase_data.get('waves', {})
        for wave_id, wave_data in waves.items():
            if wave_data.get('status') != 'COMPLETED':
                return False, f"Wave {wave_id} in Phase {phase_num} incomplete"

    # Check 3: No efforts in progress
    efforts_in_progress = state_file.get('efforts_in_progress', [])
    if efforts_in_progress:
        return False, f"{len(efforts_in_progress)} efforts still in progress"

    # Check 4: All integrations complete
    if not state_file.get('project_integration_complete', False):
        return False, "Project integration not complete"

    # All checks passed
    return True, "All project requirements met - ready for PROJECT_DONE"
```

## Actions in PROJECT_DONE State

1. **Generate Final Report**
   - Summary of all phases completed
   - Total efforts implemented
   - Integration branches created
   - PR organization structure

2. **Update State File**
   ```json
   {
     "current_state": "PROJECT_DONE",
     "project_complete": true,
     "completion_timestamp": "ISO-8601 timestamp",
     "final_statistics": {
       "phases_completed": 5,
       "waves_completed": "total",
       "efforts_completed": "total",
       "total_lines_added": "count"
     }
   }
   ```

3. **Create Completion Markers**
   - PROJECT-COMPLETE.md with summary
   - Final state backup
   - Completion notification

## Common Mistakes That Lead Here Incorrectly

### ❌ Mistake 1: Completing Phase 1 Only
- **Wrong**: Phase 1 done → PR_PLAN_CREATION → PROJECT_DONE
- **Right**: Phase 1 done → PR_PLAN_CREATION → START_PHASE_ITERATION (Phase 2)

### ❌ Mistake 2: Single Wave Completion
- **Wrong**: Wave complete → PROJECT_DONE
- **Right**: Wave complete → Next wave or phase

### ❌ Mistake 3: PR Plan Creation Trigger
- **Wrong**: Creating PR plan means project done
- **Right**: PR plans are created throughout project

## Exit Transitions

**PROJECT_DONE has NO exit transitions** - it is the terminal state.

Once in PROJECT_DONE:
- Project is complete
- No further automated work
- State machine terminates

## Automation Flag

```bash
# PROJECT_DONE is a terminal state - orchestration complete
echo "🎉 PROJECT IMPLEMENTATION COMPLETE!"
echo "All phases and waves successfully implemented."
echo "CONTINUE-SOFTWARE-FACTORY=COMPLETE"  # Special flag for terminal state
```

## Monitoring and Alerts

If PROJECT_DONE is reached prematurely:
- 🚨 CRITICAL ERROR
- Likely state machine misconfiguration
- Requires immediate investigation
- Check state transition history

## References
- R206: State Machine Validation
- R279: PR Plan Creation
- State Machine: /state-machines/software-factory-3.0-state-machine.json
- Project Completion Criteria: /.claude/agents/orchestrator.md
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
