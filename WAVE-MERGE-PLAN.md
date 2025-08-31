# Phase 1 Wave 1 Merge Plan

## Merge Plan Metadata
- **Created**: 2025-08-31
- **Created By**: Code Reviewer Agent
- **Target Branch**: idpbuidler-oci-go-cr/phase1/wave1/integration  
- **Integration Directory**: /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/wave1/integration-workspace/idpbuilder
- **Status**: PLAN CREATED (Not Executed)

## Critical Pre-Merge Analysis

### Size Compliance Status
| Effort | Branch | Lines | Status | Notes |
|--------|--------|-------|--------|-------|
| E1.1.1 | idpbuidler-oci-go-cr/phase1/wave1/kind-certificate-extraction | 418 | ✅ COMPLIANT | Within 800 line limit |
| E1.1.2 | idpbuidler-oci-go-cr/phase1/wave1/registry-tls-trust-integration | 905 | ⚠️ EXCEEDS LIMIT | 105 lines over - but has fixes for duplicates |

### Important Findings
1. **E1.1.2 Size Issue**: Currently at 905 lines (exceeds 800 limit by 105 lines)
   - Contains duplicate fix commit (1ca4353) that resolved integration issues
   - Architecture review flagged this as CHANGES_REQUIRED
   - Despite size violation, contains critical fixes for integration

2. **No Split Branches Found**: E1.1.2 was NOT split into separate branches
   - Splits were implemented as commits within main branch
   - Split 001 commit: 449a150 (Core Trust Store Management)
   - Split 002 commit: 0abb4a5 (Transport Configuration & Utilities)

3. **Duplicate Types Fixed**: Commit 1ca4353 removed duplicate CertificateInfo struct
   - Now properly uses shared types from E1.1.1
   - Integration should be cleaner

## Merge Strategy

### Approach: Sequential Cherry-Pick Merge
Due to the size violation but critical fixes, recommended approach:

1. **Merge E1.1.1 first** (clean, compliant)
2. **Selective merge of E1.1.2** with size awareness
3. **Post-merge validation** of size and functionality

## Detailed Merge Instructions

### Phase 1: Setup and Verification
```bash
# 1. Verify current location and branch
cd /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/wave1/integration-workspace/idpbuilder
git status
# Expected: On branch idpbuidler-oci-go-cr/phase1/wave1/integration

# 2. Ensure clean working state
git status --porcelain
# Expected: Empty output (no uncommitted changes)

# 3. Fetch latest from all remotes
git fetch --all
```

### Phase 2: Merge E1.1.1 (Kind Certificate Extraction)
```bash
# 1. Add E1.1.1 directory as remote (if not already added)
git remote add e111 /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/wave1/kind-certificate-extraction
git fetch e111

# 2. Merge E1.1.1 branch
git merge e111/idpbuidler-oci-go-cr/phase1/wave1/kind-certificate-extraction --no-ff \
  -m "feat: integrate E1.1.1 - Kind Certificate Extraction (418 lines)"

# 3. Expected outcome:
# - New files in pkg/kind/ for certificate extraction
# - CertificateInfo type in pkg/certs/types.go
# - No conflicts expected (new functionality)
```

### Phase 3: Merge E1.1.2 (Registry TLS Trust Integration)
```bash
# 1. Add E1.1.2 directory as remote (if not already added)
git remote add e112 /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/wave1/registry-tls-trust-integration
git fetch e112

# 2. Merge E1.1.2 branch (with size awareness)
git merge e112/idpbuidler-oci-go-cr/phase1/wave1/registry-tls-trust-integration --no-ff \
  -m "feat: integrate E1.1.2 - Registry TLS Trust Integration (905 lines - includes dedup fixes)"

# 3. Expected conflicts/issues:
# - CertificateInfo type might conflict (already resolved in 1ca4353)
# - Test helper functions might have duplicates
# - Resolution: Keep E1.1.2 versions as they have the fixes
```

### Phase 4: Post-Merge Validation
```bash
# 1. Verify no duplicate types
grep -r "type CertificateInfo struct" pkg/
# Expected: Only one occurrence

# 2. Check total size
$PROJECT_ROOT/tools/line-counter.sh
# Expected: Combined ~1323 lines (418 + 905)

# 3. Build verification
go build ./...
# Expected: Successful build

# 4. Run tests
go test ./pkg/... -v
# Expected: All tests pass
```

## Conflict Resolution Guidelines

### Expected Conflicts
1. **CertificateInfo struct**: Already resolved in E1.1.2 commit 1ca4353
   - Keep E1.1.2 version (has deduplication fix)
   
2. **Test Helper Functions**: createTestCertificate might appear twice
   - Keep E1.1.2 version (updated to work with shared types)

3. **Import statements**: Both efforts might add similar imports
   - Merge and deduplicate imports

### Resolution Strategy
```bash
# If conflicts occur:
# 1. Review conflict markers
git status

# 2. For type conflicts, keep E1.1.2 version:
git checkout --theirs pkg/certs/trust_store.go

# 3. For test conflicts, keep E1.1.2 version:
git checkout --theirs pkg/certs/trust_test.go

# 4. Manually review and fix imports
# 5. Complete merge
git add .
git commit
```

## Critical Reminders

### Size Violation Handling
⚠️ **WARNING**: E1.1.2 exceeds 800 line limit (905 lines)

**Options for Integration Agent**:
1. **Proceed with caution**: Merge as-is but flag for immediate post-integration split
2. **Request split first**: Do NOT merge, request E1.1.2 be properly split before integration
3. **Selective merge**: Cherry-pick only commits up to 800 lines

**Recommended**: Option 1 - Proceed but flag, as the deduplication fixes are critical for integration success.

### Verification Checklist
After merge, verify:
- [ ] Build succeeds: `go build ./...`
- [ ] Tests pass: `go test ./pkg/...`
- [ ] No duplicate types: Check CertificateInfo
- [ ] Size documented: Run line-counter.sh
- [ ] Commit message includes size warning if >800

## Merge Order Summary

1. **E1.1.1** (kind-certificate-extraction) - 418 lines ✅
2. **E1.1.2** (registry-tls-trust-integration) - 905 lines ⚠️

## Final Notes

### For Integration Agent
- This plan is created per R269 requirements
- Do NOT execute merges as Code Reviewer
- Size violation in E1.1.2 needs acknowledgment
- Consider requesting proper split of E1.1.2 before merge
- If proceeding, document size violation in commit messages

### Risk Assessment
- **Low Risk**: E1.1.1 merge (clean, compliant)
- **Medium Risk**: E1.1.2 merge (size violation but has critical fixes)
- **Mitigation**: Can proceed but must flag for post-integration split

### Success Criteria
- Both efforts merged without breaking builds
- All tests pass post-merge
- No duplicate type definitions
- Clear documentation of size violations
- Integration branch ready for testing

---
*End of Merge Plan - Created by Code Reviewer Agent*
*Integration Agent should execute this plan with awareness of size violations*