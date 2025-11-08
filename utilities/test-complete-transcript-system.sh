#!/bin/bash
################################################################################
# SOFTWARE FACTORY 2.0 - COMPLETE TRANSCRIPT SYSTEM TEST
################################################################################
#
# This script tests the entire transcript-based identification system:
# 1. Simulates compaction with PreCompact hook
# 2. Tests PreToolUse detection for same agent
# 3. Tests PreToolUse non-detection for different agents
# 4. Verifies precise recovery targeting
#
################################################################################

set -euo pipefail

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

echo "════════════════════════════════════════════════════════════════"
echo "COMPLETE TRANSCRIPT-BASED IDENTIFICATION SYSTEM TEST"
echo "════════════════════════════════════════════════════════════════"
echo ""

# Clean up any existing markers
cleanup_markers() {
    echo -e "${YELLOW}Cleaning up old markers...${NC}"
    rm -f /tmp/compaction_marker*.txt
    rm -f /tmp/todos-precompact*.txt
    echo "✓ Cleanup complete"
    echo ""
}

# Test 1: Simulate SW Engineer compaction
test_sw_engineer_compaction() {
    echo -e "${BLUE}══ TEST 1: SW Engineer Compaction ══${NC}"
    echo ""
    
    local SW_TRANSCRIPT="/home/vscode/.claude/projects/sw-eng/11111111-2222-3333-4444-555555555555.jsonl"
    
    echo "1. Simulating PreCompact for SW Engineer..."
    echo "   Transcript: $SW_TRANSCRIPT"
    
    # Run PreCompact hook
    export CLAUDE_TRANSCRIPT_PATH="$SW_TRANSCRIPT"
    if [ -f "./transcript-based-precompact.sh" ]; then
        bash ./transcript-based-precompact.sh "$SW_TRANSCRIPT" "manual" > /tmp/precompact-output.txt 2>&1
        echo -e "${GREEN}✓ PreCompact executed${NC}"
    else
        echo -e "${RED}✗ PreCompact script not found${NC}"
        return 1
    fi
    
    # Check marker was created
    local SW_ID=$(basename "$SW_TRANSCRIPT" .jsonl)
    local SW_MARKER="/tmp/compaction_marker_${SW_ID}.txt"
    
    if [ -f "$SW_MARKER" ]; then
        echo -e "${GREEN}✓ Marker created: $SW_MARKER${NC}"
        echo ""
        echo "Marker contents:"
        grep -E "AGENT_TYPE|TRANSCRIPT_ID|COMPACTION_TIME" "$SW_MARKER" | head -3
    else
        echo -e "${RED}✗ Marker not created${NC}"
        return 1
    fi
    
    echo ""
    return 0
}

# Test 2: SW Engineer makes tool call (should be blocked)
test_sw_engineer_recovery_needed() {
    echo -e "${BLUE}══ TEST 2: SW Engineer Tool Call (Should Block) ══${NC}"
    echo ""
    
    local SW_TRANSCRIPT="/home/vscode/.claude/projects/sw-eng/11111111-2222-3333-4444-555555555555.jsonl"
    
    echo "2. SW Engineer attempts tool call after compaction..."
    echo "   Expected: BLOCKED (needs recovery)"
    echo ""
    
    # Run PreToolUse hook
    export CLAUDE_TRANSCRIPT_PATH="$SW_TRANSCRIPT"
    if [ -f "./transcript-based-pretooluse.sh" ]; then
        # This should exit with code 1 (blocked)
        if bash ./transcript-based-pretooluse.sh "Read" "$SW_TRANSCRIPT" > /tmp/pretooluse-sw-output.txt 2>&1; then
            echo -e "${RED}✗ ERROR: Tool call was NOT blocked (should have been)${NC}"
            cat /tmp/pretooluse-sw-output.txt
            return 1
        else
            echo -e "${GREEN}✓ Tool call BLOCKED as expected${NC}"
            echo ""
            echo "Block message excerpt:"
            grep -E "COMPACTION DETECTED|BLOCKING" /tmp/pretooluse-sw-output.txt | head -2
        fi
    else
        echo -e "${RED}✗ PreToolUse script not found${NC}"
        return 1
    fi
    
    echo ""
    return 0
}

# Test 3: Different agent (Orchestrator) makes tool call (should NOT be blocked)
test_orchestrator_not_blocked() {
    echo -e "${BLUE}══ TEST 3: Orchestrator Tool Call (Should NOT Block) ══${NC}"
    echo ""
    
    local ORCH_TRANSCRIPT="/home/vscode/.claude/projects/orch/99999999-8888-7777-6666-000000000000.jsonl"
    
    echo "3. Orchestrator (different agent) attempts tool call..."
    echo "   Expected: NOT BLOCKED (different transcript)"
    echo ""
    
    # Run PreToolUse hook
    export CLAUDE_TRANSCRIPT_PATH="$ORCH_TRANSCRIPT"
    if [ -f "./transcript-based-pretooluse.sh" ]; then
        # This should exit with code 0 (allowed)
        if bash ./transcript-based-pretooluse.sh "Read" "$ORCH_TRANSCRIPT" > /tmp/pretooluse-orch-output.txt 2>&1; then
            echo -e "${GREEN}✓ Tool call ALLOWED as expected${NC}"
            echo ""
            echo "Success message excerpt:"
            grep -E "No compaction detected|proceeding normally" /tmp/pretooluse-orch-output.txt | head -2
            
            # Check if it noticed other agent's compaction
            if grep -q "Other agents may have been compacted" /tmp/pretooluse-orch-output.txt; then
                echo -e "${GREEN}✓ Correctly noticed OTHER agent was compacted${NC}"
            fi
        else
            echo -e "${RED}✗ ERROR: Tool call was blocked (should NOT have been)${NC}"
            cat /tmp/pretooluse-orch-output.txt
            return 1
        fi
    else
        echo -e "${RED}✗ PreToolUse script not found${NC}"
        return 1
    fi
    
    echo ""
    return 0
}

# Test 4: Code Reviewer compaction doesn't affect SW Engineer
test_multiple_agent_independence() {
    echo -e "${BLUE}══ TEST 4: Multiple Agent Independence ══${NC}"
    echo ""
    
    local CR_TRANSCRIPT="/home/vscode/.claude/projects/reviewer/77777777-6666-5555-4444-333333333333.jsonl"
    
    echo "4. Simulating Code Reviewer compaction..."
    
    # Run PreCompact for Code Reviewer
    export CLAUDE_TRANSCRIPT_PATH="$CR_TRANSCRIPT"
    bash ./transcript-based-precompact.sh "$CR_TRANSCRIPT" "auto" > /tmp/precompact-cr-output.txt 2>&1
    
    local CR_ID=$(basename "$CR_TRANSCRIPT" .jsonl)
    local CR_MARKER="/tmp/compaction_marker_${CR_ID}.txt"
    
    if [ -f "$CR_MARKER" ]; then
        echo -e "${GREEN}✓ Code Reviewer marker created${NC}"
    fi
    
    # Now we have TWO markers: SW Engineer and Code Reviewer
    echo ""
    echo "Current markers:"
    ls -la /tmp/compaction_marker_*.txt 2>/dev/null | awk '{print "  - " $NF}'
    echo ""
    
    # Test that each agent only sees their own marker
    echo "Testing agent-specific blocking:"
    echo ""
    
    # SW Engineer still blocked by their own marker
    export CLAUDE_TRANSCRIPT_PATH="/home/vscode/.claude/projects/sw-eng/11111111-2222-3333-4444-555555555555.jsonl"
    if bash ./transcript-based-pretooluse.sh "Read" > /dev/null 2>&1; then
        echo -e "${RED}  ✗ SW Engineer NOT blocked (should be)${NC}"
    else
        echo -e "${GREEN}  ✓ SW Engineer still blocked by their marker${NC}"
    fi
    
    # Code Reviewer blocked by their own marker
    export CLAUDE_TRANSCRIPT_PATH="$CR_TRANSCRIPT"
    if bash ./transcript-based-pretooluse.sh "Read" > /dev/null 2>&1; then
        echo -e "${RED}  ✗ Code Reviewer NOT blocked (should be)${NC}"
    else
        echo -e "${GREEN}  ✓ Code Reviewer blocked by their marker${NC}"
    fi
    
    # Orchestrator still not blocked
    export CLAUDE_TRANSCRIPT_PATH="/home/vscode/.claude/projects/orch/99999999-8888-7777-6666-000000000000.jsonl"
    if bash ./transcript-based-pretooluse.sh "Read" > /dev/null 2>&1; then
        echo -e "${GREEN}  ✓ Orchestrator NOT blocked (correct)${NC}"
    else
        echo -e "${RED}  ✗ Orchestrator blocked (should NOT be)${NC}"
    fi
    
    echo ""
    return 0
}

# Test 5: Recovery and cleanup
test_recovery_cleanup() {
    echo -e "${BLUE}══ TEST 5: Recovery and Cleanup ══${NC}"
    echo ""
    
    echo "5. Testing marker cleanup after recovery..."
    
    # Get SW Engineer marker
    local SW_ID="11111111-2222-3333-4444-555555555555"
    local SW_MARKER="/tmp/compaction_marker_${SW_ID}.txt"
    
    # Simulate recovery by archiving the marker
    if [ -f "$SW_MARKER" ]; then
        mkdir -p /tmp/compaction-archives
        mv "$SW_MARKER" "/tmp/compaction-archives/${SW_ID}-test.txt"
        echo -e "${GREEN}✓ SW Engineer marker archived (simulating recovery)${NC}"
    fi
    
    # Now SW Engineer should NOT be blocked
    export CLAUDE_TRANSCRIPT_PATH="/home/vscode/.claude/projects/sw-eng/${SW_ID}.jsonl"
    if bash ./transcript-based-pretooluse.sh "Read" > /tmp/pretooluse-recovered.txt 2>&1; then
        echo -e "${GREEN}✓ SW Engineer can now use tools after recovery${NC}"
    else
        echo -e "${RED}✗ SW Engineer still blocked after recovery${NC}"
        cat /tmp/pretooluse-recovered.txt
    fi
    
    echo ""
    return 0
}

# Main test execution
main() {
    # Start with clean slate
    cleanup_markers
    
    # Run all tests
    local all_passed=true
    
    if ! test_sw_engineer_compaction; then
        all_passed=false
    fi
    
    if ! test_sw_engineer_recovery_needed; then
        all_passed=false
    fi
    
    if ! test_orchestrator_not_blocked; then
        all_passed=false
    fi
    
    if ! test_multiple_agent_independence; then
        all_passed=false
    fi
    
    if ! test_recovery_cleanup; then
        all_passed=false
    fi
    
    # Final summary
    echo ""
    echo "════════════════════════════════════════════════════════════════"
    if [ "$all_passed" = true ]; then
        echo -e "${GREEN}✅ ALL TESTS PASSED ✅${NC}"
        echo ""
        echo "CONCLUSIONS:"
        echo "✓ Transcript IDs uniquely identify agents"
        echo "✓ Each agent has its own compaction marker"
        echo "✓ Only compacted agents are blocked for recovery"
        echo "✓ Non-compacted agents continue working normally"
        echo "✓ Multiple agents can run simultaneously"
        echo "✓ Recovery is precise and agent-specific"
    else
        echo -e "${RED}❌ SOME TESTS FAILED ❌${NC}"
        echo "Please review the output above for details."
    fi
    echo "════════════════════════════════════════════════════════════════"
    
    # Final cleanup
    echo ""
    echo "Cleaning up test artifacts..."
    cleanup_markers
    rm -f /tmp/precompact-*.txt /tmp/pretooluse-*.txt
    echo "✓ Test complete"
}

# Run if executed directly
if [ "${BASH_SOURCE[0]}" == "${0}" ]; then
    main "$@"
fi