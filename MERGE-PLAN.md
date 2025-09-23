# Phase 1 Wave 1 Integration Merge Plan

## Executive Summary
This document provides a comprehensive merge plan for integrating all three efforts from Phase 1 Wave 1 of the idpbuilder-push project into the integration branch.

**Created**: 2025-09-23 14:41:48 UTC
**Planner**: Code Reviewer Agent
**Target Branch**: phase1/wave1/integration
**Source Branches**:
- idpbuilderpush/phase1/wave1/command-tests (Effort 1.1.1)
- idpbuilderpush/phase1/wave1/command-skeleton (Effort 1.1.2)
- idpbuilderpush/phase1/wave1/integration-tests (Effort 1.1.3)

## Merge Sequence Strategy

### Rationale for Merge Order
The merge sequence has been carefully designed based on:
1. **Dependency relationships**: Tests before implementation
2. **File conflict minimization**: Merging in order of increasing complexity
3. **Feature completeness**: Each merge adds complete functionality

### Recommended Merge Sequence

#### Step 1: Merge Effort 1.1.1 (Write Command Tests)
**Branch**: idpbuilderpush/phase1/wave1/command-tests
**Priority**: FIRST
**Rationale**: Foundation tests with no dependencies on other efforts

**Files Added/Modified**:
- `cmd/push/root_test.go` (NEW - no conflicts expected)
- `EFFORT-PLAN.md` (effort-specific, no conflict)
- `IMPLEMENTATION-COMPLETE.marker` (effort-specific)
- `work-log.md` (effort-specific)

**Expected Conflicts**: NONE - This effort only adds new test files

#### Step 2: Merge Effort 1.1.2 (Command Skeleton)
**Branch**: idpbuilderpush/phase1/wave1/command-skeleton
**Priority**: SECOND
**Rationale**: Core implementation that tests will validate

**Files Added/Modified**:
- `cmd/push/config.go` (NEW)
- `cmd/push/root.go` (NEW)
- `cmd/push/root_test.go` (POTENTIAL CONFLICT with 1.1.1)
- Various effort-specific files

**Expected Conflicts**:
- `cmd/push/root_test.go`: Both 1.1.1 and 1.1.2 modify this file
  - Resolution Strategy: Keep both test sets, merge test functions

#### Step 3: Merge Effort 1.1.3 (Integration Tests)
**Branch**: idpbuilderpush/phase1/wave1/integration-tests
**Priority**: THIRD
**Rationale**: Integration tests build on both previous efforts

**Files Added/Modified**:
- `cmd/push/integration_test.go` (NEW)
- `cmd/push/test_harness.go` (NEW)
- `cmd/push/config.go` (CONFLICT with 1.1.2)
- `cmd/push/root.go` (CONFLICT with 1.1.2)
- `run_tests.sh` (NEW)

**Expected Conflicts**:
- `cmd/push/config.go`: Different PushConfig structures
  - 1.1.2 has: RegistryURL field
  - 1.1.3 missing: RegistryURL field
  - Resolution: Keep RegistryURL from 1.1.2 (more complete)
- `cmd/push/root.go`: Potential implementation differences
  - Resolution: Merge functionality, keeping all command features

## Detailed Merge Instructions

### Pre-Merge Validation
```bash
# 1. Ensure clean working directory
cd /home/vscode/workspaces/idpbuilder-push/efforts/phase1/wave1/integration-workspace
git status
# Expected: "nothing to commit, working tree clean"

# 2. Update integration branch
git checkout phase1/wave1/integration
git pull origin phase1/wave1/integration

# 3. Fetch all effort branches
git fetch origin idpbuilderpush/phase1/wave1/command-tests:idpbuilderpush/phase1/wave1/command-tests
git fetch origin idpbuilderpush/phase1/wave1/command-skeleton:idpbuilderpush/phase1/wave1/command-skeleton
git fetch origin idpbuilderpush/phase1/wave1/integration-tests:idpbuilderpush/phase1/wave1/integration-tests

# 4. Create backup tag
git tag phase1-wave1-pre-integration-$(date +%Y%m%d-%H%M%S)
```

### Merge Execution

#### Merge 1: Command Tests (1.1.1)
```bash
# 1. Start merge
git merge --no-ff idpbuilderpush/phase1/wave1/command-tests \
  -m "feat(phase1/wave1): integrate effort 1.1.1 - write command tests"

# 2. Verify no conflicts
git status
# Expected: "All conflicts fixed but you are still merging" should NOT appear

# 3. Run tests
go test ./cmd/push/... -v
# Expected: Tests pass

# 4. Commit if needed
git commit --no-edit  # Only if merge paused
```

#### Merge 2: Command Skeleton (1.1.2)
```bash
# 1. Start merge
git merge --no-ff idpbuilderpush/phase1/wave1/command-skeleton \
  -m "feat(phase1/wave1): integrate effort 1.1.2 - command skeleton"

# 2. Handle expected conflict in root_test.go
# If conflict occurs:
git status
# Edit cmd/push/root_test.go to keep both test sets
# Ensure all test functions are preserved
vim cmd/push/root_test.go  # Or preferred editor

# 3. After resolving conflicts (if any)
git add cmd/push/root_test.go
git commit --no-edit

# 4. Verify build
go build ./cmd/push/...
# Expected: Build succeeds

# 5. Run tests
go test ./cmd/push/... -v
# Expected: All tests pass
```

#### Merge 3: Integration Tests (1.1.3)
```bash
# 1. Start merge
git merge --no-ff idpbuilderpush/phase1/wave1/integration-tests \
  -m "feat(phase1/wave1): integrate effort 1.1.3 - integration tests"

# 2. Handle expected conflicts
git status
# Expected conflicts in:
# - cmd/push/config.go
# - cmd/push/root.go

# 3. Resolve config.go conflict
# Keep RegistryURL field from 1.1.2
# Merge other fields appropriately
vim cmd/push/config.go

# 4. Resolve root.go conflict
# Merge command implementations
# Keep all functionality from both branches
vim cmd/push/root.go

# 5. Stage resolved files
git add cmd/push/config.go cmd/push/root.go
git commit --no-edit

# 6. Run full test suite
go test ./cmd/push/... -v
./run_tests.sh  # If present
# Expected: All tests pass
```

### Post-Merge Validation

```bash
# 1. Verify all files are present
ls -la cmd/push/
# Should contain:
# - config.go
# - root.go
# - root_test.go
# - integration_test.go
# - test_harness.go

# 2. Run comprehensive tests
go test ./... -v -cover
# Expected: All tests pass with good coverage

# 3. Build verification
go build -o idpbuilder-push ./main.go
./idpbuilder-push push --help
# Expected: Help text displays correctly

# 4. Check for compilation warnings
go vet ./...
# Expected: No warnings

# 5. Run linter if available
golangci-lint run ./...  # If installed
# Expected: No critical issues

# 6. Create integration report
cat > INTEGRATION-REPORT.md << EOF
# Phase 1 Wave 1 Integration Report
- Date: $(date -Iseconds)
- All efforts merged successfully
- Tests: PASSING
- Build: SUCCESSFUL
- Conflicts resolved: 3 files
EOF

# 7. Push integration branch
git push origin phase1/wave1/integration
```

## Conflict Resolution Guidelines

### For config.go Conflicts
```go
// Correct merged structure should include:
type PushConfig struct {
    RegistryURL string  // From 1.1.2 - KEEP THIS
    Username    string  // Common to both
    Password    string  // Common to both
    Namespace   string  // Common to both
    Dir         string  // Common to both
    Insecure    bool    // Common to both
    PlainHTTP   bool    // Common to both
}
```

### For root_test.go Conflicts
- Keep ALL test functions from both branches
- Ensure no duplicate function names
- Merge TestMain if both exist
- Preserve all test data and fixtures

### For root.go Conflicts
- Merge command initialization
- Keep all flags from both implementations
- Preserve all RunE implementations
- Ensure proper error handling from both

## Rollback Strategy

If critical issues are encountered during integration:

### Immediate Rollback
```bash
# 1. Reset to pre-integration state
git reset --hard phase1-wave1-pre-integration-[timestamp]

# 2. Document the issue
cat > INTEGRATION-FAILURE-REPORT.md << EOF
# Integration Failure Report
- Date: $(date -Iseconds)
- Failed at: [Step X]
- Error: [Description]
- Action Required: [Fix needed]
EOF

# 3. Push the reset
git push --force origin phase1/wave1/integration
```

### Partial Rollback (if only last merge failed)
```bash
# 1. Undo last merge
git reset --hard HEAD~1

# 2. Re-attempt with different conflict resolution
# Or return to orchestrator for guidance
```

## Success Criteria

Integration is considered successful when:
- [ ] All three efforts are merged
- [ ] No unresolved conflicts remain
- [ ] All tests pass (unit and integration)
- [ ] Build completes without errors
- [ ] Command runs with --help flag
- [ ] No compilation warnings
- [ ] Integration branch is pushed to origin

## Risk Mitigation

### Identified Risks
1. **Config Structure Divergence**: Different PushConfig definitions
   - Mitigation: Use superset of fields from all efforts

2. **Test Function Collisions**: Same test names in different efforts
   - Mitigation: Rename conflicting tests with effort prefixes

3. **Import Path Issues**: Different import organizations
   - Mitigation: Standardize imports post-merge

### Contingency Plans
- If merge conflicts are more complex than expected: Document specific conflicts and request developer input
- If tests fail post-merge: Isolate failing tests and determine if issue is from merge or original implementation
- If build fails: Check for missing dependencies and update go.mod if needed

## Notes for Integration Agent

1. **Always work in the integration-workspace directory**
2. **Create detailed logs of each merge step**
3. **If uncertain about conflict resolution, document and escalate**
4. **Test after EACH merge, not just at the end**
5. **Preserve effort-specific documentation files**
6. **Use --no-ff to maintain merge history**

## Appendix: Quick Command Reference

```bash
# Check current state
git status
git log --oneline -10

# View conflicts
git diff --name-only --diff-filter=U

# Test commands
go test ./cmd/push/... -v
go build ./cmd/push/...

# Abort merge if needed
git merge --abort

# View branch graph
git log --graph --oneline --all -20
```

---
END OF MERGE PLAN