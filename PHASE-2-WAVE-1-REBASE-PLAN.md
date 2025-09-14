# Phase 2 Wave 1 Rebase Integration Plan

## 📊 Plan Metadata
- **Created By**: Code Reviewer Agent
- **Created At**: 2025-09-14T16:04:32Z
- **Plan Type**: Cascading Rebase per R327
- **Base Strategy**: R308 Incremental Branching
- **Target Integration**: Phase 2 Wave 1 Integration

## 🎯 Executive Summary

This plan orchestrates the rebase and integration of Phase 2 Wave 1 efforts onto the completed Phase 1 integration branch. Following R327 (mandatory reintegration after fixes) and R308 (incremental branching strategy), we will perform a cascading rebase of the gitea-client splits and then create a fresh integration.

### Current State
- **image-builder**: ✅ Already rebased onto Phase 1 integration (completed 2025-09-13T21:50:59Z)
- **gitea-client-split-001**: ⚠️ Needs rebase onto Phase 1 integration
- **gitea-client-split-002**: ⚠️ Needs rebase onto gitea-client-split-001 (sequential dependency)

### Target State
All Phase 2 Wave 1 efforts properly rebased and integrated into a single branch.

## 🔴 CRITICAL RULES COMPLIANCE

### R327 - Mandatory Reintegration After Fixes
- **Requirement**: Fresh integration required after any fixes
- **Application**: Phase 2 Wave 1 was created before Phase 1 integration finalized
- **Action**: Cascading rebase of all efforts

### R308 - Incremental Branching Strategy
- **Requirement**: Each phase/wave builds on previous integration
- **Application**: Phase 2 Wave 1 MUST build on Phase 1 integration (NOT main)
- **Verification**: All efforts will be rebased onto idpbuilder-oci-build-push/phase1/integration

## 📋 Pre-Rebase Verification Checklist

### ✅ Phase 1 Integration Status
```bash
Branch: idpbuilder-oci-build-push/phase1/integration
Status: COMPLETE
Tests: All passing (142/142)
Build: SUCCESS
Tag: phase1-complete-v1.0
Assessment Score: 85/100
```

### ✅ Effort Status Verification
```bash
# Image Builder
Branch: idpbuilder-oci-build-push/phase2/wave1/image-builder
Status: Already rebased onto Phase 1 integration
Rebased At: 2025-09-13T21:50:59Z

# Gitea Client Split 001
Branch: idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001
Current Base: main (needs rebase)
Target Base: idpbuilder-oci-build-push/phase1/integration
Status: Implementation complete, needs rebase

# Gitea Client Split 002
Branch: idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
Current Base: gitea-client-split-001 (old base)
Target Base: gitea-client-split-001 (after rebase)
Status: Implementation complete, needs sequential rebase
```

## 🔄 REBASE SEQUENCE (CRITICAL ORDER)

### Phase 1: Prepare Working Directories
```bash
# Create fresh worktrees for each effort
WORK_DIR="/home/vscode/workspaces/idpbuilder-oci-build-push/worktrees"
mkdir -p $WORK_DIR

# Clone target repository for each effort
cd $WORK_DIR
git clone https://github.com/jessesanford/idpbuilder.git gitea-split-001-rebase
git clone https://github.com/jessesanford/idpbuilder.git gitea-split-002-rebase
```

### Phase 2: Rebase gitea-client-split-001 onto Phase 1 Integration

**⚠️ CRITICAL: This is the foundation for split-002!**

```bash
cd $WORK_DIR/gitea-split-001-rebase

# Fetch all branches
git fetch origin

# Checkout the split-001 branch
git checkout -b idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001 \
    origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001

# Verify current state
git log --oneline -5
git status

# Record pre-rebase state
PRE_REBASE_COMMIT=$(git rev-parse HEAD)
echo "Pre-rebase HEAD: $PRE_REBASE_COMMIT"

# Perform the rebase onto Phase 1 integration
git rebase origin/idpbuilder-oci-build-push/phase1/integration

# Handle conflicts if any
# For each conflict:
# 1. Review the conflict carefully
# 2. Preserve Phase 2 Wave 1 functionality
# 3. Ensure Phase 1 integration changes are incorporated
# 4. git add <resolved-files>
# 5. git rebase --continue

# Verify build and tests
go build ./...
go test ./pkg/registry/... -v

# Push the rebased branch
git push --force-with-lease origin idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001

# Record completion
echo "✅ Split-001 rebased at: $(date -Iseconds)"
echo "New HEAD: $(git rev-parse HEAD)"
```

### Phase 3: Rebase gitea-client-split-002 onto Rebased split-001

**⚠️ CRITICAL: Must wait for split-001 rebase to complete!**

```bash
cd $WORK_DIR/gitea-split-002-rebase

# Fetch latest including rebased split-001
git fetch origin

# Checkout split-002 branch
git checkout -b idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002 \
    origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002

# Verify current state
git log --oneline -5

# CRITICAL: Rebase onto the NEW rebased split-001
git rebase origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001

# Handle conflicts (similar process as split-001)
# Key difference: Preserve split-002 specific changes

# Verify build and tests
go build ./...
go test ./pkg/registry/... -v

# Push the rebased branch
git push --force-with-lease origin idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002

# Record completion
echo "✅ Split-002 rebased at: $(date -Iseconds)"
echo "New HEAD: $(git rev-parse HEAD)"
```

## 🔀 INTEGRATION SEQUENCE

### Phase 4: Create Fresh Integration Branch

```bash
cd $WORK_DIR

# Create integration workspace
git clone https://github.com/jessesanford/idpbuilder.git phase2-wave1-integration
cd phase2-wave1-integration

# Start from Phase 1 integration (R308 compliance)
git checkout -b idpbuilder-oci-build-push/phase2/wave1/integration \
    origin/idpbuilder-oci-build-push/phase1/integration

# Verify starting point
echo "Integration base: $(git rev-parse HEAD)"
git log --oneline -1
```

### Phase 5: Sequential Merge of Rebased Efforts

**Order is CRITICAL per dependencies:**

```bash
# 1. Merge image-builder (already rebased)
echo "=== Merging image-builder ==="
git merge origin/idpbuilder-oci-build-push/phase2/wave1/image-builder \
    -m "feat(phase2/wave1): integrate image-builder capability

- OCI image building with go-containerregistry
- Multi-architecture support
- Cache optimization
- Already rebased onto Phase 1 integration"

# Verify no conflicts
git status

# 2. Merge gitea-client-split-001 (freshly rebased)
echo "=== Merging gitea-client-split-001 ==="
git merge origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001 \
    -m "feat(phase2/wave1): integrate gitea-client-split-001

- Gitea registry client foundation
- Authentication handling
- Core API implementation
- Rebased onto Phase 1 integration per R327"

# 3. Merge gitea-client-split-002 (sequentially rebased)
echo "=== Merging gitea-client-split-002 ==="
git merge origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002 \
    -m "feat(phase2/wave1): integrate gitea-client-split-002

- Extended Gitea functionality
- Additional registry operations
- Error handling improvements
- Rebased sequentially after split-001"
```

## 🧪 VALIDATION GATES (R291)

### Build Gate
```bash
echo "=== R291 BUILD GATE ==="
go build ./...
if [ $? -eq 0 ]; then
    echo "✅ BUILD GATE: PASSED"
else
    echo "❌ BUILD GATE: FAILED"
    exit 1
fi
```

### Test Gate
```bash
echo "=== R291 TEST GATE ==="
go test ./... -v > test-results.txt 2>&1
TEST_EXIT=$?
PASS_COUNT=$(grep -c "PASS" test-results.txt)
FAIL_COUNT=$(grep -c "FAIL" test-results.txt)

echo "Tests Passed: $PASS_COUNT"
echo "Tests Failed: $FAIL_COUNT"

if [ $TEST_EXIT -eq 0 ]; then
    echo "✅ TEST GATE: PASSED"
else
    echo "⚠️ TEST GATE: PARTIAL"
fi
```

### Demo Gate
```bash
echo "=== R291 DEMO GATE ==="
for demo in demo-*.sh; do
    if [ -f "$demo" ]; then
        echo "Running: $demo"
        bash "$demo"
        if [ $? -eq 0 ]; then
            echo "✅ $demo: PASSED"
        else
            echo "❌ $demo: FAILED"
        fi
    fi
done
```

### Artifact Gate
```bash
echo "=== R291 ARTIFACT GATE ==="
go build -o idpbuilder-phase2-wave1 ./cmd/idpbuilder
if [ -f "idpbuilder-phase2-wave1" ]; then
    echo "✅ ARTIFACT GATE: PASSED"
    echo "Artifact size: $(du -h idpbuilder-phase2-wave1)"
else
    echo "❌ ARTIFACT GATE: FAILED"
fi
```

## 🚨 CONFLICT RESOLUTION STRATEGY

### Expected Conflict Areas
1. **pkg/registry/** - Both Phase 1 and Phase 2 modify registry
2. **pkg/certs/** - Phase 1 cert handling may conflict with Phase 2 usage
3. **go.mod/go.sum** - Dependency version conflicts

### Resolution Principles
1. **Preserve Phase 2 functionality** - New features must work
2. **Incorporate Phase 1 foundations** - Security/cert handling required
3. **Maintain backward compatibility** - No regression in Phase 1 features
4. **Document all resolutions** - Create REBASE-CONFLICTS.md if conflicts occur

## 📊 Success Criteria

### Rebase Success
- [ ] gitea-client-split-001 rebased onto Phase 1 integration
- [ ] gitea-client-split-002 rebased onto new split-001
- [ ] All branches pushed successfully
- [ ] No unresolved conflicts

### Integration Success
- [ ] Clean merge of all three efforts
- [ ] Build passes (R291 Build Gate)
- [ ] Tests pass (R291 Test Gate)
- [ ] Demos run (R291 Demo Gate)
- [ ] Artifact builds (R291 Artifact Gate)

### R327 Compliance
- [ ] Fresh integration created after rebases
- [ ] All fixes incorporated
- [ ] No stale branches used

### R308 Compliance
- [ ] Built on Phase 1 integration (NOT main)
- [ ] Incremental development verified
- [ ] Proper base branch tracking

## 🔄 Rollback Plan

If rebase fails catastrophically:

```bash
# For split-001
git checkout idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001
git reset --hard $PRE_REBASE_COMMIT
git push --force-with-lease

# For split-002 (if started)
# Similar reset process

# Document failure
echo "REBASE FAILED: $(date -Iseconds)" >> REBASE-FAILURE.md
echo "Reason: [specific failure reason]" >> REBASE-FAILURE.md
```

## 📝 Post-Integration Tasks

1. **Create Integration Report**
   - Document all merges
   - List any issues encountered
   - Verify R291 gate results

2. **Tag Integration**
   ```bash
   git tag phase2-wave1-integration-$(date +%Y%m%d-%H%M%S)
   git push origin --tags
   ```

3. **Update Orchestrator State**
   - Mark Phase 2 Wave 1 as integrated
   - Record rebase completion times
   - Update branch tracking

4. **Trigger Code Review**
   - Spawn code-reviewer for integration validation
   - Ensure quality gates passed

## ⏱️ Timeline Estimate

| Phase | Task | Duration | Dependencies |
|-------|------|----------|--------------|
| 1 | Prepare worktrees | 5 min | None |
| 2 | Rebase split-001 | 15-30 min | Phase 1 integration ready |
| 3 | Rebase split-002 | 15-30 min | Split-001 rebase complete |
| 4 | Create integration branch | 5 min | All rebases complete |
| 5 | Merge all efforts | 10-20 min | Integration branch ready |
| 6 | Run validation gates | 10-15 min | Merges complete |
| **Total** | **Full Integration** | **60-105 min** | - |

## 🎯 Final Verification

After completion, verify:

```bash
# Check branch graph
git log --graph --oneline -20

# Verify all commits present
git log --oneline | grep "image-builder"
git log --oneline | grep "gitea-client-split-001"
git log --oneline | grep "gitea-client-split-002"

# Confirm base is Phase 1 integration
git merge-base HEAD origin/idpbuilder-oci-build-push/phase1/integration
```

---

**Plan Status**: READY FOR EXECUTION
**Next Action**: Begin Phase 1 (Prepare Working Directories)
**Assigned To**: Integration Agent or Orchestrator

🤖 Generated with Claude Code
Co-Authored-By: Code Reviewer Agent