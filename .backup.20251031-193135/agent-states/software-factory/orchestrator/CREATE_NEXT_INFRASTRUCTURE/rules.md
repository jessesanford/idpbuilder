# Orchestrator - CREATE_NEXT_INFRASTRUCTURE State Rules

## ⚠️⚠️⚠️ CRITICAL: Upstream Tracking Required ⚠️⚠️⚠️

**EVERY `git push` command in this state MUST use the `-u` flag:**

```bash
git push -u origin "$BRANCH_NAME"
```

**Why `-u` is mandatory:**
- Sets up upstream tracking for the local branch
- Enables `git pull` and `git push` without arguments
- Required by VALIDATE_INFRASTRUCTURE state validation
- **Missing `-u` flag = validation failure → recovery loop → wasted time and money**

**What happens without `-u`:**
1. Branch is pushed to remote (data transfer succeeds) ✅
2. Local branch does NOT track remote branch ❌
3. `git rev-parse --abbrev-ref @{upstream}` fails ❌
4. VALIDATE_INFRASTRUCTURE detects missing tracking ❌
5. Recovery loop triggered → return to CREATE_NEXT_INFRASTRUCTURE ❌
6. Costs ~$0.40-0.50 per retry iteration ❌

**Bottom line:** Always use `git push -u origin "$BRANCH_NAME"` - no exceptions!

---

## State Manager Bookend Pattern (MANDATORY)

**BEFORE this state**:
- State Manager validated transition to this state via STARTUP_CONSULTATION
- You are here because State Manager directed you here

**DURING this state**:
- Perform state-specific work
- NEVER update state files directly
- NEVER call update_state function
- Prepare results for State Manager

**AFTER this state**:
- Spawn State Manager SHUTDOWN_CONSULTATION
- Provide results and proposed next state
- State Manager decides actual next state
- Transition to State Manager's required_next_state

**PROHIBITED**:
- ❌ Calling update_state directly
- ❌ Updating orchestrator-state-v3.json directly
- ❌ Setting validated_by: "orchestrator"
- ❌ Bypassing State Manager consultation

---

## 👉 R322 CLARIFICATION FOR CREATE_NEXT_INFRASTRUCTURE STATE 👉

**IMPORTANT**: Creating infrastructure and transitioning to SPAWN_SW_ENGINEERS is NORMAL automation flow!

### NORMAL TRANSITIONS (CONTINUE-SOFTWARE-FACTORY=TRUE):
- ✅ CREATE_NEXT_INFRASTRUCTURE → SPAWN_SW_ENGINEERS (infrastructure ready - normal)
- ✅ CREATE_NEXT_INFRASTRUCTURE → WAVE_COMPLETE (no more infrastructure - normal)

**R322 NOTE**: The stop is ONLY required AFTER spawning agents in the next state, not here!

### EXCEPTIONAL CASES (CONTINUE-SOFTWARE-FACTORY=FALSE):
- ❌ Unable to create branch/directory
- ❌ Git operations failing
- ❌ Pre-planned infrastructure missing/corrupted
- ❌ Wrong working directory

**DEFAULT**: Use CONTINUE-SOFTWARE-FACTORY=TRUE after successfully creating infrastructure!

---

## 🔴🔴🔴 R504 - PRE-INFRASTRUCTURE PLANNING PROTOCOL 🔴🔴🔴

**SUPREME LAW - ALL INFRASTRUCTURE PRE-CALCULATED, NO RUNTIME DECISIONS**

### FUNDAMENTAL CHANGE:
- ❌ NO MORE runtime naming/pathing decisions
- ❌ NO MORE calculating branches at creation time
- ✅ ALL infrastructure pre-planned in orchestrator-state-v3.json
- ✅ Creation is 100% MECHANICAL EXECUTION

### Pre-Infrastructure Validation Required:
```bash
# MUST verify pre_planned_infrastructure exists and is validated
if ! yq '.pre_planned_infrastructure.validated == true' orchestrator-state-v3.json; then
    echo "🚨 CRITICAL: No validated pre-planned infrastructure!"
    echo "Run: bash utilities/pre-calculate-infrastructure.sh"
    exit 504
fi
```

### Infrastructure Creation = READ AND EXECUTE:
```bash
# Just read from pre_planned_infrastructure
effort_config=$(yq ".pre_planned_infrastructure.efforts.$effort_id" orchestrator-state-v3.json)
full_path=$(echo "$effort_config" | yq '.full_path')
branch_name=$(echo "$effort_config" | yq '.branch_name')
remote_branch=$(echo "$effort_config" | yq '.remote_branch')

# No decisions, just execution
mkdir -p "$full_path"
cd "$full_path"
git clone "$TARGET_REPO" .
git checkout -b "$branch_name"
git push -u origin "$branch_name"
```

## 🔴🔴🔴 R360 - JUST-IN-TIME INFRASTRUCTURE EXECUTION 🔴🔴🔴

**TIMING ONLY - NOT NAMING (R504 handles naming)**

### Infrastructure Execution Timing:
```
Execute pre-planned infrastructure ONLY when:
- All dependencies are complete ✓
- Effort/split is next to be implemented ✓
- About to spawn agents for implementation ✓
```

---

## PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:
- R006: Orchestrator cannot write code
- R287: TODO persistence requirements
- R288: State file update protocol
- R322: Mandatory stop before state transitions
- R501: Progressive trunk-based development (SUPREME - CASCADE LAW)
- R504: Pre-infrastructure planning protocol (SUPREME)
- R509: Mandatory base branch validation (SUPREME - VALIDATION LAW)
- R510: Infrastructure creation protocol (SUPREME - CREATION LAW)
- R360: Just-in-time infrastructure execution
- R312: Git config immutability
- R340: Planning file metadata tracking
- R302: Comprehensive split tracking
- R502: Mandatory plan validation gates (CRITICAL)

## State Context

CREATE_NEXT_INFRASTRUCTURE = You are ACTIVELY creating infrastructure for the NEXT effort or split that needs it!

This could be:
- A regular effort at wave start
- A dependent effort after its dependency completes
- A split after previous split completes
- An independent effort that can run in parallel

## 🚨🚨🚨 IMMEDIATE ACTIONS UPON ENTERING STATE 🚨🚨🚨

**THE INSTANT YOU ENTER THIS STATE:**

### 0. MANDATORY PLAN VALIDATION (R502)

```bash
echo "🔍 VALIDATING REQUIRED PLANS EXIST (R502)..."

# Extract current phase and wave from state
CURRENT_PHASE=$(jq -r '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
CURRENT_WAVE=$(jq -r '.project_progression.current_wave.wave_number' orchestrator-state-v3.json)

# Validate phase plans exist using R550 canonical paths
# Method 1: Check orchestrator state tracking first (R550)
PHASE_ARCH_PLAN=$(jq -r ".planning_files.phases.phase${CURRENT_PHASE}.architecture_plan // empty" orchestrator-state-v3.json)
PHASE_TEST_PLAN=$(jq -r ".planning_files.phases.phase${CURRENT_PHASE}.test_plan // empty" orchestrator-state-v3.json)

# Method 2: Fallback to standard R550 paths if not tracked
if [ -z "$PHASE_ARCH_PLAN" ]; then
    PHASE_ARCH_PLAN="$CLAUDE_PROJECT_DIR/planning/phase${CURRENT_PHASE}/PHASE-ARCHITECTURE-PLAN.md"
fi
if [ -z "$PHASE_TEST_PLAN" ]; then
    PHASE_TEST_PLAN="$CLAUDE_PROJECT_DIR/planning/phase${CURRENT_PHASE}/PHASE-TEST-PLAN.md"
fi

# Validate phase plans exist
if [ ! -f "$PHASE_ARCH_PLAN" ]; then
    echo "🚨🚨🚨 CRITICAL: Phase ${CURRENT_PHASE} architecture plan missing!"
    echo "❌ Expected at: $PHASE_ARCH_PLAN"
    echo "Cannot create infrastructure without phase architecture plan!"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Missing phase ${CURRENT_PHASE} architecture plan (R550 violation)"
    ERROR_OCCURRED="true"
    exit 550
fi

# Validate wave plans exist using R550 canonical paths
# Method 1: Check orchestrator state tracking first (R550)
WAVE_IMPL_PLAN=$(jq -r ".planning_files.phases.phase${CURRENT_PHASE}.waves.wave${CURRENT_WAVE}.implementation_plan // empty" orchestrator-state-v3.json)
WAVE_TEST_PLAN=$(jq -r ".planning_files.phases.phase${CURRENT_PHASE}.waves.wave${CURRENT_WAVE}.test_plan // empty" orchestrator-state-v3.json)

# Method 2: Fallback to standard R550 paths if not tracked
if [ -z "$WAVE_IMPL_PLAN" ]; then
    WAVE_IMPL_PLAN="$CLAUDE_PROJECT_DIR/planning/phase${CURRENT_PHASE}/wave${CURRENT_WAVE}/WAVE-IMPLEMENTATION-PLAN.md"
fi
if [ -z "$WAVE_TEST_PLAN" ]; then
    WAVE_TEST_PLAN="$CLAUDE_PROJECT_DIR/planning/phase${CURRENT_PHASE}/wave${CURRENT_WAVE}/WAVE-TEST-PLAN.md"
fi

# Validate wave plans exist
if [ ! -f "$WAVE_IMPL_PLAN" ]; then
    echo "🚨🚨🚨 CRITICAL: Wave ${CURRENT_WAVE} implementation plan missing!"
    echo "❌ Expected at: $WAVE_IMPL_PLAN"
    echo "Cannot create infrastructure without wave implementation plan!"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Missing wave ${CURRENT_WAVE} implementation plan (R550 violation)"
    ERROR_OCCURRED="true"
    exit 550
fi

echo "✅ All required plans validated - proceeding with infrastructure creation"
echo "  Phase Architecture: $PHASE_ARCH_PLAN"
echo "  Wave Implementation: $WAVE_IMPL_PLAN"

# Document plan paths in orchestrator state if not already tracked (R550)
echo "📊 Documenting planning file paths (R550)..."
if [ -z "$(jq -r ".planning_files.phases.phase${CURRENT_PHASE}.architecture_plan // empty" orchestrator-state-v3.json)" ]; then
    jq ".planning_files.phases.phase${CURRENT_PHASE}.architecture_plan = \"planning/phase${CURRENT_PHASE}/PHASE-ARCHITECTURE-PLAN.md\"" \
       orchestrator-state-v3.json > orchestrator-state-v3.json.tmp
    mv orchestrator-state-v3.json.tmp orchestrator-state-v3.json
fi
if [ -z "$(jq -r ".planning_files.phases.phase${CURRENT_PHASE}.waves.wave${CURRENT_WAVE}.implementation_plan // empty" orchestrator-state-v3.json)" ]; then
    jq ".planning_files.phases.phase${CURRENT_PHASE}.waves.wave${CURRENT_WAVE}.implementation_plan = \"planning/phase${CURRENT_PHASE}/wave${CURRENT_WAVE}/WAVE-IMPLEMENTATION-PLAN.md\"" \
       orchestrator-state-v3.json > orchestrator-state-v3.json.tmp
    mv orchestrator-state-v3.json.tmp orchestrator-state-v3.json
fi
yq eval ".planning_files.planning_repo_files.phase_plans[\"phase${CURRENT_PHASE}\"].implementation.exists = true" -i orchestrator-state-v3.json
yq eval ".planning_files.planning_repo_files.phase_plans[\"phase${CURRENT_PHASE}\"].implementation.last_validated = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state-v3.json

yq eval ".planning_files.planning_repo_files.wave_plans[\"phase${CURRENT_PHASE}_wave${CURRENT_WAVE}\"].architecture.exists = true" -i orchestrator-state-v3.json
yq eval ".planning_files.planning_repo_files.wave_plans[\"phase${CURRENT_PHASE}_wave${CURRENT_WAVE}\"].architecture.last_validated = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state-v3.json
yq eval ".planning_files.planning_repo_files.wave_plans[\"phase${CURRENT_PHASE}_wave${CURRENT_WAVE}\"].implementation.exists = true" -i orchestrator-state-v3.json
yq eval ".planning_files.planning_repo_files.wave_plans[\"phase${CURRENT_PHASE}_wave${CURRENT_WAVE}\"].implementation.last_validated = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state-v3.json

echo "✅ Planning document tracking updated"
```

### 1. Validate Pre-Planned Infrastructure (R504)

```bash
echo "🔍 VALIDATING PRE-PLANNED INFRASTRUCTURE (R504)..."

# Check if pre_planned_infrastructure exists and is validated
if ! yq '.pre_planned_infrastructure.validated == true' orchestrator-state-v3.json > /dev/null 2>&1; then
    echo "🚨🚨🚨 CRITICAL: No validated pre-planned infrastructure!"
    echo "Infrastructure must be pre-calculated per R504"
    echo "Run: bash $CLAUDE_PROJECT_DIR/utilities/pre-calculate-infrastructure.sh"
    echo "Then: bash $CLAUDE_PROJECT_DIR/utilities/validate-infrastructure-naming.sh"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Missing pre-planned infrastructure (R504 violation)"
    ERROR_OCCURRED="true"
    exit 504
fi

echo "✅ Pre-planned infrastructure validated"
```

### 2. Determine What Needs Infrastructure

```bash
echo "🔧 DETERMINING NEXT INFRASTRUCTURE TO CREATE..."

# Get current phase and wave
CURRENT_PHASE=$(yq '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
CURRENT_WAVE=$(yq '.project_progression.current_wave.wave_number' orchestrator-state-v3.json)

# Find next uncreated effort from pre_planned_infrastructure
next_effort_id=$(yq '.pre_planned_infrastructure.efforts | to_entries[] |
  select(.value.phase == "phase'$CURRENT_PHASE'" and
         .value.wave == "wave'$CURRENT_WAVE'" and
         .value.created == false) |
  .key' orchestrator-state-v3.json | head -1)

if [ -n "$next_effort_id" ]; then
    echo "Creating infrastructure for effort: $next_effort_id"
    infrastructure_type="effort"
    infrastructure_target="$next_effort_id"
else
    echo "All pre-planned infrastructure for Phase $CURRENT_PHASE Wave $CURRENT_WAVE created!"
    # Transition to appropriate completion state
fi
```

### 3. Read Pre-Calculated Infrastructure (NO DECISIONS!)

```bash
echo "📖 READING PRE-CALCULATED INFRASTRUCTURE (R504)..."

# Read ALL configuration from pre_planned_infrastructure
effort_config=$(yq ".pre_planned_infrastructure.efforts.\"$infrastructure_target\"" orchestrator-state-v3.json)

# Extract pre-calculated values (NO RUNTIME DECISIONS!)
FULL_PATH=$(echo "$effort_config" | yq '.full_path')
BRANCH_NAME=$(echo "$effort_config" | yq '.branch_name')
REMOTE_BRANCH=$(echo "$effort_config" | yq '.remote_branch')
TARGET_REMOTE=$(echo "$effort_config" | yq '.target_remote')
TARGET_REPO_URL=$(echo "$effort_config" | yq '.target_repo_url')
PLANNING_REMOTE=$(echo "$effort_config" | yq '.planning_remote')
INTEGRATE_WAVE_EFFORTS_BRANCH=$(echo "$effort_config" | yq '.integration_branch')

echo "✅ Pre-calculated configuration loaded:"
echo "  Path: $FULL_PATH"
echo "  Branch: $BRANCH_NAME"
echo "  Remote: $REMOTE_BRANCH"
echo "  Integration: $INTEGRATE_WAVE_EFFORTS_BRANCH"
```

### 4. Execute Pre-Planned Infrastructure (CASCADE-AWARE PER R510!)

```bash
echo "🤖 EXECUTING CASCADE INFRASTRUCTURE (R501/R510)..."
echo "📁 Creating infrastructure at: $FULL_PATH"

# CRITICAL R510: Read base branch from pre_planned_infrastructure
BASE_BRANCH=$(echo "$effort_config" | yq '.base_branch')

if [ -z "$BASE_BRANCH" ] || [ "$BASE_BRANCH" = "null" ]; then
    echo "🚨 R510 VIOLATION: No base_branch in pre_planned_infrastructure!"
    echo "Infrastructure must specify cascade parent"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Missing base_branch (R510 violation)"
    ERROR_OCCURRED="true"
    exit 510
fi

echo "🔗 CASCADE: Creating from base branch: $BASE_BRANCH"

# R509: Validate cascade pattern
EFFORT_INDEX=$(echo "$effort_config" | yq '.index // 1')
if [[ "$CURRENT_PHASE" -eq 1 && "$CURRENT_WAVE" -eq 1 && "$EFFORT_INDEX" -eq 1 ]]; then
    # ONLY first effort of P1W1 can be from main
    if [ "$BASE_BRANCH" != "main" ]; then
        echo "🚨 R509 VIOLATION: First effort must be from main!"
        exit 509
    fi
else
    # ALL other efforts must follow cascade (NOT from main)
    if [ "$BASE_BRANCH" = "main" ]; then
        echo "🚨 R509 VIOLATION: Non-first effort cannot branch from main!"
        echo "Must follow cascade pattern per R501"
        exit 509
    fi
fi

# Create directory
mkdir -p "$FULL_PATH"

# R510: Clone ONLY the base branch (--single-branch flag MANDATORY!)
echo "📦 R510: Cloning ONLY base branch: $BASE_BRANCH"
git clone -b "$BASE_BRANCH" --single-branch "$TARGET_REPO_URL" "$FULL_PATH"

if [ $? -ne 0 ]; then
    echo "🚨 R510 VIOLATION: Cannot clone base branch $BASE_BRANCH"
    echo "Base branch might not exist or is inaccessible"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Cannot clone base branch (R510)"
    ERROR_OCCURRED="true"
    exit 510
fi

cd "$FULL_PATH"

# R509: Verify we got the right branch
CLONED_BRANCH=$(git branch --show-current)
if [ "$CLONED_BRANCH" != "$BASE_BRANCH" ]; then
    echo "🚨 R510 VIOLATION: Cloned wrong branch!"
    echo "Expected: $BASE_BRANCH, Got: $CLONED_BRANCH"
    exit 510
fi

# Create new branch FROM CASCADE PARENT
echo "🌿 Creating branch $BRANCH_NAME from cascade parent $BASE_BRANCH"
git checkout -b "$BRANCH_NAME"

# R509: Validate cascade relationship
if [ "$BASE_BRANCH" != "main" ]; then
    if ! git merge-base --is-ancestor "origin/$BASE_BRANCH" HEAD; then
        echo "🚨 R509 VIOLATION: Branch not based on cascade parent!"
        exit 509
    fi
fi

# ⚠️ CRITICAL: Push with -u flag to set upstream tracking
# Missing -u = validation failure and expensive recovery loop!
echo "📤 Pushing branch to remote WITH upstream tracking..."
git push -u origin "$BRANCH_NAME"

if [ $? -ne 0 ]; then
    echo "❌ CRITICAL: Failed to push branch to remote!"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Cannot push branch $BRANCH_NAME to remote"
    ERROR_OCCURRED="true"
    exit 1
fi

# Verify upstream tracking was set
UPSTREAM=$(git rev-parse --abbrev-ref --symbolic-full-name @{upstream} 2>/dev/null || echo "NONE")
if [ "$UPSTREAM" = "NONE" ]; then
    echo "❌ CRITICAL: Upstream tracking NOT configured despite -u flag!"
    echo "This should not happen - investigate git configuration"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Upstream tracking failed for $BRANCH_NAME"
    ERROR_OCCURRED="true"
    exit 1
fi

echo "✅ CASCADE INFRASTRUCTURE CREATED:"
echo "  Branch: $BRANCH_NAME"
echo "  Based on: $BASE_BRANCH"
echo "  Upstream: $UPSTREAM"
echo "  Following R501 progressive cascade"

# Lock git config (R312)
echo "🔒 Locking git config per R312..."
chmod 444 .git/config

# Install branch validation pre-commit hook (MANDATORY)
echo "🔒 Installing branch validation pre-commit hook..."
if [ -f "$CLAUDE_PROJECT_DIR/utilities/install-branch-validation-hook.sh" ]; then
    bash "$CLAUDE_PROJECT_DIR/utilities/install-branch-validation-hook.sh" "$INFRA_DIR" single false
    if [ $? -eq 0 ]; then
        echo "✅ Branch validation hook installed successfully"
    else
        echo "⚠️ WARNING: Branch validation hook installation failed"
    fi
else
    echo "⚠️ WARNING: Branch validation hook installer not found"
fi

# Copy plan if it exists (R340)
if [ "$infrastructure_type" = "effort" ]; then
    plan_location=$(jq -r ".effort_repo_files.effort_plans.$infrastructure_target" orchestrator-state-v3.json)
else
    plan_location=$(jq -r ".effort_repo_files.split_plans.$infrastructure_target" orchestrator-state-v3.json)
fi

if [ -n "$plan_location" ] && [ "$plan_location" != "null" ]; then
    echo "📋 Plan will be used from: $plan_location"
    # SW Engineer will read the plan from tracked location
fi
```

### 5. MANDATORY POST-CREATION VALIDATIONS (BLOCKING)

```bash
echo "🔍 PERFORMING MANDATORY POST-CREATION VALIDATIONS..."

# VALIDATION 1: R508 Repository URL Validation (BLOCKING)
echo "🔍 R508: Validating repository URL..."
EXPECTED_REPO=$(yq '.pre_planned_infrastructure.target_repo_url' orchestrator-state-v3.json)
EFFORT_REPO=$(yq ".pre_planned_infrastructure.efforts.\"$infrastructure_target\".target_repo_url" orchestrator-state-v3.json)

if [ "$EFFORT_REPO" != "$EXPECTED_REPO" ]; then
    echo "❌ R508 SUPREME LAW VIOLATION: Wrong repository URL in effort infrastructure"
    echo "   Expected: $EXPECTED_REPO"
    echo "   Found: $EFFORT_REPO"
    echo "   Effort: $infrastructure_target"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="R508 violation - wrong repository URL (expected $EXPECTED_REPO, found $EFFORT_REPO)"
    ERROR_OCCURRED="true"
    exit 911  # R508 catastrophic failure
fi

echo "✅ R508: Repository URL validated ($EXPECTED_REPO)"

# VALIDATION 2: Git Remote Configuration Validation (BLOCKING)
echo "🔍 Validating git remotes are configured..."
cd "$FULL_PATH"

ORIGIN_REMOTE=$(git remote get-url origin 2>/dev/null || echo "MISSING")
TARGET_REMOTE=$(git remote get-url "$TARGET_REMOTE" 2>/dev/null || echo "MISSING")

if [ "$ORIGIN_REMOTE" = "MISSING" ] || [ "$TARGET_REMOTE" = "MISSING" ]; then
    echo "❌ VALIDATION FAILURE: Git remotes not configured in $FULL_PATH"
    echo "   Origin remote: $ORIGIN_REMOTE"
    echo "   Target remote: $TARGET_REMOTE"
    echo "   Expected origin: (push URL)"
    echo "   Expected target: $EXPECTED_REPO"
    # DO NOT mark created=true until remotes exist
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Git remotes not configured for $infrastructure_target"
    ERROR_OCCURRED="true"
    exit 1
fi

echo "✅ Git remotes validated:"
echo "   Origin: $ORIGIN_REMOTE"
echo "   Target: $TARGET_REMOTE"

# VALIDATION 3: Branch Name Validation (BLOCKING)
echo "🔍 Validating correct branch is checked out..."
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
EXPECTED_BRANCH="$BRANCH_NAME"

if [ "$CURRENT_BRANCH" != "$EXPECTED_BRANCH" ]; then
    echo "❌ VALIDATION FAILURE: Wrong branch checked out"
    echo "   Expected: $EXPECTED_BRANCH"
    echo "   Current: $CURRENT_BRANCH"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Wrong branch checked out in $infrastructure_target"
    ERROR_OCCURRED="true"
    exit 1
fi

echo "✅ Branch validated: $CURRENT_BRANCH"

# VALIDATION 4: Remote Branch Exists (BLOCKING)
echo "🔍 Validating remote branch was pushed..."
if ! git ls-remote --heads origin "$BRANCH_NAME" | grep -q "$BRANCH_NAME"; then
    echo "❌ VALIDATION FAILURE: Remote branch not found"
    echo "   Expected remote branch: origin/$BRANCH_NAME"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Remote branch not pushed for $infrastructure_target"
    ERROR_OCCURRED="true"
    exit 1
fi

echo "✅ Remote branch validated: origin/$BRANCH_NAME"

echo "✅ ALL POST-CREATION VALIDATIONS PASSED"
echo ""
```

**CRITICAL**: ONLY AFTER ALL 4 VALIDATIONS PASS may orchestrator proceed to update tracking.

### 6. Update Tracking

```bash
# Mark as created AND validated in pre_planned_infrastructure (R504)
echo "📝 Updating pre_planned_infrastructure tracking..."
yq -i ".pre_planned_infrastructure.efforts.\"$infrastructure_target\".created = true" orchestrator-state-v3.json
yq -i ".pre_planned_infrastructure.efforts.\"$infrastructure_target\".validated = true" orchestrator-state-v3.json
yq -i ".pre_planned_infrastructure.efforts.\"$infrastructure_target\".validation_failure_reason = null" orchestrator-state-v3.json

# Also update legacy tracking for compatibility
if [ "$infrastructure_type" = "effort" ]; then
    # Extract effort name from effort_id (e.g., phase1_wave1_effort-name -> effort-name)
    effort_name=$(echo "$infrastructure_target" | sed 's/^phase[0-9]*_wave[0-9]*_//')

    # Update effort tracking if it exists
    if yq ".effort_dependencies.\"$effort_name\"" orchestrator-state-v3.json | grep -q "null"; then
        echo "Note: No legacy effort_dependencies entry for $effort_name"
    else
        yq -i ".effort_dependencies.\"$effort_name\".infrastructure_created = true |
               .effort_dependencies.\"$effort_name\".branch = \"$BRANCH_NAME\" |
               .effort_dependencies.\"$effort_name\".status = \"ready\"" orchestrator-state-v3.json
    fi
fi

# Commit state changes
git add orchestrator-state-v3.json
git commit -m "state: Created infrastructure for $infrastructure_target [R360]"
git push
```

## ❌ Common Mistakes (Anti-Patterns)

### WRONG - Missing `-u` Flag
```bash
git push origin "$BRANCH_NAME"  # ❌ NO! Missing -u flag!
```
**Result:** Branch pushed but no upstream tracking → validation failure → recovery loop → wasted $0.40-0.50

**Impact:** Test 2 took 4 iterations and 68 minutes due to this mistake (should have been 1 iteration, 40 minutes)

### WRONG - Separate Tracking Setup
```bash
git push origin "$BRANCH_NAME"
git branch --set-upstream-to=origin/"$BRANCH_NAME"  # ❌ Extra step, error-prone
```
**Result:** Two commands where one would do - complexity increases failure risk

### RIGHT - Single Command with `-u`
```bash
git push -u origin "$BRANCH_NAME"  # ✅ YES! One command, sets everything up
```
**Result:** Branch pushed AND upstream tracking configured in one atomic operation

**Validation:** Verify with `git branch -vv` to see `[origin/branch-name]` in output

---

## State Exit Criteria

Infrastructure is successfully created when:
1. ✅ Repository cloned to correct location
2. ✅ Branch created from correct base (dependency or wave integration)
3. ✅ Branch pushed to remote **WITH upstream tracking** (`-u` flag used)
4. ✅ Git config locked (R312)
5. ✅ Tracking updated in orchestrator-state-v3.json
6. ✅ State committed and pushed

## Next State Transitions

### VALID TRANSITIONS FROM CREATE_NEXT_INFRASTRUCTURE (from state-machines/software-factory-3.0-state-machine.json):
1. **ANALYZE_CODE_REVIEWER_PARALLELIZATION** - When need to determine how to parallelize effort planning
2. **SPAWN_SW_ENGINEERS** - When ready to spawn software engineers for implementation

### Decision Logic:
After successfully creating infrastructure:
- **If creating effort infrastructure AND effort plans don't exist** → ANALYZE_CODE_REVIEWER_PARALLELIZATION
- **If creating effort infrastructure AND effort plans exist** → SPAWN_SW_ENGINEERS
- **If split infrastructure created** → SPAWN_SW_ENGINEERS (splits already have plans)
- **If no more infrastructure needed** → WAVE_COMPLETE

## Common Mistakes to Avoid

1. ❌ **MOST COMMON:** Forgetting `-u` flag when pushing branches (causes validation failures)
2. ❌ Making ANY naming/pathing decisions at runtime (R504 violation)
3. ❌ Not validating pre_planned_infrastructure exists
4. ❌ Calculating branch names instead of reading them
5. ❌ Creating paths instead of using pre-calculated ones
6. ❌ Forgetting to mark infrastructure as created in pre_planned_infrastructure
7. ❌ Not running pre-calculate-infrastructure.sh before this state
8. ❌ Not verifying upstream tracking after push

## Validation

Before transitioning to next state:
```bash
# Verify infrastructure was created
if [ ! -d "$INFRA_DIR/.git" ]; then
    echo "ERROR: Infrastructure directory not created!"
    exit 1
fi

# Verify correct branch
cd "$INFRA_DIR"
current_branch=$(git branch --show-current)
if [ "$current_branch" != "$BRANCH_NAME" ]; then
    echo "ERROR: Wrong branch!"
    exit 1
fi

# ⚠️ CRITICAL: Verify upstream tracking is configured
UPSTREAM=$(git rev-parse --abbrev-ref --symbolic-full-name @{upstream} 2>/dev/null || echo "NONE")
if [ "$UPSTREAM" = "NONE" ]; then
    echo "ERROR: Upstream tracking NOT configured!"
    echo "This means you forgot the -u flag in git push!"
    echo "VALIDATION WILL FAIL - This causes expensive recovery loops!"
    exit 1
fi
echo "✅ Upstream tracking verified: $UPSTREAM"

# Verify git config is locked
if [ -w .git/config ]; then
    echo "ERROR: Git config not locked!"
    exit 1
fi
```

## Error Recovery

If infrastructure creation fails:
1. Check target repository access
2. Verify base branch exists
3. Ensure no duplicate branches
4. Clean up partial infrastructure
5. Update state to ERROR_RECOVERY


## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete CREATE_NEXT_INFRASTRUCTURE:**

### Phase 1: Prepare Results (Steps 1-2)

#### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

#### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, set variables for State Manager
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="[REASON_FOR_TRANSITION]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for State Transition (R288 - SUPREME LAW)
```bash
# State Manager handles ALL state file updates atomically
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE_NAME" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"

# State Manager will:
# 1. Validate the transition against state machine
# 2. Update all 4 state tracking locations atomically:
#    - orchestrator-state-v3.json
#    - orchestrator-state-demo.json
#    - .cascade-state-backup.json
#    - .orchestrator-state-v3.json
# 3. Commit and push all changes
# 4. Return control

echo "✅ State Manager completed transition"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before exit (R287 trigger)
save_todos "STATE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - state complete [R287]"; then
    echo "❌ ERROR: Failed to commit TODO files"
    echo "This is non-fatal but TODOs may be lost in compaction"
    echo "Proceeding with state execution..."
    # Don't exit - TODO commit failure is not fatal
fi

git push || echo "⚠️ WARNING: TODO push failed - committed locally"
echo "✅ TODOs saved and committed"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 5: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 6: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No State Manager spawn = state machine broken (R288 violation, -100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)


## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS (SF 3.0)

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="NEXT_STATE"
TRANSITION_REASON="State work complete"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"
# State Manager updates all 4 state files atomically

# 4. Save TODOs
save_todos "R322_CHECKPOINT"

# 5. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# 6. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

