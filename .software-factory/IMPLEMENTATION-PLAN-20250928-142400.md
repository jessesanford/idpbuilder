# Implementation Plan for Builder Interface Wrapping Buildah

**Created**: 2025-09-28T14:24:00Z
**Location**: /home/vscode/workspaces/idpbuilder-gitea-push/efforts/phase1/wave2/builder-interface/.software-factory/
**Phase**: 1
**Wave**: 2
**Effort**: P1W2-E1 - Builder Interface Wrapping Buildah

## Effort Metadata

**Branch**: `phase1/wave2/builder-interface`
**Base Branch**: `phase1-wave1-integration` (or main if Wave 1 not integrated)
**Can Parallelize**: Yes
**Parallel With**: All other Wave 2 efforts (registry-client, certificate-manager, stack-mapper, progress-reporter, error-handler)
**Size Estimate**: ~400 lines (NEW code added to repository)
**Dependencies**: Wave 1 types (if available) - providers.Layer, providers.Artifact

## Pre-Planning Research Results (R374 MANDATORY)

### Existing Interfaces Found
| Interface | Location | Signature | Must Implement |
|-----------|----------|-----------|----------------|
| IProvider | pkg/kind/cluster.go:47 | Provider interface for kind | No - different domain |
| gitProvider | pkg/controllers/gitrepository/git_repository.go:43 | Git provider interface | No - different domain |

### Existing Types to Consider
| Type | Location | Purpose | Relationship |
|------|----------|---------|-------------|
| Package | pkg/printer/types/internal_types.go:33 | Package representation | May use for artifact metadata |

### FORBIDDEN DUPLICATIONS (R373)
- DO NOT create alternative provider interfaces that conflict with existing ones
- DO NOT reimplement package types that already exist
- DO NOT create competing builder patterns if any exist

### REQUIRED INTEGRATIONS (R373)
- MUST wrap Buildah library behind clean interface
- MUST use context.Context for all operations
- MUST follow existing error patterns in the codebase

## 🚨🚨🚨 R355 PRODUCTION READINESS - ZERO TOLERANCE 🚨🚨🚨

This implementation MUST be production-ready from the first commit:
- ❌ NO STUBS or placeholder implementations
- ❌ NO MOCKS except in test directories
- ❌ NO hardcoded credentials or secrets
- ❌ NO static configuration values
- ❌ NO TODO/FIXME markers in code
- ❌ NO returning nil or empty for "later implementation"
- ❌ NO panic("not implemented") patterns
- ❌ NO fake or dummy data

**VIOLATION = -100% AUTOMATIC FAILURE**

## Overview

This effort creates a clean abstraction layer over Buildah for container build operations. The interface will provide a testable, mockable API that hides Buildah implementation details while exposing the core build functionality needed by the IDPBuilder OCI management system.

## EXPLICIT SCOPE (R311 MANDATORY)

### IMPLEMENT EXACTLY:
1. **Interface Definition** - `pkg/build/builder.go` (~50 lines)
   - `Builder` interface with 3 methods
   - `NewBuilder` constructor function
   - Context support for all operations

2. **Buildah Adapter** - `pkg/build/buildah_adapter.go` (~200 lines)
   - `buildahAdapter` struct implementing Builder interface
   - `Build` method - wraps Buildah build operations
   - `AddLayer` method - adds layers to container
   - `Finalize` method - completes build and returns result
   - Proper error handling with context

3. **Type Definitions** - `pkg/build/types.go` (~100 lines)
   - `BuildConfig` struct with 6 fields
   - `BuildResult` struct with 4 fields
   - `LayerSpec` struct with 5 fields
   - `BuildOptions` struct with 3 fields
   - NO methods on any structs

4. **Error Handling** - `pkg/build/errors.go` (~50 lines)
   - `BuildError` struct with error wrapping
   - 5 error constant definitions
   - `WrapBuildError` function
   - `IsBuildError` helper function

**TOTAL ESTIMATED**: ~400 lines

### DO NOT IMPLEMENT:
- ❌ Actual container building logic (just interface wrapping)
- ❌ Registry push operations (separate effort E1.2.2)
- ❌ Stack configuration parsing (E1.2.4 handles this)
- ❌ Authentication mechanisms
- ❌ Caching layers
- ❌ Comprehensive validation (minimal only)
- ❌ Complex error recovery
- ❌ Progress reporting (E1.2.5 handles this)
- ❌ Certificate handling (E1.2.3 handles this)
- ❌ Image manifest manipulation beyond basic
- ❌ Multi-platform builds
- ❌ Build optimization strategies

## File Structure

```
efforts/phase1/wave2/builder-interface/
├── .software-factory/
│   └── IMPLEMENTATION-PLAN-20250928-142400.md (this file)
├── pkg/
│   └── build/
│       ├── builder.go         # Main interface definition (~50 lines)
│       ├── buildah_adapter.go # Buildah wrapper implementation (~200 lines)
│       ├── types.go           # Type definitions (~100 lines)
│       └── errors.go          # Error handling (~50 lines)
└── pkg/build/testdata/        # Test fixtures (not counted)
```

## Implementation Details

### 1. builder.go - Interface Definition

```go
package build

import (
    "context"
    "io"
)

// Builder defines the interface for container build operations
type Builder interface {
    // Build creates a new container build context
    Build(ctx context.Context, config *BuildConfig) (*BuildContext, error)

    // AddLayer adds a layer to the build
    AddLayer(ctx context.Context, buildCtx *BuildContext, layer *LayerSpec) error

    // Finalize completes the build and returns the result
    Finalize(ctx context.Context, buildCtx *BuildContext) (*BuildResult, error)
}

// NewBuilder creates a new Builder instance
func NewBuilder(opts *BuildOptions) (Builder, error) {
    // Validate options
    // Return buildahAdapter instance
}

// BuildContext holds the state of an in-progress build
type BuildContext struct {
    ID         string
    WorkingDir string
    // Internal fields for Buildah state
}
```

### 2. buildah_adapter.go - Buildah Wrapper

```go
package build

import (
    "context"
    "github.com/containers/buildah"
    // Other Buildah imports
)

type buildahAdapter struct {
    store    storage.Store
    options  *BuildOptions
}

func (b *buildahAdapter) Build(ctx context.Context, config *BuildConfig) (*BuildContext, error) {
    // Initialize Buildah builder
    // Set up build context
    // Return wrapped context
}

func (b *buildahAdapter) AddLayer(ctx context.Context, buildCtx *BuildContext, layer *LayerSpec) error {
    // Convert LayerSpec to Buildah operations
    // Add layer using Buildah APIs
    // Handle errors with proper wrapping
}

func (b *buildahAdapter) Finalize(ctx context.Context, buildCtx *BuildContext) (*BuildResult, error) {
    // Commit the container
    // Generate manifest
    // Return result with image ID
}
```

### 3. types.go - Type Definitions

```go
package build

import "time"

// BuildConfig specifies the configuration for a build
type BuildConfig struct {
    BaseImage    string            // Base image reference
    WorkingDir   string            // Working directory for build
    Env          map[string]string // Environment variables
    Labels       map[string]string // Container labels
    Entrypoint   []string          // Container entrypoint
    Cmd          []string          // Container command
}

// LayerSpec defines a layer to add to the container
type LayerSpec struct {
    Source      string            // Source path or URL
    Destination string            // Destination in container
    Type        LayerType         // Type of layer operation
    Permissions string            // File permissions
    Owner       string            // File ownership
}

// BuildResult contains the result of a successful build
type BuildResult struct {
    ImageID     string            // Generated image ID
    Digest      string            // Image digest
    Size        int64             // Image size in bytes
    CreatedAt   time.Time         // Build timestamp
}

// BuildOptions configures the builder behavior
type BuildOptions struct {
    StoragePath string            // Path to storage directory
    RunRoot     string            // Runtime root directory
    Debug       bool              // Enable debug logging
}

// LayerType defines the type of layer operation
type LayerType string

const (
    LayerTypeCopy LayerType = "copy"
    LayerTypeAdd  LayerType = "add"
    LayerTypeRun  LayerType = "run"
)
```

### 4. errors.go - Error Handling

```go
package build

import (
    "errors"
    "fmt"
)

// Error constants
var (
    ErrInvalidConfig     = errors.New("invalid build configuration")
    ErrBuildFailed       = errors.New("build operation failed")
    ErrLayerAddFailed    = errors.New("failed to add layer")
    ErrFinalizeFailed    = errors.New("failed to finalize build")
    ErrStorageInit       = errors.New("failed to initialize storage")
)

// BuildError wraps build-related errors with context
type BuildError struct {
    Op      string // Operation that failed
    Err     error  // Underlying error
    Context string // Additional context
}

func (e *BuildError) Error() string {
    return fmt.Sprintf("build error in %s: %v (context: %s)", e.Op, e.Err, e.Context)
}

func (e *BuildError) Unwrap() error {
    return e.Err
}

// WrapBuildError wraps an error with build context
func WrapBuildError(op string, err error, context string) error {
    if err == nil {
        return nil
    }
    return &BuildError{
        Op:      op,
        Err:     err,
        Context: context,
    }
}

// IsBuildError checks if an error is a BuildError
func IsBuildError(err error) bool {
    var buildErr *BuildError
    return errors.As(err, &buildErr)
}
```

## Configuration Requirements (R355 Mandatory)

### CORRECT Production-Ready Patterns

All configuration MUST come from environment or explicit parameters:

```go
// ✅ CORRECT - Configuration from options
func NewBuilder(opts *BuildOptions) (Builder, error) {
    if opts.StoragePath == "" {
        opts.StoragePath = os.Getenv("BUILDAH_STORAGE_PATH")
        if opts.StoragePath == "" {
            return nil, errors.New("storage path not configured")
        }
    }
    // Full implementation required
}

// ❌ WRONG - Hardcoded values
storagePath := "/var/lib/containers/storage"  // VIOLATION!

// ❌ WRONG - Stub implementation
func (b *buildahAdapter) Build(ctx context.Context, config *BuildConfig) (*BuildContext, error) {
    // TODO: implement later
    return nil, nil  // VIOLATION!
}
```

## Testing Strategy

### Unit Tests Required (60% minimum coverage):
- Interface contract tests
- Error handling validation
- Type marshaling/unmarshaling
- Mock implementation for testing other components

### Files to Create:
- `pkg/build/builder_test.go`
- `pkg/build/buildah_adapter_test.go`
- `pkg/build/mock.go` (mock implementation for other efforts)

## Size Management

**Current Estimate**: ~400 lines
**Measurement Tool**: `$CLAUDE_PROJECT_DIR/tools/line-counter.sh`
**Check Frequency**: After each file completion
**Split Threshold**: 700 lines (warning), 800 lines (stop)

## Size Limit Clarification (R359)

- The 800-line limit applies to NEW CODE ADDED
- Repository will grow by ~400 lines (EXPECTED)
- NEVER delete existing code to meet size limits
- Current codebase size is baseline, we ADD to it

## Atomic PR Design (R220 Mandatory)

### PR Summary
**Single PR implementing**: Builder interface abstraction over Buildah

### Can Merge to Main Alone
**Yes** - This PR introduces only interfaces and wrappers with no breaking changes

### R355 Production Ready Checklist
- ✅ No hardcoded values - all config from options/env
- ✅ All config from env - BuildOptions pattern
- ✅ No stub implementations - full Buildah wrapping
- ✅ No TODO markers - complete implementation
- ✅ All functions complete - no placeholders

### Feature Flags Needed
**None** - This is a foundational interface that doesn't affect existing functionality

### Interface Implementations
- **Builder**: Core interface defined
- **buildahAdapter**: Full implementation wrapping Buildah
- **Mock**: Test double for unit testing

### PR Verification
- Tests pass alone: Yes
- Build remains working: Yes
- No external dependencies on unmerged code: Yes
- Backward compatible: Yes (new functionality)

## Success Criteria

### Completion Checklist
- [ ] All 4 files created with specified line counts
- [ ] Builder interface with 3 methods defined
- [ ] Buildah adapter fully implements interface
- [ ] All types defined without methods
- [ ] Error handling complete with wrapping
- [ ] No TODO/FIXME markers in code
- [ ] No hardcoded configuration values
- [ ] All functions return proper values (no nil stubs)
- [ ] 60% test coverage achieved
- [ ] Mock implementation available
- [ ] Total lines under 800 (target ~400)

## Dependencies and Integration

### Wave 1 Dependencies (if available)
- providers.Layer type (if exists, otherwise define minimal version)
- providers.Artifact type (if exists, otherwise define minimal version)

### Provides to Wave 3
- Builder interface for MVP implementation
- Buildah operations abstraction
- Clean testing boundaries

## Next Steps for SW Engineer

1. Navigate to effort directory
2. Read this plan from `.software-factory/` directory
3. Implement files in order: types.go → builder.go → errors.go → buildah_adapter.go
4. Write tests achieving 60% coverage
5. Measure with line-counter.sh regularly
6. Commit and push when complete

---

**Plan Created By**: code-reviewer
**Agent ID**: code-reviewer-phase1-wave2
**Timestamp**: 2025-09-28T14:24:00Z
**State**: EFFORT_PLAN_CREATION