# 🚨 CRITICAL RESTART STATE 🚨

## CURRENT STATE: SPAWN_CODE_REVIEWER_MERGE_PLAN

## ON RESTART - IMMEDIATE ACTION REQUIRED:
**Spawn code-reviewer agent to create Wave 1 re-integration merge plan**

## WHY THIS IS CRITICAL:
- **Problem Discovered**: All Phase 1 reviews were done on STALE integrations
- **R321 Fixes Timeline**: Completed 2025-09-12 01:21:45 UTC
- **Wave 2 Review Done**: 2025-09-12 00:35:00 UTC (46 minutes BEFORE fixes!)
- **Result**: NO integration branches have R321 fixes properly cascaded

## R327 CASCADE SEQUENCE:
1. ✅ State set to SPAWN_CODE_REVIEWER_MERGE_PLAN
2. ⏳ **ON RESTART**: Spawn code-reviewer for Wave 1 merge plan
3. ⏳ Implement Wave 1 re-integration (with R321 fixes)
4. ⏳ Review Wave 1 integration
5. ⏳ Cascade: Wave 2 re-integration (built on fixed Wave 1)
6. ⏳ Cascade: Phase 1 re-integration (built on fixed Wave 2)
7. ⏳ Get NEW architect reviews on properly integrated branches

## WAVE 1 EFFORTS TO RE-INTEGRATE:
- E1.1.1-kind-cert-extraction (needs shared testutil)
- E1.1.2-registry-tls-trust (BASE - has shared testutil)
- E1.1.3-registry-auth-types-split-001 (fixed contains())
- E1.1.3-registry-auth-types-split-002 (needs shared testutil)

## CRITICAL REQUIREMENTS:
- **Integration Order**: registry-tls-trust FIRST (has shared pkg/testutil)
- **Base Branch**: main (per R308)
- **Workspace**: efforts/phase1/wave1/integration-workspace-r327
- **Update When Done**: stale_integrations.wave1.recreation_completed = true

## DO NOT SKIP THIS:
This re-integration cascade is MANDATORY per R327. The architect reviews were done on integrations that lack critical R321 fixes. We cannot proceed to Phase 2 without proper re-integration and new reviews.

---
State prepared: 2025-09-13T04:35:00Z
Ready for restart and immediate code-reviewer spawn.