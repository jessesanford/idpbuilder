#!/bin/bash

# Fidelity Gradient Validation Script
# Validates that planning artifacts conform to the correct fidelity level
# Usage: ./validate-fidelity-gradient.sh <file> <expected_level>
# Expected levels: pseudocode, concrete, exact

set -euo pipefail

FILE="$1"
EXPECTED_LEVEL="$2"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Validation functions
validate_pseudocode() {
    local file="$1"
    local errors=0

    echo "Validating PSEUDOCODE fidelity (Phase Architecture)..."

    # Should have conceptual patterns
    if ! grep -qi "pattern\|concept\|pseudo\|high-level" "$file"; then
        echo -e "${RED}❌ FAILED: No pseudocode/conceptual language found${NC}"
        ((errors++))
    fi

    # Should NOT have real function definitions
    if grep -q "def .*(" "$file"; then
        echo -e "${RED}❌ FAILED: Contains real function definitions (should be pseudocode only)${NC}"
        ((errors++))
    fi

    # Should NOT have class implementations
    if grep -q "class .*:$" "$file"; then
        echo -e "${RED}❌ FAILED: Contains real class definitions (should be pseudocode only)${NC}"
        ((errors++))
    fi

    # Should have library choices
    if ! grep -qi "library\|framework\|choice\|justification" "$file"; then
        echo -e "${YELLOW}⚠️  WARNING: No library/framework choices documented${NC}"
    fi

    # Should have adaptation notes
    if ! grep -qi "adaptation\|lesson\|previous" "$file"; then
        echo -e "${YELLOW}⚠️  WARNING: No adaptation notes found${NC}"
    fi

    if [ $errors -eq 0 ]; then
        echo -e "${GREEN}✅ PASSED: Pseudocode fidelity validated${NC}"
        return 0
    else
        echo -e "${RED}FAILED: $errors pseudocode validation errors${NC}"
        return 1
    fi
}

validate_concrete() {
    local file="$1"
    local errors=0

    echo "Validating CONCRETE fidelity (Wave Architecture)..."

    # Should have real code (function definitions)
    if ! grep -q "def .*(" "$file"; then
        echo -e "${RED}❌ FAILED: No real function definitions found (expected concrete code)${NC}"
        ((errors++))
    fi

    # Should have class definitions
    if ! grep -q "class .*:" "$file"; then
        echo -e "${RED}❌ FAILED: No real class definitions found (expected concrete code)${NC}"
        ((errors++))
    fi

    # Should have docstrings
    if ! grep -q '"""' "$file"; then
        echo -e "${RED}❌ FAILED: No docstrings found (expected documented code)${NC}"
        ((errors++))
    fi

    # Should have type hints
    if ! grep -q " -> " "$file" && ! grep -q ": str\|: int\|: bool\|: Optional\|: List\|: Dict" "$file"; then
        echo -e "${YELLOW}⚠️  WARNING: No type hints found (recommended for concrete code)${NC}"
    fi

    # Should have working examples
    if ! grep -qi "example\|usage" "$file"; then
        echo -e "${YELLOW}⚠️  WARNING: No usage examples found${NC}"
    fi

    # Should have adaptation notes
    if ! grep -qi "adaptation\|lesson\|previous" "$file"; then
        echo -e "${YELLOW}⚠️  WARNING: No adaptation notes found${NC}"
    fi

    if [ $errors -eq 0 ]; then
        echo -e "${GREEN}✅ PASSED: Concrete code fidelity validated${NC}"
        return 0
    else
        echo -e "${RED}FAILED: $errors concrete validation errors${NC}"
        return 1
    fi
}

validate_exact() {
    local file="$1"
    local errors=0

    echo "Validating EXACT fidelity (Wave Implementation)..."

    # Should have R213 metadata
    if ! grep -q "effort_id\|EFFORT_" "$file"; then
        echo -e "${RED}❌ FAILED: No effort_id (R213 metadata) found${NC}"
        ((errors++))
    fi

    # Should have exact file lists
    if ! grep -qi "files to create\|files to modify\|new files\|modified files" "$file"; then
        echo -e "${RED}❌ FAILED: No exact file lists found${NC}"
        ((errors++))
    fi

    # Should have line count estimates
    if ! grep -qE "[0-9]+ lines" "$file"; then
        echo -e "${RED}❌ FAILED: No line count estimates found${NC}"
        ((errors++))
    fi

    # Should have real code specifications
    if ! grep -q "def .*(" "$file"; then
        echo -e "${RED}❌ FAILED: No real code specifications found${NC}"
        ((errors++))
    fi

    # Should have test specifications
    if ! grep -qi "test\|coverage" "$file"; then
        echo -e "${RED}❌ FAILED: No test specifications found${NC}"
        ((errors++))
    fi

    # Should have dependencies
    if ! grep -qi "depend\|upstream\|downstream" "$file"; then
        echo -e "${YELLOW}⚠️  WARNING: No dependencies documented${NC}"
    fi

    # Should have acceptance criteria
    if ! grep -qi "acceptance\|criteria\|definition of done" "$file"; then
        echo -e "${YELLOW}⚠️  WARNING: No acceptance criteria found${NC}"
    fi

    if [ $errors -eq 0 ]; then
        echo -e "${GREEN}✅ PASSED: Exact specification fidelity validated${NC}"
        return 0
    else
        echo -e "${RED}FAILED: $errors exact specification validation errors${NC}"
        return 1
    fi
}

# Main validation logic
main() {
    if [ $# -ne 2 ]; then
        echo "Usage: $0 <file> <expected_level>"
        echo "Expected levels: pseudocode, concrete, exact"
        exit 1
    fi

    if [ ! -f "$FILE" ]; then
        echo -e "${RED}ERROR: File not found: $FILE${NC}"
        exit 1
    fi

    echo "========================================="
    echo "Fidelity Gradient Validator"
    echo "========================================="
    echo "File: $FILE"
    echo "Expected Fidelity: $EXPECTED_LEVEL"
    echo "========================================="
    echo ""

    case "$EXPECTED_LEVEL" in
        pseudocode)
            validate_pseudocode "$FILE"
            exit $?
            ;;
        concrete)
            validate_concrete "$FILE"
            exit $?
            ;;
        exact)
            validate_exact "$FILE"
            exit $?
            ;;
        *)
            echo -e "${RED}ERROR: Invalid fidelity level: $EXPECTED_LEVEL${NC}"
            echo "Valid levels: pseudocode, concrete, exact"
            exit 1
            ;;
    esac
}

main "$@"
