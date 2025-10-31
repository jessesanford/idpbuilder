# Transcript-Based Compaction Recovery System

## Overview
Successfully implemented and tested a precise agent identification system that uses transcript paths as unique identifiers to ensure only the compacted agent needs recovery while other agents continue working normally.

## Problem Solved
Previously, the compaction recovery system used a generic marker (`/tmp/compaction_marker.txt`) which could block ALL agents when ANY agent was compacted. This created unnecessary disruption in multi-agent workflows.

## Solution: Transcript-Based Identification

### Key Innovation
Each Claude Code agent session has a unique transcript path containing a UUID:
```
/home/vscode/.claude/projects/{context}/{UUID}.jsonl
```

This UUID serves as a unique identifier for each agent session, enabling precise tracking of which specific agent was compacted.

### Implementation

#### 1. PreCompact Hook (`transcript-based-precompact.sh`)
- Extracts transcript ID from `CLAUDE_TRANSCRIPT_PATH`
- Creates agent-specific marker: `/tmp/compaction_marker_{TRANSCRIPT_ID}.txt`
- Saves agent-specific TODOs: `/tmp/todos-precompact-{TRANSCRIPT_ID}.txt`
- Records all recovery information in the marker file

#### 2. PreToolUse Hook (`transcript-based-pretooluse.sh`)
- Extracts transcript ID from current session
- Checks ONLY for matching marker file
- Blocks ONLY if this specific transcript was compacted
- Shows other compacted agents as informational only
- Archives processed markers for debugging

#### 3. Test Suite (`test-transcript-identification.sh`, `test-complete-transcript-system.sh`)
- Proves transcript paths uniquely identify agents
- Verifies precise blocking (only compacted agent)
- Confirms non-compacted agents continue normally
- Tests multi-agent independence

## Benefits

### Precision
- **Before**: Generic marker blocked ALL agents
- **After**: Each agent has unique marker, only affected agent blocked

### Multi-Agent Safety
- Multiple agents can run simultaneously
- Compaction of one agent doesn't affect others
- Each agent recovers independently

### Debugging
- Markers archived with timestamps for troubleshooting
- Clear audit trail of compaction events
- Agent type and context preserved

## Test Results
All tests passed successfully:
✓ Transcript IDs uniquely identify agents
✓ Each agent has its own compaction marker
✓ Only compacted agents are blocked for recovery
✓ Non-compacted agents continue working normally
✓ Multiple agents can run simultaneously
✓ Recovery is precise and agent-specific

## Files Created/Modified

### New Files
1. `utilities/transcript-based-precompact.sh` - Improved PreCompact hook
2. `utilities/transcript-based-pretooluse.sh` - Improved PreToolUse hook
3. `utilities/test-transcript-identification.sh` - Unit tests for theory
4. `utilities/test-complete-transcript-system.sh` - Integration tests

### Key Features
- Automatic transcript ID extraction
- Agent-specific marker files
- Graceful fallback to generic marker if no transcript
- Marker archival after recovery
- Informational notices about other compacted agents

## Usage

### For Hook Configuration
Replace existing hooks with transcript-based versions:
```bash
# PreCompact hook
cp utilities/transcript-based-precompact.sh ~/.claude/hooks/pre-compact.sh

# PreToolUse hook  
cp utilities/transcript-based-pretooluse.sh ~/.claude/hooks/pre-tool-use.sh
```

### Marker File Format
```
/tmp/compaction_marker_{TRANSCRIPT_ID}.txt
```
Where `TRANSCRIPT_ID` is the UUID from the agent's transcript path.

## Backward Compatibility
- Falls back to generic marker if `CLAUDE_TRANSCRIPT_PATH` unavailable
- Continues to support TODO recovery workflows
- Compatible with existing state machine and recovery protocols

## Conclusion
The transcript-based identification system provides a robust, precise solution for multi-agent compaction recovery, ensuring minimal disruption to Software Factory 2.0 workflows while maintaining complete recovery capabilities for affected agents.