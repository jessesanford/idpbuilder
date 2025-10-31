# SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE State Rules

## State Context

**Current Phase**: Read from `orchestrator-state-v3.json`
**Current Wave**: Read from `orchestrator-state-v3.json`
**Entry From**: INTEGRATE_WAVE_EFFORTS, CASCADE_INTEGRATE_WAVE_EFFORTS
**Exit To**: PERFORM_INTEGRATE_WAVE_EFFORTS

## 🔴🔴🔴 SUPREME LAW: R504 PRE-PLANNED INFRASTRUCTURE ENFORCEMENT 🔴🔴🔴

**THIS STATE MUST ONLY USE PRE-PLANNED INFRASTRUCTURE - NO RUNTIME DECISIONS!**

- ❌ **FORBIDDEN**: Making ANY naming decisions in this state
- ❌ **FORBIDDEN**: Calculating integration branch names at runtime
- ❌ **FORBIDDEN**: Determining paths dynamically
- ❌ **FORBIDDEN**: Creating infrastructure not in pre_planned_infrastructure
- ✅ **REQUIRED**: Use ONLY pre_planned_infrastructure.integrations from orchestrator-state-v3.json
- ✅ **REQUIRED**: Verify integration infrastructure was pre-planned
- ✅ **REQUIRED**: Fail immediately if integration not pre-planned

## 🚨🚨🚨 PRIMARY OBJECTIVE 🚨🚨🚨

**SETUP integration branch infrastructure using ONLY pre_planned_infrastructure per R504.**

This state creates the integration branch and directory structure for merging effort branches.

## Required Inputs

### 1. Verify Pre-Planned Integration Infrastructure (R504 ENFORCEMENT)
```bash
echo "🔍 Verifying pre-planned integration infrastructure per R504..."

PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)
INTEGRATE_WAVE_EFFORTS_KEY="phase${PHASE}_wave${WAVE}"

# Check that integration infrastructure exists in pre_planned_infrastructure.integrations.wave_integrations
INTEGRATE_WAVE_EFFORTS_CONFIG=$(jq -r ".pre_planned_infrastructure.integrations.wave_integrations.\"${INTEGRATE_WAVE_EFFORTS_KEY}\" // empty" orchestrator-state-v3.json)

if [ -z "$INTEGRATE_WAVE_EFFORTS_CONFIG" ] || [ "$INTEGRATE_WAVE_EFFORTS_CONFIG" == "null" ]; then
    echo "❌ FATAL: No pre-planned integration infrastructure for ${INTEGRATE_WAVE_EFFORTS_KEY}!"
    echo "  R504 VIOLATION: Integration infrastructure must be pre-planned"
    echo "  This state CANNOT make infrastructure decisions!"
    echo "  Checked: .pre_planned_infrastructure.integrations.wave_integrations.${INTEGRATE_WAVE_EFFORTS_KEY}"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

# Extract pre-planned values
INTEGRATE_WAVE_EFFORTS_BRANCH=$(echo "$INTEGRATE_WAVE_EFFORTS_CONFIG" | jq -r '.branch_name')
INTEGRATE_WAVE_EFFORTS_REMOTE=$(echo "$INTEGRATE_WAVE_EFFORTS_CONFIG" | jq -r '.remote_branch')
INTEGRATE_WAVE_EFFORTS_BASE=$(echo "$INTEGRATE_WAVE_EFFORTS_CONFIG" | jq -r '.base_branch')
TARGET_REPO_URL=$(echo "$INTEGRATE_WAVE_EFFORTS_CONFIG" | jq -r '.target_repo_url')
INTEGRATE_WAVE_EFFORTS_DIR=$(echo "$INTEGRATE_WAVE_EFFORTS_CONFIG" | jq -r '.directory // empty')
COMPONENT_EFFORTS=$(echo "$INTEGRATE_WAVE_EFFORTS_CONFIG" | jq -r '.component_efforts[]')

if [ -z "$INTEGRATE_WAVE_EFFORTS_BRANCH" ] || [ "$INTEGRATE_WAVE_EFFORTS_BRANCH" == "null" ]; then
    echo "❌ FATAL: No branch_name in pre-planned integration infrastructure!"
    echo "  R504 VIOLATION: All branch names must be pre-planned"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

if [ -z "$TARGET_REPO_URL" ] || [ "$TARGET_REPO_URL" == "null" ]; then
    echo "❌ FATAL: No target_repo_url in pre-planned integration infrastructure!"
    echo "  R504 VIOLATION: Target repository must be pre-planned"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

echo "✅ Pre-planned integration infrastructure found:"
echo "  Branch: $INTEGRATE_WAVE_EFFORTS_BRANCH"
echo "  Remote: $INTEGRATE_WAVE_EFFORTS_REMOTE"
echo "  Base: $INTEGRATE_WAVE_EFFORTS_BASE"
echo "  Target Repo: $TARGET_REPO_URL"
echo "  Directory: $INTEGRATE_WAVE_EFFORTS_DIR"
echo "  Component efforts: $COMPONENT_EFFORTS"
```

### 2. Defensive Pre-Creation Validation (R327 ENFORCEMENT)
```bash
echo "🔍 R327 DEFENSE: Verifying infrastructure is clean (no remnants)..."

VALIDATION_FAILED=false
FAILURE_REASONS=()

# Check 1: Local branch must NOT exist
if git show-ref --verify --quiet "refs/heads/$INTEGRATE_WAVE_EFFORTS_BRANCH"; then
    echo "❌ FATAL: Local branch already exists: $INTEGRATE_WAVE_EFFORTS_BRANCH"
    FAILURE_REASONS+=("Local branch '$INTEGRATE_WAVE_EFFORTS_BRANCH' exists (should be deleted by CASCADE_REINTEGRATION)")
    VALIDATION_FAILED=true
fi

# Check 2: Working directory must NOT exist (or must be empty)
if [ -n "$INTEGRATE_WAVE_EFFORTS_DIR" ] && [ "$INTEGRATE_WAVE_EFFORTS_DIR" != "null" ]; then
    # Convert to absolute path
    local abs_dir
    if [[ "$INTEGRATE_WAVE_EFFORTS_DIR" = /* ]]; then
        abs_dir="$INTEGRATE_WAVE_EFFORTS_DIR"
    else
        abs_dir="$CLAUDE_PROJECT_DIR/$INTEGRATE_WAVE_EFFORTS_DIR"
    fi

    if [ -d "$abs_dir" ]; then
        # Check if directory is empty
        if [ "$(ls -A "$abs_dir" 2>/dev/null)" ]; then
            echo "❌ FATAL: Working directory exists and is not empty: $abs_dir"
            FAILURE_REASONS+=("Directory '$abs_dir' exists with content (should be removed by CASCADE_REINTEGRATION)")
            VALIDATION_FAILED=true
        else
            echo "⚠️ Working directory exists but is empty: $abs_dir (will use it)"
        fi
    fi
fi

# Check 3: State file must show created=false
CREATED_STATUS=$(echo "$INTEGRATE_WAVE_EFFORTS_CONFIG" | jq -r '.created // false')
if [ "$CREATED_STATUS" = "true" ]; then
    echo "❌ FATAL: State file shows integration already created"
    FAILURE_REASONS+=("State shows created=true (should be reset to false by CASCADE_REINTEGRATION)")
    VALIDATION_FAILED=true
fi

# If any validation failed, FAIL HARD with clear guidance
if [ "$VALIDATION_FAILED" = true ]; then
    echo ""
    echo "🔴🔴🔴 R327 DEFENSIVE CHECK FAILED 🔴🔴🔴"
    echo ""
    echo "SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE expects CLEAN infrastructure."
    echo "Old infrastructure remnants detected - this indicates a BUG!"
    echo ""
    echo "Failure reasons:"
    for reason in "${FAILURE_REASONS[@]}"; do
        echo "  ❌ $reason"
    done
    echo ""
    echo "🔴 ROOT CAUSE: CASCADE_REINTEGRATION did not clean up properly!"
    echo ""
    echo "REQUIRED ACTIONS:"
    echo "  1. CASCADE_REINTEGRATION MUST delete:"
    echo "     - Remote branch (git push target --delete)"
    echo "     - Local branch (git branch -D)"
    echo "     - Working directory (rm -rf)"
    echo "     - State tracking (created=false)"
    echo ""
    echo "  2. Fix CASCADE_REINTEGRATION cleanup_integration_infrastructure() function"
    echo "  3. Re-run CASCADE_REINTEGRATION to properly clean up"
    echo ""
    echo "CANNOT PROCEED - Infrastructure collision will cause failures!"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

echo "✅ R327 DEFENSIVE CHECK PASSED: Infrastructure is clean"
echo "   - No local branch exists"
echo "   - Working directory is clean"
echo "   - State shows created=false"
echo "   - Ready for fresh creation"
```

### 3. Create Integration Directory (FROM PRE-PLANNED DATA ONLY)
```bash
echo "📁 Creating integration directory from pre-planned data..."

if [ -n "$INTEGRATE_WAVE_EFFORTS_DIR" ] && [ "$INTEGRATE_WAVE_EFFORTS_DIR" != "null" ]; then
    if [ ! -d "$INTEGRATE_WAVE_EFFORTS_DIR" ]; then
        mkdir -p "$INTEGRATE_WAVE_EFFORTS_DIR"
        echo "✅ Created integration directory: $INTEGRATE_WAVE_EFFORTS_DIR"
    else
        echo "⚠️ Integration directory already exists: $INTEGRATE_WAVE_EFFORTS_DIR"
    fi

    cd "$INTEGRATE_WAVE_EFFORTS_DIR"
else
    # Use project root if no specific directory
    cd "$CLAUDE_PROJECT_DIR"
fi
```

### 3. Create Integration Branch (FROM PRE-PLANNED DATA ONLY)
```bash
echo "🌿 Creating integration branch from pre-planned data..."

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

# Create integration branch from main (or from pre-planned base)
if ! git show-ref --verify --quiet "refs/heads/$INTEGRATE_WAVE_EFFORTS_BRANCH"; then
    git checkout -b "$INTEGRATE_WAVE_EFFORTS_BRANCH" "${TARGET_REMOTE}/main"
    git push -u "$TARGET_REMOTE" "$INTEGRATE_WAVE_EFFORTS_BRANCH"
    echo "✅ Created and pushed integration branch: $INTEGRATE_WAVE_EFFORTS_BRANCH"
else
    echo "⚠️ Integration branch already exists: $INTEGRATE_WAVE_EFFORTS_BRANCH"
    git checkout "$INTEGRATE_WAVE_EFFORTS_BRANCH"
fi

# Mark as created in pre_planned_infrastructure
jq --arg key "$INTEGRATE_WAVE_EFFORTS_KEY" \
   '.pre_planned_infrastructure.integrations.wave_integrations[$key].created = true' \
   "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" > tmp.json && \
   mv tmp.json "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
```

## Validation Requirements

### Post-Creation Validation
```bash
echo "🔍 Validating integration infrastructure..."

VALIDATION_FAILED=false

# Check branch exists
if ! git show-ref --verify --quiet "refs/heads/$INTEGRATE_WAVE_EFFORTS_BRANCH"; then
    echo "❌ Integration branch not created: $INTEGRATE_WAVE_EFFORTS_BRANCH"
    VALIDATION_FAILED=true
fi

# Check remote tracking
if ! git branch -vv | grep "$INTEGRATE_WAVE_EFFORTS_BRANCH" | grep -q "\[${TARGET_REMOTE}/"; then
    echo "❌ Integration branch not tracking remote"
    VALIDATION_FAILED=true
fi

# Check directory exists (if specified)
if [ -n "$INTEGRATE_WAVE_EFFORTS_DIR" ] && [ "$INTEGRATE_WAVE_EFFORTS_DIR" != "null" ]; then
    if [ ! -d "$INTEGRATE_WAVE_EFFORTS_DIR" ]; then
        echo "❌ Integration directory not created: $INTEGRATE_WAVE_EFFORTS_DIR"
        VALIDATION_FAILED=true
    fi
fi

if [ "$VALIDATION_FAILED" = true ]; then
    echo "❌ FATAL: Integration infrastructure validation failed"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

echo "✅ Integration infrastructure validation passed"
```

## Exit Criteria

Before transitioning to PERFORM_INTEGRATE_WAVE_EFFORTS:
- ✅ Integration branch created from pre-planned data
- ✅ Branch tracking correct remote
- ✅ Integration directory created (if specified)
- ✅ pre_planned_infrastructure marked as created
- ✅ All validation checks passed

## Common Issues

### Issue: No Pre-Planned Infrastructure
**Detection**: pre_planned_infrastructure missing integration config
**Resolution**: FAIL - must transition to ERROR_RECOVERY, infrastructure must be pre-planned

### Issue: Branch Name Conflict
**Detection**: Different branch name already exists
**Resolution**: FAIL - pre-planned name must be used, no runtime decisions allowed

## Automation Flag - CASCADE INTEGRATE_WAVE_EFFORTS

### 🔴🔴🔴 SETUP DURING CASCADE - ALWAYS USE TRUE 🔴🔴🔴

**This state is often entered FROM CASCADE_REINTEGRATION as part of R327 cascade recreation.**

### Decision Tree for SETUP Exit

```
Why did you enter SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE?
├─ FROM CASCADE_REINTEGRATION (cascade mode active)
│  └─ Infrastructure created successfully?
│     ├─ YES → Returning to CASCADE for next step
│     │         └─ Use CONTINUE-SOFTWARE-FACTORY=TRUE
│     └─ NO (unrecoverable) → Cannot create infrastructure
│               └─ Use CONTINUE-SOFTWARE-FACTORY=FALSE
│
└─ FROM other state (normal integration)
   └─ Infrastructure created successfully?
      ├─ YES → Proceeding to PERFORM_INTEGRATE_WAVE_EFFORTS
      │         └─ Use CONTINUE-SOFTWARE-FACTORY=TRUE
      └─ NO (unrecoverable) → Cannot create infrastructure
                └─ Use CONTINUE-SOFTWARE-FACTORY=FALSE
```

### Standard Exit Pattern

```bash
# After successfully creating integration infrastructure:
echo "✅ Integration infrastructure setup from pre-planned data"

# Check if in cascade mode
if jq -e '.cascade_mode == true' orchestrator-state-v3.json > /dev/null 2>&1; then
    echo "🔄 CASCADE MODE: Returning to CASCADE_REINTEGRATION for next step"

    # Update state to return to cascade
    update_state "CASCADE_REINTEGRATION" "Completed infrastructure setup, continuing cascade"

    # R322 checkpoint
    save_todos "SETUP_COMPLETE_CASCADE_CONTINUE"
    git add todos/*.todo orchestrator-state-v3.json
    git commit -m "state: SETUP → CASCADE (infrastructure created, cascade continues)"
    git push

    # 🔴 CASCADE CONTINUATION - ALWAYS TRUE! 🔴
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Cascade continues automatically!
    exit 0
else
    echo "➡️ NORMAL MODE: Proceeding to PERFORM_INTEGRATE_WAVE_EFFORTS"

    # Update state to integration
    update_state "PERFORM_INTEGRATE_WAVE_EFFORTS" "Infrastructure ready, starting integration"

    # R322 checkpoint
    save_todos "SETUP_COMPLETE_INTEGRATE_WAVE_EFFORTS_START"
    git add todos/*.todo orchestrator-state-v3.json
    git commit -m "state: SETUP → PERFORM_INTEGRATE_WAVE_EFFORTS"
    git push

    # 🔴 NORMAL WORKFLOW - ALWAYS TRUE! 🔴
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Integration is automated!
    exit 0
fi
```

### Why TRUE for SETUP Exits?

**SETUP → CASCADE (cascade mode):**
- System knows current state (CASCADE_REINTEGRATION from state file)
- System knows cascade chain (remaining integrations to recreate)
- System knows next action (delete next integration or complete cascade)
- NO HUMAN INTERVENTION NEEDED!

**SETUP → PERFORM_INTEGRATE_WAVE_EFFORTS (normal mode):**
- System knows current state (PERFORM_INTEGRATE_WAVE_EFFORTS from state file)
- System knows what to do (merge efforts per integration protocol)
- System knows next steps (build, test, review)
- NO HUMAN INTERVENTION NEEDED!

### ❌ ONLY Use FALSE When (0.01% of cases)

```bash
# Truly unrecoverable infrastructure failure
if [ "$VALIDATION_FAILED" = true ] && [ "$CASCADE_CLEANUP_FAILED" = true ]; then
    echo "❌ CRITICAL: Cannot create infrastructure - cleanup failed"
    echo "   Old infrastructure remnants detected"
    echo "   CASCADE_REINTEGRATION cleanup did not work"
    echo "   Manual debugging required"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi
```

### Summary

- ✅ **TRUE when exiting to CASCADE**: Cascade continues automatically
- ✅ **TRUE when exiting to INTEGRATE_WAVE_EFFORTS**: Integration is automated workflow
- ❌ **FALSE**: ONLY for unrecoverable infrastructure failures
- 🛑 **R322 stops**: INDEPENDENT from flag value!

**SETUP infrastructure creation is ALWAYS normal workflow - use TRUE!**

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
