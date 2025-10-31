#!/bin/bash

# Test script to verify transcript-based agent identification theory
# This tests if we can use transcript paths to precisely identify which agent was compacted

echo "========================================"
echo "TRANSCRIPT-BASED AGENT IDENTIFICATION TEST"
echo "========================================"
echo ""

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test 1: Extract transcript ID from a path
test_extract_transcript_id() {
    echo -e "${BLUE}TEST 1: Extracting Transcript ID from Path${NC}"
    echo "----------------------------------------"
    
    # Sample transcript paths (as seen in hooks)
    local path1="/home/vscode/.claude/projects/agent1/12345678-90ab-cdef-1234-567890abcdef.jsonl"
    local path2="/home/user/.claude/projects/-workspaces/41205a4d-093d-4be3-97dd-02e549eb3067.jsonl"
    local path3="/home/vscode/.claude/projects/orchestrator/98765432-10ab-cdef-4321-fedcba098765.jsonl"
    
    for path in "$path1" "$path2" "$path3"; do
        echo "Path: $path"
        
        # Extract just the UUID part (transcript ID)
        transcript_id=$(basename "$path" .jsonl)
        echo -e "  Extracted ID: ${GREEN}$transcript_id${NC}"
        
        # Alternative: Extract with directory context
        transcript_with_context=$(echo "$path" | sed 's|.*/projects/||; s|\.jsonl$||')
        echo -e "  With context: ${YELLOW}$transcript_with_context${NC}"
        echo ""
    done
}

# Test 2: Simulate PreCompact hook creating marker
test_precompact_marker_creation() {
    echo -e "${BLUE}TEST 2: Simulating PreCompact Marker Creation${NC}"
    echo "----------------------------------------"
    
    # Simulate what PreCompact does
    local TRANSCRIPT_PATH="$1"
    local AGENT_TYPE="$2"
    
    if [ -z "$TRANSCRIPT_PATH" ]; then
        TRANSCRIPT_PATH="/home/vscode/.claude/projects/-workspaces/41205a4d-093d-4be3-97dd-02e549eb3067.jsonl"
        AGENT_TYPE="sw-engineer"
    fi
    
    # Extract transcript ID
    local TRANSCRIPT_ID=$(basename "$TRANSCRIPT_PATH" .jsonl)
    
    # Create marker file with transcript ID in the name
    local MARKER_FILE="/tmp/compaction_marker_${TRANSCRIPT_ID}.txt"
    
    echo "Creating marker for:"
    echo "  Agent: $AGENT_TYPE"
    echo "  Transcript: $TRANSCRIPT_PATH"
    echo "  Transcript ID: $TRANSCRIPT_ID"
    echo -e "  Marker file: ${GREEN}$MARKER_FILE${NC}"
    
    # Create the marker
    cat > "$MARKER_FILE" << EOF
COMPACTION_TIME:$(date -u +"%Y-%m-%d %H:%M:%S UTC")
COMPACTION_TYPE:manual
AGENT_TYPE:$AGENT_TYPE
TRANSCRIPT_PATH:$TRANSCRIPT_PATH
TRANSCRIPT_ID:$TRANSCRIPT_ID
WORKING_DIR:$(pwd)
GIT_BRANCH:$(git branch --show-current 2>/dev/null || echo "unknown")
EOF
    
    echo ""
    echo "Marker contents:"
    cat "$MARKER_FILE"
    echo ""
}

# Test 3: Simulate PreToolUse detection
test_pretooluse_detection() {
    echo -e "${BLUE}TEST 3: Simulating PreToolUse Detection${NC}"
    echo "----------------------------------------"
    
    # This would be passed to PreToolUse hook
    local CURRENT_TRANSCRIPT="$1"
    
    if [ -z "$CURRENT_TRANSCRIPT" ]; then
        CURRENT_TRANSCRIPT="/home/vscode/.claude/projects/-workspaces/41205a4d-093d-4be3-97dd-02e549eb3067.jsonl"
    fi
    
    local CURRENT_TRANSCRIPT_ID=$(basename "$CURRENT_TRANSCRIPT" .jsonl)
    
    echo "Current session transcript: $CURRENT_TRANSCRIPT"
    echo "Current transcript ID: $CURRENT_TRANSCRIPT_ID"
    echo ""
    
    # Look for matching compaction marker
    echo "Searching for matching compaction marker..."
    
    local MARKER_PATTERN="/tmp/compaction_marker_${CURRENT_TRANSCRIPT_ID}.txt"
    
    if [ -f "$MARKER_PATTERN" ]; then
        echo -e "${GREEN}✓ MATCH FOUND!${NC}"
        echo "This agent WAS compacted (same transcript ID)"
        echo ""
        echo "Marker contents:"
        grep -E "AGENT_TYPE|TRANSCRIPT_ID|COMPACTION_TIME" "$MARKER_PATTERN"
        
        # Extract agent type for recovery
        local AGENT_TYPE=$(grep "AGENT_TYPE:" "$MARKER_PATTERN" | cut -d: -f2)
        echo ""
        echo -e "Recovery needed for: ${YELLOW}$AGENT_TYPE${NC}"
        
        return 0
    else
        echo -e "${RED}✗ No matching marker found${NC}"
        echo "This agent was NOT compacted (different transcript)"
        
        # Show what markers exist
        echo ""
        echo "Existing markers:"
        ls -la /tmp/compaction_marker_*.txt 2>/dev/null || echo "  None found"
        
        return 1
    fi
}

# Test 4: Full simulation
test_full_workflow() {
    echo -e "${BLUE}TEST 4: Full Workflow Simulation${NC}"
    echo "========================================"
    
    # Clean up old markers
    rm -f /tmp/compaction_marker_*.txt
    
    # Scenario 1: SW Engineer gets compacted
    echo -e "${YELLOW}Scenario 1: SW Engineer Compaction${NC}"
    echo "------------------------------------"
    local SW_TRANSCRIPT="/home/vscode/.claude/projects/sw-eng/11111111-2222-3333-4444-555555555555.jsonl"
    
    echo "1. PreCompact hook fires for SW Engineer:"
    test_precompact_marker_creation "$SW_TRANSCRIPT" "sw-engineer"
    
    echo "2. Later, SW Engineer makes a tool call:"
    test_pretooluse_detection "$SW_TRANSCRIPT"
    echo ""
    
    # Scenario 2: Different agent (orchestrator) makes call
    echo -e "${YELLOW}Scenario 2: Different Agent (Orchestrator) Not Compacted${NC}"
    echo "--------------------------------------------------------"
    local ORCH_TRANSCRIPT="/home/vscode/.claude/projects/orch/99999999-8888-7777-6666-000000000000.jsonl"
    
    echo "Orchestrator (different transcript) makes tool call:"
    test_pretooluse_detection "$ORCH_TRANSCRIPT"
    echo ""
    
    # Scenario 3: Same agent makes call after compaction
    echo -e "${YELLOW}Scenario 3: Same SW Engineer After Compaction${NC}"
    echo "----------------------------------------------"
    echo "SW Engineer (same transcript) makes tool call again:"
    test_pretooluse_detection "$SW_TRANSCRIPT"
}

# Test 5: Proposed new hook structure
test_proposed_structure() {
    echo -e "${BLUE}TEST 5: Proposed Hook Structure${NC}"
    echo "========================================"
    
    cat << 'EOF'
PROPOSED PRECOMPACT HOOK:
-------------------------
#!/bin/bash
TRANSCRIPT_PATH="$1"
AGENT_TYPE="$2"  # Could be derived from path or passed

# Extract unique transcript ID
TRANSCRIPT_ID=$(basename "$TRANSCRIPT_PATH" .jsonl)

# Create agent-specific marker using transcript ID
MARKER="/tmp/compaction_marker_${TRANSCRIPT_ID}.txt"

# Save recovery information
cat > "$MARKER" << MARKER_CONTENT
COMPACTION_TIME:$(date -u +"%Y-%m-%d %H:%M:%S UTC")
TRANSCRIPT_ID:$TRANSCRIPT_ID
TRANSCRIPT_PATH:$TRANSCRIPT_PATH
AGENT_TYPE:$AGENT_TYPE
# ... other recovery data
MARKER_CONTENT

PROPOSED PRETOOLUSE HOOK:
--------------------------
#!/bin/bash
TOOL_NAME="$1"
TRANSCRIPT_PATH="$2"  # Or derived from environment

# Extract transcript ID
TRANSCRIPT_ID=$(basename "$TRANSCRIPT_PATH" .jsonl)

# Check for THIS transcript's marker
MARKER="/tmp/compaction_marker_${TRANSCRIPT_ID}.txt"

if [ -f "$MARKER" ]; then
    # This specific agent/transcript was compacted
    AGENT_TYPE=$(grep "AGENT_TYPE:" "$MARKER" | cut -d: -f2)
    
    # Block and require recovery
    echo "🔴 Agent $AGENT_TYPE (transcript: $TRANSCRIPT_ID) needs recovery!"
    
    # Only THIS agent is blocked, others continue working
    exit 1
fi

# No marker = this agent wasn't compacted, proceed normally
EOF
    
    echo ""
    echo -e "${GREEN}Benefits of Transcript-Based Identification:${NC}"
    echo "1. Precise agent identification - no false positives"
    echo "2. Multiple agents can run simultaneously"
    echo "3. Only compacted agent needs recovery"
    echo "4. Other agents continue unaffected"
    echo "5. Automatic cleanup when transcript changes"
}

# Run all tests
main() {
    test_extract_transcript_id
    echo ""
    
    test_full_workflow
    echo ""
    
    test_proposed_structure
    echo ""
    
    echo "========================================"
    echo -e "${GREEN}TEST COMPLETE${NC}"
    echo "========================================"
    echo ""
    echo "CONCLUSION:"
    echo "-----------"
    echo "✓ Transcript paths CAN uniquely identify agents"
    echo "✓ Markers can be named with transcript IDs"
    echo "✓ PreToolUse can detect its specific marker"
    echo "✓ This enables precise, per-agent recovery"
    echo ""
    echo "RECOMMENDATION:"
    echo "---------------"
    echo "Update hooks to use transcript-based identification"
    echo "for precise agent-specific compaction recovery."
}

# Run if executed directly
if [ "${BASH_SOURCE[0]}" == "${0}" ]; then
    main "$@"
fi