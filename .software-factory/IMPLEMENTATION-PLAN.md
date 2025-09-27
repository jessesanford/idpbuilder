# P1W1-E4: CLI Interface Contracts Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Effort**: E1.1.4 - Core Interfaces
**Branch**: `igp/phase1/wave1/effort-1.1.4-core-interfaces`
**Can Parallelize**: Yes
**Parallel With**: [E1.1.1, E1.1.2, E1.1.3]
**Size Estimate**: ~350 lines
**Dependencies**: None (can start immediately)

## Overview
- **Effort**: Core interfaces and types for the entire project
- **Phase**: 1 (CLI Foundation & Authentication)
- **Wave**: 1 (CLI Foundation & Push Command)
- **Estimated Size**: ~350 lines
- **Implementation Time**: 3-4 hours

## Purpose
This effort establishes the foundational interface contracts for the IDPBuilder OCI support. These interfaces will define the contracts for CLI commands, registry operations, and the Buildah abstraction layer. All subsequent efforts will depend on these interface definitions.

## File Structure
```
pkg/
├── types/
│   ├── interfaces.go        # ~100 lines - Core interfaces for all operations
│   ├── errors.go            # ~60 lines  - Error types and handling
│   └── config.go            # ~60 lines  - Configuration types
├── buildah/
│   └── interfaces.go        # ~70 lines  - Buildah abstraction interfaces
└── registry/
    └── interfaces.go        # ~60 lines  - Registry client interfaces
```

## Implementation Steps

### Step 1: Core Type Interfaces (pkg/types/interfaces.go)
**Size**: ~100 lines

Define the foundational interfaces that all components will implement:
```go
// Core interfaces for OCI operations
type OCIClient interface {
    Push(ctx context.Context, image ImageReference, opts PushOptions) error
    Pull(ctx context.Context, image ImageReference, opts PullOptions) error
    List(ctx context.Context, repo Repository, opts ListOptions) ([]ImageReference, error)
    Delete(ctx context.Context, image ImageReference, opts DeleteOptions) error
}

type ImageBuilder interface {
    Build(ctx context.Context, src BuildSource, opts BuildOptions) (ImageReference, error)
    Tag(ctx context.Context, image ImageReference, tags []string) error
    Inspect(ctx context.Context, image ImageReference) (*ImageInfo, error)
}

type Authenticator interface {
    GetCredentials(ctx context.Context, registry string) (*Credentials, error)
    RefreshToken(ctx context.Context, creds *Credentials) (*Credentials, error)
    Validate(ctx context.Context, creds *Credentials) error
}

type ConfigProvider interface {
    GetRegistryConfig(name string) (*RegistryConfig, error)
    GetAuthConfig(registry string) (*AuthConfig, error)
    GetTLSConfig(registry string) (*TLSConfig, error)
}
```

### Step 2: Error Types (pkg/types/errors.go)
**Size**: ~60 lines

Define standard error types for consistent error handling:
```go
// Common error types for OCI operations
type RegistryError struct {
    Registry string
    Code     int
    Message  string
    Cause    error
}

type AuthenticationError struct {
    Registry string
    Realm    string
    Scope    string
    Cause    error
}

type ConfigurationError struct {
    Field   string
    Value   string
    Reason  string
}

type BuildError struct {
    Stage   string
    Message string
    Cause   error
}
```

### Step 3: Configuration Types (pkg/types/config.go)
**Size**: ~60 lines

Define configuration structures:
```go
type RegistryConfig struct {
    Name        string
    URL         string
    Type        string // "docker", "harbor", "gitea", "generic"
    Insecure    bool
    SkipVerify  bool
    CABundle    string
}

type AuthConfig struct {
    Username      string
    Password      string
    Token         string
    IdentityToken string
    AuthType      string // "basic", "bearer", "oauth2"
}

type TLSConfig struct {
    CAFile     string
    CertFile   string
    KeyFile    string
    SkipVerify bool
}

type PushOptions struct {
    Insecure   bool
    SkipVerify bool
    Force      bool
    Quiet      bool
}
```

### Step 4: Buildah Interfaces (pkg/buildah/interfaces.go)
**Size**: ~70 lines

Define Buildah abstraction interfaces:
```go
type BuildahClient interface {
    // Container operations
    NewContainer(ctx context.Context, base string, opts ContainerOptions) (Container, error)
    CommitContainer(ctx context.Context, container Container, ref ImageReference) error
    DeleteContainer(ctx context.Context, container Container) error

    // Image operations
    PushImage(ctx context.Context, src, dest ImageReference, opts PushOptions) error
    PullImage(ctx context.Context, ref ImageReference, opts PullOptions) error
    InspectImage(ctx context.Context, ref ImageReference) (*ImageMetadata, error)
}

type Container interface {
    AddFile(src, dest string) error
    AddDirectory(src, dest string) error
    SetEnv(key, value string) error
    SetWorkDir(path string) error
    SetCmd(cmd []string) error
    SetEntrypoint(entrypoint []string) error
    Run(cmd []string) error
    Mount() (string, error)
    Unmount() error
}

type ImageMetadata struct {
    Digest      string
    Size        int64
    Created     time.Time
    Architecture string
    OS          string
    Layers      []LayerInfo
}
```

### Step 5: Registry Interfaces (pkg/registry/interfaces.go)
**Size**: ~60 lines

Define registry client interfaces:
```go
type RegistryClient interface {
    // Connection management
    Connect(ctx context.Context, config *RegistryConfig) error
    Disconnect(ctx context.Context) error
    Ping(ctx context.Context) error

    // Repository operations
    ListRepositories(ctx context.Context, opts ListOptions) ([]Repository, error)
    GetRepository(ctx context.Context, name string) (*Repository, error)
    CreateRepository(ctx context.Context, repo Repository) error
    DeleteRepository(ctx context.Context, name string) error

    // Tag operations
    ListTags(ctx context.Context, repo string, opts ListOptions) ([]string, error)
    GetManifest(ctx context.Context, ref ImageReference) (*Manifest, error)
    PutManifest(ctx context.Context, ref ImageReference, manifest *Manifest) error
    DeleteManifest(ctx context.Context, ref ImageReference) error
}

type Repository struct {
    Name        string
    Description string
    Public      bool
    Tags        []string
}

type Manifest struct {
    MediaType string
    Size      int64
    Digest    string
    Config    json.RawMessage
    Layers    []Layer
}
```

### Step 6: Testing
**Size**: Included in line count

Create basic interface compliance tests:
- Verify interfaces are properly defined
- Create mock implementations for testing
- Ensure no compilation errors

### Step 7: Integration Points
Document how these interfaces will be used:
- Push command will use OCIClient and ImageBuilder
- Auth flags will implement Authenticator
- TLS config will use ConfigProvider
- All efforts will import and implement these interfaces

## Size Management
- **Estimated Lines**: ~350 lines
- **Measurement Tool**: `${PROJECT_ROOT}/tools/line-counter.sh`
- **Check Frequency**: After each major interface definition
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

## Test Requirements
- **Unit Tests**: Interface compliance tests
- **Compilation Tests**: Ensure all interfaces compile cleanly
- **Mock Generation**: Provide mock implementations for testing
- **Test Coverage Target**: N/A for interfaces (will be tested via implementations)

## Pattern Compliance
- **Go Patterns**: Follow Go interface best practices
  - Small, focused interfaces
  - Interface segregation principle
  - Accept interfaces, return structs
- **Naming Conventions**: Use standard Go naming (e.g., `Reader`, `Writer`, `Client`)
- **Documentation**: Comprehensive godoc comments for all interfaces
- **Error Handling**: Use wrapped errors with context

## Dependencies
- **External Dependencies**: None (interfaces only)
- **Internal Dependencies**: None (foundational effort)
- **Import Requirements**: Only standard library imports

## Integration Notes
- All subsequent efforts in Phase 1 will import these interfaces
- Implementations will be provided in later efforts
- Changes to interfaces after this effort require coordination
- These interfaces form the contract for the entire OCI implementation

## Success Criteria
✅ All interfaces compile without errors
✅ Clear separation of concerns between interfaces
✅ Comprehensive error types defined
✅ Configuration structures support all use cases
✅ Mock implementations can be generated
✅ Documentation is clear and complete
✅ Total size under 400 lines

## Risk Mitigation
- **Risk**: Interface changes after implementation starts
  - **Mitigation**: Thorough review before implementation begins
- **Risk**: Missing interface methods discovered later
  - **Mitigation**: Interface extension pattern (new interfaces extending old)
- **Risk**: Circular dependencies
  - **Mitigation**: Careful package structure, interfaces in separate packages

## EFFORT INFRASTRUCTURE METADATA
**EFFORT_NAME**: effort-1.1.4-core-interfaces
**BRANCH**: igp/phase1/wave1/effort-1.1.4-core-interfaces
**WORKING_DIRECTORY**: /home/vscode/workspaces/idpbuilder-gitea-push/efforts/phase1/wave1/effort-1.1.4-core-interfaces
**REMOTE**: origin (https://github.com/jessesanford/idpbuilder.git)
**BASE_BRANCH**: main
**CAN_PARALLELIZE**: true
**PARALLEL_WITH**: [effort-1.1.1-push-command-skeleton, effort-1.1.2-auth-flags, effort-1.1.3-tls-config]

## Implementation Checklist
- [ ] Create pkg/types directory
- [ ] Implement interfaces.go with core interfaces
- [ ] Implement errors.go with error types
- [ ] Implement config.go with configuration types
- [ ] Create pkg/buildah directory
- [ ] Implement buildah/interfaces.go
- [ ] Create pkg/registry directory
- [ ] Implement registry/interfaces.go
- [ ] Add godoc comments to all interfaces
- [ ] Run go fmt and go vet
- [ ] Measure with line-counter.sh
- [ ] Commit and push changes