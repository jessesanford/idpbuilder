#!/bin/bash
# Validates that orchestrator states default to TRUE properly

set -euo pipefail

echo "🔍 Validating Continuation Flag Defaults"
echo "=========================================="
echo ""

VIOLATIONS=0
WARNINGS=0

# Check all orchestrator states (both locations)
STATE_DIRS=(
    "/home/vscode/software-factory-template/agent-states/orchestrator"
    "/home/vscode/software-factory-template/agent-states/software-factory/orchestrator"
)

for STATE_BASE_DIR in "${STATE_DIRS[@]}"; do
    if [[ ! -d "$STATE_BASE_DIR" ]]; then
        continue
    fi

    echo "📂 Checking states in: $STATE_BASE_DIR"
    echo ""

    for state_file in "$STATE_BASE_DIR"/*/rules.md; do
        if [[ ! -f "$state_file" ]]; then
            continue
        fi

        state=$(basename $(dirname "$state_file"))

        # Check for Exit Conditions section
        if ! grep -q "Exit Conditions and Continuation Flag" "$state_file"; then
            echo "⚠️  $state: Missing Exit Conditions section"
            WARNINGS=$((WARNINGS + 1))
        fi

        # Check for reference to master guide
        if ! grep -q "R405-CONTINUATION-FLAG-MASTER-GUIDE" "$state_file"; then
            echo "⚠️  $state: No reference to R405-CONTINUATION-FLAG-MASTER-GUIDE"
            WARNINGS=$((WARNINGS + 1))
        fi

        # Check for suspicious "awaiting user" language that suggests FALSE for normal ops
        if grep -i "await.*user\|user.*review.*required\|stop.*user" "$state_file" | grep -v "exceptional\|ERROR\|ONLY\|unrecoverable" >/dev/null 2>&1; then
            echo "🔴 $state: Suspicious language suggesting FALSE for normal ops"
            VIOLATIONS=$((VIOLATIONS + 1))
        fi

        # Check that TRUE section comes before FALSE section (proper emphasis)
        if grep -q "CONTINUE-SOFTWARE-FACTORY=TRUE" "$state_file"; then
            TRUE_LINE=$(grep -n "CONTINUE-SOFTWARE-FACTORY=TRUE" "$state_file" | head -1 | cut -d: -f1)
            if grep -q "CONTINUE-SOFTWARE-FACTORY=FALSE" "$state_file"; then
                FALSE_LINE=$(grep -n "CONTINUE-SOFTWARE-FACTORY=FALSE" "$state_file" | head -1 | cut -d: -f1)
                if [ "$FALSE_LINE" -lt "$TRUE_LINE" ]; then
                    echo "🔴 $state: FALSE section before TRUE section (wrong emphasis)"
                    VIOLATIONS=$((VIOLATIONS + 1))
                fi
            fi
        fi

        # Check for "DEFAULT TO TRUE" language
        if ! grep -q "DEFAULT TO TRUE\|default.*TRUE\|99.9%.*TRUE" "$state_file" 2>/dev/null; then
            echo "⚠️  $state: Missing DEFAULT TO TRUE emphasis"
            WARNINGS=$((WARNINGS + 1))
        fi

        # Check for R322 vs flag distinction
        if grep -q "R322" "$state_file" && ! grep -q "R322.*stop.*≠.*FALSE\|R322.*not mean.*FALSE\|checkpoint.*≠.*FALSE" "$state_file" 2>/dev/null; then
            echo "⚠️  $state: Has R322 but doesn't clarify stop ≠ FALSE"
            WARNINGS=$((WARNINGS + 1))
        fi
    done
done

echo ""
echo "=========================================="
echo "📊 VALIDATION SUMMARY"
echo "=========================================="

if [ $VIOLATIONS -gt 0 ]; then
    echo "🔴 Found $VIOLATIONS critical issues requiring immediate fix"
fi

if [ $WARNINGS -gt 0 ]; then
    echo "⚠️  Found $WARNINGS warnings (recommended improvements)"
fi

if [ $VIOLATIONS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo "✅ All states have proper continuation flag guidance"
    exit 0
else
    echo ""
    echo "RECOMMENDATIONS:"
    echo "1. Add Exit Conditions section to states missing it"
    echo "2. Reference R405-CONTINUATION-FLAG-MASTER-GUIDE.md in all states"
    echo "3. Emphasize DEFAULT TO TRUE in all states"
    echo "4. Clarify R322 stop ≠ FALSE flag in states mentioning R322"
    echo "5. Remove language suggesting FALSE for normal operations"

    if [ $VIOLATIONS -gt 0 ]; then
        exit 1
    else
        exit 0
    fi
fi
