# Phase 2 Wave 1 Integration Merge Plan

## 📋 Integration Overview
**Phase**: 2 - Build & Push Implementation  
**Wave**: 1 - Core OCI Operations  
**Integration Branch**: `idpbuilder-oci-go-cr/phase2/wave1-integration-20250904-212505`  
**Base Branch**: `idpbuilder-oci-go-cr/phase1-integration-20250902-194557`  
**Created**: 2025-09-04 21:25:05 UTC  
**Total Efforts**: 2  

## ✅ Pre-Merge Validation Status
| Effort | Branch | Lines | Review Status |
|--------|--------|-------|--------------|
| E2.1.1: go-containerregistry-image-builder | `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder` | 756 | ✅ PASSED |
| E2.1.2: gitea-registry-client | `idpbuilder-oci-go-cr/phase2/wave1/gitea-registry-client` | 689 | ✅ PASSED |

## 🔀 Merge Order and Rationale

### Merge Sequence
1. **E2.1.1: go-containerregistry-image-builder** (FIRST)
2. **E2.1.2: gitea-registry-client** (SECOND)

### Rationale for Order
1. **E2.1.1 First**: 
   - Implements core image building functionality
   - Created first (2025-09-02 22:47:36 UTC)
   - Provides foundational OCI operations
   - No dependency on E2.1.2
   - Size: 756 lines (within warning threshold but compliant)

2. **E2.1.2 Second**: 
   - Implements registry client operations
   - Created after E2.1.1 (2025-09-02 22:52:06 UTC)
   - May benefit from E2.1.1's OCI structures
   - Size: 689 lines (fully compliant)

## 📝 Detailed Merge Instructions

### Step 1: Pre-Merge Verification
```bash
# Ensure we're on the integration branch
cd /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase2/wave1/integration-workspace
git branch --show-current
# Expected: idpbuilder-oci-go-cr/phase2/wave1-integration-20250904-212505

# Verify clean working tree
git status
# Expected: working tree clean

# Fetch latest changes
git fetch origin

# Verify base commit matches Phase 1 integration
git merge-base HEAD origin/idpbuilder-oci-go-cr/phase1-integration-20250902-194557
```

### Step 2: Merge E2.1.1 (go-containerregistry-image-builder)
```bash
# Add the effort directory as a remote (if not already added)
git remote add effort-e211 ../go-containerregistry-image-builder/.git 2>/dev/null || true

# Fetch the effort branch
git fetch effort-e211 idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder:refs/remotes/effort-e211/go-containerregistry-image-builder

# Merge E2.1.1
git merge effort-e211/go-containerregistry-image-builder \
  --no-ff \
  -m "integrate(phase2/wave1): Merge E2.1.1 go-containerregistry-image-builder (756 lines)

- Core OCI image building functionality
- go-containerregistry library integration
- Daemonless image assembly from directories
- OCI tarball support
- Integration with Phase 1 certificate infrastructure"

# Verify merge success
if [ $? -ne 0 ]; then
    echo "⚠️ CONFLICT DETECTED - Manual resolution required"
    echo "Resolve conflicts, then run: git add . && git commit"
else
    echo "✅ E2.1.1 merged successfully"
fi
```

### Step 3: Test After First Merge
```bash
# Run tests for merged code
cd pkg/build
go test ./... -v

# Check compilation
cd ../..
go build ./pkg/...

# Verify no regressions
go test ./pkg/certs/... -v  # Phase 1 components should still work
```

### Step 4: Merge E2.1.2 (gitea-registry-client)
```bash
# Add the effort directory as a remote (if not already added)
git remote add effort-e212 ../gitea-registry-client/.git 2>/dev/null || true

# Fetch the effort branch
git fetch effort-e212 idpbuilder-oci-go-cr/phase2/wave1/gitea-registry-client:refs/remotes/effort-e212/gitea-registry-client

# Merge E2.1.2
git merge effort-e212/gitea-registry-client \
  --no-ff \
  -m "integrate(phase2/wave1): Merge E2.1.2 gitea-registry-client (689 lines)

- Gitea-specific registry client implementation
- OCI image push and pull operations
- Authentication mechanisms for Gitea registry
- Progress tracking for push operations
- Integration with Phase 1 TrustStoreManager"

# Verify merge success
if [ $? -ne 0 ]; then
    echo "⚠️ CONFLICT DETECTED - Manual resolution required"
    echo "Resolve conflicts, then run: git add . && git commit"
else
    echo "✅ E2.1.2 merged successfully"
fi
```

### Step 5: Post-Merge Integration Testing
```bash
# Full test suite after both merges
go test ./... -v

# Check for any compilation issues
go build ./...

# Verify all packages are properly integrated
ls -la pkg/
# Expected directories:
# - pkg/build/        (from E2.1.1)
# - pkg/registry/     (from E2.1.2)
# - pkg/certs/        (from Phase 1)
# - pkg/api/          (if contracts were shared)

# Run integration tests
go test ./tests/integration/... -v 2>/dev/null || echo "No integration tests yet"

# Check line count of integrated branch
/home/vscode/workspaces/idpbuilder-oci-go-cr/tools/line-counter.sh
# Expected: ~1445 lines total (756 + 689)
```

## 🔧 Conflict Resolution Strategy

### Expected Conflict Points
1. **go.mod/go.sum**: Both efforts may add dependencies
   - Resolution: Accept both sets of dependencies
   - Run `go mod tidy` after resolution

2. **pkg/api/interfaces.go**: If both created shared contracts
   - Resolution: Merge interface definitions, avoid duplicates
   - Ensure both BuildService and RegistryService are present

3. **README or documentation**: Both may have updated docs
   - Resolution: Combine both sets of changes chronologically

### Conflict Resolution Process
```bash
# If conflicts occur:
1. git status                    # See conflicted files
2. vim <conflicted-file>         # Resolve conflicts
3. git add <resolved-file>       # Stage resolution
4. git commit                    # Complete merge
5. go test ./...                 # Verify resolution didn't break anything
```

## 📊 Success Criteria

### Integration Complete When:
- [ ] Both effort branches merged successfully
- [ ] No unresolved conflicts
- [ ] All tests passing (`go test ./...`)
- [ ] Code compiles without errors (`go build ./...`)
- [ ] Total line count under limits (use line-counter.sh)
- [ ] Integration branch pushed to origin

### Post-Integration Checklist:
```bash
# Verify integration success
echo "=== Integration Validation Checklist ==="
echo "1. Check merged commits:"
git log --oneline -10

echo "2. Verify both efforts are included:"
git log --oneline --grep="E2.1.1" -1
git log --oneline --grep="E2.1.2" -1

echo "3. Test compilation:"
go build ./...

echo "4. Run tests:"
go test ./... -count=1

echo "5. Check line count:"
/home/vscode/workspaces/idpbuilder-oci-go-cr/tools/line-counter.sh

echo "6. Push integration branch:"
git push origin HEAD
```

## 🚨 Important Notes

### Do NOT Execute These Commands
This plan is for the Integration Agent to execute. As Code Reviewer, you should:
- ✅ CREATE this plan
- ❌ NOT execute the git merge commands
- ❌ NOT modify the integration branch

### Split Branches
No split branches were created for Wave 1 efforts. Both efforts remained under the 800-line hard limit:
- E2.1.1: 756 lines (warning threshold but compliant)
- E2.1.2: 689 lines (fully compliant)

### Dependencies
Both efforts were developed in parallel as specified in the wave plan:
- They share Phase 1 as a common base
- No inter-dependencies between E2.1.1 and E2.1.2
- Both integrate with Phase 1's certificate infrastructure

## 📝 Integration Agent Instructions

To execute this plan:
1. Navigate to integration workspace
2. Follow steps 1-5 in sequence
3. Handle any conflicts using the resolution strategy
4. Complete the post-integration checklist
5. Report status back to orchestrator

**Integration Command Summary**:
```bash
# Quick execution (for Integration Agent only)
cd /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase2/wave1/integration-workspace
git merge ../go-containerregistry-image-builder --no-ff
git merge ../gitea-registry-client --no-ff
go test ./...
git push origin HEAD
```

---
**Plan Created By**: Code Reviewer Agent  
**Date**: 2025-09-04  
**Status**: Ready for Integration Agent Execution