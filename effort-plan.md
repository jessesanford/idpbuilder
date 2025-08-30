# Integration Testing Implementation Plan

## Effort Infrastructure Metadata (R209)
**EFFORT_NAME**: integration-testing
**BRANCH**: idpbuilder-oci-mvp/phase2/wave2/integration-testing
**WORKING_DIRECTORY**: /home/vscode/workspaces/idpbuilder-oci-mvp/efforts/phase2/wave2/integration-testing
**PARENT_WAVE**: Phase 2 Wave 2
**SEQUENTIAL_POSITION**: 2 of 2 (DEPENDS on cli-commands completion)

## Critical Effort Metadata (FROM WAVE PLAN - R211)
**Branch**: `idpbuilder-oci-mvp/phase2/wave2/integration-testing`
**Can Parallelize**: No
**Parallel With**: None
**Size Estimate**: 500 lines (MUST be <800)
**Dependencies**: 
  - Effort 2.2.1 (cli-commands) - MUST BE COMPLETE
  - Phase 1: All certificate efforts (for test validation)
  - Phase 2 Wave 1: Build/Registry components (indirect via CLI)
**Execution Order**: MUST run after cli-commands effort completes

## Overview
- **Effort**: Comprehensive integration and end-to-end testing of CLI functionality
- **Phase**: 2, Wave: 2
- **Estimated Size**: 500 lines (450 target, 500 soft limit)
- **Implementation Time**: 1.5 days

## Dependency Analysis (R219 Compliance)

### Direct Dependency: CLI Commands (Effort 2.2.1)
**What We Test**:
```go
// From cli-commands effort (must be completed first)
// We test the actual CLI binary/commands, not import packages
// Testing targets:
// - idpbuilder build command
// - idpbuilder push command
// - Certificate auto-configuration
// - Error handling and recovery
```

**How It Influences Implementation**:
- Tests execute actual CLI commands via exec.Command()
- Validate command output and exit codes
- Test both success and failure scenarios
- Verify certificate auto-configuration works end-to-end

### Indirect Dependencies (Via CLI)
**Phase 1 Certificate Infrastructure**:
- Test that CLI properly auto-extracts certificates
- Verify trust store configuration works
- Validate fallback to --insecure mode

**Phase 2 Wave 1 Build/Registry**:
- Test that CLI successfully builds images
- Verify push to Gitea registry works
- Validate authentication handling

## File Structure
```
integration-testing/
├── pkg/
│   ├── tests/
│   │   ├── integration/
│   │   │   ├── setup.go              # Test environment setup (~60 lines)
│   │   │   ├── fixtures.go           # Test data and fixtures (~40 lines)
│   │   │   ├── build_test.go         # Build command tests (~150 lines)
│   │   │   └── push_test.go          # Push command tests (~150 lines)
│   │   └── e2e/
│   │       └── workflow_test.go      # Complete workflow tests (~100 lines)
│   └── testutil/
│       ├── docker.go                  # Docker/Kind utilities (~30 lines)
│       ├── cli.go                     # CLI execution helpers (~40 lines)
│       └── assertions.go              # Test assertions (~30 lines)
└── test-data/
    ├── dockerfiles/
    │   ├── simple.Dockerfile          # Simple test Dockerfile
    │   ├── multistage.Dockerfile      # Multi-stage test
    │   └── invalid.Dockerfile         # Invalid syntax test
    └── contexts/
        └── simple-app/                # Test build context
            └── main.go                # Simple Go app for testing
```

## Implementation Steps

### Step 1: Test Infrastructure Setup (100 lines)
**Files**: `pkg/tests/integration/setup.go`, `pkg/tests/integration/fixtures.go`

1. Implement setup.go:
   ```go
   package integration

   type TestEnvironment struct {
       ClusterName    string
       RegistryURL    string
       TestNamespace  string
       CLIPath        string
   }

   func SetupTestEnvironment(t *testing.T) *TestEnvironment {
       // 1. Verify Kind cluster exists or create one
       // 2. Verify Gitea is running and accessible
       // 3. Locate or build CLI binary
       // 4. Create test namespace/workspace
       // 5. Initialize test environment struct
   }

   func (env *TestEnvironment) Cleanup() {
       // Clean up test artifacts
       // Remove test images
       // Reset certificate configurations
   }
   ```

2. Implement fixtures.go:
   ```go
   func GetTestDockerfile(name string) string {
       // Return path to test Dockerfile
   }

   func GetTestContext(name string) string {
       // Return path to test build context
   }

   func GenerateTestImageTag() string {
       // Generate unique image tag for test
   }
   ```

**Size checkpoint**: ~100 lines

### Step 2: Build Command Tests (150 lines)
**Files**: `pkg/tests/integration/build_test.go`

1. Test successful build scenarios:
   ```go
   func TestBuildSimpleDockerfile(t *testing.T) {
       env := SetupTestEnvironment(t)
       defer env.Cleanup()

       // Execute: idpbuilder build --file simple.Dockerfile --context . --tag test:latest
       cmd := exec.Command(env.CLIPath, "build",
           "--file", GetTestDockerfile("simple"),
           "--context", GetTestContext("simple-app"),
           "--tag", GenerateTestImageTag())
       
       output, err := cmd.CombinedOutput()
       require.NoError(t, err)
       assert.Contains(t, string(output), "Successfully built")
   }

   func TestBuildMultistageDockerfile(t *testing.T) {
       // Test multi-stage build
   }

   func TestBuildWithCertificateAutoConfig(t *testing.T) {
       // Verify certificates are auto-configured
       // Check trust store after build
   }
   ```

2. Test error scenarios:
   ```go
   func TestBuildInvalidDockerfile(t *testing.T) {
       // Should fail with clear error message
   }

   func TestBuildMissingContext(t *testing.T) {
       // Should report missing context error
   }

   func TestBuildWithoutCluster(t *testing.T) {
       // Should handle missing Kind cluster gracefully
   }
   ```

3. Test build options:
   ```go
   func TestBuildWithPlatform(t *testing.T) {
       // Test --platform flag
   }
   ```

**Test scenarios covered**:
- Simple Dockerfile build
- Multi-stage Dockerfile build
- Invalid Dockerfile (error handling)
- Missing build context
- Certificate auto-configuration
- Platform selection

**Size checkpoint**: ~250 lines total

### Step 3: Push Command Tests (150 lines)
**Files**: `pkg/tests/integration/push_test.go`

1. Test successful push scenarios:
   ```go
   func TestPushToGitea(t *testing.T) {
       env := SetupTestEnvironment(t)
       defer env.Cleanup()

       // First build an image
       buildCmd := exec.Command(env.CLIPath, "build",
           "--file", GetTestDockerfile("simple"),
           "--context", GetTestContext("simple-app"),
           "--tag", "gitea.local:443/test/app:latest")
       require.NoError(t, buildCmd.Run())

       // Then push it
       pushCmd := exec.Command(env.CLIPath, "push",
           "gitea.local:443/test/app:latest")
       
       output, err := pushCmd.CombinedOutput()
       require.NoError(t, err)
       assert.Contains(t, string(output), "Successfully pushed")
   }

   func TestPushWithInsecureFlag(t *testing.T) {
       // Test --insecure bypass mode
   }

   func TestPushWithAuthentication(t *testing.T) {
       // Test with --username and --password
   }
   ```

2. Test error scenarios:
   ```go
   func TestPushWithoutCluster(t *testing.T) {
       // Should fail gracefully
   }

   func TestPushInvalidCredentials(t *testing.T) {
       // Should report auth error
   }

   func TestPushMissingImage(t *testing.T) {
       // Should report image not found
   }

   func TestPushNetworkInterruption(t *testing.T) {
       // Simulate network issues
       // Verify retry/recovery behavior
   }
   ```

**Test scenarios covered**:
- Push with auto-configured certificates
- Push with --insecure flag
- Push without cluster (error)
- Invalid credentials (error)
- Missing image (error)
- Network interruption recovery

**Size checkpoint**: ~400 lines total

### Step 4: End-to-End Workflow Tests (100 lines)
**Files**: `pkg/tests/e2e/workflow_test.go`

1. Test complete workflows:
   ```go
   func TestCompleteWorkflow(t *testing.T) {
       env := SetupTestEnvironment(t)
       defer env.Cleanup()

       // 1. Build image
       buildCmd := exec.Command(env.CLIPath, "build",
           "--file", GetTestDockerfile("simple"),
           "--context", GetTestContext("simple-app"),
           "--tag", "gitea.local:443/test/e2e:v1")
       
       buildOutput, err := buildCmd.CombinedOutput()
       require.NoError(t, err, "Build failed: %s", buildOutput)

       // 2. Push image
       pushCmd := exec.Command(env.CLIPath, "push",
           "gitea.local:443/test/e2e:v1")
       
       pushOutput, err := pushCmd.CombinedOutput()
       require.NoError(t, err, "Push failed: %s", pushOutput)

       // 3. Verify image exists in registry
       // Could use registry API or pull to verify
   }

   func TestFreshInstallationFlow(t *testing.T) {
       // Test with no pre-existing certificates
       // Verify auto-configuration works
   }

   func TestCertificateRotation(t *testing.T) {
       // Simulate certificate change
       // Verify CLI handles it properly
   }

   func TestRecoveryFromFailure(t *testing.T) {
       // Test recovery from various failure modes
   }

   func TestConcurrentOperations(t *testing.T) {
       // Run multiple builds/pushes in parallel
       // Verify no conflicts or race conditions
   }
   ```

**Workflow scenarios**:
- Fresh installation (no certificates)
- Complete build and push workflow
- Certificate rotation handling
- Recovery from failed operations
- Concurrent operations

**Size checkpoint**: ~500 lines total

### Step 5: Test Utilities (100 lines)
**Files**: `pkg/testutil/docker.go`, `pkg/testutil/cli.go`, `pkg/testutil/assertions.go`

1. Docker/Kind utilities (docker.go):
   ```go
   func IsKindClusterRunning(name string) bool {
       // Check if Kind cluster exists
   }

   func GetGiteaURL() string {
       // Get Gitea registry URL
   }
   ```

2. CLI execution helpers (cli.go):
   ```go
   func RunCLICommand(args ...string) (string, error) {
       // Helper to run CLI commands
   }

   func RunCLICommandWithEnv(env map[string]string, args ...string) (string, error) {
       // Run with custom environment
   }
   ```

3. Test assertions (assertions.go):
   ```go
   func AssertImageExists(t *testing.T, imageRef string) {
       // Verify image exists locally or in registry
   }

   func AssertCertificateConfigured(t *testing.T) {
       // Verify certificate trust is configured
   }
   ```

**Final size**: ~500 lines

## Size Management
- **Estimated Lines**: 500
- **Measurement Tool**: /home/vscode/workspaces/idpbuilder-oci-mvp/tools/line-counter.sh
- **Check Frequency**: After each step completion
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

### Size Monitoring Points
1. After Step 1 (Test Infrastructure): Check = 100 lines
2. After Step 2 (Build Tests): Check = 250 lines
3. After Step 3 (Push Tests): Check = 400 lines
4. After Step 4 (E2E Tests): Check = 500 lines

### If Approaching Limit
- Focus on critical path tests only
- Combine similar test scenarios
- Reduce test fixture complexity
- Defer edge case testing to post-MVP
- Use table-driven tests to reduce code

## Test Requirements

### Test Matrix
| Test Type | Coverage Required | Estimated Lines |
|-----------|------------------|-----------------|
| Integration Tests | 100% happy paths | 400 lines |
| E2E Tests | Primary workflows | 100 lines |
| **Total** | - | **500 lines** |

### Test Execution Requirements

1. **Environment Prerequisites**:
   - Kind cluster with idpbuilder installed
   - Gitea running with registry enabled
   - CLI binary built from cli-commands effort
   - Docker/Buildah available locally
   - Network connectivity to gitea.local:443

2. **Validation Criteria**:
   - All tests must pass before marking effort complete
   - No flaky tests allowed
   - Clear failure diagnostics required
   - Tests must be repeatable and idempotent

3. **Performance Targets**:
   - Integration test suite: < 5 minutes
   - E2E test suite: < 10 minutes
   - Individual test timeout: 30 seconds
   - Parallel test execution where possible

## Pattern Compliance
- **Test Organization**: Use Go testing conventions
- **Test Naming**: Test{Function}{Scenario} pattern
- **Assertions**: Use testify/assert and testify/require
- **Cleanup**: Always clean up test resources
- **Isolation**: Tests must not affect each other
- **Diagnostics**: Include helpful error messages

## Success Criteria
1. ✅ All integration tests passing
2. ✅ All E2E tests passing
3. ✅ CLI build command validated end-to-end
4. ✅ CLI push command validated end-to-end
5. ✅ Certificate auto-configuration tested
6. ✅ Error scenarios properly tested
7. ✅ Size under 500 lines (measured with line-counter.sh)
8. ✅ No flaky or timing-dependent tests

## Risk Mitigation
1. **Environment Risk**: Document setup requirements clearly
2. **Flaky Test Risk**: Use proper wait conditions and retries
3. **Size Risk**: Pre-identify tests to defer if needed
4. **Dependency Risk**: Verify CLI binary exists before tests

## Integration Points Summary
```
Direct Testing Target:
└── CLI Binary (from cli-commands effort)
    ├── idpbuilder build command
    └── idpbuilder push command

Indirect Validation:
├── Phase 1 Certificate Components (via CLI)
│   ├── Certificate extraction
│   ├── Trust store configuration
│   └── Fallback handling
└── Wave 1 Build/Registry (via CLI)
    ├── Image building
    └── Registry push
```

## Test Data Requirements
```
test-data/
├── dockerfiles/
│   ├── simple.Dockerfile      # FROM alpine \n CMD ["echo", "test"]
│   ├── multistage.Dockerfile  # Multi-stage build test
│   └── invalid.Dockerfile     # Syntax error for testing
└── contexts/
    └── simple-app/
        └── main.go            # package main \n func main() { println("test") }
```

## Execution Instructions for SW Engineer

1. **Prerequisites**:
   - Ensure cli-commands effort is complete
   - Verify CLI binary is available
   - Confirm Kind cluster is running
   - Verify Gitea is accessible

2. **Implementation Order**:
   - Set up test infrastructure first
   - Implement build tests
   - Implement push tests
   - Add E2E workflow tests
   - Run line-counter.sh after each step

3. **Testing the Tests**:
   - Run tests locally during development
   - Ensure all tests pass consistently
   - Verify cleanup works properly
   - Check test execution time

## Next Steps After Implementation
1. Run all tests to ensure they pass
2. Run line-counter.sh to verify size compliance
3. Code review by Code Reviewer agent
4. Fix any identified issues
5. Re-run tests after fixes
6. Commit and push to branch
7. Signal Wave 2 completion

---

**Document Version**: 1.0
**Created**: 2025-08-30
**Created By**: Code Reviewer Agent
**State**: EFFORT_PLAN_CREATION
**Status**: READY_FOR_IMPLEMENTATION