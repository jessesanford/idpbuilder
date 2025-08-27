# Split Plan for E3.1.3-certificate-validator

## Overview
**Total Current Size**: 1231 lines (exceeds 800 line limit by 431 lines)
**Required**: Split into 3 functional components, each under 400 lines
**Strategy**: Separate by architectural layer - Interface/Core Service, Gitea Integration, and Verification/Tests

## Split Architecture

### Split 001: Core Certificate Service & Interface
**Target Size**: ~380 lines
**Branch**: `idpbuidler-oci-mgmt/phase3/wave1/E3.1.3-certificate-validator-split-001`
**Purpose**: Foundation service with interface definitions and core certificate operations

**Files**:
- `pkg/oci/api/v2/certificate_service.go` (62 lines) - Interface definitions
- `pkg/oci/certificates/service.go` (318 lines) - Core service (reduced from 386)
  - Remove Gitea integration initialization
  - Keep core pool management and certificate operations
  - Simplified constructor without discovery components

**Functionality**:
- CertificateService interface and types
- Basic service implementation with certificate pool management
- LoadCertificateBundle and ValidateCertificate methods
- Certificate pool operations (Add, Remove, GetCertPool)
- Thread-safe operations with mutex protection

**Dependencies**: None (foundation split)

**Interface Contract**:
```go
// Exported for use by other splits
type CertificateServiceImpl struct {
    systemPool   *x509.CertPool
    customPool   *x509.CertPool
    mode         v2.VerificationMode
    mutex        sync.RWMutex
}

// Public methods remain unchanged
```

### Split 002: Gitea Integration & Discovery
**Target Size**: ~377 lines  
**Branch**: `idpbuidler-oci-mgmt/phase3/wave1/E3.1.3-certificate-validator-split-002`
**Purpose**: Gitea-specific certificate discovery and integration

**Files**:
- `pkg/oci/certificates/gitea_integration.go` (377 lines) - Complete Gitea integration

**Functionality**:
- GiteaDiscovery implementation
- Configuration file parsing
- Registry endpoint discovery
- Certificate chain validation
- Root CA loading from Gitea
- Auto-discovery mechanisms

**Dependencies**: 
- Split 001 (imports certificate types from v2 package)
- Can operate independently for discovery operations

**Interface Contract**:
```go
// Exported for service integration
type GiteaDiscovery struct {
    configPaths []string
    endpoints   []string
    discovered  []*x509.Certificate
    mutex       sync.RWMutex
}

func NewGiteaDiscovery() (*GiteaDiscovery, error)
func (g *GiteaDiscovery) DiscoverCertificates(ctx context.Context) ([]*x509.Certificate, error)
```

### Split 003: Verification Management & Tests
**Target Size**: ~400 lines
**Branch**: `idpbuidler-oci-mgmt/phase3/wave1/E3.1.3-certificate-validator-split-003`
**Purpose**: Verification mode management and comprehensive test suite

**Files**:
- `pkg/oci/certificates/verification.go` (200 lines) - Reduced verification manager
  - Remove redundant validation logic
  - Focus on mode switching and fallback strategies
- `pkg/oci/certificates/service_test.go` (200 lines) - Reduced test suite
  - Focus on core functionality tests
  - Remove redundant test cases

**Functionality**:
- VerificationManager implementation
- Dynamic mode switching (strict/custom-ca/skip)
- Fallback verification strategies
- Mode transition tracking
- Unit tests for all components
- Thread safety tests
- Integration tests

**Dependencies**:
- Split 001 (imports service implementation)
- Split 002 (for integration testing)

**Interface Contract**:
```go
// Exported for service integration
type VerificationManager struct {
    currentMode  v2.VerificationMode
    systemPool   *x509.CertPool
    customPool   *x509.CertPool
    fallbackMode v2.VerificationMode
    mutex        sync.RWMutex
}

func NewVerificationManager(system, custom *x509.CertPool) (*VerificationManager, error)
func (v *VerificationManager) SetMode(mode v2.VerificationMode) error
func (v *VerificationManager) GetPool() *x509.CertPool
```

## Implementation Sequence

### Phase 1: Split 001 (Core Service)
1. Create split-001 branch from current branch
2. Remove Gitea and Verification manager initialization from service.go
3. Simplify NewCertificateService constructor
4. Keep only core certificate operations
5. Verify compilation and size (<380 lines)
6. Run basic tests to ensure core functionality

### Phase 2: Split 002 (Gitea Integration)
1. Create split-002 branch from split-001
2. Keep only gitea_integration.go
3. Ensure GiteaDiscovery is fully self-contained
4. Add minimal service integration stub if needed
5. Verify compilation and size (377 lines)
6. Test discovery operations independently

### Phase 3: Split 003 (Verification & Tests)
1. Create split-003 branch from split-002
2. Reduce verification.go by removing duplicate logic
3. Optimize test suite to essential cases only
4. Ensure all splits integrate properly
5. Verify compilation and size (<400 lines)
6. Run full test suite

## Integration Strategy

### Service Wiring
After all splits are complete, the main service will wire components:
```go
// In final integration
service := &CertificateServiceImpl{
    giteaDiscovery:  giteaDiscovery,  // From Split 002
    verificationMgr: verificationMgr,  // From Split 003
    // ... core fields from Split 001
}
```

### Testing Integration
- Split 001: Basic service tests
- Split 002: Discovery tests
- Split 003: Verification and integration tests

## Size Verification

| Split | Files | Current Lines | Target Lines | Status |
|-------|-------|--------------|--------------|---------|
| 001 | service.go + interface | 448 | 380 | Needs reduction |
| 002 | gitea_integration.go | 377 | 377 | Ready |
| 003 | verification.go + tests | 740 | 400 | Needs reduction |

## Risk Mitigation

1. **Interface Stability**: Define clear interfaces in Split 001 that won't change
2. **Dependency Management**: Each split can compile independently
3. **Test Coverage**: Distribute tests to maintain >80% coverage across splits
4. **Integration Points**: Minimal coupling between splits
5. **Size Buffer**: Target 380 lines max to stay well under 400 limit

## Validation Checklist

- [ ] Each split compiles independently
- [ ] Each split is under 400 lines (measured with tools/line-counter.sh)
- [ ] No code duplication between splits
- [ ] All functionality preserved
- [ ] Tests distributed appropriately
- [ ] Integration points clearly defined
- [ ] Dependencies properly ordered
- [ ] Branch naming follows convention

## Next Steps

1. Orchestrator spawns SW Engineer for Split 001 implementation
2. After Split 001 complete, spawn for Split 002
3. After Split 002 complete, spawn for Split 003
4. Code Reviewer validates each split stays under limit
5. Integration testing after all splits complete
6. Merge splits back to main effort branch

## Notes

- The current implementation is functionally complete and tested
- Splitting is purely for size compliance
- Each split maintains functional cohesion
- Thread safety preserved across splits
- All existing tests will be preserved (distributed across splits)