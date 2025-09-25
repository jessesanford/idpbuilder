# INTEGRATION_CODE_REVIEW State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR INTEGRATION_CODE_REVIEW STATE

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
Spawn Code Reviewer to perform quality review of integrated code after wave integration is complete. This state ensures the merged code maintains quality standards and catches integration-introduced issues before architect review.

## Entry Criteria
- **From**: MONITORING_INTEGRATION
- **Condition**: Integration Agent has completed wave integration successfully
- **Required**: 
  - Integration workspace exists with merged code
  - All effort branches merged into wave integration branch
  - Basic merge conflicts resolved by Integration Agent

## State Actions

### 1. IMMEDIATE: Spawn Code Reviewer for Integration Review
```bash
# Spawn Code Reviewer to review integrated code
/spawn agent-code-reviewer INTEGRATION_CODE_REVIEW \
  --wave "W${current_wave}" \
  --branch "${wave_integration_branch}" \
  --focus "integration-quality"
```

### 2. Code Reviewer Responsibilities
The spawned Code Reviewer will:
- Review the quality of integrated code
- Check for integration-introduced bugs
- Validate merge conflict resolutions
- Ensure code consistency across merged efforts
- Verify tests still pass after integration
- Check for duplicate code or functionality conflicts
- **Verify demo scripts are present and executable (R330/R291)**
- **Check demo coverage of integrated features**
- **Validate demo documentation completeness**
- **Review demo execution results from integration testing**
- Create INTEGRATION_CODE_REVIEW_REPORT.md

### 3. Update State File
```json
{
  "current_state": "INTEGRATION_CODE_REVIEW",
  "phase": "integration",
  "review_status": "code_review_in_progress",
  "integration_review": {
    "reviewer": "agent-code-reviewer",
    "branch": "${wave_integration_branch}",
    "focus": "integration_quality",
    "started_at": "timestamp"
  }
}
```

## Exit Criteria

### Success Path → WAITING_FOR_INTEGRATION_CODE_REVIEW
- Code Reviewer spawned successfully
- State file updated with review details
- Transition to waiting state for results

### Failure Scenarios
- **Spawn Failure** → ERROR_RECOVERY
- **Invalid Integration State** → ERROR_RECOVERY

## Key Differences from Other Reviews
- **Focus**: Integration quality, not individual effort quality
- **Scope**: Merged code interactions and conflicts
- **Timing**: After integration, before architect review
- **Purpose**: Catch integration-specific issues

## Rules Enforced
- R233: Immediate action upon state entry
- R313: Stop after spawning agent
- R238: Monitor for review reports
- R321: Any fixes require immediate backport

## Report Expected
The Code Reviewer will create:
- `INTEGRATION_CODE_REVIEW_REPORT.md` with:
  - Integration quality assessment
  - Merge conflict resolution review
  - Cross-effort consistency check
  - Test status after integration
  - **Demo verification results (R330/R291)**
    - Demo scripts present: YES/NO
    - Demo scripts executable: YES/NO
    - Demo coverage adequate: YES/NO
    - Demo documentation complete: YES/NO
    - Demo execution results: PASSED/FAILED
  - Issues found (if any)
  - Recommendation: PASS/FAIL

### Demo Verification Checklist (R330/R291)
```markdown
## Demo Verification (R330/R291 Compliance)
- [ ] All effort demo scripts present per R330 plans
- [ ] Wave integration demo script exists
- [ ] Demo documentation (DEMO.md files) present
- [ ] Demo scripts are executable (chmod +x)
- [ ] Demo test data files available
- [ ] Demo execution logs captured in demo-results/
- [ ] Demo failures documented with root cause
- [ ] Demo coverage matches integrated features
```

## Transition Rules
- **ALWAYS** → WAITING_FOR_INTEGRATION_CODE_REVIEW (after spawn)
- **NEVER** skip directly to WAVE_REVIEW
- **NEVER** proceed without Code Reviewer spawn