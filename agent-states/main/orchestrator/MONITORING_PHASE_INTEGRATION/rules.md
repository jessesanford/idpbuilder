# Orchestrator - MONITORING_PHASE_INTEGRATION State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR MONITORING_PHASE_INTEGRATION STATE

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

5. **🔴🔴🔴 R233** - Immediate Action On State Entry
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-immediate-action-on-state-entry.md`
   - Criticality: SUPREME LAW - Must act immediately on entering state
   - Summary: Monitoring states require active checking, not passive waiting

6. **🔴🔴🔴 R358** - Integration Completion Detection and Automatic Transition
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R358-integration-completion-detection.md`
   - Criticality: SUPREME LAW - Must detect completion and transition within 30 seconds
   - Summary: CANNOT remain in monitoring state after integration completes

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

# ORCHESTRATOR STATE: MONITORING_PHASE_INTEGRATION

## 🔴🔴🔴 SUPREME DIRECTIVE: PHASE INTEGRATION FEEDBACK ENFORCEMENT 🔴🔴🔴

**YOU MUST CHECK FOR PHASE INTEGRATION REPORTS AND ACT ON FAILURES!**

## State Overview

In MONITORING_PHASE_INTEGRATION, you are monitoring the Integration Agent's phase-level integration progress and MUST check for phase integration reports.

## Required Actions

### 1. Monitor Phase Integration Agent Progress with R358 Enforcement

#### 🔴🔴🔴 R358: CONTINUOUS COMPLETION DETECTION REQUIRED 🔴🔴🔴

```bash
# R358 ENFORCEMENT: Continuous phase integration completion detection
monitor_phase_integration_completion() {
    while true; do
        # Check every 30 seconds per R358
        sleep 30

        # Method 1: Check integration_status in state file
        PHASE=$(jq '.current_phase' orchestrator-state.json)
        PHASE_STATUS=$(jq -r ".phase_integration_status.phase_${PHASE}.status // \"unknown\"" orchestrator-state.json)
        COMPLETED_AT=$(jq -r ".phase_integration_status.phase_${PHASE}.completed_at // \"\"" orchestrator-state.json)

        # Method 2: Check for phase integration report
        REPORT_FILE=$(jq -r ".metadata_locations.phase_integration_reports.\"phase${PHASE}\".file_path // \"\"" orchestrator-state.json)

        # Method 3: Check integration agent process
        INTEGRATION_PID=$(pgrep -f "integration-agent.*phase" || echo "")

        # DETERMINE IF PHASE INTEGRATION IS COMPLETE
        INTEGRATION_COMPLETE=false

        if [[ "$PHASE_STATUS" == "completed" ]] || [[ -n "$COMPLETED_AT" ]]; then
            echo "🔴 R358: Phase integration marked as completed in state file!"
            INTEGRATION_COMPLETE=true
        elif [[ -f "$REPORT_FILE" ]]; then
            echo "🔴 R358: Phase integration report exists - integration complete!"
            INTEGRATION_COMPLETE=true
        elif [[ -z "$INTEGRATION_PID" ]] && [[ -f "$REPORT_FILE" ]]; then
            echo "🔴 R358: Integration agent not running and report exists - complete!"
            INTEGRATION_COMPLETE=true
        fi

        # TRANSITION IMMEDIATELY IF COMPLETE
        if [[ "$INTEGRATION_COMPLETE" == "true" ]]; then
            echo "🔴🔴🔴 R358: PHASE INTEGRATION COMPLETE - TRANSITIONING NOW 🔴🔴🔴"
            return 0  # Exit loop to proceed with transition logic
        fi

        # Still running - continue monitoring
        echo "📊 Phase integration still in progress (PID: $INTEGRATION_PID)..."
    done
}

# START MONITORING ON STATE ENTRY
echo "📊 Starting R358 phase integration monitoring loop..."
monitor_phase_integration_completion
```

### 2. 🔴🔴🔴 CHECK PHASE INTEGRATION REPORT AND ENFORCE R291 GATES 🔴🔴🔴
```bash
# MANDATORY: Check for phase integration report AND enforce build/test gates
PHASE=$(jq '.current_phase' orchestrator-state.json)
REPORT_FILE="efforts/phase${PHASE}/phase-integration/PHASE_INTEGRATION_REPORT.md"

if [ ! -f "$REPORT_FILE" ]; then
    echo "🔴 CRITICAL: No phase integration report found at $REPORT_FILE"
    # NO REPORT = IMMEDIATE ERROR_RECOVERY
    UPDATE_STATE="ERROR_RECOVERY"
    UPDATE_REASON="No phase integration report - R291 gates cannot be verified"
else
    echo "✅ Found phase integration report, enforcing R291 gates..."
    
    # Extract status from report
    INTEGRATION_STATUS=$(grep "^Phase Integration Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    BUILD_STATUS=$(grep "^Phase Build Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    TEST_STATUS=$(grep "^Phase Test Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    DEMO_STATUS=$(grep "^Phase Demo Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ' || echo "NOT_RUN")
    CONFLICTS=$(grep "^Unresolved Conflicts:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    
    echo "🔍 Phase Gate Status Check (R291 Enforcement):"
    echo "  Phase Build: $BUILD_STATUS"
    echo "  Phase Tests: $TEST_STATUS"
    echo "  Phase Demo: $DEMO_STATUS"
    echo "  Phase Integration: $INTEGRATION_STATUS"
    echo "  Conflicts: $CONFLICTS"
    
    # 🔴🔴🔴 R291 SUPREME GATE ENFORCEMENT FOR PHASE 🔴🔴🔴
    
    # BUILD GATE CHECK - PHASE LEVEL
    if [[ "$BUILD_STATUS" != "PASSING" ]] && [[ "$BUILD_STATUS" != "SUCCESS" ]]; then
        echo "🔴🔴🔴 PHASE BUILD GATE FAILED - MANDATORY ERROR_RECOVERY 🔴🔴🔴"
        echo "R291 VIOLATION: Phase build did not pass ($BUILD_STATUS)"
        echo "Cannot proceed to phase assessment without successful build!"
        UPDATE_STATE="ERROR_RECOVERY"
        UPDATE_REASON="R291 PHASE BUILD GATE FAILURE: $BUILD_STATUS"
        
    # TEST GATE CHECK - PHASE LEVEL
    elif [[ "$TEST_STATUS" != "PASSING" ]] && [[ "$TEST_STATUS" != "SUCCESS" ]]; then
        echo "🔴🔴🔴 PHASE TEST GATE FAILED - MANDATORY ERROR_RECOVERY 🔴🔴🔴"
        echo "R291 VIOLATION: Phase tests did not pass ($TEST_STATUS)"
        echo "Cannot proceed to phase assessment without passing tests!"
        UPDATE_STATE="ERROR_RECOVERY"
        UPDATE_REASON="R291 PHASE TEST GATE FAILURE: $TEST_STATUS"
        
    # DEMO GATE CHECK - PHASE LEVEL (if present)
    elif [[ "$DEMO_STATUS" != "NOT_RUN" ]] && \
         [[ "$DEMO_STATUS" != "PASSING" ]] && \
         [[ "$DEMO_STATUS" != "SUCCESS" ]]; then
        echo "🔴🔴🔴 PHASE DEMO GATE FAILED - MANDATORY ERROR_RECOVERY 🔴🔴🔴"
        echo "R291 VIOLATION: Phase demo did not pass ($DEMO_STATUS)"
        UPDATE_STATE="ERROR_RECOVERY"
        UPDATE_REASON="R291 PHASE DEMO GATE FAILURE: $DEMO_STATUS"
        
    # CONFLICTS CHECK - R321 ENFORCEMENT
    elif [[ "$CONFLICTS" != "0" ]] && [[ -n "$CONFLICTS" ]]; then
        echo "🔴🔴🔴 R321 ENFORCEMENT: Conflicts require immediate source fixes! 🔴🔴🔴"
        echo "Integration branches are READ-ONLY - fixes must go to source branches"
        UPDATE_STATE="IMMEDIATE_BACKPORT_REQUIRED"
        UPDATE_REASON="R321: Phase has $CONFLICTS conflicts - fix in source branches immediately"
        
    # INTEGRATION STATUS CHECK - R321 ENFORCEMENT
    elif [[ "$INTEGRATION_STATUS" != "SUCCESS" ]]; then
        echo "🔴🔴🔴 R321 ENFORCEMENT: Integration failure requires immediate source fixes! 🔴🔴🔴"
        echo "Cannot fix in integration branch - must fix in wave/effort branches"
        UPDATE_STATE="IMMEDIATE_BACKPORT_REQUIRED"
        UPDATE_REASON="R321: Phase integration failed ($INTEGRATION_STATUS) - fix sources immediately"
        
    # ALL GATES PASSED - PHASE NEEDS CODE REVIEW
    else
        echo "✅✅✅ ALL PHASE R291 GATES PASSED ✅✅✅"
        echo "  ✅ Phase Build: $BUILD_STATUS"
        echo "  ✅ Phase Tests: $TEST_STATUS"
        echo "  ✅ Phase Demo: $DEMO_STATUS (if run)"
        echo "  ✅ Phase Integration: $INTEGRATION_STATUS"
        echo "  ✅ Conflicts: $CONFLICTS"
        echo "Phase integration successful - needs code review!"
        UPDATE_STATE="PHASE_INTEGRATION_CODE_REVIEW"
        UPDATE_REASON="All phase gates passed - need code review of integrated phase"
    fi
fi

echo ""
echo "🎯 DECISION: Transitioning to $UPDATE_STATE"
echo "📝 REASON: $UPDATE_REASON"
```

### 3. Update State File with Phase Integration Results
```bash
# Update orchestrator state
jq ".current_state = \"$UPDATE_STATE\"" -i orchestrator-state.json
jq ".phase_integration_status.phase${PHASE} = \"$INTEGRATION_STATUS\"" -i orchestrator-state.json
jq ".state_transition_history += [{\"from\": \"MONITORING_PHASE_INTEGRATION\", \"to\": \"$UPDATE_STATE\", \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\", \"reason\": \"$UPDATE_REASON\"}]" -i orchestrator-state.json

# Commit state change
git add orchestrator-state.json
git commit -m "state: MONITORING_PHASE_INTEGRATION → $UPDATE_STATE - $UPDATE_REASON"
git push
```

## Valid Transitions

Based on phase integration report analysis:

1. **SUCCESS Path**: `MONITORING_PHASE_INTEGRATION` → `SPAWN_ARCHITECT_PHASE_ASSESSMENT`
   - When: Phase integration, build, tests all pass, no conflicts
   
2. **FAILURE Path (R321)**: `MONITORING_PHASE_INTEGRATION` → `IMMEDIATE_BACKPORT_REQUIRED`
   - When: Phase integration failed, conflicts exist, or build/tests fail
   - MUST fix in source branches (wave/effort branches) immediately
   
3. **ERROR Path**: `MONITORING_PHASE_INTEGRATION` → `ERROR_RECOVERY`
   - When: No report found or unexpected status

## 🔴🔴🔴 MANDATORY PHASE RE-INTEGRATION PROTOCOL (R321) 🔴🔴🔴

**When phase integration fails, the cycle MUST be:**

```
MONITORING_PHASE_INTEGRATION (detects failure)
    ↓
IMMEDIATE_BACKPORT_REQUIRED (R321: fix in source branches)
    ↓
SPAWN_ENGINEERS_FOR_FIXES (fix wave/effort branches)
    ↓
MONITORING_FIX_PROGRESS (monitor source fixes)
    ↓
SPAWN_CODE_REVIEWERS_FOR_REVIEW (review fixes)
    ↓
MONITOR_REVIEWS (all fixes reviewed)
    ↓
PHASE_INTEGRATION (DELETE old, create FRESH integration)
    ↓
SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN (new merge plan)
    ↓
SPAWN_INTEGRATION_AGENT_PHASE (re-merge ALL waves with fixes)
    ↓
MONITORING_PHASE_INTEGRATION (check if NOW it works)
```

**CRITICAL POINTS:**
- Integration branches are READ-ONLY (R321)
- ALL fixes go to source branches FIRST
- Phase integration must be DELETED and RE-CREATED
- Re-merge ALL waves to get fixed code

## Grading Criteria

- ✅ **+20%**: Properly check for phase integration completion (R358)
- ✅ **+20%**: Correctly parse phase-level status fields
- ✅ **+20%**: Check for unresolved conflicts
- ✅ **+20%**: Transition to IMMEDIATE_BACKPORT_REQUIRED on failures (R321)
- ✅ **+20%**: Detect completion within 30 seconds and transition immediately

## Common Violations

- ❌ **-100%**: Remaining in MONITORING_PHASE_INTEGRATION after completion (R358 violation)
- ❌ **-50%**: Not checking for phase integration completion continuously
- ❌ **-50%**: Proceeding to phase assessment with conflicts
- ❌ **-30%**: Manual checking instead of automated loop
- ❌ **-30%**: Not updating phase integration status

## Related Rules

- R358: Integration Completion Detection (SUPREME - must detect and transition immediately)
- R321: Immediate Backport During Integration (for failures)
- R282: Phase Integration Protocol
- R259: Mandatory Phase Integration After Fixes
- R260: Integration Agent Core Requirements
- R263: Integration Documentation Requirements
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
