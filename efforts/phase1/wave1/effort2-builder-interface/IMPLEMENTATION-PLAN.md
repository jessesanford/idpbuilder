# Implementation Plan for E1.1.2: Builder Interface
Created: 2025-08-22 08:04:33 UTC
Created by: @agent-code-reviewer
Phase: 1 - MVP Core
Wave: 1 - Essential API Contracts

## Context Analysis

### Completed Related Efforts
- E1.1.1: Expected to complete minimal build types (BuildRequest, BuildResponse structs with validation)

### Established Patterns
- Using pkg/ directory structure for organizational packages
- Simple validation methods on request types
- Standard Go interface patterns
- Minimal error types without complex validation

### Integration Points
- API: Must be compatible with E1.1.1 BuildRequest/BuildResponse types
- Services: Will be implemented by Buildah client in Wave 1.2
- Registry: Interface will handle registry configuration

## Requirements (from Phase Plan)

### Primary Requirements
1. Define Builder interface with BuildAndPush method only
2. Define BuilderConfig struct with registry configuration
3. Create DefaultConfig function returning MVP configuration
4. Interface must be compatible with E1.1.1 types
5. NO implementations - pure interface contracts only

### Derived Requirements
1. Interface must use context.Context for cancellation support
2. Interface must return BuildResponse compatible with E1.1.1
3. Configuration must support insecure TLS for MVP development
4. Registry and namespace must be hardcoded for MVP

### Non-Functional Requirements
- Performance: Minimal overhead from interface abstraction
- Security: Support insecure TLS for development environment
- Scalability: Single method interface suitable for future extension

## Implementation Strategy

### Approach
Create minimal interface definitions that establish the contract for build operations without any implementation logic. Focus on compatibility with existing types and future Buildah integration.

### Design Decisions
1. **Single Method Interface**: BuildAndPush combines build and push operations for MVP simplicity
2. **Configuration Struct**: Separate BuilderConfig allows for future extension without breaking interface
3. **Default Configuration**: DefaultConfig() provides MVP-ready configuration
4. **Context Support**: Use context.Context for cancellation and timeout support

### Patterns to Follow
- Standard Go interface patterns (interface{} with behavior methods)
- Configuration struct pattern for dependency injection
- Default constructor pattern for common configurations

## Implementation Steps

### Step 1: Create Builder Interface File
**Action**: Create pkg/build/builder/interface.go with Builder interface definition
**Files**: pkg/build/builder/interface.go
**Validation**: Interface compiles and can be imported

### Step 2: Define Builder Interface
**Action**: Define Builder interface with single BuildAndPush method
**Files**: pkg/build/builder/interface.go
**Validation**: Method signature matches requirements from phase plan

### Step 3: Create Configuration Support
**Action**: Add BuilderConfig struct and DefaultConfig function
**Files**: pkg/build/builder/interface.go  
**Validation**: Default config returns correct MVP values

### Step 4: Create Registry Interface File
**Action**: Create pkg/build/registry/interface.go with Registry interface
**Files**: pkg/build/registry/interface.go
**Validation**: Interface compiles and defines Push method

### Step 5: Create Interface Tests
**Action**: Create basic tests proving interfaces can be implemented
**Files**: pkg/build/builder/interface_test.go, pkg/build/registry/interface_test.go
**Validation**: Tests pass, interface implementations work

## Files to Create/Modify

### New Files
```
pkg/
├── build/
│   ├── builder/
│   │   ├── interface.go           # ~25 lines
│   │   └── interface_test.go      # ~20 lines  
│   └── registry/
│       ├── interface.go           # ~15 lines
│       └── interface_test.go      # ~15 lines
```

### Modified Files
None - all new interface definitions

## Code Templates

### Builder Interface
```go
// pkg/build/builder/interface.go
package builder

import (
    "context"
    "idpbuilder/pkg/build/api"
)

// Builder defines the interface for container build operations
type Builder interface {
    // BuildAndPush builds a container image and pushes to registry
    BuildAndPush(ctx context.Context, req api.BuildRequest) (*api.BuildResponse, error)
}

// BuilderConfig holds configuration for builder instances
type BuilderConfig struct {
    // Registry is the target OCI registry
    Registry string
    
    // Namespace is the registry namespace
    Namespace string
    
    // InsecureSkipTLSVerify skips TLS verification for development
    InsecureSkipTLSVerify bool
}

// DefaultConfig returns MVP configuration for development
func DefaultConfig() BuilderConfig {
    return BuilderConfig{
        Registry:              "gitea.cnoe.localtest.me", 
        Namespace:             "giteaadmin",
        InsecureSkipTLSVerify: true,
    }
}
```

### Registry Interface  
```go
// pkg/build/registry/interface.go
package registry

import "context"

// Registry defines the interface for registry operations
type Registry interface {
    // Push pushes an image to the registry
    Push(ctx context.Context, imageRef string) error
}
```

## Testing Requirements

### Unit Tests
- [ ] Test that Builder interface can be implemented by mock struct
- [ ] Test that Registry interface can be implemented by mock struct  
- [ ] Test DefaultConfig returns expected values
- [ ] Test BuilderConfig struct can be constructed

### Integration Tests
Not applicable - pure interface definitions with no external dependencies

### Coverage Target
- Minimum: 80% (high due to simple interface code)
- Target: 90%

### Test File Structure
```
pkg/build/builder/
├── interface.go
└── interface_test.go
pkg/build/registry/
├── interface.go  
└── interface_test.go
```

## Size Management

### Estimated Size
- Builder interface: ~25 lines
- Registry interface: ~15 lines
- Tests: ~35 lines
- Total: ~75 lines

### Size Limit
- Maximum: 800 lines
- Current estimate: 75 lines
- Safety margin: 725 lines

### Split Strategy (if needed)
Not needed - well under size limit. If interfaces grow unexpectedly:
1. Keep core Builder interface in first split
2. Move Registry interface to separate effort if needed

## Success Criteria

### Functional
- [ ] Builder interface defines required BuildAndPush method
- [ ] Registry interface defines required Push method
- [ ] BuilderConfig provides necessary configuration fields
- [ ] DefaultConfig returns MVP-appropriate values
- [ ] Interfaces compile without errors
- [ ] Compatible with E1.1.1 api types

### Quality
- [ ] Tests pass demonstrating interface implementations
- [ ] Coverage >= 80%
- [ ] Lint clean
- [ ] Build successful
- [ ] No hardcoded values (use constants/config)

### Size
- [ ] Under 800 lines per line-counter.sh
- [ ] Estimated ~75 lines actual

### Documentation
- [ ] Interface methods have clear godoc comments
- [ ] Configuration struct fields documented
- [ ] Package documentation explains purpose

## Integration Notes

### Dependencies
- Depends on: E1.1.1 (Build Types) - requires api.BuildRequest and api.BuildResponse types
- Required by: E1.2.1 (Buildah Client) - will implement these interfaces

### API Contracts
- Interface: Builder.BuildAndPush(context.Context, api.BuildRequest) (*api.BuildResponse, error)
- Interface: Registry.Push(context.Context, string) error
- Config: DefaultConfig() returns ready-to-use configuration

### Breaking Changes
- None expected - new interface definitions

## Critical Implementation Notes

1. **Import Path**: Must import api types from E1.1.1 effort
2. **Context Usage**: All interface methods must accept context.Context as first parameter
3. **Error Handling**: Interfaces define error returns but don't specify error types
4. **Configuration**: BuilderConfig struct supports future extension without breaking interface
5. **Registry Hardcoding**: MVP values hardcoded in DefaultConfig for development environment

## Effort-Specific Constraints

- **Interface Only**: NO implementation logic whatsoever
- **Compatibility**: Must work with api.BuildRequest/BuildResponse from E1.1.1
- **Simplicity**: Single method interface for MVP, extensible for future phases
- **Development Focus**: Configuration optimized for gitea.cnoe.localtest.me environment