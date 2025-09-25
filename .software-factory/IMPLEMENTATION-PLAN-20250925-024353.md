# Implementation Plan for Effort 3.1.2 - Implement OCI Client

**Created**: 2025-09-25T02:43:53Z
**Location**: /home/vscode/workspaces/idpbuilder-push/efforts/phase3/wave1/implement-oci-client
**Phase**: 3
**Wave**: 1
**Branch**: idpbuilderpush/phase3/wave1/implement-oci-client
**Effort**: implement-oci-client

## Effort Metadata
- **Description**: GREEN phase implementation of OCI registry client using go-containerregistry
- **Size Estimate**: 400 lines of implementation code
- **Dependencies**: Effort 3.1.1 (Client Interface Tests) - COMPLETE
- **TDD Approach**: Tests already written in effort 3.1.1, now implementing to make them pass

## 🚨🚨🚨 R355 PRODUCTION READINESS - ZERO TOLERANCE 🚨🚨🚨

This implementation MUST be production-ready from the first commit:
- ❌ NO STUBS or placeholder implementations
- ❌ NO MOCKS except in test directories
- ❌ NO hardcoded credentials or secrets
- ❌ NO static configuration values
- ❌ NO TODO/FIXME markers in code
- ❌ NO returning nil or empty for "later implementation"
- ❌ NO panic("not implemented") patterns
- ❌ NO fake or dummy data

VIOLATION = -100% AUTOMATIC FAILURE

## Pre-Planning Research Results (R374 MANDATORY)

### Existing Interfaces Found
| Interface | Location | Signature | Must Implement |
|-----------|----------|-----------|----------------|
| `RegistryClient` | pkg/oci/client_test.go | Lines 11-17 | YES - PRIMARY |
| `Authenticator` | pkg/oci/types.go | Lines 9-18 | NO - Already exists |

### Existing Types to Reuse
| Component | Location | Purpose | How to Use |
|-----------|----------|---------|------------|
| `Credentials` | pkg/oci/types.go:21-27 | Auth creds struct | Import and use |
| `DefaultAuthenticator` | pkg/oci/auth.go:17-23 | Auth implementation | Use for auth |
| `ClientCredentials` | pkg/oci/client_test.go:20-27 | Test interface | Convert from/to |

### FORBIDDEN DUPLICATIONS (R373)
- DO NOT create new `Credentials` type - use existing from types.go
- DO NOT create new `Authenticator` interface - use existing from types.go
- DO NOT reimplement authentication logic - use `DefaultAuthenticator`

### REQUIRED INTEGRATIONS (R373)
- MUST implement `RegistryClient` interface from client_test.go with EXACT signature
- MUST use existing `Credentials` type from types.go
- MUST integrate with existing `DefaultAuthenticator` for auth

## EXPLICIT SCOPE (R311 MANDATORY)

### IMPLEMENT EXACTLY:
- Type: `OCIClient` struct with fields (~30 lines)
  - registry string
  - insecure bool
  - transport http.RoundTripper
  - auth Authenticator
  - mu sync.RWMutex
- Function: `NewRegistryClient() RegistryClient` (~15 lines)
- Method: `Connect(ctx context.Context, registry string) error` (~50 lines)
- Method: `Authenticate(credentials *ClientCredentials) error` (~40 lines)
- Method: `SetInsecure(insecure bool)` (~10 lines)
- Method: `GetTransport() http.RoundTripper` (~10 lines)
- Method: `Close() error` (~15 lines)
- Helper: `convertCredentials(cc *ClientCredentials) *Credentials` (~20 lines)
- Helper: `configureTransport(insecure bool) *http.Transport` (~30 lines)
- Helper: `validateURL(registry string) error` (~25 lines)

**TOTAL**: ~245 lines implementation + ~100 lines imports/comments = ~345 lines

### DO NOT IMPLEMENT:
- ❌ Push/Pull operations (future efforts)
- ❌ Manifest operations (future efforts)
- ❌ Layer operations (future efforts)
- ❌ Blob storage (future efforts)
- ❌ Registry catalog listing (future efforts)
- ❌ Tag operations (future efforts)
- ❌ Complex retry logic (keep simple exponential backoff)
- ❌ Caching beyond basic transport pooling
- ❌ Advanced proxy configuration
- ❌ Custom DNS resolution
- ❌ Certificate management
- ❌ Metrics or logging frameworks

## Size Limit Clarification (R359)
- The 800-line limit applies to NEW CODE YOU ADD
- This effort adds ~400 lines of new implementation
- NEVER delete existing code to meet size limits
- Repository will grow by ~400 lines (EXPECTED)

## Implementation Strategy

### Step 1: Add go-containerregistry dependency (10 lines)
```bash
cd /home/vscode/workspaces/idpbuilder-push/efforts/phase3/wave1/implement-oci-client
go get github.com/google/go-containerregistry@v0.20.2
```

### Step 2: Create pkg/oci/client.go (345 lines)

#### Structure:
```go
package oci

import (
    "context"
    "crypto/tls"
    "errors"
    "fmt"
    "net/http"
    "net/url"
    "strings"
    "sync"
    "time"

    "github.com/google/go-containerregistry/pkg/authn"
    "github.com/google/go-containerregistry/pkg/name"
    "github.com/google/go-containerregistry/pkg/v1/remote"
)

// OCIClient implements the RegistryClient interface
type OCIClient struct {
    registry  string
    insecure  bool
    transport http.RoundTripper
    auth      Authenticator
    mu        sync.RWMutex
}

// Implementation methods follow...
```

### Step 3: Implement Core Methods

#### 3.1 Constructor (15 lines)
```go
func NewRegistryClient() RegistryClient {
    return &OCIClient{
        auth: NewAuthenticator(nil), // Use existing DefaultAuthenticator
    }
}
```

#### 3.2 Connect Method (50 lines)
- Validate URL format
- Configure transport based on insecure flag
- Set up connection pooling
- Handle context cancellation

#### 3.3 Authenticate Method (40 lines)
- Convert ClientCredentials to internal Credentials
- Store in authenticator cache
- Validate credentials format

#### 3.4 Helper Methods (85 lines total)
- `convertCredentials`: Map test interface types to internal types
- `configureTransport`: Set up HTTP transport with TLS/pooling
- `validateURL`: Ensure URL is valid HTTP/HTTPS

### Step 4: Run Tests and Verify GREEN Phase
```bash
# Remove Skip() from tests and run
go test ./pkg/oci -v -run TestRegistryClient
```

## Configuration Requirements (R355 Mandatory)

### WRONG - Will fail review:
```go
// ❌ VIOLATION - Hardcoded values
registryURL := "https://docker.io"
timeout := 30

// ❌ VIOLATION - Stub implementation
func Connect(ctx context.Context, registry string) error {
    // TODO: implement
    return nil
}
```

### CORRECT - Production ready:
```go
// ✅ From parameters
func Connect(ctx context.Context, registry string) error {
    if err := validateURL(registry); err != nil {
        return fmt.Errorf("invalid registry URL: %w", err)
    }
    // Full implementation...
}

// ✅ Configurable timeout from context
timeout := 30 * time.Second
if deadline, ok := ctx.Deadline(); ok {
    timeout = time.Until(deadline)
}
```

## Atomic PR Design (R220)

### PR Summary
**Single PR**: "feat: implement OCI registry client with go-containerregistry integration"

### Can Merge to Main Alone
✅ YES - This PR can merge independently:
- Implements complete RegistryClient interface
- All methods fully functional
- No dependencies on unmerged work
- Tests pass independently

### R355 Production Ready Checklist
- ✅ No hardcoded values - all config from parameters
- ✅ All functions complete - no stubs
- ✅ No TODO markers
- ✅ Full error handling
- ✅ Context cancellation support

### Interface Implementations
- **RegistryClient**: Complete implementation with all methods
- **Integration with Authenticator**: Uses existing auth system
- **Production Ready**: Fully functional, not a stub

### PR Verification
- Tests pass alone: ✅ (effort 3.1.1 tests will turn GREEN)
- Build remains working: ✅
- No external dependencies: ✅ (uses existing auth)
- Backward compatible: ✅ (new functionality)

### Files in This PR
- `pkg/oci/client.go` (NEW - ~345 lines)
- `go.mod` (MODIFIED - add go-containerregistry)
- `go.sum` (MODIFIED - dependency hashes)

## Test Requirements

### Unit Tests Status
- Tests already written in effort 3.1.1 (client_test.go)
- Currently all tests have `t.Skip("TDD RED: Client implementation does not exist yet")`
- Implementation will make tests pass (GREEN phase)
- Test coverage target: >80%

### Test Execution Plan
1. Remove `t.Skip()` calls from client_test.go
2. Run tests to verify implementation
3. All 15 test cases should pass

## Size Management
- **Estimated Lines**: ~400 (well under 800 limit)
- **Measurement Tool**: Use line-counter.sh after implementation
- **Check Command**: `$PROJECT_ROOT/tools/line-counter.sh`
- **Split Threshold**: Not needed (400 < 800)

## Implementation Checklist
- [ ] Add go-containerregistry dependency to go.mod
- [ ] Create pkg/oci/client.go with OCIClient struct
- [ ] Implement NewRegistryClient constructor
- [ ] Implement Connect method with URL validation
- [ ] Implement Authenticate with credential conversion
- [ ] Implement SetInsecure flag handling
- [ ] Implement GetTransport accessor
- [ ] Implement Close cleanup method
- [ ] Add helper functions (convertCredentials, configureTransport, validateURL)
- [ ] Remove t.Skip() from tests
- [ ] Verify all tests pass
- [ ] Run line-counter.sh to verify size

## Success Criteria
✅ All 15 tests from effort 3.1.1 pass
✅ No hardcoded values or stubs
✅ Implementation under 400 lines
✅ Full production-ready code
✅ Integrates with existing auth system
✅ Uses go-containerregistry properly

## Next Steps
1. SW Engineer reads this plan from .software-factory/
2. Implements client.go following TDD GREEN phase
3. Makes all tests from effort 3.1.1 pass
4. Commits and pushes implementation
5. Ready for code review
