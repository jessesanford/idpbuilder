## Summary
Implements certificate extraction from Kind cluster to enable secure communication with the builtin Gitea OCI registry. This is the foundation for solving the self-signed certificate problem that prevents reliable OCI image operations.

## Changes
- Added `KindCertExtractor` interface for certificate extraction
- Implemented detection of Kind cluster and Gitea pod
- Added certificate copying from Gitea pod (`/data/gitea/https/cert.pem`)
- Implemented local trust store management at `~/.idpbuilder/certs/`
- Added comprehensive error handling for missing clusters/pods

## Testing
- Unit tests: ✅ All passing with mocked kubectl commands
- Integration tests: ✅ Verified in project-integration branch
- Build status: ✅ Successful
- Test coverage: 85%

## Breaking Changes
None - This is new functionality

## Related Issues
Addresses the self-signed certificate issue preventing OCI registry operations

## Verification
After merging, verify:
1. Build passes on main
2. Certificate extraction API is available
3. No conflicts with subsequent PRs
