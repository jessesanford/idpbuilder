# SPLIT-PLAN-002.md
## Split 002 of 2: Stack Package and Documentation
**Planner**: Code Reviewer @agent-code-reviewer-1756072217 (same for ALL splits)
**Parent Effort**: oci-types
**Branch Name**: phase1/wave1/oci-types-split-002
**Size**: 352 lines (well under 800 limit)

## Boundaries
- **Previous Split**: Split 001 (OCI package) - MUST BE MERGED FIRST
- **This Split Coverage**: Complete Stack package + package documentation
- **Next Split**: None (final split)

## Files in This Split (EXCLUSIVE - no overlap with other splits)

### Production Code (188 lines)
1. **pkg/stack/types.go** (107 lines)
   - StackConfiguration struct with Name, Version, Description
   - Components []StackComponent array
   - Dependencies []StackDependency array
   - StackComponent with Name, Type, Version, OCIReference
   - Configuration map[string]interface{} for flexible config
   - StackDependency type definitions
   - StackStatus enum type

2. **pkg/stack/constants.go** (42 lines)
   - Component type constants
   - Status values (Ready, Pending, Failed, etc.)
   - Default configuration keys
   - Version constraint constants

3. **pkg/doc.go** (39 lines)
   - Package overview documentation
   - Usage examples in comments
   - Interface documentation
   - Package-level constants documentation

### Test Code (164 lines)
4. **pkg/stack/types_test.go** (164 lines)
   - Stack configuration validation tests
   - Component creation tests
   - Dependency resolution tests
   - Status transition tests
   - Edge cases and error conditions

## Directory Structure
```
efforts/phase1/wave1/oci-types-split-002/
├── pkg/
│   ├── stack/
│   │   ├── types.go       (107 lines)
│   │   ├── constants.go   (42 lines)
│   │   └── types_test.go  (164 lines)
│   └── doc.go             (39 lines)
├── go.mod
├── go.sum
└── Makefile
```

## Functionality Scope
This split implements the stack configuration system:
- Stack definitions for multi-component applications
- Component management with OCI references
- Dependency tracking and resolution
- Status management for stack lifecycle
- Package-level documentation
- Comprehensive test coverage (>80%)

## Dependencies
- **External**: None (standard library only)
- **Internal**: Requires Split 001 (imports oci.OCIReference type)
- **Required by**: None (terminal split)

## CRITICAL PREREQUISITE
⚠️ **Split 001 MUST be merged to main before starting this split**
- This split references `oci.OCIReference` from Split 001
- Cannot compile without Split 001's types available
- Orchestrator must ensure proper sequencing

## Implementation Instructions

### 1. Verify Split 001 is Merged
```bash
cd /home/vscode/workspaces/idpbuilder-oci-mgmt
git checkout main
git pull origin main
# Verify pkg/oci types exist in main branch
ls pkg/oci/types.go || echo "ERROR: Split 001 not merged yet!"
```

### 2. Create New Branch
```bash
git checkout main
git pull origin main
git checkout -b phase1/wave1/oci-types-split-002
```

### 3. Create Sparse Checkout (ONLY these files)
```bash
# Create effort directory structure
mkdir -p efforts/phase1/wave1/oci-types-split-002/pkg/stack

# Copy ONLY the Stack package files
cp efforts/phase1/wave1/oci-types/pkg/stack/types.go \
   efforts/phase1/wave1/oci-types-split-002/pkg/stack/

cp efforts/phase1/wave1/oci-types/pkg/stack/constants.go \
   efforts/phase1/wave1/oci-types-split-002/pkg/stack/

cp efforts/phase1/wave1/oci-types/pkg/stack/types_test.go \
   efforts/phase1/wave1/oci-types-split-002/pkg/stack/

# Copy package documentation
cp efforts/phase1/wave1/oci-types/pkg/doc.go \
   efforts/phase1/wave1/oci-types-split-002/pkg/

# Copy supporting files
cp efforts/phase1/wave1/oci-types/go.mod \
   efforts/phase1/wave1/oci-types-split-002/

cp efforts/phase1/wave1/oci-types/go.sum \
   efforts/phase1/wave1/oci-types-split-002/

cp efforts/phase1/wave1/oci-types/Makefile \
   efforts/phase1/wave1/oci-types-split-002/
```

### 4. Update Import Paths
```bash
# Ensure imports reference the merged OCI package from Split 001
# In pkg/stack/types.go, imports should include:
# "github.com/idpbuilder/idpbuilder-oci-mgmt/pkg/oci"
```

### 5. Verify Compilation
```bash
cd efforts/phase1/wave1/oci-types-split-002
go build ./pkg/stack/...
go test ./pkg/stack/... -v
```

### 6. Measure Size
```bash
/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh
# Must show ≤352 lines (excluding work-log.md and plan files)
```

### 7. Run Tests
```bash
go test ./pkg/stack/... -cover
# Must achieve >80% coverage
```

## Quality Requirements
- ✅ All code must compile without errors
- ✅ All tests must pass
- ✅ Test coverage must exceed 80%
- ✅ Total size must not exceed 352 lines
- ✅ All exported types must have godoc comments
- ✅ Proper integration with Split 001 types

## Commit Message
```
feat(stack): implement stack configuration types (split 002/002)

- Add stack configuration types (StackConfiguration, StackComponent)
- Add dependency management types
- Add status enum and constants
- Add package-level documentation
- Add comprehensive unit tests (>80% coverage)
- Size: 352 lines (under 800 limit)

Completes oci-types effort split (depends on split 001)
```

## Review Checklist
- [ ] Only Stack package + doc.go files included (no OCI files)
- [ ] All 4 files present and correct
- [ ] Split 001 merged to main first
- [ ] Imports correctly reference merged OCI types
- [ ] Compilation successful
- [ ] Tests passing with >80% coverage
- [ ] Size verified with line-counter.sh
- [ ] Ready for merge

## Integration Notes
- This split depends on types from Split 001
- StackComponent.OCIReference field uses oci.OCIReference type
- Both splits together complete the original 974-line effort
- After merge, full oci-types functionality is available

## Success Metrics
1. Compilation successful with OCI types available
2. All tests pass including integration points
3. Size remains under 352 lines
4. Combined with Split 001 = complete 974-line implementation
5. No code duplication or gaps