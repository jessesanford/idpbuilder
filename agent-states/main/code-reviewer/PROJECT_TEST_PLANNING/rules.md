# CODE REVIEWER STATE: PROJECT_TEST_PLANNING

## Overview
This state creates comprehensive project-level tests from the master architecture, implementing TDD at the highest level before any phase begins.

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

### 7. Complete and Report
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

## Related Rules
- R341: TDD - tests before implementation
- R342: Tests will be stored in early-created branch
- R340: Planning file metadata tracking