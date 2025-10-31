# IMMEDIATE_BACKPORT_REQUIRED - State-Specific Rules

**State**: IMMEDIATE_BACKPORT_REQUIRED
**Agent**: Orchestrator
**Iteration Level**: Generic (Wave/Phase/Project)
**Type**: FIX_ENFORCEMENT
**Checkpoint**: ✅ YES (R322 - user must approve backport plan)

## Purpose

This state enforces R321 (Immediate Backport During Integration Protocol). When integration issues are detected (build failures, test failures, or conflicts), this state IMMEDIATELY stops integration work and forces fixes to be applied to source branches BEFORE retrying integration.

**CRITICAL**: This is NOT a deferred fix state. R321 mandates IMMEDIATE backporting - no deferrals, no "fix it later". Integration branches are READ-ONLY for code changes.

## Entry Conditions

You enter this state when:

1. **Integration Issue Detected**: Any of:
   - Build failure during integration (`npm run build` fails)
   - Test failure during integration (`npm test` fails)
   - Merge conflicts requiring code changes (not just conflict resolution)
   - Integration validation failures
2. **From Integration States**: Transitioning from:
   - `INTEGRATE_WAVE_EFFORTS` (wave-level integration issue)
   - `INTEGRATE_PHASE_WAVES` (phase-level integration issue)
   - `INTEGRATE_PROJECT_PHASES` (project-level integration issue)
3. **Bug Tracking Updated**: Issue recorded in `bug-tracking.json` with R321 enforcement flag

## Critical Rules

### R321: Immediate Backport During Integration Protocol (SUPREME LAW)
**Violation = -50% to -100% AUTOMATIC FAILURE**

**Core Requirements**:
1. **NO CODE CHANGES IN INTEGRATE_WAVE_EFFORTS BRANCHES**: All fixes MUST go to source branches
2. **NO DEFERRALS**: Backport is IMMEDIATE (not "later" or "after review")
3. **SOURCE VALIDATION**: Source branches MUST build and test independently before retry
4. **READ-ONLY INTEGRATE_WAVE_EFFORTS**: Integration branches can only receive merges, never edits

## Required Actions (Summary)

1. Identify integration issue type
2. Identify affected source branches
3. Record backport requirement in bug-tracking.json
4. Create backport plan
5. Present plan to user (R322 checkpoint)
6. Spawn SW Engineers for source fixes
7. Monitor fix completion
8. Validate source branches independently (R321 CRITICAL)
9. Update bug tracking
10. Transition to retry integration

## State File Updates (R288)

Update orchestrator-state-v3.json and bug-tracking.json per R288 atomic update protocol.

## Exit Conditions

Exit when all backports complete, sources validated, and ready to retry integration.

## Allowed Transitions

- START_WAVE_ITERATION, START_PHASE_ITERATION, START_PROJECT_ITERATION
- ERROR_RECOVERY

## Related Rules

- R321 (SUPREME LAW - enforced by this state)
- R361 (Integration Conflict Resolution Only)
- R300 (Effort Branch is Source of Truth)
- R322 (Critical Checkpoints)
- R288 (Atomic State Updates)
