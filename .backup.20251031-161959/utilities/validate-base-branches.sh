#!/bin/bash

# R509 Enforcement: Validate all base branches
# This script validates that all effort branches follow the cascade pattern

set -e

echo "🔍 R509: Mandatory Base Branch Validation"
echo "=========================================="
echo "Validating cascade pattern per R501/R509/R510"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Exit codes
EXIT_CODE=0

# Load orchestrator state
STATE_FILE="${1:-/workspaces/software-factory-2.0/orchestrator-state-v3.json}"

if [ ! -f "$STATE_FILE" ]; then
    # Try default location
    STATE_FILE="/home/vscode/software-factory-template/orchestrator-state-v3.json"
    if [ ! -f "$STATE_FILE" ]; then
        echo -e "${RED}🚨 FATAL: No orchestrator-state-v3.json found!${NC}"
        echo "Searched:"
        echo "  - /workspaces/software-factory-2.0/orchestrator-state-v3.json"
        echo "  - /home/vscode/software-factory-template/orchestrator-state-v3.json"
        exit 509
    fi
fi

echo "Using state file: $STATE_FILE"
echo ""

# Check if pre_planned_infrastructure exists
if ! jq -e '.pre_planned_infrastructure' "$STATE_FILE" > /dev/null 2>&1; then
    echo -e "${RED}🚨 R504 VIOLATION: No pre_planned_infrastructure in state file!${NC}"
    echo "Run: bash utilities/pre-calculate-infrastructure.sh"
    exit 504
fi

# Validate CASCADE pattern in pre_planned_infrastructure
echo "📊 Validating CASCADE PATTERN in pre_planned_infrastructure..."
echo "--------------------------------------------------------------"

# Get all efforts sorted by phase, wave, and index
jq -r '.pre_planned_infrastructure.efforts | to_entries[] |
       "\(.key):\(.value.branch_name):\(.value.base_branch // "null"):\(.value.full_path // "null"):\(.value.phase // "null"):\(.value.wave // "null"):\(.value.index // "1"):\(.value.created // false)"' \
"$STATE_FILE" | sort -t: -k5,5 -k6,6 -k7,7n | while IFS=: read -r effort_id branch base path phase wave index created; do

    echo ""
    echo "Validating: $effort_id"
    echo "  Phase/Wave: $phase/$wave (index: $index)"
    echo "  Branch: $branch"
    echo "  Expected base: $base"
    echo "  Path: $path"
    echo "  Created: $created"

    # R510: Check base_branch exists
    if [ "$base" = "null" ] || [ -z "$base" ]; then
        echo -e "  ${RED}❌ R510 VIOLATION: No base_branch specified!${NC}"
        EXIT_CODE=510
        continue
    fi

    # R509: Validate cascade pattern
    if [[ "$phase" = "phase1" && "$wave" = "wave1" && "$index" = "1" ]]; then
        # First effort of P1W1 MUST be from main
        if [ "$base" != "main" ]; then
            echo -e "  ${RED}❌ R509 VIOLATION: First effort must be from main!${NC}"
            echo -e "  ${RED}   Got base: $base${NC}"
            EXIT_CODE=509
        else
            echo -e "  ${GREEN}✅ Correctly set to branch from main (first effort)${NC}"
        fi
    else
        # All other efforts MUST NOT be from main (cascade pattern)
        if [ "$base" = "main" ]; then
            echo -e "  ${RED}❌ R509 VIOLATION: Non-first effort cannot branch from main!${NC}"
            echo -e "  ${RED}   Must follow cascade pattern per R501${NC}"
            EXIT_CODE=509
        else
            echo -e "  ${GREEN}✅ Correctly set to cascade from: $base${NC}"
        fi
    fi

    # If infrastructure is created, validate it
    if [ "$created" = "true" ] && [ -d "$path/.git" ] 2>/dev/null; then
        echo "  📁 Infrastructure exists, validating..."

        cd "$path" 2>/dev/null || {
            echo -e "  ${YELLOW}⚠️ Cannot access path: $path${NC}"
            continue
        }

        # Check current branch
        CURRENT=$(git branch --show-current 2>/dev/null || echo "UNKNOWN")
        if [ "$CURRENT" != "$branch" ]; then
            echo -e "  ${RED}❌ WRONG BRANCH! Current: $CURRENT${NC}"
            EXIT_CODE=509
        else
            echo -e "  ${GREEN}✅ Correct branch checked out${NC}"
        fi

        # Check base branch relationship
        if [ "$base" = "main" ]; then
            # Should be based directly on main
            BASE_COMMIT=$(git merge-base HEAD origin/main 2>/dev/null || echo "UNKNOWN")
            MAIN_COMMIT=$(git rev-parse origin/main 2>/dev/null || echo "UNKNOWN")

            if [ "$BASE_COMMIT" != "$MAIN_COMMIT" ]; then
                echo -e "  ${RED}❌ NOT BASED ON MAIN!${NC}"
                EXIT_CODE=509
            else
                echo -e "  ${GREEN}✅ Correctly based on main${NC}"
            fi
        else
            # Should be based on previous effort
            if git show-ref --verify --quiet "refs/remotes/origin/$base" 2>/dev/null; then
                if git merge-base --is-ancestor "origin/$base" HEAD 2>/dev/null; then
                    echo -e "  ${GREEN}✅ Correctly based on $base${NC}"
                else
                    echo -e "  ${RED}❌ NOT BASED ON $base!${NC}"
                    EXIT_CODE=509
                fi
            else
                echo -e "  ${YELLOW}⚠️ Base branch $base not found in origin${NC}"
            fi
        fi
    elif [ "$created" = "true" ]; then
        echo -e "  ${YELLOW}⚠️ Marked as created but infrastructure not found${NC}"
    fi
done

echo ""
echo "============================================"

# Check final_merge_plan if it exists
if jq -e '.final_merge_plan' "$STATE_FILE" > /dev/null 2>&1; then
    echo ""
    echo "📊 Validating final_merge_plan CASCADE..."
    echo "----------------------------------------"

    # Validate merge sequence follows cascade
    PREV_BRANCH=""
    jq -r '.final_merge_plan.merge_sequence[] |
           "\(.order):\(.branch):\(.base_branch)"' \
    "$STATE_FILE" | while IFS=: read -r order branch base; do

        echo "Order $order: $branch (base: $base)"

        if [[ $order -eq 1 ]]; then
            # First must be from main
            if [ "$base" != "main" ]; then
                echo -e "  ${RED}❌ First merge not from main!${NC}"
                EXIT_CODE=509
            else
                echo -e "  ${GREEN}✅ Correctly from main${NC}"
            fi
        else
            # Others must cascade from previous
            if [ -n "$PREV_BRANCH" ] && [ "$base" != "$PREV_BRANCH" ]; then
                echo -e "  ${RED}❌ Not cascaded from previous ($PREV_BRANCH)!${NC}"
                EXIT_CODE=509
            else
                echo -e "  ${GREEN}✅ Correctly cascaded from $base${NC}"
            fi
        fi

        PREV_BRANCH="$branch"
    done
fi

# Summary
echo ""
echo "============================================"
if [ "${EXIT_CODE}" -ne 0 ]; then
    echo -e "${RED}🚨🚨🚨 R509/R510 VIOLATIONS DETECTED!${NC}"
    echo -e "${RED}CASCADE PATTERN BROKEN!${NC}"
    echo ""
    echo "REQUIRED ACTIONS:"
    echo "1. Fix base_branch values in pre_planned_infrastructure"
    echo "2. Ensure first effort of P1W1 branches from main"
    echo "3. Ensure all other efforts cascade from previous"
    echo "4. Re-run infrastructure creation if needed"
    echo ""
    echo "See rules:"
    echo "  - R501: Progressive trunk-based development"
    echo "  - R509: Mandatory base branch validation"
    echo "  - R510: Infrastructure creation protocol"
    exit ${EXIT_CODE}
else
    echo -e "${GREEN}✅ All base branches validated successfully${NC}"
    echo -e "${GREEN}✅ CASCADE PATTERN INTACT${NC}"
    echo ""
    echo "Progressive trunk-based development structure confirmed!"
fi