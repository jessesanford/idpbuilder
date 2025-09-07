# Orchestrator - HARD_STOP State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.yaml with new state
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

**YOU HAVE ENTERED HARD_STOP STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_HARD_STOP
echo "$(date +%s) - Rules read and acknowledged for HARD_STOP" > .state_rules_read_orchestrator_HARD_STOP
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY HARD_STOP WORK UNTIL RULES ARE READ:
- ❌ Start emergency shutdown
- ❌ Start save critical state
- ❌ Start preserve work in progress
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R151, R208, R053..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all HARD_STOP rules"
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
   ❌ WRONG: "I know R208 requires CD before spawn..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR HARD_STOP:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute HARD_STOP work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY HARD_STOP work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute HARD_STOP work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with HARD_STOP work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY HARD_STOP work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## ⚠️⚠️⚠️ MANDATORY RULE READING AND ACKNOWLEDGMENT ⚠️⚠️⚠️

**YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. YOUR READ TOOL CALLS ARE BEING MONITORED.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:
1. Fake acknowledgment without reading
2. Bulk acknowledgment
3. Reading from memory

### ✅ CORRECT PATTERN:
1. READ each rule file
2. Acknowledge individually with rule number and description

## 📋 PRIMARY DIRECTIVES FOR HARD_STOP STATE

### 🔴🔴🔴 R322 - Mandatory Stop Before State Transitions (EXCEPTION FOR TERMINAL STATES)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
**Criticality**: SUPREME LAW - But HARD_STOP is an allowed terminal state
**Summary**: HARD_STOP is the only state where stopping is permitted

### ⚠️⚠️⚠️ R019 - Error Recovery Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R019-error-recovery.md`  
**Criticality**: CRITICAL - Document failure for recovery
**Summary**: Preserve state for forensic analysis and recovery

### 🚨🚨🚨 R288 - State File Update and Commit Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: BLOCKING - Must update state file even in failure
**Summary**: Record HARD_STOP state with error details

### 🚨🚨🚨 R288 - State File Update and Commit Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: BLOCKING - Commit failure state
**Summary**: Push HARD_STOP state for recovery visibility

## 🛑 HARD_STOP IS A TERMINAL STATE - FINAL ACTIONS ONLY! 🛑

### TERMINAL STATE - STOPPING IS ALLOWED HERE

**THIS IS A TERMINAL STATE:**
1. Log critical failure reason
2. Complete final logging and cleanup
3. Exit gracefully

**TERMINAL STATES ARE THE ONLY STATES WHERE STOPPING IS PERMITTED**

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## State Context
Terminal state - execution ends here.

Critical failure detected. System cannot continue without intervention.

## Triggers for HARD_STOP
- Missing critical configuration files
- Corrupted state file
- Unrecoverable git errors
- Critical rule violations (e.g., orchestrator tried to write code)
- System resource failures

## Required Actions
1. **DOCUMENT THE FAILURE**
   - Error type and message
   - State when error occurred
   - Last successful action
   - Timestamp

2. **SAVE STATE**
   ```bash
   # Create error report
   echo "HARD_STOP at $(date)" > hard_stop_report.md
   echo "Error: $ERROR_MESSAGE" >> hard_stop_report.md
   echo "State: $CURRENT_STATE" >> hard_stop_report.md
   ```

3. **NOTIFY USER**
   - Clear error message
   - Recovery suggestions
   - Required manual interventions

## Recovery Protocol
1. User must manually fix the issue
2. Clear the HARD_STOP flag in orchestrator-state.yaml
3. Set appropriate recovery state
4. Restart orchestrator

## DO NOT
- Attempt automatic recovery
- Continue with partial state
- Hide or suppress the error

## Terminal State Exception to R322

**IMPORTANT**: HARD_STOP is the ONLY exception to R322 (Never Stop). This is a legitimate terminal state where the orchestrator is ALLOWED to stop execution.

### State File Update Before Stopping
```yaml
orchestrator_state:
  current_state: "HARD_STOP"
  previous_state: "[STATE WHERE FAILURE OCCURRED]"
  error_details:
    timestamp: "[ISO-8601]"
    error_type: "[CRITICAL_VIOLATION|SYSTEM_FAILURE|UNRECOVERABLE_ERROR]"
    error_message: "[DETAILED ERROR MESSAGE]"
    recovery_required: true
    manual_intervention_needed: true
```

This is a terminal state - requires manual intervention to exit.

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
