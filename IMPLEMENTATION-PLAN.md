# E1.2.1: Certificate Validation Pipeline - Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Effort ID**: E1.2.1  
**Branch**: `phase1/wave2/cert-validation`  
**Can Parallelize**: Yes  
**Parallel With**: [E1.2.2 - Fallback Strategies]  
**Size Estimate**: 400 lines  
**Dependencies**: [E1.1.1 - Kind Certificate Extraction, E1.1.2 - Registry TLS Trust]  
**Feature Flag**: `CERT_VALIDATION_ENABLED`  

## Overview
- **Effort**: Certificate validation pipeline to validate cert chains and provide diagnostics
- **Phase**: 1, Wave: 2
- **Estimated Size**: 400 lines (target: 350 lines with 50-line buffer)
- **Implementation Time**: 6-8 hours

## Technical Architecture

### Core Components
1. **Certificate Validator**: Main validation engine implementing the `CertificateValidator` interface
2. **Chain Validator**: Validates complete certificate chains from root to leaf
3. **Diagnostics Generator**: Produces detailed diagnostic reports for troubleshooting
4. **Validation Modes**: Support for STRICT, LENIENT, and PERMISSIVE validation
5. **Error Aggregator**: Collects and categorizes validation errors

### Integration Points
- **Wave 1 Dependencies**:
  - `KindCertExtractor` from E1.1.1 for certificate retrieval
  - `TrustStoreManager` from E1.1.2 for trusted root access
  - `CertificateStorage` from E1.1.1 for cert persistence
- **External Libraries**:
  - Standard library `crypto/x509` for certificate operations
  - Standard library `crypto/tls` for TLS validation
  - `encoding/pem` for PEM format handling

## File Structure

### Production Code (350 lines target)
- `pkg/certs/validator.go` (120 lines)
  - Core `CertificateValidator` interface implementation
  - Main validation orchestration logic
  - Validation mode management
  
- `pkg/certs/chain_validator.go` (100 lines)
  - Certificate chain validation logic
  - Root CA trust verification
  - Intermediate certificate handling
  
- `pkg/certs/diagnostics.go` (80 lines)
  - `CertDiagnostics` struct and methods
  - Diagnostic report generation
  - Human-readable error formatting
  
- `pkg/certs/validation_errors.go` (50 lines)
  - Custom error types for validation failures
  - Error categorization (EXPIRED, UNTRUSTED, INVALID_CHAIN, etc.)
  - Error severity levels

### Test Code (400 lines target)
- `pkg/certs/validator_test.go` (150 lines)
  - Unit tests for core validation logic
  - Mock certificate generation
  - Edge case testing
  
- `pkg/certs/chain_validator_test.go` (100 lines)
  - Chain validation test scenarios
  - Self-signed certificate tests
  - Trust chain tests
  
- `pkg/certs/diagnostics_test.go` (100 lines)
  - Diagnostic generation tests
  - Error message formatting tests
  - Report completeness tests
  
- `pkg/certs/validation_errors_test.go` (50 lines)
  - Error type validation
  - Error categorization tests
  - Severity level tests

## Implementation Steps

### Step 1: Core Validator Interface (120 lines)
1. Create `pkg/certs/validator.go`
2. Define `CertificateValidator` interface:
   ```go
   type CertificateValidator interface {
       ValidateChain(certs []*x509.Certificate) error
       ValidateCertificate(cert *x509.Certificate) error
       VerifyHostname(cert *x509.Certificate, hostname string) error
       GenerateDiagnostics() (*CertDiagnostics, error)
       SetValidationMode(mode ValidationMode)
   }
   ```
3. Implement `DefaultCertificateValidator` struct
4. Add validation mode support (STRICT, LENIENT, PERMISSIVE)
5. Integrate with Wave 1's `TrustStoreManager` for root CA access
6. **Line estimate**: 120 lines

### Step 2: Chain Validation Logic (100 lines)
1. Create `pkg/certs/chain_validator.go`
2. Implement chain building from leaf to root
3. Add intermediate certificate handling
4. Verify each certificate in chain:
   - Signature validation
   - Time validity
   - Key usage constraints
   - Basic constraints
5. Support partial chain validation
6. **Line estimate**: 100 lines

### Step 3: Diagnostics System (80 lines)
1. Create `pkg/certs/diagnostics.go`
2. Define `CertDiagnostics` struct:
   ```go
   type CertDiagnostics struct {
       Subject          string
       Issuer           string
       NotBefore        time.Time
       NotAfter         time.Time
       ChainLength      int
       ValidationErrors []string
       IsExpired        bool
       IsSelfSigned     bool
   }
   ```
3. Implement diagnostic collection during validation
4. Add human-readable formatting methods
5. Include recommendations for fixing issues
6. **Line estimate**: 80 lines

### Step 4: Error Handling System (50 lines)
1. Create `pkg/certs/validation_errors.go`
2. Define custom error types:
   - `ExpiredCertError`
   - `UntrustedRootError`
   - `InvalidChainError`
   - `HostnameMismatchError`
3. Add error categorization and severity
4. Implement error aggregation for multiple issues
5. **Line estimate**: 50 lines

### Step 5: Unit Tests (400 lines)
1. Create comprehensive test suites for each component
2. Generate test certificates with various issues:
   - Expired certificates
   - Self-signed certificates
   - Invalid chains
   - Wrong hostname
3. Test all validation modes
4. Verify diagnostic output accuracy
5. Test error handling and aggregation
6. **Line estimate**: 400 lines

### Step 6: Integration Verification
1. Test with Wave 1's `KindCertExtractor` output
2. Verify compatibility with `TrustStoreManager`
3. Ensure feature flag controls activation
4. Document integration points

## Dependencies from Wave 1

### From E1.1.1 (Kind Certificate Extraction)
- Import `github.com/jessesanford/idpbuilder/pkg/certs` for:
  - `KindCertExtractor` (for testing integration)
  - `CertificateStorage` interface
  - Certificate retrieval utilities

### From E1.1.2 (Registry TLS Trust)
- Import for:
  - `TrustStoreManager` interface
  - Root CA certificate access
  - Trust configuration

## Size Management Strategy

### Line Count Tracking
- **Target**: 350 lines of production code
- **Buffer**: 50 lines for unexpected complexity
- **Hard Limit**: 400 lines (800 line limit / 2 for safety)

### Measurement Protocol
```bash
# Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do 
    [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Measure after each component
$PROJECT_ROOT/tools/line-counter.sh

# Check points:
# - After validator.go: expect ~120 lines
# - After chain_validator.go: expect ~220 lines
# - After diagnostics.go: expect ~300 lines
# - After validation_errors.go: expect ~350 lines
```

### Split Contingency
If approaching 400 lines:
1. Move diagnostics to separate effort
2. Simplify error handling
3. Defer advanced validation modes

## Testing Requirements

### Unit Test Coverage
- **Target**: 85% code coverage
- **Minimum**: 80% code coverage
- All public methods must have tests
- Error paths must be tested

### Test Scenarios
1. **Valid Certificates**
   - Valid chain with trusted root
   - Valid self-signed certificate
   - Valid certificate with SANs

2. **Invalid Certificates**
   - Expired certificate
   - Not yet valid certificate
   - Untrusted root CA
   - Broken chain
   - Wrong hostname
   - Invalid signature

3. **Edge Cases**
   - Empty certificate list
   - Nil certificate
   - Circular chain
   - Multiple validation errors

### Integration Tests
- Validate real certificates from Kind cluster
- Test with go-containerregistry integration
- Verify feature flag activation/deactivation

## Pattern Compliance

### Code Style
- Follow Go idioms and best practices
- Use meaningful variable names
- Keep functions focused and small
- Document all exported types and methods

### Error Handling
- Return errors, don't panic
- Wrap errors with context
- Use custom error types for specific failures
- Aggregate multiple validation errors

### Security Considerations
- Never accept invalid certificates by default
- Require explicit configuration for lenient modes
- Log all validation failures
- Sanitize diagnostic output to prevent information leakage

## Validation Checkpoints

### During Implementation
1. **After Step 1**: Verify interface matches phase plan requirements
2. **After Step 2**: Test with sample certificate chains
3. **After Step 3**: Verify diagnostic output is helpful
4. **After Step 4**: Ensure all error types are covered
5. **After Step 5**: Run full test suite, check coverage

### Before Completion
- [ ] All tests passing
- [ ] Code coverage ≥80%
- [ ] Line count <400 (production code)
- [ ] No compilation errors
- [ ] No critical TODOs
- [ ] Feature flag properly implemented
- [ ] Integration with Wave 1 verified

## Risk Mitigation

### Technical Risks
1. **Complexity in chain validation**: Use standard library where possible
2. **Cross-platform differences**: Test on Linux/Mac/Windows
3. **Performance with large chains**: Add caching if needed

### Size Risks
1. **Monitor line count after each component**
2. **Prepare to simplify if approaching limit**
3. **Have clear split points identified**

## Success Criteria

### Functional Requirements
- ✅ Validates certificate chains correctly
- ✅ Provides detailed diagnostics for failures
- ✅ Supports multiple validation modes
- ✅ Integrates with Wave 1 components
- ✅ Controlled by feature flag

### Non-Functional Requirements
- ✅ Performance: <50ms per validation
- ✅ Memory: <10MB for validation operations
- ✅ Security: No bypass of validation without explicit config
- ✅ Maintainability: Clear, documented code
- ✅ Testability: >80% coverage

## Implementation Notes

### For SW Engineer
1. Start with the interface definition to ensure contract is correct
2. Use test-driven development for validation logic
3. Generate test certificates using `crypto/x509` utilities
4. Check line count frequently with the designated tool
5. Commit after each major component
6. If approaching size limit, consult with Code Reviewer immediately

### Integration Considerations
- This effort can proceed in parallel with E1.2.2 (Fallback Strategies)
- Both efforts will integrate into the same `pkg/certs` package
- Coordinate on shared error types if needed
- Feature flags should be independent

---

**Plan Version**: 1.0  
**Created**: 2025-09-07 06:47:01  
**Effort**: E1.2.1 - Certificate Validation Pipeline  
**Phase**: 1, Wave: 2  
**Status**: Ready for Implementation