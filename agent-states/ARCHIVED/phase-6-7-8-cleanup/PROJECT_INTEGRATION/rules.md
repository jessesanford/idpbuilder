# Orchestrator - PROJECT_INTEGRATE_WAVE_EFFORTS State Rules

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

**YOU HAVE ENTERED PROJECT_INTEGRATE_WAVE_EFFORTS STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_PROJECT_INTEGRATE_WAVE_EFFORTS
echo "$(date +%s) - Rules read and acknowledged for PROJECT_INTEGRATE_WAVE_EFFORTS" > .state_rules_read_orchestrator_PROJECT_INTEGRATE_WAVE_EFFORTS
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## 📋 PRIMARY DIRECTIVES FOR PROJECT_INTEGRATE_WAVE_EFFORTS STATE

### 🚨🚨🚨 R283 - Project Integration Protocol [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R283-project-integration-protocol.md`
**Criticality**: BLOCKING - Mandatory for project completion
**Summary**: Final project integration MUST occur in isolated workspace with all phases merged

### 🚨🚨🚨 R006 - Orchestrator NEVER Writes Code [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality**: BLOCKING - Any code operation = -100% IMMEDIATE FAILURE
**Summary**: Orchestrator is a COORDINATOR ONLY - never writes, edits, or modifies code

### 🚨🚨🚨 R269 - Merge Plan Requirement [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R269-code-reviewer-merge-plan-no-execution.md`
**Criticality**: BLOCKING - Must have formal merge plan
**Summary**: Code Reviewer must create detailed merge plan before integration

### 🚨🚨🚨 R270 - Merge Order Protocol [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R270-no-integration-branches-as-sources.md`
**Criticality**: BLOCKING - Merge order is critical
**Summary**: Phases must be merged in dependency order (Phase 1 → Phase 2 → ...)

### 🚨🚨🚨 R009 - Mandatory Wave/Phase Integration Protocol [SUPREME LAW]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R009-integration-branch-creation.md`
**Criticality**: BLOCKING - Must use target repository
**Summary**: Integration happens in target repository, NOT software-factory

### 🚨🚨🚨 R321 - Immediate Backport During Integration Protocol [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-during-integration.md`
**Criticality**: BLOCKING - Fixes found during integration must be backported immediately
**Summary**: When integration issues are found, must stop and backport fixes to original branches

### 🚨🚨🚨 R280 - Main Branch Protection Protocol [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R280-main-branch-protection.md`
**Criticality**: BLOCKING - Direct commits to main/master are forbidden
**Summary**: All changes must go through PR process with proper reviews

### 🚨🚨🚨 R307 - Branch Mergeability Check [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R307-independent-branch-mergeability.md`
**Criticality**: BLOCKING - Must verify branches are mergeable before attempting
**Summary**: Check for conflicts and mergeability before integration operations

### 🚨🚨🚨 R328 - Integration Freshness Validation [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R328-integration-freshness-validation.md`
**Criticality**: BLOCKING - Stale integrations cause failed merges and lost fixes
**Summary**: MUST verify all phase branches are fresh before creating project integration

## 🚨 PROJECT_INTEGRATE_WAVE_EFFORTS IS A VERB - COORDINATE PROJECT INTEGRATE_WAVE_EFFORTS NOW! 🚨

### IMMEDIATE ACTIONS UPON ENTERING PROJECT_INTEGRATE_WAVE_EFFORTS

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Check if project integration infrastructure exists NOW
2. If NO infrastructure: Transition to SETUP_PROJECT_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE
3. If infrastructure EXISTS: Transition to SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN
4. Update state file with the appropriate next state
5. Stop per R322 for state transition

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in PROJECT_INTEGRATE_WAVE_EFFORTS" [stops]
- ❌ "Successfully entered PROJECT_INTEGRATE_WAVE_EFFORTS state" [waits]
- ❌ "Ready to set up project integration" [pauses]
- ❌ "I'm in PROJECT_INTEGRATE_WAVE_EFFORTS state" [does nothing]
- ❌ Creating infrastructure yourself (PROJECT_INTEGRATE_WAVE_EFFORTS only coordinates!)
- ❌ Merging branches yourself (R329 violation!)

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "PROJECT_INTEGRATE_WAVE_EFFORTS STATE: Checking for existing project integration infrastructure..."
- ✅ "No infrastructure found, transitioning to SETUP_PROJECT_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE..."
- ✅ "Infrastructure exists, transitioning to SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN..."

## State Context

**Prerequisites:**
- COMPLETE_PHASE has been reached (last phase completed)
- All phases have been individually integrated
- Architect has approved the final phase
- Ready for project-wide integration

**Purpose:**
**THIS STATE IS FOR COORDINATION ONLY!**

The PROJECT_INTEGRATE_WAVE_EFFORTS state is a decision point that:
1. **CHECKS** if project integration infrastructure exists
2. **TRANSITIONS** to SETUP_PROJECT_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE if no infrastructure
3. **TRANSITIONS** to SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN if infrastructure exists

**THIS STATE NEVER:**
- ❌ Creates project integration workspace itself
- ❌ Sets up branches or directories itself
- ❌ Performs any actual integration work
- ❌ Executes git merges (R329 violation!)

You are the COORDINATOR of project integration flow.

## 🔴🔴🔴 CRITICAL: VERIFY PHASE BRANCH FRESHNESS FIRST! 🔴🔴🔴

**Before creating project integration, MUST verify all phase branches are fresh:**

```bash
verify_phase_branch_freshness() {
    echo "🔍 Checking phase branch freshness..."
    
    # Check if any effort branches have been updated after phase integrations
    STALE_PHASES=""
    
    for phase_num in $(seq 1 $(echo "✅ State file updated to: $NEXT_STATE"
```

---

### ✅ Step 4: Validate State File (R324)
```bash
# Validate state file before committing
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
echo "✅ State file validated"
```

---

### ✅ Step 5: Commit State File (R288)
```bash
# Commit and push state file immediately
git add orchestrator-state-v3.json
git commit -m "state: PROJECT_INTEGRATE_WAVE_EFFORTS → $NEXT_STATE - PROJECT_INTEGRATE_WAVE_EFFORTS complete [R288]"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "PROJECT_INTEGRATE_WAVE_EFFORTS_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo
git commit -m "todo: orchestrator - PROJECT_INTEGRATE_WAVE_EFFORTS complete [R287]"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 7: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 8: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No state update = state machine broken (R288 violation, -100%)
- Missing Step 4: Invalid state = corruption (R324 violation)
- Missing Step 5: No commit = state lost on compaction (R288 violation, -100%)
- Missing Step 6: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 7: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 8: No exit = R322 violation (-100%)

**ALL 8 STEPS ARE MANDATORY - NO EXCEPTIONS**

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

