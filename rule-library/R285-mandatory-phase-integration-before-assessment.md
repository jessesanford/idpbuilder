# 🚨🚨🚨 BLOCKING RULE R285: Mandatory Phase Integration Before Assessment

## Rule Statement
Phase integration MUST occur BEFORE architect phase assessment in ALL cases. The architect must assess the INTEGRATED phase, not individual wave branches.

## Criticality: 🚨🚨🚨 BLOCKING

## Enforcement

### Required State Transitions
The orchestrator MUST follow this exact sequence when the last wave of a phase completes:

```
WAVE_REVIEW (last wave) 
    → PHASE_INTEGRATION (create integration branch)
    → SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN (plan merges)
    → WAITING_FOR_PHASE_MERGE_PLAN (wait for plan)
    → SPAWN_INTEGRATION_AGENT_PHASE (execute merges)
    → MONITORING_PHASE_INTEGRATION (monitor integration)
    → SPAWN_ARCHITECT_PHASE_ASSESSMENT (assess integrated phase)
```

### Forbidden Transitions
- ❌ `WAVE_REVIEW` → `SPAWN_ARCHITECT_PHASE_ASSESSMENT` (skips integration)
- ❌ `WAVE_REVIEW` → `PHASE_COMPLETE` (skips both integration and assessment)
- ❌ `PHASE_INTEGRATION` → `SPAWN_ARCHITECT_PHASE_ASSESSMENT` (skips merge execution)

## Context Detection

The PHASE_INTEGRATION state must detect its entry context:

1. **From WAVE_REVIEW** (normal flow):
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
5. Create PHASE-{N}-INTEGRATION.md report
6. Push integrated branch to remote
7. Record in orchestrator-state.json

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
current_state: WAVE_REVIEW
wave_review_decision: PROCEED_PHASE_INTEGRATION

# Transition to phase integration
current_state: PHASE_INTEGRATION
phase_integration_branch: phase-2-integration

# Create merge plan
current_state: SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN

# Execute integration
current_state: SPAWN_INTEGRATION_AGENT_PHASE

# Monitor integration
current_state: MONITORING_PHASE_INTEGRATION
integration_status: SUCCESS

# NOW architect can assess
current_state: SPAWN_ARCHITECT_PHASE_ASSESSMENT
assessment_target: phase-2-integration  # Integrated branch!
```

### Incorrect Flow (FORBIDDEN):
```yaml
# Wave 3 completes
current_state: WAVE_REVIEW
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
- Direct WAVE_REVIEW → SPAWN_ARCHITECT_PHASE_ASSESSMENT: **-100% FAILURE**
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

1. Update WAVE_REVIEW decision logic to include `PROCEED_PHASE_INTEGRATION`
2. Ensure PHASE_INTEGRATION handles both entry contexts correctly
3. Verify integration agent supports phase-level merges
4. Update orchestrator monitoring to track phase integration progress
5. Ensure architect prompt includes integrated branch reference

## Verification

```bash
# Check for forbidden transition
grep "WAVE_REVIEW.*SPAWN_ARCHITECT_PHASE_ASSESSMENT" orchestrator-state.json
# Should return NOTHING if compliant

# Verify integration happened
grep "phase_integration_branches" orchestrator-state.json
# Should show integration branch for current phase

# Confirm architect assessed integrated branch
grep "assessment_target.*phase.*integration" phase-assessments/*/PHASE-*-ASSESSMENT-REPORT.md
# Should show integrated branch was assessed
```

## Summary

This rule ensures logical consistency in the Software Factory flow. Just as waves are integrated before review, phases MUST be integrated before assessment. This guarantees the architect evaluates the actual integrated codebase that will be deployed, not a theoretical combination of separate branches.

**Remember**: Integration reveals issues that isolation hides. Always integrate before assessment!