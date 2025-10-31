# 🔴🔴🔴 RULE R508 - Target Repository Enforcement [SUPREME LAW]

## Category: Repository Safety

## Priority: ABSOLUTE SUPREME

## Description
ALL efforts, splits, and integrations MUST use the repository specified in target-repo-config.yaml. Using ANY other repository is an IMMEDIATE and CATASTROPHIC failure.

## SUPREME LAW DECLARATION
```
THIS IS SUPREME LAW - NO EXCEPTIONS
All code MUST go to the configured target repository
Using wrong repository = IMMEDIATE -100% FAILURE
```

## Requirements

### REPOSITORY CONFIGURATION SOURCE
```yaml
# ONLY source of truth: target-repo-config.yaml
repository_url: "https://github.com/owner/correct-repo.git"
```

### ENFORCEMENT CHECKPOINTS
1. **Pre-Infrastructure Planning (R504)**
   - MUST read target-repo-config.yaml
   - MUST use this URL for ALL infrastructure

2. **Infrastructure Creation**
   - MUST set remote to target repository
   - NO hardcoded URLs allowed
   - **MUST use `-u` flag when pushing to set upstream tracking**

3. **Infrastructure Validation (R507)**
   - MUST verify remote matches config
   - MUST verify upstream tracking is configured
   - ANY deviation = ERROR_RECOVERY

### VALIDATION COMMANDS
```bash
# Extract configured repository
CONFIGURED_REPO=$(yq '.repository_url' $CLAUDE_PROJECT_DIR/target-repo-config.yaml)

# Check every effort/split/integration
cd $effort_dir
ACTUAL_REMOTE=$(git remote get-url origin 2>/dev/null || git remote get-url target 2>/dev/null)

if [ "$ACTUAL_REMOTE" != "$CONFIGURED_REPO" ]; then
    echo "🔴🔴🔴 CATASTROPHIC FAILURE: WRONG REPOSITORY!"
    echo "Expected: $CONFIGURED_REPO"
    echo "Actual: $ACTUAL_REMOTE"
    exit 911  # Emergency exit code
fi

# Verify upstream tracking is configured (prevents validation failures)
UPSTREAM=$(git rev-parse --abbrev-ref --symbolic-full-name @{upstream} 2>/dev/null || echo "NONE")
if [ "$UPSTREAM" = "NONE" ]; then
    echo "❌ CRITICAL: Upstream tracking NOT configured!"
    echo "This means the -u flag was missing from git push"
    echo "Correct command: git push -u origin \"\$BRANCH_NAME\""
    exit 1
fi
```

## Penalties
- Using wrong repository: -100% IMMEDIATE FAILURE
- Hardcoding repository URLs: -75%
- Not checking target-repo-config.yaml: -50%
- Allowing wrong remote to persist: -100%
- Missing `-u` flag in git push: -15% (causes expensive validation failures and recovery loops)

## Examples

### CORRECT: Using Target Repository with Upstream Tracking
```bash
# Read from config
REPO=$(yq '.repository_url' target-repo-config.yaml)
git remote add origin "$REPO"

# Push with -u flag to set upstream tracking
git push -u origin "$BRANCH_NAME"
```

### CATASTROPHIC: Using Wrong Repository
```bash
# 🔴 NEVER EVER DO THIS
git remote add origin "https://github.com/wrong/repo.git"
# This is -100% IMMEDIATE FAILURE
```

### CATASTROPHIC: Hardcoding Repository
```bash
# 🔴 NEVER HARDCODE
git clone https://github.com/some/repo.git
# MUST use target-repo-config.yaml
```

### WRONG: Missing `-u` Flag (Common Mistake)
```bash
# ❌ Missing -u flag
git push origin "$BRANCH_NAME"
# Result: Branch pushed but no upstream tracking
# Impact: VALIDATE_INFRASTRUCTURE fails → recovery loop → wasted time and money
# Fix: ALWAYS use: git push -u origin "$BRANCH_NAME"
```

## Detection and Recovery
1. **Detection**: validate-infrastructure.sh checks all remotes
2. **Alert**: IMMEDIATE stop with exit code 911
3. **Recovery**: ERROR_RECOVERY to fix remotes
4. **Prevention**: NEVER bypass validation

## Related Rules
- R507: Mandatory Infrastructure Validation
- R504: Pre-Infrastructure Planning Protocol
- R360: Just-in-Time Infrastructure

## Source
Created as SUPREME LAW to prevent code going to wrong repositories.

## Metadata
- Created: 2025-09-27
- Criticality: SUPREME LAW
- Priority: ABSOLUTE HIGHEST
- Enforcement: ZERO TOLERANCE