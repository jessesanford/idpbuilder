# Certificate Manager Abstraction Implementation Plan

**Created**: 2025-09-28T14:25:29Z
**Phase**: 1
**Wave**: 2
**Effort**: P1W2-E3: Certificate Manager Abstraction
**Branch**: `phase1/wave2/certificate-manager`
**Directory**: `/home/vscode/workspaces/idpbuilder-gitea-push/efforts/phase1/wave2/certificate-manager`
**Agent**: code-reviewer
**State**: EFFORT_PLAN_CREATION

## Executive Summary

This effort creates a clean abstraction layer for TLS certificate handling and validation, wrapping the standard library's crypto/x509 package behind testable interfaces. The abstraction enables mockability for testing and provides a foundation for secure registry connections in Wave 3.

## Pre-Planning Research Results (R374 MANDATORY)

### Existing Interfaces Found
| Interface | Location | Signature | Must Implement |
|-----------|----------|-----------|----------------|
| None specific to certificates | N/A | N/A | N/A |

### Existing Implementations to Reuse
| Component | Location | Purpose | How to Use |
|-----------|----------|---------|------------|
| TLSConfig struct | efforts/phase1/wave1/integration-workspace/pkg/config/auth.go | TLS configuration | Reference for config structure |
| Insecure flag | efforts/phase1/wave1/integration-workspace/pkg/config/registry.go | Skip TLS verification | Support this flag in manager |

### APIs Already Defined
| API | Method | Signature | Notes |
|-----|--------|-----------|-------|
| Registry config | Insecure bool | `config.Insecure` | Must support this for dev environments |

### FORBIDDEN DUPLICATIONS (R373)
- DO NOT create alternative TLSConfig structures
- DO NOT reimplement crypto/x509 functionality directly
- DO NOT create competing registry config types

### REQUIRED INTEGRATIONS (R373)
- MUST support Wave 1's Insecure flag from registry config
- MUST integrate cleanly with future registry client needs
- MUST provide mockable interfaces for testing

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

## EXPLICIT SCOPE (R311 MANDATORY)

### IMPLEMENT EXACTLY:
- **Interface**: `Manager` interface with 3 methods (~30 lines)
  - `LoadSystemCerts(ctx context.Context) (*x509.CertPool, error)`
  - `ValidateCertificate(ctx context.Context, cert *x509.Certificate) error`
  - `CreateTLSConfig(ctx context.Context, insecure bool) (*tls.Config, error)`
- **Implementation**: `X509Manager` struct implementing Manager (~150 lines)
  - Full implementation using crypto/x509
  - Context support for cancellation
  - Proper error handling
- **Type**: `Validator` interface with 1 method (~10 lines)
  - `Validate(cert *x509.Certificate, opts x509.VerifyOptions) error`
- **Implementation**: `DefaultValidator` struct (~70 lines)
  - Standard x509 validation logic
- **Type**: `Store` interface with 2 methods (~10 lines)
  - `GetPool(ctx context.Context) (*x509.CertPool, error)`
  - `AddCert(ctx context.Context, cert *x509.Certificate) error`
- **Implementation**: `MemoryStore` struct (~30 lines)
  - In-memory certificate storage
- **Error Type**: `CertError` struct (~20 lines)
  - Error code enumeration
  - Context preservation
- **Tests**: Unit tests for all implementations (~100 lines)
  - Test certificate loading
  - Test validation logic
  - Test TLS config creation

**TOTAL**: ~420 lines (well under 800)

### DO NOT IMPLEMENT:
- ❌ Certificate generation (not needed)
- ❌ Complex PKI operations
- ❌ Certificate rotation logic (future effort)
- ❌ File-based certificate storage (future effort)
- ❌ Custom CA management (future effort)
- ❌ Certificate parsing utilities
- ❌ HTTP client with certificate handling
- ❌ Registry-specific certificate logic (Wave 3)
- ❌ Certificate chain building
- ❌ OCSP checking

## Size Limit Clarification (R359):
- The 800-line limit applies to NEW CODE YOU ADD
- Repository will grow by ~420 lines (EXPECTED)
- NEVER delete existing code to meet size limits
- Current codebase size: minimal (new effort directory)
- Expected total after: ~420 lines

## Configuration Requirements (R355 Mandatory)

### WRONG - Will fail review:
```go
// ❌ VIOLATION - Hardcoded certificate path
certPath := "/etc/ssl/certs/ca.crt"

// ❌ VIOLATION - Stub implementation
func LoadSystemCerts(ctx context.Context) (*x509.CertPool, error) {
    // TODO: implement later
    return nil, nil
}

// ❌ VIOLATION - Static configuration
insecure := true // Always skip verification
```

### CORRECT - Production ready:
```go
// ✅ From system certificate store
func (m *X509Manager) LoadSystemCerts(ctx context.Context) (*x509.CertPool, error) {
    pool, err := x509.SystemCertPool()
    if err != nil {
        return nil, fmt.Errorf("failed to load system certs: %w", err)
    }
    return pool, nil
}

// ✅ Full implementation required
func (m *X509Manager) CreateTLSConfig(ctx context.Context, insecure bool) (*tls.Config, error) {
    config := &tls.Config{
        InsecureSkipVerify: insecure,
    }
    if !insecure {
        pool, err := m.LoadSystemCerts(ctx)
        if err != nil {
            return nil, fmt.Errorf("failed to load certs: %w", err)
        }
        config.RootCAs = pool
    }
    return config, nil
}

// ✅ Configurable from external input
func NewManager(insecure bool) *X509Manager {
    return &X509Manager{
        insecure: insecure,
    }
}
```

## File Structure

```
efforts/phase1/wave2/certificate-manager/
├── .software-factory/
│   └── IMPLEMENTATION-PLAN-20250928-142529.md  (this file)
└── pkg/
    └── certs/
        ├── manager.go          # Manager interface (~30 lines)
        ├── x509_adapter.go     # X509Manager implementation (~150 lines)
        ├── validator.go        # Validator interface and DefaultValidator (~80 lines)
        ├── store.go            # Store interface and MemoryStore (~40 lines)
        ├── errors.go           # CertError type (~20 lines)
        ├── manager_test.go     # Manager tests (~50 lines)
        ├── validator_test.go   # Validator tests (~30 lines)
        └── store_test.go       # Store tests (~20 lines)
```

## Implementation Steps

### Step 1: Create Manager Interface (30 lines)
```go
// pkg/certs/manager.go
package certs

import (
    "context"
    "crypto/tls"
    "crypto/x509"
)

// Manager handles TLS certificate operations
type Manager interface {
    // LoadSystemCerts loads the system certificate pool
    LoadSystemCerts(ctx context.Context) (*x509.CertPool, error)

    // ValidateCertificate validates a single certificate
    ValidateCertificate(ctx context.Context, cert *x509.Certificate) error

    // CreateTLSConfig creates a TLS configuration
    CreateTLSConfig(ctx context.Context, insecure bool) (*tls.Config, error)
}
```

### Step 2: Implement X509Manager (150 lines)
- Full working implementation wrapping crypto/x509
- Context support for all operations
- Proper error handling with wrapped errors
- Support for insecure mode (dev environments)

### Step 3: Create Validator Interface (80 lines total)
- Define Validator interface
- Implement DefaultValidator with standard x509 validation
- Include verify options configuration

### Step 4: Create Store Interface (40 lines total)
- Define Store interface for certificate storage
- Implement MemoryStore for in-memory certificate pool
- Thread-safe operations with mutex

### Step 5: Define Error Types (20 lines)
- CertError struct with error codes
- Error code constants (InvalidCert, Expired, etc.)
- Error wrapping for context preservation

### Step 6: Write Tests (100 lines)
- Unit tests for Manager implementation
- Tests for Validator logic
- Tests for Store operations
- Mock certificate creation for testing

## Atomic PR Design (R220)

```yaml
effort_atomic_pr_design:
  pr_summary: "feat: add certificate manager abstraction for TLS handling"
  can_merge_to_main_alone: true  # MUST be true

  r355_production_ready_checklist:
    no_hardcoded_values: true
    all_config_from_env: true
    no_stub_implementations: true
    no_todo_markers: true
    all_functions_complete: true

  configuration_approach:
    - name: "System cert pool"
      wrong: 'pool := loadFromFile("/etc/ssl/certs")'
      correct: 'pool, err := x509.SystemCertPool()'
    - name: "Insecure mode"
      wrong: 'insecure := true // hardcoded'
      correct: 'insecure := config.Insecure // from external config'

  feature_flags_needed: []  # No feature flags needed - complete implementation

  interface_implementations:
    - interface: "Manager"
      implementation: "X509Manager"
      production_ready: true
      notes: "Fully functional certificate manager"
    - interface: "Validator"
      implementation: "DefaultValidator"
      production_ready: true
      notes: "Complete x509 validation"
    - interface: "Store"
      implementation: "MemoryStore"
      production_ready: true
      notes: "Thread-safe in-memory storage"

  pr_verification:
    tests_pass_alone: true
    build_remains_working: true
    flags_tested_both_ways: false  # No flags
    no_external_dependencies: false  # Uses crypto/x509
    backward_compatible: true

  example_pr_structure:
    files_added:
      - "pkg/certs/manager.go"
      - "pkg/certs/x509_adapter.go"
      - "pkg/certs/validator.go"
      - "pkg/certs/store.go"
      - "pkg/certs/errors.go"
      - "pkg/certs/manager_test.go"
      - "pkg/certs/validator_test.go"
      - "pkg/certs/store_test.go"
    tests_included:
      - "Unit tests for certificate loading"
      - "Unit tests for validation"
      - "Unit tests for TLS config creation"
      - "Tests with insecure mode on/off"
    documentation:
      - "Interface documentation"
      - "Usage examples in comments"
```

## Dependencies

### Wave 1 Dependencies
- `Insecure` flag from registry configuration (for skip TLS verification)

### External Dependencies
- Standard library: `crypto/x509`, `crypto/tls`
- Standard library: `context`, `fmt`, `errors`, `sync`

### No Dependencies On
- Other Wave 2 efforts (fully parallel)
- Registry client implementation
- Build operations

## Testing Strategy

### Unit Tests (60% minimum coverage)
1. **Manager Tests** (manager_test.go):
   - Test system cert pool loading
   - Test TLS config creation with insecure=true
   - Test TLS config creation with insecure=false
   - Test certificate validation

2. **Validator Tests** (validator_test.go):
   - Test valid certificate acceptance
   - Test expired certificate rejection
   - Test invalid signature detection

3. **Store Tests** (store_test.go):
   - Test certificate addition
   - Test pool retrieval
   - Test concurrent access safety

### Mock Implementations
- Create test certificates using crypto/x509 test utilities
- No external mocks needed (self-contained)

## Success Criteria

- ✅ All interfaces fully implemented (no stubs)
- ✅ Production-ready code from first commit
- ✅ 60% test coverage minimum
- ✅ All tests passing
- ✅ Under 800 lines total
- ✅ Clean abstraction over crypto/x509
- ✅ Support for insecure mode (development)
- ✅ Context support for all operations
- ✅ Proper error handling with context

## Risk Mitigation

### Technical Risks
1. **Certificate format variations**:
   - Use standard x509 parsing
   - Handle errors gracefully
   - Document supported formats

2. **System certificate pool issues**:
   - Fallback to empty pool if system pool fails
   - Clear error messages
   - Document platform differences

### Process Risks
1. **Size overrun**:
   - Current estimate: ~420 lines
   - Buffer: 380 lines remaining
   - Monitor with line-counter.sh

## Integration Points

### Wave 3 Usage
- Registry client will use Manager.CreateTLSConfig()
- Push operations will use certificate validation
- All registry connections will go through this abstraction

### Future Enhancements (NOT in this effort)
- File-based certificate storage
- Certificate rotation
- Custom CA management
- OCSP checking
- Certificate chain validation

## Out of Scope (R311 Enforcement)

This effort does NOT include:
- Registry-specific certificate handling
- HTTP transport implementation
- Authentication mechanisms
- Certificate generation
- Complex PKI operations
- Integration with registry client
- Command-line certificate management
- Certificate caching beyond memory

## Validation Checklist

Before starting implementation:
- [ ] Plan specifies EXACTLY 3 interfaces
- [ ] Plan specifies EXACTLY 3 implementations
- [ ] Total line count estimate: ~420 lines
- [ ] No stub implementations allowed
- [ ] Production-ready from first commit
- [ ] Clear boundaries defined
- [ ] Dependencies identified
- [ ] Test strategy defined

## 📋 PLANNING FILE CREATED

**Type**: effort_plan
**Path**: /home/vscode/workspaces/idpbuilder-gitea-push/efforts/phase1/wave2/certificate-manager/.software-factory/IMPLEMENTATION-PLAN-20250928-142529.md
**Effort**: certificate-manager
**Phase**: 1
**Wave**: 2
**Target Branch**: phase1/wave2/certificate-manager
**Created At**: 2025-09-28T14:25:29Z

ORCHESTRATOR: Please update planning_files.effort_plans["certificate-manager"] in state file per R340

---

**Plan Created By**: @agent-code-reviewer
**Plan Status**: COMPLETE
**Ready for**: SW-Engineer Implementation