# ORCHESTRATOR STATE: CREATE_WAVE_INTEGRATION_BRANCH_EARLY

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
## STATE PURPOSE
Create wave integration branch IMMEDIATELY after test planning (R342 enforcement). Early integration branch stores wave tests.

## KEY RULES
- R342 (Early Integration Branch): MUST create branch immediately after test planning
- R336 (Wave Integration): Branch from previous wave or phase integration branch
- R288 (State File Update): Update orchestrator-state-v3.json with branch info
- R287 (TODO Persistence): Save before transitioning

## ACTIONS
1. Determine base branch (phase-integration for wave 1, previous wave for wave N)
2. Create wave-integration branch from base
3. Add wave tests to integration branch
4. Commit and push tests (R342: immediate tracking)
5. Update orchestrator-state-v3.json with branch info
6. Transition to SPAWN_CODE_REVIEWER_WAVE_IMPL

## NEXT STATE
SPAWN_CODE_REVIEWER_WAVE_IMPL

## BRANCH NAMING
Format: `phase-${PHASE}-wave-${WAVE}-integration`
Example: `phase-1-wave-2-integration`

## WHY R342 MATTERS
Wave tests must be tracked immediately to ensure:
- Tests available for effort integration
- No test loss between waves
- R341 TDD compliance enforced
- Cumulative test accumulation (R308/R336)

See: rule-library/R342-early-integration-branch-creation.md
See: docs/PROGRESSIVE-TEST-PLANNING-ARCHITECTURE.md

---

## 🔴🔴🔴 R405 AUTOMATION CONTINUATION FLAG 🔴🔴🔴

### CRITICAL FOR THIS STATE

**THIS IS AN R322 CHECKPOINT STATE - USE TRUE!**

After completing the work in this state:
- Wave integration branch created ✅
- Wave tests committed to branch ✅
- State updated successfully ✅
- Ready for next state ✅

**YOU MUST OUTPUT:**
```
CONTINUE-SOFTWARE-FACTORY=TRUE
```

**DO NOT OUTPUT FALSE!**

**WHY TRUE?**
- Your state work completed successfully
- Creating branch and committing tests = normal workflow
- Ready to proceed to SPAWN_CODE_REVIEWER_WAVE_IMPL
- Stopping at checkpoint ≠ failure

**R322 CHECKPOINT = DESIGNED WORKFLOW = TRUE FLAG**

See: rule-library/R405-automation-continuation-flag.md (read the decision flowchart!)

---
