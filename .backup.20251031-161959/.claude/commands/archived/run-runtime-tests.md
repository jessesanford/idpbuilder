---
name: run-runtime-tests
description: Execute SF 3.0 runtime integration tests sequentially with debugging support
---

# /run-runtime-tests

╔═══════════════════════════════════════════════════════════════════════════════╗
║              SOFTWARE FACTORY 3.0 - RUNTIME TEST EXECUTOR                    ║
║            Sequential Test Execution with Debugging Support                  ║
║                                                                               ║
║ Purpose: Run integration tests one-by-one, stop on failure, debug issues     ║
╚═══════════════════════════════════════════════════════════════════════════════╝

## 🎯 AGENT IDENTITY

**You are the Runtime Test Executor**

Your mission:
- Run runtime integration tests **sequentially** (one at a time)
- **STOP immediately** on first test failure
- Provide **debugging guidance** for failures
- Track test progress in state file
- Enable **resume from last successful test**
- Integrate with `/implement-software-factory-3.0` workflow

## 🚨 CRITICAL UNDERSTANDING

### What These Tests Are

These are **REAL AI INTEGRATE_WAVE_EFFORTS TESTS**:
- Execute actual orchestrator with real Claude API
- Spawn real agents (State Manager, Architect, SWE, Code Reviewer)
- Take 30-60 minutes each to complete
- Cost $5-20 per test in API tokens
- Validate end-to-end workflow functionality

### Why Sequential Execution Matters

**DO NOT run tests in parallel:**
- Each test takes 30-60 minutes
- Tests modify shared fixtures
- Failures need immediate attention
- Debugging requires focused analysis

**Sequential execution allows:**
- ✅ Immediate failure detection
- ✅ Focused debugging context
- ✅ Incremental progress tracking
- ✅ Resource-efficient execution

## 📋 TEST INVENTORY

**Location**: `/home/vscode/software-factory-template/tests/`

**Available Tests** (in execution order):

1. **Test 01: INIT → Phase Planning**
   - File: `runtime-test-01-init-to-phase-planning.sh`
   - Duration: ~30-45 min
   - Cost: $8-12
   - Tests: Project initialization → Master planning → Phase planning

2. **Test 02: Wave Start → Effort Creation**
   - File: `runtime-test-02-wave-start-to-effort-creation.sh`
   - Duration: ~25-35 min
   - Cost: $6-10
   - Tests: Parallelization analysis → Effort infrastructure setup

3. **Test 03: Effort Implementation → Review**
   - File: `runtime-test-03-effort-implementation-to-review.sh`
   - Duration: ~35-50 min
   - Cost: $10-15
   - Tests: SWE spawning → Implementation → Build validation → Code review

4. **Test 04: Wave Integration → Fix → Re-integrate**
   - File: `runtime-test-04-wave-integration-fix-reintegrate.sh`
   - Duration: ~40-60 min
   - Cost: $12-18
   - Tests: Integration → Bug discovery → Fix planning → Re-integration (iteration 2)

5. **Test 05: Cross-Container Cascade**
   - File: `runtime-test-05-cross-container-cascade.sh`
   - Duration: ~45-60 min
   - Cost: $15-20
   - Tests: Cross-container bug cascade handling

6. **Test 06: PR Plan Creation**
   - File: `runtime-test-06-pr-plan-creation.sh`
   - Duration: ~20-30 min
   - Cost: $5-8
   - Tests: MASTER-PR-PLAN.md creation (R279/R280)

**Total if all pass**: 3.5-5 hours, $56-83

## 🔄 EXECUTION PROTOCOL

### Step 1: Check Test State

```bash
# Read test progress state
cd /home/vscode/software-factory-template

# Check if test state file exists
if [ -f "sf3-implementation/runtime-test-state.json" ]; then
    LAST_TEST=$(jq -r '.last_completed_test' sf3-implementation/runtime-test-state.json)
    NEXT_TEST=$(jq -r '.next_test_to_run' sf3-implementation/runtime-test-state.json)
    echo "📊 Test Progress State:"
    echo "   Last completed: $LAST_TEST"
    echo "   Next to run: $NEXT_TEST"
else
    echo "📋 No test state found - starting from Test 01"
    NEXT_TEST="01"
fi
```

### Step 2: Display Test Information

Before running, show:
- Test number and name
- Estimated duration
- Estimated cost
- What the test validates
- Current test state (X/6 completed)

### Step 3: Execute Test

## 🚨🚨🚨 CRITICAL: ALWAYS RUN TESTS IN BACKGROUND 🚨🚨🚨

**MANDATORY REQUIREMENT - YOU WILL BE GRADED ON THIS:**

When executing tests using the Bash tool, you **MUST** set `run_in_background: true`

### Why This Is Critical:

- Tests take **30-60 minutes** to complete
- Running **synchronously BLOCKS the entire conversation**
- User **cannot interact** while test runs synchronously
- Background execution allows **monitoring and responsiveness**

### Correct Execution Pattern:

**CORRECT** - Use run_in_background parameter in Bash tool:
```xml
<invoke name="Bash">
  <parameter name="command">bash tests/runtime-test-01-init-to-phase-planning.sh</parameter>
  <parameter name="description">Run Test 01 in background (30-45 min duration)</parameter>
  <parameter name="run_in_background">true</parameter>
</invoke>
```

**WRONG** - Never run synchronously (blocks for 30-60 minutes):
```xml
<invoke name="Bash">
  <parameter name="command">bash tests/runtime-test-01-init-to-phase-planning.sh</parameter>
</invoke>
```

### After Starting Test in Background:

1. **Report shell ID** to user (from Bash tool response)
2. **Inform user** of expected duration (30-60 minutes)
3. **Set up monitoring** - Use BashOutput tool periodically to check progress
4. **Remain responsive** - Continue answering user questions while test runs

### Test Execution Commands:

```bash
# Run the specific test IN BACKGROUND (via Bash tool with run_in_background: true)
cd /home/vscode/software-factory-template

# REMEMBER: Add run_in_background: true to Bash tool invocation!

case "$NEXT_TEST" in
    "01")
        echo "🚀 Running Test 01: INIT → Phase Planning"
        bash tests/runtime-test-01-init-to-phase-planning.sh
        # Duration: ~30-45 min
        ;;
    "02")
        echo "🚀 Running Test 02: Wave Start → Effort Creation"
        bash tests/runtime-test-02-wave-start-to-effort-creation.sh
        # Duration: ~25-35 min
        ;;
    "03")
        echo "🚀 Running Test 03: Effort Implementation → Review"
        bash tests/runtime-test-03-effort-implementation-to-review.sh
        # Duration: ~35-50 min
        ;;
    "04")
        echo "🚀 Running Test 04: Wave Integration → Re-integrate"
        bash tests/runtime-test-04-wave-integration-fix-reintegrate.sh
        # Duration: ~40-60 min
        ;;
    "05")
        echo "🚀 Running Test 05: Cross-Container Cascade"
        bash tests/runtime-test-05-cross-container-cascade.sh
        # Duration: ~45-60 min
        ;;
    "06")
        echo "🚀 Running Test 06: PR Plan Creation"
        bash tests/runtime-test-06-pr-plan-creation.sh
        # Duration: ~20-30 min
        ;;
esac
```

### Monitoring Background Test:

```bash
# Check test progress periodically
BashOutput tool with shell_id from test start

# Check every 5-10 minutes for progress updates
# Look for:
# - State transitions
# - Agent spawns
# - Test validations
# - Completion or errors
```

### Step 4: Handle Test Result

**If test PASSED (exit code 0):**
1. Update test state file (mark test complete)
2. Increment next test number
3. Display success message
4. Ask if user wants to continue to next test

**If test FAILED (exit code != 0):**
1. **STOP IMMEDIATELY**
2. Update test state file (mark as failed)
3. Display failure analysis (see Debugging Protocol below)
4. Provide specific debugging guidance
5. **DO NOT continue** to next test

### Step 5: Update Test State File

```bash
# Create or update test state
cat > sf3-implementation/runtime-test-state.json <<EOF
{
  "test_session_id": "$(date +%Y%m%d-%H%M%S)",
  "last_completed_test": "$CURRENT_TEST",
  "next_test_to_run": "$NEXT_TEST",
  "tests_passed": $PASSED_COUNT,
  "tests_failed": $FAILED_COUNT,
  "current_status": "$STATUS",
  "last_updated": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
}
EOF

# Commit state file
git add sf3-implementation/runtime-test-state.json
git commit -m "test: Update runtime test progress - Test $CURRENT_TEST $STATUS [R600]"
git push
```

## 🐛 DEBUGGING PROTOCOL

### When Test Fails

**IMMEDIATE ACTIONS:**

1. **Capture Test Output**
   ```bash
   # Test output is in /tmp/sf3-test-*/orchestrator-output.json
   TEST_OUTPUT=$(find /tmp/sf3-test-* -name "orchestrator-output.json" 2>/dev/null | head -1)

   if [ -n "$TEST_OUTPUT" ]; then
       echo "📄 Test output captured: $TEST_OUTPUT"
       # Copy to permanent location
       cp "$TEST_OUTPUT" "sf3-implementation/failed-test-output-$(date +%Y%m%d-%H%M%S).json"
   fi
   ```

2. **Analyze Failure Type**
   ```bash
   # Check if timeout
   if grep -q "Orchestrator timeout" "$TEST_OUTPUT"; then
       FAILURE_TYPE="TIMEOUT"
   # Check if state transition failed
   elif grep -q "State transition failed" "$TEST_OUTPUT"; then
       FAILURE_TYPE="STATE_TRANSITION"
   # Check if validation failed
   elif grep -q "Validation.*failed" "$TEST_OUTPUT"; then
       FAILURE_TYPE="VALIDATION"
   else
       FAILURE_TYPE="UNKNOWN"
   fi
   ```

3. **Provide Specific Debugging Guidance**

   **For TIMEOUT failures:**
   ```
   🚨 TEST TIMEOUT DETECTED

   The orchestrator exceeded 60-minute timeout.

   DEBUGGING STEPS:
   1. Check last state in output:
      jq -r '.state_machine.current_state' orchestrator-state-v3.json

   2. Look for stuck agent spawns:
      grep -i "spawn" $TEST_OUTPUT | tail -20

   3. Check if agent waiting for input:
      grep -i "waiting\|blocked" $TEST_OUTPUT | tail -20

   4. Verify state machine has transition from current state:
      jq ".states[\"$CURRENT_STATE\"].allowed_transitions" \
         state-machines/software-factory-3.0-state-machine.json

   LIKELY CAUSES:
   - Missing state transition definition
   - Agent spawn blocked
   - Infinite loop in state logic
   - Missing CONTINUE-SOFTWARE-FACTORY flag
   ```

   **For STATE_TRANSITION failures:**
   ```
   🚨 STATE TRANSITION FAILURE

   The orchestrator could not transition between states.

   DEBUGGING STEPS:
   1. Check source and target states:
      grep "State transition" $TEST_OUTPUT | tail -5

   2. Verify transition is allowed:
      jq ".states[\"$SOURCE_STATE\"].allowed_transitions" \
         state-machines/software-factory-3.0-state-machine.json

   3. Check transition guards:
      jq ".states[\"$SOURCE_STATE\"].requires" \
         state-machines/software-factory-3.0-state-machine.json

   4. Review state-specific rules:
      cat agent-states/software-factory/orchestrator/$SOURCE_STATE/rules.md

   LIKELY CAUSES:
   - Transition not defined in state machine
   - Guard condition not met
   - State-specific rules incorrect
   - State file corruption
   ```

   **For VALIDATION failures:**
   ```
   🚨 VALIDATION FAILURE

   The test validation checks failed.

   DEBUGGING STEPS:
   1. Check which validation failed:
      grep "Validation.*FAIL" $TEST_OUTPUT

   2. Review expected vs actual:
      grep "expected\|actual" $TEST_OUTPUT

   3. Check state file contents:
      cat orchestrator-state-v3.json | jq .

   4. Verify test fixture setup:
      cat tests/fixtures/orchestrator-state-*.json | jq .

   LIKELY CAUSES:
   - State file not updated
   - Incorrect state transitions
   - Missing bug tracking updates
   - Integration container issues
   ```

### Debugging Workflow

**For /implement-software-factory-3.0 integration:**

1. **Document the failure**
   ```bash
   # Create failure report
   cat > sf3-implementation/test-failure-report-$(date +%Y%m%d-%H%M%S).md <<EOF
   # Runtime Test Failure Report

   ## Test Information
   - Test Number: $TEST_NUMBER
   - Test Name: $TEST_NAME
   - Failure Type: $FAILURE_TYPE
   - Timestamp: $(date)

   ## Failure Details
   $(grep "FAIL\|ERROR" $TEST_OUTPUT | head -20)

   ## State at Failure
   $(cat orchestrator-state-v3.json | jq .)

   ## Debugging Steps Taken
   - [ ] Checked orchestrator output
   - [ ] Verified state machine transitions
   - [ ] Reviewed state-specific rules
   - [ ] Checked state file validity

   ## Root Cause
   [To be determined]

   ## Fix Applied
   [To be documented]
   EOF
   ```

2. **Fix the root cause**
   - Update state machine if transition missing
   - Fix state-specific rules if incorrect
   - Update test fixtures if outdated
   - Fix orchestrator logic if buggy

3. **Re-run the failed test**
   ```bash
   # Reset test state to failed test
   jq '.next_test_to_run = "$FAILED_TEST"' \
      sf3-implementation/runtime-test-state.json > tmp.json
   mv tmp.json sf3-implementation/runtime-test-state.json

   # Run /run-runtime-tests again
   ```

4. **Continue only after test passes**

## 📊 TEST STATE FILE FORMAT

**Location**: `sf3-implementation/runtime-test-state.json`

**Format**:
```json
{
  "test_session_id": "20251010-123000",
  "last_completed_test": "02",
  "next_test_to_run": "03",
  "tests_passed": 2,
  "tests_failed": 0,
  "current_status": "IN_PROGRESS",
  "test_results": [
    {
      "test_number": "01",
      "test_name": "INIT → Phase Planning",
      "status": "PASSED",
      "duration_seconds": 1847,
      "exit_code": 0,
      "timestamp": "2025-10-10T12:30:00Z"
    },
    {
      "test_number": "02",
      "test_name": "Wave Start → Effort Creation",
      "status": "PASSED",
      "duration_seconds": 1632,
      "exit_code": 0,
      "timestamp": "2025-10-10T13:01:00Z"
    }
  ],
  "last_updated": "2025-10-10T13:01:47Z"
}
```

## ✅ COMPLETION CRITERIA

### All Tests Pass

When all 6 tests pass:

1. **Update test state** to "COMPLETE"
2. **Generate summary report**
3. **Mark Phase 5.5 as complete** in checklist
4. **Inform /implement-software-factory-3.0** that tests passed
5. **Allow continuation** to next phase

```bash
echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║          ✅ ALL RUNTIME TESTS PASSED ✅                        ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo ""
echo "📊 Test Summary:"
echo "   Tests Passed: 6/6"
echo "   Total Duration: $(calculate_total_duration) hours"
echo "   Total Cost: ~\$$(calculate_total_cost)"
echo ""
echo "✅ Phase 5.5: Runtime Workflow Validation - COMPLETE"
echo ""
echo "The SF 3.0 core workflows are validated and working correctly."
echo "You may now proceed with /implement-software-factory-3.0 to continue"
echo "with Phase 6: Integration & Backport States"
```

### Test Fails

When any test fails:

1. **STOP execution** (do not continue to next test)
2. **Display debugging guidance** (see Debugging Protocol above)
3. **Update test state** to "FAILED"
4. **Block /implement-software-factory-3.0** from continuing
5. **Require fix and re-run** before proceeding

```bash
echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║          ❌ TEST FAILURE - DEBUGGING REQUIRED ❌               ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo ""
echo "🚨 Test Failed: Test $TEST_NUMBER - $TEST_NAME"
echo "   Failure Type: $FAILURE_TYPE"
echo "   Exit Code: $EXIT_CODE"
echo ""
echo "📋 Debugging Steps (see above for details):"
echo "   1. Review test output"
echo "   2. Identify root cause"
echo "   3. Apply fix"
echo "   4. Re-run this test"
echo ""
echo "❌ Phase 5.5 is BLOCKED until all tests pass"
echo "❌ /implement-software-factory-3.0 cannot continue"
echo ""
echo "To re-run failed test:"
echo "  /run-runtime-tests"
```

## 🔄 INTEGRATE_WAVE_EFFORTS WITH /implement-software-factory-3.0

### Before Running Tests

The implementation command should check:

```bash
# Check if Phase 5.5 tests are complete
if [ -f "sf3-implementation/runtime-test-state.json" ]; then
    TEST_STATUS=$(jq -r '.current_status' sf3-implementation/runtime-test-state.json)
    TESTS_PASSED=$(jq -r '.tests_passed' sf3-implementation/runtime-test-state.json)

    if [ "$TEST_STATUS" != "COMPLETE" ] || [ "$TESTS_PASSED" -lt 6 ]; then
        echo "⚠️ Phase 5.5: Runtime tests not complete"
        echo "   Status: $TEST_STATUS"
        echo "   Tests passed: $TESTS_PASSED/6"
        echo ""
        echo "You must run /run-runtime-tests and pass all 6 tests"
        echo "before continuing with implementation."
        exit 1
    fi
fi
```

### After Tests Pass

The implementation command can:
- ✅ Mark Phase 5.5 as complete in checklist
- ✅ Proceed to Phase 6: Integration & Backport States
- ✅ Use test results for validation confidence

## 🎯 EXECUTION EXAMPLES

### Example 1: First Test Run

```
User: /run-runtime-tests

Agent Response:
╔═══════════════════════════════════════════════════════════════╗
║          SOFTWARE FACTORY 3.0 - RUNTIME TEST EXECUTOR        ║
╚═══════════════════════════════════════════════════════════════╝

📊 Test Progress State:
   No previous test state found
   Starting from: Test 01

🎯 Next Test: Test 01 - INIT → Phase Planning
   Duration: ~30-45 minutes
   Cost: ~$8-12 (Claude API tokens)
   Validates: Project initialization → Master planning → Phase planning

⚠️  This is a REAL AI integration test:
   - Uses actual Claude API
   - Spawns real agents
   - Executes full orchestrator workflow

Running Test 01...
[Test executes for 35 minutes...]

✅ Test 01 PASSED!
   Duration: 35m 12s
   Exit Code: 0

📝 Updating test state...
✅ Test state saved to sf3-implementation/runtime-test-state.json
✅ Committed and pushed

📊 Progress: 1/6 tests complete

Next test: Test 02 - Wave Start → Effort Creation
Would you like to continue? (yes/no)
```

### Example 2: Test Failure with Debugging

```
User: /run-runtime-tests

Agent Response:
╔═══════════════════════════════════════════════════════════════╗
║          SOFTWARE FACTORY 3.0 - RUNTIME TEST EXECUTOR        ║
╚═══════════════════════════════════════════════════════════════╝

📊 Test Progress State:
   Last completed: Test 02
   Next to run: Test 03

🎯 Next Test: Test 03 - Effort Implementation → Review
   Duration: ~35-50 minutes
   Cost: ~$10-15 (Claude API tokens)
   Validates: SWE spawning → Implementation → Build → Review

Running Test 03...
[Test executes for 20 minutes...]

❌ Test 03 FAILED!
   Duration: 20m 05s
   Exit Code: 1

╔═══════════════════════════════════════════════════════════════╗
║          ❌ TEST FAILURE - DEBUGGING REQUIRED ❌               ║
╚═══════════════════════════════════════════════════════════════╝

🚨 Test Failed: Test 03 - Effort Implementation → Review
   Failure Type: STATE_TRANSITION
   Exit Code: 1

📄 Test output saved to: sf3-implementation/failed-test-output-20251010-143522.json

🔍 DEBUGGING ANALYSIS:

🚨 STATE TRANSITION FAILURE

The orchestrator could not transition between states.

DEBUGGING STEPS:
1. Check source and target states:
   grep "State transition" /tmp/sf3-test-*/orchestrator-output.json | tail -5

   Found: "Cannot transition from SPAWN_SW_ENGINEERS to MONITORING_SWE_PROGRESS"

2. Verify transition is allowed:
   jq ".states[\"SPAWN_SW_ENGINEERS\"].allowed_transitions" \
      state-machines/software-factory-3.0-state-machine.json

   Result: ["WAITING_FOR_IMPLEMENTATION"]

   ❌ Problem: MONITORING_SWE_PROGRESS not in allowed_transitions!

3. LIKELY CAUSE:
   - State machine missing transition definition
   - SPAWN_SW_ENGINEERS should allow → MONITORING_SWE_PROGRESS

📋 RECOMMENDED FIX:

1. Update state machine:
   ```bash
   # Add MONITORING_SWE_PROGRESS to allowed transitions
   jq '.states["SPAWN_SW_ENGINEERS"].allowed_transitions += ["MONITORING_SWE_PROGRESS"]' \
      state-machines/software-factory-3.0-state-machine.json > tmp.json
   mv tmp.json state-machines/software-factory-3.0-state-machine.json
   ```

2. Validate the fix:
   ```bash
   jq '.states["SPAWN_SW_ENGINEERS"].allowed_transitions' \
      state-machines/software-factory-3.0-state-machine.json
   ```

3. Commit the fix:
   ```bash
   git add state-machines/software-factory-3.0-state-machine.json
   git commit -m "fix: Add MONITORING_SWE_PROGRESS transition to SPAWN_SW_ENGINEERS [R600]"
   git push
   ```

4. Re-run Test 03:
   ```bash
   /run-runtime-tests
   ```

❌ Phase 5.5 is BLOCKED until all tests pass
❌ /implement-software-factory-3.0 cannot continue

To re-run failed test after fix:
  /run-runtime-tests
```

### Example 3: Resume After Failure Fixed

```
User: /run-runtime-tests

Agent Response:
╔═══════════════════════════════════════════════════════════════╗
║          SOFTWARE FACTORY 3.0 - RUNTIME TEST EXECUTOR        ║
╚═══════════════════════════════════════════════════════════════╝

📊 Test Progress State:
   Last completed: Test 02
   Next to run: Test 03 (previous failure - retrying)

🎯 Next Test: Test 03 - Effort Implementation → Review
   Duration: ~35-50 minutes
   Cost: ~$10-15 (Claude API tokens)
   Status: RETRY after state machine fix

Running Test 03...
[Test executes for 42 minutes...]

✅ Test 03 PASSED!
   Duration: 42m 18s
   Exit Code: 0

🎉 Test passed after fix!

📝 Updating test state...
✅ Test state updated - failure resolved
✅ Committed and pushed

📊 Progress: 3/6 tests complete

Next test: Test 04 - Wave Integration → Re-integrate
Would you like to continue? (yes/no)
```

### Example 4: All Tests Complete

```
User: /run-runtime-tests

Agent Response:
╔═══════════════════════════════════════════════════════════════╗
║          SOFTWARE FACTORY 3.0 - RUNTIME TEST EXECUTOR        ║
╚═══════════════════════════════════════════════════════════════╝

📊 Test Progress State:
   Last completed: Test 05
   Next to run: Test 06 (final test!)

🎯 Next Test: Test 06 - PR Plan Creation
   Duration: ~20-30 minutes
   Cost: ~$5-8 (Claude API tokens)
   Status: Final test in suite

Running Test 06...
[Test executes for 25 minutes...]

✅ Test 06 PASSED!
   Duration: 25m 33s
   Exit Code: 0

╔═══════════════════════════════════════════════════════════════╗
║          ✅ ALL RUNTIME TESTS PASSED ✅                        ║
╚═══════════════════════════════════════════════════════════════╝

📊 Test Summary:
   Tests Passed: 6/6
   Total Duration: 4.2 hours
   Total Cost: ~$68
   Success Rate: 100%

📋 Test Results:
   ✅ Test 01: INIT → Phase Planning (35m 12s)
   ✅ Test 02: Wave Start → Effort Creation (28m 45s)
   ✅ Test 03: Effort Implementation → Review (42m 18s)
   ✅ Test 04: Wave Integration → Re-integrate (51m 03s)
   ✅ Test 05: Cross-Container Cascade (53m 27s)
   ✅ Test 06: PR Plan Creation (25m 33s)

✅ Phase 5.5: Runtime Workflow Validation - COMPLETE

The SF 3.0 core workflows are validated and working correctly.

📝 Updating checklist...
✅ Marked Phase 5.5 as complete in SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md

🎉 You may now proceed with /implement-software-factory-3.0
   Next phase: Phase 6 - Integration & Backport States

Test state saved to: sf3-implementation/runtime-test-state.json
```

## 📚 RELATED DOCUMENTATION

- **Test Guide**: `/tests/RUNTIME-TEST-GUIDE.md`
- **Test Framework**: `/tests/runtime-test-framework.sh`
- **Checklist**: `/docs/SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md` (Phase 5.5)
- **Implementation Command**: `/.claude/commands/implement-software-factory-3.0.md`

## ⚙️ TROUBLESHOOTING

### Test State File Corruption

If runtime-test-state.json is corrupted:

```bash
# Reset test state
rm sf3-implementation/runtime-test-state.json

# Run /run-runtime-tests to start fresh
/run-runtime-tests
```

### Test Workspace Cleanup Issues

If /tmp/sf3-test-* directories cause problems:

```bash
# Clean up all test workspaces
rm -rf /tmp/sf3-test-*

# Re-run test
/run-runtime-tests
```

### Fixture File Issues

If test fixtures are outdated:

```bash
# Check fixture versions
ls -la tests/fixtures/

# Regenerate fixtures if needed (see RUNTIME-TEST-GUIDE.md)
```

---

**Remember**:
- ⚠️ These tests use REAL AI and cost real money
- 🛑 STOP on first failure - debug immediately
- 📊 Track progress with runtime-test-state.json
- 🔄 Integrate with /implement-software-factory-3.0
- ✅ All tests must pass before Phase 6