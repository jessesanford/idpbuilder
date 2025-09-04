# E2.1.1 Split-002: Certificate Management

## Split Details
**Size Target**: 600 lines ✅ (Under 700 soft limit)
**Branch**: idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-002
**Base Branch**: idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-001

## Files to Implement
```
pkg/certs/errors.go          - 41 lines
pkg/certs/types.go           - 217 lines
pkg/certs/trust_store.go     - 172 lines
pkg/certs/validator.go       - 167 lines
pkg/certs/transport.go       - 239 lines (partial)
pkg/util/                    - Subset of utilities
```

## Dependencies
- Split 001 (uses core types from pkg/builder/)
- Feature Flag: ENABLE_CERT_MANAGEMENT=true

## Implementation Notes
- This is SEQUENTIAL implementation (after Split 001)
- Focus on certificate management functionality
- Must integrate with builder types from split-001

## 🚨 SPLIT INFRASTRUCTURE METADATA (Added by Orchestrator)
**WORKING_DIRECTORY**: /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase2/wave1/go-containerregistry-image-builder-SPLIT-002
**BRANCH**: idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-002
**REMOTE**: origin/idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-002
**BASE_BRANCH**: idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-001
**SPLIT_NUMBER**: 002
**CREATED_AT**: 2025-09-04 04:13:00

### SW Engineer Instructions
1. READ this metadata FIRST
2. cd to WORKING_DIRECTORY above
3. Verify branch matches BRANCH above
4. ONLY THEN proceed with implementation
5. Use feature flag ENABLE_CERT_MANAGEMENT
6. Size limit: 600 lines
