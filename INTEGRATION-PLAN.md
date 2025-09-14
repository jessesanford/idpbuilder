# Integration Plan - Phase 1 Wave 1 R327 CASCADE
Date: 2025-09-14
CASCADE ID: WAVE1-CASCADE-20250914
Target Branch: idpbuilder-oci-build-push/phase1/wave1-integration

## R327 CASCADE Requirements
- This is a MANDATORY CASCADE re-integration
- Previous integration from 2025-09-12 is STALE
- Fixes have been applied via R321 backports
- Complete re-integration required from fresh base

## Branches to Integrate (ordered per R306)
1. **idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction**
   - Parent: main
   - Contains: Docker API fixes (CRITICAL)
   - Lines: 841 (was split, but branch exists)
   - Status: Has R321 fixes applied

2. **idpbuilder-oci-build-push/phase1/wave1/registry-tls-trust**
   - Parent: main  
   - Lines: 700
   - Status: Ready for integration

3. **idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-001**
   - Parent: main (split from registry-auth-types)
   - Lines: 800
   - Status: First split, foundation types

4. **idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-002**
   - Parent: main (split from registry-auth-types)
   - Lines: 800
   - Status: Second split, additional types

## Merge Strategy
- Order: As listed above (dependency-aware per R306)
- Method: --no-ff to preserve history (R262)
- Conflicts: Document and resolve maintaining functionality
- NO cherry-picking (R262 supreme law)
- NO modification of original branches (R262 supreme law)

## Expected Outcome
- Fully integrated Wave 1 with all fixes included
- Build passing (R291 gate)
- Tests passing (R291 gate)
- Demo scripts functional (R291 gate)
- Ready for Wave 2 incremental merge (R308)

## R291 Gates to Enforce
1. BUILD GATE: go build must succeed
2. TEST GATE: go test must pass
3. DEMO GATE: Demo scripts must execute
4. ARTIFACT GATE: Build outputs must exist
