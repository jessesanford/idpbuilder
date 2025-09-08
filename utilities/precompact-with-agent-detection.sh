#!/bin/bash
# PreCompact hook that detects agent type from transcript

# Read the PreCompact input from stdin
PRECOMPACT_INPUT=$(cat)

# Extract the transcript_path and compaction_type from the JSON input
TRANSCRIPT_PATH=$(echo "$PRECOMPACT_INPUT" | python3 -c "
import json
import sys
try:
    data = json.load(sys.stdin)
    print(data.get('transcript_path', ''))
except:
    print('')
")

COMPACTION_TYPE=$(echo "$PRECOMPACT_INPUT" | python3 -c "
import json
import sys
try:
    data = json.load(sys.stdin)
    print(data.get('compaction_type', 'auto'))
except:
    print('auto')
")

# Detect agent type from transcript
AGENT_TYPE="unknown"
if [ -n "$TRANSCRIPT_PATH" ] && [ -f "$TRANSCRIPT_PATH" ]; then
    # Search from bottom to top for subagent_type in the JSONL file
    # Using tac to reverse file, then grep to find first (most recent) match
    AGENT_JSON_LINE=$(tac "$TRANSCRIPT_PATH" 2>/dev/null | grep -m1 '"subagent_type"' || true)
    
    if [ -n "$AGENT_JSON_LINE" ]; then
        # Extract agent type from JSON using sed
        # Pattern: "subagent_type":"orchestrator"
        DETECTED_TYPE=$(echo "$AGENT_JSON_LINE" | sed -n 's/.*"subagent_type"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')
        
        # Map the subagent types to our standard names
        case "$DETECTED_TYPE" in
            "orchestrator")
                AGENT_TYPE="orchestrator"
                ;;
            "software-engineer")
                AGENT_TYPE="sw-engineer"
                ;;
            "code-reviewer")
                AGENT_TYPE="code-reviewer"
                ;;
            "architect")
                AGENT_TYPE="architect"
                ;;
            *)
                # If we found something but it's not recognized, fall back to pattern matching
                if [ -z "$DETECTED_TYPE" ]; then
                    # Fallback: Look for agent indicators in the last 100 lines
                    TRANSCRIPT_TAIL=$(tail -100 "$TRANSCRIPT_PATH" 2>/dev/null)
                    
                    # Check for orchestrator patterns (most specific first)
                    if echo "$TRANSCRIPT_TAIL" | grep -q "orchestrator\|ORCHESTRATOR\|continue-orchestrating\|orchestrator-state.json\|/orchestrator"; then
                        AGENT_TYPE="orchestrator"
                    # Check for code reviewer patterns
                    elif echo "$TRANSCRIPT_TAIL" | grep -q "code.reviewer\|CODE.REVIEWER\|code review\|Starting code review\|Reviewing code\|/reviewer"; then
                        AGENT_TYPE="code-reviewer"
                    # Check for architect patterns
                    elif echo "$TRANSCRIPT_TAIL" | grep -q "architect\|ARCHITECT\|architecture.review\|wave review\|phase assessment\|/architect"; then
                        AGENT_TYPE="architect"
                    # Check for SW engineer patterns
                    elif echo "$TRANSCRIPT_TAIL" | grep -q "sw.eng\|SW.ENGINEER\|SOFTWARE.ENGINEER\|Starting implementation\|pkg/.*\.go"; then
                        AGENT_TYPE="sw-engineer"
                    fi
                else
                    # Use the detected type even if not in our standard list
                    AGENT_TYPE="$DETECTED_TYPE"
                fi
                ;;
        esac
    else
        # No subagent_type found, fall back to pattern matching
        TRANSCRIPT_TAIL=$(tail -100 "$TRANSCRIPT_PATH" 2>/dev/null)
        
        if echo "$TRANSCRIPT_TAIL" | grep -q "orchestrator\|ORCHESTRATOR\|continue-orchestrating\|orchestrator-state.json\|/orchestrator"; then
            AGENT_TYPE="orchestrator"
        elif echo "$TRANSCRIPT_TAIL" | grep -q "code.reviewer\|CODE.REVIEWER\|code review\|Starting code review\|Reviewing code\|/reviewer"; then
            AGENT_TYPE="code-reviewer"
        elif echo "$TRANSCRIPT_TAIL" | grep -q "architect\|ARCHITECT\|architecture.review\|wave review\|phase assessment\|/architect"; then
            AGENT_TYPE="architect"
        elif echo "$TRANSCRIPT_TAIL" | grep -q "sw.eng\|SW.ENGINEER\|SOFTWARE.ENGINEER\|Starting implementation\|pkg/.*\.go"; then
            AGENT_TYPE="sw-engineer"
        fi
    fi
fi

# Create agent-specific compaction marker file
MARKER_FILE="/tmp/compaction_marker_${AGENT_TYPE}.txt"
echo "🔴🔴🔴 ${COMPACTION_TYPE^^} COMPACTION DETECTED FOR $AGENT_TYPE 🔴🔴🔴"
echo "COMPACTION_TIME:$(date '+%Y-%m-%d %H:%M:%S %Z')" > "$MARKER_FILE"
echo "COMPACTION_TYPE:$COMPACTION_TYPE" >> "$MARKER_FILE"
echo "AGENT_TYPE:$AGENT_TYPE" >> "$MARKER_FILE"
echo "WORKING_DIR:$(pwd)" >> "$MARKER_FILE"
echo "GIT_BRANCH:$(git branch --show-current 2>/dev/null || echo 'none')" >> "$MARKER_FILE"
echo "TRANSCRIPT_PATH:$TRANSCRIPT_PATH" >> "$MARKER_FILE"

# Try to save TODOs based on agent type
# Use CLAUDE_PROJECT_DIR if set, otherwise use default
PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/workspaces/software-factory-2.0-template}"
TODO_DIR="${PROJECT_DIR}/todos"

# Check if project todos directory exists
if [ -d "$TODO_DIR" ] && [ "$(ls -A $TODO_DIR 2>/dev/null)" ]; then
    # Look for agent-specific TODO files
    AGENT_PATTERN="${TODO_DIR}/${AGENT_TYPE}-*.todo"
    LATEST_TODO=$(ls -t $AGENT_PATTERN 2>/dev/null | head -1)
    
    if [ -f "$LATEST_TODO" ]; then
        echo "TODO_STATE_SAVED:/tmp/todos-precompact-$AGENT_TYPE.txt" >> "$MARKER_FILE"
        cp "$LATEST_TODO" "/tmp/todos-precompact-$AGENT_TYPE.txt"
        echo "✅ Saved TODO for $AGENT_TYPE: $(basename "$LATEST_TODO")"
    else
        # No agent-specific TODO found, try to find any recent TODO
        LATEST_TODO=$(ls -t $TODO_DIR/*.todo 2>/dev/null | head -1)
        if [ -f "$LATEST_TODO" ]; then
            echo "TODO_STATE_SAVED:/tmp/todos-precompact-$AGENT_TYPE.txt" >> "$MARKER_FILE"
            cp "$LATEST_TODO" "/tmp/todos-precompact-$AGENT_TYPE.txt"
            echo "✅ Saved TODO (non-agent-specific): $(basename "$LATEST_TODO")"
        else
            echo "NO_TODOS_FOUND" >> "$MARKER_FILE"
            echo "⚠️ No TODO files found in $TODO_DIR"
        fi
    fi
else
    echo "NO_TODOS_FOUND" >> "$MARKER_FILE"
    echo "⚠️ TODO directory does not exist or is empty: $TODO_DIR"
fi

echo "✅ Compaction marker created for $AGENT_TYPE at $MARKER_FILE"