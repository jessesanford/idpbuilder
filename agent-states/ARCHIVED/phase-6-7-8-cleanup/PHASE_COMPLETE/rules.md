# Orchestrator - COMPLETE_PHASE State Rules

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
## 🟢 NORMAL OPERATION - STOP WITH CONTINUATION FLAG TRUE 🟢

**CRITICAL: You MUST stop after this state, but set flag to TRUE for normal operations**

### MANDATORY STOP PROTOCOL (R322):
1. Complete all phase validation work
2. Update orchestrator-state-v3.json to next state
3. Commit and push the state file
4. Output: CONTINUE-SOFTWARE-FACTORY=TRUE (for normal operations)
5. **STOP PROCESSING** (mandatory for context preservation)

### KEY UNDERSTANDING:
- ✅ Stopping is MANDATORY (prevents context overflow)
- ✅ TRUE flag means "all is well, restart me automatically"
- ✅ External automation sees TRUE and restarts Claude Code
- ✅ The stop preserves context while allowing continuous operation

### R322 CLARIFICATION:
- ✅ COMPLETE_PHASE requires a stop for context preservation
- ✅ Moving to next phase is NORMAL (use TRUE flag)
- ✅ TRUE flag enables automatic continuation after restart
- ✅ Only use FALSE if this is the FINAL phase (Phase 5) of entire project

### NORMAL COMPLETION PROTOCOL:
```bash
# Complete phase validation
echo "✅ Phase $PHASE complete - all waves integrated"

# Determine next state based on project structure
CURRENT_PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
TOTAL_PHASES=$(jq -r '.phases_planned // .total_phases' orchestrator-state-v3.json)
TOTAL_WAVES=$(jq -r '.waves_completed | length' orchestrator-state-v3.json)

# Phase Type Detection
if [ "$CURRENT_PHASE" -lt "$TOTAL_PHASES" ]; then
    # More phases remaining - continue to next phase
    NEXT_STATE="START_PHASE_ITERATION"
    CONTINUE_FLAG="TRUE"
    echo "Multi-phase project: Phase $CURRENT_PHASE of $TOTAL_PHASES complete, continuing to next phase"

elif [ "$TOTAL_PHASES" -gt 1 ]; then
    # Multi-phase project, final phase - MUST do project integration (R283)
    NEXT_STATE="PROJECT_INTEGRATE_WAVE_EFFORTS"
    CONTINUE_FLAG="TRUE"
    echo "CRITICAL: Final phase of multi-phase project - PROJECT INTEGRATE_WAVE_EFFORTS REQUIRED (R283)"
    echo "Project-level integration will verify all $TOTAL_PHASES phases work together"

elif [ "$TOTAL_PHASES" -eq 1 ] && [ "$TOTAL_WAVES" -gt 1 ]; then
    # Single-phase but multi-wave - still need project demo
    NEXT_STATE="PROJECT_INTEGRATE_WAVE_EFFORTS"
    CONTINUE_FLAG="TRUE"
    echo "Single-phase, multi-wave project - project-level demo required"

else
    # Single-phase, single-wave - simple project, can complete
    NEXT_STATE="PROJECT_DONE"
    CONTINUE_FLAG="FALSE"
    echo "Single-phase, single-wave project - proceeding to PROJECT_DONE"
fi

# VALIDATION: Prevent illegal direct PROJECT_DONE transition for multi-phase
if [ "$NEXT_STATE" = "PROJECT_DONE" ] && [ "$TOTAL_PHASES" -gt 1 ]; then
    echo "🔴🔴🔴 CRITICAL ERROR: Attempted illegal COMPLETE_PHASE → PROJECT_DONE transition!"
    echo "Multi-phase projects MUST transition to PROJECT_INTEGRATE_WAVE_EFFORTS (R283)"
    echo "This violation would result in -100% grading penalty"
    exit 283
fi

# Update state

# Commit and continue
git add orchestrator-state-v3.json
git commit -m "state: phase $PHASE complete, transitioning to $NEXT_STATE"
git push

echo "CONTINUE-SOFTWARE-FACTORY=$CONTINUE_FLAG"
exit 0  # STOP processing (mandatory for context preservation)
```

**STOP PROCESSING but set TRUE flag to signal automatic restart by external automation!**

---


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED COMPLETE_PHASE STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

### 🚨🚨🚨 R040 - Documentation Requirements
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R040-documentation-requirements.md`

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="COMPLETE_PHASE complete - [accomplishment description]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for SHUTDOWN_CONSULTATION
```bash
# State Manager validates transition and updates state files (SF 3.0 Pattern)
echo "🔄 Spawning State Manager for SHUTDOWN_CONSULTATION..."

# Prepare work results summary
WORK_RESULTS=$(cat <<EOF
{
  "state_completed": "COMPLETE_PHASE",
  "work_accomplished": [
    "Created phase-level integration branch",
    "Merged all wave integration branches",
    "Generated phase completion metrics",
    "Determined if more phases exist"
  ],
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "COMPLETE_PHASE" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON" \
  --work-results "$WORK_RESULTS"

# State Manager will:
# 1. Validate PROPOSED_NEXT_STATE exists and transition is valid
# 2. Update all 4 state files atomically (R288)
# 3. Commit and push state files
# 4. Return REQUIRED_NEXT_STATE (usually same as proposed unless invalid)

echo "✅ State Manager consultation complete"
echo "✅ State files updated by State Manager"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "COMPLETE_PHASE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo
git commit -m "todo: orchestrator - COMPLETE_PHASE complete [R287]"
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
- Missing Step 2: No proposed next state = State Manager can't proceed
- Missing Step 3: No State Manager consultation = bypassing bookend pattern (-100%)
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
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

