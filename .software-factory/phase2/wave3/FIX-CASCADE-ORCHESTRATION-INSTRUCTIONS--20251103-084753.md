# FIX CASCADE ORCHESTRATION INSTRUCTIONS

## Executive Summary

**The Integration Agent incorrectly set CONTINUE-SOFTWARE-FACTORY=FALSE for a normal, recoverable build failure.**

**Correct Value:** `CONTINUE-SOFTWARE-FACTORY=TRUE`

## Situation Analysis

### What Happened
1. Wave 2.3 Iteration 1 integration attempted
2. Build failed with BUG-020 (function redeclarations)
3. Integration Agent documented bug correctly
4. **Integration Agent INCORRECTLY set FALSE**
5. System stopped when it should continue

### Why FALSE Was Wrong

**From R405 Continuation Flag Master Guide:**

#### Recoverable Issue Scenarios → TRUE
```bash
- Tests failing (can fix) → TRUE
- Review found issues (can fix) → TRUE
- Build errors (can debug) → TRUE       # ← THIS CASE
- Integration conflicts (can resolve) → TRUE
- Size violations (can split) → TRUE
```

#### Decision Tree
```
Is something broken?
├─ NO → TRUE
└─ YES → Can system recover automatically?
         ├─ YES → TRUE (use recovery protocols)  # ← THIS CASE
         └─ NO (truly stuck) → FALSE
```

#### Current State Check
- ✅ Iteration 1 of 10 (NOT overflow)
- ✅ Bug identified (BUG-020)
- ✅ Fix location known (effort-2 branch)
- ✅ Fix protocol exists (R300)
- ✅ Recovery path clear (spawn SW Engineer → fix → re-integrate)
- ❌ NO state corruption
- ❌ NO divergence
- ❌ NO unrecoverable error

**Conclusion: This is NORMAL OPERATION. System designed for this.**

### User's Correct Assessment

**User stated:** "WE MUST NOT STOP THE SYSTEM JUST BECAUSE BUGS ARE FOUND! THAT IS NORMAL AND THE SYSTEM MUST CONTINUE!"

**User is 100% correct!** This is the DESIGNED workflow:
1. Integration finds bugs → Document in bug-tracking.json ✅
2. System transitions to ERROR_RECOVERY ✅
3. Create fix plans ✅
4. SW Engineers spawn to fix effort branches ✅
5. System re-integrates in next iteration ✅
6. Repeat until bugs → 0 ✅

**Iteration containers EXPECT this pattern!**

## Fix Cascade Metadata Created

### Fix Plan Location
```
.software-factory/phase2/wave3/effort-2-error-system/fix-cascade/FIX-PLAN-BUG-020--20251103-084753.md
```

### Fix Plan Summary

**Bug:** BUG-020-VALIDATOR-REDECLARATIONS
**Severity:** CRITICAL (P0)
**Effort:** 2.3.2 (error-system)
**Fix Branch:** `idpbuilder-oci-push/phase2/wave3/effort-2-error-system`

**Problem:**
- Stub file `pkg/validator/validator.go` redeclares functions from Effort 2.3.1
- Build fails with redeclaration errors
- Blocks Wave 2.3 integration

**Solution:**
- Remove conflicting stub file from effort-2 branch
- Verify build passes
- Push to remote effort branch
- Re-run integration in iteration 2

**Estimated Time:** 15 minutes
**Risk:** VERY_LOW

## Orchestrator Continuation Instructions

### Current State
```json
{
  "current_state": "ERROR_RECOVERY",
  "phase": 2,
  "wave": 3,
  "iteration": 1,
  "bug_count": 1,
  "bugs": ["BUG-020-VALIDATOR-REDECLARATIONS"]
}
```

### Recommended State Flow

**From ERROR_RECOVERY:**
```
ERROR_RECOVERY
  ↓ (consult State Manager per R517)
CREATE_WAVE_FIX_PLAN (may exist as state)
  ↓
SPAWN_SW_ENGINEERS (for fixing)
  ↓
MONITORING_SWE_PROGRESS
  ↓
START_WAVE_ITERATION (iteration 2)
  ↓
INTEGRATE_WAVE_EFFORTS (retry with fixed efforts)
  ↓
[Convergence: bugs decrease 1→0]
  ↓
WAVE_COMPLETE
```

### Immediate Actions Required

1. **Consult State Manager (R517 MANDATORY)**
   ```bash
   # NEVER update state files directly!
   # ALWAYS spawn State Manager for state transitions

   /spawn state-manager SHUTDOWN_CONSULTATION \
     --current-state "ERROR_RECOVERY" \
     --proposed-next-state "CREATE_WAVE_FIX_PLAN" \
     --work-summary "BUG-020 fix plan created, ready to spawn SW Engineer"
   ```

2. **Spawn SW Engineer for Fix**
   ```bash
   @agent-sw-engineer

   State: FIX_ISSUES
   Working Directory: efforts/phase2/wave3/effort-2-error-system
   Branch: idpbuilder-oci-push/phase2/wave3/effort-2-error-system

   Instructions:
   - Read fix plan: .software-factory/phase2/wave3/effort-2-error-system/fix-cascade/FIX-PLAN-BUG-020--20251103-084753.md
   - Remove conflicting stub file per plan
   - Verify build and tests pass
   - Commit and push to effort branch
   - Create FIX_COMPLETE.marker
   ```

3. **Verify Fix Completion**
   ```bash
   # Check for completion marker
   ls efforts/phase2/wave3/effort-2-error-system/FIX_COMPLETE.marker

   # Verify fix in effort branch (R300 compliance)
   cd efforts/phase2/wave3/effort-2-error-system
   git log --oneline -n 5 | grep -i "BUG-020"

   # Verify pushed to remote
   git fetch origin
   git log origin/idpbuilder-oci-push/phase2/wave3/effort-2-error-system --oneline -n 5
   ```

4. **Start Iteration 2**
   ```bash
   # After fix verified, transition to next iteration
   # State Manager will update iteration counter: 1 → 2

   Transition: ERROR_RECOVERY → START_WAVE_ITERATION
   Reason: "BUG-020 fixed in effort-2 branch, ready for iteration 2"
   ```

5. **Re-Run Integration**
   ```bash
   # START_WAVE_ITERATION will:
   # - Increment iteration counter (1 → 2)
   # - Check convergence (bugs: 1 → expected 0)
   # - Verify within max iterations (2 < 10) ✅
   # - Transition to INTEGRATE_WAVE_EFFORTS

   # INTEGRATE_WAVE_EFFORTS will:
   # - Create FRESH integration branch from main
   # - Merge ALL effort branches (including fixed effort-2)
   # - Run build (should succeed now)
   # - Run tests (should pass)
   # - Complete if clean, or loop if bugs found
   ```

## R300 Compliance Checklist

Before re-integration, MUST verify:
- [ ] Fix exists in effort-2 branch (NOT integration branch)
- [ ] Fix commit has message referencing BUG-020
- [ ] Fix pushed to remote effort branch
- [ ] Build passes in effort workspace
- [ ] Tests pass in effort workspace
- [ ] FIX_COMPLETE.marker exists
- [ ] No direct fixes in integration branch
- [ ] Bug status updated to FIXED

## Continuation Flag Guidance

**When orchestrator completes ERROR_RECOVERY work:**

```bash
# At R322 checkpoint (mandatory stop)
echo "🛑 R322: Checkpoint after ERROR_RECOVERY work complete"

# Continuation flag (independent decision)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # ✅ CORRECT

# Why TRUE?
# - Bug fix protocol successfully executed
# - System knows next state (from state machine)
# - No manual intervention needed
# - Normal iteration convergence workflow
# - Iteration 2 < 10 (within limits)
```

**NEVER use FALSE unless:**
- ❌ Iteration overflow (>10 iterations)
- ❌ State corruption preventing continuation
- ❌ Critical files missing with no recovery path
- ❌ Divergence detected (bugs increasing)

**This scenario is TEXTBOOK TRUE!**

## Expected Outcome

### Iteration 2 Result
```
Bugs Found: 0 (convergence achieved)
Build: SUCCESS
Tests: PASS
Status: WAVE_COMPLETE
Next: Proceed to next wave or phase integration
```

### Metrics
- Iteration Count: 2 (well within 10 limit)
- Convergence: 1 bug → 0 bugs ✅
- Time to Resolution: ~15 minutes
- System Continuity: Maintained ✅

## Key Learnings

1. **Integration build failures are NORMAL**
   - Designed into iteration container pattern
   - Expected in first few iterations
   - System has automatic protocols

2. **CONTINUE-SOFTWARE-FACTORY=TRUE for recoverable issues**
   - System can auto-restart conversation
   - Known recovery protocols exist
   - Within operational limits

3. **R300 compliance is critical**
   - Fixes MUST go to effort branches
   - Integration branches are temporary
   - Effort branches become upstream PRs

4. **State Manager consultation is mandatory**
   - NEVER update state files directly
   - ALWAYS spawn for transitions
   - Atomic 4-file commits

## References

- **R405**: Continuation Flag Master Guide
- **R300**: Comprehensive Fix Management Protocol
- **R517**: Universal State Manager Consultation Law
- **R322**: Mandatory Stop Before State Transitions
- **Bug Tracking**: `bug-tracking.json` (BUG-020)
- **Fix Plan**: `.software-factory/phase2/wave3/effort-2-error-system/fix-cascade/FIX-PLAN-BUG-020--20251103-084753.md`

---

**Created By:** Factory Manager (software-factory-manager agent)
**Timestamp:** 2025-11-03T08:47:53Z
**Purpose:** Guide orchestrator through correct fix cascade continuation
**Compliance:** R383 (metadata in .software-factory with timestamps)
