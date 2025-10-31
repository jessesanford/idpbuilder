# R174.0.0 - Context Recovery Detection

## Rule Statement
All agents MUST check for `/tmp/compaction_marker.txt` on startup to detect if context compaction occurred.

## Rationale
The marker file is the primary indicator that context was lost and recovery procedures are needed.

## Implementation

### Detection Code (Required in all agents)
```bash
# At agent startup
if [ -f /tmp/compaction_marker.txt ]; then
    echo "⚠️ CONTEXT COMPACTION DETECTED"
    cat /tmp/compaction_marker.txt
    
    # Extract key information
    COMPACTION_TIME=$(grep "COMPACTION_TIME:" /tmp/compaction_marker.txt | cut -d: -f2)
    WORKING_DIR=$(grep "WORKING_DIR:" /tmp/compaction_marker.txt | cut -d: -f2)
    GIT_BRANCH=$(grep "GIT_BRANCH:" /tmp/compaction_marker.txt | cut -d: -f2)
    
    # Check for saved TODOs
    if grep -q "TODO_STATE_SAVED" /tmp/compaction_marker.txt; then
        echo "📋 TODO state was preserved"
        # Load TODOs from todos/ directory
    fi
    
    # Clean up marker
    rm -f /tmp/compaction_marker.txt
    
    # Initiate recovery
    echo "🔄 INITIATING CONTEXT RECOVERY..."
fi
```

### Recovery Actions
1. Read marker file contents
2. Determine last known state
3. Load preserved TODOs if available
4. Read necessary context files
5. Resume from appropriate state

## Enforcement
- MUST be in agent pre-flight checks
- MUST be in CLAUDE.md global instructions
- Setup.sh MUST ensure agents include this

## Validation
```bash
# Create test marker
echo "COMPACTION_TIME:2025-01-01 12:00:00" > /tmp/compaction_marker.txt
echo "WORKING_DIR:/workspaces/test" >> /tmp/compaction_marker.txt

# Agent should detect and report
# Then verify cleanup
ls -la /tmp/compaction_marker.txt 2>&1 | grep "No such file"
```

## Common Mistakes
- ❌ Not checking for marker
- ❌ Not cleaning up marker after reading
- ❌ Not loading preserved state
- ✅ Check, read, recover, cleanup