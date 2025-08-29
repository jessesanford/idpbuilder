# Split Plan 002 - Insecure Mode Implementation

## Split Metadata
- **Split Number**: 002
- **Parent Effort**: fallback-strategies
- **Original Branch**: idpbuilder-oci-mvp/phase1/wave2/fallback-strategies
- **Target Size**: 650 lines (max 800)
- **Created**: 2025-08-29 14:21:00

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: Split 001 of phase1/wave2/fallback-strategies
  - Path: efforts/phase1/wave2/fallback-strategies/split-001/
  - Branch: phase1/wave2/fallback-strategies-split-001
  - Summary: Core fallback mechanisms and strategy pattern
- **This Split**: Split 002 of phase1/wave2/fallback-strategies
  - Path: efforts/phase1/wave2/fallback-strategies/split-002/
  - Branch: phase1/wave2/fallback-strategies-split-002
- **Next Split**: Split 003 of phase1/wave2/fallback-strategies
  - Path: efforts/phase1/wave2/fallback-strategies/split-003/
  - Branch: phase1/wave2/fallback-strategies-split-003

## Implementation Scope

### Files to Create/Modify
1. `pkg/certs/insecure.go` (179 lines)
   - Insecure mode bypass logic
   - Development/testing certificate handling
2. `pkg/certs/insecure_test.go` (180 lines)
   - Complete tests for insecure mode
3. Complete `pkg/certs/fallback_test.go` (add ~110 lines)
   - Integration tests with insecure mode
4. Integration helpers (~181 lines)
   - Helper functions for insecure mode integration

### Functionality to Implement
- Insecure mode interface and implementation
- Certificate bypass mechanisms
- Warning/logging for insecure operations
- Configuration flags for insecure mode
- Integration with fallback chain from Split-001
- Security warnings and documentation
- Development environment detection

### Excluded from This Split
- Recovery mechanisms (Split-003)
- Complete test coverage (Split-004)
- Advanced Wave 1 integration (Split-004)

## Technical Requirements

### Dependencies
- External dependencies:
  - github.com/sirupsen/logrus (for security warnings)
  - os package (for environment detection)
- From previous splits:
  - FallbackStrategy interface from Split-001
  - FallbackChain from Split-001
  - ValidationInput/Result types from Split-001

### Interfaces to Provide
```go
// InsecureMode - bypass interface for development
type InsecureMode interface {
    IsEnabled() bool
    ShouldBypass(err error) bool
    LogSecurityWarning(context string)
}

// InsecureStrategy - fallback strategy for insecure mode
type InsecureStrategy struct {
    FallbackStrategy
    mode InsecureMode
}

// InsecureConfig - configuration for insecure operations
type InsecureConfig struct {
    Enabled bool
    LogWarnings bool
    AllowProduction bool
}
```

### Interfaces to Consume
- FallbackStrategy from Split-001
- FallbackChain from Split-001
- ValidationInput/Result from Split-001

## Implementation Instructions

### Step 1: Setup
1. Verify you're in the split-002 directory
2. Confirm branch is `phase1/wave2/fallback-strategies-split-002`
3. Copy/merge code from Split-001 as base

### Step 2: Implementation
1. Create insecure mode interface (insecure.go)
2. Implement InsecureMode with:
   - Environment detection (dev/staging/prod)
   - Bypass decision logic
   - Security warning system
3. Create InsecureStrategy as FallbackStrategy
4. Add configuration management
5. Integrate with FallbackChain from Split-001
6. Complete fallback_test.go with insecure tests
7. Write comprehensive insecure_test.go

### Step 3: Testing
- Test insecure mode enable/disable
- Test bypass decision logic
- Test security warnings
- Test environment detection
- Test integration with fallback chain
- Verify production safeguards

### Step 4: Integration
- Integrate InsecureStrategy into FallbackChain
- Ensure proper priority ordering
- Document security implications
- Add clear warnings in code
- Prepare for recovery integration (Split-003)

## Size Management
- Target: 650 lines
- Buffer: 150 lines (implement up to 650 lines)
- Measurement: Use line-counter.sh before committing
- Files well-distributed in size

## Success Criteria
- [ ] Insecure mode fully implemented
- [ ] Size under 800 lines (measured with line-counter.sh)
- [ ] All tests passing
- [ ] Security warnings functional
- [ ] Environment detection working
- [ ] Integration with Split-001 complete
- [ ] Production safeguards in place

## Notes for SW Engineer
- **SECURITY CRITICAL**: Add clear warnings everywhere
- Never enable insecure mode in production by default
- Log all insecure operations for audit
- Make bypass decisions transparent
- Consider adding metrics for insecure usage
- Document when insecure mode is appropriate
- Ensure this can be completely disabled in production builds