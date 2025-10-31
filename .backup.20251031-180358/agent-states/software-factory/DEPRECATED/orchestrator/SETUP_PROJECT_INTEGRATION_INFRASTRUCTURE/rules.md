# SETUP_PROJECT_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE State Rules

## State Context

**Entry From**: PROJECT_INTEGRATE_WAVE_EFFORTS
**Exit To**: PERFORM_PROJECT_INTEGRATE_WAVE_EFFORTS

## 🔴🔴🔴 SUPREME LAW: R504 PRE-PLANNED INFRASTRUCTURE ENFORCEMENT 🔴🔴🔴

**THIS STATE MUST ONLY USE PRE-PLANNED INFRASTRUCTURE - NO RUNTIME DECISIONS!**

- ❌ **FORBIDDEN**: Making ANY naming decisions in this state
- ❌ **FORBIDDEN**: Calculating project integration branch names at runtime
- ❌ **FORBIDDEN**: Determining paths dynamically
- ❌ **FORBIDDEN**: Creating infrastructure not in pre_planned_infrastructure
- ✅ **REQUIRED**: Use ONLY pre_planned_infrastructure.project_integration from orchestrator-state-v3.json
- ✅ **REQUIRED**: Verify project integration infrastructure was pre-planned
- ✅ **REQUIRED**: Fail immediately if project integration not pre-planned

## 🚨🚨🚨 PRIMARY OBJECTIVE 🚨🚨🚨

**SETUP project integration branch infrastructure using ONLY pre_planned_infrastructure per R504.**

This state creates the project-level integration branch for merging all phase integration branches.

## Required Inputs

### 1. Verify Pre-Planned Project Integration Infrastructure (R504 ENFORCEMENT)
```bash
echo "🔍 Verifying pre-planned project integration infrastructure per R504..."

# Check that project integration infrastructure exists in pre_planned_infrastructure
PROJECT_CONFIG=$(jq -r ".pre_planned_infrastructure.project_integration // empty" orchestrator-state-v3.json)

if [ -z "$PROJECT_CONFIG" ] || [ "$PROJECT_CONFIG" == "null" ]; then
    echo "❌ FATAL: No pre-planned project integration infrastructure!"
    echo "  R504 VIOLATION: Project integration infrastructure must be pre-planned"
    echo "  This state CANNOT make infrastructure decisions!"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

# Extract pre-planned values
PROJECT_BRANCH=$(echo "$PROJECT_CONFIG" | jq -r '.branch_name')
PROJECT_DIR=$(echo "$PROJECT_CONFIG" | jq -r '.directory // empty')
COMPONENT_PHASES=$(echo "$PROJECT_CONFIG" | jq -r '.component_phases[]')

if [ -z "$PROJECT_BRANCH" ] || [ "$PROJECT_BRANCH" == "null" ]; then
    echo "❌ FATAL: No branch_name in pre-planned project integration infrastructure!"
    echo "  R504 VIOLATION: All branch names must be pre-planned"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

echo "✅ Pre-planned project integration infrastructure found:"
echo "  Branch: $PROJECT_BRANCH"
echo "  Directory: $PROJECT_DIR"
echo "  Component phases: $COMPONENT_PHASES"
```

### 2. Create Project Integration Directory (FROM PRE-PLANNED DATA ONLY)
```bash
echo "📁 Creating project integration directory from pre-planned data..."

if [ -n "$PROJECT_DIR" ] && [ "$PROJECT_DIR" != "null" ]; then
    if [ ! -d "$PROJECT_DIR" ]; then
        mkdir -p "$PROJECT_DIR"
        echo "✅ Created project integration directory: $PROJECT_DIR"
    else
        echo "⚠️ Project integration directory already exists: $PROJECT_DIR"
    fi

    cd "$PROJECT_DIR"
else
    # Use project root if no specific directory
    cd "$CLAUDE_PROJECT_DIR"
fi
```

### 3. Create Project Integration Branch (FROM PRE-PLANNED DATA ONLY)
```bash
echo "🌿 Creating project integration branch from pre-planned data..."

# Get target repository from config
TARGET_URL=$(yq -r '.url' "$CLAUDE_PROJECT_DIR/target-repo-config.yaml")
TARGET_REMOTE="target"

# Ensure we have the target remote
if ! git remote | grep -q "^${TARGET_REMOTE}$"; then
    git remote add "$TARGET_REMOTE" "$TARGET_URL"
    echo "✅ Added remote: $TARGET_REMOTE"
fi

# Fetch latest from target
git fetch "$TARGET_REMOTE"

# Create project integration branch from main (or from pre-planned base)
if ! git show-ref --verify --quiet "refs/heads/$PROJECT_BRANCH"; then
    git checkout -b "$PROJECT_BRANCH" "${TARGET_REMOTE}/main"
    git push -u "$TARGET_REMOTE" "$PROJECT_BRANCH"
    echo "✅ Created and pushed project integration branch: $PROJECT_BRANCH"
else
    echo "⚠️ Project integration branch already exists: $PROJECT_BRANCH"
    git checkout "$PROJECT_BRANCH"
fi

# Mark as created in pre_planned_infrastructure
jq '.pre_planned_infrastructure.project_integration.created = true' \
   "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" > tmp.json && \
   mv tmp.json "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
```

## Validation Requirements

### Post-Creation Validation
```bash
echo "🔍 Validating project integration infrastructure..."

VALIDATION_FAILED=false

# Check branch exists
if ! git show-ref --verify --quiet "refs/heads/$PROJECT_BRANCH"; then
    echo "❌ Project integration branch not created: $PROJECT_BRANCH"
    VALIDATION_FAILED=true
fi

# Check remote tracking
if ! git branch -vv | grep "$PROJECT_BRANCH" | grep -q "\[${TARGET_REMOTE}/"; then
    echo "❌ Project integration branch not tracking remote"
    VALIDATION_FAILED=true
fi

# Check directory exists (if specified)
if [ -n "$PROJECT_DIR" ] && [ "$PROJECT_DIR" != "null" ]; then
    if [ ! -d "$PROJECT_DIR" ]; then
        echo "❌ Project integration directory not created: $PROJECT_DIR"
        VALIDATION_FAILED=true
    fi
fi

if [ "$VALIDATION_FAILED" = true ]; then
    echo "❌ FATAL: Project integration infrastructure validation failed"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

echo "✅ Project integration infrastructure validation passed"
```

## Exit Criteria

Before transitioning to PERFORM_PROJECT_INTEGRATE_WAVE_EFFORTS:
- ✅ Project integration branch created from pre-planned data
- ✅ Branch tracking correct remote
- ✅ Project integration directory created (if specified)
- ✅ pre_planned_infrastructure marked as created
- ✅ All validation checks passed

## Common Issues

### Issue: No Pre-Planned Infrastructure
**Detection**: pre_planned_infrastructure missing project integration config
**Resolution**: FAIL - must transition to ERROR_RECOVERY, infrastructure must be pre-planned

### Issue: Branch Name Conflict
**Detection**: Different branch name already exists
**Resolution**: FAIL - pre-planned name must be used, no runtime decisions allowed

## Automation Flag

```bash
# After successfully creating project integration infrastructure:
echo "✅ Project integration infrastructure setup from pre-planned data"
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
