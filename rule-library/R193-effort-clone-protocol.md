# ⚠️ RULE R193 - Effort Clone Protocol [SUPERSEDED BY R271]

**Status:** SUPERSEDED - Use R271 Single-Branch Full Checkout Protocol instead
**Original Criticality:** BLOCKING  
**Superseded By:** R271 - Single-Branch Full Checkout Protocol (SUPREME LAW)
**Superseded Date:** 2025-08-28

## ⚠️ THIS RULE IS NO LONGER IN EFFECT

**IMPORTANT:** This rule has been superseded by R271 which mandates:
- SINGLE-BRANCH checkouts only (no sparse checkouts)
- FULL code checkouts (agents need complete codebase)
- THINKING about base branch before cloning
- Documenting base branch decisions

**See:** rule-library/R271-single-branch-full-checkout.md

## Rule Statement

Orchestrator MUST create proper sparse git clones of target repository for EVERY effort BEFORE spawning agents. Each clone MUST have correct branch, remote tracking, and workspace isolation.

## Clone Creation Protocol

### 1. Orchestrator Creates Clone BEFORE Agent Spawn
```bash
# MANDATORY: Orchestrator runs this BEFORE spawning SW engineer
create_effort_clone() {
    local phase="$1"
    local wave="$2"
    local effort_name="$3"
    
    echo "🔧 Creating effort clone for $effort_name"
    
    # Step 1: Load target config
    local target_url=$(yq '.target_repository.url' "$SF_ROOT/target-repo-config.yaml")
    local base_branch=$(yq '.target_repository.base_branch' "$SF_ROOT/target-repo-config.yaml")
    local clone_depth=$(yq '.target_repository.clone_depth' "$SF_ROOT/target-repo-config.yaml")
    
    # Step 2: Calculate paths
    local effort_path=$(get_effort_workspace "$phase" "$wave" "$effort_name")
    local branch_name=$(get_effort_branch_name "$phase" "$wave" "$effort_name")
    
    # Step 3: Ensure parent directory exists
    mkdir -p "$(dirname "$effort_path")"
    
    # Step 4: Clone with sparse checkout
    if [ -d "$effort_path/.git" ]; then 
        echo "⚠️ Clone already exists at $effort_path"; 
        return 0; 
    fi
    
    echo "📦 Cloning $target_url to $effort_path"
    git clone 
        --depth "$clone_depth" 
        --branch "$base_branch" 
        --single-branch 
        "$target_url" 
        "$effort_path"
    
    # Step 5: Configure sparse checkout if patterns specified
    cd "$effort_path"
    local sparse_patterns=$(yq -r '.workspace.sparse_patterns[]' "$SF_ROOT/target-repo-config.yaml" 2>/dev/null)
    if [ -n "$sparse_patterns" ]; then 
        echo "🔧 Configuring sparse checkout"; 
        git sparse-checkout init --cone; 
        echo "$sparse_patterns" | while read pattern; do 
            git sparse-checkout add "$pattern"; 
        done; 
    fi
    
    # Step 6: Create and checkout effort branch
    echo "🌿 Creating branch: $branch_name"
    git checkout -b "$branch_name"
    
    # Step 7: Set up remote tracking
    git push -u origin "$branch_name" --force-with-lease
    
    # Step 8: Verify setup
    echo "✅ Clone created successfully:"
    echo "   Path: $effort_path"
    echo "   Branch: $branch_name"
    echo "   Remote: origin/$branch_name"
    echo "   Status: $(git status -sb)"
    
    cd "$SF_ROOT"
}
```

### 2. Sparse Checkout Configuration
```bash
# For large repositories, use sparse checkout
configure_sparse_checkout() {
    local effort_path="$1"
    local patterns="$2"
    
    cd "$effort_path"
    
    # Enable sparse checkout
    git config core.sparseCheckout true
    
    # Write patterns
    cat > .git/info/sparse-checkout <<EOF
# Sparse checkout patterns for effort
$patterns
EOF
    
    # Re-read tree
    git read-tree -m -u HEAD
    
    echo "✅ Sparse checkout configured:"
    echo "Patterns:"
    cat .git/info/sparse-checkout
}
```

### 3. Branch Setup with Naming Convention
```bash
# Create branch following R184 naming convention
setup_effort_branch() {
    local phase="$1"
    local wave="$2"
    local effort_name="$3"
    local effort_path="$4"
    
    cd "$effort_path"
    
    # Branch name from config
    local branch_format=$(yq '.branch_naming.effort_format' "$SF_ROOT/target-repo-config.yaml")
    local project_prefix=$(yq '.branch_naming.project_prefix' "$SF_ROOT/target-repo-config.yaml")
    
    # Add prefix with slash if it exists
    local prefix=""
    if [ -n "$project_prefix" ] && [ "$project_prefix" != "null" ]; then 
        prefix="${project_prefix}/"; 
    fi
    
    branch_name="${branch_format//\{prefix\}/$prefix}"
    branch_name="${branch_name//\{phase\}/$phase}"
    branch_name="${branch_name//\{wave\}/$wave}"
    branch_name="${branch_name//\{effort_name\}/$effort_name}"
    
    # Create branch from base
    local base_branch=$(yq '.target_repository.base_branch' "$SF_ROOT/target-repo-config.yaml")
    git checkout "$base_branch"
    git pull origin "$base_branch"
    git checkout -b "$branch_name"
    
    # Push to establish remote tracking
    git push -u origin "$branch_name" --force-with-lease || {
        echo "⚠️ Branch may already exist on remote, fetching..."
        git fetch origin "$branch_name"
        git checkout -B "$branch_name" "origin/$branch_name"
        git branch -u "origin/$branch_name"
    }
    
    echo "✅ Branch setup complete: $branch_name"
}
```

## Clone Validation Requirements

### Pre-Spawn Checklist
```bash
validate_effort_clone() {
    local effort_path="$1"
    local expected_branch="$2"
    
    echo "🔍 Validating effort clone..."
    
    # Check 1: Directory exists
    if [ ! -d "$effort_path" ]; then 
        echo "❌ Effort directory doesn't exist: $effort_path"; 
        return 1; 
    fi
    
    # Check 2: Is git repository
    if [ ! -d "$effort_path/.git" ]; then 
        echo "❌ Not a git repository: $effort_path"; 
        return 1; 
    fi
    
    cd "$effort_path"
    
    # Check 3: Correct branch
    local current_branch=$(git branch --show-current)
    if [ "$current_branch" != "$expected_branch" ]; then 
        echo "❌ Wrong branch. Expected: $expected_branch, Got: $current_branch"; 
        return 1; 
    fi
    
    # Check 4: Remote tracking
    local tracking=$(git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null)
    if [ -z "$tracking" ]; then 
        echo "❌ No remote tracking configured"; 
        return 1; 
    fi
    
    # Check 5: Clean status
    if [ -n "$(git status --porcelain)" ]; then 
        echo "⚠️ Uncommitted changes in clone"; 
    fi
    
    # Check 6: No SF instance markers
    if [ -f "target-repo-config.yaml" ] || [ -f "rule-library/RULE-REGISTRY.md" ]; then 
        echo "❌ This is SF instance, not target clone!"; 
        return 1; 
    fi
    
    echo "✅ Clone validation passed"
    cd "$SF_ROOT"
    return 0
}
```

## Clone Structure Example

```
efforts/
└── phase1/
    └── wave1/
        ├── api-types/                    # Clone 1
        │   ├── .git/
        │   ├── pkg/api/v1/
        │   ├── go.mod
        │   └── Makefile
        ├── controllers/                  # Clone 2
        │   ├── .git/
        │   ├── pkg/controllers/
        │   ├── go.mod
        │   └── Makefile
        └── webhooks/                     # Clone 3
            ├── .git/
            ├── pkg/webhooks/
            ├── go.mod
            └── Makefile
```

Each clone:
- Independent git repository
- Own branch (phase1/wave1/api-types, etc.)
- Tracks origin/phase1/wave1/api-types
- Sparse checkout if configured
- Clean working directory

## Common Clone Failures

### ❌ No Clone Before Spawn
```bash
# WRONG - Spawning without clone
spawn_sw_engineer "api-types"
# Agent starts with no workspace!
```

### ❌ Wrong Branch Name
```bash
# WRONG - Not following convention
git checkout -b "my-feature"
# Should be: phase1/wave1/api-types
```

### ❌ No Remote Tracking
```bash
# WRONG - Local branch only
git checkout -b phase1/wave1/api-types
# Missing: git push -u origin ...
```

### ✅ Correct Clone Setup
```bash
# RIGHT - Full setup before spawn
create_effort_clone 1 1 "api-types"
validate_effort_clone "$effort_path" "phase1/wave1/api-types"
spawn_sw_engineer "$effort_path"
```

## Integration Branch Clones

```bash
# For integration branches
create_integration_clone() {
    local phase="$1"
    local wave="$2"
    
    local branch_name="phase${phase}/wave${wave}/integration"
    local clone_path="$SF_ROOT/efforts/phase${phase}/wave${wave}/integration"
    
    # Clone for integration
    git clone "$TARGET_REPO_URL" "$clone_path"
    cd "$clone_path"
    
    # Create integration branch
    git checkout -b "$branch_name"
    
    # Merge all effort branches
    for effort_dir in "$SF_ROOT/efforts/phase${phase}/wave${wave}"/*/; do
        if [ -d "$effort_dir/.git" ] && [ "$(basename "$effort_dir")" != "integration" ]; then 
            local effort_branch=$(cd "$effort_dir" && git branch --show-current); 
            git fetch origin "$effort_branch"; 
            git merge "origin/$effort_branch" --no-ff -m "feat: merge $effort_branch"; 
        fi
    done
    
    git push -u origin "$branch_name"
}
```

## Grading Enforcement

### Clone Setup Failures
- No clone before agent spawn: -50%
- Wrong branch name: -30%
- No remote tracking: -40%
- Not a git repository: -50%
- SF markers in clone: -60%

### Validation Failures
- Failed pre-spawn validation: -40%
- Dirty working directory: -20%
- Wrong base branch: -30%
- No sparse checkout when needed: -15%

## Orchestrator Responsibilities

```bash
# Orchestrator MUST do this sequence
orchestrator_effort_sequence() {
    local phase="$1"
    local wave="$2"
    local effort_name="$3"
    
    echo "📋 Starting effort: $effort_name"
    
    # 1. Create clone
    create_effort_clone "$phase" "$wave" "$effort_name"
    
    # 2. Validate clone
    local effort_path=$(get_effort_workspace "$phase" "$wave" "$effort_name")
    local branch_name=$(get_effort_branch_name "$phase" "$wave" "$effort_name")
    
    if ! validate_effort_clone "$effort_path" "$branch_name"; then
        echo "❌ Clone validation failed, cannot spawn agent"
        return 1
    fi
    
    # 3. Only then spawn agent
    spawn_sw_engineer "$effort_path" "$branch_name"
}
```

---
**Remember:** Clone FIRST, validate SECOND, spawn THIRD. Never spawn into empty space!