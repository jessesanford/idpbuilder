#!/bin/bash

# validate-infrastructure-target-url.sh - Validates that all efforts have target_repo_url
# Usage: tools/validate-infrastructure-target-url.sh [path/to/orchestrator-state-v3.json]
# Returns: 0 for valid, 1 for invalid

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get the script directory (tools/)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Default state file location
DEFAULT_STATE_FILE="$PROJECT_ROOT/orchestrator-state-v3.json"

# Get state file from argument or use default
STATE_FILE="${1:-$DEFAULT_STATE_FILE}"

# Check if state file exists
if [ ! -f "$STATE_FILE" ]; then
    echo -e "${RED}❌ State file not found: $STATE_FILE${NC}"
    exit 1
fi

echo -e "${BLUE}🔍 Validating target_repo_url in pre_planned_infrastructure (R504)${NC}"
echo "────────────────────────────────────────────────────────────"

# Track validation status
VALIDATION_FAILED=false

# Check if pre_planned_infrastructure exists
HAS_PRE_PLANNED=$(jq -r 'has("pre_planned_infrastructure")' "$STATE_FILE")

if [ "$HAS_PRE_PLANNED" != "true" ]; then
    echo -e "${YELLOW}⚠️  No pre_planned_infrastructure found (might be using dynamic strategy)${NC}"
    exit 0
fi

# Check parent-level target_repo_url
PARENT_URL=$(jq -r '.pre_planned_infrastructure.target_repo_url // "NOT_SET"' "$STATE_FILE")

if [ "$PARENT_URL" = "NOT_SET" ] || [ "$PARENT_URL" = "null" ]; then
    echo -e "${RED}❌ BLOCKING: Parent-level target_repo_url missing in pre_planned_infrastructure${NC}"
    VALIDATION_FAILED=true
else
    echo -e "${GREEN}✅ Parent-level target_repo_url: $PARENT_URL${NC}"
fi

# Get all effort keys
EFFORT_KEYS=$(jq -r '.pre_planned_infrastructure.efforts // {} | keys[]' "$STATE_FILE" 2>/dev/null || echo "")

if [ -z "$EFFORT_KEYS" ]; then
    echo -e "${YELLOW}⚠️  No efforts found in pre_planned_infrastructure${NC}"
else
    echo ""
    echo "Checking efforts:"

    # Check each effort
    for effort_key in $EFFORT_KEYS; do
        EFFORT_URL=$(jq -r ".pre_planned_infrastructure.efforts[\"$effort_key\"].target_repo_url // \"NOT_SET\"" "$STATE_FILE")

        if [ "$EFFORT_URL" = "NOT_SET" ] || [ "$EFFORT_URL" = "null" ]; then
            echo -e "${RED}  ❌ $effort_key: Missing target_repo_url${NC}"
            VALIDATION_FAILED=true
        elif [ "$PARENT_URL" != "NOT_SET" ] && [ "$EFFORT_URL" != "$PARENT_URL" ]; then
            echo -e "${YELLOW}  ⚠️  $effort_key: URL mismatch (effort: $EFFORT_URL != parent: $PARENT_URL)${NC}"
            VALIDATION_FAILED=true
        else
            echo -e "${GREEN}  ✅ $effort_key: Has target_repo_url${NC}"
        fi
    done
fi

echo ""
echo "────────────────────────────────────────────────────────────"

if [ "$VALIDATION_FAILED" = true ]; then
    echo -e "${RED}❌ VALIDATION FAILED: target_repo_url requirements not met (R504)${NC}"
    echo ""
    echo -e "${YELLOW}Required fixes:${NC}"
    echo "1. Ensure pre_planned_infrastructure has target_repo_url at parent level"
    echo "2. Ensure each effort has matching target_repo_url"
    echo "3. All effort URLs must match the parent URL"
    echo ""
    echo "Example structure:"
    echo '  "pre_planned_infrastructure": {'
    echo '    "target_repo_url": "https://github.com/user/repo.git",'
    echo '    "efforts": {'
    echo '      "effort_key": {'
    echo '        "target_repo_url": "https://github.com/user/repo.git",'
    echo '        ...'
    echo '      }'
    echo '    }'
    echo '  }'
    exit 1
else
    echo -e "${GREEN}✅ All target_repo_url validations passed!${NC}"
    exit 0
fi