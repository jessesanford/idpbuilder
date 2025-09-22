# ⚠️ DEPRECATED - Subsumed by R287
This rule has been consolidated into R287-todo-persistence-comprehensive.md
Please refer to R287 for current TODO persistence requirements.

# 🚨🚨🚨 RULE R189 - TODO File Commit Protocol [DEPRECATED]

**Criticality:** BLOCKING - Uncommitted saves = Inaccessible after compaction  
**Grading Impact:** -25% for uncommitted TODO files  
**Enforcement:** IMMEDIATE - Must commit within 60 seconds of save

## Rule Statement

EVERY TODO save MUST be committed and pushed to git within 60 seconds, ensuring persistence beyond local filesystem.

## Commit Requirements

### 1. Immediate Commit After Save
```bash
save_todos() {
    local trigger="$1"
    local todo_file="$PROJECT_ROOT/todos/${AGENT_NAME}-${CURRENT_STATE}-$(date '+%Y%m%d-%H%M%S').todo"
    
    # Save TODOs
    echo "# TODO State - $trigger" > "$todo_file"
    # ... write TODO content ...
    
    # MANDATORY: Commit within 60 seconds
    cd "$PROJECT_ROOT"
    git add "todos/*.todo"
    git add "todos/audit.log"
    
    # Structured commit message
    git commit -m "todo: ${AGENT_NAME} - ${trigger} - $(date '+%H:%M:%S')" 
               -m "State: ${CURRENT_STATE}" 
               -m "Phase: ${CURRENT_PHASE}, Wave: ${CURRENT_WAVE}" 
               -m "Trigger: ${trigger}" 
               -m "TODO counts: ${IN_PROGRESS} in_progress, ${PENDING} pending, ${COMPLETED} completed"
    
    # MANDATORY: Push immediately
    git push origin "$(git branch --show-current)"
    
    # Verify success
    if [ $? -eq 0 ]; then
        echo "✅ TODOs saved, committed, and pushed successfully"
    else
        echo "❌ CRITICAL: Failed to push TODO commit!"
        echo "🚨 VIOLATION: R189 - TODOs not persisted to remote"
        return 1
    fi
}
```

### 2. Commit Message Standards

**Format:**
```
todo: [agent-name] - [trigger-reason] - [timestamp]

State: [current-state-machine-state]
Phase: [X], Wave: [Y]
Trigger: [save-trigger-from-R187/R188]
TODO counts: [X] in_progress, [Y] pending, [Z] completed
```

**Examples:**
```bash
# Good commit messages
git commit -m "todo: orchestrator - STATE_TRANSITION - 14:30:22" 
           -m "State: WAVE_COMPLETE" 
           -m "Phase: 1, Wave: 2" 
           -m "Trigger: Transitioning to INTEGRATION_REVIEW" 
           -m "TODO counts: 1 in_progress, 5 pending, 12 completed"

git commit -m "todo: sw-engineer - 10_MESSAGE_CHECKPOINT - 15:45:10" 
           -m "State: IMPLEMENTATION" 
           -m "Phase: 2, Wave: 1" 
           -m "Trigger: 10 message checkpoint" 
           -m "TODO counts: 2 in_progress, 3 pending, 8 completed"
```

### 3. Push Requirements

**MUST push to correct remote:**
```bash
# Verify remote exists
if ! git remote | grep -q "origin"; then
    echo "❌ No remote configured - TODOs cannot be persisted!"
    exit 1
fi

# Push with verification
git push origin "$(git branch --show-current)" || {
    echo "❌ Push failed - attempting force push"
    git push --force-with-lease origin "$(git branch --show-current)" || {
        echo "🚨 CRITICAL: Cannot push TODOs to remote!"
        exit 1
    }
}
```

## Multi-Agent Coordination

### Prevent Conflicts
```bash
# Before committing, pull latest
git pull --rebase origin "$(git branch --show-current)"

# If conflicts in todos/ directory
git status --porcelain | grep "^UU todos/" && {
    echo "⚠️ TODO conflicts detected - using theirs"
    git checkout --theirs todos/
    git add todos/
    git rebase --continue
}
```

### Agent Namespacing
Each agent MUST use unique filenames:
- `orchestrator-*.todo`
- `sw-engineer-*.todo`
- `code-reviewer-*.todo`
- `architect-*.todo`

## Cleanup Protocol

### Old File Management
```bash
# After successful commit, clean old files (keep last 5 per agent)
cleanup_old_todos() {
    for agent in orchestrator sw-engineer code-reviewer architect; do
        # List files older than 5 most recent
        ls -t "$PROJECT_ROOT/todos/${agent}-"*.todo 2>/dev/null | tail -n +6 | while read old_file; do
            git rm "$old_file"
        done
    done
    
    # Commit cleanup if files were removed
    if git status --porcelain | grep -q "^D"; then
        git commit -m "cleanup: remove old TODO files (keeping last 5 per agent)"
        git push
    fi
}
```

## Audit Requirements

### Commit Log Trail
```bash
# Maintain audit log of all commits
echo "[$(date '+%Y-%m-%d %H:%M:%S')] COMMIT: ${AGENT_NAME} - $(git rev-parse --short HEAD)" >> "$PROJECT_ROOT/todos/commit-audit.log"

# Include in commit
git add todos/commit-audit.log
```

### Verification Commands
```bash
# Agent can verify TODO persistence
verify_todo_persistence() {
    # Check local commits
    local commits=$(git log --oneline --grep="^todo:" -n 5)
    echo "Recent TODO commits:"
    echo "$commits"
    
    # Check remote sync
    local behind=$(git rev-list --count origin/$(git branch --show-current)..HEAD)
    if [ "$behind" -gt 0 ]; then
        echo "⚠️ WARNING: $behind TODO commits not pushed!"
        return 1
    fi
    
    echo "✅ All TODO commits pushed to remote"
}
```

## Recovery Scenarios

### After Compaction
```bash
# TODOs are recoverable from git
recover_todos_from_git() {
    local agent="$1"
    
    # Find latest TODO file for agent from git
    local latest=$(git ls-files "todos/${agent}-*.todo" | tail -1)
    
    if [ -n "$latest" ]; then
        echo "✅ Recovering TODOs from: $latest"
        git checkout HEAD -- "$latest"
        cat "$latest"
    else
        echo "❌ No TODO files found in git for $agent"
        return 1
    fi
}
```

### Network Failures
```bash
# If push fails due to network
handle_push_failure() {
    local retry_count=0
    local max_retries=3
    
    while [ $retry_count -lt $max_retries ]; do
        if git push origin "$(git branch --show-current)"; then
            echo "✅ Push succeeded on retry $((retry_count + 1))"
            return 0
        fi
        
        retry_count=$((retry_count + 1))
        echo "⚠️ Push failed, retry $retry_count of $max_retries"
        sleep 5
    done
    
    echo "🚨 CRITICAL: Cannot push after $max_retries retries"
    echo "📝 TODOs saved locally but not persisted to remote"
    return 1
}
```

## Grading Enforcement

### Violations and Penalties
- TODO file not committed: -25%
- Commit not pushed: -25%
- Missing commit message details: -10%
- Cleanup not performed: -5%
- Audit trail gaps: -10%

### Critical Failures
- No commits before compaction: -50%
- Remote out of sync >30 minutes: -40%
- Lost TODOs due to no commits: -75%

## Integration Checklist

```bash
# After every TODO save:
[ ] File created in todos/ directory
[ ] Git add performed
[ ] Commit with proper message
[ ] Push to remote
[ ] Audit log updated
[ ] Old files cleaned up
[ ] Verification passed
```

---
**Remember:** Local saves are NOT enough. Only committed and pushed TODOs survive compaction!