# Quick Start Guide for idpbuilder

## Setup Complete! 🎉

Your software factory has been configured with:
- **Language**: go
- **Complexity**: Level 2
- **Location**: /home/vscode/workspaces/idpbuilder

## Next Steps

### 1. Review and Customize Your Plan
Edit `orchestrator/PROJECT-IMPLEMENTATION-PLAN.md` to define your specific:
- Phases and goals
- Waves and groupings
- Individual efforts (keep each under 800 lines!)

### 2. Configure Your Agents
The following agents are ready:
- Orchestrator: `.claude/agents/orchestrator-task-master.md`
- Code Reviewer: `.claude/agents/code-reviewer.md`
- SW Engineer: `.claude/agents/sw-engineer-go.md`
- Architect: `.claude/agents/architect-reviewer.md`

### 3. Start Orchestration
```bash
cd /home/vscode/workspaces/idpbuilder
# In Claude Code, run:
/continue-orchestrating
```

### 4. Monitor Progress
- State tracking: `orchestrator/orchestrator-state.yaml`
- TODOs: `todos/` directory
- Work logs: Each effort directory

## Key Commands

### Measure Lines (Always Use This!)
```bash
/home/vscode/workspaces/idpbuilder/tools/line-counter.sh -c branch-name
```

### Check State
```bash
cat orchestrator/orchestrator-state.yaml
```

## Critical Configuration

### 🔴 settings.json is ESSENTIAL
The  file enables:
- Compaction recovery
- TODO state preservation  
- Context maintenance

**Without it, you WILL lose work during compaction!**

## Important Rules

1. **NEVER exceed 800 lines per effort** - Split if needed
2. **ALWAYS review before proceeding** - No skipping reviews
3. **SEQUENTIAL splits only** - Never parallel
4. **Orchestrator NEVER writes code** - Only coordinates

## Troubleshooting

If context is lost:
1. Check `orchestrator/orchestrator-state.yaml`
2. Look for TODO files in `todos/`
3. Read `.claude/CLAUDE.md` sections 7-9

## Support

Refer to:
- Main documentation: `README.md`
- State machine: `core/SOFTWARE-FACTORY-STATE-MACHINE.md`
- Operations guide: `core/ORCHESTRATOR-MASTER-OPERATIONS-GUIDE.md`

Happy coding with your new software factory! 🏭
