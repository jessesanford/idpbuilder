# Split Plan 004 - Tests and Integration

## Split Metadata
- **Split Number**: 004
- **Parent Effort**: fallback-strategies
- **Original Branch**: idpbuilder-oci-mvp/phase1/wave2/fallback-strategies
- **Target Size**: 576 lines (max 800)
- **Created**: 2025-08-29 14:23:00

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: Split 003 of phase1/wave2/fallback-strategies
  - Path: efforts/phase1/wave2/fallback-strategies/split-003/
  - Branch: phase1/wave2/fallback-strategies-split-003
  - Summary: Core recovery strategies and manager
- **This Split**: Split 004 of phase1/wave2/fallback-strategies
  - Path: efforts/phase1/wave2/fallback-strategies/split-004/
  - Branch: phase1/wave2/fallback-strategies-split-004
- **Next Split**: None (final split)

## Implementation Scope

### Files to Create/Modify
1. Complete `pkg/certs/recovery.go` (remaining ~323 lines of 723 total)
   - Advanced recovery strategies
   - Complex recovery scenarios
   - Wave 1 integration code
2. Complete `pkg/certs/recovery_test.go` (remaining ~198 lines of 448 total)
   - Integration tests
   - Complex scenario tests
3. Wave 1 integration tests (~55 lines)
   - End-to-end validation
   - Compatibility tests

### Functionality to Implement
- Advanced recovery strategies:
  - Adaptive recovery (learns from failures)
  - Cascading recovery (multiple attempts)
  - Hybrid recovery (combines strategies)
- Complex recovery scenarios
- Wave 1 integration layer
- Performance optimizations
- Comprehensive test coverage
- Documentation and examples

### Excluded from This Split
- None (this completes the implementation)

## Technical Requirements

### Dependencies
- External dependencies:
  - All from previous splits
  - github.com/stretchr/testify (for advanced testing)
- From previous splits:
  - All interfaces from Split-001, 002, 003
  - RecoveryManager from Split-003
  - FallbackChain from Split-001
  - InsecureMode from Split-002

### Interfaces to Provide
```go
// AdaptiveRecovery - learns from past failures
type AdaptiveRecovery interface {
    RecoveryStrategy
    Learn(failure RecoveryFailure)
    GetSuccessRate() float64
}

// CascadingRecovery - tries multiple approaches
type CascadingRecovery interface {
    RecoveryStrategy
    AddStage(strategy RecoveryStrategy)
    GetStageResults() []StageResult
}

// Wave1Integration - compatibility layer
type Wave1Integration interface {
    WrapValidator(validator ChainValidator) ChainValidator
    EnableFallback(chain FallbackChain)
    EnableRecovery(manager RecoveryManager)
}
```

### Interfaces to Consume
- All interfaces from Splits 001-003
- Wave 1 ChainValidator interface

## Implementation Instructions

### Step 1: Setup
1. Verify you're in the split-004 directory
2. Confirm branch is `phase1/wave2/fallback-strategies-split-004`
3. Copy/merge all code from Splits 001-003 as base

### Step 2: Implementation
1. Complete recovery.go with advanced strategies:
   - AdaptiveRecovery implementation
   - CascadingRecovery implementation
   - HybridRecovery implementation
2. Add performance optimizations:
   - Caching mechanisms
   - Parallel recovery attempts
   - Resource pooling
3. Implement Wave 1 integration layer
4. Complete all test files
5. Add integration tests
6. Create usage examples

### Step 3: Testing
- Test adaptive learning
- Test cascading recovery
- Test hybrid approaches
- Test Wave 1 compatibility
- Performance benchmarks
- End-to-end scenarios
- Stress testing

### Step 4: Integration
- Final integration with all components
- Wave 1 compatibility verification
- Documentation completion
- Example code creation
- Performance validation

## Size Management
- Target: 576 lines
- Buffer: 224 lines (implement up to 576 lines)
- Measurement: Use line-counter.sh before committing
- This is the smallest split, well under limit

## Success Criteria
- [ ] All advanced strategies implemented
- [ ] Size under 800 lines (measured with line-counter.sh)
- [ ] 100% test coverage achieved
- [ ] Wave 1 integration working
- [ ] Performance benchmarks passing
- [ ] All previous splits integrated
- [ ] Documentation complete

## Notes for SW Engineer
- This completes the fallback-strategies implementation
- Focus on integration quality
- Ensure all edge cases are tested
- Validate Wave 1 compatibility thoroughly
- Document performance characteristics
- Create clear usage examples
- Consider future extensibility
- Clean up any technical debt from splits