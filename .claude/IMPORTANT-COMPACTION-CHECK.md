# 🚨 CRITICAL: How to Check for Compaction

## For All Agents - The Scripts Are Now in Utilities:

The compaction check scripts are installed in multiple locations by the setup scripts:
- `$HOME/.claude/utilities/check-compaction-agent.sh`
- `/home/user/.claude/utilities/check-compaction-agent.sh`  
- `./utilities/check-compaction-agent.sh`

## Use This Command Pattern:

```bash
# This checks all three locations automatically:
if [ -f "$HOME/.claude/utilities/check-compaction-agent.sh" ]; then
    bash "$HOME/.claude/utilities/check-compaction-agent.sh" [your-agent-type]
elif [ -f "/home/user/.claude/utilities/check-compaction-agent.sh" ]; then
    bash "/home/user/.claude/utilities/check-compaction-agent.sh" [your-agent-type]
elif [ -f "./utilities/check-compaction-agent.sh" ]; then
    bash "./utilities/check-compaction-agent.sh" [your-agent-type]
else
    echo "Compaction check script not found"
fi

# Agent type examples: orchestrator, sw-engineer, code-reviewer, architect
```

## Why This Matters

The agent configuration files contain bash commands with backslashes (`\`) for multi-line formatting. When you copy and run these as single lines, the backslashes cause syntax errors.

## DO NOT Use These (They Will Fail):

❌ Don't copy multi-line bash with backslashes from agent configs
❌ Don't try to run commands like: `if [ -f /tmp/compaction_marker.txt ]; then \`

## If Script Is Not Available

Use this simple single-line command instead:
```bash
if [ -f /tmp/compaction_marker.txt ]; then echo "COMPACTION DETECTED"; cat /tmp/compaction_marker.txt; rm -f /tmp/compaction_marker.txt; echo "RECOVER TODOs NOW"; exit 0; else echo "No compaction"; fi
```

## Remember

1. Check for compaction FIRST before any other work
2. If compaction detected, recover TODOs before proceeding
3. Use the script approach - it's safer and more reliable