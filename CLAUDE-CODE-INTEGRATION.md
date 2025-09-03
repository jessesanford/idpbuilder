# Claude Code Integration - How SF 2.0 Actually Works

## The Reality Check

Let's be clear about what Claude Code can and cannot do:

### ✅ What Claude Code DOES Support

1. **PreCompact Hook** - Automatically runs before context compaction
   - Configured in `.claude/settings.json`
   - Executes inline bash commands
   - We use this to create `/tmp/compaction_marker.txt`

2. **Other Available Hooks** (that we DON'T use):
   - PreToolUse/PostToolUse - Before/after tool execution
   - UserPromptSubmit - When user submits prompt
   - Stop/SubagentStop - When agents finish
   - SessionStart/SessionEnd - Session lifecycle
   - Notification - For permissions/idle prompts

### ❌ What Are NOT Real Hooks

All the scripts in `utilities/` directory:
- `pre-compact.sh` - Manual state preservation
- `post-compact.sh` - Manual recovery helper
- `todo-preservation.sh` - Manual TODO management
- `state-snapshot.sh` - Manual checkpoint creation
- `recovery-assistant.sh` - Manual recovery wizard

**These must be run MANUALLY - they are NOT hooks!**

## How SF 2.0 Actually Integrates

### 1. Automatic Integration (Minimal)

Only ONE thing happens automatically:
```json
// .claude/settings.json
{
  "hooks": {
    "PreCompact": [{
      "matcher": "auto|manual",
      "hooks": [{
        "type": "command",
        "command": "creates /tmp/compaction_marker.txt with context info"
      }]
    }]
  }
}
```

This creates a marker file when compaction occurs.

### 2. Manual Workflow (Everything Else)

Everything else requires manual execution:

```bash
# You must manually run utilities when needed:
./utilities/pre-compact.sh          # Before expected compaction
./utilities/todo-preservation.sh    # During state transitions
./utilities/state-snapshot.sh       # At milestones
./utilities/post-compact.sh         # After context loss
./utilities/recovery-assistant.sh   # When confused
```

### 3. Agent Configuration

Agents are configured via:
- `.claude/agents/` directory - Agent personality files
- Agent configs include pre-flight checks (but agents must execute them)
- No automatic enforcement - relies on agent compliance

## Directory Structure After Setup

```
your-project/
├── .claude/
│   ├── settings.json        # Contains ONLY PreCompact hook
│   └── agents/              # Agent configurations
│       ├── orchestrator.md
│       ├── sw-engineer.md
│       ├── code-reviewer.md
│       └── architect-reviewer.md
├── utilities/               # MANUAL helper scripts
│   ├── pre-compact.sh
│   ├── post-compact.sh
│   ├── todo-preservation.sh
│   ├── state-snapshot.sh
│   └── recovery-assistant.sh
├── todos/                   # TODO state files
├── checkpoints/            # State snapshots
└── snapshots/              # Recovery points
```

## What Setup.sh Does

1. **Copies agent configurations** to `.claude/agents/`
2. **Configures PreCompact hook** in `.claude/settings.json`
3. **Copies utility scripts** to `utilities/`
4. **Makes utilities executable**
5. **Creates support directories**
6. **Updates paths** to be project-specific

## What You Must Do Manually

1. **Run utilities** at appropriate times
2. **Check for marker** (`/tmp/compaction_marker.txt`) after resumption
3. **Execute pre-flight checks** (agents should do this, but verify)
4. **Save TODO states** during transitions
5. **Create snapshots** at milestones

## Best Practices for Real Usage

### 1. Create Aliases
```bash
# Add to your .bashrc or project setup
alias save-state='./utilities/pre-compact.sh'
alias save-todos='./utilities/todo-preservation.sh save'
alias recover='./utilities/recovery-assistant.sh'
```

### 2. Include in Workflow
Tell agents to run utilities:
```markdown
Before any major state transition:
1. Run: ./utilities/todo-preservation.sh save orchestrator WAVE_COMPLETE
2. Proceed with transition
```

### 3. Regular Maintenance
```bash
# Clean old TODO files periodically
find todos/ -name "*.todo" -mtime +7 -delete

# Check marker on resume
cat /tmp/compaction_marker.txt 2>/dev/null
```

### 4. Understand Limitations
- No automatic state preservation (except minimal PreCompact)
- No automatic recovery
- No automatic utility execution
- Manual intervention required for most features

## The Truth About "Automation"

SF 2.0 provides:
- **Structure** - Organized approach to complex projects
- **Templates** - Consistent planning documents
- **Utilities** - Helpful manual tools
- **Minimal Automation** - Only PreCompact hook

It does NOT provide:
- Automatic state management (beyond marker file)
- Automatic recovery
- Automatic utility execution
- Enforcement of rules (relies on agent compliance)

## Conclusion

SF 2.0 is primarily a **methodology and toolkit**, not an automated system. Success depends on:
1. Understanding what's manual vs automatic
2. Running utilities proactively
3. Following the methodology consistently
4. Not expecting magic automation

The power comes from the structured approach and helpful utilities, not from automation that doesn't exist.