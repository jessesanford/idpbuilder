# Image Operations Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (R054 REQUIREMENT)

**Effort ID**: E2.2.2-B
**Effort Name**: image-operations
**Phase**: 2, **Wave**: 2
**Branch**: `idpbuilder-oci-build-push/phase2/wave2/image-operations`
**Base Branch**: `idpbuilder-oci-build-push/phase2/wave2/credential-management`
**Can Parallelize**: No
**Parallel With**: None
**Size Estimate**: 500 lines
**Dependencies**: ["cli-commands", "credential-management"]
**Creation Time**: 2025-09-15 21:29:57 UTC

## 🎯 EFFORT INFRASTRUCTURE METADATA (R303 REQUIREMENT)

**WORKING_DIRECTORY**: /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave2/image-operations
**REMOTE**: https://github.com/jessesanford/idpbuilder.git
**INTEGRATION_BRANCH**: phase2/wave2/image-operations
**PROJECT_ROOT**: /home/vscode/workspaces/idpbuilder-oci-build-push

## 📋 OVERVIEW

This effort removes ALL remaining placeholder code and makes the idpbuilder OCI build/push functionality production-ready. After completion, the binary will be 100% functional with NO stubs, placeholders, or feature flags.

### Scope
- Real OCI image loading from Docker daemon
- Actual manifest generation with proper SHA256 digests
- Production-grade progress tracking
- Complete removal of all feature flags
- Real error handling throughout

### Key Deliverables
1. **Image Loader**: Real Docker daemon integration for image loading
2. **Manifest Generation**: Proper OCI manifest with computed digests
3. **Progress Tracking**: Real-time progress reporting
4. **Production Ready**: All placeholders removed, all features enabled

## 🏗️ TECHNICAL ARCHITECTURE

### Component Structure
```
pkg/gitea/
├── image_loader.go         # NEW: Docker daemon integration
├── image_loader_test.go    # NEW: Tests for image loader
├── progress.go             # NEW: Real progress tracking
├── progress_test.go        # NEW: Progress tests
├── client.go              # MODIFIED: Replace placeholder manifest
└── client_test.go         # MODIFIED: Update tests

pkg/cmd/
├── build.go               # MODIFIED: Remove feature flags
└── push.go                # MODIFIED: Update for real operations

pkg/build/
└── feature_flags.go       # DELETE: No longer needed
```

### Key Dependencies
- `github.com/docker/docker/client`: Docker daemon API
- `github.com/docker/docker/api/types`: Docker types
- `github.com/opencontainers/go-digest`: OCI digest calculation
- `github.com/opencontainers/image-spec`: OCI image specifications
- Existing packages from Phase 2 Wave 1 and Wave 2A

### Integration Points
1. **Docker Daemon**: Connect to local Docker daemon for image operations
2. **Registry Client**: Use existing registry.Registry interface for push
3. **Credential Manager**: Leverage credential-management effort's CredentialManager
4. **Certificate Manager**: Use Phase 1 certificate infrastructure

## 📁 FILE STRUCTURE

### New Files to Create

#### 1. `pkg/gitea/image_loader.go` (200 lines)
```go
package gitea

import (
    "context"
    "io"
    "github.com/docker/docker/client"
    "github.com/docker/docker/api/types"
    "github.com/opencontainers/go-digest"
)

// ImageLoader handles loading images from Docker daemon
type ImageLoader struct {
    client *client.Client
}

// NewImageLoader creates a new Docker image loader
func NewImageLoader() (*ImageLoader, error)

// LoadImage loads an image from the Docker daemon
func (il *ImageLoader) LoadImage(ctx context.Context, imageRef string) (*ImageManifest, error)

// GetImageManifest retrieves the OCI manifest for an image
func (il *ImageLoader) GetImageManifest(ctx context.Context, imageID string) (*ImageManifest, error)

// GetImageContent returns a reader for the image content
func (il *ImageLoader) GetImageContent(ctx context.Context, imageID string) (io.ReadCloser, error)

// CalculateDigest computes the SHA256 digest of image content
func (il *ImageLoader) CalculateDigest(content []byte) digest.Digest

// ImageManifest represents an OCI image manifest
type ImageManifest struct {
    SchemaVersion int
    MediaType     string
    Config        ManifestConfig
    Layers        []ManifestLayer
}
```

#### 2. `pkg/gitea/image_loader_test.go` (80 lines)
```go
package gitea

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewImageLoader(t *testing.T)
func TestLoadImage(t *testing.T)
func TestGetImageManifest(t *testing.T)
func TestCalculateDigest(t *testing.T)
```

#### 3. `pkg/gitea/progress.go` (50 lines)
```go
package gitea

import (
    "sync"
    "time"
)

// ProgressTracker tracks real image push progress
type ProgressTracker struct {
    mu            sync.RWMutex
    totalBytes    int64
    uploadedBytes int64
    layers        []LayerProgress
    startTime     time.Time
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(totalSize int64) *ProgressTracker

// UpdateProgress updates the progress for a specific layer
func (pt *ProgressTracker) UpdateProgress(layerID string, bytes int64)

// GetProgress returns current progress information
func (pt *ProgressTracker) GetProgress() PushProgress

// LayerProgress tracks progress for individual layers
type LayerProgress struct {
    ID           string
    Size         int64
    Uploaded     int64
    Status       string
}
```

#### 4. `pkg/gitea/progress_test.go` (30 lines)
```go
package gitea

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestProgressTracker(t *testing.T)
func TestGetProgress(t *testing.T)
```

### Files to Modify

#### 1. `pkg/gitea/client.go` (~100 lines modified)
**Remove (lines 121-146):**
- Entire `placeholderManifest` block
- Simulated progress loop (lines 101-113)
- Comment about placeholder implementation

**Add:**
```go
// Push pushes an image to the registry with real progress reporting
func (c *Client) Push(imageRef string, progressChan chan<- PushProgress) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    // Initialize image loader
    loader, err := NewImageLoader()
    if err != nil {
        return fmt.Errorf("failed to create image loader: %w", err)
    }

    // Load image from Docker daemon
    manifest, err := loader.LoadImage(ctx, imageRef)
    if err != nil {
        return fmt.Errorf("failed to load image: %w", err)
    }

    // Get image content reader
    contentReader, err := loader.GetImageContent(ctx, imageRef)
    if err != nil {
        return fmt.Errorf("failed to get image content: %w", err)
    }
    defer contentReader.Close()

    // Initialize real progress tracker
    tracker := NewProgressTracker(manifest.TotalSize)

    // Start real progress reporting
    if progressChan != nil {
        go c.reportProgress(tracker, progressChan)
    }

    // Perform actual push with progress tracking
    return c.registry.Push(ctx, imageRef, contentReader)
}

// reportProgress reports real progress updates
func (c *Client) reportProgress(tracker *ProgressTracker, progressChan chan<- PushProgress) {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    for range ticker.C {
        progress := tracker.GetProgress()
        select {
        case progressChan <- progress:
            if progress.Percentage >= 100 {
                return
            }
        default:
            // Channel full, skip this update
        }
    }
}
```

#### 2. `pkg/cmd/build.go` (~20 lines modified)
**Remove:**
- Any references to `IsImageBuilderEnabled()`
- Feature flag checks

**Modify:**
```go
func runBuild(cmd *cobra.Command, args []string) error {
    // Remove feature flag check - always enabled
    // Existing validation code remains

    // Direct builder usage without feature flag
    builder, err := build.NewBuilder(storageDir)
    if err != nil {
        return fmt.Errorf("failed to create builder: %w", err)
    }

    // Rest of the function remains the same
}
```

#### 3. `pkg/cmd/push.go` (~20 lines modified)
**Update:**
- Progress reporting to use real tracker
- Remove any "simulated" or "placeholder" references

### Files to Delete

#### 1. `pkg/build/feature_flags.go`
- Completely remove this file
- No longer needed as all features are production-ready

## 🔄 IMPLEMENTATION SEQUENCE

### Step 1: Docker Integration Foundation (150 lines)
1. Create `pkg/gitea/image_loader.go` structure
2. Implement Docker client initialization
3. Add image loading from daemon
4. Implement manifest generation with real digests
5. Add error handling

### Step 2: Progress Tracking (80 lines)
1. Create `pkg/gitea/progress.go`
2. Implement real progress tracking
3. Add layer-by-layer progress
4. Create progress reporting goroutine

### Step 3: Client Integration (100 lines)
1. Modify `pkg/gitea/client.go`
2. Remove placeholder manifest (lines 121-146)
3. Remove simulated progress (lines 101-113)
4. Integrate ImageLoader
5. Add real progress reporting

### Step 4: Command Updates (40 lines)
1. Update `pkg/cmd/build.go`
2. Remove feature flag checks
3. Update `pkg/cmd/push.go`
4. Ensure production-ready behavior

### Step 5: Cleanup (30 lines)
1. Delete `pkg/build/feature_flags.go`
2. Remove all TODO comments
3. Remove all placeholder strings
4. Update documentation strings

### Step 6: Testing (100 lines)
1. Create `pkg/gitea/image_loader_test.go`
2. Create `pkg/gitea/progress_test.go`
3. Update existing tests
4. Add integration tests

## 📏 SIZE MANAGEMENT STRATEGY

### Line Count Tracking
- **Target**: 500 lines total
- **Current Estimate**: 500 lines
- **Buffer**: 300 lines before warning

### Measurement Points
1. After Docker integration (150 lines) - MEASURE
2. After progress tracking (230 lines) - MEASURE
3. After client integration (330 lines) - MEASURE
4. After command updates (370 lines) - MEASURE
5. After cleanup (400 lines) - MEASURE
6. After testing (500 lines) - FINAL MEASURE

### Measurement Command
```bash
# Find project root and use line counter
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
$PROJECT_ROOT/tools/line-counter.sh
```

## 🧪 TESTING REQUIREMENTS

### Unit Test Coverage
- **Target**: 85% coverage
- **Required Tests**:
  - ImageLoader initialization
  - Docker client connection
  - Manifest generation
  - Digest calculation
  - Progress tracking
  - Error handling

### Integration Tests
- Push to local registry
- Progress reporting accuracy
- Credential management integration
- Certificate handling

### Manual Testing Checklist
1. [ ] Build image from context directory
2. [ ] Push image to Gitea registry
3. [ ] Verify progress reporting
4. [ ] Test with invalid credentials
5. [ ] Test with insecure mode
6. [ ] Verify manifest correctness

## ✅ VALIDATION CHECKPOINTS

### Checkpoint 1: Docker Integration (After Step 1)
- [ ] Docker client connects successfully
- [ ] Images load from daemon
- [ ] Manifests generate with real digests
- [ ] No placeholder strings remain

### Checkpoint 2: Progress Implementation (After Step 2)
- [ ] Real progress tracking works
- [ ] Layer-by-layer updates
- [ ] Accurate percentage calculation

### Checkpoint 3: Full Integration (After Step 3)
- [ ] Push operation works end-to-end
- [ ] No simulated behavior
- [ ] Proper error handling

### Checkpoint 4: Production Ready (After Step 5)
- [ ] All feature flags removed
- [ ] No TODO comments
- [ ] No placeholder code
- [ ] All tests pass

### Final Validation
```bash
# Build the binary
make build

# Run all tests
make test

# Check for any remaining placeholders
grep -r "placeholder\|TODO\|stub\|simulated\|feature.*flag" pkg/

# Verify line count
$PROJECT_ROOT/tools/line-counter.sh
```

## ⚠️ RISK ASSESSMENT

### Technical Risks
1. **Docker Daemon Dependency**
   - Risk: Docker daemon not available
   - Mitigation: Clear error messages, fallback to buildah if available

2. **Manifest Complexity**
   - Risk: Incorrect manifest format
   - Mitigation: Follow OCI spec strictly, validate against schema

3. **Progress Accuracy**
   - Risk: Inaccurate progress reporting
   - Mitigation: Track actual bytes transferred

### Integration Risks
1. **Registry Compatibility**
   - Risk: Registry rejects manifest
   - Mitigation: Test with Gitea registry, follow API exactly

2. **Credential Flow**
   - Risk: Credentials not properly passed
   - Mitigation: Use existing CredentialManager

## 🔗 INTEGRATION STRATEGY

### Dependencies Integration
1. **From cli-commands (E2.2.1)**:
   - Command structure
   - Helper functions
   - Error formatting

2. **From credential-management (E2.2.2-A)**:
   - CredentialManager
   - Credential providers
   - Authentication flow

3. **From Phase 1**:
   - Certificate management
   - Trust store
   - Registry interface

### Testing Integration
1. Start Docker daemon for tests
2. Create test images
3. Use test registry
4. Verify end-to-end flow

### Merge Strategy
1. Complete implementation
2. Run full test suite
3. Verify no placeholders remain
4. Measure final size
5. Create PR to base branch

## 📊 SUCCESS CRITERIA

### Must Have
- [x] All placeholder code removed
- [x] Real Docker daemon integration
- [x] Proper OCI manifest generation
- [x] Real progress tracking
- [x] All feature flags removed
- [x] Under 800 lines (target: 500)

### Should Have
- [x] 85% test coverage
- [x] Clear error messages
- [x] Efficient image loading
- [x] Accurate progress reporting

### Nice to Have
- [ ] Buildah fallback support
- [ ] Multi-platform manifest support
- [ ] Compression options
- [ ] Layer caching

## 🎯 DELIVERABLES

1. **Working Binary**: idpbuilder with full OCI build/push capabilities
2. **No Placeholders**: Zero stub implementations
3. **Production Ready**: All features enabled by default
4. **Tested**: >85% coverage with passing tests
5. **Size Compliant**: Under 800 lines (target: 500)

## 📝 NOTES

- This effort completes the Phase 2 Wave 2 implementation
- After this, the system is fully functional
- No feature flags or placeholders remain
- Ready for production use

---
*Plan created: 2025-09-15 21:29:57 UTC*
*Effort: E2.2.2-B (image-operations)*
*Phase 2, Wave 2 - Final Production Implementation*