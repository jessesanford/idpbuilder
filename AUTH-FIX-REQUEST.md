# ERROR RECOVERY: Authentication Bug Fix Request

## 🔴 CRITICAL BUG IDENTIFIED

**Date**: 2025-09-16
**Severity**: CRITICAL - Blocks core functionality
**State Transition**: WAITING_FOR_PHASE_ASSESSMENT → ERROR_RECOVERY

## Bug Description

The `idpbuilder push` command accepts `--username` and `--token` flags but **does not actually use them** when authenticating with the registry. This is a **stub implementation** that violates R320 (no stub implementations).

## Location

**File**: `pkg/registry/push.go`
**Lines**: 93-96
**Effort Origin**: `phase2/wave1/gitea-client-split-002`

```go
// Add authentication if available
if authHeader, err := r.authMgr.GetAuthHeader(ctx); err == nil && authHeader != "" {
    // Note: We'd need to convert the authHeader to proper remote.Option
    // For now, just use empty options - this is a simplified implementation
}
```

## Root Cause

The authentication header is retrieved but never converted to `remote.Option` for the `remote.Write()` call. The code has a TODO comment admitting it's "simplified" (stub) implementation.

## Required Fix

Per R321 (immediate backport), this fix must be applied to the **source effort branch** before re-integration.

### Fix Implementation

```go
// pkg/registry/push.go lines 89-100 should be:

// Create remote options with authentication
options := []remote.Option{}

// Add authentication if available
if authHeader, err := r.authMgr.GetAuthHeader(ctx); err == nil && authHeader != "" {
    // Parse the auth header to get credentials
    username := r.config.Username
    password := r.config.Token

    if username != "" && password != "" {
        auth := &authn.Basic{
            Username: username,
            Password: password,
        }
        options = append(options, remote.WithAuth(auth))
    }
}

// For insecure registries
if r.config.Insecure {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    options = append(options, remote.WithTransport(tr))
}
```

## Affected Branches

1. **Source Branch**: `idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002`
2. **Wave Integration**: `idpbuilder-oci-build-push/phase2/wave1/integration-*`
3. **Phase Integration**: `idpbuilder-oci-build-push/phase2-integration-*`

## Recovery Process (R321 Compliance)

1. **Fix at Source**: Apply fix to `gitea-client-split-002` branch
2. **Re-test**: Verify push works with credentials
3. **Re-integrate Wave 1**: Merge fixed branch into wave integration
4. **Re-integrate Wave 2**: Rebase Wave 2 onto fixed Wave 1
5. **Re-integrate Phase**: Create new phase integration with fix
6. **Resume Assessment**: Return to WAITING_FOR_PHASE_ASSESSMENT

## Validation Test

```bash
# This test MUST pass after the fix
docker build -t gitea.cnoe.localtest.me:8443/test:v1 .

idpbuilder push \
    --username giteaadmin \
    --token f2416711c42be7670a1034ab3dfb77e1ea5d9810 \
    --insecure \
    gitea.cnoe.localtest.me:8443/test:v1

# Should succeed without "no authentication credentials available" error
```

## Grading Impact

- **R320 Violation**: Stub implementation found (-50% penalty)
- **R321 Compliance**: Must fix at source, not in integration
- **R323 Impact**: Final artifact doesn't work as intended

## Next Steps

1. Spawn SW Engineer to fix `gitea-client-split-002`
2. Spawn Code Reviewer to verify fix
3. Spawn Integration Agent to re-integrate
4. Resume phase assessment after fix verified