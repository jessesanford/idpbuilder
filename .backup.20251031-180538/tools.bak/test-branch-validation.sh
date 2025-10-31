#!/bin/bash

# Test script for branch naming validation
source /home/vscode/software-factory-template/tools/enforce-state-validation.sh

# Test the validate_branch_format function
test_branch() {
    local branch="$1"
    local effort="$2"
    local phase="$3"
    local wave="$4"
    local expected="$5"

    validate_branch_format "$branch" "$effort" "$phase" "$wave"
    local result=$?

    if [ "$result" -eq 0 ]; then
        if [ "$expected" = "PASS" ]; then
            echo "✅ PASS: $branch"
        else
            echo "❌ FAIL: $branch (should have failed)"
        fi
    else
        if [ "$expected" = "FAIL" ]; then
            echo "✅ PASS: $branch (correctly rejected)"
        else
            echo "❌ FAIL: $branch (should have passed, error code: $result)"
        fi
    fi
}

echo "========================================="
echo "BRANCH NAMING VALIDATION TESTS"
echo "========================================="
echo ""

# Test with project prefix
export PROJECT_PREFIX="idpbuilderpush"

echo "=== WITH PROJECT PREFIX: $PROJECT_PREFIX ==="
echo ""

echo "--- Standard Effort Branches ---"
test_branch "idpbuilderpush/phase1/wave1/command-tests" "command-tests" "1" "1" "PASS"
test_branch "idpbuilderpush/phase2/wave3/authentication" "authentication" "2" "3" "PASS"
test_branch "phase1/wave1/command-tests" "command-tests" "1" "1" "FAIL"  # Missing prefix
test_branch "idpbuilderpush/phase1/wave1/Command-Tests" "command-tests" "1" "1" "FAIL"  # Uppercase

echo ""
echo "--- Integration Branches ---"
test_branch "idpbuilderpush/phase1/wave1-integration" "integration" "1" "1" "PASS"
test_branch "idpbuilderpush/phase1-integration" "integration" "1" "" "PASS"
test_branch "idpbuilderpush/integration" "integration" "" "" "PASS"
test_branch "idpbuilderpush/phase1/wave1/integration" "integration" "1" "1" "FAIL"  # Wrong format
test_branch "phase1/wave1-integration" "integration" "1" "1" "FAIL"  # Missing prefix

echo ""
echo "--- Split Branches ---"
test_branch "idpbuilderpush/phase1/wave1/effort-1.1.1--split-001" "effort-1.1.1" "1" "1" "PASS"
test_branch "idpbuilderpush/phase1/wave1/effort-1.1.1--split-002" "effort-1.1.1" "1" "1" "PASS"
test_branch "idpbuilderpush/phase1/wave1/api-types--split-123" "api-types" "1" "1" "PASS"
test_branch "idpbuilderpush/phase1/wave1/effort-1.1.1-split-1" "effort-1.1.1" "1" "1" "FAIL"  # Wrong format
test_branch "idpbuilderpush/phase1/wave1/effort-1.1.1/split-001" "effort-1.1.1" "1" "1" "FAIL"  # Wrong format

echo ""
echo "--- Fix Branches ---"
test_branch "idpbuilderpush/phase1/wave1/command-tests-fix" "command-tests" "1" "1" "PASS"
test_branch "idpbuilderpush/phase2/wave2/api-gateway-fix" "api-gateway" "2" "2" "PASS"
test_branch "idpbuilderpush/phase1/wave1/command-tests--fix" "command-tests" "1" "1" "FAIL"  # Double dash

echo ""
echo "=== WITHOUT PROJECT PREFIX ==="
export PROJECT_PREFIX=""

echo ""
echo "--- Standard Effort Branches ---"
test_branch "phase1/wave1/command-tests" "command-tests" "1" "1" "PASS"
test_branch "phase2/wave3/authentication" "authentication" "2" "3" "PASS"
test_branch "feature/command-tests" "command-tests" "1" "1" "FAIL"  # Wrong format

echo ""
echo "--- Integration Branches ---"
test_branch "phase1/wave1-integration" "integration" "1" "1" "PASS"
test_branch "phase1-integration" "integration" "1" "" "PASS"
test_branch "integration" "integration" "" "" "PASS"
test_branch "phase1/wave1/integration" "integration" "1" "1" "FAIL"  # Wrong format

echo ""
echo "--- Split Branches ---"
test_branch "phase1/wave1/effort-1.1.1--split-001" "effort-1.1.1" "1" "1" "PASS"
test_branch "phase1/wave1/effort-1.1.1--split-999" "effort-1.1.1" "1" "1" "PASS"

echo ""
echo "--- Phase/Wave Mismatch Tests ---"
test_branch "phase1/wave1/command-tests" "command-tests" "2" "1" "FAIL"  # Phase mismatch
test_branch "phase1/wave1/command-tests" "command-tests" "1" "2" "FAIL"  # Wave mismatch

echo ""
echo "========================================="
echo "TEST COMPLETE"
echo "========================================="