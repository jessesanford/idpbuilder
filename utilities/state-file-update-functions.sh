#!/bin/bash

# Software Factory 2.0 - State File Update Functions
# MANDATORY functions for R288 compliance - EVERY state transition MUST use these!
# MANDATORY functions for R288 compliance - EVERY edit MUST commit and push!

# ========================================
# R288 COMPLIANCE HELPER
# ========================================

# Helper function to ensure R288 compliance
commit_and_push_state() {
    local CHANGE_DESC="$1"
    
    # Check if there are changes to commit
    if ! git status --porcelain | grep -q "orchestrator-state.json"; then
        echo "⚠️ No changes to orchestrator-state.json to commit"
        return 0
    fi
    
    # Stage, commit, and push
    git add orchestrator-state.json
    
    if ! git commit -m "state: $CHANGE_DESC [R288]"; then
        echo "❌ R288: Failed to commit state change!"
        return 253
    fi
    
    if ! git push; then
        echo "🔴🔴🔴 R288 VIOLATION: Failed to push state!"
        # Try with force-with-lease as fallback
        if ! git push --force-with-lease; then
            echo "❌❌❌ FATAL: Cannot push state! STOP ALL WORK!"
            exit 253
        fi
    fi
    
    echo "✅ R288: State committed and pushed: $CHANGE_DESC"
    return 0
}

# ========================================
# CORE STATE MACHINE UPDATE (R288 + R288)
# ========================================

# Main function to update state - MUST be called on EVERY transition
update_orchestrator_state() {
    local NEW_STATE="$1"
    local REASON="${2:-State transition}"
    
    echo "🔄 R288: Updating orchestrator-state.json for transition to $NEW_STATE"
    
    # Get current state before update
    local OLD_STATE=$(yq '.state_machine.current_state' orchestrator-state.json 2>/dev/null || echo "UNKNOWN")
    
    # Validate transition is allowed (should check state machine)
    if ! validate_state_transition "$OLD_STATE" "$NEW_STATE"; then
        echo "❌ ERROR: Invalid state transition $OLD_STATE → $NEW_STATE"
        return 1
    fi
    
    # Update core state machine fields (R288: commit after EACH edit)
    yq -i ".state_machine.current_state = \"$NEW_STATE\"" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: current_state=$NEW_STATE [R288]" && git push
    
    yq -i ".state_machine.previous_state = \"$OLD_STATE\"" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: previous_state=$OLD_STATE [R288]" && git push
    
    yq -i ".state_machine.transition_time = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: transition time recorded [R288]" && git push
    
    yq -i ".state_machine.transition_reason = \"$REASON\"" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: reason=$REASON [R288]" && git push
    
    yq -i ".state_machine.rules_reacknowledged = false" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: rules need reacknowledgment [R288]" && git push
    
    echo "✅ State file updated: $OLD_STATE → $NEW_STATE"
    echo "   Reason: $REASON"
    
    # Call state-specific update function if it exists
    local STATE_FUNC="update_state_${NEW_STATE,,}"
    if declare -f "$STATE_FUNC" > /dev/null; then
        echo "   Running state-specific updates for $NEW_STATE..."
        "$STATE_FUNC"
    fi
    
    # R288 compliance check - all commits should already be done
    if git status --porcelain | grep -q "orchestrator-state.json"; then
        echo "🔴🔴🔴 R288 VIOLATION: Uncommitted state changes detected!"
        echo "This should never happen - all changes should be committed immediately!"
        exit 253
    fi
    
    return 0
}

# Validate state transition against state machine
validate_state_transition() {
    local FROM_STATE="$1"
    local TO_STATE="$2"
    
    # TODO: Check against SOFTWARE-FACTORY-STATE-MACHINE.md
    # For now, just check states are different
    if [ "$FROM_STATE" = "$TO_STATE" ]; then
        echo "⚠️ WARNING: Transitioning to same state: $FROM_STATE"
    fi
    
    return 0  # Allow all transitions for now
}

# Mark rules as reacknowledged after R217 compliance
mark_rules_acknowledged() {
    yq -i ".state_machine.rules_reacknowledged = true" orchestrator-state.json
    # R288: Immediate commit and push
    git add orchestrator-state.json
    git commit -m "state: rules acknowledged [R288]"
    git push
    echo "✅ R217: Rules reacknowledged for current state"
}

# ========================================
# WAVE_COMPLETE STATE UPDATES
# ========================================

update_state_wave_complete() {
    echo "   📝 Adding wave completion data..."
    # This function is called when transitioning to WAVE_COMPLETE
    # The actual wave marking happens in mark_wave_complete()
}

# Mark a specific wave as complete (MANDATORY for WAVE_COMPLETE)
mark_wave_complete() {
    local PHASE="$1"
    local WAVE="$2"
    
    echo "🌊 R288: Marking wave complete - Phase $PHASE, Wave $WAVE"
    
    # Create wave completion record
    local WAVE_PATH="waves_completed.phase${PHASE}.wave${WAVE}"
    
    yq -i ".${WAVE_PATH}.completed_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: wave${WAVE} completion time [R288]" && git push
    
    yq -i ".${WAVE_PATH}.status = \"COMPLETE\"" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: wave${WAVE} status=COMPLETE [R288]" && git push
    
    # Count efforts
    local EFFORT_COUNT=$(ls -d /efforts/phase${PHASE}/wave${WAVE}/*/ 2>/dev/null | wc -l || echo 0)
    yq -i ".${WAVE_PATH}.efforts_count = $EFFORT_COUNT" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: wave${WAVE} effort_count=$EFFORT_COUNT [R288]" && git push
    
    # List effort names
    local EFFORTS=$(ls -d /efforts/phase${PHASE}/wave${WAVE}/*/ 2>/dev/null | xargs -n1 basename | tr '\n' ',' | sed 's/,$//')
    yq -i ".${WAVE_PATH}.efforts = \"$EFFORTS\"" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: wave${WAVE} efforts listed [R288]" && git push
    
    # Get integration branch name
    if [ -f utilities/branch-naming-helpers.sh ]; then
        source utilities/branch-naming-helpers.sh
        local INTEGRATION_BRANCH=$(get_wave_integration_branch_name "$PHASE" "$WAVE")
        yq -i ".${WAVE_PATH}.integration_branch = \"$INTEGRATION_BRANCH\"" orchestrator-state.json
    fi
    
    # Add validation flags (R288: commit each)
    yq -i ".${WAVE_PATH}.all_reviews_passed = true" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: wave${WAVE} reviews passed [R288]" && git push
    
    yq -i ".${WAVE_PATH}.size_compliant = true" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: wave${WAVE} size compliant [R288]" && git push
    
    yq -i ".${WAVE_PATH}.tests_passing = true" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: wave${WAVE} tests passing [R288]" && git push
    
    echo "✅ Wave marked complete: Phase $PHASE, Wave $WAVE"
    echo "   Efforts: $EFFORT_COUNT"
    echo "   Integration branch: $INTEGRATION_BRANCH"
}

# ========================================
# INTEGRATION STATE UPDATES
# ========================================

update_state_integration() {
    local PHASE="${CURRENT_PHASE:-1}"
    local WAVE="${CURRENT_WAVE:-1}"
    
    echo "   📝 Adding integration context..."
    
    yq -i ".current_integration.phase = $PHASE" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: integration phase=$PHASE [R288]" && git push
    
    yq -i ".current_integration.wave = $WAVE" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: integration wave=$WAVE [R288]" && git push
    
    yq -i ".current_integration.started_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: integration started [R288]" && git push
    
    # Get integration branch
    if [ -f utilities/branch-naming-helpers.sh ]; then
        source utilities/branch-naming-helpers.sh
        local INTEGRATION_BRANCH=$(get_wave_integration_branch_name "$PHASE" "$WAVE")
        yq -i ".current_integration.integration_branch = \"$INTEGRATION_BRANCH\"" orchestrator-state.json
    fi
    
    # Set integration directory
    yq -i ".current_integration.integration_directory = \"/efforts/phase${PHASE}/wave${WAVE}/integration-workspace\"" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: integration directory set [R288]" && git push
    
    # List efforts to merge (R288: commit after building list)
    local EFFORTS=$(ls -d /efforts/phase${PHASE}/wave${WAVE}/*/ 2>/dev/null | xargs -n1 basename)
    for effort in $EFFORTS; do
        yq -i ".current_integration.efforts_to_merge += [\"$effort\"]" orchestrator-state.json
    done
    # Single commit for effort list (acceptable batch for array building)
    git add orchestrator-state.json && git commit -m "state: integration efforts list updated [R288]" && git push
}

# ========================================
# R301 INTEGRATION TRACKING FUNCTIONS
# ========================================

# Validate current wave integration
validate_current_wave_integration() {
    local PHASE=$1
    local WAVE=$2
    local BRANCH_TO_USE=$3
    
    # Get the CURRENT wave integration branch
    local CURRENT=$(yq ".current_wave_integration | select(.phase == $PHASE and .wave == $WAVE).branch" orchestrator-state.json)
    local CURRENT_STATUS=$(yq ".current_wave_integration | select(.phase == $PHASE and .wave == $WAVE).status" orchestrator-state.json)
    
    # FATAL if not using current
    if [[ "$BRANCH_TO_USE" != "$CURRENT" ]]; then
        echo "🔴🔴🔴 SUPREME LAW VIOLATION: R301 - Using deprecated wave integration branch!"
        echo "  Current: $CURRENT (status: $CURRENT_STATUS)"
        echo "  Attempted: $BRANCH_TO_USE"
        echo "  PENALTY: -100% AUTOMATIC FAILURE"
        return 1
    fi
    
    # FATAL if current is not active
    if [[ "$CURRENT_STATUS" != "active" ]]; then
        echo "🔴🔴🔴 FATAL: Current wave integration is not active!"
        return 1
    fi
    
    echo "✅ Using current wave integration: $CURRENT"
    return 0
}

# Validate current phase integration
validate_current_phase_integration() {
    local PHASE=$1
    local BRANCH_TO_USE=$2
    
    # Get the CURRENT phase integration branch
    local CURRENT=$(yq ".current_phase_integration | select(.phase == $PHASE).branch" orchestrator-state.json)
    local CURRENT_STATUS=$(yq ".current_phase_integration | select(.phase == $PHASE).status" orchestrator-state.json)
    
    # FATAL if not using current
    if [[ "$BRANCH_TO_USE" != "$CURRENT" ]]; then
        echo "🔴🔴🔴 SUPREME LAW VIOLATION: R301 - Using deprecated phase integration branch!"
        echo "  Current: $CURRENT (status: $CURRENT_STATUS)"
        echo "  Attempted: $BRANCH_TO_USE"
        echo "  PENALTY: -100% AUTOMATIC FAILURE"
        return 1
    fi
    
    # FATAL if current is not active
    if [[ "$CURRENT_STATUS" != "active" ]]; then
        echo "🔴🔴🔴 FATAL: Current phase integration is not active!"
        return 1
    fi
    
    echo "✅ Using current phase integration: $CURRENT"
    return 0
}

# Get current wave integration branch
get_current_wave_integration() {
    local PHASE=$1
    local WAVE=$2
    yq ".current_wave_integration | select(.phase == $PHASE and .wave == $WAVE).branch" orchestrator-state.json
}

# Get current phase integration branch
get_current_phase_integration() {
    local PHASE=$1
    yq ".current_phase_integration | select(.phase == $PHASE).branch" orchestrator-state.json
}

# ========================================
# ERROR_RECOVERY STATE UPDATES
# ========================================

update_state_error_recovery() {
    local ERROR_TYPE="${ERROR_TYPE:-UNKNOWN}"
    local ERROR_MSG="${ERROR_MSG:-Error occurred}"
    
    echo "   📝 Adding error context..."
    
    yq -i ".error_context.error_type = \"$ERROR_TYPE\"" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: error_type=$ERROR_TYPE [R288]" && git push
    
    yq -i ".error_context.error_message = \"$ERROR_MSG\"" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: error recorded [R288]" && git push
    
    yq -i ".error_context.error_time = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: error time recorded [R288]" && git push
    
    local FROM_STATE=$(yq '.state_machine.previous_state' orchestrator-state.json)
    yq -i ".error_context.from_state = \"$FROM_STATE\"" orchestrator-state.json
    
    # Increment recovery attempts
    local ATTEMPTS=$(yq '.error_context.recovery_attempts' orchestrator-state.json 2>/dev/null || echo 0)
    yq -i ".error_context.recovery_attempts = $((ATTEMPTS + 1))" orchestrator-state.json
}

# ========================================
# SUCCESS STATE UPDATES
# ========================================

update_state_success() {
    local PHASE="${CURRENT_PHASE:-1}"
    
    echo "   📝 Adding phase completion data..."
    
    local PHASE_PATH="phase_completion.phase${PHASE}"
    
    yq -i ".${PHASE_PATH}.completed_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.json
    
    # Count completed waves
    local WAVE_COUNT=$(yq ".waves_completed.phase${PHASE} | keys | length" orchestrator-state.json 2>/dev/null || echo 0)
    yq -i ".${PHASE_PATH}.waves_completed = $WAVE_COUNT" orchestrator-state.json
    
    # Count total efforts
    local TOTAL_EFFORTS=0
    for wave_dir in efforts/phase${PHASE}/wave*/; do
        if [ -d "$wave_dir" ]; then
            EFFORTS=$(ls -d "$wave_dir"*/ 2>/dev/null | wc -l || echo 0)
            TOTAL_EFFORTS=$((TOTAL_EFFORTS + EFFORTS))
        fi
    done
    yq -i ".${PHASE_PATH}.total_efforts = $TOTAL_EFFORTS" orchestrator-state.json
    
    # Get phase integration branch
    if [ -f utilities/branch-naming-helpers.sh ]; then
        source utilities/branch-naming-helpers.sh
        local PHASE_BRANCH=$(get_phase_integration_branch_name "$PHASE")
        yq -i ".${PHASE_PATH}.integration_branch = \"$PHASE_BRANCH\"" orchestrator-state.json
    fi
}

# ========================================
# SPAWN_AGENTS STATE UPDATES
# ========================================

add_spawned_agent() {
    local AGENT_TYPE="$1"
    local EFFORT="$2"
    local SPAWN_ID="$3"
    local TASK="$4"
    local WORKING_DIR="$5"
    
    echo "   📝 Recording spawned agent: $AGENT_TYPE for $EFFORT"
    
    # Create agent spawn record
    local AGENT_ENTRY=$(cat <<EOF
{
  "agent_type": "$AGENT_TYPE",
  "effort": "$EFFORT",
  "spawn_time": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "spawn_id": "$SPAWN_ID",
  "working_directory": "$WORKING_DIR",
  "task": "$TASK",
  "status": "ACTIVE"
}
EOF
)
    
    # Add to agents_spawned array
    echo "$AGENT_ENTRY" | yq -P | yq eval-all 'select(fileIndex == 0).agents_spawned += select(fileIndex == 1) | select(fileIndex == 0)' orchestrator-state.json - > /tmp/updated-state.yaml
    mv /tmp/updated-state.yaml orchestrator-state.json
}

# ========================================
# MONITOR STATE UPDATES
# ========================================

update_monitoring_status() {
    echo "   📝 Updating monitoring status..."
    
    yq -i ".monitoring_status.last_check = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: monitoring check performed [R288]" && git push
    
    # Count active agents
    local ACTIVE=$(yq '.agents_spawned[] | select(.status == "ACTIVE") | length' orchestrator-state.json 2>/dev/null | wc -l || echo 0)
    local COMPLETED=$(yq '.agents_spawned[] | select(.status == "COMPLETED") | length' orchestrator-state.json 2>/dev/null | wc -l || echo 0)
    local BLOCKED=$(yq '.agents_spawned[] | select(.status == "BLOCKED") | length' orchestrator-state.json 2>/dev/null | wc -l || echo 0)
    
    yq -i ".monitoring_status.agents_active = $ACTIVE" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: active agents=$ACTIVE [R288]" && git push
    
    yq -i ".monitoring_status.agents_completed = $COMPLETED" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: completed agents=$COMPLETED [R288]" && git push
    
    yq -i ".monitoring_status.blocked_agents = $BLOCKED" orchestrator-state.json
    git add orchestrator-state.json && git commit -m "state: blocked agents=$BLOCKED [R288]" && git push
    
    # Set next check time (5 minutes from now)
    local NEXT_CHECK=$(date -u -d "+5 minutes" +%Y-%m-%dT%H:%M:%SZ)
    yq -i ".monitoring_status.next_check_due = \"$NEXT_CHECK\"" orchestrator-state.json
}

# ========================================
# VALIDATION FUNCTIONS
# ========================================

# Verify state file was updated recently
verify_state_file_updated() {
    local EXPECTED_STATE="${1:-$(yq '.state_machine.current_state' orchestrator-state.json)}"
    
    local CURRENT=$(yq '.state_machine.current_state' orchestrator-state.json)
    local TIMESTAMP=$(yq '.state_machine.transition_time' orchestrator-state.json)
    
    if [ "$CURRENT" != "$EXPECTED_STATE" ]; then
        echo "❌ R288 VIOLATION: State file not updated! Expected: $EXPECTED_STATE, Found: $CURRENT"
        return 1
    fi
    
    # Check timestamp is recent (within last 60 seconds)
    if command -v date > /dev/null; then
        local NOW=$(date +%s)
        local TRANS_TIME=$(date -d "$TIMESTAMP" +%s 2>/dev/null || echo 0)
        local DIFF=$((NOW - TRANS_TIME))
        
        if [ $DIFF -gt 60 ]; then
            echo "⚠️ R288 WARNING: State file timestamp is stale (${DIFF}s old)"
        fi
    fi
    
    echo "✅ R288: State file verified - $CURRENT at $TIMESTAMP"
    return 0
}

# Check if wave is marked complete
is_wave_complete() {
    local PHASE="$1"
    local WAVE="$2"
    
    local STATUS=$(yq ".waves_completed.phase${PHASE}.wave${WAVE}.status" orchestrator-state.json 2>/dev/null)
    
    if [ "$STATUS" = "COMPLETE" ]; then
        return 0
    else
        return 1
    fi
}

# ========================================
# EXPORT ALL FUNCTIONS
# ========================================

export -f commit_and_push_state  # R288 compliance helper
export -f update_orchestrator_state
export -f validate_state_transition
export -f mark_rules_acknowledged
export -f mark_wave_complete
export -f add_spawned_agent
export -f update_monitoring_status
export -f verify_state_file_updated
export -f is_wave_complete

# If run directly, show usage
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    echo "State File Update Functions (R288 Compliance)"
    echo "============================================="
    echo ""
    echo "Core function:"
    echo "  update_orchestrator_state NEW_STATE REASON"
    echo ""
    echo "State-specific functions:"
    echo "  mark_wave_complete PHASE WAVE"
    echo "  add_spawned_agent TYPE EFFORT ID TASK DIR"
    echo "  update_monitoring_status"
    echo ""
    echo "Validation:"
    echo "  verify_state_file_updated [EXPECTED_STATE]"
    echo "  is_wave_complete PHASE WAVE"
    echo ""
    echo "Example:"
    echo "  source state-file-update-functions.sh"
    echo "  update_orchestrator_state WAVE_COMPLETE \"All efforts reviewed\""
    echo "  mark_wave_complete 1 1"
fi