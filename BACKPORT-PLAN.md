# Comprehensive Backport Plan for Integration Fixes

## 🔴 R321 Compliance: Immediate Backport to Source Branches

**Created By**: Code Reviewer Agent
**Date**: 2025-09-12T20:22:00Z
**Purpose**: Fix all integration issues in source branches per R321 protocol

## Executive Summary

This plan addresses THREE critical issues discovered during Phase 1 integration:
1. **KIND Test Compilation Error** (Wave 1) - Blocks all test execution
2. **Certificate Test Setup Failures** (Wave 2) - Prevents cert validation
3. **Duplicate Test Helpers** (Wave 1) - Already documented, needs completion

All fixes MUST be applied to source effort branches, NOT integration branches, per R321.

---

## PART 1: KIND Test Compilation Fix (HIGHEST PRIORITY)

### Affected Effort: kind-cert-extraction (Wave 1)

**Branch**: `idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction`
**Working Directory**: `efforts/phase1/wave1/kind-cert-extraction`
**Issue**: undefined: types.ContainerListOptions at line 232
**Root Cause**: Docker API types have been reorganized; ContainerListOptions moved to container subpackage

### Implementation Instructions for SW Engineer:

```bash
# Step 1: Navigate to the effort directory
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave1/kind-cert-extraction

# Step 2: Checkout and update the branch
git checkout idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction
git pull origin idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction

# Step 3: Fix the import in cluster_test.go
```

**File to Modify**: `pkg/kind/cluster_test.go`

**EXACT CODE CHANGE**:
```go
// OLD (line 9-11 approximately):
import (
    // ... other imports ...
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    // ... other imports ...
)

// NEW - Add the container types import:
import (
    // ... other imports ...
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"  // ADD THIS LINE
    "github.com/docker/docker/client"
    // ... other imports ...
)
```

**Also modify line 232**:
```go
// OLD (line 232):
func (m *DockerClientMock) ContainerList(ctx context.Context, listOptions types.ContainerListOptions) ([]types.Container, error) {

// NEW:
func (m *DockerClientMock) ContainerList(ctx context.Context, listOptions container.ListOptions) ([]types.Container, error) {
```

### Verification Steps:

```bash
# Step 4: Verify the fix compiles
go build ./pkg/kind/...

# Step 5: Run tests to ensure they pass
go test ./pkg/kind/... -v

# Step 6: Commit and push the fix
git add pkg/kind/cluster_test.go
git commit -m "fix(R321): correct Docker API types import for ContainerListOptions

- Add container subpackage import
- Update ContainerListOptions to container.ListOptions
- Fixes compilation error in cluster_test.go:232"
git push origin idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction

# Step 7: Mark this fix as complete
echo "KIND test compilation fixed at $(date)" > KIND_FIX_COMPLETE.flag
```

---

## PART 2: Certificate Test Setup Fixes (Wave 2)

### Affected Efforts:
1. cert-validation-split-001
2. cert-validation-split-002  
3. cert-validation-split-003
4. fallback-strategies

**Root Cause Analysis**: Test setup failures likely due to:
- Missing test fixtures or data files
- Working directory issues during test execution
- Package initialization problems

### Fix for cert-validation-split-001:

**Branch**: `idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-001`
**Working Directory**: `efforts/phase1/wave2/cert-validation-split-001`

```bash
# Step 1: Navigate and checkout
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/cert-validation-split-001
git checkout idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-001
git pull origin idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-001

# Step 2: Check for test data requirements
find . -name "testdata" -type d
ls -la pkg/certs/

# Step 3: Create missing test fixtures if needed
mkdir -p pkg/certs/testdata
```

**Create test fixture file**: `pkg/certs/testdata/test-cert.pem`
```
-----BEGIN CERTIFICATE-----
MIIDXTCCAkWgAwIBAgIJAKLdQVPy90WjMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV
BAYTAlVTMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMjQwMTAxMDAwMDAwWhcNMjUwMTAxMDAwMDAwWjBF
MQswCQYDVQQGEwJVUzETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAx1xYoqhvO2bMDDFjE3WhI1dbsH+2mjPpDo9CdEozan3H9dfwK7ZvHdFD
X3B3M7z3Vr6C2TDpT9WvCaIvVzaBNdGXq1ofQk3YT9pbLCsRaFonfbhoyZT23nLl
5DhCZ3qS0LZkwDcFfy8TnqRd8bAW5FKnVQdVtXL8OodI8D5ePXL23brGDHDc7JXg
HgTdvRW8vt6J4K3BvQrXQV5qp/zlQNXQHwO3V6OkVKlpKBXs5t5Q2dPj3P5XvJwd
xL9Qpf3Zh6EJ2M5ATlF8qIrT4Z5PD4KqV3Dt4L8E5Q9wqO5Lzbmr9Y1wQZDM5QfH
qwIDAQABo1AwTjAdBgNVHQ4EFgQUqO3XKL3SJhE4FlDwQKF7LfFhTHAwHwYDVR0j
BBgwFoAUqO3XKL3SJhE4FlDwQKF7LfFhTHAwDAYDVR0TBAUwAwEB/zANBgkqhkiG
9w0BAQsFAAOCAQEAKvWq37S8BqkEHWDgaGFvHsXgJt9Wc9HQHV3L8OqjxK7gdG/3
-----END CERTIFICATE-----
```

**Fix test initialization**: Update any test files with setup issues:

```go
// If tests are failing due to working directory, add to test files:
func TestMain(m *testing.M) {
    // Ensure we're in the right directory for test data
    _, filename, _, _ := runtime.Caller(0)
    dir := path.Dir(filename)
    os.Chdir(dir)
    
    // Run tests
    os.Exit(m.Run())
}
```

### Verification and Commit for cert-validation-split-001:

```bash
# Step 4: Test the fixes
go test ./pkg/certs/... -v

# Step 5: If tests pass, commit
git add -A
git commit -m "fix(R321): resolve certificate test setup issues

- Add missing test fixtures in testdata directory
- Fix working directory issues in test initialization
- Ensure tests can locate required test data"
git push origin idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-001
```

### Fix for cert-validation-split-002:

**Branch**: `idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-002`
**Working Directory**: `efforts/phase1/wave2/cert-validation-split-002`

```bash
# Repeat similar process
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/cert-validation-split-002
git checkout idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-002
git pull origin idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-002

# Apply same test fixture and initialization fixes
mkdir -p pkg/certvalidation/testdata
# Copy test certificates if needed

# Test and commit
go test ./pkg/certvalidation/... -v
git add -A
git commit -m "fix(R321): resolve certificate validation test setup

- Add test fixtures for validation tests
- Fix test initialization issues"
git push origin idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-002
```

### Fix for cert-validation-split-003:

**Branch**: `idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-003`
**Working Directory**: `efforts/phase1/wave2/cert-validation-split-003`

```bash
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/cert-validation-split-003
git checkout idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-003
git pull origin idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-003

# Apply fixes for split-003 specific tests
# Test and commit
go test ./... -v
git add -A
git commit -m "fix(R321): resolve test setup in split-003

- Ensure all test dependencies available
- Fix any initialization issues"
git push origin idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-003
```

### Fix for fallback-strategies:

**Branch**: `idpbuilder-oci-build-push/phase1/wave2/fallback-strategies`
**Working Directory**: `efforts/phase1/wave2/fallback-strategies`

```bash
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/fallback-strategies
git checkout idpbuilder-oci-build-push/phase1/wave2/fallback-strategies
git pull origin idpbuilder-oci-build-push/phase1/wave2/fallback-strategies

# Check for test issues in fallback and insecure packages
go test ./pkg/fallback/... -v
go test ./pkg/insecure/... -v

# Fix any test setup issues found
# Common fix for fallback tests might be ensuring mock data:
mkdir -p pkg/fallback/testdata
mkdir -p pkg/insecure/testdata

# Commit fixes
git add -A
git commit -m "fix(R321): resolve fallback strategy test setup

- Add test fixtures for fallback scenarios
- Fix insecure mode test initialization"
git push origin idpbuilder-oci-build-push/phase1/wave2/fallback-strategies
```

---

## PART 3: Complete Wave 1 Duplicate Test Helpers Fix

### Status Check Required:

The file `R321-IMMEDIATE-BACKPORT-WAVE1-FIX-PLAN.md` indicates this was already planned. Check if `R321_FIXES_COMPLETE.flag` exists:

```bash
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave1/
if [ -f "R321_FIXES_COMPLETE.flag" ]; then
    echo "Wave 1 duplicate helpers already fixed"
    cat R321_FIXES_COMPLETE.flag
else
    echo "Wave 1 fixes still needed - follow existing plan"
fi
```

If not complete, follow the existing plan in `R321-IMMEDIATE-BACKPORT-WAVE1-FIX-PLAN.md`.

---

## Implementation Order and Dependencies

### CRITICAL SEQUENCING:

1. **FIRST**: Fix KIND test compilation (blocks ALL testing)
2. **SECOND**: Fix Wave 1 duplicate helpers (if not already done)
3. **THIRD**: Fix Wave 2 certificate test setups
4. **PARALLEL**: Wave 2 splits can be fixed simultaneously by different engineers

### Dependency Matrix:

| Fix | Depends On | Can Parallelize With | Priority |
|-----|-----------|---------------------|----------|
| KIND compilation | None | None | HIGHEST |
| Wave 1 duplicates | KIND fix | None | HIGH |
| cert-split-001 | KIND fix | Other Wave 2 | MEDIUM |
| cert-split-002 | KIND fix | Other Wave 2 | MEDIUM |
| cert-split-003 | KIND fix | Other Wave 2 | MEDIUM |
| fallback-strategies | KIND fix | Other Wave 2 | MEDIUM |

---

## Verification Protocol

### After ALL Fixes Complete:

```bash
# Step 1: Verify each effort builds independently
for wave in wave1 wave2; do
    for effort in $(ls efforts/phase1/$wave/); do
        if [ -d "efforts/phase1/$wave/$effort" ] && [ "$effort" != "*-workspace" ]; then
            echo "Testing $wave/$effort..."
            cd efforts/phase1/$wave/$effort
            go build ./...
            go test ./... -v
            if [ $? -ne 0 ]; then
                echo "FAILED: $wave/$effort"
                exit 1
            fi
        fi
    done
done

# Step 2: Create completion marker
echo "All backports complete at $(date)" > BACKPORTS_COMPLETE.flag
echo "Ready for re-integration" >> BACKPORTS_COMPLETE.flag
```

---

## Risk Assessment and Mitigations

### Risks:
1. **Test Data Dependencies**: Some tests may require specific test data files
   - **Mitigation**: Create minimal valid test fixtures as shown above

2. **Import Conflicts**: Docker API changes may affect other files
   - **Mitigation**: Search for all uses of ContainerListOptions globally

3. **Sequential Dependencies**: Later splits may depend on earlier ones
   - **Mitigation**: Fix splits in order (001, 002, 003)

### Global Search Commands:

```bash
# Find all potential Docker API issues
grep -r "types.ContainerListOptions" efforts/
grep -r "ContainerList" efforts/ --include="*.go"

# Find all test setup functions that might fail
grep -r "TestMain\|func setup" efforts/ --include="*_test.go"

# Find missing testdata directories
for dir in efforts/phase1/*/*/pkg/*/; do
    if ls $dir/*_test.go 2>/dev/null | grep -q test; then
        if [ ! -d "$dir/testdata" ]; then
            echo "Missing testdata in: $dir"
        fi
    fi
done
```

---

## Success Criteria Checklist

Before declaring backports complete:

- [ ] KIND test compilation error fixed in kind-cert-extraction
- [ ] All Wave 1 duplicate test helpers resolved
- [ ] cert-validation-split-001 tests pass
- [ ] cert-validation-split-002 tests pass
- [ ] cert-validation-split-003 tests pass
- [ ] fallback-strategies tests pass
- [ ] Each effort branch builds independently (`go build ./...`)
- [ ] Each effort branch tests pass independently (`go test ./...`)
- [ ] All fixes committed to source branches
- [ ] All fixes pushed to remote
- [ ] NO fixes exist only in integration branches
- [ ] BACKPORTS_COMPLETE.flag created

---

## Communication Template for SW Engineers

When assigning fixes to SW Engineers, use this template:

```
BACKPORT ASSIGNMENT: [Effort Name]
Branch: [exact branch name]
Priority: [HIGHEST/HIGH/MEDIUM]
Issue: [brief description]
Fix Location: [exact file and line]
Verification: go test ./... must pass
Deadline: IMMEDIATE per R321
See BACKPORT-PLAN.md section: [PART X]
```

---

## Notes for Orchestrator

1. **Spawn Strategy**: 
   - ONE engineer for KIND fix (highest priority)
   - PARALLEL engineers for Wave 2 fixes (after KIND complete)
   - Monitor R321_FIXES_COMPLETE.flag for Wave 1 status

2. **Validation Gates**:
   - Do NOT proceed to re-integration until ALL fixes verified
   - Each effort must pass tests independently
   - Use this plan's success criteria as gate conditions

3. **R321 Compliance**:
   - All fixes go to source branches only
   - Integration branches remain untouched
   - Re-integration happens after fixes complete

---

*End of Backport Plan*
*Generated by Code Reviewer Agent*
*R321 Compliant*