# Orchestrator - SPAWN_CODE_REVIEWERS_EFFORT_PLANNING State Rules

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

# Orchestrator - SPAWN_CODE_REVIEWERS_EFFORT_PLANNING State Rules

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_CODE_REVIEWERS_EFFORT_PLANNING STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
echo "$(date +%s) - Rules read and acknowledged for SPAWN_CODE_REVIEWERS_EFFORT_PLANNING" > .state_rules_read_orchestrator_SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWN_CODE_REVIEWERS_EFFORT_PLANNING WORK UNTIL RULES ARE READ:
- ❌ Start spawn code reviewer agents
- ❌ Start request effort plans
- ❌ Start assign efforts to reviewers
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
   ❌ WRONG: "I acknowledge all SPAWN_CODE_REVIEWERS_EFFORT_PLANNING rules"
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

### ✅ CORRECT PATTERN FOR SPAWN_CODE_REVIEWERS_EFFORT_PLANNING:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute SPAWN_CODE_REVIEWERS_EFFORT_PLANNING work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SPAWN_CODE_REVIEWERS_EFFORT_PLANNING work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute SPAWN_CODE_REVIEWERS_EFFORT_PLANNING work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with SPAWN_CODE_REVIEWERS_EFFORT_PLANNING work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SPAWN_CODE_REVIEWERS_EFFORT_PLANNING work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 📋 PRIMARY DIRECTIVES - MANDATORY RULE READING

### Rules You MUST Read Before ANY Work in This State:
1. **READ:** $CLAUDE_PROJECT_DIR/rule-library/R251-REPOSITORY-SEPARATION-LAW.md
   - R251: Universal Repository Separation Law [PARAMOUNT - -100% for violation]
   - **CRITICAL**: Ensure code reviewers understand SF vs Target repo separation
2. **READ:** $CLAUDE_PROJECT_DIR/rule-library/R309-never-create-efforts-in-sf-repo.md
   - R309: Never Create Efforts in SF Repo [PARAMOUNT - -100% for violation]
   - **CRITICAL**: Verify effort directories are TARGET clones before spawning
3. **READ:** $CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md
   - R234: Mandatory State Traversal Supreme Law [SUPREME LAW - Part of sequence]
4. **READ:** $CLAUDE_PROJECT_DIR/rule-library/R208-orchestrator-spawn-directory-protocol.md
   - R208: Orchestrator Spawn Directory Protocol [SUPREME LAW - MANDATORY for every spawn]
5. **READ:** $CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-spawning-timing.md
   - R151: Parallel Spawn Timing [CRITICAL - <5s delta required]
6. **READ:** $CLAUDE_PROJECT_DIR/rule-library/R218-orchestrator-parallel-code-reviewer-spawning.md
   - R218: Orchestrator Parallel Code Reviewer Spawning [MANDATORY - From analysis]
7. **READ:** $CLAUDE_PROJECT_DIR/rule-library/R054-implementation-plan-creation.md
   - R054: Implementation Plan Creation [BLOCKING - Understanding deliverables]
   - **Why:** You're spawning reviewers to CREATE these plans - know what they must produce
8. **READ:** $CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md
   - R287: TODO Save Triggers [BLOCKING - Save after spawning]
9. **READ:** $CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md
   - R288: State File Update and Commit Protocol [SUPREME LAW]

## 📋 RULE SUMMARY FOR SPAWN_CODE_REVIEWERS_EFFORT_PLANNING STATE

### Rules Enforced in This State:
- R251: Universal Repository Separation Law [PARAMOUNT - -100% for violation]
- R309: Never Create Efforts in SF Repo [PARAMOUNT - -100% for violation]
- R234: Mandatory State Traversal Supreme Law [SUPREME LAW - Part of sequence]
- R208: Orchestrator Spawn Directory Protocol [SUPREME LAW - MANDATORY for every spawn]
- R151: Parallel Spawn Timing [CRITICAL - <5s delta required]
- R218: Orchestrator Parallel Code Reviewer Spawning [MANDATORY - From analysis]
- R054: Implementation Plan Creation [BLOCKING - Know the deliverable]
- R287: TODO Save Triggers [BLOCKING - Save after spawning]
- R288: State File Update and Commit Protocol [SUPREME LAW]

### Critical Requirements:
1. Use parallelization plan from state file - Penalty: -50%
2. CD to correct directory for EACH spawn - Penalty: -100%
3. Spawn parallel reviewers in ONE message - Penalty: -50%
4. Save spawn times to state file - Penalty: -20%
5. Transition to WAITING_FOR_EFFORT_PLANS - Penalty: -100%

### Success Criteria:
- ✅ Parallelization plan loaded from state
- ✅ R208 CD protocol followed for each spawn
- ✅ All parallel spawns <5s delta
- ✅ Spawn times recorded
- ✅ All Code Reviewers spawned per plan

### Failure Triggers:
- ❌ Spawn without CD = -100% R208 VIOLATION
- ❌ Skip to SPAWN_AGENTS = -100% R234 VIOLATION
- ❌ Ignore parallelization plan = -50% penalty
- ❌ Spawn sequentially when should be parallel = -50%

## 🚨 SPAWN_CODE_REVIEWERS_EFFORT_PLANNING IS A VERB - START SPAWNING CODE REVIEWERS IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING SPAWN_CODE_REVIEWERS_EFFORT_PLANNING

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Spawn Code Reviewer for first effort NOW per parallelization plan
2. Use spawn sequence from orchestrator-state.yaml immediately
3. Check TodoWrite for pending items and process them
4. Follow R151 timing requirements for parallel spawns

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in SPAWN_CODE_REVIEWERS_EFFORT_PLANNING" [stops]
- ❌ "Successfully entered SPAWN_CODE_REVIEWERS_EFFORT_PLANNING state" [waits]
- ❌ "Ready to start spawning code reviewers" [pauses]
- ❌ "I'm in SPAWN_CODE_REVIEWERS_EFFORT_PLANNING state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering SPAWN_CODE_REVIEWERS_EFFORT_PLANNING, Spawn Code Reviewer for first effort NOW per parallelization plan..."
- ✅ "START SPAWNING CODE REVIEWERS, use spawn sequence from orchestrator-state.yaml immediately..."
- ✅ "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING: Follow R151 timing requirements for parallel spawns..."

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## State Context
You are spawning Code Reviewers to create individual effort implementation plans in their prepared directories.

## 🔴🔴🔴 SUPREME LAW R234 - MANDATORY SEQUENCE CONTINUES 🔴🔴🔴

### YOUR POSITION IN THE MANDATORY SEQUENCE:
```
SETUP_EFFORT_INFRASTRUCTURE (✓ completed)
    ↓
ANALYZE_CODE_REVIEWER_PARALLELIZATION (✓ completed)
    ↓
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING (👈 YOU ARE HERE)
    ↓ (MUST GO HERE NEXT)
WAITING_FOR_EFFORT_PLANS
    ↓
ANALYZE_IMPLEMENTATION_PARALLELIZATION
    ↓
SPAWN_AGENTS
```

**CRITICAL:** You analyzed parallelization (correct!)
**NOW:** Spawn Code Reviewers using that analysis
**NEXT:** You MUST go to WAITING_FOR_EFFORT_PLANS
**FORBIDDEN:** Skipping to SPAWN_AGENTS = -100% FAILURE

## 🔴🔴🔴 PREREQUISITES 🔴🔴🔴

**BEFORE ENTERING THIS STATE, YOU MUST HAVE:**
1. ✅ Wave Implementation Plan created
2. ✅ ALL effort directories created (SETUP_EFFORT_INFRASTRUCTURE)
3. ✅ ALL branches pushed to remote with tracking
4. ✅ work-log.md files initialized in each directory
5. ✅ **PARALLELIZATION ANALYSIS COMPLETE (ANALYZE_CODE_REVIEWER_PARALLELIZATION)**
6. ✅ **Parallelization plan in orchestrator-state.yaml**

**IF PARALLELIZATION NOT ANALYZED, GO BACK TO ANALYZE_CODE_REVIEWER_PARALLELIZATION!**
**IF INFRASTRUCTURE IS NOT READY, GO BACK TO SETUP_EFFORT_INFRASTRUCTURE!**

### 🚨 RULE R151 - Parallel Spawning with Directory Protocol
**Source:** rule-library/R151-parallel-agent-spawning-timing.md
**Criticality:** CRITICAL - 50% of orchestrator grade

SPAWNING REQUIREMENTS:
1. Check parallelization metadata in Wave Implementation Plan
2. CD to correct directory for EACH Code Reviewer (R208 SUPREME LAW)
3. Output pwd to VERIFY directory before spawn (R208 SUPREME LAW)
4. Spawn all parallel reviewers in ONE message (<5s delta)
5. Include effort-specific context in each spawn
6. Return to orchestrator directory after each spawn (R208 SUPREME LAW)

## 🔴🔴🔴 CRITICAL: Effort Plan Storage Instructions 🔴🔴🔴

### MANDATORY: Instruct Code Reviewers About Plan Location

When spawning Code Reviewers, you MUST explicitly tell them:
```markdown
**CRITICAL INSTRUCTION FOR EFFORT PLAN STORAGE:**
Per R303 and your EFFORT_PLAN_CREATION state rules, you MUST save the effort plan in:
`.software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/IMPLEMENTATION-PLAN-YYYYMMDD-HHMMSS.md`

DO NOT save in the root directory. The orchestrator will look for plans in the .software-factory subdirectory structure.
```

### Example Spawn Message:
```bash
cd /efforts/phase1/wave2/buildah-builder-interface
pwd  # Verify directory per R208

@agent-code-reviewer Please create the implementation plan for buildah-builder-interface.
Phase: 1, Wave: 2

**CRITICAL**: Save your plan at:
.software-factory/phase1/wave2/buildah-builder-interface/IMPLEMENTATION-PLAN-[TIMESTAMP].md

The orchestrator will look for your plan in this exact location.
```

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
