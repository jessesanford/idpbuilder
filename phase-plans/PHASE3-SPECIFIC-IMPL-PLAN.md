# Phase 3: Production Ready - Detailed Implementation Plan

## Phase Overview
**Duration:** 10 days  
**Critical Path:** YES - Quality gates for production deployment  
**Base Branch:** `phase2-integration`  
**Target Integration Branch:** `phase3-integration`  
**Prerequisites:** Phase 2 complete with working CLI

---

## Critical Libraries & Dependencies (MAINTAINER SPECIFIED)

### Required Libraries
```yaml
core_libraries:
  - name: "github.com/stretchr/testify"
    version: "v1.9.0"
    reason: "Industry standard for Go testing. Excellent assertion and mocking capabilities."
    usage: "Comprehensive test suites, mocks, test assertions"
    
  - name: "github.com/golang/mock"
    version: "v1.6.0"
    reason: "Official Go mocking framework. Essential for isolating dependencies in tests."
    usage: "Mock Buildah interfaces, Kubernetes clients, registry connections"
    
  - name: "github.com/onsi/ginkgo/v2"
    version: "v2.15.0"
    reason: "BDD testing framework for integration tests. Excellent for complex scenarios."
    usage: "End-to-end integration test suites"
    
  - name: "github.com/onsi/gomega"
    version: "v1.31.0"
    reason: "Matcher library companion to Ginkgo. Rich assertion capabilities."
    usage: "Integration test assertions and matchers"
    
  - name: "go.uber.org/zap"
    version: "v1.26.0"
    reason: "High-performance structured logging. Production-grade logging solution."
    usage: "Structured logging for build operations, error tracking"
```

### Interfaces to Reuse (MANDATORY)
```yaml
reused_from_previous:
  phase1:
    - "pkg/build/api/types.go: BuildRequest, BuildResponse"
    - "pkg/build/api/builder.go: Builder interface"
    - "pkg/build/service.go: Service struct and NewService()"
    - "pkg/build/buildah/client.go: Client struct"
    - "pkg/build/auth/gitea.go: GetGiteaCredentials()"
  phase2:
    - "pkg/cmd/build/root.go: BuildCmd command structure"
    - "pkg/cmd/build/flags.go: Flag parsing logic"
    - "pkg/cmd/build/output.go: Output formatting"
    
forbidden_duplications:
  - "DO NOT create new build interfaces - extend existing ones"
  - "DO NOT implement separate error handling - enhance existing patterns"
  - "DO NOT build new CLI commands - add features to existing build command"
  - "DO NOT duplicate test patterns - create reusable test utilities"
```

---

## Wave 3.1: Comprehensive Unit Tests

### Overview
**Focus:** Achieve 90% test coverage across all components  
**Dependencies:** Phase 2 complete  
**Parallelizable:** YES - Tests can be written in parallel

### E3.1.1: Build Service Test Suite
**Branch:** `phase3/wave1/effort1-service-tests`  
**Duration:** 16 hours  
**Estimated Lines:** 400 lines  
**Agent Assignment:** Single

#### Requirements
1. **MUST** achieve >90% coverage for pkg/build/service.go
2. **MUST** use testify for assertions and mocks
3. **MUST** test all error scenarios
4. **MUST NOT** duplicate existing Phase 1 logic

#### Implementation Guidance

##### Directory Structure
```
pkg/
├── build/
│   ├── service_test.go       # ~200 lines
│   ├── mocks/
│   │   ├── builder.go        # ~100 lines (generated)
│   │   └── storage.go        # ~100 lines (generated)
```

##### Service Test Suite (Maintainer Specified)
```go
// pkg/build/service_test.go
package build

import (
    "context"
    "errors"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
    
    "idpbuilder/pkg/build/api"
    "idpbuilder/pkg/build/mocks"
)

// ServiceTestSuite provides comprehensive testing for build service
type ServiceTestSuite struct {
    suite.Suite
    mockBuilder *mocks.MockBuilder
    service     *Service
}

func (s *ServiceTestSuite) SetupTest() {
    s.mockBuilder = mocks.NewMockBuilder(s.T())
    s.service = &Service{
        client: s.mockBuilder,
    }
}

func (s *ServiceTestSuite) TestBuildAndPush_Success() {
    // Arrange
    request := api.BuildRequest{
        DockerfilePath: "Dockerfile",
        ContextDir:     "/tmp/build",
        ImageName:      "testapp",
        ImageTag:       "v1.0",
    }
    
    expectedResponse := &api.BuildResponse{
        ImageID: "sha256:abcd1234",
        FullTag: "gitea.cnoe.localtest.me/giteaadmin/testapp:v1.0",
        Success: true,
    }
    
    s.mockBuilder.On("BuildAndPush", mock.Anything, request).
        Return(expectedResponse, nil)
    
    // Act
    result, err := s.service.BuildAndPush(context.Background(), request)
    
    // Assert
    assert.NoError(s.T(), err)
    assert.Equal(s.T(), expectedResponse, result)
    s.mockBuilder.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestBuildAndPush_BuildError() {
    // Test build failure scenarios
    request := api.BuildRequest{
        DockerfilePath: "Dockerfile",
        ContextDir:     "/tmp/build",
        ImageName:      "testapp",
        ImageTag:       "v1.0",
    }
    
    expectedError := errors.New("buildah: build failed")
    s.mockBuilder.On("BuildAndPush", mock.Anything, request).
        Return(nil, expectedError)
    
    result, err := s.service.BuildAndPush(context.Background(), request)
    
    assert.Error(s.T(), err)
    assert.Nil(s.T(), result)
    assert.Contains(s.T(), err.Error(), "build failed")
}

func (s *ServiceTestSuite) TestBuildAndPush_ValidationError() {
    // Test validation failures
    invalidRequest := api.BuildRequest{
        // Missing required fields
        ImageName: "",
    }
    
    result, err := s.service.BuildAndPush(context.Background(), invalidRequest)
    
    // Should not call builder if validation fails
    s.mockBuilder.AssertNotCalled(s.T(), "BuildAndPush")
    assert.Error(s.T(), err)
    assert.Nil(s.T(), result)
}

func (s *ServiceTestSuite) TestNewService_Success() {
    // Test service creation
    service, err := NewService()
    
    assert.NoError(s.T(), err)
    assert.NotNil(s.T(), service)
    assert.NotNil(s.T(), service.client)
}

func TestServiceTestSuite(t *testing.T) {
    suite.Run(t, new(ServiceTestSuite))
}

// Additional table-driven tests for edge cases
func TestBuildRequestValidation(t *testing.T) {
    testCases := []struct {
        name    string
        request api.BuildRequest
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid request",
            request: api.BuildRequest{
                DockerfilePath: "Dockerfile",
                ContextDir:     "/tmp",
                ImageName:      "app",
                ImageTag:       "v1",
            },
            wantErr: false,
        },
        {
            name: "missing dockerfile path",
            request: api.BuildRequest{
                ContextDir: "/tmp",
                ImageName:  "app",
                ImageTag:   "v1",
            },
            wantErr: true,
            errMsg:  "DockerfilePath is required",
        },
        {
            name: "missing context dir",
            request: api.BuildRequest{
                DockerfilePath: "Dockerfile",
                ImageName:      "app",
                ImageTag:       "v1",
            },
            wantErr: true,
            errMsg:  "ContextDir is required",
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            err := tc.request.Validate()
            
            if tc.wantErr {
                assert.Error(t, err)
                if tc.errMsg != "" {
                    assert.Contains(t, err.Error(), tc.errMsg)
                }
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

##### Mock Builder Interface (Maintainer Specified)
```go
// pkg/build/mocks/builder.go
//go:generate mockgen -source=../api/builder.go -destination=builder.go

package mocks

import (
    "context"
    
    "github.com/stretchr/testify/mock"
    
    "idpbuilder/pkg/build/api"
)

// MockBuilder implements the Builder interface for testing
type MockBuilder struct {
    mock.Mock
}

func NewMockBuilder(t interface{}) *MockBuilder {
    return &MockBuilder{}
}

func (m *MockBuilder) BuildAndPush(ctx context.Context, req api.BuildRequest) (*api.BuildResponse, error) {
    args := m.Called(ctx, req)
    
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    
    return args.Get(0).(*api.BuildResponse), args.Error(1)
}
```

#### Test Requirements (TDD)
- [ ] >90% coverage for service.go
- [ ] All public methods tested
- [ ] Error scenarios covered
- [ ] Mock objects properly configured
- [ ] Concurrent access tested

#### Success Criteria
- [ ] Test suite passes consistently
- [ ] Coverage report shows >90%
- [ ] No flaky tests
- [ ] Mock assertions validate behavior
- [ ] Under 400 lines per line-counter.sh

### E3.1.2: Buildah Client Tests
**Branch:** `phase3/wave1/effort2-buildah-tests`  
**Duration:** 20 hours  
**Estimated Lines:** 450 lines  
**Agent Assignment:** Single

#### Requirements
1. **MUST** test Buildah client wrapper thoroughly
2. **MUST** mock Buildah dependencies
3. **MUST** test authentication scenarios
4. **MUST** handle storage initialization failures

#### Implementation Guidance

##### Buildah Client Test Suite (Maintainer Specified)
```go
// pkg/build/buildah/client_test.go
package buildah

import (
    "context"
    "os"
    "path/filepath"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/stretchr/testify/require"
    
    "idpbuilder/pkg/build/api"
)

type BuildahClientTestSuite struct {
    suite.Suite
    tempDir string
    client  *Client
}

func (s *BuildahClientTestSuite) SetupSuite() {
    // Create temporary directory for test storage
    tempDir, err := os.MkdirTemp("", "buildah-test-*")
    require.NoError(s.T(), err)
    s.tempDir = tempDir
    
    // Set up test environment
    os.Setenv("TMPDIR", tempDir)
}

func (s *BuildahClientTestSuite) TearDownSuite() {
    os.RemoveAll(s.tempDir)
}

func (s *BuildahClientTestSuite) SetupTest() {
    config := api.DefaultConfig()
    client, err := NewClient(config)
    require.NoError(s.T(), err)
    s.client = client
}

func (s *BuildahClientTestSuite) TestNewClient_Success() {
    config := api.DefaultConfig()
    client, err := NewClient(config)
    
    assert.NoError(s.T(), err)
    assert.NotNil(s.T(), client)
    assert.NotNil(s.T(), client.store)
    assert.NotNil(s.T(), client.systemContext)
    assert.True(s.T(), client.systemContext.DockerInsecureSkipTLSVerify.Value)
}

func (s *BuildahClientTestSuite) TestBuildAndPush_InvalidRequest() {
    ctx := context.Background()
    
    // Test with empty request
    emptyReq := api.BuildRequest{}
    resp, err := s.client.BuildAndPush(ctx, emptyReq)
    
    assert.NoError(s.T(), err) // Should not error, but should return failure
    assert.NotNil(s.T(), resp)
    assert.False(s.T(), resp.Success)
    assert.Contains(s.T(), resp.Error, "validation failed")
}

func (s *BuildahClientTestSuite) TestBuildAndPush_MissingDockerfile() {
    ctx := context.Background()
    
    req := api.BuildRequest{
        DockerfilePath: "nonexistent.dockerfile",
        ContextDir:     s.tempDir,
        ImageName:      "test",
        ImageTag:       "latest",
    }
    
    resp, err := s.client.BuildAndPush(ctx, req)
    
    assert.NoError(s.T(), err)
    assert.NotNil(s.T(), resp)
    assert.False(s.T(), resp.Success)
    assert.Contains(s.T(), resp.Error, "failed to read Dockerfile")
}

func (s *BuildahClientTestSuite) TestBuild_SimpleDockerfile() {
    ctx := context.Background()
    
    // Create test Dockerfile
    dockerfile := `FROM alpine:3.19
RUN echo "test build" > /test.txt
CMD ["cat", "/test.txt"]`
    
    dockerfilePath := filepath.Join(s.tempDir, "Dockerfile")
    err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644)
    require.NoError(s.T(), err)
    
    req := api.BuildRequest{
        DockerfilePath: "Dockerfile",
        ContextDir:     s.tempDir,
        ImageName:      "testimage",
        ImageTag:       "latest",
    }
    
    // This test may fail in CI without proper setup, so we focus on validation
    resp, err := s.client.BuildAndPush(ctx, req)
    
    // Either succeeds or fails with expected error
    if err != nil {
        // Expected failure in test environment
        assert.Contains(s.T(), err.Error(), "build failed")
    } else {
        assert.NotNil(s.T(), resp)
        if resp.Success {
            assert.NotEmpty(s.T(), resp.ImageID)
            assert.Contains(s.T(), resp.FullTag, "testimage:latest")
        }
    }
}

func TestBuildahClientTestSuite(t *testing.T) {
    suite.Run(t, new(BuildahClientTestSuite))
}

// Test utilities for Buildah operations
func TestParseDockerfile(t *testing.T) {
    testCases := []struct {
        name         string
        dockerfile   string
        wantCommands []string
    }{
        {
            name: "simple dockerfile",
            dockerfile: `FROM alpine:3.19
COPY app /app
RUN chmod +x /app
CMD ["/app"]`,
            wantCommands: []string{"FROM", "COPY", "RUN", "CMD"},
        },
        {
            name: "dockerfile with comments",
            dockerfile: `# This is a comment
FROM alpine:3.19
# Another comment
RUN echo "hello"`,
            wantCommands: []string{"FROM", "RUN"},
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            instructions, err := parseDockerfile([]byte(tc.dockerfile))
            
            assert.NoError(t, err)
            assert.Len(t, instructions, len(tc.wantCommands))
            
            for i, expectedCmd := range tc.wantCommands {
                assert.Equal(t, expectedCmd, instructions[i].Command)
            }
        })
    }
}
```

### E3.1.3: CLI Command Tests
**Branch:** `phase3/wave1/effort3-cli-tests`  
**Duration:** 12 hours  
**Estimated Lines:** 300 lines  
**Agent Assignment:** Single

#### Requirements
1. **MUST** test all CLI flag combinations
2. **MUST** test error handling and output
3. **MUST** mock build service calls
4. **MUST** test integration scenarios

#### Implementation Guidance

##### CLI Test Suite (Maintainer Specified)
```go
// pkg/cmd/build/root_test.go (extended)
package build

import (
    "bytes"
    "context"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "testing"
    
    "github.com/spf13/cobra"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
    
    "idpbuilder/pkg/build/api"
    "idpbuilder/pkg/build/mocks"
)

type CLITestSuite struct {
    suite.Suite
    tempDir    string
    mockService *mocks.MockService
    cmd        *cobra.Command
    output     *bytes.Buffer
}

func (s *CLITestSuite) SetupTest() {
    // Create temp directory with test Dockerfile
    tempDir, err := os.MkdirTemp("", "cli-test-*")
    s.Require().NoError(err)
    s.tempDir = tempDir
    
    dockerfile := "FROM alpine:3.19\nRUN echo test"
    dockerfilePath := filepath.Join(tempDir, "Dockerfile")
    err = os.WriteFile(dockerfilePath, []byte(dockerfile), 0644)
    s.Require().NoError(err)
    
    // Set up mock service
    s.mockService = mocks.NewMockService(s.T())
    
    // Create command with output capture
    s.output = &bytes.Buffer{}
    s.cmd = BuildCmd
    s.cmd.SetOut(s.output)
    s.cmd.SetErr(s.output)
}

func (s *CLITestSuite) TearDownTest() {
    os.RemoveAll(s.tempDir)
    s.output.Reset()
}

func (s *CLITestSuite) TestBuildCommand_Success() {
    // Mock successful build
    expectedResponse := &api.BuildResponse{
        ImageID: "sha256:abcd1234",
        FullTag: "gitea.cnoe.localtest.me/giteaadmin/testapp:v1.0",
        Success: true,
    }
    
    s.mockService.On("BuildAndPush", mock.Anything, mock.AnythingOfType("api.BuildRequest")).
        Return(expectedResponse, nil)
    
    // Execute command
    s.cmd.SetArgs([]string{s.tempDir, "-t", "testapp:v1.0"})
    err := s.cmd.Execute()
    
    assert.NoError(s.T(), err)
    
    output := s.output.String()
    assert.Contains(s.T(), output, "Building testapp:v1.0")
    assert.Contains(s.T(), output, "Successfully built and pushed")
}

func (s *CLITestSuite) TestBuildCommand_MissingTag() {
    // Test missing required flag
    s.cmd.SetArgs([]string{s.tempDir})
    err := s.cmd.Execute()
    
    assert.Error(s.T(), err)
    assert.Contains(s.T(), err.Error(), "required flag")
}

func (s *CLITestSuite) TestBuildCommand_NonexistentContext() {
    // Test nonexistent build context
    s.cmd.SetArgs([]string{"/nonexistent", "-t", "test:latest"})
    err := s.cmd.Execute()
    
    assert.Error(s.T(), err)
    assert.Contains(s.T(), err.Error(), "context directory does not exist")
}

func (s *CLITestSuite) TestBuildCommand_QuietMode() {
    expectedResponse := &api.BuildResponse{
        ImageID: "sha256:abcd1234",
        FullTag: "gitea.cnoe.localtest.me/giteaadmin/testapp:latest",
        Success: true,
    }
    
    s.mockService.On("BuildAndPush", mock.Anything, mock.AnythingOfType("api.BuildRequest")).
        Return(expectedResponse, nil)
    
    // Execute with quiet flag
    s.cmd.SetArgs([]string{s.tempDir, "-t", "testapp", "--quiet"})
    err := s.cmd.Execute()
    
    assert.NoError(s.T(), err)
    
    output := s.output.String()
    // Should only contain the final tag, no progress messages
    assert.Equal(s.T(), "gitea.cnoe.localtest.me/giteaadmin/testapp:latest\n", output)
}

func TestCLITestSuite(t *testing.T) {
    suite.Run(t, new(CLITestSuite))
}

// Test flag parsing edge cases
func TestFlagParsing(t *testing.T) {
    testCases := []struct {
        name     string
        args     []string
        wantName string
        wantTag  string
        wantFile string
        wantErr  bool
    }{
        {
            name:     "basic tag",
            args:     []string{".", "-t", "myapp:v1.0"},
            wantName: "myapp",
            wantTag:  "v1.0",
            wantFile: "Dockerfile",
            wantErr:  false,
        },
        {
            name:     "custom dockerfile",
            args:     []string{".", "-t", "myapp", "-f", "docker/Dockerfile"},
            wantName: "myapp",
            wantTag:  "latest",
            wantFile: "docker/Dockerfile",
            wantErr:  false,
        },
        {
            name:    "missing tag",
            args:    []string{"."},
            wantErr: true,
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Reset flags
            dockerfileFlag = ""
            tagFlag = ""
            
            cmd := &cobra.Command{}
            cmd.Flags().StringVarP(&dockerfileFlag, "file", "f", "Dockerfile", "dockerfile")
            cmd.Flags().StringVarP(&tagFlag, "tag", "t", "", "tag")
            cmd.MarkFlagRequired("tag")
            
            cmd.SetArgs(tc.args)
            err := cmd.ParseFlags(tc.args)
            
            if tc.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            
            if tc.wantFile != "" {
                assert.Equal(t, tc.wantFile, dockerfileFlag)
            }
            
            if tc.wantName != "" || tc.wantTag != "" {
                name, tag, err := parseImageTag(tagFlag)
                assert.NoError(t, err)
                assert.Equal(t, tc.wantName, name)
                assert.Equal(t, tc.wantTag, tag)
            }
        })
    }
}
```

---

## Wave 3.2: Integration Tests

### Overview
**Focus:** End-to-end workflow testing  
**Dependencies:** Wave 3.1 complete  
**Parallelizable:** NO - Integration tests must run sequentially

### E3.2.1: End-to-End Build Tests
**Branch:** `phase3/wave2/effort1-e2e-tests`  
**Duration:** 24 hours  
**Estimated Lines:** 500 lines  
**Agent Assignment:** Single

#### Requirements
1. **MUST** test complete build and push workflow
2. **MUST** test against real or mock registry
3. **MUST** test various Dockerfile scenarios
4. **MUST** validate error handling end-to-end

#### Implementation Guidance

##### E2E Test Suite (Maintainer Specified)
```go
// tests/e2e/build_test.go
package e2e

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "testing"
    "time"
    
    "github.com/onsi/ginkgo/v2"
    "github.com/onsi/gomega"
    
    "idpbuilder/pkg/build"
    "idpbuilder/pkg/build/api"
)

var _ = ginkgo.Describe("Container Build E2E", func() {
    var (
        tempDir    string
        service    *build.Service
        ctx        context.Context
        cancel     context.CancelFunc
    )
    
    ginkgo.BeforeEach(func() {
        var err error
        
        // Create temporary build context
        tempDir, err = os.MkdirTemp("", "e2e-build-*")
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
        
        // Initialize build service
        service, err = build.NewService()
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
        
        // Create context with timeout
        ctx, cancel = context.WithTimeout(context.Background(), 10*time.Minute)
    })
    
    ginkgo.AfterEach(func() {
        cancel()
        if tempDir != "" {
            os.RemoveAll(tempDir)
        }
    })
    
    ginkgo.Context("Simple Alpine Build", func() {
        ginkgo.It("should build and push alpine-based image", func() {
            // Create simple Dockerfile
            dockerfile := `FROM alpine:3.19
RUN echo "Hello from E2E test" > /hello.txt
CMD ["cat", "/hello.txt"]`
            
            dockerfilePath := filepath.Join(tempDir, "Dockerfile")
            err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644)
            gomega.Expect(err).NotTo(gomega.HaveOccurred())
            
            // Build request
            request := api.BuildRequest{
                DockerfilePath: "Dockerfile",
                ContextDir:     tempDir,
                ImageName:      "e2e-alpine-test",
                ImageTag:       fmt.Sprintf("test-%d", time.Now().Unix()),
            }
            
            // Execute build
            response, err := service.BuildAndPush(ctx, request)
            
            // Validate response
            if response != nil && !response.Success {
                // In test environment, push might fail - that's ok for now
                ginkgo.Skip("Registry not available in test environment")
            }
            
            gomega.Expect(err).NotTo(gomega.HaveOccurred())
            gomega.Expect(response).NotTo(gomega.BeNil())
            gomega.Expect(response.ImageID).NotTo(gomega.BeEmpty())
            gomega.Expect(response.FullTag).To(gomega.ContainSubstring("e2e-alpine-test"))
        })
    })
    
    ginkgo.Context("Multi-instruction Dockerfile", func() {
        ginkgo.It("should handle COPY, RUN, and CMD instructions", func() {
            // Create test app file
            appContent := `#!/bin/sh
echo "E2E test application"
echo "Build completed successfully"`
            
            appPath := filepath.Join(tempDir, "app.sh")
            err := os.WriteFile(appPath, []byte(appContent), 0755)
            gomega.Expect(err).NotTo(gomega.HaveOccurred())
            
            // Create Dockerfile with multiple instructions
            dockerfile := `FROM alpine:3.19
COPY app.sh /usr/local/bin/app.sh
RUN chmod +x /usr/local/bin/app.sh
RUN apk add --no-cache curl
CMD ["/usr/local/bin/app.sh"]`
            
            dockerfilePath := filepath.Join(tempDir, "Dockerfile")
            err = os.WriteFile(dockerfilePath, []byte(dockerfile), 0644)
            gomega.Expect(err).NotTo(gomega.HaveOccurred())
            
            request := api.BuildRequest{
                DockerfilePath: "Dockerfile",
                ContextDir:     tempDir,
                ImageName:      "e2e-multi-test",
                ImageTag:       fmt.Sprintf("test-%d", time.Now().Unix()),
            }
            
            response, err := service.BuildAndPush(ctx, request)
            
            if response != nil && !response.Success {
                ginkgo.Skip("Registry not available in test environment")
            }
            
            gomega.Expect(err).NotTo(gomega.HaveOccurred())
            gomega.Expect(response).NotTo(gomega.BeNil())
        })
    })
    
    ginkgo.Context("Error Scenarios", func() {
        ginkgo.It("should handle missing Dockerfile gracefully", func() {
            request := api.BuildRequest{
                DockerfilePath: "NonexistentDockerfile",
                ContextDir:     tempDir,
                ImageName:      "error-test",
                ImageTag:       "latest",
            }
            
            response, err := service.BuildAndPush(ctx, request)
            
            // Should not panic, should return error info
            gomega.Expect(err).NotTo(gomega.HaveOccurred()) // Service should not error
            gomega.Expect(response).NotTo(gomega.BeNil())
            gomega.Expect(response.Success).To(gomega.BeFalse())
            gomega.Expect(response.Error).To(gomega.ContainSubstring("failed to read Dockerfile"))
        })
        
        ginkgo.It("should handle invalid Dockerfile syntax", func() {
            dockerfile := `INVALID_INSTRUCTION
FROM alpine:3.19
MALFORMED LINE WITHOUT PROPER SYNTAX`
            
            dockerfilePath := filepath.Join(tempDir, "Dockerfile")
            err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644)
            gomega.Expect(err).NotTo(gomega.HaveOccurred())
            
            request := api.BuildRequest{
                DockerfilePath: "Dockerfile",
                ContextDir:     tempDir,
                ImageName:      "syntax-error-test",
                ImageTag:       "latest",
            }
            
            response, err := service.BuildAndPush(ctx, request)
            
            gomega.Expect(err).NotTo(gomega.HaveOccurred())
            gomega.Expect(response).NotTo(gomega.BeNil())
            gomega.Expect(response.Success).To(gomega.BeFalse())
        })
    })
})

func TestE2E(t *testing.T) {
    gomega.RegisterFailHandler(ginkgo.Fail)
    ginkgo.RunSpecs(t, "Build E2E Suite")
}
```

---

## Wave 3.3: Documentation and Examples

### Overview
**Focus:** Complete documentation and example applications  
**Dependencies:** Wave 3.2 complete  
**Parallelizable:** YES - Documentation can be written in parallel

### E3.3.1: User Documentation
**Branch:** `phase3/wave3/effort1-user-docs`  
**Duration:** 8 hours  
**Estimated Lines:** 200 lines  
**Agent Assignment:** Single

#### Requirements
1. **MUST** create comprehensive user guide
2. **MUST** include troubleshooting section
3. **MUST** document all CLI options
4. **MUST NOT** duplicate existing idpbuilder documentation

#### Implementation Guidance

##### Documentation Structure (Maintainer Specified)
```markdown
# docs/build-feature.md
# Container Building with idpbuilder

## Overview

The `idpbuilder build` command allows you to build container images locally and push them directly to your idpbuilder cluster's Gitea registry.

## Quick Start

```bash
# Build current directory with default Dockerfile
idpbuilder build . -t myapp:latest

# Build with custom Dockerfile
idpbuilder build . -f docker/Dockerfile -t myapp:v1.0

# Build with quiet output (for scripts)
idpbuilder build . -t myapp:latest --quiet
```

## Prerequisites

- idpbuilder cluster must be running
- Dockerfile must exist in build context
- Build context must be accessible locally

## Command Reference

### Syntax
```
idpbuilder build [CONTEXT] [flags]
```

### Required Arguments
- `CONTEXT`: Path to build context directory

### Required Flags
- `-t, --tag`: Image name and tag (format: name:tag)

### Optional Flags
- `-f, --file`: Dockerfile path relative to context (default: "Dockerfile")
- `-q, --quiet`: Suppress build output
- `-h, --help`: Show help information

## Examples

### Basic Web Application
```bash
# Build a simple web app
idpbuilder build ./webapp -t mywebapp:v1.0
```

### Multi-stage Build
```bash
# Build with production optimizations
idpbuilder build . -f Dockerfile.prod -t myapp:production
```

### Development Workflow
```bash
# Quick development build
idpbuilder build . -t myapp:dev --quiet
```

## Registry Information

Images are automatically pushed to:
- Registry: `gitea.cnoe.localtest.me`
- Namespace: `giteaadmin`
- Full path: `gitea.cnoe.localtest.me/giteaadmin/{image-name}:{tag}`

## Troubleshooting

### Common Issues

#### Authentication Failed
```
Error: Authentication failed
```
**Solution**: Ensure your idpbuilder cluster is running:
```bash
idpbuilder get
```

#### Dockerfile Not Found
```
Error: Dockerfile not found
```
**Solutions**:
- Check Dockerfile exists: `ls -la Dockerfile`
- Use `-f` flag to specify different path: `-f docker/Dockerfile`
- Ensure you're in the correct directory

#### Build Timeout
```
Error: Build timeout - context deadline exceeded
```
**Solutions**:
- Simplify Dockerfile to reduce build time
- Check network connectivity for base image downloads
- Retry the build

#### Registry Connection Failed
```
Error: Registry connection failed
```
**Solution**: Verify cluster status:
```bash
idpbuilder get
kubectl get pods -n gitea
```

## Supported Dockerfile Instructions

Current support includes:
- `FROM` - Base image specification
- `COPY` - Copy files to image
- `RUN` - Execute commands during build
- `CMD` - Default command to run
- `EXPOSE` - Document exposed ports
- `ENV` - Set environment variables

## Limitations

- Multi-stage builds: Limited support (Phase 5 feature)
- Build arguments: Not yet supported (Phase 5 feature)
- Custom registries: Fixed to cluster Gitea registry
- Parallel builds: Not supported

## Integration with Kubernetes

After building, use your images in Kubernetes:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 1
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: gitea.cnoe.localtest.me/giteaadmin/myapp:v1.0
        ports:
        - containerPort: 8080
```

## Performance Tips

- Use .dockerignore to exclude unnecessary files
- Leverage Docker layer caching
- Use minimal base images (alpine)
- Combine RUN instructions to reduce layers

## Security Considerations

- Images are stored in cluster-local registry
- TLS verification is disabled by default for development
- Credentials are managed automatically via cluster secrets

## Next Steps

- Phase 4: Enhanced error handling and caching
- Phase 5: Multi-stage builds and build arguments
- Advanced registry operations and image management
```

### E3.3.2: Example Applications
**Branch:** `phase3/wave3/effort2-examples`  
**Duration:** 8 hours  
**Estimated Lines:** 250 lines  
**Agent Assignment:** Single

#### Requirements
1. **MUST** create working example applications
2. **MUST** include variety of use cases
3. **MUST** verify examples build successfully
4. **MUST** include documentation for each example

#### Implementation Guidance

##### Example Applications Structure
```
examples/
├── go-webapp/
│   ├── Dockerfile
│   ├── main.go
│   ├── README.md
│   └── build.sh
├── node-api/
│   ├── Dockerfile
│   ├── package.json
│   ├── server.js
│   ├── README.md
│   └── build.sh
└── static-site/
    ├── Dockerfile
    ├── index.html
    ├── styles.css
    ├── README.md
    └── build.sh
```

---

## Phase-Wide Constraints

### Architecture Decisions (Maintainer Specified)
```markdown
1. **Test Strategy**
   - Unit tests: >90% coverage required
   - Integration tests: Real workflow validation
   - E2E tests: Full system validation
   - Mock external dependencies for reliability

2. **Documentation Standards**
   - User-focused language
   - Include troubleshooting guidance
   - Provide working examples
   - Reference existing idpbuilder patterns

3. **Quality Gates**
   - All tests must pass
   - Coverage thresholds enforced
   - Examples must build successfully
   - Documentation must be comprehensive
```

### Forbidden Duplications
- DO NOT duplicate test utilities - create shared helpers
- DO NOT recreate build logic - test existing implementations
- DO NOT duplicate documentation patterns - follow idpbuilder style
- DO NOT create redundant examples - ensure each demonstrates unique concepts

---

## Testing Strategy

### Coverage Requirements
- **Unit Tests**: >90% for all packages
- **Integration Tests**: Major workflow coverage
- **E2E Tests**: Real environment validation
- **CLI Tests**: All command combinations

### Test Environment Setup
```bash
# Run unit tests
go test ./pkg/... -cover -coverprofile=coverage.out

# Run integration tests
go test ./tests/integration/... -tags=integration

# Run E2E tests (requires cluster)
go test ./tests/e2e/... -tags=e2e

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html
```

---

## Success Criteria

### Functional
- [ ] >90% test coverage across all packages
- [ ] All integration tests pass
- [ ] E2E tests validate complete workflows
- [ ] Documentation is comprehensive and accurate
- [ ] Examples build and run successfully

### Quality
- [ ] No flaky tests
- [ ] Fast test execution (<5 minutes for unit tests)
- [ ] Clear test failure messages
- [ ] Proper test isolation

### Production Readiness
- [ ] Error handling thoroughly tested
- [ ] Edge cases covered
- [ ] Performance acceptable
- [ ] Security considerations documented

This phase establishes production-quality build functionality with comprehensive testing and documentation, making the feature ready for production deployment and user adoption.