# Compaction Detection Fix Documentation

## Problem Statement
Agent configuration files contain bash commands with backslash line continuations (`\`) for readability. However, when agents copy and execute these as single-line commands, the backslashes cause syntax errors like:
```
syntax error near unexpected token `then'
```

## Solution

### 1. Script-Based Approach (RECOMMENDED)
Instead of inline bash with backslashes, use dedicated scripts:

```bash
# For general compaction check:
bash /workspaces/software-factory-2.0-template/.claude/scripts/check-compaction.sh

# For agent-specific checks:
bash /workspaces/software-factory-2.0-template/.claude/scripts/check-compaction-agent.sh orchestrator
bash /workspaces/software-factory-2.0-template/.claude/scripts/check-compaction-agent.sh sw-engineer
bash /workspaces/software-factory-2.0-template/.claude/scripts/check-compaction-agent.sh code-reviewer
bash /workspaces/software-factory-2.0-template/.claude/scripts/check-compaction-agent.sh architect
```

### 2. Single-Line Command (FALLBACK)
If scripts are unavailable, use this simplified single-line version WITHOUT backslashes:

```bash
if [ -f /tmp/compaction_marker.txt ]; then echo "⚠️ COMPACTION DETECTED"; cat /tmp/compaction_marker.txt; rm -f /tmp/compaction_marker.txt; echo "MUST RECOVER TODOs"; exit 0; else echo "No compaction detected"; fi
```

## What Was Fixed

1. **Created standalone scripts** that don't require backslash escaping
2. **Provided single-line alternatives** without backslashes for emergency use
3. **Documented the issue** so agents understand why backslashes cause problems

## For Agents

**DO NOT** copy multi-line bash with backslashes and run as single line.

**DO** one of these:
- Run the provided scripts: `bash /workspaces/software-factory-2.0-template/.claude/scripts/check-compaction-agent.sh [agent-type]`
- Use the simplified single-line version without backslashes
- Keep the multi-line format intact (copy as multi-line, paste as multi-line)

## Testing

Test the scripts work correctly:
```bash
# Create test marker
echo "TEST" > /tmp/compaction_marker.txt

# Test script
bash /workspaces/software-factory-2.0-template/.claude/scripts/check-compaction.sh

# Should detect compaction and remove marker
```