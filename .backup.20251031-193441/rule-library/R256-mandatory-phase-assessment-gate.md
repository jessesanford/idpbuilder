# 🚨🚨🚨 RULE R256 - Mandatory Phase Assessment Gate [BLOCKING]

## Rule Statement
**NO PHASE CAN BE MARKED COMPLETE WITHOUT ARCHITECT PHASE-LEVEL ASSESSMENT**

The orchestrator MUST NOT transition to PROJECT_DONE without going through:
1. SPAWN_ARCHITECT_PHASE_ASSESSMENT (request assessment)
2. WAITING_FOR_PHASE_ASSESSMENT (receive decision)  
3. COMPLETE_PHASE (finalize phase)
4. Only then → PROJECT_DONE

## Why This Rule Exists

Previously, the orchestrator could transition directly from REVIEW_WAVE_ARCHITECTURE to PROJECT_DONE when the last wave completed. This was WRONG because:
- No phase-level architectural validation occurred
- No comprehensive feature completeness check
- No API stability assessment
- No phase-wide integration validation
- Premature celebration without true completion

## Enforcement

### ❌ FORBIDDEN (Automatic Failure):
```
REVIEW_WAVE_ARCHITECTURE → PROJECT_DONE  # NEVER ALLOWED
INTEGRATE_WAVE_EFFORTS → PROJECT_DONE  # NEVER ALLOWED
WAVE_COMPLETE → PROJECT_DONE  # NEVER ALLOWED
```

### ✅ REQUIRED Flow:
```
REVIEW_WAVE_ARCHITECTURE (last wave) → SPAWN_ARCHITECT_PHASE_ASSESSMENT
SPAWN_ARCHITECT_PHASE_ASSESSMENT → WAITING_FOR_PHASE_ASSESSMENT
WAITING_FOR_PHASE_ASSESSMENT → COMPLETE_PHASE (if passed)
COMPLETE_PHASE → PROJECT_DONE (phase truly complete)
```

## Implementation Requirements

### In REVIEW_WAVE_ARCHITECTURE State:
```bash
# Determine if this is the last wave
if [ "$CURRENT_WAVE" == "$LAST_WAVE_OF_PHASE" ]; then
    # MUST transition to phase assessment, NOT PROJECT_DONE
    transition_to "SPAWN_ARCHITECT_PHASE_ASSESSMENT"
else
    transition_to "WAVE_START"  # Continue to next wave
fi
```

### In SPAWN_ARCHITECT_PHASE_ASSESSMENT:
- Spawn architect with COMPLETE phase context
- Include ALL wave integration branches
- Request comprehensive phase validation
- This assessment GATES the PROJECT_DONE state

### In WAITING_FOR_PHASE_ASSESSMENT:
- Process architect's phase-level decision
- If PASS → COMPLETE_PHASE
- If FAIL → ERROR_RECOVERY
- Never skip to PROJECT_DONE

### In COMPLETE_PHASE:
- Create final phase integration branch
- Document all phase achievements
- Generate completion metrics
- Only NOW can transition to PROJECT_DONE

## Validation

The orchestrator state machine MUST:
1. Remove any direct paths from non-COMPLETE_PHASE states to PROJECT_DONE
2. Add SPAWN_ARCHITECT_PHASE_ASSESSMENT as mandatory for last wave
3. Ensure COMPLETE_PHASE is the ONLY state that transitions to PROJECT_DONE

## Violations

**Severity**: BLOCKING - System design failure

**Detection**:
- Any transition to PROJECT_DONE not from COMPLETE_PHASE
- Missing phase assessment for completed phase
- Orchestrator celebrating without architect approval

**Consequences**:
- Immediate orchestrator failure
- Phase marked incomplete
- Manual intervention required
- Grading: -100 points (critical failure)

## Related Rules
- R057 - Wave Review Authority (wave-level)
- R058 - Phase Assessment Responsibility (phase-level)
- R206 - State Machine Validation
- R233 - All States Require Immediate Action

## Audit Command

```bash
# Verify no direct PROJECT_DONE transitions exist
grep -r "transition.*PROJECT_DONE" agent-states/software-factory/orchestrator/ | \
  grep -v "COMPLETE_PHASE" && echo "❌ VIOLATION: Direct PROJECT_DONE transition found!"

# Verify REVIEW_WAVE_ARCHITECTURE includes phase assessment logic
grep "SPAWN_ARCHITECT_PHASE_ASSESSMENT" \
  agent-states/software-factory/orchestrator/REVIEW_WAVE_ARCHITECTURE/rules.md || \
  echo "❌ VIOLATION: REVIEW_WAVE_ARCHITECTURE missing phase assessment transition!"
```

## Summary

**The Phase Assessment Gate is ABSOLUTE**:
- Every phase MUST be assessed by architect
- PROJECT_DONE is BLOCKED without phase assessment
- No shortcuts, no exceptions, no bypassing
- The architect is the gatekeeper of phase completion

This ensures architectural integrity, feature completeness, and true phase readiness before any celebration of PROJECT_DONE.