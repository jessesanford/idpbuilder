#!/bin/bash
# Agent-Aware Compaction Recovery Hook for PreToolUse
# Reads PreToolUse input from stdin to determine current agent

# Read the PreToolUse input from stdin
PRETOOL_INPUT=$(cat)

# Extract the transcript_path from the JSON input
TRANSCRIPT_PATH=$(echo "$PRETOOL_INPUT" | python3 -c "
import json
import sys
try:
    data = json.load(sys.stdin)
    print(data.get('transcript_path', ''))
except:
    print('')
")

# Debug: Log what we found
echo "DEBUG: Transcript path: $TRANSCRIPT_PATH" >&2

# Try to determine the agent from the transcript
CURRENT_AGENT=""
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
                CURRENT_AGENT="orchestrator"
                ;;
            "software-engineer")
                CURRENT_AGENT="sw-engineer"
                ;;
            "code-reviewer")
                CURRENT_AGENT="code-reviewer"
                ;;
            "architect")
                CURRENT_AGENT="architect"
                ;;
            *)
                # If we found something but it's not recognized, fall back to pattern matching
                if [ -z "$DETECTED_TYPE" ]; then
                    # Fallback: Look for agent indicators in the last 100 lines
                    TRANSCRIPT_TAIL=$(tail -100 "$TRANSCRIPT_PATH" 2>/dev/null)
                    
                    # Check for orchestrator patterns (most specific first)
                    if echo "$TRANSCRIPT_TAIL" | grep -q "orchestrator\|ORCHESTRATOR\|continue-orchestrating\|orchestrator-state.json"; then
                        CURRENT_AGENT="orchestrator"
                    # Check for code reviewer patterns
                    elif echo "$TRANSCRIPT_TAIL" | grep -q "code.reviewer\|CODE.REVIEWER\|code review\|Starting code review\|Reviewing code"; then
                        CURRENT_AGENT="code-reviewer"
                    # Check for architect patterns
                    elif echo "$TRANSCRIPT_TAIL" | grep -q "architect\|ARCHITECT\|architecture.review\|wave review\|phase assessment"; then
                        CURRENT_AGENT="architect"
                    # Check for SW engineer patterns
                    elif echo "$TRANSCRIPT_TAIL" | grep -q "sw.eng\|SW.ENGINEER\|SOFTWARE.ENGINEER\|Starting implementation\|pkg/.*\.go"; then
                        CURRENT_AGENT="sw-engineer"
                    fi
                else
                    # Use the detected type even if not in our standard list
                    CURRENT_AGENT="$DETECTED_TYPE"
                fi
                ;;
        esac
    else
        # No subagent_type found, fall back to pattern matching
        TRANSCRIPT_TAIL=$(tail -100 "$TRANSCRIPT_PATH" 2>/dev/null)
        
        # Check for orchestrator patterns
        if echo "$TRANSCRIPT_TAIL" | grep -q "orchestrator\|ORCHESTRATOR\|continue-orchestrating\|orchestrator-state.json"; then
            CURRENT_AGENT="orchestrator"
        elif echo "$TRANSCRIPT_TAIL" | grep -q "code.reviewer\|CODE.REVIEWER\|code review\|Starting code review\|Reviewing code"; then
            CURRENT_AGENT="code-reviewer"
        elif echo "$TRANSCRIPT_TAIL" | grep -q "architect\|ARCHITECT\|architecture.review\|wave review\|phase assessment"; then
            CURRENT_AGENT="architect"
        elif echo "$TRANSCRIPT_TAIL" | grep -q "sw.eng\|SW.ENGINEER\|SOFTWARE.ENGINEER\|Starting implementation\|pkg/.*\.go"; then
            CURRENT_AGENT="sw-engineer"
        fi
        
        # Also check for agent-specific commands
        if [ -z "$CURRENT_AGENT" ]; then
            if echo "$TRANSCRIPT_TAIL" | grep -q "/orchestrator\|/continue-orchestrating"; then
                CURRENT_AGENT="orchestrator"
            fi
        fi
    fi
fi

# Debug: Log detected agent
echo "DEBUG: Detected agent: $CURRENT_AGENT" >&2

# Default to unknown if we can't detect
if [ -z "$CURRENT_AGENT" ]; then
    CURRENT_AGENT="unknown"
fi

# Check for agent-specific compaction marker
MARKER_FILE="/tmp/compaction_marker_${CURRENT_AGENT}.txt"

# Also check for legacy non-agent-specific marker
LEGACY_MARKER="/tmp/compaction_marker.txt"

if [ -f "$MARKER_FILE" ]; then
    echo "DEBUG: Found agent-specific marker for $CURRENT_AGENT" >&2
    # This marker is for us! Block and require recovery
    COMPACTION_TIME=$(grep "COMPACTION_TIME:" "$MARKER_FILE" | cut -d':' -f2-)
    WORKING_DIR=$(grep "WORKING_DIR:" "$MARKER_FILE" | cut -d':' -f2)
    BRANCH=$(grep "GIT_BRANCH:" "$MARKER_FILE" | cut -d':' -f2)
    TODO_STATE=$(grep "TODO_STATE_SAVED:" "$MARKER_FILE" | cut -d':' -f2)
    TRANSCRIPT_FROM_MARKER=$(grep "TRANSCRIPT_PATH:" "$MARKER_FILE" | cut -d':' -f2-)
elif [ -f "$LEGACY_MARKER" ]; then
    # Check if legacy marker is for us
    MARKER_AGENT=$(grep "AGENT_TYPE:" "$LEGACY_MARKER" 2>/dev/null | cut -d':' -f2 | tr -d ' ')
    
    if [ -n "$MARKER_AGENT" ] && [ "$MARKER_AGENT" != "$CURRENT_AGENT" ]; then
        # Legacy marker is for a different agent
        echo "DEBUG: Legacy marker is for $MARKER_AGENT, but we are $CURRENT_AGENT - allowing" >&2
        echo '{"decision": "approve"}'
        exit 0
    fi
    
    # Legacy marker is for us or unspecified (assume orchestrator)
    if [ -z "$MARKER_AGENT" ] && [ "$CURRENT_AGENT" != "orchestrator" ]; then
        echo '{"decision": "approve"}'
        exit 0
    fi
    
    echo "DEBUG: Legacy marker applies to us" >&2
    COMPACTION_TIME=$(grep "COMPACTION_TIME:" "$LEGACY_MARKER" | cut -d':' -f2-)
    WORKING_DIR=$(grep "WORKING_DIR:" "$LEGACY_MARKER" | cut -d':' -f2)
    BRANCH=$(grep "GIT_BRANCH:" "$LEGACY_MARKER" | cut -d':' -f2)
    TODO_STATE=$(grep "TODO_STATE_SAVED:" "$LEGACY_MARKER" | cut -d':' -f2)
    TRANSCRIPT_FROM_MARKER=$(grep "TRANSCRIPT_PATH:" "$LEGACY_MARKER" | cut -d':' -f2-)
    MARKER_FILE="$LEGACY_MARKER" # For cleanup instructions
else
    # No compaction marker for us - allow tool to proceed
    echo '{"decision": "approve"}'
    exit 0
fi
    
    # Output the decision control JSON to BLOCK the tool
    cat << EOF
{
  "decision": "block",
  "message": "🔴🔴🔴 COMPACTION RECOVERY REQUIRED FOR ${CURRENT_AGENT:-AGENT}! 🔴🔴🔴\n\n⚠️ AUTO-COMPACTION DETECTED - Your context was reset!\nAgent: ${CURRENT_AGENT:-unknown}\nCompaction Time: ${COMPACTION_TIME}\nPrevious Directory: ${WORKING_DIR}\nPrevious Branch: ${BRANCH}\nTranscript: ${TRANSCRIPT_FROM_MARKER}\nSaved TODOs: ${TODO_STATE}\n\n📋 MANDATORY RECOVERY STEPS:\n\n1. CHECK COMPACTION:\n   cat ${MARKER_FILE}\n\n2. RECOVER TODOS (if saved):\n   - Read: ${TODO_STATE}\n   - Load into TodoWrite tool (MUST use TodoWrite!)\n\n3. RE-READ YOUR CONTEXT:\n   - Your agent config file\n   - Current state files\n   - Work in progress\n\n4. CLEAN UP:\n   rm -f ${MARKER_FILE}\n\n5. RETRY your original action\n\n⚠️ This is blocking only YOUR agent (${CURRENT_AGENT:-unknown}).\nOnce you complete recovery and remove the marker, your tools will work."
}
EOF
    exit 0
fi

# No compaction marker (or not for us) - allow tool to proceed
echo '{"decision": "approve"}'