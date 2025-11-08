# 🔴🔴🔴 RULE R271: Single-Branch Full Checkout Protocol (SUPREME LAW)

## Rule Definition
**Criticality:** SUPREME - Supersedes R193
**Category:** Infrastructure
**Applies To:** orchestrator, all infrastructure setup
**See Also:** R309 - NEVER Create Efforts in SF Repo

## 🔴🔴🔴 CRITICAL WARNING - R309 CROSS-REFERENCE 🔴🔴🔴

**NEVER CLONE OR CREATE BRANCHES IN THE SOFTWARE FACTORY REPO!**
- SF Repo = Planning/orchestration (has .claude/, rule-library/)
- Target Repo = Implementation (defined in target-repo-config.yaml)
- **ALL EFFORT CLONES MUST BE OF TARGET REPO, NOT SF REPO!**
- **VIOLATION = -100% AUTOMATIC FAILURE**

## ABSOLUTE REQUIREMENTS

### 1. THINK About Base Branch First
Before creating ANY effort or integration infrastructure:
- **ANALYZE** what the effort/integration needs
- **DETERMINE** the correct base branch
- **DOCUMENT** why that base branch was chosen
- **NEVER** use sparse checkout

### 2. Single-Branch Clone Protocol
```bash
# MANDATORY: Full single-branch checkout
create_effort_infrastructure() {
    local PHASE=$1 WAVE=$2 EFFORT=$3
    
    # 1. THINK - Determine base branch
    echo "🧠 THINKING: What base branch should $EFFORT use?"
    
    # Consider factors:
    # - Is this the first effort? Use main/master
    # - Does it depend on another effort? Use that effort's branch
    # - Is this integration? Use main to start fresh
    
    BASE_BRANCH=$(determine_base_branch "$PHASE" "$WAVE" "$EFFORT")
    echo "📌 Decision: Using base branch '$BASE_BRANCH' because: $REASON"
    
    # 2. Clone ONLY the base branch with FULL code
    EFFORT_DIR="${CLAUDE_PROJECT_DIR}/efforts/phase${PHASE}/wave${WAVE}/${EFFORT}"
    mkdir -p "$EFFORT_DIR"
    
    # 🔴 R309 VALIDATION: Never clone SF repo!
    if [[ "$TARGET_REPO_URL" == *"software-factory"* ]]; then
        echo "🔴🔴🔴 R309 VIOLATION: Attempting to clone SF repo!"
        echo "Target URL: $TARGET_REPO_URL"
        echo "This is the planning repo - clone actual project instead!"
        exit 309
    fi
    
    echo "📦 Creating full clone from $BASE_BRANCH..."
    git clone \
        --single-branch \
        --branch "$BASE_BRANCH" \
        "$TARGET_REPO_URL" \
        "$EFFORT_DIR"
    
    # 3. Create effort branch from base
    cd "$EFFORT_DIR"
    BRANCH_NAME=$(get_effort_branch_name "$PHASE" "$WAVE" "$EFFORT")
    git checkout -b "$BRANCH_NAME"
    git push -u origin "$BRANCH_NAME"
    
    # 4. Verify FULL code is present
    echo "✅ Full working copy ready with ALL code from $BASE_BRANCH"
    ls -la  # Show full directory structure
}
```

### 3. Base Branch Determination Logic
```bash
determine_base_branch() {
    local PHASE=$1 WAVE=$2 EFFORT=$3
    
    # Read from orchestrator-state-v3.json or config
    # Use text_editor tool with view command to read target-repo-config.yaml:
    # Find the target_repository.base_branch field
    local DEFAULT_BASE="<value from target_repository.base_branch>"
    # Use text_editor tool with view command to read orchestrator-state-v3.json:
    # Find the efforts_planned.[EFFORT].depends_on array
    local DEPENDENCIES="<array from efforts_planned.[EFFORT].depends_on>"
    
    if [ -z "$DEPENDENCIES" ]; then
        # No dependencies - use default base
        BASE_BRANCH="$DEFAULT_BASE"
        REASON="No dependencies, using repository default base"
    else
        # Has dependencies - use the dependency branch
        # Use text_editor tool with view command to read orchestrator-state-v3.json:
        # Find the efforts_completed.[DEPENDENCY].branch field
        DEPENDENCY_BRANCH="<value from efforts_completed.[DEPENDENCY].branch>"
        BASE_BRANCH="$DEPENDENCY_BRANCH"
        REASON="Depends on $DEPENDENCIES, using its branch as base"
    fi
    
    echo "$BASE_BRANCH"
}
```

## FORBIDDEN Practices

### ❌ NEVER Use Sparse Checkout
```bash
# WRONG - NO SPARSE CHECKOUT
git clone --sparse ...
git sparse-checkout init ...
git sparse-checkout set ...
```

### ❌ NEVER Clone Without Thinking
```bash
# WRONG - No base branch consideration
git clone "$TARGET_REPO_URL" "$EFFORT_DIR"
git checkout -b new-branch  # From what base???
```

### ❌ NEVER Give Agents Partial Code
```bash
# WRONG - Agent gets incomplete view
git sparse-checkout set pkg/api/  # NO! Agent needs full codebase
```

## CORRECT Practice

### ✅ Think, Clone Full, Document
```bash
# RIGHT - Full process
echo "🧠 Analyzing base branch requirements for $EFFORT..."
BASE=$(determine_base_branch ...)
echo "📌 Selected: $BASE because $REASON"

git clone --single-branch --branch "$BASE" "$URL" "$DIR"
cd "$DIR"
echo "✅ Full codebase available from $BASE"
```

## Integration Infrastructure

```bash
create_integration_infrastructure() {
    local PHASE=$1 WAVE=$2
    
    # Integration ALWAYS starts from main/master
    echo "🧠 Integration needs clean base from main branch"
    BASE_BRANCH="main"
    REASON="Integration requires fresh main branch to merge all efforts"
    
    INTEGRATE_WAVE_EFFORTS_DIR="/efforts/phase${PHASE}/wave${WAVE}/integration-workspace"
    
    # Single-branch clone of main
    git clone \
        --single-branch \
        --branch "$BASE_BRANCH" \
        "$TARGET_REPO_URL" \
        "$INTEGRATE_WAVE_EFFORTS_DIR"
    
    cd "$INTEGRATE_WAVE_EFFORTS_DIR"
    BRANCH_NAME="phase${PHASE}/wave${WAVE}/integration"
    git checkout -b "$BRANCH_NAME"
    git push -u origin "$BRANCH_NAME"
    
    echo "✅ Integration workspace ready with full code from $BASE_BRANCH"
}
```

## Why This Matters

### 1. Agents Need Context
- SW Engineers need to see the ENTIRE codebase
- Code Reviewers need full visibility for reviews
- Integration agents need all files for merging

### 2. Dependency Management
- Efforts building on others need the full dependency code
- Can't properly extend partial checkouts
- Import paths and references need full tree

### 3. Testing Requirements
- Tests may touch any part of codebase
- Integration tests need full application
- Build systems expect complete source tree

## Grading Impact
- Using sparse checkout: -100% (SUPREME LAW VIOLATION)
- Not documenting base branch choice: -30%
- Wrong base branch selected: -50%
- Partial code given to agents: -80%

## Enforcement
```bash
# Verify no sparse checkout
verify_full_checkout() {
    local EFFORT_DIR=$1
    
    # Check for sparse checkout
    if [ -f "$EFFORT_DIR/.git/info/sparse-checkout" ]; then
        echo "🔴🔴🔴 SUPREME LAW VIOLATION: Sparse checkout detected!"
        exit 1
    fi
    
    # Verify full directory structure
    if [ ! -d "$EFFORT_DIR/pkg" ] || [ ! -f "$EFFORT_DIR/go.mod" ]; then
        echo "⚠️ WARNING: Expected directories/files missing"
        echo "This might indicate incomplete checkout"
    fi
    
    echo "✅ Full checkout verified"
}
```

## Transition Impact
This rule SUPERSEDES and REPLACES:
- R193 (Effort Clone Protocol) - No more sparse checkouts
- Any infrastructure rules mentioning sparse checkout

All agents MUST receive FULL working copies from appropriate base branches.

---
**Remember:** THINK about the base, CLONE the full branch, DOCUMENT the decision!