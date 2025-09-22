#!/bin/bash
################################################################################
# TEST HOOK INTEGRATION
################################################################################
# This script tests that the transcript-based hooks are properly configured
# and will be called when compaction or tool use occurs.
################################################################################

set -euo pipefail

echo "================================"
echo "TESTING HOOK INTEGRATION"
echo "================================"
echo ""

# Clean up any existing markers
echo "1. Cleaning up old markers..."
rm -f /tmp/compaction_marker*.txt
rm -f /tmp/todos-precompact*.txt
echo "✓ Cleanup complete"
echo ""

# Test PreCompact hook directly
echo "2. Testing PreCompact hook directly..."
echo "   Simulating manual compaction..."

# The hook should be configured in settings.json
TRANSCRIPT_PATH="/home/vscode/.claude/test/test-agent-$(date +%s).jsonl"
export CLAUDE_TRANSCRIPT_PATH="$TRANSCRIPT_PATH"

# Call the hook directly as it would be called by Claude Code
CLAUDE_PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/workspaces/software-factory-2.0-template}"
# The hook might exit non-zero if no TODOs found, but that's okay for testing
bash $CLAUDE_PROJECT_DIR/utilities/transcript-based-precompact.sh "$TRANSCRIPT_PATH" "manual" 2>&1 | head -20
echo "✓ PreCompact hook executed (may not find TODOs in test environment)"

# Check if marker was created
TRANSCRIPT_ID=$(basename "$TRANSCRIPT_PATH" .jsonl)
MARKER_FILE="/tmp/compaction_marker_${TRANSCRIPT_ID}.txt"

if [ -f "$MARKER_FILE" ]; then
    echo "✓ Marker file created: $MARKER_FILE"
    echo ""
    echo "Marker contents:"
    head -5 "$MARKER_FILE"
else
    echo "✗ Marker file not created"
    exit 1
fi
echo ""

# Test PreToolUse hook
echo "3. Testing PreToolUse hook..."
echo "   Simulating tool call after compaction..."

# This should detect the marker and block
if bash $CLAUDE_PROJECT_DIR/utilities/transcript-based-pretooluse.sh "Read" "$TRANSCRIPT_PATH"; then
    echo "✗ ERROR: Tool call was not blocked (should have been)"
    exit 1
else
    echo "✓ Tool call was blocked as expected (recovery needed)"
fi
echo ""

# Test with different transcript (should not block)
echo "4. Testing different agent (should not block)..."
DIFFERENT_TRANSCRIPT="/home/vscode/.claude/test/different-agent-$(date +%s).jsonl"
export CLAUDE_TRANSCRIPT_PATH="$DIFFERENT_TRANSCRIPT"

if bash $CLAUDE_PROJECT_DIR/utilities/transcript-based-pretooluse.sh "Read" "$DIFFERENT_TRANSCRIPT"; then
    echo "✓ Different agent not blocked (correct behavior)"
else
    echo "✗ ERROR: Different agent was blocked (should not have been)"
    exit 1
fi
echo ""

# Clean up
echo "5. Cleaning up test artifacts..."
rm -f /tmp/compaction_marker*.txt
rm -f /tmp/todos-precompact*.txt
echo "✓ Cleanup complete"
echo ""

echo "================================"
echo "✅ HOOK INTEGRATION TEST PASSED"
echo "================================"
echo ""
echo "The transcript-based hooks are properly configured and working!"
echo ""
echo "Next steps to test with real compaction:"
echo "1. Use Claude Code normally until you approach context limit"
echo "2. Or trigger manual compaction with: /compact"
echo "3. The hooks should automatically:"
echo "   - Save your TODO state"
echo "   - Create agent-specific marker"
echo "   - Block only THIS agent on next tool use"
echo "   - Guide recovery process"