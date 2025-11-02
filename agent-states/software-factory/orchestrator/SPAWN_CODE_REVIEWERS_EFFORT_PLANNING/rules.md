# Orchestrator - SPAWN_CODE_REVIEWERS_EFFORT_PLANNING State Rules

# PRIMARY DIRECTIVES

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R251** - Initial Effort Planning Protocol (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R251-REPOSITORY-SEPARATION-LAW.md`

3. **R309** - Primary Implementation Effort Planning (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R309-never-create-efforts-in-sf-repo.md`

4. **R287** - TODO Persistence Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`

5. **R288** - State File Update Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`

6. **R304** - Mandatory Line Counter Usage (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`

7. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`

8. **R324** - State Transition Validation (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R324-mandatory-line-counter-auto-detection.md`


## 🛑🛑🛑 R322 MANDATORY STOP + R405 CONTINUATION FLAG 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### 🚨 CRITICAL DISTINCTION: STOP vs CONTINUATION FLAG 🚨

**R322 "STOP" = End conversation (`exit 0`) for context preservation**
**R405 FLAG = Can automation restart? (TRUE for spawning = YES!)**

**THESE ARE INDEPENDENT! You MUST do BOTH:**

### ✅ REQUIRED PATTERN FOR THIS SPAWN STATE:
```bash
# 1. Complete spawning work
echo "✅ Spawned Code Reviewers for effort planning"

# 2. Update state file per R324/R288
update_state "WAITING_FOR_EFFORT_PLANS"
commit_state_files_per_r288()

# 3. Save TODOs per R287
save_todos "SPAWNED_CODE_REVIEWERS"

# 4. R322: Stop conversation for context preservation
echo "🛑 R322: Stopping after spawn for context preservation"

# 5. R405: CONTINUATION FLAG (MUST BE TRUE FOR SPAWNING!)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# 6. Exit to end conversation
exit 0
```

**WHY TRUE IS CORRECT:**
- ✅ Spawning is NORMAL operation (not an error!)
- ✅ System knows next state (WAITING_FOR_EFFORT_PLANS)
- ✅ Automation can continue (designed workflow!)
- ✅ Context preservation ≠ manual intervention needed!

**STOP MEANS:**
- End THIS conversation turn (`exit 0`)
- Preserve context for next turn
- Allow state file commits

**TRUE MEANS:**
- Automation CAN restart this conversation
- System knows what to do next
- Normal operation proceeding

### ❌ WRONG PATTERN (CAUSES TEST FAILURES):
```bash
# This KILLS automation and causes Test 2 to stop!
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"  # ❌ WRONG for spawning!
exit 0
```

**FALSE should ONLY be used for:**
- ❌ System corruption (state files invalid)
- ❌ Critical infrastructure destroyed
- ❌ Unrecoverable errors requiring human intervention
- ❌ **NOT for normal spawning operations!**

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


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_CODE_REVIEWERS_EFFORT_PLANNING STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
mkdir -p markers/state-verification && touch "markers/state-verification/state_rules_read_orchestrator_SPAWN_CODE_REVIEWERS_EFFORT_PLANNING-$(date +%Y%m%d-%H%M%S)"
echo "$(date +%s) - Rules read and acknowledged for SPAWN_CODE_REVIEWERS_EFFORT_PLANNING" > "markers/state-verification/state_rules_read_orchestrator_SPAWN_CODE_REVIEWERS_EFFORT_PLANNING-$(date +%Y%m%d-%H%M%S)"
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

### ✅ CORRECT PATTERN FOR SPAWN_CODE_REVIEWERS_EFFORT_PLANNING:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SPAWN_CODE_REVIEWERS_EFFORT_PLANNING work until:**
### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SPAWN_CODE_REVIEWERS_EFFORT_PLANNING work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING_SWE_PROGRESS YOUR READ TOOL CALLS!**

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
- ❌ Skip to SPAWN_SW_ENGINEERS = -100% R234 VIOLATION
- ❌ Ignore parallelization plan = -50% penalty
- ❌ Spawn sequentially when should be parallel = -50%

## 🚨 SPAWN_CODE_REVIEWERS_EFFORT_PLANNING IS A VERB - START SPAWNING CODE REVIEWERS IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING SPAWN_CODE_REVIEWERS_EFFORT_PLANNING

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Spawn Code Reviewer for first effort NOW per parallelization plan
2. Use spawn sequence from orchestrator-state-v3.json immediately
3. Check TodoWrite for pending items and process them
4. Follow R151 timing requirements for parallel spawns

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in SPAWN_CODE_REVIEWERS_EFFORT_PLANNING" [stops]
- ❌ "Successfully entered SPAWN_CODE_REVIEWERS_EFFORT_PLANNING state" [waits]
- ❌ "Ready to start spawning code reviewers" [pauses]
- ❌ "I'm in SPAWN_CODE_REVIEWERS_EFFORT_PLANNING state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering SPAWN_CODE_REVIEWERS_EFFORT_PLANNING, Spawn Code Reviewer for first effort NOW per parallelization plan..."
- ✅ "START SPAWNING CODE REVIEWERS, use spawn sequence from orchestrator-state-v3.json immediately..."
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
CREATE_NEXT_INFRASTRUCTURE (✓ completed)
    ↓
ANALYZE_CODE_REVIEWER_PARALLELIZATION (✓ completed)
    ↓
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING (👈 YOU ARE HERE)
    ↓ (MUST GO HERE NEXT)
WAITING_FOR_EFFORT_PLANS
    ↓
ANALYZE_IMPLEMENTATION_PARALLELIZATION
    ↓
SPAWN_SW_ENGINEERS
```

**CRITICAL:** You analyzed parallelization (correct!)
**NOW:** Spawn Code Reviewers using that analysis
**NEXT:** You MUST go to WAITING_FOR_EFFORT_PLANS
**FORBIDDEN:** Skipping to SPAWN_SW_ENGINEERS = -100% FAILURE

## 🔴🔴🔴 PREREQUISITES 🔴🔴🔴

**BEFORE ENTERING THIS STATE, YOU MUST HAVE:**
1. ✅ Wave Implementation Plan created
2. ✅ ALL effort directories created (CREATE_NEXT_INFRASTRUCTURE)
3. ✅ ALL branches pushed to remote with tracking
4. ✅ work-log.md files initialized in each directory
5. ✅ **PARALLELIZATION ANALYSIS COMPLETE (ANALYZE_CODE_REVIEWER_PARALLELIZATION)**
6. ✅ **Parallelization plan in orchestrator-state-v3.json**

**IF PARALLELIZATION NOT ANALYZED, GO BACK TO ANALYZE_CODE_REVIEWER_PARALLELIZATION!**
**IF INFRASTRUCTURE IS NOT READY, GO BACK TO CREATE_NEXT_INFRASTRUCTURE!**

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

## 🚨🚨🚨 CRITICAL BUG FIX: SEQUENTIAL EFFORT SPAWNING (R208 ENFORCEMENT) 🚨🚨🚨

### SUPREME LAW VIOLATION PREVENTION

**DOCUMENTED BUG**: When spawning Code Reviewers for SEQUENTIAL efforts (effort 2, effort 3, etc.), the orchestrator has historically FAILED to CD to the correct directory for the SECOND+ effort, causing the Code Reviewer to create the plan in the WRONG location.

**ROOT CAUSE**: After spawning Code Reviewer for effort 1, the orchestrator remains in effort 1's directory. When spawning for effort 2, the orchestrator FORGETS to CD to effort 2's directory, violating R208.

### MANDATORY PATTERN FOR SEQUENTIAL EFFORTS:

```bash
# CORRECT PATTERN FOR SEQUENTIAL EFFORT SPAWNING

# Spawn effort 1
echo "=== SPAWNING CODE REVIEWER FOR EFFORT 1 ==="
cd efforts/phase${PHASE}/wave${WAVE}/${EFFORT_1_DIR}
pwd  # MUST output: efforts/phase${PHASE}/wave${WAVE}/${EFFORT_1_DIR}
@agent-code-reviewer [spawn instructions for effort 1]

# CRITICAL: Return to orchestrator directory BEFORE spawning effort 2
cd $CLAUDE_PROJECT_DIR
pwd  # MUST output: /workspaces/[project-name]

# Spawn effort 2 (DO NOT SKIP THIS CD STEP!)
echo "=== SPAWNING CODE REVIEWER FOR EFFORT 2 ==="
cd efforts/phase${PHASE}/wave${WAVE}/${EFFORT_2_DIR}  # ← CRITICAL! DO NOT SKIP!
pwd  # MUST output: efforts/phase${PHASE}/wave${WAVE}/${EFFORT_2_DIR}
@agent-code-reviewer [spawn instructions for effort 2]

# Return to orchestrator directory
cd $CLAUDE_PROJECT_DIR
pwd  # MUST output: /workspaces/[project-name]
```

### VALIDATION CHECKLIST (MANDATORY):

Before spawning EACH Code Reviewer in sequential mode, verify:

```bash
# FOR EACH EFFORT IN SEQUENTIAL LIST:
for effort_dir in "${EFFORT_DIRECTORIES[@]}"; do
    echo "🔍 R208 PRE-SPAWN VALIDATION FOR: $effort_dir"

    # 1. Return to orchestrator directory first
    cd "$CLAUDE_PROJECT_DIR"
    echo "  ✓ Returned to orchestrator: $(pwd)"

    # 2. CD to effort directory
    cd "efforts/phase${PHASE}/wave${WAVE}/${effort_dir}" || {
        echo "❌ CRITICAL: Failed to CD to $effort_dir"
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=R208_DIRECTORY_FAILURE"
        exit 208
    }

    # 3. Verify we're in correct directory
    CURRENT_DIR=$(basename "$(pwd)")
    if [[ "$CURRENT_DIR" != "$effort_dir" ]]; then
        echo "❌ CRITICAL R208 VIOLATION: Directory mismatch!"
        echo "   Expected: $effort_dir"
        echo "   Actual:   $CURRENT_DIR"
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=R208_VIOLATION"
        exit 208
    fi
    echo "  ✓ Directory validated: $(pwd)"

    # 4. Output PWD before spawn (MANDATORY)
    pwd

    # 5. Spawn Code Reviewer
    echo "  🚀 Spawning Code Reviewer in $(pwd)..."
    @agent-code-reviewer [spawn instructions]

    # 6. Return to orchestrator directory
    cd "$CLAUDE_PROJECT_DIR"
    echo "  ✓ Returned to orchestrator: $(pwd)"
done
```

### HISTORICAL BUG EXAMPLE (DO NOT REPRODUCE):

```bash
# ❌ WRONG - This caused bug 2.2.2 implementation plan failure!

cd efforts/phase2/wave2/effort-1-registry-override-viper
pwd  # Output: efforts/phase2/wave2/effort-1-registry-override-viper
@agent-code-reviewer [spawn for effort 1]

# BUG: Forgot to CD back to orchestrator, then CD to effort 2!
# Still in effort-1 directory when spawning effort 2 reviewer!
@agent-code-reviewer [spawn for effort 2]  # ← WRONG DIRECTORY!

# Result: Code Reviewer for effort 2 uses basename(pwd) = effort-1
#         Plan created in .software-factory/phase2/wave2/effort-1/...
#         Instead of .software-factory/phase2/wave2/effort-2/...
```

### ENFORCEMENT:

**EVERY sequential effort spawn MUST:**
1. Start from orchestrator directory (`cd $CLAUDE_PROJECT_DIR`)
2. CD to specific effort directory (`cd efforts/phase${PHASE}/wave${WAVE}/${EFFORT_DIR}`)
3. Verify directory with `pwd` output (visible in conversation)
4. Spawn Code Reviewer (they will use `basename $(pwd)` for effort name)
5. Return to orchestrator directory (`cd $CLAUDE_PROJECT_DIR`)
6. Repeat for next effort

**FAILURE TO FOLLOW THIS PATTERN = -100% R208 VIOLATION**

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



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete SPAWN_CODE_REVIEWERS_EFFORT_PLANNING:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="SPAWN_CODE_REVIEWERS_EFFORT_PLANNING complete - [describe what was accomplished]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Update State File (R288 - SUPREME LAW)
```bash
# State Manager validates transition and updates state files (SF 3.0 Pattern)
echo "🔄 Spawning State Manager for SHUTDOWN_CONSULTATION..."

# Prepare work results summary
WORK_RESULTS=$(cat <<EOF
{
  "state_completed": "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING",
  "work_accomplished": [
    "[List accomplishments from state work]"
  ],
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON" \
  --work-results "$WORK_RESULTS"

# State Manager will:
# 1. Validate PROPOSED_NEXT_STATE exists and transition is valid
# 2. Update all 4 state files atomically (R288)
# 3. Commit and push state files
# 4. Return REQUIRED_NEXT_STATE (usually same as proposed unless invalid)

echo "✅ State Manager consultation complete"
echo "✅ State files updated by State Manager"
```

---

### ✅ Step 4: Validate State File (R324)
```bash
# Validate state file before committing
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
echo "✅ State file validated"
```

---

### ✅ Step 5: Commit State File (R288)
```bash
# Commit and push state file immediately
git add orchestrator-state-v3.json

if ! git commit -m "state: SPAWN_CODE_REVIEWERS_EFFORT_PLANNING → $NEXT_STATE - SPAWN_CODE_REVIEWERS_EFFORT_PLANNING complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
    echo "Attempted transition from: SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
    echo ""
    echo "Common causes:"
    echo "  - Schema validation failure (check pre-commit hook output above)"
    echo "  - Missing required fields in JSON files"
    echo "  - Invalid JSON syntax"
    echo ""
    echo "🛑 Cannot proceed - manual intervention required"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=SCHEMA_VALIDATION"
    exit 1
fi

git push || echo "⚠️ WARNING: Push failed - committed locally"
echo "✅ State file committed and pushed"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - SPAWN_CODE_REVIEWERS_EFFORT_PLANNING complete [R287]"; then
    echo "❌ ERROR: Failed to commit TODO files"
    echo "This is non-fatal but TODOs may be lost in compaction"
    echo "Proceeding with state execution..."
    # Don't exit - TODO commit failure is not fatal
fi

git push || echo "⚠️ WARNING: TODO push failed - committed locally"
echo "✅ TODOs saved and committed"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 7: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 8: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No state update = state machine broken (R288 violation, -100%)
- Missing Step 4: Invalid state = corruption (R324 violation)
- Missing Step 5: No commit = state lost on compaction (R288 violation, -100%)
- Missing Step 6: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 7: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 8: No exit = R322 violation (-100%)

**ALL 8 STEPS ARE MANDATORY - NO EXCEPTIONS**

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS (SF 3.0)

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="WAITING_FOR_EFFORT_PLANS"
TRANSITION_REASON="State work complete"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"
# State Manager updates all 4 state files atomically

# 4. Save TODOs
save_todos "R322_CHECKPOINT"

# 5. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# 6. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
### 🚨 SPAWN STATE PATTERN - R322 + R405 USAGE 🚨

**Spawning operations require R322 stop for context preservation:**
```bash
# After spawning agent(s)
echo "✅ Spawned agents for work"

# R322 checkpoint (context preservation)
echo "🛑 R322: Stopping after spawn for context preservation"

# Flag? → MUST BE TRUE (normal operation!)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# Stop inference
exit 0
```

**Why TRUE is correct:**
- Spawning is NORMAL operation
- System knows next state
- Automation can continue
- **Context preservation ≠ manual intervention needed!**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

