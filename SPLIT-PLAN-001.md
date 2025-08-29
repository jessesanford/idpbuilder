# Split Plan 001 - Core Certificate Validation

## Split Metadata
- **Split Number**: 001
- **Parent Effort**: certificate-validation
- **Original Branch**: idpbuilder-oci-mvp/phase1/wave2/certificate-validation
- **Target Size**: 700 lines (max 800)
- **Created**: 2025-08-29 14:12:00

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: None (first split of THIS effort)
  - Path: N/A (this is Split 001)
  - Branch: N/A
- **This Split**: Split 001 of phase1/wave2/certificate-validation
  - Path: efforts/phase1/wave2/certificate-validation/split-001/
  - Branch: phase1/wave2/certificate-validation-split-001
- **Next Split**: Split 002 of phase1/wave2/certificate-validation
  - Path: efforts/phase1/wave2/certificate-validation/split-002/
  - Branch: phase1/wave2/certificate-validation-split-002
- **File Boundaries**:
  - This Split Start: pkg/certs/chain_validator.go
  - This Split End: pkg/certs/wave1_interfaces.go
  - Next Split Start: pkg/certs/audit/interface.go

## Implementation Scope

### Files to Create/Modify
1. `pkg/certs/chain_validator.go` (20 lines)
   - Main interface definition for certificate chain validation
2. `pkg/certs/chain_validator_impl.go` (564 lines)
   - Complete implementation of the ChainValidator interface
3. `pkg/certs/errors.go` (21 lines)
   - Error types and constants for validation failures
4. `pkg/certs/types_chain.go` (171 lines)
   - Type definitions for certificate chains and validation results
5. `pkg/certs/wave1_interfaces.go` (78 lines)
   - Compatibility interfaces for Wave 1 integration

### Functionality to Implement
- Core certificate chain validation logic
- X.509 certificate parsing and verification
- Chain building and path validation
- Trust anchor management
- Error handling for validation failures
- Wave 1 compatibility layer

### Excluded from This Split
- Audit logging infrastructure (handled in Split-002)
- Persistence mechanisms (handled in Split-002)
- Example usage code (handled in Split-002)
- Full test coverage (partial in this split, completed in Split-002)

## Technical Requirements

### Dependencies
- External dependencies:
  - Standard library crypto/x509
  - Standard library crypto/tls
  - github.com/pkg/errors (for error wrapping)
- From previous splits:
  - None (this is the foundational split)

### Interfaces to Provide
```go
// ChainValidator - main validation interface
type ChainValidator interface {
    ValidateChain(certs []*x509.Certificate) (*ValidationResult, error)
    SetTrustAnchors(anchors []*x509.Certificate) error
}

// ValidationResult - validation outcome
type ValidationResult struct {
    Valid bool
    Chain []*x509.Certificate
    Errors []ValidationError
}
```

### Interfaces to Consume
- Wave 1 interfaces from existing codebase (read-only compatibility)

## Implementation Instructions

### Step 1: Setup
1. Verify you're in the split-001 directory (not the too-large directory)
2. Confirm branch is `phase1/wave2/certificate-validation-split-001`
3. Ensure clean working directory

### Step 2: Implementation
1. Create the pkg/certs directory structure
2. Implement the ChainValidator interface first (chain_validator.go)
3. Define error types and constants (errors.go)
4. Create type definitions (types_chain.go)
5. Implement the core validation logic (chain_validator_impl.go)
6. Add Wave 1 compatibility layer (wave1_interfaces.go)
7. Write basic unit tests (subset of chain_validator_test.go)

### Step 3: Testing
- Write unit tests for core validation logic
- Test error handling paths
- Verify Wave 1 interface compatibility
- Ensure all tests pass independently

### Step 4: Integration
- Verify the implementation compiles
- Check that interfaces are properly exported
- Document public APIs
- Prepare for Split-002 integration

## Size Management
- Target: 700 lines
- Buffer: 100 lines (implement up to 700 lines)
- Measurement: Use line-counter.sh before committing
- Critical files: chain_validator_impl.go is the largest at 564 lines

## Success Criteria
- [ ] All specified files implemented
- [ ] Size under 800 lines (measured with line-counter.sh)
- [ ] Core validation tests passing
- [ ] Wave 1 interfaces compatible
- [ ] No functionality regression
- [ ] Clean compilation without warnings
- [ ] Proper error handling implemented

## Notes for SW Engineer
- Focus on core validation logic first
- Keep audit hooks minimal (just interface points)
- Ensure thread-safety in the implementation
- Use standard library crypto packages where possible
- Document any assumptions about certificate formats
- Maintain compatibility with existing Wave 1 code
- Do NOT implement audit logging in this split