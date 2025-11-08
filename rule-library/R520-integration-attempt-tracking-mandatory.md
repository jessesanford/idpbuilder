# 🚨🚨🚨 RULE R520: MANDATORY INTEGRATE_WAVE_EFFORTS ATTEMPT TRACKING 🚨🚨🚨

**Criticality:** BLOCKING
**Applies To:** Orchestrator, Integration Agents, CASCADE Agents
**Enforcement:** Pre-spawn validation, schema validation
**Created:** 2025-10-06
**Reason:** Prevent infinite integration loops caused by binary state tracking

---

## PROBLEM STATEMENT

**The 6-Hour Loop Bug:**

Prior to R520, orchestrator-state-v3.json tracked integration status as binary:
- `integration_complete: true` or `false`

This caused CRITICAL FAILURES:
1. **Cannot distinguish states:**
   - "Never attempted integration" → SHOULD spawn agent
   - "Attempted, failed, fixes in progress" → SHOULD wait
   - "Attempted, failed, fixes complete" → SHOULD retry
   - ALL THREE LOOKED IDENTICAL: `integration_complete: false`

2. **Real-world impact:**
   - idpbuilder-push-oci stuck in 6+ hour loop
   - Orchestrator repeatedly spawned integration agents
   - Integration already attempted, failed, CASCADE running
   - No metadata to detect "already attempted, waiting for fixes"

3. **Root cause:** Tracking OUTCOMES not PROCESS

---

## RULE REQUIREMENTS

### 🚨 MANDATORY SCHEMA FIELDS 🚨

ALL wave_integrations, phase_integrations, and project_integration objects MUST have:

#### Minimal Loop Prevention Fields (REQUIRED)
```json
{
  "last_integration_attempt": "2025-10-06T12:00:00Z" | null,
  "merge_commit_sha": "abc123def" | null,
  "last_attempt_result": "BUILD_GATE_FAILURE" | "PROJECT_DONE" | "MERGE_FAILED" | null,
  "bugs_pending_fix": ["BUG-001", "BUG-002"],
  "ready_for_retry": false,
  "retry_reason": "Waiting for CASCADE fixes to complete"
}
```

#### Progress Tracking Fields (REQUIRED)
```json
{
  "integration_started_at": "2025-10-06T12:00:00Z" | null,
  "integration_duration_seconds": 7200 | null,
  "attempt_count": 1,
  "max_attempts": 5,
  "time_limit_hours": 4,
  "last_progress_check": "2025-10-06T14:00:00Z" | null,
  "progress_warnings": ["Integration exceeds 4 hour limit"]
}
```

#### Decision Support Fields (REQUIRED)
```json
{
  "next_action": "WAIT_FOR_CASCADE_FIXES" | "RETRY_INTEGRATE_WAVE_EFFORTS" | "MANUAL_INTERVENTION" | null,
  "manual_override": "Integration attempted but failed. CASCADE running."
}
```

#### Audit Trail (STRONGLY RECOMMENDED)
```json
{
  "integration_attempts": [
    {
      "attempt_number": 1,
      "started_at": "2025-10-06T12:00:00Z",
      "completed_at": "2025-10-06T12:05:00Z",
      "merge_status": "PROJECT_DONE",
      "merge_commit_sha": "abc123",
      "build_status": "FAILED",
      "test_status": "NOT_RUN",
      "demo_status": "NOT_RUN",
      "overall_result": "BUILD_GATE_FAILURE",
      "bugs_created": ["BUG-001", "BUG-002"],
      "cascade_triggered": true
    }
  ]
}
```

---

## ENFORCEMENT RULES

### 🚨 Rule R520.1: Pre-Spawn Validation (BLOCKING)

Before spawning integration agent, orchestrator MUST validate:

```bash
# VIOLATION: Spawn when integration_complete = true
if [ "$integration.integration_complete" == "true" ]; then
  echo "❌ R520.1 VIOLATION: Cannot spawn - already complete"
  exit 1
fi

# VIOLATION: Spawn when attempt in progress
if [ "$integration.last_attempt_result" == "IN_PROGRESS" ]; then
  echo "❌ R520.1 VIOLATION: Cannot spawn - attempt in progress"
  exit 1
fi

# VIOLATION: Spawn when not ready for retry
if [ "$integration.last_attempt_result" != "null" ] && \
   [ "$integration.last_attempt_result" != "PROJECT_DONE" ] && \
   [ "$integration.ready_for_retry" != "true" ]; then
  echo "❌ R520.1 VIOLATION: Cannot spawn - not ready for retry"
  echo "   Reason: $integration.retry_reason"
  echo "   Next action: $integration.next_action"
  exit 1
fi
```

**PENALTY:** -100% (System corruption - infinite loops)

### 🚨 Rule R520.2: Attempt Counter Enforcement (BLOCKING)

```bash
# VIOLATION: Exceed max attempts
if [ $integration.attempt_count -ge $integration.max_attempts ]; then
  echo "❌ R520.2 VIOLATION: Exceeded max attempts"
  echo "   Attempts: $integration.attempt_count / $integration.max_attempts"
  echo "   REQUIRED ACTION: ESCALATE to manual intervention"
  exit 1
fi
```

**PENALTY:** -100% (Prevents infinite retry loops)

### 🚨 Rule R520.3: Time Limit Enforcement (BLOCKING)

```bash
# VIOLATION: Exceed time limit
duration_hours=$(awk "BEGIN {printf \"%.2f\", $integration.integration_duration_seconds / 3600}")
if (( $(awk "BEGIN {print ($duration_hours >= $integration.time_limit_hours)}") )); then
  echo "❌ R520.3 VIOLATION: Exceeded time limit"
  echo "   Duration: ${duration_hours}h / $integration.time_limit_hours h"
  echo "   REQUIRED ACTION: ESCALATE to manual intervention"
  exit 1
fi
```

**PENALTY:** -100% (Prevents hanging integrations)

### ⚠️⚠️⚠️ Rule R520.4: Metadata Update (WARNING) ⚠️⚠️⚠️

Integration agent MUST update ALL fields on completion:

**On Success:**
```json
{
  "last_attempt_result": "PROJECT_DONE",
  "integration_complete": true,
  "ready_for_retry": false,
  "next_action": "COMPLETE"
}
```

**On Failure:**
```json
{
  "last_attempt_result": "BUILD_GATE_FAILURE",
  "integration_complete": false,
  "ready_for_retry": false,
  "next_action": "WAIT_FOR_CASCADE_FIXES",
  "retry_reason": "Waiting for fixes to bugs [...]",
  "bugs_pending_fix": ["BUG-001", "BUG-002"]
}
```

**PENALTY:** -50% (Incomplete state tracking)

### ⚠️⚠️⚠️ Rule R520.5: CASCADE Update (WARNING) ⚠️⚠️⚠️

CASCADE agent MUST update integration metadata after fixes:

**After Effort Fixes:**
```json
{
  "bugs_pending_fix": [],
  "associated_bugs": {
    "active": [],
    "fixed": ["BUG-001", "BUG-002"],
    "cascade_pending": true
  },
  "ready_for_retry": true,
  "retry_reason": "All bugs fixed. CASCADE reintegration in progress.",
  "next_action": "WAIT_FOR_CASCADE_FIXES"
}
```

**After CASCADE Reintegration:**
```json
{
  "associated_bugs": {
    "cascade_pending": false
  },
  "ready_for_retry": true,
  "retry_reason": "CASCADE complete. All fixes integrated. Ready to retry.",
  "next_action": "RETRY_INTEGRATE_WAVE_EFFORTS"
}
```

**PENALTY:** -50% (Deadlock - integration never retries)

---

## DECISION LOGIC

The orchestrator MUST use this decision tree:

```
1. Check: integration.last_integration_attempt
   NULL? → SPAWN_INTEGRATION_AGENT (first attempt)
   NOT NULL? → Continue to step 2

2. Check: integration.last_attempt_result
   "PROJECT_DONE"? → COMPLETE (integration done)
   "IN_PROGRESS"? → WAIT (agent still running)
   "MERGE_FAILED"? → MANUAL_INTERVENTION (conflicts)
   Other failure? → Continue to step 3

3. Check: integration.ready_for_retry
   FALSE? → Execute integration.next_action (usually WAIT_FOR_CASCADE_FIXES)
   TRUE? → Continue to step 4

4. Validate limits:
   attempt_count >= max_attempts? → ESCALATE
   duration >= time_limit_hours? → ESCALATE
   ELSE? → SPAWN_INTEGRATION_AGENT (retry)
```

**See:** `/home/vscode/software-factory-template/INTEGRATE_WAVE_EFFORTS-LOOP-PREVENTION.md` for complete logic

---

## MONITORING_SWE_PROGRESS REQUIREMENTS

### 🚨 Mandatory Progress Checks

1. **On every state transition:**
   - Update `last_progress_check`
   - Recalculate `integration_duration_seconds`
   - Check time/attempt limits
   - Add warnings if limits exceeded

2. **Every 30 minutes during integration:**
   - Check for stuck state (no progress > 2 hours)
   - Verify CASCADE status if `next_action = "WAIT_FOR_CASCADE_FIXES"`
   - Detect deadlocks (ready_for_retry but CASCADE not running)

3. **On integration agent spawn:**
   - Increment `attempt_count`
   - Create new entry in `integration_attempts[]`
   - Set `last_attempt_result = "IN_PROGRESS"`

4. **On integration agent completion:**
   - Update attempt record with full details
   - Update summary fields (`last_attempt_result`, etc.)
   - Set `ready_for_retry` and `next_action` based on outcome

---

## VALIDATION QUERIES

### Detect Missing Metadata
```bash
jq '.pre_planned_infrastructure.integrations.wave_integrations | to_entries[] |
  select(.value.last_integration_attempt == null and .value.integration_complete == null) |
  {id: .key, error: "Missing R520 required fields"}
' orchestrator-state-v3.json
```

### Find Stuck Integrations
```bash
jq --arg now "$(date -u +%Y-%m-%dT%H:%M:%SZ)" '
  .pre_planned_infrastructure.integrations.wave_integrations | to_entries[] |
  select(.value.integration_complete == false and
         .value.last_integration_attempt != null and
         ((now|fromdateiso8601) - (.value.last_integration_attempt|fromdateiso8601)) > 7200) |
  {id: .key, stuck_hours: (((now|fromdateiso8601) - (.value.last_integration_attempt|fromdateiso8601))/3600)}
' orchestrator-state-v3.json
```

### Check for Infinite Loops
```bash
jq '.pre_planned_infrastructure.integrations.wave_integrations | to_entries[] |
  select(.value.attempt_count >= 3 and .value.integration_complete == false) |
  {id: .key, attempts: .value.attempt_count, results: [.value.integration_attempts[].overall_result]}
' orchestrator-state-v3.json
```

---

## MIGRATION FOR EXISTING PROJECTS

### If Your State File Lacks R520 Fields

**Option 1: Add Minimal Fields (5 minutes)**
```bash
jq '.pre_planned_infrastructure.integrations.wave_integrations |=
  with_entries(.value += {
    "last_integration_attempt": null,
    "last_attempt_result": null,
    "ready_for_retry": false,
    "next_action": null,
    "attempt_count": 0,
    "max_attempts": 5,
    "time_limit_hours": 4,
    "bugs_pending_fix": [],
    "integration_attempts": []
  })
' orchestrator-state-v3.json > orchestrator-state-updated.json
```

**Option 2: Infer from Git History (30 minutes)**
- Parse git log for integration merge commits
- Populate `integration_attempts[]` with historical data
- Set current state based on latest attempt

**Option 3: Manual Override (1 minute)**
- Add `manual_override` field documenting current state
- Orchestrator can read this to understand context
- Example: "Integration attempted on Oct 5. Merge succeeded but build failed. CASCADE in progress. DO NOT retry until CASCADE updates ready_for_retry."

---

## CRITICAL REMINDERS

1. **Binary state is INSUFFICIENT:**
   - `integration_complete: true/false` alone → LOOPS
   - MUST track attempt history and readiness

2. **ready_for_retry is a GATE:**
   - `false` → DO NOT spawn agent, wait for condition
   - `true` → Check limits, then spawn

3. **next_action is GUIDANCE:**
   - Tells orchestrator what to do when `ready_for_retry = false`
   - Examples: WAIT_FOR_CASCADE_FIXES, MANUAL_INTERVENTION, ESCALATE

4. **Limits prevent infinite loops:**
   - max_attempts: Don't retry forever
   - time_limit_hours: Don't wait forever
   - Both exceeded → ESCALATE to human

5. **Progress monitoring is MANDATORY:**
   - Without monitoring, can't detect stuck states
   - Update `last_progress_check` at EVERY checkpoint

---

## REFERENCES

- **Schema:** `/home/vscode/software-factory-template/orchestrator-state.schema.json`
- **Documentation:** `/home/vscode/software-factory-template/INTEGRATE_WAVE_EFFORTS-LOOP-PREVENTION.md`
- **Real-world bug:** idpbuilder-push-oci 6-hour loop (2025-10-05)
- **Related Rules:** R291 (build gates), R327 (CASCADE), R500 (HEAD tracking)

---

**ABSOLUTELY CRITICAL:** Every integration MUST use R520 fields. The old binary state WILL cause loops and waste hours of compute time. This is NON-NEGOTIABLE.

## State Manager Coordination (SF 3.0)

State Manager tracks integration attempts through iteration containers:
- **integration-containers.json** `.active_integrations[].current_iteration` field
- **Each integration attempt** increments iteration counter
- **Convergence tracking**: Iterations continue until bugs_found = 0
- **Atomic updates**: Integration attempt + bug count + iteration number updated together

Integration iteration is the core convergence mechanism in SF 3.0.

See: `integration-containers.json`, `agent-states/software-factory/orchestrator/START_WAVE_ITERATION/rules.md`, R327 (iteration container management)
