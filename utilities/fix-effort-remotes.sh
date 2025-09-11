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
TARGET_REPO="https://github.com/jessesanford/idpbuilder.git"
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
    
    # Fix if pointing to planning repo
    if [[ "$current_origin" == *"idpbuilder-oci-build-push"* ]]; then
        log "${YELLOW}" "  ⚠️  WRONG REPO! Fixing..."
        git remote set-url origin "$TARGET_REPO"
        log "${GREEN}" "  ✅ Fixed origin to TARGET repo"
    elif [[ "$current_origin" == *"idpbuilder.git"* ]]; then
        log "${GREEN}" "  ✅ Already pointing to TARGET repo"
    else
        log "${YELLOW}" "  ⚠️  Unknown remote, setting to TARGET"
        git remote set-url origin "$TARGET_REPO"
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