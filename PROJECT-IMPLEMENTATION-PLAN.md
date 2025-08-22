# Project Implementation Plan: Container Build and Push for idpbuilder

## Project Overview
**Project Name**: idpbuilder Container Build and Push Feature
**Goal**: Enable idpbuilder to build container images locally and push them directly to the Gitea OCI registry at gitea.cnoe.localtest.me, eliminating the need for external Docker daemon configuration
**Total Phases**: 5
**Total Waves**: 15
**Total Efforts**: 38
**Estimated Completion**: 6-7 weeks

## Implementation Strategy
- **Phase 1**: MVP - Get basic build and push working
- **Phase 2**: CLI Integration - Make it usable
- **Phase 3**: Production Ready - Tests and documentation
- **Phase 4**: Enhanced Features - Better error handling, caching
- **Phase 5**: Advanced Features - Multi-stage, build args, optimizations

## Key Design Decisions
- **Build Backend**: Buildah using native Go libraries (github.com/containers/buildah)
- **Build Context**: Local directories only
- **Registry**: Fixed to gitea.cnoe.localtest.me (idpbuilder's Gitea instance)
- **Build Trigger**: Manual CLI commands only
- **Certificate Handling**: InsecureSkipVerify by default

## Configuration
```yaml
size_limits:
  warning_threshold: 700
  error_threshold: 800
  
test_coverage:
  phase_1: 70%  # MVP - focus on functionality
  phase_2: 80%  # CLI integration
  phase_3: 90%  # Production ready with tests
  phase_4: 85%  # Enhanced features
  phase_5: 80%  # Advanced features
  
parallelization:
  max_parallel_efforts: 3
  allow_parallel_waves: false
  
review_requirements:
  code_review: mandatory
  architect_review: mandatory
  security_review: optional
```

## Phase 1: MVP Core - Minimal Working Build and Push
**Goal**: Achieve minimal viable product - build a simple Dockerfile and push to Gitea
**Duration**: 1.5 weeks
**Success Criteria**: 
- Can build a basic Dockerfile
- Can push to Gitea registry
- Works with self-signed certificates

### Wave 1: Essential API Contracts
**Dependencies**: None
**Can Parallelize**: Yes

#### E1.1.1: Minimal Build Types
- **Description**: Define only essential types for MVP
- **Requirements**:
  ```go
  // pkg/build/api/types.go
  type BuildRequest struct {
    DockerfilePath string
    ContextDir     string
    ImageName      string
    ImageTag       string
  }
  
  type BuildResponse struct {
    ImageID   string
    FullTag   string  // gitea.cnoe.localtest.me/namespace/name:tag
    Success   bool
    Error     string
  }
  ```
- **Test Requirements**: Basic validation tests
- **Estimated Size**: 100 lines

#### E1.1.2: Builder Interface
- **Description**: Simple interface for build operations
- **Requirements**:
  ```go
  // pkg/build/api/builder.go
  type Builder interface {
    BuildAndPush(ctx context.Context, req BuildRequest) (*BuildResponse, error)
  }
  ```
- **Test Requirements**: Interface compliance tests
- **Estimated Size**: 50 lines

### Wave 2: Core Libraries
**Dependencies**: Wave 1 complete
**Can Parallelize**: No

#### E1.2.1: Buildah Client Wrapper
- **Description**: Minimal Buildah client setup
- **Requirements**:
  ```go
  // pkg/build/buildah/client.go
  import (
    "github.com/containers/buildah"
    "github.com/containers/image/v5/types"
  )
  
  type Client struct {
    store         storage.Store
    systemContext *types.SystemContext
  }
  
  func NewClient() (*Client, error) {
    // Minimal setup with InsecureSkipVerify
    sc := &types.SystemContext{
      DockerInsecureSkipTLSVerify: types.OptionalBoolTrue,
    }
    // Initialize storage
    store, err := storage.GetStore(storage.StoreOptions{})
    return &Client{store: store, systemContext: sc}, err
  }
  ```
- **Test Requirements**: Client initialization tests
- **Estimated Size**: 200 lines

#### E1.2.2: Gitea Auth Helper
- **Description**: Get Gitea credentials from cluster
- **Requirements**:
  ```go
  // pkg/build/auth/gitea.go
  func GetGiteaCredentials(ctx context.Context) (username, password string, err error) {
    // Read from gitea-credential secret
    // Return hardcoded giteaAdmin if k8s client unavailable (for testing)
    return "giteaAdmin", getPasswordFromSecret(), nil
  }
  ```
- **Test Requirements**: Auth fetching tests
- **Estimated Size**: 150 lines

### Wave 3: MVP Implementation
**Dependencies**: Wave 2 complete
**Can Parallelize**: No

#### E1.3.1: Basic Build Implementation
- **Description**: Implement minimal Dockerfile build
- **Requirements**:
  ```go
  // pkg/build/buildah/build.go
  func (c *Client) Build(ctx context.Context, dockerfilePath, contextDir string) (string, error) {
    // 1. Create builder from scratch
    builder, err := buildah.NewBuilder(ctx, c.store, buildah.BuilderOptions{
      FromImage: "scratch",
      SystemContext: c.systemContext,
    })
    
    // 2. Simple Dockerfile execution (FROM, COPY, RUN only)
    // 3. Commit the image
    imageID, _, _, err := builder.Commit(ctx, imageRef, buildah.CommitOptions{})
    return imageID, err
  }
  ```
- **Test Requirements**: Build simple Dockerfile test
- **Estimated Size**: 400 lines

#### E1.3.2: Basic Push Implementation
- **Description**: Push built image to Gitea
- **Requirements**:
  ```go
  // pkg/build/buildah/push.go
  func (c *Client) Push(ctx context.Context, imageRef string) error {
    // Configure auth
    c.systemContext.DockerAuthConfig = &types.DockerAuthConfig{
      Username: "giteaAdmin",
      Password: getGiteaPassword(),
    }
    
    // Push using buildah
    _, _, err := buildah.Push(ctx, imageRef, nil, buildah.PushOptions{
      SystemContext: c.systemContext,
      Store: c.store,
    })
    return err
  }
  ```
- **Test Requirements**: Push to test registry
- **Estimated Size**: 250 lines

#### E1.3.3: BuildAndPush Orchestrator
- **Description**: Combine build and push into single operation
- **Requirements**:
  ```go
  // pkg/build/service.go
  func (s *BuildService) BuildAndPush(ctx context.Context, req BuildRequest) (*BuildResponse, error) {
    // 1. Format image tag
    fullTag := fmt.Sprintf("gitea.cnoe.localtest.me/giteaadmin/%s:%s", 
                          req.ImageName, req.ImageTag)
    
    // 2. Build
    imageID, err := s.client.Build(ctx, req.DockerfilePath, req.ContextDir)
    
    // 3. Tag
    err = s.client.Tag(imageID, fullTag)
    
    // 4. Push
    err = s.client.Push(ctx, fullTag)
    
    return &BuildResponse{
      ImageID: imageID,
      FullTag: fullTag,
      Success: true,
    }, nil
  }
  ```
- **Test Requirements**: End-to-end build and push test
- **Estimated Size**: 300 lines

## Phase 2: CLI Integration
**Goal**: Make the MVP usable through idpbuilder CLI
**Duration**: 1 week
**Success Criteria**:
- CLI command works
- Clear output and error messages
- Can specify basic options

### Wave 1: CLI Command Structure
**Dependencies**: Phase 1 complete
**Can Parallelize**: No

#### E2.1.1: Build Command Definition
- **Description**: Define build command structure
- **Requirements**:
  ```go
  // pkg/cmd/build/root.go
  var BuildCmd = &cobra.Command{
    Use:   "build [context]",
    Short: "Build and push container image to Gitea registry",
    Args:  cobra.ExactArgs(1),
    RunE:  runBuild,
  }
  
  func init() {
    BuildCmd.Flags().StringP("file", "f", "Dockerfile", "Dockerfile path")
    BuildCmd.Flags().StringP("tag", "t", "", "Image name:tag (required)")
    BuildCmd.MarkFlagRequired("tag")
  }
  ```
- **Test Requirements**: Command parsing tests
- **Estimated Size**: 200 lines

### Wave 2: CLI Implementation
**Dependencies**: Wave 1 complete
**Can Parallelize**: No

#### E2.2.1: Build Command Implementation
- **Description**: Implement the build command logic
- **Requirements**:
  ```go
  // pkg/cmd/build/execute.go
  func runBuild(cmd *cobra.Command, args []string) error {
    contextDir := args[0]
    dockerfilePath, _ := cmd.Flags().GetString("file")
    tag, _ := cmd.Flags().GetString("tag")
    
    // Parse image name and tag
    parts := strings.Split(tag, ":")
    imageName := parts[0]
    imageTag := "latest"
    if len(parts) > 1 {
      imageTag = parts[1]
    }
    
    // Initialize build service
    service := build.NewBuildService()
    
    // Execute build
    fmt.Printf("Building %s:%s...\n", imageName, imageTag)
    resp, err := service.BuildAndPush(ctx, build.BuildRequest{
      DockerfilePath: dockerfilePath,
      ContextDir:     contextDir,
      ImageName:      imageName,
      ImageTag:       imageTag,
    })
    
    if err != nil {
      return fmt.Errorf("build failed: %w", err)
    }
    
    fmt.Printf("Successfully built and pushed: %s\n", resp.FullTag)
    return nil
  }
  ```
- **Test Requirements**: CLI execution tests
- **Estimated Size**: 300 lines

#### E2.2.2: Error Handling and Output
- **Description**: Improve error messages and progress output
- **Requirements**:
  ```go
  // pkg/cmd/build/output.go
  func printBuildProgress(stage string) {
    fmt.Printf("=> %s\n", stage)
  }
  
  func handleBuildError(err error) error {
    if strings.Contains(err.Error(), "authentication") {
      return fmt.Errorf("authentication failed - ensure idpbuilder cluster is running")
    }
    if strings.Contains(err.Error(), "not found") {
      return fmt.Errorf("Dockerfile not found - use -f to specify path")
    }
    return err
  }
  ```
- **Test Requirements**: Error message tests
- **Estimated Size**: 200 lines

### Wave 3: Integration with Main CLI
**Dependencies**: Wave 2 complete
**Can Parallelize**: No

#### E2.3.1: Add Build Command to Root
- **Description**: Integrate build command into idpbuilder CLI
- **Requirements**:
  ```go
  // pkg/cmd/root.go (modification)
  func init() {
    // ... existing commands
    rootCmd.AddCommand(build.BuildCmd)
  }
  ```
- **Test Requirements**: Integration tests
- **Estimated Size**: 50 lines

## Phase 3: Production Ready
**Goal**: Add tests, documentation, and production quality
**Duration**: 1.5 weeks
**Success Criteria**:
- Comprehensive test coverage
- Clear documentation
- Handles edge cases

### Wave 1: Unit Tests
**Dependencies**: Phase 2 complete
**Can Parallelize**: Yes

#### E3.1.1: Build Service Tests
- **Description**: Unit tests for build service
- **Requirements**:
  ```go
  // pkg/build/service_test.go
  func TestBuildAndPush(t *testing.T) {
    // Test successful build
    // Test missing Dockerfile
    // Test invalid context
    // Test push failure
  }
  ```
- **Test Requirements**: 80% coverage
- **Estimated Size**: 400 lines

#### E3.1.2: Buildah Client Tests
- **Description**: Tests for Buildah wrapper
- **Requirements**:
  ```go
  // pkg/build/buildah/client_test.go
  func TestBuildahClient(t *testing.T) {
    // Test client initialization
    // Test build with simple Dockerfile
    // Test push with auth
    // Test certificate handling
  }
  ```
- **Test Requirements**: Mock Buildah calls
- **Estimated Size**: 350 lines

### Wave 2: Integration Tests
**Dependencies**: Wave 1 complete
**Can Parallelize**: No

#### E3.2.1: End-to-End Tests
- **Description**: Full workflow tests
- **Requirements**:
  ```go
  // tests/e2e/build_test.go
  func TestE2EBuildAndPush(t *testing.T) {
    // Start test cluster
    // Build test image
    // Verify in Gitea registry
    // Pull and run image
  }
  ```
- **Test Requirements**: Real cluster tests
- **Estimated Size**: 500 lines

### Wave 3: Documentation
**Dependencies**: Wave 2 complete
**Can Parallelize**: No

#### E3.3.1: User Documentation
- **Description**: User-facing documentation
- **Requirements**:
  ```markdown
  # docs/build.md
  ## Building Container Images
  
  ### Quick Start
  idpbuilder build . -t myapp:v1.0
  
  ### Prerequisites
  - idpbuilder cluster running
  - Dockerfile in build context
  
  ### Examples
  ...
  ```
- **Test Requirements**: Documentation review
- **Estimated Size**: 300 lines

#### E3.3.2: Example Applications
- **Description**: Sample apps with Dockerfiles
- **Requirements**:
  - Simple Go app
  - Node.js app
  - Static website
- **Test Requirements**: Examples must build
- **Estimated Size**: 200 lines

## Phase 4: Enhanced Features
**Goal**: Improve usability and performance
**Duration**: 1 week
**Success Criteria**:
- Better error handling
- Build caching
- Progress reporting

### Wave 1: Error Handling
**Dependencies**: Phase 3 complete
**Can Parallelize**: Yes

#### E4.1.1: Detailed Error Messages
- **Description**: Improve error reporting
- **Requirements**:
  ```go
  // pkg/build/errors/handler.go
  type BuildError struct {
    Phase   string
    Message string
    Suggestion string
  }
  
  func WrapError(err error, phase string) *BuildError {
    // Analyze error
    // Provide helpful suggestion
    // Format for user
  }
  ```
- **Test Requirements**: Error case tests
- **Estimated Size**: 250 lines

### Wave 2: Build Cache
**Dependencies**: Wave 1 complete
**Can Parallelize**: No

#### E4.2.1: Layer Caching
- **Description**: Implement build cache
- **Requirements**:
  ```go
  // pkg/build/cache/manager.go
  type CacheManager struct {
    cacheDir string
  }
  
  func (c *CacheManager) GetCachedLayer(instruction string) (string, bool) {
    // Check if layer exists
    // Return cached layer ID
  }
  ```
- **Test Requirements**: Cache hit/miss tests
- **Estimated Size**: 300 lines

### Wave 3: Progress Reporting
**Dependencies**: Wave 2 complete
**Can Parallelize**: No

#### E4.3.1: Build Progress Output
- **Description**: Show build progress to user
- **Requirements**:
  ```go
  // pkg/build/progress/reporter.go
  type ProgressReporter struct {
    steps []string
    current int
  }
  
  func (p *ProgressReporter) Report(step string) {
    fmt.Printf("[%d/%d] %s\n", p.current, len(p.steps), step)
  }
  ```
- **Test Requirements**: Progress output tests
- **Estimated Size**: 200 lines

## Phase 5: Advanced Features
**Goal**: Add nice-to-have features
**Duration**: 1 week
**Success Criteria**:
- Multi-stage builds work
- Build arguments supported
- Image listing works

### Wave 1: Multi-stage Builds
**Dependencies**: Phase 4 complete
**Can Parallelize**: Yes

#### E5.1.1: Multi-stage Dockerfile Support
- **Description**: Support multi-stage builds
- **Requirements**:
  ```go
  // pkg/build/multistage/parser.go
  func ParseMultistageDockerfile(path string) ([]Stage, error) {
    // Parse FROM instructions
    // Identify stage names
    // Handle COPY --from
  }
  ```
- **Test Requirements**: Multi-stage build tests
- **Estimated Size**: 400 lines

### Wave 2: Build Arguments
**Dependencies**: Wave 1 complete
**Can Parallelize**: Yes

#### E5.2.1: Build Arg Support
- **Description**: Support --build-arg flag
- **Requirements**:
  ```go
  // pkg/build/args/handler.go
  func ProcessBuildArgs(args []string) map[string]string {
    // Parse key=value pairs
    // Validate arguments
    // Pass to Buildah
  }
  ```
- **Test Requirements**: Build arg tests
- **Estimated Size**: 200 lines

### Wave 3: Registry Operations
**Dependencies**: Wave 2 complete
**Can Parallelize**: No

#### E5.3.1: List Images Command
- **Description**: List images in Gitea registry
- **Requirements**:
  ```go
  // pkg/cmd/images/list.go
  var ListCmd = &cobra.Command{
    Use:   "images",
    Short: "List images in Gitea registry",
    RunE: func(cmd *cobra.Command, args []string) error {
      // Connect to Gitea API
      // List container packages
      // Format as table
    }
  }
  ```
- **Test Requirements**: List command tests
- **Estimated Size**: 250 lines

## Dependency Matrix

### Critical Paths
```
Phase 1 Wave 1 (APIs) → Wave 2 (Libraries) → Wave 3 (Implementation)
                                                    ↓
                                            Phase 2 Wave 1 (CLI Structure)
                                                    ↓
                                            Phase 2 Wave 2 (CLI Implementation)
                                                    ↓
                                            Phase 2 Wave 3 (Integration)
                                                    ↓
                                            Phase 3 (Tests & Docs)
                                                    ↓
                                            Phase 4 (Enhancements)
                                                    ↓
                                            Phase 5 (Advanced)
```

### Parallel Opportunities
- Phase 1, Wave 1: Both efforts can run in parallel
- Phase 3, Wave 1: All test efforts can run in parallel
- Phase 4, Wave 1: Error handling parallel to other phases
- Phase 5: All waves can run somewhat in parallel

### Blocking Dependencies
- Phase 1 must complete before any CLI work
- Phase 2 must complete before production readiness
- Core MVP blocks everything else

## Risk Mitigation

### Technical Risks
| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Buildah API complexity | Medium | High | Start with minimal features |
| Certificate issues | Low | Medium | Default to InsecureSkipVerify |
| Storage conflicts | Low | Low | Use isolated storage directory |
| Auth failures | Medium | High | Test with real Gitea early |

### Process Risks
| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Scope creep | Medium | High | Strict MVP focus first |
| Integration issues | Low | Medium | Test against real cluster early |
| Performance problems | Low | Low | Profile after MVP works |

## Success Metrics

### Phase Metrics
- Phase 1: Basic build and push works
- Phase 2: CLI is usable
- Phase 3: Tests pass, docs complete
- Phase 4: Improved UX
- Phase 5: Advanced features work

### Overall Metrics
- MVP working: Phase 1-2 (2.5 weeks)
- Production ready: Phase 3 (4 weeks total)
- Full featured: Phase 5 (6-7 weeks)

## Implementation Notes

### MVP Focus (Phases 1-2)
- Absolute minimum to build and push
- Hardcode what you can (registry, namespace)
- Skip optimizations
- Basic error messages only

### Production Quality (Phase 3)
- Comprehensive testing
- Clear documentation
- Handle edge cases
- Proper error messages

### Enhancements (Phases 4-5)
- Only after MVP is solid
- Can be delivered incrementally
- User feedback driven

---

*This plan prioritizes getting a working MVP as quickly as possible, then incrementally adding quality and features.*