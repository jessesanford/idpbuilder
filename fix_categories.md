# BACKPORT FIX CATEGORIZATION

## Fix Categories:

### 1. K8S.IO DEPENDENCY ALIGNMENT (v0.30.5) - ALL BRANCHES WITH GO.MOD
- Apply to: ALL 15 branches that have go.mod files
- Fix: Update all k8s.io dependencies to v0.30.5

### 2. CONTROLLER-RUNTIME v0.18.5 - ALL BRANCHES WITH GO.MOD
- Apply to: ALL 15 branches that have go.mod files  
- Fix: Update controller-runtime to v0.18.5

### 3. FALLBACK-SPECIFIC FIXES - ONLY Phase 1 Wave 2 fallback branches
- Apply to: ONLY these 3 branches:
  * idpbuilder-oci-build-push/phase1/wave2/fallback-core
  * idpbuilder-oci-build-push/phase1/wave2/fallback-recommendations  
  * idpbuilder-oci-build-push/phase1/wave2/fallback-security
- Fixes:
  * NewDefaultSecurityLogger function
  * Recommendation struct fields (RiskLevel, UserGuidance)
  * GetRecommendation method
  * SecurityLogEntry field alignment

### 4. GITEA CLIENT FIXES - ONLY gitea-client-split-001
- Apply to: ONLY this 1 branch:
  * idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001
- Fixes:
  * Interface alignment
  * retryWithExponentialBackoff function

### 5. TEST FIXES - Branches with test failures
- Apply to: Branches that have these specific test issues:
  * Remove duplicate TestExtraPortMappings
  * Remove unused imports
- Need to check each branch for these issues

## CONFIRMED EFFORT BRANCHES (15 total):

**Phase 1 Wave 1 (5 branches):**
1. idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction
2. idpbuilder-oci-build-push/phase1/wave1/registry-types
3. idpbuilder-oci-build-push/phase1/wave1/registry-auth
4. idpbuilder-oci-build-push/phase1/wave1/registry-helpers
5. idpbuilder-oci-build-push/phase1/wave1/registry-tests

**Phase 1 Wave 2 (4 branches):**
6. idpbuilder-oci-build-push/phase1/wave2/cert-validation
7. idpbuilder-oci-build-push/phase1/wave2/fallback-core
8. idpbuilder-oci-build-push/phase1/wave2/fallback-recommendations
9. idpbuilder-oci-build-push/phase1/wave2/fallback-security

**Phase 2 Wave 1 (3 branches):**
10. idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001
11. idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
12. idpbuilder-oci-build-push/phase2/wave1/image-builder

**Phase 2 Wave 2 (3 branches):**
13. idpbuilder-oci-build-push/phase2/wave2/cli-commands
14. idpbuilder-oci-build-push/phase2/wave2/credential-management
15. idpbuilder-oci-build-push/phase2/wave2/image-operations

