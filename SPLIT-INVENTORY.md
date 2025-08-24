# Complete Split Plan for registry-auth-types
**Sole Planner**: Code Reviewer Instance @agent-code-reviewer
**Total Size**: 965 lines
**Splits Required**: 2
**Created**: 2025-08-24 21:48:00 UTC

## Executive Summary
The registry-auth-types effort has grown to 965 lines, exceeding the 800-line limit. The implementation includes OCI types (622 lines) and Stack types (313 lines), plus package documentation (39 lines). This split plan divides the effort into two logical, independently compilable units.

## Split Boundaries (NO OVERLAPS)
| Split | Components | Size | Files | Status |
|-------|------------|------|-------|--------|
| 001   | OCI Package + doc.go | 661 lines | 6 files | Planned |
| 002   | Stack Package | 313 lines | 3 files | Planned |

## Detailed File Distribution

### Split 001: OCI Types and Documentation (661 lines)
- `pkg/doc.go` - 39 lines (package documentation)
- `pkg/oci/types.go` - 121 lines (OCI type definitions)
- `pkg/oci/manifest.go` - 124 lines (manifest handling)
- `pkg/oci/constants.go` - 56 lines (OCI constants)
- `pkg/oci/types_test.go` - 130 lines (type tests)
- `pkg/oci/manifest_test.go` - 191 lines (manifest tests)
**Total**: 661 lines

### Split 002: Stack Types (313 lines)
- `pkg/stack/types.go` - 107 lines (stack type definitions)
- `pkg/stack/constants.go` - 42 lines (stack constants)
- `pkg/stack/types_test.go` - 164 lines (stack type tests)
**Total**: 313 lines

## Deduplication Matrix
| File/Module | Split 001 | Split 002 |
|-------------|-----------|-----------|
| pkg/doc.go | ✅ | ❌ |
| pkg/oci/* | ✅ | ❌ |
| pkg/stack/* | ❌ | ✅ |

## Dependency Analysis
- **Split 001**: Self-contained OCI types with no external dependencies
- **Split 002**: Self-contained Stack types with no dependencies on Split 001
- **Independence**: Both splits can be developed and tested in parallel

## Implementation Order
1. **Split 001** can be implemented first or in parallel
2. **Split 002** can be implemented independently
3. No strict ordering required - both are independent

## Verification Checklist
- ✅ No file appears in multiple splits
- ✅ All files from original effort covered (974 lines total)
- ✅ Each split compiles independently
- ✅ Dependencies properly analyzed (none between splits)
- ✅ Split 001: 661 lines (<800 limit)
- ✅ Split 002: 313 lines (<800 limit)
- ✅ Combined total matches original: 661 + 313 = 974 lines

## Risk Assessment
- **Low Risk**: Clean separation between OCI and Stack packages
- **No Cross-Dependencies**: Each package is self-contained
- **Test Coverage**: Both splits include their respective test files
- **Compilation**: Each split will compile independently

## Notes for Orchestrator
- These splits can be assigned to different SW Engineers for parallel work
- No coordination required between splits
- Each split has complete test coverage included
- Total effort remains under 1000 lines as originally estimated