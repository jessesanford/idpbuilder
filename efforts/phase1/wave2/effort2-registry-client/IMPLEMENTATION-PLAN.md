# Implementation Plan for E1.2.2: Registry Client
Created: 2025-08-22 13:23:51 UTC
Created by: @agent-code-reviewer
Phase: 1 - MVP Core
Wave: 2 - Core Libraries

## Context Analysis

### Completed Related Efforts
- E1.1.1: Minimal Build Types - Implemented base API types including core configuration structures and status types following Kubernetes patterns
- E1.1.2: Builder Interface - Created foundational interfaces for build operations and established patterns for client abstraction

### Established Patterns
- Kubernetes API patterns with TypeMeta and ObjectMeta
- Standard Go project structure with apis/v1alpha1 package organization
- Interface-based client abstraction for testability
- Configuration through structured types with JSON tags
- DeepCopy method generation for API types

### Integration Points
- API: Uses existing v1alpha1 API types from Wave 1
- Registry: Must integrate with Gitea container registry
- TLS: Must handle self-signed certificates for development environments
- Configuration: Integrates with existing configuration patterns

## Requirements (from Phase Plan)

### Primary Requirements
1. Create a registry client that can authenticate with Gitea registry
2. Support container image pushing operations
3. Handle self-signed TLS certificates for local development
4. Provide clean error handling and retry mechanisms
5. Follow established interface patterns from Wave 1

### Derived Requirements
1. Client interface must be mockable for testing
2. Configuration should use existing config patterns
3. Support multiple authentication methods (token, basic auth)
4. Graceful handling of network timeouts and registry errors
5. Integration with existing logging patterns

### Non-Functional Requirements
- Performance: Should handle concurrent push operations
- Security: Secure credential handling, no credentials in logs
- Scalability: Support for multiple registry endpoints
- Reliability: Retry logic with exponential backoff

## Implementation Strategy

### Approach
Create a layered registry client with interface abstraction, concrete implementation using standard container registry protocols, and comprehensive error handling. Build on existing API patterns while adding registry-specific functionality.

### Design Decisions
1. Interface Pattern: Follow builder interface patterns from E1.1.2 for consistency
2. Configuration Integration: Extend existing config types rather than creating new ones
3. TLS Handling: Create reusable TLS configuration utilities
4. Error Types: Define domain-specific error types for better error handling
5. Testing Strategy: Interface-based design enables comprehensive unit testing

### Patterns to Follow
- Interface-based client design from E1.1.2
- Configuration struct patterns from E1.1.1
- Standard Go error handling with wrapped errors
- Kubernetes-style status reporting

## Implementation Steps

### Step 1: Define Registry Client Interface
**Action**: Create core interface definition for registry operations
**Files**: pkg/registry/interface.go
**Validation**: Interface compiles and follows established patterns

### Step 2: Registry Configuration Types
**Action**: Extend existing config types with registry-specific fields
**Files**: pkg/registry/types.go
**Validation**: Types follow v1alpha1 patterns and include proper JSON tags

### Step 3: TLS Configuration Helper
**Action**: Create reusable TLS configuration utilities
**Files**: pkg/registry/tls.go
**Validation**: Can handle self-signed certificates properly

### Step 4: Core Registry Client Implementation
**Action**: Implement the registry client interface
**Files**: pkg/registry/client.go
**Validation**: Can perform basic registry authentication

### Step 5: Authentication Methods
**Action**: Implement token and basic authentication support
**Files**: pkg/registry/auth.go
**Validation**: Both auth methods work with test registry

### Step 6: Error Handling and Retry Logic
**Action**: Add comprehensive error handling with retry mechanisms
**Files**: pkg/registry/errors.go
**Validation**: Proper error classification and retry behavior

### Step 7: Image Push Operations
**Action**: Implement container image push functionality
**Files**: pkg/registry/push.go
**Validation**: Can successfully push test images

### Step 8: Comprehensive Unit Tests
**Action**: Create thorough test coverage for all components
**Files**: pkg/registry/*_test.go
**Validation**: Tests pass and coverage meets requirements

## Files to Create/Modify

### New Files
```
pkg/registry/
├── interface.go      # Registry client interface definition
├── types.go         # Configuration and data types
├── client.go        # Main client implementation
├── auth.go          # Authentication methods
├── tls.go           # TLS configuration utilities
├── push.go          # Image push operations
├── errors.go        # Domain-specific error types
├── client_test.go   # Client implementation tests
├── auth_test.go     # Authentication tests
├── tls_test.go      # TLS configuration tests
├── push_test.go     # Push operation tests
└── errors_test.go   # Error handling tests
```

### Modified Files
- No existing files require modification (clean addition)

## Code Templates

### Registry Client Interface
```go
// RegistryClient defines the interface for container registry operations
type RegistryClient interface {
    // Push pushes a container image to the registry
    Push(ctx context.Context, imageRef string, opts PushOptions) error
    
    // Authenticate authenticates with the registry
    Authenticate(ctx context.Context, config AuthConfig) error
    
    // Health checks registry connectivity
    Health(ctx context.Context) error
}

// PushOptions contains options for push operations
type PushOptions struct {
    Tags     []string
    Insecure bool
    Timeout  time.Duration
}
```

### Registry Configuration Types
```go
// RegistryConfig contains registry connection settings
type RegistryConfig struct {
    // URL is the registry base URL
    URL string `json:"url,omitempty"`
    
    // Insecure allows insecure TLS connections
    Insecure bool `json:"insecure,omitempty"`
    
    // Auth contains authentication configuration
    Auth AuthConfig `json:"auth,omitempty"`
}

// AuthConfig contains authentication settings
type AuthConfig struct {
    // Token for token-based authentication
    Token string `json:"token,omitempty"`
    
    // Username for basic authentication
    Username string `json:"username,omitempty"`
    
    // Password for basic authentication
    Password string `json:"password,omitempty"`
}
```

## Testing Requirements

### Unit Tests
- [ ] Test registry client interface compliance
- [ ] Test authentication methods (token and basic)
- [ ] Test TLS configuration with self-signed certificates
- [ ] Test push operations with various image formats
- [ ] Test error handling and retry mechanisms
- [ ] Test configuration validation and defaults

### Integration Tests
- [ ] Test against local Gitea registry
- [ ] Test end-to-end image push workflow
- [ ] Test authentication failure scenarios
- [ ] Test network timeout handling

### Coverage Target
- Minimum: 70% (Phase 1 requirement)
- Target: 85%

### Test File Structure
```
pkg/registry/
├── client_test.go
├── auth_test.go
├── tls_test.go
├── push_test.go
├── errors_test.go
└── integration_test.go
```

## Size Management

### Estimated Size
- Core implementation: ~180 lines
- Authentication logic: ~60 lines  
- TLS utilities: ~40 lines
- Error handling: ~30 lines
- Tests: ~90 lines
- Total: ~400 lines

### Size Limit
- Maximum: 800 lines
- Measurement: line-counter.sh

### Split Strategy (if needed)
If approaching limit:
1. Complete core client and auth functionality first (~280 lines)
2. Split push operations to separate effort if needed
3. Prioritize working authentication and basic connectivity

## Success Criteria

### Functional
- [ ] Registry client can authenticate with Gitea registry
- [ ] Can push container images successfully
- [ ] Handles self-signed certificates properly
- [ ] Provides clear error messages for common failures
- [ ] Follows established interface patterns

### Quality
- [ ] All tests pass
- [ ] Coverage >= 70%
- [ ] Lint clean (golangci-lint)
- [ ] Build successful
- [ ] No hardcoded credentials or URLs

### Size
- [ ] Under 800 lines per line-counter.sh
- [ ] Well-organized single effort (no split needed)

### Documentation
- [ ] All exported functions documented
- [ ] Interface documentation complete
- [ ] Example usage in tests
- [ ] Updated work-log.md

## Integration Notes

### Dependencies
- Depends on: E1.1.1 (Minimal Build Types), E1.1.2 (Builder Interface)
- Required by: E1.3.x (MVP Implementation efforts in Wave 3)

### API Contracts
- Interface: RegistryClient with Push, Authenticate, Health methods
- Configuration: Extends existing config patterns with RegistryConfig type
- Error handling: Domain-specific error types for registry operations

### Breaking Changes
- None expected (new functionality only)

## Special Considerations

### Gitea Registry Integration
- Must work with Gitea's container registry API
- Support for Gitea token authentication
- Handle Gitea-specific error responses

### Development Environment
- Support for local development with self-signed certificates
- Configurable timeouts for slow local networks
- Clear error messages for common setup issues

### Security
- Never log credentials or tokens
- Secure handling of authentication data
- Proper certificate validation options