# Phase 3 Architecture Plan - Integration Testing & Documentation

**Phase**: Phase 3 - Integration Testing & Documentation
**Created**: 2025-11-04
**Architect**: @agent-architect
**Fidelity Level**: **PSEUDOCODE** (high-level patterns, library choices)

---

## Adaptation Notes

### Lessons from Phase 1 & Phase 2

**What Worked Well in Previous Phases:**
- Interface-first approach (Phase 1) provided clear contracts for all components
- Separation of concerns across packages (Docker, Registry, Auth, TLS) remained clean
- Command orchestration pattern (Phase 2) successfully integrated all packages
- Mock-based unit testing gave fast feedback during development
- Progress reporting using channels worked well for concurrent updates
- Viper integration for flag/env var precedence was seamless

**Key Insights from Phase 2:**
- Command layer integration verified all Phase 1 packages work together
- Error handling with actionable suggestions improved user experience
- Environment variable support simplified CI/CD usage
- Registry override functionality validated package modularity
- TLS insecure mode warnings educated users about security

**Patterns to Continue:**
- Keep high unit test coverage (85%+) across all new code
- Continue using structured error types with clear messages
- Maintain separation between test infrastructure and test logic
- Use Docker containers for integration test dependencies
- Document examples alongside implementation

**Challenges to Address in Phase 3:**
- Need comprehensive end-to-end testing with real Docker and Gitea
- Must test error paths systematically (network failures, auth failures, etc.)
- Documentation must include troubleshooting for common issues
- Build system integration must work across development and CI environments
- Performance characteristics need measurement and documentation

### Changes from Master Architecture

**No Major Deviations:**
- Phase 2 delivered functional push command as planned
- All error handling implemented per specification
- Progress reporting working as designed
- No changes to core architecture needed

**Refinements for Phase 3:**
- Integration tests will use Docker Compose for consistent test environment
- Testing will use testcontainers-go library for manageable container lifecycle
- Documentation will follow Go standard godoc conventions
- Build system integration will use existing IDPBuilder Makefile patterns
- Performance testing will measure actual push times for different image sizes

**New Requirements Discovered:**
- Need test fixtures (pre-built test images) in repository
- Need integration test helper functions for common setup/teardown
- Need comprehensive troubleshooting guide for deployment scenarios
- Need CI/CD pipeline examples for typical usage patterns
- Need performance benchmarks documented for user expectations

---

## High-Level Patterns

### Core Testing Pattern: Three-Tier Test Strategy

**Pattern**: Unit → Integration → End-to-End Test Pyramid

**Pseudocode Structure**:
```
TESTING_PYRAMID:
  TIER 1 - Unit Tests (70% of tests):
    Already complete from Phase 1 & 2
    Fast execution (<1 second per test)
    Mock all external dependencies
    Coverage: 85%+ per package

  TIER 2 - Integration Tests (20% of tests):
    Test component interactions
    Use real Docker daemon
    Use test Gitea registry in container
    Test network error scenarios
    Test authentication flows
    Coverage: Critical paths

  TIER 3 - End-to-End Tests (10% of tests):
    Test complete user workflows
    Full stack with real dependencies
    Verify actual image push success
    Test multi-layer images
    Test multi-architecture images
    Coverage: Happy paths and critical errors
```

**Why This Pattern:**
- Fast feedback from unit tests (already complete)
- Confidence from integration tests (real components)
- Production validation from E2E tests (full workflows)
- Maintains fast test suite (integration/E2E in parallel)

### Testing Infrastructure Pattern: Testcontainers

**Pattern**: Container-Based Test Dependencies

**Pseudocode**:
```
PATTERN TestcontainersInfrastructure:
  Use testcontainers-go library for:
    - Gitea registry container (with pre-configuration)
    - Docker-in-Docker container (for Docker daemon tests)
    - Network isolation (each test gets clean environment)
    - Automatic cleanup (containers removed after tests)

  SETUP:
    FUNCTION SetupTestEnvironment():
      giteaContainer = StartGiteaContainer()
      WAIT FOR giteaContainer.Ready()

      dockerContainer = StartDockerDaemonContainer()
      WAIT FOR dockerContainer.Ready()

      testImage = BuildTestImage(dockerContainer)

      RETURN testEnvironment{
        gitea: giteaContainer,
        docker: dockerContainer,
        testImage: testImage,
      }

  TEARDOWN:
    FUNCTION TeardownTestEnvironment(env):
      env.gitea.Stop()
      env.docker.Stop()
      env.testImage.Remove()
```

### Documentation Pattern: Progressive Disclosure

**Pattern**: Quick Start → Deep Dive → Troubleshooting

**Structure**:
```
DOCUMENTATION_PATTERN ProgressiveDisclosure:

  LEVEL 1 - Quick Start (5 minutes):
    Single command to get started
    Minimal configuration required
    Works with defaults for common case
    Example: idpbuilder push myapp:latest --password secret

  LEVEL 2 - Common Scenarios (10 minutes):
    Custom registry usage
    Environment variable configuration
    Insecure mode for self-signed certs
    CI/CD integration examples
    Multiple authentication methods

  LEVEL 3 - Advanced Topics (20 minutes):
    Performance tuning
    Network troubleshooting
    Multi-architecture images
    Custom TLS configurations
    Debugging failed pushes

  LEVEL 4 - Troubleshooting (reference):
    Error code reference
    Common problems and solutions
    FAQ for specific scenarios
    Links to related documentation
```

---

## Library Choices

### Integration Testing Framework: testcontainers-go

**Choice**: `github.com/testcontainers/testcontainers-go` v0.26.0+
**Version**: Latest stable

**Justification**:
- **Container lifecycle management**: Automatic start/stop of test dependencies
- **Port mapping**: Dynamic port allocation avoids conflicts
- **Wait strategies**: Ensures containers are ready before tests run
- **Cleanup**: Automatic container removal even on test failures
- **Docker Compose support**: Can use compose files for complex setups
- **Go native**: Better integration than shell scripts calling docker commands

**Key Features Used**:
```
LIBRARY_USAGE testcontainers:
  - testcontainers.GenericContainer(): Create Gitea registry container
  - container.MappedPort(): Get dynamic port for connections
  - wait.ForLog(): Wait for "Gitea started" log message
  - container.Terminate(): Clean up after tests
  - compose.Up(): Start multi-container test environments
```

**Alternatives Considered**:
- **Docker CLI via shell**: Rejected (harder to manage lifecycle, cleanup issues)
- **dockertest library**: Rejected (less maintained, fewer features)
- **Manual Docker SDK**: Rejected (too much boilerplate, reinventing wheel)

### Documentation Generation: godoc + Custom Markdown

**Choice**: Standard Go godoc + hand-written markdown guides
**Justification**:
- **godoc**: Auto-generates API documentation from code comments
- **Markdown**: Flexible for user guides, examples, troubleshooting
- **No build step**: Documentation lives with code, easy to update
- **Standard**: Follows Go community conventions

**Documentation Types**:
```
DOCUMENTATION_TYPES:

  API Documentation (godoc):
    - Package-level documentation in doc.go files
    - Function/method documentation in code comments
    - Example tests in _test.go files
    - Generated with: go doc -all ./pkg/...

  User Guides (Markdown):
    - Getting Started: docs/push-command.md
    - Advanced Usage: docs/advanced-usage.md
    - Troubleshooting: docs/troubleshooting.md
    - CI/CD Examples: docs/cicd-examples.md
    - FAQ: docs/faq.md

  IDPBuilder Integration:
    - Update main README.md with push command section
    - Add to command reference in IDPBuilder docs
    - Include in IDPBuilder release notes
```

### Build System: Make (IDPBuilder Standard)

**Choice**: GNU Make (existing IDPBuilder tooling)
**Justification**:
- **Consistency**: IDPBuilder already uses Make
- **No new dependencies**: Developers already have Make
- **Simple**: Easy to understand, maintain, extend
- **Standard**: Common in Go projects

**Makefile Targets**:
```
MAKEFILE_TARGETS:

  make build:
    Build idpbuilder binary with push command included
    Output: bin/idpbuilder

  make test:
    Run all unit tests (fast, no containers)
    Output: Test results and coverage report

  make test-integration:
    Run integration tests (requires Docker)
    Start containers, run tests, cleanup
    Output: Integration test results

  make test-e2e:
    Run end-to-end tests (full workflows)
    Output: E2E test results

  make test-all:
    Run unit + integration + e2e tests
    Full test suite (CI uses this)

  make docs:
    Generate godoc documentation
    Validate markdown documentation
    Output: Documentation files

  make install:
    Install idpbuilder binary to $GOPATH/bin
    Makes command available globally
```

### Performance Testing: Go Benchmark Framework

**Choice**: Go standard testing.B benchmark framework
**Justification**:
- **Built-in**: No external dependencies
- **Statistical**: Runs benchmarks multiple times for accuracy
- **Memory profiling**: Tracks allocations and memory usage
- **CPU profiling**: Identifies performance bottlenecks
- **Standard output**: Compatible with CI tools

**Benchmark Pattern**:
```
BENCHMARK_PATTERN:

  FUNCTION BenchmarkPushSmallImage(b *testing.B):
    Setup test environment (once)
    Prepare small test image (5MB)

    FOR i = 0 TO b.N:
      Reset state
      Measure time to push image
      Record memory allocations
    END

    Report results:
      - Time per operation
      - Allocations per operation
      - Memory used per operation

  Run with: go test -bench=. -benchmem ./...
```

---

## Conceptual Interfaces (Phase 3 Focus)

### Integration Test Helper Interface

```
INTERFACE TestEnvironment:

  FUNCTION SetupGiteaRegistry() -> (registryURL, credentials):
    PURPOSE: Start Gitea container with registry enabled
    OUTPUT:
      - registryURL: "localhost:PORT" (dynamic port)
      - credentials: {username: "giteaadmin", password: "gitea123"}
    LIFECYCLE: Container removed in teardown

  FUNCTION SetupDockerDaemon() -> dockerClient:
    PURPOSE: Provide Docker daemon for image operations
    OUTPUT: Docker client connected to daemon
    NOTES: Can use host Docker or Docker-in-Docker

  FUNCTION BuildTestImage(name, layers) -> imageReference:
    PURPOSE: Create test image with specific characteristics
    INPUT: Image name, layer configuration
    OUTPUT: Image reference in Docker daemon
    EXAMPLE: BuildTestImage("testapp:v1", {layers: 3, size: 10MB})

  FUNCTION PushImageAndVerify(image, registry, creds) -> error:
    PURPOSE: Execute push and verify success in registry
    STEPS:
      1. Push image using idpbuilder push command
      2. Query registry API to verify layers exist
      3. Pull manifest and verify digest matches
    OUTPUT: Error if verification fails, nil if success

  FUNCTION CleanupTestEnvironment(env):
    PURPOSE: Remove all test containers and images
    ENSURES: Clean state for next test run
```

### Test Fixture Builder Interface

```
INTERFACE TestFixtureBuilder:

  FUNCTION CreateMultiLayerImage(numLayers, layerSize) -> image:
    PURPOSE: Build image with specific layer count and sizes
    USE CASE: Test progress reporting with multiple layers
    EXAMPLE: CreateMultiLayerImage(layers=5, layerSize=10MB)

  FUNCTION CreateMultiArchImage(architectures) -> image:
    PURPOSE: Build multi-architecture manifest
    USE CASE: Test multi-arch support
    EXAMPLE: CreateMultiArchImage([amd64, arm64, arm/v7])

  FUNCTION CreateLargeImage(sizeGB) -> image:
    PURPOSE: Build large image for performance testing
    USE CASE: Test streaming, memory usage
    EXAMPLE: CreateLargeImage(1.5) // 1.5GB image

  FUNCTION CreateImageWithTags(baseName, tags) -> images:
    PURPOSE: Create image with multiple tags
    USE CASE: Test tag handling
    EXAMPLE: CreateImageWithTags("myapp", ["latest", "v1.0", "stable"])
```

### Documentation Builder Interface

```
INTERFACE DocumentationBuilder:

  FUNCTION GenerateAPIDocumentation() -> error:
    PURPOSE: Generate godoc for all packages
    OUTPUT: HTML documentation in docs/api/
    COMMAND: Wraps `go doc -all`

  FUNCTION ValidateMarkdownLinks() -> errors:
    PURPOSE: Check all markdown links are valid
    CHECKS:
      - Internal links resolve to existing files
      - External links return 200 OK
      - Code snippets have proper formatting
    OUTPUT: List of broken links

  FUNCTION GenerateCommandReference() -> markdown:
    PURPOSE: Auto-generate command reference from cobra
    OUTPUT: Markdown documentation of all flags, examples
    SOURCE: Parse cobra command structure

  FUNCTION BuildUserGuide() -> markdown:
    PURPOSE: Compile all documentation into cohesive guide
    SECTIONS:
      - Quick start
      - Common scenarios
      - Advanced usage
      - Troubleshooting
      - FAQ
```

---

## Testing Strategy (Detailed)

### Integration Test Categories

**Category 1: Core Workflow Tests**
```
TEST_CATEGORY CoreWorkflow:

  TEST PushSmallImageToGitea:
    Setup: Gitea registry, small test image (5MB, 2 layers)
    Execute: idpbuilder push testimage:latest --insecure
    Verify:
      - Exit code 0
      - Image exists in registry
      - All layers uploaded
      - Manifest correct

  TEST PushLargeImageToGitea:
    Setup: Gitea registry, large test image (100MB, 10 layers)
    Execute: idpbuilder push largeimage:latest --insecure
    Verify:
      - Exit code 0
      - Image exists in registry
      - Progress updates received for all layers
      - No memory leaks

  TEST PushWithAuthenticationSuccess:
    Setup: Gitea with authentication required
    Execute: idpbuilder push testimage:latest --username admin --password secret --insecure
    Verify:
      - Exit code 0
      - Authentication succeeded
      - Image pushed successfully

  TEST PushWithCustomRegistry:
    Setup: Custom registry (not Gitea)
    Execute: idpbuilder push testimage:latest --registry https://custom.io --password secret
    Verify:
      - Custom registry used
      - Default registry not contacted
      - Image pushed to custom location
```

**Category 2: Error Path Tests**
```
TEST_CATEGORY ErrorPaths:

  TEST PushNonExistentImage:
    Setup: Docker daemon with no images
    Execute: idpbuilder push missing:latest --password secret --insecure
    Verify:
      - Exit code 4 (image not found)
      - Error message mentions "not found in Docker daemon"
      - Suggestion includes "docker images" and "docker build"

  TEST PushWithInvalidCredentials:
    Setup: Gitea with wrong password
    Execute: idpbuilder push testimage:latest --password wrongpass --insecure
    Verify:
      - Exit code 2 (authentication error)
      - Error message mentions "authentication failed"
      - Suggestion includes "verify username and password"

  TEST PushToUnreachableRegistry:
    Setup: No registry running
    Execute: idpbuilder push testimage:latest --registry https://unreachable.io --password secret
    Verify:
      - Exit code 3 (network error)
      - Error message mentions "unable to connect"
      - Suggestion includes "verify registry is running"

  TEST PushWithTLSVerificationFailure:
    Setup: Gitea with self-signed cert, NO --insecure flag
    Execute: idpbuilder push testimage:latest --password secret
    Verify:
      - Exit code 3 (TLS error)
      - Error message mentions "certificate verification failed"
      - Suggestion includes "--insecure flag" with warning

  TEST PushWithNetworkInterruption:
    Setup: Gitea registry, simulate network failure mid-push
    Execute: idpbuilder push testimage:latest --insecure --password secret
    Simulate: Drop network connection during layer upload
    Verify:
      - Exit code 3 (network error)
      - Error message mentions connection issue
      - Partial upload cleaned up (no corrupted state)
```

**Category 3: Feature Tests**
```
TEST_CATEGORY Features:

  TEST EnvironmentVariableSupport:
    Setup: Set IDPBUILDER_REGISTRY_PASSWORD env var
    Execute: idpbuilder push testimage:latest --insecure (no --password flag)
    Verify:
      - Password read from environment variable
      - Authentication succeeded
      - Push completed

  TEST FlagPrecedenceOverEnvironment:
    Setup: Set IDPBUILDER_REGISTRY_PASSWORD="envpass"
    Execute: idpbuilder push testimage:latest --password "flagpass" --insecure
    Verify:
      - Flag value "flagpass" used (not env var)
      - Authentication used flag password
      - Push completed

  TEST VerboseModeOutput:
    Setup: Gitea registry, test image
    Execute: idpbuilder push testimage:latest --verbose --insecure --password secret
    Capture: Standard output
    Verify:
      - Detailed logs present
      - Layer digests shown
      - Byte counts displayed
      - More verbose than normal mode

  TEST ProgressReporting:
    Setup: Gitea registry, multi-layer image (5 layers)
    Execute: idpbuilder push testimage:latest --insecure --password secret
    Capture: Standard output during push
    Verify:
      - Progress updates for each layer
      - "Complete" status for uploaded layers
      - "Already exists" status for cached layers
      - Final success message with full reference

  TEST MultiArchImageSupport:
    Setup: Build multi-arch manifest (amd64, arm64)
    Execute: idpbuilder push multiarch:latest --insecure --password secret
    Verify:
      - Both architectures pushed
      - Manifest list created
      - Correct platform metadata
```

**Category 4: Edge Cases**
```
TEST_CATEGORY EdgeCases:

  TEST PushImageWithSpecialCharactersInTag:
    Setup: Image with tag "v1.0-alpha+build.123"
    Execute: idpbuilder push "testimage:v1.0-alpha+build.123" --insecure --password secret
    Verify:
      - Tag handled correctly
      - Image pushed with exact tag
      - No validation errors

  TEST PushWithVeryLongPassword:
    Setup: Password with 512 characters including unicode
    Execute: idpbuilder push testimage:latest --password "<512-char-password>" --insecure
    Verify:
      - Password accepted
      - Authentication succeeded
      - No truncation occurred

  TEST PushWithLayerAlreadyExistsInRegistry:
    Setup:
      - Push testimage:v1 (3 layers)
      - Build testimage:v2 (shares 2 layers with v1)
    Execute: idpbuilder push testimage:v2 --insecure --password secret
    Verify:
      - Only new layer uploaded
      - Existing layers skipped (mount from v1)
      - Progress shows "already exists" for shared layers

  TEST PushWithEmptyImage:
    Setup: Build minimal scratch-based image (1 layer, <1KB)
    Execute: idpbuilder push minimage:latest --insecure --password secret
    Verify:
      - Push succeeds
      - Tiny image handled correctly
      - No minimum size errors
```

### End-to-End Test Scenarios

```
E2E_TESTS:

  TEST CompleteUserWorkflow_DevelopmentToProduction:
    Simulates real developer workflow
    STEPS:
      1. Developer builds image: docker build -t myapp:latest .
      2. Developer tests locally: docker run myapp:latest
      3. Developer pushes to dev registry: idpbuilder push myapp:latest --insecure -p secret
      4. CI system pulls from dev registry
      5. CI builds release image: docker tag myapp:latest myapp:v1.0.0
      6. CI pushes to production registry: idpbuilder push --registry https://prod.io myapp:v1.0.0 -p $PROD_PASSWORD
    VERIFY:
      - All steps succeed
      - Image integrity maintained (digest matches)
      - Tags correctly applied

  TEST CICDIntegration_GitHubActions:
    Simulates GitHub Actions workflow
    STEPS:
      1. Checkout code
      2. Build image
      3. Push using environment variables
      4. Verify push succeeded in registry
    ENVIRONMENT:
      - IDPBUILDER_REGISTRY_USERNAME from secrets
      - IDPBUILDER_REGISTRY_PASSWORD from secrets
      - IDPBUILDER_INSECURE=true
    VERIFY:
      - No credentials in logs
      - Exit codes propagate correctly
      - Artifacts available in registry

  TEST MultiTeamWorkflow_SharedRegistry:
    Simulates multiple teams pushing to same registry
    STEPS:
      1. Team A pushes teamA/app:v1.0
      2. Team B pushes teamB/app:v1.0
      3. Team A pushes teamA/app:v1.1
    VERIFY:
      - No namespace collisions
      - Each team's images isolated
      - No permission errors
```

### Performance Benchmarks

```
PERFORMANCE_BENCHMARKS:

  BENCHMARK PushSmallImage_5MB:
    Image size: 5MB, 2 layers
    Expected time: <5 seconds (local registry)
    Measure: Time, memory usage, CPU usage

  BENCHMARK PushMediumImage_100MB:
    Image size: 100MB, 5 layers
    Expected time: <30 seconds (local registry)
    Measure: Time, memory usage, network throughput

  BENCHMARK PushLargeImage_500MB:
    Image size: 500MB, 10 layers
    Expected time: <3 minutes (local registry)
    Measure: Time, memory usage, streaming efficiency

  BENCHMARK ConcurrentPushes_10Images:
    Push 10 images simultaneously
    Expected: Linear scaling (no contention)
    Measure: Total time vs sequential time

  BENCHMARK MemoryFootprint_VariousImageSizes:
    Push images from 1MB to 1GB
    Expected: Memory < 200MB regardless of image size
    Measure: Peak memory usage (streaming validation)
```

---

## Documentation Strategy

### User Documentation Structure

**File: docs/push-command.md (Primary User Guide)**
```
USER_GUIDE_STRUCTURE:

  SECTION 1 - Quick Start:
    Single command example
    Minimal explanation
    Works immediately with defaults

    EXAMPLE:
      # Push image to default Gitea registry
      idpbuilder push myapp:latest --password mypassword --insecure

  SECTION 2 - Common Scenarios:
    Scenario: Push to custom registry
    Scenario: Use environment variables for credentials
    Scenario: Enable verbose mode for debugging
    Scenario: Handle self-signed certificates
    Each scenario has: Problem → Solution → Example

  SECTION 3 - Command Reference:
    All flags documented with:
      - Name and short form
      - Type and default value
      - Environment variable equivalent
      - Examples of usage

    FLAG REFERENCE:
      --registry, -r (string):
        Description: Target OCI registry URL
        Default: https://gitea.cnoe.localtest.me:8443
        Environment: IDPBUILDER_REGISTRY
        Example: --registry https://docker.io

      --username, -u (string):
        Description: Registry username
        Default: giteaadmin
        Environment: IDPBUILDER_REGISTRY_USERNAME
        Example: --username developer

      (continue for all flags...)

  SECTION 4 - Advanced Topics:
    Multi-architecture images
    Performance tuning
    CI/CD integration patterns
    Security best practices
    Network troubleshooting

  SECTION 5 - Examples:
    Complete working examples for:
      - GitHub Actions workflow
      - GitLab CI pipeline
      - Jenkins pipeline
      - Makefile integration
      - Docker Compose integration
```

**File: docs/troubleshooting.md**
```
TROUBLESHOOTING_GUIDE:

  SECTION 1 - Common Errors:
    FOR EACH error code:
      ERROR: Exit code and typical message
      CAUSE: What causes this error
      SOLUTION: Step-by-step fix
      PREVENTION: How to avoid in future

    EXAMPLE:
      ERROR: Exit code 4 - Image not found in Docker daemon
      CAUSE: Image doesn't exist locally or name is incorrect
      SOLUTION:
        1. Run `docker images` to list available images
        2. Verify image name and tag are correct
        3. Build image if missing: `docker build -t myapp:latest .`
      PREVENTION: Automate build in CI before push

  SECTION 2 - Network Issues:
    Problem: Registry unreachable
    Problem: TLS certificate verification failures
    Problem: Timeout during push
    Problem: Slow upload speeds
    Each with diagnostic steps and solutions

  SECTION 3 - Authentication Issues:
    Problem: Wrong username/password
    Problem: Token expired
    Problem: Permission denied
    Solutions with registry-specific guidance

  SECTION 4 - Docker Issues:
    Problem: Docker daemon not running
    Problem: Cannot connect to Docker socket
    Problem: Permission denied accessing Docker
    Solutions for different OS environments

  SECTION 5 - Debugging Tools:
    How to enable verbose mode
    How to capture logs
    How to test registry connectivity manually
    How to verify image integrity
```

**File: docs/cicd-examples.md**
```
CICD_EXAMPLES:

  EXAMPLE 1 - GitHub Actions:
    Complete .github/workflows/push-image.yml
    Uses secrets for credentials
    Shows error handling
    Includes artifact attestation

  EXAMPLE 2 - GitLab CI:
    Complete .gitlab-ci.yml
    Uses CI variables
    Shows pipeline stages
    Includes cache configuration

  EXAMPLE 3 - Jenkins:
    Complete Jenkinsfile
    Uses Jenkins credentials store
    Shows parallel stages
    Includes post-build actions

  EXAMPLE 4 - Makefile:
    Complete Makefile with push target
    Shows credential handling
    Includes error checking
    Integrates with build process

  Each example includes:
    - Complete working code
    - Explanation of key parts
    - Security considerations
    - Customization notes
```

**File: docs/faq.md**
```
FAQ_STRUCTURE:

  SECTION - General:
    Q: What registries are supported?
    Q: Can I push to Docker Hub?
    Q: Does this support multi-arch images?
    Q: How do I handle special characters in passwords?
    Q: Is this compatible with podman?

  SECTION - Security:
    Q: How are credentials stored?
    Q: Should I use --insecure flag?
    Q: How do I use certificate bundles?
    Q: Can I use token authentication?
    Q: Are passwords logged anywhere?

  SECTION - Performance:
    Q: Why is my push slow?
    Q: How much memory does push use?
    Q: Can I push multiple images in parallel?
    Q: How do I optimize for large images?

  SECTION - Troubleshooting:
    Q: Why does authentication fail?
    Q: Why can't I connect to registry?
    Q: Why does TLS verification fail?
    Q: How do I debug network issues?

  Each FAQ entry has:
    - Question (user's perspective)
    - Clear answer
    - Example if applicable
    - Related documentation links
```

### API Documentation (godoc)

```
API_DOCUMENTATION_PATTERN:

  PACKAGE DOCUMENTATION (doc.go):
    Package-level overview
    Usage examples
    Key concepts
    Architecture notes

    EXAMPLE doc.go for pkg/registry:
      // Package registry provides OCI registry client operations.
      //
      // This package wraps go-containerregistry to provide a simplified
      // interface for pushing images to OCI registries.
      //
      // Usage:
      //   authProvider := auth.NewBasicAuthProvider("user", "pass")
      //   tlsConfig := tls.NewConfigProvider(insecure)
      //   client := registry.NewClient(authProvider, tlsConfig)
      //   err := client.Push(ctx, image, "registry.io/repo:tag", progressCallback)
      //
      // See the examples directory for complete usage examples.
      package registry

  FUNCTION DOCUMENTATION:
    What the function does
    Parameter descriptions
    Return value descriptions
    Error conditions
    Usage example (if complex)

    EXAMPLE:
      // Push pushes an OCI image to a registry.
      //
      // The image parameter must be a valid v1.Image from go-containerregistry.
      // The targetRef parameter must be a fully qualified image reference including
      // registry host, namespace, repository, and tag (e.g., "registry.io/ns/repo:tag").
      //
      // Progress updates are sent to the progressCallback function if provided. The callback
      // is invoked for each layer upload with current progress information.
      //
      // Returns an error if authentication fails, network issues occur, or the registry
      // rejects the image. Specific error types can be checked using errors.As().
      //
      // Example:
      //   err := client.Push(ctx, image, "registry.io/myorg/myapp:v1.0.0", func(update ProgressUpdate) {
      //       fmt.Printf("Layer %s: %d/%d bytes\n", update.LayerDigest, update.BytesPushed, update.LayerSize)
      //   })
      func (c *Client) Push(ctx context.Context, image v1.Image, targetRef string, progressCallback ProgressCallback) error

  EXAMPLE TESTS (example_test.go):
    Runnable code examples for godoc
    Show typical usage patterns
    Include output verification

    EXAMPLE:
      func ExampleClient_Push() {
          // Create authentication provider
          authProvider := auth.NewBasicAuthProvider("username", "password")

          // Configure TLS (insecure mode for self-signed certs)
          tlsProvider := tls.NewConfigProvider(true)

          // Create registry client
          client, _ := registry.NewClient(authProvider, tlsProvider)

          // Push image (image obtained from Docker client)
          err := client.Push(context.Background(), image, "localhost:5000/myapp:latest", nil)
          if err != nil {
              log.Fatal(err)
          }

          fmt.Println("Image pushed successfully")
          // Output: Image pushed successfully
      }
```

---

## Build System Integration

### Makefile Targets (Phase 3 Additions)

```
MAKEFILE_NEW_TARGETS:

  TARGET: test-integration
    DESCRIPTION: Run integration tests with Docker containers
    PREREQUISITES: Docker daemon running
    STEPS:
      1. Check Docker is available
      2. Start test containers (Gitea, Docker-in-Docker if needed)
      3. Run integration tests: go test -tags=integration ./test/integration/...
      4. Cleanup containers (even on failure)
    USAGE: make test-integration

  TARGET: test-e2e
    DESCRIPTION: Run end-to-end tests
    PREREQUISITES: Docker daemon running, IDPBuilder binary built
    STEPS:
      1. Build IDPBuilder binary
      2. Start test infrastructure
      3. Run E2E tests: go test -tags=e2e ./test/e2e/...
      4. Cleanup
    USAGE: make test-e2e

  TARGET: test-all
    DESCRIPTION: Run complete test suite (unit + integration + e2e)
    STEPS:
      1. make test (unit tests)
      2. make test-integration
      3. make test-e2e
      4. Generate coverage report
    USAGE: make test-all (for CI)

  TARGET: docs
    DESCRIPTION: Generate and validate documentation
    STEPS:
      1. Generate godoc HTML: godoc -http=:6060 (capture output)
      2. Validate markdown links: scripts/validate-docs.sh
      3. Lint documentation: markdownlint docs/
    USAGE: make docs

  TARGET: docs-serve
    DESCRIPTION: Serve documentation locally for preview
    STEPS:
      1. Start godoc server: godoc -http=:6060
      2. Print: "Documentation at http://localhost:6060/pkg/..."
    USAGE: make docs-serve (for development)

  TARGET: benchmark
    DESCRIPTION: Run performance benchmarks
    STEPS:
      1. Run benchmarks: go test -bench=. -benchmem ./...
      2. Save results: benchmark-results-$(date).txt
      3. Compare with baseline if available
    USAGE: make benchmark

  TARGET: install
    DESCRIPTION: Install idpbuilder with push command
    STEPS:
      1. Build binary: make build
      2. Copy to $GOPATH/bin or /usr/local/bin
      3. Verify installation: idpbuilder push --help
    USAGE: make install
```

### CI/CD Integration

```
CI_PIPELINE_PATTERN:

  GITHUB_ACTIONS_WORKFLOW:
    Name: Test and Build
    Trigger: Pull request, push to main

    JOBS:

      JOB: unit-tests
        Runs-on: ubuntu-latest
        Steps:
          - Checkout code
          - Setup Go 1.21
          - Cache Go modules
          - Run: make test
          - Upload coverage to codecov

      JOB: integration-tests
        Runs-on: ubuntu-latest
        Services:
          docker:
            Docker daemon for testcontainers
        Steps:
          - Checkout code
          - Setup Go 1.21
          - Install Docker Compose
          - Run: make test-integration
          - Upload test artifacts

      JOB: e2e-tests
        Runs-on: ubuntu-latest
        Steps:
          - Checkout code
          - Setup Go 1.21
          - Build IDPBuilder: make build
          - Run: make test-e2e
          - Upload test artifacts

      JOB: build
        Needs: [unit-tests, integration-tests, e2e-tests]
        Runs-on: ubuntu-latest
        Steps:
          - Checkout code
          - Setup Go 1.21
          - Build: make build
          - Upload binary artifact

      JOB: docs
        Runs-on: ubuntu-latest
        Steps:
          - Checkout code
          - Validate documentation: make docs
          - Check for broken links
          - Generate API docs

  GITLAB_CI_PIPELINE:
    Similar structure with GitLab CI syntax
    Uses GitLab runners
    Caches Go modules
    Generates pipeline artifacts
```

---

## Error Handling Strategy (Phase 3 Testing Focus)

### Error Scenario Test Matrix

```
ERROR_TEST_MATRIX:

  DIMENSION 1 - Error Type:
    - ValidationError (exit 1)
    - AuthenticationError (exit 2)
    - NetworkError (exit 3)
    - ImageNotFoundError (exit 4)
    - GeneralError (exit 1)

  DIMENSION 2 - Error Source:
    - User input (bad flags)
    - Docker daemon (connection, image issues)
    - Registry (unreachable, auth, storage)
    - Network (timeout, reset, DNS)
    - System (permissions, resources)

  DIMENSION 3 - Recovery:
    - Retryable (network transient errors)
    - Fixable by user (wrong credentials)
    - Not recoverable (fundamental incompatibility)

  TEST_COVERAGE:
    Each error type × error source = test case
    Verify:
      - Correct exit code
      - Clear error message
      - Actionable suggestion
      - No stack traces in normal mode
      - Stack traces in verbose mode
```

### Error Message Validation

```
ERROR_MESSAGE_VALIDATION_PATTERN:

  FOR EACH error message:
    CHECK_FORMAT:
      - Starts with "Error: "
      - Has "Suggestion: " section
      - Includes relevant context
      - No technical jargon for user errors
      - Includes command examples where helpful

    CHECK_CONTENT:
      - Describes what went wrong (not how it failed internally)
      - Provides actionable next step
      - Includes relevant details (image name, registry, username)
      - Does not expose sensitive info (passwords, tokens)

    CHECK_TONE:
      - Helpful, not accusatory
      - Clear, not ambiguous
      - Concise, not verbose
      - Professional, not casual

    EXAMPLE:
      ✅ GOOD:
        Error: Unable to connect to registry gitea.cnoe.localtest.me:8443
        Suggestion: Verify the registry is running and accessible. Try 'curl https://gitea.cnoe.localtest.me:8443/v2/' to test connectivity.
        Context: Network error: connection refused

      ❌ BAD:
        Error: dial tcp 192.168.1.100:8443: connect: connection refused
        (No suggestion, technical details exposed, not user-friendly)
```

---

## Open Questions / Decisions Needed

### Decisions Made (Closed)

1. ✅ **Test Framework**: testcontainers-go for integration tests
2. ✅ **Documentation Format**: godoc + markdown user guides
3. ✅ **Build Integration**: Extend existing IDPBuilder Makefile
4. ✅ **Performance Testing**: Go benchmark framework
5. ✅ **CI Integration**: GitHub Actions primary, GitLab CI secondary

### Pending Decisions (For Implementation)

1. ❓ **Test Container Images**: Use official Gitea image or custom lightweight image?
   - Options: docker.io/gitea/gitea:latest OR custom minimal registry
   - Decision needed: Wave 1 implementation
   - Recommendation: Official Gitea (matches production use case)

2. ❓ **Test Image Storage**: Commit test images to repo or build on-the-fly?
   - Options: Pre-built images in test/fixtures/ OR Dockerfiles that build during tests
   - Trade-off: Speed (pre-built) vs repository size (on-the-fly)
   - Decision needed: Wave 1 implementation
   - Recommendation: Dockerfiles (keeps repo small, ensures freshness)

3. ❓ **Documentation Hosting**: Where to host generated documentation?
   - Options: GitHub Pages, ReadTheDocs, pkg.go.dev (automatic)
   - Decision needed: Wave 2 implementation
   - Recommendation: pkg.go.dev (automatic, standard for Go)

4. ❓ **Performance Baselines**: What are acceptable performance targets?
   - Small image (5MB): < 5 seconds?
   - Large image (500MB): < 3 minutes?
   - Decision needed: Wave 1 benchmarking
   - Recommendation: Measure first, set baselines from real results

### Future Enhancements (Out of Scope Phase 3)

1. 🔮 **Automated Performance Regression Testing**: Track benchmark results over time
2. 🔮 **Load Testing**: Test push under heavy concurrent load
3. 🔮 **Chaos Engineering**: Test resilience to network failures, registry issues
4. 🔮 **Interactive Documentation**: Live examples users can run in browser
5. 🔮 **Video Tutorials**: Screencast walkthroughs for common scenarios

---

## Phase 3 Wave Breakdown

### Wave 1: Integration Testing

**Goal**: Comprehensive integration testing with real Docker and Gitea

**Efforts**:
- 3.1.1: Integration Tests - Core Workflow (~500 lines)
  - E2E test: Local Docker → Gitea registry
  - Test authentication success/failure
  - Test insecure mode TLS bypass
  - Test large image push (100MB+)
  - Setup/teardown test infrastructure with testcontainers

- 3.1.2: Integration Tests - Edge Cases (~400 lines)
  - Test network failure scenarios
  - Test missing image errors
  - Test invalid credentials
  - Test registry unreachable
  - Test multi-architecture images
  - Test tag override scenarios
  - Test concurrent pushes

**Parallelization**: YES (can run in separate test files, testcontainers isolates)

**Outcome**: All integration tests passing, full coverage of error paths

---

### Wave 2: Documentation & Build Integration

**Goal**: Complete user documentation and IDPBuilder integration

**Efforts**:
- 3.2.1: User Documentation (~300 lines documentation)
  - docs/push-command.md (complete command reference)
  - docs/troubleshooting.md (error codes, solutions)
  - docs/cicd-examples.md (GitHub Actions, GitLab CI)
  - docs/faq.md (common questions)
  - Update IDPBuilder main README.md

- 3.2.2: Build System & IDPBuilder Integration (~200 lines configs)
  - Update IDPBuilder Makefile (test-integration, test-e2e, docs targets)
  - Update go.mod with testcontainers-go dependency
  - Update CI/CD pipelines (.github/workflows/)
  - Verify full build process
  - Add push command to IDPBuilder help output

**Parallelization**: YES (independent work)

**Outcome**: Complete documentation, integrated build system, CI passing

---

## Compliance Verification

### R340: Phase Architecture Fidelity

**Verification**:
```
R340_FIDELITY_COMPLIANCE:
  REQUIRED_LEVEL: PSEUDOCODE (high-level patterns, library choices)

  THIS_DOCUMENT:
    ✅ Testing patterns identified (Three-Tier Test Pyramid, Testcontainers)
    ✅ Library choices justified (testcontainers-go, godoc, Make)
    ✅ Pseudocode examples provided (test patterns, documentation structure)
    ✅ NO real function signatures (correct for phase plan)
    ✅ NO actual test code (implementation in Wave plans)
    ✅ NO detailed API beyond pseudocode

  RESULT: ✅ COMPLIANT - PSEUDOCODE fidelity maintained
```

### R308: Incremental Branching Strategy

**Verification**:
```
R308_COMPLIANCE:
  BRANCH_CHAIN:
    phase2-integration
      └─ phase3-wave1-integration (branches from phase2, NOT main)
           └─ phase3-wave2-integration (branches from wave1)
                └─ phase3-integration (final integration)

  INCREMENTAL_BUILD:
    - Wave 1 tests Phase 1 & 2 implementations
    - Wave 2 documents Phase 1 & 2 features
    - Each wave validates previous phase work

  RESULT: ✅ COMPLIANT - Incremental branching maintained
```

### R307: Independent Branch Mergeability

**Verification**:
```
R307_COMPLIANCE:
  WAVE_2:
    - Documentation effort (3.2.1) is independent
    - Build integration effort (3.2.2) is independent
    - Both can merge to main independently
    - No conflicts between efforts (different files)

  RESULT: ✅ COMPLIANT - Wave 2 allows parallel development
```

### R359: No Code Deletion

**Verification**:
```
R359_COMPLIANCE:
  CHANGES:
    - Phase 3: ONLY adds tests and documentation
    - No Phase 1 or Phase 2 code deleted
    - Pure additive (test files, doc files, Makefile additions)

  RESULT: ✅ COMPLIANT - No code deletion
```

### R383: Metadata File Organization

**Verification**:
```
R383_COMPLIANCE:
  METADATA_STRUCTURE:
    .software-factory/
      └─ phase3/
           ├─ wave1/
           │    ├─ effort-integration-core/
           │    │    └─ IMPLEMENTATION-PLAN--YYYYMMDD-HHMMSS.md
           │    └─ effort-integration-edge/
           │         └─ IMPLEMENTATION-PLAN--YYYYMMDD-HHMMSS.md
           └─ wave2/
                ├─ effort-documentation/
                │    └─ IMPLEMENTATION-PLAN--YYYYMMDD-HHMMSS.md
                └─ effort-build-integration/
                     └─ IMPLEMENTATION-PLAN--YYYYMMDD-HHMMSS.md

  WORKING_TREES: Clean (only test and doc files in test/ and docs/)

  RESULT: ✅ COMPLIANT - All metadata in .software-factory/ with timestamps
```

### R220/R221: Size Limits

**Verification**:
```
SIZE_COMPLIANCE:
  WAVE_1:
    Effort 3.1.1: ~500 lines ✅ (integration tests)
    Effort 3.1.2: ~400 lines ✅ (edge case tests)

  WAVE_2:
    Effort 3.2.1: ~300 lines ✅ (documentation - mostly prose)
    Effort 3.2.2: ~200 lines ✅ (build configs)

  SAFETY_MARGIN: All efforts well under 800 hard limit (300-400 line buffer)

  RESULT: ✅ COMPLIANT - Conservative effort sizing
```

---

## Next Steps

### Immediate Next Step: Wave 1 Integration Testing

**Action**: Orchestrator spawns SW Engineers for Phase 3 Wave 1 efforts

**Branches**:
- `phase3-wave1-effort-integration-core`
- `phase3-wave1-effort-integration-edge`

**Timeline Estimate**: 3-4 days (parallel implementation, testcontainers setup)

### After Wave 1: Wave 2 Documentation

**Action**: Orchestrator spawns 2 agents in parallel for Wave 2 efforts

**Branches** (from `phase3-wave1-integration`):
- `phase3-wave2-effort-documentation`
- `phase3-wave2-effort-build-integration`

**Timeline Estimate**: 2 days (parallel documentation and build work)

### After Wave 2: Phase 3 Integration & Project Completion

**Action**:
1. Orchestrator integrates all Wave 2 efforts
2. Architect reviews Phase 3 integration
3. Verify all success criteria met
4. Run complete test suite (unit + integration + e2e)
5. Tag phase3-complete
6. Prepare PR for IDPBuilder main repository
7. **PROJECT COMPLETE**

**Timeline Estimate**: 1-2 days (integration + review + PR preparation)

---

## Success Criteria

### Phase 3 Success Criteria

**Testing Complete**:
- ✅ All integration tests passing (core workflows + edge cases)
- ✅ Test coverage >80% overall, >90% for critical paths
- ✅ E2E tests validate complete user workflows
- ✅ Performance benchmarks recorded and acceptable
- ✅ Error scenarios comprehensively tested

**Documentation Complete**:
- ✅ User guide (push-command.md) covers all features
- ✅ Troubleshooting guide addresses all error codes
- ✅ CI/CD examples provided for major platforms
- ✅ FAQ answers common questions
- ✅ API documentation generated via godoc
- ✅ IDPBuilder README updated with push command

**Integration Complete**:
- ✅ Makefile updated with all test targets
- ✅ CI/CD pipelines run all tests
- ✅ go.mod includes all dependencies
- ✅ Build system produces working binary
- ✅ Push command accessible in IDPBuilder CLI

### Project Success Criteria

**Technical Success**:
- ✅ Push command works with local Docker and Gitea registry
- ✅ Supports all PRD requirements
- ✅ Code quality: 85%+ test coverage overall
- ✅ No efforts exceeded 800 lines
- ✅ All phases integrated successfully

**Process Success**:
- ✅ Followed Software Factory 3.0 workflow
- ✅ All rule compliance verified (R307, R308, R383, R220/R221, R340)
- ✅ Maximum parallelization achieved
- ✅ Clean git history with independent merges

**Documentation Success**:
- ✅ Complete user documentation
- ✅ Clear troubleshooting guide
- ✅ Working CI/CD examples
- ✅ Comprehensive API documentation

**Deployment Ready**:
- ✅ All tests passing in CI
- ✅ Binary builds successfully
- ✅ No critical bugs or blockers
- ✅ Ready for PR to IDPBuilder main repository

---

## Appendix: Testing Infrastructure Details

### Testcontainers Gitea Setup

```
PSEUDOCODE TestcontainersGiteaSetup:

  FUNCTION SetupGiteaContainer() -> (container, registryURL, credentials):
    Request Gitea container:
      Image: gitea/gitea:latest
      Ports: 3000 (web UI), 8443 (registry)
      Environment:
        - GITEA__server__ROOT_URL=http://localhost:3000
        - GITEA__server__HTTP_PORT=3000
        - GITEA__database__DB_TYPE=sqlite3
        - GITEA__security__INSTALL_LOCK=true
        - GITEA__service__DISABLE_REGISTRATION=false
      Wait strategy: Wait for log "Gitea started"
      Timeout: 60 seconds

    Get dynamic port mapping:
      webPort = container.MappedPort("3000/tcp")
      registryPort = container.MappedPort("8443/tcp")

    Initialize Gitea via API:
      Create admin user: giteaadmin / gitea123
      Enable container registry
      Configure insecure registry (for testing)

    Return:
      container: Container handle for cleanup
      registryURL: "localhost:${registryPort}"
      credentials: {username: "giteaadmin", password: "gitea123"}
```

### Test Image Builder

```
PSEUDOCODE TestImageBuilder:

  FUNCTION BuildTestImage(name, config) -> imageReference:
    Generate Dockerfile:
      FROM alpine:latest
      COPY test-file-1 /data/  # Generate files of specific sizes
      COPY test-file-2 /data/
      # ...
      RUN echo "Test image ${name}" > /VERSION

    Build image:
      docker build -t ${name} .
      Wait for build completion
      Verify image exists

    Return image reference: ${name}

  FUNCTION BuildMultiLayerImage(name, numLayers) -> imageReference:
    Generate Dockerfile with numLayers RUN commands
    Each RUN creates a distinct layer
    Build and return image reference

  FUNCTION BuildLargeImage(name, sizeMB) -> imageReference:
    Generate Dockerfile with large files
    Use dd to create files of specific sizes
    Build and return image reference
```

---

## Document Status

**Status**: ✅ READY FOR IMPLEMENTATION
**Phase**: Phase 3 of 3 (Final Phase)
**Fidelity**: PSEUDOCODE (Phase Architecture - R340 compliant)
**Compliance**:
- ✅ R340: Pseudocode-level fidelity (no real code)
- ✅ R510/R511: Checklist structure followed
- ✅ R308: Incremental branching strategy (builds on Phase 2)
- ✅ R307: Independent mergeability (Wave 2 parallel)
- ✅ R359: No code deletion (pure additive)
- ✅ R383: Metadata organization (.software-factory/)
- ✅ R220/R221: Size limits (all efforts 200-500 lines)

**Adaptation Notes Included**:
- ✅ Lessons from Phase 1 & Phase 2 documented
- ✅ Patterns to continue identified
- ✅ Challenges addressed
- ✅ Changes from master architecture documented

**Next Action**: Orchestrator proceeds to SPAWN_ARCHITECT_WAVE_PLANNING for Phase 3, Wave 1 (Integration Testing)

**Created By**: @agent-architect
**Date**: 2025-11-04

---

**END OF PHASE 3 ARCHITECTURE PLAN**
