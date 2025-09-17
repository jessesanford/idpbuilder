# Phase 1 Wave 2 Integration Merge Plan

## Overview
**Created**: 2025-09-17 01:14 UTC
**Reviewer**: Code Reviewer Agent
**Phase**: 1
**Wave**: 2
**Integration Branch**: `idpbuilder-oci-build-push/phase1/wave2-integration`
**Base Branch**: Wave 1 Integration (referenced as commit 6e80b35 per R308)

## Executive Summary
This plan integrates four efforts from Phase 1 Wave 2, including three sequential cert-validation splits and the fallback-strategies effort. All efforts have been analyzed and are ready for integration following R327 requirements (fresh integration after bug fixes).

## Efforts to Integrate

### 1. cert-validation-split-001
- **Branch**: `idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-001`
- **Location**: `efforts/phase1/wave2/cert-validation-split-001`
- **Latest Commit**: 9ed933b (fix(R321): complete cert-validation-split-001 backport analysis)
- **Core Components**: Base certificate types, trust management, storage, utilities

### 2. cert-validation-split-002
- **Branch**: `idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-002`
- **Location**: `efforts/phase1/wave2/cert-validation-split-002`
- **Latest Commit**: ff14e39 (marker: R321 backport fix complete - test fixtures added)
- **Core Components**: Built on split-001, adds test fixtures and validation setup
- **Dependencies**: Requires cert-validation-split-001 to be merged first

### 3. cert-validation-split-003
- **Branch**: `idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-003`
- **Location**: `efforts/phase1/wave2/cert-validation-split-003`
- **Latest Commit**: 0e545ee (marker: PROJECT-INTEGRATION Medium Bug #4 investigation complete)
- **Core Components**: Chain validator, validation errors, extended validation logic
- **Dependencies**: Requires cert-validation-split-002 to be merged first
- **Note**: Contains Bug #4 fix (extra brace removal in chain_validator_test.go)

### 4. fallback-strategies
- **Branch**: `idpbuilder-oci-build-push/phase1/wave2/fallback-strategies`
- **Location**: `efforts/phase1/wave2/fallback-strategies`
- **Latest Commit**: f3afd41 (fix(R321): fallback strategy backport analysis complete)
- **Core Components**: Certificate validation fallback logic (pkg/certvalidation)
- **Dependencies**: Can be merged independently but logically should follow cert-validation splits

## Merge Sequence and Rationale

### Phase 1: Sequential Cert-Validation Splits (MANDATORY ORDER)
These splits MUST be merged in sequential order as they build upon each other:

1. **cert-validation-split-001** (FIRST)
   - Foundation for certificate validation
   - No dependencies on other Wave 2 efforts
   - Establishes base types and utilities

2. **cert-validation-split-002** (SECOND)
   - Extends split-001 with test fixtures
   - Requires split-001's base components
   - Adds validation test setup

3. **cert-validation-split-003** (THIRD)
   - Completes validation with chain validator
   - Requires split-002's test setup
   - Includes Bug #4 fix (syntax error resolved)

### Phase 2: Fallback Strategies
4. **fallback-strategies** (FOURTH)
   - Independent implementation but logically follows cert-validation
   - Adds pkg/certvalidation directory with fallback logic
   - Should be merged after cert-validation for logical coherence

## Pre-Merge Verification Checklist

### For Each Effort:
- [ ] Branch exists and is accessible
- [ ] All commits are present as listed above
- [ ] No uncommitted changes in effort directory
- [ ] Tests pass in isolation
- [ ] Line count within limits (<800 lines per split)

## Merge Execution Instructions

### Setup Integration Workspace
```bash
# Navigate to integration workspace
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/integration-workspace

# Ensure on integration branch
git checkout idpbuilder-oci-build-push/phase1/wave2-integration

# Verify clean state
git status
```

### Step 1: Merge cert-validation-split-001
```bash
# Add split-001 remote if needed
cd ../cert-validation-split-001
git remote -v

# Back to integration workspace
cd ../integration-workspace

# Merge split-001
git merge idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-001 --no-ff \
  -m "integrate: cert-validation-split-001 into Wave 2 integration"

# Run tests to verify
go test ./pkg/certs/...

# If conflicts occur, resolve favoring split-001 changes
```

### Step 2: Merge cert-validation-split-002
```bash
# Merge split-002 (depends on split-001)
git merge idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-002 --no-ff \
  -m "integrate: cert-validation-split-002 into Wave 2 integration"

# Verify test fixtures added
ls -la pkg/certs/testdata/

# Run tests
go test ./pkg/certs/...
```

### Step 3: Merge cert-validation-split-003
```bash
# Merge split-003 (includes Bug #4 fix)
git merge idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-003 --no-ff \
  -m "integrate: cert-validation-split-003 into Wave 2 integration (includes Bug #4 fix)"

# Verify chain validator present
ls -la pkg/certs/chain_validator.go
ls -la pkg/certs/validation_errors.go

# Run validation tests
go test ./pkg/certs/... -v
```

### Step 4: Merge fallback-strategies
```bash
# Merge fallback-strategies
git merge idpbuilder-oci-build-push/phase1/wave2/fallback-strategies --no-ff \
  -m "integrate: fallback-strategies into Wave 2 integration"

# Verify fallback package
ls -la pkg/certvalidation/

# Run all tests
go test ./pkg/...
```

## Conflict Resolution Strategy

### Expected Conflicts
1. **pkg/certs/** files may have overlapping changes between splits
   - Resolution: Accept changes from the later split (they build on earlier ones)

2. **Test files** may have duplicate test cases
   - Resolution: Keep the most comprehensive test coverage

3. **pkg/certvalidation/** should be unique to fallback-strategies
   - Resolution: Accept all additions from fallback-strategies

### Conflict Resolution Principles
- Favor additions over deletions
- Maintain sequential dependency chain for cert-validation splits
- Preserve Bug #4 fix from split-003
- Ensure all test fixtures from split-002 are retained

## Post-Merge Testing Checkpoints

### After Each Merge:
1. **Compilation Check**
   ```bash
   go build ./...
   ```

2. **Unit Tests**
   ```bash
   go test ./pkg/certs/...
   go test ./pkg/certvalidation/... # After fallback merge
   ```

3. **Line Count Verification**
   ```bash
   $CLAUDE_PROJECT_DIR/tools/line-counter.sh
   ```

### Final Integration Tests:
```bash
# Full test suite
go test ./... -v

# Verify all components integrated
ls -la pkg/certs/
ls -la pkg/certvalidation/

# Check for any missing files
git status

# Final build
go build -o idpbuilder-oci ./cmd/
```

## Success Criteria

### Integration is complete when:
- ✅ All four efforts merged in specified order
- ✅ No merge conflicts remain unresolved
- ✅ All tests pass
- ✅ Binary builds successfully
- ✅ Line count remains within limits
- ✅ Bug #4 fix is present and verified
- ✅ Integration branch pushed to remote

## Rollback Plan

If integration fails at any step:
```bash
# Record failure point
echo "Failed at: [EFFORT_NAME]" > integration-failure.log
git diff > integration-conflicts.patch

# Reset to last known good state
git reset --hard HEAD~1

# Investigate and fix issue in effort branch
cd ../[effort-directory]
# Fix issues...

# Retry from failure point
```

## Notes and Considerations

1. **R327 Compliance**: This is a fresh integration following bug fixes
2. **Bug #4 Fix**: Verified present in cert-validation-split-003
3. **Sequential Dependencies**: Cert-validation splits MUST be merged in order
4. **R308 Compliance**: Integration builds on Wave 1 integration as base
5. **Testing Priority**: Each merge must pass tests before proceeding

## Recommended Execution Timeline

1. **Pre-merge verification**: 5 minutes
2. **cert-validation-split-001 merge**: 10 minutes
3. **cert-validation-split-002 merge**: 10 minutes
4. **cert-validation-split-003 merge**: 10 minutes
5. **fallback-strategies merge**: 10 minutes
6. **Final testing and validation**: 15 minutes

**Total estimated time**: 60 minutes

## Approval

This merge plan has been created following:
- R308: Incremental branching strategy
- R327: Fresh integration requirements
- R307: Independent branch mergeability
- Software Factory 2.0 integration protocols

**Ready for execution by Integration Agent**

---
*Generated by Code Reviewer Agent*
*Timestamp: 2025-09-17 01:14:02 UTC*