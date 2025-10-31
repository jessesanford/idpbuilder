#!/bin/bash
# File: utilities/validate-state-checklists.sh
# Purpose: Validate all orchestrator states have proper R510 checklists
# Version: 1.0
# Created: 2025-10-07

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

CLAUDE_PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}"
ERRORS=0
WARNINGS=0
VALIDATED=0

echo "🔍 R510 State Checklist Validator"
echo "=================================="
echo ""

validate_state_checklist() {
    local state_dir="$1"
    local state_name=$(basename "$state_dir")
    local rules_file="${state_dir}/rules.md"
    local state_errors=0

    # Skip DEPRECATED states
    if [[ "$state_dir" == *"/DEPRECATED/"* ]]; then
        return 0
    fi

    # Check file exists
    if [ ! -f "$rules_file" ]; then
        echo -e "${RED}❌ $state_name: Missing rules.md${NC}"
        ((ERRORS++))
        return 1
    fi

    # Check checklist section exists
    if ! grep -q "## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST" "$rules_file"; then
        echo -e "${RED}❌ $state_name: Missing MANDATORY EXECUTION CHECKLIST section${NC}"
        ((ERRORS++))
        ((state_errors++))
    fi

    # Check required sections
    if ! grep -q "### BLOCKING REQUIREMENTS" "$rules_file"; then
        echo -e "${RED}❌ $state_name: Missing BLOCKING REQUIREMENTS section${NC}"
        ((ERRORS++))
        ((state_errors++))
    fi

    if ! grep -q "### STANDARD EXECUTION TASKS" "$rules_file"; then
        echo -e "${YELLOW}⚠️  $state_name: Missing STANDARD EXECUTION TASKS section${NC}"
        ((WARNINGS++))
    fi

    if ! grep -q "### EXIT REQUIREMENTS" "$rules_file"; then
        echo -e "${RED}❌ $state_name: Missing EXIT REQUIREMENTS section${NC}"
        ((ERRORS++))
        ((state_errors++))
    fi

    # Validate EXIT REQUIREMENTS content
    if grep -q "### EXIT REQUIREMENTS" "$rules_file"; then
        if ! grep -q "Update state file.*R288" "$rules_file"; then
            echo -e "${RED}❌ $state_name: EXIT REQUIREMENTS missing R288 state update${NC}"
            ((ERRORS++))
            ((state_errors++))
        fi

        if ! grep -q "Save TODOs.*R287" "$rules_file"; then
            echo -e "${RED}❌ $state_name: EXIT REQUIREMENTS missing R287 TODO save${NC}"
            ((ERRORS++))
            ((state_errors++))
        fi

        if ! grep -q "CONTINUE-SOFTWARE-FACTORY.*R405" "$rules_file"; then
            echo -e "${RED}❌ $state_name: EXIT REQUIREMENTS missing R405 flag${NC}"
            ((ERRORS++))
            ((state_errors++))
        fi

        if ! grep -q "exit 0" "$rules_file"; then
            echo -e "${YELLOW}⚠️  $state_name: EXIT REQUIREMENTS missing exit 0${NC}"
            ((WARNINGS++))
        fi
    fi

    # Check R510 rule reference
    if ! grep -q "R510" "$rules_file"; then
        echo -e "${YELLOW}⚠️  $state_name: Missing R510 rule reference${NC}"
        ((WARNINGS++))
    fi

    # Check acknowledgment protocol mentioned
    if ! grep -q "CHECKLIST\[" "$rules_file"; then
        echo -e "${YELLOW}⚠️  $state_name: Missing acknowledgment protocol examples${NC}"
        ((WARNINGS++))
    fi

    if [ $state_errors -eq 0 ]; then
        echo -e "${GREEN}✅ $state_name: Checklist valid${NC}"
        ((VALIDATED++))
    fi

    return $state_errors
}

# Find all orchestrator state directories
echo "Scanning orchestrator states..."
echo ""

for state_dir in "$CLAUDE_PROJECT_DIR/agent-states/software-factory/orchestrator"/*; do
    if [ -d "$state_dir" ]; then
        validate_state_checklist "$state_dir"
    fi
done

echo ""
echo "=================================="
echo "Validation Summary"
echo "=================================="
echo -e "States validated: ${GREEN}$VALIDATED${NC}"
echo -e "Errors found: ${RED}$ERRORS${NC}"
echo -e "Warnings found: ${YELLOW}$WARNINGS${NC}"
echo ""

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅✅✅ ALL STATE CHECKLISTS VALIDATED PROJECT_DONEFULLY ✅✅✅${NC}"
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo -e "${YELLOW}⚠️  Validation passed with $WARNINGS warnings${NC}"
    exit 0
else
    echo -e "${RED}❌❌❌ FOUND $ERRORS CHECKLIST ERRORS ❌❌❌${NC}"
    exit 1
fi
