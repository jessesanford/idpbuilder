# go-containerregistry Image Builder Implementation Plan

**Effort ID**: E2.1.1  
**Effort Name**: go-containerregistry-image-builder  
**Phase**: 2 - Build & Push Implementation  
**Wave**: 1 - Core OCI Operations  
**Created By**: Code Reviewer Agent  
**Date Created**: 2025-09-02 22:41:46  
**Assigned To**: SW Engineer  

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder`  
**Can Parallelize**: Yes  
**Parallel With**: E2.1.2 (gitea-registry-client)  
**Size Estimate**: 600 lines  
**Dependencies**: Phase 1 certificate infrastructure (E1.1.1, E1.1.2, E1.2.1, E1.2.2)  

<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## 🚨 EFFORT INFRASTRUCTURE METADATA
**WORKING_DIRECTORY**: `/home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase2/wave1/go-containerregistry-image-builder`  
**BRANCH**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder`  
**BASE_BRANCH**: `idpbuilder-oci-go-cr/phase1-integration-20250902-194557`  
**REMOTE**: `origin`  
**EFFORT_NAME**: `go-containerregistry-image-builder`  

## 📋 Effort Overview

### Description
This effort implements the core OCI image building functionality using the go-containerregistry library. It creates a daemonless image builder that can assemble OCI images from local directories, supporting both direct image objects and OCI tarball output. The implementation integrates with Phase 1's certificate infrastructure for secure registry operations.

### Size Estimate
- **Estimated Lines**: 600 (well within limit)
- **Confidence Level**: High
- **Split Risk**: Low

### Dependencies
- **Phase 1 Imports**:
  - E1.1.1 (kind-certificate-extraction): Certificate extraction for registry auth
  - E1.1.2 (registry-tls-trust-integration): TLS trust store management
  - E1.2.1 (certificate-validation-pipeline): Certificate validation
  - E1.2.2 (fallback-strategies): Insecure mode fallback
- **External**:
  - github.com/google/go-containerregistry v0.19.0
  - Standard library (archive/tar, io, path/filepath)
- **Parallel With**: E2.1.2 (gitea-registry-client) - No interdependency

## ⚡ Requirements

### Functional Requirements
- [ ] Create OCI images from directory contents
- [ ] Support single-layer tar archives
- [ ] Generate valid OCI image configurations
- [ ] Support linux/amd64 and linux/arm64 platforms
- [ ] Calculate correct digests and sizes
- [ ] Handle file permissions and timestamps
- [ ] Support base image specification
- [ ] Create OCI tarballs for offline distribution

### Non-Functional Requirements
- [ ] Performance: Build 100MB image in <30 seconds
- [ ] Memory: Handle 1GB directories without OOM
- [ ] Compatibility: OCI spec 1.0 compliant
- [ ] Extensibility: Interface-based design for future enhancements

### Acceptance Criteria
- [ ] All unit tests passing
- [ ] Test coverage ≥ 80%
- [ ] Code review approved
- [ ] Size ≤ 800 lines (measured with line-counter.sh)
- [ ] No critical TODOs
- [ ] Documentation complete
- [ ] Feature flag for incomplete functionality (R307)

## 🏗️ Technical Architecture

### Core Interfaces
```go
// pkg/builder/builder.go
type Builder interface {
    // Build creates an OCI image from a context directory
    Build(ctx context.Context, contextDir string, opts BuildOptions) (v1.Image, error)
    
    // BuildTarball creates an OCI tarball from a context directory
    BuildTarball(ctx context.Context, contextDir string, output string, opts BuildOptions) error
}

// pkg/builder/options.go
type BuildOptions struct {
    Platform    v1.Platform         // Target platform
    BaseImage   string              // Optional base image reference
    Labels      map[string]string   // OCI labels
    Env         []string           // Environment variables
    WorkingDir  string             // Working directory in container
    Entrypoint  []string          // Container entrypoint
    Cmd         []string          // Container command
    FeatureFlags map[string]bool   // Feature flags for R307 compliance
}
```

### Component Architecture
1. **SimpleBuilder**: Main builder implementation
2. **LayerFactory**: Creates layers from directories
3. **ConfigFactory**: Generates OCI configurations
4. **TarballWriter**: Exports images as OCI tarballs

## 📁 Implementation Details

### Files to Create
| File Path | Purpose | Estimated Lines |
|-----------|---------|-----------------|
| `pkg/builder/builder.go` | Main Builder interface and SimpleBuilder implementation | 180 |
| `pkg/builder/options.go` | BuildOptions and platform helpers | 60 |
| `pkg/builder/layer.go` | Layer creation from directories | 120 |
| `pkg/builder/config.go` | OCI config generation | 80 |
| `pkg/builder/tarball.go` | OCI tarball export functionality | 60 |
| `pkg/builder/builder_test.go` | Comprehensive test suite | 80 |
| `pkg/builder/testdata/Dockerfile` | Test fixture | 10 |
| `pkg/builder/testdata/content/app.txt` | Test content | 5 |
| `pkg/builder/testdata/content/config.yaml` | Test config | 5 |
| **Total** | | 600 |

### Files to Modify
None - this is a new component in the isolated effort directory.

## 🔧 Detailed Implementation Components

### Component 1: SimpleBuilder Implementation
**File**: `pkg/builder/builder.go` (180 lines)

```go
package builder

import (
    "context"
    "fmt"
    "path/filepath"
    
    "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/empty"
    "github.com/google/go-containerregistry/pkg/v1/mutate"
    "github.com/google/go-containerregistry/pkg/v1/tarball"
)

// SimpleBuilder implements Builder interface
type SimpleBuilder struct {
    platform     v1.Platform
    baseImage    v1.Image
    featureFlags map[string]bool
}

// NewBuilder creates a new builder instance
func NewBuilder(opts BuildOptions) (*SimpleBuilder, error) {
    // Initialize with platform defaults
    // Load base image if specified
    // Set feature flags for R307
    return &SimpleBuilder{
        platform:     opts.Platform,
        featureFlags: opts.FeatureFlags,
    }, nil
}

// Build creates OCI image from directory
func (b *SimpleBuilder) Build(ctx context.Context, contextDir string, opts BuildOptions) (v1.Image, error) {
    // Validate context directory
    // Create layer from directory
    // Generate OCI config
    // Combine into image
    // Return image object
}

// BuildTarball exports image as OCI tarball
func (b *SimpleBuilder) BuildTarball(ctx context.Context, contextDir string, output string, opts BuildOptions) error {
    // Build image first
    // Export to tarball format
    // Write to output path
}
```

### Component 2: Layer Creation
**File**: `pkg/builder/layer.go` (120 lines)

```go
package builder

import (
    "archive/tar"
    "io"
    "os"
    "path/filepath"
    
    "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/tarball"
)

// LayerFactory creates OCI layers from directories
type LayerFactory struct {
    preservePermissions bool
    preserveTimestamps  bool
}

// CreateLayer builds a layer from directory contents
func (f *LayerFactory) CreateLayer(contextDir string) (v1.Layer, error) {
    // Walk directory tree
    // Create tar archive
    // Add files with metadata
    // Return as v1.Layer
}

// addFileToTar adds a single file to tar archive
func (f *LayerFactory) addFileToTar(tw *tar.Writer, path string, info os.FileInfo) error {
    // Create tar header
    // Handle symlinks
    // Write file content
    // Preserve permissions if requested
}
```

### Component 3: Config Generation
**File**: `pkg/builder/config.go` (80 lines)

```go
package builder

import (
    "time"
    
    "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/mutate"
)

// ConfigFactory generates OCI configurations
type ConfigFactory struct {
    platform v1.Platform
}

// GenerateConfig creates OCI config for image
func (f *ConfigFactory) GenerateConfig(opts BuildOptions) (*v1.ConfigFile, error) {
    return &v1.ConfigFile{
        Architecture: opts.Platform.Architecture,
        OS:          opts.Platform.OS,
        Created:     v1.Time{Time: time.Now()},
        Config: v1.Config{
            Env:        opts.Env,
            Cmd:        opts.Cmd,
            WorkingDir: opts.WorkingDir,
            Entrypoint: opts.Entrypoint,
            Labels:     opts.Labels,
        },
    }, nil
}

// ApplyConfig applies configuration to image
func (f *ConfigFactory) ApplyConfig(img v1.Image, config *v1.ConfigFile) (v1.Image, error) {
    return mutate.ConfigFile(img, config)
}
```

### Component 4: Tarball Export
**File**: `pkg/builder/tarball.go` (60 lines)

```go
package builder

import (
    "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/tarball"
)

// TarballWriter exports images as OCI tarballs
type TarballWriter struct {
    format tarball.Format
}

// NewTarballWriter creates a tarball writer
func NewTarballWriter() *TarballWriter {
    return &TarballWriter{
        format: tarball.OCIFormat,
    }
}

// Write exports image to tarball file
func (w *TarballWriter) Write(img v1.Image, path string, ref string) error {
    // Open output file
    // Write OCI format tarball
    // Include manifest and config
    return tarball.WriteToFile(path, nil, img)
}
```

## 🧪 Testing Strategy

### Unit Tests Required
- [ ] Test file: `pkg/builder/builder_test.go`
- [ ] Coverage target: 80%
- [ ] Test cases:
  - [ ] Build image from simple directory
  - [ ] Build with base image
  - [ ] Build with custom platform
  - [ ] Build with environment variables
  - [ ] Build with labels
  - [ ] Tarball export
  - [ ] Empty directory handling
  - [ ] Large file handling
  - [ ] Symlink handling
  - [ ] Permission preservation

### Test Fixtures
```go
// pkg/builder/testdata/ structure:
testdata/
├── Dockerfile          # Sample Dockerfile for reference
├── content/
│   ├── app.txt        # Sample application file
│   └── config.yaml    # Sample config
└── empty/             # Empty directory test case
```

### Integration Points Testing
- [ ] Integration with Phase 1 certificate infrastructure
- [ ] Platform detection and configuration
- [ ] OCI spec compliance validation

## 🔄 Integration Points

### With Phase 1 Components
While this effort doesn't directly use certificates, it prepares images that will be pushed using the certificate infrastructure:
- Images built here will be pushed by E2.1.2 using Phase 1's TLS handling
- No direct imports from Phase 1 in this effort (clean separation)

### With E2.1.2 (gitea-registry-client)
- Provides v1.Image objects for pushing
- Shares BuildOptions structure
- Independent implementation allows parallel development

### Feature Flag Integration (R307)
```go
// Feature flags for incomplete functionality
if opts.FeatureFlags["multi-stage-build"] {
    // Future: Multi-stage build support
    return nil, fmt.Errorf("multi-stage builds not yet implemented")
}

if opts.FeatureFlags["buildkit-frontend"] {
    // Future: BuildKit frontend support
    return nil, fmt.Errorf("BuildKit frontend not yet implemented")
}
```

## 📏 Size Management Strategy

### Measurement Protocol
```bash
# Find project root and line counter
EFFORT_DIR="/home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase2/wave1/go-containerregistry-image-builder"
cd "$EFFORT_DIR"

PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ]; then 
        break
    fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Measure using official tool (R198 compliance)
BASE_BRANCH="idpbuilder-oci-go-cr/phase1-integration-20250902-194557"
CURRENT_BRANCH=$(git branch --show-current)
$PROJECT_ROOT/tools/line-counter.sh -b $BASE_BRANCH -c $CURRENT_BRANCH

# Check points:
# - After builder.go implementation (~180 lines)
# - After layer.go implementation (~300 lines)
# - After config.go implementation (~380 lines)
# - After tarball.go implementation (~440 lines)
# - After tests implementation (~520 lines)
# - Final measurement (~600 lines)
```

### Split Prevention
- Core builder logic kept concise (~180 lines)
- Layer creation separated (~120 lines)
- Config generation isolated (~80 lines)
- Tarball export minimal (~60 lines)
- Test suite compact (~80 lines)
- Total well under 800 line limit with 200 line buffer

## 🔨 Implementation Sequence

### Day 1 - Core Implementation (4 hours)

1. **Setup and Interfaces** (30 min)
   - Create pkg/builder directory structure
   - Define Builder interface
   - Create BuildOptions structure

2. **SimpleBuilder Core** (1.5 hours)
   - Implement NewBuilder constructor
   - Create Build method skeleton
   - Add platform configuration

3. **Layer Creation** (1 hour)
   - Implement LayerFactory
   - Create directory walking logic
   - Add tar archive generation

4. **Config Generation** (30 min)
   - Create ConfigFactory
   - Generate OCI configurations
   - Apply to images

5. **Tarball Export** (30 min)
   - Implement TarballWriter
   - Add OCI format export
   - Test with simple cases

### Day 2 - Testing and Polish (2 hours)

6. **Test Suite** (1 hour)
   - Create comprehensive unit tests
   - Add test fixtures
   - Achieve 80% coverage

7. **Documentation** (30 min)
   - Add godoc comments
   - Create usage examples
   - Document feature flags

8. **Integration Testing** (30 min)
   - Test with real directories
   - Validate OCI compliance
   - Performance benchmarks

## 📝 Notes for SW Engineer

### Critical Considerations
1. **OCI Compliance**: Ensure all images are OCI 1.0 compliant
2. **Platform Support**: Default to linux/amd64 but support arm64
3. **Feature Flags**: Implement R307 flags for incomplete features
4. **Clean Separation**: No direct Phase 1 imports (parallel development)
5. **Error Handling**: Clear, actionable error messages

### Example Usage
```go
// Create builder
opts := BuildOptions{
    Platform: v1.Platform{
        OS:           "linux",
        Architecture: "amd64",
    },
    Labels: map[string]string{
        "org.opencontainers.image.source": "idpbuilder",
    },
    WorkingDir: "/app",
    Cmd:       []string{"/app/start.sh"},
}

builder, err := NewBuilder(opts)
if err != nil {
    return fmt.Errorf("failed to create builder: %w", err)
}

// Build image from directory
image, err := builder.Build(ctx, "./app-directory", opts)
if err != nil {
    return fmt.Errorf("build failed: %w", err)
}

// Export as tarball
err = builder.BuildTarball(ctx, "./app-directory", "app.tar", opts)
if err != nil {
    return fmt.Errorf("tarball export failed: %w", err)
}
```

### Key Dependencies
```go
import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    
    "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/empty"
    "github.com/google/go-containerregistry/pkg/v1/mutate"
    "github.com/google/go-containerregistry/pkg/v1/tarball"
)
```

### Performance Tips
- Use streaming for large files
- Implement concurrent layer processing for multi-layer images (future)
- Cache base images locally
- Use compression for layers

## ✅ Review Checklist
- [ ] All functional requirements addressed
- [ ] Size within 600 line estimate
- [ ] Test coverage ≥ 80%
- [ ] Feature flags for R307 compliance
- [ ] No direct Phase 1 dependencies (parallel development)
- [ ] OCI spec compliance ensured
- [ ] Clear error messages
- [ ] Documentation complete
- [ ] Implementation sequence logical
- [ ] Performance requirements achievable