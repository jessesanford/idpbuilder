# Effort 1.1.2 Implementation Plan: Registry Client Interface Definition

## CRITICAL EFFORT METADATA (FROM WAVE PLAN)

**Branch**: idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.2
**Can Parallelize**: false
**Parallel With**: [] (None - sequential after 1.1.1)
**Size Estimate**: 200 lines
**Dependencies**: Effort 1.1.1 (Docker Client Interface)

## Overview

**Effort**: Registry Client Interface Definition
**Phase**: 1, Wave: 1.1
**Estimated Size**: 200 lines
**Implementation Time**: 2-3 hours

This effort creates the RegistryClient interface that will be implemented in Wave 2 to push images to OCI-compliant registries with progress reporting.

## File Structure

- `pkg/registry/interface.go`: Complete RegistryClient interface definition (200 lines)
  - RegistryClient interface (3 methods)
  - LayerStatus enum with String() method
  - ProgressUpdate struct
  - ProgressCallback function type
  - 3 error types with Error() and Unwrap() methods
  - NewRegistryClient() constructor stub
  - Complete package documentation

## Implementation Steps

### Step 1: Create Package Directory

```bash
mkdir -p pkg/registry
```

### Step 2: Create pkg/registry/interface.go

Copy EXACTLY from WAVE-1.1-ARCHITECTURE.md lines 219-606:

```go
// Package registry provides an interface to OCI-compliant container registries.
package registry

import (
	"context"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// LayerStatus represents the current status of a layer during push.
type LayerStatus int

const (
	// LayerWaiting indicates the layer is queued but not yet uploading.
	LayerWaiting LayerStatus = iota
	// LayerUploading indicates the layer is currently being uploaded.
	LayerUploading
	// LayerComplete indicates the layer upload completed successfully.
	LayerComplete
	// LayerFailed indicates the layer upload failed.
	LayerFailed
)

// String returns a string representation of the LayerStatus.
func (s LayerStatus) String() string {
	switch s {
	case LayerWaiting:
		return "Waiting"
	case LayerUploading:
		return "Uploading"
	case LayerComplete:
		return "Complete"
	case LayerFailed:
		return "Failed"
	default:
		return "Unknown"
	}
}

// ProgressUpdate represents the progress of a single layer upload.
type ProgressUpdate struct {
	// LayerDigest is the SHA256 digest of the layer (e.g., "sha256:abc123...")
	LayerDigest string

	// LayerSize is the total size of the layer in bytes.
	LayerSize int64

	// BytesUploaded is the number of bytes uploaded so far.
	BytesUploaded int64

	// Status is the current status of this layer.
	Status LayerStatus
}

// ProgressCallback is invoked during image push to report upload progress.
// The callback is called once for each layer as its status changes.
// Implementations should not block; long-running operations should be asynchronous.
type ProgressCallback func(update ProgressUpdate)

// RegistryClient handles operations with OCI-compliant container registries.
type RegistryClient interface {
	// Push uploads an image to a registry with progress reporting.
	//
	// The targetRef must be a full registry reference in the format:
	//   "registry.example.com:port/namespace/repository:tag"
	//
	// The progress callback is invoked for each layer as it uploads.
	// Pass nil for progress if callbacks are not needed.
	//
	// Example:
	//   err := client.Push(ctx, image, "gitea.example.com:8443/giteaadmin/myapp:latest", func(update ProgressUpdate) {
	//       fmt.Printf("Layer %s: %d/%d bytes (%s)\n",
	//           update.LayerDigest[:12], update.BytesUploaded, update.LayerSize, update.Status)
	//   })
	Push(ctx context.Context, image v1.Image, targetRef string, progress ProgressCallback) error

	// BuildImageReference constructs a full registry reference from a registry URL and image name.
	//
	// The registryURL should be the base URL of the registry (e.g., "https://gitea.example.com:8443").
	// The imageName should be in the format "repository:tag" or "namespace/repository:tag".
	//
	// Example:
	//   ref, err := client.BuildImageReference("https://gitea.example.com:8443", "myapp:v1")
	//   // Returns: "gitea.example.com:8443/giteaadmin/myapp:v1", nil
	BuildImageReference(registryURL, imageName string) (string, error)

	// ValidateRegistry performs connectivity and OCI compliance checks on a registry.
	// Returns an error if the registry is unreachable or does not support the OCI Distribution spec.
	//
	// Example:
	//   if err := client.ValidateRegistry(ctx, "https://gitea.example.com:8443"); err != nil {
	//       return fmt.Errorf("registry validation failed: %w", err)
	//   }
	ValidateRegistry(ctx context.Context, registryURL string) error
}

// NewRegistryClient creates a new OCI registry client.
// The auth and tls parameters customize the connection behavior.
//
// Example:
//   auth, _ := auth.NewAuthProvider("username", "password")
//   tls, _ := tls.NewTLSProvider(true) // insecure mode for local dev
//   client, err := registry.NewRegistryClient(auth, tls)
func NewRegistryClient(auth AuthProvider, tls TLSProvider) (RegistryClient, error) {
	// Implementation will be provided in Wave 2
	panic("not implemented")
}

// AuthProvider is imported from pkg/auth (defined below)
type AuthProvider interface{}

// TLSProvider is imported from pkg/tls (defined below)
type TLSProvider interface{}

// RegistryAuthError indicates that authentication to the registry failed.
type RegistryAuthError struct {
	Registry string
	Cause    error
}

func (e *RegistryAuthError) Error() string {
	if e.Cause != nil {
		return "authentication failed for registry " + e.Registry + ": " + e.Cause.Error()
	}
	return "authentication failed for registry " + e.Registry
}

func (e *RegistryAuthError) Unwrap() error {
	return e.Cause
}

// RegistryConnectionError indicates that the registry is unreachable.
type RegistryConnectionError struct {
	Registry string
	Cause    error
}

func (e *RegistryConnectionError) Error() string {
	if e.Cause != nil {
		return "cannot connect to registry " + e.Registry + ": " + e.Cause.Error()
	}
	return "cannot connect to registry " + e.Registry
}

func (e *RegistryConnectionError) Unwrap() error {
	return e.Cause
}

// LayerPushError indicates that uploading a specific layer failed.
type LayerPushError struct {
	LayerDigest string
	Cause       error
}

func (e *LayerPushError) Error() string {
	if e.Cause != nil {
		return "failed to push layer " + e.LayerDigest + ": " + e.Cause.Error()
	}
	return "failed to push layer " + e.LayerDigest
}

func (e *LayerPushError) Unwrap() error {
	return e.Cause
}
```

**CRITICAL REQUIREMENTS**:
- Use EXACT code from architecture document (no modifications)
- Import `github.com/google/go-containerregistry/pkg/v1` for v1.Image type
- LayerStatus MUST implement String() method (Stringer interface)
- All error types MUST implement error interface via Error() method
- All error types MUST implement Unwrap() for error chain
- NewRegistryClient MUST panic("not implemented") - no actual implementation yet
- AuthProvider and TLSProvider are placeholder interfaces (will be defined in Effort 1.1.3)
- All public types and methods MUST have Go documentation comments

### Step 3: Build Validation

After creating the file:

```bash
cd pkg/registry
go build .
# Expected: Build succeeds (pure interface definition)
```

### Step 4: Create Tests

Create `pkg/registry/interface_test.go` with the following tests:

```go
package registry_test

import (
	"errors"
	"testing"

	"github.com/jessesanford/idpbuilder/pkg/registry"
)

// T1.1.2-001: RegistryClient interface compiles
func TestRegistryClientInterfaceCompiles(t *testing.T) {
	var _ registry.RegistryClient = nil
}

// T1.1.2-002: LayerStatus enum compiles with String() method
func TestLayerStatus_StringMethod(t *testing.T) {
	tests := []struct {
		status   registry.LayerStatus
		expected string
	}{
		{registry.LayerWaiting, "Waiting"},
		{registry.LayerUploading, "Uploading"},
		{registry.LayerComplete, "Complete"},
		{registry.LayerFailed, "Failed"},
		{registry.LayerStatus(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.status.String(); got != tt.expected {
				t.Errorf("LayerStatus.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// T1.1.2-003: ProgressUpdate struct compiles
func TestProgressUpdate_StructValid(t *testing.T) {
	update := registry.ProgressUpdate{
		LayerDigest:   "sha256:abc123",
		LayerSize:     1000,
		BytesUploaded: 500,
		Status:        registry.LayerUploading,
	}

	if update.LayerDigest != "sha256:abc123" {
		t.Error("ProgressUpdate struct field assignment failed")
	}
}

// T1.1.2-004: ProgressCallback type signature valid
func TestProgressCallback_TypeValid(t *testing.T) {
	var callback registry.ProgressCallback = func(update registry.ProgressUpdate) {
		// Callback implementation
	}

	// Invoke callback to verify signature
	callback(registry.ProgressUpdate{
		LayerDigest:   "sha256:test",
		LayerSize:     100,
		BytesUploaded: 50,
		Status:        registry.LayerUploading,
	})
}

// T1.1.2-005: RegistryAuthError implements error with Unwrap
func TestRegistryAuthError_ImplementsError(t *testing.T) {
	cause := errors.New("401 Unauthorized")
	err := &registry.RegistryAuthError{
		Registry: "gitea.example.com:8443",
		Cause:    cause,
	}

	var _ error = err // Compile-time check

	if !errors.Is(err, cause) {
		t.Error("RegistryAuthError should unwrap to cause")
	}

	expectedMsg := "authentication failed for registry gitea.example.com:8443: 401 Unauthorized"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message %q, got %q", expectedMsg, err.Error())
	}
}

// T1.1.2-006: RegistryConnectionError implements error with Unwrap
func TestRegistryConnectionError_ImplementsError(t *testing.T) {
	cause := errors.New("connection refused")
	err := &registry.RegistryConnectionError{
		Registry: "gitea.example.com:8443",
		Cause:    cause,
	}

	var _ error = err // Compile-time check

	if !errors.Is(err, cause) {
		t.Error("RegistryConnectionError should unwrap to cause")
	}
}

// T1.1.2-007: LayerPushError implements error with Unwrap
func TestLayerPushError_ImplementsError(t *testing.T) {
	cause := errors.New("blob upload failed")
	err := &registry.LayerPushError{
		LayerDigest: "sha256:abc123...",
		Cause:       cause,
	}

	var _ error = err // Compile-time check

	if !errors.Is(err, cause) {
		t.Error("LayerPushError should unwrap to cause")
	}
}

// T1.1.2-008: NewRegistryClient constructor signature valid
func TestNewRegistryClient_SignatureValid(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "not implemented" {
				t.Errorf("Expected panic 'not implemented', got %v", r)
			}
		} else {
			t.Error("Expected NewRegistryClient to panic (not implemented)")
		}
	}()

	_, _ = registry.NewRegistryClient(nil, nil)
}
```

### Step 5: Run Tests

```bash
cd pkg/registry
go test -v -cover
# Expected: 9 tests PASS, 100% coverage
```

## Size Management

- **Estimated Lines**: 200
- **Measurement Tool**: ${PROJECT_ROOT}/tools/line-counter.sh (find project root first)
- **Check Frequency**: After file creation
- **Split Threshold (R535)**: 700 lines (warning), 900 lines (Code Reviewer enforcement)
- **Note**: SW Engineers see 800-line limit, Code Reviewers enforce at 900 (grace buffer)

**Size Verification**:
```bash
# Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ]; then break; fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Measure implementation lines
$PROJECT_ROOT/tools/line-counter.sh
# Expected: ~200 lines
```

## Test Requirements

- **Unit Tests**: 100% coverage (interface definitions are simple)
- **Integration Tests**: N/A (interfaces only, no implementation)
- **E2E Tests**: N/A (Wave 1 is interfaces only)

**Test Coverage Requirements**:
- Minimum 100% coverage (interface definitions and simple types)
- All error types tested
- LayerStatus String() method tested
- ProgressUpdate struct tested
- ProgressCallback type tested
- Interface compilation verified

## Pattern Compliance

**Go Patterns**:
- Standard error interface implementation
- Error wrapping with Unwrap() method
- Context-aware method signatures
- Stringer interface for LayerStatus enum
- go-containerregistry v1.Image integration

**Security Requirements**:
- No security concerns (pure interface definitions)
- Error types preserve error chains for debugging

**Performance Targets**:
- N/A (no implementation in this effort)

## Dependencies

**Upstream Dependencies** (must complete before this effort):
- Effort 1.1.1 (Docker Interface) - sequential ordering

**Downstream Dependencies** (efforts that depend on this):
- Effort 1.1.4 (Command Structure) imports pkg/registry

**External Dependencies**:
- `github.com/google/go-containerregistry/pkg/v1` (already in go.mod)

## Acceptance Criteria

- [ ] File `pkg/registry/interface.go` created with 200 lines
- [ ] RegistryClient interface defined with 3 methods
- [ ] LayerStatus enum defined with 4 values and String() method
- [ ] ProgressUpdate struct defined with 4 fields
- [ ] ProgressCallback function type defined
- [ ] 3 error types defined with Error() and Unwrap() methods
- [ ] NewRegistryClient() constructor stub created
- [ ] All tests passing (9 tests, 100% pass rate)
- [ ] `go build ./pkg/registry` succeeds
- [ ] `go test ./pkg/registry` succeeds
- [ ] No linting errors (`golangci-lint run ./pkg/registry`)
- [ ] All public types have documentation comments
- [ ] Line count within estimate (±15%: 170-230 lines acceptable)

## Integration with Effort 1.1.1

This effort builds on the Docker Client interface from 1.1.1:

**Reference Pattern**:
- Effort 1.1.1 defined DockerClient.GetImage() returning v1.Image
- This effort defines RegistryClient.Push() accepting v1.Image
- Integration point: Docker image → Registry push

**Usage Flow** (to be implemented in Wave 2):
```go
// Get image from Docker (1.1.1 interface)
image, err := dockerClient.GetImage(ctx, "myapp:latest")

// Push image to registry (1.1.2 interface)
err = registryClient.Push(ctx, image, targetRef, progressCallback)
```

## R340 Compliance (EXACT Fidelity)

This implementation plan provides EXACT specifications:

✅ **Concrete Code**: Real Go interface copied from architecture
✅ **Complete Types**: All structs, enums, and errors defined
✅ **Method Signatures**: Exact parameter and return types
✅ **Error Handling**: All error types with Error() and Unwrap()
✅ **Documentation**: All public types have doc comments
✅ **No Pseudocode**: All code is real Go syntax

## Notes for SW Engineer

- **Copy EXACTLY**: Do not modify the interface code from the architecture document
- **No Implementation**: Constructor panics with "not implemented" - actual logic comes in Wave 2
- **Placeholder Interfaces**: AuthProvider and TLSProvider are empty interfaces for now (defined in 1.1.3)
- **Test First**: Create tests before or alongside interface to verify correctness
- **Build Validation**: Must compile cleanly with `go build`

---

**Document Version**: 1.0
**Created**: 2025-11-11 22:43:45 UTC
**Status**: Ready for SW Engineer Implementation
**Expected Next State**: SW Engineer implements → Code Reviewer reviews → Merge to integration branch
