# Orchestrator - WAITING_FOR_IMPLEMENTATION_PLAN State Rules

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


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED WAITING_FOR_IMPLEMENTATION_PLAN STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_WAITING_FOR_IMPLEMENTATION_PLAN
echo "$(date +%s) - Rules read and acknowledged for WAITING_FOR_IMPLEMENTATION_PLAN" > .state_rules_read_orchestrator_WAITING_FOR_IMPLEMENTATION_PLAN
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY WAITING_FOR_IMPLEMENTATION_PLAN WORK UNTIL RULES ARE READ:
- ❌ Start monitor planning progress
- ❌ Start check plan status
- ❌ Start wait for implementation plans
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
   ❌ WRONG: "I acknowledge all WAITING_FOR_IMPLEMENTATION_PLAN rules"
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

### ✅ CORRECT PATTERN FOR WAITING_FOR_IMPLEMENTATION_PLAN:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute WAITING_FOR_IMPLEMENTATION_PLAN work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY WAITING_FOR_IMPLEMENTATION_PLAN work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute WAITING_FOR_IMPLEMENTATION_PLAN work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with WAITING_FOR_IMPLEMENTATION_PLAN work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY WAITING_FOR_IMPLEMENTATION_PLAN work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 🚨 WAITING_FOR_IMPLEMENTATION_PLAN IS A VERB - START ACTIVELY CHECKING IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING WAITING_FOR_IMPLEMENTATION_PLAN

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Check for implementation plan file NOW
2. Poll phase/wave directory every few seconds
3. Check TodoWrite for pending items and process them
4. Verify plan completeness when found

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in WAITING_FOR_IMPLEMENTATION_PLAN" [stops]
- ❌ "Successfully entered WAITING_FOR_IMPLEMENTATION_PLAN state" [waits]
- ❌ "Ready to start actively checking" [pauses]
- ❌ "I'm in WAITING_FOR_IMPLEMENTATION_PLAN state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering WAITING_FOR_IMPLEMENTATION_PLAN, Check for implementation plan file NOW..."
- ✅ "START ACTIVELY CHECKING, poll phase/wave directory every few seconds..."
- ✅ "WAITING_FOR_IMPLEMENTATION_PLAN: Verify plan completeness when found..."

## State Context
This is the WAITING_FOR_IMPLEMENTATION_PLAN state for the orchestrator.

## Acknowledgment Required
Thank you for reading the rules file for the WAITING_FOR_IMPLEMENTATION_PLAN state.

**IMPORTANT**: Please report that you have successfully read the WAITING_FOR_IMPLEMENTATION_PLAN rules file.

Say: "✅ Successfully read WAITING_FOR_IMPLEMENTATION_PLAN rules for orchestrator"

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## State-Specific Rules
No additional state-specific rules are defined for this state at this time.

## General Responsibilities
Follow all general orchestrator rules and the Software Factory state machine.

## Next Steps
Proceed with the standard workflow for the WAITING_FOR_IMPLEMENTATION_PLAN state as defined in the state machine.

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
