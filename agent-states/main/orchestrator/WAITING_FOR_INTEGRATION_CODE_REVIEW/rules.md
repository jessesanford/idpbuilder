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
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-requirements.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state.json before EVERY state transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322 Part A-mandatory-stop-after-spawn.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

5. **🔴🔴🔴 R304** - Mandatory Line Counting Tool
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counting-tool.md`
   - Criticality: SUPREME LAW - Line counting requirements
   - Summary: MUST use tools/line-counter.sh for ALL measurements

## State Purpose
Monitor Code Reviewer progress on integration code review and process results when complete.

## Entry Criteria
- **From**: INTEGRATION_CODE_REVIEW
- **Condition**: Code Reviewer spawned for integration review
- **Required**: State file shows review in progress

## State Actions

### 1. IMMEDIATE: Check for Review Completion
```bash
# Check for integration code review report
if [ -f "INTEGRATION_CODE_REVIEW_REPORT.md" ]; then
    echo "Integration code review complete"
    process_review_results
else
    echo "Waiting for integration code review to complete"
    exit 0  # Will be re-invoked later
fi
```

### 2. Process Review Results
When report exists:
- Parse INTEGRATION_CODE_REVIEW_REPORT.md
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
    "report": "INTEGRATION_CODE_REVIEW_REPORT.md"
  }
}
```

## Exit Criteria

### Success Path (PASS) → WAVE_REVIEW
- Review passed with no critical issues
- Integration quality verified
- Ready for architect review

### Failure Path (FAIL) → SPAWN_CODE_REVIEWER_INTEGRATION_FIX_PLAN
- Critical integration issues found
- Fixes required before proceeding
- Need plan to address issues

### Waiting Path → WAITING_FOR_INTEGRATION_CODE_REVIEW
- Review still in progress
- Exit and wait for next check

## Report Processing
Parse INTEGRATION_CODE_REVIEW_REPORT.md for:
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
CASCADE_MODE=$(jq -r '.cascade_coordination.cascade_mode // false' orchestrator-state.json)

if [[ "$CASCADE_MODE" == "true" ]]; then
    echo "🔴🔴🔴 CASCADE MODE ACTIVE 🔴🔴🔴"
    
    # In cascade mode, different transitions apply
    if [[ "$REVIEW_RESULT" == "PASS" ]]; then
        echo "Integration code review passed in cascade mode"
        echo "Returning to CASCADE_REINTEGRATION to continue cascade..."
        NEXT_STATE="CASCADE_REINTEGRATION"
    elif [[ "$REVIEW_RESULT" == "FAIL" ]]; then
        echo "Integration code review failed in cascade mode"
        echo "Must fix issues before continuing cascade"
        NEXT_STATE="SPAWN_CODE_REVIEWER_INTEGRATION_FIX_PLAN"
    else
        echo "Review still pending"
        NEXT_STATE="WAITING_FOR_INTEGRATION_CODE_REVIEW"
    fi
else
    # Normal flow - not in cascade mode
    if [[ "$REVIEW_RESULT" == "PASS" ]]; then
        NEXT_STATE="WAVE_REVIEW"
    elif [[ "$REVIEW_RESULT" == "FAIL" ]]; then
        NEXT_STATE="SPAWN_CODE_REVIEWER_INTEGRATION_FIX_PLAN"
    else
        NEXT_STATE="WAITING_FOR_INTEGRATION_CODE_REVIEW"
    fi
fi
```

### Standard Transitions (Non-Cascade Mode)
- **If PASS** → WAVE_REVIEW
- **If FAIL** → SPAWN_CODE_REVIEWER_INTEGRATION_FIX_PLAN
- **If pending** → WAITING_FOR_INTEGRATION_CODE_REVIEW (self)

### Cascade Mode Transitions
- **If PASS** → CASCADE_REINTEGRATION (continue cascade)
- **If FAIL** → SPAWN_CODE_REVIEWER_INTEGRATION_FIX_PLAN
- **If pending** → WAITING_FOR_INTEGRATION_CODE_REVIEW (self)