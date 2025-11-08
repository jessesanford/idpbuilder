# Authentication Implementation - Effort 1.2.3 Implementation Plan

**Effort**: 1.2.3 - Authentication Implementation
**Phase**: Phase 1 - Foundation & Interfaces
**Wave**: Wave 2 - Core Package Implementations
**Created**: 2025-10-29 21:33:26 UTC
**Planner**: Code Reviewer Agent (code-reviewer)

---

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)

### R213 Metadata
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

### Parallelization Info
**Can Parallelize**: Yes
**Parallel With**: Efforts 1.2.1 (Docker), 1.2.2 (Registry), 1.2.4 (TLS)
**Size Estimate**: 350 lines (implementation only, tests excluded per R007)
**Dependencies**: Wave 1 Effort 3 (Auth interface definition - COMPLETED)

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

**Implementation File**:
- **File**: `pkg/auth/basic.go`
- **Purpose**: Basic authentication provider implementation
- **Estimated Lines**: ~350 lines (including comprehensive GoDoc)
- **Content**:
  - `basicAuthProvider` struct implementation
  - `NewBasicAuthProvider()` constructor
  - `GetAuthenticator()` method returning authn.Authenticator
  - `ValidateCredentials()` validation method
  - Helper function (`containsControlChars`)

**Test File** (NOT counted in line estimates per R007):
- **File**: `pkg/auth/basic_test.go`
- **Purpose**: Comprehensive unit tests for basic authentication
- **Estimated Lines**: ~250 lines (tests excluded from line count)
- **Content**:
  - 10+ test cases covering all methods
  - Special character password tests
  - Control character username tests
  - Success and failure paths

### Modified Files

**None** - go.mod already has go-containerregistry from Wave 1

---

## Implementation Steps

### Step 1: Create Package Structure
```bash
# Navigate to effort directory
cd $EFFORT_DIR

# Create pkg/auth directory if it doesn't exist
mkdir -p pkg/auth

# Verify Wave 1 interface exists
ls pkg/auth/interface.go  # Should exist from Wave 1
ls pkg/auth/errors.go     # Should exist from Wave 1
```

### Step 2: Implement basicAuthProvider Struct

**File**: `pkg/auth/basic.go`

The implementation MUST:
1. Use go-containerregistry's `authn` package for authenticator type
2. Implement ALL 2 methods from Wave 1 `auth.Provider` interface
3. Use Wave 1 error type: `CredentialValidationError`
4. Store credentials in `Credentials` struct (from Wave 1)
5. Support ALL special characters in passwords (no restrictions except non-empty)

**Struct Definition**:
```go
// basicAuthProvider implements the Provider interface using basic username/password authentication.
type basicAuthProvider struct {
    credentials Credentials
}
```

### Step 3: Implement NewBasicAuthProvider() Constructor

**Requirements** (from Wave 2 Architecture):
- Accept username and password as strings
- Store in `Credentials` struct: `Credentials{Username: username, Password: password}`
- Return `basicAuthProvider` implementing `Provider` interface
- No validation in constructor (validation happens in ValidateCredentials)

**Function Signature**:
```go
func NewBasicAuthProvider(username, password string) Provider
```

**Key Implementation Points**:
- Simple constructor - just store credentials
- Return pointer to basicAuthProvider
- No error return (validation deferred to ValidateCredentials)

### Step 4: Implement GetAuthenticator() Method

**Requirements** (from Wave 2 Architecture):
- Call `ValidateCredentials()` first, return error if invalid
- Create `authn.Basic{Username: username, Password: password}`
- Return authenticator compatible with go-containerregistry
- Return `CredentialValidationError` if validation fails

**Function Signature**:
```go
func (p *basicAuthProvider) GetAuthenticator() (authn.Authenticator, error)
```

**Implementation Flow**:
1. Validate credentials first
2. If validation fails, propagate error
3. Create authn.Basic struct with credentials
4. Return authenticator

**Usage Example**:
```go
authenticator, err := provider.GetAuthenticator()
if err != nil {
    return fmt.Errorf("failed to get authenticator: %w", err)
}
// Use authenticator with remote.Push(ref, image, remote.WithAuth(authenticator))
```

### Step 5: Implement ValidateCredentials() Method

**Requirements** (from Wave 2 Architecture):
- Check username is not empty (return CredentialValidationError if empty)
- Check username contains no control characters (< 32 or == 127)
- Check password is not empty (return CredentialValidationError if empty)
- Allow ALL printable characters in password (including quotes, spaces, unicode)
- Return nil if valid

**Function Signature**:
```go
func (p *basicAuthProvider) ValidateCredentials() error
```

**Validation Rules**:

**Username Validation**:
- NOT empty
- NO control characters (prevents terminal escape sequence attacks)
- Control characters are: ASCII < 32 or ASCII == 127

**Password Validation**:
- NOT empty
- ALLOW everything (HTTP Basic Auth transmits as-is, base64 encoded)
- No password strength requirements (user's responsibility)
- Special characters explicitly supported:
  - Unicode: `"пароль密码🔒"`
  - Quotes: `"pass\"with'quotes"`
  - Spaces: `"pass with spaces"`
  - Special chars: `"P@ss!w0rd#123"`

**Security Considerations**:
- Username: No control characters (prevents terminal escape sequence attacks)
- Password: Allow everything (HTTP Basic Auth transmits as-is, base64 encoded)
- No password strength requirements (user's responsibility)
- No credential logging or exposure

### Step 6: Implement Helper Function

**Function**: `containsControlChars(s string) bool`

**Purpose**: Check if string contains control characters

**Implementation**:
```go
func containsControlChars(s string) bool {
    for _, r := range s {
        if r < 32 || r == 127 {
            return true
        }
    }
    return false
}
```

**Control Characters**:
- ASCII 0-31: Null, tab, newline, escape, etc.
- ASCII 127: Delete character

### Step 7: Add Comprehensive GoDoc Comments

**Requirements**:
- All public functions MUST have GoDoc comments
- Comments explain parameters, return values, behavior
- Include usage examples
- Document security considerations

**Example GoDoc Structure**:
```go
// NewBasicAuthProvider creates a basic authentication provider.
//
// Basic authentication uses username and password credentials transmitted
// via HTTP Basic Auth header to the registry.
//
// Parameters:
//   - username: Registry username (typically "giteaadmin" for Gitea)
//   - password: Registry password (supports all special characters)
//
// Returns:
//   - Provider: Authentication provider interface implementation
//
// Example:
//   provider := auth.NewBasicAuthProvider("giteaadmin", "myP@ssw0rd!")
//   if err := provider.ValidateCredentials(); err != nil {
//       return fmt.Errorf("invalid credentials: %w", err)
//   }
func NewBasicAuthProvider(username, password string) Provider {
    // implementation
}
```

---

## Exact Code Specifications (from Wave 2 Architecture)

### Complete Implementation Reference

**File**: `pkg/auth/basic.go`

```go
// Package auth provides registry authentication implementations.
package auth

import (
    "strings"

    "github.com/google/go-containerregistry/pkg/authn"
)

// basicAuthProvider implements the Provider interface using basic username/password authentication.
type basicAuthProvider struct {
    credentials Credentials
}

// NewBasicAuthProvider creates a basic authentication provider.
//
// Basic authentication uses username and password credentials transmitted
// via HTTP Basic Auth header to the registry.
//
// Parameters:
//   - username: Registry username (typically "giteaadmin" for Gitea)
//   - password: Registry password (supports all special characters)
//
// Returns:
//   - Provider: Authentication provider interface implementation
//
// Example:
//   provider := auth.NewBasicAuthProvider("giteaadmin", "myP@ssw0rd!")
//   if err := provider.ValidateCredentials(); err != nil {
//       return fmt.Errorf("invalid credentials: %w", err)
//   }
func NewBasicAuthProvider(username, password string) Provider {
    return &basicAuthProvider{
        credentials: Credentials{
            Username: username,
            Password: password,
        },
    }
}

// GetAuthenticator returns an authn.Authenticator compatible with go-containerregistry.
//
// This method converts internal credentials to the authn.Basic format expected
// by go-containerregistry's remote.Push() function.
//
// Returns:
//   - authn.Authenticator: Authenticator instance for go-containerregistry
//   - error: ValidationError if credentials are malformed
//
// Example:
//   authenticator, err := provider.GetAuthenticator()
//   if err != nil {
//       return fmt.Errorf("failed to get authenticator: %w", err)
//   }
//   // Use authenticator with remote.Push(ref, image, remote.WithAuth(authenticator))
func (p *basicAuthProvider) GetAuthenticator() (authn.Authenticator, error) {
    // Validate credentials before creating authenticator
    if err := p.ValidateCredentials(); err != nil {
        return nil, err
    }

    // Create go-containerregistry Basic authenticator
    authenticator := &authn.Basic{
        Username: p.credentials.Username,
        Password: p.credentials.Password,
    }

    return authenticator, nil
}

// ValidateCredentials performs pre-flight validation of credentials.
//
// This method checks:
//   - Username is not empty
//   - Password is not empty
//   - Username contains no control characters
//
// Note: This does NOT validate credentials with the registry. It only
// checks that credentials meet basic format requirements.
//
// Returns:
//   - error: CredentialValidationError with details if invalid, nil if valid
//
// Example:
//   if err := provider.ValidateCredentials(); err != nil {
//       return fmt.Errorf("invalid credentials: %w", err)
//   }
func (p *basicAuthProvider) ValidateCredentials() error {
    // Check username
    if p.credentials.Username == "" {
        return &CredentialValidationError{
            Field:  "username",
            Reason: "username cannot be empty",
        }
    }

    // Check for control characters in username
    if containsControlChars(p.credentials.Username) {
        return &CredentialValidationError{
            Field:  "username",
            Reason: "username contains control characters",
        }
    }

    // Check password
    if p.credentials.Password == "" {
        return &CredentialValidationError{
            Field:  "password",
            Reason: "password cannot be empty",
        }
    }

    // Password can contain ANY characters (including quotes, spaces, unicode)
    // No validation on password content

    return nil
}

// Helper functions

func containsControlChars(s string) bool {
    for _, r := range s {
        if r < 32 || r == 127 {
            return true
        }
    }
    return false
}
```

**Estimated Lines**: ~350 lines (with comprehensive GoDoc and helper functions)

---

## Test Requirements

### Test File: pkg/auth/basic_test.go

**Minimum Test Coverage**: 90% (security-critical, per Wave 2 Test Plan)

### Test Categories (from Wave 2 Test Plan)

#### A. Constructor Tests

**TC-AUTH-IMPL-001: NewBasicAuthProvider Success**
- Purpose: Verify constructor creates provider correctly
- Test: Create provider with valid username and password
- Expected: Provider created successfully, credentials stored

#### B. GetAuthenticator Tests

**TC-AUTH-IMPL-002: GetAuthenticator Success with Valid Credentials**
- Purpose: Verify authenticator creation with valid credentials
- Test: Create provider, call GetAuthenticator
- Expected: Returns authn.Authenticator, no error
- Validation: Verify authenticator has correct username/password

**TC-AUTH-IMPL-003: GetAuthenticator Fails with Empty Username**
- Purpose: Verify validation errors propagate from GetAuthenticator
- Test: Create provider with empty username, call GetAuthenticator
- Expected: Returns CredentialValidationError for username

#### C. ValidateCredentials Tests

**TC-AUTH-IMPL-004: ValidateCredentials Passes for Valid Credentials**
- Purpose: Verify validation accepts valid credentials
- Test cases:
  - Simple: `"user"` / `"pass"`
  - Special chars: `"user"` / `"P@ss!w0rd#123"`
  - Unicode: `"user"` / `"пароль密码🔒"`
  - Spaces: `"user"` / `"pass with spaces"`
  - Quotes: `"user"` / `"pass\"with'quotes"`
- Expected: All pass validation (return nil)

**TC-AUTH-IMPL-005: ValidateCredentials Fails for Empty Username**
- Purpose: Verify empty username rejected
- Test: Provider with empty username
- Expected: Returns CredentialValidationError with field="username"

**TC-AUTH-IMPL-006: ValidateCredentials Fails for Empty Password**
- Purpose: Verify empty password rejected
- Test: Provider with empty password
- Expected: Returns CredentialValidationError with field="password"

**TC-AUTH-IMPL-007: ValidateCredentials Fails for Control Characters in Username**
- Purpose: Verify control character detection
- Test cases:
  - Newline: `"user\n"`
  - Tab: `"user\t"`
  - Null byte: `"user\x00"`
  - Escape: `"user\x1b"`
- Expected: All return CredentialValidationError with field="username"

### Test Implementation Examples

**Constructor Test**:
```go
func TestNewBasicAuthProvider(t *testing.T) {
    // Given: Valid credentials
    username := "giteaadmin"
    password := "mypassword"

    // When: Creating provider
    provider := NewBasicAuthProvider(username, password)

    // Then: Provider created successfully
    assert.NotNil(t, provider)

    // Verify implements Provider interface
    var _ Provider = provider
}
```

**GetAuthenticator Test**:
```go
func TestGetAuthenticator_Success(t *testing.T) {
    // Given: Provider with valid credentials
    provider := NewBasicAuthProvider("giteaadmin", "password")

    // When: Getting authenticator
    authenticator, err := provider.GetAuthenticator()

    // Then: Returns authenticator successfully
    require.NoError(t, err)
    require.NotNil(t, authenticator)

    // Verify authenticator is correct type
    basicAuth, ok := authenticator.(*authn.Basic)
    require.True(t, ok)
    assert.Equal(t, "giteaadmin", basicAuth.Username)
    assert.Equal(t, "password", basicAuth.Password)
}
```

**Special Character Password Test**:
```go
func TestValidateCredentials_SpecialCharacters(t *testing.T) {
    testCases := []struct {
        name     string
        password string
    }{
        {"simple", "pass"},
        {"special_chars", "P@ss!w0rd#123"},
        {"unicode", "пароль密码🔒"},
        {"spaces", "pass with spaces"},
        {"quotes", "pass\"with'quotes"},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Given: Provider with special character password
            provider := NewBasicAuthProvider("user", tc.password)

            // When: Validating credentials
            err := provider.ValidateCredentials()

            // Then: Validation passes
            assert.NoError(t, err, "Password should be valid: %s", tc.password)
        })
    }
}
```

**Control Character Username Test**:
```go
func TestValidateCredentials_ControlCharactersInUsername(t *testing.T) {
    testCases := []struct {
        name     string
        username string
    }{
        {"newline", "user\n"},
        {"tab", "user\t"},
        {"null_byte", "user\x00"},
        {"escape", "user\x1b"},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Given: Provider with control character in username
            provider := NewBasicAuthProvider(tc.username, "password")

            // When: Validating credentials
            err := provider.ValidateCredentials()

            // Then: Validation fails
            require.Error(t, err)

            var valErr *CredentialValidationError
            assert.ErrorAs(t, err, &valErr)
            assert.Equal(t, "username", valErr.Field)
        })
    }
}
```

### Test Coverage Requirements
- Minimum 90% code coverage (security-critical package)
- All success paths tested
- All failure paths tested
- Special character support validated
- Security checks validated (control characters)

---

## Size Management

### Size Monitoring

**Estimated Lines**: 350 lines (implementation only, tests excluded per R007)

**Measurement Tool**: Use ONLY `${PROJECT_ROOT}/tools/line-counter.sh`

**Check Frequency**: After implementation complete

**Thresholds**:
- Target: 350 lines
- Warning: 700 lines (approaching limit)
- Stop: 800 lines (hard limit, requires split)

### Size Measurement Protocol

```bash
# Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ]; then
        break
    fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Run line counter (auto-detects base branch)
cd $EFFORT_DIR
$PROJECT_ROOT/tools/line-counter.sh

# Tool will show:
# 🎯 Detected base: idpbuilder-oci-push/phase1/wave2/integration
# 📦 Analyzing branch: idpbuilder-oci-push/phase1/wave2/effort-3-auth
# ✅ Total implementation lines: [count]
```

**NEVER DO**:
- ❌ Manual line counting (wc -l, find | xargs)
- ❌ Including test files in count
- ❌ Using wrong base branch

### Expected Size Compliance

**Status**: ✅ EXPECTED TO BE COMPLIANT
- Estimated: 350 lines
- Limit: 800 lines
- Buffer: 450 lines (128% safety margin)
- Split Required: NO

---

## Dependencies

### Upstream Dependencies (must complete before this effort)
- ✅ Wave 1 Effort 3: Auth interface definition (COMPLETED)
- ✅ Integration branch: `idpbuilder-oci-push/phase1/wave2/integration` (CREATED)

### Downstream Dependencies (efforts that depend on this)
- None (all Wave 2 efforts are parallel)
- Effort 1.2.2 (Registry Client) will USE this implementation via interface
- Wave 3 CLI will use this package

### External Library Dependencies
- `github.com/google/go-containerregistry` v0.19.0 (already in Wave 1)
  - `pkg/authn` for Authenticator types
- `github.com/stretchr/testify` v1.10.0 (already in Wave 1, for tests)

**No go.mod changes required** - all dependencies already present from Wave 1

---

## Integration Points

### Interface Implementation

This effort implements the `auth.Provider` interface from Wave 1:

```go
// From Wave 1: pkg/auth/interface.go
type Provider interface {
    GetAuthenticator() (authn.Authenticator, error)
    ValidateCredentials() error
}
```

### Usage by Other Packages

**Registry Client (Effort 1.2.2)** will use this implementation:

```go
// In registry client
authProvider := auth.NewBasicAuthProvider("giteaadmin", "password")
registryClient, err := registry.NewClient(authProvider, tlsProvider)
if err != nil {
    return err
}

// Push operation will call GetAuthenticator internally
err = registryClient.Push(ctx, image, targetRef, nil)
```

### Error Types

Uses Wave 1 error type:

```go
// From Wave 1: pkg/auth/errors.go
type CredentialValidationError struct {
    Field  string
    Reason string
}
```

---

## Acceptance Criteria

### Implementation Completeness
- [ ] All files created as specified
- [ ] All 2 interface methods implemented correctly
- [ ] Helper function implemented (containsControlChars)
- [ ] Wave 1 `Provider` interface fully satisfied

### Code Quality
- [ ] All tests passing (100% pass rate)
- [ ] Code coverage ≥90% (security-critical, per Wave 2 Test Plan)
- [ ] No linting errors (go vet, golangci-lint)
- [ ] Documentation complete (all public methods have GoDoc)

### Size Compliance
- [ ] Line count within estimate (350 lines ±15% = 298-403 lines)
- [ ] Measured using designated line counter tool ONLY
- [ ] No manual line counting performed

### Functionality
- [ ] Integration with go-containerregistry working (authn.Basic)
- [ ] Special character support validated (unicode, quotes, spaces)
- [ ] Control character detection working
- [ ] No credential exposure in logs or errors

### Security Validation
- [ ] Control character detection prevents injection attacks
- [ ] All special characters supported in passwords
- [ ] No credential logging or leakage
- [ ] Proper error messages without exposing credentials

### Testing
- [ ] All 10+ test cases passing
- [ ] Special character tests passing
- [ ] Control character tests passing
- [ ] Error path tests passing

---

## Implementation Workflow

### Development Process

1. **Pre-Implementation**:
   - [ ] Read this implementation plan completely
   - [ ] Read Wave 2 Architecture (concrete implementation reference)
   - [ ] Read Wave 2 Test Plan (test specifications)
   - [ ] Read Wave 1 interface definition (contract to implement)

2. **Setup**:
   - [ ] Verify on correct branch: `idpbuilder-oci-push/phase1/wave2/effort-3-auth`
   - [ ] Verify based on: `idpbuilder-oci-push/phase1/wave2/integration`
   - [ ] Create pkg/auth directory if needed

3. **TDD Implementation** (R341 Compliance):
   - [ ] Write test file FIRST (`pkg/auth/basic_test.go`)
   - [ ] Implement code to pass tests (`pkg/auth/basic.go`)
   - [ ] Run tests frequently: `go test ./pkg/auth -v`
   - [ ] Verify coverage: `go test ./pkg/auth -cover`

4. **Validation**:
   - [ ] All tests passing
   - [ ] Coverage ≥90%
   - [ ] Linting clean: `go vet ./pkg/auth`
   - [ ] Measure size with line counter

5. **Commit and Push**:
   - [ ] Commit regularly during development
   - [ ] Push when complete
   - [ ] Notify orchestrator for code review

### Testing Commands

```bash
# Run tests
go test ./pkg/auth -v

# Run tests with coverage
go test ./pkg/auth -cover

# Generate detailed coverage report
go test ./pkg/auth -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run linting
go vet ./pkg/auth
golangci-lint run ./pkg/auth

# Build check
go build ./pkg/auth
```

---

## Risk Mitigation

### High-Risk Areas

**1. Special Characters in Passwords**:
- **Risk**: Password encoding issues with quotes/unicode
- **Mitigation**:
  - Extensive test cases for special characters
  - Unicode password tests
  - go-containerregistry handles base64 encoding
  - No preprocessing of passwords (pass through as-is)

**2. Control Character Detection**:
- **Risk**: Missing control characters allows injection attacks
- **Mitigation**:
  - Comprehensive character range check (0-31, 127)
  - Test cases for all common control characters
  - Clear validation error messages

**3. Credential Exposure**:
- **Risk**: Credentials logged or exposed in errors
- **Mitigation**:
  - No credential logging anywhere
  - Error messages describe problem without exposing credentials
  - Review all error messages for credential leakage

### Testing Strategy

**Unit Tests with Real go-containerregistry**:
- Use actual `authn.Basic` type (not mocks)
- Validate authenticator structure
- Test actual interface compatibility

**Security Testing**:
- 90% coverage requirement (higher than normal)
- All control character variations tested
- All special character variations tested
- Boundary testing (empty strings, max lengths)

---

## Compliance Verification

### R307: Independent Branch Mergeability
- ✅ Uses frozen Wave 1 interface
- ✅ No dependencies on other Wave 2 efforts
- ✅ Can merge independently
- ✅ Build guaranteed green (implements interface exactly)

### R501: Progressive Trunk-Based Development
- ✅ Branches from Wave 2 integration branch
- ✅ Builds on Wave 1 foundation incrementally
- ✅ No branching from main (uses cascade)

### R359: No Code Deletion
- ✅ Pure addition (no Wave 1 code modified)
- ✅ No interface changes
- ✅ Only new implementation file created

### R383: Metadata File Organization
- ✅ This plan in `.software-factory/phase1/wave2/effort-3-auth/`
- ✅ Timestamped: `IMPLEMENTATION-PLAN--20251029-213326.md`
- ✅ Working tree clean (only source code visible)

### R341: Test-Driven Development
- ✅ Test Plan created BEFORE implementation
- ✅ Tests reference actual Wave 1 interfaces
- ✅ Progressive Realism approach used
- ✅ Coverage target defined (90%)

---

## Next Steps

### Immediate Actions (for SW Engineer)

1. **Checkout Branch**:
   ```bash
   # Verify you're on the correct branch
   git branch --show-current
   # Should be: idpbuilder-oci-push/phase1/wave2/effort-3-auth
   ```

2. **Read Planning Documents**:
   - [ ] This implementation plan (you are here)
   - [ ] Wave 2 Architecture (`wave-plans/WAVE-2-ARCHITECTURE.md`)
   - [ ] Wave 2 Test Plan (`wave-plans/WAVE-2-TEST-PLAN.md`)
   - [ ] Wave 1 Auth Interface (`pkg/auth/interface.go`)

3. **Start TDD Implementation**:
   - [ ] Create `pkg/auth/basic_test.go` FIRST
   - [ ] Write all test cases from test plan
   - [ ] Implement `pkg/auth/basic.go` to pass tests
   - [ ] Iterate until 90% coverage achieved

4. **Validate and Commit**:
   - [ ] Run all tests
   - [ ] Check coverage
   - [ ] Measure size with line counter
   - [ ] Commit and push

### For Code Reviewer (after SW Engineer completes)

1. **Review Checklist**:
   - [ ] All acceptance criteria met
   - [ ] Tests passing (100% pass rate)
   - [ ] Coverage ≥90%
   - [ ] Line count within limit (298-403 lines)
   - [ ] Interface correctly implemented
   - [ ] Special character support validated
   - [ ] Control character detection working
   - [ ] No credential exposure
   - [ ] Documentation complete

2. **Decision**:
   - ACCEPTED → Merge to integration branch
   - NEEDS_FIXES → Send back with specific feedback
   - NEEDS_SPLIT → Should not happen (well under 800 lines)

---

## Document Status

**Status**: ✅ READY FOR IMPLEMENTATION
**Created**: 2025-10-29 21:33:26 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**Effort**: 1.2.3 - Authentication Implementation
**Wave**: Wave 2 of Phase 1
**Phase**: Phase 1 - Foundation & Interfaces

**Compliance Summary**:
- ✅ R213: Complete metadata included
- ✅ R211: Parallelization info specified (can parallelize with 1.2.1, 1.2.2, 1.2.4)
- ✅ R341: TDD approach (test plan before implementation)
- ✅ R307: Independent branch mergeability ensured
- ✅ R359: No code deletion (pure addition)
- ✅ R383: Metadata properly organized with timestamp
- ✅ Size compliance: 350 lines < 800 line limit

**Next State Transition**: IMPLEMENTATION (SW Engineer)
- SW Engineer will read this plan
- Implement following TDD approach
- Submit for code review when complete

---

**END OF IMPLEMENTATION PLAN**
