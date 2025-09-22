# State File Update Validation Checklist (R288)

## 🔴 CRITICAL: This checklist MUST be followed for EVERY state transition!

### Pre-Transition Checklist
- [ ] Current state identified from orchestrator-state.json
- [ ] Target state validated against state machine
- [ ] Transition reason clearly defined
- [ ] State update function available in utilities/state-file-update-functions.sh

### During Transition - Core Updates (ALL STATES)
- [ ] `update_orchestrator_state()` called IMMEDIATELY
- [ ] current_state updated to new state
- [ ] previous_state captured
- [ ] transition_time set to current UTC timestamp
- [ ] transition_reason documented
- [ ] rules_reacknowledged set to false (will be true after R217)

### State-Specific Required Updates

#### WAVE_COMPLETE
- [ ] `mark_wave_complete()` called with phase and wave numbers
- [ ] waves_completed.phase{X}.wave{Y}.completed_at timestamp added
- [ ] waves_completed.phase{X}.wave{Y}.status set to "COMPLETE"
- [ ] waves_completed.phase{X}.wave{Y}.efforts_count recorded
- [ ] waves_completed.phase{X}.wave{Y}.efforts list populated
- [ ] waves_completed.phase{X}.wave{Y}.integration_branch recorded
- [ ] waves_completed.phase{X}.wave{Y}.all_reviews_passed = true
- [ ] waves_completed.phase{X}.wave{Y}.size_compliant = true
- [ ] waves_completed.phase{X}.wave{Y}.tests_passing = true

#### INTEGRATION
- [ ] current_integration.phase set
- [ ] current_integration.wave set
- [ ] current_integration.started_at timestamp
- [ ] current_integration.integration_branch recorded
- [ ] current_integration.integration_directory path set
- [ ] current_integration.efforts_to_merge array populated

#### ERROR_RECOVERY
- [ ] error_context.error_type categorized
- [ ] error_context.error_message captured
- [ ] error_context.error_time timestamp
- [ ] error_context.from_state recorded
- [ ] error_context.recovery_attempts incremented

#### SUCCESS (Phase Complete)
- [ ] phase_completion.phase{X}.completed_at timestamp
- [ ] phase_completion.phase{X}.waves_completed count
- [ ] phase_completion.phase{X}.total_efforts count
- [ ] phase_completion.phase{X}.integration_branch recorded
- [ ] phase_completion.phase{X}.metrics populated

#### SPAWN_AGENTS
- [ ] For EACH spawned agent:
  - [ ] agents_spawned entry added
  - [ ] agent_type recorded
  - [ ] effort name recorded
  - [ ] spawn_time timestamp
  - [ ] spawn_id unique identifier
  - [ ] working_directory path
  - [ ] task description
  - [ ] status set to "ACTIVE"

#### MONITOR
- [ ] monitoring_status.last_check timestamp
- [ ] monitoring_status.agents_active count
- [ ] monitoring_status.agents_completed count
- [ ] monitoring_status.blocked_agents count
- [ ] monitoring_status.next_check_due timestamp

### Post-Transition Validation
- [ ] `verify_state_file_updated()` returns success
- [ ] State file committed to git
- [ ] Timestamp is current (within 60 seconds)
- [ ] R217 rule reloading completed
- [ ] rules_reacknowledged set to true after R217

### Common Violations to Check
- [ ] ❌ State changed but file not updated
- [ ] ❌ Wave marked complete but not in state file
- [ ] ❌ Integration started but no current_integration entry
- [ ] ❌ Agents spawned but not recorded
- [ ] ❌ Timestamp is stale or missing
- [ ] ❌ Previous state not captured
- [ ] ❌ Transition reason not documented

## Example Usage

```bash
# Source the update functions
source utilities/state-file-update-functions.sh

# Example: Transitioning to WAVE_COMPLETE
update_orchestrator_state "WAVE_COMPLETE" "All Wave 1 efforts reviewed and passed"
mark_wave_complete "1" "1"
verify_state_file_updated "WAVE_COMPLETE"

# Example: Transitioning to INTEGRATION
update_orchestrator_state "INTEGRATION" "Wave 1 complete, starting integration"
# State-specific updates handled automatically by update_state_integration()
verify_state_file_updated "INTEGRATION"

# Example: Spawning an agent
add_spawned_agent "sw-engineer" "core-api-types" "core-api-sw-123456" \
    "Implement API types" "efforts/phase1/wave1/core-api-types"
```

## Verification Commands

```bash
# Check current state
jq '.state_machine.current_state' orchestrator-state.json

# Check if wave is marked complete
jq '.waves_completed.phase1.wave1.status' orchestrator-state.json

# Check transition timestamp freshness
jq '.state_machine.transition_time' orchestrator-state.json

# List all completed waves
jq '.waves_completed' orchestrator-state.json

# Check active agents
jq '.agents_spawned[] | select(.status == "ACTIVE")' orchestrator-state.json
```

## R288 Enforcement

**REMEMBER: The state file is the single source of truth!**

If it's not in orchestrator-state.json, it didn't happen.

**AUTOMATIC FAILURE** if state transitions occur without proper state file updates.