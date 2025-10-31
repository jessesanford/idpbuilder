#!/bin/bash

# 🏗️ ARCHITECTURE VALIDATION SCRIPT
# Enforces R362 - No Architectural Rewrites Without Approval

set -e

echo "🏗️ ========================================="
echo "🏗️ R362 ARCHITECTURE COMPLIANCE VALIDATION"
echo "🏗️ ========================================="
echo ""

# Configuration - Add your project's required libraries here
REQUIRED_LIBRARIES=(
    "go-containerregistry"
    # Add more required libraries as needed
)

FORBIDDEN_PATTERNS=(
    "custom.*http.*registry"
    "MyORM"
    "DirectDB"
    "stub.*implementation"
    "mock.*client"
)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

VIOLATIONS=0
WARNINGS=0

# Function to check required libraries
check_required_libraries() {
    echo "📚 Checking Required Libraries..."
    echo "================================"

    for lib in "${REQUIRED_LIBRARIES[@]}"; do
        if [ -f go.mod ]; then
            if grep -q "$lib" go.mod; then
                echo -e "${GREEN}✅ Found: $lib${NC}"
            else
                echo -e "${RED}❌ MISSING REQUIRED LIBRARY: $lib${NC}"
                ((VIOLATIONS++))
            fi
        fi

        if [ -f package.json ]; then
            if grep -q "\"$lib\"" package.json; then
                echo -e "${GREEN}✅ Found: $lib${NC}"
            else
                echo -e "${RED}❌ MISSING REQUIRED LIBRARY: $lib${NC}"
                ((VIOLATIONS++))
            fi
        fi
    done
    echo ""
}

# Function to check for forbidden patterns
check_forbidden_patterns() {
    echo "🚫 Checking for Forbidden Patterns..."
    echo "====================================="

    for pattern in "${FORBIDDEN_PATTERNS[@]}"; do
        echo -n "Checking for '$pattern'... "

        # Search in source directories
        if find . -type f \( -name "*.go" -o -name "*.js" -o -name "*.py" -o -name "*.ts" \) \
           -not -path "*/test/*" -not -path "*/tests/*" -not -path "*/vendor/*" \
           -exec grep -l "$pattern" {} \; 2>/dev/null | head -1 | grep -q .; then
            echo -e "${RED}FOUND!${NC}"
            echo -e "${RED}❌ FORBIDDEN PATTERN DETECTED: $pattern${NC}"

            # Show where it was found
            find . -type f \( -name "*.go" -o -name "*.js" -o -name "*.py" -o -name "*.ts" \) \
                -not -path "*/test/*" -not -path "*/tests/*" -not -path "*/vendor/*" \
                -exec grep -l "$pattern" {} \; 2>/dev/null | while read file; do
                echo -e "${RED}   Found in: $file${NC}"
            done

            ((VIOLATIONS++))
        else
            echo -e "${GREEN}OK${NC}"
        fi
    done
    echo ""
}

# Function to check for architectural changes
check_architectural_changes() {
    echo "🔄 Checking for Architectural Changes..."
    echo "========================================"

    if [ -z "$1" ]; then
        echo -e "${YELLOW}⚠️  No base branch specified, skipping diff checks${NC}"
        return
    fi

    BASE_BRANCH="$1"

    # Check for library removals
    if [ -f go.mod ]; then
        REMOVED_LIBS=$(git diff "$BASE_BRANCH"...HEAD -- go.mod | grep "^-.*github.com\|^-.*golang.org" | grep -v "^---" || true)
        if [ -n "$REMOVED_LIBS" ]; then
            echo -e "${RED}❌ LIBRARIES REMOVED:${NC}"
            echo "$REMOVED_LIBS"
            ((VIOLATIONS++))
        fi
    fi

    # Check for major file deletions
    DELETED_FILES=$(git diff --name-status "$BASE_BRANCH"...HEAD | grep "^D" | grep -E "\.(go|js|py|ts)$" || true)
    if [ -n "$DELETED_FILES" ]; then
        FILE_COUNT=$(echo "$DELETED_FILES" | wc -l)
        if [ "$FILE_COUNT" -gt 5 ]; then
            echo -e "${RED}❌ MAJOR FILES DELETED: $FILE_COUNT files${NC}"
            echo "$DELETED_FILES" | head -10
            ((VIOLATIONS++))
        fi
    fi

    echo ""
}

# Function to validate implementation matches plan
check_plan_compliance() {
    echo "📋 Checking Plan Compliance..."
    echo "=============================="

    # Look for effort plan files
    PLAN_FILES=$(find . -name "*-plan.md" -o -name "*-effort-plan.md" -o -name "EFFORT-PLAN.md" 2>/dev/null)

    if [ -z "$PLAN_FILES" ]; then
        echo -e "${YELLOW}⚠️  No plan files found to validate against${NC}"
    else
        echo -e "${GREEN}✅ Found plan files to validate${NC}"
        # Additional plan validation could go here
    fi

    echo ""
}

# Function to generate compliance report
generate_report() {
    echo "📊 ========================================="
    echo "📊 ARCHITECTURE COMPLIANCE REPORT"
    echo "📊 ========================================="
    echo ""

    if [ $VIOLATIONS -eq 0 ]; then
        echo -e "${GREEN}✅ ARCHITECTURE FULLY COMPLIANT${NC}"
        echo "All architectural requirements satisfied!"
    else
        echo -e "${RED}🔴🔴🔴 ARCHITECTURE VIOLATIONS DETECTED: $VIOLATIONS 🔴🔴🔴${NC}"
        echo ""
        echo "CRITICAL: This implementation violates R362!"
        echo "Changing approved architecture is ABSOLUTELY FORBIDDEN!"
        echo ""
        echo "Required Actions:"
        echo "1. STOP all implementation immediately"
        echo "2. Revert unauthorized changes"
        echo "3. If change is necessary, get EXPLICIT user approval"
        echo "4. Update planning documents with approval"
        echo "5. Only then proceed with approved changes"
    fi

    if [ $WARNINGS -gt 0 ]; then
        echo ""
        echo -e "${YELLOW}⚠️  Warnings: $WARNINGS${NC}"
    fi

    echo ""
    echo "========================================="
    echo "Validation completed at: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "========================================="
}

# Main execution
main() {
    # Parse arguments
    BASE_BRANCH="${1:-main}"

    echo "Starting architecture validation..."
    echo "Base branch: $BASE_BRANCH"
    echo ""

    # Run all checks
    check_required_libraries
    check_forbidden_patterns
    check_architectural_changes "$BASE_BRANCH"
    check_plan_compliance

    # Generate report
    generate_report

    # Exit with appropriate code
    if [ $VIOLATIONS -gt 0 ]; then
        exit 362  # R362 violation exit code
    else
        exit 0
    fi
}

# Run main function
main "$@"