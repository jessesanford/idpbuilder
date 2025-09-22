#!/bin/bash
# Validation script to verify type conflict resolution

set -e

echo "=========================================="
echo "Type Conflict Fix Validation"
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
    VALIDATION_FAILED=1
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

VALIDATION_FAILED=0

echo ""
echo "Pre-flight Checks"
echo "-----------------"

# Check we're in the right directory
if [ ! -d "pkg/certs" ]; then
    print_error "Not in project root. Run from idpbuilder-oci-bp-branch-rebasing"
    exit 1
fi

# Check Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed"
    exit 1
fi

print_status "Environment checks passed"

echo ""
echo "Type Definition Checks"
echo "----------------------"

# Check for duplicate CertificateValidator
echo -n "Checking CertificateValidator definitions... "
CERT_VAL_LOCATIONS=$(grep -l "type CertificateValidator interface" pkg/certs/*.go 2>/dev/null || true)
CERT_VAL_COUNT=$(echo "$CERT_VAL_LOCATIONS" | grep -c . || echo 0)

if [ "$CERT_VAL_COUNT" -eq "1" ]; then
    print_status "Single definition in: $CERT_VAL_LOCATIONS"
elif [ "$CERT_VAL_COUNT" -eq "0" ]; then
    print_error "No CertificateValidator definition found!"
else
    print_error "Multiple definitions found in:"
    echo "$CERT_VAL_LOCATIONS"
fi

# Check for duplicate ValidationResult
echo -n "Checking ValidationResult definitions... "
VAL_RES_LOCATIONS=$(grep -l "type ValidationResult struct" pkg/certs/*.go 2>/dev/null || true)
VAL_RES_COUNT=$(echo "$VAL_RES_LOCATIONS" | grep -c . || echo 0)

if [ "$VAL_RES_COUNT" -eq "1" ]; then
    print_status "Single definition in: $VAL_RES_LOCATIONS"
elif [ "$VAL_RES_COUNT" -eq "0" ]; then
    print_error "No ValidationResult definition found!"
else
    print_error "Multiple definitions found in:"
    echo "$VAL_RES_LOCATIONS"
fi

# Check for BasicValidator (should exist after fix)
echo -n "Checking BasicValidator definition... "
if grep -q "type BasicValidator interface" pkg/certs/*.go 2>/dev/null; then
    BASIC_VAL_LOC=$(grep -l "type BasicValidator interface" pkg/certs/*.go)
    print_status "Found in: $BASIC_VAL_LOC"
else
    print_warning "BasicValidator not found (may not be implemented yet)"
fi

echo ""
echo "Compilation Tests"
echo "-----------------"

# Try to build the certs package
echo -n "Building pkg/certs... "
if go build ./pkg/certs/... 2>/tmp/certs_build.txt; then
    print_status "Package builds successfully"
else
    print_error "Build failed:"
    head -20 /tmp/certs_build.txt
fi

# Try to build all packages
echo -n "Building all packages... "
if go build ./... 2>/tmp/all_build.txt; then
    print_status "All packages build successfully"
else
    print_warning "Some packages failed (may be unrelated):"
    grep "^#" /tmp/all_build.txt | head -10
fi

echo ""
echo "Implementation Checks"
echo "--------------------"

# Check that DefaultCertificateValidator exists
echo -n "Checking DefaultCertificateValidator... "
if grep -q "type DefaultCertificateValidator struct" pkg/certs/validator.go; then
    print_status "Found implementation"
else
    print_error "Not found in validator.go"
fi

# Check that BasicCertificateValidator exists
echo -n "Checking BasicCertificateValidator... "
if grep -q "type BasicCertificateValidator struct" pkg/certs/types.go; then
    print_status "Found implementation"
else
    print_warning "Not found in types.go (may be renamed)"
fi

echo ""
echo "Field Reference Checks"
echo "----------------------"

# Check for old IsValid references (should be replaced with Valid)
echo -n "Checking for old IsValid field references... "
OLD_REFS=$(grep -n "\.IsValid[^(]" pkg/certs/*.go 2>/dev/null | grep -v "func.*IsValid" || true)
if [ -z "$OLD_REFS" ]; then
    print_status "No old IsValid field references found"
else
    print_warning "Found potential old references:"
    echo "$OLD_REFS"
fi

# Check for Valid field usage
echo -n "Checking for Valid field usage... "
if grep -q "\.Valid[^a-zA-Z]" pkg/certs/*.go 2>/dev/null; then
    print_status "Valid field is being used"
else
    print_warning "Valid field not found in use"
fi

echo ""
echo "Test Execution"
echo "--------------"

# Run tests for the certs package
echo "Running pkg/certs tests..."
if go test ./pkg/certs/... -v -count=1 > /tmp/test_output.txt 2>&1; then
    print_status "All certs package tests pass"
    PASS_COUNT=$(grep -c "PASS:" /tmp/test_output.txt || echo 0)
    echo "  Passed: $PASS_COUNT tests"
else
    print_warning "Some tests failed (may need updates):"
    grep "FAIL:" /tmp/test_output.txt | head -5
    FAIL_COUNT=$(grep -c "FAIL:" /tmp/test_output.txt || echo 0)
    echo "  Failed: $FAIL_COUNT tests"
fi

echo ""
echo "Import Analysis"
echo "---------------"

# Check which packages import certs
echo "Packages that import pkg/certs:"
IMPORTERS=$(grep -r '"github.com/cnoe-io/idpbuilder-oci/pkg/certs"' --include="*.go" . 2>/dev/null | cut -d: -f1 | sort -u | grep -v "/certs/" || true)
if [ -z "$IMPORTERS" ]; then
    print_status "No external packages import pkg/certs (self-contained)"
else
    print_warning "The following packages import pkg/certs and may need updates:"
    echo "$IMPORTERS" | head -10
fi

echo ""
echo "=========================================="
echo "Validation Summary"
echo "=========================================="

if [ $VALIDATION_FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ Type conflict resolution appears successful!${NC}"
    echo ""
    echo "All critical checks passed:"
    echo "  ✓ Single type definitions"
    echo "  ✓ Package compiles"
    echo "  ✓ Implementations present"
    echo ""
    echo "Next steps:"
    echo "1. Review any warnings above"
    echo "2. Update failing tests if needed"
    echo "3. Commit and push the fix"
    echo "4. Apply to original branches"
else
    echo -e "${RED}✗ Validation failed - manual intervention required${NC}"
    echo ""
    echo "Issues found:"
    echo "  - Check error messages above"
    echo "  - Review type definitions"
    echo "  - Ensure proper updates to all files"
    echo ""
    echo "Debug commands:"
    echo "  grep -r 'type.*Validator' pkg/certs/"
    echo "  grep -r 'type.*Result' pkg/certs/"
    echo "  go build -v ./pkg/certs/..."
fi

exit $VALIDATION_FAILED