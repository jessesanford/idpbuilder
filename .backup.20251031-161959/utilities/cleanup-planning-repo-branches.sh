#!/bin/bash
# Remove any effort branches that accidentally got pushed to the planning repo
# Effort branches belong in TARGET repo only!

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
RESET='\033[0m'

log() {
    local level="$1"
    shift
    echo -e "${level}$*${RESET}"
}

log "${BOLD}${CYAN}" "🧹 CLEANING UP PLANNING REPO BRANCHES"
log "${YELLOW}" "⚠️  Effort branches should be in TARGET repo, not planning repo!"

# Get planning repo URL (current repo)
PLANNING_REPO=$(git remote get-url origin 2>/dev/null)
log "${BLUE}" "Planning Repo: $PLANNING_REPO"

# List all remote branches that look like effort branches
log "${CYAN}" "\n🔍 Searching for effort branches in planning repo..."

EFFORT_BRANCHES=$(git ls-remote --heads origin | grep -E "phase[0-9]/wave[0-9]/(.*-split-|image-|gitea-|registry-|cli-|buildah-|dockerfile-)" | cut -f2 | sed 's|refs/heads/||' || true)

if [ -z "$EFFORT_BRANCHES" ]; then
    log "${GREEN}" "✅ No effort branches found in planning repo - all clean!"
    exit 0
fi

log "${RED}" "❌ Found effort branches that shouldn't be here:"
echo "$EFFORT_BRANCHES" | while read branch; do
    log "${YELLOW}" "  • $branch"
done

# Ask for confirmation
log "${CYAN}" "\n🤔 Do you want to delete these branches from planning repo? (y/N)"
read -r response

if [[ "$response" =~ ^[Yy]$ ]]; then
    echo "$EFFORT_BRANCHES" | while read branch; do
        log "${RED}" "  🗑️  Deleting: $branch"
        git push origin --delete "$branch" 2>/dev/null || log "${YELLOW}" "    Already deleted or doesn't exist"
    done
    log "${GREEN}" "\n✅ Cleanup complete!"
else
    log "${YELLOW}" "\n⚠️  Skipped cleanup. Run again when ready."
fi

log "${CYAN}" "\n📝 Remember to use:"
log "${BLUE}" "  • utilities/push-efforts-from-worktrees.sh - Push from worktrees to TARGET repo"
log "${BLUE}" "  • utilities/restore-all-efforts.sh - Clone efforts from TARGET repo"