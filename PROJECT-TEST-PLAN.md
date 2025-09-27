# PROJECT TEST PLAN - IDPBuilder OCI Build & Push Feature

## 📋 Overview

**Document Type**: Test-Driven Development (TDD) Specification
**Created**: 2025-01-27
**Rule Compliance**: R341 (Tests Before Implementation)
**Project**: IDPBuilder OCI Build & Push
**Test Strategy**: Create comprehensive test specifications BEFORE implementation begins

### Test Philosophy
Following TDD principles, these tests define the expected behavior and serve as the implementation guide. Each test must be specific enough to drive development decisions while being flexible enough to allow implementation choices.

### Success Criteria
- All tests defined here must eventually pass
- No implementation begins until tests are specified
- Tests guide the architecture and design decisions
- Each phase has clear, measurable acceptance criteria

## 🎯 Test Coverage Requirements

| Phase | Unit Tests | Integration Tests | E2E Tests | Minimum Coverage |
|-------|------------|-------------------|-----------|------------------|
| Phase 1 | Required | Required | Optional | 70% |
| Phase 2 | Required | Required | Required | 80% |
| Phase 3 | Required | Required | Required | 90% |
| Phase 4 | Test Suite | Test Suite | Test Suite | 85% |
| Phase 5 | Required | Optional | Required | 80% |

## 🧪 Phase 1 Tests: CLI Foundation & Authentication

### 1.1 Unit Test Specifications

#### Stack Configuration Types
```go
// Test: Stack configuration parsing
func TestStackConfigurationParsing(t *testing.T) {
    // Given: A stack configuration with container build section
    config := `
    containerBuild:
      enabled: true
      dockerfile: ./Dockerfile
      context: .
      args:
        VERSION: "1.0.0"
    `

    // When: Configuration is parsed
    stack := ParseStackConfig(config)

    // Then: Container build settings are correctly extracted
    assert.True(t, stack.ContainerBuild.Enabled)
    assert.Equal(t, "./Dockerfile", stack.ContainerBuild.Dockerfile)
    assert.Equal(t, "1.0.0", stack.ContainerBuild.Args["VERSION"])
}

// Test: Invalid configuration handling
func TestInvalidStackConfiguration(t *testing.T) {
    // Given: Invalid configuration
    // When: Parsing attempted
    // Then: Appropriate error returned with clear message
}
```

#### Buildah Abstraction Interfaces
```go
// Test: Buildah builder interface implementation
func TestBuildahBuilderInterface(t *testing.T) {
    // Given: A mock Buildah builder
    builder := NewMockBuilder()

    // When: Build is initiated
    result, err := builder.Build(BuildOptions{
        Dockerfile: "Dockerfile",
        Context:    "/path/to/context",
        Tags:       []string{"myapp:latest"},
    })

    // Then: Build completes successfully
    assert.NoError(t, err)
    assert.NotNil(t, result.ImageID)
    assert.Contains(t, result.Logs, "Successfully built")
}

// Test: Build failure scenarios
func TestBuildahBuilderErrorHandling(t *testing.T) {
    // Test missing Dockerfile
    // Test invalid context
    // Test build step failures
}
```

#### Registry Client Interfaces
```go
// Test: Registry authentication
func TestGiteaRegistryAuthentication(t *testing.T) {
    // Given: Gitea credentials
    client := NewRegistryClient(RegistryOptions{
        URL:      "https://gitea.cnoe.localtest.me",
        Username: "admin",
        Password: "admin123",
        Insecure: true,
    })

    // When: Authentication attempted
    err := client.Authenticate()

    // Then: Successfully authenticated
    assert.NoError(t, err)
    assert.True(t, client.IsAuthenticated())
}

// Test: Certificate handling
func TestInsecureSkipVerify(t *testing.T) {
    // Given: Registry with self-signed cert
    // When: InsecureSkipVerify is true
    // Then: Connection succeeds
}
```

### 1.2 Integration Test Specifications

#### CLI Command Integration
```bash
# Test: Stack build command exists
idpbuilder stack build --help
# Expected: Command help displayed

# Test: Stack push command exists
idpbuilder stack push --help
# Expected: Command help displayed

# Test: Build with missing Dockerfile
idpbuilder stack build --name test-stack
# Expected: Error - "Dockerfile not found in stack directory"

# Test: Push without build
idpbuilder stack push --name test-stack
# Expected: Error - "No built image found for stack"
```

### 1.3 Acceptance Criteria for Phase 1

- [ ] All interfaces compile without implementation
- [ ] CLI commands are registered and accessible
- [ ] Help text is comprehensive and accurate
- [ ] Configuration parsing handles all edge cases
- [ ] Error messages are clear and actionable
- [ ] Authentication mechanism is pluggable
- [ ] Certificate handling supports self-signed certs

## 🔧 Phase 2 Tests: Core OCI Functionality

### 2.1 Unit Test Specifications

#### Build Implementation
```go
// Test: Successful container build
func TestContainerBuildSuccess(t *testing.T) {
    // Given: A valid Dockerfile and context
    dockerfile := `
    FROM alpine:latest
    RUN echo "test"
    CMD ["sh"]
    `

    // When: Build is executed
    engine := NewBuildahEngine()
    result, err := engine.Build(BuildRequest{
        Dockerfile: dockerfile,
        Context:    testContext,
        Tags:       []string{"test:latest"},
    })

    // Then: Image is built successfully
    assert.NoError(t, err)
    assert.NotEmpty(t, result.ImageID)
    assert.Contains(t, result.Manifest.Layers, "alpine")
}

// Test: Multi-stage build support
func TestMultiStageBuild(t *testing.T) {
    // Given: Multi-stage Dockerfile
    // When: Build executed
    // Then: Final stage image produced
}

// Test: Build context file inclusion
func TestBuildContextFiles(t *testing.T) {
    // Given: Context with multiple files
    // When: Build references context files
    // Then: Files are available during build
}
```

#### Registry Push Implementation
```go
// Test: Push to Gitea registry
func TestPushToGiteaRegistry(t *testing.T) {
    // Given: Built image and authenticated client
    imageID := "sha256:123abc..."
    client := AuthenticatedGiteaClient()

    // When: Push initiated
    err := client.Push(PushRequest{
        ImageID:    imageID,
        Repository: "cnoe/myapp",
        Tag:        "v1.0.0",
    })

    // Then: Image successfully pushed
    assert.NoError(t, err)

    // And: Can verify image exists in registry
    exists, err := client.ImageExists("cnoe/myapp:v1.0.0")
    assert.True(t, exists)
}

// Test: Push retry logic
func TestPushRetryOnFailure(t *testing.T) {
    // Given: Flaky network connection
    // When: Push fails initially
    // Then: Automatic retry succeeds
}

// Test: Push progress reporting
func TestPushProgressReporting(t *testing.T) {
    // Given: Large image to push
    // When: Push in progress
    // Then: Progress updates provided
}
```

### 2.2 Integration Test Specifications

```go
// Test: End-to-end build and push
func TestEndToEndBuildAndPush(t *testing.T) {
    // Given: A stack with Dockerfile
    stackDir := setupTestStack()

    // When: Build and push executed
    err := ExecuteCommand("idpbuilder", "stack", "build",
        "--name", "test-stack",
        "--tag", "v1.0.0")
    assert.NoError(t, err)

    err = ExecuteCommand("idpbuilder", "stack", "push",
        "--name", "test-stack")
    assert.NoError(t, err)

    // Then: Image available in Gitea
    // Verify via registry API
}

// Test: Stack discovery and processing
func TestStackDiscovery(t *testing.T) {
    // Given: Multiple stacks with Dockerfiles
    // When: Build all command executed
    // Then: All stacks built successfully
}
```

### 2.3 Acceptance Criteria for Phase 2

- [ ] Can build images from any valid Dockerfile
- [ ] Supports common Docker build features (ARG, ENV, COPY, etc.)
- [ ] Successfully pushes to Gitea registry
- [ ] Handles network interruptions gracefully
- [ ] Build context correctly includes local files
- [ ] Progress reporting provides meaningful feedback
- [ ] Stack integration discovers Dockerfiles automatically

## 🚀 Phase 3 Tests: Advanced Registry Features

### 3.1 Unit Test Specifications

#### Certificate Bundle Loading
```go
// Test: Load certificates from Kind cluster
func TestKindCertificateLoading(t *testing.T) {
    // Given: Kind cluster with certificates
    cluster := NewKindCluster("test-cluster")

    // When: Certificates extracted
    certs, err := ExtractClusterCertificates(cluster)

    // Then: Valid certificate bundle returned
    assert.NoError(t, err)
    assert.NotEmpty(t, certs.CACert)
    assert.True(t, certs.IsValid())
}

// Test: Dynamic certificate validation
func TestDynamicCertificateValidation(t *testing.T) {
    // Given: Registry with changing certificates
    // When: Certificate changes detected
    // Then: New certificates loaded automatically
}
```

#### Build Optimization
```go
// Test: Build cache reduces rebuild time
func TestBuildCacheOptimization(t *testing.T) {
    // Given: Previously built image layers
    firstBuild := timeExecution(func() {
        engine.Build(options)
    })

    // When: Rebuild with minimal changes
    secondBuild := timeExecution(func() {
        engine.Build(optionsWithCache)
    })

    // Then: Second build is 50% faster
    assert.Less(t, secondBuild, firstBuild*0.5)
}

// Test: Layer reuse across builds
func TestLayerReuse(t *testing.T) {
    // Given: Multiple images with common base
    // When: Building second image
    // Then: Base layers reused from cache
}
```

#### Batch Operations
```go
// Test: Batch build multiple stacks
func TestBatchBuildOperation(t *testing.T) {
    // Given: 5 stacks to build
    stacks := []string{"app1", "app2", "app3", "app4", "app5"}

    // When: Batch build initiated
    results := BatchBuild(stacks, BatchOptions{
        Parallel: 3,
        FailFast: false,
    })

    // Then: All builds complete
    assert.Len(t, results.Successful, 5)
    assert.Empty(t, results.Failed)
}

// Test: Parallel push orchestration
func TestParallelPush(t *testing.T) {
    // Given: Multiple built images
    // When: Parallel push initiated
    // Then: All images pushed concurrently
}
```

### 3.2 Integration Test Specifications

```go
// Test: Production-like certificate handling
func TestProductionCertificateHandling(t *testing.T) {
    // Given: Production-like environment
    // When: Various certificate scenarios
    // Then: All handled correctly
}

// Test: Large-scale batch operations
func TestLargeScaleBatchOperations(t *testing.T) {
    // Given: 20+ stacks
    // When: Batch build and push
    // Then: Completes within performance targets
}
```

### 3.3 Acceptance Criteria for Phase 3

- [ ] Certificate bundle loading works with Kind cluster
- [ ] Build cache provides measurable performance improvement
- [ ] Layer reuse reduces storage requirements
- [ ] Batch operations handle failures gracefully
- [ ] Parallel operations respect resource limits
- [ ] Multi-stage builds optimize final image size
- [ ] Secret management prevents credential leakage

## 🧪 Phase 4 Tests: Testing & Validation

### 4.1 Meta-Test Specifications (Tests for Tests)

```go
// Test: Unit test coverage measurement
func TestUnitTestCoverage(t *testing.T) {
    // Given: All unit tests
    // When: Coverage calculated
    // Then: Meets 85% threshold
}

// Test: Integration test reliability
func TestIntegrationTestReliability(t *testing.T) {
    // Given: Integration test suite
    // When: Run 100 times
    // Then: >95% pass rate (no flaky tests)
}
```

### 4.2 Performance Benchmarks

```go
// Benchmark: Container build performance
func BenchmarkContainerBuild(b *testing.B) {
    // Measure: Time to build typical stack
    // Target: <2 minutes for average stack
}

// Benchmark: Registry push throughput
func BenchmarkRegistryPush(b *testing.B) {
    // Measure: Push speed for various image sizes
    // Target: >10MB/s for local registry
}

// Benchmark: Memory usage during build
func BenchmarkBuildMemoryUsage(b *testing.B) {
    // Measure: Peak memory during build
    // Target: <500MB for typical build
}
```

### 4.3 Security Validation Tests

```go
// Test: No credential leakage in logs
func TestNoCredentialLeakage(t *testing.T) {
    // Given: Build with secrets
    // When: Build executes with verbose logging
    // Then: No secrets appear in output
}

// Test: Image vulnerability scanning
func TestImageVulnerabilityScanning(t *testing.T) {
    // Given: Built images
    // When: Scanned for vulnerabilities
    // Then: No critical vulnerabilities in base images
}
```

### 4.4 Acceptance Criteria for Phase 4

- [ ] Unit test coverage exceeds 85%
- [ ] Integration tests have <5% flake rate
- [ ] Performance benchmarks meet targets
- [ ] Security tests pass without violations
- [ ] E2E tests cover all user workflows
- [ ] Test execution time is reasonable (<10 minutes)

## ✨ Phase 5 Tests: Documentation & Finalization

### 5.1 Documentation Validation Tests

```bash
# Test: CLI help completeness
for cmd in build push list delete; do
    idpbuilder stack $cmd --help | grep -q "Examples:"
    # Expected: All commands have examples
done

# Test: Documentation code examples work
grep -h "```bash" docs/oci/*.md | bash -
# Expected: All examples execute successfully
```

### 5.2 User Experience Tests

```go
// Test: Error messages are helpful
func TestErrorMessageQuality(t *testing.T) {
    // Given: Various error conditions
    errors := TriggerAllErrorScenarios()

    // Then: Each error message includes:
    for _, err := range errors {
        assert.Contains(t, err.Error(), "what went wrong")
        assert.Contains(t, err.Error(), "how to fix")
    }
}

// Test: Progress feedback is meaningful
func TestProgressFeedback(t *testing.T) {
    // Given: Long-running operation
    // When: Operation in progress
    // Then: User receives regular, meaningful updates
}
```

### 5.3 Release Readiness Tests

```bash
# Test: Version command works
idpbuilder version | grep -q "OCI: enabled"
# Expected: OCI feature status shown

# Test: Metrics collection operational
idpbuilder stack build --name test --metrics
# Expected: Metrics logged to configured endpoint

# Test: Feature flags work
IDPBUILDER_OCI_ENABLED=false idpbuilder stack build
# Expected: Error - "OCI features disabled"
```

### 5.4 Acceptance Criteria for Phase 5

- [ ] All documentation examples execute correctly
- [ ] CLI help text is comprehensive
- [ ] Error messages guide users to solutions
- [ ] Telemetry collects useful metrics
- [ ] Feature can be enabled/disabled
- [ ] Release notes document all changes

## 🔄 Integration Test Scenarios

### Cross-Phase Integration Tests

```go
// Test: Phase 1-2 Integration
func TestPhase1to2Integration(t *testing.T) {
    // Given: Phase 1 interfaces
    // When: Phase 2 implements them
    // Then: No interface changes required
}

// Test: Phase 2-3 Integration
func TestPhase2to3Integration(t *testing.T) {
    // Given: Phase 2 basic functionality
    // When: Phase 3 adds advanced features
    // Then: Basic functionality still works
}

// Test: Full project integration
func TestFullProjectIntegration(t *testing.T) {
    // Given: All phases complete
    // When: End-to-end workflow executed
    // Then: Everything works together seamlessly
}
```

## 🎯 API Contract Tests

### External API Contracts

```go
// Test: Gitea Registry API compatibility
func TestGiteaRegistryAPIContract(t *testing.T) {
    // Given: Gitea registry endpoints
    endpoints := []string{
        "/v2/",
        "/v2/_catalog",
        "/v2/{name}/manifests/{reference}",
        "/v2/{name}/blobs/{digest}",
    }

    // When: Each endpoint called
    // Then: Responses match OCI Distribution Spec
}

// Test: Buildah library API stability
func TestBuildahAPIContract(t *testing.T) {
    // Given: Buildah v1.33.0 API
    // When: Our code uses the API
    // Then: No deprecated methods used
}
```

## 📊 Test Execution Strategy

### Test Execution Order
1. **Unit Tests First**: Run on every commit
2. **Integration Tests**: Run on PR creation
3. **E2E Tests**: Run before merge
4. **Performance Tests**: Run nightly
5. **Security Tests**: Run weekly

### Test Environment Requirements
```yaml
test_environments:
  unit:
    requirements:
      - Go 1.21+
      - Mock interfaces
    isolation: process

  integration:
    requirements:
      - Kind cluster
      - Gitea instance
      - Test containers
    isolation: container

  e2e:
    requirements:
      - Full IDPBuilder setup
      - Sample stacks
      - ArgoCD instance
    isolation: cluster
```

## 🚨 Test Failure Protocol

### Failure Response Matrix
| Test Type | Failure Action | Severity | Auto-Retry |
|-----------|---------------|----------|------------|
| Unit | Block commit | Critical | No |
| Integration | Block PR | High | Yes (1x) |
| E2E | Block merge | Critical | Yes (2x) |
| Performance | Warning | Medium | No |
| Security | Block release | Critical | No |

## 📈 Test Metrics

### Key Test Metrics to Track
```yaml
metrics:
  coverage:
    target: 85%
    measure: lines_covered / total_lines

  reliability:
    target: 95%
    measure: pass_rate_over_100_runs

  performance:
    build_time: < 2_minutes
    push_time: < 30_seconds
    test_suite: < 10_minutes

  quality:
    defect_escape_rate: < 5%
    test_effectiveness: > 90%
```

## 🔍 Test Verification Commands

### Quick Test Verification
```bash
# Run unit tests for specific phase
go test ./pkg/oci/... -v -cover

# Run integration tests
go test ./tests/integration/... -tags=integration

# Run E2E tests with Kind cluster
./scripts/run-e2e-tests.sh

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. -benchmem ./pkg/oci/...
```

## ✅ Test Acceptance Checklist

### Before Phase 1 Implementation Begins
- [ ] All Phase 1 test specifications reviewed
- [ ] Test interfaces compilable
- [ ] Mock implementations ready
- [ ] Test data prepared
- [ ] CI/CD pipeline configured

### Before Each Phase
- [ ] Previous phase tests passing
- [ ] New phase tests specified
- [ ] Test environment ready
- [ ] Dependencies available
- [ ] Performance baselines established

### Before Project Completion
- [ ] All tests passing
- [ ] Coverage targets met
- [ ] No flaky tests
- [ ] Performance within limits
- [ ] Security validation complete
- [ ] Documentation tests pass

---

**Document Type**: TDD Test Specification
**Rule Compliance**: R341 (Test-Driven Development)
**Created**: 2025-01-27
**Version**: 1.0
**Status**: Ready for Phase 1 Implementation

**Remember**: Tests come FIRST, implementation follows tests!