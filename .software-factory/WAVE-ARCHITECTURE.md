# Phase 1 Wave 1 Architecture

## Wave Overview
**Phase**: 1 - CLI Foundation & Authentication
**Wave**: 1 - CLI Foundation & Push Command
**Purpose**: Establish CLI structure and core interfaces for OCI push functionality

## Architectural Decisions

### 1. Command Structure
- Using Cobra CLI framework (already in IDPBuilder)
- New `push` subcommand under main CLI
- Follows existing IDPBuilder command patterns

### 2. Interface Design Philosophy
- All major components defined as interfaces first
- Concrete implementations in separate packages
- Enables easy testing with mocks
- Allows future swapping of implementations

### 3. Authentication Architecture
- Credential flags on CLI (--username, --password)
- Future: Support for Docker config.json
- Future: Support for credential helpers
- Initial: Basic auth only

### 4. TLS/Certificate Handling
- InsecureSkipVerify flag for self-signed certs
- Future: Certificate bundle loading
- Critical for Gitea integration (uses self-signed)

### 5. Core Types Location
- All interfaces in pkg/types/
- Buildah interfaces in pkg/buildah/
- Registry interfaces in pkg/registry/
- Prevents circular dependencies

## Package Structure
```
pkg/
├── cmd/
│   └── push/           # Push command implementation
├── auth/               # Authentication handling
├── tls/                # TLS configuration
├── types/              # Core interfaces and types
├── buildah/            # Buildah abstraction layer
└── registry/           # Registry client abstraction
```

## Interface Definitions

### Builder Interface
```go
type Builder interface {
    Build(ctx context.Context, dockerfile string, contextDir string, opts BuildOptions) (Image, error)
    GetImage(ctx context.Context, imageRef string) (Image, error)
}
```

### Registry Interface
```go
type Registry interface {
    Push(ctx context.Context, image Image, destination string, opts PushOptions) error
    Login(ctx context.Context, registry string, creds Credentials) error
}
```

### Image Interface
```go
type Image interface {
    Digest() string
    Tag() string
    Manifest() ([]byte, error)
}
```

## Integration Points
- Integrates with existing IDPBuilder CLI structure
- Uses same logging and error handling patterns
- Follows IDPBuilder configuration approach
- Compatible with Kind cluster setup

## Testing Strategy
- Unit tests for each package
- Interface mocks for testing
- Integration tests with real Buildah
- E2E tests with local Gitea registry