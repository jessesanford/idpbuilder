# Split Plan for Registry TLS Trust Integration

## Overview

**Effort**: E1.1.2 - Registry TLS Trust Integration
**Current Size**: 807 lines (7 lines over 800-line limit)
**Split Strategy**: 2 sequential splits
**Planner**: Code Reviewer Agent
**Date**: 2025-08-31

## Problem Statement

The implementation exceeded the 800-line hard limit by 7 lines. The code needs to be split into smaller, manageable pieces that each stay well under the limit while maintaining logical cohesion.

## Split Summary

| Split | Description | Files | Size | Dependencies |
|-------|-------------|-------|------|--------------|
| 001 | Core Trust Store Management | trust.go, partial tests | ~377 lines | None |
| 002 | Transport & Utilities | transport.go, trust_store.go, remaining tests | ~551 lines | Split 001 |

## Execution Order

**CRITICAL**: These splits must be executed SEQUENTIALLY, not in parallel:

1. **Split 001** - Implement core trust store first (provides interfaces)
2. **Split 002** - Then implement transport and utilities (uses Split 001's interfaces)

## Key Benefits of This Split

1. **Logical Separation**: Core trust management is separated from transport configuration
2. **Size Compliance**: Each split is well under 800 lines (377 and 551)
3. **Clean Dependencies**: Split 002 depends on Split 001's interfaces only
4. **Testability**: Each split can be tested independently
5. **Future Buffer**: Both splits have room for growth without exceeding limits

## Files Created

- `SPLIT-INVENTORY.md` - Complete inventory and deduplication matrix
- `SPLIT-PLAN-001.md` - Detailed plan for Split 001 (Core Trust Store)
- `SPLIT-PLAN-002.md` - Detailed plan for Split 002 (Transport & Utilities)

## Implementation Guidelines for SW Engineers

### For Split 001:
- Focus on core trust store functionality
- Ensure interface is well-defined for Split 002 to use
- Include basic tests for trust store operations
- Target: ~377 lines

### For Split 002:
- Import and use Split 001's TrustStoreManager interface
- Implement transport configuration for go-containerregistry
- Add utility functions for certificate operations
- Complete test coverage
- Target: ~551 lines

## Verification Steps

1. Each split must compile independently
2. Each split must pass its tests
3. Each split must be under 800 lines (use line-counter.sh)
4. Combined functionality must match original requirements
5. No code duplication between splits

## Risk Mitigation

- Both splits are sized conservatively (< 600 lines) to allow for minor additions
- Clear interface boundaries prevent coupling issues
- Sequential execution ensures dependencies are met
- Comprehensive test coverage in both splits

## Next Steps

1. Orchestrator assigns Split 001 to SW Engineer
2. After Split 001 completion and review, assign Split 002
3. After both splits complete, merge sequentially to parent branch
4. Verify integrated functionality meets all requirements

## Success Metrics

- [x] Both splits under 800 lines
- [x] No code duplication
- [x] Clear logical separation
- [x] Maintainable structure
- [x] All original functionality preserved
