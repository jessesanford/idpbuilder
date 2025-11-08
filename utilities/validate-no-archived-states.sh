#!/bin/bash
# Standalone Validation: Check for ARCHIVED State Contamination
# CRITICAL: ARCHIVED states should NEVER exist in active projects
# This utility can be run manually to check for contamination

set -euo pipefail

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BOLD='\033[1m'
NC='\033[0m'

# Get project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BOLD}ARCHIVED State Contamination Check${NC}"
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Check for ARCHIVED directories
ARCHIVED_DIRS=$(find agent-states -type d -name "ARCHIVED" 2>/dev/null || true)

if [ -n "$ARCHIVED_DIRS" ]; then
    echo -e "${RED}${BOLD}❌ CONTAMINATION DETECTED${NC}"
    echo ""
    echo -e "${RED}ARCHIVED state directories found:${NC}"
    echo "$ARCHIVED_DIRS" | sed 's/^/  - /'
    echo ""
    echo -e "${YELLOW}${BOLD}WHY THIS IS CRITICAL:${NC}"
    echo "  • ARCHIVED states are deprecated/obsolete"
    echo "  • They can confuse agents trying to read rules"
    echo "  • They violate rule synchronization requirements"
    echo "  • They contaminate the codebase with dead code"
    echo ""
    echo -e "${YELLOW}${BOLD}HOW TO FIX:${NC}"
    echo "  rm -rf agent-states/ARCHIVED"
    echo ""
    echo -e "${YELLOW}${BOLD}ROOT CAUSE:${NC}"
    echo "  This contamination typically comes from:"
    echo "  • tools/upgrade.sh copying from contaminated template"
    echo "  • tools/setup*.sh copying from contaminated source"
    echo ""
    echo -e "${YELLOW}${BOLD}PREVENTION:${NC}"
    echo "  • upgrade.sh and setup*.sh now exclude ARCHIVED (if fixed)"
    echo "  • Pre-commit hooks block ARCHIVED commits"
    echo "  • This utility can be run manually to check"
    echo ""
    exit 1
fi

# Check for any ARCHIVED files in git
ARCHIVED_FILES=$(find agent-states -type f -path "*/ARCHIVED/*" 2>/dev/null || true)

if [ -n "$ARCHIVED_FILES" ]; then
    FILE_COUNT=$(echo "$ARCHIVED_FILES" | wc -l)
    echo -e "${RED}${BOLD}❌ ARCHIVED FILES DETECTED${NC}"
    echo ""
    echo -e "${RED}Found $FILE_COUNT files in ARCHIVED directories${NC}"
    echo ""
    echo "First 10 files:"
    echo "$ARCHIVED_FILES" | head -10 | sed 's/^/  - /'
    if [ $FILE_COUNT -gt 10 ]; then
        echo "  ... and $((FILE_COUNT - 10)) more"
    fi
    echo ""
    echo -e "${YELLOW}${BOLD}FIX:${NC} rm -rf agent-states/ARCHIVED"
    echo ""
    exit 1
fi

# All clear!
echo -e "${GREEN}${BOLD}✅ NO CONTAMINATION DETECTED${NC}"
echo ""
echo -e "${GREEN}No ARCHIVED state directories found${NC}"
echo -e "${GREEN}Project is clean!${NC}"
echo ""
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
exit 0
