# ORCHESTRATOR STATE: CREATE_PHASE_INTEGRATION_BRANCH_EARLY

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR CREATE_PHASE_INTEGRATION_BRANCH_EARLY STATE

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

5. **🔴🔴🔴 R308** - Incremental Branching Strategy
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R308-incremental-branching-strategy.md`
   - Criticality: SUPREME LAW - Phase 2 NEVER from main!
   - Summary: Integration branches must follow incremental strategy

6. **🔴🔴🔴 R009** - Mandatory Wave/Phase Integration Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R009-integration-branch-creation.md`
   - Criticality: SUPREME LAW - Wave N+1 REQUIRES Wave N integration
   - Summary: Enforces trunk-based development

## Overview
This state creates the phase-N-integration branch immediately after phase test planning, implementing R342 for phase-level test storage.

## Entry Criteria
- Transitioned from WAITING_FOR_PHASE_TEST_PLAN
- PHASE-N-TEST-PLAN.md exists
- Phase tests have been created

## State Responsibilities

### 1. Determine Current Phase
```bash
# Get current phase number
PHASE=$(yq '.current_phase' orchestrator-state.json)
echo "📍 Creating integration branch for Phase $PHASE"
```

### 2. Create Phase Integration Workspace
```bash
# R104: Proper workspace isolation
PHASE_INT_DIR="/efforts/phase${PHASE}/integration-workspace"
mkdir -p "$PHASE_INT_DIR"
cd "$PHASE_INT_DIR"
echo "📁 Created phase $PHASE integration workspace"
```

### 3. Determine Base Branch (R308 Incremental)
```bash
# R308: Build on previous integration branch
if [ "$PHASE" -eq 1 ]; then
    BASE_BRANCH="project-integration"  # Build on project tests
    echo "🔗 Using project-integration as base (Phase 1)"
else
    PREV_PHASE=$((PHASE - 1))
    BASE_BRANCH="phase-${PREV_PHASE}-integration"
    echo "🔗 Using phase-${PREV_PHASE}-integration as base"
fi
```

### 4. Clone and Create Branch
```bash
# Clone target repository
TARGET_REPO=$(yq '.target_repository' "$CLAUDE_PROJECT_DIR/orchestrator-state.json")
git clone "$TARGET_REPO" target-repo
cd target-repo

# Create branch from appropriate base
git checkout -b "phase-${PHASE}-integration" "$BASE_BRANCH"
echo "🌿 Created phase-${PHASE}-integration branch from $BASE_BRANCH"
```

### 5. Store Phase Tests (R342 Core)
```bash
# Copy phase tests to integration branch
mkdir -p "tests/phase${PHASE}"
if [ -d "$CLAUDE_PROJECT_DIR/phase-tests/phase-${PHASE}" ]; then
    cp -r "$CLAUDE_PROJECT_DIR/phase-tests/phase-${PHASE}/"* "tests/phase${PHASE}/"
    echo "📋 Copied phase $PHASE tests to branch"
fi

# Copy test plan and harness
cp "$CLAUDE_PROJECT_DIR/PHASE-${PHASE}-TEST-PLAN.md" .
cp "$CLAUDE_PROJECT_DIR/PHASE-${PHASE}-TEST-HARNESS.sh" "tests/phase${PHASE}/" 2>/dev/null || true

# Verify accumulated tests (R308)
echo "📊 Accumulated tests in branch:"
find tests/ -type f -name "*.test.*" | wc -l
```

### 6. Commit and Push
```bash
# Add and commit phase tests
git add tests/ "PHASE-${PHASE}-TEST-PLAN.md"
git commit -m "test: add phase ${PHASE} tests (R341 TDD, R342 early branch)"
git push -u origin "phase-${PHASE}-integration"

echo "✅ R342 COMPLIANT: Phase ${PHASE} tests stored in integration branch"
```

### 7. Update State File
```bash
# Track phase integration branch
yq -i ".phase_${PHASE}_integration.branch = \"phase-${PHASE}-integration\"" orchestrator-state.json
yq -i ".phase_${PHASE}_integration.base_branch = \"$BASE_BRANCH\"" orchestrator-state.json
yq -i ".phase_${PHASE}_integration.created_at = \"$(date -Iseconds)\"" orchestrator-state.json
yq -i ".phase_${PHASE}_integration.has_tests = true" orchestrator-state.json
yq -i ".phase_${PHASE}_integration.workspace = \"$PHASE_INT_DIR\"" orchestrator-state.json
```

### 8. Transition to Implementation Planning
```bash
update_state "SPAWN_CODE_REVIEWER_PHASE_IMPL"
echo "📍 Ready for phase implementation planning with tests in place"
save_todos "Created phase ${PHASE} integration branch with tests"
git add orchestrator-state.json todos/
git commit -m "state: created phase-${PHASE}-integration branch per R342"
git push
```

## Exit Criteria
- phase-N-integration branch created
- Tests committed to branch
- Branch pushed to remote
- State transitions to SPAWN_CODE_REVIEWER_PHASE_IMPL

## Success Metrics
- ✅ Integration branch exists with correct base
- ✅ Tests stored in tests/phaseN/
- ✅ R308 incremental branching maintained
- ✅ R342 compliance achieved
- ✅ Ready for implementation planning

## Related Rules
- R342: Early integration branch creation (SUPREME LAW)
- R341: TDD - tests before implementation
- R308: Incremental branching strategy
- R104: Target repository isolation