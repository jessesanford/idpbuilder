# 🚨🚨🚨 RULE R256 - Mandatory Phase Assessment Gate [BLOCKING]

## Rule Statement
**NO PHASE CAN BE MARKED COMPLETE WITHOUT ARCHITECT PHASE-LEVEL ASSESSMENT**

The orchestrator MUST NOT transition to SUCCESS without going through:
1. SPAWN_ARCHITECT_PHASE_ASSESSMENT (request assessment)
2. WAITING_FOR_PHASE_ASSESSMENT (receive decision)  
3. PHASE_COMPLETE (finalize phase)
4. Only then → SUCCESS

## Why This Rule Exists

Previously, the orchestrator could transition directly from WAVE_REVIEW to SUCCESS when the last wave completed. This was WRONG because:
- No phase-level architectural validation occurred
- No comprehensive feature completeness check
- No API stability assessment
- No phase-wide integration validation
- Premature celebration without true completion

## Enforcement

### ❌ FORBIDDEN (Automatic Failure):
```
WAVE_REVIEW → SUCCESS  # NEVER ALLOWED
INTEGRATION → SUCCESS  # NEVER ALLOWED
WAVE_COMPLETE → SUCCESS  # NEVER ALLOWED
```

### ✅ REQUIRED Flow:
```
WAVE_REVIEW (last wave) → SPAWN_ARCHITECT_PHASE_ASSESSMENT
SPAWN_ARCHITECT_PHASE_ASSESSMENT → WAITING_FOR_PHASE_ASSESSMENT
WAITING_FOR_PHASE_ASSESSMENT → PHASE_COMPLETE (if passed)
PHASE_COMPLETE → SUCCESS (phase truly complete)
```

## Implementation Requirements

### In WAVE_REVIEW State:
```bash
# Determine if this is the last wave
if [ "$CURRENT_WAVE" == "$LAST_WAVE_OF_PHASE" ]; then
    # MUST transition to phase assessment, NOT SUCCESS
    transition_to "SPAWN_ARCHITECT_PHASE_ASSESSMENT"
else
    transition_to "WAVE_START"  # Continue to next wave
fi
```

### In SPAWN_ARCHITECT_PHASE_ASSESSMENT:
- Spawn architect with COMPLETE phase context
- Include ALL wave integration branches
- Request comprehensive phase validation
- This assessment GATES the SUCCESS state

### In WAITING_FOR_PHASE_ASSESSMENT:
- Process architect's phase-level decision
- If PASS → PHASE_COMPLETE
- If FAIL → ERROR_RECOVERY
- Never skip to SUCCESS

### In PHASE_COMPLETE:
- Create final phase integration branch
- Document all phase achievements
- Generate completion metrics
- Only NOW can transition to SUCCESS

## Validation

The orchestrator state machine MUST:
1. Remove any direct paths from non-PHASE_COMPLETE states to SUCCESS
2. Add SPAWN_ARCHITECT_PHASE_ASSESSMENT as mandatory for last wave
3. Ensure PHASE_COMPLETE is the ONLY state that transitions to SUCCESS

## Violations

**Severity**: BLOCKING - System design failure

**Detection**:
- Any transition to SUCCESS not from PHASE_COMPLETE
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
# Verify no direct SUCCESS transitions exist
grep -r "transition.*SUCCESS" agent-states/orchestrator/ | \
  grep -v "PHASE_COMPLETE" && echo "❌ VIOLATION: Direct SUCCESS transition found!"

# Verify WAVE_REVIEW includes phase assessment logic
grep "SPAWN_ARCHITECT_PHASE_ASSESSMENT" \
  agent-states/orchestrator/WAVE_REVIEW/rules.md || \
  echo "❌ VIOLATION: WAVE_REVIEW missing phase assessment transition!"
```

## Summary

**The Phase Assessment Gate is ABSOLUTE**:
- Every phase MUST be assessed by architect
- SUCCESS is BLOCKED without phase assessment
- No shortcuts, no exceptions, no bypassing
- The architect is the gatekeeper of phase completion

This ensures architectural integrity, feature completeness, and true phase readiness before any celebration of SUCCESS.