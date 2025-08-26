# Effort 4: Image Security & Signing Implementation Plan

## =¨ CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `idpbuidler-oci-mgmt/phase2/wave2/effort4-security`  
**Can Parallelize**: Yes  
**Parallel With**: Efforts 2, 3  
**Size Estimate**: 650 lines (MUST be <800)  
**Dependencies**: Effort 1 (contracts)  

## Overview
- **Effort**: Implement SecurityManager interface for image signing, verification, SBOM generation, and vulnerability scanning
- **Phase**: 2, Wave: 2
- **Estimated Size**: 650 lines
- **Implementation Time**: 8 hours
- **Key Features**: Cosign-based signing, SBOM generation, vulnerability scanner hooks

## = Dependencies Analysis (R219 Compliance)

### Dependencies from Effort 1 (Contracts)
**Location**: `../effort1-contracts/pkg/oci/api/security.go`
**What We Import**:
```go
import "github.com/cnoe-io/idpbuilder/pkg/oci/api"

// Implementing these interfaces:
- api.SecurityManager  // Main security operations interface
- api.Signer           // Digital signing interface
- api.Verifier         // Signature verification interface

// Using these models:
- api.Signature, api.SignatureBundle, api.SignatureData
- api.SBOM, api.Component, api.Tool
- api.VulnerabilityReport, api.Vulnerability
- api.Attestation, api.Policy, api.SecurityProfile
- api.Certificate
```

### How Dependencies Influence Implementation
1. **Interface Compliance**: Must implement ALL methods from SecurityManager interface
2. **Model Reuse**: Use api package models directly, no duplication
3. **Pattern Following**: Follow error handling and context patterns from contracts
4. **Integration Points**: Our implementation will be used by Effort 5 (Registry) for auto-sign/verify

## File Structure

### Implementation Files (650 lines total)
- `pkg/oci/security/manager.go`: SecurityManager implementation (150 lines)
- `pkg/oci/security/signer.go`: Cosign-based signing implementation (120 lines)
- `pkg/oci/security/verifier.go`: Signature verification logic (100 lines)
- `pkg/oci/security/sbom.go`: SBOM generation with SPDX/CycloneDX (150 lines)
- `pkg/oci/security/scanner_hooks.go`: Vulnerability scanner integration (80 lines)
- `pkg/oci/security/attestation.go`: Build attestation support (50 lines)

### Test Files
- `pkg/oci/security/manager_test.go`: SecurityManager tests
- `pkg/oci/security/signer_test.go`: Signing tests with mock keys
- `pkg/oci/security/verifier_test.go`: Verification tests
- `pkg/oci/security/sbom_test.go`: SBOM generation tests
- `pkg/oci/security/testdata/`: Test keys, certificates, sample SBOMs

## Implementation Steps

### Step 1: Create Package Structure (15 minutes)
```bash
# Create security package directory
mkdir -p pkg/oci/security
mkdir -p pkg/oci/security/testdata

# Create implementation files
touch pkg/oci/security/{manager,signer,verifier,sbom,scanner_hooks,attestation}.go
touch pkg/oci/security/{manager,signer,verifier,sbom}_test.go

# Create test data files
touch pkg/oci/security/testdata/{test-key.pem,test-cert.pem}
```

### Step 2: Implement SecurityManager Core (2 hours)
```go
// pkg/oci/security/manager.go
package security

import (
    "context"
    "fmt"
    
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// securityManager implements api.SecurityManager interface
type securityManager struct {
    signers    map[string]api.Signer
    verifiers  map[string]api.Verifier
    sbomGen    *sbomGenerator
    scanners   []ScannerPlugin
}

// NewSecurityManager creates a new security manager instance
func NewSecurityManager(opts ...Option) api.SecurityManager {
    sm := &securityManager{
        signers:   make(map[string]api.Signer),
        verifiers: make(map[string]api.Verifier),
        sbomGen:   newSBOMGenerator(),
    }
    
    for _, opt := range opts {
        opt(sm)
    }
    
    return sm
}

// SignImage signs an image using the provided signer
func (sm *securityManager) SignImage(ctx context.Context, image string, signer api.Signer) (*api.Signature, error) {
    // Implementation with Cosign integration
    // 1. Get image manifest
    // 2. Create payload
    // 3. Sign with signer
    // 4. Store signature
}

// VerifySignature verifies an image signature
func (sm *securityManager) VerifySignature(ctx context.Context, image string, verifier api.Verifier) error {
    // Implementation
}

// GenerateSBOM creates Software Bill of Materials
func (sm *securityManager) GenerateSBOM(ctx context.Context, image string) (*api.SBOM, error) {
    return sm.sbomGen.Generate(ctx, image)
}

// Additional methods...
```

### Step 3: Implement Cosign-based Signer (1.5 hours)
```go
// pkg/oci/security/signer.go
package security

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// cosignSigner implements api.Signer using Cosign-compatible signing
type cosignSigner struct {
    privateKey crypto.PrivateKey
    keyID      string
    algorithm  string
    certChain  []*api.Certificate
}

// NewCosignSigner creates a new Cosign-compatible signer
func NewCosignSigner(keyPath string, opts ...SignerOption) (api.Signer, error) {
    // Load private key
    // Parse certificate chain if provided
    // Return configured signer
}

// NewKeylessSigner creates a keyless (OIDC-based) signer
func NewKeylessSigner(provider string) (api.Signer, error) {
    // Configure OIDC provider
    // Get identity token
    // Return keyless signer
}

// Sign creates a digital signature
func (cs *cosignSigner) Sign(payload []byte) ([]byte, error) {
    hash := sha256.Sum256(payload)
    
    switch key := cs.privateKey.(type) {
    case *rsa.PrivateKey:
        return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
    default:
        return nil, fmt.Errorf("unsupported key type")
    }
}

// Additional signer methods...
```

### Step 4: Implement Signature Verifier (1.5 hours)
```go
// pkg/oci/security/verifier.go
package security

import (
    "context"
    "crypto"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// cosignVerifier implements api.Verifier
type cosignVerifier struct {
    publicKeys  map[string]crypto.PublicKey
    trustedCAs  *x509.CertPool
    trustedKeys []string
    policy      *api.Policy
}

// NewCosignVerifier creates a new verifier
func NewCosignVerifier(pubKeyPath string, opts ...VerifierOption) (api.Verifier, error) {
    // Load public keys
    // Configure trust roots
    // Return verifier
}

// Verify validates a signature
func (cv *cosignVerifier) Verify(payload []byte, signature []byte) error {
    hash := sha256.Sum256(payload)
    
    for _, pubKey := range cv.publicKeys {
        switch key := pubKey.(type) {
        case *rsa.PublicKey:
            err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hash[:], signature)
            if err == nil {
                return nil
            }
        }
    }
    
    return fmt.Errorf("signature verification failed")
}

// Additional verifier methods...
```

### Step 5: Implement SBOM Generation (2 hours)
```go
// pkg/oci/security/sbom.go
package security

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// sbomGenerator handles SBOM creation in multiple formats
type sbomGenerator struct {
    format     SBOMFormat
    scanners   []ComponentScanner
    enrichers  []SBOMEnricher
}

// SBOMFormat represents supported SBOM formats
type SBOMFormat string

const (
    FormatSPDX      SBOMFormat = "spdx"
    FormatCycloneDX SBOMFormat = "cyclonedx"
)

// Generate creates an SBOM for an image
func (sg *sbomGenerator) Generate(ctx context.Context, image string) (*api.SBOM, error) {
    sbom := &api.SBOM{
        Version:   "1.0",
        Timestamp: time.Now(),
        Image:     image,
        Components: []*api.Component{},
        Dependencies: make(map[string][]string),
    }
    
    // Scan image layers for components
    components, err := sg.scanImageComponents(ctx, image)
    if err != nil {
        return nil, fmt.Errorf("scanning components: %w", err)
    }
    
    sbom.Components = components
    
    // Build dependency graph
    sbom.Dependencies = sg.buildDependencyGraph(components)
    
    // Enrich with additional metadata
    for _, enricher := range sg.enrichers {
        enricher.Enrich(sbom)
    }
    
    return sbom, nil
}

// scanImageComponents identifies software components in the image
func (sg *sbomGenerator) scanImageComponents(ctx context.Context, image string) ([]*api.Component, error) {
    var components []*api.Component
    
    // Analyze package managers (apt, yum, npm, pip, go modules, etc.)
    // Extract version information
    // Identify licenses
    
    return components, nil
}

// FormatSPDX converts SBOM to SPDX format
func (sg *sbomGenerator) FormatSPDX(sbom *api.SBOM) ([]byte, error) {
    // Convert to SPDX JSON format
    return json.MarshalIndent(sbom, "", "  ")
}

// FormatCycloneDX converts SBOM to CycloneDX format
func (sg *sbomGenerator) FormatCycloneDX(sbom *api.SBOM) ([]byte, error) {
    // Convert to CycloneDX JSON format
    return json.MarshalIndent(sbom, "", "  ")
}
```

### Step 6: Implement Scanner Integration Hooks (1 hour)
```go
// pkg/oci/security/scanner_hooks.go
package security

import (
    "context"
    "fmt"
    
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// ScannerPlugin defines the interface for vulnerability scanner plugins
type ScannerPlugin interface {
    Name() string
    Version() string
    Scan(ctx context.Context, image string) (*api.VulnerabilityReport, error)
}

// TrivyScanner implements ScannerPlugin for Trivy integration
type TrivyScanner struct {
    executable string
    dbPath     string
}

// NewTrivyScanner creates a Trivy scanner plugin
func NewTrivyScanner(opts ...ScannerOption) ScannerPlugin {
    return &TrivyScanner{
        executable: "trivy",
        dbPath:     "/tmp/trivy-db",
    }
}

// Scan performs vulnerability scanning using Trivy
func (ts *TrivyScanner) Scan(ctx context.Context, image string) (*api.VulnerabilityReport, error) {
    // Execute trivy command
    // Parse JSON output
    // Convert to api.VulnerabilityReport
    
    report := &api.VulnerabilityReport{
        Timestamp: time.Now(),
        Image:     image,
        Scanner: &api.ScannerInfo{
            Name:    ts.Name(),
            Version: ts.Version(),
        },
        Summary: &api.VulnerabilitySummary{},
        Vulnerabilities: []*api.Vulnerability{},
    }
    
    // Populate report from scanner output
    
    return report, nil
}

// GrypeScanner implements ScannerPlugin for Grype integration
type GrypeScanner struct {
    // Grype configuration
}

// Additional scanner implementations...
```

### Step 7: Implement Attestation Support (45 minutes)
```go
// pkg/oci/security/attestation.go
package security

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// AttestationBuilder creates and verifies attestations
type AttestationBuilder struct {
    predicateTypes map[string]PredicateBuilder
}

// PredicateBuilder creates attestation predicates
type PredicateBuilder interface {
    Build(subject string, data interface{}) (map[string]interface{}, error)
}

// CreateSLSAAttestation creates a SLSA provenance attestation
func CreateSLSAAttestation(image string, buildInfo *BuildInfo) (*api.Attestation, error) {
    attestation := &api.Attestation{
        Type:      "https://slsa.dev/provenance/v0.2",
        Subject:   image,
        Timestamp: time.Now(),
        Predicate: map[string]interface{}{
            "builder": map[string]interface{}{
                "id": buildInfo.BuilderID,
            },
            "buildType": buildInfo.BuildType,
            "invocation": map[string]interface{}{
                "configSource": buildInfo.ConfigSource,
            },
        },
    }
    
    return attestation, nil
}

// CreateVulnAttestation creates a vulnerability scan attestation
func CreateVulnAttestation(image string, report *api.VulnerabilityReport) (*api.Attestation, error) {
    attestation := &api.Attestation{
        Type:      "https://cosign.dev/vulnerability/v1",
        Subject:   image,
        Timestamp: time.Now(),
        Predicate: map[string]interface{}{
            "scanner":         report.Scanner.Name,
            "scanner_version": report.Scanner.Version,
            "summary":         report.Summary,
            "timestamp":       report.Timestamp,
        },
    }
    
    return attestation, nil
}

// BuildInfo contains build metadata for attestations
type BuildInfo struct {
    BuilderID    string
    BuildType    string
    ConfigSource map[string]interface{}
}
```

### Step 8: Create Comprehensive Tests (2 hours)
```go
// pkg/oci/security/manager_test.go
package security

import (
    "context"
    "testing"
    
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

func TestSecurityManager_SignImage(t *testing.T) {
    // Test signing with different key types
    // Test keyless signing
    // Test error cases
}

func TestSecurityManager_VerifySignature(t *testing.T) {
    // Test valid signatures
    // Test invalid signatures
    // Test trust chain validation
}

func TestSecurityManager_GenerateSBOM(t *testing.T) {
    // Test SBOM generation
    // Test different formats
    // Test component detection
}

// Additional test functions...
```

## Size Management
- **Estimated Lines**: 650 lines
- **Measurement Tool**: /home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh
- **Check Frequency**: Every 100 lines during implementation
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

### Size Breakdown by File
| File | Estimated Lines | Purpose |
|------|-----------------|---------|
| manager.go | 150 | Core SecurityManager implementation |
| signer.go | 120 | Cosign signing logic |
| verifier.go | 100 | Signature verification |
| sbom.go | 150 | SBOM generation |
| scanner_hooks.go | 80 | Scanner plugins |
| attestation.go | 50 | Attestation support |
| **Total** | **650** | Under limit |

## Test Requirements

### Unit Tests
- **Coverage Target**: 85% minimum
- **Test Cases**:
  - Signing with RSA, ECDSA, ED25519 keys
  - Keyless signing flow
  - Signature verification (valid/invalid)
  - SBOM generation for different image types
  - Scanner plugin execution
  - Attestation creation and verification
  - Policy enforcement

### Integration Tests
- Test with real Cosign signatures
- Verify against public registries
- Test SBOM formats (SPDX, CycloneDX)
- Scanner integration (if available)

### Test Data
- Generate test keys and certificates
- Create sample SBOMs
- Mock vulnerability reports
- Example attestations

## Pattern Compliance

### Security Best Practices
- Never log private keys or secrets
- Use secure random for cryptographic operations
- Validate all inputs before processing
- Handle certificate expiration gracefully
- Support key rotation

### Go Patterns
- Use context for cancellation
- Return wrapped errors with context
- Use interfaces for extensibility
- Implement options pattern for configuration
- Follow Go naming conventions

### Interface Compliance
```go
// Ensure we implement the interface
var _ api.SecurityManager = (*securityManager)(nil)
var _ api.Signer = (*cosignSigner)(nil)
var _ api.Verifier = (*cosignVerifier)(nil)
```

## External Dependencies

### Required Libraries
```go
// Note: Only wrapper code counts toward our line limit
import (
    "github.com/sigstore/cosign/v2/pkg/cosign"     // For signing/verification
    "github.com/spdx/tools-golang/spdx"            // For SPDX SBOM format
    "github.com/CycloneDX/cyclonedx-go"            // For CycloneDX format
)
```

### Optional Scanner Integrations
- Trivy: via command execution
- Grype: via command execution
- Custom scanners: via plugin interface

## Integration Points

### With Wave 1
- May use build metadata for attestations
- Can reference image build configurations

### With Other Wave 2 Efforts
- **Effort 1**: Import all interfaces and models
- **Effort 5**: Will use our SecurityManager for:
  - Auto-signing on push
  - Auto-verification on pull
  - SBOM attachment to images

## Risk Mitigation

### Risk 1: Complex Cryptographic Operations
**Mitigation**: Use well-tested libraries (Cosign), extensive testing

### Risk 2: Scanner Integration Complexity
**Mitigation**: Plugin architecture, start with basic hooks

### Risk 3: Size Limit Exceeded
**Mitigation**: Focus on core functionality, defer advanced features

### Risk 4: SBOM Format Complexity
**Mitigation**: Start with basic component detection, iterate

## Success Criteria

### Must Have
- [ ] SecurityManager interface fully implemented
- [ ] Cosign-based signing working
- [ ] Signature verification functional
- [ ] Basic SBOM generation
- [ ] Scanner plugin interface defined
- [ ] All tests passing

### Should Have
- [ ] Keyless signing support
- [ ] SPDX and CycloneDX formats
- [ ] Trivy scanner integration
- [ ] Attestation support

### Could Have
- [ ] Multiple scanner plugins
- [ ] Advanced SBOM enrichment
- [ ] Policy enforcement engine

## Next Steps
1. Create package structure
2. Import Effort 1 contracts
3. Implement SecurityManager core
4. Add signing/verification
5. Implement SBOM generation
6. Create scanner hooks
7. Write comprehensive tests
8. Verify size compliance with line-counter.sh

## Document History
| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-08-26 | Code Reviewer Agent | Initial implementation plan for security & signing |