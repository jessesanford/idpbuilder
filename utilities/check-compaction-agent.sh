#!/bin/bash
# Agent-specific compaction detection script
# Usage: bash check-compaction-agent.sh [orchestrator|sw-engineer|code-reviewer|architect]

AGENT_TYPE="${1:-unknown}"

if [ -f /tmp/compaction_marker.txt ]; then
    echo "⚠️ AUTO-COMPACTION DETECTED - Context was compressed"
    echo "Reading compaction details..."
    cat /tmp/compaction_marker.txt
    
    # Try to run post-compact utility if available
    for path in "$HOME/.claude/utilities/post-compact.sh" "/home/user/.claude/utilities/post-compact.sh" "./utilities/post-compact.sh"; do
        if [ -f "$path" ]; then
            echo "🔧 Running post-compact recovery utility..."
            "$path"
            break
        fi
    done
    
    # Determine TODO file pattern based on agent type
    case "$AGENT_TYPE" in
        orchestrator)
            TODO_PATTERN="orchestrator-*.todo"
            TODO_LOAD="orchestrator"
            ;;
        sw-engineer)
            TODO_PATTERN="sw-eng-*.todo"
            TODO_LOAD="sw-eng"
            EFFORT_NAME=$(pwd | grep -oP '(?<=efforts/phase\d+/wave\d+/)[^/]+' || echo "unknown")
            echo "Current effort: $EFFORT_NAME"
            ;;
        code-reviewer)
            TODO_PATTERN="code-reviewer-*.todo"
            TODO_LOAD="code-reviewer"
            EFFORT_NAME=$(pwd | grep -oP '(?<=efforts/phase\d+/wave\d+/)[^/]+' || echo "unknown")
            echo "Current effort: $EFFORT_NAME"
            ;;
        architect)
            TODO_PATTERN="architect-*.todo"
            TODO_LOAD="architect"
            ;;
        *)
            TODO_PATTERN="*.todo"
            TODO_LOAD="unknown"
            ;;
    esac
    
    # Check for TODO preservation
    if grep -q "TODO_STATE_SAVED:" /tmp/compaction_marker.txt; then
        echo "📋 TODO STATE WAS PRESERVED - Check todos directory"
        echo "Latest TODO files for $AGENT_TYPE:"
        ls -t todos/$TODO_PATTERN 2>/dev/null | head -3 | while read f; do
            echo "  - $(basename "$f")"
        done
        echo "CRITICAL: Must load $TODO_LOAD TODOs before proceeding"
        
        # Try to use todo-preservation utility if available
        for path in "$HOME/.claude/utilities/todo-preservation.sh" "/home/user/.claude/utilities/todo-preservation.sh" "./utilities/todo-preservation.sh"; do
            if [ -f "$path" ]; then
                echo "🔧 Using todo-preservation utility to help recovery..."
                "$path" load "$TODO_LOAD"
                break
            fi
        done
    fi
    
    rm -f /tmp/compaction_marker.txt
    echo "🔄 INITIATING CONTEXT RECOVERY..."
    echo "⚠️⚠️⚠️ CRITICAL TODO RECOVERY STEPS ⚠️⚠️⚠️"
    echo "You MUST now:"
    echo "1. READ your TODO file: todos/$TODO_PATTERN"
    echo "2. LOAD INTO TODOWRITE: Use TodoWrite tool to populate"
    echo "3. DEDUPLICATE: Merge with any existing TODOs"
    echo "4. VERIFY: Confirm TodoWrite contains recovered TODOs"
    
    # Agent-specific recovery instructions
    case "$AGENT_TYPE" in
        orchestrator)
            echo "5. READ: orchestrator-state.json to understand current state"
            echo "6. CONTINUE: From appropriate state in state machine"
            ;;
        sw-engineer)
            echo "5. READ: IMPLEMENTATION-PLAN.md in your effort directory"
            echo "6. READ: work-log.md to understand progress"
            echo "7. CONTINUE: From where you left off"
            ;;
        code-reviewer)
            echo "5. Determine if planning or reviewing from TODO context"
            echo "6. READ: IMPLEMENTATION-PLAN.md if it exists"
            echo "7. CONTINUE: From appropriate review state"
            ;;
        architect)
            echo "5. READ: orchestrator-state.json to understand context"
            echo "6. Determine review type: Wave/Phase/Integration"
            echo "7. CONTINUE: From appropriate review state"
            ;;
    esac
    
    # Check for recovery assistant
    for path in "$HOME/.claude/utilities/recovery-assistant.sh" "/home/user/.claude/utilities/recovery-assistant.sh" "./utilities/recovery-assistant.sh"; do
        if [ -f "$path" ]; then
            echo "🔧 Recovery assistant available: $path"
            break
        fi
    done
    
    echo "❌ DO NOT PROCEED until TODOs are recovered!"
    exit 0
else
    echo "✅ No compaction detected - continuing normal operation"
fi