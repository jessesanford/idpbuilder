# Phase 1: MVP Core - Detailed Implementation Plan

## Phase Overview
**Duration:** 10 days  
**Critical Path:** YES - Foundation for all subsequent phases  
**Base Branch:** `main`  
**Target Integration Branch:** `phase1-integration`  
**Prerequisites:** None - Starting from scratch

---

## Critical Libraries & Dependencies (MAINTAINER SPECIFIED)

### Required Libraries
```yaml
core_libraries:
  - name: "github.com/containers/buildah"
    version: "v1.37.0"
    reason: "Latest stable version with full Dockerfile support and OCI compatibility. Proven reliable for container builds in production."
    usage: "Core build engine for all container operations"
    
  - name: "github.com/containers/image/v5"
    version: "v5.30.0"
    reason: "Required by Buildah for image manipulation. Matches Buildah compatibility matrix."
    usage: "Image handling, registry operations, format conversion"
    
  - name: "github.com/containers/storage"
    version: "v1.53.0"
    reason: "Required by Buildah for local image storage. Latest stable version."
    usage: "Local image storage and management"
    
  - name: "k8s.io/client-go"
    version: "v0.29.0"
    reason: "Access Gitea credentials from k8s secrets. Matches idpbuilder k8s version."
    usage: "Kubernetes secret access for Gitea authentication"
```

### Interfaces to Reuse (MANDATORY)
```yaml
reused_from_previous:
  none: "First phase - creates foundation interfaces"
    
forbidden_duplications:
  - "DO NOT create separate logging - use existing idpbuilder logger pattern"
  - "DO NOT implement custom context handling - use standard context.Context"
  - "DO NOT build separate config system - hardcode MVP values"
```

---

## Wave 1.1: Essential API Contracts

### Overview
**Focus:** Define minimal interfaces for build operations  
**Dependencies:** None  
**Parallelizable:** YES - Both efforts can run concurrently

### E1.1.1: Minimal Build Types
**Branch:** `phase1/wave1/effort1-build-types`  
**Duration:** 4 hours  
**Estimated Lines:** 100 lines  
**Agent Assignment:** Single

#### Requirements
1. **MUST** implement core request/response types only:
   - BuildRequest with essential fields
   - BuildResponse with success/error info
   - NO complex validation - keep it simple

2. **MUST NOT**:
   - Add configuration structures
   - Implement business logic
   - Add complex field validation

#### Implementation Guidance

##### Directory Structure
```
pkg/
├── build/
│   ├── api/
│   │   ├── types.go           # ~60 lines
│   │   └── types_test.go      # ~40 lines
```

##### Core Types (Maintainer Specified)
```go
// pkg/build/api/types.go
package api

// BuildRequest represents a container build request
type BuildRequest struct {
    // DockerfilePath is the path to the Dockerfile relative to ContextDir
    DockerfilePath string `json:"dockerfilePath"`
    
    // ContextDir is the build context directory (absolute path)
    ContextDir string `json:"contextDir"`
    
    // ImageName is the target image name (without registry)
    ImageName string `json:"imageName"`
    
    // ImageTag is the target image tag
    ImageTag string `json:"imageTag"`
}

// BuildResponse represents the result of a build operation
type BuildResponse struct {
    // ImageID is the built image ID
    ImageID string `json:"imageID"`
    
    // FullTag is the complete image reference
    // Format: gitea.cnoe.localtest.me/giteaadmin/{imageName}:{imageTag}
    FullTag string `json:"fullTag"`
    
    // Success indicates if the build completed successfully
    Success bool `json:"success"`
    
    // Error contains error details if Success is false
    Error string `json:"error,omitempty"`
}

// Validate performs basic validation on BuildRequest
func (br *BuildRequest) Validate() error {
    if br.DockerfilePath == "" {
        return fmt.Errorf("DockerfilePath is required")
    }
    if br.ContextDir == "" {
        return fmt.Errorf("ContextDir is required")
    }
    if br.ImageName == "" {
        return fmt.Errorf("ImageName is required")
    }
    if br.ImageTag == "" {
        br.ImageTag = "latest"  // Default tag
    }
    return nil
}
```

#### Test Requirements (TDD)
```go
// pkg/build/api/types_test.go
func TestBuildRequestValidation(t *testing.T) {
    testCases := []struct {
        name     string
        request  BuildRequest
        wantErr  bool
    }{
        {
            name: "valid request",
            request: BuildRequest{
                DockerfilePath: "Dockerfile",
                ContextDir:     "/tmp/build",
                ImageName:      "myapp",
                ImageTag:       "v1.0",
            },
            wantErr: false,
        },
        {
            name: "missing dockerfile path",
            request: BuildRequest{
                ContextDir: "/tmp/build",
                ImageName:  "myapp",
            },
            wantErr: true,
        },
        {
            name: "default tag applied",
            request: BuildRequest{
                DockerfilePath: "Dockerfile",
                ContextDir:     "/tmp/build",
                ImageName:      "myapp",
            },
            wantErr: false,
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            err := tc.request.Validate()
            if tc.wantErr && err == nil {
                t.Error("expected error but got none")
            }
            if !tc.wantErr && err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}
```

#### Success Criteria
- [ ] Types compile without errors
- [ ] Validation tests pass
- [ ] No external dependencies beyond standard library
- [ ] Under 100 lines per line-counter.sh

### E1.1.2: Builder Interface
**Branch:** `phase1/wave1/effort2-builder-interface`  
**Duration:** 2 hours  
**Estimated Lines:** 50 lines  
**Agent Assignment:** Single (can run parallel with E1.1.1)

#### Requirements
1. **MUST** define simple Builder interface
2. **MUST** be compatible with E1.1.1 types
3. **MUST NOT** implement the interface (just define it)

#### Implementation Guidance

##### Interface Definition (Maintainer Specified)
```go
// pkg/build/api/builder.go
package api

import "context"

// Builder defines the interface for container build operations
type Builder interface {
    // BuildAndPush builds a container image and pushes to registry
    // This is the only method needed for MVP
    BuildAndPush(ctx context.Context, req BuildRequest) (*BuildResponse, error)
}

// BuilderConfig holds configuration for builder instances
type BuilderConfig struct {
    // Registry is the target OCI registry
    Registry string
    
    // Namespace is the registry namespace (hardcoded for MVP)
    Namespace string
    
    // InsecureSkipTLSVerify skips TLS verification (default true for MVP)
    InsecureSkipTLSVerify bool
}

// DefaultConfig returns MVP configuration
func DefaultConfig() BuilderConfig {
    return BuilderConfig{
        Registry:              "gitea.cnoe.localtest.me",
        Namespace:             "giteaadmin",
        InsecureSkipTLSVerify: true,
    }
}
```

#### Test Requirements
```go
// pkg/build/api/builder_test.go
func TestBuilderInterface(t *testing.T) {
    // Test that interface can be implemented
    var builder Builder = &mockBuilder{}
    
    req := BuildRequest{
        DockerfilePath: "Dockerfile",
        ContextDir:     "/tmp",
        ImageName:      "test",
        ImageTag:       "latest",
    }
    
    _, err := builder.BuildAndPush(context.Background(), req)
    if err != nil {
        t.Errorf("interface implementation failed: %v", err)
    }
}

type mockBuilder struct{}

func (m *mockBuilder) BuildAndPush(ctx context.Context, req BuildRequest) (*BuildResponse, error) {
    return &BuildResponse{Success: true}, nil
}
```

---

## Wave 1.2: Core Libraries

### Overview
**Focus:** Implement Buildah client and authentication  
**Dependencies:** Wave 1.1 complete  
**Parallelizable:** NO - Sequential implementation required

### E1.2.1: Buildah Client Wrapper
**Branch:** `phase1/wave2/effort1-buildah-client`  
**Duration:** 16 hours  
**Estimated Lines:** 300 lines  
**Agent Assignment:** Single

#### Requirements
1. **MUST** wrap Buildah with our interface
2. **MUST** handle storage initialization
3. **MUST** configure for insecure registries
4. **MUST** implement basic build functionality

#### Implementation Guidance

##### Directory Structure
```
pkg/
├── build/
│   ├── buildah/
│   │   ├── client.go         # ~150 lines
│   │   ├── client_test.go    # ~100 lines
│   │   ├── storage.go        # ~50 lines
```

##### Buildah Client (Maintainer Specified)
```go
// pkg/build/buildah/client.go
package buildah

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    
    "github.com/containers/buildah"
    "github.com/containers/buildah/define"
    "github.com/containers/image/v5/types"
    "github.com/containers/storage"
    "github.com/containers/storage/pkg/unshare"
    
    "idpbuilder/pkg/build/api"
)

// Client implements the Builder interface using Buildah
type Client struct {
    store         storage.Store
    systemContext *types.SystemContext
    config        api.BuilderConfig
}

// NewClient creates a new Buildah client
func NewClient(config api.BuilderConfig) (*Client, error) {
    // CRITICAL: Ensure we run in user namespace for rootless
    unshare.MaybeReexecUsingUserNamespace(true)
    
    // Configure system context for insecure registries
    systemContext := &types.SystemContext{
        DockerInsecureSkipTLSVerify:    types.NewOptionalBool(config.InsecureSkipTLSVerify),
        DockerDaemonInsecureSkipTLSVerify: true,
        OCIInsecureSkipTLSVerify:       config.InsecureSkipTLSVerify,
    }
    
    // Initialize storage
    store, err := initStorage()
    if err != nil {
        return nil, fmt.Errorf("failed to initialize storage: %w", err)
    }
    
    return &Client{
        store:         store,
        systemContext: systemContext,
        config:        config,
    }, nil
}

// BuildAndPush implements the Builder interface
func (c *Client) BuildAndPush(ctx context.Context, req api.BuildRequest) (*api.BuildResponse, error) {
    // Validate request
    if err := req.Validate(); err != nil {
        return &api.BuildResponse{
            Success: false,
            Error:   fmt.Sprintf("validation failed: %v", err),
        }, nil
    }
    
    // Build the image
    imageID, err := c.build(ctx, req)
    if err != nil {
        return &api.BuildResponse{
            Success: false,
            Error:   fmt.Sprintf("build failed: %v", err),
        }, nil
    }
    
    // Create full tag
    fullTag := fmt.Sprintf("%s/%s/%s:%s", 
        c.config.Registry, c.config.Namespace, req.ImageName, req.ImageTag)
    
    // Tag the image
    if err := c.tag(ctx, imageID, fullTag); err != nil {
        return &api.BuildResponse{
            Success: false,
            Error:   fmt.Sprintf("tag failed: %v", err),
        }, nil
    }
    
    // Push to registry
    if err := c.push(ctx, fullTag); err != nil {
        return &api.BuildResponse{
            Success: false,
            Error:   fmt.Sprintf("push failed: %v", err),
        }, nil
    }
    
    return &api.BuildResponse{
        ImageID: imageID,
        FullTag: fullTag,
        Success: true,
    }, nil
}
```

#### Test Requirements (TDD)
```go
// pkg/build/buildah/client_test.go
func TestClientCreation(t *testing.T) {
    config := api.DefaultConfig()
    client, err := NewClient(config)
    
    if err != nil {
        t.Fatalf("failed to create client: %v", err)
    }
    
    if client == nil {
        t.Fatal("client is nil")
    }
    
    if client.store == nil {
        t.Error("store not initialized")
    }
    
    if client.systemContext == nil {
        t.Error("system context not initialized")
    }
}

func TestBuildAndPushValidation(t *testing.T) {
    client := &Client{} // Mock for testing
    
    req := api.BuildRequest{} // Invalid request
    resp, err := client.BuildAndPush(context.Background(), req)
    
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    
    if resp.Success {
        t.Error("expected validation failure")
    }
    
    if resp.Error == "" {
        t.Error("expected error message")
    }
}
```

### E1.2.2: Gitea Auth Helper
**Branch:** `phase1/wave2/effort2-gitea-auth`  
**Duration:** 8 hours  
**Estimated Lines:** 150 lines  
**Agent Assignment:** Single

#### Requirements
1. **MUST** fetch Gitea credentials from k8s secret
2. **MUST** fallback to hardcoded credentials for testing
3. **MUST** handle k8s client unavailability gracefully

#### Implementation Guidance

##### Authentication Helper (Maintainer Specified)
```go
// pkg/build/auth/gitea.go
package auth

import (
    "context"
    "fmt"
    
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)

const (
    // Gitea credential secret in idpbuilder cluster
    GiteaSecretName      = "gitea-credential"
    GiteaSecretNamespace = "gitea"
    GiteaUsernameKey     = "username"
    GiteaPasswordKey     = "password"
    
    // Fallback credentials for testing
    FallbackUsername = "giteaAdmin"
    FallbackPassword = "password"  // Will be read from secret in real usage
)

// Credentials represents Gitea authentication
type Credentials struct {
    Username string
    Password string
}

// GetGiteaCredentials retrieves Gitea credentials from cluster
func GetGiteaCredentials(ctx context.Context) (*Credentials, error) {
    // Try to get credentials from k8s secret first
    if creds, err := getFromKubernetes(ctx); err == nil {
        return creds, nil
    }
    
    // Fallback to default credentials for testing/development
    return &Credentials{
        Username: FallbackUsername,
        Password: FallbackPassword,
    }, nil
}

// getFromKubernetes fetches credentials from k8s secret
func getFromKubernetes(ctx context.Context) (*Credentials, error) {
    // Create in-cluster config
    config, err := rest.InClusterConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to create k8s config: %w", err)
    }
    
    // Create clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, fmt.Errorf("failed to create k8s client: %w", err)
    }
    
    // Get the secret
    secret, err := clientset.CoreV1().Secrets(GiteaSecretNamespace).
        Get(ctx, GiteaSecretName, metav1.GetOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to get secret: %w", err)
    }
    
    // Extract credentials
    username, ok := secret.Data[GiteaUsernameKey]
    if !ok {
        return nil, fmt.Errorf("username not found in secret")
    }
    
    password, ok := secret.Data[GiteaPasswordKey]
    if !ok {
        return nil, fmt.Errorf("password not found in secret")
    }
    
    return &Credentials{
        Username: string(username),
        Password: string(password),
    }, nil
}
```

---

## Wave 1.3: MVP Implementation

### Overview
**Focus:** Complete build and push functionality  
**Dependencies:** Wave 1.2 complete  
**Parallelizable:** NO - Sequential implementation required

### E1.3.1: Basic Build Implementation
**Branch:** `phase1/wave3/effort1-basic-build`  
**Duration:** 20 hours  
**Estimated Lines:** 400 lines  
**Agent Assignment:** Single

#### Requirements
1. **MUST** implement Dockerfile parsing and execution
2. **MUST** handle FROM, COPY, RUN, CMD instructions
3. **MUST** create buildable images
4. **MUST NOT** implement complex features (multi-stage, args, etc.)

#### Implementation Guidance

##### Build Implementation (Maintainer Specified)
```go
// pkg/build/buildah/build.go
func (c *Client) build(ctx context.Context, req api.BuildRequest) (string, error) {
    dockerfilePath := filepath.Join(req.ContextDir, req.DockerfilePath)
    
    // Read Dockerfile
    dockerfileContent, err := os.ReadFile(dockerfilePath)
    if err != nil {
        return "", fmt.Errorf("failed to read Dockerfile: %w", err)
    }
    
    // Parse Dockerfile instructions
    instructions, err := parseDockerfile(dockerfileContent)
    if err != nil {
        return "", fmt.Errorf("failed to parse Dockerfile: %w", err)
    }
    
    // Create builder from scratch or base image
    var builderOptions buildah.BuilderOptions
    baseImage := "scratch"
    
    // Find FROM instruction
    for _, inst := range instructions {
        if inst.Command == "FROM" {
            baseImage = inst.Args[0]
            break
        }
    }
    
    builderOptions = buildah.BuilderOptions{
        FromImage:     baseImage,
        SystemContext: c.systemContext,
    }
    
    builder, err := buildah.NewBuilder(ctx, c.store, builderOptions)
    if err != nil {
        return "", fmt.Errorf("failed to create builder: %w", err)
    }
    defer builder.Delete()
    
    // Execute instructions
    for _, inst := range instructions {
        if err := c.executeInstruction(ctx, builder, inst, req.ContextDir); err != nil {
            return "", fmt.Errorf("failed to execute %s: %w", inst.Command, err)
        }
    }
    
    // Commit the image
    imageID, _, _, err := builder.Commit(ctx, buildah.CommitOptions{
        SystemContext: c.systemContext,
    })
    if err != nil {
        return "", fmt.Errorf("failed to commit image: %w", err)
    }
    
    return imageID, nil
}

// Instruction represents a Dockerfile instruction
type Instruction struct {
    Command string
    Args    []string
}

// parseDockerfile parses Dockerfile content into instructions
func parseDockerfile(content []byte) ([]Instruction, error) {
    // Simple parser for MVP - handle basic instructions only
    lines := strings.Split(string(content), "\n")
    var instructions []Instruction
    
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        
        parts := strings.Fields(line)
        if len(parts) < 1 {
            continue
        }
        
        instructions = append(instructions, Instruction{
            Command: strings.ToUpper(parts[0]),
            Args:    parts[1:],
        })
    }
    
    return instructions, nil
}
```

### E1.3.2: Basic Push Implementation
**Branch:** `phase1/wave3/effort2-basic-push`  
**Duration:** 12 hours  
**Estimated Lines:** 250 lines  
**Agent Assignment:** Single

#### Requirements
1. **MUST** push images to Gitea registry
2. **MUST** handle authentication
3. **MUST** work with self-signed certificates

#### Implementation Guidance

##### Push Implementation (Maintainer Specified)
```go
// pkg/build/buildah/push.go
func (c *Client) push(ctx context.Context, imageRef string) error {
    // Get Gitea credentials
    creds, err := auth.GetGiteaCredentials(ctx)
    if err != nil {
        return fmt.Errorf("failed to get credentials: %w", err)
    }
    
    // Configure authentication
    c.systemContext.DockerAuthConfig = &types.DockerAuthConfig{
        Username: creds.Username,
        Password: creds.Password,
    }
    
    // Push options
    pushOptions := buildah.PushOptions{
        Store:         c.store,
        SystemContext: c.systemContext,
    }
    
    // Push to registry
    _, err = buildah.Push(ctx, imageRef, imageRef, pushOptions)
    if err != nil {
        return fmt.Errorf("push failed: %w", err)
    }
    
    return nil
}

func (c *Client) tag(ctx context.Context, imageID, tag string) error {
    // Tag the image using Buildah
    err := c.store.SetNames(imageID, []string{tag})
    if err != nil {
        return fmt.Errorf("failed to tag image: %w", err)
    }
    
    return nil
}
```

### E1.3.3: BuildAndPush Orchestrator
**Branch:** `phase1/wave3/effort3-orchestrator`  
**Duration:** 8 hours  
**Estimated Lines:** 200 lines  
**Agent Assignment:** Single

#### Requirements
1. **MUST** integrate build and push operations
2. **MUST** provide user-friendly service interface
3. **MUST** handle errors gracefully

#### Implementation Guidance

##### Service Implementation (Maintainer Specified)
```go
// pkg/build/service.go
package build

import (
    "context"
    "fmt"
    
    "idpbuilder/pkg/build/api"
    "idpbuilder/pkg/build/buildah"
)

// Service provides high-level build operations
type Service struct {
    client api.Builder
}

// NewService creates a new build service
func NewService() (*Service, error) {
    config := api.DefaultConfig()
    
    client, err := buildah.NewClient(config)
    if err != nil {
        return nil, fmt.Errorf("failed to create buildah client: %w", err)
    }
    
    return &Service{
        client: client,
    }, nil
}

// BuildAndPush builds and pushes a container image
func (s *Service) BuildAndPush(ctx context.Context, req api.BuildRequest) (*api.BuildResponse, error) {
    return s.client.BuildAndPush(ctx, req)
}

// Build builds a container image without pushing
func (s *Service) Build(ctx context.Context, req api.BuildRequest) (*api.BuildResponse, error) {
    // For MVP, we always push. Phase 2 will add build-only option
    return s.BuildAndPush(ctx, req)
}
```

---

## Phase-Wide Constraints

### Architecture Decisions (Maintainer Specified)
```markdown
1. **Storage Pattern**
   - MUST use Buildah's default storage location
   - MUST handle rootless containers properly
   - NO custom storage backends in MVP

2. **Error Handling**
   - Return errors in BuildResponse.Error field
   - NO panic/fatal - always return controlled errors
   - Log errors but don't expose internal details

3. **Registry Configuration**
   - Hardcode gitea.cnoe.localtest.me for MVP
   - Hardcode giteaadmin namespace for MVP
   - InsecureSkipVerify = true by default
```

### Forbidden Duplications
- DO NOT implement custom HTTP client - use containers/image defaults
- DO NOT create custom storage layer - use containers/storage
- DO NOT implement custom Dockerfile parser beyond basic MVP needs
- DO NOT add configuration files - hardcode values for MVP

---

## Testing Strategy

### Phase-Level Testing
1. **Unit Tests**: Each package must have >70% coverage
2. **Integration Tests**: Full build+push workflow test
3. **Manual Tests**: Real Dockerfile build and registry verification

### Test Data
```dockerfile
# test/fixtures/simple/Dockerfile
FROM alpine:3.19
COPY app /app
RUN chmod +x /app
CMD ["/app"]
```

---

## Branch Strategy

### Working Branches
```bash
# Wave 1 integration
git checkout -b phase1/wave1-integration
git merge --no-ff phase1/wave1/effort1-build-types
git merge --no-ff phase1/wave1/effort2-builder-interface

# Wave 2 integration  
git checkout -b phase1/wave2-integration
git merge --no-ff phase1/wave2/effort1-buildah-client
git merge --no-ff phase1/wave2/effort2-gitea-auth

# Wave 3 integration
git checkout -b phase1/wave3-integration
git merge --no-ff phase1/wave3/effort1-basic-build
git merge --no-ff phase1/wave3/effort2-basic-push
git merge --no-ff phase1/wave3/effort3-orchestrator

# Phase integration
git checkout -b phase1-integration
git merge --no-ff phase1/wave1-integration
git merge --no-ff phase1/wave2-integration
git merge --no-ff phase1/wave3-integration
```

---

## Size Management

### Estimated Total: 750 lines
- Wave 1: ~150 lines
- Wave 2: ~450 lines  
- Wave 3: ~850 lines (likely needs split)

### Split Strategy for Wave 3
If Wave 3 exceeds 800 lines:
1. E1.3.1: Core build implementation (~400 lines)
2. E1.3.2: Push and tag implementation (~250 lines)
3. E1.3.3: Service orchestration (~200 lines)

Each split must be independently buildable and testable.

---

## Success Criteria

### Functional
- [ ] Can build simple Dockerfile with FROM, COPY, RUN, CMD
- [ ] Can push to gitea.cnoe.localtest.me successfully
- [ ] Authentication works with k8s secrets or fallback
- [ ] Handles basic error conditions gracefully

### Quality  
- [ ] All tests pass with >70% coverage
- [ ] All packages compile without warnings
- [ ] Line count under 800 per effort (split if needed)
- [ ] Integration test demonstrates full workflow

### Integration
- [ ] Buildah libraries properly integrated
- [ ] Kubernetes client-go works for auth
- [ ] Compatible with idpbuilder cluster environment
- [ ] Foundation ready for Phase 2 CLI integration

---

## Risk Mitigation

### Technical Risks
| Risk | Mitigation | Owner |
|------|------------|-------|
| Buildah API complexity | Start with minimal feature set | SW Engineer |
| Storage permission issues | Use rootless configuration | SW Engineer |
| Registry auth failures | Implement fallback credentials | SW Engineer |
| Certificate issues | Default to InsecureSkipVerify | Maintainer Decision |

---

## Handoff Instructions

### For Orchestrator
1. Create effort working directories under `/workspaces/efforts/phase1/`
2. Ensure each effort stays under 800 lines
3. Task Code Reviewer to create IMPLEMENTATION-PLAN.md for each effort
4. Execute waves sequentially - no parallel wave execution

### For SW Engineer  
1. Follow Buildah best practices for rootless containers
2. Reuse standard library and Buildah patterns - no custom implementations
3. Implement only assigned effort scope - no feature creep
4. Include comprehensive tests as specified

### For Code Reviewer
1. Verify no duplication of Buildah functionality
2. Check line counts after every 200 lines of implementation
3. Ensure proper error handling patterns
4. Validate integration test covers full workflow

This phase establishes the foundational container build capability that all subsequent phases will build upon. Focus on getting the MVP working reliably rather than adding features.