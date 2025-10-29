# Implementation Plan: Auth & TLS Interface Definitions
## Effort 1.1.3 - Phase 1, Wave 1

**Created**: 2025-10-29T04:08:10Z
**Planner**: @agent-code-reviewer
**Effort ID**: 1.1.3
**Phase**: 1 (Foundation & Interfaces)
**Wave**: 1 (Interface & Contract Definitions)

---

## 🚨 EFFORT INFRASTRUCTURE METADATA (R360)

**EFFORT_NAME**: effort-3-auth-tls-interfaces
**BRANCH**: idpbuilder-oci-push/phase1/wave1/effort-3-auth-tls-interfaces
**BASE_BRANCH**: idpbuilder-oci-push/phase1/wave1/integration
**PHASE**: 1
**WAVE**: 1
**EFFORT_INDEX**: 3
**PARALLELIZATION**: sequential
**CAN_PARALLELIZE**: No
**PARALLEL_WITH**: None
**DEPENDENCIES**: []

---

## Overview

**Purpose**: Define authentication provider and TLS configuration provider interfaces for registry client usage.

**Scope**: Interface definitions ONLY - no implementations (Wave 1 contract definition)

**Estimated Size**: 140 lines (implementation code only, excluding tests)

**Expected Outcomes**:
- `auth.Provider` interface for authentication
- `Credentials` struct for basic auth
- `auth` error types
- `tls.ConfigProvider` interface for TLS configuration
- `Config` struct for TLS settings
- Package documentation for both packages
- Mock implementations for both interfaces
- 100% test coverage

---

## 🔴🔴🔴 REPOSITORY CONTEXT (R251/R309) 🔴🔴🔴

**CRITICAL UNDERSTANDING**:
- ✅ This plan is for the idpbuilder TARGET repo (https://github.com/jessesanford/idpbuilder.git)
- ✅ Implementation will happen in TARGET repo clone
- ✅ NOT in Software Factory planning repo
- ✅ Files reference TARGET repo structure: `pkg/`, `cmd/`, etc.

**Working Directory**: `/efforts/phase1/wave1/effort-3-auth-tls-interfaces/`
**Target Branch**: `idpbuilder-oci-push/phase1/wave1/effort-3-auth-tls-interfaces`
**Base Branch**: `idpbuilder-oci-push/phase1/wave1/integration`

---

## 🔴🔴🔴 EXPLICIT SCOPE CONTROL (R311 - SUPREME LAW) 🔴🔴🔴

### IMPLEMENT EXACTLY:

**Auth Package (5 implementation files, ~95 lines total)**:
1. `pkg/auth/interface.go` (~60 lines)
   - `Provider` interface with 2 methods
   - `Credentials` struct
   - `NewBasicAuthProvider()` function signature
2. `pkg/auth/errors.go` (~20 lines)
   - `CredentialValidationError` type
3. `pkg/auth/doc.go` (~15 lines)

**TLS Package (4 implementation files, ~75 lines total)**:
1. `pkg/tls/interface.go` (~55 lines)
   - `ConfigProvider` interface with 2 methods
   - `Config` struct
   - `NewConfigProvider()` function signature
2. `pkg/tls/doc.go` (~10 lines)

**Test Files (5 files, ~190 lines total - NOT counted)**:
- `pkg/auth/interface_test.go` (~30 lines)
- `pkg/auth/mock_test.go` (~50 lines)
- `pkg/tls/interface_test.go` (~30 lines)
- `pkg/tls/mock_test.go` (~60 lines)

**TOTAL IMPLEMENTATION**: ~140 lines (excludes tests per R007)

### DO NOT IMPLEMENT:

❌ NO actual authentication implementation (Wave 2)
❌ NO go-containerregistry integration (Wave 2)
❌ NO TLS certificate loading (Wave 2)
❌ NO HTTP transport configuration (Wave 2)
❌ NO credential storage or retrieval (Wave 2)
❌ NO additional helper functions
❌ NO logging or metrics
❌ NO additional interfaces
❌ NO additional error types

---

## File Structure

### Files to Create

**Auth Package (5 files)**:
1. `pkg/auth/interface.go` (~60 lines)
2. `pkg/auth/errors.go` (~20 lines)
3. `pkg/auth/doc.go` (~15 lines)
4. `pkg/auth/interface_test.go` (~30 lines)
5. `pkg/auth/mock_test.go` (~50 lines)

**TLS Package (4 files)**:
1. `pkg/tls/interface.go` (~55 lines)
2. `pkg/tls/doc.go` (~10 lines)
3. `pkg/tls/interface_test.go` (~30 lines)
4. `pkg/tls/mock_test.go` (~60 lines)

---

## Implementation Steps

### Step 1: Create Package Directories

```bash
cd /efforts/phase1/wave1/effort-3-auth-tls-interfaces
mkdir -p pkg/auth pkg/tls
```

### Step 2: Create pkg/auth/interface.go

**Complete file content** (from Wave Implementation Plan lines 720-767):

```go
// Package auth provides interfaces and types for registry authentication.
package auth

import (
	"github.com/google/go-containerregistry/pkg/authn"
)

// Provider defines operations for providing authentication credentials to registries.
type Provider interface {
	// GetAuthenticator returns an authn.Authenticator compatible with go-containerregistry.
	//
	// Returns:
	//   - authn.Authenticator: Authenticator instance
	//   - error: ValidationError if credentials are malformed
	GetAuthenticator() (authn.Authenticator, error)

	// ValidateCredentials performs pre-flight validation of credentials.
	//
	// Returns:
	//   - error: ValidationError with details if invalid, nil if valid
	ValidateCredentials() error
}

// Credentials holds authentication information for basic auth.
type Credentials struct {
	// Username for registry authentication.
	Username string

	// Password for registry authentication.
	// Supports ALL special characters including quotes, spaces, unicode.
	Password string
}

// NewBasicAuthProvider creates a basic authentication provider.
//
// Parameters:
//   - username: Registry username
//   - password: Registry password (supports all special characters)
//
// Returns:
//   - Provider: Authentication provider interface implementation
func NewBasicAuthProvider(username, password string) Provider {
	// Implementation will be provided in Wave 2 (pkg/auth/basic.go)
	panic("not implemented - interface definition only")
}
```

### Step 3: Create pkg/auth/errors.go

**Complete file content** (from Wave Implementation Plan lines 769-785):

```go
package auth

import "fmt"

// CredentialValidationError indicates credential validation failed.
type CredentialValidationError struct {
	Field  string // "username" or "password"
	Reason string
}

func (e *CredentialValidationError) Error() string {
	return fmt.Sprintf("credential validation failed (%s): %s", e.Field, e.Reason)
}
```

### Step 4: Create pkg/auth/doc.go

**Complete file content** (from Wave Implementation Plan lines 787-800):

```go
// Package auth provides interfaces for registry authentication.
//
// This package supports:
//   - Basic username/password authentication
//   - Credential validation
//   - Integration with go-containerregistry authn
//
// The primary interface is Provider, which supplies authentication
// to registry clients.
package auth
```

### Step 5: Create pkg/tls/interface.go

**Complete file content** (from Wave Implementation Plan lines 802-847):

```go
// Package tls provides interfaces and types for TLS configuration.
package tls

import (
	"crypto/tls"
)

// ConfigProvider defines operations for providing TLS configuration.
type ConfigProvider interface {
	// GetTLSConfig returns a tls.Config for HTTP transport.
	//
	// Returns:
	//   - *tls.Config: TLS configuration for HTTP transport
	GetTLSConfig() *tls.Config

	// IsInsecure returns whether insecure mode is enabled.
	//
	// Returns:
	//   - bool: true if --insecure flag was set, false otherwise
	IsInsecure() bool
}

// Config holds TLS configuration options.
type Config struct {
	// InsecureSkipVerify controls whether to skip TLS certificate verification.
	//
	// When true: Certificate validity NOT checked (development only)
	// When false: Full certificate verification (production)
	InsecureSkipVerify bool
}

// NewConfigProvider creates a TLS configuration provider.
//
// Parameters:
//   - insecure: Whether to enable insecure mode (skip cert verification)
//
// Returns:
//   - ConfigProvider: TLS configuration provider interface implementation
func NewConfigProvider(insecure bool) ConfigProvider {
	// Implementation will be provided in Wave 2 (pkg/tls/config.go)
	panic("not implemented - interface definition only")
}
```

### Step 6: Create pkg/tls/doc.go

**Complete file content** (from Wave Implementation Plan lines 849-859):

```go
// Package tls provides interfaces for TLS configuration.
//
// This package supports:
//   - TLS certificate verification (secure mode)
//   - Certificate verification bypass (insecure mode)
//   - HTTP transport configuration
package tls
```

### Step 7: Create Test Files

Create test files for both packages following the pattern from Effort 1.1.1.

### Step 8: Run Tests

```bash
cd /efforts/phase1/wave1/effort-3-auth-tls-interfaces
go test ./pkg/auth/... ./pkg/tls/... -v
```

### Step 9: Measure Size

```bash
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ]; then break; fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
cd /efforts/phase1/wave1/effort-3-auth-tls-interfaces
$PROJECT_ROOT/tools/line-counter.sh
```

### Step 10: Commit and Push

```bash
cd /efforts/phase1/wave1/effort-3-auth-tls-interfaces
git add pkg/auth/ pkg/tls/
git commit -m "feat(auth,tls): add Auth and TLS interface definitions

Implements Effort 1.1.3 - Auth & TLS Interface Definitions
Phase 1, Wave 1: Interface & Contract Definitions

Added:
Auth Package:
- Provider interface with 2 methods (GetAuthenticator, ValidateCredentials)
- Credentials struct for basic auth
- CredentialValidationError type
- Package documentation

TLS Package:
- ConfigProvider interface with 2 methods (GetTLSConfig, IsInsecure)
- Config struct for TLS settings
- Package documentation

Mock implementations for both interfaces
Complete test coverage (100%)

Implementation lines: ~140
Test coverage: 100%
All tests passing

Part of Wave 1 contract definition for Phase 1 Wave 2 implementations.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

git push origin idpbuilder-oci-push/phase1/wave1/effort-3-auth-tls-interfaces
```

---

## Dependencies

### Upstream Dependencies

**None** - Wave 1 efforts are independent

### Downstream Dependencies

- **Effort 1.1.2**: Registry Client (uses both interfaces)
- **Effort 1.1.4**: Command Structure (creates instances)

### External Library Dependencies

```go
require (
	github.com/google/go-containerregistry v0.19.0
	github.com/stretchr/testify v1.9.0
)
```

---

## Acceptance Criteria

- [ ] All 9 files created (5 auth + 4 tls)
- [ ] Both interfaces compile correctly
- [ ] Both mock implementations satisfy interfaces
- [ ] All tests passing (100% pass rate)
- [ ] Test coverage = 100%
- [ ] GoDoc complete for both packages
- [ ] Line count: 140±21 lines

---

## Document Status

**Status**: ✅ READY FOR IMPLEMENTATION
**Created**: 2025-10-29T04:08:10Z
**Planner**: @agent-code-reviewer
**Effort**: 1.1.3 - Auth & TLS Interface Definitions
**Phase**: 1, Wave: 1
**Fidelity**: EXACT (complete code provided)

**Lines**: 140 implementation + 190 test
**Coverage**: 100% required
**Dependencies**: None

---

**END OF IMPLEMENTATION PLAN - EFFORT 1.1.3**
