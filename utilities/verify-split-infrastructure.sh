#!/bin/bash

# verify-split-infrastructure.sh
# Utility to verify split directory structure and infrastructure
# 
# IMPORTANT: Split infrastructure is created by the ORCHESTRATOR (R204)
# - Code Reviewer creates split PLANS (SPLIT-INVENTORY.md, SPLIT-PLAN-XXX.md)
# - Orchestrator creates split INFRASTRUCTURE (directories, clones, branches)
# - SW Engineer IMPLEMENTS in the pre-created infrastructure

set -e

echo "═══════════════════════════════════════════════════════════════"
echo "🔍 SPLIT INFRASTRUCTURE VERIFICATION TOOL"
echo "═══════════════════════════════════════════════════════════════"

# Function to check split directory structure
verify_split_structure() {
    local EFFORT_NAME="$1"
    local PHASE="${2:-1}"
    local WAVE="${3:-1}"
    local EXPECTED_SPLITS="${4:-3}"
    
    echo ""
    echo "Checking split infrastructure for: $EFFORT_NAME"
    echo "Phase: $PHASE, Wave: $WAVE"
    echo "Expected splits: $EXPECTED_SPLITS"
    echo ""
    
    local BASE_DIR="efforts/phase${PHASE}/wave${WAVE}"
    local ERRORS=0
    local WARNINGS=0
    
    # Check if original effort directory exists
    if [[ -d "$BASE_DIR/${EFFORT_NAME}" ]]; then
        echo "✅ Original effort directory found: $BASE_DIR/${EFFORT_NAME}"
        
        # Check if marked as deprecated
        if [[ -f "$BASE_DIR/${EFFORT_NAME}/DEPRECATED" ]]; then
            echo "✅ Original marked as DEPRECATED"
        else
            echo "⚠️  WARNING: Original not marked as deprecated"
            WARNINGS=$((WARNINGS + 1))
        fi
        
        # Check for split plans in original
        if [[ -f "$BASE_DIR/${EFFORT_NAME}/SPLIT-INVENTORY.md" ]]; then
            echo "✅ SPLIT-INVENTORY.md found in original"
        else
            echo "⚠️  WARNING: No SPLIT-INVENTORY.md in original"
            WARNINGS=$((WARNINGS + 1))
        fi
    else
        echo "⚠️  WARNING: Original effort directory not found: $BASE_DIR/${EFFORT_NAME}"
        WARNINGS=$((WARNINGS + 1))
    fi
    
    echo ""
    echo "Checking split directories..."
    
    # Check each expected split
    for i in $(seq 1 $EXPECTED_SPLITS); do
        SPLIT_NUM=$(printf "%03d" $i)
        SPLIT_DIR="$BASE_DIR/${EFFORT_NAME}-SPLIT-${SPLIT_NUM}"
        
        echo ""
        echo "Split $SPLIT_NUM:"
        
        # Check directory exists
        if [[ ! -d "$SPLIT_DIR" ]]; then
            echo "  ❌ ERROR: Split directory missing: $SPLIT_DIR"
            ERRORS=$((ERRORS + 1))
            continue
        fi
        echo "  ✅ Directory exists: $SPLIT_DIR"
        
        # Check it's a git repository
        if [[ ! -d "$SPLIT_DIR/.git" ]]; then
            echo "  ❌ ERROR: Not a git repository (missing .git)"
            ERRORS=$((ERRORS + 1))
        else
            echo "  ✅ Is a git repository"
            
            # Check branch name
            cd "$SPLIT_DIR"
            CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "")
            if [[ -z "$CURRENT_BRANCH" ]]; then
                echo "  ❌ ERROR: Cannot determine current branch"
                ERRORS=$((ERRORS + 1))
            elif [[ "$CURRENT_BRANCH" == *"SPLIT-${SPLIT_NUM}"* ]] || [[ "$CURRENT_BRANCH" == *"split-${SPLIT_NUM}"* ]]; then
                echo "  ✅ On correct branch: $CURRENT_BRANCH"
            else
                echo "  ❌ ERROR: Wrong branch. Current: $CURRENT_BRANCH"
                echo "     Expected: branch containing 'SPLIT-${SPLIT_NUM}' or 'split-${SPLIT_NUM}'"
                ERRORS=$((ERRORS + 1))
            fi
            
            # Check remote tracking
            TRACKING=$(git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null || echo "")
            if [[ -z "$TRACKING" ]]; then
                echo "  ⚠️  WARNING: No remote tracking configured"
                WARNINGS=$((WARNINGS + 1))
            else
                echo "  ✅ Remote tracking: $TRACKING"
            fi
            
            cd - > /dev/null 2>&1
        fi
        
        # Check for split plan file
        if [[ ! -f "$SPLIT_DIR/SPLIT-PLAN-${SPLIT_NUM}.md" ]]; then
            echo "  ❌ ERROR: Split plan missing: SPLIT-PLAN-${SPLIT_NUM}.md"
            ERRORS=$((ERRORS + 1))
        else
            echo "  ✅ Split plan found: SPLIT-PLAN-${SPLIT_NUM}.md"
        fi
    done
    
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "VERIFICATION SUMMARY:"
    echo "  Errors: $ERRORS"
    echo "  Warnings: $WARNINGS"
    
    if [[ $ERRORS -eq 0 ]]; then
        if [[ $WARNINGS -eq 0 ]]; then
            echo "  Result: ✅ ALL CHECKS PASSED"
        else
            echo "  Result: ⚠️  PASSED WITH WARNINGS"
        fi
    else
        echo "  Result: ❌ FAILED - Fix errors before proceeding"
    fi
    echo "═══════════════════════════════════════════════════════════════"
    
    return $ERRORS
}

# Function to list all splits in the system
list_all_splits() {
    echo ""
    echo "📋 ALL SPLIT DIRECTORIES IN SYSTEM:"
    echo "───────────────────────────────────"
    
    local SPLIT_COUNT=0
    
    # Find all directories with -SPLIT- in the name
    while IFS= read -r split_dir; do
        if [[ -n "$split_dir" ]]; then
            SPLIT_COUNT=$((SPLIT_COUNT + 1))
            echo "$SPLIT_COUNT. $split_dir"
            
            # Check if it's a git repo and get branch
            if [[ -d "$split_dir/.git" ]]; then
                cd "$split_dir"
                BRANCH=$(git branch --show-current 2>/dev/null || echo "no-branch")
                cd - > /dev/null 2>&1
                echo "   Branch: $BRANCH"
            else
                echo "   ⚠️  Not a git repository"
            fi
        fi
    done < <(find efforts -type d -name "*-SPLIT-*" 2>/dev/null | sort)
    
    if [[ $SPLIT_COUNT -eq 0 ]]; then
        echo "No split directories found"
    else
        echo ""
        echo "Total splits found: $SPLIT_COUNT"
    fi
}

# Function to check sequential branching
verify_sequential_branching() {
    local EFFORT_NAME="$1"
    local PHASE="${2:-1}"
    local WAVE="${3:-1}"
    local NUM_SPLITS="${4:-3}"
    
    echo ""
    echo "🔗 VERIFYING SEQUENTIAL BRANCHING for $EFFORT_NAME"
    echo "───────────────────────────────────────────────────"
    
    local BASE_DIR="efforts/phase${PHASE}/wave${WAVE}"
    
    for i in $(seq 1 $NUM_SPLITS); do
        SPLIT_NUM=$(printf "%03d" $i)
        SPLIT_DIR="$BASE_DIR/${EFFORT_NAME}-SPLIT-${SPLIT_NUM}"
        
        if [[ ! -d "$SPLIT_DIR/.git" ]]; then
            echo "Split $SPLIT_NUM: ❌ Not a git repository"
            continue
        fi
        
        cd "$SPLIT_DIR"
        
        # Get current branch
        CURRENT_BRANCH=$(git branch --show-current)
        echo ""
        echo "Split $SPLIT_NUM:"
        echo "  Branch: $CURRENT_BRANCH"
        
        # Try to determine what it's based on
        if [[ $i -eq 1 ]]; then
            echo "  Expected base: Same as original (e.g., phase-integration)"
        else
            PREV_NUM=$(printf "%03d" $((i - 1)))
            echo "  Expected base: Previous split (split-${PREV_NUM})"
        fi
        
        # Check commit history for clues
        FIRST_COMMIT_MSG=$(git log --oneline -1 --reverse | cut -d' ' -f2-)
        if [[ -n "$FIRST_COMMIT_MSG" ]]; then
            echo "  First commit: $FIRST_COMMIT_MSG"
        fi
        
        cd - > /dev/null 2>&1
    done
}

# Main script logic
case "${1:-}" in
    "")
        echo "Usage: $0 <effort-name> [phase] [wave] [num-splits]"
        echo "   or: $0 --list"
        echo "   or: $0 --branching <effort-name> [phase] [wave] [num-splits]"
        echo ""
        echo "Examples:"
        echo "  $0 api-types                    # Check api-types in phase1/wave1 with 3 splits"
        echo "  $0 api-types 2 1 4              # Check api-types in phase2/wave1 with 4 splits"
        echo "  $0 --list                        # List all split directories"
        echo "  $0 --branching api-types         # Verify sequential branching"
        ;;
    
    "--list"|"-l")
        list_all_splits
        ;;
    
    "--branching"|"-b")
        shift
        if [[ -z "${1:-}" ]]; then
            echo "ERROR: Effort name required"
            exit 1
        fi
        verify_sequential_branching "$@"
        ;;
    
    *)
        verify_split_structure "$@"
        ;;
esac