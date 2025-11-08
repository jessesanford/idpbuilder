# ⚠️ DEPRECATED - Subsumed by R287
This rule has been consolidated into R287-todo-persistence-comprehensive.md
Please refer to R287 for current TODO persistence requirements.

# 🚨🚨🚨 RULE R188 - TODO Save Frequency Requirements [DEPRECATED]

**Criticality:** BLOCKING - Insufficient saves = Guaranteed work loss  
**Grading Impact:** -15% for each frequency violation  
**Enforcement:** CONTINUOUS - Tracked throughout session

## Rule Statement

EVERY agent MUST save TODOs at MINIMUM these frequencies to prevent catastrophic loss during compaction.

## Mandatory Save Frequencies

### 1. Time-Based Intervals
**EVERY 10 MESSAGES** exchanged with user/orchestrator:
```bash
MESSAGE_COUNT=$((MESSAGE_COUNT + 1))
if [ $((MESSAGE_COUNT % 10)) -eq 0 ]; then
    echo "⏰ 10-message checkpoint - MUST save TODOs"
    save_todos "10_MESSAGE_CHECKPOINT"
fi
```

**EVERY 15 MINUTES** of active work:
```bash
# Set timer at agent startup
LAST_TODO_SAVE=$(date +%s)

# Check periodically
CURRENT_TIME=$(date +%s)
ELAPSED=$((CURRENT_TIME - LAST_TODO_SAVE))
if [ $ELAPSED -ge 900 ]; then  # 900 seconds = 15 minutes
    echo "⏰ 15-minute checkpoint - MUST save TODOs"
    save_todos "15_MINUTE_CHECKPOINT"
    LAST_TODO_SAVE=$CURRENT_TIME
fi
```

### 2. Work-Based Checkpoints

**AFTER EVERY:**
- 200 lines of code written
- 3 files modified
- 1 test suite run
- 1 review cycle completed
- 1 git commit created

### 3. Critical Operation Saves

**BEFORE high-risk operations:**
```bash
# Before any operation that could trigger compaction
echo "🚨 High-memory operation ahead - saving TODOs first"
save_todos "PRE_HIGH_MEMORY_OP"

# Operations requiring pre-save:
- Large file reads (>1000 lines)
- Multiple file operations (>5 files)
- Complex tool chains
- Recursive searches
- Large output generation
```

### 4. TODO Change Threshold

**When TODOs significantly change:**
```bash
# If >30% of TODOs change status
TODO_CHANGES=$((COMPLETED_COUNT - LAST_COMPLETED_COUNT))
CHANGE_PERCENT=$((TODO_CHANGES * 100 / TOTAL_TODOS))

if [ $CHANGE_PERCENT -ge 30 ]; then
    echo "📊 Significant TODO changes (${CHANGE_PERCENT}%) - saving"
    save_todos "SIGNIFICANT_CHANGE"
fi
```

## Frequency Tracking System

```bash
# Agent must maintain frequency tracking
TODO_SAVE_TRACKER="$PROJECT_ROOT/todos/.save-tracker"

track_save() {
    local reason="$1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] SAVE: $reason" >> "$TODO_SAVE_TRACKER"
}

check_save_health() {
    local last_save=$(tail -1 "$TODO_SAVE_TRACKER" | cut -d' ' -f1-2 | tr -d '[]')
    local last_epoch=$(date -d "$last_save" +%s 2>/dev/null || echo 0)
    local current_epoch=$(date +%s)
    local elapsed=$((current_epoch - last_epoch))
    
    if [ $elapsed -gt 900 ]; then  # >15 minutes
        echo "⚠️ WARNING: No TODO save in $((elapsed/60)) minutes!"
        echo "🚨 VIOLATION: R188 - Exceeding save frequency requirement"
        return 1
    fi
}
```

## Save Quality Requirements

### Minimum Content per Save
```markdown
# MUST include:
1. At least current in_progress items
2. All pending items  
3. Recent completed items (last 5)
4. Any blocked items with reasons
5. Current context (state, phase, wave)
```

### Incremental vs Full Saves
- **Incremental** (every 10 messages): Can save just changes
- **Full** (every 15 minutes): MUST save complete state

## Automated Reminders

Agents should self-monitor:
```bash
# In agent prompt/context
"⏰ REMINDER: Last TODO save was X minutes ago"
"📊 Messages since last save: Y"
"🚨 Approaching 15-minute save deadline"
```

## Grading Penalties

### Frequency Violations
- Missing 10-message save: -15%
- Missing 15-minute save: -15%  
- Missing critical operation save: -20%
- Pattern of delays: -30%

### Catastrophic Failures
- No saves before compaction: -50%
- Lost >1 hour of work: -75%
- Complete TODO loss: -100% (IMMEDIATE FAILURE)

## Integration with Other Rules

### Works with R187 (Triggers)
- R187 defines WHEN to save (events)
- R188 defines HOW OFTEN (frequency)
- Both must be satisfied

### Works with R189 (Commits)
- Every save must be committed
- Frequency applies to commits too

### Works with R190 (Recovery)
- More frequent saves = better recovery
- Less work lost during compaction

## Example Timeline

```
09:00 - Agent starts, initial save ✅
09:10 - 10 messages exchanged, checkpoint save ✅
09:15 - 15-minute timer, full save ✅
09:18 - State transition, trigger save ✅
09:25 - 10 more messages, checkpoint save ✅
09:30 - 15-minute timer, full save ✅
09:32 - TodoWrite used, trigger save ✅
09:40 - 10 more messages, checkpoint save ✅
09:45 - 15-minute timer, full save ✅
09:47 - COMPACTION OCCURS
        → Latest save is only 2 minutes old ✅
        → Minimal work loss
```

## Bad Example Timeline

```
09:00 - Agent starts, no initial save ❌
09:20 - 20 messages, no saves ❌
09:35 - 35 minutes, no saves ❌
09:47 - COMPACTION OCCURS
        → No saves available ❌
        → TOTAL WORK LOSS
        → -100% GRADING FAILURE
```

---
**Remember:** Compaction can happen ANY TIME. Only frequent saves protect your work!