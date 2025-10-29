# Orchestrator - WAITING_FOR_EFFORT_PLANS State Rules


## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### 🛑 R322 CHECKPOINT CLARIFICATION - CRITICAL UNDERSTANDING

**AGENT STOPS ≠ FACTORY STOPS** (These are TWO INDEPENDENT operations!)

At R322 checkpoints:
- ✅ **Agent STOPS work** = Agent completes state, exits process, preserves context
- ✅ **Factory CONTINUES** = Set CONTINUE-SOFTWARE-FACTORY=TRUE (automation proceeds)
- ✅ **User resumes** = Run /continue-orchestrating to start next state
- ✅ **This is NORMAL** = Expected workflow every time!

**THE CONFUSION TO AVOID:**
- ❌ WRONG: "I must stop (R322) → Manual intervention needed → FALSE"
- ✅ RIGHT: "I must stop (R322) → Normal checkpoint → TRUE → User continues me"

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state-v3.json with new state
3. ✅ Committing and pushing the state file
4. ✅ Providing work summary
5. ✅ **Setting CONTINUE-SOFTWARE-FACTORY=TRUE** (normal operation!)

### YOU MUST NOT:
- ❌ Continue to the next state automatically
- ❌ Start work for the new state
- ❌ Spawn agents for the new state
- ❌ Assume permission to continue
- ❌ **Set CONTINUE-SOFTWARE-FACTORY=FALSE** (unless catastrophic error!)

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

CONTINUE-SOFTWARE-FACTORY=TRUE
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**
**TRUE MEANS TRUE - Factory automation continues (normal!)**

---


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED WAITING_FOR_EFFORT_PLANS STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_WAITING_FOR_EFFORT_PLANS
echo "$(date +%s) - Rules read and acknowledged for WAITING_FOR_EFFORT_PLANS" > .state_rules_read_orchestrator_WAITING_FOR_EFFORT_PLANS
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY WAITING_FOR_EFFORT_PLANS WORK UNTIL RULES ARE READ:
- ❌ Start check effort plan status
- ❌ Start monitor reviewer progress
- ❌ Start collect completed plans
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**
### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_EFFORT_PLANS STATE

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
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state-v3.json before EVERY state transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

5. **🔴🔴🔴 R233** - Immediate Action On State Entry
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-all-states-immediate-action.md`
   - Criticality: SUPREME LAW - Must act immediately on entering state
   - Summary: WAITING states require active monitoring, not passive waiting


### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R151, R208, R053..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all WAITING_FOR_EFFORT_PLANS rules"
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

### ✅ CORRECT PATTERN FOR WAITING_FOR_EFFORT_PLANS:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute WAITING_FOR_EFFORT_PLANS work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY WAITING_FOR_EFFORT_PLANS work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute WAITING_FOR_EFFORT_PLANS work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with WAITING_FOR_EFFORT_PLANS work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY WAITING_FOR_EFFORT_PLANS work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING_SWE_PROGRESS YOUR READ TOOL CALLS!**

## 📋 RULE SUMMARY FOR WAITING_FOR_EFFORT_PLANS STATE

### Rules Enforced in This State:
- R234: Mandatory State Traversal [SUPREME LAW - Part of sequence]
- R255: Post-Agent Work Verification [BLOCKING - Check every completion]
- R322: Never Stop Monitoring [SUPREME LAW - Keep checking]
- R287: TODO Save Frequency [BLOCKING - Every 10 messages/15 min]
- R288: State File Update and Commit [SUPREME LAW - Track progress]

### Critical Requirements:
1. Actively poll for plans NOW - Penalty: -30%
2. Check every 5-10 seconds - Penalty: -20%
3. Verify R255 for each completion - Penalty: -100%
4. Save TODOs every 15 minutes - Penalty: -15%
5. Transition to ANALYZE_IMPLEMENTATION_PARALLELIZATION - Penalty: -100%

### Success Criteria:
- ✅ All IMPLEMENTATION-PLAN.md files created
- ✅ All plans in correct directories (R255)
- ✅ All plans committed and pushed
- ✅ Work logs updated

### Failure Triggers:
- ❌ Skip to SPAWN_SW_ENGINEERS = -100% R234 VIOLATION
- ❌ Accept plans in wrong location = R255 VIOLATION
- ❌ Stop monitoring = R322 VIOLATION
- ❌ Forget TODO saves = -15% per violation

## 🚨 WAITING_FOR_EFFORT_PLANS IS A VERB - START ACTIVELY CHECKING IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING WAITING_FOR_EFFORT_PLANS

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Poll effort directories for IMPLEMENTATION-PLAN.md NOW
2. Check every 5-10 seconds for completion
3. Check TodoWrite for pending items and process them
4. Report status of each effort immediately

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in WAITING_FOR_EFFORT_PLANS" [stops]
- ❌ "Successfully entered WAITING_FOR_EFFORT_PLANS state" [waits]
- ❌ "Ready to start actively checking" [pauses]
- ❌ "I'm in WAITING_FOR_EFFORT_PLANS state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering WAITING_FOR_EFFORT_PLANS, Poll effort directories for IMPLEMENTATION-PLAN.md NOW..."
- ✅ "START ACTIVELY CHECKING, check every 5-10 seconds for completion..."
- ✅ "WAITING_FOR_EFFORT_PLANS: Report status of each effort immediately..."

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## SF 3.0 State Monitoring Context

This waiting state actively monitors SF 3.0 state files:
- Polls orchestrator-state-v3.json to check `spawned_agents` section for Code Reviewer completion status
- Reads `metadata_locations.effort_plans` to locate newly created implementation plans per R340
- Updates `state_machine.current_state` when all plans are detected and validated
- Tracks plan creation progress and updates orchestrator-state-v3.json atomically per R288

The state actively checks for completion (R233) rather than passively waiting, monitoring state files for changes every 5-10 seconds.

## State Context
You are waiting for Code Reviewers to complete individual effort implementation plans.

## 🔴🔴🔴 SUPREME LAW R234 - STAY IN SEQUENCE 🔴🔴🔴

### YOUR POSITION IN THE MANDATORY SEQUENCE:
```
CREATE_NEXT_INFRASTRUCTURE (✓ completed)
    ↓
ANALYZE_CODE_REVIEWER_PARALLELIZATION (✓ completed)
    ↓
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING (✓ completed)
    ↓
WAITING_FOR_EFFORT_PLANS (👈 YOU ARE HERE)
    ↓ (MUST GO HERE NEXT)
ANALYZE_IMPLEMENTATION_PARALLELIZATION
    ↓
SPAWN_SW_ENGINEERS
```

**NOW:** Actively monitor Code Reviewers
**NEXT:** You MUST go to ANALYZE_IMPLEMENTATION_PARALLELIZATION
**FORBIDDEN:** Skipping analysis to go directly to SPAWN_SW_ENGINEERS = -100%

## 🚨 CRITICAL: R356 DOES NOT SKIP INFRASTRUCTURE CREATION! 🚨

### R356 Optimization Scope - PRECISE DEFINITION

**✅ R356 APPLIES TO: Parallelization Analysis Complexity**

When Code Reviewers finish creating effort plans:
- **Single Effort**:
  - No complex dependency analysis needed
  - ANALYZE_IMPLEMENTATION_PARALLELIZATION is simpler
  - Decision is straightforward: spawn 1 SW Engineer

- **Multiple Efforts**:
  - Must analyze inter-effort dependencies
  - ANALYZE_IMPLEMENTATION_PARALLELIZATION is complex
  - Must determine sequential vs parallel spawn strategy

**❌ R356 DOES NOT APPLY TO: Infrastructure Creation**

Infrastructure creation is MANDATORY for ALL waves:
- ❌ NEVER skip CREATE_NEXT_INFRASTRUCTURE
- ❌ NEVER skip VALIDATE_INFRASTRUCTURE
- ❌ NEVER assume infrastructure exists
- ✅ ALWAYS verify branch existence before spawn
- ✅ ALWAYS go through full infrastructure sequence

### Mandatory Sequence After Effort Planning

**ALWAYS - NO EXCEPTIONS:**
```
WAITING_FOR_EFFORT_PLANS
  ↓
ANALYZE_IMPLEMENTATION_PARALLELIZATION  ← MANDATORY (even for single effort!)
  ↓                                       (checks if infrastructure exists)
CREATE_NEXT_INFRASTRUCTURE (if needed)  ← MANDATORY for new waves
  ↓
VALIDATE_INFRASTRUCTURE                 ← MANDATORY for new waves
  ↓
SPAWN_SW_ENGINEERS                      ← ONLY after infrastructure validated
```

**R356 Affects**:
- Complexity of ANALYZE_IMPLEMENTATION_PARALLELIZATION
- Time spent analyzing dependencies
- Spawn strategy decision

**R356 Does NOT Affect**:
- Whether to create infrastructure (always required)
- Whether to validate infrastructure (always required)
- Whether to go through state sequence (always required)

### Common Misunderstandings - AVOID THESE!

❌ **WRONG**: "R356 says single effort can skip to SPAWN_SW_ENGINEERS"
✅ **RIGHT**: "R356 says single effort has simpler analysis, but still must analyze"

❌ **WRONG**: "R356 optimization means skip intermediate states"
✅ **RIGHT**: "R356 optimization means faster execution of required states"

❌ **WRONG**: "Wave 1 worked by skipping infrastructure, so Wave 2 can too"
✅ **RIGHT**: "Wave 1 created infrastructure via different path, Wave 2 must create via its path"

### Why This Matters

**Wave 1 Pattern** (OLD - pre SF 3.0):
- Created infrastructure BEFORE effort planning
- Effort plans created AFTER branches exist
- Infrastructure reused for planning

**Wave 2+ Pattern** (NEW - SF 3.0):
- Create effort plans FIRST
- Then create infrastructure BASED ON plans
- More flexible, better parallelization

**The Difference**:
- Wave 1: Infrastructure → Plans → Spawn
- Wave 2: Plans → Infrastructure → Spawn

**Both patterns require infrastructure creation - just at different times!**

## Monitoring Requirements (R340 Compliant)

### 🚨🚨🚨 RULE R340: PLANNING FILE METADATA TRACKING (BLOCKING)
- **MUST** read plan locations from orchestrator-state-v3.json
- **NEVER** search directories for planning files
- **ALWAYS** use effort_repo_files.effort_plans section
- **VIOLATION = -20% for each untracked file**

```bash
# R340 Compliant: Check status of effort plans from state
check_effort_plan_status() {
    local PHASE=$1 WAVE=$2
    local ALL_COMPLETE=true
    
    echo "📊 Checking effort plan status (R340 compliant)..."
    
    # Get list of efforts for this wave from state
    EFFORTS=$(echo "✅ State file updated to: $NEXT_STATE"
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

if ! git commit -m "state: WAITING_FOR_EFFORT_PLANS → $NEXT_STATE - WAITING_FOR_EFFORT_PLANS complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: WAITING_FOR_EFFORT_PLANS"
    echo "Attempted transition from: WAITING_FOR_EFFORT_PLANS"
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
save_todos "WAITING_FOR_EFFORT_PLANS_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - WAITING_FOR_EFFORT_PLANS complete [R287]"; then
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

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### ✅ WHEN TO USE TRUE (99% of cases - including R322 checkpoints!):
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met
- ✅ **R322 checkpoint reached** (NORMAL - agent stops but factory continues!)
- ✅ **All plans tracked in state file** (PROJECT_DONE - proceed to parallelization analysis)
- ✅ **Waiting for /continue-orchestrating** (EXPECTED workflow)

### ❌ WHEN TO USE FALSE (catastrophic system failures only!):
- ❌ CATASTROPHIC error preventing ANY continuation
- ❌ State file corrupted beyond repair
- ❌ Git repository completely broken
- ❌ Data corruption that would cascade
- ❌ **NOT for R322 checkpoints** (agent stopping is NORMAL → use TRUE)
- ❌ **NOT for "manual intervention required"** (if it's a checkpoint → use TRUE)
- ❌ **NOT for "waiting for plans"** (normal monitoring → use TRUE)
- ❌ **NOT for "being cautious"** (use TRUE)

### 🔴 CRITICAL: UNDERSTAND THE DIFFERENCE
**Agent stopping (R322) ≠ Factory stopping (R405 FALSE)**

In WAITING_FOR_EFFORT_PLANS:
- When all plans are tracked → Agent stops (R322) AND Factory continues (TRUE)
- This is the NORMAL, EXPECTED workflow!
- FALSE would only be used if the system is catastrophically broken

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

