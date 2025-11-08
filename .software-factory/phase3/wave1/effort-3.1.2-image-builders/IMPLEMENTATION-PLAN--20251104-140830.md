# Effort 3.1.2: Test Image Builders - Implementation Plan

**Status**: Ready for Implementation
**Created**: 2025-11-04 14:08:30 UTC
**Created By**: code-reviewer (EFFORT_PLAN_CREATION state)
**Effort**: Phase 3, Wave 1, Effort 3.1.2

---

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)

### R213 Metadata
```yaml
effort_id: "3.1.2"
effort_name: "Test Image Builders"
estimated_lines: 300
complexity: "Medium"
dependencies: ["effort:3.1.1"]
files_touched:
  - "test/harness/image_builder.go"
  - "test/fixtures/Dockerfile.multilayer"
  - "test/fixtures/Dockerfile.minimal"
branch_name: "idpbuilder-oci-push/phase3/wave1/effort-3.1.2-image-builders"
workspace: "efforts/phase3/wave1/effort-3.1.2-image-builders"
parallelization: "SEQUENTIAL"
base_branch: "idpbuilder-oci-push/phase3/wave1/effort-3.1.1-test-harness"
```

---

## Overview

**Purpose**: Implement test image builders that create Docker images with specific characteristics (layer count, size, architecture) for comprehensive testing scenarios.

**Scope**: Create utilities to dynamically build test images with configurable properties, enabling comprehensive integration tests for the push command.

**Dependencies**:
- **Effort 3.1.1**: Test Harness Infrastructure (MUST be complete before starting)
- **Phase 1**: Docker client package (`pkg/docker`)
- **External**: Docker daemon must be available

**Integration Points**:
- **Used By**: All integration tests (Efforts 3.1.3, 3.1.4)
- **Requires**: Test harness environment from Effort 3.1.1
- **Provides**: Test image builders for various testing scenarios

---

## File Structure

### Files to Create

#### 1. `test/harness/image_builder.go` (250 lines)
**Purpose**: Build test images with configurable properties

**Key Components**:
- `ImageConfig` struct defining test image characteristics
- `BuildTestImage()` function for building images
- `generateDockerfile()` function for dynamic Dockerfile creation
- `generateTestFiles()` function for creating layer files

**Responsibilities**:
- Generate Dockerfiles dynamically based on configuration
- Create temporary build directories
- Generate appropriately sized test files for layers
- Execute Docker build operations
- Clean up build directories
- Return image references for testing

#### 2. `test/fixtures/Dockerfile.multilayer` (25 lines)
**Purpose**: Template Dockerfile for multi-layer test images

**Contains**:
- Alpine-based image with multiple RUN commands
- Multiple COPY operations for layer creation
- Test data files for each layer
- Labels for test identification

#### 3. `test/fixtures/Dockerfile.minimal` (10 lines)
**Purpose**: Template for minimal test image

**Contains**:
- Scratch-based minimal image
- Single binary/file
- Minimal layer count for edge case testing

---

## Implementation Details

### Part 1: Image Builder Core (150 lines)

**File**: `test/harness/image_builder.go`

#### Data Structures

```go
package harness

import (
    "context"
    "fmt"
    "io"
    "os"
    "path/filepath"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/pkg/archive"
    "github.com/google/go-containerregistry/pkg/v1"
)

// ImageConfig defines test image characteristics
type ImageConfig struct {
    Name       string  // Image name (without tag)
    Tag        string  // Image tag
    Layers     int     // Number of layers to create
    SizeMB     int     // Approximate total size in megabytes
    Arch       string  // Architecture (amd64, arm64, etc.)
}

// BuildResult contains information about built image
type BuildResult struct {
    ImageRef    string  // Full image reference (name:tag)
    ImageID     string  // Docker image ID
    LayerCount  int     // Actual number of layers
    SizeBytes   int64   // Actual size in bytes
}
```

#### Main Build Function

```go
// BuildTestImage creates a Docker image for testing
//
// This function generates a temporary Dockerfile and build context,
// builds the image using the Docker daemon, and returns information
// about the built image.
//
// The function will:
// 1. Create temporary build directory
// 2. Generate Dockerfile based on config
// 3. Generate test files for each layer
// 4. Execute Docker build
// 5. Clean up build directory
// 6. Return build result
//
// Example:
//   config := ImageConfig{
//       Name: "testapp",
//       Tag: "v1",
//       Layers: 3,
//       SizeMB: 10,
//       Arch: "amd64",
//   }
//   result, err := env.BuildTestImage(ctx, config)
func (env *TestEnvironment) BuildTestImage(ctx context.Context, config ImageConfig) (*BuildResult, error) {
    // Implementation steps:
    // 1. Validate configuration
    // 2. Create temporary build directory
    // 3. Generate Dockerfile
    // 4. Generate test files
    // 5. Create tar archive for build context
    // 6. Call Docker API ImageBuild()
    // 7. Stream build output to log
    // 8. Wait for build completion
    // 9. Verify image exists
    // 10. Clean up build directory
    // 11. Return build result
}
```

#### Helper Functions

```go
// generateDockerfile creates Dockerfile content with specified layers
//
// The generated Dockerfile will create the requested number of layers
// by using RUN and COPY commands appropriately.
func generateDockerfile(config ImageConfig) string {
    // Implementation:
    // - Start with FROM alpine:latest
    // - Add RUN/COPY commands to create layers
    // - Each layer adds files and runs commands
    // - Add LABEL instructions for test metadata
    // - Add CMD instruction
}

// generateTestFiles creates files for each layer
//
// Files are sized to achieve the approximate total size
// specified in the configuration.
func generateTestFiles(buildDir string, config ImageConfig) error {
    // Implementation:
    // - Calculate size per layer
    // - Create layer-N.dat files with random data
    // - Use /dev/urandom or similar for data generation
    // - Verify file sizes match requirements
}

// cleanupBuildDir removes temporary build directory
func cleanupBuildDir(buildDir string) error {
    // Remove all files and directory
}
```

### Part 2: Image Building Execution (100 lines)

#### Docker Build Integration

```go
// executeBuild runs the Docker build operation
func (env *TestEnvironment) executeBuild(ctx context.Context, buildContext io.Reader, imageName string) (string, error) {
    // Implementation:
    // 1. Configure build options
    // 2. Call Docker API ImageBuild()
    // 3. Stream build output
    // 4. Check for errors
    // 5. Return image ID
}

// verifyImageBuilt checks that the image exists in Docker daemon
func (env *TestEnvironment) verifyImageBuilt(ctx context.Context, imageRef string) (*BuildResult, error) {
    // Implementation:
    // 1. Call Docker API ImageInspect()
    // 2. Extract layer count
    // 3. Extract size
    // 4. Return BuildResult
}

// streamBuildOutput reads and logs Docker build output
func streamBuildOutput(reader io.ReadCloser) error {
    // Implementation:
    // - Read JSON messages from Docker
    // - Log progress messages
    // - Detect errors
    // - Return error if build failed
}
```

### Part 3: Fixture Dockerfiles

#### test/fixtures/Dockerfile.multilayer (25 lines)

```dockerfile
FROM alpine:latest

# Layer 1: Base utilities
RUN apk add --no-cache curl wget

# Layer 2: Test data
COPY layer-0.dat /data/
RUN echo "Layer 0 complete" > /data/layer-0.txt

# Layer 3: More test data
COPY layer-1.dat /data/
RUN echo "Layer 1 complete" > /data/layer-1.txt

# Layer 4: Final layer
COPY layer-2.dat /data/
RUN echo "Layer 2 complete" > /data/layer-2.txt

LABEL test.image.name="multilayer"
LABEL test.image.layers="3"

CMD ["sh", "-c", "echo 'Test image running'"]
```

#### test/fixtures/Dockerfile.minimal (10 lines)

```dockerfile
FROM scratch
COPY hello /
CMD ["/hello"]
```

---

## Implementation Steps (Task Checklist)

### Phase 1: Setup and Structure (30 minutes)
- [ ] Create `test/harness/image_builder.go` file
- [ ] Import required packages (docker/docker, go-containerregistry)
- [ ] Define `ImageConfig` struct with all fields
- [ ] Define `BuildResult` struct for return values
- [ ] Add package documentation comments

### Phase 2: Core Build Logic (2 hours)
- [ ] Implement `generateDockerfile()` function
  - [ ] Handle different layer counts
  - [ ] Generate appropriate COPY and RUN commands
  - [ ] Add labels for test identification
  - [ ] Return valid Dockerfile content
- [ ] Implement `generateTestFiles()` function
  - [ ] Calculate size per layer
  - [ ] Create appropriately sized files
  - [ ] Use random data from /dev/urandom
  - [ ] Verify file creation
- [ ] Implement `BuildTestImage()` main function
  - [ ] Create temporary build directory
  - [ ] Generate Dockerfile and test files
  - [ ] Create tar archive for build context
  - [ ] Call Docker build API
  - [ ] Handle errors appropriately
  - [ ] Clean up temporary files

### Phase 3: Docker Integration (1.5 hours)
- [ ] Implement `executeBuild()` function
  - [ ] Configure Docker build options
  - [ ] Call ImageBuild() API
  - [ ] Stream output to logs
  - [ ] Return image ID
- [ ] Implement `verifyImageBuilt()` function
  - [ ] Call ImageInspect() API
  - [ ] Extract layer count and size
  - [ ] Create BuildResult
- [ ] Implement `streamBuildOutput()` function
  - [ ] Parse JSON output from Docker
  - [ ] Log progress messages
  - [ ] Detect build errors

### Phase 4: Fixture Dockerfiles (30 minutes)
- [ ] Create `test/fixtures/` directory
- [ ] Create `Dockerfile.multilayer` with 3 layers
  - [ ] Add Alpine base
  - [ ] Add multiple RUN and COPY commands
  - [ ] Add test labels
- [ ] Create `Dockerfile.minimal` with scratch base
  - [ ] Minimal content for edge cases

### Phase 5: Testing (2 hours)
- [ ] Write unit tests for image builder
  - [ ] Test generateDockerfile() with various configs
  - [ ] Test generateTestFiles() with various sizes
  - [ ] Test BuildTestImage() end-to-end
- [ ] Test image builds successfully
  - [ ] Small image (5MB, 2 layers)
  - [ ] Medium image (100MB, 5 layers)
  - [ ] Large image (500MB, 10 layers)
- [ ] Verify layer counts match specifications
- [ ] Verify image sizes are approximately correct
- [ ] Test error handling for invalid configs

### Phase 6: Documentation (30 minutes)
- [ ] Add godoc comments to all exported functions
- [ ] Document ImageConfig fields
- [ ] Document BuildResult fields
- [ ] Add usage examples in comments
- [ ] Document error conditions

---

## Test Requirements

### Unit Tests

**File**: `test/harness/image_builder_test.go`

```go
func TestGenerateDockerfile(t *testing.T) {
    // Test cases:
    // - Single layer
    // - Multiple layers (2, 5, 10)
    // - Different architectures
    // - Verify Dockerfile format is valid
}

func TestGenerateTestFiles(t *testing.T) {
    // Test cases:
    // - Small files (1MB each)
    // - Large files (100MB each)
    // - Multiple files
    // - Verify sizes match configuration
}

func TestBuildTestImage(t *testing.T) {
    // Test cases:
    // - Build small image (5MB, 2 layers)
    // - Build medium image (100MB, 5 layers)
    // - Build large image (500MB, 10 layers)
    // - Verify image exists in Docker
    // - Verify layer count matches
    // - Verify size approximately matches
}

func TestBuildImageInvalidConfig(t *testing.T) {
    // Test error handling:
    // - Empty name
    // - Negative layer count
    // - Invalid architecture
    // - Negative size
}
```

### Integration Test Usage

The image builder will be used by integration tests like this:

```go
func TestPushSmallImage(t *testing.T) {
    env, err := harness.SetupGiteaRegistry(ctx)
    require.NoError(t, err)
    defer env.CleanupTestEnvironment()

    // Build test image
    config := harness.ImageConfig{
        Name: "testapp",
        Tag: "v1",
        Layers: 2,
        SizeMB: 5,
        Arch: "amd64",
    }
    result, err := env.BuildTestImage(ctx, config)
    require.NoError(t, err)

    // Use image in test
    // ... push to registry and verify ...
}
```

---

## Size Management

### Estimated Line Breakdown

| Component | Lines | Description |
|-----------|-------|-------------|
| ImageConfig & BuildResult structs | 30 | Data structures |
| BuildTestImage() main function | 60 | Core build logic |
| generateDockerfile() | 40 | Dockerfile generation |
| generateTestFiles() | 30 | Test file creation |
| executeBuild() | 30 | Docker build execution |
| verifyImageBuilt() | 20 | Image verification |
| streamBuildOutput() | 20 | Build output handling |
| Helper functions | 20 | Utilities |
| Dockerfile.multilayer | 25 | Fixture template |
| Dockerfile.minimal | 10 | Fixture template |
| Documentation comments | 15 | godoc |
| **Total** | **300** | **Within limit** |

**Status**: ✅ Estimated 300 lines - well within 800 line limit (66% buffer)

### Measurement Protocol

**Tool**: `${PROJECT_ROOT}/tools/line-counter.sh`

**Measurement Points**:
1. After completing image_builder.go core (measure)
2. After completing Docker integration (measure)
3. After completing fixtures (measure)
4. Before requesting review (final measure)

**Thresholds**:
- **Warning**: 700 lines (need to evaluate)
- **Stop**: 800 lines (must create split plan)

---

## Pattern Compliance

### Go Best Practices
- [ ] Follow effective Go guidelines
- [ ] Use appropriate error handling (wrap errors with context)
- [ ] Use context.Context for cancellation
- [ ] Defer cleanup operations
- [ ] Close all readers/writers

### Docker API Usage
- [ ] Use testcontainers-go for container management (from 3.1.1)
- [ ] Use docker/docker SDK for image operations
- [ ] Handle Docker API errors appropriately
- [ ] Stream large outputs (don't buffer in memory)
- [ ] Clean up build artifacts

### Testing Patterns
- [ ] Table-driven tests where appropriate
- [ ] Use testify/require for assertions
- [ ] Test error paths, not just happy paths
- [ ] Mock external dependencies where possible
- [ ] Use build tags for integration tests

---

## Integration Strategy

### Dependencies from Effort 3.1.1

This effort requires the test harness from 3.1.1:

```go
// TestEnvironment from effort 3.1.1
type TestEnvironment struct {
    GiteaContainer  testcontainers.Container
    RegistryURL     string
    AdminUsername   string
    AdminPassword   string
    DockerClient    docker.Client
    Cleanup         func() error
}

// We extend this with BuildTestImage method
func (env *TestEnvironment) BuildTestImage(ctx context.Context, config ImageConfig) (*BuildResult, error)
```

### Provides for Efforts 3.1.3 and 3.1.4

Integration tests will use the image builder:

```go
// Example usage in core workflow tests (3.1.3)
result, err := env.BuildTestImage(ctx, harness.ImageConfig{
    Name: "testapp",
    Tag: "latest",
    Layers: 5,
    SizeMB: 100,
    Arch: "amd64",
})

// Push the built image to Gitea
err = pushCommand.Execute(ctx, result.ImageRef, env.RegistryURL)
```

---

## Error Handling

### Expected Errors

1. **Build Directory Creation Failure**
   - Cause: Insufficient permissions or disk space
   - Handling: Return descriptive error
   - User action: Check filesystem permissions

2. **Dockerfile Generation Failure**
   - Cause: Invalid configuration
   - Handling: Validate config before generation
   - User action: Fix configuration

3. **Test File Creation Failure**
   - Cause: Disk space issues
   - Handling: Return error with disk space info
   - User action: Free up disk space

4. **Docker Build Failure**
   - Cause: Docker daemon issues, invalid Dockerfile
   - Handling: Stream error output, return wrapped error
   - User action: Check Docker daemon, review logs

5. **Image Verification Failure**
   - Cause: Image not found after build
   - Handling: Return error with build output
   - User action: Review build logs

### Error Messages

All errors should:
- Be wrapped with context using `fmt.Errorf("context: %w", err)`
- Include relevant details (image name, config)
- Suggest corrective actions
- Not expose sensitive information

---

## Performance Considerations

### Build Time Optimization

- Use Alpine as base (small, fast to download)
- Generate files efficiently (buffered writes)
- Parallelize file creation where possible
- Cache Docker layers when possible

### Resource Management

- Clean up temporary files immediately after use
- Don't keep build contexts in memory
- Stream large files instead of buffering
- Limit concurrent builds to avoid resource exhaustion

### Expected Performance

| Image Size | Layers | Expected Build Time |
|------------|--------|---------------------|
| 5MB | 2 | <5 seconds |
| 100MB | 5 | <30 seconds |
| 500MB | 10 | <2 minutes |

---

## Security Considerations

### Build Security

- [ ] Generate files in secure temporary directories
- [ ] Clean up all temporary files (don't leave artifacts)
- [ ] Use random data for file generation (not sensitive data)
- [ ] Validate all configuration inputs
- [ ] Don't expose Docker socket unnecessarily

### Image Security

- [ ] Use official base images (Alpine)
- [ ] Don't include credentials in images
- [ ] Add labels for test identification
- [ ] Tag images clearly as test images
- [ ] Document that images are for testing only

---

## Success Criteria

### Functional Requirements
- ✅ BuildTestImage() successfully creates images
- ✅ Generated images have correct layer count
- ✅ Generated images have approximately correct size
- ✅ Images can be pushed to registry (verified in 3.1.3)
- ✅ Multiple images can be built without conflicts

### Code Quality Requirements
- ✅ All functions have godoc comments
- ✅ Unit tests cover all public functions
- ✅ Error handling follows Go best practices
- ✅ Code passes go vet and golint
- ✅ Test coverage >85%

### Integration Requirements
- ✅ Integrates with TestEnvironment from 3.1.1
- ✅ Provides clean API for integration tests
- ✅ Images work with push command from Phase 2
- ✅ Cleanup is reliable (no leftover images)

---

## Validation Checklist (R510 Compliance)

### Pre-Implementation
- [x] R213 metadata complete
- [x] Dependencies identified (effort 3.1.1)
- [x] File paths specified
- [x] Line estimate within limits (300 lines)
- [x] Integration strategy defined
- [x] Test requirements documented

### During Implementation (SW Engineer)
- [ ] Create all specified files
- [ ] Implement all required functions
- [ ] Write unit tests for all components
- [ ] Test with Docker daemon
- [ ] Verify images build successfully
- [ ] Measure with line counter regularly

### Post-Implementation (Code Review)
- [ ] All files created and committed
- [ ] All tests passing
- [ ] Line count <800 (verify with tool)
- [ ] Integration with 3.1.1 verified
- [ ] Documentation complete
- [ ] Ready for use by 3.1.3 and 3.1.4

---

## Document Compliance

### Rules Compliance
- ✅ **R213**: Complete metadata block included
- ✅ **R219**: Dependencies documented (requires 3.1.1)
- ✅ **R220/R221**: Within size limits (300 lines estimated)
- ✅ **R383**: Plan in .software-factory/ with timestamp
- ✅ **R502**: Implementation plan with exact specifications
- ✅ **R510/R511**: Validation checklists included

### Template Compliance
- ✅ Used EFFORT-IMPLEMENTATION-PLAN.md template
- ✅ All required sections included
- ✅ Task checklists for SW Engineer
- ✅ Integration points documented
- ✅ Success criteria defined

---

## Next Steps

### For Orchestrator
1. Verify effort 3.1.1 is complete
2. Spawn SW Engineer for this effort
3. Provide this plan as implementation guide
4. Monitor progress every 5 messages
5. Spawn Code Reviewer when implementation complete

### For SW Engineer
1. Read this plan completely
2. Check out branch: `idpbuilder-oci-push/phase3/wave1/effort-3.1.2-image-builders`
3. Verify base branch: `idpbuilder-oci-push/phase3/wave1/effort-3.1.1-test-harness`
4. Follow task checklist sequentially
5. Measure line count after each phase
6. Request review when all tasks complete

### For Code Reviewer
1. Verify all files created
2. Verify integration with 3.1.1
3. Run unit tests
4. Test image building with Docker
5. Measure final line count
6. Approve or request fixes

---

## Document Status

**Status**: ✅ READY FOR IMPLEMENTATION
**Effort**: 3.1.2 - Test Image Builders
**Phase/Wave**: Phase 3, Wave 1
**Created**: 2025-11-04 14:08:30 UTC
**Created By**: code-reviewer (EFFORT_PLAN_CREATION state)

**Next Action**: Orchestrator spawns SW Engineer for implementation (after 3.1.1 complete)

---

**END OF IMPLEMENTATION PLAN**
