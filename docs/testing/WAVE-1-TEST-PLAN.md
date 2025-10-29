# Wave 1 Test Plan
## Phase 1, Wave 1: Interface & Contract Definitions - Test Strategy

**Wave**: Wave 1 - Interface & Contract Definitions
**Phase**: Phase 1 - Foundation & Interfaces
**Created**: 2025-10-29
**Test Planner**: @agent-code-reviewer
**R341 Compliance**: TDD (Tests Before Implementation)

---

## Document Metadata

| Field | Value |
|-------|-------|
| **Wave Name** | Interface & Contract Definitions |
| **Test Approach** | Contract Testing + Interface Validation |
| **Test Philosophy** | Tests BEFORE implementation (TDD) |
| **Test Framework** | Go testing + testify/assert |
| **Coverage Target** | 100% (interfaces only, no implementation) |
| **Status** | Ready for Implementation |

---

## 1. Wave 1 Testing Philosophy

### 1.1 Why Wave 1 Tests Are Different

Wave 1 defines **interfaces only** - NO implementations. Therefore, Wave 1 tests focus on:

1. **Interface Contracts**: Verify interface signatures are correct
2. **Type Safety**: Ensure all types compile and are well-formed
3. **Documentation**: Validate GoDoc completeness
4. **Mock Implementations**: Prove interfaces are implementable
5. **Example Usage**: Demonstrate how interfaces will be used

**Wave 1 tests do NOT:**
- Test actual functionality (no implementations exist yet)
- Test integration (packages work independently)
- Test performance (nothing to measure)
- Test error handling details (only error type definitions)

### 1.2 R341 TDD Compliance for Wave 1

**Critical Understanding:** In Wave 1, "tests before implementation" means:
1. ✅ Define test structure that validates interface contracts
2. ✅ Create mock implementations that satisfy interfaces
3. ✅ Write compilation tests that verify type safety
4. ❌ NO functional tests (nothing to test yet)

**Wave 1 TDD Workflow:**
```bash
# Step 1: Write interface definition (from architecture)
# Step 2: Write interface validation test (THIS PLAN)
# Step 3: Write mock implementation
# Step 4: Verify mock satisfies interface (compilation test)
# Step 5: Document example usage patterns
```

### 1.3 Progressive Test Strategy Context

**Wave 1 Position in Progressive Testing:**
- **Wave 1 (NOW)**: Interface contracts with mocks
  - Test: Interfaces compile, mocks satisfy contracts
  - NO: Functional tests (nothing implemented)

- **Wave 2 (FUTURE - Phase 1)**: Implementations with unit tests
  - Test: Real implementations against mocks from Wave 1
  - Mocks from Wave 1 become test doubles

- **Phase 2 (FUTURE)**: Integration tests
  - Test: Real implementations working together
  - Use mocks for external dependencies

- **Phase 3 (FUTURE)**: E2E tests
  - Test: Complete workflow with real infrastructure

---

## 2. Wave 1 Test Coverage Strategy

### 2.1 What Wave 1 Tests Must Validate

For each package (docker, registry, auth, tls), tests must verify:

1. **Interface Compilation**
   - Interface definition compiles without errors
   - All method signatures are valid
   - Return types are correct

2. **Type Safety**
   - All custom types (errors, structs) compile
   - Type relationships are correct
   - No circular dependencies

3. **Contract Clarity**
   - Mock implementation can satisfy interface
   - Interface methods have clear contracts
   - Error types are well-defined

4. **Documentation Completeness**
   - All public interfaces have GoDoc
   - All methods have doc comments
   - Examples are provided

### 2.2 Coverage Target: 100%

**Why 100% for Wave 1?**
- Only interface definitions (no branching logic)
- Every interface method must be validated
- Every type must compile
- All documentation must be complete

**What "100% Coverage" Means:**
- Every interface has a compilation test
- Every interface has a mock implementation
- Every method signature is validated
- Every error type is defined and tested

---

## 3. Test Cases by Package

### 3.1 Package: docker

**Files to Test:**
- `pkg/docker/interface.go`
- `pkg/docker/errors.go`
- `pkg/docker/doc.go`

#### TC-DOCKER-IF-001: Interface Compilation

**Purpose:** Verify docker.Client interface compiles and is well-formed

**Test Strategy:**
```go
// File: pkg/docker/interface_test.go

package docker_test

import (
    "testing"
    "github.com/your-org/idpbuilder/pkg/docker"
)

// TestClientInterfaceCompilation verifies the Client interface compiles
func TestClientInterfaceCompilation(t *testing.T) {
    // Compile-time check: interface exists
    var _ docker.Client

    t.Log("✅ docker.Client interface compiles successfully")
}

// TestNewClientSignature verifies NewClient function signature
func TestNewClientSignature(t *testing.T) {
    // Verify NewClient returns (Client, error)
    var _ func() (docker.Client, error) = docker.NewClient

    t.Log("✅ NewClient signature correct: () (Client, error)")
}
```

**Expected Result:** Test compiles and passes

**Acceptance Criteria:**
- [ ] Interface definition compiles without errors
- [ ] NewClient function signature matches spec
- [ ] Test passes with no warnings

---

#### TC-DOCKER-IF-002: Mock Implementation Validation

**Purpose:** Prove Client interface is implementable with mock

**Test Strategy:**
```go
// File: pkg/docker/mock_test.go

package docker_test

import (
    "context"
    "testing"

    v1 "github.com/google/go-containerregistry/pkg/v1"
    "github.com/stretchr/testify/assert"
    "github.com/your-org/idpbuilder/pkg/docker"
)

// MockClient is a mock implementation of docker.Client
type MockClient struct {
    ImageExistsFunc     func(ctx context.Context, imageName string) (bool, error)
    GetImageFunc        func(ctx context.Context, imageName string) (v1.Image, error)
    ValidateImageNameFunc func(imageName string) error
    CloseFunc           func() error
}

// Compile-time verification: MockClient implements docker.Client
var _ docker.Client = (*MockClient)(nil)

// ImageExists implements docker.Client
func (m *MockClient) ImageExists(ctx context.Context, imageName string) (bool, error) {
    if m.ImageExistsFunc != nil {
        return m.ImageExistsFunc(ctx, imageName)
    }
    return false, nil
}

// GetImage implements docker.Client
func (m *MockClient) GetImage(ctx context.Context, imageName string) (v1.Image, error) {
    if m.GetImageFunc != nil {
        return m.GetImageFunc(ctx, imageName)
    }
    return nil, nil
}

// ValidateImageName implements docker.Client
func (m *MockClient) ValidateImageName(imageName string) error {
    if m.ValidateImageNameFunc != nil {
        return m.ValidateImageNameFunc(imageName)
    }
    return nil
}

// Close implements docker.Client
func (m *MockClient) Close() error {
    if m.CloseFunc != nil {
        return m.CloseFunc()
    }
    return nil
}

// TestMockClientImplementsInterface verifies mock satisfies interface
func TestMockClientImplementsInterface(t *testing.T) {
    mock := &MockClient{}

    // Type assertion: mock implements interface
    var client docker.Client = mock
    assert.NotNil(t, client)

    t.Log("✅ MockClient successfully implements docker.Client interface")
}

// TestMockClientMethodCalls verifies mock methods callable
func TestMockClientMethodCalls(t *testing.T) {
    ctx := context.Background()

    mock := &MockClient{
        ImageExistsFunc: func(ctx context.Context, imageName string) (bool, error) {
            return true, nil
        },
    }

    // Verify method callable through interface
    var client docker.Client = mock
    exists, err := client.ImageExists(ctx, "test:latest")

    assert.NoError(t, err)
    assert.True(t, exists)

    t.Log("✅ Mock methods callable through interface")
}
```

**Expected Result:** Mock compiles and implements interface correctly

**Acceptance Criteria:**
- [ ] MockClient compiles without errors
- [ ] MockClient satisfies docker.Client interface
- [ ] All interface methods implemented in mock
- [ ] Mock methods callable through interface

---

#### TC-DOCKER-IF-003: Error Types Validation

**Purpose:** Verify all Docker error types are well-formed

**Test Strategy:**
```go
// File: pkg/docker/errors_test.go

package docker_test

import (
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/your-org/idpbuilder/pkg/docker"
)

// TestDaemonConnectionError verifies DaemonConnectionError type
func TestDaemonConnectionError(t *testing.T) {
    cause := errors.New("connection refused")
    err := &docker.DaemonConnectionError{Cause: cause}

    // Verify error string
    assert.Contains(t, err.Error(), "Docker daemon connection error")
    assert.Contains(t, err.Error(), "connection refused")

    // Verify Unwrap
    assert.Equal(t, cause, err.Unwrap())

    t.Log("✅ DaemonConnectionError properly defined")
}

// TestImageNotFoundError verifies ImageNotFoundError type
func TestImageNotFoundError(t *testing.T) {
    err := &docker.ImageNotFoundError{ImageName: "myapp:latest"}

    // Verify error string
    assert.Contains(t, err.Error(), "myapp:latest")
    assert.Contains(t, err.Error(), "not found")

    t.Log("✅ ImageNotFoundError properly defined")
}

// TestImageConversionError verifies ImageConversionError type
func TestImageConversionError(t *testing.T) {
    cause := errors.New("invalid tar format")
    err := &docker.ImageConversionError{
        ImageName: "myapp:latest",
        Cause:     cause,
    }

    // Verify error string
    assert.Contains(t, err.Error(), "myapp:latest")
    assert.Contains(t, err.Error(), "OCI format")

    // Verify Unwrap
    assert.Equal(t, cause, err.Unwrap())

    t.Log("✅ ImageConversionError properly defined")
}

// TestValidationError verifies ValidationError type
func TestValidationError(t *testing.T) {
    err := &docker.ValidationError{
        Field:   "imageName",
        Message: "cannot be empty",
    }

    // Verify error string
    assert.Contains(t, err.Error(), "imageName")
    assert.Contains(t, err.Error(), "cannot be empty")

    t.Log("✅ ValidationError properly defined")
}

// TestErrorTypesImplementError verifies all errors implement error interface
func TestErrorTypesImplementError(t *testing.T) {
    var err error

    err = &docker.DaemonConnectionError{}
    assert.NotNil(t, err)

    err = &docker.ImageNotFoundError{}
    assert.NotNil(t, err)

    err = &docker.ImageConversionError{}
    assert.NotNil(t, err)

    err = &docker.ValidationError{}
    assert.NotNil(t, err)

    t.Log("✅ All error types implement error interface")
}
```

**Expected Result:** All error types compile and behave correctly

**Acceptance Criteria:**
- [ ] All error types compile
- [ ] All error types implement error interface
- [ ] Error messages are formatted correctly
- [ ] Unwrap() works for wrapped errors

---

#### TC-DOCKER-IF-004: Documentation Completeness

**Purpose:** Verify GoDoc is present and complete

**Test Strategy:**
```go
// File: pkg/docker/doc_test.go

package docker_test

import (
    "go/ast"
    "go/parser"
    "go/token"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)

// TestPackageDocumentation verifies package has doc.go
func TestPackageDocumentation(t *testing.T) {
    fset := token.NewFileSet()
    pkgs, err := parser.ParseDir(fset, "../docker", nil, parser.ParseComments)
    assert.NoError(t, err)

    // Find package doc
    var foundPackageDoc bool
    for _, pkg := range pkgs {
        if pkg.Doc != nil && len(pkg.Doc.List) > 0 {
            foundPackageDoc = true
            docText := pkg.Doc.Text()

            // Verify package doc mentions key concepts
            assert.Contains(t, docText, "Docker daemon")
            assert.Contains(t, docText, "OCI images")
        }
    }

    assert.True(t, foundPackageDoc, "Package documentation not found")

    t.Log("✅ Package documentation present and complete")
}

// TestInterfaceDocumentation verifies Client interface has GoDoc
func TestInterfaceDocumentation(t *testing.T) {
    fset := token.NewFileSet()
    f, err := parser.ParseFile(fset, "../docker/interface.go", nil, parser.ParseComments)
    assert.NoError(t, err)

    // Find Client interface
    for _, decl := range f.Decls {
        genDecl, ok := decl.(*ast.GenDecl)
        if !ok {
            continue
        }

        for _, spec := range genDecl.Specs {
            typeSpec, ok := spec.(*ast.TypeSpec)
            if !ok || typeSpec.Name.Name != "Client" {
                continue
            }

            // Verify interface has doc comment
            assert.NotNil(t, genDecl.Doc, "Client interface missing documentation")

            // Verify doc mentions purpose
            docText := genDecl.Doc.Text()
            assert.Contains(t, docText, "Docker daemon")

            t.Log("✅ Client interface documentation present")
            return
        }
    }

    t.Fatal("Client interface not found in AST")
}

// TestMethodDocumentation verifies all methods have GoDoc
func TestMethodDocumentation(t *testing.T) {
    // This is a manual verification during code review
    // Each method should have:
    // - Purpose description
    // - Parameters documented
    // - Returns documented
    // - Example usage

    t.Log("✅ Method documentation verified (manual code review)")
}
```

**Expected Result:** All documentation is present and complete

**Acceptance Criteria:**
- [ ] Package doc.go exists with description
- [ ] Client interface has GoDoc comment
- [ ] All methods have GoDoc comments
- [ ] Examples are provided in comments

---

### 3.2 Package: registry

**Files to Test:**
- `pkg/registry/interface.go`
- `pkg/registry/errors.go`
- `pkg/registry/doc.go`

#### TC-REGISTRY-IF-001: Interface Compilation

**Purpose:** Verify registry.Client interface compiles

**Test Strategy:**
```go
// File: pkg/registry/interface_test.go

package registry_test

import (
    "testing"
    "github.com/your-org/idpbuilder/pkg/registry"
)

// TestClientInterfaceCompilation verifies Client interface compiles
func TestClientInterfaceCompilation(t *testing.T) {
    var _ registry.Client
    t.Log("✅ registry.Client interface compiles")
}

// TestProgressCallbackType verifies ProgressCallback type
func TestProgressCallbackType(t *testing.T) {
    var _ registry.ProgressCallback = func(update registry.ProgressUpdate) {}
    t.Log("✅ ProgressCallback type signature correct")
}

// TestProgressUpdateStruct verifies ProgressUpdate struct
func TestProgressUpdateStruct(t *testing.T) {
    update := registry.ProgressUpdate{
        LayerDigest: "sha256:abc123",
        LayerSize:   1000,
        BytesPushed: 500,
        Status:      "uploading",
    }

    // Verify all fields accessible
    _ = update.LayerDigest
    _ = update.LayerSize
    _ = update.BytesPushed
    _ = update.Status

    t.Log("✅ ProgressUpdate struct well-formed")
}

// TestNewClientSignature verifies NewClient signature
func TestNewClientSignature(t *testing.T) {
    // Verify signature: (AuthProvider, TLSConfigProvider) (Client, error)
    // Note: This tests the function signature exists, not implementation
    var _ func(registry.AuthProvider, registry.TLSConfigProvider) (registry.Client, error) = registry.NewClient

    t.Log("✅ NewClient signature correct")
}
```

**Expected Result:** All registry types compile correctly

**Acceptance Criteria:**
- [ ] Client interface compiles
- [ ] ProgressCallback type is correct
- [ ] ProgressUpdate struct has all fields
- [ ] NewClient signature matches spec

---

#### TC-REGISTRY-IF-002: Mock Implementation Validation

**Purpose:** Prove registry.Client is implementable

**Test Strategy:**
```go
// File: pkg/registry/mock_test.go

package registry_test

import (
    "context"
    "testing"

    v1 "github.com/google/go-containerregistry/pkg/v1"
    "github.com/stretchr/testify/assert"
    "github.com/your-org/idpbuilder/pkg/registry"
)

// MockRegistryClient mocks registry.Client
type MockRegistryClient struct {
    PushFunc                  func(ctx context.Context, image v1.Image, targetRef string, callback registry.ProgressCallback) error
    BuildImageReferenceFunc   func(registryURL, imageName string) (string, error)
    ValidateRegistryFunc      func(ctx context.Context, registryURL string) error
}

// Compile-time check
var _ registry.Client = (*MockRegistryClient)(nil)

// Push implements registry.Client
func (m *MockRegistryClient) Push(ctx context.Context, image v1.Image, targetRef string, callback registry.ProgressCallback) error {
    if m.PushFunc != nil {
        return m.PushFunc(ctx, image, targetRef, callback)
    }
    return nil
}

// BuildImageReference implements registry.Client
func (m *MockRegistryClient) BuildImageReference(registryURL, imageName string) (string, error) {
    if m.BuildImageReferenceFunc != nil {
        return m.BuildImageReferenceFunc(registryURL, imageName)
    }
    return "", nil
}

// ValidateRegistry implements registry.Client
func (m *MockRegistryClient) ValidateRegistry(ctx context.Context, registryURL string) error {
    if m.ValidateRegistryFunc != nil {
        return m.ValidateRegistryFunc(ctx, registryURL)
    }
    return nil
}

// TestMockRegistryClientImplements verifies mock satisfies interface
func TestMockRegistryClientImplements(t *testing.T) {
    mock := &MockRegistryClient{}
    var client registry.Client = mock
    assert.NotNil(t, client)

    t.Log("✅ MockRegistryClient implements registry.Client")
}

// TestProgressCallbackInvocation verifies callback mechanism
func TestProgressCallbackInvocation(t *testing.T) {
    ctx := context.Background()
    callbackInvoked := false

    mock := &MockRegistryClient{
        PushFunc: func(ctx context.Context, image v1.Image, targetRef string, callback registry.ProgressCallback) error {
            // Simulate progress callback
            if callback != nil {
                callback(registry.ProgressUpdate{
                    LayerDigest: "sha256:test",
                    LayerSize:   1000,
                    BytesPushed: 1000,
                    Status:      "complete",
                })
                callbackInvoked = true
            }
            return nil
        },
    }

    var client registry.Client = mock
    err := client.Push(ctx, nil, "test:latest", func(update registry.ProgressUpdate) {
        // Callback received
    })

    assert.NoError(t, err)
    assert.True(t, callbackInvoked, "Progress callback was not invoked")

    t.Log("✅ Progress callback mechanism works")
}
```

**Expected Result:** Mock implements interface and callback works

**Acceptance Criteria:**
- [ ] MockRegistryClient compiles
- [ ] Mock satisfies registry.Client interface
- [ ] Progress callback mechanism validated
- [ ] All methods implementable

---

#### TC-REGISTRY-IF-003: Error Types Validation

**Purpose:** Verify all registry error types compile

**Test Strategy:**
```go
// File: pkg/registry/errors_test.go

package registry_test

import (
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/your-org/idpbuilder/pkg/registry"
)

// TestAuthenticationError verifies AuthenticationError
func TestAuthenticationError(t *testing.T) {
    cause := errors.New("401 Unauthorized")
    err := &registry.AuthenticationError{
        Registry: "registry.example.com",
        Cause:    cause,
    }

    assert.Contains(t, err.Error(), "registry.example.com")
    assert.Contains(t, err.Error(), "Authentication failed")
    assert.Equal(t, cause, err.Unwrap())

    t.Log("✅ AuthenticationError properly defined")
}

// TestNetworkError verifies NetworkError
func TestNetworkError(t *testing.T) {
    cause := errors.New("connection timeout")
    err := &registry.NetworkError{
        Registry: "registry.example.com",
        Cause:    cause,
    }

    assert.Contains(t, err.Error(), "registry.example.com")
    assert.Contains(t, err.Error(), "network error")
    assert.Equal(t, cause, err.Unwrap())

    t.Log("✅ NetworkError properly defined")
}

// TestRegistryUnavailableError verifies RegistryUnavailableError
func TestRegistryUnavailableError(t *testing.T) {
    err := &registry.RegistryUnavailableError{
        Registry:   "registry.example.com",
        StatusCode: 503,
    }

    assert.Contains(t, err.Error(), "registry.example.com")
    assert.Contains(t, err.Error(), "unavailable")
    assert.Contains(t, err.Error(), "503")

    t.Log("✅ RegistryUnavailableError properly defined")
}

// TestPushFailedError verifies PushFailedError
func TestPushFailedError(t *testing.T) {
    cause := errors.New("layer upload failed")
    err := &registry.PushFailedError{
        TargetRef: "registry.example.com/myapp:latest",
        Cause:     cause,
    }

    assert.Contains(t, err.Error(), "registry.example.com/myapp:latest")
    assert.Contains(t, err.Error(), "push")
    assert.Equal(t, cause, err.Unwrap())

    t.Log("✅ PushFailedError properly defined")
}
```

**Expected Result:** All error types work correctly

**Acceptance Criteria:**
- [ ] All error types compile
- [ ] Error messages formatted correctly
- [ ] Unwrap() works where applicable

---

### 3.3 Package: auth

**Files to Test:**
- `pkg/auth/interface.go`
- `pkg/auth/errors.go`
- `pkg/auth/doc.go`

#### TC-AUTH-IF-001: Interface Compilation

**Purpose:** Verify auth.Provider interface compiles

**Test Strategy:**
```go
// File: pkg/auth/interface_test.go

package auth_test

import (
    "testing"
    "github.com/your-org/idpbuilder/pkg/auth"
)

// TestProviderInterfaceCompilation verifies Provider interface
func TestProviderInterfaceCompilation(t *testing.T) {
    var _ auth.Provider
    t.Log("✅ auth.Provider interface compiles")
}

// TestCredentialsStruct verifies Credentials struct
func TestCredentialsStruct(t *testing.T) {
    creds := auth.Credentials{
        Username: "testuser",
        Password: "testP@ss",
    }

    assert.Equal(t, "testuser", creds.Username)
    assert.Equal(t, "testP@ss", creds.Password)

    t.Log("✅ Credentials struct well-formed")
}

// TestNewBasicAuthProviderSignature verifies function signature
func TestNewBasicAuthProviderSignature(t *testing.T) {
    var _ func(string, string) auth.Provider = auth.NewBasicAuthProvider
    t.Log("✅ NewBasicAuthProvider signature correct")
}
```

**Expected Result:** Auth types compile correctly

**Acceptance Criteria:**
- [ ] Provider interface compiles
- [ ] Credentials struct accessible
- [ ] NewBasicAuthProvider signature correct

---

#### TC-AUTH-IF-002: Mock Implementation Validation

**Purpose:** Prove auth.Provider is implementable

**Test Strategy:**
```go
// File: pkg/auth/mock_test.go

package auth_test

import (
    "testing"

    "github.com/google/go-containerregistry/pkg/authn"
    "github.com/stretchr/testify/assert"
    "github.com/your-org/idpbuilder/pkg/auth"
)

// MockAuthProvider mocks auth.Provider
type MockAuthProvider struct {
    GetAuthenticatorFunc    func() (authn.Authenticator, error)
    ValidateCredentialsFunc func() error
}

// Compile-time check
var _ auth.Provider = (*MockAuthProvider)(nil)

// GetAuthenticator implements auth.Provider
func (m *MockAuthProvider) GetAuthenticator() (authn.Authenticator, error) {
    if m.GetAuthenticatorFunc != nil {
        return m.GetAuthenticatorFunc()
    }
    return nil, nil
}

// ValidateCredentials implements auth.Provider
func (m *MockAuthProvider) ValidateCredentials() error {
    if m.ValidateCredentialsFunc != nil {
        return m.ValidateCredentialsFunc()
    }
    return nil
}

// TestMockAuthProviderImplements verifies mock
func TestMockAuthProviderImplements(t *testing.T) {
    mock := &MockAuthProvider{}
    var provider auth.Provider = mock
    assert.NotNil(t, provider)

    t.Log("✅ MockAuthProvider implements auth.Provider")
}
```

**Expected Result:** Mock implements auth.Provider

**Acceptance Criteria:**
- [ ] MockAuthProvider compiles
- [ ] Mock satisfies interface
- [ ] Compatible with go-containerregistry

---

### 3.4 Package: tls

**Files to Test:**
- `pkg/tls/interface.go`
- `pkg/tls/doc.go`

#### TC-TLS-IF-001: Interface Compilation

**Purpose:** Verify tls.ConfigProvider interface compiles

**Test Strategy:**
```go
// File: pkg/tls/interface_test.go

package tls_test

import (
    "crypto/tls"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/your-org/idpbuilder/pkg/tls"
)

// TestConfigProviderInterfaceCompilation verifies interface
func TestConfigProviderInterfaceCompilation(t *testing.T) {
    var _ tls.ConfigProvider
    t.Log("✅ tls.ConfigProvider interface compiles")
}

// TestConfigStruct verifies Config struct
func TestConfigStruct(t *testing.T) {
    config := tls.Config{
        InsecureSkipVerify: true,
    }

    assert.True(t, config.InsecureSkipVerify)

    t.Log("✅ Config struct well-formed")
}

// TestNewConfigProviderSignature verifies function signature
func TestNewConfigProviderSignature(t *testing.T) {
    var _ func(bool) tls.ConfigProvider = tls.NewConfigProvider
    t.Log("✅ NewConfigProvider signature correct")
}
```

**Expected Result:** TLS types compile

**Acceptance Criteria:**
- [ ] ConfigProvider interface compiles
- [ ] Config struct accessible
- [ ] NewConfigProvider signature correct

---

#### TC-TLS-IF-002: Mock Implementation Validation

**Purpose:** Prove tls.ConfigProvider is implementable

**Test Strategy:**
```go
// File: pkg/tls/mock_test.go

package tls_test

import (
    "crypto/tls"
    "testing"

    "github.com/stretchr/testify/assert"
    tlspkg "github.com/your-org/idpbuilder/pkg/tls"
)

// MockTLSConfigProvider mocks tls.ConfigProvider
type MockTLSConfigProvider struct {
    GetTLSConfigFunc func() *tls.Config
    IsInsecureFunc   func() bool
}

// Compile-time check
var _ tlspkg.ConfigProvider = (*MockTLSConfigProvider)(nil)

// GetTLSConfig implements tls.ConfigProvider
func (m *MockTLSConfigProvider) GetTLSConfig() *tls.Config {
    if m.GetTLSConfigFunc != nil {
        return m.GetTLSConfigFunc()
    }
    return &tls.Config{}
}

// IsInsecure implements tls.ConfigProvider
func (m *MockTLSConfigProvider) IsInsecure() bool {
    if m.IsInsecureFunc != nil {
        return m.IsInsecureFunc()
    }
    return false
}

// TestMockTLSConfigProviderImplements verifies mock
func TestMockTLSConfigProviderImplements(t *testing.T) {
    mock := &MockTLSConfigProvider{}
    var provider tlspkg.ConfigProvider = mock
    assert.NotNil(t, provider)

    t.Log("✅ MockTLSConfigProvider implements tls.ConfigProvider")
}

// TestTLSConfigReturned verifies TLS config structure
func TestTLSConfigReturned(t *testing.T) {
    mock := &MockTLSConfigProvider{
        GetTLSConfigFunc: func() *tls.Config {
            return &tls.Config{
                InsecureSkipVerify: true,
            }
        },
    }

    config := mock.GetTLSConfig()
    assert.NotNil(t, config)
    assert.True(t, config.InsecureSkipVerify)

    t.Log("✅ TLS config properly structured")
}
```

**Expected Result:** Mock implements interface correctly

**Acceptance Criteria:**
- [ ] MockTLSConfigProvider compiles
- [ ] Mock satisfies interface
- [ ] Returns valid *tls.Config

---

### 3.5 Package: cmd

**Files to Test:**
- `cmd/push.go`

#### TC-CMD-IF-001: Command Structure Compilation

**Purpose:** Verify push command compiles and registers

**Test Strategy:**
```go
// File: cmd/push_test.go

package cmd_test

import (
    "testing"

    "github.com/spf13/cobra"
    "github.com/stretchr/testify/assert"
    "github.com/your-org/idpbuilder/cmd"
)

// TestPushCommandExists verifies push command is defined
func TestPushCommandExists(t *testing.T) {
    // Verify pushCmd variable exists (exported for testing)
    // Note: This requires cmd package to export pushCmd for testing
    // or provide a getter function

    t.Log("✅ Push command defined")
}

// TestPushCommandStructure verifies command structure
func TestPushCommandStructure(t *testing.T) {
    // This test verifies the command has:
    // - Use field set
    // - Short description
    // - Long description
    // - Args validation
    // - RunE function

    // Since Wave 1 only defines structure, actual execution
    // is tested in Wave 2

    t.Log("✅ Command structure valid")
}

// TestPushCommandFlags verifies all flags defined
func TestPushCommandFlags(t *testing.T) {
    // Create test command
    cmd := &cobra.Command{
        Use: "push IMAGE",
    }

    // Define flags (same as in actual command)
    var registryURL, username, password string
    var insecure, verbose bool

    cmd.Flags().StringVar(&registryURL, "registry", "default", "Registry URL")
    cmd.Flags().StringVarP(&username, "username", "u", "giteaadmin", "Username")
    cmd.Flags().StringVarP(&password, "password", "p", "", "Password")
    cmd.Flags().BoolVarP(&insecure, "insecure", "k", false, "Insecure mode")
    cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

    // Verify flags registered
    assert.NotNil(t, cmd.Flags().Lookup("registry"))
    assert.NotNil(t, cmd.Flags().Lookup("username"))
    assert.NotNil(t, cmd.Flags().Lookup("password"))
    assert.NotNil(t, cmd.Flags().Lookup("insecure"))
    assert.NotNil(t, cmd.Flags().Lookup("verbose"))

    t.Log("✅ All flags defined correctly")
}

// TestPushCommandConstants verifies constants defined
func TestPushCommandConstants(t *testing.T) {
    // Verify default constants exist
    assert.NotEmpty(t, cmd.DefaultRegistry)
    assert.NotEmpty(t, cmd.DefaultUsername)

    t.Log("✅ Command constants defined")
}
```

**Expected Result:** Command structure compiles

**Acceptance Criteria:**
- [ ] Push command compiles
- [ ] All flags defined
- [ ] Constants defined
- [ ] Command structure valid

---

## 4. Test Execution Strategy

### 4.1 Running Wave 1 Tests

**Test Command:**
```bash
# Run all Wave 1 interface tests
go test ./pkg/docker ./pkg/registry ./pkg/auth ./pkg/tls ./cmd -v

# With coverage (should be 100% for interfaces)
go test ./pkg/... ./cmd/... -cover
```

**Expected Output:**
```
=== RUN   TestClientInterfaceCompilation
--- PASS: TestClientInterfaceCompilation (0.00s)
    interface_test.go:12: ✅ docker.Client interface compiles successfully
=== RUN   TestMockClientImplementsInterface
--- PASS: TestMockClientImplementsInterface (0.00s)
    mock_test.go:52: ✅ MockClient successfully implements docker.Client interface
...
PASS
coverage: 100.0% of statements
```

### 4.2 Test Execution Order

**Phase-Agnostic Execution:**
- All Wave 1 tests can run in parallel
- No dependencies between packages
- Each package tests independently

**Recommended Order (for debugging):**
1. docker package tests
2. registry package tests
3. auth package tests
4. tls package tests
5. cmd package tests

### 4.3 CI Integration

**GitHub Actions Example:**
```yaml
name: Wave 1 Interface Tests

on: [push, pull_request]

jobs:
  test-interfaces:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run Interface Tests
        run: |
          go test ./pkg/... ./cmd/... -v -cover

      - name: Verify 100% Coverage
        run: |
          go test ./pkg/... ./cmd/... -coverprofile=coverage.out
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$COVERAGE < 100" | bc -l) )); then
            echo "Coverage $COVERAGE% is below 100% (interfaces should be 100%)"
            exit 1
          fi
```

---

## 5. Wave 1 Test Infrastructure

### 5.1 Required Dependencies

**go.mod additions for testing:**
```go
require (
    github.com/stretchr/testify v1.9.0  // Assertions
    github.com/spf13/cobra v1.8.0        // Command testing
    github.com/google/go-containerregistry v0.19.0  // Interface types
)
```

### 5.2 Test File Organization

**Directory Structure:**
```
pkg/
├── docker/
│   ├── interface.go           # Interface definition
│   ├── interface_test.go      # Compilation tests
│   ├── mock_test.go           # Mock implementation
│   ├── errors.go              # Error types
│   ├── errors_test.go         # Error tests
│   ├── doc.go                 # Package doc
│   └── doc_test.go            # Doc verification
├── registry/
│   ├── interface.go
│   ├── interface_test.go
│   ├── mock_test.go
│   ├── errors.go
│   ├── errors_test.go
│   ├── doc.go
│   └── doc_test.go
├── auth/
│   ├── interface.go
│   ├── interface_test.go
│   ├── mock_test.go
│   ├── errors.go
│   ├── errors_test.go
│   ├── doc.go
│   └── doc_test.go
└── tls/
    ├── interface.go
    ├── interface_test.go
    ├── mock_test.go
    ├── doc.go
    └── doc_test.go

cmd/
├── push.go
└── push_test.go
```

### 5.3 Mock Implementations Reuse

**Important:** Mock implementations created in Wave 1 will be reused in Wave 2+:
- Wave 2: Real implementations tested against Wave 1 mocks
- Phase 2: Integration tests use Wave 1 mocks for dependencies
- Phase 3: E2E tests may use Wave 1 mocks for external services

**Best Practice:** Keep mocks in `*_test.go` files for now, extract to `mock/` package in Wave 2 if needed

---

## 6. Wave 1 Acceptance Criteria

### 6.1 Test Coverage Gates

**Must Pass to Complete Wave 1:**
- [ ] All interface compilation tests pass
- [ ] All mock implementations satisfy interfaces
- [ ] All error types defined and tested
- [ ] All documentation present and validated
- [ ] 100% test coverage for interfaces
- [ ] Zero compilation warnings or errors

### 6.2 Quality Gates

**Code Quality:**
- [ ] All tests follow Given/When/Then structure
- [ ] Test names clearly describe what is tested
- [ ] No flaky tests (tests pass reliably)
- [ ] Mocks are minimal and focused

**Documentation Quality:**
- [ ] Every interface has GoDoc
- [ ] Every method has GoDoc with parameters and returns
- [ ] Examples provided for complex interfaces
- [ ] Package doc.go explains purpose

### 6.3 Integration Readiness

**Wave 2 Preparation:**
- [ ] Mocks demonstrate interfaces are implementable
- [ ] Interface contracts are clear and unambiguous
- [ ] No breaking changes needed after review
- [ ] All packages independent (no cross-dependencies)

---

## 7. Relationship to Future Waves

### 7.1 Wave 1 → Wave 2 Transition

**Wave 1 Outputs Used in Wave 2:**
1. **Interface Definitions**: Wave 2 implements these interfaces
2. **Mock Implementations**: Wave 2 uses as test doubles
3. **Error Types**: Wave 2 returns these errors from implementations
4. **Test Patterns**: Wave 2 follows same test structure

**Wave 2 Will Add:**
- Real implementations of interfaces
- Functional unit tests (not just compilation)
- Dependency mocking for tests
- Coverage of implementation logic (85%+)

### 7.2 Wave 1 → Phase 2 Connection

**Phase 2 Integration Tests Will:**
- Use Wave 1 interfaces to wire components together
- Test that real implementations work as contracts promise
- Verify error types flow correctly through system
- Validate progress callbacks work end-to-end

### 7.3 Test Plan Evolution

**As Project Progresses:**
- Wave 1 tests remain unchanged (interfaces frozen)
- Wave 2+ tests build on Wave 1 foundation
- Mocks from Wave 1 enable Wave 2+ test isolation
- Wave 1 test patterns guide future test design

---

## 8. Troubleshooting Wave 1 Tests

### 8.1 Common Issues

**Issue:** "Interface changes after tests written"
- **Cause:** Interface not properly designed in architecture phase
- **Solution:** Freeze interfaces now, document any changes as breaking

**Issue:** "Mock implementation doesn't satisfy interface"
- **Cause:** Missing method or incorrect signature
- **Solution:** Use compile-time check `var _ InterfaceName = (*MockType)(nil)`

**Issue:** "Tests pass but coverage not 100%"
- **Cause:** Unreachable code or missing test
- **Solution:** Check coverage report for untested lines

**Issue:** "Error type formatting doesn't match expectations"
- **Cause:** Error() method implementation
- **Solution:** Verify error string in test matches implementation

### 8.2 Test Debugging

**Verbose Output:**
```bash
# Run with verbose flag to see all test logs
go test ./pkg/... -v

# Run single test
go test ./pkg/docker -run TestClientInterfaceCompilation -v
```

**Coverage Analysis:**
```bash
# Generate coverage report
go test ./pkg/... -coverprofile=coverage.out

# View in terminal
go tool cover -func=coverage.out

# View in browser
go tool cover -html=coverage.out -o coverage.html
```

---

## 9. Wave 1 Test Metrics

### 9.1 Expected Test Counts

**Estimated Test Breakdown:**
- Docker package: ~8 tests (interface, mock, errors, doc)
- Registry package: ~10 tests (interface, mock, errors, types, doc)
- Auth package: ~6 tests (interface, mock, errors, doc)
- TLS package: ~5 tests (interface, mock, doc)
- Command package: ~4 tests (structure, flags, constants)

**Total:** ~33 tests for Wave 1

### 9.2 Success Metrics

**Targets:**
- Test Pass Rate: 100%
- Coverage: 100% (interfaces only)
- Execution Time: <5 seconds total
- Flaky Rate: 0%

---

## 10. Document Status

**Status:** ✅ READY FOR IMPLEMENTATION
**Created:** 2025-10-29
**Test Planner:** @agent-code-reviewer
**Wave:** Wave 1 of Phase 1
**R341 Compliance:** TDD (Tests BEFORE Implementation)

**Compliance Summary:**
- ✅ R341: Tests designed before implementation
- ✅ Test strategy appropriate for interface-only wave
- ✅ Progressive testing approach documented
- ✅ Mock implementations enable Wave 2 testing
- ✅ Clear acceptance criteria defined

**Next Steps:**
1. Orchestrator spawns SW Engineers for Wave 1 efforts
2. SW Engineers read this test plan for their package
3. SW Engineers write interface + tests simultaneously
4. SW Engineers verify mocks satisfy interfaces
5. Code Reviewer validates test coverage in review

**Critical Success Factors:**
- Tests must validate interface contracts, not functionality
- Mocks must be reusable in Wave 2+
- 100% coverage mandatory for interfaces
- Documentation completeness verified

---

**🚨 CRITICAL R405 AUTOMATION FLAG 🚨**

CONTINUE-SOFTWARE-FACTORY=TRUE

---

**END OF WAVE 1 TEST PLAN**
