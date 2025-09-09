# Orchestrator - ANALYZE_IMPLEMENTATION_PARALLELIZATION State Rules

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

**YOU HAVE ENTERED ANALYZE_IMPLEMENTATION_PARALLELIZATION STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_ANALYZE_IMPLEMENTATION_PARALLELIZATION
echo "$(date +%s) - Rules read and acknowledged for ANALYZE_IMPLEMENTATION_PARALLELIZATION" > .state_rules_read_orchestrator_ANALYZE_IMPLEMENTATION_PARALLELIZATION
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY ANALYZE_IMPLEMENTATION_PARALLELIZATION WORK UNTIL RULES ARE READ:
- ❌ Start analyze implementation dependencies
- ❌ Start determine SWE parallelization
- ❌ Start plan agent allocation
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

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
   ❌ WRONG: "I acknowledge all ANALYZE_IMPLEMENTATION_PARALLELIZATION rules"
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

### ✅ CORRECT PATTERN FOR ANALYZE_IMPLEMENTATION_PARALLELIZATION:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute ANALYZE_IMPLEMENTATION_PARALLELIZATION work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY ANALYZE_IMPLEMENTATION_PARALLELIZATION work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute ANALYZE_IMPLEMENTATION_PARALLELIZATION work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with ANALYZE_IMPLEMENTATION_PARALLELIZATION work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY ANALYZE_IMPLEMENTATION_PARALLELIZATION work before reading and acknowledging rules:**
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

## 📋 PRIMARY DIRECTIVES FOR ANALYZE_IMPLEMENTATION_PARALLELIZATION

### 🚨🚨🚨 R213 - Wave and Effort Metadata Injection
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R213-wave-and-effort-metadata-protocol.md`
**Criticality**: BLOCKING - Must inject metadata before spawning SW Engineers
**Summary**: Inject metadata into IMPLEMENTATION-PLAN.md before spawn

### 🚨🚨🚨 R219 - Code Reviewer Dependency-Aware Planning
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R219-code-reviewer-dependency-aware-effort-planning.md`
**Criticality**: BLOCKING - Must analyze effort dependencies
**Summary**: Analyze dependencies for parallelization decisions

### ⚠️⚠️⚠️ R151 - Parallel Agent Spawning Timing
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-spawning-timing.md`
**Criticality**: CRITICAL - <5s delta required
**Summary**: All parallel agents must acknowledge within 5 seconds

### 🔴🔴🔴 R234 - Mandatory State Traversal (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md`
**Criticality**: SUPREME LAW - Violation = -100% automatic failure
**Summary**: Must traverse all states in sequence, no skipping allowed

## 🚨 ANALYZE_IMPLEMENTATION_PARALLELIZATION IS A VERB - START ANALYZING IMPLEMENTATION PARALLELIZATION IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING ANALYZE_IMPLEMENTATION_PARALLELIZATION

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Check agent availability in orchestrator-state.json NOW
2. Parse effort implementation plans for dependencies immediately
3. Check TodoWrite for pending items and process them
4. Create implementation spawn sequence without delay

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in ANALYZE_IMPLEMENTATION_PARALLELIZATION" [stops]
- ❌ "Successfully entered ANALYZE_IMPLEMENTATION_PARALLELIZATION state" [waits]
- ❌ "Ready to start analyzing implementation parallelization" [pauses]
- ❌ "I'm in ANALYZE_IMPLEMENTATION_PARALLELIZATION state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering ANALYZE_IMPLEMENTATION_PARALLELIZATION, Check agent availability in orchestrator-state.json NOW..."
- ✅ "START ANALYZING IMPLEMENTATION PARALLELIZATION, parse effort implementation plans for dependencies immediately..."
- ✅ "ANALYZE_IMPLEMENTATION_PARALLELIZATION: Create implementation spawn sequence without delay..."

## State Context
You MUST analyze the individual effort implementation plans to determine parallelization strategy BEFORE spawning any SW Engineers for implementation. This is a MANDATORY GATE after Code Reviewers create effort plans.

## 🔴🔴🔴 ABSOLUTE REQUIREMENT 🔴🔴🔴

**THIS STATE IS A MANDATORY STOP!**
- You CANNOT proceed to SPAWN_AGENTS without completing this analysis
- You MUST read EACH effort's IMPLEMENTATION-PLAN.md
- You MUST create a SW Engineer parallelization plan
- You MUST acknowledge your parallelization decision BEFORE any spawning


## Mandatory Analysis Protocol

```bash
# STEP 0: INJECT EFFORT METADATA (R213 - BLOCKING)
echo "═══════════════════════════════════════════════════════════════"
echo "🔧 R213: MANDATORY - Injecting Effort Metadata BEFORE Analysis"
echo "═══════════════════════════════════════════════════════════════"
# For each effort, inject metadata into IMPLEMENTATION-PLAN.md
# This MUST happen BEFORE spawning SW Engineers!
inject_effort_metadata "$PHASE" "$WAVE" "$EFFORT_NAME" "$EFFORT_NUM" "$WORKING_DIR" "$BRANCH"

# STEP 1: READ All Effort Implementation Plans
echo "═══════════════════════════════════════════════════════════════"
echo "🔴🔴🔴 ANALYZE_IMPLEMENTATION_PARALLELIZATION STATE 🔴🔴🔴"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "📖 R219: MANDATORY - Reading ALL Effort Implementation Plans"

# For each effort directory
for effort in efforts/phase${PHASE}/wave${WAVE}/*/; do
    EFFORT_NAME=$(basename "$effort")
    IMPL_PLAN="${effort}IMPLEMENTATION-PLAN.md"
    
    echo ""
    echo "🚨 USING READ TOOL on effort plan:"
    echo "   $IMPL_PLAN"
    # READ: $IMPL_PLAN
    
    # Extract parallelization context from each plan
    echo "🔍 Extracting parallelization metadata for $EFFORT_NAME..."
done

# STEP 2: Verify Consistency with Wave Plan
echo ""
echo "🔍 Verifying parallelization consistency..."
echo "   Comparing effort plans with original wave plan metadata"

# STEP 3: Create SW Engineer Spawn Groups
cat > /tmp/sw_engineer_parallelization.yaml << EOF
sw_engineer_parallelization_plan:
  wave: ${WAVE}
  phase: ${PHASE}
  analysis_timestamp: "$(date -Iseconds)"
  
  blocking_implementations:
    # Implementations that MUST run sequentially
    - effort_id: "E3.1.1"
      name: "sync-engine-foundation"
      implementation_plan: "efforts/phase3/wave1/sync-engine-foundation/IMPLEMENTATION-PLAN.md"
      can_parallelize: false
      reason: "Core foundation - blocks all others"
      dependencies: []
      
  parallel_groups:
    # Groups that can run in parallel AFTER blocking implementations
    group_1:
      can_start_after: ["E3.1.1"]
      efforts:
        - effort_id: "E3.1.2"
          name: "webhook-framework"
          implementation_plan: "efforts/phase3/wave1/webhook-framework/IMPLEMENTATION-PLAN.md"
        - effort_id: "E3.1.3"
          name: "controller-setup"
          implementation_plan: "efforts/phase3/wave1/controller-setup/IMPLEMENTATION-PLAN.md"
        - effort_id: "E3.1.4"
          name: "validation-engine"
          implementation_plan: "efforts/phase3/wave1/validation-engine/IMPLEMENTATION-PLAN.md"
        - effort_id: "E3.1.5"
          name: "status-management"
          implementation_plan: "efforts/phase3/wave1/status-management/IMPLEMENTATION-PLAN.md"
          
  spawn_sequence:
    - step: 1
      action: "spawn_sequential"
      agent_type: "sw-engineer"
      efforts: ["E3.1.1"]
      wait_for_completion: true
      expected_duration: "2-3 hours"
    - step: 2
      action: "spawn_parallel"
      agent_type: "sw-engineer"
      efforts: ["E3.1.2", "E3.1.3", "E3.1.4", "E3.1.5"]
      r151_requirement: "All in ONE message with <5s delta"
      expected_duration: "3-4 hours"
EOF

# STEP 4: Update orchestrator-state.json
echo ""
echo "💾 Saving SW Engineer parallelization plan to orchestrator-state.json..."
```

## Mandatory Acknowledgment Output

**YOU MUST OUTPUT ALL OF THE FOLLOWING BEFORE PROCEEDING:**

```
═══════════════════════════════════════════════════════════════
📋 IMPLEMENTATION PARALLELIZATION ANALYSIS COMPLETE
═══════════════════════════════════════════════════════════════

✅ I have READ ALL effort implementation plans:
   - efforts/phase3/wave1/sync-engine-foundation/IMPLEMENTATION-PLAN.md
   - efforts/phase3/wave1/webhook-framework/IMPLEMENTATION-PLAN.md
   - efforts/phase3/wave1/controller-setup/IMPLEMENTATION-PLAN.md
   - efforts/phase3/wave1/validation-engine/IMPLEMENTATION-PLAN.md
   - efforts/phase3/wave1/status-management/IMPLEMENTATION-PLAN.md

📊 SW ENGINEER SPAWN DECISION:
   
   BLOCKING IMPLEMENTATIONS (must complete first):
   - E3.1.1: sync-engine-foundation (foundation layer)
   
   PARALLEL IMPLEMENTATIONS (can spawn together after blocking):
   - E3.1.2: webhook-framework
   - E3.1.3: controller-setup
   - E3.1.4: validation-engine
   - E3.1.5: status-management
   
🚨 SPAWN STRATEGY COMMITMENT:
   Step 1: Spawn SW Engineer for E3.1.1 ALONE and WAIT
   Step 2: Spawn SW Engineers for E3.1.2-E3.1.5 TOGETHER (R151)
   
✅ Consistency verified with wave plan parallelization
✅ This strategy is SAVED in orchestrator-state.json
✅ I WILL follow this strategy in SPAWN_AGENTS state
═══════════════════════════════════════════════════════════════
```

## Post-State Actions

After completing implementation parallelization analysis:
1. Save implementation strategy to orchestrator-state.json
2. Display state machine visualization (R230)
3. **🔴🔴🔴 R324/R325 CRITICAL: Update current_state BEFORE stopping! 🔴🔴🔴**
4. Re-acknowledge critical rules (R217)
5. Use the saved strategy to spawn SW Engineers correctly

### 🚨🚨🚨 MANDATORY STATE TRANSITION PROTOCOL (R324/R325) 🚨🚨🚨

**YOU MUST UPDATE current_state TO "SPAWN_AGENTS" BEFORE STOPPING!**

```bash
# 🔴🔴🔴 COPY THIS EXACTLY - THIS PREVENTS INFINITE LOOPS! 🔴🔴🔴

echo "✅ Implementation parallelization analysis complete"

# CRITICAL: Update state file FIRST (R324 requirement)
echo "🔴 R324: Updating current_state to prevent infinite loop..."
jq '.current_state = "SPAWN_AGENTS"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.previous_state = "ANALYZE_IMPLEMENTATION_PARALLELIZATION"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.transition_time = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

# Verify the update
echo "✅ State updated to:"
grep "current_state:" orchestrator-state.json

# Commit and push IMMEDIATELY
git add orchestrator-state.json
git commit -m "state: transition to SPAWN_AGENTS from ANALYZE_IMPLEMENTATION_PARALLELIZATION (R324)"
git push

# NOW stop per R322
echo "🛑 STATE TRANSITION CHECKPOINT: ANALYZE_IMPLEMENTATION_PARALLELIZATION → SPAWN_AGENTS"
echo "📊 State file updated to: SPAWN_AGENTS ✅"
echo "⏸️ STOPPED - Ready to continue in SPAWN_AGENTS"
echo "When restarted, will execute SPAWN_AGENTS state"
# EXIT HERE
```

**⚠️ FAILURE TO UPDATE current_state = INFINITE LOOP BUG! ⚠️**

## Validation Checks Before State Transition

```bash
validate_implementation_parallelization() {
    echo "🔍 Validating implementation parallelization analysis..."
    
    # CHECK 1: Were all implementation plans read?
    EFFORT_COUNT=$(ls efforts/phase${PHASE}/wave${WAVE}/*/IMPLEMENTATION-PLAN.md 2>/dev/null | wc -l)
    echo "Found $EFFORT_COUNT effort implementation plans"
    
    if [ $EFFORT_COUNT -eq 0 ]; then
        echo "❌ FATAL: No implementation plans found!"
        echo "Code Reviewers must create plans first!"
        exit 1
    fi
    
    # CHECK 2: Is SW Engineer parallelization plan saved?
    if ! grep -q "sw_engineer_parallelization_plan" orchestrator-state.json; then
        echo "❌ FATAL: No SW Engineer parallelization plan in orchestrator-state.json!"
        echo "Cannot proceed to SPAWN_AGENTS!"
        exit 1
    fi
    
    # CHECK 3: Verify consistency with Code Reviewer parallelization
    CR_BLOCKING=$(jq '.code_reviewer_parallelization_plan.blocking_efforts | length' orchestrator-state.json)
    SW_BLOCKING=$(jq '.sw_engineer_parallelization_plan.blocking_implementations | length' orchestrator-state.json)
    
    if [ "$CR_BLOCKING" != "$SW_BLOCKING" ]; then
        echo "⚠️ WARNING: Mismatch between Code Reviewer and SW Engineer blocking counts"
        echo "   Code Reviewer blocking: $CR_BLOCKING"
        echo "   SW Engineer blocking: $SW_BLOCKING"
    fi
    
    # CHECK 4: Verify spawn sequence exists
    SEQUENCE_COUNT=$(jq '.sw_engineer_parallelization_plan.spawn_sequence | length' orchestrator-state.json)
    echo "✅ SW Engineer spawn sequence has $SEQUENCE_COUNT steps"
    
    echo ""
    echo "✅ Implementation parallelization analysis is COMPLETE"
    echo "✅ Ready to transition to SPAWN_AGENTS"
}
```

## Consistency Verification Protocol

```bash
verify_parallelization_consistency() {
    echo "🔍 Verifying parallelization consistency across plans..."
    
    # Compare wave plan metadata with effort plans
    for effort in efforts/phase${PHASE}/wave${WAVE}/*/; do
        EFFORT_NAME=$(basename "$effort")
        IMPL_PLAN="${effort}IMPLEMENTATION-PLAN.md"
        
        # Check if effort plan preserves parallelization metadata
        if grep -q "Can Parallelize:" "$IMPL_PLAN"; then
            echo "✅ $EFFORT_NAME preserves parallelization metadata"
        else
            echo "⚠️ WARNING: $EFFORT_NAME missing parallelization metadata"
        fi
    done
    
    echo ""
    echo "✅ Consistency verification complete"
}
```

## State Transition Requirements

**BEFORE transitioning to SPAWN_AGENTS:**

1. ✅ ALL effort IMPLEMENTATION-PLAN.md files have been READ
2. ✅ Parallelization metadata extracted from each plan
3. ✅ Consistency verified with wave plan
4. ✅ SW Engineer parallelization plan saved to orchestrator-state.json
5. ✅ Acknowledgment output has been displayed
6. ✅ Validation checks have passed

**AFTER completing this state:**
```yaml
orchestrator_state:
  current_state: "SPAWN_AGENTS"
  previous_state: "ANALYZE_IMPLEMENTATION_PARALLELIZATION"
  transition_time: "{ISO-8601}"
  transition_reason: "SW Engineer parallelization analysis complete"
  
  sw_engineer_parallelization_plan:
    wave: 1
    phase: 3
    blocking_implementations: [...]
    parallel_groups: [...]
    spawn_sequence: [...]
```

## Critical Failure Conditions

**THE FOLLOWING WILL CAUSE IMMEDIATE STATE FAILURE:**

1. ❌ Not reading ANY implementation plans
2. ❌ Skipping this state and going directly to SPAWN_AGENTS
3. ❌ Not saving SW Engineer parallelization plan
4. ❌ Inconsistent parallelization with wave plan
5. ❌ Spawning all SW Engineers together when some are blocking

## Integration with Directory Protocol (R208)

```bash
prepare_spawn_with_directories() {
    echo "🗂️ R208: Preparing directory verification for spawn..."
    
    # For each effort in spawn sequence
    for effort in $EFFORTS; do
        EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${effort}"
        
        echo "   $effort will spawn in: $EFFORT_DIR"
        echo "   Branch: $(cd $EFFORT_DIR && git branch --show-current)"
        echo "   Implementation plan: $EFFORT_DIR/IMPLEMENTATION-PLAN.md"
    done
    
    echo "✅ All directories verified and ready for spawn"
}
```

## Example: Correct Execution

```bash
# 1. Enter state
echo "Transitioning to ANALYZE_IMPLEMENTATION_PARALLELIZATION"

# 2. Read ALL implementation plans
for plan in efforts/phase3/wave1/*/IMPLEMENTATION-PLAN.md; do
    READ: $plan
done

# 3. Analyze parallelization
echo "Analysis complete:"
echo "  1 blocking implementation (E3.1.1)"
echo "  4 parallel implementations (E3.1.2-E3.1.5)"

# 4. Verify consistency
verify_parallelization_consistency

# 5. Save to orchestrator-state.json
jq '.sw_engineer_parallelization_plan = ...' orchestrator-state.json

# 6. Acknowledge decision
echo "✅ SW ENGINEER SPAWN STRATEGY COMMITTED"

# 7. Validate and transition with R324 protocol
validate_implementation_parallelization

# 8. R324 CRITICAL: Update state BEFORE stopping
echo "🔴 R324: Updating current_state to SPAWN_AGENTS..."
jq '.current_state = "SPAWN_AGENTS"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.previous_state = "ANALYZE_IMPLEMENTATION_PARALLELIZATION"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
git add orchestrator-state.json
git commit -m "state: transition to SPAWN_AGENTS (R324)"
git push

# 9. NOW stop per R322
echo "🛑 Stopping - state updated to SPAWN_AGENTS"
```

## Integration with Other Rules

- **R151**: Parallel spawn timing requirements
- **R208**: Directory protocol for each spawn
- **R219**: Dependency-aware planning
- **R197**: One agent per effort requirement
- **R288**: Mandatory state file updates (includes commit/push)

## Grading Impact

**This state affects orchestrator grading:**
- Skipping this state: **-100%** (R234 SUPREME LAW VIOLATION)
- Not reading implementation plans: **-30%** (Analysis failure)
- Wrong parallelization groups: **-25%** (Strategy error)
- Inconsistent with wave plan: **-20%** (Coordination failure)
- No acknowledgment: **-15%** (Protocol violation)
- Forgetting TODO saves: **-15%** per violation (R287-R287)

### R287-R287 CHECKPOINT
```bash
# After analysis complete
echo "💾 R287: Saving TODOs after implementation analysis..."
save_todos "ANALYZE_IMPLEMENTATION_PARALLELIZATION complete"

# R287: Commit within 60 seconds
cd $CLAUDE_PROJECT_DIR
git add todos/*.todo orchestrator-state.json
git commit -m "todo: implementation parallelization analyzed"
git push

# Ready for SPAWN_AGENTS
echo "✅ Analysis complete - ready to spawn SW Engineers!"
```

## Summary

This state is a **MANDATORY GATE** that:
1. **ANALYZES** all effort implementation plans
2. **VERIFIES** consistency with wave parallelization
3. **CREATES** SW Engineer spawn strategy
4. **PREVENTS** parallelization violations
5. **ENSURES** R151 compliance before spawning
6. **ENFORCES** R234 mandatory state sequence (see rule library)
7. **REQUIRES** TODO persistence per R287-R287 (see rule library)

**THIS STATE CANNOT BE SKIPPED! R234 VIOLATION = -100% FAILURE!**

### Additional Rules Referenced:
- **R287**: `$CLAUDE_PROJECT_DIR/rule-library/R287-comprehensive-todo-persistence.md`
- **R288**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
- **R208**: `$CLAUDE_PROJECT_DIR/rule-library/R208-orchestrator-spawn-cd-protocol.md`
- **R197**: `$CLAUDE_PROJECT_DIR/rule-library/R197-one-agent-per-effort.md`

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
