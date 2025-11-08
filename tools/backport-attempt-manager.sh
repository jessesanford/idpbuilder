#!/usr/bin/env bash
# tools/backport-attempt-manager.sh
#
# Manages backport attempt counters within integration iterations per R532
#
# This script tracks how many backport fix cycles have been attempted within
# the SAME iteration. This prevents infinite loops where the iteration counter
# stays constant but backport attempts repeat forever.
#
# Usage:
#   backport-attempt-manager.sh increment_backport_attempts <LEVEL>
#   backport-attempt-manager.sh get_backport_attempt_count <LEVEL>
#   backport-attempt-manager.sh reset_backport_attempts <LEVEL>
#   backport-attempt-manager.sh check_max_backport_attempts <LEVEL>
#
# Where <LEVEL> is: WAVE, PHASE, or PROJECT

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
STATE_FILE="$PROJECT_ROOT/orchestrator-state-v3.json"

# Validate state file exists
if [ ! -f "$STATE_FILE" ]; then
    echo "ERROR: State file not found: $STATE_FILE" >&2
    exit 1
fi

# Function: increment_backport_attempts
# Increment the backport_attempts_this_iteration counter for the specified level
increment_backport_attempts() {
    local level="$1"
    local field_path

    case "$level" in
        WAVE)
            field_path=".project_progression.current_wave"
            ;;
        PHASE)
            field_path=".project_progression.current_phase"
            ;;
        PROJECT)
            field_path=".project_progression.current_project"
            ;;
        *)
            echo "ERROR: Invalid level '$level'. Must be WAVE, PHASE, or PROJECT" >&2
            exit 1
            ;;
    esac

    # Get current count (default to 0 if missing)
    local current_count
    current_count=$(jq -r "${field_path}.backport_attempts_this_iteration // 0" "$STATE_FILE")

    # Increment
    local new_count=$((current_count + 1))

    # Update state file
    local tmp_file="${STATE_FILE}.tmp.$$"
    jq "${field_path}.backport_attempts_this_iteration = $new_count" "$STATE_FILE" > "$tmp_file"
    mv "$tmp_file" "$STATE_FILE"

    echo "$new_count"
}

# Function: get_backport_attempt_count
# Get the current backport_attempts_this_iteration for the specified level
get_backport_attempt_count() {
    local level="$1"
    local field_path

    case "$level" in
        WAVE)
            field_path=".project_progression.current_wave"
            ;;
        PHASE)
            field_path=".project_progression.current_phase"
            ;;
        PROJECT)
            field_path=".project_progression.current_project"
            ;;
        *)
            echo "ERROR: Invalid level '$level'. Must be WAVE, PHASE, or PROJECT" >&2
            exit 1
            ;;
    esac

    jq -r "${field_path}.backport_attempts_this_iteration // 0" "$STATE_FILE"
}

# Function: reset_backport_attempts
# Reset backport_attempts_this_iteration to 0 (called when starting new iteration)
reset_backport_attempts() {
    local level="$1"
    local field_path

    case "$level" in
        WAVE)
            field_path=".project_progression.current_wave"
            ;;
        PHASE)
            field_path=".project_progression.current_phase"
            ;;
        PROJECT)
            field_path=".project_progression.current_project"
            ;;
        *)
            echo "ERROR: Invalid level '$level'. Must be WAVE, PHASE, or PROJECT" >&2
            exit 1
            ;;
    esac

    # Reset to 0
    local tmp_file="${STATE_FILE}.tmp.$$"
    jq "${field_path}.backport_attempts_this_iteration = 0" "$STATE_FILE" > "$tmp_file"
    mv "$tmp_file" "$STATE_FILE"

    echo "0"
}

# Function: check_max_backport_attempts
# Check if backport attempts have exceeded the maximum for this iteration
# Returns: "WITHIN_LIMIT" or "EXCEEDED"
check_max_backport_attempts() {
    local level="$1"
    local field_path

    case "$level" in
        WAVE)
            field_path=".project_progression.current_wave"
            ;;
        PHASE)
            field_path=".project_progression.current_phase"
            ;;
        PROJECT)
            field_path=".project_progression.current_project"
            ;;
        *)
            echo "ERROR: Invalid level '$level'. Must be WAVE, PHASE, or PROJECT" >&2
            exit 1
            ;;
    esac

    # Get current backport attempts (default to 0)
    local current_count
    current_count=$(jq -r "${field_path}.backport_attempts_this_iteration // 0" "$STATE_FILE")

    # Get max backport attempts (default to 3 per R532)
    local max_attempts
    max_attempts=$(jq -r "${field_path}.max_backport_attempts_per_iteration // 3" "$STATE_FILE")

    if [ "$current_count" -ge "$max_attempts" ]; then
        echo "EXCEEDED"
    else
        echo "WITHIN_LIMIT"
    fi
}

# Main command dispatcher
main() {
    if [ $# -lt 2 ]; then
        echo "Usage: $0 <command> <level>" >&2
        echo "" >&2
        echo "Commands:" >&2
        echo "  increment_backport_attempts <LEVEL>    - Increment backport attempt counter" >&2
        echo "  get_backport_attempt_count <LEVEL>     - Get current backport attempt count" >&2
        echo "  reset_backport_attempts <LEVEL>        - Reset backport attempts to 0" >&2
        echo "  check_max_backport_attempts <LEVEL>    - Check if max exceeded" >&2
        echo "" >&2
        echo "Levels: WAVE, PHASE, PROJECT" >&2
        exit 1
    fi

    local command="$1"
    local level="$2"

    case "$command" in
        increment_backport_attempts)
            increment_backport_attempts "$level"
            ;;
        get_backport_attempt_count)
            get_backport_attempt_count "$level"
            ;;
        reset_backport_attempts)
            reset_backport_attempts "$level"
            ;;
        check_max_backport_attempts)
            check_max_backport_attempts "$level"
            ;;
        *)
            echo "ERROR: Unknown command '$command'" >&2
            exit 1
            ;;
    esac
}

main "$@"
