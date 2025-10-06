# E1.2.3: Image Push Operations - Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `phase1/wave2/image-push-operations`
**Can Parallelize**: No (Dependencies: E1.2.1, E1.2.2)
**Parallel With**: None
**Size Estimate**: 700 lines
**Dependencies**: E1.2.1 (Command structure), E1.2.2 (Registry authentication)
**Base Branch**: phase1-integration (per R501/R509 cascade requirements)

## Overview
- **Effort**: Implement OCI image push operations for idpbuilder push command
- **Phase**: 1, Wave: 2
- **Estimated Size**: 700 lines
- **Implementation Time**: 6-8 hours
- **Core Library**: go-containerregistry v0.20.6 (LOCKED per R381)

## Dependency Context (R219)
This effort builds upon:
1. **E1.2.1 Command Structure**: Provides the `push` command skeleton with Cobra framework, basic validation, and flag structure
2. **E1.2.2 Registry Authentication**: Provides credential handling, authentication flow, and --insecure flag implementation

### What We Import From Dependencies:
- Command initialization and flag parsing (E1.2.1)
- Authentication configuration and credential resolver (E1.2.2)
- Registry transport setup with insecure option handling (E1.2.2)
- Error handling patterns established in previous efforts

## File Structure

### Core Implementation Files
- `pkg/push/operations.go` (250 lines): Main push operations and orchestration
  - `PushImages()`: Main entry point coordinating the push workflow
  - `discoverImages()`: Image discovery logic from build context
  - `validateImages()`: Pre-push validation and compatibility checks
  - Error aggregation and reporting

- `pkg/push/discovery.go` (150 lines): OCI image discovery logic
  - `DiscoverLocalImages()`: Scan directories for OCI artifacts
  - `LoadImageManifests()`: Parse and validate image manifests
  - `FilterPushTargets()`: Determine which images to push
  - Support for multiple image formats (tarball, directory layout)

- `pkg/push/pusher.go` (200 lines): Core push functionality using go-containerregistry
  - `ImagePusher` struct: Encapsulates push operations
  - `Push()`: Single image push with retry logic
  - `PushWithProgress()`: Push with progress indicators
  - Integration with go-containerregistry remote.Write()
  - Digest verification and manifest handling

- `pkg/push/progress.go` (50 lines): Lightweight progress indicators
  - `ProgressReporter` interface
  - `ConsoleProgressReporter`: Terminal-based progress display
  - Byte-level progress tracking
  - Layer upload status reporting

- `pkg/push/logging.go` (50 lines): Comprehensive logging utilities
  - Debug/Info/Warn/Error level logging
  - Structured logging with context
  - Registry interaction logging at debug level
  - Performance metrics logging

### Test Files
- `pkg/push/operations_test.go` (100 lines): Operations orchestration tests
- `pkg/push/discovery_test.go` (80 lines): Image discovery tests
- `pkg/push/pusher_test.go` (120 lines): Push functionality tests
- `pkg/push/testdata/`: Test fixtures and mock images

## Implementation Steps

### Step 1: Create Core Structures (50 lines)
```go
// pkg/push/operations.go
type PushOperation struct {
    Registry   string
    Repository string
    Auth       authn.Authenticator  // From E1.2.2
    Transport  *http.Transport      // From E1.2.2
    Progress   ProgressReporter
    Logger     *logr.Logger
}

type PushResult struct {
    ImageName string
    Digest    string
    Size      int64
    Duration  time.Duration
    Error     error
}
```

### Step 2: Implement Image Discovery (150 lines)
```go
// pkg/push/discovery.go
func DiscoverLocalImages(buildPath string) ([]v1.Image, error) {
    // 1. Check for .tar files (docker save format)
    // 2. Check for OCI layout directories
    // 3. Validate each discovered image
    // 4. Return list of pushable images
}

func LoadImageManifests(imagePath string) (*v1.Manifest, error) {
    // Use go-containerregistry tarball or layout packages
}
```

### Step 3: Implement Push Logic (200 lines)
```go
// pkg/push/pusher.go
func (p *ImagePusher) Push(ctx context.Context, img v1.Image, ref name.Reference) error {
    // 1. Set up remote options with auth and transport
    opts := []remote.Option{
        remote.WithAuth(p.Auth),
        remote.WithTransport(p.Transport),
    }

    // 2. Push with progress tracking
    return remote.Write(ref, img, opts...)
}
```

### Step 4: Add Progress Indicators (50 lines)
```go
// pkg/push/progress.go
type ProgressReporter interface {
    StartLayer(digest string, size int64)
    UpdateLayer(digest string, written int64)
    FinishLayer(digest string)
    Complete(digest string)
}

// Simple console implementation with percentage display
```

### Step 5: Implement Retry Logic (100 lines)
```go
// pkg/push/pusher.go
func (p *ImagePusher) PushWithRetry(ctx context.Context, img v1.Image, ref name.Reference) error {
    backoff := time.Second
    maxRetries := 5

    for attempt := 0; attempt < maxRetries; attempt++ {
        err := p.Push(ctx, img, ref)
        if err == nil {
            return nil
        }

        // Check if error is retryable
        if !isRetryable(err) {
            return err
        }

        // Exponential backoff
        time.Sleep(backoff)
        backoff *= 2
    }
    return fmt.Errorf("push failed after %d attempts", maxRetries)
}
```

### Step 6: Add Comprehensive Logging (50 lines)
```go
// pkg/push/logging.go
func LogPushStart(logger *logr.Logger, image, registry string) {
    logger.Info("Starting image push", "image", image, "registry", registry)
}

func LogLayerProgress(logger *logr.Logger, layer string, percent int) {
    logger.V(1).Info("Layer upload progress", "layer", layer, "percent", percent)
}

// Debug-level logging for registry interactions
func LogRegistryRequest(logger *logr.Logger, method, url string) {
    logger.V(2).Info("Registry request", "method", method, "url", url)
}
```

### Step 7: Wire Everything Together (100 lines)
```go
// pkg/push/operations.go
func PushImages(cmd *cobra.Command, args []string) error {
    // 1. Get configuration from command (via E1.2.1)
    config := getConfigFromCommand(cmd)

    // 2. Set up authentication (via E1.2.2)
    auth := setupAuthentication(config)

    // 3. Discover images to push
    images, err := DiscoverLocalImages(config.BuildPath)
    if err != nil {
        return fmt.Errorf("failed to discover images: %w", err)
    }

    // 4. Create pusher with progress
    pusher := NewImagePusher(auth, config.Transport, config.Logger)

    // 5. Push each image
    for _, img := range images {
        ref, _ := name.ParseReference(fmt.Sprintf("%s/%s", config.Registry, img.Name))
        if err := pusher.PushWithRetry(ctx, img, ref); err != nil {
            return fmt.Errorf("failed to push %s: %w", img.Name, err)
        }
    }

    return nil
}
```

## Library Version Requirements (R381)

### Locked Dependencies (DO NOT UPDATE)
All versions from existing go.mod are IMMUTABLE:
- `github.com/google/go-containerregistry v0.20.6` (LOCKED - primary for this effort)
- `github.com/spf13/cobra v1.9.1` (LOCKED - from E1.2.1)
- `github.com/go-logr/logr v1.4.3` (LOCKED - for logging)
- `github.com/stretchr/testify v1.10.0` (LOCKED - for testing)

### go-containerregistry Usage Patterns
Key APIs we'll use from v0.20.6:
- `remote.Write()`: Main push function
- `remote.WithAuth()`: Authentication setup
- `remote.WithTransport()`: Custom transport (for --insecure)
- `tarball.ImageFromPath()`: Load images from tar
- `layout.ImageIndexFromPath()`: Load OCI layout
- `name.ParseReference()`: Parse registry references

## Size Management
- **Estimated Lines**: 700 (within 800 limit)
- **Measurement Tool**: `/home/vscode/workspaces/idpbuilder-push-oci/tools/line-counter.sh`
- **Check Frequency**: After each major component
- **Split Threshold**: 700 lines (warning), 800 lines (stop)
- **Current Breakdown**:
  - Core implementation: 500 lines
  - Tests: 200 lines
  - Total: 700 lines

## Test Requirements

### Unit Tests (200 lines)
- **Coverage Target**: 85%
- **Key Test Cases**:
  - Image discovery with various formats
  - Push success scenarios
  - Authentication integration
  - Retry logic with backoff
  - Progress reporting accuracy
  - Error handling and recovery

### Integration Tests (via E1.1.3 framework)
- Push to real Gitea registry
- Verify image integrity after push
- Test --insecure flag with self-signed certs
- Concurrent push operations
- Large image handling (>100MB)

### Test Files Structure:
```
pkg/push/
├── operations_test.go      # Main workflow tests
├── discovery_test.go       # Discovery logic tests
├── pusher_test.go          # Push functionality tests
├── testdata/
│   ├── valid-image.tar    # Valid OCI image
│   ├── invalid-image.tar  # Corrupted image
│   └── oci-layout/        # OCI directory layout
└── mocks/
    └── registry.go         # Mock registry for unit tests
```

## Pattern Compliance

### idpbuilder Patterns (from E1.1.1 analysis)
- Command structure follows existing patterns
- Error handling matches project style
- Logging uses standard idpbuilder logger
- Configuration via viper/cobra flags

### Go Best Practices
- Proper error wrapping with `fmt.Errorf("%w")`
- Context usage for cancellation
- Interface-based design for extensibility
- Comprehensive test coverage

## Integration Points

### With E1.2.1 (Command Structure)
- Import command setup and flag definitions
- Use established configuration patterns
- Follow command execution flow

### With E1.2.2 (Registry Authentication)
- Import authentication setup functions
- Use credential resolution logic
- Leverage transport configuration with --insecure support

## Success Criteria

1. **Functionality**:
   - ✅ Discovers OCI images from build context
   - ✅ Successfully pushes to any OCI registry
   - ✅ Progress indicators show upload status
   - ✅ Comprehensive logging at appropriate levels

2. **Quality**:
   - ✅ All unit tests passing (>85% coverage)
   - ✅ Integration with E1.2.1 and E1.2.2 seamless
   - ✅ Retry logic handles transient failures
   - ✅ Clean error messages for users

3. **Performance**:
   - ✅ Efficient layer deduplication
   - ✅ Minimal memory usage during push
   - ✅ Concurrent layer uploads where possible

## Risk Mitigation

1. **Risk**: go-containerregistry API complexity
   - **Mitigation**: Start with simple remote.Write(), add features incrementally
   - **Fallback**: Use examples from go-containerregistry documentation

2. **Risk**: Image discovery ambiguity
   - **Mitigation**: Clear conventions for image location
   - **Fallback**: Allow explicit image path specification

3. **Risk**: Progress indicator performance impact
   - **Mitigation**: Keep indicators lightweight, optional verbose mode
   - **Fallback**: Simple percentage display only

## Implementation Notes

- Use go-containerregistry's high-level APIs where possible
- Log all registry interactions at debug level
- Ensure proper cleanup of temporary files
- Add TODO markers for future enhancements (rate limiting, parallel uploads)
- Follow TDD approach: write tests first for each component

## Verification Commands

```bash
# After implementation, verify with:
cd /home/vscode/workspaces/idpbuilder-push-oci/efforts/phase1/wave2/E1.2.3-image-push-operations

# Check size compliance
$PROJECT_ROOT/tools/line-counter.sh

# Run unit tests
go test ./pkg/push/... -v -cover

# Run integration tests (requires idpbuilder cluster)
go test ./test/integration/... -tags=integration

# Verify against real Gitea
idpbuilder push --registry gitea.localhost:8443 \
  --username gitea_admin --password gitea_password \
  --insecure ./test-image.tar
```

## CONTINUE-SOFTWARE-FACTORY=TRUE