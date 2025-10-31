# 🔴🔴🔴 RULE R341: Test-Driven Development Enforcement

## Classification
- **Category**: Development Process
- **Criticality Level**: 🔴🔴🔴 SUPREME
- **Enforcement**: MANDATORY at phase and wave levels
- **Penalty**: -100% for violations (AUTOMATIC FAILURE)
- **Related Rules**: R330, R291, R210, R211

## The Rule

**TEST-DRIVEN DEVELOPMENT IS MANDATORY AT PHASE AND WAVE LEVELS!**

Tests MUST be written BEFORE implementation planning begins. This ensures:
1. **Clear Success Criteria**: Tests define what "done" means
2. **Focused Implementation**: Code is written to pass specific tests
3. **Prevented Scope Creep**: Only implement what tests require
4. **Quality Assurance**: Tests exist before code, ensuring coverage

## 🔴🔴🔴 SUPREME LAW: TESTS BEFORE IMPLEMENTATION 🔴🔴🔴

**THE TDD WORKFLOW IS INVIOLABLE:**

```mermaid
graph LR
    A[Architecture] --> B[TEST CREATION]
    B --> C[Implementation Planning]
    C --> D[Implementation]
    D --> E[Tests Pass]
    
    style B fill:#ff0000,stroke:#000,stroke-width:4px
    style B color:#fff
```

**CRITICAL ENFORCEMENT POINTS:**
1. **After Phase Architecture** → MUST create phase tests FIRST
2. **After Wave Architecture** → MUST create wave tests FIRST
3. **During Effort Planning** → MUST reference existing tests
4. **During Implementation** → MUST target specific tests
5. **During Integration** → MUST validate using pre-written tests

**VIOLATION = IMMEDIATE FAILURE = -100% GRADE**

## Required State Flow

### Phase-Level TDD Flow
```
SPAWN_ARCHITECT_PHASE_PLANNING
    ↓
WAITING_FOR_ARCHITECTURE_PLAN
    ↓
SPAWN_CODE_REVIEWER_PHASE_TEST_PLANNING ← 🔴 TESTS FIRST!
    ↓
WAITING_FOR_PHASE_TEST_PLAN
    ↓
SPAWN_CODE_REVIEWER_PHASE_IMPL ← Implementation references tests
    ↓
WAITING_FOR_IMPLEMENTATION_PLAN
```

### Wave-Level TDD Flow
```
SPAWN_ARCHITECT_WAVE_PLANNING
    ↓
WAITING_FOR_ARCHITECTURE_PLAN
    ↓
SPAWN_CODE_REVIEWER_WAVE_TEST_PLANNING ← 🔴 TESTS FIRST!
    ↓
WAITING_FOR_WAVE_TEST_PLAN
    ↓
SPAWN_CODE_REVIEWER_WAVE_IMPL ← Implementation references tests
    ↓
WAITING_FOR_IMPLEMENTATION_PLAN
```

## Test Planning Requirements

### Phase Test Planning Deliverables
```
phase-tests/phase-${PHASE_NUM}/
├── PHASE-TEST-PLAN.md           # What to test
├── PHASE-TEST-HARNESS.sh        # Test runner
├── tests/
│   └── phase/
│       ├── test_core.test.js    # Actual tests
│       ├── test_api.test.py     # Multiple formats OK
│       └── test_integration.sh  # Shell tests OK
└── PHASE-DEMO-PLAN.md          # Demo scenarios (R330/R291)
```

### Wave Test Planning Deliverables
```
wave-tests/phase-${PHASE_NUM}/wave-${WAVE_NUM}/
├── WAVE-TEST-PLAN.md            # Wave-specific tests
├── WAVE-TEST-HARNESS.sh         # Wave test runner
├── tests/
│   └── wave/
│       ├── test_feature_a.test.js
│       ├── test_feature_b.test.js
│       └── test_integration.test.sh
├── WAVE-DEMO-PLAN.md           # Wave demos (R330/R291)
└── test-effort-mapping.json    # Which effort implements which tests
```

## Test Requirements

### 1. Tests Must Be Functional/Behavioral
```javascript
// ✅ CORRECT: Functional test
it('should create user and return ID', async () => {
    const result = await api.post('/users', {name: 'test'});
    expect(result.status).toBe(201);
    expect(result.body.id).toBeDefined();
});

// ❌ WRONG: Unit test (too low-level for TDD phase)
it('should call database save method', () => {
    const spy = jest.spyOn(db, 'save');
    userService.create({name: 'test'});
    expect(spy).toHaveBeenCalled();
});
```

### 2. Tests Must Fail Initially (Red Phase)
```bash
# Expected output when tests are first created:
$ ./PHASE-TEST-HARNESS.sh
🧪 PHASE 1 TEST HARNESS
================================
Mode: TDD (Tests Written First)

📋 Running Phase-Level Tests...
  Testing: test_core.test.js ... FAIL (Expected in TDD)
  Testing: test_api.test.js ... FAIL (Expected in TDD)
  Testing: test_integration.sh ... FAIL (Expected in TDD)

================================
RESULTS:
  Total Tests: 3
  Passed: 0
  Failed: 3

✅ TDD STATUS: All tests failing as expected!
   Ready for implementation planning.
```

### 3. Tests Must Guide Implementation
```markdown
## Effort Plan (MUST REFERENCE TESTS)

### Implementation Target
This effort will make the following tests pass:
- test_feature_a.test.js (from wave test plan)
- test_integration_1.test.sh (from wave test plan)

### Success Criteria
- [ ] All assigned tests pass
- [ ] No regression in existing tests
- [ ] Demo scenarios work as specified
```

## Demo Planning Consolidation

**R330 and R291 are now consolidated into test planning:**

### OLD (Separate Demo Planning)
```
Architecture → Implementation Plan → Demo Plan (R330)
                                 ↓
                        Integration Demo (R291)
```

### NEW (TDD with Integrated Demos)
```
Architecture → Test Plan (includes demos) → Implementation Plan
                    ↓
            Tests define demos
            Demos validate tests
```

**Demo requirements are now part of test planning:**
- Each test has corresponding demo scenarios
- Demo plans created alongside tests
- Integration demos validate pre-written tests

## Enforcement Protocol

### Orchestrator Enforcement
```bash
# In WAITING_FOR_ARCHITECTURE_PLAN state
validate_tdd_flow() {
    local next_state=""
    
    # Check what type of architecture we have
    if [ -f "phase-architecture.md" ]; then
        # MUST go to test planning first!
        next_state="SPAWN_CODE_REVIEWER_PHASE_TEST_PLANNING"
        echo "🔴 TDD: Creating phase tests BEFORE implementation"
    elif [ -f "wave-architecture.md" ]; then
        # MUST go to test planning first!
        next_state="SPAWN_CODE_REVIEWER_WAVE_TEST_PLANNING"
        echo "🔴 TDD: Creating wave tests BEFORE implementation"
    fi
    
    # VIOLATION CHECK
    if [[ "$next_state" == "SPAWN_CODE_REVIEWER_*_IMPL" ]]; then
        echo "❌❌❌ R341 VIOLATION: Cannot skip test planning!"
        echo "PENALTY: -100% AUTOMATIC FAILURE"
        exit 1
    fi
}
```

### Code Reviewer Enforcement
```bash
# In PHASE_IMPLEMENTATION_PLANNING or WAVE_IMPLEMENTATION_PLANNING
verify_tests_exist() {
    local test_dir=""
    
    # Determine test directory
    if [[ "$STATE" == "PHASE_IMPLEMENTATION_PLANNING" ]]; then
        test_dir="phase-tests/phase-${PHASE_NUM}"
    else
        test_dir="wave-tests/phase-${PHASE_NUM}/wave-${WAVE_NUM}"
    fi
    
    # Tests MUST exist before implementation planning
    if [ ! -f "$test_dir/TEST-PLAN.md" ]; then
        echo "❌ R341 VIOLATION: No tests found!"
        echo "Tests MUST be created BEFORE implementation planning"
        exit 1
    fi
    
    echo "✅ R341 COMPLIANT: Tests exist, referencing in plan"
}
```

### SW Engineer Enforcement
```bash
# In IMPLEMENTATION state
verify_targeting_tests() {
    # Load assigned tests from effort plan
    local assigned_tests=$(grep "Tests to Pass:" IMPLEMENTATION-PLAN.md)
    
    if [ -z "$assigned_tests" ]; then
        echo "❌ R341 VIOLATION: No test targets in plan!"
        echo "Implementation MUST target specific tests"
        exit 1
    fi
    
    echo "✅ R341 COMPLIANT: Implementing to pass assigned tests"
}
```

## Success Criteria

### Phase Level
- ✅ Phase tests created AFTER architecture, BEFORE implementation planning
- ✅ All phase architectural promises have corresponding tests
- ✅ Phase test harness is executable
- ✅ Phase demo scenarios defined alongside tests

### Wave Level
- ✅ Wave tests created AFTER architecture, BEFORE implementation planning
- ✅ Wave tests complement phase tests
- ✅ Test-to-effort assignments defined
- ✅ Wave demo scenarios defined alongside tests

### Effort Level
- ✅ Each effort references specific tests to pass
- ✅ Implementation targets assigned tests
- ✅ Success measured by test passage
- ✅ No features without corresponding tests

### Integration Level
- ✅ Integration validates using pre-written tests
- ✅ Demo scenarios prove tests pass
- ✅ No new tests during integration (all pre-written)

## Failure Conditions

### Critical Violations (-100% IMMEDIATE FAILURE)
- 🚨 Skipping test planning states
- 🚨 Implementation planning without existing tests
- 🚨 Creating tests AFTER implementation
- 🚨 Effort plans without test references
- 🚨 Integration without pre-written tests

### Major Violations (-50%)
- ⚠️ Tests not linked to architecture
- ⚠️ Demo scenarios not integrated with tests
- ⚠️ Test harness not executable
- ⚠️ Test-to-effort mapping unclear

## Examples

### ✅ CORRECT TDD Flow
```bash
# 1. Architecture complete
echo "Architecture defines WHAT to build"

# 2. Create tests FIRST
/spawn-agent code-reviewer --state PHASE_TEST_PLANNING
# Creates failing tests that define success

# 3. Create implementation plan
/spawn-agent code-reviewer --state PHASE_IMPLEMENTATION_PLANNING
# Plan references existing tests

# 4. Implement to pass tests
/spawn-agent sw-engineer --state IMPLEMENTATION
# Code written specifically to pass assigned tests

# 5. Validate with pre-written tests
./PHASE-TEST-HARNESS.sh
# Tests now pass = implementation complete
```

### ❌ WRONG: Traditional Flow
```bash
# 1. Architecture complete
# 2. Jump to implementation planning (NO TESTS!)
/spawn-agent code-reviewer --state PHASE_IMPLEMENTATION_PLANNING

# 3. Implement features
/spawn-agent sw-engineer --state IMPLEMENTATION

# 4. Write tests afterward (TOO LATE!)
# This is NOT TDD - this is a VIOLATION!
```

## Migration Notes

### For Systems Using R330/R291 Separately
1. Demo planning now happens during test planning
2. R330 requirements integrated into test plan templates
3. R291 validation happens through pre-written tests
4. No separate demo planning step needed

### State Machine Updates Required
1. Add new test planning states
2. Update transitions to enforce test-first flow
3. Update Code Reviewer states for test planning
4. Ensure orchestrator enforces new flow

## Related Rules
- **R330**: Demo Planning (now part of test planning)
- **R291**: Integration Demo (validated through tests)
- **R210**: Architecture Planning (defines what to test)
- **R211**: Implementation Planning (references tests)
- **R203**: State-Aware Startup (includes TDD states)

## Remember

**"Red, Green, Refactor"**
- **Red**: Write failing tests first
- **Green**: Implement to make tests pass
- **Refactor**: Improve code while tests stay green

**"Tests are the specification"**
**"No test, no feature"**
**"Implementation serves the tests, not vice versa"**

TEST-DRIVEN DEVELOPMENT IS NOT OPTIONAL - IT IS THE LAW!