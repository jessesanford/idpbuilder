# 🚨🚨🚨 BLOCKING RULE R285: Mandatory Phase Integration Before Assessment

## Rule Statement
Phase integration MUST occur BEFORE architect phase assessment in ALL cases. The architect must assess the INTEGRATED phase, not individual wave branches.

## Criticality: 🚨🚨🚨 BLOCKING

## Enforcement

### Required State Transitions
The orchestrator MUST follow this exact sequence when the last wave of a phase completes:

```
REVIEW_WAVE_ARCHITECTURE (last wave) 
    → INTEGRATE_PHASE_WAVES (create integration branch)
    → SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN (plan merges)
    → WAITING_FOR_PHASE_MERGE_PLAN (wait for plan)
    → SPAWN_INTEGRATION_AGENT_PHASE (execute merges)
    → MONITORING_INTEGRATE_PHASE_WAVES (monitor integration)
    → SPAWN_ARCHITECT_PHASE_ASSESSMENT (assess integrated phase)
```

### Forbidden Transitions
- ❌ `REVIEW_WAVE_ARCHITECTURE` → `SPAWN_ARCHITECT_PHASE_ASSESSMENT` (skips integration)
- ❌ `REVIEW_WAVE_ARCHITECTURE` → `COMPLETE_PHASE` (skips both integration and assessment)
- ❌ `INTEGRATE_PHASE_WAVES` → `SPAWN_ARCHITECT_PHASE_ASSESSMENT` (skips merge execution)

## Context Detection

The INTEGRATE_PHASE_WAVES state must detect its entry context:

1. **From REVIEW_WAVE_ARCHITECTURE** (normal flow):
   - This is standard phase completion
   - Integrate all wave branches for the phase
   - Branch name: `phase-{N}-integration`

2. **From ERROR_RECOVERY** (fix flow - R259):
   - This is post-assessment fix integration
   - Integrate wave branches + fix branches
   - Branch name: `phase{N}-post-fixes-integration-{TIMESTAMP}`

## Integration Requirements

### Phase Integration Must:
1. Create clean branch from main HEAD
2. Merge each wave's integration branch sequentially
3. Run tests after each wave merge
4. Handle merge conflicts if they arise
5. Create PHASE-{N}-INTEGRATE_WAVE_EFFORTS.md report
6. Push integrated branch to remote
7. Record in orchestrator-state-v3.json

### Architect Assessment Must:
1. Review the integrated phase branch (not individual waves)
2. Verify all waves work together as integrated whole
3. Check for phase-level architectural consistency
4. Validate phase-level tests pass
5. Create PHASE-{N}-ASSESSMENT-REPORT.md

## Rationale

1. **Logical Consistency**: Matches wave-level pattern where integration precedes review
2. **Quality Assurance**: Integration issues found before assessment, not after
3. **Proper Assessment**: Architect reviews actual integrated code, not theoretical combination
4. **Early Detection**: Merge conflicts and integration bugs caught immediately
5. **Architectural Validation**: Phase-level patterns only visible in integrated code

## Example Scenario

### Correct Flow:
```yaml
# Wave 3 (last wave of Phase 2) completes
current_state: REVIEW_WAVE_ARCHITECTURE
wave_review_decision: PROCEED_INTEGRATE_PHASE_WAVES

# Transition to phase integration
current_state: INTEGRATE_PHASE_WAVES
phase_integration_branch: phase-2-integration

# Create merge plan
current_state: SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN

# Execute integration
current_state: SPAWN_INTEGRATION_AGENT_PHASE

# Monitor integration
current_state: MONITORING_INTEGRATE_PHASE_WAVES
integration_status: PROJECT_DONE

# NOW architect can assess
current_state: SPAWN_ARCHITECT_PHASE_ASSESSMENT
assessment_target: phase-2-integration  # Integrated branch!
```

### Incorrect Flow (FORBIDDEN):
```yaml
# Wave 3 completes
current_state: REVIEW_WAVE_ARCHITECTURE
wave_review_decision: PROCEED_PHASE_ASSESSMENT  # WRONG!

# Skip integration - VIOLATION!
current_state: SPAWN_ARCHITECT_PHASE_ASSESSMENT
assessment_targets:  # Multiple unintegrated branches - WRONG!
  - wave-1-integration
  - wave-2-integration  
  - wave-3-integration
```

## Grading Impact

### Violations:
- Skipping phase integration entirely: **-100% AUTOMATIC FAILURE**
- Direct REVIEW_WAVE_ARCHITECTURE → SPAWN_ARCHITECT_PHASE_ASSESSMENT: **-100% FAILURE**
- Architect assessing unintegrated waves: **-50% penalty**
- Missing integration report: **-25% penalty**
- Integration branch not pushed: **-20% penalty**

### Success Criteria:
- ✅ Phase integration occurs before assessment
- ✅ All wave branches properly merged
- ✅ Integration tests pass
- ✅ Architect reviews integrated branch
- ✅ Integration tracked in state file

## Related Rules

- **R259**: Mandatory phase integration after ERROR_RECOVERY fixes
- **R258**: Wave review report requirements and decisions
- **R257**: Phase assessment report requirements
- **R282**: Phase integration protocol and isolation
- **R206**: State machine validation requirements

## Implementation Notes

1. Update REVIEW_WAVE_ARCHITECTURE decision logic to include `PROCEED_INTEGRATE_PHASE_WAVES`
2. Ensure INTEGRATE_PHASE_WAVES handles both entry contexts correctly
3. Verify integration agent supports phase-level merges
4. Update orchestrator monitoring to track phase integration progress
5. Ensure architect prompt includes integrated branch reference

## Verification

```bash
# Check for forbidden transition
grep "REVIEW_WAVE_ARCHITECTURE.*SPAWN_ARCHITECT_PHASE_ASSESSMENT" orchestrator-state-v3.json
# Should return NOTHING if compliant

# Verify integration happened
grep "phase_integration_branches" orchestrator-state-v3.json
# Should show integration branch for current phase

# Confirm architect assessed integrated branch
grep "assessment_target.*phase.*integration" phase-assessments/*/PHASE-*-ASSESSMENT-REPORT.md
# Should show integrated branch was assessed
```

## State Manager Coordination (SF 3.0)

State Manager coordinates phase integration checkpoints through iteration containers:
- **Validates** all phase waves integrated before REVIEW_PHASE_ARCHITECTURE transition
- **Checks** `integration-containers.json` for convergence status
- **Ensures** atomic state updates across all 4 files during phase completion
- **Guards** REVIEW_PHASE_ARCHITECTURE with `requires.conditions`: completed integration container

The integration container tracks iterations until convergence (all bugs resolved in integrated phase branch).

See: `integration-containers.json`, R327 (iteration container management), `state-machines/software-factory-3.0-state-machine.json` (REVIEW_PHASE_ARCHITECTURE guards)

## Summary

This rule ensures logical consistency in the Software Factory flow. Just as waves are integrated before review, phases MUST be integrated before assessment. This guarantees the architect evaluates the actual integrated codebase that will be deployed, not a theoretical combination of separate branches.

**Remember**: Integration reveals issues that isolation hides. Always integrate before assessment!