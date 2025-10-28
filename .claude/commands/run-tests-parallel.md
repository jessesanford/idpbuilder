---
name: run-tests-parallel
description: Launch multiple SF 3.0 runtime tests in parallel as independent background processes
---

# /run-tests-parallel

Launch multiple Software Factory 3.0 runtime integration tests in parallel. Each test runs as an independent background process with tracking markers.

## Usage

/run-tests-parallel <test_numbers>

## Parameters

- **test_numbers**: Space-delimited or comma-delimited test identifiers
  - Examples: `1 2 3`, `3a 3b 4`, or `1,2,3A`
  - Supports alphanumeric IDs (case-insensitive): `3a`, `3A`, `3b`, `3B`, etc.
  - Valid numeric tests: 1.5, 2, 3, 4, 5, 6, 7
  - Note: Test ID "1" is ambiguous (matches multiple files). Use "1.5" for Wave Planning test.

## Available Tests

| Test | Name | Duration | Cost |
|------|------|----------|------|
| 1.5 | Wave Planning | 20-30 min | $5-8 |
| 2 | Wave Start → Effort Creation | 25-35 min | $6-10 |
| 3 | Effort Implementation (no review) | 35-50 min | $10-15 |
| 3a | Review Clean Path | 25-35 min | $6-10 |
| 3b | Review Fix Path | 30-45 min | $8-12 |
| 4 | Wave Integration → Fix → Re-integrate | 40-60 min | $12-18 |
| 5 | Cross-Container Cascade | 45-60 min | $15-20 |
| 6 | Project Integration | 25-35 min | $6-10 |
| 7 | PR Plan Creation | 20-30 min | $5-8 |

**Note**: Test ID "1" matches multiple files and is ambiguous. Use "1.5" for Wave Planning or check available test files.

## Execution Protocol

When invoked with test numbers, you should:

### 1. Validate and Display Test Information

```bash
cd /home/vscode/software-factory-template

# Parse test numbers (handles both space and comma delimited)
TEST_NUMS="$@"

echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║        SF 3.0 PARALLEL TEST LAUNCHER                         ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo ""
echo "📋 Tests to launch: $TEST_NUMS"
echo ""
echo "⚠️  WARNING: These are REAL AI integration tests:"
echo "   - Uses actual Claude API (costs real money)"
echo "   - Each test runs independently in background"
echo "   - Total cost depends on number of tests"
echo "   - Tests continue running even if you disconnect"
echo ""
```

### 2. Safety Warning and Countdown

```bash
echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║                    ⚠️  SAFETY WARNING ⚠️                       ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo ""
echo "This will launch multiple tests simultaneously:"
echo "  • Each test makes REAL API calls to Claude"
echo "  • Each test runs independently (parallel execution)"
echo "  • Tests create separate workspaces in /tmp/"
echo "  • Tracking markers created in /tmp/sf3-runtime-tests/"
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
echo "Starting parallel test execution NOW!"
echo ""
```

### 3. Execute Parallel Test Launcher

```bash
# Run the parallel test launcher script
bash tests/run-tests-parallel.sh $TEST_NUMS
```

### 4. Report Test Information

After the script completes, it will display:

- **Test PIDs**: Process IDs for each running test
- **Workspace Paths**: Temporary workspace directories
- **Marker Files**: JSON marker files in `/tmp/sf3-runtime-tests/`
- **Log File Paths**: All log files for monitoring
- **Monitoring Commands**: Commands to track test progress

Example output:
```
╔═══════════════════════════════════════════════════════════════╗
║           TESTS RUNNING IN BACKGROUND                        ║
╚═══════════════════════════════════════════════════════════════╝

📊 Test 1:
   PID: 12345
   Name: INIT → Phase Planning
   Workspace: /tmp/fastapi-hello-sf3-test-a1b2c3d4-e5f6-7890-abcd-ef1234567890
   Marker: /tmp/sf3-runtime-tests/test-1-20251016-143000-running.json

   📂 Log Files:
      iteration_tmp: /tmp/fastapi-hello-sf3-test-*/iteration-tmp.json
      orchestrator_stream: /tmp/fastapi-hello-sf3-test-*/orchestrator-stream.log
      orchestrator_output: /tmp/fastapi-hello-sf3-test-*/orchestrator-output.json
      orchestrator_state: /tmp/fastapi-hello-sf3-test-*/orchestrator-state-v3.json
      test_execution: /tmp/sf3-test-output-12345.log

📊 Test 2:
   PID: 12346
   Name: Wave Start → Effort Creation
   Workspace: /tmp/fastapi-hello-sf3-test-f1e2d3c4-b5a6-0987-1234-567890abcdef
   Marker: /tmp/sf3-runtime-tests/test-2-20251016-143005-running.json

   📂 Log Files:
      [... similar log file paths ...]
```

### 5. Display Monitoring Commands

After launching, display these monitoring commands:

```bash
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "📊 Monitoring Commands:"
echo ""
echo "   List all running tests:"
echo "   ls /tmp/sf3-runtime-tests/*-running.json"
echo ""
echo "   List finished tests:"
echo "   ls /tmp/sf3-runtime-tests/*-finished.json"
echo ""
echo "   View test summary:"
echo "   source tests/lib/test-marker-manager.sh && display_test_summary"
echo ""
echo "   Monitor specific test (replace <workspace> with actual path):"
echo "   tail -f <workspace>/orchestrator-stream.log"
echo ""
echo "   Check test marker details:"
echo "   cat /tmp/sf3-runtime-tests/test-<num>-*-running.json | jq"
echo ""
echo "   Quick status check for all tests:"
echo "   for f in /tmp/sf3-runtime-tests/*-running.json; do"
echo "     echo \"Test: \$(jq -r '.test_number' \$f) - PID: \$(jq -r '.pid' \$f)\""
echo "   done"
echo ""
echo "═══════════════════════════════════════════════════════════════"
```

## Marker File System

### Marker File Location

All test markers are stored in: `/tmp/sf3-runtime-tests/`

### Marker File Format

**Running Test**: `test-{NUM}-{TIMESTAMP}-running.json`
**Finished Test**: `test-{NUM}-{TIMESTAMP}-finished.json`

### Marker File Structure

```json
{
  "test_number": "1",
  "test_name": "INIT → Phase Planning",
  "test_file": "runtime-test-01-init-to-phase-planning.sh",
  "pid": 12345,
  "start_time": "2025-10-16T20:35:38Z",
  "status": "running",
  "project_prefix": "fastapi-hello-sf3-test-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "workspace": "/tmp/fastapi-hello-sf3-test-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "planning_repo": {
    "local_path": "/tmp/fastapi-hello-sf3-test-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "remote_url": "https://github.com/jessesanford/software-factory-v3-test-planning-repo",
    "branch": "fastapi-hello-sf3-test-a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  },
  "effort_repos": [
    {
      "effort_id": "phase-1-wave-1-effort-1",
      "local_path": "/tmp/fastapi-hello-sf3-test-*/efforts/phase1/wave1/effort-1",
      "remote_url": "https://github.com/org/repo",
      "branch": "effort-branch-name"
    }
  ],
  "log_files": {
    "iteration_tmp": "/tmp/fastapi-hello-sf3-test-*/iteration-tmp.json",
    "orchestrator_stream": "/tmp/fastapi-hello-sf3-test-*/orchestrator-stream.log",
    "orchestrator_output": "/tmp/fastapi-hello-sf3-test-*/orchestrator-output.json",
    "orchestrator_state": "/tmp/fastapi-hello-sf3-test-*/orchestrator-state-v3.json",
    "test_execution": "/tmp/sf3-test-output-12345.log"
  },
  "marker_file": "/tmp/sf3-runtime-tests/test-1-20251016-203538-running.json"
}
```

## Quick Status Check

To see all running and finished tests at a glance:

```bash
# Show running tests
echo "🔄 RUNNING TESTS:"
ls /tmp/sf3-runtime-tests/*-running.json 2>/dev/null | while read f; do
    TEST_NUM=$(jq -r '.test_number' "$f")
    TEST_NAME=$(jq -r '.test_name' "$f")
    PID=$(jq -r '.pid' "$f")
    echo "  Test $TEST_NUM: $TEST_NAME (PID: $PID)"
done

# Show finished tests
echo ""
echo "✅ FINISHED TESTS:"
ls /tmp/sf3-runtime-tests/*-finished.json 2>/dev/null | while read f; do
    TEST_NUM=$(jq -r '.test_number' "$f")
    TEST_NAME=$(jq -r '.test_name' "$f")
    STATUS=$(jq -r '.status' "$f")
    echo "  Test $TEST_NUM: $TEST_NAME ($STATUS)"
done
```

## Examples

### Example 1: Launch Tests 1 and 2

```
User: /run-tests-parallel 1 2

Agent:
╔═══════════════════════════════════════════════════════════════╗
║        SF 3.0 PARALLEL TEST LAUNCHER                         ║
╚═══════════════════════════════════════════════════════════════╝

📋 Tests to launch: 1 2

⚠️  WARNING: These are REAL AI integration tests:
   - Uses actual Claude API (costs real money)
   - Each test runs independently in background
   - Total cost: ~$14-22 (Test 1: $8-12, Test 2: $6-10)
   - Tests continue running even if you disconnect

[... safety warning and countdown ...]

🚀 Launching Test 1: INIT → Phase Planning
   File: tests/runtime-test-01-init-to-phase-planning.sh
   PID: 12345
   Workspace: /tmp/fastapi-hello-sf3-test-a1b2c3d4
   Marker: /tmp/sf3-runtime-tests/test-1-20251016-143000-running.json
✅ Test 1 launched successfully

🚀 Launching Test 2: Wave Start → Effort Creation
   File: tests/runtime-test-02-wave-start-to-effort-creation.sh
   PID: 12346
   Workspace: /tmp/fastapi-hello-sf3-test-f1e2d3c4
   Marker: /tmp/sf3-runtime-tests/test-2-20251016-143005-running.json
✅ Test 2 launched successfully

╔═══════════════════════════════════════════════════════════════╗
║           TESTS RUNNING IN BACKGROUND                        ║
╚═══════════════════════════════════════════════════════════════╝

[... test information display ...]

📊 Monitoring Commands:
   [... monitoring command list ...]

✅ All tests launched successfully
   Tests will continue running in background
   Use the commands above to monitor progress
```

### Example 2: Launch Multiple Tests with Comma Delimiter

```
User: /run-tests-parallel 1.5,2,3

Agent:
[Same as Example 1, but launches 3 tests instead of 2]
```

### Example 3: Launch Tests with Alphanumeric IDs (Case-Insensitive)

```
User: /run-tests-parallel 3A 3B 4

Agent:
╔═══════════════════════════════════════════════════════════════╗
║        SF 3.0 PARALLEL TEST LAUNCHER                         ║
╚═══════════════════════════════════════════════════════════════╝

📋 Tests to launch: 3A 3B 4

⚠️  WARNING: These are REAL AI integration tests:
   - Uses actual Claude API (costs real money)
   - Each test runs independently in background
   - Total cost: ~$26-40 (Test 3a: $6-10, Test 3b: $8-12, Test 4: $12-18)
   - Tests continue running even if you disconnect

[... safety warning and countdown ...]

🚀 Launching Test 3a: Review Clean Path
   File: tests/runtime-test-03a-review-clean-path.sh
   PID: 12345
   Workspace: /tmp/fastapi-hello-sf3-test-a1b2c3d4
✅ Test 3a launched successfully

🚀 Launching Test 3b: Review Fix Path
   File: tests/runtime-test-03b-review-fix-path.sh
   PID: 12346
   Workspace: /tmp/fastapi-hello-sf3-test-f1e2d3c4
✅ Test 3b launched successfully

🚀 Launching Test 4: Wave Integration → Fix → Re-integrate
   File: tests/runtime-test-04-wave-integration-fix-reintegrate.sh
   PID: 12347
   Workspace: /tmp/fastapi-hello-sf3-test-c4d5e6f7
✅ Test 4 launched successfully

[... test information display ...]
```

### Example 4: Check Status of Running Tests

```
User: Can you check the status of the running tests?

Agent:
Let me check the test markers:

bash
source tests/lib/test-marker-manager.sh && display_test_summary


═══════════════════════════════════════════════════════════════
SF 3.0 RUNTIME TESTS - STATUS SUMMARY
═══════════════════════════════════════════════════════════════

🔄 RUNNING TESTS: 2

  ✅ Test 1: INIT → Phase Planning
     PID: 12345 | Started: 2025-10-16T14:30:00Z
     Workspace: /tmp/fastapi-hello-sf3-test-a1b2c3d4
     Marker: /tmp/sf3-runtime-tests/test-1-20251016-143000-running.json

  ✅ Test 2: Wave Start → Effort Creation
     PID: 12346 | Started: 2025-10-16T14:30:05Z
     Workspace: /tmp/fastapi-hello-sf3-test-f1e2d3c4
     Marker: /tmp/sf3-runtime-tests/test-2-20251016-143005-running.json

✅ FINISHED TESTS: 0

═══════════════════════════════════════════════════════════════
```

## Integration with Individual Test Scripts

Each test script automatically:

1. **Creates marker on start**: When test begins, creates `*-running.json` marker
2. **Populates infrastructure data**: Extracts repo info from state files
3. **Renames on finish**: When test completes, renames to `*-finished.json`
4. **Updates status**: Marks test as "passed" or "failed" with exit code

## Benefits of Parallel Execution

- **Time Efficiency**: Run multiple tests simultaneously instead of sequentially
- **Independent Processes**: Each test runs in complete isolation
- **Easy Monitoring**: Quick glance at marker files shows running vs finished
- **Comprehensive Metadata**: All test info (PIDs, workspaces, logs) in one place
- **Failure Resilience**: One test failure doesn't affect others

## Related Commands

- `/run-runtime-test <num>` - Run single test
- `/run-runtime-tests` - Sequential test suite execution
- `/check-status` - System health diagnostics

## Documentation

- **Marker Manager**: `/tests/lib/test-marker-manager.sh`
- **Parallel Launcher**: `/tests/run-tests-parallel.sh`
- **Test Framework**: `/tests/runtime-test-framework.sh`
