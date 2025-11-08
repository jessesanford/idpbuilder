# 🚨🚨🚨 CRITICAL RULE R530: State Transition Cycle Detection

**Criticality**: BLOCKING
**Category**: State Machine Validation / Loop Prevention
**Applies To**: State Manager Agent (all factory types)
**Performance**: O(1) time and space complexity

## Purpose

Detect and prevent infinite state transition loops that waste API costs and block progress. This rule provides a **bullet-proof, performant** mechanism to identify when the system is cycling through states without making genuine progress.

## The Problem

Without loop detection, the Software Factory can enter infinite cycles:
- **Ping-pong loops**: A → B → A → B → A (2-state cycles)
- **Same-state loops**: A → A → A (single-state cycles)
- **ERROR_RECOVERY escalation**: Multiple states failing → ERROR_RECOVERY → failing again
- **Complex cycles**: A → B → C → A (N-state cycles, harder to detect)

These loops consume:
- API tokens (expensive)
- Time (blocking development)
- Context windows (filling with repeated failures)

## Detection Thresholds

### ERROR_RECOVERY Triggers (Force immediate halt)

| Pattern | Threshold | Description |
|---------|-----------|-------------|
| **Ping-pong** | 5 cycles | A ↔ B ↔ A ↔ B ↔ A ↔ B (same 2 states alternating) |
| **Same-state** | 3 entries | A → A → A (re-entering same state) |
| **ERROR_RECOVERY** | 3 entries | Entering ERROR_RECOVERY state 3+ times |

### Why These Thresholds?

- **Ping-pong (5)**: Legitimate back-and-forth (e.g., validation failures) may happen 2-3 times. 5+ is definitely a loop.
- **Same-state (3)**: Entering the same state 3+ times consecutively indicates broken state logic.
- **ERROR_RECOVERY (3)**: If error recovery fails 3 times, manual intervention is required.

## Counter Management

### Loop Detection State Fields

Stored in `orchestrator-state-v3.json` under `state_machine.loop_detection`:

```json
{
  "loop_detection": {
    "last_two_states": ["STATE_A", "STATE_B"],    // Ring buffer (max 2)
    "ping_pong_count": 2,                         // Increments on ping-pong pattern
    "same_state_count": 1,                        // Increments on same-state re-entry
    "error_recovery_entries": 0,                  // Counts ERROR_RECOVERY entries
    "last_progress_timestamp": "2025-10-15T12:30:00Z"  // Last genuine progress
  }
}
```

### Counter Update Rules

#### 1. Ping-Pong Detection

```bash
# Check if current_state matches last_two_states[0] (ping-pong pattern)
if [[ "$current_state" == "${last_two_states[0]}" ]]; then
    ping_pong_count=$((ping_pong_count + 1))

    if [[ $ping_pong_count -ge 5 ]]; then
        echo "🚨 LOOP DETECTED: Ping-pong threshold exceeded ($ping_pong_count cycles)"
        force_state="ERROR_RECOVERY"
        loop_reason="Ping-pong loop: alternating between ${last_two_states[@]}"
    fi
fi
```

#### 2. Same-State Detection

```bash
# Check if current_state matches last state (same-state re-entry)
if [[ "$current_state" == "${last_two_states[-1]}" ]]; then
    same_state_count=$((same_state_count + 1))

    if [[ $same_state_count -ge 3 ]]; then
        echo "🚨 LOOP DETECTED: Same-state threshold exceeded ($same_state_count re-entries)"
        force_state="ERROR_RECOVERY"
        loop_reason="Same-state loop: stuck in $current_state"
    fi
fi
```

#### 3. ERROR_RECOVERY Escalation

```bash
# Count ERROR_RECOVERY entries
if [[ "$current_state" == "ERROR_RECOVERY" ]]; then
    error_recovery_entries=$((error_recovery_entries + 1))

    if [[ $error_recovery_entries -ge 3 ]]; then
        echo "🚨 LOOP DETECTED: ERROR_RECOVERY threshold exceeded ($error_recovery_entries entries)"
        force_state="ERROR_RECOVERY"
        loop_reason="ERROR_RECOVERY escalation: recovery failed $error_recovery_entries times"
    fi
fi
```

#### 4. Ring Buffer Update

```bash
# Update last_two_states ring buffer
last_two_states=("${last_two_states[-1]}" "$current_state")

# Trim to max 2 entries
if [[ ${#last_two_states[@]} -gt 2 ]]; then
    last_two_states=("${last_two_states[@]:1}")
fi
```

### Counter Reset (Genuine Progress Detection)

**Reset ALL counters to 0** when ANY of these occur:

1. **New effort created** (`pre_planned_infrastructure.efforts[].created == true`)
2. **New effort validated** (`pre_planned_infrastructure.efforts[].validated == true`)
3. **Integration completed** (wave/phase/project integration branch created)
4. **New agent spawned** (SWE, Code Reviewer, Architect)
5. **State NOT in last_two_states** (breaking the pattern)

```bash
# Example: Check for genuine progress
progress_detected=false

# Check for new effort creation
if jq -e '.pre_planned_infrastructure.efforts[] | select(.created == true and .validated == false)' "$STATE_FILE"; then
    progress_detected=true
fi

# Check for completed integration
if [[ "$current_state" == "WAVE_COMPLETE" ]] && jq -e '.integration_branches.wave_integrations | length > 0' "$STATE_FILE"; then
    progress_detected=true
fi

# Reset counters on progress
if [[ "$progress_detected" == "true" ]]; then
    echo "✅ Genuine progress detected - resetting loop counters"
    ping_pong_count=0
    same_state_count=0
    error_recovery_entries=0
    last_progress_timestamp="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
fi
```

## State Manager Responsibilities

### STARTUP_CONSULTATION State Integration

The State Manager MUST integrate loop detection into `STARTUP_CONSULTATION/rules.md`:

**Location**: After state transition validation, before accepting the transition

**Required Steps**:
1. Load loop detection counters from state file
2. Check ping-pong pattern (current_state matches last_two_states[0])
3. Check same-state pattern (current_state matches last_two_states[1])
4. Check ERROR_RECOVERY escalation
5. Check for genuine progress (resets counters)
6. Update loop detection state (ring buffer + counters)
7. Enforce ERROR_RECOVERY if thresholds exceeded

### ERROR_RECOVERY Enforcement

When a loop is detected:

```bash
# Force ERROR_RECOVERY state
current_state="ERROR_RECOVERY"
sub_state_machine="null"

# Update state file with loop metadata
jq --arg reason "$loop_reason" \
   --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
   '.state_machine.loop_detection.hard_stop_reason = $reason |
    .state_machine.loop_detection.hard_stop_timestamp = $timestamp' \
   "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"

# Log loop detection
echo "🚨🚨🚨 ERROR_RECOVERY ENFORCED: $loop_reason"
echo "📊 Loop Detection State:"
echo "  - Ping-pong count: $ping_pong_count"
echo "  - Same-state count: $same_state_count"
echo "  - ERROR_RECOVERY entries: $error_recovery_entries"
echo "  - Last two states: ${last_two_states[@]}"
```

### Recovery from ERROR_RECOVERY

Manual intervention required:
1. Investigate root cause of loop (state file, validation logic, etc.)
2. Fix underlying issue
3. Reset loop counters manually:
   ```bash
   jq '.state_machine.loop_detection = {
       "last_two_states": [],
       "ping_pong_count": 0,
       "same_state_count": 0,
       "error_recovery_entries": 0,
       "last_progress_timestamp": null
   }' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json
   ```
4. Manually set appropriate state to resume from
5. Re-run factory continuation

## Performance Characteristics

### Time Complexity: O(1)

- Counter checks: Simple integer comparisons
- Ring buffer update: Fixed 2-element array
- No history parsing, no git operations, no file scanning

### Space Complexity: O(1)

- Fixed 5 fields regardless of execution history
- Ring buffer capped at 2 states
- No growing arrays or unbounded storage

### Comparison to Alternatives

| Approach | Time | Space | Reliability |
|----------|------|-------|-------------|
| **R530 Counters** | O(1) | O(1) | ✅ Bullet-proof |
| Git history parsing | O(N) | O(N) | ❌ Fragile (requires git access) |
| State transition log | O(N) | O(N) | ⚠️ Unbounded growth |
| Manual inspection | O(human) | N/A | ❌ Not automated |

## Integration Points

### 1. Schema Definition

File: `schemas/orchestrator-state-v3.schema.json`

```json
{
  "loop_detection": {
    "type": "object",
    "properties": {
      "last_two_states": {
        "type": "array",
        "items": {"type": "string"},
        "maxItems": 2
      },
      "ping_pong_count": {"type": "integer", "minimum": 0},
      "same_state_count": {"type": "integer", "minimum": 0},
      "error_recovery_entries": {"type": "integer", "minimum": 0},
      "last_progress_timestamp": {
        "type": ["string", "null"],
        "format": "date-time"
      }
    },
    "required": ["last_two_states", "ping_pong_count", "same_state_count", "error_recovery_entries", "last_progress_timestamp"]
  }
}
```

### 2. State Initialization

File: `utilities/init-software-factory-state.sh`

When creating a new state file:

```bash
# Initialize loop detection counters
jq '.state_machine.loop_detection = {
    "last_two_states": [],
    "ping_pong_count": 0,
    "same_state_count": 0,
    "error_recovery_entries": 0,
    "last_progress_timestamp": null
}' "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
```

### 3. Test Fixtures

All test fixtures in `tests/fixtures/orchestrator-state-*.json` MUST include `loop_detection` field to prevent test failures.

## Examples

### Example 1: Ping-Pong Loop Detection

**Scenario**: Factory alternates between CREATE_NEXT_INFRASTRUCTURE and ERROR_RECOVERY

```
Transition 1: CREATE_NEXT_INFRASTRUCTURE (last_two_states: [], ping_pong: 0)
Transition 2: ERROR_RECOVERY             (last_two_states: [SETUP], ping_pong: 0)
Transition 3: CREATE_NEXT_INFRASTRUCTURE (last_two_states: [SETUP, ERROR], ping_pong: 1) ⚠️
Transition 4: ERROR_RECOVERY             (last_two_states: [ERROR, SETUP], ping_pong: 1)
Transition 5: CREATE_NEXT_INFRASTRUCTURE (last_two_states: [SETUP, ERROR], ping_pong: 2) ⚠️
Transition 6: ERROR_RECOVERY             (last_two_states: [ERROR, SETUP], ping_pong: 2)
Transition 7: CREATE_NEXT_INFRASTRUCTURE (last_two_states: [SETUP, ERROR], ping_pong: 3) ⚠️
Transition 8: ERROR_RECOVERY             (last_two_states: [ERROR, SETUP], ping_pong: 3)
Transition 9: CREATE_NEXT_INFRASTRUCTURE (last_two_states: [SETUP, ERROR], ping_pong: 4) ⚠️
Transition 10: ERROR_RECOVERY            (last_two_states: [ERROR, SETUP], ping_pong: 4)
Transition 11: CREATE_NEXT_INFRASTRUCTURE (ping_pong: 5) 🚨 ERROR_RECOVERY ENFORCED
```

**Detection**: On transition 11, `ping_pong_count` reaches 5 → ERROR_RECOVERY

### Example 2: Same-State Loop Detection

**Scenario**: Factory re-enters CREATE_NEXT_INFRASTRUCTURE repeatedly

```
Transition 1: CREATE_NEXT_INFRASTRUCTURE (last_two_states: [], same_state: 0)
Transition 2: CREATE_NEXT_INFRASTRUCTURE (last_two_states: [CREATE], same_state: 1) ⚠️
Transition 3: CREATE_NEXT_INFRASTRUCTURE (last_two_states: [CREATE, CREATE], same_state: 2) ⚠️
Transition 4: CREATE_NEXT_INFRASTRUCTURE (same_state: 3) 🚨 ERROR_RECOVERY ENFORCED
```

**Detection**: On transition 4, `same_state_count` reaches 3 → ERROR_RECOVERY

### Example 3: ERROR_RECOVERY Escalation

**Scenario**: Multiple states fail, triggering ERROR_RECOVERY repeatedly

```
Transition 1: CREATE_NEXT_INFRASTRUCTURE
Transition 2: ERROR_RECOVERY (error_recovery_entries: 1) ⚠️
Transition 3: SPAWN_SW_ENGINEERS
Transition 4: ERROR_RECOVERY (error_recovery_entries: 2) ⚠️
Transition 5: MONITOR
Transition 6: ERROR_RECOVERY (error_recovery_entries: 3) 🚨 ERROR_RECOVERY ENFORCED
```

**Detection**: On transition 6, `error_recovery_entries` reaches 3 → ERROR_RECOVERY

### Example 4: Genuine Progress (Counter Reset)

**Scenario**: Factory makes progress, resets counters

```
Transition 1: CREATE_NEXT_INFRASTRUCTURE (ping_pong: 0)
Transition 2: ERROR_RECOVERY (ping_pong: 0, error_recovery: 1)
Transition 3: CREATE_NEXT_INFRASTRUCTURE (ping_pong: 1)
Transition 4: SPAWN_SW_ENGINEERS ✅ NEW EFFORT CREATED (counters reset!)
Transition 5: MONITOR (ping_pong: 0, error_recovery: 0) ✅ Clean slate
```

**Progress Detected**: New effort created → all counters reset to 0

## Compliance Requirements

### For State Manager Agent

- ✅ MUST load loop detection state on EVERY state transition
- ✅ MUST update counters before accepting transition
- ✅ MUST enforce ERROR_RECOVERY when thresholds exceeded
- ✅ MUST reset counters on genuine progress
- ✅ MUST persist loop detection state to orchestrator-state-v3.json
- ✅ MUST log loop detection events for debugging

### For Factory Initialization

- ✅ MUST initialize loop_detection with zero counters
- ✅ MUST validate loop_detection schema on state file load
- ✅ MUST include loop_detection in all test fixtures

### For Orchestrator Agent

- ✅ MUST consult State Manager before every transition
- ✅ MUST respect ERROR_RECOVERY enforcement (no overrides)
- ✅ MUST document loop detection in state transition logs

## Grading Impact

**Compliance with R530**: +10% (prevents catastrophic failures)
**Loop detected and stopped**: +5% (system health)
**Missing loop detection**: -50% (allows infinite loops)
**Disabled loop detection**: -100% (CRITICAL FAILURE)

---

**Rule Status**: ACTIVE
**Created**: 2025-10-15
**Last Updated**: 2025-10-15
**Related Rules**: R206 (State Machine Validation), R405 (Automation Flag)
