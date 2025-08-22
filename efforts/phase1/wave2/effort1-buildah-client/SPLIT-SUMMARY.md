# Split Plan for E1.2.1: Buildah Client

## Split Analysis

### Original Effort Status
- **Current Size**: 860 lines (60 lines over 800 limit)
- **Files**: 4 files total
- **Status**: Core functionality complete, over size limit

### Size Breakdown by File
```
pkg/buildah/build.go:   271 lines  (Build execution logic)
pkg/buildah/config.go:  229 lines  (Configuration handling)
pkg/buildah/client.go:  208 lines  (Client implementation)
pkg/buildah/types.go:   152 lines  (Interface definitions)
```

## Logical Split Design

### Split 1: E1.2.1 - Core Buildah Client (360 lines)
**Purpose**: Fundamental client interface and implementation
**Files Included**:
- `pkg/buildah/types.go` (152 lines) - Core interfaces, types, and options
- `pkg/buildah/client.go` (208 lines) - Client lifecycle and basic operations

**Functionality**:
- Client interface definitions
- Basic client implementation
- Client creation and initialization  
- Resource cleanup and lifecycle management
- Binary detection and validation

**Dependencies**: None (foundation layer)
**Size**: 360 lines (✅ Under 800 limit)

### Split 2: E1.2.1a - Build Operations (500 lines)  
**Purpose**: Advanced build execution and configuration
**Files Included**:
- `pkg/buildah/config.go` (229 lines) - Configuration integration
- `pkg/buildah/build.go` (271 lines) - Build command execution

**Functionality**:
- BuildCustomizationSpec integration
- Registry configuration handling
- Network configuration
- Build argument processing
- Build command execution logic
- Progress reporting
- Error handling during builds

**Dependencies**: Requires E1.2.1 (client interface)
**Size**: 500 lines (✅ Under 800 limit)

## Split Strategy Rationale

### Why This Split Makes Sense
1. **Clean Separation**: Interface/implementation vs. operations/configuration
2. **Dependency Flow**: E1.2.1a depends on E1.2.1, not vice versa
3. **Feature Completeness**: Each split provides working functionality
4. **Size Compliance**: Both splits well under 800-line limit
5. **Logical Cohesion**: Related functionality grouped together

### What Each Split Delivers
- **E1.2.1**: Working buildah client that can be created, validated, and cleaned up
- **E1.2.1a**: Full build execution capabilities with configuration support

## Implementation Sequence

### Phase 1: Complete E1.2.1 (Core Client)
1. Create split branch: `phase1/wave2/effort1-buildah-client-part1`
2. Cherry-pick types.go and client.go 
3. Ensure client can compile and basic operations work
4. Add minimal tests for client lifecycle
5. Submit for code review

### Phase 2: Complete E1.2.1a (Build Operations) 
1. Create split branch: `phase1/wave2/effort1-buildah-client-part2`
2. Cherry-pick config.go and build.go
3. Add imports/dependencies from E1.2.1
4. Ensure build operations work end-to-end
5. Add comprehensive tests for build scenarios
6. Submit for code review

## File Migration Plan

### E1.2.1 Branch Contents
```
pkg/buildah/
├── types.go           # Complete interface definitions
├── client.go          # Complete client implementation  
└── client_test.go     # Basic client tests (new)
```

### E1.2.1a Branch Contents  
```
pkg/buildah/
├── config.go          # Configuration integration
├── build.go           # Build execution logic
├── config_test.go     # Configuration tests (new)
└── build_test.go      # Build execution tests (new)
```

## Testing Strategy

### E1.2.1 Testing (estimated 80-100 lines)
- Client creation/initialization tests
- Binary detection tests
- Cleanup and lifecycle tests
- Error handling tests

### E1.2.1a Testing (estimated 200-300 lines)
- Configuration parsing tests
- Build argument construction tests  
- Build execution tests (with mocks)
- Error handling during builds
- Integration tests with configuration

## Success Criteria

### E1.2.1 Complete When:
- [x] Client interface fully defined
- [x] Client can be created and initialized
- [x] Client can detect buildah binary
- [x] Client resources can be cleaned up
- [ ] Basic tests passing
- [ ] Size under 800 lines (✅ 360 + ~100 tests = ~460 lines)

### E1.2.1a Complete When:
- [x] Configuration integration working
- [x] Build commands can be executed
- [x] Registry configuration applied
- [x] Build arguments processed correctly
- [ ] Comprehensive tests passing
- [ ] Size under 800 lines (✅ 500 + ~300 tests = ~800 lines exactly)

## Integration Points

### E1.2.1 → E1.2.1a
- E1.2.1a imports and uses Client interface from E1.2.1
- Build operations extend client capabilities
- No circular dependencies

### Both → Existing Codebase
- Both implement builder interfaces from Wave 1
- Both integrate with v1alpha1.BuildCustomizationSpec
- Both follow project logging and error patterns

## Risk Assessment

### Low Risk
- Clean logical separation
- No circular dependencies  
- Both splits deliver working functionality
- Well under size limits

### Mitigation
- E1.2.1 must be merged before E1.2.1a work begins
- Integration tests in E1.2.1a will verify the split works correctly

## Completion Timeline

This split strategy ensures continuous progress:
1. ✅ Core implementation complete (current state)
2. 🔄 Split implementation (this plan)  
3. ⏭️ Sequential reviews and merges
4. ✅ Full buildah client functionality delivered

Both splits maintain feature completeness and stay well within size limits while delivering working functionality at each stage.