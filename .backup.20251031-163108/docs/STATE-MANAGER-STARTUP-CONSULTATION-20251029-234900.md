# STATE MANAGER STARTUP CONSULTATION
**Consultation Type:** STARTUP
**Timestamp:** 2025-10-29T23:49:00Z
**Orchestrator State:** CREATE_WAVE_FIX_PLAN
**Previous State:** WAITING_FOR_FIX_PLANS

---

## VALIDATION STATUS: ✅ PASSED

### State File Integrity
- **File:** orchestrator-state-v3.json
- **Current State:** CREATE_WAVE_FIX_PLAN (VALID)
- **Previous State:** WAITING_FOR_FIX_PLANS (VALID)
- **Phase/Wave:** Phase 1, Wave 2
- **File Size:** Large but readable
- **JSON Validity:** Valid structure confirmed

### Critical State Checks
✅ State machine current_state is valid orchestrator state
✅ Previous state is valid and matches expected transition
✅ Phase 1, Wave 2 context confirmed
✅ Fix plans directory exists and contains required files
✅ Fix plan summary file present with complete data

---

## PRIMARY DIRECTIVE

**State:** CREATE_WAVE_FIX_PLAN

**Core Objective:** Analyze bugs found during wave integration review, categorize them, identify affected upstream branches, create fix plan, and obtain user approval before proceeding to fixes.

**Critical Context:**
- Code Reviewer has already created fix plans in efforts/phase1/wave2/fix-plans/
- Fix plan summary exists: FIX-PLAN-SUMMARY--20251029-233955.yaml
- 2 efforts need fixes:
  - effort-2-registry-client (CRITICAL: R320 stub violations)
  - effort-3-auth (MINOR: R383 metadata placement)
- Both fixes can be executed in parallel
- Fix plans are ALREADY CREATED by Code Reviewer

**IMPORTANT CLARIFICATION:**
The orchestrator does NOT need to create fix plans - they already exist from the Code Reviewer agent in the SPAWN_CODE_REVIEWER_FIX_PLAN → WAITING_FOR_FIX_PLANS cycle. The CREATE_WAVE_FIX_PLAN state is for the ORCHESTRATOR to:
1. Read and understand the existing fix plans
2. Prepare the execution strategy
3. Present to user for approval (R322)
4. Prepare to spawn SW Engineers with the fix plans

---

## REQUIRED ACTIONS

The orchestrator MUST complete these actions in CREATE_WAVE_FIX_PLAN state:

### 1. READ EXISTING FIX PLANS (BLOCKING)
**Action:** Read all fix plan files created by Code Reviewer
**Files:**
- efforts/phase1/wave2/fix-plans/FIX-PLAN-SUMMARY--20251029-233955.yaml
- efforts/phase1/wave2/fix-plans/FIX-PLAN-effort-2-registry-client--20251029-233955.md
- efforts/phase1/wave2/fix-plans/FIX-PLAN-effort-3-auth--20251029-233955.md

**Purpose:** Understand what fixes are required and execution strategy

### 2. ANALYZE FIX PLAN SUMMARY (BLOCKING)
**Action:** Extract key information from fix plan summary:
- Total efforts needing fixes: 2
- Severity breakdown: 1 critical, 1 minor
- Can fix in parallel: YES
- Estimated total time: 60 minutes
- Integration blockers: 2 (both must fix before integration)

### 3. PREPARE EXECUTION STRATEGY (BLOCKING)
**Action:** Plan how to spawn SW Engineers for fixes
**Strategy from fix plan summary:**
- Spawn SW Engineer for effort-2-registry-client (CRITICAL, HIGH priority)
- Spawn SW Engineer for effort-3-auth (MINOR, MEDIUM priority)
- Both can spawn in parallel (R151 compliant)
- Each gets their specific fix plan file

### 4. PRESENT FIX PLAN TO USER (R322 CHECKPOINT - BLOCKING)
**Action:** Display comprehensive summary for user approval
**Include:**
- Number of efforts needing fixes (2)
- Severity levels (1 critical, 1 minor)
- Estimated time (60 minutes total)
- Parallel execution capability (YES)
- Integration blockers identified
- Next steps after approval

**Format:** Clear, actionable summary
**Requirement:** MUST await user approval before proceeding

### 5. UPDATE STATE FILE (R288 - BLOCKING)
**Action:** Update orchestrator-state-v3.json with:
- Fix plan analysis timestamp
- Number of efforts to fix (2)
- Affected efforts list
- Parallel execution plan
- User approval status

### 6. EXECUTE STATE TRANSITION (R288 - BLOCKING)
**Action:** Spawn State Manager for SHUTDOWN_CONSULTATION
**Next State:** SPAWN_ENGINEERS_FOR_FIXES (SF 3.0 state)
**Reason:** Fix plans analyzed and approved, ready to spawn engineers
**Validation:** State Manager validates and atomically updates all 4 state files

### 7. SAVE TODOS (R287 - BLOCKING)
**Action:** Save TODO state before transition
**Trigger:** "CREATE_WAVE_FIX_PLAN_COMPLETE"
**Format:** todos/orchestrator-CREATE_WAVE_FIX_PLAN-20251029-HHMMSS.todo
**Commit:** Within 60 seconds with proper message

### 8. SET CONTINUATION FLAG (R405 - SUPREME LAW)
**Action:** Output as ABSOLUTE LAST LINE before exit
**Value:**
- TRUE if user approves fix plan (normal case)
- FALSE if user rejects or defers decision
**Context:** This IS an R322 checkpoint - flag depends on user approval

### 9. STOP FOR CHECKPOINT (R322 - SUPREME LAW)
**Action:** Exit cleanly with exit 0
**Message:** "R322 CHECKPOINT: Wave fix plan approval required"
**Next:** /continue-software-factory will proceed to SPAWN_ENGINEERS_FOR_FIXES if approved

---

## STATE-SPECIFIC RULES TO ACKNOWLEDGE

The orchestrator MUST read and acknowledge these rules:

**Core Mandatory (All States):**
1. R006 - Orchestrator never writes code (BLOCKING)
2. R287 - TODO persistence comprehensive (SUPREME)
3. R288 - State file update requirements (SUPREME)
4. R510 - State execution checklist compliance (SUPREME)
5. R405 - Automation continuation flag (SUPREME)

**State-Specific:**
6. R322 - Mandatory checkpoints (SUPREME - user approval required)
7. R313 - Bug tracking requirements (SUPREME - if using bug-tracking.json)
8. R321 - Immediate backport protocol (SUPREME - fixes go to upstream)

---

## CRITICAL WARNINGS

### ⚠️ WARNING 1: Fix Plans Already Created
The orchestrator might think it needs to CREATE fix plans. **THIS IS WRONG.**
- Fix plans are ALREADY created by Code Reviewer
- They exist in efforts/phase1/wave2/fix-plans/
- Orchestrator's job: READ, ANALYZE, PRESENT, GET APPROVAL

### ⚠️ WARNING 2: R322 Checkpoint Required
This state CANNOT proceed without user approval:
- Must present clear summary to user
- Must await explicit approval
- CONTINUE flag depends on user response
- Do NOT auto-approve or skip checkpoint

### ⚠️ WARNING 3: Parallel Execution Allowed
Fix plan summary explicitly states can_fix_in_parallel: true
- Both SW Engineers can spawn in parallel (R151 compliant)
- Timestamps must be within 5 seconds
- This is OPTIMAL for grading (15% parallelization score)

### ⚠️ WARNING 4: State Transition to SF 3.0 State
Next state is SPAWN_ENGINEERS_FOR_FIXES (SF 3.0 naming):
- NOT "SPAWN_SW_ENGINEERS_FOR_FIXES"
- NOT "FIX_WAVE_UPSTREAM_BUGS"
- Use exact state name from state machine

---

## VALIDATION CHECKS PERFORMED

✅ orchestrator-state-v3.json exists and is readable
✅ Current state is CREATE_WAVE_FIX_PLAN (valid)
✅ Previous state is WAITING_FOR_FIX_PLANS (valid transition)
✅ Phase 1, Wave 2 context confirmed
✅ Fix plans directory exists
✅ Fix plan summary file exists and contains complete data
✅ Individual fix plan files exist (2 efforts)
✅ All required metadata present in fix plan summary

---

## SUCCESS CRITERIA FOR THIS STATE

The orchestrator will have successfully completed CREATE_WAVE_FIX_PLAN when:

1. ✅ All fix plan files have been read and understood
2. ✅ Execution strategy prepared (parallel SW Engineer spawning)
3. ✅ User presented with clear fix plan summary
4. ✅ User approval obtained (R322)
5. ✅ State file updated with fix plan details
6. ✅ State Manager SHUTDOWN_CONSULTATION completed
7. ✅ TODOs saved and committed (R287)
8. ✅ CONTINUE-SOFTWARE-FACTORY flag output (R405)
9. ✅ Clean exit with R322 checkpoint message

---

## NEXT STATE AFTER COMPLETION

**Next State:** SPAWN_ENGINEERS_FOR_FIXES
**Trigger:** User approves fix plan at R322 checkpoint
**Action:** Orchestrator will spawn 2 SW Engineers in parallel to implement fixes
**Timeline:** After approval, ~60 minutes for fixes + re-review

---

## GRADING CONSIDERATIONS

**This state tests:**
- R322 checkpoint compliance (25% workflow)
- R151 parallel spawning setup (15% parallelization)
- R287 TODO persistence (20% quality)
- R288 state file updates (20% quality)
- R405 continuation flag (automation)

**Optimal path:**
1. Read fix plans quickly
2. Present clear summary to user
3. Get approval
4. Set up parallel spawn
5. Execute clean transition

---

**State Manager Agent**
**Consultation Complete: 2025-10-29T23:49:00Z**
