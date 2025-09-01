# Fallback Strategies Implementation Plan

**Effort ID**: E1.2.2  
**Effort Name**: fallback-strategies  
**Phase**: 1 - Certificate Infrastructure  
**Wave**: 2 - Certificate Validation & Fallback  
**Created By**: Code Reviewer (code-reviewer)  
**Date Created**: 2025-09-01  
**Assigned To**: SW Engineer 2  

## CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `idpbuidler-oci-go-cr/phase1/wave2/fallback-strategies`
**Can Parallelize**: Yes
**Parallel With**: E1.2.1 (certificate-validation-pipeline)
**Size Estimate**: ~400 lines
**Dependencies**: E1.2.1 (certificate-validation-pipeline)

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Effort Overview

### Description
This effort implements comprehensive fallback strategies for certificate-related issues in registry operations. It auto-detects certificate problems (self-signed, expired, hostname mismatches, untrusted CAs), provides actionable solution recommendations, and implements the --insecure flag for development environments. This is critical for enabling developers to work with self-signed certificates in local Kind clusters.

### Size Estimate
- **Estimated Lines**: 400 (well within limit)
- **Confidence Level**: High
- **Split Risk**: Low

### Dependencies
- **Requires**: 
  - E1.2.1 (certificate-validation-pipeline) - Imports CertValidator interface and ValidationError types
- **Blocks**: 
  - E2.1.2 (gitea-registry-client) - Fallback strategies needed for push operations
- **External**: 
  - Standard library crypto/tls for insecure flag implementation
  - Standard library crypto/x509 for certificate error types

## Requirements

### Functional Requirements
- [ ] Auto-detect certificate validation failures from E1.2.1's CertValidator
- [ ] Identify specific certificate problem types (self-signed, expired, hostname mismatch, untrusted CA)
- [ ] Generate actionable solution recommendations for each problem type
- [ ] Implement --insecure flag for CLI commands to skip TLS verification
- [ ] Display prominent security warnings when --insecure is used
- [ ] Log detailed certificate chain information for debugging

### Non-Functional Requirements
- [ ] Performance: Problem detection must add <10ms overhead
- [ ] Security: --insecure flag must require explicit user consent
- [ ] UX: Error messages must be clear and actionable
- [ ] Maintainability: Clear separation between detection and recommendation logic

### Acceptance Criteria
- [ ] All unit tests passing
- [ ] Test coverage >= 80%
- [ ] Code review approved
- [ ] Size <= 800 lines (measured with line-counter.sh)
- [ ] No critical TODOs
- [ ] Documentation complete
- [ ] Integration with E1.2.1 validated

## Implementation Details

### Files to Create
| File Path | Purpose | Estimated Lines |
|-----------|---------|-----------------|
| `pkg/fallback/detector.go` | Certificate problem detection logic | 100 |
| `pkg/fallback/recommender.go` | Solution recommendation engine | 80 |
| `pkg/fallback/insecure.go` | --insecure flag implementation and TLS config | 60 |
| `pkg/fallback/logger.go` | Certificate chain logging utilities | 40 |
| `pkg/fallback/detector_test.go` | Tests for problem detection | 60 |
| `pkg/fallback/recommender_test.go` | Tests for recommendations | 40 |
| `pkg/fallback/insecure_test.go` | Tests for insecure mode | 20 |
| **Total** | | 400 |

### Files to Modify
None - this is a new component that will be integrated by consumers.

### Key Components

#### Component 1: Problem Detector
**Purpose**: Analyze certificate validation errors and identify specific problem types  
**Location**: `pkg/fallback/detector.go`  
**Lines**: ~100  

```go
// ProblemType represents specific certificate issues
type ProblemType string

const (
    ProblemSelfSigned     ProblemType = "self-signed"
    ProblemExpired        ProblemType = "expired"
    ProblemHostnameMismatch ProblemType = "hostname-mismatch"
    ProblemUntrustedCA    ProblemType = "untrusted-ca"
    ProblemUnknown        ProblemType = "unknown"
)

// CertProblem contains detailed information about a certificate issue
type CertProblem struct {
    Type        ProblemType
    Certificate *x509.Certificate
    Error       error
    Details     map[string]interface{}
}

// ProblemDetector analyzes certificate validation errors
type ProblemDetector interface {
    // DetectProblem analyzes a validation error from CertValidator
    DetectProblem(validationErr error, cert *x509.Certificate) (*CertProblem, error)
    
    // AnalyzeCertChain provides detailed chain analysis
    AnalyzeCertChain(certs []*x509.Certificate) ([]*CertProblem, error)
}

// DefaultDetector implements ProblemDetector
type DefaultDetector struct {
    validator *certs.CertValidator // From E1.2.1
}

// NewDetector creates a detector integrated with CertValidator
func NewDetector(validator *certs.CertValidator) *DefaultDetector {
    return &DefaultDetector{validator: validator}
}

// DetectProblem implementation with pattern matching
func (d *DefaultDetector) DetectProblem(validationErr error, cert *x509.Certificate) (*CertProblem, error) {
    // Check for x509.UnknownAuthorityError (self-signed or untrusted)
    // Check for x509.CertificateInvalidError (expired, not yet valid)
    // Check for x509.HostnameError (hostname mismatch)
    // Extract specific details for each problem type
    // Return structured CertProblem
}
```

#### Component 2: Solution Recommender
**Purpose**: Generate actionable recommendations for each problem type  
**Location**: `pkg/fallback/recommender.go`  
**Lines**: ~80  

```go
// Recommendation contains a suggested solution
type Recommendation struct {
    Priority    int      // 1=highest, for ordering multiple recommendations
    Title       string   // Brief description
    Command     string   // Example command to run
    Explanation string   // Detailed explanation
    Risks       []string // Security implications
}

// Recommender generates solutions for certificate problems
type Recommender interface {
    // Recommend generates solutions for a detected problem
    Recommend(problem *CertProblem) ([]*Recommendation, error)
    
    // FormatRecommendations creates user-friendly output
    FormatRecommendations(recs []*Recommendation) string
}

// DefaultRecommender implements Recommender
type DefaultRecommender struct {
    registryURL string
    insecureAllowed bool
}

// NewRecommender creates a recommender with context
func NewRecommender(registryURL string, allowInsecure bool) *DefaultRecommender {
    return &DefaultRecommender{
        registryURL: registryURL,
        insecureAllowed: allowInsecure,
    }
}

// Recommend implementation with problem-specific logic
func (r *DefaultRecommender) Recommend(problem *CertProblem) ([]*Recommendation, error) {
    switch problem.Type {
    case ProblemSelfSigned:
        // Recommend: 1) Add to trust store, 2) Use --insecure for dev
    case ProblemExpired:
        // Recommend: 1) Renew certificate, 2) Use --insecure temporarily
    case ProblemHostnameMismatch:
        // Recommend: 1) Use correct hostname, 2) Update certificate SANs
    case ProblemUntrustedCA:
        // Recommend: 1) Import CA certificate, 2) Verify CA chain
    }
}

// FormatRecommendations with clear, actionable output
func (r *DefaultRecommender) FormatRecommendations(recs []*Recommendation) string {
    // Format as numbered list with commands
    // Highlight security warnings
    // Include explanations
}
```

#### Component 3: Insecure Mode Implementation
**Purpose**: Implement --insecure flag for development environments  
**Location**: `pkg/fallback/insecure.go`  
**Lines**: ~60  

```go
// InsecureConfig manages insecure mode settings
type InsecureConfig struct {
    Enabled      bool
    WarningShown bool
    AuditLog     []string
}

// CreateInsecureTLSConfig creates TLS config that skips verification
func CreateInsecureTLSConfig() *tls.Config {
    return &tls.Config{
        InsecureSkipVerify: true,
    }
}

// ShowInsecureWarning displays prominent security warning
func ShowInsecureWarning() {
    fmt.Println("===============================================")
    fmt.Println("WARNING: TLS VERIFICATION DISABLED")
    fmt.Println("===============================================")
    fmt.Println("You are using --insecure flag which disables")
    fmt.Println("TLS certificate verification. This is:")
    fmt.Println("- DANGEROUS in production")
    fmt.Println("- Only for development/testing")
    fmt.Println("- Vulnerable to man-in-the-middle attacks")
    fmt.Println("===============================================")
}

// ApplyInsecureFlag updates configuration based on flag
func ApplyInsecureFlag(config *InsecureConfig, flagValue bool) error {
    if flagValue {
        config.Enabled = true
        if !config.WarningShown {
            ShowInsecureWarning()
            config.WarningShown = true
        }
        config.AuditLog = append(config.AuditLog, 
            fmt.Sprintf("Insecure mode enabled at %v", time.Now()))
    }
    return nil
}

// WrapTransportWithInsecure wraps HTTP transport for insecure mode
func WrapTransportWithInsecure(transport http.RoundTripper, insecure bool) http.RoundTripper {
    if insecure {
        if t, ok := transport.(*http.Transport); ok {
            t.TLSClientConfig = CreateInsecureTLSConfig()
        }
    }
    return transport
}
```

#### Component 4: Certificate Chain Logger
**Purpose**: Log detailed certificate information for debugging  
**Location**: `pkg/fallback/logger.go`  
**Lines**: ~40  

```go
// LogCertificateChain logs detailed certificate chain information
func LogCertificateChain(certs []*x509.Certificate, logger *log.Logger) {
    for i, cert := range certs {
        logger.Printf("Certificate %d of %d:", i+1, len(certs))
        logger.Printf("  Subject: %s", cert.Subject)
        logger.Printf("  Issuer: %s", cert.Issuer)
        logger.Printf("  Serial: %s", cert.SerialNumber)
        logger.Printf("  NotBefore: %v", cert.NotBefore)
        logger.Printf("  NotAfter: %v", cert.NotAfter)
        logger.Printf("  DNS Names: %v", cert.DNSNames)
        logger.Printf("  IP Addresses: %v", cert.IPAddresses)
        logger.Printf("  Is CA: %v", cert.IsCA)
        
        // Check signature algorithm
        logger.Printf("  Signature Algorithm: %v", cert.SignatureAlgorithm)
        
        // Log any certificate extensions
        if len(cert.Extensions) > 0 {
            logger.Printf("  Extensions: %d", len(cert.Extensions))
        }
    }
}

// LogValidationError logs structured validation error details
func LogValidationError(err error, cert *x509.Certificate, logger *log.Logger) {
    logger.Printf("Certificate validation failed:")
    logger.Printf("  Error Type: %T", err)
    logger.Printf("  Error Message: %v", err)
    
    // Log certificate details for context
    if cert != nil {
        logger.Printf("  Certificate Subject: %s", cert.Subject)
        logger.Printf("  Certificate Issuer: %s", cert.Issuer)
    }
}
```

## Testing Strategy

### Unit Tests Required
- [ ] Test file: `pkg/fallback/detector_test.go`
- [ ] Test file: `pkg/fallback/recommender_test.go`
- [ ] Test file: `pkg/fallback/insecure_test.go`
- [ ] Coverage target: 80%
- [ ] Test cases:
  - [ ] Detect self-signed certificate errors
  - [ ] Detect expired certificate errors
  - [ ] Detect hostname mismatch errors
  - [ ] Detect untrusted CA errors
  - [ ] Generate recommendations for each problem type
  - [ ] Format recommendations for user output
  - [ ] Create insecure TLS config
  - [ ] Display insecure mode warnings
  - [ ] Log certificate chain details
  - [ ] Handle unknown error types gracefully

### Test Fixtures
```go
// Use test certificates from E1.2.1
import (
    "github.com/idpbuilder/idpbuilder-oci-go-cr/pkg/certs/testdata"
)

// Additional test cases:
// - Error scenarios from real registry interactions
// - Various x509 error types
// - Edge cases in certificate chains
```

### Integration Tests
- [ ] Integration with CertValidator from E1.2.1
- [ ] End-to-end fallback flow with recommendations
- [ ] --insecure flag with actual registry connections

## Integration Points

### With E1.2.1 (certificate-validation-pipeline)
- Import CertValidator interface for validation
- Import ValidationError types for error analysis
- Import CertDiagnostics for detailed problem information
- Use test certificates from E1.2.1's testdata

### With Future E2.1.2 (gitea-registry-client)
- Provide ProblemDetector for analyzing push failures
- Provide Recommender for user guidance
- Provide InsecureConfig for --insecure mode
- Integrate with registry client's TLS configuration

## Size Management Strategy

### Measurement Protocol
```bash
# Measure using project line counter (R200 compliance)
cd /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/wave2/fallback-strategies
$PROJECT_ROOT/tools/line-counter.sh

# Check points:
# - After detector.go implementation (~100 lines)
# - After recommender.go implementation (~180 lines)
# - After insecure.go implementation (~240 lines)
# - After logger.go implementation (~280 lines)
# - After all tests implementation (~400 lines)
```

### Split Prevention
- Core detection logic kept focused (~100 lines)
- Recommendations separated from detection (~80 lines)
- Insecure mode isolated (~60 lines)
- Logging utilities minimal (~40 lines)
- Total well under 800 line limit with 400 line buffer

## Implementation Sequence

1. **Import Dependencies from E1.2.1** (15 min)
   - Import CertValidator interface
   - Import ValidationError types
   - Import CertDiagnostics struct

2. **Implement Problem Detection** (2 hours)
   - Create ProblemDetector interface
   - Implement error pattern matching
   - Extract problem-specific details

3. **Build Recommendation Engine** (1.5 hours)
   - Create Recommender interface
   - Implement problem-specific recommendations
   - Format user-friendly output

4. **Add Insecure Mode** (1 hour)
   - Implement --insecure flag handling
   - Create insecure TLS configuration
   - Add prominent warnings

5. **Implement Logging Utilities** (30 min)
   - Certificate chain logger
   - Validation error logger
   - Debug output formatting

6. **Create Test Suite** (2 hours)
   - Unit tests for all components
   - Integration tests with E1.2.1
   - Edge case coverage

7. **Documentation** (30 min)
   - Code comments
   - Usage examples
   - Security warnings

## Notes for SW Engineer

### Critical Considerations
1. **Import from E1.2.1**: Must import and use CertValidator interface and types
2. **Security Warnings**: --insecure mode MUST show prominent warnings
3. **Error Analysis**: Use type assertions to identify specific x509 error types
4. **Actionable Output**: Every recommendation must include a concrete action
5. **Audit Trail**: Log when insecure mode is enabled for security auditing

### Example Usage
```go
// Import validator from E1.2.1
import (
    "github.com/idpbuilder/idpbuilder-oci-go-cr/pkg/certs"
    "github.com/idpbuilder/idpbuilder-oci-go-cr/pkg/fallback"
)

// Detect and handle certificate problems
validator := certs.NewValidator(trustStore)
cert := getRegistryCertificate()

if err := validator.ValidateChain(cert); err != nil {
    // Detect the specific problem
    detector := fallback.NewDetector(validator)
    problem, _ := detector.DetectProblem(err, cert)
    
    // Generate recommendations
    recommender := fallback.NewRecommender(registryURL, true)
    recommendations, _ := recommender.Recommend(problem)
    
    // Display to user
    fmt.Println(recommender.FormatRecommendations(recommendations))
    
    // If --insecure flag is set
    if insecureFlag {
        fallback.ShowInsecureWarning()
        tlsConfig = fallback.CreateInsecureTLSConfig()
    }
}
```

### Dependencies to Import
```go
import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "log"
    "net/http"
    "time"
    
    // From E1.2.1
    "github.com/idpbuilder/idpbuilder-oci-go-cr/pkg/certs"
)
```

### CLI Flag Integration
```go
// Add to CLI command definitions
var insecureFlag = &cli.BoolFlag{
    Name:    "insecure",
    Aliases: []string{"k"},
    Usage:   "Skip TLS certificate verification (INSECURE - development only)",
    EnvVars: []string{"IDPBUILDER_INSECURE"},
}

// Check flag in command action
if c.Bool("insecure") {
    fallback.ShowInsecureWarning()
    // Configure HTTP client with insecure TLS
}
```

## Review Checklist
- [ ] All functional requirements addressed
- [ ] Size within 400 line estimate
- [ ] Test coverage plan adequate
- [ ] Integration with E1.2.1 clear
- [ ] Security warnings prominent
- [ ] Recommendations actionable
- [ ] Implementation sequence logical
- [ ] Documentation complete