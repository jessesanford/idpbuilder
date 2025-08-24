# SPLIT-PLAN-002.md
## Split 002 of 2: Stack Types
**Planner**: Code Reviewer @agent-code-reviewer (same for ALL splits)
**Parent Effort**: registry-auth-types
**Target Size**: 313 lines (well under 800 limit)

## Boundaries
- **Previous Split**: Split 001 (OCI types)
- **This Split Focus**: Stack type definitions and validation
- **Next Split**: None (final split)

## Files in This Split (EXCLUSIVE - no overlap with other splits)
```
pkg/stack/types.go       (107 lines) - Stack type definitions
pkg/stack/constants.go   (42 lines)  - Stack-related constants
pkg/stack/types_test.go  (164 lines) - Unit tests for stack types
```
**Total**: 313 lines

## Functionality Scope
### Core Components:
1. **Stack Type Definitions** (pkg/stack/types.go)
   - Stack configuration structures
   - Component definitions
   - Dependency specifications
   - Validation interfaces

2. **Constants** (pkg/stack/constants.go)
   - Stack component types
   - Status constants
   - Default values
   - Error messages

3. **Test Coverage** (pkg/stack/types_test.go)
   - Complete unit tests for stack types
   - Validation logic tests
   - Edge case coverage
   - Example usage patterns

## Dependencies
- **External**: Standard library only (encoding/json, fmt, etc.)
- **Internal**: None - completely independent of Split 001
- **Test Dependencies**: Standard testing package

## Implementation Instructions for SW Engineer

### Step 1: Create Branch
```bash
git checkout -b phase1/wave1/registry-auth-types-split-002
```

### Step 2: Sparse Checkout (if starting fresh)
```bash
# Enable sparse checkout
git sparse-checkout init --cone
git sparse-checkout set pkg/stack/
```

### Step 3: Verify Files
Ensure ONLY these files are included:
- pkg/stack/types.go
- pkg/stack/constants.go
- pkg/stack/types_test.go

### Step 4: Run Tests
```bash
go test ./pkg/stack/...
```

### Step 5: Measure Size
```bash
/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh
# Should show ~313 lines
```

### Step 6: Commit
```bash
git add pkg/stack/
git commit -m "feat: implement Stack types and validation (split 002)"
git push origin phase1/wave1/registry-auth-types-split-002
```

## Quality Checklist
- [ ] Stack types properly defined with json/yaml tags
- [ ] Validation methods comprehensive
- [ ] Constants cover all stack scenarios
- [ ] Tests achieve >80% coverage
- [ ] No dependencies on OCI package
- [ ] Clean, self-contained implementation

## Implementation Notes
### Key Type Structures:
- **Stack**: Main container for stack configuration
- **Component**: Individual stack components
- **Dependency**: Inter-component dependencies
- **Status**: Stack and component status tracking

### Validation Requirements:
- Component name uniqueness
- Dependency cycle detection
- Required field validation
- Version compatibility checks

## Merge Strategy
- This split will be merged to `phase1/wave1/registry-auth-types` branch
- Can be merged independently of Split 001
- No conflicts expected with other splits

## Parallel Development
- **Independence**: This split can be developed in parallel with Split 001
- **No Coordination**: No shared files or dependencies
- **Testing**: Can be tested independently
- **Integration**: Will combine cleanly with Split 001

## Risk Assessment
- **Low Complexity**: Simple type definitions with clear boundaries
- **Size Safety**: At 313 lines, less than 40% of limit
- **Independence**: No risk of conflicts with other splits
- **Testing**: Comprehensive test coverage included