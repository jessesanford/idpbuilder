# WAITING_FOR_PHASE_MERGE_PLAN State Rules


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_PHASE_MERGE_PLAN STATE

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

4. **🔴🔴🔴 R322** - Mandatory Stop Before State Transitions
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: Stop at specific checkpoints; stop inference ≠ FALSE flag

5. **🔴🔴🔴 R405** - Mandatory Automation Continuation Flag
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md`
   - Criticality: SUPREME LAW - Output flag as last line
   - Summary: Output CONTINUE-SOFTWARE-FACTORY=TRUE/FALSE as last line; default TRUE

### State-Specific Rules:

6. **🔴🔴🔴 R233** - Immediate Action On State Entry
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-all-states-immediate-action.md`
   - Criticality: SUPREME LAW - Must act immediately on entering state
   - Summary: WAITING states require active monitoring, not passive waiting

## State Purpose
Actively monitor Code Reviewer creating phase merge plan for integrating all wave branches into the phase integration branch. Read merge plan location from orchestrator-state-v3.json per R340. This is critical for phase-level integration.

## Critical Rules

### 🛑🛑🛑 RULE R322: MANDATORY CHECKPOINT BEFORE SPAWN_INTEGRATION_AGENT_PHASE (SUPREME LAW) 🛑🛑🛑

**THIS IS A CRITICAL R322 CHECKPOINT STATE!**

When transitioning from WAITING_FOR_PHASE_MERGE_PLAN → SPAWN_INTEGRATION_AGENT_PHASE:
- **MUST STOP** to allow user review of PHASE-MERGE-PLAN.md
- **MUST UPDATE** state file to SPAWN_INTEGRATION_AGENT_PHASE before stopping
- **MUST DISPLAY** checkpoint message with plan location
- **MUST EXIT** cleanly to preserve context
- **VIOLATION = -100% IMMEDIATE FAILURE**

See: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`

### 🔴🔴🔴 RULE R233: IMMEDIATE ACTION REQUIRED (SUPREME LAW)
- **NO PASSIVE WAITING** - Must actively check for completion
- **IMMEDIATE ACTION** - Start checking within first response
- **CONTINUOUS MONITORING_SWE_PROGRESS** - Check every 30-60 seconds
- **States are VERBS** - "WAITING" means "ACTIVELY CHECKING"

### 🔴🔴🔴 RULE R285: MANDATORY PHASE INTEGRATE_WAVE_EFFORTS (SUPREME LAW)
- Phase integration MUST happen before phase assessment
- All waves MUST be integrated into phase branch
- Integration MUST follow sequential wave order
- Failed integration triggers IMMEDIATE_BACKPORT per R321

### 🚨🚨🚨 RULE R290: STATE RULE VERIFICATION (BLOCKING)
- **MUST** verify this rules file exists and is loaded
- **MUST** acknowledge all rules before proceeding
- **MUST** validate state transitions against state machine

### 🚨🚨🚨 RULE R232: MONITOR STATE REQUIREMENTS (BLOCKING)
- **MUST** check TodoWrite for pending items BEFORE transition
- **MUST** process ALL pending items immediately
- **NO** "I will..." statements - only "I am..." with action
- **VIOLATION = AUTOMATIC FAILURE**

### 🚨🚨🚨 RULE R340: PLANNING FILE METADATA TRACKING (BLOCKING)
- **MUST** read merge plan location from orchestrator-state-v3.json
- **NEVER** search directories for planning files
- **ALWAYS** use planning_files.merge_plans.phase section
- **VIOLATION = -20% for each untracked file**

### ⚠️⚠️⚠️ RULE R269: PHASE MERGE PLAN REQUIREMENTS (WARNING)
- Plan MUST be created as PHASE-MERGE-PLAN.md
- Plan MUST list all wave branches in sequential order
- Plan MUST specify integration strategy
- Plan MUST identify cross-wave dependencies

### ⚠️⚠️⚠️ RULE R287: TODO PERSISTENCE (WARNING)
- **MUST** save TODOs every 10 messages or 15 minutes
- **MUST** save before state transition
- **MUST** commit and push TODO state

## Required Actions

1. **Initial Check (IMMEDIATE)**
   ```bash
   # Verify Code Reviewer was spawned
   grep "SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN" orchestrator-state-v3.json
   grep "spawned_agents" orchestrator-state-v3.json | tail -5
   
   # Check phase integration directory exists
   ls -la phase-*-integration/
   
   # Verify all waves completed
   grep "waves_completed" orchestrator-state-v3.json
   ```

2. **Active Monitoring Loop (R340 Compliant)**
   ```bash
   # Per R340: Read merge plan location from state file
   PHASE=$(echo "✅ State file updated to: $NEXT_STATE"
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

if ! git commit -m "state: WAITING_FOR_PHASE_MERGE_PLAN → $NEXT_STATE - WAITING_FOR_PHASE_MERGE_PLAN complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: WAITING_FOR_PHASE_MERGE_PLAN"
    echo "Attempted transition from: WAITING_FOR_PHASE_MERGE_PLAN"
    echo ""
    echo "Common causes:"
    echo "  - Schema validation failure (check pre-commit hook output above)"
    echo "  - Missing required fields in JSON files"
    echo "  - Invalid JSON syntax"
    echo ""
    echo "🛑 Cannot proceed - manual intervention required"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=SCHEMA_VALIDATION"
    exit 1
fi

git push || echo "⚠️ WARNING: Push failed - committed locally"
echo "✅ State file committed and pushed"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "WAITING_FOR_PHASE_MERGE_PLAN_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - WAITING_FOR_PHASE_MERGE_PLAN complete [R287]"; then
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

### ✅ Step 7: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
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

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**


## Common Violations to Avoid

1. **Passive waiting** - Violates R233, must actively check
2. **Not validating wave count** - Plan missing waves
3. **Ignoring sequential order** - Waves out of order (R285)
4. **Missing dependency check** - Cross-wave conflicts
5. **Not checking remote branches** - Integration will fail

## Phase vs Wave Merge Plans

Key differences:
- **Scope**: Phase plans merge waves, wave plans merge efforts
- **Complexity**: Phase has cross-wave dependencies
- **Scale**: Phase typically 3-5 waves vs 5-10 efforts
- **Risk**: Phase failures affect entire project phase

## Monitoring Pattern

```bash
# CORRECT (R340 compliant): Check state file for tracked plan
echo "Starting phase merge plan monitoring at $(date)"
PHASE=$(jq -r '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
PHASE_ID="phase${PHASE}"
CHECKS=0
MAX_CHECKS=60  # 30 minutes

while [ $CHECKS -lt $MAX_CHECKS ]; do
  CHECKS=$((CHECKS + 1))
  echo "Check #$CHECKS at $(date)"
  
  # R340: Check if plan is tracked in state
  PLAN_PATH=$(jq -r ".planning_files.merge_plans.phase[\"${PHASE_ID}\"].file_path" orchestrator-state-v3.json)
  
  if [ "$PLAN_PATH" != "null" ] && [ -f "$PLAN_PATH" ]; then
    echo "✓ Plan tracked and exists: $PLAN_PATH"
    # Validate immediately
    grep -c "wave-" "$PLAN_PATH"
    break
  fi
  
  # Check Code Reviewer state
  REVIEWER_STATE=$(jq -r '.spawned_agents[] | select(.name == "code-reviewer") | .state' orchestrator-state-v3.json)
  echo "Code Reviewer state: $REVIEWER_STATE"
  
  sleep 30
done

if [ $CHECKS -eq $MAX_CHECKS ]; then
  echo "✗ Timeout reached"
fi

# WRONG (R340 violation): Searching directories
echo "Looking for plan..."
if [ -f phase-*/PHASE-MERGE-PLAN.md ]; then  # ❌ R340 VIOLATION!
  echo "Found plan"
fi
```

## Verification Commands

```bash
# Verify state entry
echo "===================="
echo "WAITING_FOR_PHASE_MERGE_PLAN"
echo "Entered at: $(date)"
echo "===================="

# Check context
PHASE=$(jq -r '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
echo "Current phase: $PHASE"
echo "Waves completed: $(jq -r '.waves_completed | length' orchestrator-state-v3.json)"

# R340 compliant monitoring
PHASE_ID="phase${PHASE}"
timeout 1800 bash -c 'while true; do
  PLAN_PATH=$(jq -r ".planning_files.merge_plans.phase[\"'"$PHASE_ID"'\"].file_path" orchestrator-state-v3.json)
  if [ "$PLAN_PATH" != "null" ] && [ -f "$PLAN_PATH" ]; then
    echo "Plan tracked: $PLAN_PATH"
    break
  fi
  echo "Checking... $(date +%H:%M:%S)"
  sleep 30
done && echo "✓ Plan created"' || echo "✗ Timeout"

# Final validation
if [ -f phase-*/PHASE-MERGE-PLAN.md ]; then
  echo "Plan size: $(wc -l phase-*/PHASE-MERGE-PLAN.md)"
  echo "Waves in plan: $(grep -c "wave-" phase-*/PHASE-MERGE-PLAN.md)"
fi
```

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

### 🚨 WAITING STATE PATTERN - CRITICAL UNDERSTANDING 🚨

**This is a WAITING state. Common source of incorrect FALSE usage!**

**WRONG interpretation:**
> "R322 mandates stop before transition"
> "State work is complete (validation done)"
> "User needs to invoke /continue-orchestrating"
> "Therefore I must set CONTINUE-SOFTWARE-FACTORY=FALSE"

**CORRECT interpretation:**
> "R322 checkpoint is NORMAL procedure for context preservation"
> "State work completed successfully = NORMAL outcome"
> "Waiting for /continue is DESIGNED user experience"
> "System KNOWS next step from state file"
> "NO manual intervention required, just normal continuation"
> "Therefore set CONTINUE-SOFTWARE-FACTORY=TRUE"

**The key distinction:**
- **Stopping inference** (`exit 0`) = Context management (ALWAYS at R322 points)
- **Continuation flag** = Can automation proceed? (TRUE unless catastrophic failure)

**ONLY use FALSE if:**
- ❌ The thing we're waiting for completely disappeared (agents crashed with no recovery)
- ❌ Results arrived but are completely corrupted/unreadable
- ❌ State file corruption prevents determining what to wait for
- ❌ System deadlock with no automated resolution

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**

## References
- R232: rule-library/R232-enforcement-examples.md
- R233: rule-library/R233-all-states-immediate-action.md
- R269: rule-library/R269-code-reviewer-merge-plan-no-execution.md
- R285: rule-library/R285-mandatory-phase-integration-before-assessment.md
- R287: rule-library/R287-todo-persistence-comprehensive.md
- R290: rule-library/R290-state-rule-reading-verification-supreme-law.md
- R322: rule-library/R322-mandatory-stop-before-state-transitions.md
- R340: rule-library/R340-planning-file-metadata-tracking.md
- R405: rule-library/R405-automation-continuation-flag.md
