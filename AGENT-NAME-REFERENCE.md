# Software Factory 2.0 Agent Name Reference

## Official Agent Names (from YAML frontmatter)

These are the canonical agent names defined in the `.claude/agents/*.md` files:

| File | YAML Name | Agent ID for Spawning | Description |
|------|-----------|----------------------|-------------|
| `orchestrator.md` | `orchestrator` | `orchestrator` | Orchestrates AI agents, manages state |
| `sw-engineer.md` | `sw-engineer` | `software-engineer` | Implements code following plans |
| `code-reviewer.md` | `code-reviewer` | `code-reviewer` | Reviews code for quality and patterns |
| `architect.md` | `architect` | `architect` | Reviews architectural decisions |

## Correct Usage Examples

### When Spawning Agents

```bash
# CORRECT - Use the official spawning agent names
Task: orchestrator
Task: software-engineer  # NOTE: For spawning, use "software-engineer" not "sw-engineer"
Task: code-reviewer
Task: architect

# WRONG - Old or incorrect names
Task: sw-engineer  # ❌ Will fail with "Agent type 'sw-engineer' not found"
Task: @agent-orchestrator-prompt-engineer-task-master  # ❌
Task: @agent-software-engineer  # ❌
Task: @agent-kcp-go-lang-sr-sw-eng  # ❌
Task: @agent-architect-reviewer  # ❌
```

## Important Distinction: File Name vs Agent Type

**CRITICAL**: The sw-engineer agent has a mismatch between its file name and spawning type:
- **File name**: `.claude/agents/sw-engineer.md`
- **YAML name in file**: `sw-engineer`
- **Agent type for spawning**: `software-engineer` (with full word "software")

This is why spawning with `Task: sw-engineer` fails with "Agent type 'sw-engineer' not found".
Always use `Task: software-engineer` when spawning the software engineer agent.

### In Agent Acknowledgments

```bash
# Software Engineer acknowledgment
I am sw-engineer in state {CURRENT_STATE}

# NOT
I am @agent-software-engineer in state {CURRENT_STATE}  # ❌
I am @agent-sw-engineer in state {CURRENT_STATE}  # ❌
```

### In Pre-Flight Checks

```bash
echo "AGENT: sw-engineer"  # ✅
echo "AGENT: @agent-software-engineer"  # ❌
echo "AGENT: @agent-sw-engineer"  # ❌
```

## Directory Names vs Agent Names

Note that directory names may differ from agent IDs:
- Directory: `sw-engineer/` (for file organization)
- Agent ID: `@agent-software-engineer` (for spawning)

## Migration Notes

The following old names have been updated:
- `@agent-orchestrator-prompt-engineer-task-master` → `@agent-orchestrator`
- `@agent-sw-engineer` → `@agent-software-engineer`
- `@agent-architect-reviewer` → `@agent-architect`
- `@agent-code-reviewer` → `@agent-code-reviewer` (unchanged)

## Files Updated

The following files have been updated to use correct agent names:
1. `.claude/agents/sw-engineer.md` - Updated to use `@agent-software-engineer`
2. `.claude/agents/orchestrator.md` - Spawning examples updated
3. `.claude/commands/continue-orchestrating.md` - All Task references updated
4. `utilities/recovery-assistant.sh` - Agent type listings updated
5. `utilities/todo-preservation.sh` - Agent type listings updated

## Important Note

When referencing agents in any context (spawning, documentation, scripts), always use the official agent IDs listed above. These match the `name` field in the YAML frontmatter of each agent configuration file.