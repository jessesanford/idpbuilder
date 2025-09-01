# Phase 1 Wave 2 Merge Plan

## Overview
**Date Created**: 2025-09-01  
**Integration Branch**: `idpbuidler-oci-go-cr/phase1/wave2/integration`  
**Target Repository**: https://github.com/jessesanford/idpbuilder.git  
**Total Efforts**: 2  
**Total Lines**: 1,175 (431 + 744)  

## Efforts to Integrate

| Effort | Branch | Lines | Status | Base Commit |
|--------|--------|-------|--------|-------------|
| E1.2.1 | `idpbuidler-oci-go-cr/phase1/wave2/certificate-validation-pipeline` | 431 | COMPLETED & PASSED | a12268f |
| E1.2.2 | `idpbuidler-oci-go-cr/phase1/wave2/fallback-strategies` | 744 | COMPLETED & PASSED | e210954 |

## Pre-Merge Verification

### 1. Verify Integration Branch State
```bash
# Ensure we're on the correct integration branch
git checkout idpbuidler-oci-go-cr/phase1/wave2/integration
git pull origin idpbuidler-oci-go-cr/phase1/wave2/integration

# Verify clean working directory
git status
# Expected: "nothing to commit, working tree clean"

# Record starting commit
git log -1 --oneline
# Current HEAD: e210954 todo(architect): WAVE_REVIEW checkpoint at state WAVE_REVIEW [R287]
```

### 2. Fetch All Effort Branches
```bash
# Fetch the latest changes from all effort branches
git fetch origin idpbuidler-oci-go-cr/phase1/wave2/certificate-validation-pipeline:refs/remotes/origin/idpbuidler-oci-go-cr/phase1/wave2/certificate-validation-pipeline
git fetch origin idpbuidler-oci-go-cr/phase1/wave2/fallback-strategies:refs/remotes/origin/idpbuidler-oci-go-cr/phase1/wave2/fallback-strategies

# Verify branches are available
git branch -r | grep "idpbuidler-oci-go-cr/phase1/wave2"
```

## Merge Order Analysis

### Dependency Analysis:
1. **E1.2.1 (certificate-validation-pipeline)**: 
   - Base: a12268f (older)
   - Files: pkg/certs/{validator.go, diagnostics.go, testdata/, tests}
   - No dependencies on E1.2.2

2. **E1.2.2 (fallback-strategies)**:
   - Base: e210954 (newer, includes E1.2.1's base)
   - Files: pkg/fallback/*, pkg/certs/types.go
   - May depend on certificate validation interfaces

### Recommended Merge Order:
1. **FIRST**: Merge E1.2.1 (certificate-validation-pipeline)
2. **SECOND**: Merge E1.2.2 (fallback-strategies)

**Rationale**: 
- E1.2.1 has an older base commit, merging it first ensures a cleaner history
- E1.2.2's base (e210954) is already ahead and includes E1.2.1's base (a12268f)
- No file conflicts detected between efforts (different package directories)

## Merge Execution Instructions

### Step 1: Merge E1.2.1 (certificate-validation-pipeline)
```bash
# 1.1 Start merge
git merge origin/idpbuidler-oci-go-cr/phase1/wave2/certificate-validation-pipeline \
    --no-ff \
    -m "merge: integrate E1.2.1 certificate-validation-pipeline (431 lines)"

# 1.2 If conflicts occur (unlikely based on file analysis):
# - Resolve conflicts favoring the effort's implementation
# - Stage resolved files: git add <resolved-files>
# - Complete merge: git commit

# 1.3 Verify merge success
git status
git diff HEAD~1..HEAD --stat

# 1.4 Expected files added:
# - pkg/certs/validator.go
# - pkg/certs/validator_test.go
# - pkg/certs/diagnostics.go
# - pkg/certs/testdata/certs.go

# 1.5 Run tests for E1.2.1
go test ./pkg/certs/... -v
```

### Step 2: Merge E1.2.2 (fallback-strategies)
```bash
# 2.1 Start merge
git merge origin/idpbuidler-oci-go-cr/phase1/wave2/fallback-strategies \
    --no-ff \
    -m "merge: integrate E1.2.2 fallback-strategies (744 lines)"

# 2.2 If conflicts occur:
# - Most likely in documentation files (IMPLEMENTATION-PLAN.md, work-log.md)
# - For documentation conflicts: Keep both sections or merge meaningfully
# - For code conflicts: Favor the effort's implementation
# - Stage resolved files: git add <resolved-files>
# - Complete merge: git commit

# 2.3 Verify merge success
git status
git diff HEAD~1..HEAD --stat

# 2.4 Expected files added:
# - pkg/certs/types.go
# - pkg/fallback/detector.go
# - pkg/fallback/detector_test.go
# - pkg/fallback/insecure.go
# - pkg/fallback/insecure_test.go
# - pkg/fallback/logger.go
# - pkg/fallback/recommender.go
# - pkg/fallback/recommender_test.go

# 2.5 Run tests for E1.2.2
go test ./pkg/fallback/... -v
```

## Conflict Resolution Strategy

### Expected Conflicts:
1. **Documentation Files** (IMPLEMENTATION-PLAN.md, work-log.md):
   - Resolution: These are effort-specific and can be removed or archived
   - Command: `git rm IMPLEMENTATION-PLAN.md work-log.md` if conflicts occur
   
2. **TODO Files** (todos/*.todo):
   - Resolution: Keep all TODO files for audit trail
   - Command: `git add todos/*.todo` to accept all

### Unexpected Code Conflicts:
If code conflicts occur (unlikely based on analysis):

```bash
# 1. Identify conflict files
git status | grep "both modified"

# 2. For each conflicted file:
# - Open in editor
# - Look for <<<<<<< HEAD markers
# - Keep the effort's implementation (incoming changes)
# - Remove conflict markers
# - Save file

# 3. Stage resolved files
git add <resolved-file>

# 4. Complete merge
git commit
```

## Post-Merge Verification

### 1. Build Verification
```bash
# Clean build
go clean -cache
go mod tidy
go build ./...

# Expected: Build succeeds with no errors
```

### 2. Test Suite Execution
```bash
# Run all tests
go test ./... -v -cover

# Specific package tests
go test ./pkg/certs/... -v -cover
go test ./pkg/fallback/... -v -cover

# Expected: All tests pass
```

### 3. Line Count Verification
```bash
# Use project line counter if available
if [ -f "$CLAUDE_PROJECT_DIR/tools/line-counter.sh" ]; then
    $CLAUDE_PROJECT_DIR/tools/line-counter.sh
else
    # Fallback to git diff
    git diff origin/idpbuidler-oci-go-cr/phase1/wave2/integration..HEAD --stat
fi

# Expected total: ~1,175 lines added
```

### 4. Integration Validation
```bash
# Verify both efforts are integrated
git log --oneline -10 | grep -E "(certificate-validation|fallback-strategies)"

# Check file structure
ls -la pkg/certs/
ls -la pkg/fallback/

# Verify no missing files
git diff --name-only origin/idpbuidler-oci-go-cr/phase1/wave2/certificate-validation-pipeline..HEAD | grep "^pkg/"
git diff --name-only origin/idpbuidler-oci-go-cr/phase1/wave2/fallback-strategies..HEAD | grep "^pkg/"
```

## Rollback Plan

If critical issues occur during merge:

```bash
# Record current commit before starting
PRE_MERGE_COMMIT=$(git rev-parse HEAD)

# If rollback needed:
git reset --hard $PRE_MERGE_COMMIT
git clean -fd

# Re-fetch branches and retry with different strategy
```

## Final Steps

### 1. Push Integration Branch
```bash
# After successful merge and verification
git push origin idpbuidler-oci-go-cr/phase1/wave2/integration

# Verify push
git status
# Expected: "Your branch is up to date with 'origin/idpbuidler-oci-go-cr/phase1/wave2/integration'"
```

### 2. Create Integration Report
Create WAVE-INTEGRATION-REPORT.md with:
- Merge timestamps
- Final line counts
- Test results
- Any issues resolved
- Next steps for Phase 1 Wave 3

### 3. Tag Integration Point (Optional)
```bash
git tag -a "phase1-wave2-integrated" -m "Phase 1 Wave 2 Integration Complete: E1.2.1 + E1.2.2"
git push origin --tags
```

## Notes for Integration Agent

1. **DO NOT** delete effort branches after merge - keep for audit trail
2. **DO NOT** squash commits - preserve full history with --no-ff
3. **DO** run tests after each merge, not just at the end
4. **DO** document any deviations from this plan
5. **DO** create backup branch before starting: 
   ```bash
   git branch backup-pre-wave2-integration
   ```

## Success Criteria

- [ ] Both effort branches merged successfully
- [ ] Zero merge conflicts OR all conflicts resolved properly
- [ ] All tests passing (100% pass rate)
- [ ] Build succeeds without warnings
- [ ] Total line count matches expected (~1,175 lines)
- [ ] Integration branch pushed to remote
- [ ] No files lost or corrupted
- [ ] Clean git status after completion

## Contact for Issues

If critical blockers encountered:
1. Document the exact error and state
2. Create INTEGRATION-BLOCKER.md with details
3. Rollback to pre-merge state
4. Await orchestrator intervention

---
*Generated by Code Reviewer Agent*  
*Date: 2025-09-01*  
*Phase: 1, Wave: 2*