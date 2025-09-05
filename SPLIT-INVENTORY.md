# Complete Split Plan for registry-tls-trust-integration

**Sole Planner**: Code Reviewer Agent
**Full Path**: phase1/wave1/registry-tls-trust-integration
**Parent Branch**: idpbuidler-oci-go-cr/phase1/wave1/registry-tls-trust-integration
**Total Size**: 807 lines
**Splits Required**: 2
**Created**: 2025-08-31 14:54:00 UTC

⚠️ **SPLIT INTEGRITY NOTICE** ⚠️
ALL splits below belong to THIS effort ONLY: phase1/wave1/registry-tls-trust-integration
NO splits should reference efforts outside this path!

## Split Boundaries (NO OVERLAPS)

| Split | Start Line | End Line | Estimated Size | Files | Status |
|-------|------------|----------|----------------|-------|--------|
| 001   | N/A        | N/A      | 511 lines     | trust.go (317), partial tests (194) | ✅ Complete |
| 002   | N/A        | N/A      | 468 lines     | transport.go (251), trust_store.go (217) | ✅ Complete |

## File Distribution Matrix

| File | Split 001 | Split 002 | Lines |
|------|-----------|-----------|-------|
| pkg/certs/trust.go | ✅ | ❌ | 317 |
| pkg/certs/transport.go | ❌ | ✅ | 270 |
| pkg/certs/trust_store.go | ❌ | ✅ | 217 |
| pkg/certs/trust_test.go | Partial (60 lines) | Partial (64 lines) | 124 |

## Logical Separation

### Split 001: Core Trust Store Management
- **Primary Focus**: TrustStoreManager interface and implementation
- **Components**:
  - TrustStoreManager interface definition
  - trustStoreManager struct implementation
  - Certificate loading and storage
  - Insecure registry management
  - Basic unit tests for trust store operations
- **Dependencies**: None (foundational split)

### Split 002: Transport Configuration and Utilities
- **Primary Focus**: HTTP transport configuration and utility functions
- **Components**:
  - Transport configuration for go-containerregistry
  - Certificate utility functions
  - PEM encoding/decoding utilities
  - Certificate validation helpers
  - Tests for transport and utilities
- **Dependencies**: Split 001 (uses TrustStoreManager interface)

## Verification Checklist
- [x] No file appears in multiple splits (except tests which are logically divided)
- [x] All files from original effort covered
- [x] Each split can compile independently with proper interfaces
- [x] Dependencies properly ordered (001 before 002)
- [x] Each split <800 lines (target <450 lines each)
- [x] Logical cohesion maintained

## Implementation Strategy

1. **Split 001 Implementation**:
   - Implement core trust store manager
   - Provide interface for transport layer
   - Include foundational tests

2. **Split 002 Implementation**:
   - Import trust store interface from Split 001
   - Add transport configuration layer
   - Add utility functions
   - Complete test coverage

## Notes
- Tests are split logically: Split 001 tests the trust store core, Split 002 tests transport and utilities
- Each split remains under 450 lines to provide buffer for future additions
- Sequential execution required: Split 002 depends on Split 001's interface