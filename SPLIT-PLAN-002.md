# SPLIT-PLAN-002.md

## Split 002 of 2: Validation and Comprehensive Testing
**Planner**: Code Reviewer Agent (same for ALL splits)
**Parent Effort**: cert-extraction

<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (⚠️⚠️⚠️ CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: Split 001 of phase1/wave1/cert-extraction
  - Path: efforts/phase1/wave1/cert-extraction/split-001/
  - Branch: idpbuilder-oci-mvp/phase1/wave1/cert-extraction-split-001
  - Summary: Core extraction functionality (types, errors, extractor)
- **This Split**: Split 002 of phase1/wave1/cert-extraction
  - Path: efforts/phase1/wave1/cert-extraction/split-002/
  - Branch: idpbuilder-oci-mvp/phase1/wave1/cert-extraction-split-002
- **Next Split**: None (final split)
- **File Boundaries**:
  - Previous Split: types.go, errors.go, extractor.go
  - This Split: validator.go, all test files

## Files in This Split (EXCLUSIVE - no overlap with other splits)
- `pkg/certs/validator.go` (217 lines) - Certificate validation logic
- `pkg/certs/errors_test.go` (~150 lines, reduced from 204)
- `pkg/certs/extractor_test.go` (~250 lines, reduced from 391)
- `pkg/certs/validator_test.go` (~200 lines, reduced from 321)

**Target Total**: ~617 lines (well under 700 target)

## Functionality

### Certificate Validation Implementation
- Implement CertValidator interface from Split 001
- Certificate chain validation
- Expiry checking with configurable warning thresholds
- Hostname verification for certificates
- Diagnostic generation with detailed reporting
- Self-signed certificate support for Kind clusters

### Comprehensive Test Coverage
- Unit tests for all error types and handling
- Extractor tests with mocked Kubernetes clients
- Validator tests for all validation scenarios
- Edge case coverage (expired certs, invalid chains, etc.)
- Integration test helpers (for future use)

## Dependencies
- **From Split 001**: All types, interfaces, and extractor implementation
- External packages (same as Split 001 plus testing libs):
  - `github.com/stretchr/testify` - Test assertions
  - `k8s.io/client-go/fake` - Mock Kubernetes clients
  - Standard library testing packages

## Implementation Instructions

### Step 1: Checkout Split 001 Changes
```bash
# Ensure you have all changes from Split 001
git checkout idpbuilder-oci-mvp/phase1/wave1/cert-extraction-split-001
git pull
# Create new branch for Split 002
git checkout -b idpbuilder-oci-mvp/phase1/wave1/cert-extraction-split-002
```

### Step 2: Implement validator.go
- Create Validator struct implementing CertValidator interface
- Implement ValidateChain() for certificate chain validation
- Add CheckExpiry() with configurable thresholds
- Implement VerifyHostname() for hostname verification
- Add GenerateDiagnostics() for detailed reporting
- Include ValidateCertificate() method for KindExtractor

### Step 3: Create Optimized Test Suite

#### errors_test.go (~150 lines, reduced from 204)
- Combine similar test cases
- Focus on critical error paths
- Test error wrapping and suggestions
- Verify predefined errors

#### extractor_test.go (~250 lines, reduced from 391)
- Test NewKindExtractor with various configs
- Mock Kubernetes client for pod operations
- Test certificate parsing scenarios
- Verify storage operations
- Combine repetitive test cases

#### validator_test.go (~200 lines, reduced from 321)
- Test all validation methods
- Cover expiry scenarios
- Test self-signed certificate handling
- Verify diagnostic generation
- Optimize by using table-driven tests

### Step 4: Run Tests and Verify Coverage
```bash
cd pkg/certs
go test -v -cover
# Target: 65%+ coverage
```

### Step 5: Measure Final Size
```bash
# From split directory
$PROJECT_ROOT/tools/line-counter.sh
# Must be under 700 lines (target) / 800 lines (hard limit)
```

## Test Optimization Strategy
To keep under line limits while maintaining quality:
1. Use table-driven tests where possible
2. Combine similar test scenarios
3. Share test fixtures and helpers
4. Focus on critical paths and edge cases
5. Remove redundant assertions
6. Use subtests for organization

## Split Branch Strategy
- Branch: `idpbuilder-oci-mvp/phase1/wave1/cert-extraction-split-002`
- Base: `idpbuilder-oci-mvp/phase1/wave1/cert-extraction-split-001`
- Must merge to: `idpbuilder-oci-mvp/phase1/wave1/cert-extraction` after review

## Success Criteria
- ✅ Validator fully implements CertValidator interface
- ✅ Test coverage ≥65% for entire package
- ✅ All tests passing
- ✅ Total size under 700 lines (hard limit 800)
- ✅ No functionality lost from original implementation
- ✅ Maintains code quality standards

## Integration Notes
After both splits are complete and reviewed:
1. Merge Split 001 to parent branch
2. Merge Split 002 to parent branch
3. Verify complete functionality works end-to-end
4. Total implementation will be ~1110 lines (493 + 617)
5. Each split remains well under the 800-line limit

## Notes for SW Engineer
- Start from Split 001 completed code
- Focus on test quality over quantity
- Use table-driven tests to reduce line count
- Ensure validator integrates cleanly with extractor
- Maintain high code coverage despite optimizations