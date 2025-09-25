# WAITING_FOR_PROJECT_INTEGRATION_CODE_REVIEW State Rules


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_PROJECT_INTEGRATION_CODE_REVIEW STATE

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
   - Summary: WAITING states require active monitoring, not passive waiting

## State Purpose
Monitor Code Reviewer progress on project integration code review and process comprehensive results when complete.

## Entry Criteria
- **From**: PROJECT_INTEGRATION_CODE_REVIEW
- **Condition**: Code Reviewer spawned for project integration review
- **Required**: State file shows project review in progress

## State Actions

### 1. IMMEDIATE: Check for Project Review Completion
```bash
# Check for project integration code review report
if [ -f "PROJECT_INTEGRATION_CODE_REVIEW_REPORT.md" ]; then
    echo "Project integration code review complete"
    process_project_review_results
else
    echo "Waiting for project integration code review to complete"
    exit 0  # Will be re-invoked later
fi
```

### 2. Process Project Review Results
When report exists:
- Parse PROJECT_INTEGRATION_CODE_REVIEW_REPORT.md
- Extract PASS/FAIL/CONDITIONAL_PASS status
- Analyze critical issues across project
- Verify all requirements met
- Assess production readiness
- Determine next state based on comprehensive results

### 3. Update State File
```json
{
  "current_state": "WAITING_FOR_PROJECT_INTEGRATION_CODE_REVIEW",
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
    "report": "PROJECT_INTEGRATION_CODE_REVIEW_REPORT.md"
  }
}
```

## Exit Criteria

### Success Path (PASS) → SPAWN_CODE_REVIEWER_PROJECT_VALIDATION
- Review passed with no critical issues
- Project quality verified
- All requirements met
- Ready for final validation

### Conditional Pass → SPAWN_CODE_REVIEWER_PROJECT_VALIDATION
- Minor issues identified but not blocking
- Document known issues for tracking
- Proceed with validation

### Failure Path (FAIL) → SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING
- Critical project-level issues found
- Requirements not met
- Need comprehensive fix plan
- May require multiple phases of fixes

### Waiting Path → WAITING_FOR_PROJECT_INTEGRATION_CODE_REVIEW
- Review still in progress
- Exit and wait for next check

## Report Processing
Parse PROJECT_INTEGRATION_CODE_REVIEW_REPORT.md for:
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
- **If PASS** → SPAWN_CODE_REVIEWER_PROJECT_VALIDATION
- **If CONDITIONAL_PASS** → SPAWN_CODE_REVIEWER_PROJECT_VALIDATION (with notes)
- **If FAIL** → SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING
- **If pending** → WAITING_FOR_PROJECT_INTEGRATION_CODE_REVIEW (self)

## Special Handling
- This is the most critical review gate
- May require escalation for CONDITIONAL_PASS
- Document all issues for future reference
- Consider rollback if severe issues found