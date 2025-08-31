# 📍 Utility Location Strategy

## Overview

Software Factory 2.0 utilities are installed in **TWO locations** to ensure they're always accessible regardless of which directory agents are working in.

## Installation Locations

### 1. Global Location: `~/.claude/utilities/`
- **Purpose**: Consistent access point for all agents
- **Used by**: Automatic agent pre-flight checks
- **Path**: `$HOME/.claude/utilities/` or `/home/user/.claude/utilities/`
- **Why**: Agents working in different directories can always find utilities here

### 2. Project Location: `./utilities/`
- **Purpose**: Manual execution and project-specific reference
- **Used by**: Users running utilities manually
- **Path**: `{PROJECT_DIR}/utilities/`
- **Why**: Keeps utilities with the project for version control

## How Agents Find Utilities

Agents check for utilities in this order:
1. `$HOME/.claude/utilities/` (primary location)
2. `/home/user/.claude/utilities/` (fallback for different user configs)
3. `./utilities/` (last resort if in project root)

### Example from Agent Pre-Flight Check:
```bash
# CHECK 0: AUTOMATIC COMPACTION DETECTION
if [ -f /tmp/compaction_marker.txt ]; then
    # Run post-compact utility (check standard locations)
    if [ -f "$HOME/.claude/utilities/post-compact.sh" ]; then
        $HOME/.claude/utilities/post-compact.sh
    elif [ -f "/home/user/.claude/utilities/post-compact.sh" ]; then
        /home/user/.claude/utilities/post-compact.sh
    fi
fi
```

## Why This Strategy?

### Problem with Relative Paths
- **Orchestrator**: Works in project root (`/workspaces/my-project/`)
- **SW Engineer**: Works in effort directories (`/workspaces/my-project/efforts/phase1/wave1/effort1/`)
- **Code Reviewer**: Also in effort directories
- **Architect**: Could be in project root or integration branches

Relative paths like `./utilities/` or `../../utilities/` would fail depending on where the agent is working.

### Solution: Standard Global Location
- All agents know to check `~/.claude/utilities/`
- Works regardless of current working directory
- Consistent across all projects
- Survives project deletions/moves

## Installation Process

When you run `setup.sh` or `setup-noninteractive.sh`:

1. **Copies utilities to project**:
   ```bash
   cp -r utilities/ $TARGET_DIR/utilities/
   ```

2. **Installs globally for agents**:
   ```bash
   mkdir -p $HOME/.claude/utilities
   cp utilities/*.sh $HOME/.claude/utilities/
   chmod +x $HOME/.claude/utilities/*.sh
   ```

## Manual Usage

Users can run utilities from either location:

### From Project Directory:
```bash
cd /workspaces/my-project
./utilities/post-compact.sh
./utilities/todo-preservation.sh save orchestrator
```

### From Global Location:
```bash
# From anywhere
~/.claude/utilities/post-compact.sh
~/.claude/utilities/recovery-assistant.sh
```

## Utility Scripts Available

| Script | Purpose | Auto-Run by Agents |
|--------|---------|-------------------|
| `post-compact.sh` | Check compaction status | ✅ Yes |
| `todo-preservation.sh` | Save/load TODOs | ✅ Yes (load mode) |
| `recovery-assistant.sh` | Interactive recovery help | ❌ No (referenced) |
| `pre-compact.sh` | Manual state preservation | ❌ No |
| `state-snapshot.sh` | Create state snapshots | ❌ No |

## Directory Structure After Setup

```
$HOME/
├── .claude/
│   └── utilities/              # Global utilities (agents use this)
│       ├── post-compact.sh
│       ├── todo-preservation.sh
│       ├── recovery-assistant.sh
│       ├── pre-compact.sh
│       └── state-snapshot.sh

/workspaces/my-project/
├── utilities/                  # Project utilities (manual use)
│   ├── post-compact.sh
│   ├── todo-preservation.sh
│   ├── recovery-assistant.sh
│   ├── pre-compact.sh
│   └── state-snapshot.sh
├── efforts/
│   └── phase1/
│       └── wave1/
│           └── effort1/       # SW Engineer works here
│               ├── pkg/
│               └── ...
```

## Updating Utilities

If you update utility scripts:

1. **Update in template**: Edit files in software-factory-2.0-template/utilities/
2. **Reinstall globally**: 
   ```bash
   cp utilities/*.sh ~/.claude/utilities/
   chmod +x ~/.claude/utilities/*.sh
   ```
3. **Update in projects**: Copy to each project's utilities/ directory

## Troubleshooting

### Utilities Not Found
If agents can't find utilities:
1. Check if `~/.claude/utilities/` exists
2. Verify scripts are executable: `ls -la ~/.claude/utilities/`
3. Re-run setup script to reinstall

### Permission Denied
```bash
chmod +x ~/.claude/utilities/*.sh
```

### Different User Directory
Some systems use different home paths. Agents check both:
- `$HOME/.claude/utilities/`
- `/home/user/.claude/utilities/`

### Manual Installation
If setup script fails to install globally:
```bash
mkdir -p ~/.claude/utilities
cp /path/to/template/utilities/*.sh ~/.claude/utilities/
chmod +x ~/.claude/utilities/*.sh
```

## Best Practices

1. **Always use setup script**: It handles both locations automatically
2. **Don't move utilities**: Keep them in standard locations
3. **Update both locations**: When modifying utilities
4. **Check installation**: Verify `~/.claude/utilities/` exists after setup
5. **Version control**: Keep project utilities in git, global ones are local

## Summary

- **Agents automatically use**: `~/.claude/utilities/`
- **Users manually run from**: `./utilities/` or `~/.claude/utilities/`
- **Setup installs to**: Both locations
- **Why two locations**: Ensures utilities always work regardless of CWD
- **Key benefit**: Agents in any directory can recover from compaction