#!/bin/bash
# Validate state machine completeness (Phase 8 Item #1)

set -euo pipefail

STATE_MACHINE="state-machines/software-factory-3.0-state-machine.json"
STATE_DIR="agent-states/software-factory/orchestrator"

echo "🔍 SF 3.0 State Machine Completeness Validation"
echo "================================================"
echo ""

# Check 1: State count in expected range
echo "Check 1: State count validation"
STATE_COUNT=$(jq '.states | keys | length' "$STATE_MACHINE")
echo "  Current state count: $STATE_COUNT"

if [ "$STATE_COUNT" -ge 90 ] && [ "$STATE_COUNT" -le 100 ]; then
    echo "  ✅ PASS: State count in expected range (90-100)"
else
    echo "  ❌ FAIL: State count out of range (expected 90-100, got $STATE_COUNT)"
    exit 1
fi

# Check 2: No orphaned states (directories without state machine entry)
echo ""
echo "Check 2: Orphaned state directory detection"
ORPHANS=0
for dir in "$STATE_DIR"/*/; do
    state=$(basename "$dir")
    # Skip special directories
    if [[ "$state" == "DEPRECATED" ]] || [[ "$state" == "ARCHIVED" ]] || [[ "$state" == "_"* ]]; then
        continue
    fi

    # Check if state exists in state machine
    if ! jq -e ".states | has(\"$state\")" "$STATE_MACHINE" > /dev/null 2>&1; then
        echo "  ⚠️  Orphaned directory: $state"
        ORPHANS=$((ORPHANS + 1))
    fi
done

if [ "$ORPHANS" -eq 0 ]; then
    echo "  ✅ PASS: No orphaned state directories found"
else
    echo "  ❌ FAIL: Found $ORPHANS orphaned state directories"
    exit 1
fi

# Check 3: No orphaned state machine entries (states without directories)
echo ""
echo "Check 3: Orphaned state machine entry detection"
MISSING_DIRS=0
while read -r state; do
    if [ ! -d "$STATE_DIR/$state" ]; then
        echo "  ⚠️  Missing directory for state: $state"
        MISSING_DIRS=$((MISSING_DIRS + 1))
    fi
done < <(jq -r '.states | keys[]' "$STATE_MACHINE")

if [ "$MISSING_DIRS" -eq 0 ]; then
    echo "  ✅ PASS: All state machine entries have directories"
else
    echo "  ❌ FAIL: Found $MISSING_DIRS state machine entries without directories"
    exit 1
fi

# Check 4: Valid JSON syntax
echo ""
echo "Check 4: JSON syntax validation"
if jq empty "$STATE_MACHINE" 2>/dev/null; then
    echo "  ✅ PASS: Valid JSON syntax"
else
    echo "  ❌ FAIL: Invalid JSON syntax"
    exit 1
fi

echo ""
echo "================================================"
echo "✅ ALL CHECKS PASSED - State machine is complete"
echo "================================================"
exit 0
