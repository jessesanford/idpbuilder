# Orchestrator - SPAWN_ARCHITECT_WAVE_PLANNING State Rules

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

## 🔴🔴🔴 R322 MANDATORY: STOP BEFORE STATE TRANSITIONS 🔴🔴🔴

**CRITICAL REQUIREMENT PER R322:**
After spawning ANY agents in this state, you MUST:
1. Record what was spawned in state file
2. Save TODOs per R287
3. Commit and push state changes
4. Display stop message with continuation instructions
5. EXIT immediately with code 0

**VIOLATION = AUTOMATIC -100% FAILURE**

See: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`

---

# Orchestrator - SPAWN_ARCHITECT_WAVE_PLANNING State Rules

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_ARCHITECT_WAVE_PLANNING STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_ARCHITECT_WAVE_PLANNING
echo "$(date +%s) - Rules read and acknowledged for SPAWN_ARCHITECT_WAVE_PLANNING" > .state_rules_read_orchestrator_SPAWN_ARCHITECT_WAVE_PLANNING
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWN_ARCHITECT_WAVE_PLANNING WORK UNTIL RULES ARE READ:
- ❌ Start spawn architect agent
- ❌ Start request wave planning
- ❌ Start generate wave strategy
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
   ❌ WRONG: "I acknowledge all SPAWN_ARCHITECT_WAVE_PLANNING rules"
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

### ✅ CORRECT PATTERN FOR SPAWN_ARCHITECT_WAVE_PLANNING:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute SPAWN_ARCHITECT_WAVE_PLANNING work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SPAWN_ARCHITECT_WAVE_PLANNING work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute SPAWN_ARCHITECT_WAVE_PLANNING work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with SPAWN_ARCHITECT_WAVE_PLANNING work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SPAWN_ARCHITECT_WAVE_PLANNING work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 🚨 SPAWN_ARCHITECT_WAVE_PLANNING IS A VERB - START SPAWNING ARCHITECT IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING SPAWN_ARCHITECT_WAVE_PLANNING

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Spawn Architect with Wave planning request NOW
2. Include wave requirements in spawn message
3. Check TodoWrite for pending items and process them
4. Provide phase plan context immediately

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in SPAWN_ARCHITECT_WAVE_PLANNING" [stops]
- ❌ "Successfully entered SPAWN_ARCHITECT_WAVE_PLANNING state" [waits]
- ❌ "Ready to start spawning architect" [pauses]
- ❌ "I'm in SPAWN_ARCHITECT_WAVE_PLANNING state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering SPAWN_ARCHITECT_WAVE_PLANNING, Spawn Architect with Wave planning request NOW..."
- ✅ "START SPAWNING ARCHITECT, include wave requirements in spawn message..."
- ✅ "SPAWN_ARCHITECT_WAVE_PLANNING: Provide phase plan context immediately..."

## State Context
This is the SPAWN_ARCHITECT_WAVE_PLANNING state for the orchestrator.

## Acknowledgment Required
Thank you for reading the rules file for the SPAWN_ARCHITECT_WAVE_PLANNING state.

**IMPORTANT**: Please report that you have successfully read the SPAWN_ARCHITECT_WAVE_PLANNING rules file.

Say: "✅ Successfully read SPAWN_ARCHITECT_WAVE_PLANNING rules for orchestrator"

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

### R304: Mandatory Line Counter Tool Enforcement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
**Criticality**: BLOCKING - Manual counting = AUTOMATIC -100% FAILURE

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
Proceed with the standard workflow for the SPAWN_ARCHITECT_WAVE_PLANNING state as defined in the state machine.

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
