#!/usr/bin/env bash
# cleanup-completed-agents.sh
#
# Implements R610 (Agent Metadata Lifecycle Protocol) and R611 (Active Agents Cleanup Protocol)
# Automatically cleans up completed agents from active_agents and moves them to agents_history
#
# Usage:
#   bash tools/cleanup-completed-agents.sh                    # Cleanup all completed agents
#   bash tools/cleanup-completed-agents.sh --validate         # Validate only (no cleanup)
#   bash tools/cleanup-completed-agents.sh --agent-id ID      # Cleanup specific agent
#   bash tools/cleanup-completed-agents.sh --dry-run          # Show what would be cleaned
#
# Dependencies: jq, date
# Rules: R610 (BLOCKING), R611 (WARNING), R612 (STANDARD)

set -euo pipefail

# Configuration
STATE_FILE="${CLAUDE_PROJECT_DIR:-$(pwd)}/orchestrator-state-v3.json"
BACKUP_DIR=".state-backup"
SCHEMA_FILE="${CLAUDE_PROJECT_DIR:-$(pwd)}/schemas/agents-history-schema.json"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Usage
usage() {
    cat <<EOF
Usage: $(basename "$0") [OPTIONS]

Cleanup completed agents from active_agents per R610/R611.

OPTIONS:
    --validate          Validate only (no cleanup)
    --agent-id ID       Cleanup specific agent by ID
    --dry-run           Show what would be cleaned without making changes
    --help              Show this help message

EXAMPLES:
    # Cleanup all completed agents
    $(basename "$0")

    # Validate no stale agents exist
    $(basename "$0") --validate

    # Cleanup specific agent
    $(basename "$0") --agent-id swe-1.2.1-docker-client

    # Dry run to see what would be cleaned
    $(basename "$0") --dry-run

RULES:
    R610 - Agent Metadata Lifecycle Protocol (BLOCKING)
    R611 - Active Agents Cleanup Protocol (WARNING)
    R612 - Agent History Management (STANDARD)

EOF
    exit 0
}

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ️  $*${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $*${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $*${NC}"
}

log_error() {
    echo -e "${RED}❌ $*${NC}"
}

# Validate prerequisites
validate_prerequisites() {
    log_info "R610/R611: Validating prerequisites..."

    # Check jq
    if ! command -v jq &> /dev/null; then
        log_error "jq is required but not installed"
        exit 1
    fi

    # Check state file exists
    if [ ! -f "$STATE_FILE" ]; then
        log_error "State file not found: $STATE_FILE"
        exit 1
    fi

    # Validate JSON
    if ! jq empty "$STATE_FILE" 2>/dev/null; then
        log_error "State file is not valid JSON: $STATE_FILE"
        exit 1
    fi

    # Check active_agents exists
    if ! jq -e '.active_agents' "$STATE_FILE" >/dev/null 2>&1; then
        log_error "active_agents array not found in state file"
        exit 1
    fi

    # Create backup directory if needed
    mkdir -p "$BACKUP_DIR"

    log_success "Prerequisites validated"
}

# Backup state file
backup_state_file() {
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local backup_file="${BACKUP_DIR}/orchestrator-state-v3.json.backup-agent-cleanup-${timestamp}"

    cp "$STATE_FILE" "$backup_file"
    log_success "Backup created: $backup_file"
}

# Find completed agents
find_completed_agents() {
    local agent_id_filter="${1:-}"

    if [ -n "$agent_id_filter" ]; then
        # Specific agent
        jq -r --arg agent_id "$agent_id_filter" '
            .active_agents[] |
            select(.agent_id == $agent_id) |
            select(.state == "COMPLETE" or .state == "COMPLETED") |
            .agent_id
        ' "$STATE_FILE"
    else
        # All completed agents
        jq -r '
            .active_agents[] |
            select(.state == "COMPLETE" or .state == "COMPLETED") |
            .agent_id
        ' "$STATE_FILE"
    fi
}

# Extract agent metadata for history (per R612)
extract_agent_metadata() {
    local agent_id="$1"

    # Determine agent type from agent_id
    local agent_type=""
    if [[ "$agent_id" =~ ^swe- ]]; then
        agent_type="sw-engineer"
    elif [[ "$agent_id" =~ ^reviewer- ]]; then
        agent_type="code-reviewer"
    elif [[ "$agent_id" =~ ^architect- ]]; then
        agent_type="architect"
    elif [[ "$agent_id" =~ ^integration- ]]; then
        agent_type="integration"
    else
        log_warning "Unknown agent type for: $agent_id, defaulting to generic"
        agent_type="unknown"
    fi

    # Extract metadata based on agent type
    case "$agent_type" in
        "sw-engineer")
            extract_swe_metadata "$agent_id"
            ;;
        "code-reviewer")
            extract_reviewer_metadata "$agent_id"
            ;;
        "architect")
            extract_architect_metadata "$agent_id"
            ;;
        "integration")
            extract_integration_metadata "$agent_id"
            ;;
        *)
            extract_generic_metadata "$agent_id"
            ;;
    esac
}

# Extract SW-Engineer metadata
extract_swe_metadata() {
    local agent_id="$1"

    jq --arg agent_id "$agent_id" '
        .active_agents[] |
        select(.agent_id == $agent_id) |
        {
            agent_id: .agent_id,
            agent_type: "sw-engineer",
            final_state: .state,
            completed_at: (now | strftime("%Y-%m-%dT%H:%M:%SZ")),
            work_summary: {
                effort_id: .effort_id,
                effort_name: .effort_name,
                branch: .branch_name,
                outcome: .status,
                splits: (.split_count // 0)
            },
            metrics: {
                lines_added: (.line_count_tracking.final_count // 0),
                lines_removed: (.line_count_tracking.lines_removed // 0),
                files_modified: (.files_modified_count // 0),
                commits: (.commits_made // 0),
                duration_hours: (
                    if .completed_at and .spawned_at then
                        (((.completed_at | fromdate) - (.spawned_at | fromdate)) / 3600)
                    else 0 end
                )
            }
        }
    ' "$STATE_FILE"
}

# Extract Code-Reviewer metadata
extract_reviewer_metadata() {
    local agent_id="$1"

    jq --arg agent_id "$agent_id" '
        .active_agents[] |
        select(.agent_id == $agent_id) |
        {
            agent_id: .agent_id,
            agent_type: "code-reviewer",
            final_state: .state,
            completed_at: (now | strftime("%Y-%m-%dT%H:%M:%SZ")),
            work_summary: {
                reviewed_effort: .focus,
                review_type: (.review_type // "effort_review"),
                bugs_found: (.extraction_metrics.total_bugs // 0),
                review_outcome: (.review_result // "completed")
            },
            metrics: {
                review_duration_minutes: (.duration_minutes // 0),
                files_reviewed: (.files_reviewed_count // 0),
                critical_issues: (.extraction_metrics.critical_count // 0),
                blocking_issues: (.extraction_metrics.blocking_count // 0),
                warnings: (.extraction_metrics.warning_count // 0)
            }
        }
    ' "$STATE_FILE"
}

# Extract Architect metadata
extract_architect_metadata() {
    local agent_id="$1"

    jq --arg agent_id "$agent_id" '
        .active_agents[] |
        select(.agent_id == $agent_id) |
        {
            agent_id: .agent_id,
            agent_type: "architect",
            final_state: .state,
            completed_at: (now | strftime("%Y-%m-%dT%H:%M:%SZ")),
            work_summary: {
                review_scope: (.review_scope // "wave"),
                wave_id: (.wave_id // .focus),
                architecture_outcome: (.outcome // "approved"),
                concerns_raised: (.concerns_count // 0)
            },
            metrics: {
                efforts_reviewed: (.efforts_reviewed_count // 0),
                review_duration_minutes: (.duration_minutes // 0),
                recommendations_made: (.recommendations_count // 0)
            }
        }
    ' "$STATE_FILE"
}

# Extract Integration metadata
extract_integration_metadata() {
    local agent_id="$1"

    jq --arg agent_id "$agent_id" '
        .active_agents[] |
        select(.agent_id == $agent_id) |
        {
            agent_id: .agent_id,
            agent_type: "integration",
            final_state: .state,
            completed_at: (now | strftime("%Y-%m-%dT%H:%M:%SZ")),
            work_summary: {
                integration_type: (.integration_type // "wave"),
                wave_id: (.wave_id // .focus),
                efforts_integrated: (.efforts_merged_count // 0),
                outcome: (.integration_outcome // "success")
            },
            metrics: {
                merge_conflicts: (.merge_conflicts_count // 0),
                test_failures: (.test_failures_count // 0),
                integration_duration_minutes: (.duration_minutes // 0)
            }
        }
    ' "$STATE_FILE"
}

# Extract generic metadata (fallback)
extract_generic_metadata() {
    local agent_id="$1"

    jq --arg agent_id "$agent_id" '
        .active_agents[] |
        select(.agent_id == $agent_id) |
        {
            agent_id: .agent_id,
            agent_type: (.agent_type // "unknown"),
            final_state: .state,
            completed_at: (now | strftime("%Y-%m-%dT%H:%M:%SZ")),
            work_summary: {
                focus: (.focus // .effort_name // "unknown"),
                outcome: (.status // .outcome // "completed")
            },
            metrics: {}
        }
    ' "$STATE_FILE"
}

# Cleanup single agent
cleanup_agent() {
    local agent_id="$1"
    local dry_run="${2:-false}"

    log_info "Processing agent: $agent_id"

    # Extract metadata for history
    local metadata=$(extract_agent_metadata "$agent_id")

    if [ -z "$metadata" ] || [ "$metadata" = "null" ]; then
        log_error "Failed to extract metadata for agent: $agent_id"
        return 1
    fi

    if [ "$dry_run" = "true" ]; then
        log_info "[DRY RUN] Would move to agents_history:"
        echo "$metadata" | jq '.'
        return 0
    fi

    # Atomic update: add to history, remove from active
    local temp_file=$(mktemp)
    jq --arg agent_id "$agent_id" \
       --argjson metadata "$metadata" '
        # Initialize agents_history if missing
        if .agents_history == null then
            .agents_history = []
        else . end |

        # Add to history
        .agents_history += [$metadata] |

        # Remove from active
        .active_agents = [.active_agents[] | select(.agent_id != $agent_id)]
    ' "$STATE_FILE" > "$temp_file"

    # Validate result
    if ! jq empty "$temp_file" 2>/dev/null; then
        log_error "JSON validation failed after cleanup - restoring backup"
        rm "$temp_file"
        return 1
    fi

    # Replace original
    mv "$temp_file" "$STATE_FILE"

    log_success "Agent $agent_id moved to agents_history"
    return 0
}

# Cleanup all completed agents
cleanup_all_completed() {
    local dry_run="${1:-false}"

    log_info "R610/R611: Finding completed agents..."

    # Find all completed agents
    local completed_agents=$(find_completed_agents)

    if [ -z "$completed_agents" ]; then
        log_success "No completed agents found - nothing to cleanup"
        return 0
    fi

    # Count agents
    local agent_count=$(echo "$completed_agents" | wc -l)
    log_info "Found $agent_count completed agent(s)"

    # Backup before cleanup (unless dry run)
    if [ "$dry_run" != "true" ]; then
        backup_state_file
    fi

    # Cleanup each agent
    local cleanup_count=0
    local failed_count=0

    while IFS= read -r agent_id; do
        if cleanup_agent "$agent_id" "$dry_run"; then
            ((cleanup_count++))
        else
            ((failed_count++))
        fi
    done <<< "$completed_agents"

    # Report results
    if [ "$dry_run" = "true" ]; then
        log_info "[DRY RUN] Would cleanup $cleanup_count agent(s)"
    else
        log_success "Cleanup complete: $cleanup_count agent(s) moved to agents_history"

        if [ "$failed_count" -gt 0 ]; then
            log_warning "$failed_count agent(s) failed to cleanup"
            return 1
        fi
    fi

    return 0
}

# Validate no stale agents exist
validate_no_stale_agents() {
    log_info "R610/R611: Validating no stale agents..."

    local completed_agents=$(find_completed_agents)

    if [ -z "$completed_agents" ]; then
        log_success "Validation passed: No stale agents in active_agents"
        return 0
    fi

    local agent_count=$(echo "$completed_agents" | wc -l)
    log_error "Validation failed: Found $agent_count stale agent(s) in active_agents"

    echo "$completed_agents" | while IFS= read -r agent_id; do
        log_warning "  - $agent_id"
    done

    return 1
}

# Main
main() {
    local mode="cleanup"
    local agent_id_filter=""
    local dry_run=false

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --validate)
                mode="validate"
                shift
                ;;
            --agent-id)
                agent_id_filter="$2"
                shift 2
                ;;
            --dry-run)
                dry_run=true
                shift
                ;;
            --help)
                usage
                ;;
            *)
                log_error "Unknown option: $1"
                usage
                ;;
        esac
    done

    # Validate prerequisites
    validate_prerequisites

    # Execute based on mode
    case "$mode" in
        "validate")
            validate_no_stale_agents
            ;;
        "cleanup")
            if [ -n "$agent_id_filter" ]; then
                # Cleanup specific agent
                if [ "$dry_run" = "true" ]; then
                    log_info "[DRY RUN] Would cleanup agent: $agent_id_filter"
                else
                    backup_state_file
                fi
                cleanup_agent "$agent_id_filter" "$dry_run"
            else
                # Cleanup all
                cleanup_all_completed "$dry_run"
            fi
            ;;
    esac
}

# Run main
main "$@"
