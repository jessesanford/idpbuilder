# PHASE_TEST_PLANNING State Rules (Code Reviewer)

## State Purpose
Create comprehensive phase-level tests from architecture BEFORE any implementation planning begins within SF 3.0. This is the foundation of Test-Driven Development (TDD) at the phase level.

## SF 3.0 Test Planning Context

In this state, the Code Reviewer creates phase test plans:
- Reads phase architecture from orchestrator-state-v3.json `metadata_locations.phase_architecture_plans` per R340
- Creates test harness and plan that orchestrator records in `test_plans.phase` section of orchestrator-state-v3.json
- Defines test requirements that will guide implementation and validation
- Reports test plan location for orchestrator tracking with atomic updates per R288
- Test plans become the acceptance criteria for phase completion

## 🔴🔴🔴 CRITICAL: YOU ARE CREATING THE TESTS FIRST! 🔴🔴🔴

**TDD WORKFLOW:**
1. Architecture defines WHAT to build
2. **YOU create tests that verify it works** ← YOU ARE HERE
3. Implementation plans define HOW to make tests pass
4. Engineers implement to pass the tests

## Entry Conditions
- Spawned by orchestrator in SPAWN_CODE_REVIEWER_PHASE_TEST_PLANNING
- Phase architecture plan available
- No implementation planning done yet

## Required Actions

### 1. Load Phase Architecture
```bash
# Read the phase architecture
PHASE_ARCH=$(find . -name "*phase*architecture*.md" | head -1)
if [ -z "$PHASE_ARCH" ]; then
    echo "❌ ERROR: Cannot find phase architecture!"
    exit 1
fi

echo "📖 Loading phase architecture from: $PHASE_ARCH"
cat "$PHASE_ARCH"
```

### 2. Create PHASE-TEST-PLAN.md
```markdown
# PHASE ${PHASE_NUM} TEST PLAN (TDD)
Generated: [timestamp]
Status: Tests Written BEFORE Implementation

## Architecture Coverage
[Map each architectural promise to specific tests]

### Core Capabilities to Test
Based on architecture, we must verify:
1. [Capability 1] - Tested by: test_capability_1.test.js
2. [Capability 2] - Tested by: test_capability_2.test.js
3. [Capability 3] - Tested by: test_capability_3.test.js

## Test Categories

### Functional Tests (MANDATORY)
- [ ] API endpoint tests
- [ ] Data flow validation
- [ ] Integration between components
- [ ] Error handling scenarios
- [ ] Performance requirements

### Behavioral Tests
- [ ] User workflows work end-to-end
- [ ] System responds correctly to inputs
- [ ] State management is correct
- [ ] Side effects are handled

## Test Assignments for Waves
[Waves will pick specific tests to implement]

### Wave 1 Test Targets
- Tests: test_core_setup.*, test_basic_api.*
- Coverage: 30% of phase tests

### Wave 2 Test Targets  
- Tests: test_advanced_features.*, test_integrations.*
- Coverage: 40% of phase tests

### Wave 3 Test Targets
- Tests: test_edge_cases.*, test_performance.*
- Coverage: 30% of phase tests

## Success Metrics
- All tests defined: YES
- Tests executable: YES (will fail initially)
- Architecture covered: 100%
- Demo scenarios defined: YES (see PHASE-DEMO-PLAN.md)
```

### 3. Create PHASE-TEST-HARNESS.sh
```bash
cat > PHASE-TEST-HARNESS.sh << 'EOF'
#!/bin/bash
# PHASE TEST HARNESS (TDD)
# These tests WILL FAIL until implementation is complete

echo "🧪 PHASE ${PHASE_NUM} TEST HARNESS"
echo "================================"
echo "Mode: TDD (Tests Written First)"
echo ""

FAILED=0
PASSED=0
TOTAL=0

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Run all phase tests
run_phase_tests() {
    echo "📋 Running Phase-Level Tests..."
    
    for test_file in tests/phase/*.test.*; do
        if [ -f "$test_file" ]; then
            ((TOTAL++))
            echo -n "  Testing: $(basename $test_file) ... "
            
            # Run test (will fail in TDD red phase)
            if ./"$test_file" > /dev/null 2>&1; then
                echo -e "${GREEN}PASS${NC}"
                ((PASSED++))
            else
                echo -e "${RED}FAIL${NC} (Expected in TDD)"
                ((FAILED++))
            fi
        fi
    done
}

# Check implementation status
check_implementation_status() {
    if [ -d "../implementation" ]; then
        echo -e "${YELLOW}⚠️  Implementation detected - tests may pass${NC}"
    else
        echo -e "${RED}🔴 No implementation yet - all tests should fail (TDD Red Phase)${NC}"
    fi
}

# Main execution
echo "Starting test execution..."
check_implementation_status
echo ""

run_phase_tests

echo ""
echo "================================"
echo "RESULTS:"
echo "  Total Tests: $TOTAL"
echo -e "  Passed: ${GREEN}$PASSED${NC}"
echo -e "  Failed: ${RED}$FAILED${NC}"

if [ $FAILED -eq $TOTAL ]; then
    echo ""
    echo -e "${YELLOW}✅ TDD STATUS: All tests failing as expected!${NC}"
    echo "   Ready for implementation planning."
    exit 0  # This is PROJECT_DONE in TDD red phase
else
    echo ""
    echo -e "${YELLOW}⚠️  Some tests passing without implementation?${NC}"
    exit 0  # Still OK, but unexpected
fi
EOF

chmod +x PHASE-TEST-HARNESS.sh
```

### 4. Create Actual Test Files
```bash
# Create test directory
mkdir -p tests/phase

# Example functional test
cat > tests/phase/test_core_functionality.test.js << 'EOF'
#!/usr/bin/env node
/**
 * PHASE TEST: Core Functionality
 * Status: Written BEFORE implementation (TDD)
 * Expected: FAIL until implementation complete
 */

const assert = require('assert');

// This test WILL FAIL - no implementation yet
try {
    // Test that core API exists
    const api = require('../../src/api');  // Won't exist yet
    assert(api, 'API module should exist');
    
    // Test basic endpoint
    const result = api.get('/health');
    assert.equal(result.status, 'healthy');
    
    console.log('✅ Core functionality test passed');
    process.exit(0);
} catch (error) {
    console.log('❌ Core functionality test failed (expected in TDD)');
    console.log(`   Reason: ${error.message}`);
    process.exit(1);
}
EOF

chmod +x tests/phase/*.test.*
```

### 5. Create PHASE-DEMO-PLAN.md (Consolidating R330/R291)
```markdown
# PHASE ${PHASE_NUM} DEMO PLAN (TDD-Integrated)
Generated: [timestamp]
Consolidates: R330 (Demo Planning) + R291 (Integration Demo)

## Demo Objectives (From Tests)
Each demo scenario corresponds to specific tests:

1. **Core Setup Demo**
   - Test: test_core_functionality.test.js
   - Shows: System initialization and basic operations
   - Script: demo-core-setup.sh

2. **API Operations Demo**
   - Test: test_api_endpoints.test.js
   - Shows: All API endpoints working
   - Script: demo-api-operations.sh

3. **Integration Demo**
   - Test: test_integrations.test.js
   - Shows: Components working together
   - Script: demo-integrations.sh

## Demo Scenarios (R330 Compliance)

### Scenario 1: System Initialization
```bash
# Setup
./setup-environment.sh

# Action
./start-system.sh

# Verification
curl http://localhost:3000/health

# Expected
{"status": "healthy", "version": "1.0.0"}
```

### Scenario 2: Create Operations
```bash
# Action
curl -X POST http://localhost:3000/api/resource \
  -H "Content-Type: application/json" \
  -d '{"name": "test", "value": 42}'

# Expected
{"id": "uuid", "status": "created"}
```

## Integration Demo Requirements (R291)
- Build must succeed
- All tests must pass
- Demo script must execute without errors
- Features must be visibly working

## Demo Artifacts (Not counted in line limit per R007)
- demo-phase-${PHASE_NUM}.sh: Main demo script
- DEMO.md: This documentation
- test-data/: Sample data files
- Total demo lines: ~200 (excluded from 800-line limit)
```

### 6. Report Test Locations to Orchestrator (R340 Compliance)
```markdown
## 📋 TEST PLAN CREATED - METADATA FOR ORCHESTRATOR

**Type**: phase_test_plan
**Phase**: ${PHASE_NUM}
**Test Plan Path**: /phase-tests/phase-${PHASE_NUM}/PHASE-TEST-PLAN.md
**Test Harness Path**: /phase-tests/phase-${PHASE_NUM}/PHASE-TEST-HARNESS.sh
**Demo Plan Path**: /phase-tests/phase-${PHASE_NUM}/PHASE-DEMO-PLAN.md
**Test Directory**: /phase-tests/phase-${PHASE_NUM}/tests/phase
**Created At**: [ISO-8601-timestamp]
**Created By**: code-reviewer
**Status**: active
**Test Count**: [number of test files created]
**TDD Phase**: red (tests failing, no implementation yet)

ORCHESTRATOR: Please update test_plans.phase["phase${PHASE_NUM}"] in state file with:
- test_plan_path
- test_harness_path
- demo_plan_path
- test_dir
- created_at
- test_count
- tdd_phase = "red"

This enables SW Engineers to find tests per R340 pattern.
```

### 7. Commit Test Artifacts
```bash
# Add all test artifacts
git add -A

# Commit with TDD message
git commit -m "tdd: phase ${PHASE_NUM} tests created BEFORE implementation

- Created comprehensive test plan covering architecture
- Implemented failing tests (TDD red phase)
- Defined demo scenarios linked to tests
- Ready for implementation planning to target these tests
- Reported test locations for state tracking (R340)

Test-Driven Development enforced per R341"

git push
```

## Exit Conditions
- All test deliverables created
- Tests are failing (no implementation yet)
- Transition to COMPLETED state
- Ready for implementation planning phase

## Success Criteria
- ✅ PHASE-TEST-PLAN.md covers all architecture
- ✅ PHASE-TEST-HARNESS.sh is executable
- ✅ Actual test files created and failing
- ✅ PHASE-DEMO-PLAN.md links demos to tests
- ✅ Tests will guide implementation planning


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
- R341: Test-Driven Development Enforcement
- R330: Demo Planning Requirements (consolidated here)
- R291: Integration Demo Requirement (consolidated here)
- R210: Architecture Planning (what we're testing)
