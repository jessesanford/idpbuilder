# Orchestrator - CASCADE_REINTEGRATION State Rules

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
## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

---

## 🚨🚨🚨 CRITICAL: R510/R405 ENFORCEMENT - TEST 5B ROOT CAUSE FIX 🚨🚨🚨

**ROOT CAUSE**: Test 5B failed because orchestrator displayed status but never executed the 8-step exit checklist or R405 flag. Result: 92-minute timeout, -100% failure.

**BEFORE YOU DO ANYTHING - ACKNOWLEDGE:**
```
✅ CASCADE_REINTEGRATION STATE RULES READ
✅ I UNDERSTAND THE 8-STEP MANDATORY CHECKLIST (below)
✅ I WILL EXECUTE ALL 8 STEPS - NO SKIPPING
✅ I WILL NOT JUST DISPLAY STATUS - I WILL EXECUTE THE CHECKLIST
✅ I WILL OUTPUT CONTINUE-SOFTWARE-FACTORY=TRUE AS FINAL LINE
✅ R510/R405 COMPLIANCE - PENALTY: -100%
```

**CRITICAL**: Status reports are informational ONLY. You must BOTH inform user AND execute all 8 checklist steps. Omitting checklist = -100% failure.

---

## 📋 PRIMARY DIRECTIVES FOR CASCADE_REINTEGRATION STATE

### Core Mandatory Rules:

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - NEVER write, copy, move, or manipulate ANY code files

2. **🔴🔴🔴 R287** - TODO PERSISTENCE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Save TODOs within 30s after write, every 10 messages, before transitions

3. **🔴🔴🔴 R288** - STATE FILE UPDATE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
   - Update orchestrator-state-v3.json before EVERY transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - ALL spawn states require STOP after spawning

### State-Specific Rules:

5. **🔴🔴🔴 R351** - Cascade Execution Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R351-cascade-execution-protocol.md`
   - Handle cascading changes between dependent efforts

6. **🔴🔴🔴 R352** - Overlapping Cascade Protocol (SUPREME)
   - Persistent coordinator for multiple overlapping cascade chains

7. **🔴🔴🔴 R353** - Cascade Focus Protocol (SUPREME)
   - NO diversions during cascade (no size checks, splits, quality assessments)

8. **🔴🔴🔴 R354** - Post-Rebase Review Requirement (SUPREME)
   - Every rebase needs review

---

## 🔴🔴🔴 SUPREME DIRECTIVE: CASCADE IS PERSISTENT COORDINATOR 🔴🔴🔴

CASCADE_REINTEGRATION is a PERSISTENT COORDINATOR that:
1. Maintains control after EVERY operation
2. Supports MULTIPLE overlapping cascade chains
3. Checks for new fixes after each operation
4. Returns here from ALL integration states when cascade_mode=true
5. **ENFORCES R353 CASCADE FOCUS - NO DIVERSIONS!**

**YOU CANNOT LEAVE UNTIL:**
- All cascade chains complete or merged
- No pending fixes without chains
- All integrations fresh
- No new fixes in final check

**EXIT TO:**
- **REVIEW_PHASE_INTEGRATION** if chains converge at phase
- **REVIEW_PROJECT_INTEGRATION** if chains converge at project
- **CASCADE_REINTEGRATION** (loop) if new fixes or active chains
- **INTEGRATE_WAVE_EFFORTS/INTEGRATE_PHASE_WAVES/INTEGRATE_PROJECT_PHASES** for reintegration

---

## 🔴🔴🔴 R353 CASCADE FOCUS - NO DIVERSIONS 🔴🔴🔴

**DURING CASCADE:**
- ❌ NO size checks, splits, quality assessments
- ✅ ONLY validate rebases, check conflicts
- ✅ ALL agents get cascade_mode=true

---

## MANDATORY 8-STEP EXIT CHECKLIST

### ✅ Step 1: Determine Next State
```bash
# Based on cascade status, set PROPOSED_NEXT_STATE
PROPOSED_NEXT_STATE="[determined from cascade state]"
```

### ✅ Step 2: Spawn State Manager (SF 3.0)
```bash
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CASCADE_REINTEGRATION" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "Cascade checkpoint"
```

### ✅ Step 3: Save TODOs (R287)
```bash
save_todos "CASCADE_CHECKPOINT"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - CASCADE checkpoint [R287]"; then
    echo "❌ ERROR: Failed to commit TODO files"
    echo "This is non-fatal but TODOs may be lost in compaction"
    echo "Proceeding with state execution..."
    # Don't exit - TODO commit failure is not fatal
fi

git push || echo "⚠️ WARNING: TODO push failed - committed locally"
echo "✅ TODOs saved and committed"
git push
```

### ✅ Step 4: Output Continuation Flag (R405)
```bash
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

### ✅ Step 5: Stop Processing (R322)
```bash
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT

**Skipping ANY step = FAILURE:**
- Missing Step 1: No next state = stuck forever
- Missing Step 2: State Manager not consulted = R517 violation (-100%)
- Missing Step 3: No TODO save = R287 violation (-20% to -100%)
- Missing Step 4: No flag = R405 violation (-100%)
- Missing Step 5: No exit = R322 violation (-100%)

**ALL STEPS MANDATORY - NO EXCEPTIONS**

---

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - VIOLATION: -100%**

### When to use:
- **TRUE (99.9%)**: R322 checkpoints, state work complete, ready to continue
- **FALSE (0.1%)**: CATASTROPHIC unrecoverable errors only

**CASCADE operations use TRUE** - system knows what to do next from state file.

**See**: `$CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md`

---

## Related Rules

- R327: Mandatory Re-Integration After Fixes
- R348: Cascade State Transitions
- R350: Complete Cascade Dependency Graph
- R351: Cascade Execution Protocol
- R352: Overlapping Cascade Protocol (SUPREME)
- R353: Cascade Focus Protocol (SUPREME)
- R354: Post-Rebase Review Requirement (SUPREME)
- R322: Mandatory Stop After State Transitions
- R405: Automation Continuation Flag (SUPREME)
- R517: Universal State Manager Consultation Law (SUPREME)

---

## CASCADE COORDINATION MANTRA (R352)

```
Multiple fixes may arrive, while cascades are in flight.
CASCADE_REINTEGRATION stays alive, until ALL reach project site.
Chains may merge, chains may split, new fixes join the queue.
The coordinator won't quit, until EVERY fix is through!
```

---

**Violation = -100% AUTOMATIC FAILURE**
