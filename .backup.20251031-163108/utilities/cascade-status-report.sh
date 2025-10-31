#!/bin/bash
# cascade-status-report.sh
# Comprehensive R406 fix cascade status reporting
# Sourced by orchestrator states to auto-print cascade progress

# Color codes for terminal output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Main status reporting function
cascade_status_report() {
    local state_file="${1:-orchestrator-state-v3.json}"

    # Check if state file exists
    if [ ! -f "$state_file" ]; then
        echo "WARNING: orchestrator-state-v3.json not found"
        return 1
    fi

    # Check if cascade is active
    local cascade_active=$(jq -r '.fix_cascade_state.active // false' "$state_file")

    if [ "$cascade_active" != "true" ]; then
        echo "No active cascade"
        return 0
    fi

    # Extract cascade data
    local cascade_id=$(jq -r '.fix_cascade_state.cascade_id // "unknown"' "$state_file")
    local status=$(jq -r '.fix_cascade_state.status // "unknown"' "$state_file")
    local current_layer=$(jq -r '.fix_cascade_state.current_layer // 0' "$state_file")
    local total_layers=$(jq -r '.fix_cascade_state.total_layers // 0' "$state_file")
    local triggered_by=$(jq -r '.fix_cascade_state.triggered_by_integration.integration_name // "unknown"' "$state_file")
    local created_at=$(jq -r '.fix_cascade_state.created_at // "unknown"' "$state_file")

    # Bug counts
    local total_bugs=$(jq -r '.fix_cascade_state.validation.total_bugs_detected // 0' "$state_file")
    local total_fixed=$(jq -r '.fix_cascade_state.validation.total_bugs_fixed // 0' "$state_file")
    local total_pending=$(jq -r '.fix_cascade_state.validation.total_bugs_pending // 0' "$state_file")

    local pending_count=$(jq -r '[.bug_registry[]? | select(.fix_status == "pending")] | length' "$state_file")
    local in_progress_count=$(jq -r '[.bug_registry[]? | select(.fix_status == "in_progress")] | length' "$state_file")
    local fixed_count=$(jq -r '[.bug_registry[]? | select(.fix_status == "fixed")] | length' "$state_file")
    local verified_count=$(jq -r '[.bug_registry[]? | select(.fix_status == "verified")] | length' "$state_file")
    local blocked_count=$(jq -r '[.bug_registry[]? | select(.fix_status == "blocked")] | length' "$state_file")

    # Print report header
    echo ""
    echo "================================================================="
    echo "R406 FIX CASCADE STATUS (automatic report)"
    echo "================================================================="
    echo ""

    # Cascade overview
    echo "${BOLD}Cascade Overview:${NC}"
    echo "  Cascade ID: $cascade_id"
    echo "  Status: $(format_status_color "$status")"
    echo "  Layer: $current_layer / $total_layers"
    echo "  Triggered by: $triggered_by"
    echo "  Started: $created_at"
    echo ""

    # Bug registry summary
    echo "${BOLD}Bug Registry Summary:${NC}"
    echo "  Total Bugs: $total_bugs"
    if [ "$verified_count" -gt 0 ]; then
        echo "  ${GREEN}✅ Verified: $verified_count${NC}"
    fi
    if [ "$fixed_count" -gt 0 ]; then
        echo "  ${GREEN}✅ Fixed: $fixed_count${NC}"
    fi
    if [ "$in_progress_count" -gt 0 ]; then
        echo "  ${YELLOW}🔧 In Progress: $in_progress_count${NC}"
    fi
    if [ "$pending_count" -gt 0 ]; then
        echo "  ${CYAN}⏳ Pending: $pending_count${NC}"
    fi
    if [ "$blocked_count" -gt 0 ]; then
        echo "  ${RED}❌ Blocked: $blocked_count${NC}"
    fi
    echo ""

    # Integration status
    print_integration_status "$state_file"

    # Effort status
    print_effort_status "$state_file"

    # Cascade chain (layer by layer)
    print_cascade_chain "$state_file"

    # Next actions
    print_next_actions "$state_file"

    # Validation
    print_validation_status "$state_file"

    echo "================================================================="
}

# Format status with color
format_status_color() {
    local status="$1"

    case "$status" in
        "complete")
            echo "${GREEN}${BOLD}$status${NC}"
            ;;
        "fixing"|"in_progress")
            echo "${YELLOW}${BOLD}$status${NC}"
            ;;
        "failed"|"blocked"|"aborted")
            echo "${RED}${BOLD}$status${NC}"
            ;;
        *)
            echo "${CYAN}$status${NC}"
            ;;
    esac
}

# Print integration status breakdown
print_integration_status() {
    local state_file="$1"

    echo "${BOLD}Integration Status:${NC}"

    # Get all integrations with fixes
    local integrations=$(jq -r '.integration_fix_states // {} | keys[]?' "$state_file")

    if [ -z "$integrations" ]; then
        echo "  No integrations with tracked fixes"
        echo ""
        return
    fi

    echo "$integrations" | while read -r integration; do
        if [ -z "$integration" ]; then
            continue
        fi

        local int_status=$(jq -r ".integration_fix_states.\"$integration\".status // \"unknown\"" "$state_file")
        local bugs_total=$(jq -r ".integration_fix_states.\"$integration\".bugs_total_count // 0" "$state_file")
        local bugs_pending=$(jq -r ".integration_fix_states.\"$integration\".bugs_by_status.pending // 0" "$state_file")
        local bugs_in_progress=$(jq -r ".integration_fix_states.\"$integration\".bugs_by_status.in_progress // 0" "$state_file")
        local bugs_fixed=$(jq -r ".integration_fix_states.\"$integration\".bugs_by_status.fixed // 0" "$state_file")
        local bugs_verified=$(jq -r ".integration_fix_states.\"$integration\".bugs_by_status.verified // 0" "$state_file")
        local needs_reintegration=$(jq -r ".integration_fix_states.\"$integration\".requires_reintegration // false" "$state_file")

        echo ""
        echo "  ${BOLD}$integration:${NC}"
        echo "    Total Bugs: $bugs_total"

        if [ "$bugs_verified" -gt 0 ]; then
            echo "    ${GREEN}✅ Verified: $bugs_verified${NC}"
        fi
        if [ "$bugs_fixed" -gt 0 ]; then
            echo "    ${GREEN}✅ Fixed: $bugs_fixed${NC}"
        fi
        if [ "$bugs_in_progress" -gt 0 ]; then
            echo "    ${YELLOW}🔧 In Progress: $bugs_in_progress${NC}"
        fi
        if [ "$bugs_pending" -gt 0 ]; then
            echo "    ${CYAN}⏳ Pending: $bugs_pending${NC}"
        fi

        echo "    Status: $(format_status_color "$int_status")"

        if [ "$needs_reintegration" = "true" ]; then
            echo "    ${YELLOW}⚠️ Requires Reintegration: Yes${NC}"
        else
            echo "    Requires Reintegration: No"
        fi
    done

    echo ""
}

# Print effort status breakdown
print_effort_status() {
    local state_file="$1"

    echo "${BOLD}Effort Status:${NC}"

    # Get all efforts with fixes
    local efforts=$(jq -r '.effort_fix_states // {} | keys[]?' "$state_file")

    if [ -z "$efforts" ]; then
        echo "  No efforts with tracked fixes"
        echo ""
        return
    fi

    # Count efforts by status
    local efforts_with_pending=0
    local efforts_with_in_progress=0
    local efforts_complete=0

    echo "$efforts" | while read -r effort; do
        if [ -z "$effort" ]; then
            continue
        fi

        local bugs_total=$(jq -r ".effort_fix_states.\"$effort\".bugs_assigned | length" "$state_file")
        local bugs_pending=$(jq -r ".effort_fix_states.\"$effort\".bugs_by_status.pending // 0" "$state_file")
        local bugs_in_progress=$(jq -r ".effort_fix_states.\"$effort\".bugs_by_status.in_progress // 0" "$state_file")
        local bugs_fixed=$(jq -r ".effort_fix_states.\"$effort\".bugs_by_status.fixed // 0" "$state_file")
        local bugs_verified=$(jq -r ".effort_fix_states.\"$effort\".bugs_by_status.verified // 0" "$state_file")
        local ready=$(jq -r ".effort_fix_states.\"$effort\".ready_for_integration // false" "$state_file")

        local status_icon=""
        if [ "$bugs_pending" -gt 0 ]; then
            status_icon="${CYAN}⏳${NC}"
        elif [ "$bugs_in_progress" -gt 0 ]; then
            status_icon="${YELLOW}🔧${NC}"
        elif [ "$ready" = "true" ]; then
            status_icon="${GREEN}✅${NC}"
        fi

        echo "  $status_icon ${BOLD}$effort:${NC} $bugs_total bug(s)"

        if [ "$bugs_verified" -gt 0 ]; then
            echo "      ${GREEN}✅ Verified: $bugs_verified${NC}"
        fi
        if [ "$bugs_fixed" -gt 0 ]; then
            echo "      ${GREEN}✅ Fixed: $bugs_fixed${NC}"
        fi
        if [ "$bugs_in_progress" -gt 0 ]; then
            echo "      ${YELLOW}🔧 In Progress: $bugs_in_progress${NC}"
        fi
        if [ "$bugs_pending" -gt 0 ]; then
            echo "      ${CYAN}⏳ Pending: $bugs_pending${NC}"
        fi
    done

    echo ""
}

# Print cascade chain (layer progression)
print_cascade_chain() {
    local state_file="$1"

    local has_chain=$(jq -e '.fix_cascade_state.cascade_chain' "$state_file" > /dev/null 2>&1 && echo "true" || echo "false")

    if [ "$has_chain" != "true" ]; then
        return
    fi

    local chain_length=$(jq -r '.fix_cascade_state.cascade_chain | length' "$state_file")

    if [ "$chain_length" -eq 0 ]; then
        return
    fi

    echo "${BOLD}Cascade Chain (Layer Progression):${NC}"

    jq -r '.fix_cascade_state.cascade_chain[] |
        "  Layer \(.layer): \(.integration_name) (\(.type)) - \(.status)"' "$state_file" | \
    while read -r line; do
        # Color code based on status
        if echo "$line" | grep -q "complete"; then
            echo "${GREEN}$line${NC}"
        elif echo "$line" | grep -q "fixing\|in_progress\|reintegrating"; then
            echo "${YELLOW}$line${NC}"
        elif echo "$line" | grep -q "failed"; then
            echo "${RED}$line${NC}"
        else
            echo "$line"
        fi
    done

    echo ""
}

# Print next actions
print_next_actions() {
    local state_file="$1"

    echo "${BOLD}Next Actions:${NC}"

    # Determine next actions based on cascade state
    local status=$(jq -r '.fix_cascade_state.status // "unknown"' "$state_file")
    local pending=$(jq -r '[.bug_registry[]? | select(.fix_status == "pending")] | length' "$state_file")
    local in_progress=$(jq -r '[.bug_registry[]? | select(.fix_status == "in_progress")] | length' "$state_file")
    local blocked=$(jq -r '[.bug_registry[]? | select(.fix_status == "blocked")] | length' "$state_file")

    local action_num=1

    # Handle blocked bugs first
    if [ "$blocked" -gt 0 ]; then
        echo "  ${RED}${action_num}. RESOLVE BLOCKED BUGS ($blocked bugs)${NC}"
        jq -r '.bug_registry[]? | select(.fix_status == "blocked") |
            "     - \(.bug_id): \(.description // "no description")"' "$state_file"
        ((action_num++))
    fi

    # In-progress bugs
    if [ "$in_progress" -gt 0 ]; then
        echo "  ${YELLOW}${action_num}. Complete fixes in progress ($in_progress bugs)${NC}"
        jq -r '.bug_registry[]? | select(.fix_status == "in_progress") |
            "     - \(.bug_id) (\(.primary_effort)): \(.description // "no description")"' "$state_file"
        ((action_num++))
    fi

    # Pending bugs
    if [ "$pending" -gt 0 ]; then
        echo "  ${CYAN}${action_num}. Start fixes for pending bugs ($pending bugs)${NC}"
        jq -r '.bug_registry[]? | select(.fix_status == "pending") |
            "     - \(.bug_id) (\(.primary_effort)): \(.description // "no description")"' "$state_file" | head -5

        if [ "$pending" -gt 5 ]; then
            echo "     ... and $((pending - 5)) more"
        fi
        ((action_num++))
    fi

    # Check for integrations ready for reintegration
    local ready_for_reint=$(jq -r '[.integration_fix_states // {} |
        to_entries[] |
        select(.value.status == "ready_for_reintegration")] | length' "$state_file")

    if [ "$ready_for_reint" -gt 0 ]; then
        echo "  ${GREEN}${action_num}. Reintegrate fixed integrations ($ready_for_reint ready)${NC}"
        jq -r '.integration_fix_states // {} |
            to_entries[] |
            select(.value.status == "ready_for_reintegration") |
            "     - \(.key)"' "$state_file"
        ((action_num++))
    fi

    # Next layer cascade
    local current_layer=$(jq -r '.fix_cascade_state.current_layer // 0' "$state_file")
    local total_layers=$(jq -r '.fix_cascade_state.total_layers // 0' "$state_file")

    if [ "$current_layer" -lt "$total_layers" ] && [ "$pending" -eq 0 ] && [ "$in_progress" -eq 0 ]; then
        echo "  ${BLUE}${action_num}. Cascade to layer $((current_layer + 1)) of $total_layers${NC}"
        ((action_num++))
    fi

    # Completion
    if [ "$pending" -eq 0 ] && [ "$in_progress" -eq 0 ] && [ "$blocked" -eq 0 ] && \
       [ "$current_layer" -eq "$total_layers" ]; then
        echo "  ${GREEN}✅ ${action_num}. Complete cascade - all bugs fixed and integrated${NC}"
    fi

    # Default if no actions
    if [ "$action_num" -eq 1 ]; then
        echo "  No pending actions - cascade in steady state"
    fi

    echo ""
}

# Print validation status
print_validation_status() {
    local state_file="$1"

    echo "${BOLD}Validation:${NC}"

    # Count bugs in registry
    local registry_count=$(jq -r '.bug_registry | length' "$state_file")
    local validation_count=$(jq -r '.fix_cascade_state.validation.total_bugs_detected // 0' "$state_file")

    # Verify checksum
    local stored_checksum=$(jq -r '.fix_cascade_state.validation.checksum // ""' "$state_file")
    local last_validated=$(jq -r '.fix_cascade_state.validation.last_validated // "never"' "$state_file")

    # Calculate current checksum
    local bug_ids=$(jq -r '[.bug_registry[]?.bug_id] | sort | @csv' "$state_file")
    local current_checksum=$(echo "$bug_ids" | md5sum | cut -d' ' -f1)

    if [ "$registry_count" -eq "$validation_count" ]; then
        echo "  ${GREEN}✅ Bug count valid: $registry_count bugs accounted for${NC}"
    else
        echo "  ${RED}❌ Bug count MISMATCH: Registry=$registry_count, Validation=$validation_count${NC}"
    fi

    if [ "$stored_checksum" = "$current_checksum" ]; then
        echo "  ${GREEN}✅ Checksum valid: no lost bugs${NC}"
    elif [ -z "$stored_checksum" ]; then
        echo "  ${YELLOW}⚠️ No checksum stored${NC}"
    else
        echo "  ${RED}❌ Checksum MISMATCH: potential bug loss!${NC}"
    fi

    echo "  Last validated: $last_validated"
}

# Export function so it can be called from other scripts
export -f cascade_status_report

# If run directly (not sourced), execute the report
if [ "${BASH_SOURCE[0]}" -ef "$0" ]; then
    cascade_status_report "$@"
fi
