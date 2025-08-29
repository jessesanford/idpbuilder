# Split Plan 002 - Audit Infrastructure

## Split Metadata
- **Split Number**: 002
- **Parent Effort**: certificate-validation
- **Original Branch**: idpbuilder-oci-mvp/phase1/wave2/certificate-validation
- **Target Size**: 650 lines (max 800)
- **Created**: 2025-08-29 14:13:00

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: Split 001 of phase1/wave2/certificate-validation
  - Path: efforts/phase1/wave2/certificate-validation/split-001/
  - Branch: phase1/wave2/certificate-validation-split-001
  - Summary: Core certificate chain validation logic and interfaces
- **This Split**: Split 002 of phase1/wave2/certificate-validation
  - Path: efforts/phase1/wave2/certificate-validation/split-002/
  - Branch: phase1/wave2/certificate-validation-split-002
- **Next Split**: None (final split)
  - Path: N/A
  - Branch: N/A

## Implementation Scope

### Files to Create/Modify
1. `pkg/certs/audit/interface.go` (77 lines)
   - Audit logger interface definitions
2. `pkg/certs/audit/logger.go` (278 lines)
   - Complete audit logger implementation
3. `pkg/certs/audit/logger_test.go` (348 lines)
   - Comprehensive tests for audit logger
4. `pkg/certs/audit/example_usage.go` (133 lines)
   - Example code demonstrating audit integration
5. Complete `pkg/certs/chain_validator_test.go` (add ~308 lines)
   - Additional integration tests with audit logging

### Functionality to Implement
- Audit logging interface and implementation
- Persistence mechanisms for audit records
- Structured logging with contextual information
- Audit event types and formatting
- Integration with certificate validation
- Comprehensive test coverage

### Excluded from This Split
- Core validation logic (already in Split-001)
- Basic validation tests (already in Split-001)
- Wave 1 interfaces (already in Split-001)

## Technical Requirements

### Dependencies
- External dependencies:
  - github.com/sirupsen/logrus (for structured logging)
  - github.com/stretchr/testify (for testing)
- From previous splits:
  - ChainValidator interface from Split-001
  - ValidationResult types from Split-001
  - Error types from Split-001

### Interfaces to Provide
```go
// AuditLogger - audit logging interface
type AuditLogger interface {
    LogValidation(result *ValidationResult, metadata map[string]interface{}) error
    LogFailure(err error, context map[string]interface{}) error
    Flush() error
}

// AuditRecord - structured audit record
type AuditRecord struct {
    Timestamp time.Time
    EventType string
    Result    *ValidationResult
    Metadata  map[string]interface{}
}
```

### Interfaces to Consume
- ChainValidator from Split-001
- ValidationResult from Split-001
- ValidationError types from Split-001

## Implementation Instructions

### Step 1: Setup
1. Verify you're in the split-002 directory
2. Confirm branch is `phase1/wave2/certificate-validation-split-002`
3. Copy/merge code from Split-001 as the base

### Step 2: Implementation
1. Create the pkg/certs/audit directory
2. Define the AuditLogger interface (interface.go)
3. Implement the audit logger with persistence (logger.go)
4. Add structured logging capabilities
5. Create example usage demonstrations (example_usage.go)
6. Write comprehensive tests (logger_test.go)
7. Complete integration tests in chain_validator_test.go

### Step 3: Testing
- Test audit logger independently
- Test persistence mechanisms
- Verify integration with validation
- Test error scenarios and edge cases
- Ensure 80%+ code coverage

### Step 4: Integration
- Integrate audit logger with ChainValidator
- Verify audit records are properly created
- Test with various certificate scenarios
- Document integration points

## Size Management
- Target: 650 lines
- Buffer: 150 lines (implement up to 650 lines)
- Measurement: Use line-counter.sh before committing
- Critical files: logger_test.go is the largest at 348 lines

## Success Criteria
- [ ] All audit files implemented
- [ ] Size under 800 lines (measured with line-counter.sh)
- [ ] Audit tests achieving 80%+ coverage
- [ ] Integration with Split-001 working
- [ ] Persistence mechanism functional
- [ ] Example usage clear and working
- [ ] No regression in validation functionality

## Notes for SW Engineer
- Build on top of Split-001's foundation
- Ensure audit logging is non-invasive
- Make audit logger pluggable/optional
- Use structured logging for better searchability
- Include correlation IDs in audit records
- Consider performance impact of persistence
- Provide clear examples of integration
- Maintain thread-safety in logger implementation