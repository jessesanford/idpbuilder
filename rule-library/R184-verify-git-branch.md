---
name: R184 - Verify Git Branch
criticality: BLOCKING
agent: sw-engineer, code-reviewer
state: INIT
---

# 🚨🚨🚨 RULE R184 - Verify Git Branch 🚨🚨🚨

## Rule Statement
Agents MUST verify they are on the correct Git branch for their assigned work. The branch name must follow project naming conventions and match the expected pattern for the effort being implemented.

## Rationale
Working on the wrong branch causes:
- Code going to wrong location
- Merge conflicts
- Lost work
- Broken integration chains

## Implementation Requirements

### 1. Branch Verification
```bash
# CHECK 3: VERIFY GIT BRANCH (R184 + R191)
echo "Checking Git branch..."
CURRENT_BRANCH=$(git branch --show-current)
echo "Current branch: $CURRENT_BRANCH"

# Get expected branch pattern from environment
EXPECTED_PATTERN="${PROJECT_PREFIX}-phase${PHASE}-wave${WAVE}-effort${EFFORT}"

if [[ "$CURRENT_BRANCH" != *"$EXPECTED_PATTERN"* ]]; then
    echo "🚨🚨🚨 CRITICAL: Wrong branch!"
    echo "Expected pattern: $EXPECTED_PATTERN"
    echo "Current branch: $CURRENT_BRANCH"
    exit 184
fi
echo "✅ Branch verified: $CURRENT_BRANCH"
```

### 2. Branch Naming Convention (with R191)
```bash
# Project prefix required (R191)
PROJECT_PREFIX=$(grep "project_prefix:" orchestrator-state.yaml | awk '{print $2}')

# Branch format
# Feature: ${PROJECT_PREFIX}-phase${N}-wave${N}-effort${N}
# Split: ${PROJECT_PREFIX}-phase${N}-wave${N}-effort${N}-split${NNN}
# Integration: ${PROJECT_PREFIX}-phase${N}-wave${N}-integration
```

### 3. Remote Tracking Verification
```bash
# Ensure branch tracks remote
TRACKING=$(git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null)
if [ -z "$TRACKING" ]; then
    echo "⚠️⚠️⚠️ WARNING: Branch not tracking remote"
fi
```

## Enforcement
- **Trigger**: During INIT state pre-flight checks
- **Validation**: Branch name matches expected pattern
- **Failure**: Exit with code 184

## Related Rules
- R014: Branch Naming Convention
- R191: Target Repository Configuration (project prefix)
- R194: Remote Branch Tracking

---
**Status**: STUB - This rule file needs complete implementation details
**Created**: By software-factory-manager during rule library audit