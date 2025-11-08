# TLS Configuration Implementation - Implementation Plan

**Effort**: 1.2.4 - TLS Configuration Implementation
**Phase**: Phase 1 - Foundation & Interfaces
**Wave**: Wave 2 - Core Package Implementations
**Created**: 2025-10-29 21:33:16 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**State**: EFFORT_PLAN_CREATION

---

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)

**Effort ID**: 1.2.4
**Branch**: `idpbuilder-oci-push/phase1/wave2/effort-4-tls`
**Base Branch**: `idpbuilder-oci-push/phase1/wave2/integration`
**Can Parallelize**: Yes
**Parallel With**: 1.2.1 (Docker Client), 1.2.2 (Registry Client), 1.2.3 (Auth Implementation)
**Size Estimate**: 350 lines (implementation only, tests excluded per R007)
**Complexity**: Medium
**Dependencies**: Wave 1 Effort 3 (TLS interface definition) - COMPLETED

---

## Overview

### Purpose

Implement the TLS configuration provider that creates TLS configs for HTTP transports, supports secure and insecure modes, and loads system certificate pools.

### What This Effort Accomplishes

- Complete implementation of `tls.ConfigProvider` interface (frozen in Wave 1)
- TLS configuration for HTTP transports
- Secure mode with system certificate verification
- Insecure mode with certificate verification disabled (--insecure flag)
- System certificate pool loading
- Mode detection (IsInsecure flag)

### Boundaries - OUT OF SCOPE

- ❌ Certificate generation or management
- ❌ Custom CA certificate loading (system certs only)
- ❌ Client certificate authentication (mTLS)
- ❌ Certificate pinning
- ❌ TLS version or cipher suite customization

---

## File Structure

### New Files to Create

#### 1. `pkg/tls/config.go` (~350 lines)

**Purpose**: Implementation of TLS configuration provider

**Contents**:
- `tlsConfigProvider` struct implementation
- `NewConfigProvider()` constructor
- `GetTLSConfig()` method returning *tls.Config
- `IsInsecure()` mode detection method

**Key Components**:
```go
type tlsConfigProvider struct {
    config Config  // From Wave 1 interface
}
```

#### 2. `pkg/tls/config_test.go` (~200 lines, NOT counted per R007)

**Purpose**: Comprehensive unit tests for TLS configuration

**Test Categories** (from Wave 2 Test Plan):
- Constructor tests (secure/insecure modes)
- GetTLSConfig tests (system cert pool, verification flags)
- IsInsecure tests (mode detection)
- Integration tests (HTTP client compatibility)

**Total Test Cases**: 8+ tests
**Coverage Target**: ≥90% (security-critical per Wave 2 Test Plan)

### Files to Modify

**None** - This implementation uses standard library only:
- `crypto/tls` for TLS configuration
- `crypto/x509` for certificate pool

No external dependencies or go.mod changes required.

---

## Implementation Steps

### Step 1: Set Up Package Structure

**Goal**: Create pkg/tls directory and initialize package

**Tasks**:
1. Create `pkg/tls/` directory if it doesn't exist
2. Verify Wave 1 interface files are accessible:
   - `pkg/tls/interface.go` (from Wave 1 Effort 3)
   - `pkg/tls/types.go` (Config struct)

**Verification**:
```bash
# Verify Wave 1 interfaces available
ls -la pkg/tls/interface.go pkg/tls/types.go

# Should show:
# pkg/tls/interface.go (ConfigProvider interface)
# pkg/tls/types.go (Config struct)
```

### Step 2: Implement tlsConfigProvider Struct

**Goal**: Create the internal struct implementing ConfigProvider interface

**Implementation** (from Wave 2 Architecture):
```go
// Package tls provides TLS configuration implementations.
package tls

import (
    "crypto/tls"
    "crypto/x509"
)

// tlsConfigProvider implements the ConfigProvider interface.
type tlsConfigProvider struct {
    config Config  // From Wave 1: Config{InsecureSkipVerify bool}
}
```

**Key Points**:
- Stores Wave 1 `Config` struct
- Internal struct (lowercase name)
- No exported fields

### Step 3: Implement NewConfigProvider Constructor

**Goal**: Implement constructor for TLS configuration provider

**Specification** (from Wave 2 Implementation Plan):

**Function Signature**:
```go
func NewConfigProvider(insecure bool) ConfigProvider
```

**Implementation Requirements** (from Wave 2 Architecture):
```go
// NewConfigProvider creates a TLS configuration provider.
//
// Parameters:
//   - insecure: Whether to enable insecure mode (skip cert verification)
//               Typically set from --insecure / -k CLI flag
//
// Returns:
//   - ConfigProvider: TLS configuration provider interface implementation
//
// Example:
//   // Secure mode (default)
//   provider := tls.NewConfigProvider(false)
//
//   // Insecure mode (--insecure flag)
//   provider := tls.NewConfigProvider(true)
//   if provider.IsInsecure() {
//       fmt.Println("WARNING: TLS verification disabled")
//   }
func NewConfigProvider(insecure bool) ConfigProvider {
    return &tlsConfigProvider{
        config: Config{
            InsecureSkipVerify: insecure,
        },
    }
}
```

**Key Points**:
- Accept `insecure` boolean parameter
- Store in `Config` struct from Wave 1
- No validation needed (boolean parameter is always valid)
- Return interface type `ConfigProvider`

**Tests to Pass** (TC-TLS-IMPL-001, TC-TLS-IMPL-002):
- NewConfigProvider with `insecure=false` creates secure mode provider
- NewConfigProvider with `insecure=true` creates insecure mode provider

### Step 4: Implement GetTLSConfig Method

**Goal**: Implement TLS configuration generation for HTTP transport

**Specification** (from Wave 2 Implementation Plan):

**Function Signature**:
```go
func (p *tlsConfigProvider) GetTLSConfig() *tls.Config
```

**Implementation Requirements** (from Wave 2 Architecture):
```go
// GetTLSConfig returns a tls.Config for HTTP transport.
//
// Behavior depends on insecure mode:
//   - Insecure mode (--insecure flag): InsecureSkipVerify = true
//   - Secure mode (default): Uses system certificate pool
//
// The returned tls.Config is used by the HTTP client when connecting
// to the registry.
//
// Returns:
//   - *tls.Config: TLS configuration for HTTP transport
//
// Example (secure mode):
//   provider := tls.NewConfigProvider(false)
//   tlsConfig := provider.GetTLSConfig()
//   transport := &http.Transport{
//       TLSClientConfig: tlsConfig,
//   }
//
// Example (insecure mode):
//   provider := tls.NewConfigProvider(true)
//   tlsConfig := provider.GetTLSConfig()
//   // tlsConfig.InsecureSkipVerify == true
func (p *tlsConfigProvider) GetTLSConfig() *tls.Config {
    tlsConfig := &tls.Config{
        InsecureSkipVerify: p.config.InsecureSkipVerify,
    }

    // If secure mode, load system certificate pool
    if !p.config.InsecureSkipVerify {
        // Load system certificates
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

**Key Points**:
- Create `tls.Config` with `InsecureSkipVerify` from stored config
- **Secure mode** (`InsecureSkipVerify == false`):
  - Load system certificate pool using `x509.SystemCertPool()`
  - Fallback to empty pool if system certs unavailable
  - Set `RootCAs` to system cert pool
- **Insecure mode** (`InsecureSkipVerify == true`):
  - No cert pool needed (verification disabled)
- Return `*tls.Config` ready for HTTP transport

**Usage Pattern** (from Wave 2 Architecture):
```go
provider := tls.NewConfigProvider(insecure)
tlsConfig := provider.GetTLSConfig()
transport := &http.Transport{
    TLSClientConfig: tlsConfig,
}
client := &http.Client{Transport: transport}
```

**Tests to Pass** (TC-TLS-IMPL-003, TC-TLS-IMPL-004):
- GetTLSConfig in secure mode has `InsecureSkipVerify == false` and system cert pool loaded
- GetTLSConfig in insecure mode has `InsecureSkipVerify == true`

### Step 5: Implement IsInsecure Method

**Goal**: Implement mode detection method

**Specification** (from Wave 2 Implementation Plan):

**Function Signature**:
```go
func (p *tlsConfigProvider) IsInsecure() bool
```

**Implementation Requirements** (from Wave 2 Architecture):
```go
// IsInsecure returns whether insecure mode is enabled.
//
// Returns:
//   - bool: true if --insecure flag was set, false otherwise
//
// Example:
//   if provider.IsInsecure() {
//       log.Warn("TLS certificate verification disabled (insecure mode)")
//   }
func (p *tlsConfigProvider) IsInsecure() bool {
    return p.config.InsecureSkipVerify
}
```

**Key Points**:
- Simple getter method
- Returns `Config.InsecureSkipVerify` boolean
- Used by callers to warn about insecure mode

**Tests to Pass** (TC-TLS-IMPL-005, TC-TLS-IMPL-006):
- IsInsecure returns `false` for secure mode
- IsInsecure returns `true` for insecure mode

### Step 6: Write Comprehensive Unit Tests

**Goal**: Achieve ≥90% coverage with comprehensive tests

**Test File**: `pkg/tls/config_test.go`

**Test Structure** (from Wave 2 Test Plan):

#### A. Constructor Tests (2 tests)

**TC-TLS-IMPL-001: NewConfigProvider_SecureMode**
```go
func TestNewConfigProvider_SecureMode(t *testing.T) {
    // Given: Secure mode (false = verify certificates)
    insecure := false

    // When: Creating provider
    provider := NewConfigProvider(insecure)

    // Then: Provider created in secure mode
    require.NotNil(t, provider)
    var _ ConfigProvider = provider  // Interface check
    assert.False(t, provider.IsInsecure())
}
```

**TC-TLS-IMPL-002: NewConfigProvider_InsecureMode**
```go
func TestNewConfigProvider_InsecureMode(t *testing.T) {
    // Given: Insecure mode (true = skip verification)
    insecure := true

    // When: Creating provider
    provider := NewConfigProvider(insecure)

    // Then: Provider created in insecure mode
    require.NotNil(t, provider)
    assert.True(t, provider.IsInsecure())
}
```

#### B. GetTLSConfig Tests (2 tests)

**TC-TLS-IMPL-003: GetTLSConfig_SecureMode**
```go
func TestGetTLSConfig_SecureMode(t *testing.T) {
    // Given: Provider in secure mode
    provider := NewConfigProvider(false)

    // When: Getting TLS config
    tlsConfig := provider.GetTLSConfig()

    // Then: Config has proper settings
    require.NotNil(t, tlsConfig)
    assert.False(t, tlsConfig.InsecureSkipVerify)
    assert.NotNil(t, tlsConfig.RootCAs, "Should have system cert pool")
}
```

**TC-TLS-IMPL-004: GetTLSConfig_InsecureMode**
```go
func TestGetTLSConfig_InsecureMode(t *testing.T) {
    // Given: Provider in insecure mode
    provider := NewConfigProvider(true)

    // When: Getting TLS config
    tlsConfig := provider.GetTLSConfig()

    // Then: InsecureSkipVerify is true
    require.NotNil(t, tlsConfig)
    assert.True(t, tlsConfig.InsecureSkipVerify)
}
```

#### C. IsInsecure Tests (2 tests)

**TC-TLS-IMPL-005: IsInsecure_Secure**
```go
func TestIsInsecure_Secure(t *testing.T) {
    provider := NewConfigProvider(false)
    assert.False(t, provider.IsInsecure())
}
```

**TC-TLS-IMPL-006: IsInsecure_Insecure**
```go
func TestIsInsecure_Insecure(t *testing.T) {
    provider := NewConfigProvider(true)
    assert.True(t, provider.IsInsecure())
}
```

#### D. Integration Tests (1 test)

**TC-TLS-IMPL-007: TLSConfig_UsableWithHTTPClient**
```go
func TestTLSConfig_UsableWithHTTPClient(t *testing.T) {
    // Given: TLS provider
    provider := NewConfigProvider(true)
    tlsConfig := provider.GetTLSConfig()

    // When: Creating HTTP transport with TLS config
    transport := &http.Transport{
        TLSClientConfig: tlsConfig,
    }
    client := &http.Client{Transport: transport}

    // Then: Client created successfully
    require.NotNil(t, client)
    assert.Equal(t, tlsConfig, transport.TLSClientConfig)
}
```

**Total Tests**: 8+ test cases
**Coverage Target**: ≥90% (security-critical)

**Test Execution**:
```bash
# Run TLS tests with coverage
go test ./pkg/tls -v -cover

# Expected output:
# === RUN   TestNewConfigProvider_SecureMode
# --- PASS: TestNewConfigProvider_SecureMode (0.00s)
# ...
# PASS
# coverage: 91.2% of statements in ./pkg/tls
```

### Step 7: Add GoDoc Documentation

**Goal**: Complete documentation for all public methods

**Documentation Requirements** (from Wave 2 Implementation Plan):
- All public functions must have GoDoc comments
- Explain parameters, return values, behavior
- Include usage examples
- Document security considerations

**Example GoDoc Structure**:
```go
// NewConfigProvider creates a TLS configuration provider.
//
// Parameters:
//   - insecure: Whether to enable insecure mode (skip cert verification)
//
// Returns:
//   - ConfigProvider: TLS configuration provider interface implementation
//
// Example:
//   provider := tls.NewConfigProvider(false)
```

**Documentation Checklist**:
- [ ] Package-level documentation
- [ ] NewConfigProvider function
- [ ] GetTLSConfig method
- [ ] IsInsecure method
- [ ] Struct comments (if needed)

### Step 8: Verify Implementation

**Goal**: Ensure all requirements met

**Verification Checklist**:

#### Interface Compliance
- [ ] `tlsConfigProvider` implements `ConfigProvider` interface
- [ ] All 2 interface methods implemented:
  - `GetTLSConfig() *tls.Config`
  - `IsInsecure() bool`
- [ ] Uses Wave 1 `Config` struct correctly

#### Functionality
- [ ] Secure mode loads system certificate pool
- [ ] Insecure mode disables verification correctly
- [ ] TLS config compatible with HTTP transport
- [ ] No security warnings for secure mode

#### Testing
- [ ] All 8+ tests passing
- [ ] Coverage ≥90% achieved
- [ ] All test categories covered:
  - Constructor tests (2)
  - GetTLSConfig tests (2)
  - IsInsecure tests (2)
  - Integration tests (1+)

#### Code Quality
- [ ] No linting errors: `go vet ./pkg/tls`
- [ ] No golangci-lint errors: `golangci-lint run ./pkg/tls`
- [ ] Code formatted: `gofmt -s -w pkg/tls/`
- [ ] All public methods have GoDoc

#### Size Compliance
- [ ] Line count measured with: `$PROJECT_ROOT/tools/line-counter.sh`
- [ ] Within estimate: 350 lines ±15% = 298-403 lines
- [ ] Tests NOT counted (per R007)

**Final Verification Commands**:
```bash
# 1. Run tests
go test ./pkg/tls -v -cover

# 2. Check coverage
go test ./pkg/tls -coverprofile=coverage.out
go tool cover -func=coverage.out

# 3. Lint
go vet ./pkg/tls
golangci-lint run ./pkg/tls

# 4. Format
gofmt -s -w pkg/tls/

# 5. Measure size
cd /home/vscode/workspaces/idpbuilder-oci-push-planning/efforts/phase1/wave2/effort-4-tls
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
$PROJECT_ROOT/tools/line-counter.sh
```

---

## Size Management

### Estimated Lines

**Implementation**: 350 lines (config.go)
**Tests**: ~200 lines (config_test.go, NOT counted per R007)
**Total Counted**: 350 lines

### Measurement Tool

**MANDATORY**: Use only `${PROJECT_ROOT}/tools/line-counter.sh`

**Process**:
1. Find project root (where orchestrator-state.json is):
   ```bash
   PROJECT_ROOT=$(pwd)
   while [ "$PROJECT_ROOT" != "/" ]; do
       [ -f "$PROJECT_ROOT/orchestrator-state.json" ] && break
       PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
   done
   ```

2. Run line counter (NO parameters needed):
   ```bash
   $PROJECT_ROOT/tools/line-counter.sh
   ```

3. Tool will auto-detect base branch and count ONLY implementation code

**Check Frequency**: After completing each major component

### Split Threshold

- **Warning**: 700 lines (not expected for this effort)
- **Hard Stop**: 800 lines (MUST split if reached)
- **Current Estimate**: 350 lines (WELL under limit)

**Expected**: No split required - effort is well-scoped

---

## Test Requirements

### Unit Test Coverage

**Minimum Coverage**: 90% (security-critical per Wave 2 Test Plan)

**Coverage Calculation**:
```bash
go test ./pkg/tls -coverprofile=coverage.out
go tool cover -func=coverage.out

# Must show: coverage: 90.0% or higher
```

### Test Categories

**Total Tests**: 8+ test cases

**Breakdown**:
- Constructor tests: 2 tests
- GetTLSConfig tests: 2 tests
- IsInsecure tests: 2 tests
- Integration tests: 1-2 tests

### Test Execution

**Run Tests**:
```bash
# Basic test run
go test ./pkg/tls -v

# With coverage
go test ./pkg/tls -v -cover

# Coverage report
go test ./pkg/tls -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

**Expected Output**:
```
=== RUN   TestNewConfigProvider_SecureMode
--- PASS: TestNewConfigProvider_SecureMode (0.00s)
=== RUN   TestNewConfigProvider_InsecureMode
--- PASS: TestNewConfigProvider_InsecureMode (0.00s)
=== RUN   TestGetTLSConfig_SecureMode
--- PASS: TestGetTLSConfig_SecureMode (0.00s)
=== RUN   TestGetTLSConfig_InsecureMode
--- PASS: TestGetTLSConfig_InsecureMode (0.00s)
=== RUN   TestIsInsecure_Secure
--- PASS: TestIsInsecure_Secure (0.00s)
=== RUN   TestIsInsecure_Insecure
--- PASS: TestIsInsecure_Insecure (0.00s)
=== RUN   TestTLSConfig_UsableWithHTTPClient
--- PASS: TestTLSConfig_UsableWithHTTPClient (0.00s)
PASS
coverage: 91.2% of statements in ./pkg/tls
ok      github.com/cnoe-io/idpbuilder/pkg/tls   0.015s
```

---

## Pattern Compliance

### Wave 1 Interface Usage

**MUST implement** (from `pkg/tls/interface.go`):
```go
type ConfigProvider interface {
    GetTLSConfig() *tls.Config
    IsInsecure() bool
}
```

**MUST use** (from `pkg/tls/types.go`):
```go
type Config struct {
    InsecureSkipVerify bool
}
```

### Standard Library Patterns

**Imports Required**:
- `crypto/tls` - TLS configuration
- `crypto/x509` - Certificate pool management

**No External Dependencies**: This package uses standard library only

### Error Handling

**No Wave 1 Errors Used**: TLS package does not return errors from Wave 1

**Internal Error Handling**:
- System cert pool loading may fail → fallback to empty pool
- No errors returned to caller (graceful degradation)

### Security Patterns

**Secure Mode (Default)**:
- `InsecureSkipVerify = false`
- System certificate pool loaded
- Full certificate verification enabled

**Insecure Mode (Explicit Flag)**:
- `InsecureSkipVerify = true`
- No certificate verification
- Use only for development/testing

**Documentation Must Warn**:
- Insecure mode disables security
- Only use with explicit user consent
- Never default to insecure mode

---

## Dependencies

### Upstream Dependencies (COMPLETED)

**Wave 1 Effort 3**: TLS Interface Definition
- **Location**: `pkg/tls/interface.go`
- **Status**: ✅ COMPLETED
- **Provides**:
  - `ConfigProvider` interface
  - `Config` struct
- **Integration Branch**: `idpbuilder-oci-push/phase1/wave2/integration`

### Downstream Dependencies (Parallel)

**Effort 1.2.2**: Registry Client Implementation
- **Status**: Running in parallel
- **Usage**: Registry client will USE this TLS implementation via interface
- **Integration**: After both efforts complete

**Wave 3**: CLI Implementation
- **Status**: Future wave
- **Usage**: CLI will create TLS providers with --insecure flag

### External Library Dependencies

**Standard Library Only**:
- `crypto/tls` (Go 1.21+)
- `crypto/x509` (Go 1.21+)
- `net/http` (for integration tests)

**Test Dependencies**:
- `github.com/stretchr/testify` v1.10.0 (already in Wave 1)

**No go.mod Changes Required**: All dependencies already present

---

## Integration Points

### Usage by Registry Client

**Registry client will use TLS provider as follows** (from Wave 2 Architecture):

```go
// In pkg/registry/client.go

func NewClient(authProvider AuthProvider, tlsConfig TLSConfigProvider) (Client, error) {
    // Validate TLS provider
    if tlsConfig == nil {
        return nil, &ValidationError{
            Field:   "tlsConfig",
            Message: "TLS config provider cannot be nil",
        }
    }

    // Get TLS config
    httpClient := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: tlsConfig.GetTLSConfig(),  // ← Uses our implementation
        },
    }

    return &registryClient{
        authProvider: authProvider,
        tlsConfig:    tlsConfig,
        httpClient:   httpClient,
    }, nil
}
```

### Complete Push Workflow (from Wave 2 Architecture)

**How TLS fits into the full workflow**:

```go
// Step 1: Create TLS configuration provider
tlsProvider := tls.NewConfigProvider(insecure)  // ← This effort
if tlsProvider.IsInsecure() {
    fmt.Println("WARNING: TLS certificate verification disabled")
}

// Step 2: Create registry client with TLS provider
registryClient, err := registry.NewClient(authProvider, tlsProvider)
if err != nil {
    return fmt.Errorf("failed to create registry client: %w", err)
}

// Step 3: Push image (TLS config used for HTTPS connections)
err = registryClient.Push(ctx, image, targetRef, progressCallback)
```

---

## Acceptance Criteria

### Functionality Checklist

- [ ] All 2 interface methods implemented correctly
  - [ ] `GetTLSConfig()` returns valid `*tls.Config`
  - [ ] `IsInsecure()` returns correct mode flag
- [ ] Secure mode loads system certificate pool
- [ ] Insecure mode disables verification correctly
- [ ] TLS config compatible with HTTP transport
- [ ] No security warnings for secure mode

### Quality Checklist

- [ ] All tests passing (100% pass rate)
- [ ] Code coverage ≥90% (security-critical)
- [ ] No linting errors (go vet, golangci-lint)
- [ ] Documentation complete (all public methods have GoDoc)
- [ ] Code formatted with gofmt

### Compliance Checklist

- [ ] All files created/modified as specified
- [ ] Line count within estimate (350 lines ±15% = 298-403 lines)
- [ ] Uses Wave 1 interface correctly
- [ ] No unauthorized dependencies added
- [ ] Implements security best practices

### Integration Checklist

- [ ] Compatible with Wave 1 `ConfigProvider` interface
- [ ] Can be used by registry client (interface dependency)
- [ ] HTTP client integration working
- [ ] Secure/insecure modes function correctly

### Review Criteria

**Code Reviewer will verify**:
- Interface implementation complete
- Test coverage meets 90% target
- Security patterns correct (default secure, explicit insecure)
- Documentation adequate
- Size within limit

**Automated Checks**:
- Tests pass: `go test ./pkg/tls -v`
- Coverage: `go test ./pkg/tls -cover` (≥90%)
- Linting: `go vet ./pkg/tls` (no errors)
- Formatting: `gofmt -l pkg/tls/` (no output)

---

## Risk Mitigation

### Identified Risks

#### Risk 1: System Certificate Pool Unavailable

**Severity**: Low
**Impact**: TLS connections may fail in secure mode

**Mitigation**:
- Fallback to empty cert pool if `x509.SystemCertPool()` fails
- Document behavior in GoDoc
- Test both success and fallback paths

**Code Pattern**:
```go
certPool, err := x509.SystemCertPool()
if err != nil {
    // Fallback to empty pool
    certPool = x509.NewCertPool()
}
```

#### Risk 2: Insecure Mode Used Accidentally

**Severity**: High (security)
**Impact**: Connections without certificate verification

**Mitigation**:
- Default to secure mode (insecure requires explicit flag)
- IsInsecure() method allows runtime checks
- Documentation warns about insecure mode
- CLI will require explicit --insecure/-k flag

**Documentation Warning**:
```go
// WARNING: Insecure mode disables TLS certificate verification.
// Only use for development/testing with self-signed certificates.
// Never use in production without understanding security implications.
```

#### Risk 3: TLS Config Incompatibility

**Severity**: Medium
**Impact**: May not work with all HTTP clients

**Mitigation**:
- Follow standard `crypto/tls` patterns
- Test integration with `http.Client`
- Use standard `tls.Config` type (widely compatible)

### Testing Edge Cases

**Secure Mode**:
- System certs available (normal case)
- System certs unavailable (fallback)
- Valid HTTPS connections
- Invalid certificates (should fail)

**Insecure Mode**:
- Self-signed certificates (should accept)
- Expired certificates (should accept)
- Invalid certificates (should accept)
- HTTP connections (should work)

---

## Timeline Estimate

### Implementation Phases

**Phase 1: Setup** (30 minutes)
- Create package structure
- Verify Wave 1 interfaces accessible
- Set up test file structure

**Phase 2: Core Implementation** (2 hours)
- Implement tlsConfigProvider struct
- Implement NewConfigProvider constructor
- Implement GetTLSConfig method
- Implement IsInsecure method

**Phase 3: Testing** (2 hours)
- Write constructor tests
- Write GetTLSConfig tests
- Write IsInsecure tests
- Write integration tests
- Achieve 90% coverage

**Phase 4: Documentation** (1 hour)
- Add GoDoc comments
- Write usage examples
- Document security considerations

**Phase 5: Verification** (30 minutes)
- Run all tests
- Check coverage
- Run linters
- Measure size
- Final validation

**Total Estimated Time**: ~6 hours

---

## Success Metrics

### Quantitative Metrics

- **Line Count**: 298-403 lines (350 ±15%)
- **Test Coverage**: ≥90%
- **Test Pass Rate**: 100%
- **Linting Errors**: 0
- **Documentation**: 100% of public methods

### Qualitative Metrics

- **Code Quality**: Clean, readable, maintainable
- **Security**: Secure by default, explicit insecure mode
- **Usability**: Simple API, clear documentation
- **Integration**: Works seamlessly with registry client

### Review Outcomes

**Expected Outcome**: ACCEPTED
- All acceptance criteria met
- Tests passing
- Coverage target achieved
- Security patterns correct
- Documentation complete

**Unlikely Outcomes**:
- NEEDS_FIXES: Minor issues found (e.g., documentation gaps)
- NEEDS_SPLIT: Not expected (well under 800 line limit)

---

## Reference Documentation

### Wave 2 Planning Documents

**MUST READ Before Implementation**:
1. **Wave 2 Implementation Plan**: `wave-plans/WAVE-2-IMPLEMENTATION.md`
   - Section: Effort 1.2.4 (lines 571-740)
   - Contains: Exact specifications, files, dependencies
2. **Wave 2 Architecture**: `wave-plans/WAVE-2-ARCHITECTURE.md`
   - Section: Effort 1.2.4 (lines 788-889)
   - Contains: Complete implementation code with comments
3. **Wave 2 Test Plan**: `wave-plans/WAVE-2-TEST-PLAN.md`
   - Section: Package tls (lines 1341-1547)
   - Contains: All test cases with specifications

### Wave 1 Interface Files

**Reference During Implementation**:
- `pkg/tls/interface.go` - ConfigProvider interface
- `pkg/tls/types.go` - Config struct definition

### External Documentation

**Standard Library References**:
- `crypto/tls` package: https://pkg.go.dev/crypto/tls
- `crypto/x509` package: https://pkg.go.dev/crypto/x509
- TLS best practices: https://go.dev/blog/tls

---

## Commit Strategy

### Commit Frequency

**Commit after each major step**:
1. After package structure setup
2. After core implementation complete
3. After tests written and passing
4. After documentation added
5. After final verification

### Commit Message Format

**Pattern**:
```
feat(tls): [description]

[Details about changes]

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

**Examples**:
```
feat(tls): implement TLS configuration provider

- Add tlsConfigProvider struct
- Implement NewConfigProvider constructor
- Implement GetTLSConfig method
- Implement IsInsecure method
- Support secure and insecure modes
- Load system certificate pool

Implements Wave 2 Effort 1.2.4 specification.
Lines: 350 (within 298-403 limit)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

```
test(tls): add comprehensive unit tests

- Add constructor tests (secure/insecure)
- Add GetTLSConfig tests
- Add IsInsecure tests
- Add HTTP client integration test
- Coverage: 91.2% (exceeds 90% target)

All tests passing (8/8).

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

### Push Strategy

**Push frequently**: After each commit
- Ensures work is backed up
- Allows orchestrator to track progress
- Enables early feedback

---

## Additional Notes

### Parallelization Context

**This effort runs in parallel with**:
- Effort 1.2.1: Docker Client Implementation
- Effort 1.2.2: Registry Client Implementation
- Effort 1.2.3: Authentication Implementation

**No coordination required during implementation** because:
- All efforts use frozen Wave 1 interfaces
- No cross-effort dependencies
- Each implements different package
- Integration happens AFTER all complete

### Integration Timeline

**After this effort completes**:
1. Code Reviewer validates implementation
2. If approved, merge to integration branch
3. Wait for all 4 Wave 2 efforts to complete
4. Run Wave 2 integration tests
5. Architect reviews Wave 2 compliance
6. Merge integration branch to Phase 1

### Critical Success Factors

**For ACCEPTED review**:
1. ✅ Interface implementation complete and correct
2. ✅ Test coverage ≥90% achieved
3. ✅ Secure mode loads system certs properly
4. ✅ Insecure mode disables verification correctly
5. ✅ Documentation complete and clear
6. ✅ Size within limit (298-403 lines)
7. ✅ No linting errors
8. ✅ Security patterns followed

**Common Pitfalls to Avoid**:
- ❌ Not loading system cert pool in secure mode
- ❌ Defaulting to insecure mode (must be explicit)
- ❌ Returning errors from GetTLSConfig (should handle gracefully)
- ❌ Missing security warnings in documentation
- ❌ Inadequate test coverage (<90%)

---

## Document Status

**Status**: ✅ READY FOR IMPLEMENTATION
**Created**: 2025-10-29 21:33:16 UTC
**Planner**: Code Reviewer Agent (code-reviewer)
**Effort**: 1.2.4 - TLS Configuration Implementation
**Branch**: `idpbuilder-oci-push/phase1/wave2/effort-4-tls`

**Compliance Summary**:
- ✅ R356: Implementation plan creation protocol followed
- ✅ R383: Plan in .software-factory with timestamp
- ✅ R219: Dependencies analyzed (Wave 1 interfaces)
- ✅ R211: Parallelization metadata included
- ✅ Complete specifications from Wave 2 Architecture
- ✅ All test cases from Wave 2 Test Plan included
- ✅ Clear acceptance criteria defined

**Next Action**: SW Engineer implements based on this plan
- Read this plan thoroughly
- Read Wave 2 Architecture (pseudocode/examples)
- Read Wave 2 Test Plan (test specifications)
- Implement code to pass all tests
- Achieve ≥90% coverage
- Request Code Reviewer validation

---

**END OF IMPLEMENTATION PLAN**
