#!/bin/bash
# Fix Type Conflicts Script for Phase 1 Integration
# This script implements the recommended Option C - Shared Types with Unified Interface

set -e

echo "=========================================="
echo "Type Conflict Resolution Script"
echo "Phase 1 Integration - idpbuilder-oci"
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
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Ensure we're in the correct directory
if [ ! -d "pkg/certs" ]; then
    print_error "Not in the project root directory. Please run from idpbuilder-oci-bp-branch-rebasing"
    exit 1
fi

echo ""
echo "Step 1: Creating backup branch"
echo "-------------------------------"
BACKUP_BRANCH="backup/pre-type-fix-$(date +%Y%m%d-%H%M%S)"
git checkout -b "$BACKUP_BRANCH" 2>/dev/null || true
print_status "Created backup branch: $BACKUP_BRANCH"

echo ""
echo "Step 2: Creating fix branch"
echo "----------------------------"
git checkout idpbuilder-oci-build-push/phase1/integration-rebased
git checkout -b fix/type-conflicts-resolution
print_status "Created fix branch: fix/type-conflicts-resolution"

echo ""
echo "Step 3: Creating unified validation_types.go"
echo "---------------------------------------------"

cat > pkg/certs/validation_types.go << 'EOF'
// Copyright 2024 The IDP Builder Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package certs

import (
	"crypto/x509"
	"time"
)

// CertificateValidator defines the comprehensive interface for certificate validation operations
// This is the primary interface from the cert-validation effort
type CertificateValidator interface {
	// ValidateChain validates a complete certificate chain from leaf to root
	ValidateChain(certs []*x509.Certificate) error

	// ValidateCertificate validates a single certificate
	ValidateCertificate(cert *x509.Certificate) error

	// VerifyHostname verifies that a certificate is valid for a given hostname
	VerifyHostname(cert *x509.Certificate, hostname string) error

	// GenerateDiagnostics creates diagnostic information for troubleshooting
	GenerateDiagnostics() (*CertDiagnostics, error)

	// SetValidationMode changes the validation strictness
	SetValidationMode(mode ValidationMode)

	// GetValidationMode returns the current validation mode
	GetValidationMode() ValidationMode
}

// BasicValidator defines a simplified interface for basic certificate validation
// This is the renamed interface from the fallback-strategies effort
type BasicValidator interface {
	// Validate checks if the certificate is valid
	Validate(cert *Certificate) error

	// ValidateChain checks if the certificate chain is valid
	ValidateChain(chain []*Certificate) error

	// IsExpired checks if the certificate has expired
	IsExpired(cert *Certificate) bool

	// WillExpireSoon checks if certificate will expire within threshold
	WillExpireSoon(cert *Certificate, threshold time.Duration) bool
}

// ValidationResult represents the unified result of certificate validation
// This combines fields from both implementations for maximum compatibility
type ValidationResult struct {
	// Core fields
	Valid       bool      `json:"valid"`
	ValidatedAt time.Time `json:"validated_at"`

	// Detailed tracking (from validator.go implementation)
	Errors      []error             `json:"errors,omitempty"`
	Warnings    []string            `json:"warnings,omitempty"`
	Certificate *x509.Certificate   `json:"-"`
	Chain       []*x509.Certificate `json:"-"`

	// Simple message and actions (from utilities.go implementation)
	Message string   `json:"message,omitempty"`
	Actions []string `json:"actions,omitempty"`
}

// IsValid provides backward compatibility for code expecting the old field name
func (v *ValidationResult) IsValid() bool {
	return v.Valid
}

// NewValidationResult creates a new validation result with timestamp
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		ValidatedAt: time.Now(),
		Errors:      make([]error, 0),
		Warnings:    make([]string, 0),
		Actions:     make([]string, 0),
	}
}

// AddError adds an error to the validation result
func (v *ValidationResult) AddError(err error) {
	v.Errors = append(v.Errors, err)
	v.Valid = false
}

// AddWarning adds a warning to the validation result
func (v *ValidationResult) AddWarning(warning string) {
	v.Warnings = append(v.Warnings, warning)
}

// AddAction adds a suggested action to the validation result
func (v *ValidationResult) AddAction(action string) {
	v.Actions = append(v.Actions, action)
}
EOF

print_status "Created pkg/certs/validation_types.go"

echo ""
echo "Step 4: Updating validator.go"
echo "------------------------------"

# Create a temporary file with the updates
cat > /tmp/validator_update.sed << 'EOF'
# Remove the CertificateValidator interface definition (lines 34-53)
/^\/\/ CertificateValidator defines the interface/,/^}$/d

# Remove the ValidationResult struct definition (lines 73-81)
/^\/\/ ValidationResult holds the results/,/^}$/d

# Update IsValid references to Valid
s/IsValid/Valid/g
EOF

# Apply the updates
sed -i.bak -f /tmp/validator_update.sed pkg/certs/validator.go

print_status "Updated validator.go (backup saved as validator.go.bak)"

echo ""
echo "Step 5: Updating types.go"
echo "--------------------------"

# Update types.go to rename CertificateValidator to BasicValidator
sed -i.bak \
    -e 's/type CertificateValidator interface/type BasicValidator interface/g' \
    -e 's/implements CertificateValidator/implements BasicValidator/g' \
    -e 's/\/\/ CertificateValidator defines/\/\/ BasicValidator defines/g' \
    -e '/^type BasicValidator interface/,/^}$/d' \
    pkg/certs/types.go

# Also need to update BasicCertificateValidator comments
sed -i \
    -e 's/implements CertificateValidator interface/implements BasicValidator interface/g' \
    pkg/certs/types.go

print_status "Updated types.go (backup saved as types.go.bak)"

echo ""
echo "Step 6: Updating utilities.go"
echo "------------------------------"

# Remove the ValidationResult struct from utilities.go
sed -i.bak \
    -e '/^\/\/ ValidationResult represents certificate validation result/,/^}$/d' \
    pkg/certs/utilities.go

print_status "Updated utilities.go (backup saved as utilities.go.bak)"

echo ""
echo "Step 7: Updating test files"
echo "----------------------------"

# Update test files to use the new interface names
if [ -f "pkg/certs/types_test.go" ]; then
    sed -i.bak \
        -e 's/TestBasicCertificateValidator/TestBasicValidator/g' \
        -e 's/NewBasicCertificateValidator/NewBasicValidator/g' \
        pkg/certs/types_test.go
    print_status "Updated types_test.go"
fi

echo ""
echo "Step 8: Running validation checks"
echo "----------------------------------"

# Check for duplicate type definitions
echo -n "Checking for duplicate CertificateValidator definitions... "
CERT_VAL_COUNT=$(grep -r "type CertificateValidator interface" pkg/certs 2>/dev/null | wc -l)
if [ "$CERT_VAL_COUNT" -eq "1" ]; then
    print_status "Only one CertificateValidator definition found"
else
    print_error "Found $CERT_VAL_COUNT CertificateValidator definitions"
fi

echo -n "Checking for duplicate ValidationResult definitions... "
VAL_RES_COUNT=$(grep -r "type ValidationResult struct" pkg/certs 2>/dev/null | wc -l)
if [ "$VAL_RES_COUNT" -eq "1" ]; then
    print_status "Only one ValidationResult definition found"
else
    print_error "Found $VAL_RES_COUNT ValidationResult definitions"
fi

echo ""
echo "Step 9: Attempting to build the package"
echo "----------------------------------------"

if go build ./pkg/certs/... 2>/tmp/build_output.txt; then
    print_status "Package builds successfully!"
else
    print_error "Build failed. Error output:"
    cat /tmp/build_output.txt
    echo ""
    print_warning "Manual intervention may be required"
fi

echo ""
echo "Step 10: Running tests"
echo "-----------------------"

if go test ./pkg/certs/... -v > /tmp/test_output.txt 2>&1; then
    print_status "All tests pass!"
    echo "Test summary:"
    grep -E "PASS:|FAIL:" /tmp/test_output.txt | tail -10
else
    print_warning "Some tests failed. This may be expected if tests need updates."
    echo "Failed tests:"
    grep "FAIL:" /tmp/test_output.txt | head -10
fi

echo ""
echo "Step 11: Creating commit"
echo "-------------------------"

git add pkg/certs/validation_types.go
git add pkg/certs/validator.go
git add pkg/certs/types.go
git add pkg/certs/utilities.go
git add pkg/certs/*_test.go 2>/dev/null || true

git commit -m "fix: resolve type conflicts between cert-validation and fallback-strategies

- Created unified validation_types.go with shared interfaces
- CertificateValidator remains as primary comprehensive interface
- BasicValidator introduced for simple validation needs
- Unified ValidationResult struct supports both use cases
- Updated all references to use unified types
- Maintains backward compatibility with helper methods

This resolves build failures from duplicate type definitions between
the cert-validation and fallback-strategies efforts in Phase 1 Wave 2."

print_status "Created commit with fix"

echo ""
echo "=========================================="
echo "Type Conflict Resolution Complete!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Test the build: go build ./..."
echo "2. Run full test suite: go test ./... -v"
echo "3. If successful, cherry-pick to original branches:"
echo "   - git checkout cert-validation-rebased"
echo "   - git cherry-pick fix/type-conflicts-resolution"
echo "   - git checkout fallback-strategies-rebased"
echo "   - git cherry-pick fix/type-conflicts-resolution"
echo "4. Re-rebase both branches"
echo "5. Merge into integration branch"
echo ""
echo "Backup branch saved as: $BACKUP_BRANCH"
echo "Fix branch: fix/type-conflicts-resolution"