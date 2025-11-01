# Orchestrator - WAITING_FOR_FIX_PLANS State Rules

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
## ✅ FIX PLANS ARE NORMAL OPERATIONS - CONTINUE AUTOMATICALLY

**IMPORTANT CLARIFICATION PER R322 AND R405:**

Fix plans are NORMAL software development operations that should proceed automatically:
- ✅ NO user review required for fix plans
- ✅ Continue automatically to CREATE_WAVE_FIX_PLAN
- ✅ Set CONTINUE-SOFTWARE-FACTORY=TRUE when plans are ready
- ✅ This is routine error correction, not an exceptional situation

### When transitioning from WAITING_FOR_FIX_PLANS → CREATE_WAVE_FIX_PLAN:
```markdown
## ✅ Fix Plans Ready - Proceeding Automatically

### Plan Details:
- Total plans: [Number of fix plans]
- Affected efforts: [List all efforts needing fixes]
- Fix complexity: [Simple/Complex/Critical]

### Mandatory Stop Protocol:
- Current State: WAITING_FOR_FIX_PLANS ✅
- Next State: CREATE_WAVE_FIX_PLAN
- Action: Stop with TRUE flag after updating state
- STOP PROCESSING (mandatory for context preservation)

CONTINUE-SOFTWARE-FACTORY=TRUE
```

**Fix plans are NORMAL - Stop with TRUE flag for automatic restart!**

See: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
See: `$CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md`

---


## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED WAITING_FOR_FIX_PLANS STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_WAITING_FOR_FIX_PLANS
echo "$(date +%s) - Rules read and acknowledged for WAITING_FOR_FIX_PLANS" > .state_rules_read_orchestrator_WAITING_FOR_FIX_PLANS
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY WAITING WORK UNTIL RULES ARE READ:
- ❌ Check for fix plan summaries
- ❌ Verify fix plan files
- ❌ Monitor Code Reviewer progress
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**
### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_FIX_PLANS STATE

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


### PRIMARY DIRECTIVES - MANDATORY READING:

### 🛑 RULE R322 - Mandatory Stop Before State Transitions
**Source:** $CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md
**Criticality:** 🔴🔴🔴 SUPREME LAW - Violation = -100% FAILURE

After completing state work and committing state file:
1. STOP IMMEDIATELY
2. Do NOT continue to next state
3. Do NOT start new work
4. Exit and wait for user

### 🚨🚨🚨 RULE R533 - Artifact Location Reporting Protocol (BLOCKING)
**Source:** $CLAUDE_PROJECT_DIR/rule-library/R533-artifact-location-reporting-protocol.md
**Criticality:** 🚨🚨🚨 BLOCKING - Violation = -20% per untracked artifact

- **MUST** read artifact locations from orchestrator-state-v3.json
- **ALWAYS** use artifacts.fix_plans section (R533 schema)
- **ALL** artifacts must be tracked with complete metadata
- **VERIFICATION** required before proceeding

### 🚨🚨🚨 RULE R340 - Planning File Metadata Tracking (BLOCKING)
**Source:** $CLAUDE_PROJECT_DIR/rule-library/R340-planning-file-metadata-tracking.md
**Criticality:** 🚨🚨🚨 BLOCKING - Violation = -20% per untracked file

- **NEVER** search directories for planning files
- **ALWAYS** read from state file metadata (now .artifacts per R533)
- **ALL** fix plans must be tracked with metadata

## 🔴🔴🔴 SUPREME DIRECTIVE: MONITOR FIX PLAN CREATION 🔴🔴🔴

**WAIT FOR CODE REVIEWER TO COMPLETE FIX PLANS!**

## State Overview

In WAITING_FOR_FIX_PLANS, you monitor the Code Reviewer's progress in creating fix plans for integration failures.

## Required Actions

### 1. Check for Fix Plans (R340 Compliant)
```bash
# Per R340: Monitor fix plans through state file metadata
PHASE=$(jq -r '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
WAVE=$(jq -r '.project_progression.current_wave.wave_number' orchestrator-state-v3.json)

echo "📊 Checking for fix plans in state file (R340 compliant)"

# R533/R340: Check if fix plans are tracked in state
FIX_PLAN_COUNT=$(jq '.artifacts.fix_plans | length // 0' orchestrator-state-v3.json)

if [ "$FIX_PLAN_COUNT" -gt 0 ]; then
    echo "✅ Found $FIX_PLAN_COUNT fix plans tracked in state"
    
    # R340: Verify all tracked fix plans exist
    ALL_PLANS_EXIST=true
    MISSING_PLANS=()
    
    # Iterate through tracked fix plans (R533 schema)
    jq -r '.artifacts.fix_plans | to_entries[] | @json' orchestrator-state-v3.json | while IFS= read -r entry; do
        ARTIFACT_ID=$(echo "$entry" | jq -r '.key')
        PLAN_PATH=$(echo "$entry" | jq -r '.value.file_path')
        
        if [ -f "$PLAN_PATH" ]; then
            echo "✅ Fix plan exists: $ARTIFACT_ID at $PLAN_PATH"
        else
            echo "❌ CRITICAL: Tracked plan missing: $PLAN_PATH"
            ALL_PLANS_EXIST=false
            MISSING_PLANS+=("$ARTIFACT_ID")
        fi
    done
    
    # Check if all expected efforts have fix plans
    EFFORTS_WITH_FAILURES=$(jq -r '.integration_feedback.wave'"$WAVE"'.efforts_with_failures[]' orchestrator-state-v3.json 2>/dev/null)
    
    if [ -n "$EFFORTS_WITH_FAILURES" ]; then
        echo "Checking if all failed efforts have fix plans..."
        while IFS= read -r effort; do
            # R533: Check artifacts.fix_plans for this effort
            FIX_PLAN_PATH=$(jq -r ".artifacts.fix_plans[\"${effort}-fix-001\"].file_path // null" orchestrator-state-v3.json)
            if [ "$FIX_PLAN_PATH" = "null" ]; then
                echo "⚠️ No fix plan tracked for effort: $effort"
                ALL_PLANS_EXIST=false
            fi
        done <<< "$EFFORTS_WITH_FAILURES"
    fi
    
    if [ "$ALL_PLANS_EXIST" = true ]; then
        echo "✅ All fix plans tracked and verified"
        UPDATE_STATE="CREATE_WAVE_FIX_PLAN"
    else
        echo "⚠️ Some fix plans missing or not tracked - waiting..."
        sleep 10
        # Stay in WAITING_FOR_FIX_PLANS
    fi
else
    echo "⏳ No fix plans tracked yet (R340) - monitoring..."
    
    # Check Code Reviewer status
    REVIEWER_STATE=$(jq -r '.spawned_agents[] | select(.name == "code-reviewer") | .state // "UNKNOWN"' orchestrator-state-v3.json)
    
    if [ "$REVIEWER_STATE" = "COMPLETED" ]; then
        echo "⚠️ Code Reviewer completed but no fix plans tracked in state!"
        echo "Waiting for planning file metadata update..."
        sleep 10
    elif [ "$REVIEWER_STATE" = "BLOCKED" ] || [ "$REVIEWER_STATE" = "ERROR" ]; then
        echo "❌ Code Reviewer blocked/error - need intervention"
        UPDATE_STATE="ERROR_RECOVERY"
    else
        echo "Code Reviewer state: $REVIEWER_STATE"
        
        # Check timeout
        SPAWN_TIME=$(jq -r '.integration_feedback.wave'"${WAVE}"'.fix_plan_requested // null' orchestrator-state-v3.json)
        if [ "$SPAWN_TIME" != "null" ]; then
            CURRENT_TIME=$(date +%s)
            SPAWN_TIMESTAMP=$(date -d "$SPAWN_TIME" +%s 2>/dev/null || echo 0)
            ELAPSED=$((CURRENT_TIME - SPAWN_TIMESTAMP))
            
            if [ $ELAPSED -gt 600 ]; then  # 10 minute timeout
                echo "❌ Timeout waiting for fix plans (>10 minutes)"
                UPDATE_STATE="ERROR_RECOVERY"
            else
                echo "Waiting for Code Reviewer to create and track fix plans (elapsed: ${ELAPSED}s)..."
                sleep 10
                # Stay in WAITING_FOR_FIX_PLANS
            fi
        else
            echo "Waiting for Code Reviewer to track fix plans in state..."
            sleep 10
        fi
    fi
fi
```

### 2. Active Monitoring Loop (R340 Compliant)
```bash
# R340 compliant monitoring pattern
echo "Starting fix plan monitoring (R340 compliant) at $(date)"
CHECKS=0
MAX_CHECKS=60  # 10 minutes with 10s intervals

while [ $CHECKS -lt $MAX_CHECKS ]; do
    CHECKS=$((CHECKS + 1))
    echo "Check #$CHECKS at $(date +%H:%M:%S)"
    
    # R533/R340: Check state file for tracked fix plans
    FIX_PLAN_COUNT=$(jq '.artifacts.fix_plans | length // 0' orchestrator-state-v3.json)
    
    if [ "$FIX_PLAN_COUNT" -gt 0 ]; then
        echo "✓ $FIX_PLAN_COUNT fix plans tracked in state"
        
        # Verify all exist (R533 schema)
        ALL_EXIST=true
        jq -r '.artifacts.fix_plans[].file_path' orchestrator-state-v3.json | while read -r path; do
            [ ! -f "$path" ] && ALL_EXIST=false && echo "Missing: $path"
        done
        
        if [ "$ALL_EXIST" = true ]; then
            echo "✓ All fix plans verified"
            break
        fi
    fi
    
    # Check reviewer status
    REVIEWER_STATE=$(jq -r '.spawned_agents[] | select(.name == "code-reviewer") | .state' orchestrator-state-v3.json)
    echo "Code Reviewer: $REVIEWER_STATE"
    
    [ "$REVIEWER_STATE" = "ERROR" ] && echo "❌ Reviewer error" && break
    
    sleep 10
done

if [ $CHECKS -eq $MAX_CHECKS ]; then
    echo "✗ Timeout reached"
    UPDATE_STATE="ERROR_RECOVERY"
fi
```

### 2. Prepare for State Transition
```bash
# When fix plans are ready, set transition variables
if [ -n "$UPDATE_STATE" ]; then
    echo "Fix plans complete - preparing for transition to $UPDATE_STATE"

    # Set variables for State Manager (SF 3.0 pattern)
    NEXT_STATE="$UPDATE_STATE"
    TRANSITION_REASON="Fix plans ready for distribution"

    # Record fix plan completion metadata (non-state data)
    if [ "$UPDATE_STATE" = "CREATE_WAVE_FIX_PLAN" ]; then
        echo "Fix plans completed at: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
        echo "Total fix plans: $FIX_PLAN_COUNT"
        # This metadata will be recorded by State Manager during transition
    fi

    echo "✅ Ready to transition - State Manager will handle state file updates"
    echo "   Next state: $NEXT_STATE"
    echo "   Reason: $TRANSITION_REASON"
fi
```

## Valid Transitions

1. **PROJECT_DONE Path**: `WAITING_FOR_FIX_PLANS` → `CREATE_WAVE_FIX_PLAN`
   - When: All fix plans created successfully
   
2. **TIMEOUT Path**: `WAITING_FOR_FIX_PLANS` → `ERROR_RECOVERY`
   - When: Fix plan creation exceeds timeout (10 minutes)
   
3. **CONTINUE Path**: `WAITING_FOR_FIX_PLANS` → `WAITING_FOR_FIX_PLANS`
   - When: Still waiting for fix plans to complete

## Monitoring Requirements (R340 Compliant)

### 🚨🚨🚨 RULE R340: PLANNING FILE METADATA TRACKING (BLOCKING)
- **MUST** read fix plan locations from orchestrator-state-v3.json
- **NEVER** search directories for planning files  
- **ALWAYS** use effort_repo_files.fix_plans section
- **VIOLATION = -20% for each untracked file**

1. Check effort_repo_files.fix_plans in state file (R340)
2. Verify all tracked fix plan files exist
3. Monitor Code Reviewer agent state
4. Track timeout conditions
5. Transition when all plans tracked and verified

### Correct Pattern (R533/R340 Compliant):
```bash
# CORRECT: Read from state file (R533 schema)
FIX_PLANS=$(jq -r '.artifacts.fix_plans' orchestrator-state-v3.json)
for artifact_id in $(jq -r '.artifacts.fix_plans | keys[]' orchestrator-state-v3.json); do
    PLAN_PATH=$(jq -r ".artifacts.fix_plans[\"$artifact_id\"].file_path" orchestrator-state-v3.json)
    echo "Checking tracked plan: $PLAN_PATH"
done

# WRONG: Searching directories
find efforts/ -name "FIX-PLAN-*.md"  # ❌ R340 VIOLATION!
ls fix-plans/  # ❌ R340 VIOLATION!
```

## Grading Criteria

- ✅ **+25%**: Check for fix plan summary correctly
- ✅ **+25%**: Verify all fix plan files
- ✅ **+25%**: Handle timeouts properly
- ✅ **+25%**: Update state appropriately

## Common Violations

- ❌ **-100%**: Not checking for summary file
- ❌ **-50%**: Missing file verification
- ❌ **-50%**: No timeout handling
- ❌ **-30%**: Wrong state transitions



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete WAITING_FOR_FIX_PLANS:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="WAITING_FOR_FIX_PLANS complete - [describe accomplishment]"
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
  "state_completed": "WAITING_FOR_FIX_PLANS",
  "work_accomplished": [
    "Monitored Code Reviewer fix plan creation",
    "Verified all fix plans tracked in state (R340)",
    "Confirmed fix plans ready for distribution"
  ],
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "WAITING_FOR_FIX_PLANS" \
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
save_todos "WAITING_FOR_FIX_PLANS_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - WAITING_FOR_FIX_PLANS complete [R287]"; then
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

**SF 3.0 uses 6-step exit (not 8-step) - State Manager handles state file updates!**

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
PROPOSED_NEXT_STATE="ERROR_RECOVERY"
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


## Related Rules

- R340: Planning File Metadata Tracking (CRITICAL)
- R239: Fix Plan Distribution Protocol
- R008: Monitoring Frequency
- R206: State Machine Transition Validation

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**


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

