# Phase 2 Wave 2 Implementation Plan

## Overview
Phase 2 Wave 2 completes the OCI build and push functionality for idpbuilder, transforming the prototype CLI commands into production-ready features.

## Effort Structure

### E2.2.1: cli-commands (COMPLETE - Prototype)
- **Status**: Implementation Complete, Rebased on Phase 2 Wave 1
- **Branch**: `idpbuilder-oci-build-push/phase2/wave2/cli-commands`
- **Size**: 41 lines (verified)
- **Purpose**: Prototype implementation of build and push commands
- **Known Issues**:
  - Empty credential functions (returns "")
  - Placeholder manifest with fake SHA256 digests
  - Simulated progress reporting
  - TODOs for credential management

### E2.2.2-A: credential-management (PENDING)
- **Status**: Pending Creation
- **Branch**: `idpbuilder-oci-build-push/phase2/wave2/credential-management`
- **Base**: `idpbuilder-oci-build-push/phase2/wave2/cli-commands`
- **Size**: 500 lines
- **Scope**:
  - Real credential management implementation
  - Environment variable support (GITEA_USERNAME, GITEA_PASSWORD)
  - Config file parsing (~/.idpbuilder/config)
  - System keyring integration
  - Remove credential-related TODOs

**New Files**:
- `pkg/gitea/credentials.go` (150 lines)
- `pkg/gitea/credentials_test.go` (80 lines)
- `pkg/gitea/config.go` (70 lines)
- `pkg/gitea/config_test.go` (50 lines)

**Modified Files**:
- `pkg/gitea/client.go` - Update credential functions (~50 lines)
- `pkg/gitea/client_test.go` - Update tests (~50 lines)

### E2.2.2-B: image-operations (PENDING)
- **Status**: Pending Creation
- **Branch**: `idpbuilder-oci-build-push/phase2/wave2/image-operations`
- **Base**: `idpbuilder-oci-build-push/phase2/wave2/credential-management`
- **Size**: 500 lines
- **Scope**:
  - Real OCI image loading
  - Docker daemon API integration
  - Real manifest generation with actual SHA256 digests
  - Real progress tracking
  - Remove ALL placeholders and feature flags

**New Files**:
- `pkg/gitea/image_loader.go` (200 lines)
- `pkg/gitea/image_loader_test.go` (80 lines)
- `pkg/gitea/progress.go` (50 lines)

**Modified Files**:
- `pkg/gitea/client.go` - Replace placeholder code (~100 lines)
- `pkg/cmd/build.go` - Remove feature flags (~20 lines)
- `pkg/cmd/push.go` - Update for real operations (~20 lines)

**Files to Delete**:
- `pkg/build/feature_flags.go`

## Implementation Sequence

1. **E2.2.2-A Implementation**
   - Create effort infrastructure
   - Plan credential management approach
   - Implement real authentication
   - Test with real Gitea credentials
   - Verify TODOs removed

2. **E2.2.2-B Implementation**
   - Create effort infrastructure
   - Plan image operations approach
   - Implement Docker/OCI integration
   - Replace all placeholders
   - Remove feature flags
   - Verify production readiness

3. **Integration**
   - Integrate all three efforts together
   - Base: `idpbuilder-oci-build-push/phase2/wave1/integration-20250915-125755`
   - Efforts: cli-commands, credential-management, image-operations
   - Final validation against R320 (NO stubs)

## Production Completeness Requirements

### R320 Compliance (MANDATORY)
- **ZERO** TODO comments in production code
- **ZERO** placeholder strings
- **ZERO** simulated functionality
- **ZERO** stub implementations
- **ZERO** feature flags

### After Integration, Binary MUST:
- Push real Docker images to real Gitea registries
- Use real credentials from environment/config/keyring
- Report actual layer-by-layer progress
- Generate real OCI manifests with valid SHA256 digests
- Handle errors with production-grade messages

## Dependencies

```
Phase 1 Integration
    └── Phase 2 Wave 1 Integration (image-builder, gitea-client-splits)
            └── E2.2.1 (cli-commands) - COMPLETE
                    └── E2.2.2-A (credential-management) - PENDING
                            └── E2.2.2-B (image-operations) - PENDING
```

## Success Criteria

After Phase 2 Wave 2 integration:
1. `idpbuilder build` creates real OCI images from directories
2. `idpbuilder push` pushes real images to Gitea registry
3. No placeholders, stubs, or TODOs remain
4. All tests pass with real functionality
5. Binary is 100% production-ready

## Notes
- E2.2.2 was split due to size (850 lines > 800 limit)
- Sequential implementation required (E2.2.2-B depends on E2.2.2-A)
- All three efforts integrate together as one unit