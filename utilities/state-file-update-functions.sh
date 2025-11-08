#!/bin/bash

# state-file-update-functions.sh - Helper functions for safe state file updates
# Source this file in orchestrator scripts for validated state management

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to validate state file
validate_state_file() {
    local state_file="${1:-orchestrator-state-v3.json}"
    
    if [ ! -f "$state_file" ]; then
        echo -e "${RED}❌ State file not found: $state_file${NC}"
        return 1
    fi
    
    # Use the validation tool if available
    if [ -f "$CLAUDE_PROJECT_DIR/tools/validate-state.sh" ]; then
        "$CLAUDE_PROJECT_DIR/tools/validate-state.sh" "$state_file"
        return $?
    else
        # Basic validation fallback
        if ! jq empty "$state_file" 2>/dev/null; then
            echo -e "${RED}❌ State file is not valid JSON!${NC}"
            return 1
        fi
        
        # Check for required fields
        local required_fields=(
            "current_state"
            "previous_state"
            "transition_time"
            "current_phase"
            "current_wave"
        )
        
        for field in "${required_fields[@]}"; do
            if ! jq -e "has(\"$field\")" "$state_file" >/dev/null 2>&1; then
                echo -e "${RED}❌ Missing required field: $field${NC}"
                return 1
            fi
        done
        
        echo -e "${GREEN}✅ State file passed basic validation${NC}"
        return 0
    fi
}

# Function to safely update a single field in state file
safe_update_field() {
    local field="$1"
    local value="$2"
    local state_file="${3:-orchestrator-state-v3.json}"
    
    # Create backup
    cp "$state_file" "${state_file}.backup"
    
    # Update field
    if [[ "$value" =~ ^[0-9]+$ ]]; then
        # Numeric value
        jq ".$field = $value" "$state_file" > tmp.json && mv tmp.json "$state_file"
    elif [[ "$value" == "true" || "$value" == "false" || "$value" == "null" ]]; then
        # Boolean or null
        jq ".$field = $value" "$state_file" > tmp.json && mv tmp.json "$state_file"
    else
        # String value
        jq ".$field = \"$value\"" "$state_file" > tmp.json && mv tmp.json "$state_file"
    fi
    
    # Validate after update
    if ! validate_state_file "$state_file"; then
        echo -e "${RED}❌ Validation failed after updating $field${NC}"
        echo "Restoring backup..."
        mv "${state_file}.backup" "$state_file"
        return 1
    fi
    
    # Remove backup on success
    rm -f "${state_file}.backup"
    return 0
}

# Function to perform safe state transition with validation
safe_state_transition() {
    local new_state="$1"
    local reason="${2:-State transition}"
    local state_file="${3:-orchestrator-state-v3.json}"
    
    # Get current state
    local current_state=$(jq -r '.current_state' "$state_file")
    
    if [ "$current_state" == "$new_state" ]; then
        echo -e "${YELLOW}⚠️  Already in state $new_state${NC}"
        return 0
    fi
    
    echo "🔄 Transitioning: $current_state → $new_state"
    
    # Create backup before transition
    cp "$state_file" "${state_file}.pre-transition"
    
    # Update all transition fields
    local timestamp=$(date -u +%Y-%m-%dT%H:%M:%SZ)
    
    jq ".current_state = \"$new_state\" |
        .previous_state = \"$current_state\" |
        .transition_time = \"$timestamp\" |
        .transition_reason = \"$reason\"" "$state_file" > tmp.json && mv tmp.json "$state_file"
    
    # Add to transition history if the field exists
    if jq -e 'has("transition_history")' "$state_file" >/dev/null 2>&1; then
        local transition_entry="{\"from\": \"$current_state\", \"to\": \"$new_state\", \"time\": \"$timestamp\", \"reason\": \"$reason\"}"
        jq ".transition_history = (.transition_history // []) + [$transition_entry]" "$state_file" > tmp.json && mv tmp.json "$state_file"
    fi
    
    # Validate the updated state
    if ! validate_state_file "$state_file"; then
        echo -e "${RED}❌❌❌ State validation failed after transition!${NC}"
        echo "Restoring pre-transition state..."
        mv "${state_file}.pre-transition" "$state_file"
        return 1
    fi
    
    # Commit and push immediately (R288)
    git add "$state_file"
    git commit -m "state: $current_state → $new_state - $reason [R288/R324]"
    
    if ! git push; then
        echo -e "${RED}❌ Failed to push state transition!${NC}"
        echo "Attempting force push with lease..."
        git push --force-with-lease || {
            echo -e "${RED}❌❌❌ CRITICAL: Could not push state change!${NC}"
            return 1
        }
    fi
    
    # Clean up backup
    rm -f "${state_file}.pre-transition"
    
    echo -e "${GREEN}✅ State transition complete: $current_state → $new_state${NC}"
    return 0
}

# Function to update and validate effort status
update_effort_status() {
    local effort_name="$1"
    local status="$2"
    local phase="${3:-}"
    local wave="${4:-}"
    local state_file="${5:-orchestrator-state-v3.json}"
    
    echo "📝 Updating effort '$effort_name' status to '$status'"
    
    # Find and update the effort in the appropriate array
    local updated=false
    
    # Try efforts_in_progress first
    if jq -e ".efforts_in_progress[] | select(.name == \"$effort_name\")" "$state_file" >/dev/null 2>&1; then
        jq "(.efforts_in_progress[] | select(.name == \"$effort_name\") | .status) = \"$status\"" "$state_file" > tmp.json && mv tmp.json "$state_file"
        updated=true
    fi
    
    # Try efforts_completed if not found
    if [ "$updated" = false ] && jq -e ".efforts_completed[] | select(.name == \"$effort_name\")" "$state_file" >/dev/null 2>&1; then
        jq "(.efforts_completed[] | select(.name == \"$effort_name\") | .status) = \"$status\"" "$state_file" > tmp.json && mv tmp.json "$state_file"
        updated=true
    fi
    
    if [ "$updated" = false ]; then
        echo -e "${YELLOW}⚠️  Effort '$effort_name' not found in state file${NC}"
        return 1
    fi
    
    # Validate after update
    if ! validate_state_file "$state_file"; then
        echo -e "${RED}❌ State validation failed after updating effort status${NC}"
        return 1
    fi
    
    # Commit the change
    git add "$state_file"
    git commit -m "state: effort '$effort_name' status → $status [R288]"
    git push
    
    echo -e "${GREEN}✅ Effort status updated and validated${NC}"
    return 0
}

# Function to add violation record with validation
add_violation() {
    local effort="$1"
    local lines="$2"
    local limit="${3:-800}"
    local state_file="${4:-orchestrator-state-v3.json}"
    
    echo "⚠️  Recording violation for effort '$effort': $lines lines (limit: $limit)"
    
    # Create violation entry
    local violation_entry="{
        \"effort\": \"$effort\",
        \"lines\": $lines,
        \"limit\": $limit,
        \"requires_split\": true,
        \"recorded_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
    }"
    
    # Add to violations array
    jq ".violations = (.violations // []) + [$violation_entry]" "$state_file" > tmp.json && mv tmp.json "$state_file"
    
    # Validate
    if ! validate_state_file "$state_file"; then
        echo -e "${RED}❌ State validation failed after adding violation${NC}"
        return 1
    fi
    
    # Commit
    git add "$state_file"
    git commit -m "state: violation recorded for '$effort' - $lines lines [R288]"
    git push
    
    echo -e "${GREEN}✅ Violation recorded and validated${NC}"
    return 0
}

# Function to verify state machine transition is valid
verify_transition_allowed() {
    local from_state="$1"
    local to_state="$2"
    local state_machine_file="${3:-$CLAUDE_PROJECT_DIR/state-machines/software-factory-3.0-state-machine.json}"
    
    # Check if transition exists in state machine
    if grep -q "$from_state --> $to_state" "$state_machine_file"; then
        echo -e "${GREEN}✅ Transition allowed: $from_state → $to_state${NC}"
        return 0
    else
        echo -e "${RED}❌ Invalid transition: $from_state → $to_state${NC}"
        echo "This transition is not defined in the state machine!"
        return 1
    fi
}

# Function to perform emergency state recovery
emergency_state_recovery() {
    local state_file="${1:-orchestrator-state-v3.json}"
    
    echo -e "${YELLOW}🚨 Attempting emergency state recovery...${NC}"
    
    # Check for backups
    if [ -f "${state_file}.backup" ]; then
        echo "Found backup file, validating..."
        if validate_state_file "${state_file}.backup"; then
            cp "${state_file}.backup" "$state_file"
            echo -e "${GREEN}✅ Restored from backup${NC}"
            return 0
        fi
    fi
    
    # Check git history for last valid version
    echo "Checking git history for last valid state..."
    git checkout HEAD -- "$state_file"
    
    if validate_state_file "$state_file"; then
        echo -e "${GREEN}✅ Restored from git${NC}"
        return 0
    fi
    
    echo -e "${RED}❌ Could not recover valid state file${NC}"
    echo "Manual intervention required!"
    return 1
}

# R346: Atomic integration update function
atomic_integration_update() {
    local integration_type="$1"  # wave/phase/project
    local integration_id="$2"
    local status="$3"
    local state_file="${4:-orchestrator-state-v3.json}"
    local timestamp="$(date -Iseconds)"
    
    echo "🔄 R346: Performing atomic integration update for $integration_id..."
    
    # Create backup before atomic update
    cp "$state_file" "${state_file}.atomic-backup"
    
    # Perform atomic update of ALL related fields
    jq --arg type "$integration_type" \
       --arg id "$integration_id" \
       --arg status "$status" \
       --arg time "$timestamp" '
        # Update integration status
        .integration_status[$type + "_" + $id] = {
            "status": $status,
            "updated_at": $time
        } |
        
        # Clean stale tracking if successful
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
        
        # Update current integration flags
        .["current_" + $type + "_integration"] |= 
        if $status == "COMPLETE" then
            . + {
                "is_stale": false,
                "staleness_reason": null,
                "stale_since": null,
                "stale_due_to_fixes": [],
                "last_freshness_check": $time
            }
        else . end |
        
        # Add to completion tracking if complete
        if $status == "COMPLETE" then
            .[$type + "s_completed"][$id] = {
                "completed_at": $time,
                "status": "COMPLETE"
            }
        else . end
    ' "$state_file" > tmp.json && mv tmp.json "$state_file"
    
    # R346: Validate consistency after atomic update
    if [ -f "$CLAUDE_PROJECT_DIR/utilities/validate-state-consistency.sh" ]; then
        if ! "$CLAUDE_PROJECT_DIR/utilities/validate-state-consistency.sh" "$state_file"; then
            echo -e "${RED}❌ R346 VIOLATION: Atomic update created inconsistent state!${NC}"
            echo "Rolling back to backup..."
            mv "${state_file}.atomic-backup" "$state_file"
            return 346
        fi
    else
        # Fallback validation
        if ! validate_state_file "$state_file"; then
            echo -e "${RED}❌ Validation failed after atomic update${NC}"
            mv "${state_file}.atomic-backup" "$state_file"
            return 1
        fi
    fi
    
    # Remove backup on success
    rm -f "${state_file}.atomic-backup"
    
    # Commit the atomic update
    git add "$state_file"
    git commit -m "state: atomic integration update - $integration_id → $status [R346]"
    git push
    
    echo -e "${GREEN}✅ R346: Atomic integration update complete and validated${NC}"
    return 0
}

# R346: Clean stale tracking for completed integration
clean_stale_tracking_on_success() {
    local integration_id="$1"
    local state_file="${2:-orchestrator-state-v3.json}"
    local timestamp="$(date -Iseconds)"
    
    echo "🧹 R346: Cleaning stale tracking for $integration_id..."
    
    # Update stale integration record
    jq --arg id "$integration_id" --arg time "$timestamp" '
        # Mark as recreated if it exists
        (.stale_integration_tracking.stale_integrations[]? | 
         select(.integration_id == $id)) |= 
         . + {
            "recreation_completed": true, 
            "recreation_at": $time,
            "recreation_required": false
         }
    ' "$state_file" > tmp.json && mv tmp.json "$state_file"
    
    # Validate after cleanup
    if [ -f "$CLAUDE_PROJECT_DIR/utilities/validate-state-consistency.sh" ]; then
        "$CLAUDE_PROJECT_DIR/utilities/validate-state-consistency.sh" "$state_file"
    else
        validate_state_file "$state_file"
    fi
    
    echo -e "${GREEN}✅ R346: Stale tracking cleaned${NC}"
    return 0
}

# R346: Validate state consistency
validate_state_consistency() {
    local state_file="${1:-orchestrator-state-v3.json}"
    
    # Use the dedicated R346 validator if available
    if [ -f "$CLAUDE_PROJECT_DIR/utilities/validate-state-consistency.sh" ]; then
        "$CLAUDE_PROJECT_DIR/utilities/validate-state-consistency.sh" "$state_file"
        return $?
    fi
    
    # Fallback: basic consistency checks
    local issues=()
    
    # Check for contradictions in stale tracking
    if jq -e '.stale_integration_tracking.stale_integrations[]? | 
              select(.recreation_required == true and .recreation_completed == true)' \
              "$state_file" > /dev/null 2>&1; then
        issues+=("Contradiction: integration marked both required and completed")
    fi
    
    # Check integration status consistency
    if jq -e '.current_wave_integration | 
              select(.is_stale == false and .staleness_reason != null)' \
              "$state_file" > /dev/null 2>&1; then
        issues+=("Contradiction: not stale but has staleness reason")
    fi
    
    if [ ${#issues[@]} -gt 0 ]; then
        echo -e "${RED}❌ R346 VIOLATION: STATE INCONSISTENCY DETECTED:${NC}"
        printf '%s\n' "${issues[@]}"
        return 346
    fi
    
    echo -e "${GREEN}✅ R346: State consistency validated${NC}"
    return 0
}

# Export functions for use in other scripts
export -f validate_state_file
export -f safe_update_field
export -f safe_state_transition
export -f update_effort_status
export -f add_violation
export -f verify_transition_allowed
export -f emergency_state_recovery
export -f atomic_integration_update
export -f clean_stale_tracking_on_success
export -f validate_state_consistency

# ═══════════════════════════════════════════════════════════════════════════
# R349: ATOMIC STATUS FLAG CLEANUP FUNCTIONS
# ═══════════════════════════════════════════════════════════════════════════

# Function to clean rebase flags atomically
clean_rebase_flags() {
    local effort="$1"
    local state_file="${2:-orchestrator-state-v3.json}"
    local timestamp="$(date -Iseconds)"
    
    echo "🧹 R349: Cleaning rebase flags for $effort"
    
    jq --arg effort "$effort" --arg time "$timestamp" '
        # Clean base_branch_tracking in all efforts
        (.. | objects | select(has("name") and .name == $effort and has("base_branch_tracking")) | 
         .base_branch_tracking) |= 
        . + {
            "requires_rebase": false,
            "rebase_reason": null,
            "last_rebase": $time,
            "integration_eligible": true
        } |
        
        # Clean any phase_X_rebase_required sections
        (.. | objects | select(has("needs_rebase"))) |=
        if (.effort == $effort or .name == $effort) then
            . + {
                "status": "completed",
                "completed_at": $time,
                "needs_rebase": false
            }
        else . end
    ' "$state_file" > tmp.json && mv tmp.json "$state_file"
    
    echo -e "${GREEN}✅ R349: Rebase flags cleaned for $effort${NC}"
}

# Function to clean split flags atomically
clean_split_flags() {
    local effort="$1"
    local state_file="${2:-orchestrator-state-v3.json}"
    local timestamp="$(date -Iseconds)"
    
    echo "🧹 R349: Cleaning split flags for $effort"
    
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
        if has("split_tracking") and (.split_tracking[$effort]? != null) then
            .split_tracking[$effort] |= . + {
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
    ' "$state_file" > tmp.json && mv tmp.json "$state_file"
    
    echo -e "${GREEN}✅ R349: Split flags cleaned for $effort${NC}"
}

# Function to clean fix flags atomically
clean_fix_flags() {
    local effort="$1"
    local state_file="${2:-orchestrator-state-v3.json}"
    local timestamp="$(date -Iseconds)"
    
    echo "🧹 R349: Cleaning fix flags for $effort"
    
    jq --arg effort "$effort" --arg time "$timestamp" '
        # Clean effort entry
        (.. | objects | select(.name? == $effort)) |=
        if has("has_fixes_applied") or has("needs_fix") then
            . + {
                "has_fixes_applied": true,
                "fixes_applied_at": $time,
                "needs_fix": false,
                "fix_status": "applied"
            }
        else . end |
        
        # Move fixes from pending to applied if applicable
        if has("stale_integration_tracking") then
            .stale_integration_tracking.fix_tracking |=
            if (.fixes_pending_integration | length) > 0 then
                . + {
                    "fixes_applied": .fixes_applied + .fixes_pending_integration,
                    "fixes_pending_integration": []
                }
            else . end
        else . end
    ' "$state_file" > tmp.json && mv tmp.json "$state_file"
    
    echo -e "${GREEN}✅ R349: Fix flags cleaned for $effort${NC}"
}

# Function to clean staleness flags atomically
clean_staleness_flags() {
    local integration_type="$1"  # wave/phase/project
    local state_file="${2:-orchestrator-state-v3.json}"
    local timestamp="$(date -Iseconds)"
    
    echo "🧹 R349: Cleaning staleness flags for $integration_type integration"
    
    jq --arg type "$integration_type" --arg time "$timestamp" '
        # Clean current integration staleness
        if has("current_" + $type + "_integration") then
            .["current_" + $type + "_integration"] |= . + {
                "is_stale": false,
                "staleness_reason": null,
                "stale_since": null,
                "stale_due_to_fixes": []
            }
        else . end |
        
        # Clean stale_integration_tracking
        if has("stale_integration_tracking") and has("stale_integrations") then
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
        else . end
    ' "$state_file" > tmp.json && mv tmp.json "$state_file"
    
    echo -e "${GREEN}✅ R349: Staleness flags cleaned for $integration_type${NC}"
}

# R349: Validate flag consistency
validate_flag_consistency() {
    local state_file="${1:-orchestrator-state-v3.json}"
    local errors=()
    
    echo "🔍 R349: Validating flag consistency..."
    
    # Check for contradictory rebase flags
    if jq -e '.. | objects | select(
        has("completed_at") and has("needs_rebase") and 
        .completed_at != null and .needs_rebase == true
    )' "$state_file" > /dev/null 2>&1; then
        errors+=("Rebase completed but needs_rebase still true")
    fi
    
    # Check for contradictory split flags
    if jq -e '.violations[]? | select(
        has("split_completed") and .split_completed == true and .requires_split == true
    )' "$state_file" > /dev/null 2>&1; then
        errors+=("Split completed but requires_split still true")
    fi
    
    # Check for contradictory staleness flags
    if jq -e '.. | objects | select(
        has("is_stale") and has("staleness_reason") and
        .is_stale == false and .staleness_reason != null and
        .staleness_reason != ""
    )' "$state_file" > /dev/null 2>&1; then
        errors+=("Not stale but has staleness_reason")
    fi
    
    # Check for contradictory recreation flags
    if jq -e '.stale_integration_tracking.stale_integrations[]? | select(
        .recreation_completed == true and .recreation_required == true
    )' "$state_file" > /dev/null 2>&1; then
        errors+=("Recreation completed but still required")
    fi
    
    if [ ${#errors[@]} -gt 0 ]; then
        echo -e "${RED}❌ R349 VIOLATION: Flag inconsistencies detected:${NC}"
        printf '%s\n' "${errors[@]}"
        return 1
    fi
    
    echo -e "${GREEN}✅ R349: Flag consistency validated${NC}"
    return 0
}

# R349: Wrapper for state updates with automatic cleanup
update_with_cleanup() {
    local update_type="$1"
    local effort="$2"
    local state_file="${3:-orchestrator-state-v3.json}"
    
    echo "📝 R349: Performing $update_type with automatic cleanup for $effort"
    
    # Perform the update and cleanup
    case "$update_type" in
        "rebase_complete")
            clean_rebase_flags "$effort" "$state_file"
            ;;
        "split_complete")
            clean_split_flags "$effort" "$state_file"
            ;;
        "fix_applied")
            clean_fix_flags "$effort" "$state_file"
            ;;
        "integration_recreated")
            clean_staleness_flags "$effort" "$state_file"
            ;;
        *)
            echo -e "${YELLOW}⚠️ Unknown update type: $update_type${NC}"
            ;;
    esac
    
    # MANDATORY validation
    validate_flag_consistency "$state_file" || {
        echo -e "${RED}❌ R349 VIOLATION: Flags inconsistent after $update_type${NC}"
        git checkout -- "$state_file"
        return 349
    }
    
    # Commit with R349 marker
    git add "$state_file"
    git commit -m "state: $update_type for $effort with atomic flag cleanup [R349]"
    git push
    
    echo -e "${GREEN}✅ R349: Update complete with clean flags${NC}"
    return 0
}

# R349: Clean stale flags on recovery
clean_stale_flags_on_recovery() {
    local state_file="${1:-orchestrator-state-v3.json}"
    
    echo "🔍 R349: Detecting and cleaning stale flags on recovery..."
    
    # Find completed actions with stale flags
    local stale_rebases=$(jq -r '.. | objects | select(
        has("completed_at") and has("needs_rebase") and
        .completed_at != null and .needs_rebase == true
    ) | .effort // .name // "unknown"' "$state_file" 2>/dev/null | sort -u)
    
    for effort in $stale_rebases; do
        if [[ "$effort" != "unknown" && "$effort" != "null" ]]; then
            echo "⚠️ Cleaning stale rebase flag for $effort"
            clean_rebase_flags "$effort" "$state_file"
        fi
    done
    
    # Find completed splits with stale flags
    local stale_splits=$(jq -r '.violations[]? | select(
        has("split_completed") and .split_completed == true and .requires_split == true
    ) | .effort' "$state_file" 2>/dev/null | sort -u)
    
    for effort in $stale_splits; do
        echo "⚠️ Cleaning stale split flag for $effort"
        clean_split_flags "$effort" "$state_file"
    done
    
    # Clean contradictory staleness flags
    for type in wave phase project; do
        local is_stale=$(jq -r ".current_${type}_integration.is_stale // \"null\"" "$state_file")
        local reason=$(jq -r ".current_${type}_integration.staleness_reason // \"null\"" "$state_file")
        
        if [[ "$is_stale" == "false" && "$reason" != "null" && "$reason" != "" ]]; then
            echo "⚠️ Cleaning contradictory staleness for $type integration"
            clean_staleness_flags "$type" "$state_file"
        fi
    done
    
    # Validate all flags are clean
    validate_flag_consistency "$state_file" || {
        echo -e "${RED}❌ R349: Unable to clean all stale flags!${NC}"
        return 349
    }
    
    echo -e "${GREEN}✅ R349: All stale flags cleaned${NC}"
    return 0
}

# Export R349 functions
export -f clean_rebase_flags
export -f clean_split_flags
export -f clean_fix_flags
export -f clean_staleness_flags
export -f validate_flag_consistency
export -f update_with_cleanup
export -f clean_stale_flags_on_recovery

echo "✅ State file update functions loaded (including R346 atomic updates + R349 flag cleanup)"
echo "Available functions:"
echo "  - validate_state_file [file]"
echo "  - safe_update_field <field> <value> [file]"
echo "  - safe_state_transition <new_state> [reason] [file]"
echo "  - update_effort_status <effort> <status> [phase] [wave] [file]"
echo "  - add_violation <effort> <lines> [limit] [file]"
echo "  - verify_transition_allowed <from> <to> [state_machine_file]"
echo "  - emergency_state_recovery [file]"
echo "  - atomic_integration_update <type> <id> <status> [file] (R346)"
echo "  - clean_stale_tracking_on_success <id> [file] (R346)"
echo "  - validate_state_consistency [file] (R346)"
echo "  - clean_rebase_flags <effort> [file] (R349)"
echo "  - clean_split_flags <effort> [file] (R349)"
echo "  - clean_fix_flags <effort> [file] (R349)"
echo "  - clean_staleness_flags <type> [file] (R349)"
echo "  - validate_flag_consistency [file] (R349)"
echo "  - update_with_cleanup <type> <effort> [file] (R349)"
echo "  - clean_stale_flags_on_recovery [file] (R349)"