# Effort 1.2.1: Docker Client Implementation - Implementation Plan

**Created**: 2025-10-29 06:30:00 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**Effort ID**: 1.2.1
**Phase**: Phase 1 - Foundation & Interfaces
**Wave**: Wave 2 - Core Package Implementations

---

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)

### R213 Infrastructure Metadata

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

**Branch**: `idpbuilder-oci-push/phase1/wave2/effort-1-docker-client`
**Can Parallelize**: Yes
**Parallel With**: [1.2.2, 1.2.3, 1.2.4] (ALL Wave 2 efforts run simultaneously)
**Size Estimate**: 400 lines (MUST be <800)
**Dependencies**: None (all Wave 1 interfaces frozen and available)

---

## Overview

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

---

## File Structure

### New Files to Create

**Implementation Files**:
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

---

## Implementation Steps

### Step 1: Review Wave 1 Interface Definition

**MANDATORY FIRST STEP - Read frozen interfaces**:
```bash
# Read the frozen Docker interface from Wave 1
cat pkg/docker/types.go

# Understand the contract you're implementing:
# - NewClient(ctx) (Client, error)
# - ImageExists(ctx, imageName) (bool, error)
# - GetImage(ctx, imageName) (v1.Image, error)
# - ValidateImageName(imageName) error
# - Close() error
```

**Expected Interface** (from Wave 1 Effort 1):
```go
package docker

import (
    "context"
    v1 "github.com/google/go-containerregistry/pkg/v1"
)

type Client interface {
    ImageExists(ctx context.Context, imageName string) (bool, error)
    GetImage(ctx context.Context, imageName string) (v1.Image, error)
    ValidateImageName(imageName string) error
    Close() error
}
```

**Expected Error Types** (from Wave 1):
```go
type DaemonConnectionError struct { Message string }
type ImageNotFoundError struct { ImageName string }
type ImageConversionError struct { ImageName string, Cause error }
type ValidationError struct { Field string, Message string }
```

### Step 2: Update go.mod with Docker Dependencies

**Add Docker Engine API**:
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning/efforts/phase1/wave2/effort-1-docker-client
go get github.com/docker/docker@v28.2.2+incompatible
go get github.com/docker/go-connections@v0.4.0
go mod tidy
```

**Verify dependencies added**:
- `github.com/docker/docker` v28.2.2+incompatible
- `github.com/docker/go-connections` v0.4.0
- `github.com/google/go-containerregistry` v0.19.0 (already present from Wave 1)

### Step 3: Implement pkg/docker/client.go

**File: pkg/docker/client.go**

**Required Implementation Details**:

**1. Struct Definition**:
```go
type dockerClient struct {
    cli *client.Client
}
```

**2. NewClient() Implementation** (~50 lines):
- Use `client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())`
- Ping daemon to verify connectivity: `cli.Ping(ctx)`
- Return `&DaemonConnectionError{Message: "..."}` if daemon unreachable
- Store Docker Engine client in `dockerClient` struct
- Return implementation of `Client` interface

**3. ImageExists() Implementation** (~40 lines):
- Call `ValidateImageName()` first, return validation error if invalid
- Use `cli.ImageInspectWithRaw(ctx, imageName)`
- Return `(false, nil)` if `client.IsErrNotFound(err)` - NOT an error!
- Return `&DaemonConnectionError{...}` for connection issues
- Return `true` only if inspect succeeds

**4. GetImage() Implementation** (~60 lines):
- Call `ImageExists()` first
- If not exists, return `&ImageNotFoundError{ImageName: imageName}`
- Parse image reference using `daemon.NewTag(imageName)`
- Convert to OCI image using `daemon.Image(ref)`
- Return `&ImageConversionError{ImageName: imageName, Cause: err}` if conversion fails
- Return OCI `v1.Image` compatible with go-containerregistry

**5. ValidateImageName() Implementation** (~40 lines):
- Check for empty string → `&ValidationError{Field: "imageName", Message: "cannot be empty"}`
- Check for command injection attempts: `;`, `|`, `&`, `$`, `` ` ``, `(`, `)`, `<`, `>`, `\`
- Return `&ValidationError{Field: "imageName", Message: "contains invalid characters"}` if found
- Allow valid characters: alphanumeric, dots, slashes, colons, hyphens, underscores
- Helper function `containsString(s string, chars []string) bool` to check for dangerous chars

**6. Close() Implementation** (~20 lines):
- Check for nil client before closing
- Close underlying Docker Engine client connection
- Return any cleanup errors

**Complete Package Structure**:
```go
package docker

import (
    "context"
    "strings"

    "github.com/docker/docker/client"
    "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/daemon"
)

// Implementation here (~400 lines total)
```

### Step 4: Write Tests (TDD - Tests First!)

**File: pkg/docker/client_test.go**

**Test Cases** (from Wave 2 Test Plan):

**A. Constructor Tests**:
- `TestNewClient_Success`: NewClient succeeds with running daemon
- `TestNewClient_DaemonConnectionError`: NewClient fails when daemon stopped

**B. ImageExists Tests**:
- `TestImageExists_ImagePresent`: Returns true for present images (alpine:latest)
- `TestImageExists_ImageNotFound`: Returns false (NOT error) for missing images
- `TestImageExists_ValidationError`: Returns ValidationError for empty/invalid names

**C. GetImage Tests**:
- `TestGetImage_Success`: Retrieves and converts image to OCI v1.Image
- `TestGetImage_ImageNotFoundError`: Returns ImageNotFoundError for non-existent images
- `TestGetImage_ImageConversionError`: Returns ImageConversionError on conversion failure

**D. ValidateImageName Tests**:
- `TestValidateImageName_Valid`: Passes for valid names (myapp:latest, registry.io/repo:tag)
- `TestValidateImageName_CommandInjection`: Rejects dangerous names (command injection attempts)

**E. Close Tests**:
- `TestClose_Success`: Close succeeds and cleans up resources

**F. Edge Case Tests**:
- `TestContextCancellation`: Context cancellation is respected

**Test Coverage Requirements**:
- Minimum 85% code coverage
- All success paths tested
- All failure paths tested
- Edge cases (cancelled context, malformed input) tested
- Security validation (command injection) tested

**Test Setup Requirements**:
```go
package docker

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// Helper: Check if Docker daemon is running
func isDaemonRunning() bool {
    cli, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        return false
    }
    defer cli.Close()
    _, err = cli.Ping(context.Background())
    return err == nil
}

// Skip tests if daemon not available
func requireDockerDaemon(t *testing.T) {
    if !isDaemonRunning() {
        t.Skip("Docker daemon not running - skipping test")
    }
}
```

### Step 5: Size Measurement

**Measure implementation lines**:
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning
tools/line-counter.sh

# Expected output:
# 🎯 Detected base: idpbuilder-oci-push/phase1/wave2/integration
# 📦 Analyzing branch: idpbuilder-oci-push/phase1/wave2/effort-1-docker-client
# ✅ Total implementation lines: ~400
```

**Size Compliance**:
- Target: 400 lines
- Buffer: ±15% (340-460 lines acceptable)
- Hard limit: 800 lines (MUST NOT EXCEED)
- Tests NOT counted (per R007)

**If approaching 700 lines**:
- STOP IMMEDIATELY
- Report to orchestrator
- Do NOT exceed 800 lines

### Step 6: Run Tests and Coverage

**Run unit tests**:
```bash
cd pkg/docker
go test -v -cover

# Expected: All tests pass, coverage ≥85%
```

**Generate coverage report**:
```bash
go test -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
go tool cover -func=coverage.out | grep total
# Expected: total: (statements) 85.0% or higher
```

### Step 7: Linting and Documentation

**Run linters**:
```bash
go vet ./pkg/docker/...
golangci-lint run ./pkg/docker/...

# Expected: No errors
```

**Verify GoDoc**:
- All public types have GoDoc comments
- All public functions have GoDoc comments
- Examples provided for main methods

### Step 8: Commit and Push

**Commit structure**:
```bash
git add pkg/docker/client.go pkg/docker/client_test.go go.mod go.sum
git commit -m "feat(docker): implement Docker client with OCI conversion

- Implement NewClient with daemon ping
- Implement ImageExists with NotFound handling
- Implement GetImage with OCI v1.Image conversion
- Implement ValidateImageName with command injection prevention
- Implement Close with resource cleanup
- Add 12 test cases with 85%+ coverage
- Add Docker Engine API dependencies

Closes: Effort 1.2.1 - Docker Client Implementation
Lines: ~400 (within 400 ±15% estimate)
Coverage: 85%+ (meets Wave 2 Test Plan requirements)"

git push origin idpbuilder-oci-push/phase1/wave2/effort-1-docker-client
```

---

## Size Management

**Estimated Lines**: 400 lines (implementation code only)
**Measurement Tool**: `${PROJECT_ROOT}/tools/line-counter.sh` (find project root first)
**Check Frequency**: After every major function implementation (~50 lines)
**Split Threshold**:
- Warning: 700 lines (approaching limit)
- Hard stop: 800 lines (MUST NOT EXCEED)

**Measurement Commands**:
```bash
# Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ]; then break; fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Measure size
cd /home/vscode/workspaces/idpbuilder-oci-push-planning/efforts/phase1/wave2/effort-1-docker-client
$PROJECT_ROOT/tools/line-counter.sh
```

**Size Tracking**:
- After struct definition: ~50 lines
- After NewClient: ~100 lines
- After ImageExists: ~140 lines
- After GetImage: ~200 lines
- After ValidateImageName: ~240 lines
- After Close: ~260 lines
- After helpers and comments: ~400 lines (target)

---

## Test Requirements

**Minimum Coverage**: 85% (per Wave 2 Test Plan)

**Test Categories** (from Wave 2 Test Plan):

| Test Category | Test Cases | Coverage Target |
|---------------|------------|-----------------|
| Constructor | 2 tests | 100% of NewClient |
| ImageExists | 3 tests | 100% of ImageExists |
| GetImage | 3 tests | 100% of GetImage |
| ValidateImageName | 2 tests | 100% of ValidateImageName |
| Close | 1 test | 100% of Close |
| Edge Cases | 1 test | Context handling |

**Total Test Cases**: 12+ tests

**Test Execution**:
```bash
go test ./pkg/docker -v -cover
# MUST achieve ≥85% coverage
```

---

## Dependencies

### Upstream Dependencies (COMPLETED)
- ✅ Wave 1 Effort 1: Docker interface definition (FROZEN)
- ✅ Integration branch: `idpbuilder-oci-push/phase1/wave2/integration` (CREATED)

### Downstream Dependencies
- None (all Wave 2 efforts are parallel)
- Wave 3 CLI will use this package

### External Library Dependencies
**New Dependencies** (add to go.mod):
- `github.com/docker/docker` v28.2.2+incompatible (Docker Engine API)
- `github.com/docker/go-connections` v0.4.0 (Docker client helpers)

**Existing Dependencies** (from Wave 1):
- `github.com/google/go-containerregistry` v0.19.0 (OCI conversion)
- `github.com/stretchr/testify` v1.10.0 (testing)

### System Dependencies
- Docker daemon running locally
  - Unix socket: `/var/run/docker.sock`
  - Windows named pipe: `npipe:////./pipe/docker_engine`
- DOCKER_HOST environment variable (optional override)

---

## Pattern Compliance

### Go Patterns
- Interface-driven design (implement `docker.Client` interface)
- Error wrapping with custom error types
- Context propagation for cancellation
- Resource cleanup with `Close()` method

### Security Requirements
- Command injection prevention in `ValidateImageName()`
- No dangerous characters allowed in image names
- Proper error classification (don't expose internal details)

### Performance Targets
- Daemon ping should complete in <1 second
- Image inspection should complete in <2 seconds
- Image retrieval time depends on image size (no specific target)

---

## Acceptance Criteria

**MANDATORY - All must pass before Code Review**:

- [ ] All files created/modified as specified
- [ ] All 4 interface methods implemented correctly (ImageExists, GetImage, ValidateImageName, Close)
- [ ] All tests passing (100% pass rate)
- [ ] Code coverage ≥85% (per Wave 2 Test Plan)
- [ ] No linting errors (go vet, golangci-lint)
- [ ] Documentation complete (all public methods have GoDoc)
- [ ] Line count within estimate (400 lines ±15% = 340-460 lines)
- [ ] Integration with go-containerregistry working (v1.Image conversion)
- [ ] Security validation preventing command injection
- [ ] Proper error type usage (Wave 1 errors: DaemonConnectionError, ImageNotFoundError, etc.)
- [ ] Docker dependencies added to go.mod
- [ ] Code committed and pushed to effort branch

**Quality Gates**:
1. **Functionality**: All interface methods work correctly
2. **Testing**: 85%+ coverage with all paths tested
3. **Security**: Command injection prevention validated
4. **Size**: Within 400 ±15% lines (340-460)
5. **Documentation**: Complete GoDoc coverage

---

## References

**Wave 2 Planning Documents**:
- Wave Implementation Plan: `/home/vscode/workspaces/idpbuilder-oci-push-planning/wave-plans/WAVE-2-IMPLEMENTATION.md`
- Wave Architecture: `/home/vscode/workspaces/idpbuilder-oci-push-planning/wave-plans/WAVE-2-ARCHITECTURE.md`
- Wave Test Plan: `/home/vscode/workspaces/idpbuilder-oci-push-planning/wave-plans/WAVE-2-TEST-PLAN.md`

**Wave 1 Interfaces** (frozen references):
- Docker Interface: `efforts/phase1/wave1/effort-1-docker-interface/pkg/docker/types.go`
- Docker Errors: `efforts/phase1/wave1/effort-1-docker-interface/pkg/docker/errors.go`

**External Documentation**:
- Docker Engine API: https://docs.docker.com/engine/api/
- go-containerregistry: https://github.com/google/go-containerregistry
- Docker Go SDK: https://pkg.go.dev/github.com/docker/docker/client

---

## Implementation Checklist

**Pre-Implementation**:
- [ ] Read Wave 1 Docker interface definition
- [ ] Read Wave 2 Architecture (Docker section)
- [ ] Read Wave 2 Test Plan (Docker test cases)
- [ ] Checkout effort branch from integration
- [ ] Verify base branch is correct

**Implementation Phase**:
- [ ] Update go.mod with Docker dependencies
- [ ] Write test stubs (12+ test cases)
- [ ] Implement `dockerClient` struct
- [ ] Implement `NewClient()` with daemon ping
- [ ] Implement `ImageExists()` with NotFound handling
- [ ] Implement `GetImage()` with OCI conversion
- [ ] Implement `ValidateImageName()` with security checks
- [ ] Implement `Close()` with cleanup
- [ ] Implement helper functions
- [ ] Complete test implementations
- [ ] Run tests (all pass, 85%+ coverage)
- [ ] Run linters (no errors)
- [ ] Add GoDoc comments

**Validation Phase**:
- [ ] Measure size (within 340-460 lines)
- [ ] Verify coverage ≥85%
- [ ] Manual testing with real Docker daemon
- [ ] Security validation (command injection tests)
- [ ] Error handling verification
- [ ] Commit and push code

---

## Next Steps

**After Implementation Completion**:
1. Code Reviewer will be spawned for effort review
2. Code Reviewer validates all acceptance criteria
3. If approved: Merge to integration branch
4. If issues found: Fix and re-submit for review
5. After all 4 Wave 2 efforts approved: Wave integration testing

**Parallel Work**:
- This effort (1.2.1) runs in parallel with:
  - Effort 1.2.2: Registry Client Implementation
  - Effort 1.2.3: Authentication Implementation
  - Effort 1.2.4: TLS Configuration Implementation

**No coordination needed** - all efforts are independent until integration phase.

---

## Document Status

**Status**: ✅ READY FOR IMPLEMENTATION
**Created**: 2025-10-29 06:30:00 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**Effort**: 1.2.1 (Docker Client Implementation)
**Wave**: Wave 2 of Phase 1
**Branch**: `idpbuilder-oci-push/phase1/wave2/effort-1-docker-client`
**Base Branch**: `idpbuilder-oci-push/phase1/wave2/integration`

**Compliance**:
- ✅ R213: Complete metadata included
- ✅ R211: Parallelization specified (runs with 1.2.2, 1.2.3, 1.2.4)
- ✅ R341: TDD approach (test plan before implementation)
- ✅ R381: Library versions locked (Docker v28.2.2, go-containerregistry v0.19.0)
- ✅ R383: Plan stored in .software-factory with timestamp
- ✅ Size compliance: 400 lines < 800 hard limit

---

**END OF EFFORT 1.2.1 IMPLEMENTATION PLAN**
