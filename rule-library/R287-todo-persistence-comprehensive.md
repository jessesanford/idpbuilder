# 🚨🚨🚨 RULE R287 - Comprehensive TODO Persistence Protocol

**Criticality:** BLOCKING - TODO loss = Immediate failure  
**Grading Impact:** -20% to -100% depending on violation severity  
**Enforcement:** CONTINUOUS - Active throughout agent lifecycle

## Rule Statement

EVERY agent MUST persistently save, commit, and verify TODOs to prevent catastrophic work loss during compaction. This comprehensive rule consolidates and supersedes R187-R190.

## 1. MANDATORY SAVE TRIGGERS (Within 30 seconds)

### Event-Based Saves
- **TodoWrite Usage**: Save within 30 seconds after ANY TodoWrite operation
- **State Transitions**: Save BEFORE any state machine transition
- **Agent Spawning**: Save BEFORE spawning another agent
- **Completion Events**: Save AFTER completing effort/wave/review
- **Error Conditions**: Save WHEN encountering blocks/failures/violations

### Frequency-Based Saves
- **Every 10 messages** exchanged with user/orchestrator
- **Every 15 minutes** of active work (MANDATORY)
- **Every 200 lines** of code written
- **Every 3 files** modified

**Penalties**: -20% per missed trigger, -50% for lost TODOs

## 2. COMMIT & PUSH PROTOCOL (Within 60 seconds)

```bash
save_and_commit_todos() {
    local trigger="$1"
    local todo_file="$PROJECT_ROOT/todos/${AGENT_NAME}-${CURRENT_STATE}-$(date '+%Y%m%d-%H%M%S').todo"
    
    # Save TODO content
    cat > "$todo_file" <<EOF
# TODO State at $(date '+%Y-%m-%d %H:%M:%S')
# Agent: ${AGENT_NAME}
# State: ${CURRENT_STATE}
# Phase: ${CURRENT_PHASE}, Wave: ${CURRENT_WAVE}

## In Progress (${IN_PROGRESS_COUNT})
${IN_PROGRESS_ITEMS}

## Pending (${PENDING_COUNT})
${PENDING_ITEMS}

## Completed (${COMPLETED_COUNT})
${COMPLETED_ITEMS}

## Blocked (${BLOCKED_COUNT})
${BLOCKED_ITEMS}
EOF
    
    # Commit within 60 seconds
    cd "$PROJECT_ROOT"
    git add "todos/*.todo" "todos/audit.log"
    git commit -m "todo: ${AGENT_NAME} - ${trigger} - $(date '+%H:%M:%S')" \
               -m "State: ${CURRENT_STATE}" \
               -m "TODO counts: ${IN_PROGRESS_COUNT} in_progress, ${PENDING_COUNT} pending, ${COMPLETED_COUNT} completed"
    
    # Push immediately
    git push origin "$(git branch --show-current)" || {
        echo "🚨 CRITICAL: Push failed! Retrying..."
        git push --force-with-lease origin "$(git branch --show-current)"
    }
    
    # Log audit trail
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] SAVE: ${AGENT_NAME} - ${trigger}" >> "$PROJECT_ROOT/todos/audit.log"
}
```

**Penalties**: -25% for uncommitted files, -10% for delayed commits

## 3. RECOVERY VERIFICATION (After compaction)

```bash
recover_todos_after_compaction() {
    if [ -f /tmp/compaction_marker.txt ]; then
        echo "🚨 COMPACTION DETECTED - Initiating recovery"
        
        # Find latest TODO file
        LATEST_TODO=$(ls -t "$PROJECT_ROOT/todos/${AGENT_NAME}-"*.todo 2>/dev/null | head -1)
        
        if [ -z "$LATEST_TODO" ]; then
            echo "❌ CRITICAL: No TODO file found!"
            exit 1
        fi
        
        # MANDATORY: Load into TodoWrite tool (not just read!)
        TODO_CONTENT=$(cat "$LATEST_TODO")
        load_todos_to_todowrite_tool "$TODO_CONTENT"  # CRITICAL: Must use tool
        
        # Verify recovery
        verify_todo_counts_match "$LATEST_TODO"
        
        # Save post-recovery checkpoint
        save_and_commit_todos "POST_RECOVERY_CHECKPOINT"
    fi
}
```

**CRITICAL**: Must LOAD into TodoWrite tool, not just read file!  
**Penalties**: -30% for incomplete recovery, -100% for total loss

## 4. MONITORING & ENFORCEMENT

### Self-Monitoring Requirements
```bash
# Check save frequency
check_save_compliance() {
    local last_save_ago=$(($(date +%s) - LAST_TODO_SAVE))
    
    # 15-minute check
    if [ $last_save_ago -gt 900 ]; then
        echo "🚨 VIOLATION: Exceeding 15-minute requirement!"
        save_and_commit_todos "OVERDUE_CHECKPOINT"
    fi
    
    # 10-message check
    if [ $MESSAGE_COUNT -ge 10 ]; then
        save_and_commit_todos "10_MESSAGE_CHECKPOINT"
        MESSAGE_COUNT=0
    fi
}
```

### File Format
```
todos/${AGENT_NAME}-${STATE}-${YYYYMMDD-HHMMSS}.todo
```

Examples:
- `orchestrator-WAVE_COMPLETE-20250120-143000.todo`
- `sw-engineer-IMPLEMENTATION-20250120-145500.todo`

## 5. GRADING SUMMARY

### Violation Penalties
- **Save Triggers**: -20% per missed trigger
- **Save Frequency**: -15% per missed interval
- **Commit/Push**: -25% for uncommitted, -10% for delayed
- **Recovery**: -30% for incomplete recovery
- **Total Loss**: -100% IMMEDIATE FAILURE

### Critical Requirements
- ✅ Save within 30 seconds of triggers
- ✅ Save every 10 messages OR 15 minutes
- ✅ Commit within 60 seconds of save
- ✅ Push to remote immediately
- ✅ Load into TodoWrite during recovery (not just read)
- ✅ Maintain audit trail

## Integration Notes

This rule consolidates and replaces:
- R187 (Save Triggers) - DEPRECATED
- R188 (Save Frequency) - DEPRECATED  
- R189 (Commit Protocol) - DEPRECATED
- R190 (Recovery Verification) - DEPRECATED

All agents MUST follow this unified protocol. The save/commit/verify cycle is atomic and non-negotiable.

---
**Remember:** Compaction can occur ANY TIME. Only persistent TODOs survive!