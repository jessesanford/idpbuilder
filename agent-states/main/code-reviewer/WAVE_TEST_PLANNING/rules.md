# WAVE_TEST_PLANNING State Rules (Code Reviewer)

## State Purpose
Create comprehensive wave-level tests from architecture BEFORE wave implementation planning begins. Wave tests complement phase tests and guide effort planning.

## 🔴🔴🔴 CRITICAL: WAVE TESTS GUIDE EFFORT PLANNING! 🔴🔴🔴

**TDD AT WAVE LEVEL:**
1. Wave architecture defines wave-specific features
2. **YOU create tests for those features** ← YOU ARE HERE
3. Implementation plans assign tests to efforts
4. Each effort implements to pass assigned tests

## Entry Conditions
- Spawned by orchestrator in SPAWN_CODE_REVIEWER_WAVE_TEST_PLANNING
- Wave architecture plan available
- Phase tests already exist
- No wave implementation planning done yet

## Required Actions

### 1. Load Wave Architecture and Phase Tests
```bash
# Read the wave architecture
WAVE_ARCH=$(find . -name "*wave*architecture*.md" | head -1)
if [ -z "$WAVE_ARCH" ]; then
    echo "❌ ERROR: Cannot find wave architecture!"
    exit 1
fi

# Reference phase tests
PHASE_TESTS="../../../phase-tests/phase-${PHASE_NUM}/PHASE-TEST-PLAN.md"
if [ -f "$PHASE_TESTS" ]; then
    echo "📖 Referencing phase-level tests for context"
fi

echo "📖 Loading wave architecture from: $WAVE_ARCH"
cat "$WAVE_ARCH"
```

### 2. Create WAVE-TEST-PLAN.md
```markdown
# WAVE ${WAVE_NUM} TEST PLAN (TDD)
Phase: ${PHASE_NUM}
Generated: [timestamp]
Status: Tests Written BEFORE Implementation

## Wave Architecture Coverage
[Map wave-specific features to tests]

### Wave ${WAVE_NUM} Features to Test
Based on wave architecture:
1. [Feature 1] - Tested by: test_wave_feature_1.test.js
2. [Feature 2] - Tested by: test_wave_feature_2.test.js
3. [Feature 3] - Tested by: test_wave_feature_3.test.js

## Relationship to Phase Tests
- Phase tests validate: [overall functionality]
- Wave tests validate: [wave-specific features]
- Integration point: [how they work together]

## Test Categories for This Wave

### Feature Tests (Wave-Specific)
- [ ] Feature A implementation
- [ ] Feature B implementation  
- [ ] Feature integration
- [ ] Wave-specific workflows

### Integration Tests
- [ ] Integration with previous waves (if any)
- [ ] Data flow between components
- [ ] API contracts maintained
- [ ] No regression in existing features

## Effort Test Assignments (CRITICAL FOR TDD)
[Each effort will be assigned specific tests to pass]

### Effort 1 Test Targets
- Tests to pass: test_feature_a.*, test_integration_1.*
- Expected lines: ~250
- Success criteria: All assigned tests pass

### Effort 2 Test Targets
- Tests to pass: test_feature_b.*, test_workflow_1.*
- Expected lines: ~300
- Success criteria: All assigned tests pass

### Effort 3 Test Targets
- Tests to pass: test_feature_c.*, test_integration_2.*
- Expected lines: ~250
- Success criteria: All assigned tests pass

## Wave Demo Requirements (R330/R291)
- Each test has corresponding demo scenario
- Integration demo required (even for single-effort waves)
- See WAVE-DEMO-PLAN.md for details

## Success Metrics
- Wave architecture covered: 100%
- Tests defined for all features: YES
- Effort assignments clear: YES
- Demo scenarios defined: YES
```

### 3. Create WAVE-TEST-HARNESS.sh
```bash
cat > WAVE-TEST-HARNESS.sh << 'EOF'
#!/bin/bash
# WAVE TEST HARNESS (TDD)
# These tests WILL FAIL until wave implementation is complete

echo "🧪 WAVE ${WAVE_NUM} TEST HARNESS (Phase ${PHASE_NUM})"
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

# Run wave-specific tests
run_wave_tests() {
    echo "📋 Running Wave ${WAVE_NUM} Tests..."
    
    for test_file in tests/wave/*.test.*; do
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

# Check for effort implementations
check_effort_status() {
    echo ""
    echo "📊 Effort Implementation Status:"
    
    for effort_dir in ../../efforts/*/; do
        if [ -d "$effort_dir" ]; then
            effort_name=$(basename "$effort_dir")
            if [ -f "$effort_dir/implementation-complete.flag" ]; then
                echo "  ✅ $effort_name: Implementation complete"
            else
                echo "  ⏳ $effort_name: Not implemented yet"
            fi
        fi
    done
}

# Integration with phase tests
run_phase_integration() {
    echo ""
    echo "🔗 Phase Integration Check:"
    
    PHASE_HARNESS="../../../phase-tests/phase-${PHASE_NUM}/PHASE-TEST-HARNESS.sh"
    if [ -x "$PHASE_HARNESS" ]; then
        echo "  Running subset of phase tests relevant to this wave..."
        # Run only specific phase tests for this wave
    else
        echo "  Phase test harness not found"
    fi
}

# Main execution
echo "Starting wave test execution..."
echo ""

run_wave_tests
check_effort_status
run_phase_integration

echo ""
echo "================================"
echo "RESULTS:"
echo "  Total Wave Tests: $TOTAL"
echo -e "  Passed: ${GREEN}$PASSED${NC}"
echo -e "  Failed: ${RED}$FAILED${NC}"

if [ $FAILED -eq $TOTAL ]; then
    echo ""
    echo -e "${YELLOW}✅ TDD STATUS: All tests failing as expected!${NC}"
    echo "   Ready for wave implementation planning."
    exit 0  # This is SUCCESS in TDD red phase
else
    echo ""
    echo -e "${YELLOW}⚠️  Some tests passing without implementation?${NC}"
    exit 0  # Still OK, but unexpected
fi
EOF

chmod +x WAVE-TEST-HARNESS.sh
```

### 4. Create Wave-Specific Test Files
```bash
# Create test directory
mkdir -p tests/wave

# Example wave feature test
cat > tests/wave/test_wave_${WAVE_NUM}_feature.test.js << 'EOF'
#!/usr/bin/env node
/**
 * WAVE TEST: Wave-Specific Feature
 * Wave: ${WAVE_NUM}
 * Status: Written BEFORE implementation (TDD)
 * Assigned to: [Will be assigned to specific effort]
 */

const assert = require('assert');

// This test WILL FAIL - no implementation yet
try {
    // Test wave-specific feature
    const waveFeature = require('../../src/wave-feature');  // Won't exist yet
    assert(waveFeature, 'Wave feature module should exist');
    
    // Test specific functionality
    const result = waveFeature.process({type: 'wave-specific'});
    assert.equal(result.status, 'processed');
    
    console.log('✅ Wave feature test passed');
    process.exit(0);
} catch (error) {
    console.log('❌ Wave feature test failed (expected in TDD)');
    console.log(`   Reason: ${error.message}`);
    process.exit(1);
}
EOF

# Create integration test
cat > tests/wave/test_wave_${WAVE_NUM}_integration.test.sh << 'EOF'
#!/bin/bash
# WAVE INTEGRATION TEST
# Tests integration between wave components

echo "Testing wave ${WAVE_NUM} integration..."

# Will fail until implementation exists
if [ -f "../../build/wave-${WAVE_NUM}.js" ]; then
    node ../../build/wave-${WAVE_NUM}.js --test
    exit $?
else
    echo "❌ Wave implementation not found (expected in TDD)"
    exit 1
fi
EOF

chmod +x tests/wave/*.test.*
```

### 5. Create WAVE-DEMO-PLAN.md (Consolidating R330/R291)
```markdown
# WAVE ${WAVE_NUM} DEMO PLAN (TDD-Integrated)
Phase: ${PHASE_NUM}
Generated: [timestamp]
Consolidates: R330 (Demo Planning) + R291 (Integration Demo)

## Wave Demo Objectives (Linked to Tests)
Demonstrate that wave ${WAVE_NUM} features work:

1. **Feature A Demo**
   - Test: test_wave_${WAVE_NUM}_feature_a.test.js
   - Effort: Will be assigned to Effort 1
   - Demo: Show feature A working end-to-end

2. **Feature B Demo**
   - Test: test_wave_${WAVE_NUM}_feature_b.test.js
   - Effort: Will be assigned to Effort 2
   - Demo: Show feature B integration

3. **Wave Integration Demo**
   - Test: test_wave_${WAVE_NUM}_integration.test.sh
   - Shows: All wave components working together
   - Required: Even for single-effort waves (R291)

## Demo Scenarios (R330 Compliance)

### Scenario 1: Wave Feature A
```bash
# Setup
./setup-wave-${WAVE_NUM}.sh

# Action
./run-feature-a.sh --input test-data/feature-a.json

# Expected Output
{
  "status": "success",
  "feature": "a",
  "wave": ${WAVE_NUM}
}

# Verification
grep "Feature A completed" logs/wave-${WAVE_NUM}.log
```

### Scenario 2: Wave Integration
```bash
# Action - Test all wave components together
./demo-wave-integration.sh

# Expected - All components interact correctly
- Component A sends data to Component B
- Component B processes and returns result
- Integration points verified
```

## Effort Demo Assignments
Each effort will implement demos for their assigned tests:

| Effort | Tests | Demo Scripts |
|--------|-------|--------------|
| Effort 1 | test_feature_a.* | demo-effort-1.sh |
| Effort 2 | test_feature_b.* | demo-effort-2.sh |
| Effort 3 | test_integration.* | demo-effort-3.sh |

## Wave Integration Demo (R291 MANDATORY)
```bash
# Even single-effort waves need integration demo
./demo-wave-${WAVE_NUM}-integration.sh

# Must demonstrate:
- [ ] Build succeeds
- [ ] All wave tests pass
- [ ] Features work together
- [ ] No regression in previous waves
```

## Demo Size Planning (R330)
- Demo scripts: ~150 lines (excluded from 800-line limit)
- Demo documentation: ~100 lines (excluded)
- Test data: ~50 lines (excluded)
- Total demo overhead: ~300 lines (NOT counted per R007)
```

### 6. Create Test-to-Effort Mapping
```bash
# Create mapping file for effort planning
cat > test-effort-mapping.json << EOF
{
  "wave": ${WAVE_NUM},
  "phase": ${PHASE_NUM},
  "test_assignments": {
    "effort_1": {
      "tests": ["test_wave_${WAVE_NUM}_feature_a.test.js"],
      "expected_lines": 250,
      "demo_required": true
    },
    "effort_2": {
      "tests": ["test_wave_${WAVE_NUM}_feature_b.test.js"],
      "expected_lines": 300,
      "demo_required": true
    },
    "effort_3": {
      "tests": ["test_wave_${WAVE_NUM}_integration.test.sh"],
      "expected_lines": 250,
      "demo_required": true
    }
  },
  "total_wave_lines": 800,
  "integration_demo_required": true
}
EOF
```

### 7. Report Test Locations to Orchestrator (R340 Compliance)
```markdown
## 📋 WAVE TEST PLAN CREATED - METADATA FOR ORCHESTRATOR

**Type**: wave_test_plan
**Phase**: ${PHASE_NUM}
**Wave**: ${WAVE_NUM}
**Test Plan Path**: /wave-tests/phase-${PHASE_NUM}/wave-${WAVE_NUM}/WAVE-TEST-PLAN.md
**Test Harness Path**: /wave-tests/phase-${PHASE_NUM}/wave-${WAVE_NUM}/WAVE-TEST-HARNESS.sh
**Demo Plan Path**: /wave-tests/phase-${PHASE_NUM}/wave-${WAVE_NUM}/WAVE-DEMO-PLAN.md
**Test Directory**: /wave-tests/phase-${PHASE_NUM}/wave-${WAVE_NUM}/tests/wave
**Test Effort Mapping Path**: /wave-tests/phase-${PHASE_NUM}/wave-${WAVE_NUM}/test-effort-mapping.json
**Created At**: [ISO-8601-timestamp]
**Created By**: code-reviewer
**Status**: active
**Test Count**: [number of test files created]
**TDD Phase**: red (tests failing, no implementation yet)
**Effort Assignments**: [mapping of efforts to their assigned tests]

ORCHESTRATOR: Please update test_plans.wave["phase${PHASE_NUM}_wave${WAVE_NUM}"] in state file with:
- test_plan_path
- test_harness_path
- demo_plan_path
- test_dir
- test_effort_mapping_path
- created_at
- test_count
- tdd_phase = "red"
- effort_assignments (which tests each effort will implement)

This enables SW Engineers to find their assigned tests per R340 pattern.
```

### 8. Commit Wave Test Artifacts
```bash
# Add all test artifacts
git add -A

# Commit with TDD message
git commit -m "tdd: wave ${WAVE_NUM} tests created BEFORE implementation

- Created wave-specific tests from architecture
- Defined test-to-effort assignments for planning
- Implemented failing tests (TDD red phase)
- Linked demo scenarios to tests (R330/R291)
- Ready for implementation planning to assign tests to efforts
- Reported test locations for state tracking (R340)

Test-Driven Development enforced per R341"

git push
```

## Exit Conditions
- All wave test deliverables created
- Tests are failing (no implementation yet)
- Test-to-effort mapping defined
- Transition to COMPLETED state
- Ready for wave implementation planning

## Success Criteria
- ✅ WAVE-TEST-PLAN.md covers wave architecture
- ✅ WAVE-TEST-HARNESS.sh is executable
- ✅ Wave test files created and failing
- ✅ WAVE-DEMO-PLAN.md links demos to tests
- ✅ Test assignments ready for effort planning
- ✅ Integration demo planned (R291)

## Related Rules
- R341: Test-Driven Development Enforcement
- R330: Demo Planning Requirements (consolidated)
- R291: Integration Demo Requirement (consolidated)
- R210: Wave Architecture Planning