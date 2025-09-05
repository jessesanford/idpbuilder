# Phase 1 Integration Merge Plan

## Overview
**Integration Branch**: `idpbuilder-oci-go-cr/phase1-integration-20250902-194557`  
**Base Branch**: `main`  
**Total Efforts**: 4 (2 from Wave 1, 2 from Wave 2)  
**Created**: 2025-01-02  
**Merge Strategy**: Sequential by wave, preserving dependency order  

## Pre-Merge Verification Checklist

### Environment Setup
```bash
# 1. Navigate to integration workspace
cd /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/phase-integration-workspace

# 2. Ensure on integration branch
git checkout idpbuilder-oci-go-cr/phase1-integration-20250902-194557

# 3. Verify clean working directory
git status --porcelain
# Expected output: empty (no uncommitted changes)

# 4. Fetch all remote branches
git fetch --all

# 5. Verify integration branch is up-to-date with main
git merge-base HEAD main
```

### Branch Verification
```bash
# Verify all effort branches exist and are accessible
for branch in \
  "idpbuidler-oci-go-cr/phase1/wave1/kind-certificate-extraction" \
  "idpbuidler-oci-go-cr/phase1/wave1/registry-tls-trust-integration" \
  "idpbuidler-oci-go-cr/phase1/wave2/certificate-validation-pipeline" \
  "idpbuidler-oci-go-cr/phase1/wave2/fallback-strategies"; do
  echo "Checking branch: $branch"
  git rev-parse $branch || echo "ERROR: Branch $branch not found"
done
```

## Merge Order Strategy

### Phase 1, Wave 1 Merges (Execute First)

#### 1. kind-certificate-extraction (418 lines)
**Branch**: `idpbuidler-oci-go-cr/phase1/wave1/kind-certificate-extraction`  
**Dependencies**: None (foundational effort)  
**Expected Conflicts**: None (first merge)  

```bash
# Merge Wave 1, Effort 1: kind-certificate-extraction
echo "==== Merging kind-certificate-extraction ===="
git merge --no-ff -m "merge(phase1/wave1): integrate kind-certificate-extraction effort

- Extract certificates from Kind cluster nodes
- Implement cert retrieval logic via docker/kubectl
- Add foundational certificate handling interfaces
- Lines added: 418

[R269 Phase Integration Protocol]" \
  idpbuidler-oci-go-cr/phase1/wave1/kind-certificate-extraction

# Verify merge success
if [ $? -eq 0 ]; then
  echo "✅ Successfully merged kind-certificate-extraction"
  git log --oneline -1
else
  echo "❌ Merge failed - resolve conflicts and continue"
fi
```

#### 2. registry-tls-trust-integration (936 lines - 2 splits)
**Branch**: `idpbuidler-oci-go-cr/phase1/wave1/registry-tls-trust-integration`  
**Dependencies**: kind-certificate-extraction (uses certificate interfaces)  
**Expected Conflicts**: Possible in shared interfaces  

```bash
# Merge Wave 1, Effort 2: registry-tls-trust-integration
echo "==== Merging registry-tls-trust-integration ===="
git merge --no-ff -m "merge(phase1/wave1): integrate registry-tls-trust-integration effort

- Integrate TLS trust configuration with registries
- Implement trust store management
- Connect certificate extraction to registry clients
- Lines added: 936 (completed in 2 splits)

[R269 Phase Integration Protocol]" \
  idpbuidler-oci-go-cr/phase1/wave1/registry-tls-trust-integration

# Verify merge success
if [ $? -eq 0 ]; then
  echo "✅ Successfully merged registry-tls-trust-integration"
  git log --oneline -1
else
  echo "❌ Merge failed - resolve conflicts and continue"
fi
```

### Wave 1 Integration Verification
```bash
# After Wave 1 merges, verify integration
echo "==== Wave 1 Integration Verification ===="

# 1. Check all Wave 1 code is present
find pkg/ -type f -name "*.go" | xargs grep -l "wave1" | wc -l

# 2. Run Wave 1 tests
go test ./pkg/certificate/... ./pkg/trust/... -v

# 3. Verify compilation
go build ./...

# 4. Create Wave 1 integration checkpoint
git tag -a "phase1-wave1-integrated" -m "Wave 1 integration complete"
```

### Phase 1, Wave 2 Merges (Execute After Wave 1)

#### 3. certificate-validation-pipeline (431 lines)
**Branch**: `idpbuidler-oci-go-cr/phase1/wave2/certificate-validation-pipeline`  
**Dependencies**: Wave 1 efforts (builds on trust store and certificates)  
**Expected Conflicts**: Possible in validation interfaces  

```bash
# Merge Wave 2, Effort 1: certificate-validation-pipeline
echo "==== Merging certificate-validation-pipeline ===="
git merge --no-ff -m "merge(phase1/wave2): integrate certificate-validation-pipeline effort

- Implement certificate validation pipeline
- Add chain validation and expiry checks
- Integrate with Wave 1 trust store infrastructure
- Lines added: 431

Dependencies: kind-certificate-extraction, registry-tls-trust-integration
[R269 Phase Integration Protocol]" \
  idpbuidler-oci-go-cr/phase1/wave2/certificate-validation-pipeline

# Verify merge success
if [ $? -eq 0 ]; then
  echo "✅ Successfully merged certificate-validation-pipeline"
  git log --oneline -1
else
  echo "❌ Merge failed - resolve conflicts and continue"
fi
```

#### 4. fallback-strategies (658 lines)
**Branch**: `idpbuidler-oci-go-cr/phase1/wave2/fallback-strategies`  
**Dependencies**: All previous efforts (provides fallback for entire system)  
**Expected Conflicts**: Minimal (mostly additive)  

```bash
# Merge Wave 2, Effort 2: fallback-strategies
echo "==== Merging fallback-strategies ===="
git merge --no-ff -m "merge(phase1/wave2): integrate fallback-strategies effort

- Implement fallback mechanisms for certificate issues
- Add retry logic and alternative validation paths
- Provide graceful degradation options
- Lines added: 658

Dependencies: All Wave 1 efforts, certificate-validation-pipeline
[R269 Phase Integration Protocol]" \
  idpbuidler-oci-go-cr/phase1/wave2/fallback-strategies

# Verify merge success
if [ $? -eq 0 ]; then
  echo "✅ Successfully merged fallback-strategies"
  git log --oneline -1
else
  echo "❌ Merge failed - resolve conflicts and continue"
fi
```

## Conflict Resolution Strategy

### Expected Conflict Points
1. **pkg/certificate/interfaces.go** - Multiple efforts may modify interfaces
2. **pkg/trust/manager.go** - Trust store management touched by multiple efforts
3. **go.mod / go.sum** - Dependency conflicts from different efforts
4. **Test files** - Overlapping test utilities

### Conflict Resolution Process
```bash
# When conflicts occur:

# 1. Identify conflicted files
git status --porcelain | grep "^UU"

# 2. For each conflicted file, review both versions
git diff --name-only --diff-filter=U | while read file; do
  echo "=== Resolving: $file ==="
  
  # Show conflict markers
  grep -n "^<<<<<<< \|^======= \|^>>>>>>> " "$file"
  
  # Open in editor for manual resolution
  # Guideline: Preserve functionality from both branches
done

# 3. After resolving, stage the files
git add <resolved-files>

# 4. Complete the merge
git commit --no-edit

# 5. Verify resolution correctness
go test ./... -v
```

## Post-Merge Validation

### Final Integration Tests
```bash
# 1. Comprehensive build check
echo "==== Running full build ===="
go build -v ./...

# 2. Run all tests with coverage
echo "==== Running all tests ===="
go test ./... -v -cover -coverprofile=phase1-coverage.out

# 3. Generate coverage report
go tool cover -html=phase1-coverage.out -o phase1-coverage.html

# 4. Verify line count compliance
echo "==== Verifying size compliance ===="
/home/vscode/workspaces/idpbuilder-oci-go-cr/tools/line-counter.sh

# 5. Check for any uncommitted changes
git status --porcelain
```

### Integration Success Criteria
- ✅ All 4 efforts successfully merged
- ✅ No unresolved conflicts
- ✅ All tests passing
- ✅ Build successful
- ✅ Total line count within limits (2443 lines expected)
- ✅ No regression in existing functionality

## Rollback Plan

If critical issues are discovered during merging:

```bash
# 1. Abort current merge if in progress
git merge --abort

# 2. Reset to last known good state
git reset --hard HEAD

# 3. Review problematic effort branch
git log --oneline <problematic-branch> -10

# 4. Create fix branch if needed
git checkout -b fix/<effort-name>-integration-issues

# 5. Apply fixes and retry merge
```

## Final Steps

### Create Integration Tag
```bash
# After all merges complete successfully
git tag -a "phase1-integration-complete" -m "Phase 1 integration complete
- Wave 1: kind-certificate-extraction, registry-tls-trust-integration
- Wave 2: certificate-validation-pipeline, fallback-strategies
Total lines: ~2443
[R269 Protocol Complete]"

# Push integration branch and tag
git push origin idpbuilder-oci-go-cr/phase1-integration-20250902-194557
git push origin phase1-integration-complete
```

### Update Orchestrator State
```yaml
# Update orchestrator-state.yaml after successful integration
phase1_integration:
  status: COMPLETED
  branch: idpbuilder-oci-go-cr/phase1-integration-20250902-194557
  integrated_efforts:
    - kind-certificate-extraction (418 lines)
    - registry-tls-trust-integration (936 lines)
    - certificate-validation-pipeline (431 lines)
    - fallback-strategies (658 lines)
  total_lines: 2443
  integration_date: 2025-01-02
```

## Important Notes

1. **Branch Names**: The "idpbuidler" spelling (missing 'l') in branch names is CORRECT - do not fix
2. **Sequential Execution**: Execute merges in exact order listed
3. **Wave Dependencies**: ALL Wave 1 efforts must merge before ANY Wave 2 efforts
4. **No Direct Execution**: This plan is for orchestrator guidance - do not execute directly
5. **Conflict Resolution**: Manual intervention may be required for conflicts
6. **Testing**: Run tests after each merge to catch issues early

## Troubleshooting

### Common Issues and Solutions

#### Issue: Branch not found
```bash
# Solution: Fetch from correct remote
git fetch origin <branch-name>
# Or check effort directory for local branch
cd /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/wave[X]/[effort-name]
git branch --show-current
```

#### Issue: Merge conflicts in go.mod
```bash
# Solution: Regenerate after merge
go mod tidy
git add go.mod go.sum
git commit -m "fix: resolve go.mod conflicts"
```

#### Issue: Test failures after merge
```bash
# Solution: Check for interface changes
git diff HEAD~1 HEAD -- '*.go' | grep "^[+-].*interface"
# Update implementations to match new interfaces
```

---

**Document Version**: 1.0  
**Protocol**: R269 - Phase Integration Protocol  
**Created**: 2025-01-02  
**Author**: Code Reviewer Agent  