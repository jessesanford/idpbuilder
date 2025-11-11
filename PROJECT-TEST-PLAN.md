# PROJECT TEST PLAN: idpbuilder-oci-push-rebuild

**Created**: 2025-11-11
**Project**: idpbuilder OCI Push Enhancement
**Test Planning State**: PROJECT_TEST_PLANNING
**TDD Workflow**: RED phase (tests created before implementation)
**Architecture Source**: planning/project/PROJECT-ARCHITECTURE-PLAN.md

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [TDD Approach](#tdd-approach)
3. [Test Strategy](#test-strategy)
4. [Interface Contract Tests](#interface-contract-tests)
5. [Functional Integration Tests](#functional-integration-tests)
6. [End-to-End Workflow Tests](#end-to-end-workflow-tests)
7. [Quality Attribute Tests](#quality-attribute-tests)
8. [Test-to-Phase Mapping](#test-to-phase-mapping)
9. [Demo Scenarios](#demo-scenarios)
10. [Success Criteria](#success-criteria)

---

## Executive Summary

### Test Plan Purpose

This test plan defines the **comprehensive acceptance criteria** for the idpbuilder OCI push enhancement project using Test-Driven Development (TDD). Tests are created **before implementation** to serve as executable specifications that validate:

- **Interface contracts** are correctly defined and implemented
- **Integration between packages** works as designed
- **End-to-end workflows** deliver user value
- **Quality attributes** (performance, security, reliability) are met
- **Architectural promises** from master architecture are fulfilled

### Test Coverage Strategy

**Total Test Categories**: 5
1. **Interface Contract Tests** (Phase 1 validation)
2. **Functional Integration Tests** (Phase 2 validation)
3. **End-to-End Workflow Tests** (Phase 3 validation)
4. **Quality Attribute Tests** (Cross-phase validation)
5. **Demo Scenarios** (User acceptance validation)

**Expected Test Count**: ~150 tests
- Unit tests: ~80 tests (interface compliance, validation logic)
- Integration tests: ~40 tests (package interactions)
- E2E tests: ~20 tests (user workflows)
- Quality tests: ~10 tests (performance, security)

### TDD Workflow Compliance (R341)

This plan follows R341 TDD protocol:
- ✅ **RED phase**: Tests created BEFORE implementation planning
- ⏳ **GREEN phase**: Implementation will target making tests pass
- ⏳ **REFACTOR phase**: Code improvements while maintaining test passage

**Current Status**: RED phase (all tests expected to fail initially)

---

## TDD Approach

### Test-First Development Model

```
┌────────────────────────────────────────────────────────────────┐
│  Phase 0: TEST CREATION (THIS DOCUMENT)                         │
│  ────────────────────────────────────────────────────────────   │
│  - Define interface contract tests                              │
│  - Define integration tests                                     │
│  - Define E2E workflow tests                                    │
│  - Define quality attribute tests                               │
│  - Expected: ALL TESTS FAIL (red phase)                         │
└─────────────────────┬──────────────────────────────────────────┘
                      │
                      ▼
┌────────────────────────────────────────────────────────────────┐
│  Phase 1-3: IMPLEMENTATION (FUTURE)                             │
│  ────────────────────────────────────────────────────────────   │
│  - Implement features to make tests pass                        │
│  - Each phase targets specific test categories                  │
│  - Refactor while maintaining test passage                      │
│  - Goal: ALL TESTS PASS (green phase)                           │
└────────────────────────────────────────────────────────────────┘
```

### Test Categories by TDD Phase

**RED Phase (Current)**:
- All tests created but not passing
- Tests define "done" criteria
- Tests are executable specifications
- Tests validate architectural promises

**GREEN Phase (During Implementation)**:
- Phase 1: Make interface contract tests pass
- Phase 2: Make integration tests pass
- Phase 3: Make E2E and quality tests pass

**REFACTOR Phase (Post-Implementation)**:
- Improve code while maintaining test passage
- Eliminate duplication
- Enhance readability
- Optimize performance

---

## Test Strategy

### Testing Pyramid

```
                    ┌──────────────┐
                    │   E2E Tests  │  ~20 tests (13%)
                    │  (Slow, High │
                    │   Fidelity)  │
                  ┌─┴──────────────┴─┐
                  │  Integration     │  ~40 tests (27%)
                  │     Tests        │
                  │  (Medium Speed,  │
                  │   Med Fidelity)  │
                ┌─┴──────────────────┴─┐
                │    Unit Tests         │  ~80 tests (53%)
                │  (Fast, Interface     │
                │   Validation)         │
                └───────────────────────┘
                ~10 Quality Tests (7%)
```

### Test Execution Strategy

**Development Cycle** (fastest feedback):
```bash
# Run unit tests only (< 5 seconds)
go test ./pkg/... -short
```

**Pre-commit Validation** (medium feedback):
```bash
# Run unit + integration tests (< 30 seconds)
go test ./pkg/... ./tests/integration/...
```

**CI/CD Pipeline** (complete validation):
```bash
# Run all tests including E2E (< 5 minutes)
./PROJECT-TEST-HARNESS.sh
```

### Test Environment Requirements

**Unit Tests**:
- No external dependencies
- Mock all interfaces
- Run in isolation

**Integration Tests**:
- Docker daemon required (local)
- Test registry (Docker Compose)
- Network access

**E2E Tests**:
- Full IDPBuilder environment
- Gitea registry (Docker Compose)
- Real Docker images
- TLS certificates

---

## Interface Contract Tests

### Purpose
Validate that all interfaces defined in Phase 1 Wave 1 are correctly implemented and comply with architectural contracts.

### Test Category: DockerClient Interface

**Test File**: `pkg/docker/client_interface_test.go`

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

#### TC-DOCKER-001: ImageExists returns true for existing image
```go
func TestDockerClient_ImageExists_ExistingImage(t *testing.T)
```
- **Given**: Docker daemon contains image "alpine:latest"
- **When**: ImageExists("alpine:latest") called
- **Then**: Returns (true, nil)

#### TC-DOCKER-002: ImageExists returns false for missing image
```go
func TestDockerClient_ImageExists_MissingImage(t *testing.T)
```
- **Given**: Docker daemon does NOT contain "nonexistent:v1"
- **When**: ImageExists("nonexistent:v1") called
- **Then**: Returns (false, nil)

#### TC-DOCKER-003: ImageExists returns error on daemon failure
```go
func TestDockerClient_ImageExists_DaemonUnreachable(t *testing.T)
```
- **Given**: Docker daemon unreachable
- **When**: ImageExists("alpine:latest") called
- **Then**: Returns (false, DaemonConnectionError)

#### TC-DOCKER-004: GetImage returns valid v1.Image for existing image
```go
func TestDockerClient_GetImage_Success(t *testing.T)
```
- **Given**: Docker daemon contains image "nginx:latest"
- **When**: GetImage("nginx:latest") called
- **Then**: Returns valid v1.Image with layers

#### TC-DOCKER-005: GetImage returns ImageNotFoundError for missing image
```go
func TestDockerClient_GetImage_NotFound(t *testing.T)
```
- **Given**: Image "missing:v1" not in daemon
- **When**: GetImage("missing:v1") called
- **Then**: Returns ImageNotFoundError

#### TC-DOCKER-006: ValidateImageName accepts valid OCI names
```go
func TestDockerClient_ValidateImageName_Valid(t *testing.T)
```
- **Given**: Valid image names: "myapp:v1", "registry.io/ns/app:latest", "sha256:abc123"
- **When**: ValidateImageName called for each
- **Then**: Returns nil (no error)

#### TC-DOCKER-007: ValidateImageName rejects invalid OCI names
```go
func TestDockerClient_ValidateImageName_Invalid(t *testing.T)
```
- **Given**: Invalid names: "My App:v1", "app::", "registry:port/app"
- **When**: ValidateImageName called for each
- **Then**: Returns InvalidImageNameError with reason

#### TC-DOCKER-008: ValidateImageName enforces length limits
```go
func TestDockerClient_ValidateImageName_TooLong(t *testing.T)
```
- **Given**: Image name > 255 characters
- **When**: ValidateImageName called
- **Then**: Returns InvalidImageNameError (name too long)

#### TC-DOCKER-009: Close releases daemon connection
```go
func TestDockerClient_Close_ReleasesConnection(t *testing.T)
```
- **Given**: DockerClient with active connection
- **When**: Close() called
- **Then**: Connection released, subsequent calls fail gracefully

#### TC-DOCKER-010: NewDockerClient succeeds with Docker API
```go
func TestNewDockerClient_DockerAPIAvailable(t *testing.T)
```
- **Given**: Docker Engine API reachable
- **When**: NewDockerClient() called
- **Then**: Returns client using API (not subprocess)

#### TC-DOCKER-011: NewDockerClient falls back to subprocess
```go
func TestNewDockerClient_FallsBackToSubprocess(t *testing.T)
```
- **Given**: Docker Engine API unreachable
- **When**: NewDockerClient() called
- **Then**: Returns client using subprocess fallback

#### TC-DOCKER-012: GetImage handles multi-layer images
```go
func TestDockerClient_GetImage_MultiLayer(t *testing.T)
```
- **Given**: Image with 10+ layers
- **When**: GetImage called
- **Then**: Returns v1.Image with all layers accessible

#### TC-DOCKER-013: Context cancellation stops operations
```go
func TestDockerClient_ContextCancellation(t *testing.T)
```
- **Given**: Long-running GetImage operation
- **When**: Context cancelled mid-operation
- **Then**: Returns context.Canceled error immediately

#### TC-DOCKER-014: Concurrent ImageExists calls are safe
```go
func TestDockerClient_ConcurrentImageExists(t *testing.T)
```
- **Given**: DockerClient instance
- **When**: 100 concurrent ImageExists calls
- **Then**: All complete successfully without race conditions

#### TC-DOCKER-015: Error types are correctly typed
```go
func TestDockerClient_ErrorTypes(t *testing.T)
```
- **Given**: Various error scenarios
- **When**: Errors returned
- **Then**: Can type-assert to ImageNotFoundError, DaemonConnectionError, etc.

---

### Test Category: RegistryClient Interface

**Test File**: `pkg/registry/client_interface_test.go`

**Interface Under Test**:
```go
type RegistryClient interface {
    Push(ctx context.Context, image v1.Image, targetRef string, progress ProgressCallback) error
    BuildImageReference(registryURL, imageName string) (string, error)
    ValidateRegistry(ctx context.Context, registryURL string) error
}
```

**Test Cases** (20 tests):

#### TC-REGISTRY-001: Push succeeds with valid image and credentials
```go
func TestRegistryClient_Push_Success(t *testing.T)
```
- **Given**: Valid v1.Image and authenticated registry
- **When**: Push called with valid targetRef
- **Then**: Image pushed successfully, progress callbacks invoked

#### TC-REGISTRY-002: Push returns RegistryAuthError with invalid credentials
```go
func TestRegistryClient_Push_AuthFailure(t *testing.T)
```
- **Given**: Invalid credentials
- **When**: Push called
- **Then**: Returns RegistryAuthError with registry URL

#### TC-REGISTRY-003: Push returns RegistryConnectionError when registry unreachable
```go
func TestRegistryClient_Push_ConnectionFailure(t *testing.T)
```
- **Given**: Registry URL unreachable
- **When**: Push called
- **Then**: Returns RegistryConnectionError

#### TC-REGISTRY-004: Push invokes progress callbacks for each layer
```go
func TestRegistryClient_Push_ProgressCallbacks(t *testing.T)
```
- **Given**: Image with 5 layers
- **When**: Push called with progress callback
- **Then**: Callback invoked 5 times (once per layer)

#### TC-REGISTRY-005: Progress callbacks report layer status transitions
```go
func TestRegistryClient_Push_LayerStatusTransitions(t *testing.T)
```
- **Given**: Push in progress
- **When**: Progress callbacks invoked
- **Then**: Status transitions: Waiting → Uploading → Complete

#### TC-REGISTRY-006: BuildImageReference constructs valid registry reference
```go
func TestRegistryClient_BuildImageReference_Valid(t *testing.T)
```
- **Given**: registryURL="https://gitea.example.com:8443", imageName="myapp:v1"
- **When**: BuildImageReference called
- **Then**: Returns "gitea.example.com:8443/giteaadmin/myapp:v1"

#### TC-REGISTRY-007: BuildImageReference handles registry URL without scheme
```go
func TestRegistryClient_BuildImageReference_NoScheme(t *testing.T)
```
- **Given**: registryURL="gitea.example.com:8443"
- **When**: BuildImageReference called
- **Then**: Returns valid reference (defaults to https)

#### TC-REGISTRY-008: BuildImageReference validates image name format
```go
func TestRegistryClient_BuildImageReference_InvalidImageName(t *testing.T)
```
- **Given**: imageName="Invalid Name:v1"
- **When**: BuildImageReference called
- **Then**: Returns error (invalid OCI name)

#### TC-REGISTRY-009: ValidateRegistry succeeds for OCI-compliant registry
```go
func TestRegistryClient_ValidateRegistry_Compliant(t *testing.T)
```
- **Given**: OCI-compliant registry at URL
- **When**: ValidateRegistry called
- **Then**: Returns nil (validation passed)

#### TC-REGISTRY-010: ValidateRegistry fails for non-OCI registry
```go
func TestRegistryClient_ValidateRegistry_NonCompliant(t *testing.T)
```
- **Given**: Non-OCI registry (missing required endpoints)
- **When**: ValidateRegistry called
- **Then**: Returns error describing non-compliance

#### TC-REGISTRY-011: Push handles large images efficiently (streaming)
```go
func TestRegistryClient_Push_LargeImage_NoBuffering(t *testing.T)
```
- **Given**: 500MB image
- **When**: Push called
- **Then**: Memory usage stays < 100MB (streaming, not buffering)

#### TC-REGISTRY-012: Push retries transient network failures
```go
func TestRegistryClient_Push_TransientFailureRetry(t *testing.T)
```
- **Given**: Network hiccup during layer upload
- **When**: Push called
- **Then**: Retries layer upload, succeeds eventually

#### TC-REGISTRY-013: Context cancellation stops push immediately
```go
func TestRegistryClient_Push_ContextCancellation(t *testing.T)
```
- **Given**: Push in progress
- **When**: Context cancelled
- **Then**: Returns context.Canceled, stops immediately

#### TC-REGISTRY-014: NewRegistryClient requires AuthProvider
```go
func TestNewRegistryClient_RequiresAuth(t *testing.T)
```
- **Given**: nil AuthProvider
- **When**: NewRegistryClient called
- **Then**: Returns error (auth required)

#### TC-REGISTRY-015: NewRegistryClient uses TLS config from provider
```go
func TestNewRegistryClient_UsesTLSConfig(t *testing.T)
```
- **Given**: TLSProvider with custom config
- **When**: NewRegistryClient called
- **Then**: Uses provided TLS config for connections

#### TC-REGISTRY-016: Push reports LayerFailed status on upload error
```go
func TestRegistryClient_Push_LayerUploadFailure(t *testing.T)
```
- **Given**: Layer upload fails
- **When**: Progress callback invoked
- **Then**: Status = LayerFailed with error details

#### TC-REGISTRY-017: BuildImageReference handles default namespace
```go
func TestRegistryClient_BuildImageReference_DefaultNamespace(t *testing.T)
```
- **Given**: imageName without namespace "myapp:v1"
- **When**: BuildImageReference called
- **Then**: Adds default namespace "giteaadmin/myapp:v1"

#### TC-REGISTRY-018: ValidateRegistry checks required OCI endpoints
```go
func TestRegistryClient_ValidateRegistry_OciEndpoints(t *testing.T)
```
- **Given**: Registry URL
- **When**: ValidateRegistry called
- **Then**: Checks /v2/ endpoint, _catalog, blob uploads

#### TC-REGISTRY-019: Push handles manifest push after layers
```go
func TestRegistryClient_Push_ManifestAfterLayers(t *testing.T)
```
- **Given**: Image with 3 layers
- **When**: Push called
- **Then**: Layers pushed first, then manifest (correct order)

#### TC-REGISTRY-020: Progress callback receives accurate byte counts
```go
func TestRegistryClient_Push_AccurateProgressBytes(t *testing.T)
```
- **Given**: Layer of known size
- **When**: Progress callbacks invoked
- **Then**: BytesUploaded matches actual bytes sent

---

### Test Category: AuthProvider Interface

**Test File**: `pkg/auth/provider_interface_test.go`

**Test Cases** (12 tests):

#### TC-AUTH-001: GetAuthenticator returns valid authn.Authenticator
#### TC-AUTH-002: ValidateCredentials accepts valid username/password
#### TC-AUTH-003: ValidateCredentials rejects empty username
#### TC-AUTH-004: ValidateCredentials rejects empty password
#### TC-AUTH-005: ValidateCredentials enforces username length limits
#### TC-AUTH-006: ValidateCredentials enforces password length limits
#### TC-AUTH-007: NewAuthProvider accepts environment variables
#### TC-AUTH-008: NewAuthProvider prefers flags over env vars
#### TC-AUTH-009: NewAuthProvider returns MissingCredentialsError when both missing
#### TC-AUTH-010: GetAuthenticator integrates with go-containerregistry
#### TC-AUTH-011: ValidateCredentials rejects control characters
#### TC-AUTH-012: AuthProvider implements authn.Keychain interface

---

### Test Category: TLSProvider Interface

**Test File**: `pkg/tls/config_interface_test.go`

**Test Cases** (10 tests):

#### TC-TLS-001: GetTLSConfig returns secure config by default
#### TC-TLS-002: GetTLSConfig returns insecure config when flag set
#### TC-TLS-003: Insecure config has InsecureSkipVerify=true
#### TC-TLS-004: Secure config uses system certificate pool
#### TC-TLS-005: IsInsecure returns correct boolean
#### TC-TLS-006: GetWarningMessage returns warning in insecure mode
#### TC-TLS-007: GetWarningMessage returns empty string in secure mode
#### TC-TLS-008: NewTLSProvider creates secure provider by default
#### TC-TLS-009: NewTLSProvider creates insecure provider when requested
#### TC-TLS-010: GetTLSConfig is safe for concurrent calls

---

### Test Category: Command Layer Interface

**Test File**: `cmd/push_interface_test.go`

**Test Cases** (13 tests):

#### TC-CMD-001: PushCommand.Execute succeeds with valid inputs
#### TC-CMD-002: PushCommand.Execute validates image name
#### TC-CMD-003: PushCommand.Execute handles missing image
#### TC-CMD-004: PushCommand.Execute handles auth failure
#### TC-CMD-005: PushCommand.Execute handles network failure
#### TC-CMD-006: PushCommand.Execute uses custom registry flag
#### TC-CMD-007: PushCommand.Execute uses environment variables
#### TC-CMD-008: Flags override environment variables
#### TC-CMD-009: Exit code 0 on success
#### TC-CMD-010: Exit code 2 on auth failure
#### TC-CMD-011: Exit code 3 on network failure
#### TC-CMD-012: Exit code 4 on image not found
#### TC-CMD-013: Password redacted in all log output

---

## Functional Integration Tests

### Purpose
Validate that packages work together correctly to deliver integrated functionality.

### Test Category: Docker + Registry Integration

**Test File**: `tests/integration/docker_registry_test.go`

**Test Cases** (10 tests):

#### TC-INT-001: Pull from Docker, push to registry end-to-end
```go
func TestIntegration_DockerToRegistry_Success(t *testing.T)
```
- **Given**: Docker image "alpine:latest" locally
- **When**: GetImage from Docker, Push to test registry
- **Then**: Image successfully pushed, available in registry

#### TC-INT-002: Image with multiple layers transfers correctly
```go
func TestIntegration_MultiLayerImage(t *testing.T)
```
- **Given**: Image with 10 layers
- **When**: Full pull-push cycle
- **Then**: All layers present in registry

#### TC-INT-003: Large image (>100MB) transfers efficiently
```go
func TestIntegration_LargeImageTransfer(t *testing.T)
```
- **Given**: 200MB Docker image
- **When**: Push to registry
- **Then**: Completes in < 60s, memory usage < 150MB

#### TC-INT-004: Progress callbacks fire during actual transfer
```go
func TestIntegration_ProgressCallbacksRealTransfer(t *testing.T)
```
- **Given**: Real image push
- **When**: Progress callbacks monitored
- **Then**: Callbacks invoked with accurate layer counts

#### TC-INT-005: Docker API unavailable falls back to subprocess
```go
func TestIntegration_DockerAPIFallback(t *testing.T)
```
- **Given**: Docker API unreachable, CLI available
- **When**: GetImage called
- **Then**: Uses `docker save` subprocess successfully

#### TC-INT-006: Registry validates pushed image manifest
```go
func TestIntegration_ManifestValidation(t *testing.T)
```
- **Given**: Image pushed to registry
- **When**: Manifest retrieved
- **Then**: Manifest matches source image config

#### TC-INT-007: Concurrent pushes to same registry work
```go
func TestIntegration_ConcurrentPushes(t *testing.T)
```
- **Given**: 5 different images
- **When**: Push all concurrently to same registry
- **Then**: All succeed without conflicts

#### TC-INT-008: Push same image twice is idempotent
```go
func TestIntegration_IdempotentPush(t *testing.T)
```
- **Given**: Image already in registry
- **When**: Push same image again
- **Then**: Succeeds quickly (layers already exist)

#### TC-INT-009: Docker image format incompatibility detected
```go
func TestIntegration_ImageFormatMismatch(t *testing.T)
```
- **Given**: Non-OCI format image
- **When**: Push attempted
- **Then**: Returns clear error about format incompatibility

#### TC-INT-010: Registry returns error for duplicate tag
```go
func TestIntegration_TagConflict(t *testing.T)
```
- **Given**: Tag "v1.0" already exists in registry
- **When**: Push different image with same tag
- **Then**: Overwrites or returns clear conflict error

---

### Test Category: Auth + TLS + Registry Integration

**Test File**: `tests/integration/auth_tls_registry_test.go`

**Test Cases** (10 tests):

#### TC-INT-011: Basic auth credentials work with registry
#### TC-INT-012: Invalid credentials fail with RegistryAuthError
#### TC-INT-013: Insecure TLS connects to self-signed cert registry
#### TC-INT-014: Secure TLS rejects self-signed cert
#### TC-INT-015: Environment variable credentials work
#### TC-INT-016: Flag credentials override env vars
#### TC-INT-017: Auth header includes correct base64 encoding
#### TC-INT-018: TLS warning displayed in insecure mode
#### TC-INT-019: Multiple auth failures logged correctly
#### TC-INT-020: TLS handshake failure returns clear error

---

### Test Category: Full Stack Integration

**Test File**: `tests/integration/full_stack_test.go`

**Test Cases** (10 tests):

#### TC-INT-021: Command → Docker → Auth → TLS → Registry full flow
#### TC-INT-022: Flag parsing to registry push complete cycle
#### TC-INT-023: Progress reporting updates during real push
#### TC-INT-024: Error from any layer propagates correctly
#### TC-INT-025: Exit code matches failure type
#### TC-INT-026: Logs contain actionable error messages
#### TC-INT-027: Cleanup on failure (no partial state)
#### TC-INT-028: Graceful shutdown on SIGINT
#### TC-INT-029: Verbose mode produces detailed logs
#### TC-INT-030: Quiet mode suppresses non-error output

---

## End-to-End Workflow Tests

### Purpose
Validate complete user journeys from terminal command to registry verification.

### Test Category: Success Workflows

**Test File**: `tests/e2e/success_workflows_test.go`

**Test Cases** (5 tests):

#### TC-E2E-001: Happy path - build, push, verify
```bash
# Test scenario
docker build -t myapp:v1 .
idpbuilder push myapp:v1 --registry https://gitea.test:8443 --username admin --password admin --insecure
# Verify in Gitea UI or API
curl -u admin:admin https://gitea.test:8443/v2/giteaadmin/myapp/manifests/v1
```
- **Given**: Fresh Docker image built locally
- **When**: Push command executed with valid flags
- **Then**: Exit code 0, image available in Gitea registry

#### TC-E2E-002: Custom registry override
```bash
idpbuilder push myapp:v1 --registry https://custom.registry.com
```
- **Given**: Custom registry URL specified
- **When**: Push executed
- **Then**: Image pushed to custom registry (not default Gitea)

#### TC-E2E-003: Environment variable configuration
```bash
export IDPBUILDER_REGISTRY_USERNAME=admin
export IDPBUILDER_REGISTRY_PASSWORD=admin
idpbuilder push myapp:v1 --insecure
```
- **Given**: Credentials in environment variables
- **When**: Push executed without flag credentials
- **Then**: Uses env vars, push succeeds

#### TC-E2E-004: Large image push with progress reporting
```bash
# Build 500MB image
docker build -t largeapp:v1 -f Dockerfile.large .
idpbuilder push largeapp:v1 --verbose --insecure
```
- **Given**: Large multi-layer image
- **When**: Push with verbose flag
- **Then**: Progress bars show layer upload, completes successfully

#### TC-E2E-005: Multiple images pushed sequentially
```bash
idpbuilder push app1:v1 --insecure
idpbuilder push app2:v1 --insecure
idpbuilder push app3:v1 --insecure
```
- **Given**: 3 different images
- **When**: Pushed one after another
- **Then**: All succeed, all available in registry

---

### Test Category: Failure Workflows

**Test File**: `tests/e2e/failure_workflows_test.go`

**Test Cases** (10 tests):

#### TC-E2E-006: Image not found in Docker daemon
```bash
idpbuilder push nonexistent:v1 --insecure
```
- **Expected**: Exit code 4, error message suggests `docker images`

#### TC-E2E-007: Invalid credentials
```bash
idpbuilder push myapp:v1 --username admin --password wrongpass --insecure
```
- **Expected**: Exit code 2, error message says authentication failed

#### TC-E2E-008: Registry unreachable
```bash
idpbuilder push myapp:v1 --registry https://unreachable.registry.com
```
- **Expected**: Exit code 3, error message about network failure

#### TC-E2E-009: Docker daemon not running
```bash
# Stop Docker daemon
sudo systemctl stop docker
idpbuilder push myapp:v1 --insecure
```
- **Expected**: Exit code 1, error message says Docker daemon unreachable

#### TC-E2E-010: Invalid image name format
```bash
idpbuilder push "Invalid Name:v1" --insecure
```
- **Expected**: Exit code 1, error explains OCI name requirements

#### TC-E2E-011: Missing required flags (no credentials)
```bash
idpbuilder push myapp:v1 --registry https://gitea.test:8443
```
- **Expected**: Exit code 1, error says credentials required

#### TC-E2E-012: TLS certificate verification failure (secure mode)
```bash
idpbuilder push myapp:v1 --registry https://gitea.test:8443
# Note: No --insecure flag, registry has self-signed cert
```
- **Expected**: Exit code 3, error about certificate verification, suggests --insecure flag

#### TC-E2E-013: Network interruption during push
```bash
# Simulate network failure mid-push (iptables DROP during test)
idpbuilder push myapp:v1 --insecure
```
- **Expected**: Exit code 3, error about network failure, no partial push

#### TC-E2E-014: Registry quota exceeded
```bash
# Configure test registry with quota
idpbuilder push large-image:v1 --insecure
```
- **Expected**: Exit code 1, error from registry about quota

#### TC-E2E-015: Conflicting flags
```bash
idpbuilder push myapp:v1 --registry https://gitea.test:8443 --registry https://other.registry.com
```
- **Expected**: Exit code 1, error about conflicting flags

---

### Test Category: Edge Cases

**Test File**: `tests/e2e/edge_cases_test.go`

**Test Cases** (5 tests):

#### TC-E2E-016: Image with single layer
#### TC-E2E-017: Image with 50+ layers
#### TC-E2E-018: Image name with maximum length (255 chars)
#### TC-E2E-019: Password with special characters
#### TC-E2E-020: Registry URL with non-standard port

---

## Quality Attribute Tests

### Purpose
Validate non-functional requirements (performance, security, reliability).

### Test Category: Performance Tests

**Test File**: `tests/quality/performance_test.go`

**Test Cases** (3 tests):

#### TC-PERF-001: Small image push performance
```go
func BenchmarkPush_SmallImage_10MB(b *testing.B)
```
- **Target**: < 5 seconds for 10MB image
- **Metric**: Time from command start to registry confirmation

#### TC-PERF-002: Medium image push performance
```go
func BenchmarkPush_MediumImage_100MB(b *testing.B)
```
- **Target**: < 30 seconds for 100MB image (network-dependent)
- **Metric**: Time and memory usage

#### TC-PERF-003: Memory efficiency during large image push
```go
func TestPush_MemoryEfficiency_500MB(t *testing.T)
```
- **Target**: Memory usage < 150MB regardless of image size
- **Validation**: Streaming architecture (no full buffering)

---

### Test Category: Security Tests

**Test File**: `tests/quality/security_test.go`

**Test Cases** (4 tests):

#### TC-SEC-001: Password never logged in plain text
```go
func TestSecurity_PasswordRedaction_Logs(t *testing.T)
```
- **Given**: Push command with --password flag
- **When**: Logs examined (stdout, stderr, debug logs)
- **Then**: Password appears as "***REDACTED***"

#### TC-SEC-002: Password not in process arguments
```go
func TestSecurity_PasswordRedaction_ProcessArgs(t *testing.T)
```
- **Given**: Push command running
- **When**: Process arguments inspected (`ps aux`)
- **Then**: Password not visible in process list

#### TC-SEC-003: Insecure mode warning displayed prominently
```go
func TestSecurity_InsecureWarning(t *testing.T)
```
- **Given**: Push with --insecure flag
- **When**: Command executes
- **Then**: Warning about TLS verification displayed

#### TC-SEC-004: Credentials not stored in files or cache
```go
func TestSecurity_NoCredentialPersistence(t *testing.T)
```
- **Given**: Push command completed
- **When**: Filesystem searched for credentials
- **Then**: No plain text credentials found

---

### Test Category: Reliability Tests

**Test File**: `tests/quality/reliability_test.go`

**Test Cases** (3 tests):

#### TC-REL-001: Transient network failures are retried
```go
func TestReliability_NetworkRetry(t *testing.T)
```
- **Given**: Network hiccup during layer upload
- **When**: Push in progress
- **Then**: Automatically retries, eventually succeeds

#### TC-REL-002: Graceful failure with cleanup
```go
func TestReliability_FailureCleanup(t *testing.T)
```
- **Given**: Push fails mid-operation
- **When**: Error encountered
- **Then**: No partial state, resources cleaned up

#### TC-REL-003: Concurrent operations are safe
```go
func TestReliability_Concurrency(t *testing.T)
```
- **Given**: Multiple goroutines using DockerClient
- **When**: Concurrent operations
- **Then**: No race conditions (verified with `-race` flag)

---

## Test-to-Phase Mapping

### Phase 1: Foundation & Interfaces

**Primary Test Coverage**:
- All Interface Contract Tests (TC-DOCKER-*, TC-REGISTRY-*, TC-AUTH-*, TC-TLS-*, TC-CMD-*)
- Total: ~70 interface tests

**Success Criteria**:
- All interface contract tests pass
- Mock implementations validate interface signatures
- Compilation succeeds with interface usage

**Phase 1 Wave 1 Focus**:
- Interface definition validation only (no implementation tests)

**Phase 1 Wave 2 Focus**:
- Implementation of interfaces (all contract tests must pass)

---

### Phase 2: Integration & Features

**Primary Test Coverage**:
- All Functional Integration Tests (TC-INT-*)
- Command Layer Tests (TC-CMD-*)
- Total: ~40 integration tests

**Success Criteria**:
- All integration tests pass
- Full stack integration working
- Error propagation validated
- Exit codes correct

**Phase 2 Wave 1 Focus** (Command Integration):
- TC-INT-021 to TC-INT-030 (full stack tests)

**Phase 2 Wave 2 Focus** (Advanced Features):
- Custom registry tests (TC-E2E-002)
- Environment variable tests (TC-E2E-003)

**Phase 2 Wave 3 Focus** (Error Handling):
- All failure workflow tests (TC-E2E-006 to TC-E2E-015)

---

### Phase 3: Testing & Production Readiness

**Primary Test Coverage**:
- All End-to-End Workflow Tests (TC-E2E-*)
- All Quality Attribute Tests (TC-PERF-*, TC-SEC-*, TC-REL-*)
- Total: ~30 E2E + quality tests

**Success Criteria**:
- 100% E2E tests pass
- Performance targets met
- Security validations pass
- Documentation reflects tested behavior

**Phase 3 Wave 1 Focus** (Integration Testing):
- TC-E2E-001 to TC-E2E-020 (all user workflows)

**Phase 3 Wave 2 Focus** (Documentation & Build):
- Verify documented examples match test scenarios
- CI/CD integration runs all tests

---

## Demo Scenarios

### Purpose
Human-demonstrable scenarios that validate user-facing value (R330/R291 integration).

### Demo 1: First-Time User Experience

**Scenario**: Platform engineer using idpbuilder push for the first time

**Steps**:
1. Build a simple web app container:
   ```bash
   cd examples/demo-app
   docker build -t demo-app:v1 .
   ```

2. Push to IDPBuilder Gitea registry:
   ```bash
   idpbuilder push demo-app:v1 \
     --registry https://gitea.cnoe.localtest.me:8443 \
     --username giteaadmin \
     --password password \
     --insecure \
     --verbose
   ```

3. Verify in Gitea UI:
   - Navigate to https://gitea.cnoe.localtest.me:8443
   - Click "Packages" → See "demo-app:v1"
   - Click on package → See layers and metadata

**Expected Outcome**:
- Command completes in < 10 seconds
- Progress bars show layer uploads
- Success message with registry URL
- Image visible and pullable from Gitea

**Success Criteria**:
- ✅ User can complete without reading documentation
- ✅ Error messages are self-explanatory
- ✅ Feedback is immediate and clear

---

### Demo 2: Troubleshooting Common Errors

**Scenario**: User encounters typical errors and resolves them using tool guidance

**Steps**:

1. **Error: Image Not Found**
   ```bash
   idpbuilder push nonexistent:v1 --insecure
   ```
   **Expected Output**:
   ```
   Error: Image 'nonexistent:v1' not found in Docker daemon
   Suggestion: Run 'docker images' to list available images
   ```

2. **Error: Wrong Credentials**
   ```bash
   idpbuilder push demo-app:v1 --username admin --password wrong --insecure
   ```
   **Expected Output**:
   ```
   Error: Authentication failed with registry
   Suggestion: Verify username and password are correct
   Registry: https://gitea.cnoe.localtest.me:8443
   ```

3. **Error: TLS Certificate Failure**
   ```bash
   idpbuilder push demo-app:v1  # Note: No --insecure flag
   ```
   **Expected Output**:
   ```
   Error: TLS certificate verification failed
   Suggestion: For local development with self-signed certificates, use --insecure flag
   Warning: Only use --insecure with trusted registries
   ```

**Expected Outcome**:
- User can self-resolve common issues
- Error messages provide actionable next steps

---

### Demo 3: Production Workflow Simulation

**Scenario**: CI/CD pipeline pushing release images

**Steps**:

1. Set credentials via environment:
   ```bash
   export IDPBUILDER_REGISTRY_USERNAME=ci-bot
   export IDPBUILDER_REGISTRY_PASSWORD=$CI_REGISTRY_TOKEN
   export IDPBUILDER_REGISTRY=https://registry.production.com
   ```

2. Build and tag release:
   ```bash
   docker build -t myapp:v2.1.0 .
   docker tag myapp:v2.1.0 myapp:latest
   ```

3. Push both tags:
   ```bash
   idpbuilder push myapp:v2.1.0
   idpbuilder push myapp:latest
   ```

4. Verify idempotency:
   ```bash
   # Push same image again (should be fast - layers exist)
   idpbuilder push myapp:v2.1.0
   ```

**Expected Outcome**:
- First push: Normal duration
- Second push: < 2 seconds (layers already exist)
- Exit code 0 for all pushes
- No credentials visible in logs

---

### Demo 4: Multi-Image Batch Push

**Scenario**: Developer pushing multiple microservices

**Steps**:

1. Build microservices:
   ```bash
   docker build -t api-service:v1 ./api
   docker build -t web-ui:v1 ./web
   docker build -t worker:v1 ./worker
   ```

2. Push all services:
   ```bash
   for service in api-service:v1 web-ui:v1 worker:v1; do
     idpbuilder push $service --insecure
   done
   ```

3. Verify all in registry:
   ```bash
   curl -u admin:admin https://gitea.test:8443/v2/_catalog
   ```

**Expected Outcome**:
- All images pushed successfully
- Progress clearly shows which image is current
- All images pullable from registry

---

### Demo 5: Failure Recovery Demonstration

**Scenario**: Network interruption during push, graceful recovery

**Steps**:

1. Start large image push:
   ```bash
   idpbuilder push large-app:v1 --verbose --insecure
   ```

2. Simulate network failure (test environment):
   - During push, temporarily block network
   - Tool should show retry attempts

3. Restore network:
   - Push should automatically resume and complete

**Expected Outcome**:
- Transient failures are automatically retried
- User sees retry messages
- Push completes successfully without restart

---

## Success Criteria

### Project Completion Criteria (R627 Compliance)

The project is considered complete when ALL tests pass and quality metrics are met:

#### Test Passage Requirements
- ✅ **100% Interface Contract Tests Pass** (~70 tests)
- ✅ **100% Functional Integration Tests Pass** (~40 tests)
- ✅ **100% E2E Workflow Tests Pass** (~30 tests)
- ✅ **100% Quality Attribute Tests Pass** (~10 tests)
- ✅ **Total: ~150 tests passing**

#### Code Coverage Requirements
- ✅ **Unit Test Coverage ≥85%** (all packages)
- ✅ **Integration Test Coverage ≥75%** (cross-package flows)
- ✅ **E2E Test Coverage: All critical paths** (defined in TC-E2E-001 to TC-E2E-020)

#### Quality Metrics
- ✅ **Zero bugs at completion** (all found bugs fixed, R627)
- ✅ **Performance targets met**:
  - Small image (10MB): < 5 seconds
  - Medium image (100MB): < 30 seconds
  - Memory usage: < 150MB during any push
- ✅ **Security validations pass**:
  - No password leaks in logs or process args
  - TLS warnings displayed in insecure mode
- ✅ **All efforts ≤800 lines** (R220/R221)

#### Documentation & Integration
- ✅ **All demo scenarios execute successfully**
- ✅ **README updated with push command examples**
- ✅ **Integration tests run in CI/CD**
- ✅ **Build system integration complete** (`make build`, `make test`)

#### Architect Approval
- ✅ **Architect review: Grade A+** (100/100 required)
- ✅ **Phase reviews approved** (Phase 1, 2, 3)
- ✅ **Interface contracts validated** (Phase 1 Wave 1)
- ✅ **Integration coherence verified** (Phase 2/3)

---

## Test Infrastructure Files

### Test Harness Location

**Primary Test Harness**: `/PROJECT-TEST-HARNESS.sh` (project root)

**Purpose**: Execute all project-level tests in correct order

**Usage**:
```bash
# Run all tests
./PROJECT-TEST-HARNESS.sh

# Run specific category
./PROJECT-TEST-HARNESS.sh --category=interfaces
./PROJECT-TEST-HARNESS.sh --category=integration
./PROJECT-TEST-HARNESS.sh --category=e2e
./PROJECT-TEST-HARNESS.sh --category=quality

# Run with coverage
./PROJECT-TEST-HARNESS.sh --coverage

# Run in CI mode (strict)
./PROJECT-TEST-HARNESS.sh --ci
```

### Test Environment Setup

**Docker Compose Test Registry**: `tests/e2e/docker-compose.yml`

**Purpose**: Spin up local Gitea registry for E2E tests

**Services**:
- Gitea with OCI registry enabled
- Self-signed TLS certificate
- Default admin credentials
- Pre-configured with test organization

**Usage**:
```bash
cd tests/e2e
docker-compose up -d
# Run E2E tests
docker-compose down
```

### Test Data Fixtures

**Location**: `tests/fixtures/`

**Contents**:
- **small-image/**: Dockerfile for 10MB test image
- **medium-image/**: Dockerfile for 100MB test image
- **large-image/**: Dockerfile for 500MB test image
- **multi-layer-image/**: Dockerfile with 20+ layers
- **invalid-images/**: Test cases for error scenarios

---

## Test Maintenance Strategy

### Test Review Cadence

**During Implementation**:
- Review test failures daily
- Update tests if requirements change
- Add tests for newly discovered edge cases

**After Phase Completion**:
- Full test suite run
- Coverage report review
- Performance benchmark comparison

### Test Evolution

**When to Add Tests**:
- Bug found → Add regression test
- New requirement → Add acceptance test
- Performance issue → Add performance test

**When to Update Tests**:
- API contract changes → Update interface tests
- Feature enhancement → Expand test scenarios
- Optimization → Update performance targets

**When to Remove Tests**:
- Feature deprecated → Archive related tests
- Test duplicates coverage → Consolidate

---

## R340 Compliance: Planning File Metadata

**Test Plan Location**: `/home/vscode/workspaces/idpbuilder-oci-push-rebuild/planning/project/PROJECT-TEST-PLAN.md`

**Test Harness Location**: `/home/vscode/workspaces/idpbuilder-oci-push-rebuild/PROJECT-TEST-HARNESS.sh` (to be created)

**Test Files Location**: `tests/` directory structure (to be created during implementation)

**Orchestrator Tracking**: Orchestrator will record these locations in `orchestrator-state-v3.json` under:
```json
{
  "test_plans": {
    "project": {
      "plan_file": "planning/project/PROJECT-TEST-PLAN.md",
      "harness_file": "PROJECT-TEST-HARNESS.sh",
      "test_directory": "tests/",
      "created_at": "2025-11-11T01:56:43Z",
      "created_by": "code-reviewer",
      "status": "pending_implementation"
    }
  }
}
```

---

## TDD Workflow Summary

**Current State**: ✅ RED PHASE COMPLETE
- All test specifications defined
- Expected test count: ~150 tests
- All tests WILL FAIL initially (expected behavior)

**Next State**: ⏳ GREEN PHASE (during implementation)
- Phase 1: Make interface tests pass
- Phase 2: Make integration tests pass
- Phase 3: Make E2E and quality tests pass

**Final State**: ⏳ REFACTOR PHASE (post-implementation)
- All tests passing
- Code optimized while maintaining passage
- Documentation matches tested behavior

---

**Test Plan Status**: COMPLETE - Ready for Implementation Planning
**TDD Phase**: RED (tests defined, expecting failure)
**Expected Next Step**: Create PROJECT-IMPLEMENTATION-PLAN.md using this test plan as specification

---

**R341 TDD Compliance**: ✅ VERIFIED
- Tests created BEFORE implementation planning
- Tests define success criteria
- Tests are executable specifications
- Implementation will target test passage
