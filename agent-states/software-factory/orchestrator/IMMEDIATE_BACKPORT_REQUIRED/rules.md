# IMMEDIATE_BACKPORT_REQUIRED - State-Specific Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
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
