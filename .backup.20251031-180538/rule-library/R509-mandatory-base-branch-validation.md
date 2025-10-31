# 🔴🔴🔴 SUPREME RULE R509: Mandatory Base Branch Validation

## Criticality: SUPREME LAW
**Violation = -100% AUTOMATIC FAILURE**

## Description
Every effort and split branch MUST validate its base branch before ANY work begins. This rule enforces cascade branching per R501 by requiring validation that each branch is based on the correct parent according to pre_planned_infrastructure. Code reviewers MUST reject any effort with incorrect base branches. Integration agents MUST verify the cascade pattern is maintained.

## 🔴🔴🔴 THE VALIDATION LAW 🔴🔴🔴

**NO WORK ON WRONG BASE BRANCHES - VALIDATION IS MANDATORY!**

### The Validation Principle:
```
BEFORE ANY WORK:
1. Check expected base from orchestrator-state-v3.json
2. Verify actual base matches expected
3. Exit immediately if wrong
4. Never attempt to fix - STOP AND REPORT
```

## 🔴 MANDATORY VALIDATION POINTS 🔴

### 1. Orchestrator Infrastructure Creation:
```bash
# BEFORE creating ANY branch
validate_base_before_creation() {
    local EFFORT_ID="$1"

    # Get expected base from pre_planned_infrastructure
    EXPECTED_BASE=$(jq -r ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".base_branch" orchestrator-state-v3.json)

    if [ "$EXPECTED_BASE" = "null" ] || [ -z "$EXPECTED_BASE" ]; then
        echo "🚨 FATAL: No base branch specified for $EFFORT_ID"
        echo "R509 VIOLATION: Missing base branch in pre_planned_infrastructure"
        exit 509
    fi

    echo "✅ Expected base branch: $EXPECTED_BASE"

    # MUST clone ONLY the base branch
    git clone -b "$EXPECTED_BASE" --single-branch "$TARGET_REPO" "$EFFORT_DIR"

    if [ $? -ne 0 ]; then
        echo "🚨 FATAL: Cannot clone base branch $EXPECTED_BASE"
        echo "R509 VIOLATION: Base branch unavailable"
        exit 509
    fi
}
```

### 2. SW Engineer Startup Validation:
```bash
# MANDATORY at SW Engineer startup
validate_my_base_branch() {
    # Get my effort ID from context
    EFFORT_ID=$(identify_my_effort)

    # Get expected base from orchestrator-state-v3.json
    EXPECTED_BASE=$(jq -r ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".base_branch" \
                    /workspaces/software-factory-2.0/orchestrator-state-v3.json)

    # Get actual base (parent commit of my branch point)
    ACTUAL_BASE=$(git merge-base HEAD origin/main)
    EXPECTED_BASE_COMMIT=$(git rev-parse "origin/$EXPECTED_BASE")

    if [ "$ACTUAL_BASE" != "$EXPECTED_BASE_COMMIT" ]; then
        echo "🚨🚨🚨 R509 VIOLATION: WRONG BASE BRANCH!"
        echo "Expected base: $EXPECTED_BASE"
        echo "Actual base commit: $ACTUAL_BASE"
        echo "Expected base commit: $EXPECTED_BASE_COMMIT"
        echo ""
        echo "MANDATORY ACTION: STOP IMMEDIATELY!"
        echo "DO NOT ATTEMPT TO FIX!"
        echo "Report to orchestrator for infrastructure rebuild"
        exit 509
    fi

    echo "✅ Base branch validated: correctly based on $EXPECTED_BASE"
}
```

### 3. Code Reviewer Validation:
```bash
# Code Reviewer MUST validate during EFFORT_PLAN_CREATION
validate_effort_infrastructure() {
    local EFFORT_DIR="$1"
    local EFFORT_ID="$2"

    cd "$EFFORT_DIR"

    # Check branch name matches expected
    EXPECTED_BRANCH=$(jq -r ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".branch_name" \
                      /workspaces/software-factory-2.0/orchestrator-state-v3.json)
    ACTUAL_BRANCH=$(git branch --show-current)

    if [ "$ACTUAL_BRANCH" != "$EXPECTED_BRANCH" ]; then
        echo "❌ R509 VIOLATION: Wrong branch name!"
        echo "Expected: $EXPECTED_BRANCH"
        echo "Actual: $ACTUAL_BRANCH"
        mark_effort_failed "R509: Wrong branch name"
        return 509
    fi

    # Check base branch
    EXPECTED_BASE=$(jq -r ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".base_branch" \
                    /workspaces/software-factory-2.0/orchestrator-state-v3.json)

    # Verify cascade pattern
    if ! verify_cascade_pattern "$EFFORT_ID"; then
        echo "❌ R509 VIOLATION: Cascade pattern broken!"
        mark_effort_failed "R509: Not following cascade pattern"
        return 509
    fi

    echo "✅ Infrastructure validation passed"
}
```

### 4. Integration Agent Validation:
```bash
# Integration MUST verify cascade before ANY merge
validate_cascade_integrity() {
    echo "🔍 R509: Validating cascade integrity..."

    # Check every effort follows cascade
    jq -r '.final_merge_plan.merge_sequence[] |
           "\(.order):\(.branch):\(.base_branch)"' \
    orchestrator-state-v3.json | while IFS=: read -r order branch base; do

        cd "$EFFORTS_DIR/$branch"

        # Verify branch is actually based on its declared base
        if ! git merge-base --is-ancestor "origin/$base" HEAD; then
            echo "🚨 R509 VIOLATION: $branch not based on $base!"
            echo "CASCADE BROKEN - CANNOT PROCEED"
            exit 509
        fi

        echo "✅ $branch correctly based on $base"
    done

    echo "✅ Full cascade validated"
}
```

## 🔴 VALIDATION UTILITY SCRIPT 🔴

### /utilities/validate-base-branches.sh:
```bash
#!/bin/bash
# R509 Enforcement: Validate all base branches

set -e

echo "🔍 R509: Mandatory Base Branch Validation"
echo "=========================================="

# Load orchestrator state
STATE_FILE="/workspaces/software-factory-2.0/orchestrator-state-v3.json"

if [ ! -f "$STATE_FILE" ]; then
    echo "🚨 FATAL: No orchestrator-state-v3.json found!"
    exit 509
fi

# Validate each effort
jq -r '.pre_planned_infrastructure.efforts | to_entries[] |
       "\(.key):\(.value.branch_name):\(.value.base_branch):\(.value.full_path)"' \
"$STATE_FILE" | while IFS=: read -r effort_id branch base path; do

    echo ""
    echo "Validating: $effort_id"
    echo "  Branch: $branch"
    echo "  Expected base: $base"
    echo "  Path: $path"

    if [ ! -d "$path/.git" ]; then
        echo "  ⏭️ Infrastructure not yet created"
        continue
    fi

    cd "$path"

    # Check current branch
    CURRENT=$(git branch --show-current)
    if [ "$CURRENT" != "$branch" ]; then
        echo "  ❌ WRONG BRANCH! Current: $CURRENT"
        EXIT_CODE=509
        continue
    fi

    # Check base branch
    if [ "$base" = "main" ]; then
        # First effort should be from main
        BASE_COMMIT=$(git merge-base HEAD origin/main)
        MAIN_COMMIT=$(git rev-parse origin/main)

        if [ "$BASE_COMMIT" != "$MAIN_COMMIT" ]; then
            echo "  ❌ NOT BASED ON MAIN!"
            EXIT_CODE=509
        else
            echo "  ✅ Correctly based on main"
        fi
    else
        # Should be based on previous effort
        if git show-ref --verify --quiet "refs/remotes/origin/$base"; then
            if git merge-base --is-ancestor "origin/$base" HEAD; then
                echo "  ✅ Correctly based on $base"
            else
                echo "  ❌ NOT BASED ON $base!"
                EXIT_CODE=509
            fi
        else
            echo "  ❌ Base branch $base doesn't exist!"
            EXIT_CODE=509
        fi
    fi
done

if [ "${EXIT_CODE:-0}" -ne 0 ]; then
    echo ""
    echo "🚨🚨🚨 R509 VIOLATIONS DETECTED!"
    echo "CASCADE PATTERN BROKEN!"
    echo "IMMEDIATE ACTION REQUIRED!"
    exit 509
fi

echo ""
echo "✅ All base branches validated successfully"
echo "CASCADE PATTERN INTACT"
```

## 🚨 COMMON VIOLATIONS (AUTOMATIC FAILURE) 🚨

### ❌ VIOLATION 1: Creating from Wrong Base
```bash
# WRONG - Not checking base branch
git clone "$REPO" "$DIR"
git checkout -b new-branch  # From whatever default branch!

# RIGHT - Clone specific base branch
BASE=$(jq -r '.pre_planned_infrastructure.efforts.X.base_branch' state.json)
git clone -b "$BASE" --single-branch "$REPO" "$DIR"
git checkout -b new-branch  # Now correctly based
```

### ❌ VIOLATION 2: Skipping Validation
```bash
# WRONG - Starting work without validation
cd effort-dir
start_implementation

# RIGHT - Validate first
cd effort-dir
validate_my_base_branch  # Exit 509 if wrong
start_implementation
```

### ❌ VIOLATION 3: Attempting to Fix Wrong Base
```bash
# WRONG - Trying to rebase or fix
if [ "$BASE" != "$EXPECTED" ]; then
    git rebase "$EXPECTED"  # NO! NEVER DO THIS!
fi

# RIGHT - Stop immediately
if [ "$BASE" != "$EXPECTED" ]; then
    echo "R509 VIOLATION: Wrong base!"
    exit 509  # STOP! Don't fix!
fi
```

## 🔴 ENFORCEMENT MECHANISMS 🔴

### 1. Pre-Commit Hook:
```bash
# Installed in every effort directory
#!/bin/bash
# R509 Pre-commit validation

EFFORT_ID=$(basename "$PWD")
EXPECTED_BASE=$(jq -r ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".base_branch" \
                /workspaces/software-factory-2.0/orchestrator-state-v3.json)

if ! git merge-base --is-ancestor "origin/$EXPECTED_BASE" HEAD; then
    echo "🚨 R509 VIOLATION: Committing to wrong base branch!"
    echo "Expected base: $EXPECTED_BASE"
    echo "COMMIT BLOCKED"
    exit 1
fi
```

### 2. CI/CD Pipeline:
```yaml
validate-base-branch:
  stage: validate
  script:
    - bash utilities/validate-base-branches.sh
  rules:
    - if: '$CI_COMMIT_BRANCH =~ /^effort-/'
      when: always
```

### 3. Code Review Gate:
```markdown
## R509 Validation Checklist
- [ ] Branch based on correct parent per cascade
- [ ] Base branch validated at startup
- [ ] No rebasing or base changes detected
- [ ] Cascade pattern maintained
```

## 🔴 GRADING IMPACT 🔴

- **Wrong base branch**: -100% (CASCADE VIOLATION)
- **Missing validation**: -50% (Process failure)
- **Attempting to fix base**: -75% (Protocol violation)
- **Skipping validation**: -100% (Safety violation)
- **Code reviewer missing it**: -50% (Review failure)

## 🔴 WHY THIS MATTERS 🔴

### Without Base Branch Validation:
- **Cascade Breaks**: Entire progressive model fails
- **Merge Conflicts**: Branches don't build on each other
- **Lost Work**: Changes overwritten in merges
- **Integration Hell**: Can't merge sequentially

### With Base Branch Validation:
- **Cascade Protected**: Progressive model preserved
- **Clean Merges**: Each branch includes previous
- **Work Preserved**: All changes accumulate
- **Smooth Integration**: Sequential merging works

## 🔴 THE FINAL TRUTH 🔴

**VALIDATION IS NOT OPTIONAL - IT'S SURVIVAL!**

- Validate BEFORE creating branches
- Validate AT agent startup
- Validate DURING code review
- Validate BEFORE integration
- NEVER work on wrong base
- NEVER try to fix - STOP AND REPORT!

**R509 ensures the cascade pattern that makes Software Factory 2.0 work!**