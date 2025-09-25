# WAITING_FOR_MERGE_PLAN State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_MERGE_PLAN STATE

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
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-requirements.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state.json before EVERY state transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322 Part A-mandatory-stop-after-spawn.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

5. **🔴🔴🔴 R233** - Immediate Action On State Entry
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-immediate-action-on-state-entry.md`
   - Criticality: SUPREME LAW - Must act immediately on entering state
   - Summary: WAITING states require active monitoring, not passive waiting

## State Purpose
Actively monitor Code Reviewer creating wave merge plan. Read merge plan location from orchestrator-state.json per R340. Validate the merge plan meets requirements when found.

## Critical Rules

### 🛑🛑🛑 RULE R322: MANDATORY CHECKPOINT BEFORE SPAWN_INTEGRATION_AGENT (SUPREME LAW) 🛑🛑🛑

**THIS IS A CRITICAL R322 CHECKPOINT STATE!**

When transitioning from WAITING_FOR_MERGE_PLAN → SPAWN_INTEGRATION_AGENT:
- **MUST STOP** to allow user review of WAVE-MERGE-PLAN.md
- **MUST UPDATE** state file to SPAWN_INTEGRATION_AGENT before stopping
- **MUST DISPLAY** checkpoint message with plan location
- **MUST EXIT** cleanly to preserve context
- **VIOLATION = -100% IMMEDIATE FAILURE**

See: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`

### 🔴🔴🔴 RULE R233: IMMEDIATE ACTION REQUIRED (SUPREME LAW)
- **NO PASSIVE WAITING** - Must actively check for completion
- **IMMEDIATE ACTION** - Start checking within first response
- **CONTINUOUS MONITORING** - Check every 30-60 seconds
- **States are VERBS** - "WAITING" means "ACTIVELY CHECKING"

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
- **MUST** read merge plan location from orchestrator-state.json
- **NEVER** search directories for planning files
- **ALWAYS** use planning_files.merge_plans.wave section
- **VIOLATION = -20% for each untracked file**

### ⚠️⚠️⚠️ RULE R269: MERGE PLAN VALIDATION (WARNING)
- Merge plan MUST exist as WAVE-MERGE-PLAN.md
- Plan MUST list all effort branches in order
- Plan MUST specify merge strategy
- Plan MUST identify potential conflicts

### ⚠️⚠️⚠️ RULE R287: TODO PERSISTENCE (WARNING)
- **MUST** save TODOs every 10 messages or 15 minutes
- **MUST** save before state transition
- **MUST** commit and push TODO state

## Required Actions

1. **Initial Check (IMMEDIATE)**
   ```bash
   # Check if Code Reviewer was spawned
   grep "spawned_agents" orchestrator-state.json
   
   # Verify wave integration directory
   ls -la wave-*-integration/
   ```

2. **Active Monitoring Loop (R340 Compliant)**
   ```bash
   # Per R340: Read merge plan location from state file
   PHASE=$(jq -r '.current_phase' orchestrator-state.json)
   WAVE=$(jq -r '.current_wave' orchestrator-state.json)
   WAVE_ID="phase${PHASE}_wave${WAVE}"
   
   # Monitor for plan creation in state file
   while true; do
     # R340: Check if merge plan is tracked in state
     MERGE_PLAN_PATH=$(jq -r ".planning_files.merge_plans.wave[\"${WAVE_ID}\"].file_path // null" orchestrator-state.json)
     
     if [ "$MERGE_PLAN_PATH" != "null" ] && [ -n "$MERGE_PLAN_PATH" ]; then
       # Plan is tracked in state - verify it exists
       if [ -f "$MERGE_PLAN_PATH" ]; then
         echo "✓ Merge plan detected at $(date)"
         echo "📍 Location (from state): $MERGE_PLAN_PATH"
         break
       else
         echo "❌ CRITICAL: Plan tracked but file missing: $MERGE_PLAN_PATH"
         # Transition to ERROR_RECOVERY
         break
       fi
     fi
     
     # Also check if Code Reviewer reported completion
     REVIEWER_STATE=$(jq -r '.spawned_agents[] | select(.name == "code-reviewer") | .state // "UNKNOWN"' orchestrator-state.json)
     
     if [ "$REVIEWER_STATE" = "COMPLETED" ]; then
       # Reviewer claims completion but no plan in state
       echo "⚠️ Code Reviewer completed but no plan tracked in state!"
       echo "Waiting for plan metadata update..."
     elif [ "$REVIEWER_STATE" = "BLOCKED" ] || [ "$REVIEWER_STATE" = "ERROR" ]; then
       echo "✗ Code Reviewer blocked/error - need intervention"
       # Transition to ERROR_RECOVERY
       break
     fi
     
     echo "Waiting for merge plan to be tracked in state (R340)... checking again in 30s"
     sleep 30
   done
   ```

3. **Validate Merge Plan Contents (R340 Compliant)**
   ```bash
   # R340: Use the tracked plan path from state
   PHASE=$(jq -r '.current_phase' orchestrator-state.json)
   WAVE=$(jq -r '.current_wave' orchestrator-state.json)
   WAVE_ID="phase${PHASE}_wave${WAVE}"
   
   # Get plan path from state (R340 requirement)
   PLAN_FILE=$(jq -r ".planning_files.merge_plans.wave[\"${WAVE_ID}\"].file_path" orchestrator-state.json)
   
   if [ "$PLAN_FILE" = "null" ] || [ -z "$PLAN_FILE" ]; then
     echo "❌ R340 VIOLATION: No merge plan tracked in state!"
     exit 340
   fi
   
   if [ ! -f "$PLAN_FILE" ]; then
     echo "❌ Plan tracked but file missing: $PLAN_FILE"
     exit 1
   fi
   
   # Check plan has required sections
   grep -q "## Merge Order" "$PLAN_FILE" || echo "✗ Missing merge order"
   grep -q "## Effort Branches" "$PLAN_FILE" || echo "✗ Missing effort list"
   grep -q "## Merge Strategy" "$PLAN_FILE" || echo "✗ Missing strategy"
   grep -q "## Potential Conflicts" "$PLAN_FILE" || echo "✗ Missing conflict analysis"
   
   # Count effort branches in plan
   EFFORT_COUNT=$(grep -c "effort-" "$PLAN_FILE")
   echo "Plan includes $EFFORT_COUNT efforts"
   ```

4. **Check for Timeout**
   ```bash
   # Get spawn time
   SPAWN_TIME=$(grep "spawned_agents" orchestrator-state.json -A 5 | grep "timestamp" | tail -1 | cut -d'"' -f2)
   
   # Calculate elapsed time
   ELAPSED=$(($(date +%s) - $(date -d "$SPAWN_TIME" +%s)))
   
   # Timeout after 30 minutes
   if [ $ELAPSED -gt 1800 ]; then
     echo "✗ Timeout waiting for merge plan"
     # Transition to ERROR_RECOVERY
   fi
   ```

5. **Update State When Complete (R340 Compliant)**
   ```yaml
   current_state: WAITING_FOR_MERGE_PLAN
   wave_integration:
     merge_plan: (read from planning_files.merge_plans.wave)
     merge_plan_created: YYYY-MM-DD HH:MM:SS
     plan_validation: PASSED
   # R340: Planning file location already tracked in:
   planning_files:
     merge_plans:
       wave:
         phaseX_waveY:
           file_path: /absolute/path/to/WAVE-MERGE-PLAN.md
           created_by: code-reviewer
           created_at: timestamp
   ```

## Transition Rules

### 🛑🛑🛑 R322 MANDATORY CHECKPOINT - DO NOT SKIP! 🛑🛑🛑

**THIS TRANSITION REQUIRES A MANDATORY USER CHECKPOINT!**

When merge plan is ready and validated:

```bash
# STEP 1: R340 - Get merge plan path from state
PHASE=$(jq -r '.current_phase' orchestrator-state.json)
WAVE=$(jq -r '.current_wave' orchestrator-state.json)
WAVE_ID="phase${PHASE}_wave${WAVE}"
MERGE_PLAN_PATH=$(jq -r ".planning_files.merge_plans.wave[\"${WAVE_ID}\"].file_path" orchestrator-state.json)

if [ "$MERGE_PLAN_PATH" = "null" ] || [ -z "$MERGE_PLAN_PATH" ]; then
    echo "❌ R340 VIOLATION: No merge plan tracked in state!"
    exit 340
fi

if [ ! -f "$MERGE_PLAN_PATH" ]; then
    echo "❌ Cannot transition - tracked plan missing: $MERGE_PLAN_PATH"
    exit 1
fi

# STEP 2: Update state file to SPAWN_INTEGRATION_AGENT
echo "📝 R322: Updating state for checkpoint..."
jq '.current_state = "SPAWN_INTEGRATION_AGENT"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.previous_state = "WAITING_FOR_MERGE_PLAN"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.transition_time = "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
# R340: Reference the tracked location, don't hardcode path
jq --arg plan "$MERGE_PLAN_PATH" '.wave_integration.merge_plan = $plan' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

# STEP 3: Save TODOs per R287
save_todos "R322_CHECKPOINT_MERGE_PLAN_READY"

# STEP 4: Commit state changes
git add orchestrator-state.json todos/*.todo
git commit -m "state: R322 checkpoint - merge plan ready for user review"
git push

# STEP 5: Display R322 checkpoint message
cat << EOF
## 🛑 R322 MANDATORY CHECKPOINT - USER REVIEW REQUIRED

### ✅ Merge Plan Created Successfully:
- Location: $MERGE_PLAN_PATH
- Created by: Code Reviewer
- Status: Ready for review
- Tracked per R340 in state file

### 📊 State Transition Ready:
- Current State: WAITING_FOR_MERGE_PLAN ✅
- Next State: SPAWN_INTEGRATION_AGENT (pending user approval)

### 📋 CRITICAL: Review Required Before Execution
The merge plan has been created but NOT executed. Please review:
1. Merge order and dependencies
2. Conflict resolution strategy
3. Branch selection (splits vs originals)

### ⏸️ STOPPED FOR USER REVIEW
DO NOT proceed without reviewing the merge plan!
To continue after review: /continue-orchestrating
EOF

# STEP 6: EXIT CLEANLY (R322 MANDATORY)
exit 0  # STOP HERE - DO NOT CONTINUE!
```

### ❌ R322 VIOLATION - AUTOMATIC FAILURE:
```bash
# ❌❌❌ NEVER DO THIS - IMMEDIATE -100% FAILURE!
echo "Merge plan ready, spawning integration agent..."
/spawn integration-agent  # R322 VIOLATION!
```

### Valid Next States
- **SPAWN_INTEGRATION_AGENT** - After R322 checkpoint and user review
- **ERROR_RECOVERY** - Timeout or Code Reviewer blocked

### Invalid Transitions
- ❌ Skipping to MONITORING_INTEGRATION without spawn
- ❌ Going back to spawn states
- ❌ Transitioning without merge plan
- ❌ Stopping without updating current_state (CAUSES LOOPS!)

## Common Violations to Avoid

1. **Passive waiting** - Violates R233, must actively check
2. **Not checking TodoWrite** - Violates R232 before transition
3. **Missing plan validation** - Proceeding with invalid plan
4. **Ignoring timeouts** - Waiting forever for blocked agent
5. **Not saving progress** - Violates R287 persistence

## Monitoring Pattern

```bash
# CORRECT (R340 compliant): Check state file for tracked plan
echo "Checking for merge plan tracking at $(date)"
WAVE_ID="phase${PHASE}_wave${WAVE}"
while true; do
  PLAN_PATH=$(jq -r ".planning_files.merge_plans.wave[\"${WAVE_ID}\"].file_path" orchestrator-state.json)
  if [ "$PLAN_PATH" != "null" ] && [ -f "$PLAN_PATH" ]; then
    echo "✓ Merge plan tracked and exists: $PLAN_PATH"
    break
  fi
  echo "Plan not tracked yet, checking again in 30s..."
  sleep 30
done

# WRONG (R340 violation): Searching directories
echo "Looking for merge plan..."
while [ ! -f wave-*/WAVE-MERGE-PLAN.md ]; do  # ❌ R340 VIOLATION!
  sleep 30
done
```

## Verification Commands

```bash
# Verify state entry
echo "Entered WAITING_FOR_MERGE_PLAN at $(date)"
echo "Starting active monitoring for merge plan tracking (R340)"

# Check spawn record
grep "spawned_agents" orchestrator-state.json

# R340 compliant monitoring loop
PHASE=$(jq -r '.current_phase' orchestrator-state.json)
WAVE=$(jq -r '.current_wave' orchestrator-state.json)
WAVE_ID="phase${PHASE}_wave${WAVE}"

for i in {1..60}; do
  PLAN_PATH=$(jq -r ".planning_files.merge_plans.wave[\"${WAVE_ID}\"].file_path" orchestrator-state.json)
  if [ "$PLAN_PATH" != "null" ] && [ -f "$PLAN_PATH" ]; then
    echo "✓ Merge plan tracked and found: $PLAN_PATH"
    wc -l "$PLAN_PATH"
    break
  fi
  echo "Check $i: Plan not tracked yet (R340)"
  sleep 30
done

# Timeout check
if [ "$PLAN_PATH" = "null" ]; then
  echo "✗ Timeout - no merge plan tracked after 30 minutes"
fi
```

## References
- R232: rule-library/R232-monitor-state-requirements.md
- R233: rule-library/R233-all-states-immediate-action.md
- R269: rule-library/R269-merge-plan-requirements.md
- R287: rule-library/R287-todo-persistence-comprehensive.md
- R290: rule-library/R290-state-rule-verification.md
- R322: rule-library/R322-mandatory-stop-before-transition.md
- R340: rule-library/R340-planning-file-metadata-tracking.md

### 🔴🔴🔴 MANDATORY VALIDATION REQUIREMENT 🔴🔴🔴

**Per R288 and R324**: ALL state file updates MUST be validated before commit:

```bash
# After ANY update to orchestrator-state.json:
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state.json || {
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
