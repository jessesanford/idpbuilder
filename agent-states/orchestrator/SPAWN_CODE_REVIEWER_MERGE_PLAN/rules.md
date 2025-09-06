# Orchestrator - SPAWN_CODE_REVIEWER_MERGE_PLAN State Rules

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

# Orchestrator - SPAWN_CODE_REVIEWER_MERGE_PLAN State Rules

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_CODE_REVIEWER_MERGE_PLAN STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_CODE_REVIEWER_MERGE_PLAN
echo "$(date +%s) - Rules read and acknowledged for SPAWN_CODE_REVIEWER_MERGE_PLAN" > .state_rules_read_orchestrator_SPAWN_CODE_REVIEWER_MERGE_PLAN
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWN_CODE_REVIEWER_MERGE_PLAN WORK UNTIL RULES ARE READ:
- ❌ Start spawn code reviewer
- ❌ Start request merge strategy
- ❌ Start plan integration approach
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
   ❌ WRONG: "I acknowledge all SPAWN_CODE_REVIEWER_MERGE_PLAN rules"
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

### ✅ CORRECT PATTERN FOR SPAWN_CODE_REVIEWER_MERGE_PLAN:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute SPAWN_CODE_REVIEWER_MERGE_PLAN work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SPAWN_CODE_REVIEWER_MERGE_PLAN work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute SPAWN_CODE_REVIEWER_MERGE_PLAN work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with SPAWN_CODE_REVIEWER_MERGE_PLAN work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SPAWN_CODE_REVIEWER_MERGE_PLAN work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## State Definition
The orchestrator spawns a Code Reviewer agent to create a merge plan for wave integration. The orchestrator has already set up the integration infrastructure (directory and branch).

## Required Actions

### 1. Verify Integration Infrastructure
```bash
# Confirm integration workspace is ready
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
WAVE=$(yq '.current_wave' orchestrator-state.yaml)
INTEGRATION_DIR="/efforts/phase${PHASE}/wave${WAVE}/integration-workspace"

if [ ! -d "$INTEGRATION_DIR" ]; then
    echo "❌ Integration directory not found!"
    exit 1
fi

# Verify integration branch exists
cd "$INTEGRATION_DIR"
CURRENT_BRANCH=$(git branch --show-current)
if [[ ! "$CURRENT_BRANCH" =~ integration ]]; then
    echo "❌ Integration branch not created!"
    exit 1
fi

echo "✅ Integration infrastructure verified"
echo "📁 Integration directory: $INTEGRATION_DIR"
echo "🌿 Integration branch: $CURRENT_BRANCH"
```

### 2. Spawn Code Reviewer for Merge Plan
```bash
# Prepare spawn command for Code Reviewer with actual paths
cat > /tmp/code-reviewer-merge-plan-task.md << EOF
Create WAVE MERGE PLAN for Phase ${PHASE} Wave ${WAVE} integration.

CRITICAL REQUIREMENTS (R269, R270):
1. Use ONLY original effort branches - NO integration branches!
2. Analyze branch bases to determine correct merge order
3. Exclude 'too-large' branches, include only splits
4. Create WAVE-MERGE-PLAN.md with exact merge instructions
5. DO NOT execute merges - only plan them!
6. Document expected conflicts and resolution strategies

CRITICAL LOCATION REQUIREMENT:
- CD to integration directory FIRST: cd ${INTEGRATION_DIR}
- Create WAVE-MERGE-PLAN.md IN the integration directory
- Full path for the file: ${INTEGRATION_DIR}/WAVE-MERGE-PLAN.md

Integration Directory: ${INTEGRATION_DIR}
Target Branch: ${CURRENT_BRANCH}

You are spawned into state: WAVE_MERGE_PLANNING
EOF

# Spawn Code Reviewer
echo "🚀 Spawning Code Reviewer for merge plan creation..."
/spawn code-reviewer WAVE_MERGE_PLANNING "$(cat /tmp/code-reviewer-merge-plan-task.md)"
```

### 3. Update State Tracking
```yaml
# Update orchestrator-state.yaml
integration_status:
  phase: ${PHASE}
  wave: ${WAVE}
  infrastructure_ready: true
  integration_branch: "${CURRENT_BRANCH}"
  integration_dir: "${INTEGRATION_DIR}"
  merge_plan_requested: true
  merge_plan_ready: false
  waiting_for: "code-reviewer-merge-plan"
```

## Transition Rules
- Immediate transition to: WAITING_FOR_MERGE_PLAN
- Cannot skip to: SPAWN_INTEGRATION_AGENT (must wait for plan)
- Must verify infrastructure before spawning

## Success Criteria
- Integration infrastructure verified
- Code Reviewer spawned with correct state
- Clear instructions provided including R269 and R270
- State tracking updated



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

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
