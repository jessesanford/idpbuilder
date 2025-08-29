<<<<<<< HEAD
# Effort 2.1.2: Gitea Registry Client

## Overview
**Size Target**: ~250 lines
**Purpose**: Implement OCI registry operations for Gitea with full certificate support from Phase 1.
=======
# Effort 2.1.1: Buildah Build Wrapper

## Overview
**Size Target**: ~250 lines
**Purpose**: Implement container image building using Buildah Go libraries with certificate integration from Phase 1.
>>>>>>> origin/idpbuilder-oci-mvp/phase2/wave1/buildah-build-wrapper-split-001

## Key Interfaces

```go
<<<<<<< HEAD
// pkg/registry/gitea_client.go
package registry

import (
    "context"
    "github.com/cnoe-io/idpbuilder/pkg/certs"
)

// GiteaClient handles registry operations with Gitea
type GiteaClient interface {
    // Authenticate with Gitea registry
    Authenticate(ctx context.Context, creds Credentials) error
    
    // Push pushes an image to the registry
    Push(ctx context.Context, opts PushOptions) (*PushResult, error)
    
    // List lists images in a repository
    List(ctx context.Context, repository string) ([]ImageTag, error)
    
    // Pull pulls an image from the registry
    Pull(ctx context.Context, imageRef string) (*PullResult, error)
}

type Credentials struct {
    Username string
    Password string
    Token    string
}

type PushOptions struct {
    ImageID    string
    Repository string
    Tag        string
    Insecure   bool
}

type PushResult struct {
    Digest     string
    Size       int64
    PushTime   time.Duration
    Repository string
    Tag        string
=======
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
>>>>>>> origin/idpbuilder-oci-mvp/phase2/wave1/buildah-build-wrapper-split-001
}
```

## Implementation Structure
<<<<<<< HEAD
1. Authentication handling (40 lines)
2. Push operation with progress (60 lines)
3. Certificate integration from Phase 1 (40 lines)
4. List repository contents (40 lines)
5. Error handling and retries (40 lines)
6. Helper functions (30 lines)

## Phase 1 Integration
- Use CertExtractor to get Gitea certificates
- Use ChainValidator for certificate validation
- Use FallbackHandler for error recovery
- Integrate with SecurityAuditor for audit logging

## Testing Requirements
- Test authentication flows
- Test push operations with cert validation
- Test --insecure flag behavior
- Test error scenarios and recovery
=======
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
>>>>>>> origin/idpbuilder-oci-mvp/phase2/wave1/buildah-build-wrapper-split-001
- Achieve 80% code coverage
