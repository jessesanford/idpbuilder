# Wave Integration Merge Plan - Phase 2 Wave 1

## Integration Metadata
- **Integration Branch**: `igp/phase2/wave1/integration`
- **Base Branch**: `igp/phase1/integration`
- **Phase**: 2
- **Wave**: 1
- **Created By**: Code Reviewer Agent
- **Created At**: 2025-09-26T23:28:00Z
- **Rules Applied**: R269 (Plan Only), R270 (No Integration Sources)

## Summary of Efforts

All four efforts in Phase 2 Wave 1 have been successfully completed and reviewed:

1. **effort-2.1.1-build-context-management** (281 lines) - Foundation effort, ACCEPTED
2. **effort-2.1.2-multi-stage-build-support** (331 lines) - ACCEPTED
3. **effort-2.1.3-build-caching-implementation** (354 lines) - ACCEPTED
4. **effort-2.1.4-build-options-and-args** (145 lines) - ACCEPTED after formatting fix

**Total Lines**: 1,111 (all within 800 line limit per effort)

## Merge Order and Dependencies

Based on the dependency analysis from orchestrator-state.json:
- **effort-2.1.1** is the foundation (no dependencies)
- **efforts 2.1.2, 2.1.3, 2.1.4** all depend on 2.1.1

### Recommended Merge Sequence

The following merge order ensures dependencies are respected and minimizes conflicts:

#### Merge 1: Build Context Management (Foundation)
```bash
# Merge the foundation effort first
git checkout igp/phase2/wave1/integration
git merge igp/phase2/wave1/effort-2.1.1-build-context-management --no-ff \
  -m "feat: merge build context management foundation (effort 2.1.1)"
```
- **Dependencies**: None (foundation effort)
- **Expected Conflicts**: None (first merge)
- **New Files Added**:
  - `pkg/build/context.go`
  - `pkg/build/context_test.go`
  - Related test fixtures

#### Merge 2: Multi-stage Build Support
```bash
git merge igp/phase2/wave1/effort-2.1.2-multi-stage-build-support --no-ff \
  -m "feat: merge multi-stage build support (effort 2.1.2)"
```
- **Dependencies**: Requires effort-2.1.1 (build context)
- **Expected Conflicts**: Minimal - possible in `pkg/build/` if modifying shared files
- **New Files Added**:
  - `pkg/build/multistage.go`
  - `pkg/build/multistage_test.go`
  - Stage handling utilities

#### Merge 3: Build Caching Implementation
```bash
git merge igp/phase2/wave1/effort-2.1.3-build-caching-implementation --no-ff \
  -m "feat: merge build caching implementation (effort 2.1.3)"
```
- **Dependencies**: Requires effort-2.1.1 (build context)
- **Expected Conflicts**: Possible in build configuration files
- **New Files Added**:
  - `pkg/build/cache.go`
  - `pkg/build/cache_test.go`
  - Cache management utilities

#### Merge 4: Build Options and Arguments
```bash
git merge igp/phase2/wave1/effort-2.1.4-build-options-and-args --no-ff \
  -m "feat: merge build options and arguments (effort 2.1.4)"
```
- **Dependencies**: Requires effort-2.1.1 (build context)
- **Expected Conflicts**: Possible in command-line parsing or build configuration
- **New Files Added**:
  - `pkg/build/options.go`
  - `pkg/build/options_test.go`
  - Argument parsing utilities

## Conflict Resolution Strategy

### Expected Conflict Points
1. **go.mod/go.sum**: Multiple efforts may add dependencies
   - Resolution: Accept all additions (union of dependencies)

2. **pkg/build/builder.go** (if exists): Multiple efforts may modify the main builder
   - Resolution: Carefully merge all feature additions

3. **Test fixtures**: Overlapping test data
   - Resolution: Ensure all test fixtures are preserved

### Resolution Commands Template
```bash
# If conflicts occur:
git status  # Check conflicted files
git diff    # Review conflicts

# For go.mod conflicts:
go mod tidy  # After manual resolution

# For code conflicts:
# 1. Open conflicted files
# 2. Resolve keeping all functionality
# 3. Run tests to verify

# After resolution:
git add <resolved-files>
git commit --no-edit  # Use the merge commit message
```

## Validation Steps (For Integration Agent)

After each merge, validate:

### 1. Build Verification
```bash
make build || go build ./...
```

### 2. Test Execution
```bash
make test || go test ./...
```

### 3. Line Count Verification
```bash
# After all merges, verify total size
cd /home/vscode/workspaces/idpbuilder-gitea-push/efforts/phase2/wave1/integration-workspace
$CLAUDE_PROJECT_DIR/tools/line-counter.sh
# Should show ~1,111 lines of new implementation code
```

### 4. Functionality Verification
```bash
# Verify each feature is present:
grep -r "BuildContext" pkg/build/  # effort 2.1.1
grep -r "MultiStage" pkg/build/    # effort 2.1.2
grep -r "CacheManager" pkg/build/  # effort 2.1.3
grep -r "BuildOptions" pkg/build/  # effort 2.1.4
```

### 5. Integration Testing
```bash
# Run integration tests if available
make integration-test || go test -tags=integration ./...
```

## Final Integration Checklist

Before marking integration complete:

- [ ] All 4 efforts merged successfully
- [ ] No merge conflicts remain unresolved
- [ ] All tests pass
- [ ] Build succeeds
- [ ] Total line count verified (~1,111 lines)
- [ ] All functionality from each effort is present
- [ ] Integration branch pushed to remote
- [ ] Ready for architect review

## Notes for Integration Agent

1. **DO NOT modify this plan** - Execute exactly as specified (R269)
2. **Use original effort branches only** - Never merge from other integration branches (R270)
3. **Preserve commit messages** - Use the suggested merge commit messages
4. **Document any deviations** - If conflicts require different resolution
5. **Run validations after EACH merge** - Don't wait until the end

## Risk Assessment

**Low Risk**: All efforts have been reviewed and accepted, with clean implementations within size limits.

**Potential Issues**:
- Build system conflicts if multiple efforts modified Makefile
- Test fixture overlaps if efforts created similar test data
- Import conflicts in main.go if multiple efforts added commands

**Mitigation**: The merge order (foundation first, then parallel features) minimizes these risks.

---

**Plan Status**: COMPLETE
**Ready for**: Integration Agent Execution
**Compliance**: R269 (Plan Only), R270 (No Integration Sources)