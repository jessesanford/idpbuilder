#!/bin/bash
# Validates CONTINUE-SOFTWARE-FACTORY flag usage in orchestrator state rules
# Checks for proper Exit Conditions sections and likely misuse patterns

set -euo pipefail

VIOLATIONS=0
WARNINGS=0

echo "🔍 CONTINUATION FLAG VALIDATION"
echo "================================"
echo ""

# Color codes for output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

cd "$(dirname "$0")/.." || exit 1

# Function to report violations
report_violation() {
    echo -e "${RED}🔴 VIOLATION: $1${NC}"
    VIOLATIONS=$((VIOLATIONS + 1))
}

# Function to report warnings
report_warning() {
    echo -e "${YELLOW}⚠️  WARNING: $1${NC}"
    WARNINGS=$((WARNINGS + 1))
}

# Function to report success
report_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

echo "📊 Checking orchestrator state rules..."
echo ""

# Check each orchestrator state
for state_file in agent-states/software-factory/orchestrator/*/rules.md; do
    state=$(basename "$(dirname "$state_file")")

    echo "=== Checking state: $state ==="

    # Check if state has Exit Conditions section
    if ! grep -q "Exit Conditions and Continuation Flag" "$state_file" 2>/dev/null; then
        # Some states might not need it (INIT, PROJECT_DONE, ERROR_RECOVERY)
        case "$state" in
            INIT|PROJECT_DONE|ERROR_RECOVERY|CREATE_NEXT_INFRASTRUCTURE|SETUP_*_INFRASTRUCTURE|VALIDATE_INFRASTRUCTURE|VERIFY_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE|PR_PLAN_CREATION|WAVE_COMPLETE|WAVE_START|START_PHASE_ITERATION|ANALYZE_IMPLEMENTATION_PARALLELIZATION)
                report_warning "$state: No Exit Conditions section (may be acceptable for this state type)"
                ;;
            *)
                report_violation "$state: Missing 'Exit Conditions and Continuation Flag' section"
                ;;
        esac
    else
        report_success "$state: Has Exit Conditions section"
    fi

    # Check for suspicious FALSE recommendations in SPAWN/MONITOR/WAITING states
    if echo "$state" | grep -qE "SPAWN|MONITOR|WAITING"; then
        # Look for lines that recommend FALSE without proper context
        if grep -n "CONTINUE-SOFTWARE-FACTORY=FALSE" "$state_file" 2>/dev/null | \
           grep -v "ONLY\|exceptional\|ERROR\|corruption\|Exceptional Conditions ONLY\|Use FALSE ONLY"; then
            report_violation "$state: Likely incorrect FALSE usage (not marked as exceptional/error only)"
        fi

        # Check that TRUE is mentioned before FALSE
        true_line=$(grep -n "CONTINUE-SOFTWARE-FACTORY=TRUE" "$state_file" 2>/dev/null | head -1 | cut -d: -f1)
        false_line=$(grep -n "CONTINUE-SOFTWARE-FACTORY=FALSE" "$state_file" 2>/dev/null | head -1 | cut -d: -f1)

        if [ -n "$false_line" ] && [ -n "$true_line" ]; then
            if [ "$false_line" -lt "$true_line" ]; then
                report_warning "$state: FALSE documented before TRUE (should emphasize TRUE first)"
            fi
        fi
    fi

    # Check for proper context about R322 distinction
    if grep -q "CONTINUE-SOFTWARE-FACTORY" "$state_file" 2>/dev/null; then
        if ! grep -qi "R322.*stop\|stop.*R322\|Critical Distinction" "$state_file" 2>/dev/null; then
            report_warning "$state: Flag documented but no R322/stop distinction explanation"
        fi
    fi

    # Check automation flag section has proper flag
    if grep -q "## Automation Flag" "$state_file" 2>/dev/null; then
        # Get the flag value in the automation flag section
        flag_value=$(sed -n '/## Automation Flag/,/^##/p' "$state_file" | \
                     grep -o "CONTINUE-SOFTWARE-FACTORY=[A-Z]*" | head -1)

        if [ -n "$flag_value" ]; then
            if echo "$flag_value" | grep -q "FALSE"; then
                # FALSE should only be in error/setup states
                case "$state" in
                    ERROR_RECOVERY|SETUP_*_INFRASTRUCTURE|CREATE_NEXT_INFRASTRUCTURE)
                        # These might legitimately use FALSE
                        ;;
                    *)
                        report_violation "$state: Automation Flag section shows FALSE for likely normal operations"
                        ;;
                esac
            elif echo "$flag_value" | grep -q "TRUE"; then
                report_success "$state: Automation Flag correctly set to TRUE"
            fi
        fi
    fi

    echo ""
done

echo ""
echo "================================"
echo "📊 VALIDATION SUMMARY"
echo "================================"
echo -e "${RED}🔴 Violations: $VIOLATIONS${NC}"
echo -e "${YELLOW}⚠️  Warnings: $WARNINGS${NC}"

if [ $VIOLATIONS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ ALL CHECKS PASSED - No issues found${NC}"
    exit 0
elif [ $VIOLATIONS -eq 0 ]; then
    echo -e "${YELLOW}⚠️  WARNINGS ONLY - Review recommended${NC}"
    exit 0
else
    echo -e "${RED}❌ VIOLATIONS FOUND - Must fix before commit${NC}"
    echo ""
    echo "Common fixes:"
    echo "1. Add Exit Conditions section to SPAWN/MONITOR/WAITING states"
    echo "2. Ensure TRUE is emphasized as default"
    echo "3. Mark FALSE conditions as 'ONLY' exceptional cases"
    echo "4. Include R322 distinction explanation"
    exit 1
fi
