# SPLIT-002 RE-SPLIT PLAN (UPDATED)

## Problem Statement
- **Current Split-002 Size**: 1188 lines (388 lines OVER 800 limit)
- **Target**: Re-split into TWO sub-splits, each under 600 lines (safety margin)
- **Created**: 2025-09-03 06:45:00
- **Updated**: 2025-09-03 14:30:00
- **Planner**: Code Reviewer Agent (ERROR_RECOVERY state)

## Current Implementation Analysis

### File Breakdown (Total: 2987 lines including tests)
```
Production Code (1799 lines):
- doc.go: 45 lines (package documentation)
- options.go: 143 lines (build options)
- builder.go: 205 lines (builder interface)
- builder_impl.go: 228 lines (builder implementation)
- config.go: 274 lines (configuration)
- layer.go: 402 lines (layer creation)
- tarball.go: 490 lines (tarball generation)

Test Code (1188 lines):
- options_test.go: 204 lines
- config_test.go: 404 lines
- builder_test.go: 592 lines
```

## Re-Split Strategy

### Split-002a: Core Layer and Configuration (~573 lines)
**Purpose**: Foundational layer creation and configuration management
**Dependencies**: Split-001 (interfaces and base types)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder-split-002a`

**Files to Include**:
1. **doc.go** (45 lines) - Package documentation
2. **config.go** (218 lines) - Full configuration factory implementation
3. **layer.go** (310 lines) - Complete layer factory implementation
   **TOTAL**: 573 lines ✅ SAFELY UNDER 600

**No Reduction Needed**:
- These three files form a cohesive unit
- Total is already under 600 lines
- No need to split functionality

**Test Files**:
- config_test.go (partial - basic tests only)
- Basic layer tests (new file)

### Split-002b: Builder Implementation and Tarball (~615 lines)
**Purpose**: Complete builder implementation with tarball generation
**Dependencies**: Split-002a (layer and config)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder-split-002b`

**Files to Include**:
1. **tarball.go** (408 lines) - Full tarball writer implementation
2. **builder.go** (174 lines) - Builder interface and SimpleBuilder
3. **options.go** (124 lines) - Build options
   **TOTAL**: 706 lines ⚠️ Close to 700 soft limit but acceptable

**Alternative if 706 is too high**:
1. **tarball.go** (408 lines) - Tarball writer only
2. **builder_impl.go** (198 lines) - Builder implementation only
   **TOTAL**: 606 lines ✅ SAFER OPTION

**Recommended**: Use the 606-line alternative for safety

**Test Files**:
- options_test.go (204 lines)
- builder_test.go (partial - core tests)
- tarball tests (new file)

## Implementation Instructions

### Phase 1: Split-002a Implementation
1. **Setup Infrastructure**:
   ```bash
   cd efforts/phase2/wave1/go-containerregistry-image-builder-SPLIT-002
   git checkout -b idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder-split-002a
   ```

2. **Implement Core Files**:
   - Start with doc.go (package documentation)
   - Implement config.go with validation
   - Create simplified layer.go (core functionality only)
   - Write basic tests for configuration

3. **Size Verification**:
   ```bash
   PROJECT_ROOT=/home/vscode/workspaces/idpbuilder-oci-go-cr
   $PROJECT_ROOT/tools/line-counter.sh -b idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-001 -c [current-branch]
   ```

4. **Target Metrics**:
   - Production code: ~400 lines
   - Test code: ~150 lines
   - Total: ~550 lines

### Phase 2: Split-002b Implementation
1. **Setup Infrastructure**:
   ```bash
   cd efforts/phase2/wave1/go-containerregistry-image-builder-SPLIT-002
   git checkout -b idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder-split-002b
   ```

2. **Implement Builder Components**:
   - Import types from split-002a
   - Implement options.go
   - Complete builder.go and builder_impl.go
   - Implement simplified tarball.go
   - Write comprehensive tests

3. **Size Verification**:
   ```bash
   PROJECT_ROOT=/home/vscode/workspaces/idpbuilder-oci-go-cr
   $PROJECT_ROOT/tools/line-counter.sh -b idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder-split-002a -c [current-branch]
   ```

4. **Target Metrics**:
   - Production code: ~400 lines
   - Test code: ~150 lines
   - Total: ~550 lines

## Dependency Management

### Split-002a Dependencies:
- **External**: Split-001 (interfaces, base types)
- **Provides**: Configuration types, layer creation API
- **Compilation**: Must compile independently

### Split-002b Dependencies:
- **External**: Split-001, Split-002a
- **Provides**: Complete builder implementation
- **Compilation**: Requires split-002a types

## Risk Mitigation

### Size Control Measures:
1. **Continuous Monitoring**: Check size every 100 lines
2. **Early Warning**: Stop at 500 lines to review
3. **Hard Stop**: Absolute limit at 650 lines (safety margin)

### Code Organization:
1. **Minimize Duplication**: Share common utilities
2. **Clean Interfaces**: Well-defined boundaries between splits
3. **Test Strategy**: Focus on integration tests in split-002b

## Verification Checklist

### Split-002a Validation:
- [ ] Size under 600 lines (with margin)
- [ ] Compiles independently
- [ ] Tests pass
- [ ] Provides clean API for split-002b

### Split-002b Validation:
- [ ] Size under 600 lines (with margin)
- [ ] Integrates with split-002a
- [ ] Complete builder functionality
- [ ] All tests pass

## Integration Strategy

After both sub-splits complete:
1. Merge split-002a to integration branch
2. Merge split-002b to integration branch
3. Verify combined functionality
4. Continue with split-003 as originally planned

## Notes for SW Engineer

**CRITICAL**: The existing split-002 implementation MUST be refactored into these two sub-splits. Key considerations:

1. **File Location**: Files should be created in the split's root directory (not pkg/builder/)
2. **Import Paths**: Use relative imports between splits
3. **Test Coverage**: Maintain existing test coverage across both splits
4. **Incremental Development**: Split-002a must be complete before starting split-002b

**Line Counting**: Always use the project's line counter tool:
```bash
PROJECT_ROOT=/home/vscode/workspaces/idpbuilder-oci-go-cr
BASE_BRANCH="[appropriate-base]"
CURRENT_BRANCH=$(git branch --show-current)
$PROJECT_ROOT/tools/line-counter.sh -b $BASE_BRANCH -c $CURRENT_BRANCH
```

## Remaining Files Disposition

**Files NOT included in 002a or 002b**:
- **builder.go** (174 lines) - IF using safer 606-line option for 002b
- **options.go** (124 lines) - IF using safer 606-line option for 002b
- **All test files** (builder_test.go, config_test.go, options_test.go)

**These will go to Split-003 or a new Split-002c for tests**

## Critical Implementation Order

1. **STOP** all work on current oversized split-002
2. **CREATE** split-002a infrastructure first
3. **IMPLEMENT** split-002a completely (573 lines)
4. **CREATE** split-002b infrastructure
5. **IMPLEMENT** split-002b with safer option (606 lines)
6. **DEFER** remaining files to future splits

## Success Criteria

1. **Split-002a**: ≤ 600 lines (actual: 573), compiles, tests pass
2. **Split-002b**: ≤ 650 lines (actual: 606 with safer option), full functionality
3. **Combined**: Core builder functionality operational
4. **No Regression**: Integration with split-001 still works
5. **Size Compliance**: BOTH splits stay well under 800-line hard limit
## Split-002A Implementation Requirements

This is split-002a of the resplit from the oversized split-002 (1188 lines).

### Scope: Layer & Configuration (573 lines target)

Implement ONLY these files:
- layer.go (310 lines) - Layer creation functionality
- config.go (218 lines) - Configuration management
- doc.go (45 lines) - Package documentation

Total: 573 lines (well under 800 limit)

### DO NOT IMPLEMENT:
- tarball.go (reserved for split-002b)
- builder_impl.go (reserved for split-002b)  
- test files (deferred to later splits)

