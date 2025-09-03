# Orchestrator - SUCCESS State Rules

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SUCCESS STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SUCCESS
echo "$(date +%s) - Rules read and acknowledged for SUCCESS" > .state_rules_read_orchestrator_SUCCESS
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SUCCESS WORK UNTIL RULES ARE READ:
- ❌ Start finalize all work
- ❌ Start generate reports
- ❌ Start clean up resources
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
   ❌ WRONG: "I acknowledge all SUCCESS rules"
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

### ✅ CORRECT PATTERN FOR SUCCESS:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute SUCCESS work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SUCCESS work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute SUCCESS work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with SUCCESS work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SUCCESS work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 🛑 SUCCESS IS A TERMINAL STATE - FINAL ACTIONS ONLY! 🛑

### TERMINAL STATE - STOPPING IS ALLOWED HERE

**THIS IS A TERMINAL STATE:**
1. Log final success metrics
2. Complete final logging and cleanup
3. Exit gracefully

**TERMINAL STATES ARE THE ONLY STATES WHERE STOPPING IS PERMITTED**

## 🔴🔴🔴 CRITICAL: SUCCESS Prerequisites 🔴🔴🔴

**SUCCESS can ONLY be reached through:**
- PHASE_COMPLETE → SUCCESS (after architect phase assessment)
- NEVER directly from WAVE_REVIEW
- NEVER without phase-level architect approval
- NEVER bypassing SPAWN_ARCHITECT_PHASE_ASSESSMENT

## State Context
Terminal state - execution ends here.

Project or phase successfully completed. All waves integrated, all tests passing, architect approved at PHASE LEVEL.

## Completion Checklist
✅ All waves in phase completed
✅ All effort sizes compliant (<800 lines)
✅ All code reviews passed
✅ Phase integration branch created and tested
✅ Architect review: APPROVED
✅ No blocking issues

## Final Actions
1. Create phase completion report
2. Update orchestrator-state.yaml with completion timestamp
3. Document metrics (efforts completed, lines of code, test coverage)
4. Archive state for phase
5. Prepare for next phase or project completion

## State Transition
- If more phases → PLANNING (for next phase)
- If all phases complete → Project complete (terminal state)
- If issues found → ERROR_RECOVERY

## Success Criteria
- 100% effort completion
- 100% size compliance
- 100% review pass rate
- Clean integration
- Architect approval
