# PHASE 1 INTEGRATION MERGE PLAN

**Created**: 2025-01-12  
**State**: PHASE_INTEGRATION  
**Integration Branch**: `idpbuilder-oci-build-push/phase1/integration`  
**Integration Directory**: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/phase-integration-workspace`  
**Plan Type**: R269 Compliant - Plan Only (NO EXECUTION)

## CRITICAL NOTICE
This plan is created by Code Reviewer per R269. The Integration Agent will execute these commands.
DO NOT execute any merges manually - this is a PLAN ONLY document.

## 1. MERGE HIERARCHY

Per Software Factory 2.0 integration hierarchy:
```
Individual Efforts → Wave Integration → Phase Integration → Main
```

This plan merges WAVE integration branches into the PHASE integration branch.

## 2. BRANCHES TO INTEGRATE

### Primary Integration Sources (Wave Integrations):
1. **Wave 1 Integration**: `idpbuilder-oci-build-push/phase1/wave1/integration`
   - Contains: E1.1.1 (kind-cert-extraction), E1.1.2 (registry-tls-trust), E1.1.3 splits
   - Total Lines: ~2,950 lines
   - Status: Integrated and reviewed

2. **Wave 2 Integration**: `idpbuilder-oci-build-push/phase1/wave2/integration`
   - Contains: E1.2.1 splits (cert-validation), E1.2.2 (fallback-strategies)
   - Total Lines: ~2,367 lines
   - Status: Integrated and reviewed

## 3. PRE-MERGE CHECKLIST

The Integration Agent MUST verify before starting:
```bash
# Verify clean workspace
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/phase-integration-workspace
git status --porcelain  # Must be empty

# Verify correct branch
git branch --show-current  # Must show: idpbuilder-oci-build-push/phase1/integration

# Verify latest main
git fetch origin main
git log --oneline origin/main -1  # Record commit hash for reference

# Verify wave branches exist
git ls-remote origin | grep "phase1/wave1/integration"
git ls-remote origin | grep "phase1/wave2/integration"
```

## 4. MERGE SEQUENCE

### Step 1: Prepare Integration Branch
```bash
# Ensure we're starting from latest main
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/phase-integration-workspace
git checkout idpbuilder-oci-build-push/phase1/integration
git fetch origin main
git merge origin/main --no-edit  # Get latest main changes

# Push the updated base
git push origin idpbuilder-oci-build-push/phase1/integration
```

### Step 2: Merge Wave 1 Integration
```bash
# Fetch wave 1 integration
git fetch origin idpbuilder-oci-build-push/phase1/wave1/integration

# Merge wave 1 (should be clean, no conflicts expected)
git merge origin/idpbuilder-oci-build-push/phase1/wave1/integration \
  --no-ff \
  -m "feat(phase1): integrate Wave 1 - certificate extraction and registry trust

Integrates:
- E1.1.1: Kind certificate extraction (650 lines)
- E1.1.2: Registry TLS trust setup (700 lines)
- E1.1.3: Registry auth types (1,600 lines across 2 splits)

Total: ~2,950 lines of foundational certificate handling"

# Verify successful merge
if [ $? -ne 0 ]; then
    echo "ERROR: Wave 1 merge failed - see conflict resolution below"
    exit 1
fi
```

### Step 3: Test Wave 1 Integration
```bash
# Build test
make build || go build ./...

# Run unit tests
make test || go test ./...

# Verify no broken imports
go mod tidy
go mod verify

# Check for compilation issues
go vet ./...
```

### Step 4: Merge Wave 2 Integration
```bash
# Fetch wave 2 integration
git fetch origin idpbuilder-oci-build-push/phase1/wave2/integration

# Merge wave 2 (may have conflicts with wave 1)
git merge origin/idpbuilder-oci-build-push/phase1/wave2/integration \
  --no-ff \
  -m "feat(phase1): integrate Wave 2 - certificate validation and fallback

Integrates:
- E1.2.1: Certificate validation logic (1,807 lines across 3 splits)
- E1.2.2: Fallback strategies (560 lines)

Total: ~2,367 lines of validation and error handling"

# Check merge status
if [ $? -ne 0 ]; then
    echo "CONFLICTS DETECTED - proceed to conflict resolution"
fi
```

### Step 5: Test Complete Phase Integration
```bash
# Full build
make clean && make build

# Comprehensive tests
make test-all || go test -v ./...

# Integration tests if available
if [ -f "test/integration_test.go" ]; then
    go test ./test/... -tags=integration
fi

# Final verification
go mod tidy
git status --porcelain  # Should only show go.mod/go.sum if changed
```

## 5. CONFLICT RESOLUTION STRATEGIES

### Expected Conflict Points:
1. **go.mod/go.sum**: Wave 2 may add dependencies that conflict
   - Resolution: Accept both sets of dependencies, run `go mod tidy`

2. **pkg/registry/auth.go**: Both waves modify authentication
   - Resolution: Merge both changes, Wave 2 validation builds on Wave 1 types

3. **pkg/cert/validation.go**: Wave 2 extends Wave 1 structures
   - Resolution: Keep Wave 1 base types, add Wave 2 validation methods

### Conflict Resolution Commands:
```bash
# For go.mod conflicts
git checkout --theirs go.mod go.sum
go mod tidy
git add go.mod go.sum

# For source code conflicts
# Edit conflicted files to merge both changes
vim pkg/registry/auth.go  # or use preferred editor
# Resolve markers (<<<<, ====, >>>>)
git add pkg/registry/auth.go

# Verify resolution
go build ./...
go test ./...

# Complete merge
git commit --no-edit
```

## 6. ROLLBACK PROCEDURES

### If Wave 1 Merge Fails:
```bash
# Abort merge and reset
git merge --abort
git reset --hard origin/idpbuilder-oci-build-push/phase1/integration
git clean -fd

# Investigate issue
git log --oneline origin/idpbuilder-oci-build-push/phase1/wave1/integration -10
# Report to orchestrator for investigation
```

### If Wave 2 Merge Fails:
```bash
# If conflicts are unresolvable
git merge --abort

# Reset to after Wave 1 merge
git reset --hard HEAD~1  # Only if Wave 1 was successful

# Alternative: manual cherry-pick approach
git cherry-pick origin/idpbuilder-oci-build-push/phase1/wave2/integration
```

### Complete Rollback:
```bash
# Nuclear option - start over
git checkout idpbuilder-oci-build-push/phase1/integration
git reset --hard origin/main
git push --force-with-lease origin idpbuilder-oci-build-push/phase1/integration
# Then restart from Step 1
```

## 7. FINAL VALIDATION

After successful integration:
```bash
# 1. Verify all features present
grep -r "ExtractKindCertificate" pkg/  # E1.1.1
grep -r "ConfigureRegistryTLS" pkg/     # E1.1.2
grep -r "RegistryAuthConfig" pkg/       # E1.1.3
grep -r "ValidateCertificate" pkg/      # E1.2.1
grep -r "FallbackStrategy" pkg/         # E1.2.2

# 2. Line count verification
PROJECT_ROOT="/home/vscode/workspaces/idpbuilder-oci-build-push"
$PROJECT_ROOT/tools/line-counter.sh
# Expected: ~5,317 lines total (2,950 + 2,367)

# 3. Test coverage
go test -cover ./...
# Minimum expected: 70% coverage

# 4. No broken dependencies
go mod why -m all | head -20
go list -m all | wc -l  # Record dependency count

# 5. Push integrated branch
git push origin idpbuilder-oci-build-push/phase1/integration
```

## 8. SUCCESS CRITERIA

The phase integration is successful when:
- ✅ Both wave integrations merged cleanly (or conflicts resolved)
- ✅ All tests pass (unit, integration)
- ✅ Build completes without errors
- ✅ Line count is within expected range (~5,317 lines ±10%)
- ✅ All effort features are present and functional
- ✅ No regression in existing functionality
- ✅ Branch pushed to remote successfully

## 9. POST-INTEGRATION ACTIONS

After successful integration, the Integration Agent should:
1. Create PHASE-1-INTEGRATION-REPORT.md with results
2. Update orchestrator-state.json to mark phase complete
3. Tag the integration: `git tag phase1-integrated`
4. Notify orchestrator of completion
5. Prepare for Architect review

## 10. INTEGRATION AGENT EXECUTION NOTES

### Execution Command:
```bash
# Integration Agent should execute with:
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/phase-integration-workspace
bash -e  # Exit on any error
# Then copy-paste commands from sections above
```

### Error Handling:
- Any non-zero exit code should stop execution
- Report specific failure point to orchestrator
- Save git status and logs for debugging
- Do NOT attempt automatic fixes beyond documented strategies

### Time Estimate:
- Wave 1 Merge: 2-3 minutes
- Wave 1 Tests: 3-5 minutes  
- Wave 2 Merge: 2-3 minutes (or 5-10 if conflicts)
- Wave 2 Tests: 3-5 minutes
- Final Validation: 5 minutes
- **Total**: 15-30 minutes

## APPROVAL SIGNATURE

This plan is created per R269 requirements and is ready for execution by the Integration Agent.

**Plan Status**: READY FOR EXECUTION
**Created By**: Code Reviewer Agent
**Rule Compliance**: R269, R285, R307 (trunk-based mergeability)