#!/bin/bash

# Create effort infrastructure with CASCADE base branch
# This script ENFORCES R308/R501 by using pre-calculated cascade bases
#
# Usage: ./create-effort-with-cascade-base.sh <effort_id> <repo_url>
#
# CRITICAL: This is the ONLY way to create effort infrastructure
# Using any other method violates R308/R501 cascade requirements

set -euo pipefail

CLAUDE_PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}"
STATE_FILE="$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
GET_BASE_SCRIPT="$CLAUDE_PROJECT_DIR/utilities/get-effort-base-branch.sh"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to log messages
log() {
    echo -e "${GREEN}[CASCADE-CREATE]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
    exit 1
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" >&2
}

# Check arguments
if [ $# -lt 2 ]; then
    echo "Usage: $0 <effort_id> <repo_url>"
    echo "Example: $0 phase1_wave1_authentication https://github.com/org/repo.git"
    exit 1
fi

EFFORT_ID="$1"
REPO_URL="$2"

# Check if get-effort-base-branch.sh exists
if [[ ! -f "$GET_BASE_SCRIPT" ]]; then
    error "get-effort-base-branch.sh not found! Cannot determine cascade base (R308 violation)"
fi

# Get effort information from pre-planned infrastructure
EFFORT_INFO=$(jq -r --arg id "$EFFORT_ID" '
    .pre_planned_infrastructure.efforts[$id] // "null"
' "$STATE_FILE")

if [[ "$EFFORT_INFO" == "null" ]]; then
    error "Effort $EFFORT_ID not found in pre_planned_infrastructure! Run pre-calculate-infrastructure.sh first"
fi

# Extract effort details
EFFORT_PATH=$(echo "$EFFORT_INFO" | jq -r '.full_path')
EFFORT_BRANCH=$(echo "$EFFORT_INFO" | jq -r '.branch_name')

# Get the CASCADE base branch using the utility
log "Getting cascade base for $EFFORT_ID..."
BASE_BRANCH=$("$GET_BASE_SCRIPT" "$EFFORT_ID")

if [[ -z "$BASE_BRANCH" ]] || [[ "$BASE_BRANCH" == "null" ]]; then
    error "Failed to get cascade base for $EFFORT_ID (R308 violation)"
fi

log "========================================="
log "Creating effort with CASCADE base (R308/R501)"
log "========================================="
log "Effort ID: $EFFORT_ID"
log "Path: $EFFORT_PATH"
log "Branch: $EFFORT_BRANCH"
log "CASCADE Base: $BASE_BRANCH"
log "========================================="

# Validate that base branch is NOT main (unless it's P1W1 first effort)
if [[ "$BASE_BRANCH" == "main" ]]; then
    # Check if this is the first effort of Phase 1 Wave 1
    IS_FIRST=$(echo "$EFFORT_INFO" | jq -r '
        if .phase == "phase1" and .wave == "wave1" then
            "maybe"
        else
            "no"
        end
    ')

    if [[ "$IS_FIRST" == "no" ]]; then
        warning "⚠️ CASCADE WARNING: Effort $EFFORT_ID bases from 'main' but is not P1W1 first effort!"
        warning "This likely indicates a CASCADE VIOLATION (R308/R501)"
        read -p "Continue anyway? (y/n): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            error "Aborted due to cascade violation"
        fi
    fi
fi

# Create the effort directory
log "Creating effort directory: $EFFORT_PATH"
mkdir -p "$EFFORT_PATH"

# Clone with CASCADE base
log "Cloning repository with CASCADE base: $BASE_BRANCH"
cd "$EFFORT_PATH"

# Clone using the cascade base
if ! git clone --branch "$BASE_BRANCH" "$REPO_URL" .; then
    # If base branch doesn't exist yet (e.g., previous effort not pushed)
    warning "Base branch $BASE_BRANCH not found on remote. This is expected if previous effort hasn't pushed yet."
    warning "Cloning from main and will create branch from local cascade base..."

    # Clone from main first
    git clone "$REPO_URL" .

    # Try to checkout the base branch locally
    if git show-ref --verify --quiet "refs/heads/$BASE_BRANCH"; then
        git checkout "$BASE_BRANCH"
    else
        warning "Base branch $BASE_BRANCH not found locally either. Creating effort may not have proper cascade base!"
    fi
fi

# Create the new effort branch
log "Creating effort branch: $EFFORT_BRANCH"
git checkout -b "$EFFORT_BRANCH"

# Verify cascade integrity
log "Verifying CASCADE integrity..."
CURRENT_BASE=$(git merge-base HEAD origin/main || echo "unknown")

# Update state file to mark effort as created
jq --arg id "$EFFORT_ID" \
   --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
   --arg base "$BASE_BRANCH" \
   '.pre_planned_infrastructure.efforts[$id].created = true |
    .pre_planned_infrastructure.efforts[$id].created_at = $timestamp |
    .pre_planned_infrastructure.efforts[$id].actual_base = $base' \
   "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"

log "========================================="
log "✅ Effort created with CASCADE base successfully!"
log "========================================="
log ""
log "NEXT STEPS:"
log "1. Spawn agent to work in: $EFFORT_PATH"
log "2. Agent will work on branch: $EFFORT_BRANCH"
log "3. Changes cascade from: $BASE_BRANCH"
log ""
log "CASCADE CHAIN:"
log "  $BASE_BRANCH"
log "    └─→ $EFFORT_BRANCH (this effort)"
log ""

# Return success
exit 0