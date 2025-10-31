# R171.0.0 - PreCompact Hook Configuration

## Rule Statement
The PreCompact hook MUST be configured in `.claude/settings.json` as the ONLY automatic hook for state preservation.

## Rationale
Claude Code only reliably supports certain hooks, with PreCompact being the most critical for context preservation.

## Implementation

### Required Configuration
```json
{
  "hooks": {
    "PreCompact": [
      {
        "matcher": "auto",
        "hooks": [{
          "type": "command",
          "command": "echo 'marker creation command here'"
        }]
      },
      {
        "matcher": "manual",
        "hooks": [{
          "type": "command",
          "command": "echo 'marker creation command here'"
        }]
      }
    ]
  }
}
```

### Hook Behavior
1. Runs automatically before context compaction
2. Creates `/tmp/compaction_marker.txt` with context info
3. Preserves latest TODO file if available
4. Uses relative paths for portability

## Enforcement
- Setup.sh MUST copy correct settings.json
- No other hooks should be configured
- Hook command must be inline (type: "command")

## Validation
```bash
# Check hook configuration
cat .claude/settings.json | grep -A 5 "PreCompact"

# Verify marker creation after compaction
ls -la /tmp/compaction_marker.txt
```

## Common Mistakes
- ❌ Adding fictional hooks (post_compaction, etc.)
- ❌ Using external script files instead of inline commands
- ❌ Expecting automatic execution of utility scripts
- ✅ Only PreCompact hook with inline commands