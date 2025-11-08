# ⚠️ DEPRECATED - Subsumed by R287
This rule has been consolidated into R287-todo-persistence-comprehensive.md
Please refer to R287 for current TODO persistence requirements.

# 🚨🚨🚨 RULE R187 - Mandatory TODO Save Triggers [DEPRECATED]

**Criticality:** BLOCKING - Failure to save = Potential total work loss  
**Grading Impact:** -20% for each missed save trigger  
**Enforcement:** IMMEDIATE - Must save within 30 seconds of trigger

## Rule Statement

EVERY agent MUST save TODOs to disk when ANY of these triggers occur:

## Mandatory Save Triggers

### 1. TodoWrite Tool Usage (+30 seconds)
```bash
# AFTER using TodoWrite tool, within 30 seconds:
TODO_FILE="$PROJECT_ROOT/todos/${AGENT_NAME}-${CURRENT_STATE}-$(date '+%Y%m%d-%H%M%S').todo"
echo "# TODO State at $(date '+%Y-%m-%d %H:%M:%S')" > "$TODO_FILE"
echo "# Agent: ${AGENT_NAME}" >> "$TODO_FILE"
echo "# State: ${CURRENT_STATE}" >> "$TODO_FILE"
echo "" >> "$TODO_FILE"
# [Write all TODOs with status indicators]
```

### 2. State Machine Transitions
**BEFORE** transitioning between ANY states:
- WAVE_COMPLETE → INTEGRATE_WAVE_EFFORTS_REVIEW
- IMPLEMENTATION → MEASURE_SIZE
- CODE_REVIEW → CREATE_SPLIT_PLAN
- ANY_STATE → ANY_OTHER_STATE

### 3. Agent Spawn Events
**BEFORE** spawning another agent:
- Save current orchestrator TODOs
- Include pending spawn tasks
- Document expected outcomes

### 4. Completion Milestones
**IMMEDIATELY** after:
- Completing an effort
- Finishing a wave
- Passing a review
- Fixing review issues
- Creating integration branch

### 5. Error/Blocking Conditions
**WHEN** encountering:
- Line count violations (>800)
- Failed tests
- Review rejections
- Integration conflicts
- Any blocking issue

## Save Format Requirements

```markdown
# TODO State at [TIMESTAMP]
# Agent: [orchestrator|sw-engineer|code-reviewer|architect]
# State: [CURRENT_STATE_MACHINE_STATE]
# Phase: [X] Wave: [Y] Effort: [NAME]

## In Progress (1)
- [ ] Current task description [STATUS: 60% complete]

## Pending (3)
- [ ] Next task 1
- [ ] Next task 2  
- [ ] Next task 3

## Completed (5)
- [x] Completed task 1
- [x] Completed task 2
- [x] Completed task 3
- [x] Completed task 4
- [x] Completed task 5

## Blocked (1)
- [!] Blocked task [REASON: Waiting for architect review]

## Context Notes
- Current branch: phase1/wave2/effort-api-types
- Line count: 650/800
- Last measurement: 2 minutes ago
- Dependencies: effort-controllers must complete first
```

## Verification Protocol

```bash
# Agent MUST verify save success
if [ ! -f "$TODO_FILE" ]; then
    echo "❌ CRITICAL: TODO save failed!"
    echo "🚨 GRADING VIOLATION: R187 - Failed to save TODOs"
    exit 1
fi

# Verify non-empty
if [ ! -s "$TODO_FILE" ]; then
    echo "❌ CRITICAL: TODO file is empty!"
    exit 1
fi

echo "✅ TODOs saved to: $TODO_FILE"
echo "📊 File size: $(wc -l "$TODO_FILE" | awk '{print $1}') lines"
```

## Grading Enforcement

### Automatic Failures
- Missing save after TodoWrite usage: -20%
- Missing save at state transition: -20%
- Missing save before agent spawn: -30%
- Lost TODOs during compaction: -50%

### Audit Trail Required
```bash
# Every save must be logged
echo "[$(date '+%Y-%m-%d %H:%M:%S')] TODO_SAVE: ${AGENT_NAME} saved ${TODO_FILE}" >> "$PROJECT_ROOT/todos/audit.log"
```

## Recovery After Compaction

If TODOs weren't saved and compaction occurred:
1. **IMMEDIATE STOP** - Do not continue
2. **Report Loss** - Document what was lost
3. **Reconstruct** - Attempt to rebuild from git history
4. **Grading Penalty** - Automatic -50% for work loss

## Example Compliance

```bash
# GOOD: After TodoWrite
$> TodoWrite updates tasks...
$> # Within 30 seconds:
$> save_todos "POST_TODOWRITE"
✅ TODOs saved to: todos/orchestrator-WAVE_COMPLETE-20250120-143022.todo

# BAD: Delayed save
$> TodoWrite updates tasks...
$> # 5 minutes pass...
$> # Other work continues...
❌ VIOLATION: R187 - TODOs not saved within 30 seconds of TodoWrite
```

## Integration with R189

After saving (R187), MUST commit (R189):
```bash
cd "$PROJECT_ROOT"
git add todos/*.todo
git commit -m "todo: ${AGENT_NAME} - ${TRIGGER_REASON}"
git push
```

---
**Remember:** PreCompact CANNOT save your TODOs for you. Only YOU can prevent TODO loss!