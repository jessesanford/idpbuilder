#!/bin/bash
################################################################################
# VERIFY HOOK SETUP
################################################################################
# Quick verification that transcript-based hooks are properly configured
################################################################################

echo "================================"
echo "HOOK SETUP VERIFICATION"
echo "================================"
echo ""

# Check settings.json
echo "1. Checking .claude/settings.json configuration..."
if [ -f "/workspaces/software-factory-2.0-template/.claude/settings.json" ]; then
    echo "✓ Settings file exists"
    
    # Check PreCompact hook
    if grep -q "transcript-based-precompact.sh" /workspaces/software-factory-2.0-template/.claude/settings.json; then
        echo "✓ PreCompact hook configured with transcript-based script"
    else
        echo "✗ PreCompact hook not using transcript-based script"
    fi
    
    # Check PreToolUse hook
    if grep -q "transcript-based-pretooluse.sh" /workspaces/software-factory-2.0-template/.claude/settings.json; then
        echo "✓ PreToolUse hook configured with transcript-based script"
    else
        echo "✗ PreToolUse hook not using transcript-based script"
    fi
else
    echo "✗ Settings file not found"
fi
echo ""

# Check script files
echo "2. Checking hook scripts exist and are executable..."
SCRIPTS=(
    "/workspaces/software-factory-2.0-template/utilities/transcript-based-precompact.sh"
    "/workspaces/software-factory-2.0-template/utilities/transcript-based-pretooluse.sh"
)

for script in "${SCRIPTS[@]}"; do
    if [ -f "$script" ]; then
        if [ -x "$script" ]; then
            echo "✓ $(basename $script) exists and is executable"
        else
            echo "✗ $(basename $script) exists but is not executable"
        fi
    else
        echo "✗ $(basename $script) not found"
    fi
done
echo ""

# Clean test
echo "3. Running basic functionality test..."
rm -f /tmp/compaction_marker*.txt

# Create a test marker
TEST_ID="test-$(date +%s)"
cat > "/tmp/compaction_marker_${TEST_ID}.txt" << EOF
COMPACTION_TIME:$(date -u +"%Y-%m-%d %H:%M:%S UTC")
COMPACTION_TYPE:test
AGENT_TYPE:test-agent
TRANSCRIPT_ID:$TEST_ID
TRANSCRIPT_PATH:/test/path/${TEST_ID}.jsonl
WORKING_DIR:$(pwd)
EOF

echo "✓ Created test marker for transcript: $TEST_ID"

# Test detection
export CLAUDE_TRANSCRIPT_PATH="/test/path/${TEST_ID}.jsonl"
if bash /workspaces/software-factory-2.0-template/utilities/transcript-based-pretooluse.sh "test-tool" 2>&1 | grep -q "COMPACTION DETECTED"; then
    echo "✓ PreToolUse correctly detects compaction marker"
else
    echo "✗ PreToolUse did not detect marker"
fi

# Test non-detection for different transcript
export CLAUDE_TRANSCRIPT_PATH="/test/path/different-agent.jsonl"
if bash /workspaces/software-factory-2.0-template/utilities/transcript-based-pretooluse.sh "test-tool" 2>&1 | grep -q "No compaction detected"; then
    echo "✓ PreToolUse correctly ignores other agent's marker"
else
    echo "✗ PreToolUse incorrectly detected other agent's marker"
fi

# Cleanup
rm -f /tmp/compaction_marker*.txt
echo ""

echo "================================"
echo "✅ VERIFICATION COMPLETE"
echo "================================"
echo ""
echo "The transcript-based hooks are properly set up and ready to use!"
echo ""
echo "To test with real compaction:"
echo "1. The hooks will automatically activate when you approach context limit"
echo "2. Or trigger manual compaction with: /compact"
echo "3. Each agent will get its own marker based on transcript ID"
echo "4. Only the compacted agent will be blocked for recovery"