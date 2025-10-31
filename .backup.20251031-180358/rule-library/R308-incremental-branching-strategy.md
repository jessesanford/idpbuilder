# 🔴🔴🔴 RULE R308 - Progressive Trunk-Based Development (CASCADE BRANCHING) [CORE TENANT]

## Classification
- **Category**: Core Development Flow
- **Criticality Level**: 🔴🔴🔴 SUPREME LAW
- **Enforcement**: MANDATORY - Progressive trunk-based development foundation
- **Penalty**: -100% for violating progressive cascade flow

## The Rule

**EVERY effort MUST branch from the PREVIOUS EFFORT in a cascading chain. This is PROGRESSIVE TRUNK-based development where each effort builds incrementally on the previous effort's work, creating a cascade of changes that can be merged sequentially without conflicts.**

## 🔴🔴🔴 CRITICAL SCOPE CLARIFICATION 🔴🔴🔴

**THIS RULE GOVERNS EFFORT BRANCHING (DEVELOPMENT WORKFLOW) ONLY!**

**R308 applies to:** EFFORT branch creation strategy (how efforts cascade from each other)
**R308 does NOT apply to:** Integration branch base determination (governed by R009/R282/R283/R512)

### What R308 Defines:
```
EFFORT CASCADE (Development Workflow):
P1W1: main → effort1 → effort2 → effort3
P1W2: effort3 → effort4 → effort5
P2W1: effort5 → effort6 → effort7
```

### What R308 Does NOT Define:
```
INTEGRATE_WAVE_EFFORTS BASES (Testing Workflow - R009/R282/R283):
Wave integration: base = first effort of WAVE (NOT R308!)
Phase integration: base = first effort of PHASE (NOT R308!)
Project integration: base = main (NOT R308!)
```

**DO NOT USE R308 TO DETERMINE INTEGRATE_WAVE_EFFORTS BRANCH BASES!**
**Integration bases follow Sequential Rebuild Model (R009/R282/R283/R512)**

### 🔴 THE CASCADE PRINCIPLE 🔴

```
WITHIN EACH WAVE - EFFORTS CASCADE:
Phase 1, Wave 1:
  main → effort-1 → effort-2 → effort-3 → effort-4

Phase 1, Wave 2 (starts from P1W1 last effort):
  effort-4 → effort-5 → effort-6 → effort-7

Phase 2, Wave 1 (starts from P1W2 last effort):
  effort-7 → effort-8 → effort-9 → effort-10

NEVER PARALLEL BRANCHING FROM MAIN:
  ❌ main → effort-1
  ❌ main → effort-2  (WRONG! Should be from effort-1)
  ❌ main → effort-3  (WRONG! Should be from effort-2)
```

**EACH EFFORT CONTAINS ALL PREVIOUS EFFORTS' CHANGES!**

## Why This Is CRITICAL

### Traditional (WRONG) Approach:
```
main
  ├─→ ALL Phase 1 efforts (stale base!)
  ├─→ ALL Phase 2 efforts (conflicts!)
  └─→ ALL Phase 3 efforts (massive conflicts!)
```
**Result**: Integration nightmare, "big bang" merges, conflicts discovered late

### Incremental (CORRECT) Approach:
```
main
  ├─→ P1W1 efforts → integrate → P1W1-integration
                                    ├─→ P1W2 efforts → integrate → P1W2-integration
                                                                     ├─→ P1W3 efforts
```
**Result**: Conflicts detected early, smooth integration, true CI/CD

## Mandatory Implementation

### 1. Orchestrator MUST Record Base Branch in State (R337)

```bash
# 🔴🔴🔴 CRITICAL: Write to state file for CASCADE branching! 🔴🔴🔴
record_cascade_base_branch() {
    local PHASE=$1
    local WAVE=$2
    local EFFORT=$3
    local EFFORT_INDEX=$4  # Position in effort sequence (1, 2, 3, etc.)

    # Determine base per CASCADE strategy
    local BASE=""
    local REASON=""

    # FIRST effort in P1W1: Start from main
    if [[ $PHASE -eq 1 && $WAVE -eq 1 && $EFFORT_INDEX -eq 1 ]]; then
        BASE="main"
        REASON="First effort of Phase 1 Wave 1 starts from main per R308 cascade"
    # FIRST effort in new wave: From LAST effort of previous wave
    elif [[ $EFFORT_INDEX -eq 1 ]]; then
        # Get the last effort from previous wave
        if [[ $WAVE -eq 1 ]]; then
            # First wave of new phase - get last effort from previous phase's last wave
            PREV_PHASE=$((PHASE - 1))
            BASE=$(jq -r --arg p "$PREV_PHASE" '
                .efforts_completed[] |
                select(.phase == ($p | tonumber)) |
                .branch
            ' orchestrator-state-v3.json | tail -1)
            REASON="First effort of Phase $PHASE Wave 1 cascades from last effort of Phase $PREV_PHASE per R308"
        else
            # Subsequent waves - get last effort from previous wave
            PREV_WAVE=$((WAVE - 1))
            BASE=$(jq -r --arg p "$PHASE" --arg w "$PREV_WAVE" '
                .efforts_completed[] |
                select(.phase == ($p | tonumber) and .wave == ($w | tonumber)) |
                .branch
            ' orchestrator-state-v3.json | tail -1)
            REASON="First effort of Wave $WAVE cascades from last effort of Wave $PREV_WAVE per R308"
        fi
    # SUBSEQUENT efforts: From PREVIOUS effort in same wave
    else
        PREV_INDEX=$((EFFORT_INDEX - 1))
        # Get the previous effort's branch name
        BASE=$(jq -r --arg p "$PHASE" --arg w "$WAVE" --arg idx "$PREV_INDEX" '
            (.efforts_in_progress[] | select(.phase == ($p | tonumber) and .wave == ($w | tonumber) and .index == ($idx | tonumber)) | .branch) //
            (.efforts_completed[] | select(.phase == ($p | tonumber) and .wave == ($w | tonumber) and .index == ($idx | tonumber)) | .branch)
        ' orchestrator-state-v3.json)

        if [ -z "$BASE" ] || [ "$BASE" = "null" ]; then
            # If not found by index, get the most recent effort in this wave
            BASE=$(jq -r --arg p "$PHASE" --arg w "$WAVE" '
                [.efforts_in_progress[], .efforts_completed[]] |
                map(select(.phase == ($p | tonumber) and .wave == ($w | tonumber))) |
                sort_by(.branched_at) |
                .[-1].branch
            ' orchestrator-state-v3.json)
        fi
        REASON="Effort $EFFORT_INDEX cascades from previous effort per R308 progressive trunk-based development"
    fi
    
    # RECORD IN STATE FILE (R337 MANDATORY)
    jq --arg effort "$EFFORT" \
       --arg base "$BASE" \
       --arg reason "$REASON" \
       --arg timestamp "$(date -Iseconds)" \
       '.base_branch_decisions.decision_log += [{
          "timestamp": $timestamp,
          "effort": $effort,
          "decision": "Base branch for \($effort): \($base)",
          "reason": $reason,
          "decided_by": "orchestrator-R308"
       }] |
       .efforts_in_progress += [{
          "name": $effort,
          "base_branch_tracking": {
            "planned_base": $base,
            "actual_base": $base,
            "branched_at": $timestamp,
            "reason": $reason
          }
       }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
    
    echo "✅ Base branch recorded in state: $EFFORT will use $BASE"
    echo "   Reason: $REASON"
}
```

### 2. Clone MUST Use Base from State File (R337)

```bash
# WRONG - Always using main
git clone --branch main "$REPO" "$EFFORT_DIR"

# WRONG - Calculating base
BASE=$(determine_effort_base_branch $PHASE $WAVE)

# CORRECT - Reading from state file per R337
BASE=$(jq -r --arg e "$EFFORT" '
    (.efforts_in_progress[] | select(.name == $e) | .base_branch_tracking.actual_base) //
    .base_branch_decisions.current_wave_base
' orchestrator-state-v3.json)

if [ -z "$BASE" ] || [ "$BASE" = "null" ]; then
    echo "❌ FATAL: No base branch in state for $EFFORT (R337 violation)"
    exit 1
fi

git clone --branch "$BASE" "$REPO" "$EFFORT_DIR"
```

### 3. Agents MUST Verify Their Base

```bash
verify_incremental_base() {
    local expected_base="$1"
    local actual_base=$(git merge-base HEAD origin/main)
    local base_branch=$(git branch -r --contains "$actual_base" | grep -E "(main|integration)" | head -1)
    
    if [[ ! "$base_branch" =~ "$expected_base" ]]; then
        echo "❌ FATAL: Not based on incremental integration!"
        echo "Expected base: $expected_base"
        echo "Actual base: $base_branch"
        exit 1
    fi
}
```

## 🔴🔴🔴 CRITICAL: R308 IS ABOUT EFFORT BRANCHING, NOT INTEGRATE_WAVE_EFFORTS BRANCHING! 🔴🔴🔴

**THIS RULE DEFINES HOW EFFORT BRANCHES ARE CREATED (development workflow)**
**NOT how integration branches are created (testing workflow)**

### Effort Branching Flow (R308 - This Rule):
```
P1W1 efforts cascade:
main → effort1 → effort2 → effort3

P1W2 efforts cascade from P1W1 last effort:
effort3 → effort4 → effort5 → effort6

Purpose: Incremental development, each effort builds on previous
```

### Integration Branching Flow (R009/R512 - Different Rules):
```
P1W1 integration (sequential rebuild):
base: effort1, merge: effort2, effort3

P1W2 integration (sequential rebuild):
base: effort4, merge: effort5, effort6

Phase integration (sequential rebuild):
base: effort1, merge: effort2, effort3, effort4, effort5, effort6

Purpose: Test sequential mergeability at different scopes
```

### Why They're Different:

1. **Effort Cascade (R308)**: For development convenience
   - Developer working on effort4 has all effort1-3 code available
   - Incremental development workflow
   - Each effort contains cumulative changes

2. **Integration Sequential Rebuild (R009/R512)**: For merge testing
   - Tests that efforts merge cleanly in sequence
   - Validates different merge topologies
   - Creates real merge validation

**DO NOT CONFUSE THESE TWO FLOWS!**

## 🔴🔴🔴 RE-INTEGRATE_WAVE_EFFORTS HANDLING (DETERMINISTIC) 🔴🔴🔴

### When Integration Fails and Must Be Redone:

**DETERMINISTIC RE-INTEGRATE_WAVE_EFFORTS PROTOCOL:**

1. **Directory Handling:**
   ```bash
   # Archive old integration workspace
   mv integration-workspace integration-workspace-archived-N
   # Create fresh integration workspace
   mkdir -p integration-workspace
   ```

2. **Branch Handling:**
   ```bash
   # Use SAME branch name, force-push updates
   BRANCH_NAME="phase2-wave1-integration"  # Same name
   git checkout -b "$BRANCH_NAME"
   git push --force-with-lease origin "$BRANCH_NAME"
   ```

3. **Base Branch Remains Same:**
   ```bash
   # Re-integration uses SAME base as original attempt
   BASE_BRANCH=$(determine_integration_base_branch)  # Same logic
   ```

**NEVER:**
- ❌ Create numbered integration branches (phase2-wave1-integration-v2)
- ❌ Leave old integration workspace in place
- ❌ Change the base branch for re-integration

**ALWAYS:**
- ✅ Archive old workspace with numbered suffix
- ✅ Force-push to same branch name
- ✅ Use exact same base branch determination

## Serial/Dependent Efforts (Sequential Chaining)

**🔴🔴🔴 CRITICAL: Dependent efforts MUST chain sequentially like splits! 🔴🔴🔴**

### The Serial Effort Principle:

When efforts within a wave have dependencies (as determined by R219 and R053):
- **DEPENDENT EFFORTS CHAIN SEQUENTIALLY** from each other
- **PARALLEL EFFORTS BRANCH IN PARALLEL** from the wave integration base

### Parallelization Decision Flow:
```
R219 (Effort Plan) defines dependencies
    ↓
R053 determines serial vs parallel execution
    ↓
R308 enforces correct branching strategy:
    - Parallel: All from wave integration base
    - Serial: Chain from each other
```

### Serial/Dependent Effort Branching:
```
CORRECT SERIAL BRANCHING (with dependencies):
phase2-wave1-integration
    ├─→ effort-1: Authentication (branches from integration)
    │       └─→ effort-2: User Profile (DEPENDS on Auth, branches from effort-1)
    │               └─→ effort-3: Permissions (DEPENDS on Profile, branches from effort-2)
    └─→ effort-4: Independent feature (parallel, branches from integration)

CORRECT PARALLEL BRANCHING (no dependencies):
phase2-wave2-integration
    ├─→ effort-1: UI Theme (branches from integration)
    ├─→ effort-2: Documentation (branches from integration)
    └─→ effort-3: Performance (branches from integration)
```

### 🔴 WHY SERIAL CHAINING FOR DEPENDENCIES:
1. **Code Availability**: Dependent efforts need prior effort's code
2. **Conflict Prevention**: Sequential merging prevents dependency conflicts
3. **Build Integrity**: Ensures dependent features compile/run correctly
4. **Testing Continuity**: Later efforts can test with earlier features present

### Implementation for Serial Efforts:
```bash
# Orchestrator determines effort dependencies from R219 plan
determine_effort_base_for_dependencies() {
    local PHASE=$1
    local WAVE=$2
    local EFFORT_NUM=$3
    local EFFORT_PLAN="phase${PHASE}-wave${WAVE}-effort-plan.md"
    
    # Check if effort has dependencies
    local DEPENDS_ON=$(grep -A2 "effort-${EFFORT_NUM}:" "$EFFORT_PLAN" | grep "depends_on:" | cut -d: -f2 | tr -d ' ')
    
    if [[ -n "$DEPENDS_ON" ]]; then
        # This effort depends on another - branch from it
        echo "phase${PHASE}/wave${WAVE}/${DEPENDS_ON}"
    else
        # No dependencies - use wave integration base
        determine_effort_base_branch $PHASE $WAVE
    fi
}

# Example usage for serial efforts:
# Effort-1 (no dependencies)
BASE=$(determine_effort_base_branch 2 2)  # Returns phase2-wave1-integration
git clone --branch "$BASE" "$REPO" "effort-1-auth"
cd "effort-1-auth"
git checkout -b "phase2/wave2/authentication"

# Effort-2 (depends on effort-1)
BASE="phase2/wave2/authentication"  # From effort-1
git clone --branch "$BASE" "$REPO" "effort-2-profile"
cd "effort-2-profile"
git checkout -b "phase2/wave2/user-profile"

# Effort-3 (depends on effort-2)
BASE="phase2/wave2/user-profile"  # From effort-2
git clone --branch "$BASE" "$REPO" "effort-3-permissions"
cd "effort-3-permissions"
git checkout -b "phase2/wave2/permissions"
```

### Verification for Serial Dependencies:
```bash
verify_dependency_chain() {
    local EFFORT_DIR=$1
    local EXPECTED_DEPENDENCY=$2
    
    cd "$EFFORT_DIR"
    
    # Check if dependency code is present
    if ! git log --oneline | grep -q "$EXPECTED_DEPENDENCY"; then
        echo "❌ FATAL: Missing dependency code from $EXPECTED_DEPENDENCY"
        echo "This effort must be based on its dependency!"
        exit 1
    fi
    
    echo "✅ Dependency chain verified: includes $EXPECTED_DEPENDENCY"
}
```

### Orchestrator Decision Logic:
```bash
# Read effort plan to determine parallelization
for effort in $(jq -r '.efforts_in_wave[]' orchestrator-state-v3.json); do
    # Check R219 effort plan for dependencies
    if effort_has_dependencies "$effort"; then
        echo "📌 Serial execution required for $effort"
        # Branch from previous effort
        spawn_agent_serial "$effort"
    else
        echo "📌 Parallel execution allowed for $effort"
        # Branch from wave integration base
        spawn_agent_parallel "$effort"
    fi
done
```

## Split Branching (Sequential Chaining - MANDATORY)

**🔴🔴🔴 CRITICAL: Splits MUST chain sequentially! 🔴🔴🔴**

### The Sequential Split Principle:

**🚨 CRITICAL CLARIFICATION: Split-001 is based on the SAME branch the oversized effort was based on, NOT the oversized branch itself! 🚨**

```
CORRECT SPLIT BRANCHING:
phase1/wave1-integration (or other R308 base)
    ├─→ effort-foo (becomes oversized, 1200+ lines, based on integration)
    ├─→ effort-foo-split-001 (based on SAME integration branch as effort-foo, NOT on effort-foo!)
    │       └─→ effort-foo-split-002 (based on split-001, NOT integration!)
    │               └─→ effort-foo-split-003 (based on split-002, NOT split-001!)
    └─→ effort-bar (normal effort, no split needed)

WRONG SPLIT BRANCHING (causes the bug):
phase1/wave1-integration
    └─→ effort-foo (oversized, 1200+ lines)
            └─→ effort-foo-split-001 (❌ WRONG: based on effort-foo, inherits ALL 1200+ lines!)
```

### 🔴 WHY THIS MATTERS:
When split-001 is incorrectly based on the oversized branch:
1. **line-counter.sh sees ALL changes** (1200+ lines) instead of just split-001's portion
2. **Split appears oversized** even though it only contains partial work
3. **Cascading failures** as each split appears to violate size limits
4. **Integration issues** as splits don't properly isolate changes

### ✅ CORRECT Split Branching Implementation:
```bash
# When creating split-001 infrastructure:
EFFORT_NAME="effort-foo"
OVERSIZED_BRANCH="phase1/wave1/effort-foo"  # Has 1200+ lines

# DETERMINE THE BASE BRANCH FOR SPLIT-001
# Get the base that the oversized effort was based on
BASE_FOR_OVERSIZED=$(cd /efforts/phase1/wave1/effort-foo && git merge-base HEAD origin/main | xargs git branch -r --contains | grep -E "(main|integration)" | head -1 | sed 's/.*origin\///')

# OR use the same R308 determination logic:
BASE_FOR_SPLIT_001=$(determine_effort_base_branch $PHASE $WAVE)  # Returns integration branch

# Create split-001 based on SAME base as oversized effort
git clone --branch "$BASE_FOR_SPLIT_001" "$REPO" "effort-foo-split-001"
cd effort-foo-split-001
git checkout -b "phase1/wave1/effort-foo-split-001"

# Result: split-001 starts CLEAN, only contains split-001's work
```

### ❌ WRONG Split Branching (The Bug):
```bash
# WRONG - Basing split-001 on the oversized branch
git clone --branch "phase1/wave1/effort-foo" "$REPO" "effort-foo-split-001"
# Result: split-001 INHERITS all 1200+ lines from effort-foo!
```

### ✅ CORRECT Sequential Chaining:
- **split-001**: From wave/phase integration branch (SAME as oversized effort's base)
- **split-002**: From split-001 branch (contains split-001's work)
- **split-003**: From split-002 branch (contains split-001 + split-002's work)
- **split-N**: From split-(N-1) branch (contains all previous splits' work)

### ❌ WRONG (Current Bug Pattern):
```
phase1/wave1-integration
    ├─→ split-001 (correct)
    ├─→ split-002 (WRONG! Missing split-001's work!)
    └─→ split-003 (WRONG! Missing split-001 & 002's work!)
```

### Why Sequential Chaining is MANDATORY:
1. **Preserves Work**: Each split builds on previous splits
2. **Prevents Conflicts**: Later splits see earlier splits' changes
3. **Maintains Atomicity**: Each split is still independent but cumulative
4. **Enables Clean Merge**: Final integration includes all splits in order

This maintains split atomicity while building incrementally.

### Implementation Example for Splits:
```bash
# Creating Split-001 (FIRST split)
INTEGRATE_WAVE_EFFORTS_BASE="phase1/wave1-integration"  # From R308 determine_effort_base_branch
git clone --branch "$INTEGRATE_WAVE_EFFORTS_BASE" "$REPO" "split-001-dir"
cd "split-001-dir"
git checkout -b "phase1/wave1/cert-validation-split-001"

# Creating Split-002 (builds on split-001)
PREVIOUS_SPLIT="phase1/wave1/cert-validation-split-001"  # NOT integration!
git clone --branch "$PREVIOUS_SPLIT" "$REPO" "split-002-dir"
cd "split-002-dir"
git checkout -b "phase1/wave1/cert-validation-split-002"

# Creating Split-003 (builds on split-002)
PREVIOUS_SPLIT="phase1/wave1/cert-validation-split-002"  # NOT integration!
git clone --branch "$PREVIOUS_SPLIT" "$REPO" "split-003-dir"
cd "split-003-dir"
git checkout -b "phase1/wave1/cert-validation-split-003"
```

**NEVER DO THIS (loses work between splits):**
```bash
# ❌ WRONG - All splits from same base!
git clone --branch "phase1/wave1-integration" ... split-001
git clone --branch "phase1/wave1-integration" ... split-002  # Missing split-001!
git clone --branch "phase1/wave1-integration" ... split-003  # Missing 001 & 002!
```

## Verification Protocol

### Pre-Clone Verification
```bash
echo "🔍 Verifying incremental base branch..."
PHASE=$(jq '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq '.current_wave' orchestrator-state-v3.json)
BASE=$(determine_effort_base_branch $PHASE $WAVE)

echo "📌 Phase $PHASE, Wave $WAVE"
echo "📌 Incremental base: $BASE"

if ! git ls-remote --heads origin "$BASE" > /dev/null 2>&1; then
    echo "❌ FATAL: Integration branch not ready!"
    echo "Previous wave must be integrated first"
    exit 1
fi
```

### Post-Clone Verification
```bash
cd "$EFFORT_DIR"
COMMITS_SINCE_MAIN=$(git rev-list --count origin/main..HEAD)
echo "📊 Commits ahead of main: $COMMITS_SINCE_MAIN"

# Should include previous wave's work
if [[ $WAVE -gt 1 && $COMMITS_SINCE_MAIN -eq 0 ]]; then
    echo "⚠️ WARNING: No previous wave commits found!"
    echo "May not be properly based on integration"
fi
```

## Examples

### ✅ CORRECT: Parallel Efforts (No Dependencies)
```bash
# All efforts branch from the same wave integration base
BASE=$(determine_effort_base_branch 2 2)  # Returns phase2-wave1-integration

# Effort 1: UI Theme
git clone --branch "$BASE" "$REPO" "efforts/phase2/wave2/ui-theme"
cd "efforts/phase2/wave2/ui-theme"
git checkout -b "phase2/wave2/ui-theme"

# Effort 2: Documentation (parallel, same base)
git clone --branch "$BASE" "$REPO" "efforts/phase2/wave2/documentation"
cd "efforts/phase2/wave2/documentation"
git checkout -b "phase2/wave2/documentation"

# Effort 3: Performance (parallel, same base)
git clone --branch "$BASE" "$REPO" "efforts/phase2/wave2/performance"
cd "efforts/phase2/wave2/performance"
git checkout -b "phase2/wave2/performance"
```

### ✅ CORRECT: Serial Efforts (With Dependencies)
```bash
# Effort 1: Authentication (no dependencies, uses wave base)
BASE=$(determine_effort_base_branch 2 2)  # Returns phase2-wave1-integration
git clone --branch "$BASE" "$REPO" "efforts/phase2/wave2/authentication"
cd "efforts/phase2/wave2/authentication"
git checkout -b "phase2/wave2/authentication"
# ... complete authentication implementation ...

# Effort 2: User Profile (depends on authentication)
BASE="phase2/wave2/authentication"  # From effort-1, NOT wave integration
git clone --branch "$BASE" "$REPO" "efforts/phase2/wave2/user-profile"
cd "efforts/phase2/wave2/user-profile"
git checkout -b "phase2/wave2/user-profile"
# Verify includes authentication code
git log --oneline | grep "authentication" || echo "❌ Missing dependency!"

# Effort 3: Permissions (depends on user-profile)
BASE="phase2/wave2/user-profile"  # From effort-2, NOT wave integration
git clone --branch "$BASE" "$REPO" "efforts/phase2/wave2/permissions"
cd "efforts/phase2/wave2/permissions"
git checkout -b "phase2/wave2/permissions"
# Verify includes both auth and profile code
git log --oneline | grep "profile" || echo "❌ Missing dependency!"
```

### ❌ WRONG: Always Using Main
```bash
# VIOLATION - Not incremental!
git clone --branch main "$REPO" "efforts/phase2/wave2/new-feature"
# Missing all Phase 1 and Phase 2 Wave 1 work!
```

### ❌ WRONG: Using Target Config Base Only
```bash
# VIOLATION - Ignoring integration!
BASE=$(yq '.base_branch' target-repo-config.yaml)  # Always "main"
git clone --branch "$BASE" "$REPO" "$EFFORT_DIR"
# Not building on previous work!
```

## Agent-Specific Requirements

### Orchestrator
- MUST determine correct incremental base
- MUST verify integration branch exists
- MUST document base branch decision
- MUST refuse to proceed if integration missing

### SW Engineer  
- MUST verify working from correct base
- MUST confirm includes previous wave work
- MUST refuse work if not incremental

### Code Reviewer
- MUST verify effort has correct base
- MUST check for integration conflicts
- MUST validate incremental development

### Integration Agent
- MUST create integration branches
- MUST push immediately for next wave
- MUST verify branch accessible

## Visual Flow Diagram

```
CASCADE BRANCHING (CORRECT - R308):
main
    └─→ P1W1-E1
            └─→ P1W1-E2
                    └─→ P1W1-E3
                            └─→ P1W1-E4 (last of wave 1)
                                    └─→ P1W2-E1 (starts from P1W1-E4)
                                            └─→ P1W2-E2
                                                    └─→ P1W2-E3 (last of wave 2)
                                                            └─→ P2W1-E1 (starts from P1W2-E3)
                                                                    └─→ P2W1-E2
                                                                            └─→ P2W1-E3

WRONG - PARALLEL BRANCHING FROM MAIN:
main
    ├─→ P1W1-E1  ❌ (only first should branch from main)
    ├─→ P1W1-E2  ❌ (should branch from E1)
    ├─→ P1W1-E3  ❌ (should branch from E2)
    └─→ P1W1-E4  ❌ (should branch from E3)

FINAL MERGE SEQUENCE (NO CONFLICTS):
1. Merge P1W1-E1 to main (contains only E1 changes)
2. Merge P1W1-E2 to main (contains E1+E2, assumes E1 already merged)
3. Merge P1W1-E3 to main (contains E1+E2+E3, assumes E2 already merged)
4. Continue sequentially...
```

## Grading Impact

- **Using stale base (main when should use integration)**: -100%
- **Skipping incremental flow**: -100%  
- **Not verifying base branch**: -50%
- **Proceeding without integration branch**: -75%
- **Not documenting base decision**: -25%

## Integration with Other Rules

- **R196**: Base branch selection (enhanced by R308)
- **R271**: Single-branch full checkout (uses R308 base)
- **R034**: Integration requirements (creates R308 bases)
- **R302**: Split tracking (sequential within incremental)
- **R304**: Line counting (counts from R308 base)
- **R219**: Dependency-Aware Effort Planning (determines serial vs parallel)
- **R053**: Parallelization Decisions (decides execution strategy)

## Stop Work Conditions

**IMMEDIATE STOP if:**
1. Previous wave integration branch doesn't exist
2. Effort not based on correct incremental branch
3. Integration conflicts detected with base
4. Base branch has pending unmerged work

## Rule Authority

This rule is a CORE TENANT of trunk-based development and SUPERSEDES any configuration that would cause efforts to branch from stale code. The incremental principle is NON-NEGOTIABLE.

---
**Remember**: Every effort builds on what came before. This is how we achieve TRUE continuous integration!