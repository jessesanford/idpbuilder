# ORCHESTRATOR STATE: CREATE_WAVE_INTEGRATION_BRANCH_EARLY

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR CREATE_WAVE_INTEGRATION_BRANCH_EARLY STATE

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
This state creates the wave integration branch immediately after wave test planning, implementing R342 for wave-level test storage.

## Entry Criteria
- Transitioned from WAITING_FOR_WAVE_TEST_PLAN
- WAVE-M-TEST-PLAN.md exists
- Wave tests have been created

## State Responsibilities

### 1. Determine Current Phase and Wave
```bash
# Get current phase and wave numbers
PHASE=$(yq '.current_phase' orchestrator-state.json)
WAVE=$(yq '.current_wave' orchestrator-state.json)
echo "📍 Creating integration branch for Phase $PHASE Wave $WAVE"
```

### 2. Create Wave Integration Workspace
```bash
# R104: Proper workspace isolation
WAVE_INT_DIR="/efforts/phase${PHASE}/wave${WAVE}/integration-workspace"
mkdir -p "$WAVE_INT_DIR"
cd "$WAVE_INT_DIR"
echo "📁 Created phase $PHASE wave $WAVE integration workspace"
```

### 3. Determine Base Branch (R336 + R308)
```bash
# R336: Each wave builds on previous wave's integration
# R308: Incremental branching strategy
if [ "$WAVE" -eq 1 ]; then
    BASE_BRANCH="phase-${PHASE}-integration"  # First wave uses phase base
    echo "🔗 Using phase-${PHASE}-integration as base (Wave 1)"
else
    PREV_WAVE=$((WAVE - 1))
    BASE_BRANCH="phase-${PHASE}-wave-${PREV_WAVE}-integration"
    echo "🔗 Using phase-${PHASE}-wave-${PREV_WAVE}-integration as base (R336)"
fi
```

### 4. Clone and Create Branch
```bash
# Clone target repository
TARGET_REPO=$(yq '.target_repository' "$CLAUDE_PROJECT_DIR/orchestrator-state.json")
git clone "$TARGET_REPO" target-repo
cd target-repo

# Create branch from appropriate base
git checkout -b "phase-${PHASE}-wave-${WAVE}-integration" "$BASE_BRANCH"
echo "🌿 Created phase-${PHASE}-wave-${WAVE}-integration from $BASE_BRANCH"
```

### 5. Store Wave Tests (R342 Core)
```bash
# Copy wave tests to integration branch
mkdir -p "tests/phase${PHASE}/wave${WAVE}"
if [ -d "$CLAUDE_PROJECT_DIR/wave-tests/phase-${PHASE}/wave-${WAVE}" ]; then
    cp -r "$CLAUDE_PROJECT_DIR/wave-tests/phase-${PHASE}/wave-${WAVE}/"* \
          "tests/phase${PHASE}/wave${WAVE}/"
    echo "📋 Copied wave $WAVE tests to branch"
fi

# Copy test plan and harness
cp "$CLAUDE_PROJECT_DIR/WAVE-${WAVE}-TEST-PLAN.md" .
cp "$CLAUDE_PROJECT_DIR/WAVE-${WAVE}-TEST-HARNESS.sh" \
   "tests/phase${PHASE}/wave${WAVE}/" 2>/dev/null || true

# Verify accumulated tests (R308 + R336)
echo "📊 Accumulated tests in branch:"
echo "  - Project tests: $(find tests/project -type f 2>/dev/null | wc -l)"
echo "  - Phase tests: $(find tests/phase${PHASE} -type f | wc -l)"
echo "  - Total: $(find tests/ -type f -name "*.test.*" | wc -l)"
```

### 6. Commit and Push
```bash
# Add and commit wave tests
git add tests/ "WAVE-${WAVE}-TEST-PLAN.md"
git commit -m "test: add phase ${PHASE} wave ${WAVE} tests (R341 TDD, R342 early branch)"
git push -u origin "phase-${PHASE}-wave-${WAVE}-integration"

echo "✅ R342 COMPLIANT: Wave ${WAVE} tests stored in integration branch"
```

### 7. Update State File
```bash
# Track wave integration branch
yq -i ".current_wave_integration.branch = \"phase-${PHASE}-wave-${WAVE}-integration\"" orchestrator-state.json
yq -i ".current_wave_integration.base_branch = \"$BASE_BRANCH\"" orchestrator-state.json
yq -i ".current_wave_integration.created_at = \"$(date -Iseconds)\"" orchestrator-state.json
yq -i ".current_wave_integration.has_tests = true" orchestrator-state.json
yq -i ".current_wave_integration.workspace = \"$WAVE_INT_DIR\"" orchestrator-state.json

# Track in wave-specific section too
yq -i ".waves.wave${WAVE}.integration_branch = \"phase-${PHASE}-wave-${WAVE}-integration\"" orchestrator-state.json
yq -i ".waves.wave${WAVE}.tests_stored = true" orchestrator-state.json
```

### 8. Transition to Implementation Planning
```bash
update_state "SPAWN_CODE_REVIEWER_WAVE_IMPL"
echo "📍 Ready for wave implementation planning with tests in place"
save_todos "Created phase ${PHASE} wave ${WAVE} integration branch with tests"
git add orchestrator-state.json todos/
git commit -m "state: created wave ${WAVE} integration branch per R342"
git push
```

## Exit Criteria
- phase-N-wave-M-integration branch created
- Tests committed to branch
- Branch pushed to remote
- State transitions to SPAWN_CODE_REVIEWER_WAVE_IMPL

## Success Metrics
- ✅ Integration branch exists with correct base
- ✅ Tests stored in tests/phaseN/waveM/
- ✅ R336 wave chaining maintained
- ✅ R308 incremental branching maintained
- ✅ R342 compliance achieved
- ✅ Ready for implementation planning

## Related Rules
- R342: Early integration branch creation (SUPREME LAW)
- R341: TDD - tests before implementation
- R336: Mandatory wave integration before next wave
- R308: Incremental branching strategy
- R104: Target repository isolation