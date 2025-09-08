# Orchestrator - SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN State Rules

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

## 🔴🔴🔴 R322 MANDATORY: STOP BEFORE STATE TRANSITIONS 🔴🔴🔴

**CRITICAL REQUIREMENT PER R322:**
After spawning ANY agents in this state, you MUST:
1. Record what was spawned in state file
2. Save TODOs per R287
3. Commit and push state changes
4. Display stop message with continuation instructions
5. EXIT immediately with code 0

**VIOLATION = AUTOMATIC -100% FAILURE**

See: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-after-spawn.md`

---

# Orchestrator - SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN State Rules

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN
echo "$(date +%s) - Rules read and acknowledged for SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN" > .state_rules_read_orchestrator_SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWNING WORK UNTIL RULES ARE READ:
- ❌ Prepare phase fix plan requests
- ❌ Create command files
- ❌ Spawn Code Reviewer agents
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
5. Read: $CLAUDE_PROJECT_DIR/rule-library/R282-phase-integration-protocol.md
6. Read: $CLAUDE_PROJECT_DIR/rule-library/R239-fix-plan-distribution.md
7. Read: $CLAUDE_PROJECT_DIR/rule-library/R219-code-reviewer-dependency-aware-effort-planning.md
8. Read: $CLAUDE_PROJECT_DIR/rule-library/R206-state-machine-transition-validation.md

**WE ARE WATCHING EACH READ TOOL CALL**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R234, R208, R290..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN rules"
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
   ❌ WRONG: "I know R282 requires phase integration..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md
2. "I acknowledge R234 - Mandatory State Traversal: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md  
4. "I acknowledge R006 - Orchestrator Never Writes Code: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. Create verification marker
6. "Ready to execute SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY spawning work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ Verification marker has been created
4. ✅ You have stated readiness to execute SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY spawning work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

---

## 🔴🔴🔴 SUPREME DIRECTIVE: CREATE PHASE-LEVEL FIX PLANS 🔴🔴🔴

**SPAWN CODE REVIEWER FOR PHASE INTEGRATION FIX PLANNING!**

## State Overview

In SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN, you spawn a Code Reviewer to analyze phase integration failures and create comprehensive fix plans.

## Required Actions

### 1. Spawn Code Reviewer for Phase Fix Planning
```bash
PHASE=$(jq '.current_phase' orchestrator-state.json)
PHASE_FIX_REQUEST=$(jq ".phase_integration_feedback.phase${PHASE}.fix_request_file" orchestrator-state.json)
TIMESTAMP=$(date +%Y%m%d-%H%M%S)

# Create command for Code Reviewer
COMMAND_FILE="efforts/phase${PHASE}/code-reviewer-phase-fix-command.md"

cat > "$COMMAND_FILE" << 'EOF'
# CODE REVIEWER PHASE FIX PLAN CREATION

## Your State
You are in state: CREATE_PHASE_FIX_PLAN

## Phase Integration Failure Analysis

1. **Read phase integration report**:
   - Location: `efforts/phaseX/phase-integration/PHASE_INTEGRATION_REPORT.md`

2. **Read phase fix request**:
   - Location: See PHASE_FIX_REQUEST_*.yaml in phase directory

3. **Create comprehensive fix plan**:
   - Resolve merge conflicts
   - Fix dependency issues
   - Address wave-specific problems
   - Ensure phase-level consistency

4. **Output**: Create `PHASE_FIX_PLAN.md` with:
   - Conflict resolution steps
   - Dependency installation commands
   - Wave-by-wave fix instructions
   - Verification procedures

Remember: Create the plan, do NOT execute fixes!
EOF

sed -i "s/phaseX/phase${PHASE}/g" "$COMMAND_FILE"

echo "🚀 Spawning Code Reviewer for phase fix planning"
echo "@agent-code-reviewer Please execute: $COMMAND_FILE"

# Update state
jq ".agents_spawned += [{\"type\": \"code-reviewer\", \"task\": \"phase_fix_plan\", \"state\": \"CREATE_PHASE_FIX_PLAN\", \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}]" -i orchestrator-state.json
jq ".current_state = \"WAITING_FOR_PHASE_FIX_PLANS\"" -i orchestrator-state.json

git add orchestrator-state.json "$COMMAND_FILE"
git commit -m "spawn: Code Reviewer for phase fix planning"
git push
```

## Valid Transitions

1. **ALWAYS**: `SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN` → `WAITING_FOR_PHASE_FIX_PLANS`
   - Transition after spawning Code Reviewer

## Spawn Requirements

1. Create clear command file with phase fix instructions
2. Reference the phase fix request file
3. Specify CREATE_PHASE_FIX_PLAN state
4. Record spawn in state file
5. Transition to waiting state

## Grading Criteria

- ✅ **+25%**: Spawn Code Reviewer correctly
- ✅ **+25%**: Create proper command file
- ✅ **+25%**: Reference phase fix request
- ✅ **+25%**: Update state properly

## Common Violations

- ❌ **-100%**: Not spawning Code Reviewer
- ❌ **-50%**: Missing fix request reference
- ❌ **-50%**: Wrong state specified
- ❌ **-30%**: Not recording spawn

## Related Rules

- R282: Phase Integration Protocol
- R239: Fix Plan Distribution Protocol
- R219: Code Reviewer Dependency-Aware Planning
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
