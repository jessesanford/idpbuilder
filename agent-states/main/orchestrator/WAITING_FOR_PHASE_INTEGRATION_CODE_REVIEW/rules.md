# WAITING_FOR_PHASE_INTEGRATION_CODE_REVIEW State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_PHASE_INTEGRATION_CODE_REVIEW STATE

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
Monitor Code Reviewer progress on phase integration code review and process results when complete.

## Entry Criteria
- **From**: PHASE_INTEGRATION_CODE_REVIEW
- **Condition**: Code Reviewer spawned for phase integration review
- **Required**: State file shows phase review in progress

## State Actions

### 1. IMMEDIATE: Check for Phase Review Completion
```bash
# Check for phase integration code review report
if [ -f "PHASE_INTEGRATION_CODE_REVIEW_REPORT.md" ]; then
    echo "Phase integration code review complete"
    process_phase_review_results
else
    echo "Waiting for phase integration code review to complete"
    exit 0  # Will be re-invoked later
fi
```

### 2. Process Phase Review Results
When report exists:
- Parse PHASE_INTEGRATION_CODE_REVIEW_REPORT.md
- Extract PASS/FAIL status
- Identify phase-level integration issues
- Check feature completeness
- Determine next state based on results

### 3. Update State File
```json
{
  "current_state": "WAITING_FOR_PHASE_INTEGRATION_CODE_REVIEW",
  "phase_integration_review": {
    "status": "complete",
    "result": "PASS|FAIL",
    "phase_issues_found": [],
    "feature_complete": true,
    "cross_wave_conflicts": [],
    "report": "PHASE_INTEGRATION_CODE_REVIEW_REPORT.md"
  }
}
```

## Exit Criteria

### Success Path (PASS) → SPAWN_ARCHITECT_PHASE_ASSESSMENT
- Review passed with no critical issues
- Phase integration quality verified
- Ready for architect assessment

### Failure Path (FAIL) → SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN
- Critical phase integration issues found
- Cross-wave conflicts detected
- Need comprehensive fix plan

### Waiting Path → WAITING_FOR_PHASE_INTEGRATION_CODE_REVIEW
- Review still in progress
- Exit and wait for next check

## Report Processing
Parse PHASE_INTEGRATION_CODE_REVIEW_REPORT.md for:
- Overall phase status (PASS/FAIL)
- Critical cross-wave issues
- Feature completeness assessment
- Architectural violations
- Performance regressions
- Test coverage gaps

## Rules Enforced
- R233: Immediate check upon entry
- R238: Monitor for completion
- R285: Phase completeness validation
- R321: Fixes require backport planning

## Transition Rules
- **If PASS** → SPAWN_ARCHITECT_PHASE_ASSESSMENT
- **If FAIL** → SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN
- **If pending** → WAITING_FOR_PHASE_INTEGRATION_CODE_REVIEW (self)