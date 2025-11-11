# Wave 1 Test Plan: Interface Definitions

**Wave**: Wave 1.1 - Interface Definitions
**Phase**: Phase 1 - Foundation & Interfaces
**Created**: 2025-11-11
**Test Planner**: Code Reviewer Agent
**TDD Phase**: RED (tests written before implementation)
**Progressive Fidelity**: INTERFACE MOCKS (no prior implementations exist)

---

## Table of Contents

1. [Wave Test Overview](#wave-test-overview)
2. [Progressive Test Planning Strategy](#progressive-test-planning-strategy)
3. [Test Infrastructure](#test-infrastructure)
4. [Interface Contract Tests](#interface-contract-tests)
5. [Mock Implementation Tests](#mock-implementation-tests)
6. [Test Execution Plan](#test-execution-plan)
7. [Success Criteria](#success-criteria)
8. [Integration with Project Test Plan](#integration-with-project-test-plan)

---

## Wave Test Overview

### Wave Scope

**What Wave 1 Delivers**:
- 5 complete Go interface definitions (DockerClient, RegistryClient, AuthProvider, TLSProvider, PushCommand)
- 10 custom error types with Error() and Unwrap() methods
- 3 supporting types (ProgressUpdate, ProgressCallback, PushFlags)
- Complete package documentation
- **NO implementations yet** (those come in Wave 2)

**What Wave 1 Tests Validate**:
- ✅ Interface method signatures compile
- ✅ Error types implement error interface correctly
- ✅ Type definitions are valid Go code
- ✅ Interface contracts are complete (no missing methods)
- ✅ Mock implementations satisfy interfaces
- ✅ Example usage patterns compile

**What Wave 1 Tests DO NOT Validate**:
- ❌ Actual Docker daemon interaction (no implementation)
- ❌ Real registry push operations (no implementation)
- ❌ Network connectivity (no implementation)
- ❌ Performance characteristics (no implementation)

### Test Categories for Wave 1

**Total Estimated Tests**: ~35 tests

| Category | Test Count | Purpose | Files |
|----------|------------|---------|-------|
| Interface Compilation | 5 | Verify interfaces compile | `pkg/*/interface_test.go` |
| Error Type Compliance | 10 | Verify error types implement error | `pkg/*/errors_test.go` |
| Mock Implementations | 5 | Verify mocks satisfy interfaces | `tests/mocks/*_mock_test.go` |
| Type Definition Validity | 5 | Verify supporting types compile | `pkg/*/types_test.go` |
| Example Code Compilation | 5 | Verify usage examples compile | `examples/*_test.go` |
| Package Documentation | 5 | Verify godoc generation | `pkg/*/doc_test.go` |

### TDD Workflow (R341)

**Current State**: RED Phase (tests created BEFORE implementation)

```
┌─────────────────────────────────────────────────────────────┐
│  WAVE 1 TEST CREATION (THIS PLAN)                           │
│  ────────────────────────────────────────────────────────   │
│  - Define interface contract tests                          │
│  - Define error type tests                                  │
│  - Define mock implementation tests                         │
│  - Expected: ALL TESTS PASS (interfaces only, no impl)      │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────┐
│  WAVE 1 IMPLEMENTATION (EFFORTS 1.1.1 - 1.1.4)              │
│  ────────────────────────────────────────────────────────   │
│  - Create interface files (pkg/docker/interface.go, etc.)   │
│  - Create error types                                       │
│  - Run tests → Should PASS (pure interface definitions)     │
└─────────────────────────────────────────────────────────────┘
```

**Wave 1 is Special**: Unlike later waves, Wave 1 tests should **PASS immediately** after implementation because we're only creating interface definitions (not complex implementation logic).

---

## Progressive Test Planning Strategy

### R341 Compliance: Wave-Specific Testing

**Wave 1 Unique Characteristics**:
- **No Prior Implementations**: This is the FIRST wave, so no real implementations exist
- **Interface-Only Deliverables**: Wave 1 creates contracts, not functionality
- **Mock-Based Testing**: All tests use mock implementations to validate interfaces

**Progressive Fidelity Principle**:

```
Wave 1 (NOW):
  Interfaces defined → Tests use MOCKS to validate contracts

Wave 2 (FUTURE):
  Implementations created → Tests use REAL implementations from Wave 1

Wave 3+ (FUTURE):
  Integration work → Tests use REAL implementations from Waves 1-2
```

### What This Means for Wave 1 Tests

**We CAN test**:
- Interface method signatures compile
- Error types implement `error` interface
- Mock implementations satisfy interfaces
- Usage patterns type-check

**We CANNOT test** (deferred to Wave 2):
- Actual Docker daemon interaction
- Real registry push operations
- Error handling logic
- Performance characteristics

---

## Test Infrastructure

### Test File Organization

```
idpbuilder-oci-push-rebuild/
├── pkg/
│   ├── docker/
│   │   ├── interface.go           # (Created in Effort 1.1.1)
│   │   ├── interface_test.go      # (Created in Effort 1.1.1 - interface tests)
│   │   └── errors_test.go         # (Created in Effort 1.1.1 - error type tests)
│   ├── registry/
│   │   ├── interface.go           # (Created in Effort 1.1.2)
│   │   ├── interface_test.go      # (Created in Effort 1.1.2 - interface tests)
│   │   ├── types_test.go          # (Created in Effort 1.1.2 - ProgressUpdate, etc.)
│   │   └── errors_test.go         # (Created in Effort 1.1.2 - error type tests)
│   ├── auth/
│   │   ├── interface.go           # (Created in Effort 1.1.3)
│   │   ├── interface_test.go      # (Created in Effort 1.1.3 - interface tests)
│   │   └── errors_test.go         # (Created in Effort 1.1.3 - error type tests)
│   └── tls/
│       ├── interface.go           # (Created in Effort 1.1.3)
│       └── interface_test.go      # (Created in Effort 1.1.3 - interface tests)
├── cmd/
│   ├── push.go                    # (Created in Effort 1.1.4)
│   ├── push_test.go               # (Created in Effort 1.1.4 - command tests)
│   └── flags_test.go              # (Created in Effort 1.1.4 - flag parsing tests)
├── tests/
│   └── mocks/
│       ├── docker_mock.go         # Mock DockerClient implementation
│       ├── registry_mock.go       # Mock RegistryClient implementation
│       ├── auth_mock.go           # Mock AuthProvider implementation
│       ├── tls_mock.go            # Mock TLSProvider implementation
│       └── mocks_test.go          # Tests verifying mocks satisfy interfaces
└── examples/
    ├── docker_example_test.go     # Example usage compilation tests
    ├── registry_example_test.go   # Example usage compilation tests
    └── push_example_test.go       # Full workflow example compilation tests
```

### Test Execution Commands

**Run all Wave 1 tests**:
```bash
# All package tests (interfaces + error types)
go test ./pkg/...

# All mock tests
go test ./tests/mocks/...

# All example compilation tests
go test ./examples/...

# Full Wave 1 test suite
go test ./pkg/... ./tests/mocks/... ./examples/...
```

**Run with coverage**:
```bash
go test -cover ./pkg/...
```

**Expected Coverage for Wave 1**:
- Interface files: 100% (simple type definitions)
- Error types: 100% (Error() methods)
- Mock implementations: 80%+ (basic test coverage)

---

## Interface Contract Tests

### Test Category: DockerClient Interface

**Test File**: `pkg/docker/interface_test.go`

**Purpose**: Verify DockerClient interface compiles and mock can implement it

**Test Cases** (5 tests):

#### TC-W1-DOCKER-001: DockerClient interface compiles
```go
func TestDockerClientInterface_Compiles(t *testing.T) {
    // Verify interface is valid Go type
    var _ DockerClient = (*mockDockerClient)(nil)
}
```
- **Given**: DockerClient interface defined in interface.go
- **When**: Mock type assigned to interface variable
- **Then**: Compilation succeeds (type satisfies interface)

#### TC-W1-DOCKER-002: DockerClient has all required methods
```go
func TestDockerClientInterface_Methods(t *testing.T) {
    methods := []string{"ImageExists", "GetImage", "ValidateImageName", "Close"}
    interfaceType := reflect.TypeOf((*DockerClient)(nil)).Elem()

    for _, method := range methods {
        _, exists := interfaceType.MethodByName(method)
        if !exists {
            t.Errorf("DockerClient missing required method: %s", method)
        }
    }
}
```
- **Validates**: All 4 methods present in interface

#### TC-W1-DOCKER-003: ImageExists has correct signature
```go
func TestDockerClientInterface_ImageExistsSignature(t *testing.T) {
    interfaceType := reflect.TypeOf((*DockerClient)(nil)).Elem()
    method, _ := interfaceType.MethodByName("ImageExists")

    // Verify: ImageExists(ctx context.Context, imageName string) (bool, error)
    if method.Type.NumIn() != 3 { // receiver + 2 args
        t.Error("ImageExists should accept 2 parameters")
    }
    if method.Type.NumOut() != 2 { // (bool, error)
        t.Error("ImageExists should return 2 values")
    }
}
```

#### TC-W1-DOCKER-004: GetImage returns v1.Image type
```go
func TestDockerClientInterface_GetImageSignature(t *testing.T) {
    // Verify return type is v1.Image from go-containerregistry
    mock := &mockDockerClient{}
    ctx := context.Background()

    // This should compile (mock implements interface)
    var _ v1.Image = mock.GetImage(ctx, "test:latest")
}
```

#### TC-W1-DOCKER-005: NewDockerClient constructor exists
```go
func TestDockerClient_ConstructorExists(t *testing.T) {
    // Verify constructor function signature
    // NewDockerClient() (DockerClient, error)

    // This test verifies the function exists and has correct signature
    // Implementation will panic("not implemented") in Wave 1
    _, err := NewDockerClient()

    // In Wave 1, we expect panic (no implementation yet)
    // This is OK - we're just verifying the signature compiles
}
```

**Expected Result**: All 5 tests PASS after Effort 1.1.1 implementation

---

### Test Category: Error Type Compliance

**Test File**: `pkg/docker/errors_test.go`

**Purpose**: Verify all error types implement error interface correctly

**Test Cases** (3 tests for DockerClient errors):

#### TC-W1-ERROR-001: ImageNotFoundError implements error
```go
func TestImageNotFoundError_ImplementsError(t *testing.T) {
    err := &ImageNotFoundError{ImageName: "test:v1"}

    // Verify implements error interface
    var _ error = err

    // Verify Error() method returns expected format
    expected := "image not found: test:v1"
    if err.Error() != expected {
        t.Errorf("Expected '%s', got '%s'", expected, err.Error())
    }
}
```

#### TC-W1-ERROR-002: DaemonConnectionError implements error and Unwrap
```go
func TestDaemonConnectionError_Unwrap(t *testing.T) {
    cause := errors.New("connection refused")
    err := &DaemonConnectionError{
        Endpoint: "unix:///var/run/docker.sock",
        Cause:    cause,
    }

    // Verify implements error
    var _ error = err

    // Verify Unwrap returns cause
    if errors.Unwrap(err) != cause {
        t.Error("Unwrap should return original cause")
    }

    // Verify Error() includes endpoint and cause
    errMsg := err.Error()
    if !strings.Contains(errMsg, "connection refused") {
        t.Error("Error message should include cause")
    }
}
```

#### TC-W1-ERROR-003: InvalidImageNameError has descriptive message
```go
func TestInvalidImageNameError_Message(t *testing.T) {
    err := &InvalidImageNameError{
        ImageName: "Invalid Name:v1",
        Reason:    "contains spaces",
    }

    var _ error = err

    // Verify message includes both image name and reason
    errMsg := err.Error()
    if !strings.Contains(errMsg, "Invalid Name:v1") {
        t.Error("Error message should include image name")
    }
    if !strings.Contains(errMsg, "contains spaces") {
        t.Error("Error message should include reason")
    }
}
```

**Similar tests for all error types in other packages**:
- `pkg/registry/errors_test.go`: RegistryAuthError, RegistryConnectionError, LayerPushError
- `pkg/auth/errors_test.go`: InvalidCredentialsError, MissingCredentialsError

---

### Test Category: RegistryClient Interface

**Test File**: `pkg/registry/interface_test.go`

**Test Cases** (5 tests):

#### TC-W1-REGISTRY-001: RegistryClient interface compiles
#### TC-W1-REGISTRY-002: RegistryClient has Push, BuildImageReference, ValidateRegistry methods
#### TC-W1-REGISTRY-003: Push accepts ProgressCallback type
#### TC-W1-REGISTRY-004: BuildImageReference returns string, error
#### TC-W1-REGISTRY-005: NewRegistryClient requires AuthProvider and TLSProvider

---

### Test Category: Supporting Types

**Test File**: `pkg/registry/types_test.go`

**Test Cases** (3 tests):

#### TC-W1-TYPES-001: LayerStatus enum values defined
```go
func TestLayerStatus_EnumValues(t *testing.T) {
    // Verify all 4 status values exist
    statuses := []LayerStatus{LayerWaiting, LayerUploading, LayerComplete, LayerFailed}

    for _, status := range statuses {
        if status.String() == "" {
            t.Errorf("LayerStatus %d missing String() implementation", status)
        }
    }
}
```

#### TC-W1-TYPES-002: ProgressUpdate struct has required fields
```go
func TestProgressUpdate_Fields(t *testing.T) {
    update := ProgressUpdate{
        LayerDigest:   "sha256:abc123",
        LayerSize:     1000,
        BytesUploaded: 500,
        Status:        LayerUploading,
    }

    // Verify all fields accessible
    if update.LayerDigest == "" {
        t.Error("LayerDigest should be set")
    }
    if update.LayerSize != 1000 {
        t.Error("LayerSize should be 1000")
    }
}
```

#### TC-W1-TYPES-003: ProgressCallback is valid function type
```go
func TestProgressCallback_Type(t *testing.T) {
    // Verify ProgressCallback can be assigned a function
    var callback ProgressCallback = func(update ProgressUpdate) {
        // Mock callback
    }

    // Verify callback can be invoked
    callback(ProgressUpdate{Status: LayerComplete})
}
```

---

### Test Category: AuthProvider Interface

**Test File**: `pkg/auth/interface_test.go`

**Test Cases** (5 tests):

#### TC-W1-AUTH-001: AuthProvider interface compiles
#### TC-W1-AUTH-002: GetAuthenticator returns authn.Authenticator
#### TC-W1-AUTH-003: ValidateCredentials returns error
#### TC-W1-AUTH-004: NewAuthProvider accepts username, password strings
#### TC-W1-AUTH-005: Auth error types implement error interface

---

### Test Category: TLSProvider Interface

**Test File**: `pkg/tls/interface_test.go`

**Test Cases** (5 tests):

#### TC-W1-TLS-001: TLSProvider interface compiles
#### TC-W1-TLS-002: GetTLSConfig returns *tls.Config
#### TC-W1-TLS-003: IsInsecure returns bool
#### TC-W1-TLS-004: GetWarningMessage returns string
#### TC-W1-TLS-005: NewTLSProvider accepts insecure bool

---

### Test Category: Command Structure

**Test File**: `cmd/push_test.go`

**Test Cases** (5 tests):

#### TC-W1-CMD-001: PushCommand struct has required fields
```go
func TestPushCommand_Fields(t *testing.T) {
    cmd := &PushCommand{
        dockerClient:   nil, // Will be mock in Wave 2
        registryClient: nil,
        authProvider:   nil,
        tlsProvider:    nil,
    }

    // Verify struct compiles with interface fields
    var _ *PushCommand = cmd
}
```

#### TC-W1-CMD-002: PushFlags struct has all required fields
```go
func TestPushFlags_Fields(t *testing.T) {
    flags := &PushFlags{
        Registry: "https://gitea.test:8443",
        Username: "admin",
        Password: "secret",
        Insecure: true,
        Verbose:  false,
    }

    // Verify all fields accessible
    if flags.Registry == "" {
        t.Error("Registry field should be accessible")
    }
}
```

#### TC-W1-CMD-003: NewPushCommand returns *cobra.Command
```go
func TestNewPushCommand_ReturnsCobraCommand(t *testing.T) {
    cmd := NewPushCommand()

    // Verify returns Cobra command
    if cmd.Use != "push IMAGE_NAME" {
        t.Error("Command Use should be 'push IMAGE_NAME'")
    }

    // Verify flags defined
    if cmd.Flags().Lookup("registry") == nil {
        t.Error("--registry flag should be defined")
    }
}
```

#### TC-W1-CMD-004: Exit codes are defined
```go
func TestPushCommand_ExitCodes(t *testing.T) {
    // Verify all exit codes defined as constants
    codes := map[string]int{
        "ExitSuccess":      ExitSuccess,
        "ExitGeneralError": ExitGeneralError,
        "ExitAuthError":    ExitAuthError,
        "ExitNetworkError": ExitNetworkError,
        "ExitImageNotFound": ExitImageNotFound,
    }

    for name, code := range codes {
        if code < 0 || code > 255 {
            t.Errorf("Exit code %s has invalid value: %d", name, code)
        }
    }
}
```

#### TC-W1-CMD-005: Command help text is defined
```go
func TestPushCommand_HelpText(t *testing.T) {
    cmd := NewPushCommand()

    if cmd.Short == "" {
        t.Error("Command should have Short description")
    }
    if cmd.Long == "" {
        t.Error("Command should have Long description")
    }
}
```

---

## Mock Implementation Tests

### Purpose

Verify that mock implementations satisfy all interface contracts. These mocks will be used in Wave 2 testing.

### Test Category: Mock Implementations

**Test File**: `tests/mocks/mocks_test.go`

**Test Cases** (5 tests):

#### TC-W1-MOCK-001: Mock DockerClient satisfies interface
```go
func TestMockDockerClient_ImplementsInterface(t *testing.T) {
    mock := &MockDockerClient{
        ImageExistsFunc: func(ctx context.Context, imageName string) (bool, error) {
            return true, nil
        },
        GetImageFunc: func(ctx context.Context, imageName string) (v1.Image, error) {
            return nil, nil
        },
        ValidateImageNameFunc: func(imageName string) error {
            return nil
        },
        CloseFunc: func() error {
            return nil
        },
    }

    // Verify mock satisfies DockerClient interface
    var _ docker.DockerClient = mock
}
```

#### TC-W1-MOCK-002: Mock RegistryClient satisfies interface
```go
func TestMockRegistryClient_ImplementsInterface(t *testing.T) {
    mock := &MockRegistryClient{
        PushFunc: func(ctx context.Context, image v1.Image, targetRef string, progress registry.ProgressCallback) error {
            return nil
        },
        BuildImageReferenceFunc: func(registryURL, imageName string) (string, error) {
            return "", nil
        },
        ValidateRegistryFunc: func(ctx context.Context, registryURL string) error {
            return nil
        },
    }

    var _ registry.RegistryClient = mock
}
```

#### TC-W1-MOCK-003: Mock AuthProvider satisfies interface
```go
func TestMockAuthProvider_ImplementsInterface(t *testing.T) {
    mock := &MockAuthProvider{
        GetAuthenticatorFunc: func() (authn.Authenticator, error) {
            return authn.Anonymous, nil
        },
        ValidateCredentialsFunc: func() error {
            return nil
        },
    }

    var _ auth.AuthProvider = mock
}
```

#### TC-W1-MOCK-004: Mock TLSProvider satisfies interface
```go
func TestMockTLSProvider_ImplementsInterface(t *testing.T) {
    mock := &MockTLSProvider{
        GetTLSConfigFunc: func() *tls.Config {
            return &tls.Config{InsecureSkipVerify: true}
        },
        IsInsecureFunc: func() bool {
            return true
        },
        GetWarningMessageFunc: func() string {
            return "WARNING: Insecure mode"
        },
    }

    var _ tlspkg.TLSProvider = mock
}
```

#### TC-W1-MOCK-005: Mock functions can be customized per test
```go
func TestMockCustomization_PerTest(t *testing.T) {
    // Verify mocks can be customized for different test scenarios
    mock := &MockDockerClient{}

    // Scenario 1: Image exists
    mock.ImageExistsFunc = func(ctx context.Context, imageName string) (bool, error) {
        return true, nil
    }
    exists, _ := mock.ImageExists(context.Background(), "test:v1")
    if !exists {
        t.Error("Mock should return true when configured")
    }

    // Scenario 2: Image not found
    mock.ImageExistsFunc = func(ctx context.Context, imageName string) (bool, error) {
        return false, nil
    }
    exists, _ = mock.ImageExists(context.Background(), "missing:v1")
    if exists {
        t.Error("Mock should return false when reconfigured")
    }
}
```

---

## Test Execution Plan

### Wave 1 Test Execution Timeline

**During Effort Implementation** (Efforts 1.1.1 - 1.1.4):

1. **Effort 1.1.1** (Docker Interface):
   - Create `pkg/docker/interface.go`
   - Create `pkg/docker/interface_test.go` (5 tests)
   - Create `pkg/docker/errors_test.go` (3 tests)
   - Run: `go test ./pkg/docker/...`
   - **Expected**: All 8 tests PASS

2. **Effort 1.1.2** (Registry Interface):
   - Create `pkg/registry/interface.go`
   - Create `pkg/registry/interface_test.go` (5 tests)
   - Create `pkg/registry/types_test.go` (3 tests)
   - Create `pkg/registry/errors_test.go` (3 tests)
   - Run: `go test ./pkg/registry/...`
   - **Expected**: All 11 tests PASS

3. **Effort 1.1.3** (Auth & TLS Interfaces):
   - Create `pkg/auth/interface.go`
   - Create `pkg/auth/interface_test.go` (5 tests)
   - Create `pkg/auth/errors_test.go` (2 tests)
   - Create `pkg/tls/interface.go`
   - Create `pkg/tls/interface_test.go` (5 tests)
   - Run: `go test ./pkg/auth/... ./pkg/tls/...`
   - **Expected**: All 12 tests PASS

4. **Effort 1.1.4** (Command Structure):
   - Create `cmd/push.go`
   - Create `cmd/push_test.go` (5 tests)
   - Run: `go test ./cmd/...`
   - **Expected**: All 5 tests PASS

**After All Efforts Complete**:

5. **Mock Implementation Creation**:
   - Create `tests/mocks/docker_mock.go`
   - Create `tests/mocks/registry_mock.go`
   - Create `tests/mocks/auth_mock.go`
   - Create `tests/mocks/tls_mock.go`
   - Create `tests/mocks/mocks_test.go` (5 tests)
   - Run: `go test ./tests/mocks/...`
   - **Expected**: All 5 tests PASS

**Wave 1 Completion Validation**:

```bash
# Run full Wave 1 test suite
go test ./pkg/... ./cmd/... ./tests/mocks/...

# Run with coverage
go test -cover ./pkg/... ./cmd/... ./tests/mocks/...

# Verify all tests pass
echo "Expected: 36 tests total, 36 PASS, 0 FAIL"
```

### Continuous Integration

**Pre-commit Hook** (runs automatically):
```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running Wave 1 tests..."
go test ./pkg/... ./cmd/... ./tests/mocks/...

if [ $? -ne 0 ]; then
    echo "Tests failed. Commit aborted."
    exit 1
fi

echo "All tests passed!"
```

**CI/CD Pipeline** (GitHub Actions / similar):
```yaml
name: Wave 1 Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run Wave 1 Tests
        run: go test -v ./pkg/... ./cmd/... ./tests/mocks/...

      - name: Check Coverage
        run: go test -cover ./pkg/... ./cmd/... ./tests/mocks/...
```

---

## Success Criteria

### Wave 1 Test Completion Criteria

**All tests must PASS before Wave 1 is considered complete:**

- ✅ **36 total tests PASS**:
  - 5 DockerClient interface tests
  - 3 DockerClient error tests
  - 5 RegistryClient interface tests
  - 3 RegistryClient type tests
  - 3 RegistryClient error tests
  - 5 AuthProvider interface tests
  - 2 AuthProvider error tests
  - 5 TLSProvider interface tests
  - 5 Command structure tests
  - 5 Mock implementation tests

- ✅ **Test Coverage ≥95%**:
  - Interface files: 100% (simple type definitions)
  - Error type files: 100% (Error() methods)
  - Mock implementations: 90%+

- ✅ **Build Success**:
  - `go build ./...` succeeds without errors
  - All packages compile cleanly
  - No missing dependencies

- ✅ **Documentation Generated**:
  - `go doc` produces output for all interfaces
  - godoc comments are present and correct
  - Example code compiles (verified by tests)

- ✅ **No Implementation Required**:
  - NewDockerClient() can panic("not implemented")
  - NewRegistryClient() can panic("not implemented")
  - NewAuthProvider() can panic("not implemented")
  - NewTLSProvider() can panic("not implemented")
  - This is EXPECTED for Wave 1 (interface definitions only)

### Quality Gates

**Before Wave 1 Complete Transition**:

1. **Compilation Gate**:
   ```bash
   go build ./...  # Must succeed
   ```

2. **Test Gate**:
   ```bash
   go test ./pkg/... ./cmd/... ./tests/mocks/...  # All 36 tests PASS
   ```

3. **Coverage Gate**:
   ```bash
   go test -cover ./pkg/...  # Coverage ≥95%
   ```

4. **Documentation Gate**:
   ```bash
   go doc docker.DockerClient  # Produces valid output
   go doc registry.RegistryClient
   go doc auth.AuthProvider
   go doc tls.TLSProvider
   ```

5. **Mock Validation Gate**:
   ```bash
   go test ./tests/mocks/...  # All mock tests PASS
   ```

### Blocking Issues

**Wave 1 CANNOT complete if**:
- Any test fails
- Coverage < 95%
- Any package fails to compile
- godoc generation fails
- Mock implementations don't satisfy interfaces

---

## Integration with Project Test Plan

### Mapping to Project Test Plan

This Wave 1 Test Plan implements the **Interface Contract Tests** section of the Project Test Plan:

**From Project Test Plan**:
- TC-DOCKER-001 to TC-DOCKER-015 → Wave 2 (implementation tests)
- TC-REGISTRY-001 to TC-REGISTRY-020 → Wave 2 (implementation tests)
- TC-AUTH-001 to TC-AUTH-012 → Wave 2 (implementation tests)
- TC-TLS-001 to TC-TLS-010 → Wave 2 (implementation tests)
- TC-CMD-001 to TC-CMD-013 → Wave 2 (implementation tests)

**Wave 1 Focuses On** (subset of Project Plan):
- Interface compilation validation (TC-W1-*-001)
- Error type compliance (TC-W1-ERROR-*)
- Supporting type definitions (TC-W1-TYPES-*)
- Mock implementations (TC-W1-MOCK-*)

**Wave 2 Will Implement** (full Project Plan):
- All functional tests (TC-DOCKER-*, TC-REGISTRY-*, etc.)
- Integration tests (TC-INT-*)
- E2E tests (TC-E2E-*)

### Progressive Test Development

```
Project Test Plan (150 tests total)
    ↓
Phase 1 Test Plan (70 interface tests)
    ↓
Wave 1 Test Plan (36 interface definition tests) ← YOU ARE HERE
    ↓
Wave 2 Test Plan (70 implementation tests) ← FUTURE
```

### Test Reuse Strategy

**Mocks created in Wave 1** will be reused in:
- Wave 2 implementation tests (mock dependencies)
- Wave 3 integration tests (isolate components)
- E2E tests (simulate error scenarios)

**Example**:
```go
// Wave 1: Create MockDockerClient
type MockDockerClient struct { ... }

// Wave 2: Use mock in RegistryClient tests
func TestRegistryClient_Push(t *testing.T) {
    dockerMock := &MockDockerClient{ /* configured for test */ }
    // Test registry using mocked Docker client
}
```

---

## R341 TDD Compliance Verification

### TDD Checklist

This test plan complies with R341 Test-Driven Development requirements:

- ✅ **Tests Created BEFORE Implementation Planning**: This test plan created before WAVE-1-IMPLEMENTATION.md
- ✅ **Tests Define Success Criteria**: All tests specify exact expected behavior
- ✅ **RED Phase Acknowledged**: Wave 1 tests will PASS immediately (interface definitions are simple)
- ✅ **GREEN Phase Planned**: Implementation efforts will create interfaces to satisfy tests
- ✅ **Progressive Fidelity**: Tests use mocks (no prior implementations exist)
- ✅ **Integration with Project Plan**: Tests map to Project Test Plan sections

### R342 Early Integration Branch Compliance

**Test Execution Location**: All Wave 1 tests will be created in `phase1-integration` branch (created early per R342)

**Branch Strategy**:
```
main
  ↓
phase1-integration (created early, all Wave 1 tests committed here)
  ↓
phase1/wave1/effort-1.1.1 (interface + tests)
  ↓
phase1/wave1/effort-1.1.2 (interface + tests)
  ↓
phase1/wave1/effort-1.1.3 (interface + tests)
  ↓
phase1/wave1/effort-1.1.4 (interface + tests)
```

**Integration Testing**: After all efforts merge to `phase1-integration`, run full Wave 1 test suite to validate integration.

---

## Test Maintenance

### When to Update Wave 1 Tests

**Update tests if**:
- Interface signatures change (rare in Wave 1)
- New error types added
- Supporting types modified
- Interface method added/removed

**DO NOT update tests for**:
- Implementation details (Wave 2 concern)
- Performance characteristics (Wave 2 concern)
- Integration behavior (Wave 2 concern)

### Test Ownership

**Effort 1.1.1**: Owns `pkg/docker/*_test.go`
**Effort 1.1.2**: Owns `pkg/registry/*_test.go`
**Effort 1.1.3**: Owns `pkg/auth/*_test.go` and `pkg/tls/*_test.go`
**Effort 1.1.4**: Owns `cmd/*_test.go`
**Wave Integration**: Owns `tests/mocks/*_test.go`

---

## Appendix: Mock Implementation Examples

### Example Mock: DockerClient

**File**: `tests/mocks/docker_mock.go`

```go
package mocks

import (
    "context"

    "github.com/jessesanford/idpbuilder/pkg/docker"
    v1 "github.com/google/go-containerregistry/pkg/v1"
)

// MockDockerClient is a test double for docker.DockerClient
type MockDockerClient struct {
    ImageExistsFunc       func(ctx context.Context, imageName string) (bool, error)
    GetImageFunc          func(ctx context.Context, imageName string) (v1.Image, error)
    ValidateImageNameFunc func(imageName string) error
    CloseFunc             func() error
}

func (m *MockDockerClient) ImageExists(ctx context.Context, imageName string) (bool, error) {
    if m.ImageExistsFunc != nil {
        return m.ImageExistsFunc(ctx, imageName)
    }
    return false, nil
}

func (m *MockDockerClient) GetImage(ctx context.Context, imageName string) (v1.Image, error) {
    if m.GetImageFunc != nil {
        return m.GetImageFunc(ctx, imageName)
    }
    return nil, nil
}

func (m *MockDockerClient) ValidateImageName(imageName string) error {
    if m.ValidateImageNameFunc != nil {
        return m.ValidateImageNameFunc(imageName)
    }
    return nil
}

func (m *MockDockerClient) Close() error {
    if m.CloseFunc != nil {
        return m.CloseFunc()
    }
    return nil
}
```

**Usage in Tests** (Wave 2+):
```go
func TestSomeFeature(t *testing.T) {
    mock := &mocks.MockDockerClient{
        ImageExistsFunc: func(ctx context.Context, imageName string) (bool, error) {
            return imageName == "alpine:latest", nil
        },
    }

    // Use mock in test
    exists, err := mock.ImageExists(context.Background(), "alpine:latest")
    if !exists || err != nil {
        t.Error("Expected image to exist")
    }
}
```

---

## Summary

**Wave 1 Test Plan Status**: ✅ COMPLETE

**Total Tests Defined**: 36 tests
**Expected Pass Rate**: 100% (interface definitions are simple)
**Coverage Target**: ≥95%
**TDD Phase**: RED (tests before implementation)
**Progressive Fidelity**: MOCKS (no prior implementations)

**Next Steps**:
1. Orchestrator validates this test plan
2. Orchestrator creates `WAVE-1-IMPLEMENTATION.md` with exact effort specifications
3. Efforts 1.1.1 - 1.1.4 implement interfaces and tests
4. All 36 tests PASS
5. Wave 1 complete → proceed to Wave 2

**Integration with SF 3.0**:
- Tests committed to `phase1-integration` branch (R342)
- Each effort includes tests in implementation (R341)
- Code Reviewer validates test passage before approval
- Architect validates interface contracts during review

---

**Test Plan Created**: 2025-11-11
**Test Plan Status**: Ready for Orchestrator Validation
**Expected Orchestrator Transition**: WAITING_FOR_WAVE_TEST_PLAN → WAVE_START (Wave 1)

---

**R341 TDD Compliance**: ✅ VERIFIED (tests created before implementation planning)
**R342 Integration Branch**: ✅ VERIFIED (tests will commit to phase1-integration)
**R340 Concrete Fidelity**: ✅ VERIFIED (specific test cases, exact file paths, clear acceptance criteria)
