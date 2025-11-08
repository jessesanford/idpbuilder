#!/bin/bash
# Validates that effort is on latest base branch
# Usage: validate-fresh-base.sh <effort_dir> <base_branch>
#
# Part of R614 - Fresh Base Pull Before Agent Spawn Protocol
# Validates that an effort directory has pulled latest from its base branch

set -e

EFFORT_DIR="$1"
BASE_BRANCH="$2"

if [ -z "$EFFORT_DIR" ] || [ -z "$BASE_BRANCH" ]; then
    echo "Usage: $0 <effort_dir> <base_branch>"
    echo ""
    echo "Example:"
    echo "  $0 efforts/phase2/wave2/effort-2.2.2 effort-2.2.1"
    exit 1
fi

# Verify effort directory exists
if [ ! -d "$EFFORT_DIR" ]; then
    echo "❌ ERROR: Effort directory does not exist: $EFFORT_DIR"
    exit 1
fi

# CD to effort directory
cd "$EFFORT_DIR" || {
    echo "❌ ERROR: Cannot CD to $EFFORT_DIR"
    exit 1
}

echo "🔍 Validating fresh base for: $(basename $EFFORT_DIR)"
echo "   Base branch: $BASE_BRANCH"

# Fetch latest from origin
echo "   Fetching latest from origin..."
git fetch origin "$BASE_BRANCH" 2>&1 | grep -v "^From" || {
    echo "❌ ERROR: Cannot fetch $BASE_BRANCH from origin"
    exit 1
}

# Check if base branch exists on remote
if ! git ls-remote --heads origin "$BASE_BRANCH" | grep -q "$BASE_BRANCH"; then
    echo "⚠️  WARNING: Base branch $BASE_BRANCH does not exist on remote"
    echo "   This may be expected for first effort in wave"
    exit 0
fi

# Get remote HEAD and merge-base
LATEST=$(git rev-parse origin/$BASE_BRANCH 2>/dev/null || echo "")
if [ -z "$LATEST" ]; then
    echo "❌ ERROR: Cannot resolve origin/$BASE_BRANCH"
    exit 1
fi

CURRENT_BASE=$(git merge-base HEAD origin/$BASE_BRANCH 2>/dev/null || echo "")
if [ -z "$CURRENT_BASE" ]; then
    echo "❌ ERROR: Cannot find merge-base with origin/$BASE_BRANCH"
    exit 1
fi

# Compare
if [[ "$LATEST" == "$CURRENT_BASE" ]]; then
    echo "   ✅ Effort is on latest base"
    echo "   Base branch: $BASE_BRANCH"
    echo "   Latest commit: $(git rev-parse --short $LATEST)"
    echo "   Commit message: $(git log -1 --format=%s $LATEST)"
    exit 0
else
    echo "   ❌ Effort is behind base branch!"
    echo "   Base branch: $BASE_BRANCH"
    echo "   Latest commit: $(git rev-parse --short $LATEST)"
    echo "   Our merge base: $(git rev-parse --short $CURRENT_BASE)"
    echo ""
    echo "   Behind by: $(git rev-list --count $CURRENT_BASE..$LATEST) commits"
    echo ""
    echo "   Missing commits:"
    git log --oneline $CURRENT_BASE..$LATEST | head -10
    echo ""
    echo "   ⚠️  This violates R614 - effort must pull latest before agent spawn"
    echo "   Run: cd $EFFORT_DIR && git pull origin $BASE_BRANCH"
    exit 1
fi
