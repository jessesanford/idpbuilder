# Demo Retrofit Inventory

**Date**: 2025-09-10
**Purpose**: Add R330 demo requirements to all completed efforts
**State**: RETROFIT_PLANNING

## Phase 1 Efforts Requiring Demo Retrofit

### Wave 1
- **E1.1.1**: Kind Certificate Extraction
  - Package: `pkg/certs` (extraction functionality)
  - Location: Integrated in `efforts/phase2/wave1/integration-workspace/repo`
  - Status: Implementation complete, needs demo

- **E1.1.2**: Registry TLS Trust Integration  
  - Package: `pkg/certs` (trust functionality)
  - Location: Integrated in `efforts/phase2/wave1/integration-workspace/repo`
  - Status: Implementation complete, needs demo

### Wave 2
- **E1.2.1**: Certificate Validation Pipeline
  - Package: `pkg/certvalidation`
  - Location: Integrated in `efforts/phase2/wave1/integration-workspace/repo`
  - Status: Implementation complete, needs demo

- **E1.2.2**: Fallback Strategies
  - Packages: `pkg/fallback`, `pkg/insecure`
  - Location: Integrated in `efforts/phase2/wave1/integration-workspace/repo`
  - Status: Implementation complete, needs demo

## Phase 2 Efforts Requiring Demo Retrofit

### Wave 1
- **E2.1.1**: Image Builder
  - Package: `pkg/build`
  - Location: `efforts/phase2/wave1/image-builder`
  - Status: Implementation complete, needs demo

- **E2.1.2**: Gitea Client (Split into 2)
  - **Split-001**: Core Infrastructure
    - Package: `pkg/registry` (core)
    - Location: `efforts/phase2/wave1/gitea-client-split-001`
    - Status: Implementation complete, needs demo
  
  - **Split-002**: Push/List Operations
    - Package: `pkg/registry` (operations)
    - Location: `efforts/phase2/wave1/gitea-client-split-002`
    - Status: Implementation complete, needs demo

## Summary

- **Total Efforts**: 6 (4 in Phase 1, 2 in Phase 2)
- **Total Splits**: 7 (including gitea-client splits)
- **Packages Implemented**: 
  - `pkg/certs` (Phase 1)
  - `pkg/certvalidation` (Phase 1)
  - `pkg/fallback` (Phase 1)
  - `pkg/insecure` (Phase 1)
  - `pkg/build` (Phase 2)
  - `pkg/registry` (Phase 2)

## Next Steps

1. Spawn Code Reviewers for each effort/split
2. Create DEMO-RETROFIT-PLAN.md for each
3. Update IMPLEMENTATION-PLAN.md with demo requirements
4. Calculate size impact (must stay under 800 lines)
5. Define demo scenarios with exact commands
6. Ensure demos integrate at wave and phase levels