#!/usr/bin/env bash
#
# R405 COMPLIANCE CHECKER
#
# Validates that all orchestrator state rules files properly implement R405
# (Automation Continuation Flag) per the rule library specification.
#
# Returns:
#   0 - All states compliant
#   1 - Compliance violations found
#

set -euo pipefail

CLAUDE_PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"
ORCHESTRATOR_STATES_DIR="$CLAUDE_PROJECT_DIR/agent-states/software-factory/orchestrator"
RULE_LIBRARY_FILE="$CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
total_states=0
compliant_states=0
violations=0

echo "🔍 R405 COMPLIANCE CHECKER"
echo "========================================"
echo

# Check that rule library file exists
if [ ! -f "$RULE_LIBRARY_FILE" ]; then
    echo -e "${RED}❌ FATAL: Rule library file not found: $RULE_LIBRARY_FILE${NC}"
    exit 1
fi

echo "📋 Checking all orchestrator states..."
echo

# Iterate through all state directories
for state_dir in "$ORCHESTRATOR_STATES_DIR"/*/; do
    if [ ! -d "$state_dir" ]; then
        continue
    fi

    state_name=$(basename "$state_dir")
    rules_file="$state_dir/rules.md"

    if [ ! -f "$rules_file" ]; then
        echo -e "${YELLOW}⚠️  $state_name: rules.md not found${NC}"
        continue
    fi

    total_states=$((total_states + 1))
    state_compliant=true
    issues=()

    # CHECK 1: R405 mentioned in file
    if ! grep -q "R405" "$rules_file"; then
        issues+=("Missing R405 reference")
        state_compliant=false
    fi

    # CHECK 2: CONTINUE-SOFTWARE-FACTORY mentioned
    if ! grep -q "CONTINUE-SOFTWARE-FACTORY" "$rules_file"; then
        issues+=("Missing CONTINUE-SOFTWARE-FACTORY flag")
        state_compliant=false
    fi

    # CHECK 3: Has completion checklist
    if ! grep -q "STATE COMPLETION CHECKLIST" "$rules_file"; then
        issues+=("Missing STATE COMPLETION CHECKLIST")
        state_compliant=false
    fi

    # CHECK 4: Only ONE enforcement section (no duplicates)
    enforcement_count=$(grep -c "🚨 CHECKLIST ENFORCEMENT 🚨" "$rules_file" || true)
    if [ "$enforcement_count" -gt 1 ]; then
        issues+=("DUPLICATE enforcement section ($enforcement_count found)")
        state_compliant=false
    elif [ "$enforcement_count" -eq 0 ]; then
        issues+=("Missing enforcement section")
        state_compliant=false
    fi

    # CHECK 5: R405 mentioned in enforcement section
    if grep -q "🚨 CHECKLIST ENFORCEMENT 🚨" "$rules_file"; then
        if ! grep -A 20 "🚨 CHECKLIST ENFORCEMENT 🚨" "$rules_file" | grep -q "R405"; then
            issues+=("R405 not mentioned in enforcement section")
            state_compliant=false
        fi
    fi

    # CHECK 6: Has R405 detailed explanation section
    if ! grep -q "🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴" "$rules_file"; then
        issues+=("Missing R405 detailed explanation section")
        state_compliant=false
    fi

    # Report results for this state
    if $state_compliant; then
        compliant_states=$((compliant_states + 1))
        echo -e "${GREEN}✅ $state_name${NC}"
    else
        violations=$((violations + 1))
        echo -e "${RED}❌ $state_name:${NC}"
        for issue in "${issues[@]}"; do
            echo -e "${RED}   - $issue${NC}"
        done
    fi
done

echo
echo "========================================"
echo "📊 COMPLIANCE SUMMARY"
echo "========================================"
echo "Total states checked: $total_states"
echo -e "${GREEN}Compliant states: $compliant_states${NC}"
if [ $violations -gt 0 ]; then
    echo -e "${RED}States with violations: $violations${NC}"
fi
echo

if [ $violations -eq 0 ]; then
    echo -e "${GREEN}✅ ALL STATES R405 COMPLIANT${NC}"
    echo
    exit 0
else
    echo -e "${RED}❌ R405 COMPLIANCE FAILURES DETECTED${NC}"
    echo
    echo "Fix violations and re-run this checker."
    echo
    exit 1
fi
