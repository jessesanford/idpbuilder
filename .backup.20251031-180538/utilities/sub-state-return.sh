#!/bin/bash
# Sub-State Machine Return Handler - Returns control from sub-state to main orchestrator
# Usage: sub-state-return.sh <sub_state_file> <completion_status>

set -euo pipefail

CLAUDE_PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}"
MAIN_STATE_FILE="orchestrator-state-v3.json"

# Function to handle return from sub-state machine
return_from_sub_machine() {
    local sub_state_file="$1"
    local completion_status="${2:-COMPLETE}"

    echo "🔄 RETURNING FROM SUB-STATE MACHINE"
    echo "===================================="

    # Check if main state file exists
    if [[ ! -f "$MAIN_STATE_FILE" ]]; then
        echo "❌ ERROR: Main state file not found: $MAIN_STATE_FILE"
        exit 1
    fi

    # Get sub-state machine info
    local sub_machine_info=$(jq -r '.sub_state_machine' "$MAIN_STATE_FILE")

    if [[ "$sub_machine_info" == "null" ]]; then
        echo "❌ ERROR: No active sub-state machine found in main state"
        exit 1
    fi

    # Extract return information
    local machine_type=$(jq -r '.sub_state_machine.type' "$MAIN_STATE_FILE")
    local return_state=$(jq -r '.sub_state_machine.return_state' "$MAIN_STATE_FILE")
    local started_at=$(jq -r '.sub_state_machine.started_at' "$MAIN_STATE_FILE")

    echo "📋 Sub-Machine Type: $machine_type"
    echo "🔙 Return State: $return_state"
    echo "⏰ Started: $started_at"
    echo "✅ Status: $completion_status"

    # Archive the sub-state file
    if [[ -f "$sub_state_file" ]]; then
        local archive_dir="archived-states/sub-machines"
        mkdir -p "$archive_dir"
        local archive_name="${archive_dir}/$(basename "$sub_state_file" .json)-${completion_status}-$(date +%Y%m%d-%H%M%S).json"
        cp "$sub_state_file" "$archive_name"
        echo "📦 Archived sub-state to: $archive_name"
    fi

    # Update main state file based on machine type and completion
    case "$machine_type" in
        "SPLITTING")
            handle_splitting_return "$completion_status" "$return_state"
            ;;
        "FIX_CASCADE")
            handle_fix_cascade_return "$completion_status" "$return_state"
            ;;
        "INTEGRATE_WAVE_EFFORTS")
            handle_integration_return "$completion_status" "$return_state"
            ;;
        "PR_READY")
            handle_pr_ready_return "$completion_status" "$return_state"
            ;;
        "INITIALIZATION")
            handle_initialization_return "$completion_status" "$return_state"
            ;;
        *)
            echo "⚠️ Unknown sub-machine type: $machine_type"
            ;;
    esac

    # Clear sub-state machine from main state
    jq 'del(.sub_state_machine) | .current_state = $state' \
        --arg state "$return_state" \
        "$MAIN_STATE_FILE" > tmp.json && mv tmp.json "$MAIN_STATE_FILE"

    echo "✅ Successfully returned to main orchestrator"
    echo "📍 Current State: $return_state"
}

# Handle return from splitting sub-state machine
handle_splitting_return() {
    local status="$1"
    local return_state="$2"

    echo "📊 Processing splitting completion..."

    if [[ "$status" == "COMPLETE" ]]; then
        # Mark original effort as deprecated
        jq '.split_tracking |= map_values(. + {status: "SPLITS_COMPLETE"})' \
            "$MAIN_STATE_FILE" > tmp.json && mv tmp.json "$MAIN_STATE_FILE"

        echo "✅ All splits completed successfully"
    elif [[ "$status" == "ABORTED" ]]; then
        # Mark split as aborted
        jq '.split_tracking |= map_values(. + {status: "SPLIT_ABORTED"})' \
            "$MAIN_STATE_FILE" > tmp.json && mv tmp.json "$MAIN_STATE_FILE"

        echo "⚠️ Split processing aborted"
    fi
}

# Handle return from fix cascade sub-state machine
handle_fix_cascade_return() {
    local status="$1"
    local return_state="$2"

    echo "🔧 Processing fix cascade completion..."

    if [[ "$status" == "COMPLETE" ]]; then
        # Update cascade tracking
        jq '.cascade_coordination.cascade_mode = false |
            .cascade_coordination.status = "CASCADE_COMPLETE"' \
            "$MAIN_STATE_FILE" > tmp.json && mv tmp.json "$MAIN_STATE_FILE"

        echo "✅ Fix cascade completed successfully"
    fi
}

# Handle return from integration sub-state machine
handle_integration_return() {
    local status="$1"
    local return_state="$2"

    echo "🔗 Processing integration completion..."

    if [[ "$status" == "COMPLETE" ]]; then
        # Update integration status
        jq '.integration_status = "COMPLETE"' \
            "$MAIN_STATE_FILE" > tmp.json && mv tmp.json "$MAIN_STATE_FILE"

        echo "✅ Integration completed successfully"
    fi
}

# Handle return from PR ready sub-state machine
handle_pr_ready_return() {
    local status="$1"
    local return_state="$2"

    echo "📝 Processing PR ready completion..."

    if [[ "$status" == "COMPLETE" ]]; then
        # Update PR status
        jq '.pr_ready_status = "COMPLETE"' \
            "$MAIN_STATE_FILE" > tmp.json && mv tmp.json "$MAIN_STATE_FILE"

        echo "✅ PR preparation completed successfully"
    fi
}

# Handle return from initialization sub-state machine
handle_initialization_return() {
    local status="$1"
    local return_state="$2"

    echo "🚀 Processing initialization completion..."

    if [[ "$status" == "COMPLETE" ]]; then
        # Update initialization status
        jq '.initialization_complete = true' \
            "$MAIN_STATE_FILE" > tmp.json && mv tmp.json "$MAIN_STATE_FILE"

        echo "✅ Initialization completed successfully"
    fi
}

# Main execution
main() {
    local sub_state_file="${1:-}"
    local completion_status="${2:-COMPLETE}"

    if [[ -z "$sub_state_file" ]]; then
        echo "❌ ERROR: Sub-state file required"
        echo "Usage: $0 <sub_state_file> [completion_status]"
        exit 1
    fi

    return_from_sub_machine "$sub_state_file" "$completion_status"
}

# Execute if run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi