#!/bin/bash
# Push effort branches from their individual worktree directories to TARGET repo
# IMPORTANT: Efforts are in separate worktrees, NOT in the main repo!

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
RESET='\033[0m'

# Configuration
STATE_FILE="${1:-orchestrator-state.json}"
EFFORTS_ROOT="efforts"

log() {
    local level="$1"
    shift
    echo -e "${level}$*${RESET}"
}

log "${BOLD}${CYAN}" "🏭 SOFTWARE FACTORY 2.0 - EFFORT WORKTREE PUSH UTILITY"
log "${YELLOW}" "⚠️  This pushes from WORKTREE directories to TARGET repo"
log "${YELLOW}" "⚠️  Efforts should NOT exist as branches in the planning repo!"

if [ ! -f "$STATE_FILE" ]; then
    log "${RED}" "❌ State file not found: $STATE_FILE"
    exit 1
fi

# Get target repository from state file
TARGET_REPO=$(jq -r '.project_info.repository // .project_info.target_repository // ""' "$STATE_FILE")
if [ -z "$TARGET_REPO" ] || [ "$TARGET_REPO" == "null" ]; then
    log "${RED}" "❌ Could not determine target repository from state file"
    exit 1
fi

log "${BLUE}" "Target Repository: $TARGET_REPO"
log "${CYAN}" "\n📊 Processing efforts from state file..."

# Process completed efforts
jq -r '.efforts_completed[]? | "\(.working_dir // .workdir // "")"' "$STATE_FILE" | grep -v '^$' | while read workdir; do
    if [ -d "$workdir" ] && [ -d "$workdir/.git" ]; then
        log "${CYAN}" "\n📁 Processing: $workdir"
        cd "$workdir"
        
        branch=$(git branch --show-current)
        log "${BLUE}" "  Branch: $branch"
        
        # Check remote URL
        current_remote=$(git remote get-url origin 2>/dev/null || echo "none")
        if [[ "$current_remote" != "$TARGET_REPO" ]]; then
            log "${YELLOW}" "  ⚠️  Fixing remote URL..."
            git remote set-url origin "$TARGET_REPO"
        fi
        
        # Push to target repo
        log "${GREEN}" "  📤 Pushing to TARGET repo..."
        if git push -u origin HEAD 2>&1; then
            log "${GREEN}" "  ✅ Successfully pushed"
        else
            log "${YELLOW}" "  ℹ️  Already up to date or no changes"
        fi
        
        cd - > /dev/null
    else
        log "${YELLOW}" "  ⚠️  Skipping $workdir (not a git repo)"
    fi
done

# Process splits
jq -r '.split_tracking[]?.splits[]? | "\(.working_dir // "")"' "$STATE_FILE" | grep -v '^$' | while read workdir; do
    if [ -d "$workdir" ] && [ -d "$workdir/.git" ]; then
        log "${CYAN}" "\n📁 Processing split: $workdir"
        cd "$workdir"
        
        branch=$(git branch --show-current)
        log "${BLUE}" "  Branch: $branch"
        
        # Check remote URL
        current_remote=$(git remote get-url origin 2>/dev/null || echo "none")
        if [[ "$current_remote" != "$TARGET_REPO" ]]; then
            log "${YELLOW}" "  ⚠️  Fixing remote URL..."
            git remote set-url origin "$TARGET_REPO"
        fi
        
        # Push to target repo
        log "${GREEN}" "  📤 Pushing to TARGET repo..."
        if git push -u origin HEAD 2>&1; then
            log "${GREEN}" "  ✅ Successfully pushed"
        else
            log "${YELLOW}" "  ℹ️  Already up to date or no changes"
        fi
        
        cd - > /dev/null
    fi
done

log "${BOLD}${GREEN}" "\n✅ Done! All efforts pushed to TARGET repository"
log "${CYAN}" "Remember: Efforts belong in TARGET repo, not planning repo!"