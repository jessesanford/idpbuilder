# 🔴🔴🔴 SUPREME RULE R361: Integration Conflict Resolution Only Protocol

## Criticality: SUPREME LAW
**Violation = -100% AUTOMATIC FAILURE**

## Description
This rule ABSOLUTELY PROHIBITS adding new code, packages, adapters, wrappers, or ANY new files during integration. Integration is for CONFLICT RESOLUTION ONLY - selecting which version of existing code to keep when branches conflict. Any code additions, no matter how small or seemingly necessary, MUST be done in effort branches BEFORE integration.

## 🔴🔴🔴 THE ABSOLUTE LAW 🔴🔴🔴

**INTEGRATION = CONFLICT RESOLUTION ONLY. NO NEW CODE. EVER.**

### The Problem This Solves
```
❌ WRONG (WHAT HAPPENED IN YOUR VIOLATION):
Effort1 (using pkg/registry directly) ──┐
                                        ├──> Integration adds pkg/gitea/ adapter ──> BUG IN NEW CODE
Effort2 (using pkg/registry directly) ──┘         ↑
                                              VIOLATION!
                                              New package created during integration

✅ CORRECT:
Effort1 (with gitea adapter if needed) ──┐
                                         ├──> Integration (merges only) ──> NO NEW CODE
Effort2 (compatible with effort1) ──────┘
```

## Core Requirements

### 1. ABSOLUTE PROHIBITION ON NEW FILES
```bash
# MANDATORY CHECK: No new files during integration
validate_no_new_files() {
    local BASE_BRANCH=$1
    local INTEGRATION_BRANCH=$2

    # Get list of files that exist in integration but not in any source
    NEW_FILES=$(git diff --name-status $BASE_BRANCH..$INTEGRATION_BRANCH | grep "^A" | cut -f2)

    for file in $NEW_FILES; do
        # Check if file exists in ANY source branch
        FILE_EXISTS_IN_SOURCE=false
        for source_branch in $(get_source_branches); do
            if git ls-tree $source_branch --name-only | grep -q "^$file$"; then
                FILE_EXISTS_IN_SOURCE=true
                break
            fi
        done

        if [ "$FILE_EXISTS_IN_SOURCE" = false ]; then
            echo "🔴🔴🔴 R361 VIOLATION: New file created during integration!"
            echo "File: $file"
            echo "This file doesn't exist in ANY source branch!"
            echo "ALL files must originate from effort branches!"
            exit 361
        fi
    done
}
```

### 2. MAXIMUM 50 LINES CHANGE LIMIT
```bash
# ENFORCEMENT: Integration changes must be minimal
validate_integration_changes() {
    local BASE=$1
    local INTEGRATION=$2

    # Count non-merge commit changes
    NON_MERGE_CHANGES=$(git diff --shortstat $BASE..$INTEGRATION --no-merges)
    LINES_ADDED=$(echo $NON_MERGE_CHANGES | grep -oP '\d+(?= insertion)' || echo 0)
    LINES_REMOVED=$(echo $NON_MERGE_CHANGES | grep -oP '\d+(?= deletion)' || echo 0)
    TOTAL_CHANGES=$((LINES_ADDED + LINES_REMOVED))

    if [ $TOTAL_CHANGES -gt 50 ]; then
        echo "🔴🔴🔴 R361 VIOLATION: Too many changes during integration!"
        echo "Changes detected: $TOTAL_CHANGES lines"
        echo "Maximum allowed: 50 lines (for conflict markers only)"
        echo ""
        echo "This indicates NEW CODE was added, not just conflicts resolved!"
        exit 361
    fi
}
```

### 3. PROHIBITED ACTIONS DURING INTEGRATION

#### ABSOLUTELY FORBIDDEN - AUTOMATIC FAILURE:
```bash
# ❌ Creating new packages
mkdir pkg/gitea                              # R361 VIOLATION!
touch pkg/gitea/adapter.go                   # R361 VIOLATION!

# ❌ Adding adapter/wrapper layers
cat > pkg/adapters/registry_wrapper.go       # R361 VIOLATION!

# ❌ Creating "glue code"
vim pkg/integration/compatibility.go         # R361 VIOLATION!

# ❌ Adding helper functions
echo "func Helper() {}" >> utils.go          # R361 VIOLATION!

# ❌ Creating configuration files
cat > config/integration-fixes.yaml          # R361 VIOLATION!

# ❌ Adding test helpers
touch test/integration_helpers_test.go       # R361 VIOLATION!
```

#### ONLY ALLOWED ACTIONS:
```bash
# ✅ Choosing between conflicting versions
git checkout --ours file.go     # Keep our version
git checkout --theirs file.go   # Keep their version

# ✅ Manually selecting specific lines (but file must already exist)
vim existing_file.go
# Keep lines from version A
# Keep lines from version B
# BUT DO NOT ADD NEW LINES!

# ✅ Fixing import statements (minimal, <10 lines)
# ONLY if both versions already had imports
```

### 4. INTEGRATION AGENT RESTRICTIONS

The Integration Agent MUST acknowledge:
```bash
echo "🔴🔴🔴 R361 ENFORCEMENT ACTIVE 🔴🔴🔴"
echo "I am PROHIBITED from:"
echo "  ❌ Creating ANY new files"
echo "  ❌ Creating ANY new directories"
echo "  ❌ Adding ANY new packages"
echo "  ❌ Writing ANY adapter code"
echo "  ❌ Creating ANY wrapper layers"
echo "  ❌ Adding more than 50 lines total"
echo ""
echo "I can ONLY:"
echo "  ✅ Resolve merge conflicts"
echo "  ✅ Choose between existing versions"
echo "  ✅ Fix import statements (<10 lines)"
```

### 5. WHAT TO DO WHEN INTEGRATION NEEDS NEW CODE

When you discover that integration requires new code (adapters, wrappers, etc):

```bash
handle_integration_needs_new_code() {
    echo "🔴 STOP: Integration requires new code"
    echo "🔴 R361 PROHIBITS creating code during integration"

    # 1. Document what's needed
    cat > INTEGRATION-BLOCKED-NEEDS-CODE.md << 'EOF'
# Integration Blocked - New Code Required

## Issue
Integration cannot proceed without new adapter/wrapper code.

## Required Code
- Package: pkg/gitea
- Purpose: Adapter between effort1 and effort2 interfaces
- Estimated Size: ~200 lines

## Action Required
1. STOP integration immediately
2. Create new effort branch for adapter
3. Implement adapter in effort branch
4. Test adapter independently
5. THEN retry integration with adapter branch included
EOF

    # 2. Transition to state that creates new effort
    transition_to_state "SPAWN_ENGINEER_FOR_ADAPTER"

    # 3. NEVER proceed with integration
    exit 361
}
```

## Detection and Enforcement

### Pre-Integration Validation
```bash
# Run BEFORE starting any integration
pre_integration_check() {
    echo "🔍 R361 Pre-Integration Validation"

    # Verify all needed code exists in effort branches
    for effort in $(list_efforts); do
        echo "Checking $effort has all required code..."

        # Check for interface compatibility
        if ! verify_interfaces_compatible $effort; then
            echo "❌ $effort needs adapter code BEFORE integration"
            echo "Create adapter in effort branch first!"
            exit 361
        fi
    done

    echo "✅ All efforts have necessary code"
}
```

### Post-Integration Audit
```bash
# Run AFTER integration completes
post_integration_audit() {
    echo "📊 R361 Post-Integration Audit"

    # Check for new packages
    NEW_PACKAGES=$(find . -type d -name "pkg/*" -newer integration_start.timestamp)
    if [ -n "$NEW_PACKAGES" ]; then
        echo "🔴 R361 VIOLATION: New packages created!"
        echo "$NEW_PACKAGES"
        exit 361
    fi

    # Check total changes
    STATS=$(git diff --shortstat main..HEAD)
    echo "Integration changes: $STATS"

    # Verify no new files
    NEW_FILES=$(git diff --name-status main..HEAD | grep "^A")
    if [ -n "$NEW_FILES" ]; then
        echo "🔴 R361 VIOLATION: New files added!"
        echo "$NEW_FILES"
        exit 361
    fi
}
```

## Common Violations and Corrections

### ❌ VIOLATION: Creating Adapter During Integration
```bash
# WRONG - What happened in the reported violation:
cd integration-workspace
git merge effort1  # Uses pkg/registry
git merge effort2  # Also uses pkg/registry
# "Oh, these need an adapter to work together"
mkdir pkg/gitea
vim pkg/gitea/adapter.go  # R361 VIOLATION!
# Created 300 lines of new adapter code
git add pkg/gitea
git commit -m "Add gitea adapter for integration"  # VIOLATION!
```

### ✅ CORRECTION: Create Adapter in Effort Branch First
```bash
# RIGHT - How it should be done:
cd integration-workspace
git merge effort1
git merge effort2
# Discover incompatibility
echo "Integration blocked - needs adapter"

# STOP integration
cd /efforts/phase1/wave1
mkdir effort3-gitea-adapter
cd effort3-gitea-adapter
git checkout -b phase1-wave1-gitea-adapter

# Create adapter in NEW EFFORT
mkdir pkg/gitea
vim pkg/gitea/adapter.go
# Implement, test, commit, push

# THEN retry integration with adapter branch
cd integration-workspace
git merge effort3-gitea-adapter  # Now adapter exists
git merge effort1  # Can use adapter
git merge effort2  # Can use adapter
```

## Why This Rule Is Critical

### The Cascade Effect of Integration Code
1. Code added during integration has NO effort branch source
2. This code was NEVER independently tested
3. It exists ONLY in integration branch
4. When integration fails, the code is lost
5. No audit trail of who wrote it or why
6. Violates trunk-based development principles
7. Makes rollback impossible

### Proper Code Flow
1. ALL code originates in effort branches
2. Every line is traceable to an effort
3. Every package was independently tested
4. Integration ONLY combines tested code
5. Clean audit trail maintained
6. Easy rollback to any effort state

## Grading Impact

### AUTOMATIC FAILURE (-100%)
- ANY new file created during integration
- ANY new package/directory created
- More than 50 lines changed (excluding merge commits)
- Creating adapters/wrappers during integration
- Adding "glue code" or "compatibility layers"

### MAJOR VIOLATIONS (-50%)
- Attempting to hide new code as "conflict resolution"
- Not stopping when adapter needed
- Proceeding without required code in efforts

### COMPLIANCE BONUS (+25%)
- Zero new files in integration
- Less than 20 lines changed total
- Clean merge commits only
- Proper escalation when code needed

## Relationship to Other Rules

### Strengthens R321
- R321: No bug fixes during integration
- R361: No new code AT ALL during integration

### Strengthens R266
- R266: Document bugs, don't fix
- R361: Document missing code, don't create

### Enforces R262
- R262: Don't modify original branches
- R361: Don't create new code branches

## Quick Reference

### Check Integration Compliance
```bash
# Should return EMPTY (no new files)
git diff --name-status main..integration | grep "^A"

# Should be <50 lines
git diff --shortstat main..integration --no-merges

# Should return EMPTY (no new packages)
find . -type d -path "*/pkg/*" -newer integration_start.timestamp
```

### Maximum Allowed Changes
- Conflict resolution: Selecting between versions
- Import fixes: <10 lines
- Comment updates: <5 lines
- Total non-merge changes: <50 lines
- New files: 0
- New packages: 0
- New functions: 0

## Remember

**"Integration combines, never creates"**
**"Every line of code has an effort home"**
**"If it doesn't exist in efforts, it doesn't exist"**
**"Adapters are efforts, not integration glue"**

The golden rule: If you're writing code during integration, you're doing it wrong. Stop immediately and create an effort for that code.