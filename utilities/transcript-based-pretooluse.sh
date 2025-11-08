#!/bin/bash

# transcript-based-pretooluse.sh
# PreToolUse hook that blocks tool execution if compaction marker exists for current transcript

set -euo pipefail

# Read the event payload from stdin (PreToolUse hooks receive JSON via stdin)
EVENT_PAYLOAD=$(cat)

# Extract transcript_path and tool_name from the JSON payload using Python for reliability
TOOL_INFO=$(echo "$EVENT_PAYLOAD" | python3 -c "
import json
import sys
try:
    data = json.load(sys.stdin)
    transcript = data.get('transcript_path', '')
    tool = data.get('tool_name', '')
    print(f'{transcript}|{tool}')
except:
    print('|')
" 2>/dev/null || echo "|")

# Parse the output
TRANSCRIPT_PATH=$(echo "$TOOL_INFO" | cut -d'|' -f1)
TOOL_NAME=$(echo "$TOOL_INFO" | cut -d'|' -f2)

# If we can't determine the transcript, fall back to generic marker
if [ -z "$TRANSCRIPT_PATH" ] || [ "$TRANSCRIPT_PATH" = "" ]; then
    MARKER_FILE="/tmp/compaction_marker.txt"
    TRANSCRIPT_ID="generic"
else
    # Extract unique transcript ID from path
    TRANSCRIPT_ID=$(basename "$TRANSCRIPT_PATH" .jsonl)
    MARKER_FILE="/tmp/compaction_marker_${TRANSCRIPT_ID}.txt"
fi

# Function to detect agent from transcript
detect_agent_from_transcript() {
    local transcript_path="$1"
    local agent_type="unknown"
    
    if [ -f "$transcript_path" ]; then
        # Search for subagent_type in the JSONL file (most recent occurrence)
        local agent_json_line=$(tac "$transcript_path" 2>/dev/null | grep -m1 '"subagent_type"' || true)
        
        if [ -n "$agent_json_line" ]; then
            # Extract subagent_type value
            local subagent=$(echo "$agent_json_line" | grep -o '"subagent_type":"[^"]*"' | cut -d'"' -f4 || echo "")
            
            # Map subagent_type to our agent names
            case "$subagent" in
                *orchestrator*) agent_type="orchestrator" ;;
                *software*engineer*|*sw*eng*) agent_type="sw-engineer" ;;
                *code*reviewer*) agent_type="code-reviewer" ;;
                *architect*) agent_type="architect" ;;
                *) agent_type="$subagent" ;;
            esac
        fi
    fi
    
    # Fallback: try pattern matching on working directory
    if [ "$agent_type" = "unknown" ]; then
        if [ -f "orchestrator-state-v3.json" ] || [[ "$PWD" == *orchestrator* ]]; then
            agent_type="orchestrator"
        elif [ -f "IMPLEMENTATION-PLAN.md" ] || [[ "$PWD" == *effort* ]]; then
            agent_type="sw-engineer"
        elif [ -f "REVIEW-FEEDBACK.md" ] || [[ "$PWD" == *review* ]]; then
            agent_type="code-reviewer"
        elif [ -f "ARCHITECT-ASSESSMENT.md" ] || [[ "$PWD" == *architect* ]]; then
            agent_type="architect"
        fi
    fi
    
    echo "$agent_type"
}

# Check if this is a TodoWrite operation - always allow these
if [ "$TOOL_NAME" = "TodoWrite" ]; then
    # Always allow TodoWrite operations for recovery
    echo '{"decision": "approve"}'
    exit 0
fi

# Main check for compaction marker
if [ -f "$MARKER_FILE" ]; then
    # Extract recovery information
    agent_type=$(grep "^AGENT_TYPE:" "$MARKER_FILE" 2>/dev/null | cut -d: -f2 || echo "unknown")
    compaction_time=$(grep "^COMPACTION_TIME:" "$MARKER_FILE" 2>/dev/null | cut -d: -f2- || echo "unknown")
    todos_saved=$(grep "^TODO_STATE_SAVED:" "$MARKER_FILE" 2>/dev/null | cut -d: -f2 || echo "")
    
    # Check for NO_TODOS_FOUND marker
    no_todos_found=$(grep "^NO_TODOS_FOUND" "$MARKER_FILE" 2>/dev/null || echo "")
    
    # Build recovery message
    recovery_msg="🚨 CONTEXT COMPACTION DETECTED FOR THIS AGENT 🚨\\n\\n"
    recovery_msg+="Agent Type: ${agent_type}\\n"
    recovery_msg+="Transcript ID: ${TRANSCRIPT_ID}\\n"
    recovery_msg+="Compaction Time: ${compaction_time}\\n\\n"
    
    if [ -n "$no_todos_found" ]; then
        # No TODOs were found during compaction
        recovery_msg+="ℹ️ No TODO files were found during compaction.\\n"
        recovery_msg+="This is normal if you had not saved any TODOs yet.\\n\\n"
        recovery_msg+="Please check if you have TODO files in:\\n"
        recovery_msg+="- /workspaces/todos/\\n"
        recovery_msg+="- /workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/todos/\\n\\n"
        recovery_msg+="If TODOs exist, load them. Otherwise, reconstruct state from:\\n"
        recovery_msg+="- orchestrator-state-v3.json\\n"
        recovery_msg+="- Recent work logs\\n"
        recovery_msg+="- Git history\\n"
    elif [ -n "$todos_saved" ] && [ -f "$todos_saved" ]; then
        recovery_msg+="MANDATORY RECOVERY STEPS:\\n"
        recovery_msg+="1. Read the saved TODO state: ${todos_saved}\\n"
        recovery_msg+="2. Load TODOs into TodoWrite tool (CRITICAL - not just read!)\\n"
        recovery_msg+="3. Read your context files based on agent type\\n"
        recovery_msg+="4. Resume work from the state indicated in your TODOs\\n"
    else
        recovery_msg+="⚠️ WARNING: Compaction occurred but no TODOs were saved\\n\\n"
        recovery_msg+="You may need to reconstruct state from:\\n"
        recovery_msg+="- orchestrator-state-v3.json\\n"
        recovery_msg+="- Recent work logs\\n"
        recovery_msg+="- Git history\\n"
    fi
    
    # Clean up the marker
    if [ -f "$MARKER_FILE" ]; then
        archive_dir="/tmp/compaction-archives"
        mkdir -p "$archive_dir"
        archive_name="${archive_dir}/$(basename "$MARKER_FILE" .txt)-$(date +%Y%m%d-%H%M%S).txt"
        mv "$MARKER_FILE" "$archive_name" 2>/dev/null || true
    fi
    
    # Output JSON decision to block
    cat <<EOF
{
  "decision": "block",
  "message": "${recovery_msg}"
}
EOF
    exit 0
else
    # No marker for this transcript - proceed normally
    echo '{"decision": "approve"}'
    exit 0
fi