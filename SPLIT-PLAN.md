# Split Plan for buildah-build-wrapper Effort

## Overview
**Total Size**: 983 lines (exceeds 800 line limit by 183 lines)
**Splits Required**: 2
**Planner**: Code Reviewer Instance
**Created**: 2025-08-29

## Split Inventory

### Split 001: Core Buildah Wrapper Implementation
**Target Size**: ~500 lines
**Branch**: `idpbuilder-oci-mvp/phase2/wave1/buildah-build-wrapper-split-001`
**Files**:
- `pkg/build/types.go` (50 lines) - Core type definitions
- `pkg/build/builder.go` (97 lines) - Mock/fallback implementation
- `pkg/build/builder_buildah.go` (278 lines) - Main buildah implementation
- Basic unit tests for core functionality (~75 lines)

**Total**: ~500 lines

### Split 002: Test Suite and Integration
**Target Size**: ~483 lines  
**Branch**: `idpbuilder-oci-mvp/phase2/wave1/buildah-build-wrapper-split-002`
**Files**:
- `pkg/build/builder_buildah_test.go` (257 lines) - Comprehensive buildah tests
- `pkg/build/builder_mock_test.go` (229 lines) - Mock builder tests
- Integration with Phase 1 trust manager improvements
- Documentation updates

**Total**: ~483 lines

## Split Strategy

### Split 001: Core Implementation (Priority 1)
This split contains the essential buildah wrapper functionality:

1. **Type Definitions** (`types.go`):
   - Builder interface definition
   - BuildOptions structure
   - BuildResult structure
   - ImageInfo structure

2. **Mock Implementation** (`builder.go`):
   - Fallback implementation when buildah is not available
   - Uses build tags for conditional compilation
   - Provides TrustManager interface for Phase 1 integration

3. **Buildah Implementation** (`builder_buildah.go`):
   - Real buildah integration using containers/buildah libraries
   - BuildImage functionality with proper error handling
   - ListImages, RemoveImage, and TagImage operations
   - Integration with Phase 1 TrustManager for certificate handling

4. **Basic Tests**:
   - Core functionality tests
   - Happy path validation
   - Basic error handling tests

### Split 002: Test Coverage and Polish (Priority 2)
This split focuses on comprehensive testing and integration:

1. **Comprehensive Test Suite** (`builder_buildah_test.go`):
   - Full test coverage for buildah operations
   - Error condition testing
   - Edge case handling
   - Integration tests with mock storage

2. **Mock Implementation Tests** (`builder_mock_test.go`):
   - Tests for fallback implementation
   - Build tag validation
   - Mock builder behavior tests

3. **Phase 1 Integration Polish**:
   - Enhanced TrustManager integration
   - Certificate validation improvements
   - Security context configuration

## Implementation Order

1. **Split 001 First**: Implement core buildah wrapper functionality
   - Complete types and interfaces
   - Implement both mock and real buildah builders
   - Add basic tests for validation
   - Ensure it compiles and passes basic tests

2. **Split 002 Second**: Add comprehensive testing
   - Depends on Split 001 being complete
   - Adds full test coverage
   - Enhances Phase 1 integration
   - Polishes error handling

## Verification Checklist

### Split 001:
- [ ] Types compile without errors
- [ ] Mock builder works in non-buildah environments
- [ ] Buildah builder compiles with build tags
- [ ] Basic tests pass
- [ ] Under 500 lines total

### Split 002:
- [ ] All tests pass
- [ ] Test coverage >80%
- [ ] Phase 1 integration validated
- [ ] Under 500 lines total
- [ ] Combined splits = original functionality

## Integration Notes

- Both splits must be merged sequentially (001 then 002)
- Split 002 depends on Split 001 being merged first
- Final integration testing required after both splits merge
- Ensure build tags work correctly in both environments

## Risk Mitigation

- Each split independently compilable
- Split 001 provides minimum viable functionality
- Split 002 adds quality and coverage without breaking core
- Clear dependency chain prevents integration issues