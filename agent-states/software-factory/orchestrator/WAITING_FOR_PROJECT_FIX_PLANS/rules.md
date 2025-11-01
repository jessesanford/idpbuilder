# WAITING_FOR_PROJECT_FIX_PLANS State Rules

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

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_PROJECT_FIX_PLANS STATE

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
   - Summary: WAITING states require active monitoring, not passive waiting

## State Purpose
Actively monitor Code Reviewer creating fix plans for project-level integration bugs documented per R266. These are bugs found during final project integration that need to be fixed in source branches.

## Critical Rules

### 🔴🔴🔴 RULE R322: MANDATORY STOP BEFORE STATE TRANSITION (SUPREME LAW)
- **STOP** and save state before ANY transition
- **READ** orchestrator-state-v3.json to verify current state
- **VALIDATE** next state exists in software-factory-3.0-state-machine.json
- **VIOLATION = IMMEDIATE FAILURE**

### 🔴🔴🔴 RULE R233: IMMEDIATE ACTION REQUIRED (SUPREME LAW)
- **NO PASSIVE WAITING** - Must actively check for completion
- **IMMEDIATE ACTION** - Start checking within first response
- **CONTINUOUS MONITORING_SWE_PROGRESS** - Check every 30-60 seconds
- **States are VERBS** - "WAITING" means "ACTIVELY CHECKING"

### 🔴🔴🔴 RULE R321: IMMEDIATE BACKPORT REQUIRED (SUPREME LAW)
- Project fixes MUST be applied to source branches
- Integration branches are READ-ONLY for code
- Fix plans must target phase/wave/effort branches
- After fixes, project integration MUST be re-run

### 🚨🚨🚨 RULE R290: STATE RULE VERIFICATION (BLOCKING)
- **MUST** verify this rules file exists and is loaded
- **MUST** acknowledge all rules before proceeding
- **MUST** validate state transitions against state machine

### 🚨🚨🚨 RULE R232: MONITOR STATE REQUIREMENTS (BLOCKING)
- **MUST** check TodoWrite for pending items BEFORE transition
- **MUST** process ALL pending items immediately
- **NO** "I will..." statements - only "I am..." with action
- **VIOLATION = AUTOMATIC FAILURE**

### ⚠️⚠️⚠️ RULE R266: BUG DOCUMENTATION REQUIREMENTS (WARNING)
- Fix plans MUST reference bugs documented in PROJECT-INTEGRATE_WAVE_EFFORTS-REPORT.md
- Each bug from R266 documentation MUST have corresponding fix
- Plans MUST maintain bug numbering from original documentation
- Plans MUST specify exact source branch for each fix

### ⚠️⚠️⚠️ RULE R287: TODO PERSISTENCE (WARNING)
- **MUST** save TODOs every 10 messages or 15 minutes
- **MUST** save before state transition
- **MUST** commit and push TODO state

## Required Actions

1. **Initial Check (IMMEDIATE)**
   ```bash
   # Verify Code Reviewer was spawned for fix planning
   grep "spawned_agents" orchestrator-state-v3.json | grep -i "project_fix"
   
   # Check project integration directory
   ls -la project-integration/
   
   # Look for bug documentation per R266
   grep "UPSTREAM BUGS IDENTIFIED" project-integration/PROJECT-INTEGRATE_WAVE_EFFORTS-REPORT.md
   ```

2. **Active Monitoring Loop**
   ```bash
   # Monitor for fix plan creation
   while true; do
     # Check for fix plan file
     if [ -f "project-integration/.software-factory/PROJECT-FIX-PLAN--${TIMESTAMP}.md" ]; then
       echo "✓ Project fix plan detected at $(date)"
       break
     fi
     
     # Check Code Reviewer status if available
     if [ -f "project-integration/fix-planning-status.yaml" ]; then
       STATUS=$(grep "status:" project-integration/fix-planning-status.yaml | awk '{print $2}')
       echo "Fix planning status: $STATUS at $(date)"
       
       if [ "$STATUS" = "COMPLETED" ]; then
         echo "✓ Fix planning completed"
         break
       elif [ "$STATUS" = "BLOCKED" ]; then
         echo "✗ Fix planning blocked - need to handle"
         # May need to transition to ERROR_RECOVERY
         break
       fi
     fi
     
     echo "Waiting for project fix plan... checking again in 30s"
     sleep 30
   done
   ```

3. **Validate Fix Plan**
   ```bash
   PLAN="project-integration/.software-factory/PROJECT-FIX-PLAN--${TIMESTAMP}.md"
   
   if [ -f "$PLAN" ]; then
     echo "Validating project fix plan..."
     
     # Check required sections
     grep -q "## Bug Summary" "$PLAN" || echo "✗ Missing bug summary"
     grep -q "## Fix Strategy" "$PLAN" || echo "✗ Missing fix strategy"
     grep -q "### Bug #" "$PLAN" || echo "✗ Missing individual bug fixes"
     
     # Verify R266 bug references
     BUG_COUNT=$(grep "### Bug #" project-integration/PROJECT-INTEGRATE_WAVE_EFFORTS-REPORT.md | wc -l)
     FIX_COUNT=$(grep "#### Bug #" "$PLAN" | wc -l)
     
     if [ "$BUG_COUNT" != "$FIX_COUNT" ]; then
       echo "✗ ERROR: Bug count mismatch! R266 docs: $BUG_COUNT, Fix plan: $FIX_COUNT"
     fi
     
     # Verify targets source branches (R321)
     if grep -q "project-integration" "$PLAN" | grep -v "after\|from"; then
       echo "✗ ERROR: Fix plan targets integration branch (violates R321)"
     fi
     
     # Check for parallelization strategy
     grep -q "Parallel Fix Group" "$PLAN" && echo "✓ Has parallel fix strategy"
     grep -q "Sequential Fix" "$PLAN" && echo "✓ Has sequential fix handling"
     
     # Verify SW Engineer spawn instructions
     grep -q "## SW Engineer Spawn Instructions" "$PLAN" || \
       echo "⚠ Missing spawn instructions section"
   else
     echo "✗ PROJECT-FIX-PLAN--${TIMESTAMP}.md not found yet"
   fi
   ```

4. **Extract Fix Information**
   ```bash
   # Parse fix plan for orchestrator use
   if [ -f "$PLAN" ]; then
     # Count fixes by type
     PARALLEL_FIXES=$(grep -c "### Parallel Fix Group" "$PLAN")
     SEQUENTIAL_FIXES=$(grep -c "### Sequential Fix" "$PLAN")
     
     # Extract engineer assignments
     ENGINEERS_NEEDED=$(grep "Engineer [0-9]:" "$PLAN" | wc -l)
     
     echo "Fix Plan Summary:"
     echo "- Parallel fix groups: $PARALLEL_FIXES"
     echo "- Sequential fixes: $SEQUENTIAL_FIXES"
     echo "- Engineers needed: $ENGINEERS_NEEDED"
     
     # Update state with plan information
     jq ".project_fixes.plan_received = true" -i orchestrator-state-v3.json
     jq ".project_fixes.parallel_groups = $PARALLEL_FIXES" -i orchestrator-state-v3.json
     jq ".project_fixes.engineers_needed = $ENGINEERS_NEEDED" -i orchestrator-state-v3.json
   fi
   ```

5. **Check for Timeout**
   ```bash
   # Get spawn time
   SPAWN_TIME=$(grep "project_fix_plan" orchestrator-state-v3.json -A 2 | \
                grep "timestamp" | tail -1 | cut -d'"' -f2)
   
   # Calculate elapsed
   ELAPSED=$(($(date +%s) - $(date -d "$SPAWN_TIME" +%s)))
   
   # Timeout after 30 minutes (project fixes should be straightforward)
   if [ $ELAPSED -gt 1800 ]; then
     echo "✗ Timeout waiting for project fix plans"
     # Transition to SPAWN_SW_ENGINEERS with error handling
   fi
   ```

6. **Prepare for Next State**
   ```bash
   # Once plan is ready and validated
   if [ -f "project-integration/.software-factory/PROJECT-FIX-PLAN--${TIMESTAMP}.md" ]; then
     echo "✅ Project fix plan ready"

     # Update tracking fields (ALLOWED - orchestrator maintains this data)
     jq '.project_fixes.plan_location = "project-integration/.software-factory/PROJECT-FIX-PLAN--${TIMESTAMP}.md"' \
        -i orchestrator-state-v3.json

     # Set proposed next state (State Manager will update state_machine fields)
     PROPOSED_NEXT_STATE="SPAWN_SW_ENGINEERS"
     TRANSITION_REASON="Project fix plan received - ready to spawn engineers"
     # State Manager consultation happens in Step 3 of completion checklist
   fi
   ```

## Transition Rules

### Valid Next States
- **SPAWN_SW_ENGINEERS** - Fix plan ready, spawn engineers
- **ERROR_RECOVERY** - Planning failed or blocked

### Invalid Transitions
- ❌ Direct to MONITORING_EFFORT_FIXES (must spawn engineers first)
- ❌ Back to PROJECT_INTEGRATE_WAVE_EFFORTS without fixes
- ❌ Skipping fix implementation



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete WAITING_FOR_PROJECT_FIX_PLANS:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="WAITING_FOR_PROJECT_FIX_PLANS complete - [accomplishment description]"
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
  "state_completed": "WAITING_FOR_PROJECT_FIX_PLANS",
  "work_accomplished": [
    "Monitored project fix plan creation",
    "Validated fix plan completeness (R266)",
    "Verified R321 compliance (targets source branches)"
  ],
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "WAITING_FOR_PROJECT_FIX_PLANS" \
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
save_todos "WAITING_FOR_PROJECT_FIX_PLANS_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - WAITING_FOR_PROJECT_FIX_PLANS complete [R287]"; then
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
PROPOSED_NEXT_STATE="SPAWN_SW_ENGINEERS"
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
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**


## Common Violations to Avoid

1. **Passive waiting** - Violates R233, must actively monitor
2. **Not verifying R266 compliance** - Plan doesn't match documented bugs
3. **Not checking R321 compliance** - Plans target integration branch
4. **Missing timeout handling** - Waiting forever
5. **Not validating plan completeness** - Missing bugs or instructions
6. **Forgetting TODO persistence** - Violates R287

## Monitoring Pattern

```bash
# CORRECT: Active monitoring with comprehensive checks
echo "Monitoring project fix planning at $(date)"
START_TIME=$(date +%s)

while true; do
  # Check for plan file
  if [ -f "project-integration/.software-factory/PROJECT-FIX-PLAN--${TIMESTAMP}.md" ]; then
    echo "✓ Project fix plan found!"
    
    # Validate it has all bugs from R266
    R266_BUGS=$(grep -c "### Bug #" project-integration/PROJECT-INTEGRATE_WAVE_EFFORTS-REPORT.md)
    PLAN_BUGS=$(grep -c "#### Bug #" project-integration/PROJECT-FIX-PLAN--${TIMESTAMP}.md)
    
    if [ "$R266_BUGS" = "$PLAN_BUGS" ]; then
      echo "✓ All $R266_BUGS bugs addressed in plan"
      break
    else
      echo "⚠ Plan incomplete: $PLAN_BUGS/$R266_BUGS bugs"
    fi
  fi
  
  # Check status file
  test -f project-integration/fix-planning-status.yaml && \
    grep "progress:" project-integration/fix-planning-status.yaml
  
  # Check timeout
  ELAPSED=$(($(date +%s) - START_TIME))
  if [ $ELAPSED -gt 1800 ]; then
    echo "✗ Timeout after 30 minutes"
    break
  fi
  
  echo "Check at $(date): No complete plan yet"
  sleep 30
done

# WRONG: Just waiting without checks
sleep 1800  # Wait 30 minutes
ls project-integration/PROJECT-FIX-PLAN--${TIMESTAMP}.md  # Check once at end
```

## Verification Commands

```bash
# Verify state entry
echo "Entered WAITING_FOR_PROJECT_FIX_PLANS at $(date)"
echo "Project integration found bugs - monitoring fix plan creation"

# Check R266 bug documentation
echo "Bugs documented per R266:"
grep "### Bug #" project-integration/PROJECT-INTEGRATE_WAVE_EFFORTS-REPORT.md | head -5

# Active monitoring with validation
for i in {1..60}; do  # 30 minute timeout
  if [ -f "project-integration/.software-factory/PROJECT-FIX-PLAN--${TIMESTAMP}.md" ]; then
    echo "✓ Fix plan found after $i checks"
    
    # Validate R321 compliance
    if grep -q "fix.*project-integration" project-integration/PROJECT-FIX-PLAN--${TIMESTAMP}.md; then
      echo "✗ R321 VIOLATION: Plan suggests fixing integration branch!"
    else
      echo "✓ R321 compliant: All fixes target source branches"
    fi
    
    break
  fi
  echo "Check $i: No plan yet at $(date)"
  sleep 30
done

# Prepare transition
if [ -f "project-integration/.software-factory/PROJECT-FIX-PLAN--${TIMESTAMP}.md" ]; then
  echo "Transitioning to SPAWN_SW_ENGINEERS"
fi
```

## R322 Checkpoint Requirements

Before transitioning to SPAWN_SW_ENGINEERS:
1. **STOP** and display checkpoint message
2. Show the PROJECT-FIX-PLAN--${TIMESTAMP}.md for user review
3. Summarize fix strategy (parallel vs sequential)
4. List affected source branches
5. Wait for user /continue-orchestrating

```markdown
## 🛑 STATE TRANSITION CHECKPOINT: WAITING_FOR_PROJECT_FIX_PLANS → SPAWN_SW_ENGINEERS

### ✅ Fix Plan Received:
- Location: project-integration/PROJECT-FIX-PLAN--${TIMESTAMP}.md
- Bugs addressed: X/Y from R266 documentation
- Parallel fix groups: N
- Engineers needed: M

### 📋 Fix Strategy Summary:
[Show key parts of fix plan]

### ⏸️ STOPPED - Awaiting User Review
Please review the fix plan before spawning engineers.
Use /continue-orchestrating to proceed.
```

## References
- R232: rule-library/R232-enforcement-examples.md
- R233: rule-library/R233-all-states-immediate-action.md
- R266: rule-library/R266-upstream-bug-documentation.md
- R287: rule-library/R287-todo-persistence-comprehensive.md
- R290: rule-library/R290-state-rule-reading-verification-supreme-law.md
- R321: rule-library/R321-immediate-backport-during-integration.md
- R322: rule-library/R322-mandatory-stop-before-state-transitions.md

### 🔴🔴🔴 MANDATORY VALIDATION REQUIREMENT 🔴🔴🔴

**Per R288 and R324**: ALL state file updates MUST be validated before commit:

```bash
# After ANY update to orchestrator-state-v3.json:
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
```

**Use helper functions for automatic validation:**
```bash
# Source the helper functions
source "$CLAUDE_PROJECT_DIR/utilities/state-file-update-functions.sh"

# Use safe functions that include validation:
safe_state_transition "NEW_STATE" "reason"
safe_update_field "field_name" "value"
```

