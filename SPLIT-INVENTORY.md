# Split Inventory for fallback-strategies

## Overview
The fallback-strategies effort exceeded 800 lines (actual: 2526) after adding recovery mechanisms and Wave 1 integration, requiring splitting into manageable sub-efforts.

- **Original Size**: 2526 lines
- **Number of Splits**: 4
- **Date Split**: 2025-08-29
- **Split By**: Code Reviewer Agent

## Split Structure

| Split # | Name | Description | Est. Lines | Status |
|---------|------|-------------|------------|--------|
| 001 | core-fallback | Core fallback mechanism and basic strategies | 650 | Planned |
| 002 | insecure-mode | Insecure mode handling and bypass logic | 650 | Planned |
| 003 | recovery-strategies | Recovery mechanisms and error handling | 650 | Planned |
| 004 | tests-integration | Complete test coverage and Wave 1 integration | 576 | Planned |

## Integration Strategy
The splits will be implemented sequentially:
1. **Split-001 (core-fallback)**: Establishes the foundation with fallback interfaces and basic implementation
2. **Split-002 (insecure-mode)**: Adds insecure mode capabilities for development/testing scenarios
3. **Split-003 (recovery-strategies)**: Implements comprehensive recovery mechanisms for certificate failures
4. **Split-004 (tests-integration)**: Completes test coverage and ensures Wave 1 integration works

All splits build upon each other sequentially, with each adding specific capabilities.

## Files Distribution

### Split-001 (core-fallback)
Files to implement:
- `pkg/certs/fallback.go` (520 lines)
- Basic tests in `pkg/certs/fallback_test.go` (~130 lines)

**Total estimated**: 650 lines

### Split-002 (insecure-mode)
Files to implement:
- `pkg/certs/insecure.go` (179 lines)
- `pkg/certs/insecure_test.go` (180 lines)
- Complete `pkg/certs/fallback_test.go` (remaining ~110 lines)
- Integration helpers (~181 lines)

**Total estimated**: 650 lines

### Split-003 (recovery-strategies)
Files to implement:
- Start `pkg/certs/recovery.go` (first ~400 lines)
- Start `pkg/certs/recovery_test.go` (first ~250 lines)

**Total estimated**: 650 lines

### Split-004 (tests-integration)
Files to implement:
- Complete `pkg/certs/recovery.go` (remaining ~323 lines)
- Complete `pkg/certs/recovery_test.go` (remaining ~198 lines)
- Wave 1 integration tests (~55 lines)

**Total estimated**: 576 lines

## Dependencies
- Split-001: No dependencies (foundational)
- Split-002: Depends on Split-001 (requires fallback interfaces)
- Split-003: Depends on Split-001 and Split-002 (uses both fallback and insecure modes)
- Split-004: Depends on all previous splits (completes implementation)
- All splits must be completed sequentially

## Validation
Each split must:
- [ ] Stay under 800 lines (target <700 for safety)
- [ ] Pass all unit tests independently
- [ ] Integrate cleanly with previous splits
- [ ] Maintain backward compatibility with Wave 1
- [ ] Include appropriate documentation
- [ ] Handle error cases properly

## Integration Points
- Split-001 provides `FallbackStrategy` interface
- Split-002 extends with `InsecureMode` capabilities
- Split-003 adds `RecoveryManager` interface
- Split-004 completes integration and test coverage
- All use Wave 1 certificate interfaces

## Risk Mitigation
- Each split is functionally coherent
- Clear interface boundaries between splits
- Progressive enhancement pattern
- Tests distributed across splits for validation
- Recovery mechanisms tested incrementally

## Special Considerations
- Recovery.go is the largest file (723 lines) and must be split
- Test files are substantial and distributed across splits
- Insecure mode requires careful security documentation
- Wave 1 integration must be preserved throughout