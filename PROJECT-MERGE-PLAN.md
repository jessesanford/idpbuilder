# PROJECT MERGE PLAN - IDPBuilder OCI Build Push

## Executive Summary
This merge plan coordinates the integration of Phase 2, Wave 1 efforts into the project-integration branch. This is a fresh re-integration per R327 after fixing import path and API compatibility issues.

**Generated**: 2025-01-09 19:17:00 UTC
**Integration Branch**: project-integration (fresh from main)
**Target Repository**: github.com/cnoe-io/idpbuilder

## Pre-Merge State
- **Base Branch**: project-integration (clean from main)
- **Base Commit**: e210954 (todo(architect): WAVE_REVIEW checkpoint)
- **Remote**: origin (https://github.com/jessesanford/idpbuilder.git)

## Branches to Merge

### 1. idpbuilder-oci-build-push/phase2/wave1/image-builder
- **Purpose**: OCI image building functionality using go-containerregistry
- **Size**: ~601 lines (verified compliant)
- **Commits**: 33 commits since base
- **Status**: Import paths corrected to github.com/cnoe-io/idpbuilder
- **Dependencies**: None (foundational component)
- **Key Files**:
  - pkg/build/image_builder.go (core builder)
  - pkg/build/context.go (build context management)
  - pkg/build/storage.go (storage interface)

### 2. idpbuilder-oci-build-push/phase2/wave1/gitea-client
- **Purpose**: Gitea registry client implementation with Phase 1 integration
- **Size**: ~1200 lines total (split branches merged)
- **Commits**: 35 commits since base
- **Status**: Import paths and API calls corrected
- **Dependencies**: Depends on image-builder for types
- **Key Files**:
  - pkg/registry/* (registry client implementation)
  - pkg/gitea/* (Gitea-specific functionality)

## Merge Order and Strategy

### Phase 1: Dependency Resolution
```bash
# Order: Dependencies first, then dependent components
1. image-builder (no dependencies)
2. gitea-client (depends on image-builder)
```

### Phase 2: Merge Execution Steps

#### Step 1: Merge image-builder
```bash
git checkout project-integration
git pull origin project-integration
git merge origin/idpbuilder-oci-build-push/phase2/wave1/image-builder --no-ff \
  -m "feat: integrate OCI image builder functionality (Phase 2, Wave 1)"
```

**Expected Conflicts**: 
- go.mod/go.sum (dependency versions)
- pkg/certs/* (overlapping certificate handling)

**Resolution Strategy**:
- For go.mod: Keep union of dependencies, resolve version conflicts to highest
- For pkg/certs: Prefer image-builder version (more recent implementation)

#### Step 2: Validate image-builder integration
```bash
go mod tidy
go build ./pkg/build/...
go test ./pkg/build/... -v
```

#### Step 3: Merge gitea-client
```bash
git merge origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client --no-ff \
  -m "feat: integrate Gitea registry client (Phase 2, Wave 1)"
```

**Expected Conflicts**:
- go.mod/go.sum (additional dependencies)
- pkg/certs/* (significant overlaps with image-builder)
- pkg/certvalidation/* (possible duplicates)

**Resolution Strategy**:
- For go.mod: Continue union approach, verify no version downgrades
- For pkg/certs: 
  - If identical: Keep one version
  - If different: Merge functionality, prefer gitea-client for registry-specific code
- For pkg/certvalidation: Deduplicate, keeping most comprehensive implementation

#### Step 4: Final validation
```bash
go mod tidy
go build ./...
go test ./... -v
```

## Conflict Resolution Guidelines

### File-Level Conflicts

#### go.mod/go.sum
- **Strategy**: Union of dependencies
- **Version Resolution**: Always choose higher version unless breaking
- **Verification**: Run `go mod tidy` after each merge

#### pkg/certs/* (High Overlap Area)
- **Primary Source**: image-builder (merged first)
- **Enhancement Source**: gitea-client (registry-specific additions)
- **Resolution Process**:
  1. Accept image-builder base implementation
  2. Layer gitea-client enhancements
  3. Deduplicate identical functions
  4. Preserve registry-specific extensions

#### pkg/certvalidation/*
- **Strategy**: Merge unique functionality
- **Deduplication**: Remove identical implementations
- **Testing**: Ensure all cert validation tests pass

#### pkg/fallback/* and pkg/insecure/*
- **Strategy**: These appear in both branches
- **Resolution**: Keep most comprehensive implementation
- **Validation**: Security review required

## Validation Steps

### 1. Build Validation
```bash
# Clean build
go clean -cache
go mod download
go build -v ./...

# Verify no build errors
echo "Build Status: $?"
```

### 2. Test Execution
```bash
# Run all tests with coverage
go test ./... -v -cover -coverprofile=coverage.out

# Check coverage metrics
go tool cover -func=coverage.out | tail -5
```

### 3. Import Path Verification
```bash
# Verify all imports use correct path
grep -r "github.com/cnoe-io/idpbuilder" pkg/ | wc -l
grep -r "github.com/jessesanford/idpbuilder" pkg/ | wc -l  # Should be 0
```

### 4. Integration Tests
```bash
# Run integration tests if available
go test ./... -tags=integration -v
```

### 5. Dependency Audit
```bash
# Check for security vulnerabilities
go list -m all | nancy sleuth

# Verify no duplicate dependencies
go mod graph | grep "@.*@" | sort | uniq -d
```

## Success Criteria

### Must Pass (Blocking)
- ✅ All code compiles without errors
- ✅ All unit tests pass
- ✅ No import path references to jessesanford/idpbuilder
- ✅ go.mod is clean (go mod tidy produces no changes)
- ✅ No duplicate code in pkg/certs

### Should Pass (Non-Blocking)
- ⚠️ Integration tests pass (if present)
- ⚠️ Coverage >80% for new code
- ⚠️ No new security vulnerabilities
- ⚠️ Documentation updated

### Nice to Have
- 📝 CHANGELOG updated
- 📝 Example usage documented
- 📝 Performance benchmarks pass

## Rollback Plan

If merge fails catastrophically:

### Option 1: Reset to Clean State
```bash
git reset --hard origin/project-integration
git clean -fd
```

### Option 2: Create New Integration Branch
```bash
git checkout -b project-integration-v2 origin/main
# Retry merges with different strategy
```

### Option 3: Sequential Cherry-Pick
```bash
# Cherry-pick specific commits if full merge fails
git cherry-pick <commit-range>
```

## Post-Merge Actions

### 1. Documentation Updates
- Update main README with new OCI capabilities
- Document Gitea registry integration
- Add configuration examples

### 2. CI/CD Verification
- Ensure GitHub Actions pass
- Verify container builds work
- Test deployment pipelines

### 3. Communication
- Notify team of successful integration
- Document any workarounds needed
- Update project board status

## Risk Assessment

### High Risk Areas
1. **pkg/certs overlap**: ~20 files overlap between branches
   - Mitigation: Careful deduplication, comprehensive testing
   
2. **go.mod conflicts**: Different dependency versions
   - Mitigation: Use go mod tidy, test thoroughly

3. **Import path references**: Ensuring all fixed
   - Mitigation: Grep verification post-merge

### Medium Risk Areas
1. **Test compatibility**: Tests may have assumptions
   - Mitigation: Run full test suite after each merge

2. **Feature flags**: Ensuring proper isolation
   - Mitigation: Verify feature flag implementation

### Low Risk Areas
1. **Documentation files**: Non-critical conflicts
   - Mitigation: Manual review and merge

## Timeline Estimate

| Step | Duration | Notes |
|------|----------|-------|
| Pre-merge validation | 5 min | Verify clean state |
| Merge image-builder | 10 min | Including conflict resolution |
| Validate image-builder | 5 min | Build and test |
| Merge gitea-client | 15 min | More conflicts expected |
| Final validation | 10 min | Full test suite |
| Documentation | 5 min | Update merge results |
| **Total** | **~50 min** | With buffer for issues |

## Integration Agent Instructions

### For the Integration Agent executing this plan:

1. **Start State**: Verify you're on clean project-integration branch
2. **Execute Sequentially**: Follow merge order exactly
3. **Conflict Resolution**: Use strategies documented above
4. **Validation Gates**: Stop if any Must Pass criteria fail
5. **Documentation**: Create INTEGRATION-REPORT.md with results
6. **State Updates**: Update orchestrator-state.json after completion

### Expected Deliverables:
- ✅ Merged project-integration branch
- ✅ INTEGRATION-REPORT.md with results
- ✅ Updated orchestrator-state.json
- ✅ Clean go.mod with no conflicts
- ✅ All tests passing

---

**Plan Created By**: Code Reviewer Agent
**Plan Location**: /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/project/integration-workspace/PROJECT-MERGE-PLAN.md
**Ready for Execution**: YES

## Appendix: Detected Overlapping Files

The following files appear in both branches and will require careful merge resolution:

```
go.mod
go.sum
pkg/certs/chain_validator.go
pkg/certs/chain_validator_test.go
pkg/certs/diagnostics.go
pkg/certs/errors.go
pkg/certs/errors_test.go
pkg/certs/extractor.go
pkg/certs/extractor_test.go
pkg/certs/helpers.go
pkg/certs/helpers_test.go
pkg/certs/kind_client.go
pkg/certs/storage.go
pkg/certs/storage_test.go
pkg/certs/trust.go
pkg/certs/trust_test.go
pkg/certs/utilities.go
pkg/certs/utilities_test.go
pkg/certs/validation_errors.go
pkg/certs/validator_test.go
pkg/certvalidation/chain_validator.go
pkg/certvalidation/chain_validator_test.go
pkg/certvalidation/x509_utils.go
pkg/certvalidation/x509_utils_test.go
pkg/fallback/interfaces.go
pkg/fallback/manager.go
pkg/fallback/manager_test.go
pkg/fallback/strategies.go
pkg/fallback/strategies_test.go
pkg/insecure/handler.go
pkg/insecure/handler_test.go
```

Total: 31 overlapping source files requiring resolution.