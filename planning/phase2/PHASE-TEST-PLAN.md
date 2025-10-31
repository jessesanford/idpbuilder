# Phase 2 Test Plan - Core Push Functionality

**Phase**: Phase 2 - Core Push Functionality
**Created**: 2025-10-31
**Test Planner**: @agent-code-reviewer
**Test Approach**: TDD (Test-First)
**Coverage Target**: 85%+

---

## 🔴 TDD Compliance (R341)

**CRITICAL**: This test plan is created BEFORE implementation per R341.

- ✅ **Test-first**: Tests designed before any Phase 2 code
- ✅ **Phase 1 integration**: Leveraging actual Phase 1 implementations
- ✅ **R342 early integration**: Test plan committed to early integration branch
- ✅ **Progressive planning**: Using real Phase 1 fixtures and mocks

**Implementation Sequence**:
1. ✅ **NOW**: Phase 2 test plan (this document)
2. ⏭️ **NEXT**: Phase 2 implementation (guided by these tests)
3. ⏭️ **THEN**: Run tests during implementation
4. ⏭️ **FINALLY**: Achieve 85%+ coverage before completion

---

## Phase 1 Test Infrastructure (Available for Reuse)

### Test Patterns Established in Phase 1

**From `pkg/auth/basic_test.go` (31 test cases)**:
- Mock-based unit testing with `testify/assert` and `testify/require`
- Table-driven tests for credential validation
- Special character testing (unicode, quotes, spaces)
- Control character rejection in usernames
- Error type verification with `ErrorAs()`
- Given/When/Then test structure

**From `pkg/registry/client_test.go` (31 test cases)**:
- Mock providers for Auth and TLS
- `httptest.NewServer()` for registry endpoint testing
- Progress callback verification
- Empty image testing with `v1/empty.Image`
- Network error simulation
- Context cancellation testing

**From `pkg/docker/client_test.go` (G test suite)**:
- Docker daemon integration tests (requires running daemon)
- Image existence verification with `alpine:latest` prerequisite
- Command injection prevention testing
- Skippable manual tests for daemon unavailability
- Error type testing for coverage

**From `pkg/tls/config_test.go` (10 test cases)**:
- Secure vs insecure mode testing
- HTTP client integration verification
- System certificate pool validation
- Multiple calls return new instances

### Test Fixtures Available

**Mock Providers (from Phase 1)**:
```go
// From pkg/registry/client_test.go
type mockAuthProvider struct {
    authenticator authn.Authenticator
    validateErr   error
}

type mockTLSProvider struct {
    config   *tls.Config
    insecure bool
}
```

**Test Images**:
- `alpine:latest` - Phase 1 prerequisite, available for testing
- `v1/empty.Image` - go-containerregistry's empty image for unit tests
- Custom test images can be built during integration tests

**Test Servers**:
- `httptest.NewServer()` - For mocking registry HTTP endpoints
- Returns 200/401 for validation tests
- Can simulate network errors

### Testing Libraries in Use

**Core Libraries**:
- `github.com/stretchr/testify/assert` - Assertions
- `github.com/stretchr/testify/require` - Required assertions (fail fast)
- `github.com/google/go-containerregistry/pkg/v1` - OCI image types
- `github.com/google/go-containerregistry/pkg/authn` - Auth types

**Test Execution**:
- Standard Go test framework (`go test`)
- Context-aware testing
- Table-driven test patterns

---

## Phase 2 Test Categories

### Category Breakdown

| Category | Test Count Est. | Coverage Target | Phase 1 Reuse |
|----------|----------------|-----------------|---------------|
| **Unit Tests - Command Layer** | 35 | 90%+ | Mock providers |
| **Unit Tests - Progress Reporter** | 15 | 85%+ | Empty image |
| **Unit Tests - Input Validator** | 45 | 95%+ | Patterns from Phase 1 |
| **Unit Tests - Error Handling** | 20 | 90%+ | Error types |
| **Integration Tests - CLI** | 25 | 80%+ | Docker daemon + test registry |
| **Integration Tests - Component** | 15 | 85%+ | All Phase 1 packages |
| **TOTAL** | **155** | **85%+** | **Heavy reuse** |

---

## 1. Unit Tests - Command Layer (Wave 1, Effort 2.1.1)

### 1.1 Flag Parsing Tests (12 tests)

**Test File**: `cmd/push/flags_test.go`

**Pattern**: Table-driven tests similar to `pkg/docker/client_test.go`

```go
// TC-CMD-FLAG-001 through TC-CMD-FLAG-012
func TestFlagParsing_AllFlags(t *testing.T) {
    testCases := []struct {
        name     string
        args     []string
        expected PushConfig
    }{
        {
            "all flags provided",
            []string{"myapp:latest", "--registry", "https://custom.io", "--username", "user", "--password", "pass", "--insecure", "--verbose"},
            PushConfig{Registry: "https://custom.io", Username: "user", ...},
        },
        {
            "minimal flags",
            []string{"myapp:latest"},
            PushConfig{Registry: "https://gitea.cnoe.localtest.me:8443", ...}, // defaults
        },
        {
            "short flags",
            []string{"myapp:latest", "-r", "https://reg.io", "-u", "admin", "-p", "secret", "-k", "-v"},
            PushConfig{...},
        },
        // ... 9 more cases for each flag combination
    }
}
```

**Covered Scenarios**:
- ✅ All flags provided
- ✅ Minimal flags (defaults applied)
- ✅ Short flag syntax (-r, -u, -p, -k, -v)
- ✅ Long flag syntax (--registry, --username, etc.)
- ✅ Mixed short/long flags
- ✅ Missing required password (error)
- ✅ Invalid flag values
- ✅ Extra positional arguments (error)
- ✅ No image name provided (error)
- ✅ Multiple image names (error)
- ✅ Help flag (--help)
- ✅ Version flag (--version, if applicable)

### 1.2 Environment Variable Tests (8 tests)

**Test File**: `cmd/push/envvars_test.go`

**Pattern**: Similar to flag parsing, testing viper precedence

```go
// TC-CMD-ENV-001 through TC-CMD-ENV-008
func TestEnvironmentVariables(t *testing.T) {
    // Test with os.Setenv() to set IDPBUILDER_* variables
    // Verify viper reads them correctly
}
```

**Covered Scenarios**:
- ✅ `IDPBUILDER_REGISTRY` sets registry
- ✅ `IDPBUILDER_REGISTRY_USERNAME` sets username
- ✅ `IDPBUILDER_REGISTRY_PASSWORD` sets password
- ✅ `IDPBUILDER_INSECURE` sets insecure mode
- ✅ `IDPBUILDER_VERBOSE` sets verbose mode
- ✅ Flag overrides env var (precedence)
- ✅ Env var overrides default
- ✅ Missing env var uses default

### 1.3 Pipeline Orchestration Tests (10 tests)

**Test File**: `cmd/push/pipeline_test.go`

**Reuse**: Mock providers from `pkg/registry/client_test.go`

```go
// TC-CMD-PIPE-001 through TC-CMD-PIPE-010
func TestPushPipeline_Success(t *testing.T) {
    // Mock all dependencies
    mockDocker := &mockDockerClient{...}
    mockRegistry := &mockRegistryClient{...}
    mockAuth := &mockAuthProvider{...}
    mockTLS := &mockTLSProvider{...}

    // Execute pipeline
    err := executePushWorkflow(config, mockDocker, mockRegistry, mockAuth, mockTLS)

    // Verify success
    require.NoError(t, err)
}
```

**Covered Scenarios**:
- ✅ Full pipeline success (all stages pass)
- ✅ Validation stage failure (invalid config)
- ✅ Docker stage failure (image not found)
- ✅ Auth stage failure (credentials invalid)
- ✅ TLS stage success (both secure and insecure)
- ✅ Push stage failure (network error)
- ✅ Progress callback invoked correctly
- ✅ Exit code mapping (0, 1, 2, 3, 4)
- ✅ Error message formatting
- ✅ Context cancellation handling

### 1.4 Exit Code Mapping Tests (5 tests)

**Test File**: `cmd/push/exitcodes_test.go`

**Pattern**: Error type testing like `pkg/docker/client_test.go`

```go
// TC-CMD-EXIT-001 through TC-CMD-EXIT-005
func TestMapErrorToExitCode(t *testing.T) {
    testCases := []struct {
        name     string
        err      error
        expected int
    }{
        {"validation error", &ValidationError{...}, 1},
        {"auth error", &AuthenticationError{...}, 2},
        {"network error", &NetworkError{...}, 3},
        {"image not found", &ImageNotFoundError{...}, 4},
        {"generic error", errors.New("unknown"), 1},
    }
}
```

**Covered Scenarios**:
- ✅ ValidationError → exit 1
- ✅ AuthenticationError → exit 2
- ✅ NetworkError → exit 3
- ✅ ImageNotFoundError → exit 4
- ✅ Generic error → exit 1

---

## 2. Unit Tests - Progress Reporter (Wave 1, Effort 2.1.2)

### 2.1 Progress Update Handling Tests (10 tests)

**Test File**: `pkg/progress/reporter_test.go`

**Reuse**: Empty image from `pkg/registry/client_test.go`

```go
// TC-PROG-UPDATE-001 through TC-PROG-UPDATE-010
func TestProgressReporter_LayerStart(t *testing.T) {
    // Capture console output
    var buf bytes.Buffer
    reporter := NewProgressReporter(&buf, false) // normal mode

    // Send layer start update
    update := ProgressUpdate{
        Status:      "uploading",
        LayerDigest: "sha256:abc123...",
        LayerSize:   12500000,
        LayerIndex:  1,
        TotalLayers: 3,
    }
    reporter.HandleProgress(update)

    // Verify output format
    output := buf.String()
    assert.Contains(t, output, "⏳ Layer 1/3")
    assert.Contains(t, output, "sha256:abc123")
    assert.Contains(t, output, "12.5 MB")
}
```

**Covered Scenarios**:
- ✅ Layer upload started (⏳)
- ✅ Layer upload complete (✓)
- ✅ Layer already exists (✓ skipped)
- ✅ Layer upload progress (percentage in verbose)
- ✅ Manifest push started
- ✅ Manifest push complete
- ✅ Normal mode output (simple status)
- ✅ Verbose mode output (detailed progress)
- ✅ Multiple layers in sequence
- ✅ Final summary formatting

### 2.2 Verbose Mode Tests (5 tests)

**Test File**: `pkg/progress/verbose_test.go`

```go
// TC-PROG-VERB-001 through TC-PROG-VERB-005
func TestProgressReporter_VerboseMode(t *testing.T) {
    var buf bytes.Buffer
    reporter := NewProgressReporter(&buf, true) // verbose mode

    // Send progress update with percentage
    update := ProgressUpdate{
        Status:         "uploading",
        BytesPushed:    23748864,
        TotalBytes:     47497728,
        PercentComplete: 50.0,
    }
    reporter.HandleProgress(update)

    // Verify verbose output includes percentage
    output := buf.String()
    assert.Contains(t, output, "50%")
    assert.Contains(t, output, "23748864/47497728")
}
```

**Covered Scenarios**:
- ✅ Verbose mode shows percentages
- ✅ Verbose mode shows byte counts
- ✅ Verbose mode shows layer skips
- ✅ Normal mode hides detailed progress
- ✅ Mode can be toggled

---

## 3. Unit Tests - Input Validator (Wave 3, Effort 2.3.1)

### 3.1 Image Name Validation Tests (15 tests)

**Test File**: `pkg/validator/imagename_test.go`

**Reuse**: Command injection patterns from `pkg/docker/client_test.go`

```go
// TC-VAL-IMG-001 through TC-VAL-IMG-015
func TestValidateImageName_Valid(t *testing.T) {
    validCases := []string{
        "alpine:latest",
        "myapp:v1.2.3",
        "my-app_image:dev",
        "registry.io/namespace/app:tag",
        "localhost:5000/myapp:latest",
    }

    for _, imageName := range validCases {
        t.Run(imageName, func(t *testing.T) {
            err := ValidateImageName(imageName)
            assert.NoError(t, err)
        })
    }
}

func TestValidateImageName_CommandInjection(t *testing.T) {
    // Reuse patterns from pkg/docker/client_test.go
    dangerousInputs := []string{
        "alpine;rm -rf /",
        "alpine|whoami",
        "alpine&whoami",
        "alpine$USER",
        "alpine`whoami`",
        "alpine()",
        "alpine<file",
        "alpine>file",
        "alpine\\test",
    }

    for _, input := range dangerousInputs {
        t.Run(input, func(t *testing.T) {
            err := ValidateImageName(input)
            assert.Error(t, err)

            var valErr *ValidationError
            assert.ErrorAs(t, err, &valErr)
        })
    }
}
```

**Covered Scenarios**:
- ✅ Valid OCI image names (5 cases)
- ✅ Command injection prevention (9 cases - from Phase 1)
- ✅ Empty image name (error)
- ✅ Maximum length enforcement (256 chars)
- ✅ Path traversal prevention (../)
- ✅ Whitelist character enforcement

### 3.2 Registry URL Validation Tests (15 tests)

**Test File**: `pkg/validator/registryurl_test.go`

**Pattern**: Similar to image name validation

```go
// TC-VAL-REG-001 through TC-VAL-REG-015
func TestValidateRegistryURL_Valid(t *testing.T) {
    validURLs := []string{
        "https://registry.io",
        "https://gitea.cnoe.localtest.me:8443",
        "http://localhost:5000",
        "https://quay.io",
        "https://docker.io",
    }
    // Test each for successful parsing
}

func TestValidateRegistryURL_Invalid(t *testing.T) {
    invalidURLs := []string{
        "://invalid",
        "ftp://registry.io",  // wrong scheme
        "registry.io",        // missing scheme
        "https://",           // missing host
        "https:// spaces",    // spaces in URL
    }
    // Test each for validation error
}
```

**Covered Scenarios**:
- ✅ Valid HTTPS URLs (3 cases)
- ✅ Valid HTTP URLs (2 cases)
- ✅ Invalid URL schemes (FTP, file, etc.)
- ✅ Malformed URLs (5 cases)
- ✅ Missing scheme (error)
- ✅ Missing hostname (error)
- ✅ SSRF prevention (localhost warning)
- ✅ Private IP detection (10.0.0.0/8, etc.)
- ✅ Port validation
- ✅ URL parsing errors

### 3.3 Credential Validation Tests (10 tests)

**Test File**: `pkg/validator/credentials_test.go`

**Reuse**: Special character tests from `pkg/auth/basic_test.go`

```go
// TC-VAL-CRED-001 through TC-VAL-CRED-010
func TestValidateCredentials_PasswordSpecialChars(t *testing.T) {
    // Reuse from pkg/auth/basic_test.go
    validPasswords := []string{
        "P@ss!w0rd#123",
        "пароль密码🔒",
        "pass with spaces",
        "pass\"with\"quotes",
        "pass'with'quotes",
        "P@ss!w0rd#123 with \"quotes\" and 'apostrophes' пароль🔒",
    }

    for _, password := range validPasswords {
        t.Run(password, func(t *testing.T) {
            err := ValidateCredentials("user", password)
            assert.NoError(t, err)
        })
    }
}
```

**Covered Scenarios**:
- ✅ Empty username (error)
- ✅ Empty password (error)
- ✅ Both empty (error)
- ✅ Special characters in password (6 cases - from Phase 1)
- ✅ Unicode in password (valid)
- ✅ Quotes in password (valid)
- ✅ Control characters in username (error)
- ✅ Password length validation
- ✅ Username with spaces (valid or invalid based on spec)

### 3.4 Sanitization Tests (5 tests)

**Test File**: `pkg/validator/sanitize_test.go`

```go
// TC-VAL-SAN-001 through TC-VAL-SAN-005
func TestSanitizeInput(t *testing.T) {
    testCases := []struct {
        name     string
        input    string
        expected string
    }{
        {"no sanitization needed", "alpine:latest", "alpine:latest"},
        {"preserve special chars", "P@ss!w0rd#123", "P@ss!w0rd#123"},
        {"escape shell metacharacters", "test;whoami", "test\\;whoami"},
        // ... more cases
    }
}
```

**Covered Scenarios**:
- ✅ No sanitization needed (clean input)
- ✅ Preserve legitimate special chars
- ✅ Escape shell metacharacters
- ✅ Remove dangerous characters
- ✅ Handle unicode correctly

---

## 4. Unit Tests - Error Handling (Wave 3, Effort 2.3.2)

### 4.1 Error Type Tests (10 tests)

**Test File**: `cmd/push/errors_test.go`

**Reuse**: Error type patterns from `pkg/docker/client_test.go`

```go
// TC-ERR-TYPE-001 through TC-ERR-TYPE-010
func TestErrorTypes_ValidationError(t *testing.T) {
    // Similar to pkg/docker/client_test.go error type tests
    err := &ValidationError{
        Field:   "imageName",
        Message: "cannot be empty",
    }

    assert.Contains(t, err.Error(), "imageName")
    assert.Contains(t, err.Error(), "cannot be empty")
}
```

**Covered Scenarios**:
- ✅ ValidationError creation and formatting
- ✅ AuthenticationError creation and formatting
- ✅ NetworkError creation and formatting
- ✅ ImageNotFoundError creation and formatting
- ✅ PushFailedError creation and formatting
- ✅ Error wrapping with `fmt.Errorf("%w")`
- ✅ Error unwrapping with `errors.Unwrap()`
- ✅ Error type detection with `errors.As()`
- ✅ Error messages include suggestions
- ✅ Error messages include context

### 4.2 Error Message Formatting Tests (10 tests)

**Test File**: `cmd/push/errormessages_test.go`

```go
// TC-ERR-MSG-001 through TC-ERR-MSG-010
func TestFormatError_ImageNotFound(t *testing.T) {
    err := &ImageNotFoundError{ImageName: "myapp:latest"}

    formatted := FormatError(err)

    assert.Contains(t, formatted, "Error:")
    assert.Contains(t, formatted, "myapp:latest")
    assert.Contains(t, formatted, "Suggestion:")
    assert.Contains(t, formatted, "docker images")
}
```

**Covered Scenarios**:
- ✅ Image not found error formatting
- ✅ Authentication failed error formatting
- ✅ Registry unreachable error formatting
- ✅ TLS verification failed error formatting
- ✅ Validation error formatting
- ✅ Network error formatting
- ✅ Generic error formatting
- ✅ All errors include "Error:" and "Suggestion:"
- ✅ Verbose mode shows additional context
- ✅ Error redaction (passwords hidden)

---

## 5. Integration Tests - CLI (All Waves Combined)

### 5.1 End-to-End CLI Tests (10 tests)

**Test File**: `tests/integration/cli_test.go`

**Setup**: Requires Docker daemon + test Gitea instance

**Reuse**: Docker daemon integration from `pkg/docker/client_test.go`

```go
// TC-CLI-E2E-001 through TC-CLI-E2E-010
func TestCLI_PushSuccess(t *testing.T) {
    // Prerequisite: Docker daemon running, alpine:latest pulled
    // Prerequisite: Test Gitea instance running (container)

    // Execute command
    cmd := exec.Command("idpbuilder", "push", "alpine:latest",
        "--registry", testRegistryURL,
        "--username", "giteaadmin",
        "--password", "password123",
        "--insecure",
    )

    output, err := cmd.CombinedOutput()

    // Verify success
    require.NoError(t, err)
    assert.Contains(t, string(output), "Successfully pushed")
    assert.Equal(t, 0, cmd.ProcessState.ExitCode())
}
```

**Test Environment Setup**:
```bash
# Run before integration tests
docker run -d --name test-gitea \
  -p 3000:3000 \
  -p 8443:8443 \
  gitea/gitea:latest

# Wait for Gitea to start
sleep 10

# Configure test credentials
# (Use Gitea API to create test user)
```

**Covered Scenarios**:
- ✅ Push to test Gitea (success)
- ✅ Push to DockerHub (if credentials available)
- ✅ Push with custom registry
- ✅ Push with environment variables
- ✅ Push with verbose mode
- ✅ Push with insecure mode
- ✅ Push failure - image not found
- ✅ Push failure - wrong credentials
- ✅ Push failure - registry unreachable
- ✅ Push failure - TLS verification error

### 5.2 Flag Integration Tests (8 tests)

**Test File**: `tests/integration/flags_test.go`

```go
// TC-CLI-FLAG-001 through TC-CLI-FLAG-008
func TestCLI_AllFlagsIntegration(t *testing.T) {
    // Test all flags together in real execution
    // Verify viper binding works correctly
}
```

**Covered Scenarios**:
- ✅ All flags provided
- ✅ Short flags
- ✅ Long flags
- ✅ Mixed flags and env vars
- ✅ Flag precedence over env vars
- ✅ Env var precedence over defaults
- ✅ Help text display
- ✅ Invalid flag handling

### 5.3 Progress Display Integration Tests (7 tests)

**Test File**: `tests/integration/progress_test.go`

```go
// TC-CLI-PROG-001 through TC-CLI-PROG-007
func TestCLI_ProgressDisplay(t *testing.T) {
    // Execute push and capture output
    // Verify progress updates appear
}
```

**Covered Scenarios**:
- ✅ Normal mode progress (simple status)
- ✅ Verbose mode progress (detailed)
- ✅ Layer skipping shown (already exists)
- ✅ Multiple layers displayed
- ✅ Manifest push shown
- ✅ Final summary displayed
- ✅ Progress updates in real-time

---

## 6. Integration Tests - Component Integration

### 6.1 Command + Phase 1 Integration Tests (10 tests)

**Test File**: `tests/integration/component_test.go`

**Reuse**: All Phase 1 packages (docker, registry, auth, tls)

```go
// TC-COMP-INT-001 through TC-COMP-INT-010
func TestComponentIntegration_FullPipeline(t *testing.T) {
    // Use REAL Phase 1 implementations (not mocks)
    dockerClient, err := docker.NewClient()
    require.NoError(t, err)

    authProvider := auth.NewBasicAuthProvider("giteaadmin", "password")
    tlsProvider := tls.NewConfigProvider(true) // insecure for testing

    registryClient, err := registry.NewClient(authProvider, tlsProvider)
    require.NoError(t, err)

    // Execute push command with real components
    err = executePushWithRealComponents(dockerClient, registryClient, ...)

    // Verify success
    require.NoError(t, err)
}
```

**Covered Scenarios**:
- ✅ Full pipeline with all Phase 1 components
- ✅ Docker client retrieves real image
- ✅ Auth provider validates credentials
- ✅ TLS provider configures transport
- ✅ Registry client pushes successfully
- ✅ Progress callback receives updates
- ✅ Error handling across components
- ✅ Context propagation
- ✅ Resource cleanup
- ✅ Memory usage monitoring

### 6.2 Custom Registry Integration Tests (5 tests)

**Test File**: `tests/integration/registry_override_test.go`

```go
// TC-COMP-REG-001 through TC-COMP-REG-005
func TestRegistryOverride_CustomRegistry(t *testing.T) {
    // Test pushing to different registries
    registries := []string{
        "https://gitea.cnoe.localtest.me:8443",
        "http://localhost:5000",  // local registry
        // Add more if available
    }

    for _, registry := range registries {
        t.Run(registry, func(t *testing.T) {
            // Execute push to custom registry
        })
    }
}
```

**Covered Scenarios**:
- ✅ Push to default Gitea
- ✅ Push to custom registry (http://localhost:5000)
- ✅ Push to DockerHub (if credentials)
- ✅ Push to Quay.io (if credentials)
- ✅ Registry URL override works correctly

---

## Test Execution Strategy

### 7.1 Test Execution Order

**Phase 2 Development Sequence**:

1. **Wave 1, Effort 2.1.1 (Command Core)**:
   - Run unit tests: `cmd/push/*_test.go`
   - Expected: 35 tests, 90%+ coverage
   - Duration: ~30 seconds

2. **Wave 1, Effort 2.1.2 (Progress Reporter)**:
   - Run unit tests: `pkg/progress/*_test.go`
   - Expected: 15 tests, 85%+ coverage
   - Duration: ~10 seconds

3. **Wave 2, Efforts 2.2.1 & 2.2.2 (Parallel)**:
   - Run unit tests for both efforts
   - Expected: Combined coverage maintained
   - Duration: ~20 seconds

4. **Wave 3, Efforts 2.3.1 & 2.3.2**:
   - Run unit tests: `pkg/validator/*_test.go`, `cmd/push/errors*_test.go`
   - Expected: 65 tests, 90%+ coverage
   - Duration: ~30 seconds

5. **Integration Tests (All Waves Complete)**:
   - Run integration tests: `tests/integration/*_test.go`
   - Expected: 40 tests, 80%+ coverage
   - Duration: ~2-3 minutes (requires Docker + registry)

### 7.2 Test Harness Setup

**Prerequisite Setup** (before running integration tests):

```bash
#!/bin/bash
# setup-test-environment.sh

echo "Setting up Phase 2 test environment..."

# 1. Verify Docker daemon is running
if ! docker info > /dev/null 2>&1; then
    echo "ERROR: Docker daemon not running"
    exit 1
fi

# 2. Pull test images
echo "Pulling test images..."
docker pull alpine:latest

# 3. Start test Gitea registry (if not already running)
if ! docker ps | grep -q test-gitea; then
    echo "Starting test Gitea registry..."
    docker run -d --name test-gitea \
      -p 3000:3000 \
      -p 8443:8443 \
      -e USER_UID=1000 \
      -e USER_GID=1000 \
      gitea/gitea:latest

    # Wait for Gitea to start
    echo "Waiting for Gitea to start..."
    sleep 15
fi

# 4. Configure test credentials (using Gitea API or admin user)
export IDPBUILDER_REGISTRY="https://gitea.cnoe.localtest.me:8443"
export IDPBUILDER_REGISTRY_USERNAME="giteaadmin"
export IDPBUILDER_REGISTRY_PASSWORD="password123"
export IDPBUILDER_INSECURE="true"

echo "✅ Test environment ready"
echo "Run tests with: go test ./..."
```

**Cleanup Script**:

```bash
#!/bin/bash
# cleanup-test-environment.sh

echo "Cleaning up test environment..."

# Stop test Gitea
docker stop test-gitea
docker rm test-gitea

# Clean up test images (optional)
# docker rmi alpine:latest

echo "✅ Cleanup complete"
```

### 7.3 Continuous Integration

**CI Pipeline** (GitHub Actions / GitLab CI):

```yaml
# .github/workflows/phase2-tests.yml
name: Phase 2 Tests

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run unit tests
        run: |
          go test -v -coverprofile=coverage.out ./cmd/push/... ./pkg/progress/... ./pkg/validator/...
          go tool cover -html=coverage.out -o coverage.html

      - name: Upload coverage
        uses: actions/upload-artifact@v3
        with:
          name: coverage-report
          path: coverage.html

  integration-tests:
    runs-on: ubuntu-latest
    services:
      docker:
        image: docker:dind
        options: --privileged
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4

      - name: Setup test environment
        run: ./setup-test-environment.sh

      - name: Run integration tests
        run: |
          go test -v -tags=integration ./tests/integration/...

      - name: Cleanup
        if: always()
        run: ./cleanup-test-environment.sh
```

---

## Coverage Targets and Measurement

### 8.1 Coverage Targets by Component

| Component | Target Coverage | Critical Paths | Acceptable Minimum |
|-----------|----------------|----------------|-------------------|
| Command layer (cmd/push) | 90% | Flag parsing, pipeline | 85% |
| Progress reporter (pkg/progress) | 85% | Update handling | 80% |
| Input validator (pkg/validator) | 95% | Command injection, SSRF | 90% |
| Error handling (cmd/push/errors) | 90% | All error types | 85% |
| Integration (overall) | 80% | Full workflow | 75% |

### 8.2 Coverage Measurement Commands

```bash
# Measure coverage for all Phase 2 code
go test -coverprofile=coverage.out ./cmd/push/... ./pkg/progress/... ./pkg/validator/...

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# View coverage by function
go tool cover -func=coverage.out

# Check coverage threshold (85% minimum)
total_coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
if (( $(echo "$total_coverage < 85" | bc -l) )); then
    echo "ERROR: Coverage $total_coverage% is below 85% threshold"
    exit 1
fi
```

---

## Test Data and Fixtures

### 9.1 Test Images

**Required Images**:
- `alpine:latest` - Small, fast, always available (from Phase 1)
- Test images built during integration tests:
  ```dockerfile
  # tests/integration/testdata/Dockerfile.test1
  FROM scratch
  COPY hello.txt /hello.txt
  ```

**Image Sizes for Testing**:
- Small: alpine:latest (~5 MB) - Fast iteration
- Medium: Custom test image (~50 MB) - Progress testing
- Large: Not needed for Phase 2 (performance testing in Phase 3)

### 9.2 Test Credentials

**Test Gitea Credentials**:
- Username: `giteaadmin`
- Password: `password123` (or from env var)
- Registry: `https://gitea.cnoe.localtest.me:8443`

**Invalid Credentials (for error testing)**:
- Username: `wronguser`
- Password: `wrongpassword`

### 9.3 Mock Responses

**Registry HTTP Responses** (using httptest):

```go
// Success response (200 OK)
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if strings.HasSuffix(r.URL.Path, "/v2/") {
        w.WriteHeader(http.StatusOK)
    }
}))

// Auth required response (401 Unauthorized)
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusUnauthorized)
}))

// Registry unavailable (404 Not Found)
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
}))
```

---

## Test Documentation Requirements

### 10.1 Test Case Documentation

**Each test MUST include**:
- Test case ID (TC-XXX-XXX-NNN)
- Given/When/Then structure (from Phase 1 pattern)
- Clear test name describing scenario
- Expected outcome documented

**Example**:
```go
// TC-CMD-FLAG-001: All flags provided
// Given: User provides all command flags
// When: Parsing flags
// Then: All values captured correctly
func TestFlagParsing_AllFlags(t *testing.T) {
    // Test implementation
}
```

### 10.2 Test Failure Reporting

**When tests fail, report MUST include**:
- Test case ID
- Failure reason
- Expected vs actual values
- Stack trace
- Environment details (Go version, OS)

**Example Test Failure Report**:
```
FAIL: TC-CMD-FLAG-005 - Flag precedence over env vars
Reason: Flag value not overriding environment variable
Expected: registry = "https://custom.io" (from flag)
Actual: registry = "https://gitea.cnoe.localtest.me:8443" (from env)
Stack: cmd/push/flags_test.go:123
Environment: Go 1.21.5, Linux x86_64
```

---

## Risk Mitigation

### 11.1 Test Environment Risks

**Risk**: Docker daemon not running during integration tests
**Mitigation**: Skip integration tests with `t.Skip()` if daemon unavailable
**Detection**: `docker info` check in setup script

**Risk**: Test Gitea registry not available
**Mitigation**: Fall back to localhost:5000 registry or skip
**Detection**: HTTP health check before tests

**Risk**: Network isolation in CI environment
**Mitigation**: Use Docker-in-Docker (dind) service
**Detection**: Pre-flight network connectivity check

### 11.2 Test Data Risks

**Risk**: alpine:latest not pulled
**Mitigation**: Pull in setup script, verify presence
**Detection**: `docker images alpine:latest` check

**Risk**: Test credentials expired
**Mitigation**: Use environment variables, document renewal
**Detection**: Auth test before integration suite

### 11.3 Test Coverage Risks

**Risk**: Coverage drops below 85%
**Mitigation**: Coverage gate in CI pipeline, fail build
**Detection**: `go tool cover -func` threshold check

**Risk**: Critical paths not tested
**Mitigation**: Mandatory test checklist per component
**Detection**: Manual review before wave completion

---

## Test Plan Maintenance

### 12.1 Test Plan Updates

**When to update this test plan**:
- ✅ New error types added → Add error type tests
- ✅ New flags added → Add flag parsing tests
- ✅ New validation rules → Add validation tests
- ✅ Architecture changes → Review integration tests
- ✅ New Phase 1 fixtures available → Integrate into tests

### 12.2 Test Code Review Requirements

**All test code MUST be reviewed for**:
- ✅ Follows Phase 1 test patterns
- ✅ Uses existing mock providers
- ✅ Given/When/Then structure
- ✅ Clear test case IDs
- ✅ Adequate assertions
- ✅ Error type checking with `ErrorAs()`
- ✅ No flaky tests (context cancellation handled)
- ✅ Resource cleanup (defer)

---

## Test Plan Compliance Checklist

### R341 TDD Compliance

- ✅ **Test plan created BEFORE implementation**: This document created before Phase 2 code
- ✅ **Tests guide implementation**: All test categories define expected behavior
- ✅ **Coverage targets defined**: 85%+ for Phase 2, 90%+ for critical components
- ✅ **Test-first mindset**: Implementation will be driven by these tests

### R342 Early Integration Compliance

- ✅ **Test plan committed to integration branch**: Will be committed before implementation
- ✅ **Tests available to all waves**: All waves can reference this plan
- ✅ **No implementation code in this plan**: Only test specifications

### Phase 1 Integration (Progressive Planning)

- ✅ **Reuses Phase 1 test patterns**: Given/When/Then, table-driven, error types
- ✅ **Reuses Phase 1 fixtures**: Mock providers, test images, httptest servers
- ✅ **References actual Phase 1 implementations**: docker, registry, auth, tls packages
- ✅ **No hypothetical implementations**: All Phase 1 references are to real, completed code

### R510 Checklist Compliance

- ✅ **All sections have checklists**: Test categories, coverage targets, setup steps
- ✅ **Verification criteria clear**: Each test has expected outcome
- ✅ **No ambiguous requirements**: All test scenarios explicitly defined

---

## Appendix A: Test Case Summary

**Total Test Cases**: 155

**By Category**:
- Unit Tests - Command Layer: 35
- Unit Tests - Progress Reporter: 15
- Unit Tests - Input Validator: 45
- Unit Tests - Error Handling: 20
- Integration Tests - CLI: 25
- Integration Tests - Component: 15

**By Wave**:
- Wave 1 (Command + Progress): 50 tests
- Wave 2 (Registry Override + Env Vars): 0 additional tests (covered in Wave 1)
- Wave 3 (Validation + Errors): 65 tests
- Integration (All Waves): 40 tests

---

## Appendix B: Test File Structure

```
idpbuilder-oci-push/
├── cmd/
│   └── push/
│       ├── flags_test.go           # 12 tests
│       ├── envvars_test.go         # 8 tests
│       ├── pipeline_test.go        # 10 tests
│       ├── exitcodes_test.go       # 5 tests
│       ├── errors_test.go          # 10 tests
│       └── errormessages_test.go   # 10 tests
├── pkg/
│   ├── progress/
│   │   ├── reporter_test.go        # 10 tests
│   │   └── verbose_test.go         # 5 tests
│   └── validator/
│       ├── imagename_test.go       # 15 tests
│       ├── registryurl_test.go     # 15 tests
│       ├── credentials_test.go     # 10 tests
│       └── sanitize_test.go        # 5 tests
└── tests/
    └── integration/
        ├── cli_test.go             # 10 tests
        ├── flags_test.go           # 8 tests
        ├── progress_test.go        # 7 tests
        ├── component_test.go       # 10 tests
        └── registry_override_test.go # 5 tests
```

---

## Document Status

**Status**: ✅ READY FOR ORCHESTRATOR HANDOFF
**Test Planner**: @agent-code-reviewer
**Created**: 2025-10-31
**Test Count**: 155 tests
**Coverage Target**: 85%+
**Phase 1 Integration**: Heavy reuse of test patterns and fixtures

**Next Step**:
- Orchestrator proceeds to Phase 2 Wave 1 implementation
- SW Engineers implement code to pass these tests
- Tests run during implementation (TDD workflow)
- Coverage verified before wave completion

**Compliance Verified**:
- ✅ R341: TDD compliance (tests before implementation)
- ✅ R342: Early integration (test plan committed early)
- ✅ R510: Checklist structure followed
- ✅ Progressive planning: Real Phase 1 implementations referenced

---

**END OF PHASE 2 TEST PLAN**
