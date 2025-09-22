# Orchestrator - ANALYZE_CODE_REVIEWER_PARALLELIZATION State Rules

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

**YOU HAVE ENTERED ANALYZE_CODE_REVIEWER_PARALLELIZATION STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_ANALYZE_CODE_REVIEWER_PARALLELIZATION
echo "$(date +%s) - Rules read and acknowledged for ANALYZE_CODE_REVIEWER_PARALLELIZATION" > .state_rules_read_orchestrator_ANALYZE_CODE_REVIEWER_PARALLELIZATION
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY ANALYZE_CODE_REVIEWER_PARALLELIZATION WORK UNTIL RULES ARE READ:
- ❌ Start analyze effort dependencies
- ❌ Start determine parallelization strategy
- ❌ Start plan reviewer spawning
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**
### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR ANALYZE_CODE_REVIEWER_PARALLELIZATION STATE

### Core Mandatory Rules (ALL orchestrator states must have these):

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Criticality: BLOCKING - Automatic termination, 0% grade
   - Summary: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents

2. **🔴🔴🔴 R287** - TODO PERSISTENCE COMPREHENSIVE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Criticality: SUPREME - -20% to -100% penalty for violations
   - Summary: MUST save TODOs within 30s after write, every 10 messages, before transitions

3. **🔴🔴🔴 R288** - STATE FILE UPDATE REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-requirements.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state.json before EVERY state transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322 Part A-mandatory-stop-after-spawn.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

5. **🔴🔴🔴 R151** - Parallelization Requirements
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-execution-timestamp-compliance.md`
   - Criticality: SUPREME LAW - Parallel spawn requirements
   - Summary: Agents spawned in parallel must have timestamps within 5 seconds


### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R151, R208, R053..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all ANALYZE_CODE_REVIEWER_PARALLELIZATION rules"
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

### ✅ CORRECT PATTERN FOR ANALYZE_CODE_REVIEWER_PARALLELIZATION:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute ANALYZE_CODE_REVIEWER_PARALLELIZATION work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY ANALYZE_CODE_REVIEWER_PARALLELIZATION work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute ANALYZE_CODE_REVIEWER_PARALLELIZATION work


1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with ANALYZE_CODE_REVIEWER_PARALLELIZATION work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY ANALYZE_CODE_REVIEWER_PARALLELIZATION work before reading and acknowledging rules:**
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

## 📋 RULE SUMMARY FOR ANALYZE_CODE_REVIEWER_PARALLELIZATION STATE

### PRIMARY DIRECTIVES - MUST READ ALL:

### 🛑 RULE R322 - Mandatory Stop Before State Transitions
**Source:** $CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md
**Criticality:** 🔴🔴🔴 SUPREME LAW - Violation = -100% FAILURE

After completing state work and committing state file:
1. STOP IMMEDIATELY
2. Do NOT continue to next state
3. Do NOT start new work
4. Exit and wait for user
---

### 🔴🔴🔴 R234 - Mandatory State Traversal (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md`
**Criticality**: SUPREME LAW - Violation = -100% automatic failure
**Summary**: Must traverse all states in sequence, no skipping allowed

### 🚨🚨🚨 R218 - Orchestrator Parallel Code Reviewer Spawning
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R218-orchestrator-parallel-code-reviewer-spawning.md`
**Criticality**: BLOCKING - Cannot proceed without compliance
**Summary**: Mandatory parallelization analysis before spawning

### ⚠️⚠️⚠️ R151 - Parallel Agent Spawning Timing
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-spawning-timing.md`
**Criticality**: CRITICAL - <5s delta required
**Summary**: All parallel agents must acknowledge within 5 seconds

### 🔴🔴🔴 R208 - Orchestrator Working Directory Spawn Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R208-orchestrator-spawn-cd-protocol.md`
**Criticality**: SUPREME LAW - Spawn without CD = -100% failure
**Summary**: Must CD to working copy before spawning agents

### ⚠️⚠️⚠️ R219 - Code Reviewer Dependency-Aware Planning
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R219-code-reviewer-dependency-aware-effort-planning.md`
**Criticality**: CRITICAL - Consider dependencies in planning
**Summary**: Code reviewers must analyze effort dependencies

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Must update orchestrator-state.json on all transitions

### 🚨🚨🚨 R288 - State File Update and Commit Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: BLOCKING - Push immediately after update
**Summary**: Commit and push state within 60 seconds

### 🚨🚨🚨 R287 - Mandatory TODO Save Triggers
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R287-mandatory-todo-save-triggers.md`
**Criticality**: BLOCKING - Save within 30 seconds
**Summary**: Must save TODOs before state transitions

### Critical Requirements:
1. READ Wave Implementation Plan with Read tool - Penalty: -25%
2. Extract ALL parallelization metadata - Penalty: -30%
3. Create and save parallelization plan - Penalty: -20%
4. Output acknowledgment of decision - Penalty: -15%
5. Transition to SPAWN_CODE_REVIEWERS_EFFORT_PLANNING - Penalty: -100%

### Success Criteria:
- ✅ Wave plan READ with Read tool
- ✅ Blocking efforts identified
- ✅ Parallel groups created
- ✅ Plan saved to orchestrator-state.json
- ✅ Acknowledgment output displayed

### Failure Triggers:
- ❌ Skip this state = -100% R234 VIOLATION
- ❌ Not reading wave plan = R218 VIOLATION
- ❌ Skip to SPAWN_AGENTS = AUTOMATIC FAILURE
- ❌ No parallelization plan saved = Cannot proceed

## 🚨 ANALYZE_CODE_REVIEWER_PARALLELIZATION IS A VERB - START ANALYZING NOW! 🚨

### IMMEDIATE ACTIONS UPON ENTERING ANALYZE_CODE_REVIEWER_PARALLELIZATION

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. READ the Wave Implementation Plan with Read tool NOW
2. Extract "Can Parallelize" metadata immediately
3. Create parallelization groups without delay
4. Save the plan to orchestrator-state.json NOW
5. Output acknowledgment decision immediately

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in ANALYZE_CODE_REVIEWER_PARALLELIZATION" [stops]
- ❌ "Successfully entered parallelization analysis state" [waits]
- ❌ "Ready to analyze parallelization" [pauses]
- ❌ "I'm in analysis state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering ANALYZE_CODE_REVIEWER_PARALLELIZATION, reading wave plan now..."
- ✅ "Analyzing parallelization, extracting E3.1.1 metadata..."
- ✅ "ANALYZING: Found blocking effort E3.1.1, creating spawn sequence..."

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## State Context
You MUST analyze the wave implementation plan to determine parallelization strategy BEFORE spawning any Code Reviewers for effort planning. This is a MANDATORY GATE that prevents parallelization violations - DO IT NOW!

## 🔴🔴🔴 SUPREME LAW R234 - MANDATORY STATE SEQUENCE 🔴🔴🔴

**THIS STATE IS PART OF THE MANDATORY SEQUENCE - NO SKIPPING!**

See: `$CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md`

### YOUR POSITION IN THE MANDATORY SEQUENCE:
```
SETUP_EFFORT_INFRASTRUCTURE (✓ completed)
    ↓
ANALYZE_CODE_REVIEWER_PARALLELIZATION (👈 YOU ARE HERE)
    ↓ (MUST GO HERE NEXT)
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
    ↓
WAITING_FOR_EFFORT_PLANS
    ↓
ANALYZE_IMPLEMENTATION_PARALLELIZATION
    ↓
SPAWN_AGENTS
```

**CRITICAL:** You got here from SETUP_EFFORT_INFRASTRUCTURE (correct!)
**MANDATORY:** You MUST go to SPAWN_CODE_REVIEWERS_EFFORT_PLANNING next
**FORBIDDEN:** Skipping ahead to any other state = -100% FAILURE

## 🔴🔴🔴 ABSOLUTE REQUIREMENT 🔴🔴🔴

**THIS STATE IS A MANDATORY STOP!**
- You CANNOT proceed to SPAWN_CODE_REVIEWERS_EFFORT_PLANNING without completing this analysis
- You MUST create a parallelization plan and save it to orchestrator-state.json
- You MUST acknowledge your parallelization decision BEFORE any spawning


## Mandatory Analysis Protocol

```bash
# STEP 0: R356 SINGLE-EFFORT OPTIMIZATION CHECK
echo "═══════════════════════════════════════════════════════════════"
echo "🔴🔴🔴 ANALYZE_CODE_REVIEWER_PARALLELIZATION STATE 🔴🔴🔴"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "🎯 R356: Checking effort count for optimization..."

EFFORT_COUNT=$(jq '.efforts_in_progress | length' orchestrator-state.json)

if [ "$EFFORT_COUNT" -eq 1 ]; then
    EFFORT_NAME=$(jq -r '.efforts_in_progress[0].name' orchestrator-state.json)
    echo "═══════════════════════════════════════════════════════════════"
    echo "🎯 R356 OPTIMIZATION: Single Effort - No Parallelization Needed"
    echo "═══════════════════════════════════════════════════════════════"
    echo "Effort: $EFFORT_NAME"
    echo "Analysis Result: Cannot parallelize single effort"
    echo "Strategy: Spawn single Code Reviewer for the sole effort"
    echo "Next State: SPAWN_CODE_REVIEWERS_EFFORT_PLANNING (immediate)"
    echo "═══════════════════════════════════════════════════════════════"

    # Save minimal parallelization plan for single effort
    jq --arg timestamp "$(date -Iseconds)" \
       --arg effort "$EFFORT_NAME" \
       '.code_reviewer_parallelization_plan = {
          "wave": .current_wave,
          "phase": .current_phase,
          "analysis_timestamp": $timestamp,
          "single_effort_optimization": true,
          "r356_applied": true,
          "efforts": [$effort],
          "spawn_sequence": [{
            "step": 1,
            "action": "spawn_single",
            "efforts": [$effort],
            "r151_requirement": "Not applicable - single effort"
          }]
        }' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

    echo "✅ Single effort plan saved to orchestrator-state.json"
    echo "✅ Immediately transitioning to SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"

    # Immediate transition per R356
    safe_state_transition "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING" "Single effort - R356 optimization"
    exit 0  # Early exit - no further analysis needed
fi

echo "Multiple efforts detected: $EFFORT_COUNT"
echo "Proceeding with full parallelization analysis..."

# STEP 1: READ Wave Implementation Plan (MANDATORY - USE READ TOOL!)
echo ""
echo "📖 R218: MANDATORY - Reading Wave Implementation Plan for parallelization analysis"

WAVE_PLAN="phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-IMPLEMENTATION-PLAN.md"

# MANDATORY: Actually READ the file with Read tool
echo "🚨 USING READ TOOL to extract parallelization metadata from:"
echo "   $WAVE_PLAN"
# READ: $WAVE_PLAN

# STEP 2: Extract Parallelization Metadata
echo ""
echo "🔍 Analyzing effort parallelization metadata..."

# Parse the plan content to identify:
# - Efforts with "Can Parallelize: No" (blocking)
# - Efforts with "Can Parallelize: Yes" (parallel)
# - Dependencies between efforts

# STEP 3: Create Parallelization Groups
cat > /tmp/parallelization_analysis.yaml << EOF
code_reviewer_parallelization_plan:
  wave: ${WAVE}
  phase: ${PHASE}
  analysis_timestamp: "$(date -Iseconds)"
  
  blocking_efforts:
    # Efforts that MUST run sequentially
    - effort_id: "E3.1.1"
      name: "sync-engine-foundation"
      reason: "blocks all other efforts"
      can_parallelize: false
      
  parallel_groups:
    # Groups that can run in parallel AFTER blocking efforts
    group_1:
      can_start_after: ["E3.1.1"]
      efforts:
        - effort_id: "E3.1.2"
          name: "webhook-framework"
        - effort_id: "E3.1.3"
          name: "controller-setup"
        - effort_id: "E3.1.4"
          name: "validation-engine"
        - effort_id: "E3.1.5"
          name: "status-management"
          
  spawn_sequence:
    - step: 1
      action: "spawn_sequential"
      efforts: ["E3.1.1"]
      wait_for_completion: true
    - step: 2
      action: "spawn_parallel"
      efforts: ["E3.1.2", "E3.1.3", "E3.1.4", "E3.1.5"]
      r151_requirement: "All in ONE message with <5s delta"
EOF

# STEP 4: Update orchestrator-state.json
echo ""
echo "💾 Saving parallelization plan to orchestrator-state.json..."
# Update the state file with the parallelization plan
```

## Mandatory Acknowledgment Output

**YOU MUST OUTPUT ALL OF THE FOLLOWING BEFORE PROCEEDING:**

```
═══════════════════════════════════════════════════════════════
📋 PARALLELIZATION ANALYSIS COMPLETE
═══════════════════════════════════════════════════════════════

✅ I have READ the Wave Implementation Plan at:
   phase-plans/PHASE-3-WAVE-1-IMPLEMENTATION-PLAN.md

📊 PARALLELIZATION DECISION:
   
   BLOCKING EFFORTS (must complete first):
   - E3.1.1: sync-engine-foundation (blocks all others)
   
   PARALLEL EFFORTS (can spawn together after blocking):
   - E3.1.2: webhook-framework
   - E3.1.3: controller-setup  
   - E3.1.4: validation-engine
   - E3.1.5: status-management
   
🚨 SPAWN STRATEGY COMMITMENT:
   Step 1: Spawn E3.1.1 ALONE and WAIT for completion
   Step 2: Spawn E3.1.2, E3.1.3, E3.1.4, E3.1.5 TOGETHER in ONE message
   
✅ This strategy is SAVED in orchestrator-state.json
✅ I WILL follow this strategy in SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
═══════════════════════════════════════════════════════════════
```

## Post-State Actions

After completing parallelization analysis:
1. Save parallelization plan to orchestrator-state.json
2. Display state machine visualization (R230)
3. Transition to SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
4. Re-acknowledge critical rules (R217)
5. Use the saved plan to spawn correctly

## Validation Checks Before State Transition

```bash
validate_parallelization_analysis() {
    echo "🔍 Validating parallelization analysis completeness..."
    
    # CHECK 1: Was the wave plan actually read?
    if ! grep -q "code_reviewer_parallelization_plan" orchestrator-state.json; then
        echo "❌ FATAL: No parallelization plan found in orchestrator-state.json!"
        echo "Cannot proceed to SPAWN_CODE_REVIEWERS_EFFORT_PLANNING!"
        exit 1
    fi
    
    # CHECK 2: Are blocking efforts identified?
    BLOCKING_COUNT=$(jq '.code_reviewer_parallelization_plan.blocking_efforts | length' orchestrator-state.json)
    echo "✅ Found $BLOCKING_COUNT blocking efforts"
    
    # CHECK 3: Are parallel groups identified?
    PARALLEL_COUNT=$(jq '.code_reviewer_parallelization_plan.parallel_groups | length' orchestrator-state.json)
    echo "✅ Found $PARALLEL_COUNT parallel groups"
    
    # CHECK 4: Is spawn sequence defined?
    SEQUENCE_COUNT=$(jq '.code_reviewer_parallelization_plan.spawn_sequence | length' orchestrator-state.json)
    echo "✅ Spawn sequence has $SEQUENCE_COUNT steps"
    
    if [ $SEQUENCE_COUNT -eq 0 ]; then
        echo "❌ FATAL: No spawn sequence defined!"
        echo "Cannot proceed without a clear spawn strategy!"
        exit 1
    fi
    
    echo ""
    echo "✅ Parallelization analysis is COMPLETE and VALID"
    echo "✅ Ready to transition to SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
}
```

## State Transition Requirements

**BEFORE transitioning to SPAWN_CODE_REVIEWERS_EFFORT_PLANNING:**

1. ✅ Wave Implementation Plan has been READ with Read tool
2. ✅ All efforts have been analyzed for parallelization
3. ✅ Parallelization plan saved to orchestrator-state.json
4. ✅ Acknowledgment output has been displayed
5. ✅ Validation checks have passed
6. ✅ TODOs saved per R287 (before transition)

### R287-R287 TODO PERSISTENCE CHECKPOINT
```bash
# MANDATORY before state transition
echo "💾 R287: State transition trigger - saving TODOs..."
save_todos "ANALYZE_CODE_REVIEWER_PARALLELIZATION complete"

# R287: Commit within 60 seconds
cd $CLAUDE_PROJECT_DIR
git add todos/*.todo orchestrator-state.json
git commit -m "todo: parallelization analysis complete"
git push
echo "✅ Analysis and TODOs persisted"
```

**AFTER completing this state:**
```yaml
orchestrator_state:
  current_state: "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
  previous_state: "ANALYZE_CODE_REVIEWER_PARALLELIZATION"
  transition_time: "{ISO-8601}"
  transition_reason: "Parallelization analysis complete, strategy defined"
  
  code_reviewer_parallelization_plan:
    wave: 1
    phase: 3
    blocking_efforts: [...]
    parallel_groups: [...]
    spawn_sequence: [...]
```

## Critical Failure Conditions

**THE FOLLOWING WILL CAUSE IMMEDIATE STATE FAILURE:**

1. ❌ Attempting to spawn without reading the wave plan
2. ❌ Skipping this state and going directly to SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
3. ❌ Not saving parallelization plan to orchestrator-state.json
4. ❌ Spawning all efforts in parallel when some are blocking
5. ❌ Not acknowledging the parallelization decision

## Example: Correct Execution

```bash
# 1. Enter state
echo "Transitioning to ANALYZE_CODE_REVIEWER_PARALLELIZATION"

# 2. Read wave plan with Read tool
READ: phase-plans/PHASE-3-WAVE-1-IMPLEMENTATION-PLAN.md

# 3. Analyze parallelization
echo "Found 1 blocking effort: E3.1.1"
echo "Found 4 parallel efforts: E3.1.2, E3.1.3, E3.1.4, E3.1.5"

# 4. Save to orchestrator-state.json
jq '.code_reviewer_parallelization_plan = ...' orchestrator-state.json

# 5. Acknowledge decision
echo "✅ PARALLELIZATION DECISION COMMITTED"
echo "   Will spawn E3.1.1 first (blocking)"
echo "   Then spawn E3.1.2-E3.1.5 together (parallel)"

# 6. Validate and transition
validate_parallelization_analysis
echo "Transitioning to SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
```

## Integration with Other Rules

- **R218**: Mandatory reading of wave plan parallelization
- **R151**: Parallel spawn timing requirements (<5s delta)
- **R208**: Directory protocol for spawning
- **R219**: Dependency-aware effort planning
- **R288**: Mandatory state file updates (includes commit/push)

## Grading Impact

**This state affects orchestrator grading:**
- Skipping this state: **-50%** (Critical violation)
- Not reading wave plan: **-25%** (R218 violation)
- Wrong parallelization groups: **-30%** (Analysis failure)
- No acknowledgment: **-15%** (Protocol violation)
- Not saving to state file: **-20%** (R288 violation)

## Summary

This state is a **MANDATORY GATE** that:
1. **FORCES** analysis of parallelization metadata
2. **PREVENTS** spawning violations before they happen
3. **COMMITS** to a spawn strategy before execution
4. **SAVES** the decision for audit and recovery
5. **BLOCKS** progression without proper analysis

**YOU CANNOT SKIP THIS STATE!**

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**


### 🔴🔴🔴 MANDATORY VALIDATION REQUIREMENT 🔴🔴🔴

**Per R288 and R324**: ALL state file updates MUST be validated before commit:

```bash
# After ANY update to orchestrator-state.json:
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state.json || {
    echo "❌ State file validation failed!"
    exit 288
}
```

**Use helper functions for automatic validation:**
```bash
# Source the helper functions
source "$CLAUDE_PROJECT_DIR/utilities/state-file-update-functions.sh"

# Use safe functions that include validation:
safe_state_transition "NEW_STATE" "reason"
safe_update_field "field_name" "value"
```
