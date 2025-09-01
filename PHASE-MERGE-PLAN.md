# Phase 1 Integration Merge Plan

## 🚨 CRITICAL MERGE INSTRUCTIONS FOR INTEGRATION AGENT 🚨

**Phase**: 1 - Certificate Infrastructure  
**Total Waves**: 2  
**Target Branch**: `idpbuidler-oci-go-cr/phase1/integration`  
**Base Branch**: `software-factory-2.0`  
**Created By**: Code Reviewer Agent  
**Created At**: 2025-09-01  
**Plan Type**: NORMAL FLOW (from WAVE_REVIEW state)

## 📋 Executive Summary

This plan provides exact instructions for merging all Phase 1 effort branches into the phase integration branch. The Integration Agent must execute these commands in the EXACT order specified to ensure proper dependency resolution and avoid conflicts.

## 🎯 Branches to Merge (In Order)

### Wave 1 Branches (Foundational)
1. **E1.1.1**: `idpbuidler-oci-go-cr/phase1/wave1/kind-certificate-extraction`
   - Location: `/home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/wave1/kind-certificate-extraction`
   - Dependencies: None (foundational)
   
2. **E1.1.2**: `idpbuidler-oci-go-cr/phase1/wave1/registry-tls-trust-integration`
   - Location: `/home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/wave1/registry-tls-trust-integration`
   - Dependencies: None (foundational, parallel with E1.1.1)

### Wave 2 Branches (Dependent)
3. **E1.2.1**: `idpbuidler-oci-go-cr/phase1/wave2/certificate-validation-pipeline`
   - Location: `/home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/wave2/certificate-validation-pipeline`
   - Dependencies: E1.1.1, E1.1.2 (requires Wave 1 foundation)
   
4. **E1.2.2**: `idpbuidler-oci-go-cr/phase1/wave2/fallback-strategies`
   - Location: `/home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/wave2/fallback-strategies`
   - Dependencies: E1.1.1, E1.1.2 (requires Wave 1 foundation)

## 🔧 Pre-Merge Setup Commands

```bash
# 1. Navigate to phase integration workspace
cd /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/phase-integration-workspace

# 2. Verify we're on the correct integration branch
git branch --show-current
# Expected: idpbuidler-oci-go-cr/phase1/integration

# 3. Ensure clean working tree
git status --porcelain
# Should be empty. If not, stash or commit changes

# 4. Fetch all updates
git fetch origin

# 5. Ensure integration branch is up to date with base
git merge origin/software-factory-2.0 --no-edit
# This ensures we have the latest base changes
```

## 📦 Merge Execution Commands

### Step 1: Merge Wave 1 Effort 1 (Kind Certificate Extraction)

```bash
# Navigate to integration workspace
cd /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/phase-integration-workspace

# Add effort directory as remote (if not already added)
git remote add kind-cert-extraction ../wave1/kind-certificate-extraction/.git 2>/dev/null || true
git fetch kind-cert-extraction

# Merge the effort branch
git merge kind-cert-extraction/idpbuidler-oci-go-cr/phase1/wave1/kind-certificate-extraction \
  --no-ff \
  -m "merge: integrate E1.1.1 - Kind Certificate Extraction into Phase 1 integration"

# Verify successful merge
git status
```

### Step 2: Merge Wave 1 Effort 2 (Registry TLS Trust Integration)

```bash
# Add effort directory as remote
git remote add registry-tls ../wave1/registry-tls-trust-integration/.git 2>/dev/null || true
git fetch registry-tls

# Merge the effort branch
git merge registry-tls/idpbuidler-oci-go-cr/phase1/wave1/registry-tls-trust-integration \
  --no-ff \
  -m "merge: integrate E1.1.2 - Registry TLS Trust Integration into Phase 1 integration"

# Verify successful merge
git status
```

### Step 3: Merge Wave 2 Effort 1 (Certificate Validation Pipeline)

```bash
# Add effort directory as remote
git remote add cert-validation ../wave2/certificate-validation-pipeline/.git 2>/dev/null || true
git fetch cert-validation

# Merge the effort branch
git merge cert-validation/idpbuidler-oci-go-cr/phase1/wave2/certificate-validation-pipeline \
  --no-ff \
  -m "merge: integrate E1.2.1 - Certificate Validation Pipeline into Phase 1 integration"

# Verify successful merge
git status
```

### Step 4: Merge Wave 2 Effort 2 (Fallback Strategies)

```bash
# Add effort directory as remote
git remote add fallback-strat ../wave2/fallback-strategies/.git 2>/dev/null || true
git fetch fallback-strat

# Merge the effort branch
git merge fallback-strat/idpbuidler-oci-go-cr/phase1/wave2/fallback-strategies \
  --no-ff \
  -m "merge: integrate E1.2.2 - Fallback Strategies into Phase 1 integration"

# Verify successful merge
git status
```

## 🔍 Conflict Resolution Strategy

### Expected Conflicts

Based on the parallel development structure, conflicts are likely in:

1. **go.mod / go.sum**: Multiple efforts adding dependencies
   - **Resolution**: Accept both sets of dependencies, then run `go mod tidy`

2. **pkg/ directory structure**: New packages from different efforts
   - **Resolution**: Accept all new packages (no actual conflict, just additions)

3. **Import statements**: Different efforts importing different packages
   - **Resolution**: Combine all imports, remove duplicates

### Conflict Resolution Commands

```bash
# If conflicts occur during any merge:

# 1. Check conflict status
git status

# 2. For go.mod/go.sum conflicts:
# Edit go.mod to include all dependencies
git add go.mod go.sum
go mod tidy
git add go.mod go.sum

# 3. For source file conflicts:
# Open conflicted files and resolve manually
# Look for <<<<<<< HEAD markers
# Keep both changes when they're additive
# Ensure no duplicate imports or functions

# 4. After resolving all conflicts:
git add -A
git commit --no-edit

# 5. Verify build still works:
go build ./...
```

## ✅ Post-Merge Validation Steps

### Required Validation (Must Pass)

```bash
# 1. Verify all tests pass
go test ./... -v

# 2. Verify build succeeds
go build ./...

# 3. Verify no linting errors
golangci-lint run ./... || true  # Warning only, not blocking

# 4. Check for any missing files
git status --porcelain

# 5. Verify all effort code is present
ls -la pkg/certificates/
ls -la pkg/registry/
ls -la pkg/validation/
ls -la pkg/fallback/

# 6. Run integration tests if present
if [ -d "tests/integration" ]; then
    go test ./tests/integration/... -v
fi
```

### Line Count Verification

```bash
# Measure total phase size
${CLAUDE_PROJECT_DIR}/tools/line-counter.sh

# Verify no individual effort exceeded limits
# (Should already be verified, but double-check)
echo "Phase 1 Total Size Check:"
echo "Expected: < 3200 lines (4 efforts * 800 max)"
```

## 🚨 Critical Success Criteria

The Integration Agent MUST ensure:

1. ✅ All 4 effort branches successfully merged
2. ✅ No uncommitted changes remain
3. ✅ All tests pass after merge
4. ✅ Build succeeds without errors
5. ✅ Total phase size is within limits
6. ✅ Integration branch pushed to origin

## 📝 Final Integration Commands

```bash
# After all merges and validations complete:

# 1. Final status check
git status
git log --oneline -10

# 2. Push integration branch
git push origin idpbuidler-oci-go-cr/phase1/integration

# 3. Create integration summary
cat > PHASE-1-INTEGRATION-SUMMARY.md << EOF
# Phase 1 Integration Summary

**Integration Date**: $(date '+%Y-%m-%d %H:%M:%S')
**Integration Agent**: [Agent ID]
**Target Branch**: idpbuidler-oci-go-cr/phase1/integration

## Merged Efforts
- ✅ E1.1.1: Kind Certificate Extraction
- ✅ E1.1.2: Registry TLS Trust Integration  
- ✅ E1.2.1: Certificate Validation Pipeline
- ✅ E1.2.2: Fallback Strategies

## Validation Results
- Tests: PASSED
- Build: SUCCESSFUL
- Line Count: [ACTUAL] lines (Limit: 3200)

## Next Steps
- Ready for Phase 2 implementation
- Phase 1 functionality complete and integrated
EOF

git add PHASE-1-INTEGRATION-SUMMARY.md
git commit -m "doc: Phase 1 integration completed successfully"
git push origin idpbuidler-oci-go-cr/phase1/integration
```

## ⚠️ Error Recovery

If any merge fails:

1. **DO NOT PROCEED** with subsequent merges
2. **Document the failure** in MERGE-ERROR-REPORT.md
3. **Report back** to orchestrator with:
   - Which merge failed
   - Exact error message
   - Conflict details if applicable
4. **Await instructions** before continuing

## 📋 Integration Agent Checklist

Before starting:
- [ ] In phase-integration-workspace directory
- [ ] On idpbuidler-oci-go-cr/phase1/integration branch
- [ ] Working tree is clean
- [ ] All remotes fetched

During merges:
- [ ] E1.1.1 merged successfully
- [ ] E1.1.2 merged successfully
- [ ] E1.2.1 merged successfully
- [ ] E1.2.2 merged successfully
- [ ] All conflicts resolved (if any)

After merges:
- [ ] All tests pass
- [ ] Build succeeds
- [ ] No uncommitted changes
- [ ] Integration branch pushed
- [ ] Summary report created

---

**END OF PHASE MERGE PLAN**

*This plan created by Code Reviewer Agent for execution by Integration Agent*
*DO NOT MODIFY - Execute exactly as specified*