# TEST STRATEGY - Wave 3.2

## Test-Driven Development Strategy

Wave 3.2 follows strict Test-Driven Development (TDD) methodology. This document defines the comprehensive test strategy for implementing push operations.

## TDD Enforcement Protocol

### The Three Phases of TDD

#### 1. RED Phase (Effort 3.2.1)
**Goal**: Write failing tests that define expected behavior
**Duration**: 4 hours
**Output**: Comprehensive test suite that fails

```go
// Example RED phase test
func TestPushBasicOperation(t *testing.T) {
    client := NewMockClient()
    pusher := NewPusher(client)

    err := pusher.Push(context.Background(), "image.tar", "registry/repo:tag")

    // This MUST fail initially
    assert.NoError(t, err)
    assert.True(t, client.PushCalled)
}
```

#### 2. GREEN Phase (Effort 3.2.2 - Part 1)
**Goal**: Write minimal code to make tests pass
**Duration**: 6 hours
**Output**: Functional implementation

```go
// Minimal implementation to pass test
func (p *Pusher) Push(ctx context.Context, image, destination string) error {
    // Just enough code to make test pass
    return p.client.Push(ctx, image, destination)
}
```

#### 3. REFACTOR Phase (Effort 3.2.2 - Part 2)
**Goal**: Optimize while keeping tests green
**Duration**: 2 hours
**Output**: Clean, optimized code

```go
// Refactored implementation
func (p *Pusher) Push(ctx context.Context, image, destination string) error {
    // Optimized, clean implementation
    if err := p.validate(image, destination); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    // ... optimized push logic
}
```

## Test Coverage Requirements

### Overall Coverage Target: 85% Minimum

#### Coverage by Component
| Component | Required Coverage | Priority |
|-----------|------------------|----------|
| Push Operations | 90% | CRITICAL |
| Progress Reporting | 85% | HIGH |
| Validation Logic | 95% | CRITICAL |
| Error Handling | 95% | CRITICAL |
| Helper Functions | 80% | MEDIUM |

### Coverage Measurement
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./pkg/oci/...
go tool cover -func=coverage.out

# HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Enforce minimum coverage
go test -cover ./pkg/oci/... | grep -E "[8-9][5-9]\.[0-9]%|100\.0%"
```

## Test Scenarios for Effort 3.2.1

### Core Functionality Tests

#### Test 1: Basic Push Operation
```go
func TestPushBasicOperation(t *testing.T) {
    // Test successful push to registry
    // Verify image uploaded correctly
    // Check destination tag applied
}
```

#### Test 2: Push with Authentication
```go
func TestPushWithAuthentication(t *testing.T) {
    // Test push with credentials
    // Verify auth headers sent
    // Check token refresh if needed
}
```

#### Test 3: Push to Insecure Registry
```go
func TestPushToInsecureRegistry(t *testing.T) {
    // Test push with --insecure flag
    // Verify TLS verification disabled
    // Check connection successful
}
```

### Progress Reporting Tests

#### Test 4: Progress Reporting
```go
func TestPushProgressReporting(t *testing.T) {
    // Test progress callbacks
    // Verify progress updates
    // Check completion notification
}
```

#### Test 5: Large Image Progress
```go
func TestPushLargeImageProgress(t *testing.T) {
    // Test with >1GB image
    // Verify streaming without memory issues
    // Check progress granularity
}
```

### Error Handling Tests

#### Test 6: Network Failure Recovery
```go
func TestPushNetworkFailure(t *testing.T) {
    // Simulate network interruption
    // Verify retry logic
    // Check exponential backoff
}
```

#### Test 7: Invalid Image Format
```go
func TestPushInvalidImage(t *testing.T) {
    // Test with corrupted image
    // Verify validation catches issue
    // Check error message clarity
}
```

#### Test 8: Registry Rejection
```go
func TestPushRegistryRejection(t *testing.T) {
    // Test unauthorized push
    // Test quota exceeded
    // Test invalid repository
}
```

### Concurrency Tests

#### Test 9: Concurrent Push Operations
```go
func TestPushConcurrent(t *testing.T) {
    // Test multiple simultaneous pushes
    // Verify thread safety
    // Check resource management
}
```

### Edge Cases

#### Test 10: Edge Case Scenarios
```go
func TestPushEdgeCases(t *testing.T) {
    // Empty image
    // Extremely long tag names
    // Special characters in names
    // Network timeouts
}
```

## Test Data and Fixtures

### Test Image Fixtures
```
pkg/oci/testdata/
├── images/
│   ├── tiny.tar         # 1MB test image
│   ├── small.tar        # 10MB test image
│   ├── medium.tar       # 100MB test image
│   ├── large.tar        # 1GB test image
│   └── corrupted.tar    # Invalid image for error testing
├── manifests/
│   ├── v2.json          # OCI v2 manifest
│   └── list.json        # Manifest list
└── configs/
    ├── secure.yaml      # Secure registry config
    └── insecure.yaml    # Insecure registry config
```

### Mock Registry Setup
```go
// Mock registry for testing
type MockRegistry struct {
    PushCount    int
    LastImage    string
    LastTag      string
    ShouldFail   bool
    FailureType  string
}

func NewMockRegistry() *MockRegistry {
    // Initialize mock registry
}
```

## Test Execution Strategy

### Unit Test Execution
```bash
# Run all unit tests
go test ./pkg/oci/...

# Run specific test
go test -run TestPush ./pkg/oci/

# Verbose output
go test -v ./pkg/oci/...

# With race detection
go test -race ./pkg/oci/...
```

### Integration Test Execution
```bash
# Run integration tests (requires real registry)
go test -tags=integration ./pkg/oci/...

# With Docker registry
docker run -d -p 5000:5000 registry:2
go test -tags=integration -registry=localhost:5000 ./pkg/oci/...
```

### Performance Test Execution
```bash
# Run benchmarks
go test -bench=. -benchmem ./pkg/oci/...

# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=. ./pkg/oci/...

# Profile memory usage
go test -memprofile=mem.prof -bench=. ./pkg/oci/...
```

## Test Quality Standards

### Test Naming Convention
```go
// Pattern: Test<Component><Scenario><Expectation>
TestPushBasicOperationSuccess
TestPushAuthenticationFailure
TestProgressReportingAccuracy
```

### Test Structure
```go
func TestExample(t *testing.T) {
    // Arrange
    setup test data and mocks

    // Act
    perform operation

    // Assert
    verify expectations

    // Cleanup (if needed)
    defer cleanup()
}
```

### Assertion Guidelines
- Use clear, specific assertions
- One logical assertion per test
- Descriptive failure messages
- Table-driven tests for multiple scenarios

```go
// Table-driven test example
func TestPushScenarios(t *testing.T) {
    tests := []struct {
        name     string
        image    string
        tag      string
        wantErr  bool
    }{
        {"valid push", "image.tar", "repo:tag", false},
        {"invalid image", "", "repo:tag", true},
        {"invalid tag", "image.tar", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}
```

## Test Validation Criteria

### For Effort 3.2.1 (Test Development)
- [ ] All test scenarios implemented
- [ ] Tests compile without errors
- [ ] Tests fail with clear messages
- [ ] No false positives
- [ ] Test data prepared
- [ ] Mocks implemented

### For Effort 3.2.2 (Implementation)
- [ ] All tests passing
- [ ] Coverage ≥85%
- [ ] No test modifications (except bug fixes)
- [ ] Performance acceptable
- [ ] No flaky tests
- [ ] CI/CD integration ready

## Continuous Testing

### Pre-Commit Hooks
```bash
#!/bin/bash
# .git/hooks/pre-commit

# Run tests before commit
go test ./pkg/oci/... || exit 1

# Check coverage
coverage=$(go test -cover ./pkg/oci/... | grep -oE '[0-9]+\.[0-9]%')
if [ "${coverage%.*}" -lt 85 ]; then
    echo "Coverage below 85%"
    exit 1
fi
```

### CI/CD Pipeline
```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
      - run: go test -v -cover ./...
      - run: go test -race ./...
      - run: go test -tags=integration ./...
```

## Test Maintenance

### Test Review Checklist
- [ ] Tests reflect current requirements
- [ ] No duplicate test scenarios
- [ ] Tests are independent
- [ ] Tests are deterministic
- [ ] Tests complete quickly (<100ms for unit)
- [ ] Tests have clear documentation

### Test Refactoring Guidelines
- Keep tests simple and readable
- Extract common setup to helpers
- Use test fixtures consistently
- Update tests with implementation changes
- Remove obsolete tests

## Risk Mitigation Through Testing

### High-Risk Areas Requiring Extra Testing
1. **Large File Handling**: Memory/streaming tests
2. **Network Reliability**: Failure/retry tests
3. **Authentication**: Token refresh/expiry tests
4. **Concurrency**: Race condition tests
5. **Registry Compatibility**: Multi-registry tests

### Test Failure Response
```
Test Failure → Analyze → Fix Code (not test) → Verify → Continue
                  ↓
            If test wrong → Document → Get approval → Fix test
```

## Success Metrics

### Test Quality Metrics
- **Coverage**: ≥85% overall, ≥90% critical paths
- **Execution Time**: <10s for unit tests
- **Stability**: 0% flaky tests
- **Clarity**: 100% tests self-documenting
- **Independence**: 100% tests runnable in isolation

### TDD Compliance Metrics
- **Test-First**: 100% tests written before code
- **Red-Green-Refactor**: All three phases completed
- **Test Commits**: Tests committed before implementation
- **Coverage Growth**: Coverage never decreases

---

**Document Status**: COMPLETE
**Created**: 2025-09-25T17:40:29Z
**Methodology**: Test-Driven Development
**Coverage Target**: 85% minimum

*This test strategy ensures comprehensive validation of Wave 3.2 push operations through disciplined TDD approach.*