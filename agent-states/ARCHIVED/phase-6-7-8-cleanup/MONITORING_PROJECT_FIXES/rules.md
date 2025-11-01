# Orchestrator - MONITORING_EFFORT_FIXES State Rules

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

## 📋 PRIMARY DIRECTIVES FOR MONITORING_EFFORT_FIXES STATE

### Core Mandatory Rules (ALL orchestrator states must have these):

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Criticality: BLOCKING - Automatic termination, 0% grade
   - Summary: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents

2. **🔴🔴🔴 R287** - TODO PERSISTENCE COMPREHENSIVE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Criticality: SUPREME - -20% to -100% penalty for violations
   - Summary: MUST save TODOs within 30s after write, every 10 messages, before transitions

3. **🔴🔴🔴 R288** - STATE FILE UPDATE REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state-v3.json before EVERY state transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

5. **🔴🔴🔴 R233** - Immediate Action On State Entry
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-all-states-immediate-action.md`
   - Criticality: SUPREME LAW - Must act immediately on entering state
   - Summary: Monitoring states require active checking, not passive waiting

## 🛑🛑🛑 R232 MONITOR STATE REQUIREMENTS 🛑🛑🛑

**CRITICAL**: Before ANY transition from MONITOR_* states, you MUST:
1. Check TodoWrite for pending items
2. Process ALL pending items IMMEDIATELY
3. NO "I will..." statements - only "I am..." with action
4. VIOLATION = AUTOMATIC FAILURE

---

## State Context

**Purpose:**
Monitor SW Engineers as they fix bugs found during project integration.

## Primary Actions

1. **Track Fix Progress**:
   ```bash
   # For each engineer fixing bugs:
   cd /efforts/phase-X-wave-Y-effort-Z
   git status
   git log --oneline -5
   
   # Check for fix commits
   git log --grep="project integration bug"
   ```

2. **Monitor Fix Completion**:
   - Check each effort directory for completed fixes
   - Verify commits match bug fix requirements
   - Track which bugs are resolved

3. **Update Tracking**:
   ```json
   {
     "project_fixes_in_progress": [
       {
         "bug_id": 1,
         "engineer": "sw-engineer-1",
         "branch": "phase-1-wave-2-effort-3",
         "status": "completed",
         "commit": "abc123"
       }
     ]
   }
   ```

4. **Spawn Code Reviewers for Fixed Bugs**:
   - Once a bug fix is complete, spawn Code Reviewer
   - Ensure fix meets requirements
   - Verify no new issues introduced

## Monitoring Protocol

```markdown
## 📊 PROJECT FIX MONITORING_SWE_PROGRESS STATUS

### Fixes In Progress:
- Bug #1: COMPLETED ✅ (commit: abc123)
- Bug #2: IN_PROGRESS 🔄
- Bug #4: COMPLETED ✅ (commit: def456)

### Pending Reviews:
- Bug #1: Review spawned
- Bug #4: Review spawned

### Sequential Fixes Waiting:
- Bug #6: Waiting for bug #1 review to pass
- Bug #7: Waiting for bug #6 completion

### Next Actions:
- Monitor bug #2 completion
- Check review results for bugs #1 and #4
- Spawn sequential fix for bug #6 if ready
```

## Decision Logic

```python
def determine_next_state():
    all_fixes = load_project_fixes()

    # Check if all fixes are complete
    all_complete = all(fix.status == "completed" for fix in all_fixes)

    # Check if all reviews passed
    all_reviewed = all(fix.review_status == "passed" for fix in all_fixes)

    if not all_complete:
        return "MONITORING_EFFORT_FIXES"  # Keep monitoring

    if not all_reviewed:
        return "SPAWN_CODE_REVIEWERS_EFFORT_REVIEW"  # Review fixes

    # 🔴🔴🔴 R327 ENFORCEMENT: CHECK FOR STALE INTEGRATE_WAVE_EFFORTSS 🔴🔴🔴
    # CRITICAL: If fixes were applied to effort branches, ALL integrations
    # containing those efforts are now STALE and MUST be recreated via CASCADE

    stale_integrations_detected = check_for_stale_integrations()

    if stale_integrations_detected:
        # MANDATORY: Cascade re-integration required per R327
        # Cannot go directly to PROJECT_INTEGRATE_WAVE_EFFORTS - must cascade from bottom up
        return "CASCADE_REINTEGRATION"  # Enforce R327 cascade
    else:
        # No stale integrations - can proceed directly to integration
        return "PROJECT_INTEGRATE_WAVE_EFFORTS"  # Re-run FULL integration with fixed code
```

### R327 Stale Integration Detection
```bash
# MANDATORY check before transitioning from MONITORING_EFFORT_FIXES
check_for_stale_integrations() {
    echo "🔍 R327 ENFORCEMENT: Checking for stale integrations after fixes"

    local STALE_DETECTED=false

    # Check each integration level
    for integration_type in wave phase project; do
        local integration_branches=$(git branch -r | grep "${integration_type}.*-integration")

        for integration_branch in $integration_branches; do
            # Get integration creation timestamp
            local INTEGRATE_WAVE_EFFORTS_TIME=$(git log -1 --format=%ct "$integration_branch" 2>/dev/null || echo 0)

            # Get source branches for this integration
            local source_branches=$(get_source_branches_for "$integration_branch")

            # Check if any source has commits newer than integration
            for source in $source_branches; do
                local SOURCE_TIME=$(git log -1 --format=%ct "$source" 2>/dev/null || echo 0)

                if [[ $SOURCE_TIME -gt $INTEGRATE_WAVE_EFFORTS_TIME ]]; then
                    echo "❌ R327 VIOLATION DETECTED!"
                    echo "   Integration: $integration_branch"
                    echo "   Source: $source"
                    echo "   Integration created: $(date -d "@$INTEGRATE_WAVE_EFFORTS_TIME")"
                    echo "   Source last updated: $(date -d "@$SOURCE_TIME")"
                    echo "   🔴 CASCADE RE-INTEGRATE_WAVE_EFFORTS MANDATORY!"
                    STALE_DETECTED=true
                fi
            done
        done
    done

    if [[ "$STALE_DETECTED" == "true" ]]; then
        echo "🔴🔴🔴 STALE INTEGRATE_WAVE_EFFORTSS DETECTED 🔴🔴🔴"
        echo "MUST transition to CASCADE_REINTEGRATION per R327"
        return 0  # Stale detected (shell convention: 0 = true/success)
    else
        echo "✅ All integrations are current"
        return 1  # No stale integrations
    fi
}
```

## Valid State Transitions

- **FIXES_ONGOING** → MONITORING_EFFORT_FIXES (continue monitoring)
- **FIXES_COMPLETE** → SPAWN_CODE_REVIEWERS_EFFORT_REVIEW (review all fixes)
- **🔴 ALL_REVIEWED + STALE_INTEGRATE_WAVE_EFFORTSS** → CASCADE_REINTEGRATION (R327 MANDATORY: Cascade re-integration)
- **🔴 ALL_REVIEWED + NO_STALE_INTEGRATE_WAVE_EFFORTSS** → PROJECT_INTEGRATE_WAVE_EFFORTS (Re-run FULL integration with fixed code)
- **FIXES_FAILED** → ERROR_RECOVERY (unable to fix bugs)

### R327 Enforcement in State Transitions

**CRITICAL DECISION POINT:**
After all fixes are complete and reviewed, you MUST check for stale integrations:

1. **If stale integrations detected** (effort branches have commits newer than integrations):
   - ✅ MUST transition to CASCADE_REINTEGRATION
   - ❌ CANNOT go directly to PROJECT_INTEGRATE_WAVE_EFFORTS
   - Reason: R327 requires cascade deletion and recreation of ALL stale integrations

2. **If NO stale integrations** (unlikely but possible):
   - ✅ Can transition directly to PROJECT_INTEGRATE_WAVE_EFFORTS
   - All integrations already have the fixes somehow

**Default assumption:** If fixes were applied to effort branches, integrations ARE stale

## Critical Requirements

1. **R232 Compliance** - Process all TODOs before considering transition
2. **Track every fix** - Know status of each bug fix
3. **Verify in source branches** - Fixes must be in original branches
4. **Spawn reviews** - All fixes need code review
5. **Re-run integration** - After fixes, must re-integrate

## 🔴🔴🔴 MANDATORY INTEGRATE_WAVE_EFFORTS RE-RUN PROTOCOL 🔴🔴🔴

**CRITICAL: After fixes are complete, you MUST re-run the ENTIRE project integration!**

### Why Re-Integration Is MANDATORY:
- Fixes were applied to UPSTREAM branches (phase/wave/effort branches)
- The project-integration branch still has the BROKEN code
- You MUST re-merge all branches to get the fixed code into integration
- Without re-integration, the binary CANNOT be built

### The CORRECT Re-Integration Cycle:
```
MONITORING_EFFORT_FIXES (all fixes complete & reviewed)
    ↓
PROJECT_INTEGRATE_WAVE_EFFORTS (delete old integration, create fresh)
    ↓
SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN (create new merge plan)
    ↓
SPAWN_INTEGRATION_AGENT_PROJECT (re-merge ALL branches with fixes)
    ↓
MONITORING_PROJECT_INTEGRATE_WAVE_EFFORTS (check if NOW it works)
    ↓
If bugs found → SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING → WAITING_FOR_PROJECT_FIX_PLANS → SPAWN_SW_ENGINEERS → MONITORING_EFFORT_FIXES
If clean → SPAWN_CODE_REVIEWER_DEMO_VALIDATION → PROJECT_DONE
```

### What Happens During Re-Integration:
1. **Delete old broken integration branch** (it has unfixed code)
2. **Create fresh project-integration from main**
3. **Re-run ENTIRE merge plan** (all phases with their fixed branches)
4. **All fixes from upstream branches now merged in**
5. **Binary can finally be built with working code**

### NEVER DO THESE (AUTOMATIC FAILURE):
- ❌ Skip re-integration and claim "fixes are done"
- ❌ Manually copy fixes to integration branch
- ❌ Proceed to validation with broken integration branch
- ❌ Cherry-pick fixes instead of full re-merge

## Grading Impact

- **-50%** for R232 violation (not processing TODOs)
- **-30%** if skipping code reviews for fixes
- **-40%** if not re-running integration after fixes
- **+20%** for comprehensive fix tracking
- **+15%** for efficient review spawning


## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete MONITORING_EFFORT_FIXES:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, set variables for State Manager
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="[REASON_FOR_TRANSITION]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for State Transition (R288 - SUPREME LAW)
```bash
# State Manager handles ALL state file updates atomically
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE_NAME" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"

# State Manager will:
# 1. Validate the transition against state machine
# 2. Update all 4 state tracking locations atomically:
#    - orchestrator-state-v3.json
#    - orchestrator-state-demo.json
#    - .cascade-state-backup.json
#    - .orchestrator-state-v3.json
# 3. Commit and push all changes
# 4. Return control

echo "✅ State Manager completed transition"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before exit (R287 trigger)
save_todos "STATE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo
git commit -m "todo: orchestrator - state complete [R287]"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 5: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 6: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No State Manager spawn = state machine broken (R288 violation, -100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)


## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS (SF 3.0)

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="NEXT_STATE"
TRANSITION_REASON="State work complete"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"
# State Manager updates all 4 state files atomically

# 4. Save TODOs
save_todos "R322_CHECKPOINT"

# 5. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 6. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
### 🚨 MONITORING_SWE_PROGRESS STATE PATTERN - NORMAL TRANSITIONS 🚨

**Monitoring states transition to next actions automatically:**
```bash
# After monitoring completes
echo "✅ Monitoring complete, agents finished work"

# Determine next action from results
if all_succeeded; then
    transition_to "SPAWN_CODE_REVIEWERS_EFFORT_REVIEW"
elif needs_fixes; then
    transition_to "SPAWN_SW_ENGINEERS"
fi

# R322 checkpoint (if required by this transition)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # NORMAL operation!
exit 0  # If R322 checkpoint
```

**Why TRUE is correct:**
- Monitoring results drive automatic actions
- System knows what to do based on results
- **Review findings = Spawn fixes = NORMAL!**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

