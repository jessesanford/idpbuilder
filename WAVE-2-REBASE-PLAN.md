# Phase 1 Wave 2 Rebase Plan (R327 Cascade)

## Executive Summary
This plan addresses the rebase of Phase 1 Wave 2 efforts onto the newly re-integrated Phase 1 Wave 1 branch as part of the R327 cascade pattern. All Wave 2 efforts were originally based on an earlier Wave 1 integration that has since been recreated with backported fixes.

## Current State Analysis

### Wave 1 Integration (New Base)
- **Branch**: `idpbuilder-oci-build-push/phase1/wave1-integration`
- **Location**: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave1/integration-workspace`
- **Latest Commit**: `51ef23b` - docs: finalize R327 CASCADE work log
- **Status**: ✅ Stable, re-integrated with all fixes
- **R327 Compliance**: ✅ Fresh recreation completed

### Wave 2 Efforts Requiring Rebase

#### 1. cert-validation-split-001
- **Current Branch**: `idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-001`
- **Working Directory**: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/cert-validation-split-001`
- **Latest Commit**: `9d05978` - fix(R321): complete cert-validation-split-001 backport analysis
- **Status**: Has R321 fixes applied
- **Lines**: 207 (within limit)

#### 2. cert-validation-split-002
- **Current Branch**: `idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-002`
- **Working Directory**: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/cert-validation-split-002`
- **Latest Commit**: `cf8a9a3` - marker: R321 backport fix complete - test fixtures added
- **Status**: Has R321 fixes applied
- **Lines**: 800 (at limit)
- **Dependencies**: Requires split-001 to be rebased first

#### 3. cert-validation-split-003
- **Current Branch**: `idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-003`
- **Working Directory**: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/cert-validation-split-003`
- **Latest Commit**: `8aaba03` - marker: R321 backport fixes complete for split-003
- **Status**: Has R321 fixes applied
- **Lines**: 800 (at limit)
- **Dependencies**: Requires split-002 to be rebased first

#### 4. fallback-strategies
- **Current Branch**: `idpbuilder-oci-build-push/phase1/wave2/fallback-strategies`
- **Working Directory**: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/fallback-strategies`
- **Latest Commit**: `bdd84a7` - fix(R321): fallback strategy backport analysis complete
- **Status**: Has R321 fixes applied
- **Lines**: 560 (within limit)
- **Dependencies**: Can be rebased independently after cert-validation splits

## Rebase Strategy

### Order of Operations (CRITICAL - Must be Sequential)

The rebase must follow this exact sequence due to split dependencies:

1. **cert-validation-split-001** - First split, no dependencies
2. **cert-validation-split-002** - Depends on split-001
3. **cert-validation-split-003** - Depends on split-002
4. **fallback-strategies** - Independent, but should be last

### Key Principles
- **R327 Compliance**: Each effort must be independently rebased before integration
- **R321 Compliance**: Fixes remain in source branches, not integration
- **Sequential Splits**: Each split builds on the previous one
- **Clean History**: Use interactive rebase to maintain clean commit history

## Detailed Rebase Instructions

### Phase 1: Preparation

#### Step 1.1: Clone Wave 2 Efforts from Target Repository
Since the efforts are on wrong branches locally, we need to clone fresh from target:

```bash
# For each effort, create a fresh clone
for effort in cert-validation-split-001 cert-validation-split-002 cert-validation-split-003 fallback-strategies; do
    echo "=== Cloning $effort ==="

    # Navigate to effort directory
    cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/$effort

    # Backup current state
    cp -r . ../backup-$effort-$(date +%Y%m%d-%H%M%S)

    # Reset to clean state
    rm -rf .git
    git init

    # Add target remote
    git remote add origin https://github.com/jessesanford/idpbuilder.git

    # Fetch the specific branch
    git fetch origin idpbuilder-oci-build-push/phase1/wave2/$effort

    # Checkout the branch
    git checkout -b idpbuilder-oci-build-push/phase1/wave2/$effort origin/idpbuilder-oci-build-push/phase1/wave2/$effort

    # Add integration remote for rebase base
    git remote add integration /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave1/integration-workspace
    git fetch integration
done
```

### Phase 2: Sequential Rebase Execution

#### Step 2.1: Rebase cert-validation-split-001

```bash
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/cert-validation-split-001

# Verify current state
git log --oneline -n 5
git status

# Create backup branch
git branch backup-pre-rebase-$(date +%Y%m%d-%H%M%S)

# Fetch latest integration
git fetch integration
git fetch origin

# Find merge base with old integration
OLD_BASE=$(git merge-base HEAD origin/idpbuilder-oci-build-push/phase1/wave1/integration-20250912-032401)
echo "Old base: $OLD_BASE"

# Perform the rebase onto new integration
git rebase --onto integration/idpbuilder-oci-build-push/phase1/wave1-integration $OLD_BASE HEAD

# If conflicts occur:
# 1. Resolve conflicts favoring the effort's changes for feature code
# 2. Accept integration's version for infrastructure/build files
# 3. Continue with: git rebase --continue

# Verify the rebase
git log --oneline --graph -n 10

# Run build and test verification
go build ./...
go test ./...

# Force push the rebased branch
git push origin idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-001 --force-with-lease
```

#### Step 2.2: Rebase cert-validation-split-002

```bash
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/cert-validation-split-002

# CRITICAL: This split depends on split-001
# First, fetch the rebased split-001
git remote add split001 ../cert-validation-split-001
git fetch split001

# Create backup
git branch backup-pre-rebase-$(date +%Y%m%d-%H%M%S)

# Find the commit where split-002 diverged from split-001
DIVERGENCE_POINT=$(git merge-base HEAD origin/idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-001)
echo "Divergence from split-001: $DIVERGENCE_POINT"

# Rebase onto the newly rebased split-001
git rebase --onto split001/idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-001 $DIVERGENCE_POINT HEAD

# Resolve conflicts if any
# Verify build and tests
go build ./...
go test ./...

# Push rebased branch
git push origin idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-002 --force-with-lease
```

#### Step 2.3: Rebase cert-validation-split-003

```bash
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/cert-validation-split-003

# CRITICAL: This split depends on split-002
# First, fetch the rebased split-002
git remote add split002 ../cert-validation-split-002
git fetch split002

# Create backup
git branch backup-pre-rebase-$(date +%Y%m%d-%H%M%S)

# Find the commit where split-003 diverged from split-002
DIVERGENCE_POINT=$(git merge-base HEAD origin/idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-002)
echo "Divergence from split-002: $DIVERGENCE_POINT"

# Rebase onto the newly rebased split-002
git rebase --onto split002/idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-002 $DIVERGENCE_POINT HEAD

# Resolve conflicts if any
# Verify build and tests
go build ./...
go test ./...

# Push rebased branch
git push origin idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-003 --force-with-lease
```

#### Step 2.4: Rebase fallback-strategies

```bash
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/fallback-strategies

# This effort is independent but should include cert-validation changes
# Fetch the rebased split-003 (final cert-validation state)
git remote add certsplit003 ../cert-validation-split-003
git fetch certsplit003

# Create backup
git branch backup-pre-rebase-$(date +%Y%m%d-%H%M%S)

# Find merge base with old integration
OLD_BASE=$(git merge-base HEAD origin/idpbuilder-oci-build-push/phase1/wave1/integration-20250912-032401)
echo "Old base: $OLD_BASE"

# Rebase onto the cert-validation-split-003 (includes all Wave 1 + cert validation)
git rebase --onto certsplit003/idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-003 $OLD_BASE HEAD

# Resolve conflicts if any
# Verify build and tests
go build ./...
go test ./...

# Push rebased branch
git push origin idpbuilder-oci-build-push/phase1/wave2/fallback-strategies --force-with-lease
```

### Phase 3: Verification

#### Step 3.1: Verify Each Rebased Branch

```bash
for effort in cert-validation-split-001 cert-validation-split-002 cert-validation-split-003 fallback-strategies; do
    echo "=== Verifying $effort ==="
    cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/$effort

    # Check that it's based on new integration
    echo "Checking base..."
    git merge-base HEAD integration/idpbuilder-oci-build-push/phase1/wave1-integration

    # Verify build
    echo "Building..."
    go build ./...

    # Run tests
    echo "Testing..."
    go test ./...

    # Check for R321 fixes still present
    echo "Verifying R321 fixes..."
    git log --oneline | grep -i "R321\|fix\|backport"
done
```

#### Step 3.2: Integration Readiness Check

```bash
# Create temporary integration to verify all rebased branches work together
cd /tmp
rm -rf wave2-integration-test
git clone /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave1/integration-workspace wave2-integration-test
cd wave2-integration-test

# Create test integration branch
git checkout -b test-wave2-integration

# Merge each rebased effort in order
for effort in cert-validation-split-001 cert-validation-split-002 cert-validation-split-003 fallback-strategies; do
    echo "=== Merging $effort ==="
    git remote add $effort /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/$effort
    git fetch $effort
    git merge $effort/idpbuilder-oci-build-push/phase1/wave2/$effort --no-ff -m "test: merge $effort"
done

# Final build and test
go build ./...
go test ./...

echo "✅ All Wave 2 efforts successfully rebased and ready for integration"
```

## Conflict Resolution Strategy

### Expected Conflicts

1. **go.mod/go.sum files**
   - Resolution: Accept the version from the new Wave 1 integration
   - Then run `go mod tidy` to incorporate Wave 2 changes

2. **Certificate validation code**
   - Resolution: Keep Wave 2 implementation changes
   - Ensure R321 fixes are preserved

3. **Import statements**
   - Resolution: Merge both sets of imports
   - Remove duplicates

### Conflict Resolution Guidelines

```bash
# When resolving conflicts:

# 1. For infrastructure/build files (go.mod, Makefile, etc.)
git checkout --theirs <file>  # Accept from integration

# 2. For feature implementation files
git checkout --ours <file>    # Keep effort changes

# 3. For mixed changes
# Manually edit to combine both changes
# Preserve ALL R321 fixes
# Maintain feature functionality

# 4. After resolution
git add <resolved-files>
git rebase --continue
```

## R327 Compliance Checklist

### Pre-Rebase
- [ ] Wave 1 integration is fresh (recreated with all fixes)
- [ ] All Wave 2 efforts have R321 fixes applied
- [ ] Backup branches created for all efforts
- [ ] Target repository accessible

### During Rebase
- [ ] Each effort rebased independently
- [ ] Sequential order maintained for splits
- [ ] Conflicts resolved preserving fixes
- [ ] Build verification after each rebase
- [ ] Test verification after each rebase

### Post-Rebase
- [ ] All efforts build successfully
- [ ] All tests pass
- [ ] R321 fixes still present in commits
- [ ] Branches pushed to origin with --force-with-lease
- [ ] Integration test performed
- [ ] No functionality lost
- [ ] Clean commit history maintained

## Rollback Plan

If any rebase fails catastrophically:

```bash
# For each effort that needs rollback
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/$EFFORT

# Reset to backup branch
git reset --hard backup-pre-rebase-[timestamp]

# Or restore from origin
git fetch origin
git reset --hard origin/idpbuilder-oci-build-push/phase1/wave2/$EFFORT

# Restore from backup directory if needed
rm -rf .
cp -r ../backup-$EFFORT-[timestamp]/. .
```

## Next Steps After Successful Rebase

1. **Create Fresh Wave 2 Integration**
   ```bash
   # Delete old Wave 2 integration
   git push origin --delete idpbuilder-oci-build-push/phase1/wave2/integration

   # Create new integration branch from Wave 1
   git checkout -b idpbuilder-oci-build-push/phase1/wave2-integration integration/idpbuilder-oci-build-push/phase1/wave1-integration

   # Merge rebased efforts
   # Follow WAVE-MERGE-PLAN.md for proper merge sequence
   ```

2. **R327 Cascade Continues**
   - After Wave 2 rebase complete → Re-integrate Wave 2
   - After Wave 2 integration → Delete and recreate Phase 1 integration
   - Continue cascade as needed

3. **Verification Gates**
   - R291: Build, test, demo gates
   - R330: Integration quality review
   - R308: Incremental development verification

## Risk Mitigation

### High Risk Areas
1. **Split Dependencies**: Must maintain sequential order
2. **R321 Fix Preservation**: Critical that fixes aren't lost
3. **Build Compatibility**: go.mod conflicts likely

### Mitigation Strategies
1. **Incremental Verification**: Test after each rebase
2. **Backup Everything**: Multiple backup strategies
3. **Communication**: Clear status updates during process

## Timeline Estimate

- **Preparation**: 30 minutes
- **cert-validation-split-001 rebase**: 30-45 minutes
- **cert-validation-split-002 rebase**: 30-45 minutes
- **cert-validation-split-003 rebase**: 30-45 minutes
- **fallback-strategies rebase**: 30-45 minutes
- **Verification & Integration Test**: 30 minutes
- **Total Estimated Time**: 3-4 hours

## Command Summary (Quick Reference)

```bash
# Quick rebase for each effort (after preparation)
EFFORT=cert-validation-split-001  # or split-002, split-003, fallback-strategies
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/$EFFORT

# For split-001 (base on Wave 1 integration)
git rebase --onto integration/idpbuilder-oci-build-push/phase1/wave1-integration $(git merge-base HEAD origin/idpbuilder-oci-build-push/phase1/wave1/integration-20250912-032401) HEAD

# For split-002 (base on rebased split-001)
git rebase --onto split001/idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-001 $(git merge-base HEAD origin/idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-001) HEAD

# For split-003 (base on rebased split-002)
git rebase --onto split002/idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-002 $(git merge-base HEAD origin/idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-002) HEAD

# For fallback-strategies (base on rebased split-003)
git rebase --onto certsplit003/idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-003 $(git merge-base HEAD origin/idpbuilder-oci-build-push/phase1/wave1/integration-20250912-032401) HEAD

# After each rebase
go build ./... && go test ./...
git push origin idpbuilder-oci-build-push/phase1/wave2/$EFFORT --force-with-lease
```

## Status Tracking

### Current Status: PLAN CREATED
- [ ] Plan reviewed and approved
- [ ] Preparation phase started
- [ ] cert-validation-split-001 rebased
- [ ] cert-validation-split-002 rebased
- [ ] cert-validation-split-003 rebased
- [ ] fallback-strategies rebased
- [ ] Integration test passed
- [ ] Branches pushed to origin
- [ ] Ready for Wave 2 re-integration

---

**Document Created**: 2025-09-14T12:00:00Z
**Author**: Code Reviewer Agent
**Purpose**: R327 Cascade - Phase 1 Wave 2 Rebase onto Re-integrated Wave 1
**Compliance**: R327, R321, R308