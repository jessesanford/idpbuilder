# Certificate Validation Pipeline Implementation Plan

**Effort ID**: E1.2.1  
**Effort Name**: certificate-validation-pipeline  
**Phase**: 1 - Certificate Infrastructure  
**Wave**: 2 - Certificate Validation & Fallback  
**Created By**: Code Reviewer (code-reviewer)  
**Date Created**: 2025-08-31  
**Assigned To**: SW Engineer 1  

## =¨ CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `idpbuidler-oci-go-cr/phase1/wave2/certificate-validation-pipeline`
**Can Parallelize**: Yes
**Parallel With**: E1.2.2 (fallback-strategies)
**Size Estimate**: ~400 lines
**Dependencies**: E1.1.1 (kind-certificate-extraction), E1.1.2 (registry-tls-trust-integration)

<!--   ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE   -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## =Ë Effort Overview

### Description
This effort implements a comprehensive certificate validation pipeline that validates X.509 certificate chains, checks expiry dates with warnings for soon-to-expire certificates, and verifies hostname matching including wildcard support. It provides clear diagnostics for troubleshooting certificate issues, a critical component for the self-signed certificate handling in the IDPBuilder OCI MVP.

### Size Estimate
- **Estimated Lines**: 400 (well within limit)
- **Confidence Level**: High
- **Split Risk**: Low

### Dependencies
- **Requires**: 
  - E1.1.1 (kind-certificate-extraction) - Provides extracted certificates to validate
  - E1.1.2 (registry-tls-trust-integration) - Uses trust store for chain validation
- **Blocks**: 
  - E2.1.2 (gitea-registry-client) - Needs validation before push operations
- **External**: 
  - Standard library crypto/x509 for certificate operations
  - Standard library time for expiry calculations

## <¯ Requirements

### Functional Requirements
- [ ] Validate complete X.509 certificate chains against system and custom trust stores
- [ ] Check certificate expiry with configurable warning threshold (default 30 days)
- [ ] Verify hostname matches certificate CN/SAN including wildcard support
- [ ] Generate comprehensive diagnostics for all validation failures
- [ ] Support both self-signed and CA-signed certificates

### Non-Functional Requirements
- [ ] Performance: Validation must complete in <100ms for typical certificates
- [ ] Security: Never bypass validation without explicit user consent
- [ ] Maintainability: Clear separation of validation concerns
- [ ] Scalability: Support validation of multiple certificates concurrently

### Acceptance Criteria
- [ ] All unit tests passing
- [ ] Test coverage e 80%
- [ ] Code review approved
- [ ] Size d 800 lines (measured with line-counter.sh)
- [ ] No critical TODOs
- [ ] Documentation complete

## =Á Implementation Details

### Files to Create
| File Path | Purpose | Estimated Lines |
|-----------|---------|-----------------|
| `pkg/certs/validator.go` | Main validation logic and CertValidator interface | 180 |
| `pkg/certs/diagnostics.go` | Diagnostic generation and formatting | 80 |
| `pkg/certs/validator_test.go` | Comprehensive test suite | 120 |
| `pkg/certs/testdata/certs.go` | Test certificate fixtures | 20 |
| **Total** | | 400 |

### Files to Modify
None - this is a new component that will be integrated by consumers.

### Key Components

#### Component 1: CertValidator Interface
**Purpose**: Define the contract for certificate validation operations  
**Location**: `pkg/certs/validator.go`  
**Lines**: ~40  

```go
// CertValidator provides comprehensive X.509 certificate validation
type CertValidator interface {
    // ValidateChain verifies the certificate chain against trusted roots
    ValidateChain(cert *x509.Certificate) error
    
    // CheckExpiry checks if certificate is expired or expiring soon
    // Returns remaining validity duration and any warnings
    CheckExpiry(cert *x509.Certificate) (*time.Duration, error)
    
    // VerifyHostname checks if the certificate is valid for the given hostname
    VerifyHostname(cert *x509.Certificate, hostname string) error
    
    // GenerateDiagnostics creates a detailed diagnostic report for the certificate
    GenerateDiagnostics(cert *x509.Certificate) (*CertDiagnostics, error)
}

// CertDiagnostics contains detailed information about certificate validation
type CertDiagnostics struct {
    Subject         string
    Issuer          string
    SerialNumber    string
    NotBefore       time.Time
    NotAfter        time.Time
    DNSNames        []string
    IPAddresses     []net.IP
    ValidationErrors []ValidationError
    Warnings        []string
}

// ValidationError represents a specific validation failure
type ValidationError struct {
    Type    string // "chain", "expiry", "hostname", etc.
    Message string
    Detail  string
}
```

#### Component 2: DefaultValidator Implementation
**Purpose**: Concrete implementation of CertValidator with integration to trust store  
**Location**: `pkg/certs/validator.go`  
**Lines**: ~140  

```go
// DefaultValidator implements CertValidator with configurable options
type DefaultValidator struct {
    trustStore      *TrustStoreManager // From E1.1.2
    expiryWarning   time.Duration      // Default 30 days
    systemRoots     *x509.CertPool     // System CA certificates
    customRoots     *x509.CertPool     // Custom CA certificates from trust store
}

// NewValidator creates a validator with trust store integration
func NewValidator(trustStore *TrustStoreManager) (*DefaultValidator, error) {
    // Load system roots
    // Initialize custom roots from trust store
    // Set default expiry warning
    return validator, nil
}

// ValidateChain implementation with detailed error reporting
func (v *DefaultValidator) ValidateChain(cert *x509.Certificate) error {
    // Build verification options
    // Try system roots first
    // Fall back to custom roots
    // Return detailed error on failure
}

// CheckExpiry with configurable warning threshold
func (v *DefaultValidator) CheckExpiry(cert *x509.Certificate) (*time.Duration, error) {
    // Calculate time until expiry
    // Check if expired
    // Check if within warning threshold
    // Return duration and any warnings
}

// VerifyHostname with wildcard support
func (v *DefaultValidator) VerifyHostname(cert *x509.Certificate, hostname string) error {
    // Use x509.VerifyHostname
    // Provide clear error messages
    // List valid hostnames in error
}
```

#### Component 3: Diagnostic Generator
**Purpose**: Generate human-readable diagnostic reports for troubleshooting  
**Location**: `pkg/certs/diagnostics.go`  
**Lines**: ~80  

```go
// GenerateDiagnostics creates comprehensive diagnostic report
func (v *DefaultValidator) GenerateDiagnostics(cert *x509.Certificate) (*CertDiagnostics, error) {
    diag := &CertDiagnostics{
        Subject:      cert.Subject.String(),
        Issuer:       cert.Issuer.String(),
        SerialNumber: cert.SerialNumber.String(),
        NotBefore:    cert.NotBefore,
        NotAfter:     cert.NotAfter,
        DNSNames:     cert.DNSNames,
        IPAddresses:  cert.IPAddresses,
    }
    
    // Run all validations and collect errors
    // Add warnings for soon-to-expire
    // Format for human readability
    
    return diag, nil
}

// FormatDiagnostics returns human-readable diagnostic output
func FormatDiagnostics(diag *CertDiagnostics) string {
    // Format as clear, structured text
    // Include all relevant details
    // Highlight errors and warnings
}
```

## >ê Testing Strategy

### Unit Tests Required
- [ ] Test file: `pkg/certs/validator_test.go`
- [ ] Coverage target: 80%
- [ ] Test cases:
  - [ ] Valid certificate chain validation
  - [ ] Self-signed certificate validation
  - [ ] Expired certificate detection
  - [ ] Soon-to-expire warning (< 30 days)
  - [ ] Hostname match validation
  - [ ] Wildcard certificate matching
  - [ ] Hostname mismatch error
  - [ ] Chain validation with missing intermediate
  - [ ] Diagnostic generation for various scenarios

### Test Fixtures
```go
// pkg/certs/testdata/certs.go
// Generate test certificates for various scenarios:
// - Valid certificate
// - Expired certificate  
// - Self-signed certificate
// - Certificate with wildcard CN
// - Certificate with SAN entries
// - Certificate expiring in 15 days
```

### Integration Tests
- [ ] Integration with trust store from E1.1.2
- [ ] Validation of real Gitea certificates from E1.1.1
- [ ] End-to-end validation pipeline

## = Integration Points

### With E1.1.1 (kind-certificate-extraction)
- Receive extracted certificates for validation
- Use KindCertExtractor.ExtractGiteaCert() output as input

### With E1.1.2 (registry-tls-trust-integration)  
- Use TrustStoreManager for custom CA certificates
- Integrate with trust store for chain validation

### With E1.2.2 (fallback-strategies)
- Provide validation errors for fallback handler
- Generate diagnostics for fallback recommendations

### With Future E2.1.2 (gitea-registry-client)
- Called before any registry push operation
- Validation errors trigger fallback strategies

## =Ê Size Management Strategy

### Measurement Protocol
```bash
# Measure using project line counter (R200 compliance)
cd /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/wave2/certificate-validation-pipeline
$PROJECT_ROOT/tools/line-counter.sh

# Check points:
# - After validator.go implementation (~180 lines)
# - After diagnostics.go implementation (~260 lines)  
# - After tests implementation (~380 lines)
# - Final measurement (~400 lines)
```

### Split Prevention
- Core validation logic kept concise (~180 lines)
- Diagnostics separated to own file (~80 lines)
- Test fixtures minimal (~20 lines)
- Total well under 800 line limit with 400 line buffer

## =€ Implementation Sequence

1. **Define Interfaces and Types** (30 min)
   - Create CertValidator interface
   - Define CertDiagnostics struct
   - Define ValidationError type

2. **Implement Core Validation** (2 hours)
   - ValidateChain with trust store integration
   - CheckExpiry with warning threshold
   - VerifyHostname with wildcard support

3. **Add Diagnostic Generation** (1 hour)
   - GenerateDiagnostics method
   - FormatDiagnostics helper
   - Error collection and formatting

4. **Create Test Suite** (2 hours)
   - Generate test certificates
   - Unit tests for all methods
   - Integration test with trust store

5. **Documentation** (30 min)
   - Code comments
   - Package documentation
   - Usage examples

## =Ý Notes for SW Engineer

### Critical Considerations
1. **Trust Store Integration**: Must use TrustStoreManager from E1.1.2 for custom roots
2. **Clear Error Messages**: Each validation failure must have actionable error message
3. **Expiry Warning**: Default to 30 days but make configurable
4. **Wildcard Support**: Use standard x509.VerifyHostname which handles wildcards
5. **Diagnostic Output**: Should be human-readable and help troubleshooting

### Example Usage
```go
// Create validator with trust store
trustStore := getTrustStoreFromE112()
validator, err := NewValidator(trustStore)

// Validate a certificate
cert := getCertificateFromE111()
if err := validator.ValidateChain(cert); err != nil {
    // Generate diagnostics for troubleshooting
    diag, _ := validator.GenerateDiagnostics(cert)
    fmt.Println(FormatDiagnostics(diag))
    return err
}

// Check expiry
duration, err := validator.CheckExpiry(cert)
if err != nil {
    log.Warnf("Certificate expiring soon: %v", duration)
}

// Verify hostname
if err := validator.VerifyHostname(cert, "gitea.local"); err != nil {
    return fmt.Errorf("hostname verification failed: %w", err)
}
```

### Dependencies to Import
```go
import (
    "crypto/x509"
    "fmt"
    "net"
    "time"
    
    // From other efforts
    "github.com/idpbuilder/idpbuilder-oci-go-cr/pkg/certs" // TrustStoreManager from E1.1.2
)
```

##  Review Checklist
- [ ] All functional requirements addressed
- [ ] Size within 400 line estimate
- [ ] Test coverage plan adequate
- [ ] Integration points clear
- [ ] Dependencies properly identified
- [ ] Implementation sequence logical