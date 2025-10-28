#!/bin/bash
# Validate all state rules exist (Phase 8 Item #2)

set -euo pipefail

STATE_MACHINE="state-machines/software-factory-3.0-state-machine.json"
STATE_DIR="agent-states/software-factory/orchestrator"

echo "🔍 SF 3.0 State Rules Existence Validation"
echo "==========================================="
echo ""

MISSING=0
TOTAL=0

while read -r state; do
    TOTAL=$((TOTAL + 1))

    if [ ! -f "$STATE_DIR/$state/rules.md" ]; then
        echo "❌ MISSING: $state/rules.md"
        MISSING=$((MISSING + 1))
    fi
done < <(jq -r '.states | keys[]' "$STATE_MACHINE")

echo ""
echo "==========================================="
echo "Total states: $TOTAL"
echo "Missing rules.md files: $MISSING"
echo ""

if [ "$MISSING" -eq 0 ]; then
    echo "✅ ALL CHECKS PASSED - Every state has rules.md"
    echo "==========================================="
    exit 0
else
    echo "❌ VALIDATION FAILED - $MISSING states missing rules.md"
    echo "==========================================="

    # List all missing for debugging
    echo ""
    echo "Missing rules.md files:"
    while read -r state; do
        if [ ! -f "$STATE_DIR/$state/rules.md" ]; then
            echo "  - $state"
        fi
    done < <(jq -r '.states | keys[]' "$STATE_MACHINE")

    exit 1
fi
