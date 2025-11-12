# Effort 1.1.1 Implementation Plan: Docker Client Interface Definition

**Created**: 2025-11-11 20:31:52 UTC
**Code Reviewer**: Code Reviewer Agent
**Effort**: 1.1.1 - Docker Client Interface Definition
**Phase**: Phase 1
**Wave**: Wave 1.1 - Interface Definitions

---

## R213 Effort Metadata (MANDATORY - FROM WAVE PLAN)

```json
{
  "effort_id": "1.1.1",
  "effort_name": "Docker Client Interface Definition",
  "branch_name": "idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.1",
  "base_branch": "main",
  "parent_wave": "wave1.1",
  "parent_phase": "phase1",
  "depends_on": [],
  "estimated_lines": 160,
  "complexity": "low",
  "can_parallelize": false,
  "parallel_with": [],
  "tests_required": [
    "T1.1.1-001: DockerClient interface compiles",
    "T1.1.1-002: ImageNotFoundError implements error",
    "T1.1.1-003: DaemonConnectionError implements error with Unwrap",
    "T1.1.1-004: InvalidImageNameError implements error",
    "T1.1.1-005: NewDockerClient constructor signature valid",
    "T1.1.1-006: Mock DockerClient satisfies interface"
  ]
}
```

---

## Overview

**Purpose**: Create the DockerClient interface that will be implemented in Wave 2 to interact with the local Docker daemon for retrieving container images.

**Scope**: Interface definition ONLY - no implementation logic. This effort creates the contract that Wave 2 efforts will implement.

**Estimated Size**: 160 lines (interface + error types + documentation)

**Implementation Time**: ~1-2 hours (straightforward interface definition)

---

## Scope and Boundaries

### IN SCOPE (✅)

- **DockerClient interface** with 4 methods:
  - `ImageExists(ctx context.Context, imageName string) (bool, error)`
  - `GetImage(ctx context.Context, imageName string) (v1.Image, error)`
  - `ValidateImageName(imageName string) error`
  - `Close() error`

- **Constructor function signature**:
  - `NewDockerClient() (DockerClient, error)` - stub implementation that panics

- **3 custom error types**:
  - `ImageNotFoundError` - image does not exist in Docker daemon
  - `DaemonConnectionError` - Docker daemon unreachable (with Unwrap)
  - `InvalidImageNameError` - image name violates OCI spec

- **Complete package documentation**:
  - Package-level godoc comment
  - Method-level documentation with examples
  - Error type documentation

### OUT OF SCOPE (❌)

- ❌ **NO actual Docker daemon interaction** (Wave 2)
- ❌ **NO implementation logic** (Wave 2)
- ❌ **NO subprocess execution** (Wave 2)
- ❌ **NO Docker Engine API calls** (Wave 2)
- ❌ **NO real image retrieval** (Wave 2)

---

## File Structure

### New Files to Create

**File**: `pkg/docker/interface.go` (160 lines)

```
pkg/docker/interface.go
├── Package documentation comment           (5 lines)
├── Imports (context, io, v1)              (6 lines)
├── DockerClient interface definition      (45 lines)
│   ├── ImageExists method + docs          (12 lines)
│   ├── GetImage method + docs             (12 lines)
│   ├── ValidateImageName method + docs    (10 lines)
│   └── Close method + docs                (6 lines)
├── NewDockerClient constructor            (15 lines)
├── ImageNotFoundError type + methods      (10 lines)
├── DaemonConnectionError type + methods   (18 lines)
└── InvalidImageNameError type + methods   (11 lines)
```

**Total Estimated Lines**: 160 lines (matches R213 estimate)

---

## Exact Code Specification

**CRITICAL**: Copy EXACTLY from WAVE-1.1-ARCHITECTURE.md lines 52-159

The Wave Architecture provides the EXACT code to implement. Do NOT modify, do NOT improvise. Copy verbatim.

### Complete File: pkg/docker/interface.go

```go
// Package docker provides an interface to the local Docker daemon for image operations.
package docker

import (
	"context"
	"io"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// DockerClient provides access to the local Docker daemon for retrieving and validating container images.
// Implementations may use the Docker Engine API or fall back to subprocess commands.
type DockerClient interface {
	// ImageExists checks if an image exists in the local Docker daemon.
	// Returns true if the image is found, false otherwise.
	//
	// Example:
	//   exists, err := client.ImageExists(ctx, "myapp:latest")
	//   if err != nil {
	//       return fmt.Errorf("checking image: %w", err)
	//   }
	//   if !exists {
	//       return fmt.Errorf("image not found")
	//   }
	ImageExists(ctx context.Context, imageName string) (bool, error)

	// GetImage retrieves an image from the Docker daemon as a v1.Image.
	// The returned image can be pushed to an OCI registry.
	//
	// Example:
	//   image, err := client.GetImage(ctx, "myapp:latest")
	//   if err != nil {
	//       return fmt.Errorf("getting image: %w", err)
	//   }
	//   // Use image with RegistryClient.Push()
	GetImage(ctx context.Context, imageName string) (v1.Image, error)

	// ValidateImageName validates an image name against the OCI specification.
	// Returns an error if the name is invalid (wrong format, invalid characters, too long).
	//
	// Example:
	//   if err := client.ValidateImageName("My App:v1"); err != nil {
	//       // Handle invalid name (spaces not allowed)
	//   }
	ValidateImageName(imageName string) error

	// Close releases the Docker daemon connection and any associated resources.
	// Must be called when the client is no longer needed.
	//
	// Example:
	//   defer client.Close()
	Close() error
}

// NewDockerClient creates a new Docker client.
// Attempts to connect to the Docker Engine API first.
// Falls back to subprocess execution (docker save) if the API is unavailable.
//
// Returns an error if both the API and subprocess methods fail.
//
// Example:
//   client, err := docker.NewDockerClient()
//   if err != nil {
//       return fmt.Errorf("creating docker client: %w", err)
//   }
//   defer client.Close()
func NewDockerClient() (DockerClient, error) {
	// Implementation will be provided in Wave 2
	panic("not implemented")
}

// ImageNotFoundError indicates that the requested image does not exist in the Docker daemon.
type ImageNotFoundError struct {
	ImageName string
}

func (e *ImageNotFoundError) Error() string {
	return "image not found: " + e.ImageName
}

// DaemonConnectionError indicates that the Docker daemon is unreachable.
type DaemonConnectionError struct {
	Endpoint string
	Cause    error
}

func (e *DaemonConnectionError) Error() string {
	if e.Cause != nil {
		return "cannot connect to Docker daemon at " + e.Endpoint + ": " + e.Cause.Error()
	}
	return "cannot connect to Docker daemon at " + e.Endpoint
}

func (e *DaemonConnectionError) Unwrap() error {
	return e.Cause
}

// InvalidImageNameError indicates that an image name violates the OCI specification.
type InvalidImageNameError struct {
	ImageName string
	Reason    string
}

func (e *InvalidImageNameError) Error() string {
	return "invalid image name '" + e.ImageName + "': " + e.Reason
}
```

---

## Implementation Steps

### Step 1: Create Package Directory

```bash
# Navigate to effort directory (R221: MUST cd in every command)
cd /home/vscode/workspaces/idpbuilder-oci-push-rebuild/efforts/phase1/wave1/effort-1.1.1

# Create pkg/docker directory
mkdir -p pkg/docker

# Verify creation
ls -la pkg/docker/
```

### Step 2: Create interface.go

```bash
# Create the interface file
# Copy EXACT code from WAVE-1.1-ARCHITECTURE.md lines 52-159
# See "Complete File" section above
```

**CRITICAL**: Do NOT type the code manually. Copy EXACTLY from the architecture document to ensure:
- Correct method signatures
- Correct error handling patterns
- Correct documentation format
- Correct import statements

### Step 3: Verify Compilation

```bash
# Navigate to effort directory (R221)
cd /home/vscode/workspaces/idpbuilder-oci-push-rebuild/efforts/phase1/wave1/effort-1.1.1

# Verify Go code compiles
go build ./pkg/docker/

# Expected: Successful compilation (no errors)
```

### Step 4: Run Tests (Reference Wave Test Plan)

The Wave Test Plan defines 6 tests for this effort:

**Tests from WAVE-1-TEST-PLAN.md**:

1. **T1.1.1-001**: DockerClient interface compiles
   - Verify interface is valid Go type
   - Mock type can satisfy interface

2. **T1.1.1-002**: ImageNotFoundError implements error
   - Error() method returns correct format
   - Implements error interface

3. **T1.1.1-003**: DaemonConnectionError implements error with Unwrap
   - Error() method includes endpoint and cause
   - Unwrap() returns original cause
   - Supports error chain unwrapping

4. **T1.1.1-004**: InvalidImageNameError implements error
   - Error() includes image name and reason
   - Implements error interface

5. **T1.1.1-005**: NewDockerClient constructor signature valid
   - Constructor exists with correct signature
   - Returns (DockerClient, error)
   - Stub implementation panics (expected in Wave 1)

6. **T1.1.1-006**: Mock DockerClient satisfies interface
   - Mock implementation created in `tests/mocks/`
   - Mock satisfies all interface methods
   - Can be used in future tests

**Test execution**:
```bash
# Run all Docker package tests
go test ./pkg/docker/...

# Expected: All 6 tests PASS
```

**NOTE**: Tests are INTERFACE ONLY - no implementation logic tested. All tests should PASS immediately after creating interface.go because we're validating type definitions, not behavior.

### Step 5: Measure Size

```bash
# Navigate to effort directory (R221)
cd /home/vscode/workspaces/idpbuilder-oci-push-rebuild/efforts/phase1/wave1/effort-1.1.1

# Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Use ONLY the line counter tool (R304)
$PROJECT_ROOT/tools/line-counter.sh

# Expected: ~160 lines (interface + error types + docs)
# Tool will auto-detect base branch
```

**Size Decision Logic (R535)**:
- Expected: 160 lines (well within 800-line limit)
- Status: COMPLIANT
- No split required

### Step 6: Commit and Push

```bash
# Navigate to effort directory (R221)
cd /home/vscode/workspaces/idpbuilder-oci-push-rebuild/efforts/phase1/wave1/effort-1.1.1

# Stage the interface file
git add pkg/docker/interface.go

# Commit with descriptive message
git commit -m "feat(docker): Add DockerClient interface definition

- Define DockerClient interface with 4 methods
- Add ImageNotFoundError, DaemonConnectionError, InvalidImageNameError
- Add NewDockerClient constructor (stub)
- Complete package documentation
- 160 lines total

Effort: 1.1.1
Wave: 1.1 - Interface Definitions
Phase: 1

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

# Push to remote
git push origin idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.1
```

---

## Test Requirements

### From Wave Test Plan (WAVE-1-TEST-PLAN.md)

**Test Category**: Interface Contract Tests

**Total Tests**: 6 tests (all PASS expected immediately)

**Test Files to Create**:
1. `pkg/docker/interface_test.go` - Interface compilation tests (5 tests)
2. `pkg/docker/errors_test.go` - Error type tests (3 tests)

**Test Coverage Target**: 100% (interface definitions are simple)

**Test Execution**:
```bash
# Run all Docker package tests
go test ./pkg/docker/...

# Run with coverage
go test -cover ./pkg/docker/...

# Expected: 100% coverage (all lines are interface/error definitions)
```

### Test Creation Guidance

**Interface Tests** (pkg/docker/interface_test.go):
- Verify DockerClient interface compiles
- Verify all 4 methods present in interface
- Verify method signatures correct (parameters, return types)
- Verify constructor function exists
- Verify mock can satisfy interface

**Error Tests** (pkg/docker/errors_test.go):
- Verify each error type implements error interface
- Verify Error() methods return expected format
- Verify DaemonConnectionError.Unwrap() works
- Verify error messages include expected information

**Example Test** (from Wave Test Plan):
```go
func TestDockerClientInterface_Compiles(t *testing.T) {
    // Verify interface is valid Go type
    var _ DockerClient = (*mockDockerClient)(nil)
}

func TestImageNotFoundError_ImplementsError(t *testing.T) {
    err := &ImageNotFoundError{ImageName: "test:v1"}

    // Verify implements error interface
    var _ error = err

    // Verify Error() method returns expected format
    expected := "image not found: test:v1"
    if err.Error() != expected {
        t.Errorf("Expected '%s', got '%s'", expected, err.Error())
    }
}
```

---

## Acceptance Criteria

### Code Quality Checklist

- [ ] **File created**: `pkg/docker/interface.go` exists
- [ ] **Interface complete**: DockerClient has all 4 methods
- [ ] **Methods documented**: All methods have godoc comments with examples
- [ ] **Error types defined**: ImageNotFoundError, DaemonConnectionError, InvalidImageNameError
- [ ] **Error() methods**: All error types implement error interface
- [ ] **Unwrap() method**: DaemonConnectionError implements error unwrapping
- [ ] **Constructor stub**: NewDockerClient() exists and panics
- [ ] **Package documentation**: Package-level godoc comment present
- [ ] **Compilation**: `go build ./pkg/docker/` succeeds
- [ ] **Tests pass**: All 6 tests from Wave Test Plan PASS
- [ ] **Test coverage**: 100% coverage achieved
- [ ] **Size compliant**: ~160 lines (within limit)
- [ ] **Committed**: Code committed and pushed to branch
- [ ] **No implementation**: Constructor contains ONLY `panic("not implemented")`

### Code Review Criteria

**Code Reviewer will verify**:
- ✅ Exact code from architecture document (no modifications)
- ✅ All method signatures match specification
- ✅ All documentation matches architecture
- ✅ Error types implement error interface correctly
- ✅ DaemonConnectionError.Unwrap() returns Cause
- ✅ Constructor function panics (no implementation)
- ✅ Compilation succeeds
- ✅ Tests pass
- ✅ Size within limit

---

## Dependencies

### No Effort Dependencies

**This is the first effort in Wave 1** - no dependencies on other efforts.

**Can Parallelize**: FALSE (from R213 metadata)
- This effort must complete before other Wave 1 efforts can use DockerClient interface references
- However, RegistryClient (Effort 1.1.2) and AuthProvider/TLSProvider (Effort 1.1.3) can be implemented in parallel after this completes

### External Dependencies

**Go Modules Required**:
```go
// From go.mod (these should already exist in project)
github.com/google/go-containerregistry v0.19.0  // For v1.Image type
```

**Standard Library**:
- `context` - for Context parameter in methods
- `io` - for io.Reader (unused in interface, but imported for future use)

**No Additional Dependencies**: All types used are from standard library or go-containerregistry (already in project).

---

## Size Management

### Size Measurement (R304 - MANDATORY TOOL USAGE)

**ONLY valid measurement**:
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-rebuild/efforts/phase1/wave1/effort-1.1.1

# Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Use line counter tool (auto-detects base branch)
$PROJECT_ROOT/tools/line-counter.sh
```

**Expected Output**:
```
🎯 Detected base: main
📦 Analyzing branch: idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.1
✅ Total implementation lines: 160
⚠️  Note: Tests, demos, docs, configs NOT included
```

### Size Compliance (R535)

**Code Reviewer Enforcement Threshold**: 900 lines
**SW Engineer Target**: 800 lines
**This Effort Estimate**: 160 lines

**Status**: COMPLIANT (160 << 800)

**Split Decision**: NO SPLIT REQUIRED

---

## Pattern Compliance

### Go Idioms and Patterns

**Interface Design Patterns**:
- ✅ Small, focused interface (4 methods only)
- ✅ Context-aware (ImageExists, GetImage accept context.Context)
- ✅ Standard error handling (return error as last value)
- ✅ Resource cleanup (Close() method)

**Error Handling Patterns**:
- ✅ Custom error types (ImageNotFoundError, etc.)
- ✅ Error() method implementation (error interface)
- ✅ Error unwrapping (DaemonConnectionError.Unwrap())
- ✅ Descriptive error messages (include context)

**Documentation Patterns**:
- ✅ Package-level godoc comment
- ✅ Interface documentation
- ✅ Method documentation with examples
- ✅ Error type documentation

**Constructor Pattern**:
- ✅ NewDockerClient() function
- ✅ Returns (DockerClient, error)
- ✅ Stub implementation (panics in Wave 1)

### go-containerregistry Integration

**Standard Types Used**:
- `v1.Image` - container image type from go-containerregistry
- Future Wave 2 implementation will use `authn.Authenticator` (from pkg/auth)
- Future Wave 2 implementation will use `remote` package for Docker daemon operations

---

## Integration Points

### How This Effort Integrates

**Provides to Other Efforts**:
- **Effort 1.1.4 (Command Structure)**: Imports DockerClient interface for PushCommand
- **Wave 2 Efforts**: Implementation will create real DockerClient

**Consumes from Other Efforts**:
- **NONE** - This is a foundational effort with no dependencies

**Integration Validation**:
After this effort completes, Effort 1.1.4 should be able to:
```go
import "github.com/jessesanford/idpbuilder/pkg/docker"

// In cmd/push.go
type PushCommand struct {
    dockerClient docker.DockerClient  // Uses interface from this effort
    // ...
}
```

---

## Success Metrics

### Definition of Done

**This effort is complete when**:

1. ✅ **File created**: `pkg/docker/interface.go` exists
2. ✅ **Compilation**: `go build ./pkg/docker/` succeeds
3. ✅ **Tests pass**: All 6 Wave Test Plan tests PASS
4. ✅ **Coverage**: 100% test coverage achieved
5. ✅ **Size**: ~160 lines (measured with line-counter.sh)
6. ✅ **Committed**: Code committed and pushed
7. ✅ **Code review**: Code Reviewer approves (ACCEPTED status)
8. ✅ **Documentation**: godoc generation succeeds

**Code Reviewer Approval Criteria**:
- All acceptance criteria met
- No implementation logic (constructor panics)
- Exact code from architecture document
- Tests pass
- Size compliant

---

## Troubleshooting

### Common Issues

**Issue**: `package v1 is not in GOROOT`
**Solution**: Run `go mod download` to fetch go-containerregistry

**Issue**: `undefined: DockerClient`
**Solution**: Check package name is `package docker`, not `package main`

**Issue**: Tests fail with "method not found"
**Solution**: Verify all 4 methods defined in interface (ImageExists, GetImage, ValidateImageName, Close)

**Issue**: Error types don't implement error
**Solution**: Verify Error() method signature: `func (e *ErrorType) Error() string`

**Issue**: Size exceeds estimate
**Solution**: Verify copying ONLY interface.go, not creating additional files

---

## R381 Library Version Consistency

**CRITICAL**: All library versions are LOCKED per R381

### Locked Dependencies (IMMUTABLE)

From project go.mod (already present):
```
github.com/google/go-containerregistry v0.19.0  # LOCKED - DO NOT UPDATE
```

**NEVER suggest**:
- ❌ Updating to latest version
- ❌ Using version ranges (^, ~, >=)
- ❌ Changing any existing dependency version

**ONLY add NEW dependencies** (none required for this effort)

---

## Next Steps

**After This Effort Completes**:

1. **Code Review**: Code Reviewer validates implementation
2. **Merge**: Merge to wave integration branch
3. **Effort 1.1.2**: RegistryClient interface can reference DockerClient
4. **Effort 1.1.4**: Command structure can import DockerClient

**Wave 1 Completion**: After all 4 efforts (1.1.1 - 1.1.4) complete, Wave 2 can begin implementation.

---

## References

**Source Documents**:
- `/home/vscode/workspaces/idpbuilder-oci-push-rebuild/planning/phase1/wave1/WAVE-1-IMPLEMENTATION-PLAN.md` (lines 42-200)
- `/home/vscode/workspaces/idpbuilder-oci-push-rebuild/planning/phase1/wave1/WAVE-1.1-ARCHITECTURE.md` (lines 52-159)
- `/home/vscode/workspaces/idpbuilder-oci-push-rebuild/efforts/phase1/wave1/effort-1.1.1/planning/phase1/wave1/WAVE-1-TEST-PLAN.md`

**Key Rules**:
- R213: Effort Metadata Requirements
- R221: CD to effort directory in every command
- R303: Save plans in .software-factory with timestamps
- R304: Use ONLY line-counter.sh for size measurement
- R381: Library version consistency (locked versions)
- R383: ALL metadata in .software-factory with timestamps
- R502: Implementation Plan Quality Gates (EXACT fidelity)
- R535: Code Reviewer enforcement at 900 lines

---

**Plan Created**: 2025-11-11 20:31:52 UTC
**Plan Status**: Ready for SW Engineer Implementation
**Expected State Transition**: Code Reviewer creates plan → Orchestrator spawns SW Engineer
