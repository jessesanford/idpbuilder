# Effort 1.2.2: Registry Client Implementation - Implementation Plan

**Created**: 2025-10-29 06:31:00 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**Effort ID**: 1.2.2
**Phase**: Phase 1 - Foundation & Interfaces
**Wave**: Wave 2 - Core Package Implementations

---

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)

### R213 Infrastructure Metadata

```json
{
  "effort_id": "1.2.2",
  "effort_name": "Registry Client Implementation",
  "branch_name": "idpbuilder-oci-push/phase1/wave2/effort-2-registry-client",
  "base_branch": "idpbuilder-oci-push/phase1/wave2/integration",
  "parent_wave": "WAVE_2",
  "parent_phase": "PHASE_1",
  "depends_on": [],
  "estimated_lines": 450,
  "complexity": "high",
  "can_parallelize": true,
  "parallel_with": ["1.2.1", "1.2.3", "1.2.4"]
}
```

**Branch**: `idpbuilder-oci-push/phase1/wave2/effort-2-registry-client`
**Can Parallelize**: Yes
**Parallel With**: [1.2.1, 1.2.3, 1.2.4] (ALL Wave 2 efforts run simultaneously)
**Size Estimate**: 450 lines (MUST be <800)
**Dependencies**: None (all Wave 1 interfaces frozen and available)

---

## Overview

**Purpose**: Implement the registry client package that pushes OCI images to container registries, validates registry connectivity, builds fully-qualified image references, and classifies errors correctly.

**What This Effort Accomplishes**:
- Complete implementation of `registry.Client` interface (frozen in Wave 1)
- OCI image push operations using go-containerregistry
- Registry connectivity validation (/v2/ endpoint ping)
- Image reference construction (registry/namespace/repository:tag)
- Progress callback support for layer uploads
- Error classification (auth errors, network errors, push failures)
- Integration with auth and TLS providers

**Boundaries - OUT OF SCOPE**:
- Authentication implementation (auth package responsibility)
- TLS configuration implementation (tls package responsibility)
- Image building or modification (docker package responsibility)
- Registry management or administration
- Image manifest manipulation (go-containerregistry handles this)

---

## File Structure

### New Files to Create

**Implementation Files**:
- `pkg/registry/client.go` (~450 lines)
  - `registryClient` struct implementation
  - `NewClient()` constructor with provider validation
  - `Push()` method with progress callbacks
  - `BuildImageReference()` reference construction
  - `ValidateRegistry()` /v2/ endpoint check
  - Helper functions (`parseImageName`, `createProgressHandler`, `isAuthError`, `isNetworkError`)

**Test Files** (NOT counted in line estimates per R007):
- `pkg/registry/client_test.go` (~400 lines)
  - 15+ test cases covering all methods
  - Mock auth and TLS providers
  - Success paths, error paths, edge cases
  - Progress callback validation

**Modified Files**:
- None (go.mod already has go-containerregistry from Wave 1)

**Total Estimated Lines**: 450 lines (implementation only, tests excluded per R007)

---

## Implementation Steps

### Step 1: Review Wave 1 Interface Definitions

**MANDATORY FIRST STEP - Read frozen interfaces**:
```bash
# Read the frozen Registry interface from Wave 1
cat pkg/registry/types.go

# Read the Auth Provider interface (dependency)
cat pkg/auth/types.go

# Read the TLS ConfigProvider interface (dependency)
cat pkg/tls/types.go
```

**Expected Registry Interface** (from Wave 1 Effort 2):
```go
package registry

import (
    "context"
    v1 "github.com/google/go-containerregistry/pkg/v1"
    "your-module/pkg/auth"
    "your-module/pkg/tls"
)

type Client interface {
    Push(ctx context.Context, image v1.Image, targetRef string, callback ProgressCallback) error
    BuildImageReference(registryURL string, imageName string) (string, error)
    ValidateRegistry(ctx context.Context, registryURL string) error
}

type ProgressCallback func(update ProgressUpdate)

type ProgressUpdate struct {
    Layer    string
    Bytes    int64
    Total    int64
    Complete bool
}
```

**Expected Error Types** (from Wave 1):
```go
type AuthenticationError struct { Endpoint string, StatusCode int, Message string }
type NetworkError struct { Endpoint string, Cause error }
type RegistryUnavailableError struct { Endpoint string, StatusCode int }
type PushFailedError struct { Reference string, Cause error }
type ValidationError struct { Field string, Message string }
```

**Auth Provider Interface** (from Wave 1 Effort 3):
```go
package auth

type Provider interface {
    GetAuthenticator() (authn.Authenticator, error)
    ValidateCredentials() error
}
```

**TLS ConfigProvider Interface** (from Wave 1 Effort 3):
```go
package tls

import "crypto/tls"

type ConfigProvider interface {
    GetTLSConfig() *tls.Config
    IsInsecure() bool
}
```

### Step 2: Implement pkg/registry/client.go

**File: pkg/registry/client.go**

**Required Implementation Details**:

**1. Struct Definition**:
```go
type registryClient struct {
    authProvider auth.Provider
    tlsConfig    tls.ConfigProvider
    httpClient   *http.Client
}
```

**2. NewClient() Implementation** (~60 lines):
- Validate `authProvider` is not nil → `&ValidationError{Field: "authProvider", Message: "cannot be nil"}`
- Validate `tlsConfig` is not nil → `&ValidationError{Field: "tlsConfig", Message: "cannot be nil"}`
- Create HTTP client with TLS config:
  ```go
  tlsCfg := tlsConfig.GetTLSConfig()
  transport := &http.Transport{
      TLSClientConfig: tlsCfg,
  }
  httpClient := &http.Client{Transport: transport}
  ```
- Store providers in `registryClient` struct
- Return implementation of `Client` interface

**3. Push() Implementation** (~120 lines):
- Parse target reference using `name.ParseReference(targetRef)`
- Return `&ValidationError{...}` if parse fails
- Get authenticator from auth provider: `authenticator, err := authProvider.GetAuthenticator()`
- Return error if auth provider fails
- Configure remote options:
  ```go
  options := []remote.Option{
      remote.WithAuth(authenticator),
      remote.WithTransport(httpClient.Transport),
      remote.WithContext(ctx),
  }
  if callback != nil {
      options = append(options, remote.WithProgress(createProgressHandler(callback)))
  }
  ```
- Call `remote.Write(ref, image, options...)`
- Classify errors using helper functions:
  - If `isAuthError(err)` → `&AuthenticationError{...}`
  - If `isNetworkError(err)` → `&NetworkError{...}`
  - Otherwise → `&PushFailedError{Reference: targetRef, Cause: err}`

**4. BuildImageReference() Implementation** (~70 lines):
- Parse registry URL using `url.Parse(registryURL)`
- Return `&ValidationError{...}` if parse fails
- Extract host:port from parsed URL (`url.Host`)
- Parse image name:
  - Split on `:` to separate repository and tag
  - Default tag to "latest" if not specified
  - Return `&ValidationError{...}` if image name empty
- Inject "giteaadmin" as namespace (Gitea default)
- Construct reference: `{host:port}/giteaadmin/{repository}:{tag}`
- Example: `https://gitea.cnoe.localtest.me:8443` + `myapp:v1` → `gitea.cnoe.localtest.me:8443/giteaadmin/myapp:v1`

**5. ValidateRegistry() Implementation** (~60 lines):
- Parse registry URL using `url.Parse(registryURL)`
- Return `&ValidationError{...}` if parse fails
- Build /v2/ endpoint: `{scheme}://{host}/v2/`
- Create HTTP GET request with context: `http.NewRequestWithContext(ctx, "GET", endpoint, nil)`
- Execute request using HTTP client: `resp, err := httpClient.Do(req)`
- Handle responses:
  - 200 OK → Success (registry reachable)
  - 401 Unauthorized → Success (registry reachable, just needs auth)
  - Connection error → `&NetworkError{Endpoint: endpoint, Cause: err}`
  - Other status codes → `&RegistryUnavailableError{Endpoint: endpoint, StatusCode: resp.StatusCode}`

**6. Helper Functions** (~140 lines):

**parseImageName(imageName string) (repository, tag string, err error)**:
- Split image name on `:` character
- Default tag to "latest" if not specified
- Validate repository name (not empty, valid characters)
- Return repository and tag

**createProgressHandler(callback ProgressCallback) chan v1.Update**:
- Create channel for progress updates
- Start goroutine to convert v1.Update to ProgressUpdate
- Call user's callback with converted updates
- Handle channel close properly

**isAuthError(err error) bool**:
- Check error string for: "401", "403", "unauthorized", "forbidden"
- Return true if authentication-related error

**isNetworkError(err error) bool**:
- Check error string for: "connection", "timeout", "network", "dial", "i/o timeout"
- Return true if network-related error

**Complete Package Structure**:
```go
package registry

import (
    "context"
    "net/http"
    "net/url"
    "strings"

    "github.com/google/go-containerregistry/pkg/name"
    "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/remote"

    "your-module/pkg/auth"
    "your-module/pkg/tls"
)

// Implementation here (~450 lines total)
```

### Step 3: Write Tests (TDD - Tests First!)

**File: pkg/registry/client_test.go**

**Test Cases** (from Wave 2 Test Plan):

**A. Constructor Tests**:
- `TestNewClient_Success`: NewClient succeeds with valid providers
- `TestNewClient_NilAuthProvider`: NewClient fails with nil auth provider
- `TestNewClient_NilTLSProvider`: NewClient fails with nil TLS provider

**B. Push Tests**:
- `TestPush_Success`: Push succeeds with progress callbacks
- `TestPush_AuthenticationError`: Push returns AuthenticationError for 401/403
- `TestPush_NetworkError`: Push returns NetworkError for unreachable registry
- `TestPush_PushFailedError`: Push returns PushFailedError for invalid target reference

**C. BuildImageReference Tests**:
- `TestBuildImageReference_Success`: Constructs correct references
  - `https://gitea.cnoe.localtest.me:8443` + `myapp:latest` → `gitea.cnoe.localtest.me:8443/giteaadmin/myapp:latest`
  - `https://registry.io` + `myapp` → `registry.io/giteaadmin/myapp:latest` (default tag)
- `TestBuildImageReference_ValidationError`: Returns ValidationError for invalid URLs

**D. ValidateRegistry Tests**:
- `TestValidateRegistry_Success`: Succeeds for reachable registry (200 or 401)
- `TestValidateRegistry_NetworkError`: Returns NetworkError for unreachable registry

**Mock Providers**:
```go
// Mock auth provider for testing
type mockAuthProvider struct {
    authenticator authn.Authenticator
    err           error
}

func (m *mockAuthProvider) GetAuthenticator() (authn.Authenticator, error) {
    return m.authenticator, m.err
}

func (m *mockAuthProvider) ValidateCredentials() error {
    return nil
}

// Mock TLS provider for testing
type mockTLSProvider struct {
    config   *tls.Config
    insecure bool
}

func (m *mockTLSProvider) GetTLSConfig() *tls.Config {
    return m.config
}

func (m *mockTLSProvider) IsInsecure() bool {
    return m.insecure
}
```

**Test Coverage Requirements**:
- Minimum 85% code coverage
- All success paths tested
- All failure paths tested
- Error classification tested (auth vs network vs push failures)
- Progress callback functionality tested

### Step 4: Size Measurement

**Measure implementation lines**:
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning
tools/line-counter.sh

# Expected output:
# 🎯 Detected base: idpbuilder-oci-push/phase1/wave2/integration
# 📦 Analyzing branch: idpbuilder-oci-push/phase1/wave2/effort-2-registry-client
# ✅ Total implementation lines: ~450
```

**Size Compliance**:
- Target: 450 lines
- Buffer: ±15% (383-518 lines acceptable)
- Hard limit: 800 lines (MUST NOT EXCEED)
- Tests NOT counted (per R007)

**If approaching 700 lines**:
- STOP IMMEDIATELY
- Report to orchestrator
- Do NOT exceed 800 lines

### Step 5: Run Tests and Coverage

**Run unit tests**:
```bash
cd pkg/registry
go test -v -cover

# Expected: All tests pass, coverage ≥85%
```

**Generate coverage report**:
```bash
go test -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
go tool cover -func=coverage.out | grep total
# Expected: total: (statements) 85.0% or higher
```

### Step 6: Linting and Documentation

**Run linters**:
```bash
go vet ./pkg/registry/...
golangci-lint run ./pkg/registry/...

# Expected: No errors
```

**Verify GoDoc**:
- All public types have GoDoc comments
- All public functions have GoDoc comments
- Examples provided for main methods

### Step 7: Commit and Push

**Commit structure**:
```bash
git add pkg/registry/client.go pkg/registry/client_test.go
git commit -m "feat(registry): implement OCI registry client with push support

- Implement NewClient with provider validation
- Implement Push with progress callbacks and error classification
- Implement BuildImageReference with giteaadmin namespace
- Implement ValidateRegistry with /v2/ endpoint check
- Add helper functions for error classification
- Add 15 test cases with 85%+ coverage
- Support auth and TLS provider integration

Closes: Effort 1.2.2 - Registry Client Implementation
Lines: ~450 (within 450 ±15% estimate)
Coverage: 85%+ (meets Wave 2 Test Plan requirements)"

git push origin idpbuilder-oci-push/phase1/wave2/effort-2-registry-client
```

---

## Size Management

**Estimated Lines**: 450 lines (implementation code only)
**Measurement Tool**: `${PROJECT_ROOT}/tools/line-counter.sh` (find project root first)
**Check Frequency**: After every major function implementation (~60 lines)
**Split Threshold**:
- Warning: 700 lines (approaching limit)
- Hard stop: 800 lines (MUST NOT EXCEED)

**Size Tracking**:
- After struct definition: ~60 lines
- After NewClient: ~120 lines
- After Push: ~240 lines
- After BuildImageReference: ~310 lines
- After ValidateRegistry: ~370 lines
- After helpers and comments: ~450 lines (target)

---

## Test Requirements

**Minimum Coverage**: 85% (per Wave 2 Test Plan)

**Test Categories** (from Wave 2 Test Plan):

| Test Category | Test Cases | Coverage Target |
|---------------|------------|-----------------|
| Constructor | 3 tests | 100% of NewClient |
| Push | 4 tests | 100% of Push |
| BuildImageReference | 2 tests | 100% of BuildImageReference |
| ValidateRegistry | 2 tests | 100% of ValidateRegistry |
| Helper Functions | 4 tests | Error classification |

**Total Test Cases**: 15+ tests

**Test Execution**:
```bash
go test ./pkg/registry -v -cover
# MUST achieve ≥85% coverage
```

---

## Dependencies

### Upstream Dependencies (COMPLETED)
- ✅ Wave 1 Effort 2: Registry interface definition (FROZEN)
- ✅ Wave 1 Effort 3: Auth and TLS interface definitions (FROZEN)
- ✅ Integration branch: `idpbuilder-oci-push/phase1/wave2/integration` (CREATED)

### Downstream Dependencies
- None (all Wave 2 efforts are parallel)
- Effort 1.2.3 will PROVIDE auth implementation (parallel development)
- Effort 1.2.4 will PROVIDE TLS implementation (parallel development)
- Wave 3 CLI will use this package

### External Library Dependencies
**Existing Dependencies** (from Wave 1):
- `github.com/google/go-containerregistry` v0.19.0
  - `pkg/name` for reference parsing
  - `pkg/v1/remote` for push operations
- `github.com/stretchr/testify` v1.10.0 (testing)

**No New Dependencies Required**

### Internal Package Dependencies
- `pkg/auth.Provider` interface (Wave 1, will be implemented in Effort 1.2.3)
- `pkg/tls.ConfigProvider` interface (Wave 1, will be implemented in Effort 1.2.4)

**Note**: Use mock providers for testing since actual implementations are being built in parallel.

---

## Pattern Compliance

### Go Patterns
- Interface-driven design (implement `registry.Client` interface)
- Dependency injection (accept auth and TLS providers)
- Error wrapping with custom error types
- Context propagation for cancellation
- Progress callbacks for user feedback

### Security Requirements
- Error classification (don't expose internal details)
- TLS configuration from provider (secure by default)
- Authentication via provider interface
- Proper timeout handling

### Performance Targets
- Registry validation should complete in <2 seconds
- Push time depends on image size (no specific target)
- Progress callbacks should fire at reasonable intervals

---

## Acceptance Criteria

**MANDATORY - All must pass before Code Review**:

- [ ] All files created/modified as specified
- [ ] All 3 interface methods implemented correctly (Push, BuildImageReference, ValidateRegistry)
- [ ] All tests passing (100% pass rate)
- [ ] Code coverage ≥85% (per Wave 2 Test Plan)
- [ ] No linting errors (go vet, golangci-lint)
- [ ] Documentation complete (all public methods have GoDoc)
- [ ] Line count within estimate (450 lines ±15% = 383-518 lines)
- [ ] Integration with go-containerregistry working (remote.Write)
- [ ] Error classification correct (auth vs network vs push)
- [ ] Progress callbacks functional
- [ ] /v2/ endpoint validation working
- [ ] Image reference construction correct (with giteaadmin namespace)
- [ ] Mock providers working in tests
- [ ] Code committed and pushed to effort branch

**Quality Gates**:
1. **Functionality**: All interface methods work correctly
2. **Testing**: 85%+ coverage with all paths tested
3. **Error Handling**: Proper error classification
4. **Size**: Within 450 ±15% lines (383-518)
5. **Documentation**: Complete GoDoc coverage

---

## References

**Wave 2 Planning Documents**:
- Wave Implementation Plan: `/home/vscode/workspaces/idpbuilder-oci-push-planning/wave-plans/WAVE-2-IMPLEMENTATION.md`
- Wave Architecture: `/home/vscode/workspaces/idpbuilder-oci-push-planning/wave-plans/WAVE-2-ARCHITECTURE.md`
- Wave Test Plan: `/home/vscode/workspaces/idpbuilder-oci-push-planning/wave-plans/WAVE-2-TEST-PLAN.md`

**Wave 1 Interfaces** (frozen references):
- Registry Interface: `efforts/phase1/wave1/effort-2-registry-interface/pkg/registry/types.go`
- Registry Errors: `efforts/phase1/wave1/effort-2-registry-interface/pkg/registry/errors.go`
- Auth Interface: `efforts/phase1/wave1/effort-3-auth-tls-interfaces/pkg/auth/types.go`
- TLS Interface: `efforts/phase1/wave1/effort-3-auth-tls-interfaces/pkg/tls/types.go`

**External Documentation**:
- go-containerregistry: https://github.com/google/go-containerregistry
- OCI Distribution Spec: https://github.com/opencontainers/distribution-spec
- Docker Registry API: https://docs.docker.com/registry/spec/api/

---

## Implementation Checklist

**Pre-Implementation**:
- [ ] Read Wave 1 Registry interface definition
- [ ] Read Wave 1 Auth and TLS interfaces
- [ ] Read Wave 2 Architecture (Registry section)
- [ ] Read Wave 2 Test Plan (Registry test cases)
- [ ] Checkout effort branch from integration
- [ ] Verify base branch is correct

**Implementation Phase**:
- [ ] Write test stubs (15+ test cases)
- [ ] Create mock auth and TLS providers
- [ ] Implement `registryClient` struct
- [ ] Implement `NewClient()` with provider validation
- [ ] Implement `Push()` with progress callbacks
- [ ] Implement `BuildImageReference()` with namespace injection
- [ ] Implement `ValidateRegistry()` with /v2/ check
- [ ] Implement helper functions (error classification)
- [ ] Complete test implementations
- [ ] Run tests (all pass, 85%+ coverage)
- [ ] Run linters (no errors)
- [ ] Add GoDoc comments

**Validation Phase**:
- [ ] Measure size (within 383-518 lines)
- [ ] Verify coverage ≥85%
- [ ] Error classification validation
- [ ] Progress callback testing
- [ ] Reference construction verification
- [ ] Commit and push code

---

## Next Steps

**After Implementation Completion**:
1. Code Reviewer will be spawned for effort review
2. Code Reviewer validates all acceptance criteria
3. If approved: Merge to integration branch
4. If issues found: Fix and re-submit for review
5. After all 4 Wave 2 efforts approved: Wave integration testing

**Parallel Work**:
- This effort (1.2.2) runs in parallel with:
  - Effort 1.2.1: Docker Client Implementation
  - Effort 1.2.3: Authentication Implementation
  - Effort 1.2.4: TLS Configuration Implementation

**No coordination needed** - all efforts are independent until integration phase.

---

## Document Status

**Status**: ✅ READY FOR IMPLEMENTATION
**Created**: 2025-10-29 06:31:00 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**Effort**: 1.2.2 (Registry Client Implementation)
**Wave**: Wave 2 of Phase 1
**Branch**: `idpbuilder-oci-push/phase1/wave2/effort-2-registry-client`
**Base Branch**: `idpbuilder-oci-push/phase1/wave2/integration`

**Compliance**:
- ✅ R213: Complete metadata included
- ✅ R211: Parallelization specified (runs with 1.2.1, 1.2.3, 1.2.4)
- ✅ R341: TDD approach (test plan before implementation)
- ✅ R381: Library versions locked (go-containerregistry v0.19.0 from Wave 1)
- ✅ R383: Plan stored in .software-factory with timestamp
- ✅ Size compliance: 450 lines < 800 hard limit

---

**END OF EFFORT 1.2.2 IMPLEMENTATION PLAN**
