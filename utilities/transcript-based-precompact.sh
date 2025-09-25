#!/bin/bash
################################################################################
# SOFTWARE FACTORY 2.0 - TRANSCRIPT-BASED PRE-COMPACTION HOOK
################################################################################
#
# This improved version uses transcript IDs for precise agent identification
# ensuring only the compacted agent needs recovery while others continue.
#
# IMPROVEMENTS:
# - Uses transcript path as unique agent identifier
# - Creates agent-specific marker files using transcript ID
# - Enables multiple agents to run simultaneously
# - Prevents false-positive recovery requirements
# - Detects agent type from transcript subagent_type field
#
################################################################################

set -euo pipefail

# Read the event payload from stdin (PreCompact hooks receive JSON via stdin)
EVENT_PAYLOAD=$(cat)

# Extract transcript_path and type from JSON payload using Python for reliability
TRANSCRIPT_PATH=$(echo "$EVENT_PAYLOAD" | python3 -c "
import json
import sys
try:
    data = json.load(sys.stdin)
    print(data.get('transcript_path', ''))
except:
    print('')
" 2>/dev/null || echo "")

COMPACTION_TYPE=$(echo "$EVENT_PAYLOAD" | python3 -c "
import json
import sys
try:
    data = json.load(sys.stdin)
    print(data.get('type', 'auto'))
except:
    print('auto')
" 2>/dev/null || echo "auto")

# Debug: Show what we extracted
echo "📝 Extracted from PreCompact event:"
echo "   - Transcript path: $TRANSCRIPT_PATH"
echo "   - Compaction type: $COMPACTION_TYPE"

# If we can't determine the transcript, fall back to generic marker
if [ -z "$TRANSCRIPT_PATH" ] || [ "$TRANSCRIPT_PATH" = "" ]; then
    echo "⚠️ Warning: No transcript path available, using generic marker"
    MARKER_FILE="/tmp/compaction_marker.txt"
    TRANSCRIPT_ID="generic"
else
    # Extract unique transcript ID from path
    TRANSCRIPT_ID=$(basename "$TRANSCRIPT_PATH" .jsonl)
    MARKER_FILE="/tmp/compaction_marker_${TRANSCRIPT_ID}.txt"
    echo "📝 Using transcript-based marker: $MARKER_FILE"
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
    
    echo "$agent_type"
}

# Try to detect agent type from various sources
detect_agent_type() {
    local agent_type="unknown"
    
    # First try to detect from transcript
    if [ -n "$TRANSCRIPT_PATH" ] && [ -f "$TRANSCRIPT_PATH" ]; then
        agent_type=$(detect_agent_from_transcript "$TRANSCRIPT_PATH")
    fi
    
    # If still unknown, fall back to filesystem detection
    if [ "$agent_type" = "unknown" ]; then
        # Method 1: Check if orchestrator-state.json exists (orchestrator agent)
        if [ -f "orchestrator-state.json" ]; then
            agent_type="orchestrator"
        
        # Method 2: Check for IMPLEMENTATION-PLAN.md (SW engineer)
        elif [ -f "IMPLEMENTATION-PLAN.md" ] || [ -f "work-log.md" ]; then
            agent_type="sw-engineer"
        
        # Method 3: Check for REVIEW files (code reviewer)
        elif ls REVIEW*.md 2>/dev/null | head -1 > /dev/null; then
            agent_type="code-reviewer"
        
        # Method 4: Check for architecture review files
        elif ls *ARCHITECT*.md 2>/dev/null | head -1 > /dev/null; then
            agent_type="architect"
        
        # Method 5: Check if we're in an effort directory
        elif pwd | grep -q "efforts/phase[0-9]*/wave[0-9]*"; then
            if [ -f "../../../orchestrator-state.json" ]; then
                agent_type="sw-engineer"
            else
                agent_type="code-reviewer"
            fi
        
        # Method 6: Check for agent-specific TODO files
        elif ls todos/orchestrator-*.todo 2>/dev/null | head -1 > /dev/null; then
            agent_type="orchestrator"
        elif ls todos/sw-eng-*.todo 2>/dev/null | head -1 > /dev/null; then
            agent_type="sw-engineer"
        elif ls todos/code-reviewer-*.todo 2>/dev/null | head -1 > /dev/null; then
            agent_type="code-reviewer"
        elif ls todos/architect-*.todo 2>/dev/null | head -1 > /dev/null; then
            agent_type="architect"
        fi
    fi
    
    echo "$agent_type"
}

AGENT_TYPE=$(detect_agent_type)
echo "🤖 Detected agent type: $AGENT_TYPE"

# Save TODO state based on agent type
save_agent_todos() {
    local agent="$1"
    local todos_file="/tmp/todos-precompact-${TRANSCRIPT_ID}.txt"
    
    echo "💾 Saving TODO state for $agent (transcript: $TRANSCRIPT_ID)"
    
    # Use CLAUDE_PROJECT_DIR to find todos
    local todos_dir="${CLAUDE_PROJECT_DIR}/todos"
    
    if [ ! -d "$todos_dir" ]; then
        echo "  TODO directory not found at: $todos_dir"
        echo "  Trying alternate location: /workspaces/todos"
        todos_dir="/workspaces/todos"
    fi
    
    if [ ! -d "$todos_dir" ]; then
        echo "⚠️ No TODO directory found"
        echo "NO_TODOS_FOUND" >> "$MARKER_FILE"
        return
    fi
    
    echo "  Searching in: $todos_dir"
    
    # Find the most recent TODO file for this agent
    local latest_todo=$(ls -t "$todos_dir"/${agent}-*.todo 2>/dev/null | head -1)
    
    if [ -n "$latest_todo" ] && [ -f "$latest_todo" ]; then
        cp "$latest_todo" "$todos_file"
        echo "TODO_STATE_SAVED:$todos_file" >> "$MARKER_FILE"
        echo "✅ Saved TODOs from $latest_todo"
    else
        echo "NO_TODOS_FOUND" >> "$MARKER_FILE"
        echo "⚠️ No TODO files found for $agent in $todos_dir"
        echo "  Files present: $(ls "$todos_dir" 2>/dev/null | wc -l)"
    fi
}

# Create the marker file with all recovery information
create_marker() {
    echo "📋 Creating compaction marker for transcript: $TRANSCRIPT_ID"
    
    cat > "$MARKER_FILE" << EOF
COMPACTION_TIME:$(date -u +"%Y-%m-%d %H:%M:%S UTC")
COMPACTION_TYPE:$COMPACTION_TYPE
AGENT_TYPE:$AGENT_TYPE
TRANSCRIPT_ID:$TRANSCRIPT_ID
TRANSCRIPT_PATH:$TRANSCRIPT_PATH
WORKING_DIR:$(pwd)
GIT_BRANCH:$(git branch --show-current 2>/dev/null || echo "unknown")
GIT_STATUS:$(git status -s 2>/dev/null | wc -l) modified files
EOF

    # Add state machine information for orchestrator
    if [ "$AGENT_TYPE" = "orchestrator" ] && [ -f "orchestrator-state.json" ]; then
        local current_state=$(grep "current_state:" orchestrator-state.json | head -1 | awk '{print $2}')
        echo "ORCHESTRATOR_STATE:$current_state" >> "$MARKER_FILE"
    fi
    
    # Add effort information for SW engineer
    if [ "$AGENT_TYPE" = "sw-engineer" ] && pwd | grep -q "efforts/phase"; then
        local effort_path=$(pwd | sed 's|.*/efforts/||')
        echo "EFFORT_PATH:$effort_path" >> "$MARKER_FILE"
    fi
}

# Main execution
main() {
    echo "========================================="
    echo "PRE-COMPACTION HOOK (Transcript-Based)"
    echo "========================================="
    echo "Time: $(date -u +"%Y-%m-%d %H:%M:%S UTC")"
    echo "Transcript ID: $TRANSCRIPT_ID"
    echo "Agent Type: $AGENT_TYPE"
    echo "Working Dir: $(pwd)"
    echo ""
    
    # Create the marker
    create_marker
    
    # Save agent-specific TODOs
    case "$AGENT_TYPE" in
        orchestrator)
            save_agent_todos "orchestrator"
            ;;
        sw-engineer|sw-eng)
            save_agent_todos "sw-eng"
            ;;
        code-reviewer)
            save_agent_todos "code-reviewer"
            ;;
        architect)
            save_agent_todos "architect"
            ;;
        unknown)
            echo "⚠️ Unknown agent type, checking for any TODO files"
            # Try to save any TODO files found
            if [ -d "${CLAUDE_PROJECT_DIR}/todos" ] || [ -d "/workspaces/todos" ]; then
                save_agent_todos "unknown"
            fi
            ;;
        *)
            # For any other detected agent type
            save_agent_todos "$AGENT_TYPE"
            ;;
    esac
    
    # Save state files if they exist
    if [ -f "orchestrator-state.json" ]; then
        cp orchestrator-state.json "/tmp/state-precompact-${TRANSCRIPT_ID}.yaml"
        echo "STATE_FILE_SAVED:/tmp/state-precompact-${TRANSCRIPT_ID}.yaml" >> "$MARKER_FILE"
        echo "✅ Saved orchestrator state"
    fi
    
    # Display summary
    echo ""
    echo "📊 Compaction Preparation Complete:"
    echo "  - Marker: $MARKER_FILE"
    echo "  - Agent: $AGENT_TYPE"
    echo "  - Transcript: $TRANSCRIPT_ID"
    
    # Show marker contents
    echo ""
    echo "📄 Marker contents:"
    cat "$MARKER_FILE"
    
    echo ""
    echo "✅ Pre-compaction hook complete for $AGENT_TYPE (transcript: $TRANSCRIPT_ID)"
    echo "Recovery will be required only for THIS specific transcript."
}

# Run main function
main "$@"