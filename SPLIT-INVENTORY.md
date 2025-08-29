# Split Inventory for certificate-validation

## Overview
The certificate-validation effort exceeded 800 lines (actual: 1350) after adding audit persistence infrastructure and requires splitting into manageable sub-efforts.

- **Original Size**: 1350 lines
- **Number of Splits**: 2
- **Date Split**: 2025-08-29
- **Split By**: Code Reviewer Agent

## Split Structure

| Split # | Name | Description | Est. Lines | Status |
|---------|------|-------------|------------|--------|
| 001 | core-validation | Core certificate chain validation logic and interfaces | 700 | Planned |
| 002 | audit-infrastructure | Audit logging and persistence infrastructure | 650 | Planned |

## Integration Strategy
The splits will be implemented sequentially:
1. **Split-001 (core-validation)**: Establishes the foundation with validation interfaces and core implementation
2. **Split-002 (audit-infrastructure)**: Adds audit logging capabilities on top of the core validation

Both splits will integrate cleanly as Split-002 extends Split-001 without modifying core logic.

## Files Distribution

### Split-001 (core-validation)
Files to implement:
- `pkg/certs/chain_validator.go` (20 lines)
- `pkg/certs/chain_validator_impl.go` (564 lines)
- `pkg/certs/errors.go` (21 lines)
- `pkg/certs/types_chain.go` (171 lines)
- `pkg/certs/wave1_interfaces.go` (78 lines)
- Parts of `pkg/certs/chain_validator_test.go` (basic tests only, ~200 lines)

**Total estimated**: ~700 lines (excluding full test coverage)

### Split-002 (audit-infrastructure)
Files to implement:
- `pkg/certs/audit/interface.go` (77 lines)
- `pkg/certs/audit/logger.go` (278 lines)
- `pkg/certs/audit/logger_test.go` (348 lines)
- `pkg/certs/audit/example_usage.go` (133 lines)
- Remaining `pkg/certs/chain_validator_test.go` tests (~308 lines)

**Total estimated**: ~650 lines

## Dependencies
- Split-001: No dependencies (foundational)
- Split-002: Depends on Split-001 (requires chain validator interfaces)
- Both splits must be completed sequentially

## Validation
Each split must:
- [ ] Stay under 800 lines (target <700 for safety)
- [ ] Pass all unit tests independently
- [ ] Integrate cleanly with previous splits
- [ ] Maintain backward compatibility with Wave 1 interfaces
- [ ] Include appropriate documentation

## Integration Points
- Split-001 provides `ChainValidator` interface
- Split-002 extends with `AuditLogger` interface
- Both use common types from `types_chain.go`
- Test coverage split across both efforts

## Risk Mitigation
- Each split is functionally complete and testable
- Clear interface boundaries prevent coupling
- Audit infrastructure is additive, not invasive
- Tests validate integration points