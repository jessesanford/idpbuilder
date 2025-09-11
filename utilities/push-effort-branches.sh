#!/bin/bash
# Push effort branches from their worktree directories to remote

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
RESET='\033[0m'

# Configuration
STATE_FILE="${1:-orchestrator-state.json}"
EFFORTS_ROOT="efforts"
PUSH_LOG="effort-push-$(date +%Y%m%d-%H%M%S).log"

# Counters
PUSHED_COUNT=0
SKIPPED_COUNT=0
FAILED_COUNT=0

# Arrays for tracking
declare -a PUSHED_BRANCHES
declare -a SKIPPED_BRANCHES
declare -a FAILED_BRANCHES

# Print colored output
log() {
    local level="$1"
    shift
    echo -e "${level}$*${RESET}" | tee -a "$PUSH_LOG"
}

# Print header
print_header() {
    echo "═══════════════════════════════════════════════════════════════════════════" | tee "$PUSH_LOG"
    log "${BOLD}${CYAN}" "🏭 SOFTWARE FACTORY 2.0 - EFFORT BRANCH PUSH UTILITY"
    echo "═══════════════════════════════════════════════════════════════════════════" | tee -a "$PUSH_LOG"
    log "${CYAN}" "Timestamp: $(date -u '+%Y-%m-%d %H:%M:%S UTC')"
    echo "═══════════════════════════════════════════════════════════════════════════" | tee -a "$PUSH_LOG"
}

# Check prerequisites
check_prerequisites() {
    log "${CYAN}" "\n📋 Checking prerequisites..."
    
    # Check state file exists
    if [ ! -f "$STATE_FILE" ]; then
        log "${RED}" "❌ ERROR: State file not found: $STATE_FILE"
        exit 1
    fi
    
    # Check efforts directory exists
    if [ ! -d "$EFFORTS_ROOT" ]; then
        log "${RED}" "❌ ERROR: Efforts directory not found: $EFFORTS_ROOT"
        exit 1
    fi
    
    log "${GREEN}" "✅ All prerequisites met"
}

# Push branch from worktree
push_from_worktree() {
    local worktree_dir="$1"
    local expected_branch="$2"
    
    log "${CYAN}" "\n🔄 Processing: ${worktree_dir}"
    
    # Check if directory exists
    if [ ! -d "$worktree_dir" ]; then
        log "${YELLOW}" "  ⚠️  Directory does not exist: $worktree_dir"
        SKIPPED_BRANCHES+=("$expected_branch (directory not found)")
        ((SKIPPED_COUNT++))
        return 1
    fi
    
    # Check if it's a git repository
    if [ ! -d "$worktree_dir/.git" ]; then
        log "${YELLOW}" "  ⚠️  Not a git repository: $worktree_dir"
        SKIPPED_BRANCHES+=("$expected_branch (not a git repo)")
        ((SKIPPED_COUNT++))
        return 1
    fi
    
    # Enter the worktree
    cd "$worktree_dir"
    
    # Get current branch
    local current_branch=$(git branch --show-current)
    log "${BLUE}" "  Current branch: ${current_branch}"
    
    # Check if it matches expected branch
    if [ "$current_branch" != "$expected_branch" ]; then
        log "${YELLOW}" "  ⚠️  Branch mismatch. Expected: $expected_branch, Got: $current_branch"
        # Don't fail, just note it
    fi
    
    # Check for uncommitted changes
    if [ -n "$(git status --porcelain)" ]; then
        log "${YELLOW}" "  ⚠️  Uncommitted changes detected"
        log "${BLUE}" "  Committing changes..."
        git add -A
        git commit -m "chore: auto-commit before push - $(date)" || {
            log "${YELLOW}" "  No changes to commit"
        }
    fi
    
    # Check if branch exists in remote
    if git ls-remote --heads origin "$current_branch" 2>/dev/null | grep -q "$current_branch"; then
        log "${BLUE}" "  Branch exists in remote, checking for updates..."
        
        # Fetch latest
        git fetch origin "$current_branch" 2>&1 | tee -a "$PUSH_LOG"
        
        # Check if we need to push
        local ahead=$(git rev-list --count "origin/$current_branch..HEAD" 2>/dev/null || echo "0")
        if [ "$ahead" -gt 0 ]; then
            log "${GREEN}" "  📤 Pushing $ahead commits..."
            if git push origin "$current_branch" 2>&1 | tee -a "$PUSH_LOG"; then
                log "${GREEN}" "  ✅ Successfully pushed updates"
                PUSHED_BRANCHES+=("$current_branch")
                ((PUSHED_COUNT++))
            else
                log "${RED}" "  ❌ Push failed"
                FAILED_BRANCHES+=("$current_branch")
                ((FAILED_COUNT++))
            fi
        else
            log "${BLUE}" "  ℹ️  Already up to date"
            SKIPPED_BRANCHES+=("$current_branch (up to date)")
            ((SKIPPED_COUNT++))
        fi
    else
        log "${GREEN}" "  📤 Pushing new branch to remote..."
        if git push -u origin "$current_branch" 2>&1 | tee -a "$PUSH_LOG"; then
            log "${GREEN}" "  ✅ Successfully pushed new branch"
            PUSHED_BRANCHES+=("$current_branch")
            ((PUSHED_COUNT++))
        else
            log "${RED}" "  ❌ Push failed"
            FAILED_BRANCHES+=("$current_branch")
            ((FAILED_COUNT++))
        fi
    fi
    
    # Return to original directory
    cd - > /dev/null
}

# Extract efforts from state file
process_efforts_from_state() {
    log "${CYAN}" "\n📊 Processing efforts from state file..."
    
    # Get completed efforts
    local completed_efforts=$(jq -r '.efforts_completed[] | "\(.working_dir):\(.branch)"' "$STATE_FILE" 2>/dev/null || true)
    
    # Get in-progress efforts  
    local in_progress=$(jq -r '.efforts_in_progress[] | "\(.working_dir):\(.branch)"' "$STATE_FILE" 2>/dev/null || true)
    
    # Get splits
    local splits=$(jq -r '.split_tracking[]?.splits[] | "\(.working_dir):\(.branch)"' "$STATE_FILE" 2>/dev/null || true)
    
    # Combine all
    local all_efforts=$(echo -e "$completed_efforts\n$in_progress\n$splits" | grep -v '^$' | sort -u)
    
    if [ -z "$all_efforts" ]; then
        log "${YELLOW}" "No efforts found in state file"
        return
    fi
    
    local effort_count=$(echo "$all_efforts" | wc -l)
    log "${BLUE}" "Found ${effort_count} efforts to process"
    
    # Process each effort
    echo "$all_efforts" | while IFS=':' read -r working_dir branch; do
        if [ -n "$working_dir" ] && [ -n "$branch" ]; then
            push_from_worktree "$working_dir" "$branch"
        fi
    done
}

# Generate summary
generate_summary() {
    log "${BOLD}${CYAN}" "\n═══════════════════════════════════════════════════════════════════════════"
    log "${BOLD}${CYAN}" "📊 PUSH SUMMARY"
    log "${BOLD}${CYAN}" "═══════════════════════════════════════════════════════════════════════════"
    
    log "${GREEN}" "✅ Pushed: ${PUSHED_COUNT}"
    if [ ${#PUSHED_BRANCHES[@]} -gt 0 ]; then
        for branch in "${PUSHED_BRANCHES[@]}"; do
            log "${GREEN}" "    • ${branch}"
        done
    fi
    
    log "${YELLOW}" "⏭️  Skipped: ${SKIPPED_COUNT}"
    if [ ${#SKIPPED_BRANCHES[@]} -gt 0 ]; then
        for branch in "${SKIPPED_BRANCHES[@]}"; do
            log "${YELLOW}" "    • ${branch}"
        done
    fi
    
    log "${RED}" "❌ Failed: ${FAILED_COUNT}"
    if [ ${#FAILED_BRANCHES[@]} -gt 0 ]; then
        for branch in "${FAILED_BRANCHES[@]}"; do
            log "${RED}" "    • ${branch}"
        done
    fi
    
    log "${CYAN}" "\n📝 Full log saved to: ${PUSH_LOG}"
    log "${CYAN}" "⏰ Completed at: $(date -u '+%Y-%m-%d %H:%M:%S UTC')"
    
    # Provide next steps
    if [ ${#PUSHED_BRANCHES[@]} -gt 0 ]; then
        log "${BOLD}${GREEN}" "\n✨ Next Steps:"
        log "${GREEN}" "  1. Verify branches are visible in GitHub/GitLab"
        log "${GREEN}" "  2. Run restore-all-efforts.sh from another location to test restoration"
    fi
}

# Main execution
main() {
    print_header
    check_prerequisites
    
    # Save current directory
    ORIGINAL_DIR=$(pwd)
    
    # Process efforts
    process_efforts_from_state
    
    # Return to original directory
    cd "$ORIGINAL_DIR"
    
    # Generate summary
    generate_summary
    
    # Exit with appropriate code
    if [ $FAILED_COUNT -gt 0 ]; then
        exit 1
    fi
    exit 0
}

# Run main
main "$@"