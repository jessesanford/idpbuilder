# ⚠️ DEPRECATED - Subsumed by R287
This rule has been consolidated into R287-todo-persistence-comprehensive.md
Please refer to R287 for current TODO persistence requirements.

# 🚨🚨🚨 RULE R190 - TODO Recovery Verification [DEPRECATED]

**Criticality:** BLOCKING - Failed recovery = Cannot continue work  
**Grading Impact:** -30% for incomplete recovery  
**Enforcement:** IMMEDIATE - Must verify after every compaction

## Rule Statement

EVERY agent MUST successfully recover and verify TODOs after compaction, with zero tolerance for TODO loss.

## Recovery Protocol

### 1. Compaction Detection Response
```bash
# FIRST action after detecting compaction marker
if [ -f /tmp/compaction_marker.txt ]; then
    echo "🚨 COMPACTION DETECTED - Initiating TODO recovery"
    
    # Step 1: Identify agent
    AGENT_TYPE=$(identify_agent_type)  # orchestrator|sw-engineer|code-reviewer|architect
    
    # Step 2: Find latest TODO file
    TODO_DIR="$PROJECT_ROOT/todos"
    LATEST_TODO=$(ls -t "$TODO_DIR/${AGENT_TYPE}-"*.todo 2>/dev/null | head -1)
    
    if [ -z "$LATEST_TODO" ]; then
        echo "❌ CRITICAL: No TODO file found for $AGENT_TYPE"
        echo "🚨 VIOLATION: R190 - Cannot recover TODOs"
        exit 1
    fi
    
    # Step 3: Read TODO file
    echo "📖 Reading TODO file: $LATEST_TODO"
    TODO_CONTENT=$(cat "$LATEST_TODO")
    
    # Step 4: Load into TodoWrite tool
    echo "📝 Loading TODOs into TodoWrite tool..."
    # Parse and load - this is MANDATORY, not optional
    load_todos_to_memory "$TODO_CONTENT"
    
    # Step 5: Verify loading success
    verify_todo_recovery
fi
```

### 2. TODO Loading Process

**CRITICAL: Must use TodoWrite tool, not just read!**
```bash
load_todos_to_memory() {
    local content="$1"
    
    # Parse TODO file content
    local in_progress_items=$(echo "$content" | grep -A 100 "## In Progress" | grep "^- \[ \]")
    local pending_items=$(echo "$content" | grep -A 100 "## Pending" | grep "^- \[ \]")
    local completed_items=$(echo "$content" | grep -A 100 "## Completed" | grep "^- \[x\]")
    local blocked_items=$(echo "$content" | grep -A 100 "## Blocked" | grep "^- \[!\]")
    
    # MUST use TodoWrite to populate memory
    # This is NOT just reading - must actively load into tool
    
    echo "✅ Loaded:"
    echo "  - In Progress: $(echo "$in_progress_items" | wc -l)"
    echo "  - Pending: $(echo "$pending_items" | wc -l)"
    echo "  - Completed: $(echo "$completed_items" | wc -l)"
    echo "  - Blocked: $(echo "$blocked_items" | wc -l)"
}
```

### 3. Recovery Verification

**MUST verify all TODOs loaded correctly:**
```bash
verify_todo_recovery() {
    echo "🔍 Verifying TODO recovery..."
    
    # Check 1: TodoWrite tool has items
    if ! todo_write_has_items; then
        echo "❌ TodoWrite tool is empty after recovery!"
        return 1
    fi
    
    # Check 2: Counts match file
    local file_count=$(grep -c "^- \[" "$LATEST_TODO")
    local memory_count=$(get_todowrite_count)
    
    if [ "$file_count" -ne "$memory_count" ]; then
        echo "❌ Count mismatch! File: $file_count, Memory: $memory_count"
        return 1
    fi
    
    # Check 3: Critical items present
    if ! verify_critical_todos; then
        echo "❌ Critical TODOs missing after recovery!"
        return 1
    fi
    
    echo "✅ TODO recovery verified successfully"
    echo "📊 Recovered $memory_count TODOs from $LATEST_TODO"
    return 0
}
```

## Recovery Validation Checklist

### Mandatory Checks
```markdown
[ ] Compaction marker detected
[ ] Agent type identified correctly
[ ] Latest TODO file found
[ ] TODO file not empty
[ ] TODO file less than 15 minutes old (warning if older)
[ ] TodoWrite tool populated (not just read)
[ ] Counts match between file and memory
[ ] In-progress items recovered
[ ] Context information recovered
[ ] No duplicate items after merge
```

### State Recovery
```bash
# Must also recover state machine position
recover_state_context() {
    local state_line=$(grep "^# State:" "$LATEST_TODO")
    local phase_line=$(grep "^# Phase:" "$LATEST_TODO")
    
    CURRENT_STATE=$(echo "$state_line" | cut -d: -f2 | xargs)
    CURRENT_PHASE=$(echo "$phase_line" | sed 's/.*Phase: \([0-9]\+\).*/\1/')
    CURRENT_WAVE=$(echo "$phase_line" | sed 's/.*Wave: \([0-9]\+\).*/\1/')
    
    echo "📍 Recovered state context:"
    echo "  - State: $CURRENT_STATE"
    echo "  - Phase: $CURRENT_PHASE"
    echo "  - Wave: $CURRENT_WAVE"
}
```

## Recovery Failure Handling

### Failure Modes and Responses

1. **No TODO file exists**
```bash
# Attempt reconstruction from git history
echo "⚠️ No TODO file found, attempting git recovery..."
git log --grep="^todo: ${AGENT_TYPE}" -n 1 --format="%H %s"
git show HEAD:todos/ | grep "${AGENT_TYPE}-.*\.todo"
```

2. **TODO file too old**
```bash
# If file >30 minutes old
FILE_AGE=$(file_age_minutes "$LATEST_TODO")
if [ "$FILE_AGE" -gt 30 ]; then
    echo "⚠️ WARNING: TODO file is $FILE_AGE minutes old"
    echo "📝 May be missing recent work"
    echo "🔄 Attempting to reconstruct from git commits..."
fi
```

3. **Corrupt TODO file**
```bash
# If file exists but can't be parsed
if ! validate_todo_format "$LATEST_TODO"; then
    echo "❌ TODO file format invalid"
    echo "🔄 Attempting previous version..."
    PREVIOUS_TODO=$(ls -t "$TODO_DIR/${AGENT_TYPE}-"*.todo | head -2 | tail -1)
fi
```

## Deduplication Protocol

### After Loading TODOs
```bash
deduplicate_todos() {
    echo "🔄 Deduplicating TODOs..."
    
    # Remove exact duplicates
    local before_count=$(get_todowrite_count)
    remove_duplicate_todos
    local after_count=$(get_todowrite_count)
    
    if [ "$before_count" -ne "$after_count" ]; then
        echo "✅ Removed $((before_count - after_count)) duplicates"
    fi
}
```

## Grading Verification

### Success Criteria
```bash
grade_recovery() {
    local score=100
    
    # Check each requirement
    [ ! -f "$LATEST_TODO" ] && score=$((score - 30))
    [ "$(file_age_minutes "$LATEST_TODO")" -gt 15 ] && score=$((score - 15))
    ! todo_write_has_items && score=$((score - 30))
    ! verify_critical_todos && score=$((score - 25))
    
    echo "📊 Recovery Grade: ${score}%"
    
    if [ "$score" -lt 70 ]; then
        echo "❌ FAILED: Recovery below passing grade"
        return 1
    fi
}
```

## Post-Recovery Actions

### Immediate Save
```bash
# After successful recovery, immediately save current state
post_recovery_save() {
    echo "💾 Creating post-recovery checkpoint..."
    save_todos "POST_RECOVERY_CHECKPOINT"
    
    # Include recovery metadata
    echo "# Recovery completed at $(date '+%Y-%m-%d %H:%M:%S')" >> "$TODO_FILE"
    echo "# Recovered from: $LATEST_TODO" >> "$TODO_FILE"
    echo "# Recovery score: ${score}%" >> "$TODO_FILE"
    
    # Commit immediately
    git add "$TODO_FILE"
    git commit -m "todo: ${AGENT_TYPE} - post-recovery checkpoint"
    git push
}
```

## Integration Requirements

### Works with R186 (Detection)
- R186 detects compaction
- R190 handles recovery

### Works with R187-R189 (Saving)
- Better saves = better recovery
- Recent saves = less work lost

## Example Recovery Flow

```
1. Compaction occurs
2. Agent restarts
3. Detects /tmp/compaction_marker.txt (R186)
4. Finds latest TODO: orchestrator-WAVE_COMPLETE-20250120-142000.todo
5. Reads file content
6. LOADS into TodoWrite tool (not just reads!)
7. Verifies counts match
8. Deduplicates if needed
9. Confirms critical items present
10. Saves post-recovery checkpoint
11. Continues work with full context
```

## Bad Recovery Example

```
1. Compaction occurs
2. Agent restarts
3. Detects marker
4. Finds TODO file
5. ONLY READS file (doesn't load to TodoWrite) ❌
6. Continues work without TODOs in memory ❌
7. Creates new TODOs, losing all previous ❌
8. GRADING FAILURE: -30% for failed recovery
```

---
**Remember:** Reading TODOs is NOT enough. You MUST load them into TodoWrite tool!