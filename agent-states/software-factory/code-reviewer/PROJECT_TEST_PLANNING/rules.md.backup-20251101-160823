# CODE REVIEWER STATE: PROJECT_TEST_PLANNING

## Overview
This state creates comprehensive project-level tests from the master architecture within SF 3.0, implementing TDD at the highest level before any phase begins.

## SF 3.0 Project Test Planning Context

In this state, the Code Reviewer creates project-wide test framework:
- Reads master architecture from orchestrator-state-v3.json to understand overall system design
- Creates project test harness that orchestrator records in `test_plans.project` section of orchestrator-state-v3.json
- Establishes testing standards and frameworks for all phases/waves to follow
- Reports test framework location for orchestrator tracking per R340 with atomic updates per R288
- Project tests serve as ultimate acceptance criteria for entire system

## Entry Criteria
- Spawned by orchestrator from SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING
- Master architecture exists (PROJECT-ARCHITECTURE.md or MASTER-ARCHITECTURE.md)
- No existing PROJECT-TEST-PLAN.md

## State Responsibilities

### 1. Read Master Architecture
```bash
# Find and read the master architecture
if [ -f "PROJECT-ARCHITECTURE.md" ]; then
    ARCH_FILE="PROJECT-ARCHITECTURE.md"
elif [ -f "MASTER-ARCHITECTURE.md" ]; then
    ARCH_FILE="MASTER-ARCHITECTURE.md"
else
    echo "❌ No master architecture found!"
    exit 1
fi

echo "📖 Reading architecture from $ARCH_FILE"
```

### 2. Create Project Test Plan
Create PROJECT-TEST-PLAN.md with:
- Cross-phase integration tests
- End-to-end workflow tests
- System-level validation tests
- Performance benchmarks
- Security validations
- API contract tests

### 3. Create Test Harness
Create PROJECT-TEST-HARNESS.sh:
```bash
#!/bin/bash
# PROJECT-TEST-HARNESS.sh
# Runs all project-level tests

echo "🧪 Running project-level tests..."

# Integration tests
echo "📊 Cross-phase integration tests..."
./tests/project/integration.test.sh

# E2E tests
echo "🔄 End-to-end workflow tests..."
./tests/project/e2e.test.sh

# System validation
echo "✅ System validation tests..."
./tests/project/validation.test.sh

# Report results
echo "📋 Test Summary:"
echo "  - Integration: $INT_RESULT"
echo "  - E2E: $E2E_RESULT"
echo "  - Validation: $VAL_RESULT"
```

### 4. Create Test Files
Create actual test files in project-tests/:
```bash
mkdir -p project-tests/
# Create test files based on architecture
echo "Creating cross-phase integration tests..."
echo "Creating E2E workflow tests..."
echo "Creating system validation tests..."
```

### 5. Create Demo Scenarios
Create PROJECT-DEMO-SCENARIOS.md:
- User journey demonstrations
- Feature showcases
- Integration demonstrations
- Performance scenarios

### 6. Create Test-to-Phase Mapping
Create test-to-phase-mapping.json:
```json
{
  "project_tests": {
    "integration_test_1": ["phase1", "phase2"],
    "e2e_test_1": ["phase1", "phase2", "phase3"],
    "validation_test_1": ["phase3"]
  }
}
```

### 7. Report Test Plan Location (R340)
```bash
# R340: Report test plan and harness locations to orchestrator
TEST_PLAN_PATH="$(pwd)/PROJECT-TEST-PLAN.md"
TEST_HARNESS_PATH="$(pwd)/PROJECT-TEST-HARNESS.sh"
DEMO_SCENARIOS_PATH="$(pwd)/PROJECT-DEMO-SCENARIOS.md"

echo "📋 Test Plan: $TEST_PLAN_PATH"
echo "📋 Test Harness: $TEST_HARNESS_PATH"
echo "📋 Demo Scenarios: $DEMO_SCENARIOS_PATH"
echo "R340: Created project test plan at: $TEST_PLAN_PATH"
```

### 8. Complete and Report
```bash
# Mark completion
echo "✅ Project test planning complete"
echo "📋 Deliverables created:"
echo "  - PROJECT-TEST-PLAN.md"
echo "  - PROJECT-TEST-HARNESS.sh"
echo "  - project-tests/*.test.*"
echo "  - PROJECT-DEMO-SCENARIOS.md"
echo "  - test-to-phase-mapping.json"

# Update state
update_state "COMPLETED"
```

## Exit Criteria
- All project test artifacts created
- Tests cover cross-phase integration
- Test harness executable
- State transitions to COMPLETED

## Success Metrics
- ✅ PROJECT-TEST-PLAN.md created
- ✅ PROJECT-TEST-HARNESS.sh executable
- ✅ Test files in project-tests/
- ✅ Demo scenarios documented
- ✅ Test mapping defined

## Deliverables
1. **PROJECT-TEST-PLAN.md** - Comprehensive test strategy
2. **PROJECT-TEST-HARNESS.sh** - Test execution script
3. **project-tests/*.test.*** - Actual test files
4. **PROJECT-DEMO-SCENARIOS.md** - Demo scenarios
5. **test-to-phase-mapping.json** - Test dependencies


## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**


## Related Rules
- R341: TDD - tests before implementation
- R342: Tests will be stored in early-created branch
- R340: Planning file metadata tracking
