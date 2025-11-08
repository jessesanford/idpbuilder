#!/bin/bash
################################################################################
# SOFTWARE FACTORY 2.0 - TODO PRESERVATION UTILITY
################################################################################
#
# This utility provides functions for preserving and restoring TODO state
# across context compactions and agent state transitions.
#
# FEATURES:
# - Save current TODO state to persistent files
# - Load TODO state from preserved files  
# - Merge TODO lists while avoiding duplicates
# - Clean up old TODO files
# - Validate TODO state integrity
#
# USAGE:
#   ./todo-preservation.sh save [agent-name] [state]
#   ./todo-preservation.sh load [agent-name] [state]
#   ./todo-preservation.sh merge [file1] [file2]
#   ./todo-preservation.sh cleanup [agent-name]
#   ./todo-preservation.sh validate [file]
#
################################################################################

set -euo pipefail

# Constants
readonly SCRIPT_NAME="$(basename "$0")"
readonly HOOKS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Use CLAUDE_PROJECT_DIR if set, otherwise use current project root
readonly PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/workspaces/software-factory-2.0-template}"
readonly TODOS_BASE_DIR="${PROJECT_DIR}/todos"
readonly BACKUP_DIR="/tmp/todos-backup"
readonly LOG_FILE="/tmp/todo-preservation.log"

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m' 
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Logging function
log() {
    local level="$1"
    shift
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$level] $*" | tee -a "$LOG_FILE"
}

# Output functions
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Error handling
error_exit() {
    print_error "$1"
    exit 1
}

################################################################################
# USAGE AND HELP
################################################################################

usage() {
    cat << EOF
SOFTWARE FACTORY 2.0 - TODO PRESERVATION UTILITY

USAGE:
    $SCRIPT_NAME <command> [arguments]

COMMANDS:
    save <agent> <state>     Save current TODO state for agent
    load <agent> <state>     Load TODO state for agent
    merge <file1> <file2>    Merge two TODO files
    cleanup <agent>          Clean old TODO files for agent
    validate <file>          Validate TODO file format
    list [agent]             List TODO files
    help                     Show this help message

AGENT TYPES:
    orchestrator             orchestrator
    sw-engineer              sw-engineer  
    code-reviewer            code-reviewer
    architect                architect

STATE EXAMPLES:
    INIT, WAVE_START, WAVE_COMPLETE, IMPLEMENTATION, CODE_REVIEW, etc.

EXAMPLES:
    $SCRIPT_NAME save orchestrator WAVE_COMPLETE
    $SCRIPT_NAME load sw-engineer IMPLEMENTATION
    $SCRIPT_NAME cleanup code-reviewer
    $SCRIPT_NAME list orchestrator

EOF
}

################################################################################
# CORE TODO PRESERVATION FUNCTIONS
################################################################################

save_todo_state() {
    local agent="$1"
    local state="$2"
    
    log "INFO" "Saving TODO state for agent: $agent, state: $state"
    
    # Validate agent type
    validate_agent_type "$agent"
    
    # Create todos directory if it doesn't exist
    mkdir -p "$TODOS_BASE_DIR"
    
    # Generate filename with timestamp
    local timestamp=$(date '+%Y%m%d-%H%M%S')
    local todo_file="$TODOS_BASE_DIR/${agent}-${state}-${timestamp}.todo"
    
    print_info "Creating TODO state file: $(basename "$todo_file")"
    
    # Create TODO file with header
    cat > "$todo_file" << EOF
################################################################################
# SOFTWARE FACTORY 2.0 - TODO STATE PRESERVATION
################################################################################
# Agent: $agent
# State: $state
# Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')
# Working Directory: $(pwd)
# Git Branch: $(git branch --show-current 2>/dev/null || echo 'not-a-repo')
# Git Commit: $(git rev-parse HEAD 2>/dev/null || echo 'no-commit')
################################################################################

# INSTRUCTIONS FOR RECOVERY:
# 1. Use Read tool to examine this file
# 2. Use TodoWrite tool to load these items into your working TODO list
# 3. Set appropriate status for each item (pending/in_progress/completed)
# 4. Deduplicate any overlapping tasks

################################################################################
# TODO ITEMS
################################################################################

EOF
    
    # Try to extract current TODO state from various sources
    extract_todo_items "$todo_file"
    
    # Make the file readable and commit it
    chmod 644 "$todo_file"
    
    # Commit to version control
    commit_todo_file "$todo_file" "$agent" "$state"
    
    print_success "TODO state saved: $todo_file"
    return 0
}

extract_todo_items() {
    local todo_file="$1"
    
    # Add placeholder structure for manual TODO entry
    cat >> "$todo_file" << 'EOF'
## PENDING TASKS:
# Add your pending tasks here, one per line starting with "- "
# Example: - Complete Wave 2 webhook framework implementation
# Example: - Review split compliance for effort-auth-api

## IN PROGRESS TASKS:
# Add your currently active tasks here
# Example: - Implementing logical cluster fields for API types

## COMPLETED TASKS:
# Add your completed tasks for context
# Example: - Created IMPLEMENTATION-PLAN.md for Wave 1

## BLOCKED TASKS:
# Add any tasks that are waiting for dependencies
# Example: - Integration testing (blocked by code review completion)

## HIGH PRIORITY ISSUES:
# Add any critical issues that need immediate attention
# Example: - Fix architect feedback on multi-tenancy implementation

################################################################################
# CONTEXT NOTES
################################################################################

## CURRENT FOCUS:
# Describe what you were working on when this state was saved

## NEXT STEPS:
# List the immediate next actions to take

## DEPENDENCIES:
# Note what you're waiting for from other agents

## WARNINGS:
# Any critical issues or constraints to remember

EOF

    print_info "TODO template structure created. You may need to populate manually."
}

load_todo_state() {
    local agent="$1"
    local state="${2:-latest}"
    
    log "INFO" "Loading TODO state for agent: $agent, state: $state"
    
    # Validate agent type
    validate_agent_type "$agent"
    
    # Find the appropriate TODO file
    local todo_file
    if [ "$state" = "latest" ]; then
        todo_file=$(find "$TODOS_BASE_DIR" -name "${agent}-*.todo" -type f -printf '%T@ %p\n' 2>/dev/null | sort -nr | head -1 | cut -d' ' -f2-)
    else
        todo_file=$(find "$TODOS_BASE_DIR" -name "${agent}-${state}-*.todo" -type f -printf '%T@ %p\n' 2>/dev/null | sort -nr | head -1 | cut -d' ' -f2-)
    fi
    
    if [ -z "$todo_file" ] || [ ! -f "$todo_file" ]; then
        print_error "No TODO file found for agent: $agent, state: $state"
        print_info "Available files:"
        list_todo_files "$agent"
        return 1
    fi
    
    print_success "Found TODO file: $(basename "$todo_file")"
    
    # Display the file content for review
    echo
    print_info "TODO STATE CONTENT:"
    echo "----------------------------------------"
    cat "$todo_file"
    echo "----------------------------------------"
    echo
    
    print_warning "CRITICAL: You must now use the TodoWrite tool to load these items!"
    print_info "Steps to complete recovery:"
    echo "1. Parse the TODO items from the content above"
    echo "2. Use TodoWrite tool to populate your working TODO list"
    echo "3. Set appropriate status for each item"
    echo "4. Verify all tasks are loaded correctly"
    
    return 0
}

merge_todo_files() {
    local file1="$1"
    local file2="$2"
    
    log "INFO" "Merging TODO files: $file1 and $file2"
    
    if [ ! -f "$file1" ] || [ ! -f "$file2" ]; then
        error_exit "Both files must exist for merging"
    fi
    
    local merged_file="$TODOS_BASE_DIR/merged-$(date '+%Y%m%d-%H%M%S').todo"
    
    # Create merged file with header
    cat > "$merged_file" << EOF
################################################################################
# SOFTWARE FACTORY 2.0 - MERGED TODO STATE
################################################################################
# Merged from:
#   - $(basename "$file1")
#   - $(basename "$file2")
# Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')
################################################################################

EOF
    
    # Extract content from both files (skip headers)
    echo "# FROM FILE 1: $(basename "$file1")" >> "$merged_file"
    echo "################################################################################" >> "$merged_file"
    grep -E '^(- |## |#)' "$file1" | grep -v '^################################################################################' >> "$merged_file"
    echo "" >> "$merged_file"
    
    echo "# FROM FILE 2: $(basename "$file2")" >> "$merged_file"
    echo "################################################################################" >> "$merged_file"
    grep -E '^(- |## |#)' "$file2" | grep -v '^################################################################################' >> "$merged_file"
    echo "" >> "$merged_file"
    
    echo "## MERGE NOTES:" >> "$merged_file"
    echo "# Review for duplicate tasks and consolidate as needed" >> "$merged_file"
    echo "# Use TodoWrite tool to load the consolidated list" >> "$merged_file"
    
    print_success "Merged TODO file created: $merged_file"
    print_warning "Review the merged file for duplicates before loading"
    
    return 0
}

cleanup_todo_files() {
    local agent="$1"
    
    log "INFO" "Cleaning up old TODO files for agent: $agent"
    
    # Validate agent type
    validate_agent_type "$agent"
    
    # Find all TODO files for this agent
    local todo_files=($(find "$TODOS_BASE_DIR" -name "${agent}-*.todo" -type f -printf '%T@ %p\n' 2>/dev/null | sort -nr | cut -d' ' -f2-))
    
    local total_files=${#todo_files[@]}
    
    if [ $total_files -eq 0 ]; then
        print_info "No TODO files found for agent: $agent"
        return 0
    fi
    
    print_info "Found $total_files TODO files for agent: $agent"
    
    # Keep the 5 most recent files, remove older ones
    local keep_count=5
    
    if [ $total_files -le $keep_count ]; then
        print_info "All files are recent, no cleanup needed"
        return 0
    fi
    
    local files_to_remove=(${todo_files[@]:$keep_count})
    
    print_warning "Will remove ${#files_to_remove[@]} old files (keeping $keep_count most recent)"
    
    for file in "${files_to_remove[@]}"; do
        echo "  - $(basename "$file")"
    done
    
    read -p "Proceed with cleanup? (y/N): " -r
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        for file in "${files_to_remove[@]}"; do
            rm -f "$file"
            log "INFO" "Removed old TODO file: $(basename "$file")"
        done
        print_success "Cleanup complete"
    else
        print_info "Cleanup cancelled"
    fi
    
    return 0
}

validate_todo_file() {
    local todo_file="$1"
    
    log "INFO" "Validating TODO file: $todo_file"
    
    if [ ! -f "$todo_file" ]; then
        print_error "File does not exist: $todo_file"
        return 1
    fi
    
    local valid=true
    
    # Check for required header
    if ! grep -q "SOFTWARE FACTORY 2.0" "$todo_file"; then
        print_error "Missing Software Factory header"
        valid=false
    fi
    
    # Check for timestamp
    if ! grep -q "Timestamp:" "$todo_file"; then
        print_warning "Missing timestamp information"
    fi
    
    # Check for TODO structure
    if ! grep -q "## PENDING TASKS:" "$todo_file" && ! grep -q "- " "$todo_file"; then
        print_warning "No TODO items found in file"
    fi
    
    # Check file size
    local size=$(stat -c%s "$todo_file")
    if [ $size -lt 100 ]; then
        print_warning "File seems very small (${size} bytes)"
    fi
    
    if [ "$valid" = true ]; then
        print_success "TODO file validation passed"
        return 0
    else
        print_error "TODO file validation failed"
        return 1
    fi
}

list_todo_files() {
    local agent="${1:-}"
    
    if [ -n "$agent" ]; then
        validate_agent_type "$agent"
        print_info "TODO files for agent: $agent"
        local pattern="${agent}-*.todo"
    else
        print_info "All TODO files:"
        local pattern="*.todo"
    fi
    
    if [ ! -d "$TODOS_BASE_DIR" ]; then
        print_warning "TODOs directory does not exist: $TODOS_BASE_DIR"
        return 1
    fi
    
    local files=($(find "$TODOS_BASE_DIR" -name "$pattern" -type f -printf '%T@ %p\n' 2>/dev/null | sort -nr | cut -d' ' -f2-))
    
    if [ ${#files[@]} -eq 0 ]; then
        print_info "No TODO files found matching pattern: $pattern"
        return 0
    fi
    
    for file in "${files[@]}"; do
        local mod_time=$(stat -c '%Y' "$file" 2>/dev/null)
        local readable_time=$(date -d "@$mod_time" '+%Y-%m-%d %H:%M:%S' 2>/dev/null)
        local size=$(stat -c%s "$file" 2>/dev/null)
        printf "  %-40s %s (%s bytes)\n" "$(basename "$file")" "$readable_time" "$size"
    done
    
    return 0
}

################################################################################
# UTILITY FUNCTIONS
################################################################################

validate_agent_type() {
    local agent="$1"
    
    case "$agent" in
        "orchestrator"|"sw-engineer"|"code-reviewer"|"architect")
            return 0
            ;;
        *)
            error_exit "Invalid agent type: $agent. Must be one of: orchestrator, sw-engineer, code-reviewer, architect"
            ;;
    esac
}

commit_todo_file() {
    local todo_file="$1"
    local agent="$2"
    local state="$3"
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        print_warning "Not in a git repository, skipping commit"
        return 0
    fi
    
    # Navigate to the base directory for git operations
    local current_dir="$(pwd)"
    cd "${PROJECT_DIR}"
    
    # Add the TODO file
    git add "$todo_file"
    
    # Commit with descriptive message
    local commit_msg="todo: save ${agent} state at ${state} - $(date '+%Y-%m-%d %H:%M:%S')"
    git commit -m "$commit_msg" || {
        print_warning "Git commit failed, but file was saved"
        cd "$current_dir"
        return 0
    }
    
    # Try to push if remote is configured
    if git remote get-url origin > /dev/null 2>&1; then
        git push || {
            print_warning "Git push failed, but commit was successful"
        }
    fi
    
    cd "$current_dir"
    print_info "TODO file committed to version control"
}

################################################################################
# MAIN COMMAND DISPATCH
################################################################################

main() {
    # Ensure todos directory exists
    mkdir -p "$TODOS_BASE_DIR"
    
    # Parse command
    local command="${1:-help}"
    
    case "$command" in
        "save")
            if [ $# -lt 3 ]; then
                error_exit "Usage: $SCRIPT_NAME save <agent> <state>"
            fi
            save_todo_state "$2" "$3"
            ;;
        "load")
            if [ $# -lt 2 ]; then
                error_exit "Usage: $SCRIPT_NAME load <agent> [state]"
            fi
            load_todo_state "$2" "${3:-latest}"
            ;;
        "merge")
            if [ $# -lt 3 ]; then
                error_exit "Usage: $SCRIPT_NAME merge <file1> <file2>"
            fi
            merge_todo_files "$2" "$3"
            ;;
        "cleanup")
            if [ $# -lt 2 ]; then
                error_exit "Usage: $SCRIPT_NAME cleanup <agent>"
            fi
            cleanup_todo_files "$2"
            ;;
        "validate")
            if [ $# -lt 2 ]; then
                error_exit "Usage: $SCRIPT_NAME validate <file>"
            fi
            validate_todo_file "$2"
            ;;
        "list")
            list_todo_files "${2:-}"
            ;;
        "help"|"--help"|"-h")
            usage
            ;;
        *)
            print_error "Unknown command: $command"
            echo
            usage
            exit 1
            ;;
    esac
}

################################################################################
# SCRIPT EXECUTION
################################################################################

# Ensure we're being called in the right context
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    main "$@"
fi