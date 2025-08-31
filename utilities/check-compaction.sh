#!/bin/bash
# Compaction detection script for agents to source or run

if [ -f /tmp/compaction_marker.txt ]; then
    echo "⚠️ AUTO-COMPACTION DETECTED - Context was compressed"
    echo "Reading compaction details..."
    cat /tmp/compaction_marker.txt
    
    # Check for TODO preservation
    if grep -q "TODO_STATE_SAVED:" /tmp/compaction_marker.txt; then
        echo "📋 TODO STATE WAS PRESERVED - Check todos directory for latest state"
        PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/workspaces/software-factory-2.0-template}"
        echo "Latest TODO files in ${PROJECT_DIR}/todos/:"
        ls -t ${PROJECT_DIR}/todos/*.todo 2>/dev/null | head -3 | while read f; do 
            echo "  - $(basename $f)"
        done
        echo "You should load and merge TODOs from the appropriate file for your agent"
    fi
    
    rm -f /tmp/compaction_marker.txt
    echo "🔄 INITIATING CONTEXT RECOVERY..."
    echo "⚠️⚠️⚠️ CRITICAL TODO RECOVERY STEPS ⚠️⚠️⚠️"
    echo "You MUST now:"
    echo "1. Identify which agent you are (check your current prompt for @agent-*)"
    echo "2. READ your TODO file: Use Read tool on latest {your-agent-name}-*.todo from todos directory"
    echo "3. LOAD INTO TODOWRITE: Use TodoWrite tool to populate your working TODO list with those items"
    echo "4. DEDUPLICATE: Merge with any TODOs already in memory (avoid duplicates)"
    echo "5. VERIFY: Confirm your TodoWrite tool now contains all recovered TODOs"
    echo "6. Determine your current mode/state from the compaction details"
    echo "7. Jump to YOUR section below (1️⃣-5️⃣) and read ONLY the files for your current mode"
    echo "8. If unsure of state, jump to section 7️⃣ CONTEXT LOSS RECOVERY"
else
    echo "No compaction detected - proceeding normally"
fi