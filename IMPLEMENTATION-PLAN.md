# E3.1.5 Certificate Integration & Testing Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: idpbuidler-oci-mgmt/phase3/wave1/E3.1.5-integration-layer  
**Can Parallelize**: Yes (after E3.1.1)  
**Parallel With**: [E3.1.2, E3.1.3, E3.1.4]  
**Size Estimate**: 650 lines  
**Dependencies**: [E3.1.1, Phase 2 registry client]  
**Working Directory**: /home/vscode/workspaces/idpbuilder-oci-mgmt/efforts/phase3/wave1/E3.1.5-integration-layer

## Overview
- **Effort**: Certificate Integration & Testing - Complete integration with existing systems and comprehensive testing infrastructure
- **Phase**: 3, Wave: 1
- **Estimated Size**: 650 lines total
- **Implementation Time**: 6 hours
- **Purpose**: Bridge certificate management with existing registry/build systems and provide comprehensive testing

## Dependency Context (R219 Compliance)

### Dependencies Analyzed
1. **E3.1.1 (Certificate Contracts)**: Provides core interfaces we must implement mocks for
   - CertificateService interface
   - CertFormat and VerificationMode types
   - CertBundle and Certificate structures
   
2. **Phase 2 Registry Client**: Not yet available - will create mock registry for testing
   - Need to simulate registry push/pull operations
   - Certificate injection into registry authentication
   - TLS configuration for secure connections

### Integration Strategy
- Import v2 API contracts from E3.1.1
- Create comprehensive mocks for all certificate services
- Simulate Phase 2 registry behavior for testing
- Build test harnesses for certificate validation flows

## File Structure
```
pkg/oci/certificates/
├── integration.go (200 lines)
│   └── Registry client integration points
│   └── Build service integration hooks
│   └── Certificate injection middleware
│   └── TLS configuration helpers
├── mocks.go (150 lines)
│   └── MockCertificateService implementation
│   └── Mock certificate loaders
│   └── Test certificate generator
│   └── Mock registry client
└── integration_test.go (300 lines)
    └── End-to-end certificate flow tests
    └── Registry push/pull with certificates
    └── Build process with custom CA
    └── Certificate rotation scenarios
    └── Error recovery testing
```

## Implementation Steps

### Step 1: Create Integration Layer (200 lines)
**File**: `pkg/oci/certificates/integration.go`

**Key Components**:
1. **RegistryIntegration struct**
   - Wraps registry client with certificate support
   - Injects certificates into authentication flow
   - Manages TLS configuration for connections
   
2. **BuildIntegration struct**
   - Integrates certificates with build process
   - Configures buildah with custom CAs
   - Handles certificate validation during builds
   
3. **CertificateMiddleware interface**
   - Pluggable certificate injection points
   - HTTP transport configuration
   - Authentication enrichment

**Implementation Details**:
```go
// Core integration types
type RegistryIntegration struct {
    certService CertificateService
    tlsConfig   *tls.Config
    registry    RegistryClient // Will be mocked
}

type BuildIntegration struct {
    certService CertificateService
    buildConfig BuildConfiguration
}

// Methods to implement:
- InjectCertificates(transport http.RoundTripper) http.RoundTripper
- ConfigureTLS(config *tls.Config) error
- ValidateRegistryCertificate(url string) error
- SetupBuildCertificates(buildah BuildahConfig) error
```

### Step 2: Create Mock Implementations (150 lines)
**File**: `pkg/oci/certificates/mocks.go`

**Mock Components**:
1. **MockCertificateService**
   - Implements full CertificateService interface
   - Configurable responses for testing
   - Error simulation capabilities
   
2. **TestCertificateGenerator**
   - Creates valid test certificates
   - Generates certificate chains
   - Produces expired/invalid certs for testing
   
3. **MockRegistryClient**
   - Simulates Phase 2 registry behavior
   - Certificate validation hooks
   - Push/pull operation mocks

**Key Features**:
```go
// Mock service with configurable behavior
type MockCertificateService struct {
    LoadBundleFunc     func(ctx, path, format) (*CertBundle, error)
    ValidateFunc       func(cert) error
    VerificationMode   VerificationMode
    ErrorOnNthCall     int // For testing retry logic
}

// Test certificate generation utilities
func GenerateTestCA() (*x509.Certificate, crypto.PrivateKey, error)
func GenerateTestCertificate(ca *x509.Certificate) (*x509.Certificate, error)
func GenerateExpiredCertificate() (*x509.Certificate, error)
func GenerateSelfSignedCertificate() (*x509.Certificate, error)
```

### Step 3: Implement Comprehensive Test Suite (300 lines)
**File**: `pkg/oci/certificates/integration_test.go`

**Test Categories**:

1. **Certificate Loading Tests** (60 lines)
   - Test all supported formats (PEM, DER, PKCS7, PKCS12)
   - Invalid format handling
   - Corrupted certificate handling
   - Large bundle performance

2. **Registry Integration Tests** (80 lines)
   - Push with custom CA certificate
   - Pull with certificate validation
   - Registry authentication with client certs
   - TLS handshake verification
   - Certificate rotation during operations

3. **Build Integration Tests** (60 lines)
   - Build with custom CA bundle
   - Multi-stage build certificate propagation
   - Build cache with certificates
   - Certificate validation in FROM statements

4. **End-to-End Scenarios** (60 lines)
   - Complete workflow: load → inject → push → pull
   - Certificate rotation without downtime
   - Fallback to skip-verify on failure
   - Concurrent operations with shared certificates

5. **Error Recovery Tests** (40 lines)
   - Invalid certificate recovery
   - Network timeout handling
   - Certificate expiry detection
   - CA pool corruption recovery

**Test Structure Example**:
```go
func TestCertificateIntegrationE2E(t *testing.T) {
    // Setup
    certService := setupMockCertificateService()
    registry := setupMockRegistry()
    
    // Test certificate loading
    bundle, err := certService.LoadCertificateBundle(ctx, "test.pem", CertFormatPEM)
    require.NoError(t, err)
    
    // Test registry integration
    integration := NewRegistryIntegration(certService, registry)
    err = integration.Push(ctx, "test-image", bundle)
    require.NoError(t, err)
    
    // Verify certificate was used
    assert.True(t, registry.CertificateValidated())
}
```

## Size Management
- **Estimated Lines**: 650 total
  - integration.go: 200 lines
  - mocks.go: 150 lines
  - integration_test.go: 300 lines
- **Measurement Tool**: ${PROJECT_ROOT}/tools/line-counter.sh
- **Check Frequency**: After each file completion
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

## Test Requirements

### Coverage Targets
- **Unit Tests**: 85% coverage minimum
- **Integration Tests**: All integration points tested
- **E2E Tests**: 5 complete scenarios minimum
- **Mock Coverage**: 100% of external dependencies

### Test Files
1. `integration_test.go`: Main test suite
2. Test fixtures in `testdata/`:
   - Sample certificates in various formats
   - Invalid/corrupted certificates
   - Certificate chains
   - Client certificates

### Testing Strategy
1. Use table-driven tests for format variations
2. Parallel test execution where possible
3. Benchmark critical paths
4. Fuzz testing for certificate parsing
5. Property-based testing for validation logic

## Pattern Compliance

### idpbuilder-oci-mgmt Patterns
- Follow existing error handling patterns
- Use context for cancellation
- Implement proper logging with structured fields
- Follow interface segregation principle

### Security Requirements
- Never log certificate private keys
- Validate all certificate inputs
- Secure certificate storage in tests
- Clear sensitive data after tests

### Performance Targets
- Certificate loading < 100ms
- TLS handshake overhead < 50ms
- Mock operations < 10ms
- Test suite completion < 30s

## Integration Points

### With E3.1.1 (Certificate Contracts)
```go
import "github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"

// Use interfaces from E3.1.1
var certService v2.CertificateService
```

### With Other E3.1.x Efforts
- E3.1.2: Test certificate loading implementations
- E3.1.3: Validate service integration points
- E3.1.4: Test storage and management flows

### With Phase 2 Components (Mocked)
```go
// Mock Phase 2 registry client
type MockRegistryClient struct {
    v1.RegistryClient // Would import from Phase 2 if available
}
```

## Risk Management

### Identified Risks
1. **Phase 2 Integration**: Registry client not available
   - Mitigation: Comprehensive mocks based on expected interface
   
2. **Certificate Format Complexity**: Many formats to support
   - Mitigation: Focus on PEM/DER first, others in mocks
   
3. **Test Certificate Generation**: Complex certificate chains
   - Mitigation: Use existing crypto libraries, pre-generate fixtures

### Contingency Plans
- If size exceeds 650 lines: Move some tests to separate test files
- If Phase 2 changes: Update mocks to match new interfaces
- If certificate complexity grows: Split mocks into separate files

## Success Criteria

### Functional Success
- ✅ All certificate formats handled in tests
- ✅ Registry integration fully mocked
- ✅ Build integration patterns established
- ✅ Comprehensive test coverage achieved

### Quality Metrics
- ✅ Zero test flakiness
- ✅ All edge cases covered
- ✅ Performance benchmarks established
- ✅ Security best practices followed

### Deliverables
1. Complete integration layer with clear interfaces
2. Comprehensive mock implementations
3. Full test suite with >85% coverage
4. Performance benchmarks for critical paths
5. Test fixtures for all certificate scenarios

## Implementation Notes

### For Software Engineers
1. Start with integration.go to define clear boundaries
2. Build mocks incrementally as needed for tests
3. Use test-driven development for complex scenarios
4. Focus on realistic test cases from production
5. Document all mock behaviors clearly

### Critical Implementation Order
1. First: Import and understand E3.1.1 contracts
2. Second: Create basic mock certificate service
3. Third: Build integration layer with clear interfaces
4. Fourth: Expand mocks as test cases require
5. Last: Complete comprehensive test suite

### Testing Best Practices
- Keep test data in testdata/ directory
- Use subtests for better organization
- Clean up resources in test teardown
- Mock time for expiry testing
- Use golden files for complex validations

## Appendix

### Import Structure
```go
package certificates

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "testing"
    
    // Phase 3 contracts
    v2 "github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"
    
    // Standard testing libraries
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/mock"
)
```

### Sample Test Certificate Generation
```go
func GenerateTestCertificateChain() ([]*x509.Certificate, error) {
    // Root CA
    rootCA := generateRootCA()
    
    // Intermediate CA
    intermediateCA := generateIntermediateCA(rootCA)
    
    // Leaf certificate
    leafCert := generateLeafCertificate(intermediateCA)
    
    return []*x509.Certificate{leafCert, intermediateCA, rootCA}, nil
}
```

### Mock Configuration Example
```go
mockService := &MockCertificateService{
    LoadBundleFunc: func(ctx context.Context, path string, format CertFormat) (*CertBundle, error) {
        if format == CertFormatPEM {
            return testBundle, nil
        }
        return nil, ErrUnsupportedFormat
    },
    VerificationMode: VerificationModeStrict,
}
```

---

**Document Version**: 1.0  
**Created**: 2025-08-27  
**Author**: Code Reviewer Agent  
**State**: EFFORT_PLAN_CREATION  
**Review Status**: Ready for Implementation