# Wave 2 Implementation Plan

**Wave**: Wave 2 - Core Package Implementations
**Phase**: Phase 1 - Foundation & Interfaces
**Created**: 2025-10-29 06:17:00 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**Fidelity Level**: **EXACT SPECIFICATIONS** (detailed efforts with R213 metadata)

---

## Wave Overview

**Goal**: Implement all four core packages in parallel using the frozen Wave 1 interfaces, enabling complete OCI image push functionality.

**Architecture Reference**: See `wave-plans/WAVE-2-ARCHITECTURE.md` for design details

**Test Plan Reference**: See `wave-plans/WAVE-2-TEST-PLAN.md` for TDD test specifications

**Total Efforts**: 4

**Integration Branch**: `idpbuilder-oci-push/phase1/wave2/integration` (created)

---

## Effort Definitions

### Effort 1.2.1: Docker Client Implementation

#### R213 Metadata

```json
{
  "effort_id": "1.2.1",
  "effort_name": "Docker Client Implementation",
  "branch_name": "idpbuilder-oci-push/phase1/wave2/effort-1-docker-client",
  "base_branch": "idpbuilder-oci-push/phase1/wave2/integration",
  "parent_wave": "WAVE_2",
  "parent_phase": "PHASE_1",
  "depends_on": [],
  "estimated_lines": 400,
  "complexity": "medium",
  "can_parallelize": true,
  "parallel_with": ["1.2.2", "1.2.3", "1.2.4"]
}
```

#### Scope

**Purpose**: Implement the Docker client package that connects to the Docker daemon, validates image names, checks for image existence, retrieves images in OCI format, and properly cleans up resources.

**What This Effort Accomplishes**:
- Complete implementation of `docker.Client` interface (frozen in Wave 1)
- Docker daemon connectivity using Docker Engine API
- Image existence checking with proper error classification
- Image retrieval and conversion to OCI v1.Image format
- Image name validation with security checks (command injection prevention)
- Resource cleanup and connection management

**Boundaries - OUT OF SCOPE**:
- Building or creating Docker images (only retrieval)
- Docker compose or multi-container operations
- Docker network or volume management
- Image pushing (that's registry package responsibility)
- Docker registry authentication (that's auth package responsibility)

#### Files to Create/Modify

**New Files**:
- `pkg/docker/client.go` (~400 lines)
  - `dockerClient` struct implementation
  - `NewClient()` constructor with daemon ping
  - `ImageExists()` method with NotFound handling
  - `GetImage()` method with OCI conversion
  - `ValidateImageName()` method with security checks
  - `Close()` cleanup method
  - Helper functions (`containsString`)

**Test Files** (NOT counted in line estimates per R007):
- `pkg/docker/client_test.go` (~300 lines)
  - 12+ test cases covering all methods
  - Success paths, error paths, edge cases
  - Context cancellation handling
  - Security validation tests

**Modified Files**:
- `go.mod` (add Docker Engine API dependencies, +5 lines)
  ```
  github.com/docker/docker v28.2.2+incompatible
  github.com/docker/go-connections v0.4.0
  ```

**Total Estimated Lines**: 400 lines (implementation only, tests excluded per R007)

#### Exact Code Specifications

**File: pkg/docker/client.go**

The implementation MUST:
1. Use Docker Engine API client (`github.com/docker/docker/client`)
2. Implement ALL 4 methods from Wave 1 `docker.Client` interface
3. Use Wave 1 error types: `DaemonConnectionError`, `ImageNotFoundError`, `ImageConversionError`, `ValidationError`
4. Convert Docker images to OCI v1.Image using `github.com/google/go-containerregistry/pkg/v1/daemon`
5. Validate image names to prevent command injection attacks

**Key Implementation Details** (from Wave 2 Architecture):

**NewClient() requirements**:
- Use `client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())`
- Ping daemon to verify connectivity: `cli.Ping(ctx)`
- Return `DaemonConnectionError` if daemon unreachable
- Store Docker Engine client in `dockerClient` struct

**ImageExists() requirements**:
- Call `ValidateImageName()` first
- Use `cli.ImageInspectWithRaw(ctx, imageName)`
- Return `(false, nil)` if `client.IsErrNotFound(err)` - NOT an error!
- Return `DaemonConnectionError` for connection issues
- Return `true` only if inspect succeeds

**GetImage() requirements**:
- Call `ImageExists()` first, return `ImageNotFoundError` if not found
- Parse image reference using `daemon.NewTag(imageName)`
- Convert to OCI image using `daemon.Image(ref)`
- Return `ImageConversionError` if conversion fails
- Return OCI `v1.Image` compatible with go-containerregistry

**ValidateImageName() requirements**:
- Check for empty string
- Check for command injection attempts: `;`, `|`, `&`, `$`, `` ` ``, `(`, `)`, `<`, `>`, `\`
- Return `ValidationError` with field="imageName" if invalid
- Allow valid characters: alphanumeric, dots, slashes, colons, hyphens, underscores

**Close() requirements**:
- Close underlying Docker Engine client connection
- Check for nil client before closing
- Return any cleanup errors

#### Tests Required

**Test File: pkg/docker/client_test.go**

**Minimum Test Coverage**: 85% (per Wave 2 Test Plan)

**Test Categories** (from Wave 2 Test Plan):

**A. Constructor Tests**:
- TC-DOCKER-IMPL-001: NewClient success with running daemon
- TC-DOCKER-IMPL-002: NewClient fails with DaemonConnectionError when daemon stopped

**B. ImageExists Tests**:
- TC-DOCKER-IMPL-003: ImageExists returns true for present images (alpine:latest)
- TC-DOCKER-IMPL-004: ImageExists returns false (NOT error) for missing images
- TC-DOCKER-IMPL-005: ImageExists returns ValidationError for empty/invalid names

**C. GetImage Tests**:
- TC-DOCKER-IMPL-006: GetImage retrieves and converts image to OCI v1.Image
- TC-DOCKER-IMPL-007: GetImage returns ImageNotFoundError for non-existent images
- TC-DOCKER-IMPL-008: GetImage returns ImageConversionError on conversion failure (integration test)

**D. ValidateImageName Tests**:
- TC-DOCKER-IMPL-009: ValidateImageName passes for valid names (myapp:latest, registry.io/repo:tag)
- TC-DOCKER-IMPL-010: ValidateImageName rejects dangerous names (command injection attempts)

**E. Close Tests**:
- TC-DOCKER-IMPL-011: Close succeeds and cleans up resources

**F. Edge Case Tests**:
- TC-DOCKER-IMPL-012: Context cancellation is respected

**Test Coverage Requirements**:
- Minimum 85% code coverage
- All success paths tested
- All failure paths tested
- Edge cases (cancelled context, malformed input) tested
- Security validation (command injection) tested

#### Dependencies

**Upstream Dependencies** (must complete before this effort):
- Wave 1 Effort 1: Docker interface definition (COMPLETED)
- Integration branch: `idpbuilder-oci-push/phase1/wave2/integration` (CREATED)

**Downstream Dependencies** (efforts that depend on this):
- None (all Wave 2 efforts are parallel)
- Wave 3 CLI will use this package

**External Library Dependencies**:
- `github.com/docker/docker` v28.2.2+ (Docker Engine API)
- `github.com/docker/go-connections` v0.4.0 (Docker client helpers)
- `github.com/google/go-containerregistry` v0.19.0 (already in Wave 1)
- `github.com/stretchr/testify` v1.10.0 (already in Wave 1)

**System Dependencies**:
- Docker daemon running locally
  - Unix socket: `/var/run/docker.sock`
  - Windows named pipe: `npipe:////./pipe/docker_engine`
- DOCKER_HOST environment variable (optional override)

#### Acceptance Criteria

- [ ] All files created/modified as specified
- [ ] All 4 interface methods implemented correctly
- [ ] All tests passing (100% pass rate)
- [ ] Code coverage ≥85% (per Wave 2 Test Plan)
- [ ] No linting errors (go vet, golangci-lint)
- [ ] Documentation complete (all public methods have GoDoc)
- [ ] Line count within estimate (400 lines ±15% = 340-460 lines)
- [ ] Integration with go-containerregistry working (v1.Image conversion)
- [ ] Security validation preventing command injection
- [ ] Proper error type usage (Wave 1 errors)

---

### Effort 1.2.2: Registry Client Implementation

#### R213 Metadata

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

#### Scope

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
- Authentication implementation (auth package responsibility)
- TLS configuration implementation (tls package responsibility)
- Image building or modification (docker package responsibility)
- Registry management or administration
- Image manifest manipulation (go-containerregistry handles this)

#### Files to Create/Modify

**New Files**:
- `pkg/registry/client.go` (~450 lines)
  - `registryClient` struct implementation
  - `NewClient()` constructor with provider validation
  - `Push()` method with progress callbacks
  - `BuildImageReference()` reference construction
  - `ValidateRegistry()` /v2/ endpoint check
  - Helper functions (`parseImageName`, `createProgressHandler`, `isAuthError`, `isNetworkError`)

**Test Files** (NOT counted in line estimates per R007):
- `pkg/registry/client_test.go` (~400 lines)
  - 15+ test cases covering all methods
  - Mock auth and TLS providers
  - Success paths, error paths, edge cases
  - Progress callback validation

**Modified Files**:
- None (go.mod already has go-containerregistry from Wave 1)

**Total Estimated Lines**: 450 lines (implementation only, tests excluded per R007)

#### Exact Code Specifications

**File: pkg/registry/client.go**

The implementation MUST:
1. Use go-containerregistry's `remote` package for OCI push operations
2. Implement ALL 3 methods from Wave 1 `registry.Client` interface
3. Use Wave 1 error types: `AuthenticationError`, `NetworkError`, `RegistryUnavailableError`, `PushFailedError`, `ValidationError`
4. Accept `auth.Provider` and `tls.ConfigProvider` in constructor (dependency injection)
5. Classify errors correctly (401/403 → AuthenticationError, connection issues → NetworkError)

**Key Implementation Details** (from Wave 2 Architecture):

**NewClient() requirements**:
- Validate authProvider is not nil (return ValidationError if nil)
- Validate tlsConfig is not nil (return ValidationError if nil)
- Create HTTP client with TLS config from provider: `httpClient.Transport.TLSClientConfig = tlsConfig.GetTLSConfig()`
- Store providers in `registryClient` struct

**Push() requirements**:
- Parse target reference using `name.ParseReference(targetRef)`
- Get authenticator from auth provider: `authProvider.GetAuthenticator()`
- Configure remote options:
  - `remote.WithAuth(authenticator)`
  - `remote.WithTransport(httpClient.Transport)`
  - `remote.WithContext(ctx)`
  - `remote.WithProgress(createProgressHandler(callback))` if callback provided
- Call `remote.Write(ref, image, options...)`
- Classify errors:
  - Errors containing "401", "403", "unauthorized", "forbidden" → `AuthenticationError`
  - Errors containing "connection", "timeout", "network" → `NetworkError`
  - Other errors → `PushFailedError`

**BuildImageReference() requirements**:
- Parse registry URL using `url.Parse(registryURL)`
- Extract host:port from parsed URL
- Parse image name to extract repository and tag (split on `:`)
- Default tag to "latest" if not specified
- Inject "giteaadmin" as namespace (Gitea default)
- Construct reference: `{host:port}/giteaadmin/{repository}:{tag}`
- Return ValidationError for malformed URLs or empty image names

**ValidateRegistry() requirements**:
- Parse registry URL
- Build /v2/ endpoint: `{scheme}://{host}/v2/`
- Create HTTP GET request with context
- Execute request using HTTP client
- Accept 200 OK or 401 Unauthorized as success (registry is accessible)
- Return `NetworkError` for connection failures
- Return `RegistryUnavailableError` for other status codes

#### Tests Required

**Test File: pkg/registry/client_test.go**

**Minimum Test Coverage**: 85% (per Wave 2 Test Plan)

**Test Categories** (from Wave 2 Test Plan):

**A. Constructor Tests**:
- TC-REGISTRY-IMPL-001: NewClient success with valid providers
- TC-REGISTRY-IMPL-002: NewClient fails with nil auth provider
- TC-REGISTRY-IMPL-003: NewClient fails with nil TLS provider

**B. Push Tests**:
- TC-REGISTRY-IMPL-004: Push success with progress callbacks
- TC-REGISTRY-IMPL-005: Push returns AuthenticationError for 401/403
- TC-REGISTRY-IMPL-006: Push returns NetworkError for unreachable registry
- TC-REGISTRY-IMPL-007: Push returns PushFailedError for invalid target reference

**C. BuildImageReference Tests**:
- TC-REGISTRY-IMPL-008: BuildImageReference constructs correct references
  - `https://gitea.cnoe.localtest.me:8443` + `myapp:latest` → `gitea.cnoe.localtest.me:8443/giteaadmin/myapp:latest`
  - `https://registry.io` + `myapp` → `registry.io/giteaadmin/myapp:latest` (default tag)
- TC-REGISTRY-IMPL-009: BuildImageReference returns ValidationError for invalid URLs

**D. ValidateRegistry Tests**:
- TC-REGISTRY-IMPL-010: ValidateRegistry succeeds for reachable registry (200 or 401)
- TC-REGISTRY-IMPL-011: ValidateRegistry returns NetworkError for unreachable registry

**Test Coverage Requirements**:
- Minimum 85% code coverage
- All success paths tested
- All failure paths tested
- Error classification tested (auth vs network vs push failures)
- Progress callback functionality tested

#### Dependencies

**Upstream Dependencies** (must complete before this effort):
- Wave 1 Effort 2: Registry interface definition (COMPLETED)
- Wave 1 Effort 3: Auth and TLS interface definitions (COMPLETED)
- Integration branch: `idpbuilder-oci-push/phase1/wave2/integration` (CREATED)

**Downstream Dependencies** (efforts that depend on this):
- None (all Wave 2 efforts are parallel)
- Wave 3 CLI will use this package

**External Library Dependencies**:
- `github.com/google/go-containerregistry` v0.19.0 (already in Wave 1)
  - `pkg/name` for reference parsing
  - `pkg/v1/remote` for push operations
- `github.com/stretchr/testify` v1.10.0 (already in Wave 1)

**Internal Package Dependencies**:
- `pkg/auth.Provider` interface (Wave 1, implemented in Effort 1.2.3)
- `pkg/tls.ConfigProvider` interface (Wave 1, implemented in Effort 1.2.4)

#### Acceptance Criteria

- [ ] All files created/modified as specified
- [ ] All 3 interface methods implemented correctly
- [ ] All tests passing (100% pass rate)
- [ ] Code coverage ≥85% (per Wave 2 Test Plan)
- [ ] No linting errors (go vet, golangci-lint)
- [ ] Documentation complete (all public methods have GoDoc)
- [ ] Line count within estimate (450 lines ±15% = 383-518 lines)
- [ ] Integration with go-containerregistry working (remote.Write)
- [ ] Error classification correct (auth vs network vs push)
- [ ] Progress callbacks functional
- [ ] /v2/ endpoint validation working
- [ ] Image reference construction correct (with giteaadmin namespace)

---

### Effort 1.2.3: Authentication Implementation

#### R213 Metadata

```json
{
  "effort_id": "1.2.3",
  "effort_name": "Authentication Implementation",
  "branch_name": "idpbuilder-oci-push/phase1/wave2/effort-3-auth",
  "base_branch": "idpbuilder-oci-push/phase1/wave2/integration",
  "parent_wave": "WAVE_2",
  "parent_phase": "PHASE_1",
  "depends_on": [],
  "estimated_lines": 350,
  "complexity": "medium",
  "can_parallelize": true,
  "parallel_with": ["1.2.1", "1.2.2", "1.2.4"]
}
```

#### Scope

**Purpose**: Implement the basic authentication provider that validates credentials, converts them to go-containerregistry's authn.Authenticator format, and supports special characters in passwords (including quotes, spaces, unicode).

**What This Effort Accomplishes**:
- Complete implementation of `auth.Provider` interface (frozen in Wave 1)
- Basic authentication (username/password) support
- Credential validation with security checks
- Conversion to `authn.Authenticator` for go-containerregistry
- Special character support in passwords (unicode, quotes, spaces)
- Control character detection in usernames

**Boundaries - OUT OF SCOPE**:
- Token-based authentication (OAuth, bearer tokens)
- Certificate-based authentication (mTLS)
- Credential storage or management (keychain, secret stores)
- Password hashing or encryption (plaintext transmission via HTTP Basic Auth)
- Multi-factor authentication
- Session management

#### Files to Create/Modify

**New Files**:
- `pkg/auth/basic.go` (~350 lines)
  - `basicAuthProvider` struct implementation
  - `NewBasicAuthProvider()` constructor
  - `GetAuthenticator()` method returning authn.Authenticator
  - `ValidateCredentials()` validation method
  - Helper function (`containsControlChars`)

**Test Files** (NOT counted in line estimates per R007):
- `pkg/auth/basic_test.go` (~250 lines)
  - 10+ test cases covering all methods
  - Special character password tests
  - Control character username tests
  - Success and failure paths

**Modified Files**:
- None (go.mod already has go-containerregistry from Wave 1)

**Total Estimated Lines**: 350 lines (implementation only, tests excluded per R007)

#### Exact Code Specifications

**File: pkg/auth/basic.go**

The implementation MUST:
1. Use go-containerregistry's `authn` package for authenticator type
2. Implement ALL 2 methods from Wave 1 `auth.Provider` interface
3. Use Wave 1 error type: `CredentialValidationError`
4. Store credentials in `Credentials` struct (from Wave 1)
5. Support ALL special characters in passwords (no restrictions except non-empty)

**Key Implementation Details** (from Wave 2 Architecture):

**NewBasicAuthProvider() requirements**:
- Accept username and password as strings
- Store in `Credentials` struct: `Credentials{Username: username, Password: password}`
- Return `basicAuthProvider` implementing `Provider` interface
- No validation in constructor (validation happens in ValidateCredentials)

**GetAuthenticator() requirements**:
- Call `ValidateCredentials()` first, return error if invalid
- Create `authn.Basic{Username: username, Password: password}`
- Return authenticator compatible with go-containerregistry
- Return `CredentialValidationError` if validation fails

**ValidateCredentials() requirements**:
- Check username is not empty (return CredentialValidationError if empty)
- Check username contains no control characters (< 32 or == 127)
- Check password is not empty (return CredentialValidationError if empty)
- Allow ALL printable characters in password (including quotes, spaces, unicode)
- Return nil if valid

**Security Considerations**:
- Username: No control characters (prevents terminal escape sequence attacks)
- Password: Allow everything (HTTP Basic Auth transmits as-is, base64 encoded)
- No password strength requirements (user's responsibility)
- No credential logging or exposure

#### Tests Required

**Test File: pkg/auth/basic_test.go**

**Minimum Test Coverage**: 90% (security-critical, per Wave 2 Test Plan)

**Test Categories** (from Wave 2 Test Plan):

**A. Constructor Tests**:
- TC-AUTH-IMPL-001: NewBasicAuthProvider success

**B. GetAuthenticator Tests**:
- TC-AUTH-IMPL-002: GetAuthenticator success with valid credentials
- TC-AUTH-IMPL-003: GetAuthenticator fails with empty username

**C. ValidateCredentials Tests**:
- TC-AUTH-IMPL-004: ValidateCredentials passes for valid credentials
  - Simple: `"user"` / `"pass"`
  - Special chars: `"user"` / `"P@ss!w0rd#123"`
  - Unicode: `"user"` / `"пароль密码🔒"`
  - Spaces: `"user"` / `"pass with spaces"`
  - Quotes: `"user"` / `"pass\"with'quotes"`
- TC-AUTH-IMPL-005: ValidateCredentials fails for empty username
- TC-AUTH-IMPL-006: ValidateCredentials fails for empty password
- TC-AUTH-IMPL-007: ValidateCredentials fails for control characters in username
  - Newline: `"user\n"`
  - Tab: `"user\t"`
  - Null byte: `"user\x00"`
  - Escape: `"user\x1b"`

**Test Coverage Requirements**:
- Minimum 90% code coverage (security-critical)
- All success paths tested
- All failure paths tested
- Special character support validated
- Security checks validated (control characters)

#### Dependencies

**Upstream Dependencies** (must complete before this effort):
- Wave 1 Effort 3: Auth interface definition (COMPLETED)
- Integration branch: `idpbuilder-oci-push/phase1/wave2/integration` (CREATED)

**Downstream Dependencies** (efforts that depend on this):
- None (all Wave 2 efforts are parallel)
- Effort 1.2.2 will USE this implementation via interface
- Wave 3 CLI will use this package

**External Library Dependencies**:
- `github.com/google/go-containerregistry` v0.19.0 (already in Wave 1)
  - `pkg/authn` for Authenticator types
- `github.com/stretchr/testify` v1.10.0 (already in Wave 1)

#### Acceptance Criteria

- [ ] All files created/modified as specified
- [ ] All 2 interface methods implemented correctly
- [ ] All tests passing (100% pass rate)
- [ ] Code coverage ≥90% (security-critical, per Wave 2 Test Plan)
- [ ] No linting errors (go vet, golangci-lint)
- [ ] Documentation complete (all public methods have GoDoc)
- [ ] Line count within estimate (350 lines ±15% = 298-403 lines)
- [ ] Integration with go-containerregistry working (authn.Basic)
- [ ] Special character support validated (unicode, quotes, spaces)
- [ ] Control character detection working
- [ ] No credential exposure in logs or errors

---

### Effort 1.2.4: TLS Configuration Implementation

#### R213 Metadata

```json
{
  "effort_id": "1.2.4",
  "effort_name": "TLS Configuration Implementation",
  "branch_name": "idpbuilder-oci-push/phase1/wave2/effort-4-tls",
  "base_branch": "idpbuilder-oci-push/phase1/wave2/integration",
  "parent_wave": "WAVE_2",
  "parent_phase": "PHASE_1",
  "depends_on": [],
  "estimated_lines": 350,
  "complexity": "medium",
  "can_parallelize": true,
  "parallel_with": ["1.2.1", "1.2.2", "1.2.3"]
}
```

#### Scope

**Purpose**: Implement the TLS configuration provider that creates TLS configs for HTTP transports, supports secure and insecure modes, and loads system certificate pools.

**What This Effort Accomplishes**:
- Complete implementation of `tls.ConfigProvider` interface (frozen in Wave 1)
- TLS configuration for HTTP transports
- Secure mode with system certificate verification
- Insecure mode with certificate verification disabled (--insecure flag)
- System certificate pool loading
- Mode detection (IsInsecure flag)

**Boundaries - OUT OF SCOPE**:
- Certificate generation or management
- Custom CA certificate loading (system certs only)
- Client certificate authentication (mTLS)
- Certificate pinning
- TLS version or cipher suite customization

#### Files to Create/Modify

**New Files**:
- `pkg/tls/config.go` (~350 lines)
  - `tlsConfigProvider` struct implementation
  - `NewConfigProvider()` constructor
  - `GetTLSConfig()` method returning *tls.Config
  - `IsInsecure()` mode detection method

**Test Files** (NOT counted in line estimates per R007):
- `pkg/tls/config_test.go` (~200 lines)
  - 8+ test cases covering all methods
  - Secure and insecure mode tests
  - HTTP client integration tests
  - Success paths

**Modified Files**:
- None (standard library only: crypto/tls, crypto/x509)

**Total Estimated Lines**: 350 lines (implementation only, tests excluded per R007)

#### Exact Code Specifications

**File: pkg/tls/config.go**

The implementation MUST:
1. Use standard library `crypto/tls` for TLS configuration
2. Implement ALL 2 methods from Wave 1 `tls.ConfigProvider` interface
3. Use Wave 1 `Config` struct: `Config{InsecureSkipVerify: bool}`
4. Load system certificate pool using `x509.SystemCertPool()`
5. Support both secure and insecure modes

**Key Implementation Details** (from Wave 2 Architecture):

**NewConfigProvider() requirements**:
- Accept `insecure` boolean parameter
- Store in `Config` struct: `Config{InsecureSkipVerify: insecure}`
- Return `tlsConfigProvider` implementing `ConfigProvider` interface
- No validation needed (boolean parameter)

**GetTLSConfig() requirements**:
- Create `tls.Config` with `InsecureSkipVerify` from stored config
- If secure mode (InsecureSkipVerify == false):
  - Load system certificate pool: `x509.SystemCertPool()`
  - Fallback to empty pool if system certs unavailable
  - Set `RootCAs` to system cert pool
- If insecure mode (InsecureSkipVerify == true):
  - No cert pool needed (verification disabled)
- Return `*tls.Config` ready for HTTP transport

**IsInsecure() requirements**:
- Return `Config.InsecureSkipVerify` boolean
- Simple getter method

**Usage Pattern**:
```go
provider := tls.NewConfigProvider(insecure)
tlsConfig := provider.GetTLSConfig()
transport := &http.Transport{
    TLSClientConfig: tlsConfig,
}
client := &http.Client{Transport: transport}
```

#### Tests Required

**Test File: pkg/tls/config_test.go**

**Minimum Test Coverage**: 90% (security-critical, per Wave 2 Test Plan)

**Test Categories** (from Wave 2 Test Plan):

**A. Constructor Tests**:
- TC-TLS-IMPL-001: NewConfigProvider secure mode (insecure=false)
- TC-TLS-IMPL-002: NewConfigProvider insecure mode (insecure=true)

**B. GetTLSConfig Tests**:
- TC-TLS-IMPL-003: GetTLSConfig secure mode
  - Verify InsecureSkipVerify == false
  - Verify RootCAs loaded (system cert pool)
- TC-TLS-IMPL-004: GetTLSConfig insecure mode
  - Verify InsecureSkipVerify == true

**C. IsInsecure Tests**:
- TC-TLS-IMPL-005: IsInsecure returns false for secure mode
- TC-TLS-IMPL-006: IsInsecure returns true for insecure mode

**D. Integration Tests**:
- TC-TLS-IMPL-007: TLSConfig usable with http.Client
  - Create HTTP transport with TLS config
  - Verify transport has correct TLS settings

**Test Coverage Requirements**:
- Minimum 90% code coverage (security-critical)
- All success paths tested
- Secure and insecure modes tested
- System cert pool loading tested
- HTTP client integration tested

#### Dependencies

**Upstream Dependencies** (must complete before this effort):
- Wave 1 Effort 3: TLS interface definition (COMPLETED)
- Integration branch: `idpbuilder-oci-push/phase1/wave2/integration` (CREATED)

**Downstream Dependencies** (efforts that depend on this):
- None (all Wave 2 efforts are parallel)
- Effort 1.2.2 will USE this implementation via interface
- Wave 3 CLI will use this package

**External Library Dependencies**:
- Standard library only:
  - `crypto/tls` for TLS configuration
  - `crypto/x509` for certificate pool
- `github.com/stretchr/testify` v1.10.0 (already in Wave 1, for tests)

#### Acceptance Criteria

- [ ] All files created/modified as specified
- [ ] All 2 interface methods implemented correctly
- [ ] All tests passing (100% pass rate)
- [ ] Code coverage ≥90% (security-critical, per Wave 2 Test Plan)
- [ ] No linting errors (go vet, golangci-lint)
- [ ] Documentation complete (all public methods have GoDoc)
- [ ] Line count within estimate (350 lines ±15% = 298-403 lines)
- [ ] Secure mode loads system certificate pool
- [ ] Insecure mode disables verification correctly
- [ ] TLS config compatible with HTTP transport
- [ ] No security warnings for secure mode

---

## Parallelization Strategy

### Execution Plan: ALL 4 EFFORTS RUN IN PARALLEL

**Why Parallel Execution is Safe**:
1. All Wave 1 interfaces are FROZEN (no changes allowed)
2. Each effort implements a different package (no file conflicts)
3. No cross-effort dependencies during implementation
4. All efforts depend ONLY on Wave 1 interfaces (already complete)
5. Integration happens AFTER all efforts complete

**Parallel Group 1** (ALL 4 EFFORTS):
```
┌─────────────────────────────────────────────────────────────────┐
│ Effort 1.2.1: Docker Client         (400 lines)  ─────┐        │
│ Effort 1.2.2: Registry Client       (450 lines)  ─────┤        │
│ Effort 1.2.3: Auth Implementation   (350 lines)  ─────┼─ ALL   │
│ Effort 1.2.4: TLS Implementation    (350 lines)  ─────┘ PARALLEL│
└─────────────────────────────────────────────────────────────────┘
```

**Execution Timeline**:
- **T+0**: Orchestrator spawns 4 SW Engineers simultaneously
- **T+1**: All 4 engineers start implementation in parallel
- **T+N**: All 4 engineers complete and push code
- **T+N+1**: Orchestrator spawns 4 Code Reviewers for effort reviews
- **T+N+2**: After all reviews pass, orchestrator integrates to wave branch

**Total Wave Duration**: 1 implementation cycle (NOT 4 sequential cycles)

**Rationale**:
- Wave 1 provided stable interface contracts
- Each package has clear boundaries
- No implementation-time coordination needed
- Maximum efficiency through parallelization

---

## Wave Size Compliance

**Total Wave Lines**: 1,550 lines (implementation code only)

**Breakdown**:
- Effort 1.2.1 (Docker): 400 lines
- Effort 1.2.2 (Registry): 450 lines
- Effort 1.2.3 (Auth): 350 lines
- Effort 1.2.4 (TLS): 350 lines

**Size Compliance Check**:
- ✅ Individual effort limit: ALL efforts < 800 lines (hard limit)
  - Largest effort: 450 lines (Registry) - WELL UNDER LIMIT
- ✅ Wave total: 1,550 lines < 3,000 lines (wave soft limit)
- ✅ No split required for any effort
- ✅ Estimates have ±15% buffer built in

**Test Code NOT Counted** (per R007):
- Docker tests: ~300 lines (excluded)
- Registry tests: ~400 lines (excluded)
- Auth tests: ~250 lines (excluded)
- TLS tests: ~200 lines (excluded)
- Total tests: ~1,150 lines (NOT counted toward size limits)

**Status**: ✅ FULLY COMPLIANT - No size issues expected

---

## Integration Strategy

### Wave 2 Integration Process

**Integration Branch**: `idpbuilder-oci-push/phase1/wave2/integration` (CREATED)

**Integration Flow**:
```
1. ALL 4 efforts implement in parallel (from integration branch base)
2. Each effort reviewed independently by Code Reviewer
3. After ALL 4 efforts approved:
   a. Orchestrator merges Effort 1.2.1 → integration branch
   b. Orchestrator merges Effort 1.2.2 → integration branch
   c. Orchestrator merges Effort 1.2.3 → integration branch
   d. Orchestrator merges Effort 1.2.4 → integration branch
4. Run Wave 2 integration tests (all packages together)
5. Architect reviews Wave 2 architecture compliance
6. Merge integration branch → Phase 1 branch
```

**Integration Tests** (after all efforts merged):
- Cross-package integration tests
- Complete push workflow test (Docker → Registry with Auth + TLS)
- Error handling across package boundaries
- Reference Wave 2 Test Plan for integration test specifications

**Success Criteria for Integration**:
- [ ] All 4 efforts merged to integration branch
- [ ] No merge conflicts (parallel implementation ensures this)
- [ ] All unit tests still passing
- [ ] Integration tests passing
- [ ] Coverage targets met (85%/90%)
- [ ] Architect review approved

---

## Testing Strategy

### Test-Driven Development (R341 Compliance)

**Wave 2 follows TDD approach**:
1. ✅ Test Plan created BEFORE implementation (Wave 2 Test Plan completed)
2. ✅ SW Engineers read test plan first
3. ✅ Implementation must pass all tests (no test modifications)
4. ✅ Coverage targets enforced (85%/90%)

### Per-Package Test Requirements

| Package | Test File | Test Count | Coverage Target |
|---------|-----------|------------|-----------------|
| docker | `pkg/docker/client_test.go` | ~12 tests | ≥85% |
| registry | `pkg/registry/client_test.go` | ~15 tests | ≥85% |
| auth | `pkg/auth/basic_test.go` | ~10 tests | ≥90% (security) |
| tls | `pkg/tls/config_test.go` | ~8 tests | ≥90% (security) |

**Total Wave 2 Tests**: ~50 test cases

### Test Execution

**Unit Tests** (per effort):
```bash
# Run during implementation
go test ./pkg/docker -v -cover
go test ./pkg/registry -v -cover
go test ./pkg/auth -v -cover
go test ./pkg/tls -v -cover
```

**Integration Tests** (after all efforts merged):
```bash
# Run on integration branch
go test ./pkg/... -v -cover
go test -run TestCompleteWorkflow ./integration/... -v
```

**Coverage Verification**:
```bash
# Automated coverage check (from Wave 2 Test Plan)
./coverage-check.sh
# Enforces: docker/registry ≥85%, auth/tls ≥90%
```

### Test Reference

**Detailed test specifications**: See `wave-plans/WAVE-2-TEST-PLAN.md`
- 50+ test cases with exact specifications
- Progressive Realism (uses actual Wave 1 interfaces)
- Success paths, error paths, edge cases
- Security validation tests

---

## Risk Mitigation

### High-Risk Areas

**1. Docker Daemon Dependency**:
- **Risk**: Tests fail if Docker daemon not running
- **Mitigation**:
  - Use `testify/require` to skip tests if daemon unavailable
  - CI pipeline ensures Docker daemon running
  - Test helpers detect daemon availability

**2. External Registry Dependency**:
- **Risk**: Registry push tests require live registry
- **Mitigation**:
  - Use mocks for unit tests
  - Optional integration tests with test registry
  - CI pipeline can spin up registry:2 container

**3. Error Classification**:
- **Risk**: Incorrect error type returns confuse users
- **Mitigation**:
  - Comprehensive error path testing
  - Error classification helper functions (`isAuthError`, `isNetworkError`)
  - Integration tests validate end-to-end error handling

**4. Special Characters in Passwords**:
- **Risk**: Password encoding issues with quotes/unicode
- **Mitigation**:
  - Extensive test cases for special characters
  - Unicode password tests
  - go-containerregistry handles base64 encoding

**5. TLS Certificate Verification**:
- **Risk**: Insecure mode used accidentally in production
- **Mitigation**:
  - Clear documentation warnings
  - IsInsecure() flag for runtime checks
  - CLI will require explicit --insecure flag

### Complexity Hotspots

**High Complexity Areas**:
- **Registry Push**: Complex error classification, progress callbacks
  - Mitigation: Extensive testing, clear error patterns
- **Image Conversion**: Docker → OCI format conversion
  - Mitigation: go-containerregistry handles complexity, we just call daemon.Image()
- **Reference Construction**: Building fully-qualified image references
  - Mitigation: Comprehensive test cases with various formats

---

## Compliance Verification

### R307: Independent Branch Mergeability

**Verification**: Each effort can merge independently
- ✅ All efforts use frozen Wave 1 interfaces
- ✅ No cross-effort dependencies
- ✅ Each effort compiles independently
- ✅ All 4 efforts could merge to main separately (though we merge to integration)

**Result**: ✅ COMPLIANT

### R501: Progressive Trunk-Based Development (Cascade Branching)

**Verification**: Wave 2 follows cascade branching
- ✅ Integration branch created from Phase 1 branch
- ✅ All efforts branch from integration branch
- ✅ Wave 2 builds incrementally on Wave 1 foundation
- ✅ No effort branches directly from main

**Cascade Structure**:
```
main
  └─ phase1
       └─ phase1/wave1/integration (Wave 1 complete)
            └─ phase1/wave2/integration (Wave 2 base)
                 ├─ phase1/wave2/effort-1-docker-client
                 ├─ phase1/wave2/effort-2-registry-client
                 ├─ phase1/wave2/effort-3-auth
                 └─ phase1/wave2/effort-4-tls
```

**Result**: ✅ COMPLIANT

### R359: No Code Deletion

**Verification**: Wave 2 is pure addition
- ✅ No Wave 1 interfaces modified or deleted
- ✅ Only NEW implementation files created
- ✅ No existing code removed
- ✅ Purely additive enhancement

**Result**: ✅ COMPLIANT

### R383: Metadata File Organization

**Verification**: All metadata properly organized
- ✅ Wave plans in `wave-plans/` directory
- ✅ Effort plans will be in `.software-factory/phase1/wave2/effort-*/`
- ✅ Working trees clean (only source code)
- ✅ All metadata timestamped

**Result**: ✅ COMPLIANT

### R341: Test-Driven Development

**Verification**: Tests before implementation
- ✅ Wave 2 Test Plan completed before implementation
- ✅ Test specifications reference actual Wave 1 interfaces
- ✅ Progressive Realism approach used
- ✅ Coverage targets defined (85%/90%)
- ✅ SW Engineers will write code to pass tests

**Result**: ✅ COMPLIANT

---

## Next Steps

### Immediate Actions (Orchestrator)

1. **Validate Infrastructure**:
   - ✅ Integration branch created: `idpbuilder-oci-push/phase1/wave2/integration`
   - [ ] Verify go.mod ready for Docker dependencies
   - [ ] Verify Wave 1 interfaces available

2. **Spawn SW Engineers** (PARALLEL):
   ```
   Spawn 4 SW Engineers simultaneously:
   - Engineer 1: Effort 1.2.1 (Docker Client)
   - Engineer 2: Effort 1.2.2 (Registry Client)
   - Engineer 3: Effort 1.2.3 (Auth Implementation)
   - Engineer 4: Effort 1.2.4 (TLS Implementation)

   Provide to EACH engineer:
   - This implementation plan (WAVE-2-IMPLEMENTATION.md)
   - Wave 2 Architecture (WAVE-2-ARCHITECTURE.md)
   - Wave 2 Test Plan (WAVE-2-TEST-PLAN.md)
   - Their specific effort section from this plan
   - Wave 1 interface files (for reference)
   ```

3. **Monitor Progress**:
   - Track line counts during implementation
   - Watch for any effort approaching 800 line limit
   - Coordinate if any issues arise

### SW Engineer Instructions

**For EACH SW Engineer starting their effort**:

1. **Pre-Implementation Reading** (MANDATORY):
   - [ ] Read this implementation plan (your effort section)
   - [ ] Read Wave 2 Architecture (detailed design)
   - [ ] Read Wave 2 Test Plan (your package test cases)
   - [ ] Read Wave 1 interface definition (your package)

2. **Checkout and Setup**:
   ```bash
   # Checkout your effort branch (Orchestrator creates this)
   git checkout -b <your-effort-branch> origin/phase1/wave2/integration

   # Verify you're based on integration
   git log --oneline -5
   ```

3. **Implementation Workflow**:
   - Write tests FIRST (from test plan)
   - Implement code to pass tests
   - Measure size frequently: `tools/line-counter.sh`
   - Commit regularly
   - Push when complete

4. **Acceptance Criteria**:
   - All tests passing
   - Coverage target met (85% or 90%)
   - Line count within estimate (±15%)
   - No linting errors
   - Documentation complete

### Code Reviewer Instructions

**After SW Engineer completes effort**:

1. **Spawn Code Reviewer** for the completed effort
2. **Code Reviewer validates**:
   - [ ] All acceptance criteria met
   - [ ] Tests passing
   - [ ] Coverage ≥85% (or 90% for auth/tls)
   - [ ] Line count within limit
   - [ ] Interface correctly implemented
   - [ ] Error types properly used
   - [ ] Documentation complete

3. **Decision**:
   - ACCEPTED → Merge to integration branch
   - NEEDS_FIXES → Send back to SW Engineer with specific feedback
   - NEEDS_SPLIT → Should not happen (efforts well under 800 lines)

### Integration Phase

**After ALL 4 efforts approved**:

1. **Merge all efforts** to integration branch
2. **Run integration tests** (cross-package tests)
3. **Architect review** of Wave 2 compliance
4. **Final validation**:
   - All packages working together
   - Complete push workflow functional
   - Coverage targets met across wave
5. **Merge integration branch** → Phase 1 branch

---

## Document Status

**Status**: ✅ READY FOR WAVE 2 IMPLEMENTATION
**Created**: 2025-10-29 06:17:00 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**Wave**: Wave 2 of Phase 1
**Fidelity**: EXACT SPECIFICATIONS (R213 compliant)
**Integration Branch**: `idpbuilder-oci-push/phase1/wave2/integration` (CREATED)

**Compliance Summary**:
- ✅ R213: Complete metadata for all 4 efforts
- ✅ R211: Parallelization strategy specified (all 4 parallel)
- ✅ R501: Cascade branching from integration branch
- ✅ R341: TDD approach (test plan before implementation)
- ✅ R307: Independent branch mergeability ensured
- ✅ R359: No code deletion (pure addition)
- ✅ R383: Metadata properly organized
- ✅ Size compliance: All efforts < 800 lines

**Next State Transition**: SPAWN_SW_ENGINEERS (4 parallel)
- Orchestrator will spawn 4 SW Engineers simultaneously
- Each engineer receives this plan + architecture + test plan
- Implementation proceeds in parallel for maximum efficiency
- Code reviewers spawn as each effort completes

---

**🚨 CRITICAL R405 AUTOMATION FLAG 🚨**

CONTINUE-SOFTWARE-FACTORY=TRUE

---

**END OF WAVE 2 IMPLEMENTATION PLAN**
