# 🚨🚨🚨 RULE R194 - Remote Branch Tracking

**Criticality:** BLOCKING - No tracking = Lost work  
**Grading Impact:** -40% for untracked branches  
**Enforcement:** CONTINUOUS - Every git operation

## Rule Statement

EVERY effort branch MUST track its remote counterpart. NO local-only branches. ALL commits MUST be pushable to origin.

## Remote Tracking Requirements

### 1. Branch Creation MUST Include Push
```bash
# ❌ WRONG - Local branch only
git checkout -b phase1/wave1/api-types

# ✅ CORRECT - With remote tracking
git checkout -b phase1/wave1/api-types
git push -u origin phase1/wave1/api-types
```

### 2. Verify Tracking Configuration
```bash
verify_branch_tracking() {
    local branch=$(git branch --show-current)
    
    # Check if branch has upstream
    local upstream=$(git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null)
    
    if [ -z "$upstream" ]; then
        echo "❌ VIOLATION R194: Branch '$branch' has no remote tracking!"
        echo "Fix with: git push -u origin $branch"
        return 1
    fi
    
    echo "✅ Branch '$branch' tracks '$upstream'"
    
    # Check if we're behind/ahead
    local ahead=$(git rev-list --count @{u}..HEAD)
    local behind=$(git rev-list --count HEAD..@{u})
    
    if [ "$ahead" -gt 0 ]; then
        echo "⚠️ Branch is $ahead commits ahead of $upstream"
        echo "Push with: git push"
    fi
    
    if [ "$behind" -gt 0 ]; then
        echo "⚠️ Branch is $behind commits behind $upstream"
        echo "Pull with: git pull"
    fi
    
    return 0
}
```

### 3. Enforce Tracking Before Work
```bash
# SW Engineer pre-flight check
enforce_tracking_before_work() {
    echo "🔍 Checking branch tracking..."
    
    if ! verify_branch_tracking; then
        echo "❌ Cannot proceed without remote tracking"
        
        # Attempt to set up tracking
        local branch=$(git branch --show-current)
        echo "🔧 Attempting to set up tracking..."
        
        if git push -u origin "$branch" --force-with-lease; then
            echo "✅ Remote tracking established"
        else
            echo "❌ Failed to establish tracking"
            exit 1
        fi
    fi
}
```

## Push Frequency Requirements

### After Every Logical Unit
```bash
# MUST push after:
- Completing a function/method
- Adding a test
- Fixing a bug
- Creating a new file
- Major refactoring

# Example workflow
git add pkg/api/types.go
git commit -m "feat: add workspace type definition"
git push  # IMMEDIATE PUSH
```

### Maximum Unpushed Time
```bash
check_push_staleness() {
    local last_push=$(git log -1 --format=%ct origin/$(git branch --show-current) 2>/dev/null)
    local last_commit=$(git log -1 --format=%ct HEAD)
    
    if [ -z "$last_push" ]; then
        echo "❌ Branch never pushed to remote!"
        return 1
    fi
    
    local unpushed_time=$((last_commit - last_push))
    
    if [ "$unpushed_time" -gt 300 ]; then  # 5 minutes
        echo "⚠️ WARNING: Commits unpushed for $((unpushed_time/60)) minutes"
        echo "Push immediately with: git push"
    fi
}
```

## Remote Synchronization

### Pull Before Push
```bash
safe_push() {
    echo "📥 Pulling latest changes..."
    git pull --rebase origin $(git branch --show-current)
    
    echo "📤 Pushing changes..."
    if ! git push; then
        echo "⚠️ Push failed, attempting with lease..."
        git push --force-with-lease
    fi
}
```

### Handle Diverged Branches
```bash
handle_diverged_branch() {
    local branch=$(git branch --show-current)
    
    echo "⚠️ Branch has diverged from remote"
    echo "Options:"
    echo "1. Rebase: git pull --rebase origin $branch"
    echo "2. Merge: git pull origin $branch"
    echo "3. Force: git push --force-with-lease origin $branch"
    
    # Prefer rebase for clean history
    git pull --rebase origin "$branch"
    git push
}
```

## Branch Status Monitoring

### Regular Status Checks
```bash
monitor_branch_status() {
    while true; do
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "Branch Status Check"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        
        # Current branch
        echo "Branch: $(git branch --show-current)"
        
        # Tracking status
        git status -sb
        
        # Unpushed commits
        local unpushed=$(git log origin/$(git branch --show-current)..HEAD --oneline 2>/dev/null | wc -l)
        if [ "$unpushed" -gt 0 ]; then
            echo "⚠️ $unpushed unpushed commits"
        fi
        
        # Time since last push
        check_push_staleness
        
        sleep 300  # Check every 5 minutes
    done
}
```

## Integration with CI/CD

### Branch Protection Rules
```yaml
# GitHub branch protection for effort branches
protection_rules:
  pattern: "phase*/wave*/*"
  required_checks:
    - "build"
    - "test"
    - "line-count-check"
  dismiss_stale_reviews: true
  require_up_to_date: true
```

### Automated Push Validation
```bash
validate_push() {
    local branch="$1"
    
    # Check branch exists on remote
    if ! git ls-remote --heads origin "$branch" | grep -q "$branch"; then
        echo "❌ Branch doesn't exist on remote"
        return 1
    fi
    
    # Check latest commit is pushed
    local local_head=$(git rev-parse HEAD)
    local remote_head=$(git rev-parse origin/"$branch")
    
    if [ "$local_head" != "$remote_head" ]; then
        echo "❌ Local and remote heads don't match"
        echo "Local:  $local_head"
        echo "Remote: $remote_head"
        return 1
    fi
    
    echo "✅ Branch fully synchronized with remote"
}
```

## Common Tracking Failures

### ❌ No Upstream Set
```bash
$ git status
On branch phase1/wave1/api-types
nothing to commit, working tree clean

$ git push
fatal: The current branch phase1/wave1/api-types has no upstream branch
```

### ❌ Detached HEAD
```bash
$ git status
HEAD detached at abc123
# Can't track from detached HEAD!
```

### ❌ Wrong Remote
```bash
$ git branch -vv
* phase1/wave1/api-types abc123 [fork/phase1/wave1/api-types]
# Should track origin, not fork!
```

### ✅ Correct Tracking
```bash
$ git branch -vv
* phase1/wave1/api-types abc123 [origin/phase1/wave1/api-types] feat: add types
```

## Grading Enforcement

### Tracking Violations
- No remote tracking: -40%
- Wrong remote (not origin): -30%
- Unpushed >15 minutes: -20%
- Unpushed >30 minutes: -40%
- Lost commits (never pushed): -60%

### Recovery Credit
- Setting up tracking after violation: +10%
- Successfully pushing after delay: +5%
- Resolving conflicts properly: +5%

## Agent Integration

### Orchestrator Verification
```bash
# Before spawning agents
verify_all_effort_tracking() {
    for effort_dir in "$SF_ROOT/efforts/phase${phase}/wave${wave}"/*/; do
        if [ -d "$effort_dir/.git" ]; then
            cd "$effort_dir"
            if ! verify_branch_tracking; then
                echo "❌ Effort $(basename "$effort_dir") lacks tracking"
                return 1
            fi
        fi
    done
}
```

### SW Engineer Compliance
```bash
# In SW engineer startup
echo "🔍 Verifying remote tracking..."
enforce_tracking_before_work || exit 1

# After every commit
git commit -m "feat: implement feature"
git push || echo "⚠️ VIOLATION R194: Failed to push!"
```

---
**Remember:** Every branch needs a home on origin. No orphan branches!