# Split Plan 001 - Core Fallback Mechanisms

## Split Metadata
- **Split Number**: 001
- **Parent Effort**: fallback-strategies
- **Original Branch**: idpbuilder-oci-mvp/phase1/wave2/fallback-strategies
- **Target Size**: 650 lines (max 800)
- **Created**: 2025-08-29 14:20:00

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: None (first split of THIS effort)
  - Path: N/A (this is Split 001)
  - Branch: N/A
- **This Split**: Split 001 of phase1/wave2/fallback-strategies
  - Path: efforts/phase1/wave2/fallback-strategies/split-001/
  - Branch: phase1/wave2/fallback-strategies-split-001
- **Next Split**: Split 002 of phase1/wave2/fallback-strategies
  - Path: efforts/phase1/wave2/fallback-strategies/split-002/
  - Branch: phase1/wave2/fallback-strategies-split-002
- **File Boundaries**:
  - This Split Start: pkg/certs/fallback.go (line 1)
  - This Split End: pkg/certs/fallback_test.go (line ~130)
  - Next Split Start: pkg/certs/insecure.go

## Implementation Scope

### Files to Create/Modify
1. `pkg/certs/fallback.go` (520 lines)
   - Core fallback strategy interfaces and implementation
   - Strategy pattern for certificate validation fallbacks
   - Chain of responsibility for fallback execution
2. `pkg/certs/fallback_test.go` (partial ~130 lines)
   - Basic unit tests for core fallback functionality
   - Strategy pattern tests
   - Error handling tests

### Functionality to Implement
- Fallback strategy interface definition
- Primary fallback implementation
- Secondary fallback implementation
- Fallback chain manager
- Error aggregation and reporting
- Configuration management for fallback policies
- Basic integration with Wave 1 certificate validator

### Excluded from This Split
- Insecure mode handling (Split-002)
- Recovery mechanisms (Split-003)
- Complete test coverage (Split-004)
- Advanced Wave 1 integration (Split-004)

## Technical Requirements

### Dependencies
- External dependencies:
  - github.com/pkg/errors (for error wrapping)
  - Standard library crypto/x509
- From previous splits:
  - None (this is the foundational split)
- From Wave 1:
  - Certificate validator interfaces
  - Chain validator from certificate-validation effort

### Interfaces to Provide
```go
// FallbackStrategy - main fallback interface
type FallbackStrategy interface {
    Name() string
    Priority() int
    CanHandle(err error) bool
    Execute(ctx context.Context, input *ValidationInput) (*ValidationResult, error)
}

// FallbackChain - manages multiple strategies
type FallbackChain interface {
    AddStrategy(strategy FallbackStrategy)
    Execute(ctx context.Context, input *ValidationInput) (*ValidationResult, error)
}

// ValidationInput - input for fallback validation
type ValidationInput struct {
    Certificates []*x509.Certificate
    Options      map[string]interface{}
}
```

### Interfaces to Consume
- Wave 1 ChainValidator interface
- Wave 1 ValidationResult types

## Implementation Instructions

### Step 1: Setup
1. Verify you're in the split-001 directory
2. Confirm branch is `phase1/wave2/fallback-strategies-split-001`
3. Ensure clean working directory

### Step 2: Implementation
1. Create the pkg/certs directory if needed
2. Implement FallbackStrategy interface (fallback.go)
3. Create concrete strategy implementations:
   - PrimaryStrategy (strict validation)
   - SecondaryStrategy (relaxed validation)
   - TertiaryStrategy (minimal validation)
4. Implement FallbackChain manager
5. Add configuration management
6. Create error aggregation utilities
7. Write basic unit tests

### Step 3: Testing
- Test each strategy independently
- Test chain execution order
- Test error handling and aggregation
- Verify priority ordering works
- Test configuration loading

### Step 4: Integration
- Ensure interfaces are properly exported
- Document strategy patterns used
- Prepare hooks for insecure mode (Split-002)
- Create integration points for recovery (Split-003)

## Size Management
- Target: 650 lines
- Buffer: 150 lines (implement up to 650 lines)
- Measurement: Use line-counter.sh before committing
- Critical file: fallback.go at 520 lines

## Success Criteria
- [ ] All fallback interfaces implemented
- [ ] Size under 800 lines (measured with line-counter.sh)
- [ ] Basic tests passing
- [ ] Strategy pattern working correctly
- [ ] Chain of responsibility functional
- [ ] Error aggregation working
- [ ] Clean compilation

## Notes for SW Engineer
- Focus on the strategy pattern implementation
- Keep interfaces flexible for future extensions
- Document each strategy's purpose clearly
- Ensure thread-safety in chain execution
- Prepare hooks for insecure mode but don't implement
- Use clear naming for different fallback levels
- Consider performance implications of chain execution