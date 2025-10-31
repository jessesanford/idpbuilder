#!/bin/bash
# Fix all effort worktree remotes to point to the TARGET repository

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
RESET='\033[0m'

# Configuration - TARGET REPO is where code goes!
# Try to read from orchestrator-state-v3.json, fallback to command-line argument, then error
STATE_FILE="orchestrator-state-v3.json"
if [ -f "$STATE_FILE" ] && command -v jq &> /dev/null; then
    TARGET_REPO=$(jq -r '.pre_planned_infrastructure.target_repo_url // empty' "$STATE_FILE")
fi

# If not found in state file, check command-line argument
if [ -z "$TARGET_REPO" ]; then
    if [ -n "$1" ]; then
        TARGET_REPO="$1"
    else
        echo -e "${RED}ERROR: TARGET_REPO not found!${RESET}"
        echo "Please either:"
        echo "  1. Add target_repo_url to orchestrator-state-v3.json (pre_planned_infrastructure.target_repo_url)"
        echo "  2. Pass repository URL as argument: $0 <repo-url>"
        echo ""
        echo "Example: $0 https://github.com/username/project.git"
        exit 1
    fi
fi

EFFORTS_ROOT="efforts"

log() {
    local level="$1"
    shift
    echo -e "${level}$*${RESET}"
}

log "${BOLD}${CYAN}" "🔧 FIXING EFFORT REMOTES TO TARGET REPOSITORY"
log "${BLUE}" "Target Repository: ${TARGET_REPO}"
log "${YELLOW}" "⚠️  Efforts belong in TARGET repo, NOT planning repo!"

# Find all effort directories with .git
find "$EFFORTS_ROOT" -name ".git" -type d | while read git_dir; do
    effort_dir=$(dirname "$git_dir")
    
    log "${CYAN}" "\n📁 Processing: $effort_dir"
    
    cd "$effort_dir"
    
    # Get current remotes
    current_origin=$(git remote get-url origin 2>/dev/null || echo "none")
    log "${BLUE}" "  Current origin: $current_origin"
    
    # Fix if not pointing to TARGET repo
    if [[ "$current_origin" == "$TARGET_REPO" ]]; then
        log "${GREEN}" "  ✅ Already pointing to TARGET repo"
    elif [[ "$current_origin" == "none" ]]; then
        log "${YELLOW}" "  ⚠️  No origin set, adding TARGET repo"
        git remote add origin "$TARGET_REPO"
        log "${GREEN}" "  ✅ Added origin to TARGET repo"
    else
        log "${YELLOW}" "  ⚠️  Wrong remote detected, fixing..."
        git remote set-url origin "$TARGET_REPO"
        log "${GREEN}" "  ✅ Fixed origin to TARGET repo"
    fi
    
    # Fix integration remote if exists
    if git remote | grep -q integration; then
        git remote set-url integration "$TARGET_REPO"
        log "${GREEN}" "  ✅ Fixed integration remote"
    fi
    
    # Show current branch
    branch=$(git branch --show-current)
    log "${BLUE}" "  Branch: $branch"
    
    cd - > /dev/null
done

log "${BOLD}${GREEN}" "\n✅ All effort remotes fixed to TARGET repository!"
log "${CYAN}" "Now push with: ./utilities/push-effort-branches.sh"