# Wave 2 Test Plan
## Phase 1, Wave 2: Core Package Implementations - Test Strategy with Progressive Realism

**Wave**: Wave 2 - Core Package Implementations
**Phase**: Phase 1 - Foundation & Interfaces
**Created**: 2025-10-29 05:54:00 UTC
**Test Planner**: @agent-code-reviewer
**R341 Compliance**: TDD (Tests Before Implementation)
**Progressive Realism**: Uses ACTUAL Wave 1 code as test infrastructure

---

## Document Metadata

| Field | Value |
|-------|-------|
| **Wave Name** | Core Package Implementations |
| **Test Approach** | Progressive Realism + Unit Testing |
| **Test Philosophy** | Tests reference REAL Wave 1 interfaces (not imaginary) |
| **Test Framework** | Go testing + testify/assert + go-containerregistry |
| **Coverage Target** | 85%+ (90%+ for auth/tls security-critical) |
| **Status** | Ready for Wave 2 Implementation |

---

## 1. Progressive Realism Explained

### 1.1 What is Progressive Realism?

**Definition**: Test planning that uses ACTUAL code from completed waves as fixtures, mocks, and test infrastructure - NOT abstract/imaginary structures.

**Contrast**:

**❌ BAD (Abstract Test Planning)**:
```go
// Imaginary structure - doesn't exist in codebase
type MockDockerClient struct {
    SomeField string
}
```

**✅ GOOD (Progressive Realism)**:
```go
// REAL imports from Wave 1 actual code
import (
    "github.com/cnoe-io/idpbuilder/pkg/docker"  // Actual Wave 1 interface
)

// Mock using ACTUAL Wave 1 interface type discovered from real code
type mockDockerClient struct {
    docker.Client  // Real interface from Wave 1 effort-1-docker-interface
    imageExistsFunc func(ctx context.Context, imageName string) (bool, error)
}

func (m *mockDockerClient) ImageExists(ctx context.Context, imageName string) (bool, error) {
    if m.imageExistsFunc != nil {
        return m.imageExistsFunc(ctx, imageName)
    }
    return false, nil
}
```

### 1.2 Wave 1 Actual Code Discovery

**From Actual Wave 1 Efforts:**

**Effort 1 - Docker Interface** (`efforts/phase1/wave1/effort-1-docker-interface/`):
- Package: `github.com/cnoe-io/idpbuilder/pkg/docker`
- Interface: `docker.Client` (4 methods)
- Error Types: `DaemonConnectionError`, `ImageNotFoundError`, `ImageConversionError`, `ValidationError`
- Constructor: `NewClient() (Client, error)`

**Effort 2 - Registry Interface** (`efforts/phase1/wave1/effort-2-registry-interface/`):
- Package: `github.com/cnoe-io/idpbuilder/pkg/registry`
- Interface: `registry.Client` (3 methods)
- Types: `ProgressCallback`, `ProgressUpdate`
- Error Types: `AuthenticationError`, `NetworkError`, `RegistryUnavailableError`, `PushFailedError`
- Constructor: `NewClient(authProvider AuthProvider, tlsConfig TLSConfigProvider) (Client, error)`

**Effort 3 - Auth Interface** (`efforts/phase1/wave1/effort-3-auth-tls-interfaces/`):
- Package: `github.com/cnoe-io/idpbuilder/pkg/auth`
- Interface: `auth.Provider` (2 methods)
- Types: `Credentials` struct
- Error Types: `CredentialValidationError`
- Constructor: `NewBasicAuthProvider(username, password string) Provider`

**Effort 4 - TLS Interface** (`efforts/phase1/wave1/effort-3-auth-tls-interfaces/`):
- Package: `github.com/cnoe-io/idpbuilder/pkg/tls`
- Interface: `tls.ConfigProvider` (2 methods)
- Types: `Config` struct
- Constructor: `NewConfigProvider(insecure bool) ConfigProvider`

### 1.3 R341 TDD Compliance for Wave 2

**Critical Understanding:** Wave 2 tests are written BEFORE Wave 2 implementations.

**Wave 2 TDD Workflow:**
```bash
# Step 1: Read Wave 2 architecture (WHAT to build)
# Step 2: Read Wave 1 actual interfaces (TEST INFRASTRUCTURE)
# Step 3: Write comprehensive unit tests referencing Wave 1 interfaces
# Step 4: Implement Wave 2 code to pass those tests
# Step 5: Verify 85%+ coverage achieved
```

**Why This Works:**
- Wave 1 interfaces are FROZEN (can't change)
- Wave 2 tests use Wave 1 interfaces as contracts
- Implementations must satisfy Wave 1 contracts
- Tests validate correct interface usage

---

## 2. Test Infrastructure from Wave 1

### 2.1 Actual Wave 1 Packages Available

**Package: docker** (from `pkg/docker/interface.go`):
```go
type Client interface {
    ImageExists(ctx context.Context, imageName string) (bool, error)
    GetImage(ctx context.Context, imageName string) (v1.Image, error)
    ValidateImageName(imageName string) error
    Close() error
}
```

**Package: registry** (from `pkg/registry/interface.go`):
```go
type Client interface {
    Push(ctx context.Context, image v1.Image, targetRef string,
         progressCallback ProgressCallback) error
    BuildImageReference(registryURL, imageName string) (string, error)
    ValidateRegistry(ctx context.Context, registryURL string) error
}

type ProgressCallback func(update ProgressUpdate)

type ProgressUpdate struct {
    LayerDigest string
    LayerSize   int64
    BytesPushed int64
    Status      string  // "uploading", "complete", "exists"
}
```

**Package: auth** (from `pkg/auth/interface.go`):
```go
type Provider interface {
    GetAuthenticator() (authn.Authenticator, error)
    ValidateCredentials() error
}

type Credentials struct {
    Username string
    Password string
}
```

**Package: tls** (from `pkg/tls/interface.go`):
```go
type ConfigProvider interface {
    GetTLSConfig() *tls.Config
    IsInsecure() bool
}

type Config struct {
    InsecureSkipVerify bool
}
```

### 2.2 Actual Wave 1 Error Types

**Docker Errors** (from `pkg/docker/errors.go`):
```go
type DaemonConnectionError struct { Cause error }
type ImageNotFoundError struct { ImageName string }
type ImageConversionError struct { ImageName string; Cause error }
type ValidationError struct { Field string; Message string }
```

**Registry Errors** (from `pkg/registry/errors.go`):
```go
type AuthenticationError struct { Registry string; Cause error }
type NetworkError struct { Registry string; Cause error }
type RegistryUnavailableError struct { Registry string; StatusCode int }
type PushFailedError struct { TargetRef string; Cause error }
```

**Auth Errors** (from `pkg/auth/errors.go`):
```go
type CredentialValidationError struct { Field string; Reason string }
```

---

## 3. Coverage Targets

### 3.1 Per-Package Coverage Requirements

| Package | Minimum Coverage | Rationale |
|---------|-----------------|-----------|
| **pkg/docker** | 85% | Business logic + external API integration |
| **pkg/registry** | 85% | Complex push logic + error classification |
| **pkg/auth** | 90% | Security-critical credential handling |
| **pkg/tls** | 90% | Security-critical certificate verification |

**Why Different Targets?**
- Auth/TLS are security-critical (higher bar)
- Docker/Registry have external dependencies (some paths hard to test)
- Unit tests with mocks for external services

### 3.2 What Coverage Measures

**Included in Coverage:**
- All business logic functions
- Error handling paths
- Input validation
- Type conversions
- Helper functions

**Excluded from Coverage (per R007):**
- Test files (`*_test.go`)
- Generated code (none expected in Wave 2)
- Demo/example files (none expected)
- Documentation files (`.md`)

---

## 4. Test Cases by Package

### 4.1 Package: docker (Effort 1.2.1)

**Files to Test:**
- `pkg/docker/client.go` (Wave 2 implementation)

**Test File:**
- `pkg/docker/client_test.go`

#### Test Categories

**A. Constructor Tests**

##### TC-DOCKER-IMPL-001: NewClient Success

**Purpose**: Verify Docker client creation succeeds when daemon available

**Test Strategy**:
```go
package docker

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewClient_Success(t *testing.T) {
    // Given: Docker daemon is running (prerequisite)
    // When: Creating new client
    client, err := NewClient()

    // Then: Client created successfully
    require.NoError(t, err, "NewClient should succeed with running daemon")
    require.NotNil(t, client, "Client should not be nil")

    // Cleanup
    defer client.Close()

    // Verify client satisfies interface (compile-time check already done)
    var _ Client = client
}
```

**Coverage**: Constructor, daemon ping, connection initialization

---

##### TC-DOCKER-IMPL-002: NewClient DaemonNotRunning

**Purpose**: Verify proper error when daemon not running

**Test Strategy**:
```go
func TestNewClient_DaemonNotRunning(t *testing.T) {
    // Note: This test requires Docker daemon to be stopped
    // May be skipped in CI if daemon always available

    // Given: Docker daemon is NOT running
    // (Test setup would stop daemon or use mock)

    // When: Creating new client
    client, err := NewClient()

    // Then: Returns DaemonConnectionError
    require.Error(t, err, "NewClient should fail without daemon")
    assert.Nil(t, client, "Client should be nil on error")

    // Verify error type is correct (from Wave 1)
    var connErr *DaemonConnectionError
    assert.ErrorAs(t, err, &connErr, "Error should be DaemonConnectionError")
}
```

**Coverage**: Error path, daemon unavailable, error type wrapping

---

**B. ImageExists Tests**

##### TC-DOCKER-IMPL-003: ImageExists_ImagePresent

**Purpose**: Verify ImageExists returns true for present images

**Test Strategy**:
```go
func TestImageExists_ImagePresent(t *testing.T) {
    // Given: Client connected to daemon
    client, err := NewClient()
    require.NoError(t, err)
    defer client.Close()

    ctx := context.Background()

    // Given: Image exists in daemon (test uses common base image)
    imageName := "alpine:latest"  // Assume alpine pulled as test prerequisite

    // When: Checking if image exists
    exists, err := client.ImageExists(ctx, imageName)

    // Then: Returns true with no error
    require.NoError(t, err, "ImageExists should not error for valid check")
    assert.True(t, exists, "alpine:latest should exist (test prerequisite)")
}
```

**Coverage**: ImageExists happy path, Docker API call, result parsing

---

##### TC-DOCKER-IMPL-004: ImageExists_ImageNotPresent

**Purpose**: Verify ImageExists returns false (not error) for missing images

**Test Strategy**:
```go
func TestImageExists_ImageNotPresent(t *testing.T) {
    // Given: Client connected
    client, err := NewClient()
    require.NoError(t, err)
    defer client.Close()

    ctx := context.Background()

    // Given: Image definitely does NOT exist
    imageName := "nonexistent-test-image-12345:fake-tag"

    // When: Checking if image exists
    exists, err := client.ImageExists(ctx, imageName)

    // Then: Returns false (NOT an error - image just not found)
    require.NoError(t, err, "ImageExists should NOT error for missing image")
    assert.False(t, exists, "Non-existent image should return false")
}
```

**Coverage**: ImageExists with 404 NotFound handling, error classification

---

##### TC-DOCKER-IMPL-005: ImageExists_InvalidImageName

**Purpose**: Verify validation error for malformed image names

**Test Strategy**:
```go
func TestImageExists_InvalidImageName(t *testing.T) {
    // Given: Client connected
    client, err := NewClient()
    require.NoError(t, err)
    defer client.Close()

    ctx := context.Background()

    // When: Checking image with invalid name (empty string)
    exists, err := client.ImageExists(ctx, "")

    // Then: Returns ValidationError
    require.Error(t, err, "ImageExists should error for empty name")
    assert.False(t, exists, "Should return false on validation error")

    // Verify error type from Wave 1
    var valErr *ValidationError
    assert.ErrorAs(t, err, &valErr, "Should be ValidationError")
    assert.Equal(t, "imageName", valErr.Field)
}
```

**Coverage**: Input validation, ValidationError usage, early validation

---

**C. GetImage Tests**

##### TC-DOCKER-IMPL-006: GetImage_Success

**Purpose**: Verify GetImage retrieves and converts image to OCI format

**Test Strategy**:
```go
func TestGetImage_Success(t *testing.T) {
    // Given: Client and existing image
    client, err := NewClient()
    require.NoError(t, err)
    defer client.Close()

    ctx := context.Background()
    imageName := "alpine:latest"

    // When: Getting image
    image, err := client.GetImage(ctx, imageName)

    // Then: Returns v1.Image successfully
    require.NoError(t, err, "GetImage should succeed for existing image")
    require.NotNil(t, image, "Image should not be nil")

    // Verify image is valid OCI v1.Image (from go-containerregistry)
    // This validates conversion worked
    manifest, err := image.Manifest()
    require.NoError(t, err, "Should be able to get manifest from image")
    assert.NotEmpty(t, manifest.Layers, "Image should have layers")
}
```

**Coverage**: GetImage happy path, daemon.Image() call, OCI conversion

---

##### TC-DOCKER-IMPL-007: GetImage_ImageNotFound

**Purpose**: Verify ImageNotFoundError for non-existent images

**Test Strategy**:
```go
func TestGetImage_ImageNotFound(t *testing.T) {
    // Given: Client connected
    client, err := NewClient()
    require.NoError(t, err)
    defer client.Close()

    ctx := context.Background()
    imageName := "nonexistent-image-99999:missing"

    // When: Getting non-existent image
    image, err := client.GetImage(ctx, imageName)

    // Then: Returns ImageNotFoundError (Wave 1 type)
    require.Error(t, err, "GetImage should error for missing image")
    assert.Nil(t, image, "Image should be nil on error")

    var notFoundErr *ImageNotFoundError
    assert.ErrorAs(t, err, &notFoundErr, "Should be ImageNotFoundError")
    assert.Equal(t, imageName, notFoundErr.ImageName)
}
```

**Coverage**: Error path, ImageNotFoundError usage, exists check in GetImage

---

##### TC-DOCKER-IMPL-008: GetImage_ConversionError

**Purpose**: Verify ImageConversionError when conversion fails

**Test Strategy**:
```go
func TestGetImage_ConversionError(t *testing.T) {
    // Given: Client connected
    client, err := NewClient()
    require.NoError(t, err)
    defer client.Close()

    ctx := context.Background()

    // This test is harder - requires corrupted image or mock
    // May be tested via integration test or skipped if hard to reproduce

    // When: Getting image with bad format
    // (Test setup would need to create corrupted image scenario)

    // Then: Returns ImageConversionError
    // var convErr *ImageConversionError
    // assert.ErrorAs(t, err, &convErr)

    t.Skip("ImageConversionError requires special test setup - covered in integration")
}
```

**Coverage**: Conversion error path (may require mocking or integration test)

---

**D. ValidateImageName Tests**

##### TC-DOCKER-IMPL-009: ValidateImageName_Valid

**Purpose**: Verify validation passes for valid image names

**Test Strategy**:
```go
func TestValidateImageName_Valid(t *testing.T) {
    // Given: Client
    client, err := NewClient()
    require.NoError(t, err)
    defer client.Close()

    // Test cases with valid names
    validNames := []string{
        "myapp:latest",
        "registry.io/namespace/repo:tag",
        "myapp",  // no tag
        "my-app:v1.0.0",
        "my_app:latest",
        "registry.io:5000/repo:tag",
    }

    for _, name := range validNames {
        t.Run(name, func(t *testing.T) {
            // When: Validating name
            err := client.ValidateImageName(name)

            // Then: No error
            assert.NoError(t, err, "Valid name should pass: %s", name)
        })
    }
}
```

**Coverage**: Validation logic, multiple valid formats, pattern matching

---

##### TC-DOCKER-IMPL-010: ValidateImageName_Invalid

**Purpose**: Verify validation rejects invalid/dangerous names

**Test Strategy**:
```go
func TestValidateImageName_Invalid(t *testing.T) {
    // Given: Client
    client, err := NewClient()
    require.NoError(t, err)
    defer client.Close()

    // Test cases with invalid names
    testCases := []struct {
        name        string
        expectedMsg string
    }{
        {"", "empty"},
        {"myapp;rm -rf /", "semicolon"},
        {"myapp|ls", "pipe"},
        {"myapp && ls", "command injection"},
        {"myapp`whoami`", "backtick"},
        {"myapp$(whoami)", "command substitution"},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // When: Validating invalid name
            err := client.ValidateImageName(tc.name)

            // Then: Returns ValidationError
            require.Error(t, err, "Invalid name should fail: %s", tc.name)

            var valErr *ValidationError
            assert.ErrorAs(t, err, &valErr, "Should be ValidationError")
            assert.Equal(t, "imageName", valErr.Field)
        })
    }
}
```

**Coverage**: Security validation, command injection prevention, edge cases

---

**E. Close Tests**

##### TC-DOCKER-IMPL-011: Close_Success

**Purpose**: Verify client cleanup works

**Test Strategy**:
```go
func TestClose_Success(t *testing.T) {
    // Given: Client created
    client, err := NewClient()
    require.NoError(t, err)

    // When: Closing client
    err = client.Close()

    // Then: No error
    assert.NoError(t, err, "Close should succeed")

    // Note: After close, client should not be used
}
```

**Coverage**: Cleanup path, resource release

---

**F. Edge Case Tests**

##### TC-DOCKER-IMPL-012: Context_Cancellation

**Purpose**: Verify context cancellation is respected

**Test Strategy**:
```go
func TestImageExists_ContextCancellation(t *testing.T) {
    // Given: Client and cancelled context
    client, err := NewClient()
    require.NoError(t, err)
    defer client.Close()

    ctx, cancel := context.WithCancel(context.Background())
    cancel()  // Cancel immediately

    // When: Calling ImageExists with cancelled context
    exists, err := client.ImageExists(ctx, "alpine:latest")

    // Then: Returns error (context cancelled)
    require.Error(t, err, "Should error with cancelled context")
    assert.False(t, exists)
    assert.Contains(t, err.Error(), "context")
}
```

**Coverage**: Context handling, timeout/cancellation support

---

**Docker Package Test Summary:**
- Total Tests: ~12-15 test cases
- Coverage Target: 85%
- Focus: API integration, error handling, validation, security

---

### 4.2 Package: registry (Effort 1.2.2)

**Files to Test:**
- `pkg/registry/client.go` (Wave 2 implementation)

**Test File:**
- `pkg/registry/client_test.go`

#### Test Categories

**A. Constructor Tests**

##### TC-REGISTRY-IMPL-001: NewClient_Success

**Purpose**: Verify registry client creation with valid providers

**Test Strategy**:
```go
package registry

import (
    "context"
    "crypto/tls"
    "testing"

    "github.com/cnoe-io/idpbuilder/pkg/auth"
    tlspkg "github.com/cnoe-io/idpbuilder/pkg/tls"
    "github.com/google/go-containerregistry/pkg/authn"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// Mock auth provider using Wave 1 interface
type mockAuthProvider struct {
    auth.Provider
    username string
    password string
}

func (m *mockAuthProvider) GetAuthenticator() (authn.Authenticator, error) {
    return &authn.Basic{
        Username: m.username,
        Password: m.password,
    }, nil
}

func (m *mockAuthProvider) ValidateCredentials() error {
    if m.username == "" {
        return &auth.CredentialValidationError{
            Field:  "username",
            Reason: "username empty",
        }
    }
    return nil
}

// Mock TLS provider using Wave 1 interface
type mockTLSProvider struct {
    tlspkg.ConfigProvider
    insecure bool
}

func (m *mockTLSProvider) GetTLSConfig() *tls.Config {
    return &tls.Config{
        InsecureSkipVerify: m.insecure,
    }
}

func (m *mockTLSProvider) IsInsecure() bool {
    return m.insecure
}

func TestNewClient_Success(t *testing.T) {
    // Given: Valid auth and TLS providers
    authProvider := &mockAuthProvider{
        username: "testuser",
        password: "testpass",
    }
    tlsProvider := &mockTLSProvider{
        insecure: true,
    }

    // When: Creating registry client
    client, err := NewClient(authProvider, tlsProvider)

    // Then: Client created successfully
    require.NoError(t, err, "NewClient should succeed with valid providers")
    require.NotNil(t, client, "Client should not be nil")

    // Verify client satisfies interface
    var _ Client = client
}
```

**Coverage**: Constructor, provider validation, HTTP client setup

---

##### TC-REGISTRY-IMPL-002: NewClient_NilAuthProvider

**Purpose**: Verify error when auth provider is nil

**Test Strategy**:
```go
func TestNewClient_NilAuthProvider(t *testing.T) {
    // Given: Nil auth provider
    tlsProvider := &mockTLSProvider{insecure: true}

    // When: Creating client
    client, err := NewClient(nil, tlsProvider)

    // Then: Returns ValidationError
    require.Error(t, err, "Should error with nil auth provider")
    assert.Nil(t, client, "Client should be nil on error")

    var valErr *ValidationError
    assert.ErrorAs(t, err, &valErr)
    assert.Equal(t, "authProvider", valErr.Field)
}
```

**Coverage**: Nil validation, defensive programming

---

##### TC-REGISTRY-IMPL-003: NewClient_NilTLSProvider

**Purpose**: Verify error when TLS provider is nil

**Test Strategy**:
```go
func TestNewClient_NilTLSProvider(t *testing.T) {
    // Given: Nil TLS provider
    authProvider := &mockAuthProvider{
        username: "testuser",
        password: "testpass",
    }

    // When: Creating client
    client, err := NewClient(authProvider, nil)

    // Then: Returns ValidationError
    require.Error(t, err, "Should error with nil TLS provider")
    assert.Nil(t, client, "Client should be nil on error")

    var valErr *ValidationError
    assert.ErrorAs(t, err, &valErr)
    assert.Equal(t, "tlsConfig", valErr.Field)
}
```

**Coverage**: Nil validation, error construction

---

**B. Push Tests**

##### TC-REGISTRY-IMPL-004: Push_Success

**Purpose**: Verify successful image push to registry

**Test Strategy**:
```go
func TestPush_Success(t *testing.T) {
    // Note: This test requires a test registry or extensive mocking
    // May use go-containerregistry's registry/test package

    // Given: Client with valid providers
    client := setupTestClient(t)

    // Given: Test image (from test helper)
    image := createTestImage(t)

    ctx := context.Background()
    targetRef := "localhost:5000/test/myapp:latest"

    // Track progress callbacks
    progressUpdates := []ProgressUpdate{}
    callback := func(update ProgressUpdate) {
        progressUpdates = append(progressUpdates, update)
    }

    // When: Pushing image
    err := client.Push(ctx, image, targetRef, callback)

    // Then: Push succeeds
    require.NoError(t, err, "Push should succeed")

    // Verify progress callbacks were invoked
    assert.NotEmpty(t, progressUpdates, "Progress callback should be invoked")
}
```

**Coverage**: Push happy path, progress callbacks, remote.Write() integration

---

##### TC-REGISTRY-IMPL-005: Push_AuthenticationError

**Purpose**: Verify authentication error handling (401/403)

**Test Strategy**:
```go
func TestPush_AuthenticationError(t *testing.T) {
    // Given: Client with invalid credentials
    authProvider := &mockAuthProvider{
        username: "invalid",
        password: "wrong",
    }
    tlsProvider := &mockTLSProvider{insecure: true}
    client, _ := NewClient(authProvider, tlsProvider)

    // Given: Test image
    image := createTestImage(t)
    ctx := context.Background()
    targetRef := "localhost:5000/test/myapp:latest"

    // When: Pushing with bad credentials
    err := client.Push(ctx, image, targetRef, nil)

    // Then: Returns AuthenticationError (Wave 1 type)
    require.Error(t, err, "Push should fail with bad credentials")

    var authErr *AuthenticationError
    assert.ErrorAs(t, err, &authErr, "Should be AuthenticationError")
    assert.Contains(t, authErr.Registry, targetRef)
}
```

**Coverage**: Auth error path, error classification, 401/403 handling

---

##### TC-REGISTRY-IMPL-006: Push_NetworkError

**Purpose**: Verify network error handling

**Test Strategy**:
```go
func TestPush_NetworkError(t *testing.T) {
    // Given: Client pointing to unreachable registry
    client := setupTestClient(t)
    image := createTestImage(t)

    ctx := context.Background()
    targetRef := "unreachable.registry.io:9999/test/myapp:latest"

    // When: Pushing to unreachable registry
    err := client.Push(ctx, image, targetRef, nil)

    // Then: Returns NetworkError (Wave 1 type)
    require.Error(t, err, "Push should fail for unreachable registry")

    var netErr *NetworkError
    assert.ErrorAs(t, err, &netErr, "Should be NetworkError")
}
```

**Coverage**: Network error classification, timeout handling

---

##### TC-REGISTRY-IMPL-007: Push_InvalidTargetRef

**Purpose**: Verify error for invalid target reference

**Test Strategy**:
```go
func TestPush_InvalidTargetRef(t *testing.T) {
    // Given: Client
    client := setupTestClient(t)
    image := createTestImage(t)
    ctx := context.Background()

    // When: Pushing with invalid reference format
    invalidRef := "not a valid reference!@#$"
    err := client.Push(ctx, image, invalidRef, nil)

    // Then: Returns PushFailedError
    require.Error(t, err, "Push should fail for invalid reference")

    var pushErr *PushFailedError
    assert.ErrorAs(t, err, &pushErr)
    assert.Equal(t, invalidRef, pushErr.TargetRef)
}
```

**Coverage**: Input validation, name.ParseReference() error handling

---

**C. BuildImageReference Tests**

##### TC-REGISTRY-IMPL-008: BuildImageReference_Success

**Purpose**: Verify correct reference construction

**Test Strategy**:
```go
func TestBuildImageReference_Success(t *testing.T) {
    // Given: Client
    client := setupTestClient(t)

    testCases := []struct {
        registryURL string
        imageName   string
        expected    string
    }{
        {
            registryURL: "https://gitea.cnoe.localtest.me:8443",
            imageName:   "myapp:latest",
            expected:    "gitea.cnoe.localtest.me:8443/giteaadmin/myapp:latest",
        },
        {
            registryURL: "https://registry.io",
            imageName:   "myapp",  // no tag
            expected:    "registry.io/giteaadmin/myapp:latest",  // default tag
        },
        {
            registryURL: "https://localhost:5000",
            imageName:   "test-app:v1.0.0",
            expected:    "localhost:5000/giteaadmin/test-app:v1.0.0",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.imageName, func(t *testing.T) {
            // When: Building reference
            ref, err := client.BuildImageReference(tc.registryURL, tc.imageName)

            // Then: Correct reference constructed
            require.NoError(t, err)
            assert.Equal(t, tc.expected, ref)
        })
    }
}
```

**Coverage**: Reference construction, namespace injection, tag defaulting

---

##### TC-REGISTRY-IMPL-009: BuildImageReference_InvalidURL

**Purpose**: Verify error for invalid registry URL

**Test Strategy**:
```go
func TestBuildImageReference_InvalidURL(t *testing.T) {
    // Given: Client
    client := setupTestClient(t)

    invalidURLs := []string{
        "not a url",
        "://missing-scheme",
        "",
    }

    for _, url := range invalidURLs {
        t.Run(url, func(t *testing.T) {
            // When: Building reference with invalid URL
            ref, err := client.BuildImageReference(url, "myapp:latest")

            // Then: Returns ValidationError
            require.Error(t, err)
            assert.Empty(t, ref)

            var valErr *ValidationError
            assert.ErrorAs(t, err, &valErr)
            assert.Equal(t, "registryURL", valErr.Field)
        })
    }
}
```

**Coverage**: URL validation, url.Parse() error handling

---

**D. ValidateRegistry Tests**

##### TC-REGISTRY-IMPL-010: ValidateRegistry_Success

**Purpose**: Verify registry /v2/ endpoint validation

**Test Strategy**:
```go
func TestValidateRegistry_Success(t *testing.T) {
    // Given: Client and reachable registry
    client := setupTestClient(t)
    ctx := context.Background()

    // Requires test registry or mock
    registryURL := "https://localhost:5000"

    // When: Validating registry
    err := client.ValidateRegistry(ctx, registryURL)

    // Then: No error (registry reachable)
    assert.NoError(t, err, "Should succeed for reachable registry")
}
```

**Coverage**: /v2/ endpoint check, 200/401 handling

---

##### TC-REGISTRY-IMPL-011: ValidateRegistry_Unreachable

**Purpose**: Verify error for unreachable registry

**Test Strategy**:
```go
func TestValidateRegistry_Unreachable(t *testing.T) {
    // Given: Client
    client := setupTestClient(t)
    ctx := context.Background()

    // When: Validating unreachable registry
    registryURL := "https://unreachable.registry.test:9999"
    err := client.ValidateRegistry(ctx, registryURL)

    // Then: Returns NetworkError
    require.Error(t, err)

    var netErr *NetworkError
    assert.ErrorAs(t, err, &netErr)
    assert.Equal(t, registryURL, netErr.Registry)
}
```

**Coverage**: Network error detection, timeout handling

---

**Registry Package Test Summary:**
- Total Tests: ~15-18 test cases
- Coverage Target: 85%
- Focus: Push operations, error classification, reference construction

---

### 4.3 Package: auth (Effort 1.2.3)

**Files to Test:**
- `pkg/auth/basic.go` (Wave 2 implementation)

**Test File:**
- `pkg/auth/basic_test.go`

#### Test Categories

**A. Constructor Tests**

##### TC-AUTH-IMPL-001: NewBasicAuthProvider_Success

**Purpose**: Verify basic auth provider creation

**Test Strategy**:
```go
package auth

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewBasicAuthProvider_Success(t *testing.T) {
    // Given: Valid credentials
    username := "testuser"
    password := "testP@ss123!"

    // When: Creating provider
    provider := NewBasicAuthProvider(username, password)

    // Then: Provider created
    require.NotNil(t, provider)

    // Verify satisfies interface
    var _ Provider = provider
}
```

**Coverage**: Constructor, credential storage

---

**B. GetAuthenticator Tests**

##### TC-AUTH-IMPL-002: GetAuthenticator_Success

**Purpose**: Verify authenticator creation for valid credentials

**Test Strategy**:
```go
func TestGetAuthenticator_Success(t *testing.T) {
    // Given: Provider with valid credentials
    provider := NewBasicAuthProvider("testuser", "testpass")

    // When: Getting authenticator
    authenticator, err := provider.GetAuthenticator()

    // Then: Returns valid authn.Authenticator
    require.NoError(t, err)
    require.NotNil(t, authenticator)

    // Verify it's authn.Basic type (from go-containerregistry)
    config, err := authenticator.Authorization()
    require.NoError(t, err)
    assert.NotEmpty(t, config.Username)
    assert.NotEmpty(t, config.Password)
}
```

**Coverage**: Authenticator creation, authn.Basic conversion

---

##### TC-AUTH-IMPL-003: GetAuthenticator_EmptyUsername

**Purpose**: Verify error for empty username

**Test Strategy**:
```go
func TestGetAuthenticator_EmptyUsername(t *testing.T) {
    // Given: Provider with empty username
    provider := NewBasicAuthProvider("", "password")

    // When: Getting authenticator
    authenticator, err := provider.GetAuthenticator()

    // Then: Returns error
    require.Error(t, err)
    assert.Nil(t, authenticator)

    var credErr *CredentialValidationError
    assert.ErrorAs(t, err, &credErr)
    assert.Equal(t, "username", credErr.Field)
}
```

**Coverage**: Username validation, early error detection

---

**C. ValidateCredentials Tests**

##### TC-AUTH-IMPL-004: ValidateCredentials_Valid

**Purpose**: Verify validation passes for valid credentials

**Test Strategy**:
```go
func TestValidateCredentials_Valid(t *testing.T) {
    testCases := []struct {
        name     string
        username string
        password string
    }{
        {"simple", "user", "pass"},
        {"special_chars_password", "user", "P@ss!w0rd#123"},
        {"unicode_password", "user", "пароль密码🔒"},
        {"spaces_in_password", "user", "pass with spaces"},
        {"quotes_in_password", "user", "pass\"with'quotes"},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Given: Provider with credentials
            provider := NewBasicAuthProvider(tc.username, tc.password)

            // When: Validating
            err := provider.ValidateCredentials()

            // Then: No error
            assert.NoError(t, err, "Should be valid: %s / %s", tc.username, tc.password)
        })
    }
}
```

**Coverage**: Validation logic, special character support, unicode support

---

##### TC-AUTH-IMPL-005: ValidateCredentials_EmptyUsername

**Purpose**: Verify error for empty username

**Test Strategy**:
```go
func TestValidateCredentials_EmptyUsername(t *testing.T) {
    // Given: Empty username
    provider := NewBasicAuthProvider("", "password")

    // When: Validating
    err := provider.ValidateCredentials()

    // Then: Returns error
    require.Error(t, err)

    var credErr *CredentialValidationError
    assert.ErrorAs(t, err, &credErr)
    assert.Equal(t, "username", credErr.Field)
    assert.Contains(t, credErr.Reason, "empty")
}
```

**Coverage**: Username required validation

---

##### TC-AUTH-IMPL-006: ValidateCredentials_EmptyPassword

**Purpose**: Verify error for empty password

**Test Strategy**:
```go
func TestValidateCredentials_EmptyPassword(t *testing.T) {
    // Given: Empty password
    provider := NewBasicAuthProvider("username", "")

    // When: Validating
    err := provider.ValidateCredentials()

    // Then: Returns error
    require.Error(t, err)

    var credErr *CredentialValidationError
    assert.ErrorAs(t, err, &credErr)
    assert.Equal(t, "password", credErr.Field)
}
```

**Coverage**: Password required validation

---

##### TC-AUTH-IMPL-007: ValidateCredentials_ControlCharactersInUsername

**Purpose**: Verify error for control characters in username

**Test Strategy**:
```go
func TestValidateCredentials_ControlCharactersInUsername(t *testing.T) {
    // Given: Username with control characters
    controlCharUsernames := []string{
        "user\n",       // newline
        "user\t",       // tab
        "user\x00",     // null byte
        "user\x1b",     // escape
    }

    for _, username := range controlCharUsernames {
        t.Run("control_char", func(t *testing.T) {
            provider := NewBasicAuthProvider(username, "password")

            // When: Validating
            err := provider.ValidateCredentials()

            // Then: Returns error
            require.Error(t, err)

            var credErr *CredentialValidationError
            assert.ErrorAs(t, err, &credErr)
            assert.Equal(t, "username", credErr.Field)
            assert.Contains(t, credErr.Reason, "control")
        })
    }
}
```

**Coverage**: Security validation, control character detection

---

**Auth Package Test Summary:**
- Total Tests: ~10-12 test cases
- Coverage Target: 90% (security-critical)
- Focus: Credential validation, special character support, security checks

---

### 4.4 Package: tls (Effort 1.2.4)

**Files to Test:**
- `pkg/tls/config.go` (Wave 2 implementation)

**Test File:**
- `pkg/tls/config_test.go`

#### Test Categories

**A. Constructor Tests**

##### TC-TLS-IMPL-001: NewConfigProvider_SecureMode

**Purpose**: Verify secure mode configuration

**Test Strategy**:
```go
package tls

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewConfigProvider_SecureMode(t *testing.T) {
    // Given: Secure mode (false = verify certificates)
    insecure := false

    // When: Creating provider
    provider := NewConfigProvider(insecure)

    // Then: Provider created
    require.NotNil(t, provider)

    // Verify satisfies interface
    var _ ConfigProvider = provider

    // Verify insecure flag is correct
    assert.False(t, provider.IsInsecure(), "Should be in secure mode")
}
```

**Coverage**: Constructor, secure mode setup

---

##### TC-TLS-IMPL-002: NewConfigProvider_InsecureMode

**Purpose**: Verify insecure mode configuration

**Test Strategy**:
```go
func TestNewConfigProvider_InsecureMode(t *testing.T) {
    // Given: Insecure mode (true = skip verification)
    insecure := true

    // When: Creating provider
    provider := NewConfigProvider(insecure)

    // Then: Provider created in insecure mode
    require.NotNil(t, provider)
    assert.True(t, provider.IsInsecure(), "Should be in insecure mode")
}
```

**Coverage**: Constructor, insecure mode setup

---

**B. GetTLSConfig Tests**

##### TC-TLS-IMPL-003: GetTLSConfig_SecureMode

**Purpose**: Verify TLS config for secure mode

**Test Strategy**:
```go
func TestGetTLSConfig_SecureMode(t *testing.T) {
    // Given: Provider in secure mode
    provider := NewConfigProvider(false)

    // When: Getting TLS config
    tlsConfig := provider.GetTLSConfig()

    // Then: Config has proper settings
    require.NotNil(t, tlsConfig)
    assert.False(t, tlsConfig.InsecureSkipVerify, "Should verify certificates")

    // Verify system cert pool loaded
    assert.NotNil(t, tlsConfig.RootCAs, "Should have system cert pool")
}
```

**Coverage**: TLS config generation, system cert pool loading

---

##### TC-TLS-IMPL-004: GetTLSConfig_InsecureMode

**Purpose**: Verify TLS config for insecure mode

**Test Strategy**:
```go
func TestGetTLSConfig_InsecureMode(t *testing.T) {
    // Given: Provider in insecure mode
    provider := NewConfigProvider(true)

    // When: Getting TLS config
    tlsConfig := provider.GetTLSConfig()

    // Then: InsecureSkipVerify is true
    require.NotNil(t, tlsConfig)
    assert.True(t, tlsConfig.InsecureSkipVerify, "Should skip certificate verification")
}
```

**Coverage**: Insecure TLS config generation

---

**C. IsInsecure Tests**

##### TC-TLS-IMPL-005: IsInsecure_Secure

**Purpose**: Verify IsInsecure returns false for secure mode

**Test Strategy**:
```go
func TestIsInsecure_Secure(t *testing.T) {
    // Given: Provider in secure mode
    provider := NewConfigProvider(false)

    // When: Checking if insecure
    insecure := provider.IsInsecure()

    // Then: Returns false
    assert.False(t, insecure)
}
```

**Coverage**: IsInsecure getter, secure mode flag

---

##### TC-TLS-IMPL-006: IsInsecure_Insecure

**Purpose**: Verify IsInsecure returns true for insecure mode

**Test Strategy**:
```go
func TestIsInsecure_Insecure(t *testing.T) {
    // Given: Provider in insecure mode
    provider := NewConfigProvider(true)

    // When: Checking if insecure
    insecure := provider.IsInsecure()

    // Then: Returns true
    assert.True(t, insecure)
}
```

**Coverage**: IsInsecure getter, insecure mode flag

---

**D. Integration Tests**

##### TC-TLS-IMPL-007: TLSConfig_UsableWithHTTPClient

**Purpose**: Verify TLS config works with http.Client

**Test Strategy**:
```go
func TestTLSConfig_UsableWithHTTPClient(t *testing.T) {
    // Given: TLS provider
    provider := NewConfigProvider(true)
    tlsConfig := provider.GetTLSConfig()

    // When: Creating HTTP transport with TLS config
    transport := &http.Transport{
        TLSClientConfig: tlsConfig,
    }

    client := &http.Client{
        Transport: transport,
    }

    // Then: Client created successfully
    require.NotNil(t, client)

    // Verify transport has TLS config
    assert.Equal(t, tlsConfig, transport.TLSClientConfig)
}
```

**Coverage**: Integration with http.Client, real-world usage

---

**TLS Package Test Summary:**
- Total Tests: ~8-10 test cases
- Coverage Target: 90% (security-critical)
- Focus: TLS configuration, secure/insecure modes, certificate handling

---

## 5. Test Execution Strategy

### 5.1 Running Wave 2 Tests

**Test Command:**
```bash
# Run all Wave 2 tests
go test ./pkg/docker ./pkg/registry ./pkg/auth ./pkg/tls -v

# With coverage
go test ./pkg/docker ./pkg/registry ./pkg/auth ./pkg/tls -cover -coverprofile=coverage.out

# Coverage report
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

**Expected Output:**
```
=== RUN   TestNewClient_Success
--- PASS: TestNewClient_Success (0.01s)
=== RUN   TestImageExists_ImagePresent
--- PASS: TestImageExists_ImagePresent (0.02s)
...
PASS
coverage: 87.3% of statements in ./pkg/docker
coverage: 86.1% of statements in ./pkg/registry
coverage: 92.4% of statements in ./pkg/auth
coverage: 91.2% of statements in ./pkg/tls
```

### 5.2 Test Execution Order

**Sequential Package Testing (for debugging):**
1. `pkg/docker` tests (no dependencies)
2. `pkg/registry` tests (uses docker via mocks)
3. `pkg/auth` tests (used by registry)
4. `pkg/tls` tests (used by registry)

**Parallel Execution (for CI):**
```bash
# All packages can test in parallel
go test ./pkg/... -v -parallel=4
```

### 5.3 CI Integration

**GitHub Actions Example:**
```yaml
name: Wave 2 Implementation Tests

on: [push, pull_request]

jobs:
  test-wave2:
    runs-on: ubuntu-latest

    services:
      docker:
        image: docker:dind
        options: --privileged

    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'

      - name: Pull Test Images
        run: |
          docker pull alpine:latest

      - name: Run Tests
        run: |
          go test ./pkg/docker ./pkg/registry ./pkg/auth ./pkg/tls -v -cover

      - name: Check Coverage
        run: |
          go test ./pkg/... -coverprofile=coverage.out
          go tool cover -func=coverage.out

          # Verify minimum coverage
          for pkg in docker registry auth tls; do
            coverage=$(go test ./pkg/$pkg -coverprofile=temp.out -v 2>&1 | grep coverage | awk '{print $4}' | sed 's/%//')
            min=85
            if [ "$pkg" = "auth" ] || [ "$pkg" = "tls" ]; then
              min=90
            fi
            if (( $(echo "$coverage < $min" | bc -l) )); then
              echo "Coverage for pkg/$pkg is $coverage%, below minimum $min%"
              exit 1
            fi
          done
```

---

## 6. Test Infrastructure Requirements

### 6.1 Dependencies

**go.mod additions for Wave 2 testing:**
```go
require (
    // Already in Wave 1:
    github.com/stretchr/testify v1.10.0
    github.com/google/go-containerregistry v0.19.0

    // NEW for Wave 2 implementation:
    github.com/docker/docker v28.2.2+incompatible  // Docker Engine API
    github.com/docker/go-connections v0.4.0         // Docker client helpers

    // For test registry (optional):
    // github.com/google/go-containerregistry/pkg/registry v0.19.0
)
```

### 6.2 Test Prerequisites

**Environment Setup:**
```bash
# 1. Docker daemon must be running
docker info

# 2. Pull test images
docker pull alpine:latest

# 3. Optional: Start test registry
docker run -d -p 5000:5000 --name registry registry:2

# 4. Run tests
go test ./pkg/... -v
```

### 6.3 Test Helpers

**Common test utilities (create in pkg/internal/testutil/):**

```go
// testutil/docker.go
package testutil

import (
    "testing"

    "github.com/cnoe-io/idpbuilder/pkg/docker"
    v1 "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/empty"
    "github.com/google/go-containerregistry/pkg/v1/mutate"
    "github.com/google/go-containerregistry/pkg/v1/random"
)

// CreateTestImage creates a minimal test OCI image
func CreateTestImage(t *testing.T) v1.Image {
    t.Helper()

    // Create random test image
    image, err := random.Image(1024, 1)
    if err != nil {
        t.Fatalf("Failed to create test image: %v", err)
    }

    return image
}

// SetupDockerClient creates a Docker client or skips test if Docker unavailable
func SetupDockerClient(t *testing.T) docker.Client {
    t.Helper()

    client, err := docker.NewClient()
    if err != nil {
        t.Skipf("Docker daemon not available: %v", err)
    }

    t.Cleanup(func() {
        client.Close()
    })

    return client
}
```

**Common test utilities (continued):**

```go
// testutil/registry.go
package testutil

import (
    "testing"

    "github.com/cnoe-io/idpbuilder/pkg/auth"
    "github.com/cnoe-io/idpbuilder/pkg/registry"
    tlspkg "github.com/cnoe-io/idpbuilder/pkg/tls"
)

// SetupRegistryClient creates a registry client with test providers
func SetupRegistryClient(t *testing.T) registry.Client {
    t.Helper()

    authProvider := &MockAuthProvider{
        username: "testuser",
        password: "testpass",
    }

    tlsProvider := &MockTLSProvider{
        insecure: true,
    }

    client, err := registry.NewClient(authProvider, tlsProvider)
    if err != nil {
        t.Fatalf("Failed to create registry client: %v", err)
    }

    return client
}
```

---

## 7. Progressive Realism Validation

### 7.1 Verification Checklist

**Every test MUST use real Wave 1 interfaces:**

- [ ] Imports reference actual packages: `github.com/cnoe-io/idpbuilder/pkg/docker`
- [ ] Mocks implement actual Wave 1 interfaces (compile-time checked)
- [ ] Error types match Wave 1 definitions exactly
- [ ] Type conversions use real go-containerregistry types (v1.Image, authn.Authenticator)
- [ ] No imaginary/abstract structures

**Example Validation:**
```go
// ✅ CORRECT: Real interface from Wave 1
import "github.com/cnoe-io/idpbuilder/pkg/docker"

type mockDockerClient struct {
    docker.Client  // Real interface
}

// ❌ WRONG: Imaginary interface
type mockDockerClient struct {
    SomeField string  // Not based on real interface
}
```

### 7.2 Test Plan Quality Gates

**Must Pass to Approve Test Plan:**
- [ ] All test cases reference actual Wave 1 types
- [ ] All imports are to real packages (no placeholders)
- [ ] All error types match Wave 1 definitions
- [ ] Mock implementations satisfy real interfaces
- [ ] No abstract/imaginary test infrastructure

---

## 8. Coverage Verification

### 8.1 Per-Package Coverage Enforcement

**Automated Coverage Check:**
```bash
#!/bin/bash
# coverage-check.sh

PACKAGES=("docker" "registry" "auth" "tls")
FAILED=0

for pkg in "${PACKAGES[@]}"; do
    echo "Checking coverage for pkg/$pkg..."

    # Get coverage percentage
    coverage=$(go test ./pkg/$pkg -coverprofile=temp.out 2>&1 | \
               grep coverage | awk '{print $4}' | sed 's/%//')

    # Determine minimum based on package
    if [ "$pkg" = "auth" ] || [ "$pkg" = "tls" ]; then
        min=90
    else
        min=85
    fi

    # Check if meets minimum
    if (( $(echo "$coverage < $min" | bc -l) )); then
        echo "❌ FAIL: pkg/$pkg coverage $coverage% < $min%"
        FAILED=1
    else
        echo "✅ PASS: pkg/$pkg coverage $coverage% >= $min%"
    fi
done

exit $FAILED
```

### 8.2 Coverage Exclusions

**What is NOT counted (per R007):**
- Test files: `*_test.go`
- Generated code: `*.pb.go`, `*_generated.go`
- Demo files: `demo-*.go`, `example-*.go`
- Documentation: `doc.go` (though simple, not complex logic)

---

## 9. Test Plan Summary

### 9.1 Test Count Estimates

| Package | Unit Tests | Edge Cases | Integration | Total |
|---------|-----------|------------|-------------|-------|
| docker | 8 | 4 | 2 | ~14 |
| registry | 10 | 5 | 3 | ~18 |
| auth | 6 | 4 | 0 | ~10 |
| tls | 5 | 2 | 1 | ~8 |
| **TOTAL** | **29** | **15** | **6** | **~50** |

**Wave 1 Comparison:**
- Wave 1: 33 tests (interface validation)
- Wave 2: ~50 tests (implementation testing)
- Increase: ~50% more tests (expected for implementations)

### 9.2 Test Execution Time Estimates

**Per Package:**
- docker: ~10 seconds (Docker daemon calls)
- registry: ~15 seconds (network operations, may use mocks)
- auth: ~2 seconds (pure logic)
- tls: ~1 second (pure logic)

**Total**: ~30 seconds for full test suite

---

## 10. R341 TDD Compliance Summary

### 10.1 TDD Requirements Met

- ✅ Test plan created BEFORE Wave 2 implementation
- ✅ All test cases specify expected behavior
- ✅ Tests reference actual Wave 1 interfaces (Progressive Realism)
- ✅ Coverage targets defined and measurable
- ✅ Test cases cover happy paths, error paths, edge cases
- ✅ Security validation included (auth/tls)

### 10.2 Implementation Guidance

**For SW Engineers implementing Wave 2:**

1. **Read this test plan first** - understand what tests expect
2. **Reference Wave 2 architecture** - understand WHAT to build
3. **Write code to pass tests** - tests define success criteria
4. **Run tests frequently** - verify progress
5. **Check coverage** - ensure meeting targets (85%/90%)
6. **No test modifications** - tests are the contract

---

## 11. Acceptance Criteria

### 11.1 Test Plan Acceptance

**This test plan is approved when:**
- [ ] All packages have test cases defined
- [ ] All test cases reference actual Wave 1 types
- [ ] Progressive Realism verified (no imaginary structures)
- [ ] Coverage targets specified (85%+ / 90%+)
- [ ] Test execution strategy defined
- [ ] R341 TDD compliance documented

### 11.2 Implementation Acceptance

**Wave 2 implementation passes when:**
- [ ] All ~50 test cases pass
- [ ] Coverage ≥85% for docker/registry
- [ ] Coverage ≥90% for auth/tls
- [ ] No test modifications required
- [ ] CI pipeline passes

---

## 12. Document Status

**Status**: ✅ READY FOR WAVE 2 IMPLEMENTATION
**Created**: 2025-10-29 05:54:00 UTC
**Test Planner**: @agent-code-reviewer
**Wave**: Wave 2 of Phase 1
**R341 Compliance**: TDD (Tests BEFORE Implementation)
**Progressive Realism**: Uses ACTUAL Wave 1 code as infrastructure

**Compliance Summary:**
- ✅ R341: Tests designed before Wave 2 implementation
- ✅ Progressive Realism: References actual Wave 1 interfaces
- ✅ Test strategy appropriate for implementation wave
- ✅ Coverage targets defined (85%/90%)
- ✅ Clear acceptance criteria defined
- ✅ Wave 1 Test Plan patterns followed

**Next Steps:**
1. Orchestrator spawns SW Engineers for Wave 2 efforts
2. SW Engineers read this test plan + Wave 2 architecture
3. SW Engineers implement code to pass these tests
4. SW Engineers run tests frequently during development
5. Code Reviewer validates test coverage in review

**Critical Success Factors:**
- Tests use REAL Wave 1 interfaces (not imaginary)
- Implementation passes tests without test modifications
- Coverage targets achieved (85%/90%)
- Security validation passes (auth/tls)

---

**END OF WAVE 2 TEST PLAN**
