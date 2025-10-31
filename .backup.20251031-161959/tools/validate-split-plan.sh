#!/bin/bash
#
# validate-split-plan.sh - Validate split plans for R359 compliance
#
# This script validates that split plans don't violate R359 by requiring
# deletion of existing code. Splits must ADD new functionality, not divide
# existing code.
#
# Usage: ./validate-split-plan.sh <split-plan-file> [base-branch]
#
# Exit codes:
#   0 - Split plan is valid
#   359 - R359 violation detected (plans to delete code)
#   1 - General error

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# R359 violation indicators
DELETION_PATTERNS=(
    "git rm"
    "delete"
    "remove"
    "exclude"
    "take out"
    "strip"
    "eliminate"
    "discard"
    "omit"
    "drop"
    "divide existing"
    "partition existing"
    "split existing"
)

# Function to print colored output
print_color() {
    local color=$1
    shift
    echo -e "${color}$@${NC}"
}

# Function to check for R359 violations in text
check_for_violations() {
    local file=$1
    local violations_found=0

    print_color "$BLUE" "🔍 Checking for R359 violations in split plan..."

    for pattern in "${DELETION_PATTERNS[@]}"; do
        if grep -qi "$pattern" "$file"; then
            print_color "$RED" "🔴 R359 VIOLATION: Found '$pattern' in split plan!"
            grep -n -i "$pattern" "$file" | while read -r line; do
                print_color "$YELLOW" "  Line: $line"
            done
            violations_found=1
        fi
    done

    return $violations_found
}

# Function to check if plan preserves functionality
check_functionality_preservation() {
    local file=$1

    print_color "$BLUE" "🔍 Checking if plan preserves all functionality..."

    # Look for statements about preservation
    if grep -qi "all.*functionality.*preserved\|preserve.*existing\|maintains.*all\|keeps.*all" "$file"; then
        print_color "$GREEN" "✅ Plan explicitly states functionality preservation"
        return 0
    else
        print_color "$YELLOW" "⚠️  Plan should explicitly state that all functionality is preserved"
        return 1
    fi
}

# Function to check if splits add new code
check_splits_add_code() {
    local file=$1

    print_color "$BLUE" "🔍 Checking if splits add new code..."

    # Look for "adds" or "implements" keywords
    local add_count=$(grep -ci "adds\|implements\|creates\|introduces" "$file" || true)

    if [ "$add_count" -gt 0 ]; then
        print_color "$GREEN" "✅ Found $add_count references to adding/implementing new code"
        return 0
    else
        print_color "$RED" "❌ Split plan doesn't mention adding new code"
        return 1
    fi
}

# Function to validate split sizes
check_split_sizes() {
    local file=$1

    print_color "$BLUE" "🔍 Checking split sizes..."

    # Extract line counts (look for patterns like "XXX lines" or "~XXX lines")
    local sizes=$(grep -oE '[~]?[0-9]+ lines' "$file" | grep -oE '[0-9]+' || true)

    if [ -z "$sizes" ]; then
        print_color "$YELLOW" "⚠️  No explicit line counts found in split plan"
        return 1
    fi

    local violations=0
    while read -r size; do
        if [ "$size" -gt 800 ]; then
            print_color "$RED" "❌ Split with $size lines exceeds 800-line limit!"
            violations=$((violations + 1))
        else
            print_color "$GREEN" "✅ Split with $size lines is within limit"
        fi
    done <<< "$sizes"

    return $violations
}

# Function to check for R359 compliance statement
check_compliance_statement() {
    local file=$1

    print_color "$BLUE" "🔍 Checking for R359 compliance statement..."

    if grep -qi "R359" "$file"; then
        print_color "$GREEN" "✅ Split plan references R359"

        if grep -qi "R359.*compliance\|comply.*R359\|compliant.*R359" "$file"; then
            print_color "$GREEN" "✅ Split plan has R359 compliance statement"
            return 0
        else
            print_color "$YELLOW" "⚠️  Split plan mentions R359 but lacks explicit compliance statement"
            return 1
        fi
    else
        print_color "$RED" "❌ Split plan doesn't reference R359"
        return 1
    fi
}

# Function to analyze base branch strategy
check_base_branch_strategy() {
    local file=$1

    print_color "$BLUE" "🔍 Checking base branch strategy..."

    if grep -qi "same base\|from main\|from master" "$file"; then
        print_color "$GREEN" "✅ Splits use consistent base branch"
        return 0
    else
        print_color "$YELLOW" "⚠️  Base branch strategy not clearly stated"
        return 1
    fi
}

# Main validation function
validate_split_plan() {
    local plan_file=$1
    local base_branch=${2:-main}

    if [ ! -f "$plan_file" ]; then
        print_color "$RED" "❌ Split plan file not found: $plan_file"
        exit 1
    fi

    print_color "$BLUE" "═══════════════════════════════════════════════════════"
    print_color "$BLUE" "    R359 SPLIT PLAN VALIDATOR"
    print_color "$BLUE" "═══════════════════════════════════════════════════════"
    echo
    print_color "$BLUE" "📋 Validating: $plan_file"
    print_color "$BLUE" "📋 Base branch: $base_branch"
    echo

    local total_issues=0
    local critical_violations=0

    # Run validation checks
    if ! check_for_violations "$plan_file"; then
        critical_violations=$((critical_violations + 1))
    fi
    echo

    if ! check_functionality_preservation "$plan_file"; then
        total_issues=$((total_issues + 1))
    fi
    echo

    if ! check_splits_add_code "$plan_file"; then
        total_issues=$((total_issues + 1))
    fi
    echo

    if ! check_split_sizes "$plan_file"; then
        total_issues=$((total_issues + 1))
    fi
    echo

    if ! check_compliance_statement "$plan_file"; then
        total_issues=$((total_issues + 1))
    fi
    echo

    if ! check_base_branch_strategy "$plan_file"; then
        total_issues=$((total_issues + 1))
    fi
    echo

    # Summary
    print_color "$BLUE" "═══════════════════════════════════════════════════════"
    print_color "$BLUE" "    VALIDATION SUMMARY"
    print_color "$BLUE" "═══════════════════════════════════════════════════════"

    if [ "$critical_violations" -gt 0 ]; then
        print_color "$RED" "🔴🔴🔴 CRITICAL R359 VIOLATIONS DETECTED! 🔴🔴🔴"
        print_color "$RED" "This split plan would require deleting existing code!"
        print_color "$RED" "Splits must ADD new functionality, not divide existing code!"
        print_color "$RED" ""
        print_color "$RED" "REMEMBER:"
        print_color "$RED" "  • Each split implements a PORTION of NEW functionality"
        print_color "$RED" "  • All splits start from the SAME base branch"
        print_color "$RED" "  • NO split should delete existing code"
        print_color "$RED" "  • The 800-line limit applies to NEW code only"
        echo
        exit 359
    elif [ "$total_issues" -gt 0 ]; then
        print_color "$YELLOW" "⚠️  Split plan has $total_issues issues to address"
        print_color "$YELLOW" "Please review and update the plan before implementation"
        exit 1
    else
        print_color "$GREEN" "✅ Split plan is R359 COMPLIANT!"
        print_color "$GREEN" "✅ Plan correctly partitions NEW work without deleting code"
        print_color "$GREEN" "✅ Safe to proceed with implementation"
        echo
    fi
}

# Help message
show_help() {
    cat << EOF
validate-split-plan.sh - Validate split plans for R359 compliance

USAGE:
    $0 <split-plan-file> [base-branch]

DESCRIPTION:
    Validates that a split plan complies with R359 (no deletion of existing code).
    Splits must ADD new functionality, not divide existing code.

ARGUMENTS:
    split-plan-file    Path to the split plan markdown file
    base-branch        Base branch name (default: main)

EXIT CODES:
    0   - Split plan is valid and R359 compliant
    359 - R359 violation detected (plan requires deleting code)
    1   - General validation error or issues found

EXAMPLES:
    $0 SPLIT-PLAN-E1.1.md
    $0 phase1/wave1/SPLIT-PLAN.md development

R359 REMINDER:
    Splits partition NEW additions into reviewable chunks.
    They NEVER delete or divide existing code.
    The repository GROWS with each split.

EOF
}

# Main script execution
main() {
    if [ $# -eq 0 ] || [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
        show_help
        exit 0
    fi

    validate_split_plan "$@"
}

main "$@"