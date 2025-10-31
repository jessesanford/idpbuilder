#!/bin/bash
# iteration-manager.sh - Iteration tracking and enforcement for SF 3.0
#
# PURPOSE: Manage iteration counters and enforce max iteration limits for wave/phase/project containers
# USAGE:
#   bash tools/iteration-manager.sh increment_iteration WAVE|PHASE|PROJECT container_id
#   bash tools/iteration-manager.sh check_max_iterations WAVE|PHASE|PROJECT container_id
#   bash tools/iteration-manager.sh get_iteration_count WAVE|PHASE|PROJECT container_id
#
# RULES:
#   - R336: Integration Container Requirements
#   - R531: Integration Iteration Protocol
#   - R288: State File Update Requirements

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# State files
STATE_FILE="$PROJECT_ROOT/orchestrator-state-v3.json"

# Logging function
log() {
    echo -e "${BLUE}[iteration-manager]${NC} $1"
}

error() {
    echo -e "${RED}[iteration-manager ERROR]${NC} $1" >&2
}

success() {
    echo -e "${GREEN}[iteration-manager]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[iteration-manager WARNING]${NC} $1"
}

# Validate state file exists
validate_state_file() {
    if [ ! -f "$STATE_FILE" ]; then
        error "State file not found: $STATE_FILE"
        return 1
    fi

    if ! jq empty "$STATE_FILE" 2>/dev/null; then
        error "State file is not valid JSON: $STATE_FILE"
        return 1
    fi

    return 0
}

# Get current container from orchestrator-state-v3.json based on level
get_current_container() {
    local level="$1"  # WAVE, PHASE, or PROJECT

    case "$level" in
        WAVE)
            jq -r '.project_progression.current_wave' "$STATE_FILE"
            ;;
        PHASE)
            jq -r '.project_progression.current_phase' "$STATE_FILE"
            ;;
        PROJECT)
            jq -r '.project_progression.current_project' "$STATE_FILE"
            ;;
        *)
            error "Invalid level: $level (must be WAVE, PHASE, or PROJECT)"
            return 1
            ;;
    esac
}

# Increment iteration counter for a container
increment_iteration() {
    local level="$1"      # WAVE, PHASE, or PROJECT
    local container_id="$2"  # Optional: specific container ID, otherwise use current

    validate_state_file || return 1

    # If container_id not provided, get current container
    if [ -z "$container_id" ]; then
        container_id=$(get_current_container "$level")
        if [ "$container_id" = "null" ] || [ -z "$container_id" ]; then
            error "No current ${level} container found in state file"
            return 1
        fi
    fi

    log "Incrementing iteration for ${level} container: ${container_id}"

    # Get current iteration count
    local field_path
    case "$level" in
        WAVE)
            field_path=".project_progression.current_wave.iteration"
            ;;
        PHASE)
            field_path=".project_progression.current_phase.iteration"
            ;;
        PROJECT)
            field_path=".project_progression.current_project.iteration"
            ;;
        *)
            error "Invalid level: $level"
            return 1
            ;;
    esac

    local current_iteration
    current_iteration=$(jq -r "$field_path" "$STATE_FILE")

    if [ "$current_iteration" = "null" ]; then
        current_iteration=0
    fi

    local new_iteration=$((current_iteration + 1))

    log "Current iteration: ${current_iteration} → New iteration: ${new_iteration}"

    # Create backup before modification
    cp "$STATE_FILE" "$STATE_FILE.backup-iteration-increment"

    # Update iteration counter
    jq "$field_path = $new_iteration" "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"

    # Update last_progress_timestamp
    local timestamp
    timestamp=$(date -Iseconds)
    jq ".project_progression.iteration_tracking.last_progress_timestamp = \"$timestamp\"" "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"

    # Validate updated state file
    if ! jq empty "$STATE_FILE" 2>/dev/null; then
        error "State file corrupted after update, restoring backup"
        mv "$STATE_FILE.backup-iteration-increment" "$STATE_FILE"
        return 1
    fi

    # Remove backup on success
    rm -f "$STATE_FILE.backup-iteration-increment"

    success "✅ Incremented ${level} iteration: ${current_iteration} → ${new_iteration}"
    echo "$new_iteration"
    return 0
}

# Check if max iterations exceeded
check_max_iterations() {
    local level="$1"      # WAVE, PHASE, or PROJECT
    local container_id="$2"  # Optional: specific container ID, otherwise use current

    validate_state_file || return 1

    # If container_id not provided, get current container
    if [ -z "$container_id" ]; then
        container_id=$(get_current_container "$level")
        if [ "$container_id" = "null" ] || [ -z "$container_id" ]; then
            error "No current ${level} container found in state file"
            return 1
        fi
    fi

    log "Checking max iterations for ${level} container: ${container_id}"

    # Get current and max iterations
    local current_path max_path
    case "$level" in
        WAVE)
            current_path=".project_progression.current_wave.iteration"
            max_path=".project_progression.current_wave.max_iterations"
            ;;
        PHASE)
            current_path=".project_progression.current_phase.iteration"
            max_path=".project_progression.current_phase.max_iterations"
            ;;
        PROJECT)
            current_path=".project_progression.current_project.iteration"
            max_path=".project_progression.current_project.max_iterations"
            ;;
        *)
            error "Invalid level: $level"
            return 1
            ;;
    esac

    local current_iteration max_iterations
    current_iteration=$(jq -r "$current_path" "$STATE_FILE")
    max_iterations=$(jq -r "$max_path" "$STATE_FILE")

    if [ "$current_iteration" = "null" ]; then
        current_iteration=0
    fi

    if [ "$max_iterations" = "null" ]; then
        max_iterations=10  # Default per SF 3.0 Architecture
    fi

    log "Current iteration: ${current_iteration} / Max: ${max_iterations}"

    if [ "$current_iteration" -ge "$max_iterations" ]; then
        error "❌ MAX ITERATIONS EXCEEDED for ${level}"
        error "   Container: ${container_id}"
        error "   Current: ${current_iteration}"
        error "   Maximum: ${max_iterations}"
        error "   ACTION: Escalate to ERROR_RECOVERY"
        echo "EXCEEDED"
        return 1
    else
        local remaining=$((max_iterations - current_iteration))
        success "✅ Within limit: ${current_iteration}/${max_iterations} (${remaining} remaining)"
        echo "WITHIN_LIMIT"
        return 0
    fi
}

# Get current iteration count
get_iteration_count() {
    local level="$1"      # WAVE, PHASE, or PROJECT
    local container_id="$2"  # Optional: specific container ID, otherwise use current

    validate_state_file || return 1

    # If container_id not provided, get current container
    if [ -z "$container_id" ]; then
        container_id=$(get_current_container "$level")
        if [ "$container_id" = "null" ] || [ -z "$container_id" ]; then
            error "No current ${level} container found in state file"
            return 1
        fi
    fi

    # Get current iteration
    local field_path
    case "$level" in
        WAVE)
            field_path=".project_progression.current_wave.iteration"
            ;;
        PHASE)
            field_path=".project_progression.current_phase.iteration"
            ;;
        PROJECT)
            field_path=".project_progression.current_project.iteration"
            ;;
        *)
            error "Invalid level: $level"
            return 1
            ;;
    esac

    local iteration
    iteration=$(jq -r "$field_path" "$STATE_FILE")

    if [ "$iteration" = "null" ]; then
        iteration=0
    fi

    echo "$iteration"
    return 0
}

# Get total iteration count across all containers
get_total_iterations() {
    validate_state_file || return 1

    local wave_iter phase_iter project_iter
    wave_iter=$(jq -r '.project_progression.current_wave.iteration // 0' "$STATE_FILE")
    phase_iter=$(jq -r '.project_progression.current_phase.iteration // 0' "$STATE_FILE")
    project_iter=$(jq -r '.project_progression.current_project.iteration // 0' "$STATE_FILE")

    local total=$((wave_iter + phase_iter + project_iter))

    log "Total iterations: Wave=${wave_iter} + Phase=${phase_iter} + Project=${project_iter} = ${total}"
    echo "$total"
    return 0
}

# Main command dispatcher
main() {
    local command="${1:-}"

    if [ -z "$command" ]; then
        error "Usage: iteration-manager.sh <command> [args]"
        error "Commands:"
        error "  increment_iteration <WAVE|PHASE|PROJECT> [container_id]"
        error "  check_max_iterations <WAVE|PHASE|PROJECT> [container_id]"
        error "  get_iteration_count <WAVE|PHASE|PROJECT> [container_id]"
        error "  get_total_iterations"
        return 1
    fi

    case "$command" in
        increment_iteration)
            increment_iteration "${2:-}" "${3:-}"
            ;;
        check_max_iterations)
            check_max_iterations "${2:-}" "${3:-}"
            ;;
        get_iteration_count)
            get_iteration_count "${2:-}" "${3:-}"
            ;;
        get_total_iterations)
            get_total_iterations
            ;;
        *)
            error "Unknown command: $command"
            return 1
            ;;
    esac
}

# Run main if executed directly
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    main "$@"
fi
