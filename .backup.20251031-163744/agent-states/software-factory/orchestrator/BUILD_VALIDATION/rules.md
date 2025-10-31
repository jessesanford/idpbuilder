# Orchestrator - BUILD_VALIDATION State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state-v3.json with new state
3. ✅ Committing and pushing the state file  
4. ✅ Providing work summary

### YOU MUST NOT:
- ❌ Continue to the next state automatically
- ❌ Start work for the new state
- ❌ Spawn agents for the new state
- ❌ Assume permission to continue

### STOP PROTOCOL:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: CURRENT_STATE → NEXT_STATE

### ✅ Current State Work Completed:
- [List completed work]

### 📊 Current Status:
- Current State: CURRENT_STATE
- Next State: NEXT_STATE
- TODOs Completed: X/Y
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to NEXT_STATE. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED BUILD_VALIDATION STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_BUILD_VALIDATION
echo "$(date +%s) - Rules read and acknowledged for BUILD_VALIDATION" > .state_rules_read_orchestrator_BUILD_VALIDATION
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY BUILD WORK UNTIL RULES ARE READ:
- ❌ Start compilation
- ❌ Run build commands
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES

### 🚨🚨🚨 RULE R006 - Orchestrator NEVER Writes Code [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality:** BLOCKING - Any code operation = -100% IMMEDIATE FAILURE

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

**⚠️ R006 WARNING FOR BUILD_VALIDATION STATE:**
- DO NOT fix compilation errors yourself!
- DO NOT edit code to make it compile!
- DO NOT modify any source files!
- If build fails, document issues for SW Engineers
- You only run builds and report - NEVER fix code

### 🚨🚨🚨 RULE R319 - Orchestrator NEVER Measures or Assesses Code [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-code.md`
**Criticality:** BLOCKING - Any technical assessment = -100% IMMEDIATE FAILURE

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-code.md`

**⚠️ R319 WARNING FOR BUILD_VALIDATION STATE:**
- DO NOT run builds yourself!
- DO NOT execute test commands!
- DO NOT validate artifacts!
- DO NOT assess technical compliance!
- You MUST spawn Code Reviewer to handle ALL build validation
- You only coordinate and read reports - NEVER validate directly

### 🚨🚨🚨 RULE R323 - Mandatory Final Artifact Build [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R323-mandatory-final-artifact-build.md`
**Criticality:** BLOCKING - No artifact = -50% to -100% FAILURE

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R323-mandatory-final-artifact-build.md`

**⚠️ R323 ENFORCEMENT FOR BUILD_VALIDATION STATE:**
- Code Reviewer MUST build final deliverable artifact
- Code Reviewer MUST verify artifact exists
- Code Reviewer MUST document artifact path, size, type
- CANNOT transition to PROJECT_DONE without artifact confirmation
- NO project can complete without built binary/package

### ✅ Step 2: Set Proposed Next State and Transition Reason

**⚠️ CRITICAL: BUILD_VALIDATION at Wave Level MUST Transition to Phase Integration! ⚠️**

The state machine REQUIRES the following integration hierarchy:
1. **Efforts** integrate into **Waves** (INTEGRATE_WAVE_EFFORTS)
2. **Waves** integrate into **Phases** (INTEGRATE_PHASE_WAVES)
3. **Phases** integrate into **Project** (INTEGRATE_PROJECT_PHASES)
4. Only at **Project** level can PR_PLAN_CREATION occur

```bash
# CRITICAL: Determine integration level and next state
CURRENT_PHASE=$(jq -r '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
CURRENT_WAVE=$(jq -r '.project_progression.current_wave.wave_number' orchestrator-state-v3.json)

# BUILD_VALIDATION at wave level MUST go to phase infrastructure setup
if [[ "$CURRENT_WAVE" != "null" ]]; then
    echo "⚠️ Currently at Wave $CURRENT_WAVE of Phase $CURRENT_PHASE"
    echo "✅ Wave build validation complete - must proceed to PHASE integration"
    PROPOSED_NEXT_STATE="SETUP_PHASE_INFRASTRUCTURE"
    TRANSITION_REASON="Wave $CURRENT_WAVE build validation complete - proceeding to phase integration per integration hierarchy"
else
    echo "❌ ERROR: Cannot determine current integration level"
    echo "Check orchestrator-state-v3.json for current_phase and current_wave"
    exit 1
fi

echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

**❌ NEVER transition directly to PR_PLAN_CREATION from wave-level BUILD_VALIDATION!**
**✅ ALWAYS follow the integration hierarchy: Wave → Phase → Project**

---

### ✅ Step 3: Spawn State Manager for SHUTDOWN_CONSULTATION
```bash
# State Manager validates transition and updates state files (SF 3.0 Pattern)
echo "🔄 Spawning State Manager for SHUTDOWN_CONSULTATION..."

# Prepare work results summary
WORK_RESULTS=$(cat <<EOF
{
  "state_completed": "BUILD_VALIDATION",
  "work_accomplished": [
    "Spawned Code Reviewer for build validation",
    "Verified build validation report received",
    "Validated artifact generation per R323"
  ],
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "BUILD_VALIDATION" \
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
save_todos "BUILD_VALIDATION_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - BUILD_VALIDATION complete [R287]"; then
    echo "❌ ERROR: Failed to commit TODO files"
    echo "This is non-fatal but TODOs may be lost in compaction"
    echo "Proceeding with state execution..."
    # Don't exit - TODO commit failure is not fatal
fi

git push || echo "⚠️ WARNING: TODO push failed - committed locally"
echo "✅ TODOs saved and committed"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 5: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
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

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

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
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

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

