# Orchestrator - WAITING_FOR_FIX_PLANS State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.json with new state
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

### PRIMARY DIRECTIVES - MANDATORY READING:

### 🛑 RULE R322 - Mandatory Stop Before State Transitions
**Source:** $CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md
**Criticality:** 🔴🔴🔴 SUPREME LAW - Violation = -100% FAILURE

After completing state work and committing state file:
1. STOP IMMEDIATELY
2. Do NOT continue to next state
3. Do NOT start new work
4. Exit and wait for user
---

**USE THESE EXACT READ COMMANDS (IN THIS ORDER):**
1. Read: $CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md
2. Read: $CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md
3. Read: $CLAUDE_PROJECT_DIR/rule-library/R290-state-rule-reading-verification-supreme-law.md
5. Read: $CLAUDE_PROJECT_DIR/rule-library/R239-fix-plan-distribution.md
6. Read: $CLAUDE_PROJECT_DIR/rule-library/R008-monitoring-frequency.md
7. Read: $CLAUDE_PROJECT_DIR/rule-library/R206-state-machine-transition-validation.md

**WE ARE WATCHING EACH READ TOOL CALL**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R234, R208, R290..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all WAITING_FOR_FIX_PLANS rules"
   (YOU Must READ AND ACKNOWLEDGE EACH rule individually)
   ```

3. **Silent Reading**:
   ```
   ❌ WRONG: [Reads rules but doesn't acknowledge]
   "Now I've read the rules, let me start work..."
   (MUST explicitly acknowledge EACH rule)
   ```

4. **Reading From Memory**:
   ```
   ❌ WRONG: "I know R239 requires fix plan distribution..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR WAITING_FOR_FIX_PLANS:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md
2. "I acknowledge R234 - Mandatory State Traversal: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md  
4. "I acknowledge R006 - Orchestrator Never Writes Code: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. Create verification marker
6. "Ready to execute WAITING_FOR_FIX_PLANS work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY waiting work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ Verification marker has been created
4. ✅ You have stated readiness to execute WAITING_FOR_FIX_PLANS work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY waiting work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

---

## 🔴🔴🔴 SUPREME DIRECTIVE: MONITOR FIX PLAN CREATION 🔴🔴🔴

**WAIT FOR CODE REVIEWER TO COMPLETE FIX PLANS!**

## State Overview

In WAITING_FOR_FIX_PLANS, you monitor the Code Reviewer's progress in creating fix plans for integration failures.

## Required Actions

### 1. Check for Fix Plan Summary
```bash
PHASE=$(jq '.current_phase' orchestrator-state.json)
WAVE=$(jq '.current_wave' orchestrator-state.json)
FIX_PLAN_DIR="efforts/phase${PHASE}/wave${WAVE}/fix-plans"
SUMMARY_FILE="${FIX_PLAN_DIR}/FIX_PLAN_SUMMARY.yaml"

echo "⏳ Checking for fix plan summary at: $SUMMARY_FILE"

if [ -f "$SUMMARY_FILE" ]; then
    echo "✅ Fix plan summary found!"
    
    # Parse the summary
    TOTAL_EFFORTS=$(jq '.total_efforts' "$SUMMARY_FILE")
    echo "Total efforts with fix plans: $TOTAL_EFFORTS"
    
    # Verify all fix plan files exist
    ALL_PLANS_EXIST=true
    while IFS= read -r plan_file; do
        FULL_PATH="${FIX_PLAN_DIR}/${plan_file}"
        if [ ! -f "$FULL_PATH" ]; then
            echo "❌ Missing fix plan file: $FULL_PATH"
            ALL_PLANS_EXIST=false
        else
            echo "✅ Found fix plan: $plan_file"
        fi
    done < <(jq '.fix_plans_created[].plan_file' "$SUMMARY_FILE")
    
    if [ "$ALL_PLANS_EXIST" = true ]; then
        echo "✅ All fix plans created successfully"
        UPDATE_STATE="DISTRIBUTE_FIX_PLANS"
    else
        echo "⚠️ Some fix plans are missing - waiting..."
        sleep 10
        # Stay in WAITING_FOR_FIX_PLANS
    fi
else
    echo "⏳ Fix plan summary not yet created - waiting..."
    
    # Check if Code Reviewer is still working
    if pgrep -f "code-reviewer.*fix" > /dev/null; then
        echo "Code Reviewer still working on fix plans..."
        sleep 10
        # Stay in WAITING_FOR_FIX_PLANS
    else
        # Check timeout
        SPAWN_TIME=$(jq '.integration_feedback.wave'${WAVE}'.fix_plan_requested' orchestrator-state.json)
        CURRENT_TIME=$(date +%s)
        SPAWN_TIMESTAMP=$(date -d "$SPAWN_TIME" +%s 2>/dev/null || echo 0)
        ELAPSED=$((CURRENT_TIME - SPAWN_TIMESTAMP))
        
        if [ $ELAPSED -gt 600 ]; then  # 10 minute timeout
            echo "❌ Timeout waiting for fix plans (>10 minutes)"
            UPDATE_STATE="ERROR_RECOVERY"
        else
            echo "Waiting for Code Reviewer to complete fix plans..."
            sleep 10
            # Stay in WAITING_FOR_FIX_PLANS
        fi
    fi
fi
```

### 2. Update State When Ready
```bash
if [ -n "$UPDATE_STATE" ]; then
    # Record fix plans in state
    if [ "$UPDATE_STATE" = "DISTRIBUTE_FIX_PLANS" ]; then
        jq ".integration_feedback.wave${WAVE}.fix_plans_completed = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state.json
        jq '.integration_feedback.wave${WAVE}.total_fix_plans = $TOTAL_EFFORTS' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    fi
    
    # Update state
    jq ".current_state = \"$UPDATE_STATE\"" -i orchestrator-state.json
    jq ".state_transition_history += [{\"from\": \"WAITING_FOR_FIX_PLANS\", \"to\": \"$UPDATE_STATE\", \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\", \"reason\": \"Fix plans ready for distribution\"}]" -i orchestrator-state.json
    
    # Commit
    git add orchestrator-state.json
    git commit -m "state: Fix plans complete - transitioning to $UPDATE_STATE"
    git push
fi
```

## Valid Transitions

1. **SUCCESS Path**: `WAITING_FOR_FIX_PLANS` → `DISTRIBUTE_FIX_PLANS`
   - When: All fix plans created successfully
   
2. **TIMEOUT Path**: `WAITING_FOR_FIX_PLANS` → `ERROR_RECOVERY`
   - When: Fix plan creation exceeds timeout (10 minutes)
   
3. **CONTINUE Path**: `WAITING_FOR_FIX_PLANS` → `WAITING_FOR_FIX_PLANS`
   - When: Still waiting for fix plans to complete

## Monitoring Requirements

1. Check for FIX_PLAN_SUMMARY.yaml file
2. Verify all referenced fix plan files exist
3. Monitor Code Reviewer process status
4. Track timeout conditions
5. Transition when all plans ready

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

## Related Rules

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
