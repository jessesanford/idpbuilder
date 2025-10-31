# Software Factory 2.0 - Utility Scripts & Automatic Recovery

This directory contains **utility scripts** for state management. While these scripts themselves are manual helpers, **COMPACTION DETECTION IS AUTOMATIC** through agent pre-flight checks.

## 🔄 AUTOMATIC COMPACTION DETECTION SYSTEM

### How Automatic Recovery Works:
1. **PreCompact Hook** (in `.claude/settings.json`) creates `/tmp/compaction_marker.txt` automatically
2. **EVERY AGENT** checks for this marker as **CHECK 0** in pre-flight checks
3. **Detection triggers automatic recovery** steps including TODO restoration
4. **No manual intervention required** - agents detect and recover automatically

### What ARE Real Claude Code Hooks:
According to [Claude Code documentation](https://docs.anthropic.com/en/docs/claude-code/hooks), the ONLY hooks we use are:
- **PreCompact** - Runs automatically before context compaction
  - Configured in `.claude/settings.json` as inline commands
  - Creates `/tmp/compaction_marker.txt` for recovery detection
  - **Agents automatically detect this marker** via R186 compliance

### What These Scripts Actually Are:
**MANUAL UTILITY SCRIPTS** for additional state management:
- The scripts themselves do NOT run automatically
- They are NOT triggered by Claude Code events
- They must be executed manually for extra state preservation
- **BUT: Compaction detection IS automatic via agent pre-flight checks (R186)**

## Real Hook Configuration

The ONLY actual hook is in `.claude/settings.json`:
```json
{
  "hooks": {
    "PreCompact": [
      {
        "matcher": "auto",  // or "manual"
        "hooks": [{
          "type": "command",
          "command": "echo 'marker' > /tmp/compaction_marker.txt ..."
        }]
      }
    ]
  }
}
```

## Utility Scripts in This Directory

### Manual State Preservation Tools

#### pre-compact.sh
```bash
# Run manually BEFORE you expect compaction
./utilities/pre-compact.sh
```
- Comprehensive state preservation
- Creates detailed checkpoints
- Saves TODO states, git info, active files

#### post-compact.sh  
```bash
# Run manually AFTER resuming work
./utilities/post-compact.sh
```
- Checks for compaction marker
- Guides recovery process
- Helps reload context

#### todo-preservation.sh
```bash
# Run manually during state transitions
./utilities/todo-preservation.sh save orchestrator WAVE_COMPLETE
./utilities/todo-preservation.sh load orchestrator
```
- Save/load TODO states
- Critical for agent continuity
- Must be called explicitly

#### state-snapshot.sh
```bash
# Run manually at milestones
./utilities/state-snapshot.sh
```
- Creates recovery checkpoints
- Saves comprehensive state
- Good for major transitions

#### recovery-assistant.sh
```bash
# Run manually when confused
./utilities/recovery-assistant.sh
```
- Interactive recovery wizard
- Helps determine current state
- Guides through recovery steps

## How Recovery Works Automatically

### 🎯 AUTOMATIC DETECTION (No Action Required):
When compaction occurs:
1. PreCompact hook creates marker automatically
2. **Agents detect marker in pre-flight checks (R186)**
3. Agents guide you through recovery automatically
4. No need to manually run recovery scripts

### Manual Utilities (Optional Enhancement):

#### 1. During Normal Work
```bash
# Periodically save extra state
./utilities/state-snapshot.sh

# Before major transitions
./utilities/todo-preservation.sh save [agent] [state]
```

#### 2. Before Expected Compaction
```bash
# Run comprehensive preservation (optional)
./utilities/pre-compact.sh

# The automatic PreCompact hook handles basics
# Agents will detect and recover automatically
```

#### 3. Manual Recovery Assistance
```bash
# If automatic recovery needs help
./utilities/recovery-assistant.sh

# Manually reload specific TODOs
./utilities/todo-preservation.sh load [agent]
```

## Why This Confusion Exists

The legacy SF 1.0 and early documentation referred to these as "hooks" but they're really just helper scripts. Claude Code's actual hook system is limited to:
- PreCompact (we use this)
- PreToolUse/PostToolUse (we don't use)
- UserPromptSubmit (we don't use)
- Stop/SubagentStop (we don't use)
- SessionStart/SessionEnd (we don't use)
- Notification (we don't use)

## Installation

When setup.sh runs, it:
1. Copies these utilities to `{PROJECT}/utilities/`
2. Makes them executable
3. Creates support directories (`todos/`, `checkpoints/`, `snapshots/`)
4. Configures the REAL PreCompact hook in `.claude/settings.json`

## Best Practices

1. **Don't expect automatic execution** - These are manual tools
2. **Run proactively** - Don't wait for problems
3. **Create aliases** for frequently used commands:
   ```bash
   alias save-todos='./utilities/todo-preservation.sh save'
   alias snapshot='./utilities/state-snapshot.sh'
   ```
4. **Check for marker** - Look for `/tmp/compaction_marker.txt` when resuming
5. **Use scripts liberally** - Better to over-preserve than lose state

## The Bottom Line

- **Real Hook**: PreCompact in settings.json (automatic)
- **Compaction Detection**: AUTOMATIC via agent pre-flight checks (R186)
- **Recovery Process**: AUTOMATIC when agents detect marker
- **These Scripts**: Optional manual utilities for enhanced state management
- **Best Approach**: Let automatic detection handle recovery, use scripts for extra safety

**🎯 KEY POINT**: Compaction detection and basic recovery ARE automatic through agent pre-flight checks. You don't need to manually check for compaction - agents do it automatically as their FIRST action per R186!