# PLANNING State Migration Guide

## ⚠️ CRITICAL: PLANNING State Has Been Removed

The `PLANNING` state was identified as orphaned (no valid transitions) and has been removed from the Software Factory 2.0 state machine as of 2025-08-30.

## Why Was PLANNING Removed?

1. **No Valid Transitions**: The state had no transitions TO or FROM it
2. **Confusion Source**: Orchestrators were attempting to use it incorrectly
3. **Redundant**: Planning activities are properly handled by other states
4. **Dead Code**: The state could never be reached through normal flow

## Correct Planning States to Use

### For Orchestrator Planning Activities

| Planning Type | Correct State | Purpose |
|--------------|---------------|---------|
| Phase Architecture | `SPAWN_ARCHITECT_PHASE_PLANNING` | Request architect to create phase architecture |
| Wave Architecture | `SPAWN_ARCHITECT_WAVE_PLANNING` | Request architect to create wave architecture |
| Phase Implementation | `SPAWN_CODE_REVIEWER_PHASE_IMPL` | Request code reviewer to translate architecture to implementation |
| Wave Implementation | `SPAWN_CODE_REVIEWER_WAVE_IMPL` | Request code reviewer to translate wave architecture |
| Effort Planning | `SPAWN_CODE_REVIEWERS_EFFORT_PLANNING` | Spawn code reviewers to create individual effort plans |

### State Flow for Planning

```
Phase Planning:
INIT → SPAWN_ARCHITECT_PHASE_PLANNING → WAITING_FOR_ARCHITECTURE_PLAN 
     → SPAWN_CODE_REVIEWER_PHASE_IMPL → WAITING_FOR_IMPLEMENTATION_PLAN

Wave Planning:
WAVE_START → SPAWN_ARCHITECT_WAVE_PLANNING → WAITING_FOR_ARCHITECTURE_PLAN
           → SPAWN_CODE_REVIEWER_WAVE_IMPL → INJECT_WAVE_METADATA 
           → WAITING_FOR_IMPLEMENTATION_PLAN

Effort Planning:
ANALYZE_CODE_REVIEWER_PARALLELIZATION → SPAWN_CODE_REVIEWERS_EFFORT_PLANNING 
                                      → WAITING_FOR_EFFORT_PLANS
```

## Migration Instructions

### If Your State File Shows `current_state: PLANNING`

1. **Immediate Action Required**:
   ```yaml
   # Change from:
   current_state: PLANNING
   
   # To:
   current_state: ERROR_RECOVERY
   ```

2. **Assess Intent**:
   - What planning activity were you trying to perform?
   - Which agent should be doing the planning?
   - What is the correct state from the table above?

3. **Transition to Correct State**:
   ```yaml
   # After assessment, transition to appropriate state:
   current_state: SPAWN_CODE_REVIEWERS_EFFORT_PLANNING  # For effort planning
   # OR
   current_state: SPAWN_ARCHITECT_PHASE_PLANNING  # For phase architecture
   # OR
   current_state: WAVE_START  # To begin wave properly
   ```

### If Your Code References PLANNING State

1. **In Orchestrator Logic**:
   ```python
   # WRONG - Old way
   if need_planning:
       transition_to("PLANNING")
   
   # CORRECT - New way
   if need_effort_planning:
       transition_to("SPAWN_CODE_REVIEWERS_EFFORT_PLANNING")
   elif need_phase_architecture:
       transition_to("SPAWN_ARCHITECT_PHASE_PLANNING")
   elif need_wave_architecture:
       transition_to("SPAWN_ARCHITECT_WAVE_PLANNING")
   ```

2. **In State Validation**:
   ```bash
   # Remove PLANNING from valid states list
   VALID_STATES="INIT WAVE_START SPAWN_AGENTS MONITOR..."  # No PLANNING
   ```

## Validation

### Check for Orphaned State References
```bash
# Find any remaining references to PLANNING state
grep -r "current_state.*PLANNING" --include="*.yaml" .
grep -r "transition.*PLANNING" --include="*.md" --include="*.sh" .
grep -r "'PLANNING'" --include="*.py" --include="*.js" .
```

### Verify Correct Planning State Usage
```bash
# Ensure using correct planning states
grep -r "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING" orchestrator-state.yaml
grep -r "SPAWN_ARCHITECT_PHASE_PLANNING" orchestrator-state.yaml
```

## Common Mistakes to Avoid

### ❌ DON'T: Try to Create Generic Planning State
```yaml
# WRONG
current_state: PLANNING
planning_type: "effort"  # This doesn't work
```

### ✅ DO: Use Specific Planning States
```yaml
# CORRECT
current_state: SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
# State itself indicates planning type
```

### ❌ DON'T: Skip Planning States
```yaml
# WRONG - Can't skip from SETUP to SPAWN_AGENTS
SETUP_EFFORT_INFRASTRUCTURE → SPAWN_AGENTS
```

### ✅ DO: Follow Mandatory Sequence (R234)
```yaml
# CORRECT - Must traverse all states
SETUP_EFFORT_INFRASTRUCTURE 
→ ANALYZE_CODE_REVIEWER_PARALLELIZATION
→ SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
→ WAITING_FOR_EFFORT_PLANS
→ ANALYZE_IMPLEMENTATION_PARALLELIZATION
→ SPAWN_AGENTS
```

## Timeline

- **2025-08-30**: PLANNING state removed from state machine
- **2025-08-30**: State directory archived as PLANNING.DEPRECATED-20250830
- **2025-08-30**: R289 rule created for orphaned state detection
- **Immediate**: All orchestrators must stop using PLANNING state

## Support

If you encounter issues after this migration:

1. Check ERROR_RECOVERY state rules
2. Review this migration guide
3. Consult SOFTWARE-FACTORY-STATE-MACHINE.md for valid states
4. Use utilities/detect-orphaned-states.sh to validate

## Related Documentation

- R289: Orphaned State Detection and Prevention
- R234: Mandatory State Traversal
- R206: State Machine Transition Validation
- SOFTWARE-FACTORY-STATE-MACHINE.md