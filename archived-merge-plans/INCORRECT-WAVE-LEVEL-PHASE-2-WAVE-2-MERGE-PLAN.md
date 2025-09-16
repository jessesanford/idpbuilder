# Phase 2 Wave 2 Comprehensive Merge Plan

**Created**: 2025-09-16 00:45:00 UTC
**Purpose**: Prevent API mismatch issues during Wave 2 re-integration
**Critical**: This plan addresses the failures from previous integration attempts

## Executive Summary

The previous Wave 2 integration failed due to incorrect conflict resolution that introduced old API calls. This comprehensive plan ensures the Integration Agent knows EXACTLY which code to keep during conflicts.

## 🔴🔴🔴 CRITICAL API COMPATIBILITY REQUIREMENTS 🔴🔴🔴

### ✅ REQUIRED API Signatures (MUST be preserved)

These are the ONLY valid API signatures. ANY deviation must be rejected:

```go
// Certificate Management - NEW API
certs.NewTrustStore() *DefaultTrustStore                    // ✅ CORRECT
trustStore.ValidateCertificate(cert *x509.Certificate) bool // ✅ CORRECT

// Gitea Client - NEW API
gitea.NewClient(registryURL string, certManager *certs.DefaultTrustStore) (*Client, error) // ✅ CORRECT - 2 params
client.Push(imageRef string, progressChan chan<- PushProgress) error                       // ✅ CORRECT method name

// Image Operations
operations.BuildAndPush(ctx context.Context, opts BuildOptions) error // ✅ CORRECT
```

### ❌ FORBIDDEN API Calls (AUTOMATIC REJECTION)

If ANY of these appear during merge, REJECT IMMEDIATELY:

```go
// OLD/WRONG Certificate API
certs.ExtractCertificate()           // ❌ OLD API - NEVER ACCEPT
certs.GetSystemCerts()               // ❌ DOES NOT EXIST
trustStore.ExtractCertificate()      // ❌ WRONG - use ValidateCertificate

// OLD/WRONG Gitea API
gitea.NewClientAutoDetect()          // ❌ DOES NOT EXIST - REJECT
gitea.NewClient(url, user, pass, certManager) // ❌ 4 params = OLD - REJECT
client.ValidateCredentials()         // ❌ DOES NOT EXIST - REJECT
client.PushImage()                   // ❌ OLD NAME - use Push()
client.Login()                       // ❌ DOES NOT EXIST - REJECT

// OLD/WRONG Image Operations
operations.PushToRegistry()          // ❌ OLD NAME - use BuildAndPush
```

## 📋 Pre-Integration Verification

### Step 1: Verify Base Branch
```bash
# Ensure we're starting from correct base
cd /home/vscode/workspaces/idpbuilder-oci-build-push
git checkout phase2/wave1/integration-20250915-125755
git log --oneline -5  # Should show Wave 1 integration commits
```

### Step 2: Verify Effort Branches
```bash
# Check each effort branch has correct API
for effort in cli-commands credential-management image-operations; do
    echo "Checking $effort..."
    git checkout phase2/wave2/$effort

    # Verify NO old API calls exist
    ! grep -r "ExtractCertificate" pkg/ || echo "ERROR: Old API in $effort"
    ! grep -r "NewClientAutoDetect" pkg/ || echo "ERROR: Old API in $effort"
    ! grep -r "PushImage" pkg/ || echo "ERROR: Old API in $effort"

    # Verify new API exists
    grep -r "NewTrustStore" pkg/ || echo "ERROR: Missing new API in $effort"
    grep -r "Push(" pkg/gitea/ || echo "ERROR: Missing Push method in $effort"
done
```

## 🔄 Merge Sequence and Strategy

### Integration Branch Setup
```bash
# Create new integration branch from Wave 1 integration
git checkout phase2/wave1/integration-20250915-125755
git checkout -b phase2/wave2/integration-$(date +%Y%m%d-%H%M%S)
git push -u origin HEAD
```

### Merge Order (CRITICAL - Based on Dependencies)

#### 1️⃣ FIRST: cli-commands
**Why First**: Establishes CLI structure that others depend on
```bash
git merge phase2/wave2/cli-commands

# If conflicts in pkg/cmd/push.go:
# ACCEPT version with:
# - certs.NewTrustStore() NOT ExtractCertificate()
# - gitea.NewClient(url, certManager) NOT 4 params
# - client.Push() NOT PushImage()
```

**Post-Merge Verification**:
```bash
# Must pass ALL checks
grep "NewTrustStore" pkg/cmd/push.go || exit 1
! grep "ExtractCertificate" pkg/cmd/push.go || exit 1
! grep "NewClientAutoDetect" pkg/cmd/push.go || exit 1
go build ./pkg/cmd/... || exit 1
```

#### 2️⃣ SECOND: credential-management
**Why Second**: Adds authentication layer on top of CLI
```bash
git merge phase2/wave2/credential-management

# If conflicts in pkg/gitea/client.go:
# ACCEPT version with:
# - NewClient with 2 parameters
# - Push() method, NOT PushImage()
# - No ValidateCredentials() method
```

**Post-Merge Verification**:
```bash
# Must pass ALL checks
grep "func NewClient(registryURL string, certManager" pkg/gitea/client.go || exit 1
grep "func.*Push(" pkg/gitea/client.go || exit 1
! grep "func.*PushImage" pkg/gitea/client.go || exit 1
! grep "ValidateCredentials" pkg/gitea/client.go || exit 1
go build ./pkg/gitea/... || exit 1
```

#### 3️⃣ THIRD: image-operations
**Why Third**: Completes functionality using previous components
```bash
git merge phase2/wave2/image-operations

# If conflicts in pkg/operations/:
# ACCEPT version with:
# - BuildAndPush() NOT PushToRegistry()
# - Uses gitea.NewClient with 2 params
# - Uses client.Push() NOT PushImage()
```

**Post-Merge Verification**:
```bash
# Must pass ALL checks
grep "BuildAndPush" pkg/operations/*.go || exit 1
! grep "PushToRegistry" pkg/operations/*.go || exit 1
go build ./pkg/operations/... || exit 1
```

## 🔍 Conflict Resolution Decision Tree

### For ANY Conflict

```
CONFLICT DETECTED
    |
    v
Does conflicting code contain old API?
    |
    ├─ YES: ExtractCertificate, NewClientAutoDetect, PushImage, etc.
    |   └─> REJECT that version immediately
    |
    └─ NO: Both versions use new API
        |
        v
    Check effort branch directly:
        git show phase2/wave2/[effort]:path/to/file
        Use EXACTLY what's in effort branch
```

### Specific File Conflict Rules

#### pkg/cmd/push.go Conflicts
```go
// ✅ KEEP THIS VERSION:
trustStore := certs.NewTrustStore()
client, err := gitea.NewClient(registryURL, trustStore)
err = client.Push(imageRef, progressChan)

// ❌ REJECT THIS VERSION:
certManager := certs.ExtractCertificate()
client := gitea.NewClientAutoDetect(url, user, pass, certManager)
err = client.PushImage(imageRef)
```

#### pkg/gitea/client.go Conflicts
```go
// ✅ KEEP THIS VERSION:
func NewClient(registryURL string, certManager *certs.DefaultTrustStore) (*Client, error)
func (c *Client) Push(imageRef string, progressChan chan<- PushProgress) error

// ❌ REJECT THIS VERSION:
func NewClient(url, user, pass string, certManager CertificateManager) (*Client, error)
func (c *Client) PushImage(imageRef string) error
func (c *Client) ValidateCredentials() error  // Should not exist
```

## 🧪 Integration Testing Protocol

### After EACH Merge
```bash
# 1. Compile check
go build ./... || { echo "BUILD FAILED"; exit 1; }

# 2. API compatibility check
./scripts/verify-api-compatibility.sh || { echo "API CHECK FAILED"; exit 1; }

# 3. Test execution
go test ./... || { echo "TESTS FAILED"; exit 1; }

# 4. Commit the merge
git add -A
git commit -m "integrate: merge [effort-name] into Wave 2 integration"
git push
```

### Final Integration Validation
```bash
# After ALL merges complete
cd /home/vscode/workspaces/idpbuilder-oci-build-push

# 1. Full build
make clean
make build || exit 1

# 2. API verification script
cat > verify-final-api.sh << 'EOF'
#!/bin/bash
echo "Verifying final API compatibility..."

# Check for forbidden old API calls
ERRORS=0

if grep -r "ExtractCertificate" pkg/; then
    echo "ERROR: Old ExtractCertificate API found!"
    ERRORS=$((ERRORS + 1))
fi

if grep -r "NewClientAutoDetect" pkg/; then
    echo "ERROR: Non-existent NewClientAutoDetect found!"
    ERRORS=$((ERRORS + 1))
fi

if grep -r "PushImage" pkg/; then
    echo "ERROR: Old PushImage method found!"
    ERRORS=$((ERRORS + 1))
fi

if grep -r "ValidateCredentials" pkg/; then
    echo "ERROR: Non-existent ValidateCredentials found!"
    ERRORS=$((ERRORS + 1))
fi

# Check for required new API
if ! grep -r "NewTrustStore" pkg/; then
    echo "ERROR: NewTrustStore API missing!"
    ERRORS=$((ERRORS + 1))
fi

if ! grep -r "Push(" pkg/gitea/; then
    echo "ERROR: Push method missing!"
    ERRORS=$((ERRORS + 1))
fi

if [ $ERRORS -eq 0 ]; then
    echo "✅ All API checks passed!"
else
    echo "❌ Found $ERRORS API compatibility issues"
    exit 1
fi
EOF

chmod +x verify-final-api.sh
./verify-final-api.sh || exit 1

# 3. Run full test suite
go test -v ./...

# 4. Demo execution
./demo-features.sh || echo "Demo needs adjustment"
```

## 📝 Integration Agent Instructions

### CRITICAL REMINDERS FOR INTEGRATION AGENT

1. **NEVER GUESS** during conflicts - check this plan
2. **ALWAYS VERIFY** API after each merge using provided commands
3. **IF UNCERTAIN** about a conflict:
   - Check the effort branch directly: `git show phase2/wave2/[effort]:path/to/file`
   - Use EXACTLY what's in the effort branch
4. **STOP IMMEDIATELY** if you see old API calls being merged
5. **DO NOT PROCEED** to next merge if verification fails

### Conflict Resolution Checklist
- [ ] Identified which version has old API calls
- [ ] Rejected ALL old API code
- [ ] Verified new API is preserved
- [ ] Ran compilation check
- [ ] Ran API verification
- [ ] Tests pass

## 🚨 Emergency Recovery

If integration fails:
```bash
# Reset to base and start over
git checkout phase2/wave1/integration-20250915-125755
git branch -D phase2/wave2/integration-[timestamp]
git push origin --delete phase2/wave2/integration-[timestamp]

# Start fresh with this plan
```

## 📊 Success Criteria

Integration is ONLY successful when:
- ✅ All three efforts merged without old API contamination
- ✅ `go build ./...` succeeds
- ✅ `go test ./...` passes
- ✅ No ExtractCertificate() calls exist
- ✅ No NewClientAutoDetect() calls exist
- ✅ No PushImage() or ValidateCredentials() methods exist
- ✅ NewTrustStore() is used for certificates
- ✅ gitea.NewClient has exactly 2 parameters
- ✅ Push() method exists and works

## 🎯 Expected Final State

After successful integration, the codebase should have:
1. **CLI**: Working `idpbuilder oci push` command
2. **Auth**: Proper credential management with correct API
3. **Operations**: Complete build and push functionality
4. **Certificates**: Using NewTrustStore() throughout
5. **Gitea Client**: Using 2-parameter NewClient and Push() method

---

**END OF MERGE PLAN**

This plan is designed to be foolproof. Follow it exactly and the integration will succeed.