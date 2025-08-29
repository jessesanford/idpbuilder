# Effort 2.1.2: Gitea Registry Client

## Overview
**Size Target**: ~250 lines
**Purpose**: Implement OCI registry operations for Gitea with full certificate support from Phase 1.

## Key Interfaces

```go
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
}
```

## Implementation Structure
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
- Achieve 80% code coverage
