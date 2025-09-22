# 🔴🔴🔴 RULE R346: STATE METADATA SYNCHRONIZATION (SUPREME LAW)

## Summary
ALL state metadata updates MUST be atomic and synchronized across related fields to prevent inconsistent state. Integration status, stale tracking, and completion flags MUST always be consistent.

## Criticality Level
**SUPREME LAW** - Violation causes immediate system failure (-100% penalty)

## Description

The orchestrator-state.json file tracks integration status in multiple locations that MUST remain synchronized:
- `integration_status` objects
- `stale_integration_tracking` records
- `efforts_completed` entries
- `current_wave_integration` / `current_phase_integration`
- `waves_completed` / `phases_completed` tracking
- `metadata_locations` references

When ANY of these change, ALL related fields MUST be updated atomically to prevent:
- Stale tracking showing "recreation_required" for completed integrations
- Integration status conflicting with stale tracking
- Completed efforts missing from completion lists
- Metadata locations pointing to non-existent files
- State transitions losing critical decisions

## Requirements

### 1. ATOMIC UPDATE REQUIREMENT
When updating integration status, ALL these MUST be updated together:
```javascript
// REQUIRED: Atomic update of all related fields
{
  // Clear stale tracking
  "stale_integration_tracking.stale_integrations": filter_out_completed,
  "stale_integration_tracking.recreation_completed": true,
  
  // Update integration status
  "integration_status.wave_X.status": "COMPLETE",
  "integration_status.wave_X.completed_at": timestamp,
  
  // Clear current integration flags
  "current_wave_integration.is_stale": false,
  "current_wave_integration.staleness_reason": null,
  
  // Update completion tracking
  "waves_completed.phase_X.wave_Y": completion_data
}
```

### 2. STALE TRACKING CLEANUP
When integration completes successfully:
```bash
# MANDATORY: Clean up stale tracking on successful integration
clean_stale_tracking_on_success() {
    local integration_id="$1"
    
    # Mark as recreated if it was stale
    jq --arg id "$integration_id" '
        # Update stale integration record
        (.stale_integration_tracking.stale_integrations[]? | 
         select(.integration_id == $id)) |= 
         . + {
            "recreation_completed": true, 
            "recreation_at": now | todate,
            "recreation_required": false
         }
    ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
}
```

### 3. VALIDATION BEFORE COMMIT
EVERY state update MUST be validated for consistency:
```bash
validate_state_consistency() {
    # Check for contradictions
    local issues=()
    
    # Check stale tracking vs completion
    if jq -e '.stale_integration_tracking.stale_integrations[] | 
              select(.recreation_required == true and .recreation_completed == true)' \
              orchestrator-state.json > /dev/null; then
        issues+=("Contradiction: integration marked both required and completed")
    fi
    
    # Check integration status consistency
    if jq -e '.current_wave_integration | 
              select(.is_stale == false and .staleness_reason != null)' \
              orchestrator-state.json > /dev/null; then
        issues+=("Contradiction: not stale but has staleness reason")
    fi
    
    if [ ${#issues[@]} -gt 0 ]; then
        echo "❌ STATE INCONSISTENCY DETECTED:"
        printf '%s\n' "${issues[@]}"
        return 1
    fi
    
    echo "✅ State consistency validated"
    return 0
}
```

### 4. STATE TRANSITION PRESERVATION
Critical decisions made during state transitions MUST be preserved:
```javascript
// REQUIRED: Preserve transition decisions
{
  "state_transition_decisions": {
    "timestamp": "2025-01-20T10:00:00Z",
    "from_state": "MONITORING_INTEGRATION",
    "to_state": "WAVE_COMPLETE",
    "decision": "Integration completed successfully",
    "metadata": {
      "build_status": "SUCCESS",
      "test_status": "PASSING",
      "demo_status": "PASSED",
      "integration_cleared": true,
      "stale_tracking_cleaned": true
    }
  }
}
```

### 5. RECOVERY VALIDATION
On orchestrator restart, MUST validate state consistency:
```bash
validate_on_recovery() {
    echo "🔍 Validating state consistency after recovery..."
    
    # Check for orphaned stale tracking
    local orphans=$(jq -r '.stale_integration_tracking.stale_integrations[]? | 
                           select(.recreation_completed == false) | 
                           .integration_id' orchestrator-state.json)
    
    for orphan in $orphans; do
        # Check if integration actually exists and is complete
        if jq -e --arg id "$orphan" '.waves_completed | .. | 
                  select(.integration_branch? == $id)' \
                  orchestrator-state.json > /dev/null; then
            echo "⚠️ Cleaning orphaned stale tracking for completed integration: $orphan"
            clean_stale_tracking_on_success "$orphan"
        fi
    done
    
    validate_state_consistency || exit 346
}
```

## Implementation

### Helper Function for Atomic Updates
```bash
atomic_integration_update() {
    local integration_type="$1"  # wave/phase/project
    local integration_id="$2"
    local status="$3"
    local timestamp="$(date -Iseconds)"
    
    # Create atomic update transaction
    jq --arg type "$integration_type" \
       --arg id "$integration_id" \
       --arg status "$status" \
       --arg time "$timestamp" '
        # Update ALL related fields atomically
        
        # 1. Update integration status
        .integration_status[$type + "_" + $id] = {
            "status": $status,
            "updated_at": $time
        } |
        
        # 2. Clean stale tracking if successful
        if $status == "COMPLETE" then
            .stale_integration_tracking.stale_integrations |= 
            map(if .integration_id == $id then
                . + {
                    "recreation_completed": true,
                    "recreation_at": $time,
                    "recreation_required": false
                }
            else . end)
        else . end |
        
        # 3. Update current integration flags
        .["current_" + $type + "_integration"] |= 
        if $status == "COMPLETE" then
            . + {
                "is_stale": false,
                "staleness_reason": null,
                "stale_since": null,
                "stale_due_to_fixes": []
            }
        else . end |
        
        # 4. Add to completion tracking if complete
        if $status == "COMPLETE" then
            .[$type + "s_completed"][$id] = {
                "completed_at": $time,
                "status": "COMPLETE"
            }
        else . end
    ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    
    # Validate the update
    validate_state_consistency || {
        echo "❌ R346 VIOLATION: Atomic update created inconsistent state!"
        exit 346
    }
}
```

## Common Violations

### ❌ CRITICAL VIOLATIONS (-100% penalty):
1. **Partial Updates**: Updating integration_status but not stale_tracking
2. **Contradictory State**: is_stale=false but recreation_required=true
3. **Lost Decisions**: State transitions not preserving critical metadata
4. **Orphaned Tracking**: Completed integrations still in stale_integrations
5. **Validation Skip**: Committing state without consistency validation

### ❌ Example Violations:
```bash
# WRONG: Partial update
jq '.integration_status.wave_2 = "COMPLETE"' orchestrator-state.json

# WRONG: Not cleaning stale tracking
jq '.current_wave_integration.is_stale = false' orchestrator-state.json

# WRONG: No validation before commit
git add orchestrator-state.json && git commit -m "update"
```

### ✅ Correct Pattern:
```bash
# RIGHT: Atomic update with validation
atomic_integration_update "wave" "phase1-wave2-integration" "COMPLETE"
validate_state_consistency && git add orchestrator-state.json
git commit -m "state: atomic update - integration complete, stale tracking cleaned"
```

## Enforcement

### Automatic Validation Hook
Add to all state transitions:
```bash
# In state transition functions
transition_state() {
    local new_state="$1"
    local reason="$2"
    
    # Update state
    update_orchestrator_state "$new_state" "$reason"
    
    # R346: MANDATORY validation
    validate_state_consistency || {
        echo "❌ R346 VIOLATION: State inconsistency after transition!"
        echo "Rolling back state change..."
        git checkout -- orchestrator-state.json
        exit 346
    }
    
    # Commit only if valid
    git add orchestrator-state.json
    git commit -m "state: transition to $new_state (R346 validated)"
}
```

### Recovery Protocol
On every orchestrator restart:
```bash
# MANDATORY on startup
echo "🔍 R346: Validating state consistency..."
validate_on_recovery
echo "✅ R346: State consistency verified"
```

## Related Rules
- R327: Mandatory Re-Integration After Fixes
- R328: Integration Freshness Validation  
- R288: State File Update and Commit
- R324: State Transition Validation
- R325: State Machine Authority

## Penalty Structure
- **Inconsistent state detected**: -100% IMMEDIATE FAILURE
- **Partial update without validation**: -100% FAILURE
- **Lost transition metadata**: -50% penalty
- **Missing atomic update**: -50% penalty
- **Skipped validation**: -30% penalty

## Success Criteria
✅ ALL related fields updated atomically
✅ NO contradictory state ever committed
✅ Stale tracking cleaned on completion
✅ State validated before every commit
✅ Recovery validates and fixes inconsistencies
✅ Transition decisions preserved

## SUPREME LAW DECLARATION
**This rule is ABSOLUTE. Any violation of state consistency is a catastrophic failure that corrupts the entire system. There are NO exceptions, NO workarounds, and NO forgiveness for violations.**