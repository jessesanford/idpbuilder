# ⚠️ DEPRECATED - Subsumed by R288
This rule has been consolidated into R288-state-file-update-and-commit-protocol.md
Please refer to R288 for current state file update and commit requirements.

---

# 🔴🔴🔴 RULE R252 - MANDATORY STATE FILE UPDATES AT EVERY TRANSITION 🔴🔴🔴

## ⚠️⚠️⚠️ SUPREME LAW: NO STATE TRANSITION WITHOUT STATE FILE UPDATE ⚠️⚠️⚠️

## THE ABSOLUTE RULE:

**EVERY state transition MUST update orchestrator-state-v3.json IMMEDIATELY**

No exceptions. No deferrals. No "I'll update it later."

## CRITICAL COMPANION RULES:
**R253 (SUPREME LAW)**: After updating the state file, you MUST immediately commit and push EVERY SINGLE EDIT. R252 defines WHAT to update, R253 enforces WHEN to save it.
**R281 (SUPREME LAW #7)**: When creating the INITIAL state file, it MUST contain ALL phases, waves, and efforts from the implementation plan. See R281 for complete requirements.

## MANDATORY UPDATE SEQUENCE:

```bash
perform_state_transition() {
    local OLD_STATE="$1"
    local NEW_STATE="$2"
    local REASON="$3"
    
    # Step 1: Validate transition is allowed
    validate_state_transition "$OLD_STATE" "$NEW_STATE"
    
    # Step 2: 🔴 UPDATE STATE FILE IMMEDIATELY 🔴
    update_orchestrator_state_file "$OLD_STATE" "$NEW_STATE" "$REASON"
    
    # Step 3: Re-read rules (R217)
    reload_rules_for_state "$NEW_STATE"
    
    # Step 4: Continue with new state work
}
```

## REQUIRED STATE FILE FIELDS FOR EVERY TRANSITION:

```yaml
state_machine:
  current_state: "NEW_STATE"
  previous_state: "OLD_STATE"
  transition_time: "2025-08-25T12:00:00Z"
  transition_reason: "Clear explanation of why transition occurred"
  rules_reacknowledged: true  # Must be set after R217 compliance
```

## STATE-SPECIFIC REQUIRED UPDATES:

### WAVE_COMPLETE → Must Add:
```yaml
waves_completed:
  phase1:
    wave1:
      completed_at: "2025-08-25T12:00:00Z"
      efforts_count: 5
      integration_branch: "tmc-workspace/phase1/wave1-integration"
      status: "COMPLETE"
      all_reviews_passed: true
      size_compliant: true
      tests_passing: true
    wave2:
      status: "NOT_STARTED"
```

### INTEGRATE_WAVE_EFFORTS → Must Add:
```yaml
current_integration:
  phase: 1
  wave: 1
  integration_branch: "tmc-workspace/phase1/wave1-integration"
  integration_directory: "integrations/phase1/wave1/integration-workspace"
  started_at: "2025-08-25T12:00:00Z"
  efforts_to_merge:
    - "core-api-types"
    - "webhook-framework"
    - "controller-base"
```

### ERROR_RECOVERY → Must Add:
```yaml
error_context:
  error_type: "INTEGRATE_WAVE_EFFORTS_CONFLICT"
  error_message: "Merge conflict in pkg/controller/webhook.go"
  error_time: "2025-08-25T12:00:00Z"
  recovery_attempts: 1
  from_state: "INTEGRATE_WAVE_EFFORTS"
```

### PROJECT_DONE → Must Add:
```yaml
phase_completion:
  phase1:
    completed_at: "2025-08-25T12:00:00Z"
    waves_completed: [1, 2, 3]
    total_efforts: 15
    integration_branch: "tmc-workspace/phase1-integration"
    metrics:
      total_lines: 3567
      test_coverage: "82%"
      review_iterations: 2.3
```

### SPAWN_SW_ENGINEERS → Must Add:
```yaml
agents_spawned:
  - agent_type: "sw-engineer"
    effort: "core-api-types"
    spawn_time: "2025-08-25T12:00:00Z"
    spawn_id: "core-api-sw-eng-1234567890"
    working_directory: "efforts/phase1/wave1/core-api-types"
    task: "Implement API types"
```

### MONITOR → Must Update:
```yaml
monitoring_status:
  last_check: "2025-08-25T12:00:00Z"
  agents_active: 3
  agents_completed: 2
  blocked_agents: 0
  next_check_due: "2025-08-25T12:05:00Z"
```

## VALIDATION CHECKLIST:

Before ANY state transition is complete, verify:
- [ ] State file updated with new current_state
- [ ] Previous state captured in previous_state
- [ ] Transition time recorded with ISO timestamp
- [ ] Transition reason clearly documented
- [ ] State-specific fields added/updated
- [ ] File committed to git
- [ ] R217 rules reloaded and acknowledged

## ENFORCEMENT FUNCTIONS:

```bash
# Function to update state file (MANDATORY)
update_orchestrator_state() {
    local NEW_STATE="$1"
    local REASON="$2"
    
    # Get current state before update
    local OLD_STATE=$(jq '.state_machine.current_state' orchestrator-state-v3.json)
    
    # Update core state machine fields
    jq ".state_machine.current_state = \"$NEW_STATE\"" orchestrator-state-v3.json
    jq ".state_machine.previous_state = \"$OLD_STATE\"" orchestrator-state-v3.json
    jq ".state_machine.transition_time = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state-v3.json
    jq ".state_machine.transition_reason = \"$REASON\"" orchestrator-state-v3.json
    jq ".state_machine.rules_reacknowledged = false" orchestrator-state-v3.json
    
    echo "✅ State file updated: $OLD_STATE → $NEW_STATE"
}

# Function to mark wave as complete (MANDATORY for WAVE_COMPLETE)
mark_wave_complete() {
    local PHASE="$1"
    local WAVE="$2"
    
    # Create wave completion record
    jq ".waves_completed.phase${PHASE}.wave${WAVE}.completed_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state-v3.json
    jq ".waves_completed.phase${PHASE}.wave${WAVE}.status = \"COMPLETE\"" orchestrator-state-v3.json
    
    # Count efforts
    local EFFORT_COUNT=$(ls -d efforts/phase${PHASE}/wave${WAVE}/*/ 2>/dev/null | wc -l)
    jq ".waves_completed.phase${PHASE}.wave${WAVE}.efforts_count = $EFFORT_COUNT" orchestrator-state-v3.json
    
    # Record integration branch
    source utilities/branch-naming-helpers.sh
    local INTEGRATE_WAVE_EFFORTS_BRANCH=$(get_wave_integration_branch_name "$PHASE" "$WAVE")
    jq ".waves_completed.phase${PHASE}.wave${WAVE}.integration_branch = \"$INTEGRATE_WAVE_EFFORTS_BRANCH\"" orchestrator-state-v3.json
    
    echo "✅ Wave marked complete: Phase $PHASE, Wave $WAVE"
}

# Function to verify state file was updated
verify_state_file_updated() {
    local EXPECTED_STATE="$1"
    
    local CURRENT=$(jq '.state_machine.current_state' orchestrator-state-v3.json)
    local TIMESTAMP=$(jq '.state_machine.transition_time' orchestrator-state-v3.json)
    
    if [ "$CURRENT" != "$EXPECTED_STATE" ]; then
        echo "❌ ERROR: State file not updated! Expected: $EXPECTED_STATE, Found: $CURRENT"
        exit 1
    fi
    
    # Check timestamp is recent (within last 60 seconds)
    local NOW=$(date +%s)
    local TRANS_TIME=$(date -d "$TIMESTAMP" +%s 2>/dev/null || echo 0)
    local DIFF=$((NOW - TRANS_TIME))
    
    if [ $DIFF -gt 60 ]; then
        echo "⚠️ WARNING: State file timestamp is stale (${DIFF}s old)"
    fi
    
    echo "✅ State file verified: $CURRENT at $TIMESTAMP"
}
```

## COMMON VIOLATIONS TO AVOID:

### ❌ WRONG: Transitioning without update
```bash
# NO! Missing state file update
echo "Moving to INTEGRATE_WAVE_EFFORTS state"
cd integrations/
# Start integration work...
```

### ✅ CORRECT: Update THEN transition
```bash
# YES! Update state file first
update_orchestrator_state "INTEGRATE_WAVE_EFFORTS" "All wave efforts complete"
mark_wave_complete "1" "1"
echo "State updated, now transitioning to INTEGRATE_WAVE_EFFORTS"
cd integrations/
```

### ❌ WRONG: Updating state file later
```bash
# NO! Deferred update
echo "Starting integration work"
do_integration_stuff
echo "Oh, I should update the state file now..."  # TOO LATE!
```

### ✅ CORRECT: Update immediately
```bash
# YES! Immediate update
update_orchestrator_state "INTEGRATE_WAVE_EFFORTS" "Wave complete"
echo "State updated, proceeding with integration"
do_integration_stuff
```

## GRADING IMPACT:

**AUTOMATIC FAILURE** if:
- State transition occurs without state file update
- Wave marked complete without waves_completed entry
- Phase marked complete without phase_completion entry
- Timestamp is missing or stale (>60s old)
- Previous state not captured

## THE GOLDEN RULE:

**No state transition is complete until orchestrator-state-v3.json reflects it.**

The state file is the single source of truth. If it's not in the state file, it didn't happen.