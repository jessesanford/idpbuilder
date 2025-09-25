# SPAWN_CODE_REVIEWER_WAVE_TEST_PLANNING State Rules

# PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R355** - Code Reviewer Test Planning
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R355-code-reviewer-test-planning.md`

4. **R287** - TODO Persistence Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`

5. **R288** - State File Update Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-requirements.md`

6. **R304** - Mandatory Line Counter Usage (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-usage.md`

7. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-after-spawn.md`

8. **R324** - State Transition Validation (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R324-state-transition-validation.md`


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

## State Purpose
Spawn Code Reviewer to create wave-level tests BEFORE implementation planning begins. This enforces Test-Driven Development (TDD) at the wave level.

## 🔴🔴🔴 CRITICAL: WAVE-LEVEL TDD ENFORCEMENT 🔴🔴🔴

**WAVE TESTS COMPLEMENT PHASE TESTS!**
- Phase tests validate overall functionality
- Wave tests validate wave-specific features
- Both must exist before implementation

## Entry Conditions
- Wave architecture plan exists (from Architect)
- Current state is WAITING_FOR_ARCHITECTURE_PLAN
- Wave-level test planning not yet done

## Required Actions

### 1. Prepare Wave Test Planning Context
```bash
# Load wave architecture
WAVE_ARCH_FILE="wave-plans/phase-${PHASE_NUM}/wave-${WAVE_NUM}-architecture.md"

# Create test planning directory
mkdir -p wave-tests/phase-${PHASE_NUM}/wave-${WAVE_NUM}

# Prepare test planning instructions
cat > test-planning-instructions.md << 'EOF'
# WAVE TEST PLANNING INSTRUCTIONS (TDD)

## Your Mission
Create comprehensive tests for Wave ${WAVE_NUM} of Phase ${PHASE_NUM} BEFORE implementation begins.

## Required Deliverables
1. WAVE-TEST-PLAN.md - Wave-specific functionality to test
2. WAVE-TEST-HARNESS.sh - Test runner for wave
3. tests/wave/*.test.* - Actual functional tests (failing initially)
4. WAVE-DEMO-PLAN.md - Demo scenarios for wave features

## Test Requirements
- Tests must be FUNCTIONAL/BEHAVIORAL (not unit tests)
- Tests must validate wave-specific capabilities from architecture
- Tests must integrate with phase-level tests
- Tests must FAIL initially (no implementation yet)
- Tests must guide effort planning

## Effort Test Guidance
- Each effort plan must reference which wave tests it will make pass
- Effort success is measured by passing assigned tests
- Test assignments prevent scope creep

## Demo Integration (R330/R291 Consolidation)
- Wave demos must showcase integrated functionality
- Each wave test should have demo scenarios
- Single-effort waves still need integration demos (R291)

## Success Criteria
- All wave architectural promises have tests
- Test harness integrates with phase harness
- Demo scenarios clearly defined
- Tests ready for effort planning to reference
EOF
```

### 2. Spawn Code Reviewer for Wave Test Planning
```bash
/spawn-agent code-reviewer \
    --state WAVE_TEST_PLANNING \
    --working-dir wave-tests/phase-${PHASE_NUM}/wave-${WAVE_NUM} \
    --context "Create wave-level tests from architecture" \
    --deliverables "WAVE-TEST-PLAN.md,WAVE-TEST-HARNESS.sh,tests/,WAVE-DEMO-PLAN.md" \
    --timeout 20m
```

### 3. Update State File
```json
{
    "current_state": "WAITING_FOR_WAVE_TEST_PLAN",
    "wave_test_planning": {
        "phase": "${PHASE_NUM}",
        "wave": "${WAVE_NUM}",
        "spawned_at": "timestamp",
        "code_reviewer_id": "agent_id",
        "expected_deliverables": [
            "WAVE-TEST-PLAN.md",
            "WAVE-TEST-HARNESS.sh",
            "tests/wave/*.test.*",
            "WAVE-DEMO-PLAN.md"
        ]
    }
}
```

## Exit Conditions
- Transition to WAITING_FOR_WAVE_TEST_PLAN
- Code Reviewer spawned with TDD instructions
- State file updated with wave test context

## Success Criteria
- ✅ Wave tests complement phase tests
- ✅ Test planning happens BEFORE implementation planning
- ✅ Demo requirements integrated (R330/R291)
- ✅ Tests will guide effort planning

## Related Rules
- R341: Test-Driven Development Enforcement (new)
- R330: Demo Planning Requirements (consolidated)
- R291: Integration Demo Requirement (consolidated)
- R210: Wave Architecture Planning
