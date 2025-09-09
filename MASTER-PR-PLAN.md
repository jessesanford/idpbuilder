# MASTER PR PLAN
Generated: 2025-09-09T22:55:00Z
By: Software Factory 2.0

## 🚨 CRITICAL INSTRUCTIONS FOR HUMANS 🚨

**SOFTWARE FACTORY HAS COMPLETED ALL IMPLEMENTATION AND TESTING**

This document contains the exact steps to create Pull Requests for merging all completed work into the main branch of the idpbuilder repository.

### Prerequisites
1. Ensure you have write access to https://github.com/jessesanford/idpbuilder.git
2. Verify all branches listed below exist and are pushed
3. Have GitHub CLI (`gh`) installed and authenticated OR use GitHub web interface
4. Ensure you have the latest main branch locally

## 📋 VERIFICATION CHECKLIST

Before starting PR creation:
- [ ] All effort branches are pushed to remote
- [ ] Project integration branch shows successful build
- [ ] No uncommitted changes in any effort directory
- [ ] You understand the merge order below is MANDATORY

## 🔄 PR MERGE ORDER (CRITICAL - DO NOT DEVIATE)

**IMPORTANT**: PRs must be created and merged in this EXACT order to avoid conflicts.

The implementation follows the phases from the implementation plan:
- Phase 1: Certificate Infrastructure (Core foundation)
- Phase 2: Build & Push Implementation (Depends on Phase 1)

### Phase 1: Certificate Infrastructure PRs

#### PR #1: Kind Certificate Extraction
**Branch**: `idpbuilder-oci-build-push/phase1/wave1/kind-certificate-extraction`
**Depends On**: None (merge first)

**PR Title**:
```
feat(certs): implement Kind cluster certificate extraction for Gitea OCI registry
```

**PR Body**: See PR-BODY-1.md

**Merge Instructions**:
1. Create PR from `idpbuilder-oci-build-push/phase1/wave1/kind-certificate-extraction` to `main`
2. Wait for CI/CD checks to pass
3. Request code review if required
4. Merge using "Squash and merge" or project standard
5. Delete branch after merge

---

#### PR #2: Registry TLS Trust Integration
**Branch**: `idpbuilder-oci-build-push/phase1/wave1/registry-tls-trust-integration`
**Depends On**: PR #1 (MUST merge after certificate extraction)

**PR Title**:
```
feat(certs): add registry TLS trust store management for go-containerregistry
```

**PR Body**: See PR-BODY-2.md

---

#### PR #3: Certificate Validation Pipeline
**Branch**: `idpbuilder-oci-build-push/phase1/wave2/certificate-validation-pipeline`
**Depends On**: PRs #1 and #2

**PR Title**:
```
feat(validation): add certificate validation and monitoring pipeline
```

**PR Body**: See PR-BODY-3.md

---

#### PR #4: Fallback Strategies
**Branch**: `idpbuilder-oci-build-push/phase1/wave2/fallback-strategies`
**Depends On**: PRs #1, #2, and #3

**PR Title**:
```
feat(certs): implement fallback strategies for certificate operations
```

**PR Body**: See PR-BODY-4.md

---

### Phase 2: Build & Push Implementation PRs

#### PR #5: Image Builder
**Branch**: `idpbuilder-oci-build-push/phase2/wave1/image-builder`
**Depends On**: All Phase 1 PRs (#1-#4)

**PR Title**:
```
feat(build): implement OCI image builder with go-containerregistry
```

**PR Body**: See PR-BODY-5.md

---

#### PR #6: Gitea Client Implementation
**Branch**: `idpbuilder-oci-build-push/phase2/wave1/gitea-client`
**Depends On**: PR #5

**PR Title**:
```
feat(push): implement Gitea OCI registry client with certificate support
```

**PR Body**: See PR-BODY-6.md

---

## 📊 PR CREATION COMMANDS

### Using GitHub CLI (Recommended)
```bash
# PR #1: Kind Certificate Extraction
gh pr create \
  --base main \
  --head idpbuilder-oci-build-push/phase1/wave1/kind-certificate-extraction \
  --title "feat(certs): implement Kind cluster certificate extraction for Gitea OCI registry" \
  --body-file PR-BODY-1.md

# Continue for all PRs...
```

### Using GitHub Web Interface
1. Navigate to https://github.com/jessesanford/idpbuilder
2. Click "Pull requests" → "New pull request"
3. Select base: `main`, compare: `[effort-branch-name]`
4. Copy title and body from templates
5. Create pull request

## ⚠️ CRITICAL WARNINGS

### DO NOT:
- Skip the merge order - will cause conflicts
- Force merge with failing tests
- Merge multiple PRs simultaneously in wrong order
- Delete branches before confirming merge success

### DO:
- Follow the exact order specified
- Wait for each PR to fully merge before starting next
- Run verification steps after each merge
- Keep this document for reference

## 📈 PROGRESS TRACKING

Track your progress here:

### Merge Status
- [ ] PR #1: Kind Certificate Extraction - Status: [CREATED/APPROVED/MERGED]
- [ ] PR #2: Registry TLS Trust Integration - Status: [WAITING/CREATED/APPROVED/MERGED]
- [ ] PR #3: Certificate Validation Pipeline - Status: [WAITING/CREATED/APPROVED/MERGED]
- [ ] PR #4: Fallback Strategies - Status: [WAITING/CREATED/APPROVED/MERGED]
- [ ] PR #5: Image Builder - Status: [WAITING/CREATED/APPROVED/MERGED]
- [ ] PR #6: Gitea Client - Status: [WAITING/CREATED/APPROVED/MERGED]

## 🏁 COMPLETION

Once all PRs are merged:
1. Verify main branch builds successfully
2. Run full test suite on main
3. Test complete workflow: extract cert → build image → push to Gitea
4. Verify zero certificate errors during normal operation
5. Tag release if appropriate
6. Archive this PR plan for future reference

---

END OF MASTER PR PLAN

Software Factory 2.0 - Implementation Complete ✅
