# TODOs Directory

## Purpose
This directory stores TODO state files for agents to persist their task lists across context switches, compaction events, and state machine transitions.

## File Naming Convention
```
{agent-name}-{STATE}-{YYYYMMDD-HHMMSS}.todo
```

### Examples:
- `orchestrator-WAVE_COMPLETE-20250121-143000.todo`
- `sw-eng-IMPLEMENTATION-20250121-145500.todo`
- `code-reviewer-CREATE_SPLIT_PLAN-20250121-150000.todo`
- `architect-WAVE_REVIEW-20250121-153000.todo`

## When Files Are Created

### Automatic Triggers
1. **State Transitions**: When moving between state machine states
2. **Before Spawning Agents**: Orchestrator saves pending tasks
3. **Context Compaction**: System saves before memory compression
4. **Major Milestones**: Completing waves, phases, or efforts

### Manual Triggers
1. After significant TODO changes
2. Every 10-15 messages in long conversations
3. Before any risky operations
4. When switching between major task types

## File Format Example
```markdown
# ORCHESTRATOR TODOs at WAVE_COMPLETE
## In Progress:
- Create integration branch for wave 1
- Run architect review

## Pending:
- Begin wave 2 planning
- Update state file with review results

## Context:
- Phase 1, Wave 1 complete
- 6 efforts successfully implemented
- Awaiting architectural approval
```

## Recovery Process

### After Compaction
1. Check `/tmp/compaction_marker.txt` for notification
2. Find latest TODO file for your agent type
3. READ file with Read tool
4. LOAD into TodoWrite tool (not just read!)
5. Verify all TODOs loaded correctly

### After Context Loss
1. Identify which agent you are
2. Find most recent TODO file
3. Load and merge with any in-memory TODOs
4. Continue from recovered state

## Maintenance

### Cleanup Rules
- Keep last 5 files per agent
- Delete files older than 24 hours
- Commit important state files to git

### Git Integration
```bash
# Always commit TODO files after creating
cd /home/vscode/workspaces/idpbuilder
git add todos/*.todo
git commit -m "todo: save {agent} state at {STATE}"
git push
```

## Related Documentation
- See `.claude/CLAUDE.md` Section 8 for TODO state management
- See `.claude/CLAUDE.md` Section 9 for pre-compaction saving
- See `SOFTWARE-FACTORY-STATE-MACHINE.md` for state definitions

## Critical Rules
1. **Always use TodoWrite tool** to load TODOs (not just read)
2. **Never delete active TODO files** during recovery
3. **Commit important states** to preserve in git
4. **Each agent manages only their own** TODO files
5. **Deduplicate when loading** to avoid duplicate tasks