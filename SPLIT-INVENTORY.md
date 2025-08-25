# Complete Split Plan for oci-stack-types
**Sole Planner**: Code Reviewer Instance code-reviewer-1756148460
**Full Path**: phase1/wave1/oci-stack-types
**Parent Branch**: idpbuidler-oci-mgmt/phase1/wave1/oci-stack-types
**Total Size**: 915 lines (excluding tests)
**Splits Required**: 2
**Created**: 2025-08-25 19:01:10 UTC

## CRITICAL SPLIT INTEGRITY NOTICE
ALL splits below belong to THIS effort ONLY: phase1/wave1/oci-stack-types
NO splits should reference efforts outside this path!

## Split Strategy: Functional Separation
The split strategy divides the effort into core contracts/types and advanced types/validation:
- **Split 001**: Core interfaces and base types that establish the foundation
- **Split 002**: Stack-specific types, advanced configurations, and all validation logic

## Split Boundaries (NO OVERLAPS)
| Split | Lines | Files | Focus | Status |
|-------|-------|-------|-------|--------|
| 001   | ~460  | interfaces.go, base types from types.go | Core contracts & foundational types | Planned |
| 002   | ~455  | remaining types.go, validation.go | Stack types & validation | Planned |

## File Distribution Matrix
| File | Split 001 | Split 002 | Notes |
|------|-----------|-----------|-------|
| pkg/oci/api/interfaces.go | ✅ (149 lines) | ❌ | Complete file |
| pkg/oci/api/types.go (part 1) | ✅ (~311 lines) | ❌ | Lines 1-311: Core types |
| pkg/oci/api/types.go (part 2) | ❌ | ✅ (~141 lines) | Lines 312-453: Stack & progress types |
| pkg/oci/api/validation.go | ❌ | ✅ (314 lines) | Complete file |
| pkg/oci/api/types_test.go (split 1) | ✅ (~100 lines) | ❌ | Tests for core types |
| pkg/oci/api/validation_test.go | ❌ | ✅ (~100 lines) | Validation tests |

## Deduplication Matrix
| Type/Interface | Split 001 | Split 002 | Notes |
|----------------|-----------|-----------|-------|
| OCIBuildService | ✅ | ❌ | Interface definition |
| OCIRegistryService | ✅ | ❌ | Interface definition |
| StackOCIManager | ✅ | ❌ | Interface definition |
| ProgressReporter | ✅ | ❌ | Interface definition |
| LayerProcessor | ✅ | ❌ | Interface definition |
| BuildConfig | ✅ | ❌ | Core config type |
| RegistryConfig | ✅ | ❌ | Core config type |
| BuildRequest | ✅ | ❌ | Core request type |
| BuildResult | ✅ | ❌ | Core response type |
| BuildStatus | ✅ | ❌ | Core status type |
| BuildPhase | ✅ | ❌ | Enum type |
| BuildOptions | ✅ | ❌ | Options type |
| PushOptions | ✅ | ❌ | Options type |
| PullOptions | ✅ | ❌ | Options type |
| LayerInfo | ✅ | ❌ | Core info type |
| ImageInfo | ✅ | ❌ | Core info type |
| StackOCIConfig | ❌ | ✅ | Stack-specific |
| StackImageInfo | ❌ | ✅ | Stack-specific |
| StackHistoryEntry | ❌ | ✅ | Stack-specific |
| ProgressEvent | ❌ | ✅ | Progress tracking |
| All validation functions | ❌ | ✅ | All validation logic |

## Compilation Independence
- **Split 001**: Fully compilable standalone - contains all interfaces and core types
- **Split 002**: Depends on Split 001 for base types but adds stack-specific extensions

## Integration Strategy
1. Split 001 merged first - establishes foundation
2. Split 002 imports Split 001's types and adds stack-specific functionality
3. Both splits combined provide complete type system for dependent efforts

## Verification Checklist
- [x] No file appears in multiple splits (types.go logically divided)
- [x] All files from original effort covered
- [x] Each split independently compilable
- [x] Dependencies properly ordered (001 before 002)
- [x] Each split <800 lines (001: ~460, 002: ~455)
- [x] Clear functional separation maintained
- [x] No duplication of types between splits

## Next Steps
1. Create SPLIT-PLAN-001.md for core interfaces and base types
2. Create SPLIT-PLAN-002.md for stack types and validation
3. SW Engineer implements Split 001 first
4. After Split 001 passes review, implement Split 002
5. Integrate both splits back to parent branch