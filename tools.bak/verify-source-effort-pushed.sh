#!/bin/bash
# verify-source-effort-pushed.sh
# Verifies that source effort has pushed all changes before dependent can pull
# Part of R614 - Fresh Base Pull Before Agent Spawn
#
# Usage: verify-source-effort-pushed.sh <source_effort_id> [orchestrator_state_file]
#
# Exit codes:
#   0: Source pushed (or doesn't exist locally) - safe to pull
#   614: Source has unpushed/uncommitted changes - BLOCKING
#   1: Invalid arguments or state file error

set -euo pipefail

# Arguments
SOURCE_EFFORT_ID="${1:-}"
STATE_FILE="${2:-orchestrator-state-v3.json}"

# Validate arguments
if [[ -z "$SOURCE_EFFORT_ID" ]]; then
    echo "Usage: $0 <source_effort_id> [orchestrator_state_file]" >&2
    echo "Example: $0 2.2.1 orchestrator-state-v3.json" >&2
    exit 1
fi

if [[ ! -f "$STATE_FILE" ]]; then
    echo "❌ ERROR: State file not found: $STATE_FILE" >&2
    exit 1
fi

echo "🔍 R614: Verifying source effort has pushed all changes"
echo "   Source Effort ID: $SOURCE_EFFORT_ID"
echo "   State File: $STATE_FILE"

# Get source effort directory from state file
SOURCE_DIR=$(jq -r ".pre_planned_infrastructure.efforts[] |
    select(.effort_id == \"$SOURCE_EFFORT_ID\") |
    .effort_dir" "$STATE_FILE" 2>/dev/null)

# Check if effort exists in state
if [[ -z "$SOURCE_DIR" || "$SOURCE_DIR" == "null" ]]; then
    echo "❌ ERROR: Cannot find effort $SOURCE_EFFORT_ID in state file" >&2
    echo "   Checked: $STATE_FILE" >&2
    echo "   Section: .pre_planned_infrastructure.efforts[]" >&2
    exit 1
fi

echo "   Source Directory: $SOURCE_DIR"

# If directory doesn't exist locally, assume it's on remote only
if [[ ! -d "$SOURCE_DIR" ]]; then
    echo "📝 Source effort directory not found locally"
    echo "   This is OK - source may exist on remote only"
    echo "   Scenarios:"
    echo "     - Different iteration/developer"
    echo "     - Previous wave already integrated"
    echo "     - Infrastructure from different machine"
    echo ""
    echo "✅ Assuming source exists on remote - proceeding with pull"
    exit 0
fi

echo "📁 Source effort directory found locally"
echo "   Performing comprehensive push verification..."
echo ""

# Save current directory
ORIGINAL_DIR=$(pwd)

# CD to source effort directory
if ! cd "$SOURCE_DIR" 2>/dev/null; then
    echo "❌ ERROR: Cannot access source effort directory: $SOURCE_DIR" >&2
    exit 1
fi

# Get current branch
CURRENT_BRANCH=$(git branch --show-current)

if [[ -z "$CURRENT_BRANCH" ]]; then
    echo "❌ ERROR: Source effort is in detached HEAD state" >&2
    echo "   Directory: $SOURCE_DIR" >&2
    echo "   This indicates infrastructure corruption" >&2
    cd "$ORIGINAL_DIR"
    exit 614
fi

echo "   Current Branch: $CURRENT_BRANCH"

# Check if remote tracking branch exists
if ! git rev-parse --verify "origin/$CURRENT_BRANCH" >/dev/null 2>&1; then
    echo "❌ CRITICAL: Remote tracking branch does not exist!" >&2
    echo "   Branch: $CURRENT_BRANCH" >&2
    echo "   Remote: origin/$CURRENT_BRANCH" >&2
    echo ""
    echo "🛑 Source branch has never been pushed to origin!" >&2
    echo ""
    echo "Fix:" >&2
    echo "   cd $SOURCE_DIR" >&2
    echo "   git push -u origin $CURRENT_BRANCH" >&2
    echo "" >&2
    cd "$ORIGINAL_DIR"
    exit 614
fi

# Check for unpushed commits
echo "   Checking for unpushed commits..."
UNPUSHED=$(git log "origin/$CURRENT_BRANCH"..HEAD --oneline 2>/dev/null || echo "")

if [[ -n "$UNPUSHED" ]]; then
    echo "" >&2
    echo "❌ R614 CRITICAL VIOLATION: Source effort has UNPUSHED commits!" >&2
    echo "   Source: $SOURCE_DIR" >&2
    echo "   Branch: $CURRENT_BRANCH" >&2
    echo "   Effort ID: $SOURCE_EFFORT_ID" >&2
    echo "" >&2
    echo "Unpushed commits:" >&2
    echo "$UNPUSHED" >&2
    echo "" >&2
    echo "🛑 BLOCKING: Cannot pull from source until it pushes!" >&2
    echo "" >&2
    echo "Fix:" >&2
    echo "   cd $SOURCE_DIR" >&2
    echo "   git push origin $CURRENT_BRANCH" >&2
    echo "" >&2
    cd "$ORIGINAL_DIR"
    exit 614
fi

echo "   ✅ No unpushed commits"

# Check for uncommitted changes
echo "   Checking for uncommitted changes..."
UNCOMMITTED=$(git status --porcelain)

if [[ -n "$UNCOMMITTED" ]]; then
    echo "" >&2
    echo "❌ R614 CRITICAL VIOLATION: Source effort has UNCOMMITTED changes!" >&2
    echo "   Source: $SOURCE_DIR" >&2
    echo "   Branch: $CURRENT_BRANCH" >&2
    echo "   Effort ID: $SOURCE_EFFORT_ID" >&2
    echo "" >&2
    echo "Uncommitted changes:" >&2
    git status --short >&2
    echo "" >&2
    echo "🛑 BLOCKING: Cannot pull from source until it commits and pushes!" >&2
    echo "" >&2
    echo "Fix:" >&2
    echo "   cd $SOURCE_DIR" >&2
    echo "   git add <files>" >&2
    echo "   git commit -m \"...\"" >&2
    echo "   git push origin $CURRENT_BRANCH" >&2
    echo "" >&2
    cd "$ORIGINAL_DIR"
    exit 614
fi

echo "   ✅ No uncommitted changes"

# Verify local and remote are in sync
echo "   Verifying local matches remote..."
LOCAL_COMMIT=$(git rev-parse HEAD)
REMOTE_COMMIT=$(git rev-parse "origin/$CURRENT_BRANCH")

if [[ "$LOCAL_COMMIT" != "$REMOTE_COMMIT" ]]; then
    echo "⚠️ WARNING: Local and remote commits differ" >&2
    echo "   Local:  $LOCAL_COMMIT" >&2
    echo "   Remote: $REMOTE_COMMIT" >&2

    # Check if we're ahead or behind
    AHEAD=$(git rev-list "origin/$CURRENT_BRANCH"..HEAD --count 2>/dev/null || echo "0")
    BEHIND=$(git rev-list HEAD.."origin/$CURRENT_BRANCH" --count 2>/dev/null || echo "0")

    if [[ "$AHEAD" -gt 0 ]]; then
        echo "   ⚠️ Local is $AHEAD commits AHEAD - unpushed work!" >&2
        echo "   This should have been caught above - investigating..." >&2
        cd "$ORIGINAL_DIR"
        exit 614
    fi

    if [[ "$BEHIND" -gt 0 ]]; then
        echo "   ℹ️ Local is $BEHIND commits BEHIND remote" >&2
        echo "   This is OK - source just needs to pull before we pull from it" >&2
        echo "   Proceeding (we'll get remote version)" >&2
    fi
fi

echo "   ✅ Local and remote are in sync"

# Get some stats
COMMIT_COUNT=$(git rev-list HEAD --count 2>/dev/null || echo "unknown")
LAST_COMMIT_DATE=$(git log -1 --format='%ci' 2>/dev/null || echo "unknown")
LAST_COMMIT_MSG=$(git log -1 --format='%s' 2>/dev/null || echo "unknown")

echo ""
echo "✅✅✅ Source effort verification PASSED ✅✅✅"
echo ""
echo "Source Effort Status:"
echo "  Effort ID: $SOURCE_EFFORT_ID"
echo "  Directory: $SOURCE_DIR"
echo "  Branch: $CURRENT_BRANCH"
echo "  Commits: $COMMIT_COUNT"
echo "  Last Commit: $LAST_COMMIT_DATE"
echo "  Last Message: $LAST_COMMIT_MSG"
echo ""
echo "Verification Results:"
echo "  ✅ All commits pushed to origin"
echo "  ✅ No uncommitted changes"
echo "  ✅ Local matches remote (or remote is ahead)"
echo "  ✅ Safe to pull from origin/$CURRENT_BRANCH"
echo ""

# Return to original directory
cd "$ORIGINAL_DIR"

exit 0
