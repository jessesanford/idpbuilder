# Implementation Plan for E1.1.1: Minimal Build Types
Created: 2025-08-22 08:02:10 UTC
Created by: @agent-code-reviewer
Phase: 1 - MVP Core
Wave: 1 - Essential API Contracts

## Context Analysis

### Completed Related Efforts
- None (This is the first effort in Phase 1)

### Established Patterns
- None yet (Foundation effort for all build types)

### Integration Points
- API: Will establish the core API types for all build operations
- Services: Types will be consumed by future builder implementations
- Testing: Must establish testing patterns for type validation

## Requirements (from Phase Plan)

### Primary Requirements
1. MUST implement core request/response types only
2. BuildRequest with essential fields
3. BuildResponse with success/error info
4. NO complex validation - keep it simple

### Derived Requirements
1. Must include basic validation method for BuildRequest
2. Types must be JSON serializable for future API usage
3. Error handling approach must be consistent
4. Default values must be sensible (e.g., "latest" tag)

### Non-Functional Requirements
- Performance: Simple structs with minimal validation overhead
- Security: No sensitive data exposure in types
- Scalability: Types must support future extensions without breaking changes

## Implementation Strategy

### Approach
Create minimal, foundational types that establish the core data structures for container build operations. Focus on simplicity and essential fields only, avoiding any complex validation or business logic.

### Design Decisions
1. **JSON Tags**: Include json tags for future API compatibility
2. **Validation Method**: Single Validate() method on BuildRequest for essential checks
3. **Default Tag**: Apply "latest" as default when ImageTag is empty
4. **Error Handling**: Simple string-based error reporting in BuildResponse
5. **Immutable Design**: Structs are value types, no pointer fields

### Patterns to Follow
- Standard Go struct patterns with exported fields
- JSON serialization compatibility
- Simple validation with early returns
- Clear error messages

## Implementation Steps

### Step 1: Create Core Types
**Action**: Create BuildRequest and BuildResponse structs with essential fields
**Files**: pkg/build/api/types.go
**Validation**: Structs compile and have correct JSON tags

### Step 2: Add Basic Validation
**Action**: Implement Validate() method on BuildRequest with essential field checks
**Files**: pkg/build/api/types.go (add to existing)
**Validation**: Method correctly identifies missing required fields and applies defaults

### Step 3: Create Comprehensive Tests
**Action**: Implement unit tests covering all validation scenarios
**Files**: pkg/build/api/types_test.go
**Validation**: All tests pass and cover edge cases

### Step 4: Verify Build and Integration
**Action**: Ensure package builds cleanly and types integrate properly
**Files**: All created files
**Validation**: go build succeeds, no lint issues

## Files to Create/Modify

### New Files
```
pkg/
├── build/
│   └── api/
│       ├── types.go           # ~60 lines - Core type definitions
│       └── types_test.go      # ~40 lines - Validation tests
```

### Modified Files
- None (Creating new package structure)

## Code Templates

### BuildRequest Type
```go
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
```

### BuildResponse Type
```go
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
```

### Validation Method
```go
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

## Testing Requirements

### Unit Tests
- [ ] Test valid BuildRequest validation passes
- [ ] Test missing DockerfilePath fails validation
- [ ] Test missing ContextDir fails validation
- [ ] Test missing ImageName fails validation
- [ ] Test default ImageTag behavior when empty
- [ ] Test BuildResponse JSON serialization

### Integration Tests
- [ ] Test types integrate with json.Marshal/Unmarshal
- [ ] Test validation works in realistic scenarios

### Coverage Target
- Minimum: 80%
- Target: 90%

### Test File Structure
```
pkg/build/api/
├── types.go
└── types_test.go
```

## Size Management

### Estimated Size
- Core implementation: ~60 lines
- Tests: ~40 lines
- Total: ~100 lines

### Size Limit
- Maximum: 100 lines (as specified in requirements)
- Measurement: line-counter.sh

### Split Strategy (if needed)
Not anticipated - this is a minimal types-only implementation. If somehow over limit:
1. Core types in types.go
2. Validation logic in separate file
3. Tests remain together

## Success Criteria

### Functional
- [ ] BuildRequest struct with all required fields implemented
- [ ] BuildResponse struct with success/error reporting
- [ ] Basic validation on BuildRequest with sensible defaults
- [ ] JSON serialization works correctly

### Quality
- [ ] Tests pass with coverage >= 80%
- [ ] Package builds without warnings
- [ ] Code follows Go best practices
- [ ] Documentation comments on exported types

### Size
- [ ] Under 100 lines per line-counter.sh
- [ ] Clean, focused implementation

### Documentation
- [ ] All exported types have doc comments
- [ ] Field purposes clearly documented
- [ ] Validation behavior documented

## Integration Notes

### Dependencies
- Depends on: None (foundation effort)
- Required by: E1.1.2 (Builder Interface) - will import these types

### API Contracts
- Structs must be JSON serializable
- Validation must be consistent and predictable
- Error messages must be clear and actionable

### Breaking Changes
- None expected (new package)

## Risk Mitigation

### Technical Risks
| Risk | Mitigation | Plan |
|------|------------|------|
| Field names not future-proof | Use clear, descriptive names | Review with maintainer specs |
| Validation too complex | Keep to absolute minimum | Only check required fields |
| JSON compatibility issues | Test serialization explicitly | Include marshal/unmarshal tests |

### Process Risks
| Risk | Mitigation | Plan |
|------|------------|------|
| Scope creep | Stick to exact phase plan specs | Reject any additions beyond MVP |
| Over-engineering | No interfaces, no complex patterns | Simple structs only |

## Notes

1. **Simplicity Focus**: This is the foundation - keep it as simple as possible
2. **Future Compatibility**: Design choices should not prevent future extensions
3. **Testing Foundation**: Establish good testing patterns for future efforts
4. **Documentation**: Clear comments are critical for maintainability

This implementation creates the essential data structures that all subsequent build functionality will depend on. Success here enables the rest of Phase 1.