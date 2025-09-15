# Duplicate Declarations Fix Plan

## Issue Summary
- **Problem**: Duplicate type declarations preventing Phase 2 Wave 1 integration build
- **Location**: pkg/certs package has duplicate TLSConfig struct and DefaultTLSConfig function
- **Severity**: CRITICAL - Blocks all builds
- **Discovery**: Phase 2 Wave 1 integration exposed latent Phase 1 bug

## Root Cause Analysis

### Finding
The duplicate declarations appear in the Phase 2 Wave 1 integration workspace but NOT in the original Phase 1 effort branches. This indicates the duplicates were introduced during the integration process itself.

### Evidence
1. **Phase 2 Wave 1 Integration** (HAS DUPLICATES):
   - `pkg/certs/utilities.go:130-139`: Contains TLSConfig struct and DefaultTLSConfig function
   - `pkg/certs/types.go:55`: Contains TLSConfig struct (more comprehensive)
   - `pkg/certs/types.go:150`: Contains DefaultTLSConfig function

2. **Phase 1 Integration Worktree** (NO DUPLICATES):
   - Already fixed, utilities.go has no TLSConfig declarations

3. **Original Effort Branches** (NO DUPLICATES):
   - registry-tls-trust has utilities.go but no TLSConfig in it
   - registry-auth-types-split-002 introduced types.go with TLSConfig

### Root Cause
The duplicates were likely introduced during a merge conflict resolution where someone accidentally kept both versions instead of removing the duplicate from utilities.go.

## Files Requiring Fixes

### File 1: pkg/certs/utilities.go
**Location**: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/integration-workspace/repo/pkg/certs/utilities.go`

**Changes Required**:
1. **REMOVE lines 130-139** (entire TLSConfig struct declaration)
2. **REMOVE lines 139-145** (entire DefaultTLSConfig function)

**Lines to Remove**:
```go
// TLSConfig holds TLS configuration
type TLSConfig struct {
	Registry           string
	InsecureSkipVerify bool
	MinVersion         uint16
	ValidateHostname   bool
	Timeout            time.Duration
}

// DefaultTLSConfig returns default TLS config
func DefaultTLSConfig() *TLSConfig {
	return &TLSConfig{
		MinVersion:       tls.VersionTLS12,
		ValidateHostname: true,
		Timeout:          10 * time.Second,
	}
}
```

### File 2: pkg/certs/types.go
**Location**: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/integration-workspace/repo/pkg/certs/types.go`

**Changes Required**: NONE - Keep all declarations in this file as they are more comprehensive

## Fix Application Strategy

### Per R300 Compliance
Since the duplicates don't exist in the original effort branches but only in the integration, the fix must be applied to the integration branch where the problem exists.

### Target Branch
**Branch to Fix**: `idpbuilder-oci-build-push/phase2/wave1/integration`
**Location**: Phase 2 Wave 1 integration workspace

### Execution Steps

1. **Navigate to Integration Workspace**
   ```bash
   cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/integration-workspace/repo
   ```

2. **Remove Duplicate Declarations from utilities.go**
   - Edit pkg/certs/utilities.go
   - Delete lines 130-145 (TLSConfig struct and DefaultTLSConfig function)
   - Keep all other content intact

3. **Verify Fix**
   ```bash
   # Attempt compilation
   cd pkg/certs
   go build

   # Run tests
   go test ./...
   ```

4. **Commit Fix**
   ```bash
   git add pkg/certs/utilities.go
   git commit -m "fix: remove duplicate TLSConfig declarations from utilities.go

   - Removed duplicate TLSConfig struct (was in both utilities.go and types.go)
   - Removed duplicate DefaultTLSConfig function (was in both files)
   - Kept comprehensive versions in types.go only
   - Fixes Phase 2 Wave 1 integration build failure"
   ```

## Verification Steps

1. **Build Verification**
   ```bash
   # From integration workspace root
   go build ./...
   ```
   Expected: Build succeeds without duplicate declaration errors

2. **Test Verification**
   ```bash
   go test ./pkg/certs/...
   go test ./pkg/certvalidation/...
   ```
   Expected: All tests pass

3. **Demo Verification**
   ```bash
   ./demo-image-builder.sh status
   ./demo-gitea-client.sh auth
   ```
   Expected: Demos execute successfully

## Prevention Measures

### For Future Integrations
1. When resolving merge conflicts, always check for duplicate type declarations
2. Run `go build ./...` after every merge to catch compilation errors immediately
3. When seeing duplicate declaration errors, remove from utilities files, keep in types files

### R327 Cascade Implications
After this fix is applied:
1. Phase 2 Wave 1 integration can complete successfully
2. No need to cascade back to Phase 1 as the issue only exists in Phase 2 integration
3. Can proceed with Phase 2 Wave 2 planning

## Summary

The duplicate declarations are a merge artifact that only exists in the Phase 2 Wave 1 integration branch. The fix is straightforward: remove the duplicate declarations from utilities.go while keeping the comprehensive versions in types.go. This will unblock the Phase 2 Wave 1 build and allow the project to proceed.

**Action Required**: SW Engineer should apply this fix to the Phase 2 Wave 1 integration branch immediately.