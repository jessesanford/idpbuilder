# Effort 1.2.4: TLS Configuration Implementation - Implementation Plan

**Created**: 2025-10-29 06:33:00 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**Effort ID**: 1.2.4
**Phase**: Phase 1 - Foundation & Interfaces
**Wave**: Wave 2 - Core Package Implementations

---

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)

### R213 Infrastructure Metadata

```json
{
  "effort_id": "1.2.4",
  "effort_name": "TLS Configuration Implementation",
  "branch_name": "idpbuilder-oci-push/phase1/wave2/effort-4-tls",
  "base_branch": "idpbuilder-oci-push/phase1/wave2/integration",
  "parent_wave": "WAVE_2",
  "parent_phase": "PHASE_1",
  "depends_on": [],
  "estimated_lines": 350,
  "complexity": "medium",
  "can_parallelize": true,
  "parallel_with": ["1.2.1", "1.2.2", "1.2.3"]
}
```

**Branch**: `idpbuilder-oci-push/phase1/wave2/effort-4-tls`
**Can Parallelize**: Yes
**Parallel With**: [1.2.1, 1.2.2, 1.2.3] (ALL Wave 2 efforts run simultaneously)
**Size Estimate**: 350 lines (MUST be <800)
**Dependencies**: None (all Wave 1 interfaces frozen and available)

---

## Overview

**Purpose**: Implement the TLS configuration provider that creates TLS configs for HTTP transports, supports secure and insecure modes, and loads system certificate pools.

**What This Effort Accomplishes**:
- Complete implementation of `tls.ConfigProvider` interface (frozen in Wave 1)
- TLS configuration for HTTP transports
- Secure mode with system certificate verification
- Insecure mode with certificate verification disabled (--insecure flag)
- System certificate pool loading
- Mode detection (IsInsecure flag)

**Boundaries - OUT OF SCOPE**:
- Certificate generation or management
- Custom CA certificate loading (system certs only)
- Client certificate authentication (mTLS)
- Certificate pinning
- TLS version or cipher suite customization

---

## File Structure

### New Files to Create

**Implementation Files**:
- `pkg/tls/config.go` (~350 lines)
  - `tlsConfigProvider` struct implementation
  - `NewConfigProvider()` constructor
  - `GetTLSConfig()` method returning *tls.Config
  - `IsInsecure()` mode detection method

**Test Files** (NOT counted in line estimates per R007):
- `pkg/tls/config_test.go` (~200 lines)
  - 8+ test cases covering all methods
  - Secure and insecure mode tests
  - HTTP client integration tests
  - Success paths

**Modified Files**:
- None (standard library only: crypto/tls, crypto/x509)

**Total Estimated Lines**: 350 lines (implementation only, tests excluded per R007)

---

## Implementation Steps

### Step 1: Review Wave 1 Interface Definition

**MANDATORY FIRST STEP - Read frozen interfaces**:
```bash
# Read the frozen TLS interface from Wave 1
cat pkg/tls/types.go
```

**Expected TLS Interface** (from Wave 1 Effort 3):
```go
package tls

import "crypto/tls"

type ConfigProvider interface {
    GetTLSConfig() *tls.Config
    IsInsecure() bool
}

type Config struct {
    InsecureSkipVerify bool
}
```

### Step 2: Implement pkg/tls/config.go

**File: pkg/tls/config.go**

**Required Implementation Details**:

**1. Struct Definition**:
```go
type tlsConfigProvider struct {
    config Config
}
```

**2. NewConfigProvider() Implementation** (~50 lines):
- Accept `insecure` boolean parameter
- Store in `Config` struct: `Config{InsecureSkipVerify: insecure}`
- Return `&tlsConfigProvider{config: cfg}`
- No validation needed (boolean parameter)
- Return implementation of `ConfigProvider` interface

**Example**:
```go
func NewConfigProvider(insecure bool) ConfigProvider {
    return &tlsConfigProvider{
        config: Config{
            InsecureSkipVerify: insecure,
        },
    }
}
```

**3. GetTLSConfig() Implementation** (~120 lines):
- Create `tls.Config` with `InsecureSkipVerify` from stored config
- **If secure mode** (InsecureSkipVerify == false):
  - Load system certificate pool: `certPool, err := x509.SystemCertPool()`
  - If system cert pool unavailable (err != nil):
    - Create new empty cert pool: `certPool = x509.NewCertPool()`
    - Log warning (optional): "System certificate pool unavailable, using empty pool"
  - Set `RootCAs` to cert pool
- **If insecure mode** (InsecureSkipVerify == true):
  - No cert pool needed (verification disabled)
  - RootCAs can be nil
- Return `*tls.Config` ready for HTTP transport

**Example**:
```go
func (p *tlsConfigProvider) GetTLSConfig() *tls.Config {
    tlsConfig := &tls.Config{
        InsecureSkipVerify: p.config.InsecureSkipVerify,
    }

    if !p.config.InsecureSkipVerify {
        // Secure mode: Load system certificates
        certPool, err := x509.SystemCertPool()
        if err != nil {
            // Fallback to empty pool if system certs unavailable
            certPool = x509.NewCertPool()
        }
        tlsConfig.RootCAs = certPool
    }

    return tlsConfig
}
```

**4. IsInsecure() Implementation** (~20 lines):
- Return `config.InsecureSkipVerify` boolean
- Simple getter method
- No validation needed

**Example**:
```go
func (p *tlsConfigProvider) IsInsecure() bool {
    return p.config.InsecureSkipVerify
}
```

**Usage Pattern**:
```go
// Secure mode (default)
provider := tls.NewConfigProvider(false)
tlsConfig := provider.GetTLSConfig()

// Use with HTTP client
transport := &http.Transport{
    TLSClientConfig: tlsConfig,
}
client := &http.Client{Transport: transport}

// Insecure mode (--insecure flag)
insecureProvider := tls.NewConfigProvider(true)
insecureTLSConfig := insecureProvider.GetTLSConfig()
```

**Complete Package Structure**:
```go
package tls

import (
    "crypto/tls"
    "crypto/x509"
)

// Implementation here (~350 lines total including comments and spacing)
```

### Step 3: Write Tests (TDD - Tests First!)

**File: pkg/tls/config_test.go**

**Test Cases** (from Wave 2 Test Plan):

**A. Constructor Tests**:
- `TestNewConfigProvider_SecureMode`: NewConfigProvider creates secure provider (insecure=false)
- `TestNewConfigProvider_InsecureMode`: NewConfigProvider creates insecure provider (insecure=true)

**B. GetTLSConfig Tests**:
- `TestGetTLSConfig_SecureMode`: GetTLSConfig in secure mode
  - Verify `InsecureSkipVerify == false`
  - Verify `RootCAs` loaded (system cert pool)
  - Verify `RootCAs` not nil
- `TestGetTLSConfig_InsecureMode`: GetTLSConfig in insecure mode
  - Verify `InsecureSkipVerify == true`
  - RootCAs may be nil (no verification needed)

**C. IsInsecure Tests**:
- `TestIsInsecure_SecureMode`: IsInsecure returns false for secure mode
- `TestIsInsecure_InsecureMode`: IsInsecure returns true for insecure mode

**D. Integration Tests**:
- `TestTLSConfig_HTTPClientIntegration`: TLS config usable with http.Client
  - Create HTTP transport with TLS config
  - Verify transport has correct TLS settings
  - Test both secure and insecure modes

**Test Coverage Requirements**:
- Minimum 90% code coverage (security-critical, per Wave 2 Test Plan)
- All success paths tested
- Secure and insecure modes tested
- System cert pool loading tested
- HTTP client integration tested

**Test Examples**:
```go
func TestGetTLSConfig_SecureMode(t *testing.T) {
    provider := NewConfigProvider(false)
    config := provider.GetTLSConfig()

    assert.NotNil(t, config)
    assert.False(t, config.InsecureSkipVerify)
    assert.NotNil(t, config.RootCAs, "Secure mode should load certificate pool")
}

func TestGetTLSConfig_InsecureMode(t *testing.T) {
    provider := NewConfigProvider(true)
    config := provider.GetTLSConfig()

    assert.NotNil(t, config)
    assert.True(t, config.InsecureSkipVerify)
}

func TestTLSConfig_HTTPClientIntegration(t *testing.T) {
    provider := NewConfigProvider(false)
    tlsConfig := provider.GetTLSConfig()

    transport := &http.Transport{
        TLSClientConfig: tlsConfig,
    }
    client := &http.Client{Transport: transport}

    assert.NotNil(t, client)
    assert.NotNil(t, client.Transport)
    assert.Equal(t, tlsConfig, transport.TLSClientConfig)
}
```

### Step 4: Size Measurement

**Measure implementation lines**:
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning
tools/line-counter.sh

# Expected output:
# 🎯 Detected base: idpbuilder-oci-push/phase1/wave2/integration
# 📦 Analyzing branch: idpbuilder-oci-push/phase1/wave2/effort-4-tls
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
cd pkg/tls
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
go vet ./pkg/tls/...
golangci-lint run ./pkg/tls/...

# Expected: No errors
```

**Verify GoDoc**:
- All public types have GoDoc comments
- All public functions have GoDoc comments
- Security warnings documented (insecure mode risks)

**Example GoDoc**:
```go
// NewConfigProvider creates a new TLS configuration provider.
//
// Parameters:
//   - insecure: If true, disables certificate verification (INSECURE - use only for testing)
//
// Returns a ConfigProvider that can generate *tls.Config for HTTP transports.
//
// WARNING: Setting insecure=true disables ALL certificate verification.
// This makes connections vulnerable to man-in-the-middle attacks.
// Only use insecure mode for testing with self-signed certificates.
func NewConfigProvider(insecure bool) ConfigProvider {
    // ...
}
```

### Step 7: Commit and Push

**Commit structure**:
```bash
git add pkg/tls/config.go pkg/tls/config_test.go
git commit -m "feat(tls): implement TLS configuration provider

- Implement NewConfigProvider constructor
- Implement GetTLSConfig with secure/insecure modes
- Implement IsInsecure mode detection
- Load system certificate pool for secure mode
- Support insecure mode for self-signed certificates
- Add 8 test cases with 90%+ coverage (security-critical)
- Add HTTP client integration tests

Closes: Effort 1.2.4 - TLS Configuration Implementation
Lines: ~350 (within 350 ±15% estimate)
Coverage: 90%+ (meets Wave 2 Test Plan requirements)"

git push origin idpbuilder-oci-push/phase1/wave2/effort-4-tls
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
- After struct definition: ~50 lines
- After NewConfigProvider: ~100 lines
- After GetTLSConfig: ~220 lines
- After IsInsecure: ~240 lines
- After comments and documentation: ~350 lines (target)

---

## Test Requirements

**Minimum Coverage**: 90% (security-critical, per Wave 2 Test Plan)

**Test Categories** (from Wave 2 Test Plan):

| Test Category | Test Cases | Coverage Target |
|---------------|------------|-----------------|
| Constructor | 2 tests | 100% of NewConfigProvider |
| GetTLSConfig | 2 tests | 100% of GetTLSConfig |
| IsInsecure | 2 tests | 100% of IsInsecure |
| Integration | 2 tests | HTTP client integration |

**Total Test Cases**: 8+ tests

**Test Execution**:
```bash
go test ./pkg/tls -v -cover
# MUST achieve ≥90% coverage (security-critical)
```

---

## Dependencies

### Upstream Dependencies (COMPLETED)
- ✅ Wave 1 Effort 3: TLS interface definition (FROZEN)
- ✅ Integration branch: `idpbuilder-oci-push/phase1/wave2/integration` (CREATED)

### Downstream Dependencies
- None (all Wave 2 efforts are parallel)
- Effort 1.2.2 (Registry Client) will USE this implementation via interface
- Wave 3 CLI will use this package

### External Library Dependencies
**Standard Library Only**:
- `crypto/tls` for TLS configuration
- `crypto/x509` for certificate pool

**Test Dependencies**:
- `github.com/stretchr/testify` v1.10.0 (already in Wave 1, for tests)

**No New Dependencies Required**

---

## Pattern Compliance

### Go Patterns
- Interface-driven design (implement `tls.ConfigProvider` interface)
- Struct-based provider implementation
- No global state (all data in struct)
- Standard library usage (no external dependencies)

### Security Requirements
- **Secure by default**: Default mode should be secure (insecure=false)
- **System certificates**: Load and use system cert pool
- **Insecure mode warnings**: Document risks clearly
- **Fallback handling**: Graceful fallback if system certs unavailable

### Performance Targets
- TLS config creation should complete in <1 millisecond
- System cert pool loading is one-time cost
- No expensive operations during normal usage

---

## Acceptance Criteria

**MANDATORY - All must pass before Code Review**:

- [ ] All files created/modified as specified
- [ ] All 2 interface methods implemented correctly (GetTLSConfig, IsInsecure)
- [ ] All tests passing (100% pass rate)
- [ ] Code coverage ≥90% (security-critical, per Wave 2 Test Plan)
- [ ] No linting errors (go vet, golangci-lint)
- [ ] Documentation complete (all public methods have GoDoc)
- [ ] Line count within estimate (350 lines ±15% = 298-403 lines)
- [ ] Secure mode loads system certificate pool
- [ ] Insecure mode disables verification correctly
- [ ] TLS config compatible with HTTP transport
- [ ] No security warnings for secure mode
- [ ] Security warnings documented for insecure mode
- [ ] Code committed and pushed to effort branch

**Quality Gates**:
1. **Functionality**: All interface methods work correctly
2. **Testing**: 90%+ coverage with all paths tested
3. **Security**: Secure mode is default and properly configured
4. **Size**: Within 350 ±15% lines (298-403)
5. **Documentation**: Complete GoDoc coverage with security warnings

---

## References

**Wave 2 Planning Documents**:
- Wave Implementation Plan: `/home/vscode/workspaces/idpbuilder-oci-push-planning/wave-plans/WAVE-2-IMPLEMENTATION.md`
- Wave Architecture: `/home/vscode/workspaces/idpbuilder-oci-push-planning/wave-plans/WAVE-2-ARCHITECTURE.md`
- Wave Test Plan: `/home/vscode/workspaces/idpbuilder-oci-push-planning/wave-plans/WAVE-2-TEST-PLAN.md`

**Wave 1 Interfaces** (frozen references):
- TLS Interface: `efforts/phase1/wave1/effort-3-auth-tls-interfaces/pkg/tls/types.go`

**External Documentation**:
- Go crypto/tls: https://pkg.go.dev/crypto/tls
- Go crypto/x509: https://pkg.go.dev/crypto/x509
- TLS Best Practices: https://golang.org/doc/go1.17#crypto/tls

---

## Implementation Checklist

**Pre-Implementation**:
- [ ] Read Wave 1 TLS interface definition
- [ ] Read Wave 2 Architecture (TLS section)
- [ ] Read Wave 2 Test Plan (TLS test cases)
- [ ] Checkout effort branch from integration
- [ ] Verify base branch is correct

**Implementation Phase**:
- [ ] Write test stubs (8+ test cases)
- [ ] Implement `tlsConfigProvider` struct
- [ ] Implement `NewConfigProvider()` constructor
- [ ] Implement `GetTLSConfig()` method (secure mode)
- [ ] Implement `GetTLSConfig()` method (insecure mode)
- [ ] Implement `IsInsecure()` method
- [ ] Complete test implementations
- [ ] Run tests (all pass, 90%+ coverage)
- [ ] Run linters (no errors)
- [ ] Add GoDoc comments with security warnings

**Validation Phase**:
- [ ] Measure size (within 298-403 lines)
- [ ] Verify coverage ≥90%
- [ ] System cert pool loading validation
- [ ] Secure/insecure mode validation
- [ ] HTTP client integration validation
- [ ] Security warning documentation check
- [ ] Commit and push code

---

## Security Considerations

### Secure Mode (Default)
- ✅ **System certificates loaded**: Uses OS trust store
- ✅ **Certificate verification enabled**: Validates server certificates
- ✅ **Man-in-the-middle protection**: Full TLS security
- ✅ **Production ready**: Safe for production use

### Insecure Mode (--insecure flag)
- ⚠️ **Certificate verification DISABLED**: No certificate validation
- ⚠️ **Vulnerable to MITM attacks**: No protection against interception
- ⚠️ **Testing only**: Should NEVER be used in production
- ⚠️ **Explicit opt-in**: Requires explicit flag to enable

**Documentation MUST clearly warn about insecure mode risks!**

---

## Next Steps

**After Implementation Completion**:
1. Code Reviewer will be spawned for effort review
2. Code Reviewer validates all acceptance criteria
3. If approved: Merge to integration branch
4. If issues found: Fix and re-submit for review
5. After all 4 Wave 2 efforts approved: Wave integration testing

**Parallel Work**:
- This effort (1.2.4) runs in parallel with:
  - Effort 1.2.1: Docker Client Implementation
  - Effort 1.2.2: Registry Client Implementation
  - Effort 1.2.3: Authentication Implementation

**No coordination needed** - all efforts are independent until integration phase.

---

## Document Status

**Status**: ✅ READY FOR IMPLEMENTATION
**Created**: 2025-10-29 06:33:00 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**Effort**: 1.2.4 (TLS Configuration Implementation)
**Wave**: Wave 2 of Phase 1
**Branch**: `idpbuilder-oci-push/phase1/wave2/effort-4-tls`
**Base Branch**: `idpbuilder-oci-push/phase1/wave2/integration`

**Compliance**:
- ✅ R213: Complete metadata included
- ✅ R211: Parallelization specified (runs with 1.2.1, 1.2.2, 1.2.3)
- ✅ R341: TDD approach (test plan before implementation)
- ✅ R383: Plan stored in .software-factory with timestamp
- ✅ Size compliance: 350 lines < 800 hard limit
- ✅ Security-critical: 90% coverage requirement
- ✅ Standard library only: No external dependencies

---

**END OF EFFORT 1.2.4 IMPLEMENTATION PLAN**
