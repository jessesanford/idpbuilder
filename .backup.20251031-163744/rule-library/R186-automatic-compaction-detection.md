---
name: R186 - Automatic Compaction Detection
criticality: BLOCKING
agent: all
state: INIT
---

# 🚨🚨🚨 RULE R186 - Automatic Compaction Detection 🚨🚨🚨

## Rule Statement
ALL agents MUST check for context compaction markers during startup and after significant operations. If compaction is detected, agents must immediately recover their TODO state before continuing work.

## Rationale
Context compaction can occur at any time due to:
- Large file processing
- Extended conversations
- Memory constraints
- System limitations

Without proper detection and recovery:
- TODO lists are lost
- Work progress disappears
- Agents restart from scratch
- Duplication of effort occurs

## Implementation Requirements

### 1. Compaction Detection Script
```bash
# Use the check-compaction-agent.sh utility script
if [ -f "$HOME/.claude/utilities/check-compaction-agent.sh" ]; then
    bash "$HOME/.claude/utilities/check-compaction-agent.sh" ${AGENT_NAME}
elif [ -f "/home/user/.claude/utilities/check-compaction-agent.sh" ]; then
    bash "/home/user/.claude/utilities/check-compaction-agent.sh" ${AGENT_NAME}
elif [ -f "./utilities/check-compaction-agent.sh" ]; then
    bash "./utilities/check-compaction-agent.sh" ${AGENT_NAME}
else
    echo "⚠️⚠️⚠️ Compaction check script not found, using fallback"
    # Fallback detection
    if [ -f /tmp/compaction_marker.txt ]; then
        echo "🚨🚨🚨 COMPACTION DETECTED!"
        cat /tmp/compaction_marker.txt
        rm -f /tmp/compaction_marker.txt
        echo "RECOVER TODOs NOW - See R190"
        # Trigger R190 recovery
    else
        echo "No compaction detected"
    fi
fi
```

### 2. Compaction Markers
```bash
# Standard marker locations
/tmp/compaction_marker.txt
/tmp/${AGENT_NAME}_compaction.flag
./compaction_detected.flag

# Marker format
echo "COMPACTION_TIME: $(date '+%Y%m%d-%H%M%S')" > /tmp/compaction_marker.txt
echo "AGENT: ${AGENT_NAME}" >> /tmp/compaction_marker.txt
echo "LAST_STATE: ${CURRENT_STATE}" >> /tmp/compaction_marker.txt
```

### 3. Recovery Trigger
```bash
# If compaction detected, must recover TODOs (R190)
if [ "$COMPACTION_DETECTED" = "true" ]; then
    echo "🚨🚨🚨 R186: Compaction detected - triggering recovery"
    recover_todos_after_compaction  # R190 function
    echo "✅ R186: Recovery complete, continuing work"
fi
```

## Enforcement
- **Trigger**: At startup and after major operations
- **Detection**: Check for marker files
- **Recovery**: Invoke R190 TODO recovery
- **Verification**: Ensure TODOs loaded before continuing

## Related Rules
- R171: PreCompact Hook
- R174: Recovery Detection
- R190: TODO Recovery Verification

---
**Status**: STUB - This rule file needs complete implementation details
**Created**: By software-factory-manager during rule library audit