# 🚨🚨🚨 RULE R250: Integration Isolation Requirement

## Classification
- **Category**: Integration Management
- **Criticality Level**: 🚨🚨🚨 BLOCKING
- **Enforcement**: MANDATORY for all integration work
- **Penalty**: -50% to -100% for violations

## The Rule

**ALL integration work MUST occur in completely isolated workspaces to prevent contamination of orchestrator and effort environments.**

## Requirements

### 1. Workspace Location (MANDATORY - EXACT SPECIFICATIONS)

#### 🔴🔴🔴 DETERMINISTIC DIRECTORY STRUCTURE 🔴🔴🔴

**EXACT STRUCTURE (NO AMBIGUITY):**
```bash
# SF_INSTANCE_DIR is the Software Factory instance directory
SF_INSTANCE_DIR="/home/vscode/software-factory-template"  # Or wherever SF is installed

# Integration workspace is ALWAYS a subdirectory, target repo cloned AS subdirectory
INTEGRATION_WORKSPACE="${SF_INSTANCE_DIR}/efforts/phase{X}/wave{Y}/integration-workspace"

# The actual clone goes INTO integration-workspace as 'repo' subdirectory
INTEGRATION_REPO="${INTEGRATION_WORKSPACE}/repo"

# EXACT STRUCTURE:
${SF_INSTANCE_DIR}/
└── efforts/
    └── phase2/
        └── wave1/
            └── integration-workspace/     # Created by mkdir -p
                └── repo/                   # Target repo cloned here
                    ├── .git/
                    ├── src/
                    └── ... (target repo contents)
```

**NEVER perform integration in:**
- ❌ The orchestrator's main directory (`/home/vscode/software-factory-template/`)
- ❌ Any effort development workspace
- ❌ The Software Factory instance directory root
- ❌ Any shared or common directories
- ❌ Directly as `integration-workspace` (must be `integration-workspace/repo`)

### 2. Fresh Clone Requirements
The integration workspace MUST be:
- A fresh clone of the target repository
- Created specifically for this integration
- Based on the correct base branch (main/develop)
- Completely separate from all effort workspaces

### 3. Directory Structure (DETERMINISTIC)

#### First Integration Attempt:
```
${SF_INSTANCE_DIR}/efforts/
├── phase1/
│   └── wave1/
│       ├── effort1/                    # Development workspace
│       ├── effort2/                    # Development workspace  
│       └── integration-workspace/      # Integration directory
│           └── repo/                   # Target repo cloned here
│               ├── .git/
│               ├── INTEGRATION-METADATA.md
│               └── ... (project files)
```

#### Re-integration After Fixes (DETERMINISTIC HANDLING):
```
${SF_INSTANCE_DIR}/efforts/
├── phase1/
│   └── wave1/
│       ├── effort1/                           # Development workspace
│       ├── effort2/                           # Development workspace  
│       ├── integration-workspace-archived-1/  # Previous attempt (renamed)
│       │   └── repo/                         # Old integration preserved
│       └── integration-workspace/             # Fresh re-integration
│           └── repo/                         # New clean clone
│               ├── .git/
│               ├── INTEGRATION-METADATA.md
│               └── ... (project files)
```

**RE-INTEGRATION PROTOCOL (DETERMINISTIC):**
1. If `integration-workspace` exists, rename to `integration-workspace-archived-N`
2. Create fresh `integration-workspace` directory
3. Clone target repo as `integration-workspace/repo`
4. Use SAME branch name (force-push if exists)

### 4. Integration Setup Protocol (DETERMINISTIC AND EXACT)
```bash
# 🔴🔴🔴 DETERMINISTIC INTEGRATION SETUP 🔴🔴🔴
setup_integration_infrastructure() {
    local PHASE=$1
    local WAVE=$2
    local INTEGRATION_TYPE=$3  # "wave", "phase", or "project"
    
    # EXACT PATHS (NO AMBIGUITY)
    SF_INSTANCE_DIR="$(pwd)"  # Must be in SF instance root
    INTEGRATION_BASE_DIR="${SF_INSTANCE_DIR}/efforts"
    
    # Determine exact integration path based on type
    case "$INTEGRATION_TYPE" in
        "wave")
            INTEGRATION_DIR="${INTEGRATION_BASE_DIR}/phase${PHASE}/wave${WAVE}/integration-workspace"
            BRANCH_NAME="phase${PHASE}-wave${WAVE}-integration"
            ;;
        "phase")
            INTEGRATION_DIR="${INTEGRATION_BASE_DIR}/phase${PHASE}/phase-integration-workspace"
            BRANCH_NAME="phase${PHASE}-integration"
            ;;
        "project")
            INTEGRATION_DIR="${INTEGRATION_BASE_DIR}/project-integration-workspace"
            BRANCH_NAME="project-integration"
            ;;
    esac
    
    # Handle re-integration (DETERMINISTIC)
    if [ -d "$INTEGRATION_DIR" ]; then
        # Archive old integration attempt
        ARCHIVE_NUM=1
        while [ -d "${INTEGRATION_DIR}-archived-${ARCHIVE_NUM}" ]; do
            ARCHIVE_NUM=$((ARCHIVE_NUM + 1))
        done
        echo "📦 Archiving previous integration attempt to ${INTEGRATION_DIR}-archived-${ARCHIVE_NUM}"
        mv "$INTEGRATION_DIR" "${INTEGRATION_DIR}-archived-${ARCHIVE_NUM}"
    fi
    
    # Create fresh integration directory
    mkdir -p "$INTEGRATION_DIR"
    
    # Read target repository configuration
    TARGET_REPO_URL=$(yq '.target_repository.url' "${SF_INSTANCE_DIR}/target-repo-config.yaml")
    
    # Determine base branch per R308
    BASE_BRANCH=$(determine_integration_base_branch "$INTEGRATION_TYPE" "$PHASE" "$WAVE")
    
    # Clone INTO integration-workspace as 'repo' subdirectory (DETERMINISTIC)
    echo "📦 Cloning target repo INTO ${INTEGRATION_DIR}/repo"
    git clone \
        --single-branch \
        --branch "$BASE_BRANCH" \
        "$TARGET_REPO_URL" \
        "${INTEGRATION_DIR}/repo"
    
    # Change to the repo directory
    cd "${INTEGRATION_DIR}/repo"
    
    # CRITICAL SAFETY CHECK - Verify correct repository
    REMOTE_URL=$(git remote get-url origin)
    if [[ "$REMOTE_URL" == *"software-factory"* ]]; then
        echo "❌ CRITICAL: Cloned orchestrator repository instead of target!"
        echo "Expected: Target project repository"
        echo "Got: $REMOTE_URL"
        exit 250
    fi
    
    # Create or force-update integration branch
    git checkout -b "$BRANCH_NAME"
    
    # Force push if branch exists (re-integration case)
    if git ls-remote --heads origin "$BRANCH_NAME" | grep -q "$BRANCH_NAME"; then
        echo "⚠️ Branch exists, will force-push after integration"
        git push --force-with-lease -u origin "$BRANCH_NAME"
    else
        git push -u origin "$BRANCH_NAME"
    fi
    
    echo "✅ Integration infrastructure ready at: ${INTEGRATION_DIR}/repo"
    echo "✅ Branch: $BRANCH_NAME (base: $BASE_BRANCH)"
}
```

### 5. Prohibited Actions
- ❌ Using orchestrator directory for integration
- ❌ Using any effort workspace for integration
- ❌ Reusing old integration workspaces
- ❌ Mixing development and integration in same directory
- ❌ Creating integration branches in effort workspaces

## Validation Checks
```bash
# Verify integration isolation
if [[ "$PWD" == *"/integration-workspace"* ]] && [[ "$PWD" == *"/efforts/phase"* ]]; then
    echo "✅ Integration workspace properly isolated"
else
    echo "🚨 VIOLATION: Integration not in isolated workspace!"
    exit 1
fi

# Verify it's a git repository
if [ ! -d ".git" ]; then
    echo "🚨 VIOLATION: Integration workspace is not a git repository!"
    exit 1
fi

# CRITICAL: Verify correct repository (NOT software-factory)
REMOTE_URL=$(git remote get-url origin)
if [[ "$REMOTE_URL" == *"software-factory"* ]]; then
    echo "❌ CRITICAL: Working in orchestrator repository instead of target!"
    echo "Expected: Target project repository"
    echo "Got: $REMOTE_URL"
    exit 1
fi
```

## Penalties
- Using orchestrator directory: **-100% grade** (CRITICAL FAILURE)
- Using effort workspace: **-50% grade** + corrupted state
- Reusing old workspace: **-40% grade** + merge conflicts
- Missing isolation: **-100% grade** (COMPLETE FAILURE)

## Why This Matters

Integration isolation prevents:
1. **Workspace Pollution**: Mixing integration with development corrupts both
2. **Branch Confusion**: Integration branches in wrong repositories
3. **State Corruption**: Orchestrator state mixed with integration state
4. **Merge Disasters**: Conflicting changes from multiple sources

## Related Rules
- R034: Integration Requirements
- R271: Full Checkout Requirement (no sparse checkouts)
- R014: Branch Naming Convention
- R288: State File Updates

## Common Violations to Avoid

### ❌ WRONG: Integration in orchestrator directory
```bash
cd /home/vscode/software-factory-template
git checkout -b integration  # WRONG LOCATION!
```

### ❌ WRONG: Integration in effort workspace
```bash
cd /efforts/phase1/wave1/effort1
git checkout -b integration  # CONTAMINATING EFFORT!
```

### ✅ CORRECT: Isolated integration workspace with TARGET repository
```bash
cd $CLAUDE_PROJECT_DIR/efforts/phase1/wave1/integration-workspace/[target-repo-name]
git checkout -b wave-1-integration  # CORRECT - Per R104!
```

## Remember
**Isolation is NOT optional** - it's a BLOCKING requirement. The integration workspace must be completely separate from all other workspaces to prevent catastrophic merge failures and state corruption.