# PR_PLAN_CREATION State Rules

## 🔴🔴🔴 CRITICAL: PR PLANS ARE NOT PROJECT COMPLETION! 🔴🔴🔴

PR plans are created THROUGHOUT the project, not just at the end!

## State Entry Conditions
- All validation complete for current work
- Build tests passing
- Code review completed
- Ready to document PR strategy

## Primary Responsibilities
1. Generate MASTER-PR-PLAN.md for human PRs
2. Document PR organization strategy
3. Determine next state based on project progress
4. NEVER assume project completion without verification

## 🚨🚨🚨 MANDATORY STATE TRANSITION LOGIC 🚨🚨🚨

### When to Transition to Each State

#### → WAVE_START
Transition to WAVE_START when:
- More waves exist in current phase
- Current wave is complete but phase continues
- Next wave infrastructure needs creation

#### → COMPLETE_PHASE
Transition to COMPLETE_PHASE when:
- All waves in current phase are done
- Phase integration is needed
- Current phase < total phases

#### → START_PHASE_ITERATION
Transition to START_PHASE_ITERATION when:
- Starting a new phase (Phase 2, 3, 4, or 5)
- Previous phase fully integrated
- More phases remain in project

#### → PROJECT_DONE (ONLY WHEN PROJECT IS 100% COMPLETE!)
Transition to PROJECT_DONE **ONLY** when ALL of these are true:
- ✅ current_phase == project_info.total_phases (typically 5)
- ✅ All waves in final phase are complete
- ✅ All efforts in all phases are complete
- ✅ All integration branches merged
- ✅ No pending work remains
- ✅ MASTER-PR-PLAN.md contains ALL project PRs

**🚨 NEVER transition to PROJECT_DONE if ANY phases remain! 🚨**

## Decision Logic Implementation

```python
# Pseudo-code for state transition decision
def determine_next_state(state_file):
    current_phase = state_file['current_phase']
    current_wave = state_file['current_wave']
    total_phases = state_file['project_info']['total_phases']

    # Get phase plan to check wave counts
    phase_plan = read_phase_plan(current_phase)
    total_waves_in_phase = len(phase_plan['waves'])

    # CHECK 1: Are we done with entire project?
    if current_phase == total_phases:
        if current_wave == total_waves_in_phase:
            if all_efforts_complete():
                return "PROJECT_DONE"  # ONLY valid path to PROJECT_DONE

    # CHECK 2: More waves in current phase?
    if current_wave < total_waves_in_phase:
        return "WAVE_START"

    # CHECK 3: Current phase complete, more phases remain?
    if current_wave == total_waves_in_phase:
        if current_phase < total_phases:
            return "COMPLETE_PHASE"

    # Default to ERROR_RECOVERY if logic unclear
    return "ERROR_RECOVERY"
```

## Required Actions Before State Exit

1. **Generate PR Plan**: Create or update MASTER-PR-PLAN.md
2. **Document Progress**: Record which phase/wave PRs are included
3. **Validate Next State**: Use decision logic above
4. **Update State File**: Record chosen transition
5. **Log Rationale**: Document WHY this transition was chosen

## Common Mistakes to Avoid

### ❌ FORBIDDEN: Premature Project Completion
- NEVER transition to PROJECT_DONE after Phase 1
- NEVER transition to PROJECT_DONE without all 5 phases done
- NEVER assume PR plan creation means project end

### ❌ FORBIDDEN: Skipping Phases
- NEVER jump from Phase 1 to PROJECT_DONE
- NEVER skip phase integration steps
- NEVER ignore remaining waves/phases

### ✅ CORRECT: Incremental PR Planning
- Create PR plans after each phase
- Update PR plans after major milestones
- Keep PR plans current with progress
- Only mark PROJECT_DONE when truly complete

## Validation Checklist

Before transitioning from PR_PLAN_CREATION:
- [ ] Check current_phase vs total_phases
- [ ] Check current_wave vs waves_in_phase
- [ ] Verify all efforts in current scope complete
- [ ] Confirm PR plan accurately reflects work done
- [ ] Validate next state makes logical sense
- [ ] Document transition rationale

## Error Conditions

Transition to ERROR_RECOVERY if:
- Cannot determine project progress
- State file missing critical fields
- Phase plans cannot be loaded
- Conflicting state information

## References
- R279: PR Plan Creation Protocol
- R206: State Machine Validation
- State Machine: /state-machines/software-factory-3.0-state-machine.json
- Orchestrator Config: /.claude/agents/orchestrator.md

## Automation Flag

```bash
# After creating PR plan, determine next action
echo "✅ PR Plan created successfully"

# Check if more work remains
if [ "$MORE_WAVES" = true ] || [ "$MORE_PHASES" = true ]; then
    echo "Continuing with next wave/phase..."
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue to next wave/phase
else
    echo "All implementation complete, moving to PROJECT_DONE"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue to PROJECT_DONE state
fi
```
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
