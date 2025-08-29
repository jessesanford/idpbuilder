# Split Plan 003 - Recovery Strategies

## Split Metadata
- **Split Number**: 003
- **Parent Effort**: fallback-strategies
- **Original Branch**: idpbuilder-oci-mvp/phase1/wave2/fallback-strategies
- **Target Size**: 650 lines (max 800)
- **Created**: 2025-08-29 14:22:00

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: Split 002 of phase1/wave2/fallback-strategies
  - Path: efforts/phase1/wave2/fallback-strategies/split-002/
  - Branch: phase1/wave2/fallback-strategies-split-002
  - Summary: Insecure mode implementation and bypass logic
- **This Split**: Split 003 of phase1/wave2/fallback-strategies
  - Path: efforts/phase1/wave2/fallback-strategies/split-003/
  - Branch: phase1/wave2/fallback-strategies-split-003
- **Next Split**: Split 004 of phase1/wave2/fallback-strategies
  - Path: efforts/phase1/wave2/fallback-strategies/split-004/
  - Branch: phase1/wave2/fallback-strategies-split-004

## Implementation Scope

### Files to Create/Modify
1. `pkg/certs/recovery.go` (first ~400 lines of 723 total)
   - Core recovery manager interface
   - Basic recovery strategies
   - Retry mechanisms
2. `pkg/certs/recovery_test.go` (first ~250 lines of 448 total)
   - Unit tests for recovery manager
   - Tests for basic strategies

### Functionality to Implement
- Recovery manager interface and core implementation
- Retry strategy with exponential backoff
- Circuit breaker pattern implementation
- Health check mechanisms
- Recovery state management
- Error classification for recovery decisions
- Basic integration with fallback chain

### Excluded from This Split
- Advanced recovery strategies (Split-004)
- Complete test coverage (Split-004)
- Wave 1 integration tests (Split-004)
- Complex recovery scenarios (Split-004)

## Technical Requirements

### Dependencies
- External dependencies:
  - time package (for retry delays)
  - context package (for cancellation)
  - sync package (for state management)
- From previous splits:
  - FallbackStrategy from Split-001
  - FallbackChain from Split-001
  - InsecureMode from Split-002

### Interfaces to Provide
```go
// RecoveryManager - manages recovery strategies
type RecoveryManager interface {
    RegisterStrategy(name string, strategy RecoveryStrategy)
    Recover(ctx context.Context, err error, input *ValidationInput) (*ValidationResult, error)
    GetState() RecoveryState
}

// RecoveryStrategy - individual recovery approach
type RecoveryStrategy interface {
    CanRecover(err error) bool
    Attempt(ctx context.Context, input *ValidationInput) (*ValidationResult, error)
    GetMetrics() RecoveryMetrics
}

// RecoveryState - current recovery system state
type RecoveryState struct {
    Healthy bool
    LastRecovery time.Time
    FailureCount int
    CircuitOpen bool
}
```

### Interfaces to Consume
- FallbackStrategy from Split-001
- ValidationInput/Result from Split-001
- InsecureMode from Split-002 (optional integration)

## Implementation Instructions

### Step 1: Setup
1. Verify you're in the split-003 directory
2. Confirm branch is `phase1/wave2/fallback-strategies-split-003`
3. Copy/merge code from Split-001 and Split-002 as base

### Step 2: Implementation
1. Create RecoveryManager interface (recovery.go)
2. Implement core recovery manager with:
   - Strategy registration
   - Recovery orchestration
   - State management
3. Implement basic recovery strategies:
   - RetryStrategy (with exponential backoff)
   - CircuitBreakerStrategy
   - HealthCheckStrategy
4. Add error classification system
5. Implement recovery metrics collection
6. Create first portion of tests (recovery_test.go)

### Step 3: Testing
- Test recovery manager operations
- Test retry with backoff
- Test circuit breaker behavior
- Test health check mechanisms
- Test error classification
- Verify state management

### Step 4: Integration
- Integrate with FallbackChain
- Connect with InsecureMode when appropriate
- Document recovery patterns
- Prepare for advanced strategies (Split-004)

## Size Management
- Target: 650 lines
- Buffer: 150 lines (implement up to 650 lines)
- Measurement: Use line-counter.sh before committing
- Split recovery.go at logical boundary (~400 lines)

## Success Criteria
- [ ] Recovery manager implemented
- [ ] Basic strategies functional
- [ ] Size under 800 lines (measured with line-counter.sh)
- [ ] Circuit breaker working
- [ ] Retry logic tested
- [ ] State management functional
- [ ] Integration points ready

## Notes for SW Engineer
- Focus on core recovery patterns first
- Keep recovery strategies modular
- Implement proper timeout handling
- Use context for cancellation support
- Document recovery decision logic
- Consider observability needs
- Prepare hooks for advanced strategies
- Split recovery.go at a logical boundary (e.g., after basic strategies)