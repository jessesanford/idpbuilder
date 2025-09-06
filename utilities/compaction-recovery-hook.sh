#!/bin/bash
# Compaction Recovery Hook for PreToolUse
# This script checks for compaction markers and forces recovery before tool execution

# Check if compaction marker exists
if [ -f /tmp/compaction_marker.txt ]; then
    # Compaction detected! Block the tool and force recovery
    
    # Read the marker contents
    COMPACTION_TIME=$(grep "COMPACTION_TIME:" /tmp/compaction_marker.txt | cut -d':' -f2-)
    WORKING_DIR=$(grep "WORKING_DIR:" /tmp/compaction_marker.txt | cut -d':' -f2)
    BRANCH=$(grep "GIT_BRANCH:" /tmp/compaction_marker.txt | cut -d':' -f2)
    TODO_STATE=$(grep "TODO_STATE_SAVED:" /tmp/compaction_marker.txt | cut -d':' -f2)
    
    # Output the decision control JSON to BLOCK the tool
    cat << EOF
{
  "decision": "block",
  "message": "🔴🔴🔴 CRITICAL: COMPACTION RECOVERY REQUIRED! 🔴🔴🔴\n\n⚠️ AUTO-COMPACTION DETECTED - Context was reset!\nCompaction Time: ${COMPACTION_TIME}\nPrevious Directory: ${WORKING_DIR}\nPrevious Branch: ${BRANCH}\nSaved TODOs: ${TODO_STATE}\n\n📋 MANDATORY RECOVERY STEPS (per CLAUDE.md rules):\n\n1. CHECK COMPACTION MARKER:\n   bash \$HOME/.claude/utilities/check-compaction.sh\n\n2. IF TODOS WERE SAVED:\n   - Read the saved TODO file: ${TODO_STATE}\n   - Load them into TodoWrite tool (not just read!)\n   - Verify TodoWrite contains recovered tasks\n\n3. RE-READ CRITICAL CONTEXT:\n   - Read your instruction files based on your agent type\n   - Read current state files\n   - Read any in-progress work files\n\n4. CLEAN UP MARKER:\n   rm -f /tmp/compaction_marker.txt\n\n5. THEN RETRY YOUR ORIGINAL ACTION\n\n⚠️ YOU MUST COMPLETE RECOVERY BEFORE PROCEEDING!\nThis tool execution has been BLOCKED until you recover from compaction."
}
EOF
    
    exit 0
fi

# No compaction marker - allow tool to proceed
echo '{"decision": "allow"}'