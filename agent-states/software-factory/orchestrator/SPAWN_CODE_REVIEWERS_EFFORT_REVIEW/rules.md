# Orchestrator - SPAWN_CODE_REVIEWERS_EFFORT_REVIEW State Rules

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


# 🔴🔴🔴 MANDATORY: R322 STOP + R405 CONTINUATION FLAG 🔴🔴🔴

**CRITICAL FOR SPAWN STATES - READ THIS FIRST OR FAIL TEST 2!**

## 🚨 THE PATTERN THAT FAILED TEST 2 🚨

**WHAT HAPPENED IN TEST 2:**
- Orchestrator spawned Code Reviewers ✅ (correct)
- Orchestrator stopped per R322 ✅ (correct)
- Orchestrator **DID NOT emit `CONTINUE-SOFTWARE-FACTORY=TRUE`** ❌ (WRONG!)
- Test framework saw no continuation flag → stopped automation
- Test 2 FAILED at iteration 8

**ROOT CAUSE:** Confusion between R322 "stop" and R405 continuation flag

## 🔴 CRITICAL DISTINCTION: TWO INDEPENDENT DECISIONS 🔴

### Decision 1: Should Agent Stop? (R322 - Context Preservation)
**YES - ALWAYS stop after spawning for context preservation**

- **Purpose**: Prevent context overflow between states
- **Action**: `exit 0` to end conversation
- **User Experience**: User sees "/continue-orchestrating" as next step
- **This is NORMAL!** Not an error!

### Decision 2: Should Factory Continue? (R405 - Automation Control)
**YES - ALWAYS emit TRUE for normal spawning operations**

- **Purpose**: Tell automation whether it CAN restart
- **Action**: `echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"` (LAST output before exit)
- **Automation**: Framework will auto-restart orchestrator
- **This is NORMAL!** Designed behavior!

## ✅ REQUIRED PATTERN FOR ALL SPAWN STATES

```bash
# 1. Complete spawning work
echo "✅ Spawned [agent type] for [purpose]"

# 2. Update state file per R324/R288
update_state "[NEXT_STATE]"
commit_state_files_per_r288()

# 3. Save TODOs per R287
save_todos "SPAWNED_[AGENT_TYPE]"

# 4. R322: Stop conversation (context preservation)
echo "🛑 R322: Stopping after spawn for context preservation"

# 5. R405: CONTINUATION FLAG - MUST BE TRUE FOR SPAWNING!
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# 6. Exit to end conversation
exit 0
```

## ❌ WRONG PATTERN (CAUSES TEST FAILURES)

```bash
# ❌ THIS KILLS AUTOMATION - DO NOT DO THIS!
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
exit 0

# Result: Test framework stops, Test 2 fails at iteration 8
```

## 🎯 WHY TRUE IS CORRECT FOR SPAWNING

**Spawning is NORMAL operation:**
- ✅ System knows next state (from state machine)
- ✅ Automation can continue (designed workflow)
- ✅ No manual intervention needed
- ✅ Context preservation ≠ error condition

**The orchestrator stopping (`exit 0`) is for:**
- Preserving context between conversation turns
- Allowing state file commits
- Creating clean state boundaries

**The TRUE flag indicates:**
- Automation CAN restart the conversation
- System knows what to do next (check state file)
- Normal operation is proceeding

## 🔴 WHEN TO USE FALSE (NOT FOR SPAWNING!)

**FALSE should ONLY be used for catastrophic failures:**
- ❌ State file corrupted beyond parsing
- ❌ Critical infrastructure destroyed
- ❌ Unrecoverable system errors
- ❌ **NEVER for normal spawning operations!**

## 📋 SPAWN STATE CHECKLIST

**Before exiting this spawn state, verify:**
1. [ ] All agents spawned successfully
2. [ ] State file updated to next state per R324
3. [ ] State files committed per R288
4. [ ] TODOs saved per R287
5. [ ] R322 stop message displayed
6. [ ] **CONTINUE-SOFTWARE-FACTORY=TRUE emitted** ← Critical!
7. [ ] Exited with `exit 0`

**Missing step 6 = Test 2 failure = -100% grade**

---


## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state-v3.json with new state
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


## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_CODE_REVIEWERS_EFFORT_REVIEW STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
mkdir -p markers/state-verification && touch "markers/state-verification/state_rules_read_orchestrator_SPAWN_CODE_REVIEWERS_EFFORT_REVIEW-$(date +%Y%m%d-%H%M%S)"
echo "$(date +%s) - Rules read and acknowledged for SPAWN_CODE_REVIEWERS_EFFORT_REVIEW" > "markers/state-verification/state_rules_read_orchestrator_SPAWN_CODE_REVIEWERS_EFFORT_REVIEW-$(date +%Y%m%d-%H%M%S)"
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWN_CODE_REVIEWERS_EFFORT_REVIEW WORK UNTIL RULES ARE READ:
- ❌ Start spawn code reviewer agents
- ❌ Start assign review work
- ❌ Start distribute review tasks
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
   ❌ WRONG: "I acknowledge all SPAWN_CODE_REVIEWERS_EFFORT_REVIEW rules"
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

### ✅ CORRECT PATTERN FOR SPAWN_CODE_REVIEWERS_EFFORT_REVIEW:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SPAWN_CODE_REVIEWERS_EFFORT_REVIEW work until:**
### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SPAWN_CODE_REVIEWERS_EFFORT_REVIEW work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING_SWE_PROGRESS YOUR READ TOOL CALLS!**

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

### ⚠️⚠️⚠️ CRITICAL MESSAGE FOR CODE REVIEWERS ⚠️⚠️⚠️
**TELL CODE REVIEWERS: "YOU MUST USE tools/line-counter.sh - NO PARAMETERS NEEDED!"**

**ABSOLUTE REQUIREMENTS FOR CODE REVIEWERS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ✅ Tool AUTO-DETECTS correct base branch - NO PARAMETERS NEEDED!
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER use git diff for line counting
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ❌ NEVER use old -b or -c parameters (tool updated!)

**MESSAGE TO SPAWN: "MUST use tools/line-counter.sh for ALL measurements per R304. Tool auto-detects base - just run it!"**

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

### State-Specific Rules (NOT in orchestrator.md):

1. **🚨🚨🚨 R151** - Parallel Spawning Timestamp Requirement
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-spawning-timing.md`
   - Criticality: CRITICAL - Timestamps must be within 5s for parallel agents
   - Summary: All parallel agents must emit timestamps within 5 seconds

2. **R108** - Code Review Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R108-code-review-protocol.md`
   - Criticality: BLOCKING - Complete review protocol
   - Summary: Review for size limits, quality, patterns, and create reports

3. **R222** - Code Review Gate
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R222-code-review-gate.md`
   - Criticality: BLOCKING - Must create standardized reports
   - Summary: Generate CODE-REVIEW-REPORT.md for all findings

4. **R255** - Post-Agent Work Verification
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R255-POST-AGENT-WORK-VERIFICATION.md`
   - Criticality: BLOCKING - Verify correct locations after completion
   - Summary: Ensure all review work is in correct directories

5. **🔴🔴🔴 R208** - CD Before Agent Spawn (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R208-orchestrator-spawn-cd-protocol.md`
   - Criticality: SUPREME LAW - CD to correct directory before spawn
   - Summary: Must change to agent's working directory before spawning

6. **🚨🚨🚨 R338** - Mandatory Line Count State Tracking
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R338-mandatory-line-count-state-tracking.md`
   - Criticality: BLOCKING - -50% per missing tracking, -100% if none
   - Summary: TELL Code Reviewers to use standardized SIZE MEASUREMENT REPORT format
   - **CRITICAL**: Tell Code Reviewers: "Report 'Implementation Lines:' in standardized format per R338"

7. **🔴🔴🔴 R353** - Cascade Focus Protocol (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R353-cascade-focus-protocol.md`
   - Criticality: SUPREME LAW - NO diversions during CASCADE operations
   - Summary: During CASCADE mode, Code Reviewers MUST skip size checks and splits
   - **CRITICAL**: Check cascade_mode in state file, pass context to Code Reviewers

**Note**: Additional TODO persistence rules (R287) apply from orchestrator.md.

## 🚨 SPAWN_CODE_REVIEWERS_EFFORT_REVIEW IS A VERB - START SPAWNING IMMEDIATELY! 🚨

**See R151 for immediate action requirements when entering this state.**

The SPAWN_CODE_REVIEWERS_EFFORT_REVIEW state requires IMMEDIATE spawning action - no pausing or waiting.

## State Context
You are spawning Code Reviewer agents to review:
- Completed implementation work from SW Engineers
- **INTEGRATE_WAVE_EFFORTS FIXES** when coming from MONITORING_EFFORT_FIXES
- Any code that needs review before proceeding

## 🔴🔴🔴 PREREQUISITES FOR SPAWN_CODE_REVIEWERS_EFFORT_REVIEW 🔴🔴🔴

**BEFORE ENTERING THIS STATE, YOU MUST ALREADY HAVE:**
1. ✅ All SW Engineers completed their implementation
2. ✅ All code committed and pushed by SW Engineers
3. ✅ Size measurements completed and within limits
4. ✅ All effort directories contain implemented code
5. ✅ **PARALLELIZATION ANALYSIS COMPLETE (ANALYZE_CODE_REVIEWER_PARALLELIZATION)**
6. ✅ **Code Reviewer parallelization plan in orchestrator-state-v3.json**

**IF PARALLELIZATION NOT ANALYZED, GO BACK TO ANALYZE_CODE_REVIEWER_PARALLELIZATION!**

## Review Assignment Protocol

### 🔴🔴🔴 CASCADE MODE CHECK (R353) 🔴🔴🔴

**CRITICAL: Check if we're in CASCADE mode BEFORE spawning!**
```bash
# Check CASCADE mode per R353
CASCADE_MODE=$(jq -r '.cascade_coordination.cascade_mode // false' orchestrator-state-v3.json)

if [[ "$CASCADE_MODE" == "true" ]]; then
    echo "🔴🔴🔴 CASCADE MODE ACTIVE - R353 CASCADE FOCUS PROTOCOL 🔴🔴🔴"
    echo "📋 Code Reviewers will:"
    echo "  - SKIP size measurements"
    echo "  - SKIP split evaluations"
    echo "  - ONLY validate rebases"
    echo "  - ONLY check for conflicts/build issues"
    CASCADE_CONTEXT="--cascade-mode=true --skip-size-checks --skip-split-evaluation --rebase-validation-only"
else
    echo "📊 Normal review mode - full validation enabled"
    CASCADE_CONTEXT=""
fi
```

### For Each Code Reviewer to Spawn:
1. **CD to effort directory** (R208 SUPREME LAW)
2. **Spawn with clear review scope AND CASCADE CONTEXT**:
   - Which efforts to review
   - What type of review (CASCADE vs NORMAL)
   - Pass CASCADE_CONTEXT if cascade_mode=true
   - Where to create reports
3. **Track spawn timestamps** (R151 requirement)
4. **Monitor completion** via orchestrator-state-v3.json

### Parallel vs Sequential Spawning:
- **Parallel**: When reviewing independent efforts
- **Sequential**: When reviewing split efforts or dependencies
- **Decision made in**: ANALYZE_CODE_REVIEWER_PARALLELIZATION state

## Expected Deliverables

Each Code Reviewer must produce:
1. **CODE-REVIEW-REPORT.md** in effort directory
2. **Size compliance verification**
3. **Quality assessment**
4. **Recommendations for fixes if needed**

## 🚨🚨🚨 CASCADE MODE SPAWNING (R353 + R354) 🚨🚨🚨

**When CASCADE_MODE is active:**
1. **TELL Code Reviewers**: "CASCADE MODE ACTIVE per R353 - skip size checks"
2. **PASS cascade context**: Include cascade_mode=true in spawn instructions
3. **EXPECT different output**: CASCADE_VALIDATION result, not size measurements
4. **NO split transitions**: Even if reviewer mentions size, NO SPLITS during cascade

**🔴🔴🔴 R354 POST-REBASE REVIEW ENFORCEMENT 🔴🔴🔴**

```bash
# Check for pending post-rebase reviews (R354)
PENDING_REBASE_REVIEWS=$(jq -r '
    .cascade_coordination.pending_reviews[]? |
    select(.review_type == "post_rebase" and .review_status == "pending") |
    .effort' orchestrator-state-v3.json)

if [[ -n "$PENDING_REBASE_REVIEWS" ]]; then
    echo "🔴🔴🔴 R354 POST-REBASE REVIEWS REQUIRED 🔴🔴🔴"
    echo "Spawning Code Reviewers for post-rebase validation:"

    for effort in $PENDING_REBASE_REVIEWS; do
        echo "  📋 $effort - needs post-rebase review"

        # Get rebase details
        REBASED_TO=$(jq -r --arg e "$effort" '
            .cascade_coordination.pending_reviews[] |
            select(.effort == $e and .review_type == "post_rebase") |
            .rebased_to' orchestrator-state-v3.json)

        echo "Spawning Code Reviewer for $effort (rebased to $REBASED_TO)"

        /spawn-agent code-reviewer \
            --cascade-mode=true \
            --review-type=post-rebase \
            --r354-enforcement=true \
            --skip-size-checks \
            --skip-quality-checks \
            --integration-validation-only \
            --effort="$effort" \
            --rebased-to="$REBASED_TO" \
            --message="R354 POST-REBASE REVIEW: Validate integration after rebase to $REBASED_TO. Focus on build/test success ONLY per R353/R354."
    done

    echo "✅ Post-rebase reviewers spawned per R354"
fi
```

**Example spawn with CASCADE + POST-REBASE context:**
```bash
if [[ "$CASCADE_MODE" == "true" && "$REVIEW_TYPE" == "post_rebase" ]]; then
    echo "🔴 Spawning Code Reviewer for R354 POST-REBASE REVIEW"
    /spawn-agent code-reviewer \
        --cascade-mode=true \
        --review-type=post-rebase \
        --r354-mandated=true \
        --skip-size-checks \
        --skip-split-evaluation \
        --rebase-validation-only \
        --effort="$EFFORT_PATH" \
        --message="R354 POST-REBASE REVIEW: Validate rebase success ONLY. Skip all quality/size checks per R353/R354."
elif [[ "$CASCADE_MODE" == "true" ]]; then
    echo "Spawning Code Reviewer with CASCADE FOCUS (R353)"
    /spawn-agent code-reviewer \
        --cascade-mode=true \
        --skip-size-checks \
        --rebase-validation-only \
        --effort="$EFFORT_PATH" \
        --message="CASCADE MODE: Only validate rebase success, skip size checks per R353"
fi
```

## 🚨🚨🚨 SPECIAL CASE: REVIEWING INTEGRATE_WAVE_EFFORTS FIXES 🚨🚨🚨

**When coming from MONITORING_EFFORT_FIXES:**
1. You are reviewing FIXES to integration issues
2. Focus review on:
   - Did the fixes resolve the integration problems?
   - Are the fixes properly implemented?
   - Do the fixes maintain code quality?
3. After successful review of fixes:
   - Transition to MONITOR state
   - Then to WAVE_COMPLETE
   - Then BACK TO INTEGRATE_WAVE_EFFORTS for full re-run
4. **NEVER skip directly to MONITORING_INTEGRATE_WAVE_EFFORTS**

## 🔴🔴🔴 R206 - STATE MACHINE VALIDATION (SUPREME LAW) 🔴🔴🔴

**CRITICAL**: ALL state transitions MUST be validated against state machine!

### Allowed Transitions (from state machine)

Per `state-machines/software-factory-3.0-state-machine.json`:

```json
"SPAWN_CODE_REVIEWERS_EFFORT_REVIEW": {
  "allowed_transitions": [
    "MONITORING_EFFORT_REVIEWS",
    "ERROR_RECOVERY"
  ]
}
```

- ✅ **MONITORING_EFFORT_REVIEWS** - ONLY valid normal transition (monitor reviews, evaluate bugs)
- ✅ **ERROR_RECOVERY** - Only if catastrophic failure occurs
- ❌ **NOT ALLOWED**: WAVE_COMPLETE (must go through MONITORING_EFFORT_REVIEWS first!)
- ❌ **NOT ALLOWED**: Any other state

### Next State Determination

**MANDATORY**: Next state is ALWAYS `MONITORING_EFFORT_REVIEWS` after spawning reviewers.

The MONITORING_EFFORT_REVIEWS state will:
1. Monitor code reviewer progress
2. Collect review results
3. Evaluate bugs found
4. Decide next action based on review outcomes:
   - If bugs_found > 0 → CREATE_EFFORT_FIX_PLAN
   - If bugs_found == 0 → WAVE_COMPLETE

**THIS STATE DOES NOT DECIDE WAVE COMPLETION!** That decision belongs to MONITORING_EFFORT_REVIEWS.

### Validation Code

```bash
# MANDATORY validation before proposing next state
CURRENT_STATE="SPAWN_CODE_REVIEWERS_EFFORT_REVIEW"
NEXT_STATE="MONITORING_EFFORT_REVIEWS"

# Verify transition is allowed
if ! jq -e ".states.\"$CURRENT_STATE\".allowed_transitions | contains([\"$NEXT_STATE\"])" \
        /home/vscode/software-factory-template/state-machines/software-factory-3.0-state-machine.json > /dev/null; then
    echo "❌ CRITICAL: Illegal state transition!"
    echo "Current: $CURRENT_STATE"
    echo "Attempted: $NEXT_STATE"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=ILLEGAL_TRANSITION"
    exit 1
fi

echo "✅ Transition validated: $CURRENT_STATE → $NEXT_STATE"
```

## State Transitions

From SPAWN_CODE_REVIEWERS_EFFORT_REVIEW (per state machine JSON):
- **MONITORING_EFFORT_REVIEWS** ← ONLY valid normal path
- **ERROR_RECOVERY** ← Only if catastrophic error

**NEVER** transition directly to:
- ❌ WAVE_COMPLETE (must monitor reviews first per state machine)
- ❌ SPAWN_SW_ENGINEERS (handled by fix workflow in later states)
- ❌ Any other state not in allowed_transitions

## Critical Rule Enforcement Order

1. **R290**: Read these state rules FIRST and create verification marker
3. **R208**: CD to correct directory BEFORE spawning
4. **R151**: Ensure parallel timestamps within 5s
5. **R053**: Follow complete review protocol
6. **R054**: Generate standardized reports
7. **R255**: Verify work in correct locations

**Remember**: Code Reviewers check for size violations, quality issues, and create actionable reports for the orchestrator.

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete SPAWN_CODE_REVIEWERS_EFFORT_REVIEW:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="SPAWN_CODE_REVIEWERS_EFFORT_REVIEW complete - [describe what was accomplished]"
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
  "state_completed": "SPAWN_CODE_REVIEWERS_EFFORT_REVIEW",
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
  --current-state "SPAWN_CODE_REVIEWERS_EFFORT_REVIEW" \
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

if ! git commit -m "state: SPAWN_CODE_REVIEWERS_EFFORT_REVIEW → $NEXT_STATE - SPAWN_CODE_REVIEWERS_EFFORT_REVIEW complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: SPAWN_CODE_REVIEWERS_EFFORT_REVIEW"
    echo "Attempted transition from: SPAWN_CODE_REVIEWERS_EFFORT_REVIEW"
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
save_todos "SPAWN_CODE_REVIEWERS_EFFORT_REVIEW_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - SPAWN_CODE_REVIEWERS_EFFORT_REVIEW complete [R287]"; then
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
PROPOSED_NEXT_STATE="MONITORING_EFFORT_REVIEWS"
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

