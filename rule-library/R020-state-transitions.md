# Rule R020: State Transitions Protocol

## Rule Statement
Agents MUST follow the SOFTWARE-FACTORY-STATE-MACHINE.md for all state transitions. Transitions must be atomic, validated, logged, and immediately followed by work continuation. NEVER stop after a transition - states are waypoints, not destinations.

## Criticality Level
**BLOCKING** - Invalid transitions break the entire factory workflow

## Enforcement Mechanism
- **Technical**: Validate transitions against state machine
- **Behavioral**: Atomic transition with immediate continuation
- **Grading**: -50% for invalid transitions, -100% for stopping after transition

## Core Principle

```
State Transition = Validate → Update → Commit → Continue IMMEDIATELY
NEVER transition to undefined states
NEVER stop after transitioning
ALWAYS continue with state-specific work
States are VERBS - they mean DO THE WORK NOW!
```

## Detailed Requirements

### State Transition Protocol

1. **Pre-Transition Validation**
   ```bash
   # Read current state
   current_state=$(yq '.current_state' orchestrator-state.json)
   
   # Validate next state exists
   grep -q "STATE: ${next_state}" SOFTWARE-FACTORY-STATE-MACHINE.md || {
       echo "❌ INVALID STATE: ${next_state}"
       exit 1
   }
   
   # Verify transition is allowed
   validate_transition "${current_state}" "${next_state}"
   ```

2. **Atomic Transition**
   ```bash
   # Update state file
   yq -i '.current_state = "'${next_state}'"' orchestrator-state.json
   yq -i '.transition_timestamp = "'$(date -Iseconds)'"' orchestrator-state.json
   
   # Commit immediately
   git add orchestrator-state.json
   git commit -m "state: transition from ${current_state} to ${next_state}"
   git push
   ```

3. **Immediate Continuation**
   ```markdown
   ✅ Transitioned to ${next_state}
   📋 Loading ${next_state} rules...
   🚀 Starting ${next_state} work immediately...
   ```

### Valid State Sequences

**Orchestrator States:**
- INIT → PLANNING
- PLANNING → SETUP_EFFORT_INFRASTRUCTURE  
- SETUP_EFFORT_INFRASTRUCTURE → SPAWN_AGENTS
- SPAWN_AGENTS → MONITOR
- MONITOR → WAVE_COMPLETE or ERROR_RECOVERY
- WAVE_COMPLETE → INTEGRATION
- INTEGRATION → PHASE_COMPLETE or SPAWN_AGENTS

### Forbidden Practices

- ❌ NEVER create custom states
- ❌ NEVER skip states in sequence
- ❌ NEVER stop after transition
- ❌ NEVER transition without validation
- ❌ NEVER leave state file uncommitted

### State Interpretation

States are ACTIONS, not resting points:
- **PLANNING** = Create plans NOW
- **SPAWN_AGENTS** = Spawn agents NOW
- **MONITOR** = Monitor progress NOW
- **INTEGRATION** = Perform integration NOW

## Relationship to Other Rules
- **R206**: State machine transition validation
- **R233**: All states immediate action
- **R234**: Mandatory state traversal
- **R288**: State file update and commit
- **R322**: Mandatory stop before transitions

## Implementation Notes
- Transitions must complete within 30 seconds
- State files must be pushed within 60 seconds
- Work must begin within 10 seconds of transition
- Invalid transitions trigger immediate ERROR_RECOVERY