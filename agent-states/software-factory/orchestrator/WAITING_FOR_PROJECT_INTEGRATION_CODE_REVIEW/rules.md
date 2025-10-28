# WAITING_FOR_PROJECT_REVIEW_WAVE_INTEGRATION State Rules



## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_PROJECT_REVIEW_WAVE_INTEGRATION STATE

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

6. **🔴🔴🔴 R358** - 30-Second Project Review Completion Detection
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R358-completion-detection-30-second.md`
   - Criticality: SUPREME LAW - Must detect project review completion within 30 seconds
   - Summary: Active monitoring loop required to detect project review completion immediately

## State Purpose
Monitor Code Reviewer progress on project integration code review and process comprehensive results when complete.

## Entry Criteria
- **From**: PROJECT_REVIEW_WAVE_INTEGRATION
- **Condition**: Code Reviewer spawned for project integration review
- **Required**: State file shows project review in progress

## Active Monitoring Requirements

**🔴 R233 + R358: Active Monitoring Loop (SUPREME LAW)**

This WAITING state must ACTIVELY monitor for project review completion, not passively wait. The monitoring loop must:

1. **Check every 30 seconds** (R358 requirement)
2. **Detect completion immediately** when PROJECT_REVIEW_WAVE_INTEGRATION_REPORT.md appears
3. **Process results automatically** without delay
4. **Transition to next state** based on review outcome

### Monitoring Loop Implementation
```bash
# R233/R358 Monitoring Loop for Project Review Completion
while true; do
    # Check for project review report
    if [ -f "PROJECT_REVIEW_WAVE_INTEGRATION_REPORT.md" ]; then
        echo "✅ Project integration review complete (detected within 30s per R358)"
        process_project_review_results
        break
    fi

    # Check Code Reviewer agent status
    REVIEWER_PID=$(get_code_reviewer_pid)
    if ! ps -p "$REVIEWER_PID" > /dev/null 2>&1; then
        echo "⚠️ Code Reviewer agent completed or failed"
        if [ -f "PROJECT_REVIEW_WAVE_INTEGRATION_REPORT.md" ]; then
            echo "✅ Report found - proceeding"
            process_project_review_results
            break
        else
            echo "❌ No report found - Code Reviewer may have failed"
            # Create blocker and transition to ERROR_RECOVERY
            create_blocker "Code Reviewer produced no project review report"
            PROPOSED_NEXT_STATE="ERROR_RECOVERY"
            break
        fi
    fi

    # R358: Check every 30 seconds
    echo "Waiting for project review completion... (checking every 30s per R358)"
    sleep 30
done
```

### R233 Immediate Action Requirement
Upon entering this state, IMMEDIATELY:
1. Check if PROJECT_REVIEW_WAVE_INTEGRATION_REPORT.md already exists
2. If exists: Process immediately (review finished before we checked)
3. If not: Enter 30-second monitoring loop (R358)
4. **NEVER** exit without taking action or scheduling next check

## State Actions

### 1. IMMEDIATE: Check for Project Review Completion
```bash
# Check for project integration code review report
if [ -f "PROJECT_REVIEW_WAVE_INTEGRATION_REPORT.md" ]; then
    echo "Project integration code review complete"
    process_project_review_results
else
    echo "Waiting for project integration code review to complete"
    exit 0  # Will be re-invoked later
fi
```

### 2. Process Project Review Results
When report exists:
- Parse PROJECT_REVIEW_WAVE_INTEGRATION_REPORT.md
- Extract PASS/FAIL/CONDITIONAL_PASS status
- Analyze critical issues across project
- Verify all requirements met
- Assess production readiness
- Determine next state based on comprehensive results

### 3. Update State File
```json
{
  "current_state": "WAITING_FOR_PROJECT_REVIEW_WAVE_INTEGRATION",
  "project_integration_review": {
    "status": "complete",
    "result": "PASS|FAIL|CONDITIONAL_PASS",
    "critical_issues": [],
    "requirements_met": true,
    "production_ready": true,
    "cross_phase_conflicts": [],
    "performance_status": "acceptable",
    "security_status": "passed",
    "technical_debt": "low",
    "report": "PROJECT_REVIEW_WAVE_INTEGRATION_REPORT.md"
  }
}
```

## Exit Criteria

### Success Path (PASS) → SPAWN_CODE_REVIEWER_DEMO_VALIDATION
- Review passed with no critical issues
- Project quality verified
- All requirements met
- Ready for final validation

### Conditional Pass → SPAWN_CODE_REVIEWER_DEMO_VALIDATION
- Minor issues identified but not blocking
- Document known issues for tracking
- Proceed with validation

### Failure Path (FAIL) → SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING
- Critical project-level issues found
- Requirements not met
- Need comprehensive fix plan
- May require multiple phases of fixes

### Waiting Path → WAITING_FOR_PROJECT_REVIEW_WAVE_INTEGRATION
- Review still in progress
- Exit and wait for next check

## Report Processing
Parse PROJECT_REVIEW_WAVE_INTEGRATION_REPORT.md for:
- Overall project status
- Critical blocking issues
- Requirements completion percentage
- Cross-phase integration problems
- Performance regression analysis
- Security vulnerabilities
- Technical debt assessment
- Test coverage gaps
- Production readiness score

## Decision Matrix
```
PASS: All criteria met, <5 minor issues
CONDITIONAL_PASS: 5-10 minor issues, no blockers
FAIL: Any critical issue OR >10 minor issues
```

## Rules Enforced
- R233: Immediate check upon entry
- R238: Monitor for completion
- R283: Project completeness validation
- R266: Comprehensive validation before completion
- R321: Fixes require backport to all phases

## Transition Rules
- **If PASS** → SPAWN_CODE_REVIEWER_DEMO_VALIDATION
- **If CONDITIONAL_PASS** → SPAWN_CODE_REVIEWER_DEMO_VALIDATION (with notes)
- **If FAIL** → SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING
- **If pending** → WAITING_FOR_PROJECT_REVIEW_WAVE_INTEGRATION (self)

## Special Handling
- This is the most critical review gate
- May require escalation for CONDITIONAL_PASS
- Document all issues for future reference
- Consider rollback if severe issues found

## Exit Conditions and Continuation Flag

**⚠️ READ THIS:** R405-CONTINUATION-FLAG-MASTER-GUIDE.md before setting flag!

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (Standard Operation - DEFAULT)

Use TRUE when:
- ✅ Review report found and processed
- ✅ Review result: PASS (success, proceed to validation)
- ✅ Review result: CONDITIONAL_PASS (proceed with notes)
- ✅ Review result: FAIL/NEEDS_FIXES (enter fix protocol - THIS IS NORMAL!)
- ✅ Still waiting for review (normal workflow)
- ✅ Transitioning to next appropriate state

**CRITICAL: ALL REVIEW OUTCOMES USE TRUE!**

The entire PURPOSE of code review is to find issues so they can be fixed through
the AUTOMATIC review-fix cycle. This is DESIGNED BEHAVIOR.

**Review results and correct flags:**
- Review PASS → TRUE (success, proceed to validation)
- Review CONDITIONAL_PASS → TRUE (proceed with documented issues)
- Review FAIL/NEEDS_FIXES → TRUE (normal, enter fix protocol)
- Review still in progress → TRUE (normal waiting)

**ALL of these are normal operations with defined recovery paths!**

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (Exceptional - EXTREMELY RARE)

Use FALSE ONLY when:
- ❌ Review report file missing AND Code Reviewer process dead/failed
- ❌ Review report completely corrupt (cannot parse ANY content)
- ❌ State machine corruption detected
- ❌ No way to determine review outcome
- ❌ Truly unrecoverable error

**DO NOT set FALSE because:**
- ❌ Review found issues (THIS IS NORMAL! Use TRUE!)
- ❌ Review failed/needs fixes (THIS IS EXPECTED! Use TRUE!)
- ❌ Code needs improvement (THIS IS THE POINT! Use TRUE!)
- ❌ R322 checkpoint (stop ≠ FALSE flag!)
- ❌ "User might want to see review" (NO! Use TRUE!)

## The Fix Protocol - Why TRUE is Correct

**USER'S WORDS: "WE HAVE A FIX PROTOCOL THAT IS AUTOMATIC FOR THIS REASON!"**

**When review finds issues, the system automatically:**
1. Transitions to fix planning state (automatic)
2. Spawns Code Reviewer to plan fixes (automatic)
3. Spawns SW Engineers to implement fixes (automatic)
4. Re-reviews after fixes (automatic)
5. Repeats until approved (automatic)

**This is the DESIGNED WORKFLOW. It's AUTOMATIC.**

Setting FALSE breaks this automation and requires manual intervention when
the system is perfectly capable of handling it automatically.

## R322 vs Continuation Flag

**R322 requires:**
- Stop conversation (`exit 0`) ✅
- Save state ✅
- Emit flag ✅

**R322 does NOT require:**
- Setting FALSE for normal operations ❌
- Human review of standard workflow ❌
- Halting the fix-review cycle ❌

## Examples

**✅ CORRECT - Review needs fixes:**
```bash
REVIEW_STATUS="NEEDS_FIXES"
# This is NORMAL - system has fix protocol
echo "Review found fixable issues - entering fix workflow"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"  # Let automation work!
```

**✅ CORRECT - Review passed:**
```bash
REVIEW_STATUS="PASS"
# Success - proceed to validation
echo "Review approved - proceeding to validation"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"  # Continue normally
```

**✅ CORRECT - Still waiting:**
```bash
# Review not ready yet
echo "Still waiting for review completion"
exit 0  # Stop for context preservation
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"  # Normal waiting operation
```

**🔴 WRONG - Setting FALSE for fixes needed:**
```bash
REVIEW_STATUS="NEEDS_FIXES"
# VIOLATION: Stopping automation for normal workflow
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"  # WRONG! System can handle this!
```

## Grading Impact

- Using FALSE when review finds issues: -20% (defeats automation)
- Pattern of FALSE for normal review outcomes: -50%
- Complete automation defeat: -100%

**The user is FRUSTRATED:** "WE HAVE A FIX PROTOCOL THAT IS AUTOMATIC FOR THIS REASON!"

Let it work!

## References

- R405-CONTINUATION-FLAG-MASTER-GUIDE.md (THE definitive guide)
- R405: Automation Flag Continuation Principle
- R322: Mandatory Stop Before State Transitions


## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete WAITING_FOR_PROJECT_REVIEW_WAVE_INTEGRATION:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, set variables for State Manager
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="[REASON_FOR_TRANSITION]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for State Transition (R288 - SUPREME LAW)
```bash
# State Manager handles ALL state file updates atomically
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE_NAME" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"

# State Manager will:
# 1. Validate the transition against state machine
# 2. Update all 4 state tracking locations atomically:
#    - orchestrator-state-v3.json
#    - orchestrator-state-demo.json
#    - .cascade-state-backup.json
#    - .orchestrator-state-v3.json
# 3. Commit and push all changes
# 4. Return control

echo "✅ State Manager completed transition"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before exit (R287 trigger)
save_todos "STATE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - state complete [R287]"; then
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

### ✅ Step 5: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 6: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No State Manager spawn = state machine broken (R288 violation, -100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)


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
PROPOSED_NEXT_STATE="NEXT_STATE"
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
### 🚨 WAITING STATE PATTERN - CRITICAL UNDERSTANDING 🚨

**This is a WAITING state. Common source of incorrect FALSE usage!**

**WRONG interpretation:**
> "R322 mandates stop before transition"
> "State work is complete (validation done)"
> "User needs to invoke /continue-orchestrating"
> "Therefore I must set CONTINUE-SOFTWARE-FACTORY=FALSE"

**CORRECT interpretation:**
> "R322 checkpoint is NORMAL procedure for context preservation"
> "State work completed successfully = NORMAL outcome"
> "Waiting for /continue is DESIGNED user experience"
> "System KNOWS next step from state file"
> "NO manual intervention required, just normal continuation"
> "Therefore set CONTINUE-SOFTWARE-FACTORY=TRUE"

**The key distinction:**
- **Stopping inference** (`exit 0`) = Context management (ALWAYS at R322 points)
- **Continuation flag** = Can automation proceed? (TRUE unless catastrophic failure)

**ONLY use FALSE if:**
- ❌ The thing we're waiting for completely disappeared (agents crashed with no recovery)
- ❌ Results arrived but are completely corrupted/unreadable
- ❌ State file corruption prevents determining what to wait for
- ❌ System deadlock with no automated resolution
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

