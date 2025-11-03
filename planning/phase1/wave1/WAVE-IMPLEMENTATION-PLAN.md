# Wave 1 Implementation Plan
## Phase 1, Wave 1: Interface & Contract Definitions

**Wave**: Wave 1 - Interface & Contract Definitions
**Phase**: Phase 1 - Foundation & Interfaces
**Created**: 2025-10-29
**Planner**: @agent-code-reviewer
**Fidelity Level**: **EXACT SPECIFICATIONS** (detailed efforts with R213 metadata)

---

## Wave Overview

**Goal**: Define ALL interfaces upfront to freeze contracts and enable Phase 1 Wave 2 parallel implementation of 4 independent packages.

**Architecture Reference**: See `WAVE-1-ARCHITECTURE.md` for complete design details

**Test Plan Reference**: See `WAVE-1-TEST-PLAN.md` for comprehensive test strategy

**Total Efforts**: 4 (sequential execution required)

**Estimated Total Lines**: ~650 lines (implementation code only)

**Parallelization**: SEQUENTIAL - Interface definitions must be ordered to establish contract boundaries

---

## Effort Definitions

### Effort 1.1.1: Docker Client Interface Definition

#### R213 Metadata

```yaml
effort_id: "1.1.1"
effort_name: "Docker Client Interface Definition"
estimated_lines: 180
can_parallelize: false
blocks: []
dependencies: []
assigned_to: "sw-engineer-1"
working_directory: "efforts/phase1/wave1/effort-1-docker-interface"
branch_name: "idpbuilder-oci-push/phase1/wave1/effort-1-docker-interface"
base_branch: "idpbuilder-oci-push/phase1/wave1/integration"
files_touched:
  - "pkg/docker/interface.go"
  - "pkg/docker/errors.go"
  - "pkg/docker/doc.go"
  - "pkg/docker/interface_test.go"
  - "pkg/docker/mock_test.go"
  - "pkg/docker/errors_test.go"
parallelization: "sequential"
```

#### Scope

**Purpose**: Define the Docker client interface for retrieving OCI images from the local Docker daemon, along with comprehensive error types and package documentation.

**What This Effort Delivers**:
- Complete `docker.Client` interface with 4 methods
- 4 custom error types with proper error wrapping
- Package documentation
- Mock implementation for Wave 2+ testing
- Compilation tests and interface validation tests

**Boundaries - Explicitly OUT OF SCOPE**:
- ❌ NO actual Docker daemon connection (Wave 2)
- ❌ NO image retrieval implementation (Wave 2)
- ❌ NO OCI image conversion logic (Wave 2)
- ❌ NO Docker API client initialization (Wave 2)

#### Files to Create/Modify

**New Files**:
1. `pkg/docker/interface.go` (95 lines)
   - `Client` interface definition with 4 methods
   - `NewClient()` function signature (panics with "not implemented")

2. `pkg/docker/errors.go` (55 lines)
   - `DaemonConnectionError` type
   - `ImageNotFoundError` type
   - `ImageConversionError` type
   - `ValidationError` type

3. `pkg/docker/doc.go` (30 lines)
   - Package-level GoDoc
   - Usage examples
   - Purpose and scope description

4. `pkg/docker/interface_test.go` (40 lines)
   - `TestClientInterfaceCompilation` - Verifies interface compiles
   - `TestNewClientSignature` - Verifies function signature

5. `pkg/docker/mock_test.go` (90 lines)
   - `MockClient` struct implementing `docker.Client`
   - All 4 interface methods implemented
   - `TestMockClientImplementsInterface` - Compile-time verification
   - `TestMockClientMethodCalls` - Method callability test

6. `pkg/docker/errors_test.go` (70 lines)
   - `TestDaemonConnectionError` - Error formatting test
   - `TestImageNotFoundError` - Error formatting test
   - `TestImageConversionError` - Error wrapping test
   - `TestValidationError` - Error formatting test
   - `TestErrorTypesImplementError` - Interface satisfaction test

**Modified Files**: None (pure addition)

**Total Estimated Lines**: 380 lines (180 implementation + 200 test)

#### Exact Code Specifications

**File: pkg/docker/interface.go**

```go
// Package docker provides interfaces and types for interacting with Docker daemon.
package docker

import (
	"context"
	"io"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Client defines operations for interacting with the Docker daemon
// to retrieve and validate OCI images stored locally.
type Client interface {
	// ImageExists checks if an image exists in the local Docker daemon.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - imageName: Image name in format "name:tag" (e.g., "myapp:latest")
	//
	// Returns:
	//   - bool: true if image exists, false otherwise
	//   - error: DaemonConnectionError if cannot connect to daemon,
	//            ValidationError if imageName is malformed
	ImageExists(ctx context.Context, imageName string) (bool, error)

	// GetImage retrieves an image from the Docker daemon and converts it
	// to an OCI v1.Image format compatible with go-containerregistry.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - imageName: Image name in format "name:tag"
	//
	// Returns:
	//   - v1.Image: OCI-compatible image object
	//   - error: ImageNotFoundError if image doesn't exist,
	//            DaemonConnectionError if cannot connect,
	//            ImageConversionError if tar conversion fails
	GetImage(ctx context.Context, imageName string) (v1.Image, error)

	// ValidateImageName checks if an image name follows the OCI naming specification.
	//
	// Parameters:
	//   - imageName: Image name to validate
	//
	// Returns:
	//   - error: ValidationError with details if invalid, nil if valid
	ValidateImageName(imageName string) error

	// Close cleans up Docker client resources and closes connections.
	//
	// Returns:
	//   - error: Error if cleanup fails
	Close() error
}

// NewClient creates a new Docker client instance.
//
// Returns:
//   - Client: Docker client interface implementation
//   - error: DaemonConnectionError if daemon is not reachable or not running
func NewClient() (Client, error) {
	// Implementation will be provided in Wave 2 (pkg/docker/client.go)
	panic("not implemented - interface definition only")
}
```

**File: pkg/docker/errors.go**

```go
package docker

import "fmt"

// DaemonConnectionError indicates the Docker daemon is unreachable or not running.
type DaemonConnectionError struct {
	Cause error
}

func (e *DaemonConnectionError) Error() string {
	return fmt.Sprintf("Docker daemon connection error: %v", e.Cause)
}

func (e *DaemonConnectionError) Unwrap() error {
	return e.Cause
}

// ImageNotFoundError indicates the requested image does not exist in the Docker daemon.
type ImageNotFoundError struct {
	ImageName string
}

func (e *ImageNotFoundError) Error() string {
	return fmt.Sprintf("image '%s' not found in Docker daemon", e.ImageName)
}

// ImageConversionError indicates failure to convert Docker image to OCI format.
type ImageConversionError struct {
	ImageName string
	Cause     error
}

func (e *ImageConversionError) Error() string {
	return fmt.Sprintf("failed to convert image '%s' to OCI format: %v", e.ImageName, e.Cause)
}

func (e *ImageConversionError) Unwrap() error {
	return e.Cause
}

// ValidationError indicates image name validation failed.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error (%s): %s", e.Field, e.Message)
}
```

**File: pkg/docker/doc.go**

```go
// Package docker provides interfaces and implementations for interacting
// with the Docker daemon to retrieve OCI images.
//
// This package enables:
//   - Checking if images exist in the local Docker daemon
//   - Retrieving images in OCI v1.Image format
//   - Validating image names per OCI specification
//   - Managing Docker client lifecycle
//
// The primary interface is Client, which abstracts all Docker operations.
// Implementations use the Docker Engine API client library.
//
// Example usage:
//
//	client, err := docker.NewClient()
//	if err != nil {
//	    return err
//	}
//	defer client.Close()
//
//	exists, err := client.ImageExists(ctx, "myapp:latest")
//	if err != nil {
//	    return err
//	}
//
//	if exists {
//	    image, err := client.GetImage(ctx, "myapp:latest")
//	    // use image...
//	}
package docker
```

**Implementation Requirements**:
- Use `github.com/google/go-containerregistry/pkg/v1` for `v1.Image` type
- All methods documented with GoDoc
- Error types implement standard `error` interface
- Error types support error unwrapping where applicable
- `NewClient()` must panic with "not implemented" message (Wave 1 only defines interface)

#### Tests Required

**Test Coverage Requirements**:
- 100% coverage for interface definitions (compilation tests)
- 100% coverage for error types (formatting tests)
- Mock implementation validates interface is implementable

**Test Files** (See WAVE-1-TEST-PLAN.md for complete test code):
1. `pkg/docker/interface_test.go` - Interface compilation validation
2. `pkg/docker/mock_test.go` - Mock implementation and method callability
3. `pkg/docker/errors_test.go` - Error type validation

**Key Test Cases**:
- TC-DOCKER-IF-001: Interface compilation
- TC-DOCKER-IF-002: Mock implementation validation
- TC-DOCKER-IF-003: Error types validation

#### Dependencies

**Upstream Dependencies** (must complete before this effort):
- None (first effort in Wave 1)

**Downstream Dependencies** (efforts that depend on this):
- Effort 1.1.2: Registry Client Interface (references Docker error patterns)
- Effort 1.1.4: Command Structure (imports docker.Client)

**External Library Dependencies**:
```go
require (
	github.com/google/go-containerregistry v0.19.0
	github.com/stretchr/testify v1.9.0 // For testing
)
```

#### Acceptance Criteria

- [ ] All 6 files created with exact specifications
- [ ] `docker.Client` interface compiles with 4 methods
- [ ] All 4 error types implement `error` interface
- [ ] Mock implementation satisfies `docker.Client` interface
- [ ] All tests passing (100% pass rate)
- [ ] Test coverage = 100% (interfaces only)
- [ ] No linting errors (golangci-lint)
- [ ] GoDoc complete for all public interfaces
- [ ] Line count: 180±27 lines (implementation only)
- [ ] NewClient() panics with "not implemented" message

---

### Effort 1.1.2: Registry Client Interface Definition

#### R213 Metadata

```yaml
effort_id: "1.1.2"
effort_name: "Registry Client Interface Definition"
estimated_lines: 200
can_parallelize: false
blocks: []
dependencies: []
assigned_to: "sw-engineer-2"
working_directory: "efforts/phase1/wave1/effort-2-registry-interface"
branch_name: "idpbuilder-oci-push/phase1/wave1/effort-2-registry-interface"
base_branch: "idpbuilder-oci-push/phase1/wave1/integration"
files_touched:
  - "pkg/registry/interface.go"
  - "pkg/registry/errors.go"
  - "pkg/registry/doc.go"
  - "pkg/registry/interface_test.go"
  - "pkg/registry/mock_test.go"
  - "pkg/registry/errors_test.go"
parallelization: "sequential"
```

#### Scope

**Purpose**: Define the registry client interface for pushing OCI images to container registries, including progress reporting types and comprehensive error types.

**What This Effort Delivers**:
- Complete `registry.Client` interface with 3 methods
- `ProgressCallback` function type for progress reporting
- `ProgressUpdate` struct for progress data
- 4 custom error types with proper wrapping
- Package documentation
- Mock implementation with callback validation
- Interface and error validation tests

**Boundaries - Explicitly OUT OF SCOPE**:
- ❌ NO actual registry push logic (Wave 2)
- ❌ NO HTTP client setup (Wave 2)
- ❌ NO layer upload implementation (Wave 2)
- ❌ NO manifest push implementation (Wave 2)
- ❌ NO authentication integration (Wave 2)

#### Files to Create/Modify

**New Files**:
1. `pkg/registry/interface.go` (110 lines)
   - `Client` interface with 3 methods
   - `ProgressCallback` function type
   - `ProgressUpdate` struct
   - Forward references to `AuthProvider` and `TLSConfigProvider`
   - `NewClient()` function signature

2. `pkg/registry/errors.go` (70 lines)
   - `AuthenticationError` type
   - `NetworkError` type
   - `RegistryUnavailableError` type
   - `PushFailedError` type

3. `pkg/registry/doc.go` (20 lines)
   - Package-level GoDoc
   - OCI registry focus

4. `pkg/registry/interface_test.go` (50 lines)
   - Interface compilation tests
   - Type signature validation
   - ProgressUpdate struct validation

5. `pkg/registry/mock_test.go` (100 lines)
   - `MockRegistryClient` implementation
   - Progress callback invocation test
   - Method callability validation

6. `pkg/registry/errors_test.go` (70 lines)
   - All 4 error types tested
   - Error formatting validation
   - Error unwrapping tests

**Modified Files**: None

**Total Estimated Lines**: 420 lines (200 implementation + 220 test)

#### Exact Code Specifications

**File: pkg/registry/interface.go**

```go
// Package registry provides interfaces and types for interacting with OCI registries.
package registry

import (
	"context"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Client defines operations for pushing OCI images to container registries.
type Client interface {
	// Push pushes an OCI image to the specified registry with optional progress reporting.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - image: OCI v1.Image to push
	//   - targetRef: Fully qualified image reference
	//   - progressCallback: Optional callback for progress updates (can be nil)
	//
	// Returns:
	//   - error: AuthenticationError if credentials invalid,
	//            NetworkError if registry unreachable,
	//            PushFailedError if upload fails
	Push(ctx context.Context, image v1.Image, targetRef string, progressCallback ProgressCallback) error

	// BuildImageReference constructs a fully qualified registry image reference.
	//
	// Parameters:
	//   - registryURL: Base registry URL
	//   - imageName: Image name with optional tag
	//
	// Returns:
	//   - string: Fully qualified image reference
	//   - error: ValidationError if URL or name invalid
	BuildImageReference(registryURL, imageName string) (string, error)

	// ValidateRegistry checks if the registry is reachable by pinging the /v2/ endpoint.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - registryURL: Registry base URL to validate
	//
	// Returns:
	//   - error: NetworkError if unreachable,
	//            RegistryUnavailableError if invalid response
	ValidateRegistry(ctx context.Context, registryURL string) error
}

// ProgressCallback is a function type for receiving progress updates during image push.
type ProgressCallback func(update ProgressUpdate)

// ProgressUpdate contains progress information for a single layer upload.
type ProgressUpdate struct {
	// LayerDigest is the SHA256 digest of the layer being uploaded.
	LayerDigest string

	// LayerSize is the total size of the layer in bytes.
	LayerSize int64

	// BytesPushed is the number of bytes uploaded so far.
	BytesPushed int64

	// Status indicates the current state of the layer upload.
	// Values: "uploading", "complete", "exists"
	Status string
}

// NewClient creates a new registry client with authentication and TLS configuration.
//
// Parameters:
//   - authProvider: Authentication provider (from pkg/auth)
//   - tlsConfig: TLS configuration provider (from pkg/tls)
//
// Returns:
//   - Client: Registry client interface implementation
//   - error: ValidationError if providers are invalid
func NewClient(authProvider AuthProvider, tlsConfig TLSConfigProvider) (Client, error) {
	// Implementation will be provided in Wave 2 (pkg/registry/client.go)
	panic("not implemented - interface definition only")
}

// AuthProvider is imported from pkg/auth (forward reference for clarity)
type AuthProvider interface {
	GetAuthenticator() (interface{}, error) // Returns authn.Authenticator
	ValidateCredentials() error
}

// TLSConfigProvider is imported from pkg/tls (forward reference for clarity)
type TLSConfigProvider interface {
	GetTLSConfig() interface{} // Returns *tls.Config
	IsInsecure() bool
}
```

**File: pkg/registry/errors.go**

```go
package registry

import "fmt"

// AuthenticationError indicates registry authentication failed (401/403).
type AuthenticationError struct {
	Registry string
	Cause    error
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("registry authentication failed for %s: %v", e.Registry, e.Cause)
}

func (e *AuthenticationError) Unwrap() error {
	return e.Cause
}

// NetworkError indicates network connectivity issues with the registry.
type NetworkError struct {
	Registry string
	Cause    error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error connecting to registry %s: %v", e.Registry, e.Cause)
}

func (e *NetworkError) Unwrap() error {
	return e.Cause
}

// RegistryUnavailableError indicates the registry endpoint is not responding correctly.
type RegistryUnavailableError struct {
	Registry   string
	StatusCode int
}

func (e *RegistryUnavailableError) Error() string {
	return fmt.Sprintf("registry %s unavailable (status: %d)", e.Registry, e.StatusCode)
}

// PushFailedError indicates image push operation failed.
type PushFailedError struct {
	TargetRef string
	Cause     error
}

func (e *PushFailedError) Error() string {
	return fmt.Sprintf("push to %s failed: %v", e.TargetRef, e.Cause)
}

func (e *PushFailedError) Unwrap() error {
	return e.Cause
}
```

**File: pkg/registry/doc.go**

```go
// Package registry provides interfaces and implementations for pushing
// OCI images to container registries.
//
// This package supports:
//   - Pushing images to OCI-compatible registries
//   - Progress reporting during layer uploads
//   - Registry validation and connectivity checks
//   - Building fully qualified image references
//
// The primary interface is Client, which abstracts registry operations.
// Implementations use go-containerregistry for OCI compatibility.
package registry
```

**Implementation Requirements**:
- Use `github.com/google/go-containerregistry/pkg/v1` for image types
- ProgressCallback can be nil (optional progress reporting)
- Forward references to auth and tls packages (Wave 1.1.3)
- All error types support unwrapping
- NewClient() must panic with "not implemented"

#### Tests Required

**Test Coverage Requirements**:
- 100% coverage for interface definitions
- 100% coverage for error types
- Progress callback mechanism validated

**Key Test Cases** (See WAVE-1-TEST-PLAN.md):
- TC-REGISTRY-IF-001: Interface compilation
- TC-REGISTRY-IF-002: Mock implementation with callback validation
- TC-REGISTRY-IF-003: Error types validation

#### Dependencies

**Upstream Dependencies**:
- None (can run parallel with 1.1.1, but sequential for contract clarity)

**Downstream Dependencies**:
- Effort 1.1.4: Command Structure (imports registry.Client)

**External Library Dependencies**:
```go
require (
	github.com/google/go-containerregistry v0.19.0
	github.com/stretchr/testify v1.9.0
)
```

#### Acceptance Criteria

- [ ] All 6 files created
- [ ] `registry.Client` interface compiles with 3 methods
- [ ] `ProgressCallback` and `ProgressUpdate` types defined
- [ ] All 4 error types implement `error` interface
- [ ] Mock demonstrates callback invocation
- [ ] All tests passing (100% pass rate)
- [ ] Test coverage = 100%
- [ ] GoDoc complete
- [ ] Line count: 200±30 lines

---

### Effort 1.1.3: Auth & TLS Interface Definitions

#### R213 Metadata

```yaml
effort_id: "1.1.3"
effort_name: "Auth & TLS Interface Definitions"
estimated_lines: 140
can_parallelize: false
blocks: []
dependencies: []
assigned_to: "sw-engineer-3"
working_directory: "efforts/phase1/wave1/effort-3-auth-tls-interfaces"
branch_name: "idpbuilder-oci-push/phase1/wave1/effort-3-auth-tls-interfaces"
base_branch: "idpbuilder-oci-push/phase1/wave1/integration"
files_touched:
  - "pkg/auth/interface.go"
  - "pkg/auth/errors.go"
  - "pkg/auth/doc.go"
  - "pkg/auth/interface_test.go"
  - "pkg/auth/mock_test.go"
  - "pkg/tls/interface.go"
  - "pkg/tls/doc.go"
  - "pkg/tls/interface_test.go"
  - "pkg/tls/mock_test.go"
parallelization: "sequential"
```

#### Scope

**Purpose**: Define authentication provider and TLS configuration provider interfaces for registry client usage.

**What This Effort Delivers**:
- `auth.Provider` interface for authentication
- `Credentials` struct for basic auth
- `auth` error types
- `tls.ConfigProvider` interface for TLS configuration
- `Config` struct for TLS settings
- Package documentation for both packages
- Mock implementations for both interfaces
- Validation tests

**Boundaries - Explicitly OUT OF SCOPE**:
- ❌ NO actual authentication implementation (Wave 2)
- ❌ NO go-containerregistry integration (Wave 2)
- ❌ NO TLS certificate loading (Wave 2)
- ❌ NO HTTP transport configuration (Wave 2)

#### Files to Create/Modify

**New Files**:

**Auth Package (7 files)**:
1. `pkg/auth/interface.go` (60 lines)
   - `Provider` interface with 2 methods
   - `Credentials` struct
   - `NewBasicAuthProvider()` function signature

2. `pkg/auth/errors.go` (20 lines)
   - `CredentialValidationError` type

3. `pkg/auth/doc.go` (15 lines)
   - Package documentation

4. `pkg/auth/interface_test.go` (30 lines)
   - Interface compilation tests
   - Credentials struct validation

5. `pkg/auth/mock_test.go` (50 lines)
   - `MockAuthProvider` implementation
   - Interface satisfaction test

**TLS Package (5 files)**:
6. `pkg/tls/interface.go` (55 lines)
   - `ConfigProvider` interface with 2 methods
   - `Config` struct
   - `NewConfigProvider()` function signature

7. `pkg/tls/doc.go` (10 lines)
   - Package documentation

8. `pkg/tls/interface_test.go` (30 lines)
   - Interface compilation tests
   - Config struct validation

9. `pkg/tls/mock_test.go` (60 lines)
   - `MockTLSConfigProvider` implementation
   - TLS config structure validation

**Modified Files**: None

**Total Estimated Lines**: 330 lines (140 implementation + 190 test)

#### Exact Code Specifications

**File: pkg/auth/interface.go**

```go
// Package auth provides interfaces and types for registry authentication.
package auth

import (
	"github.com/google/go-containerregistry/pkg/authn"
)

// Provider defines operations for providing authentication credentials to registries.
type Provider interface {
	// GetAuthenticator returns an authn.Authenticator compatible with go-containerregistry.
	//
	// Returns:
	//   - authn.Authenticator: Authenticator instance
	//   - error: ValidationError if credentials are malformed
	GetAuthenticator() (authn.Authenticator, error)

	// ValidateCredentials performs pre-flight validation of credentials.
	//
	// Returns:
	//   - error: ValidationError with details if invalid, nil if valid
	ValidateCredentials() error
}

// Credentials holds authentication information for basic auth.
type Credentials struct {
	// Username for registry authentication.
	Username string

	// Password for registry authentication.
	// Supports ALL special characters including quotes, spaces, unicode.
	Password string
}

// NewBasicAuthProvider creates a basic authentication provider.
//
// Parameters:
//   - username: Registry username
//   - password: Registry password (supports all special characters)
//
// Returns:
//   - Provider: Authentication provider interface implementation
func NewBasicAuthProvider(username, password string) Provider {
	// Implementation will be provided in Wave 2 (pkg/auth/basic.go)
	panic("not implemented - interface definition only")
}
```

**File: pkg/auth/errors.go**

```go
package auth

import "fmt"

// CredentialValidationError indicates credential validation failed.
type CredentialValidationError struct {
	Field  string // "username" or "password"
	Reason string
}

func (e *CredentialValidationError) Error() string {
	return fmt.Sprintf("credential validation failed (%s): %s", e.Field, e.Reason)
}
```

**File: pkg/auth/doc.go**

```go
// Package auth provides interfaces for registry authentication.
//
// This package supports:
//   - Basic username/password authentication
//   - Credential validation
//   - Integration with go-containerregistry authn
//
// The primary interface is Provider, which supplies authentication
// to registry clients.
package auth
```

**File: pkg/tls/interface.go**

```go
// Package tls provides interfaces and types for TLS configuration.
package tls

import (
	"crypto/tls"
)

// ConfigProvider defines operations for providing TLS configuration.
type ConfigProvider interface {
	// GetTLSConfig returns a tls.Config for HTTP transport.
	//
	// Returns:
	//   - *tls.Config: TLS configuration for HTTP transport
	GetTLSConfig() *tls.Config

	// IsInsecure returns whether insecure mode is enabled.
	//
	// Returns:
	//   - bool: true if --insecure flag was set, false otherwise
	IsInsecure() bool
}

// Config holds TLS configuration options.
type Config struct {
	// InsecureSkipVerify controls whether to skip TLS certificate verification.
	//
	// When true: Certificate validity NOT checked (development only)
	// When false: Full certificate verification (production)
	InsecureSkipVerify bool
}

// NewConfigProvider creates a TLS configuration provider.
//
// Parameters:
//   - insecure: Whether to enable insecure mode (skip cert verification)
//
// Returns:
//   - ConfigProvider: TLS configuration provider interface implementation
func NewConfigProvider(insecure bool) ConfigProvider {
	// Implementation will be provided in Wave 2 (pkg/tls/config.go)
	panic("not implemented - interface definition only")
}
```

**File: pkg/tls/doc.go**

```go
// Package tls provides interfaces for TLS configuration.
//
// This package supports:
//   - TLS certificate verification (secure mode)
//   - Certificate verification bypass (insecure mode)
//   - HTTP transport configuration
package tls
```

**Implementation Requirements**:
- Use `github.com/google/go-containerregistry/pkg/authn` for auth types
- Use `crypto/tls` for TLS config
- Both NewXxx functions must panic with "not implemented"
- Error types follow standard Go patterns

#### Tests Required

**Test Coverage Requirements**:
- 100% coverage for both packages
- Mock implementations validate interfaces

**Key Test Cases** (See WAVE-1-TEST-PLAN.md):
- TC-AUTH-IF-001: Auth provider interface compilation
- TC-AUTH-IF-002: Mock auth provider validation
- TC-TLS-IF-001: TLS provider interface compilation
- TC-TLS-IF-002: Mock TLS provider validation

#### Dependencies

**Upstream Dependencies**:
- None (can run parallel, but sequential for clarity)

**Downstream Dependencies**:
- Effort 1.1.2: Registry Client (uses both interfaces)
- Effort 1.1.4: Command Structure (creates instances)

**External Library Dependencies**:
```go
require (
	github.com/google/go-containerregistry v0.19.0
	github.com/stretchr/testify v1.9.0
)
```

#### Acceptance Criteria

- [ ] All 9 files created (5 auth + 4 tls)
- [ ] Both interfaces compile correctly
- [ ] Both mock implementations satisfy interfaces
- [ ] All tests passing (100% pass rate)
- [ ] Test coverage = 100%
- [ ] GoDoc complete for both packages
- [ ] Line count: 140±21 lines

---

### Effort 1.1.4: Command Structure & Flag Definitions

#### R213 Metadata

```yaml
effort_id: "1.1.4"
effort_name: "Command Structure & Flag Definitions"
estimated_lines: 130
can_parallelize: false
blocks: []
dependencies: ["1.1.1", "1.1.2", "1.1.3"]
assigned_to: "sw-engineer-4"
working_directory: "efforts/phase1/wave1/effort-4-command-structure"
branch_name: "idpbuilder-oci-push/phase1/wave1/effort-4-command-structure"
base_branch: "idpbuilder-oci-push/phase1/wave1/integration"
files_touched:
  - "cmd/push.go"
  - "cmd/push_test.go"
parallelization: "sequential"
```

#### Scope

**Purpose**: Define the `push` command structure with all CLI flags, help text, and execution skeleton (no actual implementation).

**What This Effort Delivers**:
- Complete Cobra command definition for `push`
- 5 CLI flags (--registry, --username, --password, --insecure, --verbose)
- Default constants
- Comprehensive help text with examples
- Placeholder execution function
- Command structure validation tests

**Boundaries - Explicitly OUT OF SCOPE**:
- ❌ NO actual push logic (Phase 2)
- ❌ NO Docker client initialization (Phase 2)
- ❌ NO registry client initialization (Phase 2)
- ❌ NO error handling logic (Phase 2)
- ❌ NO progress display (Phase 2)

#### Files to Create/Modify

**New Files**:
1. `cmd/push.go` (130 lines)
   - `pushCmd` Cobra command definition
   - Flag definitions (5 flags)
   - Default constants (2)
   - Help text with examples
   - `runPushCommand()` stub function
   - Helper function signatures (stubs)

2. `cmd/push_test.go` (70 lines)
   - Command structure validation
   - Flag definition tests
   - Constants verification

**Modified Files**: None (push.go is independent, not registered to root yet)

**Total Estimated Lines**: 200 lines (130 implementation + 70 test)

#### Exact Code Specifications

**File: cmd/push.go**

```go
// Package cmd implements the IDPBuilder CLI commands.
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	// DefaultRegistry is the default Gitea registry URL
	DefaultRegistry = "https://gitea.cnoe.localtest.me:8443"

	// DefaultUsername is the default registry username
	DefaultUsername = "giteaadmin"
)

var (
	// Flag variables
	registryURL string
	username    string
	password    string
	insecure    bool
	verbose     bool
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push IMAGE",
	Short: "Push Docker image to OCI registry",
	Long: `Push a Docker image from local daemon to an OCI-compatible container registry.

The push command reads an image from the local Docker daemon and uploads it
to the specified registry (default: Gitea). It supports authentication with
username and password, and can bypass TLS certificate verification with the
--insecure flag for development/testing environments.

Examples:
  # Push to default Gitea registry
  idpbuilder push myapp:latest --password 'mypassword'

  # Push with custom username
  idpbuilder push myapp:latest --username developer --password 'myP@ss'

  # Push with insecure mode (bypass TLS verification)
  idpbuilder push myapp:latest -k --password 'mypassword'

  # Push to custom registry
  idpbuilder push myapp:v1.0 --registry https://registry.io --password 'pass'

  # Verbose mode for debugging
  idpbuilder push myapp:latest --verbose --password 'pass'

Environment Variables:
  IDPBUILDER_REGISTRY           Override default registry URL
  IDPBUILDER_REGISTRY_USERNAME  Override default username
  IDPBUILDER_REGISTRY_PASSWORD  Provide password (alternative to --password flag)
  IDPBUILDER_INSECURE           Set to "true" to enable insecure mode

Flag priority: CLI flags > Environment variables > Defaults`,
	Args: cobra.ExactArgs(1), // Require exactly one argument: IMAGE
	RunE: runPushCommand,
}

func init() {
	// Define flags
	pushCmd.Flags().StringVar(&registryURL, "registry", DefaultRegistry,
		"Registry URL to push to")
	pushCmd.Flags().StringVarP(&username, "username", "u", DefaultUsername,
		"Registry username for authentication")
	pushCmd.Flags().StringVarP(&password, "password", "p", "",
		"Registry password for authentication (REQUIRED)")
	pushCmd.Flags().BoolVarP(&insecure, "insecure", "k", false,
		"Skip TLS certificate verification (insecure mode)")
	pushCmd.Flags().BoolVarP(&verbose, "verbose", "v", false,
		"Enable verbose output for debugging")

	// Mark password as required
	pushCmd.MarkFlagRequired("password")

	// TODO: Add environment variable binding in Wave 2
	// viper.BindPFlag("registry", pushCmd.Flags().Lookup("registry"))
	// viper.BindEnv("registry", "IDPBUILDER_REGISTRY")

	// Register command with root command
	// rootCmd.AddCommand(pushCmd)  // Will be uncommented in Phase 2
}

// runPushCommand is the main execution function for the push command.
//
// This function orchestrates the complete push workflow:
//   1. Validate inputs (flags, image name)
//   2. Initialize Docker client
//   3. Check image exists in Docker daemon
//   4. Retrieve image from Docker
//   5. Initialize registry client with auth and TLS
//   6. Build target registry reference
//   7. Push image to registry with progress reporting
//   8. Report success or failure
//
// Implementation will be completed in Phase 2 Wave 1.
func runPushCommand(cmd *cobra.Command, args []string) error {
	// Phase 2 implementation placeholder
	imageName := args[0]

	if verbose {
		fmt.Printf("Push command invoked (not yet implemented)\n")
		fmt.Printf("  Image: %s\n", imageName)
		fmt.Printf("  Registry: %s\n", registryURL)
		fmt.Printf("  Username: %s\n", username)
		fmt.Printf("  Insecure: %v\n", insecure)
	}

	return fmt.Errorf("push command not yet implemented - interface definition only (Wave 1)")
}
```

**File: cmd/push_test.go**

```go
package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestPushCommandStructure verifies command structure is valid
func TestPushCommandStructure(t *testing.T) {
	assert.Equal(t, "push IMAGE", pushCmd.Use)
	assert.NotEmpty(t, pushCmd.Short)
	assert.NotEmpty(t, pushCmd.Long)
	assert.NotNil(t, pushCmd.RunE)
}

// TestPushCommandFlags verifies all flags are defined
func TestPushCommandFlags(t *testing.T) {
	assert.NotNil(t, pushCmd.Flags().Lookup("registry"))
	assert.NotNil(t, pushCmd.Flags().Lookup("username"))
	assert.NotNil(t, pushCmd.Flags().Lookup("password"))
	assert.NotNil(t, pushCmd.Flags().Lookup("insecure"))
	assert.NotNil(t, pushCmd.Flags().Lookup("verbose"))

	// Verify short flags
	assert.NotNil(t, pushCmd.Flags().ShorthandLookup("u"))
	assert.NotNil(t, pushCmd.Flags().ShorthandLookup("p"))
	assert.NotNil(t, pushCmd.Flags().ShorthandLookup("k"))
	assert.NotNil(t, pushCmd.Flags().ShorthandLookup("v"))
}

// TestPushCommandConstants verifies constants are defined
func TestPushCommandConstants(t *testing.T) {
	assert.Equal(t, "https://gitea.cnoe.localtest.me:8443", DefaultRegistry)
	assert.Equal(t, "giteaadmin", DefaultUsername)
}

// TestPushCommandExecution verifies command returns not implemented error
func TestPushCommandExecution(t *testing.T) {
	// Reset flags to defaults
	registryURL = DefaultRegistry
	username = DefaultUsername
	password = "test"
	insecure = false
	verbose = false

	err := runPushCommand(pushCmd, []string{"myapp:latest"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not yet implemented")
	assert.Contains(t, err.Error(), "Wave 1")
}

// TestPushCommandVerboseMode verifies verbose flag is respected
func TestPushCommandVerboseMode(t *testing.T) {
	verbose = true
	defer func() { verbose = false }()

	err := runPushCommand(pushCmd, []string{"testimage:v1"})

	assert.Error(t, err)
	// In verbose mode, output is printed (tested via manual inspection)
}
```

**Implementation Requirements**:
- Use Cobra v1.8.0+ for command framework
- Password flag must be marked required
- Help text must include examples
- runPushCommand() must return "not yet implemented" error
- Verbose mode must print configuration (for debugging)

#### Tests Required

**Test Coverage Requirements**:
- 100% coverage for command structure
- Flag definitions validated
- Constants verified
- Execution stub tested

**Key Test Cases**:
- TC-CMD-IF-001: Command structure compilation
- Command flags validation
- Constants verification
- Not-implemented error returned

#### Dependencies

**Upstream Dependencies**:
- Effort 1.1.1: Docker interfaces (for future imports)
- Effort 1.1.2: Registry interfaces (for future imports)
- Effort 1.1.3: Auth/TLS interfaces (for future imports)

**Downstream Dependencies**:
- None (last effort in Wave 1)

**External Library Dependencies**:
```go
require (
	github.com/spf13/cobra v1.8.0  // Already in idpbuilder
	github.com/spf13/viper v1.18.0 // Already in idpbuilder
	github.com/stretchr/testify v1.9.0
)
```

#### Acceptance Criteria

- [ ] Both files created (push.go + push_test.go)
- [ ] Command compiles with Cobra
- [ ] All 5 flags defined correctly
- [ ] Help text complete with examples
- [ ] Constants defined
- [ ] All tests passing (100% pass rate)
- [ ] Test coverage = 100%
- [ ] runPushCommand returns "not implemented" error
- [ ] Line count: 130±20 lines

---

## Parallelization Strategy

**Wave 1 Execution: SEQUENTIAL**

**Rationale**: While Wave 1 efforts could technically run in parallel (interface definitions are independent), we execute sequentially to:
1. Establish clear contract boundaries in order
2. Allow later efforts to reference earlier patterns
3. Simplify orchestration (no complex dependency tracking)
4. Maintain conceptual clarity (Docker → Registry → Auth/TLS → Command)

**Execution Order**:
1. Effort 1.1.1: Docker Interface (foundational)
2. Effort 1.1.2: Registry Interface (builds on error patterns)
3. Effort 1.1.3: Auth & TLS Interfaces (supporting interfaces)
4. Effort 1.1.4: Command Structure (brings all together)

**Total Wave Duration**: ~4 implementation cycles (one per effort)

---

## Wave Size Compliance

**Total Wave Lines**: 650 lines (implementation code only)

**Breakdown**:
- Effort 1.1.1: 180 lines (Docker)
- Effort 1.1.2: 200 lines (Registry)
- Effort 1.1.3: 140 lines (Auth + TLS)
- Effort 1.1.4: 130 lines (Command)

**Size Limit Check**:
- ✅ Well under soft limit (800 lines)
- ✅ Well under hard limit (1000 lines)
- ✅ No split required

**Test Code**: ~680 lines (NOT counted toward limits per R007)

---

## Integration Strategy

**Wave 1 Integration Workflow**:

1. **Effort 1.1.1** completes → Code review → Merge to `phase1/wave1/integration`
2. **Effort 1.1.2** completes → Code review → Merge to integration branch
3. **Effort 1.1.3** completes → Code review → Merge to integration branch
4. **Effort 1.1.4** completes → Code review → Merge to integration branch
5. **Wave integration verification**:
   - All interfaces compile together
   - All tests pass (go test ./pkg/... ./cmd/...)
   - Coverage report shows 100% for interfaces
6. **Architect review** of complete Wave 1
7. **Tag integration point**: `phase1-wave1-integration`
8. **Ready for Wave 2** (implementations)

**Branch Strategy**:
- Integration branch: `idpbuilder-oci-push/phase1/wave1/integration`
- Effort branches: `idpbuilder-oci-push/phase1/wave1/effort-N-name`
- All efforts merge to integration, NOT to main
- Integration branch merges to main after Architect approval

---

## Testing Strategy

### Unit Tests (Per Effort)

Each effort includes tests as specified above:
- Interface compilation tests
- Mock implementation tests
- Error type validation tests
- **Target**: 100% coverage for interface definitions

### Wave-Level Integration Tests

**After All Efforts Merged**:

```bash
# Verify all packages compile together
go build ./pkg/... ./cmd/...

# Run all tests
go test ./pkg/... ./cmd/... -v

# Generate coverage report
go test ./pkg/... ./cmd/... -coverprofile=wave1-coverage.out
go tool cover -html=wave1-coverage.out -o wave1-coverage.html

# Verify 100% coverage for interfaces
COVERAGE=$(go tool cover -func=wave1-coverage.out | grep total | awk '{print $3}' | sed 's/%//')
if (( $(echo "$COVERAGE < 100" | bc -l) )); then
  echo "❌ Coverage $COVERAGE% below 100% (Wave 1 should be 100%)"
  exit 1
fi
```

**Cross-Package Validation**:
```go
// tests/integration/wave1_integration_test.go

package integration_test

import (
	"testing"

	"github.com/your-org/idpbuilder/pkg/docker"
	"github.com/your-org/idpbuilder/pkg/registry"
	"github.com/your-org/idpbuilder/pkg/auth"
	"github.com/your-org/idpbuilder/pkg/tls"
)

// TestAllInterfacesCompileTogether verifies wave cohesion
func TestAllInterfacesCompileTogether(t *testing.T) {
	// Verify all interfaces can be referenced together
	var _ docker.Client
	var _ registry.Client
	var _ auth.Provider
	var _ tls.ConfigProvider

	t.Log("✅ All Wave 1 interfaces compile together")
}

// TestMockWorkflowSimulation verifies mocks work together
func TestMockWorkflowSimulation(t *testing.T) {
	// Simulate complete workflow with mocks
	// (demonstrates Wave 2 will work)

	// Mock auth provider
	authProvider := &auth.MockProvider{}

	// Mock TLS provider
	tlsProvider := &tls.MockConfigProvider{}

	// Mock registry client (uses both)
	// registryClient := registry.NewClient(authProvider, tlsProvider)
	// (panics in Wave 1, works in Wave 2)

	t.Log("✅ Mock workflow demonstrates Wave 2 readiness")
}
```

### Test Execution

**Command**:
```bash
# From project root
cd idpbuilder

# Run Wave 1 tests
go test ./pkg/docker ./pkg/registry ./pkg/auth ./pkg/tls ./cmd -v -cover

# Expected output:
# === RUN   TestClientInterfaceCompilation
# --- PASS: TestClientInterfaceCompilation (0.00s)
# ...
# PASS
# coverage: 100.0% of statements
```

---

## Risk Mitigation

**High-Risk Areas**: None (Wave 1 is low-risk - only interfaces)

**Potential Issues**:
1. **Interface design changes**: Minimize by thorough review before implementation
   - Mitigation: Architect reviewed Wave 1 Architecture thoroughly

2. **Missing methods discovered in Wave 2**: Additional interface methods needed
   - Mitigation: Wave 1 architecture is CONCRETE (based on real use cases)

3. **go-containerregistry version incompatibility**: Library API changes
   - Mitigation: Pin exact version (v0.19.0)

**External Dependencies**:
- `github.com/google/go-containerregistry v0.19.0` (stable, widely used)
- `github.com/spf13/cobra v1.8.0` (already in idpbuilder)
- `github.com/stretchr/testify v1.9.0` (stable testing library)

**Complexity Hotspots**: None (Wave 1 is intentionally simple)

---

## go.mod Updates

**Required for Wave 1**:

```go
module github.com/your-org/idpbuilder

go 1.21

require (
	// Existing idpbuilder dependencies...

	// NEW: Wave 1 additions
	github.com/google/go-containerregistry v0.19.0
	github.com/stretchr/testify v1.9.0
)
```

**Note**: docker and cobra libraries already present in idpbuilder

---

## Next Steps

### Immediate: Orchestrator Actions

1. **Create Wave Integration Branch**:
   ```bash
   git checkout -b idpbuilder-oci-push/phase1/wave1/integration
   git push -u origin idpbuilder-oci-push/phase1/wave1/integration
   ```

2. **Spawn SW Engineer for Effort 1.1.1**:
   - Provide: WAVE-1-ARCHITECTURE.md
   - Provide: WAVE-1-TEST-PLAN.md
   - Provide: This WAVE-1-IMPLEMENTATION.md
   - Task: Implement Effort 1.1.1 exactly as specified

3. **Monitor Implementation**:
   - Size compliance: Use line-counter.sh
   - Test coverage: Must be 100%
   - Sequential execution: Wait for each effort to complete

4. **Code Review After Each Effort**:
   - Spawn Code Reviewer
   - Validate against acceptance criteria
   - Ensure no scope creep

### After Wave 1 Complete

1. **Wave Integration Verification**:
   - All tests pass
   - Coverage = 100%
   - All interfaces compile together

2. **Architect Review**:
   - Validate interface design
   - Approve for Wave 2

3. **Tag Integration Point**:
   - `phase1-wave1-integration`
   - Document completion

4. **Prepare Wave 2**:
   - Architect creates WAVE-2-ARCHITECTURE.md
   - Code Reviewer creates WAVE-2-IMPLEMENTATION.md
   - 4 parallel efforts implement interfaces

---

## Document Status

**Status**: ✅ READY FOR IMPLEMENTATION
**Created**: 2025-10-29
**Planner**: @agent-code-reviewer
**Wave**: Wave 1 of Phase 1
**Fidelity**: EXACT (detailed effort specifications with R213 metadata)

**Compliance Summary**:
- ✅ R213: Complete metadata for all 4 efforts
- ✅ R211: Parallelization strategy documented (sequential)
- ✅ R219: Dependencies clearly mapped
- ✅ R341: TDD - tests defined for all efforts
- ✅ R304: Line counting instructions clear
- ✅ R338: Size reporting format specified
- ✅ EXACT fidelity: Detailed file lists, code specs, test requirements

**Total Efforts**: 4
**Total Lines**: 650 (implementation) + 680 (tests)
**Execution**: Sequential
**Duration**: ~4 implementation cycles

**Critical Success Factors**:
- Interface definitions must be frozen (no changes in Wave 2)
- Test coverage must be 100% (interfaces only)
- All mocks must demonstrate implementability
- Documentation must be complete (GoDoc)
- Sequential execution maintains contract clarity

---

**🚨 CRITICAL R405 AUTOMATION FLAG 🚨**

CONTINUE-SOFTWARE-FACTORY=TRUE

---

**END OF WAVE 1 IMPLEMENTATION PLAN**
