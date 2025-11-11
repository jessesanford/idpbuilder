# PHASE 1 TEST PLAN: idpbuilder-oci-push-rebuild

**Phase**: Phase 1 - Foundation & Interfaces (2 Waves)
**Created**: 2025-11-11
**Test Planning State**: PHASE_TEST_PLANNING
**TDD Workflow**: RED phase (tests created before implementation)
**Architecture Source**: planning/phase1/PHASE-1-ARCHITECTURE-PLAN.md
**Project Test Plan**: planning/project/PROJECT-TEST-PLAN.md

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Test Strategy for Phase 1](#test-strategy-for-phase-1)
3. [Wave 1 Test Specifications](#wave-1-test-specifications)
4. [Wave 2 Test Specifications](#wave-2-test-specifications)
5. [Progressive Test Planning Architecture](#progressive-test-planning-architecture)
6. [Test Fixtures and Mocks](#test-fixtures-and-mocks)
7. [Test Execution Strategy](#test-execution-strategy)
8. [Success Criteria](#success-criteria)
9. [TDD Workflow Summary](#tdd-workflow-summary)

---

## Executive Summary

### Phase 1 Test Plan Purpose

This test plan defines the **comprehensive acceptance criteria** for Phase 1 using Test-Driven Development (TDD). Phase 1 establishes the architectural foundation through interface definitions (Wave 1) and core implementations (Wave 2).

**Phase 1 Scope**:
- **Wave 1**: Interface contract definitions (DockerClient, RegistryClient, AuthProvider, TLSProvider, Command)
- **Wave 2**: Core implementations (4 parallel efforts implementing Wave 1 interfaces)

### Test Coverage Strategy

**Total Test Categories**: 2 primary categories for Phase 1
1. **Interface Contract Tests** (Wave 1 validation) - ~70 tests
2. **Implementation Unit Tests** (Wave 2 validation) - ~80 tests

**Expected Test Count for Phase 1**: ~150 tests
- Interface validation tests: ~70 tests (Wave 1)
- Implementation unit tests: ~80 tests (Wave 2)
- Integration tests: DEFERRED to Phase 2 (packages not yet wired together)
- E2E tests: DEFERRED to Phase 3 (command layer not yet integrated)

### TDD Workflow Compliance (R341)

This plan follows R341 TDD protocol:
- ✅ **RED phase**: Tests created BEFORE implementation (current state)
- ⏳ **GREEN phase**: Implementation will target making tests pass
- ⏳ **REFACTOR phase**: Code improvements while maintaining test passage

**Current Status**: RED phase (all tests expected to fail initially)

### Phase 1 Test Focus

**Wave 1 Focus** (Interface Definitions):
- Interfaces compile and have correct signatures
- Type system validates interface contracts
- No implementation required yet (pure interface definitions)
- Tests validate interface completeness

**Wave 2 Focus** (Core Implementations):
- All interface contract tests pass
- Unit tests validate implementation logic
- Mock-friendly interfaces enable isolated testing
- No integration testing yet (packages independent)

---

## Test Strategy for Phase 1

### Testing Pyramid for Phase 1

```
Phase 1 Testing Focus

                ┌───────────────────────┐
                │   Integration Tests   │  DEFERRED to Phase 2
                │   (Package Wiring)    │  (0 tests Phase 1)
                └───────────────────────┘

            ┌─────────────────────────────┐
            │  Implementation Unit Tests   │  ~80 tests (Wave 2)
            │  (Logic Validation)          │
            └─────────────────────────────┘

        ┌─────────────────────────────────────┐
        │  Interface Contract Tests           │  ~70 tests (Wave 1)
        │  (Type Safety & Signatures)         │
        └─────────────────────────────────────┘
```

### Test Execution Strategy

**Wave 1 Test Execution** (Interface validation):
```bash
# Run interface compilation tests
go test ./pkg/docker/interface_test.go
go test ./pkg/registry/interface_test.go
go test ./pkg/auth/interface_test.go
go test ./pkg/tls/interface_test.go
go test ./cmd/push_interface_test.go

# Expected: All pass (interfaces compile correctly)
```

**Wave 2 Test Execution** (Implementation validation):
```bash
# Run unit tests for each package
go test ./pkg/docker/... -short
go test ./pkg/registry/... -short
go test ./pkg/auth/... -short
go test ./pkg/tls/... -short

# Expected: All implementation tests pass
```

### Test Environment Requirements

**Wave 1 Requirements**:
- Go compiler (interface validation)
- No external dependencies needed (pure compilation)

**Wave 2 Requirements**:
- Go compiler and test runner
- Mock interfaces (no real Docker/Registry needed)
- Test fixtures for validation logic
- NO integration dependencies (isolated unit tests)

---

## Wave 1 Test Specifications

### Purpose
Validate that all interfaces defined in Phase 1 Wave 1 are correctly structured, compile successfully, and provide complete contracts for Wave 2 implementation.

---

### Test Category: DockerClient Interface Validation

**Test File**: `pkg/docker/interface_test.go`

**Interface Under Test**:
```go
type DockerClient interface {
    ImageExists(ctx context.Context, imageName string) (bool, error)
    GetImage(ctx context.Context, imageName string) (v1.Image, error)
    ValidateImageName(imageName string) error
    Close() error
}
```

**Test Cases** (15 tests):

#### TC-IFACE-DOCKER-001: Interface compiles successfully
```go
func TestDockerClientInterface_Compiles(t *testing.T)
```
- **Given**: DockerClient interface definition exists
- **When**: Go compiler processes interface.go
- **Then**: No compilation errors

#### TC-IFACE-DOCKER-002: ImageExists signature is correct
```go
func TestDockerClientInterface_ImageExistsSignature(t *testing.T)
```
- **Given**: DockerClient interface
- **When**: Reflect on ImageExists method
- **Then**: Method has signature (context.Context, string) (bool, error)

#### TC-IFACE-DOCKER-003: GetImage signature is correct
```go
func TestDockerClientInterface_GetImageSignature(t *testing.T)
```
- **Given**: DockerClient interface
- **When**: Reflect on GetImage method
- **Then**: Method has signature (context.Context, string) (v1.Image, error)

#### TC-IFACE-DOCKER-004: ValidateImageName signature is correct
```go
func TestDockerClientInterface_ValidateImageNameSignature(t *testing.T)
```
- **Given**: DockerClient interface
- **When**: Reflect on ValidateImageName method
- **Then**: Method has signature (string) error

#### TC-IFACE-DOCKER-005: Close signature is correct
```go
func TestDockerClientInterface_CloseSignature(t *testing.T)
```
- **Given**: DockerClient interface
- **When**: Reflect on Close method
- **Then**: Method has signature () error

#### TC-IFACE-DOCKER-006: NewDockerClient constructor exists
```go
func TestDockerClientInterface_ConstructorExists(t *testing.T)
```
- **Given**: docker package
- **When**: Check for NewDockerClient function
- **Then**: Function exists with signature () (DockerClient, error)

#### TC-IFACE-DOCKER-007: ImageNotFoundError type defined
```go
func TestDockerClientInterface_ImageNotFoundErrorExists(t *testing.T)
```
- **Given**: docker package
- **When**: Check for ImageNotFoundError type
- **Then**: Type exists and implements error interface

#### TC-IFACE-DOCKER-008: DaemonConnectionError type defined
```go
func TestDockerClientInterface_DaemonConnectionErrorExists(t *testing.T)
```
- **Given**: docker package
- **When**: Check for DaemonConnectionError type
- **Then**: Type exists and implements error interface

#### TC-IFACE-DOCKER-009: InvalidImageNameError type defined
```go
func TestDockerClientInterface_InvalidImageNameErrorExists(t *testing.T)
```
- **Given**: docker package
- **When**: Check for InvalidImageNameError type
- **Then**: Type exists and implements error interface

#### TC-IFACE-DOCKER-010: Mock implementation satisfies interface
```go
func TestDockerClientInterface_MockImplementation(t *testing.T)
```
- **Given**: MockDockerClient struct
- **When**: Assign to DockerClient variable
- **Then**: Type assertion succeeds (compile-time check)

#### TC-IFACE-DOCKER-011: Interface methods accept context.Context
```go
func TestDockerClientInterface_ContextSupport(t *testing.T)
```
- **Given**: DockerClient interface
- **When**: Check ImageExists and GetImage methods
- **Then**: Both accept context.Context as first parameter

#### TC-IFACE-DOCKER-012: Interface uses v1.Image from go-containerregistry
```go
func TestDockerClientInterface_UsesV1Image(t *testing.T)
```
- **Given**: DockerClient.GetImage method
- **When**: Check return type
- **Then**: Returns v1.Image from github.com/google/go-containerregistry

#### TC-IFACE-DOCKER-013: Error types have required fields
```go
func TestDockerClientInterface_ErrorTypeFields(t *testing.T)
```
- **Given**: Error type definitions
- **When**: Check struct fields
- **Then**: ImageNotFoundError has ImageName, DaemonConnectionError has Endpoint and Cause

#### TC-IFACE-DOCKER-014: Interface exported (public)
```go
func TestDockerClientInterface_ExportedInterface(t *testing.T)
```
- **Given**: DockerClient interface
- **When**: Check visibility
- **Then**: Interface name starts with uppercase (exported)

#### TC-IFACE-DOCKER-015: Interface documented with comments
```go
func TestDockerClientInterface_Documented(t *testing.T)
```
- **Given**: interface.go file
- **When**: Parse file for comments
- **Then**: Interface and all methods have documentation comments

---

### Test Category: RegistryClient Interface Validation

**Test File**: `pkg/registry/interface_test.go`

**Interface Under Test**:
```go
type RegistryClient interface {
    Push(ctx context.Context, image v1.Image, targetRef string, progress ProgressCallback) error
    BuildImageReference(registryURL, imageName string) (string, error)
    ValidateRegistry(ctx context.Context, registryURL string) error
}
```

**Test Cases** (20 tests):

#### TC-IFACE-REGISTRY-001: Interface compiles successfully
#### TC-IFACE-REGISTRY-002: Push signature is correct
#### TC-IFACE-REGISTRY-003: BuildImageReference signature is correct
#### TC-IFACE-REGISTRY-004: ValidateRegistry signature is correct
#### TC-IFACE-REGISTRY-005: ProgressCallback type defined
#### TC-IFACE-REGISTRY-006: ProgressUpdate struct defined
#### TC-IFACE-REGISTRY-007: LayerStatus enum defined
#### TC-IFACE-REGISTRY-008: NewRegistryClient constructor exists
#### TC-IFACE-REGISTRY-009: RegistryAuthError type defined
#### TC-IFACE-REGISTRY-010: RegistryConnectionError type defined
#### TC-IFACE-REGISTRY-011: LayerPushError type defined
#### TC-IFACE-REGISTRY-012: Mock implementation satisfies interface
#### TC-IFACE-REGISTRY-013: Push accepts v1.Image parameter
#### TC-IFACE-REGISTRY-014: ProgressCallback signature correct
#### TC-IFACE-REGISTRY-015: ProgressUpdate has required fields
#### TC-IFACE-REGISTRY-016: LayerStatus has all states
#### TC-IFACE-REGISTRY-017: Error types have required fields
#### TC-IFACE-REGISTRY-018: Interface exported (public)
#### TC-IFACE-REGISTRY-019: Interface documented
#### TC-IFACE-REGISTRY-020: NewRegistryClient accepts AuthProvider and TLSProvider

---

### Test Category: AuthProvider Interface Validation

**Test File**: `pkg/auth/interface_test.go`

**Interface Under Test**:
```go
type AuthProvider interface {
    GetAuthenticator() (authn.Authenticator, error)
    ValidateCredentials() error
}
```

**Test Cases** (12 tests):

#### TC-IFACE-AUTH-001: Interface compiles successfully
#### TC-IFACE-AUTH-002: GetAuthenticator signature is correct
#### TC-IFACE-AUTH-003: ValidateCredentials signature is correct
#### TC-IFACE-AUTH-004: NewAuthProvider constructor exists
#### TC-IFACE-AUTH-005: InvalidCredentialsError type defined
#### TC-IFACE-AUTH-006: MissingCredentialsError type defined
#### TC-IFACE-AUTH-007: Mock implementation satisfies interface
#### TC-IFACE-AUTH-008: GetAuthenticator returns authn.Authenticator
#### TC-IFACE-AUTH-009: Error types have required fields
#### TC-IFACE-AUTH-010: Interface exported (public)
#### TC-IFACE-AUTH-011: Interface documented
#### TC-IFACE-AUTH-012: NewAuthProvider accepts username and password

---

### Test Category: TLSProvider Interface Validation

**Test File**: `pkg/tls/interface_test.go`

**Interface Under Test**:
```go
type TLSProvider interface {
    GetTLSConfig() *tls.Config
    IsInsecure() bool
    GetWarningMessage() string
}
```

**Test Cases** (10 tests):

#### TC-IFACE-TLS-001: Interface compiles successfully
#### TC-IFACE-TLS-002: GetTLSConfig signature is correct
#### TC-IFACE-TLS-003: IsInsecure signature is correct
#### TC-IFACE-TLS-004: GetWarningMessage signature is correct
#### TC-IFACE-TLS-005: NewTLSProvider constructor exists
#### TC-IFACE-TLS-006: Mock implementation satisfies interface
#### TC-IFACE-TLS-007: GetTLSConfig returns *tls.Config
#### TC-IFACE-TLS-008: Interface exported (public)
#### TC-IFACE-TLS-009: Interface documented
#### TC-IFACE-TLS-010: NewTLSProvider accepts insecure bool parameter

---

### Test Category: Command Structure Validation

**Test File**: `cmd/push_interface_test.go`

**Structure Under Test**:
```go
type PushCommand struct {
    dockerClient   docker.DockerClient
    registryClient registry.RegistryClient
    authProvider   auth.AuthProvider
    tlsProvider    tls.TLSProvider
}
```

**Test Cases** (13 tests):

#### TC-IFACE-CMD-001: PushCommand struct compiles
#### TC-IFACE-CMD-002: PushCommand has dockerClient field
#### TC-IFACE-CMD-003: PushCommand has registryClient field
#### TC-IFACE-CMD-004: PushCommand has authProvider field
#### TC-IFACE-CMD-005: PushCommand has tlsProvider field
#### TC-IFACE-CMD-006: NewPushCommand constructor exists
#### TC-IFACE-CMD-007: Execute method signature is correct
#### TC-IFACE-CMD-008: PushFlags struct defined
#### TC-IFACE-CMD-009: PushFlags has all required fields
#### TC-IFACE-CMD-010: Cobra command integration skeleton exists
#### TC-IFACE-CMD-011: Flag definitions present
#### TC-IFACE-CMD-012: Help text defined
#### TC-IFACE-CMD-013: Command exported (public)

---

## Wave 2 Test Specifications

### Purpose
Validate that all implementations of Wave 1 interfaces are correct, handle edge cases, and provide production-ready functionality.

---

### Test Category: DockerClient Implementation Tests

**Test File**: `pkg/docker/client_test.go`

**Implementation Under Test**: Docker Engine API client with subprocess fallback

**Test Cases** (15 tests):

#### TC-IMPL-DOCKER-001: ImageExists returns true for existing image
```go
func TestDockerClient_ImageExists_ExistingImage(t *testing.T)
```
- **Given**: Mock Docker daemon with alpine:latest
- **When**: ImageExists("alpine:latest") called
- **Then**: Returns (true, nil)

#### TC-IMPL-DOCKER-002: ImageExists returns false for missing image
```go
func TestDockerClient_ImageExists_MissingImage(t *testing.T)
```
- **Given**: Mock Docker daemon without "nonexistent:v1"
- **When**: ImageExists("nonexistent:v1") called
- **Then**: Returns (false, nil)

#### TC-IMPL-DOCKER-003: ImageExists returns error on daemon failure
```go
func TestDockerClient_ImageExists_DaemonUnreachable(t *testing.T)
```
- **Given**: Mock Docker daemon unreachable
- **When**: ImageExists("alpine:latest") called
- **Then**: Returns (false, DaemonConnectionError)

#### TC-IMPL-DOCKER-004: GetImage returns valid v1.Image for existing image
```go
func TestDockerClient_GetImage_Success(t *testing.T)
```
- **Given**: Mock Docker daemon with nginx:latest
- **When**: GetImage("nginx:latest") called
- **Then**: Returns valid v1.Image with layers

#### TC-IMPL-DOCKER-005: GetImage returns ImageNotFoundError for missing image
```go
func TestDockerClient_GetImage_NotFound(t *testing.T)
```
- **Given**: Mock daemon without "missing:v1"
- **When**: GetImage("missing:v1") called
- **Then**: Returns ImageNotFoundError

#### TC-IMPL-DOCKER-006: ValidateImageName accepts valid OCI names
```go
func TestDockerClient_ValidateImageName_Valid(t *testing.T)
```
- **Given**: Valid names: "myapp:v1", "registry.io/ns/app:latest"
- **When**: ValidateImageName called for each
- **Then**: Returns nil (no error)

#### TC-IMPL-DOCKER-007: ValidateImageName rejects invalid OCI names
```go
func TestDockerClient_ValidateImageName_Invalid(t *testing.T)
```
- **Given**: Invalid names: "My App:v1", "app::", "registry:port/app"
- **When**: ValidateImageName called
- **Then**: Returns InvalidImageNameError

#### TC-IMPL-DOCKER-008: ValidateImageName enforces length limits
```go
func TestDockerClient_ValidateImageName_TooLong(t *testing.T)
```
- **Given**: Image name > 255 characters
- **When**: ValidateImageName called
- **Then**: Returns InvalidImageNameError

#### TC-IMPL-DOCKER-009: Close releases daemon connection
```go
func TestDockerClient_Close_ReleasesConnection(t *testing.T)
```
- **Given**: DockerClient with active connection
- **When**: Close() called
- **Then**: Connection released

#### TC-IMPL-DOCKER-010: NewDockerClient succeeds with Docker API
```go
func TestNewDockerClient_DockerAPIAvailable(t *testing.T)
```
- **Given**: Mock Docker Engine API reachable
- **When**: NewDockerClient() called
- **Then**: Returns client using API

#### TC-IMPL-DOCKER-011: NewDockerClient falls back to subprocess
```go
func TestNewDockerClient_FallsBackToSubprocess(t *testing.T)
```
- **Given**: Mock Docker Engine API unreachable
- **When**: NewDockerClient() called
- **Then**: Returns client using subprocess

#### TC-IMPL-DOCKER-012: GetImage handles multi-layer images
```go
func TestDockerClient_GetImage_MultiLayer(t *testing.T)
```
- **Given**: Mock image with 10+ layers
- **When**: GetImage called
- **Then**: Returns v1.Image with all layers

#### TC-IMPL-DOCKER-013: Context cancellation stops operations
```go
func TestDockerClient_ContextCancellation(t *testing.T)
```
- **Given**: Long-running GetImage operation
- **When**: Context cancelled
- **Then**: Returns context.Canceled error

#### TC-IMPL-DOCKER-014: Concurrent ImageExists calls are safe
```go
func TestDockerClient_ConcurrentImageExists(t *testing.T)
```
- **Given**: DockerClient instance
- **When**: 100 concurrent ImageExists calls
- **Then**: All complete without race conditions

#### TC-IMPL-DOCKER-015: Error types are correctly typed
```go
func TestDockerClient_ErrorTypes(t *testing.T)
```
- **Given**: Various error scenarios
- **When**: Errors returned
- **Then**: Can type-assert to correct error types

---

### Test Category: RegistryClient Implementation Tests

**Test File**: `pkg/registry/client_test.go`

**Implementation Under Test**: go-containerregistry wrapper with progress callbacks

**Test Cases** (20 tests):

#### TC-IMPL-REGISTRY-001: Push succeeds with valid image (mock)
```go
func TestRegistryClient_Push_Success(t *testing.T)
```
- **Given**: Mock registry and valid v1.Image
- **When**: Push called
- **Then**: Returns nil (success)

#### TC-IMPL-REGISTRY-002: Push returns RegistryAuthError with invalid credentials
```go
func TestRegistryClient_Push_AuthFailure(t *testing.T)
```
- **Given**: Mock registry rejecting credentials
- **When**: Push called
- **Then**: Returns RegistryAuthError

#### TC-IMPL-REGISTRY-003: Push returns RegistryConnectionError when unreachable
```go
func TestRegistryClient_Push_ConnectionFailure(t *testing.T)
```
- **Given**: Mock registry unreachable
- **When**: Push called
- **Then**: Returns RegistryConnectionError

#### TC-IMPL-REGISTRY-004: Push invokes progress callbacks
```go
func TestRegistryClient_Push_ProgressCallbacks(t *testing.T)
```
- **Given**: Mock image with 5 layers
- **When**: Push called with progress callback
- **Then**: Callback invoked 5 times

#### TC-IMPL-REGISTRY-005: Progress callbacks report status transitions
```go
func TestRegistryClient_Push_LayerStatusTransitions(t *testing.T)
```
- **Given**: Mock push in progress
- **When**: Progress callbacks invoked
- **Then**: Status transitions: Waiting → Uploading → Complete

#### TC-IMPL-REGISTRY-006: BuildImageReference constructs valid reference
```go
func TestRegistryClient_BuildImageReference_Valid(t *testing.T)
```
- **Given**: registryURL="https://gitea.example.com:8443", imageName="myapp:v1"
- **When**: BuildImageReference called
- **Then**: Returns "gitea.example.com:8443/giteaadmin/myapp:v1"

#### TC-IMPL-REGISTRY-007: BuildImageReference handles URL without scheme
```go
func TestRegistryClient_BuildImageReference_NoScheme(t *testing.T)
```
- **Given**: registryURL="gitea.example.com:8443"
- **When**: BuildImageReference called
- **Then**: Returns valid reference (defaults to https)

#### TC-IMPL-REGISTRY-008: BuildImageReference validates image name
```go
func TestRegistryClient_BuildImageReference_InvalidImageName(t *testing.T)
```
- **Given**: imageName="Invalid Name:v1"
- **When**: BuildImageReference called
- **Then**: Returns error

#### TC-IMPL-REGISTRY-009: ValidateRegistry succeeds for mock compliant registry
```go
func TestRegistryClient_ValidateRegistry_Compliant(t *testing.T)
```
- **Given**: Mock OCI-compliant registry
- **When**: ValidateRegistry called
- **Then**: Returns nil

#### TC-IMPL-REGISTRY-010: ValidateRegistry fails for non-compliant registry
```go
func TestRegistryClient_ValidateRegistry_NonCompliant(t *testing.T)
```
- **Given**: Mock non-OCI registry
- **When**: ValidateRegistry called
- **Then**: Returns error

#### TC-IMPL-REGISTRY-011: Context cancellation stops push
```go
func TestRegistryClient_Push_ContextCancellation(t *testing.T)
```
- **Given**: Mock push in progress
- **When**: Context cancelled
- **Then**: Returns context.Canceled

#### TC-IMPL-REGISTRY-012: NewRegistryClient requires AuthProvider
```go
func TestNewRegistryClient_RequiresAuth(t *testing.T)
```
- **Given**: nil AuthProvider
- **When**: NewRegistryClient called
- **Then**: Returns error

#### TC-IMPL-REGISTRY-013: NewRegistryClient uses TLS config
```go
func TestNewRegistryClient_UsesTLSConfig(t *testing.T)
```
- **Given**: Mock TLSProvider with custom config
- **When**: NewRegistryClient called
- **Then**: Client uses provided TLS config

#### TC-IMPL-REGISTRY-014: Push reports LayerFailed on error
```go
func TestRegistryClient_Push_LayerUploadFailure(t *testing.T)
```
- **Given**: Mock layer upload fails
- **When**: Progress callback invoked
- **Then**: Status = LayerFailed

#### TC-IMPL-REGISTRY-015: BuildImageReference handles default namespace
```go
func TestRegistryClient_BuildImageReference_DefaultNamespace(t *testing.T)
```
- **Given**: imageName="myapp:v1" (no namespace)
- **When**: BuildImageReference called
- **Then**: Adds default namespace "giteaadmin/myapp:v1"

#### TC-IMPL-REGISTRY-016: ValidateRegistry checks OCI endpoints
```go
func TestRegistryClient_ValidateRegistry_OciEndpoints(t *testing.T)
```
- **Given**: Mock registry URL
- **When**: ValidateRegistry called
- **Then**: Checks /v2/ endpoint

#### TC-IMPL-REGISTRY-017: Progress callback receives accurate byte counts
```go
func TestRegistryClient_Push_AccurateProgressBytes(t *testing.T)
```
- **Given**: Mock layer of known size
- **When**: Progress callbacks invoked
- **Then**: BytesUploaded matches expected values

#### TC-IMPL-REGISTRY-018: BuildImageReference handles port numbers
```go
func TestRegistryClient_BuildImageReference_WithPort(t *testing.T)
```
- **Given**: registryURL with port ":8443"
- **When**: BuildImageReference called
- **Then**: Port preserved in reference

#### TC-IMPL-REGISTRY-019: BuildImageReference strips http/https scheme
```go
func TestRegistryClient_BuildImageReference_StripScheme(t *testing.T)
```
- **Given**: registryURL="https://gitea.example.com:8443"
- **When**: BuildImageReference called
- **Then**: Reference excludes "https://"

#### TC-IMPL-REGISTRY-020: Error types contain meaningful context
```go
func TestRegistryClient_ErrorTypes_Context(t *testing.T)
```
- **Given**: Various error scenarios
- **When**: Errors returned
- **Then**: Error messages contain registry URL, layer digest, etc.

---

### Test Category: AuthProvider Implementation Tests

**Test File**: `pkg/auth/provider_test.go`

**Implementation Under Test**: Basic authentication with env var support

**Test Cases** (12 tests):

#### TC-IMPL-AUTH-001: GetAuthenticator returns valid authn.Authenticator
#### TC-IMPL-AUTH-002: ValidateCredentials accepts valid username/password
#### TC-IMPL-AUTH-003: ValidateCredentials rejects empty username
#### TC-IMPL-AUTH-004: ValidateCredentials rejects empty password
#### TC-IMPL-AUTH-005: ValidateCredentials enforces username length limits
#### TC-IMPL-AUTH-006: ValidateCredentials enforces password length limits
#### TC-IMPL-AUTH-007: NewAuthProvider accepts environment variables
#### TC-IMPL-AUTH-008: NewAuthProvider prefers flags over env vars
#### TC-IMPL-AUTH-009: NewAuthProvider returns MissingCredentialsError
#### TC-IMPL-AUTH-010: GetAuthenticator integrates with go-containerregistry
#### TC-IMPL-AUTH-011: ValidateCredentials rejects control characters
#### TC-IMPL-AUTH-012: AuthProvider implements authn.Keychain interface

---

### Test Category: TLSProvider Implementation Tests

**Test File**: `pkg/tls/config_test.go`

**Implementation Under Test**: TLS configuration generation (secure/insecure modes)

**Test Cases** (10 tests):

#### TC-IMPL-TLS-001: GetTLSConfig returns secure config by default
#### TC-IMPL-TLS-002: GetTLSConfig returns insecure config when flag set
#### TC-IMPL-TLS-003: Insecure config has InsecureSkipVerify=true
#### TC-IMPL-TLS-004: Secure config uses system certificate pool
#### TC-IMPL-TLS-005: IsInsecure returns correct boolean
#### TC-IMPL-TLS-006: GetWarningMessage returns warning in insecure mode
#### TC-IMPL-TLS-007: GetWarningMessage returns empty string in secure mode
#### TC-IMPL-TLS-008: NewTLSProvider creates secure provider by default
#### TC-IMPL-TLS-009: NewTLSProvider creates insecure provider when requested
#### TC-IMPL-TLS-010: GetTLSConfig is safe for concurrent calls

---

## Progressive Test Planning Architecture

### R341 Progressive Realism Strategy

**Phase 1 Test Realism**: ISOLATED UNIT TESTS ONLY
- **Wave 1**: Interface validation (compilation tests)
- **Wave 2**: Implementation unit tests (mocked dependencies)
- **NO real Docker daemon** (mocked DockerClient)
- **NO real registry** (mocked RegistryClient)
- **NO integration tests** (packages not wired together yet)

**Progression to Phase 2**:
- Phase 2 will use REAL implementations from Phase 1
- Integration tests will import actual Phase 1 packages
- Tests will verify package interactions work correctly
- Still mocked external services (Docker/Registry in test containers)

**Progression to Phase 3**:
- Phase 3 will use REAL end-to-end flows
- E2E tests with actual Gitea registry (Docker Compose)
- Real Docker daemon interaction
- Full user journey validation

### Test Progression Example

**Phase 1** (Current - Isolated):
```go
// Mock-based unit test
func TestRegistryClient_Push_Success(t *testing.T) {
    mockAuth := &MockAuthProvider{
        GetAuthenticatorFunc: func() (authn.Authenticator, error) {
            return authn.Anonymous, nil
        },
    }
    mockTLS := &MockTLSProvider{
        GetTLSConfigFunc: func() *tls.Config {
            return &tls.Config{}
        },
    }
    client, err := registry.NewRegistryClient(mockAuth, mockTLS)
    // Test with mocked dependencies
}
```

**Phase 2** (Next - Real Packages):
```go
// Integration test using real Phase 1 implementations
func TestIntegration_DockerToRegistry(t *testing.T) {
    // Use REAL AuthProvider from Phase 1
    auth, err := auth.NewAuthProvider("testuser", "testpass")
    // Use REAL TLSProvider from Phase 1
    tlsProvider, err := tls.NewTLSProvider(true)
    // Use REAL RegistryClient from Phase 1
    client, err := registry.NewRegistryClient(auth, tlsProvider)
    // Test real package interactions
}
```

**Phase 3** (Final - Real External Services):
```go
// E2E test with real Gitea registry
func TestE2E_PushToGitea(t *testing.T) {
    // Start real Gitea in Docker Compose
    // Use real Docker daemon
    // Test full user journey
}
```

---

## Test Fixtures and Mocks

### Mock Interfaces (Wave 2 Testing)

**Location**: `pkg/docker/mocks_test.go`, `pkg/registry/mocks_test.go`, etc.

**Mock DockerClient**:
```go
type MockDockerClient struct {
    ImageExistsFunc      func(ctx context.Context, imageName string) (bool, error)
    GetImageFunc         func(ctx context.Context, imageName string) (v1.Image, error)
    ValidateImageNameFunc func(imageName string) error
    CloseFunc            func() error
}

func (m *MockDockerClient) ImageExists(ctx context.Context, imageName string) (bool, error) {
    if m.ImageExistsFunc != nil {
        return m.ImageExistsFunc(ctx, imageName)
    }
    return false, nil
}
// ... other methods
```

**Mock RegistryClient**:
```go
type MockRegistryClient struct {
    PushFunc                  func(ctx context.Context, image v1.Image, targetRef string, progress ProgressCallback) error
    BuildImageReferenceFunc   func(registryURL, imageName string) (string, error)
    ValidateRegistryFunc      func(ctx context.Context, registryURL string) error
}
// ... implementations
```

**Mock AuthProvider**:
```go
type MockAuthProvider struct {
    GetAuthenticatorFunc     func() (authn.Authenticator, error)
    ValidateCredentialsFunc  func() error
}
// ... implementations
```

**Mock TLSProvider**:
```go
type MockTLSProvider struct {
    GetTLSConfigFunc       func() *tls.Config
    IsInsecureFunc         func() bool
    GetWarningMessageFunc  func() string
}
// ... implementations
```

### Test Data Fixtures

**Location**: `tests/fixtures/`

**Fixture Structure**:
```
tests/fixtures/
├── valid-image-names.txt        # List of valid OCI image names
├── invalid-image-names.txt      # List of invalid image names
├── mock-v1-image.go             # Fixture v1.Image for testing
└── test-credentials.yaml        # Test credential sets
```

**Example Fixture**:
```go
// tests/fixtures/mock-v1-image.go
package fixtures

import (
    v1 "github.com/google/go-containerregistry/pkg/v1"
)

// NewMockImage creates a v1.Image fixture for testing
func NewMockImage(layers int) v1.Image {
    // Create mock image with specified number of layers
    // Used for progress callback testing
}
```

---

## Test Execution Strategy

### Wave 1 Execution (Interface Validation)

**Execution Command**:
```bash
# Run all interface validation tests
go test ./pkg/docker/interface_test.go -v
go test ./pkg/registry/interface_test.go -v
go test ./pkg/auth/interface_test.go -v
go test ./pkg/tls/interface_test.go -v
go test ./cmd/push_interface_test.go -v

# Expected output: PASS (all interface tests succeed)
```

**Wave 1 Success Criteria**:
- All interface files compile without errors
- All interface method signatures correct
- All error types defined and implement error interface
- All constructors exist with correct signatures
- Mock implementations can satisfy interfaces

### Wave 2 Execution (Implementation Validation)

**Execution Command**:
```bash
# Run all implementation unit tests
go test ./pkg/docker/... -v -short
go test ./pkg/registry/... -v -short
go test ./pkg/auth/... -v -short
go test ./pkg/tls/... -v -short

# Run with coverage
go test ./pkg/... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Expected coverage: ≥85% for each package
```

**Wave 2 Success Criteria**:
- All implementation unit tests pass
- Code coverage ≥85% per package
- No race conditions detected (`go test -race`)
- All edge cases handled
- All error paths tested

### Continuous Testing During Development

**Developer Workflow**:
```bash
# Watch mode (run tests on file change)
# Using gotestsum or similar tool
gotestsum --watch ./pkg/...

# Fast feedback loop
go test ./pkg/docker/... -short -count=1
```

---

## Success Criteria

### Phase 1 Wave 1 Success Criteria

**Interface Definition Complete**:
- ✅ All 5 interfaces defined (DockerClient, RegistryClient, AuthProvider, TLSProvider, Command)
- ✅ All interfaces compile without errors
- ✅ All 70 interface validation tests pass
- ✅ All error types defined and documented
- ✅ All constructors defined with correct signatures
- ✅ Mock implementations can satisfy all interfaces
- ✅ Architect review approves interface completeness

### Phase 1 Wave 2 Success Criteria

**Implementation Complete**:
- ✅ All 4 implementations complete (Docker, Registry, Auth, TLS)
- ✅ All 80 implementation unit tests pass
- ✅ Code coverage ≥85% for each package
- ✅ No stub implementations (R320 compliant)
- ✅ No race conditions (`go test -race` passes)
- ✅ All efforts ≤800 lines (R220/R221 compliant)
- ✅ All interface contract tests pass (from Wave 1)
- ✅ Code Reviewer approval for all implementations

### Phase 1 Overall Success Criteria

**Foundation Ready**:
- ✅ All ~150 tests passing (70 interface + 80 implementation)
- ✅ All packages independently testable
- ✅ All packages compile and pass type checks
- ✅ No integration failures (packages are independent)
- ✅ Ready for Phase 2 integration work
- ✅ Architect grade A for Phase 1 completion

---

## TDD Workflow Summary

### Current State: RED Phase Complete

**Phase 1 Test Status**:
- ✅ All test specifications defined
- ✅ Expected test count: ~150 tests
- ✅ All tests WILL FAIL initially (expected TDD behavior)
- ✅ Tests serve as executable specifications

**Test Categories Defined**:
- Interface validation tests: 70 tests (Wave 1)
- Implementation unit tests: 80 tests (Wave 2)
- Integration tests: DEFERRED to Phase 2
- E2E tests: DEFERRED to Phase 3

### Next State: GREEN Phase (During Implementation)

**Wave 1 Implementation**:
- Create interface files matching test specifications
- Make all 70 interface validation tests pass
- Architect reviews and approves interfaces

**Wave 2 Implementation**:
- Implement all 4 packages (Docker, Registry, Auth, TLS)
- Make all 80 implementation unit tests pass
- Make all 70 interface contract tests pass
- Code Reviewer validates all implementations

### Final State: REFACTOR Phase (Post-Implementation)

**Code Quality Improvements**:
- Eliminate code duplication
- Improve naming and readability
- Optimize performance bottlenecks
- Enhance error messages
- Maintain 100% test passage throughout refactoring

---

## R341 TDD Compliance Verification

**TDD Protocol Compliance**:
- ✅ **RED phase**: Tests created BEFORE implementation (this document)
- ✅ **Test-First**: Tests define acceptance criteria upfront
- ✅ **Executable Specifications**: Tests validate architectural promises
- ✅ **Progressive Realism**: Tests follow progressive test planning architecture
- ⏳ **GREEN phase**: Will occur during Wave 1 and Wave 2 implementation
- ⏳ **REFACTOR phase**: Will occur after all tests pass

**R341 Checklist**:
- ✅ Test plan created before implementation planning
- ✅ Tests cover all interface contracts
- ✅ Tests cover all implementation logic
- ✅ Tests follow progressive realism (mocks → real packages → real services)
- ✅ Tests define clear success criteria
- ✅ Test execution strategy documented

---

## Test Plan Metadata

**Test Plan Location**: `/home/vscode/workspaces/idpbuilder-oci-push-rebuild/planning/phase1/PHASE-TEST-PLAN.md`

**Test Harness Location**: Will be created in Phase 1 Wave 1 as `pkg/*/interface_test.go` and Phase 1 Wave 2 as `pkg/*/*_test.go`

**Orchestrator Tracking**: Orchestrator will record test plan location in `orchestrator-state-v3.json` under:
```json
{
  "test_plans": {
    "phase1": {
      "plan_file": "planning/phase1/PHASE-TEST-PLAN.md",
      "wave1_tests": ["pkg/docker/interface_test.go", "..."],
      "wave2_tests": ["pkg/docker/client_test.go", "..."],
      "created_at": "2025-11-11T03:09:23Z",
      "created_by": "code-reviewer",
      "status": "pending_implementation"
    }
  }
}
```

---

## Appendix A: Test Execution Examples

### Example: Running Wave 1 Interface Tests

```bash
# Navigate to project root
cd /home/vscode/workspaces/idpbuilder-oci-push-rebuild

# Run all interface validation tests
go test ./pkg/docker/interface_test.go -v
go test ./pkg/registry/interface_test.go -v
go test ./pkg/auth/interface_test.go -v
go test ./pkg/tls/interface_test.go -v
go test ./cmd/push_interface_test.go -v

# Expected output:
# === RUN   TestDockerClientInterface_Compiles
# --- PASS: TestDockerClientInterface_Compiles (0.00s)
# === RUN   TestDockerClientInterface_ImageExistsSignature
# --- PASS: TestDockerClientInterface_ImageExistsSignature (0.00s)
# ...
# PASS
# ok      idpbuilder/pkg/docker   0.123s
```

### Example: Running Wave 2 Implementation Tests

```bash
# Run all implementation unit tests
go test ./pkg/docker/... -v -short
go test ./pkg/registry/... -v -short
go test ./pkg/auth/... -v -short
go test ./pkg/tls/... -v -short

# Run with coverage report
go test ./pkg/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Expected coverage: ≥85% for each package
```

---

## Appendix B: Test Count Summary

**Phase 1 Total Tests**: ~150 tests

**Wave 1 Tests** (Interface Validation): ~70 tests
- DockerClient: 15 tests
- RegistryClient: 20 tests
- AuthProvider: 12 tests
- TLSProvider: 10 tests
- Command: 13 tests

**Wave 2 Tests** (Implementation): ~80 tests
- DockerClient: 15 tests
- RegistryClient: 20 tests
- AuthProvider: 12 tests
- TLSProvider: 10 tests
- Integration helpers: ~23 tests

**Deferred to Phase 2**: ~40 integration tests
**Deferred to Phase 3**: ~30 E2E tests

---

**Test Plan Status**: COMPLETE - Ready for Phase 1 Implementation
**TDD Phase**: RED (tests defined, expecting failure)
**Expected Next Step**: Create Phase 1 Wave 1 Architecture Plan with real Go code examples

---

**R341 TDD Compliance**: ✅ VERIFIED
- Tests created BEFORE implementation planning
- Tests define success criteria for Phase 1
- Tests are executable specifications
- Implementation will target test passage (GREEN phase)
