## Summary
Implements TLS trust store management to configure go-containerregistry with custom CA certificates extracted from Kind/Gitea. This enables secure HTTPS communication with the self-signed Gitea registry.

## Changes
- Added `TrustStoreManager` interface for certificate trust management
- Implemented custom CA loading into x509.CertPool
- Configured go-containerregistry remote transport with TLS
- Added certificate rotation support
- Implemented `--insecure` registry override option
- Added registry-specific certificate management

## Testing
- Unit tests: ✅ All passing including CA pool loading tests
- Integration tests: ✅ Verified in project-integration branch
- Build status: ✅ Successful
- Test coverage: 82%

## Dependencies
- Requires PR #1 (certificate extraction) to be merged first

## Breaking Changes
None - This is new functionality

## Verification
After merging, verify:
1. Build still passes
2. Integration with certificate extraction works
3. TLS configuration is properly applied to registry operations
