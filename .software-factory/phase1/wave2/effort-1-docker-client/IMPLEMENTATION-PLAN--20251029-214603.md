# Effort 1.2.1: Docker Client Implementation - Detailed Implementation Plan

**Effort ID**: 1.2.1
**Effort Name**: Docker Client Implementation
**Phase**: Phase 1 - Foundation & Interfaces
**Wave**: Wave 2 - Core Package Implementations
**Created**: 2025-10-29 21:46:03 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**Status**: READY FOR IMPLEMENTATION

---

## 🚨 CRITICAL EFFORT INFRASTRUCTURE METADATA (R213)

### Branch Information
**WORKING_DIRECTORY**: `efforts/phase1/wave2/effort-1-docker-client/`
**BRANCH**: `idpbuilder-oci-push/phase1/wave2/effort-1-docker-client`
**BASE_BRANCH**: `idpbuilder-oci-push/phase1/wave2/integration`
**REMOTE**: `https://github.com/jessesanford/idpbuilder.git`

### Effort Metadata
**EFFORT_ID**: `1.2.1`
**EFFORT_NAME**: `Docker Client Implementation`
**PARENT_WAVE**: `WAVE_2`
**PARENT_PHASE**: `PHASE_1`

### Parallelization Strategy
**CAN_PARALLELIZE**: `true`
**PARALLEL_WITH**: `["1.2.2", "1.2.3", "1.2.4"]`
**REASON**: All Wave 2 efforts implement different packages using frozen Wave 1 interfaces. No cross-effort dependencies exist during implementation.

### Dependencies
**DEPENDS_ON**: `[]`
**BLOCKS**: `[]`
**INTEGRATES_WITH**: `["1.2.2-registry-client"]` (registry uses docker images)

### Size Estimates
**ESTIMATED_LINES**: `400`
**COMPLEXITY**: `medium`
**ESTIMATED_DURATION**: `4-6 hours`

---

## Overview

### Purpose

Implement the Docker client package that connects to the Docker daemon, validates image names, checks for image existence, retrieves images in OCI format, and properly cleans up resources.

### What This Effort Accomplishes

- Complete implementation of `docker.Client` interface (frozen in Wave 1)
- Docker daemon connectivity using Docker Engine API
- Image existence checking with proper error classification
- Image retrieval and conversion to OCI v1.Image format
- Image name validation with security checks (command injection prevention)
- Resource cleanup and connection management

### Boundaries - OUT OF SCOPE

- ❌ Building or creating Docker images (only retrieval)
- ❌ Docker compose or multi-container operations
- ❌ Docker network or volume management
- ❌ Image pushing (that's registry package responsibility)
- ❌ Docker registry authentication (that's auth package responsibility)

---

## Detailed Implementation Specifications

### Files to Create

#### 1. `pkg/docker/client.go` (~400 lines)

**Purpose**: Main implementation of Docker client interface

**Package Declaration**:
```go
// Package docker provides Docker daemon integration for OCI image operations.
package docker
```

**Required Imports**:
```go
import (
    "context"
    "io"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    v1 "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/daemon"
)
```

**Struct Definition**:
```go
// dockerClient implements the Client interface using Docker Engine API.
type dockerClient struct {
    cli *client.Client
}
```

**Methods to Implement** (ALL 4 required by Wave 1 interface):

##### Method 1: NewClient()

**Signature**:
```go
func NewClient() (Client, error)
```

**Requirements**:
- Use `client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())`
- Ping daemon to verify connectivity: `cli.Ping(ctx)`
- Return `DaemonConnectionError` if daemon unreachable
- Store Docker Engine client in `dockerClient` struct

**Implementation Details**:
- Connect to Docker daemon using DOCKER_HOST env var or defaults:
  - Unix: `unix:///var/run/docker.sock`
  - Windows: `npipe:////./pipe/docker_engine`
- Verify daemon is running with ping
- Return wrapped error if connection fails

**GoDoc**:
```go
// NewClient creates a new Docker client instance.
//
// The client connects to the Docker daemon using:
//   - DOCKER_HOST environment variable (if set)
//   - Default Unix socket: unix:///var/run/docker.sock
//   - Default Windows named pipe: npipe:////./pipe/docker_engine
//
// Returns:
//   - Client: Docker client interface implementation
//   - error: DaemonConnectionError if daemon is not reachable or not running
//
// Example:
//   client, err := docker.NewClient()
//   if err != nil {
//       return fmt.Errorf("failed to create Docker client: %w", err)
//   }
//   defer client.Close()
```

##### Method 2: ImageExists()

**Signature**:
```go
func (c *dockerClient) ImageExists(ctx context.Context, imageName string) (bool, error)
```

**Requirements**:
- Call `ValidateImageName()` first
- Use `cli.ImageInspectWithRaw(ctx, imageName)`
- Return `(false, nil)` if `client.IsErrNotFound(err)` - NOT an error!
- Return `DaemonConnectionError` for connection issues
- Return `true` only if inspect succeeds

**Critical Behavior**:
- Image not found is NOT an error condition
- Returns `(false, nil)` for non-existent images
- Only return error for actual failures (daemon down, network issues, invalid input)

**GoDoc**:
```go
// ImageExists checks if an image exists in the local Docker daemon.
//
// This method uses the Docker Engine API's ImageInspectWithRaw to check
// for image existence. A NotFound error from Docker indicates the image
// doesn't exist (returns false, nil).
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - imageName: Image name in format "name:tag" (e.g., "myapp:latest")
//
// Returns:
//   - bool: true if image exists, false otherwise
//   - error: DaemonConnectionError if cannot connect to daemon,
//            ValidationError if imageName is malformed
//
// Example:
//   exists, err := client.ImageExists(ctx, "myapp:latest")
//   if err != nil {
//       return fmt.Errorf("failed to check image: %w", err)
//   }
//   if !exists {
//       return fmt.Errorf("image not found in Docker daemon")
//   }
```

##### Method 3: GetImage()

**Signature**:
```go
func (c *dockerClient) GetImage(ctx context.Context, imageName string) (v1.Image, error)
```

**Requirements**:
- Call `ImageExists()` first, return `ImageNotFoundError` if not found
- Parse image reference using `daemon.NewTag(imageName)`
- Convert to OCI image using `daemon.Image(ref)`
- Return `ImageConversionError` if conversion fails
- Return OCI `v1.Image` compatible with go-containerregistry

**Implementation Details**:
- Validate image name first
- Check existence (reuse ImageExists)
- Use go-containerregistry's daemon package for conversion
- The daemon package internally uses Docker's SaveImage API

**GoDoc**:
```go
// GetImage retrieves an image from the Docker daemon and converts it
// to an OCI v1.Image format compatible with go-containerregistry.
//
// This method uses go-containerregistry's daemon package to handle the
// conversion from Docker image format to OCI v1.Image. The daemon package
// internally uses Docker's SaveImage API to export the image.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - imageName: Image name in format "name:tag"
//
// Returns:
//   - v1.Image: OCI-compatible image object
//   - error: ImageNotFoundError if image doesn't exist,
//            DaemonConnectionError if cannot connect,
//            ImageConversionError if conversion fails
//
// Example:
//   image, err := client.GetImage(ctx, "myapp:latest")
//   if err != nil {
//       return fmt.Errorf("failed to retrieve image: %w", err)
//   }
//   // image can now be pushed to registry
```

##### Method 4: ValidateImageName()

**Signature**:
```go
func (c *dockerClient) ValidateImageName(imageName string) error
```

**Requirements**:
- Check for empty string
- Check for command injection attempts: `;`, `|`, `&`, `$`, `` ` ``, `(`, `)`, `<`, `>`, `\`
- Return `ValidationError` with field="imageName" if invalid
- Allow valid characters: alphanumeric, dots, slashes, colons, hyphens, underscores

**Security Considerations**:
- Prevent command injection attacks
- Image names are often user-provided input
- Basic validation - Docker daemon does full OCI spec validation

**GoDoc**:
```go
// ValidateImageName checks if an image name follows the OCI naming specification.
//
// This method validates:
//   - Image name is not empty
//   - No command injection attempts (no semicolons, pipes, etc.)
//   - Basic format check (allows alphanumeric, dots, slashes, colons, hyphens, underscores)
//
// Note: Full OCI spec validation is complex. This provides basic safety checks.
// Docker daemon will perform additional validation.
//
// Parameters:
//   - imageName: Image name to validate
//
// Returns:
//   - error: ValidationError with details if invalid, nil if valid
//
// Example:
//   if err := client.ValidateImageName("myapp:latest"); err != nil {
//       return fmt.Errorf("invalid image name: %w", err)
//   }
```

##### Method 5: Close()

**Signature**:
```go
func (c *dockerClient) Close() error
```

**Requirements**:
- Close underlying Docker Engine client connection
- Check for nil client before closing
- Return any cleanup errors

**GoDoc**:
```go
// Close cleans up Docker client resources and closes connections.
//
// This method closes the underlying HTTP client connection to the Docker daemon.
// It should be called when the client is no longer needed, typically via defer.
//
// Returns:
//   - error: Error if cleanup fails
//
// Example:
//   client, err := NewClient()
//   if err != nil {
//       return err
//   }
//   defer client.Close()
```

**Helper Functions**:

```go
// containsString checks if substr exists in s
func containsString(s, substr string) bool {
    for i := 0; i < len(s); i++ {
        if i+len(substr) <= len(s) && s[i:i+len(substr)] == substr {
            return true
        }
    }
    return false
}
```

---

### Test Files to Create

#### 1. `pkg/docker/client_test.go` (~300 lines, NOT counted per R007)

**Purpose**: Comprehensive unit tests for Docker client implementation

**Test Categories** (from Wave 2 Test Plan):

**A. Constructor Tests**:
- `TestNewClient_Success` - Verify client creation with running daemon
- `TestNewClient_DaemonNotRunning` - Verify DaemonConnectionError when daemon stopped

**B. ImageExists Tests**:
- `TestImageExists_ImagePresent` - Returns true for present images (alpine:latest)
- `TestImageExists_ImageNotPresent` - Returns false (NOT error) for missing images
- `TestImageExists_InvalidImageName` - Returns ValidationError for empty/invalid names

**C. GetImage Tests**:
- `TestGetImage_Success` - Retrieves and converts image to OCI v1.Image
- `TestGetImage_ImageNotFound` - Returns ImageNotFoundError for non-existent images
- `TestGetImage_ConversionError` - Returns ImageConversionError on conversion failure

**D. ValidateImageName Tests**:
- `TestValidateImageName_Valid` - Passes for valid names (myapp:latest, registry.io/repo:tag)
- `TestValidateImageName_Invalid` - Rejects dangerous names (command injection attempts)

**E. Close Tests**:
- `TestClose_Success` - Succeeds and cleans up resources

**F. Edge Case Tests**:
- `TestImageExists_ContextCancellation` - Context cancellation is respected

**Test Coverage Target**: ≥85% (per Wave 2 Test Plan)

**Test Infrastructure**:
- Requires Docker daemon running
- Uses testify/require and testify/assert
- Tests with real Docker client (not mocked)
- Uses alpine:latest as test image (must be pulled)

**See**: `wave-plans/WAVE-2-TEST-PLAN.md` Section 4.1 for detailed test specifications

---

### Files to Modify

#### 1. `go.mod` (+5 lines)

**Purpose**: Add Docker Engine API dependencies

**Additions Required**:
```go
require (
    // NEW for Wave 2:
    github.com/docker/docker v28.2.2+incompatible  // Docker Engine API
    github.com/docker/go-connections v0.4.0        // Docker client helpers

    // Already present from Wave 1:
    github.com/google/go-containerregistry v0.19.0
    github.com/stretchr/testify v1.10.0
)
```

**Commands**:
```bash
go get github.com/docker/docker@v28.2.2
go get github.com/docker/go-connections@v0.4.0
go mod tidy
```

---

## Error Handling Strategy

### Wave 1 Error Types to Use

**From `pkg/docker/errors.go` (Wave 1)**:

#### DaemonConnectionError
**When to use**: Docker daemon is unreachable or not running
```go
return nil, &DaemonConnectionError{
    Cause: err,
}
```

#### ImageNotFoundError
**When to use**: Image doesn't exist in daemon (GetImage only, NOT ImageExists)
```go
return nil, &ImageNotFoundError{
    ImageName: imageName,
}
```

#### ImageConversionError
**When to use**: Conversion from Docker format to OCI v1.Image fails
```go
return nil, &ImageConversionError{
    ImageName: imageName,
    Cause:     err,
}
```

#### ValidationError
**When to use**: Invalid input (empty name, command injection attempt)
```go
return &ValidationError{
    Field:   "imageName",
    Message: "image name cannot be empty",
}
```

### Error Classification Logic

```go
// Example: ImageExists error handling
_, _, err := c.cli.ImageInspectWithRaw(ctx, imageName)
if err != nil {
    // Check if image simply doesn't exist (NOT AN ERROR!)
    if client.IsErrNotFound(err) {
        return false, nil  // Normal case - image not found
    }
    // Any other error is a connection/daemon issue
    return false, &DaemonConnectionError{Cause: err}
}
return true, nil
```

---

## Dependencies

### Upstream Dependencies (COMPLETED)

- ✅ Wave 1 Effort 1: Docker interface definition (`pkg/docker/interface.go`)
- ✅ Integration branch: `idpbuilder-oci-push/phase1/wave2/integration`

### External Library Dependencies

**Required**:
- `github.com/docker/docker` v28.2.2+ (Docker Engine API)
- `github.com/docker/go-connections` v0.4.0 (Docker client helpers)
- `github.com/google/go-containerregistry` v0.19.0 (already present from Wave 1)
- `github.com/stretchr/testify` v1.10.0 (already present from Wave 1)

### System Dependencies

**Required**:
- Docker daemon running locally
  - Unix socket: `/var/run/docker.sock`
  - Windows named pipe: `npipe:////./pipe/docker_engine`
- DOCKER_HOST environment variable (optional override)
- Test images pulled: `docker pull alpine:latest`

---

## Implementation Workflow

### Step 1: Pre-Implementation Reading (MANDATORY)

**Read these documents BEFORE writing ANY code**:
1. ✅ This implementation plan (current document)
2. 📖 Wave 2 Architecture: `wave-plans/WAVE-2-ARCHITECTURE.md`
3. 📖 Wave 2 Test Plan: `wave-plans/WAVE-2-TEST-PLAN.md` Section 4.1
4. 📖 Wave 1 Docker Interface: `efforts/phase1/wave1/effort-1-docker-interface/pkg/docker/interface.go`
5. 📖 Wave 1 Docker Errors: `efforts/phase1/wave1/effort-1-docker-interface/pkg/docker/errors.go`

### Step 2: Environment Setup

```bash
# Navigate to effort directory (R221 compliance)
cd efforts/phase1/wave2/effort-1-docker-client

# Verify correct branch
git branch --show-current  # Should be: idpbuilder-oci-push/phase1/wave2/effort-1-docker-client

# Verify based on integration branch
git log --oneline -5

# Ensure Docker daemon running
docker info

# Pull test images
docker pull alpine:latest

# Add dependencies
go get github.com/docker/docker@v28.2.2
go get github.com/docker/go-connections@v0.4.0
go mod tidy
```

### Step 3: TDD Implementation (R341 Compliance)

**Test-Driven Development Workflow**:

1. **Write Tests FIRST** (from Wave 2 Test Plan):
   ```bash
   # Create test file
   touch pkg/docker/client_test.go

   # Write test cases from Wave 2 Test Plan Section 4.1
   # Start with TestNewClient_Success
   ```

2. **Implement Code to Pass Tests**:
   ```bash
   # Create implementation file
   touch pkg/docker/client.go

   # Implement methods one at a time
   # Run tests after each method
   ```

3. **Iterate**:
   ```bash
   # Run tests frequently
   go test ./pkg/docker -v

   # Check coverage
   go test ./pkg/docker -cover
   ```

### Step 4: Size Measurement (MANDATORY R304)

**Measure size regularly during implementation**:

```bash
# Find project root (R304 requirement)
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state.json" ]; then break; fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
echo "Project root: $PROJECT_ROOT"

# Measure with line counter (R304 - ONLY valid method)
$PROJECT_ROOT/tools/line-counter.sh

# Expected output: ~400 lines
# If approaching 700 lines: WARNING
# If exceeds 800 lines: STOP IMMEDIATELY
```

**Measurement Frequency** (R200):
- After implementing each method
- Every 200 lines written
- Before committing
- Before requesting review

### Step 5: Commit Regularly

```bash
# Commit after each logical unit
git add pkg/docker/client.go
git commit -m "feat(docker): implement NewClient and ImageExists methods

- Add Docker Engine API integration
- Implement daemon ping verification
- Add image existence checking with NotFound handling
- Proper error classification (DaemonConnectionError, ValidationError)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

git push
```

### Step 6: Final Validation

**Before requesting review**:
```bash
# All tests must pass
go test ./pkg/docker -v
# Expected: PASS

# Coverage must be ≥85%
go test ./pkg/docker -cover
# Expected: coverage: 85.x% or higher

# No linting errors
go vet ./pkg/docker
golangci-lint run ./pkg/docker

# Size within estimate
$PROJECT_ROOT/tools/line-counter.sh
# Expected: 340-460 lines (400 ±15%)

# Documentation complete
# All public methods have GoDoc comments
```

---

## Acceptance Criteria

### Functional Requirements

- [ ] All 4 interface methods implemented correctly:
  - [ ] `NewClient()` creates client and pings daemon
  - [ ] `ImageExists()` returns true/false correctly (false is NOT error)
  - [ ] `GetImage()` converts to OCI v1.Image format
  - [ ] `ValidateImageName()` prevents command injection
  - [ ] `Close()` cleans up resources

### Test Requirements

- [ ] All tests passing (100% pass rate)
- [ ] Code coverage ≥85% (per Wave 2 Test Plan)
- [ ] Test categories complete:
  - [ ] Constructor tests (2 tests)
  - [ ] ImageExists tests (3 tests)
  - [ ] GetImage tests (3 tests)
  - [ ] ValidateImageName tests (2 tests)
  - [ ] Close tests (1 test)
  - [ ] Edge case tests (1 test)

### Code Quality Requirements

- [ ] No linting errors (go vet, golangci-lint)
- [ ] Documentation complete (all public methods have GoDoc)
- [ ] Line count within estimate (400 lines ±15% = 340-460 lines)
- [ ] Proper error type usage (Wave 1 errors)
- [ ] Security validation functional (command injection prevention)

### Integration Requirements

- [ ] Integration with go-containerregistry working (v1.Image conversion)
- [ ] Integration with Docker Engine API working (daemon connectivity)
- [ ] Dependency resolution complete (go.mod updated)

---

## Size Compliance

**Estimated Lines**: 400 lines (implementation only)
**Hard Limit**: 800 lines (R220)
**Soft Limit**: 700 lines (warning threshold)

**Breakdown**:
- `client.go`: ~400 lines
  - Imports: ~20 lines
  - Struct: ~5 lines
  - NewClient(): ~30 lines
  - ImageExists(): ~25 lines
  - GetImage(): ~40 lines
  - ValidateImageName(): ~35 lines
  - Close(): ~15 lines
  - Helper functions: ~20 lines
  - GoDoc comments: ~210 lines

**Test Code** (NOT counted per R007):
- `client_test.go`: ~300 lines (excluded from size limits)

**Status**: ✅ COMPLIANT - Well under 800 line limit

---

## Integration with Other Efforts

### Registry Client (Effort 1.2.2)

**Integration Point**: Registry client will use `docker.Client` to get images

**Example Usage**:
```go
// In registry client implementation
dockerClient, err := docker.NewClient()
if err != nil {
    return err
}
defer dockerClient.Close()

image, err := dockerClient.GetImage(ctx, imageName)
if err != nil {
    return err
}

// Now push image to registry
err = registryClient.Push(ctx, image, targetRef, progressCallback)
```

**No direct dependency during implementation**: Efforts run in parallel

### Wave 3 CLI

**Integration Point**: CLI will create Docker client to retrieve local images

**Example Usage**:
```go
// In CLI push command
dockerClient, err := docker.NewClient()
if err != nil {
    fmt.Fprintf(os.Stderr, "Error connecting to Docker daemon: %v\n", err)
    os.Exit(1)
}
defer dockerClient.Close()

exists, err := dockerClient.ImageExists(ctx, imageName)
if err != nil {
    fmt.Fprintf(os.Stderr, "Error checking image: %v\n", err)
    os.Exit(1)
}
if !exists {
    fmt.Fprintf(os.Stderr, "Image not found: %s\n", imageName)
    os.Exit(1)
}
```

---

## Reference Materials

### Architecture Reference

**See**: `wave-plans/WAVE-2-ARCHITECTURE.md` Lines 84-313 for complete Docker client architecture with:
- Detailed method specifications
- Real code examples
- Error handling patterns
- Integration with go-containerregistry

### Test Specifications

**See**: `wave-plans/WAVE-2-TEST-PLAN.md` Lines 230-652 for complete Docker client test plan with:
- 12+ test cases with exact specifications
- Expected behaviors for each test
- Coverage requirements (≥85%)
- Test infrastructure setup

### Wave 1 Interface Contract

**See**: `efforts/phase1/wave1/effort-1-docker-interface/pkg/docker/interface.go` for:
- `Client` interface definition (4 methods)
- Method signatures (MUST match exactly)
- GoDoc specifications
- Wave 1 contract (FROZEN - cannot change)

### Wave 1 Error Types

**See**: `efforts/phase1/wave1/effort-1-docker-interface/pkg/docker/errors.go` for:
- `DaemonConnectionError`
- `ImageNotFoundError`
- `ImageConversionError`
- `ValidationError`
- Error struct definitions
- Error unwrapping support

---

## Troubleshooting

### Common Issues

#### Issue: Docker daemon not running
**Symptom**: `TestNewClient_Success` fails with connection error
**Solution**:
```bash
# Check Docker daemon status
docker info

# Start Docker daemon (systemd)
sudo systemctl start docker

# Start Docker daemon (Docker Desktop)
# Open Docker Desktop application
```

#### Issue: Test image not found
**Symptom**: `TestImageExists_ImagePresent` fails
**Solution**:
```bash
# Pull required test image
docker pull alpine:latest

# Verify image exists
docker images | grep alpine
```

#### Issue: Coverage below 85%
**Symptom**: Coverage check fails in review
**Solution**:
- Add tests for uncovered code paths
- Focus on error paths (often missed)
- Test edge cases (empty strings, nil contexts)
- Run coverage report to find gaps:
  ```bash
  go test ./pkg/docker -coverprofile=coverage.out
  go tool cover -html=coverage.out
  ```

#### Issue: Approaching size limit
**Symptom**: Line counter shows >700 lines
**Solution**:
- Review for unnecessary code
- Remove verbose comments (keep GoDoc)
- Extract helper functions if needed
- DO NOT delete test code (tests don't count)

---

## Next Steps After Completion

### After Implementation Complete

1. **Push all code**:
   ```bash
   git add -A
   git commit -m "feat(docker): complete Docker client implementation

   - Implement all 4 interface methods
   - Add comprehensive unit tests (85%+ coverage)
   - Docker daemon integration working
   - Image validation with security checks
   - OCI v1.Image conversion functional

   🤖 Generated with [Claude Code](https://claude.com/claude-code)

   Co-Authored-By: Claude <noreply@anthropic.com>"
   git push
   ```

2. **Notify orchestrator**: Implementation complete, ready for review

3. **Code Reviewer will**:
   - Verify all acceptance criteria met
   - Check test coverage ≥85%
   - Validate line count within limit
   - Verify interface correctly implemented
   - Test functionality manually

4. **After review approved**:
   - Code merged to integration branch
   - Waits for other Wave 2 efforts to complete
   - Integration testing after all 4 efforts merged

---

## Document Status

**Status**: ✅ READY FOR IMPLEMENTATION
**Created**: 2025-10-29 21:46:03 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**Effort**: 1.2.1 - Docker Client Implementation
**Phase/Wave**: Phase 1 / Wave 2

**Compliance Summary**:
- ✅ R213: Complete metadata with branch, dependencies, parallelization
- ✅ R211: Parallelization strategy specified (parallel with 1.2.2, 1.2.3, 1.2.4)
- ✅ R219: Dependency analysis complete (no blocking dependencies)
- ✅ R341: TDD approach (test plan exists before implementation)
- ✅ R383: Metadata in correct location (.software-factory/ with timestamp)
- ✅ R381: No library version changes (using existing go-containerregistry)
- ✅ R304: Line counter usage mandated for size measurement
- ✅ Size compliance: 400 lines << 800 line limit

**Next Action**: SW Engineer begins implementation following this plan

---

**🚨 CRITICAL R405 AUTOMATION FLAG 🚨**

CONTINUE-CODE-REVIEWER=TRUE REASON=EFFORT_PLAN_COMPLETE

---

**END OF EFFORT 1.2.1 IMPLEMENTATION PLAN**
