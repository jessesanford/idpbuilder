# IDPBUILDER PUSH COMMAND - COMPREHENSIVE REQUIREMENTS ANALYSIS

## Executive Summary
This document provides a comprehensive requirements analysis for implementing a `push` command for the idpbuilder CLI tool. The command will enable pushing OCI images to the Gitea registry at https://gitea.cnoe.localtest.me:8443/ using the go-containerregistry library with TEST-DRIVEN DEVELOPMENT methodology.

## Project Overview

### Core Requirements
- **Command**: `idpbuilder push [image] [registry]`
- **Authentication**: Support `--username` and `--password` flags
- **Security**: Support `--insecure` flag for self-signed certificates
- **Library**: Use google/go-containerregistry for OCI operations
- **Methodology**: TEST-DRIVEN DEVELOPMENT (TDD) - Red-Green-Refactor cycle
- **Integration**: Seamless integration with existing idpbuilder codebase
- **Standards**: Use standard OCI image and transport mechanisms only

### Success Criteria
1. Tests written BEFORE implementation (RED phase)
2. Minimal code to pass tests (GREEN phase)
3. Refactored for quality (REFACTOR phase)
4. 80%+ test coverage for all new code
5. Zero regression in existing functionality
6. Full integration with existing get secrets command for defaults
7. Compliant with OCI standards (no Docker/Gitea specific code)

## Architecture Analysis

### Existing idpbuilder Architecture

#### Command Structure (Cobra-based)
```
idpbuilder/
├── cmd/
│   ├── root.go          # Main command entry
│   ├── create/          # Create command
│   ├── delete/          # Delete command
│   ├── get/             # Get command group
│   │   ├── clusters.go
│   │   ├── packages.go
│   │   └── secrets.go   # Can retrieve gitea credentials
│   └── version/         # Version command
```

#### Key Integration Points
1. **Command Framework**: Uses Cobra for CLI structure
2. **Secret Management**: Existing `get secrets` retrieves gitea credentials
   - Secret name: `gitea-credential`
   - Contains: username, token fields
3. **Logging**: Uses logr for structured logging
4. **Testing**: Uses stretchr/testify for assertions

### Proposed Architecture for Push Command

#### Component Design
```
idpbuilder/
├── cmd/
│   └── push/
│       ├── root.go           # Main push command
│       └── root_test.go      # Command-level tests
├── pkg/
│   └── oci/
│       ├── client.go         # OCI registry client interface
│       ├── client_test.go    # Client unit tests
│       ├── auth.go           # Authentication handling
│       ├── auth_test.go      # Auth unit tests
│       ├── push.go           # Push implementation
│       ├── push_test.go      # Push unit tests
│       └── transport.go      # Transport configuration
│       └── transport_test.go # Transport tests
```

#### Core Interfaces
```go
// RegistryClient defines the contract for OCI registry operations
type RegistryClient interface {
    Push(ctx context.Context, image, destination string, opts ...Option) error
    Authenticate(username, password string) error
    SetInsecure(insecure bool)
}

// Authenticator handles credential management
type Authenticator interface {
    GetCredentials(ctx context.Context) (*Credentials, error)
    ValidateCredentials(username, password string) error
}

// ImagePusher handles the actual push operation
type ImagePusher interface {
    ValidateImage(imagePath string) error
    PushImage(ctx context.Context, src, dst string, auth Authenticator) error
}
```

## Test Strategy (TDD Approach)

### Test Pyramid
```
        E2E Tests (10%)
       /            \
      /              \
   Integration Tests (30%)
    /                  \
   /                    \
Unit Tests (60%)
```

### Test Phases by Development Cycle

#### Phase 1: Foundation (RED-GREEN-REFACTOR)
1. **RED**: Write failing tests for command structure
   - Test command registration
   - Test flag parsing
   - Test help text
2. **GREEN**: Implement minimal command skeleton
3. **REFACTOR**: Clean up command structure

#### Phase 2: Authentication (RED-GREEN-REFACTOR)
1. **RED**: Write failing tests for authentication
   - Test credential retrieval from secrets
   - Test flag override behavior
   - Test validation logic
2. **GREEN**: Implement authentication
3. **REFACTOR**: Extract auth interfaces

#### Phase 3: OCI Operations (RED-GREEN-REFACTOR)
1. **RED**: Write failing tests for OCI operations
   - Test image validation
   - Test registry connection
   - Test push operation
2. **GREEN**: Implement OCI client
3. **REFACTOR**: Optimize and clean

### Test Frameworks and Tools
- **Unit Testing**: Go standard testing package
- **Assertions**: stretchr/testify
- **Mocking**: gomock for interface mocking
- **Integration**: Docker test containers
- **Coverage**: go test -cover (target 80%+)

## Implementation Plan Structure

### Phase 1: Foundation & Command Structure
**Estimated LOC**: 800 lines
**Duration**: Wave 1 (2-3 days)

#### Wave 1.1: Command Skeleton (TDD)
- **Effort 1**: Write command tests (150 lines)
  - Test command registration
  - Test flag definitions
  - Test help/usage
- **Effort 2**: Implement command (200 lines)
  - Create push command structure
  - Add to root command
  - Define flags
- **Effort 3**: Integration tests (150 lines)
  - Test command in CLI context
  - Verify flag behavior

#### Wave 1.2: Validation & Error Handling (TDD)
- **Effort 1**: Write validation tests (150 lines)
  - Test input validation
  - Test error scenarios
- **Effort 2**: Implement validation (150 lines)
  - Input sanitization
  - Error handling

### Phase 2: Authentication & Credentials
**Estimated LOC**: 1200 lines
**Duration**: Waves 2-3 (3-4 days)

#### Wave 2.1: Credential Management (TDD)
- **Effort 1**: Write auth tests (200 lines)
  - Test secret retrieval
  - Test credential validation
- **Effort 2**: Implement auth module (300 lines)
  - Secret integration
  - Credential management
- **Effort 3**: Mock auth for testing (200 lines)
  - Create test doubles
  - Integration test setup

#### Wave 2.2: Authentication Flow (TDD)
- **Effort 1**: Write flow tests (200 lines)
  - Test auth precedence
  - Test fallback behavior
- **Effort 2**: Implement auth flow (300 lines)
  - Flag override logic
  - Default credential handling

### Phase 3: OCI Registry Integration
**Estimated LOC**: 1500 lines
**Duration**: Waves 4-5 (4-5 days)

#### Wave 3.1: OCI Client Implementation (TDD)
- **Effort 1**: Write client tests (300 lines)
  - Test registry connection
  - Test transport setup
- **Effort 2**: Implement OCI client (400 lines)
  - go-containerregistry integration
  - Transport configuration
- **Effort 3**: Insecure mode handling (200 lines)
  - TLS configuration
  - Certificate validation bypass

#### Wave 3.2: Push Operation (TDD)
- **Effort 1**: Write push tests (200 lines)
  - Test image validation
  - Test push scenarios
- **Effort 2**: Implement push (400 lines)
  - Image handling
  - Push execution
  - Progress reporting

### Phase 4: Integration & Polish
**Estimated LOC**: 700 lines
**Duration**: Wave 6 (2-3 days)

#### Wave 4.1: E2E Testing & Documentation
- **Effort 1**: E2E test suite (400 lines)
  - Full workflow tests
  - Error recovery tests
- **Effort 2**: Documentation (300 lines)
  - Code documentation
  - User documentation
  - Examples

## Library Dependencies

### Core Dependencies
```go
// go.mod additions
require (
    github.com/google/go-containerregistry v0.20.2
    github.com/stretchr/testify v1.9.0  // existing
    github.com/golang/mock v1.6.0       // for mocking
)
```

### go-containerregistry Usage Pattern
```go
import (
    "github.com/google/go-containerregistry/pkg/authn"
    "github.com/google/go-containerregistry/pkg/name"
    "github.com/google/go-containerregistry/pkg/v1/remote"
    "github.com/google/go-containerregistry/pkg/v1/remote/transport"
)

// Example usage pattern
func (c *ociClient) Push(image, dest string) error {
    ref, err := name.ParseReference(dest)
    img, err := loadImage(image)

    opts := []remote.Option{
        remote.WithAuth(&authn.Basic{
            Username: c.username,
            Password: c.password,
        }),
    }

    if c.insecure {
        tr := &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true,
            },
        }
        opts = append(opts, remote.WithTransport(tr))
    }

    return remote.Write(ref, img, opts...)
}
```

## Integration Approach

### 1. Command Integration
```go
// pkg/cmd/push/root.go
var PushCmd = &cobra.Command{
    Use:   "push [image] [registry]",
    Short: "Push OCI images to a registry",
    Long:  `Push OCI images to a registry with authentication support`,
    Args:  cobra.ExactArgs(2),
    RunE:  pushE,
}

func init() {
    PushCmd.Flags().StringP("username", "u", "", "Registry username")
    PushCmd.Flags().StringP("password", "p", "", "Registry password")
    PushCmd.Flags().Bool("insecure", false, "Allow insecure registry connections")
}
```

### 2. Secret Integration
```go
// Reuse existing secret retrieval
func getDefaultCredentials(ctx context.Context, kubeClient client.Client) (*Credentials, error) {
    secret, err := getCorePackageSecret(ctx, kubeClient, "gitea", giteaAdminSecretName)
    if err != nil {
        return nil, err
    }
    return &Credentials{
        Username: string(secret.Data["username"]),
        Password: string(secret.Data["token"]),
    }, nil
}
```

### 3. Error Handling Pattern
```go
// Consistent with existing codebase
if err := validateInputs(image, registry); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}

if err := client.Push(ctx, image, registry); err != nil {
    return fmt.Errorf("push failed: %w", err)
}
```

## Risk Analysis & Mitigation

### Technical Risks

#### 1. Registry Compatibility
- **Risk**: Gitea registry might have quirks
- **Mitigation**: Use standard OCI operations only, extensive integration testing
- **Test Strategy**: Test against multiple registry implementations

#### 2. Large Image Handling
- **Risk**: Memory issues with large images
- **Mitigation**: Stream-based operations, chunked uploads
- **Test Strategy**: Test with various image sizes

#### 3. Network Reliability
- **Risk**: Network interruptions during push
- **Mitigation**: Implement retry logic with exponential backoff
- **Test Strategy**: Chaos testing with network simulation

### Security Risks

#### 1. Credential Exposure
- **Risk**: Passwords in command line/logs
- **Mitigation**: Mask sensitive data, use environment variables
- **Test Strategy**: Security scanning, log analysis

#### 2. Insecure Mode
- **Risk**: MITM attacks when --insecure used
- **Mitigation**: Clear warnings, require explicit confirmation
- **Test Strategy**: Security review, penetration testing

## Validation & Success Metrics

### Functional Validation
- [ ] Push command successfully pushes images
- [ ] Authentication works with flags
- [ ] Authentication works with defaults from secrets
- [ ] Insecure mode handles self-signed certificates
- [ ] Error messages are clear and actionable

### Quality Metrics
- [ ] Test coverage > 80%
- [ ] All tests written BEFORE implementation
- [ ] Zero regression in existing functionality
- [ ] Performance: Push 100MB image < 30 seconds
- [ ] Memory usage < 500MB for typical operations

### TDD Compliance Metrics
- [ ] Every feature has tests written first
- [ ] Test failures documented before fixes
- [ ] Refactoring phase completed for each cycle
- [ ] Mock objects used for external dependencies
- [ ] Integration tests cover all workflows

## Configuration Requirements

### setup-config.yaml
```yaml
project:
  name: "idpbuilder-push"
  description: "Push command for idpbuilder to upload OCI images to Gitea registry"
  target_dir: "/home/vscode/workspaces/idpbuilder-push"

target_repository:
  url: "https://github.com/jessesanford/idpbuilder.git"
  base_branch: "main"
  clone_depth: 100
  auth_method: "https"

technology:
  primary_language: "Go"
  frameworks:
    - "Cobra CLI"
    - "go-containerregistry"
    - "Kubernetes Client-Go"

agents:
  selected:
    - "Software Engineer"
    - "Code Reviewer"
    - "Architect"
  expertise:
    - "API Design"
    - "Testing Strategies"
    - "Cloud Architecture"

implementation:
  plan_type: "generate"
  project_type: "CLI Tool"
  estimated_loc: 4200
  num_phases: 4
  test_coverage: 80

constraints:
  max_lines_per_effort: 700
  max_parallel_agents: 3
  code_review: "mandatory"
  security_level: 1

tdd:
  enforce_tdd: true
  require_tests_first: true
  test_types:
    - "unit"
    - "integration"
    - "e2e"
  min_coverage_before_merge: 80
```

### target-repo-config.yaml
```yaml
repository_path: "/home/vscode/workspaces/idpbuilder-push/efforts/idpbuilder"
repository_name: "idpbuilder"
default_branch: "main"

target_repository:
  url: "https://github.com/jessesanford/idpbuilder.git"
  base_branch: "main"
  auth_method: "https"
  clone_depth: 100

branch_naming:
  project_prefix: "idpbuilder-push"
  effort_format: "{prefix}/phase{phase}/wave{wave}/{effort_name}"
  integration_format: "{prefix}/phase{phase}/wave{wave}/integration"
  phase_integration_format: "{prefix}/phase{phase}/integration"

workspace:
  efforts_root: "efforts"
  effort_path: "phase{phase}/wave{wave}/{effort_name}"

testing:
  test_command: "go test ./pkg/cmd/push/... ./pkg/oci/..."
  coverage_command: "go test -cover ./pkg/cmd/push/... ./pkg/oci/..."
  min_coverage: 80
  test_before_merge: true
```

## Recommended Approach for Orchestrator

### Initialization Sequence
1. Clone idpbuilder repository
2. Create feature branch structure
3. Set up test infrastructure
4. Begin Phase 1 Wave 1 with TDD approach

### Agent Deployment Strategy
1. **Software Engineer**: Implement test-first development
2. **Code Reviewer**: Ensure TDD compliance and quality
3. **Architect**: Validate design patterns and integration

### Critical Checkpoints
- After each RED phase: Verify tests fail correctly
- After each GREEN phase: Verify minimal implementation
- After each REFACTOR: Verify no regression
- End of each wave: Full test suite execution
- End of each phase: Integration testing

## Conclusion

This requirements analysis provides a comprehensive foundation for implementing the idpbuilder push command using TEST-DRIVEN DEVELOPMENT. The phased approach ensures incremental delivery of value while maintaining high quality through rigorous testing. The architecture leverages existing idpbuilder patterns while introducing clean abstractions for OCI operations.

Key success factors:
1. Strict adherence to TDD methodology
2. Incremental development with continuous integration
3. Comprehensive testing at all levels
4. Clean separation of concerns
5. Reuse of existing idpbuilder patterns and infrastructure

The implementation should result in a robust, well-tested push command that seamlessly integrates with the existing idpbuilder ecosystem while maintaining OCI standards compliance.