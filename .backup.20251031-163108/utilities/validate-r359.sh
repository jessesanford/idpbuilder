#!/bin/bash
# R359 Validation Script - Detect and prevent code deletion for size limits
# This is SUPREME LAW #6 - PENALTY: IMMEDIATE TERMINATION (-1000%)

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔴🔴🔴 R359 VALIDATION: CODE DELETION PROHIBITION CHECK 🔴🔴🔴"
echo "================================================================"

# Function to check deletions in a branch
check_branch_deletions() {
    local branch="$1"
    local base="${2:-main}"

    echo -e "\n📊 Checking branch: ${YELLOW}$branch${NC}"
    echo -e "   Against base: ${YELLOW}$base${NC}"

    # Check if branch exists
    if ! git show-ref --verify --quiet "refs/heads/$branch"; then
        echo -e "${YELLOW}⚠️  Branch $branch not found locally${NC}"
        return 1
    fi

    # Get deletion statistics
    local deleted_lines=$(git diff --numstat "$base..$branch" 2>/dev/null | awk '{sum+=$2} END {print sum}' || echo "0")

    if [ "$deleted_lines" -gt 100 ]; then
        echo -e "${RED}🔴🔴🔴 R359 VIOLATION DETECTED!${NC}"
        echo -e "${RED}Branch deletes $deleted_lines lines of code!${NC}"
        echo ""
        echo -e "${RED}CRITICAL FILES DELETED:${NC}"
        git diff --name-status "$base..$branch" | grep "^D" | head -20
        echo ""
        echo -e "${RED}⛔ THIS IS A CATASTROPHIC VIOLATION!${NC}"
        echo -e "${RED}Deleting code to meet size limits is ABSOLUTELY FORBIDDEN!${NC}"
        echo ""
        echo "CORRECT APPROACH:"
        echo "✅ Split NEW work into 800-line pieces"
        echo "✅ Keep ALL existing code"
        echo "✅ Each split ADDS to the codebase"
        echo ""
        return 359
    elif [ "$deleted_lines" -gt 50 ]; then
        echo -e "${YELLOW}⚠️  WARNING: $deleted_lines lines deleted - review needed${NC}"
        echo "   Verify these are legitimate refactoring deletions"
    else
        echo -e "${GREEN}✅ Branch deletions within acceptable range ($deleted_lines lines)${NC}"
    fi

    # Check for critical file deletions
    local critical_deletions=$(git diff --name-status "$base..$branch" 2>/dev/null | grep "^D" | grep -E "main\.(go|py|js|ts)|LICENSE|README|Makefile" || true)
    if [ -n "$critical_deletions" ]; then
        echo -e "${RED}🔴🔴🔴 CRITICAL FILE DELETIONS DETECTED!${NC}"
        echo "$critical_deletions"
        echo -e "${RED}Never delete main files, LICENSE, or README!${NC}"
        return 359
    fi

    return 0
}

# Function to check all effort branches
check_all_efforts() {
    echo -e "\n🔍 Scanning all effort branches for R359 violations..."

    local violation_count=0
    local branches=$(git branch -r | grep -E "effort|split" | sed 's/origin\///' || true)

    if [ -z "$branches" ]; then
        echo "No effort or split branches found"
        return 0
    fi

    for branch in $branches; do
        if ! check_branch_deletions "$branch"; then
            ((violation_count++))
        fi
    done

    if [ "$violation_count" -gt 0 ]; then
        echo -e "\n${RED}🔴 FOUND $violation_count R359 VIOLATIONS!${NC}"
        echo -e "${RED}THESE MUST BE FIXED IMMEDIATELY!${NC}"
        return 359
    else
        echo -e "\n${GREEN}✅ All branches comply with R359${NC}"
    fi
}

# Function to provide education about R359
educate_r359() {
    echo ""
    echo "📚 R359 EDUCATION - UNDERSTANDING SIZE LIMITS"
    echo "=============================================="
    echo ""
    echo "THE RULE:"
    echo "  The 800-line limit applies to NEW CODE ONLY"
    echo ""
    echo "WHAT THIS MEANS:"
    echo "  ✅ If repo has 10,000 lines and you add 800, total becomes 10,800"
    echo "  ✅ Repository WILL grow with each effort (EXPECTED)"
    echo "  ❌ NEVER delete existing code to fit within limits"
    echo ""
    echo "SPLITTING CORRECTLY:"
    echo "  Effort needs 2,000 lines of NEW code?"
    echo "  Split 1: Add 800 lines → Repo grows to 10,800"
    echo "  Split 2: Add 800 lines → Repo grows to 11,600"
    echo "  Split 3: Add 400 lines → Repo grows to 12,000"
    echo ""
    echo "THE DISASTER THAT CREATED THIS RULE:"
    echo "  Agents deleted 9,552 lines including:"
    echo "  - Entire pkg/ directory"
    echo "  - main.go, LICENSE, README"
    echo "  - All to make 595 lines 'fit' in 800"
    echo ""
    echo "REMEMBER: Size limits prevent PR review burden, NOT repository growth!"
}

# Main execution
main() {
    case "${1:-check}" in
        check)
            if [ -n "${2:-}" ]; then
                check_branch_deletions "$2" "${3:-main}"
            else
                check_all_efforts
            fi
            ;;
        educate)
            educate_r359
            ;;
        help)
            echo "Usage: $0 [command] [options]"
            echo ""
            echo "Commands:"
            echo "  check [branch] [base]  - Check specific branch or all efforts"
            echo "  educate               - Show R359 education material"
            echo "  help                  - Show this help"
            echo ""
            echo "Examples:"
            echo "  $0                    - Check all effort branches"
            echo "  $0 check effort-001   - Check specific branch"
            echo "  $0 check split-002 split-001  - Check split against previous"
            echo "  $0 educate           - Learn about R359"
            ;;
        *)
            echo "Unknown command: $1"
            echo "Run '$0 help' for usage"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"