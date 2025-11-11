# Effort 1.1.3 Implementation Plan: Auth & TLS Provider Interface Definitions

**Effort ID**: 1.1.3
**Effort Name**: Auth & TLS Provider Interface Definitions
**Phase**: Phase 1 - Foundation & Interfaces
**Wave**: Wave 1.1 - Interface Definitions
**Created**: 2025-11-11 20:42:13 UTC
**Planner**: Code Reviewer Agent
**State**: EFFORT_PLAN_CREATION

---

## EFFORT INFRASTRUCTURE METADATA (R360)

**EFFORT_NAME**: effort-1.1.3
**PHASE**: phase1
**WAVE**: wave1
**EFFORT_INDEX**: 3
**WORKING_DIRECTORY**: /home/vscode/workspaces/idpbuilder-oci-push-rebuild/efforts/phase1/wave1/effort-1.1.3
**BRANCH_NAME**: idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.3
**BASE_BRANCH**: idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.2
**REMOTE**: origin (https://github.com/jessesanford/idpbuilder.git)

---

## R213 Metadata (EXACT from Wave Plan)

```json
{
  "effort_id": "1.1.3",
  "effort_name": "Auth & TLS Interfaces Definition",
  "branch_name": "idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.3",
  "base_branch": "idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.2",
  "parent_wave": "wave1.1",
  "parent_phase": "phase1",
  "depends_on": ["1.1.2"],
  "estimated_lines": 150,
  "complexity": "low",
  "can_parallelize": false,
  "parallel_with": [],
  "tests_required": [
    "T1.1.3-001: AuthProvider interface compiles",
    "T1.1.3-002: InvalidCredentialsError implements error",
    "T1.1.3-003: MissingCredentialsError implements error",
    "T1.1.3-004: NewAuthProvider constructor signature valid",
    "T1.1.3-005: TLSProvider interface compiles",
    "T1.1.3-006: NewTLSProvider constructor signature valid",
    "T1.1.3-007: Mock AuthProvider satisfies interface",
    "T1.1.3-008: Mock TLSProvider satisfies interface"
  ]
}
```

---

## Overview

**Purpose**: Create the AuthProvider and TLSProvider interfaces that will be implemented in Wave 2 to handle registry authentication and TLS configuration.

**Scope**: Interface definitions ONLY (no implementation logic)

**Key Deliverables**:
- `pkg/auth/interface.go` - Authentication provider interface with 2 methods and 2 error types
- `pkg/tls/interface.go` - TLS configuration provider interface with 3 methods
- Complete test coverage (8 tests, 100% coverage expected)

**Estimated Size**: 150 lines total (75 lines auth + 75 lines tls)

**Estimated Implementation Time**: 2-3 hours (interface definitions are straightforward)

---

## Boundaries (What's In/Out of Scope)

### ✅ In Scope (MUST Implement)

**Authentication (pkg/auth)**:
- AuthProvider interface definition (2 methods)
- NewAuthProvider() constructor function (stub - panics with "not implemented")
- 2 error types: InvalidCredentialsError, MissingCredentialsError
- Complete Go documentation comments
- Import `github.com/google/go-containerregistry/pkg/authn` for authn.Authenticator type

**TLS Configuration (pkg/tls)**:
- TLSProvider interface definition (3 methods)
- NewTLSProvider() constructor function (stub - panics with "not implemented")
- Complete Go documentation comments
- Import `crypto/tls` for tls.Config type

**Testing**:
- 8 unit tests validating interface compilation and error types
- 100% test coverage (interface definitions are simple)
- All tests must pass

### ❌ Out of Scope (MUST NOT Implement)

- NO actual authentication logic (Wave 2)
- NO credential validation implementation (Wave 2)
- NO environment variable reading (Wave 2)
- NO TLS certificate handling (Wave 2)
- NO HTTP client configuration (Wave 2)
- NO integration with go-containerregistry (Wave 2)
- NO mock implementations beyond test compilation checks

---

## File Structure

This effort creates TWO new files in separate packages:

```
pkg/
├── auth/
│   └── interface.go              (75 lines - NEW)
└── tls/
    └── interface.go              (75 lines - NEW)
```

**Total New Lines**: 150 (within 800-line limit, no split required)

---

## Dependencies

### Upstream Dependencies (Must Complete First)

**Effort 1.1.2** (Registry Client Interface Definition):
- Status: MUST be complete before starting this effort
- Reason: Sequential ordering (Wave 1 is fully sequential)
- Verification: Check that branch `idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.2` exists and is merged to integration branch

### Downstream Dependencies (Depend on This Effort)

**Effort 1.1.4** (Command Structure Definition):
- Will import `pkg/auth` and `pkg/tls` interfaces
- Will use AuthProvider and TLSProvider in PushCommand struct
- Cannot start until this effort is complete

### External Dependencies

**Go Libraries** (already in project):
- `github.com/google/go-containerregistry/pkg/authn` - for authn.Authenticator type
- `crypto/tls` - standard library, for tls.Config type

**Verification**:
```bash
# Check dependencies are available
go list -m github.com/google/go-containerregistry
# Expected: github.com/google/go-containerregistry v0.19.0 (or later)
```

---

## Implementation Steps

### Step 1: Create pkg/auth/interface.go (EXACT CODE COPY)

**Action**: Copy EXACTLY from `WAVE-1.1-ARCHITECTURE.md` lines 455-523

**File**: `pkg/auth/interface.go`

**Line Count**: 75 lines

**Key Components**:
1. Package documentation comment
2. AuthProvider interface (2 methods):
   - GetAuthenticator() (authn.Authenticator, error)
   - ValidateCredentials() error
3. NewAuthProvider(username, password string) constructor (stub - panics)
4. InvalidCredentialsError struct with Error() method
5. MissingCredentialsError struct with Error() method

**CRITICAL**: Use EXACT code from architecture document - NO modifications!

**Code to Copy**:

```go
// Package auth provides authentication credential management for OCI registries.
package auth

import (
	"github.com/google/go-containerregistry/pkg/authn"
)

// AuthProvider supplies authentication credentials for registry operations.
// Implementations may read credentials from flags, environment variables, or configuration files.
type AuthProvider interface {
	// GetAuthenticator returns a go-containerregistry compatible authenticator.
	// The authenticator is used by RegistryClient for HTTP Basic Authentication.
	//
	// Example:
	//   authenticator, err := authProvider.GetAuthenticator()
	//   if err != nil {
	//       return fmt.Errorf("getting authenticator: %w", err)
	//   }
	//   // Use authenticator with go-containerregistry's remote.Write()
	GetAuthenticator() (authn.Authenticator, error)

	// ValidateCredentials checks that the credentials meet format and security requirements.
	// Returns an error if username or password is empty, too long, or contains invalid characters.
	//
	// Example:
	//   if err := authProvider.ValidateCredentials(); err != nil {
	//       return fmt.Errorf("invalid credentials: %w", err)
	//   }
	ValidateCredentials() error
}

// NewAuthProvider creates a new authentication provider.
// Credentials are read from the provided username and password parameters first.
// If empty, falls back to environment variables: IDPBUILDER_REGISTRY_USERNAME, IDPBUILDER_REGISTRY_PASSWORD.
//
// Returns an error if no valid credentials are available.
//
// Example:
//   // From flags:
//   provider, err := auth.NewAuthProvider("admin", "secretpass")
//
//   // From environment variables:
//   os.Setenv("IDPBUILDER_REGISTRY_USERNAME", "admin")
//   os.Setenv("IDPBUILDER_REGISTRY_PASSWORD", "secretpass")
//   provider, err := auth.NewAuthProvider("", "")
func NewAuthProvider(username, password string) (AuthProvider, error) {
	// Implementation will be provided in Wave 2
	panic("not implemented")
}

// InvalidCredentialsError indicates that credentials do not meet requirements.
type InvalidCredentialsError struct {
	Reason string
}

func (e *InvalidCredentialsError) Error() string {
	return "invalid credentials: " + e.Reason
}

// MissingCredentialsError indicates that required credentials are not provided.
type MissingCredentialsError struct {
	Field string // "username" or "password"
}

func (e *MissingCredentialsError) Error() string {
	return "missing required credential: " + e.Field
}
```

**Validation After Creation**:
```bash
cd pkg/auth
go build .
# Expected: Build succeeds (pure interface definition)
```

### Step 2: Create pkg/tls/interface.go (EXACT CODE COPY)

**Action**: Copy EXACTLY from `WAVE-1.1-ARCHITECTURE.md` lines 585-655

**File**: `pkg/tls/interface.go`

**Line Count**: 75 lines

**Key Components**:
1. Package documentation comment
2. TLSProvider interface (3 methods):
   - GetTLSConfig() *tls.Config
   - IsInsecure() bool
   - GetWarningMessage() string
3. NewTLSProvider(insecure bool) constructor (stub - panics)

**CRITICAL**: Use EXACT code from architecture document - NO modifications!

**Code to Copy**:

```go
// Package tls provides TLS configuration for secure registry connections.
package tls

import (
	"crypto/tls"
)

// TLSProvider generates TLS configurations for HTTPS connections to OCI registries.
type TLSProvider interface {
	// GetTLSConfig returns a TLS configuration for registry HTTPS connections.
	//
	// In secure mode (default):
	//   - Uses the system certificate pool
	//   - Verifies certificate chains
	//   - Checks that hostnames match certificates
	//
	// In insecure mode:
	//   - Sets InsecureSkipVerify = true
	//   - Accepts self-signed certificates
	//   - Should only be used for local development
	//
	// Example:
	//   tlsConfig := tlsProvider.GetTLSConfig()
	//   transport := &http.Transport{TLSClientConfig: tlsConfig}
	GetTLSConfig() *tls.Config

	// IsInsecure returns true if certificate verification is disabled.
	//
	// Example:
	//   if tlsProvider.IsInsecure() {
	//       fmt.Println("WARNING: Certificate verification disabled")
	//   }
	IsInsecure() bool

	// GetWarningMessage returns a user-facing warning message for insecure mode.
	// Returns an empty string if in secure mode.
	//
	// Example:
	//   if warning := tlsProvider.GetWarningMessage(); warning != "" {
	//       fmt.Println(warning)
	//   }
	GetWarningMessage() string
}

// NewTLSProvider creates a new TLS configuration provider.
//
// If insecure is true:
//   - Certificate verification is disabled (InsecureSkipVerify = true)
//   - Self-signed certificates are accepted
//   - A warning message is generated
//   - Use only for local development with Gitea
//
// If insecure is false:
//   - Standard TLS verification is enabled
//   - System certificate pool is used
//   - Production-ready configuration
//
// Example (secure mode):
//   provider, err := tls.NewTLSProvider(false)
//
// Example (insecure mode for local Gitea):
//   provider, err := tls.NewTLSProvider(true)
//   if warning := provider.GetWarningMessage(); warning != "" {
//       log.Println(warning)
//   }
func NewTLSProvider(insecure bool) (TLSProvider, error) {
	// Implementation will be provided in Wave 2
	panic("not implemented")
}
```

**Validation After Creation**:
```bash
cd pkg/tls
go build .
# Expected: Build succeeds (pure interface definition)
```

### Step 3: Create pkg/auth/interface_test.go (Test Coverage)

**Action**: Create comprehensive tests for auth package

**File**: `pkg/auth/interface_test.go`

**Test Coverage**: 4 tests (T1.1.3-001 through T1.1.3-004)

**Code to Implement**:

```go
package auth_test

import (
	"testing"

	"github.com/jessesanford/idpbuilder/pkg/auth"
)

// T1.1.3-001: AuthProvider interface compiles
func TestAuthProviderInterfaceCompiles(t *testing.T) {
	var _ auth.AuthProvider = nil
}

// T1.1.3-002: InvalidCredentialsError implements error
func TestInvalidCredentialsError_ImplementsError(t *testing.T) {
	err := &auth.InvalidCredentialsError{Reason: "password too short"}
	var _ error = err

	expected := "invalid credentials: password too short"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}
}

// T1.1.3-003: MissingCredentialsError implements error
func TestMissingCredentialsError_ImplementsError(t *testing.T) {
	err := &auth.MissingCredentialsError{Field: "username"}
	var _ error = err

	expected := "missing required credential: username"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}
}

// T1.1.3-004: NewAuthProvider constructor signature valid
func TestNewAuthProvider_SignatureValid(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "not implemented" {
				t.Errorf("Expected panic 'not implemented', got %v", r)
			}
		} else {
			t.Error("Expected NewAuthProvider to panic (not implemented)")
		}
	}()

	_, _ = auth.NewAuthProvider("", "")
}
```

**Run Tests**:
```bash
cd pkg/auth
go test -v -cover
# Expected: 4/4 tests PASS, coverage 100%
```

### Step 4: Create pkg/tls/interface_test.go (Test Coverage)

**Action**: Create comprehensive tests for tls package

**File**: `pkg/tls/interface_test.go`

**Test Coverage**: 2 tests (T1.1.3-005 through T1.1.3-006)

**Code to Implement**:

```go
package tls_test

import (
	"testing"

	"github.com/jessesanford/idpbuilder/pkg/tls"
)

// T1.1.3-005: TLSProvider interface compiles
func TestTLSProviderInterfaceCompiles(t *testing.T) {
	var _ tls.TLSProvider = nil
}

// T1.1.3-006: NewTLSProvider constructor signature valid
func TestNewTLSProvider_SignatureValid(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "not implemented" {
				t.Errorf("Expected panic 'not implemented', got %v", r)
			}
		} else {
			t.Error("Expected NewTLSProvider to panic (not implemented)")
		}
	}()

	_, _ = tls.NewTLSProvider(false)
}
```

**Run Tests**:
```bash
cd pkg/tls
go test -v -cover
# Expected: 2/2 tests PASS, coverage 100%
```

### Step 5: Build Validation (Both Packages)

**Action**: Verify both packages build together without errors

```bash
# From effort root directory
go build ./pkg/auth ./pkg/tls

# Expected: Build succeeds (no output)
```

### Step 6: Full Test Suite Execution

**Action**: Run all tests for this effort

```bash
# From effort root directory
go test ./pkg/auth ./pkg/tls -v -cover

# Expected Output:
# pkg/auth:
#   - T1.1.3-001: PASS
#   - T1.1.3-002: PASS
#   - T1.1.3-003: PASS
#   - T1.1.3-004: PASS
#   Coverage: 100.0%
#
# pkg/tls:
#   - T1.1.3-005: PASS
#   - T1.1.3-006: PASS
#   Coverage: 100.0%
#
# Total: 6/6 tests PASS, 100% coverage
```

**Note**: Tests T1.1.3-007 and T1.1.3-008 (Mock implementations) will be added during integration testing after all Wave 1 efforts are complete.

### Step 7: Linting and Quality Checks

**Action**: Run golangci-lint to ensure code quality

```bash
# From effort root directory
golangci-lint run ./pkg/auth ./pkg/tls

# Expected: No linting errors (interface definitions are clean)
```

### Step 8: Size Measurement

**Action**: Verify implementation is within size limit

```bash
# Find project root first
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ]; then break; fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Run line counter tool
$PROJECT_ROOT/tools/line-counter.sh

# Expected Output:
# Total implementation lines: ~150 lines
# Status: ✅ COMPLIANT (well under 800-line enforcement threshold)
```

### Step 9: Commit and Push

**Action**: Commit all changes to the effort branch

```bash
# Stage all new files
git add pkg/auth/interface.go pkg/auth/interface_test.go
git add pkg/tls/interface.go pkg/tls/interface_test.go
git add .software-factory/

# Commit with descriptive message
git commit -m "feat(interfaces): Add AuthProvider and TLSProvider interface definitions

- Create pkg/auth/interface.go with AuthProvider interface (2 methods)
- Create pkg/tls/interface.go with TLSProvider interface (3 methods)
- Add 2 auth error types: InvalidCredentialsError, MissingCredentialsError
- Add constructor stubs (panic 'not implemented')
- Add comprehensive test coverage (6 tests, 100% coverage)
- All builds and tests passing

Effort: 1.1.3 - Auth & TLS Interfaces
Size: 150 lines (within limit)
Tests: 6/6 PASS
Coverage: 100%

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

# Push to remote
git push origin idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.3
```

---

## Size Management (R304 Compliance)

### Estimated Size Breakdown

**pkg/auth/interface.go**: 75 lines
- Package doc: 2 lines
- Imports: 4 lines
- AuthProvider interface: 20 lines (with docs)
- NewAuthProvider(): 15 lines (with docs)
- InvalidCredentialsError: 10 lines
- MissingCredentialsError: 10 lines
- Error() methods: 14 lines

**pkg/tls/interface.go**: 75 lines
- Package doc: 2 lines
- Imports: 4 lines
- TLSProvider interface: 30 lines (with docs)
- NewTLSProvider(): 25 lines (with docs)
- Implementation stub: 14 lines

**Total Implementation Lines**: 150 lines

### Size Limit Compliance

**Soft Limit**: 700 lines (target)
**Hard Limit (R535)**: 900 lines (Code Reviewer enforcement threshold)
**Current Size**: 150 lines
**Status**: ✅ **COMPLIANT** (well under limits)
**Split Required**: NO

### Measurement Tool

**MUST USE**: `$PROJECT_ROOT/tools/line-counter.sh`
**NEVER USE**: Manual counting, wc -l, cloc, or any other method

**Measurement Frequency**:
- After completing both interface files
- Before final commit
- During code review

---

## Test Requirements (From Wave Test Plan)

### Test Summary

**Total Tests for Effort 1.1.3**: 8 tests (6 unit tests + 2 mock tests)

**Unit Tests** (implemented in this effort): 6 tests
- T1.1.3-001: AuthProvider interface compiles
- T1.1.3-002: InvalidCredentialsError implements error
- T1.1.3-003: MissingCredentialsError implements error
- T1.1.3-004: NewAuthProvider constructor signature valid
- T1.1.3-005: TLSProvider interface compiles
- T1.1.3-006: NewTLSProvider constructor signature valid

**Integration Tests** (wave-level, after all efforts merged): 2 tests
- T1.1.3-007: Mock AuthProvider satisfies interface
- T1.1.3-008: Mock TLSProvider satisfies interface

### Coverage Requirements

**Target Coverage**: 100% (interface definitions are simple)
**Minimum Acceptable**: 100% (no complex logic to exclude)

**Coverage Verification**:
```bash
go test ./pkg/auth ./pkg/tls -cover
# Expected: coverage: 100.0% of statements
```

### Test Quality Standards

**Each test MUST**:
- Have clear, descriptive name matching test ID (e.g., TestAuthProviderInterfaceCompiles)
- Include comments referencing test ID (e.g., // T1.1.3-001)
- Test one specific behavior
- Use table-driven tests where appropriate
- Have meaningful assertions with clear error messages

**Test Execution Requirements**:
- All tests MUST pass before code review
- No skipped tests allowed
- No flaky tests (must be deterministic)
- Tests run in <1 second (interface tests are fast)

---

## Integration Points

### Upstream Integration (Effort 1.1.2)

**Dependency**: Registry Client Interface Definition

**Integration**:
- No direct code dependency (different packages)
- Sequential ordering ensures proper foundation
- Registry package references auth.AuthProvider and tls.TLSProvider as placeholder interfaces (defined in Effort 1.1.2)
- This effort provides the REAL definitions

**Verification**:
```bash
# Check that Effort 1.1.2 is complete
git log --oneline | grep "1.1.2"
# Expected: See commit for Effort 1.1.2 merged
```

### Downstream Integration (Effort 1.1.4)

**Dependency**: Command Structure Definition

**Integration**:
- Effort 1.1.4 imports `pkg/auth` and `pkg/tls`
- PushCommand struct references AuthProvider and TLSProvider interfaces
- NewPushCommand() uses these interfaces for dependency injection
- No implementation calls (Wave 2)

**Package Import Path**:
```go
import (
    "github.com/jessesanford/idpbuilder/pkg/auth"
    "github.com/jessesanford/idpbuilder/pkg/tls"
)
```

### Cross-Package Dependencies

**AuthProvider Used By**:
- pkg/registry.NewRegistryClient() - accepts AuthProvider as parameter
- cmd.PushCommand - holds AuthProvider field

**TLSProvider Used By**:
- pkg/registry.NewRegistryClient() - accepts TLSProvider as parameter
- cmd.PushCommand - holds TLSProvider field

**Interface Contract**:
- These interfaces define the contract for Wave 2 implementations
- Wave 2 will create concrete implementations (BasicAuthProvider, InsecureTLSProvider)
- Wave 2 implementations MUST satisfy these interfaces exactly

---

## Acceptance Criteria (Checklist)

### File Creation

- [ ] File `pkg/auth/interface.go` created with 75 lines
- [ ] File `pkg/tls/interface.go` created with 75 lines
- [ ] File `pkg/auth/interface_test.go` created
- [ ] File `pkg/tls/interface_test.go` created

### Interface Definitions

- [ ] AuthProvider interface defined with 2 methods:
  - [ ] GetAuthenticator() (authn.Authenticator, error)
  - [ ] ValidateCredentials() error
- [ ] TLSProvider interface defined with 3 methods:
  - [ ] GetTLSConfig() *tls.Config
  - [ ] IsInsecure() bool
  - [ ] GetWarningMessage() string

### Error Types

- [ ] InvalidCredentialsError struct defined with Error() method
- [ ] MissingCredentialsError struct defined with Error() method
- [ ] Both error types implement error interface

### Constructor Stubs

- [ ] NewAuthProvider(username, password string) function created
- [ ] NewTLSProvider(insecure bool) function created
- [ ] Both constructors panic("not implemented")

### Documentation

- [ ] All public types have Go documentation comments
- [ ] Package-level documentation present for both packages
- [ ] Method documentation includes examples
- [ ] Constructor documentation explains parameters

### Build & Tests

- [ ] `go build ./pkg/auth` succeeds
- [ ] `go build ./pkg/tls` succeeds
- [ ] `go build ./pkg/auth ./pkg/tls` succeeds
- [ ] All 6 unit tests passing (T1.1.3-001 through T1.1.3-006)
- [ ] Test coverage 100% for both packages
- [ ] `go test ./pkg/auth ./pkg/tls` succeeds

### Code Quality

- [ ] No linting errors (`golangci-lint run ./pkg/auth ./pkg/tls`)
- [ ] Code follows Go style guidelines
- [ ] No hardcoded values (all are interface definitions)
- [ ] No TODO/FIXME comments in production code

### Size Compliance

- [ ] Line count within estimate (128-173 lines acceptable per R502, ±15%)
- [ ] Measured with `$PROJECT_ROOT/tools/line-counter.sh`
- [ ] Size documented in code review report
- [ ] Within 900-line enforcement threshold (R535)

### Version Control

- [ ] All files committed to branch `idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.3`
- [ ] Commit message follows convention
- [ ] Changes pushed to remote origin
- [ ] No uncommitted changes remaining

### Integration Readiness

- [ ] Effort 1.1.2 (Registry Interface) complete and merged
- [ ] Packages can be imported by Effort 1.1.4 (Command Structure)
- [ ] No circular dependencies detected
- [ ] All imports resolve correctly

---

## Pattern Compliance

### Go Interface Design Best Practices

**Interface Segregation**:
- ✅ AuthProvider has 2 focused methods (single responsibility)
- ✅ TLSProvider has 3 focused methods (configuration + introspection)
- ✅ No "god interfaces" with too many methods

**Error Handling**:
- ✅ Custom error types for domain-specific errors
- ✅ All error types implement error interface via Error() method
- ✅ Error messages are descriptive and actionable

**Documentation**:
- ✅ All public types documented
- ✅ Examples provided for all interfaces
- ✅ Package-level documentation explains purpose

**Dependency Injection**:
- ✅ Interfaces allow for easy mocking in tests
- ✅ Implementations can be swapped without changing consumers
- ✅ Constructor functions accept parameters (not global state)

### IDPBuilder Project Patterns

**Package Structure**:
- ✅ Follows `pkg/[domain]/interface.go` pattern
- ✅ Consistent with existing IDPBuilder packages
- ✅ Clear package boundaries (auth, tls, docker, registry, cmd)

**Naming Conventions**:
- ✅ Interface names end with action/role (Provider, Client)
- ✅ Error types end with "Error"
- ✅ Constructor functions named New[Type]()

**Import Paths**:
- ✅ Uses `github.com/jessesanford/idpbuilder/pkg/...`
- ✅ Consistent with existing IDPBuilder imports
- ✅ No relative imports

---

## Security Considerations

### Authentication Security

**What This Effort Does**:
- ✅ Defines interface for credential management
- ✅ Specifies error types for invalid/missing credentials
- ✅ Documents that credentials come from flags OR environment variables

**What This Effort Does NOT Do** (Wave 2):
- ❌ NO actual credential storage (Wave 2)
- ❌ NO plaintext credential logging (Wave 2)
- ❌ NO credential validation logic (Wave 2)

**Security Requirements for Wave 2 Implementation**:
- Credentials MUST NOT be logged
- Credentials MUST be stored in memory only (no disk persistence)
- Environment variables preferred over command-line flags (flags visible in ps)

### TLS Security

**What This Effort Does**:
- ✅ Defines interface for TLS configuration
- ✅ Specifies IsInsecure() method to check security mode
- ✅ Specifies GetWarningMessage() for insecure mode alerting

**What This Effort Does NOT Do** (Wave 2):
- ❌ NO TLS certificate validation (Wave 2)
- ❌ NO certificate pool management (Wave 2)
- ❌ NO InsecureSkipVerify implementation (Wave 2)

**Security Requirements for Wave 2 Implementation**:
- Insecure mode MUST display prominent warning
- Insecure mode MUST be opt-in (secure by default)
- Production usage MUST NOT allow insecure mode

---

## Performance Considerations

### Interface Performance

**Interface Definitions Are Zero-Cost**:
- ✅ No runtime overhead (pure interfaces)
- ✅ No memory allocations (interface types only)
- ✅ No network calls (stubs panic)

**Wave 2 Performance Considerations**:
- GetAuthenticator() should cache authenticator (no repeated creation)
- GetTLSConfig() should cache TLS config (no repeated allocation)
- Credential validation should be fast (<1ms)

### Test Performance

**Expected Test Execution Time**: <100ms total
- Interface compilation tests: instant
- Error type tests: <1ms each
- Constructor signature tests: <1ms each (panic recovery)

---

## Risk Assessment

### Low-Risk Effort

**Risk Level**: **LOW** (interface definitions only)

**Factors**:
- ✅ No complex logic (pure interface definitions)
- ✅ Code copied from validated architecture document (proven design)
- ✅ Small size (150 lines, well under limit)
- ✅ Simple tests (100% pass rate expected)
- ✅ No external service dependencies
- ✅ No database or file I/O
- ✅ No concurrency

### Potential Issues

**1. Dependency Version Conflicts**:
- **Risk**: go-containerregistry version incompatibility
- **Probability**: LOW (already in use by IDPBuilder)
- **Mitigation**: Version pinned in architecture (v0.19.0+)
- **Impact**: Build failure (easy to detect)

**2. Interface Compilation Errors**:
- **Risk**: Syntax errors in interface definitions
- **Probability**: VERY LOW (code copied exactly)
- **Mitigation**: Copy code EXACTLY from architecture document
- **Impact**: Build failure (easy to detect and fix)

**3. Test Failures**:
- **Risk**: Tests don't pass after interface creation
- **Probability**: VERY LOW (simple compilation checks)
- **Mitigation**: Tests are straightforward (interface validation)
- **Impact**: Test failure (easy to fix)

### Rollback Plan

**If This Effort Fails**:
1. Revert commit: `git revert HEAD`
2. Push revert: `git push origin idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.3`
3. Review errors and fix
4. Recommit with fixes

**Impact of Rollback**: LOW (no downstream dependencies yet)

---

## Next Steps (After This Effort)

### Immediate Next Steps

1. **SW Engineer Implementation**: Follow this plan to create both interface files
2. **Self-Review**: SW Engineer validates against acceptance criteria
3. **Code Review**: Code Reviewer validates implementation
4. **Size Measurement**: Code Reviewer measures with line-counter.sh
5. **Merge to Integration**: Orchestrator merges to wave integration branch

### Wave 1 Continuation

**Effort 1.1.4** (Command Structure Definition):
- **Depends On**: This effort (1.1.3)
- **Will Import**: pkg/auth and pkg/tls
- **Estimated Start**: Immediately after this effort merges
- **Estimated Duration**: 3-4 hours

### Wave 1 Completion

After all 4 efforts (1.1.1 through 1.1.4) complete:
1. Wave integration tests run
2. Architect reviews all interfaces
3. Wave merged to main (or Phase 1 integration)
4. Wave 2 planning begins (parallel implementations)

---

## References

### Source Documents

**Wave Implementation Plan**:
- Path: `../../../../planning/phase1/wave1/WAVE-1-IMPLEMENTATION-PLAN.md`
- Section: "Effort 1.1.3: Auth & TLS Interfaces Definition"
- Lines: 799-1147

**Wave Architecture**:
- Path: `../../../../planning/phase1/wave1/WAVE-1.1-ARCHITECTURE.md`
- Auth Section: Lines 455-523
- TLS Section: Lines 585-655

**Wave Test Plan**:
- Path: `../../../../planning/phase1/wave1/WAVE-1-TEST-PLAN.md`
- Tests: T1.1.3-001 through T1.1.3-008

### Rules Referenced

**Critical Rules**:
- R213: Effort Metadata Requirements (metadata at top of plan)
- R303: Save effort plans in .software-factory with timestamps (this file)
- R304: Use only line-counter.sh for size measurement
- R355: Production-ready code only (no stubs beyond constructor panics)
- R383: All metadata in .software-factory with timestamps
- R502: Implementation Plan Quality Gates (EXACT fidelity)
- R535: Code Reviewer enforcement threshold (900 lines)

**Workflow Rules**:
- R221: CD to effort directory in every bash command
- R232: Start immediately (no waiting)
- R287: TODO persistence (save before state transitions)
- R405: End with CONTINUE-SOFTWARE-FACTORY flag

---

## Implementation Checklist for SW Engineer

Use this checklist while implementing:

### Before Starting
- [ ] Read this entire implementation plan
- [ ] Verify Effort 1.1.2 is complete and merged
- [ ] Verify in correct directory: `/home/vscode/workspaces/idpbuilder-oci-push-rebuild/efforts/phase1/wave1/effort-1.1.3`
- [ ] Verify on correct branch: `idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.3`
- [ ] Run `git status` to confirm clean working tree

### During Implementation
- [ ] Create pkg/auth/interface.go (copy EXACTLY from architecture)
- [ ] Create pkg/tls/interface.go (copy EXACTLY from architecture)
- [ ] Build after each file: `go build ./pkg/auth` and `go build ./pkg/tls`
- [ ] Create pkg/auth/interface_test.go (4 tests)
- [ ] Create pkg/tls/interface_test.go (2 tests)
- [ ] Run tests after each file: `go test ./pkg/auth -v` and `go test ./pkg/tls -v`
- [ ] Verify 100% test coverage
- [ ] Run linter: `golangci-lint run ./pkg/auth ./pkg/tls`
- [ ] Measure size with line-counter.sh

### Before Committing
- [ ] All 6 tests passing
- [ ] Build succeeds for both packages
- [ ] Linter shows no errors
- [ ] Size within estimate (128-173 lines)
- [ ] No uncommitted files
- [ ] No TODO/FIXME in production code

### After Committing
- [ ] Push to remote
- [ ] Notify orchestrator (effort ready for review)
- [ ] Create work log entry

---

## CONTINUE-SOFTWARE-FACTORY Flag

**Status**: Plan creation successful
**Next State**: WAITING_FOR_IMPLEMENTATION
**Orchestrator Action**: Spawn SW Engineer for implementation

CONTINUE-SOFTWARE-FACTORY=TRUE
