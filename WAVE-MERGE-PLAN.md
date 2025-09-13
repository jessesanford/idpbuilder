# Phase 1 Wave 1 Integration Merge Plan - R327 CASCADE RE-INTEGRATION

## R327 CASCADE RE-INTEGRATION NOTICE
**Type**: Cascade Re-integration per R327
**Reason**: Previous integration invalidated due to fixes applied to source branches
**Strategy**: Clean re-integration from scratch with updated branches
**Date**: 2025-09-13
**Integration Branch**: idpbuilder-oci-build-push/phase1/wave1/integration
**Cascade ID**: WAVE1-CASCADE-20250913

## Executive Summary

This plan outlines the merge strategy for integrating Phase 1, Wave 1 efforts into a unified integration branch following the R327 cascade re-integration protocol. All efforts have been reviewed, R321 fixes applied, and branches are ready for integration. The original `registry-auth-types` branch has been properly split and only the split branches will be merged per R269/R270.

## Integration Metadata
- **Integration Branch**: `idpbuilder-oci-build-push/phase1/wave1/integration`
- **Base Branch**: `main`
- **Integration Directory**: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave1/integration-workspace/repo`
- **Created By**: Code Reviewer Agent (WAVE_MERGE_PLANNING state)
- **Compliance**: R269, R270, R307, R327

## Branches to Merge

### Included Branches (4 total)

#### 1. kind-cert-extraction
- **Branch**: `origin/idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction`
- **Base Commit**: `e210954` (main)
- **Latest Commit**: `0ed41bf` - fix(R321): add testutil import to git_repository_test.go
- **Functionality**: Certificate extraction from Kind clusters
- **Key Files**: `pkg/certs/extractor.go`, `pkg/certs/kind_client.go`, `pkg/certs/storage.go`
- **Status**: R321 fixes applied, ready for merge

#### 2. registry-auth-types-split-001
- **Branch**: `origin/idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-001`
- **Base Commit**: `67b4b08` (main)
- **Latest Commit**: `02fc6c7` - fix(R321): consolidate duplicate contains() function
- **Functionality**: OCI types and manifest handling
- **Key Files**: `pkg/oci/types.go`, `pkg/oci/manifest.go`, `pkg/oci/constants.go`
- **Status**: R321 fixes applied, ready for merge

#### 3. registry-auth-types-split-002
- **Branch**: `origin/idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-002`
- **Base Commit**: From split-001 (`301bf14`)
- **Latest Commit**: `c0634a5` - fix(R321): remove major duplicate test helpers from split-002
- **Functionality**: Certificate types and TLS configuration
- **Key Files**: `pkg/certs/types.go`, `pkg/certs/constants.go`, `pkg/certs/test_helpers.go`
- **Status**: R321 fixes applied, depends on split-001

#### 4. registry-tls-trust
- **Branch**: `origin/idpbuilder-oci-build-push/phase1/wave1/registry-tls-trust`
- **Base Commit**: `e210954` (main)
- **Latest Commit**: `676e220` - fix(R321): create shared testutil package for Wave 1
- **Functionality**: TLS trust store integration
- **Key Files**: `pkg/certs/trust.go`, `pkg/certs/utilities.go`
- **Status**: R321 fixes applied, ready for merge

### Excluded Branches
- **registry-auth-types** (original): EXCLUDED per R269/R270 - Replaced by split branches

## Dependency Analysis

### Functional Dependencies
1. **registry-auth-types-split-001** (Foundation)
   - No dependencies on other Wave 1 efforts
   - Provides core OCI types needed by others
   - Must merge first to establish OCI foundation

2. **registry-auth-types-split-002** (Extends split-001)
   - Direct dependency on split-001 (branched from it)
   - Adds certificate types on top of OCI types
   - Must merge immediately after split-001

3. **kind-cert-extraction** (Independent with type usage)
   - Can use types from splits 1 & 2
   - Provides extraction functionality
   - Can merge after auth type splits

4. **registry-tls-trust** (Integration layer)
   - Uses types from auth splits
   - Integrates with extraction functionality
   - Should merge last

### R321 Fix Consolidation
All branches have received R321 fixes for:
- Duplicate test helper consolidation
- Shared `pkg/testutil/helpers.go` package
- Import path updates in test files

## Merge Order and Strategy

### CRITICAL: Sequential Merge Order
Based on dependency analysis and R269/R270 split requirements:

```
┌─────────────────────────────────┐
│ 1. registry-auth-types-split-001 │ ← Foundation (OCI types)
└─────────────────────────────────┘
                ↓
┌─────────────────────────────────┐
│ 2. registry-auth-types-split-002 │ ← Builds on split-001
└─────────────────────────────────┘
                ↓
┌─────────────────────────────────┐
│ 3. kind-cert-extraction          │ ← Uses types from splits
└─────────────────────────────────┘
                ↓
┌─────────────────────────────────┐
│ 4. registry-tls-trust            │ ← Integration layer
└─────────────────────────────────┘
```

## Expected Conflicts and Resolution

### 1. Test Helper Consolidation (`pkg/testutil/helpers.go`)
**Affected merges**: Steps 3 and 4
**Conflict type**: Multiple branches create/modify this file
**Resolution strategy**:
```bash
# During merge conflicts:
# 1. Keep ALL unique helper functions
# 2. Remove true duplicates only
# 3. Ensure consistent package declaration
# Example resolution:
git checkout --theirs pkg/testutil/helpers.go  # If identical
# OR manually merge keeping all unique functions
```

### 2. Test File Imports
**Affected files**: `pkg/cmd/get/secrets_test.go`, `pkg/controllers/localbuild/argo_test.go`
**Conflict type**: Updated imports from R321 fixes
**Resolution strategy**:
```go
// Keep consolidated testutil imports
import "github.com/cnoe-io/idpbuilder/pkg/testutil"
// Remove any duplicate local test helpers
```

### 3. Go Module Dependencies
**Affected merges**: All steps
**Resolution strategy**:
```bash
# After each merge:
go mod tidy
go mod verify
# Commit the updated go.mod and go.sum
```

## Merge Execution Commands

```bash
# Prerequisites
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave1/integration-workspace/repo
git checkout idpbuilder-oci-build-push/phase1/wave1/integration
git pull origin idpbuilder-oci-build-push/phase1/wave1/integration

# R327 Cascade Tracking
echo "R327 CASCADE: Starting Wave 1 re-integration" > R327-CASCADE.log
date >> R327-CASCADE.log

# Step 1: Merge registry-auth-types-split-001 (Foundation)
echo "=== Step 1: Merging registry-auth-types-split-001 ===" | tee -a R327-CASCADE.log
git merge origin/idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-001 \
  --no-ff -m "merge: registry-auth-types-split-001 - OCI types foundation (R327 cascade)"
# Expected: Clean merge, no conflicts
# Verify: Check pkg/oci/ directory created with types

# Step 2: Merge registry-auth-types-split-002 (Certificate Types)
echo "=== Step 2: Merging registry-auth-types-split-002 ===" | tee -a R327-CASCADE.log
git merge origin/idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-002 \
  --no-ff -m "merge: registry-auth-types-split-002 - certificate types (R327 cascade)"
# Expected: Possible testutil conflicts
# If conflicts in pkg/testutil/helpers.go:
#   git status
#   # Manually resolve keeping all unique functions
#   git add pkg/testutil/helpers.go
#   git commit --no-edit

# Step 3: Merge kind-cert-extraction (Extraction Logic)
echo "=== Step 3: Merging kind-cert-extraction ===" | tee -a R327-CASCADE.log
git merge origin/idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction \
  --no-ff -m "merge: kind-cert-extraction - certificate extraction (R327 cascade)"
# Expected: testutil conflicts likely
# Resolution:
#   # Check if testutil/helpers.go has conflicts
#   # Keep all unique functions from both sides
#   # Remove only true duplicates
#   git add -A
#   git commit --no-edit

# Step 4: Merge registry-tls-trust (Trust Integration)
echo "=== Step 4: Merging registry-tls-trust ===" | tee -a R327-CASCADE.log
git merge origin/idpbuilder-oci-build-push/phase1/wave1/registry-tls-trust \
  --no-ff -m "merge: registry-tls-trust - trust integration (R327 cascade)"
# Expected: Final testutil consolidation needed
# Resolution:
#   # Final cleanup of test helpers
#   # Ensure no duplicate functions remain
#   git add -A
#   git commit --no-edit

# Post-merge cleanup
echo "=== Post-merge cleanup ===" | tee -a R327-CASCADE.log
go mod tidy
go mod verify

# Validation
echo "=== Running validation ===" | tee -a R327-CASCADE.log
go build ./...
go test ./pkg/certs/... -v
go test ./pkg/oci/... -v

# Final status
echo "=== Integration complete ===" | tee -a R327-CASCADE.log
date >> R327-CASCADE.log
git log --oneline -10 >> R327-CASCADE.log
```

## Conflict Resolution Detailed Guide

### When encountering pkg/testutil/helpers.go conflicts:

1. **Examine both versions**:
```bash
git diff HEAD...MERGE_HEAD -- pkg/testutil/helpers.go
```

2. **If files are identical** (likely after R321 fixes):
```bash
git checkout --ours pkg/testutil/helpers.go
git add pkg/testutil/helpers.go
```

3. **If files differ**, manually merge:
```go
package testutil

// Keep all unique functions from both sides
// Remove only exact duplicates
// Ensure consistent formatting
```

4. **Verify no duplicate symbols**:
```bash
go build ./pkg/testutil
```

## Validation Checklist

### During Merge Process
- [ ] Each merge completes without stopping
- [ ] Conflicts resolved preserving all functionality
- [ ] Test helpers consolidated without duplicates
- [ ] No files accidentally deleted

### Post-Merge Validation
- [ ] `go mod tidy` runs without errors
- [ ] `go build ./...` compiles successfully
- [ ] `go test ./pkg/certs/...` passes
- [ ] `go test ./pkg/oci/...` passes
- [ ] `go vet ./...` reports no issues
- [ ] No uncommitted changes remain

### Integration Completeness
- [ ] All 4 branches merged
- [ ] R327-CASCADE.log created with tracking
- [ ] Integration branch pushed to origin
- [ ] Ready for Phase 1 completion

## R327 Cascade Compliance

### Cascade Requirements Met
- ✅ Previous integration properly abandoned
- ✅ Clean integration branch initialized
- ✅ All source branches have R321 fixes
- ✅ Split branches used instead of original (R269/R270)
- ✅ Merge order respects dependencies
- ✅ Cascade tracking log maintained

### Cascade Validation Points
1. **Pre-merge**: All branches at latest R321-fixed commits
2. **During merge**: Conflicts resolved without losing fixes
3. **Post-merge**: All tests pass with consolidated code
4. **Final**: Integration branch ready for phase completion

## Risk Mitigation

### Before Starting
- Ensure integration branch is at correct base
- Verify all source branches fetched
- Create backup branch: `git branch integration-backup-$(date +%Y%m%d)`

### During Merges
- Run tests after each merge
- Commit immediately after successful merge
- Document any unexpected conflicts in R327-CASCADE.log

### If Major Issues Occur
- Stop the merge process
- Reset to backup: `git reset --hard integration-backup-YYYYMMDD`
- Document issues in R327-CASCADE-ISSUES.md
- Escalate to orchestrator for guidance

## Success Criteria

The R327 cascade re-integration is successful when:

1. **All branches merged**: 4 branches integrated in order
2. **Code compiles**: `go build ./...` succeeds
3. **Tests pass**: All test suites execute successfully
4. **No duplicates**: Test helpers properly consolidated
5. **Clean working tree**: No uncommitted changes
6. **Cascade tracked**: R327-CASCADE.log documents process
7. **Branch pushed**: Integration branch available on origin

## Notes and Warnings

### Critical Reminders
- **DO NOT** merge the original `registry-auth-types` branch
- **DO NOT** change the merge order (splits must be sequential)
- **DO NOT** skip validation steps between merges
- **DO** document any deviations in R327-CASCADE.log

### R321 Fix Preservation
All R321 fixes must be preserved during merge:
- Test helper consolidation to pkg/testutil
- Import path updates in test files
- Removal of duplicate test functions
- Consistent test utility usage

### Expected Timeline
- Merge execution: ~15 minutes
- Conflict resolution: ~10-20 minutes
- Validation: ~10 minutes
- Total: ~45 minutes

---

**Plan Created**: 2025-09-13 04:05:00 UTC
**Created By**: Code Reviewer Agent (WAVE_MERGE_PLANNING state)
**Compliance**: R269, R270, R307, R327
**Cascade ID**: WAVE1-CASCADE-20250913