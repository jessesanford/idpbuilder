# ORCHESTRATOR STATE: CREATE_PROJECT_INTEGRATION_BRANCH_EARLY

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR CREATE_PROJECT_INTEGRATION_BRANCH_EARLY STATE

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
This state creates the project-integration branch immediately after test planning, storing tests in their permanent location per R342.

## Entry Criteria
- Transitioned from WAITING_FOR_PROJECT_TEST_PLAN
- PROJECT-TEST-PLAN.md exists
- Project tests have been created

## State Responsibilities

### 1. Create Project Integration Workspace
```bash
# R104: Proper workspace isolation
PROJECT_INT_DIR="/efforts/project/integration-workspace"
mkdir -p "$PROJECT_INT_DIR"
cd "$PROJECT_INT_DIR"
echo "📁 Created project integration workspace"
```

### 2. Clone Target Repository
```bash
# Get target repository from config
TARGET_REPO=$(yq '.target_repository' "$CLAUDE_PROJECT_DIR/orchestrator-state.json")
if [ -z "$TARGET_REPO" ]; then
    echo "❌ No target repository configured!"
    exit 1
fi

git clone "$TARGET_REPO" target-repo
cd target-repo
echo "✅ Cloned target repository"
```

### 3. Create Project Integration Branch
```bash
# Create branch from main
git checkout -b project-integration
echo "🌿 Created project-integration branch"
```

### 4. Store Project Tests (R342 Core Requirement)
```bash
# Copy tests to integration branch
mkdir -p tests/project
if [ -d "$CLAUDE_PROJECT_DIR/project-tests" ]; then
    cp -r "$CLAUDE_PROJECT_DIR/project-tests/"* tests/project/
    echo "📋 Copied project tests to branch"
fi

# Copy test plan and harness
cp "$CLAUDE_PROJECT_DIR/PROJECT-TEST-PLAN.md" .
cp "$CLAUDE_PROJECT_DIR/PROJECT-TEST-HARNESS.sh" tests/project/ 2>/dev/null || true

# Verify tests exist
ls -la tests/project/
```

### 5. Commit and Push Tests
```bash
# Add and commit tests
git add tests/ PROJECT-TEST-PLAN.md
git commit -m "test: add project-level tests (R341 TDD, R342 early branch)"
git push -u origin project-integration

echo "✅ R342 COMPLIANT: Tests stored in project-integration branch"
```

### 6. Update State File with Branch Info
```bash
# Track integration branch creation
yq -i '.project_integration.branch = "project-integration"' orchestrator-state.json
yq -i '.project_integration.created_at = "'$(date -Iseconds)'"' orchestrator-state.json
yq -i '.project_integration.has_tests = true' orchestrator-state.json
yq -i '.project_integration.workspace = "'$PROJECT_INT_DIR'"' orchestrator-state.json
```

### 7. Transition Back to INIT for Phase 1
```bash
update_state "INIT"
echo "📍 Ready to begin Phase 1 with project tests in place"
save_todos "Created project integration branch with tests"
git add orchestrator-state.json todos/
git commit -m "state: created project-integration branch per R342"
git push
```

## Exit Criteria
- project-integration branch created
- Tests committed to branch
- Branch pushed to remote
- State transitions to INIT for Phase 1 start

## Success Metrics
- ✅ Integration branch exists
- ✅ Tests stored in tests/project/
- ✅ Branch pushed to remote
- ✅ R342 compliance achieved
- ✅ Ready for Phase 1

## Related Rules
- R342: Early integration branch creation (SUPREME LAW)
- R341: TDD - tests before implementation
- R104: Target repository isolation
- R308: Incremental branching strategy