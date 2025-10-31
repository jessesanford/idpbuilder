#!/bin/bash

# Software Factory 2.0 - Effort Scope Validation Tool
# Enforces R371 (Effort Scope Immutability) and R372 (Theme Enforcement)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "======================================"
echo "📋 EFFORT SCOPE VALIDATION TOOL"
echo "Enforcing R371 & R372 Supreme Laws"
echo "======================================"

# Check if we're in an effort directory
if [ ! -f ".software-factory/IMPLEMENTATION-PLAN.md" ]; then
    echo -e "${RED}❌ FATAL: Not in an effort directory!${NC}"
    echo "Missing .software-factory/IMPLEMENTATION-PLAN.md"
    exit 1
fi

# Function to validate scope
validate_scope() {
    echo -e "\n${YELLOW}📊 R371: SCOPE IMMUTABILITY CHECK${NC}"
    echo "----------------------------------------"

    # Get all changed files
    CHANGED_FILES=$(git diff --name-only origin/main 2>/dev/null || git diff --name-only main 2>/dev/null || echo "")

    if [ -z "$CHANGED_FILES" ]; then
        echo -e "${YELLOW}⚠️ No changes detected yet${NC}"
        return 0
    fi

    # Count files
    FILE_COUNT=$(echo "$CHANGED_FILES" | wc -l)
    echo "Files changed: $FILE_COUNT"

    # Check each file against plan
    VIOLATIONS=0
    VIOLATION_LIST=""

    for file in $CHANGED_FILES; do
        # Skip .software-factory directory itself
        if [[ "$file" == .software-factory/* ]]; then
            continue
        fi

        # Check if file is in plan
        if ! grep -q "$file" .software-factory/IMPLEMENTATION-PLAN.md 2>/dev/null; then
            echo -e "${RED}  ❌ SCOPE VIOLATION: $file NOT IN PLAN!${NC}"
            VIOLATIONS=$((VIOLATIONS + 1))
            VIOLATION_LIST="$VIOLATION_LIST\n  - $file"
        else
            echo -e "${GREEN}  ✅ $file is in scope${NC}"
        fi
    done

    if [ $VIOLATIONS -gt 0 ]; then
        echo -e "\n${RED}🔴🔴🔴 R371 CRITICAL VIOLATION!${NC}"
        echo -e "${RED}Found $VIOLATIONS files outside effort scope!${NC}"
        echo -e "${RED}Violations:$VIOLATION_LIST${NC}"
        echo -e "\n${RED}STOPPING: Adding unplanned files is FORBIDDEN!${NC}"
        return 371
    else
        echo -e "${GREEN}✅ All files within scope!${NC}"
    fi

    # Warn about potential scope creep
    if [ $FILE_COUNT -gt 20 ]; then
        echo -e "${YELLOW}⚠️ WARNING: >20 files suggests possible scope creep${NC}"
    fi

    if [ $FILE_COUNT -gt 100 ]; then
        echo -e "${RED}🔴 CRITICAL: >100 files is CATASTROPHIC scope violation!${NC}"
        echo -e "${RED}This is what destroyed the gitea-client-split-001 branch!${NC}"
        return 371
    fi
}

# Function to validate theme coherence
validate_theme() {
    echo -e "\n${YELLOW}🎯 R372: THEME COHERENCE CHECK${NC}"
    echo "----------------------------------------"

    # Count unique packages/directories
    PACKAGE_COUNT=$(git diff --name-only origin/main 2>/dev/null |
                    cut -d'/' -f1-2 |
                    sort -u |
                    wc -l || echo "0")

    echo "Unique packages/concerns touched: $PACKAGE_COUNT"

    if [ $PACKAGE_COUNT -gt 3 ]; then
        echo -e "${RED}🔴🔴🔴 R372 KITCHEN SINK DETECTED!${NC}"
        echo -e "${RED}Effort modifies $PACKAGE_COUNT different packages/concerns!${NC}"
        echo -e "${RED}This violates the ONE THEME PER EFFORT law!${NC}"
        return 372
    fi

    # Check for mixed concerns
    HAS_BUILD=$(git diff --name-only origin/main 2>/dev/null | grep -E "(Makefile|go.mod|package.json)" | wc -l || echo "0")
    HAS_CODE=$(git diff --name-only origin/main 2>/dev/null | grep -E "pkg/|src/|lib/" | wc -l || echo "0")
    HAS_INFRA=$(git diff --name-only origin/main 2>/dev/null | grep -E "docker|devcontainer|k8s|helm" | wc -l || echo "0")
    HAS_DOCS=$(git diff --name-only origin/main 2>/dev/null | grep -E "\.md$|docs/" | wc -l || echo "0")

    CONCERN_COUNT=0
    [ $HAS_BUILD -gt 0 ] && CONCERN_COUNT=$((CONCERN_COUNT + 1)) && echo "  - Build system changes detected"
    [ $HAS_CODE -gt 0 ] && CONCERN_COUNT=$((CONCERN_COUNT + 1)) && echo "  - Code changes detected"
    [ $HAS_INFRA -gt 0 ] && CONCERN_COUNT=$((CONCERN_COUNT + 1)) && echo "  - Infrastructure changes detected"
    [ $HAS_DOCS -gt 3 ] && CONCERN_COUNT=$((CONCERN_COUNT + 1)) && echo "  - Major documentation changes detected"

    if [ $CONCERN_COUNT -gt 1 ]; then
        echo -e "${YELLOW}⚠️ WARNING: Multiple concerns ($CONCERN_COUNT) detected!${NC}"
        echo "Verify these all support a SINGLE theme"
    fi

    echo -e "${GREEN}✅ Theme coherence check passed${NC}"
}

# Function to check for split sanity
check_split_sanity() {
    echo -e "\n${YELLOW}🔍 SPLIT SANITY CHECK${NC}"
    echo "----------------------------------------"

    # Check if this is a split branch
    BRANCH_NAME=$(git rev-parse --abbrev-ref HEAD)
    if [[ "$BRANCH_NAME" == *"split"* ]]; then
        echo "This appears to be a split branch: $BRANCH_NAME"

        # Get file count
        FILE_COUNT=$(git diff --name-only origin/main | wc -l)

        # Check if there's info about original size
        if [ -f ".software-factory/SPLIT-INFO.md" ]; then
            ORIGINAL_FILES=$(grep -oP "Original files: \K\d+" .software-factory/SPLIT-INFO.md || echo "unknown")
            if [ "$ORIGINAL_FILES" != "unknown" ] && [ $FILE_COUNT -gt $ORIGINAL_FILES ]; then
                echo -e "${RED}🔴🔴🔴 CATASTROPHIC SPLIT VIOLATION!${NC}"
                echo -e "${RED}Split has MORE files ($FILE_COUNT) than original ($ORIGINAL_FILES)!${NC}"
                echo -e "${RED}This is EXACTLY what happened to gitea-client-split-001!${NC}"
                return 371
            fi
        fi

        if [ $FILE_COUNT -gt 50 ]; then
            echo -e "${RED}⚠️ WARNING: Split with >50 files is suspicious${NC}"
        fi
    fi
}

# Main execution
main() {
    echo -e "\nStarting validation..."

    # Run all validations
    validate_scope || exit $?
    validate_theme || exit $?
    check_split_sanity || exit $?

    echo -e "\n${GREEN}════════════════════════════════════════${NC}"
    echo -e "${GREEN}✅ ALL SCOPE VALIDATIONS PASSED!${NC}"
    echo -e "${GREEN}════════════════════════════════════════${NC}"
    echo "Your effort maintains:"
    echo "  ✅ Scope immutability (R371)"
    echo "  ✅ Theme coherence (R372)"
    echo "  ✅ Split sanity"
    echo ""
    echo "You may proceed with confidence!"
}

# Run main function
main