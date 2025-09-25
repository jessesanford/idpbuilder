# ORCHESTRATOR STATE: WAITING_FOR_PROJECT_TEST_PLAN


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_PROJECT_TEST_PLAN STATE

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

## Overview
This state waits for the Code Reviewer to complete project-level test planning, then enforces R342 by transitioning to early branch creation.

## Entry Criteria
- Transitioned from SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING
- Code Reviewer spawned for project test planning

## State Responsibilities

### 1. Check for Test Plan Completion
```bash
# Look for completed project test plan
if [ -f "PROJECT-TEST-PLAN.md" ]; then
    echo "✅ Project test plan completed"
    TEST_READY=true
else
    echo "⏳ Waiting for project test plan..."
    TEST_READY=false
fi
```

### 2. Verify Test Artifacts
When tests are ready, verify:
```bash
# Check for required test artifacts
[ -f "PROJECT-TEST-PLAN.md" ] || { echo "❌ Missing test plan"; exit 1; }
[ -f "PROJECT-TEST-HARNESS.sh" ] || { echo "⚠️ Missing test harness"; }
[ -d "project-tests/" ] || { echo "⚠️ Missing test directory"; }
```

### 3. Update Planning File Tracking (R340)
```bash
# Record test plan location and metadata
yq -i '.planning_files.project_test_plan = "PROJECT-TEST-PLAN.md"' orchestrator-state.json
yq -i '.planning_files.project_test_author = "code-reviewer"' orchestrator-state.json
yq -i '.planning_files.project_test_timestamp = "'$(date -Iseconds)'"' orchestrator-state.json
yq -i '.planning_files.project_test_harness = "PROJECT-TEST-HARNESS.sh"' orchestrator-state.json
```

### 4. Enforce R342 - Transition to Early Branch Creation
```bash
if [ "$TEST_READY" = true ]; then
    echo "🔴🔴🔴 R342 ENFORCEMENT: Creating integration branch for test storage"
    update_state "CREATE_PROJECT_INTEGRATION_BRANCH_EARLY"
    echo "📍 Tests will be stored in project-integration branch"
fi
```

## Exit Criteria
- Project test plan exists
- Test artifacts verified
- R340 metadata tracking complete
- State transitions to CREATE_PROJECT_INTEGRATION_BRANCH_EARLY (R342)

## Success Metrics
- ✅ PROJECT-TEST-PLAN.md created
- ✅ Test harness present
- ✅ R340 planning file tracking complete
- ✅ R342 early branch creation enforced

## Related Rules
- R340: Planning file metadata tracking
- R341: TDD - tests created before implementation
- R342: Early integration branch creation (SUPREME LAW)
- R287: TODO persistence