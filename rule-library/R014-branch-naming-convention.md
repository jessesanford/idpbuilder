# Rule R014: Branch Naming Convention

## Rule Statement
ALL branches created by the Software Factory system MUST follow standardized naming conventions that include the project prefix when configured. This applies to effort branches, integration branches (wave and phase), split branches, and fix branches.

## Criticality Level
**MANDATORY** - Incorrect branch naming causes CI/CD failures and PR routing issues

## Enforcement Mechanism
- **Technical**: Use branch-naming-helpers.sh functions exclusively
- **Validation**: Verify branch names include project prefix before creation
- **Grading**: -20% for incorrect branch naming

## Core Principle

Branch names must be consistent, predictable, and include the project prefix to:
1. Enable automated CI/CD pipeline triggers
2. Support PR routing and review assignments
3. Maintain clear project ownership in multi-project repositories
4. Enable branch protection rules

## Detailed Requirements

### MANDATORY: Use Branch Naming Helper Functions

**NEVER hardcode branch names. ALWAYS use the helper functions:**

```bash
# Source the helpers (from SF instance directory)
source "$SF_INSTANCE_DIR/utilities/branch-naming-helpers.sh"

# Get project prefix from config
# Use text_editor tool with view command to read target-repo-config.yaml:
# Find the branch_naming.project_prefix field
PROJECT_PREFIX="<value from branch_naming.project_prefix>"
```

### Branch Types and Functions

#### 1. Effort Branches
```bash
# Function: get_effort_branch_name
EFFORT_BRANCH=$(get_effort_branch_name "$PHASE" "$WAVE" "$EFFORT_NAME")

# Example with prefix: tmc-workspace/phase3/wave2/effort1-api-gateway
# Example without: phase3/wave2/effort1-api-gateway
```

#### 2. Wave Integration Branches
```bash
# Function: get_wave_integration_branch_name
WAVE_INTEGRATION=$(get_wave_integration_branch_name "$PHASE" "$WAVE")

# Example with prefix: tmc-workspace/phase3/wave2-integration
# Example without: phase3/wave2-integration
```

#### 3. Phase Integration Branches
```bash
# Function: get_phase_integration_branch_name
PHASE_INTEGRATION=$(get_phase_integration_branch_name "$PHASE")

# Example with prefix: tmc-workspace/phase3-integration
# Example without: phase3-integration
```

#### 4. Split Branches
```bash
# Function: get_split_branch_name
SPLIT_BRANCH=$(get_split_branch_name "$ORIGINAL_BRANCH" "$SPLIT_NUMBER")

# Example: tmc-workspace/phase3/wave2/effort3-controller--split-001
# The original branch already includes prefix if configured
```

### Project Prefix Configuration

The project prefix comes from `target-repo-config.yaml`:

```yaml
branch_naming:
  project_prefix: "tmc-workspace"  # Or empty string for no prefix
  effort_format: "{prefix}phase{phase}/wave{wave}/{effort_name}"
  integration_format: "{prefix}phase{phase}/wave{wave}-integration"
  phase_integration_format: "{prefix}phase{phase}-integration"
```

### Critical Implementation Rules

#### ❌ FORBIDDEN - Hardcoding Branch Names
```bash
# WRONG - Never hardcode branch names!
BRANCH="phase3/wave2/effort1"  # Missing prefix!
BRANCH="phase3-integration-$(date +%Y%m%d)"  # Wrong format!

# WRONG - Creating branches without helpers
git checkout -b "phase${PHASE}-integration"  # No prefix!
```

#### ✅ REQUIRED - Always Use Helpers
```bash
# RIGHT - Always use helper functions
source utilities/branch-naming-helpers.sh
BRANCH=$(get_effort_branch_name "$PHASE" "$WAVE" "$EFFORT_NAME")
git checkout -b "$BRANCH"

# RIGHT - Phase integration with helpers
PHASE_BRANCH=$(get_phase_integration_branch_name "$PHASE")
git checkout -b "$PHASE_BRANCH"
```

### Validation Requirements

Before creating any branch, validate it includes the prefix:

```bash
validate_branch_name() {
    local branch_name="$1"
    # Use text_editor tool with view command to read target-repo-config.yaml:
    # Find the branch_naming.project_prefix field
    local project_prefix="<value from branch_naming.project_prefix>"
    
    if [ -n "$project_prefix" ] && [ "$project_prefix" != "null" ]; then
        if [[ ! "$branch_name" =~ ^"$project_prefix"/ ]]; then
            echo "❌ ERROR: Branch missing project prefix: $project_prefix"
            echo "Branch name: $branch_name"
            exit 1
        fi
    fi
    
    echo "✅ Branch name validated: $branch_name"
}
```

## Integration with Other Rules

### R250 - Integration Isolation
Integration branches MUST follow R014 naming conventions when created in integration workspaces.

### R208 - Spawn Directory Protocol  
Effort branches created during spawn MUST use get_effort_branch_name function.

### R269-R270 - Three-Agent Integration
Integration Agent must verify branch names follow R014 before executing merges.

## Common Violations to Avoid

### 1. Phase Integration Without Prefix
```bash
# ❌ WRONG
BRANCH="phase3-integration-20250827"

# ✅ RIGHT
source utilities/branch-naming-helpers.sh
BRANCH=$(get_phase_integration_branch_name "3")
```

### 2. Manual Timestamp Suffixes
```bash
# ❌ WRONG - Don't add timestamps to standard branches
BRANCH="phase3/wave2-integration-$(date +%Y%m%d)"

# ✅ RIGHT - Use standard format
BRANCH=$(get_wave_integration_branch_name "3" "2")
```

### 3. Forgetting to Source Helpers
```bash
# ❌ WRONG - Helpers not sourced
BRANCH=$(get_effort_branch_name "3" "2" "api")  # Will fail!

# ✅ RIGHT - Always source first
source "$SF_INSTANCE_DIR/utilities/branch-naming-helpers.sh"
BRANCH=$(get_effort_branch_name "3" "2" "api")
```

## Grading Impact

- **Missing project prefix**: -20% (breaks CI/CD)
- **Wrong format**: -15% (routing failures)
- **Hardcoded names**: -10% (maintenance issues)
- **No validation**: -5% (potential errors)

## Implementation Priority

**IMMEDIATE** - All agents creating branches must:
1. Source branch-naming-helpers.sh
2. Use appropriate helper function
3. Never hardcode branch names
4. Validate before creation

## Summary

R014 ensures consistent branch naming across the entire Software Factory system. Every branch - whether effort, integration, or split - must use the helper functions that automatically include the project prefix when configured. This is not optional; it's mandatory for system functionality.

**Remember**: If you're typing a branch name manually instead of calling a helper function, you're violating R014!