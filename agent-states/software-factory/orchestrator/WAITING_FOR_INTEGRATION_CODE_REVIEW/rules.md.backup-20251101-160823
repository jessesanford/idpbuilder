# WAITING_FOR_INTEGRATION_CODE_REVIEW State Rules


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_INTEGRATION_CODE_REVIEW STATE

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
   - Summary: Monitoring states require active checking, not passive waiting

6. **🔴🔴🔴 R358** - Integration Completion Detection and Automatic Transition
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R358-integration-completion-detection.md`
   - Criticality: SUPREME LAW - Must detect completion and transition within 30 seconds
   - Summary: CANNOT remain in WAITING_FOR_INTEGRATION_CODE_REVIEW after integration completes

7. **🔴🔴🔴 R304** - Mandatory Line Counting Tool
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
   - Criticality: SUPREME LAW - Line counting requirements
   - Summary: MUST use tools/line-counter.sh for ALL measurements

## State Purpose
Monitor Code Reviewer progress on integration code review and process results when complete.

## Entry Criteria
- **From**: REVIEW_WAVE_INTEGRATION
- **Condition**: Code Reviewer spawned for integration review
- **Required**: State file shows review in progress

## Active Monitoring Requirements

### 🔴🔴🔴 R233/R358 ENFORCEMENT: ACTIVE MONITORING REQUIRED 🔴🔴🔴

**MANDATORY**: This is NOT a passive waiting state - you MUST actively monitor for completion!

```bash
# START THE R358 MONITORING LOOP IMMEDIATELY upon entering this state
echo "📊 Entering WAITING_FOR_INTEGRATION_CODE_REVIEW - starting R358 monitoring loop"

# Start the completion detection function
monitor_integration_review_completion &
MONITOR_PID=$!

echo "✅ R358 monitoring loop started (PID: $MONITOR_PID)"
echo "Will check for completion every 30 seconds and transition immediately when detected"
```

### R358 Continuous Completion Detection

**MUST check for integration review completion every 30 seconds and transition IMMEDIATELY!**

```bash
# R358 ENFORCEMENT: Continuous completion detection loop
monitor_integration_review_completion() {
    while true; do
        # Check every 30 seconds per R358
        sleep 30

        # Method 1: Check for report file
        if [ -f "REVIEW_WAVE_INTEGRATION_REPORT.md" ]; then
            echo "✅ Integration code review complete - report found"
            process_review_results
            transition_immediately
            break
        fi

        # Method 2: Check integration_review status in state file
        REVIEW_STATUS=$(jq -r '.integration_review.status // "in_progress"' orchestrator-state-v3.json)

        if [ "$REVIEW_STATUS" = "complete" ]; then
            echo "✅ Integration code review marked complete in state file"
            process_review_results
            transition_immediately
            break
        fi

        echo "⏳ Integration code review still in progress..."
    done
}
```

**NEVER** just wait passively - the monitoring must be ACTIVE and CONTINUOUS per R233!

## State Actions

### 1. IMMEDIATE: Check for Review Completion
```bash
# Check for integration code review report
if [ -f "REVIEW_WAVE_INTEGRATION_REPORT.md" ]; then
    echo "Integration code review complete"
    process_review_results
else
    echo "Waiting for integration code review to complete"
    exit 0  # Will be re-invoked later
fi
```

### 2. Process Review Results
When report exists:
- Parse REVIEW_WAVE_INTEGRATION_REPORT.md
- Extract PASS/FAIL status
- Identify any integration issues
- Determine next state based on results

### 3. Update State File
```json
{
  "current_state": "WAITING_FOR_INTEGRATION_CODE_REVIEW",
  "integration_review": {
    "status": "complete",
    "result": "PASS|FAIL",
    "issues_found": [],
    "report": "REVIEW_WAVE_INTEGRATION_REPORT.md"
  }
}
```

## Exit Criteria

### Success Path (PASS) → REVIEW_WAVE_ARCHITECTURE
- Review passed with no critical issues
- Integration quality verified
- Ready for architect review

### Failure Path (FAIL) → SPAWN_CODE_REVIEWER_INTEGRATE_WAVE_EFFORTS_FIX_PLAN
- Critical integration issues found
- Fixes required before proceeding
- Need plan to address issues

### Waiting Path → WAITING_FOR_INTEGRATION_CODE_REVIEW
- Review still in progress
- Exit and wait for next check

## Report Processing
Parse REVIEW_WAVE_INTEGRATION_REPORT.md for:
- Overall status (PASS/FAIL)
- Critical issues list
- Merge conflict concerns
- Test failures after integration
- Cross-effort conflicts

## Rules Enforced
- R233: Immediate check upon entry
- R238: Monitor for completion
- R321: Fixes require backport planning

## Transition Rules

### 🔴🔴🔴 CASCADE MODE CHECK - R351 ENFORCEMENT 🔴🔴🔴
```bash
# Check if we're in cascade mode before determining next state
CASCADE_MODE=$(echo "✅ State file updated to: $NEXT_STATE"
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

if ! git commit -m "state: WAITING_FOR_INTEGRATION_CODE_REVIEW → $NEXT_STATE - WAITING_FOR_INTEGRATION_CODE_REVIEW complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: WAITING_FOR_INTEGRATION_CODE_REVIEW"
    echo "Attempted transition from: WAITING_FOR_INTEGRATION_CODE_REVIEW"
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
save_todos "WAITING_FOR_INTEGRATION_CODE_REVIEW_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - WAITING_FOR_INTEGRATION_CODE_REVIEW complete [R287]"; then
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
PROPOSED_NEXT_STATE="REVIEW_WAVE_ARCHITECTURE"
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

