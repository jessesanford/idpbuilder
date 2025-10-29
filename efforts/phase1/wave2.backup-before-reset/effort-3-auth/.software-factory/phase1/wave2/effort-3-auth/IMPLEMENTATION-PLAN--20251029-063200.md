# Effort 1.2.3: Authentication Implementation - Implementation Plan

**Created**: 2025-10-29 06:32:00 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**Effort ID**: 1.2.3
**Phase**: Phase 1 - Foundation & Interfaces
**Wave**: Wave 2 - Core Package Implementations

---

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)

### R213 Infrastructure Metadata

```json
{
  "effort_id": "1.2.3",
  "effort_name": "Authentication Implementation",
  "branch_name": "idpbuilder-oci-push/phase1/wave2/effort-3-auth",
  "base_branch": "idpbuilder-oci-push/phase1/wave2/integration",
  "parent_wave": "WAVE_2",
  "parent_phase": "PHASE_1",
  "depends_on": [],
  "estimated_lines": 350,
  "complexity": "medium",
  "can_parallelize": true,
  "parallel_with": ["1.2.1", "1.2.2", "1.2.4"]
}
```

**Branch**: `idpbuilder-oci-push/phase1/wave2/effort-3-auth`
**Can Parallelize**: Yes
**Parallel With**: [1.2.1, 1.2.2, 1.2.4] (ALL Wave 2 efforts run simultaneously)
**Size Estimate**: 350 lines (MUST be <800)
**Dependencies**: None (all Wave 1 interfaces frozen and available)

---

## Overview

**Purpose**: Implement the basic authentication provider that validates credentials, converts them to go-containerregistry's authn.Authenticator format, and supports special characters in passwords (including quotes, spaces, unicode).

**What This Effort Accomplishes**:
- Complete implementation of `auth.Provider` interface (frozen in Wave 1)
- Basic authentication (username/password) support
- Credential validation with security checks
- Conversion to `authn.Authenticator` for go-containerregistry
- Special character support in passwords (unicode, quotes, spaces)
- Control character detection in usernames

**Boundaries - OUT OF SCOPE**:
- Token-based authentication (OAuth, bearer tokens)
- Certificate-based authentication (mTLS)
- Credential storage or management (keychain, secret stores)
- Password hashing or encryption (plaintext transmission via HTTP Basic Auth)
- Multi-factor authentication
- Session management

---

## File Structure

### New Files to Create

**Implementation Files**:
- `pkg/auth/basic.go` (~350 lines)
  - `basicAuthProvider` struct implementation
  - `NewBasicAuthProvider()` constructor
  - `GetAuthenticator()` method returning authn.Authenticator
  - `ValidateCredentials()` validation method
  - Helper function (`containsControlChars`)

**Test Files** (NOT counted in line estimates per R007):
- `pkg/auth/basic_test.go` (~250 lines)
  - 10+ test cases covering all methods
  - Special character password tests
  - Control character username tests
  - Success and failure paths

**Modified Files**:
- None (go.mod already has go-containerregistry from Wave 1)

**Total Estimated Lines**: 350 lines (implementation only, tests excluded per R007)

---

## Implementation Steps

### Step 1: Review Wave 1 Interface Definition

**MANDATORY FIRST STEP - Read frozen interfaces**:
```bash
# Read the frozen Auth interface from Wave 1
cat pkg/auth/types.go
```

**Expected Auth Interface** (from Wave 1 Effort 3):
```go
package auth

import "github.com/google/go-containerregistry/pkg/authn"

type Provider interface {
    GetAuthenticator() (authn.Authenticator, error)
    ValidateCredentials() error
}

type Credentials struct {
    Username string
    Password string
}
```

**Expected Error Types** (from Wave 1):
```go
type CredentialValidationError struct {
    Field   string
    Message string
}
```

### Step 2: Implement pkg/auth/basic.go

**File: pkg/auth/basic.go**

**Required Implementation Details**:

**1. Struct Definition**:
```go
type basicAuthProvider struct {
    credentials Credentials
}
```

**2. NewBasicAuthProvider() Implementation** (~40 lines):
- Accept `username` and `password` as string parameters
- Store in `Credentials` struct: `Credentials{Username: username, Password: password}`
- Return `&basicAuthProvider{credentials: creds}`
- NO validation in constructor (validation happens in ValidateCredentials)
- Return implementation of `Provider` interface

**3. GetAuthenticator() Implementation** (~60 lines):
- Call `ValidateCredentials()` first
- Return `CredentialValidationError` if validation fails
- Create `authn.Basic{Username: username, Password: password}` from go-containerregistry
- Return authenticator compatible with go-containerregistry
- The authenticator will handle base64 encoding of credentials

**Example**:
```go
func (b *basicAuthProvider) GetAuthenticator() (authn.Authenticator, error) {
    if err := b.ValidateCredentials(); err != nil {
        return nil, err
    }

    return &authn.Basic{
        Username: b.credentials.Username,
        Password: b.credentials.Password,
    }, nil
}
```

**4. ValidateCredentials() Implementation** (~80 lines):
- Check username is not empty:
  - If empty → `&CredentialValidationError{Field: "username", Message: "cannot be empty"}`
- Check username contains no control characters:
  - Control chars: ASCII < 32 (includes newline, tab, null) or == 127 (DEL)
  - If found → `&CredentialValidationError{Field: "username", Message: "contains control characters"}`
  - Use helper function `containsControlChars(s string) bool`
- Check password is not empty:
  - If empty → `&CredentialValidationError{Field: "password", Message: "cannot be empty"}`
- Allow ALL printable characters in password:
  - Including quotes (', ")
  - Including spaces
  - Including unicode characters (пароль, 密码, 🔒)
  - NO restrictions except non-empty
- Return nil if valid

**5. Helper Function** (~30 lines):

**containsControlChars(s string) bool**:
- Iterate through each rune in string
- Check if rune < 32 (control characters) or == 127 (DEL)
- Return true if any control character found
- Return false otherwise

**Security Considerations**:
- **Username**: No control characters (prevents terminal escape sequence attacks)
- **Password**: Allow everything (HTTP Basic Auth transmits as-is, base64 encoded)
- **No password strength requirements** (user's responsibility)
- **No credential logging or exposure** (never log passwords)

**Complete Package Structure**:
```go
package auth

import (
    "github.com/google/go-containerregistry/pkg/authn"
)

// Implementation here (~350 lines total including comments and spacing)
```

### Step 3: Write Tests (TDD - Tests First!)

**File: pkg/auth/basic_test.go**

**Test Cases** (from Wave 2 Test Plan):

**A. Constructor Tests**:
- `TestNewBasicAuthProvider_Success`: NewBasicAuthProvider creates provider

**B. GetAuthenticator Tests**:
- `TestGetAuthenticator_Success`: GetAuthenticator succeeds with valid credentials
- `TestGetAuthenticator_EmptyUsername`: GetAuthenticator fails with empty username
- `TestGetAuthenticator_EmptyPassword`: GetAuthenticator fails with empty password

**C. ValidateCredentials Tests**:
- `TestValidateCredentials_Valid_SimpleCredentials`: Passes for simple credentials (`"user"` / `"pass"`)
- `TestValidateCredentials_Valid_SpecialChars`: Passes for password with special chars (`"P@ss!w0rd#123"`)
- `TestValidateCredentials_Valid_Unicode`: Passes for unicode password (`"пароль密码🔒"`)
- `TestValidateCredentials_Valid_Spaces`: Passes for password with spaces (`"pass with spaces"`)
- `TestValidateCredentials_Valid_Quotes`: Passes for password with quotes (`"pass\"with'quotes"`)
- `TestValidateCredentials_EmptyUsername`: Fails for empty username
- `TestValidateCredentials_EmptyPassword`: Fails for empty password
- `TestValidateCredentials_ControlCharsInUsername`: Fails for control characters in username
  - Test with: `"user\n"` (newline)
  - Test with: `"user\t"` (tab)
  - Test with: `"user\x00"` (null byte)
  - Test with: `"user\x1b"` (escape)

**Test Coverage Requirements**:
- Minimum 90% code coverage (security-critical, per Wave 2 Test Plan)
- All success paths tested
- All failure paths tested
- Special character support validated
- Security checks validated (control characters)

**Test Examples**:
```go
func TestValidateCredentials_Valid_Unicode(t *testing.T) {
    provider := NewBasicAuthProvider("user", "пароль密码🔒")
    err := provider.ValidateCredentials()
    assert.NoError(t, err)
}

func TestValidateCredentials_ControlCharsInUsername(t *testing.T) {
    testCases := []struct {
        name     string
        username string
    }{
        {"newline", "user\n"},
        {"tab", "user\t"},
        {"null", "user\x00"},
        {"escape", "user\x1b"},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            provider := NewBasicAuthProvider(tc.username, "pass")
            err := provider.ValidateCredentials()
            assert.Error(t, err)
            assert.Contains(t, err.Error(), "control characters")
        })
    }
}
```

### Step 4: Size Measurement

**Measure implementation lines**:
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning
tools/line-counter.sh

# Expected output:
# 🎯 Detected base: idpbuilder-oci-push/phase1/wave2/integration
# 📦 Analyzing branch: idpbuilder-oci-push/phase1/wave2/effort-3-auth
# ✅ Total implementation lines: ~350
```

**Size Compliance**:
- Target: 350 lines
- Buffer: ±15% (298-403 lines acceptable)
- Hard limit: 800 lines (MUST NOT EXCEED)
- Tests NOT counted (per R007)

**If approaching 700 lines**:
- STOP IMMEDIATELY
- Report to orchestrator
- Do NOT exceed 800 lines

### Step 5: Run Tests and Coverage

**Run unit tests**:
```bash
cd pkg/auth
go test -v -cover

# Expected: All tests pass, coverage ≥90%
```

**Generate coverage report**:
```bash
go test -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
go tool cover -func=coverage.out | grep total
# Expected: total: (statements) 90.0% or higher
```

### Step 6: Linting and Documentation

**Run linters**:
```bash
go vet ./pkg/auth/...
golangci-lint run ./pkg/auth/...

# Expected: No errors
```

**Verify GoDoc**:
- All public types have GoDoc comments
- All public functions have GoDoc comments
- Security warnings documented (control characters, credential handling)

### Step 7: Commit and Push

**Commit structure**:
```bash
git add pkg/auth/basic.go pkg/auth/basic_test.go
git commit -m "feat(auth): implement basic authentication provider

- Implement NewBasicAuthProvider constructor
- Implement GetAuthenticator returning authn.Basic
- Implement ValidateCredentials with security checks
- Support special characters in passwords (unicode, quotes, spaces)
- Prevent control characters in usernames (security)
- Add helper function for control character detection
- Add 10 test cases with 90%+ coverage (security-critical)

Closes: Effort 1.2.3 - Authentication Implementation
Lines: ~350 (within 350 ±15% estimate)
Coverage: 90%+ (meets Wave 2 Test Plan requirements)"

git push origin idpbuilder-oci-push/phase1/wave2/effort-3-auth
```

---

## Size Management

**Estimated Lines**: 350 lines (implementation code only)
**Measurement Tool**: `${PROJECT_ROOT}/tools/line-counter.sh` (find project root first)
**Check Frequency**: After every major function implementation (~50 lines)
**Split Threshold**:
- Warning: 700 lines (approaching limit)
- Hard stop: 800 lines (MUST NOT EXCEED)

**Size Tracking**:
- After struct definition: ~40 lines
- After NewBasicAuthProvider: ~80 lines
- After GetAuthenticator: ~140 lines
- After ValidateCredentials: ~220 lines
- After helper function: ~250 lines
- After comments and documentation: ~350 lines (target)

---

## Test Requirements

**Minimum Coverage**: 90% (security-critical, per Wave 2 Test Plan)

**Test Categories** (from Wave 2 Test Plan):

| Test Category | Test Cases | Coverage Target |
|---------------|------------|-----------------|
| Constructor | 1 test | 100% of NewBasicAuthProvider |
| GetAuthenticator | 3 tests | 100% of GetAuthenticator |
| ValidateCredentials | 9 tests | 100% of ValidateCredentials |
| Special Characters | 4 tests | Unicode, quotes, spaces |
| Security | 4 tests | Control character detection |

**Total Test Cases**: 10+ tests

**Test Execution**:
```bash
go test ./pkg/auth -v -cover
# MUST achieve ≥90% coverage (security-critical)
```

---

## Dependencies

### Upstream Dependencies (COMPLETED)
- ✅ Wave 1 Effort 3: Auth interface definition (FROZEN)
- ✅ Integration branch: `idpbuilder-oci-push/phase1/wave2/integration` (CREATED)

### Downstream Dependencies
- None (all Wave 2 efforts are parallel)
- Effort 1.2.2 (Registry Client) will USE this implementation via interface
- Wave 3 CLI will use this package

### External Library Dependencies
**Existing Dependencies** (from Wave 1):
- `github.com/google/go-containerregistry` v0.19.0
  - `pkg/authn` for Authenticator types
- `github.com/stretchr/testify` v1.10.0 (testing)

**No New Dependencies Required**

---

## Pattern Compliance

### Go Patterns
- Interface-driven design (implement `auth.Provider` interface)
- Error wrapping with custom error types
- Struct-based provider implementation
- No global state (all data in struct)

### Security Requirements
- **Control character prevention** in usernames (terminal escape attacks)
- **Special character support** in passwords (full unicode support)
- **No credential logging** or exposure in error messages
- **Proper validation** before creating authenticator

### Performance Targets
- Validation should complete in <1 millisecond
- No expensive operations (just string validation)

---

## Acceptance Criteria

**MANDATORY - All must pass before Code Review**:

- [ ] All files created/modified as specified
- [ ] All 2 interface methods implemented correctly (GetAuthenticator, ValidateCredentials)
- [ ] All tests passing (100% pass rate)
- [ ] Code coverage ≥90% (security-critical, per Wave 2 Test Plan)
- [ ] No linting errors (go vet, golangci-lint)
- [ ] Documentation complete (all public methods have GoDoc)
- [ ] Line count within estimate (350 lines ±15% = 298-403 lines)
- [ ] Integration with go-containerregistry working (authn.Basic)
- [ ] Special character support validated (unicode, quotes, spaces)
- [ ] Control character detection working
- [ ] No credential exposure in logs or errors
- [ ] Code committed and pushed to effort branch

**Quality Gates**:
1. **Functionality**: All interface methods work correctly
2. **Testing**: 90%+ coverage with all paths tested
3. **Security**: Control character prevention validated
4. **Size**: Within 350 ±15% lines (298-403)
5. **Documentation**: Complete GoDoc coverage with security warnings

---

## References

**Wave 2 Planning Documents**:
- Wave Implementation Plan: `/home/vscode/workspaces/idpbuilder-oci-push-planning/wave-plans/WAVE-2-IMPLEMENTATION.md`
- Wave Architecture: `/home/vscode/workspaces/idpbuilder-oci-push-planning/wave-plans/WAVE-2-ARCHITECTURE.md`
- Wave Test Plan: `/home/vscode/workspaces/idpbuilder-oci-push-planning/wave-plans/WAVE-2-TEST-PLAN.md`

**Wave 1 Interfaces** (frozen references):
- Auth Interface: `efforts/phase1/wave1/effort-3-auth-tls-interfaces/pkg/auth/types.go`
- Auth Errors: `efforts/phase1/wave1/effort-3-auth-tls-interfaces/pkg/auth/errors.go`

**External Documentation**:
- go-containerregistry authn: https://pkg.go.dev/github.com/google/go-containerregistry/pkg/authn
- HTTP Basic Auth: https://datatracker.ietf.org/doc/html/rfc7617
- Unicode Security: https://www.unicode.org/reports/tr36/

---

## Implementation Checklist

**Pre-Implementation**:
- [ ] Read Wave 1 Auth interface definition
- [ ] Read Wave 2 Architecture (Auth section)
- [ ] Read Wave 2 Test Plan (Auth test cases)
- [ ] Checkout effort branch from integration
- [ ] Verify base branch is correct

**Implementation Phase**:
- [ ] Write test stubs (10+ test cases)
- [ ] Implement `basicAuthProvider` struct
- [ ] Implement `NewBasicAuthProvider()` constructor
- [ ] Implement `GetAuthenticator()` method
- [ ] Implement `ValidateCredentials()` method
- [ ] Implement `containsControlChars()` helper
- [ ] Complete test implementations
- [ ] Run tests (all pass, 90%+ coverage)
- [ ] Run linters (no errors)
- [ ] Add GoDoc comments with security warnings

**Validation Phase**:
- [ ] Measure size (within 298-403 lines)
- [ ] Verify coverage ≥90%
- [ ] Special character validation (unicode, quotes, spaces)
- [ ] Security validation (control characters)
- [ ] No credential exposure verification
- [ ] Commit and push code

---

## Next Steps

**After Implementation Completion**:
1. Code Reviewer will be spawned for effort review
2. Code Reviewer validates all acceptance criteria
3. If approved: Merge to integration branch
4. If issues found: Fix and re-submit for review
5. After all 4 Wave 2 efforts approved: Wave integration testing

**Parallel Work**:
- This effort (1.2.3) runs in parallel with:
  - Effort 1.2.1: Docker Client Implementation
  - Effort 1.2.2: Registry Client Implementation
  - Effort 1.2.4: TLS Configuration Implementation

**No coordination needed** - all efforts are independent until integration phase.

---

## Document Status

**Status**: ✅ READY FOR IMPLEMENTATION
**Created**: 2025-10-29 06:32:00 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**Effort**: 1.2.3 (Authentication Implementation)
**Wave**: Wave 2 of Phase 1
**Branch**: `idpbuilder-oci-push/phase1/wave2/effort-3-auth`
**Base Branch**: `idpbuilder-oci-push/phase1/wave2/integration`

**Compliance**:
- ✅ R213: Complete metadata included
- ✅ R211: Parallelization specified (runs with 1.2.1, 1.2.2, 1.2.4)
- ✅ R341: TDD approach (test plan before implementation)
- ✅ R381: Library versions locked (go-containerregistry v0.19.0 from Wave 1)
- ✅ R383: Plan stored in .software-factory with timestamp
- ✅ Size compliance: 350 lines < 800 hard limit
- ✅ Security-critical: 90% coverage requirement

---

**END OF EFFORT 1.2.3 IMPLEMENTATION PLAN**
