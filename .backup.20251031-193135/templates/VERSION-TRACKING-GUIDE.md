# Version Tracking Guide for R381/R382 Compliance

## 🔴🔴🔴 CRITICAL: Version Immutability is LAW 🔴🔴🔴

This guide explains how to track and enforce library version consistency per R381 and R382.

## Quick Start

1. **Copy the template** to your orchestrator-state-v3.json:
```bash
# Add version_tracking section from template
cat templates/version-tracking-template.json | jq '.version_tracking' >> orchestrator-state-v3.json
```

2. **Document ALL existing versions** immediately:
```bash
# For Go projects
grep -E "^\s+\S+" go.mod | while read pkg version; do
    echo "Lock: $pkg at $version"
done

# For Node projects
jq '.dependencies' package.json

# For Python projects
grep "==" requirements.txt
```

## Version Tracking Structure

### Core Fields

```json
{
  "version_tracking": {
    "locked_versions": {
      "[language]": {
        "dependencies": {
          "[package_name]": {
            "version": "[exact_version]",
            "locked_by": "[phase-wave-effort]",
            "locked_date": "[ISO8601]",
            "update_status": "LOCKED",
            "update_requires": "R382_cascade"
          }
        }
      }
    }
  }
}
```

### Required Information Per Package

1. **version**: EXACT version (no ranges!)
2. **locked_by**: Which effort first added it
3. **locked_date**: When it was locked
4. **purpose**: Why this library is used
5. **update_status**: Always "LOCKED" unless in cascade
6. **update_requires**: Always "R382_cascade"

## Enforcement Workflow

### 1. Code Reviewer - During Planning

```bash
# Check what's already locked
check_locked_versions() {
    jq '.version_tracking.locked_versions' orchestrator-state-v3.json
}

# Document in implementation plan
cat >> IMPLEMENTATION-PLAN.md << 'EOF'
## Version Requirements (R381)
### Locked Dependencies (IMMUTABLE)
- cobra: v1.7.0 (locked phase1-wave1)
- go-containerregistry: v0.16.1 (locked phase1-wave1)

### New Dependencies Allowed
- prometheus/client_golang: v1.17.0 (new)
EOF
```

### 2. SW Engineer - During Implementation

```bash
# BEFORE adding any dependency
before_adding_dependency() {
    local package="$1"

    # Check if already locked
    if jq -e ".version_tracking.locked_versions.go.dependencies.\"$package\"" orchestrator-state-v3.json; then
        echo "🔴 LOCKED! Use existing version!"
        jq ".version_tracking.locked_versions.go.dependencies.\"$package\".version" orchestrator-state-v3.json
        return 1
    fi

    echo "✅ Not locked, can add new"
}

# BEFORE committing
verify_no_updates() {
    git diff HEAD -- go.mod | grep "^-" | grep -v "^---" && {
        echo "🔴 R381 VIOLATION: Version update detected!"
        exit 381
    }
}
```

### 3. Integration Agent - During Merging

```bash
# Check version consistency across branches
check_branch_versions() {
    local branches=("$@")

    for branch in "${branches[@]}"; do
        git show "$branch:go.mod" | grep -E "^\s+\S+" > "/tmp/$branch.versions"
    done

    if ! diff -q /tmp/*.versions; then
        echo "🔴 R381 VIOLATION: Branches have different versions!"
        exit 381
    fi
}
```

### 4. Orchestrator - State Tracking

```bash
# Update state when new dependencies added
update_version_tracking() {
    local package="$1"
    local version="$2"
    local effort="$3"

    jq ".version_tracking.locked_versions.go.dependencies.\"$package\" = {
        \"version\": \"$version\",
        \"locked_by\": \"$effort\",
        \"locked_date\": \"$(date -Iseconds)\",
        \"update_status\": \"LOCKED\",
        \"update_requires\": \"R382_cascade\"
    }" orchestrator-state-v3.json > temp.json && mv temp.json orchestrator-state-v3.json
}
```

## R382 Cascade Process

### When Version Update is REQUIRED

1. **Document the trigger**:
```json
{
  "cascade_id": "CASCADE-001",
  "trigger_reason": "CVE-2025-12345",
  "package": "vulnerable/lib",
  "old_version": "v1.0.0",
  "new_version": "v1.0.1",
  "approval": "user@example.com",
  "started": "2025-01-22T10:00:00Z"
}
```

2. **Find all affected branches**:
```bash
find_affected_branches() {
    local package="$1"
    local old_version="$2"

    git branch -r | while read branch; do
        if git show "$branch:go.mod" 2>/dev/null | grep -q "$package $old_version"; then
            echo "$branch"
        fi
    done
}
```

3. **Execute cascade testing**:
```bash
# Test each branch with new version
for branch in $AFFECTED_BRANCHES; do
    git checkout -b "cascade-test-$branch" "$branch"
    # Update version
    go get "$package@$new_version"
    # Run tests
    go test ./...
    # Document result
done
```

4. **Create compatibility matrix**:
```markdown
## Cascade Compatibility Matrix

| Branch | Old Version | New Version | Tests Pass | Fix Required |
|--------|------------|-------------|------------|--------------|
| phase1-wave1-effort1 | v1.0.0 | v1.0.1 | ✅ | No |
| phase1-wave1-effort2 | v1.0.0 | v1.0.1 | ❌ | Yes - API change |
```

## Common Violations and Penalties

### 🔴 -100% Instant Failures
- Updating locked version without approval
- Skipping R382 cascade when required
- Using "latest" tag instead of exact version

### ⚠️ -50% Major Violations
- Using version ranges (^, ~, >=)
- Not documenting version choices
- Incomplete cascade execution

### 📝 -20% Process Violations
- Not checking before adding dependency
- Missing version in state tracking
- Poor cascade documentation

## Recovery Procedures

### If Version Was Updated Without Authorization

1. **STOP all work immediately**
2. **Revert the update**:
```bash
git revert [commit_with_update]
```
3. **Document the violation**:
```json
{
  "violation_date": "2025-01-22T10:00:00Z",
  "package": "example/lib",
  "unauthorized_change": "v1.0.0 -> v1.1.0",
  "reverted_in": "commit-hash",
  "penalty_applied": -100
}
```
4. **If update is actually needed, trigger R382 cascade**

## Monitoring Commands

### Daily Version Audit
```bash
audit_versions() {
    echo "📊 Version Consistency Audit"
    echo "============================"

    # Check for unauthorized changes
    git log --since="1 day ago" -- go.mod package.json requirements.txt |
        grep -E "^\+" | grep -v "^+++" && {
            echo "⚠️ Version changes detected in last 24h"
        }

    # Verify all locked versions still match
    jq -r '.version_tracking.locked_versions.go.dependencies | to_entries[] |
        "\(.key) should be \(.value.version)"' orchestrator-state-v3.json |
    while read pkg should_be version; do
        actual=$(grep "$pkg" go.mod | awk '{print $2}')
        if [ "$actual" != "$version" ]; then
            echo "🔴 VIOLATION: $pkg is $actual, should be $version"
        fi
    done
}
```

### Pre-PR Version Check
```bash
pre_pr_version_check() {
    echo "🔍 Pre-PR Version Compliance Check"

    # No version changes allowed
    if git diff main...HEAD -- go.mod package.json | grep -q "^-"; then
        echo "🔴 FAIL: Version changes detected"
        echo "PR will be REJECTED per R381"
        exit 381
    fi

    echo "✅ PASS: No version changes"
}
```

## Integration with CI/CD

Add these checks to your CI pipeline:

```yaml
# .github/workflows/version-check.yml
name: R381 Version Compliance
on: [pull_request]

jobs:
  version-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Check for version changes
        run: |
          git diff origin/main...HEAD -- go.mod package.json requirements.txt |
          grep "^-" | grep -v "^---" && {
              echo "::error::R381 VIOLATION: Unauthorized version change detected"
              exit 381
          }
          echo "✅ No version changes detected"
```

---

**REMEMBER**: Version consistency enables trunk-based development. One broken version update can cascade failures through the entire project. RESPECT THE LOCKS!