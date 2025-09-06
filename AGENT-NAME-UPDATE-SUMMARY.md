# Agent Name Update Summary

## What Was Fixed

The agent names in acknowledgments and spawning commands were inconsistent with the YAML frontmatter names. This has been corrected throughout the Software Factory 2.0 template.

## Correct Agent Names (from YAML frontmatter)

| Agent | YAML Name | Correct Agent ID | Old/Wrong ID |
|-------|-----------|------------------|--------------|
| Orchestrator | `orchestrator` | `@agent-orchestrator` | ~~@agent-orchestrator-prompt-engineer-task-master~~ |
| Software Engineer | `software-engineer` | `@agent-software-engineer` | ~~@agent-sw-engineer~~ |
| Code Reviewer | `code-reviewer` | `@agent-code-reviewer` | (unchanged) |
| Architect | `architect` | `@agent-architect` | ~~@agent-architect-reviewer~~ |

## Files Updated

### Agent Configuration Files
1. **`.claude/agents/sw-engineer.md`**
   - Pre-flight check: `echo "AGENT: @agent-software-engineer"`
   - Identity section: `@agent-software-engineer`
   - Acknowledgment: `I am @agent-software-engineer`

2. **`.claude/agents/orchestrator.md`**
   - Spawning examples updated to use `@agent-software-engineer`

3. **`.claude/agents/architect.md`**
   - Already using correct `@agent-architect`

4. **`.claude/agents/code-reviewer.md`**
   - Already using correct `@agent-code-reviewer`

### Command Files
5. **`.claude/commands/continue-orchestrating.md`**
   - All `Task: @agent-sw-engineer` → `Task: @agent-software-engineer`

### Utility Scripts
6. **`utilities/recovery-assistant.sh`**
   - Agent type listings updated with correct IDs
   - Prompt choices updated

7. **`utilities/todo-preservation.sh`**
   - Agent type listings updated with correct IDs

### Quick Reference Files
8. **`quick-reference/orchestrator-quick-ref.md`**
   - Spawning example updated

9. **`quick-reference/orchestrator-workspace-setup-quick-ref.md`**
   - All spawning examples updated

## Why This Matters

1. **Consistency**: Agent IDs now match their YAML configuration names
2. **Claude Code Recognition**: Claude Code uses the YAML `name` field to identify agents
3. **Spawning Success**: Using correct agent IDs ensures proper agent spawning
4. **Clear Documentation**: No confusion between different naming conventions

## Going Forward

Always use these agent IDs when:
- Spawning agents: `Task: @agent-software-engineer`
- In acknowledgments: `I am @agent-software-engineer`
- In documentation: Refer to `@agent-software-engineer`
- In scripts: Use the correct agent IDs

## Note on Directory Names

Directory names (like `sw-engineer/`) may still use abbreviated forms for file organization. This is fine - what matters is that agent IDs used for spawning and identification match the YAML frontmatter names.