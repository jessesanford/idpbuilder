# Wave Merge Plan - Phase 2 Wave 1

## Overview
**Phase**: 2
**Wave**: 1
**Integration Branch**: `idpbuilderpush/phase2/wave1/integration`
**Base Branch**: `phase1/wave1/integration` (from previous wave)
**Created**: 2025-09-23
**Reviewer**: Code Reviewer Agent

## Pre-Merge Validation

### Environment Requirements
- Working directory: `/home/vscode/workspaces/idpbuilder-push/efforts/phase2/wave1/integration-workspace`
- Current branch: `idpbuilderpush/phase2/wave1/integration`
- Clean working tree (no uncommitted changes)

### Pre-Merge Checklist
```bash
# Verify clean working tree
git status

# Ensure on correct branch
git branch --show-current
# Expected: idpbuilderpush/phase2/wave1/integration

# Verify base is phase1 integration
git log --oneline -1
# Should show: 10cf3ce chore: initialize Phase 2 Wave 1 integration infrastructure
```

## Effort Branches Analysis

### Effort 2.1.1: Auth Interface Tests
- **Branch**: `idpbuilderpush/phase2/wave1/auth-interface-tests`
- **Base Commit**: `28a302c` (feat: Add Claude Code and GitHub CLI to devcontainer)
- **Files Added**:
  - `pkg/oci/auth_test.go` - Test interfaces for authentication
  - `pkg/oci/testdata/fixtures.go` - Test data fixtures
- **Dependencies**: Requires types from auth-implementation
- **Conflicts Expected**: None (new files only)

### Effort 2.1.2: Auth Implementation
- **Branch**: `idpbuilderpush/phase2/wave1/auth-implementation`
- **Base Commit**: `28a302c` (feat: Add Claude Code and GitHub CLI to devcontainer)
- **Files Added**:
  - `pkg/oci/auth.go` - Core authentication logic
  - `pkg/oci/errors.go` - Error definitions
  - `pkg/oci/types.go` - Type definitions
- **Dependencies**: None (foundational)
- **Conflicts Expected**: None (new files only)

### Effort 2.1.3: Auth Mocks
- **Branch**: `idpbuilderpush/phase2/wave1/auth-mocks`
- **Base Commit**: `28a302c` (feat: Add Claude Code and GitHub CLI to devcontainer)
- **Files Added**:
  - `pkg/oci/auth_mock.go` - Mock authentication for testing
  - `pkg/oci/testutil/helpers.go` - Test utility helpers
- **Dependencies**: Uses types from auth-implementation
- **Conflicts Expected**: None (new files only)

## Merge Order Strategy

Based on dependency analysis, the optimal merge order is:

1. **Auth Implementation** (2.1.2) - FIRST
   - Foundational code with type definitions
   - No dependencies on other efforts
   - Provides types needed by other efforts

2. **Auth Interface Tests** (2.1.1) - SECOND
   - Depends on types from auth-implementation
   - Tests the interfaces defined in implementation
   - No conflicts with implementation

3. **Auth Mocks** (2.1.3) - THIRD
   - Depends on types from auth-implementation
   - May reference test interfaces
   - Completes the testing infrastructure

## Merge Commands

**WARNING**: These commands are for the orchestrator to execute. DO NOT run these during planning!

### Step 1: Prepare Integration Branch
```bash
cd /home/vscode/workspaces/idpbuilder-push/efforts/phase2/wave1/integration-workspace
git checkout idpbuilderpush/phase2/wave1/integration
git pull origin idpbuilderpush/phase2/wave1/integration
```

### Step 2: Add Effort Remotes (if needed)
```bash
# Add remotes for each effort workspace
git remote add auth-impl ../auth-implementation/.git || true
git remote add auth-tests ../auth-interface-tests/.git || true
git remote add auth-mocks ../auth-mocks/.git || true

# Fetch all remotes
git fetch auth-impl
git fetch auth-tests
git fetch auth-mocks
```

### Step 3: Merge Auth Implementation
```bash
# Merge implementation first (foundational)
git merge auth-impl/idpbuilderpush/phase2/wave1/auth-implementation \
  --no-ff \
  -m "feat(phase2/wave1): integrate effort 2.1.2 - auth implementation

- Core authentication logic
- Error definitions
- Type definitions
- Foundation for TDD GREEN phase"

# Verify successful merge
git status
git diff --name-status HEAD~1
```

### Step 4: Merge Auth Interface Tests
```bash
# Merge interface tests second
git merge auth-tests/idpbuilderpush/phase2/wave1/auth-interface-tests \
  --no-ff \
  -m "feat(phase2/wave1): integrate effort 2.1.1 - auth interface tests

- Test interfaces for authentication
- Test data fixtures
- TDD RED phase implementation"

# Verify successful merge
git status
git diff --name-status HEAD~1
```

### Step 5: Merge Auth Mocks
```bash
# Merge mocks last
git merge auth-mocks/idpbuilderpush/phase2/wave1/auth-mocks \
  --no-ff \
  -m "feat(phase2/wave1): integrate effort 2.1.3 - auth mocks

- Mock authentication for testing
- Test utility helpers
- Completes testing infrastructure"

# Verify successful merge
git status
git diff --name-status HEAD~1
```

### Step 6: Final Validation
```bash
# Ensure all files are present
ls -la pkg/oci/

# Expected files:
# - auth.go (from implementation)
# - errors.go (from implementation)
# - types.go (from implementation)
# - auth_test.go (from interface tests)
# - testdata/fixtures.go (from interface tests)
# - auth_mock.go (from mocks)
# - testutil/helpers.go (from mocks)

# Run build to verify compilation
go build ./pkg/oci/...

# Run tests to verify integration
go test ./pkg/oci/...
```

### Step 7: Push Integration Branch
```bash
# Push the integrated changes
git push origin idpbuilderpush/phase2/wave1/integration
```

## Conflict Resolution Strategy

### Expected Conflicts
- **None Expected**: All efforts add new files without modifying existing ones
- All efforts work in `pkg/oci/` namespace but touch different files

### If Conflicts Occur
1. **File-level conflicts**: Should not happen (different files)
2. **Import conflicts**: Resolve by including all imports
3. **Type conflicts**: Unlikely, but favor auth-implementation types

### Rollback Procedure
If merge fails:
```bash
# Reset to pre-merge state
git reset --hard 10cf3ce  # Initial integration commit

# Investigate issue
git diff auth-impl/idpbuilderpush/phase2/wave1/auth-implementation
```

## Post-Merge Verification

### Compilation Check
```bash
# Build the project
go build ./...

# Run linter
golangci-lint run ./pkg/oci/...
```

### Test Execution
```bash
# Run all tests
go test -v ./pkg/oci/...

# Run with coverage
go test -cover ./pkg/oci/...
```

### Integration Validation
```bash
# Verify all effort files are present
find pkg/oci -type f -name "*.go" | sort

# Check for any uncommitted changes
git status

# View final integration commit graph
git log --graph --oneline -10
```

## Success Criteria

✅ All three effort branches merged successfully
✅ No merge conflicts encountered
✅ All files from efforts present in integration
✅ Code compiles without errors
✅ Tests pass (or fail as expected in TDD RED phase)
✅ Integration branch pushed to origin

## Notes

### Dependency Graph
```
auth-implementation (types.go, errors.go, auth.go)
    ├── auth-interface-tests (uses types)
    └── auth-mocks (uses types)
```

### Risk Assessment
- **Low Risk**: All efforts add new files without modifying existing code
- **Medium Risk**: Dependencies between efforts require specific merge order
- **Mitigation**: Follow merge order strictly (implementation → tests → mocks)

## Troubleshooting Guide

### Issue: Remote not found
**Solution**: Ensure effort directories exist and contain git repositories
```bash
ls -la ../auth-implementation/.git
ls -la ../auth-interface-tests/.git
ls -la ../auth-mocks/.git
```

### Issue: Branch not found
**Solution**: Verify branch names in effort directories
```bash
cd ../auth-implementation && git branch --show-current
cd ../auth-interface-tests && git branch --show-current
cd ../auth-mocks && git branch --show-current
```

### Issue: Merge conflicts
**Solution**: Should not occur, but if they do:
1. Check file paths for unexpected overlaps
2. Resolve favoring implementation over tests
3. Ensure all imports are preserved

## Approval

This merge plan has been reviewed and validated according to:
- **R269**: Plan created only, no merges executed
- **R270**: Using original effort branches as sources
- **R308**: Incremental branching from phase1/wave1/integration

**Status**: READY FOR EXECUTION by Orchestrator

---

*Generated by Code Reviewer Agent*
*State: WAVE_MERGE_PLANNING*
*Date: 2025-09-23*