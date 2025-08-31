# The Truth About "Hooks" in Software Factory 2.0

## TL;DR
- **Real Hook**: Only PreCompact (in `.claude/settings.json`)
- **Manual Scripts**: Everything in `utilities/` directory
- **No Magic**: Most things require manual execution

## What Actually Exists

### ✅ REAL Claude Code Hook (Automatic)

**PreCompact** - The ONLY hook we use:
```json
// .claude/settings.json
{
  "hooks": {
    "PreCompact": [
      {
        "matcher": "auto|manual",
        "hooks": [{
          "type": "command",
          "command": "bash command that creates /tmp/compaction_marker.txt"
        }]
      }
    ]
  }
}
```

This automatically:
- Runs before context compaction
- Creates `/tmp/compaction_marker.txt` with context info
- Preserves latest TODO file to `/tmp/todos-precompact.txt`

### ❌ NOT Hooks (Manual Scripts)

Everything in `utilities/` directory:
- `pre-compact.sh` - Run manually for comprehensive preservation
- `post-compact.sh` - Run manually to check for compaction
- `todo-preservation.sh` - Run manually to save/load TODOs
- `state-snapshot.sh` - Run manually to create snapshots
- `recovery-assistant.sh` - Run manually for recovery help

**These do NOT run automatically!**

## Common Misconceptions

### Myth 1: "Hooks run automatically"
**Reality**: Only PreCompact runs automatically. Everything else is manual.

### Myth 2: "There are many hooks"
**Reality**: Claude Code supports several hooks, but SF 2.0 only uses PreCompact.

### Myth 3: "State is preserved automatically"
**Reality**: Only minimal marker file is automatic. Comprehensive preservation requires manual script execution.

### Myth 4: "Recovery is automatic"
**Reality**: You must manually run recovery scripts and reload context.

## How to Use SF 2.0 Correctly

### 1. Understand What's Automatic
- PreCompact hook creates marker file - that's it!

### 2. Run Scripts Manually
```bash
# Before expected compaction
./utilities/pre-compact.sh

# During state transitions
./utilities/todo-preservation.sh save orchestrator WAVE_COMPLETE

# After context loss
./utilities/post-compact.sh
./utilities/recovery-assistant.sh
```

### 3. Check for Marker
```bash
# After resuming work
if [ -f /tmp/compaction_marker.txt ]; then
    echo "Compaction detected!"
    cat /tmp/compaction_marker.txt
fi
```

### 4. Create Workflow Habits
- Save TODOs at transitions
- Create snapshots at milestones
- Run utilities proactively

## Why This Design?

1. **Claude Code Limitations**: Limited hook support
2. **Flexibility**: Manual scripts can be run anytime
3. **Simplicity**: One real hook is easier to maintain
4. **Reliability**: Manual execution is predictable

## Installation Verification

After running `setup.sh`, verify:

1. **Check hooks configuration**:
```bash
cat .claude/settings.json
# Should show ONLY PreCompact hook
```

2. **Check utilities**:
```bash
ls -la utilities/
# Should show executable .sh scripts
```

3. **Test PreCompact**:
```bash
# Run /compact command in Claude Code
# Then check:
ls -la /tmp/compaction_marker.txt
```

4. **Test utilities**:
```bash
./utilities/pre-compact.sh
# Should run without errors
```

## The Bottom Line

SF 2.0 provides:
- **One automatic hook** (PreCompact)
- **Helpful manual scripts** (utilities)
- **Structured methodology** (phases/waves/efforts)
- **Clear templates** (planning documents)

It does NOT provide:
- Automatic state management
- Multiple automatic hooks
- Magic recovery
- Hands-free operation

Success requires understanding this reality and working with it, not against it.