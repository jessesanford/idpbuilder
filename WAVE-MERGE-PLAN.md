# Phase 1 Wave 1 Integration Merge Plan

## 🎯 Merge Planning Metadata
- **Plan Created By**: Code Reviewer Agent (WAVE_MERGE_PLANNING state)
- **Plan Creation Date**: 2025-08-29 05:44:00 UTC
- **Phase**: 1
- **Wave**: 1
- **Integration Branch**: `idpbuilder-oci-mvp/phase1/wave1/integration-20250829-054225`
- **Integration Directory**: `/home/vscode/workspaces/idpbuilder-oci-mvp/efforts/phase1/wave1/integration-workspace`
- **Target Repository**: https://github.com/jessesanford/idpbuilder.git

## 📊 Effort Status Analysis

### 1. cert-extraction
- **Status**: ACCEPTED_WITH_SPLITS (merged back to single branch)
- **Source Branch**: `idpbuilder-oci-mvp/phase1/wave1/cert-extraction`
- **Location**: `/home/vscode/workspaces/idpbuilder-oci-mvp/efforts/phase1/wave1/cert-extraction`
- **Implementation Size**: 3639 lines (originally split into 602 + 799 lines)
- **Current Branch Status**: Already has integration commits
- **Files Modified**: 22 files
- **Key Components**:
  - `pkg/certs/types.go` - Type definitions for cert extraction
  - `pkg/certs/extractor.go` - Certificate extraction logic
  - `pkg/certs/validator.go` - Certificate validation
  - Comprehensive test coverage

### 2. trust-store
- **Status**: ACCEPTED
- **Source Branch**: `idpbuilder-oci-mvp/phase1/wave1/trust-store`
- **Location**: `/home/vscode/workspaces/idpbuilder-oci-mvp/efforts/phase1/wave1/trust-store`
- **Implementation Size**: 677 lines (within limit)
- **Files Modified**: 9 files
- **Key Components**:
  - `pkg/certs/types.go` - Type definitions for trust store
  - `pkg/certs/registry.go` - Registry trust store management
  - `pkg/certs/buildah.go` - Buildah integration
  - Test files with good coverage

## 🔀 Merge Order Strategy

### Recommended Merge Sequence:

1. **FIRST: cert-extraction**
   - **Reason**: Foundational certificate extraction functionality
   - **Branch to use**: `idpbuilder-oci-mvp/phase1/wave1/cert-extraction`
   - **Note**: Contains the complete implementation (both splits merged)

2. **SECOND: trust-store**
   - **Reason**: Builds on cert functionality, will have types.go conflict
   - **Branch to use**: `idpbuilder-oci-mvp/phase1/wave1/trust-store`
   - **Conflict Resolution Required**: YES (see below)

## ⚠️ Expected Conflicts and Resolution Strategies

### 1. pkg/certs/types.go Conflict
**Location**: `pkg/certs/types.go`
**Conflict Type**: Both efforts create this file with different content
**Severity**: HIGH - Requires careful manual resolution

**Resolution Strategy**:
```go
// Merged types.go structure:
package certs

import (
    "context"
    "crypto/x509"
    "io/fs"
    "time"
)

// === CERT-EXTRACTION TYPES ===

// KindCertExtractor defines the interface for extracting certificates
type KindCertExtractor interface {
    ExtractGiteaCert(ctx context.Context) (*x509.Certificate, error)
    GetClusterName() string
}

// CertValidator defines the interface for validating certificates
type CertValidator interface {
    ValidateCertificate(cert *x509.Certificate) (*ValidationResult, error)
    CheckExpiry(cert *x509.Certificate, warnDays int) (*ExpiryResult, error)
}

// ExtractorConfig holds configuration for the certificate extractor
type ExtractorConfig struct {
    ClusterName string
    Namespace   string
    PodSelector string
    CertPath    string
}

// === TRUST-STORE TYPES ===

// TrustStoreLocation defines where trust store certificates are stored
type TrustStoreLocation int

const (
    UserTrustStore TrustStoreLocation = iota
    SystemTrustStore
)

// CertificateInfo represents metadata about a certificate
type CertificateInfo struct {
    Subject       string
    Issuer        string
    SerialNumber  string
    NotBefore     string
    NotAfter      string
    Fingerprint   string
}
```

### 2. Potential Import Conflicts
**Impact**: LOW
**Resolution**: Ensure all imports are combined and deduplicated

## 📝 Integration Agent Execution Instructions

### Pre-merge Verification Commands:
```bash
# 1. Ensure you're in the integration workspace
cd /home/vscode/workspaces/idpbuilder-oci-mvp/efforts/phase1/wave1/integration-workspace

# 2. Verify you're on the integration branch
git branch --show-current
# Expected: idpbuilder-oci-mvp/phase1/wave1/integration-20250829-054225

# 3. Ensure clean working tree
git status --porcelain
# Should be empty

# 4. Fetch latest from all effort branches
git fetch origin
```

### Merge Execution Commands:

#### Step 1: Merge cert-extraction
```bash
# Navigate to cert-extraction effort
cd /home/vscode/workspaces/idpbuilder-oci-mvp/efforts/phase1/wave1/cert-extraction

# Ensure latest commits are available
git log --oneline -1

# Return to integration workspace
cd /home/vscode/workspaces/idpbuilder-oci-mvp/efforts/phase1/wave1/integration-workspace

# Merge cert-extraction (using --allow-unrelated-histories if needed)
git merge --no-ff ../cert-extraction/.git/refs/heads/idpbuilder-oci-mvp/phase1/wave1/cert-extraction \
  -m "integrate: cert-extraction effort (Phase 1 Wave 1)"
  
# Alternative if above doesn't work:
# Add cert-extraction as remote and merge
git remote add cert-extraction ../cert-extraction/.git
git fetch cert-extraction
git merge --no-ff cert-extraction/idpbuilder-oci-mvp/phase1/wave1/cert-extraction \
  -m "integrate: cert-extraction effort (Phase 1 Wave 1)"
```

#### Step 2: Merge trust-store with Conflict Resolution
```bash
# Navigate to trust-store effort
cd /home/vscode/workspaces/idpbuilder-oci-mvp/efforts/phase1/wave1/trust-store

# Ensure latest commits
git log --oneline -1

# Return to integration workspace
cd /home/vscode/workspaces/idpbuilder-oci-mvp/efforts/phase1/wave1/integration-workspace

# Start merge (expect conflict)
git merge --no-ff ../trust-store/.git/refs/heads/idpbuilder-oci-mvp/phase1/wave1/trust-store \
  -m "integrate: trust-store effort (Phase 1 Wave 1)"

# Alternative method:
git remote add trust-store ../trust-store/.git
git fetch trust-store
git merge --no-ff trust-store/idpbuilder-oci-mvp/phase1/wave1/trust-store

# When conflict occurs in pkg/certs/types.go:
# 1. Open the file and manually combine both type definitions
# 2. Keep ALL types from both efforts
# 3. Organize with clear section comments
# 4. Ensure no duplicate type names

# After resolving conflicts:
git add pkg/certs/types.go
git commit -m "resolve: merge conflict in types.go between cert-extraction and trust-store"
```

### Post-merge Validation:
```bash
# 1. Verify all files are present
ls -la pkg/certs/
# Should show files from both efforts

# 2. Run build to ensure compilation
go build ./pkg/certs/...

# 3. Run tests
go test ./pkg/certs/...

# 4. Check final line count
find pkg/certs -name "*.go" | xargs wc -l

# 5. Commit final integration
git commit -m "complete: Phase 1 Wave 1 integration"

# 6. Push to remote
git push -u origin idpbuilder-oci-mvp/phase1/wave1/integration-20250829-054225
```

## 🎬 Integration Completion Criteria

✅ Both effort branches successfully merged
✅ types.go conflict resolved with all types preserved
✅ All tests passing
✅ Code compiles without errors
✅ Integration branch pushed to remote
✅ No files lost during merge
✅ Clear commit history showing merge progression

## 📋 Troubleshooting Guide

### If merge fails with "unrelated histories":
```bash
git merge --allow-unrelated-histories [branch] -m "merge: [effort-name] with unrelated histories"
```

### If remote tracking is incorrect:
```bash
# Add effort as local remote
git remote add [effort-name] ../[effort-name]/.git
git fetch [effort-name]
git merge [effort-name]/[branch-name]
```

### If types.go conflict is complex:
1. Save both versions to temporary files
2. Create new types.go from scratch
3. Copy types from both versions
4. Ensure proper imports
5. Test compilation after each addition

## 🚨 Critical Warnings

1. **DO NOT** skip the types.go conflict resolution
2. **DO NOT** delete any type definitions during merge
3. **DO NOT** change the merge order (cert-extraction must be first)
4. **DO NOT** force push without verification
5. **ALWAYS** run tests after each merge

## 📝 Notes for Integration Agent

This plan is created per R269 requirements. The Integration Agent should:
1. Execute these commands sequentially
2. Handle conflicts as described
3. Report any deviations from expected behavior
4. Complete all validation steps before marking integration complete

The most critical aspect is the types.go conflict resolution - both sets of types must be preserved and organized cleanly in the merged file.

---
*Plan created by Code Reviewer Agent in WAVE_MERGE_PLANNING state*
*R269 Compliant: Plan only, no execution performed*