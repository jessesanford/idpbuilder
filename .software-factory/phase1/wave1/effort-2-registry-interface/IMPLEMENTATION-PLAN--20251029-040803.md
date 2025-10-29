# Implementation Plan: Registry Client Interface Definition
## Effort 1.1.2 - Phase 1, Wave 1

**Created**: 2025-10-29T04:08:03Z
**Planner**: @agent-code-reviewer
**Effort ID**: 1.1.2
**Phase**: 1 (Foundation & Interfaces)
**Wave**: 1 (Interface & Contract Definitions)

---

## 🚨 EFFORT INFRASTRUCTURE METADATA (R360)

**EFFORT_NAME**: effort-2-registry-interface
**BRANCH**: idpbuilder-oci-push/phase1/wave1/effort-2-registry-interface
**BASE_BRANCH**: idpbuilder-oci-push/phase1/wave1/integration
**PHASE**: 1
**WAVE**: 1
**EFFORT_INDEX**: 2
**PARALLELIZATION**: sequential
**CAN_PARALLELIZE**: No
**PARALLEL_WITH**: None
**DEPENDENCIES**: []

---

## Overview

**Purpose**: Define the registry client interface for pushing OCI images to container registries, including progress reporting types and comprehensive error types.

**Scope**: Interface definitions ONLY - no implementations (Wave 1 contract definition)

**Estimated Size**: 200 lines (implementation code only, excluding tests)

**Expected Outcomes**:
- Complete `registry.Client` interface with 3 methods
- `ProgressCallback` function type for progress reporting
- `ProgressUpdate` struct for progress data
- 4 custom error types with proper wrapping
- Package documentation (GoDoc)
- Mock implementation with callback validation
- 100% test coverage for interfaces

---

## 🔴🔴🔴 REPOSITORY CONTEXT (R251/R309) 🔴🔴🔴

**CRITICAL UNDERSTANDING**:
- ✅ This plan is for the idpbuilder TARGET repo (https://github.com/jessesanford/idpbuilder.git)
- ✅ Implementation will happen in TARGET repo clone
- ✅ NOT in Software Factory planning repo
- ✅ Files reference TARGET repo structure: `pkg/`, `cmd/`, etc.
- ✅ NOT SF structure: `.claude/`, `rule-library/`, etc.

**Working Directory**: `/efforts/phase1/wave1/effort-2-registry-interface/`
**Target Branch**: `idpbuilder-oci-push/phase1/wave1/effort-2-registry-interface`
**Base Branch**: `idpbuilder-oci-push/phase1/wave1/integration`

---

## 🔴🔴🔴 PRE-PLANNING RESEARCH RESULTS (R374 MANDATORY) 🔴🔴🔴

### Existing Interfaces Found

**Search Conducted**: Searched current wave branches and previous waves for existing registry-related interfaces.

| Interface | Location | Signature | Must Implement |
|-----------|----------|-----------|----------------|
| None found | N/A | N/A | N/A |

**Conclusion**: This is the SECOND effort in Wave 1. No existing registry interfaces to reuse.

### Existing Implementations to Reuse

| Component | Location | Purpose | How to Use |
|-----------|----------|---------|------------|
| None | N/A | N/A | N/A |

**Conclusion**: Wave 1 is interface-only. No implementations exist yet.

### APIs Already Defined

| API | Method | Signature | Notes |
|-----|--------|-----------|-------|
| None | N/A | N/A | Second effort establishes Registry API |

### FORBIDDEN DUPLICATIONS (R373)

**No duplications detected** - This effort establishes the registry client contract.

### REQUIRED INTEGRATIONS (R373)

**External Libraries**:
- MUST use `github.com/google/go-containerregistry/pkg/v1` for OCI image types
- MUST NOT create custom OCI image types (use standard v1.Image)

**Interface References**:
- MUST use forward references to `AuthProvider` and `TLSConfigProvider` (defined in Effort 1.1.3)
- Forward references document expected interfaces for clarity

---

## 🔴🔴🔴 EXPLICIT SCOPE CONTROL (R311 - SUPREME LAW) 🔴🔴🔴

### IMPLEMENT EXACTLY:

**Interfaces (1 interface)**:
- `Client` interface with EXACTLY 3 methods (~60 lines including comments)
  - `Push(ctx context.Context, image v1.Image, targetRef string, progressCallback ProgressCallback) error`
  - `BuildImageReference(registryURL, imageName string) (string, error)`
  - `ValidateRegistry(ctx context.Context, registryURL string) error`

**Types (3 types)**:
- `ProgressCallback` function type (~5 lines with comments)
- `ProgressUpdate` struct with 4 fields (~25 lines with comments)
- Forward reference types: `AuthProvider`, `TLSConfigProvider` (~20 lines)

**Functions (1 function signature)**:
- `NewClient(authProvider AuthProvider, tlsConfig TLSConfigProvider) (Client, error)` - MUST panic with "not implemented" (~10 lines)

**Error Types (4 types, ~70 lines total)**:
- `AuthenticationError` struct with Error() and Unwrap() methods
- `NetworkError` struct with Error() and Unwrap() methods
- `RegistryUnavailableError` struct with Error() method
- `PushFailedError` struct with Error() and Unwrap() methods

**Documentation (1 file, ~20 lines)**:
- `doc.go` with package-level documentation

**Tests (3 test files, ~220 lines total - NOT counted toward 800)**:
- `interface_test.go`: Interface compilation and type validation (~50 lines)
- `mock_test.go`: Mock implementation + callback tests (~100 lines)
- `errors_test.go`: 5 error type tests (~70 lines)

**TOTAL IMPLEMENTATION**: ~200 lines (excludes tests per R007)

### DO NOT IMPLEMENT:

❌ NO actual registry push logic (Wave 2)
❌ NO HTTP client setup (Wave 2)
❌ NO layer upload implementation (Wave 2)
❌ NO manifest push implementation (Wave 2)
❌ NO authentication integration (Wave 2)
❌ NO TLS configuration logic (Wave 2)
❌ NO progress tracking implementation (Wave 2)
❌ NO additional helper functions beyond specified
❌ NO logging or metrics
❌ NO configuration management
❌ NO additional interfaces
❌ NO additional error types beyond the 4 specified

---

## File Structure

### Files to Create

**Implementation Files** (count toward 800-line limit):

1. **`pkg/registry/interface.go`** (~110 lines)
   - Package declaration and imports
   - `Client` interface with 3 methods
   - `ProgressCallback` function type
   - `ProgressUpdate` struct
   - Forward references to `AuthProvider` and `TLSConfigProvider`
   - `NewClient()` function signature

2. **`pkg/registry/errors.go`** (~70 lines)
   - 4 error type structs with Error() methods
   - Unwrap() methods for errors that wrap causes

3. **`pkg/registry/doc.go`** (~20 lines)
   - Package-level GoDoc
   - OCI registry focus

**TOTAL IMPLEMENTATION**: ~200 lines

**Test Files** (excluded from 800-line limit per R007):

4. **`pkg/registry/interface_test.go`** (~50 lines)
5. **`pkg/registry/mock_test.go`** (~100 lines)
6. **`pkg/registry/errors_test.go`** (~70 lines)

**TOTAL TESTS**: ~220 lines (NOT counted per R007)

---

## Implementation Steps

### Step 1: Create pkg/registry Directory

```bash
cd /efforts/phase1/wave1/effort-2-registry-interface
mkdir -p pkg/registry
```

### Step 2: Create pkg/registry/interface.go

**Complete file content** (copy EXACTLY - from Wave Implementation Plan lines 405-500):

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

### Step 3: Create pkg/registry/errors.go

**Complete file content** (copy EXACTLY - from Wave Implementation Plan lines 502-560):

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

### Step 4: Create pkg/registry/doc.go

**Complete file content** (copy EXACTLY - from Wave Implementation Plan lines 562-577):

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

### Step 5: Create Test Files

Create all 3 test files following the pattern from Effort 1.1.1, adapted for registry interface.

### Step 6: Run Tests

```bash
cd /efforts/phase1/wave1/effort-2-registry-interface
go test ./pkg/registry/... -v
```

### Step 7: Measure Size

```bash
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ]; then break; fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
cd /efforts/phase1/wave1/effort-2-registry-interface
$PROJECT_ROOT/tools/line-counter.sh
```

### Step 8: Commit and Push

```bash
cd /efforts/phase1/wave1/effort-2-registry-interface
git add pkg/registry/
git commit -m "feat(registry): add Registry client interface definition

Implements Effort 1.1.2 - Registry Client Interface Definition
Phase 1, Wave 1: Interface & Contract Definitions

Added:
- Client interface with 3 methods (Push, BuildImageReference, ValidateRegistry)
- ProgressCallback function type for progress reporting
- ProgressUpdate struct with 4 fields
- 4 error types (AuthenticationError, NetworkError, RegistryUnavailableError, PushFailedError)
- Forward references to AuthProvider and TLSConfigProvider
- Package documentation
- Mock implementation for testing
- Complete test coverage (100%)

Implementation lines: ~200
Test coverage: 100%
All tests passing

Part of Wave 1 contract definition for Phase 1 Wave 2 implementations.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

git push origin idpbuilder-oci-push/phase1/wave1/effort-2-registry-interface
```

---

## Dependencies

### Upstream Dependencies

**None** - Wave 1 efforts are independent (sequential for clarity only)

### Downstream Dependencies

- **Effort 1.1.4**: Command Structure (imports registry.Client interface)

### External Library Dependencies

```go
require (
	github.com/google/go-containerregistry v0.19.0
	github.com/stretchr/testify v1.9.0
)
```

---

## Acceptance Criteria

- [ ] All 6 files created as specified
- [ ] `registry.Client` interface compiles with 3 methods
- [ ] `ProgressCallback` and `ProgressUpdate` types defined
- [ ] All 4 error types implement `error` interface
- [ ] Mock demonstrates callback invocation
- [ ] All tests passing (100% pass rate)
- [ ] Test coverage = 100%
- [ ] GoDoc complete
- [ ] Line count: 200±30 lines

---

## Document Status

**Status**: ✅ READY FOR IMPLEMENTATION
**Created**: 2025-10-29T04:08:03Z
**Planner**: @agent-code-reviewer
**Effort**: 1.1.2 - Registry Client Interface Definition
**Phase**: 1, Wave: 1
**Fidelity**: EXACT (complete code provided)

**Lines**: 200 implementation + 220 test
**Coverage**: 100% required
**Dependencies**: None

---

**END OF IMPLEMENTATION PLAN - EFFORT 1.1.2**
