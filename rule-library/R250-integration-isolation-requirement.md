# 🚨🚨🚨 RULE R250: Integration Isolation Requirement

## Classification
- **Category**: Integration Management
- **Criticality Level**: 🚨🚨🚨 BLOCKING
- **Enforcement**: MANDATORY for all integration work
- **Penalty**: -50% to -100% for violations

## The Rule

**ALL integration work MUST occur in completely isolated workspaces to prevent contamination of orchestrator and effort environments.**

## Requirements

### 1. Workspace Location (MANDATORY)
Integration MUST happen in isolated workspace under:
```
/efforts/phase{X}/wave{Y}/integration-workspace/
```

**NEVER perform integration in:**
- ❌ The orchestrator's main directory (`/home/vscode/software-factory-template/`)
- ❌ Any effort development workspace
- ❌ The Software Factory instance directory
- ❌ Any shared or common directories

### 2. Fresh Clone Requirements
The integration workspace MUST be:
- A fresh clone of the target repository
- Created specifically for this integration
- Based on the correct base branch (main/develop)
- Completely separate from all effort workspaces

### 3. Directory Structure
```
/efforts/
├── phase1/
│   └── wave1/
│       ├── effort1/              # Development workspace
│       ├── effort2/              # Development workspace  
│       └── integration-workspace/ # ISOLATED integration (fresh clone)
```

### 4. Integration Setup Protocol (FOLLOWS R104)
```bash
# CORRECT approach - Per R104, integration branches in TARGET repository
PHASE=1
WAVE=1
INTEGRATION_DIR="$CLAUDE_PROJECT_DIR/efforts/phase${PHASE}/wave${WAVE}/integration-workspace"

# Read target repository configuration (R104 requirement)
TARGET_CONFIG="$CLAUDE_PROJECT_DIR/target-repo-config.yaml"
TARGET_REPO_PATH=$(yq '.repository_path' "$TARGET_CONFIG")
TARGET_REPO_NAME=$(yq '.repository_name' "$TARGET_CONFIG")

# Clean any previous integration
rm -rf "$INTEGRATION_DIR"

# Create fresh clone of TARGET repository (per R104)
mkdir -p "$INTEGRATION_DIR"
cd "$INTEGRATION_DIR"

# Clone the TARGET repository (not software-factory!)
git clone "$TARGET_REPO_PATH" "$TARGET_REPO_NAME"
cd "$TARGET_REPO_NAME"

# CRITICAL SAFETY CHECK - Verify correct repository
REMOTE_URL=$(git remote get-url origin)
if [[ "$REMOTE_URL" == *"software-factory"* ]]; then
    echo "❌ CRITICAL: Cloned orchestrator repository instead of target!"
    echo "Expected: Target project repository"
    echo "Got: $REMOTE_URL"
    exit 1
fi

git checkout -b "wave-${WAVE}-integration"  # Per R104 naming convention
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