#!/bin/bash
# R550 Plan Path Consistency Validation Script
# Validates that all planning paths follow R550 canonical naming and are properly tracked

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "🔍 R550 Plan Path Consistency Validation"
echo "=========================================="

VIOLATIONS=0

# Check 1: No phase-plans/ references in state rules
echo ""
echo "Check 1: No legacy phase-plans/ references in orchestrator states..."
if grep -r "phase-plans/" "$PROJECT_ROOT/agent-states/software-factory/orchestrator" 2>/dev/null; then
    echo "❌ VIOLATION: Found phase-plans/ references (should be planning/)"
    echo "   Fix: Replace all 'phase-plans/' with 'planning/'"
    ((VIOLATIONS++))
else
    echo "✅ PASS: No phase-plans/ references found in orchestrator states"
fi

# Check 2: All planning files use canonical naming (no timestamps in planning/)
echo ""
echo "Check 2: Canonical naming in planning/ directory..."
if [ -d "$PROJECT_ROOT/planning" ]; then
    # Check for timestamp patterns in planning/ filenames
    WRONG_NAMES=$(find "$PROJECT_ROOT/planning" -type f -name "*--[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]*" 2>/dev/null || true)
    if [ -n "$WRONG_NAMES" ]; then
        echo "❌ VIOLATION: Found timestamps in planning/ filenames (R550 violation):"
        echo "$WRONG_NAMES" | sed 's/^/   /'
        echo "   Fix: Remove timestamps from planning/ directory files"
        ((VIOLATIONS++))
    else
        echo "✅ PASS: No timestamps in planning/ filenames"
    fi

    # Check for generic PLAN names without type specification
    GENERIC_PLANS=$(find "$PROJECT_ROOT/planning" -type f -name "*PLAN.md" ! -name "*-PLAN.md" 2>/dev/null || true)
    if [ -n "$GENERIC_PLANS" ]; then
        echo "❌ VIOLATION: Found generic PLAN names without type (R550 violation):"
        echo "$GENERIC_PLANS" | sed 's/^/   /'
        echo "   Fix: Use explicit types: IMPLEMENTATION-PLAN, ARCHITECTURE-PLAN, TEST-PLAN"
        ((VIOLATIONS++))
    else
        echo "✅ PASS: All plan files use explicit type names"
    fi
else
    echo "⚠️  SKIP: planning/ directory does not exist yet"
fi

# Check 3: Schema includes planning_files property
echo ""
echo "Check 3: Schema has planning_files tracking..."
if [ -f "$PROJECT_ROOT/schemas/orchestrator-state-v3.schema.json" ]; then
    if ! grep -q "planning_files" "$PROJECT_ROOT/schemas/orchestrator-state-v3.schema.json" 2>/dev/null; then
        echo "❌ VIOLATION: Schema missing planning_files property"
        echo "   Fix: Add planning_files property to orchestrator-state-v3.schema.json"
        ((VIOLATIONS++))
    else
        echo "✅ PASS: Schema includes planning_files tracking"
    fi
else
    echo "⚠️  SKIP: Schema file not found"
fi

# Check 4: Example state includes planning_files
echo ""
echo "Check 4: Example state has planning_files..."
if [ -f "$PROJECT_ROOT/orchestrator-state-v3.json.example" ]; then
    if ! jq -e '.planning_files' "$PROJECT_ROOT/orchestrator-state-v3.json.example" >/dev/null 2>&1; then
        echo "❌ VIOLATION: Example state missing planning_files"
        echo "   Fix: Add planning_files section to orchestrator-state-v3.json.example"
        ((VIOLATIONS++))
    else
        echo "✅ PASS: Example state includes planning_files"
    fi
else
    echo "⚠️  SKIP: Example state file not found"
fi

# Check 5: No filesystem searching commands (ls -t, find with -t) in orchestrator states
echo ""
echo "Check 5: No filesystem searching for plans in orchestrator states..."
SEARCH_VIOLATIONS=$(grep -rn "ls -t.*planning/\|ls -t.*PLAN" "$PROJECT_ROOT/agent-states/software-factory/orchestrator" 2>/dev/null || true)
if [ -n "$SEARCH_VIOLATIONS" ]; then
    echo "❌ VIOLATION: Found filesystem searching for plans (violates R340/R550):"
    echo "$SEARCH_VIOLATIONS" | sed 's/^/   /'
    echo "   Fix: Use orchestrator state planning_files tracking instead"
    ((VIOLATIONS++))
else
    echo "✅ PASS: No filesystem searching detected"
fi

# Check 6: Verify planning/ directory structure if it exists
echo ""
echo "Check 6: Planning directory structure compliance..."
if [ -d "$PROJECT_ROOT/planning" ]; then
    # Check for phase directories with hyphens (should be phase1, not phase-1)
    HYPHEN_PHASES=$(find "$PROJECT_ROOT/planning" -maxdepth 1 -type d -name "phase-[0-9]*" 2>/dev/null || true)
    if [ -n "$HYPHEN_PHASES" ]; then
        echo "❌ VIOLATION: Found phase directories with hyphens:"
        echo "$HYPHEN_PHASES" | sed 's/^/   /'
        echo "   Fix: Rename to phase1, phase2, etc. (no hyphen before number)"
        ((VIOLATIONS++))
    else
        echo "✅ PASS: Phase directory naming correct"
    fi

    # Check for wave directories with hyphens (should be wave1, not wave-1)
    HYPHEN_WAVES=$(find "$PROJECT_ROOT/planning" -type d -name "wave-[0-9]*" 2>/dev/null || true)
    if [ -n "$HYPHEN_WAVES" ]; then
        echo "❌ VIOLATION: Found wave directories with hyphens:"
        echo "$HYPHEN_WAVES" | sed 's/^/   /'
        echo "   Fix: Rename to wave1, wave2, etc. (no hyphen before number)"
        ((VIOLATIONS++))
    else
        echo "✅ PASS: Wave directory naming correct"
    fi
else
    echo "⚠️  SKIP: planning/ directory does not exist yet"
fi

# Summary
echo ""
echo "=========================================="
if [ $VIOLATIONS -eq 0 ]; then
    echo "✅ R550 VALIDATION PASSED - All checks successful"
    exit 0
else
    echo "❌ R550 VALIDATION FAILED: $VIOLATIONS violation(s) detected"
    echo ""
    echo "R550 Summary:"
    echo "  - All planning files must be in planning/ directory (NOT phase-plans/)"
    echo "  - No timestamps in planning/ filenames"
    echo "  - Explicit plan types (IMPLEMENTATION-PLAN, ARCHITECTURE-PLAN, TEST-PLAN)"
    echo "  - All plan paths tracked in orchestrator-state-v3.json planning_files"
    echo "  - No filesystem searching with ls -t or find"
    echo ""
    echo "See: rule-library/R550-plan-path-consistency-and-discovery.md"
    exit 550
fi
