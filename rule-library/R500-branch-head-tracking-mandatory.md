# 🚨🚨🚨 RULE R500 - Branch HEAD Tracking Synchronization [BLOCKING]

**FAIL = STATE CORRUPTION** | Source: [R500](rule-library/R500-branch-head-tracking-mandatory.md)

## Purpose
Prevent branch tracking issues, incorrect base branches, and stale HEAD references that led to catastrophic branch deletion incidents.

## Requirements

### 1. MANDATORY HEAD SYNCHRONIZATION
**EVERY STATE TRANSITION MUST:**
- Update `branch_head_tracking.branches` for ALL active branches
- Record current HEAD commit for each branch
- Update `branch_head_tracking.last_sync` timestamp
- Increment `branch_head_tracking.sync_count`

### 2. SYNCHRONIZATION TRIGGERS (ALL MANDATORY)
```bash
# MUST sync when:
1. Before ANY state transition
2. After creating ANY new branch
3. After ANY commit operation
4. After ANY merge operation
5. Before spawning ANY agent
6. After completing ANY effort/wave/phase
7. Every 30 minutes during long operations
8. Before and after ANY integration
```

### 3. IMPLEMENTATION
```bash
# Function EVERY agent MUST call
sync_branch_heads() {
    local state_file="$PROJECT_ROOT/orchestrator-state-v3.json"
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

    # Get all branches
    for branch in $(git branch -r | grep -v HEAD | sed 's/.*origin\///'); do
        local head_commit=$(git rev-parse --short "origin/$branch" 2>/dev/null || continue)

        # Update in orchestrator-state-v3.json
        jq --arg branch "$branch" \
           --arg head "$head_commit" \
           --arg ts "$timestamp" \
           '.branch_head_tracking.branches[$branch].head_commit = $head |
            .branch_head_tracking.branches[$branch].last_updated = $ts' \
           "$state_file" > tmp.json && mv tmp.json "$state_file"
    done

    # Update sync metadata
    jq --arg ts "$timestamp" \
       '.branch_head_tracking.last_sync = $ts |
        .branch_head_tracking.sync_count += 1' \
       "$state_file" > tmp.json && mv tmp.json "$state_file"
}

# MUST be called:
sync_branch_heads || exit 1  # Exit if sync fails
```

### 4. VALIDATION RULES

#### Base Branch Validation
- **First effort in wave**: MUST be based on main or previous integration
- **Subsequent efforts**: MUST be based on previous effort (R308 cascade)
- **First split**: MUST be based on oversized effort branch
- **Subsequent splits**: MUST be based on previous split
- **Wave integration**: MUST be based on first effort of WAVE (R009 Sequential Rebuild)
- **Phase integration**: MUST be based on first effort of PHASE (R282 Sequential Rebuild)
- **Project integration**: MUST be based on main (R283 Sequential Rebuild)

#### Validation Enforcement
```bash
# Pre-commit hook BLOCKS commits when:
- Branch is based on wrong parent
- Branch missing files from expected base
- HEAD tracking is stale (>30 minutes)
- Base branch cannot be determined
```

### 5. RECOVERY PROCEDURES

#### When Branch Has Wrong Base
1. **STOP** - Do not attempt to fix with rebase
2. **NOTIFY** - Alert orchestrator immediately
3. **RECREATE** - Branch must be recreated from correct base
4. **VERIFY** - Run validation before any commits

#### When HEAD Tracking Is Stale
```bash
# Emergency sync procedure
emergency_sync() {
    echo "🚨 EMERGENCY: Synchronizing all branch HEADs..."
    sync_branch_heads
    git add orchestrator-state-v3.json
    git commit -m "emergency: sync branch HEAD tracking [R500]"
    git push
}
```

### 6. GRADING IMPACT

**VIOLATIONS RESULT IN:**
- Missing sync at state transition: **-20% grade**
- Stale HEAD tracking (>30 min): **-15% grade**
- Wrong base branch committed: **-30% grade**
- Branch corruption from wrong base: **-50% grade**
- Loss of branch tracking data: **-100% IMMEDIATE FAILURE**

### 7. MONITORING_SWE_PROGRESS

#### Orchestrator Must Track
```json
{
  "branch_head_tracking": {
    "last_sync": "ISO timestamp",
    "sync_count": "increments each sync",
    "branches": {
      "branch-name": {
        "head_commit": "current HEAD",
        "last_updated": "ISO timestamp",
        "base_branch": "parent branch",
        "validation_status": "valid|invalid_base|needs_rebase"
      }
    }
  }
}
```

#### Compliance Check
```bash
# Check sync freshness
check_sync_compliance() {
    local last_sync=$(jq -r '.branch_head_tracking.last_sync' orchestrator-state-v3.json)
    local now=$(date +%s)
    local sync_time=$(date -d "$last_sync" +%s)
    local age=$((now - sync_time))

    if [ $age -gt 1800 ]; then  # >30 minutes
        echo "🚨 R500 VIOLATION: HEAD tracking is $((age/60)) minutes old!"
        return 1
    fi
}
```

## Related Rules
- R287: TODO persistence and state management
- R360: Just-in-time infrastructure creation
- R206: State machine validation
- R322: Mandatory stop conditions

## Implementation Files
- `/tools/hooks/pre-commit-base-validation` - Pre-commit validation hook
- `/tools/test-base-validation-hook.sh` - Testing script
- `orchestrator-state.schema.json` - Schema with HEAD tracking
- `upgrade.sh` - Hook installation in all branches

## Critical Notes
- **NEVER** attempt to fix base branch issues with rebase/cherry-pick
- **ALWAYS** sync before state transitions
- **IMMEDIATELY** stop work if validation fails
- **TRACK** every branch modification in orchestrator-state-v3.json

---

**Remember**: This rule exists because we lost branches due to incorrect base tracking. The cost of not following this rule is catastrophic data loss.