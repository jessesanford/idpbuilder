#!/bin/bash
# Validate New Test Script
# Purpose: Ensure new runtime tests follow all safeguards and best practices
# Usage: bash tools/validate-new-test.sh tests/runtime-test-XX-description.sh

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

TEST_SCRIPT="$1"

# Usage check
if [ -z "$TEST_SCRIPT" ]; then
    echo "Usage: $0 <test-script.sh>"
    echo ""
    echo "Example:"
    echo "  $0 tests/runtime-test-07-new-test.sh"
    exit 2
fi

# File exists check
if [ ! -f "$TEST_SCRIPT" ]; then
    echo -e "${RED}ERROR: Test script not found: $TEST_SCRIPT${NC}"
    exit 2
fi

# Extract test number
TEST_NUM=$(basename "$TEST_SCRIPT" | grep -oP 'runtime-test-\K\d+' || echo "")

echo ""
echo "========================================="
echo "  Test Validation: $(basename "$TEST_SCRIPT")"
echo "========================================="
echo ""

VALIDATION_PASSED=0

# Check 1: Sources runtime-test-framework.sh
echo -n "Check 1/10: Sources runtime-test-framework.sh ... "
if grep -q "source.*runtime-test-framework.sh" "$TEST_SCRIPT"; then
    echo -e "${GREEN}PASS${NC}"
else
    echo -e "${RED}FAIL${NC}"
    echo "  ❌ Test must source runtime-test-framework.sh"
    echo "  Add: source \"\$(dirname \"\${BASH_SOURCE[0]}\")/runtime-test-framework.sh\""
    VALIDATION_PASSED=1
fi

# Check 2: No manual fixture creation (cat/echo to JSON)
echo -n "Check 2/10: No manual fixture creation ... "
if grep -qE "cat\s+(>|<<).*\.json|echo.*>.*\.json" "$TEST_SCRIPT"; then
    echo -e "${RED}FAIL${NC}"
    echo "  ❌ Test contains manual fixture creation (cat/echo to .json file)"
    echo "  Use template system instead: generate_fixture_from_template()"
    VALIDATION_PASSED=1
else
    echo -e "${GREEN}PASS${NC}"
fi

# Check 3: No jq heredocs for fixture creation
echo -n "Check 3/10: No jq heredocs ... "
if grep -qE "cat\s+<<.*EOF.*jq|jq.*<<.*EOF" "$TEST_SCRIPT"; then
    echo -e "${RED}FAIL${NC}"
    echo "  ❌ Test contains jq heredoc (manual fixture creation)"
    echo "  Use template system instead: generate_fixture_from_template()"
    VALIDATION_PASSED=1
else
    echo -e "${GREEN}PASS${NC}"
fi

# Check 4: Uses template system or framework fixture installation
echo -n "Check 4/10: Uses template system ... "
if grep -qE "install_fixture_files|generate_fixture_from_template|prepare_fixtures" "$TEST_SCRIPT"; then
    echo -e "${GREEN}PASS${NC}"
else
    echo -e "${YELLOW}WARNING${NC}"
    echo "  ⚠️  Test may not use template system"
    echo "  Ensure test uses install_fixture_files() or generate_fixture_from_template()"
fi

# Check 5: Has test_setup() function
echo -n "Check 5/10: Has test_setup() function ... "
if grep -q "test_setup()" "$TEST_SCRIPT"; then
    echo -e "${GREEN}PASS${NC}"
else
    echo -e "${RED}FAIL${NC}"
    echo "  ❌ Test missing test_setup() function"
    echo "  Required for test initialization"
    VALIDATION_PASSED=1
fi

# Check 6: Has test_validate() function
echo -n "Check 6/10: Has test_validate() function ... "
if grep -q "test_validate()" "$TEST_SCRIPT"; then
    echo -e "${GREEN}PASS${NC}"
else
    echo -e "${RED}FAIL${NC}"
    echo "  ❌ Test missing test_validate() function"
    echo "  Required for test validation logic"
    VALIDATION_PASSED=1
fi

# Check 7: Calls run_test()
echo -n "Check 7/10: Calls run_test() ... "
if grep -q "run_test" "$TEST_SCRIPT"; then
    echo -e "${GREEN}PASS${NC}"
else
    echo -e "${RED}FAIL${NC}"
    echo "  ❌ Test does not call run_test()"
    echo "  Required to execute test with framework support"
    VALIDATION_PASSED=1
fi

# Check 8: Defines nominal_path array (for fail-fast)
echo -n "Check 8/10: Defines nominal_path array ... "
if grep -q "nominal_path=(" "$TEST_SCRIPT"; then
    echo -e "${GREEN}PASS${NC}"
else
    echo -e "${YELLOW}WARNING${NC}"
    echo "  ⚠️  Test does not define nominal_path array"
    echo "  Recommended for fail-fast state validation"
fi

# Check 9: No hardcoded "hello-world-fullstack" values
echo -n "Check 9/10: No hardcoded project names ... "
if grep -qi "hello-world-fullstack" "$TEST_SCRIPT"; then
    echo -e "${RED}FAIL${NC}"
    echo "  ❌ Test contains hardcoded 'hello-world-fullstack'"
    echo "  Use \$PROJECT_PREFIX instead"
    VALIDATION_PASSED=1
else
    echo -e "${GREEN}PASS${NC}"
fi

# Check 10: Has documentation header
echo -n "Check 10/10: Has documentation header ... "
if head -20 "$TEST_SCRIPT" | grep -q "# Runtime Test"; then
    echo -e "${GREEN}PASS${NC}"
else
    echo -e "${YELLOW}WARNING${NC}"
    echo "  ⚠️  Test missing documentation header"
    echo "  Add header with purpose, nominal path, duration, cost"
fi

echo ""
echo "========================================="
if [ $VALIDATION_PASSED -eq 0 ]; then
    echo -e "${GREEN}✅ VALIDATION PASSED${NC}"
    echo "========================================="
    echo ""
    echo "Test is ready for execution and commit."
    echo ""
    exit 0
else
    echo -e "${RED}❌ VALIDATION FAILED${NC}"
    echo "========================================="
    echo ""
    echo "Please fix the issues above before committing."
    echo ""
    echo "Reference:"
    echo "  - Template system: tests/fixtures/templates/README.md"
    echo "  - Protection system: docs/SF-3.0-TEST-FIXTURE-PROTECTION-SYSTEM.md"
    echo "  - Framework: tests/runtime-test-framework.sh"
    echo ""
    exit 1
fi
