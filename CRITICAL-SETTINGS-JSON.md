# 🔴 CRITICAL: settings.json Configuration

## THIS FILE IS ESSENTIAL FOR THE SOFTWARE FACTORY TO FUNCTION PROPERLY

The `.claude/settings.json` file contains **critical hooks** that enable:
1. **Compaction recovery** - Preserves context when memory limits are reached
2. **TODO state management** - Saves TODOs before context is compressed
3. **State preservation** - Maintains working directory and branch information

## ⚠️ WITHOUT THIS FILE, YOU WILL LOSE WORK DURING COMPACTION

### What is Compaction?

Claude Code has memory limits. When these are reached, it "compacts" the conversation, potentially losing:
- Current TODO lists
- Working directory context
- Active branch information
- In-progress work state

### How settings.json Prevents Data Loss

The `PreCompact` hooks in settings.json run BEFORE compaction occurs, creating:
1. A marker file (`/tmp/compaction_marker.txt`) with context
2. A backup of the latest TODO file
3. Working directory and branch information

This allows agents to recover their state after compaction.

## Required Configuration

### File Location
```
/workspaces/[project]/.claude/settings.json
```

### File Contents
```json
{
  "model": "opus",
  "hooks": {
    "PreCompact": [
      {
        "matcher": "auto",
        "hooks": [
          {
            "type": "command",
            "command": "[PreCompact command - see full file]"
          }
        ]
      },
      {
        "matcher": "manual",
        "hooks": [
          {
            "type": "command",
            "command": "[PreCompact command - see full file]"
          }
        ]
      }
    ]
  }
}
```

### What the Hooks Do

1. **Detect Compaction Type** - Auto or Manual
2. **Create Marker File** - `/tmp/compaction_marker.txt` with:
   - Timestamp
   - Working directory
   - Current Git branch
   - Active files
3. **Save TODO State** - Copies latest TODO file to `/tmp/todos-precompact.txt`
4. **Enable Recovery** - CLAUDE.md checks for marker on startup

## Setup Instructions

### For New Projects

When using the setup script:
```bash
./setup.sh
# The script automatically creates settings.json
```

### For Manual Setup

1. **Copy the settings.json file**:
```bash
cp /workspaces/software-factory-template/.claude/settings.json /workspaces/[project]/.claude/
```

2. **Customize the TODO directory path** (if needed):
```bash
# Edit settings.json and change:
TODO_DIR='./todos'
# To your project's TODO directory:
TODO_DIR='/workspaces/[project]/todos'
```

3. **Verify it's working**:
```bash
# Check that hooks are registered
cat .claude/settings.json | grep PreCompact
# Should show PreCompact configuration
```

## How Recovery Works

### 1. Compaction Occurs
- PreCompact hooks run automatically
- Marker and TODO backup created

### 2. Next Agent Response
- CLAUDE.md checks for `/tmp/compaction_marker.txt`
- If found, recovery process starts

### 3. Recovery Process
```bash
# Agent detects marker
if [ -f /tmp/compaction_marker.txt ]; then
    # Read saved context
    cat /tmp/compaction_marker.txt
    
    # Load saved TODOs
    ls -t todos/*.todo | head -1
    # Use TodoWrite to restore TODO list
    
    # Continue from saved state
fi
```

## Integration with Other Systems

### TODO State Management Protocol
The settings.json hooks work with:
- `/protocols/TODO-STATE-MANAGEMENT-PROTOCOL.md`
- TODO files in `/todos/` directory
- TodoWrite tool for state restoration

### CLAUDE.md Compaction Detection
CLAUDE.md Section 1 checks for compaction:
- Detects marker file
- Initiates recovery
- Loads saved TODOs
- Restores context

## Troubleshooting

### Hooks Not Running
```bash
# Verify settings.json exists
ls -la .claude/settings.json

# Check hook configuration
cat .claude/settings.json | python -m json.tool
```

### TODOs Not Saving
```bash
# Verify todos directory exists
ls -la todos/

# Check for TODO files
ls -t todos/*.todo

# Verify path in settings.json matches
grep TODO_DIR .claude/settings.json
```

### Recovery Not Working
```bash
# Check for marker file after compaction
ls -la /tmp/compaction_marker.txt

# Verify CLAUDE.md is checking for it
grep compaction_marker .claude/CLAUDE.md
```

## Critical Points

### ⚠️ DO NOT:
- Delete or rename settings.json
- Modify the PreCompact hook structure
- Change the marker file location (`/tmp/compaction_marker.txt`)
- Disable hooks

### ✅ ALWAYS:
- Include settings.json in your project setup
- Keep todos directory accessible
- Test recovery after setup
- Commit settings.json to version control

## Why This Matters

Without proper compaction handling:
- **Lost TODOs**: In-progress work disappears
- **Lost Context**: Agent forgets what it was doing
- **Wasted Time**: Work must be redone
- **Broken State**: Orchestration flow disrupted

With settings.json properly configured:
- **Seamless Recovery**: Work continues after compaction
- **State Preserved**: TODOs and context maintained
- **No Lost Work**: Everything recoverable
- **Continuous Flow**: Orchestration uninterrupted

## Summary

The `.claude/settings.json` file is **NOT OPTIONAL**. It is a critical component that:
1. Prevents data loss during compaction
2. Enables state recovery
3. Maintains orchestration continuity
4. Preserves TODO lists

**Every Software Factory project MUST have a properly configured settings.json file.**