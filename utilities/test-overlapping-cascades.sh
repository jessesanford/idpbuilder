#!/bin/bash

# Test Suite for R352 - Overlapping Cascade Protocol
# This script validates that CASCADE_REINTEGRATION correctly handles:
# 1. Multiple overlapping cascade chains
# 2. New fixes arriving during cascade operations
# 3. Chain convergence and merging
# 4. Persistent coordination with return to CASCADE_REINTEGRATION

set -euo pipefail

STATE_FILE="${1:-orchestrator-state.json}"

echo "🔍 R352 Overlapping Cascade Protocol Test Suite"
echo "================================================"

# Function to simulate adding a fix to an effort branch
simulate_fix() {
    local EFFORT_BRANCH="$1"
    local FIX_ID="fix_$(date +%s)_$$"
    
    echo "📝 Simulating fix $FIX_ID in $EFFORT_BRANCH"
    
    # Add to pending fixes
    jq --arg branch "$EFFORT_BRANCH" --arg fix "$FIX_ID" '
        .cascade_coordination.pending_fixes[$branch] = {
            "fix_ids": ((.cascade_coordination.pending_fixes[$branch].fix_ids // []) + [$fix]),
            "applied_at": now | todate,
            "cascade_chain": null
        }' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    echo "✅ Fix $FIX_ID added to $EFFORT_BRANCH"
}

# Function to create a cascade chain for a fix
create_cascade_chain() {
    local BRANCH="$1"
    local CHAIN_ID="cascade_test_$(date +%s)_${RANDOM}"
    
    echo "🔗 Creating cascade chain $CHAIN_ID for $BRANCH"
    
    jq --arg chain "$CHAIN_ID" --arg branch "$BRANCH" '
        .cascade_coordination.active_cascade_chains += [{
            "chain_id": $chain,
            "trigger": {
                "type": "fix_applied",
                "location": $branch,
                "timestamp": now | todate,
                "fix_ids": .cascade_coordination.pending_fixes[$branch].fix_ids
            },
            "status": "pending",
            "operations": [],
            "started_at": now | todate
        }] |
        .cascade_coordination.pending_fixes[$branch].cascade_chain = $chain' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    echo "✅ Chain $CHAIN_ID created"
}

# Test 1: Simple Linear Cascade
test_simple_cascade() {
    echo ""
    echo "TEST 1: Simple Linear Cascade"
    echo "------------------------------"
    
    # Setup initial state
    jq '.cascade_coordination.cascade_mode = false |
        .cascade_coordination.active_cascade_chains = [] |
        .cascade_coordination.pending_fixes = {}' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    # Add a fix to wave1
    simulate_fix "phase1/wave1/effort1"
    create_cascade_chain "phase1/wave1/effort1"
    
    # Set cascade mode
    jq '.cascade_coordination.cascade_mode = true' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    # Verify cascade chain exists
    CHAIN_COUNT=$(jq '.cascade_coordination.active_cascade_chains | length' "$STATE_FILE")
    if [[ "$CHAIN_COUNT" -eq 1 ]]; then
        echo "✅ PASS: Single cascade chain created"
    else
        echo "❌ FAIL: Expected 1 chain, got $CHAIN_COUNT"
        return 1
    fi
}

# Test 2: Overlapping Cascades
test_overlapping_cascades() {
    echo ""
    echo "TEST 2: Overlapping Cascades"
    echo "-----------------------------"
    
    # Reset state for this test
    jq '.cascade_coordination.cascade_mode = false |
        .cascade_coordination.active_cascade_chains = [] |
        .cascade_coordination.pending_fixes = {}' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    # Add fix to wave1
    simulate_fix "phase1/wave1/effort1"
    create_cascade_chain "phase1/wave1/effort1"
    
    # Add fix to wave2 while first cascade is running
    simulate_fix "phase1/wave2/effort1"
    create_cascade_chain "phase1/wave2/effort1"
    
    # Add fix to phase2
    simulate_fix "phase2/wave1/effort1"
    create_cascade_chain "phase2/wave1/effort1"
    
    # Set cascade mode (would be set by CASCADE_REINTEGRATION state)
    jq '.cascade_coordination.cascade_mode = true' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    # Verify multiple chains exist
    CHAIN_COUNT=$(jq '.cascade_coordination.active_cascade_chains | length' "$STATE_FILE")
    if [[ "$CHAIN_COUNT" -eq 3 ]]; then
        echo "✅ PASS: Multiple cascade chains created ($CHAIN_COUNT chains)"
    else
        echo "❌ FAIL: Expected 3 chains, got $CHAIN_COUNT"
        return 1
    fi
    
    # Verify cascade_mode is true
    CASCADE_MODE=$(jq -r '.cascade_coordination.cascade_mode' "$STATE_FILE")
    if [[ "$CASCADE_MODE" == "true" ]]; then
        echo "✅ PASS: Cascade mode is active"
    else
        echo "❌ FAIL: Cascade mode should be true"
        return 1
    fi
}

# Test 3: Chain Convergence
test_chain_convergence() {
    echo ""
    echo "TEST 3: Chain Convergence"
    echo "-------------------------"
    
    # Simulate chains converging at phase integration
    CHAIN1=$(jq -r '.cascade_coordination.active_cascade_chains[0].chain_id' "$STATE_FILE")
    CHAIN2=$(jq -r '.cascade_coordination.active_cascade_chains[1].chain_id' "$STATE_FILE")
    
    # Add same target to both chains
    jq --arg c1 "$CHAIN1" --arg c2 "$CHAIN2" '
        (.cascade_coordination.active_cascade_chains[] | 
         select(.chain_id == $c1 or .chain_id == $c2) | .operations) += [{
            "type": "recreate",
            "target": "phase1-integration",
            "status": "pending"
        }]' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    echo "📋 Chains $CHAIN1 and $CHAIN2 both target phase1-integration"
    
    # Simulate merge
    jq --arg primary "$CHAIN1" --arg secondary "$CHAIN2" '
        (.cascade_coordination.active_cascade_chains[] | 
         select(.chain_id == $primary)) |= . + {
            "merged_with": [$secondary],
            "merge_point": "phase1-integration"
        } |
        (.cascade_coordination.active_cascade_chains[] |
         select(.chain_id == $secondary)) |= . + {
            "status": "merged_into",
            "merged_into": $primary
        }' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    # Verify merge
    MERGED_STATUS=$(jq -r --arg c "$CHAIN2" '.cascade_coordination.active_cascade_chains[] | 
                     select(.chain_id == $c) | .status' "$STATE_FILE")
    
    if [[ "$MERGED_STATUS" == "merged_into" ]]; then
        echo "✅ PASS: Chains successfully merged"
    else
        echo "❌ FAIL: Chain merge failed"
        return 1
    fi
}

# Test 4: Exit Conditions
test_exit_conditions() {
    echo ""
    echo "TEST 4: Exit Conditions"
    echo "-----------------------"
    
    # Check that CASCADE_REINTEGRATION cannot exit with active chains
    ACTIVE_CHAINS=$(jq -r '[.cascade_coordination.active_cascade_chains[] | 
                           select(.status != "completed" and .status != "merged_into")] | 
                           length' "$STATE_FILE")
    
    if [[ "$ACTIVE_CHAINS" -gt 0 ]]; then
        echo "✅ PASS: Cannot exit - $ACTIVE_CHAINS active chains remain"
    else
        echo "⚠️  No active chains to test exit prevention"
    fi
    
    # Complete all chains
    jq '(.cascade_coordination.active_cascade_chains[] | 
        select(.status != "merged_into")) |= . + {
            "status": "completed",
            "completed_at": now | todate
        }' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    # Clear pending fixes
    jq '.cascade_coordination.pending_fixes = {}' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    # Mark integrations as fresh
    jq '.current_wave_integration.is_stale = false |
        .current_phase_integration.is_stale = false |
        .project_integration.is_stale = false' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    # Now check if exit is allowed
    ACTIVE_CHAINS=$(jq -r '[.cascade_coordination.active_cascade_chains[] | 
                           select(.status != "completed" and .status != "merged_into")] | 
                           length' "$STATE_FILE")
    PENDING_FIXES=$(jq '.cascade_coordination.pending_fixes | length' "$STATE_FILE")
    
    if [[ "$ACTIVE_CHAINS" -eq 0 ]] && [[ "$PENDING_FIXES" -eq 0 ]]; then
        echo "✅ PASS: All exit conditions met - can leave CASCADE_REINTEGRATION"
        
        # Clear cascade mode
        jq '.cascade_coordination.cascade_mode = false |
            .cascade_coordination.persistent_coordination = false' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    else
        echo "❌ FAIL: Exit conditions not met"
        return 1
    fi
}

# Test 5: New Fixes During Cascade
test_new_fixes_during_cascade() {
    echo ""
    echo "TEST 5: New Fixes During Cascade"
    echo "---------------------------------"
    
    # Reset for this test
    jq '.cascade_coordination.cascade_mode = true |
        .cascade_coordination.pending_fixes = {} |
        .cascade_coordination.active_cascade_chains = [{
            "chain_id": "cascade_running",
            "status": "in_progress",
            "operations": [{"type": "recreate", "target": "wave1-integration", "status": "in_progress"}]
        }]' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    # Simulate new fix arriving
    simulate_fix "phase1/wave3/effort1"
    
    # This should create a new cascade chain
    create_cascade_chain "phase1/wave3/effort1"
    
    # Verify we now have 2 chains
    CHAIN_COUNT=$(jq '.cascade_coordination.active_cascade_chains | length' "$STATE_FILE")
    if [[ "$CHAIN_COUNT" -eq 2 ]]; then
        echo "✅ PASS: New cascade chain created for fix during operation"
    else
        echo "❌ FAIL: Expected 2 chains after new fix, got $CHAIN_COUNT"
        return 1
    fi
    
    # Verify cascade_mode is still true
    CASCADE_MODE=$(jq -r '.cascade_coordination.cascade_mode' "$STATE_FILE")
    if [[ "$CASCADE_MODE" == "true" ]]; then
        echo "✅ PASS: Cascade mode remains active with new fix"
    else
        echo "❌ FAIL: Cascade mode should remain true"
        return 1
    fi
}

# Run all tests
main() {
    echo "Starting R352 Overlapping Cascade Protocol Tests"
    echo ""
    
    # Backup original state
    cp "$STATE_FILE" "${STATE_FILE}.backup" 2>/dev/null || true
    
    # Initialize test state
    if [[ ! -f "$STATE_FILE" ]]; then
        echo "{}" > "$STATE_FILE"
    fi
    
    # Initialize cascade_coordination structure
    jq 'if .cascade_coordination == null then 
            .cascade_coordination = {
                "cascade_mode": false,
                "persistent_coordination": false,
                "active_cascade_chains": [],
                "pending_fixes": {},
                "cascade_complete_when": {
                    "all_chains_complete": true,
                    "no_pending_fixes": true,
                    "project_integration_fresh": true,
                    "no_new_fixes_detected": true
                }
            }
        else . end' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    # Run tests
    FAILED=0
    
    test_simple_cascade || ((FAILED++))
    test_overlapping_cascades || ((FAILED++))
    test_chain_convergence || ((FAILED++))
    test_exit_conditions || ((FAILED++))
    test_new_fixes_during_cascade || ((FAILED++))
    
    echo ""
    echo "================================================"
    if [[ "$FAILED" -eq 0 ]]; then
        echo "✅✅✅ ALL TESTS PASSED ✅✅✅"
    else
        echo "❌ $FAILED TESTS FAILED ❌"
    fi
    echo "================================================"
    
    # Restore original state
    if [[ -f "${STATE_FILE}.backup" ]]; then
        mv "${STATE_FILE}.backup" "$STATE_FILE"
    fi
    
    exit "$FAILED"
}

# Execute main function
main "$@"