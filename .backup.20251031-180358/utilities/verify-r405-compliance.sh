#!/bin/bash
# Script to verify R405 automation continuation flag compliance across all state rule files

set -euo pipefail

CLAUDE_PROJECT_DIR="/home/vscode/software-factory-template"
cd "$CLAUDE_PROJECT_DIR"

echo "🔍 R405 COMPLIANCE VERIFICATION REPORT"
echo "======================================"
echo ""
echo "Date: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Statistics
TOTAL=0
COMPLIANT=0
NON_COMPLIANT=0
ERRORS=0

# Find all rule files
RULE_FILES=$(find /home/vscode/software-factory-template/agent-states -name "rules.md" -type f | sort)

echo "📂 Scanning all agent state rule files..."
echo ""

# Array to store non-compliant files
declare -a NON_COMPLIANT_FILES

for file in $RULE_FILES; do
    TOTAL=$((TOTAL + 1))

    # Extract agent and state info from path for better reporting
    AGENT_STATE=$(echo "$file" | sed 's|.*/agent-states/||' | sed 's|/rules.md||')

    # Check if file is readable
    if [ ! -r "$file" ]; then
        echo -e "${RED}❌ ERROR: Cannot read file: $AGENT_STATE${NC}"
        ERRORS=$((ERRORS + 1))
        continue
    fi

    # Check for R405 presence
    if grep -q "R405" "$file"; then
        # Verify it has the correct content (not just any R405 mention)
        # Check for either format (case insensitive for "automation continuation flag")
        if grep -qi "automation.*continuation.*flag\|CONTINUE-SOFTWARE-FACTORY" "$file"; then
            echo -e "${GREEN}✓${NC} Compliant: $AGENT_STATE"
            COMPLIANT=$((COMPLIANT + 1))
        else
            echo -e "${YELLOW}⚠${NC}  Has R405 but missing proper content: $AGENT_STATE"
            NON_COMPLIANT=$((NON_COMPLIANT + 1))
            NON_COMPLIANT_FILES+=("$file")
        fi
    else
        echo -e "${RED}✗${NC} Missing R405: $AGENT_STATE"
        NON_COMPLIANT=$((NON_COMPLIANT + 1))
        NON_COMPLIANT_FILES+=("$file")
    fi
done

echo ""
echo "============================================================"
echo "📊 Compliance Report"
echo "============================================================"
echo "Total state rule files:    $TOTAL"
echo -e "${GREEN}Compliant files:          $COMPLIANT${NC}"
echo -e "${RED}Non-compliant files:      $NON_COMPLIANT${NC}"
echo -e "${YELLOW}Errors encountered:       $ERRORS${NC}"
echo ""

# Calculate compliance percentage
if [ $TOTAL -gt 0 ]; then
    COMPLIANCE_PCT=$((COMPLIANT * 100 / TOTAL))
    echo "Compliance Rate: ${COMPLIANCE_PCT}%"

    if [ $COMPLIANCE_PCT -eq 100 ]; then
        echo -e "${GREEN}🎉 FULL COMPLIANCE ACHIEVED!${NC}"
        echo "All agent states will output automation continuation flags."
    elif [ $COMPLIANCE_PCT -ge 95 ]; then
        echo -e "${YELLOW}⚠️  Nearly compliant, but some files need attention.${NC}"
    else
        echo -e "${RED}❌ Significant compliance gaps detected.${NC}"
    fi
else
    echo -e "${RED}No state rule files found!${NC}"
fi

# List non-compliant files if any
if [ ${#NON_COMPLIANT_FILES[@]} -gt 0 ]; then
    echo ""
    echo "❗ Non-compliant files requiring attention:"
    echo "--------------------------------------------"
    for file in "${NON_COMPLIANT_FILES[@]}"; do
        echo "  - $file"
    done
    echo ""
    echo "To fix these files, run: ./add-r405-to-all-states.sh"
fi

# Verify R405 rule file exists
echo ""
echo "============================================================"
echo "🔍 Checking R405 Rule Definition"
echo "============================================================"

R405_FILE="/home/vscode/software-factory-template/rule-library/R405-automation-continuation-flag.md"
if [ -f "$R405_FILE" ]; then
    echo -e "${GREEN}✓${NC} R405 rule definition exists"

    # Check if it's in the registry
    if grep -q "R405" "/home/vscode/software-factory-template/rule-library/RULE-REGISTRY.md"; then
        echo -e "${GREEN}✓${NC} R405 is registered in RULE-REGISTRY.md"
    else
        echo -e "${RED}✗${NC} R405 is NOT in RULE-REGISTRY.md"
    fi
else
    echo -e "${RED}✗${NC} R405 rule definition file not found!"
fi

# Exit code based on compliance
if [ $NON_COMPLIANT -eq 0 ] && [ $ERRORS -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ All compliance checks passed!${NC}"
    exit 0
else
    echo ""
    echo -e "${RED}❌ Compliance issues detected. Please review and fix.${NC}"
    exit 1
fi