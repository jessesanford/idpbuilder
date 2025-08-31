# Integration Plan - Phase 1 Wave 1 (RE-RUN)
Date: 2025-08-31 17:15:00 UTC
Target Branch: idpbuidler-oci-go-cr/phase1/wave1/integration-v2-20250831-171415
Integration Type: Clean merge after duplicate type fixes

## Context
This is a RE-RUN of the integration after fixing duplicate types:
- E1.1.2 had duplicate CertificateInfo struct (now fixed)
- E1.1.2 now properly imports types from E1.1.1
- Clean integration expected

## Branches to Integrate (ordered by dependency)
1. E1.1.1: kind-certificate-extraction (defines base types in pkg/certs/types.go)
   - Source: /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/wave1/kind-certificate-extraction/
   - Branch: feature/extract-kind-certificates
   - Contains: Certificate extraction functionality and base types

2. E1.1.2: registry-tls-trust-integration (imports types from E1.1.1)
   - Source: /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/wave1/registry-tls-trust-integration/
   - Branch: feature/registry-tls-trust
   - Contains: Registry TLS trust functionality (using imported types)

## Merge Strategy
1. Add both effort directories as git remotes
2. Fetch all branches from both remotes
3. Merge E1.1.1 FIRST (contains base type definitions)
4. Merge E1.1.2 SECOND (depends on types from E1.1.1)
5. Copy pkg/certs directories from both efforts
6. Verify no duplicate types exist
7. Run build and tests

## Expected Outcome
- Fully integrated branch with both features
- No duplicate type definitions
- All tests passing
- Clean build with no compilation errors
- Complete documentation of integration process