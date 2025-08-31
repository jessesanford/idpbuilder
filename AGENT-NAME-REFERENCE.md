# Software Factory 2.0 Agent Name Reference

## Official Agent Names (from YAML frontmatter)

These are the canonical agent names defined in the `.claude/agents/*.md` files:

| File | YAML Name | Agent ID | Description |
|------|-----------|----------|-------------|
| `orchestrator.md` | `orchestrator` | `orchestrator` | Orchestrates AI agents, manages state |
| `sw-engineer.md` | `sw-engineer` | `sw-engineer` | Implements code following plans |
| `code-reviewer.md` | `code-reviewer` | `code-reviewer` | Reviews code for quality and patterns |
| `architect.md` | `architect` | `architect` | Reviews architectural decisions |

## Correct Usage Examples

### When Spawning Agents

```bash
# CORRECT - Use the official agent names
Task: orchestrator
Task: sw-engineer
Task: code-reviewer
Task: architect

# WRONG - Old or incorrect names
Task: @agent-orchestrator-prompt-engineer-task-master  # ❌
Task: @agent-software-engineer  # ❌
Task: @agent-kcp-go-lang-sr-sw-eng  # ❌
Task: @agent-architect-reviewer  # ❌
```

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