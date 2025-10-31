---
name: R182 - Verify Git Repository Exists
criticality: BLOCKING
agent: all
state: INIT
---

# 🚨🚨🚨 RULE R182 - Verify Git Repository Exists 🚨🚨🚨

## Rule Statement
ALL agents MUST verify they are in a valid Git repository before beginning any work. This verification must occur during the INIT state pre-flight checks.

## Rationale
Working outside a Git repository will cause:
- Loss of version control
- Inability to create branches
- Failure to track changes
- Integration problems

## Implementation Requirements

### 1. Git Repository Check
```bash
# CHECK 2: VERIFY GIT REPOSITORY EXISTS
echo "Checking Git repository..."
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "🚨🚨🚨 CRITICAL: Not in a Git repository!"
    echo "Cannot proceed without Git repository"
    exit 182
fi
echo "✅ Git repository verified"
```

### 2. Repository Information
```bash
# Get repository details
REPO_ROOT=$(git rev-parse --show-toplevel)
GIT_DIR=$(git rev-parse --git-dir)
echo "Repository root: $REPO_ROOT"
echo "Git directory: $GIT_DIR"
```

### 3. Remote Configuration Check
```bash
# Verify remotes configured
if ! git remote -v | grep -q origin; then
    echo "⚠️⚠️⚠️ WARNING: No origin remote configured"
fi
```

## Enforcement
- **Trigger**: During INIT state pre-flight checks
- **Validation**: git rev-parse succeeds
- **Failure**: Exit with code 182

## Related Rules
- R184: Verify Git Branch
- R191: Target Repository Configuration
- R194: Remote Branch Tracking

---
**Status**: STUB - This rule file needs complete implementation details
**Created**: By software-factory-manager during rule library audit