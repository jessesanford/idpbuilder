# E3.1.3 Certificate Service Implementation Plan

## CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `idpbuidler-oci-mgmt/phase3/wave1/E3.1.3-certificate-validator`
**Can Parallelize**: Yes (after E3.1.1)
**Parallel With**: [E3.1.2, E3.1.4, E3.1.5]
**Size Estimate**: 750 lines
**Dependencies**: [E3.1.1]

## Overview
- **Effort**: Certificate Service Implementation
- **Phase**: 3, Wave: 1
- **Estimated Size**: 750 lines
- **Implementation Time**: 8-10 hours

### Purpose
Implement the complete CertificateService interface defined in E3.1.1, providing production-ready certificate management with Gitea-specific integration. This service will handle dynamic verification mode switching, thread-safe certificate pool management, and automatic Gitea certificate discovery.

## Dependencies from E3.1.1

### Interfaces to Implement
- `CertificateService` interface from `pkg/oci/api/v2/certificate_service.go`
- Certificate data structures and verification modes
- Validation methods for certificate configurations

### Key Methods Required
- `LoadCertificateBundle()` - Load certificates from various sources
- `SetVerificationMode()` - Dynamic mode switching (strict/custom-ca/skip)
- `ValidateCertificate()` - Certificate validation logic
- `LoadGiteaCertificate()` - Gitea-specific certificate handling
- `GetCertPool()` - Thread-safe pool access
- `AddCertificate()` - Add certificates to the pool
- `RemoveCertificate()` - Remove certificates from the pool

## File Structure

```
pkg/oci/certificates/
├── service.go (350 lines)
│   └── CertificateServiceImpl struct
│   └── NewCertificateService()
│   └── LoadCertificateBundle()
│   └── SetVerificationMode()
│   └── ValidateCertificate()
│   └── LoadGiteaCertificate()
│   └── GetCertPool()
│   └── AddCertificate()
│   └── RemoveCertificate()
│
├── gitea_integration.go (200 lines)
│   └── GiteaDiscovery struct
│   └── DiscoverGiteaCertificates()
│   └── LoadGiteaRootCA()
│   └── LoadGiteaRegistryCert()
│   └── ValidateGiteaCertChain()
│   └── ExtractCertFromGiteaConfig()
│
├── verification.go (150 lines)
│   └── VerificationManager struct
│   └── InitVerificationMode()
│   └── SwitchVerificationMode()
│   └── CreateCustomCAPool()
│   └── ValidateWithMode()
│   └── HandleVerificationFallback()
│
└── service_test.go (50 lines)
    └── TestServiceCreation
    └── TestVerificationModes
    └── TestThreadSafety
```

## Implementation Steps

### Step 1: Core Service Implementation (350 lines)
**File**: `pkg/oci/certificates/service.go`

1. **Define CertificateServiceImpl Structure**
   ```go
   type CertificateServiceImpl struct {
       mu                sync.RWMutex
       certPool          *x509.CertPool
       systemPool        *x509.CertPool
       verificationMode  VerificationMode
       giteaDiscovery    *GiteaDiscovery
       verificationMgr   *VerificationManager
       certificates      map[string]*x509.Certificate
       bundlePaths       []string
   }
   ```

2. **Implement Constructor**
   - Initialize with system certificate pool
   - Set up default verification mode
   - Create gitea discovery instance
   - Initialize thread-safe maps

3. **Implement LoadCertificateBundle()**
   - Support multiple bundle formats (PEM, DER)
   - Parse certificate chains
   - Add to appropriate pool
   - Handle errors gracefully
   - Log certificate subjects for debugging

4. **Implement SetVerificationMode()**
   - Thread-safe mode switching
   - Validate mode transitions
   - Update internal state
   - Rebuild certificate pools if needed
   - Emit mode change events

5. **Implement ValidateCertificate()**
   - Check certificate validity period
   - Verify signature chain
   - Check key usage constraints
   - Validate against current pool
   - Return detailed validation errors

6. **Implement Pool Management**
   - GetCertPool() with read lock
   - AddCertificate() with write lock
   - RemoveCertificate() with cleanup
   - Thread-safe pool operations

### Step 2: Gitea Integration (200 lines)
**File**: `pkg/oci/certificates/gitea_integration.go`

1. **Define GiteaDiscovery Structure**
   ```go
   type GiteaDiscovery struct {
       configPaths    []string
       registryURL    string
       discoveryCache map[string]*x509.Certificate
       mu            sync.RWMutex
   }
   ```

2. **Implement DiscoverGiteaCertificates()**
   - Search common Gitea configuration locations
   - Parse Gitea configuration files
   - Extract certificate references
   - Load discovered certificates
   - Cache for performance

3. **Implement LoadGiteaRootCA()**
   - Load Gitea's root CA certificate
   - Validate CA constraints
   - Add to certificate pool
   - Handle missing CA gracefully

4. **Implement LoadGiteaRegistryCert()**
   - Discover registry endpoint certificate
   - Support both embedded and external registry
   - Validate registry certificate
   - Handle self-signed certificates

5. **Implement ValidateGiteaCertChain()**
   - Verify complete certificate chain
   - Check intermediate certificates
   - Validate against Gitea root CA
   - Report chain issues clearly

### Step 3: Verification Management (150 lines)
**File**: `pkg/oci/certificates/verification.go`

1. **Define VerificationManager Structure**
   ```go
   type VerificationManager struct {
       currentMode    VerificationMode
       customCAPool   *x509.CertPool
       strictPool     *x509.CertPool
       fallbackMode   VerificationMode
       mu            sync.RWMutex
   }
   ```

2. **Implement InitVerificationMode()**
   - Set initial verification mode
   - Create appropriate certificate pools
   - Configure fallback behavior
   - Set up mode-specific validators

3. **Implement SwitchVerificationMode()**
   - Atomic mode switching
   - Preserve existing certificates
   - Update validation strategy
   - Log mode transitions

4. **Implement CreateCustomCAPool()**
   - Build custom CA pool
   - Merge system and custom CAs
   - Handle pool conflicts
   - Optimize pool size

5. **Implement ValidateWithMode()**
   - Mode-specific validation logic
   - Strict mode: system CAs only
   - Custom CA mode: system + custom
   - Skip verify: minimal validation
   - Return mode-appropriate errors

6. **Implement HandleVerificationFallback()**
   - Detect verification failures
   - Apply fallback strategy
   - Log fallback events
   - Maintain security posture

### Step 4: Testing (50 lines)
**File**: `pkg/oci/certificates/service_test.go`

1. **TestServiceCreation**
   - Verify proper initialization
   - Check default settings
   - Validate pool creation

2. **TestVerificationModes**
   - Test all three modes
   - Verify mode switching
   - Check fallback behavior

3. **TestThreadSafety**
   - Concurrent certificate operations
   - Parallel mode switching
   - Race condition detection

## Size Management
- **Estimated Lines**: 750 (350 + 200 + 150 + 50)
- **Measurement Tool**: `/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh`
- **Check Frequency**: After each file completion
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

### Size Breakdown by Component
- Service core logic: 350 lines
- Gitea integration: 200 lines
- Verification logic: 150 lines
- Tests: 50 lines

## Test Requirements

### Unit Tests (90% coverage)
- Certificate loading from various sources
- Verification mode switching
- Pool management operations
- Gitea certificate discovery
- Validation logic for all modes

### Integration Tests
- Mock Gitea server interaction
- Certificate rotation scenarios
- Mode switching under load
- Registry authentication with certificates
- Error recovery paths

### Thread Safety Tests
- Concurrent certificate additions
- Parallel validation requests
- Mode switching during operations
- Pool access contention
- Memory leak detection

### Performance Tests
- Certificate loading speed
- Validation throughput
- Pool lookup performance
- Mode switching latency

## Pattern Compliance

### idpbuilder-oci-mgmt Patterns
- Follow existing service interface patterns
- Use consistent error handling
- Implement proper logging with context
- Follow naming conventions from Phase 2

### Security Requirements
- Never log private keys
- Sanitize certificate subjects in logs
- Validate all certificate inputs
- Implement secure defaults
- Support security compliance modes

### Concurrency Patterns
- Use RWMutex for read-heavy operations
- Minimize lock contention
- Avoid nested locks
- Implement timeouts for lock acquisition

## Integration Points

### With E3.1.1 Contracts
- Implement all defined interfaces
- Use standard data structures
- Follow validation patterns
- Maintain API compatibility

### With E3.1.2 Certificate Loader
- Coordinate certificate formats
- Share parsing utilities
- Consistent error codes
- Unified logging approach

### With Phase 2 Components
- Registry client integration
- Build service certificate injection
- Authentication flow updates
- Error propagation patterns

## Error Handling Strategy

### Certificate Errors
- `ErrInvalidCertificate`: Malformed certificate data
- `ErrExpiredCertificate`: Certificate past validity
- `ErrUntrustedCertificate`: Cannot verify chain
- `ErrCertificateNotFound`: Missing required certificate
- `ErrVerificationFailed`: Validation failed

### Recovery Strategies
- Graceful degradation to skip-verify (if configured)
- Certificate reload on failure
- Automatic Gitea discovery retry
- Clear error messages for operators

## Success Criteria

### Functional Requirements
- ✅ All CertificateService methods implemented
- ✅ Three verification modes working
- ✅ Gitea certificates auto-discovered
- ✅ Thread-safe operations verified
- ✅ Graceful fallback functioning

### Performance Requirements
- ✅ Certificate validation < 10ms
- ✅ Mode switching < 1ms
- ✅ Pool operations < 1ms
- ✅ No memory leaks detected
- ✅ Concurrent operation support

### Quality Requirements
- ✅ 90% test coverage achieved
- ✅ No race conditions
- ✅ Clear error messages
- ✅ Comprehensive logging
- ✅ Production-ready code

## Implementation Notes

### Critical Considerations
1. **Thread Safety**: All public methods must be thread-safe
2. **Mode Transitions**: Must not drop certificates during switches
3. **Gitea Discovery**: Should not fail if Gitea is unavailable
4. **Fallback Logic**: Must maintain security while enabling operation
5. **Pool Management**: Must handle large certificate sets efficiently

### Common Pitfalls to Avoid
1. Don't parse certificates in critical paths
2. Avoid blocking operations in locks
3. Don't leak goroutines in discovery
4. Prevent certificate pool corruption
5. Handle nil pools gracefully

### Performance Optimizations
1. Cache parsed certificates
2. Use read locks for lookups
3. Batch certificate operations
4. Lazy-load Gitea certificates
5. Pre-compile validation rules

## Deliverables Checklist

- [ ] `service.go` implementing CertificateService (350 lines)
- [ ] `gitea_integration.go` with discovery logic (200 lines)
- [ ] `verification.go` with mode management (150 lines)
- [ ] `service_test.go` with comprehensive tests (50 lines)
- [ ] All tests passing with >90% coverage
- [ ] Thread safety verified with race detector
- [ ] Integration with E3.1.1 contracts validated
- [ ] Size verified under 800 lines using line counter
- [ ] Documentation comments on all public methods
- [ ] Error handling comprehensive and clear

---

**Document Version**: 1.0  
**Created**: 2025-08-27  
**Author**: Code Reviewer Agent  
**Based On**: Phase 3 Wave 1 Implementation Plan  
**Effort**: E3.1.3-certificate-validator