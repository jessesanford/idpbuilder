# SPLIT-PLAN-002.md

## Split 002 of 2: Transport Configuration and Utilities

**Planner**: Code Reviewer Agent
**Parent Effort**: registry-tls-trust-integration
**Effort Phase/Wave**: Phase 1, Wave 1

<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (⚠️⚠️⚠️ CRITICAL: All splits MUST reference SAME effort!)

- **Previous Split**: Split 001 of phase1/wave1/registry-tls-trust-integration
  - Path: efforts/phase1/wave1/registry-tls-trust-integration/split-001/
  - Branch: idpbuidler-oci-go-cr/phase1/wave1/registry-tls-trust-integration-split-001
  - Summary: Implemented core TrustStoreManager interface and trust store functionality
  
- **This Split**: Split 002 of phase1/wave1/registry-tls-trust-integration
  - Path: efforts/phase1/wave1/registry-tls-trust-integration/split-002/
  - Branch: idpbuidler-oci-go-cr/phase1/wave1/registry-tls-trust-integration-split-002
  
- **Next Split**: None (final split of THIS effort)

⚠️ NEVER reference splits from different efforts!
✅ RIGHT: All references are to phase1/wave1/registry-tls-trust-integration splits

## Files in This Split (EXCLUSIVE - no overlap with Split 001)

- `pkg/certs/transport.go` (270 lines) - Complete file
- `pkg/certs/trust_store.go` (217 lines) - Complete file  
- `pkg/certs/trust_test.go` (64 lines) - Transport and utility tests only

**Total Estimated Lines**: ~551 lines (well under 800 limit)

## Functionality

### Core Components:

1. **Transport Configuration** (`transport.go`):
   - TransportConfig struct with timeout and connection pool settings
   - DefaultTransportConfig() for standard settings
   - ConfigureTransport() - Creates remote.Option with TLS config
   - ConfigureTransportWithConfig() - Custom transport configuration
   - Integration with go-containerregistry remote package
   - TLS configuration based on trust store

2. **Trust Store Utilities** (`trust_store.go`):
   - TrustStoreUtils struct for utility operations
   - LoadCertificateFromPEM() - Load single certificate
   - LoadCertificatesFromPEM() - Load multiple certificates
   - SaveCertificateToPEM() - Save certificate to PEM format
   - ValidateCertificate() - Certificate validation
   - GetCertificateInfo() - Extract certificate metadata
   - Certificate expiry checking utilities

3. **Extended Tests** (`trust_test.go` additions):
   - Transport configuration tests
   - Utility function tests
   - Integration tests with trust store

## Dependencies

### From Split 001:
- TrustStoreManager interface (imported)
- Core trust store functionality

### External Dependencies:
- github.com/google/go-containerregistry/pkg/v1/remote
- Standard library: crypto/tls, net/http, time

## Implementation Instructions

### Step 1: Environment Setup
```bash
# Create split directory structure
cd efforts/phase1/wave1/registry-tls-trust-integration  
mkdir -p split-002/pkg/certs

# Copy Split 001's interface (for import)
# The SW Engineer should ensure Split 001 is accessible
```

### Step 2: Implement Transport Configuration
1. Create `pkg/certs/transport.go` with:
   - TransportConfig struct definition
   - Default configuration function
   - Transport option creators for go-containerregistry
   - TLS configuration based on trust store state
   - Connection pooling configuration

### Step 3: Implement Trust Store Utilities
1. Create `pkg/certs/trust_store.go` with:
   - TrustStoreUtils struct
   - Certificate loading utilities
   - PEM encoding/decoding functions
   - Certificate validation helpers
   - Certificate information extractors

### Step 4: Create Transport and Utility Tests
1. Create/extend `pkg/certs/trust_test.go` with tests for:
   - Transport configuration with various settings
   - TLS configuration verification
   - Utility function correctness
   - PEM loading/saving operations
   - Certificate validation logic

### Step 5: Verification
1. Ensure compilation: `go build ./pkg/certs/`
2. Run tests: `go test ./pkg/certs/`
3. Measure with line counter tool
4. Verify under 800 lines
5. Verify integration with Split 001 interfaces

## Integration Points

### Using Split 001:
```go
// Import the trust store manager from Split 001
import "github.com/cnoe-io/idpbuilder/pkg/certs"

// Use the TrustStoreManager interface
func (m *trustStoreManager) ConfigureTransport(registry string) (remote.Option, error) {
    // Access trust store methods from Split 001
    pool, err := m.GetCertPool(registry)
    // Configure transport with the cert pool
}
```

## Success Criteria

- [x] Transport configuration fully implemented
- [x] All utility functions implemented
- [x] Integration with Split 001 trust store working
- [x] Tests cover transport and utilities
- [x] Implementation compiles without errors
- [x] Total lines < 800 (target ~551)
- [x] Proper use of Split 001 interfaces

## Notes for SW Engineer

- Import and use the TrustStoreManager interface from Split 001
- Focus on transport layer and utilities only
- Ensure proper integration with go-containerregistry
- Test transport configuration thoroughly
- Maintain clean separation from Split 001 code
- Do not duplicate any Split 001 functionality

## Merge Strategy

After both splits are complete and reviewed:
1. Merge Split 001 to parent branch first
2. Merge Split 002 to parent branch second
3. Verify combined functionality works correctly
4. Total combined size should be ~928 lines (split across two PRs)