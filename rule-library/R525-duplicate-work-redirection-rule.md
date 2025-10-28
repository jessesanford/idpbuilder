# 🚨 RULE R525: Duplicate Work Redirection Rule

## Criticality: CRITICAL
**Violation Penalty: -60% for working on duplicates, -30% for missing redirection**

## Description
This rule ensures that agents NEVER work on duplicate bugs directly. All work must be redirected to the canonical bug. This prevents wasted effort, ensures fixes are comprehensive, and maintains bug registry integrity.

## The Problem This Solves

### Without R525 (WASTEFUL):
```
Agent assigned BUG-001 (duplicate of BUG-007)
Agent spends 2 hours fixing BUG-001
BUG-007 (canonical) still OPEN - needs separate fix
Result: Wasted effort, incomplete fix
```

### With R525 (EFFICIENT):
```
Agent assigned BUG-001 (duplicate of BUG-007)
R525 detects: BUG-001 is duplicate
R525 redirects: Work on BUG-007 instead
Agent fixes BUG-007 (canonical)
R524 propagates: BUG-001 auto-updated to FIXED
Result: Efficient, complete fix
```

## 🚦 MANDATORY PRE-WORK CHECK

### Before Starting ANY Bug Fix
```bash
# EVERY agent MUST run this before working on a bug
check_bug_before_work() {
    local bug_id="$1"

    echo "🔍 R525: Pre-work bug validation for $bug_id"

    # Check 1: Does bug exist?
    BUG_EXISTS=$(jq -r --arg bug "$bug_id" '
        .bug_registry[] |
        select(.bug_id == $bug) |
        .bug_id
    ' orchestrator-state-v3.json)

    if [[ -z "$BUG_EXISTS" ]]; then
        echo "❌ R525: Bug $bug_id does not exist!"
        return 1
    fi

    # Check 2: Is this a duplicate?
    IS_DUPLICATE=$(jq -r --arg bug "$bug_id" '
        .bug_registry[] |
        select(.bug_id == $bug) |
        .is_duplicate // false
    ' orchestrator-state-v3.json)

    if [[ "$IS_DUPLICATE" == "true" ]]; then
        # Get canonical bug
        CANONICAL=$(jq -r --arg bug "$bug_id" '
            .bug_registry[] |
            select(.bug_id == $bug) |
            .duplicate_of
        ' orchestrator-state-v3.json)

        echo "⚠️ R525: $bug_id is a DUPLICATE of $CANONICAL"
        echo "   🔄 REDIRECTING work to canonical bug $CANONICAL"
        echo "   ❌ DO NOT work on $bug_id - status will auto-propagate via R524"

        # Return canonical bug ID
        echo "$CANONICAL"
        return 2  # Special code: redirect needed
    fi

    # Check 3: Is bug already fixed?
    FIX_STATUS=$(jq -r --arg bug "$bug_id" '
        .bug_registry[] |
        select(.bug_id == $bug) |
        .fix_status
    ' orchestrator-state-v3.json)

    if [[ "$FIX_STATUS" == "FIXED" ]] || [[ "$FIX_STATUS" == "FIXED_AS_DUPLICATE" ]]; then
        echo "✅ R525: Bug $bug_id already FIXED (status: $FIX_STATUS)"
        echo "   No work needed"
        return 3  # Already fixed
    fi

    echo "✅ R525: $bug_id is canonical and OPEN - safe to work on"
    return 0  # Safe to proceed
}
```

### Redirect Work Flow
```bash
# When duplicate detected, redirect to canonical
redirect_to_canonical() {
    local duplicate_bug="$1"

    echo "🔄 R525: Redirecting from duplicate to canonical"

    # Get canonical bug
    CANONICAL=$(jq -r --arg bug "$duplicate_bug" '
        .bug_registry[] |
        select(.bug_id == $bug) |
        .duplicate_of
    ' orchestrator-state-v3.json)

    if [[ -z "$CANONICAL" ]]; then
        echo "❌ R525 VIOLATION: Bug marked as duplicate but has no duplicate_of!"
        return 1
    fi

    # Verify canonical is not itself a duplicate (should never happen per R523)
    CANONICAL_IS_DUP=$(jq -r --arg bug "$CANONICAL" '
        .bug_registry[] |
        select(.bug_id == $bug) |
        .is_duplicate // false
    ' orchestrator-state-v3.json)

    if [[ "$CANONICAL_IS_DUP" == "true" ]]; then
        echo "❌ R523 VIOLATION: Canonical $CANONICAL is itself a duplicate!"
        echo "   This is a transitive duplicate - R523 should prevent this"
        return 1
    fi

    echo "✅ R525: Canonical bug is $CANONICAL"
    echo "   All work should target $CANONICAL"
    echo "   Fix status will propagate to $duplicate_bug via R524"

    # Return canonical for agent to use
    echo "$CANONICAL"
}
```

## 🎯 AGENT-SPECIFIC IMPLEMENTATIONS

### SW Engineer (Before Implementation)
```bash
# SW Engineer startup for bug fix effort
sw_engineer_bug_fix_startup() {
    local assigned_bug="$1"

    echo "🚀 SW Engineer starting bug fix for $assigned_bug"

    # MANDATORY: Check bug before work
    check_bug_before_work "$assigned_bug"
    CHECK_RESULT=$?

    case $CHECK_RESULT in
        0)
            # Safe to proceed
            echo "✅ Working on $assigned_bug (canonical bug)"
            WORK_BUG="$assigned_bug"
            ;;
        2)
            # Redirect needed
            CANONICAL=$(redirect_to_canonical "$assigned_bug")
            echo "🔄 Redirected from $assigned_bug to $CANONICAL"
            WORK_BUG="$CANONICAL"

            # Update effort metadata to reflect redirection
            jq --arg effort "$EFFORT_ID" \
               --arg original "$assigned_bug" \
               --arg canonical "$CANONICAL" \
               '
               (.efforts[$effort]) |= {
                   originally_assigned_bug: $original,
                   actual_bug_fixed: $canonical,
                   redirected_per_r525: true,
                   .
               }
               ' orchestrator-state-v3.json > /tmp/state.json

            mv /tmp/state.json orchestrator-state-v3.json
            ;;
        3)
            # Already fixed
            echo "✅ Bug $assigned_bug already fixed - no work needed"
            return 0
            ;;
        *)
            # Error
            echo "❌ Cannot proceed with $assigned_bug"
            return 1
            ;;
    esac

    # Proceed with WORK_BUG (either original or canonical)
    implement_bug_fix "$WORK_BUG"
}
```

### Orchestrator (Assigning Bug Fixes)
```bash
# When orchestrator assigns bug to SW engineer, validate first
orchestrator_assign_bug_fix() {
    local bug_id="$1"
    local effort_id="$2"

    echo "📋 Orchestrator: Assigning $bug_id to $effort_id"

    # Pre-validate bug
    check_bug_before_work "$bug_id"
    CHECK_RESULT=$?

    if [[ $CHECK_RESULT -eq 2 ]]; then
        # Bug is duplicate - assign canonical instead
        CANONICAL=$(redirect_to_canonical "$bug_id")

        echo "⚠️ R525: $bug_id is duplicate - assigning canonical $CANONICAL instead"

        # Update assignment to canonical
        bug_id="$CANONICAL"
    elif [[ $CHECK_RESULT -eq 3 ]]; then
        # Already fixed - don't assign
        echo "✅ Bug $bug_id already fixed - skipping assignment"
        return 0
    fi

    # Create fix effort for canonical bug
    create_fix_effort "$bug_id" "$effort_id"
}
```

### Code Reviewer (Review Bug Fixes)
```bash
# When reviewing bug fix, verify work was on canonical
code_reviewer_bug_fix_validation() {
    local effort_id="$1"

    echo "🔍 Code Reviewer: Validating bug fix effort $effort_id"

    # Get bug that was fixed
    BUG_FIXED=$(jq -r --arg effort "$effort_id" '
        .efforts[$effort].actual_bug_fixed // .efforts[$effort].assigned_bug
    ' orchestrator-state-v3.json)

    # Verify it's not a duplicate
    IS_DUPLICATE=$(jq -r --arg bug "$BUG_FIXED" '
        .bug_registry[] |
        select(.bug_id == $bug) |
        .is_duplicate // false
    ' orchestrator-state-v3.json)

    if [[ "$IS_DUPLICATE" == "true" ]]; then
        echo "❌ R525 VIOLATION: Fixed a duplicate bug directly!"
        echo "   Bug: $BUG_FIXED"
        echo "   This should have been redirected to canonical per R525"
        return 1
    fi

    echo "✅ R525: Bug fix worked on canonical bug"
}
```

## 🚨 ENFORCEMENT AND VALIDATION

### Detect Work on Duplicates
```bash
# Scan for any efforts that worked on duplicate bugs
detect_duplicate_work_violations() {
    echo "📊 R525: Scanning for duplicate work violations"

    # Find efforts that fixed duplicate bugs
    VIOLATIONS=$(jq -r '
        .efforts | to_entries[] |
        select(.value.actual_bug_fixed != null) as $effort |
        .bug_registry[] |
        select(.bug_id == $effort.value.actual_bug_fixed and .is_duplicate == true) |
        "VIOLATION: Effort \($effort.key) fixed duplicate bug \(.bug_id) (should be \(.duplicate_of))"
    ' orchestrator-state-v3.json)

    if [[ -n "$VIOLATIONS" ]]; then
        echo "❌ R525 VIOLATIONS DETECTED:"
        echo "$VIOLATIONS"
        return 1
    fi

    echo "✅ R525: No duplicate work violations found"
}
```

### Validate Redirection Records
```bash
# Verify all redirections were recorded properly
validate_redirection_records() {
    echo "📊 R525: Validating redirection records"

    # Find efforts that were redirected
    REDIRECTIONS=$(jq -r '
        .efforts | to_entries[] |
        select(.value.redirected_per_r525 == true) |
        "\(.key): \(.value.originally_assigned_bug) → \(.value.actual_bug_fixed)"
    ' orchestrator-state-v3.json)

    if [[ -n "$REDIRECTIONS" ]]; then
        echo "✅ R525: Found redirections (properly handled):"
        echo "$REDIRECTIONS"

        # Verify each redirection was valid
        while IFS=: read -r effort_id redirection; do
            ORIGINAL=$(echo "$redirection" | cut -d→ -f1 | xargs)
            CANONICAL=$(echo "$redirection" | cut -d→ -f2 | xargs)

            # Verify original is duplicate of canonical
            DUPLICATE_OF=$(jq -r --arg bug "$ORIGINAL" '
                .bug_registry[] |
                select(.bug_id == $bug) |
                .duplicate_of
            ' orchestrator-state-v3.json)

            if [[ "$DUPLICATE_OF" != "$CANONICAL" ]]; then
                echo "❌ R525: Invalid redirection in $effort_id"
                echo "   $ORIGINAL is not duplicate of $CANONICAL"
                return 1
            fi
        done <<< "$REDIRECTIONS"

        echo "✅ All redirections validated"
    else
        echo "ℹ️ No redirections found (no duplicate bugs encountered)"
    fi
}
```

## 📋 EFFORT METADATA STRUCTURE

### Redirection Tracking in orchestrator-state-v3.json
```json
{
  "efforts": {
    "E1.2.3-fix-bug-001": {
      "originally_assigned_bug": "BUG-001",
      "actual_bug_fixed": "BUG-007",
      "redirected_per_r525": true,
      "redirection_reason": "BUG-001 is duplicate of BUG-007",
      "redirection_timestamp": "2025-10-06T10:30:00Z"
    }
  }
}
```

## Common Scenarios

### Scenario 1: Agent Assigned Duplicate Bug
```bash
# Orchestrator assigns BUG-001 to SW Engineer
sw_engineer_bug_fix_startup "BUG-001"

# Output: 🚀 SW Engineer starting bug fix for BUG-001
#         🔍 R525: Pre-work bug validation for BUG-001
#         ⚠️ R525: BUG-001 is a DUPLICATE of BUG-007
#         🔄 REDIRECTING work to canonical bug BUG-007
#         🔄 Redirected from BUG-001 to BUG-007
#         ✅ Working on BUG-007 (canonical bug)

# SW Engineer implements fix for BUG-007
# When BUG-007 is fixed, R524 propagates status to BUG-001
# Both bugs marked as fixed from single fix effort
```

### Scenario 2: Manual Bug Assignment
```bash
# Human assigns BUG-006 to effort
orchestrator_assign_bug_fix "BUG-006" "E1.2.5"

# Output: 📋 Orchestrator: Assigning BUG-006 to E1.2.5
#         🔍 R525: Pre-work bug validation for BUG-006
#         ⚠️ R525: BUG-006 is duplicate - assigning canonical BUG-007 instead
#         ✅ Created fix effort E1.2.5 for BUG-007

# Effort metadata records: originally_assigned_bug = BUG-006
#                          actual_bug_fixed = BUG-007
#                          redirected_per_r525 = true
```

### Scenario 3: Code Review Catches Violation
```bash
# Effort E1.2.6 incorrectly fixed BUG-001 directly (bypassed R525)
code_reviewer_bug_fix_validation "E1.2.6"

# Output: 🔍 Code Reviewer: Validating bug fix effort E1.2.6
#         ❌ R525 VIOLATION: Fixed a duplicate bug directly!
#         Bug: BUG-001
#         This should have been redirected to canonical per R525

# Code Reviewer: REJECT review, require rework
```

## 🔧 UTILITY FUNCTIONS

### Get All Canonical Bugs (For Assignment)
```bash
# Get list of canonical bugs that need fixing
get_canonical_bugs_needing_fix() {
    jq -r '
        .bug_registry[] |
        select(.is_duplicate == false and
               (.fix_status == "OPEN" or .fix_status == "IN_PROGRESS")) |
        .bug_id
    ' orchestrator-state-v3.json
}
```

### Check If Bug Assignment Is Safe
```bash
# Before assignment, verify bug is canonical and open
is_bug_safe_for_assignment() {
    local bug_id="$1"

    IS_CANONICAL=$(jq -r --arg bug "$bug_id" '
        .bug_registry[] |
        select(.bug_id == $bug) |
        .is_duplicate == false
    ' orchestrator-state-v3.json)

    IS_OPEN=$(jq -r --arg bug "$bug_id" '
        .bug_registry[] |
        select(.bug_id == $bug) |
        .fix_status | IN("OPEN", "IN_PROGRESS")
    ' orchestrator-state-v3.json)

    if [[ "$IS_CANONICAL" == "true" ]] && [[ "$IS_OPEN" == "true" ]]; then
        return 0  # Safe
    else
        return 1  # Not safe
    fi
}
```

## Grading Impact

### COMPLIANCE BONUS (+20%)
- All work redirected from duplicates to canonical
- Complete redirection audit trail
- No wasted effort on duplicates

### VIOLATIONS
- Working on duplicate bug directly: **-60%** (WASTED EFFORT)
- Missing redirection when duplicate detected: **-30%**
- Assigning duplicate bug without redirection: **-40%**
- No redirection metadata in effort: **-25%**

## Related Rules
- R522: Duplicate Bug Detection Protocol (prevent creation)
- R523: Duplicate Bug Linking Protocol (link duplicates)
- R524: Bug Status Propagation Protocol (propagate fix status to duplicates)

## Remember

**"Check before work - redirect if duplicate"**
**"Never work on duplicates - always canonical"**
**"One fix on canonical = all duplicates auto-updated via R524"**
**"Record all redirections for audit trail"**

R525 ensures that no agent ever wastes time fixing a duplicate bug. All work is automatically redirected to the canonical bug, and when that canonical bug is fixed, R524 propagates the fix status to all duplicates. This is maximum efficiency!

## Date Added
2025-10-06

## Changelog
- 2025-10-06: Initial creation to prevent wasted effort on duplicate bugs through automatic work redirection
