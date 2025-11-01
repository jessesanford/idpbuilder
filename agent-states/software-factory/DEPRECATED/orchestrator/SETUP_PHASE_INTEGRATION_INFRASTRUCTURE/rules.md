# SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## State Context

**Current Phase**: Read from `orchestrator-state-v3.json`
**Entry From**: INTEGRATE_PHASE_WAVES
**Exit To**: PERFORM_INTEGRATE_PHASE_WAVES

## 🔴🔴🔴 SUPREME LAW: R504 PRE-PLANNED INFRASTRUCTURE ENFORCEMENT 🔴🔴🔴

**THIS STATE MUST ONLY USE PRE-PLANNED INFRASTRUCTURE - NO RUNTIME DECISIONS!**

- ❌ **FORBIDDEN**: Making ANY naming decisions in this state
- ❌ **FORBIDDEN**: Calculating phase integration branch names at runtime
- ❌ **FORBIDDEN**: Determining paths dynamically
- ❌ **FORBIDDEN**: Creating infrastructure not in pre_planned_infrastructure
- ✅ **REQUIRED**: Use ONLY pre_planned_infrastructure.phase_integrations from orchestrator-state-v3.json
- ✅ **REQUIRED**: Verify phase integration infrastructure was pre-planned
- ✅ **REQUIRED**: Fail immediately if phase integration not pre-planned

## 🚨🚨🚨 PRIMARY OBJECTIVE 🚨🚨🚨

**SETUP phase integration branch infrastructure using ONLY pre_planned_infrastructure per R504.**

This state creates the phase-level integration branch for merging all wave integration branches.

## Required Inputs

### 1. Verify Pre-Planned Phase Integration Infrastructure (R504 ENFORCEMENT)
```bash
echo "🔍 Verifying pre-planned phase integration infrastructure per R504..."

PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
PHASE_KEY="phase${PHASE}"

# Check that phase integration infrastructure exists in pre_planned_infrastructure.integrations.phase_integrations
PHASE_CONFIG=$(jq -r ".pre_planned_infrastructure.integrations.phase_integrations.\"${PHASE_KEY}\" // empty" orchestrator-state-v3.json)

if [ -z "$PHASE_CONFIG" ] || [ "$PHASE_CONFIG" == "null" ]; then
    echo "❌ FATAL: No pre-planned phase integration infrastructure for ${PHASE_KEY}!"
    echo "  R504 VIOLATION: Phase integration infrastructure must be pre-planned"
    echo "  This state CANNOT make infrastructure decisions!"
    echo "  Checked: .pre_planned_infrastructure.integrations.phase_integrations.${PHASE_KEY}"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

# Extract pre-planned values
PHASE_BRANCH=$(echo "$PHASE_CONFIG" | jq -r '.branch_name')
PHASE_REMOTE=$(echo "$PHASE_CONFIG" | jq -r '.remote_branch')
PHASE_BASE=$(echo "$PHASE_CONFIG" | jq -r '.base_branch')
TARGET_REPO_URL=$(echo "$PHASE_CONFIG" | jq -r '.target_repo_url')
PHASE_DIR=$(echo "$PHASE_CONFIG" | jq -r '.directory // empty')
COMPONENT_WAVES=$(echo "$PHASE_CONFIG" | jq -r '.component_waves[]')

if [ -z "$PHASE_BRANCH" ] || [ "$PHASE_BRANCH" == "null" ]; then
    echo "❌ FATAL: No branch_name in pre-planned phase integration infrastructure!"
    echo "  R504 VIOLATION: All branch names must be pre-planned"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

if [ -z "$TARGET_REPO_URL" ] || [ "$TARGET_REPO_URL" == "null" ]; then
    echo "❌ FATAL: No target_repo_url in pre-planned phase integration infrastructure!"
    echo "  R504 VIOLATION: Target repository must be pre-planned"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

echo "✅ Pre-planned phase integration infrastructure found:"
echo "  Branch: $PHASE_BRANCH"
echo "  Remote: $PHASE_REMOTE"
echo "  Base: $PHASE_BASE"
echo "  Target Repo: $TARGET_REPO_URL"
echo "  Directory: $PHASE_DIR"
echo "  Component waves: $COMPONENT_WAVES"
```

### 2. Create Phase Integration Directory (FROM PRE-PLANNED DATA ONLY)
```bash
echo "📁 Creating phase integration directory from pre-planned data..."

if [ -n "$PHASE_DIR" ] && [ "$PHASE_DIR" != "null" ]; then
    if [ ! -d "$PHASE_DIR" ]; then
        mkdir -p "$PHASE_DIR"
        echo "✅ Created phase integration directory: $PHASE_DIR"
    else
        echo "⚠️ Phase integration directory already exists: $PHASE_DIR"
    fi

    cd "$PHASE_DIR"
else
    # Use project root if no specific directory
    cd "$CLAUDE_PROJECT_DIR"
fi
```

### 3. Create Phase Integration Branch (FROM PRE-PLANNED DATA ONLY)
```bash
echo "🌿 Creating phase integration branch from pre-planned data..."

# Get target repository from config and verify it matches pre-planned
TARGET_URL_CONFIG=$(yq -r '.url' "$CLAUDE_PROJECT_DIR/target-repo-config.yaml")
TARGET_REMOTE="target"

# CRITICAL: Verify pre-planned URL matches config
if [ "$TARGET_REPO_URL" != "$TARGET_URL_CONFIG" ]; then
    echo "❌ FATAL: Pre-planned target_repo_url doesn't match target-repo-config.yaml!"
    echo "  Pre-planned: $TARGET_REPO_URL"
    echo "  Config:      $TARGET_URL_CONFIG"
    echo "  This indicates corrupted pre-planning!"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

# Ensure we have the target remote pointing to the correct URL
if ! git remote | grep -q "^${TARGET_REMOTE}$"; then
    git remote add "$TARGET_REMOTE" "$TARGET_REPO_URL"
    echo "✅ Added remote: $TARGET_REMOTE -> $TARGET_REPO_URL"
else
    # Verify existing remote points to correct URL
    CURRENT_URL=$(git remote get-url "$TARGET_REMOTE")
    if [ "$CURRENT_URL" != "$TARGET_REPO_URL" ]; then
        echo "⚠️ Updating remote URL from $CURRENT_URL to $TARGET_REPO_URL"
        git remote set-url "$TARGET_REMOTE" "$TARGET_REPO_URL"
    fi
fi

# Fetch latest from target
git fetch "$TARGET_REMOTE"

# Create phase integration branch from main (or from pre-planned base)
if ! git show-ref --verify --quiet "refs/heads/$PHASE_BRANCH"; then
    git checkout -b "$PHASE_BRANCH" "${TARGET_REMOTE}/main"
    git push -u "$TARGET_REMOTE" "$PHASE_BRANCH"
    echo "✅ Created and pushed phase integration branch: $PHASE_BRANCH"
else
    echo "⚠️ Phase integration branch already exists: $PHASE_BRANCH"
    git checkout "$PHASE_BRANCH"
fi

# Mark as created in pre_planned_infrastructure
jq --arg key "$PHASE_KEY" \
   '.pre_planned_infrastructure.phase_integrations[$key].created = true' \
   "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" > tmp.json && \
   mv tmp.json "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
```

## Validation Requirements

### Post-Creation Validation
```bash
echo "🔍 Validating phase integration infrastructure..."

VALIDATION_FAILED=false

# Check branch exists
if ! git show-ref --verify --quiet "refs/heads/$PHASE_BRANCH"; then
    echo "❌ Phase integration branch not created: $PHASE_BRANCH"
    VALIDATION_FAILED=true
fi

# Check remote tracking
if ! git branch -vv | grep "$PHASE_BRANCH" | grep -q "\[${TARGET_REMOTE}/"; then
    echo "❌ Phase integration branch not tracking remote"
    VALIDATION_FAILED=true
fi

# Check directory exists (if specified)
if [ -n "$PHASE_DIR" ] && [ "$PHASE_DIR" != "null" ]; then
    if [ ! -d "$PHASE_DIR" ]; then
        echo "❌ Phase integration directory not created: $PHASE_DIR"
        VALIDATION_FAILED=true
    fi
fi

if [ "$VALIDATION_FAILED" = true ]; then
    echo "❌ FATAL: Phase integration infrastructure validation failed"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

echo "✅ Phase integration infrastructure validation passed"
```

## Exit Criteria

Before transitioning to PERFORM_INTEGRATE_PHASE_WAVES:
- ✅ Phase integration branch created from pre-planned data
- ✅ Branch tracking correct remote
- ✅ Phase integration directory created (if specified)
- ✅ pre_planned_infrastructure marked as created
- ✅ All validation checks passed

## Common Issues

### Issue: No Pre-Planned Infrastructure
**Detection**: pre_planned_infrastructure missing phase integration config
**Resolution**: FAIL - must transition to ERROR_RECOVERY, infrastructure must be pre-planned

### Issue: Branch Name Conflict
**Detection**: Different branch name already exists
**Resolution**: FAIL - pre-planned name must be used, no runtime decisions allowed

## Automation Flag

```bash
# After successfully creating phase integration infrastructure:
echo "✅ Phase integration infrastructure setup from pre-planned data"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```

## Related Rules
- R504: Pre-Infrastructure Planning (SUPREME LAW)
- R507: Infrastructure Validation
- R508: Repository Validation
- R327: Cascade Branching Conventions

---

**REMEMBER**: This state MUST NOT make ANY naming or pathing decisions. ALL infrastructure details MUST come from pre_planned_infrastructure!
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
