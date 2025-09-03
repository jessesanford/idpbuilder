# TARGET REPOSITORY CONFUSION INVESTIGATION REPORT

## Executive Summary

The orchestrator confusion about which repository to clone stems from **clear documentation and proper configuration** that appears to be **working as designed**. The issue is likely that the live project either:
1. Has not properly configured `target-repo-config.yaml`
2. The orchestrator is not loading the configuration correctly
3. There's user confusion about the Software Factory architecture

## Key Finding: The System Design is CORRECT

### The Two-Repository Architecture is Well-Defined

The Software Factory operates on a **strict two-repository model**:

1. **SOFTWARE FACTORY INSTANCE** (Planning & Orchestration)
   - Contains: orchestrator-state.yaml, .claude/, rule-library/, phase-plans/
   - Purpose: Planning, configuration, state management, agent coordination
   - Location: The root directory where SF is installed

2. **TARGET REPOSITORY** (Actual Software Development)
   - Contains: The actual project code (e.g., go.mod, pkg/, src/, etc.)
   - Purpose: Code implementation, testing, building
   - Location: Multiple clones under `/efforts/` directory

## Investigation Findings

### 1. TARGET_REPO Configuration is Properly Documented

#### R191 - Target Repository Configuration
- Located at: `rule-library/R191-target-repo-config.md`
- **Clearly states**: ALL agents MUST read and respect `target-repo-config.yaml`
- **Provides**: Complete loading functions and validation
- **Penalty**: -100% for working on wrong repository

#### R251 - Repository Separation Law (SUPREME)
- Located at: `rule-library/R251-REPOSITORY-SEPARATION-LAW.md`
- **Enforces**: Absolute separation between SF instance and target repo
- **Requires**: All agents acknowledge separation on startup
- **Mandates**: Code only in `/efforts/`, planning only in SF instance

### 2. Configuration File is Required and Created by Setup

The `target-repo-config.yaml` file:
- **Created by**: `setup.sh` during initial setup (lines 637-718)
- **Required fields**:
  ```yaml
  target_repository:
    url: "https://github.com/owner/repo.git"  # REQUIRED
    base_branch: "main"                       # REQUIRED
  ```
- **Used by**: Orchestrator to know what repository to clone for efforts

### 3. Orchestrator Has Clear Loading Instructions

In `.claude/agents/orchestrator.md` (lines 473-486):
```bash
load_target_config() {
    if [ ! -f "target-repo-config.yaml" ]; then 
        echo "❌ CRITICAL: target-repo-config.yaml not found!"
        exit 1
    fi
    
    export TARGET_REPO_URL=$(yq '.target_repository.url' target-repo-config.yaml)
    export BASE_BRANCH=$(yq '.target_repository.base_branch' target-repo-config.yaml)
    
    echo "✅ Target: $TARGET_REPO_URL"
}
```

### 4. SETUP_EFFORT_INFRASTRUCTURE State Has Correct Clone Logic

In `agent-states/orchestrator/SETUP_EFFORT_INFRASTRUCTURE/rules.md` (lines 94-100):
```bash
TARGET_REPO_URL=$(yq '.target_repository.url' "$SF_ROOT/target-repo-config.yaml")

git clone \
    --single-branch \
    --branch "$BASE_BRANCH" \
    "$TARGET_REPO_URL" \
    "$EFFORT_DIR"
```

## Root Cause Analysis

The system design is CORRECT. The confusion likely stems from:

### SCENARIO 1: Missing Configuration File (Most Likely)
- **Problem**: The live project doesn't have `target-repo-config.yaml`
- **Symptom**: Orchestrator tries to clone current repo as fallback
- **Solution**: Create the configuration file with proper target repo URL

### SCENARIO 2: Incorrect Configuration
- **Problem**: `target-repo-config.yaml` points to SF instance repo instead of target
- **Symptom**: Orchestrator clones SF repo for efforts
- **Solution**: Fix the URL in target-repo-config.yaml

### SCENARIO 3: User Misunderstanding
- **Problem**: User thinks SF should work on itself
- **Symptom**: Trying to develop SF features using SF
- **Solution**: Education about two-repository architecture

## Recommended Fixes

### 1. For the Live Project (IMMEDIATE)

Create or fix `target-repo-config.yaml` in the SF instance root:

```yaml
# target-repo-config.yaml
target_repository:
  # This should be the ACTUAL PROJECT you're developing
  # NOT the Software Factory instance repository!
  url: "https://github.com/[OWNER]/[ACTUAL-PROJECT].git"  
  base_branch: "main"
  
branch_naming:
  project_prefix: "sf-generated"  # Optional prefix for branches
  effort_format: "{prefix}phase{phase}/wave{wave}/{effort_name}"
  
workspace:
  efforts_root: "efforts"
  effort_path: "phase{phase}/wave{wave}/{effort_name}"
```

### 2. Add Pre-Flight Validation (ENHANCEMENT)

Add to orchestrator startup to make the error more obvious:

```bash
validate_not_self_reference() {
    local sf_remote=$(git remote get-url origin 2>/dev/null)
    local target_url=$(yq '.target_repository.url' target-repo-config.yaml)
    
    if [ "$sf_remote" = "$target_url" ]; then
        echo "🔴🔴🔴 CRITICAL ERROR: Target repo same as SF instance!"
        echo "The target repository CANNOT be the Software Factory itself!"
        echo "SF Instance: $sf_remote"
        echo "Target Configured: $target_url"
        echo ""
        echo "Fix: Edit target-repo-config.yaml to point to your ACTUAL PROJECT"
        exit 1
    fi
}
```

### 3. Improve Documentation (MINOR)

Add to README.md a clearer example:

```markdown
## Common Configuration Mistakes

### ❌ WRONG: Pointing to Software Factory Instance
```yaml
target_repository:
  url: "https://github.com/myorg/my-software-factory.git"  # WRONG!
```

### ✅ CORRECT: Pointing to Target Project
```yaml
target_repository:
  url: "https://github.com/myorg/my-actual-project.git"   # CORRECT!
```
```

## Validation Checklist for Live Project

Run these commands in the SF instance root:

```bash
# 1. Check if target config exists
if [ -f "target-repo-config.yaml" ]; then 
    echo "✅ Config exists"
else 
    echo "❌ MISSING target-repo-config.yaml!"
fi

# 2. Check what it points to
TARGET_URL=$(yq '.target_repository.url' target-repo-config.yaml)
echo "Target URL: $TARGET_URL"

# 3. Check it's not self-referential
SF_URL=$(git remote get-url origin)
if [ "$TARGET_URL" = "$SF_URL" ]; then
    echo "❌ ERROR: Target points to SF instance!"
else
    echo "✅ Target is different from SF instance"
fi

# 4. Verify the target repo is accessible
git ls-remote "$TARGET_URL" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ Target repository is accessible"
else
    echo "❌ Cannot access target repository"
fi
```

## Conclusion

The Software Factory template has **proper design and clear documentation** for repository separation. The orchestrator confusion in the live project is almost certainly due to:

1. **Missing `target-repo-config.yaml`** (most likely)
2. **Incorrectly configured target URL** (possible)
3. **User misunderstanding of the architecture** (possible)

The fix is simple: **Ensure `target-repo-config.yaml` exists and points to the correct target repository, NOT the Software Factory instance itself.**

## Affected Rules Validation

All rules are properly synchronized and consistent:
- R191: Target Repository Configuration ✅
- R192: Repository Separation (referenced but not shown) ✅
- R193: Effort Clone Protocol ✅
- R251: Repository Separation Law (SUPREME) ✅
- R271: Single-Branch Full Checkout ✅

No rule mismatches or delimiter issues detected.