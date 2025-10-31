---
name: run-runtime-test
description: Execute a specific SF 3.0 runtime integration test by number
---

# /run-runtime-test

Execute a single Software Factory 3.0 runtime integration test by test number.

## Usage

/run-runtime-test <test_number>

## Parameters

- **test_number**: Test number to execute (1-6, or use "1.5" for test 01.5)

## Available Tests

| Test | Name | Duration | Cost | File |
|------|------|----------|------|------|
| 1 | INIT → Phase Planning | 30-45 min | $8-12 | runtime-test-01-init-to-phase-planning.sh |
| 1.5 | Wave Planning | 20-30 min | $5-8 | runtime-test-01-5-wave-planning.sh |
| 2 | Wave Start → Effort Creation | 25-35 min | $6-10 | runtime-test-02-wave-start-to-effort-creation.sh |
| 3 | Effort Implementation → Review | 35-50 min | $10-15 | runtime-test-03-effort-implementation-to-review.sh |
| 4 | Wave Integration → Fix → Re-integrate | 40-60 min | $12-18 | runtime-test-04-wave-integration-fix-reintegrate.sh |
| 5 | Cross-Container Cascade | 45-60 min | $15-20 | runtime-test-05-cross-container-cascade.sh |
| 6 | PR Plan Creation | 20-30 min | $5-8 | runtime-test-06-pr-plan-creation.sh |

## Execution Protocol

When invoked with a test number, you should:

### 1. Parse and Validate Test Number

```bash
# Extract test number from parameters
TEST_NUM="$1"

# Validate test number
case "$TEST_NUM" in
    1|01)
        TEST_FILE="runtime-test-01-init-to-phase-planning.sh"
        TEST_NAME="INIT → Phase Planning"
        TEST_DURATION="30-45 min"
        TEST_COST="\$8-12"
        ;;
    1.5|01-5)
        TEST_FILE="runtime-test-01-5-wave-planning.sh"
        TEST_NAME="Wave Planning"
        TEST_DURATION="20-30 min"
        TEST_COST="\$5-8"
        ;;
    2|02)
        TEST_FILE="runtime-test-02-wave-start-to-effort-creation.sh"
        TEST_NAME="Wave Start → Effort Creation"
        TEST_DURATION="25-35 min"
        TEST_COST="\$6-10"
        ;;
    3|03)
        TEST_FILE="runtime-test-03-effort-implementation-to-review.sh"
        TEST_NAME="Effort Implementation → Review"
        TEST_DURATION="35-50 min"
        TEST_COST="\$10-15"
        ;;
    4|04)
        TEST_FILE="runtime-test-04-wave-integration-fix-reintegrate.sh"
        TEST_NAME="Wave Integration → Fix → Re-integrate"
        TEST_DURATION="40-60 min"
        TEST_COST="\$12-18"
        ;;
    5|05)
        TEST_FILE="runtime-test-05-cross-container-cascade.sh"
        TEST_NAME="Cross-Container Cascade"
        TEST_DURATION="45-60 min"
        TEST_COST="\$15-20"
        ;;
    6|06)
        TEST_FILE="runtime-test-06-pr-plan-creation.sh"
        TEST_NAME="PR Plan Creation"
        TEST_DURATION="20-30 min"
        TEST_COST="\$5-8"
        ;;
    *)
        echo "❌ Invalid test number: $TEST_NUM"
        echo "Valid tests: 1, 1.5, 2, 3, 4, 5, 6"
        exit 1
        ;;
esac
```

### 2. Display Test Information

Show detailed test metadata before execution:

```bash
echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║        SF 3.0 RUNTIME TEST - Test $TEST_NUM                     ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo ""
echo "📋 Test Information:"
echo "   Test Number: $TEST_NUM"
echo "   Test Name: $TEST_NAME"
echo "   Test File: $TEST_FILE"
echo ""
echo "⏱️  Estimated Duration: $TEST_DURATION"
echo "💰 Estimated Cost: $TEST_COST (Claude API tokens)"
echo ""
echo "⚠️  WARNING: This is a REAL AI integration test:"
echo "   - Uses actual Claude API (costs real money)"
echo "   - Spawns real agents (Orchestrator, State Manager, etc.)"
echo "   - Executes full workflow automation"
echo "   - Takes significant time to complete"
echo ""
```

### 3. Verify Test File Exists

```bash
cd /home/vscode/software-factory-template

if [ ! -f "tests/$TEST_FILE" ]; then
    echo "❌ ERROR: Test file not found!"
    echo "   Expected: tests/$TEST_FILE"
    echo "   PWD: $(pwd)"
    exit 1
fi

echo "✅ Test file located: tests/$TEST_FILE"
echo ""
```

### 4. Display Safety Warning

```bash
echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║                    ⚠️  SAFETY WARNING ⚠️                       ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo ""
echo "This test will:"
echo "  • Make REAL API calls to Claude (costs actual money)"
echo "  • Execute for $TEST_DURATION"
echo "  • Cost approximately $TEST_COST"
echo "  • Create temporary test workspace in /tmp/sf3-test-*"
echo "  • Spawn multiple AI agents"
echo ""
echo "Starting in 10 seconds..."
echo "Press Ctrl+C to cancel"
echo ""

# 10-second countdown
for i in {10..1}; do
    echo -n "$i... "
    sleep 1
done
echo ""
echo "Starting test execution NOW!"
echo ""
```

### 5. Execute Test in Background (Non-Blocking)

## 🚨🚨🚨 CRITICAL: ALWAYS USE run_in_background PARAMETER 🚨🚨🚨

**MANDATORY REQUIREMENT - YOU WILL BE GRADED ON THIS:**

When executing the test using the Bash tool, you **MUST** set `run_in_background: true`

### Why This Is Critical:

- Tests take **20-60 minutes** to complete (depending on test number)
- Running **synchronously BLOCKS the entire conversation** for that duration
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

**WRONG** - Never run synchronously (blocks for 20-60 minutes):
```xml
<invoke name="Bash">
  <parameter name="command">bash tests/runtime-test-01-init-to-phase-planning.sh</parameter>
</invoke>
```

### After Starting Test:

1. **Report shell ID** to user (from Bash tool response)
2. **Inform user** of expected duration (from $TEST_DURATION variable)
3. **Set up monitoring** - Use BashOutput tool periodically to check progress
4. **Remain responsive** - Continue answering user questions while test runs

### Test Execution Script:

```bash
# Record start time
START_TIME=$(date +%s)
START_TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S %Z')

echo "🚀 Starting test in background: bash tests/$TEST_FILE"
echo "   Started: $START_TIMESTAMP"
echo ""

# Run the test in background
bash "tests/$TEST_FILE" > /tmp/sf3-test-output-$$.log 2>&1 &
TEST_PID=$!
echo $TEST_PID > /tmp/sf3-test-pid-latest

echo "⏳ Waiting for test workspace to be created..."
sleep 5

# Find the workspace directory
TEST_WORKSPACE=$(find /tmp -maxdepth 1 -name "sf3-test-*" -type d -newer /tmp 2>/dev/null | tail -1)

if [ -n "$TEST_WORKSPACE" ]; then
    echo ""
    echo "╔═══════════════════════════════════════════════════════════════╗"
    echo "║           ✅ TEST RUNNING IN BACKGROUND ✅                     ║"
    echo "╚═══════════════════════════════════════════════════════════════╝"
    echo ""
    echo "📊 Test Status:"
    echo "   PID: $TEST_PID"
    echo "   Started: $START_TIMESTAMP"
    echo "   Workspace: $TEST_WORKSPACE"
    echo ""
    echo "📂 Log Files (tail these to monitor progress):"
    echo ""
    echo "   Real-time orchestrator output (streaming):"
    echo "   → $TEST_WORKSPACE/iteration-tmp.json"
    echo ""
    echo "   Completed iteration output:"
    echo "   → $TEST_WORKSPACE/iteration-output.json"
    echo ""
    echo "   Full orchestrator stream (appended after each iteration):"
    echo "   → $TEST_WORKSPACE/orchestrator-stream.log"
    echo ""
    echo "   Complete JSON output (all iterations):"
    echo "   → $TEST_WORKSPACE/orchestrator-output.json"
    echo ""
    echo "   Test execution log:"
    echo "   → /tmp/sf3-test-output-$$.log"
    echo ""
    echo "   State file (updated after each transition):"
    echo "   → $TEST_WORKSPACE/orchestrator-state-v3.json"
    echo ""
    echo "📊 Monitoring Commands:"
    echo ""
    echo "   Watch real-time orchestrator output:"
    echo "   tail -f $TEST_WORKSPACE/iteration-tmp.json | jq -r 'select(.type==\"assistant\") | .message.content[]? | select(.type==\"text\") | .text'"
    echo ""
    echo "   Watch stream log (iteration summaries):"
    echo "   tail -f $TEST_WORKSPACE/orchestrator-stream.log"
    echo ""
    echo "   Watch test execution:"
    echo "   tail -f /tmp/sf3-test-output-$$.log"
    echo ""
    echo "   Use monitoring script:"
    echo "   bash tests/runtime-test-monitor-realtime.sh $TEST_WORKSPACE"
    echo ""
    echo "🔍 Status Commands:"
    echo ""
    echo "   Check if test is running:"
    echo "   ps -p $TEST_PID > /dev/null && echo 'Running' || echo 'Completed'"
    echo ""
    echo "   Check current state:"
    echo "   jq -r '.state_machine.current_state' $TEST_WORKSPACE/orchestrator-state-v3.json"
    echo ""
    echo "   Get elapsed time:"
    echo "   echo \"\$(((\$(date +%s) - $START_TIME) / 60)) minutes elapsed\""
    echo ""
    echo "⏹️  Control Commands:"
    echo ""
    echo "   Stop test:"
    echo "   kill $TEST_PID"
    echo ""
    echo "   Check results after completion:"
    echo "   cat /tmp/sf3-test-output-$$.log | tail -50"
    echo ""
    echo "╔═══════════════════════════════════════════════════════════════╗"
    echo "║  Test will run for $TEST_DURATION (~$TEST_COST cost)           ║"
    echo "║  Monitor using the commands above                             ║"
    echo "║  I will NOT poll for updates - check manually when ready     ║"
    echo "╚═══════════════════════════════════════════════════════════════╝"
    echo ""
else
    echo "⚠️  Warning: Could not find test workspace (test may have failed during setup)"
    echo "   Test PID: $TEST_PID"
    echo "   Test output: /tmp/sf3-test-output-$$.log"
    echo "   Check if test is still running: ps -p $TEST_PID"
    echo ""
fi

# DO NOT WAIT - Exit immediately so user can monitor manually
echo "✅ Test launched successfully. Use the monitoring commands above to track progress."
echo ""
exit 0
```

### 6. Check Results Later (User-Initiated)

The test runs in background and exits immediately. When you want to check results later:

```bash
# Get the test PID
TEST_PID=$(cat /tmp/sf3-test-pid-latest)

# Check if test is still running
if ps -p $TEST_PID > /dev/null 2>&1; then
    echo "⏳ Test is still running (PID: $TEST_PID)"

    # Find workspace
    TEST_WORKSPACE=$(find /tmp -maxdepth 1 -name "sf3-test-*" -type d 2>/dev/null | tail -1)

    # Show current state
    if [ -n "$TEST_WORKSPACE" ]; then
        CURRENT_STATE=$(jq -r '.state_machine.current_state' $TEST_WORKSPACE/orchestrator-state-v3.json 2>/dev/null || echo "UNKNOWN")
        echo "   Current state: $CURRENT_STATE"

        # Show elapsed time
        START_TIME=$(stat -c %Y /tmp/sf3-test-pid-latest)
        ELAPSED=$(( $(date +%s) - START_TIME ))
        echo "   Elapsed: $(($ELAPSED / 60))m $(($ELAPSED % 60))s"
    fi
else
    echo "✅ Test completed (PID: $TEST_PID)"

    # Check exit code from test output log
    TEST_OUTPUT=$(cat /tmp/sf3-test-output-*.log 2>/dev/null | tail -100)

    if echo "$TEST_OUTPUT" | grep -q "TEST PASSED"; then
        echo "   Status: ✅ PASSED"
    elif echo "$TEST_OUTPUT" | grep -q "TEST FAILED"; then
        echo "   Status: ❌ FAILED"
    else
        echo "   Status: ❓ UNKNOWN (check logs)"
    fi

    # Show workspace location
    TEST_WORKSPACE=$(find /tmp -maxdepth 1 -name "sf3-test-*" -type d 2>/dev/null | tail -1)
    if [ -n "$TEST_WORKSPACE" ]; then
        echo "   Workspace: $TEST_WORKSPACE"
        echo "   Output: $TEST_WORKSPACE/orchestrator-output.json"
        echo "   Stream: $TEST_WORKSPACE/orchestrator-stream.log"
    fi
fi
```

## Integration Notes

### Background Execution (Non-Blocking)

**All tests run in background and exit immediately** - No polling or waiting. After starting a test:

1. **Test starts in background** - You receive PID, workspace location, and ALL log file paths
2. **Command exits immediately** - Returns control to you instantly (no 20-60 minute wait)
3. **Monitor manually** - Use provided `tail -f` commands to watch logs yourself
4. **Check when ready** - Ask for status update when you want to check progress

**Key Benefits**:
- Command returns immediately (< 10 seconds)
- No automatic polling (you control when to check)
- All log file locations displayed clearly
- Can monitor multiple logs simultaneously
- Easy to stop tests if needed (`kill <PID>`)
- Full output and artifacts preserved in workspace

**Log Files Provided**:
- `iteration-tmp.json` - Real-time streaming output (watch this for live updates)
- `orchestrator-stream.log` - Iteration summaries (appended after each iteration)
- `orchestrator-output.json` - Complete JSON output (all iterations)
- `orchestrator-state-v3.json` - Current state file
- `/tmp/sf3-test-output-$$.log` - Test execution log

### For Manual Testing

You can run individual tests without sequential execution:

```bash
# Run Test 1
/run-runtime-test 1

# Run Test 3
/run-runtime-test 3

# Run Test 1.5 (wave planning)
/run-runtime-test 1.5
```

### For Sequential Execution

Use `/run-runtime-tests` (plural) for sequential test suite execution with:
- Automatic continuation on pass
- Stop on first failure
- Test state tracking
- Resume capability

### Test Artifacts

Each test creates a temporary workspace in `/tmp/sf3-test-<timestamp>/`:
- `orchestrator-output.json`: Complete API response stream
- `orchestrator-stream.log`: Real-time execution log
- `orchestrator-state-v3.json`: Final state file
- Other state files (integration-containers.json, bug-tracking.json, etc.)

### Cost and Duration

Tests are expensive in both time and cost:
- **Shortest test**: Test 1.5 (~20 min, $5)
- **Longest test**: Test 4 or 5 (~60 min, $18-20)
- **Full suite**: 3.5-5 hours, $56-83

### Safety Features

- **10-second countdown**: Gives time to cancel
- **Clear cost warnings**: Shows estimated cost before execution
- **Duration estimates**: Sets expectations for runtime
- **Exit code propagation**: Returns test exit code for scripting

## Examples

### Example 1: Run Test 1

```
User: /run-runtime-test 1

Agent:
╔═══════════════════════════════════════════════════════════════╗
║        SF 3.0 RUNTIME TEST - Test 1                           ║
╚═══════════════════════════════════════════════════════════════╝

📋 Test Information:
   Test Number: 1
   Test Name: INIT → Phase Planning
   Test File: runtime-test-01-init-to-phase-planning.sh

⏱️  Estimated Duration: 30-45 min
💰 Estimated Cost: $8-12 (Claude API tokens)

⚠️  WARNING: This is a REAL AI integration test:
   - Uses actual Claude API (costs real money)
   - Spawns real agents (Orchestrator, State Manager, etc.)
   - Executes full workflow automation
   - Takes significant time to complete

✅ Test file located: tests/runtime-test-01-init-to-phase-planning.sh

╔═══════════════════════════════════════════════════════════════╗
║                    ⚠️  SAFETY WARNING ⚠️                       ║
╚═══════════════════════════════════════════════════════════════╝

This test will:
  • Make REAL API calls to Claude (costs actual money)
  • Execute for 30-45 min
  • Cost approximately $8-12
  • Create temporary test workspace in /tmp/sf3-test-*
  • Spawn multiple AI agents

Starting in 10 seconds...
Press Ctrl+C to cancel

10... 9... 8... 7... 6... 5... 4... 3... 2... 1...
Starting test execution NOW!

🚀 Starting test in background: bash tests/runtime-test-01-init-to-phase-planning.sh
   Started: 2025-10-11 14:30:00 UTC

⏳ Waiting for test workspace to be created...

✅ Test running in background
   PID: 12345
   Workspace: /tmp/sf3-test-20251011-143000

📊 Monitor test in real-time:
   bash tests/runtime-test-monitor-realtime.sh /tmp/sf3-test-20251011-143000

📄 Simple log monitoring:
   tail -f /tmp/sf3-test-20251011-143000/orchestrator-stream.log

🔍 Check if test is still running:
   ps -p 12345 > /dev/null && echo 'Running' || echo 'Completed'

⏹️  Stop test (if needed):
   kill 12345

⏳ Waiting for test to complete...
   (You can monitor progress in another terminal using the commands above)

[Test executes for 35 minutes in background...]

   Ended: 2025-10-11 15:05:12 UTC
   Duration: 35m 12s

╔═══════════════════════════════════════════════════════════════╗
║                  ✅ TEST PASSED ✅                             ║
╚═══════════════════════════════════════════════════════════════╝

✅ Test 1: INIT → Phase Planning - PASSED
   Duration: 35m 12s
   Exit Code: 0

📊 Test artifacts:
   Workspace: /tmp/sf3-test-20251011-143000
   Output: /tmp/sf3-test-20251011-143000/orchestrator-output.json
   Stream log: /tmp/sf3-test-20251011-143000/orchestrator-stream.log

🎉 Test completed successfully!
```

### Example 2: Run Test 5 (Failure)

```
User: /run-runtime-test 5

Agent:
╔═══════════════════════════════════════════════════════════════╗
║        SF 3.0 RUNTIME TEST - Test 5                           ║
╚═══════════════════════════════════════════════════════════════╝

📋 Test Information:
   Test Number: 5
   Test Name: Cross-Container Cascade
   Test File: runtime-test-05-cross-container-cascade.sh

⏱️  Estimated Duration: 45-60 min
💰 Estimated Cost: $15-20 (Claude API tokens)

[... safety warning and countdown ...]

🚀 Executing: bash tests/runtime-test-05-cross-container-cascade.sh
   Started: 2025-10-11 16:00:00 UTC

[Test executes for 25 minutes then fails...]

   Ended: 2025-10-11 16:25:33 UTC
   Duration: 25m 33s

╔═══════════════════════════════════════════════════════════════╗
║                  ❌ TEST FAILED ❌                             ║
╚═══════════════════════════════════════════════════════════════╝

❌ Test 5: Cross-Container Cascade - FAILED
   Duration: 25m 33s
   Exit Code: 1

🔍 Debugging Information:
   Workspace: /tmp/sf3-test-20251011-160000
   Output: /tmp/sf3-test-20251011-160000/orchestrator-output.json
   Stream log: /tmp/sf3-test-20251011-160000/orchestrator-stream.log

📄 Recent errors from output:
   State transition validation failed: CASCADE_REINTEGRATION -> INTEGRATE_WAVE_EFFORTS not allowed
   Missing transition definition in state machine
   ERROR: Orchestrator stuck in CASCADE_REINTEGRATION state

📄 Validation failures from stream log:
   ❌ Validation 2/10: Unexpected state: CASCADE_REINTEGRATION (expected progression)
   ❌ Validation 3/10: State progression FAILED (<50% coverage)

📚 Debugging Steps:
   1. Review test output in workspace directory
   2. Check orchestrator-stream.log for real-time execution log
   3. Examine orchestrator-output.json for detailed API responses
   4. Look for state transition errors or validation failures
   5. Refer to /tests/RUNTIME-TEST-GUIDE.md for troubleshooting

💡 Common Failure Causes:
   - Missing state machine transitions
   - Incorrect state-specific rules
   - State file validation errors
   - Agent spawn issues
   - Timeout (>60 minutes)
```

## Related Commands

- `/run-runtime-tests` - Sequential test execution with state tracking
- `/check-status` - System health diagnostics
- `/continue-software-factory` - Main SF 3.0 orchestrator

## Documentation

- **Test Guide**: `/tests/RUNTIME-TEST-GUIDE.md`
- **Test Framework**: `/tests/runtime-test-framework.sh`
- **Real-time Monitor**: `/tests/runtime-test-monitor-realtime.sh`
