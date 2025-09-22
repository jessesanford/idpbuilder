# 🔴🔴🔴 RULE R349: ATOMIC STATUS FLAG CLEANUP (SUPREME LAW)

## Summary
ALL status flags and metadata MUST be atomically cleaned up when their corresponding actions complete. NO stale flags may persist after action completion.

## Criticality Level
**SUPREME LAW** - Violation causes system confusion and wrong decisions (-100% penalty)

## Description

Status flags in orchestrator-state.json track pending actions (rebases, splits, fixes, recreations). When these actions complete, ALL related flags MUST be atomically updated to prevent:
- Confusion from stale "needs_rebase" when rebase is complete
- Duplicate work from stale "needs_split" when split is done
- Wrong decisions from stale "needs_fix" when fixes are applied
- Integration failures from stale "recreation_required" flags
- Contradictory state where completed_at exists but status is still pending

## Requirements

### 1. ATOMIC FLAG CLEANUP ON ACTION COMPLETION

#### A. REBASE COMPLETION
When a rebase completes, ALL these MUST be updated atomically:
```javascript
{
  // In base_branch_tracking
  "requires_rebase": false,
  "rebase_reason": null,
  "last_rebase": "2025-01-20T10:00:00Z",
  "integration_eligible": true,
  
  // In phase_X_rebase_required (if exists)
  "status": "completed",
  "completed_at": "2025-01-20T10:00:00Z",
  "needs_rebase": false  // CRITICAL: Clear this flag!
}
```

#### B. SPLIT COMPLETION
When a split completes, ALL these MUST be updated atomically:
```javascript
{
  // In violations entry
  "requires_split": false,
  "split_planned": false,
  "split_completed": true,
  "split_completed_at": "2025-01-20T11:00:00Z",
  
  // In split_tracking
  "status": "COMPLETED",
  "completed_at": "2025-01-20T11:00:00Z",
  
  // In effort status
  "status": "SPLIT_COMPLETED",
  "needs_split": false  // CRITICAL: Clear this flag!
}
```

#### C. FIX APPLICATION
When fixes are applied, ALL these MUST be updated atomically:
```javascript
{
  // In effort entry
  "has_fixes_applied": true,
  "fixes_applied_at": "2025-01-20T12:00:00Z",
  "needs_fix": false,  // CRITICAL: Clear this flag!
  "fix_status": "applied",
  
  // In fix_tracking
  "fixes_applied": [...],
  "fixes_pending_integration": []
}
```

#### D. INTEGRATION RECREATION
When integration is recreated, ALL these MUST be updated atomically:
```javascript
{
  // In stale_integration_tracking
  "recreation_required": false,
  "recreation_completed": true,
  "recreation_at": "2025-01-20T13:00:00Z",
  
  // In current integration
  "is_stale": false,
  "staleness_reason": null,
  "stale_since": null,
  "stale_due_to_fixes": []
}
```

### 2. MANDATORY CLEANUP FUNCTIONS

Every state update function MUST use these atomic cleanup patterns:

```bash
# Function to clean rebase flags atomically
clean_rebase_flags() {
    local effort="$1"
    local timestamp="$(date -Iseconds)"
    
    jq --arg effort "$effort" --arg time "$timestamp" '
        # Clean base_branch_tracking
        (.. | .base_branch_tracking? | select(. != null)) |= 
        . + {
            "requires_rebase": false,
            "rebase_reason": null,
            "last_rebase": $time,
            "integration_eligible": true
        } |
        
        # Clean any phase_X_rebase_required sections
        (.. | objects | select(has("needs_rebase"))) |=
        . + {
            "status": "completed",
            "completed_at": $time,
            "needs_rebase": false
        }
    ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
}

# Function to clean split flags atomically
clean_split_flags() {
    local effort="$1"
    local timestamp="$(date -Iseconds)"
    
    jq --arg effort "$effort" --arg time "$timestamp" '
        # Clean violations entry
        (.violations[]? | select(.effort == $effort)) |=
        . + {
            "requires_split": false,
            "split_planned": false,
            "split_completed": true,
            "split_completed_at": $time
        } |
        
        # Clean split_tracking
        (.split_tracking[$effort]?) |=
        if . != null then
            . + {
                "status": "COMPLETED",
                "completed_at": $time
            }
        else . end |
        
        # Clean effort status flags
        (.. | objects | select(.name? == $effort and has("needs_split"))) |=
        . + {
            "needs_split": false,
            "split_status": "completed"
        }
    ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
}

# Function to clean fix flags atomically
clean_fix_flags() {
    local effort="$1"
    local timestamp="$(date -Iseconds)"
    
    jq --arg effort "$effort" --arg time "$timestamp" '
        # Clean effort entry
        (.. | objects | select(.name? == $effort)) |=
        . + {
            "has_fixes_applied": true,
            "fixes_applied_at": $time,
            "needs_fix": false,
            "fix_status": "applied"
        } |
        
        # Move fixes from pending to applied
        .stale_integration_tracking.fix_tracking |=
        if .fixes_pending_integration | length > 0 then
            . + {
                "fixes_applied": .fixes_applied + .fixes_pending_integration,
                "fixes_pending_integration": []
            }
        else . end
    ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
}

# Function to clean staleness flags atomically
clean_staleness_flags() {
    local integration_type="$1"  # wave/phase
    local timestamp="$(date -Iseconds)"
    
    jq --arg type "$integration_type" --arg time "$timestamp" '
        # Clean current integration staleness
        .["current_" + $type + "_integration"] |=
        . + {
            "is_stale": false,
            "staleness_reason": null,
            "stale_since": null,
            "stale_due_to_fixes": []
        } |
        
        # Clean stale_integration_tracking
        .stale_integration_tracking.stale_integrations |=
        map(
            if .integration_type == $type then
                . + {
                    "recreation_required": false,
                    "recreation_completed": true,
                    "recreation_at": $time
                }
            else . end
        )
    ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
}
```

### 3. VALIDATION REQUIREMENTS

#### A. Pre-Commit Validation
EVERY state update MUST validate flag consistency:
```bash
validate_flag_consistency() {
    local errors=()
    
    # Check for contradictory rebase flags
    if jq -e '.. | objects | select(
        has("completed_at") and has("needs_rebase") and 
        .completed_at != null and .needs_rebase == true
    )' orchestrator-state.json > /dev/null; then
        errors+=("Rebase completed but needs_rebase still true")
    fi
    
    # Check for contradictory split flags
    if jq -e '.violations[]? | select(
        .split_completed == true and .requires_split == true
    )' orchestrator-state.json > /dev/null; then
        errors+=("Split completed but requires_split still true")
    fi
    
    # Check for contradictory staleness flags
    if jq -e '.. | objects | select(
        has("is_stale") and has("staleness_reason") and
        .is_stale == false and .staleness_reason != null
    )' orchestrator-state.json > /dev/null; then
        errors+=("Not stale but has staleness_reason")
    fi
    
    # Check for contradictory recreation flags
    if jq -e '.stale_integration_tracking.stale_integrations[]? | select(
        .recreation_completed == true and .recreation_required == true
    )' orchestrator-state.json > /dev/null; then
        errors+=("Recreation completed but still required")
    fi
    
    if [ ${#errors[@]} -gt 0 ]; then
        echo "❌ R349 VIOLATION: Flag inconsistencies detected:"
        printf '%s\n' "${errors[@]}"
        return 1
    fi
    
    return 0
}
```

#### B. State Transition Hook
ALL state transitions MUST clean relevant flags:
```bash
transition_with_cleanup() {
    local new_state="$1"
    local reason="$2"
    
    # Clean flags based on transition
    case "$new_state" in
        "WAVE_COMPLETE")
            # Clean all wave-level flags
            clean_staleness_flags "wave"
            ;;
        "PHASE_COMPLETE")
            # Clean all phase-level flags
            clean_staleness_flags "phase"
            ;;
        "SPLIT_COMPLETE")
            # Clean split flags for completed efforts
            for effort in $(jq -r '.split_tracking | keys[]' orchestrator-state.json); do
                if [[ $(jq -r ".split_tracking.$effort.status" orchestrator-state.json) == "COMPLETED" ]]; then
                    clean_split_flags "$effort"
                fi
            done
            ;;
    esac
    
    # Perform transition
    safe_state_transition "$new_state" "$reason"
    
    # Validate after transition
    validate_flag_consistency || exit 349
}
```

### 4. RECOVERY PROTOCOL

On orchestrator startup, MUST detect and clean stale flags:
```bash
clean_stale_flags_on_recovery() {
    echo "🔍 R349: Detecting and cleaning stale flags..."
    
    # Find completed actions with stale flags
    local stale_rebases=$(jq -r '.. | objects | select(
        has("completed_at") and has("needs_rebase") and
        .completed_at != null and .needs_rebase == true
    ) | .effort // .name // "unknown"' orchestrator-state.json)
    
    for effort in $stale_rebases; do
        echo "⚠️ Cleaning stale rebase flag for $effort"
        clean_rebase_flags "$effort"
    done
    
    # Find completed splits with stale flags
    local stale_splits=$(jq -r '.violations[]? | select(
        .split_completed == true and .requires_split == true
    ) | .effort' orchestrator-state.json)
    
    for effort in $stale_splits; do
        echo "⚠️ Cleaning stale split flag for $effort"
        clean_split_flags "$effort"
    done
    
    # Validate all flags are clean
    validate_flag_consistency || {
        echo "❌ R349: Unable to clean all stale flags!"
        exit 349
    }
    
    echo "✅ R349: All stale flags cleaned"
}
```

## Implementation Checklist

### Required Updates:
1. ✅ validate-state-consistency.sh - Add R349 flag validation
2. ✅ state-file-update-functions.sh - Add atomic cleanup functions
3. ✅ All state transition points - Add cleanup calls
4. ✅ Orchestrator startup - Add recovery protocol
5. ✅ Integration completion - Clean all flags
6. ✅ Split completion - Clean split flags
7. ✅ Fix application - Clean fix flags
8. ✅ Rebase completion - Clean rebase flags

## Common Violations

### ❌ CRITICAL VIOLATIONS (-100% penalty):
1. **Stale Flag Persistence**: completed_at exists but needs_X flag still true
2. **Contradictory State**: recreation_completed=true but recreation_required=true
3. **Partial Cleanup**: Some flags cleaned but not others
4. **No Validation**: Skipping flag consistency check
5. **Lost Updates**: Flag cleanup not persisted to disk

### ❌ Example Violations:
```bash
# WRONG: Partial flag update
jq '.phase_2_rebase_required.completed_at = "'$(date -Iseconds)'"' state.json
# But needs_rebase flag not cleared!

# WRONG: Non-atomic update
jq '.is_stale = false' state.json
jq '.staleness_reason = null' state.json
# Multiple commands = non-atomic!

# WRONG: No validation
clean_rebase_flags "effort"
git add state.json && git commit -m "cleanup"
# No validate_flag_consistency call!
```

### ✅ Correct Pattern:
```bash
# RIGHT: Atomic cleanup with validation
clean_rebase_flags "image-builder"
validate_flag_consistency || exit 349
git add orchestrator-state.json
git commit -m "state: atomic cleanup of rebase flags for image-builder [R349]"
git push
```

## Enforcement

### Automatic Enforcement in All Updates:
```bash
# Wrapper for all state updates
update_with_cleanup() {
    local update_type="$1"
    shift
    local args="$@"
    
    # Perform the update
    case "$update_type" in
        "rebase_complete")
            mark_rebase_complete "$args"
            clean_rebase_flags "$args"
            ;;
        "split_complete")
            mark_split_complete "$args"
            clean_split_flags "$args"
            ;;
        "fix_applied")
            mark_fix_applied "$args"
            clean_fix_flags "$args"
            ;;
        "integration_recreated")
            mark_integration_recreated "$args"
            clean_staleness_flags "$args"
            ;;
    esac
    
    # MANDATORY validation
    validate_flag_consistency || {
        echo "❌ R349 VIOLATION: Flags inconsistent after $update_type"
        git checkout -- orchestrator-state.json
        exit 349
    }
    
    # Commit with R349 marker
    git add orchestrator-state.json
    git commit -m "state: $update_type with atomic flag cleanup [R349]"
    git push
}
```

## Related Rules
- R346: State Metadata Synchronization (parent rule)
- R327: Mandatory Re-Integration After Fixes
- R328: Integration Freshness Validation
- R288: State File Update and Commit
- R324: State Transition Validation

## Penalty Structure
- **Stale flag detected**: -100% IMMEDIATE FAILURE
- **Non-atomic update**: -100% FAILURE
- **Contradictory flags**: -100% FAILURE
- **Missing validation**: -50% penalty
- **Incomplete cleanup**: -50% penalty

## Success Criteria
✅ NO stale flags persist after action completion
✅ ALL flag updates are atomic
✅ Flag consistency validated before every commit
✅ Recovery protocol cleans all stale flags
✅ No contradictory state ever exists
✅ All transitions trigger appropriate cleanup

## SUPREME LAW DECLARATION
**This rule is ABSOLUTE. Stale flags cause cascade failures throughout the system. There are NO exceptions to atomic cleanup requirements.**