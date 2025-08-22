# Implementation Plan for E1.2.1: Buildah Client

Created: 2025-08-22 13:24:00 UTC
Created by: @agent-code-reviewer
Phase: 1 - MVP Core
Wave: 2 - Core Libraries

## Context Analysis

### Completed Related Efforts
- E1.1.1: Minimal Build Types - Core API types and interfaces for container building
- E1.1.2: Builder Interface - Abstract interface definitions for builder implementations

### Established Patterns
- Package structure: `pkg/{domain}/{component}.go`
- API types defined in `api/v1alpha1/`
- Client libraries follow interface patterns established in Wave 1
- Configuration uses `BuildCustomizationSpec` from localbuild types
- Error handling with proper context propagation
- Test-driven development with table-driven tests

### Integration Points
- API: Uses existing `localbuild_types.go` and builder interfaces
- Build System: Integrates with existing `pkg/build/` package
- Configuration: Uses `BuildCustomizationSpec` for settings
- Logging: Uses controller-runtime logging patterns

## Requirements (from Phase Plan)

### Primary Requirements
1. Create a Buildah client wrapper that provides container building capabilities
2. Implement interface from Wave 1 builder interfaces
3. Support basic Dockerfile-based container builds
4. Handle registry authentication and configuration
5. Provide progress reporting and error handling
6. Support build context management

### Derived Requirements
1. Must integrate with existing build workflow in `pkg/build/build.go`
2. Configuration should use existing `BuildCustomizationSpec` patterns
3. Error handling must follow project conventions
4. Client must support containerized environments (running in pods)
5. Build context must handle both local and remote sources

### Non-Functional Requirements
- Performance: Handle builds efficiently without blocking
- Security: Proper handling of credentials and build contexts
- Scalability: Support concurrent builds
- Reliability: Robust error handling and cleanup

## Implementation Strategy

### Approach
Create a Buildah client library that implements the builder interfaces established in Wave 1. The client will provide a Go wrapper around Buildah commands, handling the complexities of container building while providing a clean interface for the idpbuilder system.

### Design Decisions
1. **Command Wrapper Pattern**: Use exec.Command to invoke buildah binary rather than embedding buildah as a library to avoid dependency complexity
2. **Interface Implementation**: Implement the builder interfaces from Wave 1 to ensure compatibility
3. **Context Management**: Use Go context for cancellation and timeout handling
4. **Build Isolation**: Each build gets its own working directory and cleanup

### Patterns to Follow
- Client interface pattern from existing codebase
- Context-based cancellation from controller-runtime
- Structured logging with named loggers
- Configuration injection pattern from build package

## Implementation Steps

### Step 1: Create Buildah Client Interface
**Action**: Define the buildah client interface and basic types
**Files**: `pkg/buildah/types.go`, `pkg/buildah/client.go`
**Validation**: Interface compiles and matches builder interface from Wave 1

### Step 2: Implement Build Command Execution
**Action**: Create buildah command wrapper with exec.Command
**Files**: `pkg/buildah/build.go`
**Validation**: Can execute basic buildah build commands

### Step 3: Add Configuration Support
**Action**: Integrate with existing BuildCustomizationSpec configuration
**Files**: `pkg/buildah/config.go`
**Validation**: Configuration properly parsed and applied to builds

### Step 4: Implement Progress Reporting
**Action**: Add build progress tracking and logging
**Files**: `pkg/buildah/progress.go`
**Validation**: Progress events properly emitted during builds

### Step 5: Add Error Handling and Cleanup
**Action**: Implement robust error handling and resource cleanup
**Files**: `pkg/buildah/errors.go`, update existing files
**Validation**: Failed builds properly cleaned up, errors properly categorized

### Step 6: Create Integration Layer
**Action**: Create integration with existing build package
**Files**: `pkg/buildah/integration.go`
**Validation**: Buildah client integrates cleanly with pkg/build/build.go

### Step 7: Comprehensive Testing
**Action**: Create comprehensive unit and integration tests
**Files**: `pkg/buildah/*_test.go`
**Validation**: All tests pass, coverage meets requirements

## Files to Create/Modify

### New Files
```
pkg/buildah/
├── client.go           # Main client implementation
├── types.go            # Type definitions and interfaces
├── build.go            # Build command execution
├── config.go           # Configuration handling
├── progress.go         # Progress reporting
├── errors.go           # Error handling utilities
├── integration.go      # Integration with existing build system
├── client_test.go      # Client tests
├── build_test.go       # Build execution tests
├── config_test.go      # Configuration tests
└── integration_test.go # Integration tests
```

### Modified Files
- `pkg/build/build.go`: Add buildah client integration point
- `go.mod`: Add buildah dependencies if needed

## Code Templates

### Buildah Client Interface
```go
// Package buildah provides container building capabilities using Buildah
package buildah

import (
    "context"
    "io"
    
    "github.com/cnoe-io/idpbuilder/api/v1alpha1"
)

// Client provides buildah-based container building
type Client interface {
    // Build builds a container image from a Dockerfile
    Build(ctx context.Context, opts BuildOptions) (*BuildResult, error)
    
    // Close cleans up client resources
    Close() error
}

// BuildOptions configures a container build
type BuildOptions struct {
    ContextDir    string
    Dockerfile    string
    Tag           string
    BuildArgs     map[string]string
    Progress      io.Writer
    Config        v1alpha1.BuildCustomizationSpec
}

// BuildResult contains the result of a build operation
type BuildResult struct {
    ImageID   string
    Tag       string
    Size      int64
    Duration  time.Duration
}
```

### Client Implementation Structure
```go
// client.go
type client struct {
    logger     logr.Logger
    workdir    string
    buildahCmd string
}

func NewClient(opts ClientOptions) (Client, error) {
    // Initialize buildah client with validation
    // Check buildah binary availability
    // Set up working directory
}

func (c *client) Build(ctx context.Context, opts BuildOptions) (*BuildResult, error) {
    // Validate build options
    // Prepare build context
    // Execute buildah build command
    // Parse results and return
}
```

## Testing Requirements

### Unit Tests
- [ ] Test buildah client creation and initialization
- [ ] Test build command generation and execution
- [ ] Test configuration parsing and validation
- [ ] Test error handling for various failure scenarios
- [ ] Test progress reporting functionality
- [ ] Test resource cleanup on success and failure

### Integration Tests
- [ ] Test full build workflow with sample Dockerfile
- [ ] Test integration with existing build package
- [ ] Test configuration integration with BuildCustomizationSpec
- [ ] Test concurrent builds

### Coverage Target
- Minimum: 80%
- Target: 85%

### Test File Structure
```
pkg/buildah/
├── client_test.go      # Client lifecycle tests
├── build_test.go       # Build execution tests
├── config_test.go      # Configuration tests
├── progress_test.go    # Progress reporting tests
├── errors_test.go      # Error handling tests
└── integration_test.go # Integration tests
```

## Size Management

### Estimated Size
- Core implementation: ~250 lines
- Configuration handling: ~80 lines
- Error handling: ~60 lines
- Tests: ~300 lines
- Total: ~690 lines

### Size Limit
- Maximum: 800 lines
- Measurement: line-counter.sh

### Split Strategy (if needed)
If approaching limit:
1. Complete core buildah client functionality first (Steps 1-3)
2. Split progress reporting and advanced error handling to separate effort
3. Prioritize basic build capability over advanced features

## Success Criteria

### Functional
- [ ] All requirements implemented
- [ ] Integration with existing build system works
- [ ] No hardcoded values or magic constants
- [ ] Proper buildah binary detection and error handling

### Quality
- [ ] Tests pass with 85% coverage minimum
- [ ] Lint clean (golangci-lint)
- [ ] Build successful without warnings
- [ ] Code follows project patterns and conventions

### Size
- [ ] Under 800 lines per line-counter.sh
- [ ] Properly organized if split needed

### Documentation
- [ ] Code comments for complex logic
- [ ] Interface documentation
- [ ] Integration examples
- [ ] Updated work-log.md throughout development

## Integration Notes

### Dependencies
- Depends on: E1.1.1 (Build Types), E1.1.2 (Builder Interface)
- Required by: E1.2.2 (Registry Client), E1.3.x (MVP Implementation efforts)

### API Contracts
- Implements: Builder interfaces defined in Wave 1
- Uses: BuildCustomizationSpec from localbuild types
- Integrates: With pkg/build/build.go workflow

### Breaking Changes
- None expected - implements existing interfaces

## Special Considerations

### Buildah Binary Dependency
- Must handle cases where buildah is not installed
- Provide clear error messages for missing dependencies
- Consider containerized execution environment

### Security Considerations
- Proper handling of build contexts and secrets
- Secure cleanup of temporary files
- Validation of build arguments and paths

### Performance Considerations
- Efficient handling of build contexts
- Proper resource cleanup
- Support for build cancellation