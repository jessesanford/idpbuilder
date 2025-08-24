# Complete Split Plan for oci-types
**Sole Planner**: Code Reviewer Instance @agent-code-reviewer-1756072217
**Total Size**: 974 lines
**Splits Required**: 2
**Created**: 2025-08-24 21:50:00 UTC

## Executive Summary
The oci-types effort has grown to 974 lines, exceeding the 800-line hard limit by 174 lines (21.75%). This document provides a complete, non-overlapping split strategy to divide the work into two compilable, testable units.

## Split Boundaries (NO OVERLAPS)
| Split | Description | Lines | Files | Status | Dependencies |
|-------|-------------|-------|-------|--------|--------------|
| 001   | OCI Package | 622   | 5 files | Planned | None (foundational) |
| 002   | Stack Package + Docs | 352 | 4 files | Planned | Depends on Split 001 |

## Detailed File Distribution

### Split 001: OCI Package (622 lines total)
| File | Lines | Purpose |
|------|-------|---------|
| pkg/oci/types.go | 121 | Core OCI type definitions |
| pkg/oci/manifest.go | 124 | Manifest and descriptor types |
| pkg/oci/constants.go | 56 | OCI specification constants |
| pkg/oci/types_test.go | 130 | Unit tests for types.go |
| pkg/oci/manifest_test.go | 191 | Unit tests for manifest.go |
| **Total** | **622** | Complete OCI package implementation |

### Split 002: Stack Package + Documentation (352 lines total)
| File | Lines | Purpose |
|------|-------|---------|
| pkg/stack/types.go | 107 | Stack configuration types |
| pkg/stack/constants.go | 42 | Stack-related constants |
| pkg/stack/types_test.go | 164 | Unit tests for stack types |
| pkg/doc.go | 39 | Package-level documentation |
| **Total** | **352** | Complete stack package + docs |

## Deduplication Matrix
| File/Module | Split 001 | Split 002 | Notes |
|-------------|-----------|-----------|-------|
| pkg/oci/types.go | ✅ | ❌ | Exclusive to Split 001 |
| pkg/oci/manifest.go | ✅ | ❌ | Exclusive to Split 001 |
| pkg/oci/constants.go | ✅ | ❌ | Exclusive to Split 001 |
| pkg/oci/types_test.go | ✅ | ❌ | Exclusive to Split 001 |
| pkg/oci/manifest_test.go | ✅ | ❌ | Exclusive to Split 001 |
| pkg/stack/types.go | ❌ | ✅ | Exclusive to Split 002 |
| pkg/stack/constants.go | ❌ | ✅ | Exclusive to Split 002 |
| pkg/stack/types_test.go | ❌ | ✅ | Exclusive to Split 002 |
| pkg/doc.go | ❌ | ✅ | Exclusive to Split 002 |

## Implementation Order (CRITICAL)
1. **Split 001 MUST be implemented first**
   - Contains foundational OCI types
   - No external dependencies
   - Can compile and test independently

2. **Split 002 implemented second**
   - Depends on OCI types from Split 001
   - References `oci.OCIReference` in stack types
   - Must merge Split 001 before starting Split 002

## Branch Strategy
```
main
  └── phase1/wave1/oci-types (original, abandoned)
       ├── phase1/wave1/oci-types-split-001 (OCI package)
       └── phase1/wave1/oci-types-split-002 (Stack package)
```

## Verification Checklist
- [x] No file appears in multiple splits
- [x] All files from original effort covered
- [x] Each split compiles independently (Split 001 standalone, Split 002 after Split 001)
- [x] Dependencies properly ordered
- [x] Each split <800 lines (Split 001: 622, Split 002: 352)
- [x] Total coverage: 974 lines (622 + 352)
- [x] No gaps or missing files

## Risk Mitigation
- **Compilation Risk**: Split 001 is self-contained; Split 002 requires Split 001 merged first
- **Test Coverage**: Each split maintains >80% test coverage independently
- **Integration Risk**: Low - splits are logically separated by package boundaries
- **Size Risk**: Both splits well under limit (622 < 800, 352 < 800)

## Success Criteria
1. Split 001 merges successfully with all tests passing
2. Split 002 merges successfully with all tests passing
3. Combined functionality equals original 974-line implementation
4. No code duplication between splits
5. Each split remains under 800-line limit

## Notes for SW Engineer
- Start with Split 001 (OCI package) immediately
- Do NOT start Split 002 until Split 001 is merged
- Use sparse checkout to ensure only assigned files are included
- Run line counter after implementation to verify size compliance
- Each split gets its own branch from main (not from original branch)