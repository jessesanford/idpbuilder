# CREATE_NEXT_INFRASTRUCTURE State Rules

## State Context

**Current Phase**: Read from `orchestrator-state-v3.json`
**Current Wave**: Read from `orchestrator-state-v3.json`
**Entry From**: WAVE_START, WAITING_FOR_ARCHITECTURE_PLAN
**Exit To**: VALIDATE_INFRASTRUCTURE (MANDATORY - never skip per R507/R508!)

## 🚨🚨🚨 PRIMARY OBJECTIVE 🚨🚨🚨

**CREATE all git infrastructure (branches, directories, metadata) for efforts in the current wave using ONLY pre_planned_infrastructure per R504.**

This state is the FOUNDATION for the entire wave - it creates the git branches and directory structures that all agents will work in.

## 🔴🔴🔴 SUPREME LAW: R504 PRE-PLANNED INFRASTRUCTURE ENFORCEMENT 🔴🔴🔴

**THIS STATE MUST ONLY USE PRE-PLANNED INFRASTRUCTURE - NO RUNTIME DECISIONS!**

- ❌ **FORBIDDEN**: Making ANY naming decisions in this state
- ❌ **FORBIDDEN**: Calculating paths at runtime
- ❌ **FORBIDDEN**: Determining branch names dynamically
- ❌ **FORBIDDEN**: Creating infrastructure not in pre_planned_infrastructure
- ✅ **REQUIRED**: Use ONLY pre_planned_infrastructure from orchestrator-state-v3.json
- ✅ **REQUIRED**: Verify pre_planned_infrastructure exists and is validated
- ✅ **REQUIRED**: Fail immediately if infrastructure not pre-planned

## Required Inputs

### 1. Wave Plan (Must Exist)
```bash
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

# Find wave implementation plan
WAVE_PLAN=$(ls -t "$CLAUDE_PROJECT_DIR/phase-plans/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-PLAN--"*.md 2>/dev/null | head -1)

if [ -z "$WAVE_PLAN" ] || [ ! -f "$WAVE_PLAN" ]; then
    echo "❌ FATAL: Wave plan not found"
    echo "  Expected: phase-plans/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-PLAN--*.md"
    echo "  This state requires wave plan to exist first"
    exit 1
fi

echo "✅ Wave plan found: $WAVE_PLAN"
```

### 2. Extract Efforts from Wave Plan
```bash
# Parse efforts from wave plan
# Look for effort definitions in parallelization strategy or effort breakdown sections
EFFORTS=$(grep -E "^\\s*-\\s+effort[-_]" "$WAVE_PLAN" | sed -E 's/.*effort[-_]([0-9]+).*/effort-\1/' | sort -u)

if [ -z "$EFFORTS" ]; then
    echo "❌ FATAL: No efforts found in wave plan"
    exit 1
fi

echo "📋 Efforts to create infrastructure for:"
echo "$EFFORTS"
```

### 3. Target Repository Configuration
```bash
if [ ! -f "$CLAUDE_PROJECT_DIR/target-repo-config.yaml" ]; then
    echo "❌ FATAL: target-repo-config.yaml not found"
    exit 1
fi

TARGET_REPO=$(yq -r '.target_repository' "$CLAUDE_PROJECT_DIR/target-repo-config.yaml")
PROJECT_PREFIX=$(yq -r '.project_prefix' "$CLAUDE_PROJECT_DIR/target-repo-config.yaml")
```

## 🔴🔴🔴 INFRASTRUCTURE CREATION PROTOCOL 🔴🔴🔴

### Step 1: Verify Pre-Planned Infrastructure Exists (R504 ENFORCEMENT)
```bash
echo "🔍 Verifying pre-planned infrastructure per R504..."

# Check that pre_planned_infrastructure exists and is validated
PRE_PLANNED=$(jq -r '.pre_planned_infrastructure // empty' orchestrator-state-v3.json)

if [ -z "$PRE_PLANNED" ]; then
    echo "❌ FATAL: No pre_planned_infrastructure found in orchestrator-state-v3.json!"
    echo "  R504 VIOLATION: Infrastructure must be pre-planned during PROJECT_PLANNING/PHASE_PLANNING/WAVE_PLANNING"
    echo "  This state CANNOT make infrastructure decisions!"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

# Check if validated
VALIDATED=$(jq -r '.pre_planned_infrastructure.validated // false' orchestrator-state-v3.json)

if [ "$VALIDATED" != "true" ]; then
    echo "❌ FATAL: pre_planned_infrastructure not validated!"
    echo "  R504 VIOLATION: Infrastructure must be validated before creation"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

echo "✅ Pre-planned infrastructure found and validated"
```

### Step 2: Extract Pre-Planned Infrastructure for Current Wave
```bash
echo "📖 Reading pre-planned infrastructure for Phase ${PHASE} Wave ${WAVE}..."

# Get all efforts for current phase/wave from pre_planned_infrastructure
EFFORT_KEYS=$(jq -r ".pre_planned_infrastructure.efforts | to_entries[] | select(.value.phase == \"phase${PHASE}\" and .value.wave == \"wave${WAVE}\") | .key" orchestrator-state-v3.json)

if [ -z "$EFFORT_KEYS" ]; then
    echo "❌ FATAL: No pre-planned efforts found for Phase ${PHASE} Wave ${WAVE}"
    echo "  R504 VIOLATION: All infrastructure must be pre-planned"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

echo "📋 Pre-planned efforts to create infrastructure for:"
echo "$EFFORT_KEYS"
```

### Step 3: Create Effort Directories (FROM PRE-PLANNED DATA ONLY)
```bash
echo "📁 Creating effort directories from pre-planned infrastructure..."

for effort_key in $EFFORT_KEYS; do
    # READ from pre-planned infrastructure - NO CALCULATIONS!
    EFFORT_DIR=$(jq -r ".pre_planned_infrastructure.efforts.\"${effort_key}\".full_path" orchestrator-state-v3.json)
    EFFORT_NAME=$(jq -r ".pre_planned_infrastructure.efforts.\"${effort_key}\".effort_name" orchestrator-state-v3.json)
    EFFORT_BRANCH=$(jq -r ".pre_planned_infrastructure.efforts.\"${effort_key}\".branch_name" orchestrator-state-v3.json)

    if [ "$EFFORT_DIR" == "null" ] || [ -z "$EFFORT_DIR" ]; then
        echo "❌ FATAL: No full_path in pre_planned_infrastructure for $effort_key"
        echo "  R504 VIOLATION: All paths must be pre-planned"
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
        exit 1
    fi

    if [ ! -d "$EFFORT_DIR" ]; then
        mkdir -p "$EFFORT_DIR"
        echo "✅ Created: $EFFORT_DIR"
    else
        echo "⚠️ Already exists: $EFFORT_DIR"
    fi

    # Create .software-factory metadata directory (path from pre-planned)
    METADATA_DIR="$EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
    mkdir -p "$METADATA_DIR"
    echo "✅ Created metadata directory: $METADATA_DIR"

    # Create work-log file
    TIMESTAMP=$(date +%Y%m%d-%H%M%S)
    WORKLOG="$METADATA_DIR/work-log--${TIMESTAMP}.log"

    cat > "$WORKLOG" << WORKLOG_EOF
# Work Log: ${EFFORT_NAME}
Phase: ${PHASE}
Wave: ${WAVE}
Effort: ${EFFORT_NAME}
Started: $(date -Iseconds)
Branch: ${EFFORT_BRANCH}
Pre-planned: true

## Status
INFRASTRUCTURE_CREATED

## Activity Log
- $(date -Iseconds): Infrastructure created by orchestrator from pre-planned data
WORKLOG_EOF

    echo "✅ Created work log: $WORKLOG"

    # Mark as created in pre_planned_infrastructure
    jq --arg key "$effort_key" \
       '.pre_planned_infrastructure.efforts[$key].created = true' \
       orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
done
```

### Step 4: Create and Push Git Branches (FROM PRE-PLANNED DATA ONLY)
```bash
echo "🌿 Creating git branches from pre-planned infrastructure..."

# Read from pre-planned infrastructure - NO CALCULATIONS!
for effort_key in $EFFORT_KEYS; do
    # ALL DATA MUST COME FROM PRE_PLANNED_INFRASTRUCTURE
    EFFORT_BRANCH=$(jq -r ".pre_planned_infrastructure.efforts.\"${effort_key}\".branch_name" orchestrator-state-v3.json)
    REMOTE_BRANCH=$(jq -r ".pre_planned_infrastructure.efforts.\"${effort_key}\".remote_branch" orchestrator-state-v3.json)
    TARGET_REMOTE=$(jq -r ".pre_planned_infrastructure.efforts.\"${effort_key}\".target_remote" orchestrator-state-v3.json)
    EFFORT_DIR=$(jq -r ".pre_planned_infrastructure.efforts.\"${effort_key}\".full_path" orchestrator-state-v3.json)

    if [ "$EFFORT_BRANCH" == "null" ] || [ -z "$EFFORT_BRANCH" ]; then
        echo "❌ FATAL: No branch_name in pre_planned_infrastructure for $effort_key"
        echo "  R504 VIOLATION: All branch names must be pre-planned"
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
        exit 1
    fi

    echo "Creating branch: $EFFORT_BRANCH"

    # Navigate to effort directory to create git repo there
    cd "$EFFORT_DIR"

    # Initialize git repo if not exists
    if [ ! -d ".git" ]; then
        git init
        echo "✅ Initialized git repository in $EFFORT_DIR"
    fi

    # Add target remote from pre-planned data
    if [ "$TARGET_REMOTE" != "null" ] && [ -n "$TARGET_REMOTE" ]; then
        # Get target URL from pre-planned infrastructure (R504 requirement)
        TARGET_URL=$(jq -r ".pre_planned_infrastructure.efforts.\"${effort_key}\".target_repo_url" orchestrator-state-v3.json)

        if [ "$TARGET_URL" == "null" ] || [ -z "$TARGET_URL" ]; then
            echo "❌ FATAL: No target_repo_url in pre_planned_infrastructure for $effort_key"
            echo "  R504 VIOLATION: Target repository URL must be pre-planned for each effort"
            echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
            exit 1
        fi

        if ! git remote | grep -q "^${TARGET_REMOTE}$"; then
            git remote add "$TARGET_REMOTE" "$TARGET_URL"
            echo "✅ Added remote: $TARGET_REMOTE -> $TARGET_URL"
        fi

        # Fetch from target
        git fetch "$TARGET_REMOTE"

        # Create branch tracking the pre-planned remote branch
        if ! git show-ref --verify --quiet "refs/heads/$EFFORT_BRANCH"; then
            # Determine base branch (usually main for first effort, or integration branch)
            BASE_BRANCH="${TARGET_REMOTE}/main"  # Default to main

            # Create and checkout branch
            git checkout -b "$EFFORT_BRANCH" "$BASE_BRANCH"

            # Push to remote with tracking
            git push -u "$TARGET_REMOTE" "$EFFORT_BRANCH"
            echo "✅ Created and pushed: $EFFORT_BRANCH -> $REMOTE_BRANCH"
        else
            echo "⚠️ Branch already exists: $EFFORT_BRANCH"
        fi
    else
        echo "❌ FATAL: No target_remote in pre_planned_infrastructure for $effort_key"
        echo "  R504 VIOLATION: Remote configuration must be pre-planned"
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
        exit 1
    fi
done

# Return to project root
cd "$CLAUDE_PROJECT_DIR"
```

### Step 4.5: Install Git Hooks in Working Copies (MANDATORY - R383/R343/R506)
```bash
echo "🔒 Installing pre-commit hooks in effort working copies..."

# Install hooks in each effort working copy
for effort_key in $EFFORT_KEYS; do
    EFFORT_DIR=$(jq -r ".pre_planned_infrastructure.efforts.\"${effort_key}\".full_path" orchestrator-state-v3.json)
    EFFORT_NAME=$(jq -r ".pre_planned_infrastructure.efforts.\"${effort_key}\".effort_name" orchestrator-state-v3.json)

    # Verify effort has git repository
    if [ ! -d "$EFFORT_DIR/.git" ]; then
        echo "❌ FATAL: No git repository in $EFFORT_DIR"
        echo "  Git initialization should have occurred in Step 4"
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
        exit 1
    fi

    # Install hooks using dedicated installer script
    echo "Installing hooks for: $EFFORT_NAME"
    if bash "$CLAUDE_PROJECT_DIR/tools/install-effort-hooks.sh" "$EFFORT_DIR"; then
        echo "✅ Hooks installed: $EFFORT_DIR"

        # Mark hooks as installed in pre_planned_infrastructure
        jq --arg key "$effort_key" \
           '.pre_planned_infrastructure.efforts[$key].hooks_installed = true |
            .pre_planned_infrastructure.efforts[$key].hooks_installed_timestamp = (now | todate)' \
           orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
    else
        echo "❌ FATAL: Hook installation failed for $EFFORT_NAME ($EFFORT_DIR)"
        echo "  R383/R343/R506 ENFORCEMENT FAILURE"
        echo "  Cannot proceed without pre-commit validation hooks"
        echo ""
        echo "🔴🔴🔴 CRITICAL: Effort working copies MUST have hooks installed!"
        echo "Without hooks:"
        echo "  - R383 metadata placement validation: DISABLED"
        echo "  - R343 metadata directory enforcement: DISABLED"
        echo "  - R506 pre-commit bypass protection: DISABLED"
        echo "  - Metadata pollution risk: VERY HIGH"
        echo "  - Merge conflict risk: VERY HIGH"
        echo ""
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
        exit 1
    fi
done

echo "✅ All effort working copies have pre-commit hooks installed"
echo "   R383 metadata placement validation: ACTIVE"
echo "   R343 metadata directory standardization: ACTIVE"
echo "   R506 pre-commit bypass protection: ACTIVE"
```

### Step 5: Update State File with Infrastructure
```bash
echo "💾 Updating orchestrator state file..."

# Build infrastructure_created object
INFRA_JSON=$(jq -n '{
    phase: ($phase | tonumber),
    wave: ($wave | tonumber),
    timestamp: $timestamp,
    efforts: {}
}' \
    --arg phase "$PHASE" \
    --arg wave "$WAVE" \
    --arg timestamp "$(date -Iseconds)")

# Add each effort's infrastructure details
for effort in $EFFORTS; do
    BASE_BRANCH=$(jq -r ".efforts.\"${effort}\".base_branch" "$INFRA_PLAN")
    EFFORT_BRANCH=$(jq -r ".efforts.\"${effort}\".effort_branch" "$INFRA_PLAN")
    EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${effort}"

    INFRA_JSON=$(echo "$INFRA_JSON" | jq \
        --arg effort "$effort" \
        --arg base "$BASE_BRANCH" \
        --arg branch "$EFFORT_BRANCH" \
        --arg dir "$EFFORT_DIR" \
        '.efforts[$effort] = {
            base_branch: $base,
            effort_branch: $branch,
            directory: $dir,
            status: "READY_FOR_PLANNING"
        }')
done

# Update orchestrator-state-v3.json
jq --argjson infra "$INFRA_JSON" \
   '.infrastructure_created = $infra |
    .state_machine.current_state = "VALIDATE_INFRASTRUCTURE" |
    .state_machine.previous_state = "CREATE_NEXT_INFRASTRUCTURE" |
    .state_transition_log += [{
        "from": "CREATE_NEXT_INFRASTRUCTURE",
        "to": "VALIDATE_INFRASTRUCTURE",
        "timestamp": $infra.timestamp,
        "reason": "Infrastructure created, proceeding to mandatory validation (R507/R508)"
    }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

echo "✅ State file updated with infrastructure details"
```

## Validation Requirements

### Pre-Creation Validation
- ✅ Wave plan exists and is readable
- ✅ Efforts can be extracted from wave plan
- ✅ Target repo config exists
- ✅ Working directory is project root
- ✅ Git repository is clean

### Post-Creation Validation
```bash
echo "🔍 Validating infrastructure creation..."

VALIDATION_FAILED=false

for effort in $EFFORTS; do
    # Check directory exists
    EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${effort}"
    if [ ! -d "$EFFORT_DIR" ]; then
        echo "❌ Missing directory: $EFFORT_DIR"
        VALIDATION_FAILED=true
    fi

    # Check metadata directory exists
    METADATA_DIR="$EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}"
    if [ ! -d "$METADATA_DIR" ]; then
        echo "❌ Missing metadata directory: $METADATA_DIR"
        VALIDATION_FAILED=true
    fi

    # Check work log exists
    if [ -z "$(ls ${METADATA_DIR}/work-log--*.log 2>/dev/null)" ]; then
        echo "❌ Missing work log in: $METADATA_DIR"
        VALIDATION_FAILED=true
    fi

    # Check branch exists
    EFFORT_BRANCH="phase${PHASE}/wave${WAVE}/${effort}"
    if ! git show-ref --verify --quiet "refs/heads/$EFFORT_BRANCH"; then
        echo "❌ Missing branch: $EFFORT_BRANCH"
        VALIDATION_FAILED=true
    fi

    # Check remote tracking
    if ! git branch -vv | grep -q "$EFFORT_BRANCH.*\\[origin/"; then
        echo "❌ Branch not tracking remote: $EFFORT_BRANCH"
        VALIDATION_FAILED=true
    fi
done

if [ "$VALIDATION_FAILED" = true ]; then
    echo "❌ FATAL: Infrastructure validation failed"
    exit 1
fi

echo "✅ All infrastructure validation checks passed"
```

## Integration with Rules

- **R504**: Pre-Planned Infrastructure Protocol
- **R501**: Infrastructure Planning Requirements
- **R308**: Cascade Base Branch Calculation Algorithm
- **R383**: Metadata File Timestamp Requirements (enforced via hooks)
- **R343**: Metadata Directory Standardization (enforced via hooks)
- **R506**: Absolute Prohibition on Pre-Commit Bypass (enforced via hooks)
- **R193**: Git Branch Infrastructure
- **R209**: Effort Directory Isolation
- **R507/R508**: Mandatory Infrastructure Validation (NEXT STATE)

## Exit Criteria

Before transitioning to VALIDATE_INFRASTRUCTURE (MANDATORY NEXT STATE):
- ✅ All effort directories created
- ✅ All metadata directories created
- ✅ Work logs initialized
- ✅ Git branches created and pushed
- ✅ Remote tracking configured
- ✅ **Pre-commit hooks installed in all effort working copies (R383/R343/R506)**
- ✅ Infrastructure plan saved to state file
- ✅ All validation checks passed
- ⚠️ NEXT: Must go to VALIDATE_INFRASTRUCTURE for R507/R508 checks!

## Common Issues

### Issue: Base Branch Missing
**Detection**: Cannot find base branch for first effort
**Resolution**: Ensure phase integration branch exists, or use "main" for Phase 1 Wave 1

### Issue: Branch Already Exists
**Detection**: Git branch creation fails due to existing branch
**Resolution**: Verify branch state, potentially continue if already correct

### Issue: Permission Denied
**Detection**: Cannot push branch to remote
**Resolution**: Check git credentials, remote permissions

## Automation Flag

```bash
# After successfully creating all infrastructure:
echo "✅ Infrastructure created for all efforts"
echo "🔍 Next: VALIDATE_INFRASTRUCTURE (mandatory R507/R508 checks)"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue to validation
```

---

**REMEMBER**: This infrastructure MUST be created BEFORE any agent work begins. All agents depend on these branches and directories existing!

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
