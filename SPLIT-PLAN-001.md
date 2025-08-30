# SPLIT-PLAN-001.md

## Split 001 of 2: Core Certificate Extraction
**Planner**: Code Reviewer Agent (same for ALL splits)
**Parent Effort**: cert-extraction

<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (⚠️⚠️⚠️ CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: None (first split of THIS effort)
  - Path: N/A (this is Split 001)
  - Branch: N/A
- **This Split**: Split 001 of phase1/wave1/cert-extraction
  - Path: efforts/phase1/wave1/cert-extraction/split-001/
  - Branch: idpbuilder-oci-mvp/phase1/wave1/cert-extraction-split-001
- **Next Split**: Split 002 of phase1/wave1/cert-extraction
  - Path: efforts/phase1/wave1/cert-extraction/split-002/
  - Branch: idpbuilder-oci-mvp/phase1/wave1/cert-extraction-split-002
- **File Boundaries**:
  - This Split: types.go, errors.go, extractor.go (~493 lines total)
  - Next Split: validator.go, test files

## Files in This Split (EXCLUSIVE - no overlap with other splits)
- `pkg/certs/types.go` (68 lines) - Core interfaces and types
- `pkg/certs/errors.go` (104 lines) - Error handling framework  
- `pkg/certs/extractor.go` (321 lines) - Main extraction logic

**Total Lines**: ~493 lines

## Functionality
### Core Certificate Extraction
- Define KindCertExtractor and CertValidator interfaces
- Implement ExtractorConfig and diagnostic types
- Create comprehensive error handling with suggestions
- Implement full certificate extraction from Kind clusters:
  - Kind cluster detection and connection
  - Gitea pod discovery using Kubernetes client
  - Certificate extraction via kubectl exec
  - Certificate parsing from PEM format
  - Storage to ~/.idpbuilder/certs/ directory

### Key Methods to Implement
- `NewKindExtractor()` - Create extractor with Kubernetes client
- `ExtractGiteaCert()` - Main extraction orchestration
- `findGiteaPod()` - Locate Gitea pod in cluster
- `extractCertFromPod()` - Execute command to get cert
- `parseCertificate()` - Parse PEM to x509
- `storeCertificate()` - Save to local filesystem
- `GetClusterName()` - Return configured cluster name

## Dependencies
- External packages needed:
  - `k8s.io/client-go` - Kubernetes client libraries
  - `k8s.io/api/core/v1` - Core Kubernetes types
  - `k8s.io/apimachinery` - Kubernetes API machinery
  - `sigs.k8s.io/controller-runtime` - Logging framework
  - Standard library crypto/x509 for certificates

## Implementation Instructions

### Step 1: Create Package Structure
```bash
mkdir -p pkg/certs
cd pkg/certs
```

### Step 2: Implement types.go
- Define all interfaces (KindCertExtractor, CertValidator)
- Create ExtractorConfig struct with defaults
- Define CertDiagnostics for reporting

### Step 3: Implement errors.go  
- Create CertificateError type with context
- Define predefined errors (ErrClusterNotFound, etc.)
- Implement error wrapping and suggestions

### Step 4: Implement extractor.go
- Create KindExtractor struct
- Implement NewKindExtractor with client setup
- Add ExtractGiteaCert main logic
- Implement helper methods for pod operations
- Add certificate parsing and storage

### Step 5: Basic Testing (Optional for Split 001)
- Create minimal smoke test to verify compilation
- Full test suite will be in Split 002

### Step 6: Measure and Verify
```bash
# From split directory
$PROJECT_ROOT/tools/line-counter.sh
# Must be under 700 lines (target) / 800 lines (hard limit)
```

## Split Branch Strategy
- Branch: `idpbuilder-oci-mvp/phase1/wave1/cert-extraction-split-001`
- Base: `main`
- Must merge to: `idpbuilder-oci-mvp/phase1/wave1/cert-extraction` after review

## Success Criteria
- ✅ All three files implemented completely
- ✅ Code compiles without errors
- ✅ Basic extraction functionality works
- ✅ Total size under 700 lines (hard limit 800)
- ✅ Follows idpbuilder patterns and conventions
- ✅ Ready for Split 002 to add validation and tests

## Notes for SW Engineer
- Focus on core functionality only - no validation logic
- Validation interface defined but not implemented (Split 002)
- Keep error messages clear and actionable
- Use established idpbuilder logging patterns
- This split provides the foundation for Split 002