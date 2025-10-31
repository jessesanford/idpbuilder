# RECOMMENDED RULE CHANGES

**Date**: 2025-10-29
**Context**: Invalid INTEGRATE_WAVE_EFFORTS transition incident
**Purpose**: Prevent future semantic validation failures

---

## EXECUTIVE SUMMARY

The investigation revealed that State Manager validates SYNTACTIC correctness (allowed_transitions) but NOT SEMANTIC correctness (business logic preconditions). This gap allowed an invalid transition to INTEGRATE_WAVE_EFFORTS when no efforts existed to integrate.

**Recommended**: Create new rule R540 requiring semantic precondition validation.

---

## PROPOSED NEW RULE: R540

### R540: Semantic Precondition Validation (BLOCKING)

**File**: rule-library/R540-semantic-precondition-validation.md

```markdown
# RULE R540: Semantic Precondition Validation

**Criticality**: 🚨🚨🚨 BLOCKING
**Agent**: State Manager
**Domain**: State transition validation
**Created**: 2025-10-29
**Reason**: Prevent invalid state transitions due to unmet business logic prerequisites

---

## Rule Statement

State Manager MUST validate both SYNTACTIC and SEMANTIC correctness before approving state transitions:

1. **Syntactic Validation** (existing):
   - Current state exists in state machine ✅
   - Proposed state exists in state machine ✅
   - Transition is in allowed_transitions list ✅
   - Mandatory sequence enforcement ✅

2. **Semantic Validation** (NEW):
   - Target state's requires.conditions are met ✅
   - Business logic preconditions satisfied ✅
   - Data prerequisites exist ✅
   - State-specific guards pass ✅

---

## Rationale

**Problem**: State machine defines allowed_transitions (syntax) but cannot express runtime business logic (semantics).

**Example Failure**:
- Transition: START_WAVE_ITERATION → INTEGRATE_WAVE_EFFORTS
- Syntactically valid: INTEGRATE_WAVE_EFFORTS is in allowed_transitions ✅
- Semantically invalid: efforts_to_integrate array is EMPTY ❌
- Result: State Manager approved invalid transition

**Solution**: Add semantic validation layer that checks state-specific preconditions.

---

## Implementation

### Phase 1: Core State Validations

Add semantic checks for these critical states:

**INTEGRATE_WAVE_EFFORTS**:
```bash
# Preconditions
[ $(jq '.entries[] | select(.iteration_level=="wave" and .active==true) | .efforts_to_integrate | length' integration-containers.json) -gt 0 ]
[ $(jq -r '.project_progression.current_wave.status' orchestrator-state-v3.json) != "RESET" ]
```

**INTEGRATE_PHASE_EFFORTS**:
```bash
[ $(jq '.entries[] | select(.iteration_level=="phase" and .active==true) | .waves_to_integrate | length' integration-containers.json) -gt 0 ]
[ $(jq -r '.current_phase' orchestrator-state-v3.json) -ge 1 ]
```

**INTEGRATE_PROJECT_EFFORTS**:
```bash
[ $(jq '.entries[] | select(.iteration_level=="project" and .active==true) | .phases_to_integrate | length' integration-containers.json) -gt 0 ]
```

**SPAWN_SW_ENGINEERS**:
```bash
[ -f "wave-plans/WAVE-${WAVE}-IMPLEMENTATION.md" ]
[ $(jq -r '.project_progression.efforts_pending | length' orchestrator-state-v3.json) -gt 0 ]
```

### Phase 2: Validation Function Library

**Location**: lib/semantic-validation-lib.sh

```bash
#!/bin/bash
# Semantic validation functions for state transitions

validate_integrate_wave_efforts() {
    local efforts_count=$(jq -r '.entries[] | select(.iteration_level=="wave" and .active==true) | .efforts_to_integrate | length' integration-containers.json)
    
    if [ "$efforts_count" -eq 0 ]; then
        echo "❌ Cannot integrate wave: efforts_to_integrate is empty"
        return 1
    fi
    
    local wave_status=$(jq -r '.project_progression.current_wave.status' orchestrator-state-v3.json)
    if [ "$wave_status" == "RESET" ]; then
        echo "❌ Cannot integrate wave: wave status is RESET"
        return 1
    fi
    
    echo "✅ Wave integration preconditions satisfied"
    return 0
}

validate_spawn_sw_engineers() {
    local wave=$(jq -r '.current_wave' orchestrator-state-v3.json)
    local wave_plan="wave-plans/WAVE-${wave}-IMPLEMENTATION.md"
    
    if [ ! -f "$wave_plan" ]; then
        echo "❌ Cannot spawn SW Engineers: $wave_plan does not exist"
        return 1
    fi
    
    local efforts_pending=$(jq -r '.project_progression.efforts_pending | length' orchestrator-state-v3.json)
    if [ "$efforts_pending" -eq 0 ]; then
        echo "❌ Cannot spawn SW Engineers: no efforts pending"
        return 1
    fi
    
    echo "✅ SW Engineer spawn preconditions satisfied"
    return 0
}
```

### Phase 3: State Manager Integration

**Location**: .claude/agents/state-manager.md (after Step 2b, before Step 3)

```bash
# Step 2c: Validate Semantic Preconditions
source "$CLAUDE_PROJECT_DIR/lib/semantic-validation-lib.sh"

case "$PROPOSED_NEXT_STATE" in
    INTEGRATE_WAVE_EFFORTS)
        if ! validate_integrate_wave_efforts; then
            DECISION="ERROR_RECOVERY"
            PROPOSAL_REJECTED=true
            PROPOSAL_REJECTED_REASON="Semantic preconditions not met for INTEGRATE_WAVE_EFFORTS"
            CONTINUE_SOFTWARE_FACTORY="FALSE"
        fi
        ;;
    
    INTEGRATE_PHASE_EFFORTS)
        if ! validate_integrate_phase_efforts; then
            DECISION="ERROR_RECOVERY"
            PROPOSAL_REJECTED=true
            PROPOSAL_REJECTED_REASON="Semantic preconditions not met for INTEGRATE_PHASE_EFFORTS"
            CONTINUE_SOFTWARE_FACTORY="FALSE"
        fi
        ;;
    
    SPAWN_SW_ENGINEERS)
        if ! validate_spawn_sw_engineers; then
            DECISION="ERROR_RECOVERY"
            PROPOSAL_REJECTED=true
            PROPOSAL_REJECTED_REASON="Semantic preconditions not met for SPAWN_SW_ENGINEERS"
            CONTINUE_SOFTWARE_FACTORY="FALSE"
        fi
        ;;
esac
```

---

## Testing Requirements

### Test Case 1: Empty Efforts Array
```bash
# Setup
orchestrator-state-v3.json: current_state = START_WAVE_ITERATION
integration-containers.json: efforts_to_integrate = []

# Action
Orchestrator proposes INTEGRATE_WAVE_EFFORTS

# Expected (with R540)
State Manager REJECTS
DECISION = ERROR_RECOVERY
PROPOSAL_REJECTED = true
REASON = "efforts_to_integrate is empty"
```

### Test Case 2: Wave Status RESET
```bash
# Setup
current_wave.status = "RESET"
integration-containers.json: efforts_to_integrate = []

# Action
Orchestrator proposes INTEGRATE_WAVE_EFFORTS

# Expected (with R540)
State Manager REJECTS
DECISION = ERROR_RECOVERY
REASON = "wave status is RESET"
```

### Test Case 3: Valid Transition
```bash
# Setup
integration-containers.json: efforts_to_integrate = ["1.2.1", "1.2.2"]
current_wave.status = "IN_PROGRESS"

# Action
Orchestrator proposes INTEGRATE_WAVE_EFFORTS

# Expected (with R540)
State Manager ACCEPTS
DECISION = INTEGRATE_WAVE_EFFORTS
PROPOSAL_REJECTED = false
```

---

## Success Criteria

R540 is successful when:

✅ State Manager validates semantic preconditions for all critical states
✅ Invalid transitions blocked even if syntactically valid
✅ Validation library is reusable and extensible
✅ Clear error messages explain WHY transition rejected
✅ Zero incidents of invalid state transitions
✅ Graceful degradation (fall back to ERROR_RECOVERY)

---

## Rollout Plan

### Phase 1: Core Validations (Week 1)
- Create semantic-validation-lib.sh
- Add INTEGRATE_WAVE_EFFORTS validation
- Add INTEGRATE_PHASE_EFFORTS validation
- Add INTEGRATE_PROJECT_EFFORTS validation
- Test in isolated environment

### Phase 2: State Manager Integration (Week 2)
- Update state-manager.md with Step 2c
- Add function library sourcing
- Add case statement for state-specific checks
- Test with real state transitions

### Phase 3: Extended Validations (Week 3)
- Add SPAWN_SW_ENGINEERS validation
- Add CREATE_NEXT_INFRASTRUCTURE validation
- Add BUILD_VALIDATION validation
- Comprehensive testing

### Phase 4: Documentation and Training (Week 4)
- Document all semantic checks
- Update agent training materials
- Create troubleshooting guide
- Post-implementation review

---

## Related Rules

- **R288**: State File Update and Commit (complementary)
- **R506**: Absolute Prohibition on Pre-Commit Bypass (complementary)
- **R516**: State Creation and Design Protocol (foundation)
- **R517**: State Manager Consultation Protocol (integration point)

---

## Maintenance

**Adding New Semantic Checks**:

1. Identify state requiring validation
2. Add function to semantic-validation-lib.sh
3. Add case to state-manager.md Step 2c
4. Write test cases
5. Document in this rule
6. Update state machine JSON (requires.conditions)

**Review Schedule**: Quarterly review of semantic validations to identify gaps

---

**Status**: PROPOSED - Ready for approval and implementation
```

---

## PROPOSED STATE MACHINE ENHANCEMENTS

### Update INTEGRATE_WAVE_EFFORTS Definition

**Current**:
```json
{
  "INTEGRATE_WAVE_EFFORTS": {
    "requires": {
      "conditions": [
        "All effort branches ready for integration",
        "Iteration started"
      ]
    }
  }
}
```

**Proposed**:
```json
{
  "INTEGRATE_WAVE_EFFORTS": {
    "requires": {
      "conditions": [
        "efforts_to_integrate array is non-empty (at least 1 effort)",
        "All effort branches exist in repository",
        "All effort branches have completed implementation",
        "Iteration started (iteration > 0)",
        "Wave status is IN_PROGRESS or READY (not RESET or ERROR)"
      ]
    },
    "guards": {
      "entry": "efforts_to_integrate.length > 0 && wave_status != 'RESET'",
      "validation": "semantic_validation.validate_integrate_wave_efforts()"
    }
  }
}
```

### Add Guards Section to State Machine

**New Section**:
```json
{
  "guards": {
    "description": "Runtime guards that must pass before state entry",
    "implementation": "lib/semantic-validation-lib.sh",
    "enforcement": "State Manager validates guards during SHUTDOWN_CONSULTATION"
  }
}
```

---

## IMPACT ANALYSIS

### Benefits

1. **Prevention**: Blocks invalid transitions before they occur
2. **Clarity**: Clear error messages explain failures
3. **Defense in Depth**: Multiple validation layers
4. **Maintainability**: Centralized validation logic
5. **Extensibility**: Easy to add new validations

### Risks

1. **Performance**: Additional validation overhead (~100ms per transition)
2. **Complexity**: More code to maintain
3. **False Positives**: Over-strict validation might block valid transitions

**Mitigation**:
- Performance: Acceptable for human-driven workflow
- Complexity: Well-structured library with clear documentation
- False Positives: Careful design + comprehensive testing

### Alternatives Considered

**Alternative 1**: Orchestrator-only validation
- ❌ No defense in depth
- ❌ Easy to bypass
- ❌ Duplicated logic across agent types

**Alternative 2**: Pre-commit hook validation
- ❌ Too late (change already attempted)
- ❌ Harder to provide good error messages
- ❌ Can't prevent transition proposal

**Alternative 3**: State machine code generation
- ❌ High complexity
- ❌ Requires new tooling
- ❌ Harder to debug

**Recommendation**: Implement R540 (State Manager semantic validation) - Best balance of simplicity, effectiveness, and maintainability.

---

## APPROVAL STATUS

**Proposed By**: Software Factory Manager
**Date**: 2025-10-29
**Approval Required From**:
- [ ] Architect (system design impact)
- [ ] Orchestrator Lead (workflow impact)
- [ ] State Manager Team (implementation responsibility)

**Estimated Effort**: 2-3 developer days
**Priority**: HIGH (prevents critical failures)

---

**Next Steps**:
1. Review and approve R540
2. Create semantic-validation-lib.sh
3. Update state-manager.md
4. Test in isolated environment
5. Deploy to production
6. Monitor for 1 sprint
7. Extend to additional states
