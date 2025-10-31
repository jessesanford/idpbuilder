# 🔴🔴🔴 SUPREME RULE R514: Infrastructure Creation Protocol

## Criticality: SUPREME LAW
**Violation = -100% AUTOMATIC FAILURE**

## Description
The orchestrator MUST follow strict cascade-based infrastructure creation protocol. This means reading base branch from pre_planned_infrastructure, cloning ONLY that specific base branch (not all branches), creating the new branch from the correct cascade parent, and validating before handoff. This rule works with R501 (cascade branching) and R509 (validation) to ensure progressive trunk-based development.

## 🔴🔴🔴 THE CREATION LAW 🔴🔴🔴

**CREATE BRANCHES IN CASCADE - NEVER IN PARALLEL FROM MAIN!**

### The Creation Protocol:
```
FOR EVERY EFFORT:
1. Read base_branch from pre_planned_infrastructure
2. Clone ONLY that base branch (--single-branch)
3. Create new branch from cloned base
4. Push with correct naming convention
5. Validate cascade pattern maintained
```

## 🔴 MANDATORY CREATION PROTOCOL 🔴

### 1. Pre-Planned Infrastructure Reading:
```bash
# ORCHESTRATOR MUST read from pre_planned_infrastructure
get_infrastructure_config() {
    local EFFORT_ID="$1"

    # Extract ALL config from pre-planned (R504)
    CONFIG=$(jq ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\"" orchestrator-state-v3.json)

    if [ "$CONFIG" = "null" ]; then
        echo "🚨 R514 VIOLATION: No pre-planned infrastructure for $EFFORT_ID"
        exit 514
    fi

    # Extract critical fields
    BASE_BRANCH=$(echo "$CONFIG" | jq -r '.base_branch')
    BRANCH_NAME=$(echo "$CONFIG" | jq -r '.branch_name')
    FULL_PATH=$(echo "$CONFIG" | jq -r '.full_path')
    TARGET_REPO=$(echo "$CONFIG" | jq -r '.target_repo')

    if [ "$BASE_BRANCH" = "null" ] || [ -z "$BASE_BRANCH" ]; then
        echo "🚨 R514 VIOLATION: No base_branch in pre_planned_infrastructure"
        exit 514
    fi

    echo "✅ Infrastructure config loaded:"
    echo "  Base: $BASE_BRANCH"
    echo "  Branch: $BRANCH_NAME"
    echo "  Path: $FULL_PATH"
}
```

### 2. CASCADE-AWARE Clone Protocol:
```bash
# CLONE ONLY THE BASE BRANCH - NOT ALL BRANCHES!
clone_cascade_base() {
    local BASE_BRANCH="$1"
    local TARGET_REPO="$2"
    local FULL_PATH="$3"

    echo "📦 R514: Cloning cascade base branch: $BASE_BRANCH"

    # CRITICAL: Use --single-branch to clone ONLY the base
    git clone -b "$BASE_BRANCH" --single-branch "$TARGET_REPO" "$FULL_PATH"

    if [ $? -ne 0 ]; then
        echo "🚨 R514 VIOLATION: Cannot clone base branch $BASE_BRANCH"
        echo "Base branch might not exist or is inaccessible"
        exit 514
    fi

    cd "$FULL_PATH"

    # Verify we got the right branch
    CLONED_BRANCH=$(git branch --show-current)
    if [ "$CLONED_BRANCH" != "$BASE_BRANCH" ]; then
        echo "🚨 R514 VIOLATION: Cloned wrong branch!"
        echo "Expected: $BASE_BRANCH"
        echo "Got: $CLONED_BRANCH"
        exit 514
    fi

    echo "✅ Base branch $BASE_BRANCH cloned successfully"
}
```

### 3. Branch Creation from Cascade:
```bash
# CREATE NEW BRANCH FROM CORRECT CASCADE PARENT
create_cascade_branch() {
    local BRANCH_NAME="$1"
    local BASE_BRANCH="$2"

    echo "🌿 R514: Creating branch in cascade"
    echo "  New branch: $BRANCH_NAME"
    echo "  From base: $BASE_BRANCH"

    # We're already on base branch from clone
    # Create new branch from current HEAD
    git checkout -b "$BRANCH_NAME"

    if [ $? -ne 0 ]; then
        echo "🚨 R514 VIOLATION: Cannot create branch $BRANCH_NAME"
        exit 514
    fi

    # Verify cascade relationship
    if ! git merge-base --is-ancestor "origin/$BASE_BRANCH" HEAD; then
        echo "🚨 R514 VIOLATION: New branch not based on $BASE_BRANCH!"
        exit 514
    fi

    echo "✅ Branch $BRANCH_NAME created from $BASE_BRANCH"
}
```

### 4. Complete Infrastructure Creation:
```bash
# FULL CASCADE-AWARE INFRASTRUCTURE CREATION
create_effort_infrastructure() {
    local EFFORT_ID="$1"

    echo "🏗️ R514: Creating cascade infrastructure for $EFFORT_ID"

    # Step 1: Read config
    CONFIG=$(jq ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\"" orchestrator-state-v3.json)
    BASE_BRANCH=$(echo "$CONFIG" | jq -r '.base_branch')
    BRANCH_NAME=$(echo "$CONFIG" | jq -r '.branch_name')
    FULL_PATH=$(echo "$CONFIG" | jq -r '.full_path')
    TARGET_REPO=$(echo "$CONFIG" | jq -r '.target_repo // empty')

    # Use target repo from config or fallback
    if [ -z "$TARGET_REPO" ]; then
        TARGET_REPO=$(yq '.target_repository' /workspaces/software-factory-2.0/target-repo-config.yaml)
    fi

    # Step 2: Validate cascade pattern (first effort from main, others cascade)
    PHASE=$(echo "$CONFIG" | jq -r '.phase')
    WAVE=$(echo "$CONFIG" | jq -r '.wave')
    INDEX=$(echo "$CONFIG" | jq -r '.index // 1')

    if [[ "$PHASE" = "phase1" && "$WAVE" = "wave1" && "$INDEX" = "1" ]]; then
        # ONLY first effort of P1W1 can be from main
        if [ "$BASE_BRANCH" != "main" ]; then
            echo "🚨 R514 VIOLATION: First effort must be from main"
            exit 514
        fi
    else
        # ALL others must cascade (not from main)
        if [ "$BASE_BRANCH" = "main" ]; then
            echo "🚨 R514 VIOLATION: Non-first effort cannot branch from main!"
            echo "Must follow cascade pattern per R501"
            exit 514
        fi
    fi

    # Step 3: Clone ONLY base branch
    echo "📦 Cloning base branch: $BASE_BRANCH"
    mkdir -p "$FULL_PATH"
    git clone -b "$BASE_BRANCH" --single-branch "$TARGET_REPO" "$FULL_PATH"

    if [ $? -ne 0 ]; then
        echo "🚨 R514 VIOLATION: Clone failed"
        exit 514
    fi

    # Step 4: Create branch from cascade
    cd "$FULL_PATH"
    git checkout -b "$BRANCH_NAME"

    # Step 5: Push to remote
    git push -u origin "$BRANCH_NAME"

    if [ $? -ne 0 ]; then
        echo "🚨 R514 VIOLATION: Push failed"
        exit 514
    fi

    # Step 6: Lock config (R312)
    chmod 444 .git/config

    # Step 7: Update tracking
    jq ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".created = true |
        .pre_planned_infrastructure.efforts.\"$EFFORT_ID\".created_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" \
        /workspaces/software-factory-2.0/orchestrator-state-v3.json > tmp.json && \
    mv tmp.json /workspaces/software-factory-2.0/orchestrator-state-v3.json

    echo "✅ CASCADE INFRASTRUCTURE CREATED PROJECT_DONEFULLY"
    echo "  Effort: $EFFORT_ID"
    echo "  Branch: $BRANCH_NAME (from $BASE_BRANCH)"
    echo "  Path: $FULL_PATH"
}
```

## 🔴 SPLIT INFRASTRUCTURE PROTOCOL 🔴

### Splits Also Follow Cascade:
```bash
# SPLITS CASCADE FROM EACH OTHER
create_split_infrastructure() {
    local SPLIT_ID="$1"  # e.g., "phase1_wave1_effort-X-split-001"

    CONFIG=$(jq ".pre_planned_infrastructure.splits.\"$SPLIT_ID\"" orchestrator-state-v3.json)
    BASE_BRANCH=$(echo "$CONFIG" | jq -r '.base_branch')
    SPLIT_NUM=$(echo "$CONFIG" | jq -r '.split_number')

    # Validation: First split uses original effort's base
    # Subsequent splits cascade from previous split
    if [ "$SPLIT_NUM" = "1" ]; then
        echo "First split: using original effort base"
    else
        # Base should be previous split
        EXPECTED_PATTERN="split-$(printf "%03d" $((SPLIT_NUM - 1)))"
        if [[ ! "$BASE_BRANCH" =~ $EXPECTED_PATTERN ]]; then
            echo "🚨 R514 VIOLATION: Split not cascading from previous!"
            exit 514
        fi
    fi

    # Same cascade-aware creation
    create_effort_infrastructure "$SPLIT_ID"
}
```

## 🔴 VALIDATION AFTER CREATION 🔴

### Post-Creation Validation:
```bash
# ALWAYS VALIDATE AFTER CREATING
validate_created_infrastructure() {
    local EFFORT_ID="$1"
    local EFFORT_PATH="$2"

    cd "$EFFORT_PATH"

    # Check branch exists and is correct
    CURRENT_BRANCH=$(git branch --show-current)
    EXPECTED_BRANCH=$(jq -r ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".branch_name" \
                      /workspaces/software-factory-2.0/orchestrator-state-v3.json)

    if [ "$CURRENT_BRANCH" != "$EXPECTED_BRANCH" ]; then
        echo "🚨 R514 VIOLATION: Created wrong branch!"
        exit 514
    fi

    # Check cascade relationship
    BASE_BRANCH=$(jq -r ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".base_branch" \
                  /workspaces/software-factory-2.0/orchestrator-state-v3.json)

    if [ "$BASE_BRANCH" != "main" ]; then
        if ! git merge-base --is-ancestor "origin/$BASE_BRANCH" HEAD; then
            echo "🚨 R514 VIOLATION: Not based on cascade parent!"
            exit 514
        fi
    fi

    echo "✅ Infrastructure validated: cascade intact"
}
```

## 🚨 COMMON VIOLATIONS (AUTOMATIC FAILURE) 🚨

### ❌ VIOLATION 1: Cloning All Branches
```bash
# WRONG - Clones all branches (defaults to main usually)
git clone "$REPO" "$DIR"
cd "$DIR"
git checkout -b new-branch  # Based on main!

# RIGHT - Clone specific base branch
git clone -b "$BASE_BRANCH" --single-branch "$REPO" "$DIR"
cd "$DIR"
git checkout -b new-branch  # Based on cascade parent
```

### ❌ VIOLATION 2: Not Reading Base from State
```bash
# WRONG - Hardcoding or guessing base
if [ "$WAVE" -eq 1 ]; then
    BASE="main"
else
    BASE="some-integration-branch"  # WRONG!
fi

# RIGHT - Read from pre_planned_infrastructure
BASE=$(jq -r ".pre_planned_infrastructure.efforts.\"$ID\".base_branch" state.json)
```

### ❌ VIOLATION 3: Creating Parallel Branches
```bash
# WRONG - All efforts branch from main
for effort in effort1 effort2 effort3; do
    git checkout main
    git checkout -b "$effort"  # All from main!
done

# RIGHT - Cascade branches
git checkout main
git checkout -b effort1
git checkout effort1
git checkout -b effort2
git checkout effort2
git checkout -b effort3
```

## 🔴 ENFORCEMENT MECHANISMS 🔴

### 1. Orchestrator Validation:
```bash
# In CREATE_NEXT_INFRASTRUCTURE state
before_creating_infrastructure() {
    # Validate R514 compliance
    if ! command -v validate_r510_compliance &> /dev/null; then
        source /workspaces/software-factory-2.0/utilities/r510-enforcement.sh
    fi

    validate_r510_compliance "$EFFORT_ID"
}
```

### 2. Pre-Creation Checks:
```bash
# Prevent wrong creation
pre_creation_checks() {
    local EFFORT_ID="$1"

    # Check cascade pattern will be maintained
    BASE=$(jq -r ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".base_branch" state.json)

    # Verify base exists
    if ! git ls-remote --heads origin "$BASE" | grep -q "$BASE"; then
        echo "🚨 R514: Base branch $BASE doesn't exist!"
        exit 514
    fi
}
```

### 3. Post-Creation Audit:
```bash
# Audit all created infrastructure
audit_infrastructure() {
    jq -r '.pre_planned_infrastructure.efforts | to_entries[] |
           select(.value.created == true) | .key' orchestrator-state-v3.json |
    while read -r effort_id; do
        validate_created_infrastructure "$effort_id" "$(get_effort_path "$effort_id")"
    done
}
```

## 🔴 STATE MACHINE INTEGRATE_WAVE_EFFORTS 🔴

### CREATE_NEXT_INFRASTRUCTURE Must:
1. Read base_branch from pre_planned_infrastructure
2. Clone with --single-branch flag
3. Create branch from cascade parent
4. Validate before proceeding
5. Update final_merge_plan (R501)

### VALIDATE_INFRASTRUCTURE Must:
1. Check all branches follow cascade
2. Verify no parallel branching
3. Confirm base branches correct
4. Report violations immediately

## 🔴 GRADING IMPACT 🔴

- **Cloning wrong base**: -100% (CASCADE VIOLATION)
- **Not using --single-branch**: -75% (Risk of wrong base)
- **Creating from main (non-first)**: -100% (Parallel branching)
- **Missing base_branch in state**: -100% (No cascade tracking)
- **Not validating after creation**: -50% (Process failure)

## 🔴 WHY THIS MATTERS 🔴

### Without Proper Creation Protocol:
- **Parallel Branches**: All from main = merge hell
- **Lost Cascade**: Can't merge sequentially
- **Conflicts Everywhere**: Branches don't build on each other
- **Integration Nightmare**: Big-bang merges fail

### With R514 Protocol:
- **Perfect Cascade**: Each builds on previous
- **Clean Merges**: Sequential integration works
- **No Conflicts**: Changes accumulate properly
- **Smooth Flow**: Progressive development achieved

## 🔴 THE FINAL TRUTH 🔴

**INFRASTRUCTURE CREATION IS THE FOUNDATION OF CASCADE!**

- ALWAYS read base from pre_planned_infrastructure
- ALWAYS clone with --single-branch
- ALWAYS create from cascade parent
- ALWAYS validate the cascade
- NEVER create parallel branches from main
- NEVER guess or hardcode bases

**R514 ensures every branch is born in the right place in the cascade!**
## State Manager Coordination (SF 3.0)

State Manager coordinates infrastructure creation with state file updates:
- **SETUP_WAVE_INFRASTRUCTURE state**: Create directories + update state file
- **Atomic operation**: Filesystem changes + state file changes committed together
- **Rollback on failure**: Cleanup directories if state update fails
- **Validation**: Verify infrastructure matches state file after creation

Infrastructure creation is never complete until state file reflects it (and is committed).

See: R505 (infrastructure sync), R281 (initial creation), `agent-states/software-factory/orchestrator/SETUP_WAVE_INFRASTRUCTURE/rules.md`
