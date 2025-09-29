# Effort Implementation Plan: Unit Test Framework

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Effort**: 1.1.2 - TDD - Create unit test framework
**Branch**: `phase1/wave1/unit-test-framework`
**Base Branch**: `phase1/wave1/analyze-existing-structure` (CASCADE from E1.1.1)
**Base Branch Reason**: This is Wave 1 Effort 2, cascading from E1.1.1 per R501/R509 requirements
**Can Parallelize**: Yes
**Parallel With**: E1.1.1, E1.1.3
**Size Estimate**: 450 lines (MUST be <800)
**Dependencies**: None (but builds on E1.1.1 structure analysis)
**Dependent Efforts**: E1.1.3 (integration tests), E1.2.1, E1.2.2, E1.2.3 (all Wave 2 efforts)
**Atomic PR**: ✅ This effort = ONE PR to main (R220 REQUIREMENT)

## 📋 Source Information
**Wave Plan**: Phase 1 Wave 1 Implementation (from project plan)
**Effort Section**: Effort 1.1.2
**Created By**: Code Reviewer Agent
**Date**: 2025-09-29
**Extracted**: 2025-09-29T05:17:20Z

## 🔴 BASE BRANCH VALIDATION (R337 MANDATORY)
**The orchestrator-state.json is the SOLE SOURCE OF TRUTH for base branches!**
- Base branch MUST be `phase1/wave1/analyze-existing-structure`
- Base branch MUST match what's in orchestrator-state.json (validated)
- Reason: CASCADE pattern - E1.1.2 builds on E1.1.1's analysis
- Orchestrator MUST record this in state file before creating infrastructure

## 🚀 Parallelization Context
**Can Parallelize**: Yes
**Parallel With**: E1.1.1-analyze-existing-structure, E1.1.3-integration-test-setup
**Blocking Status**: No - this effort doesn't block others in Wave 1
**Parallel Group**: Wave 1 Parallel Group (all 3 efforts)
**Orchestrator Guidance**: Can spawn immediately with E1.1.1 and E1.1.3 (R151: <5s delta)

## 🚨 EXPLICIT SCOPE DEFINITION (R311 MANDATORY)

### IMPLEMENT EXACTLY (BE SPECIFIC!)

#### Functions to Create (EXACTLY 8 - NO MORE)
```go
1. NewMockRegistry() *MockRegistry                      // ~40 lines - Creates mock registry for testing
2. NewMockAuthTransport(username, password) Transport   // ~30 lines - Mock auth transport
3. SetupTestFixtures(t *testing.T) *TestFixtures       // ~60 lines - Initialize test environment
4. CreateTestImage(name, tag string) v1.Image          // ~50 lines - Generate test OCI images
5. AssertPushSucceeds(t *testing.T, img v1.Image)      // ~40 lines - Verify push operations
6. AssertAuthRequired(t *testing.T, endpoint string)   // ~35 lines - Test auth scenarios
7. MockInsecureTransport() http.RoundTripper           // ~25 lines - Mock insecure cert handling
8. CleanupTestFixtures(fixtures *TestFixtures)         // ~20 lines - Cleanup test resources
// STOP HERE - DO NOT ADD MORE FUNCTIONS
```

#### Types/Structs to Define (EXACTLY 5)
```go
// Type 1: Mock registry server
type MockRegistry struct {
    Server     *httptest.Server  // Test server instance
    Images     map[string][]byte // Stored images
    AuthConfig *AuthConfig       // Optional auth configuration
    // EXACTLY these fields, NO methods in this effort
}

// Type 2: Authentication config
type AuthConfig struct {
    Username string  // Expected username
    Password string  // Expected password
    Required bool    // Whether auth is required
    // NO additional fields
}

// Type 3: Test fixtures container
type TestFixtures struct {
    Registry *MockRegistry       // Mock registry instance
    Client   *http.Client        // Test HTTP client
    TempDir  string              // Temporary test directory
    // ONLY these fields
}

// Type 4: Push test case
type PushTestCase struct {
    Name     string              // Test case name
    Image    v1.Image            // Image to push
    WantErr  bool                // Expected error
    ErrorMsg string              // Expected error message
    // NO additional fields
}

// Type 5: Mock transport for auth testing
type MockAuthTransport struct {
    Username string              // Auth username
    Password string              // Auth password
    Base     http.RoundTripper   // Base transport
    // EXACTLY these fields
}
```

### 🛑 DO NOT IMPLEMENT (SCOPE BOUNDARIES)

**EXPLICITLY FORBIDDEN IN THIS EFFORT:**
- ❌ DO NOT implement actual registry client (Wave 2)
- ❌ DO NOT add real OCI push functionality (Wave 2)
- ❌ DO NOT create integration with actual registries (E1.1.3)
- ❌ DO NOT implement retry logic (Wave 2)
- ❌ DO NOT add progress indicators (Wave 2)
- ❌ DO NOT create CLI command structures (E1.2.1)
- ❌ DO NOT implement real authentication (E1.2.2)
- ❌ DO NOT add comprehensive error handling (Wave 2)
- ❌ DO NOT write performance tests or benchmarks
- ❌ DO NOT implement rate limiting stubs

### 📊 REALISTIC SIZE CALCULATION

```
Component Breakdown:
- Mock registry struct & methods:        40 lines
- Auth transport mock:                   30 lines
- Test fixtures setup:                   60 lines
- Test image creation:                   50 lines
- Push assertion helper:                 40 lines
- Auth assertion helper:                 35 lines
- Insecure transport mock:              25 lines
- Cleanup function:                      20 lines
- Type definitions (5 × 15):            75 lines
- Basic test file setup:                 30 lines
- Sample test cases (3 × 15):           45 lines

TOTAL ESTIMATE: 450 lines (must be <800)
BUFFER: 350 lines for unforeseen needs
```

## 🔴🔴🔴 PRE-PLANNING RESEARCH RESULTS (R374 MANDATORY) 🔴🔴🔴

### Existing Interfaces Found
| Interface | Location | Signature | Must Implement |
|-----------|----------|-----------|----------------|
| testing.TB | Go standard library | Standard test interface | YES - use in helpers |
| http.RoundTripper | Go standard library | RoundTrip(req) (resp, error) | YES - for transport mocks |

### Existing Implementations to Reuse
| Component | Location | Purpose | How to Use |
|-----------|----------|---------|------------|
| testify/assert | E1.1.1 go.mod v1.9.0 | Test assertions | Import and use for test validations |
| httptest | Go standard library | HTTP test server | Use for mock registry server |
| testing.T | Go standard library | Test context | Pass to all test functions |

### APIs Already Defined
| API | Method | Signature | Notes |
|-----|--------|-----------|-------|
| N/A | N/A | N/A | No push APIs exist yet - Wave 2 will create them |

### FORBIDDEN DUPLICATIONS (R373)
- ❌ DO NOT recreate testify assertion functions (use v1.9.0 from E1.1.1)
- ❌ DO NOT implement custom HTTP test server (use httptest)
- ❌ DO NOT create alternative test runners (use go test)

### REQUIRED INTEGRATIONS (R373)
- ✅ MUST use testify v1.9.0 from E1.1.1's go.mod (R381 version lock)
- ✅ MUST follow existing test patterns from E1.1.1 test files
- ✅ MUST use standard Go testing package conventions

## 📁 Files to Create

### Primary Implementation Files
```yaml
new_files:
  - path: pkg/phase1/wave1/test/push/mock_registry.go
    lines: ~200 MAX
    purpose: Mock OCI registry for unit testing
    contains:
      - NewMockRegistry function
      - MockRegistry type
      - AuthConfig type
      - MockAuthTransport type

  - path: pkg/phase1/wave1/test/push/test_helpers.go
    lines: ~150 MAX
    purpose: Test helper functions and fixtures
    contains:
      - SetupTestFixtures function
      - CreateTestImage function
      - CleanupTestFixtures function
      - TestFixtures type

  - path: pkg/phase1/wave1/test/push/assertions.go
    lines: ~100 MAX
    purpose: Test assertion helpers
    contains:
      - AssertPushSucceeds function
      - AssertAuthRequired function
      - PushTestCase type
```

### Test Files
```yaml
test_files:
  - path: pkg/phase1/wave1/test/push/framework_test.go
    lines: ~100 MAX
    coverage_target: 80%
    test_functions:
      - TestMockRegistryCreation  # ~25 lines
      - TestAuthTransport         # ~25 lines
      - TestImageCreation         # ~25 lines
      - TestCleanup              # ~25 lines
      # NO edge cases, NO benchmarks
```

## 📦 Files to Import/Reuse

### From Previous Efforts (This Wave)
```yaml
this_wave_imports:
  - source: go.mod
    from_effort: E1.1.1-analyze-existing-structure
    usage: Use testify v1.9.0 for assertions (R381 version lock)

  - source: Test patterns
    from_effort: E1.1.1-analyze-existing-structure
    usage: Follow established testing conventions
```

### From Standard Library
```yaml
standard_imports:
  - source: net/http/httptest
    usage: Mock HTTP server for registry

  - source: testing
    usage: Go test framework

  - source: io
    usage: Handle image data streams

  - source: encoding/json
    usage: Mock registry responses
```

## 🔗 Dependencies

### Effort Dependencies
- **Must Complete First**: None (can run in parallel)
- **Can Run in Parallel With**: E1.1.1-analyze-existing-structure, E1.1.3-integration-test-setup
- **Blocks**: All Wave 2 efforts need this test framework

### Technical Dependencies
- Go standard library (testing, httptest, net/http)
- testify v1.9.0 (from E1.1.1's go.mod - DO NOT UPDATE per R381)
- No external OCI libraries yet (Wave 2 will add go-containerregistry)

## 🔴 ATOMIC PR REQUIREMENTS (R220 - SUPREME LAW)

### 🔴🔴🔴 PARAMOUNT: Independent Mergeability (R307) 🔴🔴🔴
**This effort MUST be mergeable at ANY time, even YEARS later:**
- ✅ Must compile when merged alone to main
- ✅ Must NOT break any existing functionality
- ✅ Test framework is self-contained
- ✅ No external dependencies on Wave 2 work
- ✅ Tests can be skipped if push command doesn't exist yet

### Feature Flags for This Effort
```yaml
feature_flags:
  - flag: "PUSH_TEST_FRAMEWORK_ENABLED"
    location: "pkg/phase1/wave1/test/push/config.go"
    default: false
    purpose: "Enable push test framework when push command exists"
    activation: "Set true when Wave 2 push implementation begins"
```

### 🚨🚨🚨 R355 PRODUCTION READY CODE (SUPREME LAW #5) 🚨🚨🚨

**ALL CODE MUST BE PRODUCTION READY - NO EXCEPTIONS**

#### ❌ ABSOLUTELY FORBIDDEN:
- NO stubs returning "not implemented"
- NO TODO/FIXME markers in code
- NO hardcoded test credentials
- NO panic("not implemented") patterns
- NO placeholder test data

#### ✅ REQUIRED PATTERNS:
```go
// ❌ WRONG - Stub implementation
func CreateTestImage(name string) v1.Image {
    panic("not implemented") // AUTOMATIC FAILURE
}

// ✅ CORRECT - Full mock implementation
func CreateTestImage(name, tag string) v1.Image {
    // Create actual test image with manifest
    manifest := v1.Manifest{
        SchemaVersion: 2,
        MediaType:     types.DockerManifestSchema2,
        // ... complete manifest
    }
    return &testImage{manifest: manifest}
}
```

### PR Mergeability Checklist
- [ ] PR can merge to main independently
- [ ] Build passes with just this PR
- [ ] All tests pass in isolation
- [ ] Test framework is complete and functional
- [ ] No dependency on Wave 2 implementations
- [ ] Feature flag controls activation
- [ ] No breaking changes to existing code

## 🔴 MANDATORY ADHERENCE CHECKPOINTS (R311)

### Before Starting:
```bash
echo "EFFORT SCOPE LOCKED:"
echo "✓ Functions: EXACTLY 8 (listed above)"
echo "✓ Types: EXACTLY 5 (MockRegistry, AuthConfig, TestFixtures, PushTestCase, MockAuthTransport)"
echo "✓ Test Files: EXACTLY 1 framework test file"
echo "✓ Tests: EXACTLY 4 basic framework tests"
echo "✗ Real registry: NONE (E1.1.3 does that)"
echo "✗ Actual push: NONE (Wave 2)"
echo "✗ CLI integration: NONE (E1.2.1)"
```

## 📝 Implementation Instructions

### Step-by-Step Guide
1. **Scope Acknowledgment**
   - Read and acknowledge DO NOT IMPLEMENT section
   - Count: EXACTLY 8 functions, 5 types, 4 tests
   - Create .scope-acknowledgment file

2. **Implementation Order**
   - Start with mock_registry.go - define MockRegistry and AuthConfig types
   - Implement test_helpers.go - TestFixtures and setup functions
   - Add assertions.go - test assertion helpers
   - Write framework_test.go - validate the framework itself works

3. **Key Implementation Details**
   ```go
   // Mock registry implementation
   type MockRegistry struct {
       Server     *httptest.Server
       Images     map[string][]byte
       AuthConfig *AuthConfig
   }

   func NewMockRegistry() *MockRegistry {
       reg := &MockRegistry{
           Images: make(map[string][]byte),
       }
       // Set up httptest server with handlers
       reg.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           // Handle registry API endpoints
       }))
       return reg
   }
   ```

4. **Integration Points**
   - Use testify v1.9.0 from E1.1.1's go.mod
   - Follow test patterns from E1.1.1's existing tests
   - Prepare for Wave 2 push command integration

## ✅ Test Requirements

### Coverage Requirements
- **Minimum Coverage**: 80%
- **Critical Paths**: Mock creation, auth handling, cleanup
- **Framework Tests**: Validate test helpers work correctly

### Test Categories
```yaml
required_tests:
  unit_tests:
    - Mock registry creation and shutdown
    - Auth transport configuration
    - Test image generation
    - Fixture cleanup

  validation_tests:
    - Verify mock responds to registry API calls
    - Confirm auth headers are checked when configured
    - Ensure cleanup removes all test artifacts
```

## 📏 Size Constraints
**Target Size**: 450 lines
**Maximum Size**: 800 lines (HARD LIMIT)
**Current Size**: 0 lines

### Size Monitoring Protocol
```bash
# Check size after each component
cd /home/vscode/workspaces/idpbuilder-push-oci/efforts/phase1/wave1/E1.1.2-unit-test-framework
PROJECT_ROOT=/home/vscode/workspaces/idpbuilder-push-oci
$PROJECT_ROOT/tools/line-counter.sh

# If approaching 700 lines:
# 1. Alert Code Reviewer
# 2. Stop adding features
# 3. Focus on completing current functionality
```

## 🏁 Completion Criteria

### Implementation Checklist
- [ ] Mock registry with configurable auth
- [ ] Test fixture setup and teardown
- [ ] Test image creation utilities
- [ ] Assertion helpers for push operations
- [ ] Framework validation tests
- [ ] Size verified under 800 lines

### Quality Checklist
- [ ] Test coverage ≥80%
- [ ] All tests passing
- [ ] No linting errors
- [ ] No hardcoded values
- [ ] Feature flag implemented

### Documentation Checklist
- [ ] Code comments for complex mock logic
- [ ] API documentation for exported functions
- [ ] Usage examples in comments
- [ ] Work log updated with progress

### Review Checklist
- [ ] Self-review completed
- [ ] Code committed and pushed
- [ ] Ready for Code Reviewer assessment
- [ ] No blocking issues

## ⚠️ Important Notes

### Parallelization Reminder
- This effort can run simultaneously with E1.1.1 and E1.1.3
- No dependencies between Wave 1 efforts
- Ensure no shared state with parallel efforts
- R151: All parallel spawns must have <5s timestamp delta

### Common Pitfalls to Avoid (R311 ENFORCEMENT)
1. **SCOPE CREEP**: Don't implement real registry client (Wave 2)
2. **OVER-ENGINEERING**: Keep mocks simple and focused
3. **VERSION UPDATES**: Use testify v1.9.0 - DO NOT update (R381)
4. **REAL IMPLEMENTATION**: This is test framework ONLY
5. **INTEGRATION TESTS**: Those go in E1.1.3, not here
6. **CLI MOCKING**: Wave 2 handles command structure

### Success Criteria Checklist
- [ ] Implemented EXACTLY 8 functions (no more)
- [ ] Created EXACTLY 5 types (no more)
- [ ] Wrote EXACTLY 4 framework tests (no more)
- [ ] Total lines under 800
- [ ] NO real registry interactions
- [ ] NO actual push implementation
- [ ] Followed all scope boundaries

## 📚 References

### Source Documents
- [Master Implementation Plan](/home/vscode/workspaces/idpbuilder-push-oci/IMPLEMENTATION-PLAN.md)
- [E1.1.1 Structure Analysis](../E1.1.1-analyze-existing-structure/)
- [Project Push Requirements](/home/vscode/workspaces/idpbuilder-push-oci/ARCHITECT-PROMPT-IDPBUILDER-OCI.md)

### Code Examples
- [E1.1.1 Test Patterns](../E1.1.1-analyze-existing-structure/pkg/kind/cluster_test.go)
- [Go httptest Examples](https://pkg.go.dev/net/http/httptest)

### Standards
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [testify Documentation](https://github.com/stretchr/testify)

## 🔴🔴🔴 CASCADE VALIDATION (R501/R509) 🔴🔴🔴

**Branch Infrastructure Confirmation:**
- Current Branch: phase1/wave1/unit-test-framework
- Base Branch: phase1/wave1/analyze-existing-structure (E1.1.1)
- Integration Target: phase1-wave1-integration
- CASCADE verified: E1.1.2 correctly cascades from E1.1.1

**R509 Compliance:**
- ✅ Not first effort in P1W1 (E1.1.1 is first)
- ✅ Correctly based on previous effort (E1.1.1)
- ✅ Will be base for E1.1.3
- ✅ Follows cascade pattern

---

**Remember**: This is a TEST FRAMEWORK ONLY effort. No actual push implementation. The framework will be used by Wave 2 efforts when they implement the real push functionality. Stay within the 450-line estimate and focus on creating robust, reusable test utilities.