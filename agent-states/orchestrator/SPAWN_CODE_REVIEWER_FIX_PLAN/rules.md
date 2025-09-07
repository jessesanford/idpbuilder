# Orchestrator - SPAWN_CODE_REVIEWER_FIX_PLAN State Rules

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

See: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-after-spawn.md`

---

# Orchestrator - SPAWN_CODE_REVIEWER_FIX_PLAN State Rules

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_CODE_REVIEWER_FIX_PLAN STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_CODE_REVIEWER_FIX_PLAN
echo "$(date +%s) - Rules read and acknowledged for SPAWN_CODE_REVIEWER_FIX_PLAN" > .state_rules_read_orchestrator_SPAWN_CODE_REVIEWER_FIX_PLAN
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWNING WORK UNTIL RULES ARE READ:
- ❌ Prepare fix plan requests
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
5. Read: $CLAUDE_PROJECT_DIR/rule-library/R239-fix-plan-distribution.md
6. Read: $CLAUDE_PROJECT_DIR/rule-library/R219-code-reviewer-dependency-aware-effort-planning.md
7. Read: $CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md
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
   ❌ WRONG: "I acknowledge all SPAWN_CODE_REVIEWER_FIX_PLAN rules"
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

### ✅ CORRECT PATTERN FOR SPAWN_CODE_REVIEWER_FIX_PLAN:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md
2. "I acknowledge R234 - Mandatory State Traversal: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md  
4. "I acknowledge R006 - Orchestrator Never Writes Code: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. Create verification marker
6. "Ready to execute SPAWN_CODE_REVIEWER_FIX_PLAN work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SPAWN_CODE_REVIEWER_FIX_PLAN work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ Verification marker has been CREATED
4. ✅ You have stated readiness to execute SPAWN_CODE_REVIEWER_FIX_PLAN work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SPAWN_CODE_REVIEWER_FIX_PLAN work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**After reading ALL rules, acknowledge them:**
□ I have read R234 - Mandatory State Traversal (SUPREME LAW #1)
□ I have read R208 - No Implementation in Orchestrator (SUPREME LAW #2)
□ I have read R290 - State Rule Reading and Verification (SUPREME LAW #3)
□ I have read R239 - Fix Plan Distribution Protocol
□ I have read R219 - Code Reviewer Dependency-Aware Planning
□ I have read R006 - Orchestrator Never Writes Code
□ I have read R206 - State Machine Transition Validation

**CRITICAL**: You must have made 8 actual Read tool calls. Count them!

---

## 🔴🔴🔴 SUPREME DIRECTIVE: CREATE FIX PLANS FOR FAILED INTEGRATIONS 🔴🔴🔴

**SPAWN CODE REVIEWER TO ANALYZE AND CREATE FIX PLANS!**

## State Overview

In SPAWN_CODE_REVIEWER_FIX_PLAN, you spawn a Code Reviewer agent to analyze integration failures and create detailed fix plans for engineers.

## Required Actions

### 1. Prepare Fix Plan Request
```bash
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
WAVE=$(yq '.current_wave' orchestrator-state.yaml)
FIX_REQUEST_FILE=$(yq ".integration_feedback.wave${WAVE}.fix_request_file" orchestrator-state.yaml)
TIMESTAMP=$(date +%Y%m%d-%H%M%S)

# Create fix plan workspace
FIX_PLAN_DIR="efforts/phase${PHASE}/wave${WAVE}/fix-plans"
mkdir -p "$FIX_PLAN_DIR"

echo "📋 Preparing fix plan request for Code Reviewer"
```

### 2. Spawn Code Reviewer for Fix Planning
```bash
# Create command file for Code Reviewer
COMMAND_FILE="${FIX_PLAN_DIR}/code-reviewer-fix-plan-command.md"

cat > "$COMMAND_FILE" << 'EOF'
# CODE REVIEWER FIX PLAN CREATION TASK

## Your State
You are in state: CREATE_INTEGRATION_FIX_PLAN

## Integration Failure Analysis Required

1. **Read the integration report**: 
   - Location: `efforts/phaseX/waveY/integration-workspace/INTEGRATION_REPORT.md`
   
2. **Read the fix request metadata**:
   - Location: `efforts/phaseX/waveY/FIX_REQUEST_*.yaml`

3. **For each failed effort, create a fix plan**:
   - Analyze the specific failure
   - Determine root cause
   - Create step-by-step fix instructions
   - Include dependency installation if needed
   - Specify test commands to verify fixes

4. **Output Format**:
   Create: `efforts/phaseX/waveY/fix-plans/FIX_PLAN_[effort].md`
   
   ```markdown
   # Fix Plan for [Effort Name]
   
   ## Issue Summary
   [Brief description of what failed]
   
   ## Root Cause
   [Analysis of why it failed]
   
   ## Fix Instructions
   1. [Step 1]
   2. [Step 2]
   ...
   
   ## Verification Steps
   1. Run: `[test command]`
   2. Verify: [what to check]
   
   ## Dependencies to Install (if any)
   - [dependency 1]
   - [dependency 2]
   ```

5. **Create summary file**: `FIX_PLAN_SUMMARY.yaml`
   ```yaml
   fix_plans_created:
     - effort: effort1
       plan_file: FIX_PLAN_effort1.md
       estimated_time: 30m
     - effort: effort2
       plan_file: FIX_PLAN_effort2.md
       estimated_time: 45m
   total_efforts: 2
   timestamp: [ISO timestamp]
   ```

Remember: You are NOT executing fixes, only creating the plans!
EOF

# Replace placeholders
sed -i "s/phaseX/phase${PHASE}/g" "$COMMAND_FILE"
sed -i "s/waveY/wave${WAVE}/g" "$COMMAND_FILE"

echo "✅ Command file created: $COMMAND_FILE"

# Spawn the Code Reviewer
echo "🚀 Spawning Code Reviewer for fix plan creation"
echo "@agent-code-reviewer Please execute the task in: $COMMAND_FILE"
```

### 3. Update State File
```bash
# Record spawning
yq eval ".agents_spawned += [{\"type\": \"code-reviewer\", \"task\": \"create_fix_plans\", \"state\": \"CREATE_INTEGRATION_FIX_PLAN\", \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\", \"command_file\": \"$COMMAND_FILE\"}]" -i orchestrator-state.yaml

# Update current state
yq eval ".current_state = \"WAITING_FOR_FIX_PLANS\"" -i orchestrator-state.yaml
yq eval ".integration_feedback.wave${WAVE}.fix_plan_requested = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state.yaml

# Commit
git add orchestrator-state.yaml "$COMMAND_FILE"
git commit -m "spawn: Code Reviewer for integration fix plans - wave ${WAVE}"
git push
```

## Valid Transitions

1. **ALWAYS**: `SPAWN_CODE_REVIEWER_FIX_PLAN` → `WAITING_FOR_FIX_PLANS`
   - Transition immediately after spawning

## Spawn Requirements

The Code Reviewer MUST be given:
- Access to integration report
- Fix request metadata file
- Clear output location for fix plans
- Instructions NOT to execute fixes

## Grading Criteria

- ✅ **+25%**: Create proper command file
- ✅ **+25%**: Include all failure information
- ✅ **+25%**: Spawn Code Reviewer correctly
- ✅ **+25%**: Update state file properly

## Common Violations

- ❌ **-100%**: Not spawning Code Reviewer
- ❌ **-50%**: Missing integration report reference
- ❌ **-50%**: Asking Code Reviewer to execute fixes
- ❌ **-30%**: Not creating fix plan directory

## Related Rules

- R239: Fix Plan Distribution Protocol
- R219: Code Reviewer Dependency-Aware Planning
- R006: Orchestrator Never Writes Code
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
