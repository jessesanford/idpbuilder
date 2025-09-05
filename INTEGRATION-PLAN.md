# Integration Testing Plan

## Effort Inventory

### Phase 1 - Certificate Infrastructure
#### Wave 1 - Certificate Management Core
- **E1.1.1**: kind-certificate-extraction (418 lines) - COMPLETED
  - Branch: idpbuilder-oci-go-cr/phase1/wave1/kind-certificate-extraction
- **E1.1.2**: registry-tls-trust-integration (936 lines, 2 splits) - COMPLETED
  - Branch: idpbuilder-oci-go-cr/phase1/wave1/registry-tls-trust-integration

#### Wave 2 - Certificate Validation & Fallback
- **E1.2.1**: certificate-validation-pipeline (431 lines) - COMPLETED
  - Branch: idpbuilder-oci-go-cr/phase1/wave2/certificate-validation-pipeline
- **E1.2.2**: fallback-strategies (658 lines) - COMPLETED
  - Branch: idpbuilder-oci-go-cr/phase1/wave2/fallback-strategies

### Phase 2 - Build & Push Implementation
#### Wave 1 - Core Build Implementation
- **E2.1.1**: go-containerregistry-image-builder (756 lines) - COMPLETED
  - Branch: idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder
- **E2.1.2**: gitea-registry-client (689 lines) - COMPLETED
  - Branch: idpbuilder-oci-go-cr/phase2/wave1/gitea-registry-client

#### Wave 2 - CLI Integration
- **E2.2.1**: cli-commands (800 lines) - COMPLETED
  - Branch: idpbuilder-oci-go-cr/phase2/wave2/cli-commands

## Merge Sequence

The efforts will be merged in dependency order:

1. **Phase 1 Wave 1 Efforts** (Foundation - no dependencies)
   - [ ] E1.1.1: kind-certificate-extraction
   - [ ] E1.1.2: registry-tls-trust-integration

2. **Phase 1 Wave 2 Efforts** (Depends on Wave 1)
   - [ ] E1.2.1: certificate-validation-pipeline
   - [ ] E1.2.2: fallback-strategies

3. **Phase 2 Wave 1 Efforts** (Depends on Phase 1 completion)
   - [ ] E2.1.1: go-containerregistry-image-builder
   - [ ] E2.1.2: gitea-registry-client

4. **Phase 2 Wave 2 Efforts** (Depends on Phase 2 Wave 1)
   - [ ] E2.2.1: cli-commands

## Validation Checkpoints

- [ ] All Phase 1 efforts merged
- [ ] Phase 1 build successful
- [ ] Phase 1 tests passing
- [ ] All Phase 2 efforts merged
- [ ] Phase 2 build successful
- [ ] Phase 2 tests passing
- [ ] Full integration tests passing
- [ ] Feature demos complete

## Conflict Resolution Log

[Conflicts will be documented here as they are encountered and resolved]

## Total Lines Implemented
- Phase 1: 2,443 lines
- Phase 2: 2,245 lines
- **Total**: 4,688 lines