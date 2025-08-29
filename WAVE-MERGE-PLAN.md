# Phase 2 Wave 1 Integration Merge Plan

## Overview
**Phase**: 2  
**Wave**: 1  
**Integration Branch**: `idpbuilder-oci-mvp/phase2/wave1/integration`  
**Integration Directory**: `/home/vscode/workspaces/idpbuilder-oci-mvp/efforts/phase2/wave1/integration-workspace/idpbuilder`  
**Created**: 2025-01-29  
**Planner**: Code Reviewer Agent  

## Critical Instructions (R269 Compliance)
- This plan is for the Integration Agent to EXECUTE
- Code Reviewer creates the plan ONLY - does NOT execute merges
- Use ORIGINAL effort branches, NOT integration branches
- Exclude 'too-large' branches, include only their splits

## Branches to Merge

### Summary
Total branches to merge: **3**
1. `idpbuilder-oci-mvp/phase2/wave1/buildah-build-wrapper-split-001` (516 lines - Split 1 of 2)
2. `idpbuilder-oci-mvp/phase2/wave1/buildah-build-wrapper-split-002` (484 lines - Split 2 of 2)
3. `idpbuilder-oci-mvp/phase2/wave1/gitea-registry-client` (736 lines - No splits needed)

### Excluded Branches
- `idpbuilder-oci-mvp/phase2/wave1/buildah-build-wrapper` - EXCLUDED (983 lines exceeded limit, replaced by splits)

## Merge Order Analysis

### Branch Dependencies
All three branches share the same merge base: `67b4b08` (feat: upgrade ingress-nginx (#537))

### File Overlap Analysis
**CRITICAL CONFLICT DETECTED**: Splits 001 and 002 both modify the same files:
- `pkg/build/builder.go`
- `pkg/build/builder_buildah.go`
- `pkg/build/types.go`
- `go.mod` and `go.sum`

**No Conflicts Expected**: Gitea Registry Client modifies separate files:
- `pkg/registry/gitea_client.go` (new file)
- `pkg/registry/gitea_client_test.go` (new file)
- `pkg/registry/types.go` (new file)

### Recommended Merge Order
1. **First**: `gitea-registry-client` - Independent, no conflicts expected
2. **Second**: `buildah-build-wrapper-split-001` - First part of buildah implementation
3. **Third**: `buildah-build-wrapper-split-002` - Builds on split-001, will have conflicts to resolve

## Detailed Merge Instructions

### Pre-Merge Verification
```bash
# Ensure you're in the integration directory
cd /home/vscode/workspaces/idpbuilder-oci-mvp/efforts/phase2/wave1/integration-workspace/idpbuilder

# Verify you're on the integration branch
git branch --show-current
# Expected: idpbuilder-oci-mvp/phase2/wave1/integration

# Ensure branch is clean
git status
# Should show no uncommitted changes

# Fetch latest changes
git fetch origin
```

### Merge 1: Gitea Registry Client
```bash
# Merge the gitea-registry-client branch
git merge origin/idpbuilder-oci-mvp/phase2/wave1/gitea-registry-client --no-ff \
  -m "merge: integrate gitea-registry-client (736 lines) - Phase 2 Wave 1"

# Expected result: Clean merge, no conflicts
# New files added:
#   - pkg/registry/gitea_client.go
#   - pkg/registry/gitea_client_test.go
#   - pkg/registry/types.go

# Validate the merge
git status
git diff HEAD~1 --stat
# Expected: ~1338 insertions, 197 deletions

# Run tests to ensure functionality
go test ./pkg/registry/...
```

### Merge 2: Buildah Build Wrapper Split 001
```bash
# Merge the first split
git merge origin/idpbuilder-oci-mvp/phase2/wave1/buildah-build-wrapper-split-001 --no-ff \
  -m "merge: integrate buildah-build-wrapper-split-001 (516 lines) - Phase 2 Wave 1"

# Expected result: Clean merge or minor conflicts in go.mod/go.sum
# Files modified:
#   - pkg/build/builder.go
#   - pkg/build/builder_buildah.go (new)
#   - pkg/build/types.go
#   - pkg/build/builder_basic_test.go

# If conflicts in go.mod/go.sum:
# 1. Keep both sets of dependencies
# 2. Run: go mod tidy
# 3. Stage resolved files: git add go.mod go.sum

# Validate the merge
git status
git diff HEAD~1 --stat
# Expected: ~1473 insertions

# Run tests
go test ./pkg/build/...
```

### Merge 3: Buildah Build Wrapper Split 002
```bash
# Merge the second split
git merge origin/idpbuilder-oci-mvp/phase2/wave1/buildah-build-wrapper-split-002 --no-ff \
  -m "merge: integrate buildah-build-wrapper-split-002 (484 lines) - Phase 2 Wave 1"

# EXPECTED CONFLICTS in:
#   - pkg/build/builder.go
#   - pkg/build/builder_buildah.go
#   - pkg/build/types.go
#   - go.mod and go.sum

# Conflict Resolution Strategy:
# 1. For pkg/build files:
#    - Split 002 contains the COMPLETE implementation
#    - Accept changes from split-002 (incoming) for all build files
#    - This replaces split-001's partial implementation with the full version
#
# 2. For go.mod/go.sum:
#    - Keep all dependencies from both branches
#    - Run: go mod tidy after resolution

# Resolution commands:
git status  # List conflicted files

# For each pkg/build file, accept split-002 version:
git checkout --theirs pkg/build/builder.go
git checkout --theirs pkg/build/builder_buildah.go
git checkout --theirs pkg/build/types.go
git add pkg/build/builder.go pkg/build/builder_buildah.go pkg/build/types.go

# For go.mod and go.sum, manually merge or:
# Edit go.mod to include all dependencies, then:
go mod tidy
git add go.mod go.sum

# Complete the merge
git commit --no-edit

# Validate the merge
git status
git diff HEAD~1 --stat
# Expected: ~1887 insertions total from split-002

# Run comprehensive tests
go test ./pkg/build/...
go test ./pkg/registry/...
```

## Post-Merge Validation

### 1. Verify All Changes Integrated
```bash
# Check total changes from all merges
git diff 67b4b08..HEAD --stat

# Expected total:
# - buildah-build-wrapper functionality (1000 lines combined from splits)
# - gitea-registry-client functionality (736 lines)
# - Total: ~1736 lines of implementation
```

### 2. Run Full Test Suite
```bash
# Run all tests
go test ./...

# Expected: All tests pass
# - Build tests from splits
# - Registry tests from gitea-client
```

### 3. Build Verification
```bash
# Ensure the project builds
go build ./...

# Run linting if configured
golangci-lint run ./...
```

### 4. Documentation Check
```bash
# Verify all implementation plans are present
ls -la *IMPLEMENTATION-PLAN.md *REVIEW-REPORT.md

# Check for any merge markers
grep -r "<<<<<<< HEAD" pkg/
# Should return nothing
```

## Expected Final State

After all merges complete:
- **Total Implementation**: ~1736 lines (within combined limits)
- **Buildah Build Wrapper**: Complete implementation from both splits
- **Gitea Registry Client**: Full OCI registry integration
- **All Tests**: Passing
- **No Merge Conflicts**: Remaining

## Troubleshooting Guide

### If Merge Conflicts Persist
1. Check the exact conflict markers
2. For Split conflicts: Always prefer split-002 (complete implementation)
3. For dependency conflicts: Keep all unique dependencies, run `go mod tidy`
4. Never mix partial implementations from different splits

### If Tests Fail After Merge
1. Check if all files from both splits are present
2. Verify no duplicate function definitions
3. Ensure types.go has complete type definitions
4. Run `go mod download` to ensure all dependencies available

### Recovery from Failed Merge
```bash
# If a merge goes wrong, abort and restart:
git merge --abort
git reset --hard HEAD
# Then retry the merge with more careful conflict resolution
```

## Integration Agent Next Steps

After successful completion of all merges:
1. Push the integration branch to remote
2. Create pull request to main branch
3. Document integration completion in orchestrator state
4. Report back to orchestrator with success status

---

**Important Notes**:
- This plan follows R269: Code Reviewer creates plan only, does not execute
- All line counts are from actual branch measurements
- Conflict resolution strategy prioritizes complete implementations over partial
- The order minimizes conflicts while respecting dependencies