# 🚨🚨🚨 RULE R195 - Branch Push Verification

**Criticality:** BLOCKING - Unpushed = Unpersisted = Lost  
**Grading Impact:** -35% for push verification failures  
**Enforcement:** AFTER every commit, BEFORE every transition

## Rule Statement

EVERY commit MUST be verified as pushed to remote. Agents MUST confirm push success and handle failures. NO work transitions without push verification.

## Push Verification Protocol

### 1. Post-Commit Verification
```bash
# MANDATORY after EVERY commit
verify_commit_pushed() {
    local commit_hash=$(git rev-parse HEAD)
    local branch=$(git branch --show-current)
    
    echo "🔍 Verifying commit $commit_hash is pushed..."
    
    # First, attempt push
    if ! git push; then
        echo "⚠️ Initial push failed, investigating..."
        
        # Check if branch exists on remote
        if ! git ls-remote --heads origin "$branch" | grep -q "$branch"; then
            echo "🔧 Branch doesn't exist on remote, creating..."
            git push -u origin "$branch"
        else
            echo "🔧 Branch exists, attempting force-with-lease..."
            git push --force-with-lease
        fi
    fi
    
    # Verify commit exists on remote
    if git branch -r --contains "$commit_hash" | grep -q "origin/$branch"; then
        echo "✅ Commit $commit_hash verified on origin/$branch"
        return 0
    else
        echo "❌ VIOLATION R195: Commit not on remote after push!"
        return 1
    fi
}
```

### 2. Batch Commit Verification
```bash
# For multiple commits
verify_all_commits_pushed() {
    local branch=$(git branch --show-current)
    local unpushed=$(git log origin/"$branch"..HEAD --oneline 2>/dev/null)
    
    if [ -z "$unpushed" ]; then
        echo "✅ All commits pushed to origin/$branch"
        return 0
    fi
    
    echo "⚠️ Unpushed commits detected:"
    echo "$unpushed"
    
    echo "📤 Pushing all commits..."
    if git push; then
        echo "✅ All commits now pushed"
        return 0
    else
        echo "❌ Push failed! Attempting recovery..."
        return 1
    fi
}
```

### 3. Pre-Transition Verification
```bash
# BEFORE any state transition or agent spawn
pre_transition_push_check() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "PRE-TRANSITION PUSH VERIFICATION"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    # Check all effort directories
    for effort_dir in "$SF_ROOT/efforts/phase${CURRENT_PHASE}/wave${CURRENT_WAVE}"/*/; do
        if [ -d "$effort_dir/.git" ]; then
            echo "Checking $(basename "$effort_dir")..."
            cd "$effort_dir"
            
            # Verify clean and pushed
            if [ -n "$(git status --porcelain)" ]; then
                echo "❌ Uncommitted changes in $(basename "$effort_dir")"
                git status --short
                return 1
            fi
            
            if ! verify_all_commits_pushed; then
                echo "❌ Unpushed commits in $(basename "$effort_dir")"
                return 1
            fi
        fi
    done
    
    echo "✅ All efforts verified pushed"
    cd "$SF_ROOT"
}
```

## Push Recovery Strategies

### 1. Network Failure Recovery
```bash
handle_push_network_failure() {
    local max_retries=5
    local retry_delay=10
    local attempt=0
    
    while [ $attempt -lt $max_retries ]; do
        echo "🔄 Push attempt $((attempt + 1))/$max_retries..."
        
        if git push; then
            echo "✅ Push succeeded on attempt $((attempt + 1))"
            return 0
        fi
        
        echo "⚠️ Push failed, waiting ${retry_delay}s..."
        sleep $retry_delay
        attempt=$((attempt + 1))
        retry_delay=$((retry_delay * 2))  # Exponential backoff
    done
    
    echo "❌ Push failed after $max_retries attempts"
    echo "📝 Saving commits locally for manual recovery"
    git format-patch origin/$(git branch --show-current)
    return 1
}
```

### 2. Conflict Resolution
```bash
handle_push_conflict() {
    local branch=$(git branch --show-current)
    
    echo "⚠️ Push conflict detected"
    
    # Try rebase first
    echo "🔧 Attempting rebase..."
    if git pull --rebase origin "$branch"; then
        echo "✅ Rebase successful, pushing..."
        git push
    else
        echo "❌ Rebase failed, conflicts need resolution"
        
        # Show conflict status
        git status
        
        # Agent must resolve or escalate
        echo "🚨 BLOCKING: Conflicts must be resolved before continuing"
        return 1
    fi
}
```

### 3. Force Push Safety
```bash
safe_force_push() {
    local branch=$(git branch --show-current)
    
    echo "⚠️ Considering force push for $branch"
    
    # Never force push integration branches
    if [[ "$branch" =~ integration ]]; then
        echo "❌ NEVER force push integration branches!"
        return 1
    fi
    
    # Check if we're the only one on this branch
    git fetch origin "$branch"
    local others=$(git log HEAD..origin/"$branch" --oneline)
    
    if [ -n "$others" ]; then
        echo "⚠️ Remote has commits not in local:"
        echo "$others"
        echo "❌ Cannot force push - would lose others' work"
        return 1
    fi
    
    # Safe to force with lease
    echo "🔧 Force pushing with lease..."
    git push --force-with-lease origin "$branch"
}
```

## Verification Automation

### Git Hooks Integration
```bash
# .git/hooks/post-commit
#!/bin/bash
echo "🔍 R195: Verifying push..."
if ! verify_commit_pushed; then
    echo "⚠️ WARNING: Commit not pushed!"
    echo "Run: git push"
fi
```

### Continuous Monitoring
```bash
# Background push monitor
monitor_push_status() {
    while true; do
        for effort_dir in "$SF_ROOT/efforts/"*/*/; do
            if [ -d "$effort_dir/.git" ]; then
                cd "$effort_dir"
                local unpushed=$(git log @{u}..HEAD --oneline 2>/dev/null | wc -l)
                if [ "$unpushed" -gt 0 ]; then
                    echo "⚠️ $(basename "$effort_dir"): $unpushed unpushed commits"
                fi
            fi
        done
        sleep 60  # Check every minute
    done
}
```

## Push Status Dashboard

```bash
show_push_dashboard() {
    echo "╔══════════════════════════════════════════════╗"
    echo "║           PUSH STATUS DASHBOARD              ║"
    echo "╠══════════════════════════════════════════════╣"
    
    for effort_dir in "$SF_ROOT/efforts/"*/*/; do
        if [ -d "$effort_dir/.git" ]; then
            cd "$effort_dir"
            local effort=$(basename "$effort_dir")
            local branch=$(git branch --show-current)
            local unpushed=$(git log @{u}..HEAD --oneline 2>/dev/null | wc -l)
            
            if [ "$unpushed" -eq 0 ]; then
                echo "║ ✅ $effort: All pushed"
            else
                echo "║ ⚠️  $effort: $unpushed unpushed commits"
            fi
        fi
    done
    
    echo "╚══════════════════════════════════════════════╝"
    cd "$SF_ROOT"
}
```

## Integration with Other Rules

### R194 Dependency
- R194 ensures tracking exists
- R195 ensures commits reach remote

### R189 Similarity
- R189 for TODO commits
- R195 for code commits
- Both require push verification

## Common Push Failures

### ❌ Authentication Failure
```bash
$ git push
remote: Invalid username or password
fatal: Authentication failed
# Fix: Update credentials or use SSH
```

### ❌ Non-Fast-Forward
```bash
$ git push
! [rejected]        branch -> branch (non-fast-forward)
# Fix: Pull and merge/rebase first
```

### ❌ Protected Branch
```bash
$ git push
remote: error: GH006: Protected branch update failed
# Fix: Create PR or use different branch
```

### ✅ Successful Push
```bash
$ git push
Enumerating objects: 5, done.
Counting objects: 100% (5/5), done.
Writing objects: 100% (3/3), 420 bytes | 420.00 KiB/s, done.
   abc123..def456  phase1/wave1/api-types -> phase1/wave1/api-types
```

## Grading Enforcement

### Push Verification Failures
- Commit not pushed within 5 min: -20%
- Commit not pushed within 15 min: -35%
- State transition with unpushed: -40%
- Lost commits (never pushed): -60%
- Force push losing others' work: -80%

### Recovery Success
- Successful retry after failure: +5%
- Proper conflict resolution: +10%
- Clean push history maintained: +5%

## Agent Requirements

### Orchestrator
```bash
# Before spawning agents
pre_transition_push_check || {
    echo "❌ Cannot spawn agents with unpushed commits"
    exit 1
}
```

### SW Engineer
```bash
# After implementation
git add -A
git commit -m "feat: implement feature"
verify_commit_pushed || exit 1
```

### Code Reviewer
```bash
# After review
git add REVIEW-FEEDBACK.md
git commit -m "docs: add review feedback"
verify_commit_pushed || exit 1
```

---
**Remember:** If it's not pushed, it doesn't exist. Push early, push often!