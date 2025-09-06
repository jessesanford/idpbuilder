# Phase 1 Wave 1 - Merge Plan

## 🚨 CRITICAL RULES (R269 & R296 Compliance)
- **USE ONLY ORIGINAL EFFORT BRANCHES** - NO integration branch merges!
- **VERIFY NO DEPRECATED BRANCHES** - Check for "-deprecated-split" suffix
- **FOLLOW EXACT SEQUENCE** - Based on dependencies and bases
- **VALIDATE AFTER EACH MERGE** - Tests must pass

## 📊 Merge Plan Summary

**Generated**: 2025-09-06 19:30:00 UTC  
**Phase**: 1 - Certificate Infrastructure  
**Wave**: 1 - Certificate Management Core  
**Target Branch**: `idpbuilder-oci-build-push/phase1/wave1-integration`  
**Target Repository**: https://github.com/jessesanford/idpbuilder.git  
**Integration Directory**: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave1/integration-workspace`

## 🔍 Branch Analysis Results

### Effort Branches Analyzed
1. **E1.1.1 - Kind Certificate Extraction**
   - Branch: `phase1/wave1/effort-kind-cert-extraction`
   - Status: COMPLETED
   - Review: ACCEPTED
   - Merge Base: e210954 (common ancestor with main)
   - Files Modified: 9 files in pkg/certs/ (no conflicts)
   - Latest Commit: 2757831

2. **E1.1.2 - Registry TLS Trust Integration**
   - Branch: `phase1/wave1/effort-registry-tls-trust`
   - Status: COMPLETED
   - Review: APPROVED_WITH_WARNINGS
   - Merge Base: e210954 (common ancestor with main)
   - Files Modified: 8 files in pkg/certs/ (no conflicts)
   - Latest Commit: 5344541

### Deprecated Branches Check
✅ **NO DEPRECATED BRANCHES FOUND**
- Searched for pattern: `*-deprecated-split*`
- Result: None found

### Conflict Analysis
✅ **NO FILE CONFLICTS DETECTED**

**E1.1.1 Files**:
- pkg/certs/errors.go
- pkg/certs/errors_test.go
- pkg/certs/extractor.go
- pkg/certs/extractor_test.go
- pkg/certs/helpers.go
- pkg/certs/helpers_test.go
- pkg/certs/kind_client.go
- pkg/certs/storage.go
- pkg/certs/storage_test.go

**E1.1.2 Files**:
- pkg/certs/trust.go
- pkg/certs/trust_test.go
- pkg/certs/utilities.go
- pkg/certs/utilities_test.go
- go.mod
- go.sum

**Conclusion**: No overlapping source files between efforts. Both modify different components within pkg/certs/.

## 📋 Merge Sequence

### Order Justification
Both efforts have:
- Same merge base (e210954)
- No interdependencies (per Wave plan)
- No file conflicts
- Can be merged in any order

**Recommended Order**: E1.1.1 → E1.1.2
- Rationale: Follow effort numbering for consistency

## 🔧 Execution Instructions for Integration Agent

### Prerequisites
```bash
# Navigate to integration workspace
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave1/integration-workspace

# Verify on correct branch
git branch --show-current
# Expected: idpbuilder-oci-build-push/phase1/wave1-integration

# Ensure clean working directory
git status
# Expected: nothing to commit, working tree clean

# Fetch latest
git fetch origin
```

### Merge Step 1: E1.1.1 - Kind Certificate Extraction
```bash
# Navigate to integration workspace
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave1/integration-workspace

# Fetch the effort branch
git fetch origin phase1/wave1/effort-kind-cert-extraction:phase1/wave1/effort-kind-cert-extraction

# Merge E1.1.1
git merge phase1/wave1/effort-kind-cert-extraction --no-ff \
  -m "feat(phase1/wave1): integrate Kind Certificate Extraction (E1.1.1)

- Adds certificate extraction from Kind clusters
- Implements KindCertExtractor with storage capabilities
- Includes comprehensive error handling and testing
- Files: pkg/certs/{extractor,storage,kind_client,helpers,errors}*.go

Effort: E1.1.1
Review Status: ACCEPTED
Branch: phase1/wave1/effort-kind-cert-extraction"

# Validate compilation
go build ./...

# Run tests
go test ./pkg/certs/... -v

# Verify no test failures
if [ $? -ne 0 ]; then
    echo "❌ Tests failed after E1.1.1 merge"
    # Integration Agent should stop and report
fi
```

### Merge Step 2: E1.1.2 - Registry TLS Trust Integration
```bash
# Continue in integration workspace
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave1/integration-workspace

# Fetch the effort branch
git fetch origin phase1/wave1/effort-registry-tls-trust:phase1/wave1/effort-registry-tls-trust

# Merge E1.1.2
git merge phase1/wave1/effort-registry-tls-trust --no-ff \
  -m "feat(phase1/wave1): integrate Registry TLS Trust (E1.1.2)

- Adds TLS trust management for registry connections
- Implements certificate trust configuration
- Updates go.mod dependencies
- Files: pkg/certs/{trust,utilities}*.go, go.mod, go.sum

Effort: E1.1.2
Review Status: APPROVED_WITH_WARNINGS
Branch: phase1/wave1/effort-registry-tls-trust"

# Validate compilation
go build ./...

# Run tests
go test ./pkg/certs/... -v

# Verify no test failures
if [ $? -ne 0 ]; then
    echo "❌ Tests failed after E1.1.2 merge"
    # Integration Agent should stop and report
fi
```

### Post-Merge Validation
```bash
# Final comprehensive test
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave1/integration-workspace

# Run all tests
go test ./... -v

# Check test coverage
go test ./pkg/certs/... -cover

# Verify both features integrated
ls -la pkg/certs/ | grep -E "(extractor|storage|trust|utilities)"

# Build the final binary (if applicable)
go build -o idpbuilder-oci ./cmd/...

# Commit validation marker
echo "Wave 1 integration completed: $(date)" > .wave1-integration-complete
git add .wave1-integration-complete
git commit -m "chore: mark Wave 1 integration complete"

# Push the integration branch
git push origin idpbuilder-oci-build-push/phase1/wave1-integration
```

## ⚠️ Warning Notes

1. **Review Status**: E1.1.2 has APPROVED_WITH_WARNINGS - ensure warnings were non-critical
2. **Dependencies**: Both efforts update go.mod - verify no version conflicts after merge
3. **Test Coverage**: Ensure combined test coverage meets requirements (>80% unit tests)

## 🎯 Success Criteria

✅ All merges complete without conflicts  
✅ All tests pass after each merge  
✅ Final compilation successful  
✅ Integration branch pushed successfully  
✅ No deprecated branches were merged  
✅ Only original effort branches were used (R269 compliance)  

## 📝 Next Steps

After successful integration:
1. Notify Orchestrator of completion
2. Update orchestrator-state.yaml with integration status
3. Prepare for Phase 1 Wave 2 planning (if applicable)
4. Archive completed effort branches if needed

---
**Generated by**: Code Reviewer Agent  
**State**: WAVE_MERGE_PLANNING  
**Compliance**: R269, R296