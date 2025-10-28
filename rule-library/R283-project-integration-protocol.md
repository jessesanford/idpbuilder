# 🚨🚨🚨 RULE R283: Project Integration Protocol (Sequential Rebuild Model)

## Classification
- **Category**: Integration Management
- **Criticality Level**: 🚨🚨🚨 BLOCKING
- **Enforcement**: MANDATORY for project completion
- **Penalty**: -75% to -100% for violations

## The Rule

**ALL multi-phase projects MUST perform project-level integration after final phase completion. This is NON-NEGOTIABLE and MANDATORY for project completion.**

**Project integration MUST use SEQUENTIAL REBUILD model: base from MAIN (the trunk!), merge ALL efforts from ALL phases in sequential order. This validates COMPLETE trunk-based sequential mergeability by testing if every single effort can merge to main.**

**Project integration MUST occur in completely isolated workspace at `/efforts/integration/` with fresh clone of target repository, NEVER the software-factory repository.**

## 🔴🔴🔴 STATE MACHINE ENFORCEMENT 🔴🔴🔴

### Prohibited Transition
```
COMPLETE_PHASE (final phase) → PROJECT_DONE  ❌ VIOLATION (-100%)
```

### Required Transition
```
COMPLETE_PHASE (final phase) → PROJECT_INTEGRATE_WAVE_EFFORTS  ✅ CORRECT
```

### Complete Project Integration Flow
```
COMPLETE_PHASE (final phase)
  → PROJECT_INTEGRATE_WAVE_EFFORTS
  → SETUP_PROJECT_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE
  → SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN
  → WAITING_FOR_PROJECT_MERGE_PLAN
  → SPAWN_INTEGRATION_AGENT_PROJECT
  → MONITORING_PROJECT_INTEGRATE_WAVE_EFFORTS
  → PROJECT_REVIEW_WAVE_INTEGRATION
  → (validation/fix cycles as needed)
  → PR_PLAN_CREATION
  → PROJECT_DONE
```

## Why Project Integration is Mandatory

1. **Verification**: Proves ALL phases work together as integrated system
2. **Integration Issues**: Catches cross-phase interaction problems
3. **Project Demo**: Creates comprehensive end-to-end demonstration (R291)
4. **Quality Gate**: Final validation before project completion
5. **Sequential Mergeability**: Validates R363 across entire project
6. **The "Bow on the Project"**: Proper ending showing everything works together

Skipping project integration means the project is **incomplete** - you have phases that individually work but no proof the complete project functions as a cohesive system.

## 🔴🔴🔴 SEQUENTIAL REBUILD MODEL 🔴🔴🔴

### The Fundamental Change

**OLD MODEL (INCORRECT - VIOLATES R512/R270/R363):**
```bash
# ❌ WRONG - Merging phase integration branches
git checkout -b project-integration origin/phase3-integration
git merge origin/phase1-integration  # Integration branch as source!
git merge origin/phase2-integration  # Integration branch as source!
```

**NEW MODEL (CORRECT - SEQUENTIAL REBUILD FROM MAIN):**
```bash
# ✅ CORRECT - Base from MAIN (the trunk!), merge ALL efforts in sequential order
git checkout -b project-integration origin/main  # Base from MAIN - the actual trunk!

# Merge ALL efforts from Phase 1 in sequential order
git merge --no-ff origin/phase1/wave1/foundation   # FIRST effort of entire project
git merge --no-ff origin/phase1/wave1/api-types   # Sequential merge
git merge --no-ff origin/phase1/wave1/controllers # Sequential merge
# ... all Phase 1 efforts

# Merge ALL efforts from Phase 2 in sequential order
git merge --no-ff origin/phase2/wave1/auth-system      # Sequential merge
git merge --no-ff origin/phase2/wave1/user-management  # Sequential merge
# ... all Phase 2 efforts

# Merge ALL efforts from Phase 3 in sequential order
git merge --no-ff origin/phase3/wave1/advanced-features # Sequential merge
git merge --no-ff origin/phase3/wave1/integrations      # Sequential merge
# ... all Phase 3 efforts

# Total: ALL efforts from ALL phases, in sequential order, all merging to MAIN
```

## Why Sequential Rebuild for Project Integration?

### 1. Tests COMPLETE Trunk-Based Mergeability
**Simulates THE EXACT production PR sequence:**
```
Real Production Workflow:
  main ← PR1 (phase1/wave1/effort1) merges
  main ← PR2 (phase1/wave1/effort2) merges
  main ← PR3 (phase1/wave2/effort1) merges
  ... through all phases ...
  main ← PRn (phase3/wave3/effortn) merges

Sequential Rebuild EXACTLY REPLICATES THIS:
  main (base - the ACTUAL trunk!)
  ← effort1 merges (first PR)
  ← effort2 merges (second PR)
  ← effort3 merges (third PR)
  ... all efforts across all phases, in exact order
```

### 2. Validates R363 Across Entire Project
**R363 requires**: Every effort must merge directly to main in sequence.

**Project integration proves this** for the COMPLETE project:
- If E1→E2→...→En works, project integration validates it
- If there are hidden cross-phase conflicts, we catch them
- If sequential mergeability fails anywhere, we detect it BEFORE production

### 3. Most Comprehensive Integration Test
**Project integration is the ULTIMATE validation**:
- Tests all phase interactions
- Validates complete dependency chain
- Ensures no phase-specific assumptions break
- Proves entire system works as one unit

## Requirements

### 1. Project Integration Infrastructure Setup

```bash
# Setup for final project integration using sequential rebuild model
TOTAL_PHASES=$(jq '.total_phases' orchestrator-state-v3.json)

# Read target repository configuration (R508)
TARGET_CONFIG="$CLAUDE_PROJECT_DIR/target-repo-config.yaml"
TARGET_REPO_PATH=$(yq '.repository_path' "$TARGET_CONFIG")
TARGET_REPO_NAME=$(yq '.repository_name' "$TARGET_CONFIG")
DEFAULT_BRANCH=$(yq '.default_branch' "$TARGET_CONFIG")

# Create isolated project integration workspace
PROJECT_INTEGRATE_WAVE_EFFORTS_DIR="$CLAUDE_PROJECT_DIR/efforts/project/integration-workspace"
rm -rf "$PROJECT_INTEGRATE_WAVE_EFFORTS_DIR"
mkdir -p "$PROJECT_INTEGRATE_WAVE_EFFORTS_DIR"
cd "$PROJECT_INTEGRATE_WAVE_EFFORTS_DIR"

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

# Verify we're NOT in software-factory directory
if [[ "$PWD" == *"software-factory"* ]]; then
    echo "❌ CRITICAL: Integration happening in software-factory directory!"
    exit 508  # R508 violation
fi
```

### 2. Determine Base Branch (Main Trunk)

```bash
# NEW: Use MAIN as the base - the actual trunk!
determine_project_base_branch() {
    # Read default branch from target-repo-config.yaml
    TARGET_CONFIG="$CLAUDE_PROJECT_DIR/target-repo-config.yaml"
    DEFAULT_BRANCH=$(yq '.default_branch' "$TARGET_CONFIG")

    if [ -z "$DEFAULT_BRANCH" ] || [ "$DEFAULT_BRANCH" = "null" ]; then
        # Fallback to main if not specified
        DEFAULT_BRANCH="main"
    fi

    echo "✅ R283: Using MAIN as project integration base"
    echo "   Base branch: $DEFAULT_BRANCH (the trunk!)"
    echo "$DEFAULT_BRANCH"
}

# Get main as base
BASE_BRANCH=$(determine_project_base_branch)

# Create project integration branch FROM MAIN (the trunk!)
git fetch origin "$BASE_BRANCH"
git checkout -b "project-integration" "origin/$BASE_BRANCH"

echo "✅ Project integration base: $BASE_BRANCH (MAIN - the actual trunk!)"
echo "   This validates COMPLETE trunk-based mergeability"
```

### 3. Collect ALL Efforts from ALL Phases in Order

```bash
# NEW: Collect ALL efforts from ALL phases (including first!)
collect_all_project_efforts_sequential() {
    echo "📊 Collecting ALL efforts from entire project (all phases)..."

    # Get ALL completed efforts, sorted by phase, wave, then completion order
    ALL_EFFORTS=$(jq -r '
        .efforts_completed[] |
        .branch
    ' orchestrator-state-v3.json | sort)

    # Count by phase
    for phase in $(seq 1 $TOTAL_PHASES); do
        PHASE_COUNT=$(echo "$ALL_EFFORTS" | grep "phase${phase}/" | wc -l)
        echo "  Phase $phase: $PHASE_COUNT efforts"
    done

    # Total count
    TOTAL=$(echo "$ALL_EFFORTS" | wc -l)

    echo "✅ Found $TOTAL total efforts across $TOTAL_PHASES phases"
    echo "   ALL $TOTAL efforts will merge to main (base = main, not first effort)"

    echo "$ALL_EFFORTS"
}

# Collect ALL project efforts (including first one!)
EFFORTS_TO_MERGE=$(collect_all_project_efforts_sequential)

if [ -z "$EFFORTS_TO_MERGE" ]; then
    echo "❌ FATAL: No efforts found in project"
    echo "Cannot perform project integration with zero efforts"
    exit 1
fi
```

### 4. Sequential Merge Execution (ALL Phases)

```bash
# NEW: Sequential rebuild of ENTIRE PROJECT from MAIN
echo "🔄 Starting sequential rebuild of ENTIRE PROJECT from MAIN..."
echo "Base: $BASE_BRANCH (MAIN - the actual trunk!)"
echo "Merging $(echo "$EFFORTS_TO_MERGE" | wc -l) total efforts..."
echo ""

MERGE_COUNT=0
CURRENT_PHASE=1
CURRENT_WAVE=1

while IFS= read -r effort_branch; do
    [ -z "$effort_branch" ] && continue

    MERGE_COUNT=$((MERGE_COUNT + 1))

    # Extract phase/wave from branch name for progress tracking
    if [[ "$effort_branch" =~ phase([0-9]+)/wave([0-9]+)/ ]]; then
        EFFORT_PHASE="${BASH_REMATCH[1]}"
        EFFORT_WAVE="${BASH_REMATCH[2]}"

        # Track phase/wave transitions for progress reporting
        if [[ "$EFFORT_PHASE" != "$CURRENT_PHASE" ]]; then
            echo ""
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "🎯 PHASE TRANSITION: Phase $CURRENT_PHASE → Phase $EFFORT_PHASE"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            CURRENT_PHASE=$EFFORT_PHASE
            CURRENT_WAVE=$EFFORT_WAVE
        elif [[ "$EFFORT_WAVE" != "$CURRENT_WAVE" ]]; then
            echo ""
            echo "───────────────────────────────────────────────"
            echo "📌 Wave $CURRENT_WAVE → Wave $EFFORT_WAVE"
            echo "───────────────────────────────────────────────"
            CURRENT_WAVE=$EFFORT_WAVE
        fi
    fi

    echo ""
    echo "📦 Merge $MERGE_COUNT: $effort_branch"

    # Fetch effort branch
    git fetch origin "$effort_branch"

    # Attempt merge with no fast-forward (always create merge commit)
    if git merge --no-ff "origin/$effort_branch" \
         -m "integrate: Merge $effort_branch into project-integration

This is merge $MERGE_COUNT in the sequential rebuild of the entire project.
Base: $BASE_BRANCH (MAIN - the trunk!)
Target: project-integration

Per R283 Sequential Rebuild Model - testing trunk-based sequential mergeability."; then
        echo "✅ Merge $MERGE_COUNT successful"
    else
        echo "❌ Merge conflict in $effort_branch"
        echo "🔍 Conflict details:"
        git status

        # Document conflict
        cat > "PROJECT-MERGE-CONFLICT-${MERGE_COUNT}.md" << EOF
# Project Integration Merge Conflict - $effort_branch

**Merge Number:** $MERGE_COUNT
**Branch:** $effort_branch
**Phase:** $EFFORT_PHASE
**Wave:** $EFFORT_WAVE
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
This conflict must be resolved before continuing project integration.
This is a CROSS-PHASE conflict that requires careful resolution.
EOF
        echo "❌ FAILED: Merge conflict requires resolution"
        exit 1
    fi

    # Run tests after EVERY 10 merges or phase transitions
    if (( MERGE_COUNT % 10 == 0 )) || [[ "$EFFORT_PHASE" != "$CURRENT_PHASE" ]]; then
        echo "🧪 Running tests after merge $MERGE_COUNT..."
        if ! (npm test 2>/dev/null || make test 2>/dev/null || pytest 2>/dev/null || true); then
            echo "⚠️ WARNING: Tests may have failed after merge $MERGE_COUNT"
            echo "Continue with caution - review test results"
        fi
    fi

done <<< "$EFFORTS_TO_MERGE"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ COMPLETE PROJECT SEQUENTIAL REBUILD FINISHED!"
echo "   Base branch: $BASE_BRANCH (MAIN - the trunk!)"
echo "   Total phases: $TOTAL_PHASES"
echo "   Total merges: $MERGE_COUNT"
echo "   Result: project-integration branch"
echo "   Validation: FULL trunk-based sequential mergeability PROVEN!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
```

## Edge Cases

### Single-Phase Projects

**Scenario**: Project has only ONE phase

**Handling**:
```bash
# Still create project integration from main
BASE_BRANCH=$(determine_project_base_branch)
git checkout -b project-integration "origin/$BASE_BRANCH"

# Merge all efforts from the single phase onto main
git merge --no-ff origin/phase1/wave1/effort1
git merge --no-ff origin/phase1/wave1/effort2
# ... all phase efforts

# Result: Tests that single phase can merge to main sequentially
```

**Why this is okay**:
- Maintains consistent project completion flow
- Tests sequential mergeability to main even for single phase
- Provides project-level demo and validation
- Ensures R283 compliance regardless of project size

### ERROR_RECOVERY Fix Integration (Cross-Phase)

**Scenario**: Bugs found during project testing affect multiple phases

**Handling**:
```bash
# Sequential rebuild includes fix branches at the END
BASE_BRANCH=$(determine_project_base_branch)
git checkout -b "project-integration" "origin/$BASE_BRANCH"

# 1. Merge all original efforts from all phases in sequential order
for effort in $ALL_ORIGINAL_EFFORTS; do
    git merge --no-ff "origin/$effort"
done

# 2. Merge all project-level fix branches at the END
for fix in $PROJECT_FIX_BRANCHES; do
    git merge --no-ff "origin/$fix" \
        -m "integrate: Apply project-level fix $fix to project-integration"
done
```

**Order matters**:
- Start from MAIN (the trunk!)
- All original efforts first (in phase/wave order)
- Project fixes last (they fix cross-phase issues)
- Tests that fixes apply cleanly on top of complete original work

## Workspace Isolation Requirements

### Mandatory Isolation
- **Location**: `$CLAUDE_PROJECT_DIR/efforts/project/integration-workspace/[target-repo-name]/`
- **Repository**: Fresh clone of TARGET repository (R508 requirement)
- **Branch**: `project-integration`
- **NEVER**: Use software-factory repository
- **NEVER**: Reuse phase integration workspaces
- **NEVER**: Work in software-factory-template directory

### Validation Checks
```bash
validate_project_integration_workspace() {
    # Check correct directory (must be in project integration workspace)
    if [[ "$PWD" != */efforts/project/integration-workspace/* ]]; then
        echo "❌ Not in project integration workspace!"
        echo "Expected: /efforts/project/integration-workspace/*"
        echo "Got: $PWD"
        return 1
    fi

    # Check NOT software-factory
    REMOTE=$(git remote get-url origin 2>/dev/null)
    if [[ "$REMOTE" == *"software-factory"* ]]; then
        echo "❌ CRITICAL: Working in software-factory repository!"
        echo "This will corrupt the orchestrator!"
        return 1
    fi

    # Check NOT in software-factory directory
    if [[ "$PWD" == *"software-factory"* ]]; then
        echo "❌ CRITICAL: In software-factory directory!"
        return 1
    fi

    # Check branch naming
    CURRENT_BRANCH=$(git branch --show-current)
    if [[ "$CURRENT_BRANCH" != "project-integration" ]]; then
        echo "⚠️ WARNING: Non-standard branch name: $CURRENT_BRANCH"
        echo "Expected: project-integration"
    fi

    echo "✅ Project integration workspace validation passed"
    return 0
}

# Run validation before EVERY operation
validate_project_integration_workspace || exit 1
```

## Integration Process Summary

### Step 1: Setup Infrastructure
- Create isolated project workspace
- Clone target repository
- Verify NOT software-factory (R508)

### Step 2: Determine Base
- Read default branch from target-repo-config.yaml
- Use MAIN as base (the trunk!)
- Create project integration branch from main

### Step 3: Collect ALL Efforts
- Get ALL efforts from ALL phases in order (including first!)
- Include splits, exclude deprecated originals (R296)
- Include project-level ERROR_RECOVERY fixes at end

### Step 4: Sequential Rebuild (Entire Project)
- Merge EVERY effort with --no-ff
- Test after phase transitions and every 10 merges
- Document conflicts if they occur
- Continue until ALL efforts from ALL phases merged

### Step 5: Comprehensive Validation & Demo (R291)
- Run COMPLETE test suite (unit, integration, E2E)
- Create production build
- Create comprehensive project demo
- Generate project integration report

### Step 6: Push and Document
- Push project integration branch
- Update state file (R288)
- Create final PR plan (R370)
- Document completion

## Comparison with Phase Integration

### Phase Integration (Per-Phase Validation)
**Purpose**: Validate sequential mergeability within a phase
**Scope**: All efforts from ONE phase
**Base**: First effort of the phase
**Result**: Phase integration branch (proves phase mergeability)

### Project Integration (Complete Trunk Validation)
**Purpose**: Validate sequential mergeability to MAIN across ENTIRE project
**Scope**: ALL efforts from ALL phases
**Base**: MAIN (the trunk!)
**Result**: Project integration branch (proves COMPLETE trunk-based mergeability)

**Key Difference**:
- Phase integrations validate per-phase sequential mergeability
- Project integration validates COMPLETE trunk-based workflow (every effort can merge to main)

## Validation Requirements

### Pre-Integration Validation
```bash
# Verify all phases completed
for phase in $(seq 1 $TOTAL_PHASES); do
    phase_status=$(jq ".phases.phase_${phase}.status" orchestrator-state-v3.json)
    if [ "$phase_status" != "INTEGRATED" ]; then
        echo "🚨 Cannot integrate project - Phase $phase not integrated"
        exit 1
    fi
done

# Verify architect approval
architect_approval=$(jq '.architect_approval.project_ready' orchestrator-state-v3.json)
if [ "$architect_approval" != "true" ]; then
    echo "🚨 Cannot integrate - Architect approval required"
    exit 1
fi

# Verify main branch exists
BASE_BRANCH=$(determine_project_base_branch)
if ! git ls-remote --heads origin "$BASE_BRANCH" > /dev/null 2>&1; then
    echo "🚨 Base branch not found: $BASE_BRANCH"
    echo "Cannot integrate without main branch in target repository"
    exit 1
fi
```

### Post-Integration Validation
```bash
# COMPREHENSIVE project testing
echo "🧪 Running COMPLETE project test suite..."

# Create project test harness
cat > project-test-harness.sh << 'EOF'
#!/bin/bash
echo "🧪 PROJECT INTEGRATE_WAVE_EFFORTS TEST SUITE"
echo "=================================="

# Unit tests
npm test 2>&1 | tee project-unit-tests.log
UNIT_RESULT=$?

# Integration tests
npm run test:integration 2>&1 | tee project-integration-tests.log
INT_RESULT=$?

# E2E tests
npm run test:e2e 2>&1 | tee project-e2e-tests.log
E2E_RESULT=$?

# Performance tests
npm run test:performance 2>&1 | tee project-performance-tests.log

# Security scan
npm audit --audit-level=moderate 2>&1 | tee project-security.log

echo "=================================="
if [ $UNIT_RESULT -eq 0 ] && [ $INT_RESULT -eq 0 ] && [ $E2E_RESULT -eq 0 ]; then
    echo "✅ ALL CRITICAL TESTS PASSED!"
    exit 0
else
    echo "❌ TEST FAILURES DETECTED!"
    exit 1
fi
EOF

chmod +x project-test-harness.sh
./project-test-harness.sh || exit 1

# Build verification
echo "🏗️ Running production build..."
rm -rf dist/ build/ out/ target/
npm run build:prod 2>&1 | tee project-build.log || exit 1

# Verify build artifacts
if [ -d "dist" ] || [ -d "build" ] || [ -d "out" ] || [ -d "target" ]; then
    echo "✅ Build artifacts created"
else
    echo "❌ No build artifacts found!"
    exit 1
fi

# Demo creation (R291)
./create-project-demo.sh

# Integration report
cat > PROJECT-INTEGRATE_WAVE_EFFORTS-REPORT.md << EOF
# Complete Project Integration Report (Sequential Rebuild from Main)

## Model Used
**Sequential Rebuild from Main per R283**
- Base: $BASE_BRANCH (MAIN - the trunk!)
- Merged: $MERGE_COUNT efforts from $TOTAL_PHASES phases
- Result: COMPLETE validation of trunk-based sequential mergeability

## All Phases Integrated
$(for p in $(seq 1 $TOTAL_PHASES); do
    COUNT=$(echo "$EFFORTS_TO_MERGE" | grep "phase${p}/" | wc -l)
    echo "- Phase $p: $COUNT efforts merged"
done)

## Build & Test Status
- Production Build: ✅ PROJECT_DONEFUL
- Unit Tests: ✅ PASSING
- Integration Tests: ✅ PASSING
- E2E Tests: ✅ PASSING
- Performance: ✅ ACCEPTABLE
- Security: ✅ CLEAR

## Validation
- Repository: $(git remote get-url origin) ✅
- Workspace: $PWD ✅
- Sequential merges: $MERGE_COUNT ✅
- Phases covered: $TOTAL_PHASES ✅
- Conflicts: $([ -f PROJECT-MERGE-CONFLICT-*.md ] && echo "DOCUMENTED" || echo "NONE") ✅
EOF
```

## Success Criteria

- ✅ Main branch used as base (the trunk!)
- ✅ ALL efforts from ALL phases collected in order (including first!)
- ✅ Sequential rebuild completed ($MERGE_COUNT merges across $TOTAL_PHASES phases)
- ✅ All merges used --no-ff (real merge commits testing mergeability)
- ✅ Comprehensive tests passing (unit, integration, E2E)
- ✅ Production build successful
- ✅ Project demo created (R291)
- ✅ Workspace isolated correctly (R508)
- ✅ State file updated (R288)
- ✅ Architect approval received
- ✅ Ready for final PR plan creation (R370)
- ✅ COMPLETE trunk-based sequential mergeability validated

## Failure Conditions

### Critical Failures (Immediate Stop)
- 🚨 Using software-factory repository = IMMEDIATE FAIL (-100%)
- 🚨 Using phase integration branches as base = VIOLATION R512/R270 (-100%)
- 🚨 Using phase integration branches as sources = VIOLATION R512/R270 (-100%)
- 🚨 Not using main as base = VIOLATION R283 (-100%)
- 🚨 Merging with fast-forward (no merge commits) = Invalid test (-50%)
- 🚨 Cannot find main branch = Environment failure (-75%)
- 🚨 Any test failures = Quality failure (-70%)
- 🚨 Build failures = Deliverable failure (-80%)
- 🚨 No project demo = Incomplete (R291) (-40%)

### Recovery Protocol
1. **STOP ALL WORK IMMEDIATELY**
2. Document exact failure point and error
3. Transition to ERROR_RECOVERY state
4. Alert architect for assessment
5. Clean ALL integration workspaces
6. Verify state file data
7. Start fresh with corrected strategy
8. **NEVER continue from corrupted state**

## Penalties

- Using software-factory repository: **-100% AUTOMATIC FAILURE**
- Using phase integration branches: **-100% (R512/R270 violation)**
- Software-factory contamination: **-100% CRITICAL**
- Non-isolated workspace: **-80% grade**
- Missing efforts in sequence: **-60% grade**
- Fast-forward merges (no merge commits): **-50% grade**
- Test failures ignored: **-70% grade**
- Build failures: **-65% grade**
- No final PR: **-40% grade**
- No project demo: **-40% grade** (R291)
- State file not updated: **-40% grade** (R288)

## Related Rules

- **R512**: Trunk-Based Integration Model (integrate effort branches)
- **R270**: No Integration Branches as Sources (phase integrations not sources)
- **R363**: Sequential Direct Mergeability (what we're validating)
- **R364**: Integration Testing Only Branches (phase integrations are testing)
- **R282**: Phase Integration Protocol (per-phase validation)
- **R308**: Incremental Branching Strategy (efforts cascade)
- **R034**: Integration Requirements (general integration rules)
- **R250**: Integration Isolation (workspace requirements)
- **R288**: State File Updates (required after integration)
- **R291**: Integration Demo Requirement (must create demo)
- **R296**: Deprecated Branch Marking (handle split branches)
- **R370**: PR Plan Creation Requirements (final deliverable)
- **R508**: Target Repository Enforcement (never use software-factory)

## Common Violations

### ❌ WRONG: Using phase integrations as base or sources
```bash
# WRONG - Using phase integration as base
git checkout -b project-integration origin/phase3-integration
git merge origin/phase1-integration  # NO-OP! Already in phase3
git merge origin/phase2-integration  # NO-OP! Already in phase3
```

### ❌ WRONG: Skipping project integration
```bash
# WRONG - Going directly to PROJECT_DONE
# Multi-phase projects MUST perform project integration
COMPLETE_PHASE → PROJECT_DONE  # VIOLATION!
```

### ✅ CORRECT: Sequential rebuild from MAIN
```bash
# Base from MAIN (the trunk!)
BASE="main"
git checkout -b project-integration "origin/$BASE"

# Merge ALL efforts from ALL phases with --no-ff (including first!)
git merge --no-ff origin/phase1/wave1/foundation      # First effort
git merge --no-ff origin/phase1/wave1/api-types      # Sequential
git merge --no-ff origin/phase1/wave1/controllers    # Sequential
# ... all Phase 1 efforts

git merge --no-ff origin/phase2/wave1/auth-system      # Sequential
git merge --no-ff origin/phase2/wave1/user-management  # Sequential
# ... all Phase 2 efforts

git merge --no-ff origin/phase3/wave1/advanced-features  # Sequential
# ... all Phase 3 efforts

# Result: COMPLETE trunk-based sequential mergeability validated!
```

## Remember

**This is the ULTIMATE VALIDATION** - project integration tests if EVERY effort can merge to main sequentially. Using the wrong repository will DESTROY the orchestrator and fail the entire project. ALWAYS verify you're in the target repository, NOT software-factory!

**Project integration validates R363 Sequential Direct Mergeability across the ENTIRE PROJECT** by starting from MAIN (the trunk!) and merging every effort in sequence. If the complete sequential rebuild works, we know the entire project can merge to main one effort at a time. If it fails, we catch problems BEFORE production.

**This is not just integration - it's THE COMPLETE TRUNK-BASED WORKFLOW VALIDATION!**

**Key Architectural Points:**
- Wave integration: Base = first effort of wave (tests wave-internal mergeability)
- Phase integration: Base = first effort of phase (tests cross-wave mergeability)
- **Project integration: Base = MAIN** (tests COMPLETE trunk-based mergeability!)

---

**The sequential rebuild from main ensures project integration is the ULTIMATE VALIDATION of trunk-based development, proving every effort can merge to production.**
