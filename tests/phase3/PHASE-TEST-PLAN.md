# Phase 3 Test Plan - Integration Testing & Documentation

**Phase**: Phase 3 of 3
**Created**: 2025-11-04
**Created By**: code-reviewer
**Test Philosophy**: Progressive Test Planning (TDD)
**Status**: Tests Defined BEFORE Implementation

---

## Test Philosophy

This test plan follows **progressive test planning** principles:
- Build on REAL test infrastructure from Phases 1 and 2
- Import actual types and fixtures from completed phases
- Use real implementations, not assumptions
- Define tests FIRST (TDD), implement to make them pass
- Integration-first testing strategy

**Critical Distinction**:
- **NOT writing new code** - verifying existing code works
- **NOT unit testing** - that's complete in Phases 1-2
- **YES integration testing** - components working together
- **YES end-to-end testing** - complete user workflows

---

## Test Infrastructure from Previous Phases

### From Phase 1: Foundation & Interfaces (REAL)

**Available Test Utilities**:
```go
// REAL package: pkg/docker/client_test.go
import (
    "github.com/idpbuilder/pkg/docker"
    "github.com/docker/docker/client"
)

// Mock Docker client helper (already exists)
func createTestDockerClient() (docker.Client, error)
func createMockImage(name string, layers int) (*v1.Image, error)
```

**Available Fixtures**:
```go
// REAL package: pkg/registry/client_test.go
import (
    "github.com/idpbuilder/pkg/registry"
    "github.com/google/go-containerregistry/pkg/v1"
)

// Mock registry helpers (already exist)
func createTestRegistryClient() (registry.Client, error)
func createMockAuthProvider(username, password string) auth.Provider
func createMockTLSProvider(insecure bool) tls.ConfigProvider
```

**Available Types** (from Phase 1 - ACTUAL IMPORTS):
```go
// pkg/docker/interface.go - REAL interface
type Client interface {
    ImageExists(ctx context.Context, imageName string) (bool, error)
    GetImage(ctx context.Context, imageName string) (v1.Image, error)
    ValidateImageName(imageName string) error
    Close() error
}

// pkg/registry/interface.go - REAL interface
type Client interface {
    Push(ctx context.Context, image v1.Image, targetRef string,
         progressCallback ProgressCallback) error
    BuildImageReference(registryURL, imageName string) (string, error)
    ValidateRegistry(ctx context.Context, registryURL string) error
}

// pkg/registry/interface.go - REAL type
type ProgressUpdate struct {
    LayerDigest  string
    LayerSize    int64
    BytesPushed  int64
    Status       string  // "uploading" | "complete" | "exists"
}

// pkg/auth/interface.go - REAL interface
type Provider interface {
    GetAuthenticator() (authn.Authenticator, error)
    ValidateCredentials() error
}

// pkg/tls/interface.go - REAL interface
type ConfigProvider interface {
    GetTLSConfig() *tls.Config
    IsInsecure() bool
}
```

**Test Coverage from Phase 1**:
- Unit tests: 85%+ coverage (Docker, Registry, Auth, TLS packages)
- Mock patterns established
- Error types defined
- Validation patterns proven

### From Phase 2: Core Push Functionality (REAL)

**Available Test Utilities**:
```go
// REAL package: cmd/push_test.go (assumed exists)
import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

// Command test helpers (if they exist)
func executeCommand(root *cobra.Command, args ...string) (output string, err error)
func resetViper()  // Clear environment/flags between tests
```

**Available Mock Implementations**:
```go
// From Phase 2 tests - REAL mocks
type MockProgressReporter struct {
    Updates []registry.ProgressUpdate
    Verbose bool
}

func (m *MockProgressReporter) HandleProgress(update registry.ProgressUpdate) {
    m.Updates = append(m.Updates, update)
}
```

**Integration Points from Phase 2**:
```go
// cmd/push.go - REAL command structure
var pushCmd = &cobra.Command{
    Use:   "push [IMAGE]",
    Short: "Push a Docker image to an OCI registry",
    RunE:  executePush,  // REAL function to test
}

// REAL flags (defined in Phase 2)
pushCmd.Flags().StringP("registry", "r", defaultRegistry, "Target registry URL")
pushCmd.Flags().StringP("username", "u", "giteaadmin", "Registry username")
pushCmd.Flags().StringP("password", "p", "", "Registry password (REQUIRED)")
pushCmd.Flags().BoolP("insecure", "k", false, "Skip TLS certificate verification")
pushCmd.Flags().Bool("verbose", false, "Enable verbose output")
```

**Test Coverage from Phase 2**:
- Unit tests: 85%+ coverage (Command layer, validators, progress reporter)
- Flag parsing tested
- Environment variable precedence tested
- Error handling tested
- Mock integrations working

### What Phase 3 ADDS (NOT Duplicating)

Phase 3 tests focus on:
1. **Integration Tests**: Real Docker + Real Gitea (not mocked)
2. **End-to-End Tests**: Complete user workflows
3. **Error Path Integration**: Real failure scenarios
4. **Performance Benchmarks**: Actual timing measurements
5. **Documentation Tests**: Validate docs match implementation

**Phase 3 does NOT duplicate**:
- ❌ Unit tests for Docker package (done in Phase 1)
- ❌ Unit tests for Registry package (done in Phase 1)
- ❌ Unit tests for Command layer (done in Phase 2)
- ❌ Mock-based tests (done in Phases 1-2)

---

## Phase 3 Test Requirements (From Architecture)

### Test Pyramid Structure (From PHASE-ARCHITECTURE-PLAN.md)

```
TIER 1 - Unit Tests (70% of tests):
  ✅ Already complete from Phase 1 & 2
  ✅ Fast execution (<1 second per test)
  ✅ Mock all external dependencies
  ✅ Coverage: 85%+ per package

TIER 2 - Integration Tests (20% of tests): ← PHASE 3 WAVE 1
  - Test component interactions
  - Use real Docker daemon
  - Use test Gitea registry in container
  - Test network error scenarios
  - Test authentication flows
  - Coverage: Critical paths

TIER 3 - End-to-End Tests (10% of tests): ← PHASE 3 WAVE 1
  - Test complete user workflows
  - Full stack with real dependencies
  - Verify actual image push success
  - Test multi-layer images
  - Test multi-architecture images
  - Coverage: Happy paths and critical errors
```

### Testing Infrastructure: testcontainers-go (REAL Library)

**Actual Dependency** (from PHASE-ARCHITECTURE-PLAN.md):
```
Library: github.com/testcontainers/testcontainers-go v0.26.0+
Purpose: Container lifecycle management for integration tests
```

**Real Usage Pattern**:
```go
import (
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

// REAL helper function (to be implemented in Phase 3)
func SetupGiteaContainer(ctx context.Context) (testcontainers.Container, string, error) {
    req := testcontainers.ContainerRequest{
        Image:        "gitea/gitea:latest",
        ExposedPorts: []string{"3000/tcp", "8443/tcp"},
        Env: map[string]string{
            "GITEA__server__ROOT_URL":            "http://localhost:3000",
            "GITEA__server__HTTP_PORT":           "3000",
            "GITEA__database__DB_TYPE":           "sqlite3",
            "GITEA__security__INSTALL_LOCK":      "true",
            "GITEA__service__DISABLE_REGISTRATION": "false",
        },
        WaitingFor: wait.ForLog("Gitea started"),
    }

    container, err := testcontainers.GenericContainer(ctx,
        testcontainers.GenericContainerRequest{
            ContainerRequest: req,
            Started:          true,
        })
    if err != nil {
        return nil, "", err
    }

    // Get dynamic port for registry
    mappedPort, err := container.MappedPort(ctx, "8443")
    if err != nil {
        return nil, "", err
    }

    registryURL := fmt.Sprintf("localhost:%s", mappedPort.Port())
    return container, registryURL, nil
}
```

---

## Test Harness Design (Using REAL Imports)

### Master Test Harness (test/harness/setup.go)

**Purpose**: Centralize test environment setup with real dependencies

```go
package harness

import (
    "context"
    "testing"

    // REAL imports from Phase 1
    "github.com/idpbuilder/pkg/docker"
    "github.com/idpbuilder/pkg/registry"
    "github.com/idpbuilder/pkg/auth"
    "github.com/idpbuilder/pkg/tls"

    // REAL imports for testing
    "github.com/testcontainers/testcontainers-go"
    "github.com/google/go-containerregistry/pkg/v1"
    dockerclient "github.com/docker/docker/client"
)

// TestEnvironment provides complete test infrastructure
type TestEnvironment struct {
    // Real Docker client (connects to host Docker)
    DockerClient docker.Client

    // Testcontainer for Gitea registry
    GiteaContainer testcontainers.Container
    RegistryURL    string
    Credentials    struct {
        Username string
        Password string
    }

    // Test image built for testing
    TestImage      v1.Image
    TestImageName  string

    // Cleanup function
    Cleanup func()
}

// SetupIntegrationTest creates full integration test environment
func SetupIntegrationTest(t *testing.T) *TestEnvironment {
    ctx := context.Background()
    env := &TestEnvironment{}

    // 1. Setup Gitea container (REAL testcontainer)
    giteaContainer, registryURL, err := SetupGiteaContainer(ctx)
    if err != nil {
        t.Fatalf("Failed to setup Gitea: %v", err)
    }
    env.GiteaContainer = giteaContainer
    env.RegistryURL = registryURL
    env.Credentials.Username = "giteaadmin"
    env.Credentials.Password = "gitea123"

    // 2. Connect to Docker daemon (REAL Docker client)
    dockerClient, err := docker.NewClient()
    if err != nil {
        t.Fatalf("Failed to create Docker client: %v", err)
    }
    env.DockerClient = dockerClient

    // 3. Build test image (REAL image build)
    testImage, testImageName, err := BuildTestImage(ctx, dockerClient)
    if err != nil {
        t.Fatalf("Failed to build test image: %v", err)
    }
    env.TestImage = testImage
    env.TestImageName = testImageName

    // 4. Setup cleanup (REAL cleanup)
    env.Cleanup = func() {
        dockerClient.Close()
        giteaContainer.Terminate(ctx)
    }

    return env
}

// BuildTestImage creates a real test image in Docker daemon
func BuildTestImage(ctx context.Context, dockerClient docker.Client) (v1.Image, string, error) {
    // Build actual Docker image (not mocked)
    // Dockerfile embedded or generated
    imageName := "test-idpbuilder-push:latest"

    // Real docker build execution
    // Returns v1.Image compatible with go-containerregistry
    // Implementation details in test/harness/image_builder.go

    return image, imageName, nil
}
```

### Test Image Builder (test/harness/image_builder.go)

**Purpose**: Build real test images with specific characteristics

```go
package harness

import (
    "context"
    "fmt"
    "io"
    "os"

    "github.com/docker/docker/client"
    "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/tarball"
)

// ImageSpec defines test image characteristics
type ImageSpec struct {
    NumLayers int
    LayerSize int64  // Size in bytes per layer
    MultiArch bool    // If true, build multi-arch manifest
}

// BuildImageFromSpec builds a real Docker image matching spec
func BuildImageFromSpec(ctx context.Context, spec ImageSpec) (v1.Image, string, error) {
    // Generate Dockerfile with specified layers
    dockerfile := generateDockerfile(spec)

    // Build image using Docker daemon
    imageName := fmt.Sprintf("test-image-l%d-s%d:latest", spec.NumLayers, spec.LayerSize)

    // Real Docker build API call
    dockerClient, _ := client.NewClientWithOpts(client.FromEnv)
    buildContext := createBuildContext(dockerfile, spec)

    // Execute build
    resp, err := dockerClient.ImageBuild(ctx, buildContext, types.ImageBuildOptions{
        Tags: []string{imageName},
    })
    if err != nil {
        return nil, "", err
    }
    defer resp.Body.Close()

    // Wait for build to complete
    io.Copy(os.Stdout, resp.Body)

    // Convert to v1.Image (using REAL go-containerregistry)
    image, err := convertDockerImageToV1(imageName, dockerClient)
    if err != nil {
        return nil, "", err
    }

    return image, imageName, nil
}

func generateDockerfile(spec ImageSpec) string {
    // Generate Dockerfile with N layers
    // Each RUN command creates a layer
    dockerfile := "FROM alpine:latest\n"
    for i := 0; i < spec.NumLayers; i++ {
        // Create layer with specified size
        dockerfile += fmt.Sprintf("RUN dd if=/dev/zero of=/layer%d bs=%d count=1\n",
                                   i, spec.LayerSize)
    }
    return dockerfile
}
```

---

## Test Categories

### 1. Integration Tests - Core Workflow (Wave 1, Effort 3.1.1)

**Test File**: `test/integration/core_workflow_test.go`

**Using REAL implementations from Phases 1-2**:

```go
package integration

import (
    "context"
    "testing"

    // REAL imports from our codebase
    "github.com/idpbuilder/pkg/docker"
    "github.com/idpbuilder/pkg/registry"
    "github.com/idpbuilder/pkg/auth"
    "github.com/idpbuilder/pkg/tls"

    // Test harness
    "github.com/idpbuilder/test/harness"
)

func TestPushSmallImageToGitea(t *testing.T) {
    // Setup REAL test environment
    env := harness.SetupIntegrationTest(t)
    defer env.Cleanup()

    // Use REAL implementations (not mocks!)
    authProvider := auth.NewBasicAuthProvider(env.Credentials.Username, env.Credentials.Password)
    tlsProvider := tls.NewConfigProvider(true)  // Insecure mode for test
    registryClient, err := registry.NewClient(authProvider, tlsProvider)
    if err != nil {
        t.Fatalf("Failed to create registry client: %v", err)
    }

    // Build target reference using REAL method
    targetRef, err := registryClient.BuildImageReference(
        "https://"+env.RegistryURL,
        env.TestImageName,
    )
    if err != nil {
        t.Fatalf("Failed to build reference: %v", err)
    }

    // Execute REAL push
    ctx := context.Background()
    err = registryClient.Push(ctx, env.TestImage, targetRef, nil)
    if err != nil {
        t.Errorf("Push failed: %v", err)
    }

    // Verify image exists in Gitea (REAL verification)
    if !verifyImageInGitea(t, env.RegistryURL, targetRef) {
        t.Errorf("Image not found in registry after push")
    }
}

func TestPushWithProgressReporting(t *testing.T) {
    env := harness.SetupIntegrationTest(t)
    defer env.Cleanup()

    // Build multi-layer test image (REAL image with 5 layers)
    multiLayerImage, imageName, err := harness.BuildImageFromSpec(context.Background(),
        harness.ImageSpec{
            NumLayers: 5,
            LayerSize: 10 * 1024 * 1024,  // 10MB per layer
        })
    if err != nil {
        t.Fatalf("Failed to build multi-layer image: %v", err)
    }

    // Capture progress updates (REAL callback)
    var progressUpdates []registry.ProgressUpdate
    progressCallback := func(update registry.ProgressUpdate) {
        progressUpdates = append(progressUpdates, update)
        t.Logf("Progress: Layer %s - %s (%d/%d bytes)",
               update.LayerDigest[:12], update.Status,
               update.BytesPushed, update.LayerSize)
    }

    // Push with progress monitoring
    authProvider := auth.NewBasicAuthProvider(env.Credentials.Username, env.Credentials.Password)
    tlsProvider := tls.NewConfigProvider(true)
    registryClient, _ := registry.NewClient(authProvider, tlsProvider)

    targetRef, _ := registryClient.BuildImageReference("https://"+env.RegistryURL, imageName)

    err = registryClient.Push(context.Background(), multiLayerImage, targetRef, progressCallback)
    if err != nil {
        t.Errorf("Push failed: %v", err)
    }

    // Verify progress updates received for all 5 layers
    if len(progressUpdates) == 0 {
        t.Errorf("No progress updates received")
    }

    // Count "complete" status updates (should be 5 for 5 layers)
    completeCount := 0
    for _, update := range progressUpdates {
        if update.Status == "complete" {
            completeCount++
        }
    }
    if completeCount < 5 {
        t.Errorf("Expected 5 complete updates, got %d", completeCount)
    }
}
```

**Additional Core Workflow Tests** (using same pattern):
- `TestPushLargeImage` - 100MB+ image
- `TestPushWithAuthenticationSuccess` - Valid credentials
- `TestPushWithCustomRegistry` - Override default registry
- `TestPushMultipleImages` - Sequential pushes

**Estimated**: ~500 lines total for core workflow integration tests

### 2. Integration Tests - Error Paths (Wave 1, Effort 3.1.2)

**Test File**: `test/integration/error_paths_test.go`

**Testing REAL error handling from Phases 1-2**:

```go
package integration

import (
    "context"
    "testing"

    "github.com/idpbuilder/pkg/docker"
    "github.com/idpbuilder/pkg/registry"
    "github.com/idpbuilder/pkg/auth"
    "github.com/idpbuilder/test/harness"
)

func TestPushNonExistentImage(t *testing.T) {
    env := harness.SetupIntegrationTest(t)
    defer env.Cleanup()

    // Try to get non-existent image (REAL error path)
    ctx := context.Background()
    _, err := env.DockerClient.GetImage(ctx, "does-not-exist:latest")

    // Verify REAL error type (from Phase 1)
    if err == nil {
        t.Errorf("Expected ImageNotFoundError, got nil")
    }

    // Type assert to specific error (REAL error type from Phase 1)
    if _, ok := err.(*docker.ImageNotFoundError); !ok {
        t.Errorf("Expected ImageNotFoundError, got %T: %v", err, err)
    }
}

func TestPushWithInvalidCredentials(t *testing.T) {
    env := harness.SetupIntegrationTest(t)
    defer env.Cleanup()

    // Use WRONG password (REAL authentication failure)
    authProvider := auth.NewBasicAuthProvider("giteaadmin", "wrongpassword")
    tlsProvider := tls.NewConfigProvider(true)
    registryClient, _ := registry.NewClient(authProvider, tlsProvider)

    targetRef, _ := registryClient.BuildImageReference("https://"+env.RegistryURL, env.TestImageName)

    // Expect authentication error (REAL error from Phase 1)
    err := registryClient.Push(context.Background(), env.TestImage, targetRef, nil)

    if err == nil {
        t.Errorf("Expected AuthenticationError, got nil")
    }

    // Verify error type
    if _, ok := err.(*registry.AuthenticationError); !ok {
        t.Errorf("Expected AuthenticationError, got %T: %v", err, err)
    }
}

func TestPushWithTLSVerificationFailure(t *testing.T) {
    env := harness.SetupIntegrationTest(t)
    defer env.Cleanup()

    // Use SECURE mode (will fail with self-signed cert)
    authProvider := auth.NewBasicAuthProvider(env.Credentials.Username, env.Credentials.Password)
    tlsProvider := tls.NewConfigProvider(false)  // Secure mode (NO --insecure)
    registryClient, _ := registry.NewClient(authProvider, tlsProvider)

    targetRef, _ := registryClient.BuildImageReference("https://"+env.RegistryURL, env.TestImageName)

    // Expect TLS error (REAL error)
    err := registryClient.Push(context.Background(), env.TestImage, targetRef, nil)

    if err == nil {
        t.Errorf("Expected TLS error with secure mode and self-signed cert")
    }

    // Error should mention certificate verification
    if err != nil && !strings.Contains(err.Error(), "certificate") {
        t.Errorf("Expected certificate error, got: %v", err)
    }
}

func TestPushToUnreachableRegistry(t *testing.T) {
    // Don't setup test environment (no registry running)

    dockerClient, _ := docker.NewClient()
    defer dockerClient.Close()

    // Build test image
    testImage, testImageName, _ := harness.BuildTestImage(context.Background(), dockerClient)

    // Try to push to non-existent registry
    authProvider := auth.NewBasicAuthProvider("user", "pass")
    tlsProvider := tls.NewConfigProvider(true)
    registryClient, _ := registry.NewClient(authProvider, tlsProvider)

    targetRef, _ := registryClient.BuildImageReference("https://unreachable-registry.invalid:9999", testImageName)

    // Expect network error (REAL error)
    err := registryClient.Push(context.Background(), testImage, targetRef, nil)

    if err == nil {
        t.Errorf("Expected NetworkError, got nil")
    }

    if _, ok := err.(*registry.NetworkError); !ok {
        t.Errorf("Expected NetworkError, got %T: %v", err, err)
    }
}
```

**Additional Error Path Tests**:
- `TestPushWithNetworkInterruption` - Simulate connection drop
- `TestPushWithInvalidImageName` - Validation errors
- `TestPushWithInvalidRegistryURL` - URL validation
- `TestPushWithLayerUploadFailure` - Partial upload cleanup

**Estimated**: ~400 lines total for error path tests

### 3. End-to-End Tests - User Workflows (Wave 1, Effort 3.1.1)

**Test File**: `test/e2e/user_workflows_test.go`

**Testing COMPLETE command execution** (cmd/push.go from Phase 2):

```go
package e2e

import (
    "os"
    "os/exec"
    "testing"

    "github.com/idpbuilder/test/harness"
)

func TestCompleteUserWorkflow_BuildAndPush(t *testing.T) {
    env := harness.SetupIntegrationTest(t)
    defer env.Cleanup()

    // Step 1: Build image (REAL docker build command)
    buildCmd := exec.Command("docker", "build", "-t", "myapp:latest", "test/fixtures/sample-app")
    if err := buildCmd.Run(); err != nil {
        t.Fatalf("Docker build failed: %v", err)
    }

    // Step 2: Push using REAL idpbuilder command
    pushCmd := exec.Command("./idpbuilder", "push", "myapp:latest",
        "--registry", "https://"+env.RegistryURL,
        "--password", env.Credentials.Password,
        "--insecure")

    output, err := pushCmd.CombinedOutput()
    if err != nil {
        t.Errorf("Push command failed: %v\nOutput: %s", err, output)
    }

    // Verify output contains success message
    if !strings.Contains(string(output), "Successfully pushed") {
        t.Errorf("Expected success message, got: %s", output)
    }

    // Step 3: Verify image in registry
    if !harness.VerifyImageInGitea(t, env.RegistryURL, "giteaadmin/myapp:latest") {
        t.Errorf("Image not found in registry after push")
    }
}

func TestEnvironmentVariableWorkflow(t *testing.T) {
    env := harness.SetupIntegrationTest(t)
    defer env.Cleanup()

    // Setup environment variables (REAL env var precedence from Phase 2)
    os.Setenv("IDPBUILDER_REGISTRY", "https://"+env.RegistryURL)
    os.Setenv("IDPBUILDER_REGISTRY_PASSWORD", env.Credentials.Password)
    os.Setenv("IDPBUILDER_INSECURE", "true")
    defer func() {
        os.Unsetenv("IDPBUILDER_REGISTRY")
        os.Unsetenv("IDPBUILDER_REGISTRY_PASSWORD")
        os.Unsetenv("IDPBUILDER_INSECURE")
    }()

    // Push WITHOUT flags (using env vars)
    pushCmd := exec.Command("./idpbuilder", "push", env.TestImageName)

    output, err := pushCmd.CombinedOutput()
    if err != nil {
        t.Errorf("Push with env vars failed: %v\nOutput: %s", err, output)
    }
}

func TestVerboseModeOutput(t *testing.T) {
    env := harness.SetupIntegrationTest(t)
    defer env.Cleanup()

    // Push with --verbose flag
    pushCmd := exec.Command("./idpbuilder", "push", env.TestImageName,
        "--registry", "https://"+env.RegistryURL,
        "--password", env.Credentials.Password,
        "--insecure",
        "--verbose")

    output, err := pushCmd.CombinedOutput()
    if err != nil {
        t.Errorf("Verbose push failed: %v", err)
    }

    // Verify verbose output contains detailed logs
    outputStr := string(output)
    if !strings.Contains(outputStr, "[INFO]") {
        t.Errorf("Verbose output missing [INFO] tags")
    }
    if !strings.Contains(outputStr, "[PROGRESS]") {
        t.Errorf("Verbose output missing [PROGRESS] tags")
    }
    if !strings.Contains(outputStr, "Layer") {
        t.Errorf("Verbose output missing layer information")
    }
}
```

**Additional E2E Tests**:
- `TestMultiArchImagePush` - Multi-architecture support
- `TestCICDWorkflow` - Simulated CI/CD pipeline
- `TestFlagPrecedenceOverEnvironment` - CLI flags override env vars

**Estimated**: ~300 lines total for E2E tests

### 4. Performance Benchmarks (Wave 1, Effort 3.1.1)

**Test File**: `test/benchmark/push_performance_test.go`

**Using Go's REAL benchmark framework**:

```go
package benchmark

import (
    "context"
    "testing"

    "github.com/idpbuilder/test/harness"
)

func BenchmarkPushSmallImage_5MB(b *testing.B) {
    // Setup test environment ONCE (outside benchmark loop)
    env := harness.SetupIntegrationTest(&testing.T{})
    defer env.Cleanup()

    // Build small test image (5MB, 2 layers)
    smallImage, imageName, _ := harness.BuildImageFromSpec(context.Background(),
        harness.ImageSpec{
            NumLayers: 2,
            LayerSize: 2.5 * 1024 * 1024,
        })

    // Setup push components
    authProvider := auth.NewBasicAuthProvider(env.Credentials.Username, env.Credentials.Password)
    tlsProvider := tls.NewConfigProvider(true)
    registryClient, _ := registry.NewClient(authProvider, tlsProvider)
    targetRef, _ := registryClient.BuildImageReference("https://"+env.RegistryURL, imageName)

    // Reset timer (exclude setup time)
    b.ResetTimer()

    // Benchmark loop
    for i := 0; i < b.N; i++ {
        // Push image (REAL push operation)
        err := registryClient.Push(context.Background(), smallImage, targetRef, nil)
        if err != nil {
            b.Fatalf("Push failed: %v", err)
        }
    }

    // Report memory allocations
    b.ReportAllocs()
}

func BenchmarkPushMediumImage_100MB(b *testing.B) {
    env := harness.SetupIntegrationTest(&testing.T{})
    defer env.Cleanup()

    // Build medium test image (100MB, 5 layers)
    mediumImage, imageName, _ := harness.BuildImageFromSpec(context.Background(),
        harness.ImageSpec{
            NumLayers: 5,
            LayerSize: 20 * 1024 * 1024,  // 20MB per layer
        })

    authProvider := auth.NewBasicAuthProvider(env.Credentials.Username, env.Credentials.Password)
    tlsProvider := tls.NewConfigProvider(true)
    registryClient, _ := registry.NewClient(authProvider, tlsProvider)
    targetRef, _ := registryClient.BuildImageReference("https://"+env.RegistryURL, imageName)

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        err := registryClient.Push(context.Background(), mediumImage, targetRef, nil)
        if err != nil {
            b.Fatalf("Push failed: %v", err)
        }
    }

    b.ReportAllocs()
}

// Run benchmarks:
// go test -bench=. -benchmem ./test/benchmark/...
```

**Additional Benchmarks**:
- `BenchmarkPushLargeImage_500MB` - Large image performance
- `BenchmarkMemoryFootprint` - Memory usage measurement
- `BenchmarkConcurrentPushes` - Parallel push operations

**Estimated**: ~200 lines total for performance benchmarks

### 5. Documentation Tests (Wave 2, Effort 3.2.1)

**Test File**: `test/docs/documentation_test.go`

**Verify documentation matches REAL implementation**:

```go
package docs

import (
    "os"
    "os/exec"
    "strings"
    "testing"
)

func TestHelpTextCompleteness(t *testing.T) {
    // Get REAL help output from command
    cmd := exec.Command("./idpbuilder", "push", "--help")
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Failed to get help text: %v", err)
    }

    helpText := string(output)

    // Verify all flags documented
    requiredFlags := []string{
        "--registry",
        "--username",
        "--password",
        "--insecure",
        "--verbose",
    }

    for _, flag := range requiredFlags {
        if !strings.Contains(helpText, flag) {
            t.Errorf("Help text missing flag: %s", flag)
        }
    }

    // Verify environment variables mentioned
    requiredEnvVars := []string{
        "IDPBUILDER_REGISTRY",
        "IDPBUILDER_REGISTRY_USERNAME",
        "IDPBUILDER_REGISTRY_PASSWORD",
    }

    for _, envVar := range requiredEnvVars {
        if !strings.Contains(helpText, envVar) {
            t.Errorf("Help text missing env var: %s", envVar)
        }
    }
}

func TestDocumentationExamplesWork(t *testing.T) {
    // Extract examples from docs/push-command.md
    docFile := "docs/push-command.md"
    content, err := os.ReadFile(docFile)
    if err != nil {
        t.Fatalf("Failed to read documentation: %v", err)
    }

    // Parse code blocks from markdown
    examples := extractCodeBlocks(string(content))

    // Verify each example is valid
    for i, example := range examples {
        // Skip non-command examples
        if !strings.HasPrefix(example, "idpbuilder push") {
            continue
        }

        // Verify command syntax is valid (without executing)
        if !isValidCommandSyntax(example) {
            t.Errorf("Example %d has invalid syntax: %s", i, example)
        }
    }
}

func TestMarkdownLinksValid(t *testing.T) {
    // Validate all links in documentation
    docFiles := []string{
        "docs/push-command.md",
        "docs/troubleshooting.md",
        "docs/cicd-examples.md",
        "docs/faq.md",
    }

    for _, docFile := range docFiles {
        content, err := os.ReadFile(docFile)
        if err != nil {
            t.Errorf("Cannot read %s: %v", docFile, err)
            continue
        }

        // Extract and validate links
        links := extractMarkdownLinks(string(content))
        for _, link := range links {
            if !isValidLink(link) {
                t.Errorf("Broken link in %s: %s", docFile, link)
            }
        }
    }
}
```

**Additional Documentation Tests**:
- `TestAPIDocumentationAccuracy` - godoc matches implementation
- `TestTroubleshootingGuideCompleteness` - All error codes documented
- `TestCICDExamplesValid` - CI/CD examples are syntactically correct

**Estimated**: ~300 lines total for documentation tests

---

## Test Execution Strategy

### Local Development Testing

```bash
# Run unit tests (from Phases 1-2) - FAST
make test
# Output: All packages tested, 85%+ coverage, <30 seconds

# Run integration tests (Phase 3) - REQUIRES DOCKER
make test-integration
# Output: Integration tests with Gitea container, ~2 minutes

# Run E2E tests (Phase 3) - REQUIRES BUILT BINARY
make test-e2e
# Output: Full workflow tests, ~3 minutes

# Run all tests
make test-all
# Output: Complete test suite, ~5 minutes

# Run benchmarks
make benchmark
# Output: Performance metrics
```

### Continuous Integration Testing

```yaml
# .github/workflows/test.yml
name: Test Suite

on: [pull_request, push]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: make test

  integration-tests:
    runs-on: ubuntu-latest
    services:
      docker:
        image: docker:dind
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: make test-integration

  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: make build
      - run: make test-e2e

  coverage:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: make test-all
      - run: go test -coverprofile=coverage.out ./...
      - uses: codecov/codecov-action@v3
```

---

## Success Criteria

### Test Completeness Checklist

- [x] Unit tests complete from Phases 1-2 (85%+ coverage)
- [ ] Integration tests cover all critical workflows
- [ ] Error paths comprehensively tested
- [ ] E2E tests validate complete user journeys
- [ ] Performance benchmarks establish baselines
- [ ] Documentation tests ensure accuracy

### Test Quality Criteria

- [ ] All tests use REAL implementations (not new mocks)
- [ ] Test fixtures reuse Phase 1-2 test infrastructure
- [ ] Test environment cleanup is reliable
- [ ] Tests are reproducible and deterministic
- [ ] Test output is clear and actionable

### Coverage Targets

- **Overall Coverage**: 85%+ (maintained from Phases 1-2)
- **Integration Coverage**: All critical paths tested
- **Error Coverage**: All error types triggered and verified
- **Documentation Coverage**: All user-facing docs validated

---

## Test Artifacts (R007 Compliance)

**Note**: Test files are NOT counted in 800-line limit per R007

### Test Directory Structure

```
test/
├── harness/                    # Test infrastructure (not counted)
│   ├── setup.go               # Test environment setup
│   ├── image_builder.go       # Test image builders
│   └── verification.go        # Verification helpers
├── integration/               # Integration tests (not counted)
│   ├── core_workflow_test.go
│   └── error_paths_test.go
├── e2e/                       # E2E tests (not counted)
│   └── user_workflows_test.go
├── benchmark/                 # Benchmarks (not counted)
│   └── push_performance_test.go
└── docs/                      # Doc tests (not counted)
    └── documentation_test.go

Total test lines: ~2000 lines (EXCLUDED from size limit per R007)
```

---

## Dependencies Required (Real)

```go
// go.mod additions for Phase 3
require (
    github.com/testcontainers/testcontainers-go v0.26.0  // Integration test infrastructure
    github.com/docker/docker v24.0.0+incompatible        // Already in Phase 1
    github.com/google/go-containerregistry v0.19.0      // Already in Phase 1
    github.com/spf13/cobra v1.x                          // Already in Phase 2
)
```

---

## Orchestrator Instructions

**Test Plan Created**: planning/phase3/PHASE-TEST-PLAN.md

**Next Actions**:
1. Orchestrator updates orchestrator-state-v3.json:
   - Add to `test_plans.phase.phase3` section
   - Record test_plan_path, test_dir, created_at
2. Proceed to Wave 1 implementation (integration tests)
3. Tests should FAIL initially (TDD red phase)
4. Implementation makes tests pass (TDD green phase)

**Test Plan Metadata for State File**:
```json
{
  "test_plan_path": "planning/phase3/PHASE-TEST-PLAN.md",
  "test_harness_path": "test/harness/setup.go",
  "test_dir": "test/",
  "created_at": "2025-11-04T07:05:16Z",
  "created_by": "code-reviewer",
  "test_count": 18,
  "tdd_phase": "red"
}
```

---

**Test Plan Complete**
**Created By**: code-reviewer (PHASE_TEST_PLANNING state)
**Date**: 2025-11-04
**Compliance**: R341 (TDD), R340 (Progressive Realism), R007 (Test Exclusion)
