#!/bin/bash

# Get effort base branch from pre-planned infrastructure
# This script reads the pre-calculated cascade base per R308/R501
# Usage: ./get-effort-base-branch.sh <effort_id>
#
# CRITICAL: This enforces R308/R501 progressive trunk-based development
# by ensuring all efforts use their pre-calculated cascade bases

set -euo pipefail

CLAUDE_PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}"
STATE_FILE="$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to log messages
log() {
    echo -e "${GREEN}[CASCADE-BASE]${NC} $1" >&2
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
    exit 1
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" >&2
}

# Check if effort_id was provided
if [ $# -lt 1 ]; then
    echo "Usage: $0 <effort_id>"
    echo "Example: $0 phase1_wave1_authentication"
    exit 1
fi

EFFORT_ID="$1"

# Check if state file exists
if [[ ! -f "$STATE_FILE" ]]; then
    error "orchestrator-state-v3.json not found at $STATE_FILE"
fi

# Check if pre_planned_infrastructure exists
HAS_INFRASTRUCTURE=$(jq -r '.pre_planned_infrastructure // "null"' "$STATE_FILE")
if [[ "$HAS_INFRASTRUCTURE" == "null" ]]; then
    error "No pre_planned_infrastructure found! Run pre-calculate-infrastructure.sh first (R504)"
fi

# Check if cascade_bases_calculated is true
CASCADE_CALCULATED=$(jq -r '.pre_planned_infrastructure.validation.cascade_bases_calculated // false' "$STATE_FILE")
if [[ "$CASCADE_CALCULATED" != "true" ]]; then
    error "CASCADE bases not calculated! Run updated pre-calculate-infrastructure.sh (R308/R501)"
fi

# Get the base branch for the effort
BASE_BRANCH=$(jq -r --arg id "$EFFORT_ID" '
    .pre_planned_infrastructure.efforts[$id].base_branch //
    .pre_planned_infrastructure.cascade_bases[$id].base_branch //
    "null"
' "$STATE_FILE")

if [[ "$BASE_BRANCH" == "null" ]] || [[ -z "$BASE_BRANCH" ]]; then
    error "No cascade base found for effort $EFFORT_ID (R308 violation)"
fi

# Get the cascade reason
CASCADE_REASON=$(jq -r --arg id "$EFFORT_ID" '
    .pre_planned_infrastructure.efforts[$id].cascade_reason //
    .pre_planned_infrastructure.cascade_bases[$id].reason //
    "No reason recorded"
' "$STATE_FILE")

# Log the cascade information
log "Effort: $EFFORT_ID"
log "Base branch: $BASE_BRANCH"
log "Reason: $CASCADE_REASON"

# Output just the base branch for use in scripts
echo "$BASE_BRANCH"