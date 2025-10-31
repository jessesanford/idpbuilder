#!/bin/bash

# get-latest-plan.sh - Helper script to find the most recent plan file
# Supports both timestamped (new) and non-timestamped (legacy) formats
#
# Usage: ./get-latest-plan.sh [PLAN_TYPE] [DIRECTORY]
#   PLAN_TYPE: IMPLEMENTATION-PLAN, CODE-REVIEW-REPORT, SPLIT-PLAN, FIX-INSTRUCTIONS, INTEGRATE_WAVE_EFFORTS-REPORT
#   DIRECTORY: Optional directory to search in (defaults to current directory)
#
# Example: ./get-latest-plan.sh IMPLEMENTATION-PLAN /efforts/phase1/wave1/api-types

PLAN_TYPE="${1:-IMPLEMENTATION-PLAN}"
SEARCH_DIR="${2:-.}"

# Function to find the latest plan
find_latest_plan() {
    local plan_prefix="$1"
    local dir="$2"
    
    # Look for timestamped versions first (newest first)
    local timestamped=$(ls -t "$dir"/${plan_prefix}-*.md 2>/dev/null | grep -v COMPLETED | head -n1)
    
    # Fallback to legacy format if no timestamped versions
    local legacy="$dir/${plan_prefix}.md"
    
    if [ -n "$timestamped" ]; then
        echo "$timestamped"
        return 0
    elif [ -f "$legacy" ]; then
        echo "⚠️ Using legacy format: $legacy" >&2
        echo "$legacy"
        return 0
    else
        echo "❌ ERROR: No ${plan_prefix}*.md found in $dir" >&2
        return 1
    fi
}

# Main logic
case "$PLAN_TYPE" in
    IMPLEMENTATION-PLAN|implementation-plan)
        find_latest_plan "IMPLEMENTATION-PLAN" "$SEARCH_DIR"
        ;;
    CODE-REVIEW-REPORT|review-report|review)
        find_latest_plan "CODE-REVIEW-REPORT" "$SEARCH_DIR"
        ;;
    SPLIT-PLAN|split-plan|split)
        # For split plans, might need to handle numbered ones
        if find_latest_plan "SPLIT-PLAN" "$SEARCH_DIR"; then
            exit 0
        else
            # Try numbered split plans
            latest_split=$(ls -t "$SEARCH_DIR"/SPLIT-PLAN-[0-9]*.md 2>/dev/null | head -n1)
            if [ -n "$latest_split" ]; then
                echo "$latest_split"
            else
                exit 1
            fi
        fi
        ;;
    FIX-INSTRUCTIONS|fix-instructions|fix)
        find_latest_plan "FIX-INSTRUCTIONS" "$SEARCH_DIR"
        ;;
    INTEGRATE_WAVE_EFFORTS-REPORT|integration-report|integration)
        find_latest_plan "INTEGRATE_WAVE_EFFORTS-REPORT" "$SEARCH_DIR"
        ;;
    INTEGRATE_WAVE_EFFORTS-DEMO|integration-demo|demo)
        find_latest_plan "INTEGRATE_WAVE_EFFORTS-DEMO" "$SEARCH_DIR"
        ;;
    SPLIT-INVENTORY|split-inventory|inventory)
        find_latest_plan "SPLIT-INVENTORY" "$SEARCH_DIR"
        ;;
    *)
        echo "Usage: $0 [PLAN_TYPE] [DIRECTORY]" >&2
        echo "Valid PLAN_TYPEs: IMPLEMENTATION-PLAN, CODE-REVIEW-REPORT, SPLIT-PLAN, FIX-INSTRUCTIONS, INTEGRATE_WAVE_EFFORTS-REPORT" >&2
        exit 1
        ;;
esac