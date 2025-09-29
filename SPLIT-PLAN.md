# Split Plan for E1.1.2-unit-test-framework

## Overview
**Effort**: E1.1.2-unit-test-framework
**Current Size**: 836 lines
**Excess**: 36 lines over 800-line limit
**Split Strategy**: Decompose into 2 focused splits
**Sole Planner**: Code Reviewer Agent (R199 compliance)

## 📊 Current Size Breakdown
```
File Analysis:
- mock_registry.go:        216 lines
- test_helpers.go:         241 lines
- assertions.go:           226 lines
- framework_test.go:       220 lines (test file)
- framework_usage_test.go: 179 lines (test file)

Total Implementation:      836 lines (excludes test files in count)
```

## Split Strategy

### Split 001: Core Mock Registry Infrastructure
**Target Size**: ~400 lines
**Branch**: phase1/wave1/unit-test-framework-split-001
**Base Branch**: phase1/wave1/analyze-existing-structure (same as original)

#### Files to Include:
1. **mock_registry.go** (216 lines)
   - MockRegistry type
   - AuthConfig type
   - MockAuthTransport type
   - NewMockRegistry function
   - Core registry handlers

2. **Basic test_helpers.go** (~150 lines - REDUCED from 241)
   - TestFixtures type
   - SetupTestFixtures function
   - CleanupTestFixtures function
   - Remove CreateTestImage (moves to Split 002)

3. **Basic framework validation test** (~50 lines)
   - TestMockRegistryCreation
   - TestAuthTransport

**Total Estimate**: ~416 lines

#### Implementation Instructions:
1. Copy full mock_registry.go
2. Create simplified test_helpers.go with only fixture management
3. Write minimal tests to validate mock registry works
4. Ensure independent compilation and testing

### Split 002: Test Utilities and Assertions
**Target Size**: ~400 lines
**Branch**: phase1/wave1/unit-test-framework-split-002
**Base Branch**: phase1/wave1/unit-test-framework-split-001 (CASCADE from Split 001)

#### Files to Include:
1. **Image creation utilities** (~100 lines)
   - CreateTestImage function
   - testImage type implementation
   - Image helper methods

2. **assertions.go** (226 lines - keep full)
   - PushTestCase type
   - AssertPushSucceeds function
   - AssertAuthRequired function
   - All assertion helpers

3. **Comprehensive tests** (~100 lines)
   - TestImageCreation
   - TestAssertions
   - Framework usage examples

**Total Estimate**: ~426 lines

#### Implementation Instructions:
1. Import types from Split 001
2. Add CreateTestImage and related image utilities
3. Include full assertions.go
4. Write tests that use both splits together

## 🔗 Dependencies Between Splits

### Split 001 → Split 002 Flow:
```
Split 001 (Core Infrastructure):
  └── Exports: MockRegistry, AuthConfig, TestFixtures

Split 002 (Utilities & Assertions):
  └── Imports: Types from Split 001
  └── Adds: Image creation, assertions, integration
```

## Verification Checklist

### For Split 001:
- [ ] Under 800 lines (target: ~416)
- [ ] Compiles independently
- [ ] Tests pass in isolation
- [ ] Mock registry fully functional
- [ ] Can be merged to main alone

### For Split 002:
- [ ] Under 800 lines (target: ~426)
- [ ] Properly imports Split 001 types
- [ ] All assertions working
- [ ] Tests demonstrate full framework
- [ ] Can be merged after Split 001

## CASCADE Branching Plan (R501/R509)

```
main
  └── phase1/wave1/analyze-existing-structure (E1.1.1)
      └── phase1/wave1/unit-test-framework-split-001
          └── phase1/wave1/unit-test-framework-split-002
              └── phase1/wave1/integration-test-setup (E1.1.3)
```

## Implementation Sequence

### Phase 1: Split 001 Implementation
1. SW Engineer spawned for Split 001
2. Implements core mock registry
3. Size verified < 800 lines
4. Code Reviewer validates
5. Ready for integration

### Phase 2: Split 002 Implementation
1. SW Engineer spawned for Split 002
2. Builds on Split 001 foundation
3. Adds utilities and assertions
4. Size verified < 800 lines
5. Code Reviewer validates
6. Ready for integration

### Phase 3: Integration
1. Split 001 merged first
2. Split 002 merged second
3. E1.1.3 can proceed with complete framework

## Risk Mitigation

### Size Control:
- Each split has 380-line buffer below limit
- Clear file boundaries defined
- No scope creep allowed

### Functionality Preservation:
- All original functionality maintained
- No features lost in splitting
- Tests ensure correctness

### Integration Safety:
- Splits follow CASCADE pattern
- Each split independently mergeable
- No circular dependencies

## Success Criteria
- [ ] Both splits under 800 lines
- [ ] Combined functionality equals original
- [ ] All tests passing (93%+ coverage maintained)
- [ ] Independent mergeability verified
- [ ] CASCADE branching correct

## Next Steps
1. Orchestrator creates split infrastructure
2. Spawn SW Engineer for Split 001
3. After Split 001 complete, spawn for Split 002
4. Sequential implementation ensures dependency order
5. Both splits reviewed independently

## Notes
- This plan follows R199: Single reviewer handles all splits
- CASCADE pattern maintained per R501/R509
- Each split can merge independently per R307
- No deletion of functionality per R359
- Production-ready code maintained per R355