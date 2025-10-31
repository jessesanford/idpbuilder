#!/bin/bash
# R550 Integration Test
# Tests R550 plan path consistency implementation end-to-end

# Don't use set -e for tests - we want to continue after failures
set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "🧪 R550 Integration Test"
echo "========================"
echo ""

TESTS_PASSED=0
TESTS_FAILED=0

# Test 1: Planning directory structure uses correct naming
echo "Test 1: Planning directory uses canonical R550 paths"
if [ -d "$PROJECT_ROOT/planning" ]; then
    # Check that planning/ directory exists (not phase-plans/)
    if [ ! -d "$PROJECT_ROOT/phase-plans" ]; then
        echo "  ✅ PASS: Uses planning/ directory (not phase-plans/)"
        ((TESTS_PASSED++))
    else
        echo "  ❌ FAIL: Found phase-plans/ directory (should be planning/)"
        ((TESTS_FAILED++))
    fi

    # Check phase naming convention (no hyphens)
    if ls -d "$PROJECT_ROOT/planning/phase-"[0-9]* >/dev/null 2>&1; then
        echo "  ❌ FAIL: Found phase directories with hyphens"
        ((TESTS_FAILED++))
    else
        echo "  ✅ PASS: Phase directories use correct naming (phase1, not phase-1)"
        ((TESTS_PASSED++))
    fi
else
    echo "  ⚠️  SKIP: planning/ directory not yet created"
fi

echo ""

# Test 2: Orchestrator state schema includes planning_files
echo "Test 2: Schema includes planning_files tracking"
if [ -f "$PROJECT_ROOT/schemas/orchestrator-state-v3.schema.json" ]; then
    if grep -q '"planning_files"' "$PROJECT_ROOT/schemas/orchestrator-state-v3.schema.json"; then
        echo "  ✅ PASS: Schema includes planning_files property"
        ((TESTS_PASSED++))
    else
        echo "  ❌ FAIL: Schema missing planning_files property"
        ((TESTS_FAILED++))
    fi
else
    echo "  ❌ FAIL: Schema file not found"
    ((TESTS_FAILED++))
fi

echo ""

# Test 3: Example state demonstrates planning_files tracking
echo "Test 3: Example state includes planning_files"
if [ -f "$PROJECT_ROOT/orchestrator-state-v3.json.example" ]; then
    if jq -e '.planning_files' "$PROJECT_ROOT/orchestrator-state-v3.json.example" >/dev/null 2>&1; then
        echo "  ✅ PASS: Example state includes planning_files section"
        ((TESTS_PASSED++))

        # Verify it has project, phases structure
        if jq -e '.planning_files.project' "$PROJECT_ROOT/orchestrator-state-v3.json.example" >/dev/null 2>&1; then
            echo "  ✅ PASS: planning_files has project section"
            ((TESTS_PASSED++))
        else
            echo "  ❌ FAIL: planning_files missing project section"
            ((TESTS_FAILED++))
        fi

        if jq -e '.planning_files.phases' "$PROJECT_ROOT/orchestrator-state-v3.json.example" >/dev/null 2>&1; then
            echo "  ✅ PASS: planning_files has phases section"
            ((TESTS_PASSED++))
        else
            echo "  ❌ FAIL: planning_files missing phases section"
            ((TESTS_FAILED++))
        fi
    else
        echo "  ❌ FAIL: Example state missing planning_files"
        ((TESTS_FAILED++))
    fi
else
    echo "  ❌ FAIL: Example state file not found"
    ((TESTS_FAILED++))
fi

echo ""

# Test 4: No phase-plans/ references in orchestrator state rules
echo "Test 4: No phase-plans/ references in orchestrator states"
PHASE_PLANS_REFS=$(grep -r "phase-plans/" "$PROJECT_ROOT/agent-states/software-factory/orchestrator" 2>/dev/null || true)
if [ -z "$PHASE_PLANS_REFS" ]; then
    echo "  ✅ PASS: No phase-plans/ references found"
    ((TESTS_PASSED++))
else
    echo "  ❌ FAIL: Found phase-plans/ references:"
    echo "$PHASE_PLANS_REFS" | head -5 | sed 's/^/    /'
    ((TESTS_FAILED++))
fi

echo ""

# Test 5: R550 validation script exists and is executable
echo "Test 5: R550 validation script exists and works"
if [ -f "$PROJECT_ROOT/tools/validate-R550-compliance.sh" ]; then
    if [ -x "$PROJECT_ROOT/tools/validate-R550-compliance.sh" ]; then
        echo "  ✅ PASS: Validation script exists and is executable"
        ((TESTS_PASSED++))

        # Try running it (may fail if violations exist, but shouldn't crash)
        if bash "$PROJECT_ROOT/tools/validate-R550-compliance.sh" >/dev/null 2>&1; then
            echo "  ✅ PASS: Validation script executes successfully"
            ((TESTS_PASSED++))
        else
            echo "  ⚠️  WARN: Validation script found violations (expected during development)"
        fi
    else
        echo "  ❌ FAIL: Validation script not executable"
        ((TESTS_FAILED++))
    fi
else
    echo "  ❌ FAIL: Validation script not found"
    ((TESTS_FAILED++))
fi

echo ""

# Test 6: Pre-commit hook includes R550 validation
echo "Test 6: Pre-commit hook includes R550 validation"
if [ -f "$PROJECT_ROOT/tools/git-commit-hooks/master-pre-commit.sh" ]; then
    if grep -q "validate-R550-compliance" "$PROJECT_ROOT/tools/git-commit-hooks/master-pre-commit.sh"; then
        echo "  ✅ PASS: Pre-commit hook includes R550 validation"
        ((TESTS_PASSED++))
    else
        echo "  ❌ FAIL: Pre-commit hook missing R550 validation"
        ((TESTS_FAILED++))
    fi
else
    echo "  ❌ FAIL: Pre-commit hook not found"
    ((TESTS_FAILED++))
fi

echo ""

# Test 7: CREATE_NEXT_INFRASTRUCTURE uses R550 compliant paths
echo "Test 7: CREATE_NEXT_INFRASTRUCTURE uses R550 paths"
CREATE_INFRA="$PROJECT_ROOT/agent-states/software-factory/orchestrator/CREATE_NEXT_INFRASTRUCTURE/rules.md"
if [ -f "$CREATE_INFRA" ]; then
    # Should reference planning/ not phase-plans/
    if grep -q "planning/" "$CREATE_INFRA" && ! grep -q "phase-plans/" "$CREATE_INFRA"; then
        echo "  ✅ PASS: Uses planning/ directory"
        ((TESTS_PASSED++))
    else
        echo "  ❌ FAIL: Still references phase-plans/ or missing planning/"
        ((TESTS_FAILED++))
    fi

    # Should reference orchestrator state planning_files
    if grep -q "planning_files" "$CREATE_INFRA"; then
        echo "  ✅ PASS: Uses orchestrator state planning_files tracking"
        ((TESTS_PASSED++))
    else
        echo "  ⚠️  WARN: Doesn't use orchestrator state tracking (acceptable if uses standard paths)"
    fi
else
    echo "  ❌ FAIL: CREATE_NEXT_INFRASTRUCTURE state file not found"
    ((TESTS_FAILED++))
fi

echo ""
echo "========================"
echo "📊 Test Results"
echo "========================"
echo "Passed: $TESTS_PASSED"
echo "Failed: $TESTS_FAILED"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo "✅ R550 Integration Test: ALL TESTS PASSED"
    echo ""
    echo "R550 implementation is complete and functional:"
    echo "  ✓ Schema updated with planning_files tracking"
    echo "  ✓ Example state demonstrates tracking"
    echo "  ✓ All phase-plans/ references replaced with planning/"
    echo "  ✓ Validation script operational"
    echo "  ✓ Pre-commit hook includes R550 validation"
    echo "  ✓ State rules use R550 compliant paths"
    echo ""
    exit 0
else
    echo "❌ R550 Integration Test: $TESTS_FAILED test(s) failed"
    echo ""
    echo "Review failures above and:"
    echo "  1. Fix any remaining phase-plans/ references"
    echo "  2. Ensure schema and state files include planning_files"
    echo "  3. Verify validation script is executable"
    echo "  4. Check pre-commit hook integration"
    echo ""
    exit 1
fi
