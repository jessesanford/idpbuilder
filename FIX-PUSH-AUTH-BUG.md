# Push Authentication Bug Fix

## Bug Description
The `idpbuilder push` command was failing with "no authentication credentials available" even when `--username` and `--token` flags were provided.

## Root Cause
The bug was in the integration-level code (`pkg/gitea/client.go`):
- When `NewInsecureClient()` creates a `GiteaRegistry`, it initializes an `AuthManager` with empty credentials
- Later, `SetCredentials()` updates the config and credential manager BUT doesn't update the already-created registry's `authMgr`
- Result: The registry continues using empty credentials despite CLI flags being provided

## Fix Applied
1. Added `UpdateCredentials()` method to `GiteaRegistry` in `pkg/registry/gitea.go`
2. Modified `SetCredentials()` in `pkg/gitea/client.go` to call `UpdateCredentials()` on the registry

## Files Modified
- `pkg/gitea/client.go` - Added registry auth update in SetCredentials
- `pkg/registry/gitea.go` - Added UpdateCredentials method

## Testing
After applying the fix, the push command should work with:
```bash
./idpbuilder push sample-app:v1.0 \
  --insecure \
  --username giteaadmin \
  --token <your-token>
```

## Branch
- Fix applied to: `fix-push-auth-credentials`
- Based on: `idpbuilder-oci-build-push/project-integration-cascade-20250919-210602`
- PR URL: https://github.com/jessesanford/idpbuilder/pull/new/fix-push-auth-credentials

## Important Notes
- This is an INTEGRATION-LEVEL fix (pkg/gitea/ doesn't exist in effort branches)
- No backport to effort branches needed
- The pkg/gitea/ package was added during integration as an adapter layer