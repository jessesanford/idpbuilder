#!/bin/bash

# 🏭 SOFTWARE FACTORY 2.0 - LOCAL EFFORT PUSH UTILITY
# ═══════════════════════════════════════════════════════════════════════════
# Purpose: Push local effort branches from worktrees to the remote TARGET repository
# 
# IMPORTANT: Software Factory 2.0 dual-repository architecture:
#   • Planning Repository: orchestrator state, rules, agents
#   • Target Repository: actual code implementations on branches
#
# This script helps push local worktree branches to the remote target repository.
#
# Usage: ./push-local-efforts.sh [orchestrator-state.json] [target-worktree-path]
# 
# Example: ./push-local-efforts.sh orchestrator-state.json /home/vscode/workspaces/idpbuilder-oci-build-push
# ═══════════════════════════════════════════════════════════════════════════

set -euo pipefail

# ANSI Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Default values
STATE_FILE="${1:-orchestrator-state.json}"
TARGET_WORKTREE="${2:-}"
PUSH_LOG="effort-push.log"
PUSHED_BRANCHES=()
FAILED_BRANCHES=()
SKIPPED_BRANCHES=()

# Timestamp function
timestamp() {
    date '+%Y-%m-%d %H:%M:%S %Z'
}

# Logging function
log() {
    local level="$1"
    shift
    local message="$*"
    echo -e "${level}${message}${NC}" | tee -a "$PUSH_LOG"
}

# Print header
print_header() {
    echo "═══════════════════════════════════════════════════════════════════════════" | tee -a "$PUSH_LOG"
    echo "🏭 SOFTWARE FACTORY 2.0 - LOCAL EFFORT PUSH UTILITY" | tee -a "$PUSH_LOG"
    echo "═══════════════════════════════════════════════════════════════════════════" | tee -a "$PUSH_LOG"
    echo "Timestamp: $(timestamp)" | tee -a "$PUSH_LOG"
    echo "═══════════════════════════════════════════════════════════════════════════" | tee -a "$PUSH_LOG"
}

# Check prerequisites
check_prerequisites() {
    log "${CYAN}" "\n📋 Checking prerequisites..."
    
    # Check if jq is installed
    if ! command -v jq &> /dev/null; then
        log "${RED}" "❌ ERROR: jq is not installed. Please install jq first."
        exit 1
    fi
    
    # Check if state file exists
    if [ ! -f "$STATE_FILE" ]; then
        log "${RED}" "❌ ERROR: State file '$STATE_FILE' not found!"
        exit 1
    fi
    
    # Auto-detect target worktree if not provided
    if [ -z "$TARGET_WORKTREE" ]; then
        log "${YELLOW}" "⚠️  No target worktree path provided. Attempting to detect..."
        
        # Try to get from branch names in state file
        local sample_branch=$(jq -r '.efforts_completed[0].branch // .efforts_in_progress[0].branch // ""' "$STATE_FILE" 2>/dev/null)
        if [ -n "$sample_branch" ] && [ "$sample_branch" != "null" ]; then
            local project_name=$(echo "$sample_branch" | cut -d'/' -f1)
            TARGET_WORKTREE="/home/vscode/workspaces/${project_name}"
            
            if [ -d "$TARGET_WORKTREE" ]; then
                log "${GREEN}" "✅ Detected target worktree: $TARGET_WORKTREE"
            else
                log "${RED}" "❌ ERROR: Detected path does not exist: $TARGET_WORKTREE"
                exit 1
            fi
        else
            log "${RED}" "❌ ERROR: Could not detect target worktree path"
            log "${YELLOW}" "  Please provide as second argument: $0 $STATE_FILE <worktree-path>"
            exit 1
        fi
    fi
    
    # Verify target worktree exists and is a git repo
    if [ ! -d "$TARGET_WORKTREE" ]; then
        log "${RED}" "❌ ERROR: Target worktree does not exist: $TARGET_WORKTREE"
        exit 1
    fi
    
    if [ ! -d "$TARGET_WORKTREE/.git" ]; then
        log "${RED}" "❌ ERROR: Target path is not a git repository: $TARGET_WORKTREE"
        exit 1
    fi
    
    log "${GREEN}" "✅ All prerequisites met"
}

# Extract unique branches from state file
extract_branches_from_state() {
    # Don't log inside this function as output is captured
    local branches=()
    
    # Get completed effort branches
    local completed=$(jq -r '.efforts_completed[]?.branch // empty' "$STATE_FILE" 2>/dev/null | grep -v '^$')
    if [ -n "$completed" ]; then
        branches+=($completed)
    fi
    
    # Get in-progress effort branches
    local in_progress=$(jq -r '.efforts_in_progress[]?.branch // empty' "$STATE_FILE" 2>/dev/null | grep -v '^$')
    if [ -n "$in_progress" ]; then
        branches+=($in_progress)
    fi
    
    # Get split tracking branches
    local splits=$(jq -r '.split_tracking[]?.splits[]?.branch // empty' "$STATE_FILE" 2>/dev/null | grep -v '^$')
    if [ -n "$splits" ]; then
        branches+=($splits)
    fi
    
    # Get integration branches (filter out null and empty)
    local integrations=$(jq -r '.integration_branches[]?.branch | select(. != null and . != "")' "$STATE_FILE" 2>/dev/null | grep -v '^$' || true)
    if [ -n "$integrations" ]; then
        branches+=($integrations)
    fi
    
    # Remove duplicates, empty lines, and sort
    printf '%s\n' "${branches[@]}" | grep -v '^$' | sort -u
}

# Check if branch exists locally
check_local_branch() {
    local branch="$1"
    cd "$TARGET_WORKTREE"
    
    # Check if branch exists locally
    if git show-ref --verify --quiet "refs/heads/$branch"; then
        return 0
    else
        return 1
    fi
}

# Check if branch exists in remote
check_remote_branch() {
    local branch="$1"
    cd "$TARGET_WORKTREE"
    
    # Check if branch exists in remote
    if git ls-remote --heads origin "$branch" 2>/dev/null | grep -q "$branch"; then
        return 0
    else
        return 1
    fi
}

# Push branch to remote
push_branch() {
    local branch="$1"
    
    log "${CYAN}" "\n🔄 Processing branch: ${branch}"
    
    cd "$TARGET_WORKTREE"
    
    # Check if branch exists locally
    if ! check_local_branch "$branch"; then
        log "${YELLOW}" "  ⚠️  Branch does not exist locally: $branch"
        SKIPPED_BRANCHES+=("$branch (not found locally)")
        return 1
    fi
    
    # Switch to the branch
    log "${BLUE}" "  📍 Switching to branch: $branch"
    if ! git checkout "$branch" 2>&1 | tee -a "$PUSH_LOG"; then
        log "${RED}" "  ❌ Failed to checkout branch"
        FAILED_BRANCHES+=("$branch (checkout failed)")
        return 1
    fi
    
    # Check if there are commits to push
    local unpushed=$(git rev-list --count @{u}..HEAD 2>/dev/null || echo "0")
    
    if [ "$unpushed" -eq 0 ]; then
        # Check if branch exists in remote
        if check_remote_branch "$branch"; then
            log "${GREEN}" "  ✅ Branch already up to date in remote"
            SKIPPED_BRANCHES+=("$branch (already up to date)")
        else
            # Branch doesn't exist in remote, push it
            log "${BLUE}" "  📤 Pushing new branch to remote..."
            if git push -u origin "$branch" 2>&1 | tee -a "$PUSH_LOG"; then
                log "${GREEN}" "  ✅ Successfully pushed new branch: $branch"
                PUSHED_BRANCHES+=("$branch (new)")
            else
                log "${RED}" "  ❌ Failed to push branch"
                FAILED_BRANCHES+=("$branch (push failed)")
                return 1
            fi
        fi
    else
        log "${BLUE}" "  📤 Pushing $unpushed commits to remote..."
        if git push origin "$branch" 2>&1 | tee -a "$PUSH_LOG"; then
            log "${GREEN}" "  ✅ Successfully pushed: $branch ($unpushed commits)"
            PUSHED_BRANCHES+=("$branch ($unpushed commits)")
        else
            log "${RED}" "  ❌ Failed to push branch"
            FAILED_BRANCHES+=("$branch (push failed)")
            return 1
        fi
    fi
    
    return 0
}

# Generate summary report
generate_summary() {
    log "${BOLD}${CYAN}" "\n═══════════════════════════════════════════════════════════════════════════"
    log "${BOLD}${CYAN}" "📊 PUSH SUMMARY"
    log "${BOLD}${CYAN}" "═══════════════════════════════════════════════════════════════════════════"
    
    log "${GREEN}" "\n✅ Pushed: ${#PUSHED_BRANCHES[@]} branches"
    for branch in "${PUSHED_BRANCHES[@]}"; do
        log "${GREEN}" "    • ${branch}"
    done
    
    if [ ${#SKIPPED_BRANCHES[@]} -gt 0 ]; then
        log "${YELLOW}" "\n⏭️  Skipped: ${#SKIPPED_BRANCHES[@]} branches"
        for branch in "${SKIPPED_BRANCHES[@]}"; do
            log "${YELLOW}" "    • ${branch}"
        done
    fi
    
    if [ ${#FAILED_BRANCHES[@]} -gt 0 ]; then
        log "${RED}" "\n❌ Failed: ${#FAILED_BRANCHES[@]} branches"
        for branch in "${FAILED_BRANCHES[@]}"; do
            log "${RED}" "    • ${branch}"
        done
    fi
    
    log "${CYAN}" "\n📝 Full log saved to: ${PUSH_LOG}"
    log "${CYAN}" "⏰ Completed at: $(timestamp)"
    
    # Provide next steps
    if [ ${#PUSHED_BRANCHES[@]} -gt 0 ]; then
        log "${BOLD}${GREEN}" "\n✨ Next Steps:"
        log "${GREEN}" "  1. Run restore-all-efforts.sh to clone these branches elsewhere"
        log "${GREEN}" "  2. Verify all branches are accessible from the remote repository"
    fi
}

# Main execution
main() {
    # Initialize log
    echo "Starting effort push at $(timestamp)" > "$PUSH_LOG"
    
    print_header
    check_prerequisites
    
    log "${BOLD}${CYAN}" "\n🏭 Starting effort push from: ${STATE_FILE}"
    log "${BLUE}" "📍 Target worktree: ${TARGET_WORKTREE}"
    
    log "${CYAN}" "\n📋 Extracting branches from state file..."
    
    # Extract branches from state file
    local branches=$(extract_branches_from_state)
    local branch_count=$(echo "$branches" | grep -c . || echo "0")
    
    log "${BLUE}" "\n📊 Found ${branch_count} unique branches in state file"
    
    # Process each branch
    while IFS= read -r branch; do
        if [ -n "$branch" ]; then
            push_branch "$branch"
        fi
    done <<< "$branches"
    
    # Generate summary
    generate_summary
    
    # Exit with appropriate code
    if [ ${#FAILED_BRANCHES[@]} -gt 0 ]; then
        exit 1
    else
        exit 0
    fi
}

# Help function
show_help() {
    cat << EOF
═══════════════════════════════════════════════════════════════════════════
🏭 SOFTWARE FACTORY 2.0 - LOCAL EFFORT PUSH UTILITY
═══════════════════════════════════════════════════════════════════════════

USAGE:
    $0 [orchestrator-state.json] [target-worktree-path]

DESCRIPTION:
    Pushes local effort branches from worktrees to the remote TARGET repository.
    Software Factory 2.0 uses two repositories:
    
    • Planning Repository: orchestrator state, rules, agents
    • Target Repository: actual code implementations on branches
    
    This script ensures all local effort branches are pushed to remote.

ARGUMENTS:
    orchestrator-state.json  Path to the state file (default: orchestrator-state.json)
    target-worktree-path    Path to the target repository worktree
                           (will attempt auto-detection if not provided)

EXAMPLES:
    # Basic usage with auto-detection
    $0
    
    # Specify state file
    $0 my-state.json
    
    # Full specification
    $0 orchestrator-state.json /home/vscode/workspaces/idpbuilder-oci-build-push

FEATURES:
    • Reads orchestrator-state.json to identify all effort branches
    • Checks each branch for local existence
    • Pushes new branches and updates to remote
    • Handles integration branches and split efforts
    • Comprehensive progress reporting
    • Safe handling of already-pushed branches

OUTPUT FILES:
    • effort-push.log - Full push operation log

EXIT CODES:
    0 - All branches pushed successfully (or already up to date)
    1 - One or more branches failed to push

═══════════════════════════════════════════════════════════════════════════
EOF
}

# Parse command line arguments
if [ "$#" -gt 0 ] && [[ "$1" == "-h" || "$1" == "--help" || "$1" == "help" ]]; then
    show_help
    exit 0
fi

# Run main function
main "$@"