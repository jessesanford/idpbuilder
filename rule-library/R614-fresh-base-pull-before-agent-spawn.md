# R614 - Fresh Base Pull Before Agent Spawn Protocol

## Rule ID: R614
## Criticality: 🔴🔴🔴 SUPREME LAW
## State Scope: SPAWN_SW_ENGINEERS, SPAWN_CODE_REVIEWERS_EFFORT_REVIEW, all spawn states
## Agents: orchestrator
## Related: R208 (CD Before Spawn), R603 (Sequential Timing), R501 (Cascade), R514 (Infrastructure Creation)

## Overview

Before spawning ANY agent to work on an effort, the orchestrator MUST pull latest changes from the effort's base branch to ensure the agent works on fresh code. This prevents agents from working on stale branches and eliminates the need for downstream rebases.

**CRITICAL**: This rule extends R208 (CD Before Spawn) by adding a MANDATORY pull step. The complete spawn sequence is now:
1. CD to effort directory (R208)
2. **Pull latest from base branch (R614) ← NEW**
3. Verify fresh base (R614)
4. Spawn agent

## Problem Statement

### What Happens Without R614

**Scenario**: Sequential efforts with bug fixes applied after initial infrastructure creation:

```
Day 1: Orchestrator creates infrastructure for effort 2.2.1
       - git clone --single-branch origin wave2.1-integration
       - Creates: efforts/phase2/wave2/effort-2.2.1/

Day 2: SW Engineer implements 2.2.1, completes, submits for review

Day 3: Code Reviewer finds bugs BUG-007, BUG-008 in 2.2.1
       - Bugs fixed in effort-2.2.1 branch
       - Fixes pushed: git push origin effort-2.2.1

Day 4: Orchestrator creates infrastructure for effort 2.2.2 (depends on 2.2.1)
       - R603: Correctly waits for 2.2.1 approval ✅
       - R514: Uses base_branch: effort-2.2.1 ✅
       - Creates: efforts/phase2/wave2/effort-2.2.2/
       - **BUT DIRECTORY ALREADY EXISTS from Day 1!**

Day 5: Orchestrator spawns SW Engineer for 2.2.2
       - R208: CD to effort-2.2.2 directory ✅
       - WITHOUT R614: Does NOT pull latest from effort-2.2.1 ❌
       - SW Engineer works on STALE base (missing Day 3 fixes) ❌

Result:
- 2.2.2 builds on code WITHOUT bug fixes
- Either rediscovers BUG-007, BUG-008 (wasted effort)
- OR requires rebase later to incorporate fixes (technical debt)
- CASCADE PATTERN BROKEN
```

### What Happens With R614

```
Day 1-4: Same as above...

Day 5: Orchestrator spawns SW Engineer for 2.2.2
       - R208: CD to effort-2.2.2 directory ✅
       - R614: Pull latest from effort-2.2.1 branch ✅ ← NEW!
       - SW Engineer works on FRESH base (includes Day 3 fixes) ✅

Result:
- 2.2.2 builds on code WITH bug fixes ✅
- No duplicate bug discovery ✅
- No rebase needed ✅
- CASCADE PATTERN MAINTAINED ✅
```

## Requirements

### 1. MANDATORY Pull Before EVERY Spawn

Before spawning ANY agent (SW-Engineer, Code-Reviewer, etc.), orchestrator MUST:

```bash
# MANDATORY SEQUENCE - NO EXCEPTIONS

# Step 1: CD to effort directory (R208)
cd "$EFFORT_DIR" || {
    echo "🚨 FATAL: Cannot CD to $EFFORT_DIR"
    exit 208
}

# Step 2: Identify base branch (from pre_planned_infrastructure)
BASE_BRANCH=$(jq -r ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".base_branch" \
              orchestrator-state-v3.json)

if [ -z "$BASE_BRANCH" ] || [ "$BASE_BRANCH" = "null" ]; then
    echo "🚨 FATAL: No base_branch found for $EFFORT_ID"
    exit 614
fi

echo "🔄 R614: Pulling latest from base branch: $BASE_BRANCH"

# Step 3: Fetch latest from origin
git fetch origin "$BASE_BRANCH" || {
    echo "⚠️ WARNING: Cannot fetch $BASE_BRANCH from origin"
    echo "This may cause stale branch issues"
    # Decision: Continue with warning or fail hard?
    # Recommendation: FAIL HARD to prevent stale work
    exit 614
}

# Step 4: Pull latest from base branch (CRITICAL!)
git pull origin "$BASE_BRANCH" || {
    echo "🚨 FATAL: Cannot pull from $BASE_BRANCH"
    echo "Possible causes:"
    echo "  - Uncommitted changes in working directory"
    echo "  - Merge conflicts between local and remote"
    echo "  - Network issues"
    exit 614
}

# Step 5: Verify we have latest commits
REMOTE_HEAD=$(git rev-parse origin/$BASE_BRANCH)
LOCAL_HEAD=$(git rev-parse HEAD)
MERGE_BASE=$(git merge-base HEAD origin/$BASE_BRANCH)

if [[ "$MERGE_BASE" != "$REMOTE_HEAD" ]]; then
    echo "⚠️ WARNING: Local branch may be ahead of base branch"
    echo "Remote HEAD: $REMOTE_HEAD"
    echo "Local HEAD:  $LOCAL_HEAD"
    echo "Merge base:  $MERGE_BASE"
    # This is OK if we have local commits, just verify we HAVE the latest from base
fi

# Step 6: Verify clean working directory
if ! git diff-index --quiet HEAD --; then
    echo "🚨 FATAL: Working directory has uncommitted changes"
    echo "Cannot pull latest while changes are uncommitted"
    git status
    exit 614
fi

echo "✅ R614: Successfully pulled latest from $BASE_BRANCH"
echo "   Current HEAD: $(git rev-parse --short HEAD)"
echo "   Remote HEAD:  $(git rev-parse --short origin/$BASE_BRANCH)"

# Step 7: NOW spawn agent (R208)
echo "🚀 Spawning agent in fresh environment..."
# task spawn agent-name ...
```

### 2. When to Pull Latest

**MANDATORY for these spawns:**
- SPAWN_SW_ENGINEERS (before implementation)
- SPAWN_CODE_REVIEWERS_EFFORT_REVIEW (before review - should review fresh code)

**RECOMMENDED for these spawns:**
- SPAWN_CODE_REVIEWERS_EFFORT_PLANNING (planning on fresh base)
- Any spawn that reads code from base branch

**NOT REQUIRED for these spawns:**
- SPAWN_ARCHITECT_MASTER_PLANNING (works at project level, not effort level)
- Spawns that don't interact with effort code

### 3. Handling Pull Failures

**Uncommitted changes:**
```bash
if ! git diff-index --quiet HEAD --; then
    echo "🚨 FATAL: Cannot pull - uncommitted changes detected"
    echo "This should NEVER happen - all work should be committed per R220"
    echo "Transition to ERROR_RECOVERY for investigation"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="R614 violation: uncommitted changes before pull"
    exit 614
fi
```

**Merge conflicts:**
```bash
if git pull origin "$BASE_BRANCH" 2>&1 | grep -q "CONFLICT"; then
    echo "🚨 FATAL: Merge conflict pulling from $BASE_BRANCH"
    echo "This should NEVER happen in sequential cascade pattern"
    echo "Indicates catastrophic failure in cascade integrity"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="R614 violation: merge conflict in cascade"
    exit 614
fi
```

**Network failures:**
```bash
# Retry with exponential backoff
for attempt in 1 2 3; do
    if git fetch origin "$BASE_BRANCH"; then
        break
    else
        echo "⚠️ Fetch attempt $attempt failed, retrying..."
        sleep $((2 ** attempt))
    fi
done

if ! git fetch origin "$BASE_BRANCH"; then
    echo "🚨 FATAL: Cannot fetch after 3 attempts"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="R614 violation: network failure fetching base branch"
    exit 614
fi
```

### 4. First Effort in Cascade (Edge Case)

**Special case**: When effort is FIRST in wave (base_branch = wave-X.Y-integration):

```bash
# First effort may not have wave integration branch yet
if [[ "$BASE_BRANCH" =~ integration ]] && ! git ls-remote --heads origin "$BASE_BRANCH" | grep -q "$BASE_BRANCH"; then
    echo "📝 Note: Integration branch $BASE_BRANCH doesn't exist yet"
    echo "This is expected for first effort in wave"
    echo "Skipping pull (nothing to pull from)"
else
    # Normal case: pull from base branch
    git pull origin "$BASE_BRANCH" || exit 614
fi
```

### 5. Verification After Pull

**Verify fresh base:**
```bash
verify_fresh_base() {
    local base_branch="$1"

    # Get latest commit on remote base branch
    local remote_base_head=$(git rev-parse origin/$base_branch)

    # Get merge-base between our HEAD and remote base
    local merge_base=$(git merge-base HEAD origin/$base_branch)

    if [[ "$merge_base" == "$remote_base_head" ]]; then
        echo "✅ Branch is based on latest commit from $base_branch"
        echo "   Remote base HEAD: $(git rev-parse --short $remote_base_head)"
        return 0
    else
        echo "⚠️ WARNING: Branch may not have latest commits from $base_branch"
        echo "   Remote base HEAD: $(git rev-parse --short $remote_base_head)"
        echo "   Our merge base:   $(git rev-parse --short $merge_base)"

        # Show commits we're missing
        echo "   Missing commits:"
        git log --oneline $merge_base..$remote_base_head | head -5

        return 1
    fi
}

verify_fresh_base "$BASE_BRANCH" || {
    echo "🚨 WARNING: Not on fresh base - continuing but may cause issues"
}
```

## Integration with Existing Rules

### Extends R208 (CD Before Spawn)

**R208 OLD sequence:**
```bash
1. CD to effort directory
2. Spawn agent
```

**R208 + R614 NEW sequence:**
```bash
1. CD to effort directory (R208)
2. Pull latest from base branch (R614) ← NEW
3. Verify fresh base (R614) ← NEW
4. Spawn agent (R208)
```

### Works with R603 (Sequential Timing)

**R603** ensures infrastructure is created AFTER dependencies complete.
**R614** ensures we pull latest FROM those dependencies before spawning.

```
R603: "Create 2.2.2 infrastructure AFTER 2.2.1 approved" ✅
R614: "Pull latest FROM 2.2.1 BEFORE spawning for 2.2.2" ✅
```

Together they maintain CASCADE integrity.

### Works with R514 (Infrastructure Creation)

**R514** creates infrastructure with initial base branch clone.
**R614** updates that clone before each agent spawn.

```
R514: "Clone base branch when creating infrastructure" (Day 1)
R614: "Pull latest from base branch before spawning" (Day 5)
```

### Works with R501 (Progressive Trunk-Based Development)

**R501** defines the CASCADE pattern (each effort builds on previous).
**R614** enforces the pattern by ensuring fresh base commits.

## Enforcement

### Exit Code
- Exit code 614 for violations
- ERROR_RECOVERY transition when pull fails

### Grading Penalties
- **Missing pull before spawn:** -50% (stale work, defeats cascade)
- **Pull failure ignored:** -75% (allows work on incorrect base)
- **Uncommitted changes before pull:** -100% (R220 violation)
- **Merge conflict in cascade:** -100% (CASCADE CORRUPTION)

### Validation Points
- Pre-spawn checklist in SPAWN_SW_ENGINEERS
- Pre-spawn checklist in SPAWN_CODE_REVIEWERS_EFFORT_REVIEW
- Post-spawn verification: agent has latest base commits
- State machine validation: base branch freshness

### Automated Validation

**Validation script**: `tools/validate-fresh-base.sh`

```bash
#!/bin/bash
# Validates that effort is on latest base branch
# Usage: validate-fresh-base.sh <effort_dir> <base_branch>

set -e

EFFORT_DIR="$1"
BASE_BRANCH="$2"

if [ -z "$EFFORT_DIR" ] || [ -z "$BASE_BRANCH" ]; then
    echo "Usage: $0 <effort_dir> <base_branch>"
    exit 1
fi

cd "$EFFORT_DIR" || exit 1

# Fetch latest
git fetch origin "$BASE_BRANCH" || {
    echo "❌ Cannot fetch $BASE_BRANCH"
    exit 1
}

# Check merge-base
LATEST=$(git rev-parse origin/$BASE_BRANCH)
CURRENT_BASE=$(git merge-base HEAD origin/$BASE_BRANCH)

if [[ "$LATEST" == "$CURRENT_BASE" ]]; then
    echo "✅ Effort is on latest base"
    echo "   Base: $BASE_BRANCH"
    echo "   Commit: $(git rev-parse --short $LATEST)"
    exit 0
else
    echo "❌ Effort is behind base branch!"
    echo "   Base: $BASE_BRANCH"
    echo "   Latest: $(git rev-parse --short $LATEST)"
    echo "   Our base: $(git rev-parse --short $CURRENT_BASE)"
    echo "   Behind by: $(git rev-list --count $CURRENT_BASE..$LATEST) commits"
    exit 1
fi
```

## Benefits

### With R614 Compliance

✅ Sequential efforts always work on latest code
✅ Bug fixes automatically incorporated into downstream work
✅ No duplicate bug discovery
✅ No downstream rebases needed
✅ CASCADE pattern maintained throughout project
✅ Clean merge history
✅ Reduced technical debt

### Without R614 (Current Problem)

❌ Stale base branches
❌ Duplicate bug discovery
❌ Rebases required later
❌ CASCADE pattern broken
❌ Technical debt accumulation
❌ Lost productivity

## Examples

### Example 1: Bug Fixes in Sequential Efforts (The User's Case)

**Scenario**: Wave 2.2 with two sequential efforts

**Timeline**:
```
Day 1: Create infrastructure for effort 2.2.1
       git clone origin wave2.1-integration → effort-2.2.1/

Day 2: Spawn SWE for 2.2.1, implement, complete

Day 3: Code review finds BUG-007, BUG-008
       Fix in effort-2.2.1 branch, push fixes
       git push origin effort-2.2.1

Day 4: Create infrastructure for effort 2.2.2
       base_branch: effort-2.2.1 (per R603)

Day 5: SPAWN_SW_ENGINEERS for effort 2.2.2

       WITHOUT R614:
       - CD to effort-2.2.2/
       - Spawn agent
       - Agent works on OLD code (no bug fixes)

       WITH R614:
       - CD to effort-2.2.2/
       - git fetch origin effort-2.2.1
       - git pull origin effort-2.2.1 ← GETS BUG FIXES!
       - Spawn agent
       - Agent works on FRESH code (includes bug fixes)
```

### Example 2: Parallel vs Sequential Behavior

**Parallel efforts** (no dependencies):
```
Effort 3.1.1: base = wave3.0-integration
Effort 3.1.2: base = wave3.0-integration

Both created simultaneously, both pull from same base.
R614 ensures both have latest integration commits.
```

**Sequential efforts** (with dependencies):
```
Effort 3.2.1: base = wave3.1-integration
Effort 3.2.2: base = effort-3.2.1 ← DEPENDS ON 3.2.1!

3.2.1 created first, completes, gets fixes.
3.2.2 created second.
R614 ensures 3.2.2 pulls latest from 3.2.1 (including fixes).
```

### Example 3: Integration Dependencies

**Effort depends on integration branch**:
```bash
# Effort 2.1.1 depends on integration:wave2.0-integration
BASE_BRANCH="wave2.0-integration"

# R614 protocol:
cd efforts/phase2/wave1/effort-2.1.1/
git fetch origin wave2.0-integration
git pull origin wave2.0-integration

# Ensures effort has all changes from previous wave's integration
# Critical for wave-to-wave CASCADE
```

## Anti-Patterns

### ❌ WRONG: Spawning Without Pull

```bash
# This violates R614
cd "$EFFORT_DIR"
task spawn sw-engineer "implement effort"
# Agent works on stale base!
```

### ❌ WRONG: Pulling After Spawn

```bash
# This violates R614
cd "$EFFORT_DIR"
task spawn sw-engineer "cd . && git pull origin $BASE_BRANCH && implement"
# Pull happens INSIDE agent, not before spawn
# Preflight checks run on stale code
```

### ❌ WRONG: Assuming Fresh Base

```bash
# This violates R614
cd "$EFFORT_DIR"
# No pull!
task spawn sw-engineer "implement effort"
# Assumption: "Infrastructure was just created, must be fresh"
# Reality: Infrastructure may be days old with many upstream fixes!
```

### ✅ RIGHT: Full R614 Protocol

```bash
# This complies with R614
cd "$EFFORT_DIR" || exit 208
BASE_BRANCH=$(jq -r ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".base_branch" orchestrator-state-v3.json)
git fetch origin "$BASE_BRANCH" || exit 614
git pull origin "$BASE_BRANCH" || exit 614
verify_fresh_base "$BASE_BRANCH" || echo "Warning: stale base"
task spawn sw-engineer "implement effort"
```

## State Integration

### SPAWN_SW_ENGINEERS State

**Add to BLOCKING CHECKLIST**:
```markdown
## BLOCKING REQUIREMENTS (R510)

Before spawning ANY SW Engineer:

[X] Checklist Item N: Pull latest from base branch per R614
    - CD to effort directory (R208)
    - Identify base_branch from pre_planned_infrastructure
    - git fetch origin <base_branch>
    - git pull origin <base_branch>
    - Verify working directory clean
    - Verify fresh base (merge-base check)
    - Proof: Git log shows latest commits from base branch
```

### SPAWN_CODE_REVIEWERS_EFFORT_REVIEW State

**Add to STANDARD CHECKLIST**:
```markdown
## STANDARD REQUIREMENTS (R510)

Recommended before spawning Code Reviewer:

[ ] Checklist Item N: Pull latest from base branch per R614
    - Ensures reviewer sees most current code
    - Includes any fixes applied since implementation
    - Recommended but not blocking (reviewer can work on older code)
```

## Monitoring and Metrics

### Freshness Metrics

Track base branch freshness:
```json
{
  "effort_metadata": {
    "effort-2.2.2": {
      "base_branch_freshness": {
        "base_branch": "effort-2.2.1",
        "last_pull_at": "2025-01-20T14:30:00Z",
        "last_pull_commit": "abc123def456",
        "remote_head_at_pull": "abc123def456",
        "commits_behind_at_spawn": 0,
        "freshness_status": "FRESH"
      }
    }
  }
}
```

### Alert Conditions

- If commits_behind_at_spawn > 0: WARNING (stale base)
- If last_pull_at > 6 hours ago: WARNING (old pull)
- If pull fails: CRITICAL (cannot spawn)

## Migration Path

**Phase 1**: Add R614 to rule library ✅ (this file)
**Phase 2**: Update SPAWN_SW_ENGINEERS state rules with R614 checklist
**Phase 3**: Update SPAWN_CODE_REVIEWERS_EFFORT_REVIEW state rules
**Phase 4**: Update orchestrator.md agent config with R614
**Phase 5**: Add R614 to RULE-REGISTRY.md
**Phase 6**: Create validation script tools/validate-fresh-base.sh
**Phase 7**: Test with sequential efforts workflow
**Phase 8**: Deploy to all projects via upgrade.sh

## Summary

**R614 is SUPREME LAW #52**: Fresh Base Pull Before Agent Spawn

**Principle**: ALWAYS pull latest from base branch before spawning agents.

**Integration**: Extends R208 (CD) with mandatory pull step.

**Purpose**: Maintains CASCADE integrity by ensuring fresh base branches.

**Impact**: Eliminates stale work, prevents duplicate bugs, removes rebase need.

**Enforcement**: Exit code 614, -50% to -100% penalties, ERROR_RECOVERY on failure.

**MANDATORY in**: SPAWN_SW_ENGINEERS, SPAWN_CODE_REVIEWERS_EFFORT_REVIEW

**Result**: Clean cascade, fresh code, efficient workflow, no technical debt.
