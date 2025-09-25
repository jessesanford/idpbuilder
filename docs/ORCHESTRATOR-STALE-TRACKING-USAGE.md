# ORCHESTRATOR GUIDE: Using Stale Integration Tracking

## Quick Start

The stale integration tracking mechanism provides comprehensive tracking of when and why integrations become stale, enabling proper R327 cascade enforcement.

## Key Commands for Orchestrator

### 1. Check All Integrations for Staleness
```bash
# Run comprehensive check
./utilities/stale-integration-manager.sh check

# Check specific integration
./utilities/stale-integration-manager.sh check wave phase1-wave2-integration
```

### 2. Track When Fixes Are Applied
```bash
# When SW Engineer reports a fix
./utilities/stale-integration-manager.sh track-fix \
    "FIX-001" \
    "abc123def" \
    "phase1/wave2/effort3" \
    "auth-module" \
    "build_fix" \
    "Fixed missing import causing build failure"
```

### 3. Mark Integration as Recreated
```bash
# After successfully recreating an integration
./utilities/stale-integration-manager.sh mark-recreated wave phase1-wave2-integration
```

### 4. Generate Staleness Report
```bash
# Create comprehensive report
./utilities/stale-integration-manager.sh report
# Creates: STALENESS-REPORT.md
```

## Integration Points in State Machine

### State: MONITOR_FIXES
```bash
# When monitoring fixes, track them
if [[ "$fix_completed" == "true" ]]; then
    ./utilities/stale-integration-manager.sh track-fix \
        "$fix_id" "$commit" "$branch" "$effort" "$type" "$description"
    
    # Check for staleness impact
    ./utilities/stale-integration-manager.sh check
fi
```

### State: INTEGRATION
```bash
# Before starting integration, check freshness
if ! ./utilities/stale-integration-manager.sh check wave; then
    echo "❌ Wave integration is STALE - must recreate"
    # Transition to recreation flow
fi
```

### State: PHASE_INTEGRATION
```bash
# Check all wave integrations are fresh
for wave in $(get_completed_waves); do
    if ! ./utilities/stale-integration-manager.sh check wave "$wave-integration"; then
        echo "❌ Wave $wave is stale - cascade recreation required"
        exit 1
    fi
done
```

## Reading Staleness Data

### Check if Integration is Stale
```bash
is_stale=$(jq -r '.current_wave_integration.is_stale' orchestrator-state.json)
if [[ "$is_stale" == "true" ]]; then
    stale_reason=$(jq -r '.current_wave_integration.staleness_reason' orchestrator-state.json)
    echo "Integration is stale: $stale_reason"
fi
```

### Get Fixes That Made Integration Stale
```bash
stale_fixes=$(jq -r '.current_wave_integration.stale_due_to_fixes[]' orchestrator-state.json)
for fix in $stale_fixes; do
    echo "Fix $fix made this integration stale"
done
```

### Check Cascade Requirements
```bash
pending_cascades=$(jq -r '
    .stale_integration_tracking.staleness_cascade[] |
    select(.cascade_status != "completed") |
    "Cascade from " + .trigger.branch + " is " + .cascade_status
' orchestrator-state.json)
```

### Find Which Efforts Received Fixes
```bash
efforts_with_fixes=$(jq -r '
    .stale_integration_tracking.fix_tracking.fixes_applied[] |
    .effort_name + ": " + .description
' orchestrator-state.json | sort -u)
```

## Enforcement Workflow

### 1. Continuous Monitoring
```bash
# In MONITOR state, periodically check
while in_monitor_state; do
    ./utilities/stale-integration-manager.sh check
    
    if [[ $? -ne 0 ]]; then
        echo "⚠️ Stale integrations detected!"
        transition_to_state "INTEGRATION"  # Start recreation
    fi
    
    sleep 300  # Check every 5 minutes
done
```

### 2. Pre-Integration Validation
```bash
# Before ANY integration operation
validate_before_integration() {
    local integration_type=$1
    
    # Check freshness
    if ! ./utilities/stale-integration-manager.sh check $integration_type; then
        echo "❌ R327 VIOLATION: Attempting to use stale integration"
        echo "Must recreate $integration_type integration first"
        return 1
    fi
    
    return 0
}
```

### 3. Cascade Enforcement
```bash
# When recreating due to staleness
handle_cascade_recreation() {
    # Get cascade requirements
    local cascades=$(jq -r '
        .stale_integration_tracking.staleness_cascade[] |
        select(.cascade_status == "pending") |
        .cascade_sequence[].integration
    ' orchestrator-state.json)
    
    for integration in $cascades; do
        echo "Recreating: $integration (R327 cascade)"
        recreate_integration "$integration"
        ./utilities/stale-integration-manager.sh mark-recreated \
            "$(get_type $integration)" "$integration"
    done
}
```

## State File Fields Reference

### Integration Objects
```json
"current_wave_integration": {
    "is_stale": false,              // Quick boolean check
    "staleness_reason": null,       // Human-readable reason
    "stale_since": null,           // When it became stale
    "stale_due_to_fixes": [],      // List of fix IDs
    "last_freshness_check": "..."   // Last validation time
}
```

### Tracking Structure
```json
"stale_integration_tracking": {
    "stale_integrations": [],       // List of all stale integrations
    "staleness_cascade": [],        // Required cascade recreations
    "fix_tracking": {               // Comprehensive fix tracking
        "fixes_applied": [],        // All fixes with details
        "fixes_pending_integration": []  // Fixes not yet integrated
    },
    "validation_history": []        // Audit trail of checks
}
```

## Common Queries

### List All Stale Integrations
```bash
jq -r '.stale_integration_tracking.stale_integrations[] |
       select(.recreation_completed == false) |
       .integration_id' orchestrator-state.json
```

### Get Fix Details for Effort
```bash
effort="auth-module"
jq -r --arg e "$effort" '
    .stale_integration_tracking.fix_tracking.fixes_applied[] |
    select(.effort_name == $e) |
    "Fix " + .fix_id + ": " + .description
' orchestrator-state.json
```

### Check Integration Status
```bash
jq -r '.stale_integration_tracking.fix_tracking.fixes_applied[] |
       "Fix " + .fix_id + " integrated: " +
       "Wave=" + (.integrated_into.wave|tostring) +
       " Phase=" + (.integrated_into.phase|tostring) +
       " Project=" + (.integrated_into.project|tostring)
' orchestrator-state.json
```

## Benefits for Orchestrator

1. **Automated Detection**: No manual checking needed
2. **Clear Accountability**: Know exactly which fix caused staleness
3. **Cascade Management**: Automatic tracking of required recreations
4. **Audit Trail**: Complete history for grading/debugging
5. **Recovery Path**: Clear steps to resolve staleness
6. **R327 Compliance**: Automatic enforcement of cascade rules

## Example: Complete Flow

```bash
# 1. SW Engineer fixes an issue
echo "SW Engineer applied fix to effort3"

# 2. Track the fix
./utilities/stale-integration-manager.sh track-fix \
    "FIX-123" "def456" "phase1/wave2/effort3" \
    "payment-module" "build_fix" "Fixed null pointer"

# 3. Check for staleness
./utilities/stale-integration-manager.sh check
# Output: Wave integration is STALE!

# 4. Generate report
./utilities/stale-integration-manager.sh report

# 5. Recreate integrations (cascade)
recreate_wave_integration
./utilities/stale-integration-manager.sh mark-recreated wave phase1-wave2-integration

recreate_phase_integration  
./utilities/stale-integration-manager.sh mark-recreated phase phase1-integration

recreate_project_integration
./utilities/stale-integration-manager.sh mark-recreated project project-integration

# 6. Verify all fresh
./utilities/stale-integration-manager.sh check
# Output: All integrations are fresh
```

This tracking mechanism ensures complete R327 compliance and provides full visibility into integration staleness.