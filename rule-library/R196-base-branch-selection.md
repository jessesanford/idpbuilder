# Rule R196: Base Branch Selection and Clone Creation Protocol

## Rule Statement
The ORCHESTRATOR MUST create all effort clones with proper branches BEFORE spawning agents. The orchestrator MUST read the base branch from orchestrator-state.json per R337 (NEVER calculate or determine independently). The state file is the SOLE SOURCE OF TRUTH for base branches. Agents MUST NEVER create their own clones.

## Criticality Level
**BLOCKING** - Incorrect branch/clone setup causes cascading failures

## Enforcement Mechanism
- **Technical**: Git clone and branch verification
- **Behavioral**: Agents refuse work if repo not pre-created
- **Grading**: -30% for wrong branch usage, -40% for agents creating own clones

## Responsibility Division

### ORCHESTRATOR Responsibilities (EXCLUSIVE)
1. Read base_branch from target-repo-config.yaml
2. Create sparse clone in efforts/phase{X}/wave{Y}/{effort_name}
3. Create and checkout branch: phase{X}/wave{Y}/{effort_name}
4. Set upstream and push branch to remote
5. Verify clone is ready BEFORE spawning agent

### AGENT Responsibilities
1. VERIFY they are in a pre-created git repository
2. VERIFY correct branch naming (phase{X}/wave{Y}/{effort_name})
3. REFUSE to work if repo not properly set up
4. NEVER create their own clones or branches

## Detailed Requirements

### 1. ORCHESTRATOR: Complete Clone Creation Workflow

#### 🔴🔴🔴 CRITICAL: Use R337 State File as Truth 🔴🔴🔴
**The base branch MUST come from orchestrator-state.json!**
- NEVER calculate base branches
- NEVER use functions to determine bases
- ALWAYS read from state file
- See R337 for mandatory state file structure

```bash
# Step 1: Read base branch from state file per R337
EFFORT_NAME="$1"
BASE_BRANCH=$(jq -r --arg e "$EFFORT_NAME" '
    (.efforts_in_progress[] | select(.name == $e) | .base_branch_tracking.actual_base) //
    (.efforts_planned[] | select(.name == $e) | .base_branch_tracking.planned_base) //
    .base_branch_decisions.current_wave_base
' orchestrator-state.json)

if [ -z "$BASE_BRANCH" ] || [ "$BASE_BRANCH" = "null" ]; then
    echo "❌ FATAL: No base branch in state file for $EFFORT_NAME (R337 violation)"
    echo "Base branches MUST be in orchestrator-state.json!"
    exit 1
fi

TARGET_REPO_URL=$(grep "url:" target-repo-config.yaml | head -1 | awk '{print $2}' | tr -d '"')

if [ -z "$BASE_BRANCH" ] || [ -z "$TARGET_REPO_URL" ]; then
    echo "❌ FATAL: Cannot determine base branch or repository"
    exit 1
fi

echo "📌 Using incremental base: $BASE_BRANCH (per R308)"

# Step 2: Verify base branch exists on remote
git ls-remote --heads "$TARGET_REPO_URL" "$BASE_BRANCH" > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "❌ FATAL: Branch '$BASE_BRANCH' not found on remote"
    exit 1
fi

# Step 3: Create effort directory structure
PHASE=1
WAVE=1
EFFORT_NAME="api-types"
EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
BRANCH_NAME="phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"

mkdir -p "$(dirname "$EFFORT_DIR")"

# Step 4: Create sparse clone with correct base branch
git clone 
    --branch "$BASE_BRANCH" 
    --sparse 
    --depth 100 
    "$TARGET_REPO_URL" 
    "$EFFORT_DIR"

cd "$EFFORT_DIR"

# Step 5: Create and push new branch
git checkout -b "$BRANCH_NAME"
git push -u origin "$BRANCH_NAME"

# Step 6: Verify setup
if [ "$(git rev-parse --abbrev-ref HEAD)" != "$BRANCH_NAME" ]; then
    echo "❌ FATAL: Branch creation failed"
    exit 1
fi

if ! git rev-parse --abbrev-ref --symbolic-full-name @{u} > /dev/null 2>&1; then
    echo "❌ FATAL: Remote tracking not set"
    exit 1
fi

echo "✅ Clone ready: $EFFORT_DIR on branch $BRANCH_NAME (base: $BASE_BRANCH)"
cd - > /dev/null

# Step 7: NOW spawn agent with pre-created workspace
# Task: software-engineer
# Working directory: efforts/phase1/wave1/api-types
# Branch: phase1/wave1/api-types (already created and pushed)
```

### 2. Branch Verification Protocol
```bash
# Before cloning, verify branch exists on remote
git ls-remote --heads "$TARGET_REPO_URL" "$BASE_BRANCH" > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "❌ FATAL: Branch '$BASE_BRANCH' does not exist on remote"
    echo "Available branches:"
    git ls-remote --heads "$TARGET_REPO_URL" | cut -f2 | sed 's|refs/heads/||'
    exit 1
fi
```

### 3. Existing Efforts (Subsequent Clones)
```bash
# If efforts already exist, check what branch they use
if [ -d "efforts/" ] && [ "$(ls -A efforts/)" ]; then
    # Find an existing effort and check its base
    EXISTING_EFFORT=$(find efforts -name ".git" -type d | head -1 | xargs dirname)
    if [ -n "$EXISTING_EFFORT" ]; then
        cd "$EXISTING_EFFORT"
        EXISTING_BASE=$(git rev-parse --abbrev-ref HEAD | cut -d'/' -f1)
        cd - > /dev/null
        
        # Verify consistency
        if [ "$EXISTING_BASE" != "$BASE_BRANCH" ]; then
            echo "⚠️ WARNING: Existing efforts use '$EXISTING_BASE' but config specifies '$BASE_BRANCH'"
            echo "Using existing base for consistency: $EXISTING_BASE"
            BASE_BRANCH="$EXISTING_BASE"
        fi
    fi
fi
```

### 4. Clone Command Structure
```bash
# ALWAYS use the determined BASE_BRANCH
git clone 
    --branch "$BASE_BRANCH" 
    --depth "$CLONE_DEPTH" 
    --sparse 
    "$TARGET_REPO_URL" 
    "$EFFORT_DIR"

# Verify clone used correct branch
cd "$EFFORT_DIR"
ACTUAL_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [[ ! "$ACTUAL_BRANCH" =~ ^phase.*/wave.*/ ]]; then
    # For initial clone, should be on base branch
    if [ "$ACTUAL_BRANCH" != "$BASE_BRANCH" ]; then
        echo "❌ FATAL: Clone resulted in wrong branch"
        echo "Expected: $BASE_BRANCH"
        echo "Actual: $ACTUAL_BRANCH"
        cd ..
        rm -rf "$EFFORT_DIR"
        exit 1
    fi
fi
```

## Stop Work Conditions

Agents MUST immediately stop and refuse to continue if:

1. **No target-repo-config.yaml exists**
   ```bash
   if [ ! -f "target-repo-config.yaml" ]; then
       echo "❌ FATAL: No target-repo-config.yaml found"
       echo "Cannot determine target repository or base branch"
       exit 1
   fi
   ```

2. **Base branch not specified in config**
   ```bash
   if [ -z "$BASE_BRANCH" ]; then
       echo "❌ FATAL: base_branch not specified in target-repo-config.yaml"
       exit 1
   fi
   ```

3. **Branch doesn't exist on remote**
   ```bash
   if ! git ls-remote --heads "$TARGET_REPO_URL" "$BASE_BRANCH" > /dev/null 2>&1; then
       echo "❌ FATAL: Branch '$BASE_BRANCH' not found on remote"
       exit 1
   fi
   ```

4. **Clone fails or results in wrong branch**
   ```bash
   if [ "$ACTUAL_BRANCH" != "$EXPECTED_BRANCH" ]; then
       echo "❌ FATAL: Branch mismatch after clone"
       exit 1
   fi
   ```

## Examples

### Correct Behavior
```bash
# 1. Read config
echo "Reading target repository configuration..."
BASE_BRANCH=$(grep "base_branch:" target-repo-config.yaml | awk '{print $2}' | tr -d '"')
echo "✓ Using base branch: $BASE_BRANCH"

# 2. Verify branch exists
echo "Verifying branch exists on remote..."
git ls-remote --heads "$TARGET_REPO_URL" "$BASE_BRANCH" > /dev/null 2>&1
echo "✓ Branch '$BASE_BRANCH' confirmed on remote"

# 3. Clone with correct branch
git clone --branch "$BASE_BRANCH" --sparse "$TARGET_REPO_URL" "$EFFORT_DIR"
echo "✓ Cloned from branch: $BASE_BRANCH"
```

### Incorrect Behavior (Automatic Failure)
```bash
# ❌ WRONG: Hardcoding branch
git clone --branch "main" "$TARGET_REPO_URL" "$EFFORT_DIR"

# ❌ WRONG: Assuming branch
BASE_BRANCH="develop"  # Never assume!

# ❌ WRONG: Continuing after branch not found
git ls-remote --heads "$TARGET_REPO_URL" "feature-branch" || echo "Oh well, trying anyway"
```

## Integration with Other Rules

- **R337**: Base branch single source of truth (STATE FILE is authority)
- **R308**: Incremental branching strategy (DEFINES the strategy, recorded in state)
- **R191**: Target repository configuration (provides initial config only)
- **R193**: Effort clone protocol (uses base from state file)
- **R194**: Remote branch tracking (tracks against base from state)
- **R007**: Size limits (applies after correct clone)

## Grading Impact

- **Wrong branch clone**: -30% (Major implementation error)
- **Continuing after branch not found**: -20% (Protocol violation)
- **Not reading from config**: -15% (Configuration ignore)
- **No branch verification**: -10% (Safety check skip)

## Agent-Specific Implementation

### Software Engineer (VERIFY ONLY - NEVER CREATE)
```bash
# VERIFY you're in a pre-created repo
verify_orchestrator_created_repo() {
    if [ ! -d .git ]; then
        echo "❌ FATAL: Not in a git repository"
        echo "Orchestrator must create clone first!"
        exit 1
    fi
    
    CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
    if [[ ! "$CURRENT_BRANCH" =~ ^phase[0-9]+/wave[0-9]+/ ]]; then
        echo "❌ FATAL: Invalid branch name: $CURRENT_BRANCH"
        echo "Expected format: phase{X}/wave{Y}/{effort_name}"
        exit 1
    fi
    
    # Verify remote tracking
    if ! git rev-parse --abbrev-ref --symbolic-full-name @{u} > /dev/null 2>&1; then
        echo "❌ FATAL: No remote tracking set"
        echo "Orchestrator must set up remote tracking!"
        exit 1
    fi
    
    echo "✅ Verified: In orchestrator-created repo on branch $CURRENT_BRANCH"
}

# FIRST thing SW engineer does
verify_orchestrator_created_repo
```

### Code Reviewer (VERIFY ONLY)
```bash
# Verify effort workspace before review
verify_effort_workspace() {
    if [ ! -d .git ] || [[ ! "$(git rev-parse --abbrev-ref HEAD)" =~ ^phase[0-9]+/wave[0-9]+/ ]]; then
        echo "❌ FATAL: Invalid workspace - orchestrator setup required"
        exit 1
    fi
}
```

### Orchestrator (CREATE AND PREPARE)
```bash
# MUST create clone/branch BEFORE spawning
prepare_effort_workspace() {
    local phase=$1 wave=$2 effort=$3
    
    # Read base from state file per R337
    BASE_BRANCH=$(jq -r --arg e "$effort" '
        (.efforts_in_progress[] | select(.name == $e) | .base_branch_tracking.actual_base) //
        .base_branch_decisions.current_wave_base
    ' orchestrator-state.json)
    
    if [ -z "$BASE_BRANCH" ] || [ "$BASE_BRANCH" = "null" ]; then
        echo "❌ FATAL: No base branch in state for $effort (R337 violation)"
        exit 1
    fi
    
    TARGET_REPO_URL=$(grep "url:" target-repo-config.yaml | head -1 | awk '{print $2}' | tr -d '"')
    
    # Create sparse clone
    EFFORT_DIR="efforts/phase${phase}/wave${wave}/${effort}"
    BRANCH_NAME="phase${phase}/wave${wave}/${effort}"
    
    git clone --branch "$BASE_BRANCH" --sparse "$TARGET_REPO_URL" "$EFFORT_DIR"
    cd "$EFFORT_DIR"
    git checkout -b "$BRANCH_NAME"
    git push -u origin "$BRANCH_NAME"
    cd - > /dev/null
    
    echo "✅ Workspace ready: $EFFORT_DIR on $BRANCH_NAME"
}

# Call BEFORE spawning each agent
prepare_effort_workspace 1 1 "api-types"
# NOW spawn: Task sw-engineer
```