# 🚨🚨🚨 RULE R259 - Mandatory Phase Re-Integration After Fixes (SF 3.0)

**Rule ID:** R259
**Category:** Integration Requirements
**Criticality:** BLOCKING
**Introduced:** v2.0.0
**Last Updated:** 2025-10-20 (SF 3.0 migration)

## Rule Statement

After completing phase-level fixes (`FIX_PHASE_UPSTREAM_BUGS` state in SF 3.0), the orchestrator MUST transition through the phase iteration sequence to re-integrate all waves before requesting reassessment. Specifically: `START_PHASE_ITERATION` → `INTEGRATE_PHASE_WAVES` → `REVIEW_PHASE_INTEGRATION`.

**SF 3.0 States**: `START_PHASE_ITERATION` → `INTEGRATE_PHASE_WAVES` → `REVIEW_PHASE_INTEGRATION` (iteration container pattern)

## Rationale

Phase assessment failures typically involve cross-wave issues that require comprehensive re-integration. After applying fixes to source branches (`FIX_PHASE_UPSTREAM_BUGS`), the phase iteration must restart with fresh integration of all waves. Without this re-integration, subsequent review would assess a stale integration branch, not the actual fixed code.

In SF 3.0's iteration container architecture, R259 enforces the critical step of re-integrating after fixes, ensuring convergence detection operates on current code state.

## Requirements

### 1. Mandatory State Transition (SF 3.0)
When `FIX_PHASE_UPSTREAM_BUGS` completes phase assessment fixes:
- ✅ MUST transition to `START_PHASE_ITERATION` (begin new iteration)
- ✅ Then to `INTEGRATE_PHASE_WAVES` (re-merge all waves)
- ✅ Then to `REVIEW_PHASE_INTEGRATION` (code review)
- ❌ MUST NOT go directly to `COMPLETE_PHASE`
- ❌ MUST NOT skip `INTEGRATE_PHASE_WAVES`

**Historical Note**: In SF 2.0, this used the deprecated `INTEGRATE_PHASE_WAVES` standalone state. SF 3.0 replaced this with the three-state iteration container pattern above.

### 2. Phase Integration Branch Creation (SF 3.0)
The `INTEGRATE_PHASE_WAVES` state MUST:
- Merge ALL wave integration branches sequentially into phase integration branch
- Use existing phase integration branch (from `integration-containers.json`)
- Resolve any merge conflicts
- Run build and validation
- Record integration outcome in iteration container
- Update `integration-containers.json` with iteration results

**Historical Note**: SF 2.0 created a new integration branch per attempt. SF 3.0 uses a single persistent phase integration branch (from `integration-containers.json`) that is reset for each iteration.

### 3. Verification Requirements
Before transitioning to reassessment:
- All Priority 1 issues from assessment report addressed
- Phase-level tests passing
- No unresolved merge conflicts
- Integration summary document created
- State file updated with branch details

### 4. State Flow

**SF 3.0 Flow**:
```
FIX_PHASE_UPSTREAM_BUGS → START_PHASE_ITERATION → INTEGRATE_PHASE_WAVES →
  REVIEW_PHASE_INTEGRATION → REVIEW_PHASE_ARCHITECTURE → (converged? → COMPLETE_PHASE)
```

**Historical Note**: SF 2.0 used `ERROR_RECOVERY` for all fixes and a standalone `INTEGRATE_PHASE_WAVES` state. SF 3.0 uses dedicated fix states (`FIX_PHASE_UPSTREAM_BUGS`) and the iteration container pattern.

## Implementation

### Orchestrator FIX_PHASE_UPSTREAM_BUGS State (SF 3.0)
```python
# In FIX_PHASE_UPSTREAM_BUGS state
# After completing all fixes
complete_phase_fixes()

# CRITICAL: Must return to START_PHASE_ITERATION
next_state = 'START_PHASE_ITERATION'  # NOT 'COMPLETE_PHASE'
update_state(next_state)
```

### START_PHASE_ITERATION State Actions (SF 3.0)
```bash
#!/bin/bash
# Immediate actions in START_PHASE_ITERATION

PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
ITERATION=$(jq -r '.phase_containers[] | select(.phase == '$PHASE') | .current_iteration' integration-containers.json)

# Increment iteration counter
NEW_ITERATION=$((ITERATION + 1))
jq --arg phase "$PHASE" --arg iter "$NEW_ITERATION" \
  '(.phase_containers[] | select(.phase == $phase) | .current_iteration) = ($iter | tonumber)' \
  integration-containers.json > tmp.json && mv tmp.json integration-containers.json

# Get existing phase integration branch from container
BRANCH=$(jq -r '.phase_containers[] | select(.phase == '$PHASE') | .integration_branch' integration-containers.json)

# Reset branch to clean state for re-iteration (if needed)
git checkout "$BRANCH"
git reset --hard origin/main

# Update state to proceed to integration
jq '.state_machine.current_state = "INTEGRATE_PHASE_WAVES"' orchestrator-state-v3.json
```

## Validation

### Detecting Violations
```python
def validate_r259_compliance(state_history):
    """Check for R259 violations in state transitions (SF 3.0)"""

    for i, state in enumerate(state_history[:-1]):
        next_state = state_history[i+1]

        # Check if FIX_PHASE_UPSTREAM_BUGS completed
        if state['name'] == 'FIX_PHASE_UPSTREAM_BUGS':
            # Next state MUST be START_PHASE_ITERATION
            if next_state['name'] != 'START_PHASE_ITERATION':
                return {
                    'compliant': False,
                    'violation': 'Skipped START_PHASE_ITERATION after phase fixes',
                    'found': f"{state['name']} → {next_state['name']}",
                    'expected': f"{state['name']} → START_PHASE_ITERATION"
                }

    return {'compliant': True}
```

### Verification Script
```bash
#!/bin/bash
# verify-r259-compliance.sh (SF 3.0)

# Check if phase integration container exists after fixes
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
CONTAINER=$(jq -r '.phase_containers[] | select(.phase == '$PHASE')' integration-containers.json)

if [ -z "$CONTAINER" ]; then
    echo "❌ R259 VIOLATION: No phase integration container found!"
    exit 1
fi

# Check that iteration counter incremented (indicating re-integration)
ITERATION=$(echo "$CONTAINER" | jq -r '.current_iteration')
if [ "$ITERATION" -lt 2 ]; then
    echo "❌ R259 VIOLATION: No re-integration detected (iteration=$ITERATION)!"
    exit 1
fi

# Check re_integrated flag is true in latest iteration
RE_INTEGRATED=$(echo "$CONTAINER" | jq -r '.iteration_history[-1].re_integrated')
if [ "$RE_INTEGRATED" != "true" ]; then
    echo "❌ R259 VIOLATION: Re-integration not completed!"
    exit 1
fi

echo "✅ R259 Compliant: Phase re-integration completed (iteration $ITERATION)"
```

## Error Messages

### Violation Detection (SF 3.0)
```
❌ R259 VIOLATION: Direct transition from FIX_PHASE_UPSTREAM_BUGS to COMPLETE_PHASE
Required: FIX_PHASE_UPSTREAM_BUGS → START_PHASE_ITERATION → INTEGRATE_PHASE_WAVES
Found: FIX_PHASE_UPSTREAM_BUGS → COMPLETE_PHASE
Action: Follow SF 3.0 iteration container pattern
```

### Missing Re-Integration
```
❌ R259 VIOLATION: No phase re-integration after fixes
Phase: 3
Container Iteration: 1 (should be 2+)
Re-integrated Flag: false
Action: Transition to START_PHASE_ITERATION to re-integrate
```

## Grading Impact

- **Violation Penalty**: -50% on phase completion score
- **Missing Branch**: Automatic phase failure
- **Skipped State**: -100% on state transition compliance

## Related Rules

- **R256**: Mandatory Phase Assessment Gate
- **R257**: Mandatory Phase Assessment Report
- **R258**: Mandatory Wave Review Report  
- **R009**: Integration Branch Creation
- **R233**: All States Require Immediate Action

## Examples

### Correct Implementation (SF 3.0)
```yaml
# State transitions for phase integration fixes
transitions:
  - state: REVIEW_PHASE_ARCHITECTURE
    timestamp: "2025-08-27T14:00:00Z"

  - state: CREATE_PHASE_FIX_PLAN
    timestamp: "2025-08-27T14:05:00Z"
    reason: "Bugs found in phase integration"

  - state: FIX_PHASE_UPSTREAM_BUGS
    timestamp: "2025-08-27T14:15:00Z"
    action: "Fixing bugs in upstream branches"

  - state: START_PHASE_ITERATION  # ✅ CORRECT - R259 enforcement
    timestamp: "2025-08-27T14:25:00Z"
    iteration: 2

  - state: INTEGRATE_PHASE_WAVES  # ✅ CORRECT - Re-integration
    timestamp: "2025-08-27T14:26:00Z"
    action: "Re-integrating all waves with fixes"

  - state: REVIEW_PHASE_INTEGRATION
    timestamp: "2025-08-27T14:30:00Z"
    action: "Code review of re-integrated phase"
```

### Violation Example (SF 3.0)
```yaml
# INCORRECT - Skips re-integration
transitions:
  - state: FIX_PHASE_UPSTREAM_BUGS
    timestamp: "2025-08-27T14:15:00Z"

  - state: COMPLETE_PHASE  # ❌ VIOLATION
    timestamp: "2025-08-27T14:25:00Z"
    # Missing START_PHASE_ITERATION → INTEGRATE_PHASE_WAVES!
```

## Enforcement

This rule is enforced at:
1. **State Machine Level**: Transition validation
2. **Orchestrator Logic**: State transition decisions
3. **Grading System**: Compliance checks
4. **Architect Review**: Requires integration branch

## Notes

- Similar to wave-level integration but for entire phase
- Ensures all cross-wave issues are properly integrated
- Provides clean baseline for reassessment
- Critical for maintaining phase integrity

## SF 3.0 Iteration Container Architecture

This rule operates within the **PHASE-level iteration container** (defined in `integration-containers.json`):

### Phase Container Convergence Pattern

**Iteration Loop**:
```
PHASE ITERATION 1:
  START_PHASE_ITERATION → INTEGRATE_PHASE_WAVES → REVIEW_PHASE_INTEGRATION →
    REVIEW_PHASE_ARCHITECTURE → (bugs found) → CREATE_PHASE_FIX_PLAN →
    FIX_PHASE_UPSTREAM_BUGS
                    ↓
PHASE ITERATION 2 (R259 ENFORCEMENT POINT):
  START_PHASE_ITERATION → INTEGRATE_PHASE_WAVES → REVIEW_PHASE_INTEGRATION →
        ↑                         ↑
        └─────────────────────────┘
  R259 ENSURES THIS RE-INTEGRATE_WAVE_EFFORTS HAPPENS

  REVIEW_PHASE_ARCHITECTURE → (no bugs) → COMPLETE_PHASE ✅
```

**Critical R259 Enforcement**:
After `FIX_PHASE_UPSTREAM_BUGS`, orchestrator MUST NOT skip directly to review. The mandatory_sequence "phase_iteration_cycle" in `software-factory-3.0-state-machine.json` enforces the full re-integration path.

### Container State Management (SF 3.0)

The phase re-integration after fixes is tracked in `integration-containers.json`:

```json
{
  "phase_containers": [{
    "container_id": "phase-2-integration",
    "phase": 2,
    "current_iteration": 2,
    "convergence_status": "INTEGRATING",
    "integration_branch": "idpbuilder-oci-mgmt/phase2-integration",
    "iteration_history": [
      {
        "iteration": 1,
        "status": "FAILED",
        "bugs_found": 8,
        "review_completed_at": "2025-10-12T14:30:00Z"
      },
      {
        "iteration": 2,
        "status": "IN_PROGRESS",
        "started_at": "2025-10-12T15:00:00Z",
        "fixes_applied": true,
        "re_integrated": true
      }
    ]
  }]
}
```

**R259 Impact on Container**:
- Sets `re_integrated: true` when `INTEGRATE_PHASE_WAVES` completes
- Without re-integration, `re_integrated` stays `false`
- Convergence detection requires `re_integrated: true`
- Skipping integration would cause container to report invalid state

### Validation Implementation (SF 3.0)

**State Machine Enforcement**:
```json
{
  "FIX_PHASE_UPSTREAM_BUGS": {
    "allowed_transitions": ["START_PHASE_ITERATION", "ERROR_RECOVERY"]
  },
  "START_PHASE_ITERATION": {
    "allowed_transitions": ["INTEGRATE_PHASE_WAVES", "ERROR_RECOVERY"]
  },
  "INTEGRATE_PHASE_WAVES": {
    "allowed_transitions": [
      "REVIEW_PHASE_INTEGRATION",
      "CASCADE_REINTEGRATION",
      "ERROR_RECOVERY"
    ]
  }
}
```

**R259 Validation Check**:
```bash
# After FIX_PHASE_UPSTREAM_BUGS completes
CURRENT_STATE=$(jq -r '.state_machine.current_state' orchestrator-state-v3.json)
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)

if [ "$CURRENT_STATE" = "FIX_PHASE_UPSTREAM_BUGS" ]; then
    # Check next state MUST be START_PHASE_ITERATION
    NEXT_STATE=$(determine_next_state)

    if [ "$NEXT_STATE" != "START_PHASE_ITERATION" ]; then
        echo "❌ R259 VIOLATION: After phase fixes, must start new iteration"
        echo "   Required: FIX_PHASE_UPSTREAM_BUGS → START_PHASE_ITERATION"
        echo "   Found: FIX_PHASE_UPSTREAM_BUGS → $NEXT_STATE"
        exit 259
    fi

    # Check that iteration will include INTEGRATE_PHASE_WAVES
    # (enforced by mandatory_sequence in state machine)
fi
```

### R259 in SF 3.0 vs SF 2.0 Differences

| Aspect | SF 2.0 | SF 3.0 |
|--------|--------|--------|
| **Fix State** | ERROR_RECOVERY (generic) | FIX_PHASE_UPSTREAM_BUGS (dedicated) |
| **Re-integration State** | INTEGRATE_PHASE_WAVES (standalone) | START_PHASE_ITERATION → INTEGRATE_PHASE_WAVES |
| **Container** | No formal container | Phase-level iteration container |
| **Iteration Tracking** | Manual restart | Automatic via iteration counter |
| **State Tracking** | State file only | integration-containers.json |
| **Convergence Detection** | Manual assessment | Container-based detection |
| **Integration Branch** | New branch per attempt | Single branch, reset per iteration |
| **Sequence Enforcement** | Optional (state machine) | Mandatory via mandatory_sequences |

### Why R259 Matters More in SF 3.0

**SF 2.0**: Skipping integration meant manual intervention required
**SF 3.0**: Skipping integration breaks the entire iteration container:
- Container cannot detect convergence
- State machine validation fails
- Integration branch becomes stale
- Downstream states receive invalid input

**Enforcement is STRONGER in SF 3.0** due to:
1. `mandatory_sequences` in state machine
2. State Manager validation (R517)
3. Integration container tracking
4. Atomic state updates (R288)