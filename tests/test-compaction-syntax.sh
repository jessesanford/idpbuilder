#!/bin/bash

# Test that all agent compaction detection blocks have correct syntax

echo "Testing compaction detection syntax in agent files..."
echo "=================================================="

# Test function
test_compaction_block() {
    local AGENT="$1"
    local FILE="$2"
    
    echo ""
    echo "Testing $AGENT compaction detection..."
    
    # Extract the compaction detection block
    sed -n '/if \[ -f \/tmp\/compaction_marker\.txt \]/,/^fi$/p' "$FILE" > /tmp/test-block-$$.sh
    
    # Check if block was found
    if [ ! -s /tmp/test-block-$$.sh ]; then
        echo "❌ No compaction detection block found!"
        return 1
    fi
    
    # Convert to single line (as agents would run it)
    SINGLE_LINE=$(cat /tmp/test-block-$$.sh | tr '\n' ' ')
    
    # Test with marker present
    touch /tmp/test_compaction_marker.txt
    echo "TODO_STATE_SAVED:test" > /tmp/test_compaction_marker.txt
    
    # Run the command (redirecting output to avoid clutter)
    if bash -c "$SINGLE_LINE" >/dev/null 2>&1; then
        echo "✅ Syntax valid with marker present"
    else
        echo "❌ Syntax error with marker present!"
        echo "Command: $SINGLE_LINE"
        return 1
    fi
    
    # Test without marker (else clause)
    rm -f /tmp/test_compaction_marker.txt
    
    if bash -c "$SINGLE_LINE" >/dev/null 2>&1; then
        echo "✅ Syntax valid without marker (else clause works)"
    else
        echo "❌ Syntax error without marker!"
        return 1
    fi
    
    # Clean up
    rm -f /tmp/test-block-$$.sh
    
    return 0
}

# Test each agent file
AGENTS=(
    "Orchestrator:/workspaces/software-factory-2.0-template/.claude/agents/orchestrator.md"
    "SW-Engineer:/workspaces/software-factory-2.0-template/.claude/agents/sw-engineer.md"
    "Code-Reviewer:/workspaces/software-factory-2.0-template/.claude/agents/code-reviewer.md"
    "Architect:/workspaces/software-factory-2.0-template/.claude/agents/architect.md"
)

FAILED=0

for AGENT_SPEC in "${AGENTS[@]}"; do
    IFS=':' read -r AGENT FILE <<< "$AGENT_SPEC"
    if ! test_compaction_block "$AGENT" "$FILE"; then
        FAILED=$((FAILED + 1))
    fi
done

echo ""
echo "=================================================="
if [ $FAILED -eq 0 ]; then
    echo "✅ All agent compaction detection blocks have valid syntax!"
else
    echo "❌ $FAILED agent(s) have syntax errors in compaction detection"
    exit 1
fi

# Also test the main CLAUDE.md file
echo ""
echo "Testing main CLAUDE.md compaction detection..."

# Extract and test the block from CLAUDE.md
CLAUDE_MD="/home/vscode/.claude/CLAUDE.md"
if [ -f "$CLAUDE_MD" ]; then
    # Extract the bash block
    sed -n '/^if \[ -f \/tmp\/compaction_marker\.txt \]/,/^fi$/p' "$CLAUDE_MD" > /tmp/claude-block-$$.sh
    
    if [ -s /tmp/claude-block-$$.sh ]; then
        # Convert to single line
        SINGLE_LINE=$(cat /tmp/claude-block-$$.sh | tr '\n' ' ')
        
        # Test execution
        touch /tmp/test_compaction_marker.txt
        echo "TODO_STATE_SAVED:test" > /tmp/test_compaction_marker.txt
        
        if bash -c "$SINGLE_LINE" >/dev/null 2>&1; then
            echo "✅ CLAUDE.md syntax valid with marker"
        else
            echo "❌ CLAUDE.md syntax error with marker!"
        fi
        
        rm -f /tmp/test_compaction_marker.txt
        
        if bash -c "$SINGLE_LINE" >/dev/null 2>&1; then
            echo "✅ CLAUDE.md syntax valid without marker"
        else
            echo "❌ CLAUDE.md syntax error without marker!"
        fi
    else
        echo "⚠️ No compaction block found in CLAUDE.md"
    fi
    
    rm -f /tmp/claude-block-$$.sh
else
    echo "⚠️ CLAUDE.md not found at expected location"
fi

# Clean up any test markers
rm -f /tmp/test_compaction_marker.txt