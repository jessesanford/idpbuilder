# WAVE MERGE PLAN - Phase 1 Wave 1

**Created**: 2025-09-16 19:09 UTC
**Agent**: code-reviewer
**State**: WAVE_MERGE_PLANNING
**Integration Branch**: idpbuilder-oci-build-push/phase1/wave1/integration
**Base Branch**: main
**Target Repository**: https://github.com/jessesanford/idpbuilder.git

## 🚨 CRITICAL RULES COMPLIANCE

This merge plan follows:
- **R269**: Using ONLY original effort branches - NO integration branches
- **R270**: Analyzing branch bases to determine correct merge order
- **R307**: Ensuring independent branch mergeability
- **Excluding**: Original 'registry-auth-types' branch (replaced by splits)
- **Including**: Only split branches for registry-auth-types

## 📊 EFFORT ANALYSIS

### Phase 1 Wave 1 Efforts Identified

| Effort | Branch | Type | Commits | Base | Status |
|--------|--------|------|---------|------|--------|
| kind-cert-extraction | idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction | Original | 9 | main | Ready |
| registry-auth-types | idpbuilder-oci-build-push/phase1/wave1/registry-auth-types | Original (TOO LARGE) | 8 | main | EXCLUDED |
| registry-auth-types-split-001 | idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-001 | Split-001 | 6 | main | Ready |
| registry-auth-types-split-002 | idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-002 | Split-002 | 7 | split-001 | Ready |
| registry-tls-trust | idpbuilder-oci-build-push/phase1/wave1/registry-tls-trust | Original | 8 | main | Ready |

### Branch Base Analysis

```
main (e210954a0aa81afd110ab47e5a6239fd228c5ce0)
├── kind-cert-extraction (based on main)
├── registry-auth-types (EXCLUDED - replaced by splits)
├── registry-auth-types-split-001 (based on main via 67b4b08e)
│   └── registry-auth-types-split-002 (based on split-001)
└── registry-tls-trust (based on main)
```

## 🔀 MERGE SEQUENCE

Based on R270 analysis and dependency chains, the correct merge order is:

### Step 1: Merge kind-cert-extraction
```bash
git merge origin/idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction --no-ff \
  -m "feat(phase1-wave1): merge kind-cert-extraction effort"
```
- **Files Modified**: Certificate extraction utilities in pkg/certs/
- **Conflict Risk**: LOW - New functionality
- **Dependencies**: None

### Step 2: Merge registry-auth-types-split-001
```bash
git merge origin/idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-001 --no-ff \
  -m "feat(phase1-wave1): merge registry-auth-types-split-001"
```
- **Files Modified**: Authentication constants and base types in pkg/auth/
- **Conflict Risk**: MEDIUM - May conflict with kind-cert-extraction in pkg/certs area
- **Dependencies**: None

### Step 3: Merge registry-auth-types-split-002
```bash
git merge origin/idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-002 --no-ff \
  -m "feat(phase1-wave1): merge registry-auth-types-split-002"
```
- **Files Modified**: Certificate types and TLS configuration in pkg/certs/
- **Conflict Risk**: HIGH - Contains test helpers that may conflict with split-001 and kind-cert-extraction
- **Dependencies**: Requires split-001 to be merged first

### Step 4: Merge registry-tls-trust
```bash
git merge origin/idpbuilder-oci-build-push/phase1/wave1/registry-tls-trust --no-ff \
  -m "feat(phase1-wave1): merge registry-tls-trust effort"
```
- **Files Modified**: TLS trust management functionality
- **Conflict Risk**: MEDIUM - May interact with certificate types from splits
- **Dependencies**: Benefits from having cert types merged first

## ⚠️ EXPECTED CONFLICTS

### Likely Conflict Points

1. **pkg/certs/constants.go**
   - Modified by: kind-cert-extraction, registry-auth-types-split-002
   - Resolution: Combine both sets of constants

2. **pkg/certs/types.go**
   - Modified by: registry-auth-types-split-002, potentially registry-tls-trust
   - Resolution: Ensure all types are included, avoid duplicates

3. **pkg/certs/test_helpers.go**
   - Modified by: registry-auth-types-split-002
   - Resolution: Keep consolidated test helpers

4. **Test files**
   - Various test files may conflict
   - Resolution: Merge test cases, ensure no duplicate test names

## 📋 CONFLICT RESOLUTION STRATEGY

### General Principles
1. **Preserve all functionality** - Never delete code from either side
2. **Combine imports** - Merge import blocks carefully
3. **Avoid duplicates** - Check for duplicate constants/types/functions
4. **Test after each merge** - Run `go build ./...` and `go test ./...`

### Specific Resolution Patterns

#### For Constants Files:
```go
// Combine both sets of constants
const (
    // From kind-cert-extraction
    CertExtractorVersion = "1.0.0"

    // From registry-auth-types-split-002
    DefaultCertPath = "/etc/ssl/certs"
    // ... more constants
)
```

#### For Type Definitions:
```go
// Check for duplicates before adding
type TLSConfig struct {
    // Merge fields from both branches if different
    // Keep the most complete version
}
```

## 🧪 VALIDATION STEPS

After each merge:

1. **Compile Check**:
   ```bash
   go build ./...
   ```

2. **Test Execution**:
   ```bash
   go test ./... -v
   ```

3. **Verify No Lost Code**:
   ```bash
   git diff HEAD~1 HEAD --stat
   ```

4. **Check for Unresolved Conflicts**:
   ```bash
   git status --porcelain | grep "^UU"
   ```

## 🎯 SUCCESS CRITERIA

The integration is successful when:
- ✅ All 4 branches merged (excluding the original registry-auth-types)
- ✅ No unresolved conflicts remain
- ✅ Code compiles successfully
- ✅ All tests pass (or known failures documented)
- ✅ No code lost during merges
- ✅ Integration branch pushed to remote

## 🚫 DO NOT MERGE

The following should NOT be merged per R269:
- ❌ idpbuilder-oci-build-push/phase1/wave1/registry-auth-types (original, replaced by splits)
- ❌ Any integration branches from other efforts
- ❌ Any branches not listed in the merge sequence above

## 📝 EXECUTION INSTRUCTIONS

For the Integration Agent:

1. **Ensure clean state**:
   ```bash
   cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave1/integration-workspace/repo
   git status  # Should be clean
   git branch --show-current  # Should be: idpbuilder-oci-build-push/phase1/wave1/integration
   ```

2. **Execute merges in order**:
   - Follow Steps 1-4 exactly
   - Resolve conflicts as they arise
   - Test after each merge

3. **Document any issues**:
   - Create INTEGRATION-ISSUES.md if problems arise
   - Note any test failures
   - Document resolution decisions

4. **Push when complete**:
   ```bash
   git push origin idpbuilder-oci-build-push/phase1/wave1/integration
   ```

## 🔍 MERGE VERIFICATION CHECKLIST

- [ ] Pre-merge: Branch is clean and on correct base
- [ ] Step 1: kind-cert-extraction merged
- [ ] Step 2: registry-auth-types-split-001 merged
- [ ] Step 3: registry-auth-types-split-002 merged
- [ ] Step 4: registry-tls-trust merged
- [ ] Post-merge: All code compiles
- [ ] Post-merge: Tests executed (pass or failures documented)
- [ ] Post-merge: Integration branch pushed

---

**END OF MERGE PLAN**

This plan is ready for execution by the Integration Agent.