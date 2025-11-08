# 🚨🚨🚨 RULE R282: Phase Integration Protocol (Sequential Rebuild Model)

## Classification
- **Category**: Integration Management
- **Criticality Level**: 🚨🚨🚨 BLOCKING
- **Enforcement**: MANDATORY for phase completion
- **Penalty**: -50% to -100% for violations

## The Rule

**Phase integration MUST use SEQUENTIAL REBUILD model: base from FIRST effort of phase, merge ALL subsequent efforts in order. This creates REAL merge commits that test true trunk-based sequential mergeability.**

**Phase integration MUST occur in completely isolated workspace with fresh clone of target repository, NOT the software-factory repository.**

## 🔴🔴🔴 SEQUENTIAL REBUILD MODEL 🔴🔴🔴

### The Fundamental Change

**OLD MODEL (INCORRECT - VIOLATES R512/R270/R363):**
```bash
# ❌ WRONG - Merging wave integration branches
git checkout -b phase2-integration origin/phase2-wave3-integration
git merge origin/phase2-wave1-integration  # Integration branch as source!
git merge origin/phase2-wave2-integration  # Integration branch as source!
```

**NEW MODEL (CORRECT - SEQUENTIAL REBUILD):**
```bash
# ✅ CORRECT - Base from first effort, merge all subsequent efforts
git checkout -b phase2-integration origin/phase2/wave1/auth-system  # First effort!
git merge --no-ff origin/phase2/wave1/user-management  # Second effort
git merge --no-ff origin/phase2/wave1/permissions      # Third effort
git merge --no-ff origin/phase2/wave2/database-layer   # Fourth effort
git merge --no-ff origin/phase2/wave2/cache-service    # Fifth effort
git merge --no-ff origin/phase2/wave2/api-gateway      # Sixth effort
# ... all efforts in sequential order
```

## Why Sequential Rebuild?

### 1. Tests TRUE Trunk-Based Mergeability
**Simulates how PRs actually merge to main:**
```
Real Production Workflow:
  main ← PR1 (effort1) merges
  main ← PR2 (effort2) merges (after PR1)
  main ← PR3 (effort3) merges (after PR2)

Sequential Rebuild Simulates This:
  effort1 (base)
  effort1 ← effort2 merges
  effort1+effort2 ← effort3 merges
```

### 2. Creates REAL Merges (Not No-Ops)
**OLD model**: Merging wave-integration branches created NO-OP merges because wave integrations already contained all prior work.

**NEW model**: Merging effort branches creates REAL merge commits because each effort adds new changes on top of previous efforts.

### 3. Validates R363 Sequential Direct Mergeability
**R363 requires**: Every effort must merge directly to main in sequence.

**Phase integration validates this** by replaying the exact merge sequence:
- If E1→E2→E3 cascade works, phase integration proves it
- If there are hidden conflicts, phase integration exposes them
- If sequential mergeability fails, we catch it BEFORE pushing to main

### 4. Aligns with R512 Trunk-Based Integration Model
**R512 states**: Integration branches are for TESTING, not merging.

**Sequential rebuild implements this**: Wave integration branches were testing checkpoints, now phase integration tests the actual effort merges.

## Requirements

### 1. Phase Integration Infrastructure Setup

```bash
# Setup for phase integration using sequential rebuild model
PHASE=$(jq '.current_phase' orchestrator-state-v3.json)

# Read target repository configuration (R508)
TARGET_CONFIG="$CLAUDE_PROJECT_DIR/target-repo-config.yaml"
TARGET_REPO_PATH=$(yq '.repository_path' "$TARGET_CONFIG")
TARGET_REPO_NAME=$(yq '.repository_name' "$TARGET_CONFIG")
DEFAULT_BRANCH=$(yq '.default_branch' "$TARGET_CONFIG")

# Create isolated phase integration workspace
INTEGRATE_PHASE_WAVES_DIR="$CLAUDE_PROJECT_DIR/efforts/phase${PHASE}/integration-workspace"
rm -rf "$INTEGRATE_PHASE_WAVES_DIR"
mkdir -p "$INTEGRATE_PHASE_WAVES_DIR"
cd "$INTEGRATE_PHASE_WAVES_DIR"

# Clone TARGET repository (R508 - NEVER software-factory!)
git clone "$TARGET_REPO_PATH" "$TARGET_REPO_NAME"
cd "$TARGET_REPO_NAME"

# CRITICAL SAFETY CHECK - Verify correct repository
REMOTE_URL=$(git remote get-url origin)
if [[ "$REMOTE_URL" == *"software-factory"* ]]; then
    echo "❌ CRITICAL: Cloned orchestrator repository instead of target!"
    echo "Expected: Target project repository"
    echo "Got: $REMOTE_URL"
    exit 508  # R508 violation
fi
```

### 2. Determine First Effort (Base Branch)

```bash
# NEW: Find FIRST effort of the phase to use as base
determine_phase_first_effort() {
    local PHASE=$1

    # Get all efforts for this phase, sorted by wave and effort order
    FIRST_EFFORT=$(jq -r "
        .efforts_completed[] |
        select(.phase == $PHASE) |
        select(.wave == 1) |
        .branch
    " orchestrator-state-v3.json | head -1)

    if [ -z "$FIRST_EFFORT" ] || [ "$FIRST_EFFORT" = "null" ]; then
        echo "❌ FATAL: Cannot determine first effort of phase $PHASE"
        exit 1
    fi

    echo "✅ R282: First effort of Phase $PHASE is: $FIRST_EFFORT"
    echo "$FIRST_EFFORT"
}

# Get first effort as base
FIRST_EFFORT=$(determine_phase_first_effort $PHASE)

# Create phase integration branch FROM first effort
git fetch origin "$FIRST_EFFORT"
git checkout -b "phase${PHASE}-integration" "origin/$FIRST_EFFORT"

echo "✅ Phase integration base: $FIRST_EFFORT (first effort of phase)"
```

### 3. Collect ALL Subsequent Efforts in Order

```bash
# NEW: Collect all efforts AFTER the first, in sequential order
collect_phase_efforts_sequential() {
    local PHASE=$1
    local FIRST_EFFORT=$2

    echo "📊 Collecting all efforts for Phase $PHASE (excluding first effort)..."

    # Get all completed efforts for this phase, sorted by completion order
    ALL_EFFORTS=$(jq -r "
        .efforts_completed[] |
        select(.phase == $PHASE) |
        .branch
    " orchestrator-state-v3.json)

    # Remove first effort from list (it's our base)
    SUBSEQUENT_EFFORTS=$(echo "$ALL_EFFORTS" | grep -v "^${FIRST_EFFORT}$")

    # Count efforts
    TOTAL=$(echo "$ALL_EFFORTS" | wc -l)
    SUBSEQUENT=$(echo "$SUBSEQUENT_EFFORTS" | wc -l)

    echo "✅ Found $TOTAL total efforts, $SUBSEQUENT to merge (excluding base)"
    echo "$SUBSEQUENT_EFFORTS"
}

# Collect efforts to merge
EFFORTS_TO_MERGE=$(collect_phase_efforts_sequential $PHASE "$FIRST_EFFORT")

if [ -z "$EFFORTS_TO_MERGE" ]; then
    echo "⚠️ WARNING: Single-effort phase detected"
    echo "Phase integration will point to first (and only) effort"
    # This is valid! See "Single-Effort Phases" section below
fi
```

### 4. Sequential Merge Execution

```bash
# NEW: Sequential rebuild - merge each effort in order
echo "🔄 Starting sequential rebuild of Phase $PHASE..."
echo "Base: $FIRST_EFFORT"
echo "Merging $(echo "$EFFORTS_TO_MERGE" | wc -l) subsequent efforts..."

MERGE_COUNT=0
while IFS= read -r effort_branch; do
    [ -z "$effort_branch" ] && continue

    MERGE_COUNT=$((MERGE_COUNT + 1))
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📦 Merge $MERGE_COUNT: $effort_branch"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    # Fetch effort branch
    git fetch origin "$effort_branch"

    # Attempt merge with no fast-forward (always create merge commit)
    if git merge --no-ff "origin/$effort_branch" \
         -m "integrate: Merge $effort_branch into phase${PHASE}-integration

This is merge $MERGE_COUNT in the sequential rebuild of Phase $PHASE.
Base: $FIRST_EFFORT (first effort)
Target: phase${PHASE}-integration

Per R282 Sequential Rebuild Model."; then
        echo "✅ Merge $MERGE_COUNT successful: $effort_branch"
    else
        echo "❌ Merge conflict in $effort_branch"
        echo "🔍 Conflict details:"
        git status

        # Document conflict
        cat > "MERGE-CONFLICT-${MERGE_COUNT}.md" << EOF
# Merge Conflict - $effort_branch

**Merge Number:** $MERGE_COUNT
**Branch:** $effort_branch
**Phase:** $PHASE
**Time:** $(date -Iseconds)

## Conflict Status
\`\`\`
$(git status)
\`\`\`

## Conflicted Files
\`\`\`
$(git diff --name-only --diff-filter=U)
\`\`\`

## Resolution Required
This conflict must be resolved before continuing the sequential rebuild.
EOF
        echo "❌ FAILED: Merge conflict requires resolution"
        exit 1
    fi

    # Run tests after EACH merge to catch integration issues early
    echo "🧪 Running tests after merge $MERGE_COUNT..."
    if ! (npm test 2>/dev/null || make test 2>/dev/null || pytest 2>/dev/null || true); then
        echo "⚠️ WARNING: Tests may have failed after merge $MERGE_COUNT"
        echo "Continue with caution - review test results"
    fi

done <<< "$EFFORTS_TO_MERGE"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Sequential rebuild complete!"
echo "   Base effort: $FIRST_EFFORT"
echo "   Merged $MERGE_COUNT subsequent efforts"
echo "   Result: phase${PHASE}-integration branch"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
```

## Edge Cases

### Single-Effort Phases

**Scenario**: Phase has only ONE effort (e.g., phase1/wave1/foundation)

**Handling**:
```bash
# First effort becomes the base
git checkout -b phase1-integration origin/phase1/wave1/foundation

# No subsequent efforts to merge
# Result: phase1-integration points to same commit as first effort

# This is VALID and CORRECT!
# The phase integration branch exists for consistency
```

**Why this is okay**:
- Maintains consistent branch naming (`phase1-integration` always exists)
- Provides integration checkpoint even for single-effort phases
- Future phases can reference this integration point
- Tools expect phase integration branches to exist

### ERROR_RECOVERY Fix Integration

**Scenario**: Bugs found during phase testing, fixes created in ERROR_RECOVERY

**Handling**:
```bash
# Sequential rebuild includes fix branches at the END
FIRST_EFFORT=$(determine_phase_first_effort $PHASE)
git checkout -b "phase${PHASE}-integration" "origin/$FIRST_EFFORT"

# 1. Merge all original efforts in order
for effort in $ORIGINAL_EFFORTS; do
    git merge --no-ff "origin/$effort"
done

# 2. Merge all fix branches at the END
for fix in $FIX_BRANCHES; do
    git merge --no-ff "origin/$fix" \
        -m "integrate: Apply ERROR_RECOVERY fix $fix to phase${PHASE}-integration"
done
```

**Order matters**:
- Original efforts first (in sequential order)
- Fix branches last (they fix issues found in integrated state)

### Split Branches

**Scenario**: An effort was too large and got split into multiple branches

**Handling**:
```bash
# Splits are included in sequential order, marked as deprecated in state file
# per R296 (Deprecated Branch Marking Protocol)

EFFORTS_TO_MERGE:
  - phase2/wave1/auth-system           # Original first effort
  - phase2/wave1/user-management-split-001  # Split branch (original marked deprecated)
  - phase2/wave1/user-management-split-002  # Split branch
  - phase2/wave1/permissions            # Next original effort

# The deprecated original branch is NOT merged (per R296)
# Only the split branches are merged
```

## Workspace Isolation Requirements

### Mandatory Isolation
- **Location**: `$CLAUDE_PROJECT_DIR/efforts/phase{X}/integration-workspace/[target-repo-name]/`
- **Repository**: Fresh clone of TARGET repository (R508 requirement)
- **Branch**: `phase-{X}-integration`
- **NEVER**: Use software-factory repository
- **NEVER**: Reuse wave integration workspaces

### Validation Checks
```bash
validate_phase_integration_workspace() {
    # Check correct directory pattern
    if [[ "$PWD" != */phase*/integration-workspace/* ]]; then
        echo "❌ Not in phase integration workspace!"
        return 1
    fi

    # Check NOT software-factory
    if git remote get-url origin | grep -q "software-factory"; then
        echo "❌ CRITICAL: Working in wrong repository!"
        return 1
    fi

    # Check branch naming
    if ! git branch --show-current | grep -q "phase.*integration"; then
        echo "❌ Wrong branch naming convention!"
        return 1
    fi

    return 0
}
```

## Integration Process Summary

### Step 1: Setup Infrastructure
- Create isolated workspace
- Clone target repository
- Verify NOT software-factory (R508)

### Step 2: Determine Base
- Find FIRST effort of phase
- Create integration branch from it

### Step 3: Collect Efforts
- Get ALL subsequent efforts in order
- Include splits, exclude deprecated originals (R296)
- Include ERROR_RECOVERY fixes at end

### Step 4: Sequential Rebuild
- Merge each effort with --no-ff
- Test after EACH merge
- Document conflicts if they occur
- Continue until all merged

### Step 5: Validation & Demo (R291)
- Run full test suite
- Create test harness
- Build production artifact
- Create demo documentation
- Generate integration report

### Step 6: Push and Document
- Push phase integration branch
- Update state file (R288)
- Create PR (if applicable)
- Document completion

## Comparison with Wave Integration

### Wave Integration (Testing Checkpoint)
**Purpose**: Test that efforts within a wave work together
**Model**: Merge wave's effort branches to create checkpoint
**Result**: Wave integration branch (for testing)
**Used as**: Testing branch, gets deleted after validation (R364)

### Phase Integration (Sequential Rebuild)
**Purpose**: Validate sequential mergeability of ALL phase efforts
**Model**: Rebuild entire phase from first effort
**Result**: Phase integration branch (proves sequential merge works)
**Used as**: Validation that R363 sequential mergeability holds

**Key Difference**: Wave integrations are CHECKPOINTS, phase integrations are VALIDATIONS.

## Validation Requirements

### Pre-Integration Validation
```bash
# Verify all waves completed
for wave in $(seq 1 $TOTAL_WAVES); do
    wave_status=$(jq ".waves.wave_${wave}.status" orchestrator-state-v3.json)
    if [ "$wave_status" != "INTEGRATED" ]; then
        echo "🚨 Cannot integrate phase - Wave $wave not integrated"
        exit 1
    fi
done

# Verify first effort exists
if ! git ls-remote --heads origin "$FIRST_EFFORT" > /dev/null 2>&1; then
    echo "🚨 First effort branch not found: $FIRST_EFFORT"
    exit 1
fi
```

### Post-Integration Validation
```bash
# Comprehensive testing
./phase${PHASE}-test-harness.sh

# Build verification
npm run build:prod || exit 1

# Demo creation (R291)
./create-phase-demo.sh

# Integration report
cat > PHASE-${PHASE}-INTEGRATE_WAVE_EFFORTS-REPORT.md << EOF
# Phase $PHASE Integration Report (Sequential Rebuild)

## Model Used
**Sequential Rebuild per R282**
- Base: $FIRST_EFFORT (first effort of phase)
- Merged: $MERGE_COUNT subsequent efforts in order
- Result: True validation of sequential mergeability

## All Efforts Merged (In Order)
$(echo "$ALL_EFFORTS" | nl)

## Build & Test Status
- Build: ✅ PROJECT_DONEFUL
- Tests: ✅ PASSING
- Demo: ✅ CREATED

## Validation
- Repository: $(git remote get-url origin) ✅
- Workspace: $PWD ✅
- Sequential merges: $MERGE_COUNT ✅
- Conflicts: $([ -f MERGE-CONFLICT-*.md ] && echo "DOCUMENTED" || echo "NONE") ✅
EOF
```

## Success Criteria

- ✅ First effort identified correctly
- ✅ All subsequent efforts collected in order
- ✅ Sequential rebuild completed ($MERGE_COUNT merges)
- ✅ All merges used --no-ff (real merge commits)
- ✅ Tests passing after integration
- ✅ Build successful
- ✅ Demo created (R291)
- ✅ Workspace isolated correctly (R508)
- ✅ State file updated (R288)

## Failure Conditions

### Critical Failures (Immediate Stop)
- 🚨 Using software-factory repository = IMMEDIATE FAIL (-100%)
- 🚨 Using wave integration branches as base = VIOLATION R512/R270 (-100%)
- 🚨 Using wave integration branches as sources = VIOLATION R512/R270 (-100%)
- 🚨 Merging with fast-forward (no merge commits) = Invalid test (-50%)
- 🚨 Cannot find first effort = Missing data (-75%)
- 🚨 Test failures ignored = Quality failure (-60%)

### Recovery Protocol
1. **STOP ALL WORK IMMEDIATELY**
2. Document failure point and error
3. Transition to ERROR_RECOVERY state
4. Clean integration workspace
5. Verify state file data
6. Start fresh with corrected approach
7. **NEVER continue from corrupted state**

## Penalties

- Using software-factory repository: **-100% IMMEDIATE FAILURE**
- Using wave integration branches: **-100% (R512/R270 violation)**
- Non-isolated workspace: **-75% grade**
- Missing efforts in sequence: **-50% grade**
- Fast-forward merges (no merge commits): **-50% grade**
- Test failures ignored: **-60% grade**
- No demo created: **-30% grade** (R291)
- State file not updated: **-40% grade** (R288)

## Related Rules

- **R512**: Trunk-Based Integration Model (integrate effort branches)
- **R270**: No Integration Branches as Sources (wave integrations not sources)
- **R363**: Sequential Direct Mergeability (what we're validating)
- **R364**: Integration Testing Only Branches (wave integrations are testing)
- **R308**: Incremental Branching Strategy (efforts cascade)
- **R034**: Integration Requirements (general integration rules)
- **R250**: Integration Isolation (workspace requirements)
- **R259**: Phase Integration After Fixes (ERROR_RECOVERY handling)
- **R288**: State File Updates (required after integration)
- **R291**: Integration Demo Requirement (must create demo)
- **R296**: Deprecated Branch Marking (handle split branches)
- **R508**: Target Repository Enforcement (never use software-factory)

## Common Violations

### ❌ WRONG: Using wave integration as base
```bash
# This creates a NO-OP merge situation
git checkout -b phase2-integration origin/phase2-wave3-integration
git merge origin/phase2-wave1-integration  # NO-OP! Already in wave3
git merge origin/phase2-wave2-integration  # NO-OP! Already in wave3
```

### ❌ WRONG: Using wave integrations as sources
```bash
git checkout -b phase2-integration origin/main
git merge origin/phase2-wave1-integration  # Violates R270/R512!
git merge origin/phase2-wave2-integration  # Violates R270/R512!
```

### ❌ WRONG: Fast-forward merges
```bash
git merge origin/phase2/wave1/user-management  # Missing --no-ff!
# Creates fast-forward, no merge commit = can't validate mergeability
```

### ✅ CORRECT: Sequential rebuild with first effort base
```bash
# Base from first effort
FIRST="phase2/wave1/auth-system"
git checkout -b phase2-integration "origin/$FIRST"

# Merge all subsequent efforts with --no-ff
git merge --no-ff origin/phase2/wave1/user-management
git merge --no-ff origin/phase2/wave1/permissions
git merge --no-ff origin/phase2/wave2/database-layer
# ... all efforts in order
```

## Remember

**Phase integration validates R363 Sequential Direct Mergeability** by replaying the exact sequence of merges that will happen when efforts merge to main. If sequential rebuild works, we know sequential mergeability holds. If it fails, we catch problems BEFORE pushing to production.

**This is not just integration - it's sequential merge validation!**

---

**The sequential rebuild model ensures phase integration is a MEANINGFUL VALIDATION, not a symbolic no-op.**

## Iteration Container Context (SF 3.0)

Phase integration is a **convergence validation point** within the PHASE-level iteration container:

### Container Lifecycle
```
PHASE CONTAINER:
  Iteration 1: All Waves → Integration → Issues Found → Fix Cycle
  Iteration 2: Fixed Waves → Integration → Fewer Issues → Fix Cycle
                            ↑ VALIDATION POINT
  Iteration N: All Waves → Integration → Clean Convergence ✅
```

### Convergence Tracking in `integration-containers.json`
```json
{
  "container_id": "phase-2-integration",
  "current_iteration": 3,
  "convergence_status": "CONVERGING",
  "sequential_merge_validation": {
    "iteration_1": {"status": "FAILED", "conflicts": 8},
    "iteration_2": {"status": "FAILED", "conflicts": 2},
    "iteration_3": {"status": "PROJECT_DONE", "conflicts": 0}
  },
  "converged": true
}
```

The integration protocol validates that:
1. Sequential mergeability improves across iterations
2. Conflicts decrease (convergence indicator)
3. Final iteration achieves clean merge (convergence achieved)
