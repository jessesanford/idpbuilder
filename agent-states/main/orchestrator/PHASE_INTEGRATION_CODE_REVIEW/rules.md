# PHASE_INTEGRATION_CODE_REVIEW State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR PHASE_INTEGRATION_CODE_REVIEW STATE

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
Spawn Code Reviewer to perform quality review of integrated code after phase integration is complete. This state ensures the merged phase maintains quality standards and catches phase-level integration issues before architect assessment.

## Entry Criteria
- **From**: MONITORING_PHASE_INTEGRATION
- **Condition**: Integration Agent has completed phase integration successfully
- **Required**: 
  - Phase integration workspace exists with merged waves
  - All wave integration branches merged into phase branch
  - Phase-level merge conflicts resolved

## State Actions

### 1. IMMEDIATE: Spawn Code Reviewer for Phase Integration Review
```bash
# Spawn Code Reviewer to review phase integrated code
/spawn agent-code-reviewer PHASE_INTEGRATION_CODE_REVIEW \
  --phase "P${current_phase}" \
  --branch "${phase_integration_branch}" \
  --focus "phase-integration-quality"
```

### 2. Code Reviewer Responsibilities
The spawned Code Reviewer will:
- Review quality of phase-integrated code
- Check for cross-wave integration issues
- Validate wave merge conflict resolutions
- Ensure architectural consistency across waves
- Verify all tests pass at phase level
- Check for feature completeness per phase plan
- Identify any missing integrations
- **Verify phase-level demo orchestration (R330/R291)**
- **Check all wave demos integrate properly**
- **Validate phase demo comprehensiveness**
- **Review phase demo execution results**
- Create PHASE_INTEGRATION_CODE_REVIEW_REPORT.md

### 3. Update State File
```json
{
  "current_state": "PHASE_INTEGRATION_CODE_REVIEW",
  "phase": "integration",
  "review_status": "phase_code_review_in_progress",
  "phase_integration_review": {
    "reviewer": "agent-code-reviewer",
    "branch": "${phase_integration_branch}",
    "focus": "phase_integration_quality",
    "waves_integrated": ["W1", "W2", "W3"],
    "started_at": "timestamp"
  }
}
```

## Exit Criteria

### Success Path → WAITING_FOR_PHASE_INTEGRATION_CODE_REVIEW
- Code Reviewer spawned successfully
- State file updated with phase review details
- Transition to waiting state for results

### Failure Scenarios
- **Spawn Failure** → ERROR_RECOVERY
- **Invalid Phase State** → ERROR_RECOVERY

## Key Differences from Wave Integration Review
- **Scope**: Entire phase (multiple waves)
- **Focus**: Cross-wave consistency and completeness
- **Complexity**: Higher-level architectural alignment
- **Validation**: Phase requirements met

## Rules Enforced
- R233: Immediate action upon state entry
- R313: Stop after spawning agent
- R238: Monitor for review reports
- R285: Phase must include all waves
- R321: Any fixes require immediate backport

## Report Expected
The Code Reviewer will create:
- `PHASE_INTEGRATION_CODE_REVIEW_REPORT.md` with:
  - Phase integration quality assessment
  - Cross-wave consistency verification
  - Feature completeness check
  - Architectural alignment review
  - Test coverage at phase level
  - Performance impact assessment
  - Issues found (if any)
  - Recommendation: PASS/FAIL

## Transition Rules
- **ALWAYS** → WAITING_FOR_PHASE_INTEGRATION_CODE_REVIEW (after spawn)
- **NEVER** skip directly to SPAWN_ARCHITECT_PHASE_ASSESSMENT
- **NEVER** proceed without Code Reviewer spawn