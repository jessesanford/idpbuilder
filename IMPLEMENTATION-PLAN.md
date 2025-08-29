# Effort 2.1.1: Buildah Build Wrapper

## Overview
**Size Target**: ~250 lines
**Purpose**: Implement container image building using Buildah Go libraries with certificate integration from Phase 1.

## Key Interfaces

```go
// pkg/build/builder.go
package build

import (
    "context"
    "github.com/cnoe-io/idpbuilder/pkg/certs/trust"
)

// Builder handles container image building operations
type Builder interface {
    // BuildImage builds a container image from a Dockerfile
    BuildImage(ctx context.Context, opts BuildOptions) (*BuildResult, error)
    
    // ListImages lists available images
    ListImages(ctx context.Context) ([]ImageInfo, error)
    
    // RemoveImage removes an image by ID
    RemoveImage(ctx context.Context, imageID string) error
    
    // TagImage tags an existing image
    TagImage(ctx context.Context, source, target string) error
}

type BuildOptions struct {
    DockerfilePath string
    ContextDir     string
    Tag            string
    Args           map[string]string
    NoCache        bool
}

type BuildResult struct {
    ImageID      string
    Repository   string
    Tag          string
    Digest       string
    Size         int64
    BuildTime    time.Duration
}
```

## Implementation Structure
1. Parse Dockerfile (50 lines)
2. Prepare build context (40 lines)
3. Configure Buildah options (40 lines)
4. Execute build with progress tracking (60 lines)
5. Handle errors and cleanup (30 lines)
6. Tag management (30 lines)

## Phase 1 Integration
- Use TrustManager from Phase 1 for certificate handling
- Integrate with audit logging for security events
- Leverage fallback strategies for registry access

## Testing Requirements
- Test successful builds with simple Dockerfile
- Test build failures and error handling
- Test certificate integration
- Test tag operations
- Achieve 80% code coverage
