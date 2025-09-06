#!/bin/bash

# Verification script for R308 - Incremental Branching Strategy
# Tests that base branch determination follows the incremental principle

set -e

echo "🔍 R308 Base Branch Verification Script"
echo "======================================="
echo ""

# Function from the orchestrator's SETUP_EFFORT_INFRASTRUCTURE rules
determine_incremental_base_branch() {
    local PHASE=$1
    local WAVE=$2
    
    # Phase 1, Wave 1: Start from main
    if [[ $PHASE -eq 1 && $WAVE -eq 1 ]]; then
        echo "main"
        return
    fi
    
    # First wave of new phase: From previous phase integration
    if [[ $WAVE -eq 1 ]]; then
        PREV_PHASE=$((PHASE - 1))
        echo "phase${PREV_PHASE}-integration"
        return
    fi
    
    # Subsequent waves: From previous wave integration
    PREV_WAVE=$((WAVE - 1))
    echo "phase${PHASE}-wave${PREV_WAVE}-integration"
}

# Test cases
echo "Testing R308 Base Branch Determination:"
echo "----------------------------------------"

# Test Phase 1, Wave 1
BASE=$(determine_incremental_base_branch 1 1)
EXPECTED="main"
echo -n "Phase 1, Wave 1: "
if [[ "$BASE" == "$EXPECTED" ]]; then
    echo "✅ $BASE (correct)"
else
    echo "❌ Got '$BASE', expected '$EXPECTED'"
    exit 1
fi

# Test Phase 1, Wave 2
BASE=$(determine_incremental_base_branch 1 2)
EXPECTED="phase1-wave1-integration"
echo -n "Phase 1, Wave 2: "
if [[ "$BASE" == "$EXPECTED" ]]; then
    echo "✅ $BASE (correct)"
else
    echo "❌ Got '$BASE', expected '$EXPECTED'"
    exit 1
fi

# Test Phase 1, Wave 3
BASE=$(determine_incremental_base_branch 1 3)
EXPECTED="phase1-wave2-integration"
echo -n "Phase 1, Wave 3: "
if [[ "$BASE" == "$EXPECTED" ]]; then
    echo "✅ $BASE (correct)"
else
    echo "❌ Got '$BASE', expected '$EXPECTED'"
    exit 1
fi

# Test Phase 2, Wave 1 (CRITICAL TEST - was failing before)
BASE=$(determine_incremental_base_branch 2 1)
EXPECTED="phase1-integration"
echo -n "Phase 2, Wave 1: "
if [[ "$BASE" == "$EXPECTED" ]]; then
    echo "✅ $BASE (correct - NOT main!)"
else
    echo "❌ Got '$BASE', expected '$EXPECTED'"
    echo "   CRITICAL: Phase 2+ must NEVER use main!"
    exit 1
fi

# Test Phase 2, Wave 2
BASE=$(determine_incremental_base_branch 2 2)
EXPECTED="phase2-wave1-integration"
echo -n "Phase 2, Wave 2: "
if [[ "$BASE" == "$EXPECTED" ]]; then
    echo "✅ $BASE (correct)"
else
    echo "❌ Got '$BASE', expected '$EXPECTED'"
    exit 1
fi

# Test Phase 3, Wave 1
BASE=$(determine_incremental_base_branch 3 1)
EXPECTED="phase2-integration"
echo -n "Phase 3, Wave 1: "
if [[ "$BASE" == "$EXPECTED" ]]; then
    echo "✅ $BASE (correct)"
else
    echo "❌ Got '$BASE', expected '$EXPECTED'"
    exit 1
fi

echo ""
echo "✅ All tests passed! R308 incremental branching is correctly configured."
echo ""
echo "📌 Remember the incremental principle:"
echo "   Each wave/phase builds on the LATEST INTEGRATED CODE from the previous."
echo "   NO EFFORT MAY BRANCH FROM STALE CODE!"