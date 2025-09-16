# Phase 2 Wave 1 Integration Merge Plan

## 🎉 UPDATE: Split-002 Successfully Rebased - Integration Simplified!

**Key Change**: `gitea-client-split-002` has been successfully rebased onto the current `gitea-client-split-001`, creating a clean linear history. This dramatically simplifies the integration process - what was previously a complex merge is now straightforward!

## Summary
This plan outlines the merge sequence for integrating Phase 2 Wave 1 efforts into the integration branch.

**Created**: 2025-01-15
**Updated**: 2025-01-15 (Split-002 rebase completed)
**Created By**: Code Reviewer Agent (WAVE_MERGE_PLANNING state)
**Target Branch**: `idpbuilder-oci-build-push/phase2/wave1/integration`
**Base Branch**: `idpbuilder-oci-build-push/phase1/integration`
**Integration Directory**: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/integration-workspace/repo`

### Integration Complexity: ✅ REDUCED
- **Before**: Complex merge required for split-002 (based on old split-001)
- **After**: Simple, clean merges for all branches (linear history achieved)

## Critical Requirements (R269, R270)
- ✅ Use ONLY original effort branches (no integration branches)
- ✅ Exclude parent 'too-large' branches for split efforts
- ✅ Include only split branches for efforts that were split
- ✅ Determine merge order based on dependencies
- ✅ Document conflict resolution strategies

## Branches to Integrate

### Original Effort Branches (per R269/R270)
1. **image-builder**: `idpbuilder-oci-build-push/phase2/wave1/image-builder`
   - Status: Ready for merge
   - Size: Within limits (no splits needed)
   - Base: phase1/integration (with rebase completed)
   - Last commit: 8c68910 (marker: feature flag fix complete)

2. **gitea-client-split-001**: `idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001`
   - Status: Split from original gitea-client (too large)
   - Size: Within limits after split
   - Base: phase1/integration (rebased)
   - Last commit: 4fb2931 (marker: rebase onto phase1/integration complete)

3. **gitea-client-split-002**: `idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002`
   - Status: Second split of gitea-client
   - Size: Within limits
   - Base: **UPDATED** - Now rebased onto CURRENT split-001 (commit 4fb2931)
   - Last commit: Updated with cherry-picked fixes from split-001
   - Contains commits: 4aeb3ec, 2a407b7, 1304ef9 (cherry-picked versions of split-001's recent fixes)

### Excluded Branches
- **gitea-client**: `idpbuilder-oci-build-push/phase2/wave1/gitea-client`
  - Reason: Original effort exceeded size limits, replaced by splits
  - Status: DO NOT MERGE (parent branch of splits)

## Dependency Analysis

### Branch Relationships
1. **image-builder**: Independent, can be merged first
   - Implements build command functionality
   - No dependencies on other Wave 1 efforts

2. **gitea-client-split-001**: Independent from image-builder
   - Provides core registry infrastructure
   - Authentication and basic operations

3. **gitea-client-split-002**: Depends on split-001
   - ✅ **RESOLVED**: Now properly rebased onto the CURRENT split-001 (commit 4fb2931)
   - Enhances push/list operations from split-001
   - Adds retry logic and test improvements
   - Contains all split-001 fixes via cherry-picked commits

### Integration Status Update
✅ **SIMPLIFIED**: gitea-client-split-002 has been successfully rebased onto the current split-001. The commit history is now LINEAR (phase1/integration → split-001 → split-002), eliminating the complex merge requirements.

## Merge Sequence

### Order of Operations
1. **First**: image-builder (independent, clean base)
2. **Second**: gitea-client-split-001 (rebased, clean)
3. **Third**: gitea-client-split-002 (simple merge - now properly based on split-001)

## Detailed Merge Instructions

### Pre-merge Setup
```bash
# 1. Navigate to integration directory
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/integration-workspace/repo

# 2. Ensure clean working directory
git status
# Expected: "nothing to commit, working tree clean"

# 3. Verify current branch
git branch --show-current
# Expected: "idpbuilder-oci-build-push/phase2/wave1/integration"

# 4. Fetch all remote branches
git fetch origin

# 5. Create backup tag before starting
git tag phase2-wave1-integration-backup-$(date +%Y%m%d-%H%M%S)

# 6. Verify base is phase1/integration
git log --oneline -1
# Should show phase1/integration as base
```

### Merge 1: image-builder
```bash
# 1. Start merge
git merge origin/idpbuilder-oci-build-push/phase2/wave1/image-builder --no-ff \
  -m "merge: integrate image-builder effort into Phase 2 Wave 1"

# 2. Expected conflicts (if any)
# - demo-features.sh (likely different demo implementations)
# - work-log.md (different work histories)
# - DEMO-IMPLEMENTATION-COMPLETE.marker
# - DEMO.md, DEMO-RETROFIT-PLAN.md

# 3. Conflict resolution strategy
# For demo-features.sh: Keep both sets of demos, merge content
# For work-log.md: Append both histories chronologically
# For marker files: Keep all markers
# For DEMO files: Merge content, keeping all demonstrations

# 4. If conflicts occur:
# git status  # Check conflicted files
# [resolve conflicts manually]
# git add [resolved files]
# git commit --no-edit

# 5. Verify after merge
git diff HEAD~1 --stat
# Should show image-builder implementation files in pkg/build/

# 6. Test compilation
go build ./...
```

### Merge 2: gitea-client-split-001
```bash
# 1. Start merge
git merge origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001 --no-ff \
  -m "merge: integrate gitea-client-split-001 into Phase 2 Wave 1"

# 2. Expected conflicts
# - demo-features.sh (both branches add demos)
# - work-log.md (both have different histories)
# - go.mod (potential dependency conflicts)
# - DEMO-related files (both efforts have demos)

# 3. Conflict resolution strategy
# For demo-features.sh:
#   - Create sections for image-builder and gitea-client demos
#   - Preserve all demo functionality
# For work-log.md:
#   - Append histories chronologically
#   - Keep all entries for audit trail
# For go.mod:
#   - Keep all dependencies
#   - Resolve version conflicts by taking latest version

# 4. If conflicts occur:
# git status
# [resolve conflicts]
# git add [resolved files]
# git commit --no-edit

# 5. Verify after merge
git diff HEAD~1 --stat
# Should show new files in pkg/registry/
# Should show auth.go, gitea.go, interface.go, etc.

# 6. Test compilation
go build ./...
go test ./pkg/registry/...
```

### Merge 3: gitea-client-split-002 (SIMPLE MERGE - Now Rebased!)
```bash
# 1. ✅ SIMPLIFIED: This merge is now straightforward
# Split-002 has been rebased onto the current split-001
# Linear history: phase1/integration → split-001 → split-002
# Conflicts should be minimal or non-existent

# 2. Start merge (expect clean or minimal conflicts)
git merge origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002 --no-ff \
  -m "merge: integrate gitea-client-split-002 into Phase 2 Wave 1"

# 3. Expected conflicts (MINIMAL)
# Since split-002 now contains all split-001 commits:
# - Possibly demo-features.sh (if modified in multiple branches)
# - Possibly work-log.md (different histories)
# - Most code conflicts should be resolved by the rebase

# 4. Resolution strategy (if any conflicts)
# For demo-features.sh:
#   - Merge demo sections from all branches
#   - Keep all demo functionality
# For work-log.md:
#   - Append histories chronologically
#
# Code conflicts should be rare since split-002 includes split-001

# 5. If minimal conflicts occur:
git status  # Should show few if any conflicts
# [resolve any remaining conflicts]
git add [resolved files if any]
git commit --no-edit

# 6. Verify the merge
git diff HEAD~1 --stat
# Should show:
# - Enhanced push/list in pkg/registry/
# - New retry.go file
# - Mocks moved to test files (R320 compliance)
# - Additional test coverage from split-002

# 7. Test to confirm everything works
go build ./...
go test ./pkg/registry/... -v
```

## Conflict Resolution Guidelines

### General Principles
1. **Functionality First**: Ensure all implemented features are preserved
2. **No Duplication**: Remove duplicate implementations, especially from older bases
3. **Test Coverage**: Keep all test cases from all branches
4. **Demo Completeness**: Merge all demo scripts into organized sections
5. **R320 Compliance**: Ensure no stub implementations remain

### Specific File Strategies

#### demo-features.sh
```bash
# Merge strategy: Combine all demos with clear sections
# Structure:
echo "=== Image Builder Demos ==="
# [image-builder demo code]

echo "=== Gitea Registry Client Demos ==="
# [gitea-client demo code]

echo "=== Integration Demos ==="
# [combined functionality demos]
```

#### work-log.md
```bash
# Merge strategy: Chronological append
# Keep all entries for complete audit trail
# Format: newest entries at top or bottom (be consistent)
```

#### go.mod
```bash
# Merge strategy: Latest versions win
# After merge:
go mod tidy
# This will resolve any dependency conflicts
```

#### pkg/registry/* files
```bash
# Merge strategy for split conflicts:
# 1. split-001 provides base infrastructure (already merged)
# 2. split-002 provides enhancements
# 3. When the same function exists in both:
#    - Take split-002's version (more complete)
# 4. For new files in split-002 (like retry.go):
#    - Keep entirely
```

## Verification Steps

### After Each Merge
```bash
# 1. Compilation check
go build ./...

# 2. Test execution
go test ./... -v

# 3. Line count verification
PROJECT_ROOT=/home/vscode/workspaces/idpbuilder-oci-build-push
$PROJECT_ROOT/tools/line-counter.sh

# 4. Demo validation (if demos exist)
if [ -f demo-features.sh ]; then
  bash demo-features.sh
fi

# 5. Commit verification
git log --oneline -5
```

### Final Integration Verification
```bash
# 1. Full test suite with coverage
go test ./... -v -cover

# 2. Specific package tests
go test ./pkg/build/... -v     # image-builder
go test ./pkg/registry/... -v  # gitea-client

# 3. Size compliance check
$PROJECT_ROOT/tools/line-counter.sh
# Total should be reasonable for Wave 1

# 4. Check all branches are integrated
git log --oneline --graph -20
# Should show merge commits for all three branches

# 5. Verify no stub implementations (R320)
grep -r "panic.*not.*implemented" pkg/
grep -r "TODO.*implement" pkg/
# Should return no results

# 6. Push to remote
git push origin idpbuilder-oci-build-push/phase2/wave1/integration
```

## Rollback Procedures

### If Merge Fails
```bash
# 1. Abort current merge
git merge --abort

# 2. Reset to backup tag
git reset --hard phase2-wave1-integration-backup-[timestamp]

# 3. Investigate issue
git diff origin/[branch-name] --stat
git log origin/[branch-name] --oneline -10

# 4. Retry with adjusted strategy
```

### If Tests Fail After Merge
```bash
# 1. Identify failing tests
go test ./... -v 2>&1 | grep -E "FAIL|Error"

# 2. Check recent changes
git diff HEAD~1

# 3. Decision point:
# Option A: Fix issues in place
# [make fixes]
git add [fixed files]
git commit -m "fix: resolve test failures after merge"

# Option B: Revert if critical
git revert HEAD
```

## Expected Outcomes

### Success Criteria
✅ All three branches merged successfully
✅ No compilation errors
✅ All tests passing
✅ Demo scripts functional (if present)
✅ Size within limits per effort
✅ Clean git history with clear merge commits
✅ No stub implementations (R320 compliance)

### Files Expected After Integration
```
pkg/
├── build/           # Image builder implementation
│   ├── builder.go
│   ├── command.go
│   └── ...
├── registry/        # Complete Gitea registry client
│   ├── auth.go
│   ├── gitea.go
│   ├── interface.go
│   ├── list.go      # Enhanced version from split-002
│   ├── push.go      # Enhanced version from split-002
│   ├── retry.go     # New from split-002
│   ├── mocks_test.go # Mocks moved to tests (R320)
│   └── ...
└── [other packages]

demo-features.sh     # Combined demo suite
go.mod              # Updated dependencies
work-log.md         # Complete history
```

## Risk Assessment

### ✅ UPDATED: Risk Level Reduced
1. **Split-002 merge**: ~~Based on older split-001, high conflict probability~~
   - **RESOLVED**: Now rebased on current split-001, LOW conflict probability
   - Expect straightforward, clean merge

2. **Registry functionality**: Split-002 properly extends split-001
   - Since split-002 now has all split-001 changes, no overlap issues
   - Test to confirm all functionality works together

3. **Test mock refactoring**: Split-002 moves mocks to comply with R320
   - This refactoring is clean and should merge without issues
   - Verify all tests still pass after merge

### Medium Risk Areas
1. **Demo script conflicts**: Multiple demo implementations
   - Mitigation: Organize into clear sections
   - Test each demo section independently

2. **Dependency versions**: Different go.mod requirements
   - Mitigation: Use go mod tidy to resolve
   - Test with final dependency set

## Notes for Orchestrator

✅ **UPDATED GUIDANCE**:
1. This is a PLAN only - actual merge execution by orchestrator/integration agent
2. **Split-002 is now SIMPLE to merge** - properly rebased on split-001
3. Test after each merge to confirm functionality
4. Document the smooth integration in your report
5. The rebase has eliminated the need for complex manual merging

## Appendix: Branch Analysis Details

### image-builder Key Changes
- Build command implementation (~4000 lines total)
- Demo implementations for R291 compliance
- Feature flag fixes (removed blocking flags)
- TLS configuration improvements
- Fixed duplicate TLSConfig struct issues

### gitea-client-split-001 Key Changes
- Core registry infrastructure (~5000 lines added)
- Authentication implementation (auth.go)
- Basic Gitea client (gitea.go)
- Interface definitions (interface.go)
- Basic push/pull operations
- Remote options configuration

### gitea-client-split-002 Key Changes
- Enhanced push/list operations (improved versions)
- Retry logic implementation (retry.go - new file)
- Mock refactoring for R320 compliance (moved to test files)
- Additional test coverage
- Removed redundant oci/types.go
- Test infrastructure improvements
- **NOW INCLUDES**: All fixes from split-001 via rebase
  - Commit 4aeb3ec: Cherry-picked fix from split-001
  - Commit 2a407b7: Cherry-picked fix from split-001
  - Commit 1304ef9: Cherry-picked fix from split-001

### Known Issues Resolved
- image-builder: Feature flag blocking issue fixed
- image-builder: Duplicate TLSConfig resolved
- gitea-client-split-001: Rebased onto phase1/integration
- gitea-client-split-002: R320 compliance (no stubs in production)

---
*End of Phase 2 Wave 1 Merge Plan*
*Next Steps: Execute merges according to this plan*