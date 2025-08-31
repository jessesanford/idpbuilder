# SPLIT-PLAN-001.md

## Split 001 of 2: Core Trust Store Management

**Planner**: Code Reviewer Agent
**Parent Effort**: registry-tls-trust-integration
**Effort Phase/Wave**: Phase 1, Wave 1

<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (⚠️⚠️⚠️ CRITICAL: All splits MUST reference SAME effort!)

- **Previous Split**: None (first split of THIS effort)
  - Path: N/A (this is Split 001)
  - Branch: N/A
  
- **This Split**: Split 001 of phase1/wave1/registry-tls-trust-integration
  - Path: efforts/phase1/wave1/registry-tls-trust-integration/split-001/
  - Branch: idpbuidler-oci-go-cr/phase1/wave1/registry-tls-trust-integration-split-001
  
- **Next Split**: Split 002 of phase1/wave1/registry-tls-trust-integration
  - Path: efforts/phase1/wave1/registry-tls-trust-integration/split-002/
  - Branch: idpbuidler-oci-go-cr/phase1/wave1/registry-tls-trust-integration-split-002

## Files in This Split (EXCLUSIVE - no overlap with other splits)

- `pkg/certs/trust.go` (317 lines) - Complete file
- `pkg/certs/trust_test.go` (60 lines) - Core trust store tests only

**Total Estimated Lines**: ~377 lines (well under 800 limit)

## Functionality

### Core Components:
1. **TrustStoreManager Interface**:
   - AddCertificate - Add certificates for registries
   - RemoveCertificate - Remove registry certificates
   - SetInsecureRegistry - Mark registries as insecure
   - GetTrustedCerts - Retrieve trusted certificates
   - GetCertPool - Get configured cert pool
   - IsInsecure - Check insecure status
   - LoadFromDisk - Load persistent certificates
   - SaveToDisk - Persist certificates

2. **trustStoreManager Implementation**:
   - Thread-safe certificate storage with sync.RWMutex
   - In-memory certificate management
   - Insecure registry tracking
   - Disk persistence for certificates
   - Certificate validation and expiry checking
   
3. **Core Functionality**:
   - NewTrustStoreManager constructor
   - Certificate lifecycle management
   - Registry-specific trust configuration
   - System certificate pool integration

## Dependencies

- **External**: None (foundational split)
- **Standard Library**:
  - crypto/x509 - Certificate handling
  - encoding/pem - PEM encoding
  - sync - Thread safety
  - os, io/ioutil - File operations
  - path/filepath - Path management

## Implementation Instructions

### Step 1: Environment Setup
```bash
# Create split directory structure
cd efforts/phase1/wave1/registry-tls-trust-integration
mkdir -p split-001/pkg/certs
```

### Step 2: Implement Core Trust Store
1. Create `pkg/certs/trust.go` with:
   - TrustStoreManager interface definition
   - trustStoreManager struct
   - All methods from the interface
   - Constructor function
   - Certificate validation logic

### Step 3: Create Basic Tests
1. Create `pkg/certs/trust_test.go` with tests for:
   - AddCertificate functionality
   - RemoveCertificate functionality
   - Insecure registry management
   - Certificate pool creation
   - Thread safety verification

### Step 4: Verification
1. Ensure compilation: `go build ./pkg/certs/`
2. Run tests: `go test ./pkg/certs/`
3. Measure with line counter tool
4. Verify under 800 lines

## Interface Contracts (for Split 002)

Split 002 will depend on the following exported types from this split:
- `TrustStoreManager` interface
- `NewTrustStoreManager()` constructor

These must remain stable to ensure Split 002 can build upon them.

## Success Criteria

- [x] Trust store manager fully implemented
- [x] Interface properly defined and exported
- [x] Basic test coverage provided
- [x] Implementation compiles without errors
- [x] Total lines < 800 (target ~377)
- [x] No dependencies on Split 002 components

## Notes for SW Engineer

- Focus on core trust store functionality only
- Do not implement transport or utility functions (those go in Split 002)
- Ensure thread safety with proper mutex usage
- Keep interface minimal but complete for transport layer usage
- Test file should only test core trust store operations