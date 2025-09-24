# Phase 2 Wave 1 Backport Plan

## Executive Summary

This backport plan addresses the TDD GREEN phase test failures identified during Phase 2 Wave 1 integration. The authentication tests (written in effort 2.1.1) define three constructor functions that need implementation in the auth-implementation effort (2.1.2). These constructors are essential for the authentication system to function properly and will enable the tests to pass.

### What Needs Backporting
Three constructor functions for the `DefaultAuthenticator` type that enable creation from different credential sources:
1. Secret-based initialization (Kubernetes secrets)
2. Flag-based initialization (CLI parameters)
3. Environment-based initialization (environment variables)

### Why These Changes Are Required
- Tests were written first (TDD approach) and define the expected interface
- Constructor functions are the primary mechanism for creating authenticator instances
- Current code has the base structure but lacks these convenience constructors
- Integration build succeeds but tests fail with "undefined" errors

### Expected Outcome
After implementing these constructors:
- Build will continue to pass ✅
- Tests will progress from "undefined" errors to actual test assertions
- Authentication system will be usable via multiple initialization patterns
- TDD GREEN phase will be complete, ready for RED phase refinement

## Effort-Specific Implementation Instructions

### Effort 2.1.2: Authentication Implementation

#### Branch Information
- **Branch Name**: `idpbuilderpush/phase2/wave1/auth-implementation`
- **Working Directory**: `/home/vscode/workspaces/idpbuilder-push/efforts/phase2/wave1/auth-implementation`
- **Base Branch**: `phase1/wave1/integration`

#### File to Modify
**Path**: `pkg/oci/auth.go`

#### Functions to Implement

##### 1. NewAuthenticatorFromSecrets
```go
// NewAuthenticatorFromSecrets creates an authenticator from Kubernetes secret data
// The secretData map should contain "username" and "password" keys
func NewAuthenticatorFromSecrets(secretData map[string][]byte) (*DefaultAuthenticator, error) {
    // Implementation requirements:
    // 1. Validate secretData is not nil
    // 2. Check for required keys: "username" and "password"
    // 3. Convert []byte values to strings
    // 4. Create AuthConfig with credentials
    // 5. Return NewAuthenticator with config
    // 6. Return error if validation fails
}
```

**Implementation Guidelines**:
- Extract username from `secretData["username"]`
- Extract password from `secretData["password"]`
- Return error if either key is missing or empty
- Use the existing `NewAuthenticator` function with appropriate config
- Error messages should be descriptive (e.g., "missing username in secret data")

##### 2. NewAuthenticatorFromFlags
```go
// NewAuthenticatorFromFlags creates an authenticator from CLI flag values
func NewAuthenticatorFromFlags(username, password string) (*DefaultAuthenticator, error) {
    // Implementation requirements:
    // 1. Validate username is not empty
    // 2. Validate password is not empty
    // 3. Create AuthConfig with provided credentials
    // 4. Return NewAuthenticator with config
    // 5. Return error if validation fails
}
```

**Implementation Guidelines**:
- Both parameters are required (non-empty)
- Return descriptive errors for missing/empty values
- Consider trimming whitespace from inputs
- Use the existing `NewAuthenticator` function

##### 3. NewAuthenticatorFromEnv
```go
// NewAuthenticatorFromEnv creates an authenticator from environment variables
// Reads OCI_USERNAME and OCI_PASSWORD environment variables
func NewAuthenticatorFromEnv() (*DefaultAuthenticator, error) {
    // Implementation requirements:
    // 1. Read OCI_USERNAME environment variable
    // 2. Read OCI_PASSWORD environment variable
    // 3. Validate both are present and non-empty
    // 4. Create AuthConfig with env credentials
    // 5. Return NewAuthenticator with config
    // 6. Return error if env vars missing/empty
}
```

**Implementation Guidelines**:
- Use `os.Getenv("OCI_USERNAME")` and `os.Getenv("OCI_PASSWORD")`
- Return error if either variable is not set or empty
- Error message should indicate which env var is missing
- Use the existing `NewAuthenticator` function

#### Size Constraints
- Current auth.go size: ~5603 bytes
- Estimated addition: ~50-80 lines per constructor
- Total estimated: ~150-250 lines
- **Well within effort limits** (< 800 lines)

## Implementation Sequence

### Order of Implementation
1. **NewAuthenticatorFromFlags** (simplest, direct parameters)
2. **NewAuthenticatorFromEnv** (simple env var reading)
3. **NewAuthenticatorFromSecrets** (most complex, map handling)

### Prerequisites
- Ensure you're on the correct branch (`idpbuilderpush/phase2/wave1/auth-implementation`)
- Verify pkg/oci/auth.go exists and has the base DefaultAuthenticator type
- Confirm you have the existing `NewAuthenticator` function to build upon

### Dependencies
- No external dependencies required
- Uses standard library (os, fmt, errors packages)
- Builds on existing DefaultAuthenticator structure

## Verification Steps

### For Each Constructor Implementation

#### Build Verification
```bash
# Navigate to effort directory
cd /home/vscode/workspaces/idpbuilder-push/efforts/phase2/wave1/auth-implementation

# Verify clean working directory
git status

# Build the package
go build ./pkg/oci/...

# Expected: Successful build with no errors
```

#### Test Verification
```bash
# Run the specific test for each constructor
go test -v ./pkg/oci/... -run TestNewAuthenticatorFromFlags
go test -v ./pkg/oci/... -run TestNewAuthenticatorFromEnv
go test -v ./pkg/oci/... -run TestNewAuthenticatorFromSecrets

# Run all auth tests
go test ./pkg/oci/...

# Expected: Tests should now execute (may fail on assertions, but not "undefined")
```

#### Integration Verification
```bash
# After all constructors implemented, test full integration
go test -v ./pkg/oci/auth_test.go

# Verify no regressions in existing code
go test ./...
```

### Expected Results After Implementation
1. **Build**: Continues to pass without errors
2. **Undefined Errors**: Should be resolved
3. **Test Execution**: Tests run and provide meaningful feedback
4. **No Deletions**: Verify no existing code was removed (git diff should show only additions)

## Success Criteria

### Mandatory Requirements
✅ All three constructor functions implemented in `pkg/oci/auth.go`
✅ Build succeeds without compilation errors
✅ Test failures change from "undefined: NewAuthenticator*" to actual test execution
✅ NO existing code deleted (R359 compliance)
✅ Implementation stays within 800-line limit (R304 compliance)
✅ Functions match exact signatures expected by tests

### Quality Standards
✅ Proper error handling with descriptive messages
✅ Input validation for all parameters
✅ Consistent with existing code patterns
✅ Clean, readable implementation
✅ No stub/mock code in production (R355)

## Risk Assessment

### Potential Issues
1. **Risk**: Constructor signatures might need adjustment
   - **Mitigation**: Check test file for exact usage patterns
   - **Action**: Adjust return types if needed (e.g., *DefaultAuthenticator vs Authenticator interface)

2. **Risk**: Tests might expect specific error messages
   - **Mitigation**: Review test assertions for error string matching
   - **Action**: Align error messages with test expectations

3. **Risk**: Missing helper methods or fields
   - **Mitigation**: Check if DefaultAuthenticator needs additional fields
   - **Action**: Add necessary fields to support constructor patterns

### Rollback Plan
If implementation causes issues:
1. Git reset to previous commit
2. Review test expectations more carefully
3. Implement incrementally, one constructor at a time
4. Test each constructor independently before proceeding

## Implementation Checklist for SW Engineer

### Pre-Implementation
- [ ] Navigate to correct working directory
- [ ] Checkout correct branch (`idpbuilderpush/phase2/wave1/auth-implementation`)
- [ ] Verify pkg/oci/auth.go exists
- [ ] Review test file for exact expectations
- [ ] Commit any pending changes

### Implementation
- [ ] Implement `NewAuthenticatorFromFlags`
- [ ] Test and verify flags constructor
- [ ] Implement `NewAuthenticatorFromEnv`
- [ ] Test and verify env constructor
- [ ] Implement `NewAuthenticatorFromSecrets`
- [ ] Test and verify secrets constructor

### Post-Implementation
- [ ] Run full test suite
- [ ] Verify build passes
- [ ] Check no code was deleted (git diff)
- [ ] Measure implementation size with line counter
- [ ] Commit with clear message
- [ ] Push to remote branch
- [ ] Update effort status

## Notes for SW Engineer

### Important Considerations
- This is TDD GREEN phase - tests define the contract
- Constructor functions should be simple wrappers around `NewAuthenticator`
- Focus on input validation and error handling
- Keep implementations clean and testable
- Do NOT modify test files - they are the specification

### Common Patterns to Follow
```go
// Example pattern for validation
if username == "" {
    return nil, fmt.Errorf("username cannot be empty")
}

// Example pattern for creating config
config := &AuthConfig{
    // Set appropriate fields based on input
}

// Example pattern for returning
return NewAuthenticator(config), nil
```

### Testing Tips
- Run tests frequently during implementation
- Start with one constructor at a time
- Use -v flag for verbose test output
- Check test file if unsure about expectations

## Compliance Verification

### Rule Compliance Checklist
- **R321**: ✅ Fixes target source branch, not integration
- **R359**: ✅ No deletion of existing code permitted
- **R304**: ✅ Line counter to be used for verification
- **R006**: ✅ This plan contains no implementation code
- **R355**: ✅ No stubs/mocks in production code
- **R405**: ✅ Automation flag will be set at completion

## Approval

This backport plan is **READY FOR SW ENGINEER IMPLEMENTATION**.

The instructions are clear, unambiguous, and provide all necessary information for successful implementation of the missing constructor functions.

---
*Generated by Code Reviewer Agent*
*State: BACKPORT_PLAN_CREATION*
*Date: 2025-09-24*
*Plan Version: 1.0*