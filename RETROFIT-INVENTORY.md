# Demo Retrofit Inventory - Phase 1 & Phase 2 Completed Efforts

## Purpose
This inventory tracks all completed efforts from Phase 1 and Phase 2 that require demo retrofitting per R330 requirements.

## R330 Demo Requirements Overview
Each effort must have:
- 3-5 Demo Objectives (what to demonstrate)
- 2-4 Demo Scenarios with exact commands
- Demo Size Impact calculation (typically 100-150 lines)
- Demo Deliverables list (demo-features.sh, DEMO.md, test-data/)

## Phase 1 Efforts

### Wave 1 (Completed)
- **E1.1.1-kind-certificate-extraction**
  - Status: COMPLETE
  - Location: Integrated into main branch
  - Requires: Demo retrofit for certificate extraction features
  
- **E1.1.2-registry-tls-trust-integration**
  - Status: COMPLETE
  - Location: Integrated into main branch
  - Requires: Demo retrofit for TLS trust integration

### Wave 2 (Completed)
- **E1.2.1-certificate-validation-pipeline**
  - Status: COMPLETE
  - Location: Integrated into main branch
  - Requires: Demo retrofit for validation pipeline features
  
- **E1.2.2-fallback-strategies**
  - Status: COMPLETE
  - Location: Integrated into main branch
  - Requires: Demo retrofit for fallback mechanism demonstrations

## Phase 2 Efforts

### Wave 1 (Completed)
- **E2.1.1-image-builder**
  - Status: COMPLETE
  - Location: /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/image-builder
  - Implementation Plan: EXISTS
  - Requires: Demo retrofit for image building features
  
- **E2.1.2-gitea-client** (Split into 2 parts)
  - Status: COMPLETE (2 splits)
  - Locations:
    - Split 001: /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/gitea-client-split-001
    - Split 002: /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/gitea-client-split-002
  - Implementation Plans: EXIST for both splits
  - Requires: Demo retrofit for Gitea client functionality

## Code Reviewer Tasks

For each effort above, a Code Reviewer will be spawned to:
1. Analyze existing implementation
2. Create Demo Requirements section per R330
3. Update IMPLEMENTATION-PLAN.md with demo section
4. Create DEMO-RETROFIT-PLAN.md with detailed demo scenarios
5. Verify demos are achievable with existing code
6. Calculate total size impact (must stay under 800 lines)

## Expected Deliverables

### Per Effort:
- Updated IMPLEMENTATION-PLAN.md with "## 🎬 Demo Requirements (R330)" section
- DEMO-RETROFIT-PLAN.md with detailed demo scenarios
- Size impact assessment ensuring <800 lines total

### Summary Documents:
- PHASE1-DEMO-RETROFIT-SUMMARY.md
- PHASE2-DEMO-RETROFIT-SUMMARY.md

## Tracking

| Effort ID | Phase | Wave | Retrofit Status | Code Reviewer Spawned | Demo Plan Created |
|-----------|-------|------|-----------------|----------------------|-------------------|
| E1.1.1 | 1 | 1 | PENDING | NO | NO |
| E1.1.2 | 1 | 1 | PENDING | NO | NO |
| E1.2.1 | 1 | 2 | PENDING | NO | NO |
| E1.2.2 | 1 | 2 | PENDING | NO | NO |
| E2.1.1 | 2 | 1 | PENDING | NO | NO |
| E2.1.2-S001 | 2 | 1 | PENDING | NO | NO |
| E2.1.2-S002 | 2 | 1 | PENDING | NO | NO |

## Next Steps
1. Identify actual directories/branches for Phase 1 efforts
2. Spawn Code Reviewers for each effort with retrofit instructions
3. Monitor Code Reviewer progress
4. Aggregate retrofit plans into phase summaries
5. Transition to READY_FOR_DEMO_IMPLEMENTATION state