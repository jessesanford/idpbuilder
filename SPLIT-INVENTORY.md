# Complete Split Plan for cert-extraction

**Sole Planner**: Code Reviewer Agent
**Full Path**: phase1/wave1/cert-extraction
**Parent Branch**: idpbuilder-oci-mvp/phase1/wave1/cert-extraction
**Total Size**: 836 lines
**Splits Required**: 2
**Created**: 2025-08-29 00:10:00

⚠️ **SPLIT INTEGRITY NOTICE** ⚠️
ALL splits below belong to THIS effort ONLY: phase1/wave1/cert-extraction
NO splits should reference efforts outside this path!

## Split Boundaries (NO OVERLAPS)
| Split | Start Line | End Line | Size | Files | Status |
|-------|------------|----------|------|-------|--------|
| 001   | 1          | 493      | ~493 | types.go, errors.go, extractor.go | Planned |
| 002   | 494        | 836      | ~343 | validator.go, all tests | Planned |

## File Distribution Matrix
| File | Lines | Split 001 | Split 002 |
|------|-------|-----------|-----------|
| types.go | 68 | ✅ | ❌ |
| errors.go | 104 | ✅ | ❌ |
| extractor.go | 321 | ✅ | ❌ |
| validator.go | 217 | ❌ | ✅ |
| errors_test.go | 204 | ❌ | ✅ (reduced) |
| extractor_test.go | 391 | ❌ | ✅ (reduced) |
| validator_test.go | 321 | ❌ | ✅ (reduced) |

## Deduplication Matrix
| Component | Split 001 | Split 002 | Notes |
|-----------|-----------|-----------|-------|
| Core Types/Interfaces | ✅ | ❌ | Foundation for all functionality |
| Error Handling | ✅ | ❌ | Required by both but in Split 001 |
| Certificate Extraction | ✅ | ❌ | Main extraction logic |
| Certificate Validation | ❌ | ✅ | Depends on types from Split 001 |
| Unit Tests | ❌ | ✅ | All tests in Split 002 |

## Implementation Strategy

### Split 001: Core Extraction Functionality (~493 lines)
- **Purpose**: Implement core certificate extraction from Kind clusters
- **Components**:
  - Type definitions and interfaces
  - Error handling framework
  - Main extraction logic
- **Dependencies**: None (foundational split)
- **Testing**: Minimal - focus on core functionality

### Split 002: Validation and Testing (~343 lines)
- **Purpose**: Add validation logic and comprehensive test coverage
- **Components**:
  - Certificate validation implementation
  - Reduced but comprehensive test suite
  - Test utilities and mocks
- **Dependencies**: Split 001 (requires types and extractor)
- **Testing**: Full test coverage for entire effort

## Risk Mitigation
- Each split compiles independently
- Split 001 provides working extraction (without validation)
- Split 002 adds validation and ensures quality
- Test reduction strategy: Focus on critical paths, combine similar tests

## Verification Checklist
- [x] No file appears in multiple splits
- [x] All files from original effort covered
- [x] Each split will compile independently
- [x] Dependencies properly ordered (Split 002 depends on 001)
- [x] Each split under 800 lines (target <700)
- [x] Total functionality preserved across splits

## Notes
- Test files will be strategically reduced in Split 002 to stay under limit
- Consider combining similar test cases to reduce line count
- Both splits must be implemented sequentially (002 depends on 001)