# STATE MACHINE UPDATES FOR TEST-DRIVEN DEVELOPMENT

## Overview

This document specifies the exact state machine modifications needed to implement the project-level testing and early integration branch creation strategy.

## NEW STATES TO ADD

### 1. SPAWN_ARCHITECT_MASTER_PLANNING (New)
```yaml
state: SPAWN_ARCHITECT_MASTER_PLANNING
agent: orchestrator
description: Spawn Architect to create master architecture
valid_from:
  - INIT (when no master plan exists)
valid_to:
  - WAITING_FOR_MASTER_ARCHITECTURE
```

### 2. WAITING_FOR_MASTER_ARCHITECTURE (New)
```yaml
state: WAITING_FOR_MASTER_ARCHITECTURE
agent: orchestrator
description: Wait for Architect to complete master architecture
valid_from:
  - SPAWN_ARCHITECT_MASTER_PLANNING
valid_to:
  - SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING
```

### 3. SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING (New)
```yaml
state: SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING
agent: orchestrator
description: Spawn Code Reviewer to create project-level tests
valid_from:
  - WAITING_FOR_MASTER_ARCHITECTURE
valid_to:
  - WAITING_FOR_PROJECT_TEST_PLAN
```

### 4. WAITING_FOR_PROJECT_TEST_PLAN (New)
```yaml
state: WAITING_FOR_PROJECT_TEST_PLAN
agent: orchestrator
description: Wait for Code Reviewer to complete project tests
valid_from:
  - SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING
valid_to:
  - CREATE_PROJECT_INTEGRATION_BRANCH_EARLY
```

### 5. CREATE_PROJECT_INTEGRATION_BRANCH_EARLY (New)
```yaml
state: CREATE_PROJECT_INTEGRATION_BRANCH_EARLY
agent: orchestrator
description: Create project-integration branch with tests
valid_from:
  - WAITING_FOR_PROJECT_TEST_PLAN
valid_to:
  - INIT (proceed to Phase 1)
```

### 6. CREATE_PHASE_INTEGRATION_BRANCH_EARLY (New)
```yaml
state: CREATE_PHASE_INTEGRATION_BRANCH_EARLY
agent: orchestrator
description: Create phase-N-integration branch with tests
valid_from:
  - WAITING_FOR_PHASE_TEST_PLAN
valid_to:
  - SPAWN_CODE_REVIEWER_PHASE_IMPL
```

### 7. CREATE_WAVE_INTEGRATION_BRANCH_EARLY (New)
```yaml
state: CREATE_WAVE_INTEGRATION_BRANCH_EARLY
agent: orchestrator
description: Create wave integration branch with tests
valid_from:
  - WAITING_FOR_WAVE_TEST_PLAN
valid_to:
  - SPAWN_CODE_REVIEWER_WAVE_IMPL
```

## MODIFIED STATE FLOWS

### Project Start Flow (Modified)

**OLD FLOW:**
```
INIT → [Phase 1 Planning]
```

**NEW FLOW:**
```
INIT → SPAWN_ARCHITECT_MASTER_PLANNING → WAITING_FOR_MASTER_ARCHITECTURE →
SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING → WAITING_FOR_PROJECT_TEST_PLAN →
CREATE_PROJECT_INTEGRATION_BRANCH_EARLY → INIT (Phase 1)
```

### Phase Planning Flow (Modified per R341)

**OLD FLOW:**
```
SPAWN_ARCHITECT_PHASE_PLANNING → WAITING_FOR_ARCHITECTURE_PLAN →
SPAWN_CODE_REVIEWER_PHASE_IMPL → WAITING_FOR_IMPLEMENTATION_PLAN
```

**NEW FLOW (R341 + Early Branches):**
```
SPAWN_ARCHITECT_PHASE_PLANNING → WAITING_FOR_ARCHITECTURE_PLAN →
SPAWN_CODE_REVIEWER_PHASE_TEST_PLANNING → WAITING_FOR_PHASE_TEST_PLAN →
CREATE_PHASE_INTEGRATION_BRANCH_EARLY →
SPAWN_CODE_REVIEWER_PHASE_IMPL → WAITING_FOR_IMPLEMENTATION_PLAN
```

### Wave Planning Flow (Modified per R341)

**OLD FLOW:**
```
SPAWN_ARCHITECT_WAVE_PLANNING → WAITING_FOR_ARCHITECTURE_PLAN →
SPAWN_CODE_REVIEWER_WAVE_IMPL → WAITING_FOR_IMPLEMENTATION_PLAN
```

**NEW FLOW (R341 + Early Branches):**
```
SPAWN_ARCHITECT_WAVE_PLANNING → WAITING_FOR_ARCHITECTURE_PLAN →
SPAWN_CODE_REVIEWER_WAVE_TEST_PLANNING → WAITING_FOR_WAVE_TEST_PLAN →
CREATE_WAVE_INTEGRATION_BRANCH_EARLY →
SPAWN_CODE_REVIEWER_WAVE_IMPL → WAITING_FOR_IMPLEMENTATION_PLAN
```

## STATE MACHINE VALIDATION RULES

### New Validation Rules to Add

```yaml
validation_rules:
  - name: project_test_before_phase
    description: Project tests must exist before Phase 1
    check: |
      if transitioning_to(INIT) for Phase 1:
        assert project-integration branch exists
        assert PROJECT-TEST-PLAN.md exists
        
  - name: phase_test_before_implementation
    description: Phase tests must exist before implementation planning
    check: |
      if transitioning_to(SPAWN_CODE_REVIEWER_PHASE_IMPL):
        assert phase-N-integration branch exists
        assert PHASE-TEST-PLAN.md exists
        
  - name: wave_test_before_implementation  
    description: Wave tests must exist before implementation planning
    check: |
      if transitioning_to(SPAWN_CODE_REVIEWER_WAVE_IMPL):
        assert wave-integration branch exists
        assert WAVE-TEST-PLAN.md exists
        
  - name: integration_branch_has_tests
    description: Integration branches must contain tests
    check: |
      if branch_name contains "integration":
        assert tests/ directory exists
        assert TEST-HARNESS.sh exists
```

## TRANSITION MATRIX UPDATES

### Add to Valid Transitions

```yaml
# Project-level testing transitions
"INIT" → "SPAWN_ARCHITECT_MASTER_PLANNING"  # When no master plan
"SPAWN_ARCHITECT_MASTER_PLANNING" → "WAITING_FOR_MASTER_ARCHITECTURE"
"WAITING_FOR_MASTER_ARCHITECTURE" → "SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING"
"SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING" → "WAITING_FOR_PROJECT_TEST_PLAN"
"WAITING_FOR_PROJECT_TEST_PLAN" → "CREATE_PROJECT_INTEGRATION_BRANCH_EARLY"
"CREATE_PROJECT_INTEGRATION_BRANCH_EARLY" → "INIT"  # Start Phase 1

# Phase-level testing transitions (R341 compliant)
"WAITING_FOR_ARCHITECTURE_PLAN" → "SPAWN_CODE_REVIEWER_PHASE_TEST_PLANNING"
"SPAWN_CODE_REVIEWER_PHASE_TEST_PLANNING" → "WAITING_FOR_PHASE_TEST_PLAN"
"WAITING_FOR_PHASE_TEST_PLAN" → "CREATE_PHASE_INTEGRATION_BRANCH_EARLY"
"CREATE_PHASE_INTEGRATION_BRANCH_EARLY" → "SPAWN_CODE_REVIEWER_PHASE_IMPL"

# Wave-level testing transitions (R341 compliant)
"WAITING_FOR_ARCHITECTURE_PLAN" → "SPAWN_CODE_REVIEWER_WAVE_TEST_PLANNING"
"SPAWN_CODE_REVIEWER_WAVE_TEST_PLANNING" → "WAITING_FOR_WAVE_TEST_PLAN"
"WAITING_FOR_WAVE_TEST_PLAN" → "CREATE_WAVE_INTEGRATION_BRANCH_EARLY"
"CREATE_WAVE_INTEGRATION_BRANCH_EARLY" → "SPAWN_CODE_REVIEWER_WAVE_IMPL"
```

### Remove Invalid Transitions

```yaml
# These transitions are now INVALID (violate R341)
"WAITING_FOR_ARCHITECTURE_PLAN" → "SPAWN_CODE_REVIEWER_PHASE_IMPL"  # Must test first
"WAITING_FOR_ARCHITECTURE_PLAN" → "SPAWN_CODE_REVIEWER_WAVE_IMPL"   # Must test first
```

## ORCHESTRATOR STATE IMPLEMENTATIONS

### CREATE_PROJECT_INTEGRATION_BRANCH_EARLY

```bash
handle_create_project_integration_branch_early() {
    echo "🏗️ Creating project integration branch with tests..."
    
    # Create integration workspace
    PROJECT_INT_DIR="/efforts/project/integration-workspace"
    mkdir -p "$PROJECT_INT_DIR"
    cd "$PROJECT_INT_DIR"
    
    # Clone target repository
    TARGET_REPO=$(yq '.repository_path' target-repo-config.yaml)
    git clone "$TARGET_REPO" target-repo
    cd target-repo
    
    # Create project-integration branch
    git checkout -b project-integration
    
    # Copy project tests from test planning
    mkdir -p tests/project
    cp -r "$CLAUDE_PROJECT_DIR/project-tests/"* tests/project/
    
    # Commit tests
    git add tests/
    git commit -m "test: add project-level tests (TDD per R341)"
    git push -u origin project-integration
    
    # Update state
    update_state "INIT"  # Ready for Phase 1
    
    echo "✅ Project integration branch created with tests"
    echo "📍 Location: project-integration branch"
    echo "🧪 Tests: tests/project/"
}
```

### CREATE_PHASE_INTEGRATION_BRANCH_EARLY

```bash
handle_create_phase_integration_branch_early() {
    local PHASE=$(yq '.current_phase' orchestrator-state.json)
    
    echo "🏗️ Creating phase $PHASE integration branch with tests..."
    
    # Determine base branch (R308)
    if [ $PHASE -eq 1 ]; then
        BASE_BRANCH="project-integration"  # Build on project
    else
        BASE_BRANCH="phase-$((PHASE-1))-integration"  # Previous phase
    fi
    
    # Create phase integration workspace
    PHASE_INT_DIR="/efforts/phase${PHASE}/integration-workspace"
    mkdir -p "$PHASE_INT_DIR"
    cd "$PHASE_INT_DIR"
    
    # Clone and create branch
    git clone "$TARGET_REPO" target-repo
    cd target-repo
    git checkout -b "phase-${PHASE}-integration" "$BASE_BRANCH"
    
    # Copy phase tests
    mkdir -p "tests/phase${PHASE}"
    cp -r "$CLAUDE_PROJECT_DIR/phase-tests/phase-${PHASE}/"* \
          "tests/phase${PHASE}/"
    
    # Commit tests
    git add tests/
    git commit -m "test: add phase ${PHASE} tests (TDD per R341)"
    git push -u origin "phase-${PHASE}-integration"
    
    # Update state
    update_state "SPAWN_CODE_REVIEWER_PHASE_IMPL"
    
    echo "✅ Phase ${PHASE} integration branch created with tests"
}
```

### CREATE_WAVE_INTEGRATION_BRANCH_EARLY

```bash
handle_create_wave_integration_branch_early() {
    local PHASE=$(yq '.current_phase' orchestrator-state.json)
    local WAVE=$(yq '.current_wave' orchestrator-state.json)
    
    echo "🏗️ Creating phase $PHASE wave $WAVE integration branch..."
    
    # Determine base branch (R308/R336)
    if [ $WAVE -eq 1 ]; then
        BASE_BRANCH="phase-${PHASE}-integration"  # Phase base
    else
        BASE_BRANCH="phase-${PHASE}-wave-$((WAVE-1))-integration"  # Previous wave
    fi
    
    # Create wave integration workspace
    WAVE_INT_DIR="/efforts/phase${PHASE}/wave${WAVE}/integration-workspace"
    mkdir -p "$WAVE_INT_DIR"
    cd "$WAVE_INT_DIR"
    
    # Clone and create branch
    git clone "$TARGET_REPO" target-repo
    cd target-repo
    git checkout -b "phase-${PHASE}-wave-${WAVE}-integration" "$BASE_BRANCH"
    
    # Copy wave tests
    mkdir -p "tests/phase${PHASE}/wave${WAVE}"
    cp -r "$CLAUDE_PROJECT_DIR/wave-tests/phase-${PHASE}/wave-${WAVE}/"* \
          "tests/phase${PHASE}/wave${WAVE}/"
    
    # Commit tests
    git add tests/
    git commit -m "test: add phase ${PHASE} wave ${WAVE} tests (TDD per R341)"
    git push -u origin "phase-${PHASE}-wave-${WAVE}-integration"
    
    # Update state
    update_state "SPAWN_CODE_REVIEWER_WAVE_IMPL"
    
    echo "✅ Wave ${WAVE} integration branch created with tests"
}
```

## CODE REVIEWER STATE IMPLEMENTATIONS

### PROJECT_TEST_PLANNING State

```yaml
state: PROJECT_TEST_PLANNING
agent: code-reviewer
deliverables:
  - PROJECT-TEST-PLAN.md
  - PROJECT-TEST-HARNESS.sh
  - tests/project/*.test.*
  - PROJECT-DEMO-SCENARIOS.md
  - test-to-phase-mapping.json
```

### PHASE_TEST_PLANNING State (per R341)

```yaml
state: PHASE_TEST_PLANNING
agent: code-reviewer
deliverables:
  - PHASE-${N}-TEST-PLAN.md
  - PHASE-${N}-TEST-HARNESS.sh
  - tests/phase${N}/*.test.*
  - PHASE-${N}-DEMO-SCENARIOS.md
  - test-to-wave-mapping.json
```

### WAVE_TEST_PLANNING State (per R341)

```yaml
state: WAVE_TEST_PLANNING
agent: code-reviewer
deliverables:
  - WAVE-${M}-TEST-PLAN.md
  - WAVE-${M}-TEST-HARNESS.sh
  - tests/phase${N}/wave${M}/*.test.*
  - WAVE-${M}-DEMO-SCENARIOS.md
  - test-to-effort-mapping.json
```

## GRADING UPDATES

### New Grading Criteria

```yaml
test_planning_grading:
  project_tests_missing: -50%  # No project tests created
  project_tests_after_phase: -40%  # Created too late
  phase_tests_missing: -40%  # Violates R341
  wave_tests_missing: -40%  # Violates R341
  tests_not_in_branch: -30%  # Storage violation
  tests_created_after_impl: -100%  # Critical R341 violation
  
integration_branch_grading:
  branch_created_late: -20%  # Should be early
  tests_missing_from_branch: -30%  # Storage failure
  branch_not_incremental: -25%  # Violates R308
  wrong_base_branch: -35%  # Violates R336
```

## BACKWARD COMPATIBILITY

### For Projects Without Project Tests

```bash
# Detection in INIT state
if [ ! -f "PROJECT-TEST-PLAN.md" ]; then
    if [ $CURRENT_PHASE -eq 1 ] && [ $CURRENT_WAVE -eq 1 ]; then
        # New project - require project tests
        transition_to "SPAWN_ARCHITECT_MASTER_PLANNING"
    else
        # Existing project - continue without project tests
        echo "⚠️ Legacy project - skipping project test requirement"
        # Continue with phase/wave tests only
    fi
fi
```

## SUMMARY OF CHANGES

### Total New States: 7
1. SPAWN_ARCHITECT_MASTER_PLANNING
2. WAITING_FOR_MASTER_ARCHITECTURE
3. SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING
4. WAITING_FOR_PROJECT_TEST_PLAN
5. CREATE_PROJECT_INTEGRATION_BRANCH_EARLY
6. CREATE_PHASE_INTEGRATION_BRANCH_EARLY
7. CREATE_WAVE_INTEGRATION_BRANCH_EARLY

### Modified Flows: 3
1. Project initialization flow
2. Phase planning flow (R341 compliance)
3. Wave planning flow (R341 compliance)

### New Rules Needed: 1
- R342: Early Integration Branch Creation Protocol

### Updated Rules: 2
- R341: Add project-level testing requirements
- R308: Clarify incremental branching with early creation

This completes the state machine updates needed for implementing the project-level testing and early integration branch creation strategy.