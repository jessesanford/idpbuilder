# Effort Implementation Plan: E1.1.3 - Integration Test Setup

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Effort**: E1.1.3 - Integration Test Setup
**Branch**: `phase1/wave1/integration-test-setup`
**Base Branch**: `phase1/wave1/unit-test-framework` (CASCADE from E1.1.2 per R501)
**Base Branch Reason**: This is effort #3 in phase1/wave1, cascading from effort #2 per R308 incremental strategy
**Can Parallelize**: Yes (per orchestrator state analysis)
**Parallel With**: E1.1.1, E1.1.2 (all Wave 1 efforts can run in parallel)
**Size Estimate**: 650 lines (MUST be <800)
**Dependencies**: None (can start immediately)
**Dependent Efforts**: E1.2.3 (Image push operations will use integration tests)
**Atomic PR**: ✅ This effort = ONE PR to main (R220 REQUIREMENT)

## 📋 Source Information
**Wave Plan**: Phase 1, Wave 1 - Project Analysis & Test Infrastructure
**Effort Section**: Effort 1.1.3
**Created By**: Code Reviewer Agent
**Date**: 2025-09-29
**Extracted**: 2025-09-29T05:20:18Z

## 🔴 BASE BRANCH VALIDATION (R337 MANDATORY)
**The orchestrator-state.json is the SOLE SOURCE OF TRUTH for base branches!**
- Base branch: `phase1/wave1/unit-test-framework` ✅ VERIFIED in orchestrator-state.json
- This branches from E1.1.2, NOT from main or E1.1.1
- Reason: CASCADE pattern per R501 - each effort builds on previous in wave
- Integration branch: `phase1-wave1-integration` (for final wave merge)

## 🚀 Parallelization Context
**Can Parallelize**: Yes
**Parallel With**: E1.1.1-analyze-existing-structure, E1.1.2-unit-test-framework
**Blocking Status**: Non-blocking - can run simultaneously with other Wave 1 efforts
**Parallel Group**: Wave 1 Group (all 3 efforts)
**Orchestrator Guidance**: Spawn immediately with E1.1.1 and E1.1.2 (R151: <5s timestamp delta)

## 🚨 EXPLICIT SCOPE DEFINITION (R311 MANDATORY)

### IMPLEMENT EXACTLY (BE SPECIFIC!)

#### Functions to Create (EXACTLY 8 - NO MORE)
```go
1. SetupTestRegistry(t *testing.T) (*testcontainers.Container, string)  // ~80 lines - Start Gitea container with registry
2. CreateTestCluster(t *testing.T) (*ClusterInfo, error)                // ~60 lines - Create idpbuilder cluster
3. CleanupTestCluster(t *testing.T, cluster *ClusterInfo)              // ~40 lines - Cleanup idpbuilder cluster
4. PushTestImage(registry string, image string) error                   // ~50 lines - Push test image to registry
5. PullTestImage(registry string, image string) error                   // ~50 lines - Pull test image from registry
6. GenerateTestCredentials() (*Credentials, error)                      // ~30 lines - Generate test auth credentials
7. VerifyImageInRegistry(registry string, image string) bool           // ~40 lines - Verify image exists in registry
8. SetupInsecureCertTest() (*tls.Config, error)                       // ~50 lines - Setup self-signed cert testing
// STOP HERE - DO NOT ADD MORE FUNCTIONS
```

#### Types/Structs to Define (EXACTLY 3)
```go
// Type 1: Cluster information
type ClusterInfo struct {
    Name      string     // Cluster name
    Namespace string     // Target namespace
    Context   string     // Kubeconfig context
    Cleanup   func()     // Cleanup function
    // EXACTLY these fields, NO additional fields
}

// Type 2: Test credentials
type Credentials struct {
    Username string     // Registry username
    Password string     // Registry password
    Token    string     // Optional token
    // NO additional fields
}

// Type 3: Test configuration
type IntegrationTestConfig struct {
    RegistryURL    string         // Registry endpoint
    InsecureMode   bool          // Allow insecure connections
    TestImagePath  string        // Path to test image
    Timeout        time.Duration // Test timeout
    // NO additional fields or methods in this effort
}
```

#### Test Scenarios to Create (EXACTLY 5 integration tests)
```go
// Integration test files only - actual push command comes in E1.2.3
TestIntegration_RegistryConnection    // ~60 lines - Test registry connectivity
TestIntegration_AuthenticationFlow    // ~70 lines - Test auth with credentials
TestIntegration_InsecureCertHandling // ~80 lines - Test self-signed cert support
TestIntegration_ImagePushPull        // ~90 lines - Test push/pull operations
TestIntegration_ClusterLifecycle    // ~70 lines - Test cluster create/destroy
// NO additional test scenarios
```

### 🛑 DO NOT IMPLEMENT (SCOPE BOUNDARIES)

**EXPLICITLY FORBIDDEN IN THIS EFFORT:**
- ❌ DO NOT implement the actual `idpbuilder push` command (E1.2.1)
- ❌ DO NOT implement production authentication logic (E1.2.2)
- ❌ DO NOT implement the actual push operation (E1.2.3)
- ❌ DO NOT add rate limiting or retry logic (future effort)
- ❌ DO NOT implement progress indicators (E1.2.3)
- ❌ DO NOT create CLI command structure (E1.2.1)
- ❌ DO NOT add comprehensive error handling beyond test needs
- ❌ DO NOT implement logging framework (use basic t.Log)
- ❌ DO NOT write performance benchmarks
- ❌ DO NOT refactor existing idpbuilder code
- ❌ DO NOT create documentation (Phase 2)

### 📊 REALISTIC SIZE CALCULATION

```
Component Breakdown:
- SetupTestRegistry:              80 lines
- CreateTestCluster:              60 lines
- CleanupTestCluster:             40 lines
- PushTestImage:                  50 lines
- PullTestImage:                  50 lines
- GenerateTestCredentials:        30 lines
- VerifyImageInRegistry:          40 lines
- SetupInsecureCertTest:          50 lines
- Type definitions (3):           30 lines
- Integration tests (5):         370 lines
- Test utilities/helpers:         50 lines

TOTAL ESTIMATE: 650 lines (must be <800)
BUFFER: 150 lines for unforeseen needs
```

## 🔴🔴🔴 PRE-PLANNING RESEARCH RESULTS (R374 MANDATORY) 🔴🔴🔴

### Existing Interfaces Found
| Interface | Location | Signature | Must Implement |
|-----------|----------|-----------|----------------|
| N/A - New effort | N/A | N/A | NO |

### Existing Implementations to Reuse
| Component | Location | Purpose | How to Use |
|-----------|----------|---------|------------|
| testcontainers-go | external library | Container management | Import for Gitea setup |
| idpbuilder CLI | main branch | Cluster management | Execute via exec.Command |
| go-containerregistry | external library | OCI operations | Import for test image operations |

### APIs Already Defined
| API | Method | Signature | Notes |
|-----|--------|-----------|-------|
| N/A - Test infrastructure only | N/A | N/A | Creating new test APIs |

### FORBIDDEN DUPLICATIONS (R373)
- ❌ DO NOT duplicate any existing idpbuilder test utilities if found
- ❌ DO NOT recreate container management (use testcontainers-go)
- ❌ DO NOT implement custom OCI client (use go-containerregistry)

### REQUIRED INTEGRATIONS (R373)
- ✅ MUST use testcontainers-go for Gitea registry container
- ✅ MUST use existing idpbuilder binary for cluster operations
- ✅ MUST use go-containerregistry for image operations

## 📁 Files to Create

### Primary Implementation Files
```yaml
new_files:
  - path: pkg/integration/registry_setup.go
    lines: ~200 MAX
    purpose: Registry container setup and management
    contains:
      - SetupTestRegistry (NO additional functions)
      - ClusterInfo type
      - Basic container lifecycle management

  - path: pkg/integration/cluster_helpers.go
    lines: ~150 MAX
    purpose: idpbuilder cluster test helpers
    contains:
      - CreateTestCluster
      - CleanupTestCluster
      - Credentials type
      - GenerateTestCredentials

  - path: pkg/integration/image_helpers.go
    lines: ~140 MAX
    purpose: Image push/pull test utilities
    contains:
      - PushTestImage
      - PullTestImage
      - VerifyImageInRegistry

  - path: pkg/integration/tls_helpers.go
    lines: ~100 MAX
    purpose: TLS/insecure mode testing
    contains:
      - SetupInsecureCertTest
      - IntegrationTestConfig type
```

### Test Files
```yaml
test_files:
  - path: pkg/integration/integration_test.go
    lines: ~370 MAX
    coverage_target: N/A (these ARE the tests)
    test_functions:
      - TestIntegration_RegistryConnection    # ~60 lines
      - TestIntegration_AuthenticationFlow    # ~70 lines
      - TestIntegration_InsecureCertHandling # ~80 lines
      - TestIntegration_ImagePushPull        # ~90 lines
      - TestIntegration_ClusterLifecycle    # ~70 lines
      # NO edge cases, NO additional scenarios
```

## 📦 Files to Import/Reuse

### From Previous Efforts (This Wave)
```yaml
this_wave_imports:
  - source: N/A - First test infrastructure effort
    from_effort: N/A
    usage: N/A
```

### External Dependencies
```yaml
external_imports:
  - source: github.com/testcontainers/testcontainers-go
    usage: Container management for Gitea registry

  - source: github.com/google/go-containerregistry
    usage: OCI image operations in tests

  - source: testing package
    usage: Standard Go testing framework
```

## 🔗 Dependencies

### Effort Dependencies
- **Must Complete First**: None - can start immediately
- **Can Run in Parallel With**: E1.1.1, E1.1.2
- **Blocks**: E1.2.3 (Image push operations will use these integration tests)

### Technical Dependencies
- testcontainers-go library for container management
- go-containerregistry for OCI operations
- Existing idpbuilder binary for cluster operations

## 🔴 ATOMIC PR REQUIREMENTS (R220 - SUPREME LAW)

### 🔴🔴🔴 PARAMOUNT: Independent Mergeability (R307) 🔴🔴🔴
**This effort MUST be mergeable at ANY time, even YEARS later:**
- ✅ Integration tests are self-contained
- ✅ Tests can be skipped if dependencies missing
- ✅ Does NOT break any existing functionality
- ✅ Tests use build tags for conditional compilation
- ✅ Gracefully handles missing Docker/Kubernetes

### Feature Flags for This Effort
```yaml
build_tags:
  - tag: "integration"
    purpose: "Isolate integration tests from unit tests"
    usage: "go test -tags=integration ./..."
    default: "Tests excluded without tag"
```

### 🚨🚨🚨 R355 PRODUCTION READY CODE (SUPREME LAW #5) 🚨🚨🚨

**ALL CODE MUST BE PRODUCTION READY - NO EXCEPTIONS**

#### ❌ ABSOLUTELY FORBIDDEN:
- NO stub test implementations
- NO placeholder assertions
- NO hardcoded registry URLs
- NO static test credentials
- NO TODO/FIXME markers in code
- NO panic("not implemented") in tests
- NO skipped tests without proper conditionals

#### ✅ REQUIRED PATTERNS:
```go
// ❌ WRONG - Hardcoded test values
registryURL := "localhost:5000"
username := "testuser"
password := "testpass"

// ✅ CORRECT - Configuration-driven
registryURL := os.Getenv("TEST_REGISTRY_URL")
if registryURL == "" {
    registryURL = setupTestRegistry(t)
}
username := os.Getenv("TEST_USERNAME")
if username == "" {
    creds := generateTestCredentials()
    username = creds.Username
}
```

### PR Mergeability Checklist
- [ ] Tests compile without the push command existing
- [ ] Tests can run standalone
- [ ] Build passes with just this PR
- [ ] Tests gracefully skip if Docker unavailable
- [ ] No breaking changes to existing code
- [ ] Uses build tags for isolation
- [ ] Backward compatible with main

## 🔴 MANDATORY ADHERENCE CHECKPOINTS (R311)

### Before Starting:
```bash
echo "EFFORT SCOPE LOCKED:"
echo "✓ Functions: EXACTLY 8 (list them)"
echo "✓ Types: EXACTLY 3 (ClusterInfo, Credentials, IntegrationTestConfig)"
echo "✓ Integration Tests: EXACTLY 5 scenarios"
echo "✓ Test Helpers: Integration test support ONLY"
echo "✗ Push Command: NONE (E1.2.1)"
echo "✗ Production Auth: NONE (E1.2.2)"
echo "✗ Actual Push Logic: NONE (E1.2.3)"
```

### During Implementation:
```bash
# Check scope adherence after each component
cd /home/vscode/workspaces/idpbuilder-push-oci/efforts/phase1/wave1/E1.1.3-integration-test-setup
FUNC_COUNT=$(grep -c "^func [A-Z]" pkg/integration/*.go 2>/dev/null || echo 0)
if [ "$FUNC_COUNT" -gt 8 ]; then
    echo "⚠️ WARNING: Exceeding function count! Stop adding!"
fi
```

## 📝 Implementation Instructions

### Step-by-Step Guide
1. **Scope Acknowledgment**
   - Read and acknowledge DO NOT IMPLEMENT section
   - Confirm exactly 8 functions, 3 types, 5 integration tests
   - Create .scope-acknowledgment file

2. **Implementation Order**
   - Start with type definitions (ClusterInfo, Credentials, IntegrationTestConfig)
   - Implement registry_setup.go with SetupTestRegistry
   - Add cluster_helpers.go with cluster management
   - Create image_helpers.go for push/pull testing
   - Add tls_helpers.go for insecure cert testing
   - Write integration_test.go with 5 test scenarios

3. **Key Implementation Details**
   ```go
   // Use testcontainers for Gitea registry
   import "github.com/testcontainers/testcontainers-go"

   func SetupTestRegistry(t *testing.T) (*testcontainers.Container, string) {
       ctx := context.Background()
       req := testcontainers.ContainerRequest{
           Image:        "gitea/gitea:latest",
           ExposedPorts: []string{"3000/tcp"},
           WaitingFor:   wait.ForHTTP("/").WithPort("3000"),
       }
       // Start container and return URL
   }

   // Execute idpbuilder via command
   func CreateTestCluster(t *testing.T) (*ClusterInfo, error) {
       cmd := exec.Command("idpbuilder", "create", "--name", "test-cluster")
       // Execute and parse output
   }
   ```

4. **Integration Points**
   - Use testcontainers-go for container management
   - Shell out to idpbuilder binary for cluster ops
   - Use go-containerregistry for image operations
   - Keep tests isolated with build tags

## ✅ Test Requirements

### Coverage Requirements
- **Integration Tests**: 5 comprehensive scenarios
- **Test Isolation**: Use build tags
- **Cleanup**: Ensure all resources cleaned up
- **Idempotency**: Tests can run repeatedly

### Test Categories
```yaml
integration_tests:
  - Registry connectivity test
  - Authentication flow validation
  - Insecure certificate handling
  - Image push/pull operations
  - Cluster lifecycle management
```

## 📏 Size Constraints
**Target Size**: 650 lines
**Maximum Size**: 800 lines (HARD LIMIT)
**Current Size**: 0 lines

### Size Monitoring Protocol
```bash
# Check size every ~200 lines
cd /home/vscode/workspaces/idpbuilder-push-oci/efforts/phase1/wave1/E1.1.3-integration-test-setup
PROJECT_ROOT="/home/vscode/workspaces/idpbuilder-push-oci"
$PROJECT_ROOT/tools/line-counter.sh

# If approaching 700 lines:
# 1. Alert Code Reviewer
# 2. Prepare for potential split
# 3. Focus on completing current functionality
```

## 🏁 Completion Criteria

### Implementation Checklist
- [ ] All 8 functions implemented
- [ ] All 3 types defined
- [ ] All 5 integration tests written
- [ ] Size verified under 800 lines
- [ ] Build tags properly used
- [ ] No scope creep detected

### Quality Checklist
- [ ] Tests can run standalone
- [ ] Proper cleanup implemented
- [ ] No hardcoded values
- [ ] Error handling in place
- [ ] Tests are idempotent

### Documentation Checklist
- [ ] Code comments for complex logic
- [ ] Test descriptions clear
- [ ] Build tag usage documented
- [ ] No TODO/FIXME markers

### Review Checklist
- [ ] Self-review completed
- [ ] Code committed and pushed
- [ ] Ready for Code Reviewer assessment
- [ ] No blocking issues

## 📊 Progress Tracking

### Work Log
```markdown
## 2025-09-29 - Session 1
- Created implementation plan
- Defined scope boundaries
- Listed exact functions and types
- Plan ready for SW Engineer

[SW Engineer to continue updating during implementation]
```

## ⚠️ Important Notes

### Parallelization Reminder
- This effort can run simultaneously with E1.1.1 and E1.1.2
- No dependencies on other Wave 1 efforts
- Ensure no shared state with parallel efforts
- Can start immediately

### Common Pitfalls to Avoid (R311 ENFORCEMENT)
1. **SCOPE CREEP**: Adding "helpful" test utilities = AUTOMATIC FAILURE
2. **OVER-ENGINEERING**: Making tests "production-ready" = 3-5X overrun
3. **ASSUMPTIONS**: Implementing push command logic = VIOLATION
4. **Size Limit**: Monitor continuously with line-counter.sh
5. **Dependencies**: Don't depend on E1.1.1 or E1.1.2 outputs
6. **Test Coverage**: Integration tests only - no unit tests needed yet
7. **Isolation**: Use build tags to isolate integration tests
8. **Parallelization**: Don't create dependencies on parallel efforts

### Success Criteria Checklist
- [ ] Read and acknowledged DO NOT IMPLEMENT section
- [ ] Implemented EXACTLY 8 functions (no more)
- [ ] Created EXACTLY 3 types (no more)
- [ ] Wrote EXACTLY 5 integration tests (no more)
- [ ] Total lines under 800
- [ ] NO push command implementation
- [ ] NO production authentication
- [ ] Followed all scope boundaries

## 📚 References

### Source Documents
- [Master Implementation Plan](/home/vscode/workspaces/idpbuilder-push-oci/IMPLEMENTATION-PLAN.md)
- [Orchestrator State](/home/vscode/workspaces/idpbuilder-push-oci/orchestrator-state.json)
- [Target Repository](https://github.com/jessesanford/idpbuilder.git)

### Key Libraries
- [testcontainers-go Documentation](https://golang.testcontainers.org/)
- [go-containerregistry Documentation](https://github.com/google/go-containerregistry)
- [idpbuilder Documentation](https://github.com/cnoe-io/idpbuilder)

### Standards
- Go testing best practices
- Integration test patterns
- Container test management

---

**Remember**: This effort focuses ONLY on integration test infrastructure. The actual push command implementation comes in later efforts (E1.2.1-E1.2.3). Stay within scope boundaries!