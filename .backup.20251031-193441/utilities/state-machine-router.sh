#!/bin/bash
# State Machine Router - Helps orchestrator decide which sub-state machine to invoke
# Usage: state-machine-router.sh <condition> <current_state>

set -euo pipefail

CLAUDE_PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}"
STATE_FILE="${1:-orchestrator-state-v3.json}"

# Source validation functions
source "$CLAUDE_PROJECT_DIR/utilities/validation-functions.sh" 2>/dev/null || true

# Function to determine if sub-state machine should be invoked
route_to_sub_machine() {
    local current_state="$1"
    local condition="${2:-}"

    echo "📍 Current State: $current_state"
    echo "📋 Condition: $condition"

    case "$condition" in
        "NEEDS_SPLIT"|"SIZE_VIOLATION")
            echo "🔄 Routing to: SPLITTING SUB-STATE MACHINE"
            echo "📁 State Machine: state-machines/splitting-state-machine.json"
            echo "🎯 Entry Point: SPLIT_INIT"
            echo "ACTION: INVOKE_SPLITTING"
            return 0
            ;;

        "FIX_CASCADE"|"CASCADE_DETECTED")
            echo "🔄 Routing to: FIX CASCADE SUB-STATE MACHINE"
            echo "📁 State Machine: state-machines/fix-cascade-state-machine.json"
            echo "🎯 Entry Point: FIX_CASCADE_INIT"
            echo "ACTION: INVOKE_FIX_CASCADE"
            return 0
            ;;

        "INTEGRATE_WAVE_EFFORTS_NEEDED"|"WAVE_INTEGRATE_WAVE_EFFORTS"|"INTEGRATE_PHASE_WAVES"|"PROJECT_INTEGRATE_WAVE_EFFORTS")
            echo "🔄 Routing to: INTEGRATE_WAVE_EFFORTS SUB-STATE MACHINE"
            echo "📁 State Machine: state-machines/integration-state-machine.json"
            echo "🎯 Entry Point: INTEGRATE_WAVE_EFFORTS"
            echo "ACTION: INVOKE_INTEGRATE_WAVE_EFFORTS"
            return 0
            ;;

        "PR_READY"|"PR_PREPARATION")
            echo "🔄 Routing to: PR READY SUB-STATE MACHINE"
            echo "📁 State Machine: state-machines/pr-ready-state-machine.json"
            echo "🎯 Entry Point: PR_INIT"
            echo "ACTION: INVOKE_PR_READY"
            return 0
            ;;

        "PROJECT_INIT"|"NEW_PROJECT")
            echo "🔄 Routing to: INITIALIZATION SUB-STATE MACHINE"
            echo "📁 State Machine: state-machines/initialization-state-machine.json"
            echo "🎯 Entry Point: INIT_START"
            echo "ACTION: INVOKE_INITIALIZATION"
            return 0
            ;;

        *)
            echo "ℹ️ No sub-state machine routing needed"
            echo "ACTION: CONTINUE_MAIN"
            return 1
            ;;
    esac
}

# Function to check if a sub-state machine is already active
check_active_sub_machine() {
    if [[ -f "$STATE_FILE" ]]; then
        local active=$(jq -r '.sub_state_machine.active // false' "$STATE_FILE")
        if [[ "$active" == "true" ]]; then
            local type=$(jq -r '.sub_state_machine.type // "UNKNOWN"' "$STATE_FILE")
            local state=$(jq -r '.sub_state_machine.current_state // "UNKNOWN"' "$STATE_FILE")
            echo "⚠️ Sub-state machine already active!"
            echo "   Type: $type"
            echo "   Current State: $state"
            return 0
        fi
    fi
    return 1
}

# Function to prepare sub-state machine invocation
prepare_invocation() {
    local machine_type="$1"
    local entry_state="$2"
    local context="${3:-}"

    echo "🔧 Preparing invocation for: $machine_type"

    # Generate state file name
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local state_file="${machine_type,,}-state-${timestamp}.json"

    echo "📄 Sub-state file: $state_file"
    echo "🎯 Entry state: $entry_state"

    # Return the invocation details
    cat <<EOF
{
  "sub_state_machine": {
    "active": true,
    "type": "$machine_type",
    "state_file": "$state_file",
    "current_state": "$entry_state",
    "return_state": "$(jq -r '.current_state' "$STATE_FILE")",
    "started_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "context": "$context"
  }
}
EOF
}

# Main logic
main() {
    local condition="${1:-}"
    local current_state="${2:-$(jq -r '.current_state' "$STATE_FILE")}"

    echo "🏭 SOFTWARE FACTORY STATE MACHINE ROUTER"
    echo "========================================="

    # Check if sub-machine already active
    if check_active_sub_machine; then
        echo "❌ Cannot route - sub-state machine already active"
        exit 1
    fi

    # Determine routing
    if route_to_sub_machine "$current_state" "$condition"; then
        echo ""
        echo "✅ Sub-state machine routing determined"
        echo ""
        echo "To invoke, update orchestrator-state-v3.json with:"

        case "$condition" in
            "NEEDS_SPLIT"|"SIZE_VIOLATION")
                prepare_invocation "SPLITTING" "SPLIT_INIT" "$condition"
                ;;
            "FIX_CASCADE"|"CASCADE_DETECTED")
                prepare_invocation "FIX_CASCADE" "FIX_CASCADE_INIT" "$condition"
                ;;
            "INTEGRATE_WAVE_EFFORTS_NEEDED"|*"INTEGRATE_WAVE_EFFORTS")
                prepare_invocation "INTEGRATE_WAVE_EFFORTS" "INTEGRATE_WAVE_EFFORTS" "$condition"
                ;;
            "PR_READY"|"PR_PREPARATION")
                prepare_invocation "PR_READY" "PR_INIT" "$condition"
                ;;
            "PROJECT_INIT"|"NEW_PROJECT")
                prepare_invocation "INITIALIZATION" "INIT_START" "$condition"
                ;;
        esac
    else
        echo "✅ Continue with main state machine"
    fi
}

# Execute if run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi