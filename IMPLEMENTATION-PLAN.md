# OCI & Stack Configuration Types Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: idpbuilder-oci-mgmt/phase1/wave1/oci-stack-types
**Can Parallelize**: Yes
**Parallel With**: [auth-cert-types, error-progress-types]
**Size Estimate**: 500 lines
**Dependencies**: None (foundation effort)

## Overview
- **Effort**: E1.1.1 - OCI & Stack Configuration Types
- **Phase**: 1, Wave: 1
- **Estimated Size**: 500 lines
- **Implementation Time**: 2-3 hours

## Dependency Context (R219)
- **Dependencies Analyzed**: None - this is a foundation effort
- **Influence on Implementation**: As a foundation effort, this establishes base types that other efforts will import
- **Reuse Opportunities**: None from dependencies (no dependencies exist)
- **Integration Strategy**: Other efforts will import types from pkg/oci/api

## File Structure
```
efforts/phase1/wave1/oci-stack-types/
└── pkg/
    └── oci/
        └── api/
            ├── interfaces.go    # 120 lines - Service contracts
            ├── types.go         # 200 lines - Core data structures
            ├── validation.go    # 80 lines - Validation logic
            └── types_test.go    # 100 lines - Test coverage
```

## Implementation Steps

### Step 1: Create Directory Structure
```bash
mkdir -p pkg/oci/api
```

### Step 2: Implement interfaces.go (120 lines)
Define service contracts for OCI operations:
- `OCIBuildService` interface with build operations
- `OCIRegistryService` interface with registry operations
- `StackOCIManager` interface with stack-specific operations

Key methods to implement:
- Build lifecycle: Initialize, Shutdown
- Build operations: BuildImage, BuildFromDockerfile
- Registry operations: Connect, PushImage, PullImage
- Stack operations: BuildStackImage, PushStackImage

### Step 3: Implement types.go (200 lines)
Create core data structures:
- `BuildConfig` - buildah configuration with storage, runtime, network settings
- `RegistryConfig` - registry connection configuration
- `StackOCIConfig` - stack-specific OCI configuration
- `BuildRequest` - build operation request structure
- `BuildResult` - build operation result with metadata
- `BuildStatus` - current build status tracking
- Supporting types: BuildOptions, PushOptions, PullOptions, LayerInfo, etc.

Validation tags to include:
- Use struct tags for validation (e.g., `validate:"required,oneof=..."`)
- Time durations with min constraints
- URL/IP validation where appropriate

### Step 4: Implement validation.go (80 lines)
Add validation logic:
- Custom validator initialization with `NewValidator()`
- Custom validation functions:
  - `validateImageTag` - OCI image tag format validation
  - `validateSemver` - semantic version validation
  - `validatePlatform` - platform format (os/arch) validation
- Business logic validation:
  - `ValidateBuildConfig` - comprehensive config validation
  - `ValidateStackConfig` - stack configuration validation
  - Rootless mode constraints (requires vfs storage driver)

### Step 5: Implement types_test.go (100 lines)
Create comprehensive test coverage:
- `TestBuildConfig_Validation` - test various config combinations
  - Valid overlay config
  - Valid vfs config for rootless
  - Invalid storage driver
  - Rootless with overlay (should fail)
  - Negative CPU quota validation
- `TestImageTagValidation` - test image tag formats
  - Valid tags with registry prefix
  - Invalid characters and formats
  - Empty tag handling

### Step 6: Testing & Validation
```bash
# Run tests after each file
go test ./pkg/oci/api/...

# Check test coverage
go test -cover ./pkg/oci/api/...

# Ensure minimum 80% coverage
```

## Size Management
- **Estimated Lines**: 500 total
- **Measurement Tool**: /home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh
- **Check Frequency**: After implementing each file
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

### Size Breakdown by File:
| File | Estimated Lines | Purpose |
|------|----------------|---------|
| interfaces.go | 120 | Service contracts |
| types.go | 200 | Data structures |
| validation.go | 80 | Validation logic |
| types_test.go | 100 | Test coverage |
| **Total** | **500** | Within limit |

## Test Requirements
- **Unit Tests**: 80% coverage minimum
- **Integration Tests**: N/A (types only, no integration needed)
- **Test Files**: pkg/oci/api/types_test.go

### Test Coverage Areas:
1. **Validation Testing**:
   - All custom validators must be tested
   - Edge cases for each validation rule
   - Business logic constraints

2. **Type Constraints**:
   - Required fields validation
   - Enum value constraints
   - Numeric range validation

3. **Error Scenarios**:
   - Invalid configurations
   - Missing required fields
   - Conflicting settings

## Pattern Compliance
- **Go Best Practices**:
  - Interface-first design
  - Clear separation of concerns
  - Comprehensive validation
  - Proper error handling

- **OCI/Container Patterns**:
  - Standard OCI image tag formats
  - Registry URL conventions
  - Platform specification compliance

- **Testing Patterns**:
  - Table-driven tests
  - Clear test naming
  - Isolated test cases

## Dependencies
External packages required:
```go
require (
    github.com/go-playground/validator/v10 v10.15.5
)
```

## Integration Points
This effort provides foundation types that will be used by:
- Wave 2: Build service implementation (imports api types)
- Wave 3: Registry client implementation (imports api types)
- Wave 4: Stack manager implementation (imports api and uses StackOCIConfig)

## Risk Assessment
- **Low Risk**: Well-defined types with clear validation
- **Mitigation**: Comprehensive test coverage ensures type safety

## Success Criteria
- [ ] All files compile without errors
- [ ] All tests pass with >80% coverage
- [ ] Size remains under 500 lines
- [ ] Validation logic is comprehensive
- [ ] No circular dependencies
- [ ] Ready for other efforts to import

## Notes for SW Engineer
1. Start with interfaces.go to establish contracts
2. Implement types.go with all validation tags
3. Add custom validators in validation.go
4. Write comprehensive tests in types_test.go
5. Run line-counter.sh after each file to monitor size
6. Commit frequently with clear messages
7. This is a foundation effort - other efforts depend on these types being correct