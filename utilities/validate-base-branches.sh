#!/usr/bin/env bash

# Base Branch Validation Utility
# Enforces R337: orchestrator-state.json as single source of truth
# Validates all base branch tracking requirements

set -euo pipefail

STATE_FILE="${1:-orchestrator-state.json}"

if [ ! -f "$STATE_FILE" ]; then
    echo "❌ FATAL: State file not found: $STATE_FILE"
    exit 1
fi

echo "🔍 BASE BRANCH VALIDATION REPORT"
echo "================================="
echo "State File: $STATE_FILE"
echo "Timestamp: $(date -Iseconds)"
echo ""

# Track validation results
ERRORS=0
WARNINGS=0

# Function to check base branch tracking
check_base_tracking() {
    local section=$1
    local name=$2
    local tracking=$3
    
    if [ "$tracking" = "null" ] || [ -z "$tracking" ]; then
        echo "  ❌ ERROR: No base_branch_tracking for $name"
        ((ERRORS++))
        return 1
    fi
    
    # Check required fields
    local planned=$(echo "$tracking" | jq -r '.planned_base')
    local actual=$(echo "$tracking" | jq -r '.actual_base')
    local branched_at=$(echo "$tracking" | jq -r '.branched_at')
    
    if [ "$planned" = "null" ] || [ -z "$planned" ]; then
        echo "  ❌ ERROR: Missing planned_base for $name"
        ((ERRORS++))
    fi
    
    if [ "$actual" = "null" ] || [ -z "$actual" ]; then
        echo "  ❌ ERROR: Missing actual_base for $name"
        ((ERRORS++))
    fi
    
    if [ "$branched_at" = "null" ] || [ -z "$branched_at" ]; then
        echo "  ⚠️ WARNING: Missing branched_at timestamp for $name"
        ((WARNINGS++))
    fi
    
    # Check consistency
    if [ "$planned" != "null" ] && [ "$actual" != "null" ] && [ "$planned" != "$actual" ]; then
        local reason=$(echo "$tracking" | jq -r '.rebase_reason // .reason')
        if [ "$reason" = "null" ] || [ -z "$reason" ]; then
            echo "  ⚠️ WARNING: planned_base ($planned) != actual_base ($actual) without reason for $name"
            ((WARNINGS++))
        fi
    fi
    
    # Check rebase requirements
    local requires_rebase=$(echo "$tracking" | jq -r '.requires_rebase')
    if [ "$requires_rebase" = "true" ]; then
        echo "  🔄 REBASE REQUIRED: $name needs rebase"
        local rebase_reason=$(echo "$tracking" | jq -r '.rebase_reason')
        if [ "$rebase_reason" != "null" ]; then
            echo "     Reason: $rebase_reason"
        fi
    fi
    
    return 0
}

echo "## 1. EFFORTS IN PROGRESS"
echo "-------------------------"
EFFORTS_IN_PROGRESS=$(jq -c '.efforts_in_progress[]?' "$STATE_FILE")
if [ -z "$EFFORTS_IN_PROGRESS" ]; then
    echo "  (No efforts in progress)"
else
    while IFS= read -r effort; do
        NAME=$(echo "$effort" | jq -r '.name')
        echo "  📦 $NAME:"
        
        # Check base_branch field (legacy)
        BASE_BRANCH=$(echo "$effort" | jq -r '.base_branch // empty')
        if [ -n "$BASE_BRANCH" ]; then
            echo "     Legacy base_branch: $BASE_BRANCH"
        fi
        
        # Check base_branch_tracking (R337)
        TRACKING=$(echo "$effort" | jq '.base_branch_tracking // empty')
        check_base_tracking "efforts_in_progress" "$NAME" "$TRACKING"
    done <<< "$EFFORTS_IN_PROGRESS"
fi
echo ""

echo "## 2. EFFORTS COMPLETED"
echo "-----------------------"
EFFORTS_COMPLETED=$(jq -c '.efforts_completed[]?' "$STATE_FILE")
if [ -z "$EFFORTS_COMPLETED" ]; then
    echo "  (No efforts completed)"
else
    while IFS= read -r effort; do
        NAME=$(echo "$effort" | jq -r '.name')
        echo "  ✅ $NAME:"
        
        # Check base_branch field (legacy)
        BASE_BRANCH=$(echo "$effort" | jq -r '.base_branch // empty')
        if [ -n "$BASE_BRANCH" ]; then
            echo "     Legacy base_branch: $BASE_BRANCH"
        fi
        
        # Check base_branch_tracking (R337)
        TRACKING=$(echo "$effort" | jq '.base_branch_tracking // empty')
        check_base_tracking "efforts_completed" "$NAME" "$TRACKING"
    done <<< "$EFFORTS_COMPLETED"
fi
echo ""

echo "## 3. SPLIT TRACKING"
echo "--------------------"
SPLIT_EFFORTS=$(jq -r '.split_tracking | keys[]?' "$STATE_FILE")
if [ -z "$SPLIT_EFFORTS" ]; then
    echo "  (No split tracking)"
else
    for EFFORT in $SPLIT_EFFORTS; do
        echo "  🔀 $EFFORT:"
        SPLITS=$(jq -c ".split_tracking[\"$EFFORT\"].splits[]?" "$STATE_FILE")
        while IFS= read -r split; do
            NUMBER=$(echo "$split" | jq -r '.number')
            BRANCH=$(echo "$split" | jq -r '.branch')
            echo "     Split $NUMBER ($BRANCH):"
            
            # Check base_branch field (legacy)
            BASE_BRANCH=$(echo "$split" | jq -r '.base_branch // empty')
            if [ -n "$BASE_BRANCH" ]; then
                echo "        Legacy base_branch: $BASE_BRANCH"
            fi
            
            # Check base_branch_tracking (R337)
            TRACKING=$(echo "$split" | jq '.base_branch_tracking // empty')
            if [ "$TRACKING" != "null" ] && [ -n "$TRACKING" ]; then
                check_base_tracking "split" "$BRANCH" "$TRACKING"
            else
                echo "        ⚠️ WARNING: No base_branch_tracking"
                ((WARNINGS++))
            fi
        done <<< "$SPLITS"
    done
fi
echo ""

echo "## 4. BASE BRANCH DECISIONS"
echo "---------------------------"
DECISIONS=$(jq '.base_branch_decisions // empty' "$STATE_FILE")
if [ "$DECISIONS" = "null" ] || [ -z "$DECISIONS" ]; then
    echo "  ❌ ERROR: Missing base_branch_decisions section (R337 requirement)"
    ((ERRORS++))
else
    CURRENT_WAVE_BASE=$(echo "$DECISIONS" | jq -r '.current_wave_base // empty')
    CURRENT_PHASE_BASE=$(echo "$DECISIONS" | jq -r '.current_phase_base // empty')
    
    if [ -n "$CURRENT_WAVE_BASE" ]; then
        echo "  Current Wave Base: $CURRENT_WAVE_BASE"
    fi
    if [ -n "$CURRENT_PHASE_BASE" ]; then
        echo "  Current Phase Base: $CURRENT_PHASE_BASE"
    fi
    
    # Check decision log
    DECISION_COUNT=$(echo "$DECISIONS" | jq '.decision_log | length')
    echo "  Decision Log: $DECISION_COUNT entries"
    
    if [ "$DECISION_COUNT" -gt 0 ]; then
        echo "  Recent Decisions:"
        echo "$DECISIONS" | jq -r '.decision_log[-3:][] | "    - \(.timestamp): \(.decision)"' 2>/dev/null || true
    fi
fi
echo ""

echo "## 5. CASCADE TRACKING"
echo "----------------------"
# Check for efforts requiring rebase
NEEDS_REBASE=$(jq -r '[
    (.efforts_in_progress[]? | select(.base_branch_tracking.requires_rebase == true) | .name),
    (.efforts_completed[]? | select(.base_branch_tracking.requires_rebase == true) | .name)
] | unique | .[]' "$STATE_FILE" 2>/dev/null)

if [ -n "$NEEDS_REBASE" ]; then
    echo "  🚨 Efforts requiring rebase:"
    while IFS= read -r effort; do
        echo "    - $effort"
    done <<< "$NEEDS_REBASE"
else
    echo "  ✅ No efforts require rebase"
fi
echo ""

echo "## 6. VALIDATION SUMMARY"
echo "------------------------"
echo "  Errors: $ERRORS"
echo "  Warnings: $WARNINGS"

if [ $ERRORS -gt 0 ]; then
    echo ""
    echo "❌ VALIDATION FAILED: $ERRORS errors found"
    echo "Fix these R337 violations before proceeding!"
    exit 1
elif [ $WARNINGS -gt 0 ]; then
    echo ""
    echo "⚠️ VALIDATION PASSED WITH WARNINGS"
    echo "Review warnings to ensure proper tracking"
    exit 0
else
    echo ""
    echo "✅ VALIDATION PASSED"
    echo "All base branches properly tracked per R337"
    exit 0
fi