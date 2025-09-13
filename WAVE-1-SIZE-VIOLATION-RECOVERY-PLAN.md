# Wave 1 Size Violation Recovery Plan

## 🔴 CRITICAL: ALL WAVE 1 EFFORTS EXCEED 800-LINE LIMIT

**Created**: 2025-09-13T01:10:00Z
**Severity**: CRITICAL (per R156: 30-minute recovery target)
**Issue Type**: SIZE_LIMIT_VIOLATION

## Current Violations

| Effort | Current Lines | Violation % | Required Splits |
|--------|--------------|-------------|-----------------|
| kind-cert-extraction | 3170 | 396% | 4-5 splits |
| registry-tls-trust | 2413 | 301% | 3-4 splits |
| registry-auth-types | 1406 | 176% | 2-3 splits |

## Recovery Actions Required

### Phase 1: Split Planning (Parallel)
**Target**: Complete within 15 minutes

1. **Spawn Code Reviewer for kind-cert-extraction split plan**
   - Target: Create 4-5 splits, each <800 lines
   - Focus: Logical separation of cert extraction logic
   
2. **Spawn Code Reviewer for registry-tls-trust split plan**
   - Target: Create 3-4 splits, each <800 lines
   - Focus: Separate TLS setup, validation, and trust management
   
3. **Spawn Code Reviewer for registry-auth-types re-split plan**
   - Target: Re-evaluate existing splits, create proper 2-3 splits
   - Focus: Fix the oversized existing splits

### Phase 2: Implementation (Sequential per effort)
**Target**: Complete within 45 minutes

For each effort:
1. Receive split plan from Code Reviewer
2. Spawn SW Engineer to implement splits sequentially
3. Measure each split with line-counter.sh
4. Verify all splits <800 lines

### Phase 3: Validation
**Target**: Complete within 10 minutes

1. Verify all splits committed and pushed
2. Update orchestrator-state.json with split tracking
3. Prepare for wave re-integration

## R327 Compliance Requirement

Per R327 (Mandatory Re-Integration After Fixes):
- After fixing source branches, MUST delete old integration
- Create fresh integration branch with fixed splits
- This ensures fixes persist and aren't lost

## Success Criteria

✅ All original efforts properly split into <800 line chunks
✅ Each split has its own branch and working directory
✅ All splits pass code review
✅ orchestrator-state.json updated with split_tracking
✅ Ready for fresh wave integration attempt

## Next State Transition

After successful recovery:
- FROM: ERROR_RECOVERY
- TO: SETUP_INTEGRATION_INFRASTRUCTURE (to recreate integration per R327)

## Monitoring

Track recovery progress against R156 time targets:
- Start: 2025-09-13T01:10:00Z
- 50% checkpoint: 2025-09-13T01:25:00Z (escalation if not progressing)
- Target completion: 2025-09-13T01:40:00Z (30-minute CRITICAL target)