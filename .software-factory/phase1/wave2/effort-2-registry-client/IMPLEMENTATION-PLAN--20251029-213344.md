# Registry Client Implementation Plan
## Effort 1.2.2: Registry Client Implementation

**Created**: 2025-10-29 21:33:44 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**State**: EFFORT_PLAN_CREATION
**Effort ID**: 1.2.2
**Phase**: Phase 1 - Foundation & Interfaces
**Wave**: Wave 2 - Core Package Implementations

---

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)

### R213 Metadata (IMMUTABLE - DO NOT MODIFY)

```json
{
  "effort_id": "1.2.2",
  "effort_name": "Registry Client Implementation",
  "branch_name": "idpbuilder-oci-push/phase1/wave2/effort-2-registry-client",
  "base_branch": "idpbuilder-oci-push/phase1/wave2/integration",
  "parent_wave": "WAVE_2",
  "parent_phase": "PHASE_1",
  "depends_on": [],
  "estimated_lines": 450,
  "complexity": "high",
  "can_parallelize": true,
  "parallel_with": ["1.2.1", "1.2.3", "1.2.4"]
}
```

### Key Parallelization Information

**Can Parallelize**: Yes (Wave 2 ALL 4 efforts run in parallel)

**Parallel With**:
- Effort 1.2.1 (Docker Client)
- Effort 1.2.3 (Authentication)
- Effort 1.2.4 (TLS Configuration)

**Why Parallel is Safe**:
- All Wave 1 interfaces are FROZEN (no coordination needed)
- This effort implements registry package only (no file conflicts)
- No implementation-time dependencies on other Wave 2 efforts
- Depends ONLY on Wave 1 interfaces (already complete)

---

## Overview

**Purpose**: Implement the registry client package that pushes OCI images to container registries, validates registry connectivity, builds fully-qualified image references, and classifies errors correctly.

**What This Effort Accomplishes**:
- Complete implementation of `registry.Client` interface (frozen in Wave 1)
- OCI image push operations using go-containerregistry
- Registry connectivity validation (/v2/ endpoint ping)
- Image reference construction (registry/namespace/repository:tag)
- Progress callback support for layer uploads
- Error classification (auth errors, network errors, push failures)
- Integration with auth and TLS providers

**Boundaries - OUT OF SCOPE**:
- Authentication implementation (auth package responsibility - Effort 1.2.3)
- TLS configuration implementation (tls package responsibility - Effort 1.2.4)
- Image building or modification (docker package responsibility - Effort 1.2.1)
- Registry management or administration
- Image manifest manipulation (go-containerregistry handles this)

**Estimated Implementation Time**: 6-8 hours

---

## Dependencies

### Upstream Dependencies (MUST COMPLETE BEFORE THIS EFFORT)

**Wave 1 Frozen Interfaces** (COMPLETED):
- ✅ Wave 1 Effort 2: Registry interface definition (`pkg/registry/interface.go`)
  - Location: `efforts/phase1/wave1/effort-2-registry-interface/`
  - Interface: `registry.Client` with 3 methods
  - Types: `ProgressCallback`, `ProgressUpdate`
  - Error Types: `AuthenticationError`, `NetworkError`, `RegistryUnavailableError`, `PushFailedError`, `ValidationError`

- ✅ Wave 1 Effort 3: Auth interface definition (`pkg/auth/interface.go`)
  - Location: `efforts/phase1/wave1/effort-3-auth-tls-interfaces/`
  - Interface: `auth.Provider` with 2 methods
  - Used via dependency injection in constructor

- ✅ Wave 1 Effort 3: TLS interface definition (`pkg/tls/interface.go`)
  - Location: `efforts/phase1/wave1/effort-3-auth-tls-interfaces/`
  - Interface: `tls.ConfigProvider` with 2 methods
  - Used via dependency injection in constructor

**Integration Branch** (CREATED):
- ✅ Branch: `idpbuilder-oci-push/phase1/wave2/integration`
- Contains all Wave 1 interfaces ready for use

### Downstream Dependencies (EFFORTS THAT DEPEND ON THIS)

- None (all Wave 2 efforts are parallel)
- Wave 3 CLI will use this package after Wave 2 integration

### External Library Dependencies

**Required Libraries** (already in go.mod from Wave 1):
- `github.com/google/go-containerregistry` v0.19.0
  - `pkg/name` for reference parsing
  - `pkg/v1/remote` for push operations
  - `pkg/authn` for authentication types
- `github.com/stretchr/testify` v1.10.0 (for tests)

**No New Dependencies**: All required libraries already added in Wave 1

### Internal Package Dependencies (Runtime Injection)

**This implementation USES these interfaces** (implemented by other Wave 2 efforts):
- `pkg/auth.Provider` interface (will be implemented in Effort 1.2.3)
- `pkg/tls.ConfigProvider` interface (will be implemented in Effort 1.2.4)

**Critical Understanding**:
- At IMPLEMENTATION time: Use INTERFACES from Wave 1 (type checking)
- At RUNTIME: Receive concrete implementations via constructor injection
- This effort does NOT implement auth or TLS - only USES their interfaces

---

## File Structure

### New Files to Create

#### 1. `pkg/registry/client.go` (~450 lines)

**Purpose**: Core registry client implementation

**Contents**:
- `registryClient` struct (private implementation type)
- `NewClient()` constructor with provider validation
- `Push()` method with progress callbacks and error classification
- `BuildImageReference()` reference construction with "giteaadmin" namespace
- `ValidateRegistry()` /v2/ endpoint check
- Helper functions:
  - `parseImageName()` - Extract repository and tag
  - `createProgressHandler()` - Convert ProgressCallback to v1.Update channel
  - `isAuthError()` - Classify authentication failures
  - `isNetworkError()` - Classify network connectivity failures

**Package Declaration**:
```go
// Package registry provides OCI registry push operations.
package registry
```

**Imports Required**:
```go
import (
    "context"
    "fmt"
    "net/http"
    "net/url"
    "strings"

    "github.com/google/go-containerregistry/pkg/name"
    v1 "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/remote"
)
```

### Test Files (NOT counted in line estimates per R007)

#### 2. `pkg/registry/client_test.go` (~400 lines)

**Purpose**: Comprehensive unit tests for registry client

**Test Count**: 15+ test cases

**Test Categories**:
- Constructor tests (3 tests)
- Push tests (4 tests)
- BuildImageReference tests (2 tests)
- ValidateRegistry tests (2 tests)
- Helper function tests (as needed)

**Test Coverage Target**: ≥85% (per Wave 2 Test Plan)

### Modified Files

**None** - All dependencies already in go.mod from Wave 1

### File Organization

```
pkg/registry/
├── interface.go          # (Wave 1 - DO NOT MODIFY)
├── errors.go             # (Wave 1 - DO NOT MODIFY)
├── client.go             # NEW - This effort
└── client_test.go        # NEW - This effort
```

---

## Implementation Steps

### Step 1: Create Package Structure

**Actions**:
1. Navigate to effort directory
2. Create `pkg/registry/client.go`
3. Add package declaration and imports
4. Define `registryClient` struct

**Struct Definition** (from Wave 2 Architecture):
```go
// registryClient implements the Client interface using go-containerregistry.
type registryClient struct {
    authProvider auth.Provider        // Injected auth provider
    tlsConfig    tls.ConfigProvider   // Injected TLS config
    httpClient   *http.Client          // Configured with TLS
}
```

**Verification**: File compiles, struct defined

---

### Step 2: Implement Constructor (NewClient)

**Signature** (from Wave 1 interface):
```go
func NewClient(authProvider auth.Provider, tlsConfig tls.ConfigProvider) (Client, error)
```

**Implementation Requirements** (from Wave 2 Architecture):

1. **Validate authProvider is not nil**:
   - Return `ValidationError` if nil
   - Field: "authProvider"
   - Message: "authentication provider cannot be nil"

2. **Validate tlsConfig is not nil**:
   - Return `ValidationError` if nil
   - Field: "tlsConfig"
   - Message: "TLS config provider cannot be nil"

3. **Create HTTP client with TLS config**:
   ```go
   httpClient := &http.Client{
       Transport: &http.Transport{
           TLSClientConfig: tlsConfig.GetTLSConfig(),
       },
   }
   ```

4. **Return registryClient with stored providers**:
   ```go
   return &registryClient{
       authProvider: authProvider,
       tlsConfig:    tlsConfig,
       httpClient:   httpClient,
   }, nil
   ```

**GoDoc** (mandatory):
```go
// NewClient creates a new registry client with authentication and TLS configuration.
//
// The client is configured with:
//   - Authentication provider for registry credentials
//   - TLS configuration for secure/insecure mode
//   - HTTP transport with proper timeouts
//
// Parameters:
//   - authProvider: Authentication provider (from pkg/auth)
//   - tlsConfig: TLS configuration provider (from pkg/tls)
//
// Returns:
//   - Client: Registry client interface implementation
//   - error: ValidationError if providers are invalid
//
// Example:
//   authProvider := auth.NewBasicAuthProvider("giteaadmin", "password")
//   tlsProvider := tls.NewConfigProvider(insecure)
//   client, err := registry.NewClient(authProvider, tlsProvider)
//   if err != nil {
//       return fmt.Errorf("failed to create registry client: %w", err)
//   }
```

**Tests to Write** (before implementing):
- TC-REGISTRY-IMPL-001: NewClient success with valid providers
- TC-REGISTRY-IMPL-002: NewClient fails with nil auth provider
- TC-REGISTRY-IMPL-003: NewClient fails with nil TLS provider

**Verification**: Constructor validates inputs, creates HTTP client correctly

---

### Step 3: Implement Push Method

**Signature** (from Wave 1 interface):
```go
func (c *registryClient) Push(ctx context.Context, image v1.Image, targetRef string,
                              progressCallback ProgressCallback) error
```

**Implementation Requirements** (from Wave 2 Architecture):

1. **Parse target reference**:
   ```go
   ref, err := name.ParseReference(targetRef)
   if err != nil {
       return &PushFailedError{
           TargetRef: targetRef,
           Cause:     fmt.Errorf("invalid reference: %w", err),
       }
   }
   ```

2. **Get authenticator from provider**:
   ```go
   authenticator, err := c.authProvider.GetAuthenticator()
   if err != nil {
       return &AuthenticationError{
           Registry: targetRef,
           Cause:    err,
       }
   }
   ```

3. **Configure remote options**:
   ```go
   options := []remote.Option{
       remote.WithAuth(authenticator),
       remote.WithTransport(c.httpClient.Transport),
       remote.WithContext(ctx),
   }

   // Add progress callback if provided
   if progressCallback != nil {
       options = append(options, remote.WithProgress(createProgressHandler(progressCallback)))
   }
   ```

4. **Call remote.Write()**:
   ```go
   err = remote.Write(ref, image, options...)
   ```

5. **Classify errors** (CRITICAL - proper error types):
   ```go
   if err != nil {
       // Check for authentication failures
       if isAuthError(err) {
           return &AuthenticationError{
               Registry: targetRef,
               Cause:    err,
           }
       }
       // Check for network connectivity issues
       if isNetworkError(err) {
           return &NetworkError{
               Registry: targetRef,
               Cause:    err,
           }
       }
       // Other push failures
       return &PushFailedError{
           TargetRef: targetRef,
           Cause:     err,
       }
   }

   return nil
   ```

**GoDoc** (mandatory):
```go
// Push pushes an OCI image to the specified registry with optional progress reporting.
//
// This method:
//   1. Parses the target reference
//   2. Gets authenticator from auth provider
//   3. Configures remote options (auth, transport, progress)
//   4. Calls go-containerregistry's remote.Write()
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - image: OCI v1.Image to push (from Docker client)
//   - targetRef: Fully qualified image reference
//                Format: "registry-host/namespace/repository:tag"
//   - progressCallback: Optional callback for progress updates (can be nil)
//
// Returns:
//   - error: AuthenticationError if credentials invalid (401/403),
//            NetworkError if registry unreachable,
//            PushFailedError if layer upload or manifest push fails
//
// Example:
//   err := client.Push(ctx, image, "registry.io/repo/myapp:latest", func(update ProgressUpdate) {
//       fmt.Printf("Layer %s: %d/%d bytes\n", update.LayerDigest, update.BytesPushed, update.LayerSize)
//   })
```

**Tests to Write** (before implementing):
- TC-REGISTRY-IMPL-004: Push success with progress callbacks
- TC-REGISTRY-IMPL-005: Push returns AuthenticationError for 401/403
- TC-REGISTRY-IMPL-006: Push returns NetworkError for unreachable registry
- TC-REGISTRY-IMPL-007: Push returns PushFailedError for invalid target reference

**Verification**: Pushes images, classifies errors correctly, callbacks work

---

### Step 4: Implement BuildImageReference Method

**Signature** (from Wave 1 interface):
```go
func (c *registryClient) BuildImageReference(registryURL, imageName string) (string, error)
```

**Implementation Requirements** (from Wave 2 Architecture):

1. **Parse registry URL**:
   ```go
   parsedURL, err := url.Parse(registryURL)
   if err != nil {
       return "", &ValidationError{
           Field:   "registryURL",
           Message: fmt.Sprintf("invalid registry URL: %v", err),
       }
   }
   ```

2. **Extract host:port**:
   ```go
   registryHost := parsedURL.Host
   if registryHost == "" {
       return "", &ValidationError{
           Field:   "registryURL",
           Message: "registry URL must include host",
       }
   }
   ```

3. **Parse image name** (use helper):
   ```go
   repository, tag := parseImageName(imageName)
   if repository == "" {
       return "", &ValidationError{
           Field:   "imageName",
           Message: "image name cannot be empty",
       }
   }
   ```

4. **Default tag if not specified**:
   ```go
   if tag == "" {
       tag = "latest"
   }
   ```

5. **Build full reference with "giteaadmin" namespace**:
   ```go
   namespace := "giteaadmin"
   fullRef := fmt.Sprintf("%s/%s/%s:%s", registryHost, namespace, repository, tag)
   return fullRef, nil
   ```

**GoDoc** (mandatory):
```go
// BuildImageReference constructs a fully qualified registry image reference.
//
// This method:
//   1. Parses the registry URL to extract host:port
//   2. Parses image name to extract repository and tag
//   3. Constructs full reference: registry/namespace/repository:tag
//   4. Uses "giteaadmin" as default namespace for Gitea registries
//
// Parameters:
//   - registryURL: Base registry URL
//                  Examples: "https://gitea.cnoe.localtest.me:8443"
//                           "https://registry.io"
//   - imageName: Image name with optional tag
//                Examples: "myapp:latest", "myapp", "myapp:v1.0.0"
//
// Returns:
//   - string: Fully qualified image reference
//   - error: ValidationError if registry URL or image name is invalid
//
// Example:
//   ref, err := client.BuildImageReference(
//       "https://gitea.cnoe.localtest.me:8443",
//       "myapp:latest",
//   )
//   // ref = "gitea.cnoe.localtest.me:8443/giteaadmin/myapp:latest"
```

**Tests to Write** (before implementing):
- TC-REGISTRY-IMPL-008: BuildImageReference constructs correct references
  - Test case 1: `https://gitea.cnoe.localtest.me:8443` + `myapp:latest` → `gitea.cnoe.localtest.me:8443/giteaadmin/myapp:latest`
  - Test case 2: `https://registry.io` + `myapp` → `registry.io/giteaadmin/myapp:latest` (default tag)
- TC-REGISTRY-IMPL-009: BuildImageReference returns ValidationError for invalid URLs

**Verification**: Constructs references correctly, handles default tag, validates inputs

---

### Step 5: Implement ValidateRegistry Method

**Signature** (from Wave 1 interface):
```go
func (c *registryClient) ValidateRegistry(ctx context.Context, registryURL string) error
```

**Implementation Requirements** (from Wave 2 Architecture):

1. **Parse registry URL**:
   ```go
   parsedURL, err := url.Parse(registryURL)
   if err != nil {
       return &ValidationError{
           Field:   "registryURL",
           Message: fmt.Sprintf("invalid registry URL: %v", err),
       }
   }
   ```

2. **Build /v2/ endpoint URL**:
   ```go
   v2URL := fmt.Sprintf("%s://%s/v2/", parsedURL.Scheme, parsedURL.Host)
   ```

3. **Create HTTP GET request**:
   ```go
   req, err := http.NewRequestWithContext(ctx, "GET", v2URL, nil)
   if err != nil {
       return &NetworkError{
           Registry: registryURL,
           Cause:    err,
       }
   }
   ```

4. **Perform request**:
   ```go
   resp, err := c.httpClient.Do(req)
   if err != nil {
       return &NetworkError{
           Registry: registryURL,
           Cause:    err,
       }
   }
   defer resp.Body.Close()
   ```

5. **Check response status**:
   ```go
   // 200 OK = registry accessible and doesn't require auth
   // 401 Unauthorized = registry accessible and requires auth (success!)
   if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized {
       return nil
   }

   return &RegistryUnavailableError{
       Registry:   registryURL,
       StatusCode: resp.StatusCode,
   }
   ```

**GoDoc** (mandatory):
```go
// ValidateRegistry checks if the registry is reachable by pinging the /v2/ endpoint.
//
// This method performs a GET request to the registry's /v2/ endpoint to verify:
//   - Registry is accessible (network connectivity)
//   - Registry responds (service is running)
//   - Registry speaks OCI protocol (returns 200 or 401)
//
// A 401 (Unauthorized) response is considered success because it indicates
// the registry is accessible and requires authentication.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - registryURL: Registry base URL to validate
//
// Returns:
//   - error: NetworkError if unreachable,
//            RegistryUnavailableError if invalid response,
//            ValidationError if URL is malformed
//
// Example:
//   if err := client.ValidateRegistry(ctx, "https://registry.io"); err != nil {
//       return fmt.Errorf("registry validation failed: %w", err)
//   }
```

**Tests to Write** (before implementing):
- TC-REGISTRY-IMPL-010: ValidateRegistry succeeds for reachable registry (200 or 401)
- TC-REGISTRY-IMPL-011: ValidateRegistry returns NetworkError for unreachable registry

**Verification**: Validates /v2/ endpoint, accepts 200 or 401 as success

---

### Step 6: Implement Helper Functions

#### Helper 1: parseImageName

**Purpose**: Extract repository and tag from image name

**Signature**:
```go
func parseImageName(imageName string) (repository, tag string)
```

**Implementation** (from Wave 2 Architecture):
```go
func parseImageName(imageName string) (repository, tag string) {
    parts := strings.Split(imageName, ":")
    if len(parts) == 2 {
        return parts[0], parts[1]
    }
    return parts[0], ""
}
```

**Logic**:
- Split on `:` character
- If 2 parts: return (repository, tag)
- If 1 part: return (repository, "")
- Caller handles default "latest" tag

---

#### Helper 2: createProgressHandler

**Purpose**: Convert ProgressCallback to v1.Update channel for go-containerregistry

**Signature**:
```go
func createProgressHandler(callback ProgressCallback) chan v1.Update
```

**Implementation** (from Wave 2 Architecture):
```go
func createProgressHandler(callback ProgressCallback) chan v1.Update {
    updates := make(chan v1.Update, 100)
    go func() {
        for update := range updates {
            callback(ProgressUpdate{
                LayerDigest: update.Digest.String(),
                LayerSize:   update.Total,
                BytesPushed: update.Complete,
                Status:      "uploading",
            })
        }
    }()
    return updates
}
```

**Logic**:
- Create buffered channel (100 capacity for smooth progress)
- Start goroutine to convert v1.Update → ProgressUpdate
- Return channel for remote.Write() to send updates

---

#### Helper 3: isAuthError

**Purpose**: Classify authentication failures from error messages

**Signature**:
```go
func isAuthError(err error) bool
```

**Implementation** (from Wave 2 Architecture):
```go
func isAuthError(err error) bool {
    errStr := err.Error()
    return strings.Contains(errStr, "401") || strings.Contains(errStr, "403") ||
        strings.Contains(errStr, "unauthorized") || strings.Contains(errStr, "forbidden")
}
```

**Detection Patterns**:
- HTTP 401 status code
- HTTP 403 status code
- Error message contains "unauthorized"
- Error message contains "forbidden"

---

#### Helper 4: isNetworkError

**Purpose**: Classify network connectivity failures from error messages

**Signature**:
```go
func isNetworkError(err error) bool
```

**Implementation** (from Wave 2 Architecture):
```go
func isNetworkError(err error) bool {
    errStr := err.Error()
    return strings.Contains(errStr, "connection") || strings.Contains(errStr, "timeout") ||
        strings.Contains(errStr, "network")
}
```

**Detection Patterns**:
- Error message contains "connection"
- Error message contains "timeout"
- Error message contains "network"

---

### Step 7: Write Comprehensive Tests (TDD Approach)

**Test File**: `pkg/registry/client_test.go`

**Test Framework**:
- Go standard testing package
- `github.com/stretchr/testify/assert` for assertions
- `github.com/stretchr/testify/require` for fatal assertions

**Mock Implementations Needed**:

1. **Mock Auth Provider** (implements `auth.Provider`):
   ```go
   type mockAuthProvider struct {
       authenticator authn.Authenticator
       validateErr   error
   }

   func (m *mockAuthProvider) GetAuthenticator() (authn.Authenticator, error) {
       return m.authenticator, m.validateErr
   }

   func (m *mockAuthProvider) ValidateCredentials() error {
       return m.validateErr
   }
   ```

2. **Mock TLS Provider** (implements `tls.ConfigProvider`):
   ```go
   type mockTLSProvider struct {
       config   *tls.Config
       insecure bool
   }

   func (m *mockTLSProvider) GetTLSConfig() *tls.Config {
       return m.config
   }

   func (m *mockTLSProvider) IsInsecure() bool {
       return m.insecure
   }
   ```

**Test Categories and Cases** (from Wave 2 Test Plan):

**A. Constructor Tests** (3 tests):
- TC-REGISTRY-IMPL-001: NewClient success with valid providers
- TC-REGISTRY-IMPL-002: NewClient fails with nil auth provider
- TC-REGISTRY-IMPL-003: NewClient fails with nil TLS provider

**B. Push Tests** (4 tests):
- TC-REGISTRY-IMPL-004: Push success with progress callbacks
- TC-REGISTRY-IMPL-005: Push returns AuthenticationError for 401/403
- TC-REGISTRY-IMPL-006: Push returns NetworkError for unreachable registry
- TC-REGISTRY-IMPL-007: Push returns PushFailedError for invalid target reference

**C. BuildImageReference Tests** (2 tests):
- TC-REGISTRY-IMPL-008: BuildImageReference constructs correct references
  - Subtest 1: Full URL with port and tag
  - Subtest 2: Simple URL without tag (default to "latest")
- TC-REGISTRY-IMPL-009: BuildImageReference returns ValidationError for invalid URLs

**D. ValidateRegistry Tests** (2 tests):
- TC-REGISTRY-IMPL-010: ValidateRegistry succeeds for reachable registry (200 or 401)
- TC-REGISTRY-IMPL-011: ValidateRegistry returns NetworkError for unreachable registry

**Minimum Test Count**: 11 tests (more encouraged for edge cases)

**Coverage Target**: ≥85% (per Wave 2 Test Plan)

**Test Execution**:
```bash
# Run tests during development
go test ./pkg/registry -v -cover

# Check coverage
go test ./pkg/registry -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

### Step 8: Measure Implementation Size

**CRITICAL**: Use ONLY the designated line counter tool (R304)

**Measurement Protocol**:

1. **Ensure code is committed**:
   ```bash
   cd $EFFORT_DIR
   git status  # Must show "nothing to commit"
   # If uncommitted:
   git add -A
   git commit -m "feat: registry client implementation complete"
   git push
   ```

2. **Find project root**:
   ```bash
   PROJECT_ROOT=$(pwd)
   while [ "$PROJECT_ROOT" != "/" ]; do
       [ -f "$PROJECT_ROOT/orchestrator-state.json" ] && break
       PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
   done
   echo "Project root: $PROJECT_ROOT"
   ```

3. **Run line counter** (NO parameters - auto-detects base):
   ```bash
   $PROJECT_ROOT/tools/line-counter.sh
   ```

4. **Interpret output**:
   - Tool will show: "🎯 Detected base: [branch]"
   - Tool will show: "✅ Total implementation lines: [count]"
   - **This count EXCLUDES**: tests, demos, docs, configs (per R007)

**Size Compliance**:
- **Estimate**: 450 lines
- **Acceptable Range**: 383-518 lines (±15% buffer)
- **Hard Limit**: 800 lines (stop if approaching 700)
- **Action if >700 lines**: Stop immediately, notify orchestrator for split planning

**Measurement Frequency**:
- After implementing each major method
- Before committing
- After all implementation complete

---

### Step 9: Commit and Push Code

**Commit Message Format**:
```
feat(registry): implement registry client with push operations

- Implement NewClient with provider validation
- Implement Push with error classification (auth/network/push)
- Implement BuildImageReference with giteaadmin namespace
- Implement ValidateRegistry with /v2/ endpoint check
- Add helper functions for error classification
- Add comprehensive unit tests (85%+ coverage)

Lines: [count] (within estimate of 450)

Related: Effort 1.2.2, Wave 2, Phase 1
```

**Commit Protocol**:
```bash
cd $EFFORT_DIR
git add pkg/registry/client.go pkg/registry/client_test.go
git commit -m "[message above]"
git push origin idpbuilder-oci-push/phase1/wave2/effort-2-registry-client
```

---

### Step 10: Update Work Log

**Work Log File**: Create `work-log.md` in effort root

**Format**:
```markdown
# Registry Client Implementation Work Log

## 2025-10-29 [HH:MM] UTC - Implementation Started
- Read implementation plan
- Read Wave 2 Architecture
- Read Wave 2 Test Plan
- Set up package structure

## 2025-10-29 [HH:MM] UTC - Constructor Complete
- Implemented NewClient with validation
- Tests passing: TC-REGISTRY-IMPL-001 through TC-REGISTRY-IMPL-003
- Coverage: [XX]%

## 2025-10-29 [HH:MM] UTC - Push Method Complete
- Implemented Push with error classification
- Tests passing: TC-REGISTRY-IMPL-004 through TC-REGISTRY-IMPL-007
- Coverage: [XX]%

## 2025-10-29 [HH:MM] UTC - BuildImageReference Complete
- Implemented reference construction
- Tests passing: TC-REGISTRY-IMPL-008, TC-REGISTRY-IMPL-009
- Coverage: [XX]%

## 2025-10-29 [HH:MM] UTC - ValidateRegistry Complete
- Implemented /v2/ endpoint validation
- Tests passing: TC-REGISTRY-IMPL-010, TC-REGISTRY-IMPL-011
- Coverage: [XX]%

## 2025-10-29 [HH:MM] UTC - Implementation Complete
- All methods implemented
- All tests passing (11+ tests)
- Coverage: [XX]% (target: ≥85%)
- Line count: [count] lines (estimate: 450)
- Ready for code review
```

---

## Detailed Code Specifications

### Complete Implementation Reference (from Wave 2 Architecture)

**File: pkg/registry/client.go** (full implementation):

```go
// Package registry provides OCI registry push operations.
package registry

import (
    "context"
    "fmt"
    "net/http"
    "net/url"
    "strings"

    "github.com/google/go-containerregistry/pkg/name"
    v1 "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/remote"
)

// registryClient implements the Client interface using go-containerregistry.
type registryClient struct {
    authProvider auth.Provider
    tlsConfig    tls.ConfigProvider
    httpClient   *http.Client
}

// NewClient creates a new registry client with authentication and TLS configuration.
//
// [Full GoDoc from Step 2]
func NewClient(authProvider auth.Provider, tlsConfig tls.ConfigProvider) (Client, error) {
    if authProvider == nil {
        return nil, &ValidationError{
            Field:   "authProvider",
            Message: "authentication provider cannot be nil",
        }
    }
    if tlsConfig == nil {
        return nil, &ValidationError{
            Field:   "tlsConfig",
            Message: "TLS config provider cannot be nil",
        }
    }

    // Create HTTP client with TLS config
    httpClient := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: tlsConfig.GetTLSConfig(),
        },
    }

    return &registryClient{
        authProvider: authProvider,
        tlsConfig:    tlsConfig,
        httpClient:   httpClient,
    }, nil
}

// Push pushes an OCI image to the specified registry with optional progress reporting.
//
// [Full GoDoc from Step 3]
func (c *registryClient) Push(ctx context.Context, image v1.Image, targetRef string, progressCallback ProgressCallback) error {
    // Parse target reference
    ref, err := name.ParseReference(targetRef)
    if err != nil {
        return &PushFailedError{
            TargetRef: targetRef,
            Cause:     fmt.Errorf("invalid reference: %w", err),
        }
    }

    // Get authenticator
    authenticator, err := c.authProvider.GetAuthenticator()
    if err != nil {
        return &AuthenticationError{
            Registry: targetRef,
            Cause:    err,
        }
    }

    // Configure remote options
    options := []remote.Option{
        remote.WithAuth(authenticator),
        remote.WithTransport(c.httpClient.Transport),
        remote.WithContext(ctx),
    }

    // Add progress callback if provided
    if progressCallback != nil {
        options = append(options, remote.WithProgress(createProgressHandler(progressCallback)))
    }

    // Push image
    err = remote.Write(ref, image, options...)
    if err != nil {
        // Classify error type
        if isAuthError(err) {
            return &AuthenticationError{
                Registry: targetRef,
                Cause:    err,
            }
        }
        if isNetworkError(err) {
            return &NetworkError{
                Registry: targetRef,
                Cause:    err,
            }
        }
        return &PushFailedError{
            TargetRef: targetRef,
            Cause:     err,
        }
    }

    return nil
}

// BuildImageReference constructs a fully qualified registry image reference.
//
// [Full GoDoc from Step 4]
func (c *registryClient) BuildImageReference(registryURL, imageName string) (string, error) {
    // Parse registry URL
    parsedURL, err := url.Parse(registryURL)
    if err != nil {
        return "", &ValidationError{
            Field:   "registryURL",
            Message: fmt.Sprintf("invalid registry URL: %v", err),
        }
    }

    // Extract host:port
    registryHost := parsedURL.Host
    if registryHost == "" {
        return "", &ValidationError{
            Field:   "registryURL",
            Message: "registry URL must include host",
        }
    }

    // Parse image name (extract repository and tag)
    repository, tag := parseImageName(imageName)
    if repository == "" {
        return "", &ValidationError{
            Field:   "imageName",
            Message: "image name cannot be empty",
        }
    }

    // Default tag if not specified
    if tag == "" {
        tag = "latest"
    }

    // Build full reference: registry/namespace/repository:tag
    // Use "giteaadmin" as default namespace for Gitea
    namespace := "giteaadmin"
    fullRef := fmt.Sprintf("%s/%s/%s:%s", registryHost, namespace, repository, tag)

    return fullRef, nil
}

// ValidateRegistry checks if the registry is reachable by pinging the /v2/ endpoint.
//
// [Full GoDoc from Step 5]
func (c *registryClient) ValidateRegistry(ctx context.Context, registryURL string) error {
    // Parse registry URL
    parsedURL, err := url.Parse(registryURL)
    if err != nil {
        return &ValidationError{
            Field:   "registryURL",
            Message: fmt.Sprintf("invalid registry URL: %v", err),
        }
    }

    // Build /v2/ endpoint URL
    v2URL := fmt.Sprintf("%s://%s/v2/", parsedURL.Scheme, parsedURL.Host)

    // Create request
    req, err := http.NewRequestWithContext(ctx, "GET", v2URL, nil)
    if err != nil {
        return &NetworkError{
            Registry: registryURL,
            Cause:    err,
        }
    }

    // Perform request
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return &NetworkError{
            Registry: registryURL,
            Cause:    err,
        }
    }
    defer resp.Body.Close()

    // Check response status
    // 200 OK = registry accessible and doesn't require auth
    // 401 Unauthorized = registry accessible and requires auth (success!)
    if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized {
        return nil
    }

    return &RegistryUnavailableError{
        Registry:   registryURL,
        StatusCode: resp.StatusCode,
    }
}

// Helper functions

func parseImageName(imageName string) (repository, tag string) {
    parts := strings.Split(imageName, ":")
    if len(parts) == 2 {
        return parts[0], parts[1]
    }
    return parts[0], ""
}

func createProgressHandler(callback ProgressCallback) chan v1.Update {
    updates := make(chan v1.Update, 100)
    go func() {
        for update := range updates {
            callback(ProgressUpdate{
                LayerDigest: update.Digest.String(),
                LayerSize:   update.Total,
                BytesPushed: update.Complete,
                Status:      "uploading",
            })
        }
    }()
    return updates
}

func isAuthError(err error) bool {
    errStr := err.Error()
    return strings.Contains(errStr, "401") || strings.Contains(errStr, "403") ||
        strings.Contains(errStr, "unauthorized") || strings.Contains(errStr, "forbidden")
}

func isNetworkError(err error) bool {
    errStr := err.Error()
    return strings.Contains(errStr, "connection") || strings.Contains(errStr, "timeout") ||
        strings.Contains(errStr, "network")
}
```

**Estimated Lines**: ~450 lines (including comprehensive GoDoc)

---

## Size Management

### Size Estimates

**Implementation Code**: 450 lines
- Constructor: ~40 lines
- Push method: ~120 lines
- BuildImageReference: ~70 lines
- ValidateRegistry: ~80 lines
- Helper functions: ~140 lines

**Test Code**: ~400 lines (NOT counted toward size limit per R007)

**Total Package Lines**: 450 (implementation only)

### Size Compliance Strategy

**Size Limit**: 800 lines (hard limit)
**Estimate**: 450 lines
**Buffer**: ±15% = 383-518 lines acceptable range

**Measurement Tool**:
```bash
# Find project root first
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Run line counter (auto-detects base branch)
$PROJECT_ROOT/tools/line-counter.sh
```

**Measurement Frequency**:
- After implementing NewClient: ~40 lines
- After implementing Push: ~160 lines cumulative
- After implementing BuildImageReference: ~230 lines cumulative
- After implementing ValidateRegistry: ~310 lines cumulative
- After implementing helpers: ~450 lines cumulative

**Warning Thresholds**:
- ≥700 lines: Approaching limit, review for optimization
- ≥800 lines: STOP IMMEDIATELY, requires split (notify orchestrator)

**Split Trigger**: If implementation exceeds 800 lines (should not happen with estimate of 450)

---

## Test Requirements

### Test Coverage Target

**Minimum Coverage**: 85% (per Wave 2 Test Plan)

**Rationale**: Complex push logic + error classification requires thorough testing

### Test Categories (from Wave 2 Test Plan)

#### A. Constructor Tests (3 tests)

**TC-REGISTRY-IMPL-001: NewClient Success**
- **Given**: Valid auth provider and TLS config
- **When**: Creating new client
- **Then**: Client created successfully with configured HTTP transport

**TC-REGISTRY-IMPL-002: NewClient Nil Auth Provider**
- **Given**: Nil auth provider
- **When**: Creating new client
- **Then**: Returns ValidationError with field="authProvider"

**TC-REGISTRY-IMPL-003: NewClient Nil TLS Provider**
- **Given**: Nil TLS config
- **When**: Creating new client
- **Then**: Returns ValidationError with field="tlsConfig"

---

#### B. Push Tests (4 tests)

**TC-REGISTRY-IMPL-004: Push Success with Progress**
- **Given**: Valid client, image, and target reference
- **When**: Pushing with progress callback
- **Then**: Push succeeds, callback invoked with progress updates

**TC-REGISTRY-IMPL-005: Push AuthenticationError**
- **Given**: Invalid credentials (401/403 response)
- **When**: Pushing image
- **Then**: Returns AuthenticationError (not NetworkError or PushFailedError)

**TC-REGISTRY-IMPL-006: Push NetworkError**
- **Given**: Unreachable registry (connection refused)
- **When**: Pushing image
- **Then**: Returns NetworkError (not AuthenticationError)

**TC-REGISTRY-IMPL-007: Push PushFailedError**
- **Given**: Invalid target reference format
- **When**: Pushing image
- **Then**: Returns PushFailedError

---

#### C. BuildImageReference Tests (2 tests)

**TC-REGISTRY-IMPL-008: BuildImageReference Constructs Correctly**

**Test Case 1**: Full URL with port and tag
- **Input**: `https://gitea.cnoe.localtest.me:8443`, `myapp:latest`
- **Expected**: `gitea.cnoe.localtest.me:8443/giteaadmin/myapp:latest`

**Test Case 2**: Simple URL, no tag (default "latest")
- **Input**: `https://registry.io`, `myapp`
- **Expected**: `registry.io/giteaadmin/myapp:latest`

**TC-REGISTRY-IMPL-009: BuildImageReference ValidationError**
- **Given**: Malformed registry URL (invalid format)
- **When**: Building reference
- **Then**: Returns ValidationError with field="registryURL"

---

#### D. ValidateRegistry Tests (2 tests)

**TC-REGISTRY-IMPL-010: ValidateRegistry Success**
- **Given**: Reachable registry returning 200 OK or 401 Unauthorized
- **When**: Validating registry
- **Then**: Returns nil (both 200 and 401 considered success)

**TC-REGISTRY-IMPL-011: ValidateRegistry NetworkError**
- **Given**: Unreachable registry (connection refused)
- **When**: Validating registry
- **Then**: Returns NetworkError

---

### Test Execution

**Run Tests**:
```bash
cd $EFFORT_DIR
go test ./pkg/registry -v -cover
```

**Expected Output**:
```
=== RUN   TestNewClient_Success
--- PASS: TestNewClient_Success (0.01s)
=== RUN   TestNewClient_NilAuthProvider
--- PASS: TestNewClient_NilAuthProvider (0.00s)
...
PASS
coverage: 87.3% of statements
ok      github.com/cnoe-io/idpbuilder/pkg/registry    0.234s
```

**Coverage Check**:
```bash
go test ./pkg/registry -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Success Criteria**:
- ✅ All tests passing (100% pass rate)
- ✅ Coverage ≥85%
- ✅ No race conditions (go test -race)
- ✅ No linting errors (golangci-lint run)

---

## Pattern Compliance

### Go Best Practices

**Error Handling**:
- ✅ Use Wave 1 error types exclusively
- ✅ Wrap errors with context using fmt.Errorf("%w", err)
- ✅ Return early on errors (no deep nesting)

**Naming Conventions**:
- ✅ Unexported implementation type: `registryClient`
- ✅ Exported constructor: `NewClient`
- ✅ Interface methods match Wave 1 exactly
- ✅ Helper functions are unexported (lowercase)

**Documentation**:
- ✅ Package-level GoDoc
- ✅ GoDoc for ALL exported functions
- ✅ GoDoc includes examples for complex methods
- ✅ Parameter and return value documentation

**Concurrency**:
- ✅ Progress handler uses goroutine + channel (safe)
- ✅ Context passed through for cancellation
- ✅ No shared mutable state

### Security Considerations

**Input Validation**:
- ✅ Validate registry URL format
- ✅ Validate providers not nil
- ✅ Validate image name not empty

**TLS Configuration**:
- ✅ Accept TLS config via injection (not hardcoded)
- ✅ Support both secure and insecure modes
- ✅ No certificate validation bypass without explicit config

**Credential Handling**:
- ✅ Accept auth provider via injection (no credentials in code)
- ✅ No credential logging or exposure
- ✅ Use go-containerregistry's secure auth handling

**Error Messages**:
- ✅ Don't expose sensitive data in error messages
- ✅ Classify errors appropriately (auth vs network)
- ✅ Provide enough context for debugging

---

## Acceptance Criteria

### Must Have (Blocking)

- [ ] All 3 interface methods implemented correctly
  - [ ] NewClient with provider validation
  - [ ] Push with error classification
  - [ ] BuildImageReference with "giteaadmin" namespace
  - [ ] ValidateRegistry with /v2/ endpoint check

- [ ] All Wave 1 error types used correctly
  - [ ] ValidationError for input validation
  - [ ] AuthenticationError for 401/403 responses
  - [ ] NetworkError for connectivity issues
  - [ ] PushFailedError for other push failures
  - [ ] RegistryUnavailableError for invalid /v2/ responses

- [ ] All tests passing (100% pass rate)
  - [ ] 11+ test cases
  - [ ] All test categories covered
  - [ ] Success paths tested
  - [ ] Error paths tested

- [ ] Code coverage ≥85% (per Wave 2 Test Plan)

- [ ] Line count within estimate
  - [ ] Implementation: 450 lines target
  - [ ] Acceptable: 383-518 lines (±15%)
  - [ ] Hard limit: <800 lines

- [ ] No linting errors
  - [ ] go vet passes
  - [ ] golangci-lint passes

- [ ] Documentation complete
  - [ ] Package-level GoDoc
  - [ ] GoDoc for ALL public methods
  - [ ] Examples in GoDoc where appropriate

### Should Have (Important)

- [ ] Integration with go-containerregistry working
  - [ ] remote.Write() correctly configured
  - [ ] Progress callbacks functional
  - [ ] Reference parsing correct

- [ ] Error classification accurate
  - [ ] Auth errors → AuthenticationError
  - [ ] Network errors → NetworkError
  - [ ] Push failures → PushFailedError
  - [ ] Validation errors → ValidationError

- [ ] /v2/ endpoint validation robust
  - [ ] Accepts 200 OK as success
  - [ ] Accepts 401 Unauthorized as success (registry requires auth)
  - [ ] Rejects other status codes

- [ ] Image reference construction correct
  - [ ] Injects "giteaadmin" namespace
  - [ ] Defaults tag to "latest" if not specified
  - [ ] Handles registry URLs with ports correctly

### Nice to Have (Enhancement)

- [ ] Additional test cases for edge conditions
- [ ] Performance benchmarks for Push operation
- [ ] Integration tests with real registry (optional)
- [ ] Example code in separate examples/ directory

---

## Risk Mitigation

### Risk 1: External Registry Dependency

**Risk**: Push tests require live registry to validate

**Mitigation**:
- Use mocks for unit tests (Wave 1 interfaces enable this)
- Optional integration tests with test registry
- CI pipeline can spin up registry:2 container for integration tests
- Unit tests focus on error classification logic (not actual push)

---

### Risk 2: Error Classification Accuracy

**Risk**: Incorrect error type returns confuse users

**Mitigation**:
- Comprehensive error path testing
- Helper functions (`isAuthError`, `isNetworkError`) with clear patterns
- Test all error scenarios (401, 403, connection refused, timeout)
- Integration tests validate end-to-end error handling

---

### Risk 3: Progress Callback Goroutine Leak

**Risk**: Goroutine in createProgressHandler may leak if channel not closed

**Mitigation**:
- go-containerregistry closes the channel after push completes
- Test with go test -race to detect race conditions
- Document that channel closure is remote.Write() responsibility

---

### Risk 4: Reference Construction Edge Cases

**Risk**: Malformed URLs or image names cause runtime errors

**Mitigation**:
- Validate ALL inputs (registry URL, image name)
- Return ValidationError for malformed inputs
- Test with various URL formats (http, https, with/without port)
- Test with various image name formats (with/without tag)

---

### Risk 5: TLS Configuration Issues

**Risk**: Insecure mode used accidentally in production

**Mitigation**:
- TLS config provided via injection (caller's responsibility)
- IsInsecure() flag available for runtime checks
- CLI will require explicit --insecure flag (Wave 3)
- Clear documentation warnings about insecure mode

---

## Integration Points

### With Wave 1 Interfaces (Read-Only)

**This implementation USES**:
- `pkg/registry.Client` interface (implements)
- `pkg/auth.Provider` interface (dependency injection)
- `pkg/tls.ConfigProvider` interface (dependency injection)
- Wave 1 error types (returns)

**Critical**: NO Wave 1 files modified (frozen contracts)

---

### With Wave 2 Implementations (Runtime Injection)

**This implementation DEPENDS ON** (at runtime):
- Effort 1.2.3: Auth implementation (basic auth provider)
- Effort 1.2.4: TLS implementation (TLS config provider)

**Dependency Strategy**:
- At IMPLEMENTATION time: Use interfaces only (compile-time type checking)
- At RUNTIME: Receive concrete implementations via NewClient constructor
- No circular dependencies (registry doesn't implement auth/TLS)

---

### With go-containerregistry Library

**External Library Integration**:
- `pkg/name.ParseReference()` - Reference parsing
- `pkg/v1/remote.Write()` - Image push operation
- `pkg/authn.Authenticator` - Authentication type
- `pkg/v1.Image` - OCI image type
- `pkg/v1.Update` - Progress update type

**Integration Points**:
- Push method calls remote.Write() with correct options
- Progress handler converts v1.Update → ProgressUpdate
- Error classification interprets remote.Write() errors

---

## Quick Reference Commands

### Development Workflow

```bash
# Navigate to effort directory
EFFORT_DIR="/home/vscode/workspaces/idpbuilder-oci-push-planning/efforts/phase1/wave2/effort-2-registry-client"
cd "$EFFORT_DIR"

# Verify branch
git branch --show-current
# Should be: idpbuilder-oci-push/phase1/wave2/effort-2-registry-client

# Create implementation file
mkdir -p pkg/registry
touch pkg/registry/client.go

# Create test file
touch pkg/registry/client_test.go

# Run tests frequently
go test ./pkg/registry -v -cover

# Check coverage
go test ./pkg/registry -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run linting
go vet ./pkg/registry
golangci-lint run ./pkg/registry

# Measure size (find project root first)
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
$PROJECT_ROOT/tools/line-counter.sh

# Commit when ready
git add pkg/registry/client.go pkg/registry/client_test.go
git commit -m "feat(registry): implement registry client"
git push
```

---

## References

### Primary Sources

1. **Wave 2 Implementation Plan**: `wave-plans/WAVE-2-IMPLEMENTATION.md`
   - Effort scope and requirements
   - Size estimates and parallelization info
   - R213 metadata

2. **Wave 2 Architecture**: `wave-plans/WAVE-2-ARCHITECTURE.md`
   - Complete implementation pseudocode
   - Method specifications
   - Error handling patterns

3. **Wave 2 Test Plan**: `wave-plans/WAVE-2-TEST-PLAN.md`
   - Test case specifications
   - Coverage targets (85%)
   - Progressive Realism approach

4. **Wave 1 Interfaces**: `efforts/phase1/wave1/effort-2-registry-interface/`
   - Interface definition: `pkg/registry/interface.go`
   - Error types: `pkg/registry/errors.go`
   - FROZEN - do not modify

5. **Wave 1 Auth/TLS Interfaces**: `efforts/phase1/wave1/effort-3-auth-tls-interfaces/`
   - Auth interface: `pkg/auth/interface.go`
   - TLS interface: `pkg/tls/interface.go`
   - Used via dependency injection

### Rule References

- **R007**: Test code excluded from line counts
- **R213**: Effort metadata requirements
- **R304**: Mandatory line counter tool usage
- **R307**: Independent branch mergeability
- **R341**: TDD (test plan before implementation)
- **R383**: Metadata file organization with timestamps

---

## Document Status

**Status**: ✅ READY FOR IMPLEMENTATION

**Created**: 2025-10-29 21:33:44 UTC

**Planner**: Code Reviewer Agent (code-reviewer)

**Effort**: 1.2.2 - Registry Client Implementation

**Wave**: Wave 2 of Phase 1

**Compliance Summary**:
- ✅ R213: Complete metadata (effort ID, branch, dependencies, parallelization)
- ✅ R211: Parallelization clearly specified (can parallelize with 1.2.1, 1.2.3, 1.2.4)
- ✅ R341: TDD approach (references Wave 2 Test Plan)
- ✅ R307: Independent mergeability (uses frozen Wave 1 interfaces)
- ✅ R383: Metadata in .software-factory with timestamp
- ✅ Size compliance: Estimated 450 lines < 800 hard limit

**Next Steps**:
1. SW Engineer spawned by orchestrator
2. Engineer reads this plan + architecture + test plan
3. Engineer implements in TDD style
4. Line counter run during development
5. Code Reviewer validates after completion

---

**END OF IMPLEMENTATION PLAN**
