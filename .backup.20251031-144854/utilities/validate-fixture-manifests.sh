#!/bin/bash
# Validates all fixture manifests for consistency and correctness
# - Checks JSON syntax
# - Verifies template files exist
# - Validates variable references
# - Ensures consistency
# Created: 2025-10-18

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

PASS_COUNT=0
FAIL_COUNT=0

echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║           Fixture Manifest Validator                        ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo ""

# Find all manifest files
MANIFESTS=$(find "${PROJECT_ROOT}/tests/fixtures" -name "manifest.json" | sort)

if [ -z "$MANIFESTS" ]; then
    echo "❌ No manifest.json files found in tests/fixtures/"
    exit 1
fi

MANIFEST_COUNT=$(echo "$MANIFESTS" | wc -l)
echo "Found $MANIFEST_COUNT manifest(s) to validate"
echo ""

# Validate each manifest
for manifest in $MANIFESTS; do
    test_dir=$(dirname "$manifest")
    test_name=$(basename "$test_dir")

    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📋 Validating: $test_name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    errors=0
    warnings=0

    # 1. Check JSON syntax
    if ! jq empty "$manifest" 2>/dev/null; then
        echo "  ❌ Invalid JSON syntax"
        ((errors++))
    else
        echo "  ✅ Valid JSON syntax"
    fi

    # 2. Check required fields
    required_fields=("test_id" "test_name" "fixtures" "variable_definitions")
    for field in "${required_fields[@]}"; do
        if ! jq -e ".$field" "$manifest" >/dev/null 2>&1; then
            echo "  ❌ Missing required field: $field"
            ((errors++))
        fi
    done

    # 3. Validate each fixture
    fixture_count=$(jq '.fixtures | length' "$manifest" 2>/dev/null || echo "0")
    echo "  Fixtures: $fixture_count"

    if [ "$fixture_count" -eq 0 ]; then
        # Check if this is intentional (dynamic fixtures only)
        notes=$(jq -r '.notes // ""' "$manifest")
        if echo "$notes" | grep -qi "dynamic\|init_orchestrator_state"; then
            echo "    ℹ️  No fixtures (dynamic creation only - see notes)"
        else
            echo "    ⚠️  No fixtures defined (is this intentional?)"
            ((warnings++))
        fi
    else
        i=0
        while [ $i -lt $fixture_count ]; do
            template_path=$(jq -r ".fixtures[$i].template_path" "$manifest")
            dest_path=$(jq -r ".fixtures[$i].destination_path" "$manifest")
            needs_subst=$(jq -r ".fixtures[$i].requires_substitution" "$manifest")

            # Check template file exists
            if [ ! -f "${PROJECT_ROOT}/${template_path}" ]; then
                echo "    ❌ Template not found: $template_path"
                ((errors++))
            else
                echo "    ✅ $dest_path"
            fi

            # Check .template suffix consistency
            if [ "$needs_subst" = "true" ] && [[ ! "$template_path" =~ \.template$ ]]; then
                echo "       ⚠️  Warning: requires_substitution=true but no .template suffix: $template_path"
                ((warnings++))
            fi

            if [ "$needs_subst" = "false" ] && [[ "$template_path" =~ \.template$ ]]; then
                echo "       ⚠️  Warning: requires_substitution=false but has .template suffix: $template_path"
                ((warnings++))
            fi

            # Check for required_variables field existence
            if ! jq -e ".fixtures[$i].required_variables" "$manifest" >/dev/null 2>&1; then
                echo "       ⚠️  Warning: Missing required_variables field (should be array, even if empty)"
                ((warnings++))
            fi

            i=$((i+1))
        done
    fi

    # 4. Check for unreplaced old template format in source files
    if [ "$fixture_count" -gt 0 ]; then
        i=0
        while [ $i -lt $fixture_count ]; do
            template_path=$(jq -r ".fixtures[$i].template_path" "$manifest")
            if [ -f "${PROJECT_ROOT}/${template_path}" ]; then
                # Check for old {{VAR}} format
                if grep -qE '\{\{[A-Z_]+\}\}' "${PROJECT_ROOT}/${template_path}" 2>/dev/null; then
                    echo "    ❌ Old template format {{VAR}} found in: $(basename "$template_path")"
                    echo "       Should use \${VAR} format for envsubst"
                    ((errors++))
                fi
            fi
            i=$((i+1))
        done
    fi

    # 5. Check test_file field matches actual file
    test_file=$(jq -r '.test_file // ""' "$manifest")
    if [ -n "$test_file" ]; then
        if [ ! -f "${PROJECT_ROOT}/tests/$test_file" ]; then
            echo "  ⚠️  Warning: Referenced test file not found: tests/$test_file"
            ((warnings++))
        fi
    fi

    # 6. Validate test_id matches directory name pattern
    test_id=$(jq -r '.test_id' "$manifest")
    expected_pattern="test-${test_id}-"
    if [[ ! "$test_name" =~ ^${expected_pattern} ]]; then
        echo "  ⚠️  Warning: test_id '$test_id' doesn't match directory name pattern '$test_name'"
        ((warnings++))
    fi

    # Summary for this manifest
    echo ""
    if [ $errors -eq 0 ]; then
        if [ $warnings -eq 0 ]; then
            echo "  ✅ PASS: $test_name (no errors or warnings)"
        else
            echo "  ✅ PASS: $test_name ($warnings warning(s))"
        fi
        ((PASS_COUNT++))
    else
        echo "  ❌ FAIL: $test_name ($errors error(s), $warnings warning(s))"
        ((FAIL_COUNT++))
    fi
    echo ""
done

# Final summary
echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║                    VALIDATION SUMMARY                        ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo ""
echo "  Total Manifests: $MANIFEST_COUNT"
echo "  Passed: $PASS_COUNT"
echo "  Failed: $FAIL_COUNT"
echo ""

if [ $FAIL_COUNT -eq 0 ]; then
    echo "✅ All manifests validated successfully!"
    exit 0
else
    echo "❌ Some manifests have errors"
    exit 1
fi
