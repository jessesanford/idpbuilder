# 🔴🔴🔴 RULE R308 - Incremental Branching Strategy [CORE TENANT]

## Classification
- **Category**: Core Development Flow
- **Criticality Level**: 🔴🔴🔴 SUPREME LAW
- **Enforcement**: MANDATORY - Trunk-based development foundation
- **Penalty**: -100% for violating incremental flow

## The Rule

**EVERY effort MUST branch from the LATEST INTEGRATED CODE of the previous wave/phase. This is TRUE TRUNK-based development where each wave builds incrementally on the previous wave's integrated work.**

### 🔴 THE INCREMENTAL PRINCIPLE 🔴

```
Phase 1, Wave 1: Efforts branch from → main
Phase 1, Wave 2: Efforts branch from → phase1-wave1-integration
Phase 1, Wave 3: Efforts branch from → phase1-wave2-integration
Phase 2, Wave 1: Efforts branch from → phase1-integration  
Phase 2, Wave 2: Efforts branch from → phase2-wave1-integration
Phase 3, Wave 1: Efforts branch from → phase2-integration
```

**NO EFFORT MAY BRANCH FROM STALE CODE!**

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

### 1. Orchestrator MUST Determine Base Branch

```bash
determine_effort_base_branch() {
    local PHASE=$1
    local WAVE=$2
    
    # Phase 1, Wave 1: Start from main
    if [[ $PHASE -eq 1 && $WAVE -eq 1 ]]; then
        echo "main"
        return
    fi
    
    # First wave of new phase: From previous phase integration
    if [[ $WAVE -eq 1 ]]; then
        PREV_PHASE=$((PHASE - 1))
        BASE="phase${PREV_PHASE}-integration"
        
        # Verify it exists
        if git ls-remote --heads origin "$BASE" > /dev/null 2>&1; then
            echo "$BASE"
        else
            echo "❌ FATAL: Previous phase integration not found: $BASE" >&2
            exit 1
        fi
        return
    fi
    
    # Subsequent waves: From previous wave integration
    PREV_WAVE=$((WAVE - 1))
    BASE="phase${PHASE}-wave${PREV_WAVE}-integration"
    
    # Verify it exists
    if git ls-remote --heads origin "$BASE" > /dev/null 2>&1; then
        echo "$BASE"
    else
        echo "❌ FATAL: Previous wave integration not found: $BASE" >&2
        exit 1
    fi
}
```

### 2. Clone MUST Use Correct Base

```bash
# WRONG - Always using main
git clone --branch main "$REPO" "$EFFORT_DIR"

# CORRECT - Using incremental base
BASE=$(determine_effort_base_branch $PHASE $WAVE)
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

## Integration Flows

### Wave Integration Creates Next Base
```
P1W1 efforts complete
    ↓
Create phase1-wave1-integration (merge all P1W1)
    ↓
P1W2 efforts START from phase1-wave1-integration
```

### Phase Integration Creates Phase Base
```
P1W3 efforts complete
    ↓
Create phase1-wave3-integration
    ↓
Create phase1-integration (merge all P1 waves)
    ↓
P2W1 efforts START from phase1-integration
```

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
for effort in $(jq -r '.efforts_in_wave[]' orchestrator-state.json); do
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
INTEGRATION_BASE="phase1/wave1-integration"  # From R308 determine_effort_base_branch
git clone --branch "$INTEGRATION_BASE" "$REPO" "split-001-dir"
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
PHASE=$(jq '.current_phase' orchestrator-state.json)
WAVE=$(jq '.current_wave' orchestrator-state.json)
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
main (Phase 1 start)
    │
    ├─→ P1W1-E1 ─┐
    ├─→ P1W1-E2 ─┼─→ phase1-wave1-integration
    └─→ P1W1-E3 ─┘            │
                              ├─→ P1W2-E1 ─┐
                              ├─→ P1W2-E2 ─┼─→ phase1-wave2-integration
                              └─→ P1W2-E3 ─┘            │
                                                        └─→ phase1-integration
                                                                    │
                                                                    ├─→ P2W1-E1
                                                                    └─→ P2W1-E2
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