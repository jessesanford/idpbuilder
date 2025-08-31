# Orchestrator - SPAWN_CODE_REVIEWER_PHASE_IMPL State Rules

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_CODE_REVIEWER_PHASE_IMPL STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_CODE_REVIEWER_PHASE_IMPL
echo "$(date +%s) - Rules read and acknowledged for SPAWN_CODE_REVIEWER_PHASE_IMPL" > .state_rules_read_orchestrator_SPAWN_CODE_REVIEWER_PHASE_IMPL
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWN_CODE_REVIEWER_PHASE_IMPL WORK UNTIL RULES ARE READ:
- ❌ Start spawn code reviewer
- ❌ Start review phase implementation
- ❌ Start validate phase work
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
   ❌ WRONG: "I acknowledge all SPAWN_CODE_REVIEWER_PHASE_IMPL rules"
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

### ✅ CORRECT PATTERN FOR SPAWN_CODE_REVIEWER_PHASE_IMPL:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute SPAWN_CODE_REVIEWER_PHASE_IMPL work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SPAWN_CODE_REVIEWER_PHASE_IMPL work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute SPAWN_CODE_REVIEWER_PHASE_IMPL work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with SPAWN_CODE_REVIEWER_PHASE_IMPL work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SPAWN_CODE_REVIEWER_PHASE_IMPL work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 🚨 SPAWN_CODE_REVIEWER_PHASE_IMPL IS A VERB - START SPAWNING CODE REVIEWER IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING SPAWN_CODE_REVIEWER_PHASE_IMPL

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Spawn Code Reviewer to translate Phase architecture NOW
2. Include architecture plan location in spawn
3. Check TodoWrite for pending items and process them
4. Request implementation plan creation immediately

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in SPAWN_CODE_REVIEWER_PHASE_IMPL" [stops]
- ❌ "Successfully entered SPAWN_CODE_REVIEWER_PHASE_IMPL state" [waits]
- ❌ "Ready to start spawning code reviewer" [pauses]
- ❌ "I'm in SPAWN_CODE_REVIEWER_PHASE_IMPL state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering SPAWN_CODE_REVIEWER_PHASE_IMPL, Spawn Code Reviewer to translate Phase architecture NOW..."
- ✅ "START SPAWNING CODE REVIEWER, include architecture plan location in spawn..."
- ✅ "SPAWN_CODE_REVIEWER_PHASE_IMPL: Request implementation plan creation immediately..."

## State Context
This is the SPAWN_CODE_REVIEWER_PHASE_IMPL state for the orchestrator.

## Acknowledgment Required
Thank you for reading the rules file for the SPAWN_CODE_REVIEWER_PHASE_IMPL state.

**IMPORTANT**: Please report that you have successfully read the SPAWN_CODE_REVIEWER_PHASE_IMPL rules file.

Say: "✅ Successfully read SPAWN_CODE_REVIEWER_PHASE_IMPL rules for orchestrator"

## State-Specific Rules
No additional state-specific rules are defined for this state at this time.

## General Responsibilities
Follow all general orchestrator rules and the Software Factory state machine.

## Next Steps
Proceed with the standard workflow for the SPAWN_CODE_REVIEWER_PHASE_IMPL state as defined in the state machine.
