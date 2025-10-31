# R173.0.0 - State Preservation Protocol

## Rule Statement
State preservation MUST use a combination of automatic PreCompact hook (minimal) and manual utility scripts (comprehensive).

## Rationale
Automatic preservation is limited to PreCompact hook. Comprehensive state preservation requires manual intervention.

## Implementation

### Automatic Preservation (PreCompact Hook)
```bash
# Creates minimal marker file
/tmp/compaction_marker.txt containing:
- COMPACTION_TIME
- WORKING_DIR
- GIT_BRANCH
- ACTIVE_FILES
- TODO_STATE_SAVED (if todos exist)
```

### Manual Preservation (Utility Scripts)
```bash
# Comprehensive state preservation
./utilities/pre-compact.sh
# Creates:
# - checkpoints/active/checkpoint-{timestamp}.md
# - Detailed state information
# - TODO preservation
# - Git state capture
```

### Recovery Detection
```bash
# In agent startup or CLAUDE.md
if [ -f /tmp/compaction_marker.txt ]; then
    echo "Compaction detected!"
    cat /tmp/compaction_marker.txt
    # Initiate recovery procedures
fi
```

## Enforcement
- Agents MUST check for marker on startup
- Manual preservation SHOULD be run at milestones
- TODOs MUST be saved during state transitions

## Validation
```bash
# Test automatic preservation
# Run /compact in Claude Code, then:
cat /tmp/compaction_marker.txt

# Test manual preservation
./utilities/pre-compact.sh
ls -la checkpoints/active/
```

## Common Mistakes
- ❌ Relying only on automatic preservation
- ❌ Not checking for marker on resume
- ❌ Forgetting to run manual utilities
- ✅ Using both automatic and manual preservation