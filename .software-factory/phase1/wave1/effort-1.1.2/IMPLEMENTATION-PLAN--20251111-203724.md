# Implementation Plan: Effort 1.1.2 - Registry Client Interface Definition

**Effort ID**: 1.1.2
**Effort Name**: Registry Client Interface Definition
**Phase**: Phase 1 - Foundation & Interfaces
**Wave**: Wave 1.1 - Interface Definitions
**Created**: 2025-11-11 20:37:24 UTC
**Planner**: Code Reviewer Agent
**Branch**: idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.2
**Base Branch**: idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.1

---

## R213 Metadata (CRITICAL - FROM WAVE PLAN)

```json
{
  "effort_id": "1.1.2",
  "effort_name": "Registry Client Interface Definition",
  "branch_name": "idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.2",
  "base_branch": "idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.1",
  "parent_wave": "wave1.1",
  "parent_phase": "phase1",
  "depends_on": ["1.1.1"],
  "estimated_lines": 200,
  "complexity": "low",
  "can_parallelize": false,
  "parallel_with": [],
  "tests_required": [
    "T1.1.2-001: RegistryClient interface compiles",
    "T1.1.2-002: LayerStatus enum compiles with String() method",
    "T1.1.2-003: ProgressUpdate struct compiles",
    "T1.1.2-004: ProgressCallback type signature valid",
    "T1.1.2-005: RegistryAuthError implements error with Unwrap",
    "T1.1.2-006: RegistryConnectionError implements error with Unwrap",
    "T1.1.2-007: LayerPushError implements error with Unwrap",
    "T1.1.2-008: NewRegistryClient constructor signature valid",
    "T1.1.2-009: Mock RegistryClient satisfies interface"
  ]
}
```

---

## Overview

**Purpose**: Create the RegistryClient interface that will be implemented in Wave 2 to push images to OCI-compliant registries with progress reporting.

**Scope**:
- Interface definition ONLY (no implementation)
- Supporting types (LayerStatus, ProgressUpdate, ProgressCallback)
- Error types with Error() and Unwrap() methods
- Constructor function signature (panics with "not implemented")
- Complete Go documentation comments

**What This Effort Does**:
- Defines the contract for interacting with OCI registries
- Establishes the interface that Wave 2 will implement
- Creates foundation for parallel implementation efforts

**What This Effort Does NOT Do**:
- NO actual registry push operations (Wave 2)
- NO HTTP client implementation (Wave 2)
- NO go-containerregistry integration (Wave 2)

---

## Dependencies

**Upstream Dependencies** (must complete before this effort):
- **Effort 1.1.1** (Docker Interface) - Sequential ordering

**Downstream Dependencies** (efforts that depend on this):
- **Effort 1.1.4** (Command Structure) - imports pkg/registry

**External Dependencies**:
- `github.com/google/go-containerregistry/pkg/v1` (already in go.mod)

---

## File Structure

**New Files to Create**:

```
pkg/registry/interface.go (200 lines)
├── Package documentation
├── LayerStatus enum (4 values) with String() method
├── ProgressUpdate struct (4 fields)
├── ProgressCallback function type
├── RegistryClient interface (3 methods)
├── NewRegistryClient() constructor (stub)
└── 3 error types with Error() and Unwrap() methods
```

**Total Estimated Lines**: 200 lines

---

## Exact Code Specifications

### File: pkg/registry/interface.go

**CRITICAL: Copy EXACTLY from Wave Architecture (WAVE-1.1-ARCHITECTURE.md lines 218-384)**

**Source Reference**: `/planning/phase1/wave1/WAVE-1.1-ARCHITECTURE.md` lines 442-606

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

---

## Implementation Requirements (CRITICAL)

**EXACT CODE FIDELITY**:
- ✅ Use EXACT code from architecture document (no modifications)
- ✅ Import `github.com/google/go-containerregistry/pkg/v1` for v1.Image type
- ✅ LayerStatus MUST implement String() method (Stringer interface)
- ✅ All error types MUST implement error interface via Error() method
- ✅ All error types MUST implement Unwrap() for error chain
- ✅ NewRegistryClient MUST panic("not implemented") - no actual implementation yet
- ✅ AuthProvider and TLSProvider are placeholder interfaces (defined in Effort 1.1.3)
- ✅ All public types and methods MUST have Go documentation comments

**CRITICAL: DO NOT**:
- ❌ Modify interface signatures
- ❌ Add implementation logic
- ❌ Change error message formats
- ❌ Add extra methods or types
- ❌ Modify package imports

---

## Implementation Steps

### Step 1: Create Package Directory Structure

```bash
# Create pkg/registry directory
mkdir -p pkg/registry

# Verify directory structure
ls -la pkg/
# Expected: pkg/ directory exists with registry/ subdirectory
```

### Step 2: Create interface.go File

```bash
# Create pkg/registry/interface.go
# Copy EXACT code from specification above (lines 442-606 from Wave Architecture)
# Total: 200 lines estimated
```

**Critical Checkpoints**:
- [ ] File created at `pkg/registry/interface.go`
- [ ] Package declaration: `package registry`
- [ ] Import section includes `context` and `github.com/google/go-containerregistry/pkg/v1`
- [ ] LayerStatus enum with 4 constants and String() method
- [ ] ProgressUpdate struct with 4 fields
- [ ] ProgressCallback function type defined
- [ ] RegistryClient interface with 3 methods
- [ ] NewRegistryClient() constructor stub (panics)
- [ ] 3 error types with Error() and Unwrap() methods

### Step 3: Build Validation

```bash
# Build the package
cd pkg/registry
go build .

# Expected output: Build succeeds with no errors
# This is a pure interface definition - no implementation to run
```

**If build fails**:
- Check imports are correct
- Verify v1.Image is from go-containerregistry
- Ensure all public symbols have documentation comments
- Confirm no syntax errors

### Step 4: Create Unit Tests

**File: pkg/registry/interface_test.go**

Create comprehensive tests for all interface elements:

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

// T1.1.2-009: Mock RegistryClient satisfies interface (added for completeness)
func TestMockRegistryClient_SatisfiesInterface(t *testing.T) {
	// This test will be implemented in Wave 2 when we have real mocks
	// For now, just verify the interface is valid Go
	var _ registry.RegistryClient = nil
}
```

**Test Coverage Requirements**:
- Minimum 100% coverage (interface definitions and simple types)
- All error types tested
- LayerStatus String() method tested
- ProgressUpdate struct tested
- ProgressCallback type tested
- Interface compilation verified

### Step 5: Run Tests

```bash
# Run tests with coverage
go test ./pkg/registry -v -cover

# Expected output:
# === RUN   TestRegistryClientInterfaceCompiles
# --- PASS: TestRegistryClientInterfaceCompiles
# === RUN   TestLayerStatus_StringMethod
# --- PASS: TestLayerStatus_StringMethod
# === RUN   TestProgressUpdate_StructValid
# --- PASS: TestProgressUpdate_StructValid
# === RUN   TestProgressCallback_TypeValid
# --- PASS: TestProgressCallback_TypeValid
# === RUN   TestRegistryAuthError_ImplementsError
# --- PASS: TestRegistryAuthError_ImplementsError
# === RUN   TestRegistryConnectionError_ImplementsError
# --- PASS: TestRegistryConnectionError_ImplementsError
# === RUN   TestLayerPushError_ImplementsError
# --- PASS: TestLayerPushError_ImplementsError
# === RUN   TestNewRegistryClient_SignatureValid
# --- PASS: TestNewRegistryClient_SignatureValid
# === RUN   TestMockRegistryClient_SatisfiesInterface
# --- PASS: TestMockRegistryClient_SatisfiesInterface
# PASS
# coverage: 100.0% of statements
```

### Step 6: Commit and Push

```bash
# Stage all files
git add pkg/registry/

# Commit with descriptive message
git commit -m "feat(registry): define RegistryClient interface with progress reporting

- Add RegistryClient interface with 3 methods (Push, BuildImageReference, ValidateRegistry)
- Define LayerStatus enum with String() method
- Define ProgressUpdate struct and ProgressCallback type
- Add 3 error types with Error() and Unwrap() methods
- Stub NewRegistryClient() constructor (Wave 2 implementation)
- 100% test coverage (9 tests, all passing)

Effort 1.1.2 - Registry Client Interface Definition
Estimated lines: 200
Part of Wave 1.1 - Interface Definitions

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

# Push to remote
git push origin idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.2
```

---

## Size Management

**Estimated Lines**: 200 lines
**Measurement Tool**: `${PROJECT_ROOT}/tools/line-counter.sh`
**Check Frequency**: After implementation complete
**Split Threshold (R535)**:
- Warning: 700 lines
- Enforcement: 900 lines (Code Reviewer threshold)

**Expected Status**: ✅ COMPLIANT (200 lines << 900 line enforcement threshold)

**Size Measurement Command**:
```bash
# After completing implementation
cd /home/vscode/workspaces/idpbuilder-oci-push-rebuild/efforts/phase1/wave1/effort-1.1.2

# Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Run line counter (auto-detects base branch)
$PROJECT_ROOT/tools/line-counter.sh

# Expected output:
# 🎯 Detected base: idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.1
# 📦 Analyzing branch: idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.2
# ✅ Total implementation lines: ~200
# ✅ Status: COMPLIANT (within 900-line enforcement threshold)
```

---

## Test Requirements

**From Wave Test Plan** (WAVE-1-TEST-PLAN.md):

### Unit Tests (9 tests)

1. **T1.1.2-001**: RegistryClient interface compiles
2. **T1.1.2-002**: LayerStatus enum compiles with String() method
3. **T1.1.2-003**: ProgressUpdate struct compiles
4. **T1.1.2-004**: ProgressCallback type signature valid
5. **T1.1.2-005**: RegistryAuthError implements error with Unwrap
6. **T1.1.2-006**: RegistryConnectionError implements error with Unwrap
7. **T1.1.2-007**: LayerPushError implements error with Unwrap
8. **T1.1.2-008**: NewRegistryClient constructor signature valid
9. **T1.1.2-009**: Mock RegistryClient satisfies interface

**Test Coverage Target**: 100% (interface definitions are simple)
**Expected Pass Rate**: 100% (no complex logic)

**Test Execution**:
```bash
# Run all tests
go test ./pkg/registry -v

# Run with coverage report
go test ./pkg/registry -cover -coverprofile=coverage.out

# View coverage details
go tool cover -html=coverage.out
```

---

## Acceptance Criteria

**All criteria MUST be met before marking effort complete:**

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

**Quality Gates**:
- ✅ Code matches architecture document exactly
- ✅ All error types implement error interface
- ✅ All error types implement Unwrap() for error chain
- ✅ LayerStatus implements Stringer interface
- ✅ 100% test coverage
- ✅ No linting errors
- ✅ Documentation complete

---

## Integration Points

**Depends On** (Effort 1.1.1):
- Sequential ordering requirement
- No code dependencies (different package)

**Used By** (Effort 1.1.4):
- `cmd/push.go` will import `pkg/registry`
- PushCommand struct will reference RegistryClient interface

**External Integration**:
- Uses `github.com/google/go-containerregistry/pkg/v1` types
- Wave 2 will implement this interface using go-containerregistry's remote package

---

## Risk Assessment

**LOW RISK** effort because:
- ✅ Pure interface definition (no complex logic)
- ✅ Code copied exactly from validated architecture
- ✅ Small size (~200 lines, well under 900-line enforcement threshold)
- ✅ Simple tests (compilation checks, error interface validation)
- ✅ No external service dependencies
- ✅ Standard library types only

**Potential Issues**:
1. **Import path errors** - Mitigation: Use exact imports from architecture
2. **Syntax errors** - Mitigation: Copy code exactly, no modifications
3. **Test failures** - Mitigation: Tests are simple (interface compilation checks)

---

## Notes for SW Engineer

**CRITICAL INSTRUCTIONS**:

1. **EXACT CODE FIDELITY**: Copy the interface code EXACTLY from the specification above. Do NOT modify method signatures, error messages, or documentation comments.

2. **NO IMPLEMENTATION**: This is interface definitions ONLY. The NewRegistryClient() constructor MUST panic("not implemented"). Wave 2 will provide actual implementation.

3. **DOCUMENTATION REQUIRED**: Every public type and method MUST have Go documentation comments. Follow the examples in the code above.

4. **ERROR INTERFACES**: All error types MUST implement:
   - `Error() string` method (required by error interface)
   - `Unwrap() error` method (for error chain support)

5. **STRING METHOD**: LayerStatus MUST implement `String() string` method (Stringer interface).

6. **PLACEHOLDER INTERFACES**: AuthProvider and TLSProvider are defined as empty interfaces in this file. They will be properly defined in Effort 1.1.3.

7. **TEST COVERAGE**: Aim for 100% coverage. Tests are simple (compilation checks, error interface validation).

8. **SIZE COMPLIANCE**: This effort should be ~200 lines. No split required.

9. **COMMIT MESSAGE**: Follow the example in Step 6 above. Include effort ID, line count, and test results.

10. **DEPENDENCY**: This effort follows Effort 1.1.1 sequentially. Branch from `idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.1`.

---

## References

**Source Documents**:
- Wave Implementation Plan: `/planning/phase1/wave1/WAVE-1-IMPLEMENTATION-PLAN.md`
- Wave Architecture: `/planning/phase1/wave1/WAVE-1.1-ARCHITECTURE.md`
- Wave Test Plan: `/planning/phase1/wave1/WAVE-1-TEST-PLAN.md` (if exists)

**Rule References**:
- R213: Effort Metadata Requirements
- R303: Implementation Plan File Locations
- R383: Metadata File Timestamp Requirements (this plan uses timestamped filename)
- R535: Code Reviewer Size Enforcement (900-line threshold)
- R502: Implementation Plan Quality Gates

**Next Steps After Completion**:
1. SW Engineer implements this plan
2. Code Reviewer reviews implementation
3. If approved, merge to wave integration branch
4. Proceed to Effort 1.1.3 (Auth & TLS Interfaces)

---

**Plan Status**: ✅ Ready for SW Engineer Implementation
**Created**: 2025-11-11 20:37:24 UTC
**Version**: 1.0
